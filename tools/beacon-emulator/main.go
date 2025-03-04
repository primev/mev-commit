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
	"0xb58458ce256859070ebe2f8a4e49b39c58fc324f68504542c00e707ec1557b34e7f2fbadf51c794a750b88486f2f6ac0",
	"0xa1f97b2adca0919d3bf522983760940986c71d8051f409baec6c82bf3b1590031cacb0f2f8f245dd0e7a83042960cfe9",
	"0x97a70b91181a572486d0fc738de2b88eef8695f462754fddc96b02aa9f328ab11e8a5f7d325b11e345aa73eae57e44d5",
	"0xa9d91e0166a88ced138be77639057cd35563df77a49cc05002b9a51dbeba695a45e6f2c3c75ed67cae6fce23cd2c7d38",
	"0xa016da5e51be142eb53f806ee2a9ea254d4fc3704ca63e5c59ecf35a072c7e17f7d112a38112523520ecf30985b5a28f",
	"0x96aa42c8b3629dcd1c92210eb2bea519d2930488604c65696146135633929a64cb4cd8f44bdeffaec198538f529205e3",
	"0xb2d1be4a38a08f7a3ba5199f05431e5fcfba6a7a9237d10bab325a277ae0b95f15c9a6932402668a36261b7b740e439f",
	"0xb66a9d3f47a15e9e5ba3ddafb200fbc2f6bdcdc7ce2a571ee645d3187fcd227b7962fdd0684e5af0c1952242f12724d3",
	"0x8d81b47f6f7e9f7bad294da15ba733bd72dc03b7794dbc9a2abfadcc6587705867ad052762ccfb4ea5de1ea166b474f3",
	"0x80d92b61ff0e54b98e5584e201069a0e85486e41f942ca9afdc0c6b55de116ae431b88ce785b978d19bba0fe67834d67",
	"0xb62bfbb505413d3b9b29378af684abed48f3eae9dc1cf8f2237262984015ae51145dc2d3f723a1e7628191c59ef58cb5",
	"0xb6155b186308c998099e35ab3501ef35fd55bd5b7b39584c64e2b4631501e341fa19f88053e0be9b9eb3a892dea6961d",
	"0x96dd552cfac2bc3c69b662297ed2f1d6caa2679dc05f08479dd0bbcc3e46426c75e1580b284957206f478549d48d75bf",
	"0xb35639b5d41d6b4f49f59a6882f980b250552e0710aaa2f5e86cf8643af5311bf45044d324645eca220542cd895a300d",
	"0x91c6d679a2bfc6cea04429d6d8301fcb24585ba4e1d3a3398c9cb814d6d544973c804a96b0c139eed51334333e5d4449",
	"0x88d03907c1438696c1fbf58cc9f9f33b3b374d3dd305e343592ea5c896c0a899d9f943df1a62124d397fd7627c79eb1e",
	"0x9690048efcf4eb6bafe89f3a2d81f93cc58e2baafd9f9a7936867fa09935b14b501bf249691ecf7f82483aa9ae9edbd3",
	"0x88528c33657cd9c98f4300fdf9832794addd77af608e5abe798190ab405b2f70cc2e8a055db40e08f495cda388705fbb",
	"0xa931044ef42172f75fe48337ded55ed518de1f92269c0a0599e9ce3e753bc0ca3f25e63a0f735cbbda159ee5c237048f",
	"0x8c5c026664e35118661d0089ccc06d4961d84a0569c3bcf108d75ee5b7e98733e31879bb3ee2e335d2398d5d4f477d95",
	"0xb255287ef2994692772eb18a471dfe5582e438d842c2180e13e83512a65816f307c75a7c5282c7df6fdb5e26b3f30da0",
	"0x83fa44d4bd9dfa7f7c490d09d9364562d8caa161552510a9fcb6279bdb92170611faf9b24bf12d03310f461790c5eeaa",
	"0x8751da86b9fc6ab13f92eeffcff1849e4913e6954b30db49da310f420c5b145395a1d2cf7152464cddb9bea401d6a807",
	"0x962b5c5d46434dcc483fdc50e3ad4faea152b76daf5277c44b839a19ec9132d97074a8890ab4fffad413fdc414bf1478",
	"0x9639113b0d32cadd125f71154a42768194472847fb70e9f76a5d62f88a407a004315fc412248357a49158d7c729b3f7e",
	"0x8e929ab3be234c811cbacb883da8b657a5c05b201edc2ac0e1c1c0b64d9d3a2b3d806382d36001d0a0a2f8bbdbce39ba",
	"0x92a1ecf5e8ba91730277a6c4bce5ec2325476210ae15d815dffb5f25286f5677135ccd76ada30ad2fdf0e1b5634d79ff",
	"0xae6192653543a8fc71e009bc236d6634dc3b1676cc50104847ddf43de83c3276f1c2bc909efefe8bcc2569871a4fafdb",
	"0x84151867989a5dddb47d20bc5afccb098f86dd70db9e6fbe1825201c5d300423dcfa3c1131e5723f576b5b4e168f42ea",
	"0x8c28f96f8dcc8562bc94d74f1a16e87cdb8f064ca9d8aa98333f889af4b0376ef7a26d1b7af63c54c91bae2aff30d288",
	"0x8bd76d3d0c2f26c89b2b834a994599fdfee9e654bd90fdb71e26b958d7f2076d3cb662c0d78af7bcef260cb88fdd548e",
	"0xababbfe729893e69384ef1f32c7fa15902be6ace12aeaa21c56be726bc8c71e4e9b884735b82dbc619315752cffdb73e",
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

			mux.HandleFunc(
				"GET /eth/v1/validator/duties/proposer/{epoch}",
				func(w http.ResponseWriter, r *http.Request) {
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
				},
			)

			mux.HandleFunc(
				"GET /eth/v1/beacon/states/head/finality_checkpoints",
				func(w http.ResponseWriter, r *http.Request) {
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
				},
			)

			mux.HandleFunc(
				"GET /eth/v1/beacon/genesis",
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(struct {
						Data struct {
							GenesisTime          string `json:"genesis_time"`
							GenesisValidatorRoot string `json:"genesis_validators_root"`
							GenesisForkVersion   string `json:"genesis_fork_version"`
						} `json:"data"`
					}{
						Data: struct {
							GenesisTime          string `json:"genesis_time"`
							GenesisValidatorRoot string `json:"genesis_validators_root"`
							GenesisForkVersion   string `json:"genesis_fork_version"`
						}{
							GenesisTime:          "1695902400",
							GenesisValidatorRoot: "0x9143aa7c615a7f7115e2b6aac319c03529df8242ae705fba9df39b79c59fa8b1",
							GenesisForkVersion:   "0x01017000",
						},
					}); err != nil {
						http.Error(w, "Failed to encode response", http.StatusInternalServerError)
						return
					}
					w.WriteHeader(http.StatusOK)
				},
			)

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
