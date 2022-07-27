package network

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
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
	Host        Host
	Payload     uint32
}

func (mh *MessageHeader) Send(conn net.Conn) error {
	// Pack the struct into data
	codeBytes := bytes.Buffer{}
	binary.Write(&codeBytes, binary.LittleEndian, mh.NetworkCode)

	commandBytes := bytes.Buffer{}
	binary.Write(&commandBytes, binary.LittleEndian, mh.Command)

	ipBytes := make([]byte, 4)
	ipStr := strings.Split(mh.Host.Address.String(), ".")
	for idx, ipComponent := range ipStr {
		value, _ := strconv.Atoi(ipComponent)

		tmpBytes := bytes.Buffer{}
		binary.Write(&tmpBytes, binary.LittleEndian, uint8(value))
		ipBytes[idx] = tmpBytes.Bytes()[0]
	}

	portBytes := bytes.Buffer{}
	binary.Write(&portBytes, binary.LittleEndian, mh.Host.Port)

	payloadBytes := bytes.Buffer{}
	binary.Write(&payloadBytes, binary.LittleEndian, mh.Payload)

	var data []byte
	data = append(data, codeBytes.Bytes()...)
	data = append(data, commandBytes.Bytes()...)
	data = append(data, ipBytes...)
	data = append(data, portBytes.Bytes()...)
	data = append(data, payloadBytes.Bytes()...)

	// Send to other node
	_, err := conn.Write(data)
	return err

}

func (mh *MessageHeader) Recv(conn net.Conn) error {
	// Receive header from connection
	data := make([]byte, MESSAGE_HEADER_BYTES)
	_, err := conn.Read(data)
	if err != nil {
		return err
	}

	// Unpack struct
	codeBytes := data[0:4]      // NetworkCode
	commandBytes := data[4:5]   // Command
	ipBytes := data[5:9]        // Host.Address
	portBytes := data[9:11]     // Host.Port
	payloadBytes := data[11:15] // Payload

	if err := binary.Read(bytes.NewReader(codeBytes), binary.LittleEndian, &mh.NetworkCode); err != nil {
		return err
	}
	if err := binary.Read(bytes.NewReader(commandBytes), binary.LittleEndian, &mh.Command); err != nil {
		return err
	}
	if err := binary.Read(bytes.NewReader(portBytes), binary.LittleEndian, &mh.Host.Port); err != nil {
		return err
	}
	if err := binary.Read(bytes.NewReader(payloadBytes), binary.LittleEndian, &mh.Payload); err != nil {
		return err
	}

	// TODO: Should create independent function to do this ([]byte => net.IP)
	bytesStr := make([]string, 4)
	for i, b := range ipBytes {
		var value uint8

		sliceByte := make([]byte, 1)
		sliceByte[0] = b

		binary.Read(bytes.NewReader(sliceByte), binary.LittleEndian, &value)
		bytesStr[i] = strconv.Itoa(int(value))
	}

	ipStr := bytesStr[0] + "." + bytesStr[1] + "." + bytesStr[2] + "." + bytesStr[3]
	mh.Host.Address = net.ParseIP(ipStr)

	return nil
}
