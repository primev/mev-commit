package crypto

import (
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
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
//
//	A = g^a     (provider’s NIKE pubkey)
//	B = g^b     (bidder’s NIKE pubkey)
//	C = B^a     (shared secret)
//
// Steps:
//  1. k in [0, p-1], compute T1 = g^k, T2 = B^k
//  2. c = TruncatedHash(ctx || A || B || C || T1 || T2)
//  3. z = k - a·c (mod p)
//  4. discard T1, T2, return only (c,z).
func GenerateOptimizedProof(
	skA *fr.Element, // a
	pubA, // A = g^a
	pubB, // B = g^b
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
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

	T1.ScalarMultiplication(&g1Aff, &kBigInt)
	T2.ScalarMultiplication(pubB, &kBigInt)

	// 2) compute truncated challenge c
	// c, err := computeTruncatedChallenge(pubA, pubB, sharedC, &T1, &T2, context)
	// if err != nil {
	// 	return Proof{}, err
	// }
	c, err := ComputeZKChallenge(context, pubA, pubB, sharedC, &T1, &T2)
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
//	T1' = g^z · A^c
//	T2' = B^z · C^c
//	c' = TruncatedHash(ctx || A || B || C || T1' || T2')
//	accept iff c' == proof.C
//
// We'll do the group ops typically on-chain via BN254 precompiles.
func VerifyOptimizedProof(
	proof Proof,
	pubA, // = g^a
	pubB, // = g^b
	sharedC *bn254.G1Affine, // = B^a
	context []byte,
) error {
	// 1) T1' = g^z * A^c, T2' = B^z * C^c
	//    We'll do this in gnark-crypto (in on-chain code, you'd do precompiles).
	var T1p, T2p bn254.G1Affine

	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

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
	// cPrime, err := computeTruncatedChallenge(pubA, pubB, sharedC, &T1p, &T2p, context)
	// if err != nil {
	// 	return err
	// }

	cPrime, err := ComputeZKChallenge(context, pubA, pubB, sharedC, &T1p, &T2p)
	if err != nil {
		return err
	}
	// 3) check c' == proof.C
	if !cPrime.Equal(&proof.C) {
		return errors.New("invalid proof: mismatch c")
	}

	return nil
}

var BN254Mask253 = new(big.Int).Lsh(big.NewInt(1), 253) // 1 << 253

// computeTruncatedChallenge is the same method used in both generation and verification:
// it concatenates (A, B, C, T1, T2, context), does SHA-256, then truncates to BN254OrderBitLen bits.
func computeTruncatedChallenge(
	A, B, C, T1, T2 *bn254.G1Affine,
	context []byte,
) (fr.Element, error) {

	bufA := A.RawBytes() // 64 bytes
	bufB := B.RawBytes() // 64 bytes
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

func ComputeZKChallenge(
	contextHash []byte,
	providerPub *bn254.G1Affine,
	bidPub *bn254.G1Affine,
	sharedSec *bn254.G1Affine,
	a *bn254.G1Affine,
	a2 *bn254.G1Affine,
	// providerPubX, providerPubY *big.Int,
	// bidPubX, bidPubY *big.Int,
	// sharedSecX, sharedSecY *big.Int,
	// aX, aY *big.Int,
	// aX2, aY2 *big.Int,
) (fr.Element, error) {
	ctxHash := crypto.Keccak256Hash(contextHash)
	ctxHashBytes := ctxHash.Bytes()

	var providerPubX, providerPubY big.Int
	providerPub.X.BigInt(&providerPubX)
	providerPub.Y.BigInt(&providerPubY)

	var bidPubX, bidPubY big.Int
	bidPub.X.BigInt(&bidPubX)
	bidPub.Y.BigInt(&bidPubY)

	var sharedSecX, sharedSecY big.Int
	sharedSec.X.BigInt(&sharedSecX)
	sharedSec.Y.BigInt(&sharedSecY)

	var aX, aY big.Int
	a.X.BigInt(&aX)
	a.Y.BigInt(&aY)

	var aX2, aY2 big.Int
	a2.X.BigInt(&aX2)
	a2.Y.BigInt(&aY2)

	// 1) Flatten in the same order as `abi.encodePacked`.
	//    We'll manually append each big.Int as 32-byte big-endian,
	//    then the final keccak256 is your computedChallenge.
	var buf []byte

	// a) ZK_CONTEXT_HASH is already 32 bytes
	buf = append(buf, ctxHashBytes[:]...)

	// b) For each big.Int, we 0-pad to 32 bytes big-endian
	buf = append(buf, leftPad32(&providerPubX)...)
	buf = append(buf, leftPad32(&providerPubY)...)
	buf = append(buf, leftPad32(&bidPubX)...)
	buf = append(buf, leftPad32(&bidPubY)...)
	buf = append(buf, leftPad32(&sharedSecX)...)
	buf = append(buf, leftPad32(&sharedSecY)...)
	buf = append(buf, leftPad32(&aX)...)
	buf = append(buf, leftPad32(&aY)...)
	buf = append(buf, leftPad32(&aX2)...)
	buf = append(buf, leftPad32(&aY2)...)

	// 2) keccak256
	hash := crypto.Keccak256Hash(buf)
	hashVal := new(big.Int).SetBytes(hash.Bytes())

	// 3) bitmask to 253 bits
	hashVal.Mod(hashVal, BN254Mask253)

	// 4) return as fr.Element
	var challenge fr.Element
	challenge.SetBigInt(hashVal)
	return challenge, nil
}

// leftPad32 serializes the big.Int as a 32-byte, big-endian slice.
func leftPad32(x *big.Int) []byte {
	buf := make([]byte, 32)
	tmp := x.Bytes()
	copy(buf[32-len(tmp):], tmp)
	return buf
}
