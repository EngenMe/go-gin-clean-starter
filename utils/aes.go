package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// KEY represents a hardcoded AES encryption key used for encryption and decryption operations.
var (
	KEY = "00112233445566778899aabb00112233445566778899aabb"
)

// AESEncrypt encrypts a given plaintext string using AES encryption and returns the encrypted string or an error.
func AESEncrypt(stringToEncrypt string) (encryptedString string, err error) {
	key, err := hex.DecodeString(KEY)
	if err != nil {
		return "", err
	}
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// AESDecrypt decrypts a given AES-encrypted string using a keyed cipher and returns the plaintext or an error.
func AESDecrypt(encryptedString string) (decryptedString string, err error) {
	defer func() {
		if r := recover(); r != nil {
			decryptedString = ""
			err = errors.New("ciphertext too short")
		}
	}()

	key, err := hex.DecodeString(KEY)
	if err != nil {
		return "", errors.New("error in decoding key")
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", errors.New("error in decoding encrypted string")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	if len(enc) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("error decrypting")
	}

	return string(plaintext), nil
}
