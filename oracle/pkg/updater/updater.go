package updater

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	preconf "github.com/primevprotocol/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-commit/oracle/pkg/settler"
	"github.com/prometheus/client_golang/prometheus"
)

type BlockWinner struct {
	BlockNumber int64
	Winner      string
}

type WinnerRegister interface {
	SubscribeWinners(ctx context.Context) <-chan BlockWinner
	UpdateComplete(ctx context.Context, blockNum int64) error
	AddSettlement(
		ctx context.Context,
		commitmentIdx []byte,
		txHash string,
		blockNum int64,
		amount uint64,
		builder string,
		bidID []byte,
		settlementType settler.SettlementType,
		decayPercentage int64,
	) error
}

type EVMClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type Oracle interface {
	GetBuilder(builder string) (common.Address, error)
}

type Preconf interface {
	GetCommitmentsByBlockNumber(blockNum *big.Int) ([][32]byte, error)
	GetCommitment(commitmentIdx [32]byte) (preconf.PreConfCommitmentStorePreConfCommitment, error)
}

type Updater struct {
	logger               *slog.Logger
	l1Client             EVMClient
	l2Client             EVMClient
	winnerRegister       WinnerRegister
	preconfClient        Preconf
	rollupClient         Oracle
	builderIdentityCache map[string]common.Address
	metrics              *metrics
}

func NewUpdater(
	logger *slog.Logger,
	l1Client EVMClient,
	l2Client EVMClient,
	winnerRegister WinnerRegister,
	rollupClient Oracle,
	preconfClient Preconf,
) *Updater {
	return &Updater{
		logger:               logger,
		l1Client:             l1Client,
		l2Client:             l2Client,
		winnerRegister:       winnerRegister,
		preconfClient:        preconfClient,
		rollupClient:         rollupClient,
		builderIdentityCache: make(map[string]common.Address),
		metrics:              newMetrics(),
	}
}

func (u *Updater) Metrics() []prometheus.Collector {
	return u.metrics.Collectors()
}

func (u *Updater) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

	RESTART:
		cctx, unsub := context.WithCancel(ctx)
		winnerChan := u.winnerRegister.SubscribeWinners(cctx)

		for {
			select {
			case <-ctx.Done():
				unsub()
				return
			case winner, more := <-winnerChan:
				if !more {
					unsub()
					goto RESTART
				}
				u.metrics.UpdaterTriggerCount.Inc()

				err := func() error {
					var err error
					builderAddr, ok := u.builderIdentityCache[winner.Winner]
					if !ok {
						builderAddr, err = u.rollupClient.GetBuilder(winner.Winner)
						if err != nil {
							if errors.Is(err, ethereum.NotFound) {
								u.logger.Warn("builder not registered", "builder", winner.Winner)
								return u.winnerRegister.UpdateComplete(ctx, winner.BlockNumber)
							}
							return fmt.Errorf("failed to get builder address: %w", err)
						}
						u.builderIdentityCache[winner.Winner] = builderAddr
					}

					blk, err := u.l1Client.BlockByNumber(ctx, big.NewInt(winner.BlockNumber))
					if err != nil {
						return fmt.Errorf("failed to get block by number: %w", err)
					}

					txnsInBlock := make(map[string]int)
					for posInBlock, tx := range blk.Transactions() {
						txnsInBlock[strings.TrimPrefix(tx.Hash().Hex(), "0x")] = posInBlock
					}

					commitmentIndexes, err := u.preconfClient.GetCommitmentsByBlockNumber(
						big.NewInt(winner.BlockNumber),
					)
					if err != nil {
						return fmt.Errorf("failed to get commitments by block number: %w", err)
					}

					u.logger.Debug(
						"commitment indexes",
						"commitments_count", len(commitmentIndexes),
						"txns_count", len(txnsInBlock),
						"blockNumber", winner.BlockNumber,
					)

					total, rewards, slashes := 0, 0, 0
					for _, index := range commitmentIndexes {
						commitment, err := u.preconfClient.GetCommitment(index)
						if err != nil {
							return fmt.Errorf("failed to get commitment: %w", err)
						}

						l2Block, err := u.l2Client.BlockByNumber(ctx, commitment.BlockCommitedAt)
						if err != nil {
							return fmt.Errorf("failed to get L2 Block: %w", err)
						}
						decayPercentage := computeDecayPercentage(commitment.DecayStartTimeStamp, commitment.DecayEndTimeStamp, l2Block.Header().Time)

						settlementType := settler.SettlementTypeReturn

						if commitment.Commiter.Cmp(builderAddr) == 0 {
							commitmentTxnHashes := strings.Split(commitment.TxnHash, ",")
							settlementType = settler.SettlementTypeReward

							// Ensure Bundle is atomic and present in the block
							for i := 0; i < len(commitmentTxnHashes); i++ {
								posInBlock, found := txnsInBlock[commitmentTxnHashes[i]]
								if !found || posInBlock != txnsInBlock[commitmentTxnHashes[0]]+i {
									settlementType = settler.SettlementTypeSlash
									break
								}
							}
						}

						err = u.winnerRegister.AddSettlement(
							ctx,
							index[:],
							commitment.TxnHash,
							winner.BlockNumber,
							commitment.Bid,
							commitment.Commiter.Hex(),
							commitment.CommitmentHash[:],
							settlementType,
							decayPercentage,
						)
						if err != nil {
							return fmt.Errorf("failed to add settlement: %w", err)
						}

						total++
						switch settlementType {
						case settler.SettlementTypeSlash:
							slashes++
						case settler.SettlementTypeReward:
							rewards++
						}
					}

					err = u.winnerRegister.UpdateComplete(ctx, winner.BlockNumber)
					if err != nil {
						return fmt.Errorf("failed to update completion of block updates: %w", err)
					}

					u.metrics.CommitmentsCount.Add(float64(total))
					u.metrics.RewardsCount.Add(float64(rewards))
					u.metrics.SlashesCount.Add(float64(slashes))
					u.metrics.BlockCommitmentsCount.Inc()

					u.logger.Info(
						"added settlements",
						"total", total,
						"rewards", rewards,
						"slashes", slashes,
						"blockNumber", winner.BlockNumber,
						"winner", winner.Winner,
					)

					return nil
				}()

				if err != nil {
					u.logger.Error("failed to process settlements", "blockNumber", winner.BlockNumber, "winner", winner.Winner, "error", err)
					unsub()
					goto RESTART
				}
			}
		}
	}()

	return doneChan
}

// computeDecayPercentage takes startTimestamp, endTimestamp, commitTimestamp and computes a linear decay percentage
// The computation does not care what format the timestamps are in, as long as they are consistent
// (e.g they could be unix or unixMili timestamps)
func computeDecayPercentage(startTimestamp, endTimestamp, commitTimestamp uint64) int64 {
	if startTimestamp >= endTimestamp || startTimestamp > commitTimestamp {
		return 0
	}

	// Calculate the total time in seconds
	totalTime := endTimestamp - startTimestamp
	// Calculate the time passed in seconds
	timePassed := commitTimestamp - startTimestamp
	// Calculate the decay percentage
	decayPercentage := float64(timePassed) / float64(totalTime)

	return int64(math.Round(decayPercentage * 100))
}
