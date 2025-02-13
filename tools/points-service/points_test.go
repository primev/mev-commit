package main

import "testing"

// TestComputePointsForMonths verifies totalPoints & preSixMonthPoints for various month scenarios.
func TestComputePointsForMonths(t *testing.T) {
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
			name:    "1 month => (1000,1000)",
			months:  1,
			wantTot: 1000,
			wantPre: 1000,
		},
		{
			name:    "2 months => (2270,2000)",
			months:  2,
			wantTot: 2270,
			wantPre: 2000,
		},
		{
			name:    "5 months => (7670,5000)",
			months:  5,
			wantTot: 7670,
			wantPre: 5000,
		},
		{
			name:    "6 months => (10000,6000)",
			months:  6,
			wantTot: 10000,
			wantPre: 10000,
		},
		{
			name:    "7 months => (11983,6000)",
			months:  7,
			wantTot: 11983,
			wantPre: 11000,
		},
		{
			name:    "11 months => (25317,6000)",
			months:  11,
			wantTot: 25317,
			wantPre: 15000,
		},
		{
			name:    "12 months => (30000,6000)",
			months:  12,
			wantTot: 30000,
			wantPre: 30000,
		},
		{
			name:    "13 months => (34683,6000)",
			months:  13,
			wantTot: 34683,
			wantPre: 31000,
		},
		{
			name:    "14 months => (39367,6000)",
			months:  14,
			wantTot: 39367,
			wantPre: 32000,
		},
		{
			name:    "17 months => (53417,6000)",
			months:  17,
			wantTot: 53417,
			wantPre: 35000,
		},
		{
			name:    "18 months => (58100,6000)",
			months:  18,
			wantTot: 58100,
			wantPre: 58100,
		},
		{
			name:    "19 months => (58100,6000) beyond 18 => cap",
			months:  19,
			wantTot: 58100,
			wantPre: 58100,
		},
	}

	for _, tc := range testCases {
		blocksActive := tc.months * blocksInOneMonth
		gotTot, gotPre := computePointsForMonths(blocksActive)

		if gotTot != tc.wantTot || gotPre != tc.wantPre {
			t.Errorf("%s: months=%d => got (tot=%d, pre=%d), want (tot=%d, pre=%d)",
				tc.name, tc.months, gotTot, gotPre, tc.wantTot, tc.wantPre)
		}
	}
}
