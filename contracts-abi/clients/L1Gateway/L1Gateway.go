// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package l1gateway

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

// L1gatewayMetaData contains all meta data concerning the L1gateway contract.
var L1gatewayMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"counterpartyFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizationFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"finalizeTransfer\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_counterpartyIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_relayer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_finalizationFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_counterpartyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateTransfer\",\"inputs\":[{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"returnIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relayer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFinalizedIdx\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferInitiatedIdx\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferFinalized\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"counterpartyIdx\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferInitiated\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"transferIdx\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AmountTooSmall\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"counterpartyFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectEtherValueSent\",\"inputs\":[{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountExpected\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientContractBalance\",\"inputs\":[{\"name\":\"thisContractBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRequested\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidCounterpartyIndex\",\"inputs\":[{\"name\":\"counterpartyIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"transferFinalizedIdx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotRelayer\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"relayer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// L1gatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use L1gatewayMetaData.ABI instead.
var L1gatewayABI = L1gatewayMetaData.ABI

// L1gateway is an auto generated Go binding around an Ethereum contract.
type L1gateway struct {
	L1gatewayCaller     // Read-only binding to the contract
	L1gatewayTransactor // Write-only binding to the contract
	L1gatewayFilterer   // Log filterer for contract events
}

// L1gatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type L1gatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1gatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L1gatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1gatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L1gatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L1gatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L1gatewaySession struct {
	Contract     *L1gateway        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// L1gatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L1gatewayCallerSession struct {
	Contract *L1gatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// L1gatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L1gatewayTransactorSession struct {
	Contract     *L1gatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// L1gatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type L1gatewayRaw struct {
	Contract *L1gateway // Generic contract binding to access the raw methods on
}

// L1gatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L1gatewayCallerRaw struct {
	Contract *L1gatewayCaller // Generic read-only contract binding to access the raw methods on
}

// L1gatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L1gatewayTransactorRaw struct {
	Contract *L1gatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL1gateway creates a new instance of L1gateway, bound to a specific deployed contract.
func NewL1gateway(address common.Address, backend bind.ContractBackend) (*L1gateway, error) {
	contract, err := bindL1gateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L1gateway{L1gatewayCaller: L1gatewayCaller{contract: contract}, L1gatewayTransactor: L1gatewayTransactor{contract: contract}, L1gatewayFilterer: L1gatewayFilterer{contract: contract}}, nil
}

// NewL1gatewayCaller creates a new read-only instance of L1gateway, bound to a specific deployed contract.
func NewL1gatewayCaller(address common.Address, caller bind.ContractCaller) (*L1gatewayCaller, error) {
	contract, err := bindL1gateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L1gatewayCaller{contract: contract}, nil
}

// NewL1gatewayTransactor creates a new write-only instance of L1gateway, bound to a specific deployed contract.
func NewL1gatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*L1gatewayTransactor, error) {
	contract, err := bindL1gateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L1gatewayTransactor{contract: contract}, nil
}

// NewL1gatewayFilterer creates a new log filterer instance of L1gateway, bound to a specific deployed contract.
func NewL1gatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*L1gatewayFilterer, error) {
	contract, err := bindL1gateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L1gatewayFilterer{contract: contract}, nil
}

// bindL1gateway binds a generic wrapper to an already deployed contract.
func bindL1gateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L1gatewayMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1gateway *L1gatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1gateway.Contract.L1gatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1gateway *L1gatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.Contract.L1gatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1gateway *L1gatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1gateway.Contract.L1gatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L1gateway *L1gatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L1gateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L1gateway *L1gatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L1gateway *L1gatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L1gateway.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_L1gateway *L1gatewayCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_L1gateway *L1gatewaySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _L1gateway.Contract.UPGRADEINTERFACEVERSION(&_L1gateway.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_L1gateway *L1gatewayCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _L1gateway.Contract.UPGRADEINTERFACEVERSION(&_L1gateway.CallOpts)
}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_L1gateway *L1gatewayCaller) CounterpartyFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "counterpartyFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_L1gateway *L1gatewaySession) CounterpartyFee() (*big.Int, error) {
	return _L1gateway.Contract.CounterpartyFee(&_L1gateway.CallOpts)
}

// CounterpartyFee is a free data retrieval call binding the contract method 0x97599011.
//
// Solidity: function counterpartyFee() view returns(uint256)
func (_L1gateway *L1gatewayCallerSession) CounterpartyFee() (*big.Int, error) {
	return _L1gateway.Contract.CounterpartyFee(&_L1gateway.CallOpts)
}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_L1gateway *L1gatewayCaller) FinalizationFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "finalizationFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_L1gateway *L1gatewaySession) FinalizationFee() (*big.Int, error) {
	return _L1gateway.Contract.FinalizationFee(&_L1gateway.CallOpts)
}

// FinalizationFee is a free data retrieval call binding the contract method 0x78d3d576.
//
// Solidity: function finalizationFee() view returns(uint256)
func (_L1gateway *L1gatewayCallerSession) FinalizationFee() (*big.Int, error) {
	return _L1gateway.Contract.FinalizationFee(&_L1gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1gateway *L1gatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1gateway *L1gatewaySession) Owner() (common.Address, error) {
	return _L1gateway.Contract.Owner(&_L1gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L1gateway *L1gatewayCallerSession) Owner() (common.Address, error) {
	return _L1gateway.Contract.Owner(&_L1gateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1gateway *L1gatewayCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1gateway *L1gatewaySession) Paused() (bool, error) {
	return _L1gateway.Contract.Paused(&_L1gateway.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_L1gateway *L1gatewayCallerSession) Paused() (bool, error) {
	return _L1gateway.Contract.Paused(&_L1gateway.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_L1gateway *L1gatewayCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_L1gateway *L1gatewaySession) PendingOwner() (common.Address, error) {
	return _L1gateway.Contract.PendingOwner(&_L1gateway.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_L1gateway *L1gatewayCallerSession) PendingOwner() (common.Address, error) {
	return _L1gateway.Contract.PendingOwner(&_L1gateway.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1gateway *L1gatewayCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1gateway *L1gatewaySession) ProxiableUUID() ([32]byte, error) {
	return _L1gateway.Contract.ProxiableUUID(&_L1gateway.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_L1gateway *L1gatewayCallerSession) ProxiableUUID() ([32]byte, error) {
	return _L1gateway.Contract.ProxiableUUID(&_L1gateway.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_L1gateway *L1gatewayCaller) Relayer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "relayer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_L1gateway *L1gatewaySession) Relayer() (common.Address, error) {
	return _L1gateway.Contract.Relayer(&_L1gateway.CallOpts)
}

// Relayer is a free data retrieval call binding the contract method 0x8406c079.
//
// Solidity: function relayer() view returns(address)
func (_L1gateway *L1gatewayCallerSession) Relayer() (common.Address, error) {
	return _L1gateway.Contract.Relayer(&_L1gateway.CallOpts)
}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_L1gateway *L1gatewayCaller) TransferFinalizedIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "transferFinalizedIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_L1gateway *L1gatewaySession) TransferFinalizedIdx() (*big.Int, error) {
	return _L1gateway.Contract.TransferFinalizedIdx(&_L1gateway.CallOpts)
}

// TransferFinalizedIdx is a free data retrieval call binding the contract method 0xa2ff158d.
//
// Solidity: function transferFinalizedIdx() view returns(uint256)
func (_L1gateway *L1gatewayCallerSession) TransferFinalizedIdx() (*big.Int, error) {
	return _L1gateway.Contract.TransferFinalizedIdx(&_L1gateway.CallOpts)
}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_L1gateway *L1gatewayCaller) TransferInitiatedIdx(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L1gateway.contract.Call(opts, &out, "transferInitiatedIdx")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_L1gateway *L1gatewaySession) TransferInitiatedIdx() (*big.Int, error) {
	return _L1gateway.Contract.TransferInitiatedIdx(&_L1gateway.CallOpts)
}

// TransferInitiatedIdx is a free data retrieval call binding the contract method 0xe557b142.
//
// Solidity: function transferInitiatedIdx() view returns(uint256)
func (_L1gateway *L1gatewayCallerSession) TransferInitiatedIdx() (*big.Int, error) {
	return _L1gateway.Contract.TransferInitiatedIdx(&_L1gateway.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_L1gateway *L1gatewayTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_L1gateway *L1gatewaySession) AcceptOwnership() (*types.Transaction, error) {
	return _L1gateway.Contract.AcceptOwnership(&_L1gateway.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_L1gateway *L1gatewayTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _L1gateway.Contract.AcceptOwnership(&_L1gateway.TransactOpts)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_L1gateway *L1gatewayTransactor) FinalizeTransfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "finalizeTransfer", _recipient, _amount, _counterpartyIdx)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_L1gateway *L1gatewaySession) FinalizeTransfer(_recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.FinalizeTransfer(&_L1gateway.TransactOpts, _recipient, _amount, _counterpartyIdx)
}

// FinalizeTransfer is a paid mutator transaction binding the contract method 0xc40a7c82.
//
// Solidity: function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) returns()
func (_L1gateway *L1gatewayTransactorSession) FinalizeTransfer(_recipient common.Address, _amount *big.Int, _counterpartyIdx *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.FinalizeTransfer(&_L1gateway.TransactOpts, _recipient, _amount, _counterpartyIdx)
}

// Initialize is a paid mutator transaction binding the contract method 0xeb990c59.
//
// Solidity: function initialize(address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_L1gateway *L1gatewayTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "initialize", _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xeb990c59.
//
// Solidity: function initialize(address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_L1gateway *L1gatewaySession) Initialize(_owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.Initialize(&_L1gateway.TransactOpts, _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// Initialize is a paid mutator transaction binding the contract method 0xeb990c59.
//
// Solidity: function initialize(address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee) returns()
func (_L1gateway *L1gatewayTransactorSession) Initialize(_owner common.Address, _relayer common.Address, _finalizationFee *big.Int, _counterpartyFee *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.Initialize(&_L1gateway.TransactOpts, _owner, _relayer, _finalizationFee, _counterpartyFee)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_L1gateway *L1gatewayTransactor) InitiateTransfer(opts *bind.TransactOpts, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "initiateTransfer", _recipient, _amount)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_L1gateway *L1gatewaySession) InitiateTransfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.InitiateTransfer(&_L1gateway.TransactOpts, _recipient, _amount)
}

// InitiateTransfer is a paid mutator transaction binding the contract method 0xb504cd1e.
//
// Solidity: function initiateTransfer(address _recipient, uint256 _amount) payable returns(uint256 returnIdx)
func (_L1gateway *L1gatewayTransactorSession) InitiateTransfer(_recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _L1gateway.Contract.InitiateTransfer(&_L1gateway.TransactOpts, _recipient, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1gateway *L1gatewayTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1gateway *L1gatewaySession) Pause() (*types.Transaction, error) {
	return _L1gateway.Contract.Pause(&_L1gateway.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_L1gateway *L1gatewayTransactorSession) Pause() (*types.Transaction, error) {
	return _L1gateway.Contract.Pause(&_L1gateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1gateway *L1gatewayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1gateway *L1gatewaySession) RenounceOwnership() (*types.Transaction, error) {
	return _L1gateway.Contract.RenounceOwnership(&_L1gateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L1gateway *L1gatewayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L1gateway.Contract.RenounceOwnership(&_L1gateway.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1gateway *L1gatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1gateway *L1gatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1gateway.Contract.TransferOwnership(&_L1gateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L1gateway *L1gatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L1gateway.Contract.TransferOwnership(&_L1gateway.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1gateway *L1gatewayTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1gateway *L1gatewaySession) Unpause() (*types.Transaction, error) {
	return _L1gateway.Contract.Unpause(&_L1gateway.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_L1gateway *L1gatewayTransactorSession) Unpause() (*types.Transaction, error) {
	return _L1gateway.Contract.Unpause(&_L1gateway.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1gateway *L1gatewayTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1gateway.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1gateway *L1gatewaySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1gateway.Contract.UpgradeToAndCall(&_L1gateway.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_L1gateway *L1gatewayTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _L1gateway.Contract.UpgradeToAndCall(&_L1gateway.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1gateway *L1gatewayTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _L1gateway.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1gateway *L1gatewaySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L1gateway.Contract.Fallback(&_L1gateway.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L1gateway *L1gatewayTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L1gateway.Contract.Fallback(&_L1gateway.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1gateway *L1gatewayTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L1gateway.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1gateway *L1gatewaySession) Receive() (*types.Transaction, error) {
	return _L1gateway.Contract.Receive(&_L1gateway.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L1gateway *L1gatewayTransactorSession) Receive() (*types.Transaction, error) {
	return _L1gateway.Contract.Receive(&_L1gateway.TransactOpts)
}

// L1gatewayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the L1gateway contract.
type L1gatewayInitializedIterator struct {
	Event *L1gatewayInitialized // Event containing the contract specifics and raw log

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
func (it *L1gatewayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayInitialized)
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
		it.Event = new(L1gatewayInitialized)
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
func (it *L1gatewayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayInitialized represents a Initialized event raised by the L1gateway contract.
type L1gatewayInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_L1gateway *L1gatewayFilterer) FilterInitialized(opts *bind.FilterOpts) (*L1gatewayInitializedIterator, error) {

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &L1gatewayInitializedIterator{contract: _L1gateway.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_L1gateway *L1gatewayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *L1gatewayInitialized) (event.Subscription, error) {

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayInitialized)
				if err := _L1gateway.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseInitialized(log types.Log) (*L1gatewayInitialized, error) {
	event := new(L1gatewayInitialized)
	if err := _L1gateway.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the L1gateway contract.
type L1gatewayOwnershipTransferStartedIterator struct {
	Event *L1gatewayOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *L1gatewayOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayOwnershipTransferStarted)
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
		it.Event = new(L1gatewayOwnershipTransferStarted)
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
func (it *L1gatewayOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the L1gateway contract.
type L1gatewayOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_L1gateway *L1gatewayFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1gatewayOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1gatewayOwnershipTransferStartedIterator{contract: _L1gateway.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_L1gateway *L1gatewayFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *L1gatewayOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayOwnershipTransferStarted)
				if err := _L1gateway.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseOwnershipTransferStarted(log types.Log) (*L1gatewayOwnershipTransferStarted, error) {
	event := new(L1gatewayOwnershipTransferStarted)
	if err := _L1gateway.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L1gateway contract.
type L1gatewayOwnershipTransferredIterator struct {
	Event *L1gatewayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L1gatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayOwnershipTransferred)
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
		it.Event = new(L1gatewayOwnershipTransferred)
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
func (it *L1gatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayOwnershipTransferred represents a OwnershipTransferred event raised by the L1gateway contract.
type L1gatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1gateway *L1gatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L1gatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L1gatewayOwnershipTransferredIterator{contract: _L1gateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L1gateway *L1gatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L1gatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayOwnershipTransferred)
				if err := _L1gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseOwnershipTransferred(log types.Log) (*L1gatewayOwnershipTransferred, error) {
	event := new(L1gatewayOwnershipTransferred)
	if err := _L1gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the L1gateway contract.
type L1gatewayPausedIterator struct {
	Event *L1gatewayPaused // Event containing the contract specifics and raw log

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
func (it *L1gatewayPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayPaused)
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
		it.Event = new(L1gatewayPaused)
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
func (it *L1gatewayPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayPaused represents a Paused event raised by the L1gateway contract.
type L1gatewayPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_L1gateway *L1gatewayFilterer) FilterPaused(opts *bind.FilterOpts) (*L1gatewayPausedIterator, error) {

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &L1gatewayPausedIterator{contract: _L1gateway.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_L1gateway *L1gatewayFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *L1gatewayPaused) (event.Subscription, error) {

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayPaused)
				if err := _L1gateway.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParsePaused(log types.Log) (*L1gatewayPaused, error) {
	event := new(L1gatewayPaused)
	if err := _L1gateway.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayTransferFinalizedIterator is returned from FilterTransferFinalized and is used to iterate over the raw logs and unpacked data for TransferFinalized events raised by the L1gateway contract.
type L1gatewayTransferFinalizedIterator struct {
	Event *L1gatewayTransferFinalized // Event containing the contract specifics and raw log

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
func (it *L1gatewayTransferFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayTransferFinalized)
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
		it.Event = new(L1gatewayTransferFinalized)
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
func (it *L1gatewayTransferFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayTransferFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayTransferFinalized represents a TransferFinalized event raised by the L1gateway contract.
type L1gatewayTransferFinalized struct {
	Recipient       common.Address
	Amount          *big.Int
	CounterpartyIdx *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferFinalized is a free log retrieval operation binding the contract event 0x8c1419715bf9fd4efa8c0fd1833ba527cbdd86f6018aa79102af32103bbfdefd.
//
// Solidity: event TransferFinalized(address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx)
func (_L1gateway *L1gatewayFilterer) FilterTransferFinalized(opts *bind.FilterOpts, recipient []common.Address, counterpartyIdx []*big.Int) (*L1gatewayTransferFinalizedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var counterpartyIdxRule []interface{}
	for _, counterpartyIdxItem := range counterpartyIdx {
		counterpartyIdxRule = append(counterpartyIdxRule, counterpartyIdxItem)
	}

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "TransferFinalized", recipientRule, counterpartyIdxRule)
	if err != nil {
		return nil, err
	}
	return &L1gatewayTransferFinalizedIterator{contract: _L1gateway.contract, event: "TransferFinalized", logs: logs, sub: sub}, nil
}

// WatchTransferFinalized is a free log subscription operation binding the contract event 0x8c1419715bf9fd4efa8c0fd1833ba527cbdd86f6018aa79102af32103bbfdefd.
//
// Solidity: event TransferFinalized(address indexed recipient, uint256 amount, uint256 indexed counterpartyIdx)
func (_L1gateway *L1gatewayFilterer) WatchTransferFinalized(opts *bind.WatchOpts, sink chan<- *L1gatewayTransferFinalized, recipient []common.Address, counterpartyIdx []*big.Int) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	var counterpartyIdxRule []interface{}
	for _, counterpartyIdxItem := range counterpartyIdx {
		counterpartyIdxRule = append(counterpartyIdxRule, counterpartyIdxItem)
	}

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "TransferFinalized", recipientRule, counterpartyIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayTransferFinalized)
				if err := _L1gateway.contract.UnpackLog(event, "TransferFinalized", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseTransferFinalized(log types.Log) (*L1gatewayTransferFinalized, error) {
	event := new(L1gatewayTransferFinalized)
	if err := _L1gateway.contract.UnpackLog(event, "TransferFinalized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayTransferInitiatedIterator is returned from FilterTransferInitiated and is used to iterate over the raw logs and unpacked data for TransferInitiated events raised by the L1gateway contract.
type L1gatewayTransferInitiatedIterator struct {
	Event *L1gatewayTransferInitiated // Event containing the contract specifics and raw log

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
func (it *L1gatewayTransferInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayTransferInitiated)
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
		it.Event = new(L1gatewayTransferInitiated)
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
func (it *L1gatewayTransferInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayTransferInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayTransferInitiated represents a TransferInitiated event raised by the L1gateway contract.
type L1gatewayTransferInitiated struct {
	Sender      common.Address
	Recipient   common.Address
	Amount      *big.Int
	TransferIdx *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTransferInitiated is a free log retrieval operation binding the contract event 0x6abe792a4e9e702afbc17fdac3c94f6ed1d8c9a8e4917c99672474b3f775ab43.
//
// Solidity: event TransferInitiated(address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx)
func (_L1gateway *L1gatewayFilterer) FilterTransferInitiated(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, transferIdx []*big.Int) (*L1gatewayTransferInitiatedIterator, error) {

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

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "TransferInitiated", senderRule, recipientRule, transferIdxRule)
	if err != nil {
		return nil, err
	}
	return &L1gatewayTransferInitiatedIterator{contract: _L1gateway.contract, event: "TransferInitiated", logs: logs, sub: sub}, nil
}

// WatchTransferInitiated is a free log subscription operation binding the contract event 0x6abe792a4e9e702afbc17fdac3c94f6ed1d8c9a8e4917c99672474b3f775ab43.
//
// Solidity: event TransferInitiated(address indexed sender, address indexed recipient, uint256 amount, uint256 indexed transferIdx)
func (_L1gateway *L1gatewayFilterer) WatchTransferInitiated(opts *bind.WatchOpts, sink chan<- *L1gatewayTransferInitiated, sender []common.Address, recipient []common.Address, transferIdx []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "TransferInitiated", senderRule, recipientRule, transferIdxRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayTransferInitiated)
				if err := _L1gateway.contract.UnpackLog(event, "TransferInitiated", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseTransferInitiated(log types.Log) (*L1gatewayTransferInitiated, error) {
	event := new(L1gatewayTransferInitiated)
	if err := _L1gateway.contract.UnpackLog(event, "TransferInitiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the L1gateway contract.
type L1gatewayUnpausedIterator struct {
	Event *L1gatewayUnpaused // Event containing the contract specifics and raw log

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
func (it *L1gatewayUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayUnpaused)
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
		it.Event = new(L1gatewayUnpaused)
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
func (it *L1gatewayUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayUnpaused represents a Unpaused event raised by the L1gateway contract.
type L1gatewayUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_L1gateway *L1gatewayFilterer) FilterUnpaused(opts *bind.FilterOpts) (*L1gatewayUnpausedIterator, error) {

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &L1gatewayUnpausedIterator{contract: _L1gateway.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_L1gateway *L1gatewayFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *L1gatewayUnpaused) (event.Subscription, error) {

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayUnpaused)
				if err := _L1gateway.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseUnpaused(log types.Log) (*L1gatewayUnpaused, error) {
	event := new(L1gatewayUnpaused)
	if err := _L1gateway.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L1gatewayUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the L1gateway contract.
type L1gatewayUpgradedIterator struct {
	Event *L1gatewayUpgraded // Event containing the contract specifics and raw log

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
func (it *L1gatewayUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L1gatewayUpgraded)
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
		it.Event = new(L1gatewayUpgraded)
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
func (it *L1gatewayUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L1gatewayUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L1gatewayUpgraded represents a Upgraded event raised by the L1gateway contract.
type L1gatewayUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1gateway *L1gatewayFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*L1gatewayUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1gateway.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &L1gatewayUpgradedIterator{contract: _L1gateway.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_L1gateway *L1gatewayFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *L1gatewayUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _L1gateway.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L1gatewayUpgraded)
				if err := _L1gateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_L1gateway *L1gatewayFilterer) ParseUpgraded(log types.Log) (*L1gatewayUpgraded, error) {
	event := new(L1gatewayUpgraded)
	if err := _L1gateway.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
