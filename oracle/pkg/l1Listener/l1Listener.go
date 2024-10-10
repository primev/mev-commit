package l1Listener

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"math/big"
	"net/http"
	"sort"
	"strings"
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
	RecordL1Block(blockNum *big.Int, winner string) (*types.Transaction, error)
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
	relayUrls      []string
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
) *L1Listener {
	return &L1Listener{
		logger:         logger,
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		eventMgr:       evtMgr,
		recorder:       recorder,
		metrics:        newMetrics(),
		relayUrls: []string{
			"https://boost-relay.flashbots.net/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://bloxroute.max-profit.blxrbdn.com/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://bloxroute.regulated.blxrbdn.com/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://relay.edennetwork.io/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://mainnet-relay.securerpc.com/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://relay.ultrasound.money/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://agnostic-relay.net/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://aestus.live/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://relay.wenmerge.com/relay/v1/data/bidtraces/proposer_payload_delivered",
			"https://blockspace.frontier.tech/relay/v1/data/bidtraces/proposer_payload_delivered",
		},
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

	eg.Go(func() error {
		return l.watchRelays(egCtx)
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

				winner := string(bytes.ToValidUTF8(header.Extra, []byte("�")))

				l.logger.Info(
					"new L1 winner",
					"winner", winner,
					"block", header.Number.Int64(),
				)

				winnerPostingTxn, err := l.recorder.RecordL1Block(
					big.NewInt(0).SetUint64(b),
					winner,
				)
				if err != nil {
					l.logger.Error("failed to register winner for block", "block", b, "error", err)
					continue
				}

				l.metrics.WinnerPostedCount.Inc()
				l.metrics.LastSentNonce.Set(float64(winnerPostingTxn.Nonce()))

				l.logger.Info(
					"registered winner",
					"winner", winner,
					"block", header.Number.Int64(),
					"txn", winnerPostingTxn.Hash().String(),
				)
			}

			currentBlockNo = int64(blockNum)
		}
	}
}

func (l *L1Listener) watchRelays(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastBlockNumber int64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			relayData, err := l.fetchRelayData()
			if err != nil {
				l.logger.Error("failed to fetch relay data", "error", err)
				continue
			}

			for _, data := range relayData {
				if data.BlockNumber > lastBlockNumber {
					l.logger.Info("New relay data",
						"timestamp", data.Timestamp,
						"block_number", data.BlockNumber,
						"relay_name", data.RelayName,
						"builder_pubkey", data.BuilderPubkey,
						"block_hash", data.BlockHash,
						"slot", data.Slot,
					)
					lastBlockNumber = data.BlockNumber
				}
			}
		}
	}
}

func (l *L1Listener) fetchRelayData() ([]RelayData, error) {
	var allData []RelayData

	for _, url := range l.relayUrls {
		resp, err := http.Get(url)
		if err != nil {
			l.logger.Error("failed to fetch data from relay", "url", url, "error", err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.logger.Error("failed to read response body", "url", url, "error", err)
			continue
		}

		var data []map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			l.logger.Error("failed to unmarshal response", "url", url, "error", err)
			continue
		}

		if len(data) > 0 {
			latestBid := data[0]
			relayData := RelayData{
				RelayName:     url[8:strings.Index(url, "/relay")],
				BuilderPubkey: fmt.Sprintf("%.10s...", latestBid["builder_pubkey"]),
				BlockNumber:   int64(latestBid["block_number"].(float64)),
				BlockHash:     fmt.Sprintf("%.10s...", latestBid["block_hash"]),
				Slot:          fmt.Sprintf("%v", latestBid["slot"]),
				Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
			}
			allData = append(allData, relayData)
		}
	}

	sort.Slice(allData, func(i, j int) bool {
		return allData[i].BlockNumber > allData[j].BlockNumber
	})

	return allData, nil
}
