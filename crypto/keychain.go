package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

// Generates a new RSA key pair
func GenerateRSAKey(save bool) (*rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	if save {
		// Save key
		pemdata := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			},
		)

		file, err := os.Create("/home/david/workspace/go_peercloud/.peercloud/privkey.pem")
		if err != nil {
			return nil, err
		}

		file.Write(pemdata)
		file.Close()
	}

	return key, nil
}

// Gets the already generated RSA key saved in the computer (if exists)
func GetRSAKey() (*rsa.PrivateKey, error) {
	pemdata, err := ioutil.ReadFile("/home/david/workspace/go_peercloud/.peercloud/privkey.pem")
	if err != nil {
		return nil, err
	}

	p, _ := pem.Decode(pemdata)
	key, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// === TODO: Should be in another file? ===
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

// =======
