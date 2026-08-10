package tracker

import (
	"context"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
	"golang.org/x/sync/errgroup"
)

// splitSourceCap bounds the split decision memo. Each entry serves one
// bundle payment notification burst and is evicted after one hit per
// bundle member. The cap protects the memo when member messages get
// lost and their entry never reaches its eviction count.
const splitSourceCap = 1024

// resubscribeBackoff is the wait before a new subscription after a
// subscription channel closed. A dead bidder connection returns a
// channel that closes at once; the wait keeps the subscribe loop from
// spinning against it.
const resubscribeBackoff = 5 * time.Second

type BidderClient interface {
	SubscribeSettlements(ctx context.Context) <-chan bidder.SettlementMsg
	SubscribePayments(ctx context.Context) <-chan bidder.PaymentMsg
}

type Store interface {
	UpdateSettlementStatus(
		ctx context.Context,
		txnHash common.Hash,
		isSlashed bool,
		provider common.Address,
	) error
	UpdateSettlementPayment(
		ctx context.Context,
		txnHash common.Hash,
		payment *big.Int,
		refund *big.Int,
	) error
	GetBidAmountsForTxns(
		ctx context.Context,
		txnHashes []common.Hash,
	) (map[common.Hash]*big.Int, error)
}

// splitSource is the memoized split decision for one bundle payment
// notification burst. A nil amounts map marks the equal split path.
// The hits count tracks how many member messages of the burst used
// the decision; the entry is evicted after one hit per member.
type splitSource struct {
	amounts map[common.Hash]*big.Int
	total   *big.Int
	hits    int
}

type tracker struct {
	client BidderClient
	store  Store
	logger *slog.Logger

	// backoff is the wait before a resubscribe after a subscription
	// channel closed. Tests shorten it.
	backoff time.Duration

	splitMu      sync.Mutex
	splitSources map[string]*splitSource
	splitOrder   []string
}

func NewTracker(
	client BidderClient,
	store Store,
	logger *slog.Logger,
) *tracker {
	return &tracker{
		client:       client,
		store:        store,
		logger:       logger,
		backoff:      resubscribeBackoff,
		splitSources: make(map[string]*splitSource),
	}
}

func (t *tracker) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}

			sub := t.client.SubscribeSettlements(egCtx)
		SETTLEMENTS:
			for {
				select {
				case <-egCtx.Done():
					return egCtx.Err()
				case msg, more := <-sub:
					if !more {
						// The subscription channel closed. Leave the
						// receive loop, so the outer loop checks the
						// context and subscribes again. A receive
						// from the closed channel would yield endless
						// zero-value messages for the zero hash.
						t.logger.Warn("Settlement subscription closed, subscribing again")
						break SETTLEMENTS
					}
					// Process settlement message
					if err := t.store.UpdateSettlementStatus(
						egCtx,
						common.HexToHash(msg.TransactionHash),
						msg.IsSlash,
						common.HexToAddress(msg.Provider),
					); err != nil {
						t.logger.Error(
							"Failed to update settlement status",
							"error", err,
							"txnHash", msg.TransactionHash,
							"isSlash", msg.IsSlash,
							"provider", msg.Provider,
						)
					} else {
						t.logger.Info("Updated settlement status",
							"txnHash", msg.TransactionHash,
							"isSlash", msg.IsSlash,
							"provider", msg.Provider,
						)
					}
				}
			}
			// The subscription channel closed. Wait before the next
			// subscription, so a dead bidder connection does not
			// cause a tight subscribe loop.
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-time.After(t.backoff):
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}

			sub := t.client.SubscribePayments(egCtx)
		PAYMENTS:
			for {
				select {
				case <-egCtx.Done():
					return egCtx.Err()
				case msg, more := <-sub:
					if !more {
						// The subscription channel closed. Leave the
						// receive loop, so the outer loop checks the
						// context and subscribes again. A receive
						// from the closed channel would yield endless
						// zero-value messages for the zero hash.
						t.logger.Warn("Payment subscription closed, subscribing again")
						break PAYMENTS
					}
					// Process payment message
					payment, refund := t.splitPaymentShare(egCtx, msg)
					if err := t.store.UpdateSettlementPayment(
						egCtx,
						common.HexToHash(msg.TransactionHash),
						payment,
						refund,
					); err != nil {
						t.logger.Error(
							"Failed to update settlement payment",
							"error", err,
							"txnHash", msg.TransactionHash,
							"payment", payment,
							"refund", refund,
						)
					} else {
						t.logger.Info(
							"Updated settlement payment",
							"txnHash", msg.TransactionHash,
							"payment", payment,
							"refund", refund,
						)
					}
				}
			}
			// The subscription channel closed. Wait before the next
			// subscription, so a dead bidder connection does not
			// cause a tight subscribe loop.
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-time.After(t.backoff):
			}
		}
	})

	go func() {
		defer close(done)

		if err := eg.Wait(); err != nil {
			t.logger.Error("Tracker encountered an error", "error", err)
		}
	}()
	return done
}

// splitPaymentShare computes the share of a payment message for its
// transaction. A bundle bid produces one aggregate payment for all bundle
// transactions. The share of each transaction is proportional to its
// stored bid amount. If the shares cannot be loaded, the aggregate is
// split equally between the bundle members. All members of one
// notification burst use one memoized split source: the first member
// message resolves the decision and the later member messages of the
// same burst reuse it, so pro-rata and equal shares never mix within
// one burst. In both cases the member with the first bundle hash adds
// the division remainder, so the sum of the member shares equals the
// aggregate exactly within each burst and the group is never charged
// more than the aggregate. Single transaction bids keep the full
// amount.
func (t *tracker) splitPaymentShare(ctx context.Context, msg bidder.PaymentMsg) (*big.Int, *big.Int) {
	if len(msg.BundleTxnHashes) < 2 {
		return msg.Payment, msg.Refund
	}

	hashes := make([]common.Hash, 0, len(msg.BundleTxnHashes))
	for _, h := range msg.BundleTxnHashes {
		hashes = append(hashes, common.HexToHash(h))
	}
	own := common.HexToHash(msg.TransactionHash)

	src := t.splitSourceFor(ctx, msg, hashes)
	if src.amounts == nil {
		return equalShare(msg.Payment, own, hashes), equalShare(msg.Refund, own, hashes)
	}

	payment := proRataShare(msg.Payment, src.amounts, hashes, own, src.total)
	refund := proRataShare(msg.Refund, src.amounts, hashes, own, src.total)
	return payment, refund
}

// splitSourceFor returns the split decision for one bundle payment
// notification burst. The burst key combines the bundle hash list and
// the aggregate amounts. The first member message of a burst resolves
// the decision; every later member message of the same burst reuses
// it. An entry is evicted after one hit per bundle member. The memo
// is capped; when it is full, the oldest entry is evicted first.
func (t *tracker) splitSourceFor(ctx context.Context, msg bidder.PaymentMsg, hashes []common.Hash) *splitSource {
	key := strings.Join(msg.BundleTxnHashes, ",") + "|" + msg.Payment.String() + "|" + msg.Refund.String()

	t.splitMu.Lock()
	defer t.splitMu.Unlock()

	if src, found := t.splitSources[key]; found {
		src.hits++
		if src.hits >= len(msg.BundleTxnHashes) {
			t.evictSplitSource(key)
		}
		return src
	}

	src := t.resolveSplitSource(ctx, msg, hashes)
	src.hits = 1
	if len(t.splitOrder) >= splitSourceCap {
		t.evictSplitSource(t.splitOrder[0])
	}
	t.splitSources[key] = src
	t.splitOrder = append(t.splitOrder, key)
	return src
}

// evictSplitSource removes one memo entry and its insertion order
// record. The caller holds splitMu.
func (t *tracker) evictSplitSource(key string) {
	delete(t.splitSources, key)
	order := make([]string, 0, len(t.splitOrder))
	for _, k := range t.splitOrder {
		if k != key {
			order = append(order, k)
		}
	}
	t.splitOrder = order
}

// resolveSplitSource loads the stored bid amounts for one bundle and
// builds the split decision: pro-rata over the stored amounts, or the
// equal split marker when the amounts cannot be loaded in full.
func (t *tracker) resolveSplitSource(ctx context.Context, msg bidder.PaymentMsg, hashes []common.Hash) *splitSource {
	amounts, err := t.store.GetBidAmountsForTxns(ctx, hashes)
	if err != nil {
		t.logger.Error(
			"Failed to get bid amounts for bundle, using equal split",
			"error", err,
			"txnHash", msg.TransactionHash,
		)
		return &splitSource{}
	}

	total := big.NewInt(0)
	for _, h := range hashes {
		amount, found := amounts[h]
		if !found {
			t.logger.Error(
				"Missing bid amount for bundle member, using equal split",
				"txnHash", h.Hex(),
			)
			return &splitSource{}
		}
		total = new(big.Int).Add(total, amount)
	}

	if total.Sign() <= 0 {
		t.logger.Error(
			"Non-positive bid amount total for bundle, using equal split",
			"txnHash", msg.TransactionHash,
		)
		return &splitSource{}
	}

	return &splitSource{amounts: amounts, total: total}
}

// equalShare splits an aggregate amount into equal parts, one for each
// bundle hash. The member with the first bundle hash adds the division
// remainder. Every member computes its share independently and the
// shares sum to the aggregate exactly.
func equalShare(aggregate *big.Int, own common.Hash, hashes []common.Hash) *big.Int {
	count := big.NewInt(int64(len(hashes)))
	share, remainder := new(big.Int).QuoRem(aggregate, count, new(big.Int))
	if own == hashes[0] {
		return new(big.Int).Add(share, remainder)
	}
	return share
}

// proRataShare computes one member's part of an aggregate amount, in
// proportion to the stored member bid amounts. The floor divisions can
// leave a remainder; the member with the first bundle hash adds it.
// Every member computes its share independently and the shares sum to
// the aggregate exactly.
func proRataShare(
	aggregate *big.Int,
	amounts map[common.Hash]*big.Int,
	hashes []common.Hash,
	own common.Hash,
	total *big.Int,
) *big.Int {
	share := big.NewInt(0)
	floorSum := big.NewInt(0)
	for _, h := range hashes {
		memberShare := new(big.Int).Div(new(big.Int).Mul(aggregate, amounts[h]), total)
		floorSum = new(big.Int).Add(floorSum, memberShare)
		if h == own {
			share = memberShare
		}
	}
	if own == hashes[0] {
		remainder := new(big.Int).Sub(aggregate, floorSum)
		return new(big.Int).Add(share, remainder)
	}
	return share
}
