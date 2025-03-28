// timings_test.go
package validatorapi

import (
	"context"
	"time"

	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
)

// SetTestTimings sets the slot duration, epoch slots and proposer notify offset for testing.
func (s *Service) SetTestTimings(slotDuration time.Duration, epochSlots int, notifyOffset time.Duration) {
	s.slotDuration = slotDuration
	s.epochSlots = epochSlots
	s.proposerNotifyOffset = notifyOffset
}

func (s *Service) GenesisTime() time.Time {
	return s.genesisTime
}

func (s *Service) ScheduleNotificationForSlot(epoch uint64, slot uint64, info *validatorapiv1.SlotInfo) {
	s.scheduleNotificationForSlot(epoch, slot, info)
}

func (s *Service) SetGenesisTime(genesisTime time.Time) {
	s.genesisTime = genesisTime
}

func (s *Service) SetProcessEpoch(ctx context.Context, epoch uint64, epochTime int64) {
	s.processEpoch(ctx, epoch, epochTime)
}
