package core

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"peercloud/crypto"
)

const DEFAULT_TMP_PATH = ".peercloud/.tmp/"
const DEFAULT_MANIFEST_PATH = ".peercloud/.storage/"
const DEFAULT_SHARD_PATH = ".peercloud/.shards/"

type Manifest struct {
	Hash      string
	Filename  string
	Extension string
	Shards    []string
}

// Shards the file in the given filepath, saving the shards in the default path for shards
// and returning a Manifest object, which is saved in the default path for manfiests
func ShardFile(filepath string, key *rsa.PrivateKey) (*Manifest, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if err == io.EOF {
			break
		}

		// Compute hash identifier for the shard
		shardContent := append(buffer[:n], filenameHash...) // TODO: Should add more information to the shardContent

		shardHash := crypto.HashSha256(shardContent)
		shardHashString := crypto.HashAsString(shardHash)

		/*
			// Encrypt content of the shard
			encryptedBuffer, err := crypto.Encrypt(buffer[:n], key)
			if err != nil {
				return nil, err
			}
		*/

		err = saveShard(buffer[:n], shardHashString)
		if err != nil {
			return nil, err
		}
		manifest.Shards = append(manifest.Shards, shardHashString)
	}

	manifest.Hash = hashManifest(manifest)

	err = saveManifest(manifest)
	if err != nil {
		log.Fatal(err.Error())
	}

	return manifest, nil
}

func StoreShard(content []byte, hash string) error {
	shardPath := path.Join(DEFAULT_SHARD_PATH, hash)

	file, err := os.Create(shardPath)
	if err != nil {
		return err
	}

	n, _ := file.Write(content)
	fmt.Println(n)
	file.Close()

	return nil
}

func RetrieveShard(hash string) ([]byte, error) {
	return getShard(hash, DEFAULT_SHARD_PATH)
}

// Returns the content of a Shard
func GetTmpShard(hash string) ([]byte, error) {
	return getShard(hash, DEFAULT_TMP_PATH)
}

func RemoveTmpShard(hash string) error {
	shardPath := path.Join(DEFAULT_TMP_PATH, hash)
	err := os.Remove(shardPath)
	return err
}

/*
// Returns if a shard exists
func ShardExists(hash string) bool {
}
*/

// Gets the Manifest object from storage
func GetManifest(hash string) (*Manifest, error) {
	manifestPath := path.Join(DEFAULT_MANIFEST_PATH, hash)
	manifestString, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}
	err = json.Unmarshal(manifestString, manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func SearchManifestFromName(name string) (*Manifest, error) {
	files, err := ioutil.ReadDir(DEFAULT_MANIFEST_PATH)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		manifest, err := GetManifest(file.Name())
		if err != nil {
			return nil, err
		}

		if manifest.Filename == name {
			return manifest, nil
		}
	}

	return &Manifest{}, errors.New("manifest with given name does not exist")
}

/*
// Returns if a Manifest file exists
func ManifestExists(hash string) bool {
}
*/

func getShard(hash string, from string) ([]byte, error) {
	shardPath := path.Join(from, hash)

	content, err := ioutil.ReadFile(shardPath)
	if err != nil {
		return nil, err
	}

	return content, err
}

// Saves the shard's content in the default path for shards
func saveShard(content []byte, hash string) error {
	shardPath := path.Join(DEFAULT_TMP_PATH, hash)

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
	shards := make([]string, len(manifest.Shards))
	copy(shards, manifest.Shards)

	if len(shards)%2 != 0 {
		shards = append(shards, "0")
	}

	for len(shards) > 1 {
		l := len(shards)
		firstHash := shards[l-1]
		secondHash := shards[l-2]

		combinedHash := firstHash + secondHash
		combinedHash = crypto.HashAsString(crypto.HashSha256([]byte(combinedHash)))

		shards = shards[:l-2]
		shards = append(shards, combinedHash)
	}

	return shards[0]
}
