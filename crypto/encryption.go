package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"io"
)

func Encrypt(content []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	key := HashSha256(privKey.N.Bytes())

	c, err := aes.NewCipher(key)
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

func Decrypt(encryptedContent []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	key := HashSha256(privKey.N.Bytes())

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedContent) < nonceSize {
		return nil, err
	}

	nonce, encryptedContent := encryptedContent[:nonceSize], encryptedContent[nonceSize:]
	content, err := gcm.Open(nil, nonce, encryptedContent, nil)
	return content, err
}
