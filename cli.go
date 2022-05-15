package main

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/akamensky/argparse"
)

func getRpcClient() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func upload(uploadFile string) error {
	client := getRpcClient()

	var success bool
	err := client.Call("RpcListener.Upload", uploadFile, &success)
	if err != nil {
		return err
	} else if !success {
		return errors.New("Error uploading file")
	}

	return nil
}

type DownloadArgs struct {
	File       string
	OutputPath string
}

func download(file string, output string) error {
	client := getRpcClient()

	args := DownloadArgs{
		File:       file,
		OutputPath: output,
	}

	var returnPath string
	err := client.Call("RpcListener.Download", &args, &returnPath)
	if err != nil {
		return err
	}

	fmt.Println("File downloaded:", returnPath)

	return nil
}

func main() {
	parser := argparse.NewParser("peercloud client", "Interact with your peercloud node running in your machine")

	uploadCommand := parser.NewCommand("upload", "Upload a file to the network")
	uploadFile := uploadCommand.String("f", "file", &argparse.Options{Required: true, Help: "Path to the file that will be uploaded"})

	downloadCommand := parser.NewCommand("download", "Download a file from the network")
	downloadFile := downloadCommand.String("f", "file", &argparse.Options{Required: true, Help: "Name of the file to be downloaded"})
	outputFile := downloadCommand.String("o", "output", &argparse.Options{Required: true, Help: "Path where the file will be downloaded"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if len(*uploadFile) != 0 {
		err := upload(*uploadFile)
		if err != nil {
			log.Fatal("Error while uploading:", err)
		}

		fmt.Println("File uploaded successfully")
	} else if len(*downloadFile) != 0 {
		err := download(*downloadFile, *outputFile)
		if err != nil {
			log.Fatal("Error while downloading:", err)
		}
		fmt.Println("File downloaded successfully")
	}
}
