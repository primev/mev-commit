package bidder

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"sync/atomic"
	"time"

	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
)

const (
	epochNotificationTopic = "epoch_validators_opted_in"
	slotDuration           = 12 * time.Second
)

var (
	ErrNoEpochInfo       = errors.New("no epoch info available")
	ErrNoSlotInNextEpoch = errors.New("no slot available in next epoch")
)

var nowFunc = time.Now

type slotInfo struct {
	slot      uint64
	startTime time.Time
	blsKey    string
}

type epochInfo struct {
	epoch     uint64
	startTime time.Time
	slots     []slotInfo
}

type BlockNumberGetter interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type BidderClient struct {
	logger              *slog.Logger
	bidderClient        bidderapiv1.BidderClient
	topologyClient      debugapiv1.DebugServiceClient
	notificationsClient notificationsapiv1.NotificationsClient
	nextEpoch           atomic.Pointer[epochInfo]
	blkNumberGetter     BlockNumberGetter
}

func NewBidderClient(
	logger *slog.Logger,
	bidderClient bidderapiv1.BidderClient,
	topologyClient debugapiv1.DebugServiceClient,
	notificationsClient notificationsapiv1.NotificationsClient,
	blkNumberGetter BlockNumberGetter,
) *BidderClient {
	return &BidderClient{
		logger:              logger,
		bidderClient:        bidderClient,
		topologyClient:      topologyClient,
		notificationsClient: notificationsClient,
		blkNumberGetter:     blkNumberGetter,
	}
}

func (b *BidderClient) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		lastMsg := nowFunc()
	RESTART:
		sub, err := b.notificationsClient.Subscribe(ctx, &notificationsapiv1.SubscribeRequest{
			Topics: []string{epochNotificationTopic},
		})
		if err != nil {
			b.logger.Error("failed to subscribe to notifications", "error", err)
			return
		}

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

			if msg.Topic != epochNotificationTopic {
				b.logger.Error("unexpected topic", "topic", msg.Topic)
				continue
			}

			epoch, err := parseEpochInfo(msg)
			if err != nil {
				b.logger.Error("failed to parse epoch info", "error", err, "msg", msg)
				continue
			}

			b.nextEpoch.Store(epoch)
		}
	}()
	return done
}

func parseEpochInfo(msg *notificationsapiv1.Notification) (*epochInfo, error) {
	epochIdx := msg.Value.Fields["epoch"].GetNumberValue()
	if epochIdx == 0 {
		return nil, errors.New("failed to parse epoch index")
	}
	startTime := msg.Value.Fields["epoch_start_time"].GetNumberValue()
	if startTime == 0 {
		return nil, errors.New("failed to parse start time")
	}
	slots := msg.Value.Fields["slots"].GetListValue()
	if slots == nil {
		return nil, errors.New("failed to parse slots")
	}
	epoch := &epochInfo{
		epoch:     uint64(epochIdx),
		startTime: time.Unix(int64(startTime), 0),
	}
	for idx, slot := range slots.Values {
		slotIdx := slot.GetStructValue().Fields["slot"].GetNumberValue()
		if slotIdx == 0 {
			return nil, errors.New("failed to parse slot index")
		}
		blsKey := slot.GetStructValue().Fields["bls_key"].GetStringValue()
		if blsKey == "" {
			return nil, errors.New("failed to parse BLS key")
		}
		if slot.GetStructValue().Fields["opted_in"].GetBoolValue() {
			epoch.slots = append(epoch.slots, slotInfo{
				slot:      uint64(slotIdx),
				startTime: epoch.startTime.Add(time.Duration(idx) * slotDuration),
				blsKey:    blsKey,
			})
		}
	}

	return epoch, nil
}

func (b *BidderClient) Bid(
	ctx context.Context,
	bidAmount *big.Int,
	bridgeAmount *big.Int,
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

	nextSlot, err := b.getNextSlot()
	if err != nil {
		return err
	}

	bidTime := nextSlot.startTime.Add(-1 * time.Second)
	wait := bidTime.Sub(nowFunc())
	if wait > 0 {
		b.logger.Info("waiting for next slot", "wait", wait)
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	blkNumber, err := b.blkNumberGetter.BlockNumber(ctx)
	if err != nil {
		b.logger.Error("failed to get block number", "error", err)
		return err
	}

	pc, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
		Amount:              bidAmount.String(),
		BlockNumber:         int64(blkNumber + 1),
		RawTransactions:     []string{rawTx},
		DecayStartTimestamp: nowFunc().UnixMilli(),
		DecayEndTimestamp:   nowFunc().Add(12 * time.Second).UnixMilli(),
		SlashAmount:         bridgeAmount.String(),
	})
	if err != nil {
		b.logger.Error("failed to send bid", "error", err)
		return err
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
			return err
		}

		commitments = append(commitments, msg)
	}

	if len(commitments) == len(providers.Values) {
		b.logger.Info("all commitments received")
	} else {
		b.logger.Warn(
			"not all commitments received",
			"received", len(commitments),
			"expected", len(providers.Values),
		)
	}

	return nil
}

func (b *BidderClient) Estimate() (int64, error) {
	nextSlot, err := b.getNextSlot()
	if err != nil {
		return 0, err
	}

	return int64(nextSlot.startTime.Sub(nowFunc()).Seconds()), nil
}

func (b *BidderClient) getNextSlot() (slotInfo, error) {
	epochInfo := b.nextEpoch.Load()
	if epochInfo == nil {
		return slotInfo{}, ErrNoEpochInfo
	}

	now := nowFunc()
	for _, slot := range epochInfo.slots {
		if now.Before(slot.startTime) {
			return slot, nil
		}
	}

	return slotInfo{}, ErrNoSlotInNextEpoch
}
