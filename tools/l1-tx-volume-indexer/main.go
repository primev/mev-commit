package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	baseURL   = "https://api.covalenthq.com/v1"
	chainName = "eth-mainnet"

	// Canonical WETH on Ethereum mainnet
	wethAddress = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
)

// ---------- Covalent transaction_v2 types ----------

type TxResponse struct {
	Data struct {
		Items []struct {
			BlockSignedAt string     `json:"block_signed_at"`
			Value         string     `json:"value"` // wei, as string
			LogEvents     []LogEvent `json:"log_events"`
		} `json:"items"`
	} `json:"data"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
}

type LogEvent struct {
	SenderAddress          string   `json:"sender_address"`
	SenderContractDecimals int      `json:"sender_contract_decimals"`
	SenderContractSymbol   string   `json:"sender_contract_ticker_symbol"`
	SupportsERC            []string `json:"supports_erc"`
	Decoded                *Decoded `json:"decoded"`
}

type Decoded struct {
	Name   string  `json:"name"`   // "Transfer", "Deposit", "Withdrawal", "Swap", ...
	Params []Param `json:"params"` // includes "value"/"wad"
}

type Param struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Indexed bool        `json:"indexed"`
	Decoded bool        `json:"decoded"`
	Value   interface{} `json:"value"`
}

// ---------- Covalent pricing types ----------

type PricingResponse struct {
	Data []struct {
		ContractAddress string `json:"contract_address"`
		Prices          []struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"` // token price in ETH
		} `json:"prices"`
	} `json:"data"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
}

// ---------- Volume summary type ----------

type TxVolumes struct {
	TxHash         string
	BlockTime      time.Time
	TxValueEth     float64
	WethVolumeEth  float64
	TokenVolumeEth float64
	TotalVolumeEth float64
}

func main() {
	apiKey := os.Getenv("COVALENT_KEY")
	if apiKey == "" {
		log.Fatalf("COVALENT_KEY env var is required")
	}

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <tx_hash>  OR  %s -file txs.txt  OR  %s -fill-db", os.Args[0], os.Args[0], os.Args[0])
	}

	// DB fill mode: go run main.go -fill-db
	if os.Args[1] == "-fill-db" || os.Args[1] == "-db" {
		if err := fillDatabase(apiKey); err != nil {
			log.Fatalf("error filling database: %v", err)
		}
		return
	}

	// File mode: go run main.go -file txs.txt
	if os.Args[1] == "-file" || os.Args[1] == "-f" {
		if len(os.Args) < 3 {
			log.Fatalf("usage: %s -file txs.txt", os.Args[0])
		}
		filePath := os.Args[2]
		if err := processFile(filePath, apiKey); err != nil {
			log.Fatalf("error processing file: %v", err)
		}
		return
	}

	// Single-tx mode: go run main.go 0x....
	txHash := strings.TrimSpace(os.Args[1])
	vol, err := processTx(txHash, apiKey)
	if err != nil {
		log.Fatalf("error processing tx %s: %v", txHash, err)
	}
	fmt.Printf("Aggregate total volume for this tx: %.18f ETH\n", vol)
}

// ---------- DB fill mode ----------

// fillDatabase:
//  1. connects to StarRocks via MySQL
//  2. finds commitment_index/bidder/committer/l1_tx_hash for processed commitments
//  3. for each tx, calls Covalent to compute volumes + timestamp
//  4. inserts/upserts into mevcommit_57173.processed_l1_txns, filling only NULL volume/timestamp fields
func fillDatabase(apiKey string) error {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return fmt.Errorf("DB_USER, DB_PW, DB_HOST, DB_PORT, DB_NAME env vars are required for -fill-db mode")
	}

	// DSN with interpolateParams=true so StarRocks doesn't choke on prepared INSERTs
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?interpolateParams=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	log.Printf("Connected to DB %s at %s:%s as %s", dbName, dbHost, dbPort, dbUser)

	// 1) Ensure processed_l1_txns table exists
	const createTable = `
CREATE TABLE IF NOT EXISTS mevcommit_57173.processed_l1_txns (
  commitment_index VARCHAR(100),
  l1_timestamp     DATETIME,
  bidder           VARCHAR(64),
  committer        VARCHAR(64),
  l1_tx_hash       VARCHAR(100),
  total_vol_eth    DOUBLE,
  eth_vol          DOUBLE,
  weth_vol         DOUBLE,
  token_vol_eth    DOUBLE
)
ENGINE=OLAP
PRIMARY KEY(commitment_index)
DISTRIBUTED BY HASH(commitment_index) BUCKETS 1
PROPERTIES(
  "compression"             = "LZ4",
  "enable_persistent_index" = "true",
  "fast_schema_evolution"   = "true",
  "replicated_storage"      = "true",
  "replication_num"         = "1"
);`

	if _, err := db.Exec(createTable); err != nil {
		return fmt.Errorf("create processed_l1_txns table: %w", err)
	}

	// 2) Find commitments that already have non-null total_vol_eth
	existing := make(map[string]struct{})
	exRows, err := db.Query(`SELECT commitment_index FROM processed_l1_txns WHERE total_vol_eth IS NOT NULL`)
	if err != nil {
		return fmt.Errorf("select existing commitments: %w", err)
	}
	for exRows.Next() {
		var ci string
		if err := exRows.Scan(&ci); err != nil {
			exRows.Close()
			return fmt.Errorf("scan existing commitment_index: %w", err)
		}
		existing[ci] = struct{}{}
	}
	exRows.Close()
	log.Printf("Found %d commitments with existing non-null total_vol_eth", len(existing))

	// 3) Query StarRocks view for commitment_index + bidder/committer + L1 tx hash
	const starrocksQuery = `
      WITH opened AS (
        SELECT
          CASE
            WHEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 1, 2) = '0x'
              THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 3)
            ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex'))
          END AS commitment_index,
          LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.bidder'))    AS bidder,
          LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.committer')) AS committer,
          CASE
            WHEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')) LIKE '0x%'
              THEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash'))
            ELSE CONCAT('0x', LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')))
          END AS l1_tx_hash
        FROM mevcommit_57173.tx_view
        WHERE l_decoded IS NOT NULL
          AND COALESCE(l_removed,0) = 0
          AND get_json_string(CAST(l_decoded AS VARCHAR), '$.name') = 'OpenedCommitmentStored'
          AND t_chain_id IN (8855,57173)
      ),
      processed AS (
        SELECT
          CASE
            WHEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 1, 2) = '0x'
              THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 3)
            ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex'))
          END AS commitment_index
        FROM mevcommit_57173.tx_view
        WHERE l_decoded IS NOT NULL
          AND COALESCE(l_removed,0) = 0
          AND get_json_string(CAST(l_decoded AS VARCHAR), '$.name') = 'CommitmentProcessed'
          AND t_chain_id IN (8855,57173)
      )
      SELECT
        o.commitment_index,
        o.bidder,
        o.committer,
        o.l1_tx_hash
      FROM opened o
      JOIN processed p
        ON o.commitment_index = p.commitment_index
      WHERE o.bidder NOT IN (
          '0x4d41ab0e0b71677dfd6d02343afae96641a4c429',
          '0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e'
        )
        AND o.bidder IS NOT NULL
        AND o.l1_tx_hash IS NOT NULL
    `

	rows, err := db.Query(starrocksQuery)
	if err != nil {
		return fmt.Errorf("query StarRocks commitments: %w", err)
	}
	defer rows.Close()

	// 4) Plain INSERT (StarRocks PK table → behaves like upsert)
	insertStmt := `
      INSERT INTO processed_l1_txns (
        commitment_index,
        l1_timestamp,
        bidder,
        committer,
        l1_tx_hash,
        total_vol_eth,
        eth_vol,
        weth_vol,
        token_vol_eth
      ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	var processedCount, skippedExisting, successCount, errorCount int

	for rows.Next() {
		var commitmentIndex, bidder, committer, txHash string
		if err := rows.Scan(&commitmentIndex, &bidder, &committer, &txHash); err != nil {
			return fmt.Errorf("scan StarRocks row: %w", err)
		}
		processedCount++

		// Skip commitments that already have non-null total_vol_eth
		if _, ok := existing[commitmentIndex]; ok {
			skippedExisting++
			continue
		}

		vols, err := computeTxVolumes(txHash, apiKey)
		if err != nil {
			errorCount++
			log.Printf("commitment %s (tx %s) volume error: %v", commitmentIndex, txHash, err)
			continue
		}

		l1TimeStr := vols.BlockTime.UTC().Format("2006-01-02 15:04:05")

		if _, err := db.Exec(
			insertStmt,
			commitmentIndex,
			l1TimeStr,
			bidder,
			committer,
			txHash,
			vols.TotalVolumeEth,
			vols.TxValueEth,
			vols.WethVolumeEth,
			vols.TokenVolumeEth,
		); err != nil {
			errorCount++
			log.Printf("DB insert error for commitment %s (tx %s): %v", commitmentIndex, txHash, err)
			continue
		}

		successCount++
		log.Printf("Updated commitment %s (tx %s): total_vol_eth=%.8f, eth_vol=%.8f, weth_vol=%.8f, token_vol_eth=%.8f",
			commitmentIndex, txHash,
			vols.TotalVolumeEth, vols.TxValueEth, vols.WethVolumeEth, vols.TokenVolumeEth)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows.Err: %w", err)
	}

	log.Printf("fillDatabase done. Processed=%d, skippedExisting=%d, success=%d, errors=%d",
		processedCount, skippedExisting, successCount, errorCount)

	return nil
}

// ---------- Existing modes (file + single tx) ----------

// processFile reads one tx hash per line, calls processTx for each, and prints a final sum.
func processFile(path, apiKey string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	var sumVolume float64
	var count int

	for scanner.Scan() {
		lineNum++
		tx := strings.TrimSpace(scanner.Text())
		if tx == "" {
			continue
		}
		if !strings.HasPrefix(tx, "0x") {
			log.Printf("skip line %d (%q): not a tx hash", lineNum, tx)
			continue
		}
		vol, err := processTx(tx, apiKey)
		if err != nil {
			log.Printf("tx %s error: %v", tx, err)
			continue
		}
		sumVolume += vol
		count++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("=============================================\n")
	fmt.Printf("Processed %d transactions from %s\n", count, path)
	fmt.Printf("SUM total volume (ETH): %.18f\n", sumVolume)
	return nil
}

// processTx now just wraps computeTxVolumes, prints details, and returns total volume.
func processTx(txHash, apiKey string) (float64, error) {
	vols, err := computeTxVolumes(txHash, apiKey)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Tx: %s\n", txHash)
	fmt.Printf("Block time:             %s\n", vols.BlockTime.Format(time.RFC3339))
	fmt.Printf("Tx value (ETH):         %.18f\n", vols.TxValueEth)
	fmt.Printf("WETH volume (ETH):      %.18f\n", vols.WethVolumeEth)
	fmt.Printf("Token volume (ETH):     %.18f\n", vols.TokenVolumeEth)
	fmt.Printf("Total volume (ETH):     %.18f\n\n", vols.TotalVolumeEth)

	return vols.TotalVolumeEth, nil
}

// computeTxVolumes does the full volume calc for a single tx and returns all components.
func computeTxVolumes(txHash, apiKey string) (*TxVolumes, error) {
	txResp, err := fetchTransaction(txHash, apiKey)
	if err != nil {
		return nil, fmt.Errorf("fetchTransaction: %w", err)
	}
	if len(txResp.Data.Items) == 0 {
		return nil, fmt.Errorf("no items returned for tx %s", txHash)
	}
	item := txResp.Data.Items[0]

	// tx.value (wei) -> ETH
	txValueEth, err := weiStringToEth(item.Value)
	if err != nil {
		return nil, fmt.Errorf("parse tx value: %w", err)
	}

	// tx date for pricing
	blockTime, err := time.Parse(time.RFC3339, item.BlockSignedAt)
	if err != nil {
		return nil, fmt.Errorf("parse block_signed_at: %w", err)
	}
	dateStr := blockTime.UTC().Format("2006-01-02")

	// Collect non-WETH ERC-20 token addresses in amount-carrying events
	tokenSet := map[string]struct{}{}
	for _, ev := range item.LogEvents {
		if ev.Decoded == nil {
			continue
		}
		if !isERC20(ev.SupportsERC) {
			continue
		}
		if sameAddress(ev.SenderAddress, wethAddress) {
			continue
		}
		if isAmountEvent(ev.Decoded.Name) {
			addr := strings.ToLower(ev.SenderAddress)
			tokenSet[addr] = struct{}{}
		}
	}

	// Fetch token->ETH prices for those tokens at this date
	tokenPrices := map[string]float64{}
	if len(tokenSet) > 0 {
		tokenPrices, err = fetchTokenPricesETH(apiKey, tokenSet, dateStr)
		if err != nil {
			return nil, fmt.Errorf("fetchTokenPricesETH: %w", err)
		}
	}

	var wethVolumeEth float64
	var tokenVolumeEth float64

	for _, ev := range item.LogEvents {
		if ev.Decoded == nil {
			continue
		}
		eventName := ev.Decoded.Name

		// WETH: Deposit / Withdrawal / Transfer → direct ETH amount
		if sameAddress(ev.SenderAddress, wethAddress) && isAmountEvent(eventName) {
			amountBase, ok := extractAmountParam(ev.Decoded.Params)
			if !ok || amountBase <= 0 {
				continue
			}
			dec := ev.SenderContractDecimals
			if dec < 0 {
				dec = 0
			}
			amountEth := amountBase / math.Pow10(dec)
			wethVolumeEth += amountEth
			continue
		}

		// Other ERC-20 tokens: value them via price in ETH
		if isERC20(ev.SupportsERC) && isAmountEvent(eventName) {
			tokenAddr := strings.ToLower(ev.SenderAddress)
			priceEth, ok := tokenPrices[tokenAddr]
			if !ok || priceEth <= 0 {
				// no price available → skip this token
				continue
			}
			amountBase, ok := extractAmountParam(ev.Decoded.Params)
			if !ok || amountBase <= 0 {
				continue
			}
			dec := ev.SenderContractDecimals
			if dec < 0 {
				dec = 0
			}
			amountTokens := amountBase / math.Pow10(dec)
			volumeEth := amountTokens * priceEth
			tokenVolumeEth += volumeEth
		}
	}

	totalVolumeEth := txValueEth + wethVolumeEth + tokenVolumeEth

	return &TxVolumes{
		TxHash:         txHash,
		BlockTime:      blockTime,
		TxValueEth:     txValueEth,
		WethVolumeEth:  wethVolumeEth,
		TokenVolumeEth: tokenVolumeEth,
		TotalVolumeEth: totalVolumeEth,
	}, nil
}

// ---------- HTTP helpers ----------

func fetchTransaction(txHash, apiKey string) (*TxResponse, error) {
	const maxRetries = 5

	backoff := time.Second
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		url := fmt.Sprintf("%s/%s/transaction_v2/%s/?no-logs=false", baseURL, chainName, txHash)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("http request error: %w", err)
			// network error → retry
		} else {
			body, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				lastErr = fmt.Errorf("reading response body: %w", readErr)
			} else {
				// Retry on 429 / 5xx
				if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
					lastErr = fmt.Errorf("covalent HTTP %d for tx %s: %s",
						resp.StatusCode, txHash, truncateBody(body))
				} else if resp.StatusCode != http.StatusOK {
					// Non-retryable HTTP error
					return nil, fmt.Errorf("covalent HTTP %d for tx %s: %s",
						resp.StatusCode, txHash, truncateBody(body))
				} else {
					// 200 OK – try to parse JSON
					ct := resp.Header.Get("Content-Type")
					if !strings.Contains(strings.ToLower(ct), "application/json") {
						return nil, fmt.Errorf("non-JSON response for tx %s (Content-Type=%s): %s",
							txHash, ct, truncateBody(body))
					}

					var txResp TxResponse
					if err := json.Unmarshal(body, &txResp); err != nil {
						return nil, fmt.Errorf("JSON decode error for tx %s: %w; body: %s",
							txHash, err, truncateBody(body))
					}
					if txResp.Error {
						return nil, fmt.Errorf("covalent error for tx %s: %s",
							txHash, txResp.ErrorMessage)
					}
					return &txResp, nil
				}
			}
		}

		// If we get here, we’re going to retry (rate limit / 5xx / network)
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
			continue
		}
	}

	return nil, fmt.Errorf("fetchTransaction failed for tx %s after %d retries: %v", txHash, maxRetries, lastErr)
}

func truncateBody(b []byte) string {
	s := string(b)
	if len(s) > 300 {
		return s[:300] + "...(truncated)"
	}
	return s
}

// fetchTokenPricesETH gets token price in ETH for each address at dateStr (YYYY-MM-DD).
func fetchTokenPricesETH(apiKey string, tokenSet map[string]struct{}, dateStr string) (map[string]float64, error) {
	if len(tokenSet) == 0 {
		return map[string]float64{}, nil
	}

	addresses := make([]string, 0, len(tokenSet))
	for addr := range tokenSet {
		addresses = append(addresses, addr)
	}
	addrParam := strings.Join(addresses, ",")

	url := fmt.Sprintf(
		"%s/pricing/historical_by_addresses_v2/%s/ETH/%s/?from=%s&to=%s",
		baseURL, chainName, addrParam, dateStr, dateStr,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PricingResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	if pr.Error {
		return nil, fmt.Errorf("pricing error: %s", pr.ErrorMessage)
	}

	result := make(map[string]float64)
	for _, item := range pr.Data {
		if len(item.Prices) == 0 {
			continue
		}
		price := item.Prices[0].Price
		addr := strings.ToLower(item.ContractAddress)
		result[addr] = price
	}
	return result, nil
}

// ---------- utility helpers ----------

func isERC20(supports []string) bool {
	for _, v := range supports {
		if strings.EqualFold(v, "erc20") {
			return true
		}
	}
	return false
}

// Which events carry an amount param we want to count
func isAmountEvent(name string) bool {
	switch name {
	case "Transfer", "Deposit", "Withdrawal":
		return true
	default:
		return false
	}
}

// extractAmountParam pulls "value" or "wad" from params as float64 base units.
func extractAmountParam(params []Param) (float64, bool) {
	var raw string
	for _, p := range params {
		if p.Name == "value" || p.Name == "wad" {
			switch v := p.Value.(type) {
			case string:
				raw = v
			case float64:
				return v, true
			default:
				b, _ := json.Marshal(v)
				raw = strings.Trim(string(b), `"`)
			}
			break
		}
	}
	if raw == "" {
		return 0, false
	}
	bi, ok := new(big.Int).SetString(raw, 10)
	if !ok {
		return 0, false
	}
	f, _ := new(big.Rat).SetInt(bi).Float64()
	return f, true
}

// weiStringToEth converts decimal string wei -> ETH (float64).
func weiStringToEth(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return 0, fmt.Errorf("invalid wei: %s", s)
	}
	eth := new(big.Rat).SetFrac(bi, big.NewInt(1e18))
	f, _ := eth.Float64()
	return f, nil
}

func sameAddress(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
