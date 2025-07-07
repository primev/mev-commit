package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/keysigner"
)

const basefeeWiggleMultiplier = 2

type transactorAccount struct {
	keySigner keysigner.KeySigner
	chainID   *big.Int
	ethClient *ethclient.Client
}

func newTransactorAccount(chainID *big.Int, keystorePath, password string, l1RPCClient *ethclient.Client) (*transactorAccount, error) {
	keySigner, err := keysigner.NewKeystoreSigner(keystorePath, password)
	if err != nil {
		return nil, fmt.Errorf("failed creating key signer: %w", err)
	}

	return &transactorAccount{
		keySigner: keySigner,
		chainID:   chainID,
		ethClient: l1RPCClient,
	}, nil
}

func (t *transactorAccount) Address() common.Address {
	return t.keySigner.GetAddress()
}

func (t *transactorAccount) SendTransaction(
	ctx context.Context,
	recipient common.Address,
	amount *big.Int,
) (*types.Transaction, error) {
	nonce, err := t.ethClient.PendingNonceAt(ctx, t.keySigner.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("failed getting account nonce: %w", err)
	}

	gasLimit := uint64(21000) // basic gas limit for Ether transfer

	head, err := t.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed getting the latest block: %w", err)
	}

	tip, err := t.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting gas tip cap: %w", err)
	}

	feeCap := new(big.Int).Add(
		tip,
		new(big.Int).Mul(head.BaseFee, big.NewInt(basefeeWiggleMultiplier)),
	)

	txData := &types.DynamicFeeTx{
		To:        &recipient,
		Nonce:     nonce,
		GasFeeCap: feeCap,
		GasTipCap: tip,
		Gas:       gasLimit,
		Value:     amount,
	}

	tx := types.NewTx(txData)

	signedTx, err := t.keySigner.SignTx(tx, t.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction payload: %w", err)
	}

	err = t.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}
