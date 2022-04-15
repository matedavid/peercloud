package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"io"
)

func Encrypt(content []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	privKeyHash := HashSha256(privKey.N.Bytes())

	c, err := aes.NewCipher(privKeyHash)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	encryptedContent := gcm.Seal(nonce, nonce, content, nil)
	return encryptedContent, nil
}
