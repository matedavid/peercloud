package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hashes a given content with the SHA256 hashing function
func HashSha256(content []byte) string {
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:])
}
