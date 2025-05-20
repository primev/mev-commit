package transfer_test

import (
	"context"
	"log/slog"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	bridgetransfer "github.com/primev/mev-commit/bridge/standard/pkg/transfer"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/transfer"
)

type MockAccountSyncer struct {
	trigger chan struct{}
}

func (m *MockAccountSyncer) Subscribe(ctx context.Context, threshold *big.Int) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		<-m.trigger
		close(ch)
	}()
	return ch
}

type MockTransfer struct {
	mtx    sync.Mutex
	called int
	amount *big.Int
}

func (m *MockTransfer) Do(ctx context.Context) <-chan bridgetransfer.TransferStatus {
	ch := make(chan bridgetransfer.TransferStatus, 1)
	ch <- bridgetransfer.TransferStatus{
		Message: "Transfer Done",
		Error:   nil,
	}
	close(ch)
	m.mtx.Lock()
	m.called++
	m.mtx.Unlock()
	return ch
}

func (m *MockTransfer) getAmount() *big.Int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.amount
}

func (m *MockTransfer) setAmount(amount *big.Int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.amount = amount
}

func (m *MockTransfer) called() int {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.called
}

type keySigner struct {
	keysigner.KeySigner
}

func (k *keySigner) GetAddress() common.Address {
	return common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
}

func TestBridger(t *testing.T) {
	t.Parallel()

	syncer := &MockAccountSyncer{
		trigger: make(chan struct{}),
	}

	txfer := &MockTransfer{}

	done := transfer.SetTransferFunc(func(
		amount *big.Int,
		destAddress common.Address,
		signer keysigner.KeySigner,
		settlementRPCUrl string,
		l1RPCUrl string,
		l1ContractAddr common.Address,
		settlementContractAddr common.Address,
	) (bridgetransfer.Transfer, error) {
		txfer.setAmount(amount)
		return txfer, nil
	},
	)
	t.Cleanup(done)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	bridger := transfer.NewBridger(
		logger,
		syncer,
		transfer.BridgeConfig{
			Signer:                 &keySigner{},
			L1ContractAddr:         common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
			SettlementContractAddr: common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			L1RPCUrl:               "http://localhost:8545",
			SettlementRPCUrl:       "http://localhost:8546",
		},
		big.NewInt(1000000000000000000),
		big.NewInt(1000000000000000000),
	)
	ctx, cancel := context.WithCancel(context.Background())
	closed := bridger.Start(ctx)

	// Simulate the syncer triggering the event
	syncer.trigger <- struct{}{}

	start := time.Now()
	for {
		if time.Since(start) > 5*time.Second {
			t.Fatal("Timeout waiting for transfer to be called")
		}
		if txfer.called() == 1 {
			break
		}
	}

	if txfer.getAmount().Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Fatalf("Expected amount to be 1 ETH, got %s", txfer.amount.String())
	}

	syncer.trigger <- struct{}{}

	start = time.Now()
	for {
		if time.Since(start) > 5*time.Second {
			t.Fatal("Timeout waiting for transfer to be called again")
		}
		if txfer.called() == 2 {
			break
		}
	}

	if txfer.getAmount().Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Fatalf("Expected amount to be 1 ETH, got %s", txfer.amount.String())
	}

	cancel()
	<-closed

}
