package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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
	config "github.com/primev/mev-commit/contracts-abi/config"
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

	createTableValidatorRecordsQuery = `
    CREATE TABLE IF NOT EXISTS validator_records (
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

	createTableLastBlockQuery = `
    CREATE TABLE IF NOT EXISTS last_processed_block (
        id INTEGER PRIMARY KEY CHECK (id = 1),
        last_block BIGINT NOT NULL
    );
    INSERT OR IGNORE INTO last_processed_block (id, last_block) VALUES (1, 2146240);
    `

	selectActiveValidatorRecordsQuery = `
        SELECT id, opted_in_block
        FROM validator_records
        WHERE opted_out_block IS NULL
    `
	updatePointsValidatorRecordsQuery = `
            UPDATE validator_records
            SET points_accumulated = ?
            WHERE id = ?
        `
	countActiveValidatorRecordsQuery = `
            SELECT COUNT(*) FROM validator_records
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

	optinRouterAddress = &cli.StringFlag{
		Name:    "router-address",
		Usage:   "Address of the ValidatorOptInRouter contract",
		EnvVars: []string{"POINTS_ROUTER_ADDRESS"},
		Value:   config.HoleskyContracts.ValidatorOptInRouter,
	}

	optionContractAddresses = &cli.StringSliceFlag{
		Name:    "contract-addresses",
		Usage:   "Addresses of the contracts to watch",
		EnvVars: []string{"POINTS_CONTRACT_ADDRESSES"},
		Value: cli.NewStringSlice(
			config.HoleskyContracts.VanillaRegistry,
			config.HoleskyContracts.MevCommitAVS,
			config.HoleskyContracts.MevCommitMiddleware,
		),
	}
)

func initDB(logger *slog.Logger, dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if _, err := db.Exec(createTableValidatorRecordsQuery); err != nil {
		return nil, fmt.Errorf("failed to create validator_records table: %w", err)
	}
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
			if rerr := tx.Rollback(); rerr != nil {
				retErr = errors.Join(retErr, rerr)
			}
			panic(p)
		} else if retErr != nil {
			if rerr := tx.Rollback(); rerr != nil {
				retErr = errors.Join(retErr, rerr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				logger.Error("failed to commit transaction", "error", commitErr)
				retErr = errors.Join(retErr, commitErr)
			}
		}
	}()

	rows, queryErr := tx.Query(selectActiveValidatorRecordsQuery)
	if queryErr != nil {
		return fmt.Errorf("failed to query validator_records for points: %w", queryErr)
	}
	defer rows.Close()

	var count int
	if err := tx.QueryRow(countActiveValidatorRecordsQuery).Scan(&count); err != nil {
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
		if _, updErr := tx.Exec(updatePointsValidatorRecordsQuery, totalPoints, id); updErr != nil {
			logger.Error("failed to update points", "error", updErr, "id", id)
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return fmt.Errorf("rows iteration error: %w", rowsErr)
	}
	return nil
}

func StartPointsRoutine(ctx context.Context, db *sql.DB, logger *slog.Logger, interval time.Duration, ethClient *ethclient.Client) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for {
			logger.Info("Starting points accrual run")
			latestBlock, err := ethClient.BlockByNumber(context.Background(), nil)
			if err != nil {
				logger.Error("cannot fetch latest block", "error", err)
			} else {
				currBlockNum := latestBlock.NumberU64()
				if err := updatePoints(db, logger, currBlockNum); err != nil {
					logger.Error("points accrual run failed", "error", err)
				} else {
					logger.Info("points accrual run completed successfully")
				}
			}

			select {
			case <-ctx.Done():
				logger.Info("points accrual routine shutting down")
				return
			case <-ticker.C:
			}
		}
	}()
}

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
        SELECT adder FROM validator_records
        WHERE pubkey = ? AND opted_out_block IS NULL
    `, pubkey).Scan(&existingAdder)

	if err == nil && existingAdder != "" && existingAdder != adder {
		logger.Warn("pubkey already opted in by a different adder",
			"pubkey", pubkey, "existing_adder", existingAdder, "new_adder", adder)
		return
	}

	_, err = db.Exec(`
        INSERT INTO validator_records (
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

func insertOptInWithVault(db *sql.DB, logger *slog.Logger, pubkey, adder, vaultAddr, registryType, eventType string, inBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	var existingAdder string
	err := db.QueryRow(`
        SELECT adder FROM validator_records
        WHERE pubkey = ? AND opted_out_block IS NULL
    `, pubkey).Scan(&existingAdder)

	if err == nil && existingAdder != "" && existingAdder != adder {
		logger.Warn("pubkey already opted in by a different adder",
			"pubkey", pubkey, "existing_adder", existingAdder, "new_adder", adder)
		return
	}

	_, err = db.Exec(`
        INSERT INTO validator_records (
            pubkey, adder, vault, registry_type, event_type, 
            opted_in_block, opted_out_block, points_accumulated
        ) VALUES (?, ?, ?, ?, ?, ?, NULL, 0)
    `, pubkey, adder, vaultAddr, registryType, eventType, inBlock)
	if err != nil {
		logger.Warn("insertOptInWithVault likely already inserted", "error", err)
	} else {
		logger.Info("inserted opt-in interval WITH vault",
			"pubkey", pubkey, "adder", adder, "vault", vaultAddr, "block", inBlock, "event_type", eventType)
	}
}

func insertOptOut(db *sql.DB, logger *slog.Logger, pubkey, adder, eventType string, outBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	_, err := db.Exec(`
        UPDATE validator_records
        SET opted_out_block = ?
        WHERE pubkey = ? AND adder = ? AND opted_out_block IS NULL
    `, outBlock, pubkey, adder)
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
		Flags: []cli.Flag{
			optionRPCURL,
			optionDBPath,
			optionContractAddresses, // Include our new flag
		},
		Action: func(c *cli.Context) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-signalChan
				logger.Info("received termination signal, shutting down")
				cancel()
			}()

			db, err := initDB(logger, c.String(optionDBPath.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			var rowCount int
			err = db.QueryRow("SELECT COUNT(*) FROM validator_records").Scan(&rowCount)
			if err != nil {
				return fmt.Errorf("failed to query validator_records: %w", err)
			}
			logger.Info("database reachable", "validator_records_count", rowCount)

			ethClient, err := ethclient.Dial(c.String(optionRPCURL.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node: %w", err)
			}

			contractABIs, err := getContractABIs()
			if err != nil {
				return fmt.Errorf("failed to get contract ABIs: %w", err)
			}
			listener := events.NewListener(logger, contractABIs...)

			ps := &PointsService{
				logger:    logger,
				db:        db,
				ethClient: ethClient,
			}
			pub := publisher.NewHTTPPublisher(ps, logger, ethClient, listener)
			done := pub.Start(ctx)

			// Get the contract addresses from CLI and add them to the publisher
			var contractAddresses []common.Address
			for _, addrStr := range c.StringSlice(optionContractAddresses.Name) {
				contractAddresses = append(contractAddresses, common.HexToAddress(addrStr))
			}
			for _, addr := range contractAddresses {
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
						vaultAddr := ev.Vault.Hex()
						pub.AddContract(ev.Vault)
						insertOptInWithVault(db, logger, pubkey, adder, vaultAddr, "symbiotic", "ValRecordAdded", ev.Raw.BlockNumber)
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
                            FROM validator_records
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
                            FROM validator_records
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
                            FROM validator_records
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
						SELECT pubkey, adder
						FROM validator_records
						WHERE vault = ? AND opted_out_block IS NULL
						`, vaultAddr)
						if err != nil {
							logger.Error("failed to query pubkeys for vault", "vault", vaultAddr, "error", err)
							return
						}
						defer rows.Close()

						var pubkeys [][]byte
						var adders []string
						for rows.Next() {
							var pubkeyHex, adderHex string
							if err := rows.Scan(&pubkeyHex, &adderHex); err != nil {
								logger.Error("scan pubkey error", "error", err)
								continue
							}
							pubkeys = append(pubkeys, common.FromHex(pubkeyHex))
							adders = append(adders, adderHex)
						}

						checkBlockNum := ev.Raw.BlockNumber + 1
						go func(pubkeys [][]byte, checkBlock uint64) {
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
							routerAddr := common.HexToAddress(c.String(optinRouterAddress.Name))
							routerCaller, err := validatoroptinrouter.NewValidatoroptinrouterCaller(routerAddr, ethClient)
							if err != nil {
								panic(fmt.Sprintf("failed to create router caller: %v", err))
							}
							callOpts := &bind.CallOpts{BlockNumber: big.NewInt(int64(checkBlock))}
							optInStatuses, err := routerCaller.AreValidatorsOptedIn(callOpts, pubkeys)
							if err != nil {
								logger.Error("failed areValidatorsOptedIn", "error", err)
								return
							}
							for i, status := range optInStatuses {
								if !status.IsMiddlewareOptedIn {
									pubkeyHex := common.Bytes2Hex(pubkeys[i])
									logger.Info("validator is no longer opted in by any registry",
										"pubkey", pubkeyHex, "block", checkBlock,
									)
									insertOptOut(db, logger, pubkeyHex, adders[i], "OnSlashAutoOptOut", checkBlock)
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

			StartPointsRoutine(ctx, db, logger, 24*time.Hour, ethClient)

			pointsAPI := NewPointsAPI(logger, db, ps)
			go func() {
				if err := pointsAPI.StartAPIServer(ctx, ":8080"); err != nil {
					logger.Error("API server error", "error", err)
				}
			}()

			go func() {
				for err := range sub.Err() {
					logger.Error("subscription error", "error", err)
				}
			}()

			select {
			case <-ctx.Done():
				logger.Info("context canceled, shutting down")
			case <-done:
				logger.Info("publisher done channel closed")
			}

			logger.Info("graceful shutdown complete")
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
