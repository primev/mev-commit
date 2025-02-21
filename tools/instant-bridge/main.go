package main

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/tools/instant-bridge/service"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port for the HTTP server",
		EnvVars: []string{"INSTANT_BRIDGE_HTTP_PORT"},
		Value:   8080,
	}

	optionKeystorePath = &cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"INSTANT_BRIDGE_KEYSTORE_DIR"},
		Required: true,
	}

	optionKeystorePassword = &cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"INSTANT_BRIDGE_KEYSTORE_PASSWORD"},
		Required: true,
	}

	optionL1RPCUrls = &cli.StringSliceFlag{
		Name:     "l1-rpc-urls",
		Usage:    "URLs for L1 RPC",
		EnvVars:  []string{"INSTANT_BRIDGE_L1_RPC_URLS"},
		Required: true,
	}

	optionSettlementRPCUrl = &cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"INSTANT_BRIDGE_SETTLEMENT_RPC_URL"},
		Required: true,
	}

	optionBidderRPCUrl = &cli.StringFlag{
		Name:     "bidder-rpc-url",
		Usage:    "URL for mev-commit bidder RPC",
		EnvVars:  []string{"INSTANT_BRIDGE_SETTLEMENT_RPC_URL"},
		Required: true,
	}

	optionL1ContractAddr = &cli.StringFlag{
		Name:     "l1-contract-addr",
		Usage:    "address of the L1 gateway contract",
		EnvVars:  []string{"INSTANT_BRIDGE_L1_CONTRACT_ADDR"},
		Required: true,
	}

	optionSettlementThreshold = &cli.StringFlag{
		Name:    "settlement-threshold",
		Usage:   "Minimum threshold for settlement chain balance",
		EnvVars: []string{"INSTANT_BRIDGE_SETTLEMENT_THRESHOLD"},
		Value:   "5000000000000000000", // 5 ETH
	}

	optionSettlementTopup = &cli.StringFlag{
		Name:    "settlement-topup",
		Usage:   "topup for settlement",
		EnvVars: []string{"INSTANT_BRIDGE_SETTLEMENT_TOPUP"},
		Value:   "10000000000000000000", // 10 ETH
	}

	optionAutoDepositAmount = &cli.StringFlag{
		Name:    "auto-deposit-amount",
		Usage:   "auto deposit amount",
		EnvVars: []string{"INSTANT_BRIDGE_AUTO_DEPOSIT_AMOUNT"},
		Value:   "1000000000000000000", // 1 ETH
	}

	optionMinServiceFee = &cli.StringFlag{
		Name:    "min-service-fee",
		Usage:   "minimum service fee",
		EnvVars: []string{"INSTANT_BRIDGE_MIN_SERVICE_FEE"},
		Value:   "50000000000000000", // 0.05 ETH
	}

	optionGasTipCap = &cli.StringFlag{
		Name:    "gas-tip-cap",
		Usage:   "gas tip cap",
		EnvVars: []string{"INSTANT_BRIDGE_GAS_TIP_CAP"},
		Value:   "50000000", // 0.05 gWEI
	}

	optionGasFeeCap = &cli.StringFlag{
		Name:    "gas-fee-cap",
		Usage:   "gas fee cap",
		EnvVars: []string{"INSTANT_BRIDGE_GAS_FEE_CAP"},
		Value:   "60000000", // 0.06 gWEI
	}

	optionSettlementContractAddr = &cli.StringFlag{
		Name:     "settlement-contract-addr",
		Usage:    "address of the settlement gateway contract",
		EnvVars:  []string{"INSTANT_BRIDGE_SETTLEMENT_CONTRACT_ADDR"},
		Required: true,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"INSTANT_BRIDGE_LOG_FMT"},
		Value:   "text",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"INSTANT_BRIDGE_LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"INSTANT_BRIDGE_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	}
)

func main() {
	app := &cli.App{
		Name:  "instant-bridge",
		Usage: "Instant Bridge service",
		Flags: []cli.Flag{
			optionHTTPPort,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
			optionKeystorePath,
			optionKeystorePassword,
			optionL1RPCUrls,
			optionSettlementRPCUrl,
			optionBidderRPCUrl,
			optionL1ContractAddr,
			optionSettlementThreshold,
			optionSettlementTopup,
			optionMinServiceFee,
			optionGasTipCap,
			optionGasFeeCap,
			optionSettlementContractAddr,
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			minServiceFee, ok := new(big.Int).SetString(c.String(optionMinServiceFee.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse min-service-fee")
			}

			gasTipCap, ok := new(big.Int).SetString(c.String(optionGasTipCap.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse gas-tip-cap")
			}

			gasFeeCap, ok := new(big.Int).SetString(c.String(optionGasFeeCap.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse gas-fee-cap")
			}

			autoDepositAmount, ok := new(big.Int).SetString(c.String(optionAutoDepositAmount.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse auto-deposit-amount")
			}

			settlementThreshold, ok := new(big.Int).SetString(c.String(optionSettlementThreshold.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse settlement-threshold")
			}

			settlementTopup, ok := new(big.Int).SetString(c.String(optionSettlementTopup.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse settlement-topup")
			}

			signer, err := keysigner.NewKeystoreSigner(
				c.String(optionKeystorePath.Name),
				c.String(optionKeystorePassword.Name),
			)
			if err != nil {
				return fmt.Errorf("failed to create signer: %w", err)
			}

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

			config := service.Config{
				HTTPPort:               c.Int(optionHTTPPort.Name),
				Logger:                 logger,
				MinServiceFee:          minServiceFee,
				GasTipCap:              gasTipCap,
				GasFeeCap:              gasFeeCap,
				AutoDepositAmount:      autoDepositAmount,
				SettlementThreshold:    settlementThreshold,
				SettlementTopup:        settlementTopup,
				SettlementRPCUrl:       c.String(optionSettlementRPCUrl.Name),
				BidderRPC:              c.String(optionBidderRPCUrl.Name),
				L1RPCUrls:              c.StringSlice(optionL1RPCUrls.Name),
				L1ContractAddr:         common.HexToAddress(c.String(optionL1ContractAddr.Name)),
				SettlementContractAddr: common.HexToAddress(c.String(optionSettlementContractAddr.Name)),
				Signer:                 signer,
			}

			s, err := service.New(&config)
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}

			<-sigc
			logger.Info("shutting down...")

			return s.Close()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
