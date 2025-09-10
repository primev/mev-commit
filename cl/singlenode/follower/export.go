package follower

import (
	"context"

	"github.com/primev/mev-commit/cl/types"
)

func (f *Follower) PayloadCh() <-chan types.PayloadInfo {
	return f.payloadCh
}

func (f *Follower) SyncFromSharedDB(ctx context.Context) error {
	return f.syncFromSharedDB(ctx)
}

func (f *Follower) LastSignalledBlock() uint64 {
	return f.lastSignalledBlock.Load()
}

func (f *Follower) SetLastSignalledBlock(block uint64) {
	f.lastSignalledBlock.Store(block)
}

func (f *Follower) QueryPayloadsFromSharedDB(ctx context.Context) {
	f.queryPayloadsFromSharedDB(ctx)
}

func (f *Follower) GetLastProcessed(ctx context.Context) (uint64, error) {
	return f.getLastProcessed(ctx)
}

func (f *Follower) SetLastProcessed(ctx context.Context, height uint64) error {
	return f.setLastProcessed(ctx, height)
}
