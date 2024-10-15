package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand/v2"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
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
		Name:    "l1-rpc-url",
		Usage:   "URL of the L1 RPC server",
		EnvVars: []string{"TRANSACTOR_L1_RPC_URL"},
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
		Name:  "l1-transactor",
		Usage: "Issue random transactions for L1 chain from specified accounts",
		Flags: []cli.Flag{
			optionKeystorePathPassword,
			optionL1RPCURL,
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

			txtors := make([]*transactorAccount, 0)
			for _, kp := range c.StringSlice(optionKeystorePathPassword.Name) {
				parts := strings.Split(kp, ":")
				t, err := newTransactorAccount(logger, parts[0], parts[1], l1RPC)
				if err != nil {
					return err
				}
				txtors = append(txtors, t)
			}

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case <-sigc:
					var err error
					for _, t := range txtors {
						err = errors.Join(err, t.Close())
					}
					return err
				default:
				}

				rand.Shuffle(len(txtors), func(i, j int) {
					txtors[i], txtors[j] = txtors[j], txtors[i]
				})

				from := txtors[0]
				to := txtors[1]

				// random amount between 0 and 1_000_000_000_000
				amount := big.NewInt(rand.Int64N(1_000_000_000_000) + 1)

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
				if err := from.SendTransaction(ctx, to.Address(), amount); err != nil {
					logger.Error("failed to send transaction", "error", err)
				}
				cancel()

				time.Sleep(100 * time.Millisecond)
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
