package preconftracker

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	preconfcommstore "github.com/primev/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/store"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type Tracker struct {
	peerType        p2p.PeerType
	self            common.Address
	evtMgr          events.EventManager
	store           CommitmentStore
	preconfContract PreconfContract
	receiptGetter   txmonitor.BatchReceiptGetter
	optsGetter      OptsGetter
	newL1Blocks     chan *blocktracker.BlocktrackerNewL1Block
	enryptedCmts    chan *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored
	commitments     chan *preconfcommstore.PreconfcommitmentstoreCommitmentStored
	winners         map[int64]*blocktracker.BlocktrackerNewL1Block
	metrics         *metrics
	logger          *slog.Logger
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type CommitmentStore interface {
	GetCommitmentsByBlockNumber(blockNum int64) ([]*store.EncryptedPreConfirmationWithDecrypted, error)
	AddCommitment(commitment *store.EncryptedPreConfirmationWithDecrypted)
	DeleteCommitmentByBlockNumber(blockNum int64) error
	DeleteCommitmentByDigest(
		blockNum int64,
		digest [32]byte,
	) error
	SetCommitmentIndexByCommitmentDigest(
		commitmentDigest,
		commitmentIndex [32]byte,
	) error
}

type PreconfContract interface {
	OpenCommitment(
		opts *bind.TransactOpts,
		encryptedCommitmentIndex [32]byte,
		bid *big.Int,
		blockNumber uint64,
		txnHash string,
		decayStartTimeStamp uint64,
		decayEndTimeStamp uint64,
		bidSignature []byte,
		commitmentSignature []byte,
		sharedSecretKey []byte,
	) (*types.Transaction, error)
}

func NewTracker(
	peerType p2p.PeerType,
	self common.Address,
	evtMgr events.EventManager,
	store CommitmentStore,
	preconfContract PreconfContract,
	receiptGetter txmonitor.BatchReceiptGetter,
	optsGetter OptsGetter,
	logger *slog.Logger,
) *Tracker {
	return &Tracker{
		peerType:        peerType,
		self:            self,
		evtMgr:          evtMgr,
		store:           store,
		preconfContract: preconfContract,
		receiptGetter:   receiptGetter,
		optsGetter:      optsGetter,
		newL1Blocks:     make(chan *blocktracker.BlocktrackerNewL1Block),
		enryptedCmts:    make(chan *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored),
		commitments:     make(chan *preconfcommstore.PreconfcommitmentstoreCommitmentStored),
		winners:         make(map[int64]*blocktracker.BlocktrackerNewL1Block),
		metrics:         newMetrics(),
		logger:          logger,
	}
}

func (t *Tracker) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	evts := []events.EventHandler{
		events.NewEventHandler(
			"NewL1Block",
			func(newL1Block *blocktracker.BlocktrackerNewL1Block) {
				select {
				case <-egCtx.Done():
					t.logger.Info("NewL1Block context done")
				case t.newL1Blocks <- newL1Block:
				}
			},
		),
		events.NewEventHandler(
			"EncryptedCommitmentStored",
			func(ec *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored) {
				select {
				case <-egCtx.Done():
					t.logger.Info("EncryptedCommitmentStored context done")
				case t.enryptedCmts <- ec:
				}
			},
		),
		events.NewEventHandler(
			"FundsRewarded",
			func(fr *bidderregistry.BidderregistryFundsRewarded) {
				if fr.Bidder.Cmp(t.self) == 0 || fr.Provider.Cmp(t.self) == 0 {
					t.logger.Info("funds settled for bid",
						"commitmentDigest", common.BytesToHash(fr.CommitmentDigest[:]),
						"window", fr.Window,
						"amount", fr.Amount,
						"bidder", fr.Bidder,
						"provider", fr.Provider,
					)
				}
			},
		),
	}

	if t.peerType == p2p.PeerTypeBidder {
		evts = append(
			evts,
			events.NewEventHandler(
				"CommitmentStored",
				func(cs *preconfcommstore.PreconfcommitmentstoreCommitmentStored) {
					select {
					case <-egCtx.Done():
						t.logger.Info("CommitmentStored context done")
					case t.commitments <- cs:
					}
				},
			),
			events.NewEventHandler(
				"FundsRetrieved",
				func(fr *bidderregistry.BidderregistryFundsRetrieved) {
					if fr.Bidder.Cmp(t.self) == 0 {
						t.logger.Info("funds returned for bid",
							"commitmentDigest", common.BytesToHash(fr.CommitmentDigest[:]),
							"amount", fr.Amount,
							"window", fr.Window,
							"bidder", fr.Bidder,
						)
					}
				},
			),
		)
	}

	sub, err := t.evtMgr.Subscribe(evts...)
	if err != nil {
		close(doneChan)
		t.logger.Error("failed to subscribe to events", "error", err)
		return doneChan
	}

	eg.Go(func() error {
		select {
		case <-egCtx.Done():
			t.logger.Info("err listener context done")
			return nil
		case err := <-sub.Err():
			return fmt.Errorf("event subscription error: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleNewL1Block context done")
				return nil
			case newL1Block := <-t.newL1Blocks:
				if err := t.handleNewL1Block(egCtx, newL1Block); err != nil {
					return err
				}
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				t.logger.Info("handleEncryptedCommitmentStored context done")
				return nil
			case ec := <-t.enryptedCmts:
				if err := t.handleEncryptedCommitmentStored(egCtx, ec); err != nil {
					return err
				}
			}
		}
	})

	if t.peerType == p2p.PeerTypeBidder {
		eg.Go(func() error {
			for {
				select {
				case <-egCtx.Done():
					t.logger.Info("handleCommitmentStored context done")
					return nil
				case cs := <-t.commitments:
					if err := t.handleCommitmentStored(egCtx, cs); err != nil {
						return err
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
	commitment *store.EncryptedPreConfirmationWithDecrypted,
) error {
	t.store.AddCommitment(commitment)
	return nil
}

func (t *Tracker) Metrics() []prometheus.Collector {
	return t.metrics.Metrics()
}

func (t *Tracker) handleNewL1Block(
	ctx context.Context,
	newL1Block *blocktracker.BlocktrackerNewL1Block,
) error {
	t.logger.Info(
		"new L1 Block event received",
		"blockNumber", newL1Block.BlockNumber,
		"winner", newL1Block.Winner,
		"window", newL1Block.Window,
	)

	openStart := time.Now()

	if t.peerType == p2p.PeerTypeBidder {
		// Bidders should process the block 1 behind the current one. Ideally the
		// provider should open the commitment as they get the reward, so the incentive
		// for bidder to open is only in cases of slashes as he will get refund. Only one
		// of bidder or provider should open the commitment as 1 of the txns would
		// fail. This delay is to ensure this.
		t.winners[newL1Block.BlockNumber.Int64()] = newL1Block
		pastBlock, ok := t.winners[newL1Block.BlockNumber.Int64()-2]
		if !ok {
			return nil
		}
		newL1Block = pastBlock
		for k := range t.winners {
			if k < pastBlock.BlockNumber.Int64() {
				delete(t.winners, k)
			}
		}
	}

	commitments, err := t.store.GetCommitmentsByBlockNumber(newL1Block.BlockNumber.Int64())
	if err != nil {
		return err
	}

	failedCommitments := make([]common.Hash, 0)
	settled := 0
	for _, commitment := range commitments {
		if commitment.CommitmentIndex == nil {
			failedCommitments = append(failedCommitments, commitment.TxnHash)
			continue
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

		opts, err := t.optsGetter(ctx)
		if err != nil {
			t.logger.Error("failed to get transact opts", "error", err)
			continue
		}

		txHash, err := t.preconfContract.OpenCommitment(
			opts,
			commitmentIdx,
			bidAmt,
			uint64(commitment.PreConfirmation.Bid.BlockNumber),
			commitment.PreConfirmation.Bid.TxHash,
			uint64(commitment.PreConfirmation.Bid.DecayStartTimestamp),
			uint64(commitment.PreConfirmation.Bid.DecayEndTimestamp),
			commitment.PreConfirmation.Bid.Signature,
			commitment.PreConfirmation.Signature,
			commitment.PreConfirmation.SharedSecret,
		)
		if err != nil {
			t.logger.Error("failed to open commitment", "error", err)
			continue
		}
		duration := time.Since(startTime)
		t.logger.Info("opened commitment",
			"txHash", txHash, "duration", duration,
			"blockNumber", newL1Block.BlockNumber,
			"commiter", common.Bytes2Hex(commitment.ProviderAddress),
		)
		settled++
	}

	err = t.store.DeleteCommitmentByBlockNumber(newL1Block.BlockNumber.Int64())
	if err != nil {
		t.logger.Error("failed to delete commitments by block number", "error", err)
		return err
	}

	openDuration := time.Since(openStart)
	t.metrics.totalCommitmentsToOpen.Add(float64(len(commitments)))
	t.metrics.totalOpenedCommitments.Add(float64(settled))
	t.metrics.blockCommitmentProcessDuration.Set(float64(openDuration))

	t.logger.Info("commitments opened",
		"blockNumber", newL1Block.BlockNumber,
		"total", len(commitments),
		"settled", settled,
		"failed", len(failedCommitments),
		"duration", openDuration,
	)

	if len(failedCommitments) > 0 {
		receipts, err := t.receiptGetter.BatchReceipts(ctx, failedCommitments)
		if err != nil {
			t.logger.Warn("failed to get receipts for failed commitments", "error", err)
			return nil
		}
		for i, receipt := range receipts {
			t.logger.Debug("receipt for failed commitment",
				"txHash", failedCommitments[i],
				"error", receipt.Err,
			)
		}
	}

	return nil
}

func (t *Tracker) handleEncryptedCommitmentStored(
	ctx context.Context,
	ec *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored,
) error {
	t.metrics.totalEncryptedCommitments.Inc()
	return t.store.SetCommitmentIndexByCommitmentDigest(ec.CommitmentDigest, ec.CommitmentIndex)
}

func (t *Tracker) handleCommitmentStored(
	ctx context.Context,
	cs *preconfcommstore.PreconfcommitmentstoreCommitmentStored,
) error {
	// In case of bidders this event keeps track of the commitments already opened
	// by the provider.
	return t.store.DeleteCommitmentByDigest(int64(cs.BlockNumber), cs.CommitmentHash)
}
