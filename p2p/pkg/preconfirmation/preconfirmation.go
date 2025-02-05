package preconfirmation

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	preconfpb "github.com/primev/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primev/mev-commit/p2p/gen/go/providerapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/preconfirmation/store"
	providerapi "github.com/primev/mev-commit/p2p/pkg/rpc/provider"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProtocolName    = "preconfirmation"
	ProtocolVersion = "3.0.0"
)

type Preconfirmation struct {
	encryptor       Encryptor
	topo            Topology
	streamer        p2p.Streamer
	depositMgr      DepositManager
	processer       BidProcessor
	commitmentDA    PreconfContract
	tracker         Tracker
	optsGetter      OptsGetter
	logger          *slog.Logger
	metrics         *metrics
	providerTimeout time.Duration
}

type OptsGetter func(context.Context) (*bind.TransactOpts, error)

type Topology interface {
	GetPeers(topology.Query) []p2p.Peer
}

type BidProcessor interface {
	ProcessBid(context.Context, *preconfpb.Bid) (chan providerapi.ProcessedBidResponse, error)
}

type DepositManager interface {
	CheckAndDeductDeposit(
		ctx context.Context,
		ethAddress common.Address,
		bidAmount string,
		blockNumber int64,
	) (func() error, error)
}

type Tracker interface {
	TrackCommitment(ctx context.Context, cm *store.EncryptedPreConfirmationWithDecrypted) error
}

type PreconfContract interface {
	StoreUnopenedCommitment(
		opts *bind.TransactOpts,
		commitmentDigest [32]byte,
		commitmentSignature []byte,
		dispatchTimestamp uint64,
	) (*types.Transaction, error)
}

type Encryptor interface {
	ConstructEncryptedBid(bid *preconfpb.Bid) (*preconfpb.EncryptedBid, *fr.Element, error)
	ConstructEncryptedPreConfirmation(*preconfpb.Bid) (*preconfpb.PreConfirmation, *preconfpb.EncryptedPreConfirmation, error)
	VerifyBid(*preconfpb.Bid) (*common.Address, error)
	VerifyEncryptedPreConfirmation(
		bid *preconfpb.Bid,
		providerNikePK *bn254.G1Affine,
		bidderNikeSC *fr.Element,
		c *preconfpb.EncryptedPreConfirmation,
	) ([]byte, *common.Address, error)
	DecryptBidData(common.Address, *preconfpb.EncryptedBid) (*preconfpb.Bid, error)
}

func New(
	topo Topology,
	streamer p2p.Streamer,
	encryptor Encryptor,
	depositMgr DepositManager,
	processor BidProcessor,
	commitmentDA PreconfContract,
	tracker Tracker,
	optsGetter OptsGetter,
	providerTimeout time.Duration,
	logger *slog.Logger,
) *Preconfirmation {
	return &Preconfirmation{
		topo:            topo,
		streamer:        streamer,
		encryptor:       encryptor,
		depositMgr:      depositMgr,
		processer:       processor,
		commitmentDA:    commitmentDA,
		tracker:         tracker,
		optsGetter:      optsGetter,
		logger:          logger,
		metrics:         newMetrics(),
		providerTimeout: providerTimeout,
	}
}

func (p *Preconfirmation) bidStream() p2p.StreamDesc {
	return p2p.StreamDesc{
		Name:    ProtocolName,
		Version: ProtocolVersion,
		Handler: p.handleBid,
	}
}

func (p *Preconfirmation) Streams() []p2p.StreamDesc {
	return []p2p.StreamDesc{p.bidStream()}
}

// SendBid is meant to be called by the bidder to construct and send bids to the provider.
// It takes the txHash, the bid amount in wei and the maximum valid block number.
// It waits for preConfirmations from all providers and then returns.
// It returns an error if the bid is not valid.
func (p *Preconfirmation) SendBid(
	ctx context.Context,
	bid *preconfpb.Bid,
) (chan *preconfpb.PreConfirmation, error) {
	startTime := time.Now()
	encryptedBid, nikePrivateKey, err := p.encryptor.ConstructEncryptedBid(bid)
	if err != nil {
		p.logger.Error("constructing encrypted bid", "error", err, "bid", bid)
		return nil, err
	}
	p.metrics.BidConstructDurationSummary.Observe(time.Since(startTime).Seconds())

	providers := p.topo.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	if len(providers) == 0 {
		p.logger.Error("no providers available", "bid", bid)
		return nil, errors.New("no providers available")
	}

	// Create a new channel to receive preConfirmations
	preConfirmations := make(chan *preconfpb.PreConfirmation, len(providers))

	wg := sync.WaitGroup{}
	for idx := range providers {
		wg.Add(1)
		go func(provider p2p.Peer) {
			defer wg.Done()
			logger := p.logger.With("provider", provider, "bid", bid.TxHash)

			providerStream, err := p.streamer.NewStream(
				ctx,
				provider,
				nil,
				p.bidStream(),
			)
			if err != nil {
				logger.Error("creating stream", "error", err)
				return
			}

			err = providerStream.WriteMsg(ctx, encryptedBid)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("writing message", "error", err)
				return
			}
			p.metrics.SentBidsCount.Inc()

			writeToReadStartTime := time.Now()
			encryptedPreConfirmation := new(preconfpb.EncryptedPreConfirmation)
			err = providerStream.ReadMsg(ctx, encryptedPreConfirmation)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("reading message", "error", err)
				return
			}
			writeToReadDuration := time.Since(writeToReadStartTime).Seconds()

			_ = providerStream.Close()

			// Process preConfirmation as a bidder
			verifyStartTime := time.Now()
			sharedSecretKey, providerAddress, err := p.encryptor.VerifyEncryptedPreConfirmation(
				bid,
				provider.Keys.NIKEPublicKey,
				nikePrivateKey,
				encryptedPreConfirmation,
			)
			if err != nil {
				logger.Error("verifying provider signature", "error", err)
				return
			}
			verifyDuration := time.Since(verifyStartTime).Seconds()
			p.metrics.VerifyPreconfDurationSummary.Observe(verifyDuration)

			wireLatency := time.Since(time.Unix(0, encryptedPreConfirmation.DispatchTimestamp)).Seconds()
			logger.Info(
				"successfully received preconf",
				"totalDuration", writeToReadDuration,
				"wireLatency", wireLatency,
			)

			preConfirmation := &preconfpb.PreConfirmation{
				Bid:               bid,
				SharedSecret:      sharedSecretKey,
				Digest:            encryptedPreConfirmation.Commitment,
				Signature:         encryptedPreConfirmation.Signature,
				DispatchTimestamp: encryptedPreConfirmation.DispatchTimestamp,
			}

			preConfirmation.ProviderAddress = make([]byte, len(providerAddress))
			copy(preConfirmation.ProviderAddress, providerAddress[:])

			encryptedAndDecryptedPreconfirmation := &store.EncryptedPreConfirmationWithDecrypted{
				EncryptedPreConfirmation: encryptedPreConfirmation,
				PreConfirmation:          preConfirmation,
			}

			p.metrics.ReceivedPreconfsCount.Inc()
			// Track the preconfirmation
			if err := p.tracker.TrackCommitment(ctx, encryptedAndDecryptedPreconfirmation); err != nil {
				logger.Error("tracking commitment", "error", err)
				return
			}

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

	encryptedBid := new(preconfpb.EncryptedBid)
	err := stream.ReadMsg(ctx, encryptedBid)
	if err != nil {
		return err
	}

	bid, err := p.encryptor.DecryptBidData(peer.EthAddress, encryptedBid)
	if err != nil {
		return err
	}
	ethAddress, err := p.encryptor.VerifyBid(bid)
	if err != nil {
		return err
	}

	refund, err := p.depositMgr.CheckAndDeductDeposit(ctx, *ethAddress, bid.BidAmount, bid.BlockNumber)
	if err != nil {
		p.logger.Error("checking deposit", "error", err)
		return err
	}

	// Setup defer for possible refund
	successful := false
	defer func() {
		if !successful {
			// Refund the deducted amount if the bid process did not succeed
			refundErr := refund()
			if refundErr != nil {
				p.logger.Error("refunding deposit", "error", refundErr)
			}
		}
	}()

	// try to get a decision within 30 seconds
	ctx, cancel := context.WithTimeout(ctx, p.providerTimeout)
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
			constructStartTime := time.Now()
			preConfirmation, encryptedPreConfirmation, err := p.encryptor.ConstructEncryptedPreConfirmation(bid)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to constuct encrypted preconfirmation: %v", err)
			}
			constructDuration := time.Since(constructStartTime).Seconds()
			p.metrics.ConstructPreconfDurationSummary.Observe(constructDuration)

			encryptedPreConfirmation.DispatchTimestamp = st.DispatchTimestamp

			err = stream.WriteMsg(ctx, encryptedPreConfirmation)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to send preconfirmation: %v", err)
			}
			var commitmentDigest [32]byte
			copy(commitmentDigest[:], encryptedPreConfirmation.Commitment)

			opts, err := p.optsGetter(ctx)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get transact opts: %v", err)
			}

			txn, err := p.commitmentDA.StoreUnopenedCommitment(
				opts,
				commitmentDigest,
				encryptedPreConfirmation.Signature,
				uint64(st.DispatchTimestamp),
			)
			if err != nil {
				p.logger.Error("storing commitment", "error", err)
				return status.Errorf(codes.Internal, "failed to store commitments: %v", err)
			}

			encryptedAndDecryptedPreconfirmation := &store.EncryptedPreConfirmationWithDecrypted{
				TxnHash:                  txn.Hash(),
				EncryptedPreConfirmation: encryptedPreConfirmation,
				PreConfirmation:          preConfirmation,
			}

			if err := p.tracker.TrackCommitment(ctx, encryptedAndDecryptedPreconfirmation); err != nil {
				p.logger.Error("tracking commitment", "error", err)
				return status.Errorf(codes.Internal, "failed to track commitment: %v", err)
			}

			// If we reach here, the bid was successful
			successful = true

			return nil
		}
	}
	return nil
}
