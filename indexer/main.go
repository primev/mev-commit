package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var (
	ethereumURLFlag = &cli.StringFlag{
		Name:    "ethereum-url",
		EnvVars: []string{"INDEXER_ETHEREUM_URL"},
		Value:   "http://localhost:8545",
		Usage:   "Ethereum node URL",
		Action: func(_ *cli.Context, s string) error {
			_, err := url.Parse(s)
			return err
		},
	}

	elasticsearchURLFlag = &cli.StringFlag{
		Name:    "elasticsearch-url",
		EnvVars: []string{"INDEXER_ELASTICSEARCH_URL"},
		Value:   "http://127.0.0.1:9200",
		Usage:   "Elasticsearch URL",
		Action: func(_ *cli.Context, s string) error {
			_, err := url.Parse(s)
			return err
		},
	}

	elasticsearchUsernameFlag = &cli.StringFlag{
		Name:    "elasticsearch-username",
		EnvVars: []string{"INDEXER_ELASTICSEARCH_USERNAME"},
		Value:   "",
		Usage:   "Elasticsearch username",
	}

	elasticsearchPasswordFlag = &cli.StringFlag{
		Name:    "elasticsearch-password",
		EnvVars: []string{"INDEXER_ELASTICSEARCH_PASSWORD"},
		Value:   "",
		Usage:   "Elasticsearch password",
	}

	logFormatFlag = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"INDEXER_LOG_FMT"},
		Value:   "text",
		Action: func(_ *cli.Context, v string) error {
			if v != "text" && v != "json" {
				return fmt.Errorf("invalid log format: %s", v)
			}
			return nil
		},
	}

	logLevelFlag = &cli.StringFlag{
		Name:    "log-level",
		EnvVars: []string{"INDEXER_LOG_LEVEL"},
		Value:   "info",
		Usage:   "Log level (debug, info, warn, error)",
		Action: func(_ *cli.Context, v string) error {
			if err := new(slog.LevelVar).UnmarshalText([]byte(v)); err != nil {
				return fmt.Errorf("invalid log level: %w", err)
			}
			return nil
		},
	}

	logTagsFlag = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"INDEXER_LOG_TAGS"},
		Action: func(_ *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid tag %q at index %d", p, i)
				}
			}
			return nil
		},
	}
)

func main() {
	app := &cli.App{
		Name:  "indexer",
		Usage: "Index Ethereum blockchain data into Elasticsearch",
		Flags: []cli.Flag{
			ethereumURLFlag,
			elasticsearchURLFlag,
			elasticsearchUsernameFlag,
			elasticsearchPasswordFlag,
			logFormatFlag,
			logLevelFlag,
			logTagsFlag,
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)

	}
}

func run(cctx *cli.Context) error {
	var (
		closers []io.Closer
	)

	ctx, stop := signal.NotifyContext(
		cctx.Context,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	logger, err := newLogger(
		os.Stdout,
		cctx.String(logLevelFlag.Name),
		cctx.String(logFormatFlag.Name),
		cctx.String(logTagsFlag.Name),
	)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	eth, err := NewEthereumChain(
		cctx.String(ethereumURLFlag.Name),
	)
	if err != nil {
		return fmt.Errorf("create Ethereum chain: %w", err)
	}
	closers = append(closers, eth)

	esc, err := NewElasticsearchStore(
		cctx.String(elasticsearchURLFlag.Name),
		cctx.String(elasticsearchUsernameFlag.Name),
		cctx.String(elasticsearchPasswordFlag.Name),
	)
	if err != nil {
		return fmt.Errorf("create Elasticsearch store: %w", err)
	}

	if err := esc.CreateIndexes(ctx); err != nil {
		return fmt.Errorf("create indexes: %w", err)
	}

	last, err := esc.LastIndexedBlock(ctx, SortOrderAsc)
	if err != nil {
		return fmt.Errorf("retrive last indexed block: %w", err)
	}
	logger.Info("last indexed block", "number", last)

	curr, err := eth.BlockNumber()
	if err != nil {
		return fmt.Errorf("retrive block number: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		n := new(big.Int).Set(last)
		if n.Sign() == 0 {
			n.Add(curr, big.NewInt(-1))
		}
		return (&ForwardIndexer{
			Config: &Config{
				logger: logger.With("indexer", "forward"),
				chain:  eth,
				store:  esc,
			},
			lastIndexedBlock: n,
		}).Run(ctx)
	})
	g.Go(func() error {
		n := new(big.Int).Set(last)
		if n.Sign() == 0 {
			n = new(big.Int).Set(curr)
		}
		return (&BackwardIndexer{
			Config: &Config{
				logger: logger.With("indexer", "backward"),
				chain:  eth,
				store:  esc,
			},
			lastIndexedBlock: n,
		}).Run(ctx)
	})
	g.Go(func() error {
		n := new(big.Int).Set(last)
		// TODO: finish implementation...
		return (&BalanceIndexer{
			Config: &Config{
				logger: logger.With("indexer", "balance"),
				chain:  eth,
				store:  esc,
			},
			lastIndexedBlock: n,
		}).Run(ctx)
	})

	select {
	case <-ctx.Done():
		logger.Info("shutting down...")
		var errs error
		if err := g.Wait(); !errors.Is(err, context.Canceled) {
			errs = errors.Join(errs, err)
		}
		for _, c := range closers {
			errs = errors.Join(errs, c.Close())
		}
		logger.Info("shutdown complete")
		return errs
	}
}
