package hotstuff

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	go_hotstuff "go-hotstuff"
	"go-hotstuff/logging"
	pb "go-hotstuff/proto"
	"time"
)

var logger = logging.GetLogger()

const (
	PREPARE = iota
	PREPARE_VOTE
	PRECOMMIT
	PRECOMMIT_VOTE
	COMMIT
	COMMIT_TYPE
	DECIDE
)

// common hotstuff func defined in the paper
type HotStuff interface {
	Msg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg
	VoteMsg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg
	CreateLeaf(parentHash []byte, cmds []string, justify *pb.QuorumCert) *pb.Block
	QC(msg *pb.Msg) *pb.QuorumCert
	MatchingMsg(msg *pb.Msg, msgType string, viewNum uint64)
	MatchingQC(qc *pb.QuorumCert, msgType string, viewNum uint64)
	SafeNode(node *pb.Block, qc *pb.QuorumCert)
}

type HotStuffImpl struct {
	ID           uint32
	BlockStorage go_hotstuff.BlockStorage
	View         *View
	Config       HotStuffConfig
	TimeChan     *time.Timer
	CurExec      bool
	CmdSet       go_hotstuff.CmdSet
	IsPrimary    bool
}

type ReplicaInfo struct {
	ID         uint32
	Address    string `mapstructure:"listen-address"`
	PublicKey  string `mapstructure:"pubkeypath"`
	PrivateKey string `mapstructure:"privatekeypath"`
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
	NetworkType    string `mapstructure:"type"`
	BatchSize      uint64 `mapstructure:"batch-size"`
	LeaderSchedule []int  `mapstructure:"leader-schedule"`
	BatchTimeout   time.Duration
	Timeout        time.Duration
	Cluster        []*ReplicaInfo
}

// ReadConfig reads hotstuff config from yaml file
func (hsc *HotStuffConfig) ReadConfig() {
	logger.Debug("[HOTSTUFF] Read config")
	viper.AddConfigPath("/opt/hotstuff/config/")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./")
	viper.SetConfigName("hotstuff")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}
	networkType := viper.GetString("hotstuff.type")
	logger.Debugf("networkType = %s", networkType)
	hsc.NetworkType = networkType
	batchTimeout := viper.GetDuration("hotstuff.batchtimeout")
	logger.Debugf("batchtimeout = %v", batchTimeout)
	hsc.BatchTimeout = batchTimeout
	timeout := viper.GetDuration("hotstuff.timeout")
	logger.Debugf("timeout = %v", timeout)
	hsc.Timeout = timeout
	batchSize := viper.GetUint64("hotstuff.batch-size")
	logrus.Debugf("batch size = %d", batchSize)
	hsc.BatchSize = batchSize
	leaderSchedule := viper.GetIntSlice("hotstuff.leader-schedule")
	logrus.Debugf("leader schedule = %v", leaderSchedule)
	hsc.LeaderSchedule = leaderSchedule
	err = viper.UnmarshalKey("hotstuff.cluster", &hsc.Cluster)
	if err != nil {
		logger.Fatal(err)
	}
}

func (h *HotStuffImpl) Msg(msgType int, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
	msg := &pb.Msg{}
	switch msgType {
	case PREPARE:
	}
	return msg
}

func (h *HotStuffImpl) VoteMsg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
	msg := &pb.Msg{}
	return msg
}

func (h *HotStuffImpl) CreateLeaf(parentHash []byte, cmds []string, justify *pb.QuorumCert) *pb.Block {
	b := &pb.Block{
		ParentHash: parentHash,
		Hash:       nil,
		Height:     h.View.ViewNum,
		Commands:   cmds,
		Justify:    justify,
	}

	b.Hash = go_hotstuff.Hash(b)
	return b
}

func (h *HotStuffImpl) QC(msg *pb.Msg) *pb.QuorumCert {
	panic("implement me")
}

func (h *HotStuffImpl) MatchingMsg(msg *pb.Msg, msgType string, viewNum uint64) {
	panic("implement me")
}

func (h *HotStuffImpl) MatchingQC(qc *pb.QuorumCert, msgType string, viewNum uint64) {
	panic("implement me")
}

func (h *HotStuffImpl) SafeNode(node *pb.Block, qc *pb.QuorumCert) {}

// GetLeader get the leader replica in view
func (h *HotStuffImpl) GetLeader() uint32 {
	id := h.View.ViewNum % uint64(h.View.Primary)
	return uint32(id)
}

func (h *HotStuffImpl) GetSelfInfo() *ReplicaInfo {
	self := &ReplicaInfo{}
	for _, info := range h.Config.Cluster {
		if info.ID == h.ID {
			self = info
		}
		break
	}
	return self
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
