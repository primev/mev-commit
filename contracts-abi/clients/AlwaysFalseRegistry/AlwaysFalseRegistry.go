// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package alwaysfalseregistry

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

// AlwaysfalseregistryMetaData contains all meta data concerning the Alwaysfalseregistry contract.
var AlwaysfalseregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"}]",
}

// AlwaysfalseregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AlwaysfalseregistryMetaData.ABI instead.
var AlwaysfalseregistryABI = AlwaysfalseregistryMetaData.ABI

// Alwaysfalseregistry is an auto generated Go binding around an Ethereum contract.
type Alwaysfalseregistry struct {
	AlwaysfalseregistryCaller     // Read-only binding to the contract
	AlwaysfalseregistryTransactor // Write-only binding to the contract
	AlwaysfalseregistryFilterer   // Log filterer for contract events
}

// AlwaysfalseregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AlwaysfalseregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlwaysfalseregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AlwaysfalseregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlwaysfalseregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AlwaysfalseregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlwaysfalseregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AlwaysfalseregistrySession struct {
	Contract     *Alwaysfalseregistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AlwaysfalseregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AlwaysfalseregistryCallerSession struct {
	Contract *AlwaysfalseregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// AlwaysfalseregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AlwaysfalseregistryTransactorSession struct {
	Contract     *AlwaysfalseregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// AlwaysfalseregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AlwaysfalseregistryRaw struct {
	Contract *Alwaysfalseregistry // Generic contract binding to access the raw methods on
}

// AlwaysfalseregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AlwaysfalseregistryCallerRaw struct {
	Contract *AlwaysfalseregistryCaller // Generic read-only contract binding to access the raw methods on
}

// AlwaysfalseregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AlwaysfalseregistryTransactorRaw struct {
	Contract *AlwaysfalseregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAlwaysfalseregistry creates a new instance of Alwaysfalseregistry, bound to a specific deployed contract.
func NewAlwaysfalseregistry(address common.Address, backend bind.ContractBackend) (*Alwaysfalseregistry, error) {
	contract, err := bindAlwaysfalseregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Alwaysfalseregistry{AlwaysfalseregistryCaller: AlwaysfalseregistryCaller{contract: contract}, AlwaysfalseregistryTransactor: AlwaysfalseregistryTransactor{contract: contract}, AlwaysfalseregistryFilterer: AlwaysfalseregistryFilterer{contract: contract}}, nil
}

// NewAlwaysfalseregistryCaller creates a new read-only instance of Alwaysfalseregistry, bound to a specific deployed contract.
func NewAlwaysfalseregistryCaller(address common.Address, caller bind.ContractCaller) (*AlwaysfalseregistryCaller, error) {
	contract, err := bindAlwaysfalseregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AlwaysfalseregistryCaller{contract: contract}, nil
}

// NewAlwaysfalseregistryTransactor creates a new write-only instance of Alwaysfalseregistry, bound to a specific deployed contract.
func NewAlwaysfalseregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AlwaysfalseregistryTransactor, error) {
	contract, err := bindAlwaysfalseregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AlwaysfalseregistryTransactor{contract: contract}, nil
}

// NewAlwaysfalseregistryFilterer creates a new log filterer instance of Alwaysfalseregistry, bound to a specific deployed contract.
func NewAlwaysfalseregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AlwaysfalseregistryFilterer, error) {
	contract, err := bindAlwaysfalseregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AlwaysfalseregistryFilterer{contract: contract}, nil
}

// bindAlwaysfalseregistry binds a generic wrapper to an already deployed contract.
func bindAlwaysfalseregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AlwaysfalseregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Alwaysfalseregistry *AlwaysfalseregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Alwaysfalseregistry.Contract.AlwaysfalseregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Alwaysfalseregistry *AlwaysfalseregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Alwaysfalseregistry.Contract.AlwaysfalseregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Alwaysfalseregistry *AlwaysfalseregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Alwaysfalseregistry.Contract.AlwaysfalseregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Alwaysfalseregistry *AlwaysfalseregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Alwaysfalseregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Alwaysfalseregistry *AlwaysfalseregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Alwaysfalseregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Alwaysfalseregistry *AlwaysfalseregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Alwaysfalseregistry.Contract.contract.Transact(opts, method, params...)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes ) pure returns(bool)
func (_Alwaysfalseregistry *AlwaysfalseregistryCaller) IsValidatorOptedIn(opts *bind.CallOpts, arg0 []byte) (bool, error) {
	var out []interface{}
	err := _Alwaysfalseregistry.contract.Call(opts, &out, "isValidatorOptedIn", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes ) pure returns(bool)
func (_Alwaysfalseregistry *AlwaysfalseregistrySession) IsValidatorOptedIn(arg0 []byte) (bool, error) {
	return _Alwaysfalseregistry.Contract.IsValidatorOptedIn(&_Alwaysfalseregistry.CallOpts, arg0)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes ) pure returns(bool)
func (_Alwaysfalseregistry *AlwaysfalseregistryCallerSession) IsValidatorOptedIn(arg0 []byte) (bool, error) {
	return _Alwaysfalseregistry.Contract.IsValidatorOptedIn(&_Alwaysfalseregistry.CallOpts, arg0)
}
