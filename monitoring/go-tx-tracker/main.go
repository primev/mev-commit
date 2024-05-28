package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

func main() {
	// Replace with your Geth RPC URL
	rpcURL := "http://localhost:8545"

	client, err := rpc.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("Connected to Geth RPC")

	for {
		trackTransactions(client)
		time.Sleep(10 * time.Second)
	}
}

func trackTransactions(client *rpc.Client) {
	var result map[string]map[string]map[string]interface{}
	err := client.CallContext(context.Background(), &result, "txpool_content")
	if err != nil {
		log.Printf("Failed to get txpool content: %v", err)
		return
	}

	fmt.Println("Pending Transactions:")
	printTransactions(result["pending"])

	fmt.Println("\nQueued Transactions:")
	printTransactions(result["queued"])
}

func printTransactions(txs map[string]map[string]interface{}) {
	for account, accountTxs := range txs {
		for nonce, tx := range accountTxs {
			var transaction Transaction
			txBytes, err := json.Marshal(tx)
			if err != nil {
				log.Printf("Failed to marshal transaction: %v", err)
				continue
			}
			err = json.Unmarshal(txBytes, &transaction)
			if err != nil {
				log.Printf("Failed to unmarshal transaction: %v", err)
				continue
			}
			txJson, err := json.MarshalIndent(transaction, "", "  ")
			if err != nil {
				log.Printf("Failed to marshal transaction: %v", err)
				continue
			}
			fmt.Printf("Account: %s, Nonce: %s, Tx: %s\n", account, nonce, txJson)
		}
	}
}
