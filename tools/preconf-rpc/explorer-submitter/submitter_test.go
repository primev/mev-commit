package explorersubmitter

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primev/mev-commit/x/util"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	endpoint := os.Getenv("EXPLORER_API_ENDPOINT")
	apiKey := os.Getenv("EXPLORER_API_KEY")
	appCode := os.Getenv("EXPLORER_APPCODE")

	if endpoint == "" || apiKey == "" || appCode == "" {
		t.Skip("skipping integration test, flags not provided")
	}

	logger := util.NewTestLogger(os.Stdout)
	submitter := New(endpoint, apiKey, appCode, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := submitter.Start(ctx)

	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	fmt.Printf("Submitting mock transaction with hash: %s\n", tx.Hash().Hex())
	err := submitter.Submit(ctx, tx, common.Address{})
	require.NoError(t, err)

	// allow processing
	time.Sleep(2 * time.Second)

	cancel()
	<-done
}
