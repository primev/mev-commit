package updater

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/lib/pq"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type SettlementType string

type TxMetadata struct {
	PosInBlock int
	Succeeded  bool
}

const (
	SettlementTypeReward SettlementType = "reward"
	SettlementTypeSlash  SettlementType = "slash"
)

type Winner struct {
	Winner []byte
	Window int64
}

type Settlement struct {
	CommitmentIdx   []byte
	TxHash          string
	BlockNum        int64
	Builder         []byte
	Amount          *big.Int
	BidID           []byte
	Type            SettlementType
	DecayPercentage int64
}

type WinnerRegister interface {
	AddEncryptedCommitment(
		ctx context.Context,
		commitmentIdx []byte,
		committer []byte,
		commitmentHash []byte,
		commitmentSignature []byte,
		dispatchTimestamp uint64,
	) error
	IsSettled(ctx context.Context, commitmentIdx []byte) (bool, error)
	GetWinner(ctx context.Context, blockNum int64) (Winner, error)
	AddSettlement(
		ctx context.Context,
		commitmentIdx []byte,
		txHash string,
		blockNum int64,
		amount *big.Int,
		builder []byte,
		bidID []byte,
		settlementType SettlementType,
		decayPercentage int64,
		window int64,
		postingTxnHash []byte,
		nonce uint64,
	) error
}

type Oracle interface {
	ProcessBuilderCommitmentForBlockNumber(
		commitmentIdx [32]byte,
		blockNum *big.Int,
		builder common.Address,
		isSlash bool,
		residualDecay *big.Int,
	) (*types.Transaction, error)
}

type EVMClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type Updater struct {
	logger         *slog.Logger
	l1Client       EVMClient
	winnerRegister WinnerRegister
	oracle         Oracle
	evtMgr         events.EventManager
	l1BlockCache   *lru.Cache[uint64, map[string]TxMetadata]
	unopenedCmts   chan *preconf.PreconfmanagerUnopenedCommitmentStored
	openedCmts     chan *preconf.PreconfmanagerOpenedCommitmentStored
	currentWindow  atomic.Int64
	metrics        *metrics
	receiptBatcher txmonitor.BatchReceiptGetter
}

func NewUpdater(
	logger *slog.Logger,
	l1Client EVMClient,
	winnerRegister WinnerRegister,
	evtMgr events.EventManager,
	oracle Oracle,
	receiptBatcher txmonitor.BatchReceiptGetter,
) (*Updater, error) {
	logger.Info("creating new updater instance")
	l1BlockCache, err := lru.New[uint64, map[string]TxMetadata](1024)
	if err != nil {
		logger.Error("failed to create L1 block cache", "error", err)
		return nil, fmt.Errorf("failed to create L1 block cache: %w", err)
	}
	logger.Info("successfully created L1 block cache with size 1024")

	updater := &Updater{
		logger:         logger,
		l1Client:       l1Client,
		l1BlockCache:   l1BlockCache,
		winnerRegister: winnerRegister,
		evtMgr:         evtMgr,
		oracle:         oracle,
		receiptBatcher: receiptBatcher,
		metrics:        newMetrics(),
		openedCmts:     make(chan *preconf.PreconfmanagerOpenedCommitmentStored),
		unopenedCmts:   make(chan *preconf.PreconfmanagerUnopenedCommitmentStored),
	}
	logger.Info("successfully created updater instance",
		"l1Client", fmt.Sprintf("%T", l1Client),
		"winnerRegister", fmt.Sprintf("%T", winnerRegister),
		"oracle", fmt.Sprintf("%T", oracle),
		"receiptBatcher", fmt.Sprintf("%T", receiptBatcher))
	return updater, nil
}

func (u *Updater) Metrics() []prometheus.Collector {
	u.logger.Debug("retrieving metrics collectors")
	collectors := u.metrics.Collectors()
	u.logger.Debug("returning metrics collectors", "count", len(collectors))
	return collectors
}

func (u *Updater) Start(ctx context.Context) <-chan struct{} {
	u.logger.Info("starting updater")
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)
	u.logger.Debug("created error group with context")

	ev1 := events.NewEventHandler(
		"UnopenedCommitmentStored",
		func(update *preconf.PreconfmanagerUnopenedCommitmentStored) {
			u.logger.Debug("handling unopened commitment event",
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]))
			select {
			case <-egCtx.Done():
				u.logger.Debug("context cancelled while handling unopened commitment")
			case u.unopenedCmts <- update:
				u.logger.Debug("successfully sent unopened commitment to channel")
			}
		},
	)

	ev2 := events.NewEventHandler(
		"OpenedCommitmentStored",
		func(update *preconf.PreconfmanagerOpenedCommitmentStored) {
			u.logger.Debug("handling opened commitment event",
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]))
			select {
			case <-egCtx.Done():
				u.logger.Debug("context cancelled while handling opened commitment")
			case u.openedCmts <- update:
				u.logger.Debug("successfully sent opened commitment to channel")
			}
		},
	)

	ev3 := events.NewEventHandler(
		"NewWindow",
		func(update *blocktracker.BlocktrackerNewWindow) {
			oldWindow := u.currentWindow.Load()
			u.currentWindow.Store(update.Window.Int64())
			u.logger.Info("updated current window",
				"oldWindow", oldWindow,
				"newWindow", update.Window.Int64())
		},
	)

	u.logger.Info("subscribing to events")
	sub, err := u.evtMgr.Subscribe(ev1, ev2, ev3)
	if err != nil {
		u.logger.Error("failed to subscribe to events", "error", err)
		close(doneChan)
		return doneChan
	}
	u.logger.Info("successfully subscribed to events")

	eg.Go(func() error {
		u.logger.Debug("starting subscription error handler goroutine")
		defer sub.Unsubscribe()
		select {
		case <-egCtx.Done():
			u.logger.Debug("context cancelled, exiting subscription error handler")
			return nil
		case err := <-sub.Err():
			u.logger.Error("subscription error received", "error", err)
			return err
		}
	})

	eg.Go(func() error {
		u.logger.Debug("starting unopened commitments handler goroutine")
		for {
			select {
			case <-egCtx.Done():
				u.logger.Debug("context cancelled, exiting unopened commitments handler")
				return nil
			case ec := <-u.unopenedCmts:
				u.logger.Info("processing unopened commitment",
					"commitmentIdx", common.Bytes2Hex(ec.CommitmentIndex[:]))
				if err := u.handleEncryptedCommitment(egCtx, ec); err != nil {
					u.logger.Error("failed to handle encrypted commitment", "error", err)
					return err
				}
			}
		}
	})

	eg.Go(func() error {
		u.logger.Debug("starting opened commitments handler goroutine")
		for {
			select {
			case <-egCtx.Done():
				u.logger.Debug("context cancelled, exiting opened commitments handler")
				return nil
			case oc := <-u.openedCmts:
				u.logger.Info("processing opened commitment",
					"commitmentIdx", common.Bytes2Hex(oc.CommitmentIndex[:]))
				if err := u.handleOpenedCommitment(egCtx, oc); err != nil {
					u.logger.Error("failed to handle opened commitment", "error", err)
					return err
				}
			}
		}
	})

	go func() {
		u.logger.Debug("starting main error group handler goroutine")
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			u.logger.Error("updater failed, exiting", "error", err)
		}
		u.logger.Info("updater stopped")
	}()

	u.logger.Info("updater started successfully")

	return doneChan
}

func (u *Updater) handleEncryptedCommitment(
	ctx context.Context,
	update *preconf.PreconfmanagerUnopenedCommitmentStored,
) error {
	u.logger.Info("handling encrypted commitment",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"committer", update.Committer.Hex(),
		"dispatchTimestamp", update.DispatchTimestamp)

	err := u.winnerRegister.AddEncryptedCommitment(
		ctx,
		update.CommitmentIndex[:],
		update.Committer.Bytes(),
		update.CommitmentDigest[:],
		update.CommitmentSignature,
		update.DispatchTimestamp,
	)
	if err != nil {
		// ignore duplicate private key constraint
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			u.logger.Warn(
				"encrypted commitment already exists",
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
				"error", err,
			)
			return nil
		}
		u.logger.Error(
			"failed to add encrypted commitment",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}
	u.metrics.EncryptedCommitmentsCount.Inc()
	u.logger.Info(
		"successfully added encrypted commitment",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"dispatch timestamp", update.DispatchTimestamp,
		"committer", update.Committer.Hex(),
	)
	return nil
}

func (u *Updater) handleOpenedCommitment(
	ctx context.Context,
	update *preconf.PreconfmanagerOpenedCommitmentStored,
) error {
	u.logger.Info("received opened commitment",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"committer", update.Committer.Hex(),
		"blockNumber", update.BlockNumber,
		"bid", update.BidAmt.String())
	u.metrics.CommitmentsReceivedCount.Inc()

	alreadySettled, err := u.winnerRegister.IsSettled(ctx, update.CommitmentIndex[:])
	if err != nil {
		u.logger.Error(
			"failed to check if commitment is settled",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}
	if alreadySettled {
		// both bidders and providers could open commitments, so this could
		// be a duplicate event
		u.logger.Info(
			"duplicate open commitment event detected",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"committer", update.Committer.Hex(),
		)
		u.metrics.DuplicateCommitmentsCount.Inc()
		return nil
	}

	u.logger.Debug("checking winner for block number", "blockNumber", update.BlockNumber)
	winner, err := u.winnerRegister.GetWinner(ctx, int64(update.BlockNumber))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Warn("winner not found for block",
				"blockNumber", update.BlockNumber,
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]))
			u.metrics.NoWinnerCount.Inc()
			return nil
		}
		u.logger.Error(
			"failed to get winner",
			"blockNumber", update.BlockNumber,
			"error", err,
		)
		return err
	}
	u.logger.Info("retrieved winner information",
		"winner", common.Bytes2Hex(winner.Winner),
		"window", winner.Window)

	if u.currentWindow.Load() > 2 && winner.Window < u.currentWindow.Load()-2 {
		u.logger.Info(
			"commitment is too old to process",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"winner", common.Bytes2Hex(winner.Winner),
			"winnerWindow", winner.Window,
			"currentWindow", u.currentWindow.Load(),
		)
		u.metrics.CommitmentsTooOldCount.Inc()
		return nil
	}

	if common.BytesToAddress(winner.Winner).Cmp(update.Committer) != 0 {
		// The winner is not the committer of the commitment
		u.logger.Info(
			"winner does not match committer - skipping commitment",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"winner", common.Bytes2Hex(winner.Winner),
			"committer", update.Committer.Hex(),
			"blockNumber", update.BlockNumber,
		)
		return nil
	}

	u.logger.Debug("retrieving L1 transactions", "blockNumber", update.BlockNumber)
	txns, err := u.getL1Txns(ctx, update.BlockNumber)
	if err != nil {
		u.logger.Error(
			"failed to get L1 transactions",
			"blockNumber", update.BlockNumber,
			"error", err,
		)
		return err
	}
	u.logger.Debug("successfully retrieved L1 transactions",
		"blockNumber", update.BlockNumber,
		"txCount", len(txns))

	// Compute the decay percentage
	decayPercentage := u.computeDecayPercentage(
		update.DecayStartTimeStamp,
		update.DecayEndTimeStamp,
		update.DispatchTimestamp,
	)
	u.logger.Info("computed decay percentage",
		"decayPercentage", decayPercentage,
		"startTimestamp", update.DecayStartTimeStamp,
		"endTimestamp", update.DecayEndTimeStamp,
		"dispatchTimestamp", update.DispatchTimestamp)

	commitmentTxnHashes := strings.Split(update.TxnHash, ",")
	u.logger.Debug("processing commitment transaction hashes",
		"txnHashes", commitmentTxnHashes,
		"count", len(commitmentTxnHashes))

	revertableTxns := strings.Split(update.RevertingTxHashes, ",")
	u.logger.Debug("processing revertable transactions",
		"revertableTxns", revertableTxns,
		"count", len(revertableTxns))

	// Create a map for revertable transactions
	revertableTxnsMap := make(map[string]bool)
	for _, txn := range revertableTxns {
		revertableTxnsMap[txn] = true
	}
	u.logger.Debug("created revertable transactions map",
		"mapSize", len(revertableTxnsMap))

	// Ensure Bundle is atomic and present in the block
	for i := 0; i < len(commitmentTxnHashes); i++ {
		txnDetails, found := txns[commitmentTxnHashes[i]]
		if !found || txnDetails.PosInBlock != (txns[commitmentTxnHashes[0]].PosInBlock)+i || (!txnDetails.Succeeded && !revertableTxnsMap[commitmentTxnHashes[i]]) {
			u.logger.Info(
				"bundle validation failed - initiating slash",
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
				"txnHash", commitmentTxnHashes[i],
				"blockNumber", update.BlockNumber,
				"found", found,
				"posInBlock", txnDetails.PosInBlock,
				"succeeded", txnDetails.Succeeded,
				"expectedPosInBlock", txns[commitmentTxnHashes[0]].PosInBlock+i,
				"isRevertible", revertableTxnsMap[commitmentTxnHashes[i]],
			)

			// The committer did not include the transactions in the block
			// correctly, so this is a slash to be processed
			return u.settle(
				ctx,
				update,
				SettlementTypeSlash,
				decayPercentage,
				winner.Window,
			)
		}
		u.logger.Debug("transaction in bundle validated successfully",
			"txnHash", commitmentTxnHashes[i],
			"posInBlock", txnDetails.PosInBlock,
			"succeeded", txnDetails.Succeeded)
	}

	u.logger.Info("bundle validation successful - processing reward",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"blockNumber", update.BlockNumber)
	return u.settle(
		ctx,
		update,
		SettlementTypeReward,
		decayPercentage,
		winner.Window,
	)
}

func (u *Updater) settle(
	ctx context.Context,
	update *preconf.PreconfmanagerOpenedCommitmentStored,
	settlementType SettlementType,
	decayPercentage int64,
	window int64,
) error {
	u.logger.Info("initiating settlement process",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"settlementType", settlementType,
		"decayPercentage", decayPercentage,
		"window", window)

	commitmentPostingTxn, err := u.oracle.ProcessBuilderCommitmentForBlockNumber(
		update.CommitmentIndex,
		big.NewInt(0).SetUint64(update.BlockNumber),
		update.Committer,
		settlementType == SettlementTypeSlash,
		big.NewInt(decayPercentage),
	)
	if err != nil {
		u.logger.Error(
			"failed to process commitment with oracle",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}
	u.logger.Info(
		"commitment processed successfully by oracle",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"blockNumber", update.BlockNumber,
		"settlementType", settlementType,
		"txnHash", commitmentPostingTxn.Hash().Hex(),
		"nonce", commitmentPostingTxn.Nonce(),
		"decayPercentage", decayPercentage,
	)
	u.metrics.LastSentNonce.Set(float64(commitmentPostingTxn.Nonce()))
	return u.addSettlement(
		ctx,
		update,
		settlementType,
		decayPercentage,
		window,
		commitmentPostingTxn.Hash().Bytes(),
		commitmentPostingTxn.Nonce(),
	)
}

func (u *Updater) addSettlement(
	ctx context.Context,
	update *preconf.PreconfmanagerOpenedCommitmentStored,
	settlementType SettlementType,
	decayPercentage int64,
	window int64,
	postingTxnHash []byte,
	nonce uint64,
) error {
	u.logger.Info("adding settlement to winner register",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"settlementType", settlementType,
		"decayPercentage", decayPercentage,
		"window", window,
		"nonce", nonce)

	err := u.winnerRegister.AddSettlement(
		ctx,
		update.CommitmentIndex[:],
		update.TxnHash,
		int64(update.BlockNumber),
		update.BidAmt,
		update.Committer.Bytes(),
		update.CommitmentDigest[:],
		settlementType,
		decayPercentage,
		window,
		postingTxnHash,
		nonce,
	)
	if err != nil {
		u.logger.Error(
			"failed to add settlement to winner register",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}

	u.metrics.CommitmentsProcessedCount.Inc()
	switch settlementType {
	case SettlementTypeReward:
		u.metrics.RewardsCount.Inc()
		u.logger.Info("reward settlement processed successfully")
	case SettlementTypeSlash:
		u.metrics.SlashesCount.Inc()
		u.logger.Info("slash settlement processed successfully")
	}
	u.logger.Info(
		"settlement added successfully",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"type", settlementType,
		"decayPercentage", decayPercentage,
		"bid", update.BidAmt.String(),
		"blockNumber", update.BlockNumber,
	)

	return nil
}

func (u *Updater) getL1Txns(ctx context.Context, blockNum uint64) (map[string]TxMetadata, error) {
	u.logger.Debug("retrieving L1 transactions", "blockNum", blockNum)

	txns, ok := u.l1BlockCache.Get(blockNum)
	if ok {
		u.metrics.BlockTxnCacheHits.Inc()
		u.logger.Debug("cache hit for block transactions",
			"blockNum", blockNum,
			"txCount", len(txns))
		return txns, nil
	}

	u.metrics.BlockTxnCacheMisses.Inc()
	u.logger.Debug("cache miss for block transactions", "blockNum", blockNum)

	block, err := u.l1Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		u.logger.Error("failed to get block by number",
			"blockNum", blockNum,
			"error", err)
		return nil, fmt.Errorf("failed to get block by number: %w", err)
	}

	u.logger.Debug("retrieved block",
		"blockNum", blockNum,
		"blockHash", block.Hash().Hex(),
		"txCount", len(block.Transactions()))

	var txnReceipts sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	txnsArray := make([]common.Hash, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		txnsArray[i] = tx.Hash()
	}
	const bucketSize = 25 // Arbitrary number for bucket size

	numBuckets := (len(txnsArray) + bucketSize - 1) / bucketSize // Calculate the number of buckets needed, rounding up
	buckets := make([][]common.Hash, numBuckets)
	for i := 0; i < numBuckets; i++ {
		start := i * bucketSize
		end := start + bucketSize
		if end > len(txnsArray) {
			end = len(txnsArray)
		}
		buckets[i] = txnsArray[start:end]
	}

	u.logger.Info("processing transactions in buckets",
		"totalTxns", len(txnsArray),
		"bucketSize", bucketSize,
		"numBuckets", numBuckets)

	blockStart := time.Now()

	for _, bucket := range buckets {
		eg.Go(func() error {
			start := time.Now()
			u.logger.Debug("requesting batch receipts",
				"bucketSize", len(bucket))
			results, err := u.receiptBatcher.BatchReceipts(ctx, bucket)
			if err != nil {
				u.logger.Error("failed to get batch receipts",
					"error", err,
					"bucketSize", len(bucket))
				return fmt.Errorf("failed to get batch receipts: %w", err)
			}
			u.metrics.TxnReceiptRequestDuration.Observe(time.Since(start).Seconds())
			u.logger.Debug("received batch receipts",
				"duration", time.Since(start).Seconds(),
				"resultCount", len(results))
			for _, result := range results {
				if result.Err != nil {
					u.logger.Error("failed to get receipt for txn",
						"txnHash", result.Receipt.TxHash.Hex(),
						"error", result.Err)
					continue
				}

				txnReceipts.Store(result.Receipt.TxHash.Hex(), result.Receipt)
				u.logger.Debug("stored receipt",
					"txnHash", result.Receipt.TxHash.Hex(),
					"status", result.Receipt.Status)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		u.logger.Error("error while waiting for batch receipts",
			"error", err,
			"blockNum", blockNum)
		return nil, err
	}

	u.metrics.TxnReceiptRequestBlockDuration.Observe(time.Since(blockStart).Seconds())
	u.logger.Info("completed batch receipt requests for block",
		"blockNum", blockNum,
		"duration", time.Since(blockStart).Seconds())

	txnsMap := make(map[string]TxMetadata)
	for i, tx := range txnsArray {
		receipt, ok := txnReceipts.Load(tx.Hex())
		if !ok {
			u.logger.Error("receipt not found for txn",
				"txnHash", tx.Hex(),
				"index", i)
			continue
		}
		txnsMap[strings.TrimPrefix(tx.Hex(), "0x")] = TxMetadata{
			PosInBlock: i,
			Succeeded:  receipt.(*types.Receipt).Status == types.ReceiptStatusSuccessful,
		}
		u.logger.Debug("added txn to map",
			"txnHash", tx.Hex(),
			"posInBlock", i,
			"succeeded", receipt.(*types.Receipt).Status == types.ReceiptStatusSuccessful)
	}

	_ = u.l1BlockCache.Add(blockNum, txnsMap)
	u.logger.Debug("added block transactions to cache",
		"blockNum", blockNum,
		"txCount", len(txnsMap))

	return txnsMap, nil
}

// computeDecayPercentage takes startTimestamp, endTimestamp, commitTimestamp and computes a linear decay percentage
// The computation does not care what format the timestamps are in, as long as they are consistent
// (e.g they could be unix or unixMili timestamps)
func (u *Updater) computeDecayPercentage(startTimestamp, endTimestamp, commitTimestamp uint64) int64 {
	u.logger.Debug("computing decay percentage",
		"startTimestamp", startTimestamp,
		"endTimestamp", endTimestamp,
		"commitTimestamp", commitTimestamp)

	if startTimestamp >= endTimestamp || startTimestamp > commitTimestamp || endTimestamp <= commitTimestamp {
		u.logger.Debug("timestamp out of range - returning 0%",
			"startTimestamp", startTimestamp,
			"endTimestamp", endTimestamp,
			"commitTimestamp", commitTimestamp)
		return 0
	}

	// Calculate the total time in seconds
	totalTime := endTimestamp - startTimestamp
	// Calculate the time passed in seconds
	timePassed := commitTimestamp - startTimestamp
	// Calculate the decay percentage
	decayPercentage := float64(timePassed) / float64(totalTime)

	decayPercentageRound := int64(math.Round(decayPercentage * 100))
	if decayPercentageRound > 100 {
		decayPercentageRound = 100
	}
	u.logger.Debug("decay calculation complete",
		"startTimestamp", startTimestamp,
		"endTimestamp", endTimestamp,
		"commitTimestamp", commitTimestamp,
		"totalTime", totalTime,
		"timePassed", timePassed,
		"decayPercentage", decayPercentage,
		"roundedPercentage", decayPercentageRound,
	)
	return decayPercentageRound
}
