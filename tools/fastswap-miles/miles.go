package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/keysigner"
)

// serviceConfig holds references shared across the miles processing pipeline.
type serviceConfig struct {
	Logger         *slog.Logger
	DB             *sql.DB
	WETH           common.Address
	FuelURL        string
	FuelKey        string
	BarterURL      string
	BarterKey      string
	Client         *ethclient.Client
	L1Client       *ethclient.Client
	Signer         keysigner.KeySigner
	ExecutorAddr   common.Address
	HTTPClient     *http.Client
	DryRun         bool
	FastSwapURL    string
	FundsRecipient common.Address
	SettlementAddr common.Address
	MaxGasGwei     uint64
}

type ethRow struct {
	txHash     string
	user       string
	surplus    string
	gasCost    sql.NullString
	inputToken string
	blockTS    sql.NullTime
}

type erc20Row struct {
	txHash     string
	user       string
	token      string
	surplus    string
	gasCost    sql.NullString
	inputToken string
	blockTS    sql.NullTime
}

type tokenBatch struct {
	Token    string
	TotalSum *big.Int
	Txs      []erc20Row
}

// -------------------- ETH/WETH Miles --------------------

func processMiles(ctx context.Context, cfg *serviceConfig) (int, error) {
	rows, err := cfg.DB.QueryContext(ctx, `
SELECT tx_hash, user_address, surplus, gas_cost, input_token, block_timestamp
FROM mevcommit_57173.fastswap_miles
WHERE processed = false
  AND swap_type = 'eth_weth'
  AND LOWER(user_address) != LOWER(?)
`, cfg.ExecutorAddr.Hex())
	if err != nil {
		return 0, fmt.Errorf("query unprocessed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var pending []ethRow
	for rows.Next() {
		var r ethRow
		if err := rows.Scan(&r.txHash, &r.user, &r.surplus, &r.gasCost, &r.inputToken, &r.blockTS); err != nil {
			return 0, err
		}
		pending = append(pending, r)
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	allHashes := make([]string, len(pending))
	for i, r := range pending {
		allHashes[i] = r.txHash
	}
	bidMap := batchLookupBidCosts(cfg.Logger, cfg.DB, allHashes)
	fastRPCSet := batchCheckFastRPC(cfg.Logger, cfg.DB, allHashes)

	processed := 0
	for _, r := range pending {
		surplusWei, ok := new(big.Int).SetString(r.surplus, 10)
		if !ok {
			cfg.Logger.Warn("bad surplus", slog.String("surplus", r.surplus), slog.String("tx", r.txHash))
			continue
		}

		userPaysGas := strings.EqualFold(r.inputToken, zeroAddr.Hex())

		gasCostWei := big.NewInt(0)
		if !userPaysGas && r.gasCost.Valid && r.gasCost.String != "" {
			if gc, ok := new(big.Int).SetString(r.gasCost.String, 10); ok {
				gasCostWei = gc
			}
		}

		bidCostWei := getBidCost(bidMap, r.txHash)

		if bidCostWei.Sign() == 0 {
			if fastRPCSet[strings.ToLower(r.txHash)] {
				if r.blockTS.Valid && time.Since(r.blockTS.Time) > 15*time.Minute {
					cfg.Logger.Info("tx in FastRPC but bid never indexed, processing with 0 bid cost",
						slog.String("tx", r.txHash), slog.String("user", r.user))
					// fall through to normal miles calculation with bidCostWei = 0
				} else {
					cfg.Logger.Info("tx in FastRPC but bid not indexed yet, will retry",
						slog.String("tx", r.txHash), slog.String("user", r.user))
					continue
				}
			} else {
				cfg.Logger.Info("tx not in FastRPC, skipping with 0 miles",
					slog.String("tx", r.txHash), slog.String("user", r.user))
				if !cfg.DryRun {
					markProcessed(cfg.DB, r.txHash, weiToEth(surplusWei), 0, 0, "0")
				}
				processed++
				continue
			}
		}

		netProfit := new(big.Int).Sub(surplusWei, gasCostWei)
		netProfit.Sub(netProfit, bidCostWei)

		surplusEth := weiToEth(surplusWei)
		netProfitEth := weiToEth(netProfit)

		if netProfit.Sign() <= 0 {
			cfg.Logger.Info("no profit",
				slog.String("tx", r.txHash), slog.String("user", r.user),
				slog.Float64("surplus_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth),
				slog.String("gas", gasCostWei.String()), slog.String("bid", bidCostWei.String()))
			if !cfg.DryRun {
				markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
			}
			processed++
			continue
		}

		userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
		userShare.Div(userShare, big.NewInt(100))

		miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
		if miles.Sign() <= 0 {
			cfg.Logger.Info("sub-threshold",
				slog.String("tx", r.txHash), slog.String("user", r.user),
				slog.Float64("surplus_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth))
			if !cfg.DryRun {
				markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
			}
			processed++
			continue
		}

		cfg.Logger.Info("awarding miles",
			slog.Int64("miles", miles.Int64()), slog.String("user", r.user),
			slog.String("tx", r.txHash), slog.Float64("surplus_eth", surplusEth),
			slog.Float64("net_profit_eth", netProfitEth),
			slog.String("gas", gasCostWei.String()), slog.String("bid", bidCostWei.String()))

		if cfg.DryRun {
			processed++
			continue
		}

		err := submitToFuel(ctx, cfg.HTTPClient, cfg.FuelURL, cfg.FuelKey,
			common.HexToAddress(r.user),
			common.HexToHash(r.txHash),
			miles,
		)
		if err != nil {
			cfg.Logger.Error("fuel submit failed", slog.String("tx", r.txHash), slog.Any("error", err))
			continue
		}

		markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, miles.Int64(), bidCostWei.String())
		processed++
		cfg.Logger.Info("awarded miles",
			slog.Int64("miles", miles.Int64()), slog.String("user", r.user), slog.String("tx", r.txHash))
	}

	return processed, nil
}

// -------------------- ERC20 Miles --------------------

func processERC20Miles(ctx context.Context, cfg *serviceConfig) (int, error) {
	processed := 0

	rows, err := cfg.DB.QueryContext(ctx, `
SELECT tx_hash, user_address, output_token, surplus, gas_cost, input_token, block_timestamp
FROM mevcommit_57173.fastswap_miles
WHERE processed = false
  AND swap_type = 'erc20'
  AND LOWER(user_address) != LOWER(?)
`, cfg.ExecutorAddr.Hex())
	if err != nil {
		return processed, fmt.Errorf("query erc20 unprocessed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var pending []erc20Row
	for rows.Next() {
		var r erc20Row
		if err := rows.Scan(&r.txHash, &r.user, &r.token, &r.surplus, &r.gasCost, &r.inputToken, &r.blockTS); err != nil {
			return processed, err
		}
		pending = append(pending, r)
	}
	if rows.Err() != nil {
		return processed, rows.Err()
	}

	if len(pending) == 0 {
		return processed, nil
	}

	batches := make(map[string]*tokenBatch)
	for _, r := range pending {
		surplusWei, ok := new(big.Int).SetString(r.surplus, 10)
		if !ok || surplusWei.Sign() <= 0 {
			cfg.Logger.Warn("bad surplus", slog.String("surplus", r.surplus), slog.String("tx", r.txHash))
			continue
		}
		if _, exists := batches[r.token]; !exists {
			batches[r.token] = &tokenBatch{
				Token:    r.token,
				TotalSum: big.NewInt(0),
				Txs:      make([]erc20Row, 0),
			}
		}
		batches[r.token].TotalSum.Add(batches[r.token].TotalSum, surplusWei)
		batches[r.token].Txs = append(batches[r.token].Txs, r)
	}

	allErc20Hashes := make([]string, len(pending))
	for i, r := range pending {
		allErc20Hashes[i] = r.txHash
	}
	erc20BidMap := batchLookupBidCosts(cfg.Logger, cfg.DB, allErc20Hashes)
	erc20FastRPCSet := batchCheckFastRPC(cfg.Logger, cfg.DB, allErc20Hashes)

	for token, batch := range batches {
		totalOriginalGasCost := big.NewInt(0)
		totalOriginalBidCost := big.NewInt(0)

		// First pass: separate rows into ready, pending-bid, and not-in-fastrpc.
		// Pending-bid rows are excluded from the batch so they retry next cycle.
		var readyTxs []erc20Row
		var readyGasCosts []*big.Int
		var readyBidCosts []*big.Int

		for _, r := range batch.Txs {
			bidCostWei := getBidCost(erc20BidMap, r.txHash)
			if bidCostWei.Sign() == 0 {
				if erc20FastRPCSet[strings.ToLower(r.txHash)] {
					if r.blockTS.Valid && time.Since(r.blockTS.Time) > 15*time.Minute {
						cfg.Logger.Info("erc20 tx in FastRPC but bid never indexed, processing with 0 bid cost",
							slog.String("tx", r.txHash), slog.String("user", r.user))
						// fall through with bidCostWei = 0
					} else {
						cfg.Logger.Info("erc20 tx in FastRPC but bid not indexed yet, will retry",
							slog.String("tx", r.txHash), slog.String("user", r.user))
						continue
					}
				} else {
					cfg.Logger.Info("erc20 tx not in FastRPC, skipping with 0 miles",
						slog.String("tx", r.txHash), slog.String("user", r.user))
					surplusWei, _ := new(big.Int).SetString(r.surplus, 10)
					if !cfg.DryRun {
						markProcessed(cfg.DB, r.txHash, weiToEth(surplusWei), 0, 0, "0")
					}
					processed++
					continue
				}
			}

			userPaysGas := strings.EqualFold(r.inputToken, zeroAddr.Hex())

			gasCostWei := big.NewInt(0)
			if !userPaysGas && r.gasCost.Valid && r.gasCost.String != "" {
				if gc, ok := new(big.Int).SetString(r.gasCost.String, 10); ok {
					gasCostWei = gc
				}
			}

			readyTxs = append(readyTxs, r)
			readyGasCosts = append(readyGasCosts, gasCostWei)
			readyBidCosts = append(readyBidCosts, bidCostWei)

			totalOriginalGasCost.Add(totalOriginalGasCost, gasCostWei)
			totalOriginalBidCost.Add(totalOriginalBidCost, bidCostWei)
		}

		if len(readyTxs) == 0 {
			continue
		}

		// Recalculate TotalSum for only the ready rows
		readyTotalSum := big.NewInt(0)
		for _, r := range readyTxs {
			surplusWei, _ := new(big.Int).SetString(r.surplus, 10)
			readyTotalSum.Add(readyTotalSum, surplusWei)
		}

		reqBody := barterRequest{
			Source:            token,
			Target:            cfg.WETH.Hex(),
			SellAmount:        readyTotalSum.String(),
			Recipient:         cfg.ExecutorAddr.Hex(),
			Origin:            cfg.ExecutorAddr.Hex(),
			MinReturnFraction: 0.98,
			Deadline:          fmt.Sprintf("%d", time.Now().Add(10*time.Minute).Unix()),
		}

		barterResp, err := callBarter(ctx, cfg.HTTPClient, cfg.BarterURL, cfg.BarterKey, reqBody)
		if err != nil {
			cfg.Logger.Warn("callBarter failed", slog.String("token", token), slog.Any("error", err))
			continue
		}

		gasLimit, err := strconv.ParseUint(barterResp.GasLimit, 10, 64)
		if err != nil {
			cfg.Logger.Warn("invalid gasLimit from barter", slog.String("gasLimit", barterResp.GasLimit))
			continue
		}
		gasLimit += 50000

		gasPrice, err := cfg.Client.SuggestGasPrice(ctx)
		if err != nil {
			cfg.Logger.Warn("suggest gas price failed", slog.Any("error", err))
			continue
		}

		expectedGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
		expectedEthReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
		if !ok {
			cfg.Logger.Warn("invalid MinReturn from barter", slog.String("minReturn", barterResp.MinReturn))
			continue
		}

		totalSweepCosts := new(big.Int).Add(expectedGasCost, totalOriginalBidCost)
		totalSweepCosts.Add(totalSweepCosts, totalOriginalGasCost)

		if expectedEthReturn.Cmp(totalSweepCosts) <= 0 {
			cfg.Logger.Info("token sweep not yet profitable",
				slog.String("token", token),
				slog.Float64("return_eth", weiToEth(expectedEthReturn)),
				slog.Float64("total_cost_eth", weiToEth(totalSweepCosts)))
			continue
		}

		var actualEthReturn *big.Int
		var actualSwapGasCost *big.Int

		if cfg.DryRun {
			cfg.Logger.Info("simulated sweep",
				slog.String("amount", readyTotalSum.String()),
				slog.String("token", token),
				slog.Float64("return_eth", weiToEth(expectedEthReturn)),
				slog.Float64("gas_eth", weiToEth(expectedGasCost)))
			actualEthReturn = expectedEthReturn
			actualSwapGasCost = expectedGasCost
		} else {
			actualEthReturn, actualSwapGasCost, err = submitFastSwapSweep(ctx, cfg.Logger, cfg.Client, cfg.L1Client, cfg.HTTPClient, cfg.Signer, cfg.ExecutorAddr, common.HexToAddress(token), readyTotalSum, cfg.FastSwapURL, cfg.FundsRecipient, cfg.SettlementAddr, barterResp, cfg.MaxGasGwei)
			if err != nil {
				cfg.Logger.Error("failed to sweep token", slog.String("token", token), slog.Any("error", err))
				continue
			}
			cfg.Logger.Info("FastSwap sweep success",
				slog.String("token", token),
				slog.Float64("return_eth", weiToEth(actualEthReturn)),
				slog.Float64("gas_eth", weiToEth(actualSwapGasCost)))
		}

		for i, r := range readyTxs {
			surplusWei, _ := new(big.Int).SetString(r.surplus, 10)

			txGrossEth := new(big.Int).Mul(actualEthReturn, surplusWei)
			txGrossEth.Div(txGrossEth, readyTotalSum)

			txOverheadGas := new(big.Int).Mul(actualSwapGasCost, surplusWei)
			txOverheadGas.Div(txOverheadGas, readyTotalSum)

			txNetProfit := new(big.Int).Sub(txGrossEth, readyGasCosts[i])
			txNetProfit.Sub(txNetProfit, readyBidCosts[i])
			txNetProfit.Sub(txNetProfit, txOverheadGas)

			surplusEth := weiToEth(txGrossEth)
			netProfitEth := weiToEth(txNetProfit)

			if txNetProfit.Sign() <= 0 {
				cfg.Logger.Info("no profit for subset tx",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Float64("gross_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth))
				if !cfg.DryRun {
					markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, readyBidCosts[i].String())
				}
				processed++
				continue
			}

			userShare := new(big.Int).Mul(txNetProfit, big.NewInt(90))
			userShare.Div(userShare, big.NewInt(100))

			miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
			if miles.Sign() <= 0 {
				cfg.Logger.Info("sub-threshold subset tx",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Float64("gross_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth))
				if !cfg.DryRun {
					markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, readyBidCosts[i].String())
				}
				processed++
				continue
			}

			cfg.Logger.Info("awarding miles for subset tx",
				slog.Int64("miles", miles.Int64()), slog.String("user", r.user),
				slog.String("tx", r.txHash), slog.Float64("gross_eth", surplusEth),
				slog.Float64("net_profit_eth", netProfitEth))

			if cfg.DryRun {
				processed++
				continue
			}

			err := submitToFuel(ctx, cfg.HTTPClient, cfg.FuelURL, cfg.FuelKey,
				common.HexToAddress(r.user),
				common.HexToHash(r.txHash),
				miles,
			)
			if err != nil {
				cfg.Logger.Error("fuel submit failed, will retry next cycle",
					slog.String("tx", r.txHash), slog.Any("error", err))
				continue // don't mark processed — retry next cycle
			}
			markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, miles.Int64(), readyBidCosts[i].String())
			processed++
		}
	}

	return processed, nil
}

// -------------------- Bid Cost / FastRPC Lookups --------------------

func batchLookupBidCosts(logger *slog.Logger, db *sql.DB, txHashes []string) map[string]*big.Int {
	result := make(map[string]*big.Int, len(txHashes))
	if len(txHashes) == 0 {
		return result
	}

	normalized := make([]string, len(txHashes))
	for i, h := range txHashes {
		normalized[i] = strings.TrimPrefix(strings.ToLower(h), "0x")
	}

	var inClause strings.Builder
	for i, h := range normalized {
		if i > 0 {
			inClause.WriteString(", ")
		}
		inClause.WriteString("'")
		inClause.WriteString(h)
		inClause.WriteString("'")
	}

	query := fmt.Sprintf(`
SELECT
  LOWER(
    CASE
      WHEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')), 1, 2) = '0x'
        THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')), 3)
      ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash'))
    END
  ) as txn_hash,
  get_json_string(CAST(l_decoded AS VARCHAR), '$.args.bidAmt') as bid_amt
FROM mevcommit_57173.tx_view
WHERE l_decoded IS NOT NULL
  AND COALESCE(l_removed, 0) = 0
  AND get_json_string(CAST(l_decoded AS VARCHAR), '$.name') = 'OpenedCommitmentStored'
  AND t_chain_id IN (8855, 57173)
  AND LOWER(
    CASE
      WHEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')), 1, 2) = '0x'
        THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')), 3)
      ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash'))
    END
  ) IN (%s)
ORDER BY l_block_number DESC`, inClause.String())

	dbRows, err := db.Query(query)
	if err != nil {
		logger.Error("batchLookupBidCosts query error", slog.Any("error", err))
		return result
	}
	defer func() { _ = dbRows.Close() }()

	for dbRows.Next() {
		var txnHash string
		var bidAmtStr sql.NullString
		if err := dbRows.Scan(&txnHash, &bidAmtStr); err != nil {
			logger.Error("batchLookupBidCosts scan error", slog.Any("error", err))
			continue
		}
		if !bidAmtStr.Valid || bidAmtStr.String == "" || bidAmtStr.String == "null" {
			continue
		}
		// Skip if we already have a value for this tx — since results are
		// ordered by block number DESC, the first row is the most recent
		// OpenedCommitmentStored event (the one that was actually included).
		if _, exists := result[txnHash]; exists {
			continue
		}
		cleanStr := strings.Trim(bidAmtStr.String, "\"")
		v, ok := new(big.Int).SetString(cleanStr, 10)
		if !ok {
			v, ok = new(big.Int).SetString(strings.TrimPrefix(cleanStr, "0x"), 16)
			if !ok {
				logger.Error("batchLookupBidCosts parse error", slog.String("tx", txnHash), slog.String("value", bidAmtStr.String))
				continue
			}
		}
		result[txnHash] = v
	}

	logger.Debug("batchLookupBidCosts", slog.Int("found", len(result)), slog.Int("total", len(txHashes)))
	return result
}

func getBidCost(bidMap map[string]*big.Int, txHash string) *big.Int {
	hashNorm := strings.TrimPrefix(strings.ToLower(txHash), "0x")
	if v, ok := bidMap[hashNorm]; ok {
		return v
	}
	return big.NewInt(0)
}

func batchCheckFastRPC(logger *slog.Logger, db *sql.DB, txHashes []string) map[string]bool {
	result := make(map[string]bool, len(txHashes))
	if len(txHashes) == 0 {
		return result
	}

	var inClause strings.Builder
	for i, h := range txHashes {
		if i > 0 {
			inClause.WriteString(", ")
		}
		inClause.WriteString("'")
		inClause.WriteString(strings.ToLower(h))
		inClause.WriteString("'")
	}

	query := fmt.Sprintf(`
SELECT hash FROM pg_mev_commit_fastrpc.public.mctransactions_sr
WHERE LOWER(hash) IN (%s)`, inClause.String())

	dbRows, err := db.Query(query)
	if err != nil {
		logger.Error("batchCheckFastRPC query error", slog.Any("error", err))
		return result
	}
	defer func() { _ = dbRows.Close() }()

	for dbRows.Next() {
		var h string
		if err := dbRows.Scan(&h); err != nil {
			logger.Error("batchCheckFastRPC scan error", slog.Any("error", err))
			continue
		}
		result[strings.ToLower(h)] = true
	}

	logger.Debug("batchCheckFastRPC", slog.Int("found", len(result)), slog.Int("total", len(txHashes)))
	return result
}
