package hotstuff

import (
	go_hotstuff "go-hotstuff"
	"testing"
)

func TestGenerateGenesisBlock(t *testing.T) {
	block := GenerateGenesisBlock()
	t.Log(go_hotstuff.String(block))
}
