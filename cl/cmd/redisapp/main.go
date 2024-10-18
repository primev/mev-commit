package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/primev/mev-commit-geth-cl/logger"
	app "github.com/primev/mev-commit-geth-cl/redisapp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to config file",
		EnvVars: []string{"RAPP_CONFIG"},
	}

	optionInstanceID = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "instance-id",
		Usage:   "Unique instance ID for this node",
		EnvVars: []string{"RAPP_INSTANCE_ID"},
	})

	optionEthClientURL = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "eth-client-url",
		Usage:   "Ethereum client URL",
		EnvVars: []string{"RAPP_ETH_CLIENT_URL"},
		Value:   "http://localhost:8551",
	})

	optionJWTSecret = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "jwt-secret",
		Usage:   "JWT secret for Ethereum client",
		EnvVars: []string{"RAPP_JWT_SECRET"},
		Value:   "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c",
	})

	optionGenesisBlockHash = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "genesis-block-hash",
		Usage:   "Genesis block hash",
		EnvVars: []string{"RAPP_GENESIS_BLOCK_HASH"},
		Value:   "c9810c36e1e8bb2adaa677338b43870f73c3a39abebdffb582a668ca63e523d2",
	})

	optionRedisAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "redis-addr",
		Usage:   "Redis address",
		EnvVars: []string{"RAPP_REDIS_ADDR"},
		Value:   "127.0.0.1:7001",
	})
)

func main() {
	// Initialize your logger
	baseLogger := logrus.New()
	baseLogger.SetLevel(logrus.InfoLevel)
	baseLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	logger := &logger.LogrusWrapper{Logger: baseLogger}

	// Define flags
	flags := []cli.Flag{
		optionConfig,
		optionInstanceID,
		optionEthClientURL,
		optionJWTSecret,
		optionGenesisBlockHash,
		optionRedisAddr,
	}

	// Create the app
	app := &cli.App{
		Name:  "rapp",
		Usage: "Entry point for rapp",
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the rapp node",
				Flags:  flags,
				Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
				Action: func(c *cli.Context) error {
					return initializeApplication(c, logger)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error("Error running app", "error", err)
	}
}

func initializeApplication(c *cli.Context, logger *logger.LogrusWrapper) error {
	instanceID := c.String("instance-id")
	if instanceID == "" {
		// Generate a UUID if INSTANCE_ID is not set
		generatedUUID, err := uuid.NewRandom()
		if err != nil {
			logger.Error("Failed to generate UUID for InstanceID", "error", err)
			return err
		}
		instanceID = generatedUUID.String()
		logger.Info("Generated new InstanceID", "InstanceID", instanceID)
	}

	genesisBlockHash := c.String("genesis-block-hash")
	if genesisBlockHash == "" {
		logger.Warn("GENESIS_BLOCK_HASH is not set, using default value")
		genesisBlockHash = "c9810c36e1e8bb2adaa677338b43870f73c3a39abebdffb582a668ca63e523d2"
	}

	ethClientURL := c.String("eth-client-url")
	if ethClientURL == "" {
		logger.Warn("ETH_CLIENT_URL is not set, using default value")
		ethClientURL = "http://localhost:8551"
	}

	jwtSecret := c.String("jwt-secret")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET is not set, using default value")
		jwtSecret = "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c"
	}

	redisAddr := c.String("redis-addr")
	if redisAddr == "" {
		logger.Warn("REDIS_ADDR is not set, using default value")
		redisAddr = "127.0.0.1:7001"
	}

	// Initialize RappChain (MevCommitChain)
	rappChain, err := app.NewMevCommitChain(
		instanceID,
		ethClientURL,
		jwtSecret,
		genesisBlockHash,
		logger,
		redisAddr,
	)

	if err != nil {
		logger.Error("Failed to initialize RappChain", "error", err)
		return err
	}

	// Handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-signalChan

	// Call the Stop function to gracefully shutdown the app
	rappChain.Stop()

	logger.Info("Application shutdown completed")
	return nil
}
