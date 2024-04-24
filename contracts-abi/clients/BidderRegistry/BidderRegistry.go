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
)

// BidderregistryMetaData contains all meta data concerning the Bidderregistry contract.
var BidderregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_minDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blockTracker\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"BidPayment\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIBidderRegistry.State\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OpenBid\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidderRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockTrackerContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBlockTracker\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositForSpecificWindow\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"feePercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipientAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDeposit\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFeeRecipientAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProviderAmount\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lockedFunds\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minDeposit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"preConfirmationsContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"protocolFeeAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerAmount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"retrieveFunds\",\"inputs\":[{\"name\":\"windowToSettle\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeePercent\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeeRecipient\",\"inputs\":[{\"name\":\"newFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPreconfirmationsContract\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockFunds\",\"inputs\":[{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidID\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawBidderAmountFromWindow\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"window\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFeeRecipientAmount\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawProtocolFee\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawProviderAmount\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BidderRegistered\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"depositedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"windowNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BidderWithdrawal\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsRetrieved\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsRewarded\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"window\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
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
	parsed, err := abi.JSON(strings.NewReader(BidderregistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

// BidPayment is a free data retrieval call binding the contract method 0x26a769e4.
//
// Solidity: function BidPayment(bytes32 ) view returns(address bidder, uint64 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistryCaller) BidPayment(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt uint64
	State  uint8
}, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "BidPayment", arg0)

	outstruct := new(struct {
		Bidder common.Address
		BidAmt uint64
		State  uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Bidder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.BidAmt = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.State = *abi.ConvertType(out[2], new(uint8)).(*uint8)

	return *outstruct, err

}

// BidPayment is a free data retrieval call binding the contract method 0x26a769e4.
//
// Solidity: function BidPayment(bytes32 ) view returns(address bidder, uint64 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistrySession) BidPayment(arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt uint64
	State  uint8
}, error) {
	return _Bidderregistry.Contract.BidPayment(&_Bidderregistry.CallOpts, arg0)
}

// BidPayment is a free data retrieval call binding the contract method 0x26a769e4.
//
// Solidity: function BidPayment(bytes32 ) view returns(address bidder, uint64 bidAmt, uint8 state)
func (_Bidderregistry *BidderregistryCallerSession) BidPayment(arg0 [32]byte) (struct {
	Bidder common.Address
	BidAmt uint64
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

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Bidderregistry *BidderregistryCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Bidderregistry *BidderregistrySession) FeeRecipient() (common.Address, error) {
	return _Bidderregistry.Contract.FeeRecipient(&_Bidderregistry.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) FeeRecipient() (common.Address, error) {
	return _Bidderregistry.Contract.FeeRecipient(&_Bidderregistry.CallOpts)
}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) FeeRecipientAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "feeRecipientAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) FeeRecipientAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.FeeRecipientAmount(&_Bidderregistry.CallOpts)
}

// FeeRecipientAmount is a free data retrieval call binding the contract method 0xe0ae4ebd.
//
// Solidity: function feeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) FeeRecipientAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.FeeRecipientAmount(&_Bidderregistry.CallOpts)
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

// GetFeeRecipientAmount is a free data retrieval call binding the contract method 0xc286f373.
//
// Solidity: function getFeeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) GetFeeRecipientAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "getFeeRecipientAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeeRecipientAmount is a free data retrieval call binding the contract method 0xc286f373.
//
// Solidity: function getFeeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) GetFeeRecipientAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.GetFeeRecipientAmount(&_Bidderregistry.CallOpts)
}

// GetFeeRecipientAmount is a free data retrieval call binding the contract method 0xc286f373.
//
// Solidity: function getFeeRecipientAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) GetFeeRecipientAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.GetFeeRecipientAmount(&_Bidderregistry.CallOpts)
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

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) MinDeposit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "minDeposit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) MinDeposit() (*big.Int, error) {
	return _Bidderregistry.Contract.MinDeposit(&_Bidderregistry.CallOpts)
}

// MinDeposit is a free data retrieval call binding the contract method 0x41b3d185.
//
// Solidity: function minDeposit() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) MinDeposit() (*big.Int, error) {
	return _Bidderregistry.Contract.MinDeposit(&_Bidderregistry.CallOpts)
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

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Bidderregistry *BidderregistryCaller) PreConfirmationsContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "preConfirmationsContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Bidderregistry *BidderregistrySession) PreConfirmationsContract() (common.Address, error) {
	return _Bidderregistry.Contract.PreConfirmationsContract(&_Bidderregistry.CallOpts)
}

// PreConfirmationsContract is a free data retrieval call binding the contract method 0x0de05a1e.
//
// Solidity: function preConfirmationsContract() view returns(address)
func (_Bidderregistry *BidderregistryCallerSession) PreConfirmationsContract() (common.Address, error) {
	return _Bidderregistry.Contract.PreConfirmationsContract(&_Bidderregistry.CallOpts)
}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCaller) ProtocolFeeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bidderregistry.contract.Call(opts, &out, "protocolFeeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Bidderregistry *BidderregistrySession) ProtocolFeeAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.ProtocolFeeAmount(&_Bidderregistry.CallOpts)
}

// ProtocolFeeAmount is a free data retrieval call binding the contract method 0x8ec9c93b.
//
// Solidity: function protocolFeeAmount() view returns(uint256)
func (_Bidderregistry *BidderregistryCallerSession) ProtocolFeeAmount() (*big.Int, error) {
	return _Bidderregistry.Contract.ProtocolFeeAmount(&_Bidderregistry.CallOpts)
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

// OpenBid is a paid mutator transaction binding the contract method 0xbb0e0878.
//
// Solidity: function OpenBid(bytes32 commitmentDigest, uint64 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistryTransactor) OpenBid(opts *bind.TransactOpts, commitmentDigest [32]byte, bid uint64, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "OpenBid", commitmentDigest, bid, bidder, blockNumber)
}

// OpenBid is a paid mutator transaction binding the contract method 0xbb0e0878.
//
// Solidity: function OpenBid(bytes32 commitmentDigest, uint64 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistrySession) OpenBid(commitmentDigest [32]byte, bid uint64, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.Contract.OpenBid(&_Bidderregistry.TransactOpts, commitmentDigest, bid, bidder, blockNumber)
}

// OpenBid is a paid mutator transaction binding the contract method 0xbb0e0878.
//
// Solidity: function OpenBid(bytes32 commitmentDigest, uint64 bid, address bidder, uint64 blockNumber) returns()
func (_Bidderregistry *BidderregistryTransactorSession) OpenBid(commitmentDigest [32]byte, bid uint64, bidder common.Address, blockNumber uint64) (*types.Transaction, error) {
	return _Bidderregistry.Contract.OpenBid(&_Bidderregistry.TransactOpts, commitmentDigest, bid, bidder, blockNumber)
}

// DepositForSpecificWindow is a paid mutator transaction binding the contract method 0xe5e4bf4c.
//
// Solidity: function depositForSpecificWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistryTransactor) DepositForSpecificWindow(opts *bind.TransactOpts, window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "depositForSpecificWindow", window)
}

// DepositForSpecificWindow is a paid mutator transaction binding the contract method 0xe5e4bf4c.
//
// Solidity: function depositForSpecificWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistrySession) DepositForSpecificWindow(window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForSpecificWindow(&_Bidderregistry.TransactOpts, window)
}

// DepositForSpecificWindow is a paid mutator transaction binding the contract method 0xe5e4bf4c.
//
// Solidity: function depositForSpecificWindow(uint256 window) payable returns()
func (_Bidderregistry *BidderregistryTransactorSession) DepositForSpecificWindow(window *big.Int) (*types.Transaction, error) {
	return _Bidderregistry.Contract.DepositForSpecificWindow(&_Bidderregistry.TransactOpts, window)
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

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Bidderregistry *BidderregistryTransactor) SetNewFeeRecipient(opts *bind.TransactOpts, newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setNewFeeRecipient", newFeeRecipient)
}

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Bidderregistry *BidderregistrySession) SetNewFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeeRecipient(&_Bidderregistry.TransactOpts, newFeeRecipient)
}

// SetNewFeeRecipient is a paid mutator transaction binding the contract method 0xa26652ea.
//
// Solidity: function setNewFeeRecipient(address newFeeRecipient) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetNewFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetNewFeeRecipient(&_Bidderregistry.TransactOpts, newFeeRecipient)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Bidderregistry *BidderregistryTransactor) SetPreconfirmationsContract(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "setPreconfirmationsContract", contractAddress)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Bidderregistry *BidderregistrySession) SetPreconfirmationsContract(contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetPreconfirmationsContract(&_Bidderregistry.TransactOpts, contractAddress)
}

// SetPreconfirmationsContract is a paid mutator transaction binding the contract method 0xf6c7e476.
//
// Solidity: function setPreconfirmationsContract(address contractAddress) returns()
func (_Bidderregistry *BidderregistryTransactorSession) SetPreconfirmationsContract(contractAddress common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.SetPreconfirmationsContract(&_Bidderregistry.TransactOpts, contractAddress)
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

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Bidderregistry *BidderregistryTransactor) WithdrawFeeRecipientAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "withdrawFeeRecipientAmount")
}

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Bidderregistry *BidderregistrySession) WithdrawFeeRecipientAmount() (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawFeeRecipientAmount(&_Bidderregistry.TransactOpts)
}

// WithdrawFeeRecipientAmount is a paid mutator transaction binding the contract method 0x7e5713d8.
//
// Solidity: function withdrawFeeRecipientAmount() returns()
func (_Bidderregistry *BidderregistryTransactorSession) WithdrawFeeRecipientAmount() (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawFeeRecipientAmount(&_Bidderregistry.TransactOpts)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address bidder) returns()
func (_Bidderregistry *BidderregistryTransactor) WithdrawProtocolFee(opts *bind.TransactOpts, bidder common.Address) (*types.Transaction, error) {
	return _Bidderregistry.contract.Transact(opts, "withdrawProtocolFee", bidder)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address bidder) returns()
func (_Bidderregistry *BidderregistrySession) WithdrawProtocolFee(bidder common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawProtocolFee(&_Bidderregistry.TransactOpts, bidder)
}

// WithdrawProtocolFee is a paid mutator transaction binding the contract method 0x668fb6dc.
//
// Solidity: function withdrawProtocolFee(address bidder) returns()
func (_Bidderregistry *BidderregistryTransactorSession) WithdrawProtocolFee(bidder common.Address) (*types.Transaction, error) {
	return _Bidderregistry.Contract.WithdrawProtocolFee(&_Bidderregistry.TransactOpts, bidder)
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
// Solidity: event BidderRegistered(address indexed bidder, uint256 depositedAmount, uint256 windowNumber)
func (_Bidderregistry *BidderregistryFilterer) FilterBidderRegistered(opts *bind.FilterOpts, bidder []common.Address) (*BidderregistryBidderRegisteredIterator, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "BidderRegistered", bidderRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryBidderRegisteredIterator{contract: _Bidderregistry.contract, event: "BidderRegistered", logs: logs, sub: sub}, nil
}

// WatchBidderRegistered is a free log subscription operation binding the contract event 0x2ed10ffb7f7e5289e3bb91b8c3751388cb5d9b7f4533b9f0d59881a99822ddb3.
//
// Solidity: event BidderRegistered(address indexed bidder, uint256 depositedAmount, uint256 windowNumber)
func (_Bidderregistry *BidderregistryFilterer) WatchBidderRegistered(opts *bind.WatchOpts, sink chan<- *BidderregistryBidderRegistered, bidder []common.Address) (event.Subscription, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "BidderRegistered", bidderRule)
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
// Solidity: event BidderRegistered(address indexed bidder, uint256 depositedAmount, uint256 windowNumber)
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
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) FilterBidderWithdrawal(opts *bind.FilterOpts, bidder []common.Address) (*BidderregistryBidderWithdrawalIterator, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "BidderWithdrawal", bidderRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryBidderWithdrawalIterator{contract: _Bidderregistry.contract, event: "BidderWithdrawal", logs: logs, sub: sub}, nil
}

// WatchBidderWithdrawal is a free log subscription operation binding the contract event 0x2be239cccec761cb15b4070dda36677f39cb05afba45c7419fe7e27ed2c90b29.
//
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) WatchBidderWithdrawal(opts *bind.WatchOpts, sink chan<- *BidderregistryBidderWithdrawal, bidder []common.Address) (event.Subscription, error) {

	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "BidderWithdrawal", bidderRule)
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
// Solidity: event BidderWithdrawal(address indexed bidder, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) ParseBidderWithdrawal(log types.Log) (*BidderregistryBidderWithdrawal, error) {
	event := new(BidderregistryBidderWithdrawal)
	if err := _Bidderregistry.contract.UnpackLog(event, "BidderWithdrawal", log); err != nil {
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
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) FilterFundsRetrieved(opts *bind.FilterOpts, commitmentDigest [][32]byte, bidder []common.Address) (*BidderregistryFundsRetrievedIterator, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.FilterLogs(opts, "FundsRetrieved", commitmentDigestRule, bidderRule)
	if err != nil {
		return nil, err
	}
	return &BidderregistryFundsRetrievedIterator{contract: _Bidderregistry.contract, event: "FundsRetrieved", logs: logs, sub: sub}, nil
}

// WatchFundsRetrieved is a free log subscription operation binding the contract event 0x4ee0e06b2d2e4d1f06e75df9f2bad2c919d860fbf843f3b1f12de3264471a102.
//
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 window, uint256 amount)
func (_Bidderregistry *BidderregistryFilterer) WatchFundsRetrieved(opts *bind.WatchOpts, sink chan<- *BidderregistryFundsRetrieved, commitmentDigest [][32]byte, bidder []common.Address) (event.Subscription, error) {

	var commitmentDigestRule []interface{}
	for _, commitmentDigestItem := range commitmentDigest {
		commitmentDigestRule = append(commitmentDigestRule, commitmentDigestItem)
	}
	var bidderRule []interface{}
	for _, bidderItem := range bidder {
		bidderRule = append(bidderRule, bidderItem)
	}

	logs, sub, err := _Bidderregistry.contract.WatchLogs(opts, "FundsRetrieved", commitmentDigestRule, bidderRule)
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
// Solidity: event FundsRetrieved(bytes32 indexed commitmentDigest, address indexed bidder, uint256 window, uint256 amount)
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
