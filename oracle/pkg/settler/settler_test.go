package settler_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"path"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-commit/oracle/pkg/keysigner"
	"github.com/primevprotocol/mev-commit/oracle/pkg/settler"
)

type testRegister struct {
	currentNonce         atomic.Int64
	pendingTxns          atomic.Int32
	settlementChan       chan settler.Settlement
	returnsChan          chan settler.Return
	mu                   sync.Mutex
	settlementsInitiated [][]byte
	settlementsCompleted atomic.Int32
}

func (t *testRegister) LastNonce() (int64, error) {
	return t.currentNonce.Load(), nil
}

func (t *testRegister) PendingTxnCount() (int, error) {
	return int(t.pendingTxns.Load()), nil
}

func (t *testRegister) SubscribeSettlements(ctx context.Context) <-chan settler.Settlement {
	sc := make(chan settler.Settlement)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case s := <-t.settlementChan:
				sc <- s
			}
		}
	}()

	return sc
}

func (t *testRegister) SubscribeReturns(ctx context.Context, _ int) <-chan settler.Return {
	rc := make(chan settler.Return)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case r := <-t.returnsChan:
				rc <- r
			}
		}
	}()

	return rc
}

func (t *testRegister) SettlementInitiated(ctx context.Context, commitmentIdx [][]byte, txHash common.Hash, nonce uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.settlementsInitiated = append(t.settlementsInitiated, commitmentIdx...)
	return nil
}

func (t *testRegister) MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error) {
	t.settlementsCompleted.Store(int32(nonce))
	return 1, nil
}

func (t *testRegister) settlementsInitiatedCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.settlementsInitiated)
}

type testOracle struct {
	key            *ecdsa.PrivateKey
	mu             sync.Mutex
	commitmentIdxs [][32]byte
	bidIDs         [][32]byte
	slashCount     atomic.Int32
	rewardCount    atomic.Int32
}

func (t *testOracle) ProcessBuilderCommitmentForBlockNumber(
	opts *bind.TransactOpts,
	commitmentIdx [32]byte,
	blockNum *big.Int,
	builder string,
	isSlash bool,
	residualDecay *big.Int,
) (*types.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if isSlash {
		t.slashCount.Add(1)
	} else {
		t.rewardCount.Add(1)
	}

	t.commitmentIdxs = append(t.commitmentIdxs, commitmentIdx)
	return types.MustSignNewTx(
		t.key,
		types.NewLondonSigner(big.NewInt(1)),
		&types.DynamicFeeTx{
			Nonce:     opts.Nonce.Uint64(),
			GasTipCap: opts.GasTipCap,
			GasFeeCap: opts.GasFeeCap,
		},
	), nil
}

func (t *testOracle) UnlockFunds(opts *bind.TransactOpts, bidIDs [][32]byte) (*types.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.bidIDs = append(t.bidIDs, bidIDs...)
	return types.MustSignNewTx(
		t.key,
		types.NewLondonSigner(big.NewInt(1)),
		&types.DynamicFeeTx{
			Nonce:     opts.Nonce.Uint64(),
			GasTipCap: opts.GasTipCap,
			GasFeeCap: opts.GasFeeCap,
		},
	), nil
}

func (t *testOracle) commitmentIdxsCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.commitmentIdxs)
}

func (t *testOracle) bidIDsCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.bidIDs)
}

type testTransactor struct {
	currentNonce       atomic.Uint64
	currentBlockNumber atomic.Uint64
}

func (t *testTransactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return t.currentNonce.Load() + 1, nil
}

func (t *testTransactor) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}

func (t *testTransactor) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}

func (t *testTransactor) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return t.currentNonce.Load(), nil
}

func (t *testTransactor) BlockNumber(ctx context.Context) (uint64, error) {
	return t.currentBlockNumber.Load(), nil
}

func waitForCount(dur time.Duration, expected int, f func() int) error {
	start := time.Now()
	for {
		if f() == expected {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
		if time.Since(start) > dur {
			return fmt.Errorf("expected count %d, got %d", expected, f())
		}
	}
}

func TestSettler(t *testing.T) {
	t.Parallel()

	ks, err := keysigner.NewPrivateKeySigner(path.Join(t.TempDir(), "key"))
	if err != nil {
		t.Fatal(err)
	}
	ownerAddr := common.HexToAddress("0xabcd")

	key, err := ks.GetPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	orcl := &testOracle{key: key}
	reg := &testRegister{
		settlementChan: make(chan settler.Settlement),
		returnsChan:    make(chan settler.Return),
	}
	transactor := &testTransactor{}

	s := settler.NewSettler(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		ks,
		big.NewInt(1000),
		ownerAddr,
		orcl,
		reg,
		transactor,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := s.Start(ctx)

	// Test that the settler is able to process a settlement
	for i := 0; i < 10; i++ {
		var sType settler.SettlementType
		if i%2 == 0 {
			sType = settler.SettlementTypeReward
		} else {
			sType = settler.SettlementTypeSlash
		}
		reg.settlementChan <- settler.Settlement{
			CommitmentIdx: big.NewInt(int64(i + 1)).Bytes(),
			TxHash:        "0x1234",
			BlockNum:      100,
			Builder:       "0x1234",
			Amount:        1000,
			BidID:         common.HexToHash(fmt.Sprintf("0x%02d", i)).Bytes(),
			Type:          sType,
		}

		if err := waitForCount(5*time.Second, i+1, orcl.commitmentIdxsCount); err != nil {
			t.Fatal(err)
		}

		if err := waitForCount(5*time.Second, i+1, reg.settlementsInitiatedCount); err != nil {
			t.Fatal(err)
		}

		reg.currentNonce.Add(1)
	}

	if reg.settlementsCompleted.Load() != 0 {
		t.Fatalf("expected 0 settlements completed, got %d", reg.settlementsCompleted.Load())
	}

	if orcl.slashCount.Load() != 5 {
		t.Fatalf("expected 5 slashes, got %d", orcl.slashCount.Load())
	}

	if orcl.rewardCount.Load() != 5 {
		t.Fatalf("expected 5 rewards, got %d", orcl.rewardCount.Load())
	}

	transactor.currentNonce.Store(10)
	transactor.currentBlockNumber.Store(1)

	if err := waitForCount(5*time.Second, 10, func() int {
		return int(reg.settlementsCompleted.Load())
	}); err != nil {
		t.Fatal(err)
	}

	returns := settler.Return{}

	for i := 0; i < 10; i++ {
		returns.BidIDs = append(returns.BidIDs, common.HexToHash(fmt.Sprintf("0x%02d", i)))
	}
	reg.returnsChan <- returns

	if err := waitForCount(5*time.Second, 20, reg.settlementsInitiatedCount); err != nil {
		t.Fatal(err)
	}

	if err := waitForCount(5*time.Second, 10, orcl.bidIDsCount); err != nil {
		t.Fatal(err)
	}

	transactor.currentNonce.Store(11)
	transactor.currentBlockNumber.Store(2)

	if err := waitForCount(5*time.Second, 11, func() int {
		return int(reg.settlementsCompleted.Load())
	}); err != nil {
		t.Fatal(err)
	}

	reg.pendingTxns.Store(129)

	reg.settlementChan <- settler.Settlement{
		CommitmentIdx: big.NewInt(11).Bytes(),
		TxHash:        "0x1234",
		BlockNum:      100,
		Builder:       "0x1234",
		Amount:        1000,
		Type:          settler.SettlementTypeReward,
	}

	time.Sleep(2 * time.Second)
	if reg.settlementsInitiatedCount() != 20 {
		t.Fatalf("expected 20 settlements initiated, got %d", reg.settlementsInitiatedCount())
	}

	cancel()
	<-done
}
