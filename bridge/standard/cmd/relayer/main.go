package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/bridge/standard/pkg/node"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "path to relayer config file",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_CONFIG"},
	}

	optionHTTPPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen on for HTTP",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_HTTP_PORT"},
		Value:   8080,
	})

	optionKeystorePath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_KEYSTORE_DIR"},
		Required: true,
	})

	optionKeystorePassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_KEYSTORE_PASSWORD"},
		Required: true,
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

	optionL1RPCUrls = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:     "l1-rpc-urls",
		Usage:    "URLs for L1 RPC",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_L1_RPC_URLS"},
		Required: true,
	})

	optionSettlementRPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL"},
		Required: true,
	})

	optionL1ContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "l1-contract-addr",
		Usage:    "address of the L1 gateway contract",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_L1_CONTRACT_ADDR"},
		Required: true,
	})

	optionSettlementContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "settlement-contract-addr",
		Usage:    "address of the settlement gateway contract",
		EnvVars:  []string{"STANDARD_BRIDGE_RELAYER_SETTLEMENT_CONTRACT_ADDR"},
		Required: true,
	})

	optionPgHost = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PG_HOST"},
		Value:   "localhost",
	})

	optionPgPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PG_PORT"},
		Value:   5432,
	})

	optionPgUser = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PG_USER"},
		Value:   "postgres",
	})

	optionPgPassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PG_PASSWORD"},
		Value:   "postgres",
	})

	optionPgDbname = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"STANDARD_BRIDGE_RELAYER_PG_DBNAME"},
		Value:   "mev_commit_bridge",
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionHTTPPort,
		optionKeystorePath,
		optionKeystorePassword,
		optionLogFmt,
		optionLogLevel,
		optionLogTags,
		optionL1RPCUrls,
		optionSettlementRPCUrl,
		optionL1ContractAddr,
		optionSettlementContractAddr,
		optionPgHost,
		optionPgPort,
		optionPgUser,
		optionPgPassword,
		optionPgDbname,
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

	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		//nolint:errcheck
		fmt.Fprintln(app.Writer, "received interrupt signal, exiting... Force exit with Ctrl+C")
		cancel()
		<-sigc
		//nolint:errcheck
		fmt.Fprintln(app.Writer, "force exiting...")
		os.Exit(1)
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		//nolint:errcheck
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

	signer, err := keysigner.NewKeystoreSigner(c.String(optionKeystorePath.Name), c.String(optionKeystorePassword.Name))
	if err != nil {
		return fmt.Errorf("failed to create keystore signer: %w", err)
	}

	nd, err := node.NewNode(&node.Options{
		Logger:                 logger,
		HTTPPort:               c.Int(optionHTTPPort.Name),
		Signer:                 signer,
		L1RPCURLs:              c.StringSlice(optionL1RPCUrls.Name),
		L1GatewayContractAddr:  common.HexToAddress(c.String(optionL1ContractAddr.Name)),
		SettlementRPCURL:       c.String(optionSettlementRPCUrl.Name),
		SettlementContractAddr: common.HexToAddress(c.String(optionSettlementContractAddr.Name)),
		PgHost:                 c.String(optionPgHost.Name),
		PgPort:                 c.Int(optionPgPort.Name),
		PgUser:                 c.String(optionPgUser.Name),
		PgPassword:             c.String(optionPgPassword.Name),
		PgDB:                   c.String(optionPgDbname.Name),
	})
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	//nolint:staticcheck
	<-c.Context.Done()

	return nd.Close()
}
