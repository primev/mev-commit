package preconftracker

import (
	"context"
	"encoding/hex"
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
	preconfcommstore "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	"github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

const (
	allowedDelayToOpenCommitment = 10
)

type Tracker struct {
	peerType        p2p.PeerType
	self            common.Address
	evtMgr          events.EventManager
	store           CommitmentStore
	preconfContract PreconfContract
	providerNikePK  *bn254.G1Affine
	providerNikeSK  *fr.Element
	receiptGetter   txmonitor.BatchReceiptGetter
	optsGetter      OptsGetter
	newL1Blocks     chan *blocktracker.BlocktrackerNewL1Block
	unopenedCmts    chan *preconfcommstore.PreconfmanagerUnopenedCommitmentStored
	commitments     chan *preconfcommstore.PreconfmanagerOpenedCommitmentStored
	triggerOpen     chan struct{}
	metrics         *metrics
	logger          *slog.Logger
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type CommitmentStore interface {
	GetCommitments(blockNum int64) ([]*store.EncryptedPreConfirmationWithDecrypted, error)
	AddCommitment(commitment *store.EncryptedPreConfirmationWithDecrypted) error
	ClearBlockNumber(blockNum int64) error
	DeleteCommitmentByDigest(
		blockNum int64,
		bidAmt string,
		digest [32]byte,
	) error
	SetCommitmentIndexByDigest(
		commitmentDigest,
		commitmentIndex [32]byte,
	) error
	ClearCommitmentIndexes(upto int64) error
	AddWinner(winner *store.BlockWinner) error
	BlockWinners() ([]*store.BlockWinner, error)
}

type PreconfContract interface {
	OpenCommitment(
		opts *bind.TransactOpts,
		encryptedCommitmentIndex [32]byte,
		bidAmt *big.Int,
		blockNumber uint64,
		txnHash string,
		revertingTxHashes string,
		decayStartTimeStamp uint64,
		decayEndTimeStamp uint64,
		bidSignature []byte,
		sharedSecretKey []byte,
		zkProof []*big.Int,
	) (*types.Transaction, error)
}

func NewTracker(
	peerType p2p.PeerType,
	self common.Address,
	evtMgr events.EventManager,
	store CommitmentStore,
	preconfContract PreconfContract,
	receiptGetter txmonitor.BatchReceiptGetter,
	providerNikePublicKey *bn254.G1Affine,
	providerNikeSecretKey *fr.Element,
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
		providerNikePK:  providerNikePublicKey,
		providerNikeSK:  providerNikeSecretKey,
		newL1Blocks:     make(chan *blocktracker.BlocktrackerNewL1Block),
		unopenedCmts:    make(chan *preconfcommstore.PreconfmanagerUnopenedCommitmentStored),
		commitments:     make(chan *preconfcommstore.PreconfmanagerOpenedCommitmentStored),
		triggerOpen:     make(chan struct{}),
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
			"UnopenedCommitmentStored",
			func(ec *preconfcommstore.PreconfmanagerUnopenedCommitmentStored) {
				select {
				case <-egCtx.Done():
					t.logger.Info("UnopenedCommitmentStored context done")
				case t.unopenedCmts <- ec:
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
				"OpenedCommitmentStored",
				func(cs *preconfcommstore.PreconfmanagerOpenedCommitmentStored) {
					select {
					case <-egCtx.Done():
						t.logger.Info("OpenedCommitmentStored context done")
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
		}
	})

	eg.Go(func() error {
		return t.clearCommitments(egCtx)
	})

	if t.peerType == p2p.PeerTypeBidder {
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
	return t.store.AddCommitment(commitment)
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
		"window", newL1Block.Window,
	)

	return t.store.AddWinner(&store.BlockWinner{
		BlockNumber: newL1Block.BlockNumber.Int64(),
		Winner:      newL1Block.Winner,
	})
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

	failedCommitments := make([]common.Hash, 0)
	settled := 0
	for _, commitment := range commitments {
		if commitment.CommitmentIndex == nil {
			t.logger.Debug("commitment index not found", "commitment", commitment)
			if commitment.TxnHash != (common.Hash{}) {
				failedCommitments = append(failedCommitments, commitment.TxnHash)
			}
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

		pubB, err := crypto.BN254PublicKeyFromBytes(commitment.Bid.NikePublicKey)
		if err != nil {
			t.logger.Error("failed to parse bidder pubkey B", "error", err)
			continue
		}

		var bidderXBig, bidderYBig big.Int
		pubB.X.BigInt(&bidderXBig)
		pubB.Y.BigInt(&bidderYBig)

		var zkProof []*big.Int
		if t.peerType == p2p.PeerTypeProvider {

			sharedC, err := crypto.BN254PublicKeyFromBytes(commitment.PreConfirmation.SharedSecret)
			if err != nil {
				t.logger.Error("failed to parse shared secret C = B^a", "error", err)
				continue
			}

			contextData := []byte("mev-commit opening, mainnet, v1.0")
			proof, err := crypto.GenerateOptimizedProof(t.providerNikeSK, t.providerNikePK, pubB, sharedC, contextData)
			if err != nil {
				t.logger.Error("failed to generate ZK proof for openCommitment", "error", err)
				continue
			}

			var cBig, zBig big.Int
			proof.C.BigInt(&cBig)
			proof.Z.BigInt(&zBig)

			var providerXBig, providerYBig big.Int
			t.providerNikePK.X.BigInt(&providerXBig)
			t.providerNikePK.Y.BigInt(&providerYBig)

			zkProof = []*big.Int{&providerXBig, &providerYBig, &bidderXBig, &bidderYBig, &cBig, &zBig}
		} else {
			zeroInt := big.NewInt(0)
			zkProof = []*big.Int{zeroInt, zeroInt, &bidderXBig, &bidderYBig, zeroInt, zeroInt}
		}
		txn, err := t.preconfContract.OpenCommitment(
			opts,
			commitmentIdx,
			bidAmt,
			uint64(commitment.PreConfirmation.Bid.BlockNumber),
			commitment.PreConfirmation.Bid.TxHash,
			commitment.PreConfirmation.Bid.RevertingTxHashes,
			uint64(commitment.PreConfirmation.Bid.DecayStartTimestamp),
			uint64(commitment.PreConfirmation.Bid.DecayEndTimestamp),
			commitment.PreConfirmation.Bid.Signature,
			commitment.PreConfirmation.SharedSecret,
			zkProof,
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
		"failed", len(failedCommitments),
		"duration", openDuration,
	)

	t.metrics.totalCommitmentsToOpen.Add(float64(len(commitments)))
	t.metrics.totalOpenedCommitments.Add(float64(settled))
	t.metrics.blockCommitmentProcessDuration.Set(float64(openDuration))

	if len(failedCommitments) > 0 {
		t.logger.Info("processing failed commitments", "count", len(failedCommitments))
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

func (t *Tracker) clearCommitments(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
		winners, err := t.store.BlockWinners()
		if err != nil {
			return err
		}

		if len(winners) == 0 {
			continue
		}

		// clear commitment indexes for all the blocks before the oldest winner
		err = t.store.ClearCommitmentIndexes(winners[0].BlockNumber)
		if err != nil {
			t.logger.Error(
				"failed to clear commitment indexes",
				"block", winners[0].BlockNumber,
				"error", err,
			)
			continue
		}

		t.logger.Info("commitment indexes cleared", "blockNumber", winners[0].BlockNumber)
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
	// In case of bidders this event keeps track of the commitments already opened
	// by the provider.
	return t.store.DeleteCommitmentByDigest(int64(cs.BlockNumber), cs.BidAmt.String(), cs.CommitmentDigest)
}
