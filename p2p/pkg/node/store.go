package node

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/primev/mev-commit/p2p/pkg/storage"
)

const (
	progressNS           = "p2progress/"
	progressLastBlockKey = progressNS + "last_block"
)

type DurableProgressStore struct {
	contractRPC ContractRPC
	kv          storage.Storage
}

type ContractRPC interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

func NewDurableProgressStore(kv storage.Storage, contractRPC ContractRPC) *DurableProgressStore {
	return &DurableProgressStore{
		contractRPC: contractRPC,
		kv:          kv,
	}
}

func (p *DurableProgressStore) LastBlock() (uint64, error) {
	buf, err := p.kv.Get(progressLastBlockKey)
	switch {
	case err == nil:
		if len(buf) != 8 {
			return 0, fmt.Errorf("invalid %q length: got %d, want 8", progressLastBlockKey, len(buf))
		}
		return binary.BigEndian.Uint64(buf), nil
	case errors.Is(err, storage.ErrKeyNotFound):
		return p.contractRPC.BlockNumber(context.Background())
	default:
		return 0, err
	}
}

func (p *DurableProgressStore) SetLastBlock(block uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], block)
	return p.kv.Put(progressLastBlockKey, b[:])
}
