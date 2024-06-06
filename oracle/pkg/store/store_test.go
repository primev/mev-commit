package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/primev/mev-commit/oracle/pkg/store"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type blockWinner struct {
	BlockNumber int64
	Winner      []byte
	Window      int64
}

type testSettlement struct {
	CommitmentIdx   []byte
	TxHash          string
	BlockNum        int64
	Builder         []byte
	Amount          uint64
	BidID           []byte
	Type            updater.SettlementType
	DecayPercentage int64
	ChainHash       []byte
	Nonce           uint64
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

	winners := []blockWinner{
		{
			Window:      1,
			Winner:      common.HexToAddress("0x01").Bytes(),
			BlockNumber: 1,
		},
		{
			Window:      2,
			Winner:      common.HexToAddress("0x02").Bytes(),
			BlockNumber: 2,
		},
	}

	settlements := []testSettlement{
		{
			CommitmentIdx: []byte{1},
			TxHash:        common.HexToHash("0x01").String(),
			BlockNum:      1,
			Amount:        2000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x01").Bytes(),
			Type:          updater.SettlementTypeReward,
			ChainHash:     common.HexToHash("0x01").Bytes(),
			Nonce:         1,
		},
		{
			CommitmentIdx: []byte{2},
			TxHash:        common.HexToHash("0x02").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x02").Bytes(),
			Type:          updater.SettlementTypeSlash,
			ChainHash:     common.HexToHash("0x02").Bytes(),
			Nonce:         2,
		},
		{
			CommitmentIdx: []byte{3},
			TxHash:        common.HexToHash("0x03").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x03").Bytes(),
			Type:          updater.SettlementTypeReward,
			ChainHash:     common.HexToHash("0x03").Bytes(),
			Nonce:         3,
		},
		{
			CommitmentIdx: []byte{4},
			TxHash:        common.HexToHash("0x04").String(),
			BlockNum:      2,
			Amount:        2000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x04").Bytes(),
			Type:          updater.SettlementTypeSlash,
			ChainHash:     common.HexToHash("0x04").Bytes(),
			Nonce:         4,
		},
		{
			CommitmentIdx: []byte{5},
			TxHash:        common.HexToHash("0x05").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x05").Bytes(),
			Type:          updater.SettlementTypeReward,
			ChainHash:     common.HexToHash("0x05").Bytes(),
			Nonce:         5,
		},
		{
			CommitmentIdx: []byte{6},
			TxHash:        common.HexToHash("0x06").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x04").Bytes(),
			Type:          updater.SettlementTypeSlash,
			ChainHash:     common.HexToHash("0x06").Bytes(),
			Nonce:         6,
		},
	}

	t.Run("NewStore", func(t *testing.T) {
		// Create the store and tables
		_, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}
	})

	t.Run("RegisterWinner", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, winner := range winners {
			err = st.RegisterWinner(context.Background(), winner.BlockNumber, winner.Winner, winner.Window)
			if err != nil {
				t.Fatalf("Failed to register winner: %s", err)
			}
		}
	})

	t.Run("GetWinner", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, winner := range winners {
			w, err := st.GetWinner(context.Background(), winner.BlockNumber)
			if err != nil {
				t.Fatalf("Failed to get winner: %s", err)
			}
			if diff := cmp.Diff(w.Winner, winner.Winner); diff != "" {
				t.Fatalf("Unexpected winner: (-want +have):\n%s", diff)
			}
		}
	})

	t.Run("AddEncryptedCommitment", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, settlement := range settlements {
			err = st.AddEncryptedCommitment(
				context.Background(),
				settlement.CommitmentIdx,
				settlement.Builder,
				[]byte("hash"),
				[]byte("signature"),
				uint64(time.Now().Unix()),
			)
			if err != nil {
				t.Fatalf("Failed to add encrypted commitment: %s", err)
			}
		}
	})

	t.Run("AddSettlement", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for i, settlement := range settlements {
			window := int64(i/3) + 1
			err = st.AddSettlement(
				context.Background(),
				settlement.CommitmentIdx,
				settlement.TxHash,
				settlement.BlockNum,
				settlement.Amount,
				settlement.Builder,
				settlement.BidID,
				settlement.Type,
				settlement.DecayPercentage,
				window,
				settlement.ChainHash,
				settlement.Nonce,
			)
			if err != nil {
				t.Fatalf("Failed to add settlement: %s", err)
			}
		}
	})

	t.Run("Save", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, s := range settlements {
			err = st.Save(context.Background(), common.Hash(s.ChainHash), s.Nonce)
			if err != nil {
				t.Fatalf("Failed to mark txn sent: %s", err)
			}
		}
	})

	t.Run("IsSettled", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, settlement := range settlements {
			settled, err := st.IsSettled(context.Background(), settlement.CommitmentIdx)
			if err != nil {
				t.Fatalf("Failed to check if settled: %s", err)
			}
			if !settled {
				t.Fatalf("Expected settlement to be settled")
			}
		}
	})

	t.Run("Settlement", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, settlement := range settlements {
			s, err := st.Settlement(context.Background(), settlement.CommitmentIdx)
			if err != nil {
				t.Fatalf("Failed to get settlement: %s", err)
			}

			if diff := cmp.Diff(s, updater.Settlement{
				CommitmentIdx:   settlement.CommitmentIdx,
				TxHash:          settlement.TxHash,
				BlockNum:        settlement.BlockNum,
				Builder:         settlement.Builder,
				Amount:          settlement.Amount,
				BidID:           settlement.BidID,
				Type:            settlement.Type,
				DecayPercentage: settlement.DecayPercentage,
			}); diff != "" {
				t.Fatalf("Unexpected settlement: (-want +have):\n%s", diff)
			}
		}
	})

	t.Run("LastBlock and SetBlockNo", func(t *testing.T) {
		st, err := store.NewStore(db)
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
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, s := range settlements {
			err = st.Update(context.Background(), common.Hash(s.ChainHash), "success")
			if err != nil {
				t.Fatalf("Failed to mark txn sent: %s", err)
			}
		}

		pendingTxnCount, err := st.PendingTxnCount()
		if err != nil {
			t.Fatalf("Failed to get pending txn count: %s", err)
		}
		if pendingTxnCount != 0 {
			t.Fatalf("Expected pending txn count 0, got %d", pendingTxnCount)
		}
	})
}
