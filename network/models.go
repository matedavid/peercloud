package network

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

type Model interface {
	// Convert Model to []byte
	Write() []byte
	// Populate Model from []byte
	Read([]byte)
}

type GenericPayload struct {
	Content []byte
}

type UploadPayload struct {
	Hash    string
	Content []byte
}

type DownloadPayload struct {
	Hash string
}

type VersionPayload struct {
	Timestamp  int64
	Address    net.IP
	Port       uint16
	Identifier string
}

// GenericPayload
func (p *GenericPayload) Write() []byte {
	return p.Content
}

func (p *GenericPayload) Read(data []byte) {
	p.Content = data
}

// UploadPayload
func (p *UploadPayload) Write() []byte {
	return append([]byte(p.Hash), p.Content...)
}

func (p *UploadPayload) Read(data []byte) {
	p.Hash = string(data[:64])
	p.Content = data[64:]
}

// DownloadPayload
func (p *DownloadPayload) Write() []byte {
	return []byte(p.Hash)
}

func (p *DownloadPayload) Read(data []byte) {
	p.Hash = string(data[:64])
}

// VersionPayload
func (p *VersionPayload) Write() []byte {
	// TODO: Should treat errors somehow
	timestampBytes := &bytes.Buffer{}
	binary.Write(timestampBytes, binary.LittleEndian, p.Timestamp)

	portBytes := &bytes.Buffer{}
	binary.Write(portBytes, binary.LittleEndian, p.Port)

	// TODO: Should create independent function to do this (net.IP => []byte)
	addressBytes := make([]byte, 4)
	bytesStr := strings.Split(p.Address.String(), ".")
	for i, b := range bytesStr {
		intVal, _ := strconv.Atoi(b)

		bBytes := &bytes.Buffer{}
		binary.Write(bBytes, binary.LittleEndian, uint8(intVal))
		addressBytes[i] = bBytes.Bytes()[0]
	}

	var data []byte
	data = append(data, timestampBytes.Bytes()...)
	data = append(data, addressBytes...)
	data = append(data, portBytes.Bytes()...)
	data = append(data, []byte(p.Identifier)...)

	return data
}

func (p *VersionPayload) Read(data []byte) {
	timestampBytes := data[:8]
	addressBytes := data[8:12]
	portBytes := data[12:16]
	identifierBytes := data[16:]

	// TODO: Should treat errors somehow
	binary.Read(bytes.NewReader(timestampBytes), binary.LittleEndian, &p.Timestamp)
	binary.Read(bytes.NewReader(portBytes), binary.LittleEndian, &p.Port)

	// TODO: Should create independent function to do this ([]byte => net.IP)
	bytesStr := make([]string, 4)
	for i, b := range addressBytes {
		var value uint8

		sliceByte := make([]byte, 1)
		sliceByte[0] = b

		binary.Read(bytes.NewReader(sliceByte), binary.LittleEndian, &value)
		bytesStr[i] = strconv.Itoa(int(value))
	}

	addressStr := bytesStr[0] + "." + bytesStr[1] + "." + bytesStr[2] + "." + bytesStr[3]
	p.Address = net.ParseIP(addressStr)
	p.Identifier = string(identifierBytes)
}
