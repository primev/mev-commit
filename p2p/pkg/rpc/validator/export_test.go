// timings_test.go
package validatorapi

import (
	"context"
	"time"

	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
)

// SetTestTimings sets the global timing variables for testing and returns a cleanup function
// that restores the previous values.
func SetTestTimings(slotDuration time.Duration, epochSlots int, notifyOffset, fetchOffset time.Duration) func() {
	// Save the original values.
	origSlotDuration := SlotDuration
	origEpochSlots := EpochSlots
	origEpochDuration := EpochDuration
	origNotifyOffset := NotifyOffset
	origFetchOffset := FetchOffset

	// Override the globals.
	SlotDuration = slotDuration
	EpochSlots = epochSlots
	EpochDuration = SlotDuration * time.Duration(EpochSlots)
	NotifyOffset = notifyOffset
	FetchOffset = fetchOffset

	// Return a function to restore the original values.
	return func() {
		SlotDuration = origSlotDuration
		EpochSlots = origEpochSlots
		EpochDuration = origEpochDuration
		NotifyOffset = origNotifyOffset
		FetchOffset = origFetchOffset
	}
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

func (s *Service) SetProcessEpoch(ctx context.Context, epoch uint64) {
	s.processEpoch(ctx, epoch)
}
