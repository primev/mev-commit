package transfer_test

import (
	"context"
	"log/slog"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/transfer"
)

type MockEthClient struct{}

func (m *MockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

func (m *MockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func (m *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return &types.Receipt{
		Status: 1,
	}, nil
}

func (m *MockEthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte{}, nil
}

type MockKeySigner struct {
	transaction *types.Transaction
}

func (m *MockKeySigner) GetAddress() common.Address {
	return common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
}

func (m *MockKeySigner) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	m.transaction = tx
	return tx, nil
}

func TestTransferrer(t *testing.T) {
	t.Parallel()

	// Mock the EthClient and KeySigner interfaces
	client := new(MockEthClient)
	signer := new(MockKeySigner)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	gasTip := big.NewInt(1000000000)    // 1 Gwei
	gasFeeCap := big.NewInt(2000000000) // 2 Gwei
	transferer := transfer.NewTransferer(logger, client, signer, gasTip, gasFeeCap)

	// Mock the context
	ctx := context.Background()
	// Mock the address and amount
	to := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	chainID := big.NewInt(1)                  // Mainnet
	amount := big.NewInt(1000000000000000000) // 1 ETH

	// Call the Transfer method
	err := transferer.Transfer(ctx, to, chainID, amount)
	if err != nil {
		t.Fatalf("Transfer failed: %v", err)
	}

	// Check if the transaction was signed correctly
	if signer.transaction == nil {
		t.Fatal("Transaction was not signed")
	}
	// Check if the transaction was sent correctly
	if signer.transaction.To() == nil {
		t.Fatal("Transaction was not sent")
	}
	if signer.transaction.Value().Cmp(amount) != 0 {
		t.Fatalf("Transaction amount mismatch: expected %s, got %s", amount.String(), signer.transaction.Value().String())
	}
	if signer.transaction.GasTipCap().Cmp(gasTip) != 0 {
		t.Fatalf("Transaction gas tip cap mismatch: expected %s, got %s", gasTip.String(), signer.transaction.GasTipCap().String())
	}
	if signer.transaction.GasFeeCap().Cmp(gasFeeCap) != 0 {
		t.Fatalf("Transaction gas fee cap mismatch: expected %s, got %s", gasFeeCap.String(), signer.transaction.GasFeeCap().String())
	}
	// Check if the transaction was sent to the correct address
	if signer.transaction.To().Hex() != to.Hex() {
		t.Fatalf("Transaction to address mismatch: expected %s, got %s", to.Hex(), signer.transaction.To().Hex())
	}
	// Check if the transaction was sent with the correct chain ID
	if signer.transaction.ChainId().Cmp(chainID) != 0 {
		t.Fatalf("Transaction chain ID mismatch: expected %s, got %s", chainID.String(), signer.transaction.ChainId().String())
	}
}
