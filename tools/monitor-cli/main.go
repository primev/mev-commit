package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	bidderregistry "github.com/primev/mev-commit/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primev/mev-commit/contracts-abi/clients/BlockTracker"
	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
)

//go:embed templates/*
var content embed.FS

type LogEntry struct {
	Timestamp   string `json:"timestamp"`
	EventName   string `json:"eventName"`
	Details     string `json:"details"`
	BlockNumber uint64 `json:"blockNumber"`
}

type LogStore struct {
	sync.RWMutex
	logs []LogEntry
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	logStore = &LogStore{
		logs: make([]LogEntry, 0),
	}
	wsClients = make(map[*websocket.Conn]bool)
	wsLock    sync.RWMutex
)

func main() {
	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFile, err := os.OpenFile(fmt.Sprintf("provider_registry_events_%s.log", timestamp), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Set up logging to both file and console
	multiLogger := log.New(os.Stdout, "", 0)
	fileLogger := log.New(logFile, "", 0)

	// Connect to the WebSocket endpoint
	client, err := ethclient.Dial("wss://chainrpc-wss.testnet.mev-commit.xyz")
	if err != nil {
		multiLogger.Fatalf("Failed to connect to the WebSocket client: %v", err)
	}

	// Create contract instances
	providerRegistryAddr := common.HexToAddress("0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67")
	providerContract, err := providerregistry.NewProviderregistryFilterer(providerRegistryAddr, client)
	if err != nil {
		multiLogger.Fatalf("Failed to instantiate provider registry contract: %v", err)
	}

	bidderRegistryAddr := common.HexToAddress("0x948eCD70FaeF6746A30a00F30f8b9fB2659e4062")
	bidderContract, err := bidderregistry.NewBidderregistryFilterer(bidderRegistryAddr, client)
	if err != nil {
		multiLogger.Fatalf("Failed to instantiate bidder registry contract: %v", err)
	}

	blockTrackerAddr := common.HexToAddress("0x0b3b6Cf113959214E313d6Ad37Ad56831acb1776")
	blockTrackerContract, err := blocktracker.NewBlocktrackerFilterer(blockTrackerAddr, client)
	if err != nil {
		multiLogger.Fatalf("Failed to instantiate block tracker contract: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channels for all events
	blsKeyAddedChan := make(chan *providerregistry.ProviderregistryBLSKeyAdded)
	feePayoutPeriodBlocksUpdatedChan := make(chan *providerregistry.ProviderregistryFeePayoutPeriodBlocksUpdated)
	feePercentUpdatedChan := make(chan *providerregistry.ProviderregistryFeePercentUpdated)
	feeTransferChan := make(chan *providerregistry.ProviderregistryFeeTransfer)
	insufficientFundsToSlashChan := make(chan *providerregistry.ProviderregistryInsufficientFundsToSlash)
	minStakeUpdatedChan := make(chan *providerregistry.ProviderregistryMinStakeUpdated)
	ownershipTransferStartedChan := make(chan *providerregistry.ProviderregistryOwnershipTransferStarted)
	ownershipTransferredChan := make(chan *providerregistry.ProviderregistryOwnershipTransferred)
	pausedChan := make(chan *providerregistry.ProviderregistryPaused)
	preconfManagerUpdatedChan := make(chan *providerregistry.ProviderregistryPreconfManagerUpdated)
	unstakeChan := make(chan *providerregistry.ProviderregistryUnstake)
	unpausedChan := make(chan *providerregistry.ProviderregistryUnpaused)
	withdrawalDelayUpdatedChan := make(chan *providerregistry.ProviderregistryWithdrawalDelayUpdated)
	providerSlashingChan := make(chan *providerregistry.ProviderregistryFundsSlashed)
	bidderRewardChan := make(chan *bidderregistry.BidderregistryFundsRewarded)
	newL1BlockChan := make(chan *blocktracker.BlocktrackerNewL1Block)

	// Set up HTTP server and routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/logs", handleGetLogs)

	// Start HTTP server
	go func() {
		multiLogger.Println("Starting web server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			multiLogger.Fatalf("Failed to start web server: %v", err)
		}
	}()

	// Watch for events
	go watchEvents(ctx, providerContract, bidderContract, blockTrackerContract, EventChannels{
		blsKeyAddedChan:                  blsKeyAddedChan,
		feePayoutPeriodBlocksUpdatedChan: feePayoutPeriodBlocksUpdatedChan,
		feePercentUpdatedChan:            feePercentUpdatedChan,
		feeTransferChan:                  feeTransferChan,
		insufficientFundsToSlashChan:     insufficientFundsToSlashChan,
		minStakeUpdatedChan:              minStakeUpdatedChan,
		ownershipTransferStartedChan:     ownershipTransferStartedChan,
		ownershipTransferredChan:         ownershipTransferredChan,
		pausedChan:                       pausedChan,
		preconfManagerUpdatedChan:        preconfManagerUpdatedChan,
		unstakeChan:                      unstakeChan,
		unpausedChan:                     unpausedChan,
		withdrawalDelayUpdatedChan:       withdrawalDelayUpdatedChan,
		providerSlashingChan:             providerSlashingChan,
		bidderRewardChan:                 bidderRewardChan,
		newL1BlockChan:                   newL1BlockChan,
	}, multiLogger, fileLogger)

	multiLogger.Println("Started watching for contract events...")
	fileLogger.Println("Started watching for contract events...")

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	for {
		select {
		case event := <-blsKeyAddedChan:
			logEvent(ctx, client, "BLS Key Added", event.Raw, fmt.Sprintf("Provider: %s, BLS Key: %x", event.Provider.Hex(), event.BlsPublicKey), multiLogger, fileLogger)
		case event := <-feePayoutPeriodBlocksUpdatedChan:
			logEvent(ctx, client, "Fee Payout Period Blocks Updated", event.Raw, fmt.Sprintf("New Period: %s", event.NewFeePayoutPeriodBlocks), multiLogger, fileLogger)
		case event := <-feePercentUpdatedChan:
			logEvent(ctx, client, "Fee Percent Updated", event.Raw, fmt.Sprintf("New Fee Percent: %s", event.NewFeePercent), multiLogger, fileLogger)
		case event := <-feeTransferChan:
			logEvent(ctx, client, "Fee Transfer", event.Raw, fmt.Sprintf("Recipient: %s, Amount: %s", event.Recipient.Hex(), event.Amount), multiLogger, fileLogger)
		case event := <-insufficientFundsToSlashChan:
			logEvent(ctx, client, "Insufficient Funds To Slash", event.Raw, fmt.Sprintf("Provider: %s, Stake: %s, Residual: %s, Penalty: %s", event.Provider.Hex(), event.ProviderStake, event.ResidualAmount, event.PenaltyFee), multiLogger, fileLogger)
		case event := <-minStakeUpdatedChan:
			logEvent(ctx, client, "Min Stake Updated", event.Raw, fmt.Sprintf("New Min Stake: %s", event.NewMinStake), multiLogger, fileLogger)
		case event := <-ownershipTransferStartedChan:
			logEvent(ctx, client, "Ownership Transfer Started", event.Raw, fmt.Sprintf("Previous: %s, New: %s", event.PreviousOwner.Hex(), event.NewOwner.Hex()), multiLogger, fileLogger)
		case event := <-ownershipTransferredChan:
			logEvent(ctx, client, "Ownership Transferred", event.Raw, fmt.Sprintf("Previous: %s, New: %s", event.PreviousOwner.Hex(), event.NewOwner.Hex()), multiLogger, fileLogger)
		case event := <-pausedChan:
			logEvent(ctx, client, "Contract Paused", event.Raw, fmt.Sprintf("Account: %s", event.Account.Hex()), multiLogger, fileLogger)
		case event := <-preconfManagerUpdatedChan:
			logEvent(ctx, client, "Preconf Manager Updated", event.Raw, fmt.Sprintf("New Manager: %s", event.NewPreconfManager.Hex()), multiLogger, fileLogger)
		case event := <-unstakeChan:
			logEvent(ctx, client, "Unstake", event.Raw, fmt.Sprintf("Provider: %s", event.Provider.Hex()), multiLogger, fileLogger)
		case event := <-unpausedChan:
			logEvent(ctx, client, "Contract Unpaused", event.Raw, fmt.Sprintf("Account: %s", event.Account.Hex()), multiLogger, fileLogger)
		case event := <-withdrawalDelayUpdatedChan:
			logEvent(ctx, client, "Withdrawal Delay Updated", event.Raw, fmt.Sprintf("New Delay: %s", event.NewWithdrawalDelay), multiLogger, fileLogger)
		case event := <-providerSlashingChan:
			logEvent(ctx, client, "Provider Slashing", event.Raw, fmt.Sprintf("Provider: %s, Amount: %s", event.Provider.Hex(), event.Amount), multiLogger, fileLogger)
		case event := <-bidderRewardChan:
			logEvent(ctx, client, "Bidder Reward", event.Raw, fmt.Sprintf("Bidder: %s, Amount: %s", event.Bidder.Hex(), event.Amount), multiLogger, fileLogger)
		case event := <-newL1BlockChan:
			if event.Winner != [20]byte{} {
				logEvent(ctx, client, "New L1 Block", event.Raw, fmt.Sprintf("Block Number: %s, Winner: %s", event.BlockNumber, event.Winner), multiLogger, fileLogger)
			}
		case sig := <-sigChan:
			shutdownMsg := fmt.Sprintf("\nCaught signal %v. Shutting down...", sig)
			multiLogger.Println(shutdownMsg)
			fileLogger.Println(shutdownMsg)
			cancel()
			return
		}
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(content, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	wsLock.Lock()
	wsClients[conn] = true
	wsLock.Unlock()

	defer func() {
		wsLock.Lock()
		delete(wsClients, conn)
		wsLock.Unlock()
	}()

	// Keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	logStore.RLock()
	defer logStore.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logStore.logs)
}

func broadcastToWebSockets(entry LogEntry) {
	wsLock.RLock()
	defer wsLock.RUnlock()

	for client := range wsClients {
		if err := client.WriteJSON(entry); err != nil {
			log.Printf("Failed to write to websocket: %v", err)
			client.Close()
			delete(wsClients, client)
		}
	}
}

type EventChannels struct {
	blsKeyAddedChan                  chan *providerregistry.ProviderregistryBLSKeyAdded
	feePayoutPeriodBlocksUpdatedChan chan *providerregistry.ProviderregistryFeePayoutPeriodBlocksUpdated
	feePercentUpdatedChan            chan *providerregistry.ProviderregistryFeePercentUpdated
	feeTransferChan                  chan *providerregistry.ProviderregistryFeeTransfer
	insufficientFundsToSlashChan     chan *providerregistry.ProviderregistryInsufficientFundsToSlash
	minStakeUpdatedChan              chan *providerregistry.ProviderregistryMinStakeUpdated
	ownershipTransferStartedChan     chan *providerregistry.ProviderregistryOwnershipTransferStarted
	ownershipTransferredChan         chan *providerregistry.ProviderregistryOwnershipTransferred
	pausedChan                       chan *providerregistry.ProviderregistryPaused
	preconfManagerUpdatedChan        chan *providerregistry.ProviderregistryPreconfManagerUpdated
	unstakeChan                      chan *providerregistry.ProviderregistryUnstake
	unpausedChan                     chan *providerregistry.ProviderregistryUnpaused
	withdrawalDelayUpdatedChan       chan *providerregistry.ProviderregistryWithdrawalDelayUpdated
	providerSlashingChan             chan *providerregistry.ProviderregistryFundsSlashed
	bidderRewardChan                 chan *bidderregistry.BidderregistryFundsRewarded
	newL1BlockChan                   chan *blocktracker.BlocktrackerNewL1Block
}

func watchEvents(
	ctx context.Context,
	providerContract *providerregistry.ProviderregistryFilterer,
	bidderContract *bidderregistry.BidderregistryFilterer,
	blockTrackerContract *blocktracker.BlocktrackerFilterer,
	channels EventChannels,
	multiLogger, fileLogger *log.Logger,
) {
	startBlock := big.NewInt(0)
	blockIncrement := big.NewInt(1000000)

	for {
		endBlock := new(big.Int).Add(startBlock, blockIncrement)
		endBlockUint := endBlock.Uint64()

		filterOpts := &bind.FilterOpts{
			Start:   startBlock.Uint64(),
			End:     &endBlockUint,
			Context: ctx,
		}

		// Filter all events
		if logs, err := providerContract.FilterBLSKeyAdded(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.blsKeyAddedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterFeePayoutPeriodBlocksUpdated(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.feePayoutPeriodBlocksUpdatedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterFeePercentUpdated(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.feePercentUpdatedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterFeeTransfer(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.feeTransferChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterInsufficientFundsToSlash(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.insufficientFundsToSlashChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterMinStakeUpdated(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.minStakeUpdatedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterOwnershipTransferStarted(filterOpts, nil, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.ownershipTransferStartedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterOwnershipTransferred(filterOpts, nil, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.ownershipTransferredChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterPaused(filterOpts); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.pausedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterPreconfManagerUpdated(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.preconfManagerUpdatedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterUnstake(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.unstakeChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterUnpaused(filterOpts); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.unpausedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterWithdrawalDelayUpdated(filterOpts); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.withdrawalDelayUpdatedChan <- logs.Event:
				}
			}
		}

		if logs, err := providerContract.FilterFundsSlashed(filterOpts, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.providerSlashingChan <- logs.Event:
				}
			}
		}

		if logs, err := bidderContract.FilterFundsRewarded(filterOpts, nil, nil, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.bidderRewardChan <- logs.Event:
				}
			}
		}

		if logs, err := blockTrackerContract.FilterNewL1Block(filterOpts, nil, nil, nil); err == nil {
			for logs.Next() {
				select {
				case <-ctx.Done():
					return
				case channels.newL1BlockChan <- logs.Event:
				}
			}
		}

		startBlock = endBlock

		select {
		case <-ctx.Done():
			return
		default:
			// Continue to next block range
		}
	}
}

func logEvent(ctx context.Context, client *ethclient.Client, eventName string, rawLog types.Log, details string, multiLogger, fileLogger *log.Logger) {
	block, err := client.BlockByHash(ctx, rawLog.BlockHash)
	if err != nil {
		multiLogger.Printf("Failed to get block details: %v", err)
		return
	}
	blockTime := time.UnixMilli(int64(block.Time()))

	entry := LogEntry{
		Timestamp:   blockTime.Format(time.RFC3339),
		EventName:   eventName,
		Details:     details,
		BlockNumber: rawLog.BlockNumber,
	}

	// Store the log
	logStore.Lock()
	logStore.logs = append(logStore.logs, entry)
	logStore.Unlock()

	// Broadcast to WebSocket clients
	broadcastToWebSockets(entry)

	eventMsg := fmt.Sprintf("[%s] %s - %s, Block: %d",
		entry.Timestamp,
		eventName,
		details,
		rawLog.BlockNumber)
	multiLogger.Println(eventMsg)
}
