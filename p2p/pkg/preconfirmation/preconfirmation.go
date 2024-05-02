package preconfirmation

import (
	"context"
	"errors"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	preconfcontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/preconf"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	signer "github.com/primevprotocol/mev-commit/p2p/pkg/signer/preconfsigner"
	"github.com/primevprotocol/mev-commit/p2p/pkg/topology"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProtocolName    = "preconfirmation"
	ProtocolVersion = "2.0.0"
)

type Preconfirmation struct {
	signer       signer.Signer
	topo         Topology
	streamer     p2p.Streamer
	us           BidderStore
	processer    BidProcessor
	commitmentDA preconfcontract.Interface
	logger       *slog.Logger
	metrics      *metrics
}

type Topology interface {
	GetPeers(topology.Query) []p2p.Peer
}

type BidderStore interface {
	CheckBidderAllowance(context.Context, common.Address) bool
}

type BidProcessor interface {
	ProcessBid(context.Context, *preconfpb.Bid) (chan providerapiv1.ProcessedBidResponse, error)
}

func New(
	topo Topology,
	streamer p2p.Streamer,
	signer signer.Signer,
	us BidderStore,
	processor BidProcessor,
	commitmentDA preconfcontract.Interface,
	logger *slog.Logger,
) *Preconfirmation {
	return &Preconfirmation{
		topo:         topo,
		streamer:     streamer,
		signer:       signer,
		us:           us,
		processer:    processor,
		commitmentDA: commitmentDA,
		logger:       logger,
		metrics:      newMetrics(),
	}
}

func (p *Preconfirmation) preconfStream() p2p.StreamDesc {
	return p2p.StreamDesc{
		Name:    ProtocolName,
		Version: ProtocolVersion,
		Handler: p.handleBid,
	}
}

func (p *Preconfirmation) Streams() []p2p.StreamDesc {
	return []p2p.StreamDesc{p.preconfStream()}
}

// SendBid is meant to be called by the bidder to construct and send bids to the provider.
// It takes the txHash, the bid amount in wei and the maximum valid block number.
// It waits for preConfirmations from all providers and then returns.
// It returns an error if the bid is not valid.
func (p *Preconfirmation) SendBid(
	ctx context.Context,
	txHash string,
	bidAmt string,
	blockNumber int64,
	decayStartTimestamp int64,
	decayEndTimestamp int64,
) (chan *preconfpb.PreConfirmation, error) {
	signedBid, err := p.signer.ConstructSignedBid(txHash, bidAmt, blockNumber, decayStartTimestamp, decayEndTimestamp)
	if err != nil {
		p.logger.Error("constructing signed bid", "error", err, "txHash", txHash)
		return nil, err
	}
	p.logger.Info("constructed signed bid", "signedBid", signedBid)

	providers := p.topo.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	if len(providers) == 0 {
		p.logger.Error("no providers available", "txHash", txHash)
		return nil, errors.New("no providers available")
	}

	// Create a new channel to receive preConfirmations
	preConfirmations := make(chan *preconfpb.PreConfirmation, len(providers))

	wg := sync.WaitGroup{}
	for idx := range providers {
		wg.Add(1)
		go func(provider p2p.Peer) {
			defer wg.Done()

			logger := p.logger.With("provider", provider, "bid", txHash)

			providerStream, err := p.streamer.NewStream(
				ctx,
				provider,
				nil,
				p.preconfStream(),
			)
			if err != nil {
				logger.Error("creating stream", "error", err)
				return
			}

			logger.Info("sending signed bid", "signedBid", signedBid)

			err = providerStream.WriteMsg(ctx, signedBid)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("writing message", "error", err)
				return
			}
			p.metrics.SentBidsCount.Inc()

			preConfirmation := new(preconfpb.PreConfirmation)
			err = providerStream.ReadMsg(ctx, preConfirmation)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("reading message", "error", err)
				return
			}

			_ = providerStream.Close()

			// Process preConfirmation as a bidder
			providerAddress, err := p.signer.VerifyPreConfirmation(preConfirmation)
			if err != nil {
				logger.Error("verifying provider signature", "error", err)
				return
			}
			preConfirmation.ProviderAddress = make([]byte, len(providerAddress))
			copy(preConfirmation.ProviderAddress, providerAddress[:])
			logger.Info("received preconfirmation", "preConfirmation", preConfirmation)
			p.metrics.ReceivedPreconfsCount.Inc()

			select {
			case preConfirmations <- preConfirmation:
			case <-ctx.Done():
				logger.Error("context cancelled", "error", ctx.Err())
				return
			}
		}(providers[idx])
	}

	go func() {
		wg.Wait()
		close(preConfirmations)
	}()

	return preConfirmations, nil
}

var ErrInvalidBidderTypeForBid = errors.New("invalid bidder type for bid")

// handlebid is the function that is called when a bid is received
// It is meant to be used by the provider exclusively to read the bid value from the bidder.
func (p *Preconfirmation) handleBid(
	ctx context.Context,
	peer p2p.Peer,
	stream p2p.Stream,
) error {
	if peer.Type != p2p.PeerTypeBidder {
		return ErrInvalidBidderTypeForBid
	}

	bid := new(preconfpb.Bid)
	err := stream.ReadMsg(ctx, bid)
	if err != nil {
		return err
	}

	p.logger.Info("received bid", "bid", bid)

	ethAddress, err := p.signer.VerifyBid(bid)
	if err != nil {
		p.logger.Error("verifying bid", "error", err)
		return status.Errorf(codes.InvalidArgument, "invalid bid: %v", err)
	}

	if !p.us.CheckBidderAllowance(ctx, *ethAddress) {
		p.logger.Error("bidder does not have enough allowance", "ethAddress", ethAddress)
		return status.Errorf(codes.FailedPrecondition, "bidder not allowed")
	}

	bidAmt, _ := new(big.Int).SetString(bid.BidAmount, 10)

	// try to enqueue for 5 seconds
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	statusC, err := p.processer.ProcessBid(ctx, bid)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case st := <-statusC:
		switch st.Status {
		case providerapiv1.BidResponse_STATUS_REJECTED:
			return status.Errorf(codes.Internal, "bid rejected")
		case providerapiv1.BidResponse_STATUS_ACCEPTED:
			preConfirmation, err := p.signer.ConstructPreConfirmation(bid)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to construct preconfirmation: %v", err)
			}
			p.logger.Info("sending preconfirmation", "preConfirmation", preConfirmation)
			err = p.commitmentDA.StoreCommitment(
				ctx,
				bidAmt,
				uint64(preConfirmation.Bid.BlockNumber),
				preConfirmation.Bid.TxHash,
				uint64(preConfirmation.Bid.DecayStartTimestamp),
				uint64(preConfirmation.Bid.DecayEndTimestamp),
				preConfirmation.Bid.Signature,
				preConfirmation.Signature,
				uint64(st.DispatchTimestamp),
			)
			if err != nil {
				p.logger.Error("storing commitment", "error", err)
				return status.Errorf(codes.Internal, "failed to store commitment: %v", err)
			}
			return stream.WriteMsg(ctx, preConfirmation)
		}
	}

	return nil
}
