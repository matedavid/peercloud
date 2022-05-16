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
	return err
}

type DownloadArgs struct {
	File       string
	OutputPath string
}

func (l *RpcListener) Download(args *DownloadArgs, reply *string) error {
	manifest, err := core.SearchManifestFromName(args.File)
	if err != nil {
		return err
	}

	err = core.Download(manifest, args.OutputPath)
	if err != nil {
		return err
	}
	*reply = args.OutputPath

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
			Address: net.ParseIP("127.0.0.1"),
			Port:    8000,
		}

		fmt.Println(cfg.GetCompleteAddress())
		rpcServer(cfg)

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

			if mh.Command == network.Version {
				core.RecvVersion(conn, mh)
			} else if mh.Command == network.Store {
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
