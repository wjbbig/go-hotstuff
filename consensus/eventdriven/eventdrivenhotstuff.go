package eventdriven

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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

var logger = logging.GetLogger()

type Event uint8

const (
	QCFinish Event = iota
	ReceiveProposal
	ReceiveNewView
)

type EventDrivenHotStuff interface {
	Update(block *pb.Block)
	OnCommit(block *pb.Block)
	OnReceiveProposal(msg *pb.Prepare) (*tcrsa.SigShare, error)
	OnReceiveVote(partSig *tcrsa.SigShare)
	OnPropose()
}

type EventDrivenHotStuffImpl struct {
	consensus.HotStuffImpl
	lock          sync.Mutex
	pacemaker     Pacemaker
	bLeaf         *pb.Block
	bLock         *pb.Block
	bExec         *pb.Block
	qcHigh        *pb.QuorumCert
	vHeight       uint64
	waitProposal  *sync.Cond
	pendingUpdate chan *pb.Block
	cancel        context.CancelFunc
	eventChannels []chan Event
}

func NewEventDrivenHotStuff(id int, handleMethod func(string) string) *EventDrivenHotStuffImpl {
	logger.Debugf("[EVENT-DRIVEN HOTSTUFF] Generate genesis block")
	genesisBlock := consensus.GenerateGenesisBlock()
	blockStore := go_hotstuff.NewBlockStorageImpl(strconv.Itoa(id))
	err := blockStore.Put(genesisBlock)
	if err != nil {
		logger.Fatal("Store genesis block failed!")
	}
	ctx, cancel := context.WithCancel(context.Background())
	ehs := &EventDrivenHotStuffImpl{
		bLeaf:         genesisBlock,
		bLock:         genesisBlock,
		bExec:         genesisBlock,
		qcHigh:        nil,
		vHeight:       genesisBlock.Height,
		pendingUpdate: make(chan *pb.Block, 1),
		cancel:        cancel,
		eventChannels: make([]chan Event, 0),
	}
	// make view number equal to 0 to create genesis block QC
	ehs.View = consensus.NewView(0, 1)
	ehs.qcHigh = ehs.QC(pb.MsgType_PREPARE_VOTE, nil, genesisBlock.Hash)
	// view number add 1
	ehs.View.ViewNum++
	ehs.waitProposal = sync.NewCond(&ehs.lock)
	msgEntrance := make(chan *pb.Msg)
	ehs.MsgEntrance = msgEntrance
	ehs.ID = uint32(id)
	logger.WithField("replicaID", id).Debug("[EVENT-DRIVEN HOTSTUFF] Init block storage.")
	ehs.BlockStorage = blockStore
	logger.WithField("replicaID", id).Debug("[EVENT-DRIVEN HOTSTUFF] Init command cache.")
	ehs.CmdSet = go_hotstuff.NewCmdSet()

	// read config
	ehs.Config = config.HotStuffConfig{}
	ehs.Config.ReadConfig()

	// init timer and stop it
	ehs.TimeChan = go_hotstuff.NewTimer(ehs.Config.Timeout)
	ehs.TimeChan.Init()

	ehs.BatchTimeChan = go_hotstuff.NewTimer(ehs.Config.BatchTimeout)
	ehs.BatchTimeChan.Init()

	ehs.CurExec = &consensus.CurProposal{
		Node:         nil,
		DocumentHash: nil,
		PrepareVote:  make([]*tcrsa.SigShare, 0),
		HighQC:       make([]*pb.QuorumCert, 0),
	}
	privateKey, err := go_hotstuff.ReadThresholdPrivateKeyFromFile(ehs.GetSelfInfo().PrivateKey)
	if err != nil {
		logger.Fatal(err)
	}
	ehs.Config.PrivateKey = privateKey
	ehs.ProcessMethod = handleMethod
	ehs.pacemaker = NewPacemaker(ehs)
	go ehs.updateAsync(ctx)
	go ehs.receiveMsg(ctx)
	go ehs.pacemaker.Run(ctx)
	return ehs
}

func (ehs *EventDrivenHotStuffImpl) emitEvent(event Event) {
	for _, c := range ehs.eventChannels {
		c <- event
	}
}

func (ehs *EventDrivenHotStuffImpl) GetHeight() uint64 {
	return ehs.bLeaf.Height
}

func (ehs *EventDrivenHotStuffImpl) GetVHeight() uint64 {
	return ehs.vHeight
}

func (ehs *EventDrivenHotStuffImpl) GetLeaf() *pb.Block {
	ehs.lock.Lock()
	defer ehs.lock.Unlock()
	return ehs.bLeaf
}

func (ehs *EventDrivenHotStuffImpl) SetLeaf(b *pb.Block) {
	ehs.lock.Lock()
	defer ehs.lock.Unlock()
	ehs.bLeaf = b
}

func (ehs *EventDrivenHotStuffImpl) GetHighQC() *pb.QuorumCert {
	ehs.lock.Lock()
	defer ehs.lock.Unlock()
	return ehs.qcHigh
}

func (ehs *EventDrivenHotStuffImpl) GetEvents() chan Event {
	c := make(chan Event)
	ehs.eventChannels = append(ehs.eventChannels, c)
	return c
}

func (ehs *EventDrivenHotStuffImpl) SafeExit() {
	ehs.cancel()
	ehs.BlockStorage.Close()
	_ = os.RemoveAll("/opt/hotstuff/dbfile/node" + strconv.Itoa(int(ehs.ID)))
}

func (ehs *EventDrivenHotStuffImpl) receiveMsg(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ehs.MsgEntrance:
			go ehs.handleMsg(msg)
		}
	}
}

func (ehs *EventDrivenHotStuffImpl) handleMsg(msg *pb.Msg) {
	switch msg.Payload.(type) {
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.WithField("content", request.String()).Debug("[EVENT-DRIVEN HOTSTUFF] Get request msg.")
		// put the cmd into the cmdset
		ehs.CmdSet.Add(request.Cmd)
		// send the request to the leader, if the replica is not the leader
		if ehs.ID != ehs.GetLeader() {
			_ = ehs.Unicast(ehs.GetNetworkInfo()[ehs.GetLeader()], msg)
			return
		}
		break
	case *pb.Msg_Prepare:
		prepareMsg := msg.GetPrepare()
		partSig, err := ehs.OnReceiveProposal(prepareMsg)
		if err != nil {
			logger.Error(err.Error())
			break
		}
		// view change
		ehs.View.ViewNum++
		ehs.View.Primary = ehs.GetLeader()
		if ehs.View.Primary == ehs.ID {
			// vote self
			ehs.OnReceiveVote(partSig)
		} else {
			// send vote to the leader
			partSigBytes, _ := json.Marshal(partSig)
			voteMsg := ehs.VoteMsg(pb.MsgType_PREPARE_VOTE, prepareMsg.CurProposal, nil, partSigBytes)
			_ = ehs.Unicast(ehs.GetNetworkInfo()[ehs.GetLeader()], voteMsg)
			ehs.CurExec = consensus.NewCurProposal()
		}
		break
	case *pb.Msg_PrepareVote:
		prepareVoteMsg := msg.GetPrepareVote()
		partSig := &tcrsa.SigShare{}
		err := json.Unmarshal(prepareVoteMsg.PartialSig, partSig)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Unmarshal partSig failed.")
		}
		ehs.OnReceiveVote(partSig)
		break
	case *pb.Msg_NewView:
		newViewMsg := msg.GetNewView()
		// wait for 2f votes
		ehs.CurExec.HighQC = append(ehs.CurExec.HighQC, newViewMsg.PrepareQC)
		if len(ehs.CurExec.HighQC) == ehs.Config.F {
			for _, cert := range ehs.CurExec.HighQC {
				if cert.ViewNum > ehs.GetHighQC().ViewNum {
					ehs.qcHigh = cert
				}
			}
			ehs.pacemaker.OnReceiverNewView(ehs.qcHigh)
		}
		break
	default:
		logger.Warn("Receive unsupported msg")
	}
}

// updateAsync receive block
func (ehs *EventDrivenHotStuffImpl) updateAsync(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case b := <-ehs.pendingUpdate:
			ehs.Update(b)
		}
	}
}

// Update update blocks before block
func (ehs *EventDrivenHotStuffImpl) Update(block *pb.Block) {
	// block1 = b'', block2 = b', block3 = b
	block1, err := ehs.BlockStorage.BlockOf(block.Justify)
	if err != nil && err != leveldb.ErrNotFound {
		logger.Fatal(err)
	}
	if block1 == nil || block1.Committed {
		return
	}

	ehs.lock.Lock()
	defer ehs.lock.Unlock()

	logger.WithField("blockHash", hex.EncodeToString(block1.Hash)).Info("[EVENT-DRIVEN HOTSTUFF] PRE COMMIT.")
	// pre-commit block1
	ehs.pacemaker.UpdateHighQC(block.Justify)

	block2, err := ehs.BlockStorage.BlockOf(block1.Justify)
	if err != nil && err != leveldb.ErrNotFound {
		logger.Fatal(err)
	}
	if block2 == nil || block2.Committed {
		return
	}

	if block2.Height > ehs.bLock.Height {
		ehs.bLock = block2
		logger.WithField("blockHash", hex.EncodeToString(block2.Hash)).Info("[EVENT-DRIVEN HOTSTUFF] COMMIT.")
	}

	block3, err := ehs.BlockStorage.BlockOf(block2.Justify)
	if err != nil && err != leveldb.ErrNotFound {
		logger.Fatal(err)
	}
	if block3 == nil || block3.Committed {
		return
	}

	if bytes.Equal(block1.ParentHash, block2.Hash) && bytes.Equal(block2.ParentHash, block3.Hash) {
		logger.WithField("blockHash", hex.EncodeToString(block3.Hash)).Info("[EVENT-DRIVEN HOTSTUFF] DECIDE.")
		ehs.OnCommit(block3)
		ehs.bExec = block3
	}
}

func (ehs *EventDrivenHotStuffImpl) OnCommit(block *pb.Block) {
	if ehs.bExec.Height < block.Height {
		if parent, _ := ehs.BlockStorage.ParentOf(block); parent != nil {
			ehs.OnCommit(parent)
		}
		go func() {
			err := ehs.BlockStorage.UpdateState(block)
			if err != nil {
				logger.WithField("error", err.Error()).Fatal("Update block state failed")
			}
		}()
		logger.WithField("blockHash", hex.EncodeToString(block.Hash)).Info("[EVENT-DRIVEN HOTSTUFF] EXEC.")
		ehs.ProcessProposal(block.Commands)
	}
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveProposal(msg *pb.Prepare) (*tcrsa.SigShare, error) {
	newBlock := msg.CurProposal
	logger.WithField("blockHash", hex.EncodeToString(newBlock.Hash)).Info("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal.")
	// store the block
	err := ehs.BlockStorage.Put(newBlock)
	if err != nil {
		logger.WithField("error", hex.EncodeToString(newBlock.Hash)).Error("Store the new block failed.")
	}
	ehs.lock.Lock()
	qcBlock, _ := ehs.expectBlock(newBlock.Justify.BlockHash)

	if newBlock.Height <= ehs.vHeight {
		ehs.lock.Unlock()
		logger.WithFields(logrus.Fields{
			"blockHeight": newBlock.Height,
			"vHeight":     ehs.vHeight,
		}).Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: Block height less than vHeight.")
		return nil, errors.New("Block was not accepted.")
	}
	safe := false

	if qcBlock != nil && qcBlock.Height > ehs.bLock.Height {
		safe = true
	} else {
		logger.Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: liveness condition failed.")
		b := newBlock
		ok := true
		for ok && b.Height > ehs.bLock.Height+1 {
			b, _ := ehs.BlockStorage.Get(b.ParentHash)
			if b == nil {
				ok = false
			}
		}
		if ok && bytes.Equal(b.ParentHash, ehs.bLock.Hash) {
			safe = true
		} else {
			logger.Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: safety condition failed.")
		}
	}
	// unsafe, return
	if !safe {
		ehs.lock.Unlock()
		logger.Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: Block not safe.")
		return nil, errors.New("Block was not accepted.")
	}
	logger.Debug("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: Accepted block.")
	// update vHeight
	ehs.vHeight = newBlock.Height
	ehs.CmdSet.MarkProposed(newBlock.Commands...)
	ehs.lock.Unlock()
	ehs.waitProposal.Broadcast()
	ehs.emitEvent(ReceiveProposal)
	ehs.pendingUpdate <- newBlock
	marshal, _ := proto.Marshal(msg)
	ehs.CurExec.DocumentHash, _ = go_hotstuff.CreateDocumentHash(marshal, ehs.Config.PublicKey)
	ehs.CurExec.Node = newBlock
	partSig, err := go_hotstuff.TSign(ehs.CurExec.DocumentHash, ehs.Config.PrivateKey, ehs.Config.PublicKey)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: signature not verified!")
	}
	return partSig, nil
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveVote(partSig *tcrsa.SigShare) {
	// verify partSig
	err := go_hotstuff.VerifyPartSig(partSig, ehs.CurExec.DocumentHash, ehs.Config.PublicKey)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":        err.Error(),
			"documentHash": hex.EncodeToString(ehs.CurExec.DocumentHash),
			"Height":       ehs.View.ViewNum,
		}).Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveVote: signature not verified!")
		return
	}
	ehs.CurExec.PrepareVote = append(ehs.CurExec.PrepareVote, partSig)
	if len(ehs.CurExec.PrepareVote) == 2*ehs.Config.F+1 {
		// create full signature
		signature, _ := go_hotstuff.CreateFullSignature(ehs.CurExec.DocumentHash, ehs.CurExec.PrepareVote,
			ehs.Config.PublicKey)
		// create a QC
		qc := ehs.QC(pb.MsgType_PREPARE_VOTE, signature, ehs.CurExec.Node.Hash)
		// update qcHigh
		ehs.pacemaker.UpdateHighQC(qc)
		ehs.CurExec = consensus.NewCurProposal()
		ehs.emitEvent(QCFinish)
	}
}

func (ehs *EventDrivenHotStuffImpl) OnPropose() {
	logger.Info("[EVENT-DRIVEN HOTSTUFF] OnPropose")
	ehs.BatchTimeChan.SoftStartTimer()
	cmds := ehs.CmdSet.GetFirst(int(ehs.Config.BatchSize))
	if len(cmds) != 0 {
		ehs.BatchTimeChan.Stop()
	} else {
		return
	}
	// create node
	proposal := ehs.createProposal(cmds)
	// create a new prepare msg
	msg := ehs.Msg(pb.MsgType_PREPARE, proposal, nil)
	// the old leader should vote too
	ehs.MsgEntrance <- msg
	// broadcast
	err := ehs.Broadcast(msg)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Broadcast proposal failed.")
	}

}

// expectBlock looks for a block with the given Hash, or waits for the next proposal to arrive
// ehs.lock must be locked when calling this function
func (ehs *EventDrivenHotStuffImpl) expectBlock(hash []byte) (*pb.Block, error) {
	block, err := ehs.BlockStorage.Get(hash)
	if err == nil {
		return block, nil
	}
	ehs.waitProposal.Wait()
	return ehs.BlockStorage.Get(hash)
}

// createProposal create a new proposal
func (ehs *EventDrivenHotStuffImpl) createProposal(cmds []string) *pb.Block {
	// create a new block
	ehs.lock.Lock()
	block := ehs.CreateLeaf(ehs.bLeaf.Hash, cmds, ehs.qcHigh)
	ehs.lock.Unlock()
	// store the block
	err := ehs.BlockStorage.Put(block)
	if err != nil {
		logger.WithField("blockHash", hex.EncodeToString(block.Hash)).Error("Store new block failed!")
	}
	return block
}
