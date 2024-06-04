// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegationvalreg

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

// DelegationValRegDelegationInfo is an auto generated low-level Go binding around an user-defined struct.
type DelegationValRegDelegationInfo struct {
	State          uint8
	ValidatorEOA   common.Address
	Amount         *big.Int
	WithdrawHeight *big.Int
}

// DelegationvalregMetaData contains all meta data concerning the Delegationvalreg contract.
var DelegationvalregMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DEFAULT_STETH_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"changeDelegation\",\"inputs\":[{\"name\":\"newValidatorEOA\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegate\",\"inputs\":[{\"name\":\"validatorEOA\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumDelegationValReg.State\"},{\"name\":\"validatorEOA\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDelegationInfo\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDelegationValReg.DelegationInfo\",\"components\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumDelegationValReg.State\"},{\"name\":\"validatorEOA\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSenderDelegationState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumDelegationValReg.State\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_reputationValReg\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_withdrawPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_stETHToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"reputationValReg\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractReputationValReg\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"requestWithdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stETHToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Delegated\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"validatorEOA\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegationChanged\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldValidatorEOA\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newValidatorEOA\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawRequested\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"validatorEOA\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// DelegationvalregABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegationvalregMetaData.ABI instead.
var DelegationvalregABI = DelegationvalregMetaData.ABI

// Delegationvalreg is an auto generated Go binding around an Ethereum contract.
type Delegationvalreg struct {
	DelegationvalregCaller     // Read-only binding to the contract
	DelegationvalregTransactor // Write-only binding to the contract
	DelegationvalregFilterer   // Log filterer for contract events
}

// DelegationvalregCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegationvalregCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationvalregTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegationvalregTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationvalregFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegationvalregFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationvalregSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegationvalregSession struct {
	Contract     *Delegationvalreg // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegationvalregCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegationvalregCallerSession struct {
	Contract *DelegationvalregCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// DelegationvalregTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegationvalregTransactorSession struct {
	Contract     *DelegationvalregTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// DelegationvalregRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegationvalregRaw struct {
	Contract *Delegationvalreg // Generic contract binding to access the raw methods on
}

// DelegationvalregCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegationvalregCallerRaw struct {
	Contract *DelegationvalregCaller // Generic read-only contract binding to access the raw methods on
}

// DelegationvalregTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegationvalregTransactorRaw struct {
	Contract *DelegationvalregTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegationvalreg creates a new instance of Delegationvalreg, bound to a specific deployed contract.
func NewDelegationvalreg(address common.Address, backend bind.ContractBackend) (*Delegationvalreg, error) {
	contract, err := bindDelegationvalreg(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegationvalreg{DelegationvalregCaller: DelegationvalregCaller{contract: contract}, DelegationvalregTransactor: DelegationvalregTransactor{contract: contract}, DelegationvalregFilterer: DelegationvalregFilterer{contract: contract}}, nil
}

// NewDelegationvalregCaller creates a new read-only instance of Delegationvalreg, bound to a specific deployed contract.
func NewDelegationvalregCaller(address common.Address, caller bind.ContractCaller) (*DelegationvalregCaller, error) {
	contract, err := bindDelegationvalreg(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregCaller{contract: contract}, nil
}

// NewDelegationvalregTransactor creates a new write-only instance of Delegationvalreg, bound to a specific deployed contract.
func NewDelegationvalregTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegationvalregTransactor, error) {
	contract, err := bindDelegationvalreg(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregTransactor{contract: contract}, nil
}

// NewDelegationvalregFilterer creates a new log filterer instance of Delegationvalreg, bound to a specific deployed contract.
func NewDelegationvalregFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegationvalregFilterer, error) {
	contract, err := bindDelegationvalreg(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregFilterer{contract: contract}, nil
}

// bindDelegationvalreg binds a generic wrapper to an already deployed contract.
func bindDelegationvalreg(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegationvalregMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegationvalreg *DelegationvalregRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegationvalreg.Contract.DelegationvalregCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegationvalreg *DelegationvalregRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.DelegationvalregTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegationvalreg *DelegationvalregRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.DelegationvalregTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegationvalreg *DelegationvalregCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegationvalreg.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegationvalreg *DelegationvalregTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegationvalreg *DelegationvalregTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTSTETHADDRESS is a free data retrieval call binding the contract method 0xfb854a2b.
//
// Solidity: function DEFAULT_STETH_ADDRESS() view returns(address)
func (_Delegationvalreg *DelegationvalregCaller) DEFAULTSTETHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "DEFAULT_STETH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DEFAULTSTETHADDRESS is a free data retrieval call binding the contract method 0xfb854a2b.
//
// Solidity: function DEFAULT_STETH_ADDRESS() view returns(address)
func (_Delegationvalreg *DelegationvalregSession) DEFAULTSTETHADDRESS() (common.Address, error) {
	return _Delegationvalreg.Contract.DEFAULTSTETHADDRESS(&_Delegationvalreg.CallOpts)
}

// DEFAULTSTETHADDRESS is a free data retrieval call binding the contract method 0xfb854a2b.
//
// Solidity: function DEFAULT_STETH_ADDRESS() view returns(address)
func (_Delegationvalreg *DelegationvalregCallerSession) DEFAULTSTETHADDRESS() (common.Address, error) {
	return _Delegationvalreg.Contract.DEFAULTSTETHADDRESS(&_Delegationvalreg.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Delegationvalreg *DelegationvalregCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Delegationvalreg *DelegationvalregSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Delegationvalreg.Contract.UPGRADEINTERFACEVERSION(&_Delegationvalreg.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Delegationvalreg *DelegationvalregCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Delegationvalreg.Contract.UPGRADEINTERFACEVERSION(&_Delegationvalreg.CallOpts)
}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) view returns(uint8 state, address validatorEOA, uint256 amount, uint256 withdrawHeight)
func (_Delegationvalreg *DelegationvalregCaller) Delegations(opts *bind.CallOpts, arg0 common.Address) (struct {
	State          uint8
	ValidatorEOA   common.Address
	Amount         *big.Int
	WithdrawHeight *big.Int
}, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "delegations", arg0)

	outstruct := new(struct {
		State          uint8
		ValidatorEOA   common.Address
		Amount         *big.Int
		WithdrawHeight *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.State = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.ValidatorEOA = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.WithdrawHeight = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) view returns(uint8 state, address validatorEOA, uint256 amount, uint256 withdrawHeight)
func (_Delegationvalreg *DelegationvalregSession) Delegations(arg0 common.Address) (struct {
	State          uint8
	ValidatorEOA   common.Address
	Amount         *big.Int
	WithdrawHeight *big.Int
}, error) {
	return _Delegationvalreg.Contract.Delegations(&_Delegationvalreg.CallOpts, arg0)
}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) view returns(uint8 state, address validatorEOA, uint256 amount, uint256 withdrawHeight)
func (_Delegationvalreg *DelegationvalregCallerSession) Delegations(arg0 common.Address) (struct {
	State          uint8
	ValidatorEOA   common.Address
	Amount         *big.Int
	WithdrawHeight *big.Int
}, error) {
	return _Delegationvalreg.Contract.Delegations(&_Delegationvalreg.CallOpts, arg0)
}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address delegator) view returns((uint8,address,uint256,uint256))
func (_Delegationvalreg *DelegationvalregCaller) GetDelegationInfo(opts *bind.CallOpts, delegator common.Address) (DelegationValRegDelegationInfo, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "getDelegationInfo", delegator)

	if err != nil {
		return *new(DelegationValRegDelegationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(DelegationValRegDelegationInfo)).(*DelegationValRegDelegationInfo)

	return out0, err

}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address delegator) view returns((uint8,address,uint256,uint256))
func (_Delegationvalreg *DelegationvalregSession) GetDelegationInfo(delegator common.Address) (DelegationValRegDelegationInfo, error) {
	return _Delegationvalreg.Contract.GetDelegationInfo(&_Delegationvalreg.CallOpts, delegator)
}

// GetDelegationInfo is a free data retrieval call binding the contract method 0xfab46d66.
//
// Solidity: function getDelegationInfo(address delegator) view returns((uint8,address,uint256,uint256))
func (_Delegationvalreg *DelegationvalregCallerSession) GetDelegationInfo(delegator common.Address) (DelegationValRegDelegationInfo, error) {
	return _Delegationvalreg.Contract.GetDelegationInfo(&_Delegationvalreg.CallOpts, delegator)
}

// GetSenderDelegationState is a free data retrieval call binding the contract method 0x469b58be.
//
// Solidity: function getSenderDelegationState() view returns(uint8)
func (_Delegationvalreg *DelegationvalregCaller) GetSenderDelegationState(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "getSenderDelegationState")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetSenderDelegationState is a free data retrieval call binding the contract method 0x469b58be.
//
// Solidity: function getSenderDelegationState() view returns(uint8)
func (_Delegationvalreg *DelegationvalregSession) GetSenderDelegationState() (uint8, error) {
	return _Delegationvalreg.Contract.GetSenderDelegationState(&_Delegationvalreg.CallOpts)
}

// GetSenderDelegationState is a free data retrieval call binding the contract method 0x469b58be.
//
// Solidity: function getSenderDelegationState() view returns(uint8)
func (_Delegationvalreg *DelegationvalregCallerSession) GetSenderDelegationState() (uint8, error) {
	return _Delegationvalreg.Contract.GetSenderDelegationState(&_Delegationvalreg.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegationvalreg *DelegationvalregCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegationvalreg *DelegationvalregSession) Owner() (common.Address, error) {
	return _Delegationvalreg.Contract.Owner(&_Delegationvalreg.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Delegationvalreg *DelegationvalregCallerSession) Owner() (common.Address, error) {
	return _Delegationvalreg.Contract.Owner(&_Delegationvalreg.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Delegationvalreg *DelegationvalregCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Delegationvalreg *DelegationvalregSession) ProxiableUUID() ([32]byte, error) {
	return _Delegationvalreg.Contract.ProxiableUUID(&_Delegationvalreg.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Delegationvalreg *DelegationvalregCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Delegationvalreg.Contract.ProxiableUUID(&_Delegationvalreg.CallOpts)
}

// ReputationValReg is a free data retrieval call binding the contract method 0x82dcba79.
//
// Solidity: function reputationValReg() view returns(address)
func (_Delegationvalreg *DelegationvalregCaller) ReputationValReg(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "reputationValReg")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ReputationValReg is a free data retrieval call binding the contract method 0x82dcba79.
//
// Solidity: function reputationValReg() view returns(address)
func (_Delegationvalreg *DelegationvalregSession) ReputationValReg() (common.Address, error) {
	return _Delegationvalreg.Contract.ReputationValReg(&_Delegationvalreg.CallOpts)
}

// ReputationValReg is a free data retrieval call binding the contract method 0x82dcba79.
//
// Solidity: function reputationValReg() view returns(address)
func (_Delegationvalreg *DelegationvalregCallerSession) ReputationValReg() (common.Address, error) {
	return _Delegationvalreg.Contract.ReputationValReg(&_Delegationvalreg.CallOpts)
}

// StETHToken is a free data retrieval call binding the contract method 0xaae92f1e.
//
// Solidity: function stETHToken() view returns(address)
func (_Delegationvalreg *DelegationvalregCaller) StETHToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "stETHToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StETHToken is a free data retrieval call binding the contract method 0xaae92f1e.
//
// Solidity: function stETHToken() view returns(address)
func (_Delegationvalreg *DelegationvalregSession) StETHToken() (common.Address, error) {
	return _Delegationvalreg.Contract.StETHToken(&_Delegationvalreg.CallOpts)
}

// StETHToken is a free data retrieval call binding the contract method 0xaae92f1e.
//
// Solidity: function stETHToken() view returns(address)
func (_Delegationvalreg *DelegationvalregCallerSession) StETHToken() (common.Address, error) {
	return _Delegationvalreg.Contract.StETHToken(&_Delegationvalreg.CallOpts)
}

// WithdrawPeriod is a free data retrieval call binding the contract method 0x12eb4f9a.
//
// Solidity: function withdrawPeriod() view returns(uint256)
func (_Delegationvalreg *DelegationvalregCaller) WithdrawPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Delegationvalreg.contract.Call(opts, &out, "withdrawPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawPeriod is a free data retrieval call binding the contract method 0x12eb4f9a.
//
// Solidity: function withdrawPeriod() view returns(uint256)
func (_Delegationvalreg *DelegationvalregSession) WithdrawPeriod() (*big.Int, error) {
	return _Delegationvalreg.Contract.WithdrawPeriod(&_Delegationvalreg.CallOpts)
}

// WithdrawPeriod is a free data retrieval call binding the contract method 0x12eb4f9a.
//
// Solidity: function withdrawPeriod() view returns(uint256)
func (_Delegationvalreg *DelegationvalregCallerSession) WithdrawPeriod() (*big.Int, error) {
	return _Delegationvalreg.Contract.WithdrawPeriod(&_Delegationvalreg.CallOpts)
}

// ChangeDelegation is a paid mutator transaction binding the contract method 0x9f973fd5.
//
// Solidity: function changeDelegation(address newValidatorEOA) returns()
func (_Delegationvalreg *DelegationvalregTransactor) ChangeDelegation(opts *bind.TransactOpts, newValidatorEOA common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "changeDelegation", newValidatorEOA)
}

// ChangeDelegation is a paid mutator transaction binding the contract method 0x9f973fd5.
//
// Solidity: function changeDelegation(address newValidatorEOA) returns()
func (_Delegationvalreg *DelegationvalregSession) ChangeDelegation(newValidatorEOA common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.ChangeDelegation(&_Delegationvalreg.TransactOpts, newValidatorEOA)
}

// ChangeDelegation is a paid mutator transaction binding the contract method 0x9f973fd5.
//
// Solidity: function changeDelegation(address newValidatorEOA) returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) ChangeDelegation(newValidatorEOA common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.ChangeDelegation(&_Delegationvalreg.TransactOpts, newValidatorEOA)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorEOA, uint256 amount) returns()
func (_Delegationvalreg *DelegationvalregTransactor) Delegate(opts *bind.TransactOpts, validatorEOA common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "delegate", validatorEOA, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorEOA, uint256 amount) returns()
func (_Delegationvalreg *DelegationvalregSession) Delegate(validatorEOA common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Delegate(&_Delegationvalreg.TransactOpts, validatorEOA, amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address validatorEOA, uint256 amount) returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) Delegate(validatorEOA common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Delegate(&_Delegationvalreg.TransactOpts, validatorEOA, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xbe203094.
//
// Solidity: function initialize(address _owner, address _reputationValReg, uint256 _withdrawPeriod, address _stETHToken) returns()
func (_Delegationvalreg *DelegationvalregTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _reputationValReg common.Address, _withdrawPeriod *big.Int, _stETHToken common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "initialize", _owner, _reputationValReg, _withdrawPeriod, _stETHToken)
}

// Initialize is a paid mutator transaction binding the contract method 0xbe203094.
//
// Solidity: function initialize(address _owner, address _reputationValReg, uint256 _withdrawPeriod, address _stETHToken) returns()
func (_Delegationvalreg *DelegationvalregSession) Initialize(_owner common.Address, _reputationValReg common.Address, _withdrawPeriod *big.Int, _stETHToken common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Initialize(&_Delegationvalreg.TransactOpts, _owner, _reputationValReg, _withdrawPeriod, _stETHToken)
}

// Initialize is a paid mutator transaction binding the contract method 0xbe203094.
//
// Solidity: function initialize(address _owner, address _reputationValReg, uint256 _withdrawPeriod, address _stETHToken) returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) Initialize(_owner common.Address, _reputationValReg common.Address, _withdrawPeriod *big.Int, _stETHToken common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Initialize(&_Delegationvalreg.TransactOpts, _owner, _reputationValReg, _withdrawPeriod, _stETHToken)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationvalreg *DelegationvalregTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationvalreg *DelegationvalregSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.RenounceOwnership(&_Delegationvalreg.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.RenounceOwnership(&_Delegationvalreg.TransactOpts)
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0xb3423eec.
//
// Solidity: function requestWithdraw() returns()
func (_Delegationvalreg *DelegationvalregTransactor) RequestWithdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "requestWithdraw")
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0xb3423eec.
//
// Solidity: function requestWithdraw() returns()
func (_Delegationvalreg *DelegationvalregSession) RequestWithdraw() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.RequestWithdraw(&_Delegationvalreg.TransactOpts)
}

// RequestWithdraw is a paid mutator transaction binding the contract method 0xb3423eec.
//
// Solidity: function requestWithdraw() returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) RequestWithdraw() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.RequestWithdraw(&_Delegationvalreg.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationvalreg *DelegationvalregTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationvalreg *DelegationvalregSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.TransferOwnership(&_Delegationvalreg.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.TransferOwnership(&_Delegationvalreg.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Delegationvalreg *DelegationvalregTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Delegationvalreg *DelegationvalregSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.UpgradeToAndCall(&_Delegationvalreg.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.UpgradeToAndCall(&_Delegationvalreg.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Delegationvalreg *DelegationvalregTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Delegationvalreg *DelegationvalregSession) Withdraw() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Withdraw(&_Delegationvalreg.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Withdraw(&_Delegationvalreg.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Delegationvalreg *DelegationvalregTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Delegationvalreg.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Delegationvalreg *DelegationvalregSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Fallback(&_Delegationvalreg.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Fallback(&_Delegationvalreg.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegationvalreg *DelegationvalregTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationvalreg.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegationvalreg *DelegationvalregSession) Receive() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Receive(&_Delegationvalreg.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Delegationvalreg *DelegationvalregTransactorSession) Receive() (*types.Transaction, error) {
	return _Delegationvalreg.Contract.Receive(&_Delegationvalreg.TransactOpts)
}

// DelegationvalregDelegatedIterator is returned from FilterDelegated and is used to iterate over the raw logs and unpacked data for Delegated events raised by the Delegationvalreg contract.
type DelegationvalregDelegatedIterator struct {
	Event *DelegationvalregDelegated // Event containing the contract specifics and raw log

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
func (it *DelegationvalregDelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregDelegated)
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
		it.Event = new(DelegationvalregDelegated)
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
func (it *DelegationvalregDelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregDelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregDelegated represents a Delegated event raised by the Delegationvalreg contract.
type DelegationvalregDelegated struct {
	Delegator    common.Address
	ValidatorEOA common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDelegated is a free log retrieval operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) FilterDelegated(opts *bind.FilterOpts, delegator []common.Address, validatorEOA []common.Address) (*DelegationvalregDelegatedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorEOARule []interface{}
	for _, validatorEOAItem := range validatorEOA {
		validatorEOARule = append(validatorEOARule, validatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "Delegated", delegatorRule, validatorEOARule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregDelegatedIterator{contract: _Delegationvalreg.contract, event: "Delegated", logs: logs, sub: sub}, nil
}

// WatchDelegated is a free log subscription operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) WatchDelegated(opts *bind.WatchOpts, sink chan<- *DelegationvalregDelegated, delegator []common.Address, validatorEOA []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorEOARule []interface{}
	for _, validatorEOAItem := range validatorEOA {
		validatorEOARule = append(validatorEOARule, validatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "Delegated", delegatorRule, validatorEOARule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregDelegated)
				if err := _Delegationvalreg.contract.UnpackLog(event, "Delegated", log); err != nil {
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

// ParseDelegated is a log parse operation binding the contract event 0xe5541a6b6103d4fa7e021ed54fad39c66f27a76bd13d374cf6240ae6bd0bb72b.
//
// Solidity: event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) ParseDelegated(log types.Log) (*DelegationvalregDelegated, error) {
	event := new(DelegationvalregDelegated)
	if err := _Delegationvalreg.contract.UnpackLog(event, "Delegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregDelegationChangedIterator is returned from FilterDelegationChanged and is used to iterate over the raw logs and unpacked data for DelegationChanged events raised by the Delegationvalreg contract.
type DelegationvalregDelegationChangedIterator struct {
	Event *DelegationvalregDelegationChanged // Event containing the contract specifics and raw log

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
func (it *DelegationvalregDelegationChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregDelegationChanged)
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
		it.Event = new(DelegationvalregDelegationChanged)
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
func (it *DelegationvalregDelegationChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregDelegationChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregDelegationChanged represents a DelegationChanged event raised by the Delegationvalreg contract.
type DelegationvalregDelegationChanged struct {
	Delegator       common.Address
	OldValidatorEOA common.Address
	NewValidatorEOA common.Address
	Amount          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDelegationChanged is a free log retrieval operation binding the contract event 0xc9705723527d1ad7e83b1f8cf2c5d15ff9dd7eb962beef1da94a69d46198fd6c.
//
// Solidity: event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) FilterDelegationChanged(opts *bind.FilterOpts, delegator []common.Address, oldValidatorEOA []common.Address, newValidatorEOA []common.Address) (*DelegationvalregDelegationChangedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var oldValidatorEOARule []interface{}
	for _, oldValidatorEOAItem := range oldValidatorEOA {
		oldValidatorEOARule = append(oldValidatorEOARule, oldValidatorEOAItem)
	}
	var newValidatorEOARule []interface{}
	for _, newValidatorEOAItem := range newValidatorEOA {
		newValidatorEOARule = append(newValidatorEOARule, newValidatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "DelegationChanged", delegatorRule, oldValidatorEOARule, newValidatorEOARule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregDelegationChangedIterator{contract: _Delegationvalreg.contract, event: "DelegationChanged", logs: logs, sub: sub}, nil
}

// WatchDelegationChanged is a free log subscription operation binding the contract event 0xc9705723527d1ad7e83b1f8cf2c5d15ff9dd7eb962beef1da94a69d46198fd6c.
//
// Solidity: event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) WatchDelegationChanged(opts *bind.WatchOpts, sink chan<- *DelegationvalregDelegationChanged, delegator []common.Address, oldValidatorEOA []common.Address, newValidatorEOA []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var oldValidatorEOARule []interface{}
	for _, oldValidatorEOAItem := range oldValidatorEOA {
		oldValidatorEOARule = append(oldValidatorEOARule, oldValidatorEOAItem)
	}
	var newValidatorEOARule []interface{}
	for _, newValidatorEOAItem := range newValidatorEOA {
		newValidatorEOARule = append(newValidatorEOARule, newValidatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "DelegationChanged", delegatorRule, oldValidatorEOARule, newValidatorEOARule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregDelegationChanged)
				if err := _Delegationvalreg.contract.UnpackLog(event, "DelegationChanged", log); err != nil {
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

// ParseDelegationChanged is a log parse operation binding the contract event 0xc9705723527d1ad7e83b1f8cf2c5d15ff9dd7eb962beef1da94a69d46198fd6c.
//
// Solidity: event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) ParseDelegationChanged(log types.Log) (*DelegationvalregDelegationChanged, error) {
	event := new(DelegationvalregDelegationChanged)
	if err := _Delegationvalreg.contract.UnpackLog(event, "DelegationChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Delegationvalreg contract.
type DelegationvalregInitializedIterator struct {
	Event *DelegationvalregInitialized // Event containing the contract specifics and raw log

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
func (it *DelegationvalregInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregInitialized)
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
		it.Event = new(DelegationvalregInitialized)
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
func (it *DelegationvalregInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregInitialized represents a Initialized event raised by the Delegationvalreg contract.
type DelegationvalregInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Delegationvalreg *DelegationvalregFilterer) FilterInitialized(opts *bind.FilterOpts) (*DelegationvalregInitializedIterator, error) {

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DelegationvalregInitializedIterator{contract: _Delegationvalreg.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Delegationvalreg *DelegationvalregFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DelegationvalregInitialized) (event.Subscription, error) {

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregInitialized)
				if err := _Delegationvalreg.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Delegationvalreg *DelegationvalregFilterer) ParseInitialized(log types.Log) (*DelegationvalregInitialized, error) {
	event := new(DelegationvalregInitialized)
	if err := _Delegationvalreg.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Delegationvalreg contract.
type DelegationvalregOwnershipTransferredIterator struct {
	Event *DelegationvalregOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DelegationvalregOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregOwnershipTransferred)
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
		it.Event = new(DelegationvalregOwnershipTransferred)
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
func (it *DelegationvalregOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregOwnershipTransferred represents a OwnershipTransferred event raised by the Delegationvalreg contract.
type DelegationvalregOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegationvalreg *DelegationvalregFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DelegationvalregOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregOwnershipTransferredIterator{contract: _Delegationvalreg.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegationvalreg *DelegationvalregFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DelegationvalregOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregOwnershipTransferred)
				if err := _Delegationvalreg.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Delegationvalreg *DelegationvalregFilterer) ParseOwnershipTransferred(log types.Log) (*DelegationvalregOwnershipTransferred, error) {
	event := new(DelegationvalregOwnershipTransferred)
	if err := _Delegationvalreg.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Delegationvalreg contract.
type DelegationvalregUpgradedIterator struct {
	Event *DelegationvalregUpgraded // Event containing the contract specifics and raw log

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
func (it *DelegationvalregUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregUpgraded)
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
		it.Event = new(DelegationvalregUpgraded)
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
func (it *DelegationvalregUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregUpgraded represents a Upgraded event raised by the Delegationvalreg contract.
type DelegationvalregUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Delegationvalreg *DelegationvalregFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DelegationvalregUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregUpgradedIterator{contract: _Delegationvalreg.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Delegationvalreg *DelegationvalregFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DelegationvalregUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregUpgraded)
				if err := _Delegationvalreg.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Delegationvalreg *DelegationvalregFilterer) ParseUpgraded(log types.Log) (*DelegationvalregUpgraded, error) {
	event := new(DelegationvalregUpgraded)
	if err := _Delegationvalreg.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregWithdrawRequestedIterator is returned from FilterWithdrawRequested and is used to iterate over the raw logs and unpacked data for WithdrawRequested events raised by the Delegationvalreg contract.
type DelegationvalregWithdrawRequestedIterator struct {
	Event *DelegationvalregWithdrawRequested // Event containing the contract specifics and raw log

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
func (it *DelegationvalregWithdrawRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregWithdrawRequested)
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
		it.Event = new(DelegationvalregWithdrawRequested)
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
func (it *DelegationvalregWithdrawRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregWithdrawRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregWithdrawRequested represents a WithdrawRequested event raised by the Delegationvalreg contract.
type DelegationvalregWithdrawRequested struct {
	Delegator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawRequested is a free log retrieval operation binding the contract event 0x82d5fce39f90fafd1445cdc4e4851843dfe8406dbba6112d2eeb23c6c860b490.
//
// Solidity: event WithdrawRequested(address indexed delegator)
func (_Delegationvalreg *DelegationvalregFilterer) FilterWithdrawRequested(opts *bind.FilterOpts, delegator []common.Address) (*DelegationvalregWithdrawRequestedIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "WithdrawRequested", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregWithdrawRequestedIterator{contract: _Delegationvalreg.contract, event: "WithdrawRequested", logs: logs, sub: sub}, nil
}

// WatchWithdrawRequested is a free log subscription operation binding the contract event 0x82d5fce39f90fafd1445cdc4e4851843dfe8406dbba6112d2eeb23c6c860b490.
//
// Solidity: event WithdrawRequested(address indexed delegator)
func (_Delegationvalreg *DelegationvalregFilterer) WatchWithdrawRequested(opts *bind.WatchOpts, sink chan<- *DelegationvalregWithdrawRequested, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "WithdrawRequested", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregWithdrawRequested)
				if err := _Delegationvalreg.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
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

// ParseWithdrawRequested is a log parse operation binding the contract event 0x82d5fce39f90fafd1445cdc4e4851843dfe8406dbba6112d2eeb23c6c860b490.
//
// Solidity: event WithdrawRequested(address indexed delegator)
func (_Delegationvalreg *DelegationvalregFilterer) ParseWithdrawRequested(log types.Log) (*DelegationvalregWithdrawRequested, error) {
	event := new(DelegationvalregWithdrawRequested)
	if err := _Delegationvalreg.contract.UnpackLog(event, "WithdrawRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegationvalregWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Delegationvalreg contract.
type DelegationvalregWithdrawnIterator struct {
	Event *DelegationvalregWithdrawn // Event containing the contract specifics and raw log

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
func (it *DelegationvalregWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationvalregWithdrawn)
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
		it.Event = new(DelegationvalregWithdrawn)
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
func (it *DelegationvalregWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationvalregWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationvalregWithdrawn represents a Withdrawn event raised by the Delegationvalreg contract.
type DelegationvalregWithdrawn struct {
	Delegator    common.Address
	ValidatorEOA common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) FilterWithdrawn(opts *bind.FilterOpts, delegator []common.Address, validatorEOA []common.Address) (*DelegationvalregWithdrawnIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorEOARule []interface{}
	for _, validatorEOAItem := range validatorEOA {
		validatorEOARule = append(validatorEOARule, validatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.FilterLogs(opts, "Withdrawn", delegatorRule, validatorEOARule)
	if err != nil {
		return nil, err
	}
	return &DelegationvalregWithdrawnIterator{contract: _Delegationvalreg.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *DelegationvalregWithdrawn, delegator []common.Address, validatorEOA []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var validatorEOARule []interface{}
	for _, validatorEOAItem := range validatorEOA {
		validatorEOARule = append(validatorEOARule, validatorEOAItem)
	}

	logs, sub, err := _Delegationvalreg.contract.WatchLogs(opts, "Withdrawn", delegatorRule, validatorEOARule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationvalregWithdrawn)
				if err := _Delegationvalreg.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xd1c19fbcd4551a5edfb66d43d2e337c04837afda3482b42bdf569a8fccdae5fb.
//
// Solidity: event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount)
func (_Delegationvalreg *DelegationvalregFilterer) ParseWithdrawn(log types.Log) (*DelegationvalregWithdrawn, error) {
	event := new(DelegationvalregWithdrawn)
	if err := _Delegationvalreg.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
