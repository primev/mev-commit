package shared

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	sg "github.com/primevprotocol/contracts-abi/clients/SettlementGateway"
)

type SettlementFilterer struct {
	*sg.SettlementgatewayFilterer
}

func NewSettlementFilterer(
	gatewayAddr common.Address,
	client *ethclient.Client,
) (*SettlementFilterer, error) {
	f, err := sg.NewSettlementgatewayFilterer(gatewayAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement filterer: %w", err)
	}
	return &SettlementFilterer{f}, nil
}

func (f *SettlementFilterer) ObtainTransferInitiatedBySender(
	opts *bind.FilterOpts,
	sender common.Address,
) (TransferInitiatedEvent, error) {
	iter, err := f.FilterTransferInitiated(opts, []common.Address{sender}, nil, nil)
	if err != nil {
		return TransferInitiatedEvent{}, fmt.Errorf("failed to filter transfer initiated: %w", err)
	}
	if !iter.Next() {
		return TransferInitiatedEvent{}, fmt.Errorf("failed to obtain single transfer initiated event with sender: " + sender.String())
	}
	return TransferInitiatedEvent{
		Sender:      iter.Event.Sender,
		Recipient:   iter.Event.Recipient,
		Amount:      iter.Event.Amount,
		TransferIdx: iter.Event.TransferIdx,
		Chain:       Settlement,
	}, nil
}

func (f *SettlementFilterer) ObtainTransferInitiatedEvents(
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
			Chain:       Settlement,
		})
	}
	return toReturn, nil
}

func (f *SettlementFilterer) ObtainTransferFinalizedEvent(
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
			Chain:           Settlement,
		})
	}
	if len(events) == 0 {
		return TransferFinalizedEvent{}, false, nil
	}
	return events[0], true, nil
}
