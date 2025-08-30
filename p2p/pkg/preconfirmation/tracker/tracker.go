package preconftracker

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"slices"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconfcommstore "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	"github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

const (
	blockHistoryLimit            = 10000
	allowedDelayToOpenCommitment = 10
)

type Tracker struct {
	ctxChainIDData  []byte
	peerType        p2p.PeerType
	self            common.Address
	evtMgr          events.EventManager
	store           CommitmentStore
	preconfContract PreconfContract
	providerNikePK  *bn254.G1Affine
	providerNikeSK  *fr.Element
	optsGetter      OptsGetter
	watcher         Watcher
	notifier        notifications.Notifier
	newL1Blocks     chan *blocktracker.BlocktrackerNewL1Block
	unopenedCmts    chan *preconfcommstore.PreconfmanagerUnopenedCommitmentStored
	commitments     chan *preconfcommstore.PreconfmanagerOpenedCommitmentStored
	processed       chan *oracle.OracleCommitmentProcessed
	rewards         chan *bidderregistry.BidderregistryFundsRewarded
	returns         chan *bidderregistry.BidderregistryFundsUnlocked
	statusUpdate    chan statusUpdateTask
	blockOpened     chan int64
	triggerOpen     chan struct{}
	metrics         *metrics
	logger          *slog.Logger
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type CommitmentStore interface {
	GetCommitments(blockNum int64) ([]*store.Commitment, error)
	AddCommitment(commitment *store.Commitment) error
	SetCommitmentIndexByDigest(
		commitmentDigest,
		commitmentIndex [32]byte,
	) error
	SetStatus(
		blockNumber int64,
		bidAmt string,
		commitmentDigest []byte,
		status store.CommitmentStatus,
		details string,
	) error
	GetCommitmentByDigest(digest []byte) (*store.Commitment, error)
	UpdateSettlement(index []byte, isSlash bool) error
	UpdatePayment(digest []byte, payment, refund string) error
	ClearCommitmentIndexes(upto int64) error
	AddWinner(winner *store.BlockWinner) error
	BlockWinners() ([]*store.BlockWinner, error)
	ClearBlockNumber(blockNum int64) error
}

type PreconfContract interface {
	OpenCommitment(
		opts *bind.TransactOpts,
		params preconfcommstore.IPreconfManagerOpenCommitmentParams,
	) (*types.Transaction, error)
}

type Watcher interface {
	WatchTx(txnHash common.Hash, nonce uint64) <-chan txmonitor.Result
}

func NewTracker(
	chainID *big.Int,
	peerType p2p.PeerType,
	self common.Address,
	evtMgr events.EventManager,
	store CommitmentStore,
	preconfContract PreconfContract,
	watcher Watcher,
	notifier notifications.Notifier,
	providerNikePublicKey *bn254.G1Affine,
	providerNikeSecretKey *fr.Element,
	optsGetter OptsGetter,
	logger *slog.Logger,
) *Tracker {
	return &Tracker{
		ctxChainIDData:  []byte(fmt.Sprintf("mev-commit opening %s", chainID.String())),
		peerType:        peerType,
		self:            self,
		evtMgr:          evtMgr,
		store:           store,
		preconfContract: preconfContract,
		optsGetter:      optsGetter,
		watcher:         watcher,
		notifier:        notifier,
		providerNikePK:  providerNikePublicKey,
		providerNikeSK:  providerNikeSecretKey,
		newL1Blocks:     make(chan *blocktracker.BlocktrackerNewL1Block),
		unopenedCmts:    make(chan *preconfcommstore.PreconfmanagerUnopenedCommitmentStored),
		commitments:     make(chan *preconfcommstore.PreconfmanagerOpenedCommitmentStored),
		processed:       make(chan *oracle.OracleCommitmentProcessed),
		rewards:         make(chan *bidderregistry.BidderregistryFundsRewarded),
		returns:         make(chan *bidderregistry.BidderregistryFundsUnlocked),
		statusUpdate:    make(chan statusUpdateTask),
		blockOpened:     make(chan int64),
		triggerOpen:     make(chan struct{}),
		metrics:         newMetrics(),
		logger:          logger,
	}
}

func (t *Tracker) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	evts := []events.EventHandler{
		events.NewChannelEventHandler(egCtx, "NewL1Block", t.newL1Blocks),
		events.NewChannelEventHandler(egCtx, "UnopenedCommitmentStored", t.unopenedCmts),
		events.NewChannelEventHandler(egCtx, "OpenedCommitmentStored", t.commitments),
		events.NewChannelEventHandler(egCtx, "CommitmentProcessed", t.processed),
		events.NewChannelEventHandler(egCtx, "FundsRewarded", t.rewards),
	}

	if t.peerType == p2p.PeerTypeBidder {
		evts = append(
			evts,
			events.NewChannelEventHandler(egCtx, "FundsUnlocked", t.returns),
		)
	}

	sub, err := t.evtMgr.Subscribe(evts...)
	if err != nil {
		close(doneChan)
		t.logger.Error("failed to subscribe to events", "error", err)
		return doneChan
	}

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleNewL1Block context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("event subscription error: %w", err)
			case newL1Block := <-t.newL1Blocks:
				if err := t.handleNewL1Block(egCtx, newL1Block); err != nil {
					t.logger.Error("failed to handle new L1 block", "error", err)
					continue
				}
				t.triggerOpenCommitments()
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleUnopenedCommitmentStored context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("event subscription error: %w", err)
			case ec := <-t.unopenedCmts:
				if err := t.handleUnopenedCommitmentStored(egCtx, ec); err != nil {
					t.logger.Error(
						"failed to handle unopened commitment stored",
						"commitmentDigest", hex.EncodeToString(ec.CommitmentDigest[:]),
						"commitmentIndex", hex.EncodeToString(ec.CommitmentIndex[:]),
						"error", err,
					)
					continue
				}
			}
		}
	})

	eg.Go(func() error {
		tick := time.NewTicker(2 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("openCommitments context done")
				return nil
			case <-t.triggerOpen:
			case <-tick.C:
			}
			winners, err := t.store.BlockWinners()
			if err != nil {
				t.logger.Error("failed to get block winners", "error", err)
				continue
			}
			if len(winners) == 0 {
				t.logger.Debug("no winners to open commitments")
				continue
			}
			t.logger.Debug("stored block winners", "count", len(winners))
			oldBlockNos := make([]int64, 0)
			winners = slices.DeleteFunc(winners, func(item *store.BlockWinner) bool {
				// the last block is the latest, so if any of the previous blocks are
				// older than the allowed delay, we should not open the commitments
				if winners[len(winners)-1].BlockNumber-item.BlockNumber > allowedDelayToOpenCommitment {
					oldBlockNos = append(oldBlockNos, item.BlockNumber)
					return true
				}
				return false
			})
			// cleanup old state
			for _, oldBlockNo := range oldBlockNos {
				if err := t.store.ClearBlockNumber(oldBlockNo); err != nil {
					t.logger.Error("failed to delete commitments by block number", "blockNumber", oldBlockNo, "error", err)
				}
				t.logger.Info("old block commitments deleted", "blockNumber", oldBlockNo)
			}
			if t.peerType == p2p.PeerTypeBidder {
				if len(winners) > 2 {
					// Bidders should process the block 2 behind the current one. Ideally the
					// provider should open the commitment as they get the reward, so the incentive
					// for bidder to open is only in cases of slashes as he will get refund. Only one
					// of bidder or provider should open the commitment as 1 of the txns would
					// fail. This delay is to ensure this.
					t.logger.Debug("bidder detected, processing 2 blocks behind the current one")
					winners = winners[:len(winners)-2]
				} else {
					t.logger.Debug("no winners to open commitments")
					continue
				}
			}
			t.logger.Debug("opening commitments", "winners", len(winners))
			for _, winner := range winners {
				if err := t.openCommitments(egCtx, winner); err != nil {
					t.logger.Error("failed to open commitments", "error", err)
					continue
				}
			}

			t.blockOpened <- winners[len(winners)-1].BlockNumber
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleCommitmentProcessed context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("event subscription error: %w", err)
			case cp := <-t.processed:
				if err := t.store.UpdateSettlement(cp.CommitmentIndex[:], cp.IsSlash); err != nil {
					t.logger.Error("failed to update commitment index", "error", err)
					continue
				}
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleFundsRewarded context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("event subscription error: %w", err)
			case fr := <-t.rewards:
				if err := t.store.UpdatePayment(fr.CommitmentDigest[:], fr.Amount.String(), ""); err != nil {
					t.logger.Error("failed to update payment", "error", err)
					continue
				}
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleCommitmentStored context done")
				return nil
			case err := <-sub.Err():
				return fmt.Errorf("event subscription error: %w", err)
			case cs := <-t.commitments:
				if err := t.handleOpenedCommitmentStored(egCtx, cs); err != nil {
					t.logger.Error("failed to handle opened commitment stored", "error", err)
					continue
				}
			}
		}
	})

	eg.Go(func() error {
		return t.clearCommitments(egCtx, t.blockOpened)
	})

	eg.Go(func() error {
		return t.statusUpdater(egCtx, t.statusUpdate)
	})

	if t.peerType == p2p.PeerTypeBidder {
		eg.Go(func() error {
			for {
				select {
				case <-egCtx.Done():
					t.logger.Info("handleFundsUnlocked context done")
					return nil
				case err := <-sub.Err():
					return fmt.Errorf("event subscription error: %w", err)
				case fr := <-t.returns:
					if err := t.store.UpdatePayment(fr.CommitmentDigest[:], "", fr.Amount.String()); err != nil {
						t.logger.Error("failed to update payment", "error", err)
						continue
					}
				}
			}
		})
	}

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			t.logger.Error("failed to start preconfirmation", "error", err)
		}
	}()

	return doneChan
}

func (t *Tracker) TrackCommitment(
	ctx context.Context,
	commitment *store.Commitment,
	txn *types.Transaction,
) error {
	commitment.Status = store.CommitmentStatusPending

	if err := t.store.AddCommitment(commitment); err != nil {
		t.logger.Error("failed to add commitment", "error", err)
		return err
	}

	if txn != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t.statusUpdate <- statusUpdateTask{
			commitment: commitment,
			txnHash:    txn.Hash(),
			nonce:      txn.Nonce(),
			onSuccess:  store.CommitmentStatusStored,
		}:
		}
	}

	return nil
}

func (t *Tracker) Metrics() []prometheus.Collector {
	return t.metrics.Metrics()
}

func (t *Tracker) triggerOpenCommitments() {
	select {
	case t.triggerOpen <- struct{}{}:
	default:
	}
}

func (t *Tracker) handleNewL1Block(
	ctx context.Context,
	newL1Block *blocktracker.BlocktrackerNewL1Block,
) error {
	t.logger.Debug(
		"new L1 Block event received",
		"blockNumber", newL1Block.BlockNumber,
		"winner", newL1Block.Winner,
	)

	return t.store.AddWinner(&store.BlockWinner{
		BlockNumber: newL1Block.BlockNumber.Int64(),
		Winner:      newL1Block.Winner,
	})
}

type statusUpdateTask struct {
	commitment *store.Commitment
	txnHash    common.Hash
	nonce      uint64
	onSuccess  store.CommitmentStatus
}

func (t *Tracker) statusUpdater(
	ctx context.Context,
	taskCh chan statusUpdateTask,
) error {
	eg, ctx := errgroup.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			t.logger.Info("status updater context done")
			return eg.Wait()
		case task := <-taskCh:
			eg.Go(func() error {
				res := t.watcher.WatchTx(task.txnHash, task.nonce)
				select {
				case <-ctx.Done():
					t.logger.Info("watcher context done")
					return ctx.Err()
				case r := <-res:
					var (
						status  store.CommitmentStatus
						details string
					)
					switch task.onSuccess {
					case store.CommitmentStatusStored:
						if r.Err != nil {
							status = store.CommitmentStatusFailed
							details = fmt.Sprintf("failed to store commitment: %s", r.Err)
						} else {
							status = store.CommitmentStatusStored
						}
					case store.CommitmentStatusOpened:
						if r.Err != nil {
							status = store.CommitmentStatusFailed
							details = fmt.Sprintf("failed to open commitment: %s", r.Err)
						} else {
							status = store.CommitmentStatusOpened
							details = fmt.Sprintf("opened by %s", t.peerType)
						}
					}
					t.logger.Info(
						"commitment status update",
						"commitmentDigest", hex.EncodeToString(task.commitment.Commitment[:]),
						"status", status,
						"details", details,
					)
					if err := t.store.SetStatus(
						task.commitment.Bid.BlockNumber,
						task.commitment.Bid.BidAmount,
						task.commitment.Commitment,
						status,
						details,
					); err != nil {
						t.logger.Error("failed to set status", "error", err)
					}
					if status == store.CommitmentStatusFailed {
						notificationPayload := map[string]any{
							"commitmentDigest": hex.EncodeToString(task.commitment.Commitment[:]),
							"txnHash":          task.commitment.Bid.TxHash,
							"error":            r.Err.Error(),
						}
						switch task.onSuccess {
						case store.CommitmentStatusStored:
							t.notifier.Notify(
								notifications.NewNotification(
									notifications.TopicCommitmentStoreFailed,
									notificationPayload,
								),
							)
						case store.CommitmentStatusOpened:
							t.notifier.Notify(
								notifications.NewNotification(
									notifications.TopicCommitmentOpenFailed,
									notificationPayload,
								),
							)
						}
					}
				}
				return nil
			})
		}
	}
}

func (t *Tracker) openCommitments(
	ctx context.Context,
	newL1Block *store.BlockWinner,
) error {
	openStart := time.Now()

	commitments, err := t.store.GetCommitments(newL1Block.BlockNumber)
	if err != nil {
		t.logger.Error("failed to get commitments by block number", "blockNumber", newL1Block.BlockNumber, "error", err)
		return err
	}

	var settled, failed, alreadyOpened = 0, 0, 0

	for _, commitment := range commitments {
		switch commitment.Status {
		case store.CommitmentStatusPending, store.CommitmentStatusFailed:
			t.logger.Debug("commitment cannot be opened", "commitment", commitment)
			failed++
			continue
		case store.CommitmentStatusOpened, store.CommitmentStatusSettled, store.CommitmentStatusSlashed:
			t.logger.Debug("commitment already opened", "commitment", commitment)
			alreadyOpened++
			continue
		default:
			if commitment.CommitmentIndex == nil {
				t.logger.Debug("commitment index not found", "commitment", commitment)
				failed++
				continue
			}
		}
		if common.BytesToAddress(commitment.ProviderAddress).Cmp(newL1Block.Winner) != 0 {
			t.logger.Debug(
				"provider address does not match the winner",
				"providerAddress", commitment.ProviderAddress,
				"winner", newL1Block.Winner,
			)
			continue
		}
		startTime := time.Now()

		var commitmentIdx [32]byte
		copy(commitmentIdx[:], commitment.CommitmentIndex[:])

		bidAmt, ok := new(big.Int).SetString(commitment.Bid.BidAmount, 10)
		if !ok {
			t.logger.Error("failed to parse bid amount", "bidAmount", commitment.Bid.BidAmount)
			continue
		}

		slashAmt, ok := new(big.Int).SetString(commitment.Bid.SlashAmount, 10)
		if !ok {
			t.logger.Error("failed to parse slash amount", "slashAmount", commitment.Bid.SlashAmount)
			continue
		}

		opts, err := t.optsGetter(ctx)
		if err != nil {
			t.logger.Error("failed to get transact opts", "error", err)
			continue
		}

		zkProof, err := t.generateZKProof(commitment)
		if err != nil {
			t.logger.Error("failed to generate ZK proof", "error", err)
			continue
		}

		txn, err := t.preconfContract.OpenCommitment(
			opts,
			preconfcommstore.IPreconfManagerOpenCommitmentParams{
				UnopenedCommitmentIndex: commitmentIdx,
				BidAmt:                  bidAmt,
				BlockNumber:             uint64(commitment.Bid.BlockNumber),
				TxnHash:                 commitment.Bid.TxHash,
				RevertingTxHashes:       commitment.Bid.RevertingTxHashes,
				DecayStartTimeStamp:     uint64(commitment.Bid.DecayStartTimestamp),
				DecayEndTimeStamp:       uint64(commitment.Bid.DecayEndTimestamp),
				BidSignature:            commitment.Bid.Signature,
				SlashAmt:                slashAmt,
				ZkProof:                 zkProof,
			},
		)
		if err != nil {
			t.logger.Error("failed to open commitment", "error", err)
			continue
		}
		duration := time.Since(startTime)
		t.logger.Info("opened commitment",
			"txHash", txn.Hash(), "duration", duration,
			"blockNumber", newL1Block.BlockNumber,
			"committer", common.Bytes2Hex(commitment.ProviderAddress),
		)
		settled++
		select {
		case <-ctx.Done():
			t.logger.Info("openCommitments context done")
			return ctx.Err()
		case t.statusUpdate <- statusUpdateTask{
			commitment: commitment,
			txnHash:    txn.Hash(),
			nonce:      txn.Nonce(),
			onSuccess:  store.CommitmentStatusOpened,
		}:
		}
	}

	err = t.store.ClearBlockNumber(newL1Block.BlockNumber)
	if err != nil {
		t.logger.Error("failed to delete commitments by block number", "blockNumber", newL1Block.BlockNumber, "error", err)
		return err
	}

	openDuration := time.Since(openStart)
	t.logger.Info("commitments opened",
		"blockNumber", newL1Block.BlockNumber,
		"total", len(commitments),
		"settled", settled,
		"failed", failed,
		"alreadyOpened", alreadyOpened,
		"duration", openDuration,
	)

	t.metrics.totalCommitmentsToOpen.Add(float64(len(commitments)))
	t.metrics.totalOpenedCommitments.Add(float64(settled))
	t.metrics.blockCommitmentProcessDuration.Set(float64(openDuration))

	return nil
}

func (t *Tracker) clearCommitments(ctx context.Context, blockOpened <-chan int64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case blockNumber := <-blockOpened:
			var clearBlockNumber int64
			if blockNumber > blockHistoryLimit {
				clearBlockNumber = blockNumber - blockHistoryLimit
			}

			if clearBlockNumber == 0 {
				t.logger.Debug("no block numbers to clear")
				continue
			}

			// clear commitment indexes for all the blocks before the oldest winner
			err := t.store.ClearCommitmentIndexes(clearBlockNumber)
			if err != nil {
				t.logger.Error(
					"failed to clear commitment indexes",
					"block", blockNumber,
					"error", err,
				)
				continue
			}

			t.logger.Info("commitment indexes cleared", "blockNumber", clearBlockNumber)
		}
	}
}

func (t *Tracker) handleUnopenedCommitmentStored(
	ctx context.Context,
	ec *preconfcommstore.PreconfmanagerUnopenedCommitmentStored,
) error {
	t.metrics.totalEncryptedCommitments.Inc()
	return t.store.SetCommitmentIndexByDigest(ec.CommitmentDigest, ec.CommitmentIndex)
}

func (t *Tracker) handleOpenedCommitmentStored(
	ctx context.Context,
	cs *preconfcommstore.PreconfmanagerOpenedCommitmentStored,
) error {
	cmt, err := t.store.GetCommitmentByDigest(cs.CommitmentDigest[:])
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil
		}
		return fmt.Errorf("failed to get commitment by digest: %w", err)
	}

	if cmt.Status != store.CommitmentStatusOpened {
		var details string
		switch t.peerType {
		case p2p.PeerTypeBidder:
			details = fmt.Sprintf("opened by %s", p2p.PeerTypeProvider)
		case p2p.PeerTypeProvider:
			details = fmt.Sprintf("opened by %s", p2p.PeerTypeBidder)
		}
		if err := t.store.SetStatus(
			cmt.Bid.BlockNumber,
			cmt.Bid.BidAmount,
			cs.CommitmentDigest[:],
			store.CommitmentStatusOpened,
			details,
		); err != nil {
			return fmt.Errorf("failed to set status: %w", err)
		}
	}

	return nil
}

func (t *Tracker) generateZKProof(
	commitment *store.Commitment,
) ([]*big.Int, error) {
	pubB, err := crypto.BN254PublicKeyFromBytes(commitment.Bid.NikePublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bidder pubkey B: %w", err)
	}

	sharedC, err := crypto.BN254PublicKeyFromBytes(commitment.SharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to parse shared secret C: %w", err)
	}

	bidderX, bidderY := crypto.AffineToBigIntXY(pubB)
	sharedX, sharedY := crypto.AffineToBigIntXY(sharedC)

	if t.peerType == p2p.PeerTypeProvider {
		return t.generateProviderProof(pubB, sharedC, bidderX, bidderY, sharedX, sharedY)
	}

	return t.generateBidderProof(bidderX, bidderY, sharedX, sharedY), nil
}

func (t *Tracker) generateProviderProof(
	pubB *bn254.G1Affine,
	sharedC *bn254.G1Affine,
	bidderX, bidderY, sharedX, sharedY big.Int,
) ([]*big.Int, error) {
	proof, err := crypto.GenerateOptimizedProof(
		t.providerNikeSK,
		t.providerNikePK,
		pubB,
		sharedC,
		t.ctxChainIDData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate optimized proof: %w", err)
	}

	var cBig, zBig big.Int
	proof.C.BigInt(&cBig)
	proof.Z.BigInt(&zBig)

	providerX, providerY := crypto.AffineToBigIntXY(t.providerNikePK)

	return []*big.Int{
		&providerX,
		&providerY,
		&bidderX,
		&bidderY,
		&sharedX,
		&sharedY,
		&cBig,
		&zBig,
	}, nil
}

func (t *Tracker) generateBidderProof(
	bidderX, bidderY, sharedX, sharedY big.Int,
) []*big.Int {
	zeroInt := big.NewInt(0)
	return []*big.Int{
		zeroInt,
		zeroInt,
		&bidderX,
		&bidderY,
		&sharedX,
		&sharedY,
		zeroInt,
		zeroInt,
	}
}
