package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/primev/mev-commit/tools/validators-monitor/api"
)

// Message represents a notification message structure
// generalized to support different platforms.
type Message struct {
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a message attachment
type Attachment struct {
	Color      string   `json:"color,omitempty"`
	Title      string   `json:"title,omitempty"`
	Text       string   `json:"text,omitempty"`
	Fields     []Field  `json:"fields,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	TS         int64    `json:"ts,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
}

// Field represents a field in a message attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Notifier sends notifications to multiple webhook endpoints
type Notifier struct {
	webhookURLs []string
	client      *http.Client
	logger      *slog.Logger
	enabled     bool
}

// NewNotifier creates a new notifier instance
func NewNotifier(webhookURLs []string, logger *slog.Logger) *Notifier {
	enabled := len(webhookURLs) > 0

	if !enabled {
		logger.Warn("Notifications disabled - no webhook URLs provided")
	} else {
		logger.Info("Notifications enabled")
	}

	return &Notifier{
		webhookURLs: webhookURLs,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:  logger,
		enabled: enabled,
	}
}

// SendMessage sends a message to all configured webhook endpoints
func (n *Notifier) SendMessage(ctx context.Context, message Message) error {
	if !n.enabled {
		n.logger.Debug("Notification skipped (disabled)")
		return nil
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var errs []error
	for _, webhookURL := range n.webhookURLs {
		req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(messageJSON))
		if err != nil {
			n.logger.Error("Failed to create request", "webhook", webhookURL, "error", err)
			errs = append(errs, fmt.Errorf("create request (%s): %w", webhookURL, err))
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := n.client.Do(req)
		if err != nil {
			n.logger.Error("Failed to send notification", "webhook", webhookURL, "error", err)
			errs = append(errs, fmt.Errorf("send notification (%s): %w", webhookURL, err))
			continue
		}

		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			n.logger.Error("Failed to read response body", "webhook", webhookURL, "error", err)
			errs = append(errs, fmt.Errorf("read response body (%s): %w", webhookURL, err))
		}

		if err := resp.Body.Close(); err != nil {
			n.logger.Error("Failed to close response body", "webhook", webhookURL, "error", err)
			errs = append(errs, fmt.Errorf("close response body (%s): %w", webhookURL, err))
		}

		if resp.StatusCode != http.StatusOK {
			n.logger.Error("Notification API returned non-200 status code", "webhook", webhookURL, "status", resp.StatusCode)
			errs = append(errs, fmt.Errorf("non-200 status (%s): %d", webhookURL, resp.StatusCode))
			continue
		}
		n.logger.Debug("Notification sent successfully", "webhook", webhookURL)
	}

	if len(errs) > 0 {
		return fmt.Errorf("one or more notifications failed: %w", errors.Join(errs...))
	}
	return nil
}

// NotifyRelayData sends a notification about relay data for a validator
func (n *Notifier) NotifyRelayData(
	ctx context.Context,
	pubkey string,
	validatorIndex,
	blockNumber,
	slot uint64,
	mevReward *big.Int,
	feeRecipient string,
	relaysWithData,
	allRelays []string,
	dashboardInfo *api.DashboardResponse,
) error {
	color := "#36a64f"
	if len(relaysWithData) == 0 {
		color = "#ff9900"
	}

	relaysWithDataStr := "None"
	if len(relaysWithData) > 0 {
		relaysWithDataStr = formatRelayList(relaysWithData)
	}

	fields := []Field{
		{"Validator Index", fmt.Sprintf("%d", validatorIndex), true},
		{"Slot", fmt.Sprintf("%d", slot), true},
		{"Block Number", fmt.Sprintf("%d", blockNumber), true},
		{"Validator Pubkey", pubkey, false},
		{"Relays With Data", fmt.Sprintf("```%s```", relaysWithDataStr), false},
		{"Data Availability", fmt.Sprintf("%d of %d relays have data", len(relaysWithData), len(allRelays)), false},
	}

	title := "Relay Data Available for Validator"
	if len(relaysWithData) == 0 {
		title = "No Relay Data Found for Validator"
	}

	message := Message{
		Attachments: []Attachment{
			{
				Color:      color,
				Title:      title,
				Text:       "Report on relay data for opted-in validator",
				Fields:     fields,
				Footer:     "Validator Monitor",
				TS:         time.Now().Unix(),
				MarkdownIn: []string{"text", "fields"},
			},
		},
	}

	if dashboardInfo != nil {
		attach := &message.Attachments[0]

		attach.Fields = append(attach.Fields, Field{"Block Winner", dashboardInfo.Winner, false})
		attach.Fields = append(attach.Fields, Field{"Commitments", fmt.Sprintf("%d (Rewards: %d, Slashes: %d)",
			dashboardInfo.TotalOpenedCommitments,
			dashboardInfo.TotalRewards,
			dashboardInfo.TotalSlashes), false})

		if dashboardInfo.TotalAmount != "" {
			amountWei, ok := new(big.Int).SetString(dashboardInfo.TotalAmount, 10)
			if ok {
				attach.Fields = append(attach.Fields, Field{"Total Bid Amount", formatWeiToEth(amountWei), true})
			} else {
				attach.Fields = append(attach.Fields, Field{"Total Bid Amount (wei)", dashboardInfo.TotalAmount, true})
			}
		}
		attach.Fields = append(attach.Fields, Field{"MEV Reward", formatWeiToEth(mevReward), true})
		attach.Fields = append(attach.Fields, Field{"MEV Reward Recipient", feeRecipient, true})
	}

	return n.SendMessage(ctx, message)
}

func formatRelayList(relays []string) string {
	if len(relays) == 0 {
		return "None"
	}
	var result string
	for _, r := range relays {
		result += "- " + r + "\n"
	}
	return result
}

func formatWeiToEth(wei *big.Int) string {
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetFloat64(params.Ether))
	return fmt.Sprintf("%.6f ETH", ethValue)
}
