package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/primev/mev-commit/cl/singlenode"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const (
	categoryDebug = "Debug"
)

var (
	stringInCheck = func(flag string, opts []string) func(c *cli.Context, p string) error {
		return func(c *cli.Context, p string) error {
			if !slices.Contains(opts, p) {
				return fmt.Errorf("invalid %s option %q, expected one of %s", flag, p, strings.Join(opts, ", "))
			}
			return nil
		}
	}
)

// CLI Flags (subset of redisapp, without Redis-specific ones)
var (
	configFlag = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to YAML config file",
		EnvVars: []string{"SNODE_CONFIG"},
	}

	instanceIDFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "instance-id",
		Usage:    "Unique instance ID for this node (for logging/identification)",
		EnvVars:  []string{"SNODE_INSTANCE_ID"},
		Required: true,
		Action: func(_ *cli.Context, s string) error {
			if s == "" {
				return fmt.Errorf("instance-id is required")
			}
			return nil
		},
	})

	ethClientURLFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "eth-client-url",
		Usage:   "Ethereum Execution client Engine API URL (e.g., http://localhost:8551)",
		EnvVars: []string{"SNODE_ETH_CLIENT_URL"},
		Value:   "http://localhost:8551",
		Action: func(_ *cli.Context, s string) error {
			if _, err := url.Parse(s); err != nil {
				return fmt.Errorf("invalid eth-client-url: %v", err)
			}
			return nil
		},
	})

	jwtSecretFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "jwt-secret",
		Usage:   "Hex-encoded JWT secret for Ethereum Execution client Engine API",
		EnvVars: []string{"SNODE_JWT_SECRET"},
		// Example default, replace with secure generation or require user input
		Value: "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c",
		Action: func(_ *cli.Context, s string) error {
			if len(s) != 64 { // 32 bytes = 64 hex characters
				return fmt.Errorf("invalid jwt-secret: must be 64 hex characters (32 bytes)")
			}
			if _, err := hex.DecodeString(s); err != nil {
				return fmt.Errorf("invalid jwt-secret: failed to decode hex: %v", err)
			}
			return nil
		},
	})

	logFmtFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "log-fmt",
		Usage:    "Log format ('text' or 'json')",
		EnvVars:  []string{"MEV_COMMIT_LOG_FMT"}, // Keep consistent env var if desired
		Value:    "text",
		Action:   stringInCheck("log-fmt", []string{"text", "json"}),
		Category: categoryDebug,
	})

	logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "log-level",
		Usage:    "Log level ('debug', 'info', 'warn', 'error')",
		EnvVars:  []string{"MEV_COMMIT_LOG_LEVEL"}, // Keep consistent
		Value:    "info",
		Action:   stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
		Category: categoryDebug,
	})

	logTagsFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "Comma-separated <name:value> log tags (e.g., env:prod,service:snode)",
		EnvVars: []string{"MEV_COMMIT_LOG_TAGS"}, // Keep consistent
		Action: func(ctx *cli.Context, s string) error {
			if s == "" {
				return nil
			}
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
		Category: categoryDebug,
	})

	evmBuildDelayFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "evm-build-delay",
		Usage:   "Delay after initiating payload construction before calling getPayload (e.g., '200ms')",
		EnvVars: []string{"SNODE_EVM_BUILD_DELAY"},
		Value:   200 * time.Millisecond,
	})

	evmBuildDelayEmptyBlockFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "evm-build-delay-empty-block",
		Usage:   "Minimum time since last block to build an empty block (0 to disable skipping, e.g., '2s')",
		EnvVars: []string{"SNODE_EVM_BUILD_DELAY_EMPTY_BLOCK"},
		Value:   2 * time.Second,
	})

	priorityFeeReceiptFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "priority-fee-recipient", // Changed flag name for clarity
		Usage:   "Ethereum address for receiving priority fees (block proposer fee)",
		EnvVars: []string{"SNODE_PRIORITY_FEE_RECIPIENT"},
		// Value:   "0xYourFeeRecipientAddressHere", // Require this or ensure a safe default/handling
		Required: true, // Making this required is safer
		Action: func(c *cli.Context, s string) error {
			if !strings.HasPrefix(s, "0x") || len(s) != 42 {
				return fmt.Errorf("priority-fee-recipient must be a 0x-prefixed 42-character hex string")
			}
			// Basic validation, more robust hex address validation could be added
			if _, err := hex.DecodeString(s[2:]); err != nil {
				return fmt.Errorf("priority-fee-recipient is not a valid hex string: %v", err)
			}
			return nil
		},
	})
)

func main() {
	flags := []cli.Flag{
		configFlag,
		instanceIDFlag,
		ethClientURLFlag,
		jwtSecretFlag,
		logFmtFlag,
		logLevelFlag,
		logTagsFlag,
		evmBuildDelayFlag,
		evmBuildDelayEmptyBlockFlag,
		priorityFeeReceiptFlag,
	}

	app := &cli.App{
		Name:  "snode",
		Usage: "Single-node MEV-commit application",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the snode node",
				Flags: flags,
				Before: altsrc.InitInputSourceWithContext(flags,
					func(c *cli.Context) (altsrc.InputSourceContext, error) {
						configFile := c.String(configFlag.Name)
						if configFile != "" {
							return altsrc.NewYamlSourceFromFile(configFile)
						}
						return &altsrc.MapInputSource{}, nil // Empty source if no config file
					}),
				Action: func(c *cli.Context) error {
					return startSingleNodeApplication(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		// Use app.Writer for logging consistency if logger not yet initialized
		fmt.Fprintf(app.Writer, "Error running snode: %v\n", err)
		os.Exit(1)
	}
}

func startSingleNodeApplication(c *cli.Context) error {
	logger, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		c.String(logTagsFlag.Name),
		c.App.Writer, // Use CLI app's writer for logs
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	logger = logger.With("app", "snode")

	cfg := singlenode.Config{
		InstanceID:               c.String(instanceIDFlag.Name),
		EthClientURL:             c.String(ethClientURLFlag.Name),
		JWTSecret:                c.String(jwtSecretFlag.Name),
		EVMBuildDelay:            c.Duration(evmBuildDelayFlag.Name),
		EVMBuildDelayEmptyBlocks: c.Duration(evmBuildDelayEmptyBlockFlag.Name),
		PriorityFeeReceipt:       c.String(priorityFeeReceiptFlag.Name),
	}

	logger.Info("Starting snode with configuration", "config", cfg) // Be careful logging sensitive parts of config

	// Create a root context that can be cancelled for graceful shutdown
	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer rootCancel()

	snode, err := singlenode.NewSingleNodeApp(rootCtx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize SingleNodeApp", "error", err)
		return err
	}

	snode.Start() // Start the application's main loop

	// Wait for the application context to be done (e.g., OS signal)
	<-rootCtx.Done()

	logger.Info("Shutdown signal received, stopping snode...")
	snode.Stop() // Initiate graceful shutdown of the application

	logger.Info("SRApp shutdown completed.")
	return nil
}
