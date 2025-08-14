package setcode

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/primev/mev-commit/x/keysigner"
)

type SetCodeHelper struct {
	logger  *slog.Logger
	signer  keysigner.KeySigner
	backend bind.ContractBackend
	chainID *big.Int
}

func NewSetCodeHelper(
	logger *slog.Logger,
	signer keysigner.KeySigner,
	backend bind.ContractBackend,
	chainID *big.Int,
) *SetCodeHelper {
	return &SetCodeHelper{
		logger:  logger,
		signer:  signer,
		backend: backend,
		chainID: chainID,
	}
}

func (s *SetCodeHelper) SetCode(
	ctx context.Context,
	opts *bind.TransactOpts,
	to common.Address,
) (*types.Transaction, error) {

	if opts.GasTipCap == nil {
		return nil, fmt.Errorf("gas tip cap is required")
	}
	if opts.GasFeeCap == nil {
		return nil, fmt.Errorf("gas fee cap is required")
	}
	if opts.GasLimit == 0 {
		return nil, fmt.Errorf("gas limit is required")
	}

	// TODO: Create our own SetCodeAuthorization signing library compatible with KeySigner.
	// For now we use geth's EIP-7702 library that only supports ecdsa.PrivateKey
	pk, err := s.signer.GetPrivateKey()
	if err != nil {
		s.logger.Error("error getting private key from keysigner", "error", err)
		return nil, err
	}
	defer s.signer.ZeroPrivateKey(pk)

	nonce, err := s.backend.PendingNonceAt(ctx, s.signer.GetAddress())
	if err != nil {
		s.logger.Error("error getting pending nonce", "error", err)
		return nil, err
	}

	auth := types.SetCodeAuthorization{
		Address: to,
		Nonce:   nonce,
		ChainID: *uint256.MustFromBig(s.chainID),
	}

	signedAuth, err := types.SignSetCode(pk, auth)
	if err != nil {
		s.logger.Error("error signing set code authorization", "error", err)
		return nil, err
	}

	txdata := &types.SetCodeTx{
		ChainID:    uint256.MustFromBig(s.chainID),
		Nonce:      nonce,
		GasTipCap:  uint256.MustFromBig(opts.GasTipCap),
		GasFeeCap:  uint256.MustFromBig(opts.GasFeeCap),
		Gas:        opts.GasLimit,
		To:         common.Address{},
		Value:      uint256.NewInt(0),
		Data:       nil,
		AccessList: nil,
		AuthList:   []types.SetCodeAuthorization{signedAuth},
	}

	signer := types.NewPragueSigner(s.chainID)
	signedTx, err := types.SignNewTx(pk, signer, txdata)
	if err != nil {
		s.logger.Error("error signing transaction", "error", err)
		return nil, err
	}
	err = s.backend.SendTransaction(ctx, signedTx)
	if err != nil {
		s.logger.Error("error sending transaction", "error", err)
		return nil, err
	}

	return signedTx, nil
}
