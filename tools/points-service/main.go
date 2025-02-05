package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

var (
	optionRPCURL = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "URL of the Ethereum RPC server",
		EnvVars: []string{"POINTS_ETH_RPC_URL"},
		Value:   "https://eth-holesky.g.alchemy.com/v2/0DDo7YeieNEucZX3jieFfzmzOCGTKAgp",
	}
)

var (
	optionDBPath = &cli.StringFlag{
		Name:    "db-path",
		Usage:   "Path to SQLite database file",
		EnvVars: []string{"POINTS_DB_PATH"},
		Value:   "points.db",
	}
)

func initDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./points.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Create tables if they don't exist
	result, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			block_number INTEGER,
			tx_hash TEXT,
			event_type TEXT,
			address TEXT,
			points_delta INTEGER,
			timestamp INTEGER,
			pubkey TEXT NOT NULL,
			vault_address TEXT,
			pubkey_poster_address TEXT NOT NULL,
			opted_in BOOLEAN DEFAULT TRUE,
			registry_type TEXT CHECK (registry_type IN ('vanilla', 'symbiotic', 'eigenlayer'))
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Warn("failed to get rows affected", "error", err)
	}
	logger.Debug("database tables created/verified", "rows_affected", rows)
	logger.Info("database created", slog.String("path", "./points.db"))

	return db, nil
}

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

func main() {
	app := &cli.App{
		Name:  "mev-commit-points",
		Usage: "MEV Commit Points Service",
		Flags: []cli.Flag{
			optionRPCURL,
		},
		Action: func(c *cli.Context) error {
			// Initialize logger
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			// Connect to database
			db, err := initDB(logger)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			// Verify database is accessible by running a test query
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to query database: %w", err)
			}
			logger.Info("database connection verified", "events_count", count)
			// Connect to Ethereum client
			ethClient, err := ethclient.Dial(c.String(optionRPCURL.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node: %w", err)
			}

			contractABIs, err := getContractABIs()
			if err != nil {
				return fmt.Errorf("failed to get contract ABIs: %w", err)
			}

			listener := events.NewListener(logger, contractABIs...)

			// contracts := []common.Address{
			// 	// TODO: fill this out
			// 	common.HexToAddress("0xBc77233855e3274E1903771675Eb71E602D9DC2e"), // AVS
			// 	common.HexToAddress("0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9"), // Vanilla Registry
			// 	common.HexToAddress("0x21fD239311B050bbeE7F32850d99ADc224761382"), // Symbiotic
			// }

			testnetcontracts := []common.Address{
				common.HexToAddress("0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4"), // AVS
				common.HexToAddress("0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf"), // Vanilla Registry
				common.HexToAddress("0x79FeCD427e5A3e5f1a40895A0AC20A6a50C95393"), // Symbiotic
			}

			// TODO(@ckartik): we'll need to make sure when we store into the DB, that we're not storing with a new address for same pubkey
			// if a pubkey exists in the DB, and we see an event with a different address putting in pubkey, we should reject it.
			handlers := []events.EventHandler{
				// Vanilla Registry handler
				events.NewEventHandler(
					"Staked",
					func(upd *vanillaregistry.Validatorregistryv1Staked, blockNumber uint64) {
						_, err := db.Exec(`
							INSERT INTO events (
								pubkey,
								pubkey_poster_address,
								event_type,
								block_number,
								opted_in,
								registry_type
							) VALUES (?, ?, ?, ?, ?, ?)`,
							common.Bytes2Hex(upd.ValBLSPubKey),
							upd.MsgSender.Hex(),
							"Staked",
							blockNumber,
							true,
							"vanilla",
						)
						if err != nil {
							logger.Error("failed to insert Staked event", "error", err)
						}
					},
				),

				// Middleware Registry handler
				events.NewEventHandler(
					"ValRecordAdded",
					func(upd *middleware.MevcommitmiddlewareValRecordAdded, blockNumber uint64) {
						_, err := db.Exec(`
							INSERT INTO events (
								pubkey,
								pubkey_poster_address,
								vault_address,
								event_type,
								block_number,
								opted_in,
								registry_type
							) VALUES (?, ?, ?, ?, ?, ?, ?)`,
							common.Bytes2Hex(upd.BlsPubkey),
							upd.Operator.Hex(),
							upd.Vault.Hex(),
							"ValRecordAdded",
							blockNumber,
							true,
							"symbiotic",
						)
						if err != nil {
							logger.Error("failed to insert ValRecordAdded event", "error", err)
						}
					},
				),

				// AVS Registry handlers
				events.NewEventHandler(
					"ValidatorRegistered",
					func(upd *avs.MevcommitavsValidatorRegistered, blockNumber uint64) {
						_, err := db.Exec(`
							INSERT INTO events (
								pubkey,
								pubkey_poster_address,
								event_type,
								block_number,
								opted_in,
								registry_type
							) VALUES (?, ?, ?, ?, ?, ?)`,
							common.Bytes2Hex(upd.ValidatorPubKey),
							upd.PodOwner.Hex(),
							"ValidatorRegistered",
							blockNumber,
							true,
							"eigenlayer",
						)
						if err != nil {
							logger.Error("failed to insert ValidatorRegistered event", "error", err)
						}
					},
				),

				events.NewEventHandler(
					"LSTRestakerRegistered",
					func(upd *avs.MevcommitavsLSTRestakerRegistered, blockNumber uint64) {
						_, err := db.Exec(`
							INSERT INTO events (
								pubkey,
								pubkey_poster_address,
								event_type,
								block_number,
								opted_in,
								registry_type
							) VALUES (?, ?, ?, ?, ?, ?)`,
							common.Bytes2Hex(upd.ChosenValidator),
							upd.LstRestaker.Hex(),
							"LSTRestakerRegistered",
							upd.Raw.BlockNumber,
							true,
							"eigenlayer",
						)
						if err != nil {
							logger.Error("failed to insert LSTRestakerRegistered event", "error", err)
						}
					},
				),
			}

			sub, err := listener.Subscribe(handlers...)
			if err != nil {
				return fmt.Errorf("failed to subscribe to events: %w", err)
			}
			defer sub.Unsubscribe()

			pointsService := &PointsService{block: 2146241}
			publisher := publisher.NewHTTPPublisher(
				pointsService,
				logger,
				ethClient,
				listener,
			)

			done := publisher.Start(context.Background(), testnetcontracts...)

			// Monitor subscription errors
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

	return []*abi.ABI{
		&symbioticABI,
		&vanillaRegistryABI,
		&avsABI,
	}, nil
}
