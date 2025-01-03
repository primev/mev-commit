package preconftracker_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	preconftracker "github.com/primev/mev-commit/p2p/pkg/preconfirmation/tracker"
	inmemstorage "github.com/primev/mev-commit/p2p/pkg/storage/inmem"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/primev/mev-commit/x/util"
)

func TestTracker(t *testing.T) {
	t.Parallel()

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(os.Stdout),
		&btABI,
		&pcABI,
		&brABI,
	)

	st := store.New(inmemstorage.New())

	contract := &testPreconfContract{
		openedCommitments: make(chan openedCommitment, 10),
	}

	tracker := preconftracker.NewTracker(
		p2p.PeerTypeBidder,
		common.HexToAddress("0x1234"),
		evtMgr,
		st,
		contract,
		&testReceiptGetter{count: 1},
		func(context.Context) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{
				From: common.HexToAddress("0x1234"),
			}, nil
		},
		util.NewTestLogger(os.Stdout),
	)

	ctx, cancel := context.WithCancel(context.Background())
	doneChan := tracker.Start(ctx)

	winnerProvider := common.HexToAddress("0x1234")
	loserProvider := common.HexToAddress("0x5678")

	getProvider := func(blkNum int64) common.Address {
		if blkNum%2 != 0 {
			return winnerProvider
		}
		return loserProvider
	}

	getBlockNum := func(idx int) int64 {
		return int64(idx/2 + idx%2)
	}

	commitments := make([]*store.EncryptedPreConfirmationWithDecrypted, 0)

	for i := 1; i <= 10; i++ {
		digest := common.HexToHash(fmt.Sprintf("0x%x", i))

		commitments = append(commitments, &store.EncryptedPreConfirmationWithDecrypted{
			EncryptedPreConfirmation: &preconfpb.EncryptedPreConfirmation{
				Commitment: digest.Bytes(),
				Signature:  []byte(fmt.Sprintf("signature%d", i)),
			},
			PreConfirmation: &preconfpb.PreConfirmation{
				Bid: &preconfpb.Bid{
					TxHash:              common.HexToHash(fmt.Sprintf("0x%x", i)).String(),
					BidAmount:           "1000",
					BlockNumber:         getBlockNum(i),
					DecayStartTimestamp: 1,
					DecayEndTimestamp:   2,
					Digest:              []byte(fmt.Sprintf("digest%d", i)),
					Signature:           []byte(fmt.Sprintf("signature%d", i)),
					NikePublicKey:       []byte(fmt.Sprintf("nikePublicKey%d", i)),
				},
				Digest:          digest.Bytes(),
				Signature:       []byte(fmt.Sprintf("signature%d", i)),
				ProviderAddress: getProvider(getBlockNum(i)).Bytes(),
				SharedSecret:    []byte(fmt.Sprintf("sharedSecret%d", i)),
			},
			TxnHash: common.HexToHash(fmt.Sprintf("0x%x", i)),
		})
	}

	for i, c := range commitments {
		err := tracker.TrackCommitment(context.Background(), c)
		if err != nil {
			t.Fatal(err)
		}

		if i == 3 {
			// skip this to simulate transaction error
			continue
		}

		err = publishUnopenedCommitment(evtMgr, &pcABI, preconf.PreconfmanagerUnopenedCommitmentStored{
			Committer:           common.BytesToAddress(c.PreConfirmation.ProviderAddress),
			CommitmentIndex:     common.HexToHash(fmt.Sprintf("0x%x", i+1)),
			CommitmentDigest:    common.BytesToHash(c.EncryptedPreConfirmation.Commitment),
			CommitmentSignature: c.EncryptedPreConfirmation.Signature,
			DispatchTimestamp:   uint64(1),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	amount, ok := new(big.Int).SetString(commitments[4].PreConfirmation.Bid.BidAmount, 10)
	if !ok {
		t.Fatalf("failed to parse bid amount %s", commitments[4].PreConfirmation.Bid.BidAmount)
	}
	// this commitment should not be opened again
	err = publishOpenedCommitment(evtMgr, &pcABI, preconf.PreconfmanagerOpenedCommitmentStored{
		CommitmentIndex:     common.HexToHash(fmt.Sprintf("0x%x", 5)),
		Bidder:              common.HexToAddress("0x1234"),
		Committer:           common.BytesToAddress(commitments[4].PreConfirmation.ProviderAddress),
		BidAmt:              amount,
		BlockNumber:         uint64(commitments[4].PreConfirmation.Bid.BlockNumber),
		BidHash:             common.BytesToHash(commitments[4].PreConfirmation.Bid.Digest),
		DecayStartTimeStamp: uint64(commitments[4].PreConfirmation.Bid.DecayStartTimestamp),
		DecayEndTimeStamp:   uint64(commitments[4].PreConfirmation.Bid.DecayEndTimestamp),
		TxnHash:             commitments[4].PreConfirmation.Bid.TxHash,
		RevertingTxHashes:   commitments[4].PreConfirmation.Bid.RevertingTxHashes,
		CommitmentDigest:    common.BytesToHash(commitments[4].PreConfirmation.Digest),
		BidSignature:        commitments[4].PreConfirmation.Bid.Signature,
		CommitmentSignature: commitments[4].PreConfirmation.Signature,
		DispatchTimestamp:   uint64(1),
		SharedSecretKey:     commitments[4].PreConfirmation.SharedSecret,
	})
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		publishNewWinner(evtMgr, &btABI, blocktracker.BlocktrackerNewL1Block{
			BlockNumber: big.NewInt(int64(i)),
			Winner:      winnerProvider,
			Window:      big.NewInt(1),
		})
	}

	opened := []*store.EncryptedPreConfirmationWithDecrypted{
		commitments[0],
		commitments[1],
		commitments[5],
	}

	for _, c := range opened {
		oc := <-contract.openedCommitments
		if !bytes.Equal(c.EncryptedPreConfirmation.Commitment, oc.encryptedCommitmentIndex[:]) {
			t.Fatalf(
				"expected commitment index %x, got %x",
				c.EncryptedPreConfirmation.CommitmentIndex,
				oc.encryptedCommitmentIndex,
			)
		}
		if c.PreConfirmation.Bid.BidAmount != oc.bid.String() {
			t.Fatalf("expected bid %s, got %d", c.PreConfirmation.Bid.BidAmount, oc.bid)
		}
		if c.PreConfirmation.Bid.BlockNumber != int64(oc.blockNumber) {
			t.Fatalf("expected block number %d, got %d", c.PreConfirmation.Bid.BlockNumber, oc.blockNumber)
		}
		if c.PreConfirmation.Bid.TxHash != oc.txnHash {
			t.Fatalf("expected txn hash %s, got %s", c.PreConfirmation.Bid.TxHash, oc.txnHash)
		}
		if c.PreConfirmation.Bid.DecayStartTimestamp != int64(oc.decayStartTimeStamp) {
			t.Fatalf(
				"expected decay start timestamp %d, got %d",
				c.PreConfirmation.Bid.DecayStartTimestamp,
				oc.decayStartTimeStamp,
			)
		}
		if c.PreConfirmation.Bid.DecayEndTimestamp != int64(oc.decayEndTimeStamp) {
			t.Fatalf("expected decay end timestamp %d, got %d", c.PreConfirmation.Bid.DecayEndTimestamp, oc.decayEndTimeStamp)
		}
		if !bytes.Equal(c.PreConfirmation.Bid.Signature, oc.bidSignature) {
			t.Fatalf(
				"expected bid signature %x, got %x",
				c.PreConfirmation.Bid.Signature,
				oc.bidSignature,
			)
		}
		if !bytes.Equal(c.PreConfirmation.SharedSecret, oc.sharedSecretKey) {
			t.Fatalf(
				"expected shared secret key %x, got %x",
				c.PreConfirmation.SharedSecret,
				oc.sharedSecretKey,
			)
		}
	}

	select {
	case <-contract.openedCommitments:
		t.Fatal("unexpected opened commitment")
	default:
	}

	publishNewWinner(evtMgr, &btABI, blocktracker.BlocktrackerNewL1Block{
		BlockNumber: big.NewInt(6),
		Winner:      winnerProvider,
		Window:      big.NewInt(1),
	})
	publishNewWinner(evtMgr, &btABI, blocktracker.BlocktrackerNewL1Block{
		BlockNumber: big.NewInt(7),
		Winner:      winnerProvider,
		Window:      big.NewInt(1),
	})

	opened = []*store.EncryptedPreConfirmationWithDecrypted{
		commitments[8],
		commitments[9],
	}

	for _, c := range opened {
		oc := <-contract.openedCommitments
		if !bytes.Equal(c.EncryptedPreConfirmation.Commitment, oc.encryptedCommitmentIndex[:]) {
			t.Fatalf(
				"expected commitment index %x, got %x",
				c.EncryptedPreConfirmation.CommitmentIndex,
				oc.encryptedCommitmentIndex,
			)
		}
		if c.PreConfirmation.Bid.BidAmount != oc.bid.String() {
			t.Fatalf("expected bid %s, got %d", c.PreConfirmation.Bid.BidAmount, oc.bid)
		}
		if c.PreConfirmation.Bid.BlockNumber != int64(oc.blockNumber) {
			t.Fatalf("expected block number %d, got %d", c.PreConfirmation.Bid.BlockNumber, oc.blockNumber)
		}
		if c.PreConfirmation.Bid.TxHash != oc.txnHash {
			t.Fatalf("expected txn hash %s, got %s", c.PreConfirmation.Bid.TxHash, oc.txnHash)
		}
		if c.PreConfirmation.Bid.DecayStartTimestamp != int64(oc.decayStartTimeStamp) {
			t.Fatalf(
				"expected decay start timestamp %d, got %d",
				c.PreConfirmation.Bid.DecayStartTimestamp,
				oc.decayStartTimeStamp,
			)
		}
		if c.PreConfirmation.Bid.RevertingTxHashes != oc.revertingTxHashes {
			t.Fatalf("expected reverting tx hashes %s, got %s", c.PreConfirmation.Bid.RevertingTxHashes, oc.revertingTxHashes)
		}
		if c.PreConfirmation.Bid.DecayEndTimestamp != int64(oc.decayEndTimeStamp) {
			t.Fatalf("expected decay end timestamp %d, got %d", c.PreConfirmation.Bid.DecayEndTimestamp, oc.decayEndTimeStamp)
		}
		if !bytes.Equal(c.PreConfirmation.Bid.Signature, oc.bidSignature) {
			t.Fatalf(
				"expected bid signature %x, got %x",
				c.PreConfirmation.Bid.Signature,
				oc.bidSignature,
			)
		}
		if !bytes.Equal(c.PreConfirmation.SharedSecret, oc.sharedSecretKey) {
			t.Fatalf(
				"expected shared secret key %x, got %x",
				c.PreConfirmation.SharedSecret,
				oc.sharedSecretKey,
			)
		}
	}

	cancel()

	<-doneChan
}

func TestTrackerIgnoreOldBlocks(t *testing.T) {
	t.Parallel()

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		t.Fatal(err)
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	brABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	evtMgr := events.NewListener(
		util.NewTestLogger(os.Stdout),
		&btABI,
		&pcABI,
		&brABI,
	)

	st := store.New(inmemstorage.New())

	for _, b := range []int64{1, 12, 13} {
		if err := st.AddWinner(&store.BlockWinner{
			BlockNumber: b,
			Winner:      common.HexToAddress("0x1234"),
		}); err != nil {
			t.Fatal(err)
		}
	}

	contract := &testPreconfContract{
		openedCommitments: make(chan openedCommitment, 10),
	}

	tracker := preconftracker.NewTracker(
		p2p.PeerTypeProvider,
		common.HexToAddress("0x1234"),
		evtMgr,
		st,
		contract,
		&testReceiptGetter{count: 1},
		func(context.Context) (*bind.TransactOpts, error) {
			return &bind.TransactOpts{
				From: common.HexToAddress("0x1234"),
			}, nil
		},
		util.NewTestLogger(os.Stdout),
	)

	ctx, cancel := context.WithCancel(context.Background())
	doneChan := tracker.Start(ctx)

	startTime := time.Now()
	for {
		winners, err := st.BlockWinners()
		if err != nil {
			t.Fatal(err)
		}

		if len(winners) == 0 {
			break
		}

		time.Sleep(100 * time.Millisecond)
		if time.Since(startTime) > 5*time.Second {
			t.Fatal("timed out waiting for block winners to be cleared")
		}
	}

	cancel()

	<-doneChan
}

type openedCommitment struct {
	encryptedCommitmentIndex [32]byte
	bid                      *big.Int
	blockNumber              uint64
	txnHash                  string
	revertingTxHashes        string
	decayStartTimeStamp      uint64
	decayEndTimeStamp        uint64
	bidSignature             []byte
	sharedSecretKey          []byte
}

type testPreconfContract struct {
	openedCommitments chan openedCommitment
}

func (t *testPreconfContract) OpenCommitment(
	_ *bind.TransactOpts,
	encryptedCommitmentIndex [32]byte,
	bid *big.Int,
	blockNumber uint64,
	txnHash string,
	revertingTxHashes string,
	decayStartTimeStamp uint64,
	decayEndTimeStamp uint64,
	bidSignature []byte,
	sharedSecretKey []byte,
) (*types.Transaction, error) {
	t.openedCommitments <- openedCommitment{
		encryptedCommitmentIndex: encryptedCommitmentIndex,
		bid:                      bid,
		blockNumber:              blockNumber,
		txnHash:                  txnHash,
		revertingTxHashes:        revertingTxHashes,
		decayStartTimeStamp:      decayStartTimeStamp,
		decayEndTimeStamp:        decayEndTimeStamp,
		bidSignature:             bidSignature,
		sharedSecretKey:          sharedSecretKey,
	}
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}

type testReceiptGetter struct {
	count int
}

func (t *testReceiptGetter) BatchReceipts(_ context.Context, txns []common.Hash) ([]txmonitor.Result, error) {
	if t.count != len(txns) {
		return nil, fmt.Errorf("expected %d txns, got %d", t.count, len(txns))
	}
	results := make([]txmonitor.Result, 0, len(txns))
	for range txns {
		results = append(results, txmonitor.Result{
			Err: errors.New("test error"),
		})
	}
	return results, nil
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
		c.BidHash,
		c.DecayStartTimeStamp,
		c.DecayEndTimeStamp,
		c.TxnHash,
		c.RevertingTxHashes,
		c.CommitmentDigest,
		c.BidSignature,
		c.CommitmentSignature,
		c.DispatchTimestamp,
		c.SharedSecretKey,
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

func publishNewWinner(
	evtMgr events.EventManager,
	btABI *abi.ABI,
	w blocktracker.BlocktrackerNewL1Block,
) {
	event := btABI.Events["NewL1Block"]

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                        // The first topic is the hash of the event signature
			common.BigToHash(w.BlockNumber), // The next topics are the indexed event parameters
			common.HexToHash(w.Winner.Hex()),
			common.BigToHash(w.Window),
		},
		// Non-indexed parameters are stored in the Data field
		Data: nil,
	}

	evtMgr.PublishLogEvent(context.Background(), testLog)
}
