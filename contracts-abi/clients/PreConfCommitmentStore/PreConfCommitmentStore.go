// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package preconfcommitmentstore

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

// PreConfCommitmentStorePreConfCommitment is an auto generated low-level Go binding around an user-defined struct.
type PreConfCommitmentStorePreConfCommitment struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	BlockCommitedAt     *big.Int
}

// PreconfcommitmentstoreMetaData contains all meta data concerning the Preconfcommitmentstore contract.
var PreconfcommitmentstoreMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_providerRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_bidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_PRECONF\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_BID_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_COMMITMENT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"_bytesToHexString\",\"inputs\":[{\"name\":\"_bytes\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"bidderRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBidderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockCommitedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentsCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBidHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockCommitedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"commitmentUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"blockCommitedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getCommitmentsByBlockNumber\",\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentsByCommitter\",\"inputs\":[{\"name\":\"commiter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPreConfHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_bidSignature\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTxnHashFromCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initiateReward\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateSlash\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lastProcessedBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProviderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeCommitment\",\"inputs\":[{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockBidFunds\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBidderRegistry\",\"inputs\":[{\"name\":\"newBidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOracle\",\"inputs\":[{\"name\":\"newOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProviderRegistry\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyBid\",\"inputs\":[{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recoveredAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyPreConfCommitment\",\"inputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"preConfHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commiterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignatureVerified\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false}]",
}

// PreconfcommitmentstoreABI is the input ABI used to generate the binding from.
// Deprecated: Use PreconfcommitmentstoreMetaData.ABI instead.
var PreconfcommitmentstoreABI = PreconfcommitmentstoreMetaData.ABI

// Preconfcommitmentstore is an auto generated Go binding around an Ethereum contract.
type Preconfcommitmentstore struct {
	PreconfcommitmentstoreCaller     // Read-only binding to the contract
	PreconfcommitmentstoreTransactor // Write-only binding to the contract
	PreconfcommitmentstoreFilterer   // Log filterer for contract events
}

// PreconfcommitmentstoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreconfcommitmentstoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreconfcommitmentstoreSession struct {
	Contract     *Preconfcommitmentstore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PreconfcommitmentstoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreconfcommitmentstoreCallerSession struct {
	Contract *PreconfcommitmentstoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// PreconfcommitmentstoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreconfcommitmentstoreTransactorSession struct {
	Contract     *PreconfcommitmentstoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// PreconfcommitmentstoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreconfcommitmentstoreRaw struct {
	Contract *Preconfcommitmentstore // Generic contract binding to access the raw methods on
}

// PreconfcommitmentstoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreCallerRaw struct {
	Contract *PreconfcommitmentstoreCaller // Generic read-only contract binding to access the raw methods on
}

// PreconfcommitmentstoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreTransactorRaw struct {
	Contract *PreconfcommitmentstoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreconfcommitmentstore creates a new instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstore(address common.Address, backend bind.ContractBackend) (*Preconfcommitmentstore, error) {
	contract, err := bindPreconfcommitmentstore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Preconfcommitmentstore{PreconfcommitmentstoreCaller: PreconfcommitmentstoreCaller{contract: contract}, PreconfcommitmentstoreTransactor: PreconfcommitmentstoreTransactor{contract: contract}, PreconfcommitmentstoreFilterer: PreconfcommitmentstoreFilterer{contract: contract}}, nil
}

// NewPreconfcommitmentstoreCaller creates a new read-only instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreCaller(address common.Address, caller bind.ContractCaller) (*PreconfcommitmentstoreCaller, error) {
	contract, err := bindPreconfcommitmentstore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreCaller{contract: contract}, nil
}

// NewPreconfcommitmentstoreTransactor creates a new write-only instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreTransactor(address common.Address, transactor bind.ContractTransactor) (*PreconfcommitmentstoreTransactor, error) {
	contract, err := bindPreconfcommitmentstore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreTransactor{contract: contract}, nil
}

// NewPreconfcommitmentstoreFilterer creates a new log filterer instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreFilterer(address common.Address, filterer bind.ContractFilterer) (*PreconfcommitmentstoreFilterer, error) {
	contract, err := bindPreconfcommitmentstore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreFilterer{contract: contract}, nil
}

// bindPreconfcommitmentstore binds a generic wrapper to an already deployed contract.
func bindPreconfcommitmentstore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreconfcommitmentstoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfcommitmentstore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) DOMAINSEPARATORBID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "DOMAIN_SEPARATOR_BID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORBID(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORBID(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) DOMAINSEPARATORPRECONF(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "DOMAIN_SEPARATOR_PRECONF")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORPRECONF(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORPRECONF(&_Preconfcommitmentstore.CallOpts)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) EIP712BIDTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "EIP712_BID_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712BIDTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712BIDTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) EIP712COMMITMENTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "EIP712_COMMITMENT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BytesToHexString(opts *bind.CallOpts, _bytes []byte) (string, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "_bytesToHexString", _bytes)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BytesToHexString(_bytes []byte) (string, error) {
	return _Preconfcommitmentstore.Contract.BytesToHexString(&_Preconfcommitmentstore.CallOpts, _bytes)
}

// BytesToHexString is a free data retrieval call binding the contract method 0xca64db2e.
//
// Solidity: function _bytesToHexString(bytes _bytes) pure returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BytesToHexString(_bytes []byte) (string, error) {
	return _Preconfcommitmentstore.Contract.BytesToHexString(&_Preconfcommitmentstore.CallOpts, _bytes)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BidderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "bidderRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BidderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BidderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BidderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BidderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BlockCommitments(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "blockCommitments", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BlockCommitments(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.BlockCommitments(&_Preconfcommitmentstore.CallOpts, arg0, arg1)
}

// BlockCommitments is a free data retrieval call binding the contract method 0x159efb47.
//
// Solidity: function blockCommitments(uint256 , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BlockCommitments(arg0 *big.Int, arg1 *big.Int) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.BlockCommitments(&_Preconfcommitmentstore.CallOpts, arg0, arg1)
}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) CommitmentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitmentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) CommitmentCount() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentCount(&_Preconfcommitmentstore.CallOpts)
}

// CommitmentCount is a free data retrieval call binding the contract method 0xc44956d1.
//
// Solidity: function commitmentCount() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) CommitmentCount() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentCount(&_Preconfcommitmentstore.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint256 blockCommitedAt)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	BlockCommitedAt     *big.Int
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		CommitmentUsed      bool
		Bidder              common.Address
		Commiter            common.Address
		Bid                 uint64
		BlockNumber         uint64
		BidHash             [32]byte
		DecayStartTimeStamp uint64
		DecayEndTimeStamp   uint64
		TxnHash             string
		CommitmentHash      [32]byte
		BidSignature        []byte
		CommitmentSignature []byte
		BlockCommitedAt     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CommitmentUsed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Bidder = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Commiter = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Bid = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.BlockNumber = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.BidHash = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.DecayStartTimeStamp = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.DecayEndTimeStamp = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.TxnHash = *abi.ConvertType(out[8], new(string)).(*string)
	outstruct.CommitmentHash = *abi.ConvertType(out[9], new([32]byte)).(*[32]byte)
	outstruct.BidSignature = *abi.ConvertType(out[10], new([]byte)).(*[]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[11], new([]byte)).(*[]byte)
	outstruct.BlockCommitedAt = *abi.ConvertType(out[12], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint256 blockCommitedAt)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Commitments(arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	BlockCommitedAt     *big.Int
}, error) {
	return _Preconfcommitmentstore.Contract.Commitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(bool commitmentUsed, address bidder, address commiter, uint64 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint256 blockCommitedAt)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) Commitments(arg0 [32]byte) (struct {
	CommitmentUsed      bool
	Bidder              common.Address
	Commiter            common.Address
	Bid                 uint64
	BlockNumber         uint64
	BidHash             [32]byte
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	BlockCommitedAt     *big.Int
}, error) {
	return _Preconfcommitmentstore.Contract.Commitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) CommitmentsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitmentsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentsCount(&_Preconfcommitmentstore.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentsCount(&_Preconfcommitmentstore.CallOpts, arg0)
}

// GetBidHash is a free data retrieval call binding the contract method 0x7b2111f6.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetBidHash(opts *bind.CallOpts, _txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getBidHash", _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBidHash is a free data retrieval call binding the contract method 0x7b2111f6.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetBidHash(_txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetBidHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)
}

// GetBidHash is a free data retrieval call binding the contract method 0x7b2111f6.
//
// Solidity: function getBidHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetBidHash(_txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetBidHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitment", commitmentIndex)

	if err != nil {
		return *new(PreConfCommitmentStorePreConfCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(PreConfCommitmentStorePreConfCommitment)).(*PreConfCommitmentStorePreConfCommitment)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256))
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitment(commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitment(commitmentIndex [32]byte) (PreConfCommitmentStorePreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xebdb90e4.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitmentIndex(opts *bind.CallOpts, commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xebdb90e4.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitmentIndex(commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xebdb90e4.
//
// Solidity: function getCommitmentIndex((bool,address,address,uint64,uint64,bytes32,uint64,uint64,string,bytes32,bytes,bytes,uint256) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitmentIndex(commitment PreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitmentsByBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitmentsByBlockNumber", blockNumber)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitmentsByBlockNumber(blockNumber *big.Int) ([][32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentsByBlockNumber(&_Preconfcommitmentstore.CallOpts, blockNumber)
}

// GetCommitmentsByBlockNumber is a free data retrieval call binding the contract method 0x82da12de.
//
// Solidity: function getCommitmentsByBlockNumber(uint256 blockNumber) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitmentsByBlockNumber(blockNumber *big.Int) ([][32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentsByBlockNumber(&_Preconfcommitmentstore.CallOpts, blockNumber)
}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitmentsByCommitter(opts *bind.CallOpts, commiter common.Address) ([][32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitmentsByCommitter", commiter)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitmentsByCommitter(commiter common.Address) ([][32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentsByCommitter(&_Preconfcommitmentstore.CallOpts, commiter)
}

// GetCommitmentsByCommitter is a free data retrieval call binding the contract method 0xac8c8a0e.
//
// Solidity: function getCommitmentsByCommitter(address commiter) view returns(bytes32[])
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitmentsByCommitter(commiter common.Address) ([][32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentsByCommitter(&_Preconfcommitmentstore.CallOpts, commiter)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0x85358288.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetPreConfHash(opts *bind.CallOpts, _txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getPreConfHash", _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPreConfHash is a free data retrieval call binding the contract method 0x85358288.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetPreConfHash(_txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetPreConfHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0x85358288.
//
// Solidity: function getPreConfHash(string _txnHash, uint64 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetPreConfHash(_txnHash string, _bid uint64, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetPreConfHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getTxnHashFromCommitment", commitmentIndex)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfcommitmentstore.Contract.GetTxnHashFromCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfcommitmentstore.Contract.GetTxnHashFromCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) LastProcessedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "lastProcessedBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) LastProcessedBlock() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.LastProcessedBlock(&_Preconfcommitmentstore.CallOpts)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) LastProcessedBlock() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.LastProcessedBlock(&_Preconfcommitmentstore.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Oracle() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Oracle(&_Preconfcommitmentstore.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) Oracle() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Oracle(&_Preconfcommitmentstore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Owner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Owner(&_Preconfcommitmentstore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) Owner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Owner(&_Preconfcommitmentstore.CallOpts)
}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) ProviderCommitments(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "providerCommitments", arg0, arg1)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) ProviderCommitments(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.ProviderCommitments(&_Preconfcommitmentstore.CallOpts, arg0, arg1)
}

// ProviderCommitments is a free data retrieval call binding the contract method 0x91b51cda.
//
// Solidity: function providerCommitments(address , uint256 ) view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) ProviderCommitments(arg0 common.Address, arg1 *big.Int) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.ProviderCommitments(&_Preconfcommitmentstore.CallOpts, arg0, arg1)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) ProviderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "providerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) ProviderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.ProviderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) ProviderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.ProviderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// VerifyBid is a free data retrieval call binding the contract method 0x45216cb6.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) VerifyBid(opts *bind.CallOpts, bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "verifyBid", bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, bidSignature)

	outstruct := new(struct {
		MessageDigest    [32]byte
		RecoveredAddress common.Address
		Stake            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MessageDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RecoveredAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Stake = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// VerifyBid is a free data retrieval call binding the contract method 0x45216cb6.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) VerifyBid(bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyBid(&_Preconfcommitmentstore.CallOpts, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, bidSignature)
}

// VerifyBid is a free data retrieval call binding the contract method 0x45216cb6.
//
// Solidity: function verifyBid(uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, bytes bidSignature) view returns(bytes32 messageDigest, address recoveredAddress, uint256 stake)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) VerifyBid(bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
	Stake            *big.Int
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyBid(&_Preconfcommitmentstore.CallOpts, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, bidSignature)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x9a16af55.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) VerifyPreConfCommitment(opts *bind.CallOpts, txnHash string, bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "verifyPreConfCommitment", txnHash, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, bidHash, bidSignature, commitmentSignature)

	outstruct := new(struct {
		PreConfHash     [32]byte
		CommiterAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PreConfHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CommiterAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x9a16af55.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) VerifyPreConfCommitment(txnHash string, bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyPreConfCommitment(&_Preconfcommitmentstore.CallOpts, txnHash, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, bidHash, bidSignature, commitmentSignature)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0x9a16af55.
//
// Solidity: function verifyPreConfCommitment(string txnHash, uint64 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes32 bidHash, bytes bidSignature, bytes commitmentSignature) view returns(bytes32 preConfHash, address commiterAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) VerifyPreConfCommitment(txnHash string, bid uint64, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidHash [32]byte, bidSignature []byte, commitmentSignature []byte) (struct {
	PreConfHash     [32]byte
	CommiterAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyPreConfCommitment(&_Preconfcommitmentstore.CallOpts, txnHash, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, bidHash, bidSignature, commitmentSignature)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) InitiateReward(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "initiateReward", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateReward(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateReward(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) InitiateSlash(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "initiateSlash", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateSlash(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateSlash(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.RenounceOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.RenounceOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0xbfd21439.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) StoreCommitment(opts *bind.TransactOpts, bid uint64, blockNumber uint64, txnHash string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "storeCommitment", bid, blockNumber, txnHash, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0xbfd21439.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) StoreCommitment(bid uint64, blockNumber uint64, txnHash string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.StoreCommitment(&_Preconfcommitmentstore.TransactOpts, bid, blockNumber, txnHash, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature)
}

// StoreCommitment is a paid mutator transaction binding the contract method 0xbfd21439.
//
// Solidity: function storeCommitment(uint64 bid, uint64 blockNumber, string txnHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) StoreCommitment(bid uint64, blockNumber uint64, txnHash string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.StoreCommitment(&_Preconfcommitmentstore.TransactOpts, bid, blockNumber, txnHash, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.TransferOwnership(&_Preconfcommitmentstore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.TransferOwnership(&_Preconfcommitmentstore.TransactOpts, newOwner)
}

// UnlockBidFunds is a paid mutator transaction binding the contract method 0x7b37732f.
//
// Solidity: function unlockBidFunds(bytes32 commitmentDigest) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UnlockBidFunds(opts *bind.TransactOpts, commitmentDigest [32]byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "unlockBidFunds", commitmentDigest)
}

// UnlockBidFunds is a paid mutator transaction binding the contract method 0x7b37732f.
//
// Solidity: function unlockBidFunds(bytes32 commitmentDigest) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UnlockBidFunds(commitmentDigest [32]byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UnlockBidFunds(&_Preconfcommitmentstore.TransactOpts, commitmentDigest)
}

// UnlockBidFunds is a paid mutator transaction binding the contract method 0x7b37732f.
//
// Solidity: function unlockBidFunds(bytes32 commitmentDigest) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UnlockBidFunds(commitmentDigest [32]byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UnlockBidFunds(&_Preconfcommitmentstore.TransactOpts, commitmentDigest)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateBidderRegistry(opts *bind.TransactOpts, newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateBidderRegistry", newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateBidderRegistry(&_Preconfcommitmentstore.TransactOpts, newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateBidderRegistry(&_Preconfcommitmentstore.TransactOpts, newBidderRegistry)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateOracle(opts *bind.TransactOpts, newOracle common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateOracle", newOracle)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateOracle(&_Preconfcommitmentstore.TransactOpts, newOracle)
}

// UpdateOracle is a paid mutator transaction binding the contract method 0x1cb44dfc.
//
// Solidity: function updateOracle(address newOracle) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateOracle(newOracle common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateOracle(&_Preconfcommitmentstore.TransactOpts, newOracle)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateProviderRegistry(opts *bind.TransactOpts, newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateProviderRegistry", newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateProviderRegistry(&_Preconfcommitmentstore.TransactOpts, newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateProviderRegistry(&_Preconfcommitmentstore.TransactOpts, newProviderRegistry)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Fallback(&_Preconfcommitmentstore.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Fallback(&_Preconfcommitmentstore.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Receive() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Receive(&_Preconfcommitmentstore.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) Receive() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Receive(&_Preconfcommitmentstore.TransactOpts)
}

// PreconfcommitmentstoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferredIterator struct {
	Event *PreconfcommitmentstoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreOwnershipTransferred)
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
		it.Event = new(PreconfcommitmentstoreOwnershipTransferred)
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
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreOwnershipTransferred represents a OwnershipTransferred event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfcommitmentstoreOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreOwnershipTransferredIterator{contract: _Preconfcommitmentstore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreOwnershipTransferred)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseOwnershipTransferred(log types.Log) (*PreconfcommitmentstoreOwnershipTransferred, error) {
	event := new(PreconfcommitmentstoreOwnershipTransferred)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreSignatureVerifiedIterator is returned from FilterSignatureVerified and is used to iterate over the raw logs and unpacked data for SignatureVerified events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreSignatureVerifiedIterator struct {
	Event *PreconfcommitmentstoreSignatureVerified // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreSignatureVerified)
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
		it.Event = new(PreconfcommitmentstoreSignatureVerified)
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
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreSignatureVerified represents a SignatureVerified event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreSignatureVerified struct {
	Signer      common.Address
	TxnHash     string
	Bid         uint64
	BlockNumber uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSignatureVerified is a free log retrieval operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterSignatureVerified(opts *bind.FilterOpts, signer []common.Address, bid []uint64) (*PreconfcommitmentstoreSignatureVerifiedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreSignatureVerifiedIterator{contract: _Preconfcommitmentstore.contract, event: "SignatureVerified", logs: logs, sub: sub}, nil
}

// WatchSignatureVerified is a free log subscription operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchSignatureVerified(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreSignatureVerified, signer []common.Address, bid []uint64) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreSignatureVerified)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
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

// ParseSignatureVerified is a log parse operation binding the contract event 0x48db0394d84b81a6f3cb6c61ea2dceff3cad797a9b889fe499fc051f08969c4d.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, uint64 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseSignatureVerified(log types.Log) (*PreconfcommitmentstoreSignatureVerified, error) {
	event := new(PreconfcommitmentstoreSignatureVerified)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
