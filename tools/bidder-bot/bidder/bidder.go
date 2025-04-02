package bidder

import (
	"context"
	"encoding/hex"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/monitor"
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
	logger         *slog.Logger
	bidderClient   bidderapiv1.BidderClient
	topologyClient debugapiv1.DebugServiceClient
	l1Client       L1Client
	beaconClient   *beaconClient
	signer         keysigner.KeySigner
	gasTipCap      *big.Int
	gasFeeCap      *big.Int
	bidAmount      *big.Int
	proposerChan   <-chan *notifier.UpcomingProposer
	sentBidChan    chan<- *monitor.SentBid
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
	sentBidChan chan<- *monitor.SentBid,
) *Bidder {
	beaconClient := newBeaconClient(beaconRPCUrl, logger.With("component", "beacon_client"))

	return &Bidder{
		logger:         logger,
		bidderClient:   bidderClient,
		topologyClient: topologyClient,
		l1Client:       l1Client,
		beaconClient:   beaconClient,
		signer:         signer,
		gasTipCap:      gasTipCap,
		gasFeeCap:      gasFeeCap,
		bidAmount:      bidAmount,
		proposerChan:   proposerChan,
		sentBidChan:    sentBidChan,
	}
}

func (b *Bidder) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				b.logger.Info("bidder context done")
				return
			case upcomingProposer := <-b.proposerChan:
				b.logger.Debug("received upcoming proposer", "proposer", upcomingProposer)
				b.handle(ctx, upcomingProposer)
			}
		}
	}()
	return done
}

func (b *Bidder) handle(ctx context.Context, upcomingProposer *notifier.UpcomingProposer) {
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

	b.logger.Info("preparing to bid", "upcomingProposer slot", upcomingProposer.Slot, "targetBlockNumber", targetBlockNum)

	err = b.bid(bidCtx, b.bidAmount, targetBlockNum)
	if err != nil {
		b.logger.Error("bid failed", "error", err)
		return
	}
}

func (b *Bidder) bid(
	ctx context.Context,
	bidAmount *big.Int,
	targetBlockNum uint64,
) error {

	tx, err := b.selfETHTransfer()
	if err != nil {
		b.logger.Error("failed to create self ETH transfer transaction", "error", err)
		return err
	}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		b.logger.Error("failed to marshal transaction", "error", err)
		return err
	}
	txString := hex.EncodeToString(txBytes)

	bidStream, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
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
		return err
	}
	b.logger.Info("bid sent", "tx_hash", tx.Hash().Hex(), "amount", bidAmount.String(), "target_block_number", targetBlockNum)

	select {
	case b.sentBidChan <- &monitor.SentBid{
		TxHash:            tx.Hash(),
		TargetBlockNumber: targetBlockNum,
		BidStream:         bidStream,
	}:
	default:
		b.logger.Warn("failed to signal bid monitor: channel full")
	}
	return nil
}

func (b *Bidder) selfETHTransfer() (*types.Transaction, error) {

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

func (b *Bidder) logDebugInfo(ctx context.Context) {
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
