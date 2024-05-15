package keykeeper

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primev/mev-commit/p2p/pkg/keykeeper/keysigner"
)

// NewBaseKeyKeeper creates a new BaseKeyKeeper.
func NewBaseKeyKeeper(keysigner keysigner.KeySigner) *BaseKeyKeeper {
	return &BaseKeyKeeper{KeySigner: keysigner}
}

func (bkk *BaseKeyKeeper) SignHash(data []byte) ([]byte, error) {
	return bkk.KeySigner.SignHash(data)
}

func (bkk *BaseKeyKeeper) GetAddress() common.Address {
	return bkk.KeySigner.GetAddress()
}

func (bkk *BaseKeyKeeper) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return bkk.KeySigner.GetPrivateKey()
}

func (bkk *BaseKeyKeeper) ZeroPrivateKey(key *ecdsa.PrivateKey) {
	bkk.KeySigner.ZeroPrivateKey(key)
}

func NewBidderKeyKeeper(keysigner keysigner.KeySigner) (*BidderKeyKeeper, error) {
	aesKey, err := generateAESKey()
	if err != nil {
		return nil, err
	}

	bidHashesToNIKE := make(map[string]*ecdh.PrivateKey)

	return &BidderKeyKeeper{
		BaseKeyKeeper:   NewBaseKeyKeeper(keysigner),
		AESKey:          aesKey,
		BidHashesToNIKE: bidHashesToNIKE,
	}, nil
}

func (bkk *BidderKeyKeeper) GenerateNIKEKeys(bidHash []byte) (*ecdh.PublicKey, error) {
	nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	nikePublicKey := nikePrivateKey.PublicKey()
	bkk.BidHashesToNIKE[hex.EncodeToString(bidHash)] = nikePrivateKey
	return nikePublicKey, nil
}

func NewProviderKeyKeeper(keysigner keysigner.KeySigner) (*ProviderKeyKeeper, error) {
	biddersAESKeys := make(map[common.Address][]byte)

	encryptionPrivateKey, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		return nil, err
	}

	nikePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ProviderKeyKeeper{
		BaseKeyKeeper:      NewBaseKeyKeeper(keysigner),
		BiddersAESKeys:     biddersAESKeys,
		bidderAESKeysMutex: &sync.RWMutex{},
		keys: ProviderKeys{
			EncryptionPrivateKey: encryptionPrivateKey,
			EncryptionPublicKey:  &encryptionPrivateKey.PublicKey,
			NIKEPrivateKey:       nikePrivateKey,
			NIKEPublicKey:        nikePrivateKey.PublicKey(),
		},
	}, nil
}

func (pkk *ProviderKeyKeeper) GetNIKEPublicKey() *ecdh.PublicKey {
	return pkk.keys.NIKEPublicKey
}

func (pkk *ProviderKeyKeeper) GetECIESPublicKey() *ecies.PublicKey {
	return pkk.keys.EncryptionPublicKey
}

func (pkk *ProviderKeyKeeper) DecryptWithECIES(message []byte) ([]byte, error) {
	return pkk.keys.EncryptionPrivateKey.Decrypt(message, nil, nil)
}

func (pkk *ProviderKeyKeeper) GetNIKEPrivateKey() *ecdh.PrivateKey {
	return pkk.keys.NIKEPrivateKey
}

func (pkk *ProviderKeyKeeper) SetAESKey(bidder common.Address, key []byte) {
	pkk.bidderAESKeysMutex.Lock()
	defer pkk.bidderAESKeysMutex.Unlock()
	pkk.BiddersAESKeys[bidder] = key
}

func (pkk *ProviderKeyKeeper) GetAESKey(bidder common.Address) ([]byte, bool) {
	pkk.bidderAESKeysMutex.RLock()
	defer pkk.bidderAESKeysMutex.RUnlock()
	key, exists := pkk.BiddersAESKeys[bidder]
	return key, exists
}
