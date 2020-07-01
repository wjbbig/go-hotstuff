package go_hotstuff

import (
	"crypto"
	"crypto/rsa"
	"github.com/niclabs/tcrsa"
)

func CreateDocumentHash(msgBytes []byte, meta *tcrsa.KeyMeta) ([]byte, error) {
	documentHash, err := tcrsa.PrepareDocumentHash(meta.PublicKey.Size(), crypto.SHA256, msgBytes)
	if err != nil {
		return nil, err
	}
	return documentHash, nil
}

// TSign create the partial signature of replica
func TSign(msgBytes []byte, privateKey *tcrsa.KeyShare, publicKey *tcrsa.KeyMeta) (*tcrsa.SigShare, error) {
	// get the msg hash
	documentHash, err := CreateDocumentHash(msgBytes, publicKey)
	if err != nil {
		return nil, err
	}
	// sign
	partSig, err := privateKey.Sign(documentHash, crypto.SHA256, publicKey)
	if err != nil {
		return nil, err
	}
	// ensure the partial signature is correct
	err = partSig.Verify(documentHash, publicKey)
	if err != nil {
		return nil, err
	}

	return partSig, nil
}

func CreateFullSignature(msgBytes []byte, partSigs *tcrsa.SigShareList, publicKey *tcrsa.KeyMeta) (tcrsa.Signature, error) {
	// get the msg hash
	documentHash, err := CreateDocumentHash(msgBytes, publicKey)
	if err != nil {
		return nil, err
	}
	signature, err := partSigs.Join(documentHash, publicKey)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func TVerify(publicKey *tcrsa.KeyMeta, signature tcrsa.Signature, msgBytes []byte) (bool, error) {
	// get the msg hash
	documentHash, err := CreateDocumentHash(msgBytes, publicKey)
	if err != nil {
		return false, err
	}
	err = rsa.VerifyPKCS1v15(publicKey.PublicKey, crypto.SHA256, documentHash, signature)
	if err != nil {
		return false, err
	}
	return true, nil
}