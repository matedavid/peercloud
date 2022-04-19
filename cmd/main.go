package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"peercloud/network"
)

func main() {
	/*
		if err := core.Upload("files/example.txt"); err != nil {
			log.Fatal(err)
		}
	*/

	if os.Args[1] == "client" {
		conn, err := net.Dial("tcp", "127.0.0.1:8000")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		header := &network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.NetworkCommandBytes(network.Store),
			Payload:     uint32(0),
		}

		err = header.Send(conn)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else if os.Args[1] == "server" {
		listener, err := net.Listen("tcp", "127.0.0.1:8000")
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		header := &network.MessageHeader{}
		err = header.Recv(conn)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(header)
	}
}
