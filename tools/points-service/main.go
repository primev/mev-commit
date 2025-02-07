package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	avs "github.com/primev/mev-commit/contracts-abi/clients/MevCommitAVS"
	middleware "github.com/primev/mev-commit/contracts-abi/clients/MevCommitMiddleware"
	vanillaregistry "github.com/primev/mev-commit/contracts-abi/clients/VanillaRegistry"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/urfave/cli/v2"

	_ "github.com/mattn/go-sqlite3"
)

// For block-based monthly logic, assume ~216000 blocks per 30-day month.
const blocksInOneMonth = 216000

// rwLock protects access to the points database
var rwLock sync.RWMutex

// monthlyIncrements: the *cumulative* point totals at each full month.
var monthlyIncrements = []int64{
	1000,  // end of Month 1 => total 1000
	1800,  // end of Month 2 => total 1800
	2500,  // end of Month 3 => total 2500
	3500,  // end of Month 4 => total 3500
	5000,  // end of Month 5 => total 5000
	10000, // end of Month 6 => total 10000
}

var (
	optionRPCURL = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "URL of the Ethereum RPC server",
		EnvVars: []string{"POINTS_ETH_RPC_URL"},
		Value:   "https://eth-holesky.g.alchemy.com/v2/0DDo7YeieNEucZX3jieFfzmzOCGTKAgp",
	}
)

// --------------------------------------
//   1) CREATE/INIT DB
// --------------------------------------

func initDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./points.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Add "vault" column for storing vault address
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        pubkey TEXT NOT NULL,
        adder TEXT NOT NULL,
        vault TEXT,                             -- new vault column
        registry_type TEXT CHECK (registry_type IN ('vanilla', 'symbiotic', 'eigenlayer')),
        event_type TEXT,
        opted_in_block BIGINT NOT NULL,
        opted_out_block BIGINT,  -- null if still opted in
        points_accumulated BIGINT DEFAULT 0,
        UNIQUE(pubkey, adder, opted_in_block)
    );
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create events table: %w", err)
	}

	logger.Info("database setup complete", slog.String("path", "./points.db"))
	return db, nil
}

// --------------------------------------
//   2) POINTS COMPUTATION LOGIC
// --------------------------------------

// computePointsForMonths returns a final total if the validator has completed
// M full months (0-based). If partial month => no increment for that month.
func computePointsForMonths(blocksActive int64) int64 {
	// number of fully completed months
	fullMonths := blocksActive / blocksInOneMonth
	if fullMonths < 1 {
		// no full month completed => 0
		return 0
	}
	if fullMonths > int64(len(monthlyIncrements)) {
		// clamp if user is beyond the last entry
		fullMonths = int64(len(monthlyIncrements))
	}
	// monthlyIncrements is cumulative, so for fullMonths=2 => monthlyIncrements[1]
	return monthlyIncrements[fullMonths-1]
}

// updatePoints locks the DB in a transaction and recomputes points for each active row.
func updatePoints(db *sql.DB, logger *slog.Logger, currentBlock uint64) error {
	rwLock.Lock()
	defer rwLock.Unlock()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			commitErr := tx.Commit()
			if commitErr != nil {
				logger.Error("failed to commit transaction", "error", commitErr)
			}
		}
	}()

	rows, queryErr := tx.Query(`
        SELECT id, opted_in_block
        FROM events
        WHERE opted_out_block IS NULL
    `)
	if queryErr != nil {
		err = fmt.Errorf("failed to query events for points: %w", queryErr)
		return err
	}
	defer rows.Close()

	// Count rows first
	var count int
	countErr := tx.QueryRow("SELECT COUNT(*) FROM events WHERE opted_out_block IS NULL").Scan(&count)
	if countErr != nil {
		logger.Error("failed to count active validators", "error", countErr)
	}
	logger.Info("updating points",
		"current_block", currentBlock,
		"active_validators", count)

	for rows.Next() {
		var id int64
		var inBlock uint64
		if scanErr := rows.Scan(&id, &inBlock); scanErr != nil {
			logger.Error("scan error", "error", scanErr)
			continue
		}

		var blocksActive int64
		if currentBlock > inBlock {
			blocksActive = int64(currentBlock - inBlock)
		}
		totalPoints := computePointsForMonths(blocksActive)

		_, updErr := tx.Exec(`
            UPDATE events
            SET points_accumulated = ?
            WHERE id = ?
        `, totalPoints, id)
		if updErr != nil {
			logger.Error("failed to update points", "error", updErr, "id", id)
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		err = fmt.Errorf("rows iteration error: %w", rowsErr)
	}

	return err
}

// StartPointsRoutine spawns a background goroutine that runs every 'interval'
// and calls updatePoints with the latest block.
func StartPointsRoutine(db *sql.DB, logger *slog.Logger, interval time.Duration, ethClient *ethclient.Client) {
	// Run once immediately
	logger.Info("Starting initial points accrual run")
	latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
	if err != nil {
		logger.Error("cannot fetch latest block", "error", err)
	} else {
		currBlockNum := latestBlock.NumberU64()
		if err := updatePoints(db, logger, currBlockNum); err != nil {
			logger.Error("initial points accrual run failed", "error", err)
		} else {
			logger.Info("initial points accrual run completed successfully")
		}
	}

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()

		for range ticker.C {
			logger.Info("Starting points accrual run")

			// fetch latest block
			latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
			if err != nil {
				logger.Error("cannot fetch latest block", "error", err)
				continue
			}
			currBlockNum := latestBlock.NumberU64()

			// do the update inside a transaction
			if err := updatePoints(db, logger, currBlockNum); err != nil {
				logger.Error("points accrual run failed", "error", err)
			} else {
				logger.Info("points accrual run completed successfully")
			}
		}
	}()
}

// --------------------------------------
//   3) SERVICE + EVENT HANDLERS
// --------------------------------------

type PointsService struct {
	logger    *slog.Logger
	db        *sql.DB
	ethClient *ethclient.Client
	block     uint64
}

func (ps *PointsService) LastBlock() (uint64, error) {
	return ps.block, nil
}

func (ps *PointsService) SetLastBlock(block uint64) error {
	ps.block = block
	return nil
}

// Insert new row for opt-in
func insertOptIn(db *sql.DB, logger *slog.Logger, pubkey, adder, registryType, eventType string, inBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	var existingAdder string
	err := db.QueryRow(`
        SELECT adder FROM events
        WHERE pubkey = ? AND opted_out_block IS NULL
    `, pubkey).Scan(&existingAdder)

	// if pubkey is active under a different adder, ignore
	if err == nil && existingAdder != "" && existingAdder != adder {
		logger.Warn("pubkey already opted in by a different adder",
			"pubkey", pubkey, "existing_adder", existingAdder, "new_adder", adder)
		return
	}

	_, err = db.Exec(`
        INSERT INTO events (
            pubkey, adder, vault, registry_type, event_type, 
            opted_in_block, opted_out_block, points_accumulated
        ) VALUES (?, ?, NULL, ?, ?, ?, NULL, 0)
    `, pubkey, adder, registryType, eventType, inBlock)
	if err != nil {
		logger.Warn("insertOptIn likely already inserted", "error", err)
	} else {
		logger.Info("inserted opt-in interval",
			"pubkey", pubkey, "block", inBlock, "event_type", eventType, "adder", adder)
	}
}

// Insert new row for ValRecordAdded (symbiotic) with vault
// This is the same pattern, but we store vault explicitly
func insertOptInWithVault(db *sql.DB, logger *slog.Logger, pubkey, adder, vault, registryType, eventType string, inBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	var existingAdder string
	err := db.QueryRow(`
        SELECT adder FROM events
        WHERE pubkey = ? AND opted_out_block IS NULL
    `, pubkey).Scan(&existingAdder)

	if err == nil && existingAdder != "" && existingAdder != adder {
		logger.Warn("pubkey already opted in by a different adder",
			"pubkey", pubkey, "existing_adder", existingAdder, "new_adder", adder)
		return
	}

	_, err = db.Exec(`
        INSERT INTO events (
            pubkey, adder, vault, registry_type, event_type, 
            opted_in_block, opted_out_block, points_accumulated
        ) VALUES (?, ?, ?, ?, ?, ?, NULL, 0)
    `, pubkey, adder, vault, registryType, eventType, inBlock)
	if err != nil {
		logger.Warn("insertOptInWithVault likely already inserted", "error", err)
	} else {
		logger.Info("inserted opt-in interval WITH vault",
			"pubkey", pubkey, "adder", adder, "vault", vault, "block", inBlock, "event_type", eventType)
	}
}

// Mark an existing interval as opted-out
func insertOptOut(db *sql.DB, logger *slog.Logger, pubkey, adder, eventType string, outBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	_, err := db.Exec(`
        UPDATE events
        SET opted_out_block = ?, event_type = ?
        WHERE pubkey = ? AND adder = ? AND opted_out_block IS NULL
    `, outBlock, eventType, pubkey, adder)
	if err != nil {
		logger.Error("failed to opt-out", "error", err, "pubkey", pubkey)
	} else {
		logger.Info("opt-out interval updated",
			"pubkey", pubkey, "block", outBlock, "event_type", eventType, "adder", adder)
	}
}

// --------------------------------------
//   4) MAIN
// --------------------------------------

func main() {
	app := &cli.App{
		Name:  "mev-commit-points",
		Usage: "MEV Commit Points Service",
		Flags: []cli.Flag{optionRPCURL},
		Action: func(c *cli.Context) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			// 1. Connect to DB
			db, err := initDB(logger)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			// 2. Test DB
			var rowCount int
			err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&rowCount)
			if err != nil {
				return fmt.Errorf("failed to query events: %w", err)
			}
			logger.Info("database reachable", "events_count", rowCount)

			// 3. Ethereum client
			ethClient, err := ethclient.Dial(c.String(optionRPCURL.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node: %w", err)
			}

			// 4. Load ABIs
			contractABIs, err := getContractABIs()
			if err != nil {
				return fmt.Errorf("failed to get contract ABIs: %w", err)
			}
			listener := events.NewListener(logger, contractABIs...)

			// example addresses
			testnetContracts := []common.Address{
				common.HexToAddress("0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4"), // AVS
				common.HexToAddress("0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf"), // Vanilla
				common.HexToAddress("0x79FeCD427e5A3e5f1a40895A0AC20A6a50C95393"), // Symbiotic
			}

			// Start the publisher
			ps := &PointsService{
				logger:    logger,
				db:        db,
				ethClient: ethClient,
				block:     2146241, // example block
			}
			pub := publisher.NewHTTPPublisher(ps, logger, ethClient, listener)
			done := pub.Start(context.Background())
			for _, addr := range testnetContracts {
				pub.AddContract(addr)
			}

			// 5. Define event handlers
			handlers := []events.EventHandler{
				// Vanilla Staked (opt-in)
				events.NewEventHandler(
					"Staked",
					func(ev *vanillaregistry.Validatorregistryv1Staked, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.ValBLSPubKey)
						adder := ev.MsgSender.Hex()
						insertOptIn(db, logger, pubkey, adder, "vanilla", "Staked", blockNumber)
					},
				),
				// Vanilla Unstaked (opt-out)
				events.NewEventHandler(
					"Unstaked",
					func(ev *vanillaregistry.Validatorregistryv1Unstaked, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.ValBLSPubKey)
						adder := ev.MsgSender.Hex()
						insertOptOut(db, logger, pubkey, adder, "Unstaked", blockNumber)
					},
				),
				// Symbiotic: ValRecordAdded (opt-in) w/ vault
				events.NewEventHandler(
					"ValRecordAdded",
					func(ev *middleware.MevcommitmiddlewareValRecordAdded, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.BlsPubkey)
						adder := ev.Operator.Hex()
						vault := ev.Vault.Hex() // store vault
						pub.AddContract(ev.Vault)
						insertOptInWithVault(db, logger, pubkey, adder, vault, "symbiotic", "ValRecordAdded", blockNumber)
					},
				),
				// AVS: ValidatorRegistered (opt-in)
				events.NewEventHandler(
					"ValidatorRegistered",
					func(ev *avs.MevcommitavsValidatorRegistered, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.ValidatorPubKey)
						adder := ev.PodOwner.Hex()
						insertOptIn(db, logger, pubkey, adder, "eigenlayer", "ValidatorRegistered", blockNumber)
					},
				),
				// AVS: LSTRestakerRegistered (opt-in)
				events.NewEventHandler(
					"LSTRestakerRegistered",
					func(ev *avs.MevcommitavsLSTRestakerRegistered, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.ChosenValidator)
						adder := ev.LstRestaker.Hex()
						insertOptIn(db, logger, pubkey, adder, "eigenlayer", "LSTRestakerRegistered", blockNumber)
					},
				),
				// AVS: ValidatorDeregistered (opt-out)
				events.NewEventHandler(
					"ValidatorDeregistered",
					func(evt *avs.MevcommitavsValidatorDeregistered, blockNumber uint64) {
						pubkeyHex := common.Bytes2Hex(evt.ValidatorPubKey)
						adderHex := evt.PodOwner.Hex()
						insertOptOut(db, logger, pubkeyHex, adderHex, "ValidatorDeregistered", blockNumber)
					},
				),

				// VaultDeregistered => bulk opt-out
				events.NewEventHandler(
					"VaultDeregistered",
					func(ev *middleware.MevcommitmiddlewareVaultDeregistered, blockNumber uint64) {
						vaultAddr := ev.Vault.Hex()
						// We'll handle existing rows with the matching vault
						rows, err := db.Query(`
                            SELECT pubkey, adder
                            FROM events
                            WHERE vault = ?
                              AND opted_out_block IS NULL
                        `, vaultAddr)
						if err != nil {
							logger.Error("failed to query for vault", "error", err)
							return
						}
						defer rows.Close()

						for rows.Next() {
							var pubkey, adder string
							if err := rows.Scan(&pubkey, &adder); err != nil {
								logger.Error("scan error", "error", err)
								continue
							}
							insertOptOut(db, logger, pubkey, adder, "VaultDeregistered", blockNumber)
						}
					},
				),
				// OperatorDeregistered => bulk opt-out
				events.NewEventHandler(
					"OperatorDeregistered",
					func(ev *middleware.MevcommitmiddlewareOperatorDeregistered, blockNumber uint64) {
						operatorAddr := ev.Operator.Hex()
						rows, err := db.Query(`
                            SELECT pubkey
                            FROM events
                            WHERE adder = ?
                              AND opted_out_block IS NULL
                        `, operatorAddr)
						if err != nil {
							logger.Error("failed to query for operator", "error", err)
							return
						}
						defer rows.Close()

						for rows.Next() {
							var pubkey string
							if err := rows.Scan(&pubkey); err != nil {
								logger.Error("scan error", "error", err)
								continue
							}
							insertOptOut(db, logger, pubkey, operatorAddr, "OperatorDeregistered", blockNumber)
						}
					},
				),
				// ValRecordDeleted => single validator removal
				events.NewEventHandler(
					"ValRecordDeleted",
					func(ev *middleware.MevcommitmiddlewareValRecordDeleted, blockNumber uint64) {
						pubkeyHex := common.Bytes2Hex(ev.BlsPubkey)

						// find adder
						var adderHex string
						err := db.QueryRow(`
                            SELECT adder FROM events
                            WHERE pubkey = ?
                              AND opted_out_block IS NULL
                            LIMIT 1
                        `, pubkeyHex).Scan(&adderHex)
						if err != nil {
							logger.Error("failed to find active adder", "error", err, "pubkey", pubkeyHex)
							return
						}
						insertOptOut(db, logger, pubkeyHex, adderHex, "ValRecordDeleted", blockNumber)
					},
				),
			}

			// 6. Subscribe
			sub, err := listener.Subscribe(handlers...)
			if err != nil {
				return fmt.Errorf("subscribe error: %w", err)
			}
			defer sub.Unsubscribe()

			// 7. Start daily routine to recompute points
			StartPointsRoutine(db, logger, 24*time.Hour, ethClient)

			// watch for subscription errors
			go func() {
				for err := range sub.Err() {
					logger.Error("subscription error", "error", err)
				}
			}()

			// wait
			<-done
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// --------------------------------------
//   5) CONTRACT ABIs LOADER
// --------------------------------------

func getContractABIs() ([]*abi.ABI, error) {
	symbioticABI, err := abi.JSON(strings.NewReader(middleware.MevcommitmiddlewareABI))
	if err != nil {
		return nil, err
	}
	vanillaRegistryABI, err := abi.JSON(strings.NewReader(vanillaregistry.Validatorregistryv1ABI))
	if err != nil {
		return nil, err
	}
	avsABI, err := abi.JSON(strings.NewReader(avs.MevcommitavsABI))
	if err != nil {
		return nil, err
	}
	return []*abi.ABI{&symbioticABI, &vanillaRegistryABI, &avsABI}, nil
}
