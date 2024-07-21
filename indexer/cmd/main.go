package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/primev/mev-commit/indexer/pkg/indexer"
	"github.com/urfave/cli/v2"
	"log/slog"
)

func main() {
	app := &cli.App{
		Name:  "blockchain-indexer",
		Usage: "Index blockchain data into Elasticsearch",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ethereum-endpoint",
				EnvVars: []string{"ETHEREUM_ENDPOINT"},
				Value:   "http://localhost:8545",
				Usage:   "Ethereum node endpoint",
			},
			&cli.StringFlag{
				Name:    "elasticsearch-endpoint",
				EnvVars: []string{"ELASTICSEARCH_ENDPOINT"},
				Value:   "http://localhost:9200",
				Usage:   "Elasticsearch endpoint",
			},
			&cli.DurationFlag{
				Name:    "index-interval",
				EnvVars: []string{"INDEX_INTERVAL"},
				Value:   15 * time.Second,
				Usage:   "Interval between indexing operations",
			},
			&cli.StringFlag{
				Name:    "log-level",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cliCtx *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ethClient, err := indexer.NewW3EthereumClient(cliCtx.String("ethereum-endpoint"))
	if err != nil {
		slog.Error("Failed to create Ethereum client", "error", err)
		return err
	}

	esClient, err := indexer.NewESClient(cliCtx.String("elasticsearch-endpoint"))
	if err != nil {
		slog.Error("Failed to create Elasticsearch client", "error", err)
		return err
	}
	defer func() {
		if err := esClient.Close(ctx); err != nil {
			slog.Error("Failed to close Elasticsearch client", "error", err)
		}
	}()

	blockchainIndexer := indexer.NewBlockchainIndexer(
		ethClient,
		esClient,
		cliCtx.Duration("index-interval"),
	)

	// Set log level
	indexer.SetLogLevel(cliCtx.String("log-level"))

	if err := blockchainIndexer.Start(ctx); err != nil {
		slog.Error("Failed to start blockchain indexer", "error", err)
		return err
	}

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-c
	slog.Info("Shutting down gracefully...")
	cancel()
	// Wait for some time to allow ongoing operations to complete
	time.Sleep(5 * time.Second)
	return nil
}
