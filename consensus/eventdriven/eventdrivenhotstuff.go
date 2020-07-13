package eventdriven

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/niclabs/tcrsa"
	"github.com/syndtr/goleveldb/leveldb"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/config"
	"github.com/wjbbig/go-hotstuff/consensus"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"strconv"
	"sync"
)

var logger = logging.GetLogger()

type Event uint8

const (
	QCFinish Event = iota
	ReceiveProposal
	HighQCUpdate
	ReceiveNewView
)

type EventDrivenHotStuff interface {
	Update(block *pb.Block)
	OnCommit(block *pb.Block)
	OnReceiveProposal(msg *pb.Prepare) (*tcrsa.SigShare, error)
	OnReceiveVote(partSig *tcrsa.SigShare)
	OnPropose() *pb.Block
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
		logger.Fatal("store genesis block failed!")
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
	ehs.View = consensus.NewView(0, 1)
	ehs.qcHigh = ehs.QC(pb.MsgType_PREPARE_VOTE, nil, genesisBlock.Hash)
	ehs.View.ViewNum++
	ehs.waitProposal = sync.NewCond(&ehs.lock)
	msgEntrance := make(chan *pb.Msg)
	ehs.MsgEntrance = msgEntrance
	ehs.ID = uint32(id)
	logger.Debugf("[HOTSTUFF] Init block storage, replica id: %d", id)
	ehs.BlockStorage = blockStore
	logger.Debugf("[HOTSTUFF] Init command set, replica id: %d", id)
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
}

func (ehs *EventDrivenHotStuffImpl) receiveMsg(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ehs.MsgEntrance:
			ehs.handleMsg(msg)
		}
	}
}

func (ehs *EventDrivenHotStuffImpl) handleMsg(msg *pb.Msg) {
	switch msg.Payload.(type) {
	case *pb.Msg_Request:
		request := msg.GetRequest()
		logger.Debugf("[EVENT-DRIVEN HOTSTUFF] Got request msg, content:%s", request.String())
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
		_ = json.Unmarshal(prepareVoteMsg.PartialSig, partSig)
		ehs.OnReceiveVote(partSig)
		break
	case *pb.Msg_NewView:
		newViewMsg := msg.GetNewView()
		// TODO
		ehs.pacemaker.OnReceiverNewView(newViewMsg.PrepareQC)
		break
	default:
		logger.Warn("receive unsupported msg")
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
	if block1 != nil || block1.Committed {
		return
	}

	ehs.lock.Lock()
	defer ehs.lock.Unlock()

	logger.Infof("[HOTSTUFF] PRE COMMIT %s", block1.String())
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
		logger.Infof("[HOTSTUFF] COMMIT %v", block2)
	}

	block3, err := ehs.BlockStorage.BlockOf(block2.Justify)
	if err != nil && err != leveldb.ErrNotFound {
		logger.Fatal(err)
	}
	if block3 == nil || block3.Committed {
		return
	}

	if bytes.Equal(block1.ParentHash, block2.Hash) && bytes.Equal(block2.ParentHash, block3.Hash) {
		logger.Infof("[HOTSTUFF] DECIDE %v", block3)
		ehs.OnCommit(block3)
	}
}

func (ehs *EventDrivenHotStuffImpl) OnCommit(block *pb.Block) {
	if ehs.bExec.Height < block.Height {
		if parent, _ := ehs.BlockStorage.ParentOf(block); parent != nil {
			ehs.OnCommit(parent)
		}
		block.Committed = true
		logger.Infof("[EVENT-DRIVEN HOTSTUFF] EXEC. Block hash: %s", block.Hash)
		ehs.ProcessProposal(block.Commands)
	}
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveProposal(msg *pb.Prepare) (*tcrsa.SigShare, error) {
	newBlock := msg.CurProposal
	logger.Info("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: ", hex.EncodeToString(newBlock.Hash))
	// store the block
	err := ehs.BlockStorage.Put(newBlock)
	if err != nil {
		logger.Errorf("store the new block failed! block hash: %s", hex.EncodeToString(newBlock.Hash))
	}
	ehs.lock.Lock()
	qcBlock, _ := ehs.expectBlock(newBlock.Justify.BlockHash)

	if newBlock.Height <= ehs.vHeight {
		ehs.lock.Unlock()
		logger.Warn("OnReceiveProposal: Block height less than vHeight")
		return nil, errors.New("block was not accepted")
	}
	safe := false

	if qcBlock != nil && qcBlock.Height > ehs.bLock.Height {
		safe = true
	} else {
		logger.Println("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: liveness condition failed")
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
			logger.Println("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: safety condition failed")
		}
	}
	// unsafe, return
	if !safe {
		ehs.lock.Unlock()
		logger.Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: Block not safe")
		return nil, errors.New("block was not accepted")
	}
	logger.Debug("[EVENT-DRIVEN HOTSTUFF] OnReceiveProposal: Accepted block")
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
	partSig, _ := go_hotstuff.TSign(ehs.CurExec.DocumentHash, ehs.Config.PrivateKey, ehs.Config.PublicKey)
	return partSig, nil
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveVote(partSig *tcrsa.SigShare) {
	// verify partSig
	err := go_hotstuff.VerifyPartSig(partSig, ehs.CurExec.DocumentHash, ehs.Config.PublicKey)
	if err != nil {
		logger.Warn("[EVENT-DRIVEN HOTSTUFF] OnReceiveVote: signature not verified!")
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

func (ehs *EventDrivenHotStuffImpl) OnPropose() *pb.Block {
	ehs.BatchTimeChan.SoftStartTimer()
	cmds := ehs.CmdSet.GetFirst(int(ehs.Config.BatchSize))
	if len(cmds) != 0 {
		ehs.BatchTimeChan.Stop()
	} else {
		return nil
	}
	// create node
	proposal := ehs.createProposal(cmds)
	// create a new prepare msg
	msg := ehs.Msg(pb.MsgType_PREPARE, proposal, nil)
	// broadcast
	err := ehs.Broadcast(msg)
	if err != nil {
		logger.Warnf("broadcast proposal failed, error: %s", err.Error())
	}
	return proposal
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
		logger.Errorf("store new block failed!, block hash: %s", hex.EncodeToString(block.Hash))
	}
	return block
}
