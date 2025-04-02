package bidder

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/notifier"
	"github.com/primev/mev-commit/x/keysigner"
)

type L1Client interface {
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

type Bidder struct {
	logger                *slog.Logger
	bidderClient          bidderapiv1.BidderClient
	topologyClient        debugapiv1.DebugServiceClient
	l1Client              L1Client
	beaconClient          *beaconClient
	signer                keysigner.KeySigner
	gasTipCap             *big.Int
	gasFeeCap             *big.Int
	bidAmount             *big.Int
	proposerChan          <-chan *notifier.UpcomingProposer
	bidAllBlocks          bool
	regularBlockBidAmount *big.Int
	optedInSlots          sync.Map // To track slots that already have opted-in proposers
}

func NewBidder(
	logger *slog.Logger,
	bidderClient bidderapiv1.BidderClient,
	topologyClient debugapiv1.DebugServiceClient,
	l1Client L1Client,
	beaconRPCUrl string,
	signer keysigner.KeySigner,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	bidAmount *big.Int,
	proposerChan <-chan *notifier.UpcomingProposer,
	bidAllBlocks bool,
	regularBlockBidAmount *big.Int,
) *Bidder {
	beaconClient := newBeaconClient(beaconRPCUrl, logger.With("component", "beacon_client"))

	return &Bidder{
		logger:                logger,
		bidderClient:          bidderClient,
		topologyClient:        topologyClient,
		l1Client:              l1Client,
		beaconClient:          beaconClient,
		signer:                signer,
		gasTipCap:             gasTipCap,
		gasFeeCap:             gasFeeCap,
		bidAmount:             bidAmount,
		proposerChan:          proposerChan,
		bidAllBlocks:          bidAllBlocks,
		regularBlockBidAmount: regularBlockBidAmount,
		optedInSlots:          sync.Map{},
	}
}

func (b *Bidder) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	// Start the main routine for handling opted-in proposers
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("bidder context done")
				return
			case upcomingProposer := <-b.proposerChan:
				b.logger.Debug("received upcoming proposer", "proposer", upcomingProposer)

				// Track this slot as having an opted-in proposer
				b.optedInSlots.Store(upcomingProposer.Slot, true)

				// Handle with the higher bid amount for opted-in proposers
				b.handleOptedInProposer(ctx, upcomingProposer)
			}
		}
	}()

	// Start bidding for regular blocks if enabled
	if b.bidAllBlocks {
		go b.startRegularBlockBidding(ctx)
		go b.startSlotCleanup(ctx)
	}

	return done
}

// Renamed from handle to handleOptedInProposer for clarity
func (b *Bidder) handleOptedInProposer(ctx context.Context, upcomingProposer *notifier.UpcomingProposer) {
	bidCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	// Upcoming proposer slot hasn't started yet, so query block number for upcoming proposer slot - 2
	upcomingSlotMinusTwo := upcomingProposer.Slot - 2
	upcomingSlotMinusTwoBlockNum, err := b.beaconClient.getBlockNumForSlot(bidCtx, upcomingSlotMinusTwo)
	if err != nil {
		b.logger.Error("failed to get block number for upcoming proposer slot - 2", "error", err)
		return
	}

	// Assume the two slots before upcoming proposer slot are NOT missed
	targetBlockNum := upcomingSlotMinusTwoBlockNum + 2

	if b.logger.Enabled(bidCtx, slog.LevelDebug) {
		b.logDebugInfo(bidCtx)
	}

	b.logger.Info("preparing to bid for opted-in proposer", "upcomingProposer slot", upcomingProposer.Slot, "targetBlockNumber", targetBlockNum)

	pc, err := b.bid(bidCtx, b.bidAmount, targetBlockNum)
	if err != nil {
		b.logger.Error("bid failed for opted-in proposer", "error", err)
		return
	}

	err = b.watchPendingBid(bidCtx, pc)
	if err != nil {
		b.logger.Error("bid failed for opted-in proposer", "error", err)
	}
}

func (b *Bidder) startRegularBlockBidding(ctx context.Context) {
	// Running at a bit less than slot frequency to ensure we don't miss any
	ticker := time.NewTicker(10 * time.Second) // Slightly under the 12s slot time
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.bidForNextRegularBlock(ctx)
		}
	}
}

func (b *Bidder) bidForNextRegularBlock(ctx context.Context) {
	bidCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Get current slot from beacon chain
	currentSlot, err := b.beaconClient.getLatestSlot(bidCtx)
	if err != nil {
		b.logger.Error("failed to get latest slot", "error", err)
		return
	}

	// Target the next slot
	targetSlot := currentSlot + 1

	// Skip if this slot already has an opted-in proposer to avoid double bidding
	if b.isOptedInSlot(targetSlot) {
		b.logger.Debug("skipping regular bid for slot with opted-in proposer", "slot", targetSlot)
		return
	}

	// Get estimated block number for the next slot
	currentBlockNum, err := b.l1Client.BlockNumber(bidCtx)
	if err != nil {
		b.logger.Error("failed to get current block number", "error", err)
		return
	}

	// Simple prediction - assuming 1:1 mapping for next block
	targetBlockNum := currentBlockNum + 1

	b.logger.Info("preparing to bid for regular block", "targetSlot", targetSlot, "targetBlockNum", targetBlockNum)

	// Create a bid with the regular block amount
	pc, err := b.bid(bidCtx, b.regularBlockBidAmount, targetBlockNum)
	if err != nil {
		b.logger.Error("regular block bid failed", "error", err)
		return
	}

	err = b.watchPendingBid(bidCtx, pc)
	if err != nil {
		b.logger.Error("regular block bid watch failed", "error", err)
	}
}

// Helper to check if a slot already has an opted-in proposer
func (b *Bidder) isOptedInSlot(slot uint64) bool {
	_, exists := b.optedInSlots.Load(slot)
	return exists
}

func (b *Bidder) startSlotCleanup(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentSlot, err := b.beaconClient.getLatestSlot(ctx)
			if err != nil {
				continue
			}

			// Clean up slots that are in the past
			b.optedInSlots.Range(func(key, value interface{}) bool {
				slot := key.(uint64)
				if slot < currentSlot {
					b.optedInSlots.Delete(slot)
				}
				return true
			})
		}
	}
}

// The rest of the code remains unchanged
func (b *Bidder) bid(
	ctx context.Context,
	bidAmount *big.Int,
	targetBlockNum uint64,
) (bidderapiv1.Bidder_SendBidClient, error) {

	tx, err := b.selfETHTransfer()
	if err != nil {
		b.logger.Error("failed to create self ETH transfer transaction", "error", err)
		return nil, err
	}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		b.logger.Error("failed to marshal transaction", "error", err)
		return nil, err
	}
	txString := hex.EncodeToString(txBytes)

	pc, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
		TxHashes:            []string{},
		Amount:              bidAmount.String(),
		BlockNumber:         int64(targetBlockNum),
		DecayStartTimestamp: time.Now().UnixMilli(),
		DecayEndTimestamp:   time.Now().Add(12 * time.Second).UnixMilli(),
		RevertingTxHashes:   []string{},
		RawTransactions:     []string{txString},
		// Do not specify slash amount
	})
	if err != nil {
		b.logger.Error("failed to send bid", "error", err)
		return nil, err
	}
	b.logger.Info("bid sent", "tx_hash", tx.Hash().Hex(), "amount", bidAmount.String(), "target_block_number", targetBlockNum)

	return pc, nil
}

func (b *Bidder) selfETHTransfer() (*types.Transaction, error) {
	// Implementation unchanged
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	address := b.signer.GetAddress()

	// Intentionally don't use pending nonce to avoid accumulating pending l1 txs
	nonce, err := b.l1Client.NonceAt(ctx, address, nil)
	if err != nil {
		b.logger.Error("Failed to get nonce", "error", err)
		return nil, err
	}

	chainID, err := b.l1Client.ChainID(ctx)
	if err != nil {
		b.logger.Error("Failed to get network ID", "error", err)
		return nil, err
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &address,
		Value:     big.NewInt(7),
		Gas:       1_000_000,
		GasFeeCap: b.gasFeeCap,
		GasTipCap: b.gasTipCap,
	})

	signedTx, err := b.signer.SignTx(tx, chainID)
	if err != nil {
		b.logger.Error("Failed to sign transaction", "error", err)
		return nil, err
	}

	b.logger.Info("Self ETH transfer transaction created and signed", "tx_hash", signedTx.Hash().Hex())

	return signedTx, nil
}

func (b *Bidder) watchPendingBid(ctx context.Context, pc bidderapiv1.Bidder_SendBidClient) error {
	// Implementation unchanged
	topo, err := b.topologyClient.GetTopology(ctx, &debugapiv1.EmptyMessage{})
	if err != nil {
		b.logger.Error("failed to get topology", "error", err)
		return err
	}

	providers := topo.Topology.Fields["connected_providers"].GetListValue()
	if providers == nil || len(providers.Values) == 0 {
		return errors.New("no connected providers")
	}

	commitments := make([]*bidderapiv1.Commitment, 0)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := pc.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			b.logger.Error("failed to receive commitment", "error", err)
			return err
		}

		commitments = append(commitments, msg)
		b.logger.Debug("received commitment", "commitment", msg)

		if len(commitments) == len(providers.Values) {
			b.logger.Info("all commitments received")
			return nil
		} else {
			b.logger.Warn(
				"not all commitments received",
				"received", len(commitments),
				"expected", len(providers.Values),
			)
		}
	}
	return errors.New("bid timeout, not all commitments received")
}

func (b *Bidder) logDebugInfo(ctx context.Context) {
	// Implementation unchanged
	go func() {
		latestSlot, err := b.beaconClient.getLatestSlot(ctx)
		if err != nil {
			b.logger.Error("failed to get current beacon slot", "error", err)
		} else {
			b.logger.Debug("current beacon slot", "slot", latestSlot)
		}
	}()

	go func() {
		elLatestBlockNum, err := b.l1Client.BlockNumber(ctx)
		if err != nil {
			b.logger.Error("failed to get current block number", "error", err)
		} else {
			b.logger.Debug("current execution layer block number", "block_number", elLatestBlockNum)
		}
	}()
}
