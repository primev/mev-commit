package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/primev/mev-commit/cl/redisapp"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
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
	// Initialize the logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	log := slog.New(handler)

	// Define flags
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to config file",
			EnvVars: []string{"RAPP_CONFIG"},
		},
		altsrc.NewStringFlag(&cli.StringFlag{
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
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
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
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
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
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
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
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
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
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "evm-build-delay",
			Usage:   "EVM build delay",
			EnvVars: []string{"RAPP_EVM_BUILD_DELAY"},
			Value:   200 * time.Millisecond,
		}),
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
					return startApplication(c, log)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error("Error running app", "error", err)
	}
}

func startApplication(c *cli.Context, log *slog.Logger) error {
	// Load configuration
	cfg := Config{
		InstanceID:       c.String("instance-id"),
		EthClientURL:     c.String("eth-client-url"),
		JWTSecret:        c.String("jwt-secret"),
		GenesisBlockHash: c.String("genesis-block-hash"),
		RedisAddr:        c.String("redis-addr"),
		EVMBuildDelay:    c.Duration("evm-build-delay"),
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
