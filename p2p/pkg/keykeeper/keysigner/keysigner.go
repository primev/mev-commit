package keysigner

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type KeySigner interface {
	fmt.Stringer

	SignHash(data []byte) ([]byte, error)
	SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
	GetAddress() common.Address
	GetPrivateKey() (*ecdsa.PrivateKey, error)
	ZeroPrivateKey(key *ecdsa.PrivateKey)
}
