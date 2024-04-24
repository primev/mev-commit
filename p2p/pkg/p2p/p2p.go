package p2p

import (
	"context"
	"crypto/ecdh"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/primevprotocol/mev-commit/p2p/pkg/keykeeper"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// PeerType is the type of a peer
type PeerType int

const (
	// PeerTypeBootnode is a boot node
	PeerTypeBootnode PeerType = iota
	// PeerTypeProvider is a provider node
	PeerTypeProvider
	// PeerTypeBidder is a bidder node
	PeerTypeBidder
)

func (pt PeerType) String() string {
	switch pt {
	case PeerTypeBootnode:
		return "bootnode"
	case PeerTypeProvider:
		return "provider"
	case PeerTypeBidder:
		return "bidder"
	default:
		return "unknown"
	}
}

func FromString(str string) PeerType {
	switch str {
	case "bootnode":
		return PeerTypeBootnode
	case "provider":
		return PeerTypeProvider
	case "bidder":
		return PeerTypeBidder
	default:
		return -1
	}
}

var (
	ErrPeerNotFound = errors.New("peer not found")
	ErrNoAddresses  = errors.New("no addresses")
)

type Keys struct {
	PKEPublicKey  *ecies.PublicKey
	NIKEPublicKey *ecdh.PublicKey
}

type jsonKeys struct {
	PKEPublicKey  string `json:"pkePublicKey"`
	NIKEPublicKey string `json:"nikePublicKey"`
}

func (k *Keys) MarshalJSON() ([]byte, error) {
	ppk := keykeeper.SerializePublicKey(k.PKEPublicKey)
	pkePublicKeyB64 := base64.StdEncoding.EncodeToString(ppk)

	npk := k.NIKEPublicKey.Bytes()
	nikePublicKeyB64 := base64.StdEncoding.EncodeToString(npk)

	return json.Marshal(jsonKeys{
		PKEPublicKey:  pkePublicKeyB64,
		NIKEPublicKey: nikePublicKeyB64,
	})
}

func (k *Keys) UnmarshalJSON(data []byte) error {
	var jk jsonKeys
	if err := json.Unmarshal(data, &jk); err != nil {
		return err
	}

	pkePublicKeyBytes, err := base64.StdEncoding.DecodeString(jk.PKEPublicKey)
	if err != nil {
		return err
	}

	pkePublicKey, err := keykeeper.DeserializePublicKey(pkePublicKeyBytes)
	if err != nil {
		return err
	}

	nikePublicKeyBytes, err := base64.StdEncoding.DecodeString(jk.NIKEPublicKey)
	if err != nil {
		return err
	}
	nikePublicKey, err := ecdh.P256().NewPublicKey(nikePublicKeyBytes)
	if err != nil {
		return err
	}

	k.PKEPublicKey = pkePublicKey
	k.NIKEPublicKey = nikePublicKey

	return nil
}

type Peer struct {
	EthAddress common.Address
	Type       PeerType
	Keys       *Keys
}

type PeerInfo struct {
	EthAddress common.Address
	Underlay   []byte
}

// Stream is a bidirectional stream of messages between two peers per protocol.
type Stream interface {
	// ReadMsg reads the next message from the stream.
	ReadMsg(context.Context, proto.Message) error
	// WriteMsg writes a message to the stream.
	WriteMsg(context.Context, proto.Message) error

	Reset() error
	io.Closer
}

type MetadataStream interface {
	// ReadHeader reads the header from the stream.
	ReadHeader(context.Context) (Header, error)
	// WriteHeader writes the header to the stream.
	WriteHeader(context.Context, Header) error
	// WriteError writes an error to the stream.
	WriteError(context.Context, *status.Status) error
}

// Header is a map of string to structpb.Value. It is used to pass headers
// between the client and the server.
type Header map[string]*structpb.Value

// HandlerFunc is a function that handles a stream.
type HandlerFunc func(ctx context.Context, peer Peer, stream Stream) error

// HeaderFunc is a function that handles a header.
type HeaderFunc func(ctx context.Context, peer Peer, hdr Header) Header

// StreamDesc describes a stream handler.
type StreamDesc struct {
	Name    string
	Version string
	Handler HandlerFunc
	Header  HeaderFunc
}

type Addressbook interface {
	GetPeerInfo(Peer) ([]byte, error)
}

type Streamer interface {
	NewStream(context.Context, Peer, Header, StreamDesc) (Stream, error)
}

type Service interface {
	AddStreamHandlers(desc ...StreamDesc)
	Connect(ctx context.Context, info []byte) (Peer, error)
	Streamer
	Addressbook
	// Peers blocklisted by libp2p. Currently no external service needs the blocklist
	// so we don't expose it.
	BlockedPeers() []BlockedPeerInfo
	io.Closer
}

type Notifier interface {
	Connected(Peer)
	Disconnected(Peer)
}

type BlockedPeerInfo struct {
	Peer     common.Address
	Reason   string
	Duration string
}
