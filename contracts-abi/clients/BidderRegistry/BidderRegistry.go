// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bidderregistry

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

// BidderregistryMetaData contains all meta data concerning the Bidderregistry contract.
var BidderregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"PERCENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidPayment\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIBidderRegistry.State\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bidderRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockTrackerContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBlockTracker\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositForWindow\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositForWindows\",\"inputs\":[{\"name\":\"windows\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"feePercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccumulatedProtocolFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeposit\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProviderAmount\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_protocolFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blockTracker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blocksPerWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_feePayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lockedFunds\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyWithdrawProtocolFee\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxBidPerBlock\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"openBid\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"preconfManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeeTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accumulatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastPayoutBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerAmount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"retrieveFunds\",\"inputs\":[{\"name\":\"windowToSettle\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBlockTrackerContract\",\"inputs\":[{\"name\":\"newBlockTrackerContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeePayoutPeriodBlocks\",\"inputs\":[{\"name\":\"newFeePayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeePercent\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewProtocolFeeRecipient\",\"inputs\":[{\"name\":\"newProtocolFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPreconfManager\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockFunds\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"usedFunds\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawBidderAmountFromWindow\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromWindows\",\"inputs\":[{\"name\":\"windows\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawProviderAmount\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BidderRegistered\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"depositedAmount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"windowNumber\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BidderWithdrawal\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BlockTrackerUpdated\",\"inputs\":[{\"name\":\"newBlockTracker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePayoutPeriodBlocksUpdated\",\"inputs\":[{\"name\":\"newFeePayoutPeriodBlocks\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePercentUpdated\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTransfer\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsRetrieved\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsRewarded\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PreconfManagerUpdated\",\"inputs\":[{\"name\":\"newPreconfManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProtocolFeeRecipientUpdated\",\"inputs\":[{\"name\":\"newProtocolFeeRecipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BidNotPreConfirmed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BidWasNotPreConfirmed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BidderAmountIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyBidderCanWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProviderAmountIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotPreconfManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferToBidderFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferToProviderFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"WindowNotSettled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawAfterWindowSettled\",\"inputs\":[]}]",
}

// BidderregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use BidderregistryMetaData.ABI instead.
var BidderregistryABI = BidderregistryMetaData.ABI

// Bidderregistry is an auto generated Go binding around an Ethereum contract.
type Bidderregistry struct {
	BidderregistryCaller     // Read-only binding to the contract
	BidderregistryTransactor // Write-only binding to the contract
	BidderregistryFilterer   // Log filterer for contract events
}

// BidderregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type BidderregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BidderregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BidderregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BidderregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BidderregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BidderregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BidderregistrySession struct {
	Contract     *Bidderregistry   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BidderregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BidderregistryCallerSession struct {
	Contract *BidderregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BidderregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BidderregistryTransactorSession struct {
	Contract     *BidderregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BidderregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type BidderregistryRaw struct {
	Contract *Bidderregistry // Generic contract binding to access the raw methods on
}

// BidderregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BidderregistryCallerRaw struct {
	Contract *BidderregistryCaller // Generic read-only contract binding to access the raw methods on
}

// BidderregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BidderregistryTransactorRaw struct {
	Contract *BidderregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBidderregistry creates a new instance of Bidderregistry, bound to a specific deployed contract.
func NewBidderregistry(address common.Address, backend bind.ContractBackend) (*Bidderregistry, error) {
	contract, err := bindBidderregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bidderregistry{BidderregistryCaller: BidderregistryCaller{contract: contract}, BidderregistryTransactor: BidderregistryTransactor{contract: contract}, BidderregistryFilterer: BidderregistryFilterer{contract: contract}}, nil
}

// NewBidderregistryCaller creates a new read-only instance of Bidderregistry, bound to a specific deployed contract.
func NewBidderregistryCaller(address common.Address, caller bind.ContractCaller) (*BidderregistryCaller, error) {
	contract, err := bindBidderregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BidderregistryCaller{contract: contract}, nil
}

// NewBidderregistryTransactor creates a new write-only instance of Bidderregistry, bound to a specific deployed contract.
func NewBidderregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*BidderregistryTransactor, error) {
	contract, err := bindBidderregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BidderregistryTransactor{contract: contract}, nil
}

// NewBidderregistryFilterer creates a new log filterer instance of Bidderregistry, bound to a specific deployed contract.
func NewBidderregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*BidderregistryFilterer, error) {
	contract, err := bindBidderregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFilterer{contract: contract}, nil
}

// bindBidderregistry binds a generic wrapper to an already deployed contract.
func bindBidderregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BidderregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bidderregistry *BidderregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bidderregistry.Contract.BidderregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bidderregistry *BidderregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.Contract.BidderregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bidderregistry *BidderregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bidderregistry.Contract.BidderregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bidderregistry *BidderregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bidderregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bidderregistry *BidderregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bidderregistry *BidderregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bidderregistry.Contract.contract.Transact(opts, method, params...)
}

// PERCENT is a free data retrieval call binding the contract method 0xb85a8b20.
//
// Solidity: function PERCENT() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) PERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PERCENT is a free data retrieval call binding the contract method 0xb85a8b20.
//
// Solidity: function PERCENT() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) PERCENT() (*big.Int, error) {
	return _Bidderregistry.Contract.PERCENT(&_Bidderregistry.CallOpts)
}

// PERCENT is a free data retrieval call binding the contract method 0xb85a8b20.
//
// Solidity: function PERCENT() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) PERCENT() (*big.Int, error) {
	return _Bidderregistry.Contract.PERCENT(&_Bidderregistry.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) PRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) PRECISION() (*big.Int, error) {
	return _Bidderregistry.Contract.PRECISION(&_Bidderregistry.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) PRECISION() (*big.Int, error) {
	return _Bidderregistry.Contract.PRECISION(&_Bidderregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Bidderregistry *BidderregistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Bidderregistry *BidderregistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Bidderregistry.Contract.UPGRADEINTERFACEVERSION(&_Bidderregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Bidderregistry *BidderregistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Bidderregistry.Contract.UPGRADEINTERFACEVERSION(&_Bidderregistry.CallOpts)
}

// BidPayment is a free data retrieval call binding the contract method 0x5a837876.
//
// Solidity: function bidPayment(bytes32 ) view returns(address bidder, uint256 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistryCaller) BidPayment(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt *big.Int
	State  uint8
}, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "bidPayment", arg0)

	outstruct := new(struct {
		Bidder common.Address
		BidAmt *big.Int
		State  uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Bidder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.BidAmt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// BidPayment is a free data retrieval call binding the contract method 0x5a837876.
//
// Solidity: function bidPayment(bytes32 ) view returns(address bidder, uint256 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistrySession) BidPayment(arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt *big.Int
	State  uint8
}, error) {
	return _Bidderregistry.Contract.BidPayment(&_Bidderregistry.CallOpts, arg0)
}

// BidPayment is a free data retrieval call binding the contract method 0x5a837876.
//
// Solidity: function bidPayment(bytes32 ) view returns(address bidder, uint256 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistryCallerSession) BidPayment(arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt *big.Int
	State  uint8
}, error) {
	return _Bidderregistry.Contract.BidPayment(&_Bidderregistry.CallOpts, arg0)
}

// BidderRegistered is a free data retrieval call binding the contract method 0x2a0773de.
//
// Solidity: function bidderRegistered(address ) view returns(bool)
func (_Bidderregistry *BidderregistryCaller) BidderRegistered(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "bidderRegistered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// BidderRegistered is a free data retrieval call binding the contract method 0x2a0773de.
//
// Solidity: function bidderRegistered(address ) view returns(bool)
func (_Bidderregistry *BidderregistrySession) BidderRegistered(arg0 common.Address) (bool, error) {
	return _Bidderregistry.Contract.BidderRegistered(&_Bidderregistry.CallOpts, arg0)
}

// BidderRegistered is a free data retrieval call binding the contract method 0x2a0773de.
//
// Solidity: function bidderRegistered(address ) view returns(bool)
func (_Bidderregistry *BidderregistryCallerSession) BidderRegistered(arg0 common.Address) (bool, error) {
	return _Bidderregistry.Contract.BidderRegistered(&_Bidderregistry.CallOpts, arg0)
}

// BlockTrackerContract is a free data retrieval call binding the contract method 0x6d82071b.
//
// Solidity: function blockTrackerContract() view returns(address)
func (_Bidderregistry *BidderregistryCaller) BlockTrackerContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "blockTrackerContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockTrackerContract is a free data retrieval call binding the contract method 0x6d82071b.
//
// Solidity: function blockTrackerContract() view returns(address)
func (_Bidderregistry *BidderregistrySession) BlockTrackerContract() (common.Address, error) {
	return _Bidderregistry.Contract.BlockTrackerContract(&_Bidderregistry.CallOpts)
}

// BlockTrackerContract is a free data retrieval call binding the contract method 0x6d82071b.
//
// Solidity: function blockTrackerContract() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) BlockTrackerContract() (common.Address, error) {
	return _Bidderregistry.Contract.BlockTrackerContract(&_Bidderregistry.CallOpts)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) BlocksPerWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "blocksPerWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) BlocksPerWindow() (*big.Int, error) {
	return _Bidderregistry.Contract.BlocksPerWindow(&_Bidderregistry.CallOpts)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) BlocksPerWindow() (*big.Int, error) {
	return _Bidderregistry.Contract.BlocksPerWindow(&_Bidderregistry.CallOpts)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Bidderregistry *BidderregistryCaller) FeePercent(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "feePercent")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Bidderregistry *BidderregistrySession) FeePercent() (uint16, error) {
	return _Bidderregistry.Contract.FeePercent(&_Bidderregistry.CallOpts)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint16)
func (_Bidderregistry *BidderregistryCallerSession) FeePercent() (uint16, error) {
	return _Bidderregistry.Contract.FeePercent(&_Bidderregistry.CallOpts)
}

// GetAccumulatedProtocolFee is a free data retrieval call binding the contract method 0x2dde2218.
//
// Solidity: function getAccumulatedProtocolFee() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) GetAccumulatedProtocolFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "getAccumulatedProtocolFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccumulatedProtocolFee is a free data retrieval call binding the contract method 0x2dde2218.
//
// Solidity: function getAccumulatedProtocolFee() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) GetAccumulatedProtocolFee() (*big.Int, error) {
	return _Bidderregistry.Contract.GetAccumulatedProtocolFee(&_Bidderregistry.CallOpts)
}

// GetAccumulatedProtocolFee is a free data retrieval call binding the contract method 0x2dde2218.
//
// Solidity: function getAccumulatedProtocolFee() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) GetAccumulatedProtocolFee() (*big.Int, error) {
	return _Bidderregistry.Contract.GetAccumulatedProtocolFee(&_Bidderregistry.CallOpts)
}

// GetDeposit is a free data retrieval call binding the contract method 0x2726b506.
//
// Solidity: function getDeposit(address bidder, uint256 window) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) GetDeposit(opts *bind.CallOpts, bidder common.Address, window *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "getDeposit", bidder, window)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDeposit is a free data retrieval call binding the contract method 0x2726b506.
//
// Solidity: function getDeposit(address bidder, uint256 window) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) GetDeposit(bidder common.Address, window *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.GetDeposit(&_Bidderregistry.CallOpts, bidder, window)
}

// GetDeposit is a free data retrieval call binding the contract method 0x2726b506.
//
// Solidity: function getDeposit(address bidder, uint256 window) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) GetDeposit(bidder common.Address, window *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.GetDeposit(&_Bidderregistry.CallOpts, bidder, window)
}

// GetProviderAmount is a free data retrieval call binding the contract method 0x0ebe2555.
//
// Solidity: function getProviderAmount(address provider) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) GetProviderAmount(opts *bind.CallOpts, provider common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "getProviderAmount", provider)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProviderAmount is a free data retrieval call binding the contract method 0x0ebe2555.
//
// Solidity: function getProviderAmount(address provider) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) GetProviderAmount(provider common.Address) (*big.Int, error) {
	return _Bidderregistry.Contract.GetProviderAmount(&_Bidderregistry.CallOpts, provider)
}

// GetProviderAmount is a free data retrieval call binding the contract method 0x0ebe2555.
//
// Solidity: function getProviderAmount(address provider) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) GetProviderAmount(provider common.Address) (*big.Int, error) {
	return _Bidderregistry.Contract.GetProviderAmount(&_Bidderregistry.CallOpts, provider)
}

// LockedFunds is a free data retrieval call binding the contract method 0x1355d861.
//
// Solidity: function lockedFunds(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) LockedFunds(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "lockedFunds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedFunds is a free data retrieval call binding the contract method 0x1355d861.
//
// Solidity: function lockedFunds(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) LockedFunds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.LockedFunds(&_Bidderregistry.CallOpts, arg0, arg1)
}

// LockedFunds is a free data retrieval call binding the contract method 0x1355d861.
//
// Solidity: function lockedFunds(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) LockedFunds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.LockedFunds(&_Bidderregistry.CallOpts, arg0, arg1)
}

// MaxBidPerBlock is a free data retrieval call binding the contract method 0xdcb95894.
//
// Solidity: function maxBidPerBlock(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) MaxBidPerBlock(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "maxBidPerBlock", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxBidPerBlock is a free data retrieval call binding the contract method 0xdcb95894.
//
// Solidity: function maxBidPerBlock(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) MaxBidPerBlock(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.MaxBidPerBlock(&_Bidderregistry.CallOpts, arg0, arg1)
}

// MaxBidPerBlock is a free data retrieval call binding the contract method 0xdcb95894.
//
// Solidity: function maxBidPerBlock(address , uint256 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) MaxBidPerBlock(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Bidderregistry.Contract.MaxBidPerBlock(&_Bidderregistry.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bidderregistry *BidderregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bidderregistry *BidderregistrySession) Owner() (common.Address, error) {
	return _Bidderregistry.Contract.Owner(&_Bidderregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) Owner() (common.Address, error) {
	return _Bidderregistry.Contract.Owner(&_Bidderregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bidderregistry *BidderregistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bidderregistry *BidderregistrySession) Paused() (bool, error) {
	return _Bidderregistry.Contract.Paused(&_Bidderregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bidderregistry *BidderregistryCallerSession) Paused() (bool, error) {
	return _Bidderregistry.Contract.Paused(&_Bidderregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Bidderregistry *BidderregistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Bidderregistry *BidderregistrySession) PendingOwner() (common.Address, error) {
	return _Bidderregistry.Contract.PendingOwner(&_Bidderregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) PendingOwner() (common.Address, error) {
	return _Bidderregistry.Contract.PendingOwner(&_Bidderregistry.CallOpts)
}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Bidderregistry *BidderregistryCaller) PreconfManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "preconfManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Bidderregistry *BidderregistrySession) PreconfManager() (common.Address, error) {
	return _Bidderregistry.Contract.PreconfManager(&_Bidderregistry.CallOpts)
}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) PreconfManager() (common.Address, error) {
	return _Bidderregistry.Contract.PreconfManager(&_Bidderregistry.CallOpts)
}

// ProtocolFeeTracker is a free data retrieval call binding the contract method 0x291af92c.
//
// Solidity: function protocolFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Bidderregistry *BidderregistryCaller) ProtocolFeeTracker(opts *bind.CallOpts) (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "protocolFeeTracker")

	outstruct := new(struct {
		Recipient          common.Address
		AccumulatedAmount  *big.Int
		LastPayoutBlock    *big.Int
		PayoutPeriodBlocks *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Recipient = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.AccumulatedAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastPayoutBlock = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PayoutPeriodBlocks = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ProtocolFeeTracker is a free data retrieval call binding the contract method 0x291af92c.
//
// Solidity: function protocolFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Bidderregistry *BidderregistrySession) ProtocolFeeTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Bidderregistry.Contract.ProtocolFeeTracker(&_Bidderregistry.CallOpts)
}

// ProtocolFeeTracker is a free data retrieval call binding the contract method 0x291af92c.
//
// Solidity: function protocolFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Bidderregistry *BidderregistryCallerSession) ProtocolFeeTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Bidderregistry.Contract.ProtocolFeeTracker(&_Bidderregistry.CallOpts)
}

// ProviderAmount is a free data retrieval call binding the contract method 0x180d02cb.
//
// Solidity: function providerAmount(address ) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) ProviderAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "providerAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProviderAmount is a free data retrieval call binding the contract method 0x180d02cb.
//
// Solidity: function providerAmount(address ) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) ProviderAmount(arg0 common.Address) (*big.Int, error) {
	return _Bidderregistry.Contract.ProviderAmount(&_Bidderregistry.CallOpts, arg0)
}

// ProviderAmount is a free data retrieval call binding the contract method 0x180d02cb.
//
// Solidity: function providerAmount(address ) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) ProviderAmount(arg0 common.Address) (*big.Int, error) {
	return _Bidderregistry.Contract.ProviderAmount(&_Bidderregistry.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Bidderregistry *BidderregistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Bidderregistry *BidderregistrySession) ProxiableUUID() ([32]byte, error) {
	return _Bidderregistry.Contract.ProxiableUUID(&_Bidderregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Bidderregistry *BidderregistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Bidderregistry.Contract.ProxiableUUID(&_Bidderregistry.CallOpts)
}

// UsedFunds is a free data retrieval call binding the contract method 0x79c14d5f.
//
// Solidity: function usedFunds(address , uint64 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) UsedFunds(opts *bind.CallOpts, arg0 common.Address, arg1 uint64) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "usedFunds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UsedFunds is a free data retrieval call binding the contract method 0x79c14d5f.
//
// Solidity: function usedFunds(address , uint64 ) view returns(uint256)
func (_Bidderregistry *BidderregistrySession) UsedFunds(arg0 common.Address, arg1 uint64) (*big.Int, error) {
	return _Bidderregistry.Contract.UsedFunds(&_Bidderregistry.CallOpts, arg0, arg1)
}

// UsedFunds is a free data retrieval call binding the contract method 0x79c14d5f.
//
// Solidity: function usedFunds(address , uint64 ) view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) UsedFunds(arg0 common.Address, arg1 uint64) (*big.Int, error) {
	return _Bidderregistry.Contract.UsedFunds(&_Bidderregistry.CallOpts, arg0, arg1)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Bidderregistry *BidderregistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Bidderregistry *BidderregistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _Bidderregistry.Contract.AcceptOwnership(&_Bidderregistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Bidderregistry *BidderregistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Bidderregistry.Contract.AcceptOwnership(&_Bidderregistry.TransactOpts)
}

// DepositForWindow is a paid mutator transaction binding the contract method 0xde7fa558.
//
// Solidity: function depositForWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistryTransactor) DepositForWindow(opts *bind.TransactOpts, window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "depositForWindow", window)
}

// DepositForWindow is a paid mutator transaction binding the contract method 0xde7fa558.
//
// Solidity: function depositForWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistrySession) DepositForWindow(window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForWindow(&_Bidderregistry.TransactOpts, window)
}

// DepositForWindow is a paid mutator transaction binding the contract method 0xde7fa558.
//
// Solidity: function depositForWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) DepositForWindow(window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForWindow(&_Bidderregistry.TransactOpts, window)
}

// DepositForWindows is a paid mutator transaction binding the contract method 0x531a2052.
//
// Solidity: function depositForWindows(uint256[] windows) payable returns()
func (_Bidderregistry *BidderregistryTransactor) DepositForWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "depositForWindows", windows)
}

// DepositForWindows is a paid mutator transaction binding the contract method 0x531a2052.
//
// Solidity: function depositForWindows(uint256[] windows) payable returns()
func (_Bidderregistry *BidderregistrySession) DepositForWindows(windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForWindows(&_Bidderregistry.TransactOpts, windows)
}

// DepositForWindows is a paid mutator transaction binding the contract method 0x531a2052.
//
// Solidity: function depositForWindows(uint256[] windows) payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) DepositForWindows(windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForWindows(&_Bidderregistry.TransactOpts, windows)
}

// Initialize is a paid mutator transaction binding the contract method 0x83efe8aa.
//
// Solidity: function initialize(address _protocolFeeRecipient, uint16 _feePercent, address _owner, address _blockTracker, uint256 _blocksPerWindow, uint256 _feePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistryTransactor) Initialize(opts *bind.TransactOpts, _protocolFeeRecipient common.Address, _feePercent uint16, _owner common.Address, _blockTracker common.Address, _blocksPerWindow *big.Int, _feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "initialize", _protocolFeeRecipient, _feePercent, _owner, _blockTracker, _blocksPerWindow, _feePayoutPeriodBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0x83efe8aa.
//
// Solidity: function initialize(address _protocolFeeRecipient, uint16 _feePercent, address _owner, address _blockTracker, uint256 _blocksPerWindow, uint256 _feePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistrySession) Initialize(_protocolFeeRecipient common.Address, _feePercent uint16, _owner common.Address, _blockTracker common.Address, _blocksPerWindow *big.Int, _feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.Initialize(&_Bidderregistry.TransactOpts, _protocolFeeRecipient, _feePercent, _owner, _blockTracker, _blocksPerWindow, _feePayoutPeriodBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0x83efe8aa.
//
// Solidity: function initialize(address _protocolFeeRecipient, uint16 _feePercent, address _owner, address _blockTracker, uint256 _blocksPerWindow, uint256 _feePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistryTransactorSession) Initialize(_protocolFeeRecipient common.Address, _feePercent uint16, _owner common.Address, _blockTracker common.Address, _blocksPerWindow *big.Int, _feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.Initialize(&_Bidderregistry.TransactOpts, _protocolFeeRecipient, _feePercent, _owner, _blockTracker, _blocksPerWindow, _feePayoutPeriodBlocks)
}

// ManuallyWithdrawProtocolFee is a paid mutator transaction binding the contract method 0xdbf63530.
//
// Solidity: function manuallyWithdrawProtocolFee() returns()
func (_Bidderregistry *BidderregistryTransactor) ManuallyWithdrawProtocolFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "manuallyWithdrawProtocolFee")
}

// ManuallyWithdrawProtocolFee is a paid mutator transaction binding the contract method 0xdbf63530.
//
// Solidity: function manuallyWithdrawProtocolFee() returns()
func (_Bidderregistry *BidderregistrySession) ManuallyWithdrawProtocolFee() (*types.Transaction, error) {
	return _Bidderregistry.Contract.ManuallyWithdrawProtocolFee(&_Bidderregistry.TransactOpts)
}

// ManuallyWithdrawProtocolFee is a paid mutator transaction binding the contract method 0xdbf63530.
//
// Solidity: function manuallyWithdrawProtocolFee() returns()
func (_Bidderregistry *BidderregistryTransactorSession) ManuallyWithdrawProtocolFee() (*types.Transaction, error) {
	return _Bidderregistry.Contract.ManuallyWithdrawProtocolFee(&_Bidderregistry.TransactOpts)
}

// OpenBid is a paid mutator transaction binding the contract method 0x2a241d75.
//
// Solidity: function openBid(bytes32 commitmentDigest, uint256 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistryTransactor) OpenBid(opts *bind.TransactOpts, commitmentDigest [32]byte, bid *big.Int, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "openBid", commitmentDigest, bid, bidder, blockNumber)
}

// OpenBid is a paid mutator transaction binding the contract method 0x2a241d75.
//
// Solidity: function openBid(bytes32 commitmentDigest, uint256 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistrySession) OpenBid(commitmentDigest [32]byte, bid *big.Int, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.Contract.OpenBid(&_Bidderregistry.TransactOpts, commitmentDigest, bid, bidder, blockNumber)
}

// OpenBid is a paid mutator transaction binding the contract method 0x2a241d75.
//
// Solidity: function openBid(bytes32 commitmentDigest, uint256 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistryTransactorSession) OpenBid(commitmentDigest [32]byte, bid *big.Int, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.Contract.OpenBid(&_Bidderregistry.TransactOpts, commitmentDigest, bid, bidder, blockNumber)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bidderregistry *BidderregistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bidderregistry *BidderregistrySession) Pause() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Pause(&_Bidderregistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bidderregistry *BidderregistryTransactorSession) Pause() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Pause(&_Bidderregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bidderregistry *BidderregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bidderregistry *BidderregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Bidderregistry.Contract.RenounceOwnership(&_Bidderregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bidderregistry *BidderregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bidderregistry.Contract.RenounceOwnership(&_Bidderregistry.TransactOpts)
}

// RetrieveFunds is a paid mutator transaction binding the contract method 0x0c2e5b0e.
//
// Solidity: function retrieveFunds(uint256 windowToSettle, bytes32 commitmentDigest, address provider, uint256 residualBidPercentAfterDecay) returns()
func (_Bidderregistry *BidderregistryTransactor) RetrieveFunds(opts *bind.TransactOpts, windowToSettle *big.Int, commitmentDigest [32]byte, provider common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "retrieveFunds", windowToSettle, commitmentDigest, provider, residualBidPercentAfterDecay)
}

// RetrieveFunds is a paid mutator transaction binding the contract method 0x0c2e5b0e.
//
// Solidity: function retrieveFunds(uint256 windowToSettle, bytes32 commitmentDigest, address provider, uint256 residualBidPercentAfterDecay) returns()
func (_Bidderregistry *BidderregistrySession) RetrieveFunds(windowToSettle *big.Int, commitmentDigest [32]byte, provider common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.RetrieveFunds(&_Bidderregistry.TransactOpts, windowToSettle, commitmentDigest, provider, residualBidPercentAfterDecay)
}

// RetrieveFunds is a paid mutator transaction binding the contract method 0x0c2e5b0e.
//
// Solidity: function retrieveFunds(uint256 windowToSettle, bytes32 commitmentDigest, address provider, uint256 residualBidPercentAfterDecay) returns()
func (_Bidderregistry *BidderregistryTransactorSession) RetrieveFunds(windowToSettle *big.Int, commitmentDigest [32]byte, provider common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.RetrieveFunds(&_Bidderregistry.TransactOpts, windowToSettle, commitmentDigest, provider, residualBidPercentAfterDecay)
}

// SetBlockTrackerContract is a paid mutator transaction binding the contract method 0x338f61f5.
//
// Solidity: function setBlockTrackerContract(address newBlockTrackerContract) returns()
func (_Bidderregistry *BidderregistryTransactor) SetBlockTrackerContract(opts *bind.TransactOpts, newBlockTrackerContract common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setBlockTrackerContract", newBlockTrackerContract)
}

// SetBlockTrackerContract is a paid mutator transaction binding the contract method 0x338f61f5.
//
// Solidity: function setBlockTrackerContract(address newBlockTrackerContract) returns()
func (_Bidderregistry *BidderregistrySession) SetBlockTrackerContract(newBlockTrackerContract common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetBlockTrackerContract(&_Bidderregistry.TransactOpts, newBlockTrackerContract)
}

// SetBlockTrackerContract is a paid mutator transaction binding the contract method 0x338f61f5.
//
// Solidity: function setBlockTrackerContract(address newBlockTrackerContract) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetBlockTrackerContract(newBlockTrackerContract common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetBlockTrackerContract(&_Bidderregistry.TransactOpts, newBlockTrackerContract)
}

// SetNewFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xe29136c0.
//
// Solidity: function setNewFeePayoutPeriodBlocks(uint256 newFeePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistryTransactor) SetNewFeePayoutPeriodBlocks(opts *bind.TransactOpts, newFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setNewFeePayoutPeriodBlocks", newFeePayoutPeriodBlocks)
}

// SetNewFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xe29136c0.
//
// Solidity: function setNewFeePayoutPeriodBlocks(uint256 newFeePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistrySession) SetNewFeePayoutPeriodBlocks(newFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeePayoutPeriodBlocks(&_Bidderregistry.TransactOpts, newFeePayoutPeriodBlocks)
}

// SetNewFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xe29136c0.
//
// Solidity: function setNewFeePayoutPeriodBlocks(uint256 newFeePayoutPeriodBlocks) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetNewFeePayoutPeriodBlocks(newFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeePayoutPeriodBlocks(&_Bidderregistry.TransactOpts, newFeePayoutPeriodBlocks)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Bidderregistry *BidderregistryTransactor) SetNewFeePercent(opts *bind.TransactOpts, newFeePercent uint16) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setNewFeePercent", newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Bidderregistry *BidderregistrySession) SetNewFeePercent(newFeePercent uint16) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeePercent(&_Bidderregistry.TransactOpts, newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0xfb22febf.
//
// Solidity: function setNewFeePercent(uint16 newFeePercent) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetNewFeePercent(newFeePercent uint16) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeePercent(&_Bidderregistry.TransactOpts, newFeePercent)
}

// SetNewProtocolFeeRecipient is a paid mutator transaction binding the contract method 0x184ac28e.
//
// Solidity: function setNewProtocolFeeRecipient(address newProtocolFeeRecipient) returns()
func (_Bidderregistry *BidderregistryTransactor) SetNewProtocolFeeRecipient(opts *bind.TransactOpts, newProtocolFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setNewProtocolFeeRecipient", newProtocolFeeRecipient)
}

// SetNewProtocolFeeRecipient is a paid mutator transaction binding the contract method 0x184ac28e.
//
// Solidity: function setNewProtocolFeeRecipient(address newProtocolFeeRecipient) returns()
func (_Bidderregistry *BidderregistrySession) SetNewProtocolFeeRecipient(newProtocolFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewProtocolFeeRecipient(&_Bidderregistry.TransactOpts, newProtocolFeeRecipient)
}

// SetNewProtocolFeeRecipient is a paid mutator transaction binding the contract method 0x184ac28e.
//
// Solidity: function setNewProtocolFeeRecipient(address newProtocolFeeRecipient) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetNewProtocolFeeRecipient(newProtocolFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewProtocolFeeRecipient(&_Bidderregistry.TransactOpts, newProtocolFeeRecipient)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Bidderregistry *BidderregistryTransactor) SetPreconfManager(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setPreconfManager", contractAddress)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Bidderregistry *BidderregistrySession) SetPreconfManager(contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetPreconfManager(&_Bidderregistry.TransactOpts, contractAddress)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetPreconfManager(contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetPreconfManager(&_Bidderregistry.TransactOpts, contractAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bidderregistry *BidderregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bidderregistry *BidderregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.TransferOwnership(&_Bidderregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bidderregistry *BidderregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.TransferOwnership(&_Bidderregistry.TransactOpts, newOwner)
}

// UnlockFunds is a paid mutator transaction binding the contract method 0x432e707b.
//
// Solidity: function unlockFunds(uint256 window, bytes32 bidID) returns()
func (_Bidderregistry *BidderregistryTransactor) UnlockFunds(opts *bind.TransactOpts, window *big.Int, bidID [32]byte) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "unlockFunds", window, bidID)
}

// UnlockFunds is a paid mutator transaction binding the contract method 0x432e707b.
//
// Solidity: function unlockFunds(uint256 window, bytes32 bidID) returns()
func (_Bidderregistry *BidderregistrySession) UnlockFunds(window *big.Int, bidID [32]byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.UnlockFunds(&_Bidderregistry.TransactOpts, window, bidID)
}

// UnlockFunds is a paid mutator transaction binding the contract method 0x432e707b.
//
// Solidity: function unlockFunds(uint256 window, bytes32 bidID) returns()
func (_Bidderregistry *BidderregistryTransactorSession) UnlockFunds(window *big.Int, bidID [32]byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.UnlockFunds(&_Bidderregistry.TransactOpts, window, bidID)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bidderregistry *BidderregistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bidderregistry *BidderregistrySession) Unpause() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Unpause(&_Bidderregistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bidderregistry *BidderregistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Unpause(&_Bidderregistry.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Bidderregistry *BidderregistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Bidderregistry *BidderregistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.UpgradeToAndCall(&_Bidderregistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.UpgradeToAndCall(&_Bidderregistry.TransactOpts, newImplementation, data)
}

// WithdrawBidderAmountFromWindow is a paid mutator transaction binding the contract method 0xa4bf023c.
//
// Solidity: function withdrawBidderAmountFromWindow(address bidder, uint256 window) returns()
func (_Bidderregistry *BidderregistryTransactor) WithdrawBidderAmountFromWindow(opts *bind.TransactOpts, bidder common.Address, window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "withdrawBidderAmountFromWindow", bidder, window)
}

// WithdrawBidderAmountFromWindow is a paid mutator transaction binding the contract method 0xa4bf023c.
//
// Solidity: function withdrawBidderAmountFromWindow(address bidder, uint256 window) returns()
func (_Bidderregistry *BidderregistrySession) WithdrawBidderAmountFromWindow(bidder common.Address, window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawBidderAmountFromWindow(&_Bidderregistry.TransactOpts, bidder, window)
}

// WithdrawBidderAmountFromWindow is a paid mutator transaction binding the contract method 0xa4bf023c.
//
// Solidity: function withdrawBidderAmountFromWindow(address bidder, uint256 window) returns()
func (_Bidderregistry *BidderregistryTransactorSession) WithdrawBidderAmountFromWindow(bidder common.Address, window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawBidderAmountFromWindow(&_Bidderregistry.TransactOpts, bidder, window)
}

// WithdrawFromWindows is a paid mutator transaction binding the contract method 0x6745206a.
//
// Solidity: function withdrawFromWindows(uint256[] windows) returns()
func (_Bidderregistry *BidderregistryTransactor) WithdrawFromWindows(opts *bind.TransactOpts, windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "withdrawFromWindows", windows)
}

// WithdrawFromWindows is a paid mutator transaction binding the contract method 0x6745206a.
//
// Solidity: function withdrawFromWindows(uint256[] windows) returns()
func (_Bidderregistry *BidderregistrySession) WithdrawFromWindows(windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawFromWindows(&_Bidderregistry.TransactOpts, windows)
}

// WithdrawFromWindows is a paid mutator transaction binding the contract method 0x6745206a.
//
// Solidity: function withdrawFromWindows(uint256[] windows) returns()
func (_Bidderregistry *BidderregistryTransactorSession) WithdrawFromWindows(windows []*big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawFromWindows(&_Bidderregistry.TransactOpts, windows)
}

// WithdrawProviderAmount is a paid mutator transaction binding the contract method 0x9a2dd5ba.
//
// Solidity: function withdrawProviderAmount(address provider) returns()
func (_Bidderregistry *BidderregistryTransactor) WithdrawProviderAmount(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "withdrawProviderAmount", provider)
}

// WithdrawProviderAmount is a paid mutator transaction binding the contract method 0x9a2dd5ba.
//
// Solidity: function withdrawProviderAmount(address provider) returns()
func (_Bidderregistry *BidderregistrySession) WithdrawProviderAmount(provider common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawProviderAmount(&_Bidderregistry.TransactOpts, provider)
}

// WithdrawProviderAmount is a paid mutator transaction binding the contract method 0x9a2dd5ba.
//
// Solidity: function withdrawProviderAmount(address provider) returns()
func (_Bidderregistry *BidderregistryTransactorSession) WithdrawProviderAmount(provider common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawProviderAmount(&_Bidderregistry.TransactOpts, provider)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Bidderregistry *BidderregistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Bidderregistry.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Bidderregistry *BidderregistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.Fallback(&_Bidderregistry.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Bidderregistry.Contract.Fallback(&_Bidderregistry.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Bidderregistry *BidderregistryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Bidderregistry *BidderregistrySession) Receive() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Receive(&_Bidderregistry.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) Receive() (*types.Transaction, error) {
	return _Bidderregistry.Contract.Receive(&_Bidderregistry.TransactOpts)
}

// BidderregistryBidderRegisteredIterator is returned from FilterBidderRegistered and is used to iterate over the raw logs and unpacked data for BidderRegistered events raised by the Bidderregistry contract.
type BidderregistryBidderRegisteredIterator struct {
	Event *BidderregistryBidderRegistered // Event containing the contract specifics and raw log

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
func (it *BidderregistryBidderRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryBidderRegistered)
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
		it.Event = new(BidderregistryBidderRegistered)
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
func (it *BidderregistryBidderRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryBidderRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryBidderRegistered represents a BidderRegistered event raised by the Bidderregistry contract.
type BidderregistryBidderRegistered struct {
	Bidder          common.Address
	DepositedAmount *big.Int
	WindowNumber    *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBidderRegistered is a free log retrieval operation binding the contract event 0x2ed10ffb7f7e5289e3bb91b8c3751388cb5d9b7f4533b9f0d59881a99822ddb3.
//
// Solidity: event BidderRegistered(address indexed bidder, uint256 indexed depositedAmount, uint256 indexed windowNumber)
func (_Bidderregistry *BidderregistryFilterer) FilterBidderRegistered(opts *bind.FilterOpts, bidder []common.Address, depositedAmount []*big.Int, windowNumber []*big.Int) (*BidderregistryBidderRegisteredIterator, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var depositedAmountRule []interface{}
	for _, depositedAmountItem := range depositedAmount {
		depositedAmountRule = append(depositedAmountRule, depositedAmountItem)
	}
	var windowNumberRule []interface{}
	for _, windowNumberItem := range windowNumber {
		windowNumberRule = append(windowNumberRule, windowNumberItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "BidderRegistered", bidderRule, depositedAmountRule, windowNumberRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryBidderRegisteredIterator{contract: _Bidderregistry.contract, event: "BidderRegistered", logs: logs, sub: sub}, nil
}

// WatchBidderRegistered is a free log subscription operation binding the contract event 0x2ed10ffb7f7e5289e3bb91b8c3751388cb5d9b7f4533b9f0d59881a99822ddb3.
//
// Solidity: event BidderRegistered(address indexed bidder, uint256 indexed depositedAmount, uint256 indexed windowNumber)
func (_Bidderregistry *BidderregistryFilterer) WatchBidderRegistered(opts *bind.WatchOpts, sink chan<- *BidderregistryBidderRegistered, bidder []common.Address, depositedAmount []*big.Int, windowNumber []*big.Int) (event.Subscription, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var depositedAmountRule []interface{}
	for _, depositedAmountItem := range depositedAmount {
		depositedAmountRule = append(depositedAmountRule, depositedAmountItem)
	}
	var windowNumberRule []interface{}
	for _, windowNumberItem := range windowNumber {
		windowNumberRule = append(windowNumberRule, windowNumberItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "BidderRegistered", bidderRule, depositedAmountRule, windowNumberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryBidderRegistered)
				if err := _Bidderregistry.contract.UnpackLog(event, "BidderRegistered", log); err != nil {
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

// ParseBidderRegistered is a log parse operation binding the contract event 0x2ed10ffb7f7e5289e3bb91b8c3751388cb5d9b7f4533b9f0d59881a99822ddb3.
//
// Solidity: event BidderRegistered(address indexed bidder, uint256 indexed depositedAmount, uint256 indexed windowNumber)
func (_Bidderregistry *BidderregistryFilterer) ParseBidderRegistered(log types.Log) (*BidderregistryBidderRegistered, error) {
	event := new(BidderregistryBidderRegistered)
	if err := _Bidderregistry.contract.UnpackLog(event, "BidderRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryBidderWithdrawalIterator is returned from FilterBidderWithdrawal and is used to iterate over the raw logs and unpacked data for BidderWithdrawal events raised by the Bidderregistry contract.
type BidderregistryBidderWithdrawalIterator struct {
	Event *BidderregistryBidderWithdrawal // Event containing the contract specifics and raw log

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
func (it *BidderregistryBidderWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryBidderWithdrawal)
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
		it.Event = new(BidderregistryBidderWithdrawal)
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
func (it *BidderregistryBidderWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryBidderWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryBidderWithdrawal represents a BidderWithdrawal event raised by the Bidderregistry contract.
type BidderregistryBidderWithdrawal struct {
	Bidder common.Address
	Window *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBidderWithdrawal is a free log retrieval operation binding the contract event 0x2be239cccec761cb15b4070dda36677f39cb05afba45c7419fe7e27ed2c90b29.
//
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 indexed window, uint256 indexed amount)
func (_Bidderregistry *BidderregistryFilterer) FilterBidderWithdrawal(opts *bind.FilterOpts, bidder []common.Address, window []*big.Int, amount []*big.Int) (*BidderregistryBidderWithdrawalIterator, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "BidderWithdrawal", bidderRule, windowRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryBidderWithdrawalIterator{contract: _Bidderregistry.contract, event: "BidderWithdrawal", logs: logs, sub: sub}, nil
}

// WatchBidderWithdrawal is a free log subscription operation binding the contract event 0x2be239cccec761cb15b4070dda36677f39cb05afba45c7419fe7e27ed2c90b29.
//
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 indexed window, uint256 indexed amount)
func (_Bidderregistry *BidderregistryFilterer) WatchBidderWithdrawal(opts *bind.WatchOpts, sink chan<- *BidderregistryBidderWithdrawal, bidder []common.Address, window []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "BidderWithdrawal", bidderRule, windowRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryBidderWithdrawal)
				if err := _Bidderregistry.contract.UnpackLog(event, "BidderWithdrawal", log); err != nil {
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

// ParseBidderWithdrawal is a log parse operation binding the contract event 0x2be239cccec761cb15b4070dda36677f39cb05afba45c7419fe7e27ed2c90b29.
//
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 indexed window, uint256 indexed amount)
func (_Bidderregistry *BidderregistryFilterer) ParseBidderWithdrawal(log types.Log) (*BidderregistryBidderWithdrawal, error) {
	event := new(BidderregistryBidderWithdrawal)
	if err := _Bidderregistry.contract.UnpackLog(event, "BidderWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryBlockTrackerUpdatedIterator is returned from FilterBlockTrackerUpdated and is used to iterate over the raw logs and unpacked data for BlockTrackerUpdated events raised by the Bidderregistry contract.
type BidderregistryBlockTrackerUpdatedIterator struct {
	Event *BidderregistryBlockTrackerUpdated // Event containing the contract specifics and raw log

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
func (it *BidderregistryBlockTrackerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryBlockTrackerUpdated)
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
		it.Event = new(BidderregistryBlockTrackerUpdated)
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
func (it *BidderregistryBlockTrackerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryBlockTrackerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryBlockTrackerUpdated represents a BlockTrackerUpdated event raised by the Bidderregistry contract.
type BidderregistryBlockTrackerUpdated struct {
	NewBlockTracker common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBlockTrackerUpdated is a free log retrieval operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Bidderregistry *BidderregistryFilterer) FilterBlockTrackerUpdated(opts *bind.FilterOpts, newBlockTracker []common.Address) (*BidderregistryBlockTrackerUpdatedIterator, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "BlockTrackerUpdated", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryBlockTrackerUpdatedIterator{contract: _Bidderregistry.contract, event: "BlockTrackerUpdated", logs: logs, sub: sub}, nil
}

// WatchBlockTrackerUpdated is a free log subscription operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Bidderregistry *BidderregistryFilterer) WatchBlockTrackerUpdated(opts *bind.WatchOpts, sink chan<- *BidderregistryBlockTrackerUpdated, newBlockTracker []common.Address) (event.Subscription, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "BlockTrackerUpdated", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryBlockTrackerUpdated)
				if err := _Bidderregistry.contract.UnpackLog(event, "BlockTrackerUpdated", log); err != nil {
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

// ParseBlockTrackerUpdated is a log parse operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Bidderregistry *BidderregistryFilterer) ParseBlockTrackerUpdated(log types.Log) (*BidderregistryBlockTrackerUpdated, error) {
	event := new(BidderregistryBlockTrackerUpdated)
	if err := _Bidderregistry.contract.UnpackLog(event, "BlockTrackerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryFeePayoutPeriodBlocksUpdatedIterator is returned from FilterFeePayoutPeriodBlocksUpdated and is used to iterate over the raw logs and unpacked data for FeePayoutPeriodBlocksUpdated events raised by the Bidderregistry contract.
type BidderregistryFeePayoutPeriodBlocksUpdatedIterator struct {
	Event *BidderregistryFeePayoutPeriodBlocksUpdated // Event containing the contract specifics and raw log

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
func (it *BidderregistryFeePayoutPeriodBlocksUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryFeePayoutPeriodBlocksUpdated)
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
		it.Event = new(BidderregistryFeePayoutPeriodBlocksUpdated)
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
func (it *BidderregistryFeePayoutPeriodBlocksUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryFeePayoutPeriodBlocksUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryFeePayoutPeriodBlocksUpdated represents a FeePayoutPeriodBlocksUpdated event raised by the Bidderregistry contract.
type BidderregistryFeePayoutPeriodBlocksUpdated struct {
	NewFeePayoutPeriodBlocks *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterFeePayoutPeriodBlocksUpdated is a free log retrieval operation binding the contract event 0x1b8b3f7fd7594ce5b7155b4c56b19bd6a1eac8c1ec5d941635acf104c8db3571.
//
// Solidity: event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks)
func (_Bidderregistry *BidderregistryFilterer) FilterFeePayoutPeriodBlocksUpdated(opts *bind.FilterOpts, newFeePayoutPeriodBlocks []*big.Int) (*BidderregistryFeePayoutPeriodBlocksUpdatedIterator, error) {

	var newFeePayoutPeriodBlocksRule []interface{}
	for _, newFeePayoutPeriodBlocksItem := range newFeePayoutPeriodBlocks {
		newFeePayoutPeriodBlocksRule = append(newFeePayoutPeriodBlocksRule, newFeePayoutPeriodBlocksItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FeePayoutPeriodBlocksUpdated", newFeePayoutPeriodBlocksRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFeePayoutPeriodBlocksUpdatedIterator{contract: _Bidderregistry.contract, event: "FeePayoutPeriodBlocksUpdated", logs: logs, sub: sub}, nil
}

// WatchFeePayoutPeriodBlocksUpdated is a free log subscription operation binding the contract event 0x1b8b3f7fd7594ce5b7155b4c56b19bd6a1eac8c1ec5d941635acf104c8db3571.
//
// Solidity: event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks)
func (_Bidderregistry *BidderregistryFilterer) WatchFeePayoutPeriodBlocksUpdated(opts *bind.WatchOpts, sink chan<- *BidderregistryFeePayoutPeriodBlocksUpdated, newFeePayoutPeriodBlocks []*big.Int) (event.Subscription, error) {

	var newFeePayoutPeriodBlocksRule []interface{}
	for _, newFeePayoutPeriodBlocksItem := range newFeePayoutPeriodBlocks {
		newFeePayoutPeriodBlocksRule = append(newFeePayoutPeriodBlocksRule, newFeePayoutPeriodBlocksItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FeePayoutPeriodBlocksUpdated", newFeePayoutPeriodBlocksRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryFeePayoutPeriodBlocksUpdated)
				if err := _Bidderregistry.contract.UnpackLog(event, "FeePayoutPeriodBlocksUpdated", log); err != nil {
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

// ParseFeePayoutPeriodBlocksUpdated is a log parse operation binding the contract event 0x1b8b3f7fd7594ce5b7155b4c56b19bd6a1eac8c1ec5d941635acf104c8db3571.
//
// Solidity: event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks)
func (_Bidderregistry *BidderregistryFilterer) ParseFeePayoutPeriodBlocksUpdated(log types.Log) (*BidderregistryFeePayoutPeriodBlocksUpdated, error) {
	event := new(BidderregistryFeePayoutPeriodBlocksUpdated)
	if err := _Bidderregistry.contract.UnpackLog(event, "FeePayoutPeriodBlocksUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryFeePercentUpdatedIterator is returned from FilterFeePercentUpdated and is used to iterate over the raw logs and unpacked data for FeePercentUpdated events raised by the Bidderregistry contract.
type BidderregistryFeePercentUpdatedIterator struct {
	Event *BidderregistryFeePercentUpdated // Event containing the contract specifics and raw log

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
func (it *BidderregistryFeePercentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryFeePercentUpdated)
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
		it.Event = new(BidderregistryFeePercentUpdated)
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
func (it *BidderregistryFeePercentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryFeePercentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryFeePercentUpdated represents a FeePercentUpdated event raised by the Bidderregistry contract.
type BidderregistryFeePercentUpdated struct {
	NewFeePercent uint16
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeePercentUpdated is a free log retrieval operation binding the contract event 0xcbefe12795a171d23fc26333ec4f03a7478bae6ec54d70e18ff15ebca679407c.
//
// Solidity: event FeePercentUpdated(uint16 indexed newFeePercent)
func (_Bidderregistry *BidderregistryFilterer) FilterFeePercentUpdated(opts *bind.FilterOpts, newFeePercent []uint16) (*BidderregistryFeePercentUpdatedIterator, error) {

	var newFeePercentRule []interface{}
	for _, newFeePercentItem := range newFeePercent {
		newFeePercentRule = append(newFeePercentRule, newFeePercentItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FeePercentUpdated", newFeePercentRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFeePercentUpdatedIterator{contract: _Bidderregistry.contract, event: "FeePercentUpdated", logs: logs, sub: sub}, nil
}

// WatchFeePercentUpdated is a free log subscription operation binding the contract event 0xcbefe12795a171d23fc26333ec4f03a7478bae6ec54d70e18ff15ebca679407c.
//
// Solidity: event FeePercentUpdated(uint16 indexed newFeePercent)
func (_Bidderregistry *BidderregistryFilterer) WatchFeePercentUpdated(opts *bind.WatchOpts, sink chan<- *BidderregistryFeePercentUpdated, newFeePercent []uint16) (event.Subscription, error) {

	var newFeePercentRule []interface{}
	for _, newFeePercentItem := range newFeePercent {
		newFeePercentRule = append(newFeePercentRule, newFeePercentItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FeePercentUpdated", newFeePercentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryFeePercentUpdated)
				if err := _Bidderregistry.contract.UnpackLog(event, "FeePercentUpdated", log); err != nil {
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

// ParseFeePercentUpdated is a log parse operation binding the contract event 0xcbefe12795a171d23fc26333ec4f03a7478bae6ec54d70e18ff15ebca679407c.
//
// Solidity: event FeePercentUpdated(uint16 indexed newFeePercent)
func (_Bidderregistry *BidderregistryFilterer) ParseFeePercentUpdated(log types.Log) (*BidderregistryFeePercentUpdated, error) {
	event := new(BidderregistryFeePercentUpdated)
	if err := _Bidderregistry.contract.UnpackLog(event, "FeePercentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryFeeTransferIterator is returned from FilterFeeTransfer and is used to iterate over the raw logs and unpacked data for FeeTransfer events raised by the Bidderregistry contract.
type BidderregistryFeeTransferIterator struct {
	Event *BidderregistryFeeTransfer // Event containing the contract specifics and raw log

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
func (it *BidderregistryFeeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryFeeTransfer)
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
		it.Event = new(BidderregistryFeeTransfer)
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
func (it *BidderregistryFeeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryFeeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryFeeTransfer represents a FeeTransfer event raised by the Bidderregistry contract.
type BidderregistryFeeTransfer struct {
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeTransfer is a free log retrieval operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Bidderregistry *BidderregistryFilterer) FilterFeeTransfer(opts *bind.FilterOpts, recipient []common.Address) (*BidderregistryFeeTransferIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFeeTransferIterator{contract: _Bidderregistry.contract, event: "FeeTransfer", logs: logs, sub: sub}, nil
}

// WatchFeeTransfer is a free log subscription operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Bidderregistry *BidderregistryFilterer) WatchFeeTransfer(opts *bind.WatchOpts, sink chan<- *BidderregistryFeeTransfer, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryFeeTransfer)
				if err := _Bidderregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
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

// ParseFeeTransfer is a log parse operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Bidderregistry *BidderregistryFilterer) ParseFeeTransfer(log types.Log) (*BidderregistryFeeTransfer, error) {
	event := new(BidderregistryFeeTransfer)
	if err := _Bidderregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryFundsRetrievedIterator is returned from FilterFundsRetrieved and is used to iterate over the raw logs and unpacked data for FundsRetrieved events raised by the Bidderregistry contract.
type BidderregistryFundsRetrievedIterator struct {
	Event *BidderregistryFundsRetrieved // Event containing the contract specifics and raw log

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
func (it *BidderregistryFundsRetrievedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryFundsRetrieved)
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
		it.Event = new(BidderregistryFundsRetrieved)
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
func (it *BidderregistryFundsRetrievedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryFundsRetrievedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryFundsRetrieved represents a FundsRetrieved event raised by the Bidderregistry contract.
type BidderregistryFundsRetrieved struct {
	CommitmentDigest [32]byte
	Bidder           common.Address
	Window           *big.Int
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterFundsRetrieved is a free log retrieval operation binding the contract event 0x4ee0e06b2d2e4d1f06e75df9f2bad2c919d860fbf843f3b1f12de3264471a102.
//
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 indexed window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) FilterFundsRetrieved(opts *bind.FilterOpts, commitmentDigest [][32]byte, bidder []common.Address, window []*big.Int) (*BidderregistryFundsRetrievedIterator, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FundsRetrieved", commitmentDigestRule, bidderRule, windowRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFundsRetrievedIterator{contract: _Bidderregistry.contract, event: "FundsRetrieved", logs: logs, sub: sub}, nil
}

// WatchFundsRetrieved is a free log subscription operation binding the contract event 0x4ee0e06b2d2e4d1f06e75df9f2bad2c919d860fbf843f3b1f12de3264471a102.
//
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 indexed window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) WatchFundsRetrieved(opts *bind.WatchOpts, sink chan<- *BidderregistryFundsRetrieved, commitmentDigest [][32]byte, bidder []common.Address, window []*big.Int) (event.Subscription, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var windowRule []interface{}
	for _, windowItem := range window {
		windowRule = append(windowRule, windowItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FundsRetrieved", commitmentDigestRule, bidderRule, windowRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryFundsRetrieved)
				if err := _Bidderregistry.contract.UnpackLog(event, "FundsRetrieved", log); err != nil {
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

// ParseFundsRetrieved is a log parse operation binding the contract event 0x4ee0e06b2d2e4d1f06e75df9f2bad2c919d860fbf843f3b1f12de3264471a102.
//
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 indexed window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) ParseFundsRetrieved(log types.Log) (*BidderregistryFundsRetrieved, error) {
	event := new(BidderregistryFundsRetrieved)
	if err := _Bidderregistry.contract.UnpackLog(event, "FundsRetrieved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryFundsRewardedIterator is returned from FilterFundsRewarded and is used to iterate over the raw logs and unpacked data for FundsRewarded events raised by the Bidderregistry contract.
type BidderregistryFundsRewardedIterator struct {
	Event *BidderregistryFundsRewarded // Event containing the contract specifics and raw log

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
func (it *BidderregistryFundsRewardedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryFundsRewarded)
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
		it.Event = new(BidderregistryFundsRewarded)
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
func (it *BidderregistryFundsRewardedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryFundsRewardedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryFundsRewarded represents a FundsRewarded event raised by the Bidderregistry contract.
type BidderregistryFundsRewarded struct {
	CommitmentDigest [32]byte
	Bidder           common.Address
	Provider         common.Address
	Window           *big.Int
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterFundsRewarded is a free log retrieval operation binding the contract event 0xd26f9e20ff994b4298fe22216ee15de6c9b7a46164d7a5509f2c4d065d8b408a.
//
// Solidity: event FundsRewarded(bytes32 indexed commitmentDigest, address indexed bidder, address indexed provider, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) FilterFundsRewarded(opts *bind.FilterOpts, commitmentDigest [][32]byte, bidder []common.Address, provider []common.Address) (*BidderregistryFundsRewardedIterator, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FundsRewarded", commitmentDigestRule, bidderRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFundsRewardedIterator{contract: _Bidderregistry.contract, event: "FundsRewarded", logs: logs, sub: sub}, nil
}

// WatchFundsRewarded is a free log subscription operation binding the contract event 0xd26f9e20ff994b4298fe22216ee15de6c9b7a46164d7a5509f2c4d065d8b408a.
//
// Solidity: event FundsRewarded(bytes32 indexed commitmentDigest, address indexed bidder, address indexed provider, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) WatchFundsRewarded(opts *bind.WatchOpts, sink chan<- *BidderregistryFundsRewarded, commitmentDigest [][32]byte, bidder []common.Address, provider []common.Address) (event.Subscription, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FundsRewarded", commitmentDigestRule, bidderRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryFundsRewarded)
				if err := _Bidderregistry.contract.UnpackLog(event, "FundsRewarded", log); err != nil {
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

// ParseFundsRewarded is a log parse operation binding the contract event 0xd26f9e20ff994b4298fe22216ee15de6c9b7a46164d7a5509f2c4d065d8b408a.
//
// Solidity: event FundsRewarded(bytes32 indexed commitmentDigest, address indexed bidder, address indexed provider, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) ParseFundsRewarded(log types.Log) (*BidderregistryFundsRewarded, error) {
	event := new(BidderregistryFundsRewarded)
	if err := _Bidderregistry.contract.UnpackLog(event, "FundsRewarded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Bidderregistry contract.
type BidderregistryInitializedIterator struct {
	Event *BidderregistryInitialized // Event containing the contract specifics and raw log

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
func (it *BidderregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryInitialized)
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
		it.Event = new(BidderregistryInitialized)
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
func (it *BidderregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryInitialized represents a Initialized event raised by the Bidderregistry contract.
type BidderregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Bidderregistry *BidderregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*BidderregistryInitializedIterator, error) {

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BidderregistryInitializedIterator{contract: _Bidderregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Bidderregistry *BidderregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BidderregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryInitialized)
				if err := _Bidderregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParseInitialized(log types.Log) (*BidderregistryInitialized, error) {
	event := new(BidderregistryInitialized)
	if err := _Bidderregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Bidderregistry contract.
type BidderregistryOwnershipTransferStartedIterator struct {
	Event *BidderregistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *BidderregistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryOwnershipTransferStarted)
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
		it.Event = new(BidderregistryOwnershipTransferStarted)
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
func (it *BidderregistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Bidderregistry contract.
type BidderregistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Bidderregistry *BidderregistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BidderregistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryOwnershipTransferStartedIterator{contract: _Bidderregistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Bidderregistry *BidderregistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *BidderregistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryOwnershipTransferStarted)
				if err := _Bidderregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*BidderregistryOwnershipTransferStarted, error) {
	event := new(BidderregistryOwnershipTransferStarted)
	if err := _Bidderregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bidderregistry contract.
type BidderregistryOwnershipTransferredIterator struct {
	Event *BidderregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BidderregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryOwnershipTransferred)
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
		it.Event = new(BidderregistryOwnershipTransferred)
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
func (it *BidderregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Bidderregistry contract.
type BidderregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bidderregistry *BidderregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BidderregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryOwnershipTransferredIterator{contract: _Bidderregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bidderregistry *BidderregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BidderregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryOwnershipTransferred)
				if err := _Bidderregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParseOwnershipTransferred(log types.Log) (*BidderregistryOwnershipTransferred, error) {
	event := new(BidderregistryOwnershipTransferred)
	if err := _Bidderregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Bidderregistry contract.
type BidderregistryPausedIterator struct {
	Event *BidderregistryPaused // Event containing the contract specifics and raw log

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
func (it *BidderregistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryPaused)
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
		it.Event = new(BidderregistryPaused)
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
func (it *BidderregistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryPaused represents a Paused event raised by the Bidderregistry contract.
type BidderregistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Bidderregistry *BidderregistryFilterer) FilterPaused(opts *bind.FilterOpts) (*BidderregistryPausedIterator, error) {

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &BidderregistryPausedIterator{contract: _Bidderregistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Bidderregistry *BidderregistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *BidderregistryPaused) (event.Subscription, error) {

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryPaused)
				if err := _Bidderregistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParsePaused(log types.Log) (*BidderregistryPaused, error) {
	event := new(BidderregistryPaused)
	if err := _Bidderregistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryPreconfManagerUpdatedIterator is returned from FilterPreconfManagerUpdated and is used to iterate over the raw logs and unpacked data for PreconfManagerUpdated events raised by the Bidderregistry contract.
type BidderregistryPreconfManagerUpdatedIterator struct {
	Event *BidderregistryPreconfManagerUpdated // Event containing the contract specifics and raw log

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
func (it *BidderregistryPreconfManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryPreconfManagerUpdated)
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
		it.Event = new(BidderregistryPreconfManagerUpdated)
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
func (it *BidderregistryPreconfManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryPreconfManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryPreconfManagerUpdated represents a PreconfManagerUpdated event raised by the Bidderregistry contract.
type BidderregistryPreconfManagerUpdated struct {
	NewPreconfManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPreconfManagerUpdated is a free log retrieval operation binding the contract event 0xba8b46bd4ed6a7ec49dc1a4758a5ffd0cbae99c172bbe007676fcb74fefa310f.
//
// Solidity: event PreconfManagerUpdated(address indexed newPreconfManager)
func (_Bidderregistry *BidderregistryFilterer) FilterPreconfManagerUpdated(opts *bind.FilterOpts, newPreconfManager []common.Address) (*BidderregistryPreconfManagerUpdatedIterator, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "PreconfManagerUpdated", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryPreconfManagerUpdatedIterator{contract: _Bidderregistry.contract, event: "PreconfManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchPreconfManagerUpdated is a free log subscription operation binding the contract event 0xba8b46bd4ed6a7ec49dc1a4758a5ffd0cbae99c172bbe007676fcb74fefa310f.
//
// Solidity: event PreconfManagerUpdated(address indexed newPreconfManager)
func (_Bidderregistry *BidderregistryFilterer) WatchPreconfManagerUpdated(opts *bind.WatchOpts, sink chan<- *BidderregistryPreconfManagerUpdated, newPreconfManager []common.Address) (event.Subscription, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "PreconfManagerUpdated", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryPreconfManagerUpdated)
				if err := _Bidderregistry.contract.UnpackLog(event, "PreconfManagerUpdated", log); err != nil {
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

// ParsePreconfManagerUpdated is a log parse operation binding the contract event 0xba8b46bd4ed6a7ec49dc1a4758a5ffd0cbae99c172bbe007676fcb74fefa310f.
//
// Solidity: event PreconfManagerUpdated(address indexed newPreconfManager)
func (_Bidderregistry *BidderregistryFilterer) ParsePreconfManagerUpdated(log types.Log) (*BidderregistryPreconfManagerUpdated, error) {
	event := new(BidderregistryPreconfManagerUpdated)
	if err := _Bidderregistry.contract.UnpackLog(event, "PreconfManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryProtocolFeeRecipientUpdatedIterator is returned from FilterProtocolFeeRecipientUpdated and is used to iterate over the raw logs and unpacked data for ProtocolFeeRecipientUpdated events raised by the Bidderregistry contract.
type BidderregistryProtocolFeeRecipientUpdatedIterator struct {
	Event *BidderregistryProtocolFeeRecipientUpdated // Event containing the contract specifics and raw log

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
func (it *BidderregistryProtocolFeeRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryProtocolFeeRecipientUpdated)
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
		it.Event = new(BidderregistryProtocolFeeRecipientUpdated)
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
func (it *BidderregistryProtocolFeeRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryProtocolFeeRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryProtocolFeeRecipientUpdated represents a ProtocolFeeRecipientUpdated event raised by the Bidderregistry contract.
type BidderregistryProtocolFeeRecipientUpdated struct {
	NewProtocolFeeRecipient common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeRecipientUpdated is a free log retrieval operation binding the contract event 0xc1b5345cce283376356748dc57f2dfa7120431d016fc7ca9ba641bc65f91411d.
//
// Solidity: event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient)
func (_Bidderregistry *BidderregistryFilterer) FilterProtocolFeeRecipientUpdated(opts *bind.FilterOpts, newProtocolFeeRecipient []common.Address) (*BidderregistryProtocolFeeRecipientUpdatedIterator, error) {

	var newProtocolFeeRecipientRule []interface{}
	for _, newProtocolFeeRecipientItem := range newProtocolFeeRecipient {
		newProtocolFeeRecipientRule = append(newProtocolFeeRecipientRule, newProtocolFeeRecipientItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "ProtocolFeeRecipientUpdated", newProtocolFeeRecipientRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryProtocolFeeRecipientUpdatedIterator{contract: _Bidderregistry.contract, event: "ProtocolFeeRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeRecipientUpdated is a free log subscription operation binding the contract event 0xc1b5345cce283376356748dc57f2dfa7120431d016fc7ca9ba641bc65f91411d.
//
// Solidity: event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient)
func (_Bidderregistry *BidderregistryFilterer) WatchProtocolFeeRecipientUpdated(opts *bind.WatchOpts, sink chan<- *BidderregistryProtocolFeeRecipientUpdated, newProtocolFeeRecipient []common.Address) (event.Subscription, error) {

	var newProtocolFeeRecipientRule []interface{}
	for _, newProtocolFeeRecipientItem := range newProtocolFeeRecipient {
		newProtocolFeeRecipientRule = append(newProtocolFeeRecipientRule, newProtocolFeeRecipientItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "ProtocolFeeRecipientUpdated", newProtocolFeeRecipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryProtocolFeeRecipientUpdated)
				if err := _Bidderregistry.contract.UnpackLog(event, "ProtocolFeeRecipientUpdated", log); err != nil {
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

// ParseProtocolFeeRecipientUpdated is a log parse operation binding the contract event 0xc1b5345cce283376356748dc57f2dfa7120431d016fc7ca9ba641bc65f91411d.
//
// Solidity: event ProtocolFeeRecipientUpdated(address indexed newProtocolFeeRecipient)
func (_Bidderregistry *BidderregistryFilterer) ParseProtocolFeeRecipientUpdated(log types.Log) (*BidderregistryProtocolFeeRecipientUpdated, error) {
	event := new(BidderregistryProtocolFeeRecipientUpdated)
	if err := _Bidderregistry.contract.UnpackLog(event, "ProtocolFeeRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Bidderregistry contract.
type BidderregistryUnpausedIterator struct {
	Event *BidderregistryUnpaused // Event containing the contract specifics and raw log

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
func (it *BidderregistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryUnpaused)
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
		it.Event = new(BidderregistryUnpaused)
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
func (it *BidderregistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryUnpaused represents a Unpaused event raised by the Bidderregistry contract.
type BidderregistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bidderregistry *BidderregistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*BidderregistryUnpausedIterator, error) {

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &BidderregistryUnpausedIterator{contract: _Bidderregistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bidderregistry *BidderregistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *BidderregistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryUnpaused)
				if err := _Bidderregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParseUnpaused(log types.Log) (*BidderregistryUnpaused, error) {
	event := new(BidderregistryUnpaused)
	if err := _Bidderregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BidderregistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Bidderregistry contract.
type BidderregistryUpgradedIterator struct {
	Event *BidderregistryUpgraded // Event containing the contract specifics and raw log

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
func (it *BidderregistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BidderregistryUpgraded)
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
		it.Event = new(BidderregistryUpgraded)
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
func (it *BidderregistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BidderregistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BidderregistryUpgraded represents a Upgraded event raised by the Bidderregistry contract.
type BidderregistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Bidderregistry *BidderregistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*BidderregistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryUpgradedIterator{contract: _Bidderregistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Bidderregistry *BidderregistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *BidderregistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BidderregistryUpgraded)
				if err := _Bidderregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Bidderregistry *BidderregistryFilterer) ParseUpgraded(log types.Log) (*BidderregistryUpgraded, error) {
	event := new(BidderregistryUpgraded)
	if err := _Bidderregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
