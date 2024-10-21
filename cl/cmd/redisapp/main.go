package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "eth-client-url",
			Usage:   "Ethereum client URL",
			EnvVars: []string{"RAPP_ETH_CLIENT_URL"},
			Value:   "http://localhost:8551",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "jwt-secret",
			Usage:   "JWT secret for Ethereum client",
			EnvVars: []string{"RAPP_JWT_SECRET"},
			Value:   "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "genesis-block-hash",
			Usage:   "Genesis block hash",
			EnvVars: []string{"RAPP_GENESIS_BLOCK_HASH"},
			Value:   "c9810c36e1e8bb2adaa677338b43870f73c3a39abebdffb582a668ca63e523d2",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "redis-addr",
			Usage:   "Redis address",
			EnvVars: []string{"RAPP_REDIS_ADDR"},
			Value:   "127.0.0.1:7001",
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:    "evm-build-delay",
			Usage:   "EVM build delay",
			EnvVars: []string{"RAPP_EVM_BUILD_DELAY"},
			Value:   time.Second,
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

	if cfg.InstanceID == "" {
		return fmt.Errorf("instance-id is required")
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
