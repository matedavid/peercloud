package main

import (
	"fmt"
	"peercloud/crypto"
)

func main() {
	/*
		filepath := "files/example.txt"
		manifest := core.ShardFile(filepath)
		fmt.Println(manifest)
	*/

	key := crypto.GenerateRSAKey()
	fmt.Println(key.PublicKey)
}
