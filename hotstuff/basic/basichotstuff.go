package basic

import (
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"strconv"
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
	logger.Debugf("[HOTSTUFF] Generate genesis block")
	genesisBlock := hotstuff.GenerateGenesisBlock()
	err := bhs.BlockStorage.Put(genesisBlock)
	if err != nil {
		logger.Fatal("generate genesis block failed")
	}
	logger.Debugf("[HOTSTUFF] Init command set, replica id: %d", id)
	bhs.CmdSet = go_hotstuff.NewCmdSet()

	// read config
	bhs.Config = hotstuff.HotStuffConfig{}
	bhs.Config.ReadConfig()

	// init timer and stop it
	bhs.TimeChan = go_hotstuff.NewTimer(bhs.Config.Timeout)
	bhs.TimeChan.Init()

	bhs.BatchTimeChan = go_hotstuff.NewTimer(bhs.Config.BatchTimeout)
	bhs.BatchTimeChan.Init()

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
		case <-bhs.TimeChan.Timeout():
		// TODO: timout, goto next view with highQC
		case <-bhs.BatchTimeChan.Timeout():

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
		logger.Debugf("got prepare msg")
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
		bhs.CmdSet.Add(request.Cmd)
		// start batch timer
		bhs.BatchTimeChan.SoftStartTimer()
		// if the length of unprocessed cmd equals to batch size, stop timer and call handleMsg to send prepare msg
		logger.Debugf("cmd set size: %d", len(bhs.CmdSet.GetFirst(int(bhs.Config.BatchSize))))
		cmds := bhs.CmdSet.GetFirst(int(bhs.Config.BatchSize))
		if len(cmds) == int(bhs.Config.BatchSize) {
			// stop timer
			bhs.BatchTimeChan.Stop()
			// create prepare msg
			node := bhs.CreateLeaf(bhs.BlockStorage.GetLastBlockHash(), cmds, nil)
			bhs.CmdSet.MarkProposed(cmds...)
			prepareMsg := &pb.Prepare{
				CurProposal: node,
				HighQC:      nil,
			}
			// vote self

			// send prepare msg
			msg := &pb.Msg{Payload:&pb.Msg_Prepare{prepareMsg}}
			bhs.Broadcast(msg)
		}
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
