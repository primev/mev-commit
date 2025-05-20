package accountsync_test

import (
	"context"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/x/accountsync"
)

type mockBalanceGetter struct {
	balance atomic.Pointer[big.Int]
}

func (m *mockBalanceGetter) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return m.balance.Load(), nil
}

func TestAccountSync(t *testing.T) {
	t.Parallel()

	// Create a new mock balance getter.
	mockBalanceGetter := &mockBalanceGetter{}
	mockBalanceGetter.balance.Store(big.NewInt(10))

	// Create a new account sync service.
	accountSync := accountsync.NewAccountSync(common.HexToAddress("0x123"), mockBalanceGetter)

	// Create a new channel to receive the notification.
	done := accountSync.Subscribe(context.Background(), big.NewInt(100))

	// Wait for the notification.
	<-done

	mockBalanceGetter.balance.Store(big.NewInt(150))

	// Create a new channel to receive the notification.
	done = accountSync.Subscribe(context.Background(), big.NewInt(100))

	select {
	case <-done:
		t.Fatal("expected the channel to be open")
	case <-time.After(2 * time.Second):
		break
	}

	mockBalanceGetter.balance.Store(big.NewInt(50))
	<-done

	// Create a new channel to receive the notification.
	mockBalanceGetter.balance.Store(big.NewInt(150))
	ctx, cancel := context.WithCancel(context.Background())
	done = accountSync.Subscribe(ctx, big.NewInt(100))

	// Cancel the context.
	cancel()
	<-done
}
