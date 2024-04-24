package preconfencryptor

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	preconfpb "github.com/primevprotocol/mev-commit/p2p/gen/go/preconfirmation/v1"
)

var EIPVerify = eipVerify

func (e *encryptor) BidOriginator(bid *preconfpb.Bid) (*common.Address, *ecdsa.PublicKey, error) {
	_, err := e.VerifyBid(bid)
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
