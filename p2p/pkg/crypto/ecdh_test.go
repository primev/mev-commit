package crypto

import (
	"testing"
)

func TestGenerateKeyPairBN254(t *testing.T) {
	sk, pk, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatalf("GenerateKeyPairBN254() error: %v", err)
	}
	if sk == nil {
		t.Fatal("Expected non-nil secret key")
	}
	if pk == nil {
		t.Fatal("Expected non-nil public key")
	}
	// Check that the public key isn't the identity point.
	// The identity point (in affine) is usually X = 0, Y = 0 (or other representation).
	// This is a simple check to ensure pk is not trivially the identity.
	if pk.X.IsZero() && pk.Y.IsZero() {
		t.Fatal("Public key appears to be the identity point, which is unexpected")
	}
}

func TestECDHKeyMatching(t *testing.T) {
	// Party A keypair
	skA, pkA, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatal(err)
	}
	// Party B keypair
	skB, pkB, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatal(err)
	}

	// A computes shared = pkB^skA
	sharedA := DeriveSharedKey(skA, pkB)
	// B computes shared = pkA^skB
	sharedB := DeriveSharedKey(skB, pkA)

	// sharedA and sharedB should be the same group element
	if !sharedA.Equal(sharedB) {
		t.Error("Expected shared keys to match, but they differ")
	}
}

func TestBN254PublicKeySerialization(t *testing.T) {
	// Generate a key pair.
	_, pk, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatalf("GenerateKeyPairBN254() error: %v", err)
	}

	// Convert pk -> bytes -> pk2
	pubBytes := BN254PublicKeyToBytes(pk)
	pk2, err := BN254PublicKeyFromBytes(pubBytes)
	if err != nil {
		t.Fatalf("BN254PublicKeyFromBytes() error: %v", err)
	}

	// Check that pk2 == pk
	if !pk.Equal(pk2) {
		t.Error("Deserialized public key does not match the original")
	}
}

func TestBN254PublicKeyDeserializationErrors(t *testing.T) {
	// Generate a valid public key
	_, pk, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatalf("GenerateKeyPairBN254() error: %v", err)
	}

	// Convert pk -> bytes
	pubBytes := BN254PublicKeyToBytes(pk)

	// 1) Truncated bytes
	truncated := pubBytes[:len(pubBytes)-1]
	_, err = BN254PublicKeyFromBytes(truncated)
	if err == nil {
		t.Error("Expected error when deserializing truncated bytes, got nil")
	}

	// 2) Completely malformed (e.g., all zeros)
	malformed := make([]byte, len(pubBytes))
	_, err = BN254PublicKeyFromBytes(malformed)
	if err == nil {
		t.Error("Expected error when deserializing malformed bytes, got nil")
	}

	// 3) Extra bytes
	extra := append(pubBytes, 0x00, 0x01, 0x02)
	_, err = BN254PublicKeyFromBytes(extra)
	if err == nil {
		t.Error("Expected error when deserializing with extra bytes, got nil")
	}
}

func TestMultipleECDHKeyAgreements(t *testing.T) {
	// Run multiple times to ensure consistency
	for i := 0; i < 5; i++ {
		// Party A keypair
		skA, pkA, err := GenerateKeyPairBN254()
		if err != nil {
			t.Fatal(err)
		}
		// Party B keypair
		skB, pkB, err := GenerateKeyPairBN254()
		if err != nil {
			t.Fatal(err)
		}

		// A computes shared
		sharedA := DeriveSharedKey(skA, pkB)
		// B computes shared
		sharedB := DeriveSharedKey(skB, pkA)

		// Check
		if !sharedA.Equal(sharedB) {
			t.Errorf("Run %d: Shared keys do not match", i)
		}
	}
}
