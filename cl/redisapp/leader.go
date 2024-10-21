package redisapp

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/heyvito/go-leader/leader"
	"github.com/primev/mev-commit/cl/redisapp/types"
)

type Leader struct {
	InstanceID     string
	stateManager   StateManager
	stepsManager   *StepsManager
	ctx            context.Context
	cancel         context.CancelFunc
	leaderElection leader.Leader
	wg             *sync.WaitGroup
	logger         Logger
}

func (l *Leader) startLeaderLoop() {
	l.logger.Info("Starting leader loop")
	leaderCtx, leaderCancel := context.WithCancel(l.ctx)
	l.cancel = leaderCancel

	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.leaderLoop(leaderCtx)
	}()
}

func (l *Leader) stopLeaderLoop() {
	if l.cancel != nil {
		l.cancel()
		l.cancel = nil
		l.logger.Info("Leader loop stopped")
	}
}

func (l *Leader) leaderLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.logger.Info("Leader loop exiting")
			return
		default:
			bbState := l.stateManager.GetBlockBuildState(ctx)
			// Determine which step to execute next based on the current state
			currentStep := bbState.CurrentStep

			switch currentStep {
			case types.StepBuildBlock:
				l.logger.Info("Leader: Starting Step 1 - BuildBlock")
				getPayloadErr := l.stepsManager.getPayload(ctx)
				if getPayloadErr != nil {
					l.logger.Error("Leader: Failed to execute Step 1 - BuildBlock", "error", getPayloadErr)

					err := l.stateManager.ResetBlockState(ctx)
					if err != nil {
						l.logger.Error("Leader: Failed to reset leader state", "error", err)
					}

					if errors.Is(getPayloadErr, ErrFailedAfterNAttempts) {
						l.logger.Error("Leader: failed to reach geth node after max attempts, exiting")
						err := l.leaderElection.Stop()
						if err != nil {
							l.logger.Error("Leader: Failed to stop leader election", "error", err)
						}
					}

					time.Sleep(2 * time.Second)
					continue
				}

			case types.StepFinalizeBlock:
				l.logger.Info("Leader: Starting Step 2 - FinalizeBlock")
				err := l.stepsManager.finalizeBlock(ctx, bbState.PayloadID, bbState.ExecutionPayload, "")
				if err != nil {
					l.logger.Error("Leader: Failed to execute Step 2 - FinalizeBlock", "error", err)

					if errors.Is(err, ErrFailedAfterNAttempts) {
						l.logger.Error("Leader: failed to reach geth node after max attempts, exiting")
						err := l.leaderElection.Stop()
						if err != nil {
							l.logger.Error("Leader: Failed to stop leader election", "error", err)
						}
					}

					time.Sleep(2 * time.Second)
					continue
				}
				l.logger.Info("Leader: Resetting state to StepBuildBlock for next block")
				err = l.stateManager.ResetBlockState(ctx)
				if err != nil {
					l.logger.Error("Leader: Failed to reset leader state", "error", err)
				}

			default:
				l.logger.Warn("Leader: Unknown current step", "current_step", currentStep.String())
				err := l.stateManager.ResetBlockState(ctx)
				if err != nil {
					l.logger.Error("Leader: Failed to reset leader state", "error", err)
				}
			}
		}
	}
}
