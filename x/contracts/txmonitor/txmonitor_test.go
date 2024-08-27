package txmonitor_test

import (
	"context"
	"io"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/util"
)

func TestTxMonitor(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	txns := make([]*types.Transaction, 0, 10)
	for i := 1; i <= 10; i++ {
		txns = append(txns, types.MustSignNewTx(
			key,
			types.NewLondonSigner(big.NewInt(1)),
			&types.DynamicFeeTx{
				ChainID:   big.NewInt(1),
				Nonce:     uint64(i),
				GasFeeCap: big.NewInt(1),
				GasTipCap: big.NewInt(1),
				To:        &common.Address{},
			},
		))
	}

	results := make(map[common.Hash]txmonitor.Result)
	for _, tx := range txns {
		results[tx.Hash()] = txmonitor.Result{Receipt: &types.Receipt{Status: 1}}
	}

	evm := &testEVM{
		blockNumC: make(chan uint64),
		nonceC:    make(chan uint64),
	}

	saver := &testSaver{status: make(map[common.Hash]string)}

	monitor := txmonitor.New(
		common.Address{},
		evm,
		&testEVMHelper{receipts: results},
		saver,
		util.NewTestLogger(io.Discard),
		10,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := monitor.Start(ctx)

	for _, tx := range txns {
		if allow := monitor.Allow(ctx, tx.Nonce()); !allow {
			t.Fatal("tx should be allowed")
		}
		monitor.Sent(ctx, tx)
	}

	allowCtx, cancelAllow := context.WithTimeout(ctx, 200*time.Millisecond)
	if allow := monitor.Allow(allowCtx, 11); allow {
		t.Fatal("tx should not be allowed")
	}
	cancelAllow()

	// tx with same nonce to simulate cancellation
	cancelledWait := monitor.WatchTx(common.HexToHash("0x1"), 2)

	evm.blockNumC <- 1
	evm.nonceC <- 5

	for {
		if saver.count() == 4 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// tx with same nonce should be cancelled
	res := <-cancelledWait
	if res.Err != txmonitor.ErrTxnCancelled {
		t.Fatal("tx should be cancelled")
	}

	duplicateListener := monitor.WatchTx(txns[9].Hash(), 10)
	closedListener := monitor.WatchTx(common.HexToHash("0x12"), 11)

	evm.blockNumC <- 2
	evm.nonceC <- 11

	for {
		if saver.count() == 10 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// duplicate listener should be notified
	res = <-duplicateListener
	if res.Err != nil {
		t.Fatal("tx should not have error")
	}

	cancel()
	<-done

	// closed listener should be notified
	res = <-closedListener
	if res.Err != txmonitor.ErrMonitorClosed {
		t.Fatal("error should be monitor closed")
	}

	for _, status := range saver.status {
		if status != "success" {
			t.Fatal("tx should be successful")
		}
	}
}

type testEVM struct {
	blockNumC chan uint64
	nonceC    chan uint64
}

func (t *testEVM) BlockNumber(ctx context.Context) (uint64, error) {
	return <-t.blockNumC, nil
}

func (t *testEVM) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	return <-t.nonceC, nil
}

type testEVMHelper struct {
	receipts map[common.Hash]txmonitor.Result
}

func (t *testEVMHelper) BatchReceipts(ctx context.Context, txns []common.Hash) ([]txmonitor.Result, error) {
	results := make([]txmonitor.Result, 0, len(txns))
	for _, tx := range txns {
		if _, ok := t.receipts[tx]; !ok {
			results = append(results, txmonitor.Result{Err: ethereum.NotFound})
			continue
		}
		results = append(results, t.receipts[tx])
	}
	return results, nil
}

func (t *testEVMHelper) RevertReason(ctx context.Context, r *types.Receipt, from common.Address) (string, error) {
	return "dummy error", nil
}

type testSaver struct {
	mu     sync.Mutex
	status map[common.Hash]string
}

func (t *testSaver) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.status)
}

func (t *testSaver) Save(ctx context.Context, txHash common.Hash, nonce uint64) error {
	return nil
}

func (t *testSaver) Update(ctx context.Context, txHash common.Hash, status string) error {
	t.mu.Lock()
	t.status[txHash] = status
	t.mu.Unlock()
	return nil
}

func (t *testSaver) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	return nil, nil
}
