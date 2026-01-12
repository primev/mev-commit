package explorersubmitter

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/util"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	reqChan := make(chan *http.Request, 1)

	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			if q.Get("apikey") == "" {
				http.Error(w, "missing apikey", http.StatusBadRequest)
				return
			}
			if q.Get("action") != "submitTxPending" {
				http.Error(w, "wrong action", http.StatusBadRequest)
				return
			}

			jsonData := q.Get("JsonData")
			if jsonData == "" {
				http.Error(w, "missing JsonData", http.StatusBadRequest)
				return
			}

			var txData TxData
			if err := json.Unmarshal([]byte(jsonData), &txData); err != nil {
				http.Error(w, "invalid json data", http.StatusBadRequest)
				return
			}

			if txData.ChainID != "1" {
				http.Error(w, "wrong chain id", http.StatusBadRequest)
				return
			}

			reqChan <- r
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer srv.Close()

	endpoint := srv.URL
	apiKey := "test-api-key"
	appCode := "test-app-code"

	logger := util.NewTestLogger(os.Stdout)
	submitter := New(endpoint, apiKey, appCode, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := submitter.Start(ctx)

	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	err := submitter.Submit(ctx, tx, common.Address{})
	require.NoError(t, err)

	select {
	case <-reqChan:
		// Request received by server.
		// Give a moment for the server to reply and the client to process the 200 OK
		// preventing "context canceled" error in logs.
		time.Sleep(50 * time.Millisecond)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for request")
	}

	cancel()
	<-done
}
