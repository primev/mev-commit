package gwcontract_test

import (
	"context"
	"errors"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/primev/mev-commit/bridge/standard/pkg/gwcontract"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/util"
)

type testStorage struct {
	transfers []*transfer
}

type transfer struct {
	tx        *l1gateway.L1gatewayTransferInitiated
	chainhash common.Hash
	nonce     uint64
	settled   bool
}

func (s *testStorage) StoreTransfer(
	ctx context.Context,
	transferIdx *big.Int,
	amount *big.Int,
	recipient common.Address,
	nonce uint64,
	chainHash common.Hash,
) error {
	s.transfers = append(s.transfers, &transfer{
		tx: &l1gateway.L1gatewayTransferInitiated{
			Sender:      common.HexToAddress("0x1234"),
			Recipient:   recipient,
			Amount:      amount,
			TransferIdx: transferIdx,
		},
		chainhash: chainHash,
		nonce:     nonce,
	})
	return nil
}

func (s *testStorage) MarkTransferSettled(ctx context.Context, transferIdx *big.Int) error {
	for _, t := range s.transfers {
		if t.tx.TransferIdx.Cmp(transferIdx) == 0 {
			t.settled = true
			return nil
		}
	}

	return errors.New("transfer not found")
}

func (s *testStorage) IsSettled(ctx context.Context, transferIdx *big.Int) (bool, error) {
	for _, t := range s.transfers {
		if t.tx.TransferIdx.Cmp(transferIdx) == 0 {
			return t.settled, nil
		}
	}
	return false, nil
}

type testMonitor struct {
	errNonce uint64
}

func (m *testMonitor) WatchTx(hash common.Hash, nonce uint64) <-chan txmonitor.Result {
	ch := make(chan txmonitor.Result, 1)
	status := uint64(1)
	if nonce == m.errNonce {
		status = 0
	}

	ch <- txmonitor.Result{
		Receipt: &types.Receipt{
			Status: status,
		},
	}
	return ch
}

type testGatewayTransactor struct {
	nonce uint64
}

func (t *testGatewayTransactor) FinalizeTransfer(
	opts *bind.TransactOpts,
	recipient common.Address,
	amount *big.Int,
	counterpartyIdx *big.Int,
) (*types.Transaction, error) {
	newNonce := t.nonce
	t.nonce++
	return types.NewTransaction(
		newNonce,
		common.BigToAddress(big.NewInt(int64(t.nonce))),
		nil,
		0,
		nil,
		nil,
	), nil
}

func TestGateway(t *testing.T) {
	logger := util.NewTestLogger(os.Stdout)
	monitor := &testMonitor{errNonce: 3}
	st := &testStorage{}
	transactor := &testGatewayTransactor{}
	optsGetter := func(context.Context) (*bind.TransactOpts, error) {
		return nil, nil
	}

	brABI, err := abi.JSON(strings.NewReader(l1gateway.L1gatewayABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		logger,
		&brABI,
	)

	gw := gwcontract.NewGateway[l1gateway.L1gatewayTransferInitiated](
		logger,
		monitor,
		evtMgr,
		transactor,
		optsGetter,
		st,
	)

	transfers := []*l1gateway.L1gatewayTransferInitiated{
		{
			Sender:      common.HexToAddress("0x1234"),
			Recipient:   common.HexToAddress("0x5678"),
			Amount:      big.NewInt(100),
			TransferIdx: big.NewInt(1),
		},
		{
			Sender:      common.HexToAddress("0x5678"),
			Recipient:   common.HexToAddress("0x1234"),
			Amount:      big.NewInt(200),
			TransferIdx: big.NewInt(2),
		},
		{
			Sender:      common.HexToAddress("0x1234"),
			Recipient:   common.HexToAddress("0x5678"),
			Amount:      big.NewInt(300),
			TransferIdx: big.NewInt(3),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	tChan, errChan := gw.Subscribe(ctx)

	doneCollecting := make(chan struct{})
	idx := 0
	go func() {
		defer close(doneCollecting)

		for {
			if idx == len(transfers) {
				return
			}
			select {
			case tr := <-tChan:
				if diff := cmp.Diff(transfers[idx], tr, cmp.AllowUnexported(big.Int{})); diff != "" {
					t.Fatalf("unexpected transfer at index %d (-want +got):\n%s", idx, diff)
				}
				idx++
			case <-errChan:
				return
			}
		}
	}()

	for _, transfer := range transfers {
		if err := publishTransfer(evtMgr, &brABI, transfer); err != nil {
			t.Fatal(err)
		}
	}

	<-doneCollecting

	cancel()
	_, more := <-errChan
	if more {
		t.Fatal("expected channel to be closed")
	}

	ctx = context.Background()

	for idx, transfer := range transfers {
		if err := gw.FinalizeTransfer(
			ctx,
			transfer.Recipient,
			transfer.Amount,
			transfer.TransferIdx,
		); err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(st.transfers[idx].tx, transfer, cmp.Comparer(func(a, b *l1gateway.L1gatewayTransferInitiated) bool {
			return a.Recipient == b.Recipient &&
				a.Amount.Cmp(b.Amount) == 0 &&
				a.TransferIdx.Cmp(b.TransferIdx) == 0
		})); diff != "" {
			t.Fatalf("unexpected transfer at index %d (-want +got):\n%s", idx, diff)
		}
		if !st.transfers[idx].settled {
			t.Fatalf("expected transfer at index %d to be settled", idx)
		}
	}

	prevNonce := transactor.nonce
	if err := gw.FinalizeTransfer(
		ctx,
		transfers[0].Recipient,
		transfers[0].Amount,
		transfers[0].TransferIdx,
	); err != nil {
		t.Fatal(err)
	}
	if transactor.nonce != prevNonce {
		t.Fatalf("expected nonce to not be incremented")
	}

	if err := gw.FinalizeTransfer(
		ctx,
		common.HexToAddress("0x1234"),
		big.NewInt(100),
		big.NewInt(4),
	); err == nil {
		t.Fatal("expected error")
	}
}

func publishTransfer(
	evtMgr events.EventManager,
	brABI *abi.ABI,
	transfer *l1gateway.L1gatewayTransferInitiated,
) error {
	event := brABI.Events["TransferInitiated"]
	buf, err := event.Inputs.NonIndexed().Pack(
		transfer.Amount,
	)
	if err != nil {
		return err
	}

	sender := common.BytesToHash(transfer.Sender.Bytes())
	recipient := common.BytesToHash(transfer.Recipient.Bytes())

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			sender,   // The next topics are the indexed event parameters
			recipient,
			common.BigToHash(transfer.TransferIdx),
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}
