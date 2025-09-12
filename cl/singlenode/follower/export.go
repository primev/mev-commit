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

func (f *Follower) GetExecutionHead() *types.ExecutionHead {
	return f.getExecutionHead()
}
