package blocktracker_test

import (
	"context"
	"hash"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/tools/preconf-rpc/blocktracker"
	"github.com/primev/mev-commit/x/util"
	"golang.org/x/crypto/sha3"
)

type mockEthClient struct {
	blockNumber chan uint64
	blocks      map[uint64]*types.Block
}

type subscriptionNoOp struct{}

func (s *subscriptionNoOp) Unsubscribe() {}
func (s *subscriptionNoOp) Err() <-chan error {
	return nil
}

func (m *mockEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	go func() {
		for {
			select {
			case blockNo := <-m.blockNumber:
				block, exists := m.blocks[blockNo]
				if !exists {
					continue
				}
				ch <- block.Header()
			case <-ctx.Done():
				return
			}
		}
	}()
	return &subscriptionNoOp{}, nil
}

func (m *mockEthClient) BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	block, exists := m.blocks[blockNumber.Uint64()]
	if !exists {
		return nil, nil // Simulate block not found
	}
	return block, nil
}

func (m *mockEthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

type testHasher struct {
	hasher hash.Hash
}

// NewHasher returns a new testHasher instance.
func NewHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

// Reset resets the hash state.
func (h *testHasher) Reset() {
	h.hasher.Reset()
}

// Update updates the hash state with the given key and value.
func (h *testHasher) Update(key, val []byte) error {
	h.hasher.Write(key)
	h.hasher.Write(val)
	return nil
}

// Hash returns the hash value.
func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}

func TestBlockTracker(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	tx1 := types.NewTransaction(1, common.HexToAddress("0xabc"), big.NewInt(100), 21000, big.NewInt(1), nil)
	tx2 := types.NewTransaction(2, common.HexToAddress("0xdef"), big.NewInt(200), 21000, big.NewInt(1), nil)
	tx3 := types.NewTransaction(3, common.HexToAddress("0x123"), big.NewInt(300), 21000, big.NewInt(1), nil)
	tx4 := types.NewTransaction(4, common.HexToAddress("0x456"), big.NewInt(400), 21000, big.NewInt(1), nil)

	blk1 := types.NewBlock(
		&types.Header{
			Number:   big.NewInt(100),
			Time:     uint64(time.Now().Unix()),
			BaseFee:  big.NewInt(1000000000),
			GasLimit: 30_000_000,
			GasUsed:  20_000_000,
		},
		&types.Body{Transactions: []*types.Transaction{tx1, tx2}},
		nil, // No receipts
		NewHasher(),
	)

	blk2 := types.NewBlock(
		&types.Header{
			Number:   big.NewInt(101),
			Time:     uint64(time.Now().Add(12 * time.Second).Unix()),
			BaseFee:  big.NewInt(1100000000),
			GasLimit: 30_000_000,
			GasUsed:  10_000_000,
		},
		&types.Body{Transactions: []*types.Transaction{tx3}},
		nil, // No receipts
		NewHasher(),
	)

	client := &mockEthClient{
		blockNumber: make(chan uint64, 1),
		blocks: map[uint64]*types.Block{
			100: blk1,
			101: blk2,
		},
	}

	tracker, err := blocktracker.NewBlockTracker(client, util.NewTestLogger(os.Stdout))
	if err != nil {
		t.Fatalf("Failed to create block tracker: %v", err)
	}
	done := tracker.Start(ctx)

	blkNo := tracker.LatestBlockNumber()
	if blkNo != 0 {
		t.Fatalf("Expected latest block number to be 0, got %d", blkNo)
	}

	included1 := tracker.WaitForTxnInclusion(tx1.Hash())
	included2 := tracker.WaitForTxnInclusion(tx2.Hash())
	included3 := tracker.WaitForTxnInclusion(tx3.Hash())
	included4 := tracker.WaitForTxnInclusion(tx4.Hash())

	client.blockNumber <- 100

	start := time.Now()
	for {
		bidBlockNo, duration, err := tracker.NextBlockNumber()
		if err == nil {
			if bidBlockNo != 101 {
				t.Fatalf("Expected next block number to be 101, got %d", bidBlockNo)
			}
			if duration <= 0 {
				t.Fatalf("Expected positive duration, got %v", duration)
			}
			break
		} else {
			t.Logf("Waiting for next block number: %v", err)
		}
		if time.Since(start) > 5*time.Second {
			t.Fatalf("Timeout waiting for next block number")
		}
		time.Sleep(100 * time.Millisecond)
	}

	bNo1 := <-included1
	if bNo1 != 100 {
		t.Fatalf("Expected transaction %s to be included in block 100, got %d", tx1.Hash().Hex(), bNo1)
	}
	bNo2 := <-included2
	if bNo2 != 100 {
		t.Fatalf("Expected transaction %s to be included in block 100, got %d", tx2.Hash().Hex(), bNo2)
	}

	blkNo = tracker.LatestBlockNumber()
	if blkNo != 100 {
		t.Fatalf("Expected latest block number to be 100, got %d", blkNo)
	}

	latest := tracker.LatestBaseFee()
	if latest.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("Expected latest base fee to be 1000000000, got %s", latest.String())
	}

	nextBase := tracker.NextBaseFee()
	expectedNextBase := big.NewInt(0).Add(
		big.NewInt(1000000000),
		big.NewInt(0).Div(
			big.NewInt(0).Mul(big.NewInt(1000000000), big.NewInt(5_000_000)),
			big.NewInt(0).Mul(big.NewInt(15_000_000), big.NewInt(8)),
		),
	)
	if nextBase.Cmp(expectedNextBase) != 0 {
		t.Fatalf("Expected next base fee to be %s, got %s", expectedNextBase.String(), nextBase.String())
	}

	client.blockNumber <- 101

	start = time.Now()
	for {
		bidBlockNo, duration, err := tracker.NextBlockNumber()
		if err == nil {
			if bidBlockNo == 102 && duration > 0 {
				break
			}
		} else {
			t.Logf("Waiting for next block number: %v", err)
		}
		if time.Since(start) > 5*time.Second {
			t.Fatalf("Timeout waiting for next block number")
		}
		time.Sleep(100 * time.Millisecond)
	}

	bNo3 := <-included3
	if bNo3 != 101 {
		t.Fatalf("Expected transaction %s to be included in block 101, got %d", tx3.Hash().Hex(), bNo3)
	}

	select {
	case bNo4 := <-included4:
		t.Fatalf("Did not expect transaction %s to be included, but got block number %d", tx4.Hash().Hex(), bNo4)
	case <-time.After(1 * time.Second):
		// Expected timeout
	}

	cancel()
	<-done // Wait for the tracker to finish
}
