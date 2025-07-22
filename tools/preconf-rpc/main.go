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
	"github.com/primev/mev-commit/tools/preconf-rpc/service"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port for the HTTP server",
		EnvVars: []string{"PRECONF_RPC_HTTP_PORT"},
		Value:   8080,
	}

	optionKeystorePath = &cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"PRECONF_RPC_KEYSTORE_DIR"},
		Required: true,
	}

	optionKeystorePassword = &cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"PRECONF_RPC_KEYSTORE_PASSWORD"},
		Required: true,
	}

	optionPgHost = &cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"PRECONF_RPC_PG_HOST"},
		Value:   "localhost",
	}

	optionPgPort = &cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"PRECONF_RPC_PG_PORT"},
		Value:   5432,
	}

	optionPgUser = &cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"PRECONF_RPC_PG_USER"},
		Value:   "postgres",
	}

	optionPgPassword = &cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"PRECONF_RPC_PG_PASSWORD"},
		Value:   "postgres",
	}

	optionPgDbname = &cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"PRECONF_RPC_PG_DBNAME"},
		Value:   "mev_oracle",
	}

	optionPgSSL = &cli.BoolFlag{
		Name:    "pg-ssl",
		Usage:   "use SSL for PostgreSQL connection",
		EnvVars: []string{"PRECONF_RPC_PG_SSL"},
		Value:   false,
	}

	optionL1RPCUrls = &cli.StringSliceFlag{
		Name:     "l1-rpc-urls",
		Usage:    "URLs for L1 RPC",
		EnvVars:  []string{"PRECONF_RPC_L1_RPC_URLS"},
		Required: true,
	}

	optionSettlementRPCUrl = &cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"PRECONF_RPC_SETTLEMENT_RPC_URL"},
		Required: true,
	}

	optionBidderRPCUrl = &cli.StringFlag{
		Name:     "bidder-rpc-url",
		Usage:    "URL for mev-commit bidder RPC",
		EnvVars:  []string{"PRECONF_RPC_BIDDER_RPC_URL"},
		Required: true,
	}

	optionL1ContractAddr = &cli.StringFlag{
		Name:     "l1-contract-addr",
		Usage:    "address of the L1 gateway contract",
		EnvVars:  []string{"PRECONF_RPC_L1_CONTRACT_ADDR"},
		Required: true,
	}

	optionSettlementThreshold = &cli.StringFlag{
		Name:    "settlement-threshold",
		Usage:   "Minimum threshold for settlement chain balance",
		EnvVars: []string{"PRECONF_RPC_SETTLEMENT_THRESHOLD"},
		Value:   "5000000000000000000", // 5 ETH
	}

	optionSettlementTopup = &cli.StringFlag{
		Name:    "settlement-topup",
		Usage:   "topup for settlement",
		EnvVars: []string{"PRECONF_RPC_SETTLEMENT_TOPUP"},
		Value:   "10000000000000000000", // 10 ETH
	}

	optionAutoDepositAmount = &cli.StringFlag{
		Name:    "auto-deposit-amount",
		Usage:   "auto deposit amount",
		EnvVars: []string{"PRECONF_RPC_AUTO_DEPOSIT_AMOUNT"},
		Value:   "1000000000000000000", // 1 ETH
	}

	optionGasTipCap = &cli.StringFlag{
		Name:    "gas-tip-cap",
		Usage:   "gas tip cap",
		EnvVars: []string{"PRECONF_RPC_GAS_TIP_CAP"},
		Value:   "50000000", // 0.05 gWEI
	}

	optionGasFeeCap = &cli.StringFlag{
		Name:    "gas-fee-cap",
		Usage:   "gas fee cap",
		EnvVars: []string{"PRECONF_RPC_GAS_FEE_CAP"},
		Value:   "60000000", // 0.06 gWEI
	}

	optionSettlementContractAddr = &cli.StringFlag{
		Name:     "settlement-contract-addr",
		Usage:    "address of the settlement gateway contract",
		EnvVars:  []string{"PRECONF_RPC_SETTLEMENT_CONTRACT_ADDR"},
		Required: true,
	}

	optionDepositAddress = &cli.StringFlag{
		Name:     "deposit-address",
		Usage:    "address to deposit funds to",
		EnvVars:  []string{"PRECONF_RPC_DEPOSIT_ADDRESS"},
		Required: true,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid deposit address: %s", s)
			}
			return nil
		},
	}

	optionBridgeAddress = &cli.StringFlag{
		Name:     "bridge-address",
		Usage:    "address to bridge to",
		EnvVars:  []string{"PRECONF_RPC_BRIDGE_ADDRESS"},
		Required: true,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid bridge address: %s", s)
			}
			return nil
		},
	}

	optionBlocknativeAPIKey = &cli.StringFlag{
		Name:    "blocknative-api-key",
		Usage:   "Blocknative API key for transaction pricing",
		EnvVars: []string{"PRECONF_RPC_BLOCKNATIVE_API_KEY"},
		Value:   "",
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"PRECONF_RPC_LOG_FMT"},
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
		EnvVars: []string{"PRECONF_RPC_LOG_LEVEL"},
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
		EnvVars: []string{"PRECONF_RPC_LOG_TAGS"},
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
		Name:  "preconf-rpc",
		Usage: "Preconf RPC service",
		Flags: []cli.Flag{
			optionHTTPPort,
			optionPgHost,
			optionPgPort,
			optionPgUser,
			optionPgPassword,
			optionPgDbname,
			optionPgSSL,
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
			optionGasTipCap,
			optionGasFeeCap,
			optionSettlementContractAddr,
			optionAutoDepositAmount,
			optionDepositAddress,
			optionBridgeAddress,
			optionBlocknativeAPIKey,
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
				PgHost:                 c.String(optionPgHost.Name),
				PgPort:                 c.Int(optionPgPort.Name),
				PgUser:                 c.String(optionPgUser.Name),
				PgPassword:             c.String(optionPgPassword.Name),
				PgDbname:               c.String(optionPgDbname.Name),
				PgSSL:                  c.Bool(optionPgSSL.Name),
				Logger:                 logger,
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
				DepositAddress:         common.HexToAddress(c.String(optionDepositAddress.Name)),
				BridgeAddress:          common.HexToAddress(c.String(optionBridgeAddress.Name)),
				PricerAPIKey:           c.String(optionBlocknativeAPIKey.Name),
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
