package blockbuilder

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/cl/types"
	"github.com/primev/mev-commit/cl/util"
	"github.com/vmihailenco/msgpack/v5"
)

const maxAttempts = 10

var ErrEmptyBlock = errors.New("payloadID is empty")

type EngineClient interface {
	NewPayloadV4(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
		beaconRoot *common.Hash, executionRequests []hexutil.Bytes) (engine.PayloadStatusV1, error)
	ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
		payloadAttributes *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)

	GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)

	HeaderByNumber(ctx context.Context, number *big.Int) (*etypes.Header, error)
}

type stateManager interface {
	SaveBlockStateAndPublishToStream(ctx context.Context, state *types.BlockBuildState) error
	GetBlockBuildState(ctx context.Context) types.BlockBuildState
	ResetBlockState(ctx context.Context) error
}

type rpcClient interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type BlockBuilder struct {
	stateManager          stateManager
	engineCl              EngineClient
	rpcClient             rpcClient
	logger                *slog.Logger
	buildDelay            time.Duration
	buildEmptyBlocksDelay time.Duration
	buildDelayMs          uint64
	lastBlockTime         time.Time
	feeRecipient          common.Address
	executionHead         *types.ExecutionHead
}

func NewBlockBuilder(
	stateManager stateManager,
	engineCl EngineClient,
	logger *slog.Logger,
	buildDelay,
	buildDelayEmptyBlocks time.Duration,
	feeReceipt string,
	rpcClient rpcClient,
) *BlockBuilder {
	return &BlockBuilder{
		stateManager:          stateManager,
		engineCl:              engineCl,
		logger:                logger,
		buildDelay:            buildDelay,
		buildDelayMs:          uint64(buildDelay.Milliseconds()),
		buildEmptyBlocksDelay: buildDelayEmptyBlocks,
		feeRecipient:          common.HexToAddress(feeReceipt),
		lastBlockTime:         time.Now().Add(-buildDelayEmptyBlocks),
		rpcClient:             rpcClient,
	}
}

func NewMemberBlockBuilder(engineCL EngineClient, logger *slog.Logger) *BlockBuilder {
	return &BlockBuilder{
		engineCl: engineCL,
		logger:   logger,
	}
}

func (bb *BlockBuilder) SetLastCallTimeToZero() {
	bb.lastBlockTime = time.Time{}
}

func (bb *BlockBuilder) startBuild(ctx context.Context, head *types.ExecutionHead, ts uint64) (engine.ForkChoiceResponse, error) {
	hash := common.BytesToHash(head.BlockHash)

	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	bb.logger.Info("Leader: Submit new EVM payload", "timestamp", ts)

	attrs := &engine.PayloadAttributes{
		Timestamp:             ts,
		Random:                hash,            // We use head block hash as randao.
		SuggestedFeeRecipient: bb.feeRecipient, // Recipient of the priority fee (not base fee)
		Withdrawals:           []*etypes.Withdrawal{},
		BeaconRoot:            &hash,
	}

	resp, err := bb.engineCl.ForkchoiceUpdatedV3(ctx, fcs, attrs)
	if err != nil {
		return engine.ForkChoiceResponse{}, fmt.Errorf("forkchoice update, %w", err)
	}

	return resp, nil
}

func (bb *BlockBuilder) GetPayload(ctx context.Context) error {
	var (
		payloadID *engine.PayloadID
		head      *types.ExecutionHead
		err       error
	)
	currentCallTime := time.Now()

	mempoolStatus, err := bb.GetMempoolStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending transaction count: %w", err)
	}
	bb.logger.Debug("GetMempoolStatus rpc duration", "duration", time.Since(currentCallTime))

	if mempoolStatus.Pending == 0 {
		timeSinceLastBlock := currentCallTime.Sub(bb.lastBlockTime)
		if timeSinceLastBlock < bb.buildEmptyBlocksDelay {
			bb.logger.Debug(
				"Leader: Skipping empty block",
				"timeSinceLastBlock", timeSinceLastBlock,
				"pendingTxes", mempoolStatus.Pending,
				"queuedTxes", mempoolStatus.Queued,
			)
			return ErrEmptyBlock
		}
		bb.logger.Info(
			"Leader: Empty block will be created",
			"timeSinceLastBlock", timeSinceLastBlock,
			"pendingTxes", mempoolStatus.Pending,
			"queuedTxes", mempoolStatus.Queued,
		)
	}

	// Load execution head to get previous block timestamp
	err = util.RetryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		head, err = bb.loadExecutionHead(ctx)
		if err != nil {
			bb.logger.Warn(
				"Failed to load execution head, retrying...",
				"error", err,
			)
			return err // Will retry
		}
		return nil // Success
	})
	if err != nil {
		return fmt.Errorf("latest execution block: %w", err)
	}

	prevTimestamp := head.BlockTime

	var ts uint64

	if bb.lastBlockTime.IsZero() {
		// First block, initialize LastCallTime and set default timestamp
		ts = uint64(time.Now().UnixMilli()) + bb.buildDelayMs
	} else {
		// Compute diff in milliseconds
		diff := currentCallTime.Sub(bb.lastBlockTime)
		diffMillis := diff.Milliseconds()

		if uint64(diffMillis) <= bb.buildDelayMs {
			ts = prevTimestamp + bb.buildDelayMs
		} else {
			// For every multiple of buildDelay that diff exceeds, increment the block time by that multiple.
			multiples := (uint64(diffMillis) + bb.buildDelayMs - 1) / bb.buildDelayMs // Round up to next multiple of buildDelay
			ts = prevTimestamp + multiples*bb.buildDelayMs
		}
	}

	// Very low chance to happen, only after restart and time.Now is broken
	if ts <= head.BlockTime {
		ts = head.BlockTime + 1 // Subsequent blocks must have a higher timestamp.
	}

	err = util.RetryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		response, err := bb.startBuild(ctx, head, ts)
		if err != nil {
			bb.logger.Warn(
				"Failed to build new EVM payload, will retry",
				"error", err,
			)
			return err // Will retry
		} else if response.PayloadStatus.Status != engine.VALID {
			return backoff.Permanent(fmt.Errorf("invalid payload status: %s", response.PayloadStatus.Status))
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
	err = util.RetryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		var err error
		payloadResp, err = bb.engineCl.GetPayloadV4(ctx, *payloadID)
		if isUnknownPayload(err) {
			return backoff.Permanent(err)
		} else if err != nil {
			bb.logger.Warn(
				"Failed to get payload, retrying...",
				"error", err,
			)
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

	encodedPayload := base64.StdEncoding.EncodeToString(payloadData)

	payloadIDStr := payloadID.String()

	err = bb.stateManager.SaveBlockStateAndPublishToStream(ctx, &types.BlockBuildState{
		CurrentStep:      types.StepFinalizeBlock,
		PayloadID:        payloadIDStr,
		ExecutionPayload: encodedPayload,
	})
	if err != nil {
		return fmt.Errorf("failed to save state after GetPayload: %w", err)
	}

	bb.logger.Info(
		"Leader: BuildBlock completed and block is distributed",
		"PayloadID", payloadIDStr,
	)

	bb.lastBlockTime = time.Now()
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

func (bb *BlockBuilder) ProcessLastPayload(ctx context.Context) error {
	bbState := bb.stateManager.GetBlockBuildState(ctx)
	if bbState.ExecutionPayload == "" {
		return nil
	}

	// If execPayload is not empty, the app likely exited after step 1
	bb.logger.Debug("exec payload not nil, processing last payload")

	err := util.RetryWithInfiniteBackoff(ctx, bb.logger, func() error {
		err := bb.FinalizeBlock(ctx, bbState.PayloadID, bbState.ExecutionPayload, "")
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
						bb.logger.Warn(
							"Follower: Invalid block height, exit",
							"invalid_height", invalidHeight,
							"expected_height", expectedHeight,
						)
						return backoff.Permanent(err)
					}
				} else {
					// Impossible to reach, unless geth changes the error message response
					bb.logger.Warn(
						"Conversion error",
						"error1", err1,
						"error2", err2,
					)
					return backoff.Permanent(fmt.Errorf("conversion error1: %w, error2: %w", err1, err2))
				}
			} else {
				bb.logger.Warn(
					"Follower: Failed to finalize block, retrying...",
					"error", err,
				)
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

	err = util.RetryWithInfiniteBackoff(ctx, bb.logger, func() error {
		bb.logger.Info("Follower: Resetting state to StepBuildBlock for next block")
		err := bb.stateManager.ResetBlockState(ctx)
		if err != nil {
			bb.logger.Warn(
				"Follower: Failed to reset block state, retrying...",
				"error", err,
			)
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

func (bb *BlockBuilder) FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	if payloadIDStr == "" || executionPayloadStr == "" {
		return errors.New("PayloadID or ExecutionPayload is missing in build state")
	}

	executionPayloadBytes, err := base64.StdEncoding.DecodeString(executionPayloadStr)
	if err != nil {
		return fmt.Errorf("failed to decode ExecutionPayload: %w", err)
	}

	var executionPayload engine.ExecutableData
	if err := msgpack.Unmarshal(executionPayloadBytes, &executionPayload); err != nil {
		return fmt.Errorf("failed to deserialize ExecutionPayload: %w", err)
	}

	var head *types.ExecutionHead
	err = util.RetryWithBackoff(ctx, maxAttempts, bb.logger, func() error {
		head, err = bb.loadExecutionHead(ctx)
		if err != nil {
			bb.logger.Warn(
				"Failed to load execution head, retrying...",
				"error", err,
			)
			return err // Will retry
		}
		return nil // Success
	})
	if err != nil {
		return fmt.Errorf("failed to load execution head: %w", err)
	}

	if err := bb.validateExecutionPayload(executionPayload, head); err != nil {
		return fmt.Errorf("failed to validate execution payload: %w", err)
	}

	hash := common.BytesToHash(head.BlockHash)
	retryFunc := bb.selectRetryFunction(ctx, msgID)

	if err := bb.pushNewPayload(ctx, executionPayload, hash, retryFunc); err != nil {
		return fmt.Errorf("failed to push new payload: %w", err)
	}

	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      executionPayload.BlockHash,
		SafeBlockHash:      executionPayload.BlockHash,
		FinalizedBlockHash: executionPayload.BlockHash,
	}

	if err := bb.updateForkChoice(ctx, fcs, retryFunc); err != nil {
		return fmt.Errorf("failed to finalize fork choice update: %w", err)
	}

	bb.executionHead = &types.ExecutionHead{
		BlockHeight: executionPayload.Number,
		BlockHash:   executionPayload.BlockHash[:],
		BlockTime:   executionPayload.Timestamp,
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
			return util.RetryWithBackoff(ctx, maxAttempts, bb.logger, operation)
		}
	}
	return func(operation func() error) error {
		return util.RetryWithInfiniteBackoff(ctx, bb.logger, operation)
	}
}

func (bb *BlockBuilder) pushNewPayload(ctx context.Context, executionPayload engine.ExecutableData, hash common.Hash, retryFunc func(f func() error) error) error {
	emptyVersionHashes := []common.Hash{}
	return retryFunc(func() error {
		status, err := bb.engineCl.NewPayloadV4(ctx, executionPayload, emptyVersionHashes, &hash, []hexutil.Bytes{})
		bb.logger.Debug("newPayload result",
			"status", status.Status,
			"validationError", status.ValidationError,
			"latestValidHash", status.LatestValidHash)
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

func (bb *BlockBuilder) loadExecutionHead(ctx context.Context) (*types.ExecutionHead, error) {
	if bb.executionHead != nil {
		return bb.executionHead, nil
	}

	header, err := bb.engineCl.HeaderByNumber(ctx, nil) // nil for the latest block
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block header: %w", err)
	}

	bb.executionHead = &types.ExecutionHead{
		BlockHeight: header.Number.Uint64(),
		BlockHash:   header.Hash().Bytes(),
		BlockTime:   header.Time,
	}

	return bb.executionHead, nil
}

func (bb *BlockBuilder) GetExecutionHead() *types.ExecutionHead {
	return bb.executionHead
}

type MempoolStatus struct {
	Pending hexutil.Uint64 `json:"pending"`
	Queued  hexutil.Uint64 `json:"queued"`
}

func (bb *BlockBuilder) GetMempoolStatus(ctx context.Context) (*MempoolStatus, error) {
	var result MempoolStatus
	err := bb.rpcClient.CallContext(ctx, &result, "txpool_status")
	if err != nil {
		return nil, err
	}
	return &result, nil
}
