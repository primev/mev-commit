package keykeeper

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
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
	return &BidderKeyKeeper{
		BaseKeyKeeper:   NewBaseKeyKeeper(keysigner),
	}, nil
}

func NewProviderKeyKeeper(keysigner keysigner.KeySigner) (*ProviderKeyKeeper, error) {
	return &ProviderKeyKeeper{
		BaseKeyKeeper:      NewBaseKeyKeeper(keysigner),
	}, nil
}