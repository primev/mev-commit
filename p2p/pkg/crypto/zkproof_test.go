package crypto_test

import (
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/stretchr/testify/require"
)

// TestGenerateAndVerify tests the happy path: generating a valid proof and verifying it.
func TestGenerateAndVerify(t *testing.T) {
	// 1) Generate secret a, b ∈ [0, p-1]
	var a, b fr.Element
	a.SetRandom()
	b.SetRandom()

	// 2) Compute A = g^a, B = g^b
	//    and then C = B^a
	//    bn254.Generators() => we want the g1 generator
	_, _, g1Aff, _ := bn254.Generators()

	// A = g^a
	var A bn254.G1Affine
	var aBigInt big.Int
	a.BigInt(&aBigInt)
	A.ScalarMultiplication(&g1Aff, &aBigInt)

	// B = g^b
	var B bn254.G1Affine
	var bBigInt big.Int
	b.BigInt(&bBigInt)
	B.ScalarMultiplication(&g1Aff, &bBigInt)

	// C = B^a
	var C bn254.G1Affine
	C.ScalarMultiplication(&B, &aBigInt)

	// 3) Generate proof
	context := []byte("test context 1234")
	proof, err := crypto.GenerateOptimizedProof(&a, &A, &B, &C, context)
	require.NoError(t, err, "failed to generate proof")

	// 4) Verify proof (should succeed)
	err = crypto.VerifyOptimizedProof(proof, &A, &B, &C, context)
	require.NoError(t, err, "proof should verify successfully")
}

// TestProofTamperedC tries changing the proof.C to ensure it fails.
func TestProofTamperedC(t *testing.T) {
	var a, b fr.Element
	a.SetRandom()
	b.SetRandom()

	// Setup A, B, C same as above
	_, _, g1Aff, _ := bn254.Generators()

	var A, B, C bn254.G1Affine
	var aBigInt, bBigInt big.Int
	a.BigInt(&aBigInt)
	b.BigInt(&bBigInt)

	A.ScalarMultiplication(&g1Aff, &aBigInt)
	B.ScalarMultiplication(&g1Aff, &bBigInt)
	C.ScalarMultiplication(&B, &aBigInt)

	context := []byte("test context tamper c")

	// Generate proof
	proof, err := crypto.GenerateOptimizedProof(&a, &A, &B, &C, context)
	require.NoError(t, err)

	// Tamper with proof.C
	// e.g. add 1 mod p
	var one fr.Element
	one.SetOne()
	proof.C.Add(&proof.C, &one)

	// Now verify should fail
	err = crypto.VerifyOptimizedProof(proof, &A, &B, &C, context)
	require.Error(t, err, "verification must fail after tampering with c")
}

// TestProofTamperedZ tries changing the proof.Z to ensure it fails.
func TestProofTamperedZ(t *testing.T) {
	var a, b fr.Element
	a.SetRandom()
	b.SetRandom()

	// Setup A, B, C
	_, _, g1Aff, _ := bn254.Generators()
	var A, B, C bn254.G1Affine
	var aBigInt, bBigInt big.Int
	a.BigInt(&aBigInt)
	b.BigInt(&bBigInt)

	A.ScalarMultiplication(&g1Aff, &aBigInt)
	B.ScalarMultiplication(&g1Aff, &bBigInt)
	C.ScalarMultiplication(&B, &aBigInt)

	context := []byte("test context tamper z")

	// Generate proof
	proof, err := crypto.GenerateOptimizedProof(&a, &A, &B, &C, context)
	require.NoError(t, err)

	// Tamper with proof.Z
	// e.g. add 2 mod p
	var two fr.Element
	two.SetUint64(2)
	proof.Z.Add(&proof.Z, &two)

	// Verify should fail
	err = crypto.VerifyOptimizedProof(proof, &A, &B, &C, context)
	require.Error(t, err, "verification must fail after tampering with z")
}

// TestProofTamperedContext tries changing the context to ensure it fails.
func TestProofTamperedContext(t *testing.T) {
	var a, b fr.Element
	a.SetRandom()
	b.SetRandom()

	// Setup A, B, C
	_, _, g1Aff, _ := bn254.Generators()
	var A, B, C bn254.G1Affine
	var aBigInt, bBigInt big.Int
	a.BigInt(&aBigInt)
	b.BigInt(&bBigInt)

	A.ScalarMultiplication(&g1Aff, &aBigInt)
	B.ScalarMultiplication(&g1Aff, &bBigInt)
	C.ScalarMultiplication(&B, &aBigInt)

	// Original context
	context := []byte("test context original")

	// Generate proof w/ original context
	proof, err := crypto.GenerateOptimizedProof(&a, &A, &B, &C, context)
	require.NoError(t, err)

	// Now verify with a different context
	altContext := []byte("test context changed")
	err = crypto.VerifyOptimizedProof(proof, &A, &B, &C, altContext)
	require.Error(t, err, "verification must fail with a different context")
}

// TestEdgeCases tries zero secret key or repeated random trials to ensure no panic or unexpected errors.
func TestEdgeCases(t *testing.T) {
	// a=0, b random
	var a fr.Element
	a.SetZero() // secret key is zero
	var b fr.Element
	b.SetRandom()

	_, _, g1Aff, _ := bn254.Generators()

	// A = g^0 => identity
	var A bn254.G1Affine
	var zeroBigInt big.Int
	A.ScalarMultiplication(&g1Aff, &zeroBigInt)

	// B = g^b
	var B bn254.G1Affine
	var bBigInt big.Int
	b.BigInt(&bBigInt)
	B.ScalarMultiplication(&g1Aff, &bBigInt)

	// C = B^a = B^0 => identity
	var C bn254.G1Affine
	C.ScalarMultiplication(&B, &zeroBigInt)

	context := []byte("edge case zero secret")
	proof, err := crypto.GenerateOptimizedProof(&a, &A, &B, &C, context)
	require.NoError(t, err, "should handle zero secret key")

	err = crypto.VerifyOptimizedProof(proof, &A, &B, &C, context)
	require.NoError(t, err, "zero secret key proof should verify if consistent with A, C = identity")

	// optional: do repeated random pairs
	for i := 0; i < 3; i++ {
		var a2, b2 fr.Element
		a2.SetRandom()
		b2.SetRandom()

		var A2, B2, C2 bn254.G1Affine
		var a2Big, b2Big big.Int
		a2.BigInt(&a2Big)
		b2.BigInt(&b2Big)

		A2.ScalarMultiplication(&g1Aff, &a2Big)
		B2.ScalarMultiplication(&g1Aff, &b2Big)
		C2.ScalarMultiplication(&B2, &a2Big)

		ctx2 := []byte("random trial")
		p2, err2 := crypto.GenerateOptimizedProof(&a2, &A2, &B2, &C2, ctx2)
		require.NoError(t, err2, "generate failed on random trial")

		err2 = crypto.VerifyOptimizedProof(p2, &A2, &B2, &C2, ctx2)
		require.NoError(t, err2, "verify failed on random trial")
	}
}
