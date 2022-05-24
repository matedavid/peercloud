package core

import (
	"errors"
	"log"
	"net"
	"os"
	"peercloud/network"
)

func Store(conn net.Conn, header network.MessageHeader) error {
	log.Println("Store:", header)

	payload := header.Payload

	ack := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Acknowledge,
		Payload:     0,
	}
	ack.Send(conn)

	buff := make([]byte, payload)
	n, err := conn.Read(buff)
	if err != nil {
		return err
	} else if n != int(payload) {
		return errors.New("length of data received does not match payload")
	}

	log.Println("Received payload of:", n, "bytes")

	hash := string(buff[:64])
	content := buff[64:n]

	err = StoreShard(content, hash)
	if err != nil {
		return err
	}

	mh := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Stored,
		Payload:     0,
	}

	err = mh.Send(conn)
	if err != nil {
		return err
	}

	return nil
}

func Retrieve(conn net.Conn, header network.MessageHeader) error {
	log.Println("Retrieve:", header)

	ack := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Acknowledge,
		Payload:     0,
	}
	err := ack.Send(conn)
	if err != nil {
		return err
	}

	// Receive shard hash
	shardHash := make([]byte, header.Payload)
	n, err := conn.Read(shardHash)
	if err != nil {
		return err
	} else if n != int(header.Payload) {
		return errors.New("length of data received does not match payload")
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

	// Receive acknowledge message header
	err = mh.Recv(conn)
	if err != nil {
		return err
	} else if mh.Command != network.Acknowledge {
		return errors.New("did not receive acknowledge message header")
	}

	// Send content
	_, err = conn.Write(content)
	if err != nil {
		return err
	}

	return nil
}
