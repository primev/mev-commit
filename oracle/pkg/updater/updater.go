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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	lru "github.com/hashicorp/golang-lru/v2"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	preconf "github.com/primevprotocol/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-commit/x/contracts/events"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type SettlementType string

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
	Amount          uint64
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
		blockNum int64,
	) error
	IsSettled(ctx context.Context, commitmentIdx []byte) (bool, error)
	GetWinner(ctx context.Context, blockNum int64) (Winner, error)
	AddSettlement(
		ctx context.Context,
		commitmentIdx []byte,
		txHash string,
		blockNum int64,
		amount uint64,
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
	logger           *slog.Logger
	l1Client         EVMClient
	l2Client         EVMClient
	winnerRegister   WinnerRegister
	oracle           Oracle
	evtMgr           events.EventManager
	l1BlockCache     *lru.Cache[uint64, map[string]int]
	l2BlockTimeCache *lru.Cache[uint64, uint64]
	encryptedCmts    chan *preconf.PreconfcommitmentstoreEncryptedCommitmentStored
	openedCmts       chan *preconf.PreconfcommitmentstoreCommitmentStored
	currentWindow    atomic.Int64
	metrics          *metrics
}

func NewUpdater(
	logger *slog.Logger,
	l1Client EVMClient,
	winnerRegister WinnerRegister,
	evtMgr events.EventManager,
	oracle Oracle,
) (*Updater, error) {
	l1BlockCache, err := lru.New[uint64, map[string]int](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 block cache: %w", err)
	}
	l2BlockTimeCache, err := lru.New[uint64, uint64](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create L2 block time cache: %w", err)
	}
	return &Updater{
		logger:           logger,
		l1Client:         l1Client,
		l2Client:         l2Client,
		l1BlockCache:     l1BlockCache,
		l2BlockTimeCache: l2BlockTimeCache,
		winnerRegister:   winnerRegister,
		evtMgr:           evtMgr,
		oracle:           oracle,
		metrics:          newMetrics(),
		openedCmts:       make(chan *preconf.PreconfcommitmentstoreCommitmentStored),
		encryptedCmts:    make(chan *preconf.PreconfcommitmentstoreEncryptedCommitmentStored),
	}, nil
}

func (u *Updater) Metrics() []prometheus.Collector {
	return u.metrics.Collectors()
}

func (u *Updater) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	// Waits for events to be subscribed before returning
	startWg := sync.WaitGroup{}
	startWg.Add(2)

	eg.Go(func() error {
		ev1 := events.NewEventHandler(
			"EncryptedCommitmentStored",
			func(update *preconf.PreconfcommitmentstoreEncryptedCommitmentStored) {
				select {
				case <-egCtx.Done():
				case u.encryptedCmts <- update:
				}
			},
		)
		sub1, err := u.evtMgr.Subscribe(ev1)
		if err != nil {
			startWg.Done()
			return fmt.Errorf("failed to subscribe to encrypted commitments: %w", err)
		}
		defer sub1.Unsubscribe()

		ev2 := events.NewEventHandler(
			"CommitmentStored",
			func(update *preconf.PreconfcommitmentstoreCommitmentStored) {
				select {
				case <-egCtx.Done():
				case u.openedCmts <- update:
				}
			},
		)
		sub2, err := u.evtMgr.Subscribe(ev2)
		if err != nil {
			startWg.Done()
			return fmt.Errorf("failed to subscribe to opened commitments: %w", err)
		}
		defer sub2.Unsubscribe()

		ev3 := events.NewEventHandler(
			"NewWindow",
			func(update *blocktracker.BlocktrackerNewWindow) {
				u.currentWindow.Store(update.Window.Int64())
			},
		)

		sub3, err := u.evtMgr.Subscribe(ev3)
		if err != nil {
			startWg.Done()
			return fmt.Errorf("failed to subscribe to new window: %w", err)
		}
		defer sub3.Unsubscribe()

		startWg.Done()

		select {
		case <-egCtx.Done():
			return nil
		case err := <-sub1.Err():
			return err
		case err := <-sub2.Err():
			return err
		case err := <-sub3.Err():
			return err
		}
	})

	eg.Go(func() error {
		startWg.Done()

		for {
			select {
			case <-egCtx.Done():
				return nil
			case ec := <-u.encryptedCmts:
				if err := u.handleEncryptedCommitment(egCtx, ec); err != nil {
					return err
				}
			case oc := <-u.openedCmts:
				if err := u.handleOpenedCommitment(egCtx, oc); err != nil {
					return err
				}
			}
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			u.logger.Error("failed to start updater", "error", err)
		}
	}()

	startWg.Wait()

	return doneChan
}

func (u *Updater) handleEncryptedCommitment(
	ctx context.Context,
	update *preconf.PreconfcommitmentstoreEncryptedCommitmentStored,
) error {
	err := u.winnerRegister.AddEncryptedCommitment(
		ctx,
		update.CommitmentIndex[:],
		update.Commiter.Bytes(),
		update.CommitmentDigest[:],
		update.CommitmentSignature,
		update.BlockCommitedAt.Int64(),
	)
	if err != nil {
		u.logger.Error(
			"failed to add encrypted commitment",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}
	u.metrics.EncryptedCommitmentsCount.Inc()
	u.logger.Debug(
		"added encrypted commitment",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"blockNumber", update.BlockCommitedAt.Int64(),
	)
	return nil
}

func (u *Updater) handleOpenedCommitment(
	ctx context.Context,
	update *preconf.PreconfcommitmentstoreCommitmentStored,
) error {
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
			"duplicate open commitment event",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		)
		u.metrics.DuplicateCommitmentsCount.Inc()
		return nil
	}

	winner, err := u.winnerRegister.GetWinner(ctx, int64(update.BlockNumber))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("winner not found", "blockNumber", update.BlockNumber)
			u.metrics.NoWinnerCount.Inc()
			return nil
		}
		u.logger.Error(
			"failed to get winner",
			"blockNumber", update.BlockCommitedAt.Int64(),
			"error", err,
		)
		return err
	}

	if u.currentWindow.Load() > 2 && winner.Window < u.currentWindow.Load()-2 {
		u.logger.Info(
			"commitment is too old",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"winner", winner.Winner,
			"window", winner.Window,
			"currentWindow", u.currentWindow.Load(),
		)
		u.metrics.CommitmentsTooOldCount.Inc()
		return nil
	}

	if common.BytesToAddress(winner.Winner).Cmp(update.Commiter) != 0 {
		// The winner is not the committer of the commitment
		u.logger.Info(
			"winner is not the committer",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"winner", winner.Winner,
			"committer", update.Commiter.Hex(),
		)
		return nil
	}

	txns, err := u.getL1Txns(ctx, update.BlockNumber)
	if err != nil {
		u.logger.Error(
			"failed to get L1 txns",
			"blockNumber", update.BlockNumber,
			"error", err,
		)
		return err
	}

	// Compute the decay percentage
	l2BlockTime, err := u.getL2BlockTime(ctx, update.BlockCommitedAt.Uint64())
	if err != nil {
		u.logger.Error(
			"failed to get L2 block time",
			"blockNumber", update.BlockCommitedAt.Int64(),
			"error", err,
		)
		return err
	}

	decayPercentage := u.computeDecayPercentage(
		update.DecayStartTimeStamp,
		update.DecayEndTimeStamp,
		l2BlockTime,
	)

	commitmentTxnHashes := strings.Split(update.TxnHash, ",")
	// Ensure Bundle is atomic and present in the block
	for i := 0; i < len(commitmentTxnHashes); i++ {
		posInBlock, found := txns[commitmentTxnHashes[i]]
		if !found || posInBlock != txns[commitmentTxnHashes[0]]+i {
			u.logger.Info(
				"bundle is not atomic",
				"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
				"txnHash", update.TxnHash,
				"blockNumber", update.BlockNumber,
				"found", found,
				"posInBlock", posInBlock,
				"expectedPosInBlock", txns[commitmentTxnHashes[0]]+i,
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
	}

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
	update *preconf.PreconfcommitmentstoreCommitmentStored,
	settlementType SettlementType,
	decayPercentage int64,
	window int64,
) error {
	commitmentPostingTxn, err := u.oracle.ProcessBuilderCommitmentForBlockNumber(
		update.CommitmentIndex,
		big.NewInt(0).SetUint64(update.BlockNumber),
		update.Commiter,
		settlementType == SettlementTypeSlash,
		big.NewInt(decayPercentage),
	)
	if err != nil {
		u.logger.Error(
			"failed to process commitment",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}
	u.logger.Info(
		"settled commitment",
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
	update *preconf.PreconfcommitmentstoreCommitmentStored,
	settlementType SettlementType,
	decayPercentage int64,
	window int64,
	postingTxnHash []byte,
	nonce uint64,
) error {
	err := u.winnerRegister.AddSettlement(
		ctx,
		update.CommitmentIndex[:],
		update.TxnHash,
		int64(update.BlockNumber),
		update.Bid,
		update.Commiter.Bytes(),
		update.CommitmentHash[:],
		settlementType,
		decayPercentage,
		window,
		postingTxnHash,
		nonce,
	)
	if err != nil {
		u.logger.Error(
			"failed to add settlement",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}

	u.metrics.CommitmentsProcessedCount.Inc()
	switch settlementType {
	case SettlementTypeReward:
		u.metrics.RewardsCount.Inc()
	case SettlementTypeSlash:
		u.metrics.SlashesCount.Inc()
	}
	u.logger.Info(
		"added settlement",
		"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
		"type", settlementType,
		"decayPercentage", decayPercentage,
	)

	return nil
}

func (u *Updater) getL1Txns(ctx context.Context, blockNum uint64) (map[string]int, error) {
	txns, ok := u.l1BlockCache.Get(blockNum)
	if ok {
		u.metrics.BlockTxnCacheHits.Inc()
		return txns, nil
	}

	u.metrics.BlockTxnCacheMisses.Inc()

	blk, err := u.l1Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %w", err)
	}

	txnsInBlock := make(map[string]int)
	for posInBlock, tx := range blk.Transactions() {
		txnsInBlock[strings.TrimPrefix(tx.Hash().Hex(), "0x")] = posInBlock
	}
	_ = u.l1BlockCache.Add(blockNum, txnsInBlock)

	return txnsInBlock, nil
}

func (u *Updater) getL2BlockTime(ctx context.Context, blockNum uint64) (uint64, error) {
	time, ok := u.l2BlockTimeCache.Get(blockNum)
	if ok {
		u.metrics.BlockTimeCacheHits.Inc()
		return time, nil
	}

	u.metrics.BlockTimeCacheMisses.Inc()

	blk, err := u.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		return 0, fmt.Errorf("failed to get block by number: %w", err)
	}

	_ = u.l2BlockTimeCache.Add(blockNum, blk.Header().Time)

	return blk.Header().Time, nil
}

// computeDecayPercentage takes startTimestamp, endTimestamp, commitTimestamp and computes a linear decay percentage
// The computation does not care what format the timestamps are in, as long as they are consistent
// (e.g they could be unix or unixMili timestamps)
func (u *Updater) computeDecayPercentage(startTimestamp, endTimestamp, commitTimestamp uint64) int64 {
	if startTimestamp >= endTimestamp || startTimestamp > commitTimestamp || endTimestamp <= commitTimestamp {
		u.logger.Info("invalid timestamps", "startTimestamp", startTimestamp, "endTimestamp", endTimestamp, "commitTimestamp", commitTimestamp)
		return 0
	}

	// Calculate the total time in seconds
	totalTime := endTimestamp - startTimestamp
	u.logger.Info("totalTime", "totalTime", totalTime)
	// Calculate the time passed in seconds
	timePassed := commitTimestamp - startTimestamp
	u.logger.Info("timePassed", "timePassed", timePassed)
	// Calculate the decay percentage
	decayPercentage := float64(timePassed) / float64(totalTime)
	u.logger.Info("decayPercentage", "decayPercentage", decayPercentage)

	decayPercentageRound := int64(math.Round(decayPercentage * 100))
	u.logger.Info("decayPercentageRound", "decayPercentageRound", decayPercentageRound)
	return decayPercentageRound
}
