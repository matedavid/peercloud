package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hashes a given content with the SHA256 hashing function
func HashSha256(content []byte) []byte {
	hash := sha256.Sum256(content)
	return hash[:]
}

func HashAsString(hash []byte) string {
	return hex.EncodeToString(hash[:])
}
