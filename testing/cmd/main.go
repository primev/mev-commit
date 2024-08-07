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
		Name:     "settlement-rpc-endpoint",
		Usage:    "Settlement RPC endpoint",
		Required: true,
		EnvVars:  []string{"MEV_COMMIT_TEST_SETTLEMENT_RPC_ENDPOINT"},
	}

	optionL1RPCEndpoint = &cli.StringFlag{
		Name:     "l1-rpc-endpoint",
		Usage:    "L1 RPC endpoint",
		Required: true,
		EnvVars:  []string{"MEV_COMMIT_TEST_L1_RPC_ENDPOINT"},
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

	optionBidderRegistryAddress = &cli.StringFlag{
		Name:    "bidder-registry-address",
		Usage:   "Bidder registry address",
		EnvVars: []string{"MEV_COMMIT_TEST_BIDDER_REGISTRY_ADDRESS"},
		Action: func(c *cli.Context, address string) error {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("invalid bidder registry address")
			}
			return nil
		},
	}

	optionPreconfContractAddress = &cli.StringFlag{
		Name:    "preconf-contract-address",
		Usage:   "Preconfirmation contract address",
		EnvVars: []string{"MEV_COMMIT_TEST_PRECONF_CONTRACT_ADDRESS"},
		Action: func(c *cli.Context, address string) error {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("invalid preconf contract address")
			}
			return nil
		},
	}

	optionBlocktrackerContractAddress = &cli.StringFlag{
		Name:    "blocktracker-contract-address",
		Usage:   "Blocktracker contract address",
		EnvVars: []string{"MEV_COMMIT_TEST_BLOCKTRACKER_CONTRACT_ADDRESS"},
		Action: func(c *cli.Context, address string) error {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("invalid provider registry address")
			}
			return nil
		},
	}

	optionOracleContractAddress = &cli.StringFlag{
		Name:    "oracle-contract-address",
		Usage:   "Oracle contract address",
		EnvVars: []string{"MEV_COMMIT_TEST_ORACLE_CONTRACT_ADDRESS"},
		Action: func(c *cli.Context, address string) error {
			if !common.IsHexAddress(address) {
				return fmt.Errorf("invalid oracle address")
			}
			return nil
		},
	}

	optionBootnodeRPCAddresses = &cli.StringSliceFlag{
		Name:    "bootnode-rpc-addresses",
		Usage:   "Bootnode RPC addresses",
		EnvVars: []string{"MEV_COMMIT_TEST_BOOTNODE_RPC_ADDRESSES"},
		Action: func(c *cli.Context, addresses []string) error {
			if len(addresses) == 0 {
				return fmt.Errorf("at least one bootnode RPC address is required")
			}
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
			if len(addresses) == 0 {
				return fmt.Errorf("at least one provider RPC address is required")
			}
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
			if len(addresses) == 0 {
				return fmt.Errorf("at least one bidder RPC address is required")
			}
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
			optionL1RPCEndpoint,
			optionProviderRegistryAddress,
			optionBidderRegistryAddress,
			optionPreconfContractAddress,
			optionBlocktrackerContractAddress,
			optionOracleContractAddress,
			optionBootnodeRPCAddresses,
			optionProviderRPCAddresses,
			optionBidderRPCAddresses,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			return run(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func run(c *cli.Context) error {
	logger, err := util.NewLogger(
		c.String(optionLogLevel.Name),
		c.String(optionLogFmt.Name),
		c.String(optionLogTags.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	opts := orchestrator.Options{
		SettlementRPCEndpoint:       c.String(optionSettlementRPCEndpoint.Name),
		L1RPCEndpoint:               c.String(optionL1RPCEndpoint.Name),
		ProviderRegistryAddress:     common.HexToAddress(c.String(optionProviderRegistryAddress.Name)),
		BidderRegistryAddress:       common.HexToAddress(c.String(optionBidderRegistryAddress.Name)),
		PreconfContractAddress:      common.HexToAddress(c.String(optionPreconfContractAddress.Name)),
		BlockTrackerContractAddress: common.HexToAddress(c.String(optionBlocktrackerContractAddress.Name)),
		OracleContractAddress:       common.HexToAddress(c.String(optionOracleContractAddress.Name)),
		BootnodeRPCAddresses:        c.StringSlice(optionBootnodeRPCAddresses.Name),
		ProviderRPCAddresses:        c.StringSlice(optionProviderRPCAddresses.Name),
		BidderRPCAddresses:          c.StringSlice(optionBidderRPCAddresses.Name),
		Logger:                      logger,
	}

	logger.Info("running with options", "options", opts)

	o, err := orchestrator.NewOrchestrator(opts)

	if err != nil {
		return err
	}

	defer o.Close()

	// Run test cases
	for _, tc := range tests.TestCases {
		logger.Info("running test case", "name", tc.Name)
		if err := tc.Run(context.Background(), o, nil); err != nil {
			logger.Error("test case failed", "name", tc.Name, "error", err)
			return fmt.Errorf("test case %s failed: %w", tc.Name, err)
		}
		logger.Info("test case passed", "name", tc.Name)
	}

	logger.Info("all test cases passed")

	return nil
}
