package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/shared"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/transfer"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/util"
	"github.com/urfave/cli/v2"
)

var errNoPendingTransactionFound = errors.New("no pending transaction found")

func main() {
	app := &cli.App{
		Name:  "bridge-cli",
		Usage: "CLI for interacting with a custom between L1 and the mev-commit (settlement) chain",
		Commands: []*cli.Command{
			{
				Name:  "bridge-to-settlement",
				Usage: "Submit a transaction to bridge ether to the settlement chain",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "amount",
						Usage:    "Amount of ether to bridge in wei",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "dest-addr",
						Usage:    "Destination address on the mev-commit (settlement) chain",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "cancel-pending",
						Usage: "Automatically cancel existing pending transactions",
					},
				},
				Action: bridgeToSettlement,
			},
			{
				Name:  "bridge-to-l1",
				Usage: "Submit a transaction to bridge ether back to L1",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "amount",
						Usage:    "Amount of ether to bridge in wei",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "dest-addr",
						Usage:    "Destination address on L1",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "cancel-pending",
						Usage: "Automatically cancel existing pending transactions",
					},
				},
				Action: bridgeToL1,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "Exited with error: %v\n", err)
	}
}

func loadConfig() (*envConfig, error) {
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	if err := checkEnvConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

func bridgeToSettlement(c *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	logger, err := util.NewLogger(cfg.LogLevel, "text", "", os.Stdout)
	if err != nil {
		return err
	}
	config, err := preTransfer(c, cfg)
	if err != nil {
		return err
	}
	autoCancel := c.Bool("cancel-pending")
	ok, err := handlePendingTxes(c.Context, logger.With("component", "l1_eth_client"), config.PrivateKey, config.L1RPCUrl, autoCancel)
	switch {
	case err == nil && !ok:
		logger.Info("user chose not to cancel pending transactions, exiting...")
		return nil
	case errors.Is(err, errNoPendingTransactionFound):
		// Do nothing.
	case err != nil:
		return err
	}

	t, err := transfer.NewTransferToSettlement(
		logger.With("component", "settlement_transfer"),
		config.Amount,
		config.DestAddress,
		config.PrivateKey,
		config.SettlementRPCUrl,
		config.L1RPCUrl,
		config.L1ContractAddr,
		config.SettlementContractAddr,
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer to settlement: %w", err)
	}
	err = t.Start(c.Context)
	if err != nil {
		return fmt.Errorf("failed to start transfer to settlement: %w", err)
	}
	return nil
}

func bridgeToL1(c *cli.Context) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	logger, err := util.NewLogger(cfg.LogLevel, "text", "", os.Stdout)
	if err != nil {
		return err
	}
	config, err := preTransfer(c, cfg)
	if err != nil {
		return err
	}
	autoCancel := c.Bool("cancel-pending")
	ok, err := handlePendingTxes(c.Context, logger.With("component", "settlement_eth_client"), config.PrivateKey, config.SettlementRPCUrl, autoCancel)
	switch {
	case err == nil && !ok:
		logger.Info("user chose not to cancel pending transactions, exiting...")
		return nil
	case errors.Is(err, errNoPendingTransactionFound):
		// Do nothing.
	case err != nil:
		return err
	}

	t, err := transfer.NewTransferToL1(
		logger.With("component", "l1_transfer"),
		config.Amount,
		config.DestAddress,
		config.PrivateKey,
		config.SettlementRPCUrl,
		config.L1RPCUrl,
		config.L1ContractAddr,
		config.SettlementContractAddr,
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer to L1: %w", err)
	}
	err = t.Start(c.Context)
	if err != nil {
		return fmt.Errorf("failed to start transfer to L1: %w", err)
	}
	return nil
}

type preTransferConfig struct {
	Amount                 *big.Int
	DestAddress            common.Address
	PrivateKey             *ecdsa.PrivateKey
	SettlementRPCUrl       string
	L1RPCUrl               string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
}

func preTransfer(c *cli.Context, cfg *envConfig) (*preTransferConfig, error) {
	privKeyTrimmed := strings.TrimPrefix(cfg.PrivKey, "0x")
	privKey, err := crypto.HexToECDSA(privKeyTrimmed)
	if err != nil {
		return nil, errors.New("failed to load private key")
	}

	amount := c.Int("amount")
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	destAddr := c.String("dest-addr")
	if !common.IsHexAddress(destAddr) {
		return nil, errors.New("dest-addr must be a valid hex address")
	}

	return &preTransferConfig{
		Amount:                 big.NewInt(int64(amount)),
		DestAddress:            common.HexToAddress(destAddr),
		PrivateKey:             privKey,
		SettlementRPCUrl:       cfg.SettlementRPCUrl,
		L1RPCUrl:               cfg.L1RPCUrl,
		L1ContractAddr:         common.HexToAddress(cfg.L1ContractAddr),
		SettlementContractAddr: common.HexToAddress(cfg.SettlementContractAddr),
	}, nil
}

type envConfig struct {
	PrivKey                string
	LogLevel               string
	L1RPCUrl               string
	SettlementRPCUrl       string
	L1ChainID              int
	SettlementChainID      int
	L1ContractAddr         string
	SettlementContractAddr string
}

func loadConfigFromEnv() (*envConfig, error) {
	l1ChainID := os.Getenv("L1_CHAIN_ID")
	l1ChainIDInt, err := strconv.Atoi(l1ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert L1_CHAIN_ID to int: %w", err)
	}
	settlementChainID := os.Getenv("SETTLEMENT_CHAIN_ID")
	settlementChainIDInt, err := strconv.Atoi(settlementChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert SETTLEMENT_CHAIN_ID to int: %w", err)
	}
	return &envConfig{
		PrivKey:                os.Getenv("PRIVATE_KEY"),
		LogLevel:               os.Getenv("LOG_LEVEL"),
		L1RPCUrl:               os.Getenv("L1_RPC_URL"),
		SettlementRPCUrl:       os.Getenv("SETTLEMENT_RPC_URL"),
		L1ChainID:              l1ChainIDInt,
		SettlementChainID:      settlementChainIDInt,
		L1ContractAddr:         os.Getenv("L1_CONTRACT_ADDR"),
		SettlementContractAddr: os.Getenv("SETTLEMENT_CONTRACT_ADDR"),
	}, nil
}

func checkEnvConfig(cfg *envConfig) error {
	if cfg.PrivKey == "" {
		return fmt.Errorf("private_key is required")
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.L1RPCUrl == "" || cfg.SettlementRPCUrl == "" {
		return fmt.Errorf("both l1_rpc_url and settlement_rpc_url are required")
	}
	if cfg.L1ChainID != 39999 && cfg.L1ChainID != 17000 {
		return fmt.Errorf("l1_chain_id must be 39999 (local l1) or 17000 (Holesky). Only test instances are supported")
	}
	if cfg.SettlementChainID != 17864 {
		return fmt.Errorf("settlement_chain_id must be 17864. Only test chains are supported")
	}
	if !common.IsHexAddress(cfg.L1ContractAddr) || !common.IsHexAddress(cfg.SettlementContractAddr) {
		return fmt.Errorf("both l1_contract_addr and settlement_contract_addr must be valid hex addresses")
	}
	l1Client, err := ethclient.Dial(cfg.L1RPCUrl)
	if err != nil {
		return fmt.Errorf("failed to create l1 client: %v", err)
	}
	obtainedL1ChainID, err := l1Client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get l1 chain id: %v", err)
	}
	if obtainedL1ChainID.Cmp(big.NewInt(int64(cfg.L1ChainID))) != 0 {
		return fmt.Errorf("l1 chain id mismatch. Expected: %d, Obtained: %d", cfg.L1ChainID, obtainedL1ChainID)
	}
	settlementClient, err := ethclient.Dial(cfg.SettlementRPCUrl)
	if err != nil {
		return fmt.Errorf("failed to create settlement client: %v", err)
	}
	obtainedSettlementChainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get settlement chain id: %v", err)
	}
	if obtainedSettlementChainID.Cmp(big.NewInt(int64(cfg.SettlementChainID))) != 0 {
		return fmt.Errorf("settlement chain id mismatch. Expected: %d, Obtained: %d", cfg.SettlementChainID, obtainedSettlementChainID)
	}
	return nil
}

func handlePendingTxes(
	ctx context.Context,
	logger *slog.Logger,
	privateKey *ecdsa.PrivateKey,
	url string,
	autoCancel bool,
) (bool, error) {
	rawClient, err := ethclient.Dial(url)
	if err != nil {
		return false, fmt.Errorf("failed to connect to eth client: %w", err)
	}
	ethClient := shared.NewETHClient(logger, rawClient)

	exist, err := ethClient.PendingTransactionsExist(ctx, privateKey)
	if err != nil {
		return false, fmt.Errorf("failed to check pending transactions: %w", err)
	}
	if !exist {
		return false, errNoPendingTransactionFound
	}
	if autoCancel {
		if err := ethClient.CancelPendingTxes(ctx, privateKey); err != nil {
			return false, fmt.Errorf("fail to cancel pending transaction(s): %w", err)
		}
		return true, nil
	}
	fmt.Println("Pending transactions exist for signing account. Do you want to cancel them? (y/n)")
	var response string
	_, err = fmt.Scanln(&response)
	if err != nil {
		return false, fmt.Errorf("failed to read user input: %w", err)
	}
	if strings.ToLower(response) == "y" {
		if err := ethClient.CancelPendingTxes(ctx, privateKey); err != nil {
			return false, fmt.Errorf("fail to cancel pending transaction(s): %w", err)
		}
		return true, nil
	}
	return false, nil
}
