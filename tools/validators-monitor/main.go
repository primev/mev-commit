package main

import (
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/primev/mev-commit/tools/validators-monitor/config"
	"github.com/primev/mev-commit/tools/validators-monitor/service"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionBeaconApiUrls = &cli.StringFlag{
		Name:    "beacon-api-url",
		Usage:   "URLs for Beacon API endpoints",
		EnvVars: []string{"BEACON_API_URL"},
		Value:   "https://ethereum-beacon-api.publicnode.com",
	}

	optionEthereumRpcUrl = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "Ethereum RPC URL",
		EnvVars: []string{"ETHEREUM_RPC_URL"},
		Value:   "https://ethereum-rpc.publicnode.com",
	}

	optionValidatorOptInContract = &cli.StringFlag{
		Name:    "validator-opt-in-contract",
		Usage:   "Validator opt-in contract address",
		EnvVars: []string{"VALIDATOR_OPT_IN_CONTRACT"},
		Value:   "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64",
	}

	optionRelayUrls = &cli.StringSliceFlag{
		Name:    "relay-urls",
		Usage:   "URLs for MEV-Boost relay APIs (comma-separated)",
		EnvVars: []string{"RELAY_URLS"},
		Value: cli.NewStringSlice(
			"https://mainnet.aestus.live",
			"https://mainnet.titanrelay.xyz",
			"https://bloxroute.max-profit.blxrbdn.com",
		),
	}

	optionSlackWebhook = &cli.StringFlag{
		Name:    "slack-webhook",
		Usage:   "Slack webhook URL for notifications",
		EnvVars: []string{"SLACK_WEBHOOK_URL"},
	}

	optionDashboardApiUrl = &cli.StringFlag{
		Name:    "dashboard-api-url",
		Usage:   "Dashboard API URL for notifications",
		EnvVars: []string{"DASHBOARD_API_URL"},
		Value:   "http://185.26.9.11:8081/",
	}

	optionTrackMissed = &cli.BoolFlag{
		Name:    "track-missed",
		Usage:   "Whether to track missed duties",
		EnvVars: []string{"TRACK_MISSED"},
		Value:   true,
	}

	optionDBEnabled = &cli.BoolFlag{
		Name:    "db-enabled",
		Usage:   "Enable database storage",
		EnvVars: []string{"DB_ENABLED"},
		Value:   false,
	}

	optionDBHost = &cli.StringFlag{
		Name:    "db-host",
		Usage:   "Database host",
		EnvVars: []string{"DB_HOST"},
		Value:   "localhost",
	}

	optionDBPort = &cli.IntFlag{
		Name:    "db-port",
		Usage:   "Database port",
		EnvVars: []string{"DB_PORT"},
		Value:   5432,
	}

	optionDBUser = &cli.StringFlag{
		Name:    "db-user",
		Usage:   "Database user",
		EnvVars: []string{"DB_USER"},
		Value:   "postgres",
	}

	optionDBPassword = &cli.StringFlag{
		Name:    "db-password",
		Usage:   "Database password",
		EnvVars: []string{"DB_PASSWORD"},
		Value:   "postgres",
	}

	optionDBName = &cli.StringFlag{
		Name:    "db-name",
		Usage:   "Database name",
		EnvVars: []string{"DB_NAME"},
		Value:   "mev_commit_validators_monitor",
	}

	optionDBSSLMode = &cli.StringFlag{
		Name:    "db-sslmode",
		Usage:   "Database SSL mode",
		EnvVars: []string{"DB_SSLMODE"},
		Value:   "disable",
	}

	optionHealthPort = &cli.IntFlag{
		Name:    "health-port",
		Usage:   "Port for health check endpoint",
		EnvVars: []string{"HEALTH_PORT"},
		Value:   9090,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"LOG_FMT"},
		Value:   "json",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"LOG_TAGS"},
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
	}
)

func main() {
	app := &cli.App{
		Name:  "validator-monitor",
		Usage: "Monitor and log Ethereum validator proposer duties",
		Flags: []cli.Flag{
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
			optionBeaconApiUrls,
			optionEthereumRpcUrl,
			optionValidatorOptInContract,
			optionTrackMissed,
			optionSlackWebhook,
			optionDashboardApiUrl,
			optionRelayUrls,
			optionHealthPort,
			optionDBEnabled,
			optionDBHost,
			optionDBPort,
			optionDBUser,
			optionDBPassword,
			optionDBName,
			optionDBSSLMode,
		},
		Action: func(c *cli.Context) error {
			// Setup logger
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			// Create configuration
			cfg := &config.Config{
				BeaconNodeURL:          c.String(optionBeaconApiUrls.Name),
				TrackMissed:            c.Bool(optionTrackMissed.Name),
				EthereumRPCURL:         c.String(optionEthereumRpcUrl.Name),
				ValidatorOptInContract: c.String(optionValidatorOptInContract.Name),
				FetchIntervalSec:       12, // Use epoch duration
				SlackWebhookURL:        c.String(optionSlackWebhook.Name),
				DashboardApiUrl:        c.String(optionDashboardApiUrl.Name),
				RelayURLs:              c.StringSlice(optionRelayUrls.Name),
				HealthPort:             c.Int(optionHealthPort.Name),
				DB: config.DBConfig{
					Enabled:  c.Bool(optionDBEnabled.Name),
					Host:     c.String(optionDBHost.Name),
					Port:     c.Int(optionDBPort.Name),
					User:     c.String(optionDBUser.Name),
					Password: c.String(optionDBPassword.Name),
					DBName:   c.String(optionDBName.Name),
					SSLMode:  c.String(optionDBSSLMode.Name),
				},
			}

			logger.Debug(
				"service config",
				"config", cfg,
			)

			s, err := service.New(cfg, logger)
			if err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}

			<-sigc
			logger.Info("shutting down...")

			return s.Close()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
