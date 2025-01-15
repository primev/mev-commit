package keysstore

import (
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
)

const (
	aesKeysNS         = "aes/"
	eciesPrivateKeyNS = "ecies/"
	bn254PrivateKeyNS = "bn254-sk/"
	bn254PublicKeyNS  = "bn254-pk/"
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

// bn254PrivateKeyToBytes flattens a BN254 fr.Element into 32 bytes (big-endian regular form).
func bn254PrivateKeyToBytes(sk *fr.Element) []byte {
	// sk.Bytes() => returns [32]byte in big-endian *regular* (non-Montgomery) form
	arr := sk.Bytes()
	return arr[:]
}

// bn254PrivateKeyFromBytes interprets data as a 32-byte big-endian integer,
// sets the fr.Element (into Montgomery form internally), and returns it.
func bn254PrivateKeyFromBytes(data []byte) (*fr.Element, error) {
	if len(data) != 32 {
		return nil, errors.New("invalid BN254 private key length (must be 32 bytes)")
	}
	var sk fr.Element
	// SetBytes interprets data as big-endian and puts it in Montgomery form internally
	sk.SetBytes(data)
	return &sk, nil
}

func (s *Store) SetBN254PrivateKey(sk *fr.Element) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw := bn254PrivateKeyToBytes(sk)
	return s.st.Put(bn254PrivateKeyNS, raw)
}

func (s *Store) GetBN254PrivateKey() (*fr.Element, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	raw, err := s.st.Get(bn254PrivateKeyNS)
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return bn254PrivateKeyFromBytes(raw)
}

func (s *Store) SetBN254PublicKey(pub *bn254.G1Affine) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw := p2pcrypto.BN254PublicKeyToBytes(pub) // 96 bytes uncompressed
	return s.st.Put(bn254PublicKeyNS, raw)
}

func (s *Store) GetBN254PublicKey() (*bn254.G1Affine, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	raw, err := s.st.Get(bn254PublicKeyNS)
	switch {
	case errors.Is(err, storage.ErrKeyNotFound):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return p2pcrypto.BN254PublicKeyFromBytes(raw)
}
