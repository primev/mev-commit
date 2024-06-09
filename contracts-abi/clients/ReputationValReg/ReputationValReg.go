// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package reputationvalreg

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

// ReputationvalregMetaData contains all meta data concerning the Reputationvalreg contract.
var ReputationvalregMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addWhitelistedEOA\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"areValidatorsOptedIn\",\"inputs\":[{\"name\":\"consAddrs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deleteConsAddrs\",\"inputs\":[{\"name\":\"consAddrs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deleteWhitelistedEOA\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freeze\",\"inputs\":[{\"name\":\"validatorConsAddr\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getWhitelistedEOAInfo\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumReputationValReg.State\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_maxConsAddrsPerEOA\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minFreezeBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_unfreezeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isEOAWhitelisted\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxConsAddrsPerEOA\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minFreezeBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeConsAddrs\",\"inputs\":[{\"name\":\"consAddrs\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storedConsAddrs\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unfreezeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"whitelistedEOAs\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumReputationValReg.State\"},{\"name\":\"numConsAddrsStored\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"freezeHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"moniker\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"ConsAddrDeleted\",\"inputs\":[{\"name\":\"consAddr\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ConsAddrStored\",\"inputs\":[{\"name\":\"consAddr\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EOAFrozen\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EOAUnfrozen\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistedEOAAdded\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistedEOADeleted\",\"inputs\":[{\"name\":\"eoa\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"moniker\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// ReputationvalregABI is the input ABI used to generate the binding from.
// Deprecated: Use ReputationvalregMetaData.ABI instead.
var ReputationvalregABI = ReputationvalregMetaData.ABI

// Reputationvalreg is an auto generated Go binding around an Ethereum contract.
type Reputationvalreg struct {
	ReputationvalregCaller     // Read-only binding to the contract
	ReputationvalregTransactor // Write-only binding to the contract
	ReputationvalregFilterer   // Log filterer for contract events
}

// ReputationvalregCaller is an auto generated read-only Go binding around an Ethereum contract.
type ReputationvalregCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationvalregTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ReputationvalregTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationvalregFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ReputationvalregFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReputationvalregSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ReputationvalregSession struct {
	Contract     *Reputationvalreg // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ReputationvalregCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ReputationvalregCallerSession struct {
	Contract *ReputationvalregCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ReputationvalregTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ReputationvalregTransactorSession struct {
	Contract     *ReputationvalregTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ReputationvalregRaw is an auto generated low-level Go binding around an Ethereum contract.
type ReputationvalregRaw struct {
	Contract *Reputationvalreg // Generic contract binding to access the raw methods on
}

// ReputationvalregCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ReputationvalregCallerRaw struct {
	Contract *ReputationvalregCaller // Generic read-only contract binding to access the raw methods on
}

// ReputationvalregTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ReputationvalregTransactorRaw struct {
	Contract *ReputationvalregTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReputationvalreg creates a new instance of Reputationvalreg, bound to a specific deployed contract.
func NewReputationvalreg(address common.Address, backend bind.ContractBackend) (*Reputationvalreg, error) {
	contract, err := bindReputationvalreg(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reputationvalreg{ReputationvalregCaller: ReputationvalregCaller{contract: contract}, ReputationvalregTransactor: ReputationvalregTransactor{contract: contract}, ReputationvalregFilterer: ReputationvalregFilterer{contract: contract}}, nil
}

// NewReputationvalregCaller creates a new read-only instance of Reputationvalreg, bound to a specific deployed contract.
func NewReputationvalregCaller(address common.Address, caller bind.ContractCaller) (*ReputationvalregCaller, error) {
	contract, err := bindReputationvalreg(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregCaller{contract: contract}, nil
}

// NewReputationvalregTransactor creates a new write-only instance of Reputationvalreg, bound to a specific deployed contract.
func NewReputationvalregTransactor(address common.Address, transactor bind.ContractTransactor) (*ReputationvalregTransactor, error) {
	contract, err := bindReputationvalreg(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregTransactor{contract: contract}, nil
}

// NewReputationvalregFilterer creates a new log filterer instance of Reputationvalreg, bound to a specific deployed contract.
func NewReputationvalregFilterer(address common.Address, filterer bind.ContractFilterer) (*ReputationvalregFilterer, error) {
	contract, err := bindReputationvalreg(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregFilterer{contract: contract}, nil
}

// bindReputationvalreg binds a generic wrapper to an already deployed contract.
func bindReputationvalreg(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReputationvalregMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reputationvalreg *ReputationvalregRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reputationvalreg.Contract.ReputationvalregCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reputationvalreg *ReputationvalregRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.ReputationvalregTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reputationvalreg *ReputationvalregRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.ReputationvalregTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reputationvalreg *ReputationvalregCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reputationvalreg.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reputationvalreg *ReputationvalregTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reputationvalreg *ReputationvalregTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Reputationvalreg *ReputationvalregCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Reputationvalreg *ReputationvalregSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Reputationvalreg.Contract.UPGRADEINTERFACEVERSION(&_Reputationvalreg.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Reputationvalreg *ReputationvalregCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Reputationvalreg.Contract.UPGRADEINTERFACEVERSION(&_Reputationvalreg.CallOpts)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] consAddrs) view returns(bool[])
func (_Reputationvalreg *ReputationvalregCaller) AreValidatorsOptedIn(opts *bind.CallOpts, consAddrs [][]byte) ([]bool, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "areValidatorsOptedIn", consAddrs)

	if err != nil {
		return *new([]bool), err
	}

	out0 := *abi.ConvertType(out[0], new([]bool)).(*[]bool)

	return out0, err

}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] consAddrs) view returns(bool[])
func (_Reputationvalreg *ReputationvalregSession) AreValidatorsOptedIn(consAddrs [][]byte) ([]bool, error) {
	return _Reputationvalreg.Contract.AreValidatorsOptedIn(&_Reputationvalreg.CallOpts, consAddrs)
}

// AreValidatorsOptedIn is a free data retrieval call binding the contract method 0x5b88b173.
//
// Solidity: function areValidatorsOptedIn(bytes[] consAddrs) view returns(bool[])
func (_Reputationvalreg *ReputationvalregCallerSession) AreValidatorsOptedIn(consAddrs [][]byte) ([]bool, error) {
	return _Reputationvalreg.Contract.AreValidatorsOptedIn(&_Reputationvalreg.CallOpts, consAddrs)
}

// GetWhitelistedEOAInfo is a free data retrieval call binding the contract method 0xbfb2342b.
//
// Solidity: function getWhitelistedEOAInfo(address eoa) view returns(uint8, uint256, uint256, string)
func (_Reputationvalreg *ReputationvalregCaller) GetWhitelistedEOAInfo(opts *bind.CallOpts, eoa common.Address) (uint8, *big.Int, *big.Int, string, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "getWhitelistedEOAInfo", eoa)

	if err != nil {
		return *new(uint8), *new(*big.Int), *new(*big.Int), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)

	return out0, out1, out2, out3, err

}

// GetWhitelistedEOAInfo is a free data retrieval call binding the contract method 0xbfb2342b.
//
// Solidity: function getWhitelistedEOAInfo(address eoa) view returns(uint8, uint256, uint256, string)
func (_Reputationvalreg *ReputationvalregSession) GetWhitelistedEOAInfo(eoa common.Address) (uint8, *big.Int, *big.Int, string, error) {
	return _Reputationvalreg.Contract.GetWhitelistedEOAInfo(&_Reputationvalreg.CallOpts, eoa)
}

// GetWhitelistedEOAInfo is a free data retrieval call binding the contract method 0xbfb2342b.
//
// Solidity: function getWhitelistedEOAInfo(address eoa) view returns(uint8, uint256, uint256, string)
func (_Reputationvalreg *ReputationvalregCallerSession) GetWhitelistedEOAInfo(eoa common.Address) (uint8, *big.Int, *big.Int, string, error) {
	return _Reputationvalreg.Contract.GetWhitelistedEOAInfo(&_Reputationvalreg.CallOpts, eoa)
}

// IsEOAWhitelisted is a free data retrieval call binding the contract method 0x346abf9d.
//
// Solidity: function isEOAWhitelisted(address eoa) view returns(bool)
func (_Reputationvalreg *ReputationvalregCaller) IsEOAWhitelisted(opts *bind.CallOpts, eoa common.Address) (bool, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "isEOAWhitelisted", eoa)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEOAWhitelisted is a free data retrieval call binding the contract method 0x346abf9d.
//
// Solidity: function isEOAWhitelisted(address eoa) view returns(bool)
func (_Reputationvalreg *ReputationvalregSession) IsEOAWhitelisted(eoa common.Address) (bool, error) {
	return _Reputationvalreg.Contract.IsEOAWhitelisted(&_Reputationvalreg.CallOpts, eoa)
}

// IsEOAWhitelisted is a free data retrieval call binding the contract method 0x346abf9d.
//
// Solidity: function isEOAWhitelisted(address eoa) view returns(bool)
func (_Reputationvalreg *ReputationvalregCallerSession) IsEOAWhitelisted(eoa common.Address) (bool, error) {
	return _Reputationvalreg.Contract.IsEOAWhitelisted(&_Reputationvalreg.CallOpts, eoa)
}

// MaxConsAddrsPerEOA is a free data retrieval call binding the contract method 0x840d7436.
//
// Solidity: function maxConsAddrsPerEOA() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCaller) MaxConsAddrsPerEOA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "maxConsAddrsPerEOA")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxConsAddrsPerEOA is a free data retrieval call binding the contract method 0x840d7436.
//
// Solidity: function maxConsAddrsPerEOA() view returns(uint256)
func (_Reputationvalreg *ReputationvalregSession) MaxConsAddrsPerEOA() (*big.Int, error) {
	return _Reputationvalreg.Contract.MaxConsAddrsPerEOA(&_Reputationvalreg.CallOpts)
}

// MaxConsAddrsPerEOA is a free data retrieval call binding the contract method 0x840d7436.
//
// Solidity: function maxConsAddrsPerEOA() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCallerSession) MaxConsAddrsPerEOA() (*big.Int, error) {
	return _Reputationvalreg.Contract.MaxConsAddrsPerEOA(&_Reputationvalreg.CallOpts)
}

// MinFreezeBlocks is a free data retrieval call binding the contract method 0x2fe0326d.
//
// Solidity: function minFreezeBlocks() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCaller) MinFreezeBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "minFreezeBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinFreezeBlocks is a free data retrieval call binding the contract method 0x2fe0326d.
//
// Solidity: function minFreezeBlocks() view returns(uint256)
func (_Reputationvalreg *ReputationvalregSession) MinFreezeBlocks() (*big.Int, error) {
	return _Reputationvalreg.Contract.MinFreezeBlocks(&_Reputationvalreg.CallOpts)
}

// MinFreezeBlocks is a free data retrieval call binding the contract method 0x2fe0326d.
//
// Solidity: function minFreezeBlocks() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCallerSession) MinFreezeBlocks() (*big.Int, error) {
	return _Reputationvalreg.Contract.MinFreezeBlocks(&_Reputationvalreg.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Reputationvalreg *ReputationvalregCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Reputationvalreg *ReputationvalregSession) Owner() (common.Address, error) {
	return _Reputationvalreg.Contract.Owner(&_Reputationvalreg.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Reputationvalreg *ReputationvalregCallerSession) Owner() (common.Address, error) {
	return _Reputationvalreg.Contract.Owner(&_Reputationvalreg.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reputationvalreg *ReputationvalregCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reputationvalreg *ReputationvalregSession) ProxiableUUID() ([32]byte, error) {
	return _Reputationvalreg.Contract.ProxiableUUID(&_Reputationvalreg.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Reputationvalreg *ReputationvalregCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Reputationvalreg.Contract.ProxiableUUID(&_Reputationvalreg.CallOpts)
}

// StoredConsAddrs is a free data retrieval call binding the contract method 0x90932287.
//
// Solidity: function storedConsAddrs(bytes ) view returns(address)
func (_Reputationvalreg *ReputationvalregCaller) StoredConsAddrs(opts *bind.CallOpts, arg0 []byte) (common.Address, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "storedConsAddrs", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StoredConsAddrs is a free data retrieval call binding the contract method 0x90932287.
//
// Solidity: function storedConsAddrs(bytes ) view returns(address)
func (_Reputationvalreg *ReputationvalregSession) StoredConsAddrs(arg0 []byte) (common.Address, error) {
	return _Reputationvalreg.Contract.StoredConsAddrs(&_Reputationvalreg.CallOpts, arg0)
}

// StoredConsAddrs is a free data retrieval call binding the contract method 0x90932287.
//
// Solidity: function storedConsAddrs(bytes ) view returns(address)
func (_Reputationvalreg *ReputationvalregCallerSession) StoredConsAddrs(arg0 []byte) (common.Address, error) {
	return _Reputationvalreg.Contract.StoredConsAddrs(&_Reputationvalreg.CallOpts, arg0)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCaller) UnfreezeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "unfreezeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Reputationvalreg *ReputationvalregSession) UnfreezeFee() (*big.Int, error) {
	return _Reputationvalreg.Contract.UnfreezeFee(&_Reputationvalreg.CallOpts)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Reputationvalreg *ReputationvalregCallerSession) UnfreezeFee() (*big.Int, error) {
	return _Reputationvalreg.Contract.UnfreezeFee(&_Reputationvalreg.CallOpts)
}

// WhitelistedEOAs is a free data retrieval call binding the contract method 0xa9c8364e.
//
// Solidity: function whitelistedEOAs(address ) view returns(uint8 state, uint256 numConsAddrsStored, uint256 freezeHeight, string moniker)
func (_Reputationvalreg *ReputationvalregCaller) WhitelistedEOAs(opts *bind.CallOpts, arg0 common.Address) (struct {
	State              uint8
	NumConsAddrsStored *big.Int
	FreezeHeight       *big.Int
	Moniker            string
}, error) {
	var out []interface{}
	err := _Reputationvalreg.contract.Call(opts, &out, "whitelistedEOAs", arg0)

	outstruct := new(struct {
		State              uint8
		NumConsAddrsStored *big.Int
		FreezeHeight       *big.Int
		Moniker            string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.NumConsAddrsStored = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FreezeHeight = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Moniker = *abi.ConvertType(out[3], new(string)).(*string)

	return *outstruct, err

}

// WhitelistedEOAs is a free data retrieval call binding the contract method 0xa9c8364e.
//
// Solidity: function whitelistedEOAs(address ) view returns(uint8 state, uint256 numConsAddrsStored, uint256 freezeHeight, string moniker)
func (_Reputationvalreg *ReputationvalregSession) WhitelistedEOAs(arg0 common.Address) (struct {
	State              uint8
	NumConsAddrsStored *big.Int
	FreezeHeight       *big.Int
	Moniker            string
}, error) {
	return _Reputationvalreg.Contract.WhitelistedEOAs(&_Reputationvalreg.CallOpts, arg0)
}

// WhitelistedEOAs is a free data retrieval call binding the contract method 0xa9c8364e.
//
// Solidity: function whitelistedEOAs(address ) view returns(uint8 state, uint256 numConsAddrsStored, uint256 freezeHeight, string moniker)
func (_Reputationvalreg *ReputationvalregCallerSession) WhitelistedEOAs(arg0 common.Address) (struct {
	State              uint8
	NumConsAddrsStored *big.Int
	FreezeHeight       *big.Int
	Moniker            string
}, error) {
	return _Reputationvalreg.Contract.WhitelistedEOAs(&_Reputationvalreg.CallOpts, arg0)
}

// AddWhitelistedEOA is a paid mutator transaction binding the contract method 0x8236b7c8.
//
// Solidity: function addWhitelistedEOA(address eoa, string moniker) returns()
func (_Reputationvalreg *ReputationvalregTransactor) AddWhitelistedEOA(opts *bind.TransactOpts, eoa common.Address, moniker string) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "addWhitelistedEOA", eoa, moniker)
}

// AddWhitelistedEOA is a paid mutator transaction binding the contract method 0x8236b7c8.
//
// Solidity: function addWhitelistedEOA(address eoa, string moniker) returns()
func (_Reputationvalreg *ReputationvalregSession) AddWhitelistedEOA(eoa common.Address, moniker string) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.AddWhitelistedEOA(&_Reputationvalreg.TransactOpts, eoa, moniker)
}

// AddWhitelistedEOA is a paid mutator transaction binding the contract method 0x8236b7c8.
//
// Solidity: function addWhitelistedEOA(address eoa, string moniker) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) AddWhitelistedEOA(eoa common.Address, moniker string) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.AddWhitelistedEOA(&_Reputationvalreg.TransactOpts, eoa, moniker)
}

// DeleteConsAddrs is a paid mutator transaction binding the contract method 0x5166792f.
//
// Solidity: function deleteConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregTransactor) DeleteConsAddrs(opts *bind.TransactOpts, consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "deleteConsAddrs", consAddrs)
}

// DeleteConsAddrs is a paid mutator transaction binding the contract method 0x5166792f.
//
// Solidity: function deleteConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregSession) DeleteConsAddrs(consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.DeleteConsAddrs(&_Reputationvalreg.TransactOpts, consAddrs)
}

// DeleteConsAddrs is a paid mutator transaction binding the contract method 0x5166792f.
//
// Solidity: function deleteConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) DeleteConsAddrs(consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.DeleteConsAddrs(&_Reputationvalreg.TransactOpts, consAddrs)
}

// DeleteWhitelistedEOA is a paid mutator transaction binding the contract method 0x9f922d31.
//
// Solidity: function deleteWhitelistedEOA(address eoa) returns()
func (_Reputationvalreg *ReputationvalregTransactor) DeleteWhitelistedEOA(opts *bind.TransactOpts, eoa common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "deleteWhitelistedEOA", eoa)
}

// DeleteWhitelistedEOA is a paid mutator transaction binding the contract method 0x9f922d31.
//
// Solidity: function deleteWhitelistedEOA(address eoa) returns()
func (_Reputationvalreg *ReputationvalregSession) DeleteWhitelistedEOA(eoa common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.DeleteWhitelistedEOA(&_Reputationvalreg.TransactOpts, eoa)
}

// DeleteWhitelistedEOA is a paid mutator transaction binding the contract method 0x9f922d31.
//
// Solidity: function deleteWhitelistedEOA(address eoa) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) DeleteWhitelistedEOA(eoa common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.DeleteWhitelistedEOA(&_Reputationvalreg.TransactOpts, eoa)
}

// Freeze is a paid mutator transaction binding the contract method 0xbb11fb62.
//
// Solidity: function freeze(bytes validatorConsAddr) returns()
func (_Reputationvalreg *ReputationvalregTransactor) Freeze(opts *bind.TransactOpts, validatorConsAddr []byte) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "freeze", validatorConsAddr)
}

// Freeze is a paid mutator transaction binding the contract method 0xbb11fb62.
//
// Solidity: function freeze(bytes validatorConsAddr) returns()
func (_Reputationvalreg *ReputationvalregSession) Freeze(validatorConsAddr []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Freeze(&_Reputationvalreg.TransactOpts, validatorConsAddr)
}

// Freeze is a paid mutator transaction binding the contract method 0xbb11fb62.
//
// Solidity: function freeze(bytes validatorConsAddr) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) Freeze(validatorConsAddr []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Freeze(&_Reputationvalreg.TransactOpts, validatorConsAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _owner, uint256 _maxConsAddrsPerEOA, uint256 _minFreezeBlocks, uint256 _unfreezeFee) returns()
func (_Reputationvalreg *ReputationvalregTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _maxConsAddrsPerEOA *big.Int, _minFreezeBlocks *big.Int, _unfreezeFee *big.Int) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "initialize", _owner, _maxConsAddrsPerEOA, _minFreezeBlocks, _unfreezeFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _owner, uint256 _maxConsAddrsPerEOA, uint256 _minFreezeBlocks, uint256 _unfreezeFee) returns()
func (_Reputationvalreg *ReputationvalregSession) Initialize(_owner common.Address, _maxConsAddrsPerEOA *big.Int, _minFreezeBlocks *big.Int, _unfreezeFee *big.Int) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Initialize(&_Reputationvalreg.TransactOpts, _owner, _maxConsAddrsPerEOA, _minFreezeBlocks, _unfreezeFee)
}

// Initialize is a paid mutator transaction binding the contract method 0x4ec81af1.
//
// Solidity: function initialize(address _owner, uint256 _maxConsAddrsPerEOA, uint256 _minFreezeBlocks, uint256 _unfreezeFee) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) Initialize(_owner common.Address, _maxConsAddrsPerEOA *big.Int, _minFreezeBlocks *big.Int, _unfreezeFee *big.Int) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Initialize(&_Reputationvalreg.TransactOpts, _owner, _maxConsAddrsPerEOA, _minFreezeBlocks, _unfreezeFee)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Reputationvalreg *ReputationvalregTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Reputationvalreg *ReputationvalregSession) RenounceOwnership() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.RenounceOwnership(&_Reputationvalreg.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.RenounceOwnership(&_Reputationvalreg.TransactOpts)
}

// StoreConsAddrs is a paid mutator transaction binding the contract method 0x6dfd5744.
//
// Solidity: function storeConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregTransactor) StoreConsAddrs(opts *bind.TransactOpts, consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "storeConsAddrs", consAddrs)
}

// StoreConsAddrs is a paid mutator transaction binding the contract method 0x6dfd5744.
//
// Solidity: function storeConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregSession) StoreConsAddrs(consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.StoreConsAddrs(&_Reputationvalreg.TransactOpts, consAddrs)
}

// StoreConsAddrs is a paid mutator transaction binding the contract method 0x6dfd5744.
//
// Solidity: function storeConsAddrs(bytes[] consAddrs) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) StoreConsAddrs(consAddrs [][]byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.StoreConsAddrs(&_Reputationvalreg.TransactOpts, consAddrs)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Reputationvalreg *ReputationvalregTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Reputationvalreg *ReputationvalregSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.TransferOwnership(&_Reputationvalreg.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.TransferOwnership(&_Reputationvalreg.TransactOpts, newOwner)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() payable returns()
func (_Reputationvalreg *ReputationvalregTransactor) Unfreeze(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "unfreeze")
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() payable returns()
func (_Reputationvalreg *ReputationvalregSession) Unfreeze() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Unfreeze(&_Reputationvalreg.TransactOpts)
}

// Unfreeze is a paid mutator transaction binding the contract method 0x6a28f000.
//
// Solidity: function unfreeze() payable returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) Unfreeze() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Unfreeze(&_Reputationvalreg.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reputationvalreg *ReputationvalregTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reputationvalreg.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reputationvalreg *ReputationvalregSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.UpgradeToAndCall(&_Reputationvalreg.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.UpgradeToAndCall(&_Reputationvalreg.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Reputationvalreg *ReputationvalregTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Reputationvalreg.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Reputationvalreg *ReputationvalregSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Fallback(&_Reputationvalreg.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Fallback(&_Reputationvalreg.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reputationvalreg *ReputationvalregTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reputationvalreg.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reputationvalreg *ReputationvalregSession) Receive() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Receive(&_Reputationvalreg.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Reputationvalreg *ReputationvalregTransactorSession) Receive() (*types.Transaction, error) {
	return _Reputationvalreg.Contract.Receive(&_Reputationvalreg.TransactOpts)
}

// ReputationvalregConsAddrDeletedIterator is returned from FilterConsAddrDeleted and is used to iterate over the raw logs and unpacked data for ConsAddrDeleted events raised by the Reputationvalreg contract.
type ReputationvalregConsAddrDeletedIterator struct {
	Event *ReputationvalregConsAddrDeleted // Event containing the contract specifics and raw log

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
func (it *ReputationvalregConsAddrDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregConsAddrDeleted)
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
		it.Event = new(ReputationvalregConsAddrDeleted)
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
func (it *ReputationvalregConsAddrDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregConsAddrDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregConsAddrDeleted represents a ConsAddrDeleted event raised by the Reputationvalreg contract.
type ReputationvalregConsAddrDeleted struct {
	ConsAddr []byte
	Eoa      common.Address
	Moniker  string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterConsAddrDeleted is a free log retrieval operation binding the contract event 0x9e2ae08eac35d7186f0b1b5bd8e43db28585e614e009f969f0bcf2b95454c84c.
//
// Solidity: event ConsAddrDeleted(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterConsAddrDeleted(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregConsAddrDeletedIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "ConsAddrDeleted", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregConsAddrDeletedIterator{contract: _Reputationvalreg.contract, event: "ConsAddrDeleted", logs: logs, sub: sub}, nil
}

// WatchConsAddrDeleted is a free log subscription operation binding the contract event 0x9e2ae08eac35d7186f0b1b5bd8e43db28585e614e009f969f0bcf2b95454c84c.
//
// Solidity: event ConsAddrDeleted(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchConsAddrDeleted(opts *bind.WatchOpts, sink chan<- *ReputationvalregConsAddrDeleted, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "ConsAddrDeleted", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregConsAddrDeleted)
				if err := _Reputationvalreg.contract.UnpackLog(event, "ConsAddrDeleted", log); err != nil {
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

// ParseConsAddrDeleted is a log parse operation binding the contract event 0x9e2ae08eac35d7186f0b1b5bd8e43db28585e614e009f969f0bcf2b95454c84c.
//
// Solidity: event ConsAddrDeleted(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseConsAddrDeleted(log types.Log) (*ReputationvalregConsAddrDeleted, error) {
	event := new(ReputationvalregConsAddrDeleted)
	if err := _Reputationvalreg.contract.UnpackLog(event, "ConsAddrDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregConsAddrStoredIterator is returned from FilterConsAddrStored and is used to iterate over the raw logs and unpacked data for ConsAddrStored events raised by the Reputationvalreg contract.
type ReputationvalregConsAddrStoredIterator struct {
	Event *ReputationvalregConsAddrStored // Event containing the contract specifics and raw log

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
func (it *ReputationvalregConsAddrStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregConsAddrStored)
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
		it.Event = new(ReputationvalregConsAddrStored)
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
func (it *ReputationvalregConsAddrStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregConsAddrStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregConsAddrStored represents a ConsAddrStored event raised by the Reputationvalreg contract.
type ReputationvalregConsAddrStored struct {
	ConsAddr []byte
	Eoa      common.Address
	Moniker  string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterConsAddrStored is a free log retrieval operation binding the contract event 0x168f982668ea302bc64437de90db02ffceac7e3fc43aa5080fcad7301105cf28.
//
// Solidity: event ConsAddrStored(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterConsAddrStored(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregConsAddrStoredIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "ConsAddrStored", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregConsAddrStoredIterator{contract: _Reputationvalreg.contract, event: "ConsAddrStored", logs: logs, sub: sub}, nil
}

// WatchConsAddrStored is a free log subscription operation binding the contract event 0x168f982668ea302bc64437de90db02ffceac7e3fc43aa5080fcad7301105cf28.
//
// Solidity: event ConsAddrStored(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchConsAddrStored(opts *bind.WatchOpts, sink chan<- *ReputationvalregConsAddrStored, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "ConsAddrStored", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregConsAddrStored)
				if err := _Reputationvalreg.contract.UnpackLog(event, "ConsAddrStored", log); err != nil {
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

// ParseConsAddrStored is a log parse operation binding the contract event 0x168f982668ea302bc64437de90db02ffceac7e3fc43aa5080fcad7301105cf28.
//
// Solidity: event ConsAddrStored(bytes consAddr, address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseConsAddrStored(log types.Log) (*ReputationvalregConsAddrStored, error) {
	event := new(ReputationvalregConsAddrStored)
	if err := _Reputationvalreg.contract.UnpackLog(event, "ConsAddrStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregEOAFrozenIterator is returned from FilterEOAFrozen and is used to iterate over the raw logs and unpacked data for EOAFrozen events raised by the Reputationvalreg contract.
type ReputationvalregEOAFrozenIterator struct {
	Event *ReputationvalregEOAFrozen // Event containing the contract specifics and raw log

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
func (it *ReputationvalregEOAFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregEOAFrozen)
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
		it.Event = new(ReputationvalregEOAFrozen)
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
func (it *ReputationvalregEOAFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregEOAFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregEOAFrozen represents a EOAFrozen event raised by the Reputationvalreg contract.
type ReputationvalregEOAFrozen struct {
	Eoa     common.Address
	Moniker string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEOAFrozen is a free log retrieval operation binding the contract event 0x9bc205fad361922c7873545a056ca9b2f24afd8816127d12220e77ba36e8a9b0.
//
// Solidity: event EOAFrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterEOAFrozen(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregEOAFrozenIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "EOAFrozen", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregEOAFrozenIterator{contract: _Reputationvalreg.contract, event: "EOAFrozen", logs: logs, sub: sub}, nil
}

// WatchEOAFrozen is a free log subscription operation binding the contract event 0x9bc205fad361922c7873545a056ca9b2f24afd8816127d12220e77ba36e8a9b0.
//
// Solidity: event EOAFrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchEOAFrozen(opts *bind.WatchOpts, sink chan<- *ReputationvalregEOAFrozen, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "EOAFrozen", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregEOAFrozen)
				if err := _Reputationvalreg.contract.UnpackLog(event, "EOAFrozen", log); err != nil {
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

// ParseEOAFrozen is a log parse operation binding the contract event 0x9bc205fad361922c7873545a056ca9b2f24afd8816127d12220e77ba36e8a9b0.
//
// Solidity: event EOAFrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseEOAFrozen(log types.Log) (*ReputationvalregEOAFrozen, error) {
	event := new(ReputationvalregEOAFrozen)
	if err := _Reputationvalreg.contract.UnpackLog(event, "EOAFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregEOAUnfrozenIterator is returned from FilterEOAUnfrozen and is used to iterate over the raw logs and unpacked data for EOAUnfrozen events raised by the Reputationvalreg contract.
type ReputationvalregEOAUnfrozenIterator struct {
	Event *ReputationvalregEOAUnfrozen // Event containing the contract specifics and raw log

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
func (it *ReputationvalregEOAUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregEOAUnfrozen)
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
		it.Event = new(ReputationvalregEOAUnfrozen)
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
func (it *ReputationvalregEOAUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregEOAUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregEOAUnfrozen represents a EOAUnfrozen event raised by the Reputationvalreg contract.
type ReputationvalregEOAUnfrozen struct {
	Eoa     common.Address
	Moniker string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterEOAUnfrozen is a free log retrieval operation binding the contract event 0xa0502c0ee62c5d07179a742a9d932b5876fbbc0a8d3bcf2ecabff3d5ea24c13e.
//
// Solidity: event EOAUnfrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterEOAUnfrozen(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregEOAUnfrozenIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "EOAUnfrozen", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregEOAUnfrozenIterator{contract: _Reputationvalreg.contract, event: "EOAUnfrozen", logs: logs, sub: sub}, nil
}

// WatchEOAUnfrozen is a free log subscription operation binding the contract event 0xa0502c0ee62c5d07179a742a9d932b5876fbbc0a8d3bcf2ecabff3d5ea24c13e.
//
// Solidity: event EOAUnfrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchEOAUnfrozen(opts *bind.WatchOpts, sink chan<- *ReputationvalregEOAUnfrozen, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "EOAUnfrozen", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregEOAUnfrozen)
				if err := _Reputationvalreg.contract.UnpackLog(event, "EOAUnfrozen", log); err != nil {
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

// ParseEOAUnfrozen is a log parse operation binding the contract event 0xa0502c0ee62c5d07179a742a9d932b5876fbbc0a8d3bcf2ecabff3d5ea24c13e.
//
// Solidity: event EOAUnfrozen(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseEOAUnfrozen(log types.Log) (*ReputationvalregEOAUnfrozen, error) {
	event := new(ReputationvalregEOAUnfrozen)
	if err := _Reputationvalreg.contract.UnpackLog(event, "EOAUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Reputationvalreg contract.
type ReputationvalregInitializedIterator struct {
	Event *ReputationvalregInitialized // Event containing the contract specifics and raw log

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
func (it *ReputationvalregInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregInitialized)
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
		it.Event = new(ReputationvalregInitialized)
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
func (it *ReputationvalregInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregInitialized represents a Initialized event raised by the Reputationvalreg contract.
type ReputationvalregInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Reputationvalreg *ReputationvalregFilterer) FilterInitialized(opts *bind.FilterOpts) (*ReputationvalregInitializedIterator, error) {

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ReputationvalregInitializedIterator{contract: _Reputationvalreg.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Reputationvalreg *ReputationvalregFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ReputationvalregInitialized) (event.Subscription, error) {

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregInitialized)
				if err := _Reputationvalreg.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Reputationvalreg *ReputationvalregFilterer) ParseInitialized(log types.Log) (*ReputationvalregInitialized, error) {
	event := new(ReputationvalregInitialized)
	if err := _Reputationvalreg.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Reputationvalreg contract.
type ReputationvalregOwnershipTransferredIterator struct {
	Event *ReputationvalregOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ReputationvalregOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregOwnershipTransferred)
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
		it.Event = new(ReputationvalregOwnershipTransferred)
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
func (it *ReputationvalregOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregOwnershipTransferred represents a OwnershipTransferred event raised by the Reputationvalreg contract.
type ReputationvalregOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Reputationvalreg *ReputationvalregFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ReputationvalregOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregOwnershipTransferredIterator{contract: _Reputationvalreg.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Reputationvalreg *ReputationvalregFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ReputationvalregOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregOwnershipTransferred)
				if err := _Reputationvalreg.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Reputationvalreg *ReputationvalregFilterer) ParseOwnershipTransferred(log types.Log) (*ReputationvalregOwnershipTransferred, error) {
	event := new(ReputationvalregOwnershipTransferred)
	if err := _Reputationvalreg.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Reputationvalreg contract.
type ReputationvalregUpgradedIterator struct {
	Event *ReputationvalregUpgraded // Event containing the contract specifics and raw log

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
func (it *ReputationvalregUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregUpgraded)
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
		it.Event = new(ReputationvalregUpgraded)
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
func (it *ReputationvalregUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregUpgraded represents a Upgraded event raised by the Reputationvalreg contract.
type ReputationvalregUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Reputationvalreg *ReputationvalregFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ReputationvalregUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregUpgradedIterator{contract: _Reputationvalreg.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Reputationvalreg *ReputationvalregFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ReputationvalregUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregUpgraded)
				if err := _Reputationvalreg.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Reputationvalreg *ReputationvalregFilterer) ParseUpgraded(log types.Log) (*ReputationvalregUpgraded, error) {
	event := new(ReputationvalregUpgraded)
	if err := _Reputationvalreg.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregWhitelistedEOAAddedIterator is returned from FilterWhitelistedEOAAdded and is used to iterate over the raw logs and unpacked data for WhitelistedEOAAdded events raised by the Reputationvalreg contract.
type ReputationvalregWhitelistedEOAAddedIterator struct {
	Event *ReputationvalregWhitelistedEOAAdded // Event containing the contract specifics and raw log

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
func (it *ReputationvalregWhitelistedEOAAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregWhitelistedEOAAdded)
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
		it.Event = new(ReputationvalregWhitelistedEOAAdded)
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
func (it *ReputationvalregWhitelistedEOAAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregWhitelistedEOAAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregWhitelistedEOAAdded represents a WhitelistedEOAAdded event raised by the Reputationvalreg contract.
type ReputationvalregWhitelistedEOAAdded struct {
	Eoa     common.Address
	Moniker string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistedEOAAdded is a free log retrieval operation binding the contract event 0x63989ede7c05f54429c0b9389838952e6c8a69538f925e4241b0875b9d50418e.
//
// Solidity: event WhitelistedEOAAdded(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterWhitelistedEOAAdded(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregWhitelistedEOAAddedIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "WhitelistedEOAAdded", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregWhitelistedEOAAddedIterator{contract: _Reputationvalreg.contract, event: "WhitelistedEOAAdded", logs: logs, sub: sub}, nil
}

// WatchWhitelistedEOAAdded is a free log subscription operation binding the contract event 0x63989ede7c05f54429c0b9389838952e6c8a69538f925e4241b0875b9d50418e.
//
// Solidity: event WhitelistedEOAAdded(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchWhitelistedEOAAdded(opts *bind.WatchOpts, sink chan<- *ReputationvalregWhitelistedEOAAdded, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "WhitelistedEOAAdded", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregWhitelistedEOAAdded)
				if err := _Reputationvalreg.contract.UnpackLog(event, "WhitelistedEOAAdded", log); err != nil {
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

// ParseWhitelistedEOAAdded is a log parse operation binding the contract event 0x63989ede7c05f54429c0b9389838952e6c8a69538f925e4241b0875b9d50418e.
//
// Solidity: event WhitelistedEOAAdded(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseWhitelistedEOAAdded(log types.Log) (*ReputationvalregWhitelistedEOAAdded, error) {
	event := new(ReputationvalregWhitelistedEOAAdded)
	if err := _Reputationvalreg.contract.UnpackLog(event, "WhitelistedEOAAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReputationvalregWhitelistedEOADeletedIterator is returned from FilterWhitelistedEOADeleted and is used to iterate over the raw logs and unpacked data for WhitelistedEOADeleted events raised by the Reputationvalreg contract.
type ReputationvalregWhitelistedEOADeletedIterator struct {
	Event *ReputationvalregWhitelistedEOADeleted // Event containing the contract specifics and raw log

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
func (it *ReputationvalregWhitelistedEOADeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ReputationvalregWhitelistedEOADeleted)
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
		it.Event = new(ReputationvalregWhitelistedEOADeleted)
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
func (it *ReputationvalregWhitelistedEOADeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ReputationvalregWhitelistedEOADeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ReputationvalregWhitelistedEOADeleted represents a WhitelistedEOADeleted event raised by the Reputationvalreg contract.
type ReputationvalregWhitelistedEOADeleted struct {
	Eoa     common.Address
	Moniker string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistedEOADeleted is a free log retrieval operation binding the contract event 0x9546eca2754003d7e04e3343b048988ba08856f4ce2ec3817725ecfb4a7d7d5c.
//
// Solidity: event WhitelistedEOADeleted(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) FilterWhitelistedEOADeleted(opts *bind.FilterOpts, eoa []common.Address) (*ReputationvalregWhitelistedEOADeletedIterator, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.FilterLogs(opts, "WhitelistedEOADeleted", eoaRule)
	if err != nil {
		return nil, err
	}
	return &ReputationvalregWhitelistedEOADeletedIterator{contract: _Reputationvalreg.contract, event: "WhitelistedEOADeleted", logs: logs, sub: sub}, nil
}

// WatchWhitelistedEOADeleted is a free log subscription operation binding the contract event 0x9546eca2754003d7e04e3343b048988ba08856f4ce2ec3817725ecfb4a7d7d5c.
//
// Solidity: event WhitelistedEOADeleted(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) WatchWhitelistedEOADeleted(opts *bind.WatchOpts, sink chan<- *ReputationvalregWhitelistedEOADeleted, eoa []common.Address) (event.Subscription, error) {

	var eoaRule []interface{}
	for _, eoaItem := range eoa {
		eoaRule = append(eoaRule, eoaItem)
	}

	logs, sub, err := _Reputationvalreg.contract.WatchLogs(opts, "WhitelistedEOADeleted", eoaRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ReputationvalregWhitelistedEOADeleted)
				if err := _Reputationvalreg.contract.UnpackLog(event, "WhitelistedEOADeleted", log); err != nil {
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

// ParseWhitelistedEOADeleted is a log parse operation binding the contract event 0x9546eca2754003d7e04e3343b048988ba08856f4ce2ec3817725ecfb4a7d7d5c.
//
// Solidity: event WhitelistedEOADeleted(address indexed eoa, string moniker)
func (_Reputationvalreg *ReputationvalregFilterer) ParseWhitelistedEOADeleted(log types.Log) (*ReputationvalregWhitelistedEOADeleted, error) {
	event := new(ReputationvalregWhitelistedEOADeleted)
	if err := _Reputationvalreg.contract.UnpackLog(event, "WhitelistedEOADeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
