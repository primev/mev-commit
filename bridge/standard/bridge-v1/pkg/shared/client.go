package shared

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETHClient struct {
	logger *slog.Logger
	client *ethclient.Client
}

func NewETHClient(logger *slog.Logger, client *ethclient.Client) *ETHClient {
	return &ETHClient{logger: logger, client: client}
}

func (c *ETHClient) ChainID(ctx context.Context) (*big.Int, error) {
	return c.client.ChainID(ctx)
}

func (c *ETHClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.client.BlockNumber(ctx)
}

func (c *ETHClient) CreateTransactOpts(
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	srcChainID *big.Int,
) (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, srcChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	fromAddress := auth.From
	nonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending nonce: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasTip, gasPrice, err := c.SuggestGasTipCapAndPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas tip cap and price: %w", err)
	}

	auth.GasFeeCap = gasPrice
	auth.GasTipCap = gasTip
	auth.GasLimit = uint64(3000000)
	return auth, nil
}

func (c *ETHClient) SuggestGasTipCapAndPrice(ctx context.Context) (*big.Int, *big.Int, error) {
	// Returns priority fee per gas
	gasTip, err := c.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get gas tip cap: %w", err)
	}
	// Returns priority fee per gas + base fee per gas
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	return gasTip, gasPrice, nil
}

// TODO: Unit tests
func (c *ETHClient) BoostTipForTransactOpts(
	ctx context.Context,
	opts *bind.TransactOpts,
) error {
	c.logger.Debug(
		"gas params for tx that were not included",
		"gas_tip", opts.GasTipCap.String(),
		"gas_fee_cap", opts.GasFeeCap.String(),
		"base_fee", new(big.Int).Sub(opts.GasFeeCap, opts.GasTipCap).String(),
	)

	newGasTip, newFeeCap, err := c.SuggestGasTipCapAndPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas tip cap and price: %w", err)
	}

	newBaseFee := new(big.Int).Sub(newFeeCap, newGasTip)
	if newBaseFee.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("new base fee cannot be negative: %s", newBaseFee.String())
	}

	prevBaseFee := new(big.Int).Sub(opts.GasFeeCap, opts.GasTipCap)
	if prevBaseFee.Cmp(big.NewInt(0)) == -1 {
		return fmt.Errorf("base fee cannot be negative: %s", prevBaseFee.String())
	}

	var maxBaseFee *big.Int
	if newBaseFee.Cmp(prevBaseFee) == 1 {
		maxBaseFee = newBaseFee
	} else {
		maxBaseFee = prevBaseFee
	}

	var maxGasTip *big.Int
	if newGasTip.Cmp(opts.GasTipCap) == 1 {
		maxGasTip = newGasTip
	} else {
		maxGasTip = opts.GasTipCap
	}

	boostedTip := new(big.Int).Add(maxGasTip, new(big.Int).Div(maxGasTip, big.NewInt(10)))
	boostedTip = boostedTip.Add(boostedTip, big.NewInt(1))

	boostedBaseFee := new(big.Int).Add(maxBaseFee, new(big.Int).Div(maxBaseFee, big.NewInt(10)))
	boostedBaseFee = boostedBaseFee.Add(boostedBaseFee, big.NewInt(1))

	opts.GasTipCap = boostedTip
	opts.GasFeeCap = new(big.Int).Add(boostedBaseFee, boostedTip)

	c.logger.Debug("tip and base fee will be boosted by 10%")
	c.logger.Debug(
		"boosted gas",
		"get_tip_cap", opts.GasTipCap.String(),
		"gas_fee_cap", opts.GasFeeCap.String(),
		"base_fee", boostedBaseFee.String(),
	)

	return nil
}

type TxSubmitFunc func(
	ctx context.Context,
	opts *bind.TransactOpts,
) (
	tx *types.Transaction,
	err error,
)

// TODO: Unit tests
func (c *ETHClient) WaitMinedWithRetry(
	ctx context.Context,
	opts *bind.TransactOpts,
	submitTx TxSubmitFunc,
) (*types.Receipt, error) {

	const maxRetries = 10
	var err error
	var tx *types.Transaction

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Info("transaction not included within 60 seconds, boosting gas tip by 10%", "attempt", attempt)
			if err := c.BoostTipForTransactOpts(ctx, opts); err != nil {
				return nil, fmt.Errorf("failed to boost gas tip for attempt %d: %w", attempt, err)
			}
		}

		tx, err = submitTx(ctx, opts)
		if err != nil {
			if strings.Contains(err.Error(), "replacement transaction underpriced") || strings.Contains(err.Error(), "already known") {
				c.logger.Error("tx submission failed", "attempt", attempt, "error", err)
				continue
			}
			return nil, fmt.Errorf("tx submission failed on attempt %d: %w", attempt, err)
		}

		timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		receiptChan := make(chan *types.Receipt)
		errChan := make(chan error)

		go func() {
			receipt, err := bind.WaitMined(timeoutCtx, c.client, tx)
			if err != nil {
				errChan <- err
				return
			}
			receiptChan <- receipt
		}()

		select {
		case receipt := <-receiptChan:
			cancel()
			return receipt, nil
		case err := <-errChan:
			cancel()
			return nil, err
		case <-timeoutCtx.Done():
			cancel()
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("tx not included after %d attempts", maxRetries)
			}
			// Continue with boosted tip
		}
	}
	return nil, fmt.Errorf("unexpected error: control flow should not reach end of WaitMinedWithRetry")
}

func (c *ETHClient) CancelPendingTxes(ctx context.Context, privateKey *ecdsa.PrivateKey) error {
	if err := c.cancelAllPendingTransactions(ctx, privateKey); err != nil {
		return err
	}

	idx := 0
	timeoutSec := 60
	for {
		if idx >= timeoutSec {
			return fmt.Errorf("timeout: failed to cancel all pending transactions")
		}
		exist, err := c.PendingTransactionsExist(ctx, privateKey)
		if err != nil {
			return fmt.Errorf("failed to check pending transactions: %w", err)
		}
		if !exist {
			c.logger.Info("all pending transactions for signing account have been cancelled")
			return nil
		}
		time.Sleep(1 * time.Second)
		idx++
	}
}

// TODO: Use WaitMinedWithRetry
func (c *ETHClient) cancelAllPendingTransactions(
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
) error {
	chainID, err := c.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain id: %w", err)
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	currentNonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get current pending nonce: %w", err)
	}
	c.logger.Debug("current pending nonce", "nonce", currentNonce)

	latestNonce, err := c.client.NonceAt(ctx, fromAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get latest nonce: %w", err)
	}
	c.logger.Debug("latest nonce", "nonce", latestNonce)

	if currentNonce <= latestNonce {
		c.logger.Info("no pending transactions to cancel")
		return nil
	}

	suggestedGasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get suggested gas price: %w", err)
	}
	c.logger.Debug("suggested gas price", "gas_price", suggestedGasPrice.String())

	for nonce := latestNonce; nonce < currentNonce; nonce++ {
		gasPrice := new(big.Int).Set(suggestedGasPrice)
		const maxRetries = 5
		for retry := 0; retry < maxRetries; retry++ {
			if retry > 0 {
				increase := new(big.Int).Div(gasPrice, big.NewInt(10))
				gasPrice = gasPrice.Add(gasPrice, increase)
				gasPrice = gasPrice.Add(gasPrice, big.NewInt(1))
				c.logger.Debug("increased gas price for retry", "retry", retry, "gas_price", gasPrice.String())
			}

			tx := types.NewTransaction(nonce, fromAddress, big.NewInt(0), 21000, gasPrice, nil)
			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
			if err != nil {
				return fmt.Errorf("failed to sign cancellation transaction for nonce %d: %w", nonce, err)
			}

			err = c.client.SendTransaction(ctx, signedTx)
			if err != nil {
				if err.Error() == "replacement transaction underpriced" {
					c.logger.Warn("underpriced transaction, increasing gas price", "retry", retry+1, "nonce", nonce, "error", err)
					continue // Try again with a higher gas price
				}
				if err.Error() == "already known" {
					c.logger.Warn("already known transaction", "retry", retry+1, "nonce", nonce, "error", err)
					continue // Try again with a higher gas price
				}
				return fmt.Errorf("failed to send cancellation transaction for nonce %d: %w", nonce, err)
			}
			c.logger.Info("sent cancel transaction", "nonce", nonce, "tx_hash", signedTx.Hash().Hex(), "gas_price", gasPrice.String())
			break
		}
	}
	return nil
}

func (c *ETHClient) PendingTransactionsExist(ctx context.Context, privateKey *ecdsa.PrivateKey) (bool, error) {
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	currentNonce, err := c.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return false, fmt.Errorf("failed to get current pending nonce: %w", err)
	}

	latestNonce, err := c.client.NonceAt(ctx, fromAddress, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get latest nonce: %w", err)
	}

	return currentNonce > latestNonce, nil
}
