package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/primev/mev-commit/indexer/pkg/ethclient"
	"github.com/primev/mev-commit/indexer/pkg/logutil"
	"github.com/primev/mev-commit/indexer/pkg/store"
	"github.com/urfave/cli/v2"
)

var parsedAddresses []string

func main() {
	app := &cli.App{
		Name:  "blockchain-indexer",
		Usage: "Index blockchain data into Elasticsearch",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ethereum-endpoint",
				EnvVars: []string{"INDEXER_ETHEREUM_ENDPOINT"},
				Value:   "http://localhost:8545",
				Usage:   "Ethereum node endpoint",
			},
			&cli.StringFlag{
				Name:    "elasticsearch-endpoint",
				EnvVars: []string{"INDEXER_ELASTICSEARCH_ENDPOINT"},
				Value:   "http://127.0.0.1:9200",
				Usage:   "Elasticsearch endpoint",
			},
			&cli.StringFlag{
				Name:    "es-username",
				EnvVars: []string{"INDEXER_ES_USERNAME"},
				Value:   "",
				Usage:   "Elasticsearch username",
			},
			&cli.StringFlag{
				Name:    "es-password",
				EnvVars: []string{"INDEXER_ES_PASSWORD"},
				Value:   "",
				Usage:   "Elasticsearch password",
			},
			&cli.DurationFlag{
				Name:    "index-interval",
				EnvVars: []string{"INDEXER_INDEX_INTERVAL"},
				Value:   15 * time.Second,
				Usage:   "Interval between indexing operations",
			},
			&cli.StringFlag{
				Name:    "log-level",
				EnvVars: []string{"INDEXER_LOG_LEVEL"},
				Value:   "info",
				Usage:   "Log level (debug, info, warn, error)",
			},
			&cli.StringFlag{
				Name:    "log-fmt",
				Usage:   "log format to use, options are 'text' or 'json'",
				EnvVars: []string{"INDEXER_LOG_FMT"},
				Value:   "text",
				Action: func(ctx *cli.Context, v string) error {
					if v != "text" && v != "json" {
						return fmt.Errorf("invalid log format: %s. Must be 'text' or 'json'", v)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "log-tags",
				Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
				EnvVars: []string{"INDEXER_LOG_TAGS"},
				Action: func(ctx *cli.Context, s string) error {
					for i, p := range strings.Split(s, ",") {
						if len(strings.Split(p, ":")) != 2 {
							return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
						}
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "account-addresses",
				EnvVars: []string{"INDEXER_ACCOUNT_ADDRESSES"},
				Value:   "0xfA0B0f5d298d28EFE4d35641724141ef19C05684",
				Usage:   "comma-separated account addresses",
				Action: func(c *cli.Context, value string) error {
					parsedAddresses = parseAddresses(value)
					return nil
				},
			},
			&cli.UintFlag{
				Name:    "min-blocks-to-fetch-account-addrs",
				EnvVars: []string{"INDEXER_MIN_BLOCK_TO_FETCH_ACCOUNT_ADDRS"},
				Value:   10,
				Usage:   "minimum number of blocks needed to pass prior to fetching account addresses",
			},
			&cli.DurationFlag{
				Name:    "timeout-to-fetch-account-addrs",
				EnvVars: []string{"INDEXER_TIMEOUT_TO_FETCH_ACCOUNT_ADDRS"},
				Value:   5 * time.Second,
				Usage:   "timeout in seconds to fetch account addresses",
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseAddresses(input string) []string {
	addresses := strings.Split(input, ",")
	for i, addr := range addresses {
		addresses[i] = strings.TrimSpace(addr)
	}
	return addresses
}

func parseLogTags(tagString string) map[string]string {
	tags := make(map[string]string)
	for _, p := range strings.Split(tagString, ",") {
		parts := strings.Split(p, ":")
		if len(parts) == 2 {
			tags[parts[0]] = parts[1]
		}
	}
	return tags
}

func run(cliCtx *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ethClient, err := ethclient.NewW3EthereumClient(cliCtx.String("ethereum-endpoint"))
	if err != nil {
		slog.Error("failed to create Ethereum client", "error", err)
		return err
	}

	esClient, err := store.NewESClient(cliCtx.String("elasticsearch-endpoint"), cliCtx.String("es-username"), cliCtx.String("es-password"))
	if err != nil {
		slog.Error("failed to create Elasticsearch client", "error", err)
		return err
	}
	defer func() {
		if err := esClient.Close(ctx); err != nil {
			slog.Error("Failed to close Elasticsearch client", "error", err)
		}
	}()

	config := Config{
		EthClient:                        ethClient,
		Storage:                          esClient,
		IndexInterval:                    cliCtx.Duration("index-interval"),
		AccountAddresses:                 parsedAddresses,
		MinBlocksToFetchAccountAddresses: cliCtx.Uint("min-blocks-to-fetch-account-addrs"),
		TimeoutToFetchAccountAddresses:   cliCtx.Duration("timeout-to-fetch-account-addrs"),
	}
	blockchainIndexer := NewBlockchainIndexer(config)

	logTags := parseLogTags(cliCtx.String("log-tags"))
	// Set log level
	err = logutil.SetLogLevel(cliCtx.String("log-level"), cliCtx.String("log-fmt"), logTags)
	if err != nil {
		return err
	}

	if err = blockchainIndexer.Start(ctx); err != nil {
		slog.Error("failed to start blockchain indexer", "error", err)
		return err
	}

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-c
	slog.Info("shutting down gracefully...")
	cancel()
	// Wait for some time to allow ongoing operations to complete
	time.Sleep(5 * time.Second)
	return nil
}
