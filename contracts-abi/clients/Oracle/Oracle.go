// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracle

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

// OracleMetaData contains all meta data concerning the Oracle contract.
var OracleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockBuilderNameToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"preconfManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blockTrackerContract_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"providerRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleAccount_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"oracleAccount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processBuilderCommitmentForBlockNumber\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"builder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBlockTracker\",\"inputs\":[{\"name\":\"newBlockTracker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOracleAccount\",\"inputs\":[{\"name\":\"newOracleAccount\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPreconfManager\",\"inputs\":[{\"name\":\"newPreconfManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"BlockTrackerSet\",\"inputs\":[{\"name\":\"newBlockTracker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitmentProcessed\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleAccountSet\",\"inputs\":[{\"name\":\"oldOracleAccount\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOracleAccount\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PreconfManagerSet\",\"inputs\":[{\"name\":\"newPreconfManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderRegistrySet\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BuilderNotBlockWinner\",\"inputs\":[{\"name\":\"blockWinner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"builder\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOracleAccount\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleAccount\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ResidualBidPercentAfterDecayExceeds100\",\"inputs\":[{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// OracleABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleMetaData.ABI instead.
var OracleABI = OracleMetaData.ABI

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Oracle *OracleCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Oracle *OracleSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Oracle.Contract.UPGRADEINTERFACEVERSION(&_Oracle.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Oracle *OracleCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Oracle.Contract.UPGRADEINTERFACEVERSION(&_Oracle.CallOpts)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Oracle *OracleCaller) BlockBuilderNameToAddress(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "blockBuilderNameToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Oracle *OracleSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Oracle.Contract.BlockBuilderNameToAddress(&_Oracle.CallOpts, arg0)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Oracle *OracleCallerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Oracle.Contract.BlockBuilderNameToAddress(&_Oracle.CallOpts, arg0)
}

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Oracle *OracleCaller) OracleAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "oracleAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Oracle *OracleSession) OracleAccount() (common.Address, error) {
	return _Oracle.Contract.OracleAccount(&_Oracle.CallOpts)
}

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Oracle *OracleCallerSession) OracleAccount() (common.Address, error) {
	return _Oracle.Contract.OracleAccount(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracle *OracleCallerSession) Owner() (common.Address, error) {
	return _Oracle.Contract.Owner(&_Oracle.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracle *OracleCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracle *OracleSession) Paused() (bool, error) {
	return _Oracle.Contract.Paused(&_Oracle.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracle *OracleCallerSession) Paused() (bool, error) {
	return _Oracle.Contract.Paused(&_Oracle.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracle *OracleCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracle *OracleSession) PendingOwner() (common.Address, error) {
	return _Oracle.Contract.PendingOwner(&_Oracle.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracle *OracleCallerSession) PendingOwner() (common.Address, error) {
	return _Oracle.Contract.PendingOwner(&_Oracle.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Oracle *OracleCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Oracle *OracleSession) ProxiableUUID() ([32]byte, error) {
	return _Oracle.Contract.ProxiableUUID(&_Oracle.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Oracle *OracleCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Oracle.Contract.ProxiableUUID(&_Oracle.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracle *OracleTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracle *OracleSession) AcceptOwnership() (*types.Transaction, error) {
	return _Oracle.Contract.AcceptOwnership(&_Oracle.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracle *OracleTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Oracle.Contract.AcceptOwnership(&_Oracle.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address preconfManager_, address blockTrackerContract_, address providerRegistry_, address oracleAccount_, address owner_) returns()
func (_Oracle *OracleTransactor) Initialize(opts *bind.TransactOpts, preconfManager_ common.Address, blockTrackerContract_ common.Address, providerRegistry_ common.Address, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "initialize", preconfManager_, blockTrackerContract_, providerRegistry_, oracleAccount_, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address preconfManager_, address blockTrackerContract_, address providerRegistry_, address oracleAccount_, address owner_) returns()
func (_Oracle *OracleSession) Initialize(preconfManager_ common.Address, blockTrackerContract_ common.Address, providerRegistry_ common.Address, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.Initialize(&_Oracle.TransactOpts, preconfManager_, blockTrackerContract_, providerRegistry_, oracleAccount_, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0x1459457a.
//
// Solidity: function initialize(address preconfManager_, address blockTrackerContract_, address providerRegistry_, address oracleAccount_, address owner_) returns()
func (_Oracle *OracleTransactorSession) Initialize(preconfManager_ common.Address, blockTrackerContract_ common.Address, providerRegistry_ common.Address, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.Initialize(&_Oracle.TransactOpts, preconfManager_, blockTrackerContract_, providerRegistry_, oracleAccount_, owner_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracle *OracleTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracle *OracleSession) Pause() (*types.Transaction, error) {
	return _Oracle.Contract.Pause(&_Oracle.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracle *OracleTransactorSession) Pause() (*types.Transaction, error) {
	return _Oracle.Contract.Pause(&_Oracle.TransactOpts)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x40d46772.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, address builder, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleTransactor) ProcessBuilderCommitmentForBlockNumber(opts *bind.TransactOpts, commitmentIndex [32]byte, blockNumber *big.Int, builder common.Address, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "processBuilderCommitmentForBlockNumber", commitmentIndex, blockNumber, builder, isSlash, residualBidPercentAfterDecay)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x40d46772.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, address builder, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, builder common.Address, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.ProcessBuilderCommitmentForBlockNumber(&_Oracle.TransactOpts, commitmentIndex, blockNumber, builder, isSlash, residualBidPercentAfterDecay)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x40d46772.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, address builder, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleTransactorSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, builder common.Address, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.ProcessBuilderCommitmentForBlockNumber(&_Oracle.TransactOpts, commitmentIndex, blockNumber, builder, isSlash, residualBidPercentAfterDecay)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Oracle *OracleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Oracle *OracleSession) RenounceOwnership() (*types.Transaction, error) {
	return _Oracle.Contract.RenounceOwnership(&_Oracle.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Oracle *OracleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Oracle.Contract.RenounceOwnership(&_Oracle.TransactOpts)
}

// SetBlockTracker is a paid mutator transaction binding the contract method 0x8b0ebeb9.
//
// Solidity: function setBlockTracker(address newBlockTracker) returns()
func (_Oracle *OracleTransactor) SetBlockTracker(opts *bind.TransactOpts, newBlockTracker common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setBlockTracker", newBlockTracker)
}

// SetBlockTracker is a paid mutator transaction binding the contract method 0x8b0ebeb9.
//
// Solidity: function setBlockTracker(address newBlockTracker) returns()
func (_Oracle *OracleSession) SetBlockTracker(newBlockTracker common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetBlockTracker(&_Oracle.TransactOpts, newBlockTracker)
}

// SetBlockTracker is a paid mutator transaction binding the contract method 0x8b0ebeb9.
//
// Solidity: function setBlockTracker(address newBlockTracker) returns()
func (_Oracle *OracleTransactorSession) SetBlockTracker(newBlockTracker common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetBlockTracker(&_Oracle.TransactOpts, newBlockTracker)
}

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Oracle *OracleTransactor) SetOracleAccount(opts *bind.TransactOpts, newOracleAccount common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setOracleAccount", newOracleAccount)
}

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Oracle *OracleSession) SetOracleAccount(newOracleAccount common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOracleAccount(&_Oracle.TransactOpts, newOracleAccount)
}

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Oracle *OracleTransactorSession) SetOracleAccount(newOracleAccount common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetOracleAccount(&_Oracle.TransactOpts, newOracleAccount)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address newPreconfManager) returns()
func (_Oracle *OracleTransactor) SetPreconfManager(opts *bind.TransactOpts, newPreconfManager common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setPreconfManager", newPreconfManager)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address newPreconfManager) returns()
func (_Oracle *OracleSession) SetPreconfManager(newPreconfManager common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetPreconfManager(&_Oracle.TransactOpts, newPreconfManager)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address newPreconfManager) returns()
func (_Oracle *OracleTransactorSession) SetPreconfManager(newPreconfManager common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.SetPreconfManager(&_Oracle.TransactOpts, newPreconfManager)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracle *OracleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.TransferOwnership(&_Oracle.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracle *OracleTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracle *OracleSession) Unpause() (*types.Transaction, error) {
	return _Oracle.Contract.Unpause(&_Oracle.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracle *OracleTransactorSession) Unpause() (*types.Transaction, error) {
	return _Oracle.Contract.Unpause(&_Oracle.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Oracle *OracleTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Oracle *OracleSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.UpgradeToAndCall(&_Oracle.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Oracle *OracleTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Oracle.Contract.UpgradeToAndCall(&_Oracle.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Oracle *OracleTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Oracle.Contract.Fallback(&_Oracle.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Oracle *OracleTransactorSession) Receive() (*types.Transaction, error) {
	return _Oracle.Contract.Receive(&_Oracle.TransactOpts)
}

// OracleBlockTrackerSetIterator is returned from FilterBlockTrackerSet and is used to iterate over the raw logs and unpacked data for BlockTrackerSet events raised by the Oracle contract.
type OracleBlockTrackerSetIterator struct {
	Event *OracleBlockTrackerSet // Event containing the contract specifics and raw log

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
func (it *OracleBlockTrackerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleBlockTrackerSet)
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
		it.Event = new(OracleBlockTrackerSet)
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
func (it *OracleBlockTrackerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleBlockTrackerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleBlockTrackerSet represents a BlockTrackerSet event raised by the Oracle contract.
type OracleBlockTrackerSet struct {
	NewBlockTracker common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBlockTrackerSet is a free log retrieval operation binding the contract event 0x99960ca3579bdd85d13cec5a46a9f349e2d4d9079d063608f362d736a573c311.
//
// Solidity: event BlockTrackerSet(address indexed newBlockTracker)
func (_Oracle *OracleFilterer) FilterBlockTrackerSet(opts *bind.FilterOpts, newBlockTracker []common.Address) (*OracleBlockTrackerSetIterator, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "BlockTrackerSet", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return &OracleBlockTrackerSetIterator{contract: _Oracle.contract, event: "BlockTrackerSet", logs: logs, sub: sub}, nil
}

// WatchBlockTrackerSet is a free log subscription operation binding the contract event 0x99960ca3579bdd85d13cec5a46a9f349e2d4d9079d063608f362d736a573c311.
//
// Solidity: event BlockTrackerSet(address indexed newBlockTracker)
func (_Oracle *OracleFilterer) WatchBlockTrackerSet(opts *bind.WatchOpts, sink chan<- *OracleBlockTrackerSet, newBlockTracker []common.Address) (event.Subscription, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "BlockTrackerSet", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleBlockTrackerSet)
				if err := _Oracle.contract.UnpackLog(event, "BlockTrackerSet", log); err != nil {
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

// ParseBlockTrackerSet is a log parse operation binding the contract event 0x99960ca3579bdd85d13cec5a46a9f349e2d4d9079d063608f362d736a573c311.
//
// Solidity: event BlockTrackerSet(address indexed newBlockTracker)
func (_Oracle *OracleFilterer) ParseBlockTrackerSet(log types.Log) (*OracleBlockTrackerSet, error) {
	event := new(OracleBlockTrackerSet)
	if err := _Oracle.contract.UnpackLog(event, "BlockTrackerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleCommitmentProcessedIterator is returned from FilterCommitmentProcessed and is used to iterate over the raw logs and unpacked data for CommitmentProcessed events raised by the Oracle contract.
type OracleCommitmentProcessedIterator struct {
	Event *OracleCommitmentProcessed // Event containing the contract specifics and raw log

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
func (it *OracleCommitmentProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleCommitmentProcessed)
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
		it.Event = new(OracleCommitmentProcessed)
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
func (it *OracleCommitmentProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleCommitmentProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleCommitmentProcessed represents a CommitmentProcessed event raised by the Oracle contract.
type OracleCommitmentProcessed struct {
	CommitmentIndex [32]byte
	IsSlash         bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCommitmentProcessed is a free log retrieval operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash)
func (_Oracle *OracleFilterer) FilterCommitmentProcessed(opts *bind.FilterOpts, commitmentIndex [][32]byte) (*OracleCommitmentProcessedIterator, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "CommitmentProcessed", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return &OracleCommitmentProcessedIterator{contract: _Oracle.contract, event: "CommitmentProcessed", logs: logs, sub: sub}, nil
}

// WatchCommitmentProcessed is a free log subscription operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash)
func (_Oracle *OracleFilterer) WatchCommitmentProcessed(opts *bind.WatchOpts, sink chan<- *OracleCommitmentProcessed, commitmentIndex [][32]byte) (event.Subscription, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "CommitmentProcessed", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleCommitmentProcessed)
				if err := _Oracle.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
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

// ParseCommitmentProcessed is a log parse operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash)
func (_Oracle *OracleFilterer) ParseCommitmentProcessed(log types.Log) (*OracleCommitmentProcessed, error) {
	event := new(OracleCommitmentProcessed)
	if err := _Oracle.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Oracle contract.
type OracleInitializedIterator struct {
	Event *OracleInitialized // Event containing the contract specifics and raw log

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
func (it *OracleInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleInitialized)
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
		it.Event = new(OracleInitialized)
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
func (it *OracleInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleInitialized represents a Initialized event raised by the Oracle contract.
type OracleInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Oracle *OracleFilterer) FilterInitialized(opts *bind.FilterOpts) (*OracleInitializedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OracleInitializedIterator{contract: _Oracle.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Oracle *OracleFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OracleInitialized) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleInitialized)
				if err := _Oracle.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseInitialized(log types.Log) (*OracleInitialized, error) {
	event := new(OracleInitialized)
	if err := _Oracle.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleOracleAccountSetIterator is returned from FilterOracleAccountSet and is used to iterate over the raw logs and unpacked data for OracleAccountSet events raised by the Oracle contract.
type OracleOracleAccountSetIterator struct {
	Event *OracleOracleAccountSet // Event containing the contract specifics and raw log

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
func (it *OracleOracleAccountSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOracleAccountSet)
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
		it.Event = new(OracleOracleAccountSet)
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
func (it *OracleOracleAccountSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleOracleAccountSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleOracleAccountSet represents a OracleAccountSet event raised by the Oracle contract.
type OracleOracleAccountSet struct {
	OldOracleAccount common.Address
	NewOracleAccount common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOracleAccountSet is a free log retrieval operation binding the contract event 0xc44093d4ba5b256ab49bc7bebfea8c049eb2f75fff6bcb9a8da6f8b1c92249e9.
//
// Solidity: event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount)
func (_Oracle *OracleFilterer) FilterOracleAccountSet(opts *bind.FilterOpts, oldOracleAccount []common.Address, newOracleAccount []common.Address) (*OracleOracleAccountSetIterator, error) {

	var oldOracleAccountRule []interface{}
	for _, oldOracleAccountItem := range oldOracleAccount {
		oldOracleAccountRule = append(oldOracleAccountRule, oldOracleAccountItem)
	}
	var newOracleAccountRule []interface{}
	for _, newOracleAccountItem := range newOracleAccount {
		newOracleAccountRule = append(newOracleAccountRule, newOracleAccountItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OracleAccountSet", oldOracleAccountRule, newOracleAccountRule)
	if err != nil {
		return nil, err
	}
	return &OracleOracleAccountSetIterator{contract: _Oracle.contract, event: "OracleAccountSet", logs: logs, sub: sub}, nil
}

// WatchOracleAccountSet is a free log subscription operation binding the contract event 0xc44093d4ba5b256ab49bc7bebfea8c049eb2f75fff6bcb9a8da6f8b1c92249e9.
//
// Solidity: event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount)
func (_Oracle *OracleFilterer) WatchOracleAccountSet(opts *bind.WatchOpts, sink chan<- *OracleOracleAccountSet, oldOracleAccount []common.Address, newOracleAccount []common.Address) (event.Subscription, error) {

	var oldOracleAccountRule []interface{}
	for _, oldOracleAccountItem := range oldOracleAccount {
		oldOracleAccountRule = append(oldOracleAccountRule, oldOracleAccountItem)
	}
	var newOracleAccountRule []interface{}
	for _, newOracleAccountItem := range newOracleAccount {
		newOracleAccountRule = append(newOracleAccountRule, newOracleAccountItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OracleAccountSet", oldOracleAccountRule, newOracleAccountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleOracleAccountSet)
				if err := _Oracle.contract.UnpackLog(event, "OracleAccountSet", log); err != nil {
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

// ParseOracleAccountSet is a log parse operation binding the contract event 0xc44093d4ba5b256ab49bc7bebfea8c049eb2f75fff6bcb9a8da6f8b1c92249e9.
//
// Solidity: event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount)
func (_Oracle *OracleFilterer) ParseOracleAccountSet(log types.Log) (*OracleOracleAccountSet, error) {
	event := new(OracleOracleAccountSet)
	if err := _Oracle.contract.UnpackLog(event, "OracleAccountSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Oracle contract.
type OracleOwnershipTransferStartedIterator struct {
	Event *OracleOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *OracleOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOwnershipTransferStarted)
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
		it.Event = new(OracleOwnershipTransferStarted)
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
func (it *OracleOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Oracle contract.
type OracleOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleOwnershipTransferStartedIterator{contract: _Oracle.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *OracleOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleOwnershipTransferStarted)
				if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseOwnershipTransferStarted(log types.Log) (*OracleOwnershipTransferStarted, error) {
	event := new(OracleOwnershipTransferStarted)
	if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Oracle contract.
type OracleOwnershipTransferredIterator struct {
	Event *OracleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OracleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleOwnershipTransferred)
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
		it.Event = new(OracleOwnershipTransferred)
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
func (it *OracleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleOwnershipTransferred represents a OwnershipTransferred event raised by the Oracle contract.
type OracleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OracleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OracleOwnershipTransferredIterator{contract: _Oracle.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Oracle *OracleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleOwnershipTransferred)
				if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseOwnershipTransferred(log types.Log) (*OracleOwnershipTransferred, error) {
	event := new(OracleOwnershipTransferred)
	if err := _Oracle.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OraclePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Oracle contract.
type OraclePausedIterator struct {
	Event *OraclePaused // Event containing the contract specifics and raw log

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
func (it *OraclePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OraclePaused)
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
		it.Event = new(OraclePaused)
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
func (it *OraclePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OraclePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OraclePaused represents a Paused event raised by the Oracle contract.
type OraclePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Oracle *OracleFilterer) FilterPaused(opts *bind.FilterOpts) (*OraclePausedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &OraclePausedIterator{contract: _Oracle.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Oracle *OracleFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *OraclePaused) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OraclePaused)
				if err := _Oracle.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Oracle *OracleFilterer) ParsePaused(log types.Log) (*OraclePaused, error) {
	event := new(OraclePaused)
	if err := _Oracle.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OraclePreconfManagerSetIterator is returned from FilterPreconfManagerSet and is used to iterate over the raw logs and unpacked data for PreconfManagerSet events raised by the Oracle contract.
type OraclePreconfManagerSetIterator struct {
	Event *OraclePreconfManagerSet // Event containing the contract specifics and raw log

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
func (it *OraclePreconfManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OraclePreconfManagerSet)
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
		it.Event = new(OraclePreconfManagerSet)
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
func (it *OraclePreconfManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OraclePreconfManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OraclePreconfManagerSet represents a PreconfManagerSet event raised by the Oracle contract.
type OraclePreconfManagerSet struct {
	NewPreconfManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPreconfManagerSet is a free log retrieval operation binding the contract event 0xb30b787197e154c5c4d294707cb846f2999047cd1d9209d7aec68fe010584acf.
//
// Solidity: event PreconfManagerSet(address indexed newPreconfManager)
func (_Oracle *OracleFilterer) FilterPreconfManagerSet(opts *bind.FilterOpts, newPreconfManager []common.Address) (*OraclePreconfManagerSetIterator, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "PreconfManagerSet", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return &OraclePreconfManagerSetIterator{contract: _Oracle.contract, event: "PreconfManagerSet", logs: logs, sub: sub}, nil
}

// WatchPreconfManagerSet is a free log subscription operation binding the contract event 0xb30b787197e154c5c4d294707cb846f2999047cd1d9209d7aec68fe010584acf.
//
// Solidity: event PreconfManagerSet(address indexed newPreconfManager)
func (_Oracle *OracleFilterer) WatchPreconfManagerSet(opts *bind.WatchOpts, sink chan<- *OraclePreconfManagerSet, newPreconfManager []common.Address) (event.Subscription, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "PreconfManagerSet", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OraclePreconfManagerSet)
				if err := _Oracle.contract.UnpackLog(event, "PreconfManagerSet", log); err != nil {
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

// ParsePreconfManagerSet is a log parse operation binding the contract event 0xb30b787197e154c5c4d294707cb846f2999047cd1d9209d7aec68fe010584acf.
//
// Solidity: event PreconfManagerSet(address indexed newPreconfManager)
func (_Oracle *OracleFilterer) ParsePreconfManagerSet(log types.Log) (*OraclePreconfManagerSet, error) {
	event := new(OraclePreconfManagerSet)
	if err := _Oracle.contract.UnpackLog(event, "PreconfManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleProviderRegistrySetIterator is returned from FilterProviderRegistrySet and is used to iterate over the raw logs and unpacked data for ProviderRegistrySet events raised by the Oracle contract.
type OracleProviderRegistrySetIterator struct {
	Event *OracleProviderRegistrySet // Event containing the contract specifics and raw log

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
func (it *OracleProviderRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleProviderRegistrySet)
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
		it.Event = new(OracleProviderRegistrySet)
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
func (it *OracleProviderRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleProviderRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleProviderRegistrySet represents a ProviderRegistrySet event raised by the Oracle contract.
type OracleProviderRegistrySet struct {
	NewProviderRegistry common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterProviderRegistrySet is a free log retrieval operation binding the contract event 0x1f0fa16c192be2cff8a73ddaf1394dc78718ff634425739ee1df537db3b58eb1.
//
// Solidity: event ProviderRegistrySet(address indexed newProviderRegistry)
func (_Oracle *OracleFilterer) FilterProviderRegistrySet(opts *bind.FilterOpts, newProviderRegistry []common.Address) (*OracleProviderRegistrySetIterator, error) {

	var newProviderRegistryRule []interface{}
	for _, newProviderRegistryItem := range newProviderRegistry {
		newProviderRegistryRule = append(newProviderRegistryRule, newProviderRegistryItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "ProviderRegistrySet", newProviderRegistryRule)
	if err != nil {
		return nil, err
	}
	return &OracleProviderRegistrySetIterator{contract: _Oracle.contract, event: "ProviderRegistrySet", logs: logs, sub: sub}, nil
}

// WatchProviderRegistrySet is a free log subscription operation binding the contract event 0x1f0fa16c192be2cff8a73ddaf1394dc78718ff634425739ee1df537db3b58eb1.
//
// Solidity: event ProviderRegistrySet(address indexed newProviderRegistry)
func (_Oracle *OracleFilterer) WatchProviderRegistrySet(opts *bind.WatchOpts, sink chan<- *OracleProviderRegistrySet, newProviderRegistry []common.Address) (event.Subscription, error) {

	var newProviderRegistryRule []interface{}
	for _, newProviderRegistryItem := range newProviderRegistry {
		newProviderRegistryRule = append(newProviderRegistryRule, newProviderRegistryItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "ProviderRegistrySet", newProviderRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleProviderRegistrySet)
				if err := _Oracle.contract.UnpackLog(event, "ProviderRegistrySet", log); err != nil {
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

// ParseProviderRegistrySet is a log parse operation binding the contract event 0x1f0fa16c192be2cff8a73ddaf1394dc78718ff634425739ee1df537db3b58eb1.
//
// Solidity: event ProviderRegistrySet(address indexed newProviderRegistry)
func (_Oracle *OracleFilterer) ParseProviderRegistrySet(log types.Log) (*OracleProviderRegistrySet, error) {
	event := new(OracleProviderRegistrySet)
	if err := _Oracle.contract.UnpackLog(event, "ProviderRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Oracle contract.
type OracleUnpausedIterator struct {
	Event *OracleUnpaused // Event containing the contract specifics and raw log

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
func (it *OracleUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleUnpaused)
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
		it.Event = new(OracleUnpaused)
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
func (it *OracleUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleUnpaused represents a Unpaused event raised by the Oracle contract.
type OracleUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Oracle *OracleFilterer) FilterUnpaused(opts *bind.FilterOpts) (*OracleUnpausedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &OracleUnpausedIterator{contract: _Oracle.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Oracle *OracleFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *OracleUnpaused) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleUnpaused)
				if err := _Oracle.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseUnpaused(log types.Log) (*OracleUnpaused, error) {
	event := new(OracleUnpaused)
	if err := _Oracle.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Oracle contract.
type OracleUpgradedIterator struct {
	Event *OracleUpgraded // Event containing the contract specifics and raw log

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
func (it *OracleUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleUpgraded)
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
		it.Event = new(OracleUpgraded)
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
func (it *OracleUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleUpgraded represents a Upgraded event raised by the Oracle contract.
type OracleUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Oracle *OracleFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*OracleUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &OracleUpgradedIterator{contract: _Oracle.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Oracle *OracleFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *OracleUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleUpgraded)
				if err := _Oracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Oracle *OracleFilterer) ParseUpgraded(log types.Log) (*OracleUpgraded, error) {
	event := new(OracleUpgraded)
	if err := _Oracle.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
