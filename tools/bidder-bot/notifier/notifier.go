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

type Notifier struct {
	logger               *slog.Logger
	notificationsClient  notificationsapiv1.NotificationsClient
	proposerChan         chan *UpcomingProposer
	lastUpcomingProposer atomic.Pointer[UpcomingProposer]
}

func NewNotifier(
	logger *slog.Logger,
	notificationsClient notificationsapiv1.NotificationsClient,
	proposerChan chan *UpcomingProposer,
) *Notifier {
	return &Notifier{
		logger:              logger,
		notificationsClient: notificationsClient,
		proposerChan:        proposerChan,
	}
}

// TODO: unit tests validating buffering logic with the bidder worker
func (b *Notifier) Start(ctx context.Context) <-chan struct{} {
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
			if err := b.handleMsg(msg); err != nil {
				b.logger.Error("failed to handle message", "error", err)
				goto RESTART
			}
		}
	}()
	return done
}

func (b *Notifier) handleMsg(msg *notificationsapiv1.Notification) error {
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
	select {
	case b.proposerChan <- upcomingProposer:
		b.logger.Debug("sent upcoming proposer", "proposer", upcomingProposer)
	default:
		select {
		// TODO: confirm draining logic is correct
		case drainedProposer := <-b.proposerChan:
			b.logger.Warn("drained buffered upcoming proposer", "drained_proposer", drainedProposer)
		default:
		}
		b.proposerChan <- upcomingProposer
		b.logger.Warn("sent upcoming proposer after draining buffer", "proposer", upcomingProposer)
	}
	b.lastUpcomingProposer.Store(upcomingProposer)
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
