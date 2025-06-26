package sender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"golang.org/x/sync/errgroup"
)

type TxType int

const (
	TxTypeRegular TxType = iota
	TxTypeDeposit
	TxTypeInstantBridge
)

const (
	blockTime = 12 // seconds, typical Ethereum block time
)

var (
	ErrInvalidTransaction       = errors.New("invalid transaction")
	ErrUnsupportedTxType        = errors.New("unsupported transaction type")
	ErrEmptyRawTransaction      = errors.New("empty raw transaction")
	ErrEmptyTransactionTo       = errors.New("empty transaction 'to' address")
	ErrNegativeTransactionValue = errors.New("negative transaction value")
	ErrZeroGasLimit             = errors.New("zero gas limit")
)

type Transaction struct {
	*types.Transaction
	Sender common.Address
	Raw    string
	Type   TxType
}

type Store interface {
	AddQueuedTransaction(tx *Transaction) error
	GetQueuedTransactions() ([]*Transaction, error)
	GetCurrentNonce(sender common.Address) uint64
	HasBalance(ctx context.Context, sender common.Address, amount *big.Int) bool
	AddBalance(ctx context.Context, account common.Address, amount *big.Int) error
	DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error
	StorePreconfirmedTransaction(
		ctx context.Context,
		blockNumber int64,
		txn *Transaction,
		commitments []*bidderapiv1.Commitment,
	) error
}

type Bidder interface {
	Estimate() (int64, error)
	Bid(
		ctx context.Context,
		bidAmount *big.Int,
		slashAmount *big.Int,
		rawTx string,
		opts *optinbidder.BidOpts,
	) (chan optinbidder.BidStatus, error)
}

type Pricer interface {
	EstimatePrice(
		ctx context.Context,
		txn *types.Transaction,
	) (*pricer.BlockPrice, error)
}

type BlockTracker interface {
	CheckTxnInclusion(ctx context.Context, txnHash common.Hash, blockNumber uint64) (bool, error)
}

type TxSender struct {
	logger          *slog.Logger
	store           Store
	bidder          Bidder
	pricer          Pricer
	blockTracker    BlockTracker
	eg              *errgroup.Group
	egCtx           context.Context
	trigger         chan struct{}
	workerPool      chan struct{}
	inflightTxns    map[common.Hash]struct{}
	inflightAccount map[common.Address]struct{}
	inflightMu      sync.Mutex
}

func NewTxSender(
	st Store,
	bidder Bidder,
	pricer Pricer,
	blockTracker BlockTracker,
	logger *slog.Logger,
) *TxSender {
	return &TxSender{
		store:           st,
		bidder:          bidder,
		pricer:          pricer,
		blockTracker:    blockTracker,
		logger:          logger.With("component", "TxSender"),
		workerPool:      make(chan struct{}, 512),
		trigger:         make(chan struct{}, 1),
		inflightTxns:    make(map[common.Hash]struct{}),
		inflightAccount: make(map[common.Address]struct{}),
	}
}

func validateTransaction(tx *Transaction) error {
	if tx == nil || tx.Transaction == nil {
		return ErrInvalidTransaction
	}
	if tx.Type < TxTypeRegular || tx.Type > TxTypeInstantBridge {
		return ErrUnsupportedTxType
	}
	if tx.Raw == "" {
		return ErrEmptyRawTransaction
	}
	if tx.Transaction.To() == nil {
		return ErrEmptyTransactionTo
	}
	if tx.Transaction.Value().Sign() < 0 {
		return ErrNegativeTransactionValue
	}
	if tx.Transaction.Gas() == 0 {
		return ErrZeroGasLimit
	}
	return nil
}

func (t *TxSender) hasLowerNonce(tx *Transaction) bool {
	currentNonce := t.store.GetCurrentNonce(tx.Sender)
	return tx.Transaction.Nonce() < currentNonce
}

func (t *TxSender) triggerSender() {
	select {
	case t.trigger <- struct{}{}:
	default:
		// Non-blocking send, if the channel is full, we do nothing
	}
}

func (t *TxSender) Enqueue(tx *Transaction) error {
	if err := validateTransaction(tx); err != nil {
		t.logger.Error("Invalid transaction", "error", err, "transaction", tx.Raw)
		return err
	}

	if t.hasLowerNonce(tx) {
		return errors.New("transaction has a lower nonce than the current highest nonce")
	}

	if err := t.store.AddQueuedTransaction(tx); err != nil {
		return err
	}

	t.triggerSender()

	return nil
}

func (t *TxSender) Start(ctx context.Context) chan struct{} {
	t.eg, t.egCtx = errgroup.WithContext(ctx)
	done := make(chan struct{})

	t.eg.Go(func() error {
		for {
			select {
			case <-t.egCtx.Done():
				t.logger.Info("Context cancelled, stopping TxSender")
				return ctx.Err()
			case <-t.trigger:
				t.processQueuedTransactions(t.egCtx)
			}
		}
	})

	go func() {
		defer close(done)
		if err := t.eg.Wait(); err != nil {
			t.logger.Error("Error in TxSender", "error", err)
			return
		}
	}()

	return done
}

func (t *TxSender) markInflight(txn *Transaction) bool {
	t.inflightMu.Lock()
	defer t.inflightMu.Unlock()

	if _, ok := t.inflightTxns[txn.Hash()]; ok {
		t.logger.Debug("Transaction already in flight, skipping", "hash", txn.Hash().Hex())
		return false
	}
	if _, ok := t.inflightAccount[txn.Sender]; ok {
		t.logger.Debug("Transaction sender already has an inflight transaction, skipping", "sender", txn.Sender.Hex())
		t.triggerSender() // Trigger to reprocess later
		return false
	}

	t.inflightTxns[txn.Hash()] = struct{}{}
	t.inflightAccount[txn.Sender] = struct{}{}
	return true
}

func (t *TxSender) markCompleted(txn *Transaction) {
	t.inflightMu.Lock()
	defer t.inflightMu.Unlock()

	delete(t.inflightTxns, txn.Hash())
	delete(t.inflightAccount, txn.Sender)
}

func (t *TxSender) processQueuedTransactions(ctx context.Context) {
	txns, err := t.store.GetQueuedTransactions()
	if err != nil {
		t.logger.Error("Failed to get queued transactions", "error", err)
		return
	}
	if len(txns) == 0 {
		t.logger.Info("No queued transactions to process")
		return
	}
	t.logger.Info("Processing queued transactions", "count", len(txns))
	for _, txn := range txns {
		txn := txn // capture range variable
		select {
		case <-ctx.Done():
			t.logger.Info("Context cancelled, stopping transaction processing")
			return
		case t.workerPool <- struct{}{}:
			t.eg.Go(func() error {
				defer func() { <-t.workerPool }()
				if !t.markInflight(txn) {
					// Transaction is already being processed or sender has an inflight transaction
					return nil
				}
				defer t.markCompleted(txn)

				t.logger.Info("Processing transaction", "sender", txn.Sender.Hex(), "type", txn.Type)
				if err := t.processTransaction(ctx, txn); err != nil {
					t.logger.Error("Failed to process transaction", "sender", txn.Sender.Hex(), "error", err)
				}
				return nil
			})
		}
	}
}

func (t *TxSender) processTransaction(ctx context.Context, txn *Transaction) error {
	var (
		result bidResult
		err    error
	)
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		result, err = t.sendBid(ctx, txn)
		switch {
		case err != nil:
			return err
		case result.optedInSlot:
			if result.noOfProviders == len(result.commitments) {
				// This means that all builders have committed to the bid and it
				// is a primev opted in slot. We can safely proceed to inform the
				// user that the txn was successfully sent and will be processed
				break BID_LOOP
			}
		default:
		}

		// Wait for block number to be updated to confirm transaction. If failed
		// we will retry the bid process till user cancels the operation
		included, err := t.blockTracker.CheckTxnInclusion(ctx, txn.Hash(), result.blockNumber)
		if err != nil {
			t.logger.Error("Failed to check transaction inclusion", "error", err)
			return fmt.Errorf("failed to check transaction inclusion: %w", err)
		}
		if included {
			break BID_LOOP
		}
	}

	if err := t.store.StorePreconfirmedTransaction(
		ctx,
		int64(result.blockNumber),
		txn,
		result.commitments,
	); err != nil {
		return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
	}

	switch txn.Type {
	case TxTypeRegular:
		if err := t.store.DeductBalance(ctx, txn.Sender, result.bidAmount); err != nil {
			t.logger.Error("Failed to deduct balance for sender", "sender", txn.Sender.Hex(), "error", err)
			return fmt.Errorf("failed to deduct balance for sender: %w", err)
		}
	case TxTypeDeposit:
		if err := t.store.AddBalance(ctx, txn.Sender, txn.Value()); err != nil {
			t.logger.Error("Failed to add balance for sender", "sender", txn.Sender.Hex(), "error", err)
			return fmt.Errorf("failed to add balance for sender: %w", err)
		}
	case TxTypeInstantBridge:
	}

	return nil
}

type bidResult struct {
	noOfProviders int
	blockNumber   uint64
	optedInSlot   bool
	bidAmount     *big.Int
	commitments   []*bidderapiv1.Commitment
}

func (t *TxSender) sendBid(
	ctx context.Context,
	txn *Transaction,
) (bidResult, error) {
	timeToOptIn, err := t.bidder.Estimate()
	if err != nil {
		t.logger.Error("Failed to estimate time to opt-in", "error", err)
		if !errors.Is(err, optinbidder.ErrNoSlotInCurrentEpoch) && !errors.Is(err, optinbidder.ErrNoEpochInfo) {
			return bidResult{}, err
		}
		// If we cannot estimate the time to opt-in, we assume a default value and
		// proceed with the bid process. The default value should be higher than
		// the typical block time to ensure we consider the next slot as a non-opt-in slot.
		timeToOptIn = blockTime * 32
	}

	optedInSlot := timeToOptIn <= blockTime

	price, err := t.pricer.EstimatePrice(ctx, txn.Transaction)
	if err != nil {
		t.logger.Error("Failed to estimate transaction price", "error", err)
		return bidResult{}, fmt.Errorf("failed to estimate transaction price: %w", err)
	}

	if txn.Type == TxTypeRegular {
		if !t.store.HasBalance(ctx, txn.Sender, price.BidAmount) {
			t.logger.Error("Insufficient balance for sender", "sender", txn.Sender.Hex())
			return bidResult{}, fmt.Errorf("insufficient balance for sender: %s", txn.Sender.Hex())
		}
	}

	bidC, err := t.bidder.Bid(
		ctx,
		price.BidAmount,
		big.NewInt(0),
		strings.TrimPrefix(txn.Raw, "0x"),
		&optinbidder.BidOpts{
			WaitForOptIn: optedInSlot,
			BlockNumber:  uint64(price.BlockNumber),
		},
	)
	if err != nil {
		t.logger.Error("Failed to place bid", "error", err)
		return bidResult{}, fmt.Errorf("failed to place bid: %w", err)
	}

	result := bidResult{
		commitments: make([]*bidderapiv1.Commitment, 0),
		bidAmount:   price.BidAmount,
	}
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			t.logger.Info("Context cancelled while waiting for bid status")
			return bidResult{}, ctx.Err()
		case bidStatus, more := <-bidC:
			if !more {
				t.logger.Info("Bid channel closed, no more bid statuses")
				break BID_LOOP
			}
			switch bidStatus.Type {
			case optinbidder.BidStatusNoOfProviders:
				result.noOfProviders = bidStatus.Arg.(int)
			case optinbidder.BidStatusAttempted:
				result.blockNumber = bidStatus.Arg.(uint64)
			case optinbidder.BidStatusCommitment:
				result.commitments = append(result.commitments, bidStatus.Arg.(*bidderapiv1.Commitment))
			case optinbidder.BidStatusCancelled:
				t.logger.Warn("Bid context cancelled by the bidder")
				break BID_LOOP
			case optinbidder.BidStatusFailed:
				t.logger.Error("Bid failed", "error", bidStatus.Arg)
				break BID_LOOP
			}
		}
	}
	if len(result.commitments) == 0 {
		t.logger.Error("Bid completed with no commitments")
		return bidResult{}, fmt.Errorf("bid completed with no commitments")
	}
	t.logger.Info(
		"Bid successful with commitments",
		"noOfProviders", result.noOfProviders,
		"noOfCommitments", len(result.commitments),
		"blockNumber", result.blockNumber,
		"optedInSlot", optedInSlot,
	)

	result.optedInSlot = optedInSlot
	return result, nil
}
