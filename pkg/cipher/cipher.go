package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	dirPerm    = 0755
	filePerm   = 0644
	keyLength  = 16
	secretPath = "./data/AESGCMKey"
)

var (
	key []byte
)

func init() {
	var err error
	key, err = read()
	if err != nil {
		log.Fatalf("[cipher-init-1] could not load key: %v", err)
	}
}

func read() ([]byte, error) {
	keyBytes, err := os.ReadFile(secretPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		keyBytes, err = create()
		if err != nil {
			return nil, err
		}
	}
	buf := make([]byte, keyLength)
	_, err = base64.StdEncoding.Decode(buf, keyBytes)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func create() ([]byte, error) {
	keyBytes, err := generate()
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(filepath.Dir(secretPath), dirPerm)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(keyBytes)))
	base64.StdEncoding.Encode(buf, keyBytes)
	err = os.WriteFile(secretPath, buf, filePerm)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func generate() ([]byte, error) {
	randomBytes := make([]byte, keyLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, fmt.Errorf("[cipher-GenerateRandomKey-1] generating random key failed: %s", err)
	}
	return randomBytes, nil
}

func Encrypt(plaintext string) (string, error) {
	encrypter, err := aes.NewCipher(key)
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

func Decrypt(ciphertext string) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("[cipher-Decrypt-1] decryption failed: %s", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("[cipher-Decrypt-2] decryption failed: %s", err)
	}
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("[cipher-Decrypt-3] decryption failed: %s", err)
	}
	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return nil, fmt.Errorf("[cipher-Decrypt-4] decryption failed: %s", err)
	}
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("[cipher-Decrypt-5] decryption failed: %s", err)
	}
	return plaintext, nil
}
