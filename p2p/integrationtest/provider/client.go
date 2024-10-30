package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"math/big"
	"time"

	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type ProviderClient struct {
	conn         *grpc.ClientConn
	client       providerapiv1.ProviderClient
	logger       *slog.Logger
	senderC      chan *providerapiv1.BidResponse
	senderClosed chan struct{}
}

func NewProviderClient(
	serverAddr string,
	logger *slog.Logger,
) (*ProviderClient, error) {
	// Since we don't know if the server has TLS enabled on its rpc
	// endpoint, we try different strategies from most secure to
	// least secure. In the future, when only TLS-enabled servers
	// are allowed, only the TLS system pool certificate strategy
	// should be used.
	var (
		conn *grpc.ClientConn
		err  error
	)
	for _, e := range []struct {
		strategy   string
		isSecure   bool
		credential credentials.TransportCredentials
	}{
		{"TLS system pool certificate", true, credentials.NewClientTLSFromCert(nil, "")},
		{"TLS skip verification", false, credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})},
		{"TLS disabled", false, insecure.NewCredentials()},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		logger.Info("dialing to grpc server", "strategy", e.strategy)

		// nolint:staticcheck
		conn, err = grpc.DialContext(
			ctx,
			serverAddr,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(e.credential),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			logger.Warn("failed to dial grpc server", "error", err)
			cancel()
			continue
		}

		cancel()
		if !e.isSecure {
			logger.Warn("established connection with the grpc server has potential security risk")
		}
		break
	}
	if conn == nil {
		return nil, errors.New("dialing of grpc server failed")
	}

	b := &ProviderClient{
		conn:         conn,
		client:       providerapiv1.NewProviderClient(conn),
		logger:       logger,
		senderC:      make(chan *providerapiv1.BidResponse),
		senderClosed: make(chan struct{}),
	}

	if err := b.startSender(); err != nil {
		return nil, errors.Join(err, b.Close())
	}
	return b, nil
}

func (b *ProviderClient) Close() error {
	close(b.senderC)
	return b.conn.Close()
}

func (b *ProviderClient) CheckAndStake() error {
	stakeAmt, err := b.client.GetStake(context.Background(), &providerapiv1.EmptyMessage{})
	if err != nil {
		b.logger.Error("failed to get stake amount", "err", err)
		return err
	}

	b.logger.Info("stake amount", "stake", stakeAmt.Amount)

	stakedAmt, set := big.NewInt(0).SetString(stakeAmt.Amount, 10)
	if !set {
		b.logger.Error("failed to parse stake amount")
		return errors.New("failed to parse stake amount")
	}

	if stakedAmt.Cmp(big.NewInt(0)) > 0 {
		b.logger.Error("bidder already staked")
		return nil
	}

	_, err = b.client.RegisterStake(context.Background(), &providerapiv1.StakeRequest{
		Amount:       "10000000000000000000",
		BlsPublicKey: "abf1ad5ec0512cb1adabe457882fa550b4935f1f7df9658e46af882049ec16da698c323af8c98c3f1f9570ebc4042a83",
	})
	if err != nil {
		b.logger.Error("failed to register stake", "err", err)
		return err
	}

	b.logger.Info("staked 10 ETH")

	return nil

}

func (b *ProviderClient) startSender() error {
	stream, err := b.client.SendProcessedBids(context.Background())
	if err != nil {
		return err
	}

	go func() {
		defer close(b.senderClosed)
		for {
			select {
			case <-stream.Context().Done():
				b.logger.Warn("closing client conn")
				return
			case resp, more := <-b.senderC:
				if !more {
					b.logger.Warn("closed sender chan")
					return
				}
				err := stream.Send(resp)
				if err != nil {
					b.logger.Error("failed sending response", "error", err)
				}
			}
		}
	}()

	return nil
}

// ReceiveBids opens a new RPC connection with the mev-commit node to receive bids.
// Each call to this function opens a new connection and the bids are randomly
// assigned to one of the existing connections from mev-commit node. So if you run
// multiple listeners, they will get unique bids in a non-deterministic fashion.
func (b *ProviderClient) ReceiveBids() (chan *providerapiv1.Bid, error) {
	emptyMessage := &providerapiv1.EmptyMessage{}
	bidStream, err := b.client.ReceiveBids(context.Background(), emptyMessage)
	if err != nil {
		return nil, err
	}

	bidC := make(chan *providerapiv1.Bid)
	go func() {
		defer close(bidC)
		for {
			bid, err := bidStream.Recv()
			if err != nil {
				b.logger.Error("failed receiving bid", "error", err)
				return
			}
			select {
			case <-bidStream.Context().Done():
			case bidC <- bid:
			}
		}
	}()

	return bidC, nil
}

// SendBidResponse can be used to send the status of the bid back to the mev-commit
// node. The provider can use his own logic to decide upon the bid and once he is
// ready to make a decision, this status has to be sent back to mev-commit to decide
// what to do with this bid. The sender is a single global worker which sends back
// the messages on grpc.
func (b *ProviderClient) SendBidResponse(
	ctx context.Context,
	bidResponse *providerapiv1.BidResponse,
) error {

	select {
	case <-b.senderClosed:
		return errors.New("sender closed")
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.senderC <- bidResponse:
		return nil
	}
}
