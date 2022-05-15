package core

import (
	"errors"
	"net"
	"os"
	"peercloud/crypto"
	"peercloud/network"
)

func Upload(filePath string) error {
	key, err := crypto.GetRSAKey()
	if err != nil {
		return err
	}

	manifest, err := ShardFile(filePath, key)
	if err != nil {
		return err
	}

	for _, shard := range manifest.Shards {
		content, err := GetTmpShard(shard)
		if err != nil {
			return err
		}

		// TODO: Find suitable nodes
		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			return err
		}
		defer conn.Close()

		// Send shard content
		sendContent := append([]byte(shard), content...)

		header := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Store,
			Payload:     uint32(len(sendContent)), // This includes the size of the added shard
		}
		header.Send(conn)

		// Wait for Acknowledge header
		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Acknowledge {
			return errors.New("did not receive acknowledge message header")
		}

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

		err = RemoveTmpShard(shard)
		if err != nil {
			return err
		}
	}

	return nil
}

func Download(manifest *Manifest, outputPath string) error {
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
		conn, err := net.Dial("tcp", "localhost:8001")
		if err != nil {
			return err
		}
		defer conn.Close()

		header := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Retrieve,
			Payload:     64,
		}
		header.Send(conn)

		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Acknowledge {
			return errors.New("did not receive acknowledge message header")
		}

		// Send shard hash
		_, err = conn.Write([]byte(shard))
		if err != nil {
			return err
		}

		err = header.Recv(conn)
		if err != nil {
			return err
		} else if header.Command != network.Retrieved {
			return errors.New("did not receive retrieved message header")
		}

		mh := network.MessageHeader{
			NetworkCode: network.MAIN_NETWORK_CODE,
			Command:     network.Acknowledge,
			Payload:     0,
		}
		err = mh.Send(conn)
		if err != nil {
			return err
		}

		buff := make([]byte, header.Payload)
		n, err := conn.Read(buff)
		if err != nil {
			return err
		} else if n != int(header.Payload) {
			return errors.New("length of data received does not match payload")
		}

		/*
			decryptedContent, err := crypto.Decrypt(buff, key)
			if err != nil {
				log.Fatal(err.Error())
			}

			file.Write(decryptedContent)
		*/

		file.Write(buff)
	}

	return nil
}
