package crypto

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// GenerateKeyPairBN254 returns a BN254 private key (fr.Element) and the
// corresponding public key in G1 affine form.
func GenerateKeyPairBN254() (sk *fr.Element, pk *bn254.G1Affine, err error) {
	// 1) Generate random secret in Fr
	sk = &fr.Element{}
	_, err = sk.SetRandom()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate random secret: %w", err)
	}

	// 2) Set the G1 generator (1,2) (same is used in PreconfManager.sol)
	var g1Aff bn254.G1Affine
	g1Aff.X.SetOne()
	g1Aff.Y.SetUint64(2)

	// 3) Convert sk -> big.Int to call ScalarMultiplication
	var skBigInt big.Int
	sk.BigInt(&skBigInt)

	// 4) pk = g1Aff^sk
	pk = &bn254.G1Affine{}
	pk.ScalarMultiplication(&g1Aff, &skBigInt)

	return sk, pk, nil
}

// DeriveSharedKey does pkB^skA in BN254 G1 (ECDH-style).
func DeriveSharedKey(skA *fr.Element, pkB *bn254.G1Affine) *bn254.G1Affine {
	var skABig big.Int
	skA.BigInt(&skABig)

	var shared bn254.G1Affine
	shared.ScalarMultiplication(pkB, &skABig)
	return &shared
}

func BN254PublicKeyToBytes(pub *bn254.G1Affine) []byte {
	raw := pub.RawBytes() // [96]byte uncompressed
	return raw[:]
}

func BN254PublicKeyFromBytes(data []byte) (*bn254.G1Affine, error) {
	// 1) Check total length strictly
	if len(data) != bn254.SizeOfG1AffineUncompressed {
		return nil, fmt.Errorf("invalid G1 bytes: expected %d bytes, got %d",
			bn254.SizeOfG1AffineUncompressed, len(data))
	}

	var pub bn254.G1Affine
	consumed, err := pub.SetBytes(data)
	if err != nil {
		return nil, fmt.Errorf("SetBytes error: %w", err)
	}

	// 2) Ensure exactly 96 bytes were consumed
	if consumed != bn254.SizeOfG1AffineUncompressed {
		return nil, fmt.Errorf("unexpected consumed bytes. got=%d want=%d",
			consumed, bn254.SizeOfG1AffineUncompressed)
	}

	// 3) Optionally disallow the identity point
	//    Some ECDH setups consider the identity to be invalid as a public key.
	//    If pub is the identity, pub.IsInfinity() will be true.
	if pub.IsInfinity() {
		return nil, fmt.Errorf("invalid G1 point: found point at infinity")
	}

	// 4) Optional: Check that the point is on the curve
	if !pub.IsOnCurve() {
		return nil, fmt.Errorf("invalid G1 point: not on curve")
	}

	return &pub, nil
}

// BN254PrivateKeyToBytes flattens a BN254 fr.Element into 32 bytes (big-endian regular form).
func BN254PrivateKeyToBytes(sk *fr.Element) []byte {
	// sk.Bytes() => returns [32]byte in big-endian *regular* (non-Montgomery) form
	arr := sk.Bytes()
	return arr[:]
}

// BN254PrivateKeyFromBytes interprets data as a 32-byte big-endian integer,
// sets the fr.Element (into Montgomery form internally), and returns it.
func BN254PrivateKeyFromBytes(data []byte) (*fr.Element, error) {
	if len(data) != 32 {
		return nil, errors.New("invalid BN254 private key length (must be 32 bytes)")
	}
	var sk fr.Element
	// SetBytes interprets data as big-endian and puts it in Montgomery form internally
	sk.SetBytes(data)
	return &sk, nil
}

func AffineToBigIntXY(point *bn254.G1Affine) (big.Int, big.Int) {
	var x, y big.Int
	point.X.BigInt(&x)
	point.Y.BigInt(&y)
	return x, y
}
