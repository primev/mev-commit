package notifier

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/tools/bidder-bot/bidder"
)

const (
	upcomingProposerTopic = "validator_opted_in"
	slotDuration          = 12 * time.Second
)

var (
	ErrUnexpectedTopic = errors.New("unexpected msg topic")
)

type SelectiveNotifier struct {
	logger               *slog.Logger
	notificationsClient  notificationsapiv1.NotificationsClient
	beaconClient         BeaconClient
	targetBlockChan      chan bidder.TargetBlock
	lastUpcomingProposer atomic.Pointer[UpcomingProposer]
}

type BeaconClient interface {
	GetPayloadDataForSlot(ctx context.Context, slot uint64) (blockNum uint64, timestamp uint64, err error)
}

func NewSelectiveNotifier(
	logger *slog.Logger,
	notificationsClient notificationsapiv1.NotificationsClient,
	beaconClient BeaconClient,
	targetBlockChan chan bidder.TargetBlock,
) *SelectiveNotifier {
	return &SelectiveNotifier{
		logger:              logger,
		notificationsClient: notificationsClient,
		beaconClient:        beaconClient,
		targetBlockChan:     targetBlockChan,
	}
}

func (b *SelectiveNotifier) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
	RESTART:
		sub, err := b.notificationsClient.Subscribe(ctx, &notificationsapiv1.SubscribeRequest{
			Topics: []string{upcomingProposerTopic},
		})
		if err != nil {
			b.logger.Error("failed to subscribe to notifications", "error", err)
			return
		}
		b.logger.Info("subscribed to notifications", "topics", []string{upcomingProposerTopic})
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
			b.logger.Info("received message", "topic", msg.Topic)
			if err := b.handleMsg(ctx, msg); err != nil {
				b.logger.Error("failed to handle message", "error", err)
				goto RESTART
			}
		}
	}()
	return done
}

func (b *SelectiveNotifier) handleMsg(ctx context.Context, msg *notificationsapiv1.Notification) error {
	upcomingProposer, err := parseUpcomingProposer(msg)
	if err != nil {
		b.logger.Error("failed to parse upcoming proposer", "error", err)
		return err
	}
	lastUpcomingProposer := b.lastUpcomingProposer.Load()
	if lastUpcomingProposer != nil && lastUpcomingProposer.Slot >= upcomingProposer.Slot {
		b.logger.Warn("received duplicate or outdated proposer notification. Msg will be dropped", "slot", upcomingProposer.Slot)
		return nil
	}

	// Upcoming proposer slot hasn't started yet, so query block number for upcoming proposer slot - 2
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	slotBeforeTarget := upcomingProposer.Slot - 2
	var targetBlock bidder.TargetBlock
	blockNumBeforeTarget, timestampBeforeTarget, err := b.beaconClient.GetPayloadDataForSlot(timeoutCtx, slotBeforeTarget)
	if err != nil {
		b.logger.Warn("failed to get block number for upcoming proposer slot - 2. This likely indicates a missed slot", "error", err)
		b.logger.Info("retrying with upcoming proposer slot - 3")
		slotBeforeTarget = upcomingProposer.Slot - 3
		blockNumBeforeTarget, timestampBeforeTarget, err = b.beaconClient.GetPayloadDataForSlot(timeoutCtx, slotBeforeTarget)
		if err != nil {
			b.logger.Error("failed to get block number for upcoming proposer slot - 3. No more retries", "error", err)
			return err
		}
	}
	targetBlock = bidder.TargetBlock{
		Num:  blockNumBeforeTarget + 2, // Same handling for either value of slotBeforeTarget
		Time: time.Unix(int64(timestampBeforeTarget), 0).Add(2 * slotDuration),
	}

	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case b.targetBlockChan <- targetBlock:
		b.logger.Debug("sent target block", "target_block_number", targetBlock.Num, "target_block_time", targetBlock.Time)
	case <-sendCtx.Done():
		b.logger.Warn("failed to send target block", "target_block_number", targetBlock.Num, "target_block_time", targetBlock.Time)
	}
	b.lastUpcomingProposer.Store(upcomingProposer)
	b.logger.Debug("updated lastUpcomingProposer", "proposer", upcomingProposer)
	return nil
}

type UpcomingProposer struct {
	BLSKey      string
	Epoch       uint64
	Slot        uint64
	CreationUTC time.Time
}

func parseUpcomingProposer(msg *notificationsapiv1.Notification) (*UpcomingProposer, error) {
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

	return &UpcomingProposer{
		BLSKey:      blsKey,
		Epoch:       uint64(epochVal),
		Slot:        uint64(slotVal),
		CreationUTC: time.Now(),
	}, nil
}
