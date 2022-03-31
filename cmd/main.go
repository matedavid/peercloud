package main

import (
	"fmt"
	"peercloud/core"
)

func main() {
	filepath := "files/example.txt"
	manifest := core.ShardFile(filepath)
	fmt.Println(manifest)
}
