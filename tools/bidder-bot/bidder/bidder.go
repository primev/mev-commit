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
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/x/keysigner"
)

const (
	upcomingProposerTopic = "validator_opted_in"
	slotDuration          = 12 * time.Second
)

var (
	ErrUnexpectedTopic = errors.New("unexpected msg topic")
)

var nowFunc = time.Now

type L1Client interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	ChainID(ctx context.Context) (*big.Int, error)
}

// new name
type BidderClient struct {
	logger              *slog.Logger
	bidderClient        bidderapiv1.BidderClient
	topologyClient      debugapiv1.DebugServiceClient
	notificationsClient notificationsapiv1.NotificationsClient
	l1Client            L1Client
	signer              keysigner.KeySigner
	gasTipCap           *big.Int
	gasFeeCap           *big.Int
}

func NewBidderClient(
	logger *slog.Logger,
	bidderClient bidderapiv1.BidderClient,
	topologyClient debugapiv1.DebugServiceClient,
	notificationsClient notificationsapiv1.NotificationsClient,
	l1Client L1Client,
	signer keysigner.KeySigner,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
) *BidderClient {
	return &BidderClient{
		logger:              logger,
		bidderClient:        bidderClient,
		topologyClient:      topologyClient,
		notificationsClient: notificationsClient,
		l1Client:            l1Client,
		signer:              signer,
		gasTipCap:           gasTipCap,
		gasFeeCap:           gasFeeCap,
	}
}

func (b *BidderClient) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		lastMsg := nowFunc()
	RESTART:
		sub, err := b.notificationsClient.Subscribe(ctx, &notificationsapiv1.SubscribeRequest{
			Topics: []string{upcomingProposerTopic},
		})
		if err != nil {
			b.logger.Error("failed to subscribe to notifications", "error", err)
			return
		}
		b.logger.Debug("subscribed to notifications", "topics", []string{upcomingProposerTopic})

		if time.Since(lastMsg) > 15*time.Minute {
			b.logger.Error("no messages received for 15 minutes, closing subscription")
			return
		}

		for {
			select {
			case <-ctx.Done():
				b.logger.Info("context done")
				return
			default:
			}

			msg, err := sub.Recv()
			if err != nil {
				b.logger.Error("failed to receive message", "error", err)
				goto RESTART
			}

			lastMsg = nowFunc()

			b.logger.Info("received message", "topic", msg.Topic)

			upcomingProposer, err := parseUpcomingProposer(msg)
			if err != nil {
				b.logger.Error("failed to parse upcoming proposer", "error", err)
				continue
			}
			b.logger.Debug("upcoming proposer", "upcomingProposer", upcomingProposer)

			// TODO: send upcoming proposer to different worker
			go func() {
				if err := b.Bid(ctx, big.NewInt(1000000000000000000), "0x"); err != nil {
					b.logger.Error("bid failed", "error", err)
				}
			}()
		}
	}()
	return done
}

func (b *BidderClient) Bid(
	ctx context.Context,
	bidAmount *big.Int,
	rawTx string,
) error {
	topo, err := b.topologyClient.GetTopology(ctx, &debugapiv1.EmptyMessage{})
	if err != nil {
		b.logger.Error("failed to get topology", "error", err)
		return err
	}

	providers := topo.Topology.Fields["connected_providers"].GetListValue()
	if providers == nil || len(providers.Values) == 0 {
		return errors.New("no connected providers")
	}

	tx, err := b.SelfETHTransfer()
	if err != nil {
		b.logger.Error("failed to create self ETH transfer transaction", "error", err)
		return err
	}

	// TODO: sanity check tx serialization
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		b.logger.Error("failed to marshal transaction", "error", err)
		return err
	}
	txString := hex.EncodeToString(txBytes)

	blkNumber, err := b.l1Client.BlockNumber(ctx)
	if err != nil {
		b.logger.Error("failed to get block number", "error", err)
		return err
	}

	pc, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
		TxHashes:            []string{},
		Amount:              bidAmount.String(),
		BlockNumber:         int64(blkNumber + 1),
		DecayStartTimestamp: nowFunc().UnixMilli(),
		DecayEndTimestamp:   nowFunc().Add(12 * time.Second).UnixMilli(),
		RevertingTxHashes:   []string{},
		RawTransactions:     []string{txString},
		SlashAmount:         big.NewInt(0).String(), // TODO: determine slash amount
	})
	if err != nil {
		b.logger.Error("failed to send bid", "error", err)
		return err
	}

	commitments := make([]*bidderapiv1.Commitment, 0)

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

		commitments = append(commitments, msg)

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
	b.logger.Info("bid succeeded, but not all commitments received", "commitments", len(commitments))

	return nil
}

func parseUpcomingProposer(msg *notificationsapiv1.Notification) (*upcomingProposer, error) {
	if msg.Topic != upcomingProposerTopic {
		return nil, ErrUnexpectedTopic
	}

	if msg == nil || msg.Value == nil {
		return nil, errors.New("notification msg is nil")
	}

	fields := msg.Value.Fields
	if fields == nil {
		return nil, errors.New("notification value fields are nil")
	}

	blsKey := fields["bls_key"].GetStringValue()
	if blsKey == "" {
		return nil, errors.New("failed to parse BLS key")
	}

	epochVal := fields["epoch"].GetNumberValue()
	if epochVal == 0 {
		return nil, errors.New("failed to parse epoch")
	}

	slotVal := fields["slot"].GetNumberValue()
	if slotVal == 0 {
		return nil, errors.New("failed to parse slot")
	}

	return &upcomingProposer{
		BLSKey: blsKey,
		Epoch:  uint64(epochVal),
		Slot:   uint64(slotVal),
	}, nil
}

type upcomingProposer struct {
	BLSKey string
	Epoch  uint64
	Slot   uint64
}
