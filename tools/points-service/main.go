package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	avs "github.com/primev/mev-commit/contracts-abi/clients/MevCommitAVS"
	middleware "github.com/primev/mev-commit/contracts-abi/clients/MevCommitMiddleware"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	vanillaregistry "github.com/primev/mev-commit/contracts-abi/clients/VanillaRegistry"
	vault "github.com/primev/mev-commit/contracts-abi/clients/Vault"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/urfave/cli/v2"

	_ "github.com/mattn/go-sqlite3"
)

var (
	blocksInOneMonth = int64(216000)

	monthlyIncrements = []int64{
		1000,
		1800,
		2500,
		3500,
		5000,
		10000,
	}

	rwLock sync.RWMutex

	createTableEventsQuery = `
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        pubkey TEXT NOT NULL,
        adder TEXT NOT NULL,
        vault TEXT,
        registry_type TEXT CHECK (registry_type IN ('vanilla', 'symbiotic', 'eigenlayer')),
        event_type TEXT,
        opted_in_block BIGINT NOT NULL,
        opted_out_block BIGINT,
        points_accumulated BIGINT DEFAULT 0,
        UNIQUE(pubkey, adder, opted_in_block)
    );`

	// New table for storing the last processed block:
	createTableLastBlockQuery = `
    CREATE TABLE IF NOT EXISTS last_processed_block (
        id INTEGER PRIMARY KEY CHECK (id = 1),
        last_block BIGINT NOT NULL
    );
    INSERT OR IGNORE INTO last_processed_block (id, last_block) VALUES (1, 2146240);
    `

	selectActiveRowsQuery = `
        SELECT id, opted_in_block
        FROM events
        WHERE opted_out_block IS NULL
    `
	updatePointsQuery = `
            UPDATE events
            SET points_accumulated = ?
            WHERE id = ?
        `
	countActiveRowsQuery = `
            SELECT COUNT(*) FROM events
            WHERE opted_out_block IS NULL
        `

	optionRPCURL = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "URL of the Ethereum RPC server",
		EnvVars: []string{"POINTS_ETH_RPC_URL"},
		Value:   "https://eth-holesky.g.alchemy.com/v2/0DDo7YeieNEucZX3jieFfzmzOCGTKAgp",
	}

	optionDBPath = &cli.StringFlag{
		Name:    "db-path",
		Usage:   "Path to the sqlite3 database file",
		EnvVars: []string{"POINTS_DB_PATH"},
		Value:   "./points.db",
	}
)

func initDB(logger *slog.Logger, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Create main events table
	if _, err := db.Exec(createTableEventsQuery); err != nil {
		return nil, fmt.Errorf("failed to create events table: %w", err)
	}
	// Create table for last processed block
	if _, err := db.Exec(createTableLastBlockQuery); err != nil {
		return nil, fmt.Errorf("failed to create last_processed_block table: %w", err)
	}

	logger.Info("database setup complete", slog.String("path", dbPath))
	return db, nil
}

func computePointsForMonths(blocksActive int64) int64 {
	fullMonths := blocksActive / blocksInOneMonth
	if fullMonths < 1 {
		return 0
	}
	if fullMonths > int64(len(monthlyIncrements)) {
		fullMonths = int64(len(monthlyIncrements))
	}
	return monthlyIncrements[fullMonths-1]
}

func updatePoints(db *sql.DB, logger *slog.Logger, currentBlock uint64) (retErr error) {
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
		} else if retErr != nil {
			_ = tx.Rollback()
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				logger.Error("failed to commit transaction", "error", commitErr)
				retErr = commitErr
			}
		}
	}()

	rows, queryErr := tx.Query(selectActiveRowsQuery)
	if queryErr != nil {
		return fmt.Errorf("failed to query events for points: %w", queryErr)
	}
	defer rows.Close()

	var count int
	if err := tx.QueryRow(countActiveRowsQuery).Scan(&count); err != nil {
		logger.Error("failed to count active validators", "error", err)
	}

	logger.Info("updating points", "current_block", currentBlock, "active_validators", count)

	for rows.Next() {
		var id int64
		var inBlock uint64
		if err := rows.Scan(&id, &inBlock); err != nil {
			logger.Error("scan error", "error", err)
			continue
		}
		blocksActive := int64(0)
		if currentBlock > inBlock {
			blocksActive = int64(currentBlock - inBlock)
		}
		totalPoints := computePointsForMonths(blocksActive)
		if _, updErr := tx.Exec(updatePointsQuery, totalPoints, id); updErr != nil {
			logger.Error("failed to update points", "error", updErr, "id", id)
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return fmt.Errorf("rows iteration error: %w", rowsErr)
	}
	return nil
}

func StartPointsRoutine(db *sql.DB, logger *slog.Logger, interval time.Duration, ethClient *ethclient.Client) {
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
			latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
			if err != nil {
				logger.Error("cannot fetch latest block", "error", err)
				continue
			}
			currBlockNum := latestBlock.NumberU64()
			if err := updatePoints(db, logger, currBlockNum); err != nil {
				logger.Error("points accrual run failed", "error", err)
			} else {
				logger.Info("points accrual run completed successfully")
			}
		}
	}()
}

// PointsService now fetches/stores last block in DB
type PointsService struct {
	logger    *slog.Logger
	db        *sql.DB
	ethClient *ethclient.Client
}

func (ps *PointsService) LastBlock() (uint64, error) {
	var blk uint64
	err := ps.db.QueryRow(`
        SELECT last_block 
        FROM last_processed_block 
        WHERE id = 1
    `).Scan(&blk)
	if err != nil {
		return 0, err
	}
	return blk, nil
}

func (ps *PointsService) SetLastBlock(block uint64) error {
	_, err := ps.db.Exec(`
        UPDATE last_processed_block 
        SET last_block = ? 
        WHERE id = 1
    `, block)
	return err
}

func insertOptIn(db *sql.DB, logger *slog.Logger, pubkey, adder, registryType, eventType string, inBlock uint64) {
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
        ) VALUES (?, ?, NULL, ?, ?, ?, NULL, 0)
    `, pubkey, adder, registryType, eventType, inBlock)
	if err != nil {
		logger.Warn("insertOptIn likely already inserted", "error", err)
	} else {
		logger.Info("inserted opt-in interval",
			"pubkey", pubkey, "block", inBlock, "event_type", eventType, "adder", adder)
	}
}

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

func main() {
	app := &cli.App{
		Name:  "mev-commit-points",
		Usage: "MEV Commit Points Service",
		Flags: []cli.Flag{optionRPCURL, optionDBPath},
		Action: func(c *cli.Context) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			// 1. Connect to DB
			db, err := initDB(logger, c.String(optionDBPath.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			// 2. Quick test
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

			// Example addresses
			testnetContracts := []common.Address{
				common.HexToAddress("0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4"),
				common.HexToAddress("0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf"),
				common.HexToAddress("0x79FeCD427e5A3e5f1a40895A0AC20A6a50C95393"),
			}

			ps := &PointsService{
				logger:    logger,
				db:        db,
				ethClient: ethClient,
			}
			pub := publisher.NewHTTPPublisher(ps, logger, ethClient, listener)
			done := pub.Start(context.Background())
			for _, addr := range testnetContracts {
				pub.AddContract(addr)
			}

			handlers := []events.EventHandler{
				events.NewEventHandler(
					"Staked",
					func(ev *vanillaregistry.Validatorregistryv1Staked) {
						pubkey := common.Bytes2Hex(ev.ValBLSPubKey)
						adder := ev.MsgSender.Hex()
						insertOptIn(db, logger, pubkey, adder, "vanilla", "Staked", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"Unstaked",
					func(ev *vanillaregistry.Validatorregistryv1Unstaked) {
						pubkey := common.Bytes2Hex(ev.ValBLSPubKey)
						adder := ev.MsgSender.Hex()
						insertOptOut(db, logger, pubkey, adder, "Unstaked", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"ValRecordAdded",
					func(ev *middleware.MevcommitmiddlewareValRecordAdded) {
						pubkey := common.Bytes2Hex(ev.BlsPubkey)
						adder := ev.Operator.Hex()
						vault := ev.Vault.Hex()
						pub.AddContract(ev.Vault)
						logger.Info("raw log data",
							"block_number", ev.Raw.BlockNumber,
							"tx_hash", ev.Raw.TxHash.Hex(),
							"block_hash", ev.Raw.BlockHash.Hex(),
							"log_index", ev.Raw.Index,
							"tx_index", ev.Raw.TxIndex)
						insertOptInWithVault(db, logger, pubkey, adder, vault, "symbiotic", "ValRecordAdded", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"ValidatorRegistered",
					func(ev *avs.MevcommitavsValidatorRegistered) {
						pubkey := common.Bytes2Hex(ev.ValidatorPubKey)
						adder := ev.PodOwner.Hex()
						insertOptIn(db, logger, pubkey, adder, "eigenlayer", "ValidatorRegistered", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"LSTRestakerRegistered",
					func(ev *avs.MevcommitavsLSTRestakerRegistered) {
						pubkey := common.Bytes2Hex(ev.ChosenValidator)
						adder := ev.LstRestaker.Hex()
						insertOptIn(db, logger, pubkey, adder, "eigenlayer", "LSTRestakerRegistered", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"ValidatorDeregistered",
					func(ev *avs.MevcommitavsValidatorDeregistered) {
						pubkeyHex := common.Bytes2Hex(ev.ValidatorPubKey)
						adderHex := ev.PodOwner.Hex()
						insertOptOut(db, logger, pubkeyHex, adderHex, "ValidatorDeregistered", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"VaultDeregistered",
					func(ev *middleware.MevcommitmiddlewareVaultDeregistered) {
						vaultAddr := ev.Vault.Hex()
						rows, err := db.Query(`
                            SELECT pubkey, adder
                            FROM events
                            WHERE vault = ? AND opted_out_block IS NULL
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
							insertOptOut(db, logger, pubkey, adder, "VaultDeregistered", ev.Raw.BlockNumber)
						}
					},
				),
				events.NewEventHandler(
					"OperatorDeregistered",
					func(ev *middleware.MevcommitmiddlewareOperatorDeregistered) {
						operatorAddr := ev.Operator.Hex()
						rows, err := db.Query(`
                            SELECT pubkey
                            FROM events
                            WHERE adder = ? AND opted_out_block IS NULL
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
							insertOptOut(db, logger, pubkey, operatorAddr, "OperatorDeregistered", ev.Raw.BlockNumber)
						}
					},
				),
				events.NewEventHandler(
					"ValRecordDeleted",
					func(ev *middleware.MevcommitmiddlewareValRecordDeleted) {
						pubkeyHex := common.Bytes2Hex(ev.BlsPubkey)
						var adderHex string
						err := db.QueryRow(`
                            SELECT adder 
                            FROM events
                            WHERE pubkey = ? AND opted_out_block IS NULL
                            LIMIT 1
                        `, pubkeyHex).Scan(&adderHex)
						if err != nil {
							logger.Error("failed to find active adder", "error", err, "pubkey", pubkeyHex)
							return
						}
						insertOptOut(db, logger, pubkeyHex, adderHex, "ValRecordDeleted", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"OnSlash",
					func(ev *vault.VaultOnSlash) {
						vaultAddr := ev.Raw.Address.Hex()
						rows, err := db.Query(`
						SELECT pubkey
						FROM events
						WHERE vault = ? AND opted_out_block IS NULL
						`, vaultAddr)
						if err != nil {
							logger.Error("failed to query pubkeys for vault", "vault", vaultAddr, "error", err)
							return
						}
						defer rows.Close()

						// 2. Store them in a slice
						var pubkeys [][]byte
						for rows.Next() {
							var pubkeyHex string
							if err := rows.Scan(&pubkeyHex); err != nil {
								logger.Error("scan pubkey error", "error", err)
								continue
							}
							pubkeys = append(pubkeys, common.FromHex(pubkeyHex)) // convert hex to []byte
						}

						// The block we'll check is onSlashEventBlock + 1
						checkBlockNum := ev.Raw.BlockNumber + 1

						go func(pubkeys [][]byte, checkBlock uint64) {
							// Wait for the chain to reach block+1, or poll for it:
							// (Simplest approach is just loop until the block is available)
							for {
								tipBlock, err := ethClient.BlockNumber(context.Background())
								if err != nil {
									time.Sleep(time.Second * 2)
									continue
								}
								if tipBlock >= checkBlock {
									break
								}
								time.Sleep(time.Second * 2)
							}

							// 2. Prepare the contract address and create a Caller
							routerAddr := common.HexToAddress("0x251Fbc993f58cBfDA8Ad7b0278084F915aCE7fc3")
							routerCaller, err := validatoroptinrouter.NewValidatoroptinrouterCaller(routerAddr, ethClient)
							if err != nil {
								panic(fmt.Sprintf("failed to create router caller: %v", err))
							}
							// 3. Query router contract at block+1
							callOpts := &bind.CallOpts{BlockNumber: big.NewInt(int64(checkBlock))}
							optInStatuses, err := routerCaller.AreValidatorsOptedIn(callOpts, pubkeys)
							if err != nil {
								logger.Error("failed areValidatorsOptedIn", "error", err)
								return
							}

							// 4. For each pubkey, if none of the three boolean flags are true, set them “optedOut”
							for i, status := range optInStatuses {
								if !status.IsVanillaOptedIn && !status.IsAvsOptedIn && !status.IsMiddlewareOptedIn {
									pubkeyHex := common.Bytes2Hex(pubkeys[i])
									logger.Info("validator is no longer opted in by any registry",
										"pubkey", pubkeyHex, "block", checkBlock,
									)
									// Now mark in DB as optedOut
									insertOptOut(db, logger, pubkeyHex, "??adder??", "OnSlashAutoOptOut", checkBlock)
								}
							}
						}(pubkeys, checkBlockNum)
					},
				),
			}

			sub, err := listener.Subscribe(handlers...)
			if err != nil {
				return fmt.Errorf("subscribe error: %w", err)
			}
			defer sub.Unsubscribe()

			// 7. Start daily routine
			StartPointsRoutine(db, logger, 24*time.Hour, ethClient)

			go func() {
				for err := range sub.Err() {
					logger.Error("subscription error", "error", err)
				}
			}()

			<-done
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

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
	vaultABI, err := abi.JSON(strings.NewReader(vault.VaultABI))
	if err != nil {
		return nil, err
	}
	return []*abi.ABI{&symbioticABI, &vanillaRegistryABI, &avsABI, &vaultABI}, nil
}
