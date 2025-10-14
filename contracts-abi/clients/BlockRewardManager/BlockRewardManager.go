// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockrewardmanager

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

// BlockrewardmanagerMetaData contains all meta data concerning the Blockrewardmanager contract.
var BlockrewardmanagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardsPctBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payProposer\",\"inputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardsPctBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRewardsPctBps\",\"inputs\":[{\"name\":\"rewardsPctBps\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTreasury\",\"inputs\":[{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"toTreasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdrawToTreasury\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposerPaid\",\"inputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposerAmt\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"rewardAmt\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsPctBpsSet\",\"inputs\":[{\"name\":\"rewardsPctBps\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasurySet\",\"inputs\":[{\"name\":\"treasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryWithdrawn\",\"inputs\":[{\"name\":\"treasuryAmt\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoFundsToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyOwnerOrTreasury\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProposerTransferFailed\",\"inputs\":[{\"name\":\"feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RewardsPctTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TreasuryIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TreasuryTransferFailed\",\"inputs\":[{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// BlockrewardmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use BlockrewardmanagerMetaData.ABI instead.
var BlockrewardmanagerABI = BlockrewardmanagerMetaData.ABI

// Blockrewardmanager is an auto generated Go binding around an Ethereum contract.
type Blockrewardmanager struct {
	BlockrewardmanagerCaller     // Read-only binding to the contract
	BlockrewardmanagerTransactor // Write-only binding to the contract
	BlockrewardmanagerFilterer   // Log filterer for contract events
}

// BlockrewardmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BlockrewardmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockrewardmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BlockrewardmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockrewardmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BlockrewardmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BlockrewardmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BlockrewardmanagerSession struct {
	Contract     *Blockrewardmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BlockrewardmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BlockrewardmanagerCallerSession struct {
	Contract *BlockrewardmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// BlockrewardmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BlockrewardmanagerTransactorSession struct {
	Contract     *BlockrewardmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// BlockrewardmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BlockrewardmanagerRaw struct {
	Contract *Blockrewardmanager // Generic contract binding to access the raw methods on
}

// BlockrewardmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BlockrewardmanagerCallerRaw struct {
	Contract *BlockrewardmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// BlockrewardmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BlockrewardmanagerTransactorRaw struct {
	Contract *BlockrewardmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBlockrewardmanager creates a new instance of Blockrewardmanager, bound to a specific deployed contract.
func NewBlockrewardmanager(address common.Address, backend bind.ContractBackend) (*Blockrewardmanager, error) {
	contract, err := bindBlockrewardmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Blockrewardmanager{BlockrewardmanagerCaller: BlockrewardmanagerCaller{contract: contract}, BlockrewardmanagerTransactor: BlockrewardmanagerTransactor{contract: contract}, BlockrewardmanagerFilterer: BlockrewardmanagerFilterer{contract: contract}}, nil
}

// NewBlockrewardmanagerCaller creates a new read-only instance of Blockrewardmanager, bound to a specific deployed contract.
func NewBlockrewardmanagerCaller(address common.Address, caller bind.ContractCaller) (*BlockrewardmanagerCaller, error) {
	contract, err := bindBlockrewardmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerCaller{contract: contract}, nil
}

// NewBlockrewardmanagerTransactor creates a new write-only instance of Blockrewardmanager, bound to a specific deployed contract.
func NewBlockrewardmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*BlockrewardmanagerTransactor, error) {
	contract, err := bindBlockrewardmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerTransactor{contract: contract}, nil
}

// NewBlockrewardmanagerFilterer creates a new log filterer instance of Blockrewardmanager, bound to a specific deployed contract.
func NewBlockrewardmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*BlockrewardmanagerFilterer, error) {
	contract, err := bindBlockrewardmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerFilterer{contract: contract}, nil
}

// bindBlockrewardmanager binds a generic wrapper to an already deployed contract.
func bindBlockrewardmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BlockrewardmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockrewardmanager *BlockrewardmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockrewardmanager.Contract.BlockrewardmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockrewardmanager *BlockrewardmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.BlockrewardmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockrewardmanager *BlockrewardmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.BlockrewardmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Blockrewardmanager *BlockrewardmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Blockrewardmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Blockrewardmanager *BlockrewardmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Blockrewardmanager *BlockrewardmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blockrewardmanager *BlockrewardmanagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blockrewardmanager *BlockrewardmanagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Blockrewardmanager.Contract.UPGRADEINTERFACEVERSION(&_Blockrewardmanager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Blockrewardmanager.Contract.UPGRADEINTERFACEVERSION(&_Blockrewardmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerSession) Owner() (common.Address, error) {
	return _Blockrewardmanager.Contract.Owner(&_Blockrewardmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) Owner() (common.Address, error) {
	return _Blockrewardmanager.Contract.Owner(&_Blockrewardmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerSession) PendingOwner() (common.Address, error) {
	return _Blockrewardmanager.Contract.PendingOwner(&_Blockrewardmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) PendingOwner() (common.Address, error) {
	return _Blockrewardmanager.Contract.PendingOwner(&_Blockrewardmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blockrewardmanager *BlockrewardmanagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blockrewardmanager *BlockrewardmanagerSession) ProxiableUUID() ([32]byte, error) {
	return _Blockrewardmanager.Contract.ProxiableUUID(&_Blockrewardmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Blockrewardmanager.Contract.ProxiableUUID(&_Blockrewardmanager.CallOpts)
}

// RewardsPctBps is a free data retrieval call binding the contract method 0xd8dd8d11.
//
// Solidity: function rewardsPctBps() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerCaller) RewardsPctBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "rewardsPctBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardsPctBps is a free data retrieval call binding the contract method 0xd8dd8d11.
//
// Solidity: function rewardsPctBps() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerSession) RewardsPctBps() (*big.Int, error) {
	return _Blockrewardmanager.Contract.RewardsPctBps(&_Blockrewardmanager.CallOpts)
}

// RewardsPctBps is a free data retrieval call binding the contract method 0xd8dd8d11.
//
// Solidity: function rewardsPctBps() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) RewardsPctBps() (*big.Int, error) {
	return _Blockrewardmanager.Contract.RewardsPctBps(&_Blockrewardmanager.CallOpts)
}

// ToTreasury is a free data retrieval call binding the contract method 0x79900169.
//
// Solidity: function toTreasury() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerCaller) ToTreasury(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "toTreasury")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ToTreasury is a free data retrieval call binding the contract method 0x79900169.
//
// Solidity: function toTreasury() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerSession) ToTreasury() (*big.Int, error) {
	return _Blockrewardmanager.Contract.ToTreasury(&_Blockrewardmanager.CallOpts)
}

// ToTreasury is a free data retrieval call binding the contract method 0x79900169.
//
// Solidity: function toTreasury() view returns(uint256)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) ToTreasury() (*big.Int, error) {
	return _Blockrewardmanager.Contract.ToTreasury(&_Blockrewardmanager.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Blockrewardmanager.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerSession) Treasury() (common.Address, error) {
	return _Blockrewardmanager.Contract.Treasury(&_Blockrewardmanager.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Blockrewardmanager *BlockrewardmanagerCallerSession) Treasury() (common.Address, error) {
	return _Blockrewardmanager.Contract.Treasury(&_Blockrewardmanager.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.AcceptOwnership(&_Blockrewardmanager.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.AcceptOwnership(&_Blockrewardmanager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address initialOwner, uint256 rewardsPctBps, address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) Initialize(opts *bind.TransactOpts, initialOwner common.Address, rewardsPctBps *big.Int, treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "initialize", initialOwner, rewardsPctBps, treasury)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address initialOwner, uint256 rewardsPctBps, address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) Initialize(initialOwner common.Address, rewardsPctBps *big.Int, treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Initialize(&_Blockrewardmanager.TransactOpts, initialOwner, rewardsPctBps, treasury)
}

// Initialize is a paid mutator transaction binding the contract method 0xc350a1b5.
//
// Solidity: function initialize(address initialOwner, uint256 rewardsPctBps, address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) Initialize(initialOwner common.Address, rewardsPctBps *big.Int, treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Initialize(&_Blockrewardmanager.TransactOpts, initialOwner, rewardsPctBps, treasury)
}

// PayProposer is a paid mutator transaction binding the contract method 0x4256053f.
//
// Solidity: function payProposer(address feeRecipient) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) PayProposer(opts *bind.TransactOpts, feeRecipient common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "payProposer", feeRecipient)
}

// PayProposer is a paid mutator transaction binding the contract method 0x4256053f.
//
// Solidity: function payProposer(address feeRecipient) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) PayProposer(feeRecipient common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.PayProposer(&_Blockrewardmanager.TransactOpts, feeRecipient)
}

// PayProposer is a paid mutator transaction binding the contract method 0x4256053f.
//
// Solidity: function payProposer(address feeRecipient) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) PayProposer(feeRecipient common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.PayProposer(&_Blockrewardmanager.TransactOpts, feeRecipient)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.RenounceOwnership(&_Blockrewardmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.RenounceOwnership(&_Blockrewardmanager.TransactOpts)
}

// SetRewardsPctBps is a paid mutator transaction binding the contract method 0x12dfd7d2.
//
// Solidity: function setRewardsPctBps(uint256 rewardsPctBps) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) SetRewardsPctBps(opts *bind.TransactOpts, rewardsPctBps *big.Int) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "setRewardsPctBps", rewardsPctBps)
}

// SetRewardsPctBps is a paid mutator transaction binding the contract method 0x12dfd7d2.
//
// Solidity: function setRewardsPctBps(uint256 rewardsPctBps) returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) SetRewardsPctBps(rewardsPctBps *big.Int) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.SetRewardsPctBps(&_Blockrewardmanager.TransactOpts, rewardsPctBps)
}

// SetRewardsPctBps is a paid mutator transaction binding the contract method 0x12dfd7d2.
//
// Solidity: function setRewardsPctBps(uint256 rewardsPctBps) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) SetRewardsPctBps(rewardsPctBps *big.Int) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.SetRewardsPctBps(&_Blockrewardmanager.TransactOpts, rewardsPctBps)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) SetTreasury(opts *bind.TransactOpts, treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "setTreasury", treasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) SetTreasury(treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.SetTreasury(&_Blockrewardmanager.TransactOpts, treasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address treasury) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) SetTreasury(treasury common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.SetTreasury(&_Blockrewardmanager.TransactOpts, treasury)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.TransferOwnership(&_Blockrewardmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.TransferOwnership(&_Blockrewardmanager.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.UpgradeToAndCall(&_Blockrewardmanager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.UpgradeToAndCall(&_Blockrewardmanager.TransactOpts, newImplementation, data)
}

// WithdrawToTreasury is a paid mutator transaction binding the contract method 0x7e80c186.
//
// Solidity: function withdrawToTreasury() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) WithdrawToTreasury(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.Transact(opts, "withdrawToTreasury")
}

// WithdrawToTreasury is a paid mutator transaction binding the contract method 0x7e80c186.
//
// Solidity: function withdrawToTreasury() returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) WithdrawToTreasury() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.WithdrawToTreasury(&_Blockrewardmanager.TransactOpts)
}

// WithdrawToTreasury is a paid mutator transaction binding the contract method 0x7e80c186.
//
// Solidity: function withdrawToTreasury() returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) WithdrawToTreasury() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.WithdrawToTreasury(&_Blockrewardmanager.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Fallback(&_Blockrewardmanager.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Fallback(&_Blockrewardmanager.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Blockrewardmanager.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerSession) Receive() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Receive(&_Blockrewardmanager.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Blockrewardmanager *BlockrewardmanagerTransactorSession) Receive() (*types.Transaction, error) {
	return _Blockrewardmanager.Contract.Receive(&_Blockrewardmanager.TransactOpts)
}

// BlockrewardmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Blockrewardmanager contract.
type BlockrewardmanagerInitializedIterator struct {
	Event *BlockrewardmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerInitialized)
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
		it.Event = new(BlockrewardmanagerInitialized)
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
func (it *BlockrewardmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerInitialized represents a Initialized event raised by the Blockrewardmanager contract.
type BlockrewardmanagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*BlockrewardmanagerInitializedIterator, error) {

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerInitializedIterator{contract: _Blockrewardmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerInitialized)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseInitialized(log types.Log) (*BlockrewardmanagerInitialized, error) {
	event := new(BlockrewardmanagerInitialized)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Blockrewardmanager contract.
type BlockrewardmanagerOwnershipTransferStartedIterator struct {
	Event *BlockrewardmanagerOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerOwnershipTransferStarted)
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
		it.Event = new(BlockrewardmanagerOwnershipTransferStarted)
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
func (it *BlockrewardmanagerOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Blockrewardmanager contract.
type BlockrewardmanagerOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BlockrewardmanagerOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerOwnershipTransferStartedIterator{contract: _Blockrewardmanager.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerOwnershipTransferStarted)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseOwnershipTransferStarted(log types.Log) (*BlockrewardmanagerOwnershipTransferStarted, error) {
	event := new(BlockrewardmanagerOwnershipTransferStarted)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Blockrewardmanager contract.
type BlockrewardmanagerOwnershipTransferredIterator struct {
	Event *BlockrewardmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerOwnershipTransferred)
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
		it.Event = new(BlockrewardmanagerOwnershipTransferred)
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
func (it *BlockrewardmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Blockrewardmanager contract.
type BlockrewardmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BlockrewardmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerOwnershipTransferredIterator{contract: _Blockrewardmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerOwnershipTransferred)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseOwnershipTransferred(log types.Log) (*BlockrewardmanagerOwnershipTransferred, error) {
	event := new(BlockrewardmanagerOwnershipTransferred)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerProposerPaidIterator is returned from FilterProposerPaid and is used to iterate over the raw logs and unpacked data for ProposerPaid events raised by the Blockrewardmanager contract.
type BlockrewardmanagerProposerPaidIterator struct {
	Event *BlockrewardmanagerProposerPaid // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerProposerPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerProposerPaid)
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
		it.Event = new(BlockrewardmanagerProposerPaid)
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
func (it *BlockrewardmanagerProposerPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerProposerPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerProposerPaid represents a ProposerPaid event raised by the Blockrewardmanager contract.
type BlockrewardmanagerProposerPaid struct {
	FeeRecipient common.Address
	ProposerAmt  *big.Int
	RewardAmt    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProposerPaid is a free log retrieval operation binding the contract event 0x0c8bf8e02773e67b63846688d6d74c4460cba72aa13c61b5e36ce9fe3d034a56.
//
// Solidity: event ProposerPaid(address indexed feeRecipient, uint256 indexed proposerAmt, uint256 indexed rewardAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterProposerPaid(opts *bind.FilterOpts, feeRecipient []common.Address, proposerAmt []*big.Int, rewardAmt []*big.Int) (*BlockrewardmanagerProposerPaidIterator, error) {

	var feeRecipientRule []interface{}
	for _, feeRecipientItem := range feeRecipient {
		feeRecipientRule = append(feeRecipientRule, feeRecipientItem)
	}
	var proposerAmtRule []interface{}
	for _, proposerAmtItem := range proposerAmt {
		proposerAmtRule = append(proposerAmtRule, proposerAmtItem)
	}
	var rewardAmtRule []interface{}
	for _, rewardAmtItem := range rewardAmt {
		rewardAmtRule = append(rewardAmtRule, rewardAmtItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "ProposerPaid", feeRecipientRule, proposerAmtRule, rewardAmtRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerProposerPaidIterator{contract: _Blockrewardmanager.contract, event: "ProposerPaid", logs: logs, sub: sub}, nil
}

// WatchProposerPaid is a free log subscription operation binding the contract event 0x0c8bf8e02773e67b63846688d6d74c4460cba72aa13c61b5e36ce9fe3d034a56.
//
// Solidity: event ProposerPaid(address indexed feeRecipient, uint256 indexed proposerAmt, uint256 indexed rewardAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchProposerPaid(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerProposerPaid, feeRecipient []common.Address, proposerAmt []*big.Int, rewardAmt []*big.Int) (event.Subscription, error) {

	var feeRecipientRule []interface{}
	for _, feeRecipientItem := range feeRecipient {
		feeRecipientRule = append(feeRecipientRule, feeRecipientItem)
	}
	var proposerAmtRule []interface{}
	for _, proposerAmtItem := range proposerAmt {
		proposerAmtRule = append(proposerAmtRule, proposerAmtItem)
	}
	var rewardAmtRule []interface{}
	for _, rewardAmtItem := range rewardAmt {
		rewardAmtRule = append(rewardAmtRule, rewardAmtItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "ProposerPaid", feeRecipientRule, proposerAmtRule, rewardAmtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerProposerPaid)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "ProposerPaid", log); err != nil {
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

// ParseProposerPaid is a log parse operation binding the contract event 0x0c8bf8e02773e67b63846688d6d74c4460cba72aa13c61b5e36ce9fe3d034a56.
//
// Solidity: event ProposerPaid(address indexed feeRecipient, uint256 indexed proposerAmt, uint256 indexed rewardAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseProposerPaid(log types.Log) (*BlockrewardmanagerProposerPaid, error) {
	event := new(BlockrewardmanagerProposerPaid)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "ProposerPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerRewardsPctBpsSetIterator is returned from FilterRewardsPctBpsSet and is used to iterate over the raw logs and unpacked data for RewardsPctBpsSet events raised by the Blockrewardmanager contract.
type BlockrewardmanagerRewardsPctBpsSetIterator struct {
	Event *BlockrewardmanagerRewardsPctBpsSet // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerRewardsPctBpsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerRewardsPctBpsSet)
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
		it.Event = new(BlockrewardmanagerRewardsPctBpsSet)
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
func (it *BlockrewardmanagerRewardsPctBpsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerRewardsPctBpsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerRewardsPctBpsSet represents a RewardsPctBpsSet event raised by the Blockrewardmanager contract.
type BlockrewardmanagerRewardsPctBpsSet struct {
	RewardsPctBps *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRewardsPctBpsSet is a free log retrieval operation binding the contract event 0xecd960b831b304435f2e2028908a05dbb70756ab3aaf9d16927526264f64cd03.
//
// Solidity: event RewardsPctBpsSet(uint256 indexed rewardsPctBps)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterRewardsPctBpsSet(opts *bind.FilterOpts, rewardsPctBps []*big.Int) (*BlockrewardmanagerRewardsPctBpsSetIterator, error) {

	var rewardsPctBpsRule []interface{}
	for _, rewardsPctBpsItem := range rewardsPctBps {
		rewardsPctBpsRule = append(rewardsPctBpsRule, rewardsPctBpsItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "RewardsPctBpsSet", rewardsPctBpsRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerRewardsPctBpsSetIterator{contract: _Blockrewardmanager.contract, event: "RewardsPctBpsSet", logs: logs, sub: sub}, nil
}

// WatchRewardsPctBpsSet is a free log subscription operation binding the contract event 0xecd960b831b304435f2e2028908a05dbb70756ab3aaf9d16927526264f64cd03.
//
// Solidity: event RewardsPctBpsSet(uint256 indexed rewardsPctBps)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchRewardsPctBpsSet(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerRewardsPctBpsSet, rewardsPctBps []*big.Int) (event.Subscription, error) {

	var rewardsPctBpsRule []interface{}
	for _, rewardsPctBpsItem := range rewardsPctBps {
		rewardsPctBpsRule = append(rewardsPctBpsRule, rewardsPctBpsItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "RewardsPctBpsSet", rewardsPctBpsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerRewardsPctBpsSet)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "RewardsPctBpsSet", log); err != nil {
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

// ParseRewardsPctBpsSet is a log parse operation binding the contract event 0xecd960b831b304435f2e2028908a05dbb70756ab3aaf9d16927526264f64cd03.
//
// Solidity: event RewardsPctBpsSet(uint256 indexed rewardsPctBps)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseRewardsPctBpsSet(log types.Log) (*BlockrewardmanagerRewardsPctBpsSet, error) {
	event := new(BlockrewardmanagerRewardsPctBpsSet)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "RewardsPctBpsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerTreasurySetIterator is returned from FilterTreasurySet and is used to iterate over the raw logs and unpacked data for TreasurySet events raised by the Blockrewardmanager contract.
type BlockrewardmanagerTreasurySetIterator struct {
	Event *BlockrewardmanagerTreasurySet // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerTreasurySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerTreasurySet)
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
		it.Event = new(BlockrewardmanagerTreasurySet)
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
func (it *BlockrewardmanagerTreasurySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerTreasurySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerTreasurySet represents a TreasurySet event raised by the Blockrewardmanager contract.
type BlockrewardmanagerTreasurySet struct {
	Treasury common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTreasurySet is a free log retrieval operation binding the contract event 0x3c864541ef71378c6229510ed90f376565ee42d9c5e0904a984a9e863e6db44f.
//
// Solidity: event TreasurySet(address indexed treasury)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterTreasurySet(opts *bind.FilterOpts, treasury []common.Address) (*BlockrewardmanagerTreasurySetIterator, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "TreasurySet", treasuryRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerTreasurySetIterator{contract: _Blockrewardmanager.contract, event: "TreasurySet", logs: logs, sub: sub}, nil
}

// WatchTreasurySet is a free log subscription operation binding the contract event 0x3c864541ef71378c6229510ed90f376565ee42d9c5e0904a984a9e863e6db44f.
//
// Solidity: event TreasurySet(address indexed treasury)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchTreasurySet(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerTreasurySet, treasury []common.Address) (event.Subscription, error) {

	var treasuryRule []interface{}
	for _, treasuryItem := range treasury {
		treasuryRule = append(treasuryRule, treasuryItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "TreasurySet", treasuryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerTreasurySet)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "TreasurySet", log); err != nil {
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

// ParseTreasurySet is a log parse operation binding the contract event 0x3c864541ef71378c6229510ed90f376565ee42d9c5e0904a984a9e863e6db44f.
//
// Solidity: event TreasurySet(address indexed treasury)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseTreasurySet(log types.Log) (*BlockrewardmanagerTreasurySet, error) {
	event := new(BlockrewardmanagerTreasurySet)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "TreasurySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerTreasuryWithdrawnIterator is returned from FilterTreasuryWithdrawn and is used to iterate over the raw logs and unpacked data for TreasuryWithdrawn events raised by the Blockrewardmanager contract.
type BlockrewardmanagerTreasuryWithdrawnIterator struct {
	Event *BlockrewardmanagerTreasuryWithdrawn // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerTreasuryWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerTreasuryWithdrawn)
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
		it.Event = new(BlockrewardmanagerTreasuryWithdrawn)
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
func (it *BlockrewardmanagerTreasuryWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerTreasuryWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerTreasuryWithdrawn represents a TreasuryWithdrawn event raised by the Blockrewardmanager contract.
type BlockrewardmanagerTreasuryWithdrawn struct {
	TreasuryAmt *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryWithdrawn is a free log retrieval operation binding the contract event 0xdcfb70a6f0f5eab41644ac0cde62fe5f51ce0bb0a53b88ea72c4b2b78ad887bc.
//
// Solidity: event TreasuryWithdrawn(uint256 indexed treasuryAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterTreasuryWithdrawn(opts *bind.FilterOpts, treasuryAmt []*big.Int) (*BlockrewardmanagerTreasuryWithdrawnIterator, error) {

	var treasuryAmtRule []interface{}
	for _, treasuryAmtItem := range treasuryAmt {
		treasuryAmtRule = append(treasuryAmtRule, treasuryAmtItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "TreasuryWithdrawn", treasuryAmtRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerTreasuryWithdrawnIterator{contract: _Blockrewardmanager.contract, event: "TreasuryWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTreasuryWithdrawn is a free log subscription operation binding the contract event 0xdcfb70a6f0f5eab41644ac0cde62fe5f51ce0bb0a53b88ea72c4b2b78ad887bc.
//
// Solidity: event TreasuryWithdrawn(uint256 indexed treasuryAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchTreasuryWithdrawn(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerTreasuryWithdrawn, treasuryAmt []*big.Int) (event.Subscription, error) {

	var treasuryAmtRule []interface{}
	for _, treasuryAmtItem := range treasuryAmt {
		treasuryAmtRule = append(treasuryAmtRule, treasuryAmtItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "TreasuryWithdrawn", treasuryAmtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerTreasuryWithdrawn)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "TreasuryWithdrawn", log); err != nil {
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

// ParseTreasuryWithdrawn is a log parse operation binding the contract event 0xdcfb70a6f0f5eab41644ac0cde62fe5f51ce0bb0a53b88ea72c4b2b78ad887bc.
//
// Solidity: event TreasuryWithdrawn(uint256 indexed treasuryAmt)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseTreasuryWithdrawn(log types.Log) (*BlockrewardmanagerTreasuryWithdrawn, error) {
	event := new(BlockrewardmanagerTreasuryWithdrawn)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "TreasuryWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BlockrewardmanagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Blockrewardmanager contract.
type BlockrewardmanagerUpgradedIterator struct {
	Event *BlockrewardmanagerUpgraded // Event containing the contract specifics and raw log

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
func (it *BlockrewardmanagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BlockrewardmanagerUpgraded)
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
		it.Event = new(BlockrewardmanagerUpgraded)
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
func (it *BlockrewardmanagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BlockrewardmanagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BlockrewardmanagerUpgraded represents a Upgraded event raised by the Blockrewardmanager contract.
type BlockrewardmanagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BlockrewardmanagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BlockrewardmanagerUpgradedIterator{contract: _Blockrewardmanager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Blockrewardmanager *BlockrewardmanagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BlockrewardmanagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Blockrewardmanager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BlockrewardmanagerUpgraded)
				if err := _Blockrewardmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Blockrewardmanager *BlockrewardmanagerFilterer) ParseUpgraded(log types.Log) (*BlockrewardmanagerUpgraded, error) {
	event := new(BlockrewardmanagerUpgraded)
	if err := _Blockrewardmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
