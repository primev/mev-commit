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
	"github.com/primev/mev-commit/p2p/pkg/depositmanager"
	depositstore "github.com/primev/mev-commit/p2p/pkg/depositmanager/store"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/util"
)

type MockBidderRegistryContract struct {
	GetDepositFunc func(opts *bind.CallOpts, bidder common.Address, window *big.Int) (*big.Int, error)
}

func (m *MockBidderRegistryContract) GetDeposit(
	opts *bind.CallOpts,
	bidder common.Address,
	window *big.Int,
) (*big.Int, error) {
	return m.GetDepositFunc(opts, bidder, window)
}

func TestDepositManager(t *testing.T) {
	t.Parallel()

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

	st := depositstore.New(inmemstorage.New())
	bidderRegistry := &MockBidderRegistryContract{
		GetDepositFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			window *big.Int,
		) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	dm := depositmanager.NewDepositManager(10, st, evtMgr, bidderRegistry, logger)
	done := dm.Start(ctx)

	// no deposit
	_, err = dm.CheckAndDeductDeposit(
		context.Background(),
		common.HexToAddress("0x123"),
		"10",
		1,
	)
	if err == nil {
		t.Fatal("expected error")
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

	for {
		if val, err := st.GetBalance(
			common.HexToAddress("0x123"),
			big.NewInt(1),
		); err == nil && val != nil && val.Cmp(big.NewInt(10)) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	for i := int64(1); i <= 10; i++ {
		// deduct deposit
		refund, err := dm.CheckAndDeductDeposit(
			context.Background(),
			common.HexToAddress("0x123"),
			"10",
			i,
		)
		if err != nil {
			t.Fatal(err)
		}

		// not enough deposit
		_, err = dm.CheckAndDeductDeposit(
			context.Background(),
			common.HexToAddress("0x123"),
			"10",
			i,
		)
		if err == nil {
			t.Fatal("expected error")
		}

		err = refund()
		if err != nil {
			t.Fatal(err)
		}

		// deduct deposit after refund
		_, err = dm.CheckAndDeductDeposit(
			context.Background(),
			common.HexToAddress("0x123"),
			"10",
			i,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	publishNewWindow(evtMgr, &btABI, big.NewInt(12))
	for {
		count, err := st.BalanceEntries(big.NewInt(1))
		if err != nil {
			t.Fatal(err)
		}
		if count == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	cancel()
	<-done
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
	buf, err := event.Inputs.NonIndexed().Pack()
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
			common.BigToHash(br.DepositedAmount),
			common.BigToHash(br.WindowNumber),
		},
		Data: buf,
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)

	return nil
}
