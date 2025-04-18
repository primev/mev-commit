package notifier

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
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
	beaconClient         *beaconClient
	targetBlockNumChan   chan uint64
	lastUpcomingProposer atomic.Pointer[UpcomingProposer]
}

func NewSelectiveNotifier(
	logger *slog.Logger,
	notificationsClient notificationsapiv1.NotificationsClient,
	beaconRPCUrl string,
	targetBlockNumChan chan uint64,
) *SelectiveNotifier {
	return &SelectiveNotifier{
		logger:              logger,
		notificationsClient: notificationsClient,
		beaconClient:        newBeaconClient(beaconRPCUrl, logger.With("component", "beacon_client")),
		targetBlockNumChan:  targetBlockNumChan,
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
	upcomingSlotMinusTwo := upcomingProposer.Slot - 2
	upcomingSlotMinusTwoBlockNum, err := b.beaconClient.getBlockNumForSlot(timeoutCtx, upcomingSlotMinusTwo)
	if err != nil {
		b.logger.Error("failed to get block number for upcoming proposer slot - 2", "error", err)
		return err
	}

	// Assume the two slots before upcoming proposer slot are NOT missed
	targetBlockNum := upcomingSlotMinusTwoBlockNum + 2

	sendCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	select {
	case b.targetBlockNumChan <- targetBlockNum:
		b.logger.Debug("sent target block number", "target_block_number", targetBlockNum)
	case <-sendCtx.Done():
		b.logger.Warn("failed to send target block number", "target_block_number", targetBlockNum)
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
