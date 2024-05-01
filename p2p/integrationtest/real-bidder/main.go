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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	pb "github.com/primevprotocol/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	logLevel = flag.String("log-level", "debug", "Verbosity level (debug|info|warn|error)")
	httpPort = flag.Int("http-port", 8080, "The port to serve the HTTP metrics endpoint on")
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

var deposits = map[uint64]struct{}{}

func main() {
	flag.Parse()
	if *serverAddr == "" {
		fmt.Println("Please provide a valid server address with the -serverAddr flag")
		return
	}

	level := new(slog.LevelVar)
	if err := level.UnmarshalText([]byte(*logLevel)); err != nil {
		level.Set(slog.LevelDebug)
		fmt.Printf("Invalid log level: %s; using %q", err, level)
	}

	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: level},
	))

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

	conn, err := grpc.Dial(
		*serverAddr,
		grpc.WithTransportCredentials(credentials.NewTLS(
			// Integration tests take place in a controlled environment,
			// thus we do not expect machine-in-the-middle attacks.
			&tls.Config{InsecureSkipVerify: true},
		)),
	)
	if err != nil {
		logger.Error("failed to connect to server", "err", err)
		return
	}
	defer conn.Close()

	bidderClient := pb.NewBidderClient(conn)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			err = checkOrDeposit(bidderClient, logger)
			if err != nil {
				logger.Error("failed to check or stake", "err", err)
			}
			<-ticker.C
		}
	}()

	type blockWithTxns struct {
		blockNum int64
		txns     []string
	}

	newBlockChan := make(chan blockWithTxns, 1)

	wg.Add(1)
	go func(logger *slog.Logger) {
		defer wg.Done()

		currentBlkNum := int64(0)
		for {
			block, blkNum, err := RetreivedBlock(rpcClient)
			if err != nil || len(block) == 0 {
				logger.Error("failed to get block", "err", err)
			}

			if currentBlkNum == blkNum {
				time.Sleep(1 * time.Second)
				continue
			}

			currentBlkNum = blkNum
			newBlockChan <- blockWithTxns{
				blockNum: blkNum,
				txns:     block,
			}
		}
	}(logger)

	wg.Add(1)
	go func(logger *slog.Logger) {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond)
		currentBlock := blockWithTxns{}
		for {
			select {
			case block := <-newBlockChan:
				currentBlock = block
			case <-ticker.C:
			}

			if len(currentBlock.txns) == 0 {
				continue
			}

			bundleLen := rand.Intn(10)
			bundleStart := rand.Intn(len(currentBlock.txns))
			bundleEnd := bundleStart + bundleLen
			if bundleEnd > len(currentBlock.txns) {
				bundleEnd = len(currentBlock.txns) - 1
			}

			min := 1000
			max := 10000
			startTimeDiff := rand.Intn(max-min+1) + min
			endTimeDiff := rand.Intn(max-min+1) + min
			err = sendBid(
				bidderClient,
				logger,
				rpcClient,
				currentBlock.txns[bundleStart:bundleEnd],
				currentBlock.blockNum,
				(time.Now().UnixMilli())-int64(startTimeDiff),
				(time.Now().UnixMilli())+int64(endTimeDiff),
			)
			if err != nil {
				logger.Error("failed to send bid", "err", err)
			}
		}
	}(logger)

	wg.Wait()
}

func RetreivedBlock(rpcClient *ethclient.Client) ([]string, int64, error) {
	blkNum, err := rpcClient.BlockNumber(context.Background())
	if err != nil {
		return nil, -1, err
	}
	fullBlock, err := rpcClient.BlockByNumber(context.Background(), big.NewInt(int64(blkNum)))
	if err != nil {
		return nil, -1, err
	}

	blockTxns := []string{}
	txns := fullBlock.Transactions()
	for _, txn := range txns {
		blockTxns = append(blockTxns, txn.Hash().Hex()[2:])
	}

	return blockTxns, int64(blkNum), nil
}

func checkOrDeposit(
	bidderClient pb.BidderClient,
	logger *slog.Logger,
) error {
	deposit, err := bidderClient.GetDeposit(context.Background(), &pb.GetDepositRequest{})
	if err != nil {
		logger.Error("failed to get deposit", "err", err)
		return err
	}

	logger.Info("deposit", "amount", deposit.Amount)

	minDeposit, err := bidderClient.GetMinDeposit(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		logger.Error("failed to get min deposit", "err", err)
		return err
	}

	depositAmt, set := big.NewInt(0).SetString(deposit.Amount, 10)
	if !set {
		logger.Error("failed to parse deposit amount")
		return errors.New("failed to parse deposit amount")
	}

	minDepositAmt, set := big.NewInt(0).SetString(minDeposit.Amount, 10)
	if !set {
		logger.Error("failed to parse min deposit amount")
		return errors.New("failed to parse min deposit amount")
	}

	if depositAmt.Cmp(minDepositAmt) > 0 {
		logger.Error("bidder already has balance")
		return nil
	}

	topup := big.NewInt(0).Mul(minDepositAmt, big.NewInt(10))

	deposit, err = bidderClient.Deposit(context.Background(), &pb.DepositRequest{
		Amount: topup.String(),
	})
	if err != nil {
		logger.Error("failed to deposit", "err", err)
		return err
	}

	logger.Info("deposit", "amount", topup.String())

	deposits[deposit.WindowNumber.Value] = struct{}{}

	for window := range deposits {
		if window < deposit.WindowNumber.Value-2 {
			resp, err := bidderClient.Withdraw(context.Background(), &pb.WithdrawRequest{
				WindowNumber: wrapperspb.UInt64(window),
			})
			if err != nil {
				logger.Error("failed to withdraw", "err", err)
				return err
			}
			logger.Info("withdraw", "amount", resp.Amount, "window", resp.WindowNumber)
			delete(deposits, window)
		}
	}

	return nil
}

func sendBid(
	bidderClient pb.BidderClient,
	logger *slog.Logger,
	rpcClient *ethclient.Client,
	txnHashes []string,
	blkNum int64,
	decayStartTimestamp int64,
	decayEndTimestamp int64,
) error {
	amount := rand.Intn(200000)
	amount += 100000

	bid := &pb.Bid{
		TxHashes:            txnHashes,
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

	sendBidDuration.WithLabelValues(
		"success",
		fmt.Sprintf("%d", preConfCount),
	).Observe(time.Since(start).Seconds())
	return nil
}
