package mockkeysigner

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type MockKeySigner struct {
	privKey *ecdsa.PrivateKey
	address common.Address
}

func NewMockKeySigner(privKey *ecdsa.PrivateKey, address common.Address) *MockKeySigner {
	return &MockKeySigner{privKey: privKey, address: address}
}

func (m *MockKeySigner) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return tx, nil
}

func (m *MockKeySigner) SignHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, m.privKey)
}

func (m *MockKeySigner) GetAddress() common.Address {
	return m.address
}

func (m *MockKeySigner) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return m.privKey, nil
}

func (m *MockKeySigner) ZeroPrivateKey(key *ecdsa.PrivateKey) {}

func (m *MockKeySigner) String() string {
	return "mock"
}

func (m *MockKeySigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	return nil, nil
}

func (m *MockKeySigner) GetAuthWithCtx(_ context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	return nil, nil
}
