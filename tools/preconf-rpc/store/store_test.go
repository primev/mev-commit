package store_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/store"
)

func TestStore(t *testing.T) {
	t.Parallel()

	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Errorf("failed to close store: %v", err)
		}
	})

	t.Run("StorePreconfirmedTransaction", func(t *testing.T) {
		txn := types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(1000000000), // 1 Gwei
			21000,                  // gas limit
			big.NewInt(1000000000), // gas price
			nil,                    // no data
		)
		commitments := []*bidderapiv1.Commitment{
			{
				TxHashes:             []string{txn.Hash().Hex()},
				BidAmount:            big.NewInt(1000000000).String(),
				BlockNumber:          1,
				ReceivedBidDigest:    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ReceivedBidSignature: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				CommitmentDigest:     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				CommitmentSignature:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				DecayStartTimestamp:  time.Now().UnixMilli(),
				DecayEndTimestamp:    time.Now().Add(24 * time.Hour).UnixMilli(),
			},
			{
				TxHashes:             []string{txn.Hash().Hex()},
				BidAmount:            big.NewInt(1000000000).String(),
				BlockNumber:          1,
				ReceivedBidDigest:    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				ReceivedBidSignature: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				CommitmentDigest:     "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				CommitmentSignature:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				DecayStartTimestamp:  time.Now().UnixMilli(),
				DecayEndTimestamp:    time.Now().Add(24 * time.Hour).UnixMilli(),
			},
		}

		err := st.StorePreconfirmedTransaction(context.Background(), 1, txn, commitments)
		if err != nil {
			t.Errorf("failed to store preconfirmed transaction: %v", err)
		}

		storedTxn, storedCommitments, err := st.GetPreconfirmedTransaction(context.Background(), txn.Hash())
		if err != nil {
			t.Errorf("failed to get preconfirmed transaction: %v", err)
		}

		if txn.Hash().Hex() != storedTxn.Hash().Hex() {
			t.Errorf("expected transaction hash %s, got %s", txn.Hash().Hex(), storedTxn.Hash().Hex())
		}
		if len(storedCommitments) != len(commitments) {
			t.Errorf("expected %d commitments, got %d", len(commitments), len(storedCommitments))
		}

		for i, commitment := range commitments {
			if diff := cmp.Diff(commitment, storedCommitments[i], cmpopts.IgnoreUnexported(bidderapiv1.Commitment{})); diff != "" {
				t.Errorf("commitment mismatch (-want +got):\n%s", diff)
			}
		}
	})

	t.Run("Account Balance", func(t *testing.T) {
		address := common.HexToAddress("0x1234567890123456789012345678901234567890")
		initialBalance := big.NewInt(1000000000) // 1 Gwei

		err := st.AddBalance(context.Background(), address, initialBalance)
		if err != nil {
			t.Errorf("failed to add balance: %v", err)
		}

		if !st.HasBalance(context.Background(), address, initialBalance) {
			t.Errorf("expected balance %s, but has no balance", initialBalance.String())
		}

		// Check if the balance is correctly stored
		balance, err := st.GetBalance(context.Background(), address)
		if err != nil {
			t.Errorf("failed to get balance: %v", err)
		}
		if balance.Cmp(initialBalance) != 0 {
			t.Errorf("expected balance %s, got %s", initialBalance.String(), balance.String())
		}

		err = st.DeductBalance(context.Background(), address, initialBalance)
		if err != nil {
			t.Errorf("failed to deduct balance: %v", err)
		}

		if st.HasBalance(context.Background(), address, initialBalance) {
			t.Errorf("expected no balance after deduction, but still has %s", initialBalance.String())
		}
	})
}
