package l1Listener

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/oracle/pkg/store"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var checkInterval = 2 * time.Second

type L1Recorder interface {
	RecordL1Block(blockNum *big.Int, winner []byte) (*types.Transaction, error)
}

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner []byte, window int64) error
	LastWinnerBlock() (int64, error)
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type L1Listener struct {
	logger         *slog.Logger
	l1Client       EthClient
	winnerRegister WinnerRegister
	eventMgr       events.EventManager
	recorder       L1Recorder
	metrics        *metrics
	relayQuerier   RelayQuerier
	builderData    map[int64]string
}

type RelayData struct {
	RelayName     string
	BuilderPubkey string
	BlockNumber   int64
	BlockHash     string
	Slot          string
	Timestamp     string
}

func NewL1Listener(
	logger *slog.Logger,
	l1Client EthClient,
	winnerRegister WinnerRegister,
	evtMgr events.EventManager,
	recorder L1Recorder,
	relayQuerier RelayQuerier,
) *L1Listener {
	return &L1Listener{
		logger:         logger,
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		eventMgr:       evtMgr,
		recorder:       recorder,
		metrics:        newMetrics(),
		relayQuerier:   relayQuerier,
		builderData:    make(map[int64]string),
	}
}

func (l *L1Listener) Metrics() []prometheus.Collector {
	return l.metrics.Collectors()
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return l.watchL1Block(egCtx)
	})

	evt := events.NewEventHandler(
		"NewL1Block",
		func(update *blocktracker.BlocktrackerNewL1Block) {
			l.logger.Info(
				"new L1 block event",
				"block", update.BlockNumber,
				"winner", update.Winner.String(),
				"window", update.Window,
			)
			err := l.winnerRegister.RegisterWinner(
				ctx,
				update.BlockNumber.Int64(),
				update.Winner.Bytes(),
				update.Window.Int64(),
			)
			if err != nil {
				l.logger.Error(
					"failed to register winner",
					"block", update.BlockNumber,
					"winner", update.Winner.String(),
					"error", err,
				)
				return
			}
			l.metrics.WinnerCount.Inc()
			l.metrics.WinnerRoundCount.WithLabelValues(update.Winner.String()).Inc()
		},
	)

	sub, err := l.eventMgr.Subscribe(evt)
	if err != nil {
		close(doneChan)
		return doneChan
	}

	eg.Go(func() error {
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			return egCtx.Err()
		case err := <-sub.Err():
			return err
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			l.logger.Error("L1listener error", "error", err)
		}
	}()

	l.logger.Info("L1Listener started")

	return doneChan
}

func (l *L1Listener) watchL1Block(ctx context.Context) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	currentBlockNo, err := l.winnerRegister.LastWinnerBlock()
	if err != nil {
		// this is a fresh start, so start from the current block
		if errors.Is(err, store.ErrNotFound) {
			tip, err := l.l1Client.BlockNumber(ctx)
			if err != nil {
				l.logger.Error("failed to get current block number", "error", err)
				return err
			}
			currentBlockNo = int64(tip)
		} else {
			l.logger.Error("failed to get last winner block", "error", err)
			return err
		}
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			blockNum, err := l.l1Client.BlockNumber(ctx)
			if err != nil {
				l.logger.Error("failed to get block number", "error", err)
				continue
			}

			if blockNum <= uint64(currentBlockNo) {
				continue
			}

			for b := uint64(currentBlockNo) + 1; b <= blockNum; b++ {
				header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(b)))
				if err != nil {
					l.logger.Error("failed to get header", "block", b, "error", err)
					continue
				}

				winnerExtraData := string(bytes.ToValidUTF8(header.Extra, []byte("")))

				// End of changes needed to be done.
				var builderPubKey string
				startTime := time.Now()
				for time.Since(startTime) < 2*time.Second {
					builderPubKey, err = l.relayQuerier.Query(int64(b), header.Hash().String())
					if err == nil {
						break
					}
					l.logger.Warn("failed to query relay, retrying", "block", b, "error", err)
					time.Sleep(500 * time.Millisecond)
				}

				if err != nil {
					l.logger.Error("failed to query relay after retries", "block", b, "error", err)
					builderPubKey = "" // Set a default value in case of failure
				}

				l.logger.Info(
					"new L1 winner",
					"winner_extra_data", winnerExtraData,
					"block", header.Number.Int64(),
					"builder_pubkey", builderPubKey,
				)

				winnerPostingTxn, err := l.recorder.RecordL1Block(
					big.NewInt(0).SetUint64(b),
					[]byte(builderPubKey),
				)
				if err != nil {
					l.logger.Error("failed to register winner for block", "block", b, "error", err)
					continue
				}

				l.metrics.WinnerPostedCount.Inc()
				l.metrics.LastSentNonce.Set(float64(winnerPostingTxn.Nonce()))

				l.logger.Info(
					"registered winner",
					"winner_extra_data", winnerExtraData,
					"block", header.Number.Int64(),
					"txn", winnerPostingTxn.Hash().String(),
				)
			}

			currentBlockNo = int64(blockNum)
		}
	}
}

type RelayQuerier interface {
	Query(blockNumber int64, blockHash string) (string, error)
}

type MiniRelayQueryEngine struct {
	relayUrls []string
	logger    *slog.Logger
}

func NewMiniRelayQueryEngine(relayUrls []string, logger *slog.Logger) RelayQuerier {
	return &MiniRelayQueryEngine{
		relayUrls: relayUrls,
		logger:    logger,
	}
}

func (m *MiniRelayQueryEngine) Query(blockNumber int64, blockHash string) (string, error) {
	var wg sync.WaitGroup
	resultChan := make(chan string, len(m.relayUrls))

	for _, url := range m.relayUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fullUrl := fmt.Sprintf("%s?block_number=%d", url, blockNumber)
			resp, err := http.Get(fullUrl)
			if err != nil {
				m.logger.Error("failed to fetch data from relay", "url", fullUrl, "error", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				m.logger.Error("failed to read response body", "url", fullUrl, "error", err)
				return
			}

			var data []map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				m.logger.Error("failed to unmarshal response", "url", fullUrl, "error", err)
				return
			}

			for _, item := range data {
				blockNum, ok := item["block_number"].(string)
				if !ok {
					m.logger.Error("block_number is not a string", "block_number", item["block_number"])
					continue
				}

				if blockNum == fmt.Sprintf("%d", blockNumber) && item["block_hash"] == blockHash {
					resultChan <- fmt.Sprintf("%v", item["builder_pubkey"])
					return
				}
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		return result, nil
	}

	return "", errors.New("no matching block found")
}
