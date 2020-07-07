package config

import (
	"github.com/niclabs/tcrsa"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	"time"
)

var logger = logging.GetLogger()

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

type ReplicaInfo struct {
	ID         uint32
	Address    string `mapstructure:"listen-address"`
	PrivateKey string `mapstructure:"privatekeypath"`
}

// ReadConfig reads hotstuff config from yaml file
func (hsc *HotStuffConfig) ReadConfig() {
	logger.Debug("[HOTSTUFF] Read config")
	viper.AddConfigPath("/opt/hotstuff/config/")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("./config")
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
