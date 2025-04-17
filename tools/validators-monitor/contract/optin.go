package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
)

// OptInStatus matches the struct in the contract
type OptInStatus struct {
	IsVanillaOptedIn    bool
	IsAvsOptedIn        bool
	IsMiddlewareOptedIn bool
}

// ValidatorOptInChecker is responsible for checking validator opt-in status
type ValidatorOptInChecker struct {
	client      *ethclient.Client
	contractAbi abi.ABI
	address     common.Address
}

// NewValidatorOptInChecker creates a new ValidatorOptInChecker
func NewValidatorOptInChecker(rpcURL, contractAddress string) (*ValidatorOptInChecker, error) {
	// Connect to Ethereum RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %s, %v", rpcURL, err)
	}

	// Parse contract ABI
	contractAbi, err := abi.JSON(strings.NewReader(validatorrouter.ValidatoroptinrouterABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %v", err)
	}

	// Convert address string to eth common.Address
	address := common.HexToAddress(contractAddress)

	return &ValidatorOptInChecker{
		client:      client,
		contractAbi: contractAbi,
		address:     address,
	}, nil
}

// CheckValidatorsOptedIn checks which validators in a batch are opted in
func (c *ValidatorOptInChecker) CheckValidatorsOptedIn(ctx context.Context, pubkeys []string) ([]OptInStatus, error) {
	blsPubKeys := make([][]byte, len(pubkeys))
	for j, pubkey := range pubkeys {
		pubkeyStr := strings.TrimPrefix(pubkey, "0x")
		pubkeyBytes, err := hex.DecodeString(pubkeyStr)
		if err != nil {
			return nil, fmt.Errorf("error decoding pubkey %s: %w", pubkey, err)
		}
		blsPubKeys[j] = pubkeyBytes
	}

	callData, err := c.contractAbi.Pack("areValidatorsOptedIn", blsPubKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to pack contract call data: %w", err)
	}

	result, err := c.client.CallContract(ctx, ethereum.CallMsg{
		To:   &c.address,
		Data: callData,
	}, nil) // nil means latest block

	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	var optInStatuses []OptInStatus
	err = c.contractAbi.UnpackIntoInterface(&optInStatuses, "areValidatorsOptedIn", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack contract result: %w", err)
	}

	return optInStatuses, nil
}

// Close releases the ethclient connection
func (c *ValidatorOptInChecker) Close() {
	c.client.Close()
}
