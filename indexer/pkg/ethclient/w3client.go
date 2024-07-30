package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
	"github.com/primev/mev-commit/indexer/pkg/store"
)

type W3EvmClient struct {
	client *w3.Client
}

func NewW3EthereumClient(endpoint string) (*W3EvmClient, error) {
	client, err := w3.Dial(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	return &W3EvmClient{client: client}, nil
}

func (c *W3EvmClient) GetBlocks(ctx context.Context, blockNums []*big.Int) ([]*types.Block, error) {
	batchBlocksCaller := make([]w3types.RPCCaller, len(blockNums))
	blocks := make([]*types.Block, len(blockNums))

	for i, blockNum := range blockNums {
		block := new(types.Block)
		batchBlocksCaller[i] = eth.BlockByNumber(blockNum).Returns(block)
		blocks[i] = block
	}
	err := c.client.Call(batchBlocksCaller...)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (c *W3EvmClient) TxReceipts(ctx context.Context, txHashes []string) (map[string]*types.Receipt, error) {
	batchTxReceiptCaller := make([]w3types.RPCCaller, len(txHashes))
	txReceipts := make([]types.Receipt, len(txHashes))
	for i, txHash := range txHashes {
		batchTxReceiptCaller[i] = eth.TxReceipt(w3.H(txHash)).Returns(&txReceipts[i])
	}
	err := c.client.Call(batchTxReceiptCaller...)
	if err != nil {
		return map[string]*types.Receipt{}, nil
	}
	txHashToReceipt := make(map[string]*types.Receipt)
	for _, txReceipt := range txReceipts {
		txHashToReceipt[txReceipt.TxHash.Hex()] = &txReceipt
	}
	return txHashToReceipt, nil
}

func (c *W3EvmClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var block types.Block
	if err := c.client.Call(eth.BlockByNumber(number).Returns(&block)); err != nil {
		return nil, err
	}
	return &block, nil
}

func (c *W3EvmClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	var blockNumber big.Int
	if err := c.client.Call(eth.BlockNumber().Returns(&blockNumber)); err != nil {
		return nil, err
	}
	return &blockNumber, nil
}

func (c *W3EvmClient) AccountBalances(ctx context.Context, addresses []common.Address, blockNumber uint64) ([]store.AccountBalance, error) {
	batchAccountBalanceCaller := make([]w3types.RPCCaller, len(addresses))
	balances := make([]big.Int, len(addresses))
	for i, address := range addresses {
		batchAccountBalanceCaller[i] = eth.Balance(address, new(big.Int).SetUint64(blockNumber)).Returns(&balances[i])
	}
	err := c.client.Call(batchAccountBalanceCaller...)
	if err != nil {
		return nil, err
	}
	var accountBalances []store.AccountBalance
	for i, address := range addresses {
		accBalance := store.AccountBalance{
			Address:     address.Hex(),
			Balance:     balances[i].String(),
			Timestamp:   time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
			BlockNumber: blockNumber,
		}
		accountBalances = append(accountBalances, accBalance)
	}
	return accountBalances, nil
}
