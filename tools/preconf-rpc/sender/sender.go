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
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
	explorersubmitter "github.com/primev/mev-commit/tools/preconf-rpc/explorer-submitter"
	"github.com/prometheus/client_golang/prometheus"
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
	bidTimeout                   = 3 * time.Second  // timeout for bid operation
	defaultConfidence            = 90               // default confidence level for the next block
	confidenceSecondAttempt      = 95               // confidence level for the second attempt
	confidenceSubsequentAttempts = 99               // confidence level for subsequent attempts
	transactionTimeout           = 10 * time.Minute // timeout for transaction processing
	maxAttemptsPerBlock          = 10               // maximum attempts per block
	defaultRetryDelay            = 500 * time.Millisecond
)

var (
	ErrInvalidTransaction          = errors.New("invalid transaction")
	ErrUnsupportedTxType           = errors.New("unsupported transaction type")
	ErrEmptyRawTransaction         = errors.New("empty raw transaction")
	ErrEmptyTransactionTo          = errors.New("empty transaction 'to' address")
	ErrNegativeTransactionValue    = errors.New("negative transaction value")
	ErrZeroGasLimit                = errors.New("zero gas limit")
	ErrTransactionCancelled        = errors.New("transaction cancelled by user")
	ErrTimeoutExceeded             = errors.New("timeout exceeded while waiting for transaction to be processed")
	ErrMaxAttemptsPerBlockExceeded = errors.New("maximum attempts exceeded for transaction in the current block")
	ErrNonceTooHigh                = errors.New("nonce too high")
	ErrNonceTooLow                 = errors.New("nonce too low")
)

type Transaction struct {
	*types.Transaction
	Sender      common.Address
	Raw         string
	Type        TxType
	Status      TxStatus
	Details     string
	BlockNumber int64
	Constraint  *bidderapiv1.PositionConstraint
	// local fields not stored in DB
	noOfProviders int
	commitments   []*bidderapiv1.Commitment
	logs          []*types.Log
	isSwap        bool
}

type Store interface {
	AddQueuedTransaction(ctx context.Context, tx *Transaction) error
	GetQueuedTransactions(ctx context.Context) ([]*Transaction, error)
	GetCurrentNonce(ctx context.Context, sender common.Address) uint64
	HasBalance(ctx context.Context, sender common.Address, amount *big.Int) bool
	AddBalance(ctx context.Context, account common.Address, amount *big.Int) error
	DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error
	StoreTransaction(ctx context.Context, txn *Transaction, commitments []*bidderapiv1.Commitment, logs []*types.Log) error
	GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*Transaction, error)
}

type Bidder interface {
	Estimate() (int64, error)
	Bid(
		ctx context.Context,
		bidAmount *big.Int,
		slashAmount *big.Int,
		rawTx string,
		opts *bidder.BidOpts,
	) (chan bidder.BidStatus, error)
	ConnectedProviders(ctx context.Context) ([]string, error)
}

type Pricer interface {
	EstimatePrice(ctx context.Context) map[int64]float64
}

type BlockTracker interface {
	WaitForTxnInclusion(txnHash common.Hash) chan uint64
	NextBlockNumber() (uint64, time.Duration, error)
	LatestBlockNumber() uint64
	AccountNonce(ctx context.Context, account common.Address) (uint64, error)
}

type Transferer interface {
	Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error
}

type Simulator interface {
	Simulate(ctx context.Context, txRaw string) ([]*types.Log, bool, error)
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
	NotifyTransactionStatus(txn *Transaction, noOfAttempts, noOfBlocks int, timeTaken time.Duration)
}

type Backrunner interface {
	Backrun(ctx context.Context, rawTx string, commitments []*bidderapiv1.Commitment) error
}

type TxSender struct {
	logger            *slog.Logger
	store             Store
	bidder            Bidder
	pricer            Pricer
	blockTracker      BlockTracker
	transferer        Transferer
	backrunner        Backrunner
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
	simulator         Simulator
	fastTrack         func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool
	bidTimeout        time.Duration
	timeoutMtx        sync.RWMutex
	receiptSignal     map[common.Hash][]chan struct{}
	receiptMtx        sync.Mutex
	metrics           *metrics
	explorerConfig    explorersubmitter.Config
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
	simulator Simulator,
	backrunner Backrunner,
	settlementChainId *big.Int,
	explorerConfig explorersubmitter.Config,
	logger *slog.Logger,
) (*TxSender, error) {
	txnAttemptHistory, err := lru.New[common.Hash, *txnAttempt](1000)
	if err != nil {
		logger.Error("Failed to create transaction attempt history cache", "error", err)
		return nil, fmt.Errorf("failed to create transaction attempt history cache: %w", err)
	}

	return &TxSender{
		store:             st,
		bidder:            bidder,
		pricer:            pricer,
		blockTracker:      blockTracker,
		transferer:        transferer,
		backrunner:        backrunner,
		settlementChainId: settlementChainId,
		logger:            logger.With("component", "TxSender"),
		workerPool:        make(chan struct{}, 512),
		trigger:           make(chan struct{}, 1),
		inflightTxns:      make(map[common.Hash]chan struct{}),
		inflightAccount:   make(map[common.Address]struct{}),
		txnAttemptHistory: txnAttemptHistory,
		notifier:          notifier,
		simulator:         simulator,
		fastTrack:         noOpFastTrack,
		bidTimeout:        bidTimeout,
		receiptSignal:     make(map[common.Hash][]chan struct{}),
		metrics:           newMetrics(),
		explorerConfig:    explorerConfig,
	}, nil
}

func (t *TxSender) Metrics() []prometheus.Collector {
	return []prometheus.Collector{
		t.metrics.connectedProviders,
		t.metrics.queuedTransactions,
		t.metrics.inflightTransactions,
		t.metrics.preconfDurationsProvider,
		t.metrics.preconfCountsProvider,
		t.metrics.blockAttemptsToConfirmation,
		t.metrics.totalAttemptsToConfirmation,
		t.metrics.timeToConfirmation,
		t.metrics.timeToFirstPreconfirmation,
		t.metrics.bidPriorityFee,
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

func (t *TxSender) hasCorrectNonce(ctx context.Context, tx *Transaction) error {
	currentNonce := t.store.GetCurrentNonce(ctx, tx.Sender) + 1
	backendNonce, err := t.blockTracker.AccountNonce(ctx, tx.Sender)
	if err == nil {
		if backendNonce > currentNonce {
			currentNonce = backendNonce
		}
	}
	switch {
	case tx.Nonce() < currentNonce:
		return ErrNonceTooLow
	case tx.Nonce() > currentNonce:
		return ErrNonceTooHigh
	}

	return nil
}

func (t *TxSender) triggerSender() {
	select {
	case t.trigger <- struct{}{}:
	default:
		// Non-blocking send, if the channel is full, we do nothing
	}
}

func (t *TxSender) SetFastTrackFunc(fastTrack func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool) {
	t.fastTrack = fastTrack
}

func (t *TxSender) Enqueue(ctx context.Context, tx *Transaction) error {
	if err := validateTransaction(tx); err != nil {
		t.logger.Error("Invalid transaction", "error", err, "transaction", tx.Raw)
		return err
	}

	if err := t.hasCorrectNonce(ctx, tx); err != nil {
		return err
	}

	if err := t.store.AddQueuedTransaction(ctx, tx); err != nil {
		return err
	}

	t.triggerSender()

	go func() {
		// extra caution in case of errors
		defer func() {
			if r := recover(); r != nil {
				t.logger.Error("Panic in explorer submitter", "error", r)
			}
		}()

		// get tx info
		from := tx.Sender.Hex()
		to := ""
		if tx.To() != nil {
			to = tx.To().Hex()
		}

		err := explorersubmitter.Submit(
			context.Background(),
			t.explorerConfig,
			tx.Hash().Hex(),
			from,
			to,
		)
		if err != nil {
			t.logger.Error("Failed to submit tx to explorer", "error", err)
		} else {
			t.logger.Info("Successfully submitted tx to explorer", "hash", tx.Hash().Hex())
		}
	}()

	return nil
}

func (t *TxSender) WaitForReceiptAvailable(ctx context.Context, txnHash common.Hash) <-chan struct{} {
	t.receiptMtx.Lock()
	defer t.receiptMtx.Unlock()

	signal, found := t.receiptSignal[txnHash]
	if !found {
		signal = []chan struct{}{}
	}
	newSignal := make(chan struct{})
	signal = append(signal, newSignal)
	t.receiptSignal[txnHash] = signal
	return newSignal
}

func (t *TxSender) signalReceiptAvailable(txnHash common.Hash) {
	t.receiptMtx.Lock()
	defer t.receiptMtx.Unlock()

	signals, found := t.receiptSignal[txnHash]
	if !found {
		return
	}
	for _, sig := range signals {
		close(sig)
	}
	delete(t.receiptSignal, txnHash)
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
					if err := t.store.StoreTransaction(ctx, txn, nil, nil); err != nil {
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

func (t *TxSender) UpdateBidTimeout(timeout time.Duration) {
	t.timeoutMtx.Lock()
	defer t.timeoutMtx.Unlock()

	t.bidTimeout = timeout
}

func (t *TxSender) getBidTimeout() time.Duration {
	t.timeoutMtx.RLock()
	defer t.timeoutMtx.RUnlock()

	return t.bidTimeout
}

func (t *TxSender) Start(ctx context.Context) chan struct{} {
	t.eg, t.egCtx = errgroup.WithContext(ctx)
	done := make(chan struct{})

	t.eg.Go(func() error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-t.egCtx.Done():
				t.logger.Info("Context cancelled, stopping TxSender")
				return ctx.Err()
			case <-t.trigger:
				ticker.Reset(1 * time.Second)
			case <-ticker.C:
			}
			t.processMu.Lock()
			t.processQueuedTransactions(t.egCtx)
			t.processMu.Unlock()
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
	t.metrics.inflightTransactions.Inc()
	return true, cancel
}

func (t *TxSender) markCompleted(txn *Transaction) {
	t.inflightMu.Lock()
	defer t.inflightMu.Unlock()

	delete(t.inflightTxns, txn.Hash())
	delete(t.inflightAccount, txn.Sender)
	t.metrics.inflightTransactions.Dec()
}

func (t *TxSender) processQueuedTransactions(ctx context.Context) {
	txns, err := t.store.GetQueuedTransactions(ctx)
	if err != nil {
		t.logger.Error("Failed to get queued transactions", "error", err)
		return
	}
	t.metrics.queuedTransactions.Set(float64(len(txns)))
	if len(txns) == 0 {
		t.logger.Debug("No queued transactions to process")
		return
	}
	t.logger.Debug("Processing queued transactions", "count", len(txns))
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
					t.clearBlockAttemptHistory(txn, time.Now())
					defer t.signalReceiptAvailable(txn.Hash())
					return t.store.StoreTransaction(ctx, txn, nil, nil)
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
	logger := t.logger.With(
		"transactionHash", txn.Hash().Hex(),
		"sender", txn.Sender.Hex(),
		"type", txn.Type,
	)

	retryTicker := time.NewTicker(defaultRetryDelay)
	defer retryTicker.Stop()
	inclusion := t.blockTracker.WaitForTxnInclusion(txn.Hash())

BID_LOOP:
	for {
		result, err = t.sendBid(ctx, txn)
		switch {
		case err != nil:
			if retryErr, ok := err.(*errRetry); ok {
				logger.Warn(
					"Retrying bid due to error",
					"error", retryErr.err,
					"retryAfter", retryErr.retryAfter,
				)
				retryTicker.Reset(retryErr.retryAfter)
			} else if errors.Is(err, ErrMaxAttemptsPerBlockExceeded) {
				retryTicker.Reset(result.timeUntillNextBlock + 500*time.Millisecond)
			} else {
				return err
			}
		case txn.noOfProviders == len(txn.commitments):
			if result.optedInSlot {
				if txn.Status != TxStatusPreConfirmed {
					t.metrics.timeToFirstPreconfirmation.Observe(float64(time.Since(result.startTime).Milliseconds()))
				}
				// This means that all builders have committed to the bid and it
				// is a primev opted in slot. We can safely proceed to inform the
				// user that the txn was successfully sent and will be processed
				txn.Status = TxStatusPreConfirmed
				txn.BlockNumber = int64(result.blockNumber)
				logger.Info(
					"Transaction pre-confirmed",
					"blockNumber", result.blockNumber,
					"bidAmount", result.bidAmount.String(),
				)
				if err := t.store.StoreTransaction(ctx, txn, txn.commitments, txn.logs); err != nil {
					return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
				}
				t.signalReceiptAvailable(txn.Hash())
			}
			retryTicker.Reset(result.timeUntillNextBlock + 1*time.Second)
		default:
			logger.Warn(
				"Not all builders committed to the bid",
				"noOfProviders", txn.noOfProviders,
				"noOfCommitments", len(txn.commitments),
				"blockNumber", result.blockNumber,
				"bidAmount", result.bidAmount.String(),
			)
			retryTicker.Reset(defaultRetryDelay)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cancel:
			return ErrTransactionCancelled
		case <-retryTicker.C:
			// Continue to the next iteration after the retry delay
		case bNo := <-inclusion:
			if txn.Status != TxStatusPreConfirmed {
				// It could happen that the transaction got included but we got the signal
				// late and made a failed attempt. So we should update the commitments and
				// logs from the last successful bid attempt.
				txn.Status = TxStatusConfirmed
				txn.BlockNumber = int64(bNo)
				logger.Info(
					"Transaction confirmed",
					"blockNumber", bNo,
					"bidAmount", result.bidAmount.String(),
				)
				if err := t.store.StoreTransaction(ctx, txn, txn.commitments, txn.logs); err != nil {
					return fmt.Errorf("failed to store preconfirmed transaction: %w", err)
				}
				t.signalReceiptAvailable(txn.Hash())
			}
			endTime := time.Now()
			if len(txn.commitments) > 0 {
				endTime = time.UnixMilli(txn.commitments[0].DispatchTimestamp)
			}
			t.clearBlockAttemptHistory(txn, endTime)
			break BID_LOOP
		}
	}

	amount := big.NewInt(0)
	for _, cmt := range txn.commitments {
		amt, ok := new(big.Int).SetString(cmt.BidAmount, 10)
		if ok && amt.Cmp(amount) > 0 {
			amount = amt
		}
	}

	switch txn.Type {
	case TxTypeRegular:
		if err := t.store.DeductBalance(ctx, txn.Sender, amount); err != nil {
			logger.Error("Failed to deduct balance for sender", "error", err)
			return fmt.Errorf("failed to deduct balance for sender: %w", err)
		}
	case TxTypeDeposit:
		balanceToAdd := new(big.Int).Sub(txn.Value(), amount)
		if err := t.store.AddBalance(ctx, txn.Sender, balanceToAdd); err != nil {
			logger.Error("Failed to add balance for sender", "error", err)
			return fmt.Errorf("failed to add balance for sender: %w", err)
		}
	case TxTypeInstantBridge:
		amountToBridge := new(big.Int).Sub(txn.Value(), new(big.Int).Mul(amount, big.NewInt(2)))
		if err := t.transferer.Transfer(ctx, txn.Sender, t.settlementChainId, amountToBridge); err != nil {
			logger.Error("Failed to transfer funds for instant bridge", "error", err)
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
	blockNumber         uint64
	optedInSlot         bool
	bidAmount           *big.Int
}

func (t *TxSender) sendBid(
	ctx context.Context,
	txn *Transaction,
) (bidResult, error) {
	start := time.Now()
	logger := t.logger.With(
		"transactionHash", txn.Hash().Hex(),
		"sender", txn.Sender.Hex(),
		"type", txn.Type,
	)

	timeToOptIn, err := t.bidder.Estimate()
	if err != nil {
		logger.Warn("Failed to estimate time to opt-in", "error", err)
		// If we cannot estimate the time to opt-in, we assume a default value and
		// proceed with the bid process. The default value should be higher than
		// the typical block time to ensure we consider the next slot as a non-opt-in slot.
		timeToOptIn = blockTime * 32
	}

	bidBlockNo, timeUntilNextBlock, err := t.blockTracker.NextBlockNumber()
	if err != nil {
		logger.Error("Failed to get next block number", "error", err)
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("failed to get next block number: %w", err),
			retryAfter: time.Second,
		}
	}
	logger.Debug("Next block info", "bidBlockNo", bidBlockNo, "timeUntilNextBlock", timeUntilNextBlock)

	if timeUntilNextBlock <= 500*time.Millisecond {
		logger.Warn("Next block time is too short, skipping bid", "timeUntilNextBlock", timeUntilNextBlock)
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("next block time is too short: %s", timeUntilNextBlock),
			retryAfter: defaultRetryDelay,
		}
	}

	prices := t.pricer.EstimatePrice(ctx)

	// Allow for certain level of tolerance w.r.t timestamps
	optedInSlot := math.Abs(float64(timeToOptIn)-float64(timeUntilNextBlock.Seconds())) < float64(blockTime/3)

	cctx, cancel := context.WithTimeout(ctx, t.getBidTimeout())
	defer cancel()

	cost, isRetry, err := t.calculatePriceForNextBlock(txn, bidBlockNo, prices, optedInSlot)
	if err != nil {
		logger.Error("Failed to calculate price for next block", "error", err)
		if errors.Is(err, ErrTimeoutExceeded) || errors.Is(err, ErrMaxAttemptsPerBlockExceeded) {
			// We propagate these errors as is
			return bidResult{}, err
		}
		return bidResult{}, &errRetry{
			err:        fmt.Errorf("failed to calculate price: %w", err),
			retryAfter: time.Second,
		}
	}

	var ignoreProviders []string
	if isRetry && len(txn.commitments) > 0 {
		for _, cmt := range txn.commitments {
			ignoreProviders = append(ignoreProviders, cmt.ProviderAddress)
		}
		logger.Info(
			"Retrying bid, ignoring previously committed providers",
			"ignoreProviders", ignoreProviders,
		)
	}

	slashAmount := big.NewInt(0)
	switch txn.Type {
	case TxTypeRegular:
		if !t.store.HasBalance(ctx, txn.Sender, cost) {
			logger.Error("Insufficient balance for sender")
			return bidResult{}, fmt.Errorf("insufficient balance for sender: %s", txn.Sender.Hex())
		}
	case TxTypeDeposit:
		if txn.Value().Cmp(cost) < 0 {
			logger.Error(
				"Deposit amount is less than price of deposit",
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
			logger.Error(
				"Instant bridge amount is less than price of bridge",
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

	if !isRetry {
		logs, isSwap, err := t.simulator.Simulate(ctx, txn.Raw)
		if err != nil {
			logger.Error("Failed to simulate transaction", "error", err, "blockNumber", bidBlockNo)
			if len(txn.commitments) > 0 && txn.commitments[0].BlockNumber+1 == int64(bidBlockNo) {
				// Could happen that it takes time to get confirmation of txn inclusion
				// so simulation would return error but we should retry after a delay to allow
				// the transaction to be included
				return bidResult{}, &errRetry{
					err:        fmt.Errorf("failed to simulate transaction: %w", err),
					retryAfter: 2 * time.Second,
				}
			}
			return bidResult{}, fmt.Errorf("failed to simulate transaction: %w", err)
		}
		providers, err := t.bidder.ConnectedProviders(ctx)
		if err != nil {
			logger.Error("Failed to get connected providers", "error", err)
			return bidResult{}, fmt.Errorf("failed to get connected providers: %w", err)
		}
		txn.logs = logs
		txn.isSwap = isSwap
		txn.noOfProviders = len(providers)
		t.metrics.connectedProviders.Set(float64(len(providers)))
		// We could have already made a attempt on the previous block but the block
		// update hasn't happened yet. This means that the bid might fail, but
		// we should retain the previous commitments. Only clear if we get new
		// commitments for the new block.
	}

	bidStart := time.Now()
	bidC, err := t.bidder.Bid(
		cctx,
		cost,
		slashAmount,
		strings.TrimPrefix(txn.Raw, "0x"),
		&bidder.BidOpts{
			WaitForOptIn:      false,
			BlockNumber:       uint64(bidBlockNo),
			RevertingTxHashes: []string{txn.Hash().Hex()},
			DecayDuration:     t.getBidTimeout() * 2,
			Constraint:        txn.Constraint,
			IgnoreProviders:   ignoreProviders,
		},
	)
	if err != nil {
		logger.Error("Failed to place bid", "error", err)
		return bidResult{}, fmt.Errorf("failed to place bid: %w", err)
	}

	result := bidResult{
		bidAmount:           cost,
		blockNumber:         bidBlockNo,
		startTime:           start,
		timeUntillNextBlock: timeUntilNextBlock,
	}
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled while waiting for bid status")
			return bidResult{}, ctx.Err()
		case bidStatus, more := <-bidC:
			if !more {
				logger.Info("Bid channel closed, no more bid statuses")
				break BID_LOOP
			}
			switch bidStatus.Type {
			case bidder.BidStatusCommitment:
				if len(txn.commitments) > 0 {
					if txn.commitments[0].BlockNumber != int64(bidBlockNo) {
						txn.commitments = nil // clear previous commitments for new block
					}
				}
				cmt := bidStatus.Arg.(*bidderapiv1.Commitment)
				txn.commitments = append(txn.commitments, cmt)
				if t.fastTrack(txn.commitments, optedInSlot) && txn.Status != TxStatusPreConfirmed {
					txn.Status = TxStatusPreConfirmed
					txn.BlockNumber = int64(bidBlockNo)
					logger.Info(
						"Transaction fast-tracked based on commitments",
						"blockNumber", result.blockNumber,
						"bidAmount", result.bidAmount.String(),
					)
					if err := t.store.StoreTransaction(ctx, txn, txn.commitments, txn.logs); err != nil {
						logger.Error("Failed to store fast-tracked transaction", "error", err)
					}
					t.signalReceiptAvailable(txn.Hash())
					t.metrics.timeToFirstPreconfirmation.Observe(float64(time.Since(start).Milliseconds()))
				}
				t.metrics.preconfDurationsProvider.WithLabelValues(cmt.ProviderAddress).Set(float64(time.Since(bidStart).Milliseconds()))
				t.metrics.preconfCountsProvider.WithLabelValues(cmt.ProviderAddress).Inc()
			case bidder.BidStatusCancelled:
				logger.Warn("Bid context cancelled by the bidder")
				break BID_LOOP
			case bidder.BidStatusFailed:
				logger.Error("Bid failed", "error", bidStatus.Arg)
				break BID_LOOP
			}
		}
	}
	logger.Info(
		"Bid operation complete",
		"noOfProviders", txn.noOfProviders,
		"noOfCommitments", len(txn.commitments),
		"blockNumber", result.blockNumber,
		"optedInSlot", optedInSlot,
	)

	if len(txn.commitments) > 0 && txn.isSwap {
		if err := t.backrunner.Backrun(ctx, txn.Raw, txn.commitments); err != nil {
			logger.Error("Failed to backrun transaction", "error", err)
		}
		logger.Info("Backrun operation initiated for transaction", "hash", txn.Hash().Hex())
	}

	result.optedInSlot = optedInSlot
	return result, nil
}

func (t *TxSender) calculatePriceForNextBlock(
	txn *Transaction,
	bidBlockNo uint64,
	prices map[int64]float64,
	optedInSlot bool,
) (*big.Int, bool, error) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		attempts = &txnAttempt{
			txnHash:   txn.Hash(),
			startTime: time.Now(),
		}
	}

	if time.Since(attempts.startTime) > transactionTimeout {
		return nil, false, ErrTimeoutExceeded
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
			case attempts.attempts[i].attempts > maxAttemptsPerBlock:
				return nil, false, fmt.Errorf("%w: block %d", ErrMaxAttemptsPerBlockExceeded, bidBlockNo)
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
			t.metrics.bidPriorityFee.Set(price)
			return new(big.Int).Mul(big.NewInt(int64(priceInWei)), big.NewInt(int64(txn.Gas()))), isRetry, nil
		}
	}

	return nil, false, fmt.Errorf(
		"no estimated price found for block %d with confidence %d", bidBlockNo, confidence,
	)
}

func (t *TxSender) clearBlockAttemptHistory(txn *Transaction, endTime time.Time) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		return
	}

	totalAttempts := 0
	blockAttempts := len(attempts.attempts)

	for _, attempt := range attempts.attempts {
		totalAttempts += attempt.attempts
	}

	t.logger.Info(
		"Clearing block attempt history for transaction",
		"transactionHash", txn.Hash().Hex(),
		"blockNumber", txn.BlockNumber,
		"blockAttempts", blockAttempts,
		"startTime", attempts.startTime.Format(time.RFC3339),
		"startBlockNumber", attempts.attempts[0].blockNumber,
		"totalAttempts", totalAttempts,
	)

	t.metrics.blockAttemptsToConfirmation.Observe(float64(blockAttempts))
	t.metrics.totalAttemptsToConfirmation.Observe(float64(totalAttempts))

	_ = t.txnAttemptHistory.Remove(txn.Hash())

	timeTaken := endTime.Sub(attempts.startTime).Round(time.Millisecond)
	t.metrics.timeToConfirmation.Observe(float64(timeTaken.Milliseconds()))
	t.notifier.NotifyTransactionStatus(txn, totalAttempts, blockAttempts, timeTaken)
}
