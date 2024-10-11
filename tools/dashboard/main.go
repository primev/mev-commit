package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	oracle "github.com/primev/mev-commit/contracts-abi/clients/Oracle"
	preconf "github.com/primev/mev-commit/contracts-abi/clients/PreconfManager"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
	"github.com/primev/mev-commit/x/contracts/events"
	"github.com/primev/mev-commit/x/contracts/events/publisher"
	"github.com/primev/mev-commit/x/util"
	"github.com/urfave/cli/v2"
)

var (
	optionRPCURL = &cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL of the Settlement RPC server",
		EnvVars: []string{"DASH_CLI_RPC_URL"},
		Value:   "wss://chainrpc-wss.testnet.mev-commit.xyz",
	}

	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "port for the HTTP server",
		EnvVars: []string{"DASH_HTTP_PORT"},
		Value:   8080,
	}

	optionStartBlock = &cli.IntFlag{
		Name:    "start-block",
		Usage:   "start block for reading the events for the dashboard",
		EnvVars: []string{"DASH_START_BLOCK"},
		Value:   0,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"DASH_LOG_FMT"},
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
		EnvVars: []string{"DASH_LOG_LEVEL"},
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
		EnvVars: []string{"DASH_LOG_TAGS"},
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

type progStore struct {
	startBlock uint64
}

func (ps *progStore) LastBlock() (uint64, error) {
	return ps.startBlock, nil
}

func (ps *progStore) SetLastBlock(block uint64) error {
	ps.startBlock = block
	return nil
}

func main() {
	app := &cli.App{
		Name:  "mev-commit-dashboard",
		Usage: "MEV Commit Dashboard",
		Flags: []cli.Flag{
			optionRPCURL,
			optionHTTPPort,
			optionStartBlock,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			abis, err := getContractABIs()
			if err != nil {
				return err
			}

			settlementClient, err := ethclient.Dial(c.String(optionRPCURL.Name))
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

			evtMgr := events.NewListener(
				logger,
				abis...,
			)

			pb := publisher.NewWSPublisher(
				&progStore{startBlock: uint64(c.Int(optionStartBlock.Name))},
				logger,
				settlementClient,
				evtMgr,
			)

			statHdlr, err := newStatHandler(evtMgr, 10)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			pbStopped := pb.Start(ctx)

			mux := http.NewServeMux()
			registerRoutes(mux, statHdlr)

			srv := &http.Server{
				Addr: fmt.Sprintf(":%d", c.Int(optionHTTPPort.Name)),
				Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					recorder := &responseStatusRecorder{ResponseWriter: w}

					start := time.Now()
					mux.ServeHTTP(recorder, req)
					logger.Info(
						"api access",
						slog.Int("http_status", recorder.status),
						slog.String("http_method", req.Method),
						slog.String("path", req.URL.Path),
						slog.Duration("duration", time.Since(start)),
					)
				}),
			}

			go func() {
				if err := srv.ListenAndServe(); err != nil {
					fmt.Fprintf(os.Stderr, "error: %v\n", err)
				}
			}()

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-sigc:
			case <-pbStopped:
			}

			cancel()
			statHdlr.close()

			return srv.Shutdown(c.Context)
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

type responseStatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseStatusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Hijack implements http.Hijacker.
func (r *responseStatusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush implements http.Flusher.
func (r *responseStatusRecorder) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Push implements http.Pusher.
func (r *responseStatusRecorder) Push(target string, opts *http.PushOptions) error {
	return r.ResponseWriter.(http.Pusher).Push(target, opts)
}

func registerRoutes(mux *http.ServeMux, statHdlr *statHandler) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if !statHdlr.healthy() {
			http.Error(w, "not healthy", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /dashboard", func(w http.ResponseWriter, r *http.Request) {
		page, limit := parsePagination(r)

		dout := statHdlr.getDashboard(page, limit)
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /windows", func(w http.ResponseWriter, r *http.Request) {
		page, limit := parsePagination(r)

		dout := statHdlr.getWindows(page, limit)
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /window/{window}", func(w http.ResponseWriter, r *http.Request) {
		windowStr := r.PathValue("window")
		window, err := strconv.Atoi(windowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dout := statHdlr.getWindowStat(uint64(window))
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /blocks", func(w http.ResponseWriter, r *http.Request) {
		page, limit := parsePagination(r)

		dout := statHdlr.getBlocks(page, limit)
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /block/{block}", func(w http.ResponseWriter, r *http.Request) {
		blockStr := r.PathValue("block")
		block, err := strconv.Atoi(blockStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dout := statHdlr.getBlockStats(uint64(block))
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /providers", func(w http.ResponseWriter, r *http.Request) {
		dout := statHdlr.getProviders()
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /bidders", func(w http.ResponseWriter, r *http.Request) {
		dout := statHdlr.getCurrentBidders()
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /bidders/{window}", func(w http.ResponseWriter, r *http.Request) {
		windowStr := r.PathValue("window")
		window, err := strconv.Atoi(windowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dout := statHdlr.getBidders(window)
		if dout == nil {
			http.Error(w, "no data", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(dout); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}

func parsePagination(r *http.Request) (int, int) {
	limit := 10
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil {
			limit = l
		}
	}

	page := 0
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil {
			page = p
		}
	}
	return page, limit
}

func getContractABIs() ([]*abi.ABI, error) {
	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		return nil, err
	}

	pcABI, err := abi.JSON(strings.NewReader(preconf.PreconfmanagerABI))
	if err != nil {
		return nil, err
	}

	bidderRegistry, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		return nil, err
	}

	providerRegistry, err := abi.JSON(strings.NewReader(providerregistry.ProviderregistryABI))
	if err != nil {
		return nil, err
	}

	orABI, err := abi.JSON(strings.NewReader(oracle.OracleABI))
	if err != nil {
		return nil, err
	}

	return []*abi.ABI{
		&btABI,
		&pcABI,
		&bidderRegistry,
		&providerRegistry,
		&orABI,
	}, nil
}
