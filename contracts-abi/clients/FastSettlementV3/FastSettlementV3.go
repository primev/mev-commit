// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package fastsettlementv3

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

// IFastSettlementV3Intent is an auto generated low-level Go binding around an user-defined struct.
type IFastSettlementV3Intent struct {
	User        common.Address
	InputToken  common.Address
	OutputToken common.Address
	InputAmt    *big.Int
	UserAmtOut  *big.Int
	Recipient   common.Address
	Deadline    *big.Int
	Nonce       *big.Int
}

// IFastSettlementV3SwapCall is an auto generated low-level Go binding around an user-defined struct.
type IFastSettlementV3SwapCall struct {
	To    common.Address
	Value *big.Int
	Data  []byte
}

// Fastsettlementv3MetaData contains all meta data concerning the Fastsettlementv3 contract.
var Fastsettlementv3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_permit2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_weth\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"INTENT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PERMIT2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPermit2\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WETH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIWETH\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WITNESS_TYPE_STRING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowedSwapTargets\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeWithETH\",\"inputs\":[{\"name\":\"intent\",\"type\":\"tuple\",\"internalType\":\"structIFastSettlementV3.Intent\",\"components\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"inputToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"outputToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"inputAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userAmtOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"swapData\",\"type\":\"tuple\",\"internalType\":\"structIFastSettlementV3.SwapCall\",\"components\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"surplus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"executeWithPermit\",\"inputs\":[{\"name\":\"intent\",\"type\":\"tuple\",\"internalType\":\"structIFastSettlementV3.Intent\",\"components\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"inputToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"outputToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"inputAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userAmtOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"swapData\",\"type\":\"tuple\",\"internalType\":\"structIFastSettlementV3.SwapCall\",\"components\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"surplus\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"executor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_executor\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_initialSwapTargets\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rescueTokens\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setExecutor\",\"inputs\":[{\"name\":\"_newExecutor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapTargets\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"allowed\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTreasury\",\"inputs\":[{\"name\":\"_newTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"ExecutorUpdated\",\"inputs\":[{\"name\":\"oldExecutor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newExecutor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"IntentExecuted\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"inputToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"outputToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"inputAmt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"userAmtOut\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"received\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"surplus\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapTargetsUpdated\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"allowed\",\"type\":\"bool[]\",\"indexed\":false,\"internalType\":\"bool[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"oldTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadCallTarget\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadExecutor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadInputAmt\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadInputToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadNonce\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadTreasury\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadUserAmtOut\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedETHInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientOut\",\"inputs\":[{\"name\":\"received\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"userAmtOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"IntentExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidETHAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPermit2\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidWETH\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedCaller\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedExecutor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnauthorizedSwapTarget\",\"inputs\":[]}]",
}

// Fastsettlementv3ABI is the input ABI used to generate the binding from.
// Deprecated: Use Fastsettlementv3MetaData.ABI instead.
var Fastsettlementv3ABI = Fastsettlementv3MetaData.ABI

// Fastsettlementv3 is an auto generated Go binding around an Ethereum contract.
type Fastsettlementv3 struct {
	Fastsettlementv3Caller     // Read-only binding to the contract
	Fastsettlementv3Transactor // Write-only binding to the contract
	Fastsettlementv3Filterer   // Log filterer for contract events
}

// Fastsettlementv3Caller is an auto generated read-only Go binding around an Ethereum contract.
type Fastsettlementv3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Fastsettlementv3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Fastsettlementv3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Fastsettlementv3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Fastsettlementv3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Fastsettlementv3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Fastsettlementv3Session struct {
	Contract     *Fastsettlementv3 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Fastsettlementv3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Fastsettlementv3CallerSession struct {
	Contract *Fastsettlementv3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// Fastsettlementv3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Fastsettlementv3TransactorSession struct {
	Contract     *Fastsettlementv3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// Fastsettlementv3Raw is an auto generated low-level Go binding around an Ethereum contract.
type Fastsettlementv3Raw struct {
	Contract *Fastsettlementv3 // Generic contract binding to access the raw methods on
}

// Fastsettlementv3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Fastsettlementv3CallerRaw struct {
	Contract *Fastsettlementv3Caller // Generic read-only contract binding to access the raw methods on
}

// Fastsettlementv3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Fastsettlementv3TransactorRaw struct {
	Contract *Fastsettlementv3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewFastsettlementv3 creates a new instance of Fastsettlementv3, bound to a specific deployed contract.
func NewFastsettlementv3(address common.Address, backend bind.ContractBackend) (*Fastsettlementv3, error) {
	contract, err := bindFastsettlementv3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3{Fastsettlementv3Caller: Fastsettlementv3Caller{contract: contract}, Fastsettlementv3Transactor: Fastsettlementv3Transactor{contract: contract}, Fastsettlementv3Filterer: Fastsettlementv3Filterer{contract: contract}}, nil
}

// NewFastsettlementv3Caller creates a new read-only instance of Fastsettlementv3, bound to a specific deployed contract.
func NewFastsettlementv3Caller(address common.Address, caller bind.ContractCaller) (*Fastsettlementv3Caller, error) {
	contract, err := bindFastsettlementv3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3Caller{contract: contract}, nil
}

// NewFastsettlementv3Transactor creates a new write-only instance of Fastsettlementv3, bound to a specific deployed contract.
func NewFastsettlementv3Transactor(address common.Address, transactor bind.ContractTransactor) (*Fastsettlementv3Transactor, error) {
	contract, err := bindFastsettlementv3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3Transactor{contract: contract}, nil
}

// NewFastsettlementv3Filterer creates a new log filterer instance of Fastsettlementv3, bound to a specific deployed contract.
func NewFastsettlementv3Filterer(address common.Address, filterer bind.ContractFilterer) (*Fastsettlementv3Filterer, error) {
	contract, err := bindFastsettlementv3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3Filterer{contract: contract}, nil
}

// bindFastsettlementv3 binds a generic wrapper to an already deployed contract.
func bindFastsettlementv3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Fastsettlementv3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fastsettlementv3 *Fastsettlementv3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fastsettlementv3.Contract.Fastsettlementv3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fastsettlementv3 *Fastsettlementv3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Fastsettlementv3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fastsettlementv3 *Fastsettlementv3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Fastsettlementv3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Fastsettlementv3 *Fastsettlementv3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Fastsettlementv3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Fastsettlementv3 *Fastsettlementv3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Fastsettlementv3 *Fastsettlementv3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.contract.Transact(opts, method, params...)
}

// INTENTTYPEHASH is a free data retrieval call binding the contract method 0xb082a274.
//
// Solidity: function INTENT_TYPEHASH() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3Caller) INTENTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "INTENT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// INTENTTYPEHASH is a free data retrieval call binding the contract method 0xb082a274.
//
// Solidity: function INTENT_TYPEHASH() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3Session) INTENTTYPEHASH() ([32]byte, error) {
	return _Fastsettlementv3.Contract.INTENTTYPEHASH(&_Fastsettlementv3.CallOpts)
}

// INTENTTYPEHASH is a free data retrieval call binding the contract method 0xb082a274.
//
// Solidity: function INTENT_TYPEHASH() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) INTENTTYPEHASH() ([32]byte, error) {
	return _Fastsettlementv3.Contract.INTENTTYPEHASH(&_Fastsettlementv3.CallOpts)
}

// PERMIT2 is a free data retrieval call binding the contract method 0x6afdd850.
//
// Solidity: function PERMIT2() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) PERMIT2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "PERMIT2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PERMIT2 is a free data retrieval call binding the contract method 0x6afdd850.
//
// Solidity: function PERMIT2() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) PERMIT2() (common.Address, error) {
	return _Fastsettlementv3.Contract.PERMIT2(&_Fastsettlementv3.CallOpts)
}

// PERMIT2 is a free data retrieval call binding the contract method 0x6afdd850.
//
// Solidity: function PERMIT2() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) PERMIT2() (common.Address, error) {
	return _Fastsettlementv3.Contract.PERMIT2(&_Fastsettlementv3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _Fastsettlementv3.Contract.UPGRADEINTERFACEVERSION(&_Fastsettlementv3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Fastsettlementv3.Contract.UPGRADEINTERFACEVERSION(&_Fastsettlementv3.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) WETH(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "WETH")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) WETH() (common.Address, error) {
	return _Fastsettlementv3.Contract.WETH(&_Fastsettlementv3.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) WETH() (common.Address, error) {
	return _Fastsettlementv3.Contract.WETH(&_Fastsettlementv3.CallOpts)
}

// WITNESSTYPESTRING is a free data retrieval call binding the contract method 0x156e2152.
//
// Solidity: function WITNESS_TYPE_STRING() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3Caller) WITNESSTYPESTRING(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "WITNESS_TYPE_STRING")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// WITNESSTYPESTRING is a free data retrieval call binding the contract method 0x156e2152.
//
// Solidity: function WITNESS_TYPE_STRING() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3Session) WITNESSTYPESTRING() (string, error) {
	return _Fastsettlementv3.Contract.WITNESSTYPESTRING(&_Fastsettlementv3.CallOpts)
}

// WITNESSTYPESTRING is a free data retrieval call binding the contract method 0x156e2152.
//
// Solidity: function WITNESS_TYPE_STRING() view returns(string)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) WITNESSTYPESTRING() (string, error) {
	return _Fastsettlementv3.Contract.WITNESSTYPESTRING(&_Fastsettlementv3.CallOpts)
}

// AllowedSwapTargets is a free data retrieval call binding the contract method 0x1fa1fe36.
//
// Solidity: function allowedSwapTargets(address ) view returns(bool)
func (_Fastsettlementv3 *Fastsettlementv3Caller) AllowedSwapTargets(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "allowedSwapTargets", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedSwapTargets is a free data retrieval call binding the contract method 0x1fa1fe36.
//
// Solidity: function allowedSwapTargets(address ) view returns(bool)
func (_Fastsettlementv3 *Fastsettlementv3Session) AllowedSwapTargets(arg0 common.Address) (bool, error) {
	return _Fastsettlementv3.Contract.AllowedSwapTargets(&_Fastsettlementv3.CallOpts, arg0)
}

// AllowedSwapTargets is a free data retrieval call binding the contract method 0x1fa1fe36.
//
// Solidity: function allowedSwapTargets(address ) view returns(bool)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) AllowedSwapTargets(arg0 common.Address) (bool, error) {
	return _Fastsettlementv3.Contract.AllowedSwapTargets(&_Fastsettlementv3.CallOpts, arg0)
}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) Executor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "executor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) Executor() (common.Address, error) {
	return _Fastsettlementv3.Contract.Executor(&_Fastsettlementv3.CallOpts)
}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) Executor() (common.Address, error) {
	return _Fastsettlementv3.Contract.Executor(&_Fastsettlementv3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) Owner() (common.Address, error) {
	return _Fastsettlementv3.Contract.Owner(&_Fastsettlementv3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) Owner() (common.Address, error) {
	return _Fastsettlementv3.Contract.Owner(&_Fastsettlementv3.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) PendingOwner() (common.Address, error) {
	return _Fastsettlementv3.Contract.PendingOwner(&_Fastsettlementv3.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) PendingOwner() (common.Address, error) {
	return _Fastsettlementv3.Contract.PendingOwner(&_Fastsettlementv3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3Session) ProxiableUUID() ([32]byte, error) {
	return _Fastsettlementv3.Contract.ProxiableUUID(&_Fastsettlementv3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) ProxiableUUID() ([32]byte, error) {
	return _Fastsettlementv3.Contract.ProxiableUUID(&_Fastsettlementv3.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Caller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Fastsettlementv3.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3Session) Treasury() (common.Address, error) {
	return _Fastsettlementv3.Contract.Treasury(&_Fastsettlementv3.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_Fastsettlementv3 *Fastsettlementv3CallerSession) Treasury() (common.Address, error) {
	return _Fastsettlementv3.Contract.Treasury(&_Fastsettlementv3.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) AcceptOwnership() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.AcceptOwnership(&_Fastsettlementv3.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.AcceptOwnership(&_Fastsettlementv3.TransactOpts)
}

// ExecuteWithETH is a paid mutator transaction binding the contract method 0x1fb7a307.
//
// Solidity: function executeWithETH((address,address,address,uint256,uint256,address,uint256,uint256) intent, (address,uint256,bytes) swapData) payable returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Transactor) ExecuteWithETH(opts *bind.TransactOpts, intent IFastSettlementV3Intent, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "executeWithETH", intent, swapData)
}

// ExecuteWithETH is a paid mutator transaction binding the contract method 0x1fb7a307.
//
// Solidity: function executeWithETH((address,address,address,uint256,uint256,address,uint256,uint256) intent, (address,uint256,bytes) swapData) payable returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Session) ExecuteWithETH(intent IFastSettlementV3Intent, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.ExecuteWithETH(&_Fastsettlementv3.TransactOpts, intent, swapData)
}

// ExecuteWithETH is a paid mutator transaction binding the contract method 0x1fb7a307.
//
// Solidity: function executeWithETH((address,address,address,uint256,uint256,address,uint256,uint256) intent, (address,uint256,bytes) swapData) payable returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) ExecuteWithETH(intent IFastSettlementV3Intent, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.ExecuteWithETH(&_Fastsettlementv3.TransactOpts, intent, swapData)
}

// ExecuteWithPermit is a paid mutator transaction binding the contract method 0x02c52a55.
//
// Solidity: function executeWithPermit((address,address,address,uint256,uint256,address,uint256,uint256) intent, bytes signature, (address,uint256,bytes) swapData) returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Transactor) ExecuteWithPermit(opts *bind.TransactOpts, intent IFastSettlementV3Intent, signature []byte, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "executeWithPermit", intent, signature, swapData)
}

// ExecuteWithPermit is a paid mutator transaction binding the contract method 0x02c52a55.
//
// Solidity: function executeWithPermit((address,address,address,uint256,uint256,address,uint256,uint256) intent, bytes signature, (address,uint256,bytes) swapData) returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Session) ExecuteWithPermit(intent IFastSettlementV3Intent, signature []byte, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.ExecuteWithPermit(&_Fastsettlementv3.TransactOpts, intent, signature, swapData)
}

// ExecuteWithPermit is a paid mutator transaction binding the contract method 0x02c52a55.
//
// Solidity: function executeWithPermit((address,address,address,uint256,uint256,address,uint256,uint256) intent, bytes signature, (address,uint256,bytes) swapData) returns(uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) ExecuteWithPermit(intent IFastSettlementV3Intent, signature []byte, swapData IFastSettlementV3SwapCall) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.ExecuteWithPermit(&_Fastsettlementv3.TransactOpts, intent, signature, swapData)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address _owner, address _executor, address _treasury, address[] _initialSwapTargets) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _executor common.Address, _treasury common.Address, _initialSwapTargets []common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "initialize", _owner, _executor, _treasury, _initialSwapTargets)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address _owner, address _executor, address _treasury, address[] _initialSwapTargets) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) Initialize(_owner common.Address, _executor common.Address, _treasury common.Address, _initialSwapTargets []common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Initialize(&_Fastsettlementv3.TransactOpts, _owner, _executor, _treasury, _initialSwapTargets)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address _owner, address _executor, address _treasury, address[] _initialSwapTargets) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) Initialize(_owner common.Address, _executor common.Address, _treasury common.Address, _initialSwapTargets []common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Initialize(&_Fastsettlementv3.TransactOpts, _owner, _executor, _treasury, _initialSwapTargets)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) RenounceOwnership() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.RenounceOwnership(&_Fastsettlementv3.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.RenounceOwnership(&_Fastsettlementv3.TransactOpts)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x57376198.
//
// Solidity: function rescueTokens(address token, uint256 amount) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) RescueTokens(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "rescueTokens", token, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x57376198.
//
// Solidity: function rescueTokens(address token, uint256 amount) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) RescueTokens(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.RescueTokens(&_Fastsettlementv3.TransactOpts, token, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0x57376198.
//
// Solidity: function rescueTokens(address token, uint256 amount) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) RescueTokens(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.RescueTokens(&_Fastsettlementv3.TransactOpts, token, amount)
}

// SetExecutor is a paid mutator transaction binding the contract method 0x1c3c0ea8.
//
// Solidity: function setExecutor(address _newExecutor) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) SetExecutor(opts *bind.TransactOpts, _newExecutor common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "setExecutor", _newExecutor)
}

// SetExecutor is a paid mutator transaction binding the contract method 0x1c3c0ea8.
//
// Solidity: function setExecutor(address _newExecutor) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) SetExecutor(_newExecutor common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetExecutor(&_Fastsettlementv3.TransactOpts, _newExecutor)
}

// SetExecutor is a paid mutator transaction binding the contract method 0x1c3c0ea8.
//
// Solidity: function setExecutor(address _newExecutor) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) SetExecutor(_newExecutor common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetExecutor(&_Fastsettlementv3.TransactOpts, _newExecutor)
}

// SetSwapTargets is a paid mutator transaction binding the contract method 0x57d6924c.
//
// Solidity: function setSwapTargets(address[] targets, bool[] allowed) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) SetSwapTargets(opts *bind.TransactOpts, targets []common.Address, allowed []bool) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "setSwapTargets", targets, allowed)
}

// SetSwapTargets is a paid mutator transaction binding the contract method 0x57d6924c.
//
// Solidity: function setSwapTargets(address[] targets, bool[] allowed) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) SetSwapTargets(targets []common.Address, allowed []bool) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetSwapTargets(&_Fastsettlementv3.TransactOpts, targets, allowed)
}

// SetSwapTargets is a paid mutator transaction binding the contract method 0x57d6924c.
//
// Solidity: function setSwapTargets(address[] targets, bool[] allowed) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) SetSwapTargets(targets []common.Address, allowed []bool) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetSwapTargets(&_Fastsettlementv3.TransactOpts, targets, allowed)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address _newTreasury) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) SetTreasury(opts *bind.TransactOpts, _newTreasury common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "setTreasury", _newTreasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address _newTreasury) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) SetTreasury(_newTreasury common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetTreasury(&_Fastsettlementv3.TransactOpts, _newTreasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address _newTreasury) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) SetTreasury(_newTreasury common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.SetTreasury(&_Fastsettlementv3.TransactOpts, _newTreasury)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.TransferOwnership(&_Fastsettlementv3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.TransferOwnership(&_Fastsettlementv3.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.UpgradeToAndCall(&_Fastsettlementv3.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.UpgradeToAndCall(&_Fastsettlementv3.TransactOpts, newImplementation, data)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fastsettlementv3 *Fastsettlementv3Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Fastsettlementv3.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fastsettlementv3 *Fastsettlementv3Session) Receive() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Receive(&_Fastsettlementv3.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Fastsettlementv3 *Fastsettlementv3TransactorSession) Receive() (*types.Transaction, error) {
	return _Fastsettlementv3.Contract.Receive(&_Fastsettlementv3.TransactOpts)
}

// Fastsettlementv3ExecutorUpdatedIterator is returned from FilterExecutorUpdated and is used to iterate over the raw logs and unpacked data for ExecutorUpdated events raised by the Fastsettlementv3 contract.
type Fastsettlementv3ExecutorUpdatedIterator struct {
	Event *Fastsettlementv3ExecutorUpdated // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3ExecutorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3ExecutorUpdated)
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
		it.Event = new(Fastsettlementv3ExecutorUpdated)
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
func (it *Fastsettlementv3ExecutorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3ExecutorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3ExecutorUpdated represents a ExecutorUpdated event raised by the Fastsettlementv3 contract.
type Fastsettlementv3ExecutorUpdated struct {
	OldExecutor common.Address
	NewExecutor common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterExecutorUpdated is a free log retrieval operation binding the contract event 0x0ef3c7eb9dbcf33ddf032f4cce366a07eda85eed03e3172e4a90c4cc16d57886.
//
// Solidity: event ExecutorUpdated(address indexed oldExecutor, address indexed newExecutor)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterExecutorUpdated(opts *bind.FilterOpts, oldExecutor []common.Address, newExecutor []common.Address) (*Fastsettlementv3ExecutorUpdatedIterator, error) {

	var oldExecutorRule []interface{}
	for _, oldExecutorItem := range oldExecutor {
		oldExecutorRule = append(oldExecutorRule, oldExecutorItem)
	}
	var newExecutorRule []interface{}
	for _, newExecutorItem := range newExecutor {
		newExecutorRule = append(newExecutorRule, newExecutorItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "ExecutorUpdated", oldExecutorRule, newExecutorRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3ExecutorUpdatedIterator{contract: _Fastsettlementv3.contract, event: "ExecutorUpdated", logs: logs, sub: sub}, nil
}

// WatchExecutorUpdated is a free log subscription operation binding the contract event 0x0ef3c7eb9dbcf33ddf032f4cce366a07eda85eed03e3172e4a90c4cc16d57886.
//
// Solidity: event ExecutorUpdated(address indexed oldExecutor, address indexed newExecutor)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchExecutorUpdated(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3ExecutorUpdated, oldExecutor []common.Address, newExecutor []common.Address) (event.Subscription, error) {

	var oldExecutorRule []interface{}
	for _, oldExecutorItem := range oldExecutor {
		oldExecutorRule = append(oldExecutorRule, oldExecutorItem)
	}
	var newExecutorRule []interface{}
	for _, newExecutorItem := range newExecutor {
		newExecutorRule = append(newExecutorRule, newExecutorItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "ExecutorUpdated", oldExecutorRule, newExecutorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3ExecutorUpdated)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "ExecutorUpdated", log); err != nil {
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

// ParseExecutorUpdated is a log parse operation binding the contract event 0x0ef3c7eb9dbcf33ddf032f4cce366a07eda85eed03e3172e4a90c4cc16d57886.
//
// Solidity: event ExecutorUpdated(address indexed oldExecutor, address indexed newExecutor)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseExecutorUpdated(log types.Log) (*Fastsettlementv3ExecutorUpdated, error) {
	event := new(Fastsettlementv3ExecutorUpdated)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "ExecutorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Fastsettlementv3 contract.
type Fastsettlementv3InitializedIterator struct {
	Event *Fastsettlementv3Initialized // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3Initialized)
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
		it.Event = new(Fastsettlementv3Initialized)
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
func (it *Fastsettlementv3InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3Initialized represents a Initialized event raised by the Fastsettlementv3 contract.
type Fastsettlementv3Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterInitialized(opts *bind.FilterOpts) (*Fastsettlementv3InitializedIterator, error) {

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3InitializedIterator{contract: _Fastsettlementv3.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3Initialized) (event.Subscription, error) {

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3Initialized)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseInitialized(log types.Log) (*Fastsettlementv3Initialized, error) {
	event := new(Fastsettlementv3Initialized)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3IntentExecutedIterator is returned from FilterIntentExecuted and is used to iterate over the raw logs and unpacked data for IntentExecuted events raised by the Fastsettlementv3 contract.
type Fastsettlementv3IntentExecutedIterator struct {
	Event *Fastsettlementv3IntentExecuted // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3IntentExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3IntentExecuted)
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
		it.Event = new(Fastsettlementv3IntentExecuted)
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
func (it *Fastsettlementv3IntentExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3IntentExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3IntentExecuted represents a IntentExecuted event raised by the Fastsettlementv3 contract.
type Fastsettlementv3IntentExecuted struct {
	User        common.Address
	InputToken  common.Address
	OutputToken common.Address
	InputAmt    *big.Int
	UserAmtOut  *big.Int
	Received    *big.Int
	Surplus     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIntentExecuted is a free log retrieval operation binding the contract event 0x1ad6a4af59e844de3a921ec3dba60cb46f0b9051c9a106258624709dff629a87.
//
// Solidity: event IntentExecuted(address indexed user, address indexed inputToken, address indexed outputToken, uint256 inputAmt, uint256 userAmtOut, uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterIntentExecuted(opts *bind.FilterOpts, user []common.Address, inputToken []common.Address, outputToken []common.Address) (*Fastsettlementv3IntentExecutedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var inputTokenRule []interface{}
	for _, inputTokenItem := range inputToken {
		inputTokenRule = append(inputTokenRule, inputTokenItem)
	}
	var outputTokenRule []interface{}
	for _, outputTokenItem := range outputToken {
		outputTokenRule = append(outputTokenRule, outputTokenItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "IntentExecuted", userRule, inputTokenRule, outputTokenRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3IntentExecutedIterator{contract: _Fastsettlementv3.contract, event: "IntentExecuted", logs: logs, sub: sub}, nil
}

// WatchIntentExecuted is a free log subscription operation binding the contract event 0x1ad6a4af59e844de3a921ec3dba60cb46f0b9051c9a106258624709dff629a87.
//
// Solidity: event IntentExecuted(address indexed user, address indexed inputToken, address indexed outputToken, uint256 inputAmt, uint256 userAmtOut, uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchIntentExecuted(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3IntentExecuted, user []common.Address, inputToken []common.Address, outputToken []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var inputTokenRule []interface{}
	for _, inputTokenItem := range inputToken {
		inputTokenRule = append(inputTokenRule, inputTokenItem)
	}
	var outputTokenRule []interface{}
	for _, outputTokenItem := range outputToken {
		outputTokenRule = append(outputTokenRule, outputTokenItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "IntentExecuted", userRule, inputTokenRule, outputTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3IntentExecuted)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "IntentExecuted", log); err != nil {
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

// ParseIntentExecuted is a log parse operation binding the contract event 0x1ad6a4af59e844de3a921ec3dba60cb46f0b9051c9a106258624709dff629a87.
//
// Solidity: event IntentExecuted(address indexed user, address indexed inputToken, address indexed outputToken, uint256 inputAmt, uint256 userAmtOut, uint256 received, uint256 surplus)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseIntentExecuted(log types.Log) (*Fastsettlementv3IntentExecuted, error) {
	event := new(Fastsettlementv3IntentExecuted)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "IntentExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3OwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Fastsettlementv3 contract.
type Fastsettlementv3OwnershipTransferStartedIterator struct {
	Event *Fastsettlementv3OwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3OwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3OwnershipTransferStarted)
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
		it.Event = new(Fastsettlementv3OwnershipTransferStarted)
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
func (it *Fastsettlementv3OwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3OwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3OwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Fastsettlementv3 contract.
type Fastsettlementv3OwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Fastsettlementv3OwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3OwnershipTransferStartedIterator{contract: _Fastsettlementv3.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3OwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3OwnershipTransferStarted)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseOwnershipTransferStarted(log types.Log) (*Fastsettlementv3OwnershipTransferStarted, error) {
	event := new(Fastsettlementv3OwnershipTransferStarted)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Fastsettlementv3 contract.
type Fastsettlementv3OwnershipTransferredIterator struct {
	Event *Fastsettlementv3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3OwnershipTransferred)
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
		it.Event = new(Fastsettlementv3OwnershipTransferred)
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
func (it *Fastsettlementv3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3OwnershipTransferred represents a OwnershipTransferred event raised by the Fastsettlementv3 contract.
type Fastsettlementv3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Fastsettlementv3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3OwnershipTransferredIterator{contract: _Fastsettlementv3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3OwnershipTransferred)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseOwnershipTransferred(log types.Log) (*Fastsettlementv3OwnershipTransferred, error) {
	event := new(Fastsettlementv3OwnershipTransferred)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3SwapTargetsUpdatedIterator is returned from FilterSwapTargetsUpdated and is used to iterate over the raw logs and unpacked data for SwapTargetsUpdated events raised by the Fastsettlementv3 contract.
type Fastsettlementv3SwapTargetsUpdatedIterator struct {
	Event *Fastsettlementv3SwapTargetsUpdated // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3SwapTargetsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3SwapTargetsUpdated)
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
		it.Event = new(Fastsettlementv3SwapTargetsUpdated)
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
func (it *Fastsettlementv3SwapTargetsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3SwapTargetsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3SwapTargetsUpdated represents a SwapTargetsUpdated event raised by the Fastsettlementv3 contract.
type Fastsettlementv3SwapTargetsUpdated struct {
	Targets []common.Address
	Allowed []bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSwapTargetsUpdated is a free log retrieval operation binding the contract event 0xe18e0ae71e84871d203445f1d9d5c51bd93bb2e362ee0e455940a88475dc13bc.
//
// Solidity: event SwapTargetsUpdated(address[] targets, bool[] allowed)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterSwapTargetsUpdated(opts *bind.FilterOpts) (*Fastsettlementv3SwapTargetsUpdatedIterator, error) {

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "SwapTargetsUpdated")
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3SwapTargetsUpdatedIterator{contract: _Fastsettlementv3.contract, event: "SwapTargetsUpdated", logs: logs, sub: sub}, nil
}

// WatchSwapTargetsUpdated is a free log subscription operation binding the contract event 0xe18e0ae71e84871d203445f1d9d5c51bd93bb2e362ee0e455940a88475dc13bc.
//
// Solidity: event SwapTargetsUpdated(address[] targets, bool[] allowed)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchSwapTargetsUpdated(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3SwapTargetsUpdated) (event.Subscription, error) {

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "SwapTargetsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3SwapTargetsUpdated)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "SwapTargetsUpdated", log); err != nil {
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

// ParseSwapTargetsUpdated is a log parse operation binding the contract event 0xe18e0ae71e84871d203445f1d9d5c51bd93bb2e362ee0e455940a88475dc13bc.
//
// Solidity: event SwapTargetsUpdated(address[] targets, bool[] allowed)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseSwapTargetsUpdated(log types.Log) (*Fastsettlementv3SwapTargetsUpdated, error) {
	event := new(Fastsettlementv3SwapTargetsUpdated)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "SwapTargetsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3TreasuryUpdatedIterator is returned from FilterTreasuryUpdated and is used to iterate over the raw logs and unpacked data for TreasuryUpdated events raised by the Fastsettlementv3 contract.
type Fastsettlementv3TreasuryUpdatedIterator struct {
	Event *Fastsettlementv3TreasuryUpdated // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3TreasuryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3TreasuryUpdated)
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
		it.Event = new(Fastsettlementv3TreasuryUpdated)
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
func (it *Fastsettlementv3TreasuryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3TreasuryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3TreasuryUpdated represents a TreasuryUpdated event raised by the Fastsettlementv3 contract.
type Fastsettlementv3TreasuryUpdated struct {
	OldTreasury common.Address
	NewTreasury common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryUpdated is a free log retrieval operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterTreasuryUpdated(opts *bind.FilterOpts, oldTreasury []common.Address, newTreasury []common.Address) (*Fastsettlementv3TreasuryUpdatedIterator, error) {

	var oldTreasuryRule []interface{}
	for _, oldTreasuryItem := range oldTreasury {
		oldTreasuryRule = append(oldTreasuryRule, oldTreasuryItem)
	}
	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "TreasuryUpdated", oldTreasuryRule, newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3TreasuryUpdatedIterator{contract: _Fastsettlementv3.contract, event: "TreasuryUpdated", logs: logs, sub: sub}, nil
}

// WatchTreasuryUpdated is a free log subscription operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchTreasuryUpdated(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3TreasuryUpdated, oldTreasury []common.Address, newTreasury []common.Address) (event.Subscription, error) {

	var oldTreasuryRule []interface{}
	for _, oldTreasuryItem := range oldTreasury {
		oldTreasuryRule = append(oldTreasuryRule, oldTreasuryItem)
	}
	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "TreasuryUpdated", oldTreasuryRule, newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3TreasuryUpdated)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
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

// ParseTreasuryUpdated is a log parse operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseTreasuryUpdated(log types.Log) (*Fastsettlementv3TreasuryUpdated, error) {
	event := new(Fastsettlementv3TreasuryUpdated)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Fastsettlementv3UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Fastsettlementv3 contract.
type Fastsettlementv3UpgradedIterator struct {
	Event *Fastsettlementv3Upgraded // Event containing the contract specifics and raw log

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
func (it *Fastsettlementv3UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Fastsettlementv3Upgraded)
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
		it.Event = new(Fastsettlementv3Upgraded)
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
func (it *Fastsettlementv3UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Fastsettlementv3UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Fastsettlementv3Upgraded represents a Upgraded event raised by the Fastsettlementv3 contract.
type Fastsettlementv3Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*Fastsettlementv3UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &Fastsettlementv3UpgradedIterator{contract: _Fastsettlementv3.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Fastsettlementv3 *Fastsettlementv3Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *Fastsettlementv3Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Fastsettlementv3.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Fastsettlementv3Upgraded)
				if err := _Fastsettlementv3.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Fastsettlementv3 *Fastsettlementv3Filterer) ParseUpgraded(log types.Log) (*Fastsettlementv3Upgraded, error) {
	event := new(Fastsettlementv3Upgraded)
	if err := _Fastsettlementv3.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
