all: test build

test: crypto/* core/* network/*
	go test ./crypto/ ./network/ ./core/ -v 

build: 
	go build ./cmd/main.go 

run:
	go run cmd/main.go

clean:
	rm main && rm -rf .peercloud/

