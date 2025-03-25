package bidder

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"sync"
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
	ErrNoEpochInfo          = errors.New("no epoch info available")
	ErrNoSlotInCurrentEpoch = errors.New("no slot available in current epoch")
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
	bigWg               sync.WaitGroup
	bidderClient        bidderapiv1.BidderClient
	topologyClient      debugapiv1.DebugServiceClient
	notificationsClient notificationsapiv1.NotificationsClient
	currentEpoch        atomic.Pointer[epochInfo]
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

			b.logger.Debug("received message", "msg", msg)

			if msg.Topic != epochNotificationTopic {
				b.logger.Error("unexpected topic", "topic", msg.Topic)
				continue
			}

			epoch, err := parseEpochInfo(msg)
			if err != nil {
				b.logger.Error("failed to parse epoch info", "error", err, "msg", msg)
				continue
			}

			b.currentEpoch.Store(epoch)
			b.logger.Info("current epoch info updated", "epoch", epoch.epoch)
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
	baseSlot := epochIdx * 32
	for _, slot := range slots.Values {
		slotIdx := slot.GetStructValue().Fields["slot"].GetNumberValue()
		if slotIdx == 0 {
			return nil, errors.New("failed to parse slot index")
		}
		if slotIdx < baseSlot || slotIdx >= baseSlot+32 {
			return nil, errors.New("slot index out of range")
		}
		blsKey := slot.GetStructValue().Fields["bls_key"].GetStringValue()
		if blsKey == "" {
			return nil, errors.New("failed to parse BLS key")
		}
		idx := slotIdx - baseSlot
		epoch.slots = append(epoch.slots, slotInfo{
			slot:      uint64(slotIdx),
			startTime: epoch.startTime.Add(time.Duration(idx) * slotDuration),
			blsKey:    blsKey,
		})
	}

	return epoch, nil
}

type BidStatusType int

const (
	BidStatusNoOfProviders BidStatusType = iota
	BidStatusWaitSecs
	BidStatusAttempted
	BidStatusSucceeded
	BidStatusFailed
)

type BidStatus struct {
	Type BidStatusType
	Arg1 int
	Arg2 string
}

func (b *BidderClient) Bid(
	ctx context.Context,
	bidAmount *big.Int,
	bridgeAmount *big.Int,
	rawTx string,
) (chan BidStatus, error) {
	topo, err := b.topologyClient.GetTopology(ctx, &debugapiv1.EmptyMessage{})
	if err != nil {
		b.logger.Error("failed to get topology", "error", err)
		return nil, err
	}

	providers := topo.Topology.Fields["connected_providers"].GetListValue()
	if providers == nil || len(providers.Values) == 0 {
		return nil, errors.New("no connected providers")
	}

	// Channel length chosen is 3 so that sending the bid is not blocked by the first
	// status message.
	res := make(chan BidStatus, 3)
	b.bigWg.Add(1)
	go func() {
		defer close(res)
		defer b.bigWg.Done()

		res <- BidStatus{Type: BidStatusNoOfProviders, Arg1: len(providers.Values)}

		nextSlot, err := b.getNextSlot()
		if err != nil {
			b.logger.Error("failed to get next slot", "error", err)
			res <- BidStatus{Type: BidStatusFailed, Arg2: err.Error()}
			return
		}

		bidTime := nextSlot.startTime.Add(-1 * time.Second)
		wait := bidTime.Sub(nowFunc())
		res <- BidStatus{Type: BidStatusWaitSecs, Arg1: int(wait.Seconds())}

		if wait > 0 {
			b.logger.Info("waiting for next slot", "wait", wait)
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				res <- BidStatus{Type: BidStatusFailed, Arg2: ctx.Err().Error()}
				return
			}
		}

		blkNumber, err := b.blkNumberGetter.BlockNumber(ctx)
		if err != nil {
			b.logger.Error("failed to get block number", "error", err)
			res <- BidStatus{Type: BidStatusFailed, Arg2: err.Error()}
			return
		}

		res <- BidStatus{Type: BidStatusAttempted, Arg1: int(blkNumber + 1)}

		pc, err := b.bidderClient.SendBid(ctx, &bidderapiv1.Bid{
			Amount:              bidAmount.String(),
			BlockNumber:         int64(blkNumber + 1),
			RawTransactions:     []string{rawTx},
			DecayStartTimestamp: nowFunc().UnixMilli(),
			DecayEndTimestamp:   nowFunc().Add(12 * time.Second).UnixMilli(),
		})
		if err != nil {
			b.logger.Error("failed to send bid", "error", err)
			res <- BidStatus{Type: BidStatusFailed, Arg2: err.Error()}
			return
		}

		commitments := make([]*bidderapiv1.Commitment, 0)
		for {
			select {
			case <-ctx.Done():
				res <- BidStatus{Type: BidStatusFailed, Arg2: ctx.Err().Error()}
				return
			default:
			}

			msg, err := pc.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				b.logger.Error("failed to receive commitment", "error", err)
				res <- BidStatus{Type: BidStatusFailed, Arg2: err.Error()}
				return
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
		res <- BidStatus{Type: BidStatusSucceeded, Arg1: len(commitments)}
	}()

	return res, nil
}

func (b *BidderClient) Estimate() (int64, error) {
	nextSlot, err := b.getNextSlot()
	if err != nil {
		return 0, err
	}

	return int64(nextSlot.startTime.Sub(nowFunc()).Seconds()), nil
}

func (b *BidderClient) getNextSlot() (slotInfo, error) {
	epochInfo := b.currentEpoch.Load()
	if epochInfo == nil {
		return slotInfo{}, ErrNoEpochInfo
	}

	now := nowFunc()
	for _, slot := range epochInfo.slots {
		if now.Before(slot.startTime) {
			return slot, nil
		}
	}

	return slotInfo{}, ErrNoSlotInCurrentEpoch
}
