package core

import (
	"errors"
	"net"
	"os"
	"peercloud/crypto"
	"peercloud/network"
)

func Upload(filePath string, cfg *Config) error {
	/*
		key, err := crypto.GetRSAKey()
		if err != nil {
			return err
		}
	*/

	// TEMPORAL: Should use GetRSAKey()
	key, err := crypto.GenerateRSAKey(false)
	if err != nil {
		return err
	}
	///    ///    ////

	manifest, err := ShardFile(filePath, key, cfg)
	if err != nil {
		return err
	}

	for _, shard := range manifest.Shards {
		content, err := GetTmpShard(shard, cfg)
		if err != nil {
			return err
		}

		// TODO: Find suitable nodes
		conn, err := net.Dial("tcp", "localhost:8003")
		if err != nil {
			return err
		}

		// Send shard content (shard hash + shard content)
		payload := network.UploadPayload{
			Hash:    shard,
			Content: content,
		}

		// Send Store Header
		header := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Store,
			Payload:     uint32(len(payload.Write())),
		}
		err = header.Send(conn)
		if err != nil {
			return err
		}

		// Send shard content information
		err = network.SendPayload(conn, &payload)
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

		err = RemoveTmpShard(shard, cfg)
		if err != nil {
			return err
		}

		conn.Close()
	}

	return nil
}

func Download(manifest *Manifest, outputPath string, cfg *Config) error {
	/*
		key, err := crypto.GetRSAKey()
		if err != nil {
			return err
		}
	*/

	//filePath := path.Join(".peercloud/.tmp/", manifest.Filename)

	file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, shard := range manifest.Shards {
		// TODO: Find suitable nodes
		conn, err := net.Dial("tcp", "localhost:8003")
		if err != nil {
			return err
		}

		header := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Retrieve,
			Payload:     64, // The length of hash in []byte
		}
		header.Send(conn)

		payload := network.DownloadPayload{
			Hash: shard,
		}

		// Send shard hash that we want to retrieve
		err = network.SendPayload(conn, &payload)
		if err != nil {
			return err
		}

		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Retrieved {
			return errors.New("did not receive retrieved message header")
		}

		content := network.GenericPayload{}
		err = network.ReceivePayload(conn, header.Payload, &content)
		if err != nil {
			return err
		}

		/*
			decryptedContent, err := crypto.Decrypt(buff, key)
			if err != nil {
				log.Fatal(err.Error())
			}

			file.Write(decryptedContent)
		*/

		file.Write(content.Content)
		conn.Close()
	}

	return nil
}
