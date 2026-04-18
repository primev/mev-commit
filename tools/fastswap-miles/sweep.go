package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	fastsettlement "github.com/primev/mev-commit/contracts-abi/clients/FastSettlementV3"
	"github.com/primev/mev-commit/x/keysigner"
)

// -------------------- Event Indexer --------------------

func indexBatch(
	ctx context.Context,
	logger *slog.Logger,
	filterer *fastsettlement.Fastsettlementv3Filterer,
	client *ethclient.Client,
	db *sql.DB,
	weth common.Address,
	from, to uint64,
) (int, error) {
	opts := &bind.FilterOpts{
		Start:   from,
		End:     &to,
		Context: ctx,
	}

	iter, err := filterer.FilterIntentExecuted(opts, nil, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("FilterIntentExecuted: %w", err)
	}
	defer func() { _ = iter.Close() }()

	count := 0
	for iter.Next() {
		ev := iter.Event
		if ev.Surplus == nil || ev.Surplus.Sign() == 0 {
			continue
		}

		txHash := ev.Raw.TxHash.Hex()
		blockNum := ev.Raw.BlockNumber

		var receipt *types.Receipt
		for attempt := 0; attempt < 3; attempt++ {
			receipt, err = client.TransactionReceipt(ctx, ev.Raw.TxHash)
			if err == nil {
				break
			}
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "Too Many Requests") {
				time.Sleep(time.Duration(attempt+1) * 2 * time.Second)
				continue
			}
			break
		}
		if err != nil {
			logger.Warn("receipt fetch failed, skipping gas cost", slog.String("tx", txHash), slog.Any("error", err))
		}
		var gasCost *big.Int
		if receipt != nil {
			gasCost = new(big.Int).Mul(
				new(big.Int).SetUint64(receipt.GasUsed),
				receipt.EffectiveGasPrice,
			)
		}

		var header *types.Header
		for attempt := 0; attempt < 3; attempt++ {
			header, err = client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
			if err == nil {
				break
			}
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "Too Many Requests") {
				time.Sleep(time.Duration(attempt+1) * 2 * time.Second)
				continue
			}
			break
		}
		if err != nil {
			logger.Warn("header fetch failed", slog.Uint64("block", blockNum), slog.Any("error", err))
		}
		var blockTS *time.Time
		if header != nil {
			t := time.Unix(int64(header.Time), 0).UTC()
			blockTS = &t
		}

		swapType := "erc20"
		if ev.OutputToken == zeroAddr || strings.EqualFold(ev.OutputToken.Hex(), weth.Hex()) {
			swapType = "eth_weth"
		}

		err = insertEvent(db, txHash, blockNum, blockTS, ev, gasCost, swapType)
		if err != nil {
			logger.Warn("insertEvent failed", slog.String("tx", txHash), slog.Any("error", err))
			continue
		}
		count++
	}
	if iter.Error() != nil {
		return count, fmt.Errorf("iter: %w", iter.Error())
	}
	return count, nil
}

func insertEvent(
	db *sql.DB,
	txHash string,
	blockNum uint64,
	blockTS *time.Time,
	ev *fastsettlement.Fastsettlementv3IntentExecuted,
	gasCost *big.Int,
	swapType string,
) error {
	// CRITICAL: the fastswap_miles table uses StarRocks `PRIMARY KEY(tx_hash)`
	// model, which means an unconditional INSERT UPSERTS the entire row and
	// resets every column we don't specify (processed → false, miles → NULL,
	// surplus_eth → NULL, bid_cost → NULL, fuel_submitted_at → NULL, …).
	//
	// If the indexer rescans a block (pod restart with last_block reset, an
	// explicit -start-block flag going backward, meta row lost during a
	// deploy, etc.) an unconditional INSERT would clobber already-processed
	// rows back to the pending state — and the miles pipeline would then
	// re-submit each one to Fuel, double-crediting users.
	//
	// The 2026-04-16 double-credit incident was caused by exactly this: the
	// Docker-support deploy restarted the pod, the indexer re-walked historical
	// blocks, and every re-inserted event wiped `processed` on the existing
	// row. 78 events ended up re-submitted to Fuel for a single test user
	// alone; the protocol-wide overcount was much larger.
	//
	var tsVal interface{} = nil
	if blockTS != nil {
		tsVal = *blockTS
	}
	var gcStr interface{} = nil
	if gasCost != nil {
		gcStr = gasCost.String()
	}

	// Fix: check for existence before inserting. The IntentExecuted event
	// args themselves are immutable once L1-finalized, so the row's core
	// fields (user, tokens, amounts, surplus) must never be replaced.
	//
	// For rows that already exist: run a COALESCE-only UPDATE that fills in
	// gas_cost or block_timestamp IF they were previously NULL (which happens
	// when indexBatch caught a transient receipt/header RPC failure on the
	// first pass). This preserves every derived column (processed, miles,
	// surplus_eth, net_profit_eth, bid_cost, fuel_submitted_at) — so a rescan
	// can heal partial metadata without destroying pipeline state.
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM mevcommit_57173.fastswap_miles WHERE tx_hash = ?)`,
		txHash,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check existing row: %w", err)
	}
	if exists {
		// Backfill only the two fields that can legitimately arrive NULL from
		// transient RPC failures. COALESCE(col, newVal) keeps the existing
		// value if it's non-NULL, and substitutes newVal otherwise. If newVal
		// is also NULL (we still don't have fresh data), the column is
		// unchanged — no-op.
		_, err := db.Exec(`
UPDATE mevcommit_57173.fastswap_miles
SET gas_cost = COALESCE(gas_cost, ?),
    block_timestamp = COALESCE(block_timestamp, ?)
WHERE tx_hash = ?`, gcStr, tsVal, txHash)
		if err != nil {
			return fmt.Errorf("backfill null metadata: %w", err)
		}
		return nil
	}

	_, err = db.Exec(`
INSERT INTO mevcommit_57173.fastswap_miles (
  tx_hash, block_number, block_timestamp, user_address,
  input_token, output_token, input_amount, user_amt_out,
  surplus, gas_cost, swap_type, processed
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, false)`,
		txHash,
		blockNum,
		tsVal,
		strings.ToLower(ev.User.Hex()),
		strings.ToLower(ev.InputToken.Hex()),
		strings.ToLower(ev.OutputToken.Hex()),
		ev.InputAmt.String(),
		ev.UserAmtOut.String(),
		ev.Surplus.String(),
		gcStr,
		swapType,
	)
	return err
}

// -------------------- Token Sweep --------------------

// Minimal ERC20 ABI for Approve
const erc20ApproveABI = `[{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

func submitFastSwapSweep(
	ctx context.Context,
	logger *slog.Logger,
	client *ethclient.Client,
	l1Client *ethclient.Client,
	httpClient *http.Client,
	signer keysigner.KeySigner,
	executorAddr common.Address,
	tokenAddr common.Address,
	totalAmount *big.Int,
	fastswapURL string,
	fundsRecipient common.Address,
	settlementAddr common.Address,
	barterResp *BarterResponse,
	maxGasGwei uint64,
) (*big.Int, *big.Int, error) {

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("gas price: %w", err)
	}
	gasPriceGwei := new(big.Int).Div(gasPrice, big.NewInt(1_000_000_000))
	if gasPriceGwei.Uint64() > maxGasGwei {
		return nil, nil, fmt.Errorf("gas price %d gwei exceeds max %d gwei, skipping sweep", gasPriceGwei.Uint64(), maxGasGwei)
	}

	if err := ensurePermit2Approval(ctx, logger, client, l1Client, signer, executorAddr, tokenAddr, totalAmount); err != nil {
		return nil, nil, fmt.Errorf("permit2 approval: %w", err)
	}

	barterMinReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
	if !ok {
		return nil, nil, fmt.Errorf("invalid MinReturn from barter: %s", barterResp.MinReturn)
	}
	userAmtOut := new(big.Int).Mul(barterMinReturn, big.NewInt(95))
	userAmtOut.Div(userAmtOut, big.NewInt(100))

	deadline := big.NewInt(time.Now().Add(10 * time.Minute).Unix())

	nonceBuf := make([]byte, 32)
	if _, err := rand.Read(nonceBuf); err != nil {
		return nil, nil, fmt.Errorf("generate nonce: %w", err)
	}
	nonce := new(big.Int).SetBytes(nonceBuf)

	signature, err := signPermit2Witness(
		signer,
		tokenAddr,
		totalAmount,
		settlementAddr,
		nonce,
		deadline,
		executorAddr,
		tokenAddr,
		common.Address{},
		totalAmount,
		userAmtOut,
		fundsRecipient,
		deadline,
		nonce,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("sign permit2: %w", err)
	}

	swapReq := map[string]string{
		"user":        executorAddr.Hex(),
		"inputToken":  tokenAddr.Hex(),
		"outputToken": common.Address{}.Hex(),
		"inputAmt":    totalAmount.String(),
		"userAmtOut":  userAmtOut.String(),
		"recipient":   fundsRecipient.Hex(),
		"deadline":    deadline.String(),
		"nonce":       nonce.String(),
		"signature":   "0x" + hex.EncodeToString(signature),
		"slippage":    "1.0",
	}

	reqBody, err := json.Marshal(swapReq)
	if err != nil {
		return nil, nil, fmt.Errorf("marshal swap request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fastswapURL+"/fastswap", bytes.NewReader(reqBody))
	if err != nil {
		return nil, nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fastswap API request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("fastswap API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		TxHash       string `json:"txHash"`
		OutputAmount string `json:"outputAmount"`
		GasLimit     uint64 `json:"gasLimit"`
		Status       string `json:"status"`
		Error        string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("decode response: %w", err)
	}

	if result.Status != "success" {
		return nil, nil, fmt.Errorf("fastswap returned error: %s", result.Error)
	}

	logger.Info("FastSwap sweep submitted",
		slog.String("tx", result.TxHash),
		slog.String("outputAmount", result.OutputAmount),
		slog.Uint64("gasLimit", result.GasLimit))

	gasLimit, _ := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	gasLimit += 100000
	estimatedGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)

	return userAmtOut, estimatedGasCost, nil
}

// -------------------- Permit2 --------------------

func ensurePermit2Approval(
	ctx context.Context,
	logger *slog.Logger,
	client *ethclient.Client,
	l1Client *ethclient.Client,
	signer keysigner.KeySigner,
	owner common.Address,
	token common.Address,
	requiredAmount *big.Int,
) error {
	permit2 := common.HexToAddress(permit2Addr)

	parsedABI, err := abi.JSON(strings.NewReader(erc20ApproveABI))
	if err != nil {
		return fmt.Errorf("parse erc20 ABI: %w", err)
	}

	allowanceData, err := parsedABI.Pack("allowance", owner, permit2)
	if err != nil {
		return fmt.Errorf("pack allowance call: %w", err)
	}

	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &token,
		Data: allowanceData,
	}, nil)
	if err != nil {
		return fmt.Errorf("call allowance: %w", err)
	}

	currentAllowance := new(big.Int).SetBytes(result)
	if currentAllowance.Cmp(requiredAmount) >= 0 {
		logger.Debug("Permit2 allowance sufficient",
			slog.String("token", token.Hex()),
			slog.String("have", currentAllowance.String()),
			slog.String("need", requiredAmount.String()))
		return nil
	}

	logger.Info("Permit2 allowance insufficient, approving max",
		slog.String("token", token.Hex()),
		slog.String("have", currentAllowance.String()),
		slog.String("need", requiredAmount.String()))

	maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	approveData, err := parsedABI.Pack("approve", permit2, maxUint256)
	if err != nil {
		return fmt.Errorf("pack approve: %w", err)
	}

	nonce, err := client.NonceAt(ctx, owner, nil)
	if err != nil {
		return fmt.Errorf("nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("gas price: %w", err)
	}

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return fmt.Errorf("network id: %w", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      100000,
		To:       &token,
		Data:     approveData,
	})

	signedTx, err := signer.SignTx(tx, chainID)
	if err != nil {
		return fmt.Errorf("sign approve tx: %w", err)
	}

	if err := l1Client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("send approve tx: %w", err)
	}
	logger.Info("Permit2 approve tx sent", slog.String("tx", signedTx.Hash().Hex()))

	approveDeadline := time.Now().Add(15 * time.Minute)
	for time.Now().Before(approveDeadline) {
		time.Sleep(12 * time.Second)
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			continue
		}
		if receipt.Status != 1 {
			return fmt.Errorf("approve tx reverted: %s", signedTx.Hash().Hex())
		}
		logger.Info("Permit2 approve confirmed", slog.String("tx", signedTx.Hash().Hex()))
		return nil
	}
	return fmt.Errorf("approve tx not confirmed after 15 min, will retry next cycle: %s", signedTx.Hash().Hex())
}

func signPermit2Witness(
	signer keysigner.KeySigner,
	token common.Address,
	amount *big.Int,
	spender common.Address,
	permitNonce *big.Int,
	permitDeadline *big.Int,
	user common.Address,
	inputToken common.Address,
	outputToken common.Address,
	inputAmt *big.Int,
	userAmtOut *big.Int,
	recipient common.Address,
	intentDeadline *big.Int,
	intentNonce *big.Int,
) ([]byte, error) {
	permit2 := common.HexToAddress(permit2Addr)

	domainSep := crypto.Keccak256(
		crypto.Keccak256([]byte("EIP712Domain(string name,uint256 chainId,address verifyingContract)")),
		crypto.Keccak256([]byte("Permit2")),
		padTo32(big.NewInt(1)),
		padTo32Address(permit2),
	)

	tokenPermissionsTypeHash := crypto.Keccak256([]byte("TokenPermissions(address token,uint256 amount)"))
	tokenPermissionsHash := crypto.Keccak256(
		tokenPermissionsTypeHash,
		padTo32Address(token),
		padTo32(amount),
	)

	intentTypeHash := crypto.Keccak256([]byte("Intent(address user,address inputToken,address outputToken,uint256 inputAmt,uint256 userAmtOut,address recipient,uint256 deadline,uint256 nonce)"))
	witnessHash := crypto.Keccak256(
		intentTypeHash,
		padTo32Address(user),
		padTo32Address(inputToken),
		padTo32Address(outputToken),
		padTo32(inputAmt),
		padTo32(userAmtOut),
		padTo32Address(recipient),
		padTo32(intentDeadline),
		padTo32(intentNonce),
	)

	permitWitnessTypeHash := crypto.Keccak256([]byte(
		"PermitWitnessTransferFrom(TokenPermissions permitted,address spender,uint256 nonce,uint256 deadline,Intent witness)" +
			"Intent(address user,address inputToken,address outputToken,uint256 inputAmt,uint256 userAmtOut,address recipient,uint256 deadline,uint256 nonce)" +
			"TokenPermissions(address token,uint256 amount)",
	))

	structHash := crypto.Keccak256(
		permitWitnessTypeHash,
		tokenPermissionsHash,
		padTo32Address(spender),
		padTo32(permitNonce),
		padTo32(permitDeadline),
		witnessHash,
	)

	digest := crypto.Keccak256(
		[]byte{0x19, 0x01},
		domainSep,
		structHash,
	)

	sig, err := signer.SignHash(digest)
	if err != nil {
		return nil, fmt.Errorf("sign hash: %w", err)
	}

	if len(sig) == 65 && sig[64] < 27 {
		sig[64] += 27
	}

	return sig, nil
}
