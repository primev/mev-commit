package contract

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock implementation of EthClient
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	args := m.Called(ctx, call, blockNumber)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockEthClient) Close() {
	m.Called()
}

// Test for NewValidatorOptInChecker
func TestNewValidatorOptInChecker(t *testing.T) {
	// Test successful creation
	checker, err := NewValidatorOptInChecker("http://localhost:8545", "0x1234567890123456789012345678901234567890")
	if err != nil {
		t.Fatalf("Failed to create validator checker: %v", err)
	}
	defer checker.Close()

	assert.NotNil(t, checker.client)
	assert.NotNil(t, checker.contractAbi)
	assert.Equal(t, common.HexToAddress("0x1234567890123456789012345678901234567890"), checker.address)
}

// Test for CheckValidatorsOptedIn
func TestCheckValidatorsOptedIn(t *testing.T) {
	// Create a test validator checker with a mock client
	mockClient := new(MockEthClient)
	contractAbi, err := abi.JSON(strings.NewReader(validatorrouter.ValidatoroptinrouterABI))
	require.NoError(t, err)

	checker := &ValidatorOptInChecker{
		client:      mockClient,
		contractAbi: contractAbi,
		address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}

	// Test data
	pubkeys := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	decodedPubKey1, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	decodedPubKey2, _ := hex.DecodeString("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	expectedBlsPubKeys := [][]byte{decodedPubKey1, decodedPubKey2}

	callData, err := contractAbi.Pack("areValidatorsOptedIn", expectedBlsPubKeys)
	require.NoError(t, err)

	expectedStatuses := []OptInStatus{
		{IsVanillaOptedIn: true, IsAvsOptedIn: false, IsMiddlewareOptedIn: true},
		{IsVanillaOptedIn: false, IsAvsOptedIn: true, IsMiddlewareOptedIn: false},
	}

	packedResult, err := contractAbi.Methods["areValidatorsOptedIn"].Outputs.Pack(expectedStatuses)
	require.NoError(t, err)

	mockClient.On("CallContract", mock.Anything, mock.MatchedBy(func(call ethereum.CallMsg) bool {
		return reflect.DeepEqual(call.Data, callData) && *call.To == checker.address
	}), mock.Anything).Return(packedResult, nil)

	ctx := context.Background()
	statuses, err := checker.CheckValidatorsOptedIn(ctx, pubkeys)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatuses, statuses)
	mockClient.AssertExpectations(t)
}

// Test error handling for CheckValidatorsOptedIn
func TestCheckValidatorsOptedIn_Errors(t *testing.T) {
	mockClient := new(MockEthClient)
	contractAbi, err := abi.JSON(strings.NewReader(validatorrouter.ValidatoroptinrouterABI))
	require.NoError(t, err)

	checker := &ValidatorOptInChecker{
		client:      mockClient,
		contractAbi: contractAbi,
		address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}

	t.Run("Invalid pubkey format", func(t *testing.T) {
		pubkeys := []string{"invalid_pubkey"}
		_, err := checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error decoding pubkey")
	})

	t.Run("Contract call failure", func(t *testing.T) {
		pubkeys := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}
		decodedPubKey, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		expectedBlsPubKeys := [][]byte{decodedPubKey}

		callData, err := contractAbi.Pack("areValidatorsOptedIn", expectedBlsPubKeys)
		require.NoError(t, err)

		mockClient.On("CallContract", mock.Anything, mock.MatchedBy(func(call ethereum.CallMsg) bool {
			return reflect.DeepEqual(call.Data, callData)
		}), mock.Anything).Return([]byte{}, errors.New("contract call failed"))

		_, err = checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to call contract")
		mockClient.AssertExpectations(t)
	})

	t.Run("Result unpacking failure", func(t *testing.T) {
		pubkeys := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}
		decodedPubKey, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		expectedBlsPubKeys := [][]byte{decodedPubKey}

		callData, err := contractAbi.Pack("areValidatorsOptedIn", expectedBlsPubKeys)
		require.NoError(t, err)

		invalidData := []byte{0x01, 0x02, 0x03} // Too short to be valid ABI-encoded data

		mockClient.ExpectedCalls = nil

		mockClient.On("CallContract", mock.Anything, mock.MatchedBy(func(call ethereum.CallMsg) bool {
			return reflect.DeepEqual(call.Data, callData)
		}), mock.Anything).Return(invalidData, nil)

		_, err = checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unpack contract result")
		mockClient.AssertExpectations(t)
	})
}

// Test for Close method
func TestClose(t *testing.T) {
	mockClient := new(MockEthClient)
	mockClient.On("Close").Return()

	checker := &ValidatorOptInChecker{
		client:      mockClient,
		contractAbi: abi.ABI{},
		address:     common.Address{},
	}

	checker.Close()
	mockClient.AssertExpectations(t)
}
