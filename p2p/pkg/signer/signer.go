package signer

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	// create signature
	Sign(*ecdsa.PrivateKey, []byte) ([]byte, error)
	// verify signature and return signer address
	Verify([]byte, []byte) (bool, common.Address, error)
}

type signer struct{}

func New() Signer {
	return &signer{}
}

func (s *signer) Sign(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	messageHash := crypto.Keccak256Hash(message)
	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func (s *signer) Verify(signature []byte, message []byte) (bool, common.Address, error) {
	messageHash := crypto.Keccak256Hash(message)
	signaturePublicKey, err := crypto.SigToPub(messageHash.Bytes(), signature)
	if err != nil {
		return false, common.Address{}, err
	}

	signerAddress := crypto.PubkeyToAddress(*signaturePublicKey)

	verified := crypto.VerifySignature(
		crypto.FromECDSAPub(signaturePublicKey),
		messageHash.Bytes(),
		signature[:len(signature)-1],
	)
	return verified, signerAddress, nil
}
