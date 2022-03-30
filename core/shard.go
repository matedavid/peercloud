package core

const DEFAULT_SHARD_PATH = ".peercloud/tmp/"
const DEFAULT_MANIFEST_PATH = ".peercloud/storage/"

type Manifest struct {
	hash      string
	filename  string
	extension string
	shards    []string
}

// Shards the file in the given filepath, saving the shards in the default path for shards
// and returning a Manifest object, which is saved in the default path for manfiests
func ShardFile(filepath string) *Manifest {
	// TODO:
	// - Get chunks of 1024 bytes from file
	// - Save each shard
	// - Hash manifest and return it
}

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

// Saves the shard's content in the default path for shards
func saveShard(content []byte, hash string) {
}

// Saves a Manifest object in the default path for manifests
func saveManifest(manifest *Manifest) {
}

// Returns the Hash of a given file using a Merkle Tree
func hashManifest(manifest *Manifest) []byte {
}
