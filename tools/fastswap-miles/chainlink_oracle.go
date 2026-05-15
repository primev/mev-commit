package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// chainlinkFeedRegistryAddr is the Chainlink Feed Registry on Ethereum
// mainnet. A single contract that maps (base, quote) → underlying feed and
// exposes latestRoundData / decimals directly. New pairs were frozen years
// ago, so some legitimate tokens (PEPE, ARB, possibly others) revert; the
// caller routes those rows to sweep-time pricing instead.
const chainlinkFeedRegistryAddr = "0x47Fb2585D2C56Fe188D0E6ec628a38b74fCeeeDf"

// chainlinkEthDenomination is the sentinel address Chainlink uses to
// represent ETH on the quote side of a (base, quote) pair.
const chainlinkEthDenomination = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

// priceOracleCacheTTL bounds how long a fetched Chainlink rate is reused.
// Major Chainlink feeds heartbeat hourly and update on price-deviation
// thresholds; 5 min is comfortably tighter than the heartbeat and well
// within miles-grade precision.
const priceOracleCacheTTL = 5 * time.Minute

const chainlinkFeedRegistryABI = `[
	{"inputs":[{"name":"base","type":"address"},{"name":"quote","type":"address"}],
	 "name":"latestRoundData",
	 "outputs":[{"name":"roundId","type":"uint80"},{"name":"answer","type":"int256"},{"name":"startedAt","type":"uint256"},{"name":"updatedAt","type":"uint256"},{"name":"answeredInRound","type":"uint80"}],
	 "stateMutability":"view","type":"function"},
	{"inputs":[{"name":"base","type":"address"},{"name":"quote","type":"address"}],
	 "name":"decimals",
	 "outputs":[{"name":"","type":"uint8"}],
	 "stateMutability":"view","type":"function"}
]`

const erc20DecimalsABI = `[
	{"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"stateMutability":"view","type":"function"}
]`

// priceOracle resolves the ETH-wei value of a fastswap surplus at miles-award
// time. Three sources, picked per swap shape:
//
//   - Output is ETH/WETH: handled upstream (surplus is already ETH wei).
//   - ETH/WETH input + whitelisted ERC20 output: event-derived from the
//     trade's executed rate (most accurate possible — it IS the realized
//     rate of this exact swap).
//   - ERC20 input + whitelisted ERC20 output: Chainlink Feed Registry. Reads
//     are out-of-band so flash-loan manipulation doesn't apply; the oracle
//     network itself is not pool-manipulable.
//
// Non-whitelisted output tokens are always deferred to sweep time. This is
// the structural defense against the attacker-token attack: a malicious
// actor who mints their own token and controls its on-chain liquidity
// cannot extract upfront miles, because their token isn't in tokenConfigs.
type priceOracle struct {
	client *ethclient.Client
	logger *slog.Logger
	weth   common.Address

	registryAddr common.Address
	registryABI  abi.ABI
	erc20ABI     abi.ABI

	mu            sync.RWMutex
	rateCache     map[common.Address]rateEntry
	decimalsCache map[common.Address]uint8
}

type rateEntry struct {
	answer       *big.Int
	feedDecimals uint8
	fetchedAt    time.Time
}

func newPriceOracle(client *ethclient.Client, weth common.Address, logger *slog.Logger) (*priceOracle, error) {
	regABI, err := abi.JSON(strings.NewReader(chainlinkFeedRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("parse chainlink registry ABI: %w", err)
	}
	ercABI, err := abi.JSON(strings.NewReader(erc20DecimalsABI))
	if err != nil {
		return nil, fmt.Errorf("parse erc20 decimals ABI: %w", err)
	}
	return &priceOracle{
		client:        client,
		logger:        logger,
		weth:          weth,
		registryAddr:  common.HexToAddress(chainlinkFeedRegistryAddr),
		registryABI:   regABI,
		erc20ABI:      ercABI,
		rateCache:     map[common.Address]rateEntry{},
		decimalsCache: map[common.Address]uint8{},
	}, nil
}

// PriceSurplusEth returns the ETH-wei value of an ERC20-output surplus.
//
// Returns (eth wei, eligible, source). When eligible is false the caller MUST
// route the row to sweep-time pro-rata pricing instead — this is the
// attacker-token defense and the no-Chainlink-feed fallback rolled into one.
//
// Source values for logging:
//
//	"event_derived"            ETH/WETH input + whitelisted output
//	"chainlink"                ERC20 input + whitelisted output, Registry hit
//	"deferred:not_whitelisted" output token not in tokenConfigs
//	"deferred:invalid_event"   event surplus + userAmtOut sum to zero
//	"deferred:no_chainlink"    Registry call reverted or returned bad data
//	"deferred:no_token_decim"  ERC20 decimals call failed
func (o *priceOracle) PriceSurplusEth(
	ctx context.Context,
	inputToken, outputToken common.Address,
	inputAmt, userAmtOut, surplus *big.Int,
) (*big.Int, bool, string) {
	if !isWhitelisted(outputToken) {
		return nil, false, "deferred:not_whitelisted"
	}

	if inputToken == zeroAddr || inputToken == o.weth {
		v := deriveEthInputSurplusEth(inputAmt, userAmtOut, surplus)
		if v == nil {
			return nil, false, "deferred:invalid_event"
		}
		return v, true, "event_derived"
	}

	rate, feedDecimals, ok := o.getChainlinkRate(ctx, outputToken)
	if !ok {
		return nil, false, "deferred:no_chainlink"
	}
	tokenDecimals, ok := o.getTokenDecimals(ctx, outputToken)
	if !ok {
		return nil, false, "deferred:no_token_decim"
	}
	return scaleChainlinkAnswer(surplus, rate, tokenDecimals, feedDecimals), true, "chainlink"
}

// deriveEthInputSurplusEth computes surplus_eth from event data alone when
// the input was ETH or WETH (so inputAmt is denominated in ETH wei).
//
// surplus_eth = surplus × inputAmt / (userAmtOut + surplus)
//
// This IS the trade's executed exchange rate — the contract took inputAmt
// wei from the user, Barter returned (userAmtOut + surplus) tokens. No
// external oracle can be more accurate than the trade's own realized price.
// Whitelist-gated to prevent attacker-controlled-pool manipulation.
//
// Returns nil for degenerate events where the divisor is non-positive.
func deriveEthInputSurplusEth(inputAmt, userAmtOut, surplus *big.Int) *big.Int {
	if inputAmt == nil || userAmtOut == nil || surplus == nil {
		return nil
	}
	denom := new(big.Int).Add(userAmtOut, surplus)
	if denom.Sign() <= 0 {
		return nil
	}
	result := new(big.Int).Mul(surplus, inputAmt)
	result.Div(result, denom)
	return result
}

// scaleChainlinkAnswer converts a Chainlink rate into surplus_eth.
//
//	surplus_eth_wei = surplus_raw × answer × 10^(18 - feed_decimals) / 10^token_decimals
//
// Combined into one divisor exponent (token_decimals + feed_decimals - 18)
// which can be negative — handled both directions.
func scaleChainlinkAnswer(surplus, answer *big.Int, tokenDecimals, feedDecimals uint8) *big.Int {
	result := new(big.Int).Mul(surplus, answer)
	expDiv := int(tokenDecimals) + int(feedDecimals) - 18
	switch {
	case expDiv > 0:
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(expDiv)), nil)
		result.Div(result, divisor)
	case expDiv < 0:
		multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-expDiv)), nil)
		result.Mul(result, multiplier)
	}
	return result
}

func (o *priceOracle) getChainlinkRate(ctx context.Context, token common.Address) (*big.Int, uint8, bool) {
	o.mu.RLock()
	if entry, ok := o.rateCache[token]; ok && time.Since(entry.fetchedAt) < priceOracleCacheTTL {
		o.mu.RUnlock()
		return entry.answer, entry.feedDecimals, true
	}
	o.mu.RUnlock()

	answer, decimals, err := o.fetchChainlinkRate(ctx, token)
	if err != nil {
		o.logger.Warn("chainlink registry lookup failed; deferring to sweep for this token",
			slog.String("token", token.Hex()), slog.Any("error", err))
		return nil, 0, false
	}

	o.mu.Lock()
	o.rateCache[token] = rateEntry{answer: answer, feedDecimals: decimals, fetchedAt: time.Now()}
	o.mu.Unlock()
	return answer, decimals, true
}

func (o *priceOracle) fetchChainlinkRate(ctx context.Context, token common.Address) (*big.Int, uint8, error) {
	ethSentinel := common.HexToAddress(chainlinkEthDenomination)

	answer, err := o.callRegistryReturningBigInt(ctx, "latestRoundData", token, ethSentinel, 1)
	if err != nil {
		return nil, 0, err
	}
	if answer.Sign() <= 0 {
		return nil, 0, fmt.Errorf("non-positive chainlink answer: %s", answer.String())
	}

	decimals, err := o.callRegistryReturningUint8(ctx, "decimals", token, ethSentinel)
	if err != nil {
		return nil, 0, err
	}
	return answer, decimals, nil
}

// callRegistryReturningBigInt invokes a registry method that returns one or
// more values, extracting the *big.Int at outputIndex.
func (o *priceOracle) callRegistryReturningBigInt(
	ctx context.Context, method string, base, quote common.Address, outputIndex int,
) (*big.Int, error) {
	data, err := o.registryABI.Pack(method, base, quote)
	if err != nil {
		return nil, fmt.Errorf("pack %s: %w", method, err)
	}
	raw, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &o.registryAddr, Data: data}, nil)
	if err != nil {
		return nil, fmt.Errorf("call %s: %w", method, err)
	}
	out, err := o.registryABI.Unpack(method, raw)
	if err != nil {
		return nil, fmt.Errorf("unpack %s: %w", method, err)
	}
	if len(out) <= outputIndex {
		return nil, fmt.Errorf("unexpected %s output length %d", method, len(out))
	}
	v, ok := out[outputIndex].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("%s output[%d] is %T, want *big.Int", method, outputIndex, out[outputIndex])
	}
	return v, nil
}

func (o *priceOracle) callRegistryReturningUint8(
	ctx context.Context, method string, base, quote common.Address,
) (uint8, error) {
	data, err := o.registryABI.Pack(method, base, quote)
	if err != nil {
		return 0, fmt.Errorf("pack %s: %w", method, err)
	}
	raw, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &o.registryAddr, Data: data}, nil)
	if err != nil {
		return 0, fmt.Errorf("call %s: %w", method, err)
	}
	out, err := o.registryABI.Unpack(method, raw)
	if err != nil {
		return 0, fmt.Errorf("unpack %s: %w", method, err)
	}
	if len(out) < 1 {
		return 0, fmt.Errorf("empty %s output", method)
	}
	v, ok := out[0].(uint8)
	if !ok {
		return 0, fmt.Errorf("%s output is %T, want uint8", method, out[0])
	}
	return v, nil
}

func (o *priceOracle) getTokenDecimals(ctx context.Context, token common.Address) (uint8, bool) {
	o.mu.RLock()
	if dec, ok := o.decimalsCache[token]; ok {
		o.mu.RUnlock()
		return dec, true
	}
	o.mu.RUnlock()

	data, err := o.erc20ABI.Pack("decimals")
	if err != nil {
		o.logger.Warn("erc20 decimals pack failed", slog.String("token", token.Hex()), slog.Any("error", err))
		return 0, false
	}
	raw, err := o.client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: data}, nil)
	if err != nil {
		o.logger.Warn("erc20 decimals call failed", slog.String("token", token.Hex()), slog.Any("error", err))
		return 0, false
	}
	out, err := o.erc20ABI.Unpack("decimals", raw)
	if err != nil || len(out) < 1 {
		o.logger.Warn("erc20 decimals unpack failed", slog.String("token", token.Hex()), slog.Any("error", err))
		return 0, false
	}
	dec, ok := out[0].(uint8)
	if !ok {
		o.logger.Warn("erc20 decimals not uint8", slog.String("token", token.Hex()))
		return 0, false
	}

	o.mu.Lock()
	o.decimalsCache[token] = dec
	o.mu.Unlock()
	return dec, true
}
