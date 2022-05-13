package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	sampleContent := []byte("This is sample content")

	// Generate Rsa Key for testing
	key, err := GenerateRSAKey(false)
	if err != nil {
		t.Fatal(err.Error())
	}

	encryptedContent, err := Encrypt(sampleContent, key)
	if err != nil {
		t.Fatal("Error encrypting content:", err.Error())
	}

	decryptedContent, err := Decrypt(encryptedContent, key)
	if err != nil {
		t.Fatal("Error decrypting content:", err.Error())
	}

	if !bytes.Equal(sampleContent, decryptedContent) {
		t.Fatalf("Encrypted contet: %s does not match decrypted: %s\n", sampleContent, decryptedContent)
	}
}
