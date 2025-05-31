package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/primev/mev-commit/cl/singlenode"
	"github.com/primev/mev-commit/cl/singlenode/membernode"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const (
	categoryDebug    = "Debug"
	categoryDatabase = "Database"
	categoryMember   = "Member Node"
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

var (
	configFlag = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to YAML config file",
		EnvVars: []string{"LEADER_CONFIG"},
	}

	instanceIDFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "instance-id",
		Usage:    "Unique instance ID for this node (for logging/identification)",
		EnvVars:  []string{"LEADER_INSTANCE_ID"},
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
		EnvVars: []string{"LEADER_ETH_CLIENT_URL"},
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
		EnvVars: []string{"LEADER_JWT_SECRET"},
		Value:   "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c",
		Action: func(_ *cli.Context, s string) error {
			if len(s) != 64 {
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
		EnvVars:  []string{"MEV_COMMIT_LOG_FMT"},
		Value:    "text",
		Action:   stringInCheck("log-fmt", []string{"text", "json"}),
		Category: categoryDebug,
	})

	logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "log-level",
		Usage:    "Log level ('debug', 'info', 'warn', 'error')",
		EnvVars:  []string{"MEV_COMMIT_LOG_LEVEL"},
		Value:    "info",
		Action:   stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
		Category: categoryDebug,
	})

	logTagsFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "Comma-separated <name:value> log tags (e.g., env:prod,service:snode)",
		EnvVars: []string{"MEV_COMMIT_LOG_TAGS"},
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
		EnvVars: []string{"LEADER_EVM_BUILD_DELAY"},
		Value:   100 * time.Millisecond,
	})

	evmBuildDelayEmptyBlockFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "evm-build-delay-empty-block",
		Usage:   "Minimum time since last block to build an empty block (0 to disable skipping, e.g., '2s')",
		EnvVars: []string{"LEADER_EVM_BUILD_DELAY_EMPTY_BLOCK"},
		Value:   2 * time.Second,
	})

	priorityFeeReceiptFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "priority-fee-recipient",
		Usage:    "Ethereum address for receiving priority fees (block proposer fee)",
		EnvVars:  []string{"LEADER_PRIORITY_FEE_RECIPIENT"},
		Required: true,
		Action: func(c *cli.Context, s string) error {
			if !strings.HasPrefix(s, "0x") || len(s) != 42 {
				return fmt.Errorf("priority-fee-recipient must be a 0x-prefixed 42-character hex string")
			}
			// Basic validation
			if _, err := hex.DecodeString(s[2:]); err != nil {
				return fmt.Errorf("priority-fee-recipient is not a valid hex string: %v", err)
			}
			return nil
		},
	})

	healthAddrPortFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "health-addr",
		Usage:   "Address for health check endpoint (e.g., ':8080')",
		EnvVars: []string{"LEADER_HEALTH_ADDR"},
		Value:   ":8080",
		Action: func(_ *cli.Context, s string) error {
			if !strings.HasPrefix(s, ":") {
				return fmt.Errorf("health-addr must start with ':' (e.g., ':8080')")
			}
			// Validate port number
			portStr := s[1:] // Remove the ':'
			if port, err := strconv.Atoi(portStr); err != nil || port < 1 || port > 65535 {
				return fmt.Errorf("health-addr must be a valid port number (e.g., ':8080')")
			}
			return nil
		},
	})

	postgresDSNFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name: "postgres-dsn",
		Usage: "PostgreSQL DSN for storing payloads. If empty, saving to DB is disabled. " +
			"(e.g., 'postgres://user:pass@host:port/dbname?sslmode=disable')",
		EnvVars:  []string{"LEADER_POSTGRES_DSN"},
		Value:    "", // Default to empty, making it optional
		Category: categoryDatabase,
	})

	apiAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "api-addr",
		Usage:   "Address for member node API endpoint (e.g., ':9090'). If empty, API is disabled.",
		EnvVars: []string{"LEADER_API_ADDR"},
		Value:   ":9090",
		Action: func(_ *cli.Context, s string) error {
			if s == "" {
				return nil // Optional flag
			}
			if !strings.HasPrefix(s, ":") {
				return fmt.Errorf("api-addr must start with ':'")
			}
			return nil
		},
	})

	// Member node specific flags
	leaderAPIURLFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "leader-api-url",
		Usage:    "Leader node API URL for member nodes (e.g., 'http://leader:9090')",
		EnvVars:  []string{"MEMBER_LEADER_API_URL"},
		Category: categoryMember,
		Action: func(_ *cli.Context, s string) error {
			if s == "" {
				return nil // Will be validated in member command
			}
			if _, err := url.Parse(s); err != nil {
				return fmt.Errorf("invalid leader-api-url: %v", err)
			}
			return nil
		},
	})

	pollIntervalFlag = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:     "poll-interval",
		Usage:    "Interval for polling leader node for new payloads (e.g., '1s')",
		EnvVars:  []string{"MEMBER_POLL_INTERVAL"},
		Value:    1 * time.Second,
		Category: categoryMember,
	})
)

func main() {
	leaderFlags := []cli.Flag{
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
		healthAddrPortFlag,
		postgresDSNFlag,
		apiAddrFlag,
	}

	memberFlags := []cli.Flag{
		configFlag,
		instanceIDFlag,
		ethClientURLFlag,
		jwtSecretFlag,
		logFmtFlag,
		logLevelFlag,
		logTagsFlag,
		healthAddrPortFlag,
		leaderAPIURLFlag,
		pollIntervalFlag,
	}

	app := &cli.App{
		Name:  "snode",
		Usage: "Single-node MEV-commit application",
		Commands: []*cli.Command{
			{
				Name:  "leader",
				Usage: "Start as leader node (produces blocks)",
				Flags: leaderFlags,
				Before: altsrc.InitInputSourceWithContext(leaderFlags,
					func(c *cli.Context) (altsrc.InputSourceContext, error) {
						configFile := c.String(configFlag.Name)
						if configFile != "" {
							return altsrc.NewYamlSourceFromFile(configFile)
						}
						return &altsrc.MapInputSource{}, nil
					}),
				Action: func(c *cli.Context) error {
					return startLeaderNode(c)
				},
			},
			{
				Name:  "member",
				Usage: "Start as member node (follows leader)",
				Flags: memberFlags,
				Before: altsrc.InitInputSourceWithContext(memberFlags,
					func(c *cli.Context) (altsrc.InputSourceContext, error) {
						configFile := c.String(configFlag.Name)
						if configFile != "" {
							return altsrc.NewYamlSourceFromFile(configFile)
						}
						return &altsrc.MapInputSource{}, nil
					}),
				Action: func(c *cli.Context) error {
					return startMemberNode(c)
				},
			},
			// Keep the old "start" command for backward compatibility
			{
				Name:  "start",
				Usage: "Start as leader node (deprecated, use 'leader' instead)",
				Flags: leaderFlags,
				Before: altsrc.InitInputSourceWithContext(leaderFlags,
					func(c *cli.Context) (altsrc.InputSourceContext, error) {
						configFile := c.String(configFlag.Name)
						if configFile != "" {
							return altsrc.NewYamlSourceFromFile(configFile)
						}
						return &altsrc.MapInputSource{}, nil
					}),
				Action: func(c *cli.Context) error {
					return startLeaderNode(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(app.Writer, "Error running snode: %v\n", err)
		os.Exit(1)
	}
}

func startLeaderNode(c *cli.Context) error {
	logger, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		c.String(logTagsFlag.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	logger = logger.With("app", "snode", "role", "leader")

	cfg := singlenode.Config{
		InstanceID:               c.String(instanceIDFlag.Name),
		EthClientURL:             c.String(ethClientURLFlag.Name),
		JWTSecret:                c.String(jwtSecretFlag.Name),
		EVMBuildDelay:            c.Duration(evmBuildDelayFlag.Name),
		EVMBuildDelayEmptyBlocks: c.Duration(evmBuildDelayEmptyBlockFlag.Name),
		PriorityFeeReceipt:       c.String(priorityFeeReceiptFlag.Name),
		HealthAddr:               c.String(healthAddrPortFlag.Name),
		PostgresDSN:              c.String(postgresDSNFlag.Name),
		APIAddr:                  c.String(apiAddrFlag.Name),
	}

	logger.Info("Starting leader node with configuration", "config", cfg)

	// Create a root context that can be cancelled for graceful shutdown
	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer rootCancel()

	snode, err := singlenode.NewSingleNodeApp(rootCtx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize SingleNodeApp", "error", err)
		return err
	}

	snode.Start()

	<-rootCtx.Done()

	logger.Info("Shutdown signal received, stopping leader node...")
	snode.Stop()

	logger.Info("Leader node shutdown completed.")
	return nil
}

func startMemberNode(c *cli.Context) error {
	leaderURL := c.String(leaderAPIURLFlag.Name)
	if leaderURL == "" {
		return fmt.Errorf("leader-api-url is required for member nodes")
	}

	logger, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		c.String(logTagsFlag.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	logger = logger.With("app", "snode", "role", "member")

	cfg := membernode.Config{
		InstanceID:   c.String(instanceIDFlag.Name),
		LeaderAPIURL: leaderURL,
		EthClientURL: c.String(ethClientURLFlag.Name),
		JWTSecret:    c.String(jwtSecretFlag.Name),
		HealthAddr:   c.String(healthAddrPortFlag.Name),
		PollInterval: c.Duration(pollIntervalFlag.Name),
	}

	logger.Info("Starting member node", "config", cfg)

	// Create a root context that can be cancelled for graceful shutdown
	rootCtx, rootCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer rootCancel()

	memberNode, err := membernode.NewMemberNodeApp(rootCtx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize MemberNodeApp", "error", err)
		return err
	}

	memberNode.Start()

	<-rootCtx.Done()

	logger.Info("Shutdown signal received, stopping member node...")
	memberNode.Stop()

	logger.Info("Member node shutdown completed.")
	return nil
}
