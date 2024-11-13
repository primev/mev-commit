package l1Listener

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
			// We need to get the previous block number because the current block has finalized header
			blockNum = blockNum - 1

			if blockNum <= uint64(currentBlockNo) {
				continue
			}

			for b := uint64(currentBlockNo) + 1; b <= blockNum; b++ {
				header, err := l.l1Client.HeaderByNumber(ctx, new(big.Int).SetUint64(b+1))
				if err != nil {
					l.logger.Error("failed to get header", "block", b, "error", err)
					continue
				}
				// End of changes needed to be done.
				var builderPubKey string
				l.logger.Info("querying relay", "block", b, "hash", header.ParentHash.Hex())
				builderPubKey, err = l.relayQuerier.Query(ctx, int64(b), header.ParentHash.Hex())
				if err != nil {
					l.logger.Info("block not found in relay, assuming out of PBS block", "block", b, "error", err)
					builderPubKey = "" // Set a default value in case of failure
				}

				builderPubKey = strings.TrimPrefix(builderPubKey, "0x")

				builderPubKeyBytes, err := hex.DecodeString(builderPubKey)
				if err != nil {
					l.logger.Error("failed to decode builder pubkey", "block", b, "builder_pubkey", builderPubKey, "error", err)
				}

				l.logger.Info(
					"new L1 winner",
					"block", header.Number.Int64()-1,
					"builder_pubkey", builderPubKey,
				)

				winnerPostingTxn, err := l.recorder.RecordL1Block(
					big.NewInt(0).SetUint64(b),
					builderPubKeyBytes,
				)
				if err != nil {
					l.logger.Error("failed to register winner for block", "block", b, "error", err)
					continue
				}

				l.metrics.WinnerPostedCount.Inc()
				l.metrics.LastSentNonce.Set(float64(winnerPostingTxn.Nonce()))

				l.logger.Info(
					"registered winner",
					"block", header.Number.Int64()-1,
					"txn", winnerPostingTxn.Hash().String(),
				)
			}

			currentBlockNo = int64(blockNum)
		}
	}
}

type RelayQuerier interface {
	Query(ctx context.Context, blockNumber int64, blockHash string) (string, error)
}

type RelayQueryEngine struct {
	relayUrls []string
	logger    *slog.Logger
}

func NewRelayQueryEngine(relayUrls []string, logger *slog.Logger) RelayQuerier {
	return &RelayQueryEngine{
		relayUrls: relayUrls,
		logger:    logger,
	}
}

func (m *RelayQueryEngine) Query(ctx context.Context, blockNumber int64, blockHash string) (string, error) {
	var wg sync.WaitGroup
	resultChan := make(chan string, len(m.relayUrls))

	for _, u := range m.relayUrls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			path := url.PathEscape(fmt.Sprintf("/relay/v1/data/bidtraces/proposer_payload_delivered?block_number=%d", blockNumber))
			fullUrl := fmt.Sprintf("%s%s", u, path)
			m.logger.Debug("querying relay", "url", fullUrl)

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
			if err != nil {
				m.logger.Error("failed to create request", "url", fullUrl, "error", err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				m.logger.Error("failed to fetch data from relay", "url", fullUrl, "error", err)
				return
			}
			defer resp.Body.Close()
			m.logger.Info("received response from relay", "url", fullUrl, "status", resp.Status)

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
				blockNumInt, err := strconv.Atoi(blockNum)
				if err != nil {
					m.logger.Error("failed to convert block_number to int", "block_number", blockNum, "error", err)
					continue
				}
				if blockNumInt == int(blockNumber) && item["block_hash"] == blockHash {
					resultChan <- fmt.Sprintf("%v", item["builder_pubkey"])
					return
				}
			}
		}(u)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result, ok := <-resultChan:
		if !ok {
			return "", errors.New("no matching block found")
		}
		return result, nil
	}
}
