package cryptography

import (
	"crypto/rand"
	"io"
	"testing"
)

func TestEncryption(t *testing.T) {
	secret := "c"
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		t.Fatal(err)
	}

	ciph, _ := encrypt([]byte(secret), key)
	_, err = decrypt(ciph, key)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDecryption(t *testing.T) {
	secret := "c"
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		t.Fatal(err)
	}

	ciph, _ := encrypt([]byte(secret), key)
	decr, err := decrypt(ciph, key)
	if err != nil {
		t.Fatal(err)
	}

	if decr != secret {
		t.Fatalf("Decrypted %s should be equal to %s", decr, secret)
	}
}
