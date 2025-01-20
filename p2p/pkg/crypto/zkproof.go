package crypto

import (
    "crypto/sha256"
    "errors"
    "math/big"

    "github.com/consensys/gnark-crypto/ecc/bn254"
    "github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// For BN254, p ~ 2^254, so let l = 253 bits:
const BN254OrderBitLen = 253

// Proof holds the final optimized proof with NO ephemeral points:
// we only store the challenge c and response z.
type Proof struct {
    C fr.Element // truncated challenge in [0, p)
    Z fr.Element // response
}

// GenerateOptimizedProof uses the "less-obvious optimization" so the final proof is only (c,z).
// We have:
//   A = g^a     (provider’s NIKE pubkey)
//   B = g^b     (bidder’s NIKE pubkey)
//   C = B^a     (shared secret)
//
// Steps:
//   1) k in [0, p-1], compute T1 = g^k, T2 = B^k
//   2) c = TruncatedHash(ctx || A || B || C || T1 || T2)
//   3) z = k - a·c (mod p)
//   4) discard T1, T2, return only (c,z).
func GenerateOptimizedProof(
    skA *fr.Element,    // a
    pubA,               // A = g^a
    pubB,               // B = g^b
    sharedC *bn254.G1Affine, // C = B^a
    context []byte,
) (Proof, error) {

    // 1) sample random k
    var k fr.Element
    k.SetRandom() // uniformly random in [0, p-1]
	var kBigInt big.Int
	k.BigInt(&kBigInt)

    // T1 = g^k, T2 = B^k
    var T1, T2 bn254.G1Affine
	_, _, g1Aff, _ := bn254.Generators()
    T1.ScalarMultiplication(&g1Aff, &kBigInt)
    T2.ScalarMultiplication(pubB, &kBigInt)

    // 2) compute truncated challenge c
    c, err := computeTruncatedChallenge(pubA, pubB, sharedC, &T1, &T2, context)
    if err != nil {
        return Proof{}, err
    }

    // 3) z = k - a*c (mod p)
    //    (the article does "k - c*a"; it's equivalent to "k + p - c*a" mod p)
    var z fr.Element
    z.Mul(&c, skA) // z = a*c
    z.Neg(&z)      // z = - (a*c)
    z.Add(&z, &k)  // z = k - a*c

    return Proof{C: c, Z: z}, nil
}

// VerifyOptimizedProof does not receive T1, T2 from the prover. Instead, it re-derives them:
//
//   T1' = g^z · A^c
//   T2' = B^z · C^c
//   c' = TruncatedHash(ctx || A || B || C || T1' || T2')
//   accept iff c' == proof.C
//
// We'll do the group ops typically on-chain via BN254 precompiles.
func VerifyOptimizedProof(
    proof Proof,
    pubA,   // = g^a
    pubB,   // = g^b
    sharedC *bn254.G1Affine, // = B^a
    context []byte,
) error {
    // 1) T1' = g^z * A^c, T2' = B^z * C^c
    //    We'll do this in gnark-crypto (in on-chain code, you'd do precompiles).
    var T1p, T2p bn254.G1Affine

	_, _, g1Aff, _ := bn254.Generators()

    // T1' = g^z
	var zBigInt big.Int
	proof.Z.BigInt(&zBigInt)
    T1p.ScalarMultiplication(&g1Aff, &zBigInt)
    // A^c
	var cBigInt big.Int
	proof.C.BigInt(&cBigInt)
    var Ac bn254.G1Affine
    Ac.ScalarMultiplication(pubA, &cBigInt)
    // T1' = T1' * A^c
    T1p.Add(&T1p, &Ac)

    // T2' = B^z
    T2p.ScalarMultiplication(pubB, &zBigInt)
    // C^c
    var Cc bn254.G1Affine
    Cc.ScalarMultiplication(sharedC, &cBigInt)
    // T2' = T2' * C^c
    T2p.Add(&T2p, &Cc)

    // 2) c' = truncatedHash(ctx || A || B || C || T1' || T2')
    cPrime, err := computeTruncatedChallenge(pubA, pubB, sharedC, &T1p, &T2p, context)
    if err != nil {
        return err
    }

    // 3) check c' == proof.C
    if !cPrime.Equal(&proof.C) {
        return errors.New("invalid proof: mismatch c")
    }

    return nil
}

// computeTruncatedChallenge is the same method used in both generation and verification:
// it concatenates (A, B, C, T1, T2, context), does SHA-256, then truncates to BN254OrderBitLen bits.
func computeTruncatedChallenge(
    A, B, C, T1, T2 *bn254.G1Affine,
    context []byte,
) (fr.Element, error) {

    bufA := A.RawBytes()   // 64 bytes
    bufB := B.RawBytes()   // 64 bytes
    bufC := C.RawBytes()
    bufT1 := T1.RawBytes()
    bufT2 := T2.RawBytes()

    // build the preimage for hashing
    preimage := append(context, bufA[:]...)
    preimage = append(preimage, bufB[:]...)
    preimage = append(preimage, bufC[:]...)
    preimage = append(preimage, bufT1[:]...)
    preimage = append(preimage, bufT2[:]...)

    // compute sha256
    hash := sha256.Sum256(preimage)
    hashVal := new(big.Int).SetBytes(hash[:])

    // bit truncation:  hashVal = hashVal mod 2^BN254OrderBitLen
    maxVal := new(big.Int).Lsh(big.NewInt(1), BN254OrderBitLen) // 1 << 253
    hashVal.Mod(hashVal, maxVal)

    // convert to fr.Element
    var c fr.Element
    c.SetBigInt(hashVal)

	return c, nil
}
