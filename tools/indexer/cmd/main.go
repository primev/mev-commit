package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var (
	optionConfig = &cli.StringFlag{
		Name:    "config",
		Usage:   "Path to config file",
		EnvVars: []string{"INDEXER_CONFIG"},
	}
	optionDatabaseURL = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "database-url",
		Usage:    "Database connection URL",
		EnvVars:  []string{"INDEXER_DATABASE_URL"},
		Required: true,
	})
	optionOptInContract = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "opt-in-contract",
		Usage:   "Opt-in contract address",
		EnvVars: []string{"INDEXER_OPT_IN_CONTRACT"},
		Value:   "0x821798d7b9d57dF7Ed7616ef9111A616aB19ed64",
	})
	optionAlchemyRPC = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "alchemy-rpc",
		Usage:    "Alchemy RPC URL",
		EnvVars:  []string{"INDEXER_ALCHEMY_RPC"},
		Required: true,
	})
	optionBeaconBase = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "beacon-base",
		Usage:   "Beacon API base URL",
		EnvVars: []string{"INDEXER_BEACON_BASE"},
		Value:   "https://beaconcha.in/api/v1",
	})
	optionBlockInterval = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "block-interval",
		Usage:   "interval between block processing",
		EnvVars: []string{"INDEXER_BLOCK_INTERVAL"},
		Value:   12 * time.Second,
	})

	optionValidatorDelay = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "validator-delay",
		Usage:   "delay before fetching validator data",
		EnvVars: []string{"INDEXER_VALIDATOR_DELAY"},
		Value:   1500 * time.Millisecond,
	})

	optionBackfillLookback = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "backfill-lookback",
		Usage:   "number of slots to look back for backfill",
		EnvVars: []string{"INDEXER_BACKFILL_LOOKBACK"},
		Value:   50400,
	})

	optionBackfillBatch = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "backfill-batch",
		Usage:   "batch size for backfill operations",
		EnvVars: []string{"INDEXER_BACKFILL_BATCH"},
		Value:   100,
	})

	optionHTTPTimeout = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "http-timeout",
		Usage:   "HTTP client timeout",
		EnvVars: []string{"INDEXER_HTTP_TIMEOUT"},
		Value:   15 * time.Second,
	})

	optionRelayFlag = altsrc.NewBoolFlag(&cli.BoolFlag{
		Name:    "relay",
		Usage:   "Whether to run in relay mode",
		EnvVars: []string{"INDEXER_RELAY"},
		Value:   false,
	})
	optionRelaysJSON = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "relays-json",
		Usage:   "JSON array overriding default relays (fields: relay_id,name,tag,url)",
		EnvVars: []string{"INDEXER_RELAYS_JSON"},
	})
	optionRatedAPIKey = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "rated-api-key",
		Usage:   "Rated Network API key",
		EnvVars: []string{"INDEXER_RATED_API_KEY"},
		// don't mark Required here to keep runtime flexible; we'll validate only when needed
	})

	optionQuickNodeBase = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "quicknode-base",
		Usage:   "QuickNode Beacon API base, e.g. https://<subdomain>.quiknode.pro/<token>",
		EnvVars: []string{"INDEXER_QUICKNODE_BASE"},
	})
)

func createOptionsFromCLI(c *cli.Context) *config.Config {
	return &config.Config{
		BlockTick:        c.Duration("block-interval"),
		ValidatorWait:    c.Duration("validator-delay"),
		BackfillLookback: int64(c.Int("backfill-lookback")),
		BackfillBatch:    c.Int("backfill-batch"),
		HTTPTimeout:      c.Duration("http-timeout"),
		OptInContract:    c.String("opt-in-contract"),
		AlchemyRPC:       c.String("alchemy-rpc"),
		BeaconBase:       c.String("beacon-base"),
		RelayData:        c.Bool("relay"),
		RelaysJSON:       c.String("relays-json"),
		RatedAPIKey:      c.String("rated-api-key"),
		QuickNodeBase:    c.String("quicknode-base"),
	}
}

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionDatabaseURL,
		optionAlchemyRPC,
		optionBeaconBase,
		optionBlockInterval,
		optionValidatorDelay,
		optionBackfillLookback,
		optionBackfillBatch,
		optionHTTPTimeout,
		optionOptInContract,
		optionRelayFlag,
		optionRelaysJSON,
		optionRatedAPIKey,
		optionQuickNodeBase,
	}

	app := &cli.App{
		Name:  "mev-indexer",
		Usage: "Builder/observer indexer",
		Commands: []*cli.Command{{
			Name:  "start",
			Usage: "Start the indexer",
			Flags: flags,
			Before: altsrc.InitInputSourceWithContext(
				flags, altsrc.NewYamlSourceFromFlagFunc("config"),
			),
			Action: func(c *cli.Context) error {
				return startIndexer(c)
			},
		}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigc
		_, _ = fmt.Fprintln(app.Writer, "received interrupt signal, exiting... Force exit with Ctrl+C")
		cancel()
		<-sigc
		_, _ = fmt.Fprintln(app.Writer, "force exiting...")
		os.Exit(1)
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}
