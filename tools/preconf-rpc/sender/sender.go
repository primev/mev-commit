package sender

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"golang.org/x/sync/errgroup"
)

type TxType int

const (
	TxTypeRegular TxType = iota
	TxTypeDeposit
	TxTypeInstantBridge
)

type TxStatus string

const (
	TxStatusPending      TxStatus = "pending"
	TxStatusPreConfirmed TxStatus = "pre-confirmed"
	TxStatusConfirmed    TxStatus = "confirmed"
	TxStatusFailed       TxStatus = "failed"
)

const (
	blockTime                    = 12               // seconds, typical Ethereum block time
	bidTimeout                   = 3 * time.Second  // timeout for bid operations
	defaultConfidence            = 90               // default confidence level for the next block
	confidenceSecondAttempt      = 95               // confidence level for the second attempt
	confidenceSubsequentAttempts = 99               // confidence level for subsequent attempts
	transactionTimeout           = 10 * time.Minute // timeout for transaction processing
)

var (
	ErrInvalidTransaction       = errors.New("invalid transaction")
	ErrUnsupportedTxType        = errors.New("unsupported transaction type")
	ErrEmptyRawTransaction      = errors.New("empty raw transaction")
	ErrEmptyTransactionTo       = errors.New("empty transaction 'to' address")
	ErrNegativeTransactionValue = errors.New("negative transaction value")
	ErrZeroGasLimit             = errors.New("zero gas limit")
	ErrTransactionCancelled     = errors.New("transaction cancelled by user")
	ErrTimeoutExceeded          = errors.New("timeout exceeded while waiting for transaction to be processed")
)

type Transaction struct {
	*types.Transaction
	Sender      common.Address
	Raw         string
	Type        TxType
	Status      TxStatus
	Details     string
	BlockNumber int64
}

type Store interface {
	AddQueuedTransaction(ctx context.Context, tx *Transaction) error
	GetQueuedTransactions(ctx context.Context) ([]*Transaction, error)
	GetCurrentNonce(ctx context.Context, sender common.Address) uint64
	HasBalance(ctx context.Context, sender common.Address, amount *big.Int) bool
	AddBalance(ctx context.Context, account common.Address, amount *big.Int) error
	DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error
	StoreTransaction(ctx context.Context, txn *Transaction, commitments []*bidderapiv1.Commitment) error
	GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*Transaction, error)
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
	EstimatePrice(ctx context.Context) map[int64]float64
}

type BlockTracker interface {
	CheckTxnInclusion(ctx context.Context, txnHash common.Hash, blockNumber uint64) (bool, error)
	NextBlockNumber() (uint64, time.Duration, error)
}

type Transferer interface {
	Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error
}

type blockAttempt struct {
	blockNumber uint64
	attempts    int
}

type txnAttempt struct {
	txnHash   common.Hash
	startTime time.Time
	attempts  []*blockAttempt
}

type Notifier interface {
	NotifyTransactionStatus(txn *Transaction, noOfAttempts int, start time.Time)
}

type TxSender struct {
	logger            *slog.Logger
	store             Store
	bidder            Bidder
	pricer            Pricer
	blockTracker      BlockTracker
	transferer        Transferer
	settlementChainId *big.Int
	eg                *errgroup.Group
	egCtx             context.Context
	trigger           chan struct{}
	workerPool        chan struct{}
	inflightTxns      map[common.Hash]chan struct{}
	inflightAccount   map[common.Address]struct{}
	inflightMu        sync.RWMutex
	processMu         sync.RWMutex
	txnAttemptHistory *lru.Cache[common.Hash, *txnAttempt]
	notifier          Notifier
	fastTrack         func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool
}

func noOpFastTrack(_ []*bidderapiv1.Commitment, _ bool) bool {
	return false
}

func NewTxSender(
	st Store,
	bidder Bidder,
	pricer Pricer,
	blockTracker BlockTracker,
	transferer Transferer,
	notifier Notifier,
	settlementChainId *big.Int,
	fastTrack func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool,
	logger *slog.Logger,
) (*TxSender, error) {
	txnAttemptHistory, err := lru.New[common.Hash, *txnAttempt](1000)
	if err != nil {
		logger.Error("Failed to create transaction attempt history cache", "error", err)
		return nil, fmt.Errorf("failed to create transaction attempt history cache: %w", err)
	}

	if fastTrack == nil {
		fastTrack = noOpFastTrack
	}

	return &TxSender{
		store:             st,
		bidder:            bidder,
		pricer:            pricer,
		blockTracker:      blockTracker,
		transferer:        transferer,
		settlementChainId: settlementChainId,
		logger:            logger.With("component", "TxSender"),
		workerPool:        make(chan struct{}, 512),
		trigger:           make(chan struct{}, 1),
		inflightTxns:      make(map[common.Hash]chan struct{}),
		inflightAccount:   make(map[common.Address]struct{}),
		txnAttemptHistory: txnAttemptHistory,
		notifier:          notifier,
		fastTrack:         fastTrack,
	}, nil
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
	if tx.To() == nil {
		return ErrEmptyTransactionTo
	}
	if tx.Value().Sign() < 0 {
		return ErrNegativeTransactionValue
	}
	if tx.Gas() == 0 {
		return ErrZeroGasLimit
	}
	return nil
}

func (t *TxSender) hasLowerNonce(ctx context.Context, tx *Transaction) bool {
	currentNonce := t.store.GetCurrentNonce(ctx, tx.Sender)
	return tx.Nonce() < currentNonce
}

func (t *TxSender) triggerSender() {
	select {
	case t.trigger <- struct{}{}:
	default:
		// Non-blocking send, if the channel is full, we do nothing
	}
}

func (t *TxSender) Enqueue(ctx context.Context, tx *Transaction) error {
	if err := validateTransaction(tx); err != nil {
		t.logger.Error("Invalid transaction", "error", err, "transaction", tx.Raw)
		return err
	}

	if t.hasLowerNonce(ctx, tx) {
		return errors.New("transaction has a lower nonce than the current highest nonce")
	}

	if err := t.store.AddQueuedTransaction(ctx, tx); err != nil {
		return err
	}

	t.triggerSender()

	return nil
}

func (t *TxSender) CancelTransaction(ctx context.Context, txnHash common.Hash) (bool, error) {
	t.inflightMu.RLock()
	cancel, found := t.inflightTxns[txnHash]
	t.inflightMu.RUnlock()
	if !found {
		return func() (bool, error) {
			// we need to hold the processMu lock till we check as a parallel goroutine
			// might try to process the transaction and update its status
			t.processMu.RLock()
			defer t.processMu.RUnlock()

			txn, err := t.store.GetTransactionByHash(ctx, txnHash)
			if err == nil {
				// if a transaction is not yet enqueued due to nonce order, we mark it
				// cancelled directly in the store
				if txn.Status == TxStatusPending {
					txn.Status = TxStatusFailed
					txn.Details = ErrTransactionCancelled.Error()
					if err := t.store.StoreTransaction(ctx, txn, nil); err != nil {
						t.logger.Error("Failed to store cancelled transaction", "hash", txnHash.Hex(), "error", err)
						return false, fmt.Errorf("failed to store cancelled transaction: %w", err)
					}
					t.logger.Info("Transaction cancelled before processing", "hash", txnHash.Hex())
					return true, nil
				}
			}
			t.logger.Warn("Transaction not found", "hash", txnHash.Hex())
			return false, nil
		}()
	}

	t.logger.Info("Cancelling transaction", "hash", txnHash.Hex())
	close(cancel) // Signal the transaction processing to stop

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.logger.Info("Context cancelled while waiting for transaction cancellation")
			return false, ctx.Err()
		case <-ticker.C:
			t.inflightMu.RLock()
			_, stillInFlight := t.inflightTxns[txnHash]
			t.inflightMu.RUnlock()
			if !stillInFlight {
				txn, err := t.store.GetTransactionByHash(ctx, txnHash)
				switch {
				case err != nil:
					t.logger.Error("Failed to get transaction by hash", "hash", txnHash.Hex(), "error", err)
					return false, fmt.Errorf("failed to get transaction by hash: %w", err)
				case txn.Status == TxStatusFailed:
					if txn.Details == ErrTransactionCancelled.Error() {
						t.logger.Info("Transaction successfully cancelled", "hash", txnHash.Hex())
						return true, nil
					}
					t.logger.Warn(
						"Transaction failed with other error",
						"hash", txnHash.Hex(),
						"status", txn.Status,
						"details", txn.Details,
					)
					return false, fmt.Errorf("transaction failed: %s", txn.Details)
				case txn.Status == TxStatusPreConfirmed || txn.Status == TxStatusConfirmed:
					t.logger.Info("Transaction already confirmed or pre-confirmed", "hash", txnHash.Hex(), "status", txn.Status)
					return false, errors.New("transaction already confirmed or pre-confirmed")
				}
			}
		}
	}
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
				t.processMu.Lock()
				t.processQueuedTransactions(t.egCtx)
				t.processMu.Unlock()
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

func (t *TxSender) markInflight(txn *Transaction) (bool, <-chan struct{}) {
	t.inflightMu.Lock()
	defer t.inflightMu.Unlock()

	if _, ok := t.inflightTxns[txn.Hash()]; ok {
		t.logger.Debug("Transaction already in flight, skipping", "hash", txn.Hash().Hex())
		return false, nil
	}
	if _, ok := t.inflightAccount[txn.Sender]; ok {
		t.logger.Debug("Transaction sender already has an inflight transaction, skipping", "sender", txn.Sender.Hex())
		t.triggerSender() // Trigger to reprocess later
		return false, nil
	}

	cancel := make(chan struct{})
	t.inflightTxns[txn.Hash()] = cancel
	t.inflightAccount[txn.Sender] = struct{}{}
	return true, cancel
}

func (t *TxSender) markCompleted(txn *Transaction) {
	t.inflightMu.Lock()
	defer t.inflightMu.Unlock()

	delete(t.inflightTxns, txn.Hash())
	delete(t.inflightAccount, txn.Sender)
}

func (t *TxSender) processQueuedTransactions(ctx context.Context) {
	txns, err := t.store.GetQueuedTransactions(ctx)
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
				defer t.triggerSender() // Trigger to reprocess after this transaction

				canExecute, cancel := t.markInflight(txn)
				if !canExecute {
					// Transaction is already being processed or sender has an inflight transaction
					return nil
				}
				defer t.markCompleted(txn)

				t.logger.Info("Processing transaction", "sender", txn.Sender.Hex(), "type", txn.Type)
				if err := t.processTransaction(ctx, txn, cancel); err != nil {
					t.logger.Error("Failed to process transaction", "sender", txn.Sender.Hex(), "error", err)
					txn.Status = TxStatusFailed
					txn.Details = err.Error()
					t.clearBlockAttemptHistory(txn)
					return t.store.StoreTransaction(ctx, txn, nil)
				}
				return nil
			})
		}
	}
}

func (t *TxSender) processTransaction(ctx context.Context, txn *Transaction, cancel <-chan struct{}) error {
	var (
		result bidResult
		err    error
	)
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cancel:
			return ErrTransactionCancelled
		default:
		}

		result, err = t.sendBid(ctx, txn)
		switch {
		case err != nil:
			if retryErr, ok := err.(*errRetry); ok {
				t.logger.Warn(
					"Retrying bid due to error",
					"error", retryErr.err,
					"retryAfter", retryErr.retryAfter,
				)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(retryErr.retryAfter):
					// Wait for the specified retry duration before retrying
				case <-cancel:
					return ErrTransactionCancelled
				}
				continue
			}
			return err
		case t.fastTrack(result.commitments, result.optedInSlot):
			// If the commitments indicate that the transaction can be fast-tracked,
			// we consider it pre-confirmed and skip further checks
			txn.Status = TxStatusPreConfirmed
			txn.BlockNumber = int64(result.blockNumber)
			t.logger.Info(
				"Transaction fast-tracked based on commitments",
				"sender", txn.Sender.Hex(),
				"type", txn.Type,
				"blockNumber", result.blockNumber,
				"bidAmount", result.bidAmount.String(),
			)
			t.clearBlockAttemptHistory(txn)
			break BID_LOOP
		case result.optedInSlot:
			if result.noOfProviders == len(result.commitments) {
				// This means that all builders have committed to the bid and it
				// is a primev opted in slot. We can safely proceed to inform the
				// user that the txn was successfully sent and will be processed
				txn.Status = TxStatusPreConfirmed
				txn.BlockNumber = int64(result.blockNumber)
				t.logger.Info(
					"Transaction pre-confirmed",
					"sender", txn.Sender.Hex(),
					"type", txn.Type,
					"blockNumber", result.blockNumber,
					"bidAmount", result.bidAmount.String(),
				)
				t.clearBlockAttemptHistory(txn)
				break BID_LOOP
			}
		default:
		}

		if result.noOfProviders > len(result.commitments) {
			t.logger.Warn(
				"Not all builders committed to the bid",
				"noOfProviders", result.noOfProviders,
				"noOfCommitments", len(result.commitments),
				"sender", txn.Sender.Hex(),
				"type", txn.Type,
				"blockNumber", result.blockNumber,
				"bidAmount", result.bidAmount.String(),
			)
			if (result.timeUntillNextBlock - time.Second) > time.Since(result.startTime) {
				// If not all builders committed, we will retry the bid process
				// immediately if we have atleast 1 second left before the next block
				continue
			}
		}

		// Wait for block number to be updated to confirm transaction. If failed
		// we will retry the bid process till user cancels the operation
		included, err := t.blockTracker.CheckTxnInclusion(ctx, txn.Hash(), result.blockNumber)
		if err != nil {
			t.logger.Error("Failed to check transaction inclusion", "error", err)
			return fmt.Errorf("failed to check transaction inclusion: %w", err)
		}
		if included {
			txn.Status = TxStatusConfirmed
			txn.BlockNumber = int64(result.blockNumber)
			t.logger.Info(
				"Transaction confirmed for non opted-in slot",
				"sender", txn.Sender.Hex(),
				"type", txn.Type,
				"blockNumber", result.blockNumber,
				"bidAmount", result.bidAmount.String(),
			)
			t.clearBlockAttemptHistory(txn)
			break BID_LOOP
		}
	}

	if err := t.store.StoreTransaction(ctx, txn, result.commitments); err != nil {
		return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
	}

	switch txn.Type {
	case TxTypeRegular:
		if err := t.store.DeductBalance(ctx, txn.Sender, result.bidAmount); err != nil {
			t.logger.Error("Failed to deduct balance for sender", "sender", txn.Sender.Hex(), "error", err)
			return fmt.Errorf("failed to deduct balance for sender: %w", err)
		}
	case TxTypeDeposit:
		balanceToAdd := new(big.Int).Sub(txn.Value(), result.bidAmount)
		if err := t.store.AddBalance(ctx, txn.Sender, balanceToAdd); err != nil {
			t.logger.Error("Failed to add balance for sender", "sender", txn.Sender.Hex(), "error", err)
			return fmt.Errorf("failed to add balance for sender: %w", err)
		}
	case TxTypeInstantBridge:
		amountToBridge := new(big.Int).Sub(txn.Value(), new(big.Int).Mul(result.bidAmount, big.NewInt(2)))
		if err := t.transferer.Transfer(ctx, txn.Sender, t.settlementChainId, amountToBridge); err != nil {
			t.logger.Error("Failed to transfer funds for instant bridge", "sender", txn.Sender.Hex(), "error", err)
			return fmt.Errorf("failed to transfer funds for instant bridge: %w", err)
		}
	}

	return nil
}

type errRetry struct {
	err        error
	retryAfter time.Duration
}

func (e *errRetry) Error() string {
	return fmt.Sprintf("retry after %s: %v", e.retryAfter, e.err)
}

type bidResult struct {
	startTime           time.Time
	timeUntillNextBlock time.Duration
	noOfProviders       int
	blockNumber         uint64
	optedInSlot         bool
	bidAmount           *big.Int
	commitments         []*bidderapiv1.Commitment
}

func (t *TxSender) sendBid(
	ctx context.Context,
	txn *Transaction,
) (bidResult, error) {
	timeToOptIn, err := t.bidder.Estimate()
	if err != nil {
		t.logger.Warn("Failed to estimate time to opt-in", "error", err)
		// If we cannot estimate the time to opt-in, we assume a default value and
		// proceed with the bid process. The default value should be higher than
		// the typical block time to ensure we consider the next slot as a non-opt-in slot.
		timeToOptIn = blockTime * 32
	}

	start := time.Now()
	bidBlockNo, timeUntilNextBlock, err := t.blockTracker.NextBlockNumber()
	if err != nil {
		t.logger.Error("Failed to get next block number", "error", err)
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("failed to get next block number: %w", err),
			retryAfter: time.Second,
		}
	}

	if timeUntilNextBlock <= time.Second {
		t.logger.Warn("Next block time is too short, skipping bid", "timeUntilNextBlock", timeUntilNextBlock)
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("next block time is too short: %s", timeUntilNextBlock),
			retryAfter: time.Second,
		}
	}

	prices := t.pricer.EstimatePrice(ctx)

	// Allow for certain level of tolerance w.r.t timestamps
	optedInSlot := math.Abs(float64(timeToOptIn)-float64(timeUntilNextBlock.Seconds())) < float64(blockTime/3)

	cctx, cancel := context.WithTimeout(ctx, bidTimeout)
	defer cancel()

	cost, err := t.calculatePriceForNextBlock(txn, bidBlockNo, prices, optedInSlot)
	if err != nil {
		t.logger.Error("Failed to calculate price for next block", "error", err)
		if errors.Is(err, ErrTimeoutExceeded) {
			t.logger.Warn("Timeout exceeded while trying to process transaction", "txnHash", txn.Hash().Hex())
			return bidResult{}, ErrTimeoutExceeded
		}
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("failed to calculate price: %w", err),
			retryAfter: time.Second,
		}
	}

	slashAmount := big.NewInt(0)
	switch txn.Type {
	case TxTypeRegular:
		if !t.store.HasBalance(ctx, txn.Sender, cost) {
			t.logger.Error("Insufficient balance for sender", "sender", txn.Sender.Hex())
			return bidResult{}, fmt.Errorf("insufficient balance for sender: %s", txn.Sender.Hex())
		}
	case TxTypeDeposit:
		if txn.Value().Cmp(cost) < 0 {
			t.logger.Error(
				"Deposit amount is less than price of deposit",
				"sender", txn.Sender.Hex(),
				"deposit", txn.Value().String(),
				"price", cost.String(),
			)
			return bidResult{}, fmt.Errorf(
				"deposit amount is less than price of deposit: %s, deposit: %s, price: %s",
				txn.Sender.Hex(),
				txn.Value().String(),
				cost.String(),
			)
		}
	case TxTypeInstantBridge:
		costOfBridge := new(big.Int).Mul(cost, big.NewInt(2)) // 2x the price for instant bridge
		if txn.Value().Cmp(costOfBridge) < 0 {
			t.logger.Error(
				"Instant bridge amount is less than price of bridge",
				"sender", txn.Sender.Hex(),
				"bridge", txn.Value().String(),
				"price", costOfBridge.String(),
			)
			return bidResult{}, fmt.Errorf(
				"instant bridge amount is less than price of bridge: %s, bridge: %s, price: %s",
				txn.Sender.Hex(),
				txn.Value().String(),
				costOfBridge.String(),
			)
		}
		slashAmount = new(big.Int).Set(txn.Value())
	}

	bidC, err := t.bidder.Bid(
		cctx,
		cost,
		slashAmount,
		strings.TrimPrefix(txn.Raw, "0x"),
		&optinbidder.BidOpts{
			WaitForOptIn:      false,
			BlockNumber:       uint64(bidBlockNo),
			RevertingTxHashes: []string{txn.Hash().Hex()},
			DecayDuration:     bidTimeout * 2,
		},
	)
	if err != nil {
		t.logger.Error("Failed to place bid", "error", err)
		return bidResult{}, fmt.Errorf("failed to place bid: %w", err)
	}

	result := bidResult{
		commitments:         make([]*bidderapiv1.Commitment, 0),
		bidAmount:           cost,
		startTime:           start,
		timeUntillNextBlock: timeUntilNextBlock,
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
	t.logger.Info(
		"Bid operation complete",
		"noOfProviders", result.noOfProviders,
		"noOfCommitments", len(result.commitments),
		"blockNumber", result.blockNumber,
		"optedInSlot", optedInSlot,
	)

	result.optedInSlot = optedInSlot
	return result, nil
}

func (t *TxSender) calculatePriceForNextBlock(
	txn *Transaction,
	bidBlockNo uint64,
	prices map[int64]float64,
	optedInSlot bool,
) (*big.Int, error) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		attempts = &txnAttempt{
			txnHash:   txn.Hash(),
			startTime: time.Now(),
		}
	}

	if time.Since(attempts.startTime) > transactionTimeout {
		return nil, ErrTimeoutExceeded
	}

	// default confidence level for the next block
	confidence := defaultConfidence
	isRetry := false

	for i := len(attempts.attempts) - 1; i >= 0; i-- {
		if attempts.attempts[i].blockNumber < bidBlockNo {
			break // We only care about attempts for the current block
		}
		if attempts.attempts[i].blockNumber == bidBlockNo {
			isRetry = true
			attempts.attempts[i].attempts++
			switch {
			case attempts.attempts[i].attempts == 2:
				confidence = confidenceSecondAttempt
			case attempts.attempts[i].attempts > 2:
				confidence = confidenceSubsequentAttempts
			}
			break // No need to check further attempts for the same block
		}
	}

	if optedInSlot {
		confidence = confidenceSubsequentAttempts
	}

	// If this is the first attempt for the next block, we add it to the history
	if !isRetry {
		attempts.attempts = append(attempts.attempts, &blockAttempt{
			blockNumber: bidBlockNo,
			attempts:    1,
		})
	}

	_ = t.txnAttemptHistory.Add(txn.Hash(), attempts)

	for conf, price := range prices {
		if conf == int64(confidence) {
			// the gwei value is in float, so we need to convert it to wei before multiplying with gas limit
			priceInWei := price * 1e9 // Convert Gwei to Wei
			return new(big.Int).Mul(big.NewInt(int64(priceInWei)), big.NewInt(int64(txn.Gas()))), nil
		}
	}

	return nil, fmt.Errorf(
		"no estimated price found for block %d with confidence %d", bidBlockNo, confidence,
	)
}

func (t *TxSender) clearBlockAttemptHistory(txn *Transaction) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		return
	}

	totalAttempts := 0
	for _, attempt := range attempts.attempts {
		totalAttempts += attempt.attempts
	}

	t.logger.Info(
		"Clearing block attempt history for transaction",
		"hash", txn.Hash().Hex(),
		"blockAttempts", len(attempts.attempts),
		"startTime", attempts.startTime.Format(time.RFC3339),
		"startBlockNumber", attempts.attempts[0].blockNumber,
		"totalAttempts", totalAttempts,
	)

	_ = t.txnAttemptHistory.Remove(txn.Hash())

	t.notifier.NotifyTransactionStatus(txn, totalAttempts, attempts.startTime)
}
