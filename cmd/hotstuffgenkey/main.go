package main

import (
	"flag"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/logging"
	"os"
	"path"
	"strconv"
	"strings"
)

var logger = logging.GetLogger()

var filePath string
var k, l int

func init() {
	flag.StringVar(&filePath, "p", "", "where is the key generated")
	flag.IntVar(&k, "k", 3, "how many keys are needed to create a signature")
	flag.IntVar(&l, "l", 4, "how many keys are needed to generate")
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
	privateKeys, publicKey, err := go_hotstuff.GenerateThresholdKeys(k, l)
	if err != nil {
		logger.Fatal(err)
	}
	for i, key := range privateKeys {
		privateKeyPath := path.Join(filePath,"r"+strconv.Itoa(i+1)+".key")
		err = go_hotstuff.WriteThresholdPrivateKeyToFile(key, privateKeyPath)
		if err != nil {
			logger.Fatal(err)
		}
	}
	publicKeyPath := path.Join(filePath, "pub.key")
	err = go_hotstuff.WriteThresholdPublicKeyToFile(publicKey, publicKeyPath)
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
