package handshake

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core"
	handshakepb "github.com/primev/mev-commit/p2p/gen/go/handshake/v1"
	p2pcrypto "github.com/primev/mev-commit/p2p/pkg/crypto"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	"github.com/primev/mev-commit/p2p/pkg/signer"
	"github.com/primev/mev-commit/x/keysigner"
)

const (
	ProtocolName    = "handshake"
	ProtocolVersion = "7.0.0"
	StreamName      = "handshake"
)

var (
	ErrSignatureVerificationFailed = errors.New("signature verification failed")
	ErrObservedAddressMismatch     = errors.New("observed address mismatch")
	ErrInsufficientStake           = errors.New("insufficient stake")
)

type ProviderRegistry interface {
	CheckProviderRegistered(context.Context, common.Address) bool
}

// Handshake is the handshake protocol
type Service struct {
	ks            keysigner.KeySigner
	peerType      p2p.PeerType
	passcode      string
	signer        signer.Signer
	providerKeys  *p2p.Keys
	register      ProviderRegistry
	handshakeReq  *handshakepb.HandshakeReq
	getEthAddress func(core.PeerID) (common.Address, error)
}

func New(
	ks keysigner.KeySigner,
	peerType p2p.PeerType,
	passcode string,
	signer signer.Signer,
	providerKeys *p2p.Keys,
	register ProviderRegistry,
	getEthAddress func(core.PeerID) (common.Address, error),
) (*Service, error) {
	s := &Service{
		ks:            ks,
		peerType:      peerType,
		passcode:      passcode,
		signer:        signer,
		providerKeys:  providerKeys,
		register:      register,
		getEthAddress: getEthAddress,
	}

	err := s.setHandshakeReq()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (h *Service) verifyReq(
	req *handshakepb.HandshakeReq,
	peerID core.PeerID,
) (common.Address, error) {
	unsignedData := []byte(req.PeerType + req.Token)

	verified, ethAddress, err := h.signer.Verify(req.Sig, unsignedData)
	if err != nil {
		return common.Address{}, errors.Join(err, ErrSignatureVerificationFailed)
	}

	if !verified {
		return common.Address{}, ErrSignatureVerificationFailed
	}

	observedEthAddress, err := h.getEthAddress(peerID)
	if err != nil {
		return common.Address{}, err
	}

	if !bytes.Equal(observedEthAddress.Bytes(), ethAddress.Bytes()) {
		return common.Address{}, ErrObservedAddressMismatch
	}

	if req.PeerType == p2p.PeerTypeProvider.String() {
		if !h.register.CheckProviderRegistered(context.Background(), ethAddress) {
			return common.Address{}, ErrInsufficientStake
		}
	}

	return ethAddress, nil
}

func (h *Service) createSignature() ([]byte, error) {
	unsignedData := []byte(h.peerType.String() + h.passcode)
	hash := crypto.Keccak256Hash(unsignedData)
	sig, err := h.ks.SignHash(hash.Bytes())
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (h *Service) setHandshakeReq() error {
	sig, err := h.createSignature()
	if err != nil {
		return err
	}

	req := &handshakepb.HandshakeReq{
		PeerType: h.peerType.String(),
		Token:    h.passcode,
		Sig:      sig,
	}

	if h.peerType == p2p.PeerTypeProvider {
		ppk := p2pcrypto.SerializeEciesPublicKey(h.providerKeys.PKEPublicKey)
		req.Keys = &handshakepb.SerializedKeys{
			PKEPublicKey:  ppk,
			NIKEPublicKey: h.providerKeys.NIKEPublicKey.Bytes(),
		}
	}

	h.handshakeReq = req
	return nil
}

func (h *Service) verifyResp(resp *handshakepb.HandshakeResp) error {
	if !bytes.Equal(resp.ObservedAddress, h.ks.GetAddress().Bytes()) {
		return errors.New("observed address mismatch")
	}

	if resp.PeerType != h.peerType.String() {
		return errors.New("peer type mismatch")
	}

	return nil
}

func (h *Service) Handle(
	ctx context.Context,
	stream p2p.Stream,
	peerID core.PeerID,
) (*p2p.Peer, error) {
	req := new(handshakepb.HandshakeReq)
	err := stream.ReadMsg(ctx, req)
	if err != nil {
		return nil, err
	}

	ethAddress, err := h.verifyReq(req, peerID)
	if err != nil {
		return nil, err
	}

	resp := &handshakepb.HandshakeResp{
		ObservedAddress: ethAddress.Bytes(),
		PeerType:        req.PeerType,
	}

	if err := stream.WriteMsg(ctx, resp); err != nil {
		return nil, err
	}

	err = stream.WriteMsg(ctx, h.handshakeReq)
	if err != nil {
		return nil, err
	}

	ack := new(handshakepb.HandshakeResp)
	err = stream.ReadMsg(ctx, ack)
	if err != nil {
		return nil, err
	}

	if err := h.verifyResp(ack); err != nil {
		return nil, err
	}

	p := &p2p.Peer{
		EthAddress: ethAddress,
		Type:       p2p.FromString(req.PeerType),
	}

	if req.PeerType == p2p.PeerTypeProvider.String() {
		ppk, err := p2pcrypto.DeserializeEciesPublicKey(req.Keys.PKEPublicKey)
		if err != nil {
			return &p2p.Peer{}, err
		}
		npk, err := ecdh.P256().NewPublicKey(req.Keys.NIKEPublicKey)
		if err != nil {
			return &p2p.Peer{}, err
		}
		p.Keys = &p2p.Keys{
			PKEPublicKey:  ppk,
			NIKEPublicKey: npk,
		}
	}
	return p, nil
}

func (h *Service) Handshake(
	ctx context.Context,
	peerID core.PeerID,
	stream p2p.Stream,
) (*p2p.Peer, error) {
	if err := stream.WriteMsg(ctx, h.handshakeReq); err != nil {
		return nil, err
	}

	resp := new(handshakepb.HandshakeResp)
	err := stream.ReadMsg(ctx, resp)
	if err != nil {
		return nil, err
	}

	if err := h.verifyResp(resp); err != nil {
		return nil, err
	}

	ack := new(handshakepb.HandshakeReq)
	err = stream.ReadMsg(ctx, ack)
	if err != nil {
		return nil, err
	}

	ethAddress, err := h.verifyReq(ack, peerID)
	if err != nil {
		return nil, err
	}

	err = stream.WriteMsg(ctx, &handshakepb.HandshakeResp{
		ObservedAddress: ethAddress.Bytes(),
		PeerType:        ack.PeerType,
	})
	if err != nil {
		return nil, err
	}

	p := &p2p.Peer{
		EthAddress: ethAddress,
		Type:       p2p.FromString(ack.PeerType),
	}

	if ack.PeerType == p2p.PeerTypeProvider.String() {
		ppk, err := p2pcrypto.DeserializeEciesPublicKey(ack.Keys.PKEPublicKey)
		if err != nil {
			return &p2p.Peer{}, err
		}
		npk, err := ecdh.P256().NewPublicKey(ack.Keys.NIKEPublicKey)
		if err != nil {
			return &p2p.Peer{}, err
		}
		p.Keys = &p2p.Keys{
			PKEPublicKey:  ppk,
			NIKEPublicKey: npk,
		}
	}

	return p, nil
}
