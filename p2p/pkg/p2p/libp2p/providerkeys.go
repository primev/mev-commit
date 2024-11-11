package libp2p

import (
	"crypto/ecdh"
	"crypto/rand"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
)

type Store interface {
	SetECIESPrivateKey(*ecies.PrivateKey) error
	GetECIESPrivateKey() (*ecies.PrivateKey, error)
	SetNikePrivateKey(*ecdh.PrivateKey) error
	GetNikePrivateKey() (*ecdh.PrivateKey, error)
}

func getOrSetProviderKeys(store Store) (*p2p.Keys, error) {
	nikePublicKey, err := getOrSetNikePrivateKey(store)
	if err != nil {
		return nil, err
	}
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

func getOrSetNikePrivateKey(store Store) (*ecdh.PublicKey, error) {
	prvKey, err := store.GetNikePrivateKey()
	if err != nil {
		return nil, err
	}
	if prvKey == nil {
		prvKey, err = ecdh.P256().GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		err = store.SetNikePrivateKey(prvKey)
		if err != nil {
			return nil, err
		}
	}
	return prvKey.PublicKey(), nil
}
