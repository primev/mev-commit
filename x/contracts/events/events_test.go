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
	bidderregistry "github.com/primevprotocol/mev-commit/contracts-abi/clients/BidderRegistry"
	"github.com/primevprotocol/mev-commit/x/util"
)

func TestEventHandler(t *testing.T) {
	t.Parallel()

	b := bidderregistry.BidderregistryBidderRegistered{
		Bidder:          common.HexToAddress("0xabcd"),
		DepositedAmount: big.NewInt(1000),
	}

	errC := make(chan error, 1)

	evtHdlr := NewEventHandler(
		"BidderRegistered",
		func(ev *bidderregistry.BidderregistryBidderRegistered) {
			if ev.Bidder.Hex() != b.Bidder.Hex() {
				errC <- fmt.Errorf("expected bidder %s, got %s", b.Bidder.Hex(), ev.Bidder.Hex())
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

	event := bidderABI.Events["BidderRegistered"]

	evtHdlr.setTopicAndContract(event.ID, &bidderABI)

	buf, err := event.Inputs.NonIndexed().Pack(
		b.DepositedAmount,
	)
	if err != nil {
		t.Fatal(err)
	}

	bidder := common.HexToHash(b.Bidder.Hex())

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			bidder,   // The next topics are the indexed event parameters
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

	bidders := []bidderregistry.BidderregistryBidderRegistered{
		{
			Bidder:          common.HexToAddress("0xabcd"),
			DepositedAmount: big.NewInt(1000),
		},
		{
			Bidder:          common.HexToAddress("0xcdef"),
			DepositedAmount: big.NewInt(2000),
		},
	}

	logger := util.NewTestLogger(io.Discard)

	count := 0
	handlerTriggered := make(chan int, 1)
	errC := make(chan error, 1)

	evtHdlr := NewEventHandler(
		"BidderRegistered",
		func(ev *bidderregistry.BidderregistryBidderRegistered) {
			if count >= len(bidders) {
				errC <- fmt.Errorf("unexpected event")
				return
			}
			if ev.Bidder.Hex() != bidders[count].Bidder.Hex() {
				errC <- fmt.Errorf("expected bidder %s, got %s", bidders[count].Bidder.Hex(), ev.Bidder.Hex())
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

	data1, err := bidderABI.Events["BidderRegistered"].Inputs.NonIndexed().Pack(
		bidders[0].DepositedAmount,
	)
	if err != nil {
		t.Fatal(err)
	}

	data2, err := bidderABI.Events["BidderRegistered"].Inputs.NonIndexed().Pack(
		bidders[1].DepositedAmount,
	)
	if err != nil {
		t.Fatal(err)
	}

	logs := []types.Log{
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderRegistered"].ID,
				common.HexToHash(bidders[0].Bidder.Hex()),
			},
			Data:        data1,
			BlockNumber: 1,
		},
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderRegistered"].ID,
				common.HexToHash(bidders[1].Bidder.Hex()),
			},
			Data:        data2,
			BlockNumber: 2,
		},
		{
			Topics: []common.Hash{
				bidderABI.Events["BidderRegistered"].ID,
				common.HexToHash("test"),
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
