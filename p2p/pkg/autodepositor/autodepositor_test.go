package autodepositor_test

import (
	"context"
	"io"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/p2p/pkg/autodepositor"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/util"
)

type MockBidderRegistryContract struct {
	DepositForWindowsFunc   func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
	WithdrawFromWindowsFunc func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error)
}

func (m *MockBidderRegistryContract) DepositForWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return m.DepositForWindowsFunc(opts, windows)
}

func (m *MockBidderRegistryContract) WithdrawFromWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return m.WithdrawFromWindowsFunc(opts, windows)
}

func TestAutoDepositTracker_Start(t *testing.T) {
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

	addr := common.HexToAddress("0x1234")
	amount := big.NewInt(100)
	logger := util.NewTestLogger(os.Stdout)
	evtMgr := events.NewListener(logger, &btABI, &brABI)
	brContract := &MockBidderRegistryContract{
		DepositForWindowsFunc: func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
			for _, window := range windows {
				err = publishBidderRegistered(evtMgr, &brABI, &bidderregistry.BidderregistryBidderRegistered{
					Bidder:          addr,
					DepositedAmount: amount,
					WindowNumber:    window,
				})
				if err != nil {
					return nil, err
				}
			}
			return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
		},
		WithdrawFromWindowsFunc: func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
			for _, window := range windows {
				err = publishBidderWithdrawal(evtMgr, &brABI, &bidderregistry.BidderregistryBidderWithdrawal{
					Bidder: addr,
					Amount: amount,
					Window: window,
				})
				if err != nil {
					return nil, err
				}
			}
			return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
		},
	}
	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{}, nil
	}

	// Create AutoDepositTracker instance
	adt := autodepositor.New(evtMgr, brContract, optsGetter, logger)

	// Start AutoDepositTracker
	ctx := context.Background()
	startWindow := big.NewInt(2)
	err = adt.Start(ctx, startWindow, amount)
	if err != nil {
		t.Fatal(err)
	}

	assertStatus := func(t *testing.T, working bool, deposits []uint64) {
		t.Helper()

		for {
			depositsMap, status := adt.GetStatus()
			if status != working {
				t.Fatalf("expected status to be %v, got %v", working, status)
			}
			foundAll := true
			for _, deposit := range deposits {
				if !depositsMap[deposit] {
					foundAll = false
					break
				}
			}
			if foundAll && len(depositsMap) == len(deposits) {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	assertStatus(t, true, []uint64{2, 3})

	publishNewWindow(evtMgr, &btABI, big.NewInt(1))

	assertStatus(t, true, []uint64{2, 3, 4})

	publishNewWindow(evtMgr, &btABI, big.NewInt(2))

	assertStatus(t, true, []uint64{2, 3, 4, 5})

	publishNewWindow(evtMgr, &btABI, big.NewInt(3))

	assertStatus(t, true, []uint64{3, 4, 5, 6})

	// Stop AutoDepositTracker
	windowNumbers, err := adt.Stop()
	if err != nil {
		t.Fatal(err)
	}

	// Assert window numbers
	expectedWindowNumbers := []*big.Int{big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6)}
	if len(windowNumbers) != len(expectedWindowNumbers) {
		t.Fatalf("expected %d window numbers, got %d", len(expectedWindowNumbers), len(windowNumbers))
	}
	for i, wn := range windowNumbers {
		if wn.Cmp(expectedWindowNumbers[i]) != 0 {
			t.Errorf("expected window number %d to be %s, got %s", i, expectedWindowNumbers[i].String(), wn.String())
		}
	}

	assertStatus(t, false, []uint64{})
}

func TestAutoDepositTracker_Start_CancelContext(t *testing.T) {
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
	brContract := &MockBidderRegistryContract{
		DepositForWindowsFunc: func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
			return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
		},
	}
	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{}, nil
	}

	// Create AutoDepositTracker instance
	adt := autodepositor.New(evtMgr, brContract, optsGetter, logger)

	// Start AutoDepositTracker with a cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	startWindow := big.NewInt(1)
	amount := big.NewInt(100)
	cancel()
	err = adt.Start(ctx, startWindow, amount)
	if err != context.Canceled {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}

func TestAutoDepositTracker_Stop_NotRunning(t *testing.T) {
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
	brContract := &MockBidderRegistryContract{}
	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{}, nil
	}

	// Create AutoDepositTracker instance
	adt := autodepositor.New(evtMgr, brContract, optsGetter, logger)

	// Stop AutoDepositTracker when not running
	_, err = adt.Stop()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAutoDepositTracker_IsWorking(t *testing.T) {
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
	brContract := &MockBidderRegistryContract{
		DepositForWindowsFunc: func(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
			return types.NewTransaction(1, common.Address{}, nil, 0, nil, nil), nil
		},
	}
	optsGetter := func(ctx context.Context) (*bind.TransactOpts, error) {
		return &bind.TransactOpts{}, nil
	}

	// Create AutoDepositTracker instance
	adt := autodepositor.New(evtMgr, brContract, optsGetter, logger)

	// Assert initial IsWorking status
	if adt.IsWorking() {
		t.Fatal("expected IsWorking to be false, got true")
	}

	// Start AutoDepositTracker
	ctx := context.Background()
	startWindow := big.NewInt(1)
	amount := big.NewInt(100)
	err = adt.Start(ctx, startWindow, amount)
	if err != nil {
		t.Fatal(err)
	}

	// Assert IsWorking status after starting
	if !adt.IsWorking() {
		t.Fatal("expected IsWorking to be true, got false")
	}

	// Stop AutoDepositTracker
	_, err = adt.Stop()
	if err != nil {
		t.Fatal(err)
	}

	// Assert IsWorking status after stopping
	if adt.IsWorking() {
		t.Fatal("expected IsWorking to be false, got true")
	}
}

func publishNewWindow(
	evtMgr events.EventManager,
	btABI *abi.ABI,
	windowNumber *big.Int,
) {
	testLog := types.Log{
		Topics: []common.Hash{
			btABI.Events["NewWindow"].ID,
			common.BigToHash(windowNumber),
		},
		Data: []byte{},
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)
}

func publishBidderRegistered(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	br *bidderregistry.BidderregistryBidderRegistered,
) error {
	event := brABI.Events["BidderRegistered"]
	buf, err := event.Inputs.NonIndexed().Pack(
		br.DepositedAmount,
		br.WindowNumber,
	)
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
		},
		Data: buf,
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)

	return nil
}

func publishBidderWithdrawal(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	br *bidderregistry.BidderregistryBidderWithdrawal,
) error {
	event := brABI.Events["BidderWithdrawal"]
	buf, err := event.Inputs.NonIndexed().Pack(
		br.Window,
		br.Amount,
	)
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
		},
		Data: buf,
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)

	return nil
}
