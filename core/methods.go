package core

import (
	"log"
	"net"
	"os"
	"peercloud/network"
)

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
