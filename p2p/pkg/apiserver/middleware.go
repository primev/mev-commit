package apiserver

import (
	"bufio"
	"log/slog"
	"net"
	"net/http"
	"time"
)

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

func newAccessLogHandler(log *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			recorder := &responseStatusRecorder{ResponseWriter: w}

			start := time.Now()
			h.ServeHTTP(recorder, req)
			log.Info("api access",
				"http_status", recorder.status,
				"method", req.Method,
				"path", req.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}
