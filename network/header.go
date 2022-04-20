package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
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
	Unknown
)

func Command2Bytes(command NetworkCommand) [12]byte {
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

func Bytes2Command(byteCommand [12]byte) NetworkCommand {
	strCommand := string(byteCommand[:])
	strCommand = strings.Trim(strCommand, "\x00")

	if strCommand == "version" {
		return Version
	} else if strCommand == "verack" {
		return Verack
	} else if strCommand == "store" {
		return Store
	} else if strCommand == "retrieve" {
		return Retrieve
	}
	return Unknown
}

type MessageHeader struct {
	NetworkCode uint32
	Command     [12]byte
	Payload     uint32
}

func (mh *MessageHeader) Send(conn net.Conn) error {
	// Pack the struct
	payload := &bytes.Buffer{}
	err := binary.Write(payload, binary.LittleEndian, mh)
	if err != nil {
		return err
	}

	// Send to other node
	n, err := conn.Write(payload.Bytes())
	fmt.Println("Sent:", n, "bytes")
	return err
}

func (mh *MessageHeader) Recv(conn net.Conn) error {
	// Recieve header from connection
	data := make([]byte, 20)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	// Unpack struct
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, mh)
	return err
}
