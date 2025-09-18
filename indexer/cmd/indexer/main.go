package main

import (
	"context"

	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/primev/mev-commit/indexer/pkg/backfill"
	"github.com/primev/mev-commit/indexer/pkg/beacon"
	"github.com/primev/mev-commit/indexer/pkg/config"
	"github.com/primev/mev-commit/indexer/pkg/database"
	"github.com/primev/mev-commit/indexer/pkg/ethereum"
	httputil "github.com/primev/mev-commit/indexer/pkg/http"
	"github.com/primev/mev-commit/indexer/pkg/relay"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg := config.LoadConfig()

	// Validate required configuration
	if cfg.EtherscanKey == "" {
		log.Fatal("[CONFIG] ETHERSCAN_API_KEY is required")
	}
	if cfg.InfuraRPC == "" {
		log.Fatal("[CONFIG] INFURA_RPC is required")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("[CONFIG] set DATABASE_URL, e.g. user:pass@tcp(host:port)/database")
	}

	log.Printf("[INIT] Starting blockchain indexer with StarRocks database")
	log.Printf("[CONFIG] Block interval: %s, Validator delay: %s, Backfill every: %s",
		cfg.BlockTick, cfg.ValidatorWait, cfg.BackfillEvery)

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
	httpc := httputil.NewHTTPClient(cfg.HTTPTimeout)
	log.Printf("[HTTP] Client initialized with %s timeout", cfg.HTTPTimeout)

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

		latestBlock, err := ethereum.GetLatestBlockNumber(httpc, cfg.InfuraRPC)
		if err != nil {
			log.Fatalf("[INIT] Failed to get latest block from RPC: %v", err)
		}

		lastBN = latestBlock - 10 // Start 10 blocks behind to ensure data availability
		log.Printf("[INIT] Starting from block %d (latest: %d)", lastBN, latestBlock)
	}

	log.Printf(" [INIT] Starting from block number: %d", lastBN)
	log.Printf("[INIT] Indexer configuration - Lookback: %d slots, Batch size: %d",
		cfg.BackfillLookback, cfg.BackfillBatch)

	// Setup tickers
	backfillTicker := time.NewTicker(cfg.BackfillEvery)
	defer backfillTicker.Stop()

	mainTicker := time.NewTicker(cfg.BlockTick)
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
			return

		case <-backfillTicker.C:
			log.Printf("[BACKFILL] Starting backfill operations...")
			backfill.RunAll(ctx, db, httpc, cfg, relays)

		case <-mainTicker.C:
			nextBN := lastBN + 1

			// Fetch execution block data
			ei, err := beacon.FetchCombinedBlockData(httpc, cfg.InfuraRPC, cfg.BeaconBase, nextBN)
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
					time.Sleep(cfg.ValidatorWait)

					vpub, err := beacon.FetchValidatorPubkey(httpc, cfg.BeaconBase, proposerIdx)
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
					time.Sleep(cfg.ValidatorWait + 500*time.Millisecond)

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

					opted, err := ethereum.CallAreOptedInAtBlock(httpc, cfg, blockNumber, vpk)
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
