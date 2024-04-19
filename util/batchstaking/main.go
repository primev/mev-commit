package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	vr "github.com/primevprotocol/mev-commit/contracts-abi/clients/ValidatorRegistry"
)

func main() {

	privateKeyString := os.Getenv("PRIVATE_KEY")
	if privateKeyString == "" {
		privateKeyString = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		fmt.Println("PRIVATE_KEY env var not supplied. Using default account")
	}

	if privateKeyString[:2] == "0x" {
		privateKeyString = privateKeyString[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		fmt.Println("Failed to parse private key")
		os.Exit(1)
	}

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain id: %v", err)
	}
	fmt.Println("Chain ID: ", chainID)

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get account balance: %v", err)
	}
	if balance.Cmp(big.NewInt(3100000000000000000)) == -1 {
		log.Fatalf("Insufficient balance. Please fund %v with at least 3.1 ETH", fromAddress.Hex())
	}

	// TODO: make this configurable
	contractAddress := common.HexToAddress("0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512")

	vrt, err := vr.NewValidatorregistryTransactor(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create Validator Registry transactor: %v", err)
	}

	vrc, err := vr.NewValidatorregistryCaller(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create Validator Registry caller: %v", err)
	}

	amount, err := vrc.GetStakedAmount(nil, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get staked amount: %v", err)
	}

	fmt.Println("Self staked amount: ", amount)

	ec := NewETHClient(nil, client)
	opts, err := ec.CreateTransactOpts(context.Background(), privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transact opts: %v", err)
	}
	opts.Value = big.NewInt(3100000000000000000) // 3.1 ETH in wei

	submitTx := func(
		ctx context.Context,
		opts *bind.TransactOpts,
	) (*types.Transaction, error) {
		tx, err := vrt.SelfStake(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to self stake: %w", err)
		}
		fmt.Println("Self stake sent. Transaction hash: ", tx.Hash().Hex())
		return tx, nil
	}

	receipt, err := ec.WaitMinedWithRetry(context.Background(), opts, submitTx)
	if err != nil {
		log.Fatalf("Failed to wait for self stake tx to be mined: %v", err)
	}

	fmt.Println("Self stake included in block: ", receipt.BlockNumber)

	amount, err = vrc.GetStakedAmount(nil, common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"))

	if err != nil {
		log.Fatalf("Failed to get staked amount: %v", err)
	}

	fmt.Println("Staked amount after selfStake: ", amount)

	// Split stake between 50 addrs, 200 times
	for i := 0; i < 200; i++ {

		addrs := make([]common.Address, 50)
		for i := 0; i < 50; i++ {
			key, err := crypto.GenerateKey()
			if err != nil {
				log.Fatalf("Failed to generate key: %v", err)
			}
			addrs[i] = crypto.PubkeyToAddress(key.PublicKey)
		}

		for _, addr := range addrs {
			amount, err = vrc.GetStakedAmount(nil, addr)
			if err != nil {
				log.Fatalf("Failed to get staked amount: %v", err)
			}
			fmt.Println("Initial staked amount for ", addr.Hex(), " is: ", amount)
		}

		opts, err = ec.CreateTransactOpts(context.Background(), privateKey, chainID)
		if err != nil {
			log.Fatalf("Failed to create transact opts: %v", err)
		}

		// 3.1 ETH * 50 = 155 ETH
		totalAmount := new(big.Int)
		totalAmount.SetString("155000000000000000000", 10)
		opts.Value = totalAmount

		submitTx = func(
			ctx context.Context,
			opts *bind.TransactOpts,
		) (*types.Transaction, error) {
			tx, err := vrt.SplitStake(opts, addrs)
			if err != nil {
				return nil, fmt.Errorf("failed to self stake: %w", err)
			}
			fmt.Println("Split stake sent. Transaction hash: ", tx.Hash().Hex())
			return tx, nil
		}

		receipt, err = ec.WaitMinedWithRetry(context.Background(), opts, submitTx)
		if err != nil {
			log.Fatalf("Failed to wait for split stake tx to be mined: %v", err)
		}
		fmt.Println("Split stake included in block: ", receipt.BlockNumber)

		for _, addr := range addrs {
			amount, err = vrc.GetStakedAmount(nil, addr)
			if err != nil {
				log.Fatalf("Failed to get staked amount: %v", err)
			}
			fmt.Println("Final staked amount for ", addr.Hex(), " is: ", amount)
		}
		fmt.Println("-------------------")
		fmt.Println("Batch iteration completed. Idx: ", i)
		fmt.Println("-------------------")
	}

	fmt.Println("All batches completed!")
}

//
// The following should exist as a shared lib
//

// TODO: modularize with other files!!

// TODO: make issue or modularize in this PR! Include tests..

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
