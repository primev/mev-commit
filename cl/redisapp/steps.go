package redisapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/cl/redisapp/types"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/exp/rand"
)

const maxAttempts = 3

type BlockBuilder struct {
	stateManager StateManager
	engineCl     EngineClient
	logger       *slog.Logger
	buildDelay   time.Duration
	buildDelayMs uint64
	lastCallTime time.Time
	ctx          context.Context
}

func (bb *BlockBuilder) startBuild(ctx context.Context, feeRecipient common.Address, head *types.ExecutionHead, ts uint64) (engine.ForkChoiceResponse, error) {
	hash := common.BytesToHash(head.BlockHash)

	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	bb.logger.Info("Leader: Submit new EVM payload", "timestamp", ts)

	attrs := &engine.PayloadAttributes{
		Timestamp:             ts,
		Random:                hash, // We use head block hash as randao.
		SuggestedFeeRecipient: feeRecipient,
		Withdrawals:           []*etypes.Withdrawal{},
		BeaconRoot:            &hash,
	}

	resp, err := bb.engineCl.ForkchoiceUpdatedV3(ctx, fcs, attrs)
	if err != nil {
		return engine.ForkChoiceResponse{}, fmt.Errorf("forkchoice update, %w", err)
	}

	return resp, nil
}

func (bb *BlockBuilder) getPayload(ctx context.Context) error {
	var payloadID *engine.PayloadID

	currentCallTime := time.Now()

	// Load execution head to get previous block timestamp
	head, err := bb.stateManager.LoadExecutionHead(ctx)
	if err != nil {
		return fmt.Errorf("latest execution block: %w", err)
	}

	prevTimestamp := head.BlockTime

	var ts uint64

	if bb.lastCallTime.IsZero() {
		// First block, initialize lastCallTime and set default timestamp
		ts = uint64(time.Now().UnixMilli()) + bb.buildDelayMs
		bb.lastCallTime = currentCallTime
	} else {
		// Compute diff in milliseconds
		diff := currentCallTime.Sub(bb.lastCallTime)
		diffMillis := diff.Milliseconds()

		if uint64(diffMillis) <= bb.buildDelayMs {
			ts = prevTimestamp + bb.buildDelayMs
		} else {
			// For every multiple of buildDelay that diff exceeds, increment the block time by that multiple.
			multiples := (uint64(diffMillis) + bb.buildDelayMs - 1) / bb.buildDelayMs // Round up to next multiple of buildDelay
			ts = prevTimestamp + multiples*bb.buildDelayMs
		}

		bb.lastCallTime = currentCallTime
	}

	// Very low chance to happen, only after restart and time.Now is broken
	if ts <= head.BlockTime {
		ts = head.BlockTime + 1 // Subsequent blocks must have a higher timestamp.
	}

	err = retryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		response, err := bb.startBuild(ctx, common.Address{}, head, ts)
		if err != nil {
			bb.logger.Warn("Failed to build new EVM payload, will retry", "error", err)
			return err // Will retry
		} else if response.PayloadStatus.Status != engine.VALID {
			return backoff.Permanent(fmt.Errorf("invalid payload status: %bb", response.PayloadStatus.Status))
		} else if response.PayloadID == nil {
			return backoff.Permanent(errors.New("payloadID is nil"))
		}

		bb.logger.Info("Leader: GetPayload completed", "PayloadID", response.PayloadID.String())

		payloadID = response.PayloadID
		return nil // Success
	})

	if err != nil {
		return fmt.Errorf("failed to start build: %w", err)
	}

	if payloadID == nil {
		return errors.New("payloadID is nil")
	}

	waitTo := time.Now().Add(bb.buildDelay)
	select {
	case <-ctx.Done():
		bb.logger.Info("context cancelled")
		return nil
	case <-time.After(time.Until(waitTo)):
		bb.logger.Info("Leader: Waited for EVM build delay", "delay", bb.buildDelay)
	}

	var payloadResp *engine.ExecutionPayloadEnvelope
	err = retryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		var err error
		payloadResp, err = bb.engineCl.GetPayloadV3(ctx, *payloadID)
		if isUnknownPayload(err) {
			return backoff.Permanent(err)
		} else if err != nil {
			bb.logger.Warn("Failed to get payload, retrying...", "error", err)
			return err // Will retry
		}

		return nil // Success
	})

	if err != nil {
		return fmt.Errorf("failed to get payload: %w", err)
	}

	payloadData, err := msgpack.Marshal(payloadResp.ExecutionPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadIDStr := payloadID.String()

	err = bb.stateManager.SaveBlockStateAndPublishToStream(ctx, &types.BlockBuildState{
		CurrentStep:      types.StepFinalizeBlock,
		PayloadID:        payloadIDStr,
		ExecutionPayload: string(payloadData),
	})
	if err != nil {
		return fmt.Errorf("failed to save state after GetPayload: %w", err)
	}

	bb.logger.Info("Leader: BuildBlock completed and block is distributed", "PayloadID", payloadIDStr)

	return nil
}

func isUnknownPayload(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(err.Error()),
		strings.ToLower(engine.UnknownPayload.Error()),
	)
}

func (bb *BlockBuilder) processLastPayload(ctx context.Context) error {
	bbState := bb.stateManager.GetBlockBuildState(ctx)
	if bbState.ExecutionPayload == "" {
		return nil
	}

	// If execPayload is not empty, the app likely exited after step 1
	bb.logger.Debug("exec payload not nil, processing last payload")

	err := retryWithInfiniteBackoff(ctx, bb.logger, func() error {
		err := bb.finalizeBlock(ctx, bbState.PayloadID, bbState.ExecutionPayload, "")
		if err != nil {
			re := regexp.MustCompile(`invalid block height: (\d+), expected: (\d+)`)
			matches := re.FindStringSubmatch(err.Error())
			// handling edge, if leader fails to reset state, but already pushed his block to the EVM
			if len(matches) == 3 {
				invalidHeight, err1 := strconv.Atoi(matches[1])
				expectedHeight, err2 := strconv.Atoi(matches[2])

				if err1 == nil && err2 == nil {
					if invalidHeight == expectedHeight-1 {
						bb.logger.Warn("Follower: Block already pushed to EVM, resetting state to StepBuildBlock")
						return nil // Success
					} else {
						bb.logger.Warn("Follower: Invalid block height, exit", "invalid_height", invalidHeight, "expected_height", expectedHeight)
						return backoff.Permanent(err)
					}
				} else {
					// Impossible to reach, unless geth changes the error message response
					bb.logger.Warn("Conversion error", "error1", err1, "error2", err2)
					return backoff.Permanent(fmt.Errorf("conversion error1: %w, error2: %w", err1, err2))
				}
			} else {
				bb.logger.Warn("Follower: Failed to finalize block, retrying...", "error", err)
				return err // Will retry
			}
		}
		return nil // Success
	})

	// could happen, only if program exited and ctx cancelled
	if err != nil {
		bb.logger.Error("Follower: Failed to finalize block with retry, exiting")
		return err
	}

	err = retryWithInfiniteBackoff(ctx, bb.logger, func() error {
		bb.logger.Info("Follower: Resetting state to StepBuildBlock for next block")
		err := bb.stateManager.ResetBlockState(ctx)
		if err != nil {
			bb.logger.Warn("Follower: Failed to reset block state, retrying...", "error", err)
			return err // Will retry
		}
		return nil // Success
	})

	// could happen, only if program exited and ctx cancelled
	if err != nil {
		bb.logger.Warn("Follower: Failed to reset block state, exiting", "error", err)
		return err
	}

	return nil
}

func isUnknown(status engine.PayloadStatusV1) bool {
	if status.Status == engine.VALID ||
		status.Status == engine.INVALID ||
		status.Status == engine.SYNCING ||
		status.Status == engine.ACCEPTED {
		return false
	}

	return true
}

func isInvalid(status engine.PayloadStatusV1) (bool, error) {
	if status.Status != engine.INVALID {
		return false, nil
	}

	valErr := "nil"
	if status.ValidationError != nil {
		valErr = *status.ValidationError
	}

	hash := "nil"
	if status.LatestValidHash != nil {
		hash = status.LatestValidHash.Hex()
	}

	return true, fmt.Errorf("payload invalid, validation_err: %s, last_valid_hash: %s", valErr, hash)
}

func isSyncing(status engine.PayloadStatusV1) bool {
	return status.Status == engine.SYNCING || status.Status == engine.ACCEPTED
}

// temp function for testing
func sometimesFails() error {
	rand.Seed(uint64(time.Now().UnixNano()))

	chance := rand.Intn(5) // 0, 1, 2, 3, 4

	if chance == 0 {
		// Fail 1 out of 5 times (when chance == 0)
		return errors.New("failed: 1 in 5 chance")
	}

	// Otherwise succeed
	return nil
}

func (bb *BlockBuilder) finalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	if payloadIDStr == "" || executionPayloadStr == "" {
		return errors.New("PayloadID or ExecutionPayload is missing in build state")
	}

	var executionPayload engine.ExecutableData
	if err := msgpack.Unmarshal([]byte(executionPayloadStr), &executionPayload); err != nil {
		return fmt.Errorf("failed to deserialize ExecutionPayload: %w", err)
	}

	head, err := bb.stateManager.LoadExecutionHead(ctx)
	if err != nil {
		return fmt.Errorf("failed to load execution head: %w", err)
	}

	if err := bb.validateExecutionPayload(executionPayload, head); err != nil {
		return err
	}

	hash := common.BytesToHash(head.BlockHash)
	retryFunc := bb.selectRetryFunction(ctx, msgID)

	if err := bb.pushNewPayload(ctx, executionPayload, hash, retryFunc); err != nil {
		return fmt.Errorf("failed to push new payload: %w", err)
	}

	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	if err := bb.updateForkChoice(ctx, fcs, retryFunc); err != nil {
		return fmt.Errorf("failed to finalize fork choice update: %w", err)
	}

	executionHead := &types.ExecutionHead{
		BlockHeight: executionPayload.Number,
		BlockHash:   executionPayload.BlockHash[:],
		BlockTime:   executionPayload.Timestamp,
	}

	if err := bb.saveExecutionHead(ctx, executionHead, msgID); err != nil {
		return fmt.Errorf("failed to save execution head: %w", err)
	}

	return nil
}

func (bb *BlockBuilder) validateExecutionPayload(executionPayload engine.ExecutableData, head *types.ExecutionHead) error {
	if executionPayload.Number != head.BlockHeight+1 {
		return fmt.Errorf("invalid block height: %d, expected: %d", executionPayload.Number, head.BlockHeight+1)
	}
	if executionPayload.ParentHash != common.Hash(head.BlockHash) {
		return fmt.Errorf("invalid parent hash: %s, head: %s", executionPayload.ParentHash, head.BlockHash)
	}
	minTimestamp := head.BlockTime + 1
	if executionPayload.Timestamp < minTimestamp && executionPayload.Number != 1 {
		return fmt.Errorf("invalid timestamp: %d, min: %d", executionPayload.Timestamp, minTimestamp)
	}
	hash := common.BytesToHash(head.BlockHash)
	if executionPayload.Random != hash {
		return fmt.Errorf("invalid random: %s, head: %s", executionPayload.Random, head.BlockHash)
	}
	return nil
}

func (bb *BlockBuilder) selectRetryFunction(ctx context.Context, msgID string) func(operation func() error) error {
	if msgID == "" {
		return func(operation func() error) error {
			return retryWithBackoff(ctx, maxAttempts, bb.logger, operation)
		}
	}
	return func(operation func() error) error {
		return retryWithInfiniteBackoff(ctx, bb.logger, operation)
	}
}

func (bb *BlockBuilder) pushNewPayload(ctx context.Context, executionPayload engine.ExecutableData, hash common.Hash, retryFunc func(f func() error) error) error {
	emptyVersionHashes := []common.Hash{}
	return retryFunc(func() error {
		status, err := bb.engineCl.NewPayloadV3(ctx, executionPayload, emptyVersionHashes, &hash)
		if err != nil || isUnknown(status) {
			bb.logger.Error("Failed to push new payload", "error", err)
			return err // Will retry
		}
		if invalid, err := isInvalid(status); invalid {
			bb.logger.Error("Payload is not valid", "error", err)
			return backoff.Permanent(fmt.Errorf("payload is not valid: %w", err))
		}
		if isSyncing(status) {
			bb.logger.Info("Processing payload, EVM is syncing")
		}
		return nil // Success
	})
}


func (bb *BlockBuilder) updateForkChoice(ctx context.Context, fcs engine.ForkchoiceStateV1, retryFunc func(f func() error) error) error {
	return retryFunc(func() error {
		fcr, err := bb.engineCl.ForkchoiceUpdatedV3(ctx, fcs, nil)
		if err != nil || isUnknown(fcr.PayloadStatus) {
			bb.logger.Error("Failed to finalize fork choice update", "error", err)
			return err // Will retry
		}
		if invalid, err := isInvalid(fcr.PayloadStatus); invalid {
			bb.logger.Error("Payload is not valid", "error", err)
			return backoff.Permanent(fmt.Errorf("payload is not valid: %w", err))
		}
		if isSyncing(fcr.PayloadStatus) {
			bb.logger.Info("Payload is syncing")
		}
		return nil // Success
	})
}

func (bb *BlockBuilder) saveExecutionHead(ctx context.Context, executionHead *types.ExecutionHead, msgID string) error {
	if msgID == "" {
		return bb.stateManager.SaveExecutionHead(ctx, executionHead)
	}
	return bb.stateManager.SaveExecutionHeadAndAck(ctx, executionHead, msgID)
}
