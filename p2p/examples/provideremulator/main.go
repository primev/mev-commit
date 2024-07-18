package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/primev/mev-commit/p2p/examples/provideremulator/client"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
)

var (
	serverAddr = flag.String(
		"server-addr",
		"localhost:13524",
		"The server address in the format of host:port",
	)
	logLevel = flag.String(
		"log-level",
		"debug",
		"Verbosity level (debug|info|warn|error)",
	)
)

func main() {
	flag.Parse()
	if *serverAddr == "" {
		fmt.Println("Please provide a valid server address with the -server-addr flag")
		return
	}

	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		level.Set(slog.LevelDebug)
		fmt.Printf("invalid log level: %s; using %q", err, level)
	}

	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: level},
	))

	providerClient, err := client.NewProviderClient(*serverAddr, logger)
	if err != nil {
		logger.Error("failed to create provider client", "error", err)
		return
	}

	bidS, err := providerClient.ReceiveBids()
	if err != nil {
		logger.Error("failed to create bid receiver", "error", err)
		return
	}

	for {
		select {
		case bid, more := <-bidS:
			if !more {
				logger.Warn("closed bid stream")
				return
			}
			logger.Info("received new bid", "bid", bid)
			err := providerClient.SendBidResponse(context.Background(), &providerapiv1.BidResponse{
				BidDigest: bid.BidDigest,
				Status:    providerapiv1.BidResponse_STATUS_ACCEPTED,
			})
			if err != nil {
				logger.Error("failed to send bid response", "error", err)
				return
			}
			logger.Info("accepted bid")
		}
	}
}
