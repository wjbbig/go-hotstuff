package hotstuff

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/niclabs/tcrsa"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"time"
)

var logger = logging.GetLogger()

// common hotstuff func defined in the paper
type HotStuff interface {
	Msg(msgType pb.MsgType, node *pb.Block, qc *pb.QuorumCert) *pb.Msg
	VoteMsg(msgType pb.MsgType, node *pb.Block, qc *pb.QuorumCert, justify []byte) *pb.Msg
	CreateLeaf(parentHash []byte, cmds []string, justify *pb.QuorumCert) *pb.Block
	QC(msgType pb.MsgType, sig tcrsa.Signature, blockHash []byte) *pb.QuorumCert
	MatchingMsg(msg *pb.Msg, msgType pb.MsgType) bool
	MatchingQC(qc *pb.QuorumCert, msgType pb.MsgType) bool
	SafeNode(node *pb.Block, qc *pb.QuorumCert) bool
	GetMsgEntrance() chan<- *pb.Msg
	GetSelfInfo() *ReplicaInfo
	SafeExit()
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
	HighQC        *pb.QuorumCert
	PrepareQC     *pb.QuorumCert // highQC
	PreCommitQC   *pb.QuorumCert // lockQC
	CommitQC      *pb.QuorumCert
	MsgEntrance   chan *pb.Msg // receive msg
	ProcessMethod func(args string) string
}

type ReplicaInfo struct {
	ID         uint32
	Address    string `mapstructure:"listen-address"`
	PrivateKey string `mapstructure:"privatekeypath"`
}

type CurProposal struct {
	Node          *pb.Block
	DocumentHash  []byte
	PrepareVote   []*tcrsa.SigShare
	PreCommitVote []*tcrsa.SigShare
	CommitVote    []*tcrsa.SigShare
	HighQC        []*pb.QuorumCert
}

func NewCurProposal() *CurProposal {
	return &CurProposal{
		Node:          nil,
		DocumentHash:  nil,
		PrepareVote:   make([]*tcrsa.SigShare, 0),
		PreCommitVote: make([]*tcrsa.SigShare, 0),
		CommitVote:    make([]*tcrsa.SigShare, 0),
		HighQC:        make([]*pb.QuorumCert, 0),
	}
}

type View struct {
	ViewNum      uint64 // view number
	Primary      uint32 // the leader's id
	ViewChanging bool
}

func NewView(viewNum uint64, primary uint32) *View {
	return &View{
		ViewNum:      viewNum,
		Primary:      primary,
		ViewChanging: false,
	}
}

type HotStuffConfig struct {
	NetworkType    string `mapstructure:"type"`
	BatchSize      uint64 `mapstructure:"batch-size"`
	LeaderSchedule []int  `mapstructure:"leader-schedule"`
	BatchTimeout   time.Duration
	Timeout        time.Duration
	PublicKey      *tcrsa.KeyMeta
	PrivateKey     *tcrsa.KeyShare
	Cluster        []*ReplicaInfo
	N              int
	F              int
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
	hsc.N = viper.GetInt("hotstuff.N")
	hsc.F = viper.GetInt("hotstuff.f")
}

func (h *HotStuffImpl) Msg(msgType pb.MsgType, node *pb.Block, qc *pb.QuorumCert) *pb.Msg {
	msg := &pb.Msg{}
	switch msgType {
	case pb.MsgType_PREPARE:
		msg.Payload = &pb.Msg_Prepare{Prepare: &pb.Prepare{
			CurProposal: node,
			HighQC:      qc,
			ViewNum:     h.View.ViewNum,
		}}
		break
	case pb.MsgType_PRECOMMIT:
		msg.Payload = &pb.Msg_PreCommit{PreCommit: &pb.PreCommit{PrepareQC: qc, ViewNum: h.View.ViewNum}}
		break
	case pb.MsgType_COMMIT:
		msg.Payload = &pb.Msg_Commit{Commit: &pb.Commit{PreCommitQC: qc, ViewNum: h.View.ViewNum}}
		break
	case pb.MsgType_NEWVIEW:
		msg.Payload = &pb.Msg_NewView{NewView: &pb.NewView{PrepareQC: qc, ViewNum: h.View.ViewNum}}
		break
	case pb.MsgType_DECIDE:
		msg.Payload = &pb.Msg_Decide{Decide: &pb.Decide{CommitQC: qc, ViewNum: h.View.ViewNum}}
		break
	}
	return msg
}

func (h *HotStuffImpl) VoteMsg(msgType pb.MsgType, node *pb.Block, qc *pb.QuorumCert, justify []byte) *pb.Msg {
	msg := &pb.Msg{}
	switch msgType {
	case pb.MsgType_PREPARE_VOTE:
		msg.Payload = &pb.Msg_PrepareVote{PrepareVote: &pb.PrepareVote{
			BlockHash:  node.Hash,
			Qc:         qc,
			PartialSig: justify,
			ViewNum:    h.View.ViewNum,
		}}
		break
	case pb.MsgType_PRECOMMIT_VOTE:
		msg.Payload = &pb.Msg_PreCommitVote{PreCommitVote: &pb.PreCommitVote{
			BlockHash:  node.Hash,
			Qc:         qc,
			PartialSig: justify,
			ViewNum:    h.View.ViewNum,
		}}
		break
	case pb.MsgType_COMMIT_VOTE:
		msg.Payload = &pb.Msg_CommitVote{CommitVote: &pb.CommitVote{
			BlockHash:  node.Hash,
			Qc:         qc,
			PartialSig: justify,
			ViewNum:    h.View.ViewNum,
		}}
		break
	}
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

func (h *HotStuffImpl) QC(msgType pb.MsgType, sig tcrsa.Signature, blockHash []byte) *pb.QuorumCert {
	marshal, _ := json.Marshal(sig)
	return &pb.QuorumCert{
		BlockHash: blockHash,
		Type:      msgType,
		ViewNum:   h.View.ViewNum,
		Signature: marshal,
	}
}

func (h *HotStuffImpl) MatchingMsg(msg *pb.Msg, msgType pb.MsgType) bool {
	switch msgType {
	case pb.MsgType_PREPARE:
		return msg.GetPrepare() != nil && msg.GetPrepare().ViewNum == h.View.ViewNum
	case pb.MsgType_PREPARE_VOTE:
		return msg.GetPrepareVote() != nil && msg.GetPrepareVote().ViewNum == h.View.ViewNum
	case pb.MsgType_PRECOMMIT:
		return msg.GetPreCommit() != nil && msg.GetPreCommit().ViewNum == h.View.ViewNum
	case pb.MsgType_PRECOMMIT_VOTE:
		return msg.GetPreCommitVote() != nil && msg.GetPreCommitVote().ViewNum == h.View.ViewNum
	case pb.MsgType_COMMIT:
		return msg.GetCommit() != nil && msg.GetCommit().ViewNum == h.View.ViewNum
	case pb.MsgType_COMMIT_VOTE:
		return msg.GetCommitVote() != nil && msg.GetCommitVote().ViewNum == h.View.ViewNum
	case pb.MsgType_NEWVIEW:
		return msg.GetNewView() != nil && msg.GetNewView().ViewNum == h.View.ViewNum
	}
	return false
}

func (h *HotStuffImpl) MatchingQC(qc *pb.QuorumCert, msgType pb.MsgType) bool {
	return qc.Type == msgType && qc.ViewNum == h.View.ViewNum
}

func (h *HotStuffImpl) SafeNode(node *pb.Block, qc *pb.QuorumCert) bool {
	return bytes.Equal(node.ParentHash, h.PreCommitQC.BlockHash) || //safety rule
		qc.ViewNum > h.PreCommitQC.ViewNum // liveness rule
}

func (h *HotStuffImpl) GetMsgEntrance() chan<- *pb.Msg {
	return h.MsgEntrance
}

func (h *HotStuffImpl) SafeExit() {
	close(h.MsgEntrance)
	h.BlockStorage.Close()
	os.RemoveAll("/opt/hotstuff/dbfile/node" + strconv.Itoa(int(h.ID)))
}

// GetLeader get the leader replica in view
func (h *HotStuffImpl) GetLeader() uint32 {
	id := h.View.ViewNum % uint64(h.Config.N)
	if id == 0 {
		id = 4
	}
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

func (h *HotStuffImpl) GetNetworkInfo() map[uint32]string {
	networkInfo := make(map[uint32]string)
	for _, info := range h.Config.Cluster {
		if info.ID == h.ID {
			continue
		}
		networkInfo[info.ID] = info.Address
	}
	return networkInfo
}

func (h *HotStuffImpl) Broadcast(msg *pb.Msg) error {
	infos := h.GetNetworkInfo()

	for _, address := range infos {
		err := h.Unicast(address, msg)
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

func (h *HotStuffImpl) ProcessProposal(cmds []string) {
	for _, cmd := range cmds {
		result := h.ProcessMethod(cmd)
		msg := &pb.Msg{Payload: &pb.Msg_Reply{Reply: &pb.Reply{Result: result}}}
		h.Unicast("localhost:9999", msg)
	}
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
