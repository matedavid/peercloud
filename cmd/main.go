package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"peercloud/core"
	"peercloud/network"
)

type RpcListener int

func (l *RpcListener) Upload(filePath *string, reply *bool) error {
	err := core.Upload(*filePath)
	*reply = err == nil
	return nil
}

func (l *RpcListener) Download(fileName *string, reply *string) error {
	manifest, err := core.SearchManifestFromName(*fileName)
	if err != nil {
		return err
	}

	err = core.Download(manifest)
	if err != nil {
		return err
	}
	*reply = ".peercloud/.tmp/example.txt"

	return nil
}

func rpcServer(cfg *core.Config) {
	rpcListener := new(RpcListener)
	rpc.Register(rpcListener)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", cfg.GetCompleteAddress())
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(listener, nil)
}

func main() {
	if os.Args[1] == "server" {
		cfg := &core.Config{
			Address: "127.0.0.1",
			Port:    8000,
		}

		fmt.Println(cfg.GetCompleteAddress())

		rpcServer(cfg)
	} else if os.Args[1] == "client" {
		client, err := rpc.DialHTTP("tcp", "127.0.0.1:8000")
		if err != nil {
			log.Fatal(err)
		}

		var reply bool
		err = client.Call("RpcListener.Upload", "files/example.txt", &reply)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Upload:", reply)

		var dreply string
		err = client.Call("RpcListener.Download", "example.txt", &dreply)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Download:", dreply)

	} else if os.Args[1] == "tcpServer" {
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
