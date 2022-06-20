package core

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"peercloud/network"
	"time"
)

func generateRandomNumber(seed int64) uint32 {
	source := rand.NewSource(seed)
	return rand.New(source).Uint32()
}

func SendVersion(conn net.Conn, cfg *Config) error {
	currentTime := time.Now().Unix()
	randomNumber := generateRandomNumber(currentTime)

	version := network.VersionPayload{
		Timestamp:    currentTime,
		RandomNumber: randomNumber,
	}

	fmt.Println(version.Write())

	header := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Version,
		Payload:     uint32(len(version.Write())),
	}
	err := header.Send(conn)
	if err != nil {
		return err
	}

	fmt.Println("Before send:", version)

	// Send version payload
	err = network.SendPayload(conn, &version)
	if err != nil {
		return err
	}

	// Recv Verack command
	err = header.Recv(conn)
	if err != nil {
		return err
	} else if header.Command != network.Verack {
		return errors.New("did not receive 'verack' header command")
	}

	// Receive verack VersionInfo response
	verack := network.VersionPayload{}
	err = network.ReceivePayload(conn, header.Payload, &verack)
	if err != nil {
		return err
	}

	fmt.Println("SendVersion:", verack)

	if verack.Timestamp < currentTime || verack.RandomNumber-1 != randomNumber {
		return errors.New("verack information does not match with sent version command")
	}

	return nil
}

func RecvVersion(conn net.Conn, header network.MessageHeader) error {
	versionInfo := network.VersionPayload{}
	err := network.ReceivePayload(conn, header.Payload, &versionInfo)
	if err != nil {
		return err
	}

	fmt.Println("versionInfo:", versionInfo)

	verackInfo := network.VersionPayload{
		Timestamp:    time.Now().Unix(),
		RandomNumber: versionInfo.RandomNumber + 1,
	}

	fmt.Println("RecvVersion:", versionInfo)

	verack := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Verack,
		Payload:     uint32(len(verackInfo.Write())),
	}
	err = verack.Send(conn)
	if err != nil {
		return err
	}

	return network.SendPayload(conn, &verackInfo)
}

func Store(conn net.Conn, header network.MessageHeader, cfg *Config) error {
	log.Println("Store:", header)

	data := network.UploadPayload{}
	err := network.ReceivePayload(conn, header.Payload, &data)
	if err != nil {
		return err
	}

	err = StoreShard(data.Content, data.Hash, cfg)
	if err != nil {
		return err
	}

	mh := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Stored,
		Payload:     0,
	}

	return mh.Send(conn)
}

func Retrieve(conn net.Conn, header network.MessageHeader, cfg *Config) error {
	log.Println("Retrieve:", header)

	// Receive shard hash
	data := network.DownloadPayload{}
	err := network.ReceivePayload(conn, header.Payload, &data)
	if err != nil {
		return err
	}

	content, err := RetrieveShard(data.Hash, cfg)
	if err == os.ErrNotExist {
		log.Println("Shard:", data.Hash, "does not exist")
		return err
	} else if err != nil {
		return err
	}

	mh := network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Retrieved,
		Payload:     uint32(len(content)),
	}
	err = mh.Send(conn)
	if err != nil {
		return err
	}

	payload := network.GenericPayload{
		Content: content,
	}
	return network.SendPayload(conn, &payload)
}
