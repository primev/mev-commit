package notifier

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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
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

type BalanceGetter interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

func (n *Notifier) SetupLowBalanceNotification(
	ctx context.Context,
	chainDesc string,
	getter BalanceGetter,
	account common.Address,
	msg string,
	thresholdEth float64,
	checkInterval time.Duration,
	alertCooldown time.Duration,
) <-chan struct{} {
	ticker := time.NewTicker(checkInterval)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer ticker.Stop()
		var lastAlert time.Time

		for {
			select {
			case <-ctx.Done():
				n.logger.Info("Low balance notification routine stopping due to context cancellation")
				return
			case <-ticker.C:
				if time.Since(lastAlert) < alertCooldown {
					continue
				}

				balance, err := getter.BalanceAt(ctx, account, nil)
				if err != nil {
					n.logger.Error("Failed to get account balance", "account", account, "error", err)
					continue
				}
				balanceEth := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(params.Ether))
				balanceEthFloat, _ := balanceEth.Float64()
				if balanceEthFloat < thresholdEth {
					message := Message{
						Text: "⚠️ Low Balance Alert",
						Attachments: []Attachment{
							{
								Color:  "#ff0000",
								Title:  fmt.Sprintf("The following accounts have low balances on %s chain:", chainDesc),
								Text:   fmt.Sprintf("%s\n\nAccount: %s\nBalance: %s (Threshold: %.4f ETH)", msg, account, formatWeiToEth(balance), thresholdEth),
								Footer: "Preconf RPC Monitor",
								TS:     time.Now().Unix(),
							},
						},
					}
					if err := n.SendMessage(ctx, message); err != nil {
						n.logger.Error("Failed to send low balance notification", "error", err)
					} else {
						lastAlert = time.Now()
					}
				}
			}
		}
	}()
	return done
}

func (n *Notifier) SendBidderFundedNotification(
	ctx context.Context,
	chainDesc string,
	account common.Address,
	amountWei *big.Int,
) error {
	message := Message{
		Text: "🎉 Bidder Funded Alert",
		Attachments: []Attachment{
			{
				Color:  "#36a64f",
				Title:  fmt.Sprintf("A bidder account has been funded on %s chain:", chainDesc),
				Text:   fmt.Sprintf("Account: %s\nAmount: %s", account, formatWeiToEth(amountWei)),
				Footer: "Preconf RPC Monitor",
				TS:     time.Now().Unix(),
			},
		},
	}
	return n.SendMessage(ctx, message)
}

func formatWeiToEth(wei *big.Int) string {
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetFloat64(params.Ether))
	return fmt.Sprintf("%.6f ETH", ethValue)
}
