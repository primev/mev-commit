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
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"slices"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/go-sql-driver/mysql"
	fastsettlement "github.com/primev/mev-commit/contracts-abi/clients/FastSettlementV3"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

// 0.00001 ETH in wei — same rate as preconf-rpc/points.
const weiPerPoint = 10_000_000_000_000

const (
	defaultContract  = "0x084C0EC7f5C0585195c1c713ED9f06272F48cB45"
	defaultWETH      = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	defaultRecipient = "0xD5881f91270550B8850127f05BD6C8C203B3D33f"
	permit2Addr      = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
)

var zeroAddr = common.Address{}

func main() {
	app := &cli.App{
		Name:  "fastswap-miles",
		Usage: "MEV Commit FastSwap Miles Service",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "contract", Value: defaultContract, Usage: "FastSettlementV3 proxy address", EnvVars: []string{"FASTSWAP_CONTRACT"}},
			&cli.StringFlag{Name: "weth", Value: defaultWETH, Usage: "WETH contract address", EnvVars: []string{"WETH_CONTRACT"}},
			&cli.Uint64Flag{Name: "start-block", Value: 0, Usage: "block to start indexing from (0 = auto-resume from DB)", EnvVars: []string{"FASTSWAP_START_BLOCK"}},
			&cli.DurationFlag{Name: "poll", Value: 12 * time.Second, Usage: "poll interval for new blocks", EnvVars: []string{"FASTSWAP_POLL_INTERVAL"}},
			&cli.Uint64Flag{Name: "batch", Value: 2000, Usage: "blocks per eth_getLogs batch", EnvVars: []string{"FASTSWAP_BATCH_SIZE"}},
			&cli.BoolFlag{Name: "dry-run", Value: false, Usage: "index events and compute miles but skip Fuel submission and processed marking", EnvVars: []string{"FASTSWAP_DRY_RUN"}},
			&cli.StringFlag{Name: "keystore", Usage: "Path to the executor keystore file (required for production token sweeping)", EnvVars: []string{"FASTSWAP_KEYSTORE"}},
			&cli.StringFlag{Name: "passphrase", Usage: "Password for the keystore", EnvVars: []string{"FASTSWAP_PASSPHRASE"}},
			&cli.StringFlag{Name: "fastswap-url", Usage: "FastSwap API endpoint URL (e.g., https://fastrpc.primev.xyz)", EnvVars: []string{"FASTSWAP_URL"}},
			&cli.StringFlag{Name: "funds-recipient", Value: defaultRecipient, Usage: "Address to receive swept ETH", EnvVars: []string{"FASTSWAP_FUNDS_RECIPIENT"}},
			&cli.Uint64Flag{Name: "max-gas-gwei", Value: 50, Usage: "Skip sweep if L1 gas price exceeds this (gwei)", EnvVars: []string{"FASTSWAP_MAX_GAS_GWEI"}},
			&cli.StringFlag{Name: "db-user", Required: true, Usage: "Database user", EnvVars: []string{"DB_USER"}},
			&cli.StringFlag{Name: "db-pw", Required: true, Usage: "Database password", EnvVars: []string{"DB_PASSWORD"}},
			&cli.StringFlag{Name: "db-host", Value: "127.0.0.1", Usage: "Database host", EnvVars: []string{"DB_HOST"}},
			&cli.StringFlag{Name: "db-port", Value: "9030", Usage: "Database port", EnvVars: []string{"DB_PORT"}},
			&cli.StringFlag{Name: "db-name", Value: "mevcommit_57173", Usage: "Database name", EnvVars: []string{"DB_NAME"}},
			&cli.StringFlag{Name: "l1-rpc-url", Required: true, Usage: "L1 Ethereum HTTP RPC URL", EnvVars: []string{"L1_RPC_URL"}},
			&cli.StringFlag{Name: "barter-url", Required: true, Usage: "Barter API base URL", EnvVars: []string{"BARTER_URL"}},
			&cli.StringFlag{Name: "barter-api-key", Usage: "Barter API Key", EnvVars: []string{"BARTER_KEY"}},
			&cli.StringFlag{Name: "fuel-api-url", Usage: "Fuel points API URL", EnvVars: []string{"FUEL_URL"}},
			&cli.StringFlag{Name: "fuel-api-key", Usage: "Fuel points API Key", EnvVars: []string{"FUEL_API_KEY"}},
			&cli.StringFlag{
				Name:    "log-fmt",
				Usage:   "log format to use, options are 'text' or 'json'",
				EnvVars: []string{"LOG_FMT"},
				Value:   "json",
				Action: func(ctx *cli.Context, s string) error {
					if !slices.Contains([]string{"text", "json"}, s) {
						return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "info",
				Action: func(ctx *cli.Context, s string) error {
					if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
						return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:    "log-tags",
				Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
				EnvVars: []string{"LOG_TAGS"},
				Action: func(ctx *cli.Context, s string) error {
					if s == "" {
						return nil
					}
					for i, p := range strings.Split(s, ",") {
						if len(strings.Split(p, ":")) != 2 {
							return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
						}
					}
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String("log-level"),
				c.String("log-fmt"),
				c.String("log-tags"),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			dryRun := c.Bool("dry-run")
			if dryRun {
				logger.Info("DRY-RUN mode: will index events and compute miles but skip Fuel submission and processed marking")
			} else {
				if c.String("fuel-api-url") == "" || c.String("fuel-api-key") == "" {
					return fmt.Errorf("-fuel-api-url and -fuel-api-key are required in production mode")
				}
			}

			ctx, cancel := context.WithCancel(c.Context)
			defer cancel()

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigc
				logger.Info("shutting down...")
				cancel()
			}()

			client, err := ethclient.Dial(c.String("l1-rpc-url"))
			if err != nil {
				return fmt.Errorf("ethclient.Dial: %w", err)
			}

			settlementAddr := common.HexToAddress(c.String("contract"))
			weth := common.HexToAddress(c.String("weth"))

			fastSettlement, err := fastsettlement.NewFastsettlementv3Caller(settlementAddr, client)
			if err != nil {
				return fmt.Errorf("failed to bind FastSettlementV3: %w", err)
			}

			executorAddr, err := fastSettlement.Executor(&bind.CallOpts{Context: ctx})
			if err != nil {
				return fmt.Errorf("failed to fetch executor from contract: %w", err)
			}

			logger.Info("contract parameters initialized",
				slog.String("executor", executorAddr.Hex()),
				slog.String("settlement", settlementAddr.Hex()),
			)

			cfg := &serviceConfig{
				Logger:         logger,
				WETH:           weth,
				FuelURL:        c.String("fuel-api-url"),
				FuelKey:        c.String("fuel-api-key"),
				BarterURL:      c.String("barter-url"),
				BarterKey:      c.String("barter-api-key"),
				Client:         client,
				L1Client:       client,
				ExecutorAddr:   executorAddr,
				DryRun:         dryRun,
				FastSwapURL:    c.String("fastswap-url"),
				FundsRecipient: common.HexToAddress(c.String("funds-recipient")),
				SettlementAddr: settlementAddr,
				MaxGasGwei:     c.Uint64("max-gas-gwei"),
				HTTPClient: &http.Client{
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
				},
			}

			if !dryRun {
				if c.String("keystore") == "" || c.String("passphrase") == "" {
					return fmt.Errorf("-keystore and -passphrase are required for production token sweeping")
				}
				signer, err := loadKeystoreFile(c.String("keystore"), c.String("passphrase"))
				if err != nil {
					return fmt.Errorf("failed to load keystore: %w", err)
				}
				signerAddr := signer.GetAddress()
				if signerAddr != executorAddr {
					logger.Warn("loaded keystore address does NOT match the contract executor",
						slog.String("keystore_address", signerAddr.Hex()),
						slog.String("contract_executor", executorAddr.Hex()))
				}
				logger.Info("loaded executor wallet", slog.String("address", signerAddr.Hex()))
				cfg.Signer = signer
			}

			db := openDB(logger, c.String("db-user"), c.String("db-pw"), c.String("db-host"), c.String("db-port"), c.String("db-name"))
			defer func() { _ = db.Close() }()
			cfg.DB = db

			filterer, err := fastsettlement.NewFastsettlementv3Filterer(settlementAddr, client)
			if err != nil {
				return fmt.Errorf("NewFastsettlementv3Filterer: %w", err)
			}

			startBlock := c.Uint64("start-block")
			if startBlock == 0 {
				startBlock = loadLastBlock(db)
				if startBlock == 0 {
					startBlock = 21746973 // Roughly FastSettlementV3 deployment block
				} else {
					startBlock++
				}
			}

			batchSize := c.Uint64("batch")
			ticker := time.NewTicker(c.Duration("poll"))
			defer ticker.Stop()

			firstRun := make(chan struct{}, 1)
			firstRun <- struct{}{}

			for {
				select {
				case <-ctx.Done():
					return nil
				case <-firstRun:
				case <-ticker.C:
				}

				head, err := client.BlockNumber(ctx)
				if err != nil {
					logger.Error("failed to get latest block", slog.Any("error", err))
					continue
				}

				if startBlock > head {
					continue
				}

				for startBlock <= head {
					endBlock := startBlock + batchSize - 1
					if endBlock > head {
						endBlock = head
					}

					indexed, err := indexBatch(ctx, logger, filterer, client, db, weth, startBlock, endBlock)
					if err != nil {
						logger.Error("indexBatch error",
							slog.Uint64("start", startBlock),
							slog.Uint64("end", endBlock),
							slog.Any("error", err))
						break
					}

					if indexed > 0 {
						logger.Info("indexed intents",
							slog.Uint64("start", startBlock),
							slog.Uint64("end", endBlock),
							slog.Int("count", indexed))
					}

					saveLastBlock(db, endBlock)
					startBlock = endBlock + 1
				}

				// Only process miles when we've caught up to chain tip to avoid
				// hammering the Barter API while still indexing old blocks.
				if startBlock > head {
					ethProcessed, err := processMiles(ctx, cfg)
					if err != nil {
						logger.Error("processMiles error", slog.Any("error", err))
					}
					if ethProcessed > 0 {
						logger.Info("processed ETH/WETH miles", slog.Int("count", ethProcessed))
					}

					erc20Processed, err := processERC20Miles(ctx, cfg)
					if err != nil {
						logger.Error("processERC20Miles error", slog.Any("error", err))
					}
					if erc20Processed > 0 {
						logger.Info("processed ERC20 token sweeps", slog.Int("count", erc20Processed))
					}
				}
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

// -------------------- DB --------------------

func openDB(logger *slog.Logger, dbUser, dbPass, dbHost, dbPort, dbName string) *sql.DB {
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

// saveLastBlock persists the indexer's progress marker. The prior
// DELETE-then-INSERT implementation was not atomic: a pod crash, SIGTERM
// during rolling deploy, or transient DB error between the two statements
// could vanish the `last_block` row. On the next startup loadLastBlock would
// return 0, startBlock would fall back to the contract deployment block,
// and the indexer would re-walk all history — re-inserting every event and
// (before the insertEvent existence guard landed) wiping processed=true on
// every row, causing mass re-submission to Fuel. This was the underlying
// trigger for the 2026-04-16 double-credit incident.
//
// fastswap_miles_meta has PRIMARY KEY(k), so a plain INSERT is an atomic
// upsert under StarRocks PK semantics. The DELETE is unnecessary and unsafe.
func saveLastBlock(db *sql.DB, block uint64) {
	_, err := db.Exec(`INSERT INTO mevcommit_57173.fastswap_miles_meta (k, v) VALUES ('last_block', ?)`, fmt.Sprintf("%d", block))
	if err != nil {
		log.Printf("saveLastBlock: %v", err)
	}
}

// markProcessed sets processed=true and populates the derived columns. It
// intentionally does NOT touch fuel_submitted_at so that re-runs of the
// pipeline (e.g. a row flipped back to processed=false by a reset SQL) can
// rebuild the derived state without appearing to re-credit Fuel.
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

// markProcessedWithFuelSubmission is called only when submitToFuel has just
// succeeded for this row. It sets fuel_submitted_at so future pipeline runs
// (even after a reset of `processed`) skip the submitToFuel call and don't
// double-credit the user.
func markProcessedWithFuelSubmission(db *sql.DB, txHash string, surplusEth, netProfitEth float64, miles int64, bidCost string) {
	_, err := db.Exec(`
UPDATE mevcommit_57173.fastswap_miles
SET surplus_eth = ?, net_profit_eth = ?, miles = ?, bid_cost = ?, processed = true,
    fuel_submitted_at = CURRENT_TIMESTAMP
WHERE tx_hash = ?`,
		surplusEth, netProfitEth, miles, bidCost, txHash)
	if err != nil {
		log.Printf("markProcessedWithFuelSubmission %s: %v", txHash, err)
	}
}

// -------------------- Barter API --------------------

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
	MinReturnFraction float64 `json:"minReturnFraction"`
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

func weiToEth(wei *big.Int) float64 {
	f := new(big.Float).SetInt(wei)
	eth, _ := f.Quo(f, new(big.Float).SetFloat64(1e18)).Float64()
	return eth
}

func padTo32(n *big.Int) []byte {
	b := n.Bytes()
	if len(b) >= 32 {
		return b[:32]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

func padTo32Address(addr common.Address) []byte {
	padded := make([]byte, 32)
	copy(padded[12:], addr.Bytes())
	return padded
}

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

	for i := range keyjson {
		keyjson[i] = 0
	}

	return signer, nil
}
