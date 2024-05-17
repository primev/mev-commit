package evmclient

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/prometheus/client_golang/prometheus"
)

type TxRequest struct {
	To        *common.Address
	CallData  []byte
	GasPrice  *big.Int
	GasLimit  uint64
	GasFeeCap *big.Int
	Value     *big.Int
}

func (t TxRequest) String() string {
	return fmt.Sprintf(
		"To: %s, CallData: %x, GasPrice: %d, GasLimit: %d, GasFeeCap: %d, Value: %d",
		t.To.String(),
		t.CallData,
		t.GasPrice,
		t.GasLimit,
		t.GasFeeCap,
		t.Value,
	)
}

type Interface interface {
	Send(ctx context.Context, tx *TxRequest) (common.Hash, error)
	WaitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	Call(ctx context.Context, tx *TxRequest) ([]byte, error)
	CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
}

type EvmClient struct {
	mtx       sync.Mutex
	chainID   *big.Int
	ethClient EVM
	owner     common.Address
	keySigner keysigner.KeySigner
	logger    *slog.Logger
	nonce     uint64
	metrics   *metrics
	sentTxs   map[common.Hash]txnDetails
	monitor   *txmonitor
}

type txnDetails struct {
	nonce   uint64
	created time.Time
}

func New(
	keySigner keysigner.KeySigner,
	ethClient EVM,
	logger *slog.Logger,
) (*EvmClient, error) {
	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %w", err)
	}

	m := newMetrics()
	address := keySigner.GetAddress()
	monitor := newTxMonitor(
		address,
		ethClient,
		logger.With("component", "evmclient/txmonitor"),
		m,
	)

	return &EvmClient{
		chainID:   chainID,
		ethClient: ethClient,
		owner:     address,
		keySigner: keySigner,
		logger:    logger,
		metrics:   m,
		sentTxs:   make(map[common.Hash]txnDetails),
		monitor:   monitor,
	}, nil
}

func (c *EvmClient) Close() error {
	return c.monitor.Close()
}

func (c *EvmClient) suggestMaxFeeAndTipCap(
	ctx context.Context,
	gasPrice *big.Int,
) (*big.Int, *big.Int, error) {
	// Returns priority fee per gas
	gasTipCap, err := c.ethClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to suggest gas tip cap: %w", err)
	}

	// Returns priority fee per gas + base fee per gas
	if gasPrice == nil {
		gasPrice, err = c.ethClient.SuggestGasPrice(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to suggest gas price: %w", err)
		}
	}

	return gasPrice, gasTipCap, nil
}

func (c *EvmClient) newTx(ctx context.Context, req *TxRequest, nonce uint64) (*types.Transaction, error) {
	var err error

	if req.GasLimit == 0 {
		// if gas limit is not provided, estimate it
		req.GasLimit, err = c.ethClient.EstimateGas(ctx, ethereum.CallMsg{
			From:  c.owner,
			To:    req.To,
			Data:  req.CallData,
			Value: req.Value,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %w", err)
		}
	}

	gasFeeCap, gasTipCap, err := c.suggestMaxFeeAndTipCap(ctx, req.GasPrice)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest max fee and tip cap: %w", err)
	}

	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		ChainID:   c.chainID,
		To:        req.To,
		Value:     req.Value,
		Gas:       req.GasLimit,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Data:      req.CallData,
	}), nil
}

func (c *EvmClient) getNonce(ctx context.Context) (uint64, error) {
	accountNonce, err := c.ethClient.PendingNonceAt(ctx, c.owner)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %w", err)
	}

	// first nonce
	if accountNonce == 0 {
		return 0, nil
	}

	if c.nonce == 0 {
		c.nonce = accountNonce
	}

	// if external transactions were sent from the owner account, update the
	// nonce to the latest one
	if accountNonce > c.nonce {
		c.nonce = accountNonce
	}

	c.metrics.LastUsedNonce.Set(float64(c.nonce))

	return c.nonce, nil
}

func (c *EvmClient) Send(ctx context.Context, tx *TxRequest) (common.Hash, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.metrics.AttemptedTxCount.Inc()

	nonce, err := c.getNonce(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %w", err)
	}

	if !c.monitor.allowNonce(nonce) {
		return common.Hash{}, fmt.Errorf("too many pending transactions")
	}

	txnData, err := c.newTx(ctx, tx, nonce)
	if err != nil {
		c.logger.Error("failed to create tx", "err", err)
		return common.Hash{}, err
	}

	signedTx, err := c.keySigner.SignTx(txnData, c.chainID)
	if err != nil {
		c.logger.Error("failed to sign tx", "err", err)
		return common.Hash{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		c.logger.Error("failed to send tx", "err", err)
		return common.Hash{}, err
	}

	c.metrics.SentTxCount.Inc()
	c.nonce++
	c.logger.Info("sent txn", "tx", txnString(txnData), "txHash", signedTx.Hash().Hex())

	c.sentTxs[signedTx.Hash()] = txnDetails{nonce: nonce, created: time.Now()}
	c.waitForTxn(signedTx.Hash(), nonce)

	return signedTx.Hash(), nil
}

func (c *EvmClient) waitForTxn(txnHash common.Hash, nonce uint64) {
	go func() {
		res, err := c.monitor.watchTx(txnHash, nonce)
		if err != nil {
			c.logger.Error("failed to watch tx", "err", err)
			return
		}

		receipt := <-res
		if receipt.Err != nil {
			c.logger.Warn("failed to get receipt", "err", receipt.Err)
			if !errors.Is(err, ErrTxnCancelled) {
				return
			}
		} else {
			switch receipt.Receipt.Status {
			case types.ReceiptStatusSuccessful:
				c.metrics.SuccessfulTxCount.Inc()
			case types.ReceiptStatusFailed:
				c.metrics.FailedTxCount.Inc()
			}
		}
		c.mtx.Lock()
		delete(c.sentTxs, txnHash)
		c.mtx.Unlock()

		c.logger.Info("tx status", "txHash", txnHash.Hex(), "status", receipt.Receipt.Status)
	}()

}

func txnString(tx *types.Transaction) string {
	return fmt.Sprintf(
		"nonce=%d, gasPrice=%s, gasLimit=%d, gasTip=%s gasFeeCap=%s to=%s, value=%s, data=%x",
		tx.Nonce(),
		tx.GasPrice().String(),
		tx.Gas(),
		tx.GasTipCap().String(),
		tx.GasFeeCap().String(),
		tx.To().Hex(),
		tx.Value().String(),
		tx.Data(),
	)
}

func (c *EvmClient) WaitForReceipt(
	ctx context.Context,
	txHash common.Hash,
) (*types.Receipt, error) {

	c.mtx.Lock()
	d, ok := c.sentTxs[txHash]
	c.mtx.Unlock()
	if !ok {
		return nil, fmt.Errorf("tx not found")
	}

	res, err := c.monitor.watchTx(txHash, d.nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to watch tx: %w", err)
	}

	select {
	case receipt := <-res:
		if receipt.Err != nil {
			return nil, fmt.Errorf("failed to get receipt: %w", receipt.Err)
		}
		return receipt.Receipt, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *EvmClient) Call(
	ctx context.Context,
	tx *TxRequest,
) ([]byte, error) {

	msg := ethereum.CallMsg{
		From:     c.owner,
		To:       tx.To,
		Data:     tx.CallData,
		Gas:      tx.GasLimit,
		GasPrice: tx.GasPrice,
		Value:    tx.Value,
	}

	result, err := c.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		c.logger.Error("failed to call contract", "err", err)
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	return result, nil
}

func (c *EvmClient) CancelTx(ctx context.Context, txnHash common.Hash) (common.Hash, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	txn, isPending, err := c.ethClient.TransactionByHash(ctx, txnHash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			c.metrics.NotFoundDuringCancelCount.Inc()
		}
		return common.Hash{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	if !isPending {
		c.metrics.NotFoundDuringCancelCount.Inc()
		return common.Hash{}, ethereum.NotFound
	}

	gasFeeCap, gasTipCap, err := c.suggestMaxFeeAndTipCap(ctx, txn.GasPrice())
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to suggest max fee and tip cap: %w", err)
	}

	if gasFeeCap.Cmp(txn.GasFeeCap()) <= 0 {
		gasFeeCap = txn.GasFeeCap()
	}

	if gasTipCap.Cmp(txn.GasTipCap()) <= 0 {
		gasTipCap = txn.GasTipCap()
	}

	// increase gas fee cap and tip cap by 10% for better chance of replacing
	gasTipCap = new(big.Int).Div(new(big.Int).Mul(gasTipCap, big.NewInt(110)), big.NewInt(100))
	gasFeeCap.Add(gasFeeCap, gasTipCap)

	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     txn.Nonce(),
		ChainID:   c.chainID,
		To:        &c.owner,
		Value:     big.NewInt(0),
		Gas:       21000,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Data:      []byte{},
	})

	signedTx, err := c.keySigner.SignTx(tx, c.chainID)
	if err != nil {
		c.logger.Error("failed to sign cancel tx", "err", err)
		return common.Hash{}, fmt.Errorf("failed to sign cancel tx: %w", err)
	}

	err = c.ethClient.SendTransaction(ctx, signedTx)
	if err != nil {
		c.logger.Error("failed to send cancel tx", "err", err)
		return common.Hash{}, err
	}

	c.metrics.CancelledTxCount.Inc()
	c.logger.Info("sent cancel txn", "txHash", signedTx.Hash().Hex())

	c.sentTxs[signedTx.Hash()] = txnDetails{nonce: txn.Nonce(), created: time.Now()}
	c.waitForTxn(signedTx.Hash(), txn.Nonce())

	return signedTx.Hash(), nil
}

func (c *EvmClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, logsCh chan<- types.Log) (ethereum.Subscription, error) {
	return c.ethClient.SubscribeFilterLogs(ctx, query, logsCh)
}

func (c *EvmClient) BlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	return c.ethClient.BlockByNumber(ctx, blockNumber)
}

func (c *EvmClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

func (c *EvmClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return c.ethClient.FilterLogs(ctx, query)
}

type TxnInfo struct {
	Hash    string
	Nonce   uint64
	Created string
}

func (c *EvmClient) PendingTxns() []TxnInfo {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var txns []TxnInfo
	for hash, d := range c.sentTxs {
		txns = append(txns, TxnInfo{
			Hash:    hash.Hex(),
			Nonce:   d.nonce,
			Created: d.created.String(),
		})
	}

	sort.Slice(txns, func(i, j int) bool {
		return txns[i].Nonce < txns[j].Nonce
	})

	return txns
}

func (c *EvmClient) Metrics() []prometheus.Collector {
	return c.metrics.collectors()
}
