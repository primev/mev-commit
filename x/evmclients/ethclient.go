package evmclients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}
