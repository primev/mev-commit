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

// --------------------
//   INITIALIZE DB
// --------------------

func initDB(logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./points.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Single table storing intervals
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS events (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            pubkey TEXT NOT NULL,
            adder TEXT NOT NULL,
            registry_type TEXT CHECK (registry_type IN ('vanilla', 'symbiotic', 'eigenlayer')),
            event_type TEXT,
            opted_in_block BIGINT NOT NULL,
            opted_out_block BIGINT,   -- null if still opted in
            UNIQUE(pubkey, adder, opted_in_block)
        );
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to create events table: %w", err)
	}

	logger.Info("database setup complete", slog.String("path", "./points.db"))
	return db, nil
}

// --------------------
//   POINTS SERVICE
// --------------------

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

// ---------------------------
//   HELPER FUNCTIONS
// ---------------------------

// Insert a new row for opt-in (if not already active under a different adder).
func insertOptIn(db *sql.DB, logger *slog.Logger, pubkey, adder, registryType, eventType string, inBlock uint64) {
	// 1. Check if there's an active interval for this pubkey with a different adder.
	var existingAdder string
	err := db.QueryRow(`
        SELECT adder FROM events
        WHERE pubkey = ? AND opted_out_block IS NULL
    `, pubkey).Scan(&existingAdder)

	if err == nil && existingAdder != "" && existingAdder != adder {
		logger.Warn("pubkey already opted in by a different adder; ignoring new event",
			"pubkey", pubkey, "existing_adder", existingAdder, "new_adder", adder)
		return
	}

	// 2. Insert a new interval row (opted_out_block = NULL).
	_, err = db.Exec(`
        INSERT INTO events (pubkey, adder, registry_type, event_type, opted_in_block, opted_out_block)
        VALUES (?, ?, ?, ?, ?, NULL)
    `, pubkey, adder, registryType, eventType, inBlock)
	if err != nil {
		logger.Warn("insertOptIn likely already inserted, ignoring", "error", err)
	} else {
		logger.Info("inserted opt-in interval",
			"pubkey", pubkey, "block", inBlock, "event_type", eventType, "adder", adder)
	}
}

// Mark an existing interval as opted-out.
func insertOptOut(db *sql.DB, logger *slog.Logger, pubkey, adder, eventType string, outBlock uint64) {
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

// --------------------
//       MAIN
// --------------------

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

			// 2. Verify DB is reachable
			var rowCount int
			err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&rowCount)
			if err != nil {
				return fmt.Errorf("failed to query events table: %w", err)
			}
			logger.Info("database connection verified", "events_count", rowCount)

			// 3. Connect to Ethereum
			ethClient, err := ethclient.Dial(c.String(optionRPCURL.Name))
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum node: %w", err)
			}

			// Load ABIs
			contractABIs, err := getContractABIs()
			if err != nil {
				return fmt.Errorf("failed to get contract ABIs: %w", err)
			}
			listener := events.NewListener(logger, contractABIs...)

			// Example addresses
			testnetContracts := []common.Address{
				common.HexToAddress("0xEDEDB8ed37A43Fd399108A44646B85b780D85DD4"), // AVS
				common.HexToAddress("0x87D5F694fAD0b6C8aaBCa96277DE09451E277Bcf"), // Vanilla
				common.HexToAddress("0x79FeCD427e5A3e5f1a40895A0AC20A6a50C95393"), // Symbiotic
			}

			// 4. Define event handlers
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
				// Symbiotic: ValRecordAdded (opt-in)
				events.NewEventHandler(
					"ValRecordAdded",
					func(ev *middleware.MevcommitmiddlewareValRecordAdded, blockNumber uint64) {
						pubkey := common.Bytes2Hex(ev.BlsPubkey)
						adder := ev.Operator.Hex()
						insertOptIn(db, logger, pubkey, adder, "symbiotic", "ValRecordAdded", blockNumber)
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

				// --------------------
				//   EVENT HANDLERS FOR FINAL DEREGISTRATION
				// --------------------

				// 1) VaultDeregistered
				events.NewEventHandler(
					"VaultDeregistered",
					func(ev *middleware.MevcommitmiddlewareVaultDeregistered, blockNumber uint64) {
						vaultAddr := ev.Vault.Hex()

						// If your table doesn't store "vault" directly,
						// you'd need another approach. This example assumes a "vault" column exists.
						rows, err := db.Query(`
							SELECT pubkey, adder
							FROM events
							WHERE vault = ?
							  AND opted_out_block IS NULL
						`, vaultAddr)
						if err != nil {
							logger.Error("failed to query validators for vault", "error", err)
							return
						}
						defer rows.Close()

						for rows.Next() {
							var pubkey, adder string
							if err := rows.Scan(&pubkey, &adder); err != nil {
								logger.Error("failed to scan row", "error", err)
								continue
							}
							insertOptOut(db, logger, pubkey, adder, "VaultDeregistered", blockNumber)
						}
					},
				),

				// 2) OperatorDeregistered
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
							logger.Error("failed to query validators for operator", "error", err)
							return
						}
						defer rows.Close()

						for rows.Next() {
							var pubkey string
							if err := rows.Scan(&pubkey); err != nil {
								logger.Error("failed to scan operator row", "error", err)
								continue
							}
							insertOptOut(db, logger, pubkey, operatorAddr, "OperatorDeregistered", blockNumber)
						}
					},
				),

				// 3) ValRecordDeleted
				events.NewEventHandler(
					"ValRecordDeleted",
					func(ev *middleware.MevcommitmiddlewareValRecordDeleted, blockNumber uint64) {
						pubkeyHex := common.Bytes2Hex(ev.BlsPubkey)
						callerHex := ev.MsgSender.Hex()

						// If your DB only has (pubkey, adder),
						// find the correct adder for pubkey:
						err := db.QueryRow(`
							SELECT caller 
							FROM events
							WHERE pubkey = ? 
							  AND opted_out_block IS NULL
							LIMIT 1
						`, pubkeyHex).Scan(&callerHex)
						if err != nil {
							logger.Error("failed to find active adder for pubkey", "error", err, "pubkey", pubkeyHex)
							return
						}

						insertOptOut(db, logger, pubkeyHex, callerHex, "ValRecordDeleted", blockNumber)
					},
				),
			}

			// 5. Subscribe to events
			sub, err := listener.Subscribe(handlers...)
			if err != nil {
				return fmt.Errorf("failed to subscribe to events: %w", err)
			}
			defer sub.Unsubscribe()

			// Start the publisher
			ps := &PointsService{
				logger:    logger,
				db:        db,
				ethClient: ethClient,
				block:     2146241,
			}
			pub := publisher.NewHTTPPublisher(ps, logger, ethClient, listener)
			done := pub.Start(context.Background(), testnetContracts...)

			// Watch for subscription errors
			go func() {
				for err := range sub.Err() {
					logger.Error("subscription error", "error", err)
				}
			}()

			// Wait until shutdown
			<-done
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// --------------------
//
//	GET CONTRACT ABIs
//
// --------------------
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
