package eventdriven

import (
	"context"
	"github.com/wjbbig/go-hotstuff/consensus"
	pb "github.com/wjbbig/go-hotstuff/proto"
)

type Pacemaker interface {
	UpdateHighQC(qcHigh *pb.QuorumCert)
	OnBeat()
	OnNextSyncView()
	OnReceiverNewView(qc *pb.QuorumCert)
	Run(ctx context.Context)
}

type pacemakerImpl struct {
	ehs    *EventDrivenHotStuffImpl
	notify chan Event
}

func NewPacemaker(e *EventDrivenHotStuffImpl) *pacemakerImpl {
	return &pacemakerImpl{e, e.GetEvents()}
}

func (p *pacemakerImpl) UpdateHighQC(qcHigh *pb.QuorumCert) {
	logger.Info("[EVENT-DRIVEN HOTSTUFF] UpdateHighQC")
	block, _ := p.ehs.expectBlock(qcHigh.BlockHash)
	if block == nil {
		logger.Warn("Could not find block of new QC")
		return
	}
	oldQCHighBlock, _ := p.ehs.BlockStorage.BlockOf(p.ehs.qcHigh)
	if oldQCHighBlock == nil {
		logger.Error("Block from the old qcHigh missing from storage")
		return
	}

	if block.Height > oldQCHighBlock.Height {
		p.ehs.qcHigh = qcHigh
		p.ehs.bLeaf = block
		p.ehs.emitEvent(HighQCUpdate)
	}
}

func (p *pacemakerImpl) OnBeat() {
	go p.ehs.OnPropose()
}

func (p *pacemakerImpl) OnNextSyncView() {
	logger.Warn("[EVENT-DRIVEN HOTSTUFF] NewViewTimeout triggered")
	// view change
	p.ehs.View.ViewNum++
	p.ehs.View.Primary = p.ehs.GetLeader()
	// create a dummyNode
	p.ehs.CreateLeaf(p.ehs.GetLeaf().Hash, nil, nil)
	// create a new view msg
	newViewMsg := p.ehs.Msg(pb.MsgType_NEWVIEW, nil, p.ehs.GetHighQC())
	// send msg
	_ = p.ehs.Unicast(p.ehs.GetNetworkInfo()[p.ehs.GetLeader()], newViewMsg)
	// clean the current proposal
	p.ehs.CurExec = consensus.NewCurProposal()
}

func (p *pacemakerImpl) OnReceiverNewView(qc *pb.QuorumCert) {
	p.ehs.lock.Lock()
	defer p.ehs.lock.Unlock()
	logger.Info("[EVENT-DRIVEN HOTSTUFF] OnReceiveNewView")
	p.ehs.emitEvent(ReceiveNewView)
	p.UpdateHighQC(qc)
}

func (p *pacemakerImpl) Run(ctx context.Context) {
	if p.ehs.ID == p.ehs.GetLeader() {
		go p.OnBeat()
	}
	go p.startNewViewTimeout(ctx)
	defer p.ehs.TimeChan.Stop()
	defer p.ehs.BatchTimeChan.Stop()
	// get events
	n := <-p.notify

	for {
		switch n {
		case ReceiveProposal:
			p.ehs.TimeChan.HardStartTimer()
			break
		case QCFinish:
			p.OnBeat()
			break
		case ReceiveNewView:
			p.OnBeat()
			break
		}

		var ok bool
		select {
		case n, ok = <-p.notify:
			if !ok {
				return
			}
			break
		case <-ctx.Done():
			return
		}
	}
}

func (p *pacemakerImpl) startNewViewTimeout(ctx context.Context) {
	for {
		select {
		case <-p.ehs.TimeChan.Timeout():
			logger.Info("dsfdsfdsfsdf")
			// To keep liveness, multiply the timeout duration by 2
			p.ehs.Config.Timeout *= 2
			// init timer
			p.ehs.TimeChan.Init()
			// send new view msg
			//p.OnNextSyncView()
		case <-p.ehs.BatchTimeChan.Timeout():
			logger.Debug("[EVENT-DRIVEN HOTSTUFF] BatchTimeout triggered")
			p.ehs.BatchTimeChan.Init()
			go p.ehs.OnPropose()
		case <-ctx.Done():
			p.ehs.TimeChan.Stop()
			p.ehs.BatchTimeChan.Stop()
		}
	}
}
