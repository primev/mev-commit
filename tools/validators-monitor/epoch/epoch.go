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
	// Default number of months to look back
	DefaultLookbackMonths = 3
)

// Calculator handles epoch and slot calculations
type Calculator struct {
	genesisTime      time.Time
	slotDuration     time.Duration
	slotsPerEpoch    int
	epochsToOffset   int
	lookbackMonths   int // New field for months to look back
	maxEpochsToFetch int // New field to limit number of epochs fetched at once
}

// NewCalculator creates a new epoch calculator
func NewCalculator(genesisTimestamp int64, slotDurationSec, slotsPerEpoch, epochsToOffset int) *Calculator {
	return &Calculator{
		genesisTime:      time.Unix(genesisTimestamp, 0),
		slotDuration:     time.Duration(slotDurationSec) * time.Second,
		slotsPerEpoch:    slotsPerEpoch,
		epochsToOffset:   epochsToOffset,
		lookbackMonths:   DefaultLookbackMonths,
		maxEpochsToFetch: 10, // Default limit of epochs to fetch in one batch
	}
}

// SetLookbackMonths sets the number of months to look back for historical data
func (c *Calculator) SetLookbackMonths(months int) {
	c.lookbackMonths = months
}

// SetMaxEpochsToFetch sets the maximum number of epochs to fetch in one batch
func (c *Calculator) SetMaxEpochsToFetch(max int) {
	c.maxEpochsToFetch = max
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

// TimeToEpoch converts a time to its corresponding epoch
func (c *Calculator) TimeToEpoch(t time.Time) uint64 {
	// If time is before genesis, return 0
	if t.Before(c.genesisTime) {
		return 0
	}

	timeSinceGenesis := t.Sub(c.genesisTime)
	slot := uint64(timeSinceGenesis / c.slotDuration)
	return slot / uint64(c.slotsPerEpoch)
}

// TargetEpoch returns the current epoch plus the configured offset
func (c *Calculator) TargetEpoch() uint64 {
	return c.CurrentEpoch() - uint64(c.epochsToOffset)
}

// GetEpochForMonthsAgo returns the epoch from a specified number of months ago
func (c *Calculator) GetEpochForMonthsAgo(months int) uint64 {
	// Calculate the time 'months' ago
	timeMonthsAgo := time.Now().AddDate(0, -months, 0)
	return c.TimeToEpoch(timeMonthsAgo)
}

// EpochsToFetch returns the epochs that should be fetched based on the lookback period
// It returns epochs in batches to avoid overwhelming the API
func (c *Calculator) EpochsToFetch() []uint64 {
	// If we're still using the legacy method (current epoch minus offset)
	if c.lookbackMonths <= 0 {
		targetEpoch := c.CurrentEpoch() - uint64(c.epochsToOffset)
		return []uint64{targetEpoch}
	}

	// Calculate the earliest epoch we want to fetch (from X months ago)
	epochMonthsAgo := c.GetEpochForMonthsAgo(c.lookbackMonths)

	// Get the current target epoch
	targetEpoch := c.TargetEpoch()

	result := make([]uint64, 0, c.maxEpochsToFetch)
	for i := 0; i < c.maxEpochsToFetch && epochMonthsAgo+uint64(i) <= targetEpoch; i++ {
		result = append(result, epochMonthsAgo+uint64(i))
	}

	return result
}

// SlotToEpoch converts a slot number to its corresponding epoch
func (c *Calculator) SlotToEpoch(slot uint64) uint64 {
	return slot / uint64(c.slotsPerEpoch)
}

// FormatEpoch returns a human-readable string for an epoch
func FormatEpoch(epoch uint64) string {
	return fmt.Sprintf("epoch_%d", epoch)
}
