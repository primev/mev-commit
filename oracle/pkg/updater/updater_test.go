package updater_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	"github.com/primev/mev-commit/oracle/pkg/updater"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/util"
	"golang.org/x/crypto/sha3"
)

func getIdxBytes(idx int64) [32]byte {
	var idxBytes [32]byte
	big.NewInt(idx).FillBytes(idxBytes[:])
	return idxBytes
}

type testBatcher struct {
	failedReceipts map[common.Hash]bool
}

func (t *testBatcher) BatchReceipts(ctx context.Context, txns []common.Hash) ([]txmonitor.Result, error) {
	var results []txmonitor.Result
	for _, txn := range txns {
		status := types.ReceiptStatusSuccessful
		if t.failedReceipts[txn] {
			status = types.ReceiptStatusFailed
		}
		results = append(results, txmonitor.Result{
			Receipt: &types.Receipt{
				TxHash: txn,
				Status: status,
			},
			Err: nil,
		})
	}
	return results, nil
}

type testHasher struct {
	hasher hash.Hash
}

// NewHasher returns a new testHasher instance.
func NewHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

// Reset resets the hash state.
func (h *testHasher) Reset() {
	h.hasher.Reset()
}

// Update updates the hash state with the given key and value.
func (h *testHasher) Update(key, val []byte) error {
	h.hasher.Write(key)
	h.hasher.Write(val)
	return nil
}

// Hash returns the hash value.
func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}

func TestUpdater(t *testing.T) {
	t.Parallel()

	// timestamp of the First block commitment is X
	startTimestamp := time.UnixMilli(1615195200000)
	midTimestamp := startTimestamp.Add(time.Duration(2.5 * float64(time.Second)))
	endTimestamp := startTimestamp.Add(5 * time.Second)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	builderAddr := common.HexToAddress("0xabcd")
	otherBuilderAddr := common.HexToAddress("0xabdc")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := range 10 {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	unopenedCommitments := make([]preconf.PreconfmanagerUnopenedCommitmentStored, 0)
	commitments := make([]preconf.PreconfmanagerOpenedCommitmentStored, 0)

	for i, txn := range txns {
		idxBytes := getIdxBytes(int64(i))

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			TxnHash:             strings.TrimPrefix(txn.Hash().Hex(), "0x"),
			BidAmt:              big.NewInt(10),
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   "",
			// SharedSecretKey:     []byte("shared_secret_key"),
		}

		if i%2 == 0 {
			unopenedCommitment.Committer = builderAddr
			commitment.Committer = builderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		} else {
			unopenedCommitment.Committer = otherBuilderAddr
			commitment.Committer = otherBuilderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		}
	}

	// constructing bundles
	for i := range 10 {
		idxBytes := getIdxBytes(int64(i + 10))

		bundle := strings.TrimPrefix(txns[i].Hash().Hex(), "0x")
		for j := i + 1; j < 10; j++ {
			bundle += "," + strings.TrimPrefix(txns[j].Hash().Hex(), "0x")
		}

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			BidAmt:              big.NewInt(10),
			TxnHash:             bundle,
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   "",
			// SharedSecretKey:     []byte("shared_secret_key"),
		}
		unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
		commitments = append(commitments, commitment)
	}

	register := &testWinnerRegister{
		winners: []testWinner{
			{
				blockNum: 5,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 1,
				},
			},
		},
		settlements:    make(chan testSettlement, 1),
		unopenedCommit: make(chan testEncryptedCommitment, 1),
	}

	body := &types.Body{Transactions: txns, Uncles: nil}

	l1Client := &testEVMClient{
		blocks: map[int64]*types.Block{
			5: types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
		},
		receipts: make(map[string]*types.Receipt),
	}
	for _, txn := range txns {
		receipt := &types.Receipt{
			Status: types.ReceiptStatusSuccessful,
			TxHash: txn.Hash(),
		}
		l1Client.receipts[txn.Hash().Hex()] = receipt
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
		&pcABI,
	)

	oracle := &testOracle{
		commitments: make(chan processedCommitment, 1),
	}

	updtr, err := updater.NewUpdater(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		l1Client,
		register,
		evtMgr,
		oracle,
		&testBatcher{},
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(1),
	}
	publishNewWindow(evtMgr, &btABI, w)

	for _, ec := range unopenedCommitments {
		if err := publishUnopenedCommitment(evtMgr, &pcABI, ec); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case enc := <-register.unopenedCommit:
			if !bytes.Equal(enc.commitmentIdx, ec.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if !bytes.Equal(enc.committer, ec.Committer.Bytes()) {
				t.Fatal("wrong committer")
			}
			if !bytes.Equal(enc.commitmentHash, ec.CommitmentDigest[:]) {
				t.Fatal("wrong commitment hash")
			}
			if !bytes.Equal(enc.commitmentSignature, ec.CommitmentSignature) {
				t.Fatal("wrong commitment signature")
			}
			if enc.dispatchTimestamp != ec.DispatchTimestamp {
				t.Fatal("wrong dispatch timestamp")
			}
		}
	}

	for _, c := range commitments {
		if err := publishOpenedCommitment(evtMgr, &pcABI, c); err != nil {
			t.Fatal(err)
		}

		if c.Committer.Cmp(otherBuilderAddr) == 0 {
			continue
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case commitment := <-oracle.commitments:
			if !bytes.Equal(commitment.commitmentIdx[:], c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if commitment.blockNum.Cmp(big.NewInt(5)) != 0 {
				t.Fatal("wrong block number")
			}
			if commitment.builder != c.Committer {
				t.Fatal("wrong builder")
			}
			if commitment.isSlash {
				t.Fatal("wrong isSlash")
			}
			if commitment.residualDecay.Cmp(big.NewInt(50*updater.PRECISION)) != 0 {
				t.Fatal("wrong residual decay")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case settlement := <-register.settlements:
			if !bytes.Equal(settlement.commitmentIdx, c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if settlement.txHash != c.TxnHash {
				t.Fatal("wrong txn hash")
			}
			if settlement.blockNum != 5 {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(settlement.builder, c.Committer.Bytes()) {
				t.Fatal("wrong builder")
			}
			if settlement.amount.Uint64() != 10 {
				t.Fatal("wrong amount")
			}
			if settlement.settlementType != updater.SettlementTypeReward {
				t.Fatal("wrong settlement type")
			}
			if settlement.decayPercentage != 50*updater.PRECISION {
				t.Fatal("wrong decay percentage")
			}
			if settlement.window != 1 {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}
func TestUpdaterRevertedTxns(t *testing.T) {
	t.Parallel()

	// timestamp of the First block commitment is X
	startTimestamp := time.UnixMilli(1615195200000)
	midTimestamp := startTimestamp.Add(time.Duration(2.5 * float64(time.Second)))
	endTimestamp := startTimestamp.Add(5 * time.Second)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	builderAddr := common.HexToAddress("0xabcd")
	otherBuilderAddr := common.HexToAddress("0xabdc")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := range 10 {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	unopenedCommitments := make([]preconf.PreconfmanagerUnopenedCommitmentStored, 0)
	commitments := make([]preconf.PreconfmanagerOpenedCommitmentStored, 0)

	for i, txn := range txns {
		idxBytes := getIdxBytes(int64(i))

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			TxnHash:             strings.TrimPrefix(txn.Hash().Hex(), "0x"),
			BidAmt:              big.NewInt(10),
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   "",
			// SharedSecretKey:     []byte("shared_secret_key"),
		}

		if i%2 == 0 {
			unopenedCommitment.Committer = builderAddr
			commitment.Committer = builderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		} else {
			unopenedCommitment.Committer = otherBuilderAddr
			commitment.Committer = otherBuilderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		}
	}

	// constructing bundles
	for i := range 10 {
		idxBytes := getIdxBytes(int64(i + 10))

		bundle := strings.TrimPrefix(txns[i].Hash().Hex(), "0x")
		for j := i + 1; j < 10; j++ {
			bundle += "," + strings.TrimPrefix(txns[j].Hash().Hex(), "0x")
		}

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			BidAmt:              big.NewInt(10),
			TxnHash:             bundle,
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   "",
			// SharedSecretKey:     []byte("shared_secret_key"),
		}
		unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
		commitments = append(commitments, commitment)
	}

	register := &testWinnerRegister{
		winners: []testWinner{
			{
				blockNum: 5,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 1,
				},
			},
		},
		settlements:    make(chan testSettlement, 1),
		unopenedCommit: make(chan testEncryptedCommitment, 1),
	}

	body := &types.Body{Transactions: txns, Uncles: nil}

	l1Client := &testEVMClient{
		blocks: map[int64]*types.Block{
			5: types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
		},
		receipts: make(map[string]*types.Receipt),
	}
	for _, txn := range txns {
		receipt := &types.Receipt{
			Status: types.ReceiptStatusFailed,
			TxHash: txn.Hash(),
		}
		l1Client.receipts[txn.Hash().Hex()] = receipt
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
		&pcABI,
	)

	oracle := &testOracle{
		commitments: make(chan processedCommitment, 1),
	}
	testBatcher := &testBatcher{
		failedReceipts: make(map[common.Hash]bool),
	}
	for _, txn := range txns {
		testBatcher.failedReceipts[txn.Hash()] = true
	}

	updtr, err := updater.NewUpdater(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		l1Client,
		register,
		evtMgr,
		oracle,
		testBatcher,
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(1),
	}
	publishNewWindow(evtMgr, &btABI, w)

	for _, ec := range unopenedCommitments {
		if err := publishUnopenedCommitment(evtMgr, &pcABI, ec); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case enc := <-register.unopenedCommit:
			if !bytes.Equal(enc.commitmentIdx, ec.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if !bytes.Equal(enc.committer, ec.Committer.Bytes()) {
				t.Fatal("wrong committer")
			}
			if !bytes.Equal(enc.commitmentHash, ec.CommitmentDigest[:]) {
				t.Fatal("wrong commitment hash")
			}
			if !bytes.Equal(enc.commitmentSignature, ec.CommitmentSignature) {
				t.Fatal("wrong commitment signature")
			}
			if enc.dispatchTimestamp != ec.DispatchTimestamp {
				t.Fatal("wrong dispatch timestamp")
			}
		}
	}

	for _, c := range commitments {
		if err := publishOpenedCommitment(evtMgr, &pcABI, c); err != nil {
			t.Fatal(err)
		}

		if c.Committer.Cmp(otherBuilderAddr) == 0 {
			continue
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case commitment := <-oracle.commitments:
			if !bytes.Equal(commitment.commitmentIdx[:], c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if commitment.blockNum.Cmp(big.NewInt(5)) != 0 {
				t.Fatal("wrong block number")
			}
			if commitment.builder != c.Committer {
				t.Fatal("wrong builder")
			}
			if !commitment.isSlash {
				t.Fatal("wrong isSlash")
			}
			if commitment.residualDecay.Cmp(big.NewInt(50*updater.PRECISION)) != 0 {
				t.Fatal("wrong residual decay")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case settlement := <-register.settlements:
			if !bytes.Equal(settlement.commitmentIdx, c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if settlement.txHash != c.TxnHash {
				t.Fatal("wrong txn hash")
			}
			if settlement.blockNum != 5 {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(settlement.builder, c.Committer.Bytes()) {
				t.Fatal("wrong builder")
			}
			if settlement.amount.Uint64() != 10 {
				t.Fatal("wrong amount")
			}
			if settlement.settlementType != updater.SettlementTypeSlash {
				t.Fatal("wrong settlement type")
			}
			if settlement.decayPercentage != 50*updater.PRECISION {
				t.Fatal("wrong decay percentage")
			}
			if settlement.window != 1 {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestUpdaterRevertedTxnsWithRevertingHashes(t *testing.T) {
	t.Parallel()

	// timestamp of the First block commitment is X
	startTimestamp := time.UnixMilli(1615195200000)
	midTimestamp := startTimestamp.Add(time.Duration(2.5 * float64(time.Second)))
	endTimestamp := startTimestamp.Add(5 * time.Second)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	builderAddr := common.HexToAddress("0xabcd")
	otherBuilderAddr := common.HexToAddress("0xabdc")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := range 10 {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	unopenedCommitments := make([]preconf.PreconfmanagerUnopenedCommitmentStored, 0)
	commitments := make([]preconf.PreconfmanagerOpenedCommitmentStored, 0)

	for i, txn := range txns {
		idxBytes := getIdxBytes(int64(i))

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			TxnHash:             strings.TrimPrefix(txn.Hash().Hex(), "0x"),
			BidAmt:              big.NewInt(10),
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   strings.TrimPrefix(txn.Hash().Hex(), "0x"),
			// SharedSecretKey:     []byte("shared_secret_key"),
		}

		if i%2 == 0 {
			unopenedCommitment.Committer = builderAddr
			commitment.Committer = builderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		} else {
			unopenedCommitment.Committer = otherBuilderAddr
			commitment.Committer = otherBuilderAddr
			unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
			commitments = append(commitments, commitment)
		}
	}

	// constructing bundles
	for i := range 10 {
		idxBytes := getIdxBytes(int64(i + 10))

		bundle := strings.TrimPrefix(txns[i].Hash().Hex(), "0x")
		for j := i + 1; j < 10; j++ {
			bundle += "," + strings.TrimPrefix(txns[j].Hash().Hex(), "0x")
		}

		unopenedCommitment := preconf.PreconfmanagerUnopenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			CommitmentSignature: []byte("signature"),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}
		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			BidAmt:              big.NewInt(10),
			TxnHash:             bundle,
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   bundle,
			// SharedSecretKey:     []byte("shared_secret_key"),
		}
		unopenedCommitments = append(unopenedCommitments, unopenedCommitment)
		commitments = append(commitments, commitment)
	}

	register := &testWinnerRegister{
		winners: []testWinner{
			{
				blockNum: 5,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 1,
				},
			},
		},
		settlements:    make(chan testSettlement, 1),
		unopenedCommit: make(chan testEncryptedCommitment, 1),
	}

	body := &types.Body{Transactions: txns, Uncles: nil}

	l1Client := &testEVMClient{
		blocks: map[int64]*types.Block{
			5: types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
		},
		receipts: make(map[string]*types.Receipt),
	}
	for _, txn := range txns {
		receipt := &types.Receipt{
			Status: types.ReceiptStatusFailed,
			TxHash: txn.Hash(),
		}
		l1Client.receipts[txn.Hash().Hex()] = receipt
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
		&pcABI,
	)

	oracle := &testOracle{
		commitments: make(chan processedCommitment, 1),
	}
	testBatcher := &testBatcher{
		failedReceipts: make(map[common.Hash]bool),
	}
	for _, txn := range txns {
		testBatcher.failedReceipts[txn.Hash()] = true
	}

	updtr, err := updater.NewUpdater(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		l1Client,
		register,
		evtMgr,
		oracle,
		testBatcher,
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(1),
	}
	publishNewWindow(evtMgr, &btABI, w)

	for _, ec := range unopenedCommitments {
		if err := publishUnopenedCommitment(evtMgr, &pcABI, ec); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case enc := <-register.unopenedCommit:
			if !bytes.Equal(enc.commitmentIdx, ec.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if !bytes.Equal(enc.committer, ec.Committer.Bytes()) {
				t.Fatal("wrong committer")
			}
			if !bytes.Equal(enc.commitmentHash, ec.CommitmentDigest[:]) {
				t.Fatal("wrong commitment hash")
			}
			if !bytes.Equal(enc.commitmentSignature, ec.CommitmentSignature) {
				t.Fatal("wrong commitment signature")
			}
			if enc.dispatchTimestamp != ec.DispatchTimestamp {
				t.Fatal("wrong dispatch timestamp")
			}
		}
	}

	for _, c := range commitments {
		if err := publishOpenedCommitment(evtMgr, &pcABI, c); err != nil {
			t.Fatal(err)
		}

		if c.Committer.Cmp(otherBuilderAddr) == 0 {
			continue
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case commitment := <-oracle.commitments:
			if !bytes.Equal(commitment.commitmentIdx[:], c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if commitment.blockNum.Cmp(big.NewInt(5)) != 0 {
				t.Fatal("wrong block number")
			}
			if commitment.builder != c.Committer {
				t.Fatal("wrong builder")
			}
			if commitment.isSlash {
				t.Fatal("wrong isSlash")
			}
			if commitment.residualDecay.Cmp(big.NewInt(50*updater.PRECISION)) != 0 {
				t.Fatal("wrong residual decay")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case settlement := <-register.settlements:
			if !bytes.Equal(settlement.commitmentIdx, c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if settlement.txHash != c.TxnHash {
				t.Fatal("wrong txn hash")
			}
			if settlement.blockNum != 5 {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(settlement.builder, c.Committer.Bytes()) {
				t.Fatal("wrong builder")
			}
			if settlement.amount.Uint64() != 10 {
				t.Fatal("wrong amount")
			}
			if settlement.settlementType != updater.SettlementTypeReward {
				t.Fatal("wrong settlement type")
			}
			if settlement.decayPercentage != 50*updater.PRECISION {
				t.Fatal("wrong decay percentage")
			}
			if settlement.window != 1 {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}
func TestUpdaterBundlesFailure(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	startTimestamp := time.UnixMilli(1615195200000)
	midTimestamp := startTimestamp.Add(time.Duration(2.5 * float64(time.Second)))
	endTimestamp := startTimestamp.Add(5 * time.Second)

	builderAddr := common.HexToAddress("0xabcd")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := range 10 {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	commitments := make([]preconf.PreconfmanagerOpenedCommitmentStored, 0)

	// constructing bundles
	for i := 1; i < 10; i++ {
		idxBytes := getIdxBytes(int64(i))

		bundle := txns[i].Hash().Hex()
		for j := 10 - i; j > 0; j-- {
			bundle += "," + txns[j].Hash().Hex()
		}

		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			BidAmt:              big.NewInt(10),
			TxnHash:             bundle,
			BlockNumber:         5,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
			RevertingTxHashes:   "",
			// SharedSecretKey:     []byte("shared_secret_key"),
		}

		commitments = append(commitments, commitment)
	}

	register := &testWinnerRegister{
		winners: []testWinner{
			{
				blockNum: 5,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 1,
				},
			},
		},
		settlements: make(chan testSettlement, 1),
	}

	body := &types.Body{Transactions: txns, Uncles: nil}

	l1Client := &testEVMClient{
		blocks: map[int64]*types.Block{
			5: types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
		},
		receipts: make(map[string]*types.Receipt),
	}
	for _, txn := range txns {
		receipt := &types.Receipt{
			Status: types.ReceiptStatusSuccessful,
			TxHash: txn.Hash(),
		}
		l1Client.receipts[txn.Hash().Hex()] = receipt
	}

	oracle := &testOracle{
		commitments: make(chan processedCommitment, 1),
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
		&pcABI,
	)

	updtr, err := updater.NewUpdater(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		l1Client,
		register,
		evtMgr,
		oracle,
		&testBatcher{},
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(1),
	}
	publishNewWindow(evtMgr, &btABI, w)

	for _, c := range commitments {
		if err := publishOpenedCommitment(evtMgr, &pcABI, c); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case commitment := <-oracle.commitments:
			if !bytes.Equal(commitment.commitmentIdx[:], c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if commitment.blockNum.Cmp(big.NewInt(5)) != 0 {
				t.Fatal("wrong block number")
			}
			if commitment.builder != c.Committer {
				t.Fatal("wrong builder")
			}
			if !commitment.isSlash {
				t.Fatal("wrong isSlash")
			}
			if commitment.residualDecay.Cmp(big.NewInt(50*updater.PRECISION)) != 0 {
				t.Fatal("wrong residual decay")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case settlement := <-register.settlements:
			if !bytes.Equal(settlement.commitmentIdx, c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if settlement.txHash != c.TxnHash {
				t.Fatal("wrong txn hash")
			}
			if settlement.blockNum != 5 {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(settlement.builder, c.Committer.Bytes()) {
				t.Fatal("wrong builder")
			}
			if settlement.amount.Uint64() != 10 {
				t.Fatal("wrong amount")
			}
			if settlement.settlementType != updater.SettlementTypeSlash {
				t.Fatal("wrong settlement type")
			}
			if settlement.decayPercentage != 50*updater.PRECISION {
				t.Fatal("wrong decay percentage")
			}
			if settlement.window != 1 {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestUpdaterIgnoreCommitments(t *testing.T) {
	t.Parallel()

	// timestamp of the First block commitment is X
	startTimestamp := time.UnixMilli(1615195200000)
	midTimestamp := startTimestamp.Add(time.Duration(2.5 * float64(time.Second)))
	endTimestamp := startTimestamp.Add(5 * time.Second)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	builderAddr := common.HexToAddress("0xabcd")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := range 10 {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	commitments := make([]preconf.PreconfmanagerOpenedCommitmentStored, 0)

	for i, txn := range txns {
		idxBytes := getIdxBytes(int64(i))

		// block no 5 will not be settled, so we will ignore it
		// block no 8 will not be settled as no winner is registered for it
		// block no 10 will be settled
		blockNum := uint64(5)
		if i > 5 && i < 8 {
			blockNum = 8
		} else if i >= 8 {
			blockNum = 10
		}

		commitment := preconf.PreconfmanagerOpenedCommitmentStored{
			CommitmentIndex:     idxBytes,
			Committer:           builderAddr,
			BidAmt:              big.NewInt(10),
			TxnHash:             strings.TrimPrefix(txn.Hash().Hex(), "0x"),
			RevertingTxHashes:   "",
			BlockNumber:         blockNum,
			CommitmentDigest:    common.HexToHash(fmt.Sprintf("0x%02d", i)),
			// CommitmentSignature: []byte("signature"),
			DecayStartTimeStamp: uint64(startTimestamp.UnixMilli()),
			DecayEndTimeStamp:   uint64(endTimestamp.UnixMilli()),
			DispatchTimestamp:   uint64(midTimestamp.UnixMilli()),
		}

		if i == 9 {
			// duplicate commitment
			commitment.CommitmentIndex = getIdxBytes(int64(i - 1))
		}

		commitments = append(commitments, commitment)
	}

	register := &testWinnerRegister{
		winners: []testWinner{
			{
				blockNum: 5,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 1,
				},
			},
			{
				blockNum: 10,
				winner: updater.Winner{
					Winner: builderAddr.Bytes(),
					Window: 5,
				},
			},
		},
		settlements:    make(chan testSettlement, 1),
		unopenedCommit: make(chan testEncryptedCommitment, 1),
	}

	body := &types.Body{Transactions: txns, Uncles: nil}

	l1Client := &testEVMClient{
		blocks: map[int64]*types.Block{
			5:  types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
			8:  types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
			10: types.NewBlock(&types.Header{}, body, []*types.Receipt{}, trie.NewStackTrie(nil)),
		},
		receipts: make(map[string]*types.Receipt),
	}
	for _, txn := range txns {
		receipt := &types.Receipt{
			Status: types.ReceiptStatusSuccessful,
			TxHash: txn.Hash(),
		}
		l1Client.receipts[txn.Hash().Hex()] = receipt
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
		&pcABI,
	)

	oracle := &testOracle{
		commitments: make(chan processedCommitment, 1),
	}

	updtr, err := updater.NewUpdater(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		l1Client,
		register,
		evtMgr,
		oracle,
		&testBatcher{},
	)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(5),
	}
	publishNewWindow(evtMgr, &btABI, w)

	for i, c := range commitments {
		if err := publishOpenedCommitment(evtMgr, &pcABI, c); err != nil {
			t.Fatal(err)
		}

		if i < 8 {
			continue
		}

		if i == 9 {
			// duplicate commitment
			continue
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case commitment := <-oracle.commitments:
			if !bytes.Equal(commitment.commitmentIdx[:], c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if commitment.blockNum.Cmp(big.NewInt(10)) != 0 {
				t.Fatal("wrong block number", commitment.blockNum)
			}
			if commitment.builder != c.Committer {
				t.Fatal("wrong builder")
			}
			if commitment.isSlash {
				t.Fatal("wrong isSlash")
			}
			if commitment.residualDecay.Cmp(big.NewInt(50*updater.PRECISION)) != 0 {
				t.Fatal("wrong residual decay")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout")
		case settlement := <-register.settlements:
			if !bytes.Equal(settlement.commitmentIdx, c.CommitmentIndex[:]) {
				t.Fatal("wrong commitment index")
			}
			if settlement.txHash != c.TxnHash {
				t.Fatal("wrong txn hash")
			}
			if settlement.blockNum != 10 {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(settlement.builder, c.Committer.Bytes()) {
				t.Fatal("wrong builder")
			}
			if settlement.amount.Uint64() != 10 {
				t.Fatal("wrong amount")
			}
			if settlement.settlementType != updater.SettlementTypeReward {
				t.Fatal("wrong settlement type")
			}
			if settlement.decayPercentage != 50*updater.PRECISION {
				t.Fatal("wrong decay percentage")
			}
			if settlement.window != 5 {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

type testSettlement struct {
	commitmentIdx   []byte
	txHash          string
	blockNum        int64
	builder         []byte
	amount          *big.Int
	settlementType  updater.SettlementType
	decayPercentage int64
	window          int64
	chainhash       []byte
	nonce           uint64
}

type testEncryptedCommitment struct {
	commitmentIdx       []byte
	committer           []byte
	commitmentHash      []byte
	commitmentSignature []byte
	dispatchTimestamp   uint64
}

type testWinner struct {
	blockNum int64
	winner   updater.Winner
}

type testWinnerRegister struct {
	mu              sync.Mutex
	winners         []testWinner
	setttlementIdxs [][]byte
	settlements     chan testSettlement
	unopenedCommit  chan testEncryptedCommitment
}

func (t *testWinnerRegister) IsSettled(ctx context.Context, commitmentIdx []byte) (bool, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, idx := range t.setttlementIdxs {
		if bytes.Equal(idx, commitmentIdx) {
			return true, nil
		}
	}
	return false, nil
}

func (t *testWinnerRegister) GetWinner(ctx context.Context, blockNum int64) (updater.Winner, error) {
	for _, w := range t.winners {
		if w.blockNum == blockNum {
			return w.winner, nil
		}
	}
	return updater.Winner{}, sql.ErrNoRows
}

func (t *testWinnerRegister) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount *big.Int,
	builder []byte,
	_ []byte,
	settlementType updater.SettlementType,
	decayPercentage int64,
	window int64,
	chainhash []byte,
	nonce uint64,
) error {
	t.mu.Lock()
	t.setttlementIdxs = append(t.setttlementIdxs, commitmentIdx)
	t.mu.Unlock()

	t.settlements <- testSettlement{
		commitmentIdx:   commitmentIdx,
		txHash:          txHash,
		blockNum:        blockNum,
		amount:          amount,
		builder:         builder,
		settlementType:  settlementType,
		decayPercentage: decayPercentage,
		window:          window,
		chainhash:       chainhash,
		nonce:           nonce,
	}
	return nil
}

func (t *testWinnerRegister) AddEncryptedCommitment(
	ctx context.Context,
	commitmentIdx []byte,
	committer []byte,
	commitmentHash []byte,
	commitmentSignature []byte,
	dispatchTimestamp uint64,
) error {
	t.unopenedCommit <- testEncryptedCommitment{
		commitmentIdx:       commitmentIdx,
		committer:           committer,
		commitmentHash:      commitmentHash,
		commitmentSignature: commitmentSignature,
		dispatchTimestamp:   dispatchTimestamp,
	}
	return nil
}

type testEVMClient struct {
	blocks   map[int64]*types.Block
	receipts map[string]*types.Receipt
}

func (t *testEVMClient) BlockByNumber(ctx context.Context, blkNum *big.Int) (*types.Block, error) {
	blk, found := t.blocks[blkNum.Int64()]
	if !found {
		return nil, fmt.Errorf("block %d not found", blkNum.Int64())
	}
	return blk, nil
}

type processedCommitment struct {
	commitmentIdx [32]byte
	blockNum      *big.Int
	builder       common.Address
	isSlash       bool
	residualDecay *big.Int
}

type testOracle struct {
	commitments chan processedCommitment
}

func (t *testOracle) ProcessBuilderCommitmentForBlockNumber(
	commitmentIdx [32]byte,
	blockNum *big.Int,
	builder common.Address,
	isSlash bool,
	residualDecay *big.Int,
) (*types.Transaction, error) {
	t.commitments <- processedCommitment{
		commitmentIdx: commitmentIdx,
		blockNum:      blockNum,
		builder:       builder,
		isSlash:       isSlash,
		residualDecay: residualDecay,
	}
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}

func publishUnopenedCommitment(
	evtMgr events.EventManager,
	pcABI *abi.ABI,
	ec preconf.PreconfmanagerUnopenedCommitmentStored,
) error {
	event := pcABI.Events["UnopenedCommitmentStored"]
	buf, err := event.Inputs.NonIndexed().Pack(
		ec.Committer,
		ec.CommitmentDigest,
		ec.CommitmentSignature,
		ec.DispatchTimestamp,
	)
	if err != nil {
		return err
	}

	commitmentIndex := common.BytesToHash(ec.CommitmentIndex[:])

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,        // The first topic is the hash of the event signature
			commitmentIndex, // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: buf,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishOpenedCommitment(
	evtMgr events.EventManager,
	pcABI *abi.ABI,
	c preconf.PreconfmanagerOpenedCommitmentStored,
) error {
	event := pcABI.Events["OpenedCommitmentStored"]
	buf, err := event.Inputs.NonIndexed().Pack(
		c.Bidder,
		c.Committer,
		c.BidAmt,
		c.BlockNumber,
		// c.BidHash,
		c.DecayStartTimeStamp,
		c.DecayEndTimeStamp,
		c.TxnHash,
		c.RevertingTxHashes,
		c.CommitmentDigest,
		// c.BidSignature,
		// c.CommitmentSignature,
		c.DispatchTimestamp,
		// c.SharedSecretKey,
	)
	if err != nil {
		return err
	}

	commitmentIndex := common.BytesToHash(c.CommitmentIndex[:])

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,        // The first topic is the hash of the event signature
			commitmentIndex, // The next topics are the indexed event parameters
		},
		// Since there are no non-indexed parameters, Data is empty
		Data: buf,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
	return nil
}

func publishNewWindow(
	evtMgr events.EventManager,
	btABI *abi.ABI,
	w blocktracker.BlocktrackerNewWindow,
) {
	event := btABI.Events["NewWindow"]

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                   // The first topic is the hash of the event signature
			common.BigToHash(w.Window), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: nil,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
}
