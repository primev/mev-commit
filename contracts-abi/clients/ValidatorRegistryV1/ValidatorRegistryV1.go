// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatorregistryv1

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

// Validatorregistryv1MetaData contains all meta data concerning the Validatorregistryv1 contract.
var Validatorregistryv1MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegateStake\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"stakeOriginator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getBlocksTillWithdrawAllowed\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumberOfStakedValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedValidators\",\"inputs\":[{\"name\":\"start\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnstakingAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_unstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isStaked\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOriginators\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakedValsetVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeBlockNums\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unstakePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unstakingBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeWithdrawn\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// Validatorregistryv1ABI is the input ABI used to generate the binding from.
// Deprecated: Use Validatorregistryv1MetaData.ABI instead.
var Validatorregistryv1ABI = Validatorregistryv1MetaData.ABI

// Validatorregistryv1 is an auto generated Go binding around an Ethereum contract.
type Validatorregistryv1 struct {
	Validatorregistryv1Caller     // Read-only binding to the contract
	Validatorregistryv1Transactor // Write-only binding to the contract
	Validatorregistryv1Filterer   // Log filterer for contract events
}

// Validatorregistryv1Caller is an auto generated read-only Go binding around an Ethereum contract.
type Validatorregistryv1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Validatorregistryv1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Validatorregistryv1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Validatorregistryv1Session struct {
	Contract     *Validatorregistryv1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Validatorregistryv1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Validatorregistryv1CallerSession struct {
	Contract *Validatorregistryv1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// Validatorregistryv1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Validatorregistryv1TransactorSession struct {
	Contract     *Validatorregistryv1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// Validatorregistryv1Raw is an auto generated low-level Go binding around an Ethereum contract.
type Validatorregistryv1Raw struct {
	Contract *Validatorregistryv1 // Generic contract binding to access the raw methods on
}

// Validatorregistryv1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Validatorregistryv1CallerRaw struct {
	Contract *Validatorregistryv1Caller // Generic read-only contract binding to access the raw methods on
}

// Validatorregistryv1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Validatorregistryv1TransactorRaw struct {
	Contract *Validatorregistryv1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorregistryv1 creates a new instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1(address common.Address, backend bind.ContractBackend) (*Validatorregistryv1, error) {
	contract, err := bindValidatorregistryv1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1{Validatorregistryv1Caller: Validatorregistryv1Caller{contract: contract}, Validatorregistryv1Transactor: Validatorregistryv1Transactor{contract: contract}, Validatorregistryv1Filterer: Validatorregistryv1Filterer{contract: contract}}, nil
}

// NewValidatorregistryv1Caller creates a new read-only instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Caller(address common.Address, caller bind.ContractCaller) (*Validatorregistryv1Caller, error) {
	contract, err := bindValidatorregistryv1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Caller{contract: contract}, nil
}

// NewValidatorregistryv1Transactor creates a new write-only instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Transactor(address common.Address, transactor bind.ContractTransactor) (*Validatorregistryv1Transactor, error) {
	contract, err := bindValidatorregistryv1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Transactor{contract: contract}, nil
}

// NewValidatorregistryv1Filterer creates a new log filterer instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Filterer(address common.Address, filterer bind.ContractFilterer) (*Validatorregistryv1Filterer, error) {
	contract, err := bindValidatorregistryv1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Filterer{contract: contract}, nil
}

// bindValidatorregistryv1 binds a generic wrapper to an already deployed contract.
func bindValidatorregistryv1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Validatorregistryv1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistryv1.Contract.Validatorregistryv1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Validatorregistryv1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Validatorregistryv1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistryv1 *Validatorregistryv1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistryv1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistryv1 *Validatorregistryv1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistryv1 *Validatorregistryv1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatorregistryv1.Contract.UPGRADEINTERFACEVERSION(&_Validatorregistryv1.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatorregistryv1.Contract.UPGRADEINTERFACEVERSION(&_Validatorregistryv1.CallOpts)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetBlocksTillWithdrawAllowed(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getBlocksTillWithdrawAllowed", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetNumberOfStakedValidators is a free data retrieval call binding the contract method 0x07258504.
//
// Solidity: function getNumberOfStakedValidators() view returns(uint256, uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetNumberOfStakedValidators(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getNumberOfStakedValidators")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetNumberOfStakedValidators is a free data retrieval call binding the contract method 0x07258504.
//
// Solidity: function getNumberOfStakedValidators() view returns(uint256, uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetNumberOfStakedValidators() (*big.Int, *big.Int, error) {
	return _Validatorregistryv1.Contract.GetNumberOfStakedValidators(&_Validatorregistryv1.CallOpts)
}

// GetNumberOfStakedValidators is a free data retrieval call binding the contract method 0x07258504.
//
// Solidity: function getNumberOfStakedValidators() view returns(uint256, uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetNumberOfStakedValidators() (*big.Int, *big.Int, error) {
	return _Validatorregistryv1.Contract.GetNumberOfStakedValidators(&_Validatorregistryv1.CallOpts)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetStakedAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getStakedAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedValidators is a free data retrieval call binding the contract method 0x7d377923.
//
// Solidity: function getStakedValidators(uint256 start, uint256 end) view returns(bytes[], uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetStakedValidators(opts *bind.CallOpts, start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getStakedValidators", start, end)

	if err != nil {
		return *new([][]byte), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetStakedValidators is a free data retrieval call binding the contract method 0x7d377923.
//
// Solidity: function getStakedValidators(uint256 start, uint256 end) view returns(bytes[], uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetStakedValidators(start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedValidators(&_Validatorregistryv1.CallOpts, start, end)
}

// GetStakedValidators is a free data retrieval call binding the contract method 0x7d377923.
//
// Solidity: function getStakedValidators(uint256 start, uint256 end) view returns(bytes[], uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetStakedValidators(start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedValidators(&_Validatorregistryv1.CallOpts, start, end)
}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetUnstakingAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getUnstakingAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetUnstakingAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetUnstakingAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetUnstakingAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetUnstakingAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Caller) IsStaked(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "isStaked", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Session) IsStaked(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsStaked(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) IsStaked(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsStaked(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) MinStake() (*big.Int, error) {
	return _Validatorregistryv1.Contract.MinStake(&_Validatorregistryv1.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) MinStake() (*big.Int, error) {
	return _Validatorregistryv1.Contract.MinStake(&_Validatorregistryv1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Session) Owner() (common.Address, error) {
	return _Validatorregistryv1.Contract.Owner(&_Validatorregistryv1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) Owner() (common.Address, error) {
	return _Validatorregistryv1.Contract.Owner(&_Validatorregistryv1.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1Session) ProxiableUUID() ([32]byte, error) {
	return _Validatorregistryv1.Contract.ProxiableUUID(&_Validatorregistryv1.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) ProxiableUUID() ([32]byte, error) {
	return _Validatorregistryv1.Contract.ProxiableUUID(&_Validatorregistryv1.CallOpts)
}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Caller) StakeOriginators(opts *bind.CallOpts, arg0 []byte) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "stakeOriginators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Session) StakeOriginators(arg0 []byte) (common.Address, error) {
	return _Validatorregistryv1.Contract.StakeOriginators(&_Validatorregistryv1.CallOpts, arg0)
}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) StakeOriginators(arg0 []byte) (common.Address, error) {
	return _Validatorregistryv1.Contract.StakeOriginators(&_Validatorregistryv1.CallOpts, arg0)
}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) StakedValsetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "stakedValsetVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) StakedValsetVersion() (*big.Int, error) {
	return _Validatorregistryv1.Contract.StakedValsetVersion(&_Validatorregistryv1.CallOpts)
}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) StakedValsetVersion() (*big.Int, error) {
	return _Validatorregistryv1.Contract.StakedValsetVersion(&_Validatorregistryv1.CallOpts)
}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UnstakeBlockNums(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "unstakeBlockNums", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) UnstakeBlockNums(arg0 []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakeBlockNums(&_Validatorregistryv1.CallOpts, arg0)
}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UnstakeBlockNums(arg0 []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakeBlockNums(&_Validatorregistryv1.CallOpts, arg0)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UnstakePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "unstakePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakePeriodBlocks(&_Validatorregistryv1.CallOpts)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakePeriodBlocks(&_Validatorregistryv1.CallOpts)
}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UnstakingBalances(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "unstakingBalances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) UnstakingBalances(arg0 []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakingBalances(&_Validatorregistryv1.CallOpts, arg0)
}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UnstakingBalances(arg0 []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakingBalances(&_Validatorregistryv1.CallOpts, arg0)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] valBLSPubKeys, address stakeOriginator) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) DelegateStake(opts *bind.TransactOpts, valBLSPubKeys [][]byte, stakeOriginator common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "delegateStake", valBLSPubKeys, stakeOriginator)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] valBLSPubKeys, address stakeOriginator) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) DelegateStake(valBLSPubKeys [][]byte, stakeOriginator common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.DelegateStake(&_Validatorregistryv1.TransactOpts, valBLSPubKeys, stakeOriginator)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] valBLSPubKeys, address stakeOriginator) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) DelegateStake(valBLSPubKeys [][]byte, stakeOriginator common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.DelegateStake(&_Validatorregistryv1.TransactOpts, valBLSPubKeys, stakeOriginator)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "initialize", _minStake, _unstakePeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Initialize(_minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Initialize(&_Validatorregistryv1.TransactOpts, _minStake, _unstakePeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Initialize(_minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Initialize(&_Validatorregistryv1.TransactOpts, _minStake, _unstakePeriodBlocks, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.RenounceOwnership(&_Validatorregistryv1.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.RenounceOwnership(&_Validatorregistryv1.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Stake(opts *bind.TransactOpts, valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "stake", valBLSPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Stake(valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Stake(&_Validatorregistryv1.TransactOpts, valBLSPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Stake(valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Stake(&_Validatorregistryv1.TransactOpts, valBLSPubKeys)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.TransferOwnership(&_Validatorregistryv1.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.TransferOwnership(&_Validatorregistryv1.TransactOpts, newOwner)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Unstake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "unstake", blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unstake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unstake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.UpgradeToAndCall(&_Validatorregistryv1.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.UpgradeToAndCall(&_Validatorregistryv1.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Withdraw(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "withdraw", blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Withdraw(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Withdraw(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Fallback(&_Validatorregistryv1.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Fallback(&_Validatorregistryv1.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Receive() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Receive(&_Validatorregistryv1.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Receive() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Receive(&_Validatorregistryv1.TransactOpts)
}

// Validatorregistryv1InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validatorregistryv1 contract.
type Validatorregistryv1InitializedIterator struct {
	Event *Validatorregistryv1Initialized // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Initialized)
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
		it.Event = new(Validatorregistryv1Initialized)
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
func (it *Validatorregistryv1InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Initialized represents a Initialized event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterInitialized(opts *bind.FilterOpts) (*Validatorregistryv1InitializedIterator, error) {

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1InitializedIterator{contract: _Validatorregistryv1.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Initialized) (event.Subscription, error) {

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Initialized)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseInitialized(log types.Log) (*Validatorregistryv1Initialized, error) {
	event := new(Validatorregistryv1Initialized)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferredIterator struct {
	Event *Validatorregistryv1OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1OwnershipTransferred)
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
		it.Event = new(Validatorregistryv1OwnershipTransferred)
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
func (it *Validatorregistryv1OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1OwnershipTransferred represents a OwnershipTransferred event raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Validatorregistryv1OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1OwnershipTransferredIterator{contract: _Validatorregistryv1.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1OwnershipTransferred)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseOwnershipTransferred(log types.Log) (*Validatorregistryv1OwnershipTransferred, error) {
	event := new(Validatorregistryv1OwnershipTransferred)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1StakeWithdrawnIterator is returned from FilterStakeWithdrawn and is used to iterate over the raw logs and unpacked data for StakeWithdrawn events raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeWithdrawnIterator struct {
	Event *Validatorregistryv1StakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1StakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1StakeWithdrawn)
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
		it.Event = new(Validatorregistryv1StakeWithdrawn)
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
func (it *Validatorregistryv1StakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1StakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1StakeWithdrawn represents a StakeWithdrawn event raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeWithdrawn struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdrawn is a free log retrieval operation binding the contract event 0x60e5b516955702ed9d33ecfc7eaaf6b2f5cea6bd67820e5e4f0096eed587c29b.
//
// Solidity: event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterStakeWithdrawn(opts *bind.FilterOpts, txOriginator []common.Address) (*Validatorregistryv1StakeWithdrawnIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "StakeWithdrawn", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1StakeWithdrawnIterator{contract: _Validatorregistryv1.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchStakeWithdrawn is a free log subscription operation binding the contract event 0x60e5b516955702ed9d33ecfc7eaaf6b2f5cea6bd67820e5e4f0096eed587c29b.
//
// Solidity: event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1StakeWithdrawn, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "StakeWithdrawn", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1StakeWithdrawn)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
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

// ParseStakeWithdrawn is a log parse operation binding the contract event 0x60e5b516955702ed9d33ecfc7eaaf6b2f5cea6bd67820e5e4f0096eed587c29b.
//
// Solidity: event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseStakeWithdrawn(log types.Log) (*Validatorregistryv1StakeWithdrawn, error) {
	event := new(Validatorregistryv1StakeWithdrawn)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1StakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakedIterator struct {
	Event *Validatorregistryv1Staked // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1StakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Staked)
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
		it.Event = new(Validatorregistryv1Staked)
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
func (it *Validatorregistryv1StakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1StakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Staked represents a Staked event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Staked struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xf4679d394f1f97f1a3df1d73e193866ec5a813168ad5fa6958f9be21b10a594e.
//
// Solidity: event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterStaked(opts *bind.FilterOpts, txOriginator []common.Address) (*Validatorregistryv1StakedIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Staked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1StakedIterator{contract: _Validatorregistryv1.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xf4679d394f1f97f1a3df1d73e193866ec5a813168ad5fa6958f9be21b10a594e.
//
// Solidity: event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Staked, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Staked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Staked)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0xf4679d394f1f97f1a3df1d73e193866ec5a813168ad5fa6958f9be21b10a594e.
//
// Solidity: event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseStaked(log types.Log) (*Validatorregistryv1Staked, error) {
	event := new(Validatorregistryv1Staked)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UnstakedIterator struct {
	Event *Validatorregistryv1Unstaked // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Unstaked)
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
		it.Event = new(Validatorregistryv1Unstaked)
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
func (it *Validatorregistryv1UnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Unstaked represents a Unstaked event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Unstaked struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x34f6c826286e3bec36208d054dcb1ad3c49725237a7644e1a6d157a92ae7a3e1.
//
// Solidity: event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUnstaked(opts *bind.FilterOpts, txOriginator []common.Address) (*Validatorregistryv1UnstakedIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Unstaked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UnstakedIterator{contract: _Validatorregistryv1.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x34f6c826286e3bec36208d054dcb1ad3c49725237a7644e1a6d157a92ae7a3e1.
//
// Solidity: event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Unstaked, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Unstaked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Unstaked)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x34f6c826286e3bec36208d054dcb1ad3c49725237a7644e1a6d157a92ae7a3e1.
//
// Solidity: event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUnstaked(log types.Log) (*Validatorregistryv1Unstaked, error) {
	event := new(Validatorregistryv1Unstaked)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UpgradedIterator struct {
	Event *Validatorregistryv1Upgraded // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Upgraded)
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
		it.Event = new(Validatorregistryv1Upgraded)
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
func (it *Validatorregistryv1UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Upgraded represents a Upgraded event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*Validatorregistryv1UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UpgradedIterator{contract: _Validatorregistryv1.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Upgraded)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUpgraded(log types.Log) (*Validatorregistryv1Upgraded, error) {
	event := new(Validatorregistryv1Upgraded)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
