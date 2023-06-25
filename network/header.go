package network

import (
	"bytes"
	"encoding/binary"
	"net"
)

// Network codes
const MAIN_NETWORK_CODE = 0xDD13

// Network commands
type NetworkCommand uint8

const (
	Version NetworkCommand = iota
	Verack
	Store
	Stored
	Retrieve
	Retrieved
	Acknowledge
	Unknown
)

const MESSAGE_HEADER_BYTES = 15

type MessageHeader struct {
	NetworkCode uint32
	Command     NetworkCommand
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
	_, err = conn.Write(payload.Bytes())
	return err
}

func (mh *MessageHeader) Recv(conn net.Conn) error {
	// Receive header from connection
	data := make([]byte, 12)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	// Unpack struct
	err = binary.Read(bytes.NewReader(data), binary.LittleEndian, mh)
	return err
}
