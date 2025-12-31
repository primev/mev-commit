// main.go
//
// -------------------- HOW TO RUN --------------------
//
// REQUIRED ENV:
//   DB_USER, DB_PW, DB_HOST, DB_PORT, DB_NAME   (StarRocks/MySQL protocol)
//   COVALENT_KEY                                (Covalent API key)
//
// Recommended:
//   go run . -dry-run -limit 200
//   go run . -limit 200
//
// Single-tx debug (no DB writes):
//   go run . -tx 0x<hash>
//
// Force recompute/overwrite existing non-null values:
//   go run . -recompute-all -limit 500
//
// Only inserts (skip updating existing incomplete rows):
//   go run . -only-inserts -limit 500
//
// Only updates (skip discovering/inserting missing txs):
//   go run . -only-updates -limit 500
//
// ----------------------------------------------------

package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	covalentBaseURL = "https://api.covalenthq.com/v1"
	chainName       = "eth-mainnet"

	wethAddress = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"

	erc20TransferTopic0 = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	zeroAddr            = "0x0000000000000000000000000000000000000000"

	// User-provided swap topics:
	uniswapV2SwapTopic0 = "0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"
	uniswapV3SwapTopic0 = "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"
	uniswapV4SwapTopic0 = "0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f"
)

var httpClient = &http.Client{Timeout: 180 * time.Second}

// Router/settlement allowlist (used to prevent false-positive swaps)
var swapRouterAllowlist = map[string]string{
	"0x66a9893cc07d91d95644aedd05d03f95e1dba8af": "UniswapV4UniversalRouter",
	"0x0000000000001ff3684f28c67538d4d072c22734": "ZeroXAllowanceHolder",
	"0x70bf6634ee8cb27d04478f184b9b8bb13e5f4710": "ZeroXSettlerV1_6",
	"0x111111125421ca6dc452d289314280a0f8842a65": "OneInchRouterV6",
	"0x7a250d5630b4cf539739df2c5dacb4c659f2488d": "UniswapV2Router02",
	"0x888888888889758f76e7103c6cbf23abbf58f946": "PendleRouterV4",
}

// Lending pool / protocol allowlist (optional; keeps lending detection conservative)
// Fill with known pool/router addresses if you want 'touching these' to count as lending.
var lendingPoolAllowlist = map[string]string{
	"0x7d2768DE32b0b80b7a3454c06BdAc94A69DDc7A9": "AaveV2LendingPool",
	"0xEFFC18fC3b7eb8E676dac549E0c693ad50D1Ce31": "AaveV2WETHGateway",
	"0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2": "AaveV3Pool",
	"0x893411580e590D62dDBca8a703d61Cc4A8c7b2b9": "AaveV3WETHGateway",
	"0x4Ddc2D193948926D02f9B1fE9e1daa0718270ED5": "CompoundV2cETH",
	"0x5d3a536E4D6DbD6114cc1Ead35777bAB948E3643": "CompoundV2cDAI",
	"0x39AA39c021dfbaE8faC545936693aC917d5E7563": "CompoundV2cUSDC",
	"0xf650C3d88D12dB855b8bf7D11Be6C55A4e07dCC9": "CompoundV2cUSDT",
	"0xccF4429DB6322D5C611ee964527D42E5d685DD6a": "CompoundV2cWBTC",
	"0x35a18000230da775cac24873d00ff85bccded550": "CompoundV2cUNI",
	"0xC3D688B66703497DAA19211EEdff47f25384cdc3": "CompoundIIICometUSDC",
	"0x5ef30b9986345249bc32d8928B7ee64DE9435E39": "MakerDSProxyRegistry",
	"0x1476483Dd8C35F25e568113C5f70249D3976ba21": "MakerDssCdpManager",
	"0x9759A6Ac90977b93B58547b4A71c78317f391A28": "MakerDaiJoin",
}

type Candidate struct {
	Hash0x          string
	HashNorm        string
	Source          string // "events"|"rpc_only"|"existing"
	CommitmentIndex *string
	Bidder          *string
	Committer       *string
}

// Existing v2 row snapshot (for discrepancy reporting and fill-only updates)
type ExistingRow struct {
	HashNorm string
	Hash0x   string

	CommitmentIndex *string
	L1Timestamp     *time.Time
	From            *string
	To              *string
	Bidder          *string
	Committer       *string

	TotalVol *float64
	EthVol   *float64
	WethVol  *float64
	TokenVol *float64
	SwapVol  *float64

	IsSwap     *bool
	IsLending  *bool
	IsTransfer *bool
	IsApproval *bool

	PrimaryClass *string
	Protocol     *string
}

// ---------- Covalent transaction_v2 types ----------

type TxResponse struct {
	Data struct {
		Items []struct {
			BlockSignedAt string     `json:"block_signed_at"`
			Value         string     `json:"value"` // wei as string
			FromAddress   string     `json:"from_address"`
			ToAddress     string     `json:"to_address"`
			LogEvents     []LogEvent `json:"log_events"`
		} `json:"items"`
	} `json:"data"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"error_message"`
}

type LogEvent struct {
	SenderAddress          string   `json:"sender_address"`
	SenderContractDecimals int      `json:"sender_contract_decimals"`
	SupportsERC            []string `json:"supports_erc"`
	Decoded                *Decoded `json:"decoded"`

	RawLogTopics []string `json:"raw_log_topics"`
	RawLogData   string   `json:"raw_log_data"`
}

type Decoded struct {
	Name   string  `json:"name"`
	Params []Param `json:"params"`
}

type Param struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
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

type Computed struct {
	Hash0x    string
	HashNorm  string
	BlockTime time.Time
	From      string
	To        string

	TxValueEth  float64
	WethVolEth  float64
	TokenVolEth float64
	TotalVolEth float64

	// Classification
	IsSwap        bool
	IsLending     bool
	IsTransfer    bool
	IsApproval    bool
	PrimaryClass  string
	Protocol      *string
	SwapVolEth    float64
	SwapEvidence  string // "uniswap_topic"|"decoded_swap"|"router_allowlist"|"" (debug)
	SwapGuardrail string // "mint_burn_backoff"|"" (debug)
}

type workItem struct {
	Hash0x string
	Norm   string

	InsertCandidate *Candidate
	Existing        *ExistingRow
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		txHash                   = flag.String("tx", "", "single tx debug mode (no DB writes)")
		limit                    = flag.Int("limit", 0, "limit tx count (0 = no limit)")
		dryRun                   = flag.Bool("dry-run", false, "no DB writes; print counts + discrepancy summary")
		recomputeAll             = flag.Bool("recompute-all", false, "overwrite existing non-null columns with newly computed values")
		onlyInserts              = flag.Bool("only-inserts", false, "only insert missing txs; do not update existing rows")
		onlyUpdates              = flag.Bool("only-updates", false, "only update existing incomplete rows; do not insert missing txs")
		printSample              = flag.Int("print-sample", 10, "print N sample hashes")
		compareOnlyOldSwapVolGT0 = flag.Bool("compare-only-old-swapvol-gt0", false, "dry-run only: only compare discrepancies for rows where existing swap_vol_eth > 0")
		onlyOldLending           = flag.Bool(
			"only-old-lending",
			false,
			"only update/compare rows where the existing DB row has is_lending=1 (ignores incompleteness filter)",
		)
	)
	flag.Parse()

	apiKey := os.Getenv("COVALENT_KEY")
	if apiKey == "" {
		log.Fatal("COVALENT_KEY is required")
	}

	// Single tx mode (no DB writes)
	if strings.TrimSpace(*txHash) != "" {
		h := strings.ToLower(strings.TrimSpace(*txHash))
		h = ensure0x(strip0x(h))
		comp, err := computeAll(h, apiKey)
		if err != nil {
			log.Fatalf("computeAll: %v", err)
		}
		b, _ := json.MarshalIndent(comp, "", "  ")
		fmt.Println(string(b))
		return
	}

	db := mustOpenDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("db.Close: %v", err)
		}
	}()

	// 1) Discover missing txs to insert (unless onlyUpdates)
	missingToInsert := []Candidate{}
	if !*onlyUpdates {
		var err error
		missingToInsert, err = loadMissingInsertCandidates(db, *limit)
		if err != nil {
			log.Fatalf("loadMissingInsertCandidates: %v", err)
		}
	}

	// 2) Load existing rows that need update (unless onlyInserts)
	existingNeedUpdate := []ExistingRow{}
	if !*onlyInserts {
		var err error
		existingNeedUpdate, err = loadExistingNeedingUpdate(db, *limit, *onlyOldLending)
		if err != nil {
			log.Fatalf("loadExistingNeedingUpdate: %v", err)
		}
	}

	log.Printf("would_insert=%d would_update=%d (recompute_all=%v)", len(missingToInsert), len(existingNeedUpdate), *recomputeAll)

	// Print samples
	if *printSample > 0 {
		for i := 0; i < len(missingToInsert) && i < *printSample; i++ {
			ci := "<nil>"
			if missingToInsert[i].CommitmentIndex != nil {
				ci = *missingToInsert[i].CommitmentIndex
			}
			log.Printf("insert_sample[%d]=%s source=%s commitment_index=%s", i, missingToInsert[i].Hash0x, missingToInsert[i].Source, ci)
		}
		for i := 0; i < len(existingNeedUpdate) && i < *printSample; i++ {
			log.Printf("update_sample[%d]=%s", i, existingNeedUpdate[i].Hash0x)
		}
	}

	work := map[string]*workItem{}

	for i := range missingToInsert {
		c := missingToInsert[i]
		if _, ok := work[c.HashNorm]; !ok {
			work[c.HashNorm] = &workItem{Hash0x: c.Hash0x, Norm: c.HashNorm, InsertCandidate: &c}
		} else if work[c.HashNorm].InsertCandidate == nil {
			work[c.HashNorm].InsertCandidate = &c
		}
	}
	for i := range existingNeedUpdate {
		e := existingNeedUpdate[i]
		if _, ok := work[e.HashNorm]; !ok {
			work[e.HashNorm] = &workItem{Hash0x: e.Hash0x, Norm: e.HashNorm, Existing: &e}
		} else if work[e.HashNorm].Existing == nil {
			work[e.HashNorm].Existing = &e
		}
	}

	keys := make([]string, 0, len(work))
	for k := range work {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Dry-run discrepancy reporting:
	// - Compare computed vs existing for rows that exist (update list)
	// - Also useful to compute a few insert rows for sanity
	if *dryRun {
		reportDiscrepancies(db, apiKey, keys, work, *limit, *compareOnlyOldSwapVolGT0)
		log.Println("dry-run: exiting without writes")
		return
	}

	// Real run: insert missing, update existing (fill-only unless recompute-all)
	var inserted, updated, computeErr int
	for idx, k := range keys {
		if idx%100 == 0 {
			log.Printf("progress %d/%d", idx, len(keys))
		}
		w := work[k]

		comp, err := computeAll(w.Hash0x, apiKey)
		if err != nil {
			computeErr++
			log.Printf("compute error %s: %v", w.Hash0x, err)

			// If we get a"tx not found" (Covalent 404), insert a tombstone (not_found row)so we don't retry forever.
			if w.InsertCandidate != nil && !*onlyUpdates && isCovalentTxNotFound(err) {
				ok, why := shouldTombstoneNotFound(db, w.InsertCandidate.HashNorm, 75)
				if !ok {
					log.Printf("skip tombstone %s (%s)", w.Hash0x, why)
					continue
				}

				if insErr := insertV2NotFoundRow(db, *w.InsertCandidate); insErr != nil {
					log.Printf("insert not_found error %s: %v", w.Hash0x, insErr)
				} else {
					inserted++
				}
			}
			continue
		}

		// Insert missing row (if applicable)
		if w.InsertCandidate != nil && !*onlyUpdates {
			err := insertV2Row(db, *w.InsertCandidate, comp)
			if err != nil {
				// if already exists due to race, fall through to update
				log.Printf("insert error %s: %v", w.Hash0x, err)
			} else {
				inserted++
			}
		}

		// Update existing row if needed (or if recompute-all requested)
		if w.Existing != nil && !*onlyInserts {
			err := updateV2Row(db, *w.Existing, comp, *recomputeAll)
			if err != nil {
				log.Printf("update error %s: %v", w.Hash0x, err)
			} else {
				updated++
			}
		}
	}

	log.Printf("done: inserted=%d updated=%d compute_errors=%d", inserted, updated, computeErr)
}

// -------------------- DB connection --------------------

func mustOpenDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("DB_USER, DB_PW, DB_HOST, DB_PORT, DB_NAME are required")
	}
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

// -------------------- Candidate discovery --------------------

// loadMissingInsertCandidates returns tx hashes that are NOT in processed_l1_txns_v2 yet.
// It includes:
//
//	A) event-backed (OpenedCommitmentStored joined to CommitmentProcessed)
//	B) rpc-only backfill (mctransactions_sr confirmed/pre-confirmed NOT in OpenedCommitmentStored and NOT in v2)
func loadMissingInsertCandidates(db *sql.DB, limit int) ([]Candidate, error) {
	lim := ""
	if limit > 0 {
		lim = fmt.Sprintf("LIMIT %d", limit)
	}

	// A) event-backed missing from v2
	eventSQL := fmt.Sprintf(`
WITH opened AS (
  SELECT
    CASE
      WHEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 1, 2) = '0x'
        THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex')), 3)
      ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.commitmentIndex'))
    END AS commitment_index,

    CASE
      WHEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.bidder')) LIKE '0x%%'
        THEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.bidder'))
      ELSE CONCAT('0x', LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.bidder')))
    END AS bidder,

    CASE
      WHEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.committer')) LIKE '0x%%'
        THEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.committer'))
      ELSE CONCAT('0x', LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.committer')))
    END AS committer,

    CASE
      WHEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')) LIKE '0x%%'
        THEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash'))
      ELSE CONCAT('0x', LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')))
    END AS l1_tx_hash_0x,

    CASE
      WHEN LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')) LIKE '0x%%'
        THEN SUBSTR(LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash')), 3)
      ELSE LOWER(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash'))
    END AS hash_norm
  FROM mevcommit_57173.tx_view
  WHERE l_decoded IS NOT NULL
    AND COALESCE(l_removed, 0) = 0
    AND get_json_string(CAST(l_decoded AS VARCHAR), '$.name') = 'OpenedCommitmentStored'
    AND t_chain_id IN (8855, 57173)
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
    AND COALESCE(l_removed, 0) = 0
    AND get_json_string(CAST(l_decoded AS VARCHAR), '$.name') = 'CommitmentProcessed'
    AND t_chain_id IN (8855, 57173)
),
v2 AS (
  SELECT
    CASE
      WHEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 1, 2) = '0x'
        THEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 3)
      ELSE LOWER(CAST(l1_tx_hash AS VARCHAR))
    END AS hash_norm
  FROM mevcommit_57173.processed_l1_txns_v2
  WHERE l1_tx_hash IS NOT NULL
    AND CAST(l1_tx_hash AS VARCHAR) <> ''
)
SELECT
  o.l1_tx_hash_0x,
  o.hash_norm,
  o.commitment_index,
  o.bidder,
  o.committer
FROM opened o
JOIN processed p
  ON o.commitment_index = p.commitment_index
LEFT JOIN v2
  ON v2.hash_norm = o.hash_norm
WHERE o.hash_norm IS NOT NULL
  AND o.hash_norm <> ''
  AND v2.hash_norm IS NULL
  AND o.bidder IS NOT NULL
  AND CAST(o.bidder AS VARCHAR) <> ''
  AND o.bidder NOT IN (
    '0x4d41ab0e0b71677dfd6d02343afae96641a4c429',
    '0xae2885e0e7a6c5f99b93b4dbc43d206c7cf67c7e'
  )
%s;
`, lim)

	// B) rpc-only backfill missing from v2 (and missing from OpenedCommitmentStored)
	rpcSQL := fmt.Sprintf(`
WITH mc_raw AS (
  SELECT
    CAST(m.hash AS VARCHAR) AS hash_str
  FROM pg_mev_commit_fastrpc.public.mctransactions_sr m
  WHERE LOWER(CAST(m.status AS VARCHAR)) IN ('confirmed','pre-confirmed')
    AND m.hash IS NOT NULL
    AND CAST(m.hash AS VARCHAR) <> ''
),
mc AS (
  SELECT
    CASE
      WHEN SUBSTR(LOWER(hash_str), 1, 2) = '0x'
        THEN LOWER(hash_str)
      ELSE CONCAT('0x', LOWER(hash_str))
    END AS hash_0x,
    CASE
      WHEN SUBSTR(LOWER(hash_str), 1, 2) = '0x'
        THEN SUBSTR(LOWER(hash_str), 3)
      ELSE LOWER(hash_str)
    END AS hash_norm
  FROM mc_raw
),
opened_raw AS (
  SELECT
    CAST(get_json_string(CAST(l_decoded AS VARCHAR), '$.args.txnHash') AS VARCHAR) AS txn_hash_str
  FROM mevcommit_57173.tx_view
  WHERE l_decoded IS NOT NULL
    AND COALESCE(l_removed, 0) = 0
    AND CAST(get_json_string(CAST(l_decoded AS VARCHAR), '$.name') AS VARCHAR) = 'OpenedCommitmentStored'
    AND t_chain_id IN (8855, 57173)
),
opened AS (
  SELECT
    CASE
      WHEN SUBSTR(LOWER(txn_hash_str), 1, 2) = '0x'
        THEN SUBSTR(LOWER(txn_hash_str), 3)
      ELSE LOWER(txn_hash_str)
    END AS hash_norm
  FROM opened_raw
  WHERE txn_hash_str IS NOT NULL AND txn_hash_str <> ''
),
v2_raw AS (
  SELECT
    CAST(l1_tx_hash AS VARCHAR) AS l1_tx_hash_str
  FROM mevcommit_57173.processed_l1_txns_v2
  WHERE l1_tx_hash IS NOT NULL
    AND CAST(l1_tx_hash AS VARCHAR) <> ''
),
v2 AS (
  SELECT
    CASE
      WHEN SUBSTR(LOWER(l1_tx_hash_str), 1, 2) = '0x'
        THEN SUBSTR(LOWER(l1_tx_hash_str), 3)
      ELSE LOWER(l1_tx_hash_str)
    END AS hash_norm
  FROM v2_raw
)
SELECT
  mc.hash_0x,
  mc.hash_norm
FROM mc
LEFT JOIN opened
  ON mc.hash_norm = opened.hash_norm
LEFT JOIN v2
  ON mc.hash_norm = v2.hash_norm
WHERE opened.hash_norm IS NULL
  AND v2.hash_norm IS NULL
%s;
`, lim)

	out := []Candidate{}

	// Event-backed
	rows, err := db.Query(eventSQL)
	if err != nil {
		return nil, fmt.Errorf("eventSQL: %w", err)
	}
	for rows.Next() {
		var hash0x, hashNorm, ci, bidder, committer sql.NullString
		if err := rows.Scan(&hash0x, &hashNorm, &ci, &bidder, &committer); err != nil {
			_ = rows.Close()
			return nil, err
		}

		h0x := ensure0x(strip0x(strings.ToLower(strings.TrimSpace(hash0x.String))))
		hn := strip0x(strings.ToLower(strings.TrimSpace(hashNorm.String)))

		var pci, pb, pc *string
		if ci.Valid && ci.String != "" {
			s := ci.String
			pci = &s
		}
		if bidder.Valid && bidder.String != "" {
			s := bidder.String
			pb = &s
		}
		if committer.Valid && committer.String != "" {
			s := committer.String
			pc = &s
		}

		out = append(out, Candidate{
			Hash0x:          h0x,
			HashNorm:        hn,
			Source:          "events",
			CommitmentIndex: pci,
			Bidder:          pb,
			Committer:       pc,
		})
	}
	_ = rows.Close()

	// RPC-only
	rows2, err := db.Query(rpcSQL)
	if err != nil {
		return nil, fmt.Errorf("rpcSQL: %w", err)
	}
	for rows2.Next() {
		var hash0x, hashNorm sql.NullString
		if err := rows2.Scan(&hash0x, &hashNorm); err != nil {
			defer func() {
				if err := rows2.Close(); err != nil {
					log.Printf("rows2.Close: %v", err)
				}
			}()
			return nil, err
		}
		h0x := ensure0x(strip0x(strings.ToLower(strings.TrimSpace(hash0x.String))))
		hn := strip0x(strings.ToLower(strings.TrimSpace(hashNorm.String)))
		out = append(out, Candidate{
			Hash0x:   h0x,
			HashNorm: hn,
			Source:   "rpc_only",
		})
	}
	defer func() {
		if err := rows2.Close(); err != nil {
			log.Printf("rows2.Close: %v", err)
		}
	}()

	// Dedupe (prefer event-backed metadata)
	by := map[string]Candidate{}
	for _, c := range out {
		ex, ok := by[c.HashNorm]
		if !ok {
			by[c.HashNorm] = c
			continue
		}
		if ex.Source == "rpc_only" && c.Source == "events" {
			by[c.HashNorm] = c
		}
	}
	final := make([]Candidate, 0, len(by))
	for _, v := range by {
		final = append(final, v)
	}
	sort.Slice(final, func(i, j int) bool { return final[i].HashNorm < final[j].HashNorm })
	return final, nil
}

// loadExistingNeedingUpdate returns v2 rows where we still need to compute volumes/classification/from/to/timestamp.
func loadExistingNeedingUpdate(db *sql.DB, limit int, onlyOldLending bool) ([]ExistingRow, error) {
	lim := ""
	if limit > 0 {
		lim = fmt.Sprintf("LIMIT %d", limit)
	}

	q := ""
	if onlyOldLending {
		q = fmt.Sprintf(`
SELECT
  LOWER(CASE
    WHEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 1, 2) = '0x'
      THEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 3)
    ELSE LOWER(CAST(l1_tx_hash AS VARCHAR))
  END) AS hash_norm,
  CASE
    WHEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 1, 2) = '0x'
      THEN LOWER(CAST(l1_tx_hash AS VARCHAR))
    ELSE CONCAT('0x', LOWER(CAST(l1_tx_hash AS VARCHAR)))
  END AS hash0x,

  commitment_index,
  l1_timestamp,
  from_address,
  to_address,
  bidder,
  committer,

  total_vol_eth,
  eth_vol,
  weth_vol,
  token_vol_eth,
  swap_vol_eth,

  is_swap,
  is_lending,
  is_transfer,
  is_approval,

  primary_class,
  protocol

FROM mevcommit_57173.processed_l1_txns_v2
WHERE l1_tx_hash IS NOT NULL
  AND CAST(l1_tx_hash AS VARCHAR) <> ''
  AND COALESCE(is_lending, 0) = 1
%s;
`, lim)
	} else {
		q = fmt.Sprintf(`
SELECT
  LOWER(CASE
    WHEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 1, 2) = '0x'
      THEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 3)
    ELSE LOWER(CAST(l1_tx_hash AS VARCHAR))
  END) AS hash_norm,
  CASE
    WHEN SUBSTR(LOWER(CAST(l1_tx_hash AS VARCHAR)), 1, 2) = '0x'
      THEN LOWER(CAST(l1_tx_hash AS VARCHAR))
    ELSE CONCAT('0x', LOWER(CAST(l1_tx_hash AS VARCHAR)))
  END AS hash0x,

  commitment_index,
  l1_timestamp,
  from_address,
  to_address,
  bidder,
  committer,

  total_vol_eth,
  eth_vol,
  weth_vol,
  token_vol_eth,
  swap_vol_eth,

  is_swap,
  is_lending,
  is_transfer,
  is_approval,

  primary_class,
  protocol

FROM mevcommit_57173.processed_l1_txns_v2
WHERE l1_tx_hash IS NOT NULL
  AND CAST(l1_tx_hash AS VARCHAR) <> ''
  AND (primary_class IS NULL OR LOWER(CAST(primary_class AS VARCHAR)) <> 'not_found')
  AND (
       is_swap IS NULL
    OR is_lending IS NULL
    OR is_transfer IS NULL
    OR is_approval IS NULL
    OR from_address IS NULL OR CAST(from_address AS VARCHAR) = ''
    OR to_address IS NULL OR CAST(to_address AS VARCHAR) = ''
    OR l1_timestamp IS NULL
    OR total_vol_eth IS NULL
    OR swap_vol_eth IS NULL
  )
%s;
`, lim)
	}

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close: %v", err)
		}
	}()

	out := []ExistingRow{}
	for rows.Next() {
		var (
			hashNorm, hash0x sql.NullString

			ci sql.NullString
			ts sql.NullTime
			fr sql.NullString
			to sql.NullString
			bd sql.NullString
			cm sql.NullString

			total sql.NullFloat64
			ethv  sql.NullFloat64
			wethv sql.NullFloat64
			tokv  sql.NullFloat64
			swapv sql.NullFloat64

			isSwap sql.NullBool
			isLnd  sql.NullBool
			isTr   sql.NullBool
			isAp   sql.NullBool

			prim sql.NullString
			prot sql.NullString
		)

		if err := rows.Scan(
			&hashNorm, &hash0x,
			&ci, &ts, &fr, &to, &bd, &cm,
			&total, &ethv, &wethv, &tokv, &swapv,
			&isSwap, &isLnd, &isTr, &isAp,
			&prim, &prot,
		); err != nil {
			return nil, err
		}

		r := ExistingRow{
			HashNorm: hashNorm.String,
			Hash0x:   ensure0x(strip0x(hash0x.String)),
		}
		if ci.Valid && ci.String != "" {
			s := ci.String
			r.CommitmentIndex = &s
		}
		if ts.Valid {
			t := ts.Time
			r.L1Timestamp = &t
		}
		if fr.Valid && fr.String != "" {
			s := strings.ToLower(fr.String)
			r.From = &s
		}
		if to.Valid && to.String != "" {
			s := strings.ToLower(to.String)
			r.To = &s
		}
		if bd.Valid && bd.String != "" {
			s := strings.ToLower(bd.String)
			r.Bidder = &s
		}
		if cm.Valid && cm.String != "" {
			s := strings.ToLower(cm.String)
			r.Committer = &s
		}

		if total.Valid {
			v := total.Float64
			r.TotalVol = &v
		}
		if ethv.Valid {
			v := ethv.Float64
			r.EthVol = &v
		}
		if wethv.Valid {
			v := wethv.Float64
			r.WethVol = &v
		}
		if tokv.Valid {
			v := tokv.Float64
			r.TokenVol = &v
		}
		if swapv.Valid {
			v := swapv.Float64
			r.SwapVol = &v
		}

		if isSwap.Valid {
			v := isSwap.Bool
			r.IsSwap = &v
		}
		if isLnd.Valid {
			v := isLnd.Bool
			r.IsLending = &v
		}
		if isTr.Valid {
			v := isTr.Bool
			r.IsTransfer = &v
		}
		if isAp.Valid {
			v := isAp.Bool
			r.IsApproval = &v
		}

		if prim.Valid && prim.String != "" {
			s := prim.String
			r.PrimaryClass = &s
		}
		if prot.Valid && prot.String != "" {
			s := prot.String
			r.Protocol = &s
		}

		out = append(out, r)
	}
	return out, nil
}

// -------------------- Insert / Update --------------------

func insertV2Row(db *sql.DB, c Candidate, comp *Computed) error {
	// Full insert (we include candidate metadata when we have it)
	q := `
INSERT INTO mevcommit_57173.processed_l1_txns_v2 (
  l1_tx_hash,
  commitment_index,
  l1_timestamp,
  from_address,
  to_address,
  bidder,
  committer,
  total_vol_eth,
  eth_vol,
  weth_vol,
  token_vol_eth,
  swap_vol_eth,
  is_swap,
  is_lending,
  is_transfer,
  is_approval,
  primary_class,
  protocol
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`
	tsStr := comp.BlockTime.UTC().Format("2006-01-02 15:04:05")

	_, err := db.Exec(q,
		comp.HashNorm, // IMPORTANT: your v2 currently stores non-0x hashes; keep same format.
		nilOrStr(c.CommitmentIndex),
		tsStr,
		nilOrStr(ptrLower(comp.From)),
		nilOrStr(ptrLower(comp.To)),
		nilOrStr(c.Bidder),
		nilOrStr(c.Committer),
		comp.TotalVolEth,
		comp.TxValueEth,
		comp.WethVolEth,
		comp.TokenVolEth,
		comp.SwapVolEth,
		comp.IsSwap,
		comp.IsLending,
		comp.IsTransfer,
		comp.IsApproval,
		comp.PrimaryClass,
		nilOrStr(comp.Protocol),
	)
	return err
}

func insertV2NotFoundRow(db *sql.DB, c Candidate) error {
	q := `
INSERT INTO mevcommit_57173.processed_l1_txns_v2 (
  l1_tx_hash,
  commitment_index,
  bidder,
  committer,
  primary_class
) VALUES (?, ?, ?, ?, ?);
`
	primary := "not_found"

	_, err := db.Exec(q,
		c.HashNorm,
		nilOrStr(c.CommitmentIndex),
		nilOrStr(c.Bidder),
		nilOrStr(c.Committer),
		primary,
	)
	return err
}

// updateV2Row updates existing row.
// - if recomputeAll=false: only fills missing columns (NULL/empty), leaves existing non-null values unchanged
// - if recomputeAll=true: overwrites everything computed
func updateV2Row(db *sql.DB, ex ExistingRow, comp *Computed, recomputeAll bool) error {
	tsStr := comp.BlockTime.UTC().Format("2006-01-02 15:04:05")

	if recomputeAll {
		q := `
UPDATE mevcommit_57173.processed_l1_txns_v2
SET
  l1_timestamp   = ?,
  from_address   = ?,
  to_address     = ?,
  total_vol_eth  = ?,
  eth_vol        = ?,
  weth_vol       = ?,
  token_vol_eth  = ?,
  swap_vol_eth   = ?,
  is_swap        = ?,
  is_lending     = ?,
  is_transfer    = ?,
  is_approval    = ?,
  primary_class  = ?,
  protocol       = ?
WHERE LOWER(CAST(l1_tx_hash AS VARCHAR)) = ?;
`
		_, err := db.Exec(q,
			tsStr,
			strLower(comp.From),
			strLower(comp.To),
			comp.TotalVolEth,
			comp.TxValueEth,
			comp.WethVolEth,
			comp.TokenVolEth,
			comp.SwapVolEth,
			comp.IsSwap,
			comp.IsLending,
			comp.IsTransfer,
			comp.IsApproval,
			comp.PrimaryClass,
			nilOrStr(comp.Protocol),
			ex.HashNorm,
		)
		return err
	}

	// Fill-only behavior (no overwrites)
	q := `
UPDATE mevcommit_57173.processed_l1_txns_v2
SET
  l1_timestamp  = CASE WHEN l1_timestamp IS NULL THEN ? ELSE l1_timestamp END,
  from_address  = CASE WHEN from_address IS NULL OR CAST(from_address AS VARCHAR) = '' THEN ? ELSE from_address END,
  to_address    = CASE WHEN to_address IS NULL OR CAST(to_address AS VARCHAR) = '' THEN ? ELSE to_address END,

  total_vol_eth = CASE WHEN total_vol_eth IS NULL THEN ? ELSE total_vol_eth END,
  eth_vol       = CASE WHEN eth_vol IS NULL THEN ? ELSE eth_vol END,
  weth_vol      = CASE WHEN weth_vol IS NULL THEN ? ELSE weth_vol END,
  token_vol_eth = CASE WHEN token_vol_eth IS NULL THEN ? ELSE token_vol_eth END,
  swap_vol_eth  = CASE WHEN swap_vol_eth IS NULL THEN ? ELSE swap_vol_eth END,

  is_swap       = CASE WHEN is_swap IS NULL THEN ? ELSE is_swap END,
  is_lending    = CASE WHEN is_lending IS NULL THEN ? ELSE is_lending END,
  is_transfer   = CASE WHEN is_transfer IS NULL THEN ? ELSE is_transfer END,
  is_approval   = CASE WHEN is_approval IS NULL THEN ? ELSE is_approval END,

  primary_class = CASE WHEN primary_class IS NULL OR CAST(primary_class AS VARCHAR) = '' THEN ? ELSE primary_class END,
  protocol      = CASE WHEN protocol IS NULL OR CAST(protocol AS VARCHAR) = '' THEN ? ELSE protocol END

WHERE LOWER(CAST(l1_tx_hash AS VARCHAR)) = ?;
`
	_, err := db.Exec(q,
		tsStr,
		strLower(comp.From),
		strLower(comp.To),
		comp.TotalVolEth,
		comp.TxValueEth,
		comp.WethVolEth,
		comp.TokenVolEth,
		comp.SwapVolEth,
		comp.IsSwap,
		comp.IsLending,
		comp.IsTransfer,
		comp.IsApproval,
		comp.PrimaryClass,
		nilOrStr(comp.Protocol),
		ex.HashNorm,
	)
	return err
}

// -------------------- Dry-run discrepancy reporting --------------------

type deltaRow struct {
	HashNorm string

	OldTotal *float64
	NewTotal float64
	DTotal   float64

	OldSwap *float64
	NewSwap float64
	DSwap   float64

	OldClass *string
	NewClass string

	OldIsSwap *bool
	NewIsSwap bool
}

func reportDiscrepancies(db *sql.DB, apiKey string, keys []string, work map[string]*workItem, limit int, compareOnlyOldSwapVolGT0 bool) {
	// Only compare against existing rows (updates). We’ll also compute a few inserts just for sanity.
	const maxCompare = 300 // keep dry-run not crazy slow; tweak as needed
	compared := 0

	var (
		classMismatch int
		swapMismatch  int
		totalMismatch int
		swapDeltaTop  []deltaRow
		totalDeltaTop []deltaRow
	)

	for _, k := range keys {
		w := work[k]
		if w.Existing == nil && w.InsertCandidate == nil {
			continue
		}
		if compared >= maxCompare {
			break
		}
		comp, err := computeAll(w.Hash0x, apiKey)
		if err != nil {
			log.Printf("dry-run compute error %s: %v", w.Hash0x, err)
			continue
		}
		compared++

		if w.Existing == nil {
			continue // nothing to compare to
		}
		ex := w.Existing

		// Optional filter: only compare when old swap_vol_eth is non-null and > 0
		if compareOnlyOldSwapVolGT0 {
			if ex == nil || ex.SwapVol == nil || *ex.SwapVol <= 0 {
				continue
			}
		}

		// Compare classification
		if ex.PrimaryClass != nil && *ex.PrimaryClass != "" && *ex.PrimaryClass != comp.PrimaryClass {
			classMismatch++
		}
		if ex.IsSwap != nil && *ex.IsSwap != comp.IsSwap {
			swapMismatch++
		}

		// Compare total_vol_eth and swap_vol_eth with tolerance
		if ex.TotalVol != nil {
			if !almostEqual(*ex.TotalVol, comp.TotalVolEth, 1e-6, 1e-4) {
				totalMismatch++
				d := deltaRow{
					HashNorm:  k,
					OldTotal:  ex.TotalVol,
					NewTotal:  comp.TotalVolEth,
					DTotal:    comp.TotalVolEth - *ex.TotalVol,
					OldSwap:   ex.SwapVol,
					NewSwap:   comp.SwapVolEth,
					OldClass:  ex.PrimaryClass,
					NewClass:  comp.PrimaryClass,
					OldIsSwap: ex.IsSwap,
					NewIsSwap: comp.IsSwap,
				}
				totalDeltaTop = append(totalDeltaTop, d)
			}
		}
		if ex.SwapVol != nil {
			if !almostEqual(*ex.SwapVol, comp.SwapVolEth, 1e-6, 1e-4) {
				swapMismatch++ // count separately too
				d := deltaRow{
					HashNorm:  k,
					OldTotal:  ex.TotalVol,
					NewTotal:  comp.TotalVolEth,
					OldSwap:   ex.SwapVol,
					NewSwap:   comp.SwapVolEth,
					DSwap:     comp.SwapVolEth - *ex.SwapVol,
					OldClass:  ex.PrimaryClass,
					NewClass:  comp.PrimaryClass,
					OldIsSwap: ex.IsSwap,
					NewIsSwap: comp.IsSwap,
				}
				swapDeltaTop = append(swapDeltaTop, d)
			}
		}
	}

	sort.Slice(totalDeltaTop, func(i, j int) bool { return math.Abs(totalDeltaTop[i].DTotal) > math.Abs(totalDeltaTop[j].DTotal) })
	sort.Slice(swapDeltaTop, func(i, j int) bool { return math.Abs(swapDeltaTop[i].DSwap) > math.Abs(swapDeltaTop[j].DSwap) })

	log.Printf("dry-run compared=%d (cap=%d)", compared, maxCompare)
	log.Printf("discrepancies: classMismatch=%d totalVolMismatch=%d swapVolMismatch=%d", classMismatch, totalMismatch, len(swapDeltaTop))

	printTop := func(name string, arr []deltaRow, n int) {
		if len(arr) == 0 {
			return
		}
		if n > len(arr) {
			n = len(arr)
		}
		log.Printf("top %d %s deltas:", n, name)
		for i := 0; i < n; i++ {
			d := arr[i]
			oldTot := "<nil>"
			if d.OldTotal != nil {
				oldTot = fmt.Sprintf("%.8f", *d.OldTotal)
			}
			oldSwap := "<nil>"
			if d.OldSwap != nil {
				oldSwap = fmt.Sprintf("%.8f", *d.OldSwap)
			}
			oldClass := "<nil>"
			if d.OldClass != nil {
				oldClass = *d.OldClass
			}
			oldIsSwap := "<nil>"
			if d.OldIsSwap != nil {
				oldIsSwap = fmt.Sprintf("%v", *d.OldIsSwap)
			}
			log.Printf("[%d] %s old_total=%s new_total=%.8f d_total=%.8f old_swap=%s new_swap=%.8f d_swap=%.8f old_class=%s new_class=%s old_isSwap=%s new_isSwap=%v",
				i, d.HashNorm, oldTot, d.NewTotal, d.DTotal, oldSwap, d.NewSwap, d.DSwap, oldClass, d.NewClass, oldIsSwap, d.NewIsSwap)
		}
	}

	printTop("TOTAL", totalDeltaTop, 15)
	printTop("SWAP", swapDeltaTop, 15)
}

func almostEqual(old, new float64, absTol, relTol float64) bool {
	diff := math.Abs(old - new)
	if diff <= absTol {
		return true
	}
	den := math.Max(1.0, math.Max(math.Abs(old), math.Abs(new)))
	return (diff / den) <= relTol
}

// -------------------- Compute + classification --------------------

func computeAll(txHash0x string, apiKey string) (*Computed, error) {
	txHash0x = ensure0x(strip0x(strings.ToLower(strings.TrimSpace(txHash0x))))
	hashNorm := strip0x(txHash0x)

	txResp, err := fetchTransaction(txHash0x, apiKey)
	if err != nil {
		return nil, err
	}
	if len(txResp.Data.Items) == 0 {
		return nil, fmt.Errorf("no items for tx %s", txHash0x)
	}
	item := txResp.Data.Items[0]

	blockTime, err := time.Parse(time.RFC3339, item.BlockSignedAt)
	if err != nil {
		return nil, fmt.Errorf("parse block_signed_at: %w", err)
	}

	fromAddr := strings.ToLower(strings.TrimSpace(item.FromAddress))
	toAddr := strings.ToLower(strings.TrimSpace(item.ToAddress))

	txValueEth, err := weiStringToEth(item.Value)
	if err != nil {
		return nil, err
	}

	wethVol, tokenVol, err := computeWethAndTokenVolEth(item.LogEvents, blockTime, apiKey)
	if err != nil {
		return nil, err
	}
	totalVol := txValueEth + wethVol + tokenVol

	// Swap evidence:
	swapIsUni := hasAnyTopic0(item.LogEvents, uniswapV2SwapTopic0, uniswapV3SwapTopic0, uniswapV4SwapTopic0)
	hasRouter := touchedDexInfra(item.LogEvents, toAddr)
	hasSwapEvent := hasDecodedSwapLikeEvent(item.LogEvents)

	// Swap detection with false-positive protections:
	isSwap, swapVol, proto, evidence, guardrail := detectSwapAndVolume(item.LogEvents, fromAddr, toAddr, txValueEth, blockTime, apiKey, swapIsUni, hasSwapEvent, hasRouter)

	// Other classifications:
	isApproval := hasApprovalLike(item.LogEvents)
	isTransfer := isPlainTransfer(item.LogEvents, txValueEth)

	// Lending signals (stricter):
	// - decoded borrow/repay/liquidate/redeem signals
	// - OR explicit allowlist touch (optional)
	isLending := false
	if !isSwap {
		if hasLendingBorrowSignals(item.LogEvents) || touchedLendingInfra(item.LogEvents, toAddr) {
			isLending = true
		}
	}

	primary := primaryClass(isSwap, isLending, isTransfer, isApproval)

	return &Computed{
		Hash0x:        txHash0x,
		HashNorm:      hashNorm,
		BlockTime:     blockTime,
		From:          fromAddr,
		To:            toAddr,
		TxValueEth:    txValueEth,
		WethVolEth:    wethVol,
		TokenVolEth:   tokenVol,
		TotalVolEth:   totalVol,
		IsSwap:        isSwap,
		SwapVolEth:    swapVol,
		IsLending:     isLending,
		IsTransfer:    isTransfer,
		IsApproval:    isApproval,
		PrimaryClass:  primary,
		Protocol:      proto,
		SwapEvidence:  evidence,
		SwapGuardrail: guardrail,
	}, nil
}

// detectSwapAndVolume enforces your anti-false-positive requirement:
// We only label swap if there is strong evidence:
//   - Uniswap swap topic0, OR
//   - decoded swap-like event, OR
//   - router allowlist interaction
//
// AND: mint/burn guardrail backoff if we only have router evidence (no topic/event).
func detectSwapAndVolume(
	logs []LogEvent,
	fromAddr, toAddr string,
	txValueEth float64,
	blockTime time.Time,
	apiKey string,
	swapIsUni bool,
	hasSwapEvent bool,
	hasRouter bool,
) (isSwap bool, swapVol float64, protocol *string, evidence string, guardrail string) {

	if swapIsUni {
		isSwap = true
		evidence = "uniswap_topic"
		p := "uniswap"
		protocol = &p
	} else if hasSwapEvent {
		isSwap = true
		evidence = "decoded_swap"
	} else if hasRouter {
		isSwap = true
		evidence = "router_allowlist"
		// Protocol inference from router allowlist label (e.g. PendleRouterV4)
		label := swapRouterAllowlist[strings.ToLower(strings.TrimSpace(toAddr))]
		if strings.Contains(strings.ToLower(label), "pendle") {
			p := "pendle"
			protocol = &p
		}
	}

	if !isSwap {
		return false, 0, nil, "", ""
	}

	// Guardrail: mint/burn markers can appear in real swaps (wrappers, Pendle SY/PT/YT/LPT, etc).
	// Only back off when the ONLY swap evidence is a router allowlist hit (i.e. no swap topics and no decoded swap-like events).
	strongEvidence := (evidence == "uniswap_topic" || evidence == "decoded_swap")

	// Treat known Pendle router hits as strong (Pendle frequently mints/burns wrapper tokens inside swaps).
	if evidence == "router_allowlist" && protocol != nil && strings.EqualFold(*protocol, "pendle") {
		strongEvidence = true
	}

	if evidence == "router_allowlist" && hasMintBurnMarkers(logs) && !strongEvidence {
		guardrail = "mint_burn_backoff"
		return false, 0, nil, evidence, guardrail
	}

	vol, err := computeSwapInputEthHeuristic(logs, fromAddr, toAddr, txValueEth, blockTime, apiKey)
	if err != nil {
		return true, 0, protocol, evidence, guardrail
	}
	return true, vol, protocol, evidence, guardrail
}

func primaryClass(isSwap, isLending, isTransfer, isApproval bool) string {
	if isSwap {
		return "swap"
	}
	if isLending {
		return "lending"
	}
	if isTransfer {
		return "transfer"
	}
	if isApproval {
		return "approval"
	}
	return "unknown"
}

// -------------------- Covalent fetch + pricing --------------------

func fetchTransaction(txHash0x, apiKey string) (*TxResponse, error) {
	apiKey = strings.TrimSpace(apiKey)

	url := fmt.Sprintf("%s/%s/transaction_v2/%s/?no-logs=false", covalentBaseURL, chainName, txHash0x)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Covalent v1 API auth: HTTP Basic (key as username, empty password)
	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Accept", "application/json")

	start := time.Now()
	resp, err := httpClient.Do(req)
	dur := time.Since(start)
	if err != nil {
		return nil, fmt.Errorf("tx request error after %s: %w", dur, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("resp.Body.Close: %v", err)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read tx body: %w", readErr)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("covalent tx HTTP %d: %s", resp.StatusCode, truncateBody(body))
	}

	var txResp TxResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("covalent tx JSON decode: %w; body: %s", err, truncateBody(body))
	}
	if txResp.Error {
		log.Printf("fetchTransaction %s status=%d", txHash0x, resp.StatusCode)
		return nil, fmt.Errorf("covalent tx error: %s", txResp.ErrorMessage)
	}
	return &txResp, nil
}

func isCovalentTxNotFound(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	if !strings.Contains(s, "covalent tx HTTP 404:") {
		return false
	}
	ls := strings.ToLower(s)
	return strings.Contains(ls, "transaction hash:") && strings.Contains(ls, " not found")
}

func shouldTombstoneNotFound(db *sql.DB, txHashNorm string, minAgeBlocks int64) (bool, string) {
	// txHashNorm is expected WITHOUT 0x, lowercased
	txHashNorm = strings.ToLower(strings.TrimSpace(strip0x(txHashNorm)))

	// One scan over mc table:
	// - head_block = MAX(block_number)
	// - tx_block = MAX(block_number) for that hash (if present)
	q := `
SELECT
  MAX(CAST(block_number AS BIGINT)) AS head_block,
  MAX(CASE
        WHEN LOWER(REPLACE(CAST(hash AS VARCHAR), '0x', '')) = ?
        THEN CAST(block_number AS BIGINT)
      END) AS tx_block
FROM pg_mev_commit_fastrpc.public.mctransactions_sr
WHERE LOWER(CAST(status AS VARCHAR)) IN ('confirmed','pre-confirmed');
`

	var head sql.NullInt64
	var txb sql.NullInt64
	if err := db.QueryRow(q, txHashNorm).Scan(&head, &txb); err != nil {
		return false, fmt.Sprintf("age_query_error: %v", err)
	}
	if !txb.Valid || txb.Int64 == 0 {
		// We can’t prove it’s old; be conservative and don’t tombstone.
		return false, "no_tx_block"
	}
	if !head.Valid || head.Int64 == 0 {
		return false, "no_head_block"
	}

	age := head.Int64 - txb.Int64
	if age <= minAgeBlocks {
		return false, fmt.Sprintf("too_recent age_blocks=%d head=%d txb=%d", age, head.Int64, txb.Int64)
	}
	return true, fmt.Sprintf("old_enough age_blocks=%d head=%d txb=%d", age, head.Int64, txb.Int64)
}

func fetchTokenPricesETH(apiKey string, tokenSet map[string]struct{}, dateStr string) (map[string]float64, error) {
	if len(tokenSet) == 0 {
		return map[string]float64{}, nil
	}
	addrs := make([]string, 0, len(tokenSet))
	for a := range tokenSet {
		addrs = append(addrs, a)
	}
	addrParam := strings.Join(addrs, ",")

	url := fmt.Sprintf("%s/pricing/historical_by_addresses_v2/%s/ETH/%s/?from=%s&to=%s",
		covalentBaseURL, chainName, addrParam, dateStr, dateStr)

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pricing request error: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("resp.Body.Close: %v", err)
		}
	}()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read pricing body: %w", readErr)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("covalent pricing HTTP %d: %s", resp.StatusCode, truncateBody(body))
	}

	var pr PricingResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("pricing JSON decode: %w; body: %s", err, truncateBody(body))
	}
	if pr.Error {
		return nil, fmt.Errorf("pricing error: %s", pr.ErrorMessage)
	}

	out := make(map[string]float64)
	for _, item := range pr.Data {
		if len(item.Prices) == 0 {
			continue
		}
		out[strings.ToLower(item.ContractAddress)] = item.Prices[0].Price
	}
	return out, nil
}

// -------------------- Base volume calc --------------------

func computeWethAndTokenVolEth(logs []LogEvent, blockTime time.Time, apiKey string) (wethVolEth float64, tokenVolEth float64, err error) {
	dateStr := blockTime.UTC().Format("2006-01-02")

	tokenSet := map[string]struct{}{}
	for _, ev := range logs {
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
			tokenSet[strings.ToLower(ev.SenderAddress)] = struct{}{}
		}
	}

	tokenPrices := map[string]float64{}
	if len(tokenSet) > 0 {
		tokenPrices, err = fetchTokenPricesETH(apiKey, tokenSet, dateStr)
		if err != nil {
			return 0, 0, err
		}
	}

	for _, ev := range logs {
		if ev.Decoded == nil {
			continue
		}
		name := ev.Decoded.Name

		if sameAddress(ev.SenderAddress, wethAddress) && isAmountEvent(name) {
			amountBase, ok := extractAmountParam(ev.Decoded.Params)
			if !ok || amountBase <= 0 {
				continue
			}
			dec := clampDec(ev.SenderContractDecimals)
			wethVolEth += amountBase / math.Pow10(dec)
			continue
		}

		if isERC20(ev.SupportsERC) && isAmountEvent(name) {
			addr := strings.ToLower(ev.SenderAddress)
			price := tokenPrices[addr]
			if price <= 0 {
				continue
			}
			amountBase, ok := extractAmountParam(ev.Decoded.Params)
			if !ok || amountBase <= 0 {
				continue
			}
			dec := clampDec(ev.SenderContractDecimals)
			amtTokens := amountBase / math.Pow10(dec)
			tokenVolEth += amtTokens * price
		}
	}
	return wethVolEth, tokenVolEth, nil
}

// -------------------- Swap volume heuristic --------------------

func computeSwapInputEthHeuristic(logs []LogEvent, fromAddr, toAddr string, txValueEth float64, blockTime time.Time, apiKey string) (float64, error) {
	tokenSet := map[string]struct{}{}
	for _, ev := range logs {
		addr := strings.ToLower(strings.TrimSpace(ev.SenderAddress))
		if addr == "" {
			continue
		}
		if (ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params)) || isRawERC20Transfer(ev) {
			tokenSet[addr] = struct{}{}
		}
	}

	dateStr := blockTime.UTC().Format("2006-01-02")
	tokenPrices, err := fetchTokenPricesETH(apiKey, tokenSet, dateStr)
	if err != nil {
		return 0, err
	}

	getPrice := func(token string) float64 {
		token = strings.ToLower(strings.TrimSpace(token))
		if token == "" {
			return 0
		}
		if sameAddress(token, wethAddress) {
			return 1.0
		}
		if p := tokenPrices[token]; p > 0 {
			return p
		}
		return 0
	}

	// "User-like" net outflow from fromAddr (matches your earlier logic style)
	userLike := addressHasERC20OutAny(logs, fromAddr) || txValueEth > 0
	if userLike {
		type flow struct{ in, out float64 }
		flows := map[string]*flow{}

		add := func(token, f, t string, amt float64) {
			token = strings.ToLower(strings.TrimSpace(token))
			f = strings.ToLower(strings.TrimSpace(f))
			t = strings.ToLower(strings.TrimSpace(t))
			if isZero(f) || isZero(t) {
				return
			}
			if token == "" || amt <= 0 {
				return
			}
			if f != fromAddr && t != fromAddr {
				return
			}
			x := flows[token]
			if x == nil {
				x = &flow{}
				flows[token] = x
			}
			if f == fromAddr {
				x.out += amt
			}
			if t == fromAddr {
				x.in += amt
			}
		}

		for _, ev := range logs {
			token := strings.ToLower(strings.TrimSpace(ev.SenderAddress))
			if token == "" {
				continue
			}
			if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
				f, t := decodedFromTo(ev.Decoded.Params)
				if isZero(f) || isZero(t) {
					continue
				}
				base, ok := extractAmountParam(ev.Decoded.Params)
				if !ok || base <= 0 {
					continue
				}
				dec := clampDec(ev.SenderContractDecimals)
				add(token, f, t, base/math.Pow10(dec))
				continue
			}
			if isRawERC20Transfer(ev) && len(ev.RawLogTopics) >= 3 {
				f, ok1 := topicToAddress(ev.RawLogTopics[1])
				t, ok2 := topicToAddress(ev.RawLogTopics[2])
				v, ok3 := hexToBigInt(ev.RawLogData)
				if !ok1 || !ok2 || !ok3 {
					continue
				}
				if isZero(f) || isZero(t) {
					continue
				}
				dec := clampDec(ev.SenderContractDecimals)
				add(token, f, t, bigIntToScaledFloat(v, dec))
			}
		}

		swapIn := txValueEth
		for token, fl := range flows {
			price := getPrice(token)
			if price <= 0 {
				continue
			}
			netOut := fl.out - fl.in
			if netOut > 1e-12 {
				swapIn += netOut * price
			}
		}
		return swapIn, nil
	}

	// Fallback: best gross min(in,out) over candidate addresses
	cands := candidateAddresses(logs, fromAddr, toAddr)
	best := 0.0
	for _, trader := range cands {
		type flow struct{ in, out float64 }
		flows := map[string]*flow{}
		add := func(token, f, t string, amt float64) {
			token = strings.ToLower(strings.TrimSpace(token))
			f = strings.ToLower(strings.TrimSpace(f))
			t = strings.ToLower(strings.TrimSpace(t))
			if isZero(f) || isZero(t) {
				return
			}
			if token == "" || amt <= 0 {
				return
			}
			if f != trader && t != trader {
				return
			}
			x := flows[token]
			if x == nil {
				x = &flow{}
				flows[token] = x
			}
			if f == trader {
				x.out += amt
			}
			if t == trader {
				x.in += amt
			}
		}

		for _, ev := range logs {
			token := strings.ToLower(strings.TrimSpace(ev.SenderAddress))
			if token == "" {
				continue
			}
			if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
				f, t := decodedFromTo(ev.Decoded.Params)
				base, ok := extractAmountParam(ev.Decoded.Params)
				if !ok || base <= 0 {
					continue
				}
				dec := clampDec(ev.SenderContractDecimals)
				add(token, f, t, base/math.Pow10(dec))
				continue
			}
			if isRawERC20Transfer(ev) && len(ev.RawLogTopics) >= 3 {
				f, ok1 := topicToAddress(ev.RawLogTopics[1])
				t, ok2 := topicToAddress(ev.RawLogTopics[2])
				v, ok3 := hexToBigInt(ev.RawLogData)
				if !ok1 || !ok2 || !ok3 {
					continue
				}
				dec := clampDec(ev.SenderContractDecimals)
				add(token, f, t, bigIntToScaledFloat(v, dec))
			}
		}

		for token, fl := range flows {
			price := getPrice(token)
			if price <= 0 {
				continue
			}
			gross := math.Min(fl.in, fl.out)
			if gross <= 1e-12 {
				continue
			}
			eth := gross * price
			if eth > best {
				best = eth
			}
		}
	}
	return best, nil
}

// -------------------- Classification helpers --------------------

func hasAnyTopic0(logs []LogEvent, topic0s ...string) bool {
	set := map[string]struct{}{}
	for _, t := range topic0s {
		set[strings.ToLower(t)] = struct{}{}
	}
	for _, ev := range logs {
		if len(ev.RawLogTopics) == 0 {
			continue
		}
		t0 := strings.ToLower(strings.TrimSpace(ev.RawLogTopics[0]))
		if _, ok := set[t0]; ok {
			return true
		}
	}
	return false
}

func hasDecodedSwapLikeEvent(logs []LogEvent) bool {
	for _, ev := range logs {
		if ev.Decoded == nil {
			continue
		}
		if isSwapLikeEvent(ev.Decoded.Name) {
			return true
		}
	}
	return false
}

func hasApprovalLike(logs []LogEvent) bool {
	for _, ev := range logs {
		if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Approval") {
			return true
		}
	}
	return false
}

func isPlainTransfer(logs []LogEvent, txValueEth float64) bool {
	if txValueEth > 0 {
		return true
	}
	for _, ev := range logs {
		if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
			return true
		}
		if isRawERC20Transfer(ev) {
			return true
		}
	}
	return false
}

// Conservative lending/borrow signals:
// - decoded event names that commonly show up in lending markets
func hasLendingBorrowSignals(logs []LogEvent) bool {
	for _, ev := range logs {
		if ev.Decoded == nil {
			continue
		}
		n := strings.ToLower(ev.Decoded.Name)
		if strings.Contains(n, "borrow") ||
			strings.Contains(n, "repay") ||
			strings.Contains(n, "liquidat") {
			return true
		}
	}
	return false
}

// Mint/burn markers via Transfer from/to zero address.
func hasMintBurnMarkers(logs []LogEvent) bool {
	for _, ev := range logs {
		if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
			f, t := decodedFromTo(ev.Decoded.Params)
			if isZero(f) || isZero(t) {
				return true
			}
		}
		if isRawERC20Transfer(ev) && len(ev.RawLogTopics) >= 3 {
			f, ok1 := topicToAddress(ev.RawLogTopics[1])
			t, ok2 := topicToAddress(ev.RawLogTopics[2])
			if ok1 && ok2 && (isZero(f) || isZero(t)) {
				return true
			}
		}
	}
	return false
}

func touchedDexInfra(logs []LogEvent, toAddr string) bool {
	toAddr = strings.ToLower(strings.TrimSpace(toAddr))
	if _, ok := swapRouterAllowlist[toAddr]; ok {
		return true
	}
	for _, ev := range logs {
		emitter := strings.ToLower(strings.TrimSpace(ev.SenderAddress))
		if _, ok := swapRouterAllowlist[emitter]; ok {
			return true
		}
	}
	return false
}

func touchedLendingInfra(logs []LogEvent, toAddr string) bool {
	toAddr = strings.ToLower(strings.TrimSpace(toAddr))
	if _, ok := lendingPoolAllowlist[toAddr]; ok {
		return true
	}
	for _, ev := range logs {
		emitter := strings.ToLower(strings.TrimSpace(ev.SenderAddress))
		if _, ok := lendingPoolAllowlist[emitter]; ok {
			return true
		}
	}
	return false
}

// -------------------- Low-level log parsing --------------------

func isERC20(supports []string) bool {
	for _, v := range supports {
		if strings.EqualFold(v, "erc20") {
			return true
		}
	}
	return false
}

func isAmountEvent(name string) bool {
	switch name {
	case "Transfer", "Deposit", "Withdrawal":
		return true
	default:
		return false
	}
}

func isSwapLikeEvent(name string) bool {
	// Keep explicit allowlist, but also treat any event name containing "swap" (case-insensitive) as swap-like.
	n := strings.ToLower(strings.TrimSpace(name))
	if strings.Contains(n, "swap") {
		return true
	}
	switch name {
	case "Swap", "TokenExchange", "TokenExchangeUnderlying", "Swapped",
		"Trade", "OrderSettled", "Settlement",
		"Fill", "LimitOrderFilled", "RfqOrderFilled", "OrderFilled",
		"TransformERC20", "TransformedERC20", "ERC20BridgeTransfer",
		"DODOSwap", "DODOV2SellBaseToken", "DODOV2SellQuoteToken", "Buy", "Sell":
		return true
	default:
		return false
	}
}

func hasFromToValue(params []Param) bool {
	hasFrom, hasTo, hasVal := false, false, false
	for _, p := range params {
		switch p.Name {
		case "from":
			hasFrom = true
		case "to":
			hasTo = true
		case "value", "wad":
			hasVal = true
		}
	}
	return hasFrom && hasTo && hasVal
}

func decodedFromTo(params []Param) (string, string) {
	var f, t string
	for _, p := range params {
		if p.Name == "from" {
			if s, ok := p.Value.(string); ok {
				f = strings.ToLower(s)
			}
		}
		if p.Name == "to" {
			if s, ok := p.Value.(string); ok {
				t = strings.ToLower(s)
			}
		}
	}
	return f, t
}

func extractAmountParam(params []Param) (float64, bool) {
	var raw string
	for _, p := range params {
		if p.Name == "value" || p.Name == "wad" {
			switch v := p.Value.(type) {
			case string:
				raw = v
			case float64:
				raw = fmt.Sprintf("%.0f", v)
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

func isRawERC20Transfer(ev LogEvent) bool {
	return len(ev.RawLogTopics) >= 3 && strings.EqualFold(ev.RawLogTopics[0], erc20TransferTopic0)
}

func topicToAddress(topic string) (string, bool) {
	topic = strings.TrimSpace(topic)
	if !strings.HasPrefix(topic, "0x") || len(topic) != 66 {
		return "", false
	}
	return "0x" + strings.ToLower(topic[len(topic)-40:]), true
}

func hexToBigInt(hexStr string) (*big.Int, bool) {
	hexStr = strings.TrimSpace(hexStr)
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		return nil, false
	}
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, false
	}
	return new(big.Int).SetBytes(b), true
}

func bigIntToScaledFloat(v *big.Int, decimals int) float64 {
	if v == nil {
		return 0
	}
	decimals = clampDec(decimals)
	r := new(big.Rat).SetInt(v)
	den := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	r.Quo(r, den)
	f, _ := r.Float64()
	return f
}

func addressHasERC20OutAny(logs []LogEvent, addr string) bool {
	addr = strings.ToLower(strings.TrimSpace(addr))
	if addr == "" {
		return false
	}
	for _, ev := range logs {
		if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
			f, _ := decodedFromTo(ev.Decoded.Params)
			if strings.ToLower(strings.TrimSpace(f)) == addr {
				return true
			}
		}
		if isRawERC20Transfer(ev) {
			f, ok := topicToAddress(ev.RawLogTopics[1])
			if ok && strings.ToLower(f) == addr {
				return true
			}
		}
	}
	return false
}

func candidateAddresses(logs []LogEvent, fromAddr, toAddr string) []string {
	cands := []string{}
	push := func(a string) {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" {
			return
		}
		for _, x := range cands {
			if x == a {
				return
			}
		}
		cands = append(cands, a)
	}
	push(fromAddr)
	push(toAddr)

	counts := map[string]int{}
	bump := func(a string) {
		a = strings.ToLower(strings.TrimSpace(a))
		if a == "" || isZero(a) {
			return
		}
		counts[a]++
	}
	for _, ev := range logs {
		if ev.Decoded != nil && strings.EqualFold(ev.Decoded.Name, "Transfer") && hasFromToValue(ev.Decoded.Params) {
			f, t := decodedFromTo(ev.Decoded.Params)
			if isZero(f) || isZero(t) {
				continue
			}
			bump(f)
			bump(t)
			continue
		}
		if isRawERC20Transfer(ev) && len(ev.RawLogTopics) >= 3 {
			f, ok1 := topicToAddress(ev.RawLogTopics[1])
			t, ok2 := topicToAddress(ev.RawLogTopics[2])
			if ok1 && ok2 && !isZero(f) && !isZero(t) {
				bump(f)
				bump(t)
			}
		}
	}
	type kv struct {
		a string
		c int
	}
	kvs := make([]kv, 0, len(counts))
	for a, c := range counts {
		kvs = append(kvs, kv{a: a, c: c})
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].c > kvs[j].c })
	for i := 0; i < len(kvs) && i < 6; i++ {
		push(kvs[i].a)
	}
	return cands
}

// -------------------- Utilities --------------------

func weiStringToEth(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return 0, nil
	}
	bi, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return 0, fmt.Errorf("invalid wei string: %s", s)
	}
	eth := new(big.Rat).SetFrac(bi, big.NewInt(1e18))
	f, _ := eth.Float64()
	return f, nil
}

func sameAddress(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func truncateBody(b []byte) string {
	s := string(b)
	if len(s) > 300 {
		return s[:300] + "...(truncated)"
	}
	return s
}

func strip0x(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if strings.HasPrefix(s, "0x") {
		return s[2:]
	}
	return s
}

func ensure0x(hashNorm string) string {
	hashNorm = strings.ToLower(strings.TrimSpace(hashNorm))
	if hashNorm == "" {
		return ""
	}
	if strings.HasPrefix(hashNorm, "0x") {
		return hashNorm
	}
	return "0x" + hashNorm
}

func isZero(a string) bool {
	return strings.EqualFold(strings.ToLower(strings.TrimSpace(a)), zeroAddr)
}

func clampDec(d int) int {
	if d < 0 {
		return 0
	}
	if d > 36 {
		return 36
	}
	return d
}

func nilOrStr(p *string) interface{} {
	if p == nil {
		return nil
	}
	s := strings.TrimSpace(*p)
	if s == "" {
		return nil
	}
	return s
}

func strLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ptrLower(s string) *string {
	x := strings.ToLower(strings.TrimSpace(s))
	return &x
}
