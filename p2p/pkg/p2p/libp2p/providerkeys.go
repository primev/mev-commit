package libp2p

import (
	"crypto/rand"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
)

type Store interface {
	SetECIESPrivateKey(*ecies.PrivateKey) error
	GetECIESPrivateKey() (*ecies.PrivateKey, error)
	SetBN254PrivateKey(*fr.Element) error
	GetBN254PrivateKey() (*fr.Element, error)
	SetBN254PublicKey(*bn254.G1Affine) error
	GetBN254PublicKey() (*bn254.G1Affine, error)
}

func getOrSetProviderKeys(store Store) (*p2p.Keys, error) {
	nikePublicKey, err := getOrSetECDHPublicKey(store)
	if err != nil {
		return nil, err
	}
	fmt.Println("nikePublicKey: ", nikePublicKey)
	eciesPublicKey, err := getOrSetECIESPublicKey(store)
	if err != nil {
		return nil, err
	}
	providerKeys := &p2p.Keys{
		NIKEPublicKey: nikePublicKey,
		PKEPublicKey:  eciesPublicKey,
	}
	return providerKeys, nil
}

func getOrSetECIESPublicKey(store Store) (*ecies.PublicKey, error) {
	prvKey, err := store.GetECIESPrivateKey()
	if err != nil {
		return nil, err
	}
	if prvKey == nil {
		prvKey, err = ecies.GenerateKey(rand.Reader, crypto.S256(), nil)
		if err != nil {
			return nil, err
		}
		err = store.SetECIESPrivateKey(prvKey)
		if err != nil {
			return nil, err
		}
	}
	return &prvKey.PublicKey, nil
}

func getOrSetECDHPublicKey(store Store) (*bn254.G1Affine, error) {
	pk, err := store.GetBN254PublicKey()
	if err != nil {
		return nil, err
	}

	if pk == nil {
		sk, pk, err := p2pcrypto.GenerateKeyPairBN254()
		if err != nil {
			return nil, err
		}
		err = store.SetBN254PrivateKey(&sk)
		if err != nil {
			return nil, err
		}
		err = store.SetBN254PublicKey(&pk)
		if err != nil {
			return nil, err
		}
		fmt.Println("pk: ", pk)
		fmt.Println("sk: ", sk)
	}

	return pk, nil
}
