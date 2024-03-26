package shared

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	l1g "github.com/primevprotocol/contracts-abi/clients/L1Gateway"
)

type L1Filterer struct {
	*l1g.L1gatewayFilterer
}

func NewL1Filterer(
	gatewayAddr common.Address,
	client *ethclient.Client,
) (*L1Filterer, error) {
	f, err := l1g.NewL1gatewayFilterer(gatewayAddr, client)
	if err != nil {
		return nil, err
	}
	return &L1Filterer{f}, nil
}

func (f *L1Filterer) ObtainTransferInitiatedBySender(
	opts *bind.FilterOpts,
	sender common.Address,
) (TransferInitiatedEvent, error) {
	iter, err := f.FilterTransferInitiated(opts, []common.Address{sender}, nil, nil)
	if err != nil {
		return TransferInitiatedEvent{}, fmt.Errorf("failed to filter transfer initiated: %w", err)
	}
	if !iter.Next() {
		return TransferInitiatedEvent{}, fmt.Errorf("failed to obtain single transfer initiated event with sender: %s", sender.String())
	}
	return TransferInitiatedEvent{
		Sender:      iter.Event.Sender,
		Recipient:   iter.Event.Recipient,
		Amount:      iter.Event.Amount,
		TransferIdx: iter.Event.TransferIdx,
		Chain:       L1,
	}, nil
}

func (f *L1Filterer) ObtainTransferInitiatedEvents(
	opts *bind.FilterOpts,
) ([]TransferInitiatedEvent, error) {
	iter, err := f.FilterTransferInitiated(opts, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to filter transfer initiated: %w", err)
	}
	toReturn := make([]TransferInitiatedEvent, 0)
	for iter.Next() {
		toReturn = append(toReturn, TransferInitiatedEvent{
			Sender:      iter.Event.Sender,
			Recipient:   iter.Event.Recipient,
			Amount:      iter.Event.Amount,
			TransferIdx: iter.Event.TransferIdx,
			Chain:       L1,
		})
	}
	return toReturn, nil
}

func (f *L1Filterer) ObtainTransferFinalizedEvent(
	opts *bind.FilterOpts,
	counterpartyIdx *big.Int,
) (TransferFinalizedEvent, bool, error) {
	iter, err := f.FilterTransferFinalized(opts, nil, []*big.Int{counterpartyIdx})
	if err != nil {
		return TransferFinalizedEvent{}, false, fmt.Errorf("failed to filter transfer finalized: %w", err)
	}
	events := make([]TransferFinalizedEvent, 0)
	for iter.Next() {
		events = append(events, TransferFinalizedEvent{
			Recipient:       iter.Event.Recipient,
			Amount:          iter.Event.Amount,
			CounterpartyIdx: iter.Event.CounterpartyIdx,
			Chain:           L1,
		})
	}
	if len(events) == 0 {
		return TransferFinalizedEvent{}, false, nil
	}
	return events[0], true, nil
}
