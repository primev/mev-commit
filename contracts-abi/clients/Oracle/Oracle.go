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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_preConfContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_nextRequestedBlockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addBuilderAddress\",\"inputs\":[{\"name\":\"builderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"builderAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockBuilderNameToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuilder\",\"inputs\":[{\"name\":\"builderNameGrafiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNextRequestedBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"moveToNextBlock\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextRequestedBlockNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processBuilderCommitmentForBlockNumber\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockBuilderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNextBlock\",\"inputs\":[{\"name\":\"newBlockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockFunds\",\"inputs\":[{\"name\":\"bidIDs\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CommitmentProcessed\",\"inputs\":[{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"isSlash\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
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

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Oracle *OracleCaller) GetBuilder(opts *bind.CallOpts, builderNameGrafiti string) (common.Address, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getBuilder", builderNameGrafiti)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Oracle *OracleSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _Oracle.Contract.GetBuilder(&_Oracle.CallOpts, builderNameGrafiti)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Oracle *OracleCallerSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _Oracle.Contract.GetBuilder(&_Oracle.CallOpts, builderNameGrafiti)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleCaller) GetNextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "getNextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _Oracle.Contract.GetNextRequestedBlockNumber(&_Oracle.CallOpts)
}

// GetNextRequestedBlockNumber is a free data retrieval call binding the contract method 0xfce2a502.
//
// Solidity: function getNextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleCallerSession) GetNextRequestedBlockNumber() (*big.Int, error) {
	return _Oracle.Contract.GetNextRequestedBlockNumber(&_Oracle.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleCaller) NextRequestedBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Oracle.contract.Call(opts, &out, "nextRequestedBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _Oracle.Contract.NextRequestedBlockNumber(&_Oracle.CallOpts)
}

// NextRequestedBlockNumber is a free data retrieval call binding the contract method 0xc512c561.
//
// Solidity: function nextRequestedBlockNumber() view returns(uint256)
func (_Oracle *OracleCallerSession) NextRequestedBlockNumber() (*big.Int, error) {
	return _Oracle.Contract.NextRequestedBlockNumber(&_Oracle.CallOpts)
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

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Oracle *OracleTransactor) AddBuilderAddress(opts *bind.TransactOpts, builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "addBuilderAddress", builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Oracle *OracleSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.AddBuilderAddress(&_Oracle.TransactOpts, builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Oracle *OracleTransactorSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.AddBuilderAddress(&_Oracle.TransactOpts, builderName, builderAddress)
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_Oracle *OracleTransactor) MoveToNextBlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "moveToNextBlock")
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_Oracle *OracleSession) MoveToNextBlock() (*types.Transaction, error) {
	return _Oracle.Contract.MoveToNextBlock(&_Oracle.TransactOpts)
}

// MoveToNextBlock is a paid mutator transaction binding the contract method 0x32289b72.
//
// Solidity: function moveToNextBlock() returns()
func (_Oracle *OracleTransactorSession) MoveToNextBlock() (*types.Transaction, error) {
	return _Oracle.Contract.MoveToNextBlock(&_Oracle.TransactOpts)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x2e9bd11f.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleTransactor) ProcessBuilderCommitmentForBlockNumber(opts *bind.TransactOpts, commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "processBuilderCommitmentForBlockNumber", commitmentIndex, blockNumber, blockBuilderName, isSlash, residualBidPercentAfterDecay)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x2e9bd11f.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.ProcessBuilderCommitmentForBlockNumber(&_Oracle.TransactOpts, commitmentIndex, blockNumber, blockBuilderName, isSlash, residualBidPercentAfterDecay)
}

// ProcessBuilderCommitmentForBlockNumber is a paid mutator transaction binding the contract method 0x2e9bd11f.
//
// Solidity: function processBuilderCommitmentForBlockNumber(bytes32 commitmentIndex, uint256 blockNumber, string blockBuilderName, bool isSlash, uint256 residualBidPercentAfterDecay) returns()
func (_Oracle *OracleTransactorSession) ProcessBuilderCommitmentForBlockNumber(commitmentIndex [32]byte, blockNumber *big.Int, blockBuilderName string, isSlash bool, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.ProcessBuilderCommitmentForBlockNumber(&_Oracle.TransactOpts, commitmentIndex, blockNumber, blockBuilderName, isSlash, residualBidPercentAfterDecay)
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

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_Oracle *OracleTransactor) SetNextBlock(opts *bind.TransactOpts, newBlockNumber uint64) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "setNextBlock", newBlockNumber)
}

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_Oracle *OracleSession) SetNextBlock(newBlockNumber uint64) (*types.Transaction, error) {
	return _Oracle.Contract.SetNextBlock(&_Oracle.TransactOpts, newBlockNumber)
}

// SetNextBlock is a paid mutator transaction binding the contract method 0x19072add.
//
// Solidity: function setNextBlock(uint64 newBlockNumber) returns()
func (_Oracle *OracleTransactorSession) SetNextBlock(newBlockNumber uint64) (*types.Transaction, error) {
	return _Oracle.Contract.SetNextBlock(&_Oracle.TransactOpts, newBlockNumber)
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

// UnlockFunds is a paid mutator transaction binding the contract method 0x9dd4152b.
//
// Solidity: function unlockFunds(bytes32[] bidIDs) returns()
func (_Oracle *OracleTransactor) UnlockFunds(opts *bind.TransactOpts, bidIDs [][32]byte) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "unlockFunds", bidIDs)
}

// UnlockFunds is a paid mutator transaction binding the contract method 0x9dd4152b.
//
// Solidity: function unlockFunds(bytes32[] bidIDs) returns()
func (_Oracle *OracleSession) UnlockFunds(bidIDs [][32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.UnlockFunds(&_Oracle.TransactOpts, bidIDs)
}

// UnlockFunds is a paid mutator transaction binding the contract method 0x9dd4152b.
//
// Solidity: function unlockFunds(bytes32[] bidIDs) returns()
func (_Oracle *OracleTransactorSession) UnlockFunds(bidIDs [][32]byte) (*types.Transaction, error) {
	return _Oracle.Contract.UnlockFunds(&_Oracle.TransactOpts, bidIDs)
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
	CommitmentHash [32]byte
	IsSlash        bool
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitmentProcessed is a free log retrieval operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Oracle *OracleFilterer) FilterCommitmentProcessed(opts *bind.FilterOpts) (*OracleCommitmentProcessedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "CommitmentProcessed")
	if err != nil {
		return nil, err
	}
	return &OracleCommitmentProcessedIterator{contract: _Oracle.contract, event: "CommitmentProcessed", logs: logs, sub: sub}, nil
}

// WatchCommitmentProcessed is a free log subscription operation binding the contract event 0xddc1768a3a762a04e5fd3abea8ae3b60e23bcf290f4a032280e6a726611d41f5.
//
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Oracle *OracleFilterer) WatchCommitmentProcessed(opts *bind.WatchOpts, sink chan<- *OracleCommitmentProcessed) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "CommitmentProcessed")
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
// Solidity: event CommitmentProcessed(bytes32 commitmentHash, bool isSlash)
func (_Oracle *OracleFilterer) ParseCommitmentProcessed(log types.Log) (*OracleCommitmentProcessed, error) {
	event := new(OracleCommitmentProcessed)
	if err := _Oracle.contract.UnpackLog(event, "CommitmentProcessed", log); err != nil {
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
