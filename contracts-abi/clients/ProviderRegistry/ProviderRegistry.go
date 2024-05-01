// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package providerregistry

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

// ProviderregistryMetaData contains all meta data concerning the Providerregistry contract.
var ProviderregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"bidderAmount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkStake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"feePercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipientAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"preConfirmationsContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerStakes\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAndStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeePercent\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeeRecipient\",\"inputs\":[{\"name\":\"newFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPreconfirmationsContract\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slash\",\"inputs\":[{\"name\":\"amt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawBidderAmount\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFeeRecipientAmount\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawStakedAmount\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FundsDeposited\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsSlashed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderRegistered\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"stakedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]}]",
}

// ProviderregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ProviderregistryMetaData.ABI instead.
var ProviderregistryABI = ProviderregistryMetaData.ABI

// Providerregistry is an auto generated Go binding around an Ethereum contract.
type Providerregistry struct {
	ProviderregistryCaller     // Read-only binding to the contract
	ProviderregistryTransactor // Write-only binding to the contract
	ProviderregistryFilterer   // Log filterer for contract events
}

// ProviderregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProviderregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProviderregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProviderregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProviderregistrySession struct {
	Contract     *Providerregistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProviderregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProviderregistryCallerSession struct {
	Contract *ProviderregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ProviderregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProviderregistryTransactorSession struct {
	Contract     *ProviderregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ProviderregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProviderregistryRaw struct {
	Contract *Providerregistry // Generic contract binding to access the raw methods on
}

// ProviderregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProviderregistryCallerRaw struct {
	Contract *ProviderregistryCaller // Generic read-only contract binding to access the raw methods on
}

// ProviderregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProviderregistryTransactorRaw struct {
	Contract *ProviderregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProviderregistry creates a new instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistry(address common.Address, backend bind.ContractBackend) (*Providerregistry, error) {
	contract, err := bindProviderregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Providerregistry{ProviderregistryCaller: ProviderregistryCaller{contract: contract}, ProviderregistryTransactor: ProviderregistryTransactor{contract: contract}, ProviderregistryFilterer: ProviderregistryFilterer{contract: contract}}, nil
}

// NewProviderregistryCaller creates a new read-only instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryCaller(address common.Address, caller bind.ContractCaller) (*ProviderregistryCaller, error) {
	contract, err := bindProviderregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryCaller{contract: contract}, nil
}

// NewProviderregistryTransactor creates a new write-only instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ProviderregistryTransactor, error) {
	contract, err := bindProviderregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryTransactor{contract: contract}, nil
}

// NewProviderregistryFilterer creates a new log filterer instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ProviderregistryFilterer, error) {
	contract, err := bindProviderregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFilterer{contract: contract}, nil
}

// bindProviderregistry binds a generic wrapper to an already deployed contract.
func bindProviderregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProviderregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Providerregistry *ProviderregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Providerregistry.Contract.ProviderregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Providerregistry *ProviderregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.Contract.ProviderregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Providerregistry *ProviderregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Providerregistry.Contract.ProviderregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Providerregistry *ProviderregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Providerregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Providerregistry *ProviderregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Providerregistry *ProviderregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Providerregistry.Contract.contract.Transact(opts, method, params...)
}

// BidderAmount is a free data retrieval call binding the contract method 0x218a5e67.
//
// Solidity: function bidderAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) BidderAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "bidderAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BidderAmount is a free data retrieval call binding the contract method 0x218a5e67.
//
// Solidity: function bidderAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) BidderAmount(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.BidderAmount(&_Providerregistry.CallOpts, arg0)
}

// BidderAmount is a free data retrieval call binding the contract method 0x218a5e67.
//
// Solidity: function bidderAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) BidderAmount(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.BidderAmount(&_Providerregistry.CallOpts, arg0)
}

// CheckStake is a free data retrieval call binding the contract method 0x90d96d76.
//
// Solidity: function checkStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) CheckStake(opts *bind.CallOpts, provider common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "checkStake", provider)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CheckStake is a free data retrieval call binding the contract method 0x90d96d76.
//
// Solidity: function checkStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) CheckStake(provider common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.CheckStake(&_Providerregistry.CallOpts, provider)
}

// CheckStake is a free data retrieval call binding the contract method 0x90d96d76.
//
// Solidity: function checkStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) CheckStake(provider common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.CheckStake(&_Providerregistry.CallOpts, provider)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Providerregistry *ProviderregistryCaller) FeePercent(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "feePercent")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Providerregistry *ProviderregistrySession) FeePercent() (uint16, error) {
	return _Providerregistry.Contract.FeePercent(&_Providerregistry.CallOpts)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Providerregistry *ProviderregistryCallerSession) FeePercent() (uint16, error) {
	return _Providerregistry.Contract.FeePercent(&_Providerregistry.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Providerregistry *ProviderregistryCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Providerregistry *ProviderregistrySession) FeeRecipient() (common.Address, error) {
	return _Providerregistry.Contract.FeeRecipient(&_Providerregistry.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) FeeRecipient() (common.Address, error) {
	return _Providerregistry.Contract.FeeRecipient(&_Providerregistry.CallOpts)
}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) FeeRecipientAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "feeRecipientAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) FeeRecipientAmount() (*big.Int, error) {
	return _Providerregistry.Contract.FeeRecipientAmount(&_Providerregistry.CallOpts)
}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) FeeRecipientAmount() (*big.Int, error) {
	return _Providerregistry.Contract.FeeRecipientAmount(&_Providerregistry.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) MinStake() (*big.Int, error) {
	return _Providerregistry.Contract.MinStake(&_Providerregistry.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) MinStake() (*big.Int, error) {
	return _Providerregistry.Contract.MinStake(&_Providerregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistrySession) Owner() (common.Address, error) {
	return _Providerregistry.Contract.Owner(&_Providerregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) Owner() (common.Address, error) {
	return _Providerregistry.Contract.Owner(&_Providerregistry.CallOpts)
}

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Providerregistry *ProviderregistryCaller) PreConfirmationsContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "preConfirmationsContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Providerregistry *ProviderregistrySession) PreConfirmationsContract() (common.Address, error) {
	return _Providerregistry.Contract.PreConfirmationsContract(&_Providerregistry.CallOpts)
}

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) PreConfirmationsContract() (common.Address, error) {
	return _Providerregistry.Contract.PreConfirmationsContract(&_Providerregistry.CallOpts)
}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistryCaller) ProviderRegistered(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "providerRegistered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistrySession) ProviderRegistered(arg0 common.Address) (bool, error) {
	return _Providerregistry.Contract.ProviderRegistered(&_Providerregistry.CallOpts, arg0)
}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistryCallerSession) ProviderRegistered(arg0 common.Address) (bool, error) {
	return _Providerregistry.Contract.ProviderRegistered(&_Providerregistry.CallOpts, arg0)
}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) ProviderStakes(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "providerStakes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) ProviderStakes(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.ProviderStakes(&_Providerregistry.CallOpts, arg0)
}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) ProviderStakes(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.ProviderStakes(&_Providerregistry.CallOpts, arg0)
}

// DepositFunds is a paid mutator transaction binding the contract method 0xe2c41dbc.
//
// Solidity: function depositFunds() payable returns()
func (_Providerregistry *ProviderregistryTransactor) DepositFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "depositFunds")
}

// DepositFunds is a paid mutator transaction binding the contract method 0xe2c41dbc.
//
// Solidity: function depositFunds() payable returns()
func (_Providerregistry *ProviderregistrySession) DepositFunds() (*types.Transaction, error) {
	return _Providerregistry.Contract.DepositFunds(&_Providerregistry.TransactOpts)
}

// DepositFunds is a paid mutator transaction binding the contract method 0xe2c41dbc.
//
// Solidity: function depositFunds() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) DepositFunds() (*types.Transaction, error) {
	return _Providerregistry.Contract.DepositFunds(&_Providerregistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x66953e62.
//
// Solidity: function initialize(uint256 _minStake, address _feeRecipient, uint16 _feePercent, address _owner) returns()
func (_Providerregistry *ProviderregistryTransactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _feeRecipient common.Address, _feePercent uint16, _owner common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "initialize", _minStake, _feeRecipient, _feePercent, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x66953e62.
//
// Solidity: function initialize(uint256 _minStake, address _feeRecipient, uint16 _feePercent, address _owner) returns()
func (_Providerregistry *ProviderregistrySession) Initialize(_minStake *big.Int, _feeRecipient common.Address, _feePercent uint16, _owner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.Initialize(&_Providerregistry.TransactOpts, _minStake, _feeRecipient, _feePercent, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x66953e62.
//
// Solidity: function initialize(uint256 _minStake, address _feeRecipient, uint16 _feePercent, address _owner) returns()
func (_Providerregistry *ProviderregistryTransactorSession) Initialize(_minStake *big.Int, _feeRecipient common.Address, _feePercent uint16, _owner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.Initialize(&_Providerregistry.TransactOpts, _minStake, _feeRecipient, _feePercent, _owner)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistryTransactor) RegisterAndStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "registerAndStake")
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistrySession) RegisterAndStake() (*types.Transaction, error) {
	return _Providerregistry.Contract.RegisterAndStake(&_Providerregistry.TransactOpts)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) RegisterAndStake() (*types.Transaction, error) {
	return _Providerregistry.Contract.RegisterAndStake(&_Providerregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.RenounceOwnership(&_Providerregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.RenounceOwnership(&_Providerregistry.TransactOpts)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Providerregistry *ProviderregistryTransactor) SetNewFeePercent(opts *bind.TransactOpts, newFeePercent uint16) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setNewFeePercent", newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Providerregistry *ProviderregistrySession) SetNewFeePercent(newFeePercent uint16) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeePercent(&_Providerregistry.TransactOpts, newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetNewFeePercent(newFeePercent uint16) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeePercent(&_Providerregistry.TransactOpts, newFeePercent)
}

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistryTransactor) SetNewFeeRecipient(opts *bind.TransactOpts, newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setNewFeeRecipient", newFeeRecipient)
}

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistrySession) SetNewFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeeRecipient(&_Providerregistry.TransactOpts, newFeeRecipient)
}

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetNewFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeeRecipient(&_Providerregistry.TransactOpts, newFeeRecipient)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Providerregistry *ProviderregistryTransactor) SetPreconfirmationsContract(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setPreconfirmationsContract", contractAddress)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Providerregistry *ProviderregistrySession) SetPreconfirmationsContract(contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetPreconfirmationsContract(&_Providerregistry.TransactOpts, contractAddress)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetPreconfirmationsContract(contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetPreconfirmationsContract(&_Providerregistry.TransactOpts, contractAddress)
}

// Slash is a paid mutator transaction binding the contract method 0x8b6e1f8d.
//
// Solidity: function slash(uint256 amt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistryTransactor) Slash(opts *bind.TransactOpts, amt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "slash", amt, provider, bidder, residualBidPercentAfterDecay)
}

// Slash is a paid mutator transaction binding the contract method 0x8b6e1f8d.
//
// Solidity: function slash(uint256 amt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistrySession) Slash(amt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Slash(&_Providerregistry.TransactOpts, amt, provider, bidder, residualBidPercentAfterDecay)
}

// Slash is a paid mutator transaction binding the contract method 0x8b6e1f8d.
//
// Solidity: function slash(uint256 amt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistryTransactorSession) Slash(amt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Slash(&_Providerregistry.TransactOpts, amt, provider, bidder, residualBidPercentAfterDecay)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.TransferOwnership(&_Providerregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.TransferOwnership(&_Providerregistry.TransactOpts, newOwner)
}

// WithdrawBidderAmount is a paid mutator transaction binding the contract method 0x292ae408.
//
// Solidity: function withdrawBidderAmount(address bidder) returns()
func (_Providerregistry *ProviderregistryTransactor) WithdrawBidderAmount(opts *bind.TransactOpts, bidder common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "withdrawBidderAmount", bidder)
}

// WithdrawBidderAmount is a paid mutator transaction binding the contract method 0x292ae408.
//
// Solidity: function withdrawBidderAmount(address bidder) returns()
func (_Providerregistry *ProviderregistrySession) WithdrawBidderAmount(bidder common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawBidderAmount(&_Providerregistry.TransactOpts, bidder)
}

// WithdrawBidderAmount is a paid mutator transaction binding the contract method 0x292ae408.
//
// Solidity: function withdrawBidderAmount(address bidder) returns()
func (_Providerregistry *ProviderregistryTransactorSession) WithdrawBidderAmount(bidder common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawBidderAmount(&_Providerregistry.TransactOpts, bidder)
}

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Providerregistry *ProviderregistryTransactor) WithdrawFeeRecipientAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "withdrawFeeRecipientAmount")
}

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Providerregistry *ProviderregistrySession) WithdrawFeeRecipientAmount() (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawFeeRecipientAmount(&_Providerregistry.TransactOpts)
}

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Providerregistry *ProviderregistryTransactorSession) WithdrawFeeRecipientAmount() (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawFeeRecipientAmount(&_Providerregistry.TransactOpts)
}

// WithdrawStakedAmount is a paid mutator transaction binding the contract method 0xf31f3da0.
//
// Solidity: function withdrawStakedAmount(address provider) returns()
func (_Providerregistry *ProviderregistryTransactor) WithdrawStakedAmount(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "withdrawStakedAmount", provider)
}

// WithdrawStakedAmount is a paid mutator transaction binding the contract method 0xf31f3da0.
//
// Solidity: function withdrawStakedAmount(address provider) returns()
func (_Providerregistry *ProviderregistrySession) WithdrawStakedAmount(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawStakedAmount(&_Providerregistry.TransactOpts, provider)
}

// WithdrawStakedAmount is a paid mutator transaction binding the contract method 0xf31f3da0.
//
// Solidity: function withdrawStakedAmount(address provider) returns()
func (_Providerregistry *ProviderregistryTransactorSession) WithdrawStakedAmount(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawStakedAmount(&_Providerregistry.TransactOpts, provider)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.Fallback(&_Providerregistry.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.Fallback(&_Providerregistry.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistrySession) Receive() (*types.Transaction, error) {
	return _Providerregistry.Contract.Receive(&_Providerregistry.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) Receive() (*types.Transaction, error) {
	return _Providerregistry.Contract.Receive(&_Providerregistry.TransactOpts)
}

// ProviderregistryFundsDepositedIterator is returned from FilterFundsDeposited and is used to iterate over the raw logs and unpacked data for FundsDeposited events raised by the Providerregistry contract.
type ProviderregistryFundsDepositedIterator struct {
	Event *ProviderregistryFundsDeposited // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFundsDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFundsDeposited)
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
		it.Event = new(ProviderregistryFundsDeposited)
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
func (it *ProviderregistryFundsDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFundsDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFundsDeposited represents a FundsDeposited event raised by the Providerregistry contract.
type ProviderregistryFundsDeposited struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFundsDeposited is a free log retrieval operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterFundsDeposited(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryFundsDepositedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FundsDeposited", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFundsDepositedIterator{contract: _Providerregistry.contract, event: "FundsDeposited", logs: logs, sub: sub}, nil
}

// WatchFundsDeposited is a free log subscription operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchFundsDeposited(opts *bind.WatchOpts, sink chan<- *ProviderregistryFundsDeposited, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FundsDeposited", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFundsDeposited)
				if err := _Providerregistry.contract.UnpackLog(event, "FundsDeposited", log); err != nil {
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

// ParseFundsDeposited is a log parse operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseFundsDeposited(log types.Log) (*ProviderregistryFundsDeposited, error) {
	event := new(ProviderregistryFundsDeposited)
	if err := _Providerregistry.contract.UnpackLog(event, "FundsDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFundsSlashedIterator is returned from FilterFundsSlashed and is used to iterate over the raw logs and unpacked data for FundsSlashed events raised by the Providerregistry contract.
type ProviderregistryFundsSlashedIterator struct {
	Event *ProviderregistryFundsSlashed // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFundsSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFundsSlashed)
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
		it.Event = new(ProviderregistryFundsSlashed)
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
func (it *ProviderregistryFundsSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFundsSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFundsSlashed represents a FundsSlashed event raised by the Providerregistry contract.
type ProviderregistryFundsSlashed struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFundsSlashed is a free log retrieval operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterFundsSlashed(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryFundsSlashedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FundsSlashed", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFundsSlashedIterator{contract: _Providerregistry.contract, event: "FundsSlashed", logs: logs, sub: sub}, nil
}

// WatchFundsSlashed is a free log subscription operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchFundsSlashed(opts *bind.WatchOpts, sink chan<- *ProviderregistryFundsSlashed, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FundsSlashed", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFundsSlashed)
				if err := _Providerregistry.contract.UnpackLog(event, "FundsSlashed", log); err != nil {
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

// ParseFundsSlashed is a log parse operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseFundsSlashed(log types.Log) (*ProviderregistryFundsSlashed, error) {
	event := new(ProviderregistryFundsSlashed)
	if err := _Providerregistry.contract.UnpackLog(event, "FundsSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Providerregistry contract.
type ProviderregistryInitializedIterator struct {
	Event *ProviderregistryInitialized // Event containing the contract specifics and raw log

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
func (it *ProviderregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryInitialized)
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
		it.Event = new(ProviderregistryInitialized)
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
func (it *ProviderregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryInitialized represents a Initialized event raised by the Providerregistry contract.
type ProviderregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Providerregistry *ProviderregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ProviderregistryInitializedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryInitializedIterator{contract: _Providerregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Providerregistry *ProviderregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ProviderregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryInitialized)
				if err := _Providerregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseInitialized(log types.Log) (*ProviderregistryInitialized, error) {
	event := new(ProviderregistryInitialized)
	if err := _Providerregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferredIterator struct {
	Event *ProviderregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ProviderregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryOwnershipTransferred)
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
		it.Event = new(ProviderregistryOwnershipTransferred)
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
func (it *ProviderregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProviderregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryOwnershipTransferredIterator{contract: _Providerregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ProviderregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryOwnershipTransferred)
				if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseOwnershipTransferred(log types.Log) (*ProviderregistryOwnershipTransferred, error) {
	event := new(ProviderregistryOwnershipTransferred)
	if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryProviderRegisteredIterator is returned from FilterProviderRegistered and is used to iterate over the raw logs and unpacked data for ProviderRegistered events raised by the Providerregistry contract.
type ProviderregistryProviderRegisteredIterator struct {
	Event *ProviderregistryProviderRegistered // Event containing the contract specifics and raw log

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
func (it *ProviderregistryProviderRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryProviderRegistered)
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
		it.Event = new(ProviderregistryProviderRegistered)
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
func (it *ProviderregistryProviderRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryProviderRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryProviderRegistered represents a ProviderRegistered event raised by the Providerregistry contract.
type ProviderregistryProviderRegistered struct {
	Provider     common.Address
	StakedAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProviderRegistered is a free log retrieval operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) FilterProviderRegistered(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryProviderRegisteredIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryProviderRegisteredIterator{contract: _Providerregistry.contract, event: "ProviderRegistered", logs: logs, sub: sub}, nil
}

// WatchProviderRegistered is a free log subscription operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) WatchProviderRegistered(opts *bind.WatchOpts, sink chan<- *ProviderregistryProviderRegistered, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryProviderRegistered)
				if err := _Providerregistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
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

// ParseProviderRegistered is a log parse operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) ParseProviderRegistered(log types.Log) (*ProviderregistryProviderRegistered, error) {
	event := new(ProviderregistryProviderRegistered)
	if err := _Providerregistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
