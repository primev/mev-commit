package store_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/autodepositor/store"
	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

func TestStore(t *testing.T) {
	st := inmem.New()
	store := store.New(st)
	windows := []*big.Int{big.NewInt(1), big.NewInt(2)}

	t.Run("StoreDeposits", func(t *testing.T) {
		err := store.StoreDeposits(context.Background(), windows)
		if err != nil {
			t.Fatalf("StoreDeposits failed: %v", err)
		}

		for _, w := range windows {
			if !store.IsDepositMade(context.Background(), w) {
				t.Errorf("Deposit for window %s was not stored", w)
			}
		}
	})

	t.Run("ListDeposits", func(t *testing.T) {
		deposits, err := store.ListDeposits(context.Background(), big.NewInt(2))
		if err != nil {
			t.Fatalf("ListDeposits failed: %v", err)
		}

		expectedDeposits := []*big.Int{big.NewInt(1), big.NewInt(2)}
		if !reflect.DeepEqual(deposits, expectedDeposits) {
			t.Errorf("Expected deposits %+v, got %+v", expectedDeposits, deposits)
		}
	})

	t.Run("ClearDeposits", func(t *testing.T) {
		err := store.ClearDeposits(context.Background(), windows)
		if err != nil {
			t.Fatalf("ClearDeposits failed: %v", err)
		}

		for _, w := range windows {
			if store.IsDepositMade(context.Background(), w) {
				t.Errorf("Deposit for window %s was not cleared", w)
			}
		}
	})
}
