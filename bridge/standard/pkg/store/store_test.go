package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/bridge/standard/pkg/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testTransfer struct {
	Recipient   common.Address
	Amount      *big.Int
	TransferIdx *big.Int
	ChainHash   []byte
	Nonce       uint64
}

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
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL container: %s", err)
	}

	transfers := []testTransfer{
		{
			Recipient:   common.HexToAddress("0x1234"),
			Amount:      big.NewInt(100),
			TransferIdx: big.NewInt(1),
			ChainHash:   common.HexToHash("0x01").Bytes(),
			Nonce:       1,
		},
		{
			Recipient:   common.HexToAddress("0x5678"),
			Amount:      big.NewInt(200),
			TransferIdx: big.NewInt(2),
			ChainHash:   common.HexToHash("0x02").Bytes(),
			Nonce:       2,
		},
		{
			Recipient:   common.HexToAddress("0x9abc"),
			Amount:      big.NewInt(300),
			TransferIdx: big.NewInt(3),
			ChainHash:   common.HexToHash("0x03").Bytes(),
			Nonce:       3,
		},
		{
			Recipient:   common.HexToAddress("0xdef0"),
			Amount:      big.NewInt(400),
			TransferIdx: big.NewInt(4),
			ChainHash:   common.HexToHash("0x04").Bytes(),
			Nonce:       4,
		},
		{
			Recipient:   common.HexToAddress("0x1256"),
			Amount:      big.NewInt(500),
			TransferIdx: big.NewInt(5),
			ChainHash:   common.HexToHash("0x05").Bytes(),
			Nonce:       5,
		},
		{
			Recipient:   common.HexToAddress("0x789a"),
			Amount:      big.NewInt(600),
			TransferIdx: big.NewInt(6),
			ChainHash:   common.HexToHash("0x06").Bytes(),
			Nonce:       6,
		},
	}

	t.Run("NewStore", func(t *testing.T) {
		// Create the store and tables
		_, err := store.NewStore(db, "test")
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}
	})

	t.Run("StoreTransfer", func(t *testing.T) {
		st, err := store.NewStore(db, "test")
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, transfer := range transfers {
			err = st.StoreTransfer(
				context.Background(),
				transfer.TransferIdx,
				transfer.Amount,
				transfer.Recipient,
				transfer.Nonce,
				common.BytesToHash(transfer.ChainHash),
			)
			if err != nil {
				t.Fatalf("Failed to add encrypted commitment: %s", err)
			}
		}

		for _, transfer := range transfers {
			settled, err := st.IsSettled(context.Background(), transfer.TransferIdx)
			if err != nil {
				t.Fatalf("Failed to check if settled: %s", err)
			}
			if settled {
				t.Fatalf("Expected transfer to not be settled")
			}
		}

		for _, transfer := range transfers {
			err = st.MarkTransferSettled(context.Background(), transfer.TransferIdx)
			if err != nil {
				t.Fatalf("Failed to mark transfer settled: %s", err)
			}
		}

		for _, transfer := range transfers {
			settled, err := st.IsSettled(context.Background(), transfer.TransferIdx)
			if err != nil {
				t.Fatalf("Failed to check if settled: %s", err)
			}
			if !settled {
				t.Fatalf("Expected transfer to be settled")
			}
		}
	})

	t.Run("Save", func(t *testing.T) {
		st, err := store.NewStore(db, "test")
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, transfer := range transfers {
			err = st.Save(context.Background(), common.Hash(transfer.ChainHash), transfer.Nonce)
			if err != nil {
				t.Fatalf("Failed to mark txn sent: %s", err)
			}
		}
	})

	t.Run("LastBlock and SetBlockNo", func(t *testing.T) {
		st, err := store.NewStore(db, "test")
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		lastBlock, err := st.LastBlock()
		if err != nil {
			t.Fatalf("Failed to get last block: %s", err)
		}
		if lastBlock != 0 {
			t.Fatalf("Expected last block 0, got %d", lastBlock)
		}

		err = st.SetLastBlock(3)
		if err != nil {
			t.Fatalf("Failed to set block number: %s", err)
		}

		lastBlock, err = st.LastBlock()
		if err != nil {
			t.Fatalf("Failed to get last block: %s", err)
		}
		if lastBlock != 3 {
			t.Fatalf("Expected last block 3, got %d", lastBlock)
		}
	})

	t.Run("Update", func(t *testing.T) {
		st, err := store.NewStore(db, "test")
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		pendingTxns, err := st.PendingTxns()
		if err != nil {
			t.Fatalf("Failed to get pending txns: %s", err)
		}

		for _, transfer := range transfers {
			found := false
			for _, txn := range pendingTxns {
				if txn.Hash == common.BytesToHash(transfer.ChainHash) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected txn %s to be pending", common.BytesToHash(transfer.ChainHash).String())
			}
		}

		for _, transfer := range transfers {
			err = st.Update(context.Background(), common.Hash(transfer.ChainHash), "success")
			if err != nil {
				t.Fatalf("Failed to mark txn sent: %s", err)
			}
		}

		pendingTxns, err = st.PendingTxns()
		if err != nil {
			t.Fatalf("Failed to get pending txns: %s", err)
		}

		if len(pendingTxns) != 0 {
			t.Fatalf("Expected no pending txns, got %d", len(pendingTxns))
		}
	})
}
