package core

import (
	"errors"
	"net"
	"peercloud/crypto"
	"peercloud/network"
)

func Upload(filePath string) error {
	// TODO: This should get the node's key, not generate a new one
	key, err := crypto.GenerateRSAKey()
	if err != nil {
		return err
	}

	manifest, err := ShardFile(filePath, key)
	if err != nil {
		return err
	}

	for _, shard := range manifest.Shards {
		content, err := GetShard(shard)
		if err != nil {
			return err
		}

		// TODO: Find suitable nodes
		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			return err
		}
		defer conn.Close()

		header := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Store,
			Payload:     uint32(len(content) + 64), // +64 because sending shard
		}

		header.Send(conn)

		// Wait for Acknowledge header
		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Acknowledge {
			return errors.New("did not receive acknowledge message header")
		}

		// Send shard content
		sendContent := append([]byte(shard), content...)

		_, err = conn.Write(sendContent)
		if err != nil {
			return err
		}

		// Wait for Stored header
		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Stored {
			return errors.New("did not receive stored message header")
		}
	}

	return nil
}

func Download(manifest *Manifest) error {
	return nil
}
