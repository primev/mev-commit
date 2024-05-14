package depositmanager_test

import (
	"context"
	"io"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primevprotocol/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-commit/p2p/pkg/depositmanager"
	"github.com/primevprotocol/mev-commit/p2p/pkg/store"
	"github.com/primevprotocol/mev-commit/x/contracts/events"
	"github.com/primevprotocol/mev-commit/x/util"
)

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

	st := store.NewStore()
	bt := &testBlockTracker{value: big.NewInt(10)}

	ctx, cancel := context.WithCancel(context.Background())

	dm := depositmanager.NewDepositManager(bt, st, evtMgr, logger)
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
		if st.Len() == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	cancel()
	<-done
}

type testBlockTracker struct {
	value *big.Int
}

func (tbt *testBlockTracker) GetBlocksPerWindow() (*big.Int, error) {
	return tbt.value, nil
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
