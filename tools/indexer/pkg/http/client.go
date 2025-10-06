package http

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func NewHTTPClient(timeout time.Duration) *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.HTTPClient.Timeout = timeout
	client.RetryMax = 3
	client.RetryWaitMin = 200 * time.Millisecond
	client.RetryWaitMax = 2 * time.Second
	client.Logger = nil
	return client
}

func FetchJSON(ctx context.Context, client *retryablehttp.Client, url string, out any) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
