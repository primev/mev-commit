package preconfsigner

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
)

var EIPVerify = eipVerify

func (p *privateKeySigner) BidOriginator(bid *preconfpb.Bid) (*common.Address, *ecdsa.PublicKey, error) {
	_, err := p.VerifyBid(bid)
	if err != nil {
		return nil, nil, err
	}

	sig := make([]byte, len(bid.Signature))
	copy(sig, bid.Signature)
	if sig[64] >= 27 && sig[64] <= 28 {
		sig[64] -= 27
	}

	pubkey, err := crypto.SigToPub(bid.Digest, sig)
	if err != nil {
		return nil, nil, err
	}

	address := crypto.PubkeyToAddress(*pubkey)

	return &address, pubkey, nil
}

func (p *privateKeySigner) PreConfirmationOriginator(
	c *preconfpb.PreConfirmation,
) (*common.Address, *ecdsa.PublicKey, error) {
	_, err := p.VerifyPreConfirmation(c)
	if err != nil {
		return nil, nil, err
	}

	sig := make([]byte, len(c.Signature))
	copy(sig, c.Signature)
	if sig[64] >= 27 && sig[64] <= 28 {
		sig[64] -= 27
	}
	pubkey, err := crypto.SigToPub(c.Digest, sig)
	if err != nil {
		return nil, nil, err
	}

	address := crypto.PubkeyToAddress(*pubkey)

	return &address, pubkey, nil
}
