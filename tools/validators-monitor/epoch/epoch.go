// internal/epoch/epoch.go
package epoch

import (
	"fmt"
	"time"
)

// Constants for Ethereum 2.0 timing
const (
	// Genesis time for mainnet in Unix timestamp
	MainnetGenesisTime = 1606824023
)

// Calculator handles epoch and slot calculations
type Calculator struct {
	genesisTime    time.Time
	slotDuration   time.Duration
	slotsPerEpoch  int
	epochsToOffset int
}

// NewCalculator creates a new epoch calculator
func NewCalculator(genesisTimestamp int64, slotDurationSec, slotsPerEpoch, epochsToOffset int) *Calculator {
	return &Calculator{
		genesisTime:    time.Unix(genesisTimestamp, 0),
		slotDuration:   time.Duration(slotDurationSec) * time.Second,
		slotsPerEpoch:  slotsPerEpoch,
		epochsToOffset: epochsToOffset,
	}
}

// CurrentSlot returns the current slot number
func (c *Calculator) CurrentSlot() uint64 {
	timeSinceGenesis := time.Since(c.genesisTime)
	return uint64(timeSinceGenesis / c.slotDuration)
}

// CurrentEpoch returns the current epoch number
func (c *Calculator) CurrentEpoch() uint64 {
	return c.CurrentSlot() / uint64(c.slotsPerEpoch)
}

// TimeUntilNextEpoch returns the duration until the next epoch starts
func (c *Calculator) TimeUntilNextEpoch() time.Duration {
	currentSlot := c.CurrentSlot()
	currentEpochStartSlot := (currentSlot / uint64(c.slotsPerEpoch)) * uint64(c.slotsPerEpoch)
	nextEpochStartSlot := currentEpochStartSlot + uint64(c.slotsPerEpoch)

	slotsSinceGenesis := time.Duration(nextEpochStartSlot) * c.slotDuration
	nextEpochTime := c.genesisTime.Add(slotsSinceGenesis)

	return time.Until(nextEpochTime)
}

// EpochStartTime returns the start time of a given epoch
func (c *Calculator) EpochStartTime(epoch uint64) time.Time {
	epochStartSlot := epoch * uint64(c.slotsPerEpoch)
	timeSinceGenesis := time.Duration(epochStartSlot) * c.slotDuration
	return c.genesisTime.Add(timeSinceGenesis)
}

// TargetEpoch returns the current epoch plus the configured offset
func (c *Calculator) TargetEpoch() uint64 {
	return c.CurrentEpoch() - uint64(c.epochsToOffset)
}

// EpochsToFetch returns the epochs that should be fetched based on the current epoch
// and the configured offset
func (c *Calculator) EpochsToFetch() []uint64 {
	targetEpoch := c.CurrentEpoch() - uint64(c.epochsToOffset)
	// Just return the target epoch in a slice with one element
	return []uint64{targetEpoch}
}

// SlotToEpoch converts a slot number to its corresponding epoch
func (c *Calculator) SlotToEpoch(slot uint64) uint64 {
	return slot / uint64(c.slotsPerEpoch)
}

// FormatEpoch returns a human-readable string for an epoch
func FormatEpoch(epoch uint64) string {
	return fmt.Sprintf("epoch_%d", epoch)
}
