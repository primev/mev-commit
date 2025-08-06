package events

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	"github.com/primev/mev-commit/x/util"
)

func TestEventHandler(t *testing.T) {
	t.Parallel()

	b := bidderregistry.BidderregistryBidderDeposited{
		Bidder:          common.HexToAddress("0xabcd"),
		Provider:        common.HexToAddress("0x1234"),
		DepositedAmount: big.NewInt(1000),
	}

	errC := make(chan error, 1)

	evtHdlr := NewEventHandler(
		"BidderDeposited",
		func(ev *bidderregistry.BidderregistryBidderDeposited) {
			if ev.Bidder.Hex() != b.Bidder.Hex() {
				errC <- fmt.Errorf("expected bidder %s, got %s", b.Bidder.Hex(), ev.Bidder.Hex())
				return
			}
			if ev.Provider.Hex() != b.Provider.Hex() {
				errC <- fmt.Errorf("expected provider %s, got %s", b.Provider.Hex(), ev.Provider.Hex())
				return
			}
			if ev.DepositedAmount.Cmp(b.DepositedAmount) != 0 {
				errC <- fmt.Errorf("expected prepaid amount %d, got %d", b.DepositedAmount, ev.DepositedAmount)
				return
			}
			close(errC)
		},
	)

	bidderABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	event := bidderABI.Events["BidderDeposited"]

	evtHdlr.setTopicAndContract(event.ID, &bidderABI)

	buf, err := event.Inputs.NonIndexed().Pack()
	if err != nil {
		t.Fatal(err)
	}

	bidder := common.HexToHash(b.Bidder.Hex())
	provider := common.HexToHash(b.Provider.Hex())
	depositedAmount := common.BigToHash(b.DepositedAmount)

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			bidder,   // The next topics are the indexed event parameters
			provider,
			depositedAmount,
		},
		Data: buf,
	}

	if err := evtHdlr.handle(testLog); err != nil {
		t.Fatal(err)
	}

	select {
	case err, more := <-errC:
		if !more {
			break
		}
		t.Fatal(err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for handler to be triggered")
	}
}

func TestEventManager(t *testing.T) {
	t.Parallel()

	bidders := []bidderregistry.BidderregistryBidderDeposited{
		{
			Bidder:          common.HexToAddress("0xabcd"),
			Provider:        common.HexToAddress("0x1234"),
			DepositedAmount: big.NewInt(1000),
		},
		{
			Bidder:          common.HexToAddress("0xcdef"),
			Provider:        common.HexToAddress("0x5678"),
			DepositedAmount: big.NewInt(2000),
		},
	}

	logger := util.NewTestLogger(io.Discard)

	count := 0
	handlerTriggered := make(chan int, 1)
	errC := make(chan error, 1)

	evtHdlr := NewEventHandler(
		"BidderDeposited",
		func(ev *bidderregistry.BidderregistryBidderDeposited) {
			if count >= len(bidders) {
				errC <- fmt.Errorf("unexpected event")
				return
			}
			if ev.Bidder.Hex() != bidders[count].Bidder.Hex() {
				errC <- fmt.Errorf("expected bidder %s, got %s", bidders[count].Bidder.Hex(), ev.Bidder.Hex())
				return
			}
			if ev.Provider.Hex() != bidders[count].Provider.Hex() {
				errC <- fmt.Errorf("expected provider %s, got %s", bidders[count].Provider.Hex(), ev.Provider.Hex())
				return
			}
			if ev.DepositedAmount.Cmp(bidders[count].DepositedAmount) != 0 {
				errC <- fmt.Errorf("expected prepaid amount %d, got %d", bidders[count].DepositedAmount, ev.DepositedAmount)
				return
			}
			count++
			handlerTriggered <- count
		},
	)

	bidderABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	data1, err := bidderABI.Events["BidderDeposited"].Inputs.NonIndexed().Pack()
	if err != nil {
		t.Fatal(err)
	}

	data2, err := bidderABI.Events["BidderDeposited"].Inputs.NonIndexed().Pack()
	if err != nil {
		t.Fatal(err)
	}

	logs := []types.Log{
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderDeposited"].ID,
				common.HexToHash(bidders[0].Bidder.Hex()),
				common.HexToHash(bidders[0].Provider.Hex()),
				common.BigToHash(bidders[0].DepositedAmount),
			},
			Data:        data1,
			BlockNumber: 1,
		},
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderDeposited"].ID,
				common.HexToHash(bidders[1].Bidder.Hex()),
				common.HexToHash(bidders[1].Provider.Hex()),
				common.BigToHash(bidders[1].DepositedAmount),
			},
			Data:        data2,
			BlockNumber: 2,
		},
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderDeposited"].ID,
				common.HexToHash("test"),
				common.BigToHash(big.NewInt(3000)),
			},
			Data:        []byte("test"),
			BlockNumber: 3,
		},
	}

	evtMgr := NewListener(
		logger,
		&bidderABI,
	)

	sub, err := evtMgr.Subscribe(evtHdlr)
	if err != nil {
		t.Fatal(err)
	}

	evtMgr.PublishLogEvent(context.Background(), logs[0])
	select {
	case c := <-handlerTriggered:
		if c != 1 {
			t.Fatalf("expected 1, got %d", c)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for handler to be triggered")
	}

	evtMgr.PublishLogEvent(context.Background(), logs[1])
	select {
	case c := <-handlerTriggered:
		if c != 2 {
			t.Fatalf("expected 2, got %d", c)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for handler to be triggered")
	}

	evtMgr.PublishLogEvent(context.Background(), logs[2])
	select {
	case err := <-sub.Err():
		if err == nil {
			t.Fatal("expected error")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for error")
	}

	sub.Unsubscribe()
	evtMgr.PublishLogEvent(context.Background(), logs[0])
	select {
	case <-handlerTriggered:
		t.Fatal("unexpected handler trigger")
	case err, more := <-sub.Err():
		if !more {
			break
		}
		t.Fatal(err)
	case <-time.After(5 * time.Second):
	}
}
