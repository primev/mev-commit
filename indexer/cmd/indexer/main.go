package main

import (
	"context"

	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"
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
	// if cfg.EtherscanKey == "" {
	// 	log.Fatal("[CONFIG] ETHERSCAN_API_KEY is required")
	// }
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
	db, err := database.MustConnect(ctx, dbURL, 50, 10)
	if err != nil {
		log.Fatalf("[DB] Connection failed: %v", err)
	}
	defer db.Close()
	log.Printf("[DB] Connected to StarRocks database")

	// Ensure required tables exist
	if err := db.EnsureStateTable(ctx); err != nil {
		log.Fatalf("[DB] Failed to ensure state table: %v", err)
	}
	log.Printf("[DB] State table ready")

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
		log.Printf("[INIT] No previous state found, checking database for latest block...")
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

	log.Printf("[INIT] Starting from block number: %d", lastBN)
	log.Printf("[INIT] Indexer configuration - Lookback: %d slots, Batch size: %d",
		cfg.BackfillLookback, cfg.BackfillBatch)

	// Setup tickers
	backfillTicker := time.NewTicker(cfg.BackfillEvery)
	defer backfillTicker.Stop()

	mainTicker := time.NewTicker(cfg.BlockTick)
	defer mainTicker.Stop()

	log.Printf("🎉[INIT] Blockchain indexer started successfully")

	// Main processing loop
	for {
		select {
		case <-ctx.Done():
			log.Printf("[SHUTDOWN] Graceful shutdown initiated: %v", ctx.Err())
			if err := db.SaveLastBlockNumber(ctx, lastBN); err != nil {
				log.Printf("[SHUTDOWN] Failed to save last block number: %v", err)
			}
			log.Printf("[SHUTDOWN] Indexer stopped at block %d", lastBN)
			return

		case <-backfillTicker.C:
			log.Printf("[BACKFILL] Starting backfill operations...")
			backfill.RunAll(ctx, db, httpc, cfg, relays)
			go func() {
				epochCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				var epochResp struct {
					Data struct {
						Epoch int64 `json:"epoch"`
					} `json:"data"`
				}

				epochURL := fmt.Sprintf("%s/epoch/latest", cfg.BeaconBase)
				if err := httputil.FetchJSONWithRetry(epochCtx, httpc, epochURL, &epochResp, 3, 300*time.Millisecond); err != nil {
					log.Printf("[EPOCH] Failed to fetch latest epoch: %v", err)
					return
				}

				// Store/update latest epoch in a simple way
				_, err := db.Conn.ExecContext(epochCtx, `
            INSERT INTO state (key_name, value) VALUES ('latest_epoch', ?) 
            ON DUPLICATE KEY UPDATE value = ?`,
					fmt.Sprintf("%d", epochResp.Data.Epoch),
					fmt.Sprintf("%d", epochResp.Data.Epoch))

				if err != nil {
					log.Printf("[EPOCH] Failed to save latest epoch: %v", err)
				} else {
					log.Printf("[EPOCH] Latest epoch: %d", epochResp.Data.Epoch)
				}
			}()
		case <-mainTicker.C:
			nextBN := lastBN + 1

			// Fetch execution block data
			ei, err := beacon.FetchCombinedBlockData(httpc, cfg.InfuraRPC, cfg.BeaconBase, nextBN)
			if err != nil || ei == nil {
				log.Printf("[BLOCK] Block %d not available yet: %v", nextBN, err)
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
				log.Printf("[BUILDER] Builder pubkey: %s...", (*ei.BuilderHex)[:20])
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
			go func(slot int64) {
				time.Sleep(1 * time.Second) // Brief delay for block to settle

				// Fetch BeaconChain.in block/{slot} data for missing fields
				url := fmt.Sprintf("%s/block/%d", cfg.BeaconBase, slot)
				ctx2, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()

				var resp struct {
					Data struct {
						Status     string `json:"status"` // 0=missed, 1=proposed, 2=orphaned
						Graffiti   string `json:"graffiti"`
						BlockRoot  string `json:"blockroot"`
						ParentRoot string `json:"parentroot"`
						StateRoot  string `json:"stateroot"`
					} `json:"data"`
				}

				if err := httputil.FetchJSONWithRetry(ctx2, httpc, url, &resp, 3, 300*time.Millisecond); err != nil {
					log.Printf("[BEACON] Failed to fetch slot metadata for %d: %v", slot, err)
					return
				}
				// Convert string status to int
				var blockStatus int = 1 // default to proposed
				switch resp.Data.Status {
				case "0":
					blockStatus = 0
				case "1":
					blockStatus = 1
				case "2":
					blockStatus = 2
				default:
					blockStatus = 1
				}
				// Update blocks table with beacon metadata
				updateCtx, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel2()

				// _, err := db.Conn.ExecContext(updateCtx, `
				// INSERT INTO blocks (
				//     slot, block_status, graffiti, block_root, parent_root, state_root
				// ) VALUES (?, ?, ?, ?, ?, ?)`,
				// slot, blockStatus, resp.Data.Graffiti, resp.Data.BlockRoot,
				// resp.Data.ParentRoot, resp.Data.StateRoot)
				var exists bool
				err := db.Conn.QueryRowContext(updateCtx,
					`SELECT 1 FROM blocks WHERE slot = ? LIMIT 1`, slot).Scan(&exists)

				if err == nil && exists {
					// Row exists, we can do the beacon metadata insert
					_, err = db.Conn.ExecContext(updateCtx, `
        INSERT INTO blocks (slot, block_status, graffiti, block_root, parent_root, state_root) 
        VALUES (?, ?, ?, ?, ?, ?)`,
						slot, blockStatus, resp.Data.Graffiti, resp.Data.BlockRoot,
						resp.Data.ParentRoot, resp.Data.StateRoot)
				} else {
					log.Printf("[BEACON] Row for slot %d doesn't exist yet, skipping metadata update", slot)
				}

				if err != nil {
					log.Printf("[BEACON] Failed to update metadata for slot %d: %v", slot, err)
				} else {
					statusText := "proposed"
					if resp.Data.Status == string(0) {
						statusText = "missed"
					}
					if resp.Data.Status == string(2) {
						statusText = "orphaned"
					}
					log.Printf("[BEACON] Updated slot %d metadata (%s)", slot, statusText)
				}
			}(ei.Slot)
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
					log.Printf("[BIDS] Relay %d (%s) failed: %v", rr.ID, rr.URL, err)
					continue
				}

				ok, failedBids, lastError := relay.BatchInsertBids(ctx, db, ei.Slot, rr.ID, bids)
				relayBids := ok

				if mainContextCanceled {
					break
				}

				// Log summary for this relay - only one log per relay
				if failedBids > 0 {
					log.Printf("[BIDS] Relay %d: %d bids collected, %d failed (last error: %v)",
						rr.ID, relayBids, failedBids, lastError)
				} else if relayBids > 0 {
					log.Printf("[BIDS] Relay %d: %d bids collected", rr.ID, relayBids)
				}

				if relayBids > 0 {
					totalBids += relayBids
					successfulRelays++
				}
			}

			log.Printf("[SUMMARY] Block %d: %d bids from %d relays", nextBN, totalBids, successfulRelays)
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
							log.Printf("[VALIDATOR] Failed to save pubkey for slot %d: %v", slot, err)
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
						`INSERT INTO blocks (slot, validator_opted_in) VALUES (?, ?)`, slot, opted)
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
