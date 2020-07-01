package basic

import (
	go_hotstuff "go-hotstuff"
	"go-hotstuff/hotstuff"
	"go-hotstuff/logging"
	pb "go-hotstuff/proto"
	"strconv"
	"time"
)

var logger = logging.GetLogger()

type BasicHotStuff struct {
	hotstuff.HotStuffImpl
	MsgEntrance chan *pb.Msg   // receive msg
	prepareQC   *pb.QuorumCert // highQC
	preCommitQC *pb.QuorumCert // lockQC
	commitQC    *pb.QuorumCert
}

func NewBasicHotStuff(id int) *BasicHotStuff {
	msgEntrance := make(chan *pb.Msg)
	bhs := &BasicHotStuff{
		MsgEntrance: msgEntrance,
		prepareQC:   nil,
		preCommitQC: nil,
		commitQC:    nil,
	}
	bhs.ID = uint32(id)
	bhs.View = hotstuff.NewView(1, 1)
	logger.Debugf("[HOTSTUFF] Init block storage, replica id: %d", id)
	bhs.BlockStorage = go_hotstuff.NewBlockStorageImpl(strconv.Itoa(id))
	logger.Debugf("[HOTSTUFF] Init command set, replica id: %d", id)
	bhs.CmdSet = go_hotstuff.NewCmdSet()

	// read config
	bhs.Config = hotstuff.HotStuffConfig{}
	bhs.Config.ReadConfig()

	// init timer and stop it
	bhs.TimeChan = time.NewTimer(bhs.Config.Timeout)
	bhs.TimeChan.Stop()
	go bhs.receiveMsg()

	return bhs
}

func (bhs *BasicHotStuff) SafeNode(node *pb.Block, qc *pb.QuorumCert) {
	//
}

// receiveMsg receive msg from msg channel
func (bhs *BasicHotStuff) receiveMsg() {
	for {
		select {
		case msg, ok := <-bhs.MsgEntrance:
			if ok {
				bhs.handleMsg(msg)
			}
		case <-bhs.TimeChan.C:
			// TODO: timout, goto next view with highQC
		}
	}
}

// handleMsg handle different msg with different way
func (bhs *BasicHotStuff) handleMsg(msg *pb.Msg) {
	switch msg.Payload.(type) {
	case *pb.Msg_NewView:
		// check leader

		// process highqc and node

		// broadcast
		break
	case *pb.Msg_Prepare:
		break
	case *pb.Msg_PrepareVote:
		break
	case *pb.Msg_PreCommit:
		break
	case *pb.Msg_PreCommitVote:
		break
	case *pb.Msg_Commit:
		break
	case *pb.Msg_CommitVote:
		break
	case *pb.Msg_Decide:
		break
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.Debugf("[HOTSTUFF] Got request msg, content:%s", request.String())
		// put the cmd into the cmdset
		break
	default:
		logger.Warn("Unsupported msg type, drop it.")
		break
	}
}

func (bhs *BasicHotStuff) Close() {
	close(bhs.MsgEntrance)
	bhs.BlockStorage.Close()
}
