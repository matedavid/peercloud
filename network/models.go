package network

type Model interface {
	// Convert Model to []byte
	Write() []byte
	// Populate Model from []byte
	Read([]byte)
}

type GenericPayload struct {
	Content []byte
}

func (p *GenericPayload) Write() []byte {
	return p.Content
}

func (p *GenericPayload) Read(data []byte) {
	p.Content = data
}

type UploadPayload struct {
	Hash    string
	Content []byte
}

func (p *UploadPayload) Write() []byte {
	return append([]byte(p.Hash), p.Content...)
}

func (p *UploadPayload) Read(data []byte) {
	p.Hash = string(data[:64])
	p.Content = data[64:]
}

type DownloadPayload struct {
	Hash string
}

func (p *DownloadPayload) Write() []byte {
	return []byte(p.Hash)
}

func (p *DownloadPayload) Read(data []byte) {
	p.Hash = string(data[:64])
}
