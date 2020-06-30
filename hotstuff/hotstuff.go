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
	viper.AddConfigPath("../")
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

func (h *HotStuffImpl) Msg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
	panic("implement me")
}

func (h *HotStuffImpl) VoteMsg(msgType string, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
	panic("implement me")
}

func (h *HotStuffImpl) CreateLeaf(parentHash []byte, cmds []string) *pb.Block {
	panic("implement me")
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

func (h *HotStuffImpl) SafeNode(node *pb.Block, qc *pb.QuorumCert) {
	panic("implement me")
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
