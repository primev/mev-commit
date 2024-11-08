package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionL1RPCURL = &cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL of the L1 RPC server",
		EnvVars: []string{"MOCK_RELAY_L1_RPC_URL"},
	}

	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen on for HTTP requests",
		EnvVars: []string{"MOCK_RELAY_HTTP_PORT"},
		Value:   8080,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MOCK_RELAY_LOG_FMT"},
		Value:   "text",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"MOCK_RELAY_LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"MOCK_RELAY_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	}
)

func main() {
	app := &cli.App{
		Name:  "relay-emulator",
		Usage: "Emulates the required relay APIs",
		Flags: []cli.Flag{
			optionL1RPCURL,
			optionHTTPPort,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			l1RPC, err := ethclient.Dial(c.String(optionL1RPCURL.Name))
			if err != nil {
				return err
			}

			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			var (
				registeredKeys []string
				registeredLock sync.RWMutex
			)

			mux := http.NewServeMux()

			mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
				// Get BLS keys from request body
				var keys []string
				if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
					http.Error(w, "Failed to decode request body", http.StatusBadRequest)
					return
				}

				if len(keys) == 0 {
					http.Error(w, "No BLS keys provided", http.StatusBadRequest)
					return
				}

				// Validate BLS keys
				for _, key := range keys {
					keyBytes, err := hex.DecodeString(key)
					if err != nil {
						http.Error(w, "Invalid BLS key format", http.StatusBadRequest)
						return
					}
					if len(keyBytes) != 48 {
						http.Error(w, "BLS key must be 48 bytes", http.StatusBadRequest)
						return
					}
				}

				// Register BLS keys
				registeredLock.Lock()
				defer registeredLock.Unlock()

				registeredKeys = append(registeredKeys, keys...)
				logger.Info("Registered BLS keys", "keys", keys)

				w.WriteHeader(http.StatusOK)
			})

			mux.HandleFunc("GET /relay/v1/data/bidtraces/proposer_payload_delivered", func(w http.ResponseWriter, r *http.Request) {
				blockNumberStr := r.URL.Query().Get("block_number")
				if blockNumberStr == "" {
					http.Error(w, "Missing block_number parameter", http.StatusBadRequest)
					return
				}

				blockNumber, err := strconv.ParseUint(blockNumberStr, 10, 64)
				if err != nil {
					http.Error(w, "Invalid block_number", http.StatusBadRequest)
					return
				}

				block, err := l1RPC.BlockByNumber(r.Context(), big.NewInt(int64(blockNumber)))
				if err != nil {
					http.Error(w, "Failed to get block", http.StatusInternalServerError)
					return
				}

				idx := int(blockNumber) % len(registeredKeys)
				registeredLock.RLock()
				key := registeredKeys[idx]
				registeredLock.RUnlock()

				response := []map[string]interface{}{
					{
						"block_number":   blockNumberStr,
						"block_hash":     block.Hash().Hex(),
						"builder_pubkey": key,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			})

			server := &http.Server{
				Addr:    fmt.Sprintf(":%d", c.Int(optionHTTPPort.Name)),
				Handler: mux,
			}

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			go func() {
				logger.Info("Starting server", "port", c.Int(optionHTTPPort.Name))
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Failed to start server", "error", err)
				}
			}()

			<-ctx.Done()
			return server.Shutdown(context.Background())
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
