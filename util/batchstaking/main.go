package main

import (
	"bufio"
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
		fmt.Println("PRIVATE_KEY env var not supplied")
		os.Exit(1)
	}

	if privateKeyString[:2] == "0x" {
		privateKeyString = privateKeyString[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		fmt.Println("Failed to parse private key")
		os.Exit(1)
	}

	client, err := ethclient.Dial("https://chainrpc.testnet.mev-commit.xyz")
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

	contractAddress := common.HexToAddress("0xAd291fcfDBA7c9c5545af35Ca52edDe6cBdF92e5") // Accurate as of 4-19-2024

	vrt, err := vr.NewValidatorregistryTransactor(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create Validator Registry transactor: %v", err)
	}

	vrc, err := vr.NewValidatorregistryCaller(contractAddress, client)
	if err != nil {
		log.Fatalf("Failed to create Validator Registry caller: %v", err)
	}

	ec := NewETHClient(nil, client)

	publicKeyFilePath := "./keys_example.txt"
	pksAsBytes, err := readBLSPublicKeysFromFile(publicKeyFilePath)
	if err != nil {
		log.Fatalf("Failed to read public keys from file: %v", err)
	}

	// Split into batches of 50
	type Batch struct {
		pubKeys [][]byte
	}
	batches := make([]Batch, 0)
	for i := 0; i < len(pksAsBytes); i += 50 {
		end := i + 50
		if end > len(pksAsBytes) {
			end = len(pksAsBytes)
		}
		batches = append(batches, Batch{pubKeys: pksAsBytes[i:end]})
	}

	for idx, batch := range batches {

		for _, pk := range batch.pubKeys {
			amount, err := vrc.GetStakedAmount(nil, pk)
			if err != nil {
				log.Fatalf("Failed to get staked amount: %v", err)
			}
			fmt.Println("Initial staked amount for ", common.Bytes2Hex(pk), " is: ", amount)
		}

		opts, err := ec.CreateTransactOpts(context.Background(), privateKey, chainID)
		if err != nil {
			log.Fatalf("Failed to create transact opts: %v", err)
		}

		// 3.1 ETH * 50 = 155 ETH
		totalAmount := new(big.Int)
		totalAmount.SetString("155000000000000000000", 10)
		opts.Value = totalAmount

		submitTx := func(
			ctx context.Context,
			opts *bind.TransactOpts,
		) (*types.Transaction, error) {
			tx, err := vrt.Stake(opts, batch.pubKeys)
			if err != nil {
				return nil, fmt.Errorf("failed to self stake: %w", err)
			}
			fmt.Println("Split stake sent. Transaction hash: ", tx.Hash().Hex())
			return tx, nil
		}

		receipt, err := ec.WaitMinedWithRetry(context.Background(), opts, submitTx)
		if err != nil {
			log.Fatalf("Failed to wait for split stake tx to be mined: %v", err)
		}
		fmt.Println("Split stake included in block: ", receipt.BlockNumber)

		for _, pk := range batch.pubKeys {
			amount, err := vrc.GetStakedAmount(nil, pk)
			if err != nil {
				log.Fatalf("Failed to get staked amount: %v", err)
			}
			fmt.Println("Final staked amount for ", common.Bytes2Hex(pk), " is: ", amount)
		}
		fmt.Println("-------------------")
		fmt.Printf("Batch %d completed\n", idx+1)
		fmt.Println("-------------------")
	}

	fmt.Println("All batches completed!")
}

func readBLSPublicKeysFromFile(filePath string) ([][]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := scanner.Text()
		if len(key) > 2 && key[:2] == "0x" {
			key = key[2:]
		}
		keyBytes := common.Hex2Bytes(key)
		keys = append(keys, keyBytes)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return keys, nil
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
