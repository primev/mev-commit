package node

import (
	"context"
	"encoding/binary"
	"testing"

	inmem "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
)

type mockContractRPC struct {
	blockNumber uint64
}

func (m *mockContractRPC) BlockNumber(ctx context.Context) (uint64, error) {
	return m.blockNumber, nil
}

func TestDurableProgressStore_LastBlock_FallbackToRPCWhenUnset(t *testing.T) {
	kv := inmem.New()
	mockRPC := &mockContractRPC{blockNumber: 12345}

	ps := NewDurableProgressStore(kv, mockRPC)

	got, err := ps.LastBlock()
	if err != nil {
		t.Fatalf("LastBlock: %v", err)
	}
	if got != 12345 {
		t.Fatalf("LastBlock fallback mismatch: got %d want %d", got, 12345)
	}
}

func TestDurableProgressStore_SetAndGet(t *testing.T) {
	kv := inmem.New()
	mockRPC := &mockContractRPC{blockNumber: 0}

	ps := NewDurableProgressStore(kv, mockRPC)

	want := uint64(9876543210)
	if err := ps.SetLastBlock(want); err != nil {
		t.Fatalf("SetLastBlock: %v", err)
	}

	got, err := ps.LastBlock()
	if err != nil {
		t.Fatalf("LastBlock: %v", err)
	}
	if got != uint64(9876543210) {
		t.Fatalf("LastBlock persisted mismatch: got %d want %d", got, want)
	}

	raw, err := kv.Get(progressLastBlockKey)
	if err != nil {
		t.Fatalf("kv.Get: %v", err)
	}
	if len(raw) != 8 {
		t.Fatalf("stored length mismatch: got %d want 8", len(raw))
	}
	if binary.BigEndian.Uint64(raw) != uint64(9876543210) {
		t.Fatalf("stored value mismatch: got %d want %d", binary.BigEndian.Uint64(raw), want)
	}
}
