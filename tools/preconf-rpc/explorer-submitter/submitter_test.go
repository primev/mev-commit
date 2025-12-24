package explorersubmitter

import (
	"context"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	explorerEndpoint = flag.String("explorer-endpoint", "", "Explorer API endpoint")
	explorerApiKey   = flag.String("explorer-apikey", "", "Explorer API Key")
	explorerAppCode  = flag.String("explorer-appcode", "", "Explorer App Code")
)

func TestSubmit(t *testing.T) {
	if *explorerEndpoint == "" || *explorerApiKey == "" || *explorerAppCode == "" {
		t.Skip("skipping integration test, flags not provided")
	}

	config := Config{
		Endpoint: *explorerEndpoint,
		ApiKey:   *explorerApiKey,
		AppCode:  *explorerAppCode,
	}

	err := Submit(context.Background(), config, "1", "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", "0x123", "0x456")
	require.NoError(t, err)
}
