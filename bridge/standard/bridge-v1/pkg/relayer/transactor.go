package relayer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/shared"
)

type Transactor struct {
	logger            *slog.Logger
	privateKey        *ecdsa.PrivateKey
	rawClient         *shared.ETHClient
	gatewayTransactor shared.GatewayTransactor
	gatewayFilterer   shared.GatewayFilterer
	chainID           *big.Int
	chain             shared.Chain
	eventChan         <-chan shared.TransferInitiatedEvent
	mostRecentFinalized
}

type mostRecentFinalized struct {
	event shared.TransferFinalizedEvent
	opts  bind.FilterOpts
}

func (m mostRecentFinalized) String() string {
	return "Event: " + m.event.String() + ". Opts start: " + fmt.Sprint(m.opts.Start) + ". Opts end: " + fmt.Sprint(*m.opts.End)
}

func NewTransactor(
	logger *slog.Logger,
	pk *ecdsa.PrivateKey,
	gatewayAddr common.Address,
	ethClient *ethclient.Client,
	gatewayTransactor shared.GatewayTransactor,
	gatewayFilterer shared.GatewayFilterer,
	eventChan <-chan shared.TransferInitiatedEvent,
) *Transactor {
	return &Transactor{
		logger:     logger,
		privateKey: pk,
		rawClient: shared.NewETHClient(
			logger.With("component", "eth_client"),
			ethClient,
		),
		gatewayTransactor: gatewayTransactor,
		gatewayFilterer:   gatewayFilterer,
		eventChan:         eventChan,
		mostRecentFinalized: mostRecentFinalized{
			event: shared.TransferFinalizedEvent{},
			opts:  bind.FilterOpts{Start: 0, End: nil}, // TODO: cache doesn't need to start at 0 once non-syncing relayer is implemented
		},
	}
}

func (t *Transactor) Start(ctx context.Context) (<-chan struct{}, error) {
	var err error
	t.chainID, err = t.rawClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %w", err)
	}
	switch t.chainID.String() {
	case "39999":
		t.logger.Info("starting transactor for local_l1")
		t.chain = shared.L1
	case "17000":
		t.logger.Info("starting transactor for Holesky L1")
		t.chain = shared.L1
	case "17864":
		t.logger.Info("starting transactor for mev-commit chain (settlement)")
		t.chain = shared.Settlement
	default:
		return nil, fmt.Errorf("unsupported chain id: %s", t.chainID)
	}

	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		if err := t.rawClient.CancelPendingTxes(ctx, t.privateKey); err != nil {
			t.logger.Error("failed to cancel pending transactions", "error", err)
		}

		for event := range t.eventChan {
			t.logger.Debug(
				"received signal from listener to submit transfer finalization tx",
				"dst_chain", t.chain,
				"src_chain", event.Chain,
				"recipient", event.Recipient,
				"amount", event.Amount,
				"src_transfer_idx", event.TransferIdx,
			)
			opts, err := t.rawClient.CreateTransactOpts(ctx, t.privateKey, t.chainID)
			if err != nil {
				t.logger.Error("failed to create transact opts for transfer finalization tx", "error", err)
				t.logger.Warn("skipping transfer finalization tx", "src_transfer_idx", event.TransferIdx)
				continue
			}
			finalized, err := t.transferAlreadyFinalized(ctx, event.TransferIdx)
			if err != nil {
				t.logger.Error("failed to check if transfer already finalized")
				t.logger.Warn("skipping transfer finalization tx", "src_transfer_idx", event.TransferIdx)
				continue
			}
			if finalized {
				continue
			}
			receipt, err := t.sendFinalizeTransfer(ctx, opts, event)
			if err != nil {
				t.logger.Error("failed to send transfer finalization tx")
				t.logger.Warn("skipping transfer finalization tx", "src_transfer_idx", event.TransferIdx)
				continue
			}
			// Event should be obtainable to update cache
			eventBlock := receipt.BlockNumber.Uint64()
			filterOpts := &bind.FilterOpts{Start: eventBlock, End: &eventBlock, Context: ctx}
			_, found, err := t.obtainTransferFinalizedAndUpdateCache(ctx, filterOpts, event.TransferIdx)
			if err != nil {
				t.logger.Error("failed to obtain transfer finalized event after sending tx")
				continue
			}
			if !found {
				t.logger.Warn("transfer finalized event not found after sending tx")
				continue
			}
		}
		t.logger.Info("channel to transactor was closed, transactor is exiting", "chain", t.chain)
	}()
	return doneChan, nil
}

func (t *Transactor) transferAlreadyFinalized(
	ctx context.Context,
	transferIdx *big.Int,
) (bool, error) {
	const maxBlockRange = 40000
	var interEndBlock uint64 // Intermediate end block for each range

	currentBlock, err := t.rawClient.BlockNumber(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get current block number: %w", err)
	}

	startBlock := t.mostRecentFinalized.opts.Start

	for start := startBlock; start <= currentBlock; start = interEndBlock + 1 {
		interEndBlock = start + maxBlockRange
		if interEndBlock > currentBlock {
			interEndBlock = currentBlock
		}
		opts := &bind.FilterOpts{
			Start:   start,
			End:     &interEndBlock,
			Context: ctx,
		}
		event, found, err := t.obtainTransferFinalizedAndUpdateCache(ctx, opts, transferIdx)
		if err != nil {
			return false, fmt.Errorf("failed to obtain transfer finalized event: %w", err)
		}
		if found {
			t.logger.Debug(
				"transfer already finalized",
				"dst_chain", t.chain,
				"recipient", event.Recipient,
				"amount", event.Amount,
				"src_transfer_idx", event.CounterpartyIdx,
			)
			return true, nil
		}
	}
	return false, nil
}

func (t *Transactor) sendFinalizeTransfer(
	ctx context.Context,
	opts *bind.TransactOpts,
	event shared.TransferInitiatedEvent,
) (*gethtypes.Receipt, error) {

	// Capture event params in closure and define tx submission callback
	submitFinalizeTransfer := func(
		ctx context.Context,
		opts *bind.TransactOpts,
	) (*gethtypes.Transaction, error) {
		tx, err := t.gatewayTransactor.FinalizeTransfer(opts, event.Recipient, event.Amount, event.TransferIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to send finalize transfer tx: %w", err)
		}
		t.logger.Debug("transfer finalization tx sent, hash: %s, destChain: %s, recipient: %s, amount: %d, srcTransferIdx: %d",
			"hash", tx.Hash().Hex(),
			"dest_chain", t.chain,
			"recipient", event.Recipient,
			"amount", event.Amount,
			"src_transfer_idx", event.TransferIdx,
		)
		return tx, nil
	}

	receipt, err := t.rawClient.WaitMinedWithRetry(ctx, opts, submitFinalizeTransfer)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for finalize transfer tx to be mined: %w", err)
	}
	includedInBlock := receipt.BlockNumber.Uint64()
	t.logger.Info("finalizeTransfer tx included in block", "block_number", includedInBlock, "chain", t.chain)

	return receipt, nil
}

func (t *Transactor) obtainTransferFinalizedAndUpdateCache(
	ctx context.Context,
	opts *bind.FilterOpts,
	transferIdx *big.Int,
) (shared.TransferFinalizedEvent, bool, error) {
	event, found, err := t.gatewayFilterer.ObtainTransferFinalizedEvent(opts, transferIdx)
	if err != nil {
		return shared.TransferFinalizedEvent{}, false, fmt.Errorf("failed to obtain transfer finalized event: %w", err)
	}
	if found {
		t.mostRecentFinalized = mostRecentFinalized{event, *opts}
		t.logger.Debug("mostRecentFinalized cache updated", "new", fmt.Sprintf("%+v", t.mostRecentFinalized))
	}
	return event, found, nil
}
