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
	bhs.PrepareQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	bhs.PreCommitQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PRECOMMIT_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	bhs.CommitQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_COMMIT_VOTE,
		ViewNum:   0,
		Signature: nil,
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
				go bhs.handleMsg(msg)
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

		if !bytes.Equal(prepare.CurProposal.ParentHash, prepare.HighQC.BlockHash) ||
			!bhs.SafeNode(prepare.CurProposal, prepare.HighQC) {
			logger.Debugf("[HOTSTUFF PREPARE] node is not correct")
			return
		}
		// create prepare vote msg
		marshal, _ := proto.Marshal(msg)
		bhs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, bhs.Config.PublicKey)
		bhs.CurExec.Node = prepare.CurProposal
		partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		prepareVoteMsg := bhs.VoteMsg(pb.MsgType_PREPARE_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		// send msg to leader
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], prepareVoteMsg)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_PrepareVote:
		logger.Debugf("[HOTSTUFF PREPARE-VOTE] Got prepare vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[HOTSTUFF PREPARE-VOTE] Msg not match")
			return
		}
		// verify
		prepareVote := msg.GetPrepareVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(prepareVote.PartialSig, partSig)
		if err := go_hotstuff.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF PREPARE-VOTE] Partial signature is not correct")
			return
		}
		// put it into preparevote slice
		bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
		if len(bhs.CurExec.PrepareVote) == bhs.Config.F*2+1 {
			// create full signature
			signature, _ := go_hotstuff.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.PrepareVote, bhs.Config.PublicKey)
			qc := bhs.QC(pb.MsgType_PREPARE_VOTE, signature, prepareVote.BlockHash)
			bhs.PrepareQC = qc
			preCommitMsg := &pb.Msg{}
			preCommitMsg.Payload = &pb.Msg_PreCommit{PreCommit: &pb.PreCommit{
				PrepareQC: qc,
				ViewNum:   bhs.View.ViewNum,
			}}
			// broadcast msg
			bhs.Broadcast(preCommitMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_PreCommit:
		logger.Debug("[HOTSTUFF PRECOMMIT] Got precommit msg")
		if !bhs.MatchingQC(msg.GetPreCommit().PrepareQC, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[HOTSTUFF PRECOMMIT] Msg not match")
			return
		}
		bhs.PrepareQC = msg.GetPreCommit().PrepareQC
		partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		preCommitVote := bhs.VoteMsg(pb.MsgType_PRECOMMIT_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], preCommitVote)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_PreCommitVote:
		logger.Debug("[HOTSTUFF PRECOMMIT-VOTE] Got precommit vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_PRECOMMIT_VOTE) {
			logger.Warn("[HOTSTUFF PRECOMMIT-VOTE] Msg not match")
			return
		}
		// verify
		preCommitVote := msg.GetPreCommitVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(preCommitVote.PartialSig, partSig)
		if err := go_hotstuff.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF PRECOMMIT-VOTE] Partial signature is not correct")
			return
		}
		bhs.CurExec.PreCommitVote = append(bhs.CurExec.PreCommitVote, partSig)
		if len(bhs.CurExec.PreCommitVote) == 2*bhs.Config.F+1 {
			signature, _ := go_hotstuff.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.PreCommitVote, bhs.Config.PublicKey)
			preCommitQC := bhs.QC(pb.MsgType_PRECOMMIT_VOTE, signature, bhs.CurExec.Node.Hash)
			// vote self
			bhs.PreCommitQC = preCommitQC
			commitMsg := bhs.Msg(pb.MsgType_COMMIT, bhs.CurExec.Node, preCommitQC)
			bhs.Broadcast(commitMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_Commit:
		logger.Debug("[HOTSTUFF PRECOMMIT-VOTE] Got precommit vote msg")
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
				HighQC:      bhs.PrepareQC,
				ViewNum:     bhs.View.ViewNum,
			}
			msg := &pb.Msg{Payload: &pb.Msg_Prepare{prepareMsg}}
			// vote self
			marshal, _ := proto.Marshal(msg)
			bhs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, bhs.Config.PublicKey)
			partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
			bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
			// broadcast prepare msg
			bhs.Broadcast(msg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	default:
		logger.Warn("Unsupported msg type, drop it.")
		break
	}
}