package crypto

import (
	"errors"
	"io"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"golang.org/x/crypto/sha3"
)

// For BN254, p ~ 2^254, so let l = 253 bits
const BN254OrderBitLen = 253

// BN254Mask253 is 1<<253, used to truncate the hash to 253 bits
var BN254Mask253 = new(big.Int).Lsh(big.NewInt(1), 253) // = 1 << 253

// Proof holds the final optimized proof with NO ephemeral points:
// we only store the challenge c and response z.
type Proof struct {
	C fr.Element // truncated challenge in [0, p)
	Z fr.Element // response
}

// GenerateOptimizedProof uses the "less-obvious optimization" so the final proof is only (C,Z).
//
//	A = g^a     (provider’s NIKE pubkey)
//	B = g^b     (bidder’s NIKE pubkey)
//	C = B^a     (shared secret)
//
// Steps:
//  1. k in [0, p-1], compute T1 = g^k, T2 = B^k
//  2. c = TruncatedHash(ctx || A || B || C || T1 || T2)
//  3. z = k - a·c (mod p)
//  4. discard T1, T2, return only (c, z).
func GenerateOptimizedProof(
	skA *fr.Element, // a
	pubA, // A = g^a
	pubB, // B = g^b
	sharedC *bn254.G1Affine, // C = B^a
	context []byte,
) (Proof, error) {
	// 1) sample random k in [0, p-1]
	// You can also use fr.Element.SetRandom, but let's demonstrate an explicit random read.
	var k fr.Element
	if _, err := k.SetRandom(); err != nil {
		return Proof{}, err
	}

	// 2) T1 = g^k, T2 = B^k
	//    Using the built-in generator from gnark-crypto
	var G bn254.G1Affine
	// G.SetOne() // G is the conventional generator for bn254
	// var g1Aff bn254.G1Affine
	G.X.SetOne()
	G.Y.SetUint64(2)

	var kBig big.Int
	k.BigInt(&kBig)

	var T1, T2 bn254.G1Affine
	T1.ScalarMultiplication(&G, &kBig)
	T2.ScalarMultiplication(pubB, &kBig)

	// 3) compute truncated challenge c
	c, err := ComputeZKChallenge(context, pubA, pubB, sharedC, &T1, &T2)
	if err != nil {
		return Proof{}, err
	}

	// 4) z = k - a*c (mod p)
	//    z = k + ( - a*c mod p )
	var z fr.Element
	z.Mul(&c, skA).
		Neg(&z).
		Add(&z, &k) // z = k - a*c

	return Proof{C: c, Z: z}, nil
}

// VerifyOptimizedProof re-derives T1, T2:
//
//	T1' = g^z · A^c
//	T2' = B^z · C^c
//	c'  = TruncatedHash(ctx || A || B || C || T1' || T2')
//	accept iff c' == proof.C
func VerifyOptimizedProof(
	proof Proof,
	pubA, // = g^a
	pubB, // = g^b
	sharedC *bn254.G1Affine, // = B^a
	context []byte,
) error {
	// 1) T1' = g^z * A^c
	var G bn254.G1Affine
	G.X.SetOne()
	G.Y.SetUint64(2)

	var zBig, cBig big.Int
	proof.Z.BigInt(&zBig)
	proof.C.BigInt(&cBig)

	var T1p, T2p bn254.G1Affine

	// T1' = g^z
	T1p.ScalarMultiplication(&G, &zBig)
	// A^c
	var Ac bn254.G1Affine
	Ac.ScalarMultiplication(pubA, &cBig)
	// T1' = T1' + A^c  (group addition in G1)
	T1p.Add(&T1p, &Ac)

	// 2) T2' = B^z + C^c
	T2p.ScalarMultiplication(pubB, &zBig)
	var Cc bn254.G1Affine
	Cc.ScalarMultiplication(sharedC, &cBig)
	T2p.Add(&T2p, &Cc)

	// 3) c' = computeZKChallenge(...)
	cPrime, err := ComputeZKChallenge(context, pubA, pubB, sharedC, &T1p, &T2p)
	if err != nil {
		return err
	}

	// 4) check c' == proof.C
	if !cPrime.Equal(&proof.C) {
		return errors.New("invalid proof: mismatch c")
	}
	return nil
}

// ComputeZKChallenge:
//  1. keccak256( context || pubA || pubB || sharedC || a1 || a2 )
//  2. truncate to 253 bits
//  3. parse as fr.Element
func ComputeZKChallenge(
	context []byte,
	pubA, pubB, sharedC, a1, a2 *bn254.G1Affine,
) (fr.Element, error) {
	keccak := sha3.New256()

	// 1) Write context
	_, _ = keccak.Write(context)

	// 2) Write pubA, pubB, sharedC, a1, a2 as 2 x 32 bytes each (X, Y)
	if err := writeAffine(keccak, pubA); err != nil {
		return fr.Element{}, err
	}
	if err := writeAffine(keccak, pubB); err != nil {
		return fr.Element{}, err
	}
	if err := writeAffine(keccak, sharedC); err != nil {
		return fr.Element{}, err
	}
	if err := writeAffine(keccak, a1); err != nil {
		return fr.Element{}, err
	}
	if err := writeAffine(keccak, a2); err != nil {
		return fr.Element{}, err
	}

	// Finalize
	hashBytes := keccak.Sum(nil)

	// Convert to big.Int, truncate to 253 bits
	hashVal := new(big.Int).SetBytes(hashBytes)
	hashVal.Mod(hashVal, BN254Mask253)

	// Convert into fr.Element
	var challenge fr.Element
	challenge.SetBigInt(hashVal)
	return challenge, nil
}

// writeAffine serializes (X,Y) of a bn254.G1Affine in 32-byte big-endian each.
func writeAffine(w io.Writer, p *bn254.G1Affine) error {
	xInt, yInt := AffineToBigIntXY(p)

	xBytes := leftPad32(&xInt)
	if err := writeAll(w, xBytes); err != nil {
		return err
	}

	yBytes := leftPad32(&yInt)
	if err := writeAll(w, yBytes); err != nil {
		return err
	}
	return nil
}

// leftPad32 serializes the big.Int as a 32-byte, big-endian slice.
func leftPad32(x *big.Int) []byte {
	buf := make([]byte, 32)
	tmp := x.Bytes()
	copy(buf[32-len(tmp):], tmp)
	return buf
}

// writeAll is like io.WriteFull, ensures we either write all data or return an error.
func writeAll(w io.Writer, data []byte) error {
	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return io.ErrShortWrite
	}
	return nil
}
