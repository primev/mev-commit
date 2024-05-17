package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func GenerateAESKey() ([]byte, error) {
	aesKey := make([]byte, 32) // AES-256
	_, err := rand.Read(aesKey)
	if err != nil {
		return nil, err
	}
	return aesKey, nil
}

func EncryptWithAESGCM(aesKey, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func DecryptWithAESGCM(aesKey, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := ciphertext[:aesgcm.NonceSize()]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext[aesgcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
