package crypto

import (
	"testing"
)

func TestECDHKeyMatching(t *testing.T) {
	// Party A keypair
	skA, pkA, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Party B keypair
	skB, pkB, err := GenerateKeyPairBN254()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// A -> B: pkA;  B -> A: pkB
	// A computes shared = pkB^skA
	sharedA := DeriveSharedKey(skA, pkB)
	// B computes shared = pkA^skB
	sharedB := DeriveSharedKey(skB, pkA)

	// sharedA and sharedB should be the same group element
	if !sharedA.Equal(&sharedB) {
		t.Error("Expected shared keys to match")
	}
}
