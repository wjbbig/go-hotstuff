package eventdriven

import pb "github.com/wjbbig/go-hotstuff/proto"

type Pacemaker interface {
	UpdateHighQC(qcHigh *pb.QuorumCert)
	OnBeat(cmds []string)
	OnNextSyncView()
	OnReceiverNewView(msg *pb.Msg)
}

type pacemakerImpl struct {
}

func NewPacemaker() *pacemakerImpl {
	return &pacemakerImpl{}
}

func (p *pacemakerImpl) UpdateHighQC(qcHigh *pb.QuorumCert) {
	panic("implement me")
}

func (p *pacemakerImpl) OnBeat(cmds []string) {

	panic("implement me")
}

func (p *pacemakerImpl) OnNextSyncView() {
	panic("implement me")
}

func (p *pacemakerImpl) OnReceiverNewView(msg *pb.Msg) {
	panic("implement me")
}
