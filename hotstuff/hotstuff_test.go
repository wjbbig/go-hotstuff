package hotstuff

import (
	go_hotstuff "go-hotstuff"
	"testing"
)

func TestGenerateGenesisBlock(t *testing.T) {
	block := GenerateGenesisBlock()
	t.Log(go_hotstuff.String(block))
}

func TestHotStuffConfig_ReadConfig(t *testing.T) {
	hsc := &HotStuffConfig{}
	hsc.ReadConfig()
	t.Log(hsc.Cluster[0].Address)
}
