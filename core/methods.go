package core

import (
	"errors"
	"fmt"
	"net"
	"peercloud/network"
)

func Store(conn net.Conn, header network.MessageHeader) error {
	fmt.Println(header)

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

	fmt.Println("Received payload of:", n, "bytes")

	hash := string(buff[:64])
	content := buff[64:]

	// TODO: Save shard
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
