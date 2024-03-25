package evmclient_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient/mockevm"
	mockkeysigner "github.com/primevprotocol/mev-commit/p2p/pkg/keysigner/mock"
	"github.com/primevprotocol/mev-commit/p2p/pkg/util"
)

func TestSendCall(t *testing.T) {
	t.Parallel()

	owner := common.HexToAddress("0xab")
	callData := []byte("call data")
	nonce := uint64(1)
	chainID := big.NewInt(1)
	unblockMonitor := make(chan struct{})

	pk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	ks := mockkeysigner.NewMockKeySigner(pk, owner)

	evm := mockevm.NewMockEvm(
		chainID.Uint64(),
		mockevm.WithPendingNonceAtFunc(
			func(ctx context.Context, account common.Address) (uint64, error) {
				if account != owner {
					t.Errorf("expected owner to be %v, got %v", owner, account)
				}
				return nonce, nil
			},
		),
		mockevm.WithEstimateGasFunc(
			func(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
				if call.From != owner {
					return 0, fmt.Errorf("expected from to be %v, got %v", owner, call.From)
				}
				if call.To.Hex() != owner.Hex() {
					return 0, fmt.Errorf("expected to to be %v, got %v", owner, call.To)
				}
				if string(call.Data) != string(callData) {
					return 0, fmt.Errorf("expected call data to be %v, got %v", callData, call.Data)
				}
				return 21000, nil
			},
		),
		mockevm.WithSuggestGasPriceFunc(
			func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(2000000000), nil
			},
		),
		mockevm.WithSuggestGasTipCapFunc(
			func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(1000000000), nil
			},
		),
		mockevm.WithSendTransactionFunc(
			func(ctx context.Context, tx *types.Transaction) error {
				if tx.GasFeeCap().Cmp(big.NewInt(2000000000)) != 0 {
					return fmt.Errorf("expected gas price to be 2000000000, got %v", tx.GasPrice())
				}
				if tx.GasTipCap().Cmp(big.NewInt(1000000000)) != 0 {
					return fmt.Errorf("expected gas tip cap to be 1000000000, got %v", tx.GasTipCap())
				}
				if tx.Gas() != 21000 {
					return fmt.Errorf("expected gas to be 21000, got %v", tx.Gas())
				}
				return nil
			},
		),
		mockevm.WithBlockNumFunc(
			func(ctx context.Context) (uint64, error) {
				select {
				case <-unblockMonitor:
					return 1, nil
				case <-ctx.Done():
					return 0, ctx.Err()
				}
			},
		),
		mockevm.WithNonceAtFunc(
			func(ctx context.Context, account common.Address, blockNum *big.Int) (uint64, error) {
				if account != owner {
					return 0, fmt.Errorf("expected owner to be %v, got %v", owner, account)
				}
				if blockNum.Uint64() != 1 {
					return 0, fmt.Errorf("expected blockNum to be 1, got %v", blockNum)
				}
				return nonce + 1, nil
			},
		),
		mockevm.WithBatcherFunc(
			func(ctx context.Context, elems []rpc.BatchElem) error {
				if len(elems) != 1 {
					return fmt.Errorf("expected 1 batch elem, got %v", len(elems))
				}
				if elems[0].Method != "eth_getTransactionReceipt" {
					return fmt.Errorf(
						"expected method to be eth_getTransactionReceipt, got %v",
						elems[0].Method,
					)
				}
				if len(elems[0].Args) != 1 {
					return fmt.Errorf("expected 1 arg, got %v", len(elems[0].Args))
				}
				elems[0].Result.(*types.Receipt).Status = 1
				return nil
			},
		),
		mockevm.WithCallContractFunc(
			func(ctx context.Context, call ethereum.CallMsg, blockNum *big.Int) ([]byte, error) {
				if call.From != owner {
					return nil, fmt.Errorf("expected from to be %v, got %v", owner, call.From)
				}
				if call.To.Hex() != owner.Hex() {
					return nil, fmt.Errorf("expected to to be %v, got %v", owner, call.To)
				}
				if string(call.Data) != string(callData) {
					return nil, fmt.Errorf("expected call data to be %v, got %v", callData, call.Data)
				}
				return []byte("result"), nil
			},
		),
	)

	client, err := evmclient.New(ks, evm, util.NewTestLogger(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := client.Send(context.Background(), &evmclient.TxRequest{
		To:       &owner,
		CallData: callData,
		Value:    big.NewInt(0),
	})
	if err != nil {
		t.Fatal(err)
	}

	txns := client.PendingTxns()
	if len(txns) != 1 {
		t.Fatalf("expected 1 pending txn, got %v", len(txns))
	}

	if txns[0].Hash != txHash.Hex() {
		t.Errorf("expected hash to be %v, got %v", txHash, txns[0].Hash)
	}

	if txns[0].Nonce != nonce {
		t.Errorf("expected nonce to be %v, got %v", nonce, txns[0].Nonce)
	}

	close(unblockMonitor)

	start := time.Now()
	for {
		if len(client.PendingTxns()) == 0 {
			break
		}
		if time.Since(start) > 5*time.Second {
			t.Fatal("timed out waiting for pending txns to be removed")
		}
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := client.Call(context.Background(), &evmclient.TxRequest{
		To:       &owner,
		CallData: callData,
	})
	if err != nil {
		t.Fatal(err)
	}

	if string(resp) != "result" {
		t.Errorf("expected result to be %v, got %v", "result", string(resp))
	}

	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCancel(t *testing.T) {
	t.Parallel()

	owner := common.HexToAddress("0xab")
	callData := []byte("call data")
	nonce := uint64(1)
	chainID := big.NewInt(1)
	unblockMonitor := make(chan struct{})
	successHash := common.HexToHash("0x123")
	blkNum := uint64(1)

	pk, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	ks := mockkeysigner.NewMockKeySigner(pk, owner)

	evm := mockevm.NewMockEvm(
		chainID.Uint64(),
		mockevm.WithPendingNonceAtFunc(
			func(ctx context.Context, account common.Address) (uint64, error) {
				return nonce, nil
			},
		),
		mockevm.WithSuggestGasTipCapFunc(
			func(ctx context.Context) (*big.Int, error) {
				return big.NewInt(1000000000), nil
			},
		),
		mockevm.WithSendTransactionFunc(
			func(ctx context.Context, tx *types.Transaction) error {
				if bytes.Equal(tx.Data(), callData) {
					return nil
				}
				if tx.GasFeeCap().Cmp(big.NewInt(2000000000)) <= 0 {
					return fmt.Errorf("expected gas price to be 2000000000, got %v", tx.GasFeeCap())
				}
				if tx.GasTipCap().Cmp(big.NewInt(1000000000)) <= 0 {
					return fmt.Errorf("expected gas tip cap to be 1000000000, got %v", tx.GasTipCap())
				}
				return nil
			},
		),
		mockevm.WithBlockNumFunc(
			func(ctx context.Context) (uint64, error) {
				select {
				case <-unblockMonitor:
					defer func() { blkNum++ }()
					return blkNum, nil
				case <-ctx.Done():
					return 0, ctx.Err()
				}
			},
		),
		mockevm.WithNonceAtFunc(
			func(ctx context.Context, account common.Address, blockNum *big.Int) (uint64, error) {
				return nonce + 1, nil
			},
		),
		mockevm.WithBatcherFunc(
			func(ctx context.Context, elems []rpc.BatchElem) error {
				for i, elem := range elems {
					if elem.Args[0].(common.Hash) != successHash {
						elems[i].Error = ethereum.NotFound
					} else {
						elems[i].Result.(*types.Receipt).Status = 1
					}
				}
				return nil
			},
		),
		mockevm.WithTransactionByHashFunc(
			func(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
				return types.NewTx(&types.DynamicFeeTx{
					ChainID:   chainID,
					Nonce:     nonce,
					Data:      callData,
					GasFeeCap: big.NewInt(2000000000),
					GasTipCap: big.NewInt(1000000000),
					Gas:       21000,
					Value:     big.NewInt(0),
				}), true, nil
			},
		),
	)

	client, err := evmclient.New(ks, evm, util.NewTestLogger(os.Stdout))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	txHash, err := client.Send(ctx, &evmclient.TxRequest{
		To:       &owner,
		CallData: callData,
		GasLimit: 21000,
		GasPrice: big.NewInt(1000000000),
		Value:    big.NewInt(0),
	})
	if err != nil {
		t.Fatal(err)
	}

	errC := make(chan error, 1)
	go func() {
		_, err := client.WaitForReceipt(ctx, txHash)
		if err != nil {
			errC <- err
			return
		}
		errC <- nil
	}()

	cancelHash, err := client.CancelTx(ctx, txHash)
	if err != nil {
		t.Fatal(err)
	}

	successHash = cancelHash

	txns := client.PendingTxns()
	if len(txns) != 2 {
		t.Fatalf("expected 1 pending txn, got %v", len(txns))
	}
	for _, txn := range txns {
		if txn.Hash != txHash.Hex() && txn.Hash != cancelHash.Hex() {
			t.Errorf("expected hash to be %v or %v, got %v", txHash, cancelHash, txn.Hash)
		}
		if txn.Nonce != nonce {
			t.Errorf("expected nonce to be %v, got %v", nonce, txn.Nonce)
		}
	}

	close(unblockMonitor)

	res, err := client.WaitForReceipt(ctx, cancelHash)
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 1 {
		t.Errorf("expected status to be 1, got %v", res.Status)
	}

	select {
	case err := <-errC:
		if !errors.Is(err, evmclient.ErrTxnCancelled) {
			t.Fatalf("expected error to be %v, got %v", evmclient.ErrTxnCancelled, err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for receipt")
	}
}
