// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatoroptinrouter

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

// ValidatoroptinrouterMetaData contains all meta data concerning the Validatoroptinrouter contract.
var ValidatoroptinrouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"areValidatorsOptedIn\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_vanillaRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_mevCommitAVS\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mevCommitAVS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIMevCommitAVS\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMevCommitAVS\",\"inputs\":[{\"name\":\"_mevCommitAVS\",\"type\":\"address\",\"internalType\":\"contractIMevCommitAVS\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVanillaRegistry\",\"inputs\":[{\"name\":\"_vanillaRegistry\",\"type\":\"address\",\"internalType\":\"contractIVanillaRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"vanillaRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVanillaRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MevCommitAVSSet\",\"inputs\":[{\"name\":\"oldContract\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newContract\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VanillaRegistrySet\",\"inputs\":[{\"name\":\"oldContract\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newContract\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// ValidatoroptinrouterABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidatoroptinrouterMetaData.ABI instead.
var ValidatoroptinrouterABI = ValidatoroptinrouterMetaData.ABI

// Validatoroptinrouter is an auto generated Go binding around an Ethereum contract.
type Validatoroptinrouter struct {
	ValidatoroptinrouterCaller     // Read-only binding to the contract
	ValidatoroptinrouterTransactor // Write-only binding to the contract
	ValidatoroptinrouterFilterer   // Log filterer for contract events
}

// ValidatoroptinrouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatoroptinrouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinrouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatoroptinrouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinrouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatoroptinrouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinrouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatoroptinrouterSession struct {
	Contract     *Validatoroptinrouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ValidatoroptinrouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatoroptinrouterCallerSession struct {
	Contract *ValidatoroptinrouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// ValidatoroptinrouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatoroptinrouterTransactorSession struct {
	Contract     *ValidatoroptinrouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ValidatoroptinrouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatoroptinrouterRaw struct {
	Contract *Validatoroptinrouter // Generic contract binding to access the raw methods on
}

// ValidatoroptinrouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatoroptinrouterCallerRaw struct {
	Contract *ValidatoroptinrouterCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatoroptinrouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatoroptinrouterTransactorRaw struct {
	Contract *ValidatoroptinrouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatoroptinrouter creates a new instance of Validatoroptinrouter, bound to a specific deployed contract.
func NewValidatoroptinrouter(address common.Address, backend bind.ContractBackend) (*Validatoroptinrouter, error) {
	contract, err := bindValidatoroptinrouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validatoroptinrouter{ValidatoroptinrouterCaller: ValidatoroptinrouterCaller{contract: contract}, ValidatoroptinrouterTransactor: ValidatoroptinrouterTransactor{contract: contract}, ValidatoroptinrouterFilterer: ValidatoroptinrouterFilterer{contract: contract}}, nil
}

// NewValidatoroptinrouterCaller creates a new read-only instance of Validatoroptinrouter, bound to a specific deployed contract.
func NewValidatoroptinrouterCaller(address common.Address, caller bind.ContractCaller) (*ValidatoroptinrouterCaller, error) {
	contract, err := bindValidatoroptinrouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterCaller{contract: contract}, nil
}

// NewValidatoroptinrouterTransactor creates a new write-only instance of Validatoroptinrouter, bound to a specific deployed contract.
func NewValidatoroptinrouterTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatoroptinrouterTransactor, error) {
	contract, err := bindValidatoroptinrouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterTransactor{contract: contract}, nil
}

// NewValidatoroptinrouterFilterer creates a new log filterer instance of Validatoroptinrouter, bound to a specific deployed contract.
func NewValidatoroptinrouterFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatoroptinrouterFilterer, error) {
	contract, err := bindValidatoroptinrouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterFilterer{contract: contract}, nil
}

// bindValidatoroptinrouter binds a generic wrapper to an already deployed contract.
func bindValidatoroptinrouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidatoroptinrouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatoroptinrouter *ValidatoroptinrouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatoroptinrouter.Contract.ValidatoroptinrouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatoroptinrouter *ValidatoroptinrouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.ValidatoroptinrouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatoroptinrouter *ValidatoroptinrouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.ValidatoroptinrouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatoroptinrouter *ValidatoroptinrouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatoroptinrouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatoroptinrouter.Contract.UPGRADEINTERFACEVERSION(&_Validatoroptinrouter.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatoroptinrouter.Contract.UPGRADEINTERFACEVERSION(&_Validatoroptinrouter.CallOpts)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) AreValidatorsOptedIn(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([]bool, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "areValidatorsOptedIn", valBLSPubKeys)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinrouter *ValidatoroptinrouterSession) AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error) {
	return _Validatoroptinrouter.Contract.AreValidatorsOptedIn(&_Validatoroptinrouter.CallOpts, valBLSPubKeys)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error) {
	return _Validatoroptinrouter.Contract.AreValidatorsOptedIn(&_Validatoroptinrouter.CallOpts, valBLSPubKeys)
}

// MevCommitAVS is a free data retrieval call binding the contract method 0x22a3d9d6.
//
// Solidity: function mevCommitAVS() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) MevCommitAVS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "mevCommitAVS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MevCommitAVS is a free data retrieval call binding the contract method 0x22a3d9d6.
//
// Solidity: function mevCommitAVS() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) MevCommitAVS() (common.Address, error) {
	return _Validatoroptinrouter.Contract.MevCommitAVS(&_Validatoroptinrouter.CallOpts)
}

// MevCommitAVS is a free data retrieval call binding the contract method 0x22a3d9d6.
//
// Solidity: function mevCommitAVS() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) MevCommitAVS() (common.Address, error) {
	return _Validatoroptinrouter.Contract.MevCommitAVS(&_Validatoroptinrouter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) Owner() (common.Address, error) {
	return _Validatoroptinrouter.Contract.Owner(&_Validatoroptinrouter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) Owner() (common.Address, error) {
	return _Validatoroptinrouter.Contract.Owner(&_Validatoroptinrouter.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) PendingOwner() (common.Address, error) {
	return _Validatoroptinrouter.Contract.PendingOwner(&_Validatoroptinrouter.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) PendingOwner() (common.Address, error) {
	return _Validatoroptinrouter.Contract.PendingOwner(&_Validatoroptinrouter.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) ProxiableUUID() ([32]byte, error) {
	return _Validatoroptinrouter.Contract.ProxiableUUID(&_Validatoroptinrouter.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Validatoroptinrouter.Contract.ProxiableUUID(&_Validatoroptinrouter.CallOpts)
}

// VanillaRegistry is a free data retrieval call binding the contract method 0x6dfade9e.
//
// Solidity: function vanillaRegistry() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCaller) VanillaRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinrouter.contract.Call(opts, &out, "vanillaRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VanillaRegistry is a free data retrieval call binding the contract method 0x6dfade9e.
//
// Solidity: function vanillaRegistry() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterSession) VanillaRegistry() (common.Address, error) {
	return _Validatoroptinrouter.Contract.VanillaRegistry(&_Validatoroptinrouter.CallOpts)
}

// VanillaRegistry is a free data retrieval call binding the contract method 0x6dfade9e.
//
// Solidity: function vanillaRegistry() view returns(address)
func (_Validatoroptinrouter *ValidatoroptinrouterCallerSession) VanillaRegistry() (common.Address, error) {
	return _Validatoroptinrouter.Contract.VanillaRegistry(&_Validatoroptinrouter.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) AcceptOwnership() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.AcceptOwnership(&_Validatoroptinrouter.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.AcceptOwnership(&_Validatoroptinrouter.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _vanillaRegistry, address _mevCommitAVS, address _owner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) Initialize(opts *bind.TransactOpts, _vanillaRegistry common.Address, _mevCommitAVS common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "initialize", _vanillaRegistry, _mevCommitAVS, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _vanillaRegistry, address _mevCommitAVS, address _owner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) Initialize(_vanillaRegistry common.Address, _mevCommitAVS common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Initialize(&_Validatoroptinrouter.TransactOpts, _vanillaRegistry, _mevCommitAVS, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address _vanillaRegistry, address _mevCommitAVS, address _owner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) Initialize(_vanillaRegistry common.Address, _mevCommitAVS common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Initialize(&_Validatoroptinrouter.TransactOpts, _vanillaRegistry, _mevCommitAVS, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.RenounceOwnership(&_Validatoroptinrouter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.RenounceOwnership(&_Validatoroptinrouter.TransactOpts)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address _mevCommitAVS) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) SetMevCommitAVS(opts *bind.TransactOpts, _mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "setMevCommitAVS", _mevCommitAVS)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address _mevCommitAVS) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) SetMevCommitAVS(_mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.SetMevCommitAVS(&_Validatoroptinrouter.TransactOpts, _mevCommitAVS)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address _mevCommitAVS) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) SetMevCommitAVS(_mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.SetMevCommitAVS(&_Validatoroptinrouter.TransactOpts, _mevCommitAVS)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address _vanillaRegistry) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) SetVanillaRegistry(opts *bind.TransactOpts, _vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "setVanillaRegistry", _vanillaRegistry)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address _vanillaRegistry) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) SetVanillaRegistry(_vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.SetVanillaRegistry(&_Validatoroptinrouter.TransactOpts, _vanillaRegistry)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address _vanillaRegistry) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) SetVanillaRegistry(_vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.SetVanillaRegistry(&_Validatoroptinrouter.TransactOpts, _vanillaRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.TransferOwnership(&_Validatoroptinrouter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.TransferOwnership(&_Validatoroptinrouter.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.UpgradeToAndCall(&_Validatoroptinrouter.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.UpgradeToAndCall(&_Validatoroptinrouter.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Fallback(&_Validatoroptinrouter.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Fallback(&_Validatoroptinrouter.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinrouter.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterSession) Receive() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Receive(&_Validatoroptinrouter.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinrouter *ValidatoroptinrouterTransactorSession) Receive() (*types.Transaction, error) {
	return _Validatoroptinrouter.Contract.Receive(&_Validatoroptinrouter.TransactOpts)
}

// ValidatoroptinrouterInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterInitializedIterator struct {
	Event *ValidatoroptinrouterInitialized // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterInitialized)
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
		it.Event = new(ValidatoroptinrouterInitialized)
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
func (it *ValidatoroptinrouterInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterInitialized represents a Initialized event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterInitialized(opts *bind.FilterOpts) (*ValidatoroptinrouterInitializedIterator, error) {

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterInitializedIterator{contract: _Validatoroptinrouter.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterInitialized) (event.Subscription, error) {

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterInitialized)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseInitialized(log types.Log) (*ValidatoroptinrouterInitialized, error) {
	event := new(ValidatoroptinrouterInitialized)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinrouterMevCommitAVSSetIterator is returned from FilterMevCommitAVSSet and is used to iterate over the raw logs and unpacked data for MevCommitAVSSet events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterMevCommitAVSSetIterator struct {
	Event *ValidatoroptinrouterMevCommitAVSSet // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterMevCommitAVSSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterMevCommitAVSSet)
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
		it.Event = new(ValidatoroptinrouterMevCommitAVSSet)
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
func (it *ValidatoroptinrouterMevCommitAVSSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterMevCommitAVSSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterMevCommitAVSSet represents a MevCommitAVSSet event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterMevCommitAVSSet struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMevCommitAVSSet is a free log retrieval operation binding the contract event 0x447ff28b7015145d149b3cd03a9f9e162a197e114c0559fcc47e6ce791f37676.
//
// Solidity: event MevCommitAVSSet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterMevCommitAVSSet(opts *bind.FilterOpts) (*ValidatoroptinrouterMevCommitAVSSetIterator, error) {

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "MevCommitAVSSet")
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterMevCommitAVSSetIterator{contract: _Validatoroptinrouter.contract, event: "MevCommitAVSSet", logs: logs, sub: sub}, nil
}

// WatchMevCommitAVSSet is a free log subscription operation binding the contract event 0x447ff28b7015145d149b3cd03a9f9e162a197e114c0559fcc47e6ce791f37676.
//
// Solidity: event MevCommitAVSSet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchMevCommitAVSSet(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterMevCommitAVSSet) (event.Subscription, error) {

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "MevCommitAVSSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterMevCommitAVSSet)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "MevCommitAVSSet", log); err != nil {
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

// ParseMevCommitAVSSet is a log parse operation binding the contract event 0x447ff28b7015145d149b3cd03a9f9e162a197e114c0559fcc47e6ce791f37676.
//
// Solidity: event MevCommitAVSSet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseMevCommitAVSSet(log types.Log) (*ValidatoroptinrouterMevCommitAVSSet, error) {
	event := new(ValidatoroptinrouterMevCommitAVSSet)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "MevCommitAVSSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinrouterOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterOwnershipTransferStartedIterator struct {
	Event *ValidatoroptinrouterOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterOwnershipTransferStarted)
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
		it.Event = new(ValidatoroptinrouterOwnershipTransferStarted)
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
func (it *ValidatoroptinrouterOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatoroptinrouterOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterOwnershipTransferStartedIterator{contract: _Validatoroptinrouter.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterOwnershipTransferStarted)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseOwnershipTransferStarted(log types.Log) (*ValidatoroptinrouterOwnershipTransferStarted, error) {
	event := new(ValidatoroptinrouterOwnershipTransferStarted)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinrouterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterOwnershipTransferredIterator struct {
	Event *ValidatoroptinrouterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterOwnershipTransferred)
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
		it.Event = new(ValidatoroptinrouterOwnershipTransferred)
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
func (it *ValidatoroptinrouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterOwnershipTransferred represents a OwnershipTransferred event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatoroptinrouterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterOwnershipTransferredIterator{contract: _Validatoroptinrouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterOwnershipTransferred)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseOwnershipTransferred(log types.Log) (*ValidatoroptinrouterOwnershipTransferred, error) {
	event := new(ValidatoroptinrouterOwnershipTransferred)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinrouterUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterUpgradedIterator struct {
	Event *ValidatoroptinrouterUpgraded // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterUpgraded)
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
		it.Event = new(ValidatoroptinrouterUpgraded)
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
func (it *ValidatoroptinrouterUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterUpgraded represents a Upgraded event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ValidatoroptinrouterUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterUpgradedIterator{contract: _Validatoroptinrouter.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterUpgraded)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseUpgraded(log types.Log) (*ValidatoroptinrouterUpgraded, error) {
	event := new(ValidatoroptinrouterUpgraded)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinrouterVanillaRegistrySetIterator is returned from FilterVanillaRegistrySet and is used to iterate over the raw logs and unpacked data for VanillaRegistrySet events raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterVanillaRegistrySetIterator struct {
	Event *ValidatoroptinrouterVanillaRegistrySet // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinrouterVanillaRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinrouterVanillaRegistrySet)
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
		it.Event = new(ValidatoroptinrouterVanillaRegistrySet)
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
func (it *ValidatoroptinrouterVanillaRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinrouterVanillaRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinrouterVanillaRegistrySet represents a VanillaRegistrySet event raised by the Validatoroptinrouter contract.
type ValidatoroptinrouterVanillaRegistrySet struct {
	OldContract common.Address
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVanillaRegistrySet is a free log retrieval operation binding the contract event 0x3f45fe29af7577363f9fe5c59df3712e84086c072bd5fca14a1c76230b19c667.
//
// Solidity: event VanillaRegistrySet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) FilterVanillaRegistrySet(opts *bind.FilterOpts) (*ValidatoroptinrouterVanillaRegistrySetIterator, error) {

	logs, sub, err := _Validatoroptinrouter.contract.FilterLogs(opts, "VanillaRegistrySet")
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinrouterVanillaRegistrySetIterator{contract: _Validatoroptinrouter.contract, event: "VanillaRegistrySet", logs: logs, sub: sub}, nil
}

// WatchVanillaRegistrySet is a free log subscription operation binding the contract event 0x3f45fe29af7577363f9fe5c59df3712e84086c072bd5fca14a1c76230b19c667.
//
// Solidity: event VanillaRegistrySet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) WatchVanillaRegistrySet(opts *bind.WatchOpts, sink chan<- *ValidatoroptinrouterVanillaRegistrySet) (event.Subscription, error) {

	logs, sub, err := _Validatoroptinrouter.contract.WatchLogs(opts, "VanillaRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinrouterVanillaRegistrySet)
				if err := _Validatoroptinrouter.contract.UnpackLog(event, "VanillaRegistrySet", log); err != nil {
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

// ParseVanillaRegistrySet is a log parse operation binding the contract event 0x3f45fe29af7577363f9fe5c59df3712e84086c072bd5fca14a1c76230b19c667.
//
// Solidity: event VanillaRegistrySet(address oldContract, address newContract)
func (_Validatoroptinrouter *ValidatoroptinrouterFilterer) ParseVanillaRegistrySet(log types.Log) (*ValidatoroptinrouterVanillaRegistrySet, error) {
	event := new(ValidatoroptinrouterVanillaRegistrySet)
	if err := _Validatoroptinrouter.contract.UnpackLog(event, "VanillaRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
