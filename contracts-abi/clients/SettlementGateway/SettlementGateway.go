// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package settlementgateway

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SettlementgatewayMetaData contains all meta data concerning the Settlementgateway contract.
var SettlementgatewayMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocatorAddr\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"counterpartyFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizationFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeTransfer\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_counterpartyIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_allocatorAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_finalizationFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_counterpartyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateTransfer\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"returnIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFinalizedIdx\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferInitiatedIdx\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferFinalized\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"counterpartyIdx\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferInitiated\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"transferIdx\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AmountTooSmall\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"counterpartyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectEtherValueSent\",\"inputs\":[{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountExpected\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidCounterpartyIndex\",\"inputs\":[{\"name\":\"counterpartyIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transferFinalizedIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotRelayer\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// SettlementgatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use SettlementgatewayMetaData.ABI instead.
var SettlementgatewayABI = SettlementgatewayMetaData.ABI

// Settlementgateway is an auto generated Go binding around an Ethereum contract.
type Settlementgateway struct {
	SettlementgatewayCaller     // Read-only binding to the contract
	SettlementgatewayTransactor // Write-only binding to the contract
	SettlementgatewayFilterer   // Log filterer for contract events
}

// SettlementgatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettlementgatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementgatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettlementgatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementgatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettlementgatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettlementgatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettlementgatewaySession struct {
	Contract     *Settlementgateway // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SettlementgatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettlementgatewayCallerSession struct {
	Contract *SettlementgatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// SettlementgatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettlementgatewayTransactorSession struct {
	Contract     *SettlementgatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SettlementgatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettlementgatewayRaw struct {
	Contract *Settlementgateway // Generic contract binding to access the raw methods on
}

// SettlementgatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettlementgatewayCallerRaw struct {
	Contract *SettlementgatewayCaller // Generic read-only contract binding to access the raw methods on
}

// SettlementgatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettlementgatewayTransactorRaw struct {
	Contract *SettlementgatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettlementgateway creates a new instance of Settlementgateway, bound to a specific deployed contract.
func NewSettlementgateway(address common.Address, backend bind.ContractBackend) (*Settlementgateway, error) {
	contract, err := bindSettlementgateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Settlementgateway{SettlementgatewayCaller: SettlementgatewayCaller{contract: contract}, SettlementgatewayTransactor: SettlementgatewayTransactor{contract: contract}, SettlementgatewayFilterer: SettlementgatewayFilterer{contract: contract}}, nil
}

// NewSettlementgatewayCaller creates a new read-only instance of Settlementgateway, bound to a specific deployed contract.
func NewSettlementgatewayCaller(address common.Address, caller bind.ContractCaller) (*SettlementgatewayCaller, error) {
	contract, err := bindSettlementgateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayCaller{contract: contract}, nil
}

// NewSettlementgatewayTransactor creates a new write-only instance of Settlementgateway, bound to a specific deployed contract.
func NewSettlementgatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*SettlementgatewayTransactor, error) {
	contract, err := bindSettlementgateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayTransactor{contract: contract}, nil
}

// NewSettlementgatewayFilterer creates a new log filterer instance of Settlementgateway, bound to a specific deployed contract.
func NewSettlementgatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*SettlementgatewayFilterer, error) {
	contract, err := bindSettlementgateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayFilterer{contract: contract}, nil
}

// bindSettlementgateway binds a generic wrapper to an already deployed contract.
func bindSettlementgateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SettlementgatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settlementgateway *SettlementgatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Settlementgateway.Contract.SettlementgatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settlementgateway *SettlementgatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.Contract.SettlementgatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settlementgateway *SettlementgatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settlementgateway.Contract.SettlementgatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settlementgateway *SettlementgatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Settlementgateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settlementgateway *SettlementgatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settlementgateway *SettlementgatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settlementgateway.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Settlementgateway *SettlementgatewayCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Settlementgateway *SettlementgatewaySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Settlementgateway.Contract.UPGRADEINTERFACEVERSION(&_Settlementgateway.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Settlementgateway *SettlementgatewayCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Settlementgateway.Contract.UPGRADEINTERFACEVERSION(&_Settlementgateway.CallOpts)
}

// AllocatorAddr is a free data retrieval call binding the contract method 0x1043bf3a.
//
// Solidity: function allocatorAddr() view returns(address)
func (_Settlementgateway *SettlementgatewayCaller) AllocatorAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "allocatorAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AllocatorAddr is a free data retrieval call binding the contract method 0x1043bf3a.
//
// Solidity: function allocatorAddr() view returns(address)
func (_Settlementgateway *SettlementgatewaySession) AllocatorAddr() (common.Address, error) {
	return _Settlementgateway.Contract.AllocatorAddr(&_Settlementgateway.CallOpts)
}

// AllocatorAddr is a free data retrieval call binding the contract method 0x1043bf3a.
//
// Solidity: function allocatorAddr() view returns(address)
func (_Settlementgateway *SettlementgatewayCallerSession) AllocatorAddr() (common.Address, error) {
	return _Settlementgateway.Contract.AllocatorAddr(&_Settlementgateway.CallOpts)
}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCaller) CounterpartyFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "counterpartyFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewaySession) CounterpartyFee() (*big.Int, error) {
	return _Settlementgateway.Contract.CounterpartyFee(&_Settlementgateway.CallOpts)
}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCallerSession) CounterpartyFee() (*big.Int, error) {
	return _Settlementgateway.Contract.CounterpartyFee(&_Settlementgateway.CallOpts)
}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCaller) FinalizationFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "finalizationFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewaySession) FinalizationFee() (*big.Int, error) {
	return _Settlementgateway.Contract.FinalizationFee(&_Settlementgateway.CallOpts)
}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCallerSession) FinalizationFee() (*big.Int, error) {
	return _Settlementgateway.Contract.FinalizationFee(&_Settlementgateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Settlementgateway *SettlementgatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Settlementgateway *SettlementgatewaySession) Owner() (common.Address, error) {
	return _Settlementgateway.Contract.Owner(&_Settlementgateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Settlementgateway *SettlementgatewayCallerSession) Owner() (common.Address, error) {
	return _Settlementgateway.Contract.Owner(&_Settlementgateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Settlementgateway *SettlementgatewayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Settlementgateway *SettlementgatewaySession) Paused() (bool, error) {
	return _Settlementgateway.Contract.Paused(&_Settlementgateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Settlementgateway *SettlementgatewayCallerSession) Paused() (bool, error) {
	return _Settlementgateway.Contract.Paused(&_Settlementgateway.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Settlementgateway *SettlementgatewayCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Settlementgateway *SettlementgatewaySession) PendingOwner() (common.Address, error) {
	return _Settlementgateway.Contract.PendingOwner(&_Settlementgateway.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Settlementgateway *SettlementgatewayCallerSession) PendingOwner() (common.Address, error) {
	return _Settlementgateway.Contract.PendingOwner(&_Settlementgateway.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Settlementgateway *SettlementgatewayCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Settlementgateway *SettlementgatewaySession) ProxiableUUID() ([32]byte, error) {
	return _Settlementgateway.Contract.ProxiableUUID(&_Settlementgateway.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Settlementgateway *SettlementgatewayCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Settlementgateway.Contract.ProxiableUUID(&_Settlementgateway.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Settlementgateway *SettlementgatewayCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "relayer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Settlementgateway *SettlementgatewaySession) Relayer() (common.Address, error) {
	return _Settlementgateway.Contract.Relayer(&_Settlementgateway.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_Settlementgateway *SettlementgatewayCallerSession) Relayer() (common.Address, error) {
	return _Settlementgateway.Contract.Relayer(&_Settlementgateway.CallOpts)
}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCaller) TransferFinalizedIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "transferFinalizedIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewaySession) TransferFinalizedIdx() (*big.Int, error) {
	return _Settlementgateway.Contract.TransferFinalizedIdx(&_Settlementgateway.CallOpts)
}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCallerSession) TransferFinalizedIdx() (*big.Int, error) {
	return _Settlementgateway.Contract.TransferFinalizedIdx(&_Settlementgateway.CallOpts)
}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCaller) TransferInitiatedIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settlementgateway.contract.Call(opts, &out, "transferInitiatedIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewaySession) TransferInitiatedIdx() (*big.Int, error) {
	return _Settlementgateway.Contract.TransferInitiatedIdx(&_Settlementgateway.CallOpts)
}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_Settlementgateway *SettlementgatewayCallerSession) TransferInitiatedIdx() (*big.Int, error) {
	return _Settlementgateway.Contract.TransferInitiatedIdx(&_Settlementgateway.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Settlementgateway *SettlementgatewayTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Settlementgateway *SettlementgatewaySession) AcceptOwnership() (*types.Transaction, error) {
	return _Settlementgateway.Contract.AcceptOwnership(&_Settlementgateway.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Settlementgateway.Contract.AcceptOwnership(&_Settlementgateway.TransactOpts)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_Settlementgateway *SettlementgatewayTransactor) FinalizeTransfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "finalizeTransfer", _recipient, _amount, _counterpartyIdx)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_Settlementgateway *SettlementgatewaySession) FinalizeTransfer(_recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.FinalizeTransfer(&_Settlementgateway.TransactOpts, _recipient, _amount, _counterpartyIdx)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) FinalizeTransfer(_recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.FinalizeTransfer(&_Settlementgateway.TransactOpts, _recipient, _amount, _counterpartyIdx)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6b63eb8.
//
// Solidity: function initialize(address _allocatorAddr, address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_Settlementgateway *SettlementgatewayTransactor) Initialize(opts *bind.TransactOpts, _allocatorAddr common.Address, _owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "initialize", _allocatorAddr, _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6b63eb8.
//
// Solidity: function initialize(address _allocatorAddr, address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_Settlementgateway *SettlementgatewaySession) Initialize(_allocatorAddr common.Address, _owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.Initialize(&_Settlementgateway.TransactOpts, _allocatorAddr, _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6b63eb8.
//
// Solidity: function initialize(address _allocatorAddr, address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) Initialize(_allocatorAddr common.Address, _owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.Initialize(&_Settlementgateway.TransactOpts, _allocatorAddr, _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_Settlementgateway *SettlementgatewayTransactor) InitiateTransfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "initiateTransfer", _recipient, _amount)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_Settlementgateway *SettlementgatewaySession) InitiateTransfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.InitiateTransfer(&_Settlementgateway.TransactOpts, _recipient, _amount)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_Settlementgateway *SettlementgatewayTransactorSession) InitiateTransfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Settlementgateway.Contract.InitiateTransfer(&_Settlementgateway.TransactOpts, _recipient, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Settlementgateway *SettlementgatewayTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Settlementgateway *SettlementgatewaySession) Pause() (*types.Transaction, error) {
	return _Settlementgateway.Contract.Pause(&_Settlementgateway.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) Pause() (*types.Transaction, error) {
	return _Settlementgateway.Contract.Pause(&_Settlementgateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Settlementgateway *SettlementgatewayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Settlementgateway *SettlementgatewaySession) RenounceOwnership() (*types.Transaction, error) {
	return _Settlementgateway.Contract.RenounceOwnership(&_Settlementgateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Settlementgateway.Contract.RenounceOwnership(&_Settlementgateway.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Settlementgateway *SettlementgatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Settlementgateway *SettlementgatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Settlementgateway.Contract.TransferOwnership(&_Settlementgateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Settlementgateway.Contract.TransferOwnership(&_Settlementgateway.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Settlementgateway *SettlementgatewayTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Settlementgateway *SettlementgatewaySession) Unpause() (*types.Transaction, error) {
	return _Settlementgateway.Contract.Unpause(&_Settlementgateway.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) Unpause() (*types.Transaction, error) {
	return _Settlementgateway.Contract.Unpause(&_Settlementgateway.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Settlementgateway *SettlementgatewayTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Settlementgateway.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Settlementgateway *SettlementgatewaySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Settlementgateway.Contract.UpgradeToAndCall(&_Settlementgateway.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Settlementgateway *SettlementgatewayTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Settlementgateway.Contract.UpgradeToAndCall(&_Settlementgateway.TransactOpts, newImplementation, data)
}

// SettlementgatewayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Settlementgateway contract.
type SettlementgatewayInitializedIterator struct {
	Event *SettlementgatewayInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayInitialized represents a Initialized event raised by the Settlementgateway contract.
type SettlementgatewayInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Settlementgateway *SettlementgatewayFilterer) FilterInitialized(opts *bind.FilterOpts) (*SettlementgatewayInitializedIterator, error) {

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayInitializedIterator{contract: _Settlementgateway.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Settlementgateway *SettlementgatewayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SettlementgatewayInitialized) (event.Subscription, error) {

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayInitialized)
				if err := _Settlementgateway.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Settlementgateway *SettlementgatewayFilterer) ParseInitialized(log types.Log) (*SettlementgatewayInitialized, error) {
	event := new(SettlementgatewayInitialized)
	if err := _Settlementgateway.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Settlementgateway contract.
type SettlementgatewayOwnershipTransferStartedIterator struct {
	Event *SettlementgatewayOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Settlementgateway contract.
type SettlementgatewayOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SettlementgatewayOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayOwnershipTransferStartedIterator{contract: _Settlementgateway.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *SettlementgatewayOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayOwnershipTransferStarted)
				if err := _Settlementgateway.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) ParseOwnershipTransferStarted(log types.Log) (*SettlementgatewayOwnershipTransferStarted, error) {
	event := new(SettlementgatewayOwnershipTransferStarted)
	if err := _Settlementgateway.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Settlementgateway contract.
type SettlementgatewayOwnershipTransferredIterator struct {
	Event *SettlementgatewayOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayOwnershipTransferred represents a OwnershipTransferred event raised by the Settlementgateway contract.
type SettlementgatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SettlementgatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayOwnershipTransferredIterator{contract: _Settlementgateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SettlementgatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayOwnershipTransferred)
				if err := _Settlementgateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Settlementgateway *SettlementgatewayFilterer) ParseOwnershipTransferred(log types.Log) (*SettlementgatewayOwnershipTransferred, error) {
	event := new(SettlementgatewayOwnershipTransferred)
	if err := _Settlementgateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Settlementgateway contract.
type SettlementgatewayPausedIterator struct {
	Event *SettlementgatewayPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayPaused represents a Paused event raised by the Settlementgateway contract.
type SettlementgatewayPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) FilterPaused(opts *bind.FilterOpts) (*SettlementgatewayPausedIterator, error) {

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayPausedIterator{contract: _Settlementgateway.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SettlementgatewayPaused) (event.Subscription, error) {

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayPaused)
				if err := _Settlementgateway.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) ParsePaused(log types.Log) (*SettlementgatewayPaused, error) {
	event := new(SettlementgatewayPaused)
	if err := _Settlementgateway.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayTransferFinalizedIterator is returned from FilterTransferFinalized and is used to iterate over the raw logs and unpacked data for TransferFinalized events raised by the Settlementgateway contract.
type SettlementgatewayTransferFinalizedIterator struct {
	Event *SettlementgatewayTransferFinalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayTransferFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayTransferFinalized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayTransferFinalized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayTransferFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayTransferFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayTransferFinalized represents a TransferFinalized event raised by the Settlementgateway contract.
type SettlementgatewayTransferFinalized struct {
	Recipient       common.Address
	Amount          *big.Int
	CounterpartyIdx *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferFinalized is a free log retrieval operation binding the contract event 0x8c1419715bf9fd4efa8c0fd1833ba527cbdd86f6018aa79102af32103bbfdefd.
//
// Solidity: event TransferFinalized(address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx)
func (_Settlementgateway *SettlementgatewayFilterer) FilterTransferFinalized(opts *bind.FilterOpts, recipient []common.Address, counterpartyIdx []*big.Int) (*SettlementgatewayTransferFinalizedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var counterpartyIdxRule []interface{}
	for _, counterpartyIdxItem := range counterpartyIdx {
		counterpartyIdxRule = append(counterpartyIdxRule, counterpartyIdxItem)
	}

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "TransferFinalized", recipientRule, counterpartyIdxRule)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayTransferFinalizedIterator{contract: _Settlementgateway.contract, event: "TransferFinalized", logs: logs, sub: sub}, nil
}

// WatchTransferFinalized is a free log subscription operation binding the contract event 0x8c1419715bf9fd4efa8c0fd1833ba527cbdd86f6018aa79102af32103bbfdefd.
//
// Solidity: event TransferFinalized(address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx)
func (_Settlementgateway *SettlementgatewayFilterer) WatchTransferFinalized(opts *bind.WatchOpts, sink chan<- *SettlementgatewayTransferFinalized, recipient []common.Address, counterpartyIdx []*big.Int) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var counterpartyIdxRule []interface{}
	for _, counterpartyIdxItem := range counterpartyIdx {
		counterpartyIdxRule = append(counterpartyIdxRule, counterpartyIdxItem)
	}

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "TransferFinalized", recipientRule, counterpartyIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayTransferFinalized)
				if err := _Settlementgateway.contract.UnpackLog(event, "TransferFinalized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferFinalized is a log parse operation binding the contract event 0x8c1419715bf9fd4efa8c0fd1833ba527cbdd86f6018aa79102af32103bbfdefd.
//
// Solidity: event TransferFinalized(address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx)
func (_Settlementgateway *SettlementgatewayFilterer) ParseTransferFinalized(log types.Log) (*SettlementgatewayTransferFinalized, error) {
	event := new(SettlementgatewayTransferFinalized)
	if err := _Settlementgateway.contract.UnpackLog(event, "TransferFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayTransferInitiatedIterator is returned from FilterTransferInitiated and is used to iterate over the raw logs and unpacked data for TransferInitiated events raised by the Settlementgateway contract.
type SettlementgatewayTransferInitiatedIterator struct {
	Event *SettlementgatewayTransferInitiated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayTransferInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayTransferInitiated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayTransferInitiated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayTransferInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayTransferInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayTransferInitiated represents a TransferInitiated event raised by the Settlementgateway contract.
type SettlementgatewayTransferInitiated struct {
	Sender      common.Address
	Recipient   common.Address
	Amount      *big.Int
	TransferIdx *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTransferInitiated is a free log retrieval operation binding the contract event 0x6abe792a4e9e702afbc17fdac3c94f6ed1d8c9a8e4917c99672474b3f775ab43.
//
// Solidity: event TransferInitiated(address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx)
func (_Settlementgateway *SettlementgatewayFilterer) FilterTransferInitiated(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, transferIdx []*big.Int) (*SettlementgatewayTransferInitiatedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var transferIdxRule []interface{}
	for _, transferIdxItem := range transferIdx {
		transferIdxRule = append(transferIdxRule, transferIdxItem)
	}

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "TransferInitiated", senderRule, recipientRule, transferIdxRule)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayTransferInitiatedIterator{contract: _Settlementgateway.contract, event: "TransferInitiated", logs: logs, sub: sub}, nil
}

// WatchTransferInitiated is a free log subscription operation binding the contract event 0x6abe792a4e9e702afbc17fdac3c94f6ed1d8c9a8e4917c99672474b3f775ab43.
//
// Solidity: event TransferInitiated(address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx)
func (_Settlementgateway *SettlementgatewayFilterer) WatchTransferInitiated(opts *bind.WatchOpts, sink chan<- *SettlementgatewayTransferInitiated, sender []common.Address, recipient []common.Address, transferIdx []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var transferIdxRule []interface{}
	for _, transferIdxItem := range transferIdx {
		transferIdxRule = append(transferIdxRule, transferIdxItem)
	}

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "TransferInitiated", senderRule, recipientRule, transferIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayTransferInitiated)
				if err := _Settlementgateway.contract.UnpackLog(event, "TransferInitiated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferInitiated is a log parse operation binding the contract event 0x6abe792a4e9e702afbc17fdac3c94f6ed1d8c9a8e4917c99672474b3f775ab43.
//
// Solidity: event TransferInitiated(address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx)
func (_Settlementgateway *SettlementgatewayFilterer) ParseTransferInitiated(log types.Log) (*SettlementgatewayTransferInitiated, error) {
	event := new(SettlementgatewayTransferInitiated)
	if err := _Settlementgateway.contract.UnpackLog(event, "TransferInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Settlementgateway contract.
type SettlementgatewayUnpausedIterator struct {
	Event *SettlementgatewayUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayUnpaused represents a Unpaused event raised by the Settlementgateway contract.
type SettlementgatewayUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SettlementgatewayUnpausedIterator, error) {

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayUnpausedIterator{contract: _Settlementgateway.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SettlementgatewayUnpaused) (event.Subscription, error) {

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayUnpaused)
				if err := _Settlementgateway.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Settlementgateway *SettlementgatewayFilterer) ParseUnpaused(log types.Log) (*SettlementgatewayUnpaused, error) {
	event := new(SettlementgatewayUnpaused)
	if err := _Settlementgateway.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SettlementgatewayUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Settlementgateway contract.
type SettlementgatewayUpgradedIterator struct {
	Event *SettlementgatewayUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SettlementgatewayUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SettlementgatewayUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SettlementgatewayUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SettlementgatewayUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SettlementgatewayUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SettlementgatewayUpgraded represents a Upgraded event raised by the Settlementgateway contract.
type SettlementgatewayUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Settlementgateway *SettlementgatewayFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SettlementgatewayUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Settlementgateway.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SettlementgatewayUpgradedIterator{contract: _Settlementgateway.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Settlementgateway *SettlementgatewayFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SettlementgatewayUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Settlementgateway.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SettlementgatewayUpgraded)
				if err := _Settlementgateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Settlementgateway *SettlementgatewayFilterer) ParseUpgraded(log types.Log) (*SettlementgatewayUpgraded, error) {
	event := new(SettlementgatewayUpgraded)
	if err := _Settlementgateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
