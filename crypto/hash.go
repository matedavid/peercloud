package crypto

import (
	"crypto/sha256"
)

// Hashes a given content with the SHA256 hashing function
func HashSha256(content []byte) [32]byte {
	return sha256.Sum256(content)
}
