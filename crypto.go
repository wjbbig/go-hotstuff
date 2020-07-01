package go_hotstuff

import (
	"container/list"
	"sync"
)

type Signer interface {
	TSign()
	TVerify()
}

type SignatureCache struct {
	cache list.List
	lock sync.Mutex
	verifiedSignatures map[string]bool
}

func NewSignatureCache() *SignatureCache {
	return &SignatureCache{
		verifiedSignatures: make(map[string]bool),
	}
}

