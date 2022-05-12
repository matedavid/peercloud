package main

import (
	"log"
	"net"
	"os"
	"peercloud/core"
	"peercloud/network"
)

func main() {
	if os.Args[1] == "upload" {
		if err := core.Upload("files/example.txt"); err != nil {
			log.Fatal(err.Error())
		}
	} else if os.Args[1] == "store" {
		listener, err := net.Listen("tcp", "localhost:8001")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err.Error())
			}
			defer conn.Close()

			mh := network.MessageHeader{}
			mh.Recv(conn)

			core.Store(conn, mh)
		}

	}
}
