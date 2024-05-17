// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blocktracker

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

// BlocktrackerMetaData contains all meta data concerning the Blocktracker contract.
var BlocktrackerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"addBuilderAddress\",\"inputs\":[{\"name\":\"builderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"builderAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockBuilderNameToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockWinners\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockWinner\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuilder\",\"inputs\":[{\"name\":\"builderNameGrafiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWindowFromBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordL1Block\",\"inputs\":[{\"name\":\"_blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_winnerGraffiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewBlocksPerWindow\",\"inputs\":[{\"name\":\"blocksPerWindow\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewL1Block\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"winner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewWindow\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// BlocktrackerABI is the input ABI used to generate the binding from.
// Deprecated: Use BlocktrackerMetaData.ABI instead.
var BlocktrackerABI = BlocktrackerMetaData.ABI

// Blocktracker is an auto generated Go binding around an Ethereum contract.
type Blocktracker struct {
	BlocktrackerCaller     // Read-only binding to the contract
	BlocktrackerTransactor // Write-only binding to the contract
	BlocktrackerFilterer   // Log filterer for contract events
}

// BlocktrackerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlocktrackerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlocktrackerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlocktrackerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlocktrackerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlocktrackerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlocktrackerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlocktrackerSession struct {
	Contract     *Blocktracker     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BlocktrackerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlocktrackerCallerSession struct {
	Contract *BlocktrackerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BlocktrackerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlocktrackerTransactorSession struct {
	Contract     *BlocktrackerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BlocktrackerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlocktrackerRaw struct {
	Contract *Blocktracker // Generic contract binding to access the raw methods on
}

// BlocktrackerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlocktrackerCallerRaw struct {
	Contract *BlocktrackerCaller // Generic read-only contract binding to access the raw methods on
}

// BlocktrackerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlocktrackerTransactorRaw struct {
	Contract *BlocktrackerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlocktracker creates a new instance of Blocktracker, bound to a specific deployed contract.
func NewBlocktracker(address common.Address, backend bind.ContractBackend) (*Blocktracker, error) {
	contract, err := bindBlocktracker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Blocktracker{BlocktrackerCaller: BlocktrackerCaller{contract: contract}, BlocktrackerTransactor: BlocktrackerTransactor{contract: contract}, BlocktrackerFilterer: BlocktrackerFilterer{contract: contract}}, nil
}

// NewBlocktrackerCaller creates a new read-only instance of Blocktracker, bound to a specific deployed contract.
func NewBlocktrackerCaller(address common.Address, caller bind.ContractCaller) (*BlocktrackerCaller, error) {
	contract, err := bindBlocktracker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerCaller{contract: contract}, nil
}

// NewBlocktrackerTransactor creates a new write-only instance of Blocktracker, bound to a specific deployed contract.
func NewBlocktrackerTransactor(address common.Address, transactor bind.ContractTransactor) (*BlocktrackerTransactor, error) {
	contract, err := bindBlocktracker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerTransactor{contract: contract}, nil
}

// NewBlocktrackerFilterer creates a new log filterer instance of Blocktracker, bound to a specific deployed contract.
func NewBlocktrackerFilterer(address common.Address, filterer bind.ContractFilterer) (*BlocktrackerFilterer, error) {
	contract, err := bindBlocktracker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerFilterer{contract: contract}, nil
}

// bindBlocktracker binds a generic wrapper to an already deployed contract.
func bindBlocktracker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlocktrackerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blocktracker *BlocktrackerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blocktracker.Contract.BlocktrackerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blocktracker *BlocktrackerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.Contract.BlocktrackerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blocktracker *BlocktrackerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blocktracker.Contract.BlocktrackerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blocktracker *BlocktrackerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blocktracker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blocktracker *BlocktrackerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blocktracker *BlocktrackerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blocktracker.Contract.contract.Transact(opts, method, params...)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Blocktracker *BlocktrackerCaller) BlockBuilderNameToAddress(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "blockBuilderNameToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Blocktracker *BlocktrackerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Blocktracker.Contract.BlockBuilderNameToAddress(&_Blocktracker.CallOpts, arg0)
}

// BlockBuilderNameToAddress is a free data retrieval call binding the contract method 0xeebac3ac.
//
// Solidity: function blockBuilderNameToAddress(string ) view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) BlockBuilderNameToAddress(arg0 string) (common.Address, error) {
	return _Blocktracker.Contract.BlockBuilderNameToAddress(&_Blocktracker.CallOpts, arg0)
}

// BlockWinners is a free data retrieval call binding the contract method 0xe4747419.
//
// Solidity: function blockWinners(uint256 ) view returns(address)
func (_Blocktracker *BlocktrackerCaller) BlockWinners(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "blockWinners", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockWinners is a free data retrieval call binding the contract method 0xe4747419.
//
// Solidity: function blockWinners(uint256 ) view returns(address)
func (_Blocktracker *BlocktrackerSession) BlockWinners(arg0 *big.Int) (common.Address, error) {
	return _Blocktracker.Contract.BlockWinners(&_Blocktracker.CallOpts, arg0)
}

// BlockWinners is a free data retrieval call binding the contract method 0xe4747419.
//
// Solidity: function blockWinners(uint256 ) view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) BlockWinners(arg0 *big.Int) (common.Address, error) {
	return _Blocktracker.Contract.BlockWinners(&_Blocktracker.CallOpts, arg0)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCaller) BlocksPerWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "blocksPerWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerSession) BlocksPerWindow() (*big.Int, error) {
	return _Blocktracker.Contract.BlocksPerWindow(&_Blocktracker.CallOpts)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCallerSession) BlocksPerWindow() (*big.Int, error) {
	return _Blocktracker.Contract.BlocksPerWindow(&_Blocktracker.CallOpts)
}

// CurrentWindow is a free data retrieval call binding the contract method 0xba0bafb4.
//
// Solidity: function currentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCaller) CurrentWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "currentWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentWindow is a free data retrieval call binding the contract method 0xba0bafb4.
//
// Solidity: function currentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerSession) CurrentWindow() (*big.Int, error) {
	return _Blocktracker.Contract.CurrentWindow(&_Blocktracker.CallOpts)
}

// CurrentWindow is a free data retrieval call binding the contract method 0xba0bafb4.
//
// Solidity: function currentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCallerSession) CurrentWindow() (*big.Int, error) {
	return _Blocktracker.Contract.CurrentWindow(&_Blocktracker.CallOpts)
}

// GetBlockWinner is a free data retrieval call binding the contract method 0x6753ab34.
//
// Solidity: function getBlockWinner(uint256 blockNumber) view returns(address)
func (_Blocktracker *BlocktrackerCaller) GetBlockWinner(opts *bind.CallOpts, blockNumber *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getBlockWinner", blockNumber)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBlockWinner is a free data retrieval call binding the contract method 0x6753ab34.
//
// Solidity: function getBlockWinner(uint256 blockNumber) view returns(address)
func (_Blocktracker *BlocktrackerSession) GetBlockWinner(blockNumber *big.Int) (common.Address, error) {
	return _Blocktracker.Contract.GetBlockWinner(&_Blocktracker.CallOpts, blockNumber)
}

// GetBlockWinner is a free data retrieval call binding the contract method 0x6753ab34.
//
// Solidity: function getBlockWinner(uint256 blockNumber) view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) GetBlockWinner(blockNumber *big.Int) (common.Address, error) {
	return _Blocktracker.Contract.GetBlockWinner(&_Blocktracker.CallOpts, blockNumber)
}

// GetBlocksPerWindow is a free data retrieval call binding the contract method 0x8711a019.
//
// Solidity: function getBlocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCaller) GetBlocksPerWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getBlocksPerWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksPerWindow is a free data retrieval call binding the contract method 0x8711a019.
//
// Solidity: function getBlocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerSession) GetBlocksPerWindow() (*big.Int, error) {
	return _Blocktracker.Contract.GetBlocksPerWindow(&_Blocktracker.CallOpts)
}

// GetBlocksPerWindow is a free data retrieval call binding the contract method 0x8711a019.
//
// Solidity: function getBlocksPerWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCallerSession) GetBlocksPerWindow() (*big.Int, error) {
	return _Blocktracker.Contract.GetBlocksPerWindow(&_Blocktracker.CallOpts)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Blocktracker *BlocktrackerCaller) GetBuilder(opts *bind.CallOpts, builderNameGrafiti string) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getBuilder", builderNameGrafiti)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Blocktracker *BlocktrackerSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _Blocktracker.Contract.GetBuilder(&_Blocktracker.CallOpts, builderNameGrafiti)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGrafiti) view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) GetBuilder(builderNameGrafiti string) (common.Address, error) {
	return _Blocktracker.Contract.GetBuilder(&_Blocktracker.CallOpts, builderNameGrafiti)
}

// GetCurrentWindow is a free data retrieval call binding the contract method 0x0f67e7d5.
//
// Solidity: function getCurrentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCaller) GetCurrentWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getCurrentWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentWindow is a free data retrieval call binding the contract method 0x0f67e7d5.
//
// Solidity: function getCurrentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerSession) GetCurrentWindow() (*big.Int, error) {
	return _Blocktracker.Contract.GetCurrentWindow(&_Blocktracker.CallOpts)
}

// GetCurrentWindow is a free data retrieval call binding the contract method 0x0f67e7d5.
//
// Solidity: function getCurrentWindow() view returns(uint256)
func (_Blocktracker *BlocktrackerCallerSession) GetCurrentWindow() (*big.Int, error) {
	return _Blocktracker.Contract.GetCurrentWindow(&_Blocktracker.CallOpts)
}

// GetWindowFromBlockNumber is a free data retrieval call binding the contract method 0x70ef60db.
//
// Solidity: function getWindowFromBlockNumber(uint256 blockNumber) view returns(uint256)
func (_Blocktracker *BlocktrackerCaller) GetWindowFromBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getWindowFromBlockNumber", blockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWindowFromBlockNumber is a free data retrieval call binding the contract method 0x70ef60db.
//
// Solidity: function getWindowFromBlockNumber(uint256 blockNumber) view returns(uint256)
func (_Blocktracker *BlocktrackerSession) GetWindowFromBlockNumber(blockNumber *big.Int) (*big.Int, error) {
	return _Blocktracker.Contract.GetWindowFromBlockNumber(&_Blocktracker.CallOpts, blockNumber)
}

// GetWindowFromBlockNumber is a free data retrieval call binding the contract method 0x70ef60db.
//
// Solidity: function getWindowFromBlockNumber(uint256 blockNumber) view returns(uint256)
func (_Blocktracker *BlocktrackerCallerSession) GetWindowFromBlockNumber(blockNumber *big.Int) (*big.Int, error) {
	return _Blocktracker.Contract.GetWindowFromBlockNumber(&_Blocktracker.CallOpts, blockNumber)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blocktracker *BlocktrackerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blocktracker *BlocktrackerSession) Owner() (common.Address, error) {
	return _Blocktracker.Contract.Owner(&_Blocktracker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) Owner() (common.Address, error) {
	return _Blocktracker.Contract.Owner(&_Blocktracker.CallOpts)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Blocktracker *BlocktrackerTransactor) AddBuilderAddress(opts *bind.TransactOpts, builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "addBuilderAddress", builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Blocktracker *BlocktrackerSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.AddBuilderAddress(&_Blocktracker.TransactOpts, builderName, builderAddress)
}

// AddBuilderAddress is a paid mutator transaction binding the contract method 0x0bd0a9e1.
//
// Solidity: function addBuilderAddress(string builderName, address builderAddress) returns()
func (_Blocktracker *BlocktrackerTransactorSession) AddBuilderAddress(builderName string, builderAddress common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.AddBuilderAddress(&_Blocktracker.TransactOpts, builderName, builderAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Blocktracker *BlocktrackerTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "initialize", _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Blocktracker *BlocktrackerSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.Initialize(&_Blocktracker.TransactOpts, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_Blocktracker *BlocktrackerTransactorSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.Initialize(&_Blocktracker.TransactOpts, _owner)
}

// RecordL1Block is a paid mutator transaction binding the contract method 0x1d63f108.
//
// Solidity: function recordL1Block(uint256 _blockNumber, string _winnerGraffiti) returns()
func (_Blocktracker *BlocktrackerTransactor) RecordL1Block(opts *bind.TransactOpts, _blockNumber *big.Int, _winnerGraffiti string) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "recordL1Block", _blockNumber, _winnerGraffiti)
}

// RecordL1Block is a paid mutator transaction binding the contract method 0x1d63f108.
//
// Solidity: function recordL1Block(uint256 _blockNumber, string _winnerGraffiti) returns()
func (_Blocktracker *BlocktrackerSession) RecordL1Block(_blockNumber *big.Int, _winnerGraffiti string) (*types.Transaction, error) {
	return _Blocktracker.Contract.RecordL1Block(&_Blocktracker.TransactOpts, _blockNumber, _winnerGraffiti)
}

// RecordL1Block is a paid mutator transaction binding the contract method 0x1d63f108.
//
// Solidity: function recordL1Block(uint256 _blockNumber, string _winnerGraffiti) returns()
func (_Blocktracker *BlocktrackerTransactorSession) RecordL1Block(_blockNumber *big.Int, _winnerGraffiti string) (*types.Transaction, error) {
	return _Blocktracker.Contract.RecordL1Block(&_Blocktracker.TransactOpts, _blockNumber, _winnerGraffiti)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blocktracker *BlocktrackerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blocktracker *BlocktrackerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Blocktracker.Contract.RenounceOwnership(&_Blocktracker.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blocktracker *BlocktrackerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Blocktracker.Contract.RenounceOwnership(&_Blocktracker.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blocktracker *BlocktrackerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blocktracker *BlocktrackerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.TransferOwnership(&_Blocktracker.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blocktracker *BlocktrackerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.TransferOwnership(&_Blocktracker.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blocktracker *BlocktrackerTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Blocktracker.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blocktracker *BlocktrackerSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Blocktracker.Contract.Fallback(&_Blocktracker.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blocktracker *BlocktrackerTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Blocktracker.Contract.Fallback(&_Blocktracker.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blocktracker *BlocktrackerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blocktracker *BlocktrackerSession) Receive() (*types.Transaction, error) {
	return _Blocktracker.Contract.Receive(&_Blocktracker.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blocktracker *BlocktrackerTransactorSession) Receive() (*types.Transaction, error) {
	return _Blocktracker.Contract.Receive(&_Blocktracker.TransactOpts)
}

// BlocktrackerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Blocktracker contract.
type BlocktrackerInitializedIterator struct {
	Event *BlocktrackerInitialized // Event containing the contract specifics and raw log

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
func (it *BlocktrackerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerInitialized)
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
		it.Event = new(BlocktrackerInitialized)
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
func (it *BlocktrackerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerInitialized represents a Initialized event raised by the Blocktracker contract.
type BlocktrackerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Blocktracker *BlocktrackerFilterer) FilterInitialized(opts *bind.FilterOpts) (*BlocktrackerInitializedIterator, error) {

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BlocktrackerInitializedIterator{contract: _Blocktracker.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Blocktracker *BlocktrackerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BlocktrackerInitialized) (event.Subscription, error) {

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerInitialized)
				if err := _Blocktracker.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseInitialized(log types.Log) (*BlocktrackerInitialized, error) {
	event := new(BlocktrackerInitialized)
	if err := _Blocktracker.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerNewBlocksPerWindowIterator is returned from FilterNewBlocksPerWindow and is used to iterate over the raw logs and unpacked data for NewBlocksPerWindow events raised by the Blocktracker contract.
type BlocktrackerNewBlocksPerWindowIterator struct {
	Event *BlocktrackerNewBlocksPerWindow // Event containing the contract specifics and raw log

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
func (it *BlocktrackerNewBlocksPerWindowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerNewBlocksPerWindow)
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
		it.Event = new(BlocktrackerNewBlocksPerWindow)
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
func (it *BlocktrackerNewBlocksPerWindowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerNewBlocksPerWindowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerNewBlocksPerWindow represents a NewBlocksPerWindow event raised by the Blocktracker contract.
type BlocktrackerNewBlocksPerWindow struct {
	BlocksPerWindow *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewBlocksPerWindow is a free log retrieval operation binding the contract event 0xd2abfe7e41e182f2121a97e57d5133f2ccef005fb15e25ef2f09d4a6657e20e0.
//
// Solidity: event NewBlocksPerWindow(uint256 blocksPerWindow)
func (_Blocktracker *BlocktrackerFilterer) FilterNewBlocksPerWindow(opts *bind.FilterOpts) (*BlocktrackerNewBlocksPerWindowIterator, error) {

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "NewBlocksPerWindow")
	if err != nil {
		return nil, err
	}
	return &BlocktrackerNewBlocksPerWindowIterator{contract: _Blocktracker.contract, event: "NewBlocksPerWindow", logs: logs, sub: sub}, nil
}

// WatchNewBlocksPerWindow is a free log subscription operation binding the contract event 0xd2abfe7e41e182f2121a97e57d5133f2ccef005fb15e25ef2f09d4a6657e20e0.
//
// Solidity: event NewBlocksPerWindow(uint256 blocksPerWindow)
func (_Blocktracker *BlocktrackerFilterer) WatchNewBlocksPerWindow(opts *bind.WatchOpts, sink chan<- *BlocktrackerNewBlocksPerWindow) (event.Subscription, error) {

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "NewBlocksPerWindow")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerNewBlocksPerWindow)
				if err := _Blocktracker.contract.UnpackLog(event, "NewBlocksPerWindow", log); err != nil {
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

// ParseNewBlocksPerWindow is a log parse operation binding the contract event 0xd2abfe7e41e182f2121a97e57d5133f2ccef005fb15e25ef2f09d4a6657e20e0.
//
// Solidity: event NewBlocksPerWindow(uint256 blocksPerWindow)
func (_Blocktracker *BlocktrackerFilterer) ParseNewBlocksPerWindow(log types.Log) (*BlocktrackerNewBlocksPerWindow, error) {
	event := new(BlocktrackerNewBlocksPerWindow)
	if err := _Blocktracker.contract.UnpackLog(event, "NewBlocksPerWindow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerNewL1BlockIterator is returned from FilterNewL1Block and is used to iterate over the raw logs and unpacked data for NewL1Block events raised by the Blocktracker contract.
type BlocktrackerNewL1BlockIterator struct {
	Event *BlocktrackerNewL1Block // Event containing the contract specifics and raw log

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
func (it *BlocktrackerNewL1BlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerNewL1Block)
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
		it.Event = new(BlocktrackerNewL1Block)
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
func (it *BlocktrackerNewL1BlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerNewL1BlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerNewL1Block represents a NewL1Block event raised by the Blocktracker contract.
type BlocktrackerNewL1Block struct {
	BlockNumber *big.Int
	Winner      common.Address
	Window      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewL1Block is a free log retrieval operation binding the contract event 0x8323d3e5d25db513e1a772870aaa45e9b069a13d49879d72e70638b5c1c18cb7.
//
// Solidity: event NewL1Block(uint256 indexed blockNumber, address indexed winner, uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) FilterNewL1Block(opts *bind.FilterOpts, blockNumber []*big.Int, winner []common.Address, window []*big.Int) (*BlocktrackerNewL1BlockIterator, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "NewL1Block", blockNumberRule, winnerRule, windowRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerNewL1BlockIterator{contract: _Blocktracker.contract, event: "NewL1Block", logs: logs, sub: sub}, nil
}

// WatchNewL1Block is a free log subscription operation binding the contract event 0x8323d3e5d25db513e1a772870aaa45e9b069a13d49879d72e70638b5c1c18cb7.
//
// Solidity: event NewL1Block(uint256 indexed blockNumber, address indexed winner, uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) WatchNewL1Block(opts *bind.WatchOpts, sink chan<- *BlocktrackerNewL1Block, blockNumber []*big.Int, winner []common.Address, window []*big.Int) (event.Subscription, error) {

	var blockNumberRule []interface{}
	for _, blockNumberItem := range blockNumber {
		blockNumberRule = append(blockNumberRule, blockNumberItem)
	}
	var winnerRule []interface{}
	for _, winnerItem := range winner {
		winnerRule = append(winnerRule, winnerItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "NewL1Block", blockNumberRule, winnerRule, windowRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerNewL1Block)
				if err := _Blocktracker.contract.UnpackLog(event, "NewL1Block", log); err != nil {
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

// ParseNewL1Block is a log parse operation binding the contract event 0x8323d3e5d25db513e1a772870aaa45e9b069a13d49879d72e70638b5c1c18cb7.
//
// Solidity: event NewL1Block(uint256 indexed blockNumber, address indexed winner, uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) ParseNewL1Block(log types.Log) (*BlocktrackerNewL1Block, error) {
	event := new(BlocktrackerNewL1Block)
	if err := _Blocktracker.contract.UnpackLog(event, "NewL1Block", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerNewWindowIterator is returned from FilterNewWindow and is used to iterate over the raw logs and unpacked data for NewWindow events raised by the Blocktracker contract.
type BlocktrackerNewWindowIterator struct {
	Event *BlocktrackerNewWindow // Event containing the contract specifics and raw log

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
func (it *BlocktrackerNewWindowIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerNewWindow)
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
		it.Event = new(BlocktrackerNewWindow)
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
func (it *BlocktrackerNewWindowIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerNewWindowIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerNewWindow represents a NewWindow event raised by the Blocktracker contract.
type BlocktrackerNewWindow struct {
	Window *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNewWindow is a free log retrieval operation binding the contract event 0x6553ddbc9c04d825543b5b531877439a6abdb68d5825c39f2dd3798e54118870.
//
// Solidity: event NewWindow(uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) FilterNewWindow(opts *bind.FilterOpts, window []*big.Int) (*BlocktrackerNewWindowIterator, error) {

	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "NewWindow", windowRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerNewWindowIterator{contract: _Blocktracker.contract, event: "NewWindow", logs: logs, sub: sub}, nil
}

// WatchNewWindow is a free log subscription operation binding the contract event 0x6553ddbc9c04d825543b5b531877439a6abdb68d5825c39f2dd3798e54118870.
//
// Solidity: event NewWindow(uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) WatchNewWindow(opts *bind.WatchOpts, sink chan<- *BlocktrackerNewWindow, window []*big.Int) (event.Subscription, error) {

	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "NewWindow", windowRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerNewWindow)
				if err := _Blocktracker.contract.UnpackLog(event, "NewWindow", log); err != nil {
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

// ParseNewWindow is a log parse operation binding the contract event 0x6553ddbc9c04d825543b5b531877439a6abdb68d5825c39f2dd3798e54118870.
//
// Solidity: event NewWindow(uint256 indexed window)
func (_Blocktracker *BlocktrackerFilterer) ParseNewWindow(log types.Log) (*BlocktrackerNewWindow, error) {
	event := new(BlocktrackerNewWindow)
	if err := _Blocktracker.contract.UnpackLog(event, "NewWindow", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Blocktracker contract.
type BlocktrackerOwnershipTransferredIterator struct {
	Event *BlocktrackerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BlocktrackerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerOwnershipTransferred)
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
		it.Event = new(BlocktrackerOwnershipTransferred)
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
func (it *BlocktrackerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerOwnershipTransferred represents a OwnershipTransferred event raised by the Blocktracker contract.
type BlocktrackerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Blocktracker *BlocktrackerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BlocktrackerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerOwnershipTransferredIterator{contract: _Blocktracker.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Blocktracker *BlocktrackerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BlocktrackerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerOwnershipTransferred)
				if err := _Blocktracker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseOwnershipTransferred(log types.Log) (*BlocktrackerOwnershipTransferred, error) {
	event := new(BlocktrackerOwnershipTransferred)
	if err := _Blocktracker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
