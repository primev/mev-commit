package keysstore

import (
	"crypto/ecdh"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primev/mev-commit/p2p/pkg/storage"
)

const (
	aesKeysNS         = "aes/"
	eciesPrivateKeyNS = "ecies/"
	nikePrivateKeyNS  = "nike/"
)

var (
	bidderAesKey = func(bidder common.Address) string {
		return fmt.Sprintf("%s%s", aesKeysNS, bidder)
	}
)

type Store struct {
	mu sync.RWMutex
	st storage.Storage
}

func New(st storage.Storage) *Store {
	return &Store{
		st: st,
	}
}

func (s *Store) SetAESKey(bidder common.Address, key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(bidderAesKey(bidder), key)
}

func (s *Store) GetAESKey(bidder common.Address) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, err := s.st.Get(bidderAesKey(bidder))
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}
	return val, nil
}

func eciesPrivateKeyToBytes(priv *ecies.PrivateKey) []byte {
	return priv.ExportECDSA().D.Bytes()
}

func eciesPrivateKeyFromBytes(data []byte) *ecies.PrivateKey {
	curve := crypto.S256()
	priv := new(ecies.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(data)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(data)
	return priv
}

func (s *Store) SetECIESPrivateKey(key *ecies.PrivateKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(eciesPrivateKeyNS, eciesPrivateKeyToBytes(key))
}

func (s *Store) GetECIESPrivateKey() (*ecies.PrivateKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, err := s.st.Get(eciesPrivateKeyNS)
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return eciesPrivateKeyFromBytes(val), nil
}

func ecdhPrivateKeyToBytes(priv *ecdh.PrivateKey) []byte {
	return priv.Bytes()
}

func ecdhPrivateKeyFromBytes(data []byte) (*ecdh.PrivateKey, error) {
	return ecdh.P256().NewPrivateKey(data)
}

func (s *Store) SetNikePrivateKey(key *ecdh.PrivateKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.st.Put(nikePrivateKeyNS, ecdhPrivateKeyToBytes(key))
}

func (s *Store) GetNikePrivateKey() (*ecdh.PrivateKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, err := s.st.Get(nikePrivateKeyNS)
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return ecdhPrivateKeyFromBytes(val)
}
