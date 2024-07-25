package store_test

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

func TestStore_AddCommitment(t *testing.T) {
	st := store.New(inmem.New())
	commitment := &store.EncryptedPreConfirmationWithDecrypted{
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

	commitments, err := st.GetCommitments(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(commitments) != 1 {
		t.Fatalf("expected 1 commitment, got %d", len(commitments))
	}

	if !bytes.Equal(commitments[0].EncryptedPreConfirmation.Commitment, []byte("commitment")) {
		t.Fatalf("expected commitment, got %s", commitments[0].EncryptedPreConfirmation.Commitment)
	}
}

func TestStore_ClearBlockNumber(t *testing.T) {
	st := store.New(inmem.New())
	commitment := &store.EncryptedPreConfirmationWithDecrypted{
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

func TestStore_DeleteCommitmentByDigest(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))

	commitment := &store.EncryptedPreConfirmationWithDecrypted{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
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

	err = st.DeleteCommitmentByDigest(1, digest)
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

func TestStore_SetCommitmentIndexByDigest(t *testing.T) {
	st := store.New(inmem.New())
	digest := [32]byte{}
	copy(digest[:], []byte("commitment"))
	index := [32]byte{}
	copy(index[:], []byte("index"))

	commitment := &store.EncryptedPreConfirmationWithDecrypted{
		EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
			Commitment: digest[:],
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

	if !bytes.Equal(commitments[0].EncryptedPreConfirmation.CommitmentIndex, index[:]) {
		t.Fatalf("expected index, got %s", commitments[0].EncryptedPreConfirmation.CommitmentIndex)
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
