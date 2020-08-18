package config

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type KeyHolder struct {
	Key *rsa.PrivateKey
}

func (c *KeyHolder) LoadKeys(filename string) (err error) {

	log.Print("Load RSA Key")
	err = c.readKey(filename)

	if err == nil {
		return
	}

	err = nil
	reader := rand.Reader
	bitSize := 2048

	c.Key, err = rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return
	}

	err = c.saveKey(filename)
	return
}

func (c *KeyHolder) readKey(filename string) (err error) {
	outFile, err := os.Open(filename)
	if err != nil {
		return
	}
	defer outFile.Close()

	reader := bufio.NewReader(outFile)

	fileContents, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(fileContents)
	if block == nil || block.Type != "PRIVATE KEY" {
		panic(errors.New("failed to parse PEM block containing the key"))
	}
	c.Key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return
}
func (c *KeyHolder) saveKey(filename string) (err error) {
	inputFile, err := os.Create(filename)
	if err != nil {
		return
	}
	defer inputFile.Close()

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(c.Key),
	}

	err = pem.Encode(inputFile, privateKey)
	if err != nil {
		return
	}
	return
}

var keyHolderMut sync.Mutex

var keyHolder *KeyHolder

func GetKeyHolder() *KeyHolder {
	keyHolderMut.Lock()
	defer keyHolderMut.Unlock()

	if keyHolder == nil {
		keyHolder = &KeyHolder{}
	}
	return keyHolder
}
