package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/relayer"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const (
	defaultHTTPPort  = 8080
	defaultConfigDir = "~/.mev-commit-bridge"
	defaultKeyFile   = "key"
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "path to relayer config file",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_CONFIG"},
	}

	optionPrivKeyFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "priv-key-file",
		Usage:   "path to private key file",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PRIV_KEY_FILE"},
		Value:   filepath.Join(defaultConfigDir, defaultKeyFile),
	})

	optionLogFmt = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_LOG_FMT"},
		Value:   "text",
		Action: func(_ *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid value: -log-fmt=%q", s)
			}
			return nil
		},
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_LOG_LEVEL"},
		Value:   "info",
		Action: func(_ *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid value: -log-level=%q", s)
			}
			return nil
		},
	})

	optionLogTags = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	})

	optionL1RPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL for L1 RPC",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_L1_RPC_URL"},
	})

	optionSettlementRPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL for settlement RPC",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL"},
		Value:   "http://localhost:8545",
	})

	optionL1ContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "l1-contract-addr",
		Usage:   "address of the L1 gateway contract",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR"},
	})

	optionSettlementContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-contract-addr",
		Usage:   "address of the settlement gateway contract",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR"},
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionPrivKeyFile,
		optionLogFmt,
		optionLogLevel,
		optionLogTags,
		optionL1RPCUrl,
		optionSettlementRPCUrl,
		optionL1ContractAddr,
		optionSettlementContractAddr,
	}

	app := &cli.App{
		Name:  "standard-bridge-relayer",
		Usage: "Entry point for relayer of mev-commit standard bridge",
		Commands: []*cli.Command{{
			Name:   "start",
			Usage:  "Start standard bridge relayer",
			Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
			Flags:  flags,
			Action: start,
		}},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

// start is the entrypoint of the cli app.
func start(c *cli.Context) error {
	logger, err := util.NewLogger(
		c.String(optionLogLevel.Name),
		c.String(optionLogFmt.Name),
		c.String(optionLogTags.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	privKeyFile, err := resolveFilePath(c.String(optionPrivKeyFile.Name))
	if err != nil {
		return fmt.Errorf("failed to get private key file path: %w", err)
	}

	privKey, err := crypto.LoadECDSA(privKeyFile)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	r, err := relayer.NewRelayer(&relayer.Options{
		Ctx:                    c.Context,
		Logger:                 logger.With("component", "relayer"),
		PrivateKey:             privKey,
		L1RPCUrl:               c.String(optionL1RPCUrl.Name),
		SettlementRPCUrl:       c.String(optionSettlementRPCUrl.Name),
		L1ContractAddr:         common.HexToAddress(c.String(optionL1ContractAddr.Name)),
		SettlementContractAddr: common.HexToAddress(c.String(optionSettlementContractAddr.Name)),
	})
	if err != nil {
		return err
	}

	interruptSigChan := make(chan os.Signal, 1)
	signal.Notify(interruptSigChan, os.Interrupt, syscall.SIGTERM)

	// Block until interrupt signal OR context's Done channel is closed.
	select {
	case <-interruptSigChan:
	case <-c.Done():
	}
	logger.Info("shutting down...")

	closedAllSuccessfully := make(chan struct{})
	go func() {
		defer close(closedAllSuccessfully)

		err := r.TryCloseAll()
		if err != nil {
			logger.Error("failed to close all routines and db connection", "error", err)
		}
	}()
	select {
	case <-closedAllSuccessfully:
	case <-time.After(5 * time.Second):
		logger.Error("failed to close all in time", "error", err)
	}

	return nil
}

func resolveFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}
