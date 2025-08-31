package store_test

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/depositmanager/store"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

func TestStore_SetBalance(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(depositedAmount) != 0 {
		t.Fatalf("expected %s, got %s", depositedAmount.String(), val.String())
	}
}

func TestStore_GetBalance(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(depositedAmount) != 0 {
		t.Fatalf("expected %s, got %s", depositedAmount.String(), val.String())
	}
}

func TestStore_GetBalance_NoBalance(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")

	val, err := s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatalf("expected nil, got %s", val.String())
	}
}

func TestStore_RefundBalanceIfExists(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	amount := big.NewInt(20)

	err := s.RefundBalanceIfExists(bidder, amount)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "balance not found, no refund needed") {
		t.Fatalf("expected error containing 'balance not found, no refund needed', got %v", err)
	}

	err = s.SetBalance(bidder, amount)
	if err != nil {
		t.Fatal(err)
	}

	increaseAmount := big.NewInt(5)
	err = s.RefundBalanceIfExists(bidder, increaseAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	expectedAmount := new(big.Int).SetUint64(25)
	if val.Cmp(expectedAmount) != 0 {
		t.Fatalf("expected %s, got %s", expectedAmount.String(), val.String())
	}
}

func TestStore_DeleteBalance(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(depositedAmount) != 0 {
		t.Fatalf("expected %s, got %s", depositedAmount.String(), val.String())
	}

	err = s.DeleteBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}

	val, err = s.GetBalance(bidder)
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatalf("expected nil, got %s", val.String())
	}
}
