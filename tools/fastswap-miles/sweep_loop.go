package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// sweepProfitabilityNumerator is the multiplier in the profitability guard
// applied to the sweep gas cost. expected_return > sweep_gas × 1.05 is the
// minimum margin to fire a sweep — 5% buffer for gas-price slippage between
// quote and execution.
const sweepProfitabilityNumerator = 105

// sweepDustRawUnits is a per-token raw-units floor below which we don't bother
// fetching a Barter quote. Most tokens have 6+ decimals so 1000 raw units is
// effectively zero (1e-3 units of a 6-decimal token, etc).
const sweepDustRawUnits = 1000

// sweepBarterMinReturnFraction is the lower bound on Barter's MinReturn
// relative to its quoted output, used in the cadence-sweep quote request.
// Mirrors the deferred path's value (0.98).
const sweepBarterMinReturnFraction = 0.98

const erc20BalanceOfABI = `[
	{"constant":true,"inputs":[{"name":"_owner","type":"address"}],
	 "name":"balanceOf",
	 "outputs":[{"name":"balance","type":"uint256"}],
	 "type":"function"}
]`

// sweepLoop runs cadence-based per-token sweeps for tokens whose surplus is
// already credited to users via upfront miles awarding. The miles pipeline
// has marked those rows processed, so this loop's sole job is to convert
// accumulated surplus tokens to ETH on the configured cadence — fully
// decoupled from any miles processing for them.
//
// Tokens whose rows are deferred (non-whitelisted, no Chainlink feed) are
// swept by processDeferredERC20Batch instead; that path triggers its own
// sweep when its batch becomes profitable. Both paths share a profitability
// guard and the same submitFastSwapSweep underneath, so they don't conflict
// — the wallet balance is the source of truth for what's available.
type sweepLoop struct {
	cfg     *serviceConfig
	balance abi.ABI
}

func newSweepLoop(cfg *serviceConfig) (*sweepLoop, error) {
	a, err := abi.JSON(strings.NewReader(erc20BalanceOfABI))
	if err != nil {
		return nil, fmt.Errorf("parse balanceOf ABI: %w", err)
	}
	return &sweepLoop{cfg: cfg, balance: a}, nil
}

// RunOnce iterates every whitelisted token in tokenConfigs once and processes
// any that are due for a sweep. Cheap to call repeatedly; the cadence check
// is in-memory and short-circuits before any RPC traffic for tokens that
// haven't reached their cadence floor.
func (s *sweepLoop) RunOnce(ctx context.Context) {
	if s.cfg.SweepClock == nil || s.cfg.GasBuffer == nil {
		return
	}
	for token, tokCfg := range tokenConfigs {
		if ctx.Err() != nil {
			return
		}
		s.processToken(ctx, token, tokCfg)
	}
}

func (s *sweepLoop) processToken(ctx context.Context, token common.Address, tokCfg tokenConfig) {
	cfg := s.cfg
	now := time.Now()
	lastSwept := cfg.SweepClock.LastSwept(token)
	elapsed := now.Sub(lastSwept)

	// Cadence floor with high-volume override. If cadence hasn't elapsed,
	// check whether enough rows have accumulated that we can sweep at a
	// per-row cost within our estimated overhead even at current gas. If
	// so, bypass cadence (and below, the gas cap too). This is the "we
	// don't need to wait when volume already justifies sweeping" path.
	earlyOverride := false
	if tokCfg.SweepCadence > 0 && elapsed < tokCfg.SweepCadence {
		ok, nRows, estGasWei := s.shouldOverrideCadenceEarly(ctx, token, lastSwept)
		if !ok {
			return
		}
		earlyOverride = true
		cfg.Logger.Info("cadence_override_triggered",
			slog.String("token", token.Hex()),
			slog.String("tier", tokCfg.Tier.String()),
			slog.Duration("elapsed", elapsed),
			slog.Duration("cadence", tokCfg.SweepCadence),
			slog.Int("rows_since_last_sweep", nRows),
			slog.Float64("est_sweep_gas_eth", weiToEth(estGasWei)),
			slog.String("reason", "high_volume_amortizes_below_per_row_estimate"))
	}

	balance, err := s.fetchBalance(ctx, token)
	if err != nil {
		cfg.Logger.Warn("balanceOf failed", slog.String("token", token.Hex()), slog.Any("error", err))
		return
	}
	if balance.Cmp(big.NewInt(sweepDustRawUnits)) <= 0 {
		// Dust — log only at debug, this is fine.
		cfg.Logger.Debug("sweep_evaluation_skipped_dust",
			slog.String("token", token.Hex()),
			slog.String("balance", balance.String()))
		return
	}

	gasPrice, err := cfg.Client.SuggestGasPrice(ctx)
	if err != nil {
		cfg.Logger.Warn("suggest gas price failed", slog.Any("error", err))
		return
	}

	var decision sweepDecisionOutput
	if earlyOverride {
		// Volume-based override: skip the gas-cap percentile gate entirely.
		// Profitability gate still applies below — we never sweep at a loss.
		decision = sweepDecisionOutput{
			Decision: SweepForce,
			Reason:   "early_volume_override",
		}
	} else {
		decision = decideSweep(sweepDecisionInput{
			Now:           now,
			Cfg:           tokCfg,
			LastSweepAt:   lastSwept,
			CurrentGasWei: gasPrice.Uint64(),
			Buf:           cfg.GasBuffer,
		})
	}

	cfg.Logger.Debug("sweep_evaluation",
		slog.String("token", token.Hex()),
		slog.String("tier", tokCfg.Tier.String()),
		slog.Duration("since_last_sweep", elapsed),
		slog.String("balance", balance.String()),
		slog.Uint64("gas_wei", gasPrice.Uint64()),
		slog.Uint64("gas_cap_wei", decision.GasCapWei),
		slog.Bool("has_gas_data", decision.HasGasData),
		slog.String("decision", decision.Decision.String()),
		slog.String("reason", decision.Reason),
		slog.Bool("early_override", earlyOverride))

	if decision.Decision == SweepSkip {
		cfg.Logger.Info("sweep_skipped",
			slog.String("token", token.Hex()),
			slog.String("tier", tokCfg.Tier.String()),
			slog.String("reason", decision.Reason),
			slog.Uint64("gas_wei", gasPrice.Uint64()),
			slog.Uint64("gas_cap_wei", decision.GasCapWei),
			slog.Duration("since_last_sweep", elapsed))
		return
	}

	// Attempt or Force: get a fresh Barter quote for the entire balance.
	barterResp, err := callBarter(ctx, cfg.HTTPClient, cfg.BarterURL, cfg.BarterKey, barterRequest{
		Source:            token.Hex(),
		Target:            cfg.WETH.Hex(),
		SellAmount:        balance.String(),
		Recipient:         cfg.ExecutorAddr.Hex(),
		Origin:            cfg.ExecutorAddr.Hex(),
		MinReturnFraction: sweepBarterMinReturnFraction,
		Deadline:          fmt.Sprintf("%d", now.Add(10*time.Minute).Unix()),
	})
	if err != nil {
		// "Amount too low" is steady-state expected for tokens with small
		// accumulated balances; demote to debug to keep the warn stream
		// clean. Other errors (timeouts, malformed responses) still warn.
		if isBarterAmountTooLow(err) {
			cfg.Logger.Debug("cadence sweep skipped: balance under barter minimum",
				slog.String("token", token.Hex()), slog.String("balance", balance.String()))
		} else {
			cfg.Logger.Warn("cadence sweep barter quote failed",
				slog.String("token", token.Hex()), slog.Any("error", err))
		}
		return
	}

	gasLimit, err := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	if err != nil {
		cfg.Logger.Warn("invalid gasLimit from barter",
			slog.String("token", token.Hex()), slog.String("gasLimit", barterResp.GasLimit))
		return
	}
	gasLimit += 50000

	expectedGas := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
	expectedReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
	if !ok {
		cfg.Logger.Warn("invalid MinReturn from barter",
			slog.String("token", token.Hex()), slog.String("minReturn", barterResp.MinReturn))
		return
	}

	// Profitability guard: expected_return > sweep_gas × 1.05.
	// This is the absolute floor — even force-sweep applies it. We never
	// sweep at a loss.
	threshold := new(big.Int).Mul(expectedGas, big.NewInt(sweepProfitabilityNumerator))
	threshold.Div(threshold, big.NewInt(100))
	if expectedReturn.Cmp(threshold) <= 0 {
		cfg.Logger.Info("sweep_skipped",
			slog.String("token", token.Hex()),
			slog.String("tier", tokCfg.Tier.String()),
			slog.String("reason", "unprofitable"),
			slog.String("decision", decision.Decision.String()),
			slog.Float64("expected_return_eth", weiToEth(expectedReturn)),
			slog.Float64("sweep_gas_eth", weiToEth(expectedGas)),
			slog.Float64("threshold_eth", weiToEth(threshold)))
		return
	}

	if cfg.DryRun {
		cfg.Logger.Info("sweep_executed_dryrun",
			slog.String("token", token.Hex()),
			slog.String("tier", tokCfg.Tier.String()),
			slog.String("decision", decision.Decision.String()),
			slog.Float64("expected_return_eth", weiToEth(expectedReturn)),
			slog.Float64("expected_gas_eth", weiToEth(expectedGas)))
		cfg.SweepClock.MarkSwept(token, now)
		return
	}

	actualReturn, actualGas, err := submitFastSwapSweep(
		ctx, cfg.Logger, cfg.Client, cfg.L1Client, cfg.HTTPClient, cfg.Signer,
		cfg.ExecutorAddr, token, balance, cfg.FastSwapURL, cfg.FundsRecipient,
		cfg.SettlementAddr, barterResp, cfg.MaxGasGwei,
	)
	if err != nil {
		cfg.Logger.Error("cadence sweep failed",
			slog.String("token", token.Hex()), slog.Any("error", err))
		return
	}

	cfg.SweepClock.MarkSwept(token, now)
	actualNetEth := new(big.Int).Sub(actualReturn, actualGas)
	cfg.Logger.Info("sweep_executed",
		slog.String("token", token.Hex()),
		slog.String("tier", tokCfg.Tier.String()),
		slog.String("decision", decision.Decision.String()),
		slog.String("balance_swept", balance.String()),
		slog.Float64("actual_return_eth", weiToEth(actualReturn)),
		slog.Float64("actual_gas_eth", weiToEth(actualGas)),
		slog.Float64("actual_net_eth", weiToEth(actualNetEth)))

	// Compare the just-realized sweep profit against the miles obligation
	// already accrued for this token's recent rows. Surfaces the question
	// "are we awarding more miles than we're earning?" per-sweep, not just
	// in the hourly aggregate. Best-effort — a query failure shouldn't
	// hold up the sweep marking.
	if obligation, err := s.queryRecentMilesObligation(ctx, token); err == nil {
		cfg.Logger.Info("sweep_vs_obligation",
			slog.String("token", token.Hex()),
			slog.Float64("realized_net_eth", weiToEth(actualNetEth)),
			slog.Float64("recent_miles_obligation_eth", obligation.milesEth),
			slog.Int("recent_miles_rows", obligation.nRows),
			slog.Int("lookback_days", sweepObligationLookbackDays),
			slog.Float64("ratio", safeRatio(obligation.milesEth, weiToEth(actualNetEth))))
	} else {
		cfg.Logger.Debug("sweep vs obligation query failed",
			slog.String("token", token.Hex()), slog.Any("error", err))
	}
}

// sweepObligationLookbackDays bounds how far back we sum miles for the
// post-sweep comparison. Matched to the cadence-period scale; long enough
// that a stable token's last few sweeps fall in the window, short enough
// that it reflects recent estimate accuracy rather than ancient history.
const sweepObligationLookbackDays = 7

type recentObligation struct {
	milesEth float64
	nRows    int
}

func (s *sweepLoop) queryRecentMilesObligation(ctx context.Context, token common.Address) (recentObligation, error) {
	var milesSum int64
	var n int
	err := s.cfg.DB.QueryRowContext(ctx, fmt.Sprintf(`
SELECT COALESCE(SUM(miles), 0), COUNT(*)
FROM mevcommit_57173.fastswap_miles
WHERE processed = 1
  AND miles > 0
  AND swap_type = 'erc20'
  AND LOWER(output_token) = ?
  AND LOWER(user_address) != ?
  AND block_timestamp >= NOW() - INTERVAL %d DAY
`, sweepObligationLookbackDays),
		strings.ToLower(token.Hex()),
		strings.ToLower(s.cfg.ExecutorAddr.Hex()),
	).Scan(&milesSum, &n)
	if err != nil {
		return recentObligation{}, err
	}
	return recentObligation{
		milesEth: float64(milesSum) * float64(weiPerPoint) / 1e18,
		nRows:    n,
	}, nil
}

func safeRatio(num, den float64) float64 {
	if den == 0 {
		return 0
	}
	return num / den
}

// isBarterAmountTooLow reports whether a Barter error is the "Amount too low"
// response — Barter rejects inputs below its per-token minimum, and this is
// the dominant steady-state error for early/small batches. Callers use it to
// demote noise from warn to debug while keeping real errors loud.
func isBarterAmountTooLow(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Amount too low") || strings.Contains(msg, "amount too low")
}

// cadenceOverrideMargin is the safety multiplier on estimated sweep gas
// when checking whether volume justifies an early sweep. A value of 1.2
// means we require the row-count budget to exceed sweep gas by 20% before
// bypassing cadence — buffer for gas drift between estimate and execution.
const cadenceOverrideMargin = 1.2

// typicalSweepGasLimit is a conservative gas-limit estimate used purely to
// size the "is current gas affordable for sweeping right now" check
// without making an extra Barter call. Real sweeps will discover their
// actual gasLimit from the Barter response. Slightly overestimating here
// just makes the early-override harder to trigger — preferable to under-
// estimating and overriding when sweep would actually be expensive.
const typicalSweepGasLimit = 500_000

// cadenceOverrideLookbackCap bounds the COUNT(*) query when a token has
// never been swept (lastSwept is zero). 14 days mirrors the cost
// estimator's lookback so the threshold uses comparable row volumes.
const cadenceOverrideLookbackCap = 14 * 24 * time.Hour

// shouldOverrideCadenceEarly reports whether accumulated row volume since
// the last sweep is high enough that the realized sweep gas can be
// amortized below our per-row overhead estimate at current gas levels.
//
//	threshold:  N × per_row_overhead_estimate  >=  cadenceOverrideMargin × est_sweep_gas
//
// This says: we have enough rows that even at the current gas, sweeping
// would land us at or under our estimated per-row cost. No reason to wait
// for cadence — bigger batches won't materially improve the protocol P&L
// from here, and waiting just adds inventory risk on stables.
//
// Returns (override, nRows, est_sweep_gas_wei) for caller logging.
func (s *sweepLoop) shouldOverrideCadenceEarly(
	ctx context.Context, token common.Address, lastSwept time.Time,
) (bool, int, *big.Int) {
	cfg := s.cfg
	if cfg.CostEstimator == nil {
		return false, 0, nil
	}
	est := cfg.CostEstimator.Get(strings.ToLower(token.Hex()))
	if est.PerRowOverhead <= 0 {
		return false, 0, nil
	}

	// Recent gas median — cheaper than another SuggestGasPrice and reflects
	// typical gas in the window we'd actually sweep at. Cold start (no
	// samples yet): cannot evaluate; never override.
	medianGas, ok := cfg.GasBuffer.Percentile(50, 30*time.Minute)
	if !ok {
		return false, 0, nil
	}
	estSweepGasWei := new(big.Int).Mul(
		new(big.Int).SetUint64(medianGas), big.NewInt(typicalSweepGasLimit))
	estSweepGasEth := weiToEth(estSweepGasWei)

	sinceTS := lastSwept
	if sinceTS.IsZero() {
		sinceTS = time.Now().Add(-cadenceOverrideLookbackCap)
	}
	var nRows int
	err := cfg.DB.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM mevcommit_57173.fastswap_miles
WHERE swap_type = 'erc20'
  AND LOWER(output_token) = ?
  AND block_timestamp >= ?
`, strings.ToLower(token.Hex()), sinceTS).Scan(&nRows)
	if err != nil {
		cfg.Logger.Debug("cadence override row-count query failed",
			slog.String("token", token.Hex()), slog.Any("error", err))
		return false, 0, estSweepGasWei
	}
	if nRows == 0 {
		return false, 0, estSweepGasWei
	}

	return cadenceOverrideMet(nRows, est.PerRowOverhead, estSweepGasEth), nRows, estSweepGasWei
}

// cadenceOverrideMet applies the override threshold:
//
//	N × per_row_overhead >= cadenceOverrideMargin × estimated_sweep_gas_eth
//
// Separated out so the threshold math is unit-testable without needing a DB
// or gas-buffer fixture.
func cadenceOverrideMet(nRows int, perRowOverheadEth, estSweepGasEth float64) bool {
	if nRows <= 0 || perRowOverheadEth <= 0 {
		return false
	}
	budget := float64(nRows) * perRowOverheadEth
	threshold := estSweepGasEth * cadenceOverrideMargin
	return budget >= threshold
}

func (s *sweepLoop) fetchBalance(ctx context.Context, token common.Address) (*big.Int, error) {
	data, err := s.balance.Pack("balanceOf", s.cfg.ExecutorAddr)
	if err != nil {
		return nil, fmt.Errorf("pack balanceOf: %w", err)
	}
	raw, err := s.cfg.Client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("call balanceOf: %w", err)
	}
	out, err := s.balance.Unpack("balanceOf", raw)
	if err != nil {
		return nil, fmt.Errorf("unpack balanceOf: %w", err)
	}
	if len(out) < 1 {
		return nil, fmt.Errorf("empty balanceOf output")
	}
	bal, ok := out[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("balanceOf returned %T not *big.Int", out[0])
	}
	return bal, nil
}
