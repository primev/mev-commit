package shared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Chain int

const (
	Settlement Chain = iota
	L1
)

func (c Chain) String() string {
	switch c {
	case Settlement:
		return "Settlement"
	case L1:
		return "L1"
	default:
		return "unknown"
	}
}

type TransferInitiatedEvent struct {
	Sender      common.Address
	Recipient   common.Address
	Amount      *big.Int
	TransferIdx *big.Int
	Chain       Chain
}

func (t TransferInitiatedEvent) String() string {
	return "Sender: " + t.Sender.String() +
		" Recipient: " + t.Recipient.String() +
		" Amount: " + t.Amount.String() +
		" TransferIdx: " + t.TransferIdx.String() +
		" Chain: " + t.Chain.String()
}

type TransferFinalizedEvent struct {
	Recipient       common.Address
	Amount          *big.Int
	CounterpartyIdx *big.Int
	Chain           Chain
}

func (t TransferFinalizedEvent) String() string {
	return "Recipient: " + t.Recipient.String() +
		" Amount: " + t.Amount.String() +
		" CounterpartyIdx: " + t.CounterpartyIdx.String() +
		" Chain: " + t.Chain.String()
}
