package bidder

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/monitor"
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
	logger          *slog.Logger
	bidderClient    bidderapiv1.BidderClient
	topologyClient  debugapiv1.DebugServiceClient
	l1Client        L1Client
	signer          keysigner.KeySigner
	gasTipCap       *big.Int
	gasFeeCap       *big.Int
	bidAmount       *big.Int
	targetBlockChan <-chan TargetBlock
	acceptedBidChan chan<- *monitor.AcceptedBid
}

type TargetBlock struct {
	Num  uint64
	Time time.Time
}

func NewBidder(
	logger *slog.Logger,
	bidderClient bidderapiv1.BidderClient,
	topologyClient debugapiv1.DebugServiceClient,
	l1Client L1Client,
	signer keysigner.KeySigner,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	bidAmount *big.Int,
	targetBlockChan <-chan TargetBlock,
	acceptedBidChan chan<- *monitor.AcceptedBid,
) *Bidder {
	return &Bidder{
		logger:          logger,
		bidderClient:    bidderClient,
		topologyClient:  topologyClient,
		l1Client:        l1Client,
		signer:          signer,
		gasTipCap:       gasTipCap,
		gasFeeCap:       gasFeeCap,
		bidAmount:       bidAmount,
		targetBlockChan: targetBlockChan,
		acceptedBidChan: acceptedBidChan,
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
			case targetBlock := <-b.targetBlockChan:
				b.logger.Debug("received target block", "target_block_number", targetBlock.Num, "target_block_time", targetBlock.Time)
				b.handle(ctx, targetBlock)
			}
		}
	}()
	return done
}

func (b *Bidder) handle(ctx context.Context, targetBlock TargetBlock) {
	delay := time.Until(targetBlock.Time)
	bidCtx, cancel := context.WithTimeout(ctx, delay)
	defer cancel()

	b.logger.Info("preparing to bid",
		"targetBlockNumber", targetBlock.Num,
		"targetBlockTime", targetBlock.Time,
	)

	bidStream, tx, err := b.bid(bidCtx, b.bidAmount, targetBlock)
	if err != nil {
		b.logger.Error("bid failed", "error", err)
		return
	}

	err = b.watchPendingBid(bidCtx, bidStream)
	if err != nil {
		b.logger.Error("failed to watch pending bid", "error", err)
		return
	}

	select {
	case b.acceptedBidChan <- &monitor.AcceptedBid{
		TxHash:            tx.Hash(),
		TargetBlockNumber: targetBlock.Num,
	}:
	default:
		b.logger.Warn("failed to signal bid monitor: channel full")
	}
}

func (b *Bidder) bid(
	ctx context.Context,
	bidAmount *big.Int,
	targetBlock TargetBlock,
) (bidStream bidderapiv1.Bidder_SendBidClient, tx *types.Transaction, err error) {

	tx, err = b.selfETHTransfer()
	if err != nil {
		b.logger.Error("failed to create self ETH transfer transaction", "error", err)
		return nil, nil, err
	}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		b.logger.Error("failed to marshal transaction", "error", err)
		return nil, nil, err
	}
	txString := hex.EncodeToString(txBytes)

	bidStream, err = b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
		TxHashes:            []string{},
		Amount:              bidAmount.String(),
		BlockNumber:         int64(targetBlock.Num),
		DecayStartTimestamp: time.Now().UnixMilli(),
		DecayEndTimestamp:   time.Now().Add(12 * time.Second).UnixMilli(),
		RevertingTxHashes:   []string{},
		RawTransactions:     []string{txString},
		// Do not specify slash amount
	})
	if err != nil {
		b.logger.Error("failed to send bid", "error", err)
		return nil, nil, err
	}
	b.logger.Info("bid sent",
		"tx_hash", tx.Hash().Hex(),
		"amount", bidAmount.String(),
		"target_block_number", targetBlock.Num,
		"target_block_time", targetBlock.Time,
	)

	return bidStream, tx, nil
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

func (b *Bidder) watchPendingBid(ctx context.Context, pc bidderapiv1.Bidder_SendBidClient) error {
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
	if len(commitments) > 0 {
		return nil
	} else {
		return errors.New("bid timeout, no commitments received")
	}
}
