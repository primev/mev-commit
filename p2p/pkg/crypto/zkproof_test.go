package crypto_test

import (
	"fmt"
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
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

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
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

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
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)
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
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)
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

	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

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

func TestFixedPublicKeys(t *testing.T) {
	//
	// 1) Hard-code: a=1, b=1 => pubA=(1,2), pubB=(1,2), sharedC=(1,2)
	//
	var skA, skB fr.Element
	skA.SetOne() // a=1
	skB.SetOne() // b=1

	// The "public keys" are all (1,2). We won't even call GenerateKeyPairBN254,
	// because we are forcing them. In your proof.go, (1,2) is used as the generator.
	// So A=(1,2), B=(1,2), C=(1,2).
	pubA := makeAffinePoint(1, 2)
	pubB := makeAffinePoint(1, 2)
	sharedC := makeAffinePoint(1, 2)

	// This context is typically hashed in with T1, T2
	context := []byte("mev-commit opening, mainnet, v1.0")

	foundAny := false
	for kVal := 1; kVal <= 30; kVal++ {
		// We'll manually produce T1 = g^k, T2 = B^k, then c = H(...), z = k - a*c
		// If it verifies, we print out c,z
		proof, ok := generateProofWithFixedK(skA, pubA, pubB, sharedC, context, kVal)
		if !ok {
			continue
		}
		// Try verifying
		err := crypto.VerifyOptimizedProof(proof, pubA, pubB, sharedC, context)
		if err == nil {
			fmt.Println("SUCCESS: ephemeral proof with k =", kVal)
			fmt.Printf(" c = %s\n", proof.C.String())
			fmt.Printf(" z = %s\n", proof.Z.String())
			fmt.Println("So your zkProof array is: [1, 2, 1, 2, 1, 2, c, z] => ")
			fmt.Printf("  [1, 2, 1, 2, 1, 2, %s, %s]\n", proof.C.String(), proof.Z.String())
			foundAny = true
			break
		}
	}
	if !foundAny {
		t.Error("No ephemeral k in [1..30] produced a valid proof with A=B=C=(1,2).")
	} else {
		t.Log("Test complete: at least one ephemeral proof was found.")
	}
}

// makeAffinePoint(1,2)
func makeAffinePoint(x, y uint64) *bn254.G1Affine {
	var p bn254.G1Affine
	p.X.SetUint64(x)
	p.Y.SetUint64(y)
	return &p
}

// generateProofWithFixedK is basically "GenerateOptimizedProof" except we forcibly set k = kVal
// rather than randomizing it. We replicate the logic: T1 = g^k, T2 = B^k, c=H(...), z=k-a*c
func generateProofWithFixedK(
	skA fr.Element, // a=1
	pubA, pubB, sharedC *bn254.G1Affine, // (1,2)
	context []byte,
	kVal int,
) (crypto.Proof, bool) {

	var proof crypto.Proof

	// (1) T1 = g^k, T2 = B^k => but B=(1,2), so T2 = (1,2)^k
	g := makeAffinePoint(1, 2)
	var T1, T2 bn254.G1Affine

	// convert kVal -> fr, then -> big.Int for gnark-crypto
	var kEl fr.Element
	kEl.SetUint64(uint64(kVal))

	var kBig big.Int
	kEl.BigInt(&kBig)

	// T1 = (1,2)^k
	T1.ScalarMultiplication(g, &kBig)
	// T2 = B^k = (1,2)^k
	T2.ScalarMultiplication(pubB, &kBig)

	// (2) c = truncatedHash(pubA, pubB, sharedC, T1, T2, context)
	c, err := crypto.ComputeZKChallenge(context, pubA, pubB, sharedC, &T1, &T2)
	if err != nil {
		return proof, false
	}

	// (3) z = k - a*c
	var z fr.Element
	z.Mul(&c, &skA) // z = a*c => 1*c
	z.Neg(&z)       // z = -c
	z.Add(&z, &kEl) // z = kVal - c
	proof.C = c
	proof.Z = z

	return proof, true
}
