package preconf

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/armon/go-radix"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

const (
	noOfBids = 20
)

var (
	bidKey = func(txHashes []string) string {
		return fmt.Sprintf("bid/%s", strings.Join(txHashes, ","))
	}

	encryptCmtKey = func(cmtDigest []byte) string {
		return fmt.Sprintf("encr/%s", string(cmtDigest))
	}

	openCmtKey = func(cmtIdx []byte) string {
		return fmt.Sprintf("opened/%s", string(cmtIdx))
	}

	settleKey = func(cmtIdx []byte) string {
		return fmt.Sprintf("settle/%s", cmtIdx)
	}

	fundsUnlockedKey = func(cmtDigest []byte) string {
		return fmt.Sprintf("fu/%s", string(cmtDigest))
	}

	fundsRewardedKey = func(cmtDigest []byte) string {
		return fmt.Sprintf("frw/%s", string(cmtDigest))
	}

	fundsSlashedKey = func(providerAddr common.Address, amount *big.Int) string {
		return fmt.Sprintf("fs/%s/%s", providerAddr, amount)
	}

	blkKey = func(bNo uint64) string {
		return fmt.Sprintf("blk/%d", bNo)
	}

	blkWinnerKey = func(bNo uint64) string {
		return fmt.Sprintf("blkw/%d", bNo)
	}

	errNoTxnsInBlock = fmt.Errorf("no transactions in block")
)

type BidEntry struct {
	Bid         *bidderapiv1.Bid
	Accept      bool
	ShouldSlash bool
	Preconfs    []*bidderapiv1.Commitment
}

func RunPreconf(ctx context.Context, cluster orchestrator.Orchestrator, _ any) error {
	bidders := cluster.Bidders()
	logger := cluster.Logger().With("test", "preconfirmations")
	store := radix.New()

	// Listen for encrypted commitments, opened commitments, and settlements
	sub, err := cluster.Events().Subscribe(
		events.NewEventHandler(
			"UnopenedCommitmentStored",
			func(c *preconf.PreconfmanagerUnopenedCommitmentStored) {
				logger.Info("Received encrypted commitment", "digest", hex.EncodeToString(c.CommitmentDigest[:]))
				store.Insert(encryptCmtKey(c.CommitmentDigest[:]), c)
			},
		),
		events.NewEventHandler(
			"OpenedCommitmentStored",
			func(c *preconf.PreconfmanagerOpenedCommitmentStored) {
				logger.Info(
					"Received opened commitment",
					"digest", hex.EncodeToString(c.CommitmentDigest[:]),
					"index", hex.EncodeToString(c.CommitmentIndex[:]),
					"decay_start", c.DecayStartTimeStamp,
					"decay_end", c.DecayEndTimeStamp,
					"dispatch_timestamp", c.DispatchTimestamp,
					"block_number", c.BlockNumber,
				)
				store.Insert(openCmtKey(c.CommitmentIndex[:]), c)
			},
		),
		events.NewEventHandler(
			"CommitmentProcessed",
			func(c *oracle.OracleCommitmentProcessed) {
				logger.Info(
					"Received settlement",
					"index", hex.EncodeToString(c.CommitmentIndex[:]),
					"slash", c.IsSlash,
				)
				store.Insert(settleKey(c.CommitmentIndex[:]), c)
			},
		),
		events.NewEventHandler(
			"FundsUnlocked",
			func(c *bidderregistry.BidderregistryFundsUnlocked) {
				logger.Info("Unlocked funds", "digest", hex.EncodeToString(c.CommitmentDigest[:]))
				store.Insert(fundsUnlockedKey(c.CommitmentDigest[:]), c)
			},
		),
		events.NewEventHandler(
			"FundsRewarded",
			func(c *bidderregistry.BidderregistryFundsRewarded) {
				logger.Info("Rewarded funds", "digest", hex.EncodeToString(c.CommitmentDigest[:]))
				store.Insert(fundsRewardedKey(c.CommitmentDigest[:]), c)
			},
		),
		events.NewEventHandler(
			"FundsSlashed",
			func(c *providerregistry.ProviderregistryFundsSlashed) {
				logger.Info("Funds slashed", "provider", c.Provider, "amount", c.Amount)
				store.Insert(fundsSlashedKey(c.Provider, c.Amount), c)
			},
		),
		events.NewEventHandler(
			"NewL1Block",
			func(c *blocktracker.BlocktrackerNewL1Block) {
				logger.Info("Received new L1 block", "block", c.BlockNumber)
				store.Insert(blkWinnerKey(c.BlockNumber.Uint64()), c)
			},
		),
	)
	if err != nil {
		return err
	}

	defer sub.Unsubscribe()

	eg := errgroup.Group{}
	egCtx, egCancel := context.WithCancel(ctx)

	for _, p := range cluster.Providers() {
		eg.Go(func() error {
			in, err := p.ProviderAPI().ReceiveBids(egCtx, &providerapiv1.EmptyMessage{})
			if err != nil {
				return err
			}

			out, err := p.ProviderAPI().SendProcessedBids(egCtx)
			if err != nil {
				return err
			}

			for {
				select {
				case <-egCtx.Done():
					return nil
				default:
					bid, err := in.Recv()
					if err == io.EOF {
						return nil
					}
					if err != nil {
						return err
					}
					val, ok := store.Get(bidKey(bid.TxHashes))
					if !ok {
						logger.Error("Bid not found in store", "digest", bid.TxHashes)
						return fmt.Errorf("bid not found in store")
					}
					entry := val.(*BidEntry)
					if len(entry.Bid.RawTransactions) != len(bid.RawTransactions) {
						logger.Error(
							"Raw transactions length mismatch",
							"entry", entry,
							"bid", bid,
						)
						return fmt.Errorf("raw transactions length mismatch")
					}
					if len(entry.Bid.RevertingTxHashes) != len(bid.RevertingTxHashes) {
						logger.Error(
							"Reverting transactions length mismatch",
							"entry", entry,
							"bid", bid,
						)
						return fmt.Errorf("reverting transactions length mismatch")
					}
					if entry.Accept {
						logger.Info("Bid accepted", "entry", entry)
						err := out.Send(&providerapiv1.BidResponse{
							BidDigest:         bid.BidDigest,
							Status:            providerapiv1.BidResponse_STATUS_ACCEPTED,
							DispatchTimestamp: time.Now().UnixMilli() + 100,
						})
						if err != nil {
							logger.Error("Failed to send bid response", "digest", bid.BidDigest)
							return err
						}
					} else {
						logger.Info("Bid rejected", "entry", entry)
						err := out.Send(&providerapiv1.BidResponse{
							BidDigest: bid.BidDigest,
							Status:    providerapiv1.BidResponse_STATUS_REJECTED,
						})
						if err != nil {
							logger.Error("Failed to send bid response", "digest", bid.BidDigest)
							return err
						}
					}
				}
			}
		})
	}

	bidderIn := make(map[string]chan *BidEntry)
	for _, b := range bidders {
		bidderIn[b.EthAddress()] = make(chan *BidEntry)
		eg.Go(func() error {
			for {
				select {
				case <-egCtx.Done():
					return nil
				case entry := <-bidderIn[b.EthAddress()]:
					out, err := b.BidderAPI().SendBid(egCtx, entry.Bid)
					if err != nil {
						return err
					}
					var preconfs []*bidderapiv1.Commitment
					for {
						resp, err := out.Recv()
						if err == io.EOF {
							break
						}
						if err != nil {
							return err
						}
						preconfs = append(preconfs, resp)
					}
					entry.Preconfs = preconfs
					logger.Info("Received preconfs", "count", len(preconfs))
				}
			}
		})
	}

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	count := 0
	lastWinnerBlock := 0

	usedTxHashes := make(map[string]struct{})

DONE:
	for {
		select {
		case <-egCtx.Done():
			egCancel()
			return nil
		case <-tick.C:
			if count == noOfBids {
				_, ok := store.Get(blkWinnerKey(uint64(lastWinnerBlock + 3)))
				if ok {
					// allow enough time for everything to settle
					egCancel()
					logger.Info("All bids sent")
					break DONE
				}
			} else {
				for _, b := range bidders {
					entry, err := getRandomBid(ctx, cluster, store, usedTxHashes)
					if err != nil {
						if errors.Is(err, errNoTxnsInBlock) {
							logger.Info("No transactions in block")
							continue
						}
						egCancel()
						return err
					}
					bidderIn[b.EthAddress()] <- entry
					count++
					lastWinnerBlock = int(entry.Bid.BlockNumber)
					logger.Info("Sent bid", "count", count)
				}
			}
		}
	}

	if err := eg.Wait(); err != nil && !strings.Contains(err.Error(), "context canceled") {
		return err
	}

	bids := make([]*BidEntry, 0, noOfBids)
	store.WalkPrefix("bid/", func(k string, v interface{}) bool {
		bids = append(bids, v.(*BidEntry))
		return false
	})

	for _, entry := range bids {
		if !entry.Accept {
			if len(entry.Preconfs) != 0 {
				logger.Error("Bid not accepted but preconfs received", "entry", entry)
				return fmt.Errorf("bid not accepted but preconfs received")
			}
			continue
		}
		if len(entry.Preconfs) != len(cluster.Providers()) {
			logger.Error(
				"Bid accepted but not all preconfs received",
				"entry", entry,
				"preconfs", len(entry.Preconfs),
				"providers", len(cluster.Providers()),
			)
			return fmt.Errorf("bid accepted but not all preconfs received")
		}
		winner, ok := store.Get(blkWinnerKey(uint64(entry.Bid.BlockNumber)))
		if !ok {
			logger.Error("Winner not found", "block", entry.Bid.BlockNumber)
			return fmt.Errorf("winner not found")
		}
		foundCmt := false
		for _, pc := range entry.Preconfs {
			cmtDigest, err := hex.DecodeString(pc.CommitmentDigest)
			if err != nil {
				logger.Error(
					"Failed to decode commitment digest",
					"error", err,
					"entry", entry,
					"digest", pc.CommitmentDigest,
				)
				return fmt.Errorf("failed to decode commitment digest")
			}
			ec, ok := store.Get(encryptCmtKey(cmtDigest))
			if !ok {
				logger.Error(
					"Encrypted commitment not found",
					"entry", entry,
					"digest", pc.CommitmentDigest,
				)
				return fmt.Errorf("encrypted commitment not found")
			}
			providerAddr, err := hex.DecodeString(pc.ProviderAddress)
			if err != nil {
				logger.Error(
					"Failed to decode provider address",
					"error", err,
					"entry", entry,
					"address", pc.ProviderAddress,
				)
				return fmt.Errorf("failed to decode provider address")
			}
			if common.BytesToAddress(providerAddr).Cmp(winner.(*blocktracker.BlocktrackerNewL1Block).Winner) == 0 {
				foundCmt = true
				ecmt := ec.(*preconf.PreconfmanagerUnopenedCommitmentStored)
				_, ok := store.Get(openCmtKey(ecmt.CommitmentIndex[:]))
				if !ok {
					logger.Error(
						"Opened commitment not found",
						"entry", entry,
						"index", hex.EncodeToString(ecmt.CommitmentIndex[:]),
					)
					return fmt.Errorf("opened commitment not found")
				}
				pcmt, ok := store.Get(settleKey(ecmt.CommitmentIndex[:]))
				if !ok {
					logger.Error(
						"Settlement not found",
						"entry", entry,
						"index", hex.EncodeToString(ecmt.CommitmentIndex[:]),
					)
					return fmt.Errorf("settlement not found")
				}
				residualBidPercent := computeResidualAfterDecay(uint64(pc.DecayStartTimestamp), uint64(pc.DecayEndTimestamp), uint64(pc.DispatchTimestamp))
				bidAmt, _ := new(big.Int).SetString(entry.Bid.Amount, 10)
				slashAmount, _ := new(big.Int).SetString(entry.Bid.SlashAmount, 10)
				residualBidAmt := new(big.Int).Mul(residualBidPercent, bidAmt)
				residualBidAmt.Div(residualBidAmt, big.NewInt(ONE_HUNDRED_PERCENT))

				if entry.ShouldSlash {
					if !pcmt.(*oracle.OracleCommitmentProcessed).IsSlash {
						logger.Error("Provider should be slashed", "entry", entry)
						return fmt.Errorf("provider should be slashed")
					}
					_, ok := store.Get(fundsUnlockedKey(cmtDigest))
					if !ok {
						logger.Error("Funds not unlocked", "entry", entry)
						return fmt.Errorf("funds not unlocked")
					}

					penaltyFee := new(big.Int).Mul(slashAmount, big.NewInt(FEE_PERCENT))
					penaltyFee.Div(penaltyFee, big.NewInt(ONE_HUNDRED_PERCENT))
					totalSlash := new(big.Int).Add(slashAmount, penaltyFee)

					_, ok = store.Get(fundsSlashedKey(common.BytesToAddress(providerAddr), totalSlash))
					if !ok {
						logger.Error("Funds not slashed", "entry", entry, "total", totalSlash)
						return fmt.Errorf("funds not slashed")
					}
				} else {
					if pcmt.(*oracle.OracleCommitmentProcessed).IsSlash {
						logger.Error("Provider should not be slashed", "entry", entry)
						return fmt.Errorf("provider should not be slashed")
					}
					fr, ok := store.Get(fundsRewardedKey(cmtDigest))
					if !ok {
						logger.Error("Funds not rewarded", "entry", entry)
						return fmt.Errorf("funds not rewarded")
					}
					if fr.(*bidderregistry.BidderregistryFundsRewarded).Amount.Cmp(residualBidAmt) != 0 {
						logger.Error("residual bid amount mismatch", "entry", entry, "expected", residualBidAmt, "actual", fr.(*bidderregistry.BidderregistryFundsRewarded).Amount)
						return fmt.Errorf("residual bid amount mismatch")
					}
				}
			}
		}
		if !foundCmt {
			logger.Error(
				"Winner not found in preconfs",
				"entry", entry,
				"winner", winner.(*blocktracker.BlocktrackerNewL1Block).Winner.Hex(),
			)
			return fmt.Errorf("winner not found in preconfs")
		}
	}

	return nil
}

func getRandomBid(
	ctx context.Context,
	o orchestrator.Orchestrator,
	store *radix.Tree,
	usedTxHashes map[string]struct{},
) (*BidEntry, error) {
	blkNum, err := o.L1Client().BlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	blk, found := store.Get(blkKey(blkNum))
	if !found {
		blk, err = o.L1Client().BlockByNumber(ctx, big.NewInt(int64(blkNum)))
		if err != nil {
			return nil, err
		}
		store.Insert(blkKey(blkNum), blk)
	}

	transactions := blk.(*types.Block).Transactions()
	txCount := len(transactions)
	start := rand.Intn(txCount)

	switch txCount {
	case 0:
		return nil, errNoTxnsInBlock
	case 1:
		// skip
	default:
		// we select a random number of transactions to bundle starting from the start index
		// in that order
		maxBundleLen := min(4, txCount-start)
		if maxBundleLen == 1 {
			transactions = transactions[start : start+1]
		} else {
			end := start + rand.Intn(maxBundleLen) + 1
			transactions = transactions[start:end]
		}
	}

	var (
		txHashes []string
		rawTxns  []string
	)
	for _, txn := range transactions {
		txHash := strings.TrimPrefix(txn.Hash().String(), "0x")
		if _, exists := usedTxHashes[txHash]; exists {
			// Duplicate found, try getting a new bid
			return getRandomBid(ctx, o, store, usedTxHashes)
		}
		txHashes = append(
			txHashes,
			txHash,
		)
		buf, err := txn.MarshalBinary()
		if err != nil {
			return nil, err
		}
		usedTxHashes[txHash] = struct{}{}
		rawTxns = append(rawTxns, hex.EncodeToString(buf))
	}

	revertingTxnHashes, err := getRevertingTxns(
		ctx,
		o.L1Client(),
		transactions,
	)
	if err != nil {
		return nil, err
	}

	// send payload instead of hashes
	sendPayload := rand.Intn(100) < 30
	// accept 90% of bids
	accept := rand.Intn(100) < 90
	// slash 10% of accepted bids
	shouldSlash := rand.Intn(100) < 10 && !sendPayload
	// amount between 5M and 6M
	amount := 5_000_000 + rand.Intn(1_000_000)
	// slash amount between 10000 and 100000
	slashAmount := 10_000 + rand.Intn(100_000)

	var opts *bidderapiv1.BidOptions

	if shouldSlash {
		usingOpts := rand.Intn(2) == 0
		if usingOpts && (start > 1 || start+len(txHashes) < txCount-1) {
			if start > 1 {
				opts = &bidderapiv1.BidOptions{
					Options: []*bidderapiv1.BidOption{
						{
							Opt: &bidderapiv1.BidOption_PositionConstraint{
								PositionConstraint: &bidderapiv1.PositionConstraint{
									Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
									Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
									Value:  1,
								},
							},
						},
					},
				}
			} else {
				opts = &bidderapiv1.BidOptions{
					Options: []*bidderapiv1.BidOption{
						{
							Opt: &bidderapiv1.BidOption_PositionConstraint{
								PositionConstraint: &bidderapiv1.PositionConstraint{
									Anchor: bidderapiv1.PositionConstraint_ANCHOR_BOTTOM,
									Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
									Value:  1,
								},
							},
						},
					},
				}
			}
		} else {
			if len(txHashes) > 1 {
				original := slices.Clone(txHashes)
				for {
					rand.Shuffle(len(txHashes), func(i, j int) {
						txHashes[i], txHashes[j] = txHashes[j], txHashes[i]
					})
					if !reflect.DeepEqual(original, txHashes) {
						break
					}
				}
			} else {
				// get random tx hash
				randBytes := make([]byte, 32)
				_, _ = crand.Read(randBytes)
				txHashes[0] = strings.TrimPrefix(common.BytesToHash(randBytes).String(), "0x")
			}
		}
	}

	if !shouldSlash {
		useOpts := rand.Intn(2) == 0
		if useOpts {
			opts = &bidderapiv1.BidOptions{
				Options: []*bidderapiv1.BidOption{
					{
						Opt: &bidderapiv1.BidOption_PositionConstraint{
							PositionConstraint: &bidderapiv1.PositionConstraint{
								Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
								Basis:  bidderapiv1.PositionConstraint_BASIS_ABSOLUTE,
								Value:  int32(start),
							},
						},
					},
				},
			}
		}
	}

	bid := &BidEntry{
		Bid: &bidderapiv1.Bid{
			Amount:              fmt.Sprintf("%d", amount),
			SlashAmount:         fmt.Sprintf("%d", slashAmount),
			BlockNumber:         int64(blkNum),
			DecayStartTimestamp: time.Now().UnixMilli(),
			DecayEndTimestamp:   time.Now().Add(5 * time.Second).UnixMilli(),
			RevertingTxHashes:   revertingTxnHashes,
			BidOptions:          opts,
		},
		Accept:      accept,
		ShouldSlash: shouldSlash,
	}

	if sendPayload {
		bid.Bid.RawTransactions = rawTxns
	} else {
		bid.Bid.TxHashes = txHashes
	}

	store.Insert(bidKey(txHashes), bid)
	return bid, nil
}

func getRevertingTxns(
	ctx context.Context,
	client *ethclient.Client,
	txns []*types.Transaction,
) ([]string, error) {
	var revertingTxns []string
	// do batch call
	batch := make([]rpc.BatchElem, 0, len(txns))
	for _, h := range txns {
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{h.Hash()},
			Result: new(types.Receipt),
		})
	}

	err := client.Client().BatchCallContext(ctx, batch)
	if err != nil {
		return nil, err
	}

	for i, b := range batch {
		if b.Error != nil {
			return nil, b.Error
		}
		receipt := b.Result.(*types.Receipt)
		if receipt.Status != types.ReceiptStatusSuccessful {
			revertingTxns = append(revertingTxns, txns[i].Hash().String())
		}
	}

	return revertingTxns, nil
}

const (
	PRECISION           = 1e16
	ONE_HUNDRED_PERCENT = 100 * PRECISION
	FEE_PERCENT         = 5 * PRECISION
)

func computeResidualAfterDecay(startTimestamp, endTimestamp, commitTimestamp uint64) *big.Int {
	if startTimestamp >= endTimestamp || startTimestamp > commitTimestamp || endTimestamp <= commitTimestamp {
		return big.NewInt(0)
	}

	totalTime := endTimestamp - startTimestamp
	timePassed := commitTimestamp - startTimestamp
	decayPercentage := float64(timePassed) / float64(totalTime)
	residual := 1 - decayPercentage

	residualPercentageRound := math.Round(residual * ONE_HUNDRED_PERCENT)
	if residualPercentageRound > ONE_HUNDRED_PERCENT {
		residualPercentageRound = ONE_HUNDRED_PERCENT
	}

	return big.NewInt(int64(residualPercentageRound))
}
