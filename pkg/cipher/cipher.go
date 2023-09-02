package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(key string, plaintext string) (string, error) {
	encrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("[cipher-Encrypt-1] encryption failed: %s", err)
	}
	gcm, err := cipher.NewGCM(encrypter)
	if err != nil {
		return "", fmt.Errorf("[cipher-Encrypt-2] encryption failed: %s", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", fmt.Errorf("[cipher-Encrypt-3] encryption failed: %s", err)
	}
	return base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}

func Decrypt(key string, ciphertext string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("[cipher-Decrypt-1] decryption failed: %s", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("[cipher-Decrypt-2] decryption failed: %s", err)
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("[cipher-Decrypt-3] decryption failed: %s", err)
	}
	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", fmt.Errorf("[cipher-Decrypt-4] decryption failed: %s", err)
	}
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", fmt.Errorf("[cipher-Decrypt-5] decryption failed: %s", err)
	}
	return string(plaintext), nil
}

func GenerateRandomKey(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("[cipher-GenerateRandomKey-1] generating random key failed: %s", err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
