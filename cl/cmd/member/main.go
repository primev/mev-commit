package main

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/primev/mev-commit/cl/member"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	configFlag = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to config file",
		EnvVars: []string{"MEMBER_CONFIG"},
	}

	clientIDFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "client-id",
		Usage:    "Unique client ID for this member",
		EnvVars:  []string{"MEMBER_CLIENT_ID"},
		Required: true,
	})

	relayerAddrFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "relayer-addr",
		Usage:    "Relayer address",
		EnvVars:  []string{"MEMBER_RELAYER_ADDR"},
		Required: true,
		Action: func(_ *cli.Context, s string) error {
			if _, err := url.Parse(s); err != nil {
				return fmt.Errorf("invalid relayer-addr: %v", err)
			}
			return nil
		},
	})

	ethClientURLFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "eth-client-url",
		Usage:   "Ethereum client URL",
		EnvVars: []string{"MEMBER_ETH_CLIENT_URL"},
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
		EnvVars: []string{"MEMBER_JWT_SECRET"},
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

	logFmtFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "Log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MEMBER_LOG_FMT"},
		Value:   "text",
	})

	logLevelFlag = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "Log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"MEMBER_LOG_LEVEL"},
		Value:   "info",
	})
)

type Config struct {
    ClientID      string
    RelayerAddr   string
    EthClientURL  string
    JWTSecret     string
}

func main() {
	flags := []cli.Flag{
		configFlag,
		clientIDFlag,
		relayerAddrFlag,
		ethClientURLFlag,
		jwtSecretFlag,
		logFmtFlag,
		logLevelFlag,
	}

	app := &cli.App{
		Name:  "memberclient",
		Usage: "Start the member client",
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
			return startMemberClient(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error running member client:", err)
	}
}

func startMemberClient(c *cli.Context) error {
	log, err := util.NewLogger(
		c.String(logLevelFlag.Name),
		c.String(logFmtFlag.Name),
		"", // No log tags
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	cfg := Config{
		ClientID:     c.String(clientIDFlag.Name),
		RelayerAddr:  c.String(relayerAddrFlag.Name),
		EthClientURL: c.String(ethClientURLFlag.Name),
		JWTSecret:    c.String(jwtSecretFlag.Name),
	}
	
	log.Info("Starting member client with configuration", "config", cfg)

	// Initialize the MemberClient
	memberClient, err := member.NewMemberClient(cfg.ClientID, cfg.RelayerAddr, cfg.EthClientURL, cfg.JWTSecret, log)
	if err != nil {
		log.Error("Failed to initialize MemberClient", "error", err)
		return err
	}

	ctx, stop := signal.NotifyContext(c.Context, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the member client
	go func() {
		if err := memberClient.Run(ctx); err != nil {
			log.Error("Member client exited with error", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	log.Info("Member client shutdown completed")
	return nil
}
