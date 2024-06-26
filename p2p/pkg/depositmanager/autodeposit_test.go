package depositmanager_test

import (
	"context"
	"io"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/depositmanager"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/util"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type MockBidderRegistryContract struct {
	MoveDepositToWindowFunc func(opts *bind.TransactOpts, fromWindow *big.Int, toWindow *big.Int) (*types.Transaction, error)
}

func (m *MockBidderRegistryContract) MoveDepositToWindow(opts *bind.TransactOpts, fromWindow *big.Int, toWindow *big.Int) (*types.Transaction, error) {
	return m.MoveDepositToWindowFunc(opts, fromWindow, toWindow)
}

func (m *MockBidderRegistryContract) WithdrawFromSpecificWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return m.WithdrawFromSpecificWindows(opts, windows)
}
func TestAutoDepositTracker(t *testing.T) {
	t.Parallel()

	// Setup ABIs
	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	logger := util.NewTestLogger(io.Discard)
	evtMgr := events.NewListener(logger, &btABI, &brABI)

	mockContract := &MockBidderRegistryContract{
		MoveDepositToWindowFunc: func(opts *bind.TransactOpts, fromWindow *big.Int, toWindow *big.Int) (*types.Transaction, error) {
			return types.NewTx(&types.LegacyTx{}), nil
		},
	}

	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{}, nil
	}

	adt := depositmanager.NewAutoDepositTracker(evtMgr, mockContract, optsGetter, logger)

	ads := []*bidderapiv1.AutoDeposit{
		{WindowNumber: &wrapperspb.UInt64Value{Value: 1}, Amount: "100"},
		{WindowNumber: &wrapperspb.UInt64Value{Value: 2}, Amount: "100"},
		{WindowNumber: &wrapperspb.UInt64Value{Value: 3}, Amount: "100"},
	}
	ctx := context.Background()
	err = adt.DoAutoMoveToAnotherWindow(ctx, ads)
	if err != nil {
		t.Fatal(err)
	}

	br := &bidderregistry.BidderregistryBidderRegistered{
		Bidder:          common.HexToAddress("0x123"),
		DepositedAmount: big.NewInt(100),
		WindowNumber:    big.NewInt(1),
	}

	err = publishBidderRegistered(evtMgr, &brABI, br)
	if err != nil {
		t.Fatal(err)
	}

	publishNewWindow(evtMgr, &btABI, big.NewInt(int64(ads[0].WindowNumber.Value+1)))
	publishNewWindow(evtMgr, &btABI, big.NewInt(int64(ads[0].WindowNumber.Value+2)))

	deposits, status := adt.GetStatus()

	if !status {
		t.Fatalf("expected status to be true, got %v", status)
	}

	// need to wait for the goroutine to process the events
	time.Sleep(100 * time.Millisecond)
	for _, ad := range ads {
		if !deposits[ad.WindowNumber.Value+2] {
			t.Fatalf("expected deposit for window %d to be true, got %v", ad.WindowNumber.Value, deposits[ad.WindowNumber.Value+2])
		}
	}
	adt.Stop()

	// need to wait for the goroutine to stop
	time.Sleep(100 * time.Millisecond)
	_, status = adt.GetStatus()
	if status {
		t.Fatalf("expected status to be false, got %v", status)
	}
}
