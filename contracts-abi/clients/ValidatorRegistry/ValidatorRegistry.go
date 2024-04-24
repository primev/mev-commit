// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatorregistry

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

// ValidatorregistryMetaData contains all meta data concerning the Validatorregistry contract.
var ValidatorregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBlocksTillWithdrawAllowed\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumberOfStakedValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedValidators\",\"inputs\":[{\"name\":\"start\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnstakingAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_unstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isStaked\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"valBLSPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakeOriginators\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakedValsetVersion\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakeBlockNums\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unstakePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unstakingBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeWithdrawn\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"txOriginator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
}

// ValidatorregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidatorregistryMetaData.ABI instead.
var ValidatorregistryABI = ValidatorregistryMetaData.ABI

// Validatorregistry is an auto generated Go binding around an Ethereum contract.
type Validatorregistry struct {
	ValidatorregistryCaller     // Read-only binding to the contract
	ValidatorregistryTransactor // Write-only binding to the contract
	ValidatorregistryFilterer   // Log filterer for contract events
}

// ValidatorregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorregistrySession struct {
	Contract     *Validatorregistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ValidatorregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorregistryCallerSession struct {
	Contract *ValidatorregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ValidatorregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorregistryTransactorSession struct {
	Contract     *ValidatorregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ValidatorregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorregistryRaw struct {
	Contract *Validatorregistry // Generic contract binding to access the raw methods on
}

// ValidatorregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorregistryCallerRaw struct {
	Contract *ValidatorregistryCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorregistryTransactorRaw struct {
	Contract *ValidatorregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorregistry creates a new instance of Validatorregistry, bound to a specific deployed contract.
func NewValidatorregistry(address common.Address, backend bind.ContractBackend) (*Validatorregistry, error) {
	contract, err := bindValidatorregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validatorregistry{ValidatorregistryCaller: ValidatorregistryCaller{contract: contract}, ValidatorregistryTransactor: ValidatorregistryTransactor{contract: contract}, ValidatorregistryFilterer: ValidatorregistryFilterer{contract: contract}}, nil
}

// NewValidatorregistryCaller creates a new read-only instance of Validatorregistry, bound to a specific deployed contract.
func NewValidatorregistryCaller(address common.Address, caller bind.ContractCaller) (*ValidatorregistryCaller, error) {
	contract, err := bindValidatorregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryCaller{contract: contract}, nil
}

// NewValidatorregistryTransactor creates a new write-only instance of Validatorregistry, bound to a specific deployed contract.
func NewValidatorregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorregistryTransactor, error) {
	contract, err := bindValidatorregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryTransactor{contract: contract}, nil
}

// NewValidatorregistryFilterer creates a new log filterer instance of Validatorregistry, bound to a specific deployed contract.
func NewValidatorregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorregistryFilterer, error) {
	contract, err := bindValidatorregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryFilterer{contract: contract}, nil
}

// bindValidatorregistry binds a generic wrapper to an already deployed contract.
func bindValidatorregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidatorregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistry *ValidatorregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistry.Contract.ValidatorregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistry *ValidatorregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistry.Contract.ValidatorregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistry *ValidatorregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistry.Contract.ValidatorregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistry *ValidatorregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistry *ValidatorregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistry *ValidatorregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistry.Contract.contract.Transact(opts, method, params...)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) GetBlocksTillWithdrawAllowed(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "getBlocksTillWithdrawAllowed", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// GetNumberOfStakedValidators is a free data retrieval call binding the contract method 0x07258504.
//
// Solidity: function getNumberOfStakedValidators() view returns(uint256, uint256)
func (_Validatorregistry *ValidatorregistryCaller) GetNumberOfStakedValidators(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "getNumberOfStakedValidators")

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
func (_Validatorregistry *ValidatorregistrySession) GetNumberOfStakedValidators() (*big.Int, *big.Int, error) {
	return _Validatorregistry.Contract.GetNumberOfStakedValidators(&_Validatorregistry.CallOpts)
}

// GetNumberOfStakedValidators is a free data retrieval call binding the contract method 0x07258504.
//
// Solidity: function getNumberOfStakedValidators() view returns(uint256, uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) GetNumberOfStakedValidators() (*big.Int, *big.Int, error) {
	return _Validatorregistry.Contract.GetNumberOfStakedValidators(&_Validatorregistry.CallOpts)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) GetStakedAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "getStakedAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetStakedAmount(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetStakedAmount(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// GetStakedValidators is a free data retrieval call binding the contract method 0x7d377923.
//
// Solidity: function getStakedValidators(uint256 start, uint256 end) view returns(bytes[], uint256)
func (_Validatorregistry *ValidatorregistryCaller) GetStakedValidators(opts *bind.CallOpts, start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "getStakedValidators", start, end)

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
func (_Validatorregistry *ValidatorregistrySession) GetStakedValidators(start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	return _Validatorregistry.Contract.GetStakedValidators(&_Validatorregistry.CallOpts, start, end)
}

// GetStakedValidators is a free data retrieval call binding the contract method 0x7d377923.
//
// Solidity: function getStakedValidators(uint256 start, uint256 end) view returns(bytes[], uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) GetStakedValidators(start *big.Int, end *big.Int) ([][]byte, *big.Int, error) {
	return _Validatorregistry.Contract.GetStakedValidators(&_Validatorregistry.CallOpts, start, end)
}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) GetUnstakingAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "getUnstakingAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) GetUnstakingAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetUnstakingAmount(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// GetUnstakingAmount is a free data retrieval call binding the contract method 0xa812e103.
//
// Solidity: function getUnstakingAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) GetUnstakingAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.GetUnstakingAmount(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistry *ValidatorregistryCaller) IsStaked(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "isStaked", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistry *ValidatorregistrySession) IsStaked(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistry.Contract.IsStaked(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// IsStaked is a free data retrieval call binding the contract method 0xcdb513b4.
//
// Solidity: function isStaked(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistry *ValidatorregistryCallerSession) IsStaked(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistry.Contract.IsStaked(&_Validatorregistry.CallOpts, valBLSPubKey)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) MinStake() (*big.Int, error) {
	return _Validatorregistry.Contract.MinStake(&_Validatorregistry.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) MinStake() (*big.Int, error) {
	return _Validatorregistry.Contract.MinStake(&_Validatorregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistry *ValidatorregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistry *ValidatorregistrySession) Owner() (common.Address, error) {
	return _Validatorregistry.Contract.Owner(&_Validatorregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistry *ValidatorregistryCallerSession) Owner() (common.Address, error) {
	return _Validatorregistry.Contract.Owner(&_Validatorregistry.CallOpts)
}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistry *ValidatorregistryCaller) StakeOriginators(opts *bind.CallOpts, arg0 []byte) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "stakeOriginators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistry *ValidatorregistrySession) StakeOriginators(arg0 []byte) (common.Address, error) {
	return _Validatorregistry.Contract.StakeOriginators(&_Validatorregistry.CallOpts, arg0)
}

// StakeOriginators is a free data retrieval call binding the contract method 0x2e5b5fd7.
//
// Solidity: function stakeOriginators(bytes ) view returns(address)
func (_Validatorregistry *ValidatorregistryCallerSession) StakeOriginators(arg0 []byte) (common.Address, error) {
	return _Validatorregistry.Contract.StakeOriginators(&_Validatorregistry.CallOpts, arg0)
}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) StakedValsetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "stakedValsetVersion")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) StakedValsetVersion() (*big.Int, error) {
	return _Validatorregistry.Contract.StakedValsetVersion(&_Validatorregistry.CallOpts)
}

// StakedValsetVersion is a free data retrieval call binding the contract method 0xd628ee62.
//
// Solidity: function stakedValsetVersion() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) StakedValsetVersion() (*big.Int, error) {
	return _Validatorregistry.Contract.StakedValsetVersion(&_Validatorregistry.CallOpts)
}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) UnstakeBlockNums(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "unstakeBlockNums", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) UnstakeBlockNums(arg0 []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakeBlockNums(&_Validatorregistry.CallOpts, arg0)
}

// UnstakeBlockNums is a free data retrieval call binding the contract method 0x2f8836a5.
//
// Solidity: function unstakeBlockNums(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) UnstakeBlockNums(arg0 []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakeBlockNums(&_Validatorregistry.CallOpts, arg0)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) UnstakePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "unstakePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakePeriodBlocks(&_Validatorregistry.CallOpts)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakePeriodBlocks(&_Validatorregistry.CallOpts)
}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCaller) UnstakingBalances(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistry.contract.Call(opts, &out, "unstakingBalances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistrySession) UnstakingBalances(arg0 []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakingBalances(&_Validatorregistry.CallOpts, arg0)
}

// UnstakingBalances is a free data retrieval call binding the contract method 0xfe6c470c.
//
// Solidity: function unstakingBalances(bytes ) view returns(uint256)
func (_Validatorregistry *ValidatorregistryCallerSession) UnstakingBalances(arg0 []byte) (*big.Int, error) {
	return _Validatorregistry.Contract.UnstakingBalances(&_Validatorregistry.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistry *ValidatorregistryTransactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "initialize", _minStake, _unstakePeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistry *ValidatorregistrySession) Initialize(_minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Initialize(&_Validatorregistry.TransactOpts, _minStake, _unstakePeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6ab36f2.
//
// Solidity: function initialize(uint256 _minStake, uint256 _unstakePeriodBlocks, address _owner) returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) Initialize(_minStake *big.Int, _unstakePeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Initialize(&_Validatorregistry.TransactOpts, _minStake, _unstakePeriodBlocks, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistry *ValidatorregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistry *ValidatorregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistry.Contract.RenounceOwnership(&_Validatorregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistry.Contract.RenounceOwnership(&_Validatorregistry.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistry *ValidatorregistryTransactor) Stake(opts *bind.TransactOpts, valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "stake", valBLSPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistry *ValidatorregistrySession) Stake(valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Stake(&_Validatorregistry.TransactOpts, valBLSPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] valBLSPubKeys) payable returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) Stake(valBLSPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Stake(&_Validatorregistry.TransactOpts, valBLSPubKeys)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistry *ValidatorregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistry *ValidatorregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.Contract.TransferOwnership(&_Validatorregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistry.Contract.TransferOwnership(&_Validatorregistry.TransactOpts, newOwner)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistryTransactor) Unstake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "unstake", blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistrySession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Unstake(&_Validatorregistry.TransactOpts, blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Unstake(&_Validatorregistry.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistryTransactor) Withdraw(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.contract.Transact(opts, "withdraw", blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistrySession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Withdraw(&_Validatorregistry.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistry *ValidatorregistryTransactorSession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistry.Contract.Withdraw(&_Validatorregistry.TransactOpts, blsPubKeys)
}

// ValidatorregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validatorregistry contract.
type ValidatorregistryInitializedIterator struct {
	Event *ValidatorregistryInitialized // Event containing the contract specifics and raw log

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
func (it *ValidatorregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorregistryInitialized)
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
		it.Event = new(ValidatorregistryInitialized)
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
func (it *ValidatorregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorregistryInitialized represents a Initialized event raised by the Validatorregistry contract.
type ValidatorregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistry *ValidatorregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ValidatorregistryInitializedIterator, error) {

	logs, sub, err := _Validatorregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryInitializedIterator{contract: _Validatorregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistry *ValidatorregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ValidatorregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Validatorregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorregistryInitialized)
				if err := _Validatorregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validatorregistry *ValidatorregistryFilterer) ParseInitialized(log types.Log) (*ValidatorregistryInitialized, error) {
	event := new(ValidatorregistryInitialized)
	if err := _Validatorregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validatorregistry contract.
type ValidatorregistryOwnershipTransferredIterator struct {
	Event *ValidatorregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ValidatorregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorregistryOwnershipTransferred)
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
		it.Event = new(ValidatorregistryOwnershipTransferred)
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
func (it *ValidatorregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Validatorregistry contract.
type ValidatorregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistry *ValidatorregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ValidatorregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryOwnershipTransferredIterator{contract: _Validatorregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistry *ValidatorregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ValidatorregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorregistryOwnershipTransferred)
				if err := _Validatorregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validatorregistry *ValidatorregistryFilterer) ParseOwnershipTransferred(log types.Log) (*ValidatorregistryOwnershipTransferred, error) {
	event := new(ValidatorregistryOwnershipTransferred)
	if err := _Validatorregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorregistryStakeWithdrawnIterator is returned from FilterStakeWithdrawn and is used to iterate over the raw logs and unpacked data for StakeWithdrawn events raised by the Validatorregistry contract.
type ValidatorregistryStakeWithdrawnIterator struct {
	Event *ValidatorregistryStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *ValidatorregistryStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorregistryStakeWithdrawn)
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
		it.Event = new(ValidatorregistryStakeWithdrawn)
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
func (it *ValidatorregistryStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorregistryStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorregistryStakeWithdrawn represents a StakeWithdrawn event raised by the Validatorregistry contract.
type ValidatorregistryStakeWithdrawn struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdrawn is a free log retrieval operation binding the contract event 0x60e5b516955702ed9d33ecfc7eaaf6b2f5cea6bd67820e5e4f0096eed587c29b.
//
// Solidity: event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) FilterStakeWithdrawn(opts *bind.FilterOpts, txOriginator []common.Address) (*ValidatorregistryStakeWithdrawnIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.FilterLogs(opts, "StakeWithdrawn", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryStakeWithdrawnIterator{contract: _Validatorregistry.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchStakeWithdrawn is a free log subscription operation binding the contract event 0x60e5b516955702ed9d33ecfc7eaaf6b2f5cea6bd67820e5e4f0096eed587c29b.
//
// Solidity: event StakeWithdrawn(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *ValidatorregistryStakeWithdrawn, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.WatchLogs(opts, "StakeWithdrawn", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorregistryStakeWithdrawn)
				if err := _Validatorregistry.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
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
func (_Validatorregistry *ValidatorregistryFilterer) ParseStakeWithdrawn(log types.Log) (*ValidatorregistryStakeWithdrawn, error) {
	event := new(ValidatorregistryStakeWithdrawn)
	if err := _Validatorregistry.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorregistryStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Validatorregistry contract.
type ValidatorregistryStakedIterator struct {
	Event *ValidatorregistryStaked // Event containing the contract specifics and raw log

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
func (it *ValidatorregistryStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorregistryStaked)
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
		it.Event = new(ValidatorregistryStaked)
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
func (it *ValidatorregistryStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorregistryStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorregistryStaked represents a Staked event raised by the Validatorregistry contract.
type ValidatorregistryStaked struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xf4679d394f1f97f1a3df1d73e193866ec5a813168ad5fa6958f9be21b10a594e.
//
// Solidity: event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) FilterStaked(opts *bind.FilterOpts, txOriginator []common.Address) (*ValidatorregistryStakedIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.FilterLogs(opts, "Staked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryStakedIterator{contract: _Validatorregistry.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xf4679d394f1f97f1a3df1d73e193866ec5a813168ad5fa6958f9be21b10a594e.
//
// Solidity: event Staked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *ValidatorregistryStaked, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.WatchLogs(opts, "Staked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorregistryStaked)
				if err := _Validatorregistry.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_Validatorregistry *ValidatorregistryFilterer) ParseStaked(log types.Log) (*ValidatorregistryStaked, error) {
	event := new(ValidatorregistryStaked)
	if err := _Validatorregistry.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorregistryUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Validatorregistry contract.
type ValidatorregistryUnstakedIterator struct {
	Event *ValidatorregistryUnstaked // Event containing the contract specifics and raw log

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
func (it *ValidatorregistryUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorregistryUnstaked)
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
		it.Event = new(ValidatorregistryUnstaked)
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
func (it *ValidatorregistryUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorregistryUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorregistryUnstaked represents a Unstaked event raised by the Validatorregistry contract.
type ValidatorregistryUnstaked struct {
	TxOriginator common.Address
	ValBLSPubKey []byte
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x34f6c826286e3bec36208d054dcb1ad3c49725237a7644e1a6d157a92ae7a3e1.
//
// Solidity: event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) FilterUnstaked(opts *bind.FilterOpts, txOriginator []common.Address) (*ValidatorregistryUnstakedIterator, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.FilterLogs(opts, "Unstaked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorregistryUnstakedIterator{contract: _Validatorregistry.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x34f6c826286e3bec36208d054dcb1ad3c49725237a7644e1a6d157a92ae7a3e1.
//
// Solidity: event Unstaked(address indexed txOriginator, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistry *ValidatorregistryFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *ValidatorregistryUnstaked, txOriginator []common.Address) (event.Subscription, error) {

	var txOriginatorRule []interface{}
	for _, txOriginatorItem := range txOriginator {
		txOriginatorRule = append(txOriginatorRule, txOriginatorItem)
	}

	logs, sub, err := _Validatorregistry.contract.WatchLogs(opts, "Unstaked", txOriginatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorregistryUnstaked)
				if err := _Validatorregistry.contract.UnpackLog(event, "Unstaked", log); err != nil {
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
func (_Validatorregistry *ValidatorregistryFilterer) ParseUnstaked(log types.Log) (*ValidatorregistryUnstaked, error) {
	event := new(ValidatorregistryUnstaked)
	if err := _Validatorregistry.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
