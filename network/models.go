package network

import (
	"bytes"
	"encoding/binary"
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
	Timestamp    int64
	RandomNumber uint32
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
	bts := &bytes.Buffer{}
	binary.Write(bts, binary.LittleEndian, p)

	return bts.Bytes()
}

func (p *VersionPayload) Read(data []byte) {
	binary.Read(bytes.NewReader(data), binary.LittleEndian, p)
}
