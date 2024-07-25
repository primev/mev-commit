package evmclients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type winnerOverrideL1Client struct {
	EthClient
	winners []string
}

func NewWinnerOverrideL1Client(ethClient EthClient, winners []string) *winnerOverrideL1Client {
	return &winnerOverrideL1Client{
		EthClient: ethClient,
		winners:   winners,
	}
}

func (w *winnerOverrideL1Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	hdr, err := w.EthClient.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	idx := number.Int64() % int64(len(w.winners))
	hdr.Extra = []byte(w.winners[idx])

	return hdr, nil
}
