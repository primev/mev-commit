package main

import (
	"context"
	"crypto/tls"
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
	"google.golang.org/grpc/credentials"
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
	conn, err := grpc.NewClient(
		rpcURL,
		grpc.WithTransportCredentials(credentials.NewTLS(
			&tls.Config{InsecureSkipVerify: true},
		)),
	)
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
	depositAmountInt, ok := new(big.Int).SetString(depositAmount, 10)
	if !ok {
		return fmt.Errorf("failed to parse deposit amount")
	}

	status, err := b.client.DepositManagerStatus(context.Background(), &pb.DepositManagerStatusRequest{})
	if err != nil {
		return fmt.Errorf("failed to get deposit manager status: %w", err)
	}
	if !status.Enabled {
		resp, err := b.client.EnableDepositManager(context.Background(), &pb.EnableDepositManagerRequest{})
		if err != nil {
			return fmt.Errorf("failed to enable deposit manager: %w", err)
		}
		if !resp.Success {
			return fmt.Errorf("failed to enable deposit manager")
		}
	}

	validProviders, err := b.client.GetValidProviders(context.Background(), &pb.GetValidProvidersRequest{})
	if err != nil {
		return fmt.Errorf("failed to get valid providers: %w", err)
	}
	if len(validProviders.ValidProviders) == 0 {
		return fmt.Errorf("no valid providers found")
	}

	targetDeposits := make([]*pb.TargetDeposit, len(validProviders.ValidProviders))
	for i, provider := range validProviders.ValidProviders {
		targetDeposits[i] = &pb.TargetDeposit{
			Provider:      provider,
			TargetDeposit: depositAmountInt.Uint64(),
		}
	}

	resp, err := b.client.SetTargetDeposits(context.Background(), &pb.SetTargetDepositsRequest{
		TargetDeposits: targetDeposits,
	})
	if err != nil {
		return fmt.Errorf("failed to set target deposits: %w", err)
	}
	if len(resp.SuccessfullySetDeposits) != len(targetDeposits) {
		return fmt.Errorf("failed to set target deposits")
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
