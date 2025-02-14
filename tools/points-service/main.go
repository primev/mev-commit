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
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	avs "github.com/primev/mev-commit/contracts-abi/clients/MevCommitAVS"
	middleware "github.com/primev/mev-commit/contracts-abi/clients/MevCommitMiddleware"
	validatoroptinrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	vanillaregistry "github.com/primev/mev-commit/contracts-abi/clients/VanillaRegistry"
	vault "github.com/primev/mev-commit/contracts-abi/clients/Vault"
	config "github.com/primev/mev-commit/contracts-abi/config"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// ~216000 blocks is roughly one month on Ethereum (~2s block time)
var blocksInOneMonth = int64(216000)

var (
	rwLock                           sync.RWMutex
	createTableValidatorRecordsQuery = `
	CREATE TABLE IF NOT EXISTS validator_records (
		id SERIAL PRIMARY KEY,
		pubkey TEXT NOT NULL,
		adder TEXT NOT NULL,
		vault TEXT,
		registry_type TEXT CHECK (registry_type IN ('vanilla', 'symbiotic', 'eigenlayer')),
		event_type TEXT,
		opted_in_block BIGINT NOT NULL,
		opted_out_block BIGINT,
		points_accumulated BIGINT DEFAULT 0,
		pre_cliff_points BIGINT DEFAULT 0,
		UNIQUE(pubkey, adder, opted_in_block)
	);`

	createTableLastBlockQuery = `
	CREATE TABLE IF NOT EXISTS last_processed_block (
		id INT PRIMARY KEY,
		last_block BIGINT NOT NULL
	);
	`

	selectActiveValidatorRecordsQuery = `
	SELECT id, opted_in_block
	FROM validator_records
	WHERE opted_out_block IS NULL
	`

	updatePointsValidatorRecordsQuery = `
	UPDATE validator_records
	SET points_accumulated = $1, pre_cliff_points = $2
	WHERE id = $3
	`

	countActiveValidatorRecordsQuery = `
	SELECT COUNT(*) FROM validator_records
	WHERE opted_out_block IS NULL
	`

	optionRPCURL = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "URL of the Ethereum RPC server",
		EnvVars: []string{"POINTS_ETH_RPC_URL"},
		Value:   "https://eth.llamarpc.com",
	}

	optionPgHost = altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "pg-host",
		Usage:   "PostgreSQL host",
		EnvVars: []string{"POINTS_PG_HOST"},
		Value:   "localhost",
	})

	optionPgPort = &cli.IntFlag{
		Name:    "pg-port",
		Usage:   "PostgreSQL port",
		EnvVars: []string{"POINTS_PG_PORT"},
		Value:   5432,
	}

	optionPgUser = &cli.StringFlag{
		Name:    "pg-user",
		Usage:   "PostgreSQL user",
		EnvVars: []string{"POINTS_PG_USER"},
		Value:   "postgres",
	}

	optionPgPassword = &cli.StringFlag{
		Name:    "pg-password",
		Usage:   "PostgreSQL password",
		EnvVars: []string{"POINTS_PG_PASSWORD"},
		Value:   "postgres",
	}

	optionPgDbname = &cli.StringFlag{
		Name:    "pg-dbname",
		Usage:   "PostgreSQL database name",
		EnvVars: []string{"POINTS_PG_DBNAME"},
		Value:   "mev_oracle",
	}

	optionMainnet = &cli.BoolFlag{
		Name:    "mainnet",
		Usage:   "Use mainnet contracts",
		EnvVars: []string{"POINTS_MAINNET"},
		Value:   true,
	}

	optionStartBlock = &cli.Int64Flag{
		Name:    "start-block",
		Usage:   "Block number to start processing from",
		EnvVars: []string{"POINTS_START_BLOCK"},
		Value:   21344601,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"POINTS_LOG_FMT"},
		Value:   "text",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"POINTS_LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"POINTS_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	}
)

// initDB initializes the database, creates tables if needed, and ensures schema is in place.
func initDB(logger *slog.Logger, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec(createTableValidatorRecordsQuery); err != nil {
		return nil, fmt.Errorf("failed to create validator_records table: %w", err)
	}
	if _, err := db.Exec(createTableLastBlockQuery); err != nil {
		return nil, fmt.Errorf("failed to create last_processed_block table: %w", err)
	}

	logger.Info("database setup complete", slog.String("dsn", dsn))
	return db, nil
}

func computePointsForMonths(blocksActive int64) (int64, int64) {
	months := blocksActive / blocksInOneMonth
	if months < 1 {
		return 0, 0
	}
	chunk1Partial := []int64{1000, 2270, 3800, 5600, 7670, 10000}
	chunk2Partial := []int64{11983, 14506, 17570, 21173, 25317, 30000}
	chunk3Partial := []int64{34683, 39367, 44050, 48734, 53417, 58100}

	if months < 6 {
		totalPoints := chunk1Partial[months-1]
		fallbackPoints := months * 1000
		return totalPoints, fallbackPoints
	}
	if months == 6 {
		return 10000, 10000
	}
	if months <= 12 {
		totalPoints := chunk2Partial[months-7]
		if months < 12 {
			fallbackPoints := chunk1Partial[5] + (months-6)*1000
			return totalPoints, fallbackPoints
		}
		return chunk2Partial[5], chunk2Partial[5] // month=12
	}
	if months <= 18 {
		totalPoints := chunk3Partial[months-13]
		if months < 18 {
			fallbackPoints := chunk2Partial[5] + (months-12)*1000
			return totalPoints, fallbackPoints
		}
		return chunk3Partial[5], chunk3Partial[5] // month=18
	}
	return 58100, 58100
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

		totalPoints, preSixMonthPoints := computePointsForMonths(blocksActive)
		if _, updErr := tx.Exec(updatePointsValidatorRecordsQuery, totalPoints, preSixMonthPoints, id); updErr != nil {
			logger.Error("failed to update points", "error", updErr, "id", id)
		}
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return fmt.Errorf("rows iteration error: %w", rowsErr)
	}
	return nil
}

// StartPointsRoutine periodically accrues points
func StartPointsRoutine(
	ctx context.Context,
	db *sql.DB,
	logger *slog.Logger,
	interval time.Duration,
	ethClient *ethclient.Client,
	ps *PointsService,
) {
	ticker := time.NewTicker(interval)
	ps.SetPointsRoutineRunning(true)

	go func() {
		defer ticker.Stop()
		defer ps.SetPointsRoutineRunning(false)
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
	logger                  *slog.Logger
	db                      *sql.DB
	ethClient               *ethclient.Client
	mu                      sync.RWMutex
	pointsRoutineRunning    bool
	eventSubscriptionActive bool
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
        SET last_block = $1 
        WHERE id = 1
    `, block)
	return err
}

func (ps *PointsService) SetPointsRoutineRunning(running bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.pointsRoutineRunning = running
}

func (ps *PointsService) IsPointsRoutineRunning() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.pointsRoutineRunning
}

func (ps *PointsService) SetSubscriptionActive(active bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.eventSubscriptionActive = active
}

func (ps *PointsService) IsSubscriptionActive() bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.eventSubscriptionActive
}

func insertOptIn(db *sql.DB, logger *slog.Logger, pubkey, adder, registryType, eventType string, inBlock uint64) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	var existingAdder string
	err := db.QueryRow(`
        SELECT adder FROM validator_records
        WHERE pubkey = $1 AND opted_out_block IS NULL
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
        ) VALUES ($1, $2, NULL, $3, $4, $5, NULL, 0)
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
        WHERE pubkey = $1 AND opted_out_block IS NULL
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
        ) VALUES ($1, $2, $3, $4, $5, $6, NULL, 0)
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
        SET opted_out_block = $1
        WHERE pubkey = $2 AND adder = $3 AND opted_out_block IS NULL
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
			optionPgHost,
			optionPgPort,
			optionPgUser,
			optionPgPassword,
			optionPgDbname,
			optionMainnet,
			optionStartBlock,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			ctx, cancel := context.WithCancel(c.Context)
			defer cancel()

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-signalChan
				logger.Info("received termination signal, shutting down")
				cancel()
			}()

			dsn := fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				c.String(optionPgHost.Name),
				c.Int(optionPgPort.Name),
				c.String(optionPgUser.Name),
				c.String(optionPgPassword.Name),
				c.String(optionPgDbname.Name),
			)
			db, err := initDB(logger, dsn)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			// Equivalent of "INSERT OR IGNORE": use ON CONFLICT DO NOTHING
			_, err = db.Exec(`
				INSERT INTO last_processed_block (id, last_block) 
				VALUES (1, $1) 
				ON CONFLICT (id) DO NOTHING
			`, c.Int64(optionStartBlock.Name))
			if err != nil {
				return fmt.Errorf("failed to insert initial block: %w", err)
			}

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

			var contractAddresses []common.Address
			if c.Bool(optionMainnet.Name) {
				contractAddresses = []common.Address{
					common.HexToAddress(config.EthereumContracts.VanillaRegistry),
					common.HexToAddress(config.EthereumContracts.MevCommitAVS),
					common.HexToAddress(config.EthereumContracts.MevCommitMiddleware),
				}
			} else {
				contractAddresses = []common.Address{
					common.HexToAddress(config.HoleskyContracts.VanillaRegistry),
					common.HexToAddress(config.HoleskyContracts.MevCommitAVS),
					common.HexToAddress(config.HoleskyContracts.MevCommitMiddleware),
				}
			}
			pub.AddContracts(contractAddresses...)

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
					"StakeWithdrawn",
					func(ev *vanillaregistry.Validatorregistryv1StakeWithdrawn) {
						pubkey := common.Bytes2Hex(ev.ValBLSPubKey)
						adder := ev.MsgSender.Hex()
						insertOptOut(db, logger, pubkey, adder, "StakeWithdrawn", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"ValRecordAdded",
					func(ev *middleware.MevcommitmiddlewareValRecordAdded) {
						pubkey := common.Bytes2Hex(ev.BlsPubkey)
						adder := ev.Operator.Hex()
						vaultAddr := ev.Vault.Hex()
						pub.AddContracts(ev.Vault)
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
					"ValidatorDeregistrationRequested",
					func(ev *avs.MevcommitavsValidatorDeregistrationRequested) {
						pubkeyHex := common.Bytes2Hex(ev.ValidatorPubKey)
						adderHex := ev.PodOwner.Hex()
						insertOptOut(db, logger, pubkeyHex, adderHex, "ValidatorDeregistrationRequested", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"ValRecordDeleted",
					func(ev *middleware.MevcommitmiddlewareValRecordDeleted) {
						pubkey := common.Bytes2Hex(ev.BlsPubkey)
						adder := ev.MsgSender.Hex()
						insertOptOut(db, logger, pubkey, adder, "ValRecordDeleted", ev.Raw.BlockNumber)
					},
				),
				events.NewEventHandler(
					"VaultDeregistered",
					func(ev *middleware.MevcommitmiddlewareVaultDeregistered) {
						vaultAddr := ev.Vault.Hex()
						rows, err := db.Query(`
                            SELECT pubkey, adder
                            FROM validator_records
                            WHERE vault = $1 AND opted_out_block IS NULL
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
                            WHERE adder = $1 AND opted_out_block IS NULL
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
                            WHERE pubkey = $1 AND opted_out_block IS NULL
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
						WHERE vault = $1 AND opted_out_block IS NULL
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

							var routerAddr common.Address
							if c.Bool(optionMainnet.Name) {
								routerAddr = common.HexToAddress(config.EthereumContracts.ValidatorOptInRouter)
							} else {
								routerAddr = common.HexToAddress(config.HoleskyContracts.ValidatorOptInRouter)
							}
							routerCaller, err := validatoroptinrouter.NewValidatoroptinrouterCaller(routerAddr, ethClient)
							if err != nil {
								panic(fmt.Sprintf("failed to create router caller: %v", err))
							}
							callOpts := &bind.CallOpts{BlockNumber: big.NewInt(int64(checkBlock))}
							optInStatuses, err := routerCaller.AreValidatorsOptedIn(callOpts, pubkeys)
							if err != nil {
								logger.Error("failed AreValidatorsOptedIn", "error", err)
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

			ps.SetSubscriptionActive(true)
			go func() {
				for err := range sub.Err() {
					ps.SetSubscriptionActive(false)
					logger.Error("subscription error", "error", err)
				}
			}()

			StartPointsRoutine(ctx, db, logger, 30*time.Minute, ethClient, ps)

			pointsAPI := NewPointsAPI(logger, db, ps)
			go func() {
				if err := pointsAPI.StartAPIServer(ctx, ":8080"); err != nil {
					logger.Error("API server error", "error", err)
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
