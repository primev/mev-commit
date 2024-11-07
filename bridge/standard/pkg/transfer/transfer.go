package transfer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	l1gateway "github.com/primev/mev-commit/contracts-abi/clients/L1Gateway"
	settlementgateway "github.com/primev/mev-commit/contracts-abi/clients/SettlementGateway"
	"github.com/primev/mev-commit/x/keysigner"
)

type Initiator interface {
	InitiateTransfer(
		opts *bind.TransactOpts,
		recipient common.Address,
		amount *big.Int,
	) (*types.Transaction, error)
}

type Filterer interface {
	ParseTransferIdx(log types.Log) (*big.Int, error)
	WaitForTransferFinalized(ctx context.Context, startBlock uint64, index *big.Int) error
}

type TransferStatus struct {
	Message string
	Error   error
}

type Transfer struct {
	signer           keysigner.KeySigner
	amount           *big.Int
	destAddress      common.Address
	srcClient        *ethclient.Client
	srcChainID       *big.Int
	srcTransactor    Initiator
	srcFilterer      Filterer
	destInitialBlock uint64
	destFilterer     Filterer
}

func NewTransferToSettlement(
	amount *big.Int,
	destAddress common.Address,
	signer keysigner.KeySigner,
	settlementRPCUrl string,
	l1RPCUrl string,
	l1ContractAddr common.Address,
	settlementContractAddr common.Address,
) (*Transfer, error) {
	l1Client, err := ethclient.Dial(l1RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial l1 rpc: %s", err)
	}
	l1ChainID, err := l1Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get l1 chain id: %s", err)
	}
	l1Transactor, err := l1gateway.NewL1gatewayTransactor(l1ContractAddr, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 transactor: %s", err)
	}
	l1Filterer, err := l1gateway.NewL1gatewayFilterer(l1ContractAddr, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 filterer: %s", err)
	}

	settlementClient, err := ethclient.Dial(settlementRPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial settlement rpc: %s", err)
	}
	initialBlock, err := settlementClient.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get initial block: %s", err)
	}
	settlementFilterer, err := settlementgateway.NewSettlementgatewayFilterer(settlementContractAddr, settlementClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement filterer: %s", err)
	}

	return &Transfer{
		signer:           signer,
		amount:           amount,
		destAddress:      destAddress,
		srcClient:        l1Client,
		srcChainID:       l1ChainID,
		srcTransactor:    l1Transactor,
		srcFilterer:      &L1Filterer{l1Filterer},
		destInitialBlock: initialBlock,
		destFilterer:     &SettlementFilterer{settlementFilterer},
	}, nil
}

func NewTransferToL1(
	amount *big.Int,
	destAddress common.Address,
	signer keysigner.KeySigner,
	settlementRPCUrl string,
	l1RPCUrl string,
	l1ContractAddr common.Address,
	settlementContractAddr common.Address,
) (*Transfer, error) {
	l1Client, err := ethclient.Dial(l1RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial l1 rpc: %s", err)
	}
	initialBlock, err := l1Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get initial block: %s", err)
	}
	l1Filterer, err := l1gateway.NewL1gatewayFilterer(l1ContractAddr, l1Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create l1 filterer: %s", err)
	}

	settlementClient, err := ethclient.Dial(settlementRPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial settlement rpc: %s", err)
	}
	settlementChainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get settlement chain id: %s", err)
	}
	settlementTransactor, err := settlementgateway.NewSettlementgatewayTransactor(
		settlementContractAddr,
		settlementClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement transactor: %s", err)
	}
	settlementFilterer, err := settlementgateway.NewSettlementgatewayFilterer(
		settlementContractAddr,
		settlementClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create settlement filterer: %s", err)
	}

	return &Transfer{
		amount:           amount,
		destAddress:      destAddress,
		signer:           signer,
		srcClient:        settlementClient,
		srcChainID:       settlementChainID,
		srcTransactor:    settlementTransactor,
		srcFilterer:      &SettlementFilterer{settlementFilterer},
		destInitialBlock: initialBlock,
		destFilterer:     &L1Filterer{l1Filterer},
	}, nil
}

func (t *Transfer) Do(ctx context.Context) <-chan TransferStatus {
	statusChan := make(chan TransferStatus)
	go func() {
		defer close(statusChan)

		statusChan <- TransferStatus{Message: "Setting up initial transfer..."}

		opts, err := t.signer.GetAuthWithCtx(ctx, t.srcChainID)
		if err != nil {
			statusChan <- TransferStatus{
				Message: "Failed to setup initial transfer",
				Error:   fmt.Errorf("failed to get auth: %s", err),
			}
			return
		}
		opts.Value = t.amount

		statusChan <- TransferStatus{
			Message: fmt.Sprintf("Initiating transfer of %s to %s", t.amount, t.destAddress.Hex()),
		}

		tx, err := t.srcTransactor.InitiateTransfer(
			opts,
			t.destAddress,
			t.amount,
		)
		if err != nil {
			statusChan <- TransferStatus{
				Message: "Failed to initiate transfer",
				Error:   fmt.Errorf("failed to send transaction: %s", err),
			}
			return
		}

		statusChan <- TransferStatus{
			Message: fmt.Sprintf("Transfer initiated with hash %s. Waiting for it to be mined...", tx.Hash().Hex()),
		}

		receipt, err := bind.WaitMined(ctx, t.srcClient, tx)
		if err != nil {
			statusChan <- TransferStatus{
				Message: "Failed to wait for transaction to be mined",
				Error:   fmt.Errorf("failed to wait for transaction to be mined: %s", err),
			}
			return
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			statusChan <- TransferStatus{
				Message: "Transaction failed",
				Error:   fmt.Errorf("transaction failed: %d", receipt.Status),
			}
			return
		}

		var transferIdx *big.Int
		for _, log := range receipt.Logs {
			if idx, err := t.srcFilterer.ParseTransferIdx(*log); err == nil {
				transferIdx = idx
				break
			}
		}

		if transferIdx == nil {
			statusChan <- TransferStatus{
				Message: "Failed to find transfer initiated event",
				Error:   fmt.Errorf("transfer initiated event not found"),
			}
			return
		}

		statusChan <- TransferStatus{
			Message: fmt.Sprintf("Transaction mined in block %d with index %s", receipt.BlockNumber, transferIdx),
		}

		statusChan <- TransferStatus{
			Message: "Waiting for transfer to be finalized...",
		}

		switch err := t.destFilterer.WaitForTransferFinalized(ctx, t.destInitialBlock, transferIdx); {
		case err == nil:
			statusChan <- TransferStatus{Message: "Transfer finalized. Bridging complete."}
		case errors.Is(err, context.DeadlineExceeded):
			statusChan <- TransferStatus{
				Message: "Timeout while waiting for transfer finalization",
				Error:   fmt.Errorf("timeout while waiting for transfer finalization: %s", err),
			}
		default:
			statusChan <- TransferStatus{
				Message: "Failed to wait for transfer finalization",
				Error:   fmt.Errorf("failed to finalize transfer: %s", err),
			}
			return
		}
	}()
	return statusChan
}

type L1Filterer struct {
	*l1gateway.L1gatewayFilterer
}

func (f *L1Filterer) WaitForTransferFinalized(ctx context.Context, startBlock uint64, index *big.Int) error {
	opts := &bind.FilterOpts{
		Start: startBlock,
		End:   nil,
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		iter, err := f.FilterTransferFinalized(opts, nil, []*big.Int{index})
		if err != nil {
			return fmt.Errorf("error obtaining transfer finalized event: %s", err)
		}
		if iter.Next() {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func (f *L1Filterer) ParseTransferIdx(log types.Log) (*big.Int, error) {
	ev, err := f.ParseTransferInitiated(log)
	if err != nil {
		return nil, err
	}
	return ev.TransferIdx, nil
}

type SettlementFilterer struct {
	*settlementgateway.SettlementgatewayFilterer
}

func (s *SettlementFilterer) WaitForTransferFinalized(ctx context.Context, startBlock uint64, index *big.Int) error {
	opts := &bind.FilterOpts{
		Start: startBlock,
		End:   nil,
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		iter, err := s.FilterTransferFinalized(opts, nil, []*big.Int{index})
		if err != nil {
			return fmt.Errorf("error obtaining transfer finalized event: %s", err)
		}
		if iter.Next() {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func (s *SettlementFilterer) ParseTransferIdx(log types.Log) (*big.Int, error) {
	ev, err := s.ParseTransferInitiated(log)
	if err != nil {
		return nil, err
	}
	return ev.TransferIdx, nil
}
