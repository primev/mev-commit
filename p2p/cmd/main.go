package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	contracts "github.com/primev/mev-commit/contracts-abi/config"
	mevcommit "github.com/primev/mev-commit/p2p"
	"github.com/primev/mev-commit/p2p/pkg/node"
	ks "github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/primev/mev-commit/x/util/otelutil"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"go.opentelemetry.io/otel"
)

const (
	defaultP2PAddr = "0.0.0.0"
	defaultP2PPort = 13522

	defaultHTTPPort = 13523
	defaultRPCPort  = 13524

	defaultConfigDir = "~/.mev-commit"
	defaultKeyFile   = "key"
	defaultSecret    = "secret"
	defaultKeystore  = "keystore"
	defaultDataDir   = "db"

	defaultOracleWindowOffset = 1
)

const (
	categoryGlobal    = "Global options"
	categoryP2P       = "P2P options"
	categoryAPI       = "API options"
	categoryDebug     = "Debug options"
	categoryContracts = "Contracts"
	categoryBidder    = "Bidder options"
	categoryProvider  = "Provider options"
	categoryEthRPC    = "Ethereum RPC options"
)

var (
	portCheck = func(c *cli.Context, p int) error {
		if p < 0 || p > 65535 {
			return fmt.Errorf("invalid port number %d, expected 0 <= port <= 65535", p)
		}
		return nil
	}

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
	optionConfig = &cli.StringFlag{
		Name:     "config",
		Usage:    "Path to configuration file",
		EnvVars:  []string{"MEV_COMMIT_CONFIG"},
		Category: categoryGlobal,
	}

	optionDataDir = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "data-dir",
		Usage:    "Path to data directory",
		EnvVars:  []string{"MEV_COMMIT_DATA_DIR"},
		Value:    filepath.Join(defaultConfigDir, defaultDataDir),
		Category: categoryGlobal,
	})

	optionPrivKeyFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "priv-key-file",
		Usage:    "Path to private key file",
		EnvVars:  []string{"MEV_COMMIT_PRIVKEY_FILE"},
		Value:    filepath.Join(defaultConfigDir, defaultKeyFile),
		Category: categoryGlobal,
	})

	optionKeystorePassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "Password to use to access keystore",
		EnvVars:  []string{"MEV_COMMIT_KEYSTORE_PASSWORD"},
		Category: categoryGlobal,
	})

	optionKeystorePath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "keystore-path",
		Usage:    "Path to keystore location",
		EnvVars:  []string{"MEV_COMMIT_KEYSTORE_PATH"},
		Value:    filepath.Join(defaultConfigDir, defaultKeystore),
		Category: categoryGlobal,
	})

	optionPeerType = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "peer-type",
		Usage:    "Peer type to use, options are 'bidder', 'provider' or 'bootnode'",
		EnvVars:  []string{"MEV_COMMIT_PEER_TYPE"},
		Value:    "bidder",
		Action:   stringInCheck("peer-type", []string{"bidder", "provider", "bootnode"}),
		Category: categoryGlobal,
	})

	optionP2PPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "p2p-port",
		Usage:    "Port to listen for p2p connections",
		EnvVars:  []string{"MEV_COMMIT_P2P_PORT"},
		Value:    defaultP2PPort,
		Action:   portCheck,
		Category: categoryP2P,
	})

	optionP2PAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "p2p-addr",
		Usage:    "Address to bind for p2p connections",
		EnvVars:  []string{"MEV_COMMIT_P2P_ADDR"},
		Value:    defaultP2PAddr,
		Category: categoryP2P,
	})

	optionHTTPPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "http-port",
		Usage:    "Port to listen for http connections",
		EnvVars:  []string{"MEV_COMMIT_HTTP_PORT"},
		Value:    defaultHTTPPort,
		Action:   portCheck,
		Category: categoryAPI,
	})

	optionHTTPAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "http-addr",
		Usage:    "Address to bind for http connections",
		EnvVars:  []string{"MEV_COMMIT_HTTP_ADDR"},
		Value:    "",
		Category: categoryAPI,
	})

	optionRPCPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "rpc-port",
		Usage:    "Port to listen for rpc connections",
		EnvVars:  []string{"MEV_COMMIT_RPC_PORT"},
		Value:    defaultRPCPort,
		Action:   portCheck,
		Category: categoryAPI,
	})

	optionRPCAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "rpc-addr",
		Usage:    "Address to bind for RPC connections",
		EnvVars:  []string{"MEV_COMMIT_RPC_ADDR"},
		Value:    "",
		Category: categoryAPI,
	})

	optionBootnodes = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:     "bootnodes",
		Usage:    "List of bootnodes to connect to",
		EnvVars:  []string{"MEV_COMMIT_BOOTNODES"},
		Category: categoryP2P,
	})

	optionSecret = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "secret",
		Usage:    "Secret to use for authenticating node in the p2p network",
		EnvVars:  []string{"MEV_COMMIT_SECRET"},
		Value:    defaultSecret,
		Category: categoryP2P,
	})

	optionLogFmt = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "log-fmt",
		Usage:    "Log format to use, options are 'text' or 'json'",
		EnvVars:  []string{"MEV_COMMIT_LOG_FMT"},
		Value:    "text",
		Action:   stringInCheck("log-fmt", []string{"text", "json"}),
		Category: categoryDebug,
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "log-level",
		Usage:    "Log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars:  []string{"MEV_COMMIT_LOG_LEVEL"},
		Value:    "info",
		Action:   stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
		Category: categoryDebug,
	})

	optionLogTags = altsrc.NewStringFlag(&cli.StringFlag{
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

	optionBidderRegistryAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "bidder-registry-contract",
		Usage:   "Address of the bidder registry contract",
		EnvVars: []string{"MEV_COMMIT_BIDDER_REGISTRY_ADDR"},
		Value:   contracts.TestnetContracts.BidderRegistry,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid bidder registry address: %s", s)
			}
			return nil
		},
		Category: categoryContracts,
	})

	optionProviderRegistryAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "provider-registry-contract",
		Usage:   "Address of the provider registry contract",
		EnvVars: []string{"MEV_COMMIT_PROVIDER_REGISTRY_ADDR"},
		Value:   contracts.TestnetContracts.ProviderRegistry,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid provider registry address: %s", s)
			}
			return nil
		},
		Category: categoryContracts,
	})

	optionPreconfStoreAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "preconf-contract",
		Usage:   "Address of the preconfirmation commitment store contract",
		EnvVars: []string{"MEV_COMMIT_PRECONF_ADDR"},
		Value:   contracts.TestnetContracts.PreconfManager,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid preconfirmation commitment store address: %s", s)
			}
			return nil
		},
		Category: categoryContracts,
	})

	optionBlockTrackerAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "block-tracker-contract",
		Usage:   "Address of the block tracker contract",
		EnvVars: []string{"MEV_COMMIT_BLOCK_TRACKER_ADDR"},
		Value:   contracts.TestnetContracts.BlockTracker,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid block tracker address: %s", s)
			}
			return nil
		},
		Category: categoryContracts,
	})

	optionValidatorRouterAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "validator-router-contract",
		Usage:   "Address of the validator router contract",
		EnvVars: []string{"MEV_COMMIT_VALIDATOR_ROUTER_ADDR"},
		Value:   contracts.HoleskyContracts.ValidatorOptInRouter,
		Action: func(ctx *cli.Context, s string) error {
			if !common.IsHexAddress(s) {
				return fmt.Errorf("invalid validator router address: %s", s)
			}
			return nil
		},
		Category: categoryContracts,
	})

	optionAutodepositAmount = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "autodeposit-amount",
		Usage:    "Amount to auto deposit in each window in wei",
		EnvVars:  []string{"MEV_COMMIT_AUTODEPOSIT_AMOUNT"},
		Category: categoryBidder,
	})

	optionAutodepositEnabled = altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:     "autodeposit-enabled",
		Usage:    "Enable auto deposit",
		EnvVars:  []string{"MEV_COMMIT_AUTODEPOSIT_ENABLED"},
		Value:    false,
		Category: categoryBidder,
	})

	optionSettlementRPCEndpoint = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "settlement-rpc-endpoint",
		Usage:    "RPC endpoint of the settlement layer",
		EnvVars:  []string{"MEV_COMMIT_SETTLEMENT_RPC_ENDPOINT"},
		Value:    "http://localhost:8545",
		Category: categoryEthRPC,
	})

	optionSettlementWSRPCEndpoint = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "settlement-ws-rpc-endpoint",
		Usage:    "WebSocket RPC endpoint of the settlement layer",
		EnvVars:  []string{"MEV_COMMIT_SETTLEMENT_WS_RPC_ENDPOINT"},
		Category: categoryEthRPC,
	})

	optionNATAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "nat-addr",
		Usage:    "External address of the node to advertise to other peers",
		EnvVars:  []string{"MEV_COMMIT_NAT_ADDR"},
		Category: categoryP2P,
	})

	optionNATPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "nat-port",
		Usage:    "Externally mapped port for the node",
		EnvVars:  []string{"MEV_COMMIT_NAT_PORT"},
		Value:    defaultP2PPort,
		Category: categoryP2P,
	})

	optionServerTLSCert = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "server-tls-certificate",
		Usage:    "Path to the server TLS certificate",
		EnvVars:  []string{"MEV_COMMIT_SERVER_TLS_CERTIFICATE"},
		Category: categoryAPI,
	})

	optionServerTLSPrivateKey = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "server-tls-private-key",
		Usage:    "Path to the server TLS private key",
		EnvVars:  []string{"MEV_COMMIT_SERVER_TLS_PRIVATE_KEY"},
		Category: categoryAPI,
	})

	optionProviderWhitelist = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:    "provider-whitelist",
		Usage:   "List of provider addresses to whitelist for bids",
		EnvVars: []string{"MEV_COMMIT_PROVIDER_WHITELIST"},
		Action: func(ctx *cli.Context, vals []string) error {
			for i, v := range vals {
				if !common.IsHexAddress(v) {
					return fmt.Errorf("invalid provider address at index %d: %s", i, v)
				}
			}
			return nil
		},
		Category: categoryBidder,
	})

	optionGasLimit = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "gas-limit",
		Usage:    "Use predefined gas limit for transactions",
		EnvVars:  []string{"MEV_COMMIT_GAS_LIMIT"},
		Value:    2000000,
		Category: categoryGlobal,
	})

	optionGasTipCap = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "gas-tip-cap",
		Usage:    "Use predefined gas tip cap for transactions",
		EnvVars:  []string{"MEV_COMMIT_GAS_TIP_CAP"},
		Value:    "1000000000", // 1 gWEI
		Category: categoryGlobal,
	})

	optionGasFeeCap = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "gas-fee-cap",
		Usage:    "Use predefined gas fee cap for transactions",
		EnvVars:  []string{"MEV_COMMIT_GAS_FEE_CAP"},
		Value:    "2000000000", // 2 gWEI
		Category: categoryGlobal,
	})

	optionBeaconAPIURL = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "beacon-api-url",
		Usage:    "URL of the beacon chain API",
		Value:    "https://ethereum-holesky-beacon-api.publicnode.com",
		EnvVars:  []string{"MEV_COMMIT_BEACON_API_URL"},
		Category: categoryEthRPC,
	})

	optionL1RPCURL = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "l1-rpc-url",
		Usage:    "URL for L1 RPC",
		Value:    "https://ethereum-holesky-rpc.publicnode.com",
		EnvVars:  []string{"MEV_COMMIT_L1_RPC_URL"},
		Category: categoryEthRPC,
	})

	optionOTelCollectorEndpointURL = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "otel-collector-endpoint-url",
		Usage:   "URL for OpenTelemetry collector endpoint",
		EnvVars: []string{"MEV_COMMIT_OTEL_COLLECTOR_ENDPOINT_URL"},
		Action: func(_ *cli.Context, s string) error {
			_, err := url.Parse(s)
			return err
		},
		Category: categoryDebug,
	})

	optionBidderBidTimeout = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:     "bidder-bid-timeout",
		Usage:    "Timeout to use on bidder node to wait for preconfirmations after submitting bids",
		EnvVars:  []string{"MEV_COMMIT_BIDDER_BID_TIMEOUT"},
		Value:    30 * time.Second,
		Category: categoryBidder,
	})

	optionProviderDecisionTimeout = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:     "provider-decision-timeout",
		Usage:    "Timeout to use on provider node to allow for preconfirmation decision logic after receiving bids",
		EnvVars:  []string{"MEV_COMMIT_PROVIDER_DECISION_TIMEOUT"},
		Value:    30 * time.Second,
		Category: categoryProvider,
	})

	optionLaggardMode = altsrc.NewIntFlag(&cli.IntFlag{
		Name:     "laggard-mode",
		Usage:    "No of blocks to lag behind for L1 chain when fetching validator duties",
		EnvVars:  []string{"MEV_COMMIT_LAGGARD_MODE"},
		Value:    10,
		Category: categoryEthRPC,
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionDataDir,
		optionPeerType,
		optionPrivKeyFile,
		optionKeystorePassword,
		optionKeystorePath,
		optionP2PPort,
		optionP2PAddr,
		optionHTTPPort,
		optionHTTPAddr,
		optionRPCPort,
		optionRPCAddr,
		optionBootnodes,
		optionSecret,
		optionLogFmt,
		optionLogLevel,
		optionLogTags,
		optionBidderRegistryAddr,
		optionProviderRegistryAddr,
		optionPreconfStoreAddr,
		optionBlockTrackerAddr,
		optionValidatorRouterAddr,
		optionAutodepositAmount,
		optionAutodepositEnabled,
		optionSettlementRPCEndpoint,
		optionSettlementWSRPCEndpoint,
		optionNATAddr,
		optionNATPort,
		optionServerTLSCert,
		optionServerTLSPrivateKey,
		optionProviderWhitelist,
		optionGasLimit,
		optionGasTipCap,
		optionGasFeeCap,
		optionBeaconAPIURL,
		optionL1RPCURL,
		optionOTelCollectorEndpointURL,
		optionBidderBidTimeout,
		optionProviderDecisionTimeout,
	}

	app := &cli.App{
		Name:    "mev-commit",
		Usage:   "Start mev-commit node",
		Version: mevcommit.Version(),
		Flags:   flags,
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
		Action:  initializeApplication,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		fmt.Fprintln(app.Writer, "received interrupt signal, exiting... Force exit with Ctrl+C")
		cancel()
		<-sigc
		fmt.Fprintln(app.Writer, "force exiting...")
		os.Exit(1)
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Fprintln(app.Writer, "exited with error:", err)
	}

	os.Exit(0)
}

func initializeApplication(c *cli.Context) error {
	if err := verifyKeystorePasswordPresence(c); err != nil {
		return err
	}
	if err := launchNodeWithConfig(c); err != nil {
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

// launchNodeWithConfig configures and starts the p2p node based on the CLI context.
func launchNodeWithConfig(c *cli.Context) (err error) {
	logger, err := util.NewLogger(
		c.String(optionLogLevel.Name),
		c.String(optionLogFmt.Name),
		c.String(optionLogTags.Name),
		c.App.Writer,
	)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	keysigner, err := newKeySigner(c)
	if err != nil {
		return fmt.Errorf("failed to create keysigner: %w", err)
	}
	logger.Info("key signer account", "address", keysigner.GetAddress().Hex(), "url", keysigner.String())

	httpAddr := fmt.Sprintf("%s:%d", c.String(optionHTTPAddr.Name), c.Int(optionHTTPPort.Name))
	rpcAddr := fmt.Sprintf("%s:%d", c.String(optionRPCAddr.Name), c.Int(optionRPCPort.Name))
	natAddr := ""
	if c.String(optionNATAddr.Name) != "" {
		natAddr = fmt.Sprintf("%s:%d", c.String(optionNATAddr.Name), c.Int(optionNATPort.Name))
	}

	var (
		autodepositAmount *big.Int
		ok                bool
	)
	if c.String(optionAutodepositAmount.Name) != "" && c.Bool(optionAutodepositEnabled.Name) {
		autodepositAmount, ok = new(big.Int).SetString(c.String(optionAutodepositAmount.Name), 10)
		if !ok {
			return fmt.Errorf("failed to parse autodeposit amount %q", c.String(optionAutodepositAmount.Name))
		}
	}
	crtFile := c.String(optionServerTLSCert.Name)
	keyFile := c.String(optionServerTLSPrivateKey.Name)
	if (crtFile == "") != (keyFile == "") {
		return fmt.Errorf("both -%s and -%s must be provided to enable TLS", optionServerTLSCert.Name, optionServerTLSPrivateKey.Name)
	}

	whitelist := make([]common.Address, 0, len(c.StringSlice(optionProviderWhitelist.Name)))
	for _, addr := range c.StringSlice(optionProviderWhitelist.Name) {
		whitelist = append(whitelist, common.HexToAddress(addr))
	}

	var gasTipCap, gasFeeCap *big.Int
	if c.String(optionGasTipCap.Name) != "" {
		gasTipCap, ok = new(big.Int).SetString(c.String(optionGasTipCap.Name), 10)
		if !ok {
			return fmt.Errorf("failed to parse gas tip cap %q", c.String(optionGasTipCap.Name))
		}
	}
	if c.String(optionGasFeeCap.Name) != "" {
		gasFeeCap, ok = new(big.Int).SetString(c.String(optionGasFeeCap.Name), 10)
		if !ok {
			return fmt.Errorf("failed to parse gas fee cap %q", c.String(optionGasFeeCap.Name))
		}
	}

	dbPath := ""
	if c.String(optionDataDir.Name) != "" {
		dbPath, err = util.ResolveFilePath(c.String(optionDataDir.Name))
		if err != nil {
			return fmt.Errorf("failed to resolve data directory: %w", err)
		}
		if err := os.MkdirAll(dbPath, 0700); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	if c.IsSet(optionOTelCollectorEndpointURL.Name) {
		logger.Info("setting up OpenTelemetry SDK", "endpoint", c.String(optionOTelCollectorEndpointURL.Name))
		ctx, cancel := context.WithTimeout(c.Context, 5*time.Second)
		defer cancel()
		shutdown, err := otelutil.SetupOTelSDK(
			ctx,
			c.String(optionOTelCollectorEndpointURL.Name),
			c.String(optionLogTags.Name),
		)
		if err != nil {
			logger.Warn("failed to setup OpenTelemetry SDK; continuing without telemetry", "error", err)
		} else {
			otel.SetLogger(logr.FromSlogHandler(
				logger.Handler().WithAttrs([]slog.Attr{
					{Key: "component", Value: slog.StringValue("otel")},
				}),
			))
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err = errors.Join(err, shutdown(ctx))
				cancel()
			}()
		}
	}

	nd, err := node.NewNode(&node.Options{
		KeySigner:                keysigner,
		DataDir:                  dbPath,
		Secret:                   c.String(optionSecret.Name),
		PeerType:                 c.String(optionPeerType.Name),
		P2PPort:                  c.Int(optionP2PPort.Name),
		P2PAddr:                  c.String(optionP2PAddr.Name),
		HTTPAddr:                 httpAddr,
		RPCAddr:                  rpcAddr,
		Logger:                   logger,
		Bootnodes:                c.StringSlice(optionBootnodes.Name),
		PreconfContract:          c.String(optionPreconfStoreAddr.Name),
		ProviderRegistryContract: c.String(optionProviderRegistryAddr.Name),
		BidderRegistryContract:   c.String(optionBidderRegistryAddr.Name),
		BlockTrackerContract:     c.String(optionBlockTrackerAddr.Name),
		ValidatorRouterContract:  c.String(optionValidatorRouterAddr.Name),
		AutodepositAmount:        autodepositAmount,
		RPCEndpoint:              c.String(optionSettlementRPCEndpoint.Name),
		WSRPCEndpoint:            c.String(optionSettlementWSRPCEndpoint.Name),
		NatAddr:                  natAddr,
		TLSCertificateFile:       crtFile,
		TLSPrivateKeyFile:        keyFile,
		ProviderWhitelist:        whitelist,
		DefaultGasLimit:          uint64(c.Int(optionGasLimit.Name)),
		DefaultGasTipCap:         gasTipCap,
		DefaultGasFeeCap:         gasFeeCap,
		OracleWindowOffset:       big.NewInt(defaultOracleWindowOffset),
		BeaconAPIURL:             c.String(optionBeaconAPIURL.Name),
		L1RPCURL:                 c.String(optionL1RPCURL.Name),
		LaggardMode:              c.Int(optionLaggardMode.Name),
		BidderBidTimeout:         c.Duration(optionBidderBidTimeout.Name),
		ProviderDecisionTimeout:  c.Duration(optionProviderDecisionTimeout.Name),
	})
	if err != nil {
		return fmt.Errorf("failed starting node: %w", err)
	}

	<-c.Done()
	logger.Info("shutting down...")
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
		logger.Error("failed to close node in time")
	}

	return nil
}

func newKeySigner(c *cli.Context) (ks.KeySigner, error) {
	if c.IsSet(optionKeystorePath.Name) {
		return ks.NewKeystoreSigner(c.String(optionKeystorePath.Name), c.String(optionKeystorePassword.Name))
	}
	return ks.NewPrivateKeySigner(c.String(optionPrivKeyFile.Name))
}
