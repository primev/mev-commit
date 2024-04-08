package l1Listener_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-commit/oracle/pkg/l1Listener"
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

	l := l1Listener.NewL1Listener(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		ethClient,
		reg,
	)
	ctx, cancel := context.WithCancel(context.Background())

	cl := l1Listener.SetCheckInterval(100 * time.Millisecond)
	t.Cleanup(cl)

	done := l.Start(ctx)

	for i := 1; i < 10; i++ {
		ethClient.AddHeader(uint64(i), &types.Header{
			Number: big.NewInt(int64(i)),
			Extra:  []byte(fmt.Sprintf("b%d", i)),
		})

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for winner", i)
		case winner := <-reg.winners:
			if winner.blockNum != int64(i) {
				t.Fatal("wrong block number")
			}
			if winner.winner != fmt.Sprintf("b%d", i) {
				t.Fatal("wrong winner")
			}
		}
	}

	// no winner
	ethClient.AddHeader(10, &types.Header{
		Number: big.NewInt(10),
	})

	// error registering winner, ensure it is retried
	ethClient.errC <- errors.New("dummy error")
	ethClient.AddHeader(11, &types.Header{
		Number: big.NewInt(11),
		Extra:  []byte("b11"),
	})

	// ensure no winner is sent for the previous block
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for winner")
	case winner := <-reg.winners:
		if winner.blockNum != 11 {
			t.Fatal("wrong block number")
		}
		if winner.winner != "b11" {
			t.Fatal("wrong winner")
		}
	}

	cancel()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for done")
	case <-done:
	}
}

type winnerObj struct {
	blockNum int64
	winner   string
}

type testRegister struct {
	winners chan winnerObj
}

func (t *testRegister) RegisterWinner(_ context.Context, blockNum int64, winner string) error {
	t.winners <- winnerObj{blockNum: blockNum, winner: winner}
	return nil
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
