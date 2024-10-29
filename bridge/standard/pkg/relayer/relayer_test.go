package relayer_test

import (
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/primev/mev-commit/bridge/standard/pkg/relayer"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	settlementgateway "github.com/primev/mev-commit/contracts-abi/clients/SettlementGateway"
	"github.com/primev/mev-commit/x/util"
)

type Transfer struct {
	Recipient   common.Address
	Amount      *big.Int
	TransferIdx *big.Int
}

type testL1Gateway struct {
	initiated chan *l1gateway.L1gatewayTransferInitiated
	err       chan error
	mu        sync.Mutex // mu guards finalized.
	finalized []Transfer
}

func (t *testL1Gateway) Subscribe(ctx context.Context) (<-chan *l1gateway.L1gatewayTransferInitiated, <-chan error) {
	return t.initiated, t.err
}

func (t *testL1Gateway) FinalizeTransfer(ctx context.Context, recipient common.Address, amount *big.Int, transferIdx *big.Int) error {
	t.mu.Lock()
	t.finalized = append(t.finalized, Transfer{
		Recipient:   recipient,
		Amount:      amount,
		TransferIdx: transferIdx,
	})
	t.mu.Unlock()
	return nil
}

type testSettlementGateway struct {
	initiated chan *settlementgateway.SettlementgatewayTransferInitiated
	err       chan error
	mu        sync.Mutex // mu guards finalized.
	finalized []Transfer
}

func (t *testSettlementGateway) Subscribe(ctx context.Context) (<-chan *settlementgateway.SettlementgatewayTransferInitiated, <-chan error) {
	return t.initiated, t.err
}

func (t *testSettlementGateway) FinalizeTransfer(ctx context.Context, recipient common.Address, amount *big.Int, transferIdx *big.Int) error {
	t.mu.Lock()
	t.finalized = append(t.finalized, Transfer{
		Recipient:   recipient,
		Amount:      amount,
		TransferIdx: transferIdx,
	})
	t.mu.Unlock()
	return nil
}

func TestRelayer(t *testing.T) {
	l1Gateway := &testL1Gateway{
		initiated: make(chan *l1gateway.L1gatewayTransferInitiated),
		err:       make(chan error),
	}
	settlementGateway := &testSettlementGateway{
		initiated: make(chan *settlementgateway.SettlementgatewayTransferInitiated),
		err:       make(chan error),
	}

	relayer := relayer.NewRelayer(util.NewTestLogger(os.Stdout), l1Gateway, settlementGateway)

	ctx, cancel := context.WithCancel(context.Background())
	done := relayer.Start(ctx)

	expTransfers := []Transfer{
		{
			Recipient:   common.HexToAddress("0x1234"),
			Amount:      big.NewInt(100),
			TransferIdx: big.NewInt(1),
		},
		{
			Recipient:   common.HexToAddress("0x5678"),
			Amount:      big.NewInt(200),
			TransferIdx: big.NewInt(2),
		},
		{
			Recipient:   common.HexToAddress("0x9abc"),
			Amount:      big.NewInt(300),
			TransferIdx: big.NewInt(3),
		},
		{
			Recipient:   common.HexToAddress("0xdef0"),
			Amount:      big.NewInt(400),
			TransferIdx: big.NewInt(4),
		},
	}

	for _, transfer := range expTransfers {
		l1Gateway.initiated <- &l1gateway.L1gatewayTransferInitiated{
			Recipient:   transfer.Recipient,
			Amount:      transfer.Amount,
			TransferIdx: transfer.TransferIdx,
		}
	}

	start := time.Now()
	for {
		if time.Since(start) > 5*time.Second {
			t.Fatal("timeout waiting for relayer to finish")
		}
		settlementGateway.mu.Lock()
		if len(settlementGateway.finalized) == len(expTransfers) {
			settlementGateway.mu.Unlock()
			break
		}
		settlementGateway.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	settlementGateway.mu.Lock()
	if s := cmp.Diff(expTransfers, settlementGateway.finalized, cmp.AllowUnexported(big.Int{})); s != "" {
		settlementGateway.mu.Unlock()
		t.Fatalf("unexpected finalized transfers (-want +got):\n%s", s)
	}
	settlementGateway.mu.Unlock()

	for _, transfer := range expTransfers {
		settlementGateway.initiated <- &settlementgateway.SettlementgatewayTransferInitiated{
			Recipient:   transfer.Recipient,
			Amount:      transfer.Amount,
			TransferIdx: transfer.TransferIdx,
		}
	}

	start = time.Now()
	for {
		if time.Since(start) > 5*time.Second {
			t.Fatal("timeout waiting for relayer to finish")
		}
		l1Gateway.mu.Lock()
		if len(l1Gateway.finalized) == len(expTransfers) {
			l1Gateway.mu.Unlock()
			break
		}
		l1Gateway.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}

	l1Gateway.mu.Lock()
	if s := cmp.Diff(expTransfers, l1Gateway.finalized, cmp.AllowUnexported(big.Int{})); s != "" {
		l1Gateway.mu.Unlock()
		t.Fatalf("unexpected finalized transfers (-want +got):\n%s", s)
	}
	l1Gateway.mu.Unlock()

	cancel()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for relayer to finish")
	}
}
