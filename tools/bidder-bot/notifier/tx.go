package notifier

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Notifier) SelfETHTransfer() (*types.Transaction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	address := b.signer.GetAddress()

	nonce, err := b.l1Client.PendingNonceAt(ctx, address)
	if err != nil {
		b.logger.Error("Failed to get pending nonce", "error", err)
		return nil, err
	}

	chainID, err := b.l1Client.ChainID(ctx)
	if err != nil {
		b.logger.Error("Failed to get network ID", "error", err)
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &address,
		Value:     big.NewInt(7),
		Gas:       1_000_000,
		GasFeeCap: b.gasFeeCap, // TODO: sanity check fees here
		GasTipCap: b.gasTipCap,
	})

	signedTx, err := b.signer.SignTx(tx, chainID)
	if err != nil {
		b.logger.Error("Failed to sign transaction", "error", err)
		return nil, err
	}

	b.logger.Info("Self ETH transfer transaction created and signed", "tx_hash", signedTx.Hash().Hex())

	return signedTx, nil
}
