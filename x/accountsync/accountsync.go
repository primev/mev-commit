package accountsync

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type BalanceGetter interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

// AccountSync service is responsible for keeping track of the accounts and their balances.
type AccountSync struct {
	owner  common.Address
	client BalanceGetter
}

// NewAccountSync creates a new AccountSync service.
func NewAccountSync(address common.Address, client BalanceGetter) *AccountSync {
	return &AccountSync{
		owner:  address,
		client: client,
	}
}

// Subscribe call will start a goroutine that will periodically check the account balance and
// will notify the caller when the balance is below the threshold. The channel returned will be closed
// when the context is done or the account balance is below the threshold.
func (a *AccountSync) Subscribe(ctx context.Context, threshold *big.Int) <-chan struct{} {
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan struct{})
	go func() {
		defer ticker.Stop()
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				balance, err := a.client.BalanceAt(ctx, a.owner, nil)
				if err != nil {
					continue
				}
				if balance.Cmp(threshold) < 0 {
					return
				}
			}
		}
	}()

	return done
}
