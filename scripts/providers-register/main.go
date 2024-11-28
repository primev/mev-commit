package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	rpcURL             = "https://chainrpc.testnet.mev-commit.xyz/"
	newContractAddress = "0x1C2a592950E5dAd49c0E2F3A402DCF496bdf7b67"
	providersJSONFile  = "provider_registered_events.json"
)

type ProviderRegisteredEvent struct {
	Provider       common.Address `json:"provider"`
	StakedAmount   *big.Int       `json:"stakedAmount"`
	BLSPublicKey   string         `json:"blsPublicKey"`
	CurrentBalance *big.Int       `json:"currentBalance"`
}

func main() {
	privateKeyHex := os.Getenv("OWNER_PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("OWNER_PRIVATE_KEY environment variable not set")
	}
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to get public key")
	}
	ownerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Owner Address: %s\n", ownerAddress.Hex())

	events, err := readProvidersFromJSON(providersJSONFile)
	if err != nil {
		log.Fatalf("Failed to read providers from JSON: %v", err)
	}

	err = registerProviders(client, privateKey, ownerAddress, events)
	if err != nil {
		log.Fatalf("Failed to register providers: %v", err)
	}

	fmt.Println("All providers have been registered in the new contract.")
}

func readProvidersFromJSON(filename string) ([]ProviderRegisteredEvent, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var events []ProviderRegisteredEvent
	err = json.Unmarshal(fileData, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return events, nil
}

func registerProviders(client *ethclient.Client, privateKey *ecdsa.PrivateKey, ownerAddress common.Address, events []ProviderRegisteredEvent) error {
	newContractAddr := common.HexToAddress(newContractAddress)

	contractABI, err := abi.JSON(strings.NewReader(`[
		{
			"inputs": [
				{"internalType": "address", "name": "provider", "type": "address"},
				{"internalType": "bytes[]", "name": "blsPublicKeys", "type": "bytes[]"}
			],
			"name": "delegateRegisterAndStake",
			"outputs": [],
			"stateMutability": "payable",
			"type": "function"
		}
	]`))
	if err != nil {
		return fmt.Errorf("failed to parse new contract ABI: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), ownerAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get network ID: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}
	
	for _, event := range events {
		fmt.Printf("Registering provider: %s\n", event.Provider.Hex())

		blsPublicKeyBytes, err := hex.DecodeString(event.BLSPublicKey)
		if err != nil {
			return fmt.Errorf("failed to decode BLS public key: %v", err)
		}

		providerAddress := event.Provider
		blsPublicKeys := [][]byte{blsPublicKeyBytes}

		input, err := contractABI.Pack("delegateRegisterAndStake", providerAddress, blsPublicKeys)
		if err != nil {
			return fmt.Errorf("failed to pack function arguments: %v", err)
		}

		msg := ethereum.CallMsg{
			From:     ownerAddress,
			To:       &newContractAddr,
			GasPrice: gasPrice,
			Value:    event.CurrentBalance,
			Data:     input,
		}
		gasLimit, err := client.EstimateGas(context.Background(), msg)
		if err != nil {
			gasLimit = uint64(300000)
			fmt.Printf("Failed to estimate gas limit, using default: %v\n", err)
		}

		tx := types.NewTransaction(nonce, newContractAddr, event.CurrentBalance, gasLimit, gasPrice, input)

		signedTx, err := auth.Signer(ownerAddress, tx)
		if err != nil {
			return fmt.Errorf("failed to sign transaction: %v", err)
		}

		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			return fmt.Errorf("failed to send transaction: %v", err)
		}

		fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())

		nonce++

		receipt, err := waitForReceipt(client, signedTx.Hash())
		if err != nil {
			return fmt.Errorf("failed to get transaction receipt: %v", err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("transaction failed: %s", signedTx.Hash().Hex())
		}

		fmt.Printf("Provider %s successfully registered.\n", event.Provider.Hex())

		time.Sleep(1 * time.Second)
	}

	return nil
}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			if err == ethereum.NotFound {
				time.Sleep(1 * time.Second)
				continue
			}
			return nil, err
		}
		return receipt, nil
	}
}
