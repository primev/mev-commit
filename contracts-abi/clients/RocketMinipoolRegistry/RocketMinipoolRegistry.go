// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rocketminipoolregistry

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

// IRocketMinipoolRegistryValidatorRegistration is an auto generated low-level Go binding around an user-defined struct.
type IRocketMinipoolRegistryValidatorRegistration struct {
	Exists          bool
	DeregTimestamp  uint64
	FreezeTimestamp uint64
}

// RocketminipoolregistryMetaData contains all meta data concerning the Rocketminipoolregistry contract.
var RocketminipoolregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterValidators\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregistrationPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"freeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEligibleTimeForDeregistration\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinipoolFromPubkey\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeAddressFromMinipool\",\"inputs\":[{\"name\":\"minipool\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeAddressFromPubkey\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidOperatorsForKey\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatorRegInfo\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRocketMinipoolRegistry.ValidatorRegistration\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"freezeTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freezeOracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezeReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rocketStorage\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregistrationPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isMinipoolActive\",\"inputs\":[{\"name\":\"minipool\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isOperatorValidForKey\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorRegistered\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerUnfreeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidators\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestValidatorDeregistration\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rocketStorage\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractRocketStorageInterface\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDeregistrationPeriod\",\"inputs\":[{\"name\":\"newDeregistrationPeriod\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFreezeOracle\",\"inputs\":[{\"name\":\"newFreezeOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRocketStorage\",\"inputs\":[{\"name\":\"newRocketStorage\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeFee\",\"inputs\":[{\"name\":\"newUnfreezeFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeReceiver\",\"inputs\":[{\"name\":\"newUnfreezeReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unfreezeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezeReceiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validatorRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"freezeTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"nodeAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistrationRequested\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"nodeAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorFrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRegistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"nodeAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnfrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DeregRequestAlreadyExists\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"DeregRequestDoesNotExist\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"DeregistrationTooSoon\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenValidatorCannotDeregister\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expectedLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MinipoolNotActive\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NoMinipoolForKey\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMinipoolOperator\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"OnlyFreezeOracle\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RefundFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnfreezeFeeRequired\",\"inputs\":[{\"name\":\"requiredFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnfreezeTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorAlreadyFrozen\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorAlreadyRegistered\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorDeregistrationNotExpired\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotFrozen\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotRegistered\",\"inputs\":[{\"name\":\"validatorPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ZeroParam\",\"inputs\":[]}]",
}

// RocketminipoolregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use RocketminipoolregistryMetaData.ABI instead.
var RocketminipoolregistryABI = RocketminipoolregistryMetaData.ABI

// Rocketminipoolregistry is an auto generated Go binding around an Ethereum contract.
type Rocketminipoolregistry struct {
	RocketminipoolregistryCaller     // Read-only binding to the contract
	RocketminipoolregistryTransactor // Write-only binding to the contract
	RocketminipoolregistryFilterer   // Log filterer for contract events
}

// RocketminipoolregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RocketminipoolregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RocketminipoolregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RocketminipoolregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RocketminipoolregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RocketminipoolregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RocketminipoolregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RocketminipoolregistrySession struct {
	Contract     *Rocketminipoolregistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// RocketminipoolregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RocketminipoolregistryCallerSession struct {
	Contract *RocketminipoolregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// RocketminipoolregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RocketminipoolregistryTransactorSession struct {
	Contract     *RocketminipoolregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// RocketminipoolregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RocketminipoolregistryRaw struct {
	Contract *Rocketminipoolregistry // Generic contract binding to access the raw methods on
}

// RocketminipoolregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RocketminipoolregistryCallerRaw struct {
	Contract *RocketminipoolregistryCaller // Generic read-only contract binding to access the raw methods on
}

// RocketminipoolregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RocketminipoolregistryTransactorRaw struct {
	Contract *RocketminipoolregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRocketminipoolregistry creates a new instance of Rocketminipoolregistry, bound to a specific deployed contract.
func NewRocketminipoolregistry(address common.Address, backend bind.ContractBackend) (*Rocketminipoolregistry, error) {
	contract, err := bindRocketminipoolregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rocketminipoolregistry{RocketminipoolregistryCaller: RocketminipoolregistryCaller{contract: contract}, RocketminipoolregistryTransactor: RocketminipoolregistryTransactor{contract: contract}, RocketminipoolregistryFilterer: RocketminipoolregistryFilterer{contract: contract}}, nil
}

// NewRocketminipoolregistryCaller creates a new read-only instance of Rocketminipoolregistry, bound to a specific deployed contract.
func NewRocketminipoolregistryCaller(address common.Address, caller bind.ContractCaller) (*RocketminipoolregistryCaller, error) {
	contract, err := bindRocketminipoolregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryCaller{contract: contract}, nil
}

// NewRocketminipoolregistryTransactor creates a new write-only instance of Rocketminipoolregistry, bound to a specific deployed contract.
func NewRocketminipoolregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*RocketminipoolregistryTransactor, error) {
	contract, err := bindRocketminipoolregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryTransactor{contract: contract}, nil
}

// NewRocketminipoolregistryFilterer creates a new log filterer instance of Rocketminipoolregistry, bound to a specific deployed contract.
func NewRocketminipoolregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*RocketminipoolregistryFilterer, error) {
	contract, err := bindRocketminipoolregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryFilterer{contract: contract}, nil
}

// bindRocketminipoolregistry binds a generic wrapper to an already deployed contract.
func bindRocketminipoolregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RocketminipoolregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rocketminipoolregistry *RocketminipoolregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rocketminipoolregistry.Contract.RocketminipoolregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rocketminipoolregistry *RocketminipoolregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RocketminipoolregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rocketminipoolregistry *RocketminipoolregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RocketminipoolregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rocketminipoolregistry *RocketminipoolregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rocketminipoolregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rocketminipoolregistry.Contract.UPGRADEINTERFACEVERSION(&_Rocketminipoolregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rocketminipoolregistry.Contract.UPGRADEINTERFACEVERSION(&_Rocketminipoolregistry.CallOpts)
}

// DeregistrationPeriod is a free data retrieval call binding the contract method 0x1667a3b6.
//
// Solidity: function deregistrationPeriod() view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) DeregistrationPeriod(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "deregistrationPeriod")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DeregistrationPeriod is a free data retrieval call binding the contract method 0x1667a3b6.
//
// Solidity: function deregistrationPeriod() view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) DeregistrationPeriod() (uint64, error) {
	return _Rocketminipoolregistry.Contract.DeregistrationPeriod(&_Rocketminipoolregistry.CallOpts)
}

// DeregistrationPeriod is a free data retrieval call binding the contract method 0x1667a3b6.
//
// Solidity: function deregistrationPeriod() view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) DeregistrationPeriod() (uint64, error) {
	return _Rocketminipoolregistry.Contract.DeregistrationPeriod(&_Rocketminipoolregistry.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) FreezeOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "freezeOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) FreezeOracle() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.FreezeOracle(&_Rocketminipoolregistry.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) FreezeOracle() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.FreezeOracle(&_Rocketminipoolregistry.CallOpts)
}

// GetEligibleTimeForDeregistration is a free data retrieval call binding the contract method 0x371a83ad.
//
// Solidity: function getEligibleTimeForDeregistration(bytes validatorPubkey) view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetEligibleTimeForDeregistration(opts *bind.CallOpts, validatorPubkey []byte) (uint64, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getEligibleTimeForDeregistration", validatorPubkey)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetEligibleTimeForDeregistration is a free data retrieval call binding the contract method 0x371a83ad.
//
// Solidity: function getEligibleTimeForDeregistration(bytes validatorPubkey) view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetEligibleTimeForDeregistration(validatorPubkey []byte) (uint64, error) {
	return _Rocketminipoolregistry.Contract.GetEligibleTimeForDeregistration(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetEligibleTimeForDeregistration is a free data retrieval call binding the contract method 0x371a83ad.
//
// Solidity: function getEligibleTimeForDeregistration(bytes validatorPubkey) view returns(uint64)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetEligibleTimeForDeregistration(validatorPubkey []byte) (uint64, error) {
	return _Rocketminipoolregistry.Contract.GetEligibleTimeForDeregistration(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetMinipoolFromPubkey is a free data retrieval call binding the contract method 0x6dc6b1ec.
//
// Solidity: function getMinipoolFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetMinipoolFromPubkey(opts *bind.CallOpts, validatorPubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getMinipoolFromPubkey", validatorPubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMinipoolFromPubkey is a free data retrieval call binding the contract method 0x6dc6b1ec.
//
// Solidity: function getMinipoolFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetMinipoolFromPubkey(validatorPubkey []byte) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetMinipoolFromPubkey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetMinipoolFromPubkey is a free data retrieval call binding the contract method 0x6dc6b1ec.
//
// Solidity: function getMinipoolFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetMinipoolFromPubkey(validatorPubkey []byte) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetMinipoolFromPubkey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetNodeAddressFromMinipool is a free data retrieval call binding the contract method 0x1dc14943.
//
// Solidity: function getNodeAddressFromMinipool(address minipool) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetNodeAddressFromMinipool(opts *bind.CallOpts, minipool common.Address) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getNodeAddressFromMinipool", minipool)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeAddressFromMinipool is a free data retrieval call binding the contract method 0x1dc14943.
//
// Solidity: function getNodeAddressFromMinipool(address minipool) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetNodeAddressFromMinipool(minipool common.Address) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetNodeAddressFromMinipool(&_Rocketminipoolregistry.CallOpts, minipool)
}

// GetNodeAddressFromMinipool is a free data retrieval call binding the contract method 0x1dc14943.
//
// Solidity: function getNodeAddressFromMinipool(address minipool) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetNodeAddressFromMinipool(minipool common.Address) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetNodeAddressFromMinipool(&_Rocketminipoolregistry.CallOpts, minipool)
}

// GetNodeAddressFromPubkey is a free data retrieval call binding the contract method 0x8eabd1ee.
//
// Solidity: function getNodeAddressFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetNodeAddressFromPubkey(opts *bind.CallOpts, validatorPubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getNodeAddressFromPubkey", validatorPubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeAddressFromPubkey is a free data retrieval call binding the contract method 0x8eabd1ee.
//
// Solidity: function getNodeAddressFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetNodeAddressFromPubkey(validatorPubkey []byte) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetNodeAddressFromPubkey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetNodeAddressFromPubkey is a free data retrieval call binding the contract method 0x8eabd1ee.
//
// Solidity: function getNodeAddressFromPubkey(bytes validatorPubkey) view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetNodeAddressFromPubkey(validatorPubkey []byte) (common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetNodeAddressFromPubkey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetValidOperatorsForKey is a free data retrieval call binding the contract method 0xa259a713.
//
// Solidity: function getValidOperatorsForKey(bytes validatorPubkey) view returns(address, address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetValidOperatorsForKey(opts *bind.CallOpts, validatorPubkey []byte) (common.Address, common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getValidOperatorsForKey", validatorPubkey)

	if err != nil {
		return *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return out0, out1, err

}

// GetValidOperatorsForKey is a free data retrieval call binding the contract method 0xa259a713.
//
// Solidity: function getValidOperatorsForKey(bytes validatorPubkey) view returns(address, address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetValidOperatorsForKey(validatorPubkey []byte) (common.Address, common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetValidOperatorsForKey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetValidOperatorsForKey is a free data retrieval call binding the contract method 0xa259a713.
//
// Solidity: function getValidOperatorsForKey(bytes validatorPubkey) view returns(address, address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetValidOperatorsForKey(validatorPubkey []byte) (common.Address, common.Address, error) {
	return _Rocketminipoolregistry.Contract.GetValidOperatorsForKey(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,uint64,uint64))
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) GetValidatorRegInfo(opts *bind.CallOpts, valPubKey []byte) (IRocketMinipoolRegistryValidatorRegistration, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "getValidatorRegInfo", valPubKey)

	if err != nil {
		return *new(IRocketMinipoolRegistryValidatorRegistration), err
	}

	out0 := *abi.ConvertType(out[0], new(IRocketMinipoolRegistryValidatorRegistration)).(*IRocketMinipoolRegistryValidatorRegistration)

	return out0, err

}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,uint64,uint64))
func (_Rocketminipoolregistry *RocketminipoolregistrySession) GetValidatorRegInfo(valPubKey []byte) (IRocketMinipoolRegistryValidatorRegistration, error) {
	return _Rocketminipoolregistry.Contract.GetValidatorRegInfo(&_Rocketminipoolregistry.CallOpts, valPubKey)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,uint64,uint64))
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) GetValidatorRegInfo(valPubKey []byte) (IRocketMinipoolRegistryValidatorRegistration, error) {
	return _Rocketminipoolregistry.Contract.GetValidatorRegInfo(&_Rocketminipoolregistry.CallOpts, valPubKey)
}

// IsMinipoolActive is a free data retrieval call binding the contract method 0x59d48d83.
//
// Solidity: function isMinipoolActive(address minipool) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) IsMinipoolActive(opts *bind.CallOpts, minipool common.Address) (bool, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "isMinipoolActive", minipool)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMinipoolActive is a free data retrieval call binding the contract method 0x59d48d83.
//
// Solidity: function isMinipoolActive(address minipool) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) IsMinipoolActive(minipool common.Address) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsMinipoolActive(&_Rocketminipoolregistry.CallOpts, minipool)
}

// IsMinipoolActive is a free data retrieval call binding the contract method 0x59d48d83.
//
// Solidity: function isMinipoolActive(address minipool) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) IsMinipoolActive(minipool common.Address) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsMinipoolActive(&_Rocketminipoolregistry.CallOpts, minipool)
}

// IsOperatorValidForKey is a free data retrieval call binding the contract method 0x5ac0191f.
//
// Solidity: function isOperatorValidForKey(address operator, bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) IsOperatorValidForKey(opts *bind.CallOpts, operator common.Address, validatorPubkey []byte) (bool, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "isOperatorValidForKey", operator, validatorPubkey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperatorValidForKey is a free data retrieval call binding the contract method 0x5ac0191f.
//
// Solidity: function isOperatorValidForKey(address operator, bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) IsOperatorValidForKey(operator common.Address, validatorPubkey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsOperatorValidForKey(&_Rocketminipoolregistry.CallOpts, operator, validatorPubkey)
}

// IsOperatorValidForKey is a free data retrieval call binding the contract method 0x5ac0191f.
//
// Solidity: function isOperatorValidForKey(address operator, bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) IsOperatorValidForKey(operator common.Address, validatorPubkey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsOperatorValidForKey(&_Rocketminipoolregistry.CallOpts, operator, validatorPubkey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) IsValidatorOptedIn(opts *bind.CallOpts, valPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "isValidatorOptedIn", valPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsValidatorOptedIn(&_Rocketminipoolregistry.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsValidatorOptedIn(&_Rocketminipoolregistry.CallOpts, valPubKey)
}

// IsValidatorRegistered is a free data retrieval call binding the contract method 0x0c5642d3.
//
// Solidity: function isValidatorRegistered(bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) IsValidatorRegistered(opts *bind.CallOpts, validatorPubkey []byte) (bool, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "isValidatorRegistered", validatorPubkey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorRegistered is a free data retrieval call binding the contract method 0x0c5642d3.
//
// Solidity: function isValidatorRegistered(bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) IsValidatorRegistered(validatorPubkey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsValidatorRegistered(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// IsValidatorRegistered is a free data retrieval call binding the contract method 0x0c5642d3.
//
// Solidity: function isValidatorRegistered(bytes validatorPubkey) view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) IsValidatorRegistered(validatorPubkey []byte) (bool, error) {
	return _Rocketminipoolregistry.Contract.IsValidatorRegistered(&_Rocketminipoolregistry.CallOpts, validatorPubkey)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Owner() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.Owner(&_Rocketminipoolregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) Owner() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.Owner(&_Rocketminipoolregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Paused() (bool, error) {
	return _Rocketminipoolregistry.Contract.Paused(&_Rocketminipoolregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) Paused() (bool, error) {
	return _Rocketminipoolregistry.Contract.Paused(&_Rocketminipoolregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) PendingOwner() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.PendingOwner(&_Rocketminipoolregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) PendingOwner() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.PendingOwner(&_Rocketminipoolregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) ProxiableUUID() ([32]byte, error) {
	return _Rocketminipoolregistry.Contract.ProxiableUUID(&_Rocketminipoolregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Rocketminipoolregistry.Contract.ProxiableUUID(&_Rocketminipoolregistry.CallOpts)
}

// RocketStorage is a free data retrieval call binding the contract method 0x67601a8e.
//
// Solidity: function rocketStorage() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) RocketStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "rocketStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RocketStorage is a free data retrieval call binding the contract method 0x67601a8e.
//
// Solidity: function rocketStorage() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) RocketStorage() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.RocketStorage(&_Rocketminipoolregistry.CallOpts)
}

// RocketStorage is a free data retrieval call binding the contract method 0x67601a8e.
//
// Solidity: function rocketStorage() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) RocketStorage() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.RocketStorage(&_Rocketminipoolregistry.CallOpts)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) UnfreezeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "unfreezeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) UnfreezeFee() (*big.Int, error) {
	return _Rocketminipoolregistry.Contract.UnfreezeFee(&_Rocketminipoolregistry.CallOpts)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) UnfreezeFee() (*big.Int, error) {
	return _Rocketminipoolregistry.Contract.UnfreezeFee(&_Rocketminipoolregistry.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) UnfreezeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "unfreezeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) UnfreezeReceiver() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.UnfreezeReceiver(&_Rocketminipoolregistry.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) UnfreezeReceiver() (common.Address, error) {
	return _Rocketminipoolregistry.Contract.UnfreezeReceiver(&_Rocketminipoolregistry.CallOpts)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, uint64 deregTimestamp, uint64 freezeTimestamp)
func (_Rocketminipoolregistry *RocketminipoolregistryCaller) ValidatorRegistrations(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists          bool
	DeregTimestamp  uint64
	FreezeTimestamp uint64
}, error) {
	var out []interface{}
	err := _Rocketminipoolregistry.contract.Call(opts, &out, "validatorRegistrations", arg0)

	outstruct := new(struct {
		Exists          bool
		DeregTimestamp  uint64
		FreezeTimestamp uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.DeregTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.FreezeTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, uint64 deregTimestamp, uint64 freezeTimestamp)
func (_Rocketminipoolregistry *RocketminipoolregistrySession) ValidatorRegistrations(arg0 []byte) (struct {
	Exists          bool
	DeregTimestamp  uint64
	FreezeTimestamp uint64
}, error) {
	return _Rocketminipoolregistry.Contract.ValidatorRegistrations(&_Rocketminipoolregistry.CallOpts, arg0)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, uint64 deregTimestamp, uint64 freezeTimestamp)
func (_Rocketminipoolregistry *RocketminipoolregistryCallerSession) ValidatorRegistrations(arg0 []byte) (struct {
	Exists          bool
	DeregTimestamp  uint64
	FreezeTimestamp uint64
}, error) {
	return _Rocketminipoolregistry.Contract.ValidatorRegistrations(&_Rocketminipoolregistry.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.AcceptOwnership(&_Rocketminipoolregistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.AcceptOwnership(&_Rocketminipoolregistry.TransactOpts)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) DeregisterValidators(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "deregisterValidators", valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.DeregisterValidators(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.DeregisterValidators(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Freeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "freeze", valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Freeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Freeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Initialize is a paid mutator transaction binding the contract method 0x3de2598c.
//
// Solidity: function initialize(address owner, address freezeOracle, address unfreezeReceiver, address rocketStorage, uint256 unfreezeFee, uint64 deregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Initialize(opts *bind.TransactOpts, owner common.Address, freezeOracle common.Address, unfreezeReceiver common.Address, rocketStorage common.Address, unfreezeFee *big.Int, deregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "initialize", owner, freezeOracle, unfreezeReceiver, rocketStorage, unfreezeFee, deregistrationPeriod)
}

// Initialize is a paid mutator transaction binding the contract method 0x3de2598c.
//
// Solidity: function initialize(address owner, address freezeOracle, address unfreezeReceiver, address rocketStorage, uint256 unfreezeFee, uint64 deregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Initialize(owner common.Address, freezeOracle common.Address, unfreezeReceiver common.Address, rocketStorage common.Address, unfreezeFee *big.Int, deregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Initialize(&_Rocketminipoolregistry.TransactOpts, owner, freezeOracle, unfreezeReceiver, rocketStorage, unfreezeFee, deregistrationPeriod)
}

// Initialize is a paid mutator transaction binding the contract method 0x3de2598c.
//
// Solidity: function initialize(address owner, address freezeOracle, address unfreezeReceiver, address rocketStorage, uint256 unfreezeFee, uint64 deregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Initialize(owner common.Address, freezeOracle common.Address, unfreezeReceiver common.Address, rocketStorage common.Address, unfreezeFee *big.Int, deregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Initialize(&_Rocketminipoolregistry.TransactOpts, owner, freezeOracle, unfreezeReceiver, rocketStorage, unfreezeFee, deregistrationPeriod)
}

// OwnerUnfreeze is a paid mutator transaction binding the contract method 0x3ba2b274.
//
// Solidity: function ownerUnfreeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) OwnerUnfreeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "ownerUnfreeze", valPubKeys)
}

// OwnerUnfreeze is a paid mutator transaction binding the contract method 0x3ba2b274.
//
// Solidity: function ownerUnfreeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) OwnerUnfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.OwnerUnfreeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// OwnerUnfreeze is a paid mutator transaction binding the contract method 0x3ba2b274.
//
// Solidity: function ownerUnfreeze(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) OwnerUnfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.OwnerUnfreeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Pause() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Pause(&_Rocketminipoolregistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Pause() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Pause(&_Rocketminipoolregistry.TransactOpts)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0xdbd739ad.
//
// Solidity: function registerValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) RegisterValidators(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "registerValidators", valPubKeys)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0xdbd739ad.
//
// Solidity: function registerValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) RegisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RegisterValidators(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0xdbd739ad.
//
// Solidity: function registerValidators(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) RegisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RegisterValidators(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RenounceOwnership(&_Rocketminipoolregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RenounceOwnership(&_Rocketminipoolregistry.TransactOpts)
}

// RequestValidatorDeregistration is a paid mutator transaction binding the contract method 0xb0d5445c.
//
// Solidity: function requestValidatorDeregistration(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) RequestValidatorDeregistration(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "requestValidatorDeregistration", valPubKeys)
}

// RequestValidatorDeregistration is a paid mutator transaction binding the contract method 0xb0d5445c.
//
// Solidity: function requestValidatorDeregistration(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) RequestValidatorDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RequestValidatorDeregistration(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// RequestValidatorDeregistration is a paid mutator transaction binding the contract method 0xb0d5445c.
//
// Solidity: function requestValidatorDeregistration(bytes[] valPubKeys) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) RequestValidatorDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.RequestValidatorDeregistration(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// SetDeregistrationPeriod is a paid mutator transaction binding the contract method 0xaaa47b39.
//
// Solidity: function setDeregistrationPeriod(uint64 newDeregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) SetDeregistrationPeriod(opts *bind.TransactOpts, newDeregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "setDeregistrationPeriod", newDeregistrationPeriod)
}

// SetDeregistrationPeriod is a paid mutator transaction binding the contract method 0xaaa47b39.
//
// Solidity: function setDeregistrationPeriod(uint64 newDeregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) SetDeregistrationPeriod(newDeregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetDeregistrationPeriod(&_Rocketminipoolregistry.TransactOpts, newDeregistrationPeriod)
}

// SetDeregistrationPeriod is a paid mutator transaction binding the contract method 0xaaa47b39.
//
// Solidity: function setDeregistrationPeriod(uint64 newDeregistrationPeriod) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) SetDeregistrationPeriod(newDeregistrationPeriod uint64) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetDeregistrationPeriod(&_Rocketminipoolregistry.TransactOpts, newDeregistrationPeriod)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address newFreezeOracle) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) SetFreezeOracle(opts *bind.TransactOpts, newFreezeOracle common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "setFreezeOracle", newFreezeOracle)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address newFreezeOracle) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) SetFreezeOracle(newFreezeOracle common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetFreezeOracle(&_Rocketminipoolregistry.TransactOpts, newFreezeOracle)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address newFreezeOracle) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) SetFreezeOracle(newFreezeOracle common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetFreezeOracle(&_Rocketminipoolregistry.TransactOpts, newFreezeOracle)
}

// SetRocketStorage is a paid mutator transaction binding the contract method 0x3af9b8ff.
//
// Solidity: function setRocketStorage(address newRocketStorage) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) SetRocketStorage(opts *bind.TransactOpts, newRocketStorage common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "setRocketStorage", newRocketStorage)
}

// SetRocketStorage is a paid mutator transaction binding the contract method 0x3af9b8ff.
//
// Solidity: function setRocketStorage(address newRocketStorage) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) SetRocketStorage(newRocketStorage common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetRocketStorage(&_Rocketminipoolregistry.TransactOpts, newRocketStorage)
}

// SetRocketStorage is a paid mutator transaction binding the contract method 0x3af9b8ff.
//
// Solidity: function setRocketStorage(address newRocketStorage) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) SetRocketStorage(newRocketStorage common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetRocketStorage(&_Rocketminipoolregistry.TransactOpts, newRocketStorage)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 newUnfreezeFee) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) SetUnfreezeFee(opts *bind.TransactOpts, newUnfreezeFee *big.Int) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "setUnfreezeFee", newUnfreezeFee)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 newUnfreezeFee) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) SetUnfreezeFee(newUnfreezeFee *big.Int) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetUnfreezeFee(&_Rocketminipoolregistry.TransactOpts, newUnfreezeFee)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 newUnfreezeFee) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) SetUnfreezeFee(newUnfreezeFee *big.Int) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetUnfreezeFee(&_Rocketminipoolregistry.TransactOpts, newUnfreezeFee)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address newUnfreezeReceiver) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) SetUnfreezeReceiver(opts *bind.TransactOpts, newUnfreezeReceiver common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "setUnfreezeReceiver", newUnfreezeReceiver)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address newUnfreezeReceiver) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) SetUnfreezeReceiver(newUnfreezeReceiver common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetUnfreezeReceiver(&_Rocketminipoolregistry.TransactOpts, newUnfreezeReceiver)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address newUnfreezeReceiver) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) SetUnfreezeReceiver(newUnfreezeReceiver common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.SetUnfreezeReceiver(&_Rocketminipoolregistry.TransactOpts, newUnfreezeReceiver)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.TransferOwnership(&_Rocketminipoolregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.TransferOwnership(&_Rocketminipoolregistry.TransactOpts, newOwner)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Unfreeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "unfreeze", valPubKeys)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Unfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Unfreeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Unfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Unfreeze(&_Rocketminipoolregistry.TransactOpts, valPubKeys)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Unpause() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Unpause(&_Rocketminipoolregistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Unpause(&_Rocketminipoolregistry.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.UpgradeToAndCall(&_Rocketminipoolregistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.UpgradeToAndCall(&_Rocketminipoolregistry.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Fallback(&_Rocketminipoolregistry.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Fallback(&_Rocketminipoolregistry.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rocketminipoolregistry.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistrySession) Receive() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Receive(&_Rocketminipoolregistry.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rocketminipoolregistry *RocketminipoolregistryTransactorSession) Receive() (*types.Transaction, error) {
	return _Rocketminipoolregistry.Contract.Receive(&_Rocketminipoolregistry.TransactOpts)
}

// RocketminipoolregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryInitializedIterator struct {
	Event *RocketminipoolregistryInitialized // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryInitialized)
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
		it.Event = new(RocketminipoolregistryInitialized)
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
func (it *RocketminipoolregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryInitialized represents a Initialized event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*RocketminipoolregistryInitializedIterator, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryInitializedIterator{contract: _Rocketminipoolregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryInitialized)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseInitialized(log types.Log) (*RocketminipoolregistryInitialized, error) {
	event := new(RocketminipoolregistryInitialized)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryOwnershipTransferStartedIterator struct {
	Event *RocketminipoolregistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryOwnershipTransferStarted)
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
		it.Event = new(RocketminipoolregistryOwnershipTransferStarted)
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
func (it *RocketminipoolregistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RocketminipoolregistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryOwnershipTransferStartedIterator{contract: _Rocketminipoolregistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryOwnershipTransferStarted)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*RocketminipoolregistryOwnershipTransferStarted, error) {
	event := new(RocketminipoolregistryOwnershipTransferStarted)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryOwnershipTransferredIterator struct {
	Event *RocketminipoolregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryOwnershipTransferred)
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
		it.Event = new(RocketminipoolregistryOwnershipTransferred)
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
func (it *RocketminipoolregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RocketminipoolregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryOwnershipTransferredIterator{contract: _Rocketminipoolregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryOwnershipTransferred)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseOwnershipTransferred(log types.Log) (*RocketminipoolregistryOwnershipTransferred, error) {
	event := new(RocketminipoolregistryOwnershipTransferred)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryPausedIterator struct {
	Event *RocketminipoolregistryPaused // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryPaused)
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
		it.Event = new(RocketminipoolregistryPaused)
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
func (it *RocketminipoolregistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryPaused represents a Paused event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterPaused(opts *bind.FilterOpts) (*RocketminipoolregistryPausedIterator, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryPausedIterator{contract: _Rocketminipoolregistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryPaused) (event.Subscription, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryPaused)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParsePaused(log types.Log) (*RocketminipoolregistryPaused, error) {
	event := new(RocketminipoolregistryPaused)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryUnpausedIterator struct {
	Event *RocketminipoolregistryUnpaused // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryUnpaused)
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
		it.Event = new(RocketminipoolregistryUnpaused)
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
func (it *RocketminipoolregistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryUnpaused represents a Unpaused event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RocketminipoolregistryUnpausedIterator, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryUnpausedIterator{contract: _Rocketminipoolregistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryUnpaused)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseUnpaused(log types.Log) (*RocketminipoolregistryUnpaused, error) {
	event := new(RocketminipoolregistryUnpaused)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryUpgradedIterator struct {
	Event *RocketminipoolregistryUpgraded // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryUpgraded)
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
		it.Event = new(RocketminipoolregistryUpgraded)
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
func (it *RocketminipoolregistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryUpgraded represents a Upgraded event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RocketminipoolregistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryUpgradedIterator{contract: _Rocketminipoolregistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryUpgraded)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseUpgraded(log types.Log) (*RocketminipoolregistryUpgraded, error) {
	event := new(RocketminipoolregistryUpgraded)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryValidatorDeregisteredIterator is returned from FilterValidatorDeregistered and is used to iterate over the raw logs and unpacked data for ValidatorDeregistered events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorDeregisteredIterator struct {
	Event *RocketminipoolregistryValidatorDeregistered // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryValidatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryValidatorDeregistered)
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
		it.Event = new(RocketminipoolregistryValidatorDeregistered)
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
func (it *RocketminipoolregistryValidatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryValidatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryValidatorDeregistered represents a ValidatorDeregistered event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorDeregistered struct {
	ValidatorPubKey common.Hash
	NodeAddress     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistered is a free log retrieval operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterValidatorDeregistered(opts *bind.FilterOpts, validatorPubKey [][]byte, nodeAddress []common.Address) (*RocketminipoolregistryValidatorDeregisteredIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "ValidatorDeregistered", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryValidatorDeregisteredIterator{contract: _Rocketminipoolregistry.contract, event: "ValidatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistered is a free log subscription operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchValidatorDeregistered(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryValidatorDeregistered, validatorPubKey [][]byte, nodeAddress []common.Address) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "ValidatorDeregistered", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryValidatorDeregistered)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
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

// ParseValidatorDeregistered is a log parse operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseValidatorDeregistered(log types.Log) (*RocketminipoolregistryValidatorDeregistered, error) {
	event := new(RocketminipoolregistryValidatorDeregistered)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryValidatorDeregistrationRequestedIterator is returned from FilterValidatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for ValidatorDeregistrationRequested events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorDeregistrationRequestedIterator struct {
	Event *RocketminipoolregistryValidatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryValidatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryValidatorDeregistrationRequested)
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
		it.Event = new(RocketminipoolregistryValidatorDeregistrationRequested)
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
func (it *RocketminipoolregistryValidatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryValidatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryValidatorDeregistrationRequested represents a ValidatorDeregistrationRequested event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorDeregistrationRequested struct {
	ValidatorPubKey common.Hash
	NodeAddress     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterValidatorDeregistrationRequested(opts *bind.FilterOpts, validatorPubKey [][]byte, nodeAddress []common.Address) (*RocketminipoolregistryValidatorDeregistrationRequestedIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "ValidatorDeregistrationRequested", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryValidatorDeregistrationRequestedIterator{contract: _Rocketminipoolregistry.contract, event: "ValidatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistrationRequested is a free log subscription operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchValidatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryValidatorDeregistrationRequested, validatorPubKey [][]byte, nodeAddress []common.Address) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "ValidatorDeregistrationRequested", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryValidatorDeregistrationRequested)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
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

// ParseValidatorDeregistrationRequested is a log parse operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseValidatorDeregistrationRequested(log types.Log) (*RocketminipoolregistryValidatorDeregistrationRequested, error) {
	event := new(RocketminipoolregistryValidatorDeregistrationRequested)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryValidatorFrozenIterator is returned from FilterValidatorFrozen and is used to iterate over the raw logs and unpacked data for ValidatorFrozen events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorFrozenIterator struct {
	Event *RocketminipoolregistryValidatorFrozen // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryValidatorFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryValidatorFrozen)
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
		it.Event = new(RocketminipoolregistryValidatorFrozen)
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
func (it *RocketminipoolregistryValidatorFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryValidatorFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryValidatorFrozen represents a ValidatorFrozen event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorFrozen struct {
	ValidatorPubKey common.Hash
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorFrozen is a free log retrieval operation binding the contract event 0xfa1d47945b5949504666dd8477bbfc00a77b976fb2119961c96adf939e53e876.
//
// Solidity: event ValidatorFrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterValidatorFrozen(opts *bind.FilterOpts, validatorPubKey [][]byte) (*RocketminipoolregistryValidatorFrozenIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "ValidatorFrozen", validatorPubKeyRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryValidatorFrozenIterator{contract: _Rocketminipoolregistry.contract, event: "ValidatorFrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorFrozen is a free log subscription operation binding the contract event 0xfa1d47945b5949504666dd8477bbfc00a77b976fb2119961c96adf939e53e876.
//
// Solidity: event ValidatorFrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchValidatorFrozen(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryValidatorFrozen, validatorPubKey [][]byte) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "ValidatorFrozen", validatorPubKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryValidatorFrozen)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
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

// ParseValidatorFrozen is a log parse operation binding the contract event 0xfa1d47945b5949504666dd8477bbfc00a77b976fb2119961c96adf939e53e876.
//
// Solidity: event ValidatorFrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseValidatorFrozen(log types.Log) (*RocketminipoolregistryValidatorFrozen, error) {
	event := new(RocketminipoolregistryValidatorFrozen)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryValidatorRegisteredIterator is returned from FilterValidatorRegistered and is used to iterate over the raw logs and unpacked data for ValidatorRegistered events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorRegisteredIterator struct {
	Event *RocketminipoolregistryValidatorRegistered // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryValidatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryValidatorRegistered)
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
		it.Event = new(RocketminipoolregistryValidatorRegistered)
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
func (it *RocketminipoolregistryValidatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryValidatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryValidatorRegistered represents a ValidatorRegistered event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorRegistered struct {
	ValidatorPubKey common.Hash
	NodeAddress     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorRegistered is a free log retrieval operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterValidatorRegistered(opts *bind.FilterOpts, validatorPubKey [][]byte, nodeAddress []common.Address) (*RocketminipoolregistryValidatorRegisteredIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "ValidatorRegistered", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryValidatorRegisteredIterator{contract: _Rocketminipoolregistry.contract, event: "ValidatorRegistered", logs: logs, sub: sub}, nil
}

// WatchValidatorRegistered is a free log subscription operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchValidatorRegistered(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryValidatorRegistered, validatorPubKey [][]byte, nodeAddress []common.Address) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}
	var nodeAddressRule []interface{}
	for _, nodeAddressItem := range nodeAddress {
		nodeAddressRule = append(nodeAddressRule, nodeAddressItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "ValidatorRegistered", validatorPubKeyRule, nodeAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryValidatorRegistered)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
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

// ParseValidatorRegistered is a log parse operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes indexed validatorPubKey, address indexed nodeAddress)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseValidatorRegistered(log types.Log) (*RocketminipoolregistryValidatorRegistered, error) {
	event := new(RocketminipoolregistryValidatorRegistered)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RocketminipoolregistryValidatorUnfrozenIterator is returned from FilterValidatorUnfrozen and is used to iterate over the raw logs and unpacked data for ValidatorUnfrozen events raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorUnfrozenIterator struct {
	Event *RocketminipoolregistryValidatorUnfrozen // Event containing the contract specifics and raw log

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
func (it *RocketminipoolregistryValidatorUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RocketminipoolregistryValidatorUnfrozen)
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
		it.Event = new(RocketminipoolregistryValidatorUnfrozen)
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
func (it *RocketminipoolregistryValidatorUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RocketminipoolregistryValidatorUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RocketminipoolregistryValidatorUnfrozen represents a ValidatorUnfrozen event raised by the Rocketminipoolregistry contract.
type RocketminipoolregistryValidatorUnfrozen struct {
	ValidatorPubKey common.Hash
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorUnfrozen is a free log retrieval operation binding the contract event 0x83f1169418932171cc1b130d2c1e543ace8af55b47e4a1341b39b9c71d491392.
//
// Solidity: event ValidatorUnfrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) FilterValidatorUnfrozen(opts *bind.FilterOpts, validatorPubKey [][]byte) (*RocketminipoolregistryValidatorUnfrozenIterator, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.FilterLogs(opts, "ValidatorUnfrozen", validatorPubKeyRule)
	if err != nil {
		return nil, err
	}
	return &RocketminipoolregistryValidatorUnfrozenIterator{contract: _Rocketminipoolregistry.contract, event: "ValidatorUnfrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorUnfrozen is a free log subscription operation binding the contract event 0x83f1169418932171cc1b130d2c1e543ace8af55b47e4a1341b39b9c71d491392.
//
// Solidity: event ValidatorUnfrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) WatchValidatorUnfrozen(opts *bind.WatchOpts, sink chan<- *RocketminipoolregistryValidatorUnfrozen, validatorPubKey [][]byte) (event.Subscription, error) {

	var validatorPubKeyRule []interface{}
	for _, validatorPubKeyItem := range validatorPubKey {
		validatorPubKeyRule = append(validatorPubKeyRule, validatorPubKeyItem)
	}

	logs, sub, err := _Rocketminipoolregistry.contract.WatchLogs(opts, "ValidatorUnfrozen", validatorPubKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RocketminipoolregistryValidatorUnfrozen)
				if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
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

// ParseValidatorUnfrozen is a log parse operation binding the contract event 0x83f1169418932171cc1b130d2c1e543ace8af55b47e4a1341b39b9c71d491392.
//
// Solidity: event ValidatorUnfrozen(bytes indexed validatorPubKey)
func (_Rocketminipoolregistry *RocketminipoolregistryFilterer) ParseValidatorUnfrozen(log types.Log) (*RocketminipoolregistryValidatorUnfrozen, error) {
	event := new(RocketminipoolregistryValidatorUnfrozen)
	if err := _Rocketminipoolregistry.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
