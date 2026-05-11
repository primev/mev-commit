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

// orphanRetryWindow is how long we keep retrying a row whose L1 tx hasn't
// shown up in the fastrpc DB yet. ETH-path swaps are user-submitted so the
// fastrpc indexer can lag behind L1 by an unbounded amount; we treat a row as
// a definitive orphan (0 miles) only after it has been missing this long.
const orphanRetryWindow = 24 * time.Hour

// bidIndexerGrace is how long we wait for the mev-commit bid indexer to catch
// up before processing a row with bidCost=0.
const bidIndexerGrace = 15 * time.Minute

// bidCheckOutcome describes what to do with a row whose bid lookup returned 0.
type bidCheckOutcome int

const (
	// bidCheckProceed means fall through and compute miles with whatever bid
	// cost value we have (typically 0 because the indexer is behind).
	bidCheckProceed bidCheckOutcome = iota
	// bidCheckRetry means leave the row pending and reevaluate next cycle.
	bidCheckRetry
	// bidCheckOrphan means mark the row processed with 0 miles and move on —
	// the tx did not go through fastrpc so no bid was ever placed.
	bidCheckOrphan
)

// decideBidCheckOutcome encodes how we handle a row whose bid lookup returned 0.
//
//   - Permit-path rows (userPaysGas=false) are always executor-submitted via
//     fastrpc by construction. A missing fastrpc row can only mean indexer lag,
//     never a non-fastrpc submission, so they follow the bid-indexer grace path
//     regardless of whether fastrpc has caught up yet.
//   - ETH-path rows that ARE in fastrpc follow the same grace path.
//   - ETH-path rows that are NOT in fastrpc retry for orphanRetryWindow before
//     being marked as definitive orphans (the user genuinely bypassed fastrpc).
//
// hasBlockTS / txAge come from the row's block_timestamp column. When the
// timestamp is invalid we fall back to bidCheckRetry in the in-fastrpc case
// (indeterminate age; err on the side of retrying) and to bidCheckOrphan in
// the not-in-fastrpc case (matches prior behavior).
func decideBidCheckOutcome(userPaysGas, inFastRPC, hasBlockTS bool, txAge time.Duration) bidCheckOutcome {
	txInFastRPC := !userPaysGas || inFastRPC
	if txInFastRPC {
		if hasBlockTS && txAge > bidIndexerGrace {
			return bidCheckProceed
		}
		return bidCheckRetry
	}
	// ETH path, not in fastrpc
	if hasBlockTS && txAge < orphanRetryWindow {
		return bidCheckRetry
	}
	return bidCheckOrphan
}

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

	// New (sweep redesign): upfront pricing + cost estimation + sweep
	// scheduling components. All optional in the sense that nil values
	// degrade gracefully — the miles loop simply defers every row to the
	// legacy sweep-time path. Production wiring lives in main.go.
	PriceOracle   *priceOracle
	CostEstimator *costEstimator
	GasBuffer     *gasBuffer
	SweepClock    *sweepClock
}

type ethRow struct {
	txHash     string
	user       string
	surplus    string
	gasCost    sql.NullString
	inputToken string
	blockTS    sql.NullTime
	miles      sql.NullInt64
}

type erc20Row struct {
	txHash     string
	user       string
	token      string // output_token
	surplus    string
	gasCost    sql.NullString
	inputToken string
	inputAmt   string // raw input token units (or ETH wei when input is ETH)
	userAmtOut string // raw output token units delivered to user
	blockTS    sql.NullTime
	miles      sql.NullInt64
}

// -------------------- ETH/WETH Miles --------------------

func processMiles(ctx context.Context, cfg *serviceConfig) (int, error) {
	rows, err := cfg.DB.QueryContext(ctx, `
SELECT tx_hash, user_address, surplus, gas_cost, input_token, block_timestamp, miles
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
		if err := rows.Scan(&r.txHash, &r.user, &r.surplus, &r.gasCost, &r.inputToken, &r.blockTS, &r.miles); err != nil {
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
			txAge := time.Duration(0)
			if r.blockTS.Valid {
				txAge = time.Since(r.blockTS.Time)
			}
			inFastRPC := fastRPCSet[strings.ToLower(r.txHash)]

			switch decideBidCheckOutcome(userPaysGas, inFastRPC, r.blockTS.Valid, txAge) {
			case bidCheckProceed:
				cfg.Logger.Info("tx in FastRPC but bid never indexed, processing with 0 bid cost",
					slog.String("tx", r.txHash), slog.String("user", r.user))
				// fall through with bidCostWei = 0
			case bidCheckRetry:
				cfg.Logger.Info("tx bid lookup pending, will retry next cycle",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Bool("in_fastrpc", inFastRPC),
					slog.Bool("user_pays_gas", userPaysGas),
					slog.Duration("age", txAge))
				continue
			case bidCheckOrphan:
				cfg.Logger.Info("tx not in FastRPC after retry window, skipping with 0 miles",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Duration("age", txAge))
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

		// Idempotency guard: if this row already has miles set (even if
		// `processed` got flipped back to false by a reset SQL), do NOT
		// re-submit to Fuel — just refresh the derived columns. miles is only
		// ever written after a successful Fuel submission (or as 0 for the
		// no-credit terminal paths), so a non-null value means we already
		// settled this row's outcome.
		if r.miles.Valid {
			cfg.Logger.Info("tx already has miles recorded, skipping re-submission",
				slog.String("tx", r.txHash), slog.String("user", r.user),
				slog.Int64("recorded_miles", r.miles.Int64),
				slog.Int64("recomputed_miles", miles.Int64()))
			markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, r.miles.Int64, bidCostWei.String())
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

// erc20CycleStats counts the outcomes of one processERC20Miles invocation.
// Logged at the end of the cycle as `erc20_cycle_summary` so deployment can
// be monitored at a glance: a steady stream of upfront_awarded with low
// deferred_no_chainlink confirms the upfront path is doing its job; spikes
// in deferred_not_whitelisted indicate user activity on novel tokens worth
// investigating.
type erc20CycleStats struct {
	idempotentSkipped       int
	bidRetried              int
	orphaned                int
	upfrontAwarded          int
	upfrontNoProfit         int
	upfrontSubThreshold     int
	upfrontFuelFailed       int
	deferredNotWhitelisted  int
	deferredNoChainlink     int
	deferredInvalidEvent    int
	deferredNoTokenDecimals int
	deferredOtherReason     int
	deferredAwarded         int
	deferredBadSurplus      int
}

func (s *erc20CycleStats) noteDeferralSource(source string) {
	switch source {
	case "deferred:not_whitelisted":
		s.deferredNotWhitelisted++
	case "deferred:no_chainlink":
		s.deferredNoChainlink++
	case "deferred:invalid_event":
		s.deferredInvalidEvent++
	case "deferred:no_token_decim":
		s.deferredNoTokenDecimals++
	default:
		s.deferredOtherReason++
	}
}

func (s *erc20CycleStats) total() int {
	return s.idempotentSkipped + s.bidRetried + s.orphaned +
		s.upfrontAwarded + s.upfrontNoProfit + s.upfrontSubThreshold + s.upfrontFuelFailed +
		s.deferredNotWhitelisted + s.deferredNoChainlink + s.deferredInvalidEvent +
		s.deferredNoTokenDecimals + s.deferredOtherReason + s.deferredAwarded +
		s.deferredBadSurplus
}

// processERC20Miles handles all pending erc20-output rows. Each row is routed
// either to the upfront-awarding path (when the priceOracle returns a value)
// or to the deferred batch path (legacy sweep-then-pro-rata logic, used for
// non-whitelisted output tokens and whitelisted tokens that lack a Chainlink
// feed). The deferred path is unchanged behavior — it still triggers its own
// sweep when the batch becomes profitable. The cadence-based sweep loop
// (sweep_loop.go) operates independently, sweeping accumulated balance from
// upfront-awarded rows once per cadence interval.
func processERC20Miles(ctx context.Context, cfg *serviceConfig) (int, error) {
	processed := 0
	stats := erc20CycleStats{}

	rows, err := cfg.DB.QueryContext(ctx, `
SELECT tx_hash, user_address, output_token, surplus, gas_cost, input_token, input_amount, user_amt_out, block_timestamp, miles
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
		if err := rows.Scan(&r.txHash, &r.user, &r.token, &r.surplus, &r.gasCost, &r.inputToken, &r.inputAmt, &r.userAmtOut, &r.blockTS, &r.miles); err != nil {
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

	allErc20Hashes := make([]string, len(pending))
	for i, r := range pending {
		allErc20Hashes[i] = r.txHash
	}
	erc20BidMap := batchLookupBidCosts(cfg.Logger, cfg.DB, allErc20Hashes)
	erc20FastRPCSet := batchCheckFastRPC(cfg.Logger, cfg.DB, allErc20Hashes)

	// First pass: route every row through bid/orphan checks; rows that survive
	// either get awarded upfront (priceOracle eligible) or land in a per-token
	// deferred bucket for the legacy sweep-time pro-rata path.
	deferredByToken := make(map[string][]deferredErc20Row)
	for _, r := range pending {
		// Idempotency: a row with miles already recorded was settled on a
		// prior run. Just re-flip processed and preserve all derived columns;
		// we can't recompute surplus_eth / net_profit_eth without re-sweeping.
		if r.miles.Valid {
			cfg.Logger.Info("erc20 tx already has miles recorded, skipping",
				slog.String("tx", r.txHash), slog.String("user", r.user),
				slog.Int64("recorded_miles", r.miles.Int64))
			if !cfg.DryRun {
				markProcessedFlagOnly(cfg.DB, r.txHash)
			}
			stats.idempotentSkipped++
			processed++
			continue
		}

		userPaysGas := strings.EqualFold(r.inputToken, zeroAddr.Hex())
		bidCostWei := getBidCost(erc20BidMap, r.txHash)
		if bidCostWei.Sign() == 0 {
			txAge := time.Duration(0)
			if r.blockTS.Valid {
				txAge = time.Since(r.blockTS.Time)
			}
			inFastRPC := erc20FastRPCSet[strings.ToLower(r.txHash)]

			switch decideBidCheckOutcome(userPaysGas, inFastRPC, r.blockTS.Valid, txAge) {
			case bidCheckProceed:
				cfg.Logger.Info("erc20 tx in FastRPC but bid never indexed, processing with 0 bid cost",
					slog.String("tx", r.txHash), slog.String("user", r.user))
				// fall through with bidCostWei = 0
			case bidCheckRetry:
				cfg.Logger.Info("erc20 tx bid lookup pending, will retry next cycle",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Bool("in_fastrpc", inFastRPC),
					slog.Bool("user_pays_gas", userPaysGas),
					slog.Duration("age", txAge))
				stats.bidRetried++
				continue
			case bidCheckOrphan:
				cfg.Logger.Info("erc20 tx not in FastRPC after retry window, skipping with 0 miles",
					slog.String("tx", r.txHash), slog.String("user", r.user),
					slog.Duration("age", txAge))
				if !cfg.DryRun {
					// surplus_eth=0 because raw surplus is in output-token
					// units and weiToEth would yield nonsense for non-ETH
					// outputs (e.g. PEPE 14,413 ETH). The orphan's surplus
					// tokens stay in the executor wallet and become
					// protocol revenue when the cadence sweep eventually
					// fires for this token.
					markProcessed(cfg.DB, r.txHash, 0, 0, 0, "0")
				}
				stats.orphaned++
				processed++
				continue
			}
		}

		gasCostWei := big.NewInt(0)
		// CRITICAL: ETH-input swaps (userPaysGas) MUST NOT subtract gas_cost.
		// The user paid that L1 gas from their own wallet, not from protocol
		// revenue. Mirroring this exactly is what the backtest validated as
		// +18% miles to users; getting it wrong showed -12%.
		if !userPaysGas && r.gasCost.Valid && r.gasCost.String != "" {
			if gc, ok := new(big.Int).SetString(r.gasCost.String, 10); ok {
				gasCostWei = gc
			}
		}

		surplusWei, ok := new(big.Int).SetString(r.surplus, 10)
		if !ok || surplusWei.Sign() <= 0 {
			cfg.Logger.Warn("bad surplus", slog.String("surplus", r.surplus), slog.String("tx", r.txHash))
			stats.deferredBadSurplus++
			continue
		}

		// Try upfront pricing. nil priceOracle (e.g. dry-run misconfiguration)
		// causes everything to defer — same as a non-whitelisted token.
		if cfg.PriceOracle != nil {
			outputAddr := common.HexToAddress(r.token)
			inputAddr := common.HexToAddress(r.inputToken)
			inputAmtBig, _ := new(big.Int).SetString(r.inputAmt, 10)
			userAmtOutBig, _ := new(big.Int).SetString(r.userAmtOut, 10)
			surplusEthWei, eligible, source := cfg.PriceOracle.PriceSurplusEth(
				ctx, inputAddr, outputAddr, inputAmtBig, userAmtOutBig, surplusWei,
			)
			if eligible {
				outcome := awardUpfrontERC20Miles(ctx, cfg, r, surplusEthWei, gasCostWei, bidCostWei, source)
				switch outcome {
				case upfrontAwarded:
					stats.upfrontAwarded++
					processed++
				case upfrontNoProfit:
					stats.upfrontNoProfit++
					processed++
				case upfrontSubThreshold:
					stats.upfrontSubThreshold++
					processed++
				case upfrontFuelFailed:
					stats.upfrontFuelFailed++
				}
				continue
			}
			cfg.Logger.Debug("erc20 row deferred to sweep-time pricing",
				slog.String("tx", r.txHash), slog.String("token", r.token),
				slog.String("reason", source))
			stats.noteDeferralSource(source)
		} else {
			stats.deferredOtherReason++
		}

		deferredByToken[r.token] = append(deferredByToken[r.token], deferredErc20Row{
			row:     r,
			gas:     gasCostWei,
			bid:     bidCostWei,
			surplus: surplusWei,
		})
	}

	// Second pass: legacy sweep-then-pro-rata path for rows the priceOracle
	// could not value upfront. Per-token batch, exactly the prior behavior.
	for token, rows := range deferredByToken {
		n, err := processDeferredERC20Batch(ctx, cfg, token, rows)
		if err != nil {
			cfg.Logger.Warn("deferred erc20 batch failed",
				slog.String("token", token), slog.Any("error", err))
		}
		stats.deferredAwarded += n
		processed += n
	}

	// One-line cycle summary. Grep this in Groundcover to see at a glance
	// whether the upfront path is working and what's flowing through each
	// branch. Quiet cycles (nothing pending) emit nothing.
	if stats.total() > 0 {
		cfg.Logger.Info("erc20_cycle_summary",
			slog.Int("total", stats.total()),
			slog.Int("upfront_awarded", stats.upfrontAwarded),
			slog.Int("upfront_no_profit", stats.upfrontNoProfit),
			slog.Int("upfront_sub_threshold", stats.upfrontSubThreshold),
			slog.Int("upfront_fuel_failed", stats.upfrontFuelFailed),
			slog.Int("deferred_awarded", stats.deferredAwarded),
			slog.Int("deferred_not_whitelisted", stats.deferredNotWhitelisted),
			slog.Int("deferred_no_chainlink", stats.deferredNoChainlink),
			slog.Int("deferred_no_token_decim", stats.deferredNoTokenDecimals),
			slog.Int("deferred_invalid_event", stats.deferredInvalidEvent),
			slog.Int("deferred_other", stats.deferredOtherReason),
			slog.Int("deferred_bad_surplus", stats.deferredBadSurplus),
			slog.Int("idempotent_skipped", stats.idempotentSkipped),
			slog.Int("bid_retried", stats.bidRetried),
			slog.Int("orphaned", stats.orphaned))
	}

	return processed, nil
}

// deferredErc20Row carries the parsed wei values alongside the raw row, so
// the deferred batch processor doesn't have to re-parse strings.
type deferredErc20Row struct {
	row     erc20Row
	gas     *big.Int // wei; zero for ETH-input (userPaysGas)
	bid     *big.Int // wei
	surplus *big.Int // raw output-token units
}

// upfrontOutcome describes how awardUpfrontERC20Miles handled a row, used
// for cycle-summary stats.
type upfrontOutcome int

const (
	// upfrontAwarded — positive miles posted to Fuel and row marked processed.
	upfrontAwarded upfrontOutcome = iota
	// upfrontNoProfit — net was non-positive; row settled with miles=0.
	upfrontNoProfit
	// upfrontSubThreshold — net positive but below the per-mile floor;
	// row settled with miles=0.
	upfrontSubThreshold
	// upfrontFuelFailed — Fuel submission errored; row stays pending for
	// retry next cycle.
	upfrontFuelFailed
)

// awardUpfrontERC20Miles applies the new design's per-row formula and posts
// miles immediately.
//
//	net = surplus_eth - deductible_user_gas - user_bid - estimated_overhead
//	miles = max(0, net × 0.9 / weiPerPoint)
//
// estimated_overhead comes from the per-token p25 of realized sweep overhead
// (costEstimator). The protocol absorbs variance between estimate and
// realized; reconciliationMonitor catches drift over weeks.
func awardUpfrontERC20Miles(
	ctx context.Context,
	cfg *serviceConfig,
	r erc20Row,
	surplusEthWei, gasCostWei, bidCostWei *big.Int,
	pricingSource string,
) upfrontOutcome {
	overheadFloat := costEstimateLastResort
	overheadSource := "default_no_data"
	if cfg.CostEstimator != nil {
		est := cfg.CostEstimator.Get(r.token)
		overheadFloat = est.PerRowOverhead
		overheadSource = est.Source
	}
	overheadWei := ethFloatToWei(overheadFloat)

	netProfit := new(big.Int).Sub(surplusEthWei, gasCostWei)
	netProfit.Sub(netProfit, bidCostWei)
	netProfit.Sub(netProfit, overheadWei)

	surplusEth := weiToEth(surplusEthWei)
	netProfitEth := weiToEth(netProfit)

	if netProfit.Sign() <= 0 {
		cfg.Logger.Info("no upfront profit",
			slog.String("tx", r.txHash), slog.String("user", r.user),
			slog.String("token", r.token),
			slog.String("pricing", pricingSource),
			slog.String("overhead_src", overheadSource),
			slog.Float64("surplus_eth", surplusEth),
			slog.Float64("overhead_eth", overheadFloat),
			slog.Float64("net_profit_eth", netProfitEth))
		if !cfg.DryRun {
			markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
		}
		return upfrontNoProfit
	}

	userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
	userShare.Div(userShare, big.NewInt(100))

	miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
	if miles.Sign() <= 0 {
		cfg.Logger.Info("sub-threshold upfront",
			slog.String("tx", r.txHash), slog.String("user", r.user),
			slog.String("token", r.token),
			slog.String("pricing", pricingSource),
			slog.Float64("surplus_eth", surplusEth),
			slog.Float64("net_profit_eth", netProfitEth))
		if !cfg.DryRun {
			markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
		}
		return upfrontSubThreshold
	}

	cfg.Logger.Info("awarding upfront erc20 miles",
		slog.Int64("miles", miles.Int64()),
		slog.String("user", r.user), slog.String("tx", r.txHash),
		slog.String("token", r.token),
		slog.String("pricing", pricingSource),
		slog.String("overhead_src", overheadSource),
		slog.Float64("surplus_eth", surplusEth),
		slog.Float64("overhead_eth", overheadFloat),
		slog.Float64("net_profit_eth", netProfitEth))

	if cfg.DryRun {
		return upfrontAwarded
	}

	if err := submitToFuel(ctx, cfg.HTTPClient, cfg.FuelURL, cfg.FuelKey,
		common.HexToAddress(r.user),
		common.HexToHash(r.txHash),
		miles,
	); err != nil {
		cfg.Logger.Error("fuel submit failed, will retry next cycle",
			slog.String("tx", r.txHash), slog.Any("error", err))
		return upfrontFuelFailed
	}
	markProcessed(cfg.DB, r.txHash, surplusEth, netProfitEth, miles.Int64(), bidCostWei.String())
	return upfrontAwarded
}

// processDeferredERC20Batch is the legacy sweep-then-pro-rata path for rows
// the priceOracle couldn't value upfront. Behavior matches the prior
// processERC20Miles batch logic exactly — per-token Barter quote,
// profitability gate that includes user gas and bid cost, sweep submission,
// pro-rata miles per row.
func processDeferredERC20Batch(ctx context.Context, cfg *serviceConfig, token string, rows []deferredErc20Row) (int, error) {
	if len(rows) == 0 {
		return 0, nil
	}

	processed := 0
	totalOriginalGasCost := big.NewInt(0)
	totalOriginalBidCost := big.NewInt(0)
	readyTotalSum := big.NewInt(0)
	for _, d := range rows {
		totalOriginalGasCost.Add(totalOriginalGasCost, d.gas)
		totalOriginalBidCost.Add(totalOriginalBidCost, d.bid)
		readyTotalSum.Add(readyTotalSum, d.surplus)
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
		// "Amount too low" is Barter's response when the input is below its
		// per-token minimum and is the normal steady state for early/small
		// batches — log at debug so it doesn't flood the warn stream. Other
		// errors (timeouts, 5xx, malformed responses) still bubble up to the
		// caller's warn log.
		if isBarterAmountTooLow(err) {
			cfg.Logger.Debug("deferred batch under barter minimum",
				slog.String("token", token),
				slog.Int("batch_size", len(rows)),
				slog.String("amount", readyTotalSum.String()))
			return processed, nil
		}
		return processed, fmt.Errorf("callBarter: %w", err)
	}

	gasLimit, err := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	if err != nil {
		return processed, fmt.Errorf("invalid gasLimit %q: %w", barterResp.GasLimit, err)
	}
	gasLimit += 50000

	gasPrice, err := cfg.Client.SuggestGasPrice(ctx)
	if err != nil {
		return processed, fmt.Errorf("suggest gas price: %w", err)
	}

	expectedGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	expectedEthReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
	if !ok {
		return processed, fmt.Errorf("invalid MinReturn from barter: %s", barterResp.MinReturn)
	}

	totalSweepCosts := new(big.Int).Add(expectedGasCost, totalOriginalBidCost)
	totalSweepCosts.Add(totalSweepCosts, totalOriginalGasCost)

	if expectedEthReturn.Cmp(totalSweepCosts) <= 0 {
		// Debug-only: counts in erc20_cycle_summary's deferred_other_reason
		// are the info-level signal that deferred rows are accumulating
		// without being profitable.
		cfg.Logger.Debug("deferred token sweep not yet profitable",
			slog.String("token", token),
			slog.Int("batch_size", len(rows)),
			slog.Float64("return_eth", weiToEth(expectedEthReturn)),
			slog.Float64("total_cost_eth", weiToEth(totalSweepCosts)))
		return processed, nil
	}

	var actualEthReturn, actualSwapGasCost *big.Int
	if cfg.DryRun {
		cfg.Logger.Info("simulated deferred sweep",
			slog.String("amount", readyTotalSum.String()),
			slog.String("token", token),
			slog.Float64("return_eth", weiToEth(expectedEthReturn)),
			slog.Float64("gas_eth", weiToEth(expectedGasCost)))
		actualEthReturn = expectedEthReturn
		actualSwapGasCost = expectedGasCost
	} else {
		actualEthReturn, actualSwapGasCost, err = submitFastSwapSweep(ctx, cfg.Logger, cfg.Client, cfg.L1Client, cfg.HTTPClient, cfg.Signer, cfg.ExecutorAddr, common.HexToAddress(token), readyTotalSum, cfg.FastSwapURL, cfg.FundsRecipient, cfg.SettlementAddr, barterResp, cfg.MaxGasGwei)
		if err != nil {
			return processed, fmt.Errorf("submit sweep: %w", err)
		}
		cfg.Logger.Info("deferred sweep success",
			slog.String("token", token),
			slog.Int("batch_size", len(rows)),
			slog.Float64("return_eth", weiToEth(actualEthReturn)),
			slog.Float64("gas_eth", weiToEth(actualSwapGasCost)))
		// Cadence sweep loop should not double-sweep this token immediately;
		// mark its clock now so it respects the cadence floor.
		if cfg.SweepClock != nil {
			cfg.SweepClock.MarkSwept(common.HexToAddress(token), time.Now())
		}
	}

	for _, d := range rows {
		txGrossEth := new(big.Int).Mul(actualEthReturn, d.surplus)
		txGrossEth.Div(txGrossEth, readyTotalSum)

		txOverheadGas := new(big.Int).Mul(actualSwapGasCost, d.surplus)
		txOverheadGas.Div(txOverheadGas, readyTotalSum)

		txNetProfit := new(big.Int).Sub(txGrossEth, d.gas)
		txNetProfit.Sub(txNetProfit, d.bid)
		txNetProfit.Sub(txNetProfit, txOverheadGas)

		surplusEth := weiToEth(txGrossEth)
		netProfitEth := weiToEth(txNetProfit)

		if txNetProfit.Sign() <= 0 {
			cfg.Logger.Info("no profit for deferred subset tx",
				slog.String("tx", d.row.txHash), slog.String("user", d.row.user),
				slog.Float64("gross_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth))
			if !cfg.DryRun {
				markProcessed(cfg.DB, d.row.txHash, surplusEth, netProfitEth, 0, d.bid.String())
			}
			processed++
			continue
		}

		userShare := new(big.Int).Mul(txNetProfit, big.NewInt(90))
		userShare.Div(userShare, big.NewInt(100))
		miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
		if miles.Sign() <= 0 {
			cfg.Logger.Info("sub-threshold deferred subset tx",
				slog.String("tx", d.row.txHash), slog.String("user", d.row.user),
				slog.Float64("gross_eth", surplusEth), slog.Float64("net_profit_eth", netProfitEth))
			if !cfg.DryRun {
				markProcessed(cfg.DB, d.row.txHash, surplusEth, netProfitEth, 0, d.bid.String())
			}
			processed++
			continue
		}

		cfg.Logger.Info("awarding deferred miles for subset tx",
			slog.Int64("miles", miles.Int64()), slog.String("user", d.row.user),
			slog.String("tx", d.row.txHash), slog.Float64("gross_eth", surplusEth),
			slog.Float64("net_profit_eth", netProfitEth))

		if cfg.DryRun {
			processed++
			continue
		}

		if err := submitToFuel(ctx, cfg.HTTPClient, cfg.FuelURL, cfg.FuelKey,
			common.HexToAddress(d.row.user),
			common.HexToHash(d.row.txHash),
			miles,
		); err != nil {
			cfg.Logger.Error("fuel submit failed, will retry next cycle",
				slog.String("tx", d.row.txHash), slog.Any("error", err))
			continue
		}
		markProcessed(cfg.DB, d.row.txHash, surplusEth, netProfitEth, miles.Int64(), d.bid.String())
		processed++
	}

	return processed, nil
}

// ethFloatToWei converts a float ETH amount to wei (big.Int). Used to bring
// the cost-estimator's float overhead into the wei domain where the rest of
// the miles arithmetic lives.
func ethFloatToWei(eth float64) *big.Int {
	wei, _ := new(big.Float).Mul(big.NewFloat(eth), big.NewFloat(1e18)).Int(nil)
	if wei == nil {
		return big.NewInt(0)
	}
	return wei
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
