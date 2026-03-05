// fastswap-miles
//
// -------------------- HOW TO RUN --------------------
//
// REQUIRED ENV:
//   DB_USER, DB_PW, DB_HOST, DB_PORT, DB_NAME   (StarRocks / MySQL protocol)
//   L1_RPC_URL                                   (L1 Ethereum HTTP RPC)
//   FUEL_API_URL, FUEL_API_KEY                   (Fuel points API — optional in -dry-run mode)
//
// OPTIONAL CLI FLAGS:
//   -contract    FastSettlementV3 proxy address (default: 0x084C0EC7f5C0585195c1c713ED9f06272F48cB45)
//   -weth        WETH address (default: 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2)
//   -start-block block to start indexing from (default: 0 = auto-resume)
//   -poll        poll interval (default: 12s)
//   -batch       blocks per FilterLogs batch (default: 2000)
//   -dry-run     index events and compute miles but skip Fuel submission + processed marking
//
// ----------------------------------------------------

package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
	fastsettlement "github.com/primev/mev-commit/contracts-abi/clients/FastSettlementV3"
	"github.com/primev/mev-commit/x/keysigner"
)

// 0.00001 ETH in wei — same rate as preconf-rpc/points.
const weiPerPoint = 10_000_000_000_000

// WETH, zero address, and Permit2 constants.
const (
	defaultContract  = "0x084C0EC7f5C0585195c1c713ED9f06272F48cB45"
	defaultWETH      = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	defaultRecipient = "0xD5881f91270550B8850127f05BD6C8C203B3D33f"
	permit2Addr      = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
)

var zeroAddr = common.Address{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		contractAddr   = flag.String("contract", defaultContract, "FastSettlementV3 proxy address")
		wethAddr       = flag.String("weth", defaultWETH, "WETH contract address")
		startBlock     = flag.Uint64("start-block", 0, "block to start indexing from (0 = auto-resume from DB)")
		pollInterval   = flag.Duration("poll", 12*time.Second, "poll interval for new blocks")
		batchSize      = flag.Uint64("batch", 2000, "blocks per eth_getLogs batch")
		dryRun         = flag.Bool("dry-run", false, "index events and compute miles but skip Fuel submission and processed marking")
		keystorePath   = flag.String("keystore", "", "Path to the executor keystore file (required for production token sweeping)")
		passphrase     = flag.String("passphrase", "", "Password for the keystore")
		executorFlag   = flag.String("executor", "", "Executor public address for dry-run sweeping simulation")
		fastswapURL    = flag.String("fastswap-url", "", "FastSwap API endpoint URL (e.g., https://fastrpc.primev.xyz)")
		fundsRecipient = flag.String("funds-recipient", defaultRecipient, "Address to receive swept ETH")
		maxGasGwei     = flag.Uint64("max-gas-gwei", 50, "Skip sweep if L1 gas price exceeds this (gwei)")

		// Test-swap mode: run a single FastSwap sweep and exit.
		testSwap        = flag.Bool("test-swap", false, "Run a single test swap and exit")
		testInputToken  = flag.String("test-input-token", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", "Input token for test swap (default: USDC)")
		testInputAmount = flag.String("test-input-amount", "1000000", "Input amount for test swap in raw units (default: 1 USDC = 1000000)")

		dbUser    = flag.String("db-user", "", "Database user")
		dbPass    = flag.String("db-pw", "", "Database password")
		dbHost    = flag.String("db-host", "127.0.0.1", "Database host")
		dbPort    = flag.String("db-port", "9030", "Database port")
		dbName    = flag.String("db-name", "mevcommit_57173", "Database name")
		l1RPC     = flag.String("l1-rpc-url", "", "L1 Ethereum HTTP RPC URL")
		barterURL = flag.String("barter-url", "", "Barter API base URL")
		barterKey = flag.String("barter-api-key", "", "Barter API Key")
		fuelURL   = flag.String("fuel-api-url", "", "Fuel points API URL")
		fuelKey   = flag.String("fuel-api-key", "", "Fuel points API Key")
	)
	flag.Parse()

	// Required flags validation (relaxed for test-swap mode)
	if *testSwap {
		if *keystorePath == "" || *passphrase == "" {
			log.Fatal("-keystore and -passphrase are required for test-swap")
		}
		if *fastswapURL == "" {
			log.Fatal("-fastswap-url is required for test-swap")
		}
		if *l1RPC == "" || *barterURL == "" {
			log.Fatal("-l1-rpc-url and -barter-url are required for test-swap")
		}
		runTestSwap(*keystorePath, *passphrase, *fastswapURL, *l1RPC, *barterURL, *barterKey, *fundsRecipient, *contractAddr, *testInputToken, *testInputAmount, *maxGasGwei)
		return
	}

	// Regular mode: required flags validation
	if *l1RPC == "" || *barterURL == "" || *dbUser == "" || *dbPass == "" || *dbHost == "" {
		log.Fatal("missing required flags: -l1-rpc-url, -barter-url, -db-user, -db-pw, -db-host")
	}

	if *dryRun {
		log.Println("DRY-RUN mode: will index events and compute miles but skip Fuel submission")
	} else {
		if *fuelURL == "" || *fuelKey == "" {
			log.Fatal("-fuel-api-url and -fuel-api-key are required in production mode")
		}
	}

	var executorAddr common.Address
	var signer keysigner.KeySigner
	var err error

	if *dryRun {
		if *executorFlag == "" {
			log.Fatal("-executor is required for dry-run to simulate token sweeping")
		}
		executorAddr = common.HexToAddress(*executorFlag)
		log.Printf("Using executor address %s for dry-run simulation", executorAddr.Hex())
	} else {
		if *keystorePath == "" || *passphrase == "" {
			log.Fatal("-keystore and -passphrase are required for production token sweeping")
		}
		signer, err = loadKeystoreFile(*keystorePath, *passphrase)
		if err != nil {
			log.Fatalf("failed to load keystore: %v", err)
		}
		executorAddr = signer.GetAddress()
		log.Printf("Loaded executor wallet %s for production sweeping", executorAddr.Hex())
	}

	db := openDB(*dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	defer func() { _ = db.Close() }()

	client, err := ethclient.Dial(*l1RPC)
	if err != nil {
		log.Fatalf("ethclient.Dial: %v", err)
	}

	contract := common.HexToAddress(*contractAddr)
	weth := common.HexToAddress(*wethAddr)

	filterer, err := fastsettlement.NewFastsettlementv3Filterer(contract, client)
	if err != nil {
		log.Fatalf("NewFastsettlementv3Filterer: %v", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			MaxIdleConns:        64,
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     90 * time.Second,
			ForceAttemptHTTP2:   true,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: 15 * time.Second,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Determine start block.
	from := *startBlock
	if from == 0 {
		from = loadLastBlock(db)
	}
	if from == 0 {
		log.Fatal("no start block: set -start-block or ensure fastswap_miles_meta has a saved cursor")
	}

	log.Printf("starting: contract=%s weth=%s from_block=%d poll=%s batch=%d dry_run=%v",
		contract.Hex(), weth.Hex(), from, *pollInterval, *batchSize, *dryRun)

	ticker := time.NewTicker(*pollInterval)
	defer ticker.Stop()

	// Trigger the first iteration immediately
	firstRun := make(chan struct{}, 1)
	firstRun <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-firstRun:
		case <-ticker.C:
		}

		head, err := client.BlockNumber(ctx)
		if err != nil {
			log.Printf("BlockNumber: %v", err)
			continue
		}
		if from > head {
			continue
		}

		to := from + *batchSize - 1
		if to > head {
			to = head
		}

		indexed, err := indexBatch(ctx, filterer, client, db, weth, from, to)
		if err != nil {
			log.Printf("indexBatch [%d..%d]: %v", from, to, err)
			continue
		}

		// Only process miles (which involves expensive Barter API calls) if we've caught up to the chain tip
		var processed int
		if to == head {
			var errProcess error
			processed, errProcess = processMiles(ctx, db, weth, *fuelURL, *fuelKey, *barterURL, *barterKey, client, client, signer, executorAddr, httpClient, *dryRun, *fastswapURL, common.HexToAddress(*fundsRecipient), contract, *maxGasGwei)
			if errProcess != nil {
				log.Printf("processMiles: %v", errProcess)
				// Don't stop advancing — events are already persisted.
			}
		}

		saveLastBlock(db, to+1)
		if to == head {
			log.Printf("batch [%d..%d]: indexed=%d processed=%d next=%d", from, to, indexed, processed, to+1)
		} else {
			log.Printf("batch [%d..%d]: indexed=%d caught_up=false next=%d", from, to, indexed, to+1)
		}
		from = to + 1
	}
}

// -------------------- DB --------------------

func openDB(dbUser, dbPass, dbHost, dbPort, dbName string) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true&parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql.Open: %v", err)
	}
	db.SetMaxOpenConns(6)
	db.SetMaxIdleConns(6)
	db.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db.Ping: %v", err)
	}
	return db
}

func loadLastBlock(db *sql.DB) uint64 {
	var v uint64
	err := db.QueryRow(`SELECT CAST(v AS BIGINT) FROM mevcommit_57173.fastswap_miles_meta WHERE k = 'last_block'`).Scan(&v)
	if err != nil {
		return 0
	}
	return v
}

func saveLastBlock(db *sql.DB, block uint64) {
	// StarRocks doesn't support ON CONFLICT, so we delete + insert (single-row meta table).
	_, _ = db.Exec(`DELETE FROM mevcommit_57173.fastswap_miles_meta WHERE k = 'last_block'`)
	_, err := db.Exec(`INSERT INTO mevcommit_57173.fastswap_miles_meta (k, v) VALUES ('last_block', ?)`, fmt.Sprintf("%d", block))
	if err != nil {
		log.Printf("saveLastBlock: %v", err)
	}
}

// -------------------- Event Indexer --------------------

func indexBatch(
	ctx context.Context,
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

		// Get tx receipt for gas cost.
		receipt, err := client.TransactionReceipt(ctx, ev.Raw.TxHash)
		if err != nil {
			log.Printf("receipt %s: %v (skipping gas cost)", txHash, err)
		}
		var gasCost *big.Int
		if receipt != nil {
			gasCost = new(big.Int).Mul(
				new(big.Int).SetUint64(receipt.GasUsed),
				receipt.EffectiveGasPrice,
			)
		}

		// Get block timestamp.
		header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
		if err != nil {
			log.Printf("header %d: %v", blockNum, err)
		}
		var blockTS *time.Time
		if header != nil {
			t := time.Unix(int64(header.Time), 0).UTC()
			blockTS = &t
		}

		// Determine swap_type based on output token.
		swapType := "erc20"
		if ev.OutputToken == zeroAddr || strings.EqualFold(ev.OutputToken.Hex(), weth.Hex()) {
			swapType = "eth_weth"
		}

		err = insertEvent(db, txHash, blockNum, blockTS, ev, gasCost, swapType)
		if err != nil {
			log.Printf("insertEvent %s: %v", txHash, err)
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
	var tsVal interface{} = nil
	if blockTS != nil {
		tsVal = *blockTS
	}
	var gcStr interface{} = nil
	if gasCost != nil {
		gcStr = gasCost.String()
	}

	_, err := db.Exec(`
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

// -------------------- Barter API Types --------------------

type BarterResponse struct {
	To        common.Address `json:"to"`
	GasLimit  string         `json:"gasLimit"`
	Value     string         `json:"value"`
	Data      string         `json:"data"`
	MinReturn string         `json:"minReturn"`
	Route     struct {
		OutputAmount  string `json:"outputAmount"`
		GasEstimation uint64 `json:"gasEstimation"`
		BlockNumber   uint64 `json:"blockNumber"`
	} `json:"route"`
}

type barterRequest struct {
	Source            string  `json:"source"`
	Target            string  `json:"target"`
	SellAmount        string  `json:"sellAmount"`
	Recipient         string  `json:"recipient"`
	Origin            string  `json:"origin"`
	MinReturnFraction float64 `json:"minReturnFraction"` // 0.98 for 2% slippage
	Deadline          string  `json:"deadline"`
}

func callBarter(
	ctx context.Context,
	httpClient *http.Client,
	apiURL, apiKey string,
	reqBody barterRequest,
) (*BarterResponse, error) {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal barter request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL+"/swap", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("barter API request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read barter response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("barter API error %d: %s", resp.StatusCode, string(respBody))
	}

	var barterResp BarterResponse
	if err := json.Unmarshal(respBody, &barterResp); err != nil {
		return nil, fmt.Errorf("decode barter response: %w", err)
	}
	return &barterResp, nil
}

// -------------------- Miles Processor --------------------

func processMiles(
	ctx context.Context,
	db *sql.DB,
	weth common.Address,
	fuelURL, fuelKey string,
	barterURL, barterKey string,
	client *ethclient.Client,
	l1Client *ethclient.Client,
	signer keysigner.KeySigner,
	executorAddr common.Address,
	httpClient *http.Client,
	dryRun bool,
	fastswapURL string,
	fundsRecipient common.Address,
	settlementAddr common.Address,
	maxGasGwei uint64,
) (int, error) {
	// Process ETH/WETH rows only (phase 1). ERC20 deferred.
	rows, err := db.QueryContext(ctx, `
SELECT tx_hash, user_address, surplus, gas_cost, input_token, block_timestamp
FROM mevcommit_57173.fastswap_miles
WHERE processed = false
  AND swap_type = 'eth_weth'
  AND LOWER(user_address) != LOWER(?)
`, executorAddr.Hex())
	if err != nil {
		return 0, fmt.Errorf("query unprocessed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	type row struct {
		txHash     string
		user       string
		surplus    string
		gasCost    sql.NullString
		inputToken string
		blockTS    time.Time
	}
	var pending []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.txHash, &r.user, &r.surplus, &r.gasCost, &r.inputToken, &r.blockTS); err != nil {
			return 0, err
		}
		pending = append(pending, r)
	}
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	// Batch-fetch all bid costs and FastRPC status in one query each.
	allHashes := make([]string, len(pending))
	for i, r := range pending {
		allHashes[i] = r.txHash
	}
	bidMap := batchLookupBidCosts(db, allHashes)
	fastRPCSet := batchCheckFastRPC(db, allHashes)

	processed := 0
	for _, r := range pending {
		surplusWei, ok := new(big.Int).SetString(r.surplus, 10)
		if !ok {
			log.Printf("bad surplus %s for tx %s", r.surplus, r.txHash)
			continue
		}

		// When input is ETH, the user pays gas — don't deduct from our profit.
		// When input is a token, we submit the tx and pay gas.
		userPaysGas := strings.EqualFold(r.inputToken, zeroAddr.Hex())

		gasCostWei := big.NewInt(0)
		if !userPaysGas && r.gasCost.Valid && r.gasCost.String != "" {
			if gc, ok := new(big.Int).SetString(r.gasCost.String, 10); ok {
				gasCostWei = gc
			}
		}

		// Lookup bid cost from pre-fetched map.
		bidCostWei := getBidCost(bidMap, r.txHash)

		// If no bid found, check mctransactions_sr to determine if tx went through FastRPC.
		if bidCostWei.Sign() == 0 {
			if fastRPCSet[strings.ToLower(r.txHash)] {
				// Tx IS in FastRPC but bid not indexed yet in tx_view. Retry later.
				log.Printf("[dry-run=%v] tx in FastRPC but bid not indexed yet tx=%s user=%s (will retry)",
					dryRun, r.txHash, r.user)
				continue
			}
			// Tx NOT in FastRPC — didn't use our RPC.
			log.Printf("[dry-run=%v] tx not in FastRPC, skipping tx=%s user=%s (0 miles)",
				dryRun, r.txHash, r.user)
			markProcessed(db, r.txHash, weiToEth(surplusWei), 0, 0, "0")
			processed++
			continue
		}

		// Net profit = surplus - gas_cost (if we pay) - bid_cost
		netProfit := new(big.Int).Sub(surplusWei, gasCostWei)
		netProfit.Sub(netProfit, bidCostWei)

		// Convert values for logging.
		surplusEth := weiToEth(surplusWei)
		netProfitEth := weiToEth(netProfit)

		if netProfit.Sign() <= 0 {
			log.Printf("[dry-run=%v] no profit tx=%s user=%s surplus_eth=%.6f net_profit_eth=%.6f gas=%s bid=%s",
				dryRun, r.txHash, r.user, surplusEth, netProfitEth, gasCostWei.String(), bidCostWei.String())
			// Even in dry-run, mark unprofitable trades as processed to avoid infinite log loops
			markProcessed(db, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
			processed++
			continue
		}

		// Award 90% of net profit as miles basis.
		userShare := new(big.Int).Mul(netProfit, big.NewInt(90))
		userShare.Div(userShare, big.NewInt(100))

		miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
		if miles.Sign() <= 0 {
			log.Printf("[dry-run=%v] sub-threshold tx=%s user=%s surplus_eth=%.6f net_profit_eth=%.6f",
				dryRun, r.txHash, r.user, surplusEth, netProfitEth)
			// Even in dry-run, mark unprofitable trades as processed to avoid infinite log loops
			markProcessed(db, r.txHash, surplusEth, netProfitEth, 0, bidCostWei.String())
			processed++
			continue
		}

		log.Printf("[dry-run=%v] miles=%d user=%s tx=%s surplus_eth=%.6f net_profit_eth=%.6f gas=%s bid=%s",
			dryRun, miles.Int64(), r.user, r.txHash, surplusEth, netProfitEth,
			gasCostWei.String(), bidCostWei.String())

		if dryRun {
			// Actually mark it processed in dry-run so we don't spam the console in a loop
			markProcessed(db, r.txHash, surplusEth, netProfitEth, miles.Int64(), bidCostWei.String())
			processed++
			continue
		}

		// Submit to Fuel API.
		err := submitToFuel(ctx, httpClient, fuelURL, fuelKey,
			common.HexToAddress(r.user),
			common.HexToHash(r.txHash),
			miles,
		)
		if err != nil {
			log.Printf("fuel submit %s: %v", r.txHash, err)
			continue // Retry next cycle.
		}

		markProcessed(db, r.txHash, surplusEth, netProfitEth, miles.Int64(), bidCostWei.String())
		processed++
		log.Printf("awarded %d miles to %s (tx=%s surplus_eth=%.6f net_profit_eth=%.6f)",
			miles.Int64(), r.user, r.txHash, surplusEth, netProfitEth)
	}

	// Process ERC20 rows (phase 2).
	erc20Processed, err := processERC20Miles(ctx, db, weth, fuelURL, fuelKey, barterURL, barterKey, client, l1Client, signer, executorAddr, httpClient, dryRun, fastswapURL, fundsRecipient, settlementAddr, maxGasGwei)
	if err != nil {
		log.Printf("processERC20Miles: %v", err)
	}
	processed += erc20Processed

	return processed, nil
}

type erc20Row struct {
	txHash     string
	user       string
	token      string
	surplus    string
	gasCost    sql.NullString
	inputToken string
	blockTS    time.Time
}

type tokenBatch struct {
	Token    string
	TotalSum *big.Int
	Txs      []erc20Row
}

// Minimal ERC20 ABI for Approve
const erc20ApproveABI = `[{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

func processERC20Miles(
	ctx context.Context,
	db *sql.DB,
	weth common.Address,
	fuelURL, fuelKey string,
	barterURL, barterKey string,
	client *ethclient.Client,
	l1Client *ethclient.Client,
	signer keysigner.KeySigner,
	executorAddr common.Address,
	httpClient *http.Client,
	dryRun bool,
	fastswapURL string,
	fundsRecipient common.Address,
	settlementAddr common.Address,
	maxGasGwei uint64,
) (int, error) {
	processed := 0

	rows, err := db.QueryContext(ctx, `
SELECT tx_hash, user_address, output_token, surplus, gas_cost, input_token, block_timestamp
FROM mevcommit_57173.fastswap_miles
WHERE processed = false
  AND swap_type = 'erc20'
  AND LOWER(user_address) != LOWER(?)
`, executorAddr.Hex())
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
			log.Printf("bad surplus %s for tx %s", r.surplus, r.txHash)
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

	// Batch-fetch all bid costs and FastRPC status for all pending erc20 tx hashes.
	allErc20Hashes := make([]string, len(pending))
	for i, r := range pending {
		allErc20Hashes[i] = r.txHash
	}
	erc20BidMap := batchLookupBidCosts(db, allErc20Hashes)
	erc20FastRPCSet := batchCheckFastRPC(db, allErc20Hashes)

	for token, batch := range batches {
		gasCosts := make([]*big.Int, len(batch.Txs))
		bidCosts := make([]*big.Int, len(batch.Txs))
		totalOriginalGasCost := big.NewInt(0)
		totalOriginalBidCost := big.NewInt(0)

		var skippedRows []erc20Row // Track rows skipped due to no bid

		for i, r := range batch.Txs {
			// Lookup bid cost from pre-fetched map — if no bid, check FastRPC.
			bidCostWei := getBidCost(erc20BidMap, r.txHash)
			if bidCostWei.Sign() == 0 {
				if erc20FastRPCSet[strings.ToLower(r.txHash)] {
					// Tx IS in FastRPC but bid not indexed yet. Retry later.
					log.Printf("[dry-run=%v] erc20 tx in FastRPC but bid not indexed yet tx=%s user=%s (will retry)",
						dryRun, r.txHash, r.user)
				} else {
					// Tx NOT in FastRPC — didn't use our RPC.
					log.Printf("[dry-run=%v] erc20 tx not in FastRPC, skipping tx=%s user=%s (0 miles)",
						dryRun, r.txHash, r.user)
					skippedRows = append(skippedRows, r)
				}
				gasCosts[i] = big.NewInt(0)
				bidCosts[i] = big.NewInt(0)
				continue
			}

			// When input is ETH, the user pays gas — don't deduct from our profit.
			userPaysGas := strings.EqualFold(r.inputToken, zeroAddr.Hex())

			gasCostWei := big.NewInt(0)
			if !userPaysGas && r.gasCost.Valid && r.gasCost.String != "" {
				if gc, ok := new(big.Int).SetString(r.gasCost.String, 10); ok {
					gasCostWei = gc
				}
			}

			gasCosts[i] = gasCostWei
			bidCosts[i] = bidCostWei

			totalOriginalGasCost.Add(totalOriginalGasCost, gasCostWei)
			totalOriginalBidCost.Add(totalOriginalBidCost, bidCostWei)
		}

		// Mark skipped rows as processed with 0 miles.
		for _, r := range skippedRows {
			surplusWei, _ := new(big.Int).SetString(r.surplus, 10)
			markProcessed(db, r.txHash, weiToEth(surplusWei), 0, 0, "0")
			processed++
		}

		if dryRun && len(batch.Txs) > 0 {
			firstTx := batch.Txs[0]
			log.Printf("[dry-run=true] token %s breakdown for sample tx %s: gas_cost_eth=%.6f bid_cost_eth=%.6f",
				token, firstTx.txHash, weiToEth(gasCosts[0]), weiToEth(bidCosts[0]))
		}

		// Price checking and sweeping logic here
		reqBody := barterRequest{
			Source:            token,
			Target:            weth.Hex(),
			SellAmount:        batch.TotalSum.String(),
			Recipient:         executorAddr.Hex(), // Return to executor address
			Origin:            executorAddr.Hex(),
			MinReturnFraction: 0.98,
			Deadline:          fmt.Sprintf("%d", time.Now().Add(10*time.Minute).Unix()),
		}

		barterResp, err := callBarter(ctx, httpClient, barterURL, barterKey, reqBody)
		if err != nil {
			log.Printf("callBarter token %s: %v", token, err)
			continue
		}

		// Estimate Gas
		gasLimit, err := strconv.ParseUint(barterResp.GasLimit, 10, 64)
		if err != nil {
			log.Printf("invalid gasLimit %s from barter", barterResp.GasLimit)
			continue
		}
		gasLimit += 50000 // Add safety buffer

		gasPrice, err := client.SuggestGasPrice(ctx)
		if err != nil {
			log.Printf("suggest gas price: %v", err)
			continue
		}

		expectedGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)
		expectedEthReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
		if !ok {
			log.Printf("invalid MinReturn %s from barter", barterResp.MinReturn)
			continue
		}

		// Profitability: ETH Return > Swap Gas Cost + Total Underlying Bid Cost + Total Underlying Gas Cost
		totalSweepCosts := new(big.Int).Add(expectedGasCost, totalOriginalBidCost)
		totalSweepCosts.Add(totalSweepCosts, totalOriginalGasCost)

		if expectedEthReturn.Cmp(totalSweepCosts) <= 0 {
			log.Printf("[dry-run=%v] token %s sweeping not yet profitable (would return %.6f ETH vs total %.6f ETH) [breakdown: sweep_gas=%.6f, orig_tx_gas=%.6f, orig_tx_bids=%.6f]",
				dryRun, token, weiToEth(expectedEthReturn), weiToEth(totalSweepCosts),
				weiToEth(expectedGasCost), weiToEth(totalOriginalGasCost), weiToEth(totalOriginalBidCost))
			continue
		}

		var actualEthReturn *big.Int
		var actualSwapGasCost *big.Int

		if dryRun {
			log.Printf("[dry-run=true] simulated sweep of %s token %s returning %.6f ETH (gas=%.6f)",
				batch.TotalSum.String(), token, weiToEth(expectedEthReturn), weiToEth(expectedGasCost))
			actualEthReturn = expectedEthReturn
			actualSwapGasCost = expectedGasCost
		} else {
			actualEthReturn, actualSwapGasCost, err = submitFastSwapSweep(ctx, client, l1Client, httpClient, signer, executorAddr, common.HexToAddress(token), batch.TotalSum, fastswapURL, fundsRecipient, settlementAddr, barterResp, maxGasGwei)
			if err != nil {
				log.Printf("failed to sweep token %s: %v", token, err)
				continue
			}
			log.Printf("FastSwap sweep token %s success: expected %.6f ETH (est gas=%.6f)", token, weiToEth(actualEthReturn), weiToEth(actualSwapGasCost))
		}

		// Proportionally award miles
		for i, r := range batch.Txs {
			surplusWei, _ := new(big.Int).SetString(r.surplus, 10)

			// Share of gross ETH: actualEthReturn * (surplusWei / TotalSum)
			txGrossEth := new(big.Int).Mul(actualEthReturn, surplusWei)
			txGrossEth.Div(txGrossEth, batch.TotalSum)

			txOverheadGas := new(big.Int).Mul(actualSwapGasCost, surplusWei)
			txOverheadGas.Div(txOverheadGas, batch.TotalSum)

			txNetProfit := new(big.Int).Sub(txGrossEth, gasCosts[i])
			txNetProfit.Sub(txNetProfit, bidCosts[i])
			txNetProfit.Sub(txNetProfit, txOverheadGas)

			surplusEth := weiToEth(txGrossEth)
			netProfitEth := weiToEth(txNetProfit)

			if txNetProfit.Sign() <= 0 {
				log.Printf("[dry-run=%v] no profit subset tx=%s user=%s gross_eth=%.6f net_profit_eth=%.6f orig_gas=%.6f orig_bid=%.6f overhead_gas=%.6f",
					dryRun, r.txHash, r.user, surplusEth, netProfitEth, weiToEth(gasCosts[i]), weiToEth(bidCosts[i]), weiToEth(txOverheadGas))
				if !dryRun {
					markProcessed(db, r.txHash, surplusEth, netProfitEth, 0, bidCosts[i].String())
				}
				processed++
				continue
			}

			userShare := new(big.Int).Mul(txNetProfit, big.NewInt(90))
			userShare.Div(userShare, big.NewInt(100))

			miles := new(big.Int).Div(userShare, big.NewInt(weiPerPoint))
			if miles.Sign() <= 0 {
				log.Printf("[dry-run=%v] sub-threshold subset tx=%s user=%s gross_eth=%.6f net_profit_eth=%.6f",
					dryRun, r.txHash, r.user, surplusEth, netProfitEth)
				if !dryRun {
					markProcessed(db, r.txHash, surplusEth, netProfitEth, 0, bidCosts[i].String())
				}
				processed++
				continue
			}

			log.Printf("[dry-run=%v] miles=%d subset user=%s tx=%s gross_eth=%.6f net_profit_eth=%.6f orig_gas=%.6f orig_bid=%.6f overhead_gas=%.6f",
				dryRun, miles.Int64(), r.user, r.txHash, surplusEth, netProfitEth,
				weiToEth(gasCosts[i]), weiToEth(bidCosts[i]), weiToEth(txOverheadGas))

			if !dryRun {
				err := submitToFuel(ctx, httpClient, fuelURL, fuelKey,
					common.HexToAddress(r.user),
					common.HexToHash(r.txHash),
					miles,
				)
				if err != nil {
					log.Printf("fuel submit %s failed (will not retry to avoid duplicate swap): %v", r.txHash, err)
				}
				markProcessed(db, r.txHash, surplusEth, netProfitEth, miles.Int64(), bidCosts[i].String())
			}
			processed++
		}
	}

	return processed, nil
}

func submitFastSwapSweep(
	ctx context.Context,
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

	// Gas price gate: skip if L1 gas is too expensive.
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("gas price: %w", err)
	}
	gasPriceGwei := new(big.Int).Div(gasPrice, big.NewInt(1_000_000_000))
	if gasPriceGwei.Uint64() > maxGasGwei {
		return nil, nil, fmt.Errorf("gas price %d gwei exceeds max %d gwei, skipping sweep", gasPriceGwei.Uint64(), maxGasGwei)
	}

	// Ensure the executor has approved this token to the Permit2 contract.
	// Uses l1Client for sending the approve tx (no mev-commit balance needed).
	if err := ensurePermit2Approval(ctx, client, l1Client, signer, executorAddr, tokenAddr, totalAmount); err != nil {
		return nil, nil, fmt.Errorf("permit2 approval: %w", err)
	}

	// Use 95% of Barter's minReturn as userAmtOut. This gives the handler's
	// internal Barter call room to succeed (prices may shift between our quote
	// and when the handler re-quotes). Barter already applies 2% slippage.
	barterMinReturn, ok := new(big.Int).SetString(barterResp.MinReturn, 10)
	if !ok {
		return nil, nil, fmt.Errorf("invalid MinReturn from barter: %s", barterResp.MinReturn)
	}
	userAmtOut := new(big.Int).Mul(barterMinReturn, big.NewInt(95))
	userAmtOut.Div(userAmtOut, big.NewInt(100))

	// Deadline: 10 minutes from now.
	deadline := big.NewInt(time.Now().Add(10 * time.Minute).Unix())

	// Random Permit2 nonce (Permit2 uses unordered nonces via bitmap).
	nonceBuf := make([]byte, 32)
	if _, err := rand.Read(nonceBuf); err != nil {
		return nil, nil, fmt.Errorf("generate nonce: %w", err)
	}
	nonce := new(big.Int).SetBytes(nonceBuf)

	// Sign the Permit2 witness (EIP-712 typed data).
	signature, err := signPermit2Witness(
		signer,
		tokenAddr,
		totalAmount,
		settlementAddr,
		nonce,
		deadline,
		executorAddr,     // user = executor (holds the tokens)
		tokenAddr,        // inputToken
		common.Address{}, // outputToken = ETH (address(0))
		totalAmount,      // inputAmt
		userAmtOut,       // userAmtOut
		fundsRecipient,   // recipient for output ETH
		deadline,         // intent deadline
		nonce,            // intent nonce
	)
	if err != nil {
		return nil, nil, fmt.Errorf("sign permit2: %w", err)
	}

	// Build the SwapRequest and POST to /fastswap.
	swapReq := map[string]string{
		"user":        executorAddr.Hex(),
		"inputToken":  tokenAddr.Hex(),
		"outputToken": common.Address{}.Hex(), // ETH
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

	log.Printf("FastSwap sweep submitted: tx=%s outputAmount=%s gasLimit=%d", result.TxHash, result.OutputAmount, result.GasLimit)

	// Use the expected values from Barter (actual settlement happens on-chain via the executor).
	// The gas cost estimation uses Barter's gas limit * current gas price.
	gasLimit, _ := strconv.ParseUint(barterResp.GasLimit, 10, 64)
	gasLimit += 100000 // settlement overhead
	estimatedGasCost := new(big.Int).Mul(big.NewInt(int64(gasLimit)), gasPrice)

	return userAmtOut, estimatedGasCost, nil
}

// ensurePermit2Approval checks if the executor has sufficient ERC20 allowance
// to the Permit2 contract. If not, sends a max-uint256 approval transaction.
func ensurePermit2Approval(
	ctx context.Context,
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

	// Check current allowance.
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
		log.Printf("Permit2 allowance sufficient for %s (have %s, need %s)", token.Hex(), currentAllowance.String(), requiredAmount.String())
		return nil
	}

	log.Printf("Permit2 allowance insufficient for %s (have %s, need %s), approving max...",
		token.Hex(), currentAllowance.String(), requiredAmount.String())

	// Approve max uint256 to Permit2.
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
		Gas:      100000, // approve is cheap
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
	log.Printf("Permit2 approve tx sent: %s", signedTx.Hash().Hex())

	// Wait for receipt with 15 minute timeout.
	deadline := time.Now().Add(15 * time.Minute)
	for time.Now().Before(deadline) {
		time.Sleep(12 * time.Second)
		receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
		if err != nil {
			continue
		}
		if receipt.Status != 1 {
			return fmt.Errorf("approve tx reverted: %s", signedTx.Hash().Hex())
		}
		log.Printf("Permit2 approve confirmed: %s", signedTx.Hash().Hex())
		return nil
	}
	return fmt.Errorf("approve tx not confirmed after 15 min, will retry next cycle: %s", signedTx.Hash().Hex())
}

// signPermit2Witness signs the Permit2 PermitWitnessTransferFrom EIP-712 typed data.
// This authorizes the FastSettlement contract (spender) to pull tokens from the signer
// via Permit2, with the Intent as the witness data.
func signPermit2Witness(
	signer keysigner.KeySigner,
	token common.Address,
	amount *big.Int,
	spender common.Address,
	permitNonce *big.Int,
	permitDeadline *big.Int,
	// Intent witness fields:
	user common.Address,
	inputToken common.Address,
	outputToken common.Address,
	inputAmt *big.Int,
	userAmtOut *big.Int,
	recipient common.Address,
	intentDeadline *big.Int,
	intentNonce *big.Int,
) ([]byte, error) {
	// Permit2 domain
	permit2 := common.HexToAddress(permit2Addr)

	// EIP-712 domain separator
	domainSep := crypto.Keccak256(
		crypto.Keccak256([]byte("EIP712Domain(string name,uint256 chainId,address verifyingContract)")),
		crypto.Keccak256([]byte("Permit2")),
		padTo32(big.NewInt(1)), // chainId = 1 (mainnet)
		padTo32Address(permit2),
	)

	// TokenPermissions type hash
	tokenPermissionsTypeHash := crypto.Keccak256([]byte("TokenPermissions(address token,uint256 amount)"))
	tokenPermissionsHash := crypto.Keccak256(
		tokenPermissionsTypeHash,
		padTo32Address(token),
		padTo32(amount),
	)

	// Intent (witness) type hash
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

	// PermitWitnessTransferFrom type hash
	permitWitnessTypeHash := crypto.Keccak256([]byte(
		"PermitWitnessTransferFrom(TokenPermissions permitted,address spender,uint256 nonce,uint256 deadline,Intent witness)" +
			"Intent(address user,address inputToken,address outputToken,uint256 inputAmt,uint256 userAmtOut,address recipient,uint256 deadline,uint256 nonce)" +
			"TokenPermissions(address token,uint256 amount)",
	))

	// Struct hash
	structHash := crypto.Keccak256(
		permitWitnessTypeHash,
		tokenPermissionsHash,
		padTo32Address(spender),
		padTo32(permitNonce),
		padTo32(permitDeadline),
		witnessHash,
	)

	// Final EIP-712 hash: \x19\x01 + domainSep + structHash
	digest := crypto.Keccak256(
		[]byte{0x19, 0x01},
		domainSep,
		structHash,
	)

	// Sign with the keysigner
	sig, err := signer.SignHash(digest)
	if err != nil {
		return nil, fmt.Errorf("sign hash: %w", err)
	}

	// crypto.Sign returns v=0/1, but Permit2/ecrecover expects v=27/28.
	if len(sig) == 65 && sig[64] < 27 {
		sig[64] += 27
	}

	return sig, nil
}

// padTo32 pads a big.Int to 32 bytes (left-padded with zeros).
func padTo32(n *big.Int) []byte {
	b := n.Bytes()
	if len(b) >= 32 {
		return b[:32]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

// padTo32Address pads an address to 32 bytes (left-padded with zeros).
func padTo32Address(addr common.Address) []byte {
	padded := make([]byte, 32)
	copy(padded[12:], addr.Bytes())
	return padded
}

// batchLookupBidCosts queries tx_view once for all OpenedCommitmentStored events
// matching the given L1 tx hashes and returns a map of txHash -> bidAmt.
// This avoids per-tx 13M+ row scans by using a single IN-clause query.
func batchLookupBidCosts(db *sql.DB, txHashes []string) map[string]*big.Int {
	result := make(map[string]*big.Int, len(txHashes))
	if len(txHashes) == 0 {
		return result
	}

	// Build normalized hash list and IN-clause placeholders.
	normalized := make([]string, len(txHashes))
	for i, h := range txHashes {
		normalized[i] = strings.TrimPrefix(strings.ToLower(h), "0x")
	}

	// StarRocks doesn't support parameterized IN with variable-length lists well,
	// so build the IN clause directly with quoted values.
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
  ) IN (%s)`, inClause.String())

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("batchLookupBidCosts query error: %v", err)
		return result
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var txnHash string
		var bidAmtStr sql.NullString
		if err := rows.Scan(&txnHash, &bidAmtStr); err != nil {
			log.Printf("batchLookupBidCosts scan error: %v", err)
			continue
		}
		if !bidAmtStr.Valid || bidAmtStr.String == "" || bidAmtStr.String == "null" {
			continue
		}
		cleanStr := strings.Trim(bidAmtStr.String, "\"")
		v, ok := new(big.Int).SetString(cleanStr, 10)
		if !ok {
			v, ok = new(big.Int).SetString(strings.TrimPrefix(cleanStr, "0x"), 16)
			if !ok {
				log.Printf("batchLookupBidCosts parse error for %s, value: %q", txnHash, bidAmtStr.String)
				continue
			}
		}
		// Map back using the normalized hash (without 0x prefix).
		result[txnHash] = v
	}

	log.Printf("batchLookupBidCosts: found %d/%d bid costs", len(result), len(txHashes))
	return result
}

// getBidCost retrieves a bid cost from the pre-fetched map. Returns 0 if not found.
func getBidCost(bidMap map[string]*big.Int, txHash string) *big.Int {
	hashNorm := strings.TrimPrefix(strings.ToLower(txHash), "0x")
	if v, ok := bidMap[hashNorm]; ok {
		return v
	}
	return big.NewInt(0)
}

// batchCheckFastRPC queries mctransactions_sr to determine which tx hashes were
// submitted through our FastRPC. Returns a set of lowercase tx hashes that exist.
func batchCheckFastRPC(db *sql.DB, txHashes []string) map[string]bool {
	result := make(map[string]bool, len(txHashes))
	if len(txHashes) == 0 {
		return result
	}

	// Build IN clause with the hashes as-is (they include 0x prefix in mctransactions_sr).
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

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("batchCheckFastRPC query error: %v", err)
		return result
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var h string
		if err := rows.Scan(&h); err != nil {
			log.Printf("batchCheckFastRPC scan error: %v", err)
			continue
		}
		result[strings.ToLower(h)] = true
	}

	log.Printf("batchCheckFastRPC: %d/%d txns found in FastRPC", len(result), len(txHashes))
	return result
}

func markProcessed(db *sql.DB, txHash string, surplusEth, netProfitEth float64, miles int64, bidCost string) {
	_, err := db.Exec(`
UPDATE mevcommit_57173.fastswap_miles
SET surplus_eth = ?, net_profit_eth = ?, miles = ?, bid_cost = ?, processed = true
WHERE tx_hash = ?`,
		surplusEth, netProfitEth, miles, bidCost, txHash)
	if err != nil {
		log.Printf("markProcessed %s: %v", txHash, err)
	}
}

// -------------------- Fuel API --------------------

func submitToFuel(
	ctx context.Context,
	client *http.Client,
	apiURL, apiKey string,
	user common.Address,
	txHash common.Hash,
	miles *big.Int,
) error {
	body := map[string]any{
		"user": map[string]any{
			"identifier_type": "evm_address",
			"identifier":      user.Hex(),
		},
		"name": "fast-swap-surplus",
		"args": map[string]any{
			"value": map[string]any{
				"amount": miles.String(),
				"currency": map[string]any{
					"name": "POINT",
				},
			},
			"transaction_hash": txHash.Hex(),
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("fuel API status %d", resp.StatusCode)
	}
	return nil
}

// -------------------- Helpers --------------------

// loadKeystoreFile reads a single keystore JSON file, decrypts it with the
// passphrase, and returns a PrivateKeySigner. Accepts a full file path.
func loadKeystoreFile(path, passphrase string) (keysigner.KeySigner, error) {
	keyjson, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read keystore file: %w", err)
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, fmt.Errorf("decrypt keystore: %w", err)
	}

	hexKey := fmt.Sprintf("%x", crypto.FromECDSA(key.PrivateKey))
	signer, err := keysigner.NewPrivateKeySignerFromHex(hexKey)
	if err != nil {
		return nil, fmt.Errorf("create signer: %w", err)
	}

	// Zero the sensitive bytes.
	for i := range keyjson {
		keyjson[i] = 0
	}

	return signer, nil
}

func weiToEth(wei *big.Int) float64 {
	f := new(big.Float).SetInt(wei)
	eth, _ := f.Quo(f, new(big.Float).SetFloat64(1e18)).Float64()
	return eth
}
