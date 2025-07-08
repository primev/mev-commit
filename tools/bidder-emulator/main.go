package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"math/rand/v2"
	"os"
	"os/signal"
	"path"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
	"resenje.org/multex"
)

var (
	optionKeystorePathPassword = &cli.StringSliceFlag{
		Name:    "keystore-path-password",
		Usage:   "Path to the keystore file and password in the format path:password",
		EnvVars: []string{"TRANSACTOR_KEYSTORE_PATH_PASSWORD"},
		Action: func(c *cli.Context, keystores []string) error {
			for _, kp := range keystores {
				parts := strings.Split(kp, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid keystore-path-password format: %s", kp)
				}
			}
			return nil
		},
	}

	optionL1RPCURL = &cli.StringFlag{
		Name:     "l1-rpc-url",
		Usage:    "URL of the L1 RPC server",
		EnvVars:  []string{"TRANSACTOR_L1_RPC_URL"},
		Required: true,
	}

	optionBidderRPCURL = &cli.StringFlag{
		Name:     "bidder-rpc-url",
		Usage:    "URL of the bidder RPC server",
		EnvVars:  []string{"TRANSACTOR_BIDDER_RPC_URL"},
		Required: true,
	}

	optionDepositAmount = &cli.StringFlag{
		Name:    "deposit-amount",
		Usage:   "Amount to deposit in wei",
		EnvVars: []string{"TRANSACTOR_DEPOSIT_AMOUNT"},
		Value:   "1000000000000000000", // Default to 1 ETH in wei
		Action: func(ctx *cli.Context, s string) error {
			if _, ok := new(big.Int).SetString(s, 10); !ok {
				return fmt.Errorf("invalid deposit amount: %s", s)
			}
			return nil
		},
	}

	optionBidWorkers = &cli.IntFlag{
		Name:    "bid-workers",
		Usage:   "Number of bid workers to run concurrently",
		EnvVars: []string{"TRANSACTOR_BID_WORKERS"},
		Value:   1,
		Action: func(ctx *cli.Context, i int) error {
			if i < 1 {
				return fmt.Errorf("invalid bid-workers, must be at least 1")
			}
			return nil
		},
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"TRANSACTOR_LOG_FMT"},
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
		EnvVars: []string{"TRANSACTOR_LOG_LEVEL"},
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
		EnvVars: []string{"TRANSACTOR_LOG_TAGS"},
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

func main() {
	app := &cli.App{
		Name:  "bidder-emulator",
		Usage: "Issue random transactions for L1 chain from specified accounts and bid on them on mev-commit chain",
		Flags: []cli.Flag{
			optionKeystorePathPassword,
			optionL1RPCURL,
			optionBidderRPCURL,
			optionDepositAmount,
			optionBidWorkers,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			l1RPC, err := ethclient.Dial(c.String(optionL1RPCURL.Name))
			if err != nil {
				return err
			}

			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			chainID, err := l1RPC.ChainID(context.Background())
			if err != nil {
				return fmt.Errorf("failed getting chain ID: %w", err)
			}

			txtors := make([]*transactorAccount, 0)
			for _, kp := range c.StringSlice(optionKeystorePathPassword.Name) {
				parts := strings.Split(kp, ":")
				t, err := newTransactorAccount(chainID, path.Dir(parts[0]), parts[1], l1RPC)
				if err != nil {
					return err
				}
				txtors = append(txtors, t)
			}

			bidderClient, err := newBidder(c.String(optionBidderRPCURL.Name), c.String(optionDepositAmount.Name))
			if err != nil {
				return fmt.Errorf("failed to create bidder client: %w", err)
			}

			ctx, cancel := signal.NotifyContext(c.Context, syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			eg, egCtx := errgroup.WithContext(ctx)
			for i := 0; i < c.Int(optionBidWorkers.Name); i++ {
				eg.Go(func() error {
					return bidWorker(egCtx, logger, txtors, bidderClient, l1RPC)
				})
			}

			return eg.Wait()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func bidWorker(
	ctx context.Context,
	logger *slog.Logger,
	txtors []*transactorAccount,
	bidderClient *bidder,
	l1RPC *ethclient.Client,
) error {
	if len(txtors) < 2 {
		return fmt.Errorf("at least 2 transactor accounts are required")
	}

	mtx := multex.New[string]()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		start := time.Now()

		rand.Shuffle(len(txtors), func(i, j int) {
			txtors[i], txtors[j] = txtors[j], txtors[i]
		})

		from := txtors[0]
		to := txtors[1]

		// random amount between 0 and 1_000_000_000_000
		amount := big.NewInt(rand.Int64N(1_000_000_000_000) + 1)

		mtx.Lock(from.Address().Hex())
		txn, err := from.SendTransaction(ctx, to.Address(), amount)
		if err != nil {
			mtx.Unlock(from.Address().Hex())
			logger.Error("failed to send transaction", "error", err)
			continue
		}
		mtx.Unlock(from.Address().Hex())

		signingDuration := time.Since(start)

		logger.Info("transaction sent", "from", from.Address(), "to", to.Address(), "amount", amount)
		rcpt, err := bind.WaitMined(ctx, l1RPC, txn)
		if err != nil {
			logger.Error("failed to wait for transaction receipt", "error", err)
			continue
		}

		waitDuration := time.Since(start) - signingDuration

		if rcpt.Status != types.ReceiptStatusSuccessful {
			logger.Error("transaction failed", "hash", txn.Hash().Hex(), "status", rcpt.Status)
			continue
		}
		logger.Info("transaction mined", "hash", txn.Hash().Hex(), "block", rcpt.BlockNumber)

		res, err := bidderClient.SendBid(ctx, txn, rcpt.BlockNumber.Int64())
		if err != nil {
			logger.Error("failed to send bid", "error", err)
			continue
		}

		preconfDuration := time.Since(start) - signingDuration - waitDuration

		logger.Info(
			"bid sent",
			"hash", txn.Hash().Hex(),
			"bidAmount", res.bid.Amount,
			"blockNumber", res.bid.BlockNumber,
			"noOfPreconfs", len(res.preconfs),
			"signingDuration", signingDuration,
			"waitDuration", waitDuration,
			"preconfDuration", preconfDuration,
		)
	}
}
