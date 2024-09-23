package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/x/util"
	"github.com/primev/mev-commit/x/util/otelutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
)

// The following const block contains the name of the cli flags, especially
// for reuse purposes.
const (
	serverAddrFlagName       = "server-addr"
	logLevelFlagName         = "log-level"
	logFmtFlagName           = "log-fmt"
	logTagsFlagName          = "log-tags"
	httpPortFlagName         = "http-port"
	errorProbabilityFlagName = "error-probability"
)

var (
	serverAddr = flag.String(
		serverAddrFlagName,
		"localhost:13524",
		"The server address in the format of host:port",
	)
	logLevel = flag.String(
		logLevelFlagName,
		"debug",
		"Verbosity level (debug|info|warn|error)",
	)
	logFmt = flag.String(
		logFmtFlagName,
		"text",
		"Format of the log output: 'text', 'json'",
	)
	logTags = flag.String(
		logTagsFlagName,
		"",
		"Comma-separated list of <name:value> pairs that will be inserted into each log line",
	)
	httpPort = flag.Int(
		httpPortFlagName,
		8080,
		"The port to serve the HTTP metrics endpoint on",
	)
	otelCollectorEndpointURL = flag.String(
		"otel-collector-endpoint-url",
		"",
		"URL for OpenTelemetry collector endpoint",
	)
	errorProbability = flag.Int(
		errorProbabilityFlagName,
		20,
		"The probability of returning an error when sending a bid response",
	)
)

var (
	receivedBids = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mev_commit",
		Subsystem: "provider_emulator",
		Name:      "total_received_bids",
		Help:      "Total number of bids received from mev_commit nodes",
	})
	sentBids = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mev_commit",
		Subsystem: "provider_emulator",
		Name:      "total_sent_bids",
		Help:      "Total number of bids sent mev_commit nodes",
	})
	rejectedBids = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mev_commit",
		Subsystem: "provider_emulator",
		Name:      "total_rejected_bids",
		Help:      "Total number of bids rejected",
	})
)

func main() {
	flag.Parse()

	logger, err := util.NewLogger(*logLevel, *logFmt, *logTags, os.Stdout)
	if err != nil {
		fmt.Printf("failed to create logger: %v", err)
		return
	}

	if *otelCollectorEndpointURL != "" {
		logger.Info("setting up OpenTelemetry SDK", "endpoint", *otelCollectorEndpointURL)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown, err := otelutil.SetupOTelSDK(
			ctx,
			*otelCollectorEndpointURL,
			*logTags,
		)
		if err != nil {
			logger.Warn("failed to setup OpenTelemetry SDK; continuing without telemetry", "error", err)
		} else {
			otel.SetLogger(logr.FromSlogHandler(
				logger.Handler().WithAttrs([]slog.Attr{
					{Key: "component", Value: slog.StringValue("otel")},
				}),
			))
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				err = errors.Join(err, shutdown(ctx))
				cancel()
			}()
		}
	}

	if *serverAddr == "" {
		fmt.Printf("please provide a valid server address with the -%s flag\n", serverAddrFlagName)
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(receivedBids, sentBids)

	go func() {
		router := http.NewServeMux()
		router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", *httpPort),
			Handler: router,
		}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", "error", err)
		}
	}()

	providerClient, err := NewProviderClient(*serverAddr, logger)
	if err != nil {
		logger.Error("failed to create provider client", "error", err)
		return
	}
	defer providerClient.Close()

	err = providerClient.CheckAndStake()
	if err != nil {
		logger.Error("failed to check and stake", "error", err)
		return
	}

	bidS, err := providerClient.ReceiveBids()
	if err != nil {
		logger.Error("failed to create bid receiver", "error", err)
		return
	}

	for bid := range bidS {
		receivedBids.Inc()
		buf, err := json.Marshal(bid)
		if err != nil {
			logger.Error("failed to marshal bid", "error", err)
		}
		logger.Info("received new bid", "bid", string(buf))

		status := providerapiv1.BidResponse_STATUS_ACCEPTED
		if *errorProbability > 0 {
			if rand.Intn(100) < *errorProbability {
				logger.Warn("sending error response")
				status = providerapiv1.BidResponse_STATUS_REJECTED
				rejectedBids.Inc()
			}
		}
		err = providerClient.SendBidResponse(context.Background(), &providerapiv1.BidResponse{
			BidDigest:         bid.BidDigest,
			Status:            status,
			DispatchTimestamp: time.Now().UnixMilli() + 600,
		})
		if err != nil {
			logger.Error("failed to send bid response", "error", err)
			return
		}
		sentBids.Inc()
		logger.Info("sent bid", "status", status.String())
	}
}
