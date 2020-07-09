package eventdriven

import (
	"bytes"
	"context"
	"github.com/niclabs/tcrsa"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/config"
	"github.com/wjbbig/go-hotstuff/consensus"
	"github.com/wjbbig/go-hotstuff/consensus/eventdriven/pacemaker"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"strconv"
	"sync"
)

var logger = logging.GetLogger()

type EventDrivenHotStuff interface {
	Update(block *pb.Block)
	OnCommit(block *pb.Block)
	OnReceiveProposal(msg *pb.Msg)
	OnReceiveVote(msg *pb.Msg)
	OnPropose(bLeaf *pb.Block, cmds []string, qcHigh *pb.QuorumCert) *pb.Block
}

type EventDrivenHotStuffImpl struct {
	consensus.HotStuffImpl
	lock          sync.Mutex
	pacemaker     pacemaker.Pacemaker
	bLeaf         *pb.Block
	bLock         *pb.Block
	bExec         *pb.Block
	qcHigh        *pb.QuorumCert
	vHeight       uint64
	waitProposal  *sync.Cond
	pendingUpdate chan *pb.Block
	cancel        context.CancelFunc
}

func NewEventDrivenHotStuff(id int, handleMethod func(string) string) *EventDrivenHotStuffImpl {
	logger.Debugf("[HOTSTUFF] Generate genesis block")
	genesisBlock := consensus.GenerateGenesisBlock()
	blockStore := go_hotstuff.NewBlockStorageImpl(strconv.Itoa(id))
	err := blockStore.Put(genesisBlock)
	if err != nil {
		logger.Fatal("store genesis block failed!")
	}
	ctx, cancel := context.WithCancel(context.Background())
	ehs := &EventDrivenHotStuffImpl{
		pacemaker:     pacemaker.NewPacemaker(),
		bLeaf:         genesisBlock,
		bLock:         genesisBlock,
		bExec:         genesisBlock,
		qcHigh:        nil,
		vHeight:       genesisBlock.Height,
		pendingUpdate: make(chan *pb.Block, 1),
		cancel:        cancel,
	}
	ehs.waitProposal = sync.NewCond(&ehs.lock)
	msgEntrance := make(chan *pb.Msg)
	ehs.MsgEntrance = msgEntrance
	ehs.ID = uint32(id)
	ehs.View = consensus.NewView(1, 1)
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
	go ehs.updateAsync(ctx)
	go ehs.receiveMsg(ctx)
	return ehs
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
	case *pb.Msg_Prepare:
		break
	case *pb.Msg_NewView:
		ehs.pacemaker.OnReceiverNewView(msg)
		break
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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		logger.Fatal(err)
	}
	if block3 == nil || block3.Committed {
		return
	}

	if bytes.Equal(block1.ParentHash, block2.Hash) && bytes.Equal(block2.ParentHash, block3.Hash) {
		logger.Infof("[HOTSTUFF] DECIDE %v", block3)
		ehs.OnCommit(block3)
		ehs.ProcessProposal(block3.Commands)
	}
}

func (ehs *EventDrivenHotStuffImpl) OnCommit(block *pb.Block) {
	panic("implement me")
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveProposal(msg *pb.Msg) {
	panic("implement me")
}

func (ehs *EventDrivenHotStuffImpl) OnReceiveVote(msg *pb.Msg) {
	panic("implement me")
}

func (ehs *EventDrivenHotStuffImpl) OnPropose(bLeaf *pb.Block, cmds []string, qcHigh *pb.QuorumCert) *pb.Block {
	// create node
	block := ehs.CreateLeaf(bLeaf.Hash, cmds, qcHigh)
	// broadcast

	return block
}
