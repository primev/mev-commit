package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/primev/mev-commit/cl/redisapp"
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

var (
    configFlag = &cli.StringFlag{
        Name:    "config",
        Usage:   "Path to config file",
        EnvVars: []string{"RAPP_CONFIG"},
    }

    instanceIDFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:     "instance-id",
        Usage:    "Unique instance ID for this node",
        EnvVars:  []string{"RAPP_INSTANCE_ID"},
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
        Usage:   "Ethereum client URL",
        EnvVars: []string{"RAPP_ETH_CLIENT_URL"},
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
		Usage:   "JWT secret for Ethereum client",
        EnvVars: []string{"RAPP_JWT_SECRET"},
        Value:   "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c",
        Action: func(_ *cli.Context, s string) error {
            if len(s) != 64 {
                return fmt.Errorf("invalid jwt-secret: must be 64 hex characters")
            }
            if _, err := hex.DecodeString(s); err != nil {
                return fmt.Errorf("invalid jwt-secret: %v", err)
            }
            return nil
        },
    })

    genesisBlockHashFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:    "genesis-block-hash",
        Usage:   "Genesis block hash",
        EnvVars: []string{"RAPP_GENESIS_BLOCK_HASH"},
        Value:   "dfc7fa546e1268f5bb65b9ec67759307d2435ad1bf609307c7c306e9bb0edcde",
        Action: func(_ *cli.Context, s string) error {
            if len(s) != 64 {
                return fmt.Errorf("invalid genesis-block-hash: must be 64 hex characters")
            }
            if _, err := hex.DecodeString(s); err != nil {
                return fmt.Errorf("invalid genesis-block-hash: %v", err)
            }
            return nil
        },
    })

    redisAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:    "redis-addr",
        Usage:   "Redis address",
        EnvVars: []string{"RAPP_REDIS_ADDR"},
        Value:   "127.0.0.1:7001",
        Action: func(_ *cli.Context, s string) error {
            host, port, err := net.SplitHostPort(s)
            if err != nil {
                return fmt.Errorf("invalid redis-addr: %v", err)
            }
            if net.ParseIP(host) == nil {
                return fmt.Errorf("invalid redis-addr: invalid IP address")
            }
            if p, err := strconv.Atoi(port); err != nil || p <= 0 || p > 65535 {
                return fmt.Errorf("invalid redis-addr: invalid port number")
            }
            return nil
        },
    })

    logFmtFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:     "log-fmt",
        Usage:    "Log format to use, options are 'text' or 'json'",
        EnvVars:  []string{"MEV_COMMIT_LOG_FMT"},
        Value:    "text",
        Action:   stringInCheck("log-fmt", []string{"text", "json"}),
        Category: categoryDebug,
    })

    logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:     "log-level",
        Usage:    "Log level to use, options are 'debug', 'info', 'warn', 'error'",
        EnvVars:  []string{"MEV_COMMIT_LOG_LEVEL"},
        Value:    "info",
        Action:   stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
        Category: categoryDebug,
    })

    logTagsFlag = altsrc.NewStringFlag(&cli.StringFlag{
        Name:    "log-tags",
        Usage:   "Log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
        EnvVars: []string{"MEV_COMMIT_LOG_TAGS"},
        Action: func(ctx *cli.Context, s string) error {
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
        Usage:   "EVM build delay",
        EnvVars: []string{"RAPP_EVM_BUILD_DELAY"},
        Value:   200 * time.Millisecond,
    })
)

type Config struct {
	InstanceID       string
	EthClientURL     string
	JWTSecret        string
	GenesisBlockHash string
	RedisAddr        string
	EVMBuildDelay    time.Duration
}

func main() {
	flags := []cli.Flag{
        configFlag,
        instanceIDFlag,
        ethClientURLFlag,
        jwtSecretFlag,
        genesisBlockHashFlag,
        redisAddrFlag,
        logFmtFlag,
        logLevelFlag,
        logTagsFlag,
        evmBuildDelayFlag,
    }

	// Create the app
	app := &cli.App{
		Name:  "rapp",
		Usage: "Entry point for rapp",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "Start the rapp node",
				Flags: flags,
				Before: altsrc.InitInputSourceWithContext(flags,
					func(c *cli.Context) (altsrc.InputSourceContext, error) {
						configFile := c.String("config")
						if configFile != "" {
							return altsrc.NewYamlSourceFromFile(configFile)
						}
						return &altsrc.MapInputSource{}, nil
					}),
				Action: func(c *cli.Context) error {
					return startApplication(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(app.Writer, "Error running app", "error", err)
	}
}

func startApplication(c *cli.Context) error {
	log, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		c.String(logTagsFlag.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	// Load configuration
	cfg := Config{
		InstanceID:       c.String(instanceIDFlag.Name),
		EthClientURL:     c.String(ethClientURLFlag.Name),
		JWTSecret:        c.String(jwtSecretFlag.Name),
		GenesisBlockHash: c.String(genesisBlockHashFlag.Name),
		RedisAddr:        c.String(redisAddrFlag.Name),
		EVMBuildDelay:    c.Duration(evmBuildDelayFlag.Name),
	}

	log.Info("Starting application with configuration", "config", cfg)

	// Initialize the MevCommitChain
	rappChain, err := redisapp.NewMevCommitChain(
		cfg.InstanceID,
		cfg.EthClientURL,
		cfg.JWTSecret,
		cfg.GenesisBlockHash,
		log,
		cfg.RedisAddr,
		cfg.EVMBuildDelay,
	)
	if err != nil {
		log.Error("Failed to initialize RappChain", "error", err)
		return err
	}

	ctx, stop := signal.NotifyContext(c.Context, os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	rappChain.Stop()

	log.Info("Application shutdown completed")
	return nil
}
