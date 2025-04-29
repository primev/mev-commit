package epoch

import (
	"fmt"
	"time"
)

// Constants for Ethereum 2.0 timing
const (
	// Genesis time for mainnet in Unix timestamp
	MainnetGenesisTime = 1606824023

	// Slot duration in seconds
	SlotDuration = 12 * time.Second

	// Slots per epoch
	SlotsPerEpoch = uint64(32)

	// Epoch duration in seconds
	EpochDuration = SlotDuration * time.Duration(SlotsPerEpoch)
)

// Calculator handles epoch and slot calculations
type Calculator struct {
	genesisTime    time.Time
	slotDuration   time.Duration
	slotsPerEpoch  uint64
	epochDuration  time.Duration
	epochsToOffset uint64
}

// NewCalculator creates a new epoch calculator
func NewCalculator(
	genesisTimestamp int64,
	slotDuration time.Duration,
	slotsPerEpoch uint64,
	epochsToOffset uint64,
) *Calculator {
	epochDuration := slotDuration * time.Duration(slotsPerEpoch)
	return &Calculator{
		genesisTime:    time.Unix(genesisTimestamp, 0),
		slotDuration:   slotDuration,
		slotsPerEpoch:  slotsPerEpoch,
		epochDuration:  epochDuration,
		epochsToOffset: epochsToOffset,
	}
}

// NewMainnetCalculator creates a calculator with mainnet parameters
func NewMainnetCalculator(epochsToOffset uint64) *Calculator {
	return NewCalculator(
		MainnetGenesisTime,
		SlotDuration,
		SlotsPerEpoch,
		epochsToOffset,
	)
}

// CurrentSlot returns the current slot number
func (c *Calculator) CurrentSlot() uint64 {
	timeSinceGenesis := time.Since(c.genesisTime)
	return uint64(timeSinceGenesis / c.slotDuration)
}

// CurrentEpoch returns the current epoch number
func (c *Calculator) CurrentEpoch() uint64 {
	return c.CurrentSlot() / c.slotsPerEpoch
}

// FirstSlotOfEpoch returns the first slot of a given epoch
func (c *Calculator) FirstSlotOfEpoch(epoch uint64) uint64 {
	return epoch * c.slotsPerEpoch
}

// TimeUntilNextEpoch returns the duration until the next epoch starts
func (c *Calculator) TimeUntilNextEpoch() time.Duration {
	currentEpoch := c.CurrentEpoch()
	nextEpoch := currentEpoch + 1
	nextEpochTime := c.EpochStartTime(nextEpoch)
	return time.Until(nextEpochTime)
}

// EpochStartTime returns the start time of a given epoch
func (c *Calculator) EpochStartTime(epoch uint64) time.Time {
	durationSinceGenesis := time.Duration(epoch) * c.epochDuration
	return c.genesisTime.Add(durationSinceGenesis)
}

// SlotStartTime returns the start time of a given slot
func (c *Calculator) SlotStartTime(slot uint64) time.Time {
	timeSinceGenesis := time.Duration(slot) * c.slotDuration
	return c.genesisTime.Add(timeSinceGenesis)
}

// TargetEpoch returns the current epoch plus the configured offset
func (c *Calculator) TargetEpoch() uint64 {
	return c.CurrentEpoch() - c.epochsToOffset
}

// EpochsToFetch returns the epochs that should be fetched based on the current epoch
// and the configured offset
func (c *Calculator) EpochsToFetch() []uint64 {
	targetEpoch := c.CurrentEpoch() - c.epochsToOffset
	// Just return the target epoch in a slice with one element
	return []uint64{targetEpoch}
}

// SlotToEpoch converts a slot number to its corresponding epoch
func (c *Calculator) SlotToEpoch(slot uint64) uint64 {
	return slot / c.slotsPerEpoch
}

// GenesisTime returns the genesis time of the calculator
func (c *Calculator) GenesisTime() time.Time {
	return c.genesisTime
}

// SetGenesisTime sets the genesis time of the calculator
func (c *Calculator) SetGenesisTime(genesisTime time.Time) {
	c.genesisTime = genesisTime
}

// FormatEpoch returns a human-readable string for an epoch
func FormatEpoch(epoch uint64) string {
	return fmt.Sprintf("epoch_%d", epoch)
}
