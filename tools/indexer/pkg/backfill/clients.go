package backfill

import (
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// RatedAPIClient handles communication with Rated Network API
type RatedAPIClient struct {
	httpc   *retryablehttp.Client
	apiKey  string
	baseURL string
}

// NewRatedAPIClient creates a new Rated API client
func NewRatedAPIClient(httpc *retryablehttp.Client, apiKey string) *RatedAPIClient {
	return &RatedAPIClient{
		httpc:   httpc,
		apiKey:  apiKey,
		baseURL: "https://api.rated.network/v1/eth",
	}
}

type QuickNodeClient struct {
	httpc      *retryablehttp.Client
	base       string
	chunkSize  int // ids per call; default 300
	concurrent int // concurrent calls; default 6
	timeout    time.Duration
}

func NewQuickNodeClient(httpc *retryablehttp.Client, base string) *QuickNodeClient {
	base = strings.TrimRight(base, "/")
	return &QuickNodeClient{
		httpc:      httpc,
		base:       base,
		chunkSize:  300, // safe default; increase if URLs stay < 8–10KB
		concurrent: 6,   // keep it polite
		timeout:    8 * time.Second,
	}
}
