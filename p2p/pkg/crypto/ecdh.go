package crypto

import (
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

	// 2) Retrieve the G1 generator (1,2) from the bn254 package
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
	var pub bn254.G1Affine

	// G1Affine.SetBytes returns (int, error).
	// For uncompressed, we expect 96 consumed bytes if successful.
	consumed, err := pub.SetBytes(data)
	if err != nil {
		return nil, err
	}
	if consumed != bn254.SizeOfG1AffineUncompressed {
		return nil, fmt.Errorf("unexpected consumed bytes. got=%d want=%d", consumed, bn254.SizeOfG1AffineUncompressed)
	}

	return &pub, nil
}
