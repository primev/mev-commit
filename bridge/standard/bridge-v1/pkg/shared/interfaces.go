package shared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type GatewayTransactor interface {
	InitiateTransfer(opts *bind.TransactOpts, _recipient common.Address,
		amount *big.Int) (*types.Transaction, error)
	FinalizeTransfer(opts *bind.TransactOpts, _recipient common.Address,
		_amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error)
}

type GatewayFilterer interface {
	ObtainTransferInitiatedEvents(opts *bind.FilterOpts,
	) ([]TransferInitiatedEvent, error)
	ObtainTransferInitiatedBySender(opts *bind.FilterOpts, sender common.Address,
	) (TransferInitiatedEvent, error)
	ObtainTransferFinalizedEvent(opts *bind.FilterOpts, counterpartyIdx *big.Int) (
		TransferFinalizedEvent, bool, error)
}
