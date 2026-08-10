package tracker

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
)

type stubStore struct {
	amounts map[common.Hash]*big.Int
	err     error
}

func (s *stubStore) UpdateSettlementStatus(
	_ context.Context,
	_ common.Hash,
	_ bool,
	_ common.Address,
) error {
	return nil
}

func (s *stubStore) UpdateSettlementPayment(
	_ context.Context,
	_ common.Hash,
	_ *big.Int,
	_ *big.Int,
) error {
	return nil
}

func (s *stubStore) GetBidAmountsForTxns(
	_ context.Context,
	_ []common.Hash,
) (map[common.Hash]*big.Int, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.amounts, nil
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// flakyStore returns the bid amounts on the first lookup and errors on
// every later lookup.
type flakyStore struct {
	stubStore
	mu    sync.Mutex
	calls int
}

func (s *flakyStore) GetBidAmountsForTxns(
	_ context.Context,
	_ []common.Hash,
) (map[common.Hash]*big.Int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls++
	if s.calls > 1 {
		return nil, errors.New("transient lookup failure")
	}
	return s.amounts, nil
}

func (s *flakyStore) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls
}

// stubClient serves one closed subscription channel per topic first,
// and then open channels that never send. It signals when both loops
// have subscribed again after their closed channel, and records the
// time of every subscription per topic.
type stubClient struct {
	mu              sync.Mutex
	settlementSubs  int
	paymentSubs     int
	settlementTimes []time.Time
	paymentTimes    []time.Time
	resubscribed    chan struct{}
	signalled       bool
}

// maybeSignal runs with the mutex held.
func (c *stubClient) maybeSignal() {
	if !c.signalled && c.settlementSubs >= 2 && c.paymentSubs >= 2 {
		c.signalled = true
		close(c.resubscribed)
	}
}

func (c *stubClient) SubscribeSettlements(_ context.Context) <-chan bidder.SettlementMsg {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.settlementSubs++
	c.settlementTimes = append(c.settlementTimes, time.Now())
	if c.settlementSubs == 1 {
		ch := make(chan bidder.SettlementMsg)
		close(ch)
		return ch
	}
	c.maybeSignal()
	return make(chan bidder.SettlementMsg)
}

func (c *stubClient) SubscribePayments(_ context.Context) <-chan bidder.PaymentMsg {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.paymentSubs++
	c.paymentTimes = append(c.paymentTimes, time.Now())
	if c.paymentSubs == 1 {
		ch := make(chan bidder.PaymentMsg)
		close(ch)
		return ch
	}
	c.maybeSignal()
	return make(chan bidder.PaymentMsg)
}

// recordingStore records every settlement and payment update.
type recordingStore struct {
	stubStore
	mu            sync.Mutex
	statusHashes  []common.Hash
	paymentHashes []common.Hash
}

func (s *recordingStore) UpdateSettlementStatus(
	_ context.Context,
	txnHash common.Hash,
	_ bool,
	_ common.Address,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statusHashes = append(s.statusHashes, txnHash)
	return nil
}

func (s *recordingStore) UpdateSettlementPayment(
	_ context.Context,
	txnHash common.Hash,
	_ *big.Int,
	_ *big.Int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.paymentHashes = append(s.paymentHashes, txnHash)
	return nil
}

// TestSplitPaymentShareFallbackEqualSplit checks the fallback path: when
// the bid amount lookup fails, a bundle payment splits equally between
// the members, the first bundle hash carries the division remainder,
// and the member shares sum to the aggregate exactly.
func TestSplitPaymentShareFallbackEqualSplit(t *testing.T) {
	t.Parallel()

	hashes := []string{
		common.HexToHash("0x01").Hex(),
		common.HexToHash("0x02").Hex(),
		common.HexToHash("0x03").Hex(),
	}
	tr := NewTracker(nil, &stubStore{err: errors.New("lookup failed")}, testLogger())

	wantPayments := []*big.Int{big.NewInt(34), big.NewInt(33), big.NewInt(33)}
	wantRefunds := []*big.Int{big.NewInt(4), big.NewInt(3), big.NewInt(3)}

	totalPayment := big.NewInt(0)
	totalRefund := big.NewInt(0)
	for i, h := range hashes {
		payment, refund := tr.splitPaymentShare(context.Background(), bidder.PaymentMsg{
			TransactionHash: h,
			Payment:         big.NewInt(100),
			Refund:          big.NewInt(10),
			BundleTxnHashes: hashes,
		})
		if payment.Cmp(wantPayments[i]) != 0 {
			t.Fatalf("expected payment share %s for member %d, got %s", wantPayments[i], i, payment)
		}
		if refund.Cmp(wantRefunds[i]) != 0 {
			t.Fatalf("expected refund share %s for member %d, got %s", wantRefunds[i], i, refund)
		}
		totalPayment = new(big.Int).Add(totalPayment, payment)
		totalRefund = new(big.Int).Add(totalRefund, refund)
	}
	if totalPayment.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("expected payment shares to sum to 100, got %s", totalPayment)
	}
	if totalRefund.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("expected refund shares to sum to 10, got %s", totalRefund)
	}
}

// TestSplitPaymentShareMemoizedAcrossBurst checks that every member of
// one bundle payment notification burst uses the same split source: a
// transient store error after the first member message must not mix
// pro-rata and equal shares, and the shares sum to the aggregate
// exactly. It also checks that the memo entry is evicted after one
// hit per member, so a later burst with the same key resolves fresh.
func TestSplitPaymentShareMemoizedAcrossBurst(t *testing.T) {
	t.Parallel()

	hashes := []string{
		common.HexToHash("0x01").Hex(),
		common.HexToHash("0x02").Hex(),
		common.HexToHash("0x03").Hex(),
	}
	amounts := map[common.Hash]*big.Int{
		common.HexToHash(hashes[0]): big.NewInt(3),
		common.HexToHash(hashes[1]): big.NewInt(5),
		common.HexToHash(hashes[2]): big.NewInt(7),
	}
	st := &flakyStore{stubStore: stubStore{amounts: amounts}}
	tr := NewTracker(nil, st, testLogger())

	// The memoized pro-rata result of the single lookup covers the
	// full burst: 100 over 3:5:7 floors to 20, 33 and 46, and the
	// first hash adds the remainder 1. Refund 10 floors to 2, 3 and
	// 4, and the first hash adds the remainder 1.
	wantPayments := []*big.Int{big.NewInt(21), big.NewInt(33), big.NewInt(46)}
	wantRefunds := []*big.Int{big.NewInt(3), big.NewInt(3), big.NewInt(4)}

	totalPayment := big.NewInt(0)
	totalRefund := big.NewInt(0)
	for i, h := range hashes {
		payment, refund := tr.splitPaymentShare(context.Background(), bidder.PaymentMsg{
			TransactionHash: h,
			Payment:         big.NewInt(100),
			Refund:          big.NewInt(10),
			BundleTxnHashes: hashes,
		})
		if payment.Cmp(wantPayments[i]) != 0 {
			t.Fatalf("expected payment share %s for member %d, got %s", wantPayments[i], i, payment)
		}
		if refund.Cmp(wantRefunds[i]) != 0 {
			t.Fatalf("expected refund share %s for member %d, got %s", wantRefunds[i], i, refund)
		}
		totalPayment = new(big.Int).Add(totalPayment, payment)
		totalRefund = new(big.Int).Add(totalRefund, refund)
	}
	if totalPayment.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("expected payment shares to sum to 100, got %s", totalPayment)
	}
	if totalRefund.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("expected refund shares to sum to 10, got %s", totalRefund)
	}
	if count := st.callCount(); count != 1 {
		t.Fatalf("expected one store lookup for the burst, got %d", count)
	}

	// The burst used up its memo entry, so the same key resolves
	// fresh: the store errors now and the whole new burst takes the
	// equal split path.
	payment, refund := tr.splitPaymentShare(context.Background(), bidder.PaymentMsg{
		TransactionHash: hashes[0],
		Payment:         big.NewInt(100),
		Refund:          big.NewInt(10),
		BundleTxnHashes: hashes,
	})
	if payment.Cmp(big.NewInt(34)) != 0 {
		t.Fatalf("expected equal payment share 34 after eviction, got %s", payment)
	}
	if refund.Cmp(big.NewInt(4)) != 0 {
		t.Fatalf("expected equal refund share 4 after eviction, got %s", refund)
	}
	if count := st.callCount(); count != 2 {
		t.Fatalf("expected a fresh store lookup after eviction, got %d", count)
	}
}

// TestClosedSubscriptionResubscribes checks that a closed subscription
// channel does not spin the tracker: the receive loops exit, no store
// update runs for the zero hash, and the outer loops subscribe again.
func TestClosedSubscriptionResubscribes(t *testing.T) {
	t.Parallel()

	client := &stubClient{resubscribed: make(chan struct{})}
	st := &recordingStore{}
	tr := NewTracker(client, st, testLogger())
	// A short backoff keeps the test fast; the backoff itself is
	// covered by TestResubscribeWaitsBackoff.
	tr.backoff = 10 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	done := tr.Start(ctx)

	select {
	case <-client.resubscribed:
	case <-time.After(5 * time.Second):
		t.Fatal("tracker did not subscribe again after the closed channel")
	}

	cancel()
	<-done

	st.mu.Lock()
	defer st.mu.Unlock()
	if len(st.statusHashes) != 0 {
		t.Fatalf(
			"expected no settlement updates from the closed channel, got %d (first %s)",
			len(st.statusHashes), st.statusHashes[0].Hex(),
		)
	}
	if len(st.paymentHashes) != 0 {
		t.Fatalf(
			"expected no payment updates from the closed channel, got %d (first %s)",
			len(st.paymentHashes), st.paymentHashes[0].Hex(),
		)
	}
}

// TestResubscribeWaitsBackoff checks that the tracker waits for the
// resubscribe backoff after a subscription channel closed: the second
// subscription of each topic starts only after the backoff, so a dead
// bidder connection cannot cause a tight subscribe loop.
func TestResubscribeWaitsBackoff(t *testing.T) {
	t.Parallel()

	client := &stubClient{resubscribed: make(chan struct{})}
	tr := NewTracker(client, &recordingStore{}, testLogger())
	tr.backoff = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	done := tr.Start(ctx)

	select {
	case <-client.resubscribed:
	case <-time.After(5 * time.Second):
		t.Fatal("tracker did not subscribe again after the closed channel")
	}

	cancel()
	<-done

	client.mu.Lock()
	defer client.mu.Unlock()
	for topic, times := range map[string][]time.Time{
		"settlement": client.settlementTimes,
		"payment":    client.paymentTimes,
	} {
		if len(times) < 2 {
			t.Fatalf("expected at least two %s subscriptions, got %d", topic, len(times))
		}
		if gap := times[1].Sub(times[0]); gap < tr.backoff {
			t.Fatalf(
				"expected the second %s subscription after at least %s, got %s",
				topic, tr.backoff, gap,
			)
		}
	}
}

// TestSplitPaymentShareProRataDust checks the pro-rata path with
// amounts that do not divide evenly: the shares sum to the aggregate
// exactly and the first bundle hash carries the dust.
func TestSplitPaymentShareProRataDust(t *testing.T) {
	t.Parallel()

	hashes := []string{
		common.HexToHash("0x01").Hex(),
		common.HexToHash("0x02").Hex(),
		common.HexToHash("0x03").Hex(),
	}
	amounts := map[common.Hash]*big.Int{
		common.HexToHash(hashes[0]): big.NewInt(3),
		common.HexToHash(hashes[1]): big.NewInt(5),
		common.HexToHash(hashes[2]): big.NewInt(7),
	}
	tr := NewTracker(nil, &stubStore{amounts: amounts}, testLogger())

	// Payment 100 over amounts 3:5:7 floors to 20, 33 and 46; the
	// first hash adds the remainder 1. Refund 10 floors to 2, 3 and
	// 4; the first hash adds the remainder 1.
	wantPayments := []*big.Int{big.NewInt(21), big.NewInt(33), big.NewInt(46)}
	wantRefunds := []*big.Int{big.NewInt(3), big.NewInt(3), big.NewInt(4)}

	totalPayment := big.NewInt(0)
	totalRefund := big.NewInt(0)
	for i, h := range hashes {
		payment, refund := tr.splitPaymentShare(context.Background(), bidder.PaymentMsg{
			TransactionHash: h,
			Payment:         big.NewInt(100),
			Refund:          big.NewInt(10),
			BundleTxnHashes: hashes,
		})
		if payment.Cmp(wantPayments[i]) != 0 {
			t.Fatalf("expected payment share %s for member %d, got %s", wantPayments[i], i, payment)
		}
		if refund.Cmp(wantRefunds[i]) != 0 {
			t.Fatalf("expected refund share %s for member %d, got %s", wantRefunds[i], i, refund)
		}
		totalPayment = new(big.Int).Add(totalPayment, payment)
		totalRefund = new(big.Int).Add(totalRefund, refund)
	}
	if totalPayment.Cmp(big.NewInt(100)) != 0 {
		t.Fatalf("expected payment shares to sum to 100, got %s", totalPayment)
	}
	if totalRefund.Cmp(big.NewInt(10)) != 0 {
		t.Fatalf("expected refund shares to sum to 10, got %s", totalRefund)
	}
}
