package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"log"
)

// Generates a new RSA key pair
func GenerateRSAKey() (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Gets the already generated RSA key saved in the computer (if exists)
func GetRSAKey() (*rsa.PrivateKey, error) {
	// TODO: Not implemented
	return nil, errors.New("")
}

// TODO: Should be in another file?
func SignMessage(hashedMessage []byte, key *rsa.PrivateKey) ([]byte, error) {
	signature, err := rsa.SignPSS(rand.Reader, key, crypto.SHA256, hashedMessage, nil)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func VerifyMessage(hashedMessage []byte, signature []byte, publKey *rsa.PublicKey) bool {
	err := rsa.VerifyPSS(publKey, crypto.SHA256, hashedMessage, signature, nil)
	if err != nil && err != rsa.ErrVerification {
		log.Fatal(err.Error())
	}
	return err == nil
}
