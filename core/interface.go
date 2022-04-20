package core

import (
	"log"
	"net"
	"peercloud/network"
)

func Upload(filePath string) error {
	/*
		// TODO: This should get the node's key, not generate a new one
		key, err := crypto.GenerateRSAKey()
		if err != nil {
			return err
		}

		manifest, err := ShardFile(filePath, key)
		if err != nil {
			return err
		}
	*/

	// === Temporal ===
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	header := &network.MessageHeader{
		NetworkCode: network.MAIN_NETWORK_CODE,
		Command:     network.Command2Bytes(network.Store),
		Payload:     uint32(0),
	}

	header.Send(conn)

	// === ===

	/*
		for _, shard := range manifest.Shards {
			content, err := GetShard(shard)
			if err != nil {
				return err
			}

			header := network.MessageHeader{
				network.MAIN_NETWORK_CODE,
				network.NetworkCommandBytes(network.Store),
				uint32(len(content)),
			}
		}
	*/

	return nil
}
