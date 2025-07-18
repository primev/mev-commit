package sender_test

import (
	"context"
	"errors"
	"math/big"
	"os"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderapiv1 "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/primev/mev-commit/tools/preconf-rpc/pricer"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
	optinbidder "github.com/primev/mev-commit/x/opt-in-bidder"
	"github.com/primev/mev-commit/x/util"
)

type result struct {
	txn         *sender.Transaction
	commitments []*bidderapiv1.Commitment
	blockNumber int64
}

type mockStore struct {
	mu               sync.Mutex
	queued           map[common.Address][]*sender.Transaction
	nonce            map[common.Address]uint64
	balances         map[common.Address]*big.Int
	preconfirmedTxns chan result
}

func newMockStore() *mockStore {
	return &mockStore{
		queued:           make(map[common.Address][]*sender.Transaction),
		nonce:            make(map[common.Address]uint64),
		balances:         make(map[common.Address]*big.Int),
		preconfirmedTxns: make(chan result, 10),
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

func (m *mockStore) GetCurrentNonce(_ context.Context, sender common.Address) uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	nonce, exists := m.nonce[sender]
	if !exists {
		return 0
	}

	return nonce
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
) error {
	m.preconfirmedTxns <- result{
		txn:         txn,
		commitments: commitments,
		blockNumber: txn.BlockNumber,
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
	return nil
}

type bidOp struct {
	bidAmount   *big.Int
	slashAmount *big.Int
	rawTx       string
	opts        *optinbidder.BidOpts
}

type mockBidder struct {
	optinEstimate chan int64
	in            chan bidOp
	out           chan chan optinbidder.BidStatus
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
	opts *optinbidder.BidOpts,
) (chan optinbidder.BidStatus, error) {
	m.in <- bidOp{
		bidAmount:   bidAmount,
		slashAmount: slashAmount,
		rawTx:       rawTx,
		opts:        opts,
	}
	res := <-m.out

	return res, nil
}

type mockPricer struct {
	out chan *pricer.BlockPrices
}

func (m *mockPricer) EstimatePrice(
	ctx context.Context,
) (*pricer.BlockPrices, error) {
	select {
	case prices := <-m.out:
		if prices == nil {
			return nil, errors.New("nil price returned")
		}
		return prices, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type op struct {
	hash  common.Hash
	block uint64
}

type mockBlockTracker struct {
	in  chan op
	out chan bool
}

func (m *mockBlockTracker) CheckTxnInclusion(ctx context.Context, txnHash common.Hash, blockNumber uint64) (bool, error) {
	m.in <- op{
		hash:  txnHash,
		block: blockNumber,
	}
	select {
	case included := <-m.out:
		return included, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

type mockTransferer struct{}

func (m *mockTransferer) Transfer(ctx context.Context, to common.Address, chainID *big.Int, amount *big.Int) error {
	return nil
}

func TestSender(t *testing.T) {
	t.Parallel()

	st := newMockStore()
	testPricer := &mockPricer{
		out: make(chan *pricer.BlockPrices, 10),
	}
	bidder := &mockBidder{
		optinEstimate: make(chan int64, 10),
		in:            make(chan bidOp, 10),
		out:           make(chan chan optinbidder.BidStatus, 10),
	}
	blockTracker := &mockBlockTracker{
		in:  make(chan op, 10),
		out: make(chan bool, 10),
	}

	sndr, err := sender.NewTxSender(
		st,
		bidder,
		testPricer,
		blockTracker,
		&mockTransferer{},
		big.NewInt(1), // Settlement chain ID
		util.NewTestLogger(os.Stdout),
	)
	if err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := sndr.Start(ctx)

	tx1 := &sender.Transaction{
		Transaction: types.NewTransaction(
			1,
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
	bidder.optinEstimate <- 1

	// Simulate a price estimate
	testPricer.out <- &pricer.BlockPrices{
		CurrentBlockNumber: 0,
		Prices: []pricer.BlockPrice{
			{
				BlockNumber: 1,
				EstimatedPrices: []pricer.EstimatedPrice{
					{
						Confidence:            90,
						PriorityFeePerGasGwei: 1.0,
					},
					{
						Confidence:            95,
						PriorityFeePerGasGwei: 1.5,
					},
					{
						Confidence:            99,
						PriorityFeePerGasGwei: 2.0,
					},
				},
			},
		},
	}

	// Simulate a bid response
	bidOp := <-bidder.in
	if bidOp.rawTx != tx1.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
	}
	resC := make(chan optinbidder.BidStatus, 3)
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusNoOfProviders,
		Arg:  1,
	}
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusAttempted,
		Arg:  uint64(1),
	}
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     1,
			ProviderAddress: "provider1",
		},
	}
	close(resC)
	bidder.out <- resC

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
	if len(res.commitments) != 1 {
		t.Fatalf("expected 1 commitment, got %d", len(res.commitments))
	}

	tx2 := &sender.Transaction{
		Transaction: types.NewTransaction(
			2,
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
	bidder.optinEstimate <- 20

	// Simulate a price estimate
	testPricer.out <- &pricer.BlockPrices{
		CurrentBlockNumber: 1,
		Prices: []pricer.BlockPrice{
			{
				BlockNumber: 2,
				EstimatedPrices: []pricer.EstimatedPrice{
					{
						Confidence:            90,
						PriorityFeePerGasGwei: 1.0,
					},
					{
						Confidence:            95,
						PriorityFeePerGasGwei: 1.5,
					},
					{
						Confidence:            99,
						PriorityFeePerGasGwei: 2.0,
					},
				},
			},
		},
	}

	// Simulate a bid response
	bidOp = <-bidder.in
	if bidOp.rawTx != tx2.Raw[2:] {
		t.Fatalf("expected raw transaction %s, got %s", tx1.Raw, bidOp.rawTx)
	}
	resC = make(chan optinbidder.BidStatus, 3)
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusNoOfProviders,
		Arg:  1,
	}
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusAttempted,
		Arg:  uint64(2),
	}
	resC <- optinbidder.BidStatus{
		Type: optinbidder.BidStatusCommitment,
		Arg: &bidderapiv1.Commitment{
			TxHashes:        []string{tx1.Hash().Hex()},
			BidAmount:       big.NewInt(100).String(),
			BlockNumber:     2,
			ProviderAddress: "provider1",
		},
	}
	close(resC)
	bidder.out <- resC

	checkOp := <-blockTracker.in
	if checkOp.hash != tx2.Hash() {
		t.Fatalf("expected transaction hash %s, got %s", tx2.Hash().Hex(), checkOp.hash.Hex())
	}
	if checkOp.block != 2 {
		t.Fatalf("expected block number 2, got %d", checkOp.block)
	}
	// Simulate transaction inclusion
	blockTracker.out <- true

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
}
