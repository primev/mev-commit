package main

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/tools/preconf-rpc/service"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port for the HTTP server",
		EnvVars: []string{"PRECONF_RPC_HTTP_PORT"},
		Value:   8080,
	}

	optionKeystorePath = &cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"PRECONF_RPC_KEYSTORE_DIR"},
		Required: true,
	}

	optionKeystorePassword = &cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"PRECONF_RPC_KEYSTORE_PASSWORD"},
		Required: true,
	}

	optionPgHost = &cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"PRECONF_RPC_PG_HOST"},
		Value:   "localhost",
	}

	optionPgPort = &cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"PRECONF_RPC_PG_PORT"},
		Value:   5432,
	}

	optionPgUser = &cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"PRECONF_RPC_PG_USER"},
		Value:   "postgres",
	}

	optionPgPassword = &cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"PRECONF_RPC_PG_PASSWORD"},
		Value:   "postgres",
	}

	optionPgDbname = &cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"PRECONF_RPC_PG_DBNAME"},
		Value:   "mev_oracle",
	}

	optionPgSSL = &cli.BoolFlag{
		Name:    "pg-ssl",
		Usage:   "use SSL for PostgreSQL connection",
		EnvVars: []string{"PRECONF_RPC_PG_SSL"},
		Value:   false,
	}

	optionL1RPCHTTPUrl = &cli.StringFlag{
		Name:     "l1-rpc-http-url",
		Usage:    "HTTP URL for L1 RPC",
		EnvVars:  []string{"PRECONF_RPC_L1_RPC_HTTP_URL"},
		Required: true,
	}

	optionL1ReceiptsRPCUrl = &cli.StringFlag{
		Name:    "l1-receipts-rpc-url",
		Usage:   "URL for L1 receipts RPC, if different from L1 RPC HTTP URL",
		EnvVars: []string{"PRECONF_RPC_L1_RECEIPTS_RPC_URL"},
	}

	optionL1RPCWSUrl = &cli.StringFlag{
		Name:     "l1-rpc-ws-url",
		Usage:    "Websocket URL for L1 RPC",
		EnvVars:  []string{"PRECONF_RPC_L1_RPC_WS_URL"},
		Required: true,
	}

	optionSettlementRPCUrl = &cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"PRECONF_RPC_SETTLEMENT_RPC_URL"},
		Required: true,
	}

	optionBidderRPCUrl = &cli.StringFlag{
		Name:     "bidder-rpc-url",
		Usage:    "URL for mev-commit bidder RPC",
		EnvVars:  []string{"PRECONF_RPC_BIDDER_RPC_URL"},
		Required: true,
	}

	optionL1ContractAddr = &cli.StringFlag{
		Name:     "l1-contract-addr",
		Usage:    "address of the L1 gateway contract",
		EnvVars:  []string{"PRECONF_RPC_L1_CONTRACT_ADDR"},
		Required: true,
	}

	optionSettlementThreshold = &cli.StringFlag{
		Name:    "settlement-threshold",
		Usage:   "Minimum threshold for settlement chain balance",
		EnvVars: []string{"PRECONF_RPC_SETTLEMENT_THRESHOLD"},
		Value:   "2000000000000000000", // 2 ETH
	}

	optionSettlementTopup = &cli.StringFlag{
		Name:    "settlement-topup",
		Usage:   "topup for settlement",
		EnvVars: []string{"PRECONF_RPC_SETTLEMENT_TOPUP"},
		Value:   "2100000000000000000", // 2.1 ETH
	}

	optionBidderThreshold = &cli.StringFlag{
		Name:    "bidder-threshold",
		Usage:   "threshold for bidder balance on settlement chain",
		EnvVars: []string{"PRECONF_RPC_BIDDER_THRESHOLD"},
		Value:   "100000000000000000", // 0.1 ETH
	}

	optionBidderTopup = &cli.StringFlag{
		Name:    "bidder-topup",
		Usage:   "topup for bidder",
		EnvVars: []string{"PRECONF_RPC_BIDDER_TOPUP"},
		Value:   "110000000000000000", // 0.11 ETH
	}

	optionTargetDepositAmount = &cli.StringFlag{
		Name:    "target-deposit-amount",
		Usage:   "target deposit amount",
		EnvVars: []string{"PRECONF_RPC_TARGET_DEPOSIT_AMOUNT"},
		Value:   "100000000000000000", // 0.1 ETH
	}

	optionGasTipCap = &cli.StringFlag{
		Name:    "gas-tip-cap",
		Usage:   "gas tip cap",
		EnvVars: []string{"PRECONF_RPC_GAS_TIP_CAP"},
		Value:   "50000000", // 0.05 gWEI
	}

	optionGasFeeCap = &cli.StringFlag{
		Name:    "gas-fee-cap",
		Usage:   "gas fee cap",
		EnvVars: []string{"PRECONF_RPC_GAS_FEE_CAP"},
		Value:   "60000000", // 0.06 gWEI
	}

	optionSettlementContractAddr = &cli.StringFlag{
		Name:     "settlement-contract-addr",
		Usage:    "address of the settlement gateway contract",
		EnvVars:  []string{"PRECONF_RPC_SETTLEMENT_CONTRACT_ADDR"},
		Required: true,
	}

	optionDepositAddress = &cli.StringFlag{
		Name:     "deposit-address",
		Usage:    "address to deposit funds to",
		EnvVars:  []string{"PRECONF_RPC_DEPOSIT_ADDRESS"},
		Required: true,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid deposit address: %s", s)
			}
			return nil
		},
	}

	optionBridgeAddress = &cli.StringFlag{
		Name:     "bridge-address",
		Usage:    "address to bridge to",
		EnvVars:  []string{"PRECONF_RPC_BRIDGE_ADDRESS"},
		Required: true,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid bridge address: %s", s)
			}
			return nil
		},
	}

	optionBlocknativeAPIKey = &cli.StringFlag{
		Name:    "blocknative-api-key",
		Usage:   "Blocknative API key for transaction pricing",
		EnvVars: []string{"PRECONF_RPC_BLOCKNATIVE_API_KEY"},
		Value:   "",
	}

	optionWebhookURLs = &cli.StringSliceFlag{
		Name:    "webhook-urls",
		Usage:   "List of webhook URLs to send notifications to",
		EnvVars: []string{"PRECONF_RPC_WEBHOOK_URLS"},
	}

	optionAuthToken = &cli.StringFlag{
		Name:    "auth-token",
		Usage:   "authentication token for securing endpoints",
		EnvVars: []string{"PRECONF_RPC_AUTH_TOKEN"},
		Value:   "",
	}

	optionSimulationURLs = &cli.StringSliceFlag{
		Name:     "simulation-url",
		Usage:    "URL(s) for the transaction simulation service. Multiple URLs can be specified for fallback support (first URL is primary, others are fallbacks)",
		EnvVars:  []string{"PRECONF_RPC_SIMULATION_URL"},
		Required: true,
	}

	optionUseInlineSimulation = &cli.BoolFlag{
		Name:    "use-inline-simulation",
		Usage:   "Use inline simulation via debug_traceCall instead of external rethsim service. When false (default), uses external rethsim. When true, uses debug_traceCall (requires RPC with debug API support like Alchemy, Infura, or Erigon)",
		EnvVars: []string{"PRECONF_RPC_USE_INLINE_SIMULATION"},
		Value:   false,
	}

	optionBackrunnerAPIURL = &cli.StringFlag{
		Name:     "backrunner-api-url",
		Usage:    "URL for the transaction backrun service",
		EnvVars:  []string{"PRECONF_RPC_BACKRUNNER_API_URL"},
		Required: true,
	}

	optionBackrunnerRPCURL = &cli.StringFlag{
		Name:     "backrunner-rpc-url",
		Usage:    "URL for the backrun RPC",
		EnvVars:  []string{"PRECONF_RPC_BACKRUNNER_RPC_URL"},
		Required: true,
	}

	optionBackrunnerAPIKey = &cli.StringFlag{
		Name:     "backrunner-api-key",
		Usage:    "API key for the backrun service",
		EnvVars:  []string{"PRECONF_RPC_BACKRUNNER_API_KEY"},
		Required: true,
	}

	optionPointsAPIURL = &cli.StringFlag{
		Name:    "points-api-url",
		Usage:   "URL for the points tracking service",
		EnvVars: []string{"PRECONF_RPC_POINTS_API_URL"},
	}

	optionPointsAPIKey = &cli.StringFlag{
		Name:    "points-api-key",
		Usage:   "API key for the points tracking service",
		EnvVars: []string{"PRECONF_RPC_POINTS_API_KEY"},
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"PRECONF_RPC_LOG_FMT"},
		Value:   "text",
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
		EnvVars: []string{"PRECONF_RPC_LOG_LEVEL"},
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
		EnvVars: []string{"PRECONF_RPC_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	}

	optionExplorerEndpoint = &cli.StringFlag{
		Name:    "explorer-endpoint",
		Usage:   "Explorer API endpoint for submitting transactions",
		EnvVars: []string{"PRECONF_RPC_EXPLORER_API_ENDPOINT"},
	}

	optionExplorerApiKey = &cli.StringFlag{
		Name:    "explorer-apikey",
		Usage:   "Explorer API Key",
		EnvVars: []string{"PRECONF_RPC_EXPLORER_API_KEY"},
	}

	optionExplorerAppCode = &cli.StringFlag{
		Name:    "explorer-appcode",
		Usage:   "Explorer App Code",
		EnvVars: []string{"PRECONF_RPC_EXPLORER_APPCODE"},
	}

	// FastSwap configuration
	optionBarterAPIURL = &cli.StringFlag{
		Name:    "barter-api-url",
		Usage:   "Barter API URL for swap routing",
		EnvVars: []string{"PRECONF_RPC_BARTER_API_URL"},
	}

	optionBarterAPIKey = &cli.StringFlag{
		Name:    "barter-api-key",
		Usage:   "Barter API key",
		EnvVars: []string{"PRECONF_RPC_BARTER_API_KEY"},
	}

	optionFastSettlementAddress = &cli.StringFlag{
		Name:    "fast-settlement-address",
		Usage:   "FastSettlementV3 contract address",
		EnvVars: []string{"PRECONF_RPC_FAST_SETTLEMENT_ADDRESS"},
		Action: func(ctx *cli.Context, s string) error {
			if s != "" && !common.IsHexAddress(s) {
				return fmt.Errorf("invalid fast-settlement-address: %s", s)
			}
			return nil
		},
	}

	optionFastSwapKeystorePath = &cli.StringFlag{
		Name:    "fastswap-keystore-path",
		Usage:   "Path to the FastSwap executor keystore file (separate wallet for FastSwap transactions)",
		EnvVars: []string{"PRECONF_RPC_FASTSWAP_KEYSTORE_PATH"},
	}

	optionFastSwapKeystorePassword = &cli.StringFlag{
		Name:    "fastswap-keystore-password",
		Usage:   "Password for the FastSwap executor keystore",
		EnvVars: []string{"PRECONF_RPC_FASTSWAP_KEYSTORE_PASSWORD"},
	}
)

func main() {
	app := &cli.App{
		Name:  "preconf-rpc",
		Usage: "Preconf RPC service",
		Flags: []cli.Flag{
			optionHTTPPort,
			optionPgHost,
			optionPgPort,
			optionPgUser,
			optionPgPassword,
			optionPgDbname,
			optionPgSSL,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
			optionKeystorePath,
			optionKeystorePassword,
			optionL1RPCHTTPUrl,
			optionL1RPCWSUrl,
			optionSettlementRPCUrl,
			optionBidderRPCUrl,
			optionL1ContractAddr,
			optionSettlementThreshold,
			optionSettlementTopup,
			optionGasTipCap,
			optionGasFeeCap,
			optionSettlementContractAddr,
			optionTargetDepositAmount,
			optionDepositAddress,
			optionBridgeAddress,
			optionBlocknativeAPIKey,
			optionWebhookURLs,
			optionBidderThreshold,
			optionBidderTopup,
			optionAuthToken,
			optionSimulationURLs,
			optionUseInlineSimulation,
			optionBackrunnerAPIURL,
			optionBackrunnerRPCURL,
			optionBackrunnerAPIKey,
			optionExplorerEndpoint,
			optionExplorerApiKey,
			optionExplorerAppCode,
			optionPointsAPIURL,
			optionPointsAPIKey,
			optionL1ReceiptsRPCUrl,
			optionBarterAPIURL,
			optionBarterAPIKey,
			optionFastSettlementAddress,
			optionFastSwapKeystorePath,
			optionFastSwapKeystorePassword,
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			gasTipCap, ok := new(big.Int).SetString(c.String(optionGasTipCap.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse gas-tip-cap")
			}

			gasFeeCap, ok := new(big.Int).SetString(c.String(optionGasFeeCap.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse gas-fee-cap")
			}

			targetDepositAmount, ok := new(big.Int).SetString(c.String(optionTargetDepositAmount.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse target-deposit-amount")
			}

			settlementThreshold, ok := new(big.Int).SetString(c.String(optionSettlementThreshold.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse settlement-threshold")
			}

			settlementTopup, ok := new(big.Int).SetString(c.String(optionSettlementTopup.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse settlement-topup")
			}

			bidderThreshold, ok := new(big.Int).SetString(c.String(optionBidderThreshold.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse bidder-threshold")
			}

			bidderTopup, ok := new(big.Int).SetString(c.String(optionBidderTopup.Name), 10)
			if !ok {
				return fmt.Errorf("failed to parse bidder-topup")
			}

			signer, err := keysigner.NewKeystoreSigner(
				c.String(optionKeystorePath.Name),
				c.String(optionKeystorePassword.Name),
			)
			if err != nil {
				return fmt.Errorf("failed to create signer: %w", err)
			}

			// Load separate FastSwap signer if configured
			var fastSwapSigner keysigner.KeySigner
			if c.String(optionFastSwapKeystorePath.Name) != "" {
				fastSwapSigner, err = keysigner.NewKeystoreSigner(
					c.String(optionFastSwapKeystorePath.Name),
					c.String(optionFastSwapKeystorePassword.Name),
				)
				if err != nil {
					return fmt.Errorf("failed to create FastSwap signer: %w", err)
				}
				logger.Info("FastSwap executor wallet loaded", "address", fastSwapSigner.GetAddress().Hex())
			}

			l1ReceiptsURL := c.String(optionL1ReceiptsRPCUrl.Name)
			if l1ReceiptsURL == "" {
				l1ReceiptsURL = c.String(optionL1RPCHTTPUrl.Name)
			}

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

			config := service.Config{
				HTTPPort:               c.Int(optionHTTPPort.Name),
				PgHost:                 c.String(optionPgHost.Name),
				PgPort:                 c.Int(optionPgPort.Name),
				PgUser:                 c.String(optionPgUser.Name),
				PgPassword:             c.String(optionPgPassword.Name),
				PgDbname:               c.String(optionPgDbname.Name),
				PgSSL:                  c.Bool(optionPgSSL.Name),
				Logger:                 logger,
				GasTipCap:              gasTipCap,
				GasFeeCap:              gasFeeCap,
				TargetDepositAmount:    targetDepositAmount,
				SettlementThreshold:    settlementThreshold,
				SettlementTopup:        settlementTopup,
				BidderThreshold:        bidderThreshold,
				BidderTopup:            bidderTopup,
				SettlementRPCUrl:       c.String(optionSettlementRPCUrl.Name),
				BidderRPC:              c.String(optionBidderRPCUrl.Name),
				L1RPCHTTPUrl:           c.String(optionL1RPCHTTPUrl.Name),
				L1RPCWSUrl:             c.String(optionL1RPCWSUrl.Name),
				L1ContractAddr:         common.HexToAddress(c.String(optionL1ContractAddr.Name)),
				SettlementContractAddr: common.HexToAddress(c.String(optionSettlementContractAddr.Name)),
				Signer:                 signer,
				DepositAddress:         common.HexToAddress(c.String(optionDepositAddress.Name)),
				BridgeAddress:          common.HexToAddress(c.String(optionBridgeAddress.Name)),
				PricerAPIKey:           c.String(optionBlocknativeAPIKey.Name),
				Webhooks:               c.StringSlice(optionWebhookURLs.Name),
				Token:                  c.String(optionAuthToken.Name),
				SimulatorURLs:          c.StringSlice(optionSimulationURLs.Name),
				UseInlineSimulation:    c.Bool(optionUseInlineSimulation.Name),
				BackrunnerAPIURL:       c.String(optionBackrunnerAPIURL.Name),
				BackrunnerRPC:          c.String(optionBackrunnerRPCURL.Name),
				BackrunnerAPIKey:       c.String(optionBackrunnerAPIKey.Name),
				ExplorerEndpoint:       c.String(optionExplorerEndpoint.Name),
				ExplorerApiKey:         c.String(optionExplorerApiKey.Name),
				ExplorerAppCode:        c.String(optionExplorerAppCode.Name),
				PointsAPIURL:           c.String(optionPointsAPIURL.Name),
				PointsAPIKey:           c.String(optionPointsAPIKey.Name),
				L1ReceiptsRPCUrl:       l1ReceiptsURL,
				BarterAPIURL:           c.String(optionBarterAPIURL.Name),
				BarterAPIKey:           c.String(optionBarterAPIKey.Name),
				FastSettlementAddress:  common.HexToAddress(c.String(optionFastSettlementAddress.Name)),
				FastSwapSigner:         fastSwapSigner,
			}

			s, err := service.New(&config)
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
