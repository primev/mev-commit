package debugapi

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"github.com/primev/mev-commit/p2p/pkg/txnstore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	debugapiv1.UnimplementedDebugServiceServer
	store     Store
	canceller Canceller
	p2p       P2PService
	topology  Topology
}

func NewService(
	store Store,
	canceller Canceller,
	p2p P2PService,
	topology Topology,
) *Service {
	return &Service{
		store:     store,
		canceller: canceller,
		p2p:       p2p,
		topology:  topology,
	}
}

type Store interface {
	PendingTxns() ([]*txnstore.TxnDetails, error)
}

type Canceller interface {
	CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error)
}

type P2PService interface {
	Self() map[string]interface{}
	BlockedPeers() []p2p.BlockedPeerInfo
}

type Topology interface {
	GetPeers(q topology.Query) []p2p.Peer
}

func (s *Service) GetTopology(
	ctx context.Context,
	_ *debugapiv1.EmptyMessage,
) (*debugapiv1.TopologyResponse, error) {
	providers := s.topology.GetPeers(topology.Query{Type: p2p.PeerTypeProvider})
	bidders := s.topology.GetPeers(topology.Query{Type: p2p.PeerTypeBidder})
	self := s.p2p.Self()
	blocked := s.p2p.BlockedPeers()

	resp := make(map[string]interface{})

	if len(providers) > 0 {
		connectedProviders := make([]interface{}, len(providers))
		for i, p := range providers {
			connectedProviders[i] = p.EthAddress.String()
		}
		resp["connected_providers"] = connectedProviders
	}

	if len(bidders) > 0 {
		connectedBidders := make([]interface{}, len(bidders))
		for i, b := range bidders {
			connectedBidders[i] = b.EthAddress.String()
		}
		resp["connected_bidders"] = connectedBidders
	}

	if len(self) > 0 {
		resp["self"] = self
	}

	if len(blocked) > 0 {
		blockedPeers := make([]interface{}, len(blocked))
		for i, b := range blocked {
			blockedPeers[i] = map[string]interface{}{
				"peer":   b.Peer.String(),
				"reason": b.Reason,
			}
		}
		resp["blocked_peers"] = blockedPeers
	}

	structResp, err := structpb.NewStruct(resp)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "creating response: %v", err)
	}

	return &debugapiv1.TopologyResponse{Topology: structResp}, nil
}

func (s *Service) GetPendingTransactions(
	ctx context.Context,
	_ *debugapiv1.EmptyMessage,
) (*debugapiv1.PendingTransactionsResponse, error) {
	txns, err := s.store.PendingTxns()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pending transactions: %v", err)
	}

	txnsMsg := make([]*debugapiv1.TransactionInfo, len(txns))
	for i, txn := range txns {
		txnsMsg[i] = &debugapiv1.TransactionInfo{
			TxHash:  txn.Hash.Hex(),
			Nonce:   int64(txn.Nonce),
			Created: time.Unix(txn.Created, 0).String(),
		}
	}

	return &debugapiv1.PendingTransactionsResponse{PendingTransactions: txnsMsg}, nil
}

func (s *Service) CancelTransaction(
	ctx context.Context,
	cancel *debugapiv1.CancelTransactionReq,
) (*debugapiv1.CancelTransactionResponse, error) {
	txHash := common.HexToHash(cancel.TxHash)
	cHash, err := s.canceller.CancelTx(ctx, txHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cancelling transaction: %v", err)
	}

	return &debugapiv1.CancelTransactionResponse{TxHash: cHash.Hex()}, nil
}
