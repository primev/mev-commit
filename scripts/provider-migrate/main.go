package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	providerregistry "github.com/primev/mev-commit/contracts-abi/clients/ProviderRegistry"
)

// ContractABI is the ABI of the ProviderRegistry contract
var ContractABI = providerregistry.ProviderregistryABI

type Provider struct {
	Address      common.Address `json:"address"`
	BlsPublicKey []byte         `json:"bls_public_key"`
}

// Custom UnmarshalJSON method for Provider struct
func (p *Provider) UnmarshalJSON(data []byte) error {
	var temp struct {
		Address      string `json:"address"`
		BlsPublicKey string `json:"bls_public_key"`
	}

	// Unmarshal into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Convert address from string to common.Address
	p.Address = common.HexToAddress(temp.Address)

	// Convert BLS public key from string to []byte
	p.BlsPublicKey = []byte(temp.BlsPublicKey)

	return nil
}

// DelegateRegisterAndStake delegates registration and staking for a list of providers
func DelegateRegisterAndStake(providers []Provider) error {
	// ProviderRegistryAddress is the address of the deployed ProviderRegistry contract
	providerRegistryAddr := os.Getenv("PROVIDER_REGISTRY_ADDRESS")
	if providerRegistryAddr == "" {
		return fmt.Errorf("PROVIDER_REGISTRTY_ADDRESS environment variable is not set")
	}
	providerRegistryAddress := common.HexToAddress(providerRegistryAddr)

	settlementURL := os.Getenv("SETTLEMENT_URL")
	if settlementURL == "" {
		return fmt.Errorf("SETTLEMENT_URL environment variable is not set")
	}

	client, err := ethclient.Dial(settlementURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	defer client.Close()

	prvKey := os.Getenv("PRIVATE_KEY")
	if prvKey == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable is not set")
	}

	privateKey, err := crypto.HexToECDSA(prvKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed getting chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	var ok bool
	auth.Value, ok = new(big.Int).SetString("10000000000000000000", 10) // 10 ETH
	if !ok {
		return fmt.Errorf("failed to set value")
	}
	providerRegistry, err := providerregistry.NewProviderregistryTransactor(
		providerRegistryAddress,
		client,
	)
	if err != nil {
		return fmt.Errorf("failed to instantiate ProviderRegistry contract: %w", err)
	}

	for _, provider := range providers {
		// Delegate registration and staking
		tx, err := providerRegistry.DelegateRegisterAndStake(auth, provider.Address, provider.BlsPublicKey)
		if err != nil {
			log.Printf("failed to delegate register and stake for provider %s: %v", provider.Address.Hex(), err)
			continue
		}
		log.Printf("transaction submitted for provider %s: %s", provider.Address.Hex(), tx.Hash().Hex())
	}

	return nil
}

func main() {
	// Open the JSON file
	file, err := os.Open("providers.json")
	if err != nil {
		log.Fatalf("failed to open providers file: %v", err)
	}
	defer file.Close()

	// Read the file's contents
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read providers file: %v", err)
	}

	// Unmarshal the JSON data into a slice of Provider structs
	var providers []Provider
	err = json.Unmarshal(byteValue, &providers)
	if err != nil {
		log.Fatalf("failed to unmarshal providers: %v", err)
	}

	// Proceed with your logic
	err = DelegateRegisterAndStake(providers)
	if err != nil {
		log.Fatalf("failed to delegate register and stake: %v", err)
	}
}
