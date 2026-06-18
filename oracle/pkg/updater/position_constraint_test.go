package updater

import (
	"fmt"
	"testing"

	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
)

func buildTxnMap(size int) map[string]TxMetadata {
	txns := make(map[string]TxMetadata, size)
	for i := 0; i < size; i++ {
		txns[fmt.Sprintf("tx%d", i)] = TxMetadata{
			PosInBlock: i,
			Succeeded:  true,
			GasUsed:    100000,
			TotalGas:   uint64(size) * 100000,
		}
	}
	return txns
}

func TestCheckPositionConstraint_AbsoluteBottom(t *testing.T) {
	t.Parallel()

	txns := buildTxnMap(10)

	tests := []struct {
		name     string
		value    int32
		pos      int
		expected bool
	}{
		{"value=0, last position (9)", 0, 9, true},
		{"value=0, second to last (8)", 0, 8, false},
		{"value=0, first position (0)", 0, 0, false},
		{"value=1, position 8", 1, 8, true},
		{"value=1, position 9", 1, 9, true},
		{"value=1, position 7", 1, 7, false},
		{"value=2, position 7", 2, 7, true},
		{"value=2, position 6", 2, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := &bidderapiv1.PositionConstraint{
				Anchor: bidderapiv1.PositionConstraint_ANCHOR_BOTTOM,
				Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
				Value:  tt.value,
			}
			txMeta := txns[fmt.Sprintf("tx%d", tt.pos)]
			result := checkPositionConstraintSatisfied(constraint, txMeta, txns)
			if result != tt.expected {
				t.Errorf("pos=%d, value=%d: got %v, want %v", tt.pos, tt.value, result, tt.expected)
			}
		})
	}
}

func TestCheckPositionConstraint_AbsoluteTop(t *testing.T) {
	t.Parallel()

	txns := buildTxnMap(10)

	tests := []struct {
		name     string
		value    int32
		pos      int
		expected bool
	}{
		{"value=0, first position (0)", 0, 0, true},
		{"value=0, second position (1)", 0, 1, false},
		{"value=1, position 0", 1, 0, true},
		{"value=1, position 1", 1, 1, true},
		{"value=1, position 2", 1, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := &bidderapiv1.PositionConstraint{
				Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
				Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
				Value:  tt.value,
			}
			txMeta := txns[fmt.Sprintf("tx%d", tt.pos)]
			result := checkPositionConstraintSatisfied(constraint, txMeta, txns)
			if result != tt.expected {
				t.Errorf("pos=%d, value=%d: got %v, want %v", tt.pos, tt.value, result, tt.expected)
			}
		})
	}
}

func TestCheckPositionConstraint_PercentileBottom(t *testing.T) {
	t.Parallel()

	txns := buildTxnMap(10)

	tests := []struct {
		name     string
		value    int32
		pos      int
		expected bool
	}{
		{"value=0, last position (9)", 0, 9, true},
		{"value=0, position 8", 0, 8, false},
		{"value=10, position 8", 10, 8, true},
		{"value=10, position 9", 10, 9, true},
		{"value=10, position 7", 10, 7, false},
		{"value=20, position 7", 20, 7, true},
		{"value=20, position 6", 20, 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := &bidderapiv1.PositionConstraint{
				Anchor: bidderapiv1.PositionConstraint_ANCHOR_BOTTOM,
				Basis:  bidderapiv1.PositionConstraint_BASIS_PERCENTILE,
				Value:  tt.value,
			}
			txMeta := txns[fmt.Sprintf("tx%d", tt.pos)]
			result := checkPositionConstraintSatisfied(constraint, txMeta, txns)
			if result != tt.expected {
				t.Errorf("pos=%d, value=%d: got %v, want %v", tt.pos, tt.value, result, tt.expected)
			}
		})
	}
}

func TestCheckPositionConstraint_PercentileTop(t *testing.T) {
	t.Parallel()

	txns := buildTxnMap(10)

	tests := []struct {
		name     string
		value    int32
		pos      int
		expected bool
	}{
		{"value=0, first position (0)", 0, 0, true},
		{"value=0, second position (1)", 0, 1, false},
		{"value=10, position 1", 10, 1, true},
		{"value=10, position 0", 10, 0, true},
		{"value=10, position 2", 10, 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraint := &bidderapiv1.PositionConstraint{
				Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
				Basis:  bidderapiv1.PositionConstraint_BASIS_PERCENTILE,
				Value:  tt.value,
			}
			txMeta := txns[fmt.Sprintf("tx%d", tt.pos)]
			result := checkPositionConstraintSatisfied(constraint, txMeta, txns)
			if result != tt.expected {
				t.Errorf("pos=%d, value=%d: got %v, want %v", tt.pos, tt.value, result, tt.expected)
			}
		})
	}
}

func TestCheckPositionConstraint_Symmetry(t *testing.T) {
	t.Parallel()

	txns := buildTxnMap(10)

	// ANCHOR_TOP value=0 covers 1 position (pos 0)
	// ANCHOR_BOTTOM value=0 should also cover 1 position (pos 9)
	topCount := 0
	bottomCount := 0
	for i := 0; i < 10; i++ {
		txMeta := txns[fmt.Sprintf("tx%d", i)]
		topConstraint := &bidderapiv1.PositionConstraint{
			Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
			Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
			Value:  0,
		}
		if checkPositionConstraintSatisfied(topConstraint, txMeta, txns) {
			topCount++
		}
		bottomConstraint := &bidderapiv1.PositionConstraint{
			Anchor: bidderapiv1.PositionConstraint_ANCHOR_BOTTOM,
			Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
			Value:  0,
		}
		if checkPositionConstraintSatisfied(bottomConstraint, txMeta, txns) {
			bottomCount++
		}
	}
	if topCount != 1 {
		t.Errorf("ANCHOR_TOP value=0 matched %d positions, want 1", topCount)
	}
	if bottomCount != 1 {
		t.Errorf("ANCHOR_BOTTOM value=0 matched %d positions, want 1", bottomCount)
	}
}
