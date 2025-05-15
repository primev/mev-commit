package localstate

import (
	"context"
	"log/slog"
	"sync"

	"github.com/primev/mev-commit/cl/types"
)

// LocalStateManager implements the blockbuilder.StateManager interface for single-node operation.
// It manages state in-memory.
type LocalStateManager struct {
	mu              sync.RWMutex
	blockBuildState *types.BlockBuildState
	logger          *slog.Logger
}

// NewLocalStateManager creates a new LocalStateManager.
func NewLocalStateManager(logger *slog.Logger) *LocalStateManager {
	return &LocalStateManager{
		// Initialize with a default state, typically to start building a block.
		blockBuildState: &types.BlockBuildState{
			CurrentStep: types.StepBuildBlock,
		},
		logger: logger,
	}
}

// SaveBlockStateAndPublishToStream saves the block state locally.
// The "PublishToStream" aspect is a NOP for this local manager.
func (lsm *LocalStateManager) SaveBlockStateAndPublishToStream(_ context.Context, state *types.BlockBuildState) error {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	lsm.blockBuildState = state // Store the provided state
	lsm.logger.Info(
		"LocalStateManager: Saved block state",
		"step", state.CurrentStep.String(),
		"payload_id", state.PayloadID,
	)
	return nil
}

// GetBlockBuildState retrieves the current block build state.
func (lsm *LocalStateManager) GetBlockBuildState(_ context.Context) types.BlockBuildState {
	lsm.mu.RLock()
	defer lsm.mu.RUnlock()

	if lsm.blockBuildState == nil {
		// This should ideally not happen if constructor initializes it.
		lsm.logger.Error("LocalStateManager: blockBuildState is nil, returning default. This indicates an issue.")
		return types.BlockBuildState{CurrentStep: types.StepBuildBlock}
	}
	// Return a copy to prevent external modification of the internal state.
	stateCopy := *lsm.blockBuildState
	return stateCopy
}

// ResetBlockState resets the block build state to the initial step (StepBuildBlock).
func (lsm *LocalStateManager) ResetBlockState(_ context.Context) error {
	lsm.mu.Lock()
	defer lsm.mu.Unlock()

	lsm.blockBuildState = &types.BlockBuildState{
		CurrentStep: types.StepBuildBlock,
	}
	return nil
}
