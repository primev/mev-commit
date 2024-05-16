package apiserver_test

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/primev/mev-commit/p2p/pkg/apiserver"
)

func newTestLogger(w io.Writer) *slog.Logger {
	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func TestAPIServer(t *testing.T) {
	t.Parallel()

	t.Run("new and close", func(t *testing.T) {
		var logBuf bytes.Buffer
		s := apiserver.New(
			"test",
			newTestLogger(&logBuf),
		)

		srv := httptest.NewServer(s.Router())
		t.Cleanup(func() {
			srv.Close()
		})

		r, err := http.NewRequest("GET", srv.URL+"/metrics", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var b bytes.Buffer
		n, err := b.ReadFrom(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if n == 0 {
			t.Fatal("expected non-zero body")
		}

		if !strings.Contains(b.String(), "test") {
			t.Fatalf("expected body to contain 'test', got %q", b.String())
		}

		if !strings.Contains(b.String(), "go_info") {
			t.Fatalf("expected body to contain 'go_info', got %q", b.String())
		}

		if !strings.Contains(b.String(), "go_memstats") {
			t.Fatalf("expected body to contain 'go_memstats', got %q", b.String())
		}

		if !strings.Contains(b.String(), "go_gc_duration_seconds") {
			t.Fatalf("expected body to contain 'go_gc_duration_seconds', got %q", b.String())
		}

		if !strings.Contains(logBuf.String(), "api access") {
			t.Fatalf("expected log to contain 'api access', got %q", logBuf.String())
		}
	})

	t.Run("chain handlers", func(t *testing.T) {
		var logBuf bytes.Buffer
		s := apiserver.New(
			"test",
			newTestLogger(&logBuf),
		)

		srv := httptest.NewServer(s.Router())
		t.Cleanup(func() {
			srv.Close()
		})

		var orderedHandlerActions []int
		s.ChainHandlers(
			"/chain",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				orderedHandlerActions = append(orderedHandlerActions, 3)
				w.WriteHeader(http.StatusOK)
			}),
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					orderedHandlerActions = append(orderedHandlerActions, 1)
					next.ServeHTTP(w, r)
				})
			},
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					orderedHandlerActions = append(orderedHandlerActions, 2)
					next.ServeHTTP(w, r)
				})
			},
		)

		r, err := http.NewRequest("GET", srv.URL+"/chain", nil)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		if len(orderedHandlerActions) != 3 {
			t.Fatalf("expected 3 handler actions, got %d", len(orderedHandlerActions))
		}

		for i, v := range []int{1, 2, 3} {
			if orderedHandlerActions[i] != v {
				t.Fatalf("expected handler action %d, got %d", v, orderedHandlerActions[i])
			}
		}
	})
}
