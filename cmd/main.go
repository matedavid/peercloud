package main

import (
	"fmt"
	"log"
	"peercloud/crypto"
)

func main() {
	/*
		filepath := "files/example.txt"
		manifest := core.ShardFile(filepath)
		fmt.Println(manifest)
	*/

	key, err := crypto.GenerateRSAKey()
	if err != nil {
		log.Fatal(err.Error())
	}

	testContent := []byte("This is a test message")
	hash := crypto.HashSha256(testContent)

	signature, err := crypto.SignMessage(hash, key)
	if err != nil {
		log.Fatal(err.Error())
	}

	isValid := crypto.VerifyMessage(hash, signature, &key.PublicKey)
	fmt.Println(isValid)
}
