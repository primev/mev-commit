package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/primev/mev-commit/tools/validators-monitor/api"
)

// SlackMessage represents a Slack message structure
type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a Slack message attachment
type Attachment struct {
	Color      string   `json:"color,omitempty"`
	Title      string   `json:"title,omitempty"`
	Text       string   `json:"text,omitempty"`
	Fields     []Field  `json:"fields,omitempty"`
	Footer     string   `json:"footer,omitempty"`
	TS         int64    `json:"ts,omitempty"`
	MarkdownIn []string `json:"mrkdwn_in,omitempty"`
}

// Field represents a field in a Slack message attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// SlackNotifier sends notifications to Slack
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
	logger     *slog.Logger
	enabled    bool
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(webhookURL string, logger *slog.Logger) *SlackNotifier {
	enabled := webhookURL != ""

	if !enabled {
		logger.Warn("Slack notifications disabled - no webhook URL provided")
	} else {
		logger.Info("Slack notifications enabled")
	}

	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:  logger,
		enabled: enabled,
	}
}

// SendMessage sends a message to Slack
func (n *SlackNotifier) SendMessage(ctx context.Context, message SlackMessage) error {
	if !n.enabled {
		n.logger.Debug("Slack notification skipped (disabled)")
		return nil
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned non-200 status code: %d", resp.StatusCode)
	}

	n.logger.Debug("Slack notification sent successfully")
	return nil
}

// NotifyRelayData sends a notification about relay data for a validator
func (n *SlackNotifier) NotifyRelayData(
	ctx context.Context, pubkey string,
	validatorIndex, blockNumber, slot uint64, mevReward *big.Int,
	relaysWithData []string, allRelays []string,
	dashboardInfo *api.DashboardResponse,
) error {
	if !n.enabled {
		return nil
	}

	color := "#36a64f" // Green for relays having data
	if len(relaysWithData) == 0 {
		color = "#ff9900" // Orange for no data available
	}

	relaysWithDataStr := "None"
	if len(relaysWithData) > 0 {
		relaysWithDataStr = formatRelayList(relaysWithData)
	}

	fields := []Field{
		{
			Title: "Validator Index",
			Value: fmt.Sprintf("%d", validatorIndex),
			Short: true,
		},
		{
			Title: "Slot",
			Value: fmt.Sprintf("%d", slot),
			Short: true,
		},
		{
			Title: "Block Number",
			Value: fmt.Sprintf("%d", blockNumber),
			Short: true,
		},
		{
			Title: "Validator Pubkey",
			Value: pubkey,
			Short: false,
		},
		{
			Title: "Relays With Data",
			Value: fmt.Sprintf("```%s```", relaysWithDataStr),
			Short: false,
		},
		{
			Title: "Data Availability",
			Value: fmt.Sprintf("%d of %d relays have data", len(relaysWithData), len(allRelays)),
			Short: false,
		},
	}

	title := "Relay Data Available for Validator"
	if len(relaysWithData) == 0 {
		title = "No Relay Data Found for Validator"
	}

	message := SlackMessage{
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
		attachment := &message.Attachments[0]

		attachment.Fields = append(attachment.Fields, Field{
			Title: "Block Winner",
			Value: dashboardInfo.Winner,
			Short: false,
		})

		attachment.Fields = append(attachment.Fields, Field{
			Title: "Commitments",
			Value: fmt.Sprintf("%d (Rewards: %d, Slashes: %d)",
				dashboardInfo.TotalOpenedCommitments,
				dashboardInfo.TotalRewards,
				dashboardInfo.TotalSlashes),
			Short: false,
		})

		if dashboardInfo.TotalAmount != "" {
			amountWei, ok := new(big.Int).SetString(dashboardInfo.TotalAmount, 10)
			if ok {
				attachment.Fields = append(attachment.Fields, Field{
					Title: "Total Bid Amount",
					Value: formatWeiToEth(amountWei),
					Short: true,
				})
			} else {
				// If conversion fails, show raw value
				attachment.Fields = append(attachment.Fields, Field{
					Title: "Total Bid Amount (wei)",
					Value: dashboardInfo.TotalAmount,
					Short: true,
				})
			}
		}
		attachment.Fields = append(attachment.Fields, Field{
			Title: "MEV Reward",
			Value: formatWeiToEth(mevReward),
			Short: true,
		})
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
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(wei),
		new(big.Float).SetFloat64(params.Ether),
	)

	return fmt.Sprintf("%.6f ETH", ethValue)
}
