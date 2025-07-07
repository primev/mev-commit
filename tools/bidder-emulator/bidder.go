package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	pb "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type bidder struct {
	client pb.BidderClient
}

type result struct {
	txn      *types.Transaction
	bid      *pb.Bid
	preconfs []*pb.Commitment
}

func newBidder(rpcURL string, depositAmount string) (*bidder, error) {
	creds := insecure.NewCredentials()
	conn, err := grpc.NewClient(rpcURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	client := pb.NewBidderClient(conn)
	b := &bidder{
		client: client,
	}

	return b, b.setup(depositAmount)
}

func (b *bidder) setup(depositAmount string) error {
	status, err := b.client.AutoDepositStatus(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		return fmt.Errorf("failed to get auto deposit status: %w", err)
	}

	if !status.IsAutodepositEnabled {
		_, err := b.client.AutoDeposit(context.Background(), &pb.DepositRequest{
			Amount: depositAmount,
		})
		if err != nil {
			return fmt.Errorf("failed to auto deposit: %w", err)
		}
	}
	return nil
}

func (b *bidder) SendBid(ctx context.Context, txn *types.Transaction, blockNumber int64) (*result, error) {
	txBytes, err := txn.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	// Choose a random bid amount between 1 and 2 gwei
	n := rand.Intn(1_000_000_000)
	bidAmount := new(big.Int).Add(big.NewInt(int64(n)), big.NewInt(1_000_000_000))

	req := &pb.Bid{
		RawTransactions:     []string{hex.EncodeToString(txBytes)},
		BlockNumber:         blockNumber,
		Amount:              bidAmount.String(),
		DecayStartTimestamp: time.Now().Add(100 * time.Millisecond).UnixMilli(),
		DecayEndTimestamp:   time.Now().Add(12 * time.Second).UnixMilli(),
	}

	resp, err := b.client.SendBid(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send bid: %w", err)
	}

	res := &result{
		txn: txn,
		bid: req,
	}

	for {
		preconf, err := resp.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to receive preconf: %w", err)
		}
		res.preconfs = append(res.preconfs, preconf)
	}

	return res, nil
}
