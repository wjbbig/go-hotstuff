package go_hotstuff

import "testing"

func TestGeneratePrivateKey(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(key)
}

func TestWriteKeysToFile(t *testing.T) {
	key, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	err = WritePrivateKeyToFile(key, "keys/priv_sk")
	if err != nil {
		t.Fatal(err)
	}
	pubKey := key.PublicKey

	err = WritePublicKeyToFile(&pubKey, "keys/r1.pub")
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadPrivateKeyFromFile(t *testing.T) {
	privateKey, err := ReadPrivateKeyFromFile("keys/priv_sk")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(privateKey)
}

func TestReadPublicKeyFromFile(t *testing.T) {
	pubKey, err := ReadPublicKeyFromFile("keys/r1.pub")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pubKey)
}