package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"

	pb "github.com/primevprotocol/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

var (
	txHash      string
	amount      string
	blockNumber int64
)

type config struct {
	ServerAddress string `json:"server_address" yaml:"server_address"`
	LogFmt        string `json:"log_fmt" yaml:"log_fmt"`
	LogLevel      string `json:"log_level" yaml:"log_level"`
}

var (
	optionConfig = &cli.StringFlag{
		Name:     "config",
		Usage:    "path to config file",
		Required: true,
		EnvVars:  []string{"SEARCHER_CLI_CONFIG"},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "bidder-cli"
	app.Usage = "A CLI tool for interacting with a gRPC bidder server"
	app.Version = "1.0.0"

	var (
		cfg    config
		logger *slog.Logger
	)

	app.Flags = []cli.Flag{
		optionConfig,
	}

	app.Before = func(c *cli.Context) error {
		configFile := c.String(optionConfig.Name)

		buf, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file at '%s': %w", configFile, err)
		}

		if err := yaml.Unmarshal(buf, &cfg); err != nil {
			return fmt.Errorf("failed to unmarshal config file at '%s': %w", configFile, err)
		}

		if err := checkConfig(&cfg); err != nil {
			return fmt.Errorf("failed to unmarshal config file at '%s': %w", configFile, err)
		}

		logger, err = newLogger(cfg.LogLevel, cfg.LogFmt, c.App.Writer)
		if err != nil {
			return fmt.Errorf("failed to create logger: %w", err)
		}

		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:  "send-bid",
			Usage: "Send a bid to the gRPC bidder server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "txhash",
					Usage:       "Transaction hash",
					Destination: &txHash,
				},
				&cli.StringFlag{
					Name:        "amount",
					Usage:       "Bid amount",
					Destination: &amount,
				},
				&cli.Int64Flag{
					Name:        "block",
					Usage:       "Block number",
					Destination: &blockNumber,
				},
			},
			Action: func(c *cli.Context) error {
				if txHash == "" || amount == "" || blockNumber == 0 {
					return fmt.Errorf("Missing required arguments. Please provide --txhash, --amount, and --block.")
				}

				creds := insecure.NewCredentials()
				conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(creds))
				if err != nil {
					return err
				}
				defer conn.Close()

				client := pb.NewBidderClient(conn)

				bid := &pb.Bid{
					TxHashes:    []string{txHash},
					Amount:      amount,
					BlockNumber: blockNumber,
				}

				ctx := context.Background()
				stream, err := client.SendBid(ctx, bid)
				if err != nil {
					return err
				}

				preConfirmation, err := stream.Recv()
				if err != nil {
					return err
				}

				logger.Info("received preconfirmation", "preconfirmation", preConfirmation)
				return nil
			},
		},
		{
			// NOTE: (@iowar) By sending an empty Bid request, the status of the RPC
			// server is being checked. Instead, a ping request can be defined within
			// the bidder proto or a better solution can be found. Seeking the team's
			// opinion on this
			Name:  "status",
			Usage: "Check the status of the gRPC bidder server",
			Action: func(c *cli.Context) error {
				creds := insecure.NewCredentials()
				conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(creds))
				if err != nil {
					return err
				}
				defer conn.Close()

				client := pb.NewBidderClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
				defer cancel()

				_, err = client.SendBid(ctx, &pb.Bid{})
				if err != nil {
					logger.Info("gRPC bidder server is not reachable", "server", cfg.ServerAddress)
					return nil
				}

				logger.Info("gRPC bidder server is up and running", "server", cfg.ServerAddress)
				return nil
			},
		},
		{
			Name:  "send-rand-bid",
			Usage: "Send a random bid to the gRPC bidder server",
			Action: func(c *cli.Context) error {
				randSource := rand.NewSource(time.Now().UnixNano())
				randGenerator := rand.New(randSource)

				txHash = generateTxHash(randGenerator)
				amount = strconv.Itoa(randGenerator.Intn(1000) + 1)
				blockNumber = randGenerator.Int63n(100000) + 1

				creds := insecure.NewCredentials()
				conn, err := grpc.Dial(cfg.ServerAddress, grpc.WithTransportCredentials(creds))
				if err != nil {
					return err
				}
				defer conn.Close()

				client := pb.NewBidderClient(conn)

				bid := &pb.Bid{
					TxHashes:    []string{txHash},
					Amount:      amount,
					BlockNumber: blockNumber,
				}

				ctx := context.Background()
				stream, err := client.SendBid(ctx, bid)
				if err != nil {
					return err
				}

				preConfirmation, err := stream.Recv()
				if err != nil {
					return err
				}

				logger.Info("received preconfirmation", "preconfirmation", preConfirmation)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}

func checkConfig(cfg *config) error {
	if cfg.ServerAddress == "" {
		return fmt.Errorf("server_address is required")
	}

	if cfg.LogFmt == "" {
		cfg.LogFmt = "text"
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return nil
}

func newLogger(lvl, logFmt string, sink io.Writer) (*slog.Logger, error) {
	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(lvl)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	var (
		handler slog.Handler
		options = &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
	)
	switch logFmt {
	case "text":
		handler = slog.NewTextHandler(sink, options)
	case "json", "none":
		handler = slog.NewJSONHandler(sink, options)
	default:
		return nil, fmt.Errorf("invalid log format: %s", logFmt)
	}

	return slog.New(handler), nil
}

func generateTxHash(r *rand.Rand) string {
	const charset = "0123456789abcdef"
	result := make([]byte, 66)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	result = append([]byte("0x"), result...)
	return string(result)
}
