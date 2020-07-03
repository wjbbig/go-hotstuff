package basic

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/niclabs/tcrsa"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"strconv"
)

var logger = logging.GetLogger()

type BasicHotStuff struct {
	hotstuff.HotStuffImpl
}

func NewBasicHotStuff(id int) *BasicHotStuff {
	msgEntrance := make(chan *pb.Msg)
	bhs := &BasicHotStuff{}
	bhs.MsgEntrance = msgEntrance
	bhs.PrepareQC = nil
	bhs.PreCommitQC = nil
	bhs.CommitQC = nil
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

	bhs.CurExec = &hotstuff.CurProposal{
		Node:          nil,
		DocumentHash:  nil,
		PrepareVote:   make([]*tcrsa.SigShare, 0),
		PreCommitVote: make([]*tcrsa.SigShare, 0),
		CommitVote:    make([]*tcrsa.SigShare, 0),
	}
	privateKey, err := go_hotstuff.ReadThresholdPrivateKeyFromFile(bhs.GetSelfInfo().PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	bhs.Config.PrivateKey = privateKey
	go bhs.receiveMsg()

	return bhs
}

//func (bhs *BasicHotStuff) SafeNode(node *pb.Block, qc *pb.QuorumCert) bool {
//	//
//}

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
			logger.Warn("Time out, goto new view")
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
		logger.Debugf("[HOTSTUFF PREPARE] Got prepare msg")

		if !bhs.MatchingMsg(msg, pb.MsgType_PREPARE) {
			logger.Debugf("[HOTSTUFF PREPARE] msg does not match")
			return
		}
		prepare := msg.GetPrepare()
		// TODO FIX HighQC is nil
		if bytes.Equal(prepare.CurProposal.ParentHash, prepare.HighQC.BlockHash) &&
			bhs.SafeNode(prepare.CurProposal, prepare.HighQC) {
			logger.Debugf("[HOTSTUFF PREPARE] node is not correct")
			return
		}
		marshal, _ := proto.Marshal(msg)
		bhs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, bhs.Config.PublicKey)
		bhs.CurExec.Node = prepare.CurProposal
		partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		prepareVoteMsg := &pb.Msg{}
		prepareVoteMsg.Payload = &pb.Msg_PrepareVote{
			PrepareVote: &pb.PrepareVote{
				BlockHash:  prepare.CurProposal.Hash,
				PartialSig: partSigBytes,
			},
		}
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.ID], msg)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_PrepareVote:
		logger.Debugf("[HOTSTUFF PREPARE-VOTE] Got prepare vote msg")
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
			msg := &pb.Msg{Payload: &pb.Msg_Prepare{prepareMsg}}
			// vote self
			marshal, _ := proto.Marshal(msg)
			bhs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, bhs.Config.PublicKey)
			partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
			bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
			// broadcast prepare msg
			bhs.Unicast("localhost:8081", msg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	default:
		logger.Warn("Unsupported msg type, drop it.")
		break
	}
}
