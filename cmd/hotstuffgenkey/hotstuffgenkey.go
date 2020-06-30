package main

import (
	"flag"
	go_hotstuff "go-hotstuff"
	"go-hotstuff/logging"
	"os"
	"strings"
)

var logger = logging.GetLogger()

var filePath string

func init() {
	flag.StringVar(&filePath, "filepath", "", "where is the key generated")
}

func main() {
	flag.Parse()
	if filePath == "" {
		flag.Usage()
	}
	dirPath := filePath[:strings.LastIndex(filePath, "/")]
	exist, err := isDirExist(dirPath)
	if err != nil {
		logger.Fatal(err)
	}
	if !exist {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			logger.Fatal(err)
		}
	}
	privateKeyPath := filePath
	publicKeyPath := filePath + ".pub"

	privateKey, err := go_hotstuff.GeneratePrivateKey()
	if err != nil {
		logger.Fatal(err)
	}

	err = go_hotstuff.WritePrivateKeyToFile(privateKey, privateKeyPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = go_hotstuff.WritePublicKeyToFile(&privateKey.PublicKey, publicKeyPath)
	if err != nil {
		logger.Fatal(err)
	}
}

func isDirExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
