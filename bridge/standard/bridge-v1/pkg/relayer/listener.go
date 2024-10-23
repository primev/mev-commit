package relayer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/shared"
)

type Listener struct {
	logger          *slog.Logger
	rawClient       *ethclient.Client
	gatewayFilterer shared.GatewayFilterer
	sync            bool
	chain           shared.Chain
	DoneChan        chan struct{}
	EventChan       chan shared.TransferInitiatedEvent
}

func NewListener(
	logger *slog.Logger,
	client *ethclient.Client,
	gatewayFilterer shared.GatewayFilterer,
	sync bool,
) *Listener {
	return &Listener{
		logger:          logger,
		rawClient:       client,
		gatewayFilterer: gatewayFilterer,
		sync:            true,
	}
}

func (l *Listener) Start(ctx context.Context) (
	<-chan struct{}, <-chan shared.TransferInitiatedEvent, error,
) {
	chainID, err := l.rawClient.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain id: %w", err)
	}
	switch chainID.String() {
	case "39999":
		l.logger.Info("starting listener for local_l1")
		l.chain = shared.L1
	case "17000":
		l.logger.Info("starting listener for Holesky L1")
		l.chain = shared.L1
	case "17864":
		l.logger.Info("starting listener for mev-commit chain (settlement)")
		l.chain = shared.Settlement
	default:
		return nil, nil, fmt.Errorf("unsupported chain id: %s", chainID)
	}

	l.DoneChan = make(chan struct{})
	l.EventChan = make(chan shared.TransferInitiatedEvent, 10) // Buffer up to 10 events

	go func() {
		defer close(l.DoneChan)
		defer close(l.EventChan)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Blocks up to this value have been handled
		blockNumHandled := uint64(0)

		if l.sync {
			blockNumHandled, err = l.obtainFinalizedBlockNum(ctx)
			if err != nil {
				l.logger.Error("failed to obtain block number during sync", "error", err)
				return
			}
			// Most nodes limit query ranges so we fetch in 40k increments
			events, err := l.obtainTransferInitiatedEventsInBatches(
				ctx, 0, blockNumHandled)
			if err != nil {
				l.logger.Error("failed to fetch transfer initiated events during sync", "error", err)
				return
			}
			for _, event := range events {
				l.logger.Info("transfer initiated event seen by listener during sync", "event", event)
				l.EventChan <- event
			}
		}

		for {
			select {
			case <-ctx.Done():
				l.logger.Info("listener shutting down", "chain", l.chain.String())
				return
			case <-ticker.C:
			}

			currentBlockNum, err := l.obtainFinalizedBlockNum(ctx)
			if err != nil {
				// TODO: Secondary url if rpc fails. For now just start over...
				l.logger.Error("failed to obtain block number", "error", err)
				l.logger.Warn("listener restarting from block 0...")
				blockNumHandled = 0
				continue
			}
			if blockNumHandled < currentBlockNum {
				events, err := l.obtainTransferInitiatedEventsInBatches(ctx, blockNumHandled+1, currentBlockNum)
				if err != nil {
					// TODO: Secondary url if rpc fails. For now just start over...
					l.logger.Error(
						"failed to query transfer initiated events",
						"from_block", blockNumHandled+1,
						"to_block", currentBlockNum,
						"chain", l.chain.String(),
					)
					l.logger.Warn("listener restarting from block 0...")
					blockNumHandled = 0
					continue
				}
				l.logger.Debug(
					"fetched events",
					"event_count", len(events),
					"from_block", blockNumHandled+1,
					"to_block", currentBlockNum,
					"chain", l.chain.String(),
				)
				for _, event := range events {
					l.logger.Info("transfer initiated event seen by listener", "event", event)
					l.EventChan <- event
				}
				blockNumHandled = currentBlockNum
			}
		}
	}()
	return l.DoneChan, l.EventChan, nil
}

func (l *Listener) obtainFinalizedBlockNum(ctx context.Context) (uint64, error) {
	blockNum, err := l.rawClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to obtain block number: %w", err)
	}
	// Blocks 2 epochs old are considered finalized
	epochBlocks := uint64(32)
	if blockNum < 2*epochBlocks {
		return 0, nil
	}
	return blockNum - 2*epochBlocks, nil
}

func (l *Listener) obtainTransferInitiatedEventsInBatches(
	ctx context.Context,
	startBlock,
	endBlock uint64,
) ([]shared.TransferInitiatedEvent, error) {
	var totalEvents []shared.TransferInitiatedEvent
	const maxBlockRange = 40000
	for start := startBlock; start <= endBlock; start += maxBlockRange + 1 {
		end := start + maxBlockRange
		if end > endBlock {
			end = endBlock
		}
		opts := &bind.FilterOpts{Start: start, End: &end, Context: ctx}
		events, err := l.gatewayFilterer.ObtainTransferInitiatedEvents(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to obtain transfer initiated events: %w", err)
		}
		totalEvents = append(totalEvents, events...)
	}
	return totalEvents, nil
}
