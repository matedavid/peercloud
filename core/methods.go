package core

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"peercloud/network"
	"time"
)

func SendVersion(conn net.Conn, cfg *Config) error {
	version := network.VersionPayload{
		Timestamp:  time.Now().Unix(),
		Address:    cfg.Address,
		Port:       cfg.Port,
		Identifier: cfg.GetNodeIdentifier(),
	}

	header := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Version,
		Payload:     uint32(len(version.Write())),
	}
	err := header.Send(conn)
	if err != nil {
		return err
	}

	// Send version payload
	err = network.SendPayload(conn, &version)
	if err != nil {
		return err
	}

	// Recv Verack command
	err = header.Recv(conn)
	if err != nil {
		return err
	} else if header.Command != network.Verack {
		return errors.New("did not receive 'verack' header command")
	}
	return nil
}

func RecvVersion(conn net.Conn, header network.MessageHeader) error {
	versionInfo := network.VersionPayload{}
	err := network.ReceivePayload(conn, header.Payload, &versionInfo)
	if err != nil {
		return err
	}

	// TODO: Do something with the versionInfo
	fmt.Println("Version:", versionInfo)

	verack := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Verack,
		Payload:     0,
	}

	return verack.Send(conn)
}

func Store(conn net.Conn, header network.MessageHeader, cfg *Config) error {
	log.Println("Store:", header)

	data := network.UploadPayload{}
	err := network.ReceivePayload(conn, header.Payload, &data)
	if err != nil {
		return err
	}

	err = StoreShard(data.Content, data.Hash, cfg)
	if err != nil {
		return err
	}

	mh := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Stored,
		Payload:     0,
	}

	return mh.Send(conn)
}

func Retrieve(conn net.Conn, header network.MessageHeader, cfg *Config) error {
	log.Println("Retrieve:", header)

	// Receive shard hash
	data := network.DownloadPayload{}
	err := network.ReceivePayload(conn, header.Payload, &data)
	if err != nil {
		return err
	}

	content, err := RetrieveShard(data.Hash, cfg)
	if err == os.ErrNotExist {
		log.Println("Shard:", data.Hash, "does not exist")
		return err
	} else if err != nil {
		return err
	}

	mh := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Retrieved,
		Payload:     uint32(len(content)),
	}
	err = mh.Send(conn)
	if err != nil {
		return err
	}

	payload := network.GenericPayload{
		Content: content,
	}
	return network.SendPayload(conn, &payload)
}
