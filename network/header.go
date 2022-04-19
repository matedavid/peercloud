package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// Network codes
const MAIN_NETWORK_CODE = 0xDD13

// Network commands
type NetworkCommand uint32

const (
	Version NetworkCommand = iota
	Verack
	Store
	Retrieve
)

func NetworkCommandBytes(command NetworkCommand) [12]byte {
	var byteCommand [12]byte

	switch command {
	case Version:
		byteCommand = [12]byte{'v', 'e', 'r', 's', 'i', 'o', 'n'}
	case Verack:
		byteCommand = [12]byte{'v', 'e', 'r', 'a', 'c', 'k'}
	case Store:
		byteCommand = [12]byte{'s', 't', 'o', 'r', 'e'}
	case Retrieve:
		byteCommand = [12]byte{'r', 'e', 't', 'r', 'i', 'e', 'v', 'e'}
	}

	return byteCommand
}

type MessageHeader struct {
	NetworkCode uint32
	Command     [12]byte
	Payload     uint32
}

func (mh *MessageHeader) Send(conn net.Conn) {
	// Pack the struct
	payload := &bytes.Buffer{}
	err := binary.Write(payload, binary.LittleEndian, mh)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Send to other node
	n, err := conn.Write(payload.Bytes())
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Sent:", n, "bytes")
}

func (mh *MessageHeader) Recv(conn net.Conn) {
	data := make([]byte, 20)
	_, err := conn.Read(data)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, mh)
	if err != nil {
		log.Fatal(err.Error())
	}
}
