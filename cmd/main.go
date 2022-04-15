package main

import (
	"fmt"
	"log"
	"peercloud/crypto"
)

func main() {
	key, err := crypto.GenerateRSAKey()
	if err != nil {
		log.Fatal(err.Error())
	}

	address := crypto.HashSha256(key.PublicKey.N.Bytes())
	fmt.Println(key.PublicKey.N.Bytes(), len(key.PublicKey.N.Bytes()))
	fmt.Println(address, len(address))

	testContent := []byte("This is a test message")
	hash := crypto.HashSha256(testContent)

	signature, err := crypto.SignMessage(hash, key)
	if err != nil {
		log.Fatal(err.Error())
	}

	signatureValid := crypto.VerifyMessage(hash, signature, &key.PublicKey)
	fmt.Println(signatureValid)
}
