package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
)

func GenerateRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err.Error())
	}

	return key
}
