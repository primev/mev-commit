package store_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	"github.com/primev/mev-commit/p2p/pkg/storage"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

func TestStore_AddCommitment(t *testing.T) {
	st := store.New(inmem.New())
	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("commitment"),
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(commitments) != 1 {
		t.Fatalf("expected 1 commitment, got %d", len(commitments))
	}

	if !bytes.Equal(commitments[0].Commitment, []byte("commitment")) {
		t.Fatalf("expected commitment, got %s", commitments[0].Commitment)
	}
}

func TestStore_ClearBlockNumber(t *testing.T) {
	st := store.New(inmem.New())
	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("commitment"),
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.ClearBlockNumber(1)
	if err != nil {
		t.Fatal(err)
	}

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(commitments) != 0 {
		t.Fatalf("expected 0 commitments, got %d", len(commitments))
	}
}

func TestStore_ClearCommitmentIndex(t *testing.T) {
	inmemstore := inmem.New()
	st := store.New(inmemstore)

	digest := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	index := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetCommitmentIndexByDigest(digest, index)
	if err != nil {
		t.Fatal(err)
	}

	err = st.ClearCommitmentIndexes(2)
	if err != nil {
		t.Fatal(err)
	}

	entries := 0
	err = inmemstore.WalkPrefix(store.CmtIndexNS, func(_ string, _ []byte) bool {
		entries++
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
	if entries != 0 {
		t.Fatalf("expected 0 entries, got %d", entries)
	}

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(commitments) != 0 {
		t.Fatalf("expected 0 commitments, got %d", len(commitments))
	}

	err = inmemstore.WalkPrefix(store.IndexToDigestNS, func(_ string, _ []byte) bool {
		entries++
		return false
	})
	if err != nil {
		t.Fatal(err)
	}
	if entries != 0 {
		t.Fatalf("expected 0 entries, got %d", entries)
	}
}

func TestStore_SetCommitmentIndexByDigest(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))
	index := [32]byte{}
	copy(index[:], []byte("index"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetCommitmentIndexByDigest(digest, index)
	if err != nil {
		t.Fatal(err)
	}

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(commitments) != 1 {
		t.Fatalf("expected 1 commitment, got %d", len(commitments))
	}

	if !bytes.Equal(commitments[0].CommitmentIndex, index[:]) {
		t.Fatalf("expected index, got %s", commitments[0].CommitmentIndex)
	}
}

func TestStore_AddWinner(t *testing.T) {
	st := store.New(inmem.New())
	winner := &store.BlockWinner{
		BlockNumber: 1,
		Winner:      common.HexToAddress("0x123"),
	}
	err := st.AddWinner(winner)
	if err != nil {
		t.Fatal(err)
	}

	winners, err := st.BlockWinners()
	if err != nil {
		t.Fatal(err)
	}

	if len(winners) != 1 {
		t.Fatalf("expected 1 winner, got %d", len(winners))
	}

	if winners[0].BlockNumber != 1 {
		t.Fatalf("expected block number 1, got %d", winners[0].BlockNumber)
	}

	if winners[0].Winner != common.HexToAddress("0x123") {
		t.Fatalf("expected winner 0x123, got %s", winners[0].Winner.Hex())
	}
}

func TestStore_GetCommitments_Order(t *testing.T) {
	st := store.New(inmem.New())

	commitment1 := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("commitment1"),
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "300",
			},
		},
	}
	commitment2 := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("commitment2"),
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "200",
			},
		},
	}
	commitment3 := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: []byte("commitment3"),
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}

	err := st.AddCommitment(commitment1)
	if err != nil {
		t.Fatal(err)
	}
	err = st.AddCommitment(commitment3)
	if err != nil {
		t.Fatal(err)
	}
	err = st.AddCommitment(commitment2)
	if err != nil {
		t.Fatal(err)
	}

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(commitments) != 3 {
		t.Fatalf("expected 3 commitments, got %d", len(commitments))
	}

	expectedOrder := []string{"commitment1", "commitment2", "commitment3"}
	for i, commitment := range commitments {
		expectedCommitment := expectedOrder[i]
		if !bytes.Equal(commitment.Commitment, []byte(expectedCommitment)) {
			t.Fatalf("expected commitment %s at position %d, got %s", expectedCommitment, i, commitment.Commitment)
		}
	}
}

func TestStore_GetCommitmentByDigest(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	foundCommitment, err := st.GetCommitmentByDigest(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(foundCommitment.Commitment, digest[:]) {
		t.Fatalf("expected commitment %x, got %x", digest, foundCommitment.Commitment)
	}
}

func TestStore_GetCommitmentByDigest_NotFound(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("nonexistent"))

	_, err := st.GetCommitmentByDigest(digest[:])
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expectedErr := storage.ErrKeyNotFound
	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestStore_SetStatus(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetStatus(1, "100", digest[:], store.CommitmentStatusOpened, "details")
	if err != nil {
		t.Fatal(err)
	}

	foundCommitment, err := st.GetCommitmentByDigest(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	if foundCommitment.Status != store.CommitmentStatusOpened {
		t.Fatalf("expected status %s, got %s", store.CommitmentStatusOpened, foundCommitment.Status)
	}
	if foundCommitment.Details != "details" {
		t.Fatalf("expected details 'details', got '%s'", foundCommitment.Details)
	}
}

func TestStore_UpdateSettlementSettled(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	index := [32]byte{}
	copy(index[:], []byte("index"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetCommitmentIndexByDigest(digest, index)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdateSettlement(index[:], false)
	if err != nil {
		t.Fatal(err)
	}

	foundCommitment, err := st.GetCommitmentByDigest(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	if foundCommitment.Status != store.CommitmentStatusSettled {
		t.Fatalf("expected settlement status %s, got %s", store.CommitmentStatusSettled, foundCommitment.Status)
	}
}

func TestStore_UpdateSettlementSlashed(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	index := [32]byte{}
	copy(index[:], []byte("index"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.SetCommitmentIndexByDigest(digest, index)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdateSettlement(index[:], true)
	if err != nil {
		t.Fatal(err)
	}

	foundCommitment, err := st.GetCommitmentByDigest(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	if foundCommitment.Status != store.CommitmentStatusSlashed {
		t.Fatalf("expected settlement status %s, got %s", store.CommitmentStatusSlashed, foundCommitment.Status)
	}
}

func TestStore_UpdatePayment(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	commitment := &store.Commitment{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
		},
		PreConfirmation: &preconfpb.PreConfirmation{
			Bid: &preconfpb.Bid{
				BlockNumber: 1,
				BidAmount:   "100",
			},
		},
	}
	err := st.AddCommitment(commitment)
	if err != nil {
		t.Fatal(err)
	}

	err = st.UpdatePayment(digest[:], "80", "20")
	if err != nil {
		t.Fatal(err)
	}

	foundCommitment, err := st.GetCommitmentByDigest(digest[:])
	if err != nil {
		t.Fatal(err)
	}

	if foundCommitment.Payment != "80" {
		t.Fatalf("expected payment '80', got '%s'", foundCommitment.Payment)
	}
	if foundCommitment.Refund != "20" {
		t.Fatalf("expected refund '20', got '%s'", foundCommitment.Refund)
	}
}
