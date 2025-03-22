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
	"github.com/primev/mev-commit/tools/bidder-bot/notifier"
	"github.com/primev/mev-commit/x/keysigner"
)

type L1Client interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
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
	signer         keysigner.KeySigner
	gasTipCap      *big.Int
	gasFeeCap      *big.Int
	proposerChan   <-chan *notifier.UpcomingProposer
}

func NewBidder(
	logger *slog.Logger,
	bidderClient bidderapiv1.BidderClient,
	topologyClient debugapiv1.DebugServiceClient,
	l1Client L1Client,
	signer keysigner.KeySigner,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	proposerChan <-chan *notifier.UpcomingProposer,
) *Bidder {
	return &Bidder{
		logger:         logger,
		bidderClient:   bidderClient,
		topologyClient: topologyClient,
		l1Client:       l1Client,
		signer:         signer,
		gasTipCap:      gasTipCap,
		gasFeeCap:      gasFeeCap,
		proposerChan:   proposerChan,
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
			case proposer := <-b.proposerChan:
				b.logger.Debug("received upcoming proposer", "proposer", proposer)
				bidAmount := big.NewInt(5000000000000000) // 0.005 eth
				pc, err := b.Bid(ctx, bidAmount)
				if err != nil {
					b.logger.Error("bid failed", "error", err)
					continue
				}
				err = b.watchPendingBid(ctx, pc)
				if err != nil {
					b.logger.Error("bid failed", "error", err)
					continue
				}
			}
		}
	}()
	return done
}

func (b *Bidder) Bid(
	ctx context.Context,
	bidAmount *big.Int,
) (bidderapiv1.Bidder_SendBidClient, error) {

	tx, err := b.SelfETHTransfer()
	if err != nil {
		b.logger.Error("failed to create self ETH transfer transaction", "error", err)
		return nil, err
	}

	// TODO: sanity check tx serialization
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		b.logger.Error("failed to marshal transaction", "error", err)
		return nil, err
	}
	txString := hex.EncodeToString(txBytes)

	blkNumber, err := b.l1Client.BlockNumber(ctx)
	if err != nil {
		b.logger.Error("failed to get block number", "error", err)
		return nil, err
	}

	pc, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
		TxHashes:            []string{},
		Amount:              bidAmount.String(),
		BlockNumber:         int64(blkNumber + 1),
		DecayStartTimestamp: time.Now().UnixMilli(),
		DecayEndTimestamp:   time.Now().Add(12 * time.Second).UnixMilli(),
		RevertingTxHashes:   []string{},
		RawTransactions:     []string{txString},
		SlashAmount:         big.NewInt(0).String(), // TODO: determine slash amount
	})
	if err != nil {
		b.logger.Error("failed to send bid", "error", err)
		return nil, err
	}

	return pc, nil
}

func (b *Bidder) SelfETHTransfer() (*types.Transaction, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	address := b.signer.GetAddress()

	nonce, err := b.l1Client.PendingNonceAt(ctx, address)
	if err != nil {
		b.logger.Error("Failed to get pending nonce", "error", err)
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
		GasFeeCap: b.gasFeeCap, // TODO: sanity check fees here
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

// TODO: add logic where current bid cant send until the previous bid's block has been proposed as 'latest'.
// This'd prevent any some nonce issues, but to cover all edge cases, we'd need to wait til previous
// bid's relevant L1 block has been fully 'finalized', OR we've received commitments from all providers
// for the previous bid's block.
// Likely introduce a new go routine that handles 'pending bids' and monitors them till bid is finalized or failed.
// TODO: Also include db component here where restarted service still waits for pending bids to finalize before next bid.
// TODO: tracking / metrics on # commitments, and if tx lands on L1.
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

	// TODO: replace 12 second timeout with checking that the relevant L1 block is latest or finalized
	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)

	defer cancel()

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

		// TODO: confirm commitment + timeout waiting logic
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
