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
	provider := common.HexToAddress("0x456")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, provider, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, provider)
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
	provider := common.HexToAddress("0x456")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, provider, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, provider)
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
	provider := common.HexToAddress("0x456")

	val, err := s.GetBalance(bidder, provider)
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
	provider := common.HexToAddress("0x456")
	amount := big.NewInt(20)

	err := s.RefundBalanceIfExists(bidder, provider, amount)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "balance not found, no increase needed") {
		t.Fatalf("expected error containing 'balance not found, no increase needed', got %v", err)
	}

	err = s.SetBalance(bidder, provider, amount)
	if err != nil {
		t.Fatal(err)
	}

	increaseAmount := big.NewInt(5)
	err = s.RefundBalanceIfExists(bidder, provider, increaseAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, provider)
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
	provider := common.HexToAddress("0x456")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, provider, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, provider)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(depositedAmount) != 0 {
		t.Fatalf("expected %s, got %s", depositedAmount.String(), val.String())
	}

	err = s.DeleteBalance(bidder, provider)
	if err != nil {
		t.Fatal(err)
	}

	val, err = s.GetBalance(bidder, provider)
	if err != nil {
		t.Fatal(err)
	}
	if val != nil {
		t.Fatalf("expected nil, got %s", val.String())
	}
}
