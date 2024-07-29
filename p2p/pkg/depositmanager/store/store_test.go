package store_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/p2p/pkg/depositmanager/store"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

func TestStore_SetBalance(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	windowNumber := big.NewInt(1)
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, windowNumber, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, windowNumber)
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
	windowNumber := big.NewInt(1)
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder, windowNumber, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalance(bidder, windowNumber)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(depositedAmount) != 0 {
		t.Fatalf("expected %s, got %s", depositedAmount.String(), val.String())
	}
}

func TestStore_ClearBalances(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	windowNumber := big.NewInt(1)
	bidder1 := common.HexToAddress("0x123")
	bidder2 := common.HexToAddress("0x456")
	depositedAmount := big.NewInt(10)

	err := s.SetBalance(bidder1, windowNumber, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}
	err = s.SetBalance(bidder2, windowNumber, depositedAmount)
	if err != nil {
		t.Fatal(err)
	}

	windows, err := s.ClearBalances(windowNumber)
	if err != nil {
		t.Fatal(err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1, got %d", len(windows))
	}

	val1, err := s.GetBalance(bidder1, windowNumber)
	if err != nil {
		t.Fatal(err)
	}
	if val1 != nil {
		t.Fatalf("expected nil, got %s", val1.String())
	}

	val2, err := s.GetBalance(bidder2, windowNumber)
	if err != nil {
		t.Fatal(err)
	}
	if val2 != nil {
		t.Fatalf("expected nil, got %s", val2.String())
	}
}

func TestStore_GetBalanceForBlock(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	windowNumber := big.NewInt(1)
	blockNumber := int64(10)
	amount := big.NewInt(20)

	err := s.SetBalanceForBlock(bidder, windowNumber, amount, blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalanceForBlock(bidder, windowNumber, blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(amount) != 0 {
		t.Fatalf("expected %s, got %s", amount.String(), val.String())
	}
}

func TestStore_SetBalanceForBlock(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	windowNumber := big.NewInt(1)
	blockNumber := int64(10)
	amount := big.NewInt(20)

	err := s.SetBalanceForBlock(bidder, windowNumber, amount, blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalanceForBlock(bidder, windowNumber, blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	if val.Cmp(amount) != 0 {
		t.Fatalf("expected %s, got %s", amount.String(), val.String())
	}
}

func TestStore_RefundBalanceForBlock(t *testing.T) {
	st := inmem.New()
	s := store.New(st)

	bidder := common.HexToAddress("0x123")
	windowNumber := big.NewInt(1)
	blockNumber := int64(10)
	amount := big.NewInt(20)

	err := s.SetBalanceForBlock(bidder, windowNumber, amount, blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	refundAmount := big.NewInt(5)
	err = s.RefundBalanceForBlock(bidder, windowNumber, refundAmount, blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	val, err := s.GetBalanceForBlock(bidder, windowNumber, blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	expectedAmount := new(big.Int).Add(amount, refundAmount)
	if val.Cmp(expectedAmount) != 0 {
		t.Fatalf("expected %s, got %s", expectedAmount.String(), val.String())
	}
}
