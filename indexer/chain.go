package main

import (
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

// Chain defines an interface for interacting with a blockchain.
type Chain interface {
	// BlockNumber returns the current block number.
	BlockNumber() (*big.Int, error)

	// Blocks returns [start, end] blocks.
	Blocks(start, end *big.Int) ([]*types.Block, error)

	// TxReceipts returns the receipts for the given transaction hashes if they exist.
	TxReceipts(txHashes []common.Hash) (map[string]*types.Receipt, error)
}

var (
	_ Chain     = (*EthereumChain)(nil)
	_ io.Closer = (*EthereumChain)(nil)
)

// EthereumChain implements the Chain interface and
// represents a connection to the Ethereum blockchain.
type EthereumChain struct {
	client *w3.Client
}

// NewEthereumChain creates a new EthereumChain instance connected to the specified URL.
func NewEthereumChain(url string) (*EthereumChain, error) {
	client, err := w3.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthereumChain{client: client}, nil
}

// BlockNumber implements the Chain.BlockNumber method.
func (e *EthereumChain) BlockNumber() (*big.Int, error) {
	var number big.Int
	if err := e.client.Call(eth.BlockNumber().Returns(&number)); err != nil {
		return nil, err
	}
	return &number, nil
}

// Blocks implements the Chain.Blocks method.
func (e *EthereumChain) Blocks(start, end *big.Int) ([]*types.Block, error) {
	var (
		window  = new(big.Int).Sub(end, start)
		length  = new(big.Int).Abs(window).Uint64() + 1
		numbers = make([]*big.Int, length)
	)

	switch window.Sign() {
	case -1:
		for i := 0; i < len(numbers); i++ {
			numbers[i] = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(i)))
		}
	case 1:
		for i := 0; i < len(numbers); i++ {
			numbers[i] = new(big.Int).Add(start, new(big.Int).SetUint64(uint64(i)))
		}
	default:
		return nil, nil
	}

	var (
		batch  = make([]w3types.RPCCaller, len(numbers))
		blocks = make([]*types.Block, len(numbers))
	)

	for i := range numbers {
		block := new(types.Block)
		batch[i] = eth.BlockByNumber(numbers[i]).Returns(block)
		blocks[i] = block
	}

	if err := e.client.Call(batch...); err != nil {
		return nil, err
	}
	return blocks, nil
}

// TxReceipts implements the Chain.TxReceipts method.
func (e *EthereumChain) TxReceipts(txHashes []common.Hash) (map[string]*types.Receipt, error) {
	var (
		batch      = make([]w3types.RPCCaller, len(txHashes))
		txReceipts = make([]types.Receipt, len(txHashes))
	)

	for i, txHash := range txHashes {
		batch[i] = eth.TxReceipt(txHash).Returns(&txReceipts[i])
	}

	if err := e.client.Call(batch...); err != nil {
		return map[string]*types.Receipt{}, nil
	}

	res := make(map[string]*types.Receipt, len(txReceipts))
	for _, txReceipt := range txReceipts {
		res[txReceipt.TxHash.Hex()] = &txReceipt
	}
	return res, nil
}

// Close closes the connection to the Ethereum blockchain.
// It implements the io.Closer interface.
func (e *EthereumChain) Close() error {
	return e.client.Close()
}
