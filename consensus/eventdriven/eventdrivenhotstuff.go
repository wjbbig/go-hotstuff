package eventdriven

import (
	"github.com/wjbbig/go-hotstuff/consensus"
	"github.com/wjbbig/go-hotstuff/consensus/eventdriven/pacemaker"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
)

var logger = logging.GetLogger()

type EventDrivenHotStuff interface {
	Update()
	OnCommit()
	OnReceiveProposal()
	OnReceiveVote()
	OnPropose()
}


type eventDrivenHotStuffImpl struct {
	consensus.HotStuff
	pacemaker *pacemaker.Pacemaker
	bLeaf *pb.Block
	bLock *pb.Block
	bExec *pb.Block
	qcHigh *pb.QuorumCert
	vHeight uint64
}
