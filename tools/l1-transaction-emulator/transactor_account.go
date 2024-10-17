package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/keysigner"
)

const basefeeWiggleMultiplier = 2

type transactorAccount struct {
	keySigner  keysigner.KeySigner
	chainID    *big.Int
	monitor    *txmonitor.Monitor
	transactor bind.ContractTransactor
	ethClient  *ethclient.Client
	closeFn    func(context.Context) error
}

func newTransactorAccount(logger *slog.Logger, keystorePath, password string, l1RPCClient *ethclient.Client) (*transactorAccount, error) {
	chainID, err := l1RPCClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed getting chain ID: %w", err)
	}

	keySigner, err := keysigner.NewKeystoreSigner(keystorePath, password)
	if err != nil {
		return nil, fmt.Errorf("failed creating key signer: %w", err)
	}

	logger = logger.With("account", keySigner.GetAddress().Hex())

	monitor := txmonitor.New(
		keySigner.GetAddress(),
		l1RPCClient,
		txmonitor.NewEVMHelperWithLogger(l1RPCClient, logger, make(map[common.Address]*abi.ABI)),
		&testSaver{logger},
		logger,
		64,
	)

	ctx, cancel := context.WithCancel(context.Background())
	monitorClosed := monitor.Start(ctx)

	txtor := transactor.NewTransactor(l1RPCClient, monitor)

	return &transactorAccount{
		keySigner:  keySigner,
		chainID:    chainID,
		monitor:    monitor,
		transactor: txtor,
		ethClient:  l1RPCClient,
		closeFn: func(ctx context.Context) error {
			cancel()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-monitorClosed:
				return nil
			}
		},
	}, nil
}

func (t *transactorAccount) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return t.closeFn(ctx)
}

func (t *transactorAccount) Address() common.Address {
	return t.keySigner.GetAddress()
}

func (t *transactorAccount) SendTransaction(
	ctx context.Context,
	recipient common.Address,
	amount *big.Int,
) error {

	nonce, err := t.transactor.PendingNonceAt(ctx, t.keySigner.GetAddress())
	if err != nil {
		return fmt.Errorf("failed getting account nonce: %w", err)
	}

	gasLimit := uint64(21000) // basic gas limit for Ether transfer

	head, err := t.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed getting the latest block: %w", err)
	}

	tip, err := t.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return fmt.Errorf("failed getting gas tip cap: %w", err)
	}

	// Add boost % to tip
	tip = new(big.Int).Div(new(big.Int).Mul(big.NewInt(110), tip), big.NewInt(100))
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
		return fmt.Errorf("failed to sign transaction payload: %w", err)
	}

	err = t.transactor.SendTransaction(ctx, signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	return nil
}

type testSaver struct {
	logger *slog.Logger
}

func (t *testSaver) Save(_ context.Context, txHash common.Hash, nonce uint64) error {
	t.logger.Info("sent tx", "txHash", txHash.Hex(), "nonce", nonce)
	return nil
}

func (t *testSaver) Update(_ context.Context, txHash common.Hash, status string) error {
	t.logger.Info("tx status updated", "txHash", txHash.Hex(), "status", status)
	return nil
}

func (t *testSaver) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	return nil, nil
}
