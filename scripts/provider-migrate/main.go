package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL       = "https://chainrpc.testnet.mev-commit.xyz/"
	contractAddress = "0xf4F10e18244d836311508917A3B04694D88999Dd"
)

// Updated struct to include current balance
type ProviderRegisteredEvent struct {
	Provider       common.Address `json:"provider"`
	StakedAmount   *big.Int       `json:"stakedAmount"`
	BLSPublicKey   string         `json:"blsPublicKey"`
	CurrentBalance *big.Int       `json:"currentBalance"`
}

func getProviders() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address := common.HexToAddress(contractAddress)

	// Load the contract ABI
	contractABI, err := abi.JSON(strings.NewReader(`[
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "internalType": "address", "name": "provider", "type": "address"},
				{"indexed": false, "internalType": "uint256", "name": "stakedAmount", "type": "uint256"},
				{"indexed": false, "internalType": "bytes", "name": "blsPublicKey", "type": "bytes"}
			],
			"name": "ProviderRegistered",
			"type": "event"
		},
		{
			"inputs": [{"internalType": "address", "name": "", "type": "address"}],
			"name": "providerStakes",
			"outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		}
	]`))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Get the event ID (topic hash)
	eventID := contractABI.Events["ProviderRegistered"].ID

	// Build the query with the event topic
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{eventID}},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to filter logs: %v", err)
	}

	var events []ProviderRegisteredEvent
	for _, vLog := range logs {
		// Initialize the map to hold the unpacked data
		eventData := make(map[string]interface{})

		// Unpack the non-indexed event data
		err := contractABI.UnpackIntoMap(eventData, "ProviderRegistered", vLog.Data)
		if err != nil {
			log.Fatalf("Failed to unpack log data: %v", err)
		}

		// Extract the indexed field (Provider)
		event := ProviderRegisteredEvent{}
		if len(vLog.Topics) > 1 {
			event.Provider = common.HexToAddress(vLog.Topics[1].Hex())
		}

		// Assign the non-indexed fields from eventData
		if stakedAmount, ok := eventData["stakedAmount"].(*big.Int); ok {
			event.StakedAmount = stakedAmount
		} else {
			log.Fatalf("Failed to get stakedAmount from event data")
		}

		if blsPublicKey, ok := eventData["blsPublicKey"].([]byte); ok {
			event.BLSPublicKey = hex.EncodeToString(blsPublicKey)
		} else {
			log.Fatalf("Failed to get blsPublicKey from event data")
		}

		// Now, get the current balance from the providerStakes mapping
		currentBalance, err := getProviderStake(client, contractABI, address, event.Provider)
		if err != nil {
			log.Fatalf("Failed to get current balance for provider %s: %v", event.Provider.Hex(), err)
		}
		event.CurrentBalance = currentBalance

		events = append(events, event)
	}

	file, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal events to JSON: %v", err)
	}

	err = os.WriteFile("provider_registered_events.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write events to file: %v", err)
	}

	fmt.Println("Events successfully written to provider_registered_events.json")
}

// Function to get the current stake of a provider
func getProviderStake(client *ethclient.Client, contractABI abi.ABI, contractAddress common.Address, providerAddress common.Address) (*big.Int, error) {
	// Prepare the data for the call
	data, err := contractABI.Pack("providerStakes", providerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to pack parameters: %v", err)
	}

	// Call the contract function
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	// Use nil for block number to get the latest state
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	// Unpack the result
	outputs, err := contractABI.Unpack("providerStakes", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}

	// The output should be a slice with one element: the staked amount
	if len(outputs) != 1 {
		return nil, fmt.Errorf("unexpected output length: %d", len(outputs))
	}

	currentBalance, ok := outputs[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected type in output")
	}

	return currentBalance, nil
}

func main() {
	getProviders()
}
