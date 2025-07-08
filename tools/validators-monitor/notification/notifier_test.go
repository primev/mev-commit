package notification

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestFormatRelayList(t *testing.T) {
	if got := formatRelayList([]string{}); got != "None" {
		t.Errorf("formatRelayList(empty) = %q; want \"None\"", got)
	}
	relays := []string{"relay1", "relay2"}
	want := "- relay1\n- relay2\n"
	if got := formatRelayList(relays); got != want {
		t.Errorf("formatRelayList(%v) = %q; want %q", relays, got, want)
	}
}

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

func TestNewNotifier(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewNotifier([]string{}, logger)
	if n.enabled {
		t.Errorf("expected notifier disabled when webhookURLs are empty")
	}
	n2 := NewNotifier([]string{"http://example.com"}, logger)
	if !n2.enabled {
		t.Errorf("expected notifier enabled when webhookURLs provided")
	}
}

func TestSendMessage_Disabled(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewNotifier([]string{}, logger)
	if err := n.SendMessage(context.Background(), Message{Text: "hello"}); err != nil {
		t.Errorf("SendMessage disabled = error %v; want nil", err)
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
	require.True(t, strings.Contains(err.Error(), "one or more notifications failed"))
}

func TestNotifyRelayData(t *testing.T) {
	var payload Message
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewNotifier([]string{server.URL}, logger)
	relays := []string{"relay1"}
	allRelays := []string{"relay1", "relay2"}

	err := n.NotifyRelayData(context.Background(), "0xabc", 123, 456, 789, big.NewInt(2e18), "0xfee", relays, allRelays, nil, "0xbuilder")
	require.NoError(t, err)

	require.Len(t, payload.Attachments, 1)
	attachment := payload.Attachments[0]
	require.Equal(t, "#36a64f", attachment.Color)
	require.Equal(t, "Relay Data Available for Validator", attachment.Title)
}
