package api

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primev/mev-commit/p2p/pkg/apiserver"
	"github.com/primev/mev-commit/tools/instant-bridge/bidder"
	"github.com/primev/mev-commit/tools/instant-bridge/transfer"
	"github.com/primev/mev-commit/x/health"
)

type API struct {
	logger           *slog.Logger
	mux              *http.ServeMux
	port             int
	srv              *http.Server
	health           health.Health
	bidder           *bidder.BidderClient
	transferer       *transfer.Transferer
	minServiceFee    *big.Int
	status           *status
	owner            common.Address
	l1Client         *ethclient.Client
	settlementClient *ethclient.Client
}

type bid struct {
	BridgeAmount string `json:"bridgeAmount"`
	RawTx        string `json:"rawTx"`
	DestAddress  string `json:"destAddress"`
}

type status struct {
	bidsAttempted      atomic.Int64
	bidsSucceeded      atomic.Int64
	transfersAttempted atomic.Int64
	transfersSucceeded atomic.Int64
	bridgedAmount      atomic.Pointer[big.Int]
	bidAmountSpent     atomic.Pointer[big.Int]
	feesAccumulated    atomic.Pointer[big.Int]
}

func NewAPI(
	logger *slog.Logger,
	port int,
	health health.Health,
	bdr *bidder.BidderClient,
	transferer *transfer.Transferer,
	minServiceFee *big.Int,
	owner common.Address,
	l1Client *ethclient.Client,
	settlementClient *ethclient.Client,
) *API {
	a := &API{
		logger:           logger,
		mux:              http.NewServeMux(),
		port:             port,
		status:           &status{},
		health:           health,
		bidder:           bdr,
		transferer:       transferer,
		minServiceFee:    minServiceFee,
		owner:            owner,
		l1Client:         l1Client,
		settlementClient: settlementClient,
	}

	a.status.bridgedAmount.Store(big.NewInt(0))
	a.status.bidAmountSpent.Store(big.NewInt(0))
	a.status.feesAccumulated.Store(big.NewInt(0))

	a.mux.HandleFunc("GET /health", func(w http.ResponseWriter, req *http.Request) {
		err := a.health.Health()
		if err != nil {
			apiserver.WriteError(w, http.StatusServiceUnavailable, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "ok\n")
		if err != nil {
			a.logger.Error(
				"failed to write response",
				"error", err,
			)
		}
	})

	a.mux.HandleFunc("GET /estimate", func(w http.ResponseWriter, req *http.Request) {
		estimation, err := a.bidder.Estimate()
		if err != nil {
			apiserver.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		if err := apiserver.WriteResponse(w, http.StatusOK, struct {
			Seconds     int64  `json:"seconds"`
			Cost        string `json:"cost"`
			Destination string `json:"destination"`
		}{
			Seconds:     estimation,
			Cost:        a.minServiceFee.String(),
			Destination: a.owner.Hex(),
		}); err != nil {
			a.logger.Error("failed to write response", "error", err)
		}
	})

	a.mux.HandleFunc("GET /status", func(w http.ResponseWriter, req *http.Request) {
		bridgedAmount := a.status.bridgedAmount.Load()
		bidAmountSpent := a.status.bidAmountSpent.Load()
		feesAccumulated := a.status.feesAccumulated.Load()

		l1Balance, err := a.l1Client.BalanceAt(req.Context(), a.owner, nil)
		if err != nil {
			apiserver.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		settlementBalance, err := a.settlementClient.BalanceAt(req.Context(), a.owner, nil)
		if err != nil {
			apiserver.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if err := apiserver.WriteResponse(w, http.StatusOK, struct {
			BidsAttempted      int64  `json:"bidsAttempted"`
			BidsSucceeded      int64  `json:"bidsSucceeded"`
			TransfersAttempted int64  `json:"transfersAttempted"`
			TransfersSucceeded int64  `json:"transfersSucceeded"`
			BridgedAmount      string `json:"bridgedAmount"`
			BidAmountSpent     string `json:"bidAmountSpent"`
			FeesAccumulated    string `json:"feesAccumulated"`
			L1Balance          string `json:"l1Balance"`
			SettlementBalance  string `json:"settlementBalance"`
		}{
			BidsAttempted:      a.status.bidsAttempted.Load(),
			BidsSucceeded:      a.status.bidsSucceeded.Load(),
			TransfersAttempted: a.status.transfersAttempted.Load(),
			TransfersSucceeded: a.status.transfersSucceeded.Load(),
			BridgedAmount:      bridgedAmount.String(),
			BidAmountSpent:     bidAmountSpent.String(),
			FeesAccumulated:    feesAccumulated.String(),
			L1Balance:          l1Balance.String(),
			SettlementBalance:  settlementBalance.String(),
		}); err != nil {
			a.logger.Error("failed to write response", "error", err)
		}
	})

	a.mux.HandleFunc("POST /bid", func(w http.ResponseWriter, req *http.Request) {
		b, err := apiserver.BindJSON[bid](w, req)
		if err != nil {
			apiserver.WriteError(w, http.StatusBadRequest, err)
			return
		}

		if b.RawTx == "" || b.BridgeAmount == "" {
			apiserver.WriteError(w, http.StatusBadRequest, errors.New("missing fields"))
			return
		}

		tx, err := a.transferer.ValidateL1Tx(b.RawTx)
		if err != nil {
			apiserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid raw tx: %w", err))
			return
		}

		bridgeAmt, ok := new(big.Int).SetString(b.BridgeAmount, 10)
		if !ok {
			apiserver.WriteError(w, http.StatusBadRequest, errors.New("invalid bridge amount"))
			return
		}

		minCost := new(big.Int).Add(bridgeAmt, a.minServiceFee)
		if tx.Value().Cmp(minCost) < 0 {
			diff := new(big.Int).Sub(minCost, tx.Value())
			apiserver.WriteError(
				w,
				http.StatusBadRequest,
				fmt.Errorf("insufficient funds; short by %s", diff.String()),
			)
			return
		}

		fees := new(big.Int).Sub(tx.Value(), bridgeAmt)
		halfFee := big.NewInt(0).Div(fees, big.NewInt(2))

		var destAddr common.Address
		if b.DestAddress == "" {
			destAddr, err = a.transferer.Sender(tx)
			if err != nil {
				apiserver.WriteError(
					w,
					http.StatusBadRequest,
					fmt.Errorf("failed to identify sender: %w", err),
				)
				return
			}
		} else {
			destAddr = common.HexToAddress(b.DestAddress)
		}

		a.status.bidsAttempted.Add(1)
		statusC, err := a.bidder.Bid(
			req.Context(),
			halfFee,
			bridgeAmt,
			b.RawTx,
		)
		if err != nil {
			apiserver.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		for status := range statusC {
			switch status.Type {
			case bidder.BidStatusNoOfProviders:
				a.logger.Info("no of providers", "count", status.Arg1)
			case bidder.BidStatusWaitSecs:
				a.logger.Info("waiting for next slot", "seconds", status.Arg1)
			case bidder.BidStatusAttempted:
				a.logger.Info("bid attempted", "block", status.Arg1)
			case bidder.BidStatusFailed:
				apiserver.WriteError(
					w,
					http.StatusInternalServerError,
					fmt.Errorf("bid failed: %s", status.Arg2),
				)
				return
			case bidder.BidStatusSucceeded:
				a.logger.Info("bid succeeded", "block", status.Arg1)
			}
		}

		a.status.bidsSucceeded.Add(1)

		a.status.transfersAttempted.Add(1)
		err = a.transferer.TransferOnSettlement(
			req.Context(),
			destAddr,
			bridgeAmt,
		)
		if err != nil {
			apiserver.WriteError(w, http.StatusInternalServerError, err)
			return
		}
		a.status.transfersSucceeded.Add(1)
		a.status.bridgedAmount.Store(new(big.Int).Add(a.status.bridgedAmount.Load(), bridgeAmt))
		a.status.bidAmountSpent.Store(new(big.Int).Add(a.status.bidAmountSpent.Load(), halfFee))
		a.status.feesAccumulated.Store(new(big.Int).Add(a.status.feesAccumulated.Load(), halfFee))

		if err := apiserver.WriteResponse(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{
			Message: "success",
		}); err != nil {
			a.logger.Error("failed to write response", "error", err)
		}
	})

	return a
}

func (a *API) Start() {
	a.srv = &http.Server{
		Addr: fmt.Sprintf(":%d", a.port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			recorder := &responseStatusRecorder{ResponseWriter: w}

			start := time.Now()
			a.mux.ServeHTTP(recorder, req)
			a.logger.Info(
				"api access",
				slog.Int("http_status", recorder.status),
				slog.String("http_method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Duration("duration", time.Since(start)),
			)
		}),
	}

	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}()
}

func (a *API) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.srv.Shutdown(ctx)
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
