package contract

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRouterContract mocks the RouterContract interface
type MockRouterContract struct {
	mock.Mock
}

func (m *MockRouterContract) AreValidatorsOptedIn(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([]validatorrouter.IValidatorOptInRouterOptInStatus, error) {
	args := m.Called(opts, valBLSPubKeys)
	return args.Get(0).([]validatorrouter.IValidatorOptInRouterOptInStatus), args.Error(1)
}

// MockEthClient mocks the EthClient interface
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockEthClient) Close() {
	m.Called()
}

// Test for CheckValidatorsOptedIn
func TestCheckValidatorsOptedIn(t *testing.T) {
	mockClient := new(MockEthClient)
	mockRouter := new(MockRouterContract)

	mockClient.On("BlockNumber", mock.Anything).Return(uint64(100), nil)
	mockClient.On("Close").Return()

	checker := &ValidatorOptInChecker{
		client:         mockClient,
		routerContract: mockRouter,
		optsGetter: func() (*bind.CallOpts, error) {
			return &bind.CallOpts{BlockNumber: big.NewInt(90)}, nil
		},
	}

	pubkeys := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
	}

	decodedPubKey1, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	decodedPubKey2, _ := hex.DecodeString("abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	expectedBlsPubKeys := [][]byte{decodedPubKey1, decodedPubKey2}

	expectedStatuses := []validatorrouter.IValidatorOptInRouterOptInStatus{
		{IsVanillaOptedIn: true, IsAvsOptedIn: false, IsMiddlewareOptedIn: true},
		{IsVanillaOptedIn: false, IsAvsOptedIn: true, IsMiddlewareOptedIn: false},
	}

	mockRouter.On("AreValidatorsOptedIn", mock.MatchedBy(func(opts *bind.CallOpts) bool {
		return opts.BlockNumber.Int64() == int64(90)
	}), mock.MatchedBy(func(keys [][]byte) bool {
		if len(keys) != len(expectedBlsPubKeys) {
			return false
		}
		for i, key := range keys {
			if string(key) != string(expectedBlsPubKeys[i]) {
				return false
			}
		}
		return true
	})).Return(expectedStatuses, nil)

	ctx := context.Background()
	statuses, err := checker.CheckValidatorsOptedIn(ctx, pubkeys)

	assert.NoError(t, err)
	assert.Equal(t, expectedStatuses, statuses)
	mockRouter.AssertExpectations(t)
}

// Test error handling for CheckValidatorsOptedIn
func TestCheckValidatorsOptedIn_Errors(t *testing.T) {
	mockClient := new(MockEthClient)
	mockRouter := new(MockRouterContract)

	mockClient.On("BlockNumber", mock.Anything).Return(uint64(100), nil)
	mockClient.On("Close").Return()

	t.Run("Invalid pubkey format", func(t *testing.T) {
		checker := &ValidatorOptInChecker{
			client:         mockClient,
			routerContract: mockRouter,
			optsGetter: func() (*bind.CallOpts, error) {
				return &bind.CallOpts{BlockNumber: big.NewInt(90)}, nil
			},
		}

		pubkeys := []string{"invalid_pubkey"}
		_, err := checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error decoding pubkey")
	})

	t.Run("Contract call failure", func(t *testing.T) {
		checker := &ValidatorOptInChecker{
			client:         mockClient,
			routerContract: mockRouter,
			optsGetter: func() (*bind.CallOpts, error) {
				return &bind.CallOpts{BlockNumber: big.NewInt(90)}, nil
			},
		}

		pubkeys := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}
		decodedPubKey, _ := hex.DecodeString("1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		expectedBlsPubKeys := [][]byte{decodedPubKey}

		mockRouter.On("AreValidatorsOptedIn", mock.Anything, mock.MatchedBy(func(keys [][]byte) bool {
			return string(keys[0]) == string(expectedBlsPubKeys[0])
		})).Return([]validatorrouter.IValidatorOptInRouterOptInStatus{}, errors.New("contract call failed"))

		_, err := checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to call contract")
		mockRouter.AssertExpectations(t)
	})

	t.Run("optsGetter failure", func(t *testing.T) {
		checker := &ValidatorOptInChecker{
			client:         mockClient,
			routerContract: mockRouter,
			optsGetter: func() (*bind.CallOpts, error) {
				return nil, errors.New("failed to get block number")
			},
		}

		pubkeys := []string{"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"}

		_, err := checker.CheckValidatorsOptedIn(context.Background(), pubkeys)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "getting call opts")
	})
}

// Test for Close method
func TestClose(t *testing.T) {
	mockClient := new(MockEthClient)
	mockClient.On("Close").Return()

	checker := &ValidatorOptInChecker{
		client: mockClient,
	}

	checker.Close()
	mockClient.AssertExpectations(t)
}
