package core

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"peercloud/crypto"
)

const DEFAULT_SHARD_PATH = ".peercloud/.tmp/"
const DEFAULT_MANIFEST_PATH = ".peercloud/.storage/"

type Manifest struct {
	Hash      string
	Filename  string
	Extension string
	Shards    []string
}

// Shards the file in the given filepath, saving the shards in the default path for shards
// and returning a Manifest object, which is saved in the default path for manfiests
func ShardFile(filepath string, key *rsa.PrivateKey) *Manifest {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	manifest := &Manifest{}
	manifest.Filename = path.Base(filepath)
	manifest.Extension = path.Ext(filepath)

	filenameHash := crypto.HashSha256([]byte(path.Base(filepath)))

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err.Error())
		}

		if err == io.EOF {
			break
		}

		// Compute hash identifier for the shard
		shardContent := append(buffer[:], filenameHash[:]...) // TODO: Should add more information to the shardContent
		shardHash := crypto.HashSha256(shardContent)

		shardHashString := crypto.HashAsString(shardHash)

		// Encrypt content of the shard
		encryptedBuffer, err := crypto.Encrypt(buffer[:n], key)
		if err != nil {
			log.Fatal(err)
		}

		err = saveShard(encryptedBuffer, shardHashString)
		if err != nil {
			log.Fatal(err.Error())
		}
		manifest.Shards = append(manifest.Shards, shardHashString)
	}

	manifest.Hash = hashManifest(manifest)

	err = saveManifest(manifest)
	if err != nil {
		log.Fatal(err.Error())
	}

	return manifest
}

/*
// Returns the content of a Shard
func GetShard(hash string) (string, error) {
}

// Returns if a shard exists
func ShardExists(hash string) bool {
}

// Gets the Manifest object from storage
func GetManifest(hash string) *Manifest {
}

// Returns if a Manifest file exists
func ManifestExists(hash string) bool {
}
*/

// Saves the shard's content in the default path for shards
func saveShard(content []byte, hash string) error {
	shardPath := path.Join(DEFAULT_SHARD_PATH, hash)

	file, err := os.Create(shardPath)
	if err != nil {
		return err
	}

	file.Write(content)
	file.Close()

	return nil
}

// Saves a Manifest object in the default path for manifests
func saveManifest(manifest *Manifest) error {
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return err
	}

	manifestPath := path.Join(DEFAULT_MANIFEST_PATH, manifest.Hash)
	file, err := os.Create(manifestPath)
	if err != nil {
		return err
	}

	file.Write(manifestJSON)
	file.Close()

	return nil
}

// Returns the Hash of a given file using a Merkle Tree
func hashManifest(manifest *Manifest) string {
	shards := manifest.Shards

	if len(shards)%2 != 0 {
		shards = append(shards, "0")
	}

	for len(shards) > 1 {
		l := len(shards)
		firstHash := shards[l-1]
		secondHash := shards[l-2]

		combinedHash := firstHash + secondHash
		shards = shards[:l-2]

		shards = append(shards, combinedHash)
	}

	return shards[0]
}
