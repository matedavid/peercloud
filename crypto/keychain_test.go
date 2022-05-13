package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	_, err := GenerateRSAKey(false)
	if err != nil {
		t.Fatal("Error generating Private Key:", err.Error())
	}
}

func TestSaveOpenKey(t *testing.T) {
	key, err := GenerateRSAKey(true)
	if err != nil {
		t.Fatal("Error generating Private Key:", err.Error())
	}

	openKey, err := GetRSAKey()
	if err != nil {
		t.Fatal("Error opening Private Key:", err.Error())
	}

	if !bytes.Equal(key.N.Bytes(), openKey.N.Bytes()) {
		t.Fatal("Generated and opened key do not match")
	}

}
