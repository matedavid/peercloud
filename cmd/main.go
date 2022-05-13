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
	} else if os.Args[1] == "download" {
		manifest, err := core.GetManifest("b73cc7de1d94ba3ba5d0b6827b9496b75c5961331a071b0395cd3423a0029de2")
		if err != nil {
			log.Fatal(err)
		}

		if err := core.Download(manifest); err != nil {
			log.Fatal(err)
		}
	} else if os.Args[1] == "server" {
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

			mh := network.MessageHeader{}
			mh.Recv(conn)

			if mh.Command == network.Store {
				core.Store(conn, mh)
			} else if mh.Command == network.Retrieve {
				if err := core.Retrieve(conn, mh); err != nil {
					log.Fatal(err)
				}
			}

			conn.Close()
		}
	}
}
