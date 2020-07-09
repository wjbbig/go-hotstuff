package go_hotstuff

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/syndtr/goleveldb/leveldb"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"testing"
)

func TestHash(t *testing.T) {
	commands := []string{"ss", "aaa", "ddd"}
	parentHash := sha256.Sum256([]byte("parentHash"))

	block := &pb.Block{
		ParentHash: parentHash[:],
		Hash:       nil,
		Height:     1,
		Commands:   commands,
		Justify: &pb.QuorumCert{
			BlockHash: parentHash[:],
			ViewNum:   1,
		},
	}

	blockHash := Hash(block)

	t.Log(hex.EncodeToString(blockHash))
}

func TestBlock_PutAndGet(t *testing.T) {
	commands := []string{"ss", "aaa", "ddd"}
	parentHash := sha256.Sum256([]byte("parentHash"))

	block := &pb.Block{
		ParentHash: parentHash[:],
		Hash:       nil,
		Height:     1,
		Commands:   commands,
		Justify: &pb.QuorumCert{
			BlockHash: parentHash[:],
			ViewNum:   1,
		},
	}

	blockHash := Hash(block)
	block.Hash = blockHash

	impl := NewBlockStorageImpl("1")

	err := impl.Put(block)
	if err != nil {
		t.Fatal(block)
	}

	get, err := impl.Get(block.Hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(String(get))
}

func TestBlockNilParam(t *testing.T) {
	impl := NewBlockStorageImpl("1")
	get, err := impl.Get([]byte("ssss"))
	if err ==leveldb.ErrNotFound {
		t.Log("not found")
	}
	t.Log(get)
}
