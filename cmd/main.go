package main

import (
	"fmt"
	"log"
	"peercloud/core"
	"peercloud/crypto"
)

func main() {
	key, err := crypto.GenerateRSAKey()
	if err != nil {
		log.Fatal(err.Error())
	}

	manifest := core.ShardFile("files/example.txt", key)
	fmt.Println(manifest)
}
