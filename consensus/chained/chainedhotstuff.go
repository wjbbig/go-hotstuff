package chained

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/niclabs/tcrsa"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/config"
	"github.com/wjbbig/go-hotstuff/consensus"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"os"
	"strconv"
	"sync"
)

var logger *logrus.Logger

func init() {
	logger = logging.GetLogger()
}

type ChainedHotStuff struct {
	consensus.HotStuffImpl
	genericQC *pb.QuorumCert
	lockQC    *pb.QuorumCert
	cancel    context.CancelFunc
	lock      sync.Mutex
}

func NewChainedHotStuff(id int, handleMethod func(string) string) *ChainedHotStuff {
	msgEntrance := make(chan *pb.Msg)
	chs := &ChainedHotStuff{}
	chs.MsgEntrance = msgEntrance
	chs.ID = uint32(id)
	chs.View = consensus.NewView(1, 1)
	logger.Debugf("[HOTSTUFF] Init block storage, replica id: %d", id)
	chs.BlockStorage = go_hotstuff.NewBlockStorageImpl(strconv.Itoa(id))
	logger.Debugf("[HOTSTUFF] Generate genesis block")
	genesisBlock := consensus.GenerateGenesisBlock()
	err := chs.BlockStorage.Put(genesisBlock)
	if err != nil {
		logger.Fatal("generate genesis block failed")
	}
	chs.genericQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	chs.lockQC = &pb.QuorumCert{
		BlockHash: genesisBlock.Hash,
		Type:      pb.MsgType_PREPARE_VOTE,
		ViewNum:   0,
		Signature: nil,
	}
	logger.Debugf("[HOTSTUFF] Init command set, replica id: %d", id)
	chs.CmdSet = go_hotstuff.NewCmdSet()

	// read config
	chs.Config = config.HotStuffConfig{}
	chs.Config.ReadConfig()

	// init timer and stop it
	chs.TimeChan = go_hotstuff.NewTimer(chs.Config.Timeout)
	chs.TimeChan.Init()

	chs.BatchTimeChan = go_hotstuff.NewTimer(chs.Config.BatchTimeout)
	chs.BatchTimeChan.Init()

	chs.CurExec = &consensus.CurProposal{
		Node:         nil,
		DocumentHash: nil,
		PrepareVote:  make([]*tcrsa.SigShare, 0),
		HighQC:       make([]*pb.QuorumCert, 0),
	}
	privateKey, err := go_hotstuff.ReadThresholdPrivateKeyFromFile(chs.GetSelfInfo().PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	chs.Config.PrivateKey = privateKey
	chs.ProcessMethod = handleMethod
	ctx, cancel := context.WithCancel(context.Background())
	chs.cancel = cancel
	go chs.receiveMsg(ctx)
	return chs
}

func (chs *ChainedHotStuff) receiveMsg(ctx context.Context) {
	for {
		select {
		case msg := <-chs.MsgEntrance:
			chs.handleMsg(msg)
		case <-chs.BatchTimeChan.Timeout():
			logger.Debug("batch timeout")
			chs.BatchTimeChan.Init()
			chs.batchEvent(chs.CmdSet.GetFirst(int(chs.Config.BatchSize)))
		case <-chs.TimeChan.Timeout():
			chs.Config.Timeout = chs.Config.Timeout * 2
			chs.TimeChan.Init()
			logger.Warn("timeout")
			// create dummy node
			//chs.CreateLeaf()
			// send new view msg
		case <-ctx.Done():
			return
		}
	}
}

func (chs *ChainedHotStuff) handleMsg(msg *pb.Msg) {
	if msg == nil {
		return
	}
	switch msg.Payload.(type) {
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.Debugf("[CHAINED HOTSTUFF] Get request msg, content:%s", request.String())
		// put the cmd into the cmdset
		chs.CmdSet.Add(request.Cmd)
		fmt.Println(chs.GetLeader())
		if chs.GetLeader() != chs.ID {
			// redirect to the leader
			chs.Unicast(chs.GetNetworkInfo()[chs.GetLeader()], msg)
			return
		}
		if chs.CurExec.Node != nil {
			return
		}
		chs.BatchTimeChan.SoftStartTimer()
		// if the length of unprocessed cmd equals to batch size, stop timer and call handleMsg to send prepare msg
		logger.Debugf("Command cache size: %d", len(chs.CmdSet.GetFirst(int(chs.Config.BatchSize))))
		cmds := chs.CmdSet.GetFirst(int(chs.Config.BatchSize))
		if len(cmds) == int(chs.Config.BatchSize) {
			// stop timer
			chs.BatchTimeChan.Stop()
			// create prepare msg
			chs.batchEvent(cmds)
		}
		break
	case *pb.Msg_Prepare:
		logger.Info("[CHAINED HOTSTUFF] Get generic msg")
		if !chs.MatchingMsg(msg, pb.MsgType_PREPARE) {
			logger.Info("[CHAINED HOTSTUFF] Prepare msg not match")
			return
		}
		// get prepare msg
		prepare := msg.GetPrepare()

		if !chs.SafeNode(prepare.CurProposal, prepare.CurProposal.Justify) {
			logger.Warn("[CHAINED HOTSTUFF] Unsafe node")
			return
		}

		err := chs.BlockStorage.Put(prepare.CurProposal)
		if err != nil {
			logger.Error("")
		}
		// add view number and change leader
		chs.View.ViewNum++
		chs.View.Primary = chs.GetLeader()

		marshal, _ := proto.Marshal(msg)
		chs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, chs.Config.PublicKey)
		chs.CurExec.Node = prepare.CurProposal
		partSig, _ := go_hotstuff.TSign(chs.CurExec.DocumentHash, chs.Config.PrivateKey, chs.Config.PublicKey)

		if chs.View.Primary == chs.ID {
			// vote self
			chs.CurExec.PrepareVote = append(chs.CurExec.PrepareVote, partSig)
		} else {
			// send vote msg to the next leader
			partSigBytes, _ := json.Marshal(partSig)
			voteMsg := chs.VoteMsg(pb.MsgType_PREPARE_VOTE, prepare.CurProposal, nil, partSigBytes)
			chs.Unicast(chs.GetNetworkInfo()[chs.GetLeader()], voteMsg)
			chs.CurExec = consensus.NewCurProposal()
		}
		go chs.update(prepare.CurProposal)
		break
	case *pb.Msg_PrepareVote:
		logger.Info("[CHAINED HOTSTUFF] Get generic vote msg")
		if !chs.MatchingMsg(msg, pb.MsgType_PREPARE_VOTE) {
			logger.Warn("[CHAINED HOTSTUFF] Msg not match")
			return
		}
		prepareVote := msg.GetPrepareVote()
		partSig := new(tcrsa.SigShare)
		_ = json.Unmarshal(prepareVote.PartialSig, partSig)
		if err := go_hotstuff.VerifyPartSig(partSig, chs.CurExec.DocumentHash, chs.Config.PublicKey); err != nil {
			logger.Warn("[CHAINED HOTSTUFF GENERIC-VOTE] Partial signature is not correct")
			return
		}
		chs.CurExec.PrepareVote = append(chs.CurExec.PrepareVote, partSig)
		if len(chs.CurExec.PrepareVote) == chs.Config.F*2+1 {
			signature, _ := go_hotstuff.CreateFullSignature(chs.CurExec.DocumentHash, chs.CurExec.PrepareVote, chs.Config.PublicKey)
			qc := chs.QC(pb.MsgType_PREPARE_VOTE, signature, prepareVote.BlockHash)
			chs.genericQC = qc
			chs.CurExec = consensus.NewCurProposal()
			chs.BatchTimeChan.SoftStartTimer()
		}
		break
	case *pb.Msg_NewView:
		break
	default:
		logger.Warn("receive unsupported msg")
	}
}

func (chs *ChainedHotStuff) SafeNode(node *pb.Block, qc *pb.QuorumCert) bool {
	logger.Debug("safety rule ", bytes.Equal(node.ParentHash, chs.genericQC.BlockHash))
	logger.Debug("liveness rule ", qc.ViewNum, chs.genericQC.ViewNum)
	return bytes.Equal(node.ParentHash, chs.genericQC.BlockHash) || //safety rule
		qc.ViewNum > chs.genericQC.ViewNum // liveness rule
}

func (chs *ChainedHotStuff) update(block *pb.Block) {
	chs.lock.Lock()
	defer chs.lock.Unlock()
	if block.Justify == nil {
		return
	}
	// block = b*, block1 = b'', block2 = b', block3 = b
	block1, err := chs.BlockStorage.BlockOf(block.Justify)
	if err == leveldb.ErrNotFound || block1.Committed {
		return
	}
	// start pre-commit phase on block’s parent
	if bytes.Equal(block.ParentHash, block1.Hash) {
		logger.Infof("[CHAINED HOTSTUFF] Start pre-commit phase on block %s", hex.EncodeToString(block1.Hash))
		chs.genericQC = block.Justify
	}

	block2, err := chs.BlockStorage.BlockOf(block1.Justify)
	if err == leveldb.ErrNotFound || block2.Committed {
		return
	}
	// start commit phase on block1’s parent
	if bytes.Equal(block.ParentHash, block1.Hash) && bytes.Equal(block1.ParentHash, block2.Hash) {
		logger.Infof("[CHAINED HOTSTUFF] Start commit phase on block %s", hex.EncodeToString(block2.Hash))
		chs.lockQC = block1.Justify
	}

	block3, err := chs.BlockStorage.BlockOf(block2.Justify)
	if err == leveldb.ErrNotFound || block3.Committed {
		return
	}

	if bytes.Equal(block.ParentHash, block1.Hash) && bytes.Equal(block1.ParentHash, block2.Hash) &&
		bytes.Equal(block2.ParentHash, block3.Hash) {
		//decide
		logger.Infof("[CHAINED HOTSTUFF] Start decide phase on block %s", hex.EncodeToString(block3.Hash))
		chs.processProposal(block3)
	}
}

func (chs *ChainedHotStuff) SafeExit() {
	chs.cancel()
	close(chs.MsgEntrance)
	chs.BlockStorage.Close()
	os.RemoveAll("/opt/hotstuff/dbfile/node" + strconv.Itoa(int(chs.ID)))
}

func (chs *ChainedHotStuff) batchEvent(cmds []string) {
	// if batch timeout, check size
	if len(cmds) == 0 {
		chs.BatchTimeChan.SoftStartTimer()
		return
	}
	// create prepare msg
	if chs.HighQC == nil {
		chs.HighQC = chs.genericQC
	}
	node := chs.CreateLeaf(chs.BlockStorage.GetLastBlockHash(), cmds, chs.HighQC)
	chs.CurExec.Node = node
	err := chs.BlockStorage.Put(node)
	if err != nil {
		logger.Errorf("store block %s failed", hex.EncodeToString(node.Hash))
	}
	chs.CmdSet.MarkProposed(cmds...)
	prepareMsg := chs.Msg(pb.MsgType_PREPARE, node, nil)
	// broadcast prepare msg
	chs.Broadcast(prepareMsg)
	chs.TimeChan.SoftStartTimer()
	go chs.update(node)
	chs.View.ViewNum++
	chs.View.Primary = chs.GetLeader()
}

func (chs *ChainedHotStuff) processProposal(block *pb.Block) {
	// process proposal
	logger.Debugf("cmds: %v", block.Commands)
	go chs.ProcessProposal(block.Commands)
	// store block
	block.Committed = true
}
