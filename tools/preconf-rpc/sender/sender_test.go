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
}

func newMockStore() *mockStore {
	return &mockStore{
		queued:           make(map[common.Address][]*sender.Transaction),
		nonce:            make(map[common.Address]uint64),
		balances:         make(map[common.Address]*big.Int),
		preconfirmedTxns: make(chan result, 10),
		byHash:           make(map[common.Hash]*sender.Transaction),
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

	if _, exists := m.balances[account]; !exists {
		m.balances[account] = amount
	} else {
		newBalance := new(big.Int).Add(m.balances[account], amount)
		m.balances[account] = newBalance
	}

	return nil
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
	rawTx       string
	opts        *bidder.BidOpts
}

type mockBidder struct {
	optinEstimate chan int64
	in            chan bidOp
	out           chan chan bidder.BidStatus
}

func (m *mockBidder) Estimate() (int64, error) {
	estimate := <-m.optinEstimate
	return estimate, nil
}

func (m *mockBidder) Bid(
	ctx context.Context,
	bidAmount *big.Int,
	slashAmount *big.Int,
	rawTx string,
	opts *bidder.BidOpts,
) (chan bidder.BidStatus, error) {
	m.in <- bidOp{
		bidAmount:   bidAmount,
		slashAmount: slashAmount,
		rawTx:       rawTx,
		opts:        opts,
	}
	res := <-m.out

	return res, nil
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
	lbf   *big.Int
	nbf   *big.Int
}

func (m *mockBlockTracker) WaitForTxnInclusion(txnHash common.Hash) chan uint64 {
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
	if m.lbf == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(m.lbf)
}

func (m *mockBlockTracker) NextBaseFee() *big.Int {
	if m.nbf == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(m.nbf)
}

type mockTransferer struct{}

func (m *mockTransferer) Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error {
	return nil
}

type mockNotifier struct {
	notifications []string
}

func (m *mockNotifier) NotifyTransactionStatus(txn *sender.Transaction, attempts, blocks int, start time.Duration) {
	m.notifications = append(m.notifications, txn.Hash().Hex())
}

type mockSimulator struct{}

func (m *mockSimulator) Simulate(ctx context.Context, rawTx string, _ sim.SimState) ([]*types.Log, bool, error) {
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
		out:           make(chan chan bidder.BidStatus, 10),
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
		90: 1.0,
		95: 1.5,
		99: 2.0,
	}

	// Simulate transaction inclusion
	blockTracker.out <- 1

	// Simulate a bid response
	bidOp := <-bidderImpl.in
	if bidOp.rawTx != tx1.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
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
	bidderImpl.out <- resC

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
		90: 1.0,
		95: 1.5,
		99: 2.0,
	}

	// Simulate a bid response
	bidOp = <-bidderImpl.in
	if bidOp.rawTx != tx2.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
	}
	resC = make(chan bidder.BidStatus, 3)
	// Simulate retry due to incomplete commitments
	close(resC)
	bidderImpl.out <- resC

	// Simulate non opted in block
	bidderImpl.optinEstimate <- 18

	<-blockTracker.bnIn
	blockTracker.bnOut <- blockNoOp{
		block:             2,
		timeTillNextBlock: 5 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		90: 1.0,
		95: 1.5,
		99: 2.0,
	}

	// Simulate transaction inclusion
	blockTracker.out <- 2

	// Simulate a bid response
	bidOp = <-bidderImpl.in
	if bidOp.rawTx != tx2.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
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
	bidderImpl.out <- resC

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
		out:           make(chan chan bidder.BidStatus, 10),
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
		out:           make(chan chan bidder.BidStatus, 10),
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
		90: 1.0,
		95: 1.5,
		99: 2.0,
	}

	// Simulate a bid response
	bidOp := <-bidderImpl.in
	if bidOp.rawTx != tx1.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
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
	bidderImpl.out <- resC

	bidderImpl.optinEstimate <- 2

	<-blockTracker.bnIn

	blockTracker.bnOut <- blockNoOp{
		block:             1,
		timeTillNextBlock: 2 * time.Second,
	}

	// Simulate a price estimate
	testPricer.out <- map[int64]float64{
		90: 1.0,
		95: 1.5,
		99: 2.0,
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
	bidderImpl.out <- resC
	res := <-st.preconfirmedTxns
	if res.txn == nil {
		t.Fatal("expected a preconfirmed transaction, got nil")
	}

	cancel()
	<-done
}
