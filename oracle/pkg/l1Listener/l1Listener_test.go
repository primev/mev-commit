package l1Listener_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	"github.com/primev/mev-commit/oracle/pkg/l1Listener"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/util"
)

func TestL1Listener(t *testing.T) {
	t.Parallel()

	reg := &testRegister{
		winners: make(chan winnerObj),
	}
	ethClient := &testEthClient{
		headers: make(map[uint64]*types.Header),
		errC:    make(chan error, 1),
	}
	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}
	eventManager := events.NewListener(
		util.NewTestLogger(io.Discard),
		&btABI,
	)
	rec := &testRecorder{
		updates: make(chan l1Update),
	}

	testRelayQuerier := &testRelayQuerier{
		responses: map[int64]string{
			1:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			2:  "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			3:  "0x9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210fedcba",
			4:  "0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdef",
			5:  "0x123412341234123412341234123412341234123412341234123412341234123412341234123412341234123412341234",
			6:  "0x567856785678567856785678567856785678567856785678567856785678567856785678567856785678567856785678",
			7:  "0x9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc9abc",
			8:  "0xdef0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0def0",
			9:  "0x345634563456345634563456345634563456345634563456345634563456345634563456345634563456345634563456",
			10: "0x789078907890789078907890789078907890789078907890789078907890789078907890789078907890789078907890",
		},
	}

	l := l1Listener.NewL1Listener(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
		ethClient,
		reg,
		eventManager,
		rec,
		testRelayQuerier,
	)
	ctx, cancel := context.WithCancel(context.Background())

	cl := l1Listener.SetCheckInterval(100 * time.Millisecond)
	t.Cleanup(cl)

	done := l.Start(ctx)

	// no winner
	ethClient.AddHeader(2, &types.Header{
		Number: big.NewInt(2),
	})

	builderPubKeyBytes, _ := hex.DecodeString("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	if err != nil {
		t.Fatal(err)
	}

	// ensure no winner is sent for the previous block
	select {
	case <-time.After(10 * time.Second):
		t.Fatal("timeout waiting for winner")
	case update := <-rec.updates:
		if update.blockNum.Int64() != 1 {
			t.Fatalf("wrong block number: %d", update.blockNum.Int64())
		}
		if !bytes.Equal(update.winner, builderPubKeyBytes) {
			t.Fatalf("wrong winner: %s", update.winner)
		}
	}

	for i := 1; i < 10; i++ {
		addr := common.HexToAddress(fmt.Sprintf("0x%d", i))
		go func() {
			publishLog(
				eventManager,
				big.NewInt(int64(i)),
				addr,
				big.NewInt(int64(i)),
			)
			if err != nil {
				t.Error(err)
			}
		}()

		select {
		case <-time.After(10 * time.Second):
			t.Fatal("timeout waiting for winner", i)
		case winner := <-reg.winners:
			if winner.blockNum != int64(i) {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(winner.winner, addr.Bytes()) {
				t.Fatal("wrong winner")
			}
			if winner.window != int64(i) {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for done")
	case <-done:
	}
}

type testRelayQuerier struct {
	responses map[int64]string
	mu        sync.Mutex
}

func (t *testRelayQuerier) SetResponse(blockNumber int64, builderPubKey string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.responses[blockNumber] = builderPubKey
}

func (t *testRelayQuerier) Query(ctx context.Context, blockNumber int64, blockHash string) (string, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if response, ok := t.responses[blockNumber]; ok {
		return response, nil
	}
	return "", fmt.Errorf("no response set for block number %d", blockNumber)
}

type winnerObj struct {
	blockNum int64
	winner   []byte
	window   int64
}

type testRegister struct {
	winners chan winnerObj
}

func (t *testRegister) RegisterWinner(_ context.Context, blockNum int64, winner []byte, window int64) error {
	t.winners <- winnerObj{blockNum: blockNum, winner: winner, window: window}
	return nil
}

func (t *testRegister) LastWinnerBlock() (int64, error) {
	return 0, nil
}

type testEthClient struct {
	mu      sync.Mutex
	headers map[uint64]*types.Header
	errC    chan error
}

func (t *testEthClient) AddHeader(blockNum uint64, hdr *types.Header) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headers[blockNum] = hdr
}

func (t *testEthClient) BlockNumber(_ context.Context) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.headers) == 0 {
		return 0, nil
	}
	blks := make([]uint64, len(t.headers))
	for k := range t.headers {
		blks = append(blks, k)
	}

	sort.Slice(blks, func(i, j int) bool {
		return blks[i] < blks[j]
	})

	return blks[len(blks)-1], nil
}

func (t *testEthClient) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	select {
	case err := <-t.errC:
		return nil, err
	default:
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	hdr, ok := t.headers[number.Uint64()]
	if !ok {
		return nil, errors.New("header not found")
	}
	return hdr, nil
}

func (t *testEthClient) BlockByNumber(_ context.Context, number *big.Int) (*types.Block, error) {
	return nil, nil
}

func publishLog(
	eventManager events.EventManager,
	blockNum *big.Int,
	winner common.Address,
	window *big.Int,
) {
	eventSignature := []byte("NewL1Block(uint256,address,uint256)")
	hashEventSignature := crypto.Keccak256Hash(eventSignature)

	blockNumber := common.BigToHash(blockNum)
	winnerHash := common.HexToHash(winner.Hex())
	windowNumber := common.BigToHash(window)

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			hashEventSignature, // The first topic is the hash of the event signature
			blockNumber,        // The next topics are the indexed event parameters
			winnerHash,
			windowNumber,
		},
		// Since there are no non-indexed parameters, Data is empty
		Data: []byte{},
	}

	eventManager.PublishLogEvent(context.Background(), testLog)
}

type l1Update struct {
	blockNum *big.Int
	winner   []byte
}

type testRecorder struct {
	updates chan l1Update
}

func (t *testRecorder) RecordL1Block(blockNum *big.Int, winner []byte) (*types.Transaction, error) {
	t.updates <- l1Update{blockNum: blockNum, winner: winner}
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}
