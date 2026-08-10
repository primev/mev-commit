package sender_test

import (
	"context"
	"errors"
	"io"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/bidder"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	"github.com/primev/mev-commit/tools/preconf-rpc/sim"
	"github.com/primev/mev-commit/x/util"
)

type result struct {
	txn         *sender.Transaction
	commitments []*bidderapiv1.Commitment
	blockNumber int64
	logs        []*types.Log
}

type mockStore struct {
	mu               sync.Mutex
	queued           map[common.Address][]*sender.Transaction
	nonce            map[common.Address]uint64
	balances         map[common.Address]*big.Int
	byHash           map[common.Hash]*sender.Transaction
	preconfirmedTxns chan result
	addBalanceCalls  map[common.Address]int
	// honorCtx makes StoreTransaction fail after its context ends,
	// like a real database write. Set it before the sender starts.
	honorCtx bool
}

func newMockStore() *mockStore {
	return &mockStore{
		queued:           make(map[common.Address][]*sender.Transaction),
		nonce:            make(map[common.Address]uint64),
		balances:         make(map[common.Address]*big.Int),
		preconfirmedTxns: make(chan result, 10),
		byHash:           make(map[common.Hash]*sender.Transaction),
		addBalanceCalls:  make(map[common.Address]int),
	}
}

func (m *mockStore) AddQueuedTransaction(_ context.Context, tx *sender.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queued[tx.Sender] = append(m.queued[tx.Sender], tx)
	m.nonce[tx.Sender] = tx.Nonce()

	return nil
}

func (m *mockStore) GetQueuedTransactions(_ context.Context) ([]*sender.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var txns []*sender.Transaction

	for _, acctTxns := range m.queued {
		if len(acctTxns) == 0 {
			continue
		}
		txns = append(txns, acctTxns[0])
	}

	return txns, nil
}

func (m *mockStore) GetCurrentNonce(_ context.Context, sender common.Address) (uint64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	nonce, exists := m.nonce[sender]
	if !exists {
		return 0, false
	}

	return nonce, true
}

func (m *mockStore) HasBalance(ctx context.Context, sender common.Address, amount *big.Int) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	balance, exists := m.balances[sender]
	if !exists {
		return false
	}
	return balance.Cmp(amount) >= 0
}

func (m *mockStore) AddBalance(ctx context.Context, account common.Address, amount *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.addBalanceCalls[account]++
	if _, exists := m.balances[account]; !exists {
		m.balances[account] = amount
	} else {
		newBalance := new(big.Int).Add(m.balances[account], amount)
		m.balances[account] = newBalance
	}

	return nil
}

func (m *mockStore) addBalanceCount(account common.Address) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.addBalanceCalls[account]
}

func (m *mockStore) DeductBalance(ctx context.Context, account common.Address, amount *big.Int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.balances[account]; !exists {
		return errors.New("account does not exist")
	}
	newBalance := new(big.Int).Sub(m.balances[account], amount)
	if newBalance.Sign() < 0 {
		return errors.New("insufficient balance")
	}
	m.balances[account] = newBalance
	return nil
}

func (m *mockStore) StoreTransaction(
	ctx context.Context,
	txn *sender.Transaction,
	commitments []*bidderapiv1.Commitment,
	logs []*types.Log,
) error {
	if m.honorCtx && ctx.Err() != nil {
		return ctx.Err()
	}
	m.preconfirmedTxns <- result{
		txn:         txn,
		commitments: commitments,
		blockNumber: txn.BlockNumber,
		logs:        logs,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, queuedTxn := range m.queued[txn.Sender] {
		if queuedTxn.Hash() == txn.Hash() {
			// Remove the transaction from the queue
			m.queued[txn.Sender] = append(m.queued[txn.Sender][:i], m.queued[txn.Sender][i+1:]...)
			break
		}
	}
	m.byHash[txn.Hash()] = txn
	return nil
}

func (m *mockStore) GetTransactionByHash(
	ctx context.Context,
	hash common.Hash,
) (*sender.Transaction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	txn, exists := m.byHash[hash]
	if !exists {
		return nil, errors.New("transaction not found")
	}

	return txn, nil
}

func (m *mockStore) StoreReceipt(
	ctx context.Context,
	receipt *types.Receipt,
) error {
	return nil
}

type bidOp struct {
	bidAmount   *big.Int
	slashAmount *big.Int
	rawTxs      []string
	opts        *bidder.BidOpts
}

type bidResponse struct {
	statusCh chan bidder.BidStatus
	err      error
}

type mockBidder struct {
	optinEstimate chan int64
	in            chan bidOp
	out           chan bidResponse
}

func (m *mockBidder) Estimate() (int64, error) {
	estimate := <-m.optinEstimate
	return estimate, nil
}

func (m *mockBidder) Bid(
	ctx context.Context,
	bidAmount *big.Int,
	slashAmount *big.Int,
	rawTxs []string,
	opts *bidder.BidOpts,
) (chan bidder.BidStatus, error) {
	m.in <- bidOp{
		bidAmount:   bidAmount,
		slashAmount: slashAmount,
		rawTxs:      rawTxs,
		opts:        opts,
	}
	res := <-m.out

	return res.statusCh, res.err
}

func (m *mockBidder) ConnectedProviders(ctx context.Context) ([]string, error) {
	return []string{"provider1", "provider2"}, nil
}

type mockPricer struct {
	out chan map[int64]float64
}

func (m *mockPricer) EstimatePrice(ctx context.Context) map[int64]float64 {
	select {
	case prices := <-m.out:
		if prices == nil {
			return nil
		}
		return prices
	case <-ctx.Done():
		return nil
	}
}

type blockNoOp struct {
	block             uint64
	timeTillNextBlock time.Duration
}

type mockBlockTracker struct {
	out   chan uint64
	bnIn  chan struct{}
	bnOut chan blockNoOp
	bnErr chan error
	mu    sync.Mutex
	lbf   *big.Int
	nbf   *big.Int
	incl  map[common.Hash]chan uint64
}

// inclusionFor registers a dedicated inclusion channel for one
// transaction hash. WaitForTxnInclusion returns that channel for the
// hash, so a test can signal inclusion per transaction instead of
// through the shared out channel. Register before the sender starts.
func (m *mockBlockTracker) inclusionFor(txnHash common.Hash) chan uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.incl == nil {
		m.incl = make(map[common.Hash]chan uint64)
	}
	ch, found := m.incl[txnHash]
	if !found {
		ch = make(chan uint64, 1)
		m.incl[txnHash] = ch
	}
	return ch
}

func (m *mockBlockTracker) WaitForTxnInclusion(txnHash common.Hash) chan uint64 {
	m.mu.Lock()
	ch, found := m.incl[txnHash]
	m.mu.Unlock()
	if found {
		return ch
	}

	includedCh := make(chan uint64, 1)
	go func() {
		included := <-m.out
		includedCh <- included
	}()
	return includedCh
}

func (m *mockBlockTracker) NextBlockNumber() (uint64, time.Duration, error) {
	m.bnIn <- struct{}{}

	select {
	case <-m.bnErr:
		return 0, 0, errors.New("error getting next block number")
	case op := <-m.bnOut:
		return op.block, op.timeTillNextBlock, nil
	}
}

func (m *mockBlockTracker) LatestBlockNumber() uint64 {
	return 0
}

func (m *mockBlockTracker) AccountNonce(ctx context.Context, account common.Address) (uint64, error) {
	return 0, nil
}

func (m *mockBlockTracker) LatestBaseFee() *big.Int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.lbf == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(m.lbf)
}

func (m *mockBlockTracker) NextBaseFee() *big.Int {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.nbf == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(m.nbf)
}

func (m *mockBlockTracker) setNextBaseFee(fee *big.Int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nbf = fee
}

type mockTransferer struct{}

func (m *mockTransferer) Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error {
	return nil
}

type mockNotifier struct {
	mu            sync.Mutex
	notifications []string
}

func (m *mockNotifier) NotifyTransactionStatus(txn *sender.Transaction, attempts, blocks int, start time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = append(m.notifications, txn.Hash().Hex())
}

func (m *mockNotifier) notified(hash common.Hash) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range m.notifications {
		if h == hash.Hex() {
			return true
		}
	}
	return false
}

type mockSimulator struct {
	failRaw string
}

func (m *mockSimulator) Simulate(ctx context.Context, rawTx string, _ sim.SimState) ([]*types.Log, bool, error) {
	if m.failRaw != "" && rawTx == m.failRaw {
		return nil, false, errors.New("simulation failed")
	}
	return []*types.Log{}, false, nil
}

type mockBackrunner struct{}

func (m *mockBackrunner) Backrun(ctx context.Context, rawTx string, commitments []*bidderapiv1.Commitment) error {
	return nil
}

type MockExplorerSubmitter struct{}

func (m *MockExplorerSubmitter) Submit(ctx context.Context, tx *types.Transaction, from common.Address) error {
	return nil
}

func TestSender(t *testing.T) {
	t.Parallel()

	st := newMockStore()
	testPricer := &mockPricer{
		out: make(chan map[int64]float64, 10),
	}
	bidderImpl := &mockBidder{
		optinEstimate: make(chan int64, 10),
		in:            make(chan bidOp, 10),
		out:           make(chan bidResponse, 10),
	}
	blockTracker := &mockBlockTracker{
		out:   make(chan uint64, 10),
		bnIn:  make(chan struct{}, 10),
		bnOut: make(chan blockNoOp, 10),
		bnErr: make(chan error, 1),
	}
	notifier := &mockNotifier{}

	sndr, err := sender.NewTxSender(
		st,
		bidderImpl,
		testPricer,
		blockTracker,
		&mockTransferer{},
		notifier,
		&mockSimulator{},
		&mockBackrunner{},
		big.NewInt(1), // Settlement chain ID
		&MockExplorerSubmitter{},
		nil,  // no log encryption key in tests
		true, // bundle bids enabled
		util.NewTestLogger(os.Stdout),
	)
	if err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := sndr.Start(ctx)

	tx1 := &sender.Transaction{
		Transaction: types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(100),
			21000,
			big.NewInt(1),
			nil,
		),
		Sender: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Type:   sender.TxTypeRegular,
		Raw:    "0x1234567890123456789012345678901234567890",
	}

	if err := st.AddBalance(ctx, tx1.Sender, big.NewInt(5e18)); err != nil {
		t.Fatalf("failed to add balance: %v", err)
	}

	if err := sndr.Enqueue(ctx, tx1); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	waitCh := sndr.WaitForReceiptAvailable(ctx, tx1.Hash())

	// Simulate opted in block
	bidderImpl.optinEstimate <- 2

	<-blockTracker.bnIn
	blockTracker.bnErr <- errors.New("simulated error for testing")

	bidderImpl.optinEstimate <- 7

	<-blockTracker.bnIn

	blockTracker.bnOut <- blockNoOp{
		block:             1,
		timeTillNextBlock: 5 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}

	// Simulate transaction inclusion
	blockTracker.out <- 1

	// Simulate a bid response
	bidOp := <-bidderImpl.in
	if len(bidOp.rawTxs) != 1 || bidOp.rawTxs[0] != tx1.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %v", tx1.Raw, bidOp.rawTxs)
	}
	resC := make(chan bidder.BidStatus, 3)
	resC <- bidder.BidStatus{
		Type: bidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     1,
			ProviderAddress: "provider1",
		},
	}
	resC <- bidder.BidStatus{
		Type: bidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     1,
			ProviderAddress: "provider2",
		},
	}
	close(resC)
	bidderImpl.out <- bidResponse{statusCh: resC}

	<-waitCh

	res := <-st.preconfirmedTxns
	if res.txn == nil {
		t.Fatal("expected a preconfirmed transaction, got nil")
	}
	if res.blockNumber != 1 {
		t.Fatalf("expected block number 1, got %d", res.blockNumber)
	}
	if res.txn.Sender != tx1.Sender {
		t.Fatalf("expected sender %s, got %s", tx1.Sender.Hex(), res.txn.Sender.Hex())
	}
	if res.txn.Nonce() != tx1.Nonce() {
		t.Fatalf("expected nonce %d, got %d", tx1.Nonce(), res.txn.Nonce())
	}
	if res.txn.Type != tx1.Type {
		t.Fatalf("expected transaction type %d, got %d", tx1.Type, res.txn.Type)
	}
	if res.txn.Hash() != tx1.Hash() {
		t.Fatalf("expected transaction hash %s, got %s", tx1.Hash().Hex(), res.txn.Hash().Hex())
	}
	// Check that the commitments are as expected
	if len(res.commitments) != 2 {
		t.Fatalf("expected 2 commitments, got %d", len(res.commitments))
	}

	tx2 := &sender.Transaction{
		Transaction: types.NewTransaction(
			1,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(1e18),
			21000,
			big.NewInt(1),
			nil,
		),
		Sender: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Type:   sender.TxTypeDeposit,
		Raw:    "0x1234567890123456789012345678901234567890",
	}

	if err := sndr.Enqueue(ctx, tx2); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	// Simulate non opted in block
	bidderImpl.optinEstimate <- 20

	<-blockTracker.bnIn
	blockTracker.bnOut <- blockNoOp{
		block:             2,
		timeTillNextBlock: 5 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}

	// Simulate a bid response
	bidOp = <-bidderImpl.in
	if len(bidOp.rawTxs) != 1 || bidOp.rawTxs[0] != tx2.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %v", tx2.Raw, bidOp.rawTxs)
	}
	resC = make(chan bidder.BidStatus, 3)
	// Simulate retry due to incomplete commitments
	close(resC)
	bidderImpl.out <- bidResponse{statusCh: resC}

	// Simulate non opted in block
	bidderImpl.optinEstimate <- 18

	<-blockTracker.bnIn
	blockTracker.bnOut <- blockNoOp{
		block:             2,
		timeTillNextBlock: 5 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}

	// Simulate transaction inclusion
	blockTracker.out <- 2

	// Simulate a bid response
	bidOp = <-bidderImpl.in
	if len(bidOp.rawTxs) != 1 || bidOp.rawTxs[0] != tx2.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %v", tx2.Raw, bidOp.rawTxs)
	}
	resC = make(chan bidder.BidStatus, 3)
	resC <- bidder.BidStatus{
		Type: bidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     2,
			ProviderAddress: "provider1",
		},
	}
	close(resC)
	bidderImpl.out <- bidResponse{statusCh: resC}

	res = <-st.preconfirmedTxns
	if res.txn == nil {
		t.Fatal("expected a preconfirmed transaction, got nil")
	}
	if res.blockNumber != 2 {
		t.Fatalf("expected block number 2, got %d", res.blockNumber)
	}
	if res.txn.Sender != tx2.Sender {
		t.Fatalf("expected sender %s, got %s", tx2.Sender.Hex(), res.txn.Sender.Hex())
	}
	if res.txn.Nonce() != tx2.Nonce() {
		t.Fatalf("expected nonce %d, got %d", tx2.Nonce(), res.txn.Nonce())
	}
	if res.txn.Type != tx2.Type {
		t.Fatalf("expected transaction type %d, got %d", tx2.Type, res.txn.Type)
	}
	if res.txn.Hash() != tx2.Hash() {
		t.Fatalf("expected transaction hash %s, got %s", tx2.Hash().Hex(), res.txn.Hash().Hex())
	}
	// Check that the commitments are as expected
	if len(res.commitments) != 1 {
		t.Fatalf("expected 1 commitment, got %d", len(res.commitments))
	}

	cancel()
	<-done

	if len(notifier.notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(notifier.notifications))
	}
}

func TestCancelTransaction(t *testing.T) {
	t.Parallel()

	st := newMockStore()
	testPricer := &mockPricer{
		out: make(chan map[int64]float64, 10),
	}
	bidderImpl := &mockBidder{
		optinEstimate: make(chan int64),
		in:            make(chan bidOp, 10),
		out:           make(chan bidResponse, 10),
	}
	blockTracker := &mockBlockTracker{
		out:   make(chan uint64, 10),
		bnIn:  make(chan struct{}, 10),
		bnOut: make(chan blockNoOp, 10),
		bnErr: make(chan error, 3),
	}

	sndr, err := sender.NewTxSender(
		st,
		bidderImpl,
		testPricer,
		blockTracker,
		&mockTransferer{},
		&mockNotifier{},
		&mockSimulator{},
		&mockBackrunner{},
		big.NewInt(1), // Settlement chain ID
		&MockExplorerSubmitter{},
		nil,  // no log encryption key in tests
		true, // bundle bids enabled
		util.NewTestLogger(os.Stdout),
	)
	if err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := sndr.Start(ctx)

	tx1 := &sender.Transaction{
		Transaction: types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(100),
			21000,
			big.NewInt(1),
			nil,
		),
		Sender: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Type:   sender.TxTypeRegular,
		Raw:    "0x1234567890123456789012345678901234567890",
	}

	if err := st.AddBalance(ctx, tx1.Sender, big.NewInt(5e18)); err != nil {
		t.Fatalf("failed to add balance: %v", err)
	}

	if err := sndr.Enqueue(ctx, tx1); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	go func() {
		for {
			select {
			case <-blockTracker.bnIn:
			case <-ctx.Done():
				return
			}
		}
	}()

	blockTracker.bnErr <- errors.New("simulated error for testing")
	blockTracker.bnErr <- errors.New("simulated error for testing")
	bidderImpl.optinEstimate <- 2

	cancelled, err := sndr.CancelTransaction(ctx, tx1.Hash())
	if err != nil {
		t.Fatalf("failed to cancel transaction: %v", err)
	}
	if !cancelled {
		t.Fatal("expected transaction to be cancelled, but it was not")
	}

	cancel()
	<-done
}

func TestIgnoreProvidersOnRetry(t *testing.T) {
	t.Parallel()

	st := newMockStore()
	testPricer := &mockPricer{
		out: make(chan map[int64]float64, 10),
	}
	bidderImpl := &mockBidder{
		optinEstimate: make(chan int64, 10),
		in:            make(chan bidOp, 10),
		out:           make(chan bidResponse, 10),
	}
	blockTracker := &mockBlockTracker{
		out:   make(chan uint64, 10),
		bnIn:  make(chan struct{}, 10),
		bnOut: make(chan blockNoOp, 10),
		bnErr: make(chan error, 1),
	}
	notifier := &mockNotifier{}

	sndr, err := sender.NewTxSender(
		st,
		bidderImpl,
		testPricer,
		blockTracker,
		&mockTransferer{},
		notifier,
		&mockSimulator{},
		&mockBackrunner{},
		big.NewInt(1), // Settlement chain ID
		&MockExplorerSubmitter{},
		nil,  // no log encryption key in tests
		true, // bundle bids enabled
		util.NewTestLogger(io.Discard),
	)
	if err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := sndr.Start(ctx)

	tx1 := &sender.Transaction{
		Transaction: types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(100),
			21000,
			big.NewInt(1),
			nil,
		),
		Sender: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Type:   sender.TxTypeRegular,
		Raw:    "0x1234567890123456789012345678901234567890",
	}

	if err := st.AddBalance(ctx, tx1.Sender, big.NewInt(5e18)); err != nil {
		t.Fatalf("failed to add balance: %v", err)
	}

	if err := sndr.Enqueue(ctx, tx1); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	// Simulate opted in block
	bidderImpl.optinEstimate <- 2

	<-blockTracker.bnIn

	blockTracker.bnOut <- blockNoOp{
		block:             1,
		timeTillNextBlock: 5 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}

	// Simulate a bid response
	bidOp := <-bidderImpl.in
	if len(bidOp.rawTxs) != 1 || bidOp.rawTxs[0] != tx1.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %v", tx1.Raw, bidOp.rawTxs)
	}
	resC := make(chan bidder.BidStatus, 3)
	resC <- bidder.BidStatus{
		Type: bidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     1,
			ProviderAddress: "provider1",
		},
	}
	close(resC)
	bidderImpl.out <- bidResponse{statusCh: resC}

	bidderImpl.optinEstimate <- 2

	<-blockTracker.bnIn

	blockTracker.bnOut <- blockNoOp{
		block:             1,
		timeTillNextBlock: 2 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}

	bidOp = <-bidderImpl.in
	if len(bidOp.opts.IgnoreProviders) != 1 {
		t.Fatalf("expected 1 ignored provider, got %d", len(bidOp.opts.IgnoreProviders))
	}

	resC = make(chan bidder.BidStatus, 3)
	resC <- bidder.BidStatus{
		Type: bidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     1,
			ProviderAddress: "provider2",
		},
	}
	close(resC)
	bidderImpl.out <- bidResponse{statusCh: resC}
	res := <-st.preconfirmedTxns
	if res.txn == nil {
		t.Fatal("expected a preconfirmed transaction, got nil")
	}

	cancel()
	<-done
}

type senderTestEnv struct {
	st           *mockStore
	pricer       *mockPricer
	bidder       *mockBidder
	blockTracker *mockBlockTracker
	notifier     *mockNotifier
	sndr         *sender.TxSender
}

func newSenderTestEnv(t *testing.T, simulator *mockSimulator) *senderTestEnv {
	t.Helper()

	st := newMockStore()
	testPricer := &mockPricer{
		out: make(chan map[int64]float64, 20),
	}
	bidderImpl := &mockBidder{
		optinEstimate: make(chan int64, 20),
		in:            make(chan bidOp, 20),
		out:           make(chan bidResponse, 20),
	}
	blockTracker := &mockBlockTracker{
		out:   make(chan uint64, 20),
		bnIn:  make(chan struct{}, 20),
		bnOut: make(chan blockNoOp, 20),
		bnErr: make(chan error, 1),
	}
	notifier := &mockNotifier{}

	sndr, err := sender.NewTxSender(
		st,
		bidderImpl,
		testPricer,
		blockTracker,
		&mockTransferer{},
		notifier,
		simulator,
		&mockBackrunner{},
		big.NewInt(1), // Settlement chain ID
		&MockExplorerSubmitter{},
		nil,  // no log encryption key in tests
		true, // bundle bids enabled
		util.NewTestLogger(io.Discard),
	)
	if err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}

	return &senderTestEnv{
		st:           st,
		pricer:       testPricer,
		bidder:       bidderImpl,
		blockTracker: blockTracker,
		notifier:     notifier,
		sndr:         sndr,
	}
}

func newTestTxn(senderAddr common.Address, value int64, raw string, txType sender.TxType) *sender.Transaction {
	return newTestTxnWithGas(senderAddr, value, raw, 21000, txType)
}

func newTestTxnWithGas(senderAddr common.Address, value int64, raw string, gas uint64, txType sender.TxType) *sender.Transaction {
	return &sender.Transaction{
		Transaction: types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(value),
			gas,
			big.NewInt(1),
			nil,
		),
		Sender: senderAddr,
		Type:   txType,
		Raw:    raw,
	}
}

func newTestTxnWithGasPrice(senderAddr common.Address, value int64, raw string, gasPrice int64) *sender.Transaction {
	return &sender.Transaction{
		Transaction: types.NewTransaction(
			0,
			common.HexToAddress("0x1234567890123456789012345678901234567890"),
			big.NewInt(value),
			21000,
			big.NewInt(gasPrice),
			nil,
		),
		Sender: senderAddr,
		Type:   sender.TxTypeRegular,
		Raw:    raw,
	}
}

// bidCost is the cost of a transaction with the given gas limit at the
// 85 confidence price of 2.0 gwei.
func bidCost(gas int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(2_000_000_000), big.NewInt(gas))
}

func testPrices() map[int64]float64 {
	return map[int64]float64{
		70: 0.8,
		75: 1.0,
		80: 1.5,
		85: 2.0,
	}
}

// bundleShare is the cost of one 21000 gas transaction at the 85
// confidence price of 2.0 gwei.
var bundleShare = big.NewInt(42_000_000_000_000)

func TestBundleBid(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Distinct gas limits give every member a distinct cost, so a
	// share or attribution mix-up between members cannot hide behind
	// equal amounts.
	txns := []*sender.Transaction{
		newTestTxnWithGas(common.HexToAddress("0x1111111111111111111111111111111111111111"), 101, "0xaa01", 21000, sender.TxTypeRegular),
		newTestTxnWithGas(common.HexToAddress("0x2222222222222222222222222222222222222222"), 102, "0xaa02", 30000, sender.TxTypeRegular),
		newTestTxnWithGas(common.HexToAddress("0x3333333333333333333333333333333333333333"), 103, "0xaa03", 40000, sender.TxTypeRegular),
	}
	wantShare := map[common.Hash]*big.Int{
		txns[0].Hash(): bidCost(21000),
		txns[1].Hash(): bidCost(30000),
		txns[2].Hash(): bidCost(40000),
	}
	for _, txn := range txns {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// Pre-feed the single bundle attempt. Only the leader calls these.
	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 3 {
		t.Fatalf("expected 3 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}
	wantRaw := map[string]bool{"aa01": true, "aa02": true, "aa03": true}
	for _, r := range op.rawTxs {
		if !wantRaw[r] {
			t.Fatalf("unexpected raw transaction in bundle bid: %s", r)
		}
	}
	expectedTotal := big.NewInt(0)
	for _, share := range wantShare {
		expectedTotal = new(big.Int).Add(expectedTotal, share)
	}
	if op.bidAmount.Cmp(expectedTotal) != 0 {
		t.Fatalf("expected bundle bid amount %s, got %s", expectedTotal, op.bidAmount)
	}
	if op.slashAmount.Sign() != 0 {
		t.Fatalf("expected zero slash amount, got %s", op.slashAmount)
	}
	if op.opts.BlockNumber != 1 {
		t.Fatalf("expected block number 1, got %d", op.opts.BlockNumber)
	}
	if len(op.opts.RevertingTxHashes) != 3 {
		t.Fatalf("expected 3 reverting tx hashes, got %d", len(op.opts.RevertingTxHashes))
	}
	wantHashes := map[string]bool{}
	for _, txn := range txns {
		wantHashes[txn.Hash().Hex()] = true
	}
	for _, h := range op.opts.RevertingTxHashes {
		if !wantHashes[h] {
			t.Fatalf("unexpected reverting tx hash: %s", h)
		}
	}
	if op.opts.Constraint != nil {
		t.Fatal("expected no constraint on bundle bid")
	}
	if len(op.opts.IgnoreProviders) != 0 {
		t.Fatalf("expected no ignored providers, got %v", op.opts.IgnoreProviders)
	}

	resC := make(chan bidder.BidStatus, 3)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       expectedTotal.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	seen := map[common.Hash]result{}
	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		seen[res.txn.Hash()] = res
	}
	for _, txn := range txns {
		res, found := seen[txn.Hash()]
		if !found {
			t.Fatalf("missing preconfirmed result for %s", txn.Hash().Hex())
		}
		if res.blockNumber != 1 {
			t.Fatalf("expected block number 1, got %d", res.blockNumber)
		}
		if len(res.commitments) != 2 {
			t.Fatalf("expected 2 commitments, got %d", len(res.commitments))
		}
		for _, cmt := range res.commitments {
			if cmt.BidAmount != wantShare[txn.Hash()].String() {
				t.Fatalf(
					"expected member share %s for %s, got %s",
					wantShare[txn.Hash()], txn.Hash().Hex(), cmt.BidAmount,
				)
			}
		}
	}

	// Let all members observe inclusion and finish.
	for i := 0; i < 3; i++ {
		env.blockTracker.out <- 1
	}

	cancel()
	<-done
}

func TestBundleFailureSoloFallback(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	txns := []*sender.Transaction{
		newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 201, "0xbb01", sender.TxTypeRegular),
		newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 202, "0xbb02", sender.TxTypeRegular),
		newTestTxn(common.HexToAddress("0x3333333333333333333333333333333333333333"), 203, "0xbb03", sender.TxTypeRegular),
	}
	for _, txn := range txns {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// Pre-feed one bundle attempt and three solo retries.
	for i := 0; i < 4; i++ {
		env.bidder.optinEstimate <- 7
		env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
		env.pricer.out <- testPrices()
	}

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 3 {
		t.Fatalf("expected 3 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}
	env.bidder.out <- bidResponse{err: errors.New("bundle rejected")}

	soloSeen := map[string]bool{}
	for i := 0; i < 3; i++ {
		op := <-env.bidder.in
		if len(op.rawTxs) != 1 {
			t.Fatalf("expected solo bid after bundle failure, got %d raw transactions", len(op.rawTxs))
		}
		soloSeen[op.rawTxs[0]] = true

		resC := make(chan bidder.BidStatus, 3)
		for _, provider := range []string{"provider1", "provider2"} {
			resC <- bidder.BidStatus{
				Type: bidder.BidStatusCommitment,
				Arg: &bidderapiv1.Commitment{
					BidAmount:       bundleShare.String(),
					BlockNumber:     1,
					ProviderAddress: provider,
				},
			}
		}
		close(resC)
		env.bidder.out <- bidResponse{statusCh: resC}
	}
	if len(soloSeen) != 3 {
		t.Fatalf("expected 3 distinct solo bids, got %d", len(soloSeen))
	}

	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		if res.txn.Status != sender.TxStatusPreConfirmed {
			t.Fatalf("expected preconfirmed transaction, got %s", res.txn.Status)
		}
	}

	for i := 0; i < 3; i++ {
		env.blockTracker.out <- 1
	}

	cancel()
	<-done
}

func TestBundleMemberSimFailureDropped(t *testing.T) {
	t.Parallel()

	badTxn := newTestTxn(common.HexToAddress("0x3333333333333333333333333333333333333333"), 303, "0xcc03", sender.TxTypeRegular)
	env := newSenderTestEnv(t, &mockSimulator{failRaw: badTxn.Raw})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	txns := []*sender.Transaction{
		newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 301, "0xcc01", sender.TxTypeRegular),
		newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 302, "0xcc02", sender.TxTypeRegular),
		badTxn,
	}
	for _, txn := range txns {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected 2 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}
	for _, r := range op.rawTxs {
		if r == "cc03" {
			t.Fatal("failed member must not appear in the bundle bid")
		}
	}
	expectedTotal := new(big.Int).Mul(bundleShare, big.NewInt(2))
	if op.bidAmount.Cmp(expectedTotal) != 0 {
		t.Fatalf("expected bundle bid amount %s, got %s", expectedTotal, op.bidAmount)
	}

	resC := make(chan bidder.BidStatus, 3)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				BidAmount:       expectedTotal.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	preconfirmed := 0
	failed := 0
	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		switch res.txn.Status {
		case sender.TxStatusPreConfirmed:
			preconfirmed++
			if len(res.commitments) != 2 {
				t.Fatalf("expected 2 commitments, got %d", len(res.commitments))
			}
			for _, cmt := range res.commitments {
				if cmt.BidAmount != bundleShare.String() {
					t.Fatalf("expected member share %s, got %s", bundleShare, cmt.BidAmount)
				}
			}
		case sender.TxStatusFailed:
			failed++
			if res.txn.Hash() != badTxn.Hash() {
				t.Fatalf("expected failed transaction %s, got %s", badTxn.Hash().Hex(), res.txn.Hash().Hex())
			}
		default:
			t.Fatalf("unexpected transaction status %s", res.txn.Status)
		}
	}
	if preconfirmed != 2 || failed != 1 {
		t.Fatalf("expected 2 preconfirmed and 1 failed, got %d and %d", preconfirmed, failed)
	}

	for i := 0; i < 3; i++ {
		env.blockTracker.out <- 1
	}

	cancel()
	<-done
}

// TestBundleDispatchAllOrNothing checks that a batch dispatch is all or
// nothing: when the context ends while the batch waits for its second
// worker slot, no member goroutine starts, no bid goes out, and every
// inflight mark is released so the next pass forms the full bundle
// again. The worker pool keeps exactly one free slot, so the old
// per-member dispatch would start the leader goroutine with that slot
// and place a bundle bid from the pre-fed inputs below.
func TestBundleDispatchAllOrNothing(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx := context.Background()

	txns := []*sender.Transaction{
		newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 601, "0xff01", sender.TxTypeRegular),
		newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 602, "0xff02", sender.TxTypeRegular),
	}
	for _, txn := range txns {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// Leave exactly one free worker slot. The two-member batch needs
	// two slots, so the all-or-nothing acquisition blocks on the
	// second slot until the dispatch context ends.
	drain := env.sndr.FillWorkerPoolForTest(1)
	defer drain()

	// One full round of inputs. The all-or-nothing dispatch starts no
	// goroutine, so nothing consumes them here; the follow-up pass
	// does. A leader goroutine started by a per-member dispatch would
	// consume them and place a bundle bid.
	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	dispatchCtx, cancelDispatch := context.WithCancel(context.Background())
	defer cancelDispatch()
	timer := time.AfterFunc(100*time.Millisecond, cancelDispatch)
	defer timer.Stop()

	env.sndr.ProcessQueuedForTest(dispatchCtx)

	select {
	case op := <-env.bidder.in:
		t.Fatalf("expected no bid after blocked dispatch, got bid with %d raw transactions", len(op.rawTxs))
	default:
	}

	// Every inflight mark is released, so the next pass with a free
	// worker pool forms the full bundle from the pre-fed inputs.
	drain()

	env.sndr.ProcessQueuedForTest(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected the full bundle on the next pass, got %d raw transactions", len(op.rawTxs))
	}

	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       op.bidAmount.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	for i := 0; i < 2; i++ {
		res := <-env.st.preconfirmedTxns
		if res.txn.Status != sender.TxStatusPreConfirmed {
			t.Fatalf("expected preconfirmed transaction, got %s", res.txn.Status)
		}
	}

	for i := 0; i < 2; i++ {
		env.blockTracker.out <- 1
	}

	if err := env.sndr.WaitWorkersForTest(); err != nil {
		t.Fatalf("worker goroutines failed: %v", err)
	}
}

// TestBundleAbortSoloRetryFirstAttempt checks that a bundle abort with
// less than two survivors does not consume first-attempt semantics: the
// survivor's solo bid sets a real noOfProviders, so a zero-commitment
// outcome is not marked preconfirmed.
func TestBundleAbortSoloRetryFirstAttempt(t *testing.T) {
	t.Parallel()

	badTxn := newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 702, "0xab02", sender.TxTypeRegular)
	env := newSenderTestEnv(t, &mockSimulator{failRaw: badTxn.Raw})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	survivor := newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 701, "0xab01", sender.TxTypeRegular)
	for _, txn := range []*sender.Transaction{survivor, badTxn} {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// One bundle attempt that aborts before the bid, then the
	// survivor's solo attempt.
	for i := 0; i < 2; i++ {
		env.bidder.optinEstimate <- 7
		env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
		env.pricer.out <- testPrices()
	}

	done := env.sndr.Start(ctx)

	// The bad member fails simulation, the bundle collapses to one
	// member and aborts before the bid. The only bid on the wire is
	// the survivor's solo bid.
	op := <-env.bidder.in
	if len(op.rawTxs) != 1 || op.rawTxs[0] != "ab01" {
		t.Fatalf("expected solo bid for the survivor, got %v", op.rawTxs)
	}

	// Both members registered an inclusion waiter; feed one signal
	// for each before the bid response, so the survivor's loop can
	// end through the inclusion path right after the empty bid.
	for i := 0; i < 2; i++ {
		env.blockTracker.out <- 1
	}

	// The solo bid resolves with zero commitments.
	resC := make(chan bidder.BidStatus)
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	for i := 0; i < 2; i++ {
		res := <-env.st.preconfirmedTxns
		switch res.txn.Hash() {
		case survivor.Hash():
			if res.txn.Status == sender.TxStatusPreConfirmed {
				t.Fatalf("survivor preconfirmed with %d commitments after aborted bundle", len(res.commitments))
			}
			if res.txn.Status != sender.TxStatusConfirmed {
				t.Fatalf("expected confirmed survivor, got %s", res.txn.Status)
			}
		case badTxn.Hash():
			if res.txn.Status != sender.TxStatusFailed {
				t.Fatalf("expected failed member, got %s", res.txn.Status)
			}
		default:
			t.Fatalf("unexpected stored transaction %s", res.txn.Hash().Hex())
		}
	}

	cancel()
	<-done
}

// TestBundleCtxCancelKeepsFastTracked checks that a context end during
// the bundle bid status loop does not revert members that were already
// fast-tracked: they keep their stored preconfirmation, and only the
// members that are not resolved get the error path. The store honors
// the cancelled context, like a real database, so a repeat store on
// the dead context would fail: the fast-tracked member must not repeat
// the store the leader already did, and must not signal a second
// receipt. The unresolved member fails visibly with a terminal write
// that survives the cancellation.
func TestBundleCtxCancelKeepsFastTracked(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})
	env.st.honorCtx = true

	// Distinct gas limits give distinct member costs, so the fast
	// track function can single out member A by its own share.
	txnA := newTestTxnWithGas(common.HexToAddress("0x1111111111111111111111111111111111111111"), 801, "0xac01", 21000, sender.TxTypeRegular)
	txnB := newTestTxnWithGas(common.HexToAddress("0x2222222222222222222222222222222222222222"), 802, "0xac02", 42000, sender.TxTypeRegular)
	// The fast track fires for member A on the second commitment.
	// Reading A's fast-track store write below then proves that the
	// status loop consumed both commitments before the context ends.
	costA := bidCost(21000)
	env.sndr.SetFastTrackFunc(func(cmts []*bidderapiv1.Commitment, _ bool) bool {
		return len(cmts) >= 2 && cmts[0].BidAmount == costA.String()
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, txn := range []*sender.Transaction{txnA, txnB} {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	waitA := env.sndr.WaitForReceiptAvailable(ctx, txnA.Hash())

	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected 2 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}

	// Both commitments arrive; the first fast-tracks member A. The
	// channel stays open, so the status loop still runs when the
	// context ends.
	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       op.bidAmount.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	env.bidder.out <- bidResponse{statusCh: resC}

	res := <-env.st.preconfirmedTxns
	if res.txn.Hash() != txnA.Hash() {
		t.Fatalf("expected fast-tracked store write for member A, got %s", res.txn.Hash().Hex())
	}
	if res.txn.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("expected preconfirmed fast-tracked member, got %s", res.txn.Status)
	}
	if len(res.commitments) != 2 {
		t.Fatalf("expected 2 commitments on fast track, got %d", len(res.commitments))
	}
	for _, cmt := range res.commitments {
		if cmt.BidAmount != costA.String() {
			t.Fatalf("expected member A share %s, got %s", costA, cmt.BidAmount)
		}
	}

	// The fast track signaled member A's receipt once. A second
	// waiter registered now must stay open: a duplicate receipt
	// signal after the cancel would close it.
	<-waitA
	dupA := env.sndr.WaitForReceiptAvailable(ctx, txnA.Hash())

	cancel()
	<-done

	// After shutdown member A keeps its preconfirmation: the leader
	// already stored and signaled it at fast-track time, and no
	// commitment arrived after that write, so the success gate
	// repeats no store write. Member B, which is not resolved, fails
	// visibly; its terminal failed write runs on a context that
	// survives the cancellation, so it lands even though the store
	// honors context cancellation.
	var post []result
	drained := false
	for !drained {
		select {
		case r := <-env.st.preconfirmedTxns:
			post = append(post, r)
		default:
			drained = true
		}
	}
	if len(post) != 1 {
		t.Fatalf("expected one store write after cancel, got %d", len(post))
	}
	if post[0].txn.Hash() != txnB.Hash() || post[0].txn.Status != sender.TxStatusFailed {
		t.Fatalf(
			"expected the failed store write for member B, got %s with status %s",
			post[0].txn.Hash().Hex(), post[0].txn.Status,
		)
	}
	if txnA.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("fast-tracked member status reverted to %s", txnA.Status)
	}
	if txnB.Status != sender.TxStatusFailed {
		t.Fatalf("expected failed unresolved member B, got %s", txnB.Status)
	}
	select {
	case <-dupA:
		t.Fatal("duplicate receipt signal for fast-tracked member")
	default:
	}
}

// TestDepositCtxCancelFailsVisibly checks that a shutdown does not
// mask a preconfirmed deposit as a success. The balance credit runs
// only after the bid loop, and the queue only selects pending
// transactions again. So the deposit must end visibly failed, with no
// balance credit, and the user can retry it.
func TestDepositCtxCancelFailsVisibly(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})
	// The store honors context cancellation, like a real database.
	// The terminal failed write must land anyway: the sender writes
	// it on a context that survives the cancellation.
	env.st.honorCtx = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	depositTxn := newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 1e18, "0xdd11", sender.TxTypeDeposit)
	if err := env.sndr.Enqueue(ctx, depositTxn); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	// Opted-in slot, so the deposit preconfirms after the bid.
	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 1 || op.rawTxs[0] != "dd11" {
		t.Fatalf("expected solo deposit bid, got %v", op.rawTxs)
	}

	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        []string{depositTxn.Hash().Hex()},
				BidAmount:       bundleShare.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	res := <-env.st.preconfirmedTxns
	if res.txn.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("expected preconfirmed deposit, got %s", res.txn.Status)
	}

	// The context ends before inclusion. The deposit has not credited
	// its balance yet, so it must not end as a silent preconfirmed
	// success that the queue never selects again.
	cancel()
	<-done

	var writes []result
	drained := false
	for !drained {
		select {
		case r := <-env.st.preconfirmedTxns:
			writes = append(writes, r)
		default:
			drained = true
		}
	}
	if len(writes) != 1 {
		t.Fatalf("expected one store write after cancel, got %d", len(writes))
	}
	if writes[0].txn.Status != sender.TxStatusFailed {
		t.Fatalf("expected failed deposit after cancel, got %s", writes[0].txn.Status)
	}
	if depositTxn.Status != sender.TxStatusFailed {
		t.Fatalf("expected failed deposit status, got %s", depositTxn.Status)
	}
	if count := env.st.addBalanceCount(depositTxn.Sender); count != 0 {
		t.Fatalf("expected no balance credit for the cancelled deposit, got %d", count)
	}
}

// TestFastSwapCtxCancelKeepsPreconfirmed checks that a shutdown does
// not revert a preconfirmed fastswap to a failed state. A fastswap has
// no post-loop obligation, so the stored preconfirmation stands: no
// failed store write and no failure notification.
func TestFastSwapCtxCancelKeepsPreconfirmed(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})
	env.st.honorCtx = true

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	swapTxn := newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 501, "0xfa01", sender.TxTypeFastSwap)
	if err := env.sndr.Enqueue(ctx, swapTxn); err != nil {
		t.Fatalf("failed to enqueue transaction: %v", err)
	}

	// Opted-in slot, so the fastswap preconfirms after the bid.
	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 1 || op.rawTxs[0] != "fa01" {
		t.Fatalf("expected solo fastswap bid, got %v", op.rawTxs)
	}

	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        []string{swapTxn.Hash().Hex()},
				BidAmount:       bundleShare.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	res := <-env.st.preconfirmedTxns
	if res.txn.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("expected preconfirmed fastswap, got %s", res.txn.Status)
	}

	// The context ends before inclusion. The fastswap has no
	// post-loop work, so the shutdown must not revert the stored
	// preconfirmation to a failed state.
	cancel()
	<-done

	select {
	case r := <-env.st.preconfirmedTxns:
		t.Fatalf("unexpected store write after cancel with status %s", r.txn.Status)
	default:
	}
	if swapTxn.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("expected preconfirmed fastswap after cancel, got %s", swapTxn.Status)
	}
	if env.notifier.notified(swapTxn.Hash()) {
		t.Fatal("unexpected failure notification for the preconfirmed fastswap")
	}
}

// TestBundleLateCommitmentsRestored checks bundle and solo parity for
// commitments that arrive after the fast-track write: the fast-track
// write persists the early set, later commitments append in memory
// only, and the success gate stores the grown set again, on a context
// that survives the shutdown. The unresolved member still fails
// visibly with a terminal write that also survives the shutdown.
func TestBundleLateCommitmentsRestored(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})
	env.st.honorCtx = true

	// Distinct gas limits give distinct member costs, so the fast
	// track function can single out member A by its own share.
	txnA := newTestTxnWithGas(common.HexToAddress("0x1111111111111111111111111111111111111111"), 811, "0xad01", 21000, sender.TxTypeRegular)
	txnB := newTestTxnWithGas(common.HexToAddress("0x2222222222222222222222222222222222222222"), 812, "0xad02", 42000, sender.TxTypeRegular)

	// The fast track fires for member A on the first commitment, so
	// the fast-track write persists one commitment and the second
	// commitment exists in memory only. The channel signals when the
	// second commitment reached member A.
	costA := bidCost(21000)
	secondSeen := make(chan struct{})
	var secondOnce sync.Once
	env.sndr.SetFastTrackFunc(func(cmts []*bidderapiv1.Commitment, _ bool) bool {
		if len(cmts) == 0 || cmts[0].BidAmount != costA.String() {
			return false
		}
		if len(cmts) >= 2 {
			secondOnce.Do(func() { close(secondSeen) })
		}
		return true
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, txn := range []*sender.Transaction{txnA, txnB} {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected 2 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}

	// Both commitments arrive; the first fast-tracks member A. The
	// channel stays open, so the status loop still runs when the
	// context ends.
	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       op.bidAmount.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	env.bidder.out <- bidResponse{statusCh: resC}

	// The fast-track write persists exactly one commitment.
	res := <-env.st.preconfirmedTxns
	if res.txn.Hash() != txnA.Hash() {
		t.Fatalf("expected fast-track store write for member A, got %s", res.txn.Hash().Hex())
	}
	if res.txn.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("expected preconfirmed fast-tracked member, got %s", res.txn.Status)
	}
	if len(res.commitments) != 1 {
		t.Fatalf("expected 1 commitment on the fast-track write, got %d", len(res.commitments))
	}

	// The second commitment reached member A in memory only. The
	// shutdown ends the status loop.
	<-secondSeen
	cancel()
	<-done

	var writes []result
	drained := false
	for !drained {
		select {
		case r := <-env.st.preconfirmedTxns:
			writes = append(writes, r)
		default:
			drained = true
		}
	}
	if len(writes) != 2 {
		t.Fatalf("expected two store writes after the fast-track write, got %d", len(writes))
	}
	for _, w := range writes {
		switch w.txn.Hash() {
		case txnA.Hash():
			if w.txn.Status != sender.TxStatusPreConfirmed {
				t.Fatalf("expected preconfirmed re-store for member A, got %s", w.txn.Status)
			}
			if len(w.commitments) != 2 {
				t.Fatalf("expected 2 commitments on the re-store, got %d", len(w.commitments))
			}
			for _, cmt := range w.commitments {
				if cmt.BidAmount != costA.String() {
					t.Fatalf("expected member A share %s, got %s", costA, cmt.BidAmount)
				}
			}
		case txnB.Hash():
			if w.txn.Status != sender.TxStatusFailed {
				t.Fatalf("expected failed unresolved member B, got %s", w.txn.Status)
			}
		default:
			t.Fatalf("unexpected store write for %s", w.txn.Hash().Hex())
		}
	}
	if txnA.Status != sender.TxStatusPreConfirmed {
		t.Fatalf("fast-tracked member status reverted to %s", txnA.Status)
	}
}

// TestBundleBalanceDropRecordsHistory checks that a terminal pre-check
// drop whose only block attempt was rolled back still notifies the
// user and records the hash as historical: a later queue pass must not
// process the same hash again.
func TestBundleBalanceDropRecordsHistory(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Three members, so the bundle survives the drop: member B has no
	// balance and drops terminally; members A and C bid.
	txnA := newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 951, "0xba01", sender.TxTypeRegular)
	txnB := newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 952, "0xba02", sender.TxTypeRegular)
	txnC := newTestTxn(common.HexToAddress("0x3333333333333333333333333333333333333333"), 953, "0xba03", sender.TxTypeRegular)

	inclA := env.blockTracker.inclusionFor(txnA.Hash())
	env.blockTracker.inclusionFor(txnB.Hash())
	inclC := env.blockTracker.inclusionFor(txnC.Hash())

	for _, txn := range []*sender.Transaction{txnA, txnC} {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
	}
	for _, txn := range []*sender.Transaction{txnA, txnB, txnC} {
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	env.sndr.ProcessQueuedForTest(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected 2 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}
	for _, r := range op.rawTxs {
		if r == "ba02" {
			t.Fatal("member without balance must not appear in the bundle bid")
		}
	}

	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       op.bidAmount.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	inclA <- 1
	inclC <- 1

	if err := env.sndr.WaitWorkersForTest(); err != nil {
		t.Fatalf("worker goroutines failed: %v", err)
	}

	failed := 0
	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		if res.txn.Hash() == txnB.Hash() {
			if res.txn.Status != sender.TxStatusFailed {
				t.Fatalf("expected failed dropped member, got %s", res.txn.Status)
			}
			failed++
		}
	}
	if failed != 1 {
		t.Fatalf("expected one failed store write for the dropped member, got %d", failed)
	}

	// The drop rolled its only block attempt back, so the attempt
	// entry is gone. The user notification must still happen.
	if !env.notifier.notified(txnB.Hash()) {
		t.Fatal("expected a notification for the dropped member")
	}

	// The historical record keeps a later queue pass from processing
	// the same hash again: no bid goes out and no store write runs.
	if err := env.st.AddQueuedTransaction(ctx, txnB); err != nil {
		t.Fatalf("failed to re-queue transaction: %v", err)
	}
	env.sndr.ProcessQueuedForTest(ctx)
	if err := env.sndr.WaitWorkersForTest(); err != nil {
		t.Fatalf("worker goroutines failed: %v", err)
	}
	select {
	case op := <-env.bidder.in:
		t.Fatalf("historical transaction was bid again: %v", op.rawTxs)
	default:
	}
	select {
	case r := <-env.st.preconfirmedTxns:
		t.Fatalf("unexpected store write for the historical transaction with status %s", r.txn.Status)
	default:
	}
}

// TestBundlePreCheckDropSoloRetryFirstAttempt checks that a member
// dropped by a bundle pre-check before the bid keeps first-attempt
// semantics: the drop rolls its block attempt back, so the solo retry
// sets a real provider count and a zero-commitment outcome is not
// marked preconfirmed.
func TestBundlePreCheckDropSoloRetryFirstAttempt(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	txnA := newTestTxnWithGasPrice(common.HexToAddress("0x1111111111111111111111111111111111111111"), 901, "0xae01", 10)
	txnB := newTestTxnWithGasPrice(common.HexToAddress("0x2222222222222222222222222222222222222222"), 902, "0xae02", 1)
	txnC := newTestTxnWithGasPrice(common.HexToAddress("0x3333333333333333333333333333333333333333"), 903, "0xae03", 10)

	// The next base fee is above member B's fee per gas, so the fee
	// pre-check drops member B from the bundle before the bid.
	env.blockTracker.setNextBaseFee(big.NewInt(5))

	inclA := env.blockTracker.inclusionFor(txnA.Hash())
	inclB := env.blockTracker.inclusionFor(txnB.Hash())
	inclC := env.blockTracker.inclusionFor(txnC.Hash())

	for _, txn := range []*sender.Transaction{txnA, txnB, txnC} {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// Bundle round: not opted in, with a short next block time so
	// member B's solo retry runs after about one second.
	env.bidder.optinEstimate <- 20
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: time.Second}
	env.pricer.out <- testPrices()

	done := env.sndr.Start(ctx)

	op := <-env.bidder.in
	if len(op.rawTxs) != 2 {
		t.Fatalf("expected 2 raw transactions in bundle bid, got %d", len(op.rawTxs))
	}
	for _, r := range op.rawTxs {
		if r == "ae02" {
			t.Fatal("dropped member must not appear in the bundle bid")
		}
	}

	// Clear the base fee so member B's solo retry passes the fee
	// check, and feed the solo round: opted in, so a zero-commitment
	// retry that wrongly compares zero providers with zero
	// commitments would be stored as preconfirmed.
	env.blockTracker.setNextBaseFee(nil)
	env.bidder.optinEstimate <- 7
	env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
	env.pricer.out <- testPrices()

	// The bundle members A and C get both commitments. The slot is
	// not opted in, so they end through the inclusion path.
	resC := make(chan bidder.BidStatus, 2)
	for _, provider := range []string{"provider1", "provider2"} {
		resC <- bidder.BidStatus{
			Type: bidder.BidStatusCommitment,
			Arg: &bidderapiv1.Commitment{
				TxHashes:        op.opts.RevertingTxHashes,
				BidAmount:       op.bidAmount.String(),
				BlockNumber:     1,
				ProviderAddress: provider,
			},
		}
	}
	close(resC)
	env.bidder.out <- bidResponse{statusCh: resC}

	inclA <- 1
	inclC <- 1

	// Member B retries solo after the fee drop.
	opB := <-env.bidder.in
	if len(opB.rawTxs) != 1 || opB.rawTxs[0] != "ae02" {
		t.Fatalf("expected solo bid for the dropped member, got %v", opB.rawTxs)
	}

	// The solo retry resolves with zero commitments.
	resB := make(chan bidder.BidStatus)
	close(resB)
	env.bidder.out <- bidResponse{statusCh: resB}

	inclB <- 1

	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		if res.txn.Hash() == txnB.Hash() && res.txn.Status == sender.TxStatusPreConfirmed {
			t.Fatalf("dropped member preconfirmed with %d commitments", len(res.commitments))
		}
		if res.txn.Status != sender.TxStatusConfirmed {
			t.Fatalf("expected confirmed transaction %s, got %s", res.txn.Hash().Hex(), res.txn.Status)
		}
	}

	cancel()
	<-done
}

func TestIneligibleTxnsNeverBundle(t *testing.T) {
	t.Parallel()

	env := newSenderTestEnv(t, &mockSimulator{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	regularTxn := newTestTxn(common.HexToAddress("0x1111111111111111111111111111111111111111"), 401, "0xdd01", sender.TxTypeRegular)
	depositTxn := newTestTxn(common.HexToAddress("0x2222222222222222222222222222222222222222"), 1e18, "0xdd02", sender.TxTypeDeposit)
	constraintTxn := newTestTxn(common.HexToAddress("0x3333333333333333333333333333333333333333"), 403, "0xdd03", sender.TxTypeRegular)
	constraintTxn.Constraint = &bidderapiv1.PositionConstraint{
		Anchor: bidderapiv1.PositionConstraint_ANCHOR_TOP,
		Basis:  bidderapiv1.PositionConstraint_BASIS_PERCENTILE,
		Value:  10,
	}

	txns := []*sender.Transaction{regularTxn, depositTxn, constraintTxn}
	for _, txn := range txns {
		if err := env.st.AddBalance(ctx, txn.Sender, big.NewInt(5e18)); err != nil {
			t.Fatalf("failed to add balance: %v", err)
		}
		if err := env.sndr.Enqueue(ctx, txn); err != nil {
			t.Fatalf("failed to enqueue transaction: %v", err)
		}
	}

	// Pre-feed three solo attempts. No bundle forms because only one
	// transaction is eligible.
	for i := 0; i < 3; i++ {
		env.bidder.optinEstimate <- 7
		env.blockTracker.bnOut <- blockNoOp{block: 1, timeTillNextBlock: 5 * time.Second}
		env.pricer.out <- testPrices()
	}

	done := env.sndr.Start(ctx)

	for i := 0; i < 3; i++ {
		op := <-env.bidder.in
		if len(op.rawTxs) != 1 {
			t.Fatalf("expected solo bid, got %d raw transactions", len(op.rawTxs))
		}
		if op.rawTxs[0] == "dd03" && op.opts.Constraint == nil {
			t.Fatal("expected constraint on constrained transaction bid")
		}

		resC := make(chan bidder.BidStatus, 3)
		for _, provider := range []string{"provider1", "provider2"} {
			resC <- bidder.BidStatus{
				Type: bidder.BidStatusCommitment,
				Arg: &bidderapiv1.Commitment{
					BidAmount:       bundleShare.String(),
					BlockNumber:     1,
					ProviderAddress: provider,
				},
			}
		}
		close(resC)
		env.bidder.out <- bidResponse{statusCh: resC}
	}

	for i := 0; i < 3; i++ {
		res := <-env.st.preconfirmedTxns
		if res.txn.Status != sender.TxStatusPreConfirmed {
			t.Fatalf("expected preconfirmed transaction, got %s", res.txn.Status)
		}
	}

	for i := 0; i < 3; i++ {
		env.blockTracker.out <- 1
	}

	cancel()
	<-done
}
