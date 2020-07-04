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
		NewViewQC:     make([]*pb.QuorumCert, 0),
	}
	privateKey, err := go_hotstuff.ReadThresholdPrivateKeyFromFile(bhs.GetSelfInfo().PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	bhs.Config.PrivateKey = privateKey
	go bhs.receiveMsg()

	return bhs
}

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
		logger.Debug("[HOTSTUFF NEWVIEW] Got new view msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_DECIDE) {
			logger.Debug("[HOTSTUFF NEWVIEW] Msg not match")
			return
		}
		// process highqc and node
		bhs.CurExec.NewViewQC = append(bhs.CurExec.NewViewQC, msg.GetNewView().PrepareQC)

		// TODO next view

		break
	case *pb.Msg_Prepare:
		logger.Debug("[HOTSTUFF PREPARE] Got prepare msg")

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
			preCommitMsg := bhs.Msg(pb.MsgType_PRECOMMIT, bhs.CurExec.Node, qc)
			// broadcast msg
			bhs.Broadcast(preCommitMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_PreCommit:
		logger.Debug("[HOTSTUFF PRECOMMIT] Got precommit msg")
		if !bhs.MatchingQC(msg.GetPreCommit().PrepareQC, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[HOTSTUFF PRECOMMIT] QC not match")
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
		logger.Debug("[HOTSTUFF COMMIT] Got commit msg")
		commit := msg.GetCommit()
		if !bhs.MatchingQC(commit.PreCommitQC, pb.MsgType_PRECOMMIT_VOTE) {
			logger.Warn("[HOTSTUFF COMMIT] QC not match")
			return
		}
		bhs.PreCommitQC = commit.PreCommitQC
		partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
		partSigBytes, _ := json.Marshal(partSig)
		commitVoteMsg := bhs.VoteMsg(pb.MsgType_COMMIT_VOTE, bhs.CurExec.Node, nil, partSigBytes)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], commitVoteMsg)
		bhs.TimeChan.SoftStartTimer()
		break
	case *pb.Msg_CommitVote:
		logger.Debug("[HOTSTUFF COMMIT-VOTE] Got commit vote msg")
		if !bhs.MatchingMsg(msg, pb.MsgType_COMMIT_VOTE) {
			logger.Warn("[HOTSTUFF COMMIT-VOTE] Msg not match")
			return
		}
		commitVoteMsg := msg.GetCommitVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(commitVoteMsg.PartialSig, partSig)
		if err := go_hotstuff.VerifyPartSig(partSig, bhs.CurExec.DocumentHash, bhs.Config.PublicKey); err != nil {
			logger.Warn("[HOTSTUFF COMMIT-VOTE] Partial signature is not correct")
			return
		}
		bhs.CurExec.CommitVote = append(bhs.CurExec.CommitVote, partSig)
		if len(bhs.CurExec.CommitVote) == 2*bhs.Config.F+1 {
			signature, _ := go_hotstuff.CreateFullSignature(bhs.CurExec.DocumentHash, bhs.CurExec.CommitVote, bhs.Config.PublicKey)
			commitQC := bhs.QC(pb.MsgType_COMMIT_VOTE, signature, bhs.CurExec.Node.Hash)
			// vote self
			bhs.CommitQC = commitQC
			decideMsg := bhs.Msg(pb.MsgType_DECIDE, bhs.CurExec.Node, commitQC)
			bhs.Broadcast(decideMsg)

			bhs.TimeChan.Stop()

			bhs.processProposal()
		}

		break
	case *pb.Msg_Decide:
		logger.Debug("[HOTSTUFF DECIDE] Got decide msg")
		decideMsg := msg.GetDecide()
		if !bhs.MatchingQC(decideMsg.CommitQC, pb.MsgType_COMMIT_VOTE) {
			logger.Warn("[HOTSTUFF DECIDE] QC not match")
			return
		}
		bhs.CommitQC = decideMsg.CommitQC
		bhs.TimeChan.Stop()
		bhs.processProposal()
		break
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.Debugf("[HOTSTUFF] Got request msg, content:%s", request.String())
		// put the cmd into the cmdset
		bhs.CmdSet.Add(request.Cmd)
		// send request to the leader, if the replica is not the leader
		if bhs.ID != bhs.GetLeader() {
			bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], msg)
			return
		}
		if bhs.CurExec.Node != nil {
			return
		}
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
			bhs.CurExec.Node = node
			bhs.CmdSet.MarkProposed(cmds...)
			prepareMsg := bhs.Msg(pb.MsgType_PREPARE, node, bhs.PrepareQC)
			// vote self
			marshal, _ := proto.Marshal(prepareMsg)
			bhs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, bhs.Config.PublicKey)
			partSig, _ := go_hotstuff.TSign(bhs.CurExec.DocumentHash, bhs.Config.PrivateKey, bhs.Config.PublicKey)
			bhs.CurExec.PrepareVote = append(bhs.CurExec.PrepareVote, partSig)
			// broadcast prepare msg
			bhs.Broadcast(prepareMsg)
			bhs.TimeChan.SoftStartTimer()
		}
		break
	default:
		logger.Warn("Unsupported msg type, drop it.")
		break
	}
}

func (bhs *BasicHotStuff) processProposal() {
	// process proposal
	bhs.ProcessProposal(bhs.CurExec.Node.Commands)
	// store block
	bhs.BlockStorage.Put(bhs.CurExec.Node)
	// clear curExec
	bhs.CurExec = &hotstuff.CurProposal{
		Node:          nil,
		DocumentHash:  nil,
		PrepareVote:   make([]*tcrsa.SigShare, 0),
		PreCommitVote: make([]*tcrsa.SigShare, 0),
		CommitVote:    make([]*tcrsa.SigShare, 0),
		NewViewQC:     make([]*pb.QuorumCert, 0),
	}
	// add view number
	bhs.View.ViewNum++
	bhs.View.Primary = bhs.GetLeader()
	// check if self is the next leader
	if bhs.GetLeader() != bhs.ID {
		// if not, send next view mag to the next leader
		newViewMsg := bhs.Msg(pb.MsgType_NEWVIEW, nil, bhs.PrepareQC)
		bhs.Unicast(bhs.GetNetworkInfo()[bhs.GetLeader()], newViewMsg)
	}
}
