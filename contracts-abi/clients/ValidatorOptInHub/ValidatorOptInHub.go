// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatoroptinhub

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

// ValidatoroptinhubMetaData contains all meta data concerning the Validatoroptinhub contract.
var ValidatoroptinhubMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addRegistry\",\"inputs\":[{\"name\":\"registry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"areValidatorsOptedIn\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"areValidatorsOptedInList\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool[][]\",\"internalType\":\"bool[][]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_registries\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorOptedInList\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registries\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeRegistry\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"registry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRegistry\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oldRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistryAdded\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"registry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistryRemoved\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"registry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RegistryReplaced\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexRegistryMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidIndex\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// ValidatoroptinhubABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidatoroptinhubMetaData.ABI instead.
var ValidatoroptinhubABI = ValidatoroptinhubMetaData.ABI

// Validatoroptinhub is an auto generated Go binding around an Ethereum contract.
type Validatoroptinhub struct {
	ValidatoroptinhubCaller     // Read-only binding to the contract
	ValidatoroptinhubTransactor // Write-only binding to the contract
	ValidatoroptinhubFilterer   // Log filterer for contract events
}

// ValidatoroptinhubCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatoroptinhubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinhubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatoroptinhubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinhubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatoroptinhubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatoroptinhubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatoroptinhubSession struct {
	Contract     *Validatoroptinhub // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ValidatoroptinhubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatoroptinhubCallerSession struct {
	Contract *ValidatoroptinhubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ValidatoroptinhubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatoroptinhubTransactorSession struct {
	Contract     *ValidatoroptinhubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ValidatoroptinhubRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatoroptinhubRaw struct {
	Contract *Validatoroptinhub // Generic contract binding to access the raw methods on
}

// ValidatoroptinhubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatoroptinhubCallerRaw struct {
	Contract *ValidatoroptinhubCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatoroptinhubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatoroptinhubTransactorRaw struct {
	Contract *ValidatoroptinhubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatoroptinhub creates a new instance of Validatoroptinhub, bound to a specific deployed contract.
func NewValidatoroptinhub(address common.Address, backend bind.ContractBackend) (*Validatoroptinhub, error) {
	contract, err := bindValidatoroptinhub(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validatoroptinhub{ValidatoroptinhubCaller: ValidatoroptinhubCaller{contract: contract}, ValidatoroptinhubTransactor: ValidatoroptinhubTransactor{contract: contract}, ValidatoroptinhubFilterer: ValidatoroptinhubFilterer{contract: contract}}, nil
}

// NewValidatoroptinhubCaller creates a new read-only instance of Validatoroptinhub, bound to a specific deployed contract.
func NewValidatoroptinhubCaller(address common.Address, caller bind.ContractCaller) (*ValidatoroptinhubCaller, error) {
	contract, err := bindValidatoroptinhub(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubCaller{contract: contract}, nil
}

// NewValidatoroptinhubTransactor creates a new write-only instance of Validatoroptinhub, bound to a specific deployed contract.
func NewValidatoroptinhubTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatoroptinhubTransactor, error) {
	contract, err := bindValidatoroptinhub(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubTransactor{contract: contract}, nil
}

// NewValidatoroptinhubFilterer creates a new log filterer instance of Validatoroptinhub, bound to a specific deployed contract.
func NewValidatoroptinhubFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatoroptinhubFilterer, error) {
	contract, err := bindValidatoroptinhub(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubFilterer{contract: contract}, nil
}

// bindValidatoroptinhub binds a generic wrapper to an already deployed contract.
func bindValidatoroptinhub(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidatoroptinhubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatoroptinhub *ValidatoroptinhubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatoroptinhub.Contract.ValidatoroptinhubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatoroptinhub *ValidatoroptinhubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.ValidatoroptinhubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatoroptinhub *ValidatoroptinhubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.ValidatoroptinhubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatoroptinhub *ValidatoroptinhubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatoroptinhub.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatoroptinhub *ValidatoroptinhubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatoroptinhub *ValidatoroptinhubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinhub *ValidatoroptinhubCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinhub *ValidatoroptinhubSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatoroptinhub.Contract.UPGRADEINTERFACEVERSION(&_Validatoroptinhub.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatoroptinhub.Contract.UPGRADEINTERFACEVERSION(&_Validatoroptinhub.CallOpts)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubCaller) AreValidatorsOptedIn(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([]bool, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "areValidatorsOptedIn", valBLSPubKeys)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubSession) AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error) {
	return _Validatoroptinhub.Contract.AreValidatorsOptedIn(&_Validatoroptinhub.CallOpts, valBLSPubKeys)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] valBLSPubKeys) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) AreValidatorsOptedIn(valBLSPubKeys [][]byte) ([]bool, error) {
	return _Validatoroptinhub.Contract.AreValidatorsOptedIn(&_Validatoroptinhub.CallOpts, valBLSPubKeys)
}

// AreValidatorsOptedInList is a free data retrieval call binding the contract method 0xc9384003.
//
// Solidity: function areValidatorsOptedInList(bytes[] valBLSPubKeys) view returns(bool[][])
func (_Validatoroptinhub *ValidatoroptinhubCaller) AreValidatorsOptedInList(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([][]bool, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "areValidatorsOptedInList", valBLSPubKeys)

	if err != nil {
		return *new([][]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([][]bool)).(*[][]bool)

	return out0, err

}

// AreValidatorsOptedInList is a free data retrieval call binding the contract method 0xc9384003.
//
// Solidity: function areValidatorsOptedInList(bytes[] valBLSPubKeys) view returns(bool[][])
func (_Validatoroptinhub *ValidatoroptinhubSession) AreValidatorsOptedInList(valBLSPubKeys [][]byte) ([][]bool, error) {
	return _Validatoroptinhub.Contract.AreValidatorsOptedInList(&_Validatoroptinhub.CallOpts, valBLSPubKeys)
}

// AreValidatorsOptedInList is a free data retrieval call binding the contract method 0xc9384003.
//
// Solidity: function areValidatorsOptedInList(bytes[] valBLSPubKeys) view returns(bool[][])
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) AreValidatorsOptedInList(valBLSPubKeys [][]byte) ([][]bool, error) {
	return _Validatoroptinhub.Contract.AreValidatorsOptedInList(&_Validatoroptinhub.CallOpts, valBLSPubKeys)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Validatoroptinhub *ValidatoroptinhubCaller) IsValidatorOptedIn(opts *bind.CallOpts, valPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "isValidatorOptedIn", valPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Validatoroptinhub *ValidatoroptinhubSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Validatoroptinhub.Contract.IsValidatorOptedIn(&_Validatoroptinhub.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Validatoroptinhub.Contract.IsValidatorOptedIn(&_Validatoroptinhub.CallOpts, valPubKey)
}

// IsValidatorOptedInList is a free data retrieval call binding the contract method 0x374876f2.
//
// Solidity: function isValidatorOptedInList(bytes valPubKey) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubCaller) IsValidatorOptedInList(opts *bind.CallOpts, valPubKey []byte) ([]bool, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "isValidatorOptedInList", valPubKey)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// IsValidatorOptedInList is a free data retrieval call binding the contract method 0x374876f2.
//
// Solidity: function isValidatorOptedInList(bytes valPubKey) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubSession) IsValidatorOptedInList(valPubKey []byte) ([]bool, error) {
	return _Validatoroptinhub.Contract.IsValidatorOptedInList(&_Validatoroptinhub.CallOpts, valPubKey)
}

// IsValidatorOptedInList is a free data retrieval call binding the contract method 0x374876f2.
//
// Solidity: function isValidatorOptedInList(bytes valPubKey) view returns(bool[])
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) IsValidatorOptedInList(valPubKey []byte) ([]bool, error) {
	return _Validatoroptinhub.Contract.IsValidatorOptedInList(&_Validatoroptinhub.CallOpts, valPubKey)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubSession) Owner() (common.Address, error) {
	return _Validatoroptinhub.Contract.Owner(&_Validatoroptinhub.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) Owner() (common.Address, error) {
	return _Validatoroptinhub.Contract.Owner(&_Validatoroptinhub.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubSession) PendingOwner() (common.Address, error) {
	return _Validatoroptinhub.Contract.PendingOwner(&_Validatoroptinhub.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) PendingOwner() (common.Address, error) {
	return _Validatoroptinhub.Contract.PendingOwner(&_Validatoroptinhub.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinhub *ValidatoroptinhubCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinhub *ValidatoroptinhubSession) ProxiableUUID() ([32]byte, error) {
	return _Validatoroptinhub.Contract.ProxiableUUID(&_Validatoroptinhub.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Validatoroptinhub.Contract.ProxiableUUID(&_Validatoroptinhub.CallOpts)
}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCaller) Registries(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Validatoroptinhub.contract.Call(opts, &out, "registries", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubSession) Registries(arg0 *big.Int) (common.Address, error) {
	return _Validatoroptinhub.Contract.Registries(&_Validatoroptinhub.CallOpts, arg0)
}

// Registries is a free data retrieval call binding the contract method 0x6347c900.
//
// Solidity: function registries(uint256 ) view returns(address)
func (_Validatoroptinhub *ValidatoroptinhubCallerSession) Registries(arg0 *big.Int) (common.Address, error) {
	return _Validatoroptinhub.Contract.Registries(&_Validatoroptinhub.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) AcceptOwnership() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.AcceptOwnership(&_Validatoroptinhub.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.AcceptOwnership(&_Validatoroptinhub.TransactOpts)
}

// AddRegistry is a paid mutator transaction binding the contract method 0xf6b11a23.
//
// Solidity: function addRegistry(address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) AddRegistry(opts *bind.TransactOpts, registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "addRegistry", registry)
}

// AddRegistry is a paid mutator transaction binding the contract method 0xf6b11a23.
//
// Solidity: function addRegistry(address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) AddRegistry(registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.AddRegistry(&_Validatoroptinhub.TransactOpts, registry)
}

// AddRegistry is a paid mutator transaction binding the contract method 0xf6b11a23.
//
// Solidity: function addRegistry(address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) AddRegistry(registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.AddRegistry(&_Validatoroptinhub.TransactOpts, registry)
}

// Initialize is a paid mutator transaction binding the contract method 0x462d0b2e.
//
// Solidity: function initialize(address[] _registries, address _owner) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) Initialize(opts *bind.TransactOpts, _registries []common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "initialize", _registries, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x462d0b2e.
//
// Solidity: function initialize(address[] _registries, address _owner) returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) Initialize(_registries []common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Initialize(&_Validatoroptinhub.TransactOpts, _registries, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x462d0b2e.
//
// Solidity: function initialize(address[] _registries, address _owner) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) Initialize(_registries []common.Address, _owner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Initialize(&_Validatoroptinhub.TransactOpts, _registries, _owner)
}

// RemoveRegistry is a paid mutator transaction binding the contract method 0x151c5f27.
//
// Solidity: function removeRegistry(uint256 index, address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) RemoveRegistry(opts *bind.TransactOpts, index *big.Int, registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "removeRegistry", index, registry)
}

// RemoveRegistry is a paid mutator transaction binding the contract method 0x151c5f27.
//
// Solidity: function removeRegistry(uint256 index, address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) RemoveRegistry(index *big.Int, registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.RemoveRegistry(&_Validatoroptinhub.TransactOpts, index, registry)
}

// RemoveRegistry is a paid mutator transaction binding the contract method 0x151c5f27.
//
// Solidity: function removeRegistry(uint256 index, address registry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) RemoveRegistry(index *big.Int, registry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.RemoveRegistry(&_Validatoroptinhub.TransactOpts, index, registry)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.RenounceOwnership(&_Validatoroptinhub.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.RenounceOwnership(&_Validatoroptinhub.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.TransferOwnership(&_Validatoroptinhub.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.TransferOwnership(&_Validatoroptinhub.TransactOpts, newOwner)
}

// UpdateRegistry is a paid mutator transaction binding the contract method 0x28b4c446.
//
// Solidity: function updateRegistry(uint256 index, address oldRegistry, address newRegistry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) UpdateRegistry(opts *bind.TransactOpts, index *big.Int, oldRegistry common.Address, newRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "updateRegistry", index, oldRegistry, newRegistry)
}

// UpdateRegistry is a paid mutator transaction binding the contract method 0x28b4c446.
//
// Solidity: function updateRegistry(uint256 index, address oldRegistry, address newRegistry) returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) UpdateRegistry(index *big.Int, oldRegistry common.Address, newRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.UpdateRegistry(&_Validatoroptinhub.TransactOpts, index, oldRegistry, newRegistry)
}

// UpdateRegistry is a paid mutator transaction binding the contract method 0x28b4c446.
//
// Solidity: function updateRegistry(uint256 index, address oldRegistry, address newRegistry) returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) UpdateRegistry(index *big.Int, oldRegistry common.Address, newRegistry common.Address) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.UpdateRegistry(&_Validatoroptinhub.TransactOpts, index, oldRegistry, newRegistry)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.UpgradeToAndCall(&_Validatoroptinhub.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.UpgradeToAndCall(&_Validatoroptinhub.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Fallback(&_Validatoroptinhub.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Fallback(&_Validatoroptinhub.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatoroptinhub.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubSession) Receive() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Receive(&_Validatoroptinhub.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatoroptinhub *ValidatoroptinhubTransactorSession) Receive() (*types.Transaction, error) {
	return _Validatoroptinhub.Contract.Receive(&_Validatoroptinhub.TransactOpts)
}

// ValidatoroptinhubInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validatoroptinhub contract.
type ValidatoroptinhubInitializedIterator struct {
	Event *ValidatoroptinhubInitialized // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubInitialized)
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
		it.Event = new(ValidatoroptinhubInitialized)
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
func (it *ValidatoroptinhubInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubInitialized represents a Initialized event raised by the Validatoroptinhub contract.
type ValidatoroptinhubInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterInitialized(opts *bind.FilterOpts) (*ValidatoroptinhubInitializedIterator, error) {

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubInitializedIterator{contract: _Validatoroptinhub.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubInitialized) (event.Subscription, error) {

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubInitialized)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseInitialized(log types.Log) (*ValidatoroptinhubInitialized, error) {
	event := new(ValidatoroptinhubInitialized)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Validatoroptinhub contract.
type ValidatoroptinhubOwnershipTransferStartedIterator struct {
	Event *ValidatoroptinhubOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubOwnershipTransferStarted)
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
		it.Event = new(ValidatoroptinhubOwnershipTransferStarted)
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
func (it *ValidatoroptinhubOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Validatoroptinhub contract.
type ValidatoroptinhubOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatoroptinhubOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubOwnershipTransferStartedIterator{contract: _Validatoroptinhub.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubOwnershipTransferStarted)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseOwnershipTransferStarted(log types.Log) (*ValidatoroptinhubOwnershipTransferStarted, error) {
	event := new(ValidatoroptinhubOwnershipTransferStarted)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validatoroptinhub contract.
type ValidatoroptinhubOwnershipTransferredIterator struct {
	Event *ValidatoroptinhubOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubOwnershipTransferred)
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
		it.Event = new(ValidatoroptinhubOwnershipTransferred)
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
func (it *ValidatoroptinhubOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubOwnershipTransferred represents a OwnershipTransferred event raised by the Validatoroptinhub contract.
type ValidatoroptinhubOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatoroptinhubOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubOwnershipTransferredIterator{contract: _Validatoroptinhub.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubOwnershipTransferred)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseOwnershipTransferred(log types.Log) (*ValidatoroptinhubOwnershipTransferred, error) {
	event := new(ValidatoroptinhubOwnershipTransferred)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubRegistryAddedIterator is returned from FilterRegistryAdded and is used to iterate over the raw logs and unpacked data for RegistryAdded events raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryAddedIterator struct {
	Event *ValidatoroptinhubRegistryAdded // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubRegistryAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubRegistryAdded)
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
		it.Event = new(ValidatoroptinhubRegistryAdded)
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
func (it *ValidatoroptinhubRegistryAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubRegistryAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubRegistryAdded represents a RegistryAdded event raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryAdded struct {
	Index    *big.Int
	Registry common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRegistryAdded is a free log retrieval operation binding the contract event 0x08d8a372b87adbf98c2b19281866b69031c8e54ae5fa823510d19b803b11382a.
//
// Solidity: event RegistryAdded(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterRegistryAdded(opts *bind.FilterOpts, index []*big.Int, registry []common.Address) (*ValidatoroptinhubRegistryAddedIterator, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var registryRule []interface{}
	for _, registryItem := range registry {
		registryRule = append(registryRule, registryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "RegistryAdded", indexRule, registryRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubRegistryAddedIterator{contract: _Validatoroptinhub.contract, event: "RegistryAdded", logs: logs, sub: sub}, nil
}

// WatchRegistryAdded is a free log subscription operation binding the contract event 0x08d8a372b87adbf98c2b19281866b69031c8e54ae5fa823510d19b803b11382a.
//
// Solidity: event RegistryAdded(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchRegistryAdded(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubRegistryAdded, index []*big.Int, registry []common.Address) (event.Subscription, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var registryRule []interface{}
	for _, registryItem := range registry {
		registryRule = append(registryRule, registryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "RegistryAdded", indexRule, registryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubRegistryAdded)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryAdded", log); err != nil {
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

// ParseRegistryAdded is a log parse operation binding the contract event 0x08d8a372b87adbf98c2b19281866b69031c8e54ae5fa823510d19b803b11382a.
//
// Solidity: event RegistryAdded(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseRegistryAdded(log types.Log) (*ValidatoroptinhubRegistryAdded, error) {
	event := new(ValidatoroptinhubRegistryAdded)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubRegistryRemovedIterator is returned from FilterRegistryRemoved and is used to iterate over the raw logs and unpacked data for RegistryRemoved events raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryRemovedIterator struct {
	Event *ValidatoroptinhubRegistryRemoved // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubRegistryRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubRegistryRemoved)
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
		it.Event = new(ValidatoroptinhubRegistryRemoved)
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
func (it *ValidatoroptinhubRegistryRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubRegistryRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubRegistryRemoved represents a RegistryRemoved event raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryRemoved struct {
	Index    *big.Int
	Registry common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRegistryRemoved is a free log retrieval operation binding the contract event 0xabeeaf793e0f2bff8bb9ceb901ea839147a1c0282fb908f6f18142bc05ffcd1d.
//
// Solidity: event RegistryRemoved(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterRegistryRemoved(opts *bind.FilterOpts, index []*big.Int, registry []common.Address) (*ValidatoroptinhubRegistryRemovedIterator, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var registryRule []interface{}
	for _, registryItem := range registry {
		registryRule = append(registryRule, registryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "RegistryRemoved", indexRule, registryRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubRegistryRemovedIterator{contract: _Validatoroptinhub.contract, event: "RegistryRemoved", logs: logs, sub: sub}, nil
}

// WatchRegistryRemoved is a free log subscription operation binding the contract event 0xabeeaf793e0f2bff8bb9ceb901ea839147a1c0282fb908f6f18142bc05ffcd1d.
//
// Solidity: event RegistryRemoved(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchRegistryRemoved(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubRegistryRemoved, index []*big.Int, registry []common.Address) (event.Subscription, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var registryRule []interface{}
	for _, registryItem := range registry {
		registryRule = append(registryRule, registryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "RegistryRemoved", indexRule, registryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubRegistryRemoved)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryRemoved", log); err != nil {
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

// ParseRegistryRemoved is a log parse operation binding the contract event 0xabeeaf793e0f2bff8bb9ceb901ea839147a1c0282fb908f6f18142bc05ffcd1d.
//
// Solidity: event RegistryRemoved(uint256 indexed index, address indexed registry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseRegistryRemoved(log types.Log) (*ValidatoroptinhubRegistryRemoved, error) {
	event := new(ValidatoroptinhubRegistryRemoved)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubRegistryReplacedIterator is returned from FilterRegistryReplaced and is used to iterate over the raw logs and unpacked data for RegistryReplaced events raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryReplacedIterator struct {
	Event *ValidatoroptinhubRegistryReplaced // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubRegistryReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubRegistryReplaced)
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
		it.Event = new(ValidatoroptinhubRegistryReplaced)
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
func (it *ValidatoroptinhubRegistryReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubRegistryReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubRegistryReplaced represents a RegistryReplaced event raised by the Validatoroptinhub contract.
type ValidatoroptinhubRegistryReplaced struct {
	Index       *big.Int
	OldRegistry common.Address
	NewRegistry common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRegistryReplaced is a free log retrieval operation binding the contract event 0x7d7b26f7fb2a6f76146f4e2ae8c5168da44eab44cd4fbecd3dc72ad3bafe03fe.
//
// Solidity: event RegistryReplaced(uint256 indexed index, address indexed oldRegistry, address indexed newRegistry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterRegistryReplaced(opts *bind.FilterOpts, index []*big.Int, oldRegistry []common.Address, newRegistry []common.Address) (*ValidatoroptinhubRegistryReplacedIterator, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "RegistryReplaced", indexRule, oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubRegistryReplacedIterator{contract: _Validatoroptinhub.contract, event: "RegistryReplaced", logs: logs, sub: sub}, nil
}

// WatchRegistryReplaced is a free log subscription operation binding the contract event 0x7d7b26f7fb2a6f76146f4e2ae8c5168da44eab44cd4fbecd3dc72ad3bafe03fe.
//
// Solidity: event RegistryReplaced(uint256 indexed index, address indexed oldRegistry, address indexed newRegistry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchRegistryReplaced(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubRegistryReplaced, index []*big.Int, oldRegistry []common.Address, newRegistry []common.Address) (event.Subscription, error) {

	var indexRule []interface{}
	for _, indexItem := range index {
		indexRule = append(indexRule, indexItem)
	}
	var oldRegistryRule []interface{}
	for _, oldRegistryItem := range oldRegistry {
		oldRegistryRule = append(oldRegistryRule, oldRegistryItem)
	}
	var newRegistryRule []interface{}
	for _, newRegistryItem := range newRegistry {
		newRegistryRule = append(newRegistryRule, newRegistryItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "RegistryReplaced", indexRule, oldRegistryRule, newRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubRegistryReplaced)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryReplaced", log); err != nil {
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

// ParseRegistryReplaced is a log parse operation binding the contract event 0x7d7b26f7fb2a6f76146f4e2ae8c5168da44eab44cd4fbecd3dc72ad3bafe03fe.
//
// Solidity: event RegistryReplaced(uint256 indexed index, address indexed oldRegistry, address indexed newRegistry)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseRegistryReplaced(log types.Log) (*ValidatoroptinhubRegistryReplaced, error) {
	event := new(ValidatoroptinhubRegistryReplaced)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "RegistryReplaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatoroptinhubUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Validatoroptinhub contract.
type ValidatoroptinhubUpgradedIterator struct {
	Event *ValidatoroptinhubUpgraded // Event containing the contract specifics and raw log

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
func (it *ValidatoroptinhubUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatoroptinhubUpgraded)
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
		it.Event = new(ValidatoroptinhubUpgraded)
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
func (it *ValidatoroptinhubUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatoroptinhubUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatoroptinhubUpgraded represents a Upgraded event raised by the Validatoroptinhub contract.
type ValidatoroptinhubUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ValidatoroptinhubUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ValidatoroptinhubUpgradedIterator{contract: _Validatoroptinhub.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatoroptinhub *ValidatoroptinhubFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ValidatoroptinhubUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatoroptinhub.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatoroptinhubUpgraded)
				if err := _Validatoroptinhub.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Validatoroptinhub *ValidatoroptinhubFilterer) ParseUpgraded(log types.Log) (*ValidatoroptinhubUpgraded, error) {
	event := new(ValidatoroptinhubUpgraded)
	if err := _Validatoroptinhub.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
