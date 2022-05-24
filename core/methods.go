package core

import (
	"log"
	"net"
	"os"
	"peercloud/network"
)

func Store(conn net.Conn, header network.MessageHeader) error {
	log.Println("Store:", header)

	buff, err := network.ReceivePayload(conn, header.Payload)
	if err != nil {
		return err
	}

	log.Println("Received payload of:", len(buff), "bytes")

	hash := string(buff[:64])
	content := buff[64:]

	err = StoreShard(content, hash)
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

func Retrieve(conn net.Conn, header network.MessageHeader) error {
	log.Println("Retrieve:", header)

	// Receive shard hash
	shardHash, err := network.ReceivePayload(conn, header.Payload)
	if err != nil {
		return err
	}

	content, err := RetrieveShard(string(shardHash))
	if err == os.ErrNotExist {
		log.Println("Shard:", shardHash, "does not exist")
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

	return network.SendPayload(conn, content)
}
