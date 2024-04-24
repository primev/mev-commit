package keykeeper

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primevprotocol/mev-commit/p2p/pkg/keykeeper/keysigner"
)

type KeyKeeper interface {
	SignHash(data []byte) ([]byte, error)
	GetAddress() common.Address
	GetPrivateKey() (*ecdsa.PrivateKey, error)
	ZeroPrivateKey(key *ecdsa.PrivateKey)
}

type BaseKeyKeeper struct {
	KeySigner keysigner.KeySigner
}

type ProviderKeys struct {
	EncryptionPrivateKey *ecies.PrivateKey
	EncryptionPublicKey  *ecies.PublicKey
	NIKEPrivateKey       *ecdh.PrivateKey
	NIKEPublicKey        *ecdh.PublicKey
}

type ProviderKeyKeeper struct {
	*BaseKeyKeeper
	keys               ProviderKeys
	bidderAESKeysMutex *sync.RWMutex
	BiddersAESKeys     map[common.Address][]byte
}

type BidderKeyKeeper struct {
	*BaseKeyKeeper
	AESKey          []byte
	BidHashesToNIKE map[string]*ecdh.PrivateKey
}
