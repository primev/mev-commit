package crypto

import (
	"bytes"
	"testing"
)

func TestGenerateAESKey(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(key) != 32 {
		t.Errorf("Expected key length of 32, got %d", len(key))
	}
}

func TestEncryptWithAESGCM(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	plaintext := []byte("This is a test")
	ciphertext, err := EncryptWithAESGCM(key, plaintext)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Error("Ciphertext should not be equal to plaintext")
	}
}

func TestDecryptWithAESGCM(t *testing.T) {
	key, err := GenerateAESKey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	plaintext := []byte("This is a test")
	ciphertext, err := EncryptWithAESGCM(key, plaintext)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	decryptedText, err := DecryptWithAESGCM(key, ciphertext)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !bytes.Equal(plaintext, decryptedText) {
		t.Errorf("Expected decrypted text to be %s, got %s", plaintext, decryptedText)
	}
}
