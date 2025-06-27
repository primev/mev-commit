package store_test

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	_ "github.com/lib/pq"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"github.com/primev/mev-commit/tools/preconf-rpc/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestStore(t *testing.T) {
	ctx := context.Background()

	// Define the PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// Start the PostgreSQL container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %s", err)
	}
	defer func() {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			t.Errorf("Failed to terminate PostgreSQL container: %s", err)
		}
	}()

	// Retrieve the container's mapped port
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}
	// Construct the database connection string
	connStr := fmt.Sprintf("postgresql://user:password@localhost:%s/testdb?sslmode=disable", mappedPort.Port())

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL container: %s", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL container: %s", err)
	}

	st, err := store.New(db)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Errorf("failed to close store: %v", err)
		}
	})

	// Test data common for all tests
	txn1 := types.NewTransaction(
		0,
		common.HexToAddress("0x1234567890123456789012345678901234567890"),
		big.NewInt(1000000000), // 1 Gwei
		21000,                  // gas limit
		big.NewInt(1000000000), // gas price
		nil,                    // no data
	)
	rawTxn1, err := txn1.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal transaction: %v", err)
	}
	wrappedTxn1 := &sender.Transaction{
		Transaction: txn1,
		Raw:         hex.EncodeToString(rawTxn1),
		Sender:      common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		Type:        sender.TxTypeRegular,
		Status:      sender.TxStatusPending,
	}

	txn2 := types.NewTransaction(
		1,
		common.HexToAddress("0x0987654321098765432109876543210987654321"),
		big.NewInt(2000000000), // 2 Gwei
		21000,                  // gas limit
		big.NewInt(2000000000), // gas price
		nil,                    // no data
	)
	rawTxn2, err := txn2.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal second transaction: %v", err)
	}
	wrappedTxn2 := &sender.Transaction{
		Transaction: txn2,
		Raw:         hex.EncodeToString(rawTxn2),
		Sender:      common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		Type:        sender.TxTypeRegular,
		Status:      sender.TxStatusPending,
	}

	commitments := []*bidderapiv1.Commitment{
		{
			TxHashes:             []string{txn1.Hash().Hex()},
			BidAmount:            big.NewInt(1000000000).String(),
			BlockNumber:          1,
			ReceivedBidDigest:    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			ReceivedBidSignature: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			CommitmentDigest:     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			CommitmentSignature:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			DecayStartTimestamp:  time.Now().UnixMilli(),
			DecayEndTimestamp:    time.Now().Add(24 * time.Hour).UnixMilli(),
		},
		{
			TxHashes:             []string{txn1.Hash().Hex()},
			BidAmount:            big.NewInt(1000000000).String(),
			BlockNumber:          1,
			ReceivedBidDigest:    "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			ReceivedBidSignature: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			CommitmentDigest:     "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			CommitmentSignature:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			DecayStartTimestamp:  time.Now().UnixMilli(),
			DecayEndTimestamp:    time.Now().Add(24 * time.Hour).UnixMilli(),
		},
	}

	t.Run("AddQueuedTransaction", func(t *testing.T) {
		err := st.AddQueuedTransaction(context.Background(), wrappedTxn1)
		if err != nil {
			t.Fatalf("failed to add queued transaction: %v", err)
		}

		err = st.AddQueuedTransaction(context.Background(), wrappedTxn1) // Adding the same transaction again
		if err == nil {
			t.Fatalf("expected error when adding duplicate transaction, got nil")
		}

		err = st.AddQueuedTransaction(context.Background(), wrappedTxn2)
		if err != nil {
			t.Fatalf("failed to add second queued transaction: %v", err)
		}
	})

	t.Run("GetCurrentNonce", func(t *testing.T) {
		senderAddress := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
		nonce := st.GetCurrentNonce(context.Background(), senderAddress)
		if nonce != 2 {
			t.Fatalf("expected nonce 2, got %d", nonce)
		}
	})

	t.Run("GetTransactionByHash", func(t *testing.T) {
		retrievedTxn, err := st.GetTransactionByHash(context.Background(), wrappedTxn1.Transaction.Hash())
		if err != nil {
			t.Fatalf("failed to get transaction by hash: %v", err)
		}
		if diff := cmp.Diff(wrappedTxn1, retrievedTxn, cmpopts.IgnoreUnexported(sender.Transaction{}, types.Transaction{})); diff != "" {
			t.Fatalf("transaction mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetQueuedTransactions", func(t *testing.T) {
		retrievedTxns, err := st.GetQueuedTransactions(context.Background())
		if err != nil {
			t.Fatalf("failed to get queued transactions: %v", err)
		}
		if len(retrievedTxns) != 1 {
			t.Fatalf("expected 1 queued transaction, got %d", len(retrievedTxns))
		}
		if diff := cmp.Diff(wrappedTxn1, retrievedTxns[0], cmpopts.IgnoreUnexported(sender.Transaction{}, types.Transaction{})); diff != "" {
			t.Fatalf("queued transaction mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("StoreTransaction", func(t *testing.T) {
		wrappedTxn1.Status = sender.TxStatusPreConfirmed
		wrappedTxn1.BlockNumber = 1

		err := st.StoreTransaction(context.Background(), wrappedTxn1, commitments)
		if err != nil {
			t.Errorf("failed to store preconfirmed transaction: %v", err)
		}

		commitments, err := st.GetTransactionCommitments(context.Background(), wrappedTxn1.Transaction.Hash())
		if err != nil {
			t.Errorf("failed to get transaction commitments: %v", err)
		}
		if len(commitments) != 2 {
			t.Errorf("expected 2 commitments, got %d", len(commitments))
		}
		for i, commitment := range commitments {
			if diff := cmp.Diff(commitment, commitments[i], cmpopts.IgnoreUnexported(bidderapiv1.Commitment{}, types.Transaction{})); diff != "" {
				t.Errorf("commitment mismatch (-want +got):\n%s", diff)
			}
		}

		nextTxns, err := st.GetQueuedTransactions(context.Background())
		if err != nil {
			t.Errorf("failed to get queued transactions: %v", err)
		}
		if len(nextTxns) != 1 {
			t.Errorf("expected 1 queued transaction, got %d", len(nextTxns))
		}
		if diff := cmp.Diff(wrappedTxn2, nextTxns[0], cmpopts.IgnoreUnexported(sender.Transaction{}, types.Transaction{})); diff != "" {
			t.Errorf("queued transaction mismatch (-want +got):\n%s", diff)
		}

		txns, err := st.GetTransactionsForBlock(context.Background(), 1)
		if err != nil {
			t.Errorf("failed to get transactions for block: %v", err)
		}
		if len(txns) != 1 {
			t.Errorf("expected 1 transaction for block 1, got %d", len(txns))
		}
		if diff := cmp.Diff(wrappedTxn1, txns[0], cmpopts.IgnoreUnexported(sender.Transaction{}, types.Transaction{})); diff != "" {
			t.Errorf("transaction mismatch (-want +got):\n%s", diff)
		}

		wrappedTxn2.Status = sender.TxStatusFailed
		wrappedTxn2.Details = "Transaction failed due to insufficient funds"
		wrappedTxn2.BlockNumber = 2
		err = st.StoreTransaction(context.Background(), wrappedTxn2, nil)
		if err != nil {
			t.Errorf("failed to store failed transaction: %v", err)
		}

		failedTxn, err := st.GetTransactionByHash(context.Background(), wrappedTxn2.Transaction.Hash())
		if err != nil {
			t.Errorf("failed to get failed transaction by hash: %v", err)
		}

		if diff := cmp.Diff(wrappedTxn2, failedTxn, cmpopts.IgnoreUnexported(sender.Transaction{}, types.Transaction{})); diff != "" {
			t.Errorf("failed transaction mismatch (-want +got):\n%s", diff)
		}

		noTxns, err := st.GetTransactionsForBlock(context.Background(), 2)
		if err != nil {
			t.Errorf("failed to get transactions for block 2: %v", err)
		}
		if len(noTxns) != 0 {
			t.Errorf("expected no transactions for block 2, got %d", len(noTxns))
		}
	})

	t.Run("Account Balance", func(t *testing.T) {
		address := common.HexToAddress("0x1234567890123456789012345678901234567890")
		initialBalance := big.NewInt(1000000000) // 1 Gwei

		err := st.AddBalance(context.Background(), address, initialBalance)
		if err != nil {
			t.Errorf("failed to add balance: %v", err)
		}

		if !st.HasBalance(context.Background(), address, initialBalance) {
			t.Errorf("expected balance %s, but has no balance", initialBalance.String())
		}

		// Check if the balance is correctly stored
		balance, err := st.GetBalance(context.Background(), address)
		if err != nil {
			t.Errorf("failed to get balance: %v", err)
		}
		if balance.Cmp(initialBalance) != 0 {
			t.Errorf("expected balance %s, got %s", initialBalance.String(), balance.String())
		}

		err = st.DeductBalance(context.Background(), address, initialBalance)
		if err != nil {
			t.Errorf("failed to deduct balance: %v", err)
		}

		if st.HasBalance(context.Background(), address, initialBalance) {
			t.Errorf("expected no balance after deduction, but still has %s", initialBalance.String())
		}
	})
}
