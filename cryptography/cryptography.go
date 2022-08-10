package cryptography

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encrypt Takes in input a plaintext with variable length and a 32 bytes key and
// returns an encrypted key with the following structure:
// | PlainText | EncryptedContent |
// Where PlainText is a string(16) of plaintext containing the IV and
// EncryptedContent has a ratio of 1:1 with the relative decryption and is composed by
// | string(32) the HMAC | string[N = EncryptedContnent - string(32)(HMAC Size) ] the secret
func Encrypt( plaintext []byte, key []byte ) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	full, err := io.ReadFull(rand.Reader, iv)
	if err != nil || full != cap(iv) {
		return nil, err
	}

	hashedKey := HMAC(plaintext, key)
	plaintext = append(hashedKey, plaintext...)

	aesctr := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plaintext))
	aesctr.XORKeyStream(ciphertext, plaintext)
	output := append(iv, ciphertext...)

	encoded := base64.RawStdEncoding.EncodeToString(output)

	return []byte(encoded), nil
}

// Decrypt - See Encrypt for the ciphertext structure
func Decrypt( ciphertext []byte, key []byte ) ([]byte, error) {
	// Add a minimum length check
	decoded, err := base64.RawStdEncoding.DecodeString(string(ciphertext))
	ciphertext = decoded

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	aesctr := cipher.NewCTR(block, iv)
	ciphertext = ciphertext[aes.BlockSize:]

	decrypted := make([]byte, len(ciphertext))
	aesctr.XORKeyStream(decrypted, ciphertext)

	keyedHash := decrypted[:sha256.Size]
	content := decrypted[sha256.Size:]

	if bytes.Compare(keyedHash, HMAC(content, key)) != 0 {
		fmt.Printf("%20s: %+v\n%20s: %+v\n\n", "hmac", keyedHash, "extracted", HMAC(content, key))
		return nil, errors.New("decryption failed because of content integrity check: hmac are not equals")
	}

	return content, nil
}

func DeriveEncryptionKeysFromPassword(password string) (contentEncryptionKey []byte, headerEncryptionKey []byte) {
	sha512 := crypto.SHA512.New()
	hashedPassword := sha512.Sum([]byte(password))
	headerEncryptionKey = hashedPassword[32:]
	contentEncryptionKey = hashedPassword[:32]

	return
}

func HMAC(content []byte, key []byte) []byte {
	hashFunc := hmac.New(sha256.New, key)
	hashFunc.Write(content)
	return hashFunc.Sum(nil)
}