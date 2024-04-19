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

	ec := NewETHClient(nil, client)

	pks := []string{
		"0x894082d1fa8c0fa1ad258f0f03fb1dfe6f2bc0a28355e2f8f55608ebc93016474f877551451cd44b3726c1e6dd7f9fec",
		"0xa30e2aee808f5a141f036a853aaf6c6457c1d15a900e7bbc511608bf54b918b2284cc578c5074781f899c23028472a11",
		"0xa97c84a4a076e26680f33c39fc018428f7a2ab0e14b4cae776ed4d5050464fe67ea1423242c18e089b23f0d51b6c6f33",
		"0x8456186dc2c9d6f78425bed99eb35d1aa0a6ffb21eabe062e3651779ec671e52f447e506d5f492a976cd3a6e17816a31",
		"0x948e7ebb44a6b54655e5ddc893bc60880f1019737c679b90f4324bdadb2d3a0de146624d331ea41adcb8e429e98a6587",
		"0x99346362fbd4665f91df5aaf4a57a0376b9405cb35ff2a80598cb62e8ba115c1396753fa5e6b22a9d74747e43e5aee5a",
		"0xac33f519e217d18a1c86dbc06682400c806a99494329428acb8fb7b8cc8ccb56b64f02b977d973e70cb16415c448e17d",
		"0xa1aff153fbe7e3e2629ed557cc08f1c20c277708c9180704b8ea9526d5e561e3deb39d18491efb90fb3cd83eb91003b3",
		"0xb6d1fd1a5dbe7bb04328dc43c2dc3427200ce447d9391a881b8be94a50f6b42003b13a0efc1566c171eccdf1e14b2bdc",
		"0x93bb30895fa8afacab9a0d7d790db2356c729ebeb21045cd05cd7da5113f481f7be77b2908bd64b979e2500a2757ca96",
		"0xa933d0c54a5974f1c84e4ea47d4c537220f90be42fd67f244dd358daf90f3fc1b841c60ba1af04aa0b4a661a6b99822a",
		"0x99e833803f1757e8f7823c55da4dc35ed7f38bb9bed429dbba86d95bf81dd8cb82f698517e7b019b0ea4f6b1b06accc7",
		"0xb59a6109f25c4270aa3e61482bed3b6db0cbc49fb6ed72f7d55de242ab20a4b5c0b1cc5893d81576213a796463698725",
		"0xa81cdefcfb33f785eec7f5fe3465338c42043fafb7f4a6b2d603c133d7fb48d8186225caad9b59739c88508cc2d5345d",
		"0x8fcfe4e7b8334d9cef44b1cea2449b5fa881bee7fe977843dd6752a00b1e58315c3290e75084dc5e2243ea50d962db6d",
		"0xaa3b6da9757fc309b6473a95ba043b0254b352bacb7ff788f5353def4d8143ed5d4acd3b1a698121efea06d41a8a0296",
		"0xb54a37216fb48544a161b7a24e389ef3937beb207539f39960b622193a74288095974f42a995dc7f3f6e379e8b16a5e4",
		"0x934525116f3a700bb7a9834b0fd8bed942db7f7176a4921ab17ef665b9ae3fddaf5aad297dac9f6c7b1816eb8115c1f5",
		"0x8546b3c1421fdf8321273de417d86d63658250cb653b70f554a4a3dd6c3fbf6a768d62c61047b032daa4765fd714ef43",
		"0xa5642eb0ad86722a953710e0afc49cd9de5c7dd0ea305a8e147aebc03e1f62ee8e72ac65111853c3ac0cdf0a72042900",
		"0x97ffdf67b0fb9310eec15bfcdc7fc96b7c0713346d5f95cd3696a1ce690158b88c28dc830476da238511f3e3011aea96",
		"0xa9948f2bffe58516e2d04ff5dd2ec7f13e393fdaf43ef900c7bc6fad8fcb202a7958770bd3f5557690a4bc7a9f66f00c",
		"0x857e5abf12bd3d5b2da03387afb8d1cda21d64f98a069f5fd891e10693385174b45f58b65dafb3d9c096a947697dbb02",
		"0x8c010c1bc0f160340ef60c1e59bce549e87bba79f675590778d058525ccb29e044c515a4c7bd4684e0eb2115aafa0db6",
		"0x9268e508d6945661c49d8c7a7ad2d719e32336a8e8989ac59808954bfd9b958f68193a4a83ecc6125c00b9728798b508",
		"0xa4ea44cef4f6212d644d4859a3d8fa79112685b68c2e02fe6ef345f5e39f60425d34a5590ba3ff22bbcfaf2aecc9d05c",
	}

	pksAsBytes := make([][]byte, len(pks))
	for i, pk := range pks {
		if pk[:2] == "0x" {
			pk = pk[2:]
		}
		pksAsBytes[i] = common.Hex2Bytes(pk)
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
		fmt.Println("Batch iteration completed. Idx: ", idx)
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
