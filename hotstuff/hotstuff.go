package hotstuff

import (
	"context"
	"github.com/niclabs/tcrsa"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
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
	ID            uint32
	BlockStorage  go_hotstuff.BlockStorage
	View          *View
	Config        HotStuffConfig
	TimeChan      *go_hotstuff.Timer
	BatchTimeChan *go_hotstuff.Timer
	CurExec       *CurProposal
	CmdSet        go_hotstuff.CmdSet
	IsPrimary     bool
}

type ReplicaInfo struct {
	ID         uint32
	Address    string `mapstructure:"listen-address"`
	PrivateKey string `mapstructure:"privatekeypath"`
}

type CurProposal struct {
	node          *pb.Block
	prepareVote   *tcrsa.SigShare
	preCommitVote *tcrsa.SigShare
	commitVote    *tcrsa.SigShare
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
	PublicKey      *tcrsa.KeyMeta
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
	publicKeyPath := viper.GetString("hotstuff.pubkeypath")
	publicKey, err := go_hotstuff.ReadThresholdPublicKeyFromFile(publicKeyPath)
	if err != nil {
		logger.Fatal(err)
	}
	hsc.PublicKey = publicKey
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

func (h *HotStuffImpl) VoteMsg(msgType int, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
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
	return &pb.QuorumCert{
	}
}

func (h *HotStuffImpl) MatchingMsg(msg *pb.Msg, msgType int, viewNum uint64) {
	panic("implement me")
}

func (h *HotStuffImpl) MatchingQC(qc *pb.QuorumCert, msgType int, viewNum uint64) {
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
			break
		}
	}
	return self
}

func (h *HotStuffImpl) GetNetworkInfo() []*ReplicaInfo {
	networkInfo := make([]*ReplicaInfo, len(h.Config.Cluster)-1)
	for _, info := range h.Config.Cluster {
		if info.ID == h.ID {
			continue
		}
		networkInfo = append(networkInfo, info)
	}
	return networkInfo
}

func (h *HotStuffImpl) Broadcast(msg *pb.Msg) error {
	infos := h.GetNetworkInfo()

	for _, info := range infos {
		err := h.Unicast(info.Address, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HotStuffImpl) Unicast(address string, msg *pb.Msg) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewBasicHotStuffClient(conn)
	_, err = client.SendMsg(context.Background(), msg)
	if err != nil {
		return err
	}
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
