package notifier

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestFormatWeiToEth(t *testing.T) {
	cases := []struct {
		wei  *big.Int
		want string
	}{
		{big.NewInt(0), "0.000000 ETH"},
		{big.NewInt(params.Ether), "1.000000 ETH"},
		{big.NewInt(1234567890000000000), "1.234568 ETH"},
	}
	for _, c := range cases {
		if got := formatWeiToEth(c.wei); got != c.want {
			t.Errorf("formatWeiToEth(%v) = %q; want %q", c.wei, got, c.want)
		}
	}
}

func TestSendMessage_Success(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		received = body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewNotifier([]string{server.URL}, logger)
	msg := Message{Text: "test"}
	require.NoError(t, n.SendMessage(context.Background(), msg))

	var got Message
	require.NoError(t, json.Unmarshal(received, &got))
	require.Equal(t, msg.Text, got.Text)
}

func TestSendMessage_NonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewNotifier([]string{server.URL}, logger)
	err := n.SendMessage(context.Background(), Message{Text: "err"})
	require.Error(t, err)
}
