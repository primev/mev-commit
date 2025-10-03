package main

import (
	"context"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/tools/indexer/pkg/backfill"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
	"github.com/primev/mev-commit/tools/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/tools/indexer/pkg/http"
	"github.com/primev/mev-commit/tools/indexer/pkg/relay"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	BlockTick     time.Duration
	ValidatorWait time.Duration

	BackfillLookback int64
	BackfillBatch    int
	HTTPTimeout      time.Duration
	OptInContract    string
	EtherscanKey     string
	InfuraRPC        string
	BeaconBase       string
}

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
	optionEtherscanKey = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "etherscan-key",
		Usage:   "Etherscan API key",
		EnvVars: []string{"INDEXER_ETHERSCAN_KEY"},
	})
	optionInfuraRPC = altsrc.NewStringFlag(&cli.StringFlag{
		Name:     "infura-rpc",
		Usage:    "Infura RPC URL",
		EnvVars:  []string{"INDEXER_INFURA_RPC"},
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
		Value:   512,
	})

	optionBackfillBatch = altsrc.NewIntFlag(&cli.IntFlag{
		Name:    "backfill-batch",
		Usage:   "batch size for backfill operations",
		EnvVars: []string{"INDEXER_BACKFILL_BATCH"},
		Value:   5,
	})

	optionHTTPTimeout = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "http-timeout",
		Usage:   "HTTP client timeout",
		EnvVars: []string{"INDEXER_HTTP_TIMEOUT"},
		Value:   15 * time.Second,
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
		EtherscanKey:     c.String("etherscan-key"),
		InfuraRPC:        c.String("infura-rpc"),
		BeaconBase:       c.String("beacon-base"),
	}
}

func startIndexer(c *cli.Context) error {

	initLogger := slog.With("component", "init")
	dbLogger := slog.With("component", "db")
	httpLogger := slog.With("component", "http")
	relayLogger := slog.With("component", "relay")
	backfillLogger := slog.With("component", "backfill")
	blockLogger := slog.With("component", "block")
	bidsLogger := slog.With("component", "bids")
	validatorLogger := slog.With("component", "validator")
	optInLogger := slog.With("component", "opt-in")
	progressLogger := slog.With("component", "progress")
	shutdownLogger := slog.With("component", "shutdown")

	dbURL := c.String(optionDatabaseURL.Name)
	infuraRPC := c.String(optionInfuraRPC.Name)
	beaconBase := c.String(optionBeaconBase.Name)
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	initLogger.Info("starting blockchain indexer with StarRocks database")
	initLogger.Info("configuration loaded",
		"block_interval", c.Duration("block-interval"),
		"validator_delay", c.Duration("validator-delay"))

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to StarRocks database
	db, err := database.MustConnect(ctx, dbURL, 20, 5)
	if err != nil {
		dbLogger.Error("connection failed", "error", err)
	}
	defer db.Close()
	dbLogger.Info("connected to StarRocks database")

	// Ensure required tables exist
	if err := db.EnsureStateTable(ctx); err != nil {
		dbLogger.Error("failed to ensure state table", "error", err)
		return err
	}
	dbLogger.Info("state table ready")

	// Initialize HTTP client
	httpc := httputil.NewHTTPClient(c.Duration("http-timeout"))
	httpLogger.Info("client initialized", "timeout", c.Duration("http-timeout"))

	// Load relay configurations
	relays, err := relay.UpsertRelaysAndLoad(ctx, db)
	if err != nil {
		relayLogger.Error("failed to load", "error", err)
	}
	relayLogger.Info("loaded active relays", "count", len(relays))
	for _, r := range relays {
		relayLogger.Info("relay found", "id", r.ID, "url", r.URL)
	}

	// Initialize starting block number
	lastBN, found := db.LoadLastBlockNumber(ctx)
	if !found || lastBN == 0 {
		initLogger.Info("no previous state found, checking database for latest block")
		lastBN, err = db.GetMaxBlockNumber(ctx)
		if err != nil {
			initLogger.Error("database query failed", "error", err)
		}
	}

	// Replace the hardcoded block search with:
	if lastBN == 0 {
		initLogger.Info("getting latest block from Ethereum RPC...")

		latestBlock, err := ethereum.GetLatestBlockNumber(httpc.HTTPClient, infuraRPC)
		if err != nil {
			initLogger.Error("failed to get latest block from RPC", "error", err)
		}

		lastBN = latestBlock - 10 // Start 10 blocks behind to ensure data availability
		initLogger.Info("starting from block", "block", lastBN, "latest", latestBlock)
	}

	initLogger.Info("starting from block number", "block", lastBN)
	initLogger.Info("indexer configuration", "lookback", c.Int("backfill-lookback"), "batch", c.Int("backfill-batch"))

	if c.Int("backfill-lookback") > 0 {
		backfillLogger.Info("running one-time backfill",
			"lookback", c.Int("backfill-lookback"),
			"batch", c.Int("backfill-batch"))
		go backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays)
		backfillLogger.Info("completed startup backfill")
	} else {
		backfillLogger.Info("skipped", "reason", "backfill-lookback=0")
	}
	mainTicker := time.NewTicker(c.Duration("block-interval"))
	defer mainTicker.Stop()
	// initLogger.Info("blockchain indexer started successfully")
	// go backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays)

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			shutdownLogger.Info("graceful shutdown initiated", "reason", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				shutdownLogger.Error("failed to save last block number", "error", err)
			}
			shutdownLogger.Info("indexer stopped", "block", lastBN)
			return nil

		case <-mainTicker.C:
			nextBN := lastBN + 1

			// Fetch execution block data
			ei, err := beacon.FetchCombinedBlockData(httpc, infuraRPC, beaconBase, nextBN)
			if err != nil || ei == nil {
				blockLogger.Warn("block not available yet", "block", nextBN, "error", err)

				continue
			}

			// Log block details
			blockLogger.Info("processing block", "block", nextBN, "slot", ei.Slot)

			if ei.Timestamp != nil {
				blockLogger.Info("block timestamp", "block", nextBN, "timestamp", ei.Timestamp.Format(time.RFC3339))
			}
			if ei.ProposerIdx != nil {
				validatorLogger.Info("proposer index", "index", *ei.ProposerIdx)
			}
			if ei.RelayTag != nil {
				relayLogger.Info("winning relay", "tag", *ei.RelayTag)
			}
			if ei.BuilderHex != nil && len(*ei.BuilderHex) > 20 {
				blockLogger.Info("builder pubkey", "pubkey_prefix", (*ei.BuilderHex)[:20])
			}
			if ei.RewardEth != nil {
				blockLogger.Info("producer reward", "reward_eth", *ei.RewardEth)
			}

			// Save block to database
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				dbLogger.Error("failed to save block", "block", nextBN, "error", err)
				continue
			}
			dbLogger.Info("block saved successfully", "block", nextBN)

			// Fetch and store bid data from all relays
			totalBids := 0
			successfulRelays := 0
			mainContextCanceled := false
const batchSize = 500
			for _, rr := range relays {
				// Check if main context is canceled before processing each relay
				if ctx.Err() != nil {
					bidsLogger.Warn("main context canceled, stopping relay processing")
					mainContextCanceled = true
					break
				}

				bids, err := relay.FetchBuilderBlocksReceived(ctx, httpc, rr.URL, ei.Slot)
				if err != nil {
					bidsLogger.Error("relay failed", "relay_id", rr.ID, "url", rr.URL, "error", err)
					continue
				}

				relayBids := 0
				    batch := make([]database.BidRow, 0, batchSize)
				// a separate context with timeout for bid insertions
				// bidCtx, bidCancel := context.WithTimeout(ctx, 30*time.Second)

				for _, bid := range bids {
					// Check if main context is still valid
					if ctx.Err() != nil {
						bidsLogger.Warn("main context canceled, stopping bid insertion")
						mainContextCanceled = true
						break
					}

       if row, ok := relay.BuildBidInsert(ei.Slot, rr.ID, bid); ok {
            batch = append(batch, row)
			// for _, bid := range bids {
            if len(batch) >= batchSize {
				  insCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
                if err := db.InsertBidsBatch(insCtx, batch); err != nil {

                    bidsLogger.Error("batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
                } else {
                    relayBids += len(batch)
                }
                   cancel()
                batch = batch[:0]
            }
        }
    

    // final flush
    if len(batch) > 0 {
        if err := db.InsertBidsBatch(ctx, batch); err != nil {
            bidsLogger.Error("batch insert failed", "slot", ei.Slot, "relay_id", rr.ID, "count", len(batch), "error", err)
        } else {
            relayBids += len(batch)
        }
    }
}

				if mainContextCanceled {
					break
				}

				if relayBids > 0 {
					bidsLogger.Info("bids collected", "relay_id", rr.ID, "count", relayBids)
					totalBids += relayBids
					successfulRelays++
				}
			}

			bidsLogger.Info("summary", "block", nextBN, "total_bids", totalBids, "successful_relays", successfulRelays)
			// Async validator pubkey fetch
			if ei.ProposerIdx != nil {
				go func(slot int64, proposerIdx int64) {
					time.Sleep(c.Duration("validator-delay"))

					vpub, err := beacon.FetchValidatorPubkey(httpc, beaconBase, proposerIdx)
					if err != nil {
						validatorLogger.Error("failed to fetch pubkey", "proposer", proposerIdx, "error", err)
						return
					}

					if len(vpub) > 0 {
						if err := db.UpdateValidatorPubkey(ctx, slot, vpub); err != nil {
							validatorLogger.Error("failed to save pubkey", "slot", slot, "error", err)
						} else {
							validatorLogger.Info("pubkey saved", "proposer", proposerIdx, "slot", slot)
						}
					}
				}(ei.Slot, *ei.ProposerIdx)
			}

			// Async opt-in status check
			if ei.ProposerIdx != nil {
				go func(slot int64, blockNumber int64) {
					time.Sleep(c.Duration("validator-delay") + 500*time.Millisecond)

					// Wait for validator pubkey to be available
					vpk, err := db.GetValidatorPubkeyWithRetry(ctx, slot, 3, time.Second)
					if err != nil {
						optInLogger.Error("validator pubkey not available", "slot", slot, "error", err)
						return
					}

					opted, err := ethereum.CallAreOptedInAtBlock(httpc.HTTPClient, createOptionsFromCLI(c), blockNumber, vpk)
					if err != nil {
						optInLogger.Error("failed to check opt-in status", "slot", slot, "error", err)
						return
					}

					err = db.UpdateValidatorOptInStatus(ctx, slot, opted)
					if err != nil {
						optInLogger.Error("failed to save opt-in status", "slot", slot, "error", err)
					} else {
						optInLogger.Info("validator opt-in status", "slot", slot, "opted_in", opted)
					}
				}(ei.Slot, ei.BlockNumber)
			}

			lastBN = nextBN
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				progressLogger.Error("failed to save block number", "block", lastBN, "error", err)
			} else {
				progressLogger.Info("advanced to block", "block", lastBN)
			}
		}
	}
}

func main() {
	flags := []cli.Flag{
		optionConfig,
		optionDatabaseURL,
		optionInfuraRPC,
		optionBeaconBase,
		optionBlockInterval,
		optionValidatorDelay,

		optionBackfillLookback,
		optionBackfillBatch,
		optionHTTPTimeout,
		optionOptInContract,
		optionEtherscanKey,
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
