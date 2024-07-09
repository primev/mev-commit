package main

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	contracts "github.com/primev/mev-commit/contracts-abi/config"
	mevcommit "github.com/primev/mev-commit/p2p"
	"github.com/primev/mev-commit/p2p/pkg/node"
	ks "github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
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
		Name:    "config",
		Usage:   "path to config file",
		EnvVars: []string{"MEV_COMMIT_CONFIG"},
	}

	optionPrivKeyFile = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "priv-key-file",
		Usage:   "path to private key file",
		EnvVars: []string{"MEV_COMMIT_PRIVKEY_FILE"},
		Value:   filepath.Join(defaultConfigDir, defaultKeyFile),
	})

	optionKeystorePassword = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-password",
		Usage:   "use to access keystore",
		EnvVars: []string{"MEV_COMMIT_KEYSTORE_PASSWORD"},
	})

	optionKeystorePath = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "keystore-path",
		Usage:   "path to keystore location",
		EnvVars: []string{"MEV_COMMIT_KEYSTORE_PATH"},
		Value:   filepath.Join(defaultConfigDir, defaultKeystore),
	})

	optionPeerType = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "peer-type",
		Usage:   "peer type to use, options are 'bidder', 'provider' or 'bootnode'",
		EnvVars: []string{"MEV_COMMIT_PEER_TYPE"},
		Value:   "bidder",
		Action:  stringInCheck("peer-type", []string{"bidder", "provider", "bootnode"}),
	})

	optionP2PPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "p2p-port",
		Usage:   "port to listen for p2p connections",
		EnvVars: []string{"MEV_COMMIT_P2P_PORT"},
		Value:   defaultP2PPort,
		Action:  portCheck,
	})

	optionP2PAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "p2p-addr",
		Usage:   "address to bind for p2p connections",
		EnvVars: []string{"MEV_COMMIT_P2P_ADDR"},
		Value:   defaultP2PAddr,
	})

	optionHTTPPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen for http connections",
		EnvVars: []string{"MEV_COMMIT_HTTP_PORT"},
		Value:   defaultHTTPPort,
		Action:  portCheck,
	})

	optionHTTPAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "http-addr",
		Usage:   "address to bind for http connections",
		EnvVars: []string{"MEV_COMMIT_HTTP_ADDR"},
		Value:   "",
	})

	optionRPCPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "rpc-port",
		Usage:   "port to listen for rpc connections",
		EnvVars: []string{"MEV_COMMIT_RPC_PORT"},
		Value:   defaultRPCPort,
		Action:  portCheck,
	})

	optionRPCAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "rpc-addr",
		Usage:   "address to bind for RPC connections",
		EnvVars: []string{"MEV_COMMIT_RPC_ADDR"},
		Value:   "",
	})

	optionBootnodes = altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:    "bootnodes",
		Usage:   "list of bootnodes to connect to",
		EnvVars: []string{"MEV_COMMIT_BOOTNODES"},
	})

	optionSecret = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "secret",
		Usage:   "secret to use for signing",
		EnvVars: []string{"MEV_COMMIT_SECRET"},
		Value:   defaultSecret,
	})

	optionLogFmt = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MEV_COMMIT_LOG_FMT"},
		Value:   "text",
		Action:  stringInCheck("log-fmt", []string{"text", "json"}),
	})

	optionLogLevel = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"MEV_COMMIT_LOG_LEVEL"},
		Value:   "info",
		Action:  stringInCheck("log-level", []string{"debug", "info", "warn", "error"}),
	})

	optionLogTags = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"MEV_COMMIT_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	})

	optionBidderRegistryAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "bidder-registry-contract",
		Usage:   "address of the bidder registry contract",
		EnvVars: []string{"MEV_COMMIT_BIDDER_REGISTRY_ADDR"},
		Value:   contracts.TestnetContracts.BidderRegistry,
	})

	optionProviderRegistryAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "provider-registry-contract",
		Usage:   "address of the provider registry contract",
		EnvVars: []string{"MEV_COMMIT_PROVIDER_REGISTRY_ADDR"},
		Value:   contracts.TestnetContracts.ProviderRegistry,
	})

	optionPreconfStoreAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "preconf-contract",
		Usage:   "address of the preconfirmation commitment store contract",
		EnvVars: []string{"MEV_COMMIT_PRECONF_ADDR"},
		Value:   contracts.TestnetContracts.PreconfCommitmentStore,
	})

	optionBlockTrackerAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "block-tracker-contract",
		Usage:   "address of the block tracker contract",
		EnvVars: []string{"MEV_COMMIT_BLOCK_TRACKER_ADDR"},
		Value:   contracts.TestnetContracts.BlockTracker,
	})

	optionAutodepositAmount = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "autodeposit-amount",
		Usage:   "amount to auto deposit",
		EnvVars: []string{"MEV_COMMIT_AUTODEPOSIT_AMOUNT"},
	})

	optionAutodepositEnabled = altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:    "autodeposit-enabled",
		Usage:   "enable auto deposit",
		EnvVars: []string{"MEV_COMMIT_AUTODEPOSIT_ENABLED"},
		Value:   false,
	})

	optionSettlementRPCEndpoint = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-rpc-endpoint",
		Usage:   "rpc endpoint of the settlement layer",
		EnvVars: []string{"MEV_COMMIT_SETTLEMENT_RPC_ENDPOINT"},
		Value:   "http://localhost:8545",
	})

	optionSettlementWSRPCEndpoint = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "settlement-ws-rpc-endpoint",
		Usage:   "WebSocket rpc endpoint of the settlement layer",
		EnvVars: []string{"MEV_COMMIT_SETTLEMENT_WS_RPC_ENDPOINT"},
	})

	optionNATAddr = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "nat-addr",
		Usage:   "external address of the node",
		EnvVars: []string{"MEV_COMMIT_NAT_ADDR"},
	})

	optionNATPort = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "nat-port",
		Usage:   "externally mapped port for the node",
		EnvVars: []string{"MEV_COMMIT_NAT_PORT"},
		Value:   defaultP2PPort,
	})

	optionServerTLSCert = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server-tls-certificate",
		Usage:   "Path to the server TLS certificate",
		EnvVars: []string{"MEV_COMMIT_SERVER_TLS_CERTIFICATE"},
	})

	optionServerTLSPrivateKey = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "server-tls-private-key",
		Usage:   "Path to the server TLS private key",
		EnvVars: []string{"MEV_COMMIT_SERVER_TLS_PRIVATE_KEY"},
	})
)

func main() {
	flags := []cli.Flag{
		optionConfig,
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
		optionAutodepositAmount,
		optionAutodepositEnabled,
		optionSettlementRPCEndpoint,
		optionSettlementWSRPCEndpoint,
		optionNATAddr,
		optionNATPort,
		optionServerTLSCert,
		optionServerTLSPrivateKey,
	}

	app := &cli.App{
		Name:    "mev-commit",
		Usage:   "Start mev-commit node",
		Version: mevcommit.Version(),
		Flags:   flags,
		Before:  altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(optionConfig.Name)),
		Action:  initializeApplication,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(app.Writer, "exited with error:", err)
	}
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
func launchNodeWithConfig(c *cli.Context) error {
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
		ok bool
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

	nd, err := node.NewNode(&node.Options{
		KeySigner:                keysigner,
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
		AutodepositAmount:        autodepositAmount,
		RPCEndpoint:              c.String(optionSettlementRPCEndpoint.Name),
		WSRPCEndpoint:            c.String(optionSettlementWSRPCEndpoint.Name),
		NatAddr:                  natAddr,
		TLSCertificateFile:       crtFile,
		TLSPrivateKeyFile:        keyFile,
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
