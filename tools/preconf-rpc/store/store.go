package store

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
)

type rpcstore struct {
}

func (s *rpcstore) StorePreconfirmedTransaction(
	ctx context.Context,
	blockNumber int64,
	txn *types.Transaction,
	commitments []*bidderapiv1.Commitment,
) error {
	return nil
}

func (s *rpcstore) GetPreconfirmedTransaction(
	ctx context.Context,
	txnHash string,
) (*types.Transaction, []*bidderapiv1.Commitment, error) {
	return nil, nil, nil
}
