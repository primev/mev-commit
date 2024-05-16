package keykeeper

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/keykeeper/keysigner"
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

type ProviderKeyKeeper struct {
	*BaseKeyKeeper
}

type BidderKeyKeeper struct {
	*BaseKeyKeeper
}
