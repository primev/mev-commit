package txmonitor

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/keysigner"
)

type Canceller struct {
	chainID   *big.Int
	ethClient *ethclient.Client
	keySigner keysigner.KeySigner
	monitor   *Monitor
	logger    *slog.Logger
}

func NewCanceller(
	chainID *big.Int,
	ethClient *ethclient.Client,
	keySigner keysigner.KeySigner,
	monitor *Monitor,
	logger *slog.Logger,
) *Canceller {
	return &Canceller{
		chainID:   chainID,
		ethClient: ethClient,
		keySigner: keySigner,
		monitor:   monitor,
		logger:    logger,
	}
}

func (c *Canceller) suggestMaxFeeAndTipCap(
	ctx context.Context,
	gasPrice *big.Int,
) (*big.Int, *big.Int, error) {
	// Returns priority fee per gas
	gasTipCap, err := c.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to suggest gas tip cap: %w", err)
	}

	// Returns priority fee per gas + base fee per gas
	if gasPrice == nil {
		gasPrice, err = c.ethClient.SuggestGasPrice(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to suggest gas price: %w", err)
		}
	}

	return gasPrice, gasTipCap, nil
}

func (c *Canceller) CancelTx(ctx context.Context, txnHash common.Hash) (common.Hash, error) {
	txn, isPending, err := c.ethClient.TransactionByHash(ctx, txnHash)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	if !isPending {
		return common.Hash{}, ethereum.NotFound
	}

	gasFeeCap, gasTipCap, err := c.suggestMaxFeeAndTipCap(ctx, txn.GasPrice())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to suggest max fee and tip cap: %w", err)
	}

	if gasFeeCap.Cmp(txn.GasFeeCap()) <= 0 {
		gasFeeCap = txn.GasFeeCap()
	}

	if gasTipCap.Cmp(txn.GasTipCap()) <= 0 {
		gasTipCap = txn.GasTipCap()
	}

	// increase gas fee cap and tip cap by 10% for better chance of replacing
	gasTipCap = new(big.Int).Div(new(big.Int).Mul(gasTipCap, big.NewInt(110)), big.NewInt(100))
	gasFeeCap = new(big.Int).Div(new(big.Int).Mul(gasFeeCap, big.NewInt(110)), big.NewInt(100))

	owner := c.keySigner.GetAddress()

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     txn.Nonce(),
		ChainID:   c.chainID,
		To:        &owner,
		Value:     big.NewInt(0),
		Gas:       21000,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Data:      []byte{},
	})

	signedTx, err := c.keySigner.SignTx(tx, c.chainID)
	if err != nil {
		c.logger.Error("failed to sign cancel tx", "err", err)
		return common.Hash{}, fmt.Errorf("failed to sign cancel tx: %w", err)
	}

	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		c.logger.Error("failed to send cancel tx", "err", err)
		return common.Hash{}, err
	}

	c.logger.Info("sent cancel txn", "txHash", signedTx.Hash().Hex())
	c.monitor.Sent(ctx, signedTx)

	return signedTx.Hash(), nil
}
