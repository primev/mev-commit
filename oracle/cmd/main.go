package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/primevprotocol/mev-commit/contracts-abi/config"
	"github.com/primevprotocol/mev-commit/oracle/pkg/keysigner"
	"github.com/primevprotocol/mev-commit/oracle/pkg/node"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const (
	defaultHTTPPort  = 8080
	defaultConfigDir = "~/.mev-commit-oracle"
	defaultKeyFile   = "key"
	defaultKeystore  = "keystore"
)

var (
	portCheck = func(c *cli.Context, p int) error {
		if p < 0 || p > 65535 {
			return fmt.Errorf("invalid port number %d, expected 0 <= port <= 65535", p)
		}
		return nil
	}

	stringInCheck = func(flag string, opts []string) func(c *cli.Context, s string) error {
		return func(c *cli.Context, s string) error {
			if !slices.Contains(opts, s) {
				return fmt.Errorf("invalid %s option %q, expected one of %s", flag, s, strings.Join(opts, ", "))
			}
			return nil
		}
	}
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "path to config file",
		EnvVars: []string{"MEV_ORACLE_CONFIG"},
	}

	optionPrivKeyFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "priv-key-file",
		Usage:   "path to private key file",
		EnvVars: []string{"MEV_ORACLE_PRIV_KEY_FILE"},
		Value:   filepath.Join(defaultConfigDir, defaultKeyFile),
	})

	optionHTTPPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen on for HTTP requests",
		EnvVars: []string{"MEV_ORACLE_HTTP_PORT"},
		Value:   defaultHTTPPort,
		Action:  portCheck,
	})

	optionLogFmt = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MEV_ORACLE_LOG_FMT"},
		Value:   "text",
		Action:  stringInCheck("log-fmt", []string{"text", "json"}),
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"MEV_ORACLE_LOG_LEVEL"},
		Value:   "info",
		Action:  stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
	})

	optionLogTags = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"MEV_ORACLE_LOG_TAGS"},
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
		EnvVars: []string{"MEV_ORACLE_L1_RPC_URL"},
	})

	optionSettlementRPCUrl = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL for settlement RPC",
		EnvVars: []string{"MEV_ORACLE_SETTLEMENT_RPC_URL"},
		Value:   "http://localhost:8545",
	})

	optionOracleContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "oracle-contract-addr",
		Usage:   "address of the oracle contract",
		EnvVars: []string{"MEV_ORACLE_ORACLE_CONTRACT_ADDR"},
		Value:   contracts.TestnetContracts.Oracle,
	})

	optionPreconfContractAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "preconf-contract-addr",
		Usage:   "address of the preconf contract",
		EnvVars: []string{"MEV_ORACLE_PRECONF_CONTRACT_ADDR"},
		Value:   contracts.TestnetContracts.PreconfCommitmentStore,
	})

	optionPgHost = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"MEV_ORACLE_PG_HOST"},
		Value:   "localhost",
	})

	optionPgPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"MEV_ORACLE_PG_PORT"},
		Value:   5432,
	})

	optionPgUser = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"MEV_ORACLE_PG_USER"},
		Value:   "postgres",
	})

	optionPgPassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"MEV_ORACLE_PG_PASSWORD"},
		Value:   "postgres",
	})

	optionPgDbname = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"MEV_ORACLE_PG_DBNAME"},
		Value:   "mev_oracle",
	})

	optionLaggerdMode = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "laggerd-mode",
		Usage:   "No of blocks to lag behind for L1 chain",
		EnvVars: []string{"MEV_ORACLE_LAGGERD_MODE"},
		Value:   0,
	})

	optionOverrideWinners = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:    "override-winners",
		Usage:   "Override winners for testing",
		EnvVars: []string{"MEV_ORACLE_OVERRIDE_WINNERS"},
	})

	optionKeystorePassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-password",
		Usage:   "use to access keystore",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PASSWORD"},
	})

	optionKeystorePath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-path",
		Usage:   "path to keystore location",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PATH"},
		Value:   filepath.Join(defaultConfigDir, defaultKeystore),
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionPrivKeyFile,
		optionHTTPPort,
		optionLogFmt,
		optionLogLevel,
		optionLogTags,
		optionL1RPCUrl,
		optionSettlementRPCUrl,
		optionOracleContractAddr,
		optionPreconfContractAddr,
		optionPgHost,
		optionPgPort,
		optionPgUser,
		optionPgPassword,
		optionPgDbname,
		optionLaggerdMode,
		optionOverrideWinners,
		optionKeystorePath,
		optionKeystorePassword,
	}
	app := &cli.App{
		Name:  "mev-oracle",
		Usage: "Entry point for mev-oracle",
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the mev-oracle node",
				Flags:  flags,
				Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
				Action: func(c *cli.Context) error {
					return initializeApplication(c)
				},
			},
		}}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

func initializeApplication(c *cli.Context) error {
	if err := verifyKeystorePasswordPresence(c); err != nil {
		return err
	}
	if err := launchOracleWithConfig(c); err != nil {
		return err
	}
	return nil
}

// verifyKeystorePasswordPresence checks for the presence of a keystore password.
// it returns error, if keystore path is set and keystore password is not
func verifyKeystorePasswordPresence(c *cli.Context) error {
	if c.IsSet(optionKeystorePath.Name) && !c.IsSet(optionKeystorePassword.Name) {
		return cli.Exit("Password for encrypted keystore is missing", 1)
	}
	return nil
}

// launchOracleWithConfig configures and starts the oracle based on the CLI context or config.yaml file.
func launchOracleWithConfig(c *cli.Context) error {
	logger, err := newLogger(
		c.String(optionLogLevel.Name),
		c.String(optionLogFmt.Name),
		c.String(optionLogTags.Name),
		c.App.Writer,
	)

	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	keySigner, err := setupKeySigner(c)
	if err != nil {
		return fmt.Errorf("failed to setup key signer: %w", err)
	}
	logger.Info("key signer account", "address", keySigner.GetAddress().Hex(), "url", keySigner.String())

	nd, err := node.NewNode(&node.Options{
		Logger:              logger,
		KeySigner:           keySigner,
		HTTPPort:            c.Int(optionHTTPPort.Name),
		L1RPCUrl:            c.String(optionL1RPCUrl.Name),
		SettlementRPCUrl:    c.String(optionSettlementRPCUrl.Name),
		OracleContractAddr:  common.HexToAddress(c.String(optionOracleContractAddr.Name)),
		PreconfContractAddr: common.HexToAddress(c.String(optionPreconfContractAddr.Name)),
		PgHost:              c.String(optionPgHost.Name),
		PgPort:              c.Int(optionPgPort.Name),
		PgUser:              c.String(optionPgUser.Name),
		PgPassword:          c.String(optionPgPassword.Name),
		PgDbname:            c.String(optionPgDbname.Name),
		LaggerdMode:         c.Int(optionLaggerdMode.Name),
		OverrideWinners:     c.StringSlice(optionOverrideWinners.Name),
	})
	if err != nil {
		return fmt.Errorf("failed starting node: %w", err)
	}

	<-c.Done()
	fmt.Fprintf(c.App.Writer, "shutting down...\n")
	closed := make(chan struct{})

	go func() {
		defer close(closed)

		err := nd.Close()
		if err != nil {
			logger.Error("failed to close node", "error", err)
		}
	}()

	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		logger.Error("failed to close node in time", "error", err)
	}

	return nil
}

// newLogger initializes a *slog.Logger with specified level, format, and sink.
//   - lvl: string representation of slog.Level
//   - logFmt: format of the log output: "text", "json", "none" defaults to "json"
//   - tags: comma-separated list of <name:value> pairs that will be inserted into each log line
//   - sink: destination for log output (e.g., os.Stdout, file)
//
// Returns a configured *slog.Logger on success or nil on failure.
func newLogger(lvl, logFmt, tags string, sink io.Writer) (*slog.Logger, error) {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(lvl)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var (
		handler slog.Handler
		options = &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
	)
	switch logFmt {
	case "text":
		handler = slog.NewTextHandler(sink, options)
	case "json", "none":
		handler = slog.NewJSONHandler(sink, options)
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	logger := slog.New(handler)

	if tags == "" {
		return logger, nil
	}

	var args []any
	for i, p := range strings.Split(tags, ",") {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag at index %d", i)
		}
		args = append(args, strings.ToValidUTF8(kv[0], "�"), strings.ToValidUTF8(kv[1], "�"))
	}

	return logger.With(args...), nil
}

func setupKeySigner(c *cli.Context) (keysigner.KeySigner, error) {
	if c.IsSet(optionKeystorePath.Name) {
		return keysigner.NewKeystoreSigner(c.String(optionKeystorePath.Name), c.String(optionKeystorePassword.Name))
	}
	return keysigner.NewPrivateKeySigner(c.String(optionPrivKeyFile.Name))
}
