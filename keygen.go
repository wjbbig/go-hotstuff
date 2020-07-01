package go_hotstuff

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

/*
Generate private and public keys for hotstuff (ECDSA)
Deprecated， PLEASE SEE thresholdkeygen.go
*/

const (
	PRIVATEKEYFILETYPE = "HOTSTUFF PRIVATE KEY"
	PUBLICKEYFILETYPE  = "HOTSTUFF PUBLIC KEY"
)

// GeneratePrivateKey 使用ecdsa生成私钥
func GeneratePrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	curve := elliptic.P256()
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	return
}

// WritePrivateKeyToFile 将私钥写入磁盘
func WritePrivateKeyToFile(privateKey *ecdsa.PrivateKey, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("cannot open private key file")
	}
	defer file.Close()
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return errors.New("cannot marshal private key")
	}

	b := &pem.Block{
		Type:  PRIVATEKEYFILETYPE,
		Bytes: keyBytes,
	}

	err = pem.Encode(file, b)
	if err != nil {
		return errors.New("write private key to file failed")
	}
	return nil
}

// WritePublicKeyToFile 将公钥写入磁盘
func WritePublicKeyToFile(publicKey *ecdsa.PublicKey, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("open public key file failed")
	}
	defer file.Close()

	keyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return errors.New("marshal public key failed")
	}
	b := &pem.Block{
		Type:  PUBLICKEYFILETYPE,
		Bytes: keyBytes,
	}
	err = pem.Encode(file, b)
	if err != nil {
		return errors.New("write public key to file failed")
	}
	return nil
}

// ReadPrivateKeyFromFile 从文件中读取私钥
func ReadPrivateKeyFromFile(filePath string) (*ecdsa.PrivateKey, error) {
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
	privateKey, err := x509.ParseECPrivateKey(b.Bytes)
	if err != nil {
		return nil, errors.New("parse private key failed")
	}

	return privateKey, nil
}

// ReadPublicKeyFromFile 从硬盘中读取公钥
func ReadPublicKeyFromFile(filePath string) (*ecdsa.PublicKey, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("find public key failed")
	}

	b, _ := pem.Decode(file)
	if b == nil {
		return nil, errors.New("public key did not exist")
	}
	if b.Type != PUBLICKEYFILETYPE {
		return nil, errors.New("file type did not match")
	}

	k, err := x509.ParsePKIXPublicKey(b.Bytes)
	if err != nil {
		return nil, errors.New("parse private key failed")
	}

	key, ok := k.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key was of wrong type")
	}

	return key, nil
}


