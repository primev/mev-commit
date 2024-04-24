package preconfirmation

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	blocktracker "github.com/primevprotocol/mev-commit/contracts-abi/clients/BlockTracker"
	preconfcommstore "github.com/primevprotocol/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
	providerapiv1 "github.com/primevprotocol/mev-commit/p2p/gen/go/providerapi/v1"
	blocktrackercontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/block_tracker"
	preconfcontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/preconf"
	"github.com/primevprotocol/mev-commit/p2p/pkg/events"
	"github.com/primevprotocol/mev-commit/p2p/pkg/p2p"
	providerapi "github.com/primevprotocol/mev-commit/p2p/pkg/rpc/provider"
	encryptor "github.com/primevprotocol/mev-commit/p2p/pkg/signer/preconfencryptor"
	"github.com/primevprotocol/mev-commit/p2p/pkg/store"
	"github.com/primevprotocol/mev-commit/p2p/pkg/topology"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProtocolName    = "preconfirmation"
	ProtocolVersion = "3.0.0"
)

type Preconfirmation struct {
	owner        common.Address
	encryptor    encryptor.Encryptor
	topo         Topology
	streamer     p2p.Streamer
	depositMgr   DepositManager
	processer    BidProcessor
	commitmentDA preconfcontract.Interface
	blockTracker blocktrackercontract.Interface
	evtMgr       events.EventManager
	ecds         EncrDecrCommitmentStore
	newL1Blocks  chan *blocktracker.BlocktrackerNewL1Block
	enryptedCmts chan *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored
	logger       *slog.Logger
	metrics      *metrics
}

type Topology interface {
	GetPeers(topology.Query) []p2p.Peer
}

type BidProcessor interface {
	ProcessBid(context.Context, *preconfpb.Bid) (chan providerapi.ProcessedBidResponse, error)
}

type EncrDecrCommitmentStore interface {
	GetCommitmentsByBlockNumber(blockNum int64) ([]*store.EncryptedPreConfirmationWithDecrypted, error)
	GetCommitmentByHash(commitmentHash string) (*store.EncryptedPreConfirmationWithDecrypted, error)
	AddCommitment(commitment *store.EncryptedPreConfirmationWithDecrypted)
	DeleteCommitmentByBlockNumber(blockNum int64) error
	SetCommitmentIndexByCommitmentDigest(commitmentDigest, commitmentIndex [32]byte) error
}

type DepositManager interface {
	Start(ctx context.Context) <-chan struct{}
	CheckAndDeductDeposit(ctx context.Context, ethAddress common.Address, bidAmount string, blockNumber int64) (*big.Int, error)
	RefundDeposit(ethAddress common.Address, amount *big.Int, blockNumber int64) error
}

type EncrDecrCommitmentStore interface {
	GetCommitmentsByBlockNumber(blockNum int64) ([]*store.EncryptedPreConfirmationWithDecrypted, error)
	GetCommitmentByHash(commitmentHash string) (*store.EncryptedPreConfirmationWithDecrypted, error)
	AddCommitment(commitment *store.EncryptedPreConfirmationWithDecrypted)
	DeleteCommitmentByBlockNumber(blockNum int64) error
	SetCommitmentIndexByCommitmentDigest(commitmentDigest, commitmentIndex [32]byte) error
}

type DepositManager interface {
	Start(ctx context.Context) <-chan struct{}
	CheckAndDeductDeposit(ctx context.Context, ethAddress common.Address, bidAmount string, blockNumber int64) (*big.Int, error)
	RefundDeposit(ethAddress common.Address, amount *big.Int, blockNumber int64) error
}

func New(
	owner common.Address,
	topo Topology,
	streamer p2p.Streamer,
	encryptor encryptor.Encryptor,
	depositMgr DepositManager,
	processor BidProcessor,
	commitmentDA preconfcontract.Interface,
	blockTracker blocktrackercontract.Interface,
	evtMgr events.EventManager,
	edcs EncrDecrCommitmentStore,
	logger *slog.Logger,
) *Preconfirmation {
	return &Preconfirmation{
		owner:        owner,
		topo:         topo,
		streamer:     streamer,
		encryptor:    encryptor,
		depositMgr:   depositMgr,
		processer:    processor,
		commitmentDA: commitmentDA,
		blockTracker: blockTracker,
		evtMgr:       evtMgr,
		ecds:         edcs,
		newL1Blocks:  make(chan *blocktracker.BlocktrackerNewL1Block),
		enryptedCmts: make(chan *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored),
		logger:       logger,
		metrics:      newMetrics(),
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

func (p *Preconfirmation) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		ev1 := events.NewEventHandler(
			"NewL1Block",
			func(newL1Block *blocktracker.BlocktrackerNewL1Block) error {
				select {
				case <-egCtx.Done():
					return nil
				case p.newL1Blocks <- newL1Block:
					return nil
				}
			},
		)

		sub1, err := p.evtMgr.Subscribe(ev1)
		if err != nil {
			return fmt.Errorf("failed to subscribe to NewL1Block event: %w", err)
		}
		defer sub1.Unsubscribe()

		ev2 := events.NewEventHandler(
			"EncryptedCommitmentStored",
			func(ec *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored) error {
				select {
				case <-egCtx.Done():
					return nil
				case p.enryptedCmts <- ec:
					return nil
				}
			},
		)
		sub2, err := p.evtMgr.Subscribe(ev2)
		if err != nil {
			return fmt.Errorf("failed to subscribe to EncryptedCommitmentStored event: %w", err)
		}
		defer sub2.Unsubscribe()

		select {
		case <-egCtx.Done():
			return nil
		case err := <-sub1.Err():
			return fmt.Errorf("NewL1Block subscription error: %w", err)
		case err := <-sub2.Err():
			return fmt.Errorf("EncryptedCommitmentStored subscription error: %w", err)
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-egCtx.Done():
				return nil
			case newL1Block := <-p.newL1Blocks:
				if err := p.handleNewL1Block(egCtx, newL1Block); err != nil {
					return err
				}
			case ec := <-p.enryptedCmts:
				if err := p.handleEncryptedCommitmentStored(egCtx, ec); err != nil {
					return err
				}
			}
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			p.logger.Error("failed to start preconfirmation", "error", err)
		}
	}()

	return doneChan
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
	startTime := time.Now()
	bid, encryptedBid, err := p.encryptor.ConstructEncryptedBid(txHash, bidAmt, blockNumber, decayStartTimestamp, decayEndTimestamp)
	if err != nil {
		p.logger.Error("constructing encrypted bid", "error", err, "txHash", txHash)
		return nil, err
	}
	duration := time.Since(startTime)
	p.logger.Info("constructed encrypted bid", "encryptedBid", encryptedBid, "duration", duration)

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
				p.bidStream(),
			)
			if err != nil {
				logger.Error("creating stream", "error", err)
				return
			}

			logger.Info("sending encrypted bid", "encryptedBid", encryptedBid)

			err = providerStream.WriteMsg(ctx, encryptedBid)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("writing message", "error", err)
				return
			}
			p.metrics.SentBidsCount.Inc()

			encryptedPreConfirmation := new(preconfpb.EncryptedPreConfirmation)
			err = providerStream.ReadMsg(ctx, encryptedPreConfirmation)
			if err != nil {
				_ = providerStream.Reset()
				logger.Error("reading message", "error", err)
				return
			}

			_ = providerStream.Close()

			// Process preConfirmation as a bidder
			verifyStartTime := time.Now()
			sharedSecretKey, providerAddress, err := p.encryptor.VerifyEncryptedPreConfirmation(provider.Keys.NIKEPublicKey, bid.Digest, encryptedPreConfirmation)
			if err != nil {
				logger.Error("verifying provider signature", "error", err)
				return
			}
			verifyDuration := time.Since(verifyStartTime)
			logger.Info("verified encrypted preconfirmation", "duration", verifyDuration)

			preConfirmation := &preconfpb.PreConfirmation{
				Bid:          bid,
				SharedSecret: sharedSecretKey,
				Digest:       encryptedPreConfirmation.Commitment,
				Signature:    encryptedPreConfirmation.Signature,
			}

			preConfirmation.ProviderAddress = make([]byte, len(providerAddress))
			copy(preConfirmation.ProviderAddress, providerAddress[:])

			encryptedAndDecryptedPreconfirmation := &store.EncryptedPreConfirmationWithDecrypted{
				EncryptedPreConfirmation: encryptedPreConfirmation,
				PreConfirmation:          preConfirmation,
			}

			p.ecds.AddCommitment(encryptedAndDecryptedPreconfirmation)
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

	encryptedBid := new(preconfpb.EncryptedBid)
	err := stream.ReadMsg(ctx, encryptedBid)
	if err != nil {
		return err
	}

	p.logger.Info("received bid", "encryptedBid", encryptedBid)
	bid, err := p.encryptor.DecryptBidData(peer.EthAddress, encryptedBid)
	if err != nil {
		return err
	}
	ethAddress, err := p.encryptor.VerifyBid(bid)
	if err != nil {
		return err
	}

	deductedAmount, err := p.depositMgr.CheckAndDeductDeposit(ctx, *ethAddress, bid.BidAmount, bid.BlockNumber)
	if err != nil {
		p.logger.Error("checking deposit", "error", err)
		return err
	}

	// Setup defer for possible refund
	successful := false
	defer func() {
		if !successful {
			// Refund the deducted amount if the bid process did not succeed
			refundErr := p.depositMgr.RefundDeposit(*ethAddress, deductedAmount, bid.BlockNumber)
			if refundErr != nil {
				p.logger.Error("refunding deposit", "error", refundErr)
			}
		}
	}()

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
			preConfirmation, encryptedPreConfirmation, err := p.encryptor.ConstructEncryptedPreConfirmation(bid)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to constuct encrypted preconfirmation: %v", err)
			}
			p.logger.Info("sending preconfirmation", "preConfirmation", encryptedPreConfirmation)
			_, err = p.commitmentDA.StoreEncryptedCommitment(
				ctx,
				encryptedPreConfirmation.Commitment,
				encryptedPreConfirmation.Signature,
				uint64(st.DispatchTimestamp),
			)
			if err != nil {
				p.logger.Error("storing commitment", "error", err)
				return status.Errorf(codes.Internal, "failed to store commitments: %v", err)
			}

			encryptedAndDecryptedPreconfirmation := &store.EncryptedPreConfirmationWithDecrypted{
				EncryptedPreConfirmation: encryptedPreConfirmation,
				PreConfirmation:          preConfirmation,
			}

			p.ecds.AddCommitment(encryptedAndDecryptedPreconfirmation)

			// If we reach here, the bid was successful
			successful = true

			return stream.WriteMsg(ctx, encryptedPreConfirmation)
		}
	}
	return nil
}

func (p *Preconfirmation) handleNewL1Block(ctx context.Context, newL1Block *blocktracker.BlocktrackerNewL1Block) error {
	p.logger.Info("New L1 Block event received", "blockNumber", newL1Block.BlockNumber, "winner", newL1Block.Winner, "window", newL1Block.Window)
	commitments, err := p.ecds.GetCommitmentsByBlockNumber(newL1Block.BlockNumber.Int64())
	if err != nil {
		p.logger.Error("failed to get commitments by block number", "error", err)
		return err
	}
	for _, commitment := range commitments {
		if common.BytesToAddress(commitment.ProviderAddress) != newL1Block.Winner {
			p.logger.Info("provider address does not match the winner", "providerAddress", commitment.ProviderAddress, "winner", newL1Block.Winner)
			continue
		}
		startTime := time.Now()
		txHash, err := p.commitmentDA.OpenCommitment(
			ctx,
			commitment.EncryptedPreConfirmation.CommitmentIndex,
			commitment.PreConfirmation.Bid.BidAmount,
			commitment.PreConfirmation.Bid.BlockNumber,
			commitment.PreConfirmation.Bid.TxHash,
			commitment.PreConfirmation.Bid.DecayStartTimestamp,
			commitment.PreConfirmation.Bid.DecayEndTimestamp,
			commitment.PreConfirmation.Bid.Signature,
			commitment.PreConfirmation.Signature,
			commitment.PreConfirmation.SharedSecret,
		)
		if err != nil {
			// todo: retry mechanism?
			p.logger.Error("failed to open commitment", "error", err)
			continue
		}
		duration := time.Since(startTime)
		p.logger.Info("opened commitment", "txHash", txHash, "duration", duration)
	}

	err = p.ecds.DeleteCommitmentByBlockNumber(newL1Block.BlockNumber.Int64())
	if err != nil {
		p.logger.Error("failed to delete commitments by block number", "error", err)
		return err
	}
	return nil
}

func (p *Preconfirmation) handleEncryptedCommitmentStored(ctx context.Context, ec *preconfcommstore.PreconfcommitmentstoreEncryptedCommitmentStored) error {
	p.logger.Info("Encrypted Commitment Stored event received", "commitmentDigest", ec.CommitmentDigest, "commitmentIndex", ec.CommitmentIndex)
	return p.ecds.SetCommitmentIndexByCommitmentDigest(ec.CommitmentDigest, ec.CommitmentIndex)
}
