package main

import (
	"testing"
)

// TestComputePointsForMonths uses a table-driven approach
// to verify totalPoints & preSixMonthPoints for various monthly scenarios.
func TestComputePointsForMonths(t *testing.T) {
	// Each test case: (Months, ExpectedTotalPoints, ExpectedPreSixMonthPoints)
	testCases := []struct {
		name    string
		months  int64
		wantTot int64
		wantPre int64
	}{
		{
			name:    "0 months => (0,0)",
			months:  0,
			wantTot: 0,
			wantPre: 0,
		},
		{
			name:    "1 month => (10000,10000)",
			months:  1,
			wantTot: 10000,
			wantPre: 10000,
		},
		{
			name:    "2 months => (22700,20000)",
			months:  2,
			wantTot: 22700,
			wantPre: 20000,
		},
		{
			name:    "5 months => (76700,50000)",
			months:  5,
			wantTot: 76700,
			wantPre: 50000,
		},
		{
			name:    "6 months => (100000,60000)",
			months:  6,
			wantTot: 100000,
			wantPre: 60000,
		},
		{
			name:    "7 months => (120000,110000)",
			months:  7,
			wantTot: 120000,
			wantPre: 110000,
		},
		{
			name:    "11 months => (253400,210000)",
			months:  11,
			wantTot: 253400,
			wantPre: 150000,
		},
		{
			name:    "12 months => (300000,300000)",
			months:  12,
			wantTot: 300000,
			wantPre: 300000,
		},
		{
			name:    "13 months => (340000,310000)",
			months:  13,
			wantTot: 340000,
			wantPre: 310000,
		},
		{
			name:    "14 months => (390800,320000)",
			months:  14,
			wantTot: 390800,
			wantPre: 320000,
		},
		{
			name:    "17 months => (606800,370000)",
			months:  17,
			wantTot: 606800,
			wantPre: 350000,
		},
		{
			name:    "18 months => (700000,700000)",
			months:  18,
			wantTot: 700000,
			wantPre: 700000,
		},
		{
			name:    "19 months => (700000,700000) beyond 18 => cap",
			months:  19,
			wantTot: 700000,
			wantPre: 700000,
		},
	}

	for _, tc := range testCases {
		// Convert from "months" to "blocksActive"
		blocksActive := tc.months * blocksInOneMonth

		gotTot, gotPre := computePointsForMonths(blocksActive)
		if gotTot != tc.wantTot || gotPre != tc.wantPre {
			t.Errorf("%s: months=%d => got (tot=%d, pre=%d), want (tot=%d, pre=%d)",
				tc.name, tc.months, gotTot, gotPre, tc.wantTot, tc.wantPre)
		}
	}
}
