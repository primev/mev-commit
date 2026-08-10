package sender

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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
	"github.com/primev/mev-commit/tools/preconf-rpc/sim"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

type TxType int

const (
	TxTypeRegular TxType = iota
	TxTypeDeposit
	TxTypeInstantBridge
	TxTypeFastSwap // Executor-submitted fastswap transactions (skip balance check)
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
	defaultConfidence            = 75               // default confidence level for the next block
	defaultSwapConfidence        = 70               // default confidence level for fastswap transactions
	confidenceSecondAttempt      = 80               // confidence level for the second attempt
	confidenceSubsequentAttempts = 85               // confidence level for subsequent attempts
	topOfBlockSwapConfidence     = 80               // confidence level for top-of-block fastswap first attempt
	topOfBlockSubsequentAttempts = 90               // confidence level for top-of-block fastswap retries
	transactionTimeout           = 10 * time.Minute // timeout for transaction processing
	maxAttemptsPerBlock          = 10               // maximum attempts per block
	defaultRetryDelay            = 500 * time.Millisecond
	maxBundleSize                = 5 // maximum number of transactions in one bundle bid
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
	simFailed     bool // tracks if the last simulation failed, to force re-sim on retry
}

// encryptForLog encrypts plaintext using AES-256-GCM and returns a base64-encoded
// ciphertext string suitable for logging. Returns empty string if key is nil.
func encryptForLog(key []byte, plaintext string) string {
	if len(key) == 0 {
		return ""
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return ""
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

// ParseLogEncryptionKey parses a hex-encoded 32-byte AES-256 key.
// Returns nil if the input is empty.
func ParseLogEncryptionKey(hexKey string) ([]byte, error) {
	if hexKey == "" {
		return nil, nil
	}
	key, err := hex.DecodeString(strings.TrimPrefix(hexKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid log encryption key hex: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("log encryption key must be 32 bytes, got %d", len(key))
	}
	return key, nil
}

func effectiveFeePerGas(tx *types.Transaction) *big.Int {
	if tx == nil {
		return big.NewInt(0)
	}

	if tx.Type() == types.DynamicFeeTxType {
		if tx.GasFeeCap() != nil {
			return new(big.Int).Set(tx.GasFeeCap())
		}
	}

	if tx.GasPrice() != nil {
		return new(big.Int).Set(tx.GasPrice())
	}

	if tx.GasFeeCap() != nil {
		return new(big.Int).Set(tx.GasFeeCap())
	}

	return big.NewInt(0)
}

type Store interface {
	AddQueuedTransaction(ctx context.Context, tx *Transaction) error
	GetQueuedTransactions(ctx context.Context) ([]*Transaction, error)
	GetCurrentNonce(ctx context.Context, sender common.Address) (uint64, bool)
	HasBalance(ctx context.Context, sender common.Address, amount *big.Int) bool
	AddBalance(ctx context.Context, account common.Address, amount *big.Int) error
	DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error
	StoreTransaction(ctx context.Context, txn *Transaction, commitments []*bidderapiv1.Commitment, logs []*types.Log) error
	GetTransactionByHash(ctx context.Context, txnHash common.Hash) (*Transaction, error)
	StoreReceipt(ctx context.Context, receipt *types.Receipt) error
}

type Bidder interface {
	Estimate() (int64, error)
	Bid(
		ctx context.Context,
		bidAmount *big.Int,
		slashAmount *big.Int,
		rawTxs []string,
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
	LatestBaseFee() *big.Int
	NextBaseFee() *big.Int
}

type Transferer interface {
	Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error
}

type Simulator interface {
	Simulate(ctx context.Context, txRaw string, state sim.SimState) ([]*types.Log, bool, error)
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

type ExplorerSubmitter interface {
	Submit(ctx context.Context, tx *types.Transaction, from common.Address) error
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
	historicalTxns    *lru.Cache[common.Hash, struct{}]
	notifier          Notifier
	simulator         Simulator
	fastTrack         func(cmts []*bidderapiv1.Commitment, optedInSlot bool) bool
	bidTimeout        time.Duration
	timeoutMtx        sync.RWMutex
	receiptSignal     map[common.Hash][]chan struct{}
	receiptMtx        sync.Mutex
	metrics           *metrics
	explorerSubmitter ExplorerSubmitter
	logEncryptionKey  []byte
	bundleBids        bool
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
	explorerSubmitter ExplorerSubmitter,
	logEncryptionKey []byte,
	bundleBids bool,
	logger *slog.Logger,
) (*TxSender, error) {
	txnAttemptHistory, err := lru.New[common.Hash, *txnAttempt](1000)
	if err != nil {
		logger.Error("Failed to create transaction attempt history cache", "error", err)
		return nil, fmt.Errorf("failed to create transaction attempt history cache: %w", err)
	}

	historicalTxns, err := lru.New[common.Hash, struct{}](10000)
	if err != nil {
		logger.Error("Failed to create historical transactions cache", "error", err)
		return nil, fmt.Errorf("failed to create historical transactions cache: %w", err)
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
		historicalTxns:    historicalTxns,
		notifier:          notifier,
		simulator:         simulator,
		fastTrack:         noOpFastTrack,
		bidTimeout:        bidTimeout,
		receiptSignal:     make(map[common.Hash][]chan struct{}),
		metrics:           newMetrics(),
		explorerSubmitter: explorerSubmitter,
		logEncryptionKey:  logEncryptionKey,
		bundleBids:        bundleBids,
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
	if tx.Type < TxTypeRegular || tx.Type > TxTypeFastSwap {
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
	// Get backend (chain) nonce first
	backendNonce, err := t.blockTracker.AccountNonce(ctx, tx.Sender)
	if err != nil {
		return fmt.Errorf("failed to get backend nonce: %w", err)
	}

	// Get store nonce - returns (maxNonce, hasTxs)
	maxNonce, hasTxs := t.store.GetCurrentNonce(ctx, tx.Sender)

	var expectedNonce uint64
	if hasTxs {
		// Has transactions in store, next nonce is max + 1
		expectedNonce = maxNonce + 1
	} else {
		// No transactions in store, use chain nonce
		expectedNonce = backendNonce
	}

	// If chain has advanced beyond our tracking, use chain nonce
	if backendNonce > expectedNonce {
		expectedNonce = backendNonce
	}

	switch {
	case tx.Nonce() < expectedNonce:
		return ErrNonceTooLow
	case tx.Nonce() > expectedNonce:
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

	if t.txnAttemptHistory.Contains(tx.Hash()) {
		t.logger.Warn("Duplicate transaction enqueued", "hash", tx.Hash().Hex())
		return nil
	}

	if err := t.hasCorrectNonce(ctx, tx); err != nil {
		return err
	}

	if err := t.store.AddQueuedTransaction(ctx, tx); err != nil {
		return err
	}

	t.triggerSender()

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

// bundleMember is the slot of one transaction in a bundle bid batch. The
// batch leader writes the cost, result and error fields before it closes
// the batch done channel. Followers read their slot only after the done
// channel is closed. The fastTracked flag records that the member got
// enough commitments to be stored and signaled as preconfirmed while
// the bid status loop was still open.
type bundleMember struct {
	txn         *Transaction
	cancel      <-chan struct{}
	cost        *big.Int
	res         bidResult
	err         error
	fastTracked bool
	// storedCommitments is the number of commitments that the
	// fast-track write persisted. Later bundle commitments append in
	// memory only; the success gate compares this count with the live
	// count and stores the grown set again.
	storedCommitments int
}

// bundleBatch groups the transactions that share one bundle bid attempt.
// The first member is the leader. The leader performs the bid for the
// whole batch and closes the done channel when every member slot is set.
type bundleBatch struct {
	members []*bundleMember
	done    chan struct{}
}

func (b *bundleBatch) leader() *Transaction {
	return b.members[0].txn
}

func (b *bundleBatch) member(txn *Transaction) *bundleMember {
	for _, m := range b.members {
		if m.txn == txn {
			return m
		}
	}
	return nil
}

// bundleEligible reports whether a transaction can join a bundle bid.
// Only regular transactions without a position constraint are eligible.
// All other transactions follow the solo bid path. The freshness check
// (no previous commitments) runs later, after the inflight mark, because
// the commitments field is only safe to read for transactions that no
// other goroutine owns.
func (t *TxSender) bundleEligible(txn *Transaction) bool {
	return txn.Type == TxTypeRegular && txn.Constraint == nil
}

// acquireWorkerSlots takes count worker pool slots before any goroutine
// starts. It returns false if the context ended first. In that case it
// puts the slots it already took back into the pool, so the caller can
// release every inflight mark safely: no goroutine has started. The
// context check before each select makes a context that is already done
// win deterministically; a select with a free pool slot and a done
// context picks a case at random.
func (t *TxSender) acquireWorkerSlots(ctx context.Context, count int) bool {
	release := func(held int) {
		for ; held > 0; held-- {
			<-t.workerPool
		}
	}
	for held := 0; held < count; held++ {
		if ctx.Err() != nil {
			release(held)
			return false
		}
		select {
		case <-ctx.Done():
			release(held)
			return false
		case t.workerPool <- struct{}{}:
		}
	}
	return true
}

// dispatchMarked takes one worker pool slot and starts the worker
// goroutine for a transaction that is already marked inflight. It
// returns false if the context ended before a slot was free. In that
// case the caller must release the inflight marks of the members that
// did not get a goroutine.
func (t *TxSender) dispatchMarked(ctx context.Context, txn *Transaction, cancel <-chan struct{}, batch *bundleBatch) bool {
	if !t.acquireWorkerSlots(ctx, 1) {
		return false
	}
	t.startMarked(ctx, txn, cancel, batch)
	return true
}

// startMarked starts the worker goroutine for a transaction that is
// already marked inflight and already holds one worker pool slot. The
// goroutine puts the slot back when it ends.
func (t *TxSender) startMarked(ctx context.Context, txn *Transaction, cancel <-chan struct{}, batch *bundleBatch) {
	t.eg.Go(func() error {
		defer func() { <-t.workerPool }()
		defer t.triggerSender() // Trigger to reprocess after this transaction
		defer t.markCompleted(txn)

		t.logger.Info("Processing transaction", "sender", txn.Sender.Hex(), "type", txn.Type)
		if err := t.processTransaction(ctx, txn, cancel, batch); err != nil {
			t.logger.Error("Failed to process transaction", "sender", txn.Sender.Hex(), "error", err)
			txn.Status = TxStatusFailed
			txn.Details = err.Error()
			t.clearBlockAttemptHistory(txn, time.Now())
			// The terminal status write must survive the cancellation
			// that caused the failure. A write on the cancelled
			// context fails, the row keeps its old status, and the
			// queue never selects it again: the failure would stay
			// invisible.
			return t.store.StoreTransaction(context.WithoutCancel(ctx), txn, nil, nil)
		}
		return nil
	})
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

	var solo, eligible []*Transaction
	for _, txn := range txns {
		if t.bundleBids && t.bundleEligible(txn) {
			eligible = append(eligible, txn)
		} else {
			solo = append(solo, txn)
		}
	}
	if len(eligible) == 1 {
		// A single eligible transaction cannot form a bundle.
		solo = append(solo, eligible[0])
		eligible = nil
	}

	release := func(members []*bundleMember) {
		for _, m := range members {
			t.markCompleted(m.txn)
		}
	}

	for start := 0; start < len(eligible); start += maxBundleSize {
		end := min(start+maxBundleSize, len(eligible))
		// Build the roster under processMu. The historical check and the
		// inflight mark run here, in the same order the solo goroutine
		// uses. After this loop the roster is fixed.
		roster := make([]*bundleMember, 0, end-start)
		var preMarked []*bundleMember
		for _, txn := range eligible[start:end] {
			if t.historicalTxns.Contains(txn.Hash()) {
				t.logger.Warn("Transaction already processed historically, skipping", "hash", txn.Hash().Hex())
				continue
			}
			canExecute, cancel := t.markInflight(txn)
			if !canExecute {
				// Transaction is already being processed or sender has an inflight transaction
				continue
			}
			if len(txn.commitments) > 0 {
				// Not a fresh transaction. It keeps its inflight mark
				// and follows the solo path.
				preMarked = append(preMarked, &bundleMember{txn: txn, cancel: cancel})
				continue
			}
			roster = append(roster, &bundleMember{txn: txn, cancel: cancel})
		}
		if len(roster) == 1 {
			// A single member cannot form a bundle.
			preMarked = append(preMarked, roster[0])
			roster = nil
		}
		var batch *bundleBatch
		if len(roster) > 1 {
			batch = &bundleBatch{
				members: roster,
				done:    make(chan struct{}),
			}
		}
		// All-or-nothing dispatch for the batch: every roster member
		// gets its worker pool slot before any member goroutine starts.
		// A partial dispatch would let the leader bid for a member that
		// has no goroutine and whose inflight mark is released, and the
		// next tick would re-dispatch that member: a duplicate bid.
		// When the acquisition fails no goroutine has started, so every
		// inflight mark can be released and the next tick retries the
		// whole batch.
		if len(roster) > 0 {
			if !t.acquireWorkerSlots(ctx, len(roster)) {
				release(roster)
				release(preMarked)
				t.logger.Info("Context cancelled, stopping transaction processing")
				return
			}
			// The leader is started first so that the batch never runs
			// without its leader goroutine.
			for _, m := range roster {
				t.startMarked(ctx, m.txn, m.cancel, batch)
			}
		}
		for i, m := range preMarked {
			if !t.dispatchMarked(ctx, m.txn, m.cancel, nil) {
				release(preMarked[i:])
				t.logger.Info("Context cancelled, stopping transaction processing")
				return
			}
		}
	}

	for _, txn := range solo {
		txn := txn // capture range variable
		select {
		case <-ctx.Done():
			t.logger.Info("Context cancelled, stopping transaction processing")
			return
		case t.workerPool <- struct{}{}:
			t.eg.Go(func() error {
				defer func() { <-t.workerPool }()
				defer t.triggerSender() // Trigger to reprocess after this transaction

				if t.historicalTxns.Contains(txn.Hash()) {
					t.logger.Warn("Transaction already processed historically, skipping", "hash", txn.Hash().Hex())
					return nil
				}

				canExecute, cancel := t.markInflight(txn)
				if !canExecute {
					// Transaction is already being processed or sender has an inflight transaction
					return nil
				}
				defer t.markCompleted(txn)

				t.logger.Info("Processing transaction", "sender", txn.Sender.Hex(), "type", txn.Type)
				if err := t.processTransaction(ctx, txn, cancel, nil); err != nil {
					t.logger.Error("Failed to process transaction", "sender", txn.Sender.Hex(), "error", err)
					txn.Status = TxStatusFailed
					txn.Details = err.Error()
					t.clearBlockAttemptHistory(txn, time.Now())
					// The terminal status write must survive the
					// cancellation that caused the failure. A write on
					// the cancelled context fails, the row keeps its
					// old status, and the queue never selects it
					// again: the failure would stay invisible.
					return t.store.StoreTransaction(context.WithoutCancel(ctx), txn, nil, nil)
				}
				return nil
			})
		}
	}
}

func (t *TxSender) processTransaction(ctx context.Context, txn *Transaction, cancel <-chan struct{}, batch *bundleBatch) error {
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
		switch {
		case batch == nil:
			result, err = t.sendBid(ctx, txn)
		case batch.leader() == txn:
			// The leader sends one bid for the whole batch and fills
			// every member slot. The batch is spent after the first
			// iteration; later iterations use the solo path.
			result, err = t.sendBundleBid(ctx, batch)
			batch = nil
		default:
			// Followers wait for the leader to close the done channel.
			// The wait is bounded by the leader's bid timeout. Cancel
			// and context signals are handled by the select below.
			<-batch.done
			m := batch.member(txn)
			result, err = m.res, m.err
			batch = nil
		}
		switch {
		case err != nil:
			if errors.Is(err, errBundleAborted) {
				// The bundle never placed a bid for this member. The
				// attempt state is already rolled back, so the next
				// iteration runs a solo bid with first-attempt
				// semantics.
				logger.Warn("Bundle aborted before the bid, retrying solo")
				retryTicker.Reset(defaultRetryDelay)
			} else if retryErr, ok := err.(*errRetry); ok {
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
				if result.fastTracked {
					// The batch leader already stored this member's
					// preconfirmation, signaled its receipt and
					// submitted it to the explorer when the member
					// was fast-tracked. A repeat store here can run
					// with a context that is already done, and its
					// failure would flip the stored preconfirmation
					// to a failed state. So the member skips that
					// work and only resets the retry ticker below.
					logger.Info(
						"Fast-tracked bundle member keeps its stored preconfirmation",
						"blockNumber", result.blockNumber,
						"bidAmount", result.bidAmount.String(),
					)
					if len(txn.commitments) > result.fastTrackCommitments {
						// More commitments arrived after the
						// fast-track write and exist in memory only.
						// Store the grown set, so the receipt API
						// shows every commitment, like the solo path.
						// The write must survive a done context, and
						// a failure only logs a warning: it must
						// never flip the stored preconfirmation to a
						// failed state.
						if err := t.store.StoreTransaction(
							context.WithoutCancel(ctx), txn, txn.commitments, txn.logs,
						); err != nil {
							logger.Warn(
								"Failed to store late bundle commitments",
								"error", err,
							)
						}
					}
				} else {
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
					if err := t.explorerSubmitter.Submit(
						ctx,
						txn.Transaction,
						txn.Sender,
					); err != nil {
						t.logger.Error("Failed to submit tx to explorer", "error", err)
					}
				}
			}
			retryTicker.Reset(result.timeUntillNextBlock + 1*time.Second)
		default:
			warnFields := []any{
				"noOfProviders", txn.noOfProviders,
				"noOfCommitments", len(txn.commitments),
				"blockNumber", result.blockNumber,
				"bidAmount", result.bidAmount.String(),
			}
			if encrypted := encryptForLog(t.logEncryptionKey, txn.Raw); encrypted != "" {
				warnFields = append(warnFields, "encryptedRawTx", encrypted)
			}
			logger.Warn("Not all builders committed to the bid", warnFields...)
			retryTicker.Reset(defaultRetryDelay)
		}
		select {
		case <-ctx.Done():
			if txn.Status == TxStatusPreConfirmed &&
				txn.Type != TxTypeDeposit && txn.Type != TxTypeInstantBridge {
				// The preconfirmation for this transaction is stored
				// and signaled. A shutdown must not revert it to a
				// failed state, so the loop ends without an error.
				// This holds for the solo path and for bundle members
				// in the same way. Only TxTypeDeposit and
				// TxTypeInstantBridge are excluded: they run their
				// post-loop obligation, the AddBalance credit or the
				// bridge Transfer, after this loop, and the queue only
				// selects pending transactions again. A silent success
				// here would skip that step forever. Those two types
				// return the context error, become visibly failed, and
				// the user can retry them. Every other type, such as a
				// preconfirmed fastswap, has no post-loop work and
				// keeps its preconfirmation.
				return nil
			}
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
	case TxTypeDeposit:
		if err := t.store.AddBalance(ctx, txn.Sender, txn.Value()); err != nil {
			logger.Error("Failed to add balance for sender", "error", err)
			return fmt.Errorf("failed to add balance for sender: %w", err)
		}
	case TxTypeInstantBridge:
		// deduct 5% as fee
		fee := new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(5)), big.NewInt(100))
		amountToBridge := new(big.Int).Sub(amount, fee)
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

// errBundleAborted tells a bundle member that the bundle bid was
// aborted before any bid went out. The member did not use up a real
// bid attempt, so its next solo bid must run with first-attempt
// semantics: it simulates again and sets noOfProviders from the live
// provider set. This keeps the zero-commitment outcome of the solo
// bid from passing the noOfProviders == len(commitments) success gate
// as 0 == 0.
var errBundleAborted = errors.New("bundle aborted before bid placement")

type bidResult struct {
	startTime           time.Time
	timeUntillNextBlock time.Duration
	blockNumber         uint64
	optedInSlot         bool
	bidAmount           *big.Int
	// fastTracked marks a bundle member result whose preconfirmation
	// the batch leader already stored and signaled per member. The
	// success gate must not repeat that work: a repeat store can run
	// after the context ends, and its failure would flip the stored
	// preconfirmation to a failed state. The solo path never sets it.
	fastTracked bool
	// fastTrackCommitments is the number of commitments that the
	// fast-track write persisted for this member. When more
	// commitments arrive before the bid loop ends, the success gate
	// stores the grown set again, so the receipt API shows every
	// commitment, like the solo path.
	fastTrackCommitments int
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

	feePerGas := effectiveFeePerGas(txn.Transaction)
	nextBaseFee := t.blockTracker.NextBaseFee()
	latestBaseFee := t.blockTracker.LatestBaseFee()
	if nextBaseFee.Sign() == 0 {
		nextBaseFee = latestBaseFee
	}

	if nextBaseFee.Sign() > 0 && feePerGas.Cmp(nextBaseFee) < 0 {
		logger.Warn(
			"Fee per gas too low for next block",
			"feePerGas", feePerGas.String(),
			"nextBaseFee", nextBaseFee.String(),
		)
		return bidResult{}, &errRetry{
			err: fmt.Errorf(
				"fee per gas too low for next block: %s min %s",
				feePerGas.String(),
				nextBaseFee.String(),
			),
			retryAfter: timeUntilNextBlock,
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
		if txn.Value().Cmp(cost) < 0 {
			logger.Error(
				"Instant bridge amount is less than price of bridge",
				"bridge", txn.Value().String(),
				"price", cost.String(),
			)
			return bidResult{}, fmt.Errorf(
				"instant bridge amount is less than price of bridge: %s, bridge: %s, price: %s",
				txn.Sender.Hex(),
				txn.Value().String(),
				cost.String(),
			)
		}
		slashAmount = new(big.Int).Set(txn.Value())
	case TxTypeFastSwap:
		// FastSwap executor transactions skip balance check - RPC bidder pays for bids
		logger.Debug("FastSwap transaction - skipping balance check", "sender", txn.Sender.Hex())
	}

	state := sim.Latest
	if latestBaseFee.Sign() > 0 && feePerGas.Cmp(latestBaseFee) < 0 {
		state = sim.Pending
	}

	// Re-simulate on new blocks or when the previous simulation failed.
	// This prevents bidding with stale data after a sim failure while
	// avoiding redundant sim calls when sim already passed within the same block.
	if !isRetry || txn.simFailed {
		logs, isSwap, err := t.simulator.Simulate(ctx, txn.Raw, state)
		if err != nil {
			txn.simFailed = true
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
		txn.simFailed = false
		if !isRetry {
			txn.logs = logs
			txn.isSwap = isSwap
			// We could have already made a attempt on the previous block but the block
			// update hasn't happened yet. This means that the bid might fail, but
			// we should retain the previous commitments. Only clear if we get new
			// commitments for the new block.
		}
	}

	// The success gate compares noOfProviders with the commitment
	// count. A retry that never set noOfProviders must not let a
	// zero-commitment outcome pass that gate as zero equals zero. So
	// the provider count comes from the live provider set on every
	// first attempt and on every retry that has no provider count yet.
	if !isRetry || txn.noOfProviders == 0 {
		providers, err := t.bidder.ConnectedProviders(ctx)
		if err != nil {
			logger.Error("Failed to get connected providers", "error", err)
			return bidResult{}, fmt.Errorf("failed to get connected providers: %w", err)
		}
		txn.noOfProviders = len(providers)
		t.metrics.connectedProviders.Set(float64(len(providers)))
	}

	bidStart := time.Now()
	bidC, err := t.bidder.Bid(
		cctx,
		cost,
		slashAmount,
		[]string{strings.TrimPrefix(txn.Raw, "0x")},
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
					if err := t.explorerSubmitter.Submit(
						ctx,
						txn.Transaction,
						txn.Sender,
					); err != nil {
						t.logger.Error("Failed to submit fast-tracked tx to explorer", "error", err)
					}
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

// sendBundleBid sends one bid that carries the raw transactions of all
// batch members. It runs in the goroutine of the batch leader. It writes
// a result or an error into the slot of every member and closes the batch
// done channel when it returns. The return values are the slot values of
// the leader. A member that fails a pre-check drops out of the bundle and
// receives the same error the solo path produces for that failure.
func (t *TxSender) sendBundleBid(ctx context.Context, batch *bundleBatch) (bidResult, error) {
	defer close(batch.done)

	start := time.Now()
	leader := batch.members[0]
	logger := t.logger.With(
		"leaderTransactionHash", leader.txn.Hash().Hex(),
		"bundleSize", len(batch.members),
	)

	distribute := func(members []*bundleMember, err error) {
		for _, m := range members {
			m.err = err
		}
	}

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
		distribute(batch.members, &errRetry{
			err:        fmt.Errorf("failed to get next block number: %w", err),
			retryAfter: time.Second,
		})
		return leader.res, leader.err
	}
	logger.Debug("Next block info", "bidBlockNo", bidBlockNo, "timeUntilNextBlock", timeUntilNextBlock)

	if timeUntilNextBlock <= 500*time.Millisecond {
		logger.Warn("Next block time is too short, skipping bid", "timeUntilNextBlock", timeUntilNextBlock)
		distribute(batch.members, &errRetry{
			err:        fmt.Errorf("next block time is too short: %s", timeUntilNextBlock),
			retryAfter: defaultRetryDelay,
		})
		return leader.res, leader.err
	}

	prices := t.pricer.EstimatePrice(ctx)

	// Allow for certain level of tolerance w.r.t timestamps
	optedInSlot := math.Abs(float64(timeToOptIn)-float64(timeUntilNextBlock.Seconds())) < float64(blockTime/3)

	cctx, cancel := context.WithTimeout(ctx, t.getBidTimeout())
	defer cancel()

	nextBaseFee := t.blockTracker.NextBaseFee()
	latestBaseFee := t.blockTracker.LatestBaseFee()
	if nextBaseFee.Sign() == 0 {
		nextBaseFee = latestBaseFee
	}

	// Per-member pre-checks. A member that fails a check keeps its own
	// error and is left out of the bundle. This filters transactions
	// that would make the providers reject the full bundle.
	active := make([]*bundleMember, 0, len(batch.members))
	for _, m := range batch.members {
		mlogger := t.logger.With(
			"transactionHash", m.txn.Hash().Hex(),
			"sender", m.txn.Sender.Hex(),
			"type", m.txn.Type,
		)

		// Every drop below this price calculation rolls the recorded
		// block attempt back: the member leaves the bundle before any
		// bid goes out, so its solo retry must run with first-attempt
		// semantics and set noOfProviders from the live provider set.
		cost, _, err := t.calculatePriceForNextBlock(m.txn, bidBlockNo, prices, optedInSlot)
		if err != nil {
			mlogger.Error("Failed to calculate price for next block", "error", err)
			if errors.Is(err, ErrTimeoutExceeded) || errors.Is(err, ErrMaxAttemptsPerBlockExceeded) {
				// We propagate these errors as is. A timeout records
				// no attempt, so there is nothing to roll back.
				m.err = err
			} else {
				// The failed price lookup recorded an attempt for
				// this block, but no bid goes out for this member.
				t.rollbackBlockAttempt(m.txn, bidBlockNo)
				m.err = &errRetry{
					err:        fmt.Errorf("failed to calculate price: %w", err),
					retryAfter: time.Second,
				}
			}
			continue
		}

		feePerGas := effectiveFeePerGas(m.txn.Transaction)
		if nextBaseFee.Sign() > 0 && feePerGas.Cmp(nextBaseFee) < 0 {
			mlogger.Warn(
				"Fee per gas too low for next block",
				"feePerGas", feePerGas.String(),
				"nextBaseFee", nextBaseFee.String(),
			)
			t.rollbackBlockAttempt(m.txn, bidBlockNo)
			m.err = &errRetry{
				err: fmt.Errorf(
					"fee per gas too low for next block: %s min %s",
					feePerGas.String(),
					nextBaseFee.String(),
				),
				retryAfter: timeUntilNextBlock,
			}
			continue
		}

		// Bundle members are always regular transactions, so the balance
		// check mirrors the regular solo path.
		if !t.store.HasBalance(ctx, m.txn.Sender, cost) {
			mlogger.Error("Insufficient balance for sender")
			t.rollbackBlockAttempt(m.txn, bidBlockNo)
			m.err = fmt.Errorf("insufficient balance for sender: %s", m.txn.Sender.Hex())
			continue
		}

		state := sim.Latest
		if latestBaseFee.Sign() > 0 && feePerGas.Cmp(latestBaseFee) < 0 {
			state = sim.Pending
		}

		logs, isSwap, err := t.simulator.Simulate(ctx, m.txn.Raw, state)
		if err != nil {
			m.txn.simFailed = true
			mlogger.Error("Failed to simulate transaction", "error", err, "blockNumber", bidBlockNo)
			t.rollbackBlockAttempt(m.txn, bidBlockNo)
			m.err = fmt.Errorf("failed to simulate transaction: %w", err)
			continue
		}
		m.txn.simFailed = false
		m.txn.logs = logs
		m.txn.isSwap = isSwap
		m.cost = cost
		active = append(active, m)
	}

	if len(active) < 2 {
		logger.Info("Not enough members left for a bundle bid", "active", len(active))
		// The per-member price calculation consumed one block attempt
		// for each survivor, but no bid went out. Roll the attempt
		// back so the solo retry runs as a first attempt.
		for _, m := range active {
			t.rollbackBlockAttempt(m.txn, bidBlockNo)
		}
		distribute(active, errBundleAborted)
		return leader.res, leader.err
	}

	providers, err := t.bidder.ConnectedProviders(ctx)
	if err != nil {
		logger.Error("Failed to get connected providers", "error", err)
		distribute(active, fmt.Errorf("failed to get connected providers: %w", err))
		return leader.res, leader.err
	}
	t.metrics.connectedProviders.Set(float64(len(providers)))

	// The bid amount is the sum of the per-member costs. Each member
	// hash is listed as revertable, which matches the solo path where
	// each bid lists its own transaction hash.
	totalCost := big.NewInt(0)
	rawTxs := make([]string, 0, len(active))
	revertingTxHashes := make([]string, 0, len(active))
	for _, m := range active {
		m.txn.noOfProviders = len(providers)
		totalCost = new(big.Int).Add(totalCost, m.cost)
		rawTxs = append(rawTxs, strings.TrimPrefix(m.txn.Raw, "0x"))
		revertingTxHashes = append(revertingTxHashes, m.txn.Hash().Hex())
	}

	bidStart := time.Now()
	bidC, err := t.bidder.Bid(
		cctx,
		totalCost,
		big.NewInt(0),
		rawTxs,
		&bidder.BidOpts{
			WaitForOptIn:      false,
			BlockNumber:       bidBlockNo,
			RevertingTxHashes: revertingTxHashes,
			DecayDuration:     t.getBidTimeout() * 2,
		},
	)
	if err != nil {
		logger.Error("Failed to place bundle bid", "error", err)
		// The members retry on the solo path after a short delay.
		distribute(active, &errRetry{
			err:        fmt.Errorf("failed to place bundle bid: %w", err),
			retryAfter: defaultRetryDelay,
		})
		return leader.res, leader.err
	}

	commitmentCount := 0
BID_LOOP:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled while waiting for bid status")
			// A fast-tracked member already has its preconfirmation
			// stored and signaled. It gets its success result, the
			// same result the solo path returns for a fast-tracked
			// transaction. Only members that are not resolved get
			// the context error.
			for _, m := range active {
				if m.fastTracked {
					m.res = bidResult{
						bidAmount:            m.cost,
						blockNumber:          bidBlockNo,
						startTime:            start,
						timeUntillNextBlock:  timeUntilNextBlock,
						optedInSlot:          optedInSlot,
						fastTracked:          true,
						fastTrackCommitments: m.storedCommitments,
					}
					continue
				}
				m.err = ctx.Err()
			}
			return leader.res, leader.err
		case bidStatus, more := <-bidC:
			if !more {
				logger.Info("Bid channel closed, no more bid statuses")
				break BID_LOOP
			}
			switch bidStatus.Type {
			case bidder.BidStatusCommitment:
				cmt := bidStatus.Arg.(*bidderapiv1.Commitment)
				commitmentCount++
				for _, m := range active {
					t.appendMemberCommitment(ctx, m, cmt, bidBlockNo, optedInSlot, start)
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
		"Bundle bid operation complete",
		"noOfProviders", len(providers),
		"noOfCommitments", commitmentCount,
		"blockNumber", bidBlockNo,
		"optedInSlot", optedInSlot,
	)

	for _, m := range active {
		if len(m.txn.commitments) > 0 && m.txn.isSwap {
			if err := t.backrunner.Backrun(ctx, m.txn.Raw, m.txn.commitments); err != nil {
				logger.Error("Failed to backrun transaction", "error", err)
			}
			logger.Info("Backrun operation initiated for transaction", "hash", m.txn.Hash().Hex())
		}
		m.res = bidResult{
			bidAmount:            m.cost,
			blockNumber:          bidBlockNo,
			startTime:            start,
			timeUntillNextBlock:  timeUntilNextBlock,
			optedInSlot:          optedInSlot,
			fastTracked:          m.fastTracked,
			fastTrackCommitments: m.storedCommitments,
		}
	}

	return leader.res, leader.err
}

// appendMemberCommitment stores a copy of a bundle commitment on one
// member with the bid amount replaced by the member's own share. The sum
// of the member shares equals the aggregate bundle bid amount.
func (t *TxSender) appendMemberCommitment(
	ctx context.Context,
	m *bundleMember,
	cmt *bidderapiv1.Commitment,
	bidBlockNo uint64,
	optedInSlot bool,
	start time.Time,
) {
	txn := m.txn
	if len(txn.commitments) > 0 {
		if txn.commitments[0].BlockNumber != int64(bidBlockNo) {
			txn.commitments = nil // clear previous commitments for new block
		}
	}
	memberCmt, ok := proto.Clone(cmt).(*bidderapiv1.Commitment)
	if !ok {
		t.logger.Error("Failed to clone commitment", "transactionHash", txn.Hash().Hex())
		return
	}
	memberCmt.BidAmount = m.cost.String()
	txn.commitments = append(txn.commitments, memberCmt)
	if !t.fastTrack(txn.commitments, optedInSlot) {
		return
	}
	m.fastTracked = true
	if txn.Status != TxStatusPreConfirmed {
		txn.Status = TxStatusPreConfirmed
		txn.BlockNumber = int64(bidBlockNo)
		t.logger.Info(
			"Transaction fast-tracked based on commitments",
			"transactionHash", txn.Hash().Hex(),
			"blockNumber", bidBlockNo,
			"bidAmount", m.cost.String(),
		)
		if err := t.store.StoreTransaction(ctx, txn, txn.commitments, txn.logs); err != nil {
			t.logger.Error("Failed to store fast-tracked transaction", "error", err)
		}
		// Record how many commitments the fast-track write covered.
		// The success gate stores the set again when it grew after
		// this write.
		m.storedCommitments = len(txn.commitments)
		t.signalReceiptAvailable(txn.Hash())
		t.metrics.timeToFirstPreconfirmation.Observe(float64(time.Since(start).Milliseconds()))
		if err := t.explorerSubmitter.Submit(
			ctx,
			txn.Transaction,
			txn.Sender,
		); err != nil {
			t.logger.Error("Failed to submit fast-tracked tx to explorer", "error", err)
		}
	}
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
	if txn.Type == TxTypeFastSwap {
		confidence = defaultSwapConfidence
	}
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

	if txn.Constraint != nil && txn.Type == TxTypeFastSwap {
		if isRetry {
			confidence = topOfBlockSubsequentAttempts
		} else {
			confidence = topOfBlockSwapConfidence
		}
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

// rollbackBlockAttempt undoes one recorded attempt for a block. The
// bundle pre-check records an attempt through calculatePriceForNextBlock
// before the bundle bid goes out. When the bundle aborts before the bid,
// that attempt did not happen, so the next solo bid for the same block
// must see first-attempt state. The entry is removed when no attempts
// remain, which keeps clearBlockAttemptHistory safe: it reads the first
// element of a non-empty attempts list.
func (t *TxSender) rollbackBlockAttempt(txn *Transaction, bidBlockNo uint64) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		return
	}

	for i := len(attempts.attempts) - 1; i >= 0; i-- {
		if attempts.attempts[i].blockNumber != bidBlockNo {
			continue
		}
		attempts.attempts[i].attempts--
		if attempts.attempts[i].attempts <= 0 {
			attempts.attempts = append(attempts.attempts[:i], attempts.attempts[i+1:]...)
		}
		break
	}

	if len(attempts.attempts) == 0 {
		_ = t.txnAttemptHistory.Remove(txn.Hash())
		return
	}
	_ = t.txnAttemptHistory.Add(txn.Hash(), attempts)
}

func (t *TxSender) clearBlockAttemptHistory(txn *Transaction, endTime time.Time) {
	attempts, found := t.txnAttemptHistory.Get(txn.Hash())
	if !found {
		// A rollback of the only recorded attempt removes the entry,
		// so a terminal drop can land here without one. The user
		// notification and the historical record must still happen:
		// the historical record keeps the queue from processing the
		// same hash again. Only the metrics that read the entry are
		// skipped.
		t.notifier.NotifyTransactionStatus(txn, 0, 0, 0)
		_ = t.historicalTxns.Add(txn.Hash(), struct{}{})
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
	_ = t.historicalTxns.Add(txn.Hash(), struct{}{})
}
