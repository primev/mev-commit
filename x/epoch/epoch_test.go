package epoch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCalculator(t *testing.T) {
	// Test with Ethereum mainnet values
	calc := NewCalculator(MainnetGenesisTime, 12*time.Second, 32, 2)

	assert.Equal(t, time.Unix(MainnetGenesisTime, 0), calc.GenesisTime())
	assert.Equal(t, 12*time.Second, calc.slotDuration)
	assert.Equal(t, uint64(32), calc.slotsPerEpoch)
	assert.Equal(t, uint64(2), calc.epochsToOffset)
}

func TestCurrentSlot(t *testing.T) {
	// Create a calculator with a fixed genesis time for testing
	now := time.Now()
	genesisTime := now.Add(-100 * time.Second)
	calc := NewCalculator(genesisTime.Unix(), 10*time.Second, 32, 0)

	// With 10 second slots, and 100 seconds since genesis, we should be on slot 10
	expectedSlot := uint64(10)

	// We use a small delta to account for slight time differences during test execution
	actualSlot := calc.CurrentSlot()
	assert.InDelta(t, float64(expectedSlot), float64(actualSlot), 1.0)
}

func TestCurrentEpoch(t *testing.T) {
	// Create a calculator with a fixed genesis time for testing
	now := time.Now()
	genesisTime := now.Add(-1000 * time.Second)
	calc := NewCalculator(genesisTime.Unix(), 10*time.Second, 32, 0)

	// With 10 second slots, 32 slots per epoch, and 1000 seconds since genesis
	// We should be on epoch 3 (slot 100, which is epoch 3)
	expectedEpoch := uint64(3)

	// We use a small delta to account for slight time differences during test execution
	actualEpoch := calc.CurrentEpoch()
	assert.InDelta(t, float64(expectedEpoch), float64(actualEpoch), 1.0)
}

func TestTimeUntilNextEpoch(t *testing.T) {
	// Create a calculator with a fixed genesis time for testing
	now := time.Now()

	// Calculate genesis time such that we're 5 seconds into the current epoch
	offsetInEpoch := 5 * time.Second
	currentEpochStart := now.Add(-offsetInEpoch)
	genesisTime := currentEpochStart.Add(-time.Duration(10) * EpochDuration)

	calc := NewCalculator(genesisTime.Unix(), SlotDuration, SlotsPerEpoch, 0)

	// The time until the next epoch should be the epoch duration minus the offset
	expectedTimeUntil := EpochDuration - offsetInEpoch

	// Allow a 1 second delta for test execution time
	timeUntil := calc.TimeUntilNextEpoch()
	assert.InDelta(t, expectedTimeUntil.Seconds(), timeUntil.Seconds(), 1.0)
}

func TestEpochStartTime(t *testing.T) {
	// Use a fixed genesis time for testing
	genesisTime := time.Date(2020, 12, 1, 12, 0, 23, 0, time.UTC)
	calc := NewCalculator(genesisTime.Unix(), 12*time.Second, 32, 0)

	// Calculate start time for epoch 10
	epoch := uint64(10)
	expectedStartTime := genesisTime.Add(time.Duration(10*32*12) * time.Second)

	actualStartTime := calc.EpochStartTime(epoch)

	// Use Equal with time.UTC to ensure timezone doesn't affect comparison
	assert.Equal(t, expectedStartTime.UTC(), actualStartTime.UTC())
}

func TestTargetEpoch(t *testing.T) {
	// Create a calculator with a deterministic current epoch
	now := time.Now()

	// Set genesis time so current epoch is exactly 10
	calc := &Calculator{
		genesisTime:    now.Add(-time.Duration(10*32*12) * time.Second),
		slotDuration:   12 * time.Second,
		slotsPerEpoch:  32,
		epochsToOffset: 2,
	}

	// With epochsToOffset = 2, target epoch should be current epoch (10) - 2 = 8
	expected := uint64(8)
	assert.Equal(t, expected, calc.TargetEpoch())
}

func TestEpochsToFetch(t *testing.T) {
	// Create a calculator with a deterministic current epoch
	now := time.Now()

	// Set genesis time so current epoch is exactly w0
	calc := &Calculator{
		genesisTime:    now.Add(-time.Duration(10*32*12) * time.Second),
		slotDuration:   12 * time.Second,
		slotsPerEpoch:  32,
		epochsToOffset: 3,
	}

	// With epochsToOffset = 3, we should fetch epoch 7
	expected := []uint64{7}
	actual := calc.EpochsToFetch()
	assert.Equal(t, expected, actual)
}

func TestSlotToEpoch(t *testing.T) {
	calc := NewCalculator(MainnetGenesisTime, 12, 32, 0)

	testCases := []struct {
		slot          uint64
		expectedEpoch uint64
	}{
		{0, 0},     // First slot of epoch 0
		{31, 0},    // Last slot of epoch 0
		{32, 1},    // First slot of epoch 1
		{63, 1},    // Last slot of epoch 1
		{64, 2},    // First slot of epoch 2
		{100, 3},   // Random slot
		{1000, 31}, // Higher epoch
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			actualEpoch := calc.SlotToEpoch(tc.slot)
			assert.Equal(t, tc.expectedEpoch, actualEpoch)
		})
	}
}

func TestFormatEpoch(t *testing.T) {
	testCases := []struct {
		epoch        uint64
		expectedText string
	}{
		{0, "epoch_0"},
		{1, "epoch_1"},
		{100, "epoch_100"},
		{1000, "epoch_1000"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			actualText := FormatEpoch(tc.epoch)
			assert.Equal(t, tc.expectedText, actualText)
		})
	}
}

// Since we can't modify the Calculator struct to use a mockable time function,
// we'll skip the mock time test and use a different approach with fixed times
func TestWithFixedTime(t *testing.T) {
	// Use current time to create tests that are accurate when run
	startTime := time.Now()

	// Create a genesis time that's exactly 100 slots ago
	secondsAgo := 100 * SlotDuration

	genesisTime := startTime.Add(-secondsAgo)
	calc := NewCalculator(genesisTime.Unix(), SlotDuration, SlotsPerEpoch, 2)

	// We should be at exactly slot 100
	// Allow a small delta for test execution time
	assert.InDelta(t, float64(100), float64(calc.CurrentSlot()), 1.0)

	// With 32 slots per epoch, slot 100 is in epoch 3 (slot 100 / 32 = 3.125)
	assert.InDelta(t, float64(3), float64(calc.CurrentEpoch()), 0.1)

	// Target epoch is current epoch - 2
	assert.InDelta(t, float64(1), float64(calc.TargetEpoch()), 0.1)
}
