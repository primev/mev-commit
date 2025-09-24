package http

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}
func FetchJSONWithRetry(ctx context.Context, httpc *http.Client, url string, out any, attempts int, baseDelay time.Duration) error {
	var lastErr error

	for i := 0; i < attempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			continue
		}
		resp, err := httpc.Do(req)
		if err == nil && resp != nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			return json.NewDecoder(resp.Body).Decode(out)
		}
		if resp != nil {
			// 429 courtesy backoff if provided
			if resp.StatusCode == 429 {
				if ra := resp.Header.Get("Retry-After"); ra != "" {
					if secs, err := strconv.Atoi(ra); err == nil {
						select {
						case <-ctx.Done():
							resp.Body.Close()
							return ctx.Err()

						case <-time.After(time.Duration(secs) * time.Second):
						}

					}
				}
			}
			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
			}
			resp.Body.Close()
		} else if err != nil {
			lastErr = err
		}
		if i < attempts-1 {
			sleep := baseDelay * time.Duration(1<<i)
			jitter := time.Duration(rand.Int63n(int64(baseDelay / 2)))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleep + jitter):
			}
		}
	}
	return fmt.Errorf("GET %s failed after %d attempts: %v", url, attempts, lastErr)

}
