package depositmanager_test

import (
	"bytes"
	"context"
	"io"
	"math/big"
	"strings"
	"sync"
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
	GetDepositConsideringWithdrawalRequestFunc func(opts *bind.CallOpts, bidder common.Address, provider common.Address) (*big.Int, error)
}

func (m *MockBidderRegistryContract) GetDepositConsideringWithdrawalRequest(
	opts *bind.CallOpts,
	bidder common.Address,
	provider common.Address,
) (*big.Int, error) {
	return m.GetDepositConsideringWithdrawalRequestFunc(opts, bidder, provider)
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
		GetDepositConsideringWithdrawalRequestFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			provider common.Address,
		) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	providerAddress := common.HexToAddress("0x456")

	dm := depositmanager.NewDepositManager(st, evtMgr, bidderRegistry, providerAddress, logger)
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
		Bidder:             common.HexToAddress("0x123"),
		Provider:           common.HexToAddress("0x456"),
		DepositedAmount:    big.NewInt(100),
		NewAvailableAmount: big.NewInt(100),
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
		Bidder:             common.HexToAddress("0x123"),
		Provider:           common.HexToAddress("0x456"),
		DepositedAmount:    big.NewInt(777),
		NewAvailableAmount: big.NewInt(777),
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
		GetDepositConsideringWithdrawalRequestFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			provider common.Address,
		) (*big.Int, error) {
			if opts.BlockNumber.Cmp(big.NewInt(15)) != 0 {
				t.Fatal("expected block number 15")
			}
			return big.NewInt(33), nil // Existing deposit
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	providerAddress := common.HexToAddress("0x456")

	dm := depositmanager.NewDepositManager(st, evtMgr, bidderRegistry, providerAddress, logger)
	done := dm.Start(ctx)

	err = publishBidderDeposited(evtMgr, &brABI, &bidderregistry.BidderregistryBidderDeposited{
		Bidder:             common.HexToAddress("0x123"),
		Provider:           common.HexToAddress("0x456"),
		DepositedAmount:    big.NewInt(100),
		NewAvailableAmount: big.NewInt(133),
		Raw: types.Log{
			BlockNumber: 16,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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

func TestOtherProvidersEventsAreIgnored(t *testing.T) {
	t.Parallel()

	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	logBuf := &SafeBuffer{}
	logger := util.NewTestLogger(logBuf)
	evtMgr := events.NewListener(logger, &btABI, &brABI)

	st := depositstore.New(inmemstorage.New())
	bidderRegistry := &MockBidderRegistryContract{
		GetDepositConsideringWithdrawalRequestFunc: func(
			opts *bind.CallOpts,
			bidder common.Address,
			provider common.Address,
		) (*big.Int, error) {
			return big.NewInt(0), nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	providerAddress := common.HexToAddress("0x456")

	dm := depositmanager.NewDepositManager(st, evtMgr, bidderRegistry, providerAddress, logger)
	done := dm.Start(ctx)

	differentProvider := common.HexToAddress("0x789")

	err = publishBidderDeposited(evtMgr, &brABI, &bidderregistry.BidderregistryBidderDeposited{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        differentProvider,
		DepositedAmount: big.NewInt(100),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = publishBidderWithdrawalRequested(evtMgr, &brABI, &bidderregistry.BidderregistryWithdrawalRequested{
		Bidder:          common.HexToAddress("0x123"),
		Provider:        differentProvider,
		AvailableAmount: big.NewInt(100),
		EscrowedAmount:  big.NewInt(100),
		Timestamp:       big.NewInt(1000),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = publishBidderWithdrawal(evtMgr, &brABI, &bidderregistry.BidderregistryBidderWithdrawal{
		Bidder:              common.HexToAddress("0x123"),
		Provider:            differentProvider,
		AmountWithdrawn:     big.NewInt(100),
		AmountStillEscrowed: big.NewInt(100),
	})
	if err != nil {
		t.Fatal(err)
	}

	type seen struct {
		deposit           bool
		withdrawalRequest bool
		withdrawal        bool
	}
	haveSeen := seen{}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(logBuf.String(), "ignoring deposit event for different provider") {
			haveSeen.deposit = true
		}
		if strings.Contains(logBuf.String(), "ignoring withdrawal request event for different provider") {
			haveSeen.withdrawalRequest = true
		}
		if strings.Contains(logBuf.String(), "ignoring withdrawal event for different provider") {
			haveSeen.withdrawal = true
		}
		time.Sleep(1 * time.Second)
	}
	if !haveSeen.deposit || !haveSeen.withdrawalRequest || !haveSeen.withdrawal {
		t.Fatal("expected all events to be seen, but got ", haveSeen)
	}

	cancel()
	<-done
}

type SafeBuffer struct {
	mu  sync.RWMutex
	buf bytes.Buffer
}

func (b *SafeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *SafeBuffer) String() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.buf.String()
}

func (b *SafeBuffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.Reset()
}

func (b *SafeBuffer) Bytes() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]byte(nil), b.buf.Bytes()...)
}

func publishBidderDeposited(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	br *bidderregistry.BidderregistryBidderDeposited,
) error {
	event := brABI.Events["BidderDeposited"]

	newAvail := br.NewAvailableAmount
	if newAvail == nil {
		newAvail = big.NewInt(0)
	}
	buf, err := event.Inputs.NonIndexed().Pack(newAvail)
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
	buf, err := event.Inputs.NonIndexed().Pack(br.AmountStillEscrowed)
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

func TestPendingRefunds(t *testing.T) {
	t.Parallel()

	logger := util.NewTestLogger(io.Discard)
	st := depositstore.New(inmemstorage.New())
	dm := depositmanager.NewDepositManager(st, nil, nil, logger)

	digest1 := depositmanager.CommitmentDigest{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
	}

	pendingRefund, ok := dm.GetPendingRefund(digest1)
	if ok {
		t.Fatal("expected no pending refunds")
	}
	if pendingRefund != (depositmanager.PendingRefund{}) {
		t.Fatal("expected zeroed pending refund")
	}

	dm.AddPendingRefund(digest1, common.HexToAddress("0x123"), common.HexToAddress("0x456"), big.NewInt(100))

	pendingRefund, ok = dm.GetPendingRefund(digest1)
	if !ok {
		t.Fatal("expected pending refund")
	}
	if pendingRefund.Bidder != common.HexToAddress("0x123") {
		t.Fatal("expected bidder 0x123")
	}
	if pendingRefund.Provider != common.HexToAddress("0x456") {
		t.Fatal("expected provider 0x456")
	}
	if pendingRefund.Amount.Cmp(big.NewInt(100)) != 0 {
		t.Fatal("expected amount 100")
	}

	err := st.SetBalance(common.HexToAddress("0x123"), common.HexToAddress("0x456"), big.NewInt(77))
	if err != nil {
		t.Fatal(err)
	}

	balance, err := st.GetBalance(common.HexToAddress("0x123"), common.HexToAddress("0x456"))
	if err != nil {
		t.Fatal(err)
	}
	if balance.Cmp(big.NewInt(77)) != 0 {
		t.Fatal("expected balance 77")
	}

	err = dm.ApplyPendingRefund(digest1)
	if err != nil {
		t.Fatal(err)
	}

	balance, err = st.GetBalance(common.HexToAddress("0x123"), common.HexToAddress("0x456"))
	if err != nil {
		t.Fatal(err)
	}
	if balance.Cmp(big.NewInt(177)) != 0 {
		t.Fatal("expected balance 177")
	}

	digest2 := depositmanager.CommitmentDigest{
		33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64,
	}

	dm.AddPendingRefund(digest2, common.HexToAddress("0x123"), common.HexToAddress("0x456"), big.NewInt(302))

	pendingRefund, ok = dm.GetPendingRefund(digest2)
	if !ok {
		t.Fatal("expected pending refund")
	}
	if pendingRefund.Bidder != common.HexToAddress("0x123") {
		t.Fatal("expected bidder 0x123")
	}
	if pendingRefund.Provider != common.HexToAddress("0x456") {
		t.Fatal("expected provider 0x456")
	}
	if pendingRefund.Amount.Cmp(big.NewInt(302)) != 0 {
		t.Fatal("expected amount 302")
	}

	err = dm.DropPendingRefund(digest2)
	if err != nil {
		t.Fatal(err)
	}

	pendingRefund, ok = dm.GetPendingRefund(digest2)
	if ok {
		t.Fatal("expected no pending refund")
	}
}
