package main

import (
	"log"
	"peercloud/core"
)

func main() {
	if err := core.Upload("files/example.txt"); err != nil {
		log.Fatal(err)
	}
}
