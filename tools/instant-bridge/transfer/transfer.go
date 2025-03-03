package transfer

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/keysigner"
)

type Transferer struct {
	mtx               sync.Mutex
	logger            *slog.Logger
	client            *ethclient.Client
	l1ChainID         *big.Int
	settlementChainID *big.Int
	signer            keysigner.KeySigner
	gasTip            *big.Int
	gasFeeCap         *big.Int
}

func NewTransferer(
	logger *slog.Logger,
	client *ethclient.Client,
	l1ChainID *big.Int,
	settlementChainID *big.Int,
	signer keysigner.KeySigner,
	gasTip *big.Int,
	gasFeeCap *big.Int,
) *Transferer {
	return &Transferer{
		logger:            logger,
		client:            client,
		l1ChainID:         l1ChainID,
		settlementChainID: settlementChainID,
		signer:            signer,
		gasTip:            gasTip,
		gasFeeCap:         gasFeeCap,
	}
}

func (t *Transferer) TransferOnSettlement(
	ctx context.Context,
	to common.Address,
	amount *big.Int,
) error {
	// Only one transfer at a time
	t.mtx.Lock()
	defer t.mtx.Unlock()

	nonce, err := t.client.PendingNonceAt(ctx, t.signer.GetAddress())
	if err != nil {
		t.logger.Error("failed to get nonce", "error", err)
		return err
	}
	txData := &types.DynamicFeeTx{
		To:        &to,
		Nonce:     nonce,
		GasFeeCap: t.gasFeeCap,
		GasTipCap: t.gasTip,
		Gas:       21000,
		Value:     amount,
	}

	tx := types.NewTx(txData)

	signedTx, err := t.signer.SignTx(tx, t.settlementChainID)
	if err != nil {
		t.logger.Error("failed to sign tx", "error", err)
		return err
	}

	err = t.client.SendTransaction(ctx, signedTx)
	if err != nil {
		t.logger.Error("failed to send tx", "error", err)
		return err
	}

	r, err := bind.WaitMined(ctx, t.client, signedTx)
	if err != nil {
		t.logger.Error("failed to wait for tx", "error", err)
		return err
	}

	if r.Status != types.ReceiptStatusSuccessful {
		t.logger.Error("tx failed", "status", r.Status)
		return errors.New("tx failed")
	}

	return nil
}

func (t *Transferer) ValidateL1Tx(rawTx string) (*types.Transaction, error) {
	txBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		t.logger.Error("failed to decode tx", "error", err)
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		t.logger.Error("failed to decode tx", "error", err)
		return nil, err
	}

	if tx.ChainId().Cmp(t.l1ChainID) != 0 {
		t.logger.Error("tx has wrong chain ID", "chainID", tx.ChainId(), "expected", t.l1ChainID)
		return nil, errors.New("tx has wrong chain ID")
	}

	if tx.To() == nil {
		t.logger.Error("tx has no recipient")
		return nil, errors.New("tx has no recipient")
	}

	if tx.Value().Sign() <= 0 {
		t.logger.Error("tx has no value")
		return nil, errors.New("tx has no value")
	}

	return tx, nil
}

func (t *Transferer) Sender(tx *types.Transaction) (common.Address, error) {
	return types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
}
