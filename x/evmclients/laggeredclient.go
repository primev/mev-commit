package evmclients

import (
	"context"
)

type laggerdL1Client struct {
	EthClient
	amount int
}

func NewLaggerdL1Client(ethClient EthClient, amount int) EthClient {
	return &laggerdL1Client{
		EthClient: ethClient,
		amount:    amount,
	}
}

func (l *laggerdL1Client) BlockNumber(ctx context.Context) (uint64, error) {
	blkNum, err := l.EthClient.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blkNum - uint64(l.amount), nil
}
