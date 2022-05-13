all: test build

test: crypto/* core/* network/*
	go test ./crypto/ ./network/ ./core/ -v 

build: 
	go build ./cmd/main.go 

clean:
	rm main

