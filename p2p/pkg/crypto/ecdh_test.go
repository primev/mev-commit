package crypto

import (
	"testing"
)

func TestECDHKeyMatching(t *testing.T) {
	// Party A keypair
	skA, pkA := GenerateKeyPairBN254()
	// Party B keypair
	skB, pkB := GenerateKeyPairBN254()
	// A -> B: pkA;  B -> A: pkB
	// A computes shared = pkB^skA
	sharedA := DeriveSharedKey(skA, pkB)
	// B computes shared = pkA^skB
	sharedB := DeriveSharedKey(skB, pkA)

	// sharedA and sharedB should be the same group element
	if !sharedA.Equal(sharedB) {
		t.Error("Expected shared keys to match")
	}
}
