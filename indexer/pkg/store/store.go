package store

import (
	"context"
	"math/big"
)

type Storage interface {
	IndexBlock(ctx context.Context, block *IndexBlock) error
	IndexTransactions(ctx context.Context, transactions []*IndexTransaction) error
	GetLastIndexedBlock(ctx context.Context, direction string) (*big.Int, error)
	GetAddresses(ctx context.Context) ([]string, error)
	IndexAccountBalances(ctx context.Context, accountBalances []AccountBalance) error
	CreateIndices(ctx context.Context) error
	Close(ctx context.Context) error
}
