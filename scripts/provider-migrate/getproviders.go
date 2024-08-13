package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	infuraURL               = "http://69.67.149.79:8545"
	contractAddress         = "0x4FC9b98e1A0Ff10de4c2cf294656854F1d5B207D"
	eventProviderRegistered = "ProviderRegistered(address,uint256,bytes)"
)

type ProviderRegisteredEvent struct {
	Provider     common.Address `json:"provider"`
	StakedAmount *big.Int       `json:"stakedAmount"`
	BLSPublicKey []byte         `json:"blsPublicKey"`
}

func getProviders() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address := common.HexToAddress(contractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	// Load the contract ABI
	contractABI, err := abi.JSON(strings.NewReader(`[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"provider","type":"address"},{"indexed":false,"internalType":"uint256","name":"stakedAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"blsPublicKey","type":"bytes"}],"name":"ProviderRegistered","type":"event"}]`))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to filter logs: %v", err)
	}

	var events []ProviderRegisteredEvent
	for _, vLog := range logs {
		event := new(ProviderRegisteredEvent)

		err := contractABI.UnpackIntoInterface(event, "ProviderRegistered", vLog.Data)
		if err != nil {
			log.Fatalf("Failed to unpack log data: %v", err)
		}

		// Add the indexed event data
		event.Provider = common.HexToAddress(vLog.Topics[1].Hex())
		events = append(events, *event)
	}

	file, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal events to JSON: %v", err)
	}

	err = ioutil.WriteFile("provider_registered_events.json", file, 0644)
	if err != nil {
		log.Fatalf("Failed to write events to file: %v", err)
	}

	fmt.Println("Events successfully written to provider_registered_events.json")
}

