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
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/primev/mev-commit/tools/preconf-rpc/sender"
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

type txnInfo struct {
	txn          *sender.Transaction
	noOfAttempts int
	timeTaken    time.Duration
}

// Notifier sends notifications to multiple webhook endpoints
type Notifier struct {
	webhookURLs []string
	client      *http.Client
	logger      *slog.Logger
	queuedTxns  []txnInfo
	queuedMu    sync.Mutex
}

// NewNotifier creates a new notifier instance
func NewNotifier(webhookURLs []string, logger *slog.Logger) *Notifier {
	return &Notifier{
		webhookURLs: webhookURLs,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// SendMessage sends a message to all configured webhook endpoints
func (n *Notifier) SendMessage(ctx context.Context, message Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var retErr error
	for _, url := range n.webhookURLs {
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(messageJSON))
		if err != nil {
			n.logger.Error("Failed to create HTTP request", "url", url, "error", err)
			retErr = errors.Join(retErr, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := n.client.Do(req)
		if err != nil {
			n.logger.Error("Failed to send notification", "url", url, "error", err)
			retErr = errors.Join(retErr, err)
			continue
		}
		func() {
			// Ensure the body is fully read and closed to allow connection reuse
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			errMsg := fmt.Sprintf("non-2xx response: %d - %s", resp.StatusCode, string(body))
			n.logger.Error("Failed to send notification", "url", url, "error", errMsg)
			retErr = errors.Join(retErr, errors.New(errMsg))
			continue
		}

		n.logger.Info("Notification sent successfully", "url", url)
	}

	return retErr
}

type BalanceGetter interface {
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
}

func (n *Notifier) SetupLowBalanceNotification(
	ctx context.Context,
	desc string,
	getter BalanceGetter,
	account common.Address,
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
						Text: "🏦 Low Balance Alert",
						Attachments: []Attachment{
							{
								Color:  "#ff0000",
								Title:  desc,
								Text:   fmt.Sprintf("Account: %s\nBalance: %s (Threshold: %.4f ETH)", account, formatWeiToEth(balance), thresholdEth),
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
	account common.Address,
	amountWei *big.Int,
) error {
	message := Message{
		Text: "💵 Bidder Funded",
		Attachments: []Attachment{
			{
				Color:  "#36a64f",
				Title:  "Bidder account was funded",
				Text:   fmt.Sprintf("Account: %s\nAmount: %s", account, formatWeiToEth(amountWei)),
				Footer: "Preconf RPC Monitor",
				TS:     time.Now().Unix(),
			},
		},
	}
	return n.SendMessage(ctx, message)
}

func (n *Notifier) StartTransactionNotifier(
	ctx context.Context,
) <-chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(15 * time.Minute)

	go func() {
		defer close(done)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				n.logger.Info("Transaction notification routine stopping due to context cancellation")
				return
			case <-ticker.C:
				n.queuedMu.Lock()
				if len(n.queuedTxns) == 0 {
					n.logger.Debug("No queued transactions to notify about")
					n.queuedMu.Unlock()
					continue
				}
				txnsToNotify := n.queuedTxns
				n.queuedTxns = nil
				n.queuedMu.Unlock()
				// create markdown table with the txn info
				text := ""
				for _, txnInfo := range txnsToNotify {
					status := ""
					switch txnInfo.txn.Status {
					case sender.TxStatusPreConfirmed:
						status = "⚡ Pre-Confirmed"
					case sender.TxStatusConfirmed:
						status = "✅ Confirmed"
					case sender.TxStatusFailed:
						status = "❌ Failed"
						status = fmt.Sprintf("%s (Error: %s)", status, txnInfo.txn.Details)
					default:
						status = "❓ Unknown"
					}
					txType := ""
					switch txnInfo.txn.Type {
					case sender.TxTypeRegular:
						txType = "💸 ETH Transaction"
					case sender.TxTypeDeposit:
						txType = "🏦 Deposit"
					case sender.TxTypeInstantBridge:
						txType = "🌉 Instant Bridge"
					default:
						txType = "❓ Unknown"
					}
					text += fmt.Sprintf(
						"- Txn: %s | Sender: %s | Attempts: %d | Duration: %s | Type: %s | Status: %s\n",
						txnInfo.txn.Hash().Hex()[:8],
						txnInfo.txn.Sender.Hex()[:8],
						txnInfo.noOfAttempts,
						txnInfo.timeTaken,
						txType,
						status,
					)
				}
				message := Message{
					Text: "🚀 Transaction Report",
					Attachments: []Attachment{
						{
							Color:  "#FFA500",
							Title:  "The following transactions were completed in the last 15 mins",
							Text:   text,
							Footer: "Preconf RPC Monitor",
							TS:     time.Now().Unix(),
						},
					},
				}
				if err := n.SendMessage(ctx, message); err != nil {
					n.logger.Error("Failed to send 15 minute transaction notification", "error", err)
				}
			}
		}
	}()

	return done
}

func (n *Notifier) NotifyTransactionStatus(
	txn *sender.Transaction,
	noOfAttempts int,
	start time.Time,
) {
	n.queuedMu.Lock()
	defer n.queuedMu.Unlock()

	n.queuedTxns = append(n.queuedTxns, txnInfo{
		txn:          txn,
		noOfAttempts: noOfAttempts,
		timeTaken:    time.Since(start).Round(time.Millisecond),
	})
}

func formatWeiToEth(wei *big.Int) string {
	ethValue := new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetFloat64(params.Ether))
	return fmt.Sprintf("%.6f ETH", ethValue)
}
