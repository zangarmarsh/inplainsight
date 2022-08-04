package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt( plaintext []byte, key []byte ) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	iv := make([]byte, aesgcm.NonceSize())
	full, err := io.ReadFull(rand.Reader, iv)
	if err != nil || full != cap(iv) {
		return "", err
	}

	sealed := aesgcm.Seal( nil, iv, plaintext, nil )

	return base64.RawStdEncoding.EncodeToString(append(iv, sealed...)), nil
}

func Decrypt( ciphertext string, key []byte ) (string, error) {
	bCiphertext, err := base64.RawStdEncoding.DecodeString(ciphertext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	bPlainText, err := aesgcm.Open(nil, bCiphertext[:aesgcm.NonceSize()], bCiphertext[aesgcm.NonceSize():], nil )
	if err != nil {
		return "", err
	}

	return string(bPlainText), nil
}
