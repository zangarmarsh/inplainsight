package cryptography

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"testing"
)

func TestEncryptionAndDecryption(t *testing.T) {
	secret := []byte("c")
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		t.Fatal(err)
	}

	ciph, err := Encrypt(secret, key)
	if err != nil {
		t.Fatal(err)
	}

	decr, err := Decrypt(ciph, key)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(decr, secret) != 0 {
		t.Fatalf("Decrypted %s should be equal to %s", decr, secret)
	}
}

func TestHMAC(t *testing.T) {
	content := []byte("Lorem ipsum dolor sit amet")
	key := make([]byte, sha256.BlockSize)

	if bytes.Compare(HMAC(content, key), HMAC(content,key)) != 0 {
		t.Fatal("generated hmacs are not equals")
	}
}
