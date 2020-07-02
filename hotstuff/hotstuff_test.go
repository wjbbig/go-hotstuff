package hotstuff

import (
	go_hotstuff "github.com/wjbbig/go-hotstuff"
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

func TestHotStuffImpl_GetSelfInfo(t *testing.T) {
	hsc := &HotStuffConfig{}
	hsc.ReadConfig()
	h := &HotStuffImpl{}
	h.Config = *hsc
	h.ID = 1
	t.Log(h.GetSelfInfo())
}
