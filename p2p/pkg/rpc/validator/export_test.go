// timings_test.go
package validatorapi

import (
	"context"
	"time"

	validatorapiv1 "github.com/primev/mev-commit/p2p/gen/go/validatorapi/v1"
)

func (s *Service) GenesisTime() time.Time {
	return s.ec.GenesisTime()
}

func (s *Service) ScheduleNotificationForSlot(epoch uint64, slot uint64, info *validatorapiv1.SlotInfo) {
	s.scheduleNotificationForSlot(epoch, slot, info)
}

func (s *Service) SetGenesisTime(genesisTime time.Time) {
	s.ec.SetGenesisTime(genesisTime)
}

func (s *Service) SetProcessEpoch(ctx context.Context, epoch uint64) {
	s.processEpoch(ctx, epoch)
}
