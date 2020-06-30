package hotstuff

import (
	"github.com/golang/protobuf/ptypes/duration"
	go_hotstuff "go-hotstuff"
	pb "go-hotstuff/proto"
	"time"
)

// common hotstuff func defined in the paper
type HotStuff interface {
	Msg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg
	VoteMsg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg
	CreateLeaf(parentHash []byte, cmds []string) *pb.Block
	QC(msg *pb.Msg) *pb.QuorumCert
	MatchingMsg(msg *pb.Msg, msgType string, viewNum uint64)
	MatchingQC(qc *pb.QuorumCert, msgType string, viewNum uint64)
	SafeNode(node *pb.Block, qc *pb.QuorumCert)
}

type HotStuffImpl struct {
	HotStuff
	ID           uint32
	BlockStorage go_hotstuff.BlockStorage
	View         *View
	Config       HotStuffConfig
	TimeChan     *time.Timer
	CurExec      bool
	CmdSet       go_hotstuff.CmdSet
	IsPrimary    bool
}

type Replica struct {
	ID         uint32
	Address    string
	publicKey  string
	privateKey string
}

type View struct {
	ViewNum uint64 // view number
	Primary uint32 // the leader's id
}

func NewView(viewNum uint64, primary uint32) *View {
	return &View{
		ViewNum: viewNum,
		Primary: primary,
	}
}

type HotStuffConfig struct {
	networkType    string
	batchSize      int64
	leaderSchedule []uint32
	timeout        duration.Duration
	cluster        []*Replica
}

// ReadConfig reads hotstuff config from yaml file
func (hsc *HotStuffConfig) ReadConfig() *HotStuffConfig {
	// TODO: viper viper viper
	return nil
}

// GenerateGenesisBlock returns genesis block
func GenerateGenesisBlock() *pb.Block {
	genesisBlock := &pb.Block{
		ParentHash: nil,
		Hash:       nil,
		Height:     0,
		Commands:   nil,
		Justify:    nil,
	}
	hash := go_hotstuff.Hash(genesisBlock)
	genesisBlock.Hash = hash
	return genesisBlock
}
