package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-logr/logr"
	pb "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/x/util"
	"github.com/primev/mev-commit/x/util/otelutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	serverAddr = flag.String(
		"server-addr",
		"localhost:13524",
		"The server address in the format of host:port",
	)
	rpcAddr = flag.String(
		"rpc-addr",
		"localhost:13524",
		"The server address in the format of host:port",
	)
	logLevel = flag.String(
		"log-level",
		"debug",
		"Verbosity level (debug|info|warn|error)",
	)
	logFmt = flag.String(
		"log-fmt",
		"text",
		"Format of the log output: 'text', 'json'",
	)
	logTags = flag.String(
		"log-tags",
		"",
		"Comma-separated list of <name:value> pairs that will be inserted into each log line",
	)
	otelCollectorEndpointURL = flag.String(
		"otel-collector-endpoint-url",
		"",
		"URL for OpenTelemetry collector endpoint",
	)
	httpPort = flag.Int(
		"http-port",
		8080,
		"The port to serve the HTTP metrics endpoint on",
	)
	bidWorkers = flag.Int(
		"bid-workers",
		2,
		"Number of workers to send bids",
	)
)

var (
	receivedPreconfs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mev_commit",
		Subsystem: "bidder_emulator",
		Name:      "total_received_preconfs",
		Help:      "Total number of preconfs received from mev_commit nodes",
	})
	sentBids = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "mev_commit",
		Subsystem: "bidder_emulator",
		Name:      "total_sent_bids",
		Help:      "Total number of bids sent to mev_commit nodes",
	})
	sendBidDuration = *prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "mev_commit",
			Subsystem: "bidder_emulator",
			Name:      "send_bid_duration",
			Help:      "Duration of method calls.",
		},
		[]string{"status", "no_of_preconfs"},
	)
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
		fmt.Println("please provide a valid server address with the -serverAddr flag")
		return
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(receivedPreconfs, sentBids, sendBidDuration)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", "err", err)
		}
	}()

	rpcClient, err := ethclient.Dial(*rpcAddr)
	if err != nil {
		logger.Error("failed to connect to rpc", "err", err)
		return
	}

	// nolint:staticcheck
	conn, err := grpc.Dial(
		*serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(
			// Integration tests take place in a controlled environment,
			// thus we do not expect machine-in-the-middle attacks.
			&tls.Config{InsecureSkipVerify: true},
		)),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logger.Error("failed to connect to server", "err", err)
		return
	}
	defer conn.Close()

	bidderClient := pb.NewBidderClient(conn)

	wg := sync.WaitGroup{}

	// set as 1 eth
	minDeposit, set := big.NewInt(0).SetString("1000000000000000000", 10)
	if !set {
		logger.Error("failed to parse min deposit amount")
		return
	}

	status, err := bidderClient.AutoDepositStatus(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		logger.Error("failed to get auto deposit status", "err", err)
		return
	}

	if !status.IsAutodepositEnabled {
		resp, err := bidderClient.AutoDeposit(context.Background(), &pb.DepositRequest{
			Amount: minDeposit.String(),
		})
		if err != nil {
			logger.Error("failed to auto deposit", "err", err)
			return
		}
		logger.Info("auto deposit", "amount", resp.AmountPerWindow, "window", resp.StartWindowNumber)
	}

	type blockWithTxns struct {
		blockNum int64
		txns     []string
	}

	blockChans := make([]chan *blockWithTxns, *bidWorkers)
	for i := 0; i < *bidWorkers; i++ {
		blockChans[i] = make(chan *blockWithTxns, 1)
	}

	wg.Add(1)
	go func(logger *slog.Logger) {
		defer wg.Done()

		currentBlkNum := uint64(0)
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			blkNum, err := rpcClient.BlockNumber(context.Background())
			if err != nil {
				logger.Error("failed to get block number", "err", err)
				continue
			}

			if blkNum <= currentBlkNum {
				continue
			}

			block, err := RetrieveTxns(rpcClient, blkNum)
			if err != nil {
				logger.Error("failed to get block", "err", err)
				continue
			}

			currentBlkNum = blkNum
			for _, ch := range blockChans {
				txns := make([]string, len(block))
				copy(txns, block)

				ch <- &blockWithTxns{
					blockNum: int64(blkNum),
					txns:     txns,
				}
			}
		}
	}(logger)

	for i := 0; i < *bidWorkers; i++ {
		wg.Add(1)
		go func(logger *slog.Logger, newBlockChan <-chan *blockWithTxns) {
			defer wg.Done()
			ticker := time.NewTicker(200 * time.Millisecond)
			currentBlock := &blockWithTxns{}
			for {
				select {
				case block := <-newBlockChan:
					if block.blockNum <= currentBlock.blockNum {
						continue
					}
					currentBlock = block
				case <-ticker.C:
				}

				if len(currentBlock.txns) == 0 {
					continue
				}

				bundleLen := rand.Intn(10) + 1
				bundleStart := rand.Intn(len(currentBlock.txns))
				bundleEnd := bundleStart + bundleLen
				if bundleEnd >= len(currentBlock.txns) {
					bundleEnd = len(currentBlock.txns)
				}

				min := 5000
				max := 10000
				startTimeDiff := rand.Intn(max-min+1) + min
				endTimeDiff := rand.Intn(max-min+1) + min
				err = sendBid(
					bidderClient,
					logger,
					currentBlock.txns[bundleStart:bundleEnd],
					currentBlock.blockNum,
					(time.Now().UnixMilli())-int64(startTimeDiff),
					(time.Now().UnixMilli())+int64(endTimeDiff),
				)
				if err != nil {
					logger.Error("failed to send bid", "err", err)
				}
			}
		}(logger, blockChans[i])
	}

	wg.Wait()
}

func RetrieveTxns(rpcClient *ethclient.Client, blkNum uint64) ([]string, error) {
	fullBlock, err := rpcClient.BlockByNumber(context.Background(), big.NewInt(int64(blkNum)))
	if err != nil {
		return nil, err
	}

	blockTxns := []string{}
	txns := fullBlock.Transactions()
	for _, txn := range txns {
		blockTxns = append(blockTxns, strings.TrimPrefix(txn.Hash().Hex(), "0x"))
	}

	if len(blockTxns) == 0 {
		return nil, errors.New("no txns in block")
	}

	return blockTxns, nil
}

func sendBid(
	bidderClient pb.BidderClient,
	logger *slog.Logger,
	txnHashes []string,
	blkNum int64,
	decayStartTimestamp int64,
	decayEndTimestamp int64,
) error {
	if len(txnHashes) == 0 {
		return errors.New("no txns to send")
	}
	amount := rand.Intn(200000)
	amount += 100000

	hashesToSend := make([]string, len(txnHashes))
	copy(hashesToSend, txnHashes)

	bid := &pb.Bid{
		TxHashes:            hashesToSend,
		Amount:              strconv.Itoa(amount),
		BlockNumber:         int64(blkNum),
		DecayStartTimestamp: decayStartTimestamp,
		DecayEndTimestamp:   decayEndTimestamp,
	}

	logger.Info("sending bid", "bid", bid)

	start := time.Now()
	rcv, err := bidderClient.SendBid(context.Background(), bid)
	if err != nil {
		logger.Error("failed to send bid", "err", err)
		return err
	}

	sentBids.Inc()

	ctx, span := otel.Tracer("").Start(rcv.Context(), "receive")
	logger.InfoContext(ctx, "receiving preconfs")
	preConfCount := 0
	for {
		_, err := rcv.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error("failed receiving preconf", "error", err)
			sendBidDuration.WithLabelValues(
				"error",
				fmt.Sprintf("%d", preConfCount),
			).Observe(time.Since(start).Seconds())
			return err
		}
		receivedPreconfs.Inc()
		preConfCount++
	}
	logger.InfoContext(ctx, "received preconfs")
	span.End()

	sendBidDuration.WithLabelValues(
		"success",
		fmt.Sprintf("%d", preConfCount),
	).Observe(time.Since(start).Seconds())
	return nil
}
