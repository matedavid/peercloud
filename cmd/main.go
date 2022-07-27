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
	"strconv"
)

type RpcListener int

var cfg *core.Config

func (l *RpcListener) Upload(filePath *string, reply *bool) error {
	err := core.Upload(*filePath, cfg)
	*reply = err == nil
	return err
}

type DownloadArgs struct {
	File       string
	OutputPath string
}

func (l *RpcListener) Download(args *DownloadArgs, reply *string) error {
	manifest, err := core.SearchManifestFromName(args.File, cfg)
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

	listener, err := net.Listen("tcp", cfg.Node.GetCompleteAddress())
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(listener, nil)
}

func tcpServer(tcpCfg *core.Config) {
	listener, err := net.Listen("tcp", tcpCfg.Node.GetCompleteAddress())
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

		switch mh.Command {
		case network.Store: // Store
			if err := core.Store(conn, mh, cfg); err != nil {
				log.Fatal(err)
			}
		case network.Retrieve: // Retrieve
			if err := core.Retrieve(conn, mh, cfg); err != nil {
				log.Fatal(err)
			}
		case network.Version: // Version
			if err := core.RecvVersion(conn, mh); err != nil {
				log.Fatal(err)
			}

		}
		conn.Close()
	}
}

func main() {
	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	path := ".peercloud/" + ip + ":" + fmt.Sprint(port)

	cfg = &core.Config{
		Node: network.Host{
			Address: net.ParseIP(ip),
			Port:    uint16(port),
		},
		Path: path,
	}

	tcpCfg := &core.Config{
		Node: network.Host{
			Address: net.ParseIP(ip),
			Port:    uint16(port + 1),
		},
	}

	fmt.Println(cfg.Node.GetCompleteAddress(), "-", cfg.Node.GetNodeIdentifier())

	go rpcServer(cfg)
	tcpServer(tcpCfg)
}
