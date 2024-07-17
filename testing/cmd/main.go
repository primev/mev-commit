package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/testing/pkg/tests"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionSettlementRPCEndpoint = &cli.StringFlag{
		Name:    "settlement-rpc-endpoint",
		Usage:   "Settlement RPC endpoint",
		Value:   "http://localhost:8545",
		EnvVars: []string{"MEV_COMMIT_TEST_SETTLEMENT_RPC_ENDPOINT"},
	}

	optionProviderRegistryAddress = &cli.StringFlag{
		Name:    "provider-registry-address",
		Usage:   "Provider registry address",
		EnvVars: []string{"MEV_COMMIT_TEST_PROVIDER_REGISTRY_ADDRESS"},
		Action: func(c *cli.Context, address string) error {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("invalid provider registry address")
			}
			return nil
		},
	}

	optionBootnodeRPCAddresses = &cli.StringSliceFlag{
		Name:    "bootnode-rpc-addresses",
		Usage:   "Bootnode RPC addresses",
		EnvVars: []string{"MEV_COMMIT_TEST_BOOTNODE_RPC_ADDRESSES"},
		Action: func(c *cli.Context, addresses []string) error {
			for _, address := range addresses {
				if _, _, err := net.SplitHostPort(address); err != nil {
					return fmt.Errorf("invalid bootnode RPC address")
				}
			}
			return nil
		},
	}

	optionProviderRPCAddresses = &cli.StringSliceFlag{
		Name:    "provider-rpc-addresses",
		Usage:   "Provider RPC addresses",
		EnvVars: []string{"MEV_COMMIT_TEST_PROVIDER_RPC_ADDRESSES"},
		Action: func(c *cli.Context, addresses []string) error {
			for _, address := range addresses {
				if _, _, err := net.SplitHostPort(address); err != nil {
					return fmt.Errorf("invalid provider RPC address")
				}
			}
			return nil
		},
	}

	optionBidderRPCAddresses = &cli.StringSliceFlag{
		Name:    "bidder-rpc-addresses",
		Usage:   "Bidder RPC addresses",
		EnvVars: []string{"MEV_COMMIT_TEST_BIDDER_RPC_ADDRESSES"},
		Action: func(c *cli.Context, addresses []string) error {
			for _, address := range addresses {
				if _, _, err := net.SplitHostPort(address); err != nil {
					return fmt.Errorf("invalid bidder RPC address")
				}
			}
			return nil
		},
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MEV_COMMIT_TEST_LOG_FMT"},
		Value:   "text",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log format")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"MEV_COMMIT_TEST_LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log level")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"MEV_COMMIT_TEST_LOG_TAGS"},
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
		Name:  "mev-commit-test",
		Usage: "MEV commit test",
		Flags: []cli.Flag{
			optionSettlementRPCEndpoint,
			optionProviderRegistryAddress,
			optionBootnodeRPCAddresses,
			optionProviderRPCAddresses,
			optionBidderRPCAddresses,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run MEV commit test",
				Action: func(c *cli.Context) error {
					return run(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func run(c *cli.Context) error {
	settlementRPCEndpoint := c.String(optionSettlementRPCEndpoint.Name)
	providerRegistryAddress := c.String(optionProviderRegistryAddress.Name)
	bootnodeRPCAddresses := c.StringSlice(optionBootnodeRPCAddresses.Name)
	providerRPCAddresses := c.StringSlice(optionProviderRPCAddresses.Name)
	bidderRPCAddresses := c.StringSlice(optionBidderRPCAddresses.Name)

	fmt.Println("Settlement RPC endpoint:", settlementRPCEndpoint)
	fmt.Println("Provider registry address:", providerRegistryAddress)
	fmt.Println("Bootnode RPC addresses:", bootnodeRPCAddresses)
	fmt.Println("Provider RPC addresses:", providerRPCAddresses)
	fmt.Println("Bidder RPC addresses:", bidderRPCAddresses)

	logger, err := util.NewLogger(
		c.String(optionLogLevel.Name),
		c.String(optionLogFmt.Name),
		c.String(optionLogTags.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	o, err := orchestrator.NewOrchestrator(orchestrator.Options{
		SettlementRPCEndpoint:   settlementRPCEndpoint,
		ProviderRegistryAddress: common.HexToAddress(providerRegistryAddress),
		BootnodeRPCAddresses:    bootnodeRPCAddresses,
		ProviderRPCAddresses:    providerRPCAddresses,
		BidderRPCAddresses:      bidderRPCAddresses,
		Logger:                  logger,
	})

	if err != nil {
		return err
	}

	defer o.Close()

	// Run test cases
	for name, tc := range tests.TestCases {
		logger.Info("running test case", "name", name)
		if err := tc(context.Background(), o, nil); err != nil {
			logger.Error("test case failed", "name", name, "error", err)
			return fmt.Errorf("test case %s failed: %w", name, err)
		}
		logger.Info("test case passed", "name", name)
	}

	logger.Info("all test cases passed")

	return nil
}
