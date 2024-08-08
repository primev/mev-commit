package preconf

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/armon/go-radix"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconfcommitmentstore "github.com/primev/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/contracts/events"
	"golang.org/x/sync/errgroup"
)

const (
	noOfBids = 100
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

	fundsRetrievedKey = func(cmtDigest []byte) string {
		return fmt.Sprintf("fr/%s", string(cmtDigest))
	}

	fundsRewardedKey = func(cmtDigest []byte) string {
		return fmt.Sprintf("frw/%s", string(cmtDigest))
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
			func(c *preconfcommitmentstore.PreconfcommitmentstoreUnopenedCommitmentStored) {
				logger.Info("Received encrypted commitment", "digest", hex.EncodeToString(c.CommitmentDigest[:]))
				store.Insert(encryptCmtKey(c.CommitmentDigest[:]), c)
			},
		),
		events.NewEventHandler(
			"OpenedCommitmentStored",
			func(c *preconfcommitmentstore.PreconfcommitmentstoreOpenedCommitmentStored) {
				logger.Info(
					"Received opened commitment",
					"digest", hex.EncodeToString(c.CommitmentDigest[:]),
					"index", hex.EncodeToString(c.CommitmentIndex[:]),
					"decay_start", c.DecayStartTimeStamp,
					"decay_end", c.DecayEndTimeStamp,
					"dispatch_timestamp", c.DispatchTimestamp,
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
			"FundsRetrieved",
			func(c *bidderregistry.BidderregistryFundsRetrieved) {
				logger.Info("Retrieved funds", "digest", hex.EncodeToString(c.CommitmentDigest[:]))
				store.Insert(fundsRetrievedKey(c.CommitmentDigest[:]), c)
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
			"NewL1Block",
			func(c *blocktracker.BlocktrackerNewL1Block) {
				logger.Info("Received new L1 block")
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
					val, ok := store.Get(bidKey(entry.Bid.TxHashes))
					if !ok {
						logger.Error("Bid not found in store", "key", bidKey(entry.Bid.TxHashes))
						return fmt.Errorf("bid not found in store")
					}
					val.(*BidEntry).Preconfs = preconfs
					store.Insert(bidKey(entry.Bid.TxHashes), val)
					logger.Info("Received preconfs", "count", len(preconfs))
				}
			}
		})
	}

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	count := 0
	lastWinnerBlock := 0
DONE:
	for {
		select {
		case <-egCtx.Done():
			egCancel()
			return nil
		case <-tick.C:
			if count == noOfBids {
				_, ok := store.Get(blkWinnerKey(uint64(lastWinnerBlock + 5)))
				if ok {
					// allow enough time for everything to settle
					egCancel()
					logger.Info("All bids sent")
					break DONE
				}
			} else {
				for _, b := range bidders {
					entry, err := getRandomBid(cluster, store)
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
		} else {
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
					ecmt := ec.(*preconfcommitmentstore.PreconfcommitmentstoreUnopenedCommitmentStored)
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
					if entry.ShouldSlash {
						if !pcmt.(*oracle.OracleCommitmentProcessed).IsSlash {
							logger.Error("Provider should be slashed", "entry", entry)
							return fmt.Errorf("provider should be slashed")
						}
						_, ok := store.Get(fundsRetrievedKey(cmtDigest))
						if !ok {
							logger.Error("Funds not retrieved", "entry", entry)
							return fmt.Errorf("funds not retrieved")
						}
					} else {
						if pcmt.(*oracle.OracleCommitmentProcessed).IsSlash {
							// check if any of the transactions were not successful,
							// if so, the provider should not be slashed. Test doesnt
							// handle reverting transactions.
							failedTxnPresent := false
							for _, h := range entry.Bid.TxHashes {
								receipt, err := cluster.L1RPC().TransactionReceipt(
									context.Background(),
									common.HexToHash(h),
								)
								if err != nil {
									logger.Error(
										"failed getting transaction receipt",
										"error", err,
										"entry", entry,
										"hash", h,
									)
								}
								if receipt.Status != types.ReceiptStatusSuccessful {
									failedTxnPresent = true
								}
							}
							if !failedTxnPresent {
								logger.Error("Provider should not be slashed", "entry", entry)
								return fmt.Errorf("provider should not be slashed")
							}
							continue
						}
						_, ok := store.Get(fundsRewardedKey(cmtDigest))
						if !ok {
							logger.Error("Funds not rewarded", "entry", entry)
							return fmt.Errorf("funds not rewarded")
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
	}

	return nil
}

func getRandomBid(
	o orchestrator.Orchestrator,
	store *radix.Tree,
) (*BidEntry, error) {
	blkNum, err := o.L1RPC().BlockNumber(context.Background())
	if err != nil {
		return nil, err
	}

	blk, found := store.Get(blkKey(blkNum))
	if !found {
		blk, err = o.L1RPC().BlockByNumber(context.Background(), big.NewInt(int64(blkNum)))
		if err != nil {
			return nil, err
		}
		store.Insert(blkKey(blkNum), blk)
	}

	if len(blk.(*types.Block).Transactions()) == 0 {
		return nil, errNoTxnsInBlock
	}

	idx := rand.Intn(len(blk.(*types.Block).Transactions()))
	bundleLen := rand.Intn(5) + 1
	if idx+bundleLen > len(blk.(*types.Block).Transactions()) {
		bundleLen = len(blk.(*types.Block).Transactions()) - idx
	}

	txHashes := make([]string, 0, bundleLen)
	for i := idx; i < idx+bundleLen; i++ {
		txHashes = append(
			txHashes,
			strings.TrimPrefix(blk.(*types.Block).Transactions()[i].Hash().String(), "0x"),
		)
	}

	// accept 90% of bids
	accept := rand.Intn(100) < 90
	// slash 10% of accepted bids
	shouldSlash := rand.Intn(100) < 10
	// amount between 5M and 6M
	amount := 5_000_000 + rand.Intn(1_000_000)

	if shouldSlash {
		if len(txHashes) > 1 {
			rand.Shuffle(len(txHashes), func(i, j int) {
				txHashes[i], txHashes[j] = txHashes[j], txHashes[i]
			})
		} else {
			// get random tx hash
			randBytes := make([]byte, 32)
			_, _ = crand.Read(randBytes)
			txHashes[0] = strings.TrimPrefix(common.BytesToHash(randBytes).String(), "0x")
		}
	}

	bid := &BidEntry{
		Bid: &bidderapiv1.Bid{
			TxHashes:            txHashes,
			Amount:              fmt.Sprintf("%d", amount),
			BlockNumber:         int64(blkNum),
			DecayStartTimestamp: time.Now().UnixMilli(),
			DecayEndTimestamp:   time.Now().Add(5 * time.Second).UnixMilli(),
		},
		Accept:      accept,
		ShouldSlash: shouldSlash,
	}

	store.Insert(bidKey(bid.Bid.TxHashes), bid)
	return bid, nil
}
