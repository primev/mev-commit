// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rewardmanager

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

// RewardmanagerMetaData contains all meta data concerning the Rewardmanager contract.
var RewardmanagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"autoClaim\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"enabled\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"autoClaimBlacklist\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"blacklisted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"autoClaimGasLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimOrphanedRewards\",\"inputs\":[{\"name\":\"pubkeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"toPay\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"disableAutoClaim\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"enableAutoClaim\",\"inputs\":[{\"name\":\"claimExistingRewards\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"vanillaRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mevCommitAVS\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mevCommitMiddleware\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"autoClaimGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"orphanedRewards\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"overrideAddresses\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"overrideAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"overrideReceiver\",\"inputs\":[{\"name\":\"overrideAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"migrateExistingRewards\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"payProposer\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeFromAutoClaimBlacklist\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeOverrideAddress\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAutoClaimGasLimit\",\"inputs\":[{\"name\":\"autoClaimGasLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMevCommitAVS\",\"inputs\":[{\"name\":\"mevCommitAVS\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMevCommitMiddleware\",\"inputs\":[{\"name\":\"mevCommitMiddleware\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVanillaRegistry\",\"inputs\":[{\"name\":\"vanillaRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unclaimedRewards\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AutoClaimDisabled\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AutoClaimEnabled\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AutoClaimGasLimitSet\",\"inputs\":[{\"name\":\"autoClaimGasLimit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AutoClaimTransferFailed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"toPay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AutoClaimed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"toPay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MevCommitAVSSet\",\"inputs\":[{\"name\":\"newMevCommitAVS\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MevCommitMiddlewareSet\",\"inputs\":[{\"name\":\"newMevCommitMiddleware\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NoRewards\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrphanedRewardsAccumulated\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrphanedRewardsClaimed\",\"inputs\":[{\"name\":\"toPay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OverrideAddressRemoved\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OverrideAddressSet\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"overrideAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PaymentStored\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"toPay\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposerNotFound\",\"inputs\":[{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":true,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemovedFromAutoClaimBlacklist\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsMigrated\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VanillaRegistrySet\",\"inputs\":[{\"name\":\"newVanillaRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAutoClaimGasLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expectedLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoEthPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoOrphanedRewards\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoOverriddenAddressToRemove\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OrphanedRewardsClaimFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RewardsClaimFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// RewardmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use RewardmanagerMetaData.ABI instead.
var RewardmanagerABI = RewardmanagerMetaData.ABI

// Rewardmanager is an auto generated Go binding around an Ethereum contract.
type Rewardmanager struct {
	RewardmanagerCaller     // Read-only binding to the contract
	RewardmanagerTransactor // Write-only binding to the contract
	RewardmanagerFilterer   // Log filterer for contract events
}

// RewardmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardmanagerSession struct {
	Contract     *Rewardmanager    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardmanagerCallerSession struct {
	Contract *RewardmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RewardmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardmanagerTransactorSession struct {
	Contract     *RewardmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RewardmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardmanagerRaw struct {
	Contract *Rewardmanager // Generic contract binding to access the raw methods on
}

// RewardmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardmanagerCallerRaw struct {
	Contract *RewardmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// RewardmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardmanagerTransactorRaw struct {
	Contract *RewardmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardmanager creates a new instance of Rewardmanager, bound to a specific deployed contract.
func NewRewardmanager(address common.Address, backend bind.ContractBackend) (*Rewardmanager, error) {
	contract, err := bindRewardmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rewardmanager{RewardmanagerCaller: RewardmanagerCaller{contract: contract}, RewardmanagerTransactor: RewardmanagerTransactor{contract: contract}, RewardmanagerFilterer: RewardmanagerFilterer{contract: contract}}, nil
}

// NewRewardmanagerCaller creates a new read-only instance of Rewardmanager, bound to a specific deployed contract.
func NewRewardmanagerCaller(address common.Address, caller bind.ContractCaller) (*RewardmanagerCaller, error) {
	contract, err := bindRewardmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerCaller{contract: contract}, nil
}

// NewRewardmanagerTransactor creates a new write-only instance of Rewardmanager, bound to a specific deployed contract.
func NewRewardmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardmanagerTransactor, error) {
	contract, err := bindRewardmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerTransactor{contract: contract}, nil
}

// NewRewardmanagerFilterer creates a new log filterer instance of Rewardmanager, bound to a specific deployed contract.
func NewRewardmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardmanagerFilterer, error) {
	contract, err := bindRewardmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerFilterer{contract: contract}, nil
}

// bindRewardmanager binds a generic wrapper to an already deployed contract.
func bindRewardmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewardmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewardmanager *RewardmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rewardmanager.Contract.RewardmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewardmanager *RewardmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.Contract.RewardmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewardmanager *RewardmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewardmanager.Contract.RewardmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewardmanager *RewardmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rewardmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewardmanager *RewardmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewardmanager *RewardmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewardmanager.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewardmanager *RewardmanagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewardmanager *RewardmanagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rewardmanager.Contract.UPGRADEINTERFACEVERSION(&_Rewardmanager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewardmanager *RewardmanagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rewardmanager.Contract.UPGRADEINTERFACEVERSION(&_Rewardmanager.CallOpts)
}

// AutoClaim is a free data retrieval call binding the contract method 0x637d7035.
//
// Solidity: function autoClaim(address receiver) view returns(bool enabled)
func (_Rewardmanager *RewardmanagerCaller) AutoClaim(opts *bind.CallOpts, receiver common.Address) (bool, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "autoClaim", receiver)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AutoClaim is a free data retrieval call binding the contract method 0x637d7035.
//
// Solidity: function autoClaim(address receiver) view returns(bool enabled)
func (_Rewardmanager *RewardmanagerSession) AutoClaim(receiver common.Address) (bool, error) {
	return _Rewardmanager.Contract.AutoClaim(&_Rewardmanager.CallOpts, receiver)
}

// AutoClaim is a free data retrieval call binding the contract method 0x637d7035.
//
// Solidity: function autoClaim(address receiver) view returns(bool enabled)
func (_Rewardmanager *RewardmanagerCallerSession) AutoClaim(receiver common.Address) (bool, error) {
	return _Rewardmanager.Contract.AutoClaim(&_Rewardmanager.CallOpts, receiver)
}

// AutoClaimBlacklist is a free data retrieval call binding the contract method 0xeb0b73f4.
//
// Solidity: function autoClaimBlacklist(address receiver) view returns(bool blacklisted)
func (_Rewardmanager *RewardmanagerCaller) AutoClaimBlacklist(opts *bind.CallOpts, receiver common.Address) (bool, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "autoClaimBlacklist", receiver)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AutoClaimBlacklist is a free data retrieval call binding the contract method 0xeb0b73f4.
//
// Solidity: function autoClaimBlacklist(address receiver) view returns(bool blacklisted)
func (_Rewardmanager *RewardmanagerSession) AutoClaimBlacklist(receiver common.Address) (bool, error) {
	return _Rewardmanager.Contract.AutoClaimBlacklist(&_Rewardmanager.CallOpts, receiver)
}

// AutoClaimBlacklist is a free data retrieval call binding the contract method 0xeb0b73f4.
//
// Solidity: function autoClaimBlacklist(address receiver) view returns(bool blacklisted)
func (_Rewardmanager *RewardmanagerCallerSession) AutoClaimBlacklist(receiver common.Address) (bool, error) {
	return _Rewardmanager.Contract.AutoClaimBlacklist(&_Rewardmanager.CallOpts, receiver)
}

// AutoClaimGasLimit is a free data retrieval call binding the contract method 0xe492535a.
//
// Solidity: function autoClaimGasLimit() view returns(uint256)
func (_Rewardmanager *RewardmanagerCaller) AutoClaimGasLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "autoClaimGasLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AutoClaimGasLimit is a free data retrieval call binding the contract method 0xe492535a.
//
// Solidity: function autoClaimGasLimit() view returns(uint256)
func (_Rewardmanager *RewardmanagerSession) AutoClaimGasLimit() (*big.Int, error) {
	return _Rewardmanager.Contract.AutoClaimGasLimit(&_Rewardmanager.CallOpts)
}

// AutoClaimGasLimit is a free data retrieval call binding the contract method 0xe492535a.
//
// Solidity: function autoClaimGasLimit() view returns(uint256)
func (_Rewardmanager *RewardmanagerCallerSession) AutoClaimGasLimit() (*big.Int, error) {
	return _Rewardmanager.Contract.AutoClaimGasLimit(&_Rewardmanager.CallOpts)
}

// OrphanedRewards is a free data retrieval call binding the contract method 0xf82af456.
//
// Solidity: function orphanedRewards(bytes pubkey) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerCaller) OrphanedRewards(opts *bind.CallOpts, pubkey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "orphanedRewards", pubkey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OrphanedRewards is a free data retrieval call binding the contract method 0xf82af456.
//
// Solidity: function orphanedRewards(bytes pubkey) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerSession) OrphanedRewards(pubkey []byte) (*big.Int, error) {
	return _Rewardmanager.Contract.OrphanedRewards(&_Rewardmanager.CallOpts, pubkey)
}

// OrphanedRewards is a free data retrieval call binding the contract method 0xf82af456.
//
// Solidity: function orphanedRewards(bytes pubkey) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerCallerSession) OrphanedRewards(pubkey []byte) (*big.Int, error) {
	return _Rewardmanager.Contract.OrphanedRewards(&_Rewardmanager.CallOpts, pubkey)
}

// OverrideAddresses is a free data retrieval call binding the contract method 0x2136822b.
//
// Solidity: function overrideAddresses(address receiver) view returns(address overrideAddress)
func (_Rewardmanager *RewardmanagerCaller) OverrideAddresses(opts *bind.CallOpts, receiver common.Address) (common.Address, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "overrideAddresses", receiver)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OverrideAddresses is a free data retrieval call binding the contract method 0x2136822b.
//
// Solidity: function overrideAddresses(address receiver) view returns(address overrideAddress)
func (_Rewardmanager *RewardmanagerSession) OverrideAddresses(receiver common.Address) (common.Address, error) {
	return _Rewardmanager.Contract.OverrideAddresses(&_Rewardmanager.CallOpts, receiver)
}

// OverrideAddresses is a free data retrieval call binding the contract method 0x2136822b.
//
// Solidity: function overrideAddresses(address receiver) view returns(address overrideAddress)
func (_Rewardmanager *RewardmanagerCallerSession) OverrideAddresses(receiver common.Address) (common.Address, error) {
	return _Rewardmanager.Contract.OverrideAddresses(&_Rewardmanager.CallOpts, receiver)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewardmanager *RewardmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewardmanager *RewardmanagerSession) Owner() (common.Address, error) {
	return _Rewardmanager.Contract.Owner(&_Rewardmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewardmanager *RewardmanagerCallerSession) Owner() (common.Address, error) {
	return _Rewardmanager.Contract.Owner(&_Rewardmanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewardmanager *RewardmanagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewardmanager *RewardmanagerSession) Paused() (bool, error) {
	return _Rewardmanager.Contract.Paused(&_Rewardmanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewardmanager *RewardmanagerCallerSession) Paused() (bool, error) {
	return _Rewardmanager.Contract.Paused(&_Rewardmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewardmanager *RewardmanagerCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewardmanager *RewardmanagerSession) PendingOwner() (common.Address, error) {
	return _Rewardmanager.Contract.PendingOwner(&_Rewardmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewardmanager *RewardmanagerCallerSession) PendingOwner() (common.Address, error) {
	return _Rewardmanager.Contract.PendingOwner(&_Rewardmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewardmanager *RewardmanagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewardmanager *RewardmanagerSession) ProxiableUUID() ([32]byte, error) {
	return _Rewardmanager.Contract.ProxiableUUID(&_Rewardmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewardmanager *RewardmanagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Rewardmanager.Contract.ProxiableUUID(&_Rewardmanager.CallOpts)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address addr) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerCaller) UnclaimedRewards(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Rewardmanager.contract.Call(opts, &out, "unclaimedRewards", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address addr) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerSession) UnclaimedRewards(addr common.Address) (*big.Int, error) {
	return _Rewardmanager.Contract.UnclaimedRewards(&_Rewardmanager.CallOpts, addr)
}

// UnclaimedRewards is a free data retrieval call binding the contract method 0x949813b8.
//
// Solidity: function unclaimedRewards(address addr) view returns(uint256 amount)
func (_Rewardmanager *RewardmanagerCallerSession) UnclaimedRewards(addr common.Address) (*big.Int, error) {
	return _Rewardmanager.Contract.UnclaimedRewards(&_Rewardmanager.CallOpts, addr)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewardmanager *RewardmanagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewardmanager *RewardmanagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Rewardmanager.Contract.AcceptOwnership(&_Rewardmanager.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Rewardmanager.Contract.AcceptOwnership(&_Rewardmanager.TransactOpts)
}

// ClaimOrphanedRewards is a paid mutator transaction binding the contract method 0x8440e996.
//
// Solidity: function claimOrphanedRewards(bytes[] pubkeys, address toPay) returns()
func (_Rewardmanager *RewardmanagerTransactor) ClaimOrphanedRewards(opts *bind.TransactOpts, pubkeys [][]byte, toPay common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "claimOrphanedRewards", pubkeys, toPay)
}

// ClaimOrphanedRewards is a paid mutator transaction binding the contract method 0x8440e996.
//
// Solidity: function claimOrphanedRewards(bytes[] pubkeys, address toPay) returns()
func (_Rewardmanager *RewardmanagerSession) ClaimOrphanedRewards(pubkeys [][]byte, toPay common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.ClaimOrphanedRewards(&_Rewardmanager.TransactOpts, pubkeys, toPay)
}

// ClaimOrphanedRewards is a paid mutator transaction binding the contract method 0x8440e996.
//
// Solidity: function claimOrphanedRewards(bytes[] pubkeys, address toPay) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) ClaimOrphanedRewards(pubkeys [][]byte, toPay common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.ClaimOrphanedRewards(&_Rewardmanager.TransactOpts, pubkeys, toPay)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Rewardmanager *RewardmanagerTransactor) ClaimRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "claimRewards")
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Rewardmanager *RewardmanagerSession) ClaimRewards() (*types.Transaction, error) {
	return _Rewardmanager.Contract.ClaimRewards(&_Rewardmanager.TransactOpts)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x372500ab.
//
// Solidity: function claimRewards() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) ClaimRewards() (*types.Transaction, error) {
	return _Rewardmanager.Contract.ClaimRewards(&_Rewardmanager.TransactOpts)
}

// DisableAutoClaim is a paid mutator transaction binding the contract method 0x038217ab.
//
// Solidity: function disableAutoClaim() returns()
func (_Rewardmanager *RewardmanagerTransactor) DisableAutoClaim(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "disableAutoClaim")
}

// DisableAutoClaim is a paid mutator transaction binding the contract method 0x038217ab.
//
// Solidity: function disableAutoClaim() returns()
func (_Rewardmanager *RewardmanagerSession) DisableAutoClaim() (*types.Transaction, error) {
	return _Rewardmanager.Contract.DisableAutoClaim(&_Rewardmanager.TransactOpts)
}

// DisableAutoClaim is a paid mutator transaction binding the contract method 0x038217ab.
//
// Solidity: function disableAutoClaim() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) DisableAutoClaim() (*types.Transaction, error) {
	return _Rewardmanager.Contract.DisableAutoClaim(&_Rewardmanager.TransactOpts)
}

// EnableAutoClaim is a paid mutator transaction binding the contract method 0x0621037f.
//
// Solidity: function enableAutoClaim(bool claimExistingRewards) returns()
func (_Rewardmanager *RewardmanagerTransactor) EnableAutoClaim(opts *bind.TransactOpts, claimExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "enableAutoClaim", claimExistingRewards)
}

// EnableAutoClaim is a paid mutator transaction binding the contract method 0x0621037f.
//
// Solidity: function enableAutoClaim(bool claimExistingRewards) returns()
func (_Rewardmanager *RewardmanagerSession) EnableAutoClaim(claimExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.Contract.EnableAutoClaim(&_Rewardmanager.TransactOpts, claimExistingRewards)
}

// EnableAutoClaim is a paid mutator transaction binding the contract method 0x0621037f.
//
// Solidity: function enableAutoClaim(bool claimExistingRewards) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) EnableAutoClaim(claimExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.Contract.EnableAutoClaim(&_Rewardmanager.TransactOpts, claimExistingRewards)
}

// Initialize is a paid mutator transaction binding the contract method 0x530b97a4.
//
// Solidity: function initialize(address vanillaRegistry, address mevCommitAVS, address mevCommitMiddleware, uint256 autoClaimGasLimit, address owner) returns()
func (_Rewardmanager *RewardmanagerTransactor) Initialize(opts *bind.TransactOpts, vanillaRegistry common.Address, mevCommitAVS common.Address, mevCommitMiddleware common.Address, autoClaimGasLimit *big.Int, owner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "initialize", vanillaRegistry, mevCommitAVS, mevCommitMiddleware, autoClaimGasLimit, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x530b97a4.
//
// Solidity: function initialize(address vanillaRegistry, address mevCommitAVS, address mevCommitMiddleware, uint256 autoClaimGasLimit, address owner) returns()
func (_Rewardmanager *RewardmanagerSession) Initialize(vanillaRegistry common.Address, mevCommitAVS common.Address, mevCommitMiddleware common.Address, autoClaimGasLimit *big.Int, owner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.Initialize(&_Rewardmanager.TransactOpts, vanillaRegistry, mevCommitAVS, mevCommitMiddleware, autoClaimGasLimit, owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x530b97a4.
//
// Solidity: function initialize(address vanillaRegistry, address mevCommitAVS, address mevCommitMiddleware, uint256 autoClaimGasLimit, address owner) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) Initialize(vanillaRegistry common.Address, mevCommitAVS common.Address, mevCommitMiddleware common.Address, autoClaimGasLimit *big.Int, owner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.Initialize(&_Rewardmanager.TransactOpts, vanillaRegistry, mevCommitAVS, mevCommitMiddleware, autoClaimGasLimit, owner)
}

// OverrideReceiver is a paid mutator transaction binding the contract method 0x53e55e15.
//
// Solidity: function overrideReceiver(address overrideAddress, bool migrateExistingRewards) returns()
func (_Rewardmanager *RewardmanagerTransactor) OverrideReceiver(opts *bind.TransactOpts, overrideAddress common.Address, migrateExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "overrideReceiver", overrideAddress, migrateExistingRewards)
}

// OverrideReceiver is a paid mutator transaction binding the contract method 0x53e55e15.
//
// Solidity: function overrideReceiver(address overrideAddress, bool migrateExistingRewards) returns()
func (_Rewardmanager *RewardmanagerSession) OverrideReceiver(overrideAddress common.Address, migrateExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.Contract.OverrideReceiver(&_Rewardmanager.TransactOpts, overrideAddress, migrateExistingRewards)
}

// OverrideReceiver is a paid mutator transaction binding the contract method 0x53e55e15.
//
// Solidity: function overrideReceiver(address overrideAddress, bool migrateExistingRewards) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) OverrideReceiver(overrideAddress common.Address, migrateExistingRewards bool) (*types.Transaction, error) {
	return _Rewardmanager.Contract.OverrideReceiver(&_Rewardmanager.TransactOpts, overrideAddress, migrateExistingRewards)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewardmanager *RewardmanagerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewardmanager *RewardmanagerSession) Pause() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Pause(&_Rewardmanager.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) Pause() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Pause(&_Rewardmanager.TransactOpts)
}

// PayProposer is a paid mutator transaction binding the contract method 0x41d8d551.
//
// Solidity: function payProposer(bytes pubkey) payable returns()
func (_Rewardmanager *RewardmanagerTransactor) PayProposer(opts *bind.TransactOpts, pubkey []byte) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "payProposer", pubkey)
}

// PayProposer is a paid mutator transaction binding the contract method 0x41d8d551.
//
// Solidity: function payProposer(bytes pubkey) payable returns()
func (_Rewardmanager *RewardmanagerSession) PayProposer(pubkey []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.PayProposer(&_Rewardmanager.TransactOpts, pubkey)
}

// PayProposer is a paid mutator transaction binding the contract method 0x41d8d551.
//
// Solidity: function payProposer(bytes pubkey) payable returns()
func (_Rewardmanager *RewardmanagerTransactorSession) PayProposer(pubkey []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.PayProposer(&_Rewardmanager.TransactOpts, pubkey)
}

// RemoveFromAutoClaimBlacklist is a paid mutator transaction binding the contract method 0x30683883.
//
// Solidity: function removeFromAutoClaimBlacklist(address addr) returns()
func (_Rewardmanager *RewardmanagerTransactor) RemoveFromAutoClaimBlacklist(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "removeFromAutoClaimBlacklist", addr)
}

// RemoveFromAutoClaimBlacklist is a paid mutator transaction binding the contract method 0x30683883.
//
// Solidity: function removeFromAutoClaimBlacklist(address addr) returns()
func (_Rewardmanager *RewardmanagerSession) RemoveFromAutoClaimBlacklist(addr common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.RemoveFromAutoClaimBlacklist(&_Rewardmanager.TransactOpts, addr)
}

// RemoveFromAutoClaimBlacklist is a paid mutator transaction binding the contract method 0x30683883.
//
// Solidity: function removeFromAutoClaimBlacklist(address addr) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) RemoveFromAutoClaimBlacklist(addr common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.RemoveFromAutoClaimBlacklist(&_Rewardmanager.TransactOpts, addr)
}

// RemoveOverrideAddress is a paid mutator transaction binding the contract method 0xe74eda82.
//
// Solidity: function removeOverrideAddress() returns()
func (_Rewardmanager *RewardmanagerTransactor) RemoveOverrideAddress(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "removeOverrideAddress")
}

// RemoveOverrideAddress is a paid mutator transaction binding the contract method 0xe74eda82.
//
// Solidity: function removeOverrideAddress() returns()
func (_Rewardmanager *RewardmanagerSession) RemoveOverrideAddress() (*types.Transaction, error) {
	return _Rewardmanager.Contract.RemoveOverrideAddress(&_Rewardmanager.TransactOpts)
}

// RemoveOverrideAddress is a paid mutator transaction binding the contract method 0xe74eda82.
//
// Solidity: function removeOverrideAddress() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) RemoveOverrideAddress() (*types.Transaction, error) {
	return _Rewardmanager.Contract.RemoveOverrideAddress(&_Rewardmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewardmanager *RewardmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewardmanager *RewardmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewardmanager.Contract.RenounceOwnership(&_Rewardmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewardmanager.Contract.RenounceOwnership(&_Rewardmanager.TransactOpts)
}

// SetAutoClaimGasLimit is a paid mutator transaction binding the contract method 0xee0def12.
//
// Solidity: function setAutoClaimGasLimit(uint256 autoClaimGasLimit) returns()
func (_Rewardmanager *RewardmanagerTransactor) SetAutoClaimGasLimit(opts *bind.TransactOpts, autoClaimGasLimit *big.Int) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "setAutoClaimGasLimit", autoClaimGasLimit)
}

// SetAutoClaimGasLimit is a paid mutator transaction binding the contract method 0xee0def12.
//
// Solidity: function setAutoClaimGasLimit(uint256 autoClaimGasLimit) returns()
func (_Rewardmanager *RewardmanagerSession) SetAutoClaimGasLimit(autoClaimGasLimit *big.Int) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetAutoClaimGasLimit(&_Rewardmanager.TransactOpts, autoClaimGasLimit)
}

// SetAutoClaimGasLimit is a paid mutator transaction binding the contract method 0xee0def12.
//
// Solidity: function setAutoClaimGasLimit(uint256 autoClaimGasLimit) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) SetAutoClaimGasLimit(autoClaimGasLimit *big.Int) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetAutoClaimGasLimit(&_Rewardmanager.TransactOpts, autoClaimGasLimit)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address mevCommitAVS) returns()
func (_Rewardmanager *RewardmanagerTransactor) SetMevCommitAVS(opts *bind.TransactOpts, mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "setMevCommitAVS", mevCommitAVS)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address mevCommitAVS) returns()
func (_Rewardmanager *RewardmanagerSession) SetMevCommitAVS(mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetMevCommitAVS(&_Rewardmanager.TransactOpts, mevCommitAVS)
}

// SetMevCommitAVS is a paid mutator transaction binding the contract method 0x18df7701.
//
// Solidity: function setMevCommitAVS(address mevCommitAVS) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) SetMevCommitAVS(mevCommitAVS common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetMevCommitAVS(&_Rewardmanager.TransactOpts, mevCommitAVS)
}

// SetMevCommitMiddleware is a paid mutator transaction binding the contract method 0xab746a6a.
//
// Solidity: function setMevCommitMiddleware(address mevCommitMiddleware) returns()
func (_Rewardmanager *RewardmanagerTransactor) SetMevCommitMiddleware(opts *bind.TransactOpts, mevCommitMiddleware common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "setMevCommitMiddleware", mevCommitMiddleware)
}

// SetMevCommitMiddleware is a paid mutator transaction binding the contract method 0xab746a6a.
//
// Solidity: function setMevCommitMiddleware(address mevCommitMiddleware) returns()
func (_Rewardmanager *RewardmanagerSession) SetMevCommitMiddleware(mevCommitMiddleware common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetMevCommitMiddleware(&_Rewardmanager.TransactOpts, mevCommitMiddleware)
}

// SetMevCommitMiddleware is a paid mutator transaction binding the contract method 0xab746a6a.
//
// Solidity: function setMevCommitMiddleware(address mevCommitMiddleware) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) SetMevCommitMiddleware(mevCommitMiddleware common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetMevCommitMiddleware(&_Rewardmanager.TransactOpts, mevCommitMiddleware)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address vanillaRegistry) returns()
func (_Rewardmanager *RewardmanagerTransactor) SetVanillaRegistry(opts *bind.TransactOpts, vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "setVanillaRegistry", vanillaRegistry)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address vanillaRegistry) returns()
func (_Rewardmanager *RewardmanagerSession) SetVanillaRegistry(vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetVanillaRegistry(&_Rewardmanager.TransactOpts, vanillaRegistry)
}

// SetVanillaRegistry is a paid mutator transaction binding the contract method 0xf99a9e82.
//
// Solidity: function setVanillaRegistry(address vanillaRegistry) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) SetVanillaRegistry(vanillaRegistry common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.SetVanillaRegistry(&_Rewardmanager.TransactOpts, vanillaRegistry)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewardmanager *RewardmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewardmanager *RewardmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.TransferOwnership(&_Rewardmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewardmanager *RewardmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewardmanager.Contract.TransferOwnership(&_Rewardmanager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewardmanager *RewardmanagerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewardmanager *RewardmanagerSession) Unpause() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Unpause(&_Rewardmanager.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewardmanager *RewardmanagerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Unpause(&_Rewardmanager.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewardmanager *RewardmanagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewardmanager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewardmanager *RewardmanagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.UpgradeToAndCall(&_Rewardmanager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewardmanager *RewardmanagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.UpgradeToAndCall(&_Rewardmanager.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewardmanager *RewardmanagerTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Rewardmanager.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewardmanager *RewardmanagerSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.Fallback(&_Rewardmanager.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewardmanager *RewardmanagerTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rewardmanager.Contract.Fallback(&_Rewardmanager.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewardmanager *RewardmanagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewardmanager.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewardmanager *RewardmanagerSession) Receive() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Receive(&_Rewardmanager.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewardmanager *RewardmanagerTransactorSession) Receive() (*types.Transaction, error) {
	return _Rewardmanager.Contract.Receive(&_Rewardmanager.TransactOpts)
}

// RewardmanagerAutoClaimDisabledIterator is returned from FilterAutoClaimDisabled and is used to iterate over the raw logs and unpacked data for AutoClaimDisabled events raised by the Rewardmanager contract.
type RewardmanagerAutoClaimDisabledIterator struct {
	Event *RewardmanagerAutoClaimDisabled // Event containing the contract specifics and raw log

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
func (it *RewardmanagerAutoClaimDisabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerAutoClaimDisabled)
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
		it.Event = new(RewardmanagerAutoClaimDisabled)
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
func (it *RewardmanagerAutoClaimDisabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerAutoClaimDisabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerAutoClaimDisabled represents a AutoClaimDisabled event raised by the Rewardmanager contract.
type RewardmanagerAutoClaimDisabled struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAutoClaimDisabled is a free log retrieval operation binding the contract event 0xa6faef39bcfc2889eb92a5bf9f422d21e549c0d94676fe4fd879b78dd529e46b.
//
// Solidity: event AutoClaimDisabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) FilterAutoClaimDisabled(opts *bind.FilterOpts, receiver []common.Address) (*RewardmanagerAutoClaimDisabledIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "AutoClaimDisabled", receiverRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerAutoClaimDisabledIterator{contract: _Rewardmanager.contract, event: "AutoClaimDisabled", logs: logs, sub: sub}, nil
}

// WatchAutoClaimDisabled is a free log subscription operation binding the contract event 0xa6faef39bcfc2889eb92a5bf9f422d21e549c0d94676fe4fd879b78dd529e46b.
//
// Solidity: event AutoClaimDisabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) WatchAutoClaimDisabled(opts *bind.WatchOpts, sink chan<- *RewardmanagerAutoClaimDisabled, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "AutoClaimDisabled", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerAutoClaimDisabled)
				if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimDisabled", log); err != nil {
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

// ParseAutoClaimDisabled is a log parse operation binding the contract event 0xa6faef39bcfc2889eb92a5bf9f422d21e549c0d94676fe4fd879b78dd529e46b.
//
// Solidity: event AutoClaimDisabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) ParseAutoClaimDisabled(log types.Log) (*RewardmanagerAutoClaimDisabled, error) {
	event := new(RewardmanagerAutoClaimDisabled)
	if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimDisabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerAutoClaimEnabledIterator is returned from FilterAutoClaimEnabled and is used to iterate over the raw logs and unpacked data for AutoClaimEnabled events raised by the Rewardmanager contract.
type RewardmanagerAutoClaimEnabledIterator struct {
	Event *RewardmanagerAutoClaimEnabled // Event containing the contract specifics and raw log

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
func (it *RewardmanagerAutoClaimEnabledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerAutoClaimEnabled)
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
		it.Event = new(RewardmanagerAutoClaimEnabled)
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
func (it *RewardmanagerAutoClaimEnabledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerAutoClaimEnabledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerAutoClaimEnabled represents a AutoClaimEnabled event raised by the Rewardmanager contract.
type RewardmanagerAutoClaimEnabled struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAutoClaimEnabled is a free log retrieval operation binding the contract event 0xd1f03cfccc13e9f522e155a39beed6bd9b54c9cbb956c81d56f9a0e93d3a09e9.
//
// Solidity: event AutoClaimEnabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) FilterAutoClaimEnabled(opts *bind.FilterOpts, receiver []common.Address) (*RewardmanagerAutoClaimEnabledIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "AutoClaimEnabled", receiverRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerAutoClaimEnabledIterator{contract: _Rewardmanager.contract, event: "AutoClaimEnabled", logs: logs, sub: sub}, nil
}

// WatchAutoClaimEnabled is a free log subscription operation binding the contract event 0xd1f03cfccc13e9f522e155a39beed6bd9b54c9cbb956c81d56f9a0e93d3a09e9.
//
// Solidity: event AutoClaimEnabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) WatchAutoClaimEnabled(opts *bind.WatchOpts, sink chan<- *RewardmanagerAutoClaimEnabled, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "AutoClaimEnabled", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerAutoClaimEnabled)
				if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimEnabled", log); err != nil {
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

// ParseAutoClaimEnabled is a log parse operation binding the contract event 0xd1f03cfccc13e9f522e155a39beed6bd9b54c9cbb956c81d56f9a0e93d3a09e9.
//
// Solidity: event AutoClaimEnabled(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) ParseAutoClaimEnabled(log types.Log) (*RewardmanagerAutoClaimEnabled, error) {
	event := new(RewardmanagerAutoClaimEnabled)
	if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimEnabled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerAutoClaimGasLimitSetIterator is returned from FilterAutoClaimGasLimitSet and is used to iterate over the raw logs and unpacked data for AutoClaimGasLimitSet events raised by the Rewardmanager contract.
type RewardmanagerAutoClaimGasLimitSetIterator struct {
	Event *RewardmanagerAutoClaimGasLimitSet // Event containing the contract specifics and raw log

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
func (it *RewardmanagerAutoClaimGasLimitSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerAutoClaimGasLimitSet)
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
		it.Event = new(RewardmanagerAutoClaimGasLimitSet)
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
func (it *RewardmanagerAutoClaimGasLimitSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerAutoClaimGasLimitSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerAutoClaimGasLimitSet represents a AutoClaimGasLimitSet event raised by the Rewardmanager contract.
type RewardmanagerAutoClaimGasLimitSet struct {
	AutoClaimGasLimit *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterAutoClaimGasLimitSet is a free log retrieval operation binding the contract event 0x1fb681c70a3a546070a58f4686cb626a8578d62c030bc2ca5b0e92d61b4a1f08.
//
// Solidity: event AutoClaimGasLimitSet(uint256 autoClaimGasLimit)
func (_Rewardmanager *RewardmanagerFilterer) FilterAutoClaimGasLimitSet(opts *bind.FilterOpts) (*RewardmanagerAutoClaimGasLimitSetIterator, error) {

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "AutoClaimGasLimitSet")
	if err != nil {
		return nil, err
	}
	return &RewardmanagerAutoClaimGasLimitSetIterator{contract: _Rewardmanager.contract, event: "AutoClaimGasLimitSet", logs: logs, sub: sub}, nil
}

// WatchAutoClaimGasLimitSet is a free log subscription operation binding the contract event 0x1fb681c70a3a546070a58f4686cb626a8578d62c030bc2ca5b0e92d61b4a1f08.
//
// Solidity: event AutoClaimGasLimitSet(uint256 autoClaimGasLimit)
func (_Rewardmanager *RewardmanagerFilterer) WatchAutoClaimGasLimitSet(opts *bind.WatchOpts, sink chan<- *RewardmanagerAutoClaimGasLimitSet) (event.Subscription, error) {

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "AutoClaimGasLimitSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerAutoClaimGasLimitSet)
				if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimGasLimitSet", log); err != nil {
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

// ParseAutoClaimGasLimitSet is a log parse operation binding the contract event 0x1fb681c70a3a546070a58f4686cb626a8578d62c030bc2ca5b0e92d61b4a1f08.
//
// Solidity: event AutoClaimGasLimitSet(uint256 autoClaimGasLimit)
func (_Rewardmanager *RewardmanagerFilterer) ParseAutoClaimGasLimitSet(log types.Log) (*RewardmanagerAutoClaimGasLimitSet, error) {
	event := new(RewardmanagerAutoClaimGasLimitSet)
	if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimGasLimitSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerAutoClaimTransferFailedIterator is returned from FilterAutoClaimTransferFailed and is used to iterate over the raw logs and unpacked data for AutoClaimTransferFailed events raised by the Rewardmanager contract.
type RewardmanagerAutoClaimTransferFailedIterator struct {
	Event *RewardmanagerAutoClaimTransferFailed // Event containing the contract specifics and raw log

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
func (it *RewardmanagerAutoClaimTransferFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerAutoClaimTransferFailed)
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
		it.Event = new(RewardmanagerAutoClaimTransferFailed)
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
func (it *RewardmanagerAutoClaimTransferFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerAutoClaimTransferFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerAutoClaimTransferFailed represents a AutoClaimTransferFailed event raised by the Rewardmanager contract.
type RewardmanagerAutoClaimTransferFailed struct {
	Provider common.Address
	Receiver common.Address
	ToPay    common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAutoClaimTransferFailed is a free log retrieval operation binding the contract event 0x02f32ab06b518743bd01da217d5c0e5f5bc0f3115395cd18bd4b3aecb5f0f0a7.
//
// Solidity: event AutoClaimTransferFailed(address indexed provider, address indexed receiver, address indexed toPay)
func (_Rewardmanager *RewardmanagerFilterer) FilterAutoClaimTransferFailed(opts *bind.FilterOpts, provider []common.Address, receiver []common.Address, toPay []common.Address) (*RewardmanagerAutoClaimTransferFailedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "AutoClaimTransferFailed", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerAutoClaimTransferFailedIterator{contract: _Rewardmanager.contract, event: "AutoClaimTransferFailed", logs: logs, sub: sub}, nil
}

// WatchAutoClaimTransferFailed is a free log subscription operation binding the contract event 0x02f32ab06b518743bd01da217d5c0e5f5bc0f3115395cd18bd4b3aecb5f0f0a7.
//
// Solidity: event AutoClaimTransferFailed(address indexed provider, address indexed receiver, address indexed toPay)
func (_Rewardmanager *RewardmanagerFilterer) WatchAutoClaimTransferFailed(opts *bind.WatchOpts, sink chan<- *RewardmanagerAutoClaimTransferFailed, provider []common.Address, receiver []common.Address, toPay []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "AutoClaimTransferFailed", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerAutoClaimTransferFailed)
				if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimTransferFailed", log); err != nil {
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

// ParseAutoClaimTransferFailed is a log parse operation binding the contract event 0x02f32ab06b518743bd01da217d5c0e5f5bc0f3115395cd18bd4b3aecb5f0f0a7.
//
// Solidity: event AutoClaimTransferFailed(address indexed provider, address indexed receiver, address indexed toPay)
func (_Rewardmanager *RewardmanagerFilterer) ParseAutoClaimTransferFailed(log types.Log) (*RewardmanagerAutoClaimTransferFailed, error) {
	event := new(RewardmanagerAutoClaimTransferFailed)
	if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimTransferFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerAutoClaimedIterator is returned from FilterAutoClaimed and is used to iterate over the raw logs and unpacked data for AutoClaimed events raised by the Rewardmanager contract.
type RewardmanagerAutoClaimedIterator struct {
	Event *RewardmanagerAutoClaimed // Event containing the contract specifics and raw log

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
func (it *RewardmanagerAutoClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerAutoClaimed)
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
		it.Event = new(RewardmanagerAutoClaimed)
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
func (it *RewardmanagerAutoClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerAutoClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerAutoClaimed represents a AutoClaimed event raised by the Rewardmanager contract.
type RewardmanagerAutoClaimed struct {
	Provider common.Address
	Receiver common.Address
	ToPay    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAutoClaimed is a free log retrieval operation binding the contract event 0x55e13ae964e89156459daa4d62db464e276deb21a2e14d48d585a2901ccdb67c.
//
// Solidity: event AutoClaimed(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterAutoClaimed(opts *bind.FilterOpts, provider []common.Address, receiver []common.Address, toPay []common.Address) (*RewardmanagerAutoClaimedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "AutoClaimed", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerAutoClaimedIterator{contract: _Rewardmanager.contract, event: "AutoClaimed", logs: logs, sub: sub}, nil
}

// WatchAutoClaimed is a free log subscription operation binding the contract event 0x55e13ae964e89156459daa4d62db464e276deb21a2e14d48d585a2901ccdb67c.
//
// Solidity: event AutoClaimed(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchAutoClaimed(opts *bind.WatchOpts, sink chan<- *RewardmanagerAutoClaimed, provider []common.Address, receiver []common.Address, toPay []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "AutoClaimed", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerAutoClaimed)
				if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimed", log); err != nil {
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

// ParseAutoClaimed is a log parse operation binding the contract event 0x55e13ae964e89156459daa4d62db464e276deb21a2e14d48d585a2901ccdb67c.
//
// Solidity: event AutoClaimed(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParseAutoClaimed(log types.Log) (*RewardmanagerAutoClaimed, error) {
	event := new(RewardmanagerAutoClaimed)
	if err := _Rewardmanager.contract.UnpackLog(event, "AutoClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Rewardmanager contract.
type RewardmanagerInitializedIterator struct {
	Event *RewardmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *RewardmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerInitialized)
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
		it.Event = new(RewardmanagerInitialized)
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
func (it *RewardmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerInitialized represents a Initialized event raised by the Rewardmanager contract.
type RewardmanagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rewardmanager *RewardmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*RewardmanagerInitializedIterator, error) {

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RewardmanagerInitializedIterator{contract: _Rewardmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rewardmanager *RewardmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RewardmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerInitialized)
				if err := _Rewardmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParseInitialized(log types.Log) (*RewardmanagerInitialized, error) {
	event := new(RewardmanagerInitialized)
	if err := _Rewardmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerMevCommitAVSSetIterator is returned from FilterMevCommitAVSSet and is used to iterate over the raw logs and unpacked data for MevCommitAVSSet events raised by the Rewardmanager contract.
type RewardmanagerMevCommitAVSSetIterator struct {
	Event *RewardmanagerMevCommitAVSSet // Event containing the contract specifics and raw log

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
func (it *RewardmanagerMevCommitAVSSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerMevCommitAVSSet)
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
		it.Event = new(RewardmanagerMevCommitAVSSet)
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
func (it *RewardmanagerMevCommitAVSSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerMevCommitAVSSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerMevCommitAVSSet represents a MevCommitAVSSet event raised by the Rewardmanager contract.
type RewardmanagerMevCommitAVSSet struct {
	NewMevCommitAVS common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMevCommitAVSSet is a free log retrieval operation binding the contract event 0xea31ef8e67c313130e3f091a4d8405fae9a0438ff1790f754299e50f4507cb9e.
//
// Solidity: event MevCommitAVSSet(address indexed newMevCommitAVS)
func (_Rewardmanager *RewardmanagerFilterer) FilterMevCommitAVSSet(opts *bind.FilterOpts, newMevCommitAVS []common.Address) (*RewardmanagerMevCommitAVSSetIterator, error) {

	var newMevCommitAVSRule []interface{}
	for _, newMevCommitAVSItem := range newMevCommitAVS {
		newMevCommitAVSRule = append(newMevCommitAVSRule, newMevCommitAVSItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "MevCommitAVSSet", newMevCommitAVSRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerMevCommitAVSSetIterator{contract: _Rewardmanager.contract, event: "MevCommitAVSSet", logs: logs, sub: sub}, nil
}

// WatchMevCommitAVSSet is a free log subscription operation binding the contract event 0xea31ef8e67c313130e3f091a4d8405fae9a0438ff1790f754299e50f4507cb9e.
//
// Solidity: event MevCommitAVSSet(address indexed newMevCommitAVS)
func (_Rewardmanager *RewardmanagerFilterer) WatchMevCommitAVSSet(opts *bind.WatchOpts, sink chan<- *RewardmanagerMevCommitAVSSet, newMevCommitAVS []common.Address) (event.Subscription, error) {

	var newMevCommitAVSRule []interface{}
	for _, newMevCommitAVSItem := range newMevCommitAVS {
		newMevCommitAVSRule = append(newMevCommitAVSRule, newMevCommitAVSItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "MevCommitAVSSet", newMevCommitAVSRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerMevCommitAVSSet)
				if err := _Rewardmanager.contract.UnpackLog(event, "MevCommitAVSSet", log); err != nil {
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

// ParseMevCommitAVSSet is a log parse operation binding the contract event 0xea31ef8e67c313130e3f091a4d8405fae9a0438ff1790f754299e50f4507cb9e.
//
// Solidity: event MevCommitAVSSet(address indexed newMevCommitAVS)
func (_Rewardmanager *RewardmanagerFilterer) ParseMevCommitAVSSet(log types.Log) (*RewardmanagerMevCommitAVSSet, error) {
	event := new(RewardmanagerMevCommitAVSSet)
	if err := _Rewardmanager.contract.UnpackLog(event, "MevCommitAVSSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerMevCommitMiddlewareSetIterator is returned from FilterMevCommitMiddlewareSet and is used to iterate over the raw logs and unpacked data for MevCommitMiddlewareSet events raised by the Rewardmanager contract.
type RewardmanagerMevCommitMiddlewareSetIterator struct {
	Event *RewardmanagerMevCommitMiddlewareSet // Event containing the contract specifics and raw log

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
func (it *RewardmanagerMevCommitMiddlewareSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerMevCommitMiddlewareSet)
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
		it.Event = new(RewardmanagerMevCommitMiddlewareSet)
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
func (it *RewardmanagerMevCommitMiddlewareSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerMevCommitMiddlewareSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerMevCommitMiddlewareSet represents a MevCommitMiddlewareSet event raised by the Rewardmanager contract.
type RewardmanagerMevCommitMiddlewareSet struct {
	NewMevCommitMiddleware common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterMevCommitMiddlewareSet is a free log retrieval operation binding the contract event 0x03827e418018b6cb23d8b6fe48f4f8b5b273c0873cf65e080e64d75444642b11.
//
// Solidity: event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware)
func (_Rewardmanager *RewardmanagerFilterer) FilterMevCommitMiddlewareSet(opts *bind.FilterOpts, newMevCommitMiddleware []common.Address) (*RewardmanagerMevCommitMiddlewareSetIterator, error) {

	var newMevCommitMiddlewareRule []interface{}
	for _, newMevCommitMiddlewareItem := range newMevCommitMiddleware {
		newMevCommitMiddlewareRule = append(newMevCommitMiddlewareRule, newMevCommitMiddlewareItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "MevCommitMiddlewareSet", newMevCommitMiddlewareRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerMevCommitMiddlewareSetIterator{contract: _Rewardmanager.contract, event: "MevCommitMiddlewareSet", logs: logs, sub: sub}, nil
}

// WatchMevCommitMiddlewareSet is a free log subscription operation binding the contract event 0x03827e418018b6cb23d8b6fe48f4f8b5b273c0873cf65e080e64d75444642b11.
//
// Solidity: event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware)
func (_Rewardmanager *RewardmanagerFilterer) WatchMevCommitMiddlewareSet(opts *bind.WatchOpts, sink chan<- *RewardmanagerMevCommitMiddlewareSet, newMevCommitMiddleware []common.Address) (event.Subscription, error) {

	var newMevCommitMiddlewareRule []interface{}
	for _, newMevCommitMiddlewareItem := range newMevCommitMiddleware {
		newMevCommitMiddlewareRule = append(newMevCommitMiddlewareRule, newMevCommitMiddlewareItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "MevCommitMiddlewareSet", newMevCommitMiddlewareRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerMevCommitMiddlewareSet)
				if err := _Rewardmanager.contract.UnpackLog(event, "MevCommitMiddlewareSet", log); err != nil {
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

// ParseMevCommitMiddlewareSet is a log parse operation binding the contract event 0x03827e418018b6cb23d8b6fe48f4f8b5b273c0873cf65e080e64d75444642b11.
//
// Solidity: event MevCommitMiddlewareSet(address indexed newMevCommitMiddleware)
func (_Rewardmanager *RewardmanagerFilterer) ParseMevCommitMiddlewareSet(log types.Log) (*RewardmanagerMevCommitMiddlewareSet, error) {
	event := new(RewardmanagerMevCommitMiddlewareSet)
	if err := _Rewardmanager.contract.UnpackLog(event, "MevCommitMiddlewareSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerNoRewardsIterator is returned from FilterNoRewards and is used to iterate over the raw logs and unpacked data for NoRewards events raised by the Rewardmanager contract.
type RewardmanagerNoRewardsIterator struct {
	Event *RewardmanagerNoRewards // Event containing the contract specifics and raw log

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
func (it *RewardmanagerNoRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerNoRewards)
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
		it.Event = new(RewardmanagerNoRewards)
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
func (it *RewardmanagerNoRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerNoRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerNoRewards represents a NoRewards event raised by the Rewardmanager contract.
type RewardmanagerNoRewards struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNoRewards is a free log retrieval operation binding the contract event 0xc8fe3c53b1969860fd19335c08b61c8042e1062471fed38206696bdbb7edb280.
//
// Solidity: event NoRewards(address addr)
func (_Rewardmanager *RewardmanagerFilterer) FilterNoRewards(opts *bind.FilterOpts) (*RewardmanagerNoRewardsIterator, error) {

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "NoRewards")
	if err != nil {
		return nil, err
	}
	return &RewardmanagerNoRewardsIterator{contract: _Rewardmanager.contract, event: "NoRewards", logs: logs, sub: sub}, nil
}

// WatchNoRewards is a free log subscription operation binding the contract event 0xc8fe3c53b1969860fd19335c08b61c8042e1062471fed38206696bdbb7edb280.
//
// Solidity: event NoRewards(address addr)
func (_Rewardmanager *RewardmanagerFilterer) WatchNoRewards(opts *bind.WatchOpts, sink chan<- *RewardmanagerNoRewards) (event.Subscription, error) {

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "NoRewards")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerNoRewards)
				if err := _Rewardmanager.contract.UnpackLog(event, "NoRewards", log); err != nil {
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

// ParseNoRewards is a log parse operation binding the contract event 0xc8fe3c53b1969860fd19335c08b61c8042e1062471fed38206696bdbb7edb280.
//
// Solidity: event NoRewards(address addr)
func (_Rewardmanager *RewardmanagerFilterer) ParseNoRewards(log types.Log) (*RewardmanagerNoRewards, error) {
	event := new(RewardmanagerNoRewards)
	if err := _Rewardmanager.contract.UnpackLog(event, "NoRewards", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOrphanedRewardsAccumulatedIterator is returned from FilterOrphanedRewardsAccumulated and is used to iterate over the raw logs and unpacked data for OrphanedRewardsAccumulated events raised by the Rewardmanager contract.
type RewardmanagerOrphanedRewardsAccumulatedIterator struct {
	Event *RewardmanagerOrphanedRewardsAccumulated // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOrphanedRewardsAccumulatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOrphanedRewardsAccumulated)
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
		it.Event = new(RewardmanagerOrphanedRewardsAccumulated)
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
func (it *RewardmanagerOrphanedRewardsAccumulatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOrphanedRewardsAccumulatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOrphanedRewardsAccumulated represents a OrphanedRewardsAccumulated event raised by the Rewardmanager contract.
type RewardmanagerOrphanedRewardsAccumulated struct {
	Provider common.Address
	Pubkey   common.Hash
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOrphanedRewardsAccumulated is a free log retrieval operation binding the contract event 0x4bc16ba7a77bac0a237974dcb09883d423c0dff3dcaa061d93e9dc2bf637bc67.
//
// Solidity: event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterOrphanedRewardsAccumulated(opts *bind.FilterOpts, provider []common.Address, pubkey [][]byte) (*RewardmanagerOrphanedRewardsAccumulatedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OrphanedRewardsAccumulated", providerRule, pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOrphanedRewardsAccumulatedIterator{contract: _Rewardmanager.contract, event: "OrphanedRewardsAccumulated", logs: logs, sub: sub}, nil
}

// WatchOrphanedRewardsAccumulated is a free log subscription operation binding the contract event 0x4bc16ba7a77bac0a237974dcb09883d423c0dff3dcaa061d93e9dc2bf637bc67.
//
// Solidity: event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchOrphanedRewardsAccumulated(opts *bind.WatchOpts, sink chan<- *RewardmanagerOrphanedRewardsAccumulated, provider []common.Address, pubkey [][]byte) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OrphanedRewardsAccumulated", providerRule, pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOrphanedRewardsAccumulated)
				if err := _Rewardmanager.contract.UnpackLog(event, "OrphanedRewardsAccumulated", log); err != nil {
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

// ParseOrphanedRewardsAccumulated is a log parse operation binding the contract event 0x4bc16ba7a77bac0a237974dcb09883d423c0dff3dcaa061d93e9dc2bf637bc67.
//
// Solidity: event OrphanedRewardsAccumulated(address indexed provider, bytes indexed pubkey, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParseOrphanedRewardsAccumulated(log types.Log) (*RewardmanagerOrphanedRewardsAccumulated, error) {
	event := new(RewardmanagerOrphanedRewardsAccumulated)
	if err := _Rewardmanager.contract.UnpackLog(event, "OrphanedRewardsAccumulated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOrphanedRewardsClaimedIterator is returned from FilterOrphanedRewardsClaimed and is used to iterate over the raw logs and unpacked data for OrphanedRewardsClaimed events raised by the Rewardmanager contract.
type RewardmanagerOrphanedRewardsClaimedIterator struct {
	Event *RewardmanagerOrphanedRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOrphanedRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOrphanedRewardsClaimed)
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
		it.Event = new(RewardmanagerOrphanedRewardsClaimed)
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
func (it *RewardmanagerOrphanedRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOrphanedRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOrphanedRewardsClaimed represents a OrphanedRewardsClaimed event raised by the Rewardmanager contract.
type RewardmanagerOrphanedRewardsClaimed struct {
	ToPay  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOrphanedRewardsClaimed is a free log retrieval operation binding the contract event 0x96ea63e11bceeaac5e0271bd0d7e76bee392659c8336a88d73891d6f7a9541a4.
//
// Solidity: event OrphanedRewardsClaimed(address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterOrphanedRewardsClaimed(opts *bind.FilterOpts, toPay []common.Address) (*RewardmanagerOrphanedRewardsClaimedIterator, error) {

	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OrphanedRewardsClaimed", toPayRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOrphanedRewardsClaimedIterator{contract: _Rewardmanager.contract, event: "OrphanedRewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchOrphanedRewardsClaimed is a free log subscription operation binding the contract event 0x96ea63e11bceeaac5e0271bd0d7e76bee392659c8336a88d73891d6f7a9541a4.
//
// Solidity: event OrphanedRewardsClaimed(address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchOrphanedRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewardmanagerOrphanedRewardsClaimed, toPay []common.Address) (event.Subscription, error) {

	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OrphanedRewardsClaimed", toPayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOrphanedRewardsClaimed)
				if err := _Rewardmanager.contract.UnpackLog(event, "OrphanedRewardsClaimed", log); err != nil {
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

// ParseOrphanedRewardsClaimed is a log parse operation binding the contract event 0x96ea63e11bceeaac5e0271bd0d7e76bee392659c8336a88d73891d6f7a9541a4.
//
// Solidity: event OrphanedRewardsClaimed(address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParseOrphanedRewardsClaimed(log types.Log) (*RewardmanagerOrphanedRewardsClaimed, error) {
	event := new(RewardmanagerOrphanedRewardsClaimed)
	if err := _Rewardmanager.contract.UnpackLog(event, "OrphanedRewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOverrideAddressRemovedIterator is returned from FilterOverrideAddressRemoved and is used to iterate over the raw logs and unpacked data for OverrideAddressRemoved events raised by the Rewardmanager contract.
type RewardmanagerOverrideAddressRemovedIterator struct {
	Event *RewardmanagerOverrideAddressRemoved // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOverrideAddressRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOverrideAddressRemoved)
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
		it.Event = new(RewardmanagerOverrideAddressRemoved)
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
func (it *RewardmanagerOverrideAddressRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOverrideAddressRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOverrideAddressRemoved represents a OverrideAddressRemoved event raised by the Rewardmanager contract.
type RewardmanagerOverrideAddressRemoved struct {
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOverrideAddressRemoved is a free log retrieval operation binding the contract event 0x833597e98bf237734ce683696457e46672325e5a9410c3eb9a628b17671721da.
//
// Solidity: event OverrideAddressRemoved(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) FilterOverrideAddressRemoved(opts *bind.FilterOpts, receiver []common.Address) (*RewardmanagerOverrideAddressRemovedIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OverrideAddressRemoved", receiverRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOverrideAddressRemovedIterator{contract: _Rewardmanager.contract, event: "OverrideAddressRemoved", logs: logs, sub: sub}, nil
}

// WatchOverrideAddressRemoved is a free log subscription operation binding the contract event 0x833597e98bf237734ce683696457e46672325e5a9410c3eb9a628b17671721da.
//
// Solidity: event OverrideAddressRemoved(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) WatchOverrideAddressRemoved(opts *bind.WatchOpts, sink chan<- *RewardmanagerOverrideAddressRemoved, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OverrideAddressRemoved", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOverrideAddressRemoved)
				if err := _Rewardmanager.contract.UnpackLog(event, "OverrideAddressRemoved", log); err != nil {
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

// ParseOverrideAddressRemoved is a log parse operation binding the contract event 0x833597e98bf237734ce683696457e46672325e5a9410c3eb9a628b17671721da.
//
// Solidity: event OverrideAddressRemoved(address indexed receiver)
func (_Rewardmanager *RewardmanagerFilterer) ParseOverrideAddressRemoved(log types.Log) (*RewardmanagerOverrideAddressRemoved, error) {
	event := new(RewardmanagerOverrideAddressRemoved)
	if err := _Rewardmanager.contract.UnpackLog(event, "OverrideAddressRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOverrideAddressSetIterator is returned from FilterOverrideAddressSet and is used to iterate over the raw logs and unpacked data for OverrideAddressSet events raised by the Rewardmanager contract.
type RewardmanagerOverrideAddressSetIterator struct {
	Event *RewardmanagerOverrideAddressSet // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOverrideAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOverrideAddressSet)
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
		it.Event = new(RewardmanagerOverrideAddressSet)
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
func (it *RewardmanagerOverrideAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOverrideAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOverrideAddressSet represents a OverrideAddressSet event raised by the Rewardmanager contract.
type RewardmanagerOverrideAddressSet struct {
	Receiver        common.Address
	OverrideAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterOverrideAddressSet is a free log retrieval operation binding the contract event 0xfc48fb94734fd09fa0848abf47643a88319bc97992a6ca5af994eb62c9b39629.
//
// Solidity: event OverrideAddressSet(address indexed receiver, address indexed overrideAddress)
func (_Rewardmanager *RewardmanagerFilterer) FilterOverrideAddressSet(opts *bind.FilterOpts, receiver []common.Address, overrideAddress []common.Address) (*RewardmanagerOverrideAddressSetIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var overrideAddressRule []interface{}
	for _, overrideAddressItem := range overrideAddress {
		overrideAddressRule = append(overrideAddressRule, overrideAddressItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OverrideAddressSet", receiverRule, overrideAddressRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOverrideAddressSetIterator{contract: _Rewardmanager.contract, event: "OverrideAddressSet", logs: logs, sub: sub}, nil
}

// WatchOverrideAddressSet is a free log subscription operation binding the contract event 0xfc48fb94734fd09fa0848abf47643a88319bc97992a6ca5af994eb62c9b39629.
//
// Solidity: event OverrideAddressSet(address indexed receiver, address indexed overrideAddress)
func (_Rewardmanager *RewardmanagerFilterer) WatchOverrideAddressSet(opts *bind.WatchOpts, sink chan<- *RewardmanagerOverrideAddressSet, receiver []common.Address, overrideAddress []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var overrideAddressRule []interface{}
	for _, overrideAddressItem := range overrideAddress {
		overrideAddressRule = append(overrideAddressRule, overrideAddressItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OverrideAddressSet", receiverRule, overrideAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOverrideAddressSet)
				if err := _Rewardmanager.contract.UnpackLog(event, "OverrideAddressSet", log); err != nil {
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

// ParseOverrideAddressSet is a log parse operation binding the contract event 0xfc48fb94734fd09fa0848abf47643a88319bc97992a6ca5af994eb62c9b39629.
//
// Solidity: event OverrideAddressSet(address indexed receiver, address indexed overrideAddress)
func (_Rewardmanager *RewardmanagerFilterer) ParseOverrideAddressSet(log types.Log) (*RewardmanagerOverrideAddressSet, error) {
	event := new(RewardmanagerOverrideAddressSet)
	if err := _Rewardmanager.contract.UnpackLog(event, "OverrideAddressSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Rewardmanager contract.
type RewardmanagerOwnershipTransferStartedIterator struct {
	Event *RewardmanagerOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOwnershipTransferStarted)
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
		it.Event = new(RewardmanagerOwnershipTransferStarted)
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
func (it *RewardmanagerOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Rewardmanager contract.
type RewardmanagerOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rewardmanager *RewardmanagerFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardmanagerOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOwnershipTransferStartedIterator{contract: _Rewardmanager.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rewardmanager *RewardmanagerFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *RewardmanagerOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOwnershipTransferStarted)
				if err := _Rewardmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParseOwnershipTransferStarted(log types.Log) (*RewardmanagerOwnershipTransferStarted, error) {
	event := new(RewardmanagerOwnershipTransferStarted)
	if err := _Rewardmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rewardmanager contract.
type RewardmanagerOwnershipTransferredIterator struct {
	Event *RewardmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewardmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerOwnershipTransferred)
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
		it.Event = new(RewardmanagerOwnershipTransferred)
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
func (it *RewardmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Rewardmanager contract.
type RewardmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewardmanager *RewardmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerOwnershipTransferredIterator{contract: _Rewardmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewardmanager *RewardmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerOwnershipTransferred)
				if err := _Rewardmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParseOwnershipTransferred(log types.Log) (*RewardmanagerOwnershipTransferred, error) {
	event := new(RewardmanagerOwnershipTransferred)
	if err := _Rewardmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Rewardmanager contract.
type RewardmanagerPausedIterator struct {
	Event *RewardmanagerPaused // Event containing the contract specifics and raw log

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
func (it *RewardmanagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerPaused)
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
		it.Event = new(RewardmanagerPaused)
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
func (it *RewardmanagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerPaused represents a Paused event raised by the Rewardmanager contract.
type RewardmanagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rewardmanager *RewardmanagerFilterer) FilterPaused(opts *bind.FilterOpts) (*RewardmanagerPausedIterator, error) {

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RewardmanagerPausedIterator{contract: _Rewardmanager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rewardmanager *RewardmanagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RewardmanagerPaused) (event.Subscription, error) {

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerPaused)
				if err := _Rewardmanager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParsePaused(log types.Log) (*RewardmanagerPaused, error) {
	event := new(RewardmanagerPaused)
	if err := _Rewardmanager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerPaymentStoredIterator is returned from FilterPaymentStored and is used to iterate over the raw logs and unpacked data for PaymentStored events raised by the Rewardmanager contract.
type RewardmanagerPaymentStoredIterator struct {
	Event *RewardmanagerPaymentStored // Event containing the contract specifics and raw log

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
func (it *RewardmanagerPaymentStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerPaymentStored)
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
		it.Event = new(RewardmanagerPaymentStored)
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
func (it *RewardmanagerPaymentStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerPaymentStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerPaymentStored represents a PaymentStored event raised by the Rewardmanager contract.
type RewardmanagerPaymentStored struct {
	Provider common.Address
	Receiver common.Address
	ToPay    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPaymentStored is a free log retrieval operation binding the contract event 0x696199f3a1a686870e664b74805bb702d32adb07ea32f1812d1641469db5f88c.
//
// Solidity: event PaymentStored(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterPaymentStored(opts *bind.FilterOpts, provider []common.Address, receiver []common.Address, toPay []common.Address) (*RewardmanagerPaymentStoredIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "PaymentStored", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerPaymentStoredIterator{contract: _Rewardmanager.contract, event: "PaymentStored", logs: logs, sub: sub}, nil
}

// WatchPaymentStored is a free log subscription operation binding the contract event 0x696199f3a1a686870e664b74805bb702d32adb07ea32f1812d1641469db5f88c.
//
// Solidity: event PaymentStored(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchPaymentStored(opts *bind.WatchOpts, sink chan<- *RewardmanagerPaymentStored, provider []common.Address, receiver []common.Address, toPay []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var toPayRule []interface{}
	for _, toPayItem := range toPay {
		toPayRule = append(toPayRule, toPayItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "PaymentStored", providerRule, receiverRule, toPayRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerPaymentStored)
				if err := _Rewardmanager.contract.UnpackLog(event, "PaymentStored", log); err != nil {
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

// ParsePaymentStored is a log parse operation binding the contract event 0x696199f3a1a686870e664b74805bb702d32adb07ea32f1812d1641469db5f88c.
//
// Solidity: event PaymentStored(address indexed provider, address indexed receiver, address indexed toPay, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParsePaymentStored(log types.Log) (*RewardmanagerPaymentStored, error) {
	event := new(RewardmanagerPaymentStored)
	if err := _Rewardmanager.contract.UnpackLog(event, "PaymentStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerProposerNotFoundIterator is returned from FilterProposerNotFound and is used to iterate over the raw logs and unpacked data for ProposerNotFound events raised by the Rewardmanager contract.
type RewardmanagerProposerNotFoundIterator struct {
	Event *RewardmanagerProposerNotFound // Event containing the contract specifics and raw log

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
func (it *RewardmanagerProposerNotFoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerProposerNotFound)
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
		it.Event = new(RewardmanagerProposerNotFound)
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
func (it *RewardmanagerProposerNotFoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerProposerNotFoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerProposerNotFound represents a ProposerNotFound event raised by the Rewardmanager contract.
type RewardmanagerProposerNotFound struct {
	Pubkey common.Hash
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterProposerNotFound is a free log retrieval operation binding the contract event 0x739849996958cb570eef3c2841985e02c6fccd94b4473bbb4955a91685eae93f.
//
// Solidity: event ProposerNotFound(bytes indexed pubkey)
func (_Rewardmanager *RewardmanagerFilterer) FilterProposerNotFound(opts *bind.FilterOpts, pubkey [][]byte) (*RewardmanagerProposerNotFoundIterator, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "ProposerNotFound", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerProposerNotFoundIterator{contract: _Rewardmanager.contract, event: "ProposerNotFound", logs: logs, sub: sub}, nil
}

// WatchProposerNotFound is a free log subscription operation binding the contract event 0x739849996958cb570eef3c2841985e02c6fccd94b4473bbb4955a91685eae93f.
//
// Solidity: event ProposerNotFound(bytes indexed pubkey)
func (_Rewardmanager *RewardmanagerFilterer) WatchProposerNotFound(opts *bind.WatchOpts, sink chan<- *RewardmanagerProposerNotFound, pubkey [][]byte) (event.Subscription, error) {

	var pubkeyRule []interface{}
	for _, pubkeyItem := range pubkey {
		pubkeyRule = append(pubkeyRule, pubkeyItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "ProposerNotFound", pubkeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerProposerNotFound)
				if err := _Rewardmanager.contract.UnpackLog(event, "ProposerNotFound", log); err != nil {
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

// ParseProposerNotFound is a log parse operation binding the contract event 0x739849996958cb570eef3c2841985e02c6fccd94b4473bbb4955a91685eae93f.
//
// Solidity: event ProposerNotFound(bytes indexed pubkey)
func (_Rewardmanager *RewardmanagerFilterer) ParseProposerNotFound(log types.Log) (*RewardmanagerProposerNotFound, error) {
	event := new(RewardmanagerProposerNotFound)
	if err := _Rewardmanager.contract.UnpackLog(event, "ProposerNotFound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerRemovedFromAutoClaimBlacklistIterator is returned from FilterRemovedFromAutoClaimBlacklist and is used to iterate over the raw logs and unpacked data for RemovedFromAutoClaimBlacklist events raised by the Rewardmanager contract.
type RewardmanagerRemovedFromAutoClaimBlacklistIterator struct {
	Event *RewardmanagerRemovedFromAutoClaimBlacklist // Event containing the contract specifics and raw log

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
func (it *RewardmanagerRemovedFromAutoClaimBlacklistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerRemovedFromAutoClaimBlacklist)
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
		it.Event = new(RewardmanagerRemovedFromAutoClaimBlacklist)
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
func (it *RewardmanagerRemovedFromAutoClaimBlacklistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerRemovedFromAutoClaimBlacklistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerRemovedFromAutoClaimBlacklist represents a RemovedFromAutoClaimBlacklist event raised by the Rewardmanager contract.
type RewardmanagerRemovedFromAutoClaimBlacklist struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRemovedFromAutoClaimBlacklist is a free log retrieval operation binding the contract event 0x754b19e96c74531bb2405588a9cc003f89da91b013554d437def3897698f9326.
//
// Solidity: event RemovedFromAutoClaimBlacklist(address indexed addr)
func (_Rewardmanager *RewardmanagerFilterer) FilterRemovedFromAutoClaimBlacklist(opts *bind.FilterOpts, addr []common.Address) (*RewardmanagerRemovedFromAutoClaimBlacklistIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "RemovedFromAutoClaimBlacklist", addrRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerRemovedFromAutoClaimBlacklistIterator{contract: _Rewardmanager.contract, event: "RemovedFromAutoClaimBlacklist", logs: logs, sub: sub}, nil
}

// WatchRemovedFromAutoClaimBlacklist is a free log subscription operation binding the contract event 0x754b19e96c74531bb2405588a9cc003f89da91b013554d437def3897698f9326.
//
// Solidity: event RemovedFromAutoClaimBlacklist(address indexed addr)
func (_Rewardmanager *RewardmanagerFilterer) WatchRemovedFromAutoClaimBlacklist(opts *bind.WatchOpts, sink chan<- *RewardmanagerRemovedFromAutoClaimBlacklist, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "RemovedFromAutoClaimBlacklist", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerRemovedFromAutoClaimBlacklist)
				if err := _Rewardmanager.contract.UnpackLog(event, "RemovedFromAutoClaimBlacklist", log); err != nil {
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

// ParseRemovedFromAutoClaimBlacklist is a log parse operation binding the contract event 0x754b19e96c74531bb2405588a9cc003f89da91b013554d437def3897698f9326.
//
// Solidity: event RemovedFromAutoClaimBlacklist(address indexed addr)
func (_Rewardmanager *RewardmanagerFilterer) ParseRemovedFromAutoClaimBlacklist(log types.Log) (*RewardmanagerRemovedFromAutoClaimBlacklist, error) {
	event := new(RewardmanagerRemovedFromAutoClaimBlacklist)
	if err := _Rewardmanager.contract.UnpackLog(event, "RemovedFromAutoClaimBlacklist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerRewardsClaimedIterator is returned from FilterRewardsClaimed and is used to iterate over the raw logs and unpacked data for RewardsClaimed events raised by the Rewardmanager contract.
type RewardmanagerRewardsClaimedIterator struct {
	Event *RewardmanagerRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *RewardmanagerRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerRewardsClaimed)
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
		it.Event = new(RewardmanagerRewardsClaimed)
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
func (it *RewardmanagerRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerRewardsClaimed represents a RewardsClaimed event raised by the Rewardmanager contract.
type RewardmanagerRewardsClaimed struct {
	MsgSender common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardsClaimed is a free log retrieval operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed msgSender, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, msgSender []common.Address) (*RewardmanagerRewardsClaimedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "RewardsClaimed", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerRewardsClaimedIterator{contract: _Rewardmanager.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsClaimed is a free log subscription operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed msgSender, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewardmanagerRewardsClaimed, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "RewardsClaimed", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerRewardsClaimed)
				if err := _Rewardmanager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

// ParseRewardsClaimed is a log parse operation binding the contract event 0xfc30cddea38e2bf4d6ea7d3f9ed3b6ad7f176419f4963bd81318067a4aee73fe.
//
// Solidity: event RewardsClaimed(address indexed msgSender, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParseRewardsClaimed(log types.Log) (*RewardmanagerRewardsClaimed, error) {
	event := new(RewardmanagerRewardsClaimed)
	if err := _Rewardmanager.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerRewardsMigratedIterator is returned from FilterRewardsMigrated and is used to iterate over the raw logs and unpacked data for RewardsMigrated events raised by the Rewardmanager contract.
type RewardmanagerRewardsMigratedIterator struct {
	Event *RewardmanagerRewardsMigrated // Event containing the contract specifics and raw log

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
func (it *RewardmanagerRewardsMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerRewardsMigrated)
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
		it.Event = new(RewardmanagerRewardsMigrated)
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
func (it *RewardmanagerRewardsMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerRewardsMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerRewardsMigrated represents a RewardsMigrated event raised by the Rewardmanager contract.
type RewardmanagerRewardsMigrated struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardsMigrated is a free log retrieval operation binding the contract event 0x6eac28563b2ee99d5e30af1598e85ff38d4c26055247564c50f29db3264026c5.
//
// Solidity: event RewardsMigrated(address indexed from, address indexed to, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) FilterRewardsMigrated(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RewardmanagerRewardsMigratedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "RewardsMigrated", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerRewardsMigratedIterator{contract: _Rewardmanager.contract, event: "RewardsMigrated", logs: logs, sub: sub}, nil
}

// WatchRewardsMigrated is a free log subscription operation binding the contract event 0x6eac28563b2ee99d5e30af1598e85ff38d4c26055247564c50f29db3264026c5.
//
// Solidity: event RewardsMigrated(address indexed from, address indexed to, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) WatchRewardsMigrated(opts *bind.WatchOpts, sink chan<- *RewardmanagerRewardsMigrated, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "RewardsMigrated", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerRewardsMigrated)
				if err := _Rewardmanager.contract.UnpackLog(event, "RewardsMigrated", log); err != nil {
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

// ParseRewardsMigrated is a log parse operation binding the contract event 0x6eac28563b2ee99d5e30af1598e85ff38d4c26055247564c50f29db3264026c5.
//
// Solidity: event RewardsMigrated(address indexed from, address indexed to, uint256 amount)
func (_Rewardmanager *RewardmanagerFilterer) ParseRewardsMigrated(log types.Log) (*RewardmanagerRewardsMigrated, error) {
	event := new(RewardmanagerRewardsMigrated)
	if err := _Rewardmanager.contract.UnpackLog(event, "RewardsMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Rewardmanager contract.
type RewardmanagerUnpausedIterator struct {
	Event *RewardmanagerUnpaused // Event containing the contract specifics and raw log

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
func (it *RewardmanagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerUnpaused)
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
		it.Event = new(RewardmanagerUnpaused)
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
func (it *RewardmanagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerUnpaused represents a Unpaused event raised by the Rewardmanager contract.
type RewardmanagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rewardmanager *RewardmanagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RewardmanagerUnpausedIterator, error) {

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RewardmanagerUnpausedIterator{contract: _Rewardmanager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rewardmanager *RewardmanagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RewardmanagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerUnpaused)
				if err := _Rewardmanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParseUnpaused(log types.Log) (*RewardmanagerUnpaused, error) {
	event := new(RewardmanagerUnpaused)
	if err := _Rewardmanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Rewardmanager contract.
type RewardmanagerUpgradedIterator struct {
	Event *RewardmanagerUpgraded // Event containing the contract specifics and raw log

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
func (it *RewardmanagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerUpgraded)
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
		it.Event = new(RewardmanagerUpgraded)
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
func (it *RewardmanagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerUpgraded represents a Upgraded event raised by the Rewardmanager contract.
type RewardmanagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rewardmanager *RewardmanagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RewardmanagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerUpgradedIterator{contract: _Rewardmanager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rewardmanager *RewardmanagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RewardmanagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerUpgraded)
				if err := _Rewardmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Rewardmanager *RewardmanagerFilterer) ParseUpgraded(log types.Log) (*RewardmanagerUpgraded, error) {
	event := new(RewardmanagerUpgraded)
	if err := _Rewardmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardmanagerVanillaRegistrySetIterator is returned from FilterVanillaRegistrySet and is used to iterate over the raw logs and unpacked data for VanillaRegistrySet events raised by the Rewardmanager contract.
type RewardmanagerVanillaRegistrySetIterator struct {
	Event *RewardmanagerVanillaRegistrySet // Event containing the contract specifics and raw log

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
func (it *RewardmanagerVanillaRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardmanagerVanillaRegistrySet)
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
		it.Event = new(RewardmanagerVanillaRegistrySet)
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
func (it *RewardmanagerVanillaRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardmanagerVanillaRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardmanagerVanillaRegistrySet represents a VanillaRegistrySet event raised by the Rewardmanager contract.
type RewardmanagerVanillaRegistrySet struct {
	NewVanillaRegistry common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterVanillaRegistrySet is a free log retrieval operation binding the contract event 0x98be997fcfb397b31fe449613068b2538ab70e3efc11f00b337bc816388a8a3d.
//
// Solidity: event VanillaRegistrySet(address indexed newVanillaRegistry)
func (_Rewardmanager *RewardmanagerFilterer) FilterVanillaRegistrySet(opts *bind.FilterOpts, newVanillaRegistry []common.Address) (*RewardmanagerVanillaRegistrySetIterator, error) {

	var newVanillaRegistryRule []interface{}
	for _, newVanillaRegistryItem := range newVanillaRegistry {
		newVanillaRegistryRule = append(newVanillaRegistryRule, newVanillaRegistryItem)
	}

	logs, sub, err := _Rewardmanager.contract.FilterLogs(opts, "VanillaRegistrySet", newVanillaRegistryRule)
	if err != nil {
		return nil, err
	}
	return &RewardmanagerVanillaRegistrySetIterator{contract: _Rewardmanager.contract, event: "VanillaRegistrySet", logs: logs, sub: sub}, nil
}

// WatchVanillaRegistrySet is a free log subscription operation binding the contract event 0x98be997fcfb397b31fe449613068b2538ab70e3efc11f00b337bc816388a8a3d.
//
// Solidity: event VanillaRegistrySet(address indexed newVanillaRegistry)
func (_Rewardmanager *RewardmanagerFilterer) WatchVanillaRegistrySet(opts *bind.WatchOpts, sink chan<- *RewardmanagerVanillaRegistrySet, newVanillaRegistry []common.Address) (event.Subscription, error) {

	var newVanillaRegistryRule []interface{}
	for _, newVanillaRegistryItem := range newVanillaRegistry {
		newVanillaRegistryRule = append(newVanillaRegistryRule, newVanillaRegistryItem)
	}

	logs, sub, err := _Rewardmanager.contract.WatchLogs(opts, "VanillaRegistrySet", newVanillaRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardmanagerVanillaRegistrySet)
				if err := _Rewardmanager.contract.UnpackLog(event, "VanillaRegistrySet", log); err != nil {
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

// ParseVanillaRegistrySet is a log parse operation binding the contract event 0x98be997fcfb397b31fe449613068b2538ab70e3efc11f00b337bc816388a8a3d.
//
// Solidity: event VanillaRegistrySet(address indexed newVanillaRegistry)
func (_Rewardmanager *RewardmanagerFilterer) ParseVanillaRegistrySet(log types.Log) (*RewardmanagerVanillaRegistrySet, error) {
	event := new(RewardmanagerVanillaRegistrySet)
	if err := _Rewardmanager.contract.UnpackLog(event, "VanillaRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
