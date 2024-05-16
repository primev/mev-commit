package keyexchange

import (
	"errors"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/signer"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"github.com/primev/mev-commit/x/keysigner"
)

// Protocol constants.
const (
	ProtocolName        = "keyexchange"
	ProtocolHandlerName = "timestampMessage"
	ProtocolVersion     = "1.0.0"
)

// Error declarations.
var (
	ErrSignatureVerificationFailed = errors.New("signature verification failed")
	ErrObservedAddressMismatch     = errors.New("observed address mismatch")
	ErrInvalidBidderTypeForMessage = errors.New("invalid bidder type for message")
	ErrNoProvidersAvailable        = errors.New("no providers available")
)

// KeyExchange manages the key exchange process.
type KeyExchange struct {
	keySigner keysigner.KeySigner
	topo      Topology
	streamer  p2p.Streamer
	signer    signer.Signer
	store     Store
	logger    *slog.Logger
}

// Topology interface to get peers.
type Topology interface {
	GetPeers(topology.Query) []p2p.Peer
}

type Store interface {
	SetAESKey(common.Address, []byte) error
	GetAESKey(common.Address) ([]byte, error)
	GetECIESPrivateKey() (*ecies.PrivateKey, error)
}
