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
	api "github.com/primev/mev-commit/tools/validators-monitor/api"
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

func TestNewSlackNotifier(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier("", logger)
	if n.enabled {
		t.Errorf("expected notifier disabled when webhookURL is empty")
	}
	n2 := NewSlackNotifier("http://example.com", logger)
	if !n2.enabled {
		t.Errorf("expected notifier enabled when webhookURL provided")
	}
}

func TestSendMessage_Disabled(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier("", logger)
	if err := n.SendMessage(context.Background(), SlackMessage{Text: "hello"}); err != nil {
		t.Errorf("SendMessage disabled = error %v; want nil", err)
	}
}

func TestSendMessage_Success(t *testing.T) {
	var received []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST; got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json; got %s", ct)
		}
		body, _ := io.ReadAll(r.Body)
		received = body
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier(server.URL, logger)
	msg := SlackMessage{Text: "test"}
	if err := n.SendMessage(context.Background(), msg); err != nil {
		t.Fatalf("SendMessage = %v; want nil", err)
	}

	var got SlackMessage
	if err := json.Unmarshal(received, &got); err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}
	if got.Text != msg.Text {
		t.Errorf("Sent message Text = %q; want %q", got.Text, msg.Text)
	}
}

func TestSendMessage_NonOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier(server.URL, logger)
	err := n.SendMessage(context.Background(), SlackMessage{Text: "err"})
	if err == nil || !strings.Contains(err.Error(), "non-200") {
		t.Errorf("expected non-200 status error; got %v", err)
	}
}

func TestNotifyRelayData_NoDashboard(t *testing.T) {
	var payload SlackMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier(server.URL, logger)
	relays := []string{"r1", "r2"}
	allRelays := []string{"r1", "r2", "r3"}
	if err := n.NotifyRelayData(context.Background(), "0xabc", 1, 100, 10, big.NewInt(2e18), "0xabc", relays, allRelays, nil); err != nil {
		t.Fatalf("NotifyRelayData = %v; want nil", err)
	}
	if len(payload.Attachments) != 1 {
		t.Fatalf("expected 1 attachment; got %d", len(payload.Attachments))
	}
	att := payload.Attachments[0]
	if att.Color != "#36a64f" {
		t.Errorf("Color = %q; want %q", att.Color, "#36a64f")
	}
	if att.Title != "Relay Data Available for Validator" {
		t.Errorf("Title = %q; want %q", att.Title, "Relay Data Available for Validator")
	}
	fields := map[string]string{}
	for _, f := range att.Fields {
		fields[f.Title] = f.Value
	}
	wantFields := map[string]string{
		"Validator Index":   "1",
		"Slot":              "10",
		"Block Number":      "100",
		"Validator Pubkey":  "0xabc",
		"Relays With Data":  "```- r1\n- r2\n```",
		"Data Availability": "2 of 3 relays have data",
	}
	for k, v := range wantFields {
		if got, ok := fields[k]; !ok {
			t.Errorf("missing field %q", k)
		} else if got != v {
			t.Errorf("field %q = %q; want %q", k, got, v)
		}
	}
}

func TestNotifyRelayData_WithDashboard(t *testing.T) {
	var payload SlackMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	n := NewSlackNotifier(server.URL, logger)
	relays := []string{}
	all := []string{"a", "b", "c"}
	info := &api.DashboardResponse{
		Winner:                 "relayX",
		TotalOpenedCommitments: 5,
		TotalRewards:           8,
		TotalSlashes:           1,
		TotalAmount:            "5000000000000000000",
	}
	if err := n.NotifyRelayData(context.Background(), "0xdef", 2, 200, 20, big.NewInt(3e18), "0xabc", relays, all, info); err != nil {
		t.Fatalf("NotifyRelayData = %v; want nil", err)
	}
	if len(payload.Attachments) != 1 {
		t.Fatalf("expected 1 attachment; got %d", len(payload.Attachments))
	}
	att := payload.Attachments[0]
	if att.Color != "#ff9900" {
		t.Errorf("Color = %q; want %q", att.Color, "#ff9900")
	}
	if att.Title != "No Relay Data Found for Validator" {
		t.Errorf("Title = %q; want %q", att.Title, "No Relay Data Found for Validator")
	}
	fields := map[string]string{}
	for _, f := range att.Fields {
		fields[f.Title] = f.Value
	}
	wantFields := map[string]string{
		"Validator Index":      "2",
		"Slot":                 "20",
		"Block Number":         "200",
		"Validator Pubkey":     "0xdef",
		"Relays With Data":     "```None```",
		"Data Availability":    "0 of 3 relays have data",
		"Block Winner":         "relayX",
		"Commitments":          "5 (Rewards: 8, Slashes: 1)",
		"Total Bid Amount":     "5.000000 ETH",
		"MEV Reward":           "3.000000 ETH",
		"MEV Reward Recipient": "0xabc",
	}
	for k, v := range wantFields {
		if got, ok := fields[k]; !ok {
			t.Errorf("missing field %q", k)
		} else if got != v {
			t.Errorf("field %q = %q; want %q", k, got, v)
		}
	}
}
