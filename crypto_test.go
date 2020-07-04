package go_hotstuff

import (
	"github.com/niclabs/tcrsa"
	"testing"
)

func TestThresholdSignature(t *testing.T) {
	cryptoStrBytes := []byte("hello, world")
	keys, meta, err := GenerateThresholdKeys(3, 4)
	if err != nil {
		t.Fatal(err)
	}

	docHash, err := CreateDocumentHash(cryptoStrBytes, meta)
	if err != nil {
		t.Fatal(err)
	}

	sigList := make(tcrsa.SigShareList, len(keys))
	i := 0
	for i< len(keys) {
		sigList[i], err = TSign(docHash, keys[i], meta)
		if err != nil {
			t.Fatal(err)
		}
		err := VerifyPartSig(sigList[i], docHash, meta)
		if err != nil {
			t.Fatal(err)
		}
	}
	signature, err := CreateFullSignature(docHash, sigList, meta)
	if err != nil {
		t.Fatal(err)
	}

	_, err = TVerify(meta, signature, docHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
