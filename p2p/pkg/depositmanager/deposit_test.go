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
	GetDepositFunc              func(opts *bind.CallOpts, bidder common.Address, provider common.Address) (*big.Int, error)
	WithdrawalRequestExistsFunc func(opts *bind.CallOpts, bidder common.Address, provider common.Address) (bool, error)
}

func (m *MockBidderRegistryContract) GetDeposit(
	opts *bind.CallOpts,
	bidder common.Address,
	provider common.Address,
) (*big.Int, error) {
	return m.GetDepositFunc(opts, bidder, provider)
}

func (m *MockBidderRegistryContract) WithdrawalRequestExists(
	opts *bind.CallOpts,
	bidder common.Address,
	provider common.Address,
) (bool, error) {
	return m.WithdrawalRequestExistsFunc(opts, bidder, provider)
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
			provider common.Address,
		) (*big.Int, error) {
			return big.NewInt(0), nil
		},
		WithdrawalRequestExistsFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			provider common.Address,
		) (bool, error) {
			return false, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	dm := depositmanager.NewDepositManager(st, evtMgr, bidderRegistry, logger)
	done := dm.Start(ctx)

	// no deposit
	refund, err := dm.CheckAndDeductDeposit(
		context.Background(),
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
		"10",
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if refund != nil {
		t.Fatal("expected nil refund")
	}

	br := &bidderregistry.BidderregistryBidderDeposited{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        common.HexToAddress("0x456"),
		DepositedAmount: big.NewInt(100),
	}

	err = publishBidderDeposited(evtMgr, &brABI, br)
	if err != nil {
		t.Fatal(err)
	}

	for {
		if val, err := st.GetBalance(
			common.HexToAddress("0x123"),
			common.HexToAddress("0x456"),
		); err == nil && val != nil && val.Cmp(big.NewInt(100)) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// deduct deposit
	refund, err = dm.CheckAndDeductDeposit(
		context.Background(),
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
		"100",
	)
	if err != nil {
		t.Fatal(err)
	}

	// not enough deposit
	_, err = dm.CheckAndDeductDeposit(
		context.Background(),
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
		"10",
	)
	if err == nil || !strings.Contains(err.Error(), "insufficient balance") {
		t.Fatal("expected error for insufficient balance")
	}

	err = refund()
	if err != nil {
		t.Fatal(err)
	}

	// deduct deposit after refund
	_, err = dm.CheckAndDeductDeposit(
		context.Background(),
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
		"10",
	)
	if err != nil {
		t.Fatal(err)
	}

	balance, err := st.GetBalance(
		common.HexToAddress("0x123"),
		common.HexToAddress("0x456"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if balance == nil || balance.Cmp(big.NewInt(90)) != 0 {
		t.Fatal("expected balance of 90")
	}

	err = publishBidderWithdrawalRequested(evtMgr, &brABI, &bidderregistry.BidderregistryWithdrawalRequested{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        common.HexToAddress("0x456"),
		AvailableAmount: big.NewInt(10),
		EscrowedAmount:  big.NewInt(10),
		Timestamp:       big.NewInt(1000),
	})
	if err != nil {
		t.Fatal(err)
	}

	for {
		if val, err := st.GetBalance(
			common.HexToAddress("0x123"),
			common.HexToAddress("0x456"),
		); err == nil && val == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err = publishBidderWithdrawal(evtMgr, &brABI, &bidderregistry.BidderregistryBidderWithdrawal{
		Bidder:              common.HexToAddress("0x123"),
		Provider:            common.HexToAddress("0x456"),
		AmountWithdrawn:     big.NewInt(10),
		AmountStillEscrowed: big.NewInt(10),
	})
	if err != nil {
		t.Fatal(err)
	}

	for {
		count, err := st.BalanceEntries(common.HexToAddress("0x123"))
		if err != nil {
			t.Fatal(err)
		}
		if count == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	err = publishBidderDeposited(evtMgr, &brABI, &bidderregistry.BidderregistryBidderDeposited{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        common.HexToAddress("0x456"),
		DepositedAmount: big.NewInt(777),
	})
	if err != nil {
		t.Fatal(err)
	}

	for {
		if val, err := st.GetBalance(
			common.HexToAddress("0x123"),
			common.HexToAddress("0x456"),
		); err == nil && val != nil && val.Cmp(big.NewInt(777)) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	cancel()
	<-done
}

func TestStartWithBidderAlreadyDeposited(t *testing.T) {
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
			provider common.Address,
		) (*big.Int, error) {
			if opts.BlockNumber.Cmp(big.NewInt(15)) != 0 {
				t.Fatal("expected block number 15")
			}
			return big.NewInt(33), nil // Existing deposit
		},
		WithdrawalRequestExistsFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			provider common.Address,
		) (bool, error) {
			if opts.BlockNumber.Cmp(big.NewInt(15)) != 0 {
				t.Fatal("expected block number 15")
			}
			return false, nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	dm := depositmanager.NewDepositManager(st, evtMgr, bidderRegistry, logger)
	done := dm.Start(ctx)

	publishBidderDeposited(evtMgr, &brABI, &bidderregistry.BidderregistryBidderDeposited{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        common.HexToAddress("0x456"),
		DepositedAmount: big.NewInt(100),
		Raw: types.Log{
			BlockNumber: 16,
		},
	})

	for {
		if val, err := st.GetBalance(
			common.HexToAddress("0x123"),
			common.HexToAddress("0x456"),
		); err == nil && val != nil && val.Cmp(big.NewInt(133)) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	cancel()
	<-done
}

func publishBidderDeposited(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	br *bidderregistry.BidderregistryBidderDeposited,
) error {
	event := brABI.Events["BidderDeposited"]
	buf, err := event.Inputs.NonIndexed().Pack()
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
			common.HexToHash(br.Provider.Hex()),
			common.BigToHash(br.DepositedAmount),
		},
		Data:        buf,
		BlockNumber: br.Raw.BlockNumber,
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)

	return nil
}

func publishBidderWithdrawalRequested(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	br *bidderregistry.BidderregistryWithdrawalRequested,
) error {
	event := brABI.Events["WithdrawalRequested"]
	buf, err := event.Inputs.NonIndexed().Pack(br.AvailableAmount, br.EscrowedAmount)
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
			common.HexToHash(br.Provider.Hex()),
			common.BigToHash(br.Timestamp),
		},
		Data:        buf,
		BlockNumber: 1,
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
	buf, err := event.Inputs.NonIndexed().Pack(br.AmountWithdrawn)
	if err != nil {
		return err
	}

	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,
			common.HexToHash(br.Bidder.Hex()),
			common.HexToHash(br.Provider.Hex()),
			common.BigToHash(br.AmountWithdrawn),
		},
		Data: buf,
	}
	evtMgr.PublishLogEvent(context.Background(), testLog)

	return nil
}
