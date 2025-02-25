package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"

	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port to listen on for HTTP requests",
		EnvVars: []string{"MOCK_BEACON_HTTP_PORT"},
		Value:   8080,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"MOCK_BEACON_LOG_FMT"},
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
		EnvVars: []string{"MOCK_BEACON_LOG_LEVEL"},
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
		EnvVars: []string{"MOCK_BEACON_LOG_TAGS"},
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

// Valus hard-coded from current testnet
var registeredKeys = []string{
	"0x8372542ab5380465a370cc19f67b47f5746d889ea3ffb8983b1c176de245c860b660a4716450668da25582bf1ed71d17",
	"0x86d2cc9058f81613cd1a16909cbf337ad0eba54941d5b3236844d6a7f05ad214a2c012a84ed313df2c6ded96f9c4d2a5",
	"0xa8af120741a5f2d4b10fb8544620af734dd51254a368f353c4ba33c4f3e25d60b2c027afec4263ae526d9d8df5cb85fc",
	"0x92173b07178fa53b6b8e0f97b71509a7381ca471fc05a4fba41e810bb3d74a5b8c0209c6bd190b0040b4a51ce69cfc24",
	"0xa9bee05cc34c540973229d2e924c7884f23bb82cb39574a685a5df856f897abaa1a24650a9cccbd9d7dc68731d25652b",
	"0xb8ac5686f7badf23d999f6d052b1332846a9930e4b4df494aa8aa89e1497d3475be5d9f9e7aabc3f7272d3c4b755d098",
	"0x880102538af1165d5c86525158ee1f70d0312c609fbe68452be155b263c2a65213ccc9718a8c616ce2e9b916b1136cc5",
	"0xaf4d6003f9dc818689fba4e0a5ceefc091390c5f72af6c101159a527f92773be24567f58f77c413164f3af2969e3c2d3",
	"0x81e20a12d616f4802b94a978c30e8e1031fd529770215d73b504d9190ca3d2619598289480db554a4ecba496a84788df",
	"0xaad73d536a177db9e310c8236c6c3ee27437920d0c8167de10a066872e672d82b6f821d7b794ffebb8811d3583b353d4",
	"0x92173b07178fa53b6b8e0f97b71509a7381ca471fc05a4fba41e810bb3d74a5b8c0209c6bd190b0040b4a51ce69cfc24",
	"0xa8af120741a5f2d4b10fb8544620af734dd51254a368f353c4ba33c4f3e25d60b2c027afec4263ae526d9d8df5cb85fc",
	"0xb8ac5686f7badf23d999f6d052b1332846a9930e4b4df494aa8aa89e1497d3475be5d9f9e7aabc3f7272d3c4b755d098",
	"0xa9bee05cc34c540973229d2e924c7884f23bb82cb39574a685a5df856f897abaa1a24650a9cccbd9d7dc68731d25652b",
	"0x81e20a12d616f4802b94a978c30e8e1031fd529770215d73b504d9190ca3d2619598289480db554a4ecba496a84788df",
	"0x8372542ab5380465a370cc19f67b47f5746d889ea3ffb8983b1c176de245c860b660a4716450668da25582bf1ed71d17",
	"0xaf4d6003f9dc818689fba4e0a5ceefc091390c5f72af6c101159a527f92773be24567f58f77c413164f3af2969e3c2d3",
	"0x880102538af1165d5c86525158ee1f70d0312c609fbe68452be155b263c2a65213ccc9718a8c616ce2e9b916b1136cc5",
	"0xaad73d536a177db9e310c8236c6c3ee27437920d0c8167de10a066872e672d82b6f821d7b794ffebb8811d3583b353d4",
	"0x86d2cc9058f81613cd1a16909cbf337ad0eba54941d5b3236844d6a7f05ad214a2c012a84ed313df2c6ded96f9c4d2a5",
	"0x81e20a12d616f4802b94a978c30e8e1031fd529770215d73b504d9190ca3d2619598289480db554a4ecba496a84788df",
	"0x880102538af1165d5c86525158ee1f70d0312c609fbe68452be155b263c2a65213ccc9718a8c616ce2e9b916b1136cc5",
	"0x8372542ab5380465a370cc19f67b47f5746d889ea3ffb8983b1c176de245c860b660a4716450668da25582bf1ed71d17",
	"0xa9bee05cc34c540973229d2e924c7884f23bb82cb39574a685a5df856f897abaa1a24650a9cccbd9d7dc68731d25652b",
	"0xa8af120741a5f2d4b10fb8544620af734dd51254a368f353c4ba33c4f3e25d60b2c027afec4263ae526d9d8df5cb85fc",
	"0x86d2cc9058f81613cd1a16909cbf337ad0eba54941d5b3236844d6a7f05ad214a2c012a84ed313df2c6ded96f9c4d2a5",
	"0x92173b07178fa53b6b8e0f97b71509a7381ca471fc05a4fba41e810bb3d74a5b8c0209c6bd190b0040b4a51ce69cfc24",
	"0xb8ac5686f7badf23d999f6d052b1332846a9930e4b4df494aa8aa89e1497d3475be5d9f9e7aabc3f7272d3c4b755d098",
	"0xaf4d6003f9dc818689fba4e0a5ceefc091390c5f72af6c101159a527f92773be24567f58f77c413164f3af2969e3c2d3",
	"0xaad73d536a177db9e310c8236c6c3ee27437920d0c8167de10a066872e672d82b6f821d7b794ffebb8811d3583b353d4",
	"0x8147b3726a49faa0bb034b9947c3e2742881546abdffbd460754cb47d75ed1e5d0af0d6c5f80ac5c8b07078c11c0c06c",
	"0xadf76bf5182e7adae3b8617711bdaef439772754ece921fb9e3684ad68113b0b95da915c7bbe5142fe60bfb50e84185e",
}

type ProposerDutiesResponse struct {
	Data []struct {
		Pubkey string `json:"pubkey"`
		Slot   string `json:"slot"`
	} `json:"data"`
}

type FinalityCheckpointsResponse struct {
	Data struct {
		PreviousJustified struct {
			Epoch string `json:"epoch"`
		} `json:"previous_justified"`
		CurrentJustified struct {
			Epoch string `json:"epoch"`
		} `json:"current_justified"`
		Finalized struct {
			Epoch string `json:"epoch"`
		} `json:"finalized"`
	} `json:"data"`
}

func main() {
	app := &cli.App{
		Name:  "beacon-emulator",
		Usage: "Emulates the required beacon APIs",
		Flags: []cli.Flag{
			optionHTTPPort,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			mux := http.NewServeMux()

			mux.HandleFunc("GET /eth/v1/validator/duties/proposer/{epoch}", func(w http.ResponseWriter, r *http.Request) {
				epochStr := r.PathValue("epoch")
				epoch, err := strconv.Atoi(epochStr)
				if err != nil {
					http.Error(w, "Invalid epoch", http.StatusBadRequest)
					return
				}

				response := &ProposerDutiesResponse{
					Data: make([]struct {
						Pubkey string `json:"pubkey"`
						Slot   string `json:"slot"`
					}, len(registeredKeys)),
				}

				baseSlot := epoch * 32

				for i, key := range registeredKeys {
					response.Data[i] = struct {
						Pubkey string `json:"pubkey"`
						Slot   string `json:"slot"`
					}{
						Pubkey: key,
						Slot:   fmt.Sprintf("%d", baseSlot+i),
					}
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			})

			mux.Handle("GET /eth/v1/beacon/states/head/finality_checkpoints", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				response := &FinalityCheckpointsResponse{
					Data: struct {
						PreviousJustified struct {
							Epoch string `json:"epoch"`
						} `json:"previous_justified"`
						CurrentJustified struct {
							Epoch string `json:"epoch"`
						} `json:"current_justified"`
						Finalized struct {
							Epoch string `json:"epoch"`
						} `json:"finalized"`
					}{
						PreviousJustified: struct {
							Epoch string `json:"epoch"`
						}{
							Epoch: "0",
						},
						CurrentJustified: struct {
							Epoch string `json:"epoch"`
						}{
							Epoch: "0",
						},
						Finalized: struct {
							Epoch string `json:"epoch"`
						}{
							Epoch: "0",
						},
					},
				}

				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))

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
