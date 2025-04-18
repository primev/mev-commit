package main

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"

	"github.com/primev/mev-commit/tools/bidder-bot/service"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
)

var (
	optionKeystorePath = &cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"KEYSTORE_DIR"},
		Required: true,
	}

	optionKeystorePassword = &cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"KEYSTORE_PASSWORD"},
		Required: true,
	}

	optionL1RPCUrls = &cli.StringSliceFlag{
		Name:     "l1-rpc-urls",
		Usage:    "URLs for L1 RPC",
		EnvVars:  []string{"L1_RPC_URLS"},
		Required: true,
	}

	optionL1WsUrls = &cli.StringSliceFlag{
		Name:     "l1-ws-urls",
		Usage:    "URLs for L1 WebSocket",
		EnvVars:  []string{"L1_WS_URLS"},
		Required: true,
	}

	optionBeaconApiUrls = &cli.StringSliceFlag{
		Name:     "beacon-api-urls",
		Usage:    "URLs for Beacon API endpoints",
		EnvVars:  []string{"BEACON_API_URLS"},
		Required: true,
	}

	optionSettlementRPCUrl = &cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"SETTLEMENT_RPC_URL"},
		Required: true,
	}

	optionBidderNodeRPCUrl = &cli.StringFlag{
		Name:     "bidder-node-rpc-url",
		Usage:    "URL for mev-commit bidder node RPC",
		EnvVars:  []string{"BIDDER_NODE_RPC_URL"},
		Required: true,
	}

	optionAutoDepositAmount = &cli.StringFlag{
		Name:    "auto-deposit-amount",
		Usage:   "auto deposit amount",
		EnvVars: []string{"AUTO_DEPOSIT_AMOUNT"},
		Value:   "100000000000000000", // 0.1 ETH
	}

	optionBidAmount = &cli.StringFlag{
		Name:    "bid-amount",
		Usage:   "amount to use for each bid",
		EnvVars: []string{"BID_AMOUNT"},
		Value:   "5000000000000000", // 0.005 ETH
	}

	optionUseFullNotifier = &cli.BoolFlag{
		Name:    "use-full-notifier",
		Usage:   "whether to use full notifier (false = selective opted-in notifier)",
		EnvVars: []string{"USE_FULL_NOTIFIER"},
		Value:   false,
	}

	optionCheckBalances = &cli.BoolFlag{
		Name:    "check-balances",
		Usage:   "whether to periodically check account balances",
		EnvVars: []string{"CHECK_BALANCES"},
		Value:   true,
	}

	optionGasTipCap = &cli.StringFlag{
		Name:    "gas-tip-cap",
		Usage:   "gas tip cap",
		EnvVars: []string{"GAS_TIP_CAP"},
		Value:   "15000000", // 0.015 gwei
	}

	optionGasFeeCap = &cli.StringFlag{
		Name:    "gas-fee-cap",
		Usage:   "gas fee cap",
		EnvVars: []string{"GAS_FEE_CAP"},
		Value:   "1000000000", // 1 gwei
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"LOG_FMT"},
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
		EnvVars: []string{"LOG_LEVEL"},
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
		EnvVars: []string{"LOG_TAGS"},
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
		Name:  "bidder-bot",
		Usage: "Bidder bot service",
		Flags: []cli.Flag{
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
			optionKeystorePath,
			optionKeystorePassword,
			optionL1RPCUrls,
			optionL1WsUrls,
			optionBeaconApiUrls,
			optionSettlementRPCUrl,
			optionBidderNodeRPCUrl,
			optionGasTipCap,
			optionGasFeeCap,
			optionAutoDepositAmount,
			optionBidAmount,
			optionUseFullNotifier,
			optionCheckBalances,
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

			bidAmount, ok := new(big.Int).SetString(c.String(optionBidAmount.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse bid-amount")
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
				Logger:            logger,
				GasTipCap:         gasTipCap,
				GasFeeCap:         gasFeeCap,
				AutoDepositAmount: autoDepositAmount,
				BidAmount:         bidAmount,
				SettlementRPCUrl:  c.String(optionSettlementRPCUrl.Name),
				BidderNodeRPC:     c.String(optionBidderNodeRPCUrl.Name),
				L1RPCUrls:         c.StringSlice(optionL1RPCUrls.Name),
				L1WsUrls:          c.StringSlice(optionL1WsUrls.Name),
				BeaconApiUrls:     c.StringSlice(optionBeaconApiUrls.Name),
				IsFullNotifier:    c.Bool(optionUseFullNotifier.Name),
				Signer:            signer,
				CheckBalances:     c.Bool(optionCheckBalances.Name),
			}

			logger.Debug("service config", "config", config)

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
