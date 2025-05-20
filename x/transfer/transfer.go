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
)

type EthClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
}

type Signer interface {
	GetAddress() common.Address
	SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)
}

type Transferer struct {
	mtx       sync.Mutex
	logger    *slog.Logger
	client    EthClient
	signer    Signer
	gasTip    *big.Int
	gasFeeCap *big.Int
}

func NewTransferer(
	logger *slog.Logger,
	client EthClient,
	signer Signer,
	gasTip *big.Int,
	gasFeeCap *big.Int,
) *Transferer {
	return &Transferer{
		logger:    logger,
		client:    client,
		signer:    signer,
		gasTip:    gasTip,
		gasFeeCap: gasFeeCap,
	}
}

func (t *Transferer) Transfer(
	ctx context.Context,
	to common.Address,
	chainID *big.Int,
	amount *big.Int,
) error {
	// Only one transfer at a time
	t.mtx.Lock()
	defer t.mtx.Unlock()

	if to == (common.Address{}) {
		t.logger.Error("invalid address")
		return errors.New("invalid address")
	}

	if amount.Sign() <= 0 {
		t.logger.Error("invalid amount")
		return errors.New("invalid amount")
	}

	if chainID.Cmp(big.NewInt(0)) <= 0 {
		t.logger.Error("invalid chain ID")
		return errors.New("invalid chain ID")
	}

	// Check if the account is a contract
	code, err := t.client.CodeAt(ctx, to, nil)
	if err != nil {
		t.logger.Error("failed to get code", "error", err)
		return err
	}

	if len(code) > 0 {
		t.logger.Error("address is a contract")
		return errors.New("address is a contract")
	}

	nonce, err := t.client.PendingNonceAt(ctx, t.signer.GetAddress())
	if err != nil {
		t.logger.Error("failed to get nonce", "error", err)
		return err
	}
	txData := &types.DynamicFeeTx{
		To:        &to,
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: t.gasFeeCap,
		GasTipCap: t.gasTip,
		Gas:       21000,
		Value:     amount,
	}

	tx := types.NewTx(txData)

	signedTx, err := t.signer.SignTx(tx, chainID)
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

func (t *Transferer) ValidateTx(rawTx string, chainID *big.Int) (*types.Transaction, error) {
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

	if tx.ChainId().Cmp(chainID) != 0 {
		t.logger.Error("tx has wrong chain ID", "chainID", tx.ChainId(), "expected", chainID)
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
