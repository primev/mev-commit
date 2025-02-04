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
)

var (
	optionRPCURL = &cli.StringFlag{
		Name:    "ethereum-rpc-url",
		Usage:   "URL of the Ethereum RPC server",
		EnvVars: []string{"POINTS_ETH_RPC_URL"},
		Value:   "wss://eth-holesky.g.alchemy.com/v2/0DDo7YeieNEucZX3jieFfzmzOCGTKAgp",
	}
)

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

			// // Connect to database
			// db, err := sql.Open("postgres", c.String(optionDBURL.Name))
			// if err != nil {
			// 	return fmt.Errorf("failed to connect to database: %w", err)
			// }
			// defer db.Close()

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

			handlers := []events.EventHandler{
				// Vanilla Registry handler
				events.NewEventHandler(
					"Staked",
					func(upd *vanillaregistry.Validatorregistryv1Staked) {
						logger.Info("Vanilla Registry Staked event",
							"sender", upd.MsgSender.Hex(),
							"withdrawalAddress", upd.WithdrawalAddress.Hex(),
							"pubkey", common.Bytes2Hex(upd.ValBLSPubKey),
							"amount", upd.Amount.String(),
						)
					},
				),

				// Middleware Registry handler
				events.NewEventHandler(
					"ValRecordAdded",
					func(upd *middleware.MevcommitmiddlewareValRecordAdded) {
						logger.Info("Middleware ValRecordAdded event",
							"pubkey", common.Bytes2Hex(upd.BlsPubkey),
							"operator", upd.Operator.Hex(),
							"vault", upd.Vault.Hex(),
							"position", upd.Position.String(),
						)
					},
				),

				// AVS Registry handlers
				events.NewEventHandler(
					"ValidatorRegistered",
					func(upd *avs.MevcommitavsValidatorRegistered) {
						logger.Info("AVS ValidatorRegistered event",
							"pubkey", common.Bytes2Hex(upd.ValidatorPubKey),
							"podOwner", upd.PodOwner.Hex(),
						)
					},
				),
				events.NewEventHandler(
					"LSTRestakerRegistered",
					func(upd *avs.MevcommitavsLSTRestakerRegistered) {
						logger.Info("AVS LSTRestakerRegistered event",
							"pubkey", common.Bytes2Hex(upd.ChosenValidator),
							"numChosen", upd.NumChosen.String(),
							"lstRestaker", upd.LstRestaker.Hex(),
						)
					},
				),
			}

			sub, err := listener.Subscribe(handlers...)
			if err != nil {
				return fmt.Errorf("failed to subscribe to events: %w", err)
			}
			defer sub.Unsubscribe()

			pointsService := &PointsService{block: 2146241}
			publisher := publisher.NewWSPublisher(
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
