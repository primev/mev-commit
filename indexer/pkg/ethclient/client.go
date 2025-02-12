package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/indexer/pkg/store"
)

type EthereumClient interface {
	GetBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (*big.Int, error)
	TxReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error)
	AccountBalances(ctx context.Context, addresses []common.Address, blockNumber uint64) ([]store.AccountBalance, error)
}
