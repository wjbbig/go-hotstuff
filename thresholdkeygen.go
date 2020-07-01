package go_hotstuff

import (
	"crypto"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/niclabs/tcrsa"
	"io/ioutil"
	"os"
)

const (
	size = 2048
	hashType = crypto.SHA256
	privateKeyFileType = "HOTSTUFF PRIVATE KEY"
	publicKeyFileType  = "HOTSTUFF PUBLIC KEY"
)

// GenerateThresholdKeys generate threshold signature keys
// need: how many signatures we need ,the same as 2f+1
// all: how many private keys we need generate, the same as N
// it may take a bit time to generate keys
func GenerateThresholdKeys(need, all int) (shares tcrsa.KeyShareList, meta *tcrsa.KeyMeta, err error) {
	k := uint16(need)
	l := uint16(all)

	return tcrsa.NewKey(size, k, l, nil)
}

// TODO find a better way to store keys

func WriteThresholdPrivateKeyToFile(privateKey *tcrsa.KeyShare, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("cannot open private key file")
	}
	defer file.Close()
	marshal, err := json.Marshal(*privateKey)
	if err != nil {
		return errors.New("cannot marshal private key")
	}
	b := &pem.Block{
		Type:    privateKeyFileType,
		Bytes:   marshal,
	}

	err = pem.Encode(file, b)
	if err != nil {
		return errors.New("write private key to file failed")
	}
	return nil
}

func WriteThresholdPublicKeyToFile(publicKey *tcrsa.KeyMeta, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("cannot open public key file")
	}
	defer file.Close()
	marshal, err := json.Marshal(*publicKey)
	if err != nil {
		return errors.New("cannot marshal private key")
	}
	b := &pem.Block{
		Type:    publicKeyFileType,
		Bytes:   marshal,
	}

	err = pem.Encode(file, b)
	if err != nil {
		return errors.New("write public key to file failed")
	}
	return nil
}

func ReadThresholdPrivateKeyToFile(filePath string) (*tcrsa.KeyShare, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("find private key failed")
	}

	b, _ := pem.Decode(file)
	if b == nil {
		return nil, errors.New("private key did not exist")
	}

	if b.Type != PRIVATEKEYFILETYPE {
		return nil, errors.New("file type did not match")
	}
	privateKey := new(tcrsa.KeyShare)
	err = json.Unmarshal(b.Bytes, privateKey)
	if err != nil {
		return nil, errors.New("parse private key failed")
	}

	return privateKey, nil
}

//func ReadThresholdPublicKeyToFile