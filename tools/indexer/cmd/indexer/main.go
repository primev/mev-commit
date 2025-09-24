package main

import (
	"context"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/indexer/pkg/backfill"
	"github.com/primev/mev-commit/indexer/pkg/beacon"
	"github.com/primev/mev-commit/indexer/pkg/config"
	"github.com/primev/mev-commit/indexer/pkg/database"
	"github.com/primev/mev-commit/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/indexer/pkg/http"
	"github.com/primev/mev-commit/indexer/pkg/relay"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	BlockTick        time.Duration
	ValidatorWait    time.Duration
	BackfillEvery    time.Duration
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

	optionBackfillEvery = altsrc.NewDurationFlag(&cli.DurationFlag{
		Name:    "backfill-every",
		Usage:   "interval for backfill operations",
		EnvVars: []string{"INDEXER_BACKFILL_EVERY"},
		Value:   5 * time.Minute,
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
		Value:   50,
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
		BackfillEvery:    c.Duration("backfill-every"),
		BackfillLookback: int64(c.Int("backfill-lookback-slots")),
		BackfillBatch:    c.Int("backfill-batch"),
		HTTPTimeout:      c.Duration("http-timeout"),
		OptInContract:    c.String("opt-in-contract"),
		EtherscanKey:     c.String("etherscan-api-key"),
		InfuraRPC:        c.String("infura-rpc"),
		BeaconBase:       c.String("beacon-base"),
	}
}

func startIndexer(c *cli.Context) error {

	dbURL := c.String(optionDatabaseURL.Name)
	infuraRPC := c.String(optionInfuraRPC.Name)
	beaconBase := c.String(optionBeaconBase.Name)
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration

	// Validate required configuration

	log.Printf("[INIT] Starting blockchain indexer with StarRocks database")
	log.Printf("[CONFIG] Block interval: %s, Validator delay: %s, Backfill every: %s",
		c.Duration("block-interval"), c.Duration("validator-delay"), c.Duration("backfill-every"))

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to StarRocks database
	db, err := database.MustConnect(ctx, dbURL, 20, 5)
	if err != nil {
		log.Fatalf("[DB] Connection failed: %v", err)
	}
	defer db.Close()
	log.Printf("[DB] Connected to StarRocks database")

	// Ensure required tables exist
	if err := db.EnsureStateTable(ctx); err != nil {
		log.Fatalf("[DB] Failed to ensure state table: %v", err)
	}
	log.Printf(" [DB] State table ready")

	// Initialize HTTP client
	httpc := httputil.NewHTTPClient(c.Duration("http-timeout"))
	log.Printf("[HTTP] Client initialized with %s timeout", c.Duration("http-timeout"))

	// Load relay configurations
	relays, err := relay.UpsertRelaysAndLoad(ctx, db)
	if err != nil {
		log.Fatalf("[RELAYS] Failed to load: %v", err)
	}
	log.Printf("[RELAYS] Loaded %d active relays:", len(relays))
	for _, r := range relays {
		log.Printf("  - Relay ID %d: %s", r.ID, r.URL)
	}

	// Initialize starting block number
	lastBN, found := db.LoadLastBlockNumber(ctx)
	if !found || lastBN == 0 {
		log.Printf(" [INIT] No previous state found, checking database for latest block...")
		err := db.Conn.QueryRowContext(ctx, `SELECT COALESCE(MAX(block_number),0) FROM blocks`).Scan(&lastBN)
		if err != nil {
			log.Printf("[INIT] Database query failed: %v", err)
		}
	}

	// Replace the hardcoded block search with:
	if lastBN == 0 {
		log.Printf("[INIT] Getting latest block from Ethereum RPC...")

		latestBlock, err := ethereum.GetLatestBlockNumber(httpc, infuraRPC)
		if err != nil {
			log.Fatalf("[INIT] Failed to get latest block from RPC: %v", err)
		}

		lastBN = latestBlock - 10 // Start 10 blocks behind to ensure data availability
		log.Printf("[INIT] Starting from block %d (latest: %d)", lastBN, latestBlock)
	}

	log.Printf(" [INIT] Starting from block number: %d", lastBN)
	log.Printf("[INIT] Indexer configuration - Lookback: %d slots, Batch size: %d",
		c.Int("backfill-lookback-slots"), c.Int("backfill-batch"))

	// Setup tickers
	backfillTicker := time.NewTicker(c.Duration("backfill-batch"))
	defer backfillTicker.Stop()

	mainTicker := time.NewTicker(c.Duration("block-interval"))
	defer mainTicker.Stop()

	log.Printf("🎉 [INIT] Blockchain indexer started successfully")

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			log.Printf(" [SHUTDOWN] Graceful shutdown initiated: %v", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				log.Printf("[SHUTDOWN] Failed to save last block number: %v", err)
			}
			log.Printf("[SHUTDOWN] Indexer stopped at block %d", lastBN)
			return nil

		case <-backfillTicker.C:
			log.Printf("[BACKFILL] Starting backfill operations...")
			backfill.RunAll(ctx, db, httpc, createOptionsFromCLI(c), relays)

		case <-mainTicker.C:
			nextBN := lastBN + 1

			// Fetch execution block data
			ei, err := beacon.FetchCombinedBlockData(httpc, infuraRPC, beaconBase, nextBN)
			if err != nil || ei == nil {
				log.Printf("⏳ [BLOCK] Block %d not available yet: %v", nextBN, err)
				continue
			}

			// Log block details
			log.Printf("[BLOCK] Processing block %d → slot %d", nextBN, ei.Slot)
			if ei.Timestamp != nil {
				log.Printf("[BLOCK] Timestamp: %v", ei.Timestamp.Format(time.RFC3339))
			}
			if ei.ProposerIdx != nil {
				log.Printf("[VALIDATOR] Proposer index: %d", *ei.ProposerIdx)
			}
			if ei.RelayTag != nil {
				log.Printf("[RELAY] Winning relay: %s", *ei.RelayTag)
			}
			if ei.BuilderHex != nil && len(*ei.BuilderHex) > 20 {
				log.Printf("🔨 [BUILDER] Builder pubkey: %s...", (*ei.BuilderHex)[:20])
			}
			if ei.RewardEth != nil {
				log.Printf("[REWARD] Producer reward: %.6f ETH", *ei.RewardEth)
			}

			// Save block to database
			if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
				log.Printf("[DB] Failed to save block %d: %v", nextBN, err)
				continue
			}
			log.Printf("[DB] Block %d saved successfully", nextBN)

			// Fetch and store bid data from all relays
			totalBids := 0
			successfulRelays := 0
			mainContextCanceled := false

			for _, rr := range relays {
				// Check if main context is canceled before processing each relay
				if ctx.Err() != nil {
					log.Printf("[BIDS] Main context canceled, stopping relay processing")
					mainContextCanceled = true
					break
				}

				bids, err := relay.FetchBuilderBlocksReceived(httpc, rr.URL, ei.Slot)
				if err != nil {
					log.Printf(" [BIDS] Relay %d (%s) failed: %v", rr.ID, rr.URL, err)
					continue
				}

				relayBids := 0
				// Create a separate context with timeout for bid insertions
				bidCtx, bidCancel := context.WithTimeout(context.Background(), 30*time.Second)

				for _, bid := range bids {
					// Check if main context is still valid
					if ctx.Err() != nil {
						log.Printf("[BIDS] Main context canceled, stopping bid insertion")
						mainContextCanceled = true
						break
					}

					if err := relay.InsertBid(bidCtx, db, ei.Slot, rr.ID, bid); err != nil {
						log.Printf(" [BIDS] Failed to insert bid for slot %d, relay %d: %v", ei.Slot, rr.ID, err)
					} else {
						relayBids++
					}
				}
				bidCancel()

				if mainContextCanceled {
					break
				}

				if relayBids > 0 {
					log.Printf(" [BIDS] Relay %d: %d bids collected", rr.ID, relayBids)
					totalBids += relayBids
					successfulRelays++
				}
			}

			log.Printf(" [SUMMARY] Block %d: %d bids from %d relays", nextBN, totalBids, successfulRelays)
			// Async validator pubkey fetch
			if ei.ProposerIdx != nil {
				go func(slot int64, proposerIdx int64) {
					time.Sleep(c.Duration("validator-delay"))

					vpub, err := beacon.FetchValidatorPubkey(httpc, beaconBase, proposerIdx)
					if err != nil {
						log.Printf("[VALIDATOR] Failed to fetch pubkey for proposer %d: %v", proposerIdx, err)
						return
					}

					if len(vpub) > 0 {
						if err := db.UpdateValidatorPubkey(context.Background(), slot, vpub); err != nil {
							log.Printf(" [VALIDATOR] Failed to save pubkey for slot %d: %v", slot, err)
						} else {
							log.Printf("[VALIDATOR] Pubkey saved for proposer %d (slot %d)", proposerIdx, slot)
						}
					}
				}(ei.Slot, *ei.ProposerIdx)
			}

			// Async opt-in status check
			if ei.ProposerIdx != nil {
				go func(slot int64, blockNumber int64) {
					time.Sleep(c.Duration("validator-delay") + 500*time.Millisecond)

					// Wait for validator pubkey to be available
					var vpk []byte
					retries := 3
					for i := 0; i < retries; i++ {
						err := db.Conn.QueryRowContext(context.Background(),
							`SELECT validator_pubkey FROM blocks WHERE slot=?`, slot).Scan(&vpk)
						if err == nil && len(vpk) > 0 {
							break
						}
						if i < retries-1 {
							time.Sleep(time.Second)
						}
					}

					if len(vpk) == 0 {
						log.Printf("[OPT-IN] Validator pubkey not available for slot %d", slot)
						return
					}

					opted, err := ethereum.CallAreOptedInAtBlock(httpc, createOptionsFromCLI(c), blockNumber, vpk)
					if err != nil {
						log.Printf("[OPT-IN] Failed to check opt-in status for slot %d: %v", slot, err)
						return
					}

					_, err = db.Conn.ExecContext(context.Background(),
						`UPDATE blocks SET validator_opted_in=? WHERE slot=?`, opted, slot)
					if err != nil {
						log.Printf("[OPT-IN] Failed to save opt-in status for slot %d: %v", slot, err)
					} else {
						log.Printf("[OPT-IN] Slot %d validator opted-in: %t", slot, opted)
					}
				}(ei.Slot, ei.BlockNumber)
			}

			lastBN = nextBN
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				log.Printf("[PROGRESS] Failed to save block number %d: %v", lastBN, err)
			} else {
				log.Printf("[PROGRESS] Advanced to block %d", lastBN)
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
		optionBackfillEvery,
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
