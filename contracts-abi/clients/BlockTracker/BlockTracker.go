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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addBuilderAddress\",\"inputs\":[{\"name\":\"builderName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"builderAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blockBuilderNameToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockWinners\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlockWinner\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuilder\",\"inputs\":[{\"name\":\"builderNameGraffiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrentWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"blocksPerWindow_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"oracleAccount_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"oracleAccount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordL1Block\",\"inputs\":[{\"name\":\"_blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_winnerGraffiti\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOracleAccount\",\"inputs\":[{\"name\":\"newOracleAccount\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewL1Block\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"winner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NewWindow\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleAccountSet\",\"inputs\":[{\"name\":\"oldOracleAccount\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOracleAccount\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BlockNumberIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOracleAccount\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleAccount\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
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

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blocktracker *BlocktrackerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blocktracker *BlocktrackerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Blocktracker.Contract.UPGRADEINTERFACEVERSION(&_Blocktracker.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blocktracker *BlocktrackerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Blocktracker.Contract.UPGRADEINTERFACEVERSION(&_Blocktracker.CallOpts)
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
// Solidity: function getBuilder(string builderNameGraffiti) view returns(address)
func (_Blocktracker *BlocktrackerCaller) GetBuilder(opts *bind.CallOpts, builderNameGraffiti string) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "getBuilder", builderNameGraffiti)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGraffiti) view returns(address)
func (_Blocktracker *BlocktrackerSession) GetBuilder(builderNameGraffiti string) (common.Address, error) {
	return _Blocktracker.Contract.GetBuilder(&_Blocktracker.CallOpts, builderNameGraffiti)
}

// GetBuilder is a free data retrieval call binding the contract method 0x237ba8fb.
//
// Solidity: function getBuilder(string builderNameGraffiti) view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) GetBuilder(builderNameGraffiti string) (common.Address, error) {
	return _Blocktracker.Contract.GetBuilder(&_Blocktracker.CallOpts, builderNameGraffiti)
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

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Blocktracker *BlocktrackerCaller) OracleAccount(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "oracleAccount")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Blocktracker *BlocktrackerSession) OracleAccount() (common.Address, error) {
	return _Blocktracker.Contract.OracleAccount(&_Blocktracker.CallOpts)
}

// OracleAccount is a free data retrieval call binding the contract method 0xe7c59736.
//
// Solidity: function oracleAccount() view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) OracleAccount() (common.Address, error) {
	return _Blocktracker.Contract.OracleAccount(&_Blocktracker.CallOpts)
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

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Blocktracker *BlocktrackerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Blocktracker *BlocktrackerSession) Paused() (bool, error) {
	return _Blocktracker.Contract.Paused(&_Blocktracker.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Blocktracker *BlocktrackerCallerSession) Paused() (bool, error) {
	return _Blocktracker.Contract.Paused(&_Blocktracker.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blocktracker *BlocktrackerCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blocktracker *BlocktrackerSession) PendingOwner() (common.Address, error) {
	return _Blocktracker.Contract.PendingOwner(&_Blocktracker.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blocktracker *BlocktrackerCallerSession) PendingOwner() (common.Address, error) {
	return _Blocktracker.Contract.PendingOwner(&_Blocktracker.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blocktracker *BlocktrackerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Blocktracker.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blocktracker *BlocktrackerSession) ProxiableUUID() ([32]byte, error) {
	return _Blocktracker.Contract.ProxiableUUID(&_Blocktracker.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blocktracker *BlocktrackerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Blocktracker.Contract.ProxiableUUID(&_Blocktracker.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blocktracker *BlocktrackerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blocktracker *BlocktrackerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Blocktracker.Contract.AcceptOwnership(&_Blocktracker.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blocktracker *BlocktrackerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Blocktracker.Contract.AcceptOwnership(&_Blocktracker.TransactOpts)
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

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 blocksPerWindow_, address oracleAccount_, address owner_) returns()
func (_Blocktracker *BlocktrackerTransactor) Initialize(opts *bind.TransactOpts, blocksPerWindow_ *big.Int, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "initialize", blocksPerWindow_, oracleAccount_, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 blocksPerWindow_, address oracleAccount_, address owner_) returns()
func (_Blocktracker *BlocktrackerSession) Initialize(blocksPerWindow_ *big.Int, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.Initialize(&_Blocktracker.TransactOpts, blocksPerWindow_, oracleAccount_, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xb4988fd0.
//
// Solidity: function initialize(uint256 blocksPerWindow_, address oracleAccount_, address owner_) returns()
func (_Blocktracker *BlocktrackerTransactorSession) Initialize(blocksPerWindow_ *big.Int, oracleAccount_ common.Address, owner_ common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.Initialize(&_Blocktracker.TransactOpts, blocksPerWindow_, oracleAccount_, owner_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Blocktracker *BlocktrackerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Blocktracker *BlocktrackerSession) Pause() (*types.Transaction, error) {
	return _Blocktracker.Contract.Pause(&_Blocktracker.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Blocktracker *BlocktrackerTransactorSession) Pause() (*types.Transaction, error) {
	return _Blocktracker.Contract.Pause(&_Blocktracker.TransactOpts)
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

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Blocktracker *BlocktrackerTransactor) SetOracleAccount(opts *bind.TransactOpts, newOracleAccount common.Address) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "setOracleAccount", newOracleAccount)
}

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Blocktracker *BlocktrackerSession) SetOracleAccount(newOracleAccount common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.SetOracleAccount(&_Blocktracker.TransactOpts, newOracleAccount)
}

// SetOracleAccount is a paid mutator transaction binding the contract method 0x58b20365.
//
// Solidity: function setOracleAccount(address newOracleAccount) returns()
func (_Blocktracker *BlocktrackerTransactorSession) SetOracleAccount(newOracleAccount common.Address) (*types.Transaction, error) {
	return _Blocktracker.Contract.SetOracleAccount(&_Blocktracker.TransactOpts, newOracleAccount)
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

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Blocktracker *BlocktrackerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Blocktracker *BlocktrackerSession) Unpause() (*types.Transaction, error) {
	return _Blocktracker.Contract.Unpause(&_Blocktracker.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Blocktracker *BlocktrackerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Blocktracker.Contract.Unpause(&_Blocktracker.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blocktracker *BlocktrackerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blocktracker.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blocktracker *BlocktrackerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blocktracker.Contract.UpgradeToAndCall(&_Blocktracker.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blocktracker *BlocktrackerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blocktracker.Contract.UpgradeToAndCall(&_Blocktracker.TransactOpts, newImplementation, data)
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

// BlocktrackerOracleAccountSetIterator is returned from FilterOracleAccountSet and is used to iterate over the raw logs and unpacked data for OracleAccountSet events raised by the Blocktracker contract.
type BlocktrackerOracleAccountSetIterator struct {
	Event *BlocktrackerOracleAccountSet // Event containing the contract specifics and raw log

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
func (it *BlocktrackerOracleAccountSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerOracleAccountSet)
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
		it.Event = new(BlocktrackerOracleAccountSet)
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
func (it *BlocktrackerOracleAccountSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerOracleAccountSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerOracleAccountSet represents a OracleAccountSet event raised by the Blocktracker contract.
type BlocktrackerOracleAccountSet struct {
	OldOracleAccount common.Address
	NewOracleAccount common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOracleAccountSet is a free log retrieval operation binding the contract event 0xc44093d4ba5b256ab49bc7bebfea8c049eb2f75fff6bcb9a8da6f8b1c92249e9.
//
// Solidity: event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount)
func (_Blocktracker *BlocktrackerFilterer) FilterOracleAccountSet(opts *bind.FilterOpts, oldOracleAccount []common.Address, newOracleAccount []common.Address) (*BlocktrackerOracleAccountSetIterator, error) {

	var oldOracleAccountRule []interface{}
	for _, oldOracleAccountItem := range oldOracleAccount {
		oldOracleAccountRule = append(oldOracleAccountRule, oldOracleAccountItem)
	}
	var newOracleAccountRule []interface{}
	for _, newOracleAccountItem := range newOracleAccount {
		newOracleAccountRule = append(newOracleAccountRule, newOracleAccountItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "OracleAccountSet", oldOracleAccountRule, newOracleAccountRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerOracleAccountSetIterator{contract: _Blocktracker.contract, event: "OracleAccountSet", logs: logs, sub: sub}, nil
}

// WatchOracleAccountSet is a free log subscription operation binding the contract event 0xc44093d4ba5b256ab49bc7bebfea8c049eb2f75fff6bcb9a8da6f8b1c92249e9.
//
// Solidity: event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount)
func (_Blocktracker *BlocktrackerFilterer) WatchOracleAccountSet(opts *bind.WatchOpts, sink chan<- *BlocktrackerOracleAccountSet, oldOracleAccount []common.Address, newOracleAccount []common.Address) (event.Subscription, error) {

	var oldOracleAccountRule []interface{}
	for _, oldOracleAccountItem := range oldOracleAccount {
		oldOracleAccountRule = append(oldOracleAccountRule, oldOracleAccountItem)
	}
	var newOracleAccountRule []interface{}
	for _, newOracleAccountItem := range newOracleAccount {
		newOracleAccountRule = append(newOracleAccountRule, newOracleAccountItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "OracleAccountSet", oldOracleAccountRule, newOracleAccountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerOracleAccountSet)
				if err := _Blocktracker.contract.UnpackLog(event, "OracleAccountSet", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseOracleAccountSet(log types.Log) (*BlocktrackerOracleAccountSet, error) {
	event := new(BlocktrackerOracleAccountSet)
	if err := _Blocktracker.contract.UnpackLog(event, "OracleAccountSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Blocktracker contract.
type BlocktrackerOwnershipTransferStartedIterator struct {
	Event *BlocktrackerOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *BlocktrackerOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerOwnershipTransferStarted)
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
		it.Event = new(BlocktrackerOwnershipTransferStarted)
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
func (it *BlocktrackerOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Blocktracker contract.
type BlocktrackerOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Blocktracker *BlocktrackerFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BlocktrackerOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerOwnershipTransferStartedIterator{contract: _Blocktracker.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Blocktracker *BlocktrackerFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *BlocktrackerOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerOwnershipTransferStarted)
				if err := _Blocktracker.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseOwnershipTransferStarted(log types.Log) (*BlocktrackerOwnershipTransferStarted, error) {
	event := new(BlocktrackerOwnershipTransferStarted)
	if err := _Blocktracker.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// BlocktrackerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Blocktracker contract.
type BlocktrackerPausedIterator struct {
	Event *BlocktrackerPaused // Event containing the contract specifics and raw log

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
func (it *BlocktrackerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerPaused)
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
		it.Event = new(BlocktrackerPaused)
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
func (it *BlocktrackerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerPaused represents a Paused event raised by the Blocktracker contract.
type BlocktrackerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Blocktracker *BlocktrackerFilterer) FilterPaused(opts *bind.FilterOpts) (*BlocktrackerPausedIterator, error) {

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &BlocktrackerPausedIterator{contract: _Blocktracker.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Blocktracker *BlocktrackerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *BlocktrackerPaused) (event.Subscription, error) {

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerPaused)
				if err := _Blocktracker.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParsePaused(log types.Log) (*BlocktrackerPaused, error) {
	event := new(BlocktrackerPaused)
	if err := _Blocktracker.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Blocktracker contract.
type BlocktrackerUnpausedIterator struct {
	Event *BlocktrackerUnpaused // Event containing the contract specifics and raw log

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
func (it *BlocktrackerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerUnpaused)
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
		it.Event = new(BlocktrackerUnpaused)
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
func (it *BlocktrackerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerUnpaused represents a Unpaused event raised by the Blocktracker contract.
type BlocktrackerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Blocktracker *BlocktrackerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*BlocktrackerUnpausedIterator, error) {

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &BlocktrackerUnpausedIterator{contract: _Blocktracker.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Blocktracker *BlocktrackerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *BlocktrackerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerUnpaused)
				if err := _Blocktracker.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseUnpaused(log types.Log) (*BlocktrackerUnpaused, error) {
	event := new(BlocktrackerUnpaused)
	if err := _Blocktracker.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlocktrackerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Blocktracker contract.
type BlocktrackerUpgradedIterator struct {
	Event *BlocktrackerUpgraded // Event containing the contract specifics and raw log

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
func (it *BlocktrackerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlocktrackerUpgraded)
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
		it.Event = new(BlocktrackerUpgraded)
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
func (it *BlocktrackerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlocktrackerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlocktrackerUpgraded represents a Upgraded event raised by the Blocktracker contract.
type BlocktrackerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Blocktracker *BlocktrackerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BlocktrackerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Blocktracker.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BlocktrackerUpgradedIterator{contract: _Blocktracker.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Blocktracker *BlocktrackerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BlocktrackerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Blocktracker.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlocktrackerUpgraded)
				if err := _Blocktracker.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Blocktracker *BlocktrackerFilterer) ParseUpgraded(log types.Log) (*BlocktrackerUpgraded, error) {
	event := new(BlocktrackerUpgraded)
	if err := _Blocktracker.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
