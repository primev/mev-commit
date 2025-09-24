// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rewarddistributor

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

// IRewardDistributorDistribution is an auto generated low-level Go binding around an user-defined struct.
type IRewardDistributorDistribution struct {
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
}

// RewarddistributorMetaData contains all meta data concerning the Rewarddistributor contract.
var RewarddistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimDelegate\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegate\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimOnbehalfOfOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipients\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"recipients\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getKeyRecipient\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPendingRewards\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantETHRewards\",\"inputs\":[{\"name\":\"rewardList\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardDistributor.Distribution[]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"grantTokenRewards\",\"inputs\":[{\"name\":\"rewardList\",\"type\":\"tuple[]\",\"internalType\":\"structIRewardDistributor.Distribution[]\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_rewardManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrateExistingRewards\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"operatorGlobalOverride\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorKeyOverrides\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"keyhash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"overrideRecipientByPubkey\",\"inputs\":[{\"name\":\"pubkeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reclaimStipendsToOwner\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"recipients\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rewardData\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"accrued\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"claimed\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rewardManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rewardTokens\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setClaimDelegate\",\"inputs\":[{\"name\":\"delegate\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorGlobalOverride\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardManager\",\"inputs\":[{\"name\":\"_rewardManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardToken\",\"inputs\":[{\"name\":\"_rewardToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"ClaimDelegateSet\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"delegate\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ETHGranted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ETHRewardsClaimed\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorGlobalOverrideSet\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecipientSet\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardManagerSet\",\"inputs\":[{\"name\":\"rewardManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardTokenSet\",\"inputs\":[{\"name\":\"rewardToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsBatchGranted\",\"inputs\":[{\"name\":\"tokenID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsMigrated\",\"inputs\":[{\"name\":\"tokenID\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsReclaimed\",\"inputs\":[{\"name\":\"tokenID\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenRewardsClaimed\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokensGranted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressInsufficientBalance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IncorrectPaymentAmount\",\"inputs\":[{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidClaimDelegate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRewardToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTokenID\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoClaimableRewards\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwnerOrRewardManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RewardsTransferFailed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// RewarddistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use RewarddistributorMetaData.ABI instead.
var RewarddistributorABI = RewarddistributorMetaData.ABI

// Rewarddistributor is an auto generated Go binding around an Ethereum contract.
type Rewarddistributor struct {
	RewarddistributorCaller     // Read-only binding to the contract
	RewarddistributorTransactor // Write-only binding to the contract
	RewarddistributorFilterer   // Log filterer for contract events
}

// RewarddistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewarddistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewarddistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewarddistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewarddistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewarddistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewarddistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewarddistributorSession struct {
	Contract     *Rewarddistributor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RewarddistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewarddistributorCallerSession struct {
	Contract *RewarddistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// RewarddistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewarddistributorTransactorSession struct {
	Contract     *RewarddistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// RewarddistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewarddistributorRaw struct {
	Contract *Rewarddistributor // Generic contract binding to access the raw methods on
}

// RewarddistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewarddistributorCallerRaw struct {
	Contract *RewarddistributorCaller // Generic read-only contract binding to access the raw methods on
}

// RewarddistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewarddistributorTransactorRaw struct {
	Contract *RewarddistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewarddistributor creates a new instance of Rewarddistributor, bound to a specific deployed contract.
func NewRewarddistributor(address common.Address, backend bind.ContractBackend) (*Rewarddistributor, error) {
	contract, err := bindRewarddistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rewarddistributor{RewarddistributorCaller: RewarddistributorCaller{contract: contract}, RewarddistributorTransactor: RewarddistributorTransactor{contract: contract}, RewarddistributorFilterer: RewarddistributorFilterer{contract: contract}}, nil
}

// NewRewarddistributorCaller creates a new read-only instance of Rewarddistributor, bound to a specific deployed contract.
func NewRewarddistributorCaller(address common.Address, caller bind.ContractCaller) (*RewarddistributorCaller, error) {
	contract, err := bindRewarddistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorCaller{contract: contract}, nil
}

// NewRewarddistributorTransactor creates a new write-only instance of Rewarddistributor, bound to a specific deployed contract.
func NewRewarddistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*RewarddistributorTransactor, error) {
	contract, err := bindRewarddistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorTransactor{contract: contract}, nil
}

// NewRewarddistributorFilterer creates a new log filterer instance of Rewarddistributor, bound to a specific deployed contract.
func NewRewarddistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*RewarddistributorFilterer, error) {
	contract, err := bindRewarddistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorFilterer{contract: contract}, nil
}

// bindRewarddistributor binds a generic wrapper to an already deployed contract.
func bindRewarddistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewarddistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewarddistributor *RewarddistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rewarddistributor.Contract.RewarddistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewarddistributor *RewarddistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.RewarddistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewarddistributor *RewarddistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.RewarddistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rewarddistributor *RewarddistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rewarddistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rewarddistributor *RewarddistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rewarddistributor *RewarddistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewarddistributor *RewarddistributorCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewarddistributor *RewarddistributorSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rewarddistributor.Contract.UPGRADEINTERFACEVERSION(&_Rewarddistributor.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Rewarddistributor *RewarddistributorCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Rewarddistributor.Contract.UPGRADEINTERFACEVERSION(&_Rewarddistributor.CallOpts)
}

// ClaimDelegate is a free data retrieval call binding the contract method 0x281fe909.
//
// Solidity: function claimDelegate(address operator, address recipient, address delegate) view returns(bool)
func (_Rewarddistributor *RewarddistributorCaller) ClaimDelegate(opts *bind.CallOpts, operator common.Address, recipient common.Address, delegate common.Address) (bool, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "claimDelegate", operator, recipient, delegate)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ClaimDelegate is a free data retrieval call binding the contract method 0x281fe909.
//
// Solidity: function claimDelegate(address operator, address recipient, address delegate) view returns(bool)
func (_Rewarddistributor *RewarddistributorSession) ClaimDelegate(operator common.Address, recipient common.Address, delegate common.Address) (bool, error) {
	return _Rewarddistributor.Contract.ClaimDelegate(&_Rewarddistributor.CallOpts, operator, recipient, delegate)
}

// ClaimDelegate is a free data retrieval call binding the contract method 0x281fe909.
//
// Solidity: function claimDelegate(address operator, address recipient, address delegate) view returns(bool)
func (_Rewarddistributor *RewarddistributorCallerSession) ClaimDelegate(operator common.Address, recipient common.Address, delegate common.Address) (bool, error) {
	return _Rewarddistributor.Contract.ClaimDelegate(&_Rewarddistributor.CallOpts, operator, recipient, delegate)
}

// GetKeyRecipient is a free data retrieval call binding the contract method 0xa4529e70.
//
// Solidity: function getKeyRecipient(address operator, bytes pubkey) view returns(address)
func (_Rewarddistributor *RewarddistributorCaller) GetKeyRecipient(opts *bind.CallOpts, operator common.Address, pubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "getKeyRecipient", operator, pubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetKeyRecipient is a free data retrieval call binding the contract method 0xa4529e70.
//
// Solidity: function getKeyRecipient(address operator, bytes pubkey) view returns(address)
func (_Rewarddistributor *RewarddistributorSession) GetKeyRecipient(operator common.Address, pubkey []byte) (common.Address, error) {
	return _Rewarddistributor.Contract.GetKeyRecipient(&_Rewarddistributor.CallOpts, operator, pubkey)
}

// GetKeyRecipient is a free data retrieval call binding the contract method 0xa4529e70.
//
// Solidity: function getKeyRecipient(address operator, bytes pubkey) view returns(address)
func (_Rewarddistributor *RewarddistributorCallerSession) GetKeyRecipient(operator common.Address, pubkey []byte) (common.Address, error) {
	return _Rewarddistributor.Contract.GetKeyRecipient(&_Rewarddistributor.CallOpts, operator, pubkey)
}

// GetPendingRewards is a free data retrieval call binding the contract method 0x32d30141.
//
// Solidity: function getPendingRewards(address operator, address recipient, uint256 tokenID) view returns(uint128)
func (_Rewarddistributor *RewarddistributorCaller) GetPendingRewards(opts *bind.CallOpts, operator common.Address, recipient common.Address, tokenID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "getPendingRewards", operator, recipient, tokenID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingRewards is a free data retrieval call binding the contract method 0x32d30141.
//
// Solidity: function getPendingRewards(address operator, address recipient, uint256 tokenID) view returns(uint128)
func (_Rewarddistributor *RewarddistributorSession) GetPendingRewards(operator common.Address, recipient common.Address, tokenID *big.Int) (*big.Int, error) {
	return _Rewarddistributor.Contract.GetPendingRewards(&_Rewarddistributor.CallOpts, operator, recipient, tokenID)
}

// GetPendingRewards is a free data retrieval call binding the contract method 0x32d30141.
//
// Solidity: function getPendingRewards(address operator, address recipient, uint256 tokenID) view returns(uint128)
func (_Rewarddistributor *RewarddistributorCallerSession) GetPendingRewards(operator common.Address, recipient common.Address, tokenID *big.Int) (*big.Int, error) {
	return _Rewarddistributor.Contract.GetPendingRewards(&_Rewarddistributor.CallOpts, operator, recipient, tokenID)
}

// OperatorGlobalOverride is a free data retrieval call binding the contract method 0x60aa8706.
//
// Solidity: function operatorGlobalOverride(address operator) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorCaller) OperatorGlobalOverride(opts *bind.CallOpts, operator common.Address) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "operatorGlobalOverride", operator)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorGlobalOverride is a free data retrieval call binding the contract method 0x60aa8706.
//
// Solidity: function operatorGlobalOverride(address operator) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorSession) OperatorGlobalOverride(operator common.Address) (common.Address, error) {
	return _Rewarddistributor.Contract.OperatorGlobalOverride(&_Rewarddistributor.CallOpts, operator)
}

// OperatorGlobalOverride is a free data retrieval call binding the contract method 0x60aa8706.
//
// Solidity: function operatorGlobalOverride(address operator) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorCallerSession) OperatorGlobalOverride(operator common.Address) (common.Address, error) {
	return _Rewarddistributor.Contract.OperatorGlobalOverride(&_Rewarddistributor.CallOpts, operator)
}

// OperatorKeyOverrides is a free data retrieval call binding the contract method 0x27e53420.
//
// Solidity: function operatorKeyOverrides(address operator, bytes32 keyhash) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorCaller) OperatorKeyOverrides(opts *bind.CallOpts, operator common.Address, keyhash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "operatorKeyOverrides", operator, keyhash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorKeyOverrides is a free data retrieval call binding the contract method 0x27e53420.
//
// Solidity: function operatorKeyOverrides(address operator, bytes32 keyhash) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorSession) OperatorKeyOverrides(operator common.Address, keyhash [32]byte) (common.Address, error) {
	return _Rewarddistributor.Contract.OperatorKeyOverrides(&_Rewarddistributor.CallOpts, operator, keyhash)
}

// OperatorKeyOverrides is a free data retrieval call binding the contract method 0x27e53420.
//
// Solidity: function operatorKeyOverrides(address operator, bytes32 keyhash) view returns(address recipient)
func (_Rewarddistributor *RewarddistributorCallerSession) OperatorKeyOverrides(operator common.Address, keyhash [32]byte) (common.Address, error) {
	return _Rewarddistributor.Contract.OperatorKeyOverrides(&_Rewarddistributor.CallOpts, operator, keyhash)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewarddistributor *RewarddistributorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewarddistributor *RewarddistributorSession) Owner() (common.Address, error) {
	return _Rewarddistributor.Contract.Owner(&_Rewarddistributor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Rewarddistributor *RewarddistributorCallerSession) Owner() (common.Address, error) {
	return _Rewarddistributor.Contract.Owner(&_Rewarddistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewarddistributor *RewarddistributorCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewarddistributor *RewarddistributorSession) Paused() (bool, error) {
	return _Rewarddistributor.Contract.Paused(&_Rewarddistributor.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Rewarddistributor *RewarddistributorCallerSession) Paused() (bool, error) {
	return _Rewarddistributor.Contract.Paused(&_Rewarddistributor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewarddistributor *RewarddistributorCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewarddistributor *RewarddistributorSession) PendingOwner() (common.Address, error) {
	return _Rewarddistributor.Contract.PendingOwner(&_Rewarddistributor.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Rewarddistributor *RewarddistributorCallerSession) PendingOwner() (common.Address, error) {
	return _Rewarddistributor.Contract.PendingOwner(&_Rewarddistributor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewarddistributor *RewarddistributorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewarddistributor *RewarddistributorSession) ProxiableUUID() ([32]byte, error) {
	return _Rewarddistributor.Contract.ProxiableUUID(&_Rewarddistributor.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Rewarddistributor *RewarddistributorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Rewarddistributor.Contract.ProxiableUUID(&_Rewarddistributor.CallOpts)
}

// RewardData is a free data retrieval call binding the contract method 0x68a99ef6.
//
// Solidity: function rewardData(address operator, address recipient, uint256 tokenID) view returns(uint128 accrued, uint128 claimed)
func (_Rewarddistributor *RewarddistributorCaller) RewardData(opts *bind.CallOpts, operator common.Address, recipient common.Address, tokenID *big.Int) (struct {
	Accrued *big.Int
	Claimed *big.Int
}, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "rewardData", operator, recipient, tokenID)

	outstruct := new(struct {
		Accrued *big.Int
		Claimed *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Accrued = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Claimed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RewardData is a free data retrieval call binding the contract method 0x68a99ef6.
//
// Solidity: function rewardData(address operator, address recipient, uint256 tokenID) view returns(uint128 accrued, uint128 claimed)
func (_Rewarddistributor *RewarddistributorSession) RewardData(operator common.Address, recipient common.Address, tokenID *big.Int) (struct {
	Accrued *big.Int
	Claimed *big.Int
}, error) {
	return _Rewarddistributor.Contract.RewardData(&_Rewarddistributor.CallOpts, operator, recipient, tokenID)
}

// RewardData is a free data retrieval call binding the contract method 0x68a99ef6.
//
// Solidity: function rewardData(address operator, address recipient, uint256 tokenID) view returns(uint128 accrued, uint128 claimed)
func (_Rewarddistributor *RewarddistributorCallerSession) RewardData(operator common.Address, recipient common.Address, tokenID *big.Int) (struct {
	Accrued *big.Int
	Claimed *big.Int
}, error) {
	return _Rewarddistributor.Contract.RewardData(&_Rewarddistributor.CallOpts, operator, recipient, tokenID)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Rewarddistributor *RewarddistributorCaller) RewardManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "rewardManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Rewarddistributor *RewarddistributorSession) RewardManager() (common.Address, error) {
	return _Rewarddistributor.Contract.RewardManager(&_Rewarddistributor.CallOpts)
}

// RewardManager is a free data retrieval call binding the contract method 0x0f4ef8a6.
//
// Solidity: function rewardManager() view returns(address)
func (_Rewarddistributor *RewarddistributorCallerSession) RewardManager() (common.Address, error) {
	return _Rewarddistributor.Contract.RewardManager(&_Rewarddistributor.CallOpts)
}

// RewardTokens is a free data retrieval call binding the contract method 0x7bb7bed1.
//
// Solidity: function rewardTokens(uint256 id) view returns(address token)
func (_Rewarddistributor *RewarddistributorCaller) RewardTokens(opts *bind.CallOpts, id *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Rewarddistributor.contract.Call(opts, &out, "rewardTokens", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardTokens is a free data retrieval call binding the contract method 0x7bb7bed1.
//
// Solidity: function rewardTokens(uint256 id) view returns(address token)
func (_Rewarddistributor *RewarddistributorSession) RewardTokens(id *big.Int) (common.Address, error) {
	return _Rewarddistributor.Contract.RewardTokens(&_Rewarddistributor.CallOpts, id)
}

// RewardTokens is a free data retrieval call binding the contract method 0x7bb7bed1.
//
// Solidity: function rewardTokens(uint256 id) view returns(address token)
func (_Rewarddistributor *RewarddistributorCallerSession) RewardTokens(id *big.Int) (common.Address, error) {
	return _Rewarddistributor.Contract.RewardTokens(&_Rewarddistributor.CallOpts, id)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewarddistributor *RewarddistributorTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewarddistributor *RewarddistributorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.AcceptOwnership(&_Rewarddistributor.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.AcceptOwnership(&_Rewarddistributor.TransactOpts)
}

// ClaimOnbehalfOfOperator is a paid mutator transaction binding the contract method 0x61714a3b.
//
// Solidity: function claimOnbehalfOfOperator(address operator, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactor) ClaimOnbehalfOfOperator(opts *bind.TransactOpts, operator common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "claimOnbehalfOfOperator", operator, recipients, tokenID)
}

// ClaimOnbehalfOfOperator is a paid mutator transaction binding the contract method 0x61714a3b.
//
// Solidity: function claimOnbehalfOfOperator(address operator, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorSession) ClaimOnbehalfOfOperator(operator common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ClaimOnbehalfOfOperator(&_Rewarddistributor.TransactOpts, operator, recipients, tokenID)
}

// ClaimOnbehalfOfOperator is a paid mutator transaction binding the contract method 0x61714a3b.
//
// Solidity: function claimOnbehalfOfOperator(address operator, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) ClaimOnbehalfOfOperator(operator common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ClaimOnbehalfOfOperator(&_Rewarddistributor.TransactOpts, operator, recipients, tokenID)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xfad66604.
//
// Solidity: function claimRewards(address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactor) ClaimRewards(opts *bind.TransactOpts, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "claimRewards", recipients, tokenID)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xfad66604.
//
// Solidity: function claimRewards(address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorSession) ClaimRewards(recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ClaimRewards(&_Rewarddistributor.TransactOpts, recipients, tokenID)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xfad66604.
//
// Solidity: function claimRewards(address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) ClaimRewards(recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ClaimRewards(&_Rewarddistributor.TransactOpts, recipients, tokenID)
}

// GrantETHRewards is a paid mutator transaction binding the contract method 0xa990eff1.
//
// Solidity: function grantETHRewards((address,address,uint128)[] rewardList) payable returns()
func (_Rewarddistributor *RewarddistributorTransactor) GrantETHRewards(opts *bind.TransactOpts, rewardList []IRewardDistributorDistribution) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "grantETHRewards", rewardList)
}

// GrantETHRewards is a paid mutator transaction binding the contract method 0xa990eff1.
//
// Solidity: function grantETHRewards((address,address,uint128)[] rewardList) payable returns()
func (_Rewarddistributor *RewarddistributorSession) GrantETHRewards(rewardList []IRewardDistributorDistribution) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.GrantETHRewards(&_Rewarddistributor.TransactOpts, rewardList)
}

// GrantETHRewards is a paid mutator transaction binding the contract method 0xa990eff1.
//
// Solidity: function grantETHRewards((address,address,uint128)[] rewardList) payable returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) GrantETHRewards(rewardList []IRewardDistributorDistribution) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.GrantETHRewards(&_Rewarddistributor.TransactOpts, rewardList)
}

// GrantTokenRewards is a paid mutator transaction binding the contract method 0x9bf6c9f5.
//
// Solidity: function grantTokenRewards((address,address,uint128)[] rewardList, uint256 tokenID) payable returns()
func (_Rewarddistributor *RewarddistributorTransactor) GrantTokenRewards(opts *bind.TransactOpts, rewardList []IRewardDistributorDistribution, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "grantTokenRewards", rewardList, tokenID)
}

// GrantTokenRewards is a paid mutator transaction binding the contract method 0x9bf6c9f5.
//
// Solidity: function grantTokenRewards((address,address,uint128)[] rewardList, uint256 tokenID) payable returns()
func (_Rewarddistributor *RewarddistributorSession) GrantTokenRewards(rewardList []IRewardDistributorDistribution, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.GrantTokenRewards(&_Rewarddistributor.TransactOpts, rewardList, tokenID)
}

// GrantTokenRewards is a paid mutator transaction binding the contract method 0x9bf6c9f5.
//
// Solidity: function grantTokenRewards((address,address,uint128)[] rewardList, uint256 tokenID) payable returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) GrantTokenRewards(rewardList []IRewardDistributorDistribution, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.GrantTokenRewards(&_Rewarddistributor.TransactOpts, rewardList, tokenID)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "initialize", _owner, _rewardManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorSession) Initialize(_owner common.Address, _rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Initialize(&_Rewarddistributor.TransactOpts, _owner, _rewardManager)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _owner, address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) Initialize(_owner common.Address, _rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Initialize(&_Rewarddistributor.TransactOpts, _owner, _rewardManager)
}

// MigrateExistingRewards is a paid mutator transaction binding the contract method 0x4c7bf229.
//
// Solidity: function migrateExistingRewards(address from, address to, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactor) MigrateExistingRewards(opts *bind.TransactOpts, from common.Address, to common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "migrateExistingRewards", from, to, tokenID)
}

// MigrateExistingRewards is a paid mutator transaction binding the contract method 0x4c7bf229.
//
// Solidity: function migrateExistingRewards(address from, address to, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorSession) MigrateExistingRewards(from common.Address, to common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.MigrateExistingRewards(&_Rewarddistributor.TransactOpts, from, to, tokenID)
}

// MigrateExistingRewards is a paid mutator transaction binding the contract method 0x4c7bf229.
//
// Solidity: function migrateExistingRewards(address from, address to, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) MigrateExistingRewards(from common.Address, to common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.MigrateExistingRewards(&_Rewarddistributor.TransactOpts, from, to, tokenID)
}

// OverrideRecipientByPubkey is a paid mutator transaction binding the contract method 0xd5a2693e.
//
// Solidity: function overrideRecipientByPubkey(bytes[] pubkeys, address recipient) returns()
func (_Rewarddistributor *RewarddistributorTransactor) OverrideRecipientByPubkey(opts *bind.TransactOpts, pubkeys [][]byte, recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "overrideRecipientByPubkey", pubkeys, recipient)
}

// OverrideRecipientByPubkey is a paid mutator transaction binding the contract method 0xd5a2693e.
//
// Solidity: function overrideRecipientByPubkey(bytes[] pubkeys, address recipient) returns()
func (_Rewarddistributor *RewarddistributorSession) OverrideRecipientByPubkey(pubkeys [][]byte, recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.OverrideRecipientByPubkey(&_Rewarddistributor.TransactOpts, pubkeys, recipient)
}

// OverrideRecipientByPubkey is a paid mutator transaction binding the contract method 0xd5a2693e.
//
// Solidity: function overrideRecipientByPubkey(bytes[] pubkeys, address recipient) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) OverrideRecipientByPubkey(pubkeys [][]byte, recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.OverrideRecipientByPubkey(&_Rewarddistributor.TransactOpts, pubkeys, recipient)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewarddistributor *RewarddistributorTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewarddistributor *RewarddistributorSession) Pause() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Pause(&_Rewarddistributor.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) Pause() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Pause(&_Rewarddistributor.TransactOpts)
}

// ReclaimStipendsToOwner is a paid mutator transaction binding the contract method 0xcfb9a4f9.
//
// Solidity: function reclaimStipendsToOwner(address[] operators, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactor) ReclaimStipendsToOwner(opts *bind.TransactOpts, operators []common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "reclaimStipendsToOwner", operators, recipients, tokenID)
}

// ReclaimStipendsToOwner is a paid mutator transaction binding the contract method 0xcfb9a4f9.
//
// Solidity: function reclaimStipendsToOwner(address[] operators, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorSession) ReclaimStipendsToOwner(operators []common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ReclaimStipendsToOwner(&_Rewarddistributor.TransactOpts, operators, recipients, tokenID)
}

// ReclaimStipendsToOwner is a paid mutator transaction binding the contract method 0xcfb9a4f9.
//
// Solidity: function reclaimStipendsToOwner(address[] operators, address[] recipients, uint256 tokenID) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) ReclaimStipendsToOwner(operators []common.Address, recipients []common.Address, tokenID *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.ReclaimStipendsToOwner(&_Rewarddistributor.TransactOpts, operators, recipients, tokenID)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewarddistributor *RewarddistributorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewarddistributor *RewarddistributorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.RenounceOwnership(&_Rewarddistributor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.RenounceOwnership(&_Rewarddistributor.TransactOpts)
}

// SetClaimDelegate is a paid mutator transaction binding the contract method 0x537d22fe.
//
// Solidity: function setClaimDelegate(address delegate, address recipient, bool status) returns()
func (_Rewarddistributor *RewarddistributorTransactor) SetClaimDelegate(opts *bind.TransactOpts, delegate common.Address, recipient common.Address, status bool) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "setClaimDelegate", delegate, recipient, status)
}

// SetClaimDelegate is a paid mutator transaction binding the contract method 0x537d22fe.
//
// Solidity: function setClaimDelegate(address delegate, address recipient, bool status) returns()
func (_Rewarddistributor *RewarddistributorSession) SetClaimDelegate(delegate common.Address, recipient common.Address, status bool) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetClaimDelegate(&_Rewarddistributor.TransactOpts, delegate, recipient, status)
}

// SetClaimDelegate is a paid mutator transaction binding the contract method 0x537d22fe.
//
// Solidity: function setClaimDelegate(address delegate, address recipient, bool status) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) SetClaimDelegate(delegate common.Address, recipient common.Address, status bool) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetClaimDelegate(&_Rewarddistributor.TransactOpts, delegate, recipient, status)
}

// SetOperatorGlobalOverride is a paid mutator transaction binding the contract method 0x86ee57b7.
//
// Solidity: function setOperatorGlobalOverride(address recipient) returns()
func (_Rewarddistributor *RewarddistributorTransactor) SetOperatorGlobalOverride(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "setOperatorGlobalOverride", recipient)
}

// SetOperatorGlobalOverride is a paid mutator transaction binding the contract method 0x86ee57b7.
//
// Solidity: function setOperatorGlobalOverride(address recipient) returns()
func (_Rewarddistributor *RewarddistributorSession) SetOperatorGlobalOverride(recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetOperatorGlobalOverride(&_Rewarddistributor.TransactOpts, recipient)
}

// SetOperatorGlobalOverride is a paid mutator transaction binding the contract method 0x86ee57b7.
//
// Solidity: function setOperatorGlobalOverride(address recipient) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) SetOperatorGlobalOverride(recipient common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetOperatorGlobalOverride(&_Rewarddistributor.TransactOpts, recipient)
}

// SetRewardManager is a paid mutator transaction binding the contract method 0x153ee554.
//
// Solidity: function setRewardManager(address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorTransactor) SetRewardManager(opts *bind.TransactOpts, _rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "setRewardManager", _rewardManager)
}

// SetRewardManager is a paid mutator transaction binding the contract method 0x153ee554.
//
// Solidity: function setRewardManager(address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorSession) SetRewardManager(_rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetRewardManager(&_Rewarddistributor.TransactOpts, _rewardManager)
}

// SetRewardManager is a paid mutator transaction binding the contract method 0x153ee554.
//
// Solidity: function setRewardManager(address _rewardManager) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) SetRewardManager(_rewardManager common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetRewardManager(&_Rewarddistributor.TransactOpts, _rewardManager)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0xf77a3fb9.
//
// Solidity: function setRewardToken(address _rewardToken, uint256 _id) returns()
func (_Rewarddistributor *RewarddistributorTransactor) SetRewardToken(opts *bind.TransactOpts, _rewardToken common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "setRewardToken", _rewardToken, _id)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0xf77a3fb9.
//
// Solidity: function setRewardToken(address _rewardToken, uint256 _id) returns()
func (_Rewarddistributor *RewarddistributorSession) SetRewardToken(_rewardToken common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetRewardToken(&_Rewarddistributor.TransactOpts, _rewardToken, _id)
}

// SetRewardToken is a paid mutator transaction binding the contract method 0xf77a3fb9.
//
// Solidity: function setRewardToken(address _rewardToken, uint256 _id) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) SetRewardToken(_rewardToken common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.SetRewardToken(&_Rewarddistributor.TransactOpts, _rewardToken, _id)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewarddistributor *RewarddistributorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewarddistributor *RewarddistributorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.TransferOwnership(&_Rewarddistributor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.TransferOwnership(&_Rewarddistributor.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewarddistributor *RewarddistributorTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewarddistributor *RewarddistributorSession) Unpause() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Unpause(&_Rewarddistributor.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) Unpause() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Unpause(&_Rewarddistributor.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewarddistributor *RewarddistributorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewarddistributor.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewarddistributor *RewarddistributorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.UpgradeToAndCall(&_Rewarddistributor.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.UpgradeToAndCall(&_Rewarddistributor.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewarddistributor *RewarddistributorTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Rewarddistributor.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewarddistributor *RewarddistributorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Fallback(&_Rewarddistributor.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Fallback(&_Rewarddistributor.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewarddistributor *RewarddistributorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rewarddistributor.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewarddistributor *RewarddistributorSession) Receive() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Receive(&_Rewarddistributor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Rewarddistributor *RewarddistributorTransactorSession) Receive() (*types.Transaction, error) {
	return _Rewarddistributor.Contract.Receive(&_Rewarddistributor.TransactOpts)
}

// RewarddistributorClaimDelegateSetIterator is returned from FilterClaimDelegateSet and is used to iterate over the raw logs and unpacked data for ClaimDelegateSet events raised by the Rewarddistributor contract.
type RewarddistributorClaimDelegateSetIterator struct {
	Event *RewarddistributorClaimDelegateSet // Event containing the contract specifics and raw log

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
func (it *RewarddistributorClaimDelegateSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorClaimDelegateSet)
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
		it.Event = new(RewarddistributorClaimDelegateSet)
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
func (it *RewarddistributorClaimDelegateSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorClaimDelegateSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorClaimDelegateSet represents a ClaimDelegateSet event raised by the Rewarddistributor contract.
type RewarddistributorClaimDelegateSet struct {
	Operator  common.Address
	Recipient common.Address
	Delegate  common.Address
	Status    bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimDelegateSet is a free log retrieval operation binding the contract event 0x678890bc989aa73fd3b36ff01efad3e0a551d088550af66284a3b876cafc1fd2.
//
// Solidity: event ClaimDelegateSet(address indexed operator, address indexed recipient, address indexed delegate, bool status)
func (_Rewarddistributor *RewarddistributorFilterer) FilterClaimDelegateSet(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address, delegate []common.Address) (*RewarddistributorClaimDelegateSetIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "ClaimDelegateSet", operatorRule, recipientRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorClaimDelegateSetIterator{contract: _Rewarddistributor.contract, event: "ClaimDelegateSet", logs: logs, sub: sub}, nil
}

// WatchClaimDelegateSet is a free log subscription operation binding the contract event 0x678890bc989aa73fd3b36ff01efad3e0a551d088550af66284a3b876cafc1fd2.
//
// Solidity: event ClaimDelegateSet(address indexed operator, address indexed recipient, address indexed delegate, bool status)
func (_Rewarddistributor *RewarddistributorFilterer) WatchClaimDelegateSet(opts *bind.WatchOpts, sink chan<- *RewarddistributorClaimDelegateSet, operator []common.Address, recipient []common.Address, delegate []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "ClaimDelegateSet", operatorRule, recipientRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorClaimDelegateSet)
				if err := _Rewarddistributor.contract.UnpackLog(event, "ClaimDelegateSet", log); err != nil {
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

// ParseClaimDelegateSet is a log parse operation binding the contract event 0x678890bc989aa73fd3b36ff01efad3e0a551d088550af66284a3b876cafc1fd2.
//
// Solidity: event ClaimDelegateSet(address indexed operator, address indexed recipient, address indexed delegate, bool status)
func (_Rewarddistributor *RewarddistributorFilterer) ParseClaimDelegateSet(log types.Log) (*RewarddistributorClaimDelegateSet, error) {
	event := new(RewarddistributorClaimDelegateSet)
	if err := _Rewarddistributor.contract.UnpackLog(event, "ClaimDelegateSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorETHGrantedIterator is returned from FilterETHGranted and is used to iterate over the raw logs and unpacked data for ETHGranted events raised by the Rewarddistributor contract.
type RewarddistributorETHGrantedIterator struct {
	Event *RewarddistributorETHGranted // Event containing the contract specifics and raw log

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
func (it *RewarddistributorETHGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorETHGranted)
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
		it.Event = new(RewarddistributorETHGranted)
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
func (it *RewarddistributorETHGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorETHGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorETHGranted represents a ETHGranted event raised by the Rewarddistributor contract.
type RewarddistributorETHGranted struct {
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterETHGranted is a free log retrieval operation binding the contract event 0x4633f6e8264d8c05cba97be33e0f1eaae97dbf81597c4da193d52da4b7b4d1f1.
//
// Solidity: event ETHGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterETHGranted(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address, amount []*big.Int) (*RewarddistributorETHGrantedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "ETHGranted", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorETHGrantedIterator{contract: _Rewarddistributor.contract, event: "ETHGranted", logs: logs, sub: sub}, nil
}

// WatchETHGranted is a free log subscription operation binding the contract event 0x4633f6e8264d8c05cba97be33e0f1eaae97dbf81597c4da193d52da4b7b4d1f1.
//
// Solidity: event ETHGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchETHGranted(opts *bind.WatchOpts, sink chan<- *RewarddistributorETHGranted, operator []common.Address, recipient []common.Address, amount []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "ETHGranted", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorETHGranted)
				if err := _Rewarddistributor.contract.UnpackLog(event, "ETHGranted", log); err != nil {
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

// ParseETHGranted is a log parse operation binding the contract event 0x4633f6e8264d8c05cba97be33e0f1eaae97dbf81597c4da193d52da4b7b4d1f1.
//
// Solidity: event ETHGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseETHGranted(log types.Log) (*RewarddistributorETHGranted, error) {
	event := new(RewarddistributorETHGranted)
	if err := _Rewarddistributor.contract.UnpackLog(event, "ETHGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorETHRewardsClaimedIterator is returned from FilterETHRewardsClaimed and is used to iterate over the raw logs and unpacked data for ETHRewardsClaimed events raised by the Rewarddistributor contract.
type RewarddistributorETHRewardsClaimedIterator struct {
	Event *RewarddistributorETHRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *RewarddistributorETHRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorETHRewardsClaimed)
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
		it.Event = new(RewarddistributorETHRewardsClaimed)
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
func (it *RewarddistributorETHRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorETHRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorETHRewardsClaimed represents a ETHRewardsClaimed event raised by the Rewarddistributor contract.
type RewarddistributorETHRewardsClaimed struct {
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterETHRewardsClaimed is a free log retrieval operation binding the contract event 0xd74da12fa2c485462374064b211802914920859a34df550f95c1c5d5f811fe42.
//
// Solidity: event ETHRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterETHRewardsClaimed(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address, amount []*big.Int) (*RewarddistributorETHRewardsClaimedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "ETHRewardsClaimed", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorETHRewardsClaimedIterator{contract: _Rewarddistributor.contract, event: "ETHRewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchETHRewardsClaimed is a free log subscription operation binding the contract event 0xd74da12fa2c485462374064b211802914920859a34df550f95c1c5d5f811fe42.
//
// Solidity: event ETHRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchETHRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewarddistributorETHRewardsClaimed, operator []common.Address, recipient []common.Address, amount []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "ETHRewardsClaimed", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorETHRewardsClaimed)
				if err := _Rewarddistributor.contract.UnpackLog(event, "ETHRewardsClaimed", log); err != nil {
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

// ParseETHRewardsClaimed is a log parse operation binding the contract event 0xd74da12fa2c485462374064b211802914920859a34df550f95c1c5d5f811fe42.
//
// Solidity: event ETHRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseETHRewardsClaimed(log types.Log) (*RewarddistributorETHRewardsClaimed, error) {
	event := new(RewarddistributorETHRewardsClaimed)
	if err := _Rewarddistributor.contract.UnpackLog(event, "ETHRewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Rewarddistributor contract.
type RewarddistributorInitializedIterator struct {
	Event *RewarddistributorInitialized // Event containing the contract specifics and raw log

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
func (it *RewarddistributorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorInitialized)
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
		it.Event = new(RewarddistributorInitialized)
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
func (it *RewarddistributorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorInitialized represents a Initialized event raised by the Rewarddistributor contract.
type RewarddistributorInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rewarddistributor *RewarddistributorFilterer) FilterInitialized(opts *bind.FilterOpts) (*RewarddistributorInitializedIterator, error) {

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &RewarddistributorInitializedIterator{contract: _Rewarddistributor.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Rewarddistributor *RewarddistributorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *RewarddistributorInitialized) (event.Subscription, error) {

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorInitialized)
				if err := _Rewarddistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParseInitialized(log types.Log) (*RewarddistributorInitialized, error) {
	event := new(RewarddistributorInitialized)
	if err := _Rewarddistributor.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorOperatorGlobalOverrideSetIterator is returned from FilterOperatorGlobalOverrideSet and is used to iterate over the raw logs and unpacked data for OperatorGlobalOverrideSet events raised by the Rewarddistributor contract.
type RewarddistributorOperatorGlobalOverrideSetIterator struct {
	Event *RewarddistributorOperatorGlobalOverrideSet // Event containing the contract specifics and raw log

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
func (it *RewarddistributorOperatorGlobalOverrideSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorOperatorGlobalOverrideSet)
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
		it.Event = new(RewarddistributorOperatorGlobalOverrideSet)
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
func (it *RewarddistributorOperatorGlobalOverrideSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorOperatorGlobalOverrideSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorOperatorGlobalOverrideSet represents a OperatorGlobalOverrideSet event raised by the Rewarddistributor contract.
type RewarddistributorOperatorGlobalOverrideSet struct {
	Operator  common.Address
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOperatorGlobalOverrideSet is a free log retrieval operation binding the contract event 0x8fd2f6227c705a2e4c9af87a9c9e974f7c2b5c124ce9f4c418881e2de432f41f.
//
// Solidity: event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) FilterOperatorGlobalOverrideSet(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address) (*RewarddistributorOperatorGlobalOverrideSetIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "OperatorGlobalOverrideSet", operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorOperatorGlobalOverrideSetIterator{contract: _Rewarddistributor.contract, event: "OperatorGlobalOverrideSet", logs: logs, sub: sub}, nil
}

// WatchOperatorGlobalOverrideSet is a free log subscription operation binding the contract event 0x8fd2f6227c705a2e4c9af87a9c9e974f7c2b5c124ce9f4c418881e2de432f41f.
//
// Solidity: event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) WatchOperatorGlobalOverrideSet(opts *bind.WatchOpts, sink chan<- *RewarddistributorOperatorGlobalOverrideSet, operator []common.Address, recipient []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "OperatorGlobalOverrideSet", operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorOperatorGlobalOverrideSet)
				if err := _Rewarddistributor.contract.UnpackLog(event, "OperatorGlobalOverrideSet", log); err != nil {
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

// ParseOperatorGlobalOverrideSet is a log parse operation binding the contract event 0x8fd2f6227c705a2e4c9af87a9c9e974f7c2b5c124ce9f4c418881e2de432f41f.
//
// Solidity: event OperatorGlobalOverrideSet(address indexed operator, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) ParseOperatorGlobalOverrideSet(log types.Log) (*RewarddistributorOperatorGlobalOverrideSet, error) {
	event := new(RewarddistributorOperatorGlobalOverrideSet)
	if err := _Rewarddistributor.contract.UnpackLog(event, "OperatorGlobalOverrideSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Rewarddistributor contract.
type RewarddistributorOwnershipTransferStartedIterator struct {
	Event *RewarddistributorOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *RewarddistributorOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorOwnershipTransferStarted)
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
		it.Event = new(RewarddistributorOwnershipTransferStarted)
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
func (it *RewarddistributorOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Rewarddistributor contract.
type RewarddistributorOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rewarddistributor *RewarddistributorFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewarddistributorOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorOwnershipTransferStartedIterator{contract: _Rewarddistributor.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Rewarddistributor *RewarddistributorFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *RewarddistributorOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorOwnershipTransferStarted)
				if err := _Rewarddistributor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParseOwnershipTransferStarted(log types.Log) (*RewarddistributorOwnershipTransferStarted, error) {
	event := new(RewarddistributorOwnershipTransferStarted)
	if err := _Rewarddistributor.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rewarddistributor contract.
type RewarddistributorOwnershipTransferredIterator struct {
	Event *RewarddistributorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewarddistributorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorOwnershipTransferred)
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
		it.Event = new(RewarddistributorOwnershipTransferred)
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
func (it *RewarddistributorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorOwnershipTransferred represents a OwnershipTransferred event raised by the Rewarddistributor contract.
type RewarddistributorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewarddistributor *RewarddistributorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewarddistributorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorOwnershipTransferredIterator{contract: _Rewarddistributor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Rewarddistributor *RewarddistributorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewarddistributorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorOwnershipTransferred)
				if err := _Rewarddistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParseOwnershipTransferred(log types.Log) (*RewarddistributorOwnershipTransferred, error) {
	event := new(RewarddistributorOwnershipTransferred)
	if err := _Rewarddistributor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Rewarddistributor contract.
type RewarddistributorPausedIterator struct {
	Event *RewarddistributorPaused // Event containing the contract specifics and raw log

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
func (it *RewarddistributorPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorPaused)
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
		it.Event = new(RewarddistributorPaused)
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
func (it *RewarddistributorPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorPaused represents a Paused event raised by the Rewarddistributor contract.
type RewarddistributorPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rewarddistributor *RewarddistributorFilterer) FilterPaused(opts *bind.FilterOpts) (*RewarddistributorPausedIterator, error) {

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RewarddistributorPausedIterator{contract: _Rewarddistributor.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Rewarddistributor *RewarddistributorFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RewarddistributorPaused) (event.Subscription, error) {

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorPaused)
				if err := _Rewarddistributor.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParsePaused(log types.Log) (*RewarddistributorPaused, error) {
	event := new(RewarddistributorPaused)
	if err := _Rewarddistributor.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRecipientSetIterator is returned from FilterRecipientSet and is used to iterate over the raw logs and unpacked data for RecipientSet events raised by the Rewarddistributor contract.
type RewarddistributorRecipientSetIterator struct {
	Event *RewarddistributorRecipientSet // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRecipientSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRecipientSet)
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
		it.Event = new(RewarddistributorRecipientSet)
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
func (it *RewarddistributorRecipientSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRecipientSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRecipientSet represents a RecipientSet event raised by the Rewarddistributor contract.
type RewarddistributorRecipientSet struct {
	Operator  common.Address
	Pubkey    []byte
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRecipientSet is a free log retrieval operation binding the contract event 0xdd1c6108b84f1dfef6f5b45484370f1a270fad78f335e45b376a351f489b5460.
//
// Solidity: event RecipientSet(address indexed operator, bytes pubkey, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRecipientSet(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address) (*RewarddistributorRecipientSetIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RecipientSet", operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRecipientSetIterator{contract: _Rewarddistributor.contract, event: "RecipientSet", logs: logs, sub: sub}, nil
}

// WatchRecipientSet is a free log subscription operation binding the contract event 0xdd1c6108b84f1dfef6f5b45484370f1a270fad78f335e45b376a351f489b5460.
//
// Solidity: event RecipientSet(address indexed operator, bytes pubkey, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRecipientSet(opts *bind.WatchOpts, sink chan<- *RewarddistributorRecipientSet, operator []common.Address, recipient []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RecipientSet", operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRecipientSet)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RecipientSet", log); err != nil {
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

// ParseRecipientSet is a log parse operation binding the contract event 0xdd1c6108b84f1dfef6f5b45484370f1a270fad78f335e45b376a351f489b5460.
//
// Solidity: event RecipientSet(address indexed operator, bytes pubkey, address indexed recipient)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRecipientSet(log types.Log) (*RewarddistributorRecipientSet, error) {
	event := new(RewarddistributorRecipientSet)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RecipientSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRewardManagerSetIterator is returned from FilterRewardManagerSet and is used to iterate over the raw logs and unpacked data for RewardManagerSet events raised by the Rewarddistributor contract.
type RewarddistributorRewardManagerSetIterator struct {
	Event *RewarddistributorRewardManagerSet // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRewardManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRewardManagerSet)
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
		it.Event = new(RewarddistributorRewardManagerSet)
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
func (it *RewarddistributorRewardManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRewardManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRewardManagerSet represents a RewardManagerSet event raised by the Rewarddistributor contract.
type RewarddistributorRewardManagerSet struct {
	RewardManager common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRewardManagerSet is a free log retrieval operation binding the contract event 0x4b722357462a92ae2776fce05f39902ac96cb04d2fe661eb688f4a72ab069cff.
//
// Solidity: event RewardManagerSet(address indexed rewardManager)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRewardManagerSet(opts *bind.FilterOpts, rewardManager []common.Address) (*RewarddistributorRewardManagerSetIterator, error) {

	var rewardManagerRule []interface{}
	for _, rewardManagerItem := range rewardManager {
		rewardManagerRule = append(rewardManagerRule, rewardManagerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RewardManagerSet", rewardManagerRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRewardManagerSetIterator{contract: _Rewarddistributor.contract, event: "RewardManagerSet", logs: logs, sub: sub}, nil
}

// WatchRewardManagerSet is a free log subscription operation binding the contract event 0x4b722357462a92ae2776fce05f39902ac96cb04d2fe661eb688f4a72ab069cff.
//
// Solidity: event RewardManagerSet(address indexed rewardManager)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRewardManagerSet(opts *bind.WatchOpts, sink chan<- *RewarddistributorRewardManagerSet, rewardManager []common.Address) (event.Subscription, error) {

	var rewardManagerRule []interface{}
	for _, rewardManagerItem := range rewardManager {
		rewardManagerRule = append(rewardManagerRule, rewardManagerItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RewardManagerSet", rewardManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRewardManagerSet)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RewardManagerSet", log); err != nil {
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

// ParseRewardManagerSet is a log parse operation binding the contract event 0x4b722357462a92ae2776fce05f39902ac96cb04d2fe661eb688f4a72ab069cff.
//
// Solidity: event RewardManagerSet(address indexed rewardManager)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRewardManagerSet(log types.Log) (*RewarddistributorRewardManagerSet, error) {
	event := new(RewarddistributorRewardManagerSet)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RewardManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRewardTokenSetIterator is returned from FilterRewardTokenSet and is used to iterate over the raw logs and unpacked data for RewardTokenSet events raised by the Rewarddistributor contract.
type RewarddistributorRewardTokenSetIterator struct {
	Event *RewarddistributorRewardTokenSet // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRewardTokenSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRewardTokenSet)
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
		it.Event = new(RewarddistributorRewardTokenSet)
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
func (it *RewarddistributorRewardTokenSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRewardTokenSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRewardTokenSet represents a RewardTokenSet event raised by the Rewarddistributor contract.
type RewarddistributorRewardTokenSet struct {
	RewardToken common.Address
	TokenID     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardTokenSet is a free log retrieval operation binding the contract event 0x699da1ad32a331acc89644abfd5a2d9960170bd81cde59bda93fe7f53d3c1f31.
//
// Solidity: event RewardTokenSet(address indexed rewardToken, uint256 indexed tokenID)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRewardTokenSet(opts *bind.FilterOpts, rewardToken []common.Address, tokenID []*big.Int) (*RewarddistributorRewardTokenSetIterator, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RewardTokenSet", rewardTokenRule, tokenIDRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRewardTokenSetIterator{contract: _Rewarddistributor.contract, event: "RewardTokenSet", logs: logs, sub: sub}, nil
}

// WatchRewardTokenSet is a free log subscription operation binding the contract event 0x699da1ad32a331acc89644abfd5a2d9960170bd81cde59bda93fe7f53d3c1f31.
//
// Solidity: event RewardTokenSet(address indexed rewardToken, uint256 indexed tokenID)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRewardTokenSet(opts *bind.WatchOpts, sink chan<- *RewarddistributorRewardTokenSet, rewardToken []common.Address, tokenID []*big.Int) (event.Subscription, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RewardTokenSet", rewardTokenRule, tokenIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRewardTokenSet)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RewardTokenSet", log); err != nil {
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

// ParseRewardTokenSet is a log parse operation binding the contract event 0x699da1ad32a331acc89644abfd5a2d9960170bd81cde59bda93fe7f53d3c1f31.
//
// Solidity: event RewardTokenSet(address indexed rewardToken, uint256 indexed tokenID)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRewardTokenSet(log types.Log) (*RewarddistributorRewardTokenSet, error) {
	event := new(RewarddistributorRewardTokenSet)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RewardTokenSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRewardsBatchGrantedIterator is returned from FilterRewardsBatchGranted and is used to iterate over the raw logs and unpacked data for RewardsBatchGranted events raised by the Rewarddistributor contract.
type RewarddistributorRewardsBatchGrantedIterator struct {
	Event *RewarddistributorRewardsBatchGranted // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRewardsBatchGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRewardsBatchGranted)
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
		it.Event = new(RewarddistributorRewardsBatchGranted)
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
func (it *RewarddistributorRewardsBatchGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRewardsBatchGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRewardsBatchGranted represents a RewardsBatchGranted event raised by the Rewarddistributor contract.
type RewarddistributorRewardsBatchGranted struct {
	TokenID *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRewardsBatchGranted is a free log retrieval operation binding the contract event 0x5a560578c2a2576a22870a295e55da0d824316461ca74857c685b7ab72391658.
//
// Solidity: event RewardsBatchGranted(uint256 indexed tokenID, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRewardsBatchGranted(opts *bind.FilterOpts, tokenID []*big.Int, amount []*big.Int) (*RewarddistributorRewardsBatchGrantedIterator, error) {

	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RewardsBatchGranted", tokenIDRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRewardsBatchGrantedIterator{contract: _Rewarddistributor.contract, event: "RewardsBatchGranted", logs: logs, sub: sub}, nil
}

// WatchRewardsBatchGranted is a free log subscription operation binding the contract event 0x5a560578c2a2576a22870a295e55da0d824316461ca74857c685b7ab72391658.
//
// Solidity: event RewardsBatchGranted(uint256 indexed tokenID, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRewardsBatchGranted(opts *bind.WatchOpts, sink chan<- *RewarddistributorRewardsBatchGranted, tokenID []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RewardsBatchGranted", tokenIDRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRewardsBatchGranted)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsBatchGranted", log); err != nil {
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

// ParseRewardsBatchGranted is a log parse operation binding the contract event 0x5a560578c2a2576a22870a295e55da0d824316461ca74857c685b7ab72391658.
//
// Solidity: event RewardsBatchGranted(uint256 indexed tokenID, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRewardsBatchGranted(log types.Log) (*RewarddistributorRewardsBatchGranted, error) {
	event := new(RewarddistributorRewardsBatchGranted)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsBatchGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRewardsMigratedIterator is returned from FilterRewardsMigrated and is used to iterate over the raw logs and unpacked data for RewardsMigrated events raised by the Rewarddistributor contract.
type RewarddistributorRewardsMigratedIterator struct {
	Event *RewarddistributorRewardsMigrated // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRewardsMigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRewardsMigrated)
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
		it.Event = new(RewarddistributorRewardsMigrated)
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
func (it *RewarddistributorRewardsMigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRewardsMigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRewardsMigrated represents a RewardsMigrated event raised by the Rewarddistributor contract.
type RewarddistributorRewardsMigrated struct {
	TokenID  *big.Int
	Operator common.Address
	From     common.Address
	To       common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRewardsMigrated is a free log retrieval operation binding the contract event 0x74103144248cb03c4a8398b58ebd933e52aad3e743dc6d25c7e566d6129d2e05.
//
// Solidity: event RewardsMigrated(uint256 tokenID, address indexed operator, address indexed from, address indexed to, uint128 amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRewardsMigrated(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*RewarddistributorRewardsMigratedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RewardsMigrated", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRewardsMigratedIterator{contract: _Rewarddistributor.contract, event: "RewardsMigrated", logs: logs, sub: sub}, nil
}

// WatchRewardsMigrated is a free log subscription operation binding the contract event 0x74103144248cb03c4a8398b58ebd933e52aad3e743dc6d25c7e566d6129d2e05.
//
// Solidity: event RewardsMigrated(uint256 tokenID, address indexed operator, address indexed from, address indexed to, uint128 amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRewardsMigrated(opts *bind.WatchOpts, sink chan<- *RewarddistributorRewardsMigrated, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RewardsMigrated", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRewardsMigrated)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsMigrated", log); err != nil {
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

// ParseRewardsMigrated is a log parse operation binding the contract event 0x74103144248cb03c4a8398b58ebd933e52aad3e743dc6d25c7e566d6129d2e05.
//
// Solidity: event RewardsMigrated(uint256 tokenID, address indexed operator, address indexed from, address indexed to, uint128 amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRewardsMigrated(log types.Log) (*RewarddistributorRewardsMigrated, error) {
	event := new(RewarddistributorRewardsMigrated)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsMigrated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorRewardsReclaimedIterator is returned from FilterRewardsReclaimed and is used to iterate over the raw logs and unpacked data for RewardsReclaimed events raised by the Rewarddistributor contract.
type RewarddistributorRewardsReclaimedIterator struct {
	Event *RewarddistributorRewardsReclaimed // Event containing the contract specifics and raw log

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
func (it *RewarddistributorRewardsReclaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorRewardsReclaimed)
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
		it.Event = new(RewarddistributorRewardsReclaimed)
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
func (it *RewarddistributorRewardsReclaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorRewardsReclaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorRewardsReclaimed represents a RewardsReclaimed event raised by the Rewarddistributor contract.
type RewarddistributorRewardsReclaimed struct {
	TokenID   *big.Int
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardsReclaimed is a free log retrieval operation binding the contract event 0xc541b4439a9af87e81276e3d4782683d5ec3ac36e3f806c4d8b6de9eba5acd85.
//
// Solidity: event RewardsReclaimed(uint256 indexed tokenID, address indexed operator, address indexed recipient, uint256 amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterRewardsReclaimed(opts *bind.FilterOpts, tokenID []*big.Int, operator []common.Address, recipient []common.Address) (*RewarddistributorRewardsReclaimedIterator, error) {

	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "RewardsReclaimed", tokenIDRule, operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorRewardsReclaimedIterator{contract: _Rewarddistributor.contract, event: "RewardsReclaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsReclaimed is a free log subscription operation binding the contract event 0xc541b4439a9af87e81276e3d4782683d5ec3ac36e3f806c4d8b6de9eba5acd85.
//
// Solidity: event RewardsReclaimed(uint256 indexed tokenID, address indexed operator, address indexed recipient, uint256 amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchRewardsReclaimed(opts *bind.WatchOpts, sink chan<- *RewarddistributorRewardsReclaimed, tokenID []*big.Int, operator []common.Address, recipient []common.Address) (event.Subscription, error) {

	var tokenIDRule []interface{}
	for _, tokenIDItem := range tokenID {
		tokenIDRule = append(tokenIDRule, tokenIDItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "RewardsReclaimed", tokenIDRule, operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorRewardsReclaimed)
				if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsReclaimed", log); err != nil {
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

// ParseRewardsReclaimed is a log parse operation binding the contract event 0xc541b4439a9af87e81276e3d4782683d5ec3ac36e3f806c4d8b6de9eba5acd85.
//
// Solidity: event RewardsReclaimed(uint256 indexed tokenID, address indexed operator, address indexed recipient, uint256 amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseRewardsReclaimed(log types.Log) (*RewarddistributorRewardsReclaimed, error) {
	event := new(RewarddistributorRewardsReclaimed)
	if err := _Rewarddistributor.contract.UnpackLog(event, "RewardsReclaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorTokenRewardsClaimedIterator is returned from FilterTokenRewardsClaimed and is used to iterate over the raw logs and unpacked data for TokenRewardsClaimed events raised by the Rewarddistributor contract.
type RewarddistributorTokenRewardsClaimedIterator struct {
	Event *RewarddistributorTokenRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *RewarddistributorTokenRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorTokenRewardsClaimed)
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
		it.Event = new(RewarddistributorTokenRewardsClaimed)
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
func (it *RewarddistributorTokenRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorTokenRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorTokenRewardsClaimed represents a TokenRewardsClaimed event raised by the Rewarddistributor contract.
type RewarddistributorTokenRewardsClaimed struct {
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokenRewardsClaimed is a free log retrieval operation binding the contract event 0x2ecbbad36dfb8af92a78e491dd8a3a195a314c9e3d7262dd303ae70bf6378e26.
//
// Solidity: event TokenRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterTokenRewardsClaimed(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address, amount []*big.Int) (*RewarddistributorTokenRewardsClaimedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "TokenRewardsClaimed", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorTokenRewardsClaimedIterator{contract: _Rewarddistributor.contract, event: "TokenRewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchTokenRewardsClaimed is a free log subscription operation binding the contract event 0x2ecbbad36dfb8af92a78e491dd8a3a195a314c9e3d7262dd303ae70bf6378e26.
//
// Solidity: event TokenRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchTokenRewardsClaimed(opts *bind.WatchOpts, sink chan<- *RewarddistributorTokenRewardsClaimed, operator []common.Address, recipient []common.Address, amount []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "TokenRewardsClaimed", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorTokenRewardsClaimed)
				if err := _Rewarddistributor.contract.UnpackLog(event, "TokenRewardsClaimed", log); err != nil {
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

// ParseTokenRewardsClaimed is a log parse operation binding the contract event 0x2ecbbad36dfb8af92a78e491dd8a3a195a314c9e3d7262dd303ae70bf6378e26.
//
// Solidity: event TokenRewardsClaimed(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseTokenRewardsClaimed(log types.Log) (*RewarddistributorTokenRewardsClaimed, error) {
	event := new(RewarddistributorTokenRewardsClaimed)
	if err := _Rewarddistributor.contract.UnpackLog(event, "TokenRewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorTokensGrantedIterator is returned from FilterTokensGranted and is used to iterate over the raw logs and unpacked data for TokensGranted events raised by the Rewarddistributor contract.
type RewarddistributorTokensGrantedIterator struct {
	Event *RewarddistributorTokensGranted // Event containing the contract specifics and raw log

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
func (it *RewarddistributorTokensGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorTokensGranted)
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
		it.Event = new(RewarddistributorTokensGranted)
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
func (it *RewarddistributorTokensGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorTokensGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorTokensGranted represents a TokensGranted event raised by the Rewarddistributor contract.
type RewarddistributorTokensGranted struct {
	Operator  common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokensGranted is a free log retrieval operation binding the contract event 0x4b1544c19293e59b936b0280f49423a62bdf27ffa57af6fdabe8e7a4a454d80e.
//
// Solidity: event TokensGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) FilterTokensGranted(opts *bind.FilterOpts, operator []common.Address, recipient []common.Address, amount []*big.Int) (*RewarddistributorTokensGrantedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "TokensGranted", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorTokensGrantedIterator{contract: _Rewarddistributor.contract, event: "TokensGranted", logs: logs, sub: sub}, nil
}

// WatchTokensGranted is a free log subscription operation binding the contract event 0x4b1544c19293e59b936b0280f49423a62bdf27ffa57af6fdabe8e7a4a454d80e.
//
// Solidity: event TokensGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) WatchTokensGranted(opts *bind.WatchOpts, sink chan<- *RewarddistributorTokensGranted, operator []common.Address, recipient []common.Address, amount []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "TokensGranted", operatorRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorTokensGranted)
				if err := _Rewarddistributor.contract.UnpackLog(event, "TokensGranted", log); err != nil {
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

// ParseTokensGranted is a log parse operation binding the contract event 0x4b1544c19293e59b936b0280f49423a62bdf27ffa57af6fdabe8e7a4a454d80e.
//
// Solidity: event TokensGranted(address indexed operator, address indexed recipient, uint256 indexed amount)
func (_Rewarddistributor *RewarddistributorFilterer) ParseTokensGranted(log types.Log) (*RewarddistributorTokensGranted, error) {
	event := new(RewarddistributorTokensGranted)
	if err := _Rewarddistributor.contract.UnpackLog(event, "TokensGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Rewarddistributor contract.
type RewarddistributorUnpausedIterator struct {
	Event *RewarddistributorUnpaused // Event containing the contract specifics and raw log

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
func (it *RewarddistributorUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorUnpaused)
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
		it.Event = new(RewarddistributorUnpaused)
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
func (it *RewarddistributorUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorUnpaused represents a Unpaused event raised by the Rewarddistributor contract.
type RewarddistributorUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rewarddistributor *RewarddistributorFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RewarddistributorUnpausedIterator, error) {

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RewarddistributorUnpausedIterator{contract: _Rewarddistributor.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Rewarddistributor *RewarddistributorFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RewarddistributorUnpaused) (event.Subscription, error) {

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorUnpaused)
				if err := _Rewarddistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParseUnpaused(log types.Log) (*RewarddistributorUnpaused, error) {
	event := new(RewarddistributorUnpaused)
	if err := _Rewarddistributor.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewarddistributorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Rewarddistributor contract.
type RewarddistributorUpgradedIterator struct {
	Event *RewarddistributorUpgraded // Event containing the contract specifics and raw log

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
func (it *RewarddistributorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewarddistributorUpgraded)
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
		it.Event = new(RewarddistributorUpgraded)
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
func (it *RewarddistributorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewarddistributorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewarddistributorUpgraded represents a Upgraded event raised by the Rewarddistributor contract.
type RewarddistributorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rewarddistributor *RewarddistributorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*RewarddistributorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rewarddistributor.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &RewarddistributorUpgradedIterator{contract: _Rewarddistributor.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Rewarddistributor *RewarddistributorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *RewarddistributorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Rewarddistributor.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewarddistributorUpgraded)
				if err := _Rewarddistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Rewarddistributor *RewarddistributorFilterer) ParseUpgraded(log types.Log) (*RewarddistributorUpgraded, error) {
	event := new(RewarddistributorUpgraded)
	if err := _Rewarddistributor.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
