// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package depositmanager

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

// DepositmanagerMetaData contains all meta data concerning the Depositmanager contract.
var DepositmanagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_registry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minBalance\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"BIDDER_REGISTRY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_BALANCE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setTargetDeposit\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTargetDeposits\",\"inputs\":[{\"name\":\"providers\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"targetDeposits\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"topUpDeposit\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CurrentBalanceAtOrBelowMin\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currentBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"minBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CurrentDepositIsSufficient\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currentDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"targetDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DepositToppedUp\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TargetDepositDoesNotExist\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TargetDepositSet\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TopUpReduced\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"available\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalRequestExists\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotThisEOA\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"thisAddress\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// DepositmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositmanagerMetaData.ABI instead.
var DepositmanagerABI = DepositmanagerMetaData.ABI

// Depositmanager is an auto generated Go binding around an Ethereum contract.
type Depositmanager struct {
	DepositmanagerCaller     // Read-only binding to the contract
	DepositmanagerTransactor // Write-only binding to the contract
	DepositmanagerFilterer   // Log filterer for contract events
}

// DepositmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositmanagerSession struct {
	Contract     *Depositmanager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositmanagerCallerSession struct {
	Contract *DepositmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// DepositmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositmanagerTransactorSession struct {
	Contract     *DepositmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DepositmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositmanagerRaw struct {
	Contract *Depositmanager // Generic contract binding to access the raw methods on
}

// DepositmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositmanagerCallerRaw struct {
	Contract *DepositmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// DepositmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositmanagerTransactorRaw struct {
	Contract *DepositmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositmanager creates a new instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanager(address common.Address, backend bind.ContractBackend) (*Depositmanager, error) {
	contract, err := bindDepositmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depositmanager{DepositmanagerCaller: DepositmanagerCaller{contract: contract}, DepositmanagerTransactor: DepositmanagerTransactor{contract: contract}, DepositmanagerFilterer: DepositmanagerFilterer{contract: contract}}, nil
}

// NewDepositmanagerCaller creates a new read-only instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerCaller(address common.Address, caller bind.ContractCaller) (*DepositmanagerCaller, error) {
	contract, err := bindDepositmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerCaller{contract: contract}, nil
}

// NewDepositmanagerTransactor creates a new write-only instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositmanagerTransactor, error) {
	contract, err := bindDepositmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTransactor{contract: contract}, nil
}

// NewDepositmanagerFilterer creates a new log filterer instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositmanagerFilterer, error) {
	contract, err := bindDepositmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerFilterer{contract: contract}, nil
}

// bindDepositmanager binds a generic wrapper to an already deployed contract.
func bindDepositmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositmanager *DepositmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositmanager.Contract.DepositmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositmanager *DepositmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.Contract.DepositmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositmanager *DepositmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositmanager.Contract.DepositmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositmanager *DepositmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Depositmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositmanager *DepositmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositmanager *DepositmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositmanager.Contract.contract.Transact(opts, method, params...)
}

// BIDDERREGISTRY is a free data retrieval call binding the contract method 0xbf524608.
//
// Solidity: function BIDDER_REGISTRY() view returns(address)
func (_Depositmanager *DepositmanagerCaller) BIDDERREGISTRY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Depositmanager.contract.Call(opts, &out, "BIDDER_REGISTRY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BIDDERREGISTRY is a free data retrieval call binding the contract method 0xbf524608.
//
// Solidity: function BIDDER_REGISTRY() view returns(address)
func (_Depositmanager *DepositmanagerSession) BIDDERREGISTRY() (common.Address, error) {
	return _Depositmanager.Contract.BIDDERREGISTRY(&_Depositmanager.CallOpts)
}

// BIDDERREGISTRY is a free data retrieval call binding the contract method 0xbf524608.
//
// Solidity: function BIDDER_REGISTRY() view returns(address)
func (_Depositmanager *DepositmanagerCallerSession) BIDDERREGISTRY() (common.Address, error) {
	return _Depositmanager.Contract.BIDDERREGISTRY(&_Depositmanager.CallOpts)
}

// MINBALANCE is a free data retrieval call binding the contract method 0x867378c5.
//
// Solidity: function MIN_BALANCE() view returns(uint256)
func (_Depositmanager *DepositmanagerCaller) MINBALANCE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Depositmanager.contract.Call(opts, &out, "MIN_BALANCE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINBALANCE is a free data retrieval call binding the contract method 0x867378c5.
//
// Solidity: function MIN_BALANCE() view returns(uint256)
func (_Depositmanager *DepositmanagerSession) MINBALANCE() (*big.Int, error) {
	return _Depositmanager.Contract.MINBALANCE(&_Depositmanager.CallOpts)
}

// MINBALANCE is a free data retrieval call binding the contract method 0x867378c5.
//
// Solidity: function MIN_BALANCE() view returns(uint256)
func (_Depositmanager *DepositmanagerCallerSession) MINBALANCE() (*big.Int, error) {
	return _Depositmanager.Contract.MINBALANCE(&_Depositmanager.CallOpts)
}

// TargetDeposits is a free data retrieval call binding the contract method 0x77936281.
//
// Solidity: function targetDeposits(address ) view returns(uint256)
func (_Depositmanager *DepositmanagerCaller) TargetDeposits(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Depositmanager.contract.Call(opts, &out, "targetDeposits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetDeposits is a free data retrieval call binding the contract method 0x77936281.
//
// Solidity: function targetDeposits(address ) view returns(uint256)
func (_Depositmanager *DepositmanagerSession) TargetDeposits(arg0 common.Address) (*big.Int, error) {
	return _Depositmanager.Contract.TargetDeposits(&_Depositmanager.CallOpts, arg0)
}

// TargetDeposits is a free data retrieval call binding the contract method 0x77936281.
//
// Solidity: function targetDeposits(address ) view returns(uint256)
func (_Depositmanager *DepositmanagerCallerSession) TargetDeposits(arg0 common.Address) (*big.Int, error) {
	return _Depositmanager.Contract.TargetDeposits(&_Depositmanager.CallOpts, arg0)
}

// SetTargetDeposit is a paid mutator transaction binding the contract method 0x2d902b07.
//
// Solidity: function setTargetDeposit(address provider, uint256 amount) returns()
func (_Depositmanager *DepositmanagerTransactor) SetTargetDeposit(opts *bind.TransactOpts, provider common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "setTargetDeposit", provider, amount)
}

// SetTargetDeposit is a paid mutator transaction binding the contract method 0x2d902b07.
//
// Solidity: function setTargetDeposit(address provider, uint256 amount) returns()
func (_Depositmanager *DepositmanagerSession) SetTargetDeposit(provider common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetTargetDeposit(&_Depositmanager.TransactOpts, provider, amount)
}

// SetTargetDeposit is a paid mutator transaction binding the contract method 0x2d902b07.
//
// Solidity: function setTargetDeposit(address provider, uint256 amount) returns()
func (_Depositmanager *DepositmanagerTransactorSession) SetTargetDeposit(provider common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetTargetDeposit(&_Depositmanager.TransactOpts, provider, amount)
}

// SetTargetDeposits is a paid mutator transaction binding the contract method 0x5b79493e.
//
// Solidity: function setTargetDeposits(address[] providers, uint256[] amounts) returns()
func (_Depositmanager *DepositmanagerTransactor) SetTargetDeposits(opts *bind.TransactOpts, providers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "setTargetDeposits", providers, amounts)
}

// SetTargetDeposits is a paid mutator transaction binding the contract method 0x5b79493e.
//
// Solidity: function setTargetDeposits(address[] providers, uint256[] amounts) returns()
func (_Depositmanager *DepositmanagerSession) SetTargetDeposits(providers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetTargetDeposits(&_Depositmanager.TransactOpts, providers, amounts)
}

// SetTargetDeposits is a paid mutator transaction binding the contract method 0x5b79493e.
//
// Solidity: function setTargetDeposits(address[] providers, uint256[] amounts) returns()
func (_Depositmanager *DepositmanagerTransactorSession) SetTargetDeposits(providers []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetTargetDeposits(&_Depositmanager.TransactOpts, providers, amounts)
}

// TopUpDeposit is a paid mutator transaction binding the contract method 0xbe8c9cc0.
//
// Solidity: function topUpDeposit(address provider) returns()
func (_Depositmanager *DepositmanagerTransactor) TopUpDeposit(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "topUpDeposit", provider)
}

// TopUpDeposit is a paid mutator transaction binding the contract method 0xbe8c9cc0.
//
// Solidity: function topUpDeposit(address provider) returns()
func (_Depositmanager *DepositmanagerSession) TopUpDeposit(provider common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.TopUpDeposit(&_Depositmanager.TransactOpts, provider)
}

// TopUpDeposit is a paid mutator transaction binding the contract method 0xbe8c9cc0.
//
// Solidity: function topUpDeposit(address provider) returns()
func (_Depositmanager *DepositmanagerTransactorSession) TopUpDeposit(provider common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.TopUpDeposit(&_Depositmanager.TransactOpts, provider)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Depositmanager *DepositmanagerTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Depositmanager.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Depositmanager *DepositmanagerSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Depositmanager.Contract.Fallback(&_Depositmanager.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Depositmanager *DepositmanagerTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Depositmanager.Contract.Fallback(&_Depositmanager.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Depositmanager *DepositmanagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Depositmanager *DepositmanagerSession) Receive() (*types.Transaction, error) {
	return _Depositmanager.Contract.Receive(&_Depositmanager.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Depositmanager *DepositmanagerTransactorSession) Receive() (*types.Transaction, error) {
	return _Depositmanager.Contract.Receive(&_Depositmanager.TransactOpts)
}

// DepositmanagerCurrentBalanceAtOrBelowMinIterator is returned from FilterCurrentBalanceAtOrBelowMin and is used to iterate over the raw logs and unpacked data for CurrentBalanceAtOrBelowMin events raised by the Depositmanager contract.
type DepositmanagerCurrentBalanceAtOrBelowMinIterator struct {
	Event *DepositmanagerCurrentBalanceAtOrBelowMin // Event containing the contract specifics and raw log

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
func (it *DepositmanagerCurrentBalanceAtOrBelowMinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerCurrentBalanceAtOrBelowMin)
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
		it.Event = new(DepositmanagerCurrentBalanceAtOrBelowMin)
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
func (it *DepositmanagerCurrentBalanceAtOrBelowMinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerCurrentBalanceAtOrBelowMinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerCurrentBalanceAtOrBelowMin represents a CurrentBalanceAtOrBelowMin event raised by the Depositmanager contract.
type DepositmanagerCurrentBalanceAtOrBelowMin struct {
	Provider       common.Address
	CurrentBalance *big.Int
	MinBalance     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCurrentBalanceAtOrBelowMin is a free log retrieval operation binding the contract event 0x2816f654a045f3dc4fc70d7d509c97bd73066bbf09187a84950040cd3ba28079.
//
// Solidity: event CurrentBalanceAtOrBelowMin(address indexed provider, uint256 currentBalance, uint256 minBalance)
func (_Depositmanager *DepositmanagerFilterer) FilterCurrentBalanceAtOrBelowMin(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerCurrentBalanceAtOrBelowMinIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "CurrentBalanceAtOrBelowMin", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerCurrentBalanceAtOrBelowMinIterator{contract: _Depositmanager.contract, event: "CurrentBalanceAtOrBelowMin", logs: logs, sub: sub}, nil
}

// WatchCurrentBalanceAtOrBelowMin is a free log subscription operation binding the contract event 0x2816f654a045f3dc4fc70d7d509c97bd73066bbf09187a84950040cd3ba28079.
//
// Solidity: event CurrentBalanceAtOrBelowMin(address indexed provider, uint256 currentBalance, uint256 minBalance)
func (_Depositmanager *DepositmanagerFilterer) WatchCurrentBalanceAtOrBelowMin(opts *bind.WatchOpts, sink chan<- *DepositmanagerCurrentBalanceAtOrBelowMin, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "CurrentBalanceAtOrBelowMin", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerCurrentBalanceAtOrBelowMin)
				if err := _Depositmanager.contract.UnpackLog(event, "CurrentBalanceAtOrBelowMin", log); err != nil {
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

// ParseCurrentBalanceAtOrBelowMin is a log parse operation binding the contract event 0x2816f654a045f3dc4fc70d7d509c97bd73066bbf09187a84950040cd3ba28079.
//
// Solidity: event CurrentBalanceAtOrBelowMin(address indexed provider, uint256 currentBalance, uint256 minBalance)
func (_Depositmanager *DepositmanagerFilterer) ParseCurrentBalanceAtOrBelowMin(log types.Log) (*DepositmanagerCurrentBalanceAtOrBelowMin, error) {
	event := new(DepositmanagerCurrentBalanceAtOrBelowMin)
	if err := _Depositmanager.contract.UnpackLog(event, "CurrentBalanceAtOrBelowMin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerCurrentDepositIsSufficientIterator is returned from FilterCurrentDepositIsSufficient and is used to iterate over the raw logs and unpacked data for CurrentDepositIsSufficient events raised by the Depositmanager contract.
type DepositmanagerCurrentDepositIsSufficientIterator struct {
	Event *DepositmanagerCurrentDepositIsSufficient // Event containing the contract specifics and raw log

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
func (it *DepositmanagerCurrentDepositIsSufficientIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerCurrentDepositIsSufficient)
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
		it.Event = new(DepositmanagerCurrentDepositIsSufficient)
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
func (it *DepositmanagerCurrentDepositIsSufficientIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerCurrentDepositIsSufficientIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerCurrentDepositIsSufficient represents a CurrentDepositIsSufficient event raised by the Depositmanager contract.
type DepositmanagerCurrentDepositIsSufficient struct {
	Provider       common.Address
	CurrentDeposit *big.Int
	TargetDeposit  *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCurrentDepositIsSufficient is a free log retrieval operation binding the contract event 0xf53ce58639471dc76b2e6d62f3421857c3e3223a97849b7a22a6221f6423600b.
//
// Solidity: event CurrentDepositIsSufficient(address indexed provider, uint256 currentDeposit, uint256 targetDeposit)
func (_Depositmanager *DepositmanagerFilterer) FilterCurrentDepositIsSufficient(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerCurrentDepositIsSufficientIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "CurrentDepositIsSufficient", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerCurrentDepositIsSufficientIterator{contract: _Depositmanager.contract, event: "CurrentDepositIsSufficient", logs: logs, sub: sub}, nil
}

// WatchCurrentDepositIsSufficient is a free log subscription operation binding the contract event 0xf53ce58639471dc76b2e6d62f3421857c3e3223a97849b7a22a6221f6423600b.
//
// Solidity: event CurrentDepositIsSufficient(address indexed provider, uint256 currentDeposit, uint256 targetDeposit)
func (_Depositmanager *DepositmanagerFilterer) WatchCurrentDepositIsSufficient(opts *bind.WatchOpts, sink chan<- *DepositmanagerCurrentDepositIsSufficient, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "CurrentDepositIsSufficient", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerCurrentDepositIsSufficient)
				if err := _Depositmanager.contract.UnpackLog(event, "CurrentDepositIsSufficient", log); err != nil {
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

// ParseCurrentDepositIsSufficient is a log parse operation binding the contract event 0xf53ce58639471dc76b2e6d62f3421857c3e3223a97849b7a22a6221f6423600b.
//
// Solidity: event CurrentDepositIsSufficient(address indexed provider, uint256 currentDeposit, uint256 targetDeposit)
func (_Depositmanager *DepositmanagerFilterer) ParseCurrentDepositIsSufficient(log types.Log) (*DepositmanagerCurrentDepositIsSufficient, error) {
	event := new(DepositmanagerCurrentDepositIsSufficient)
	if err := _Depositmanager.contract.UnpackLog(event, "CurrentDepositIsSufficient", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerDepositToppedUpIterator is returned from FilterDepositToppedUp and is used to iterate over the raw logs and unpacked data for DepositToppedUp events raised by the Depositmanager contract.
type DepositmanagerDepositToppedUpIterator struct {
	Event *DepositmanagerDepositToppedUp // Event containing the contract specifics and raw log

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
func (it *DepositmanagerDepositToppedUpIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerDepositToppedUp)
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
		it.Event = new(DepositmanagerDepositToppedUp)
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
func (it *DepositmanagerDepositToppedUpIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerDepositToppedUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerDepositToppedUp represents a DepositToppedUp event raised by the Depositmanager contract.
type DepositmanagerDepositToppedUp struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositToppedUp is a free log retrieval operation binding the contract event 0xb5f2e468403466e947cb06c78c9b37008f6ff157cbf947e7cc8c675e128222e0.
//
// Solidity: event DepositToppedUp(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) FilterDepositToppedUp(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerDepositToppedUpIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "DepositToppedUp", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerDepositToppedUpIterator{contract: _Depositmanager.contract, event: "DepositToppedUp", logs: logs, sub: sub}, nil
}

// WatchDepositToppedUp is a free log subscription operation binding the contract event 0xb5f2e468403466e947cb06c78c9b37008f6ff157cbf947e7cc8c675e128222e0.
//
// Solidity: event DepositToppedUp(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) WatchDepositToppedUp(opts *bind.WatchOpts, sink chan<- *DepositmanagerDepositToppedUp, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "DepositToppedUp", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerDepositToppedUp)
				if err := _Depositmanager.contract.UnpackLog(event, "DepositToppedUp", log); err != nil {
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

// ParseDepositToppedUp is a log parse operation binding the contract event 0xb5f2e468403466e947cb06c78c9b37008f6ff157cbf947e7cc8c675e128222e0.
//
// Solidity: event DepositToppedUp(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) ParseDepositToppedUp(log types.Log) (*DepositmanagerDepositToppedUp, error) {
	event := new(DepositmanagerDepositToppedUp)
	if err := _Depositmanager.contract.UnpackLog(event, "DepositToppedUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerTargetDepositDoesNotExistIterator is returned from FilterTargetDepositDoesNotExist and is used to iterate over the raw logs and unpacked data for TargetDepositDoesNotExist events raised by the Depositmanager contract.
type DepositmanagerTargetDepositDoesNotExistIterator struct {
	Event *DepositmanagerTargetDepositDoesNotExist // Event containing the contract specifics and raw log

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
func (it *DepositmanagerTargetDepositDoesNotExistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerTargetDepositDoesNotExist)
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
		it.Event = new(DepositmanagerTargetDepositDoesNotExist)
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
func (it *DepositmanagerTargetDepositDoesNotExistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerTargetDepositDoesNotExistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerTargetDepositDoesNotExist represents a TargetDepositDoesNotExist event raised by the Depositmanager contract.
type DepositmanagerTargetDepositDoesNotExist struct {
	Provider common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTargetDepositDoesNotExist is a free log retrieval operation binding the contract event 0xd4a180e367ae6a1e28aa646fb6a2f3f20346488ec50bd03fb5986e3a9a1f9bfe.
//
// Solidity: event TargetDepositDoesNotExist(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) FilterTargetDepositDoesNotExist(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerTargetDepositDoesNotExistIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "TargetDepositDoesNotExist", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTargetDepositDoesNotExistIterator{contract: _Depositmanager.contract, event: "TargetDepositDoesNotExist", logs: logs, sub: sub}, nil
}

// WatchTargetDepositDoesNotExist is a free log subscription operation binding the contract event 0xd4a180e367ae6a1e28aa646fb6a2f3f20346488ec50bd03fb5986e3a9a1f9bfe.
//
// Solidity: event TargetDepositDoesNotExist(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) WatchTargetDepositDoesNotExist(opts *bind.WatchOpts, sink chan<- *DepositmanagerTargetDepositDoesNotExist, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "TargetDepositDoesNotExist", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerTargetDepositDoesNotExist)
				if err := _Depositmanager.contract.UnpackLog(event, "TargetDepositDoesNotExist", log); err != nil {
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

// ParseTargetDepositDoesNotExist is a log parse operation binding the contract event 0xd4a180e367ae6a1e28aa646fb6a2f3f20346488ec50bd03fb5986e3a9a1f9bfe.
//
// Solidity: event TargetDepositDoesNotExist(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) ParseTargetDepositDoesNotExist(log types.Log) (*DepositmanagerTargetDepositDoesNotExist, error) {
	event := new(DepositmanagerTargetDepositDoesNotExist)
	if err := _Depositmanager.contract.UnpackLog(event, "TargetDepositDoesNotExist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerTargetDepositSetIterator is returned from FilterTargetDepositSet and is used to iterate over the raw logs and unpacked data for TargetDepositSet events raised by the Depositmanager contract.
type DepositmanagerTargetDepositSetIterator struct {
	Event *DepositmanagerTargetDepositSet // Event containing the contract specifics and raw log

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
func (it *DepositmanagerTargetDepositSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerTargetDepositSet)
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
		it.Event = new(DepositmanagerTargetDepositSet)
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
func (it *DepositmanagerTargetDepositSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerTargetDepositSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerTargetDepositSet represents a TargetDepositSet event raised by the Depositmanager contract.
type DepositmanagerTargetDepositSet struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTargetDepositSet is a free log retrieval operation binding the contract event 0x6cb9e47917968530985729bd0855f1d40e32bfe78c5f738734c4f1f4e7e5e0ba.
//
// Solidity: event TargetDepositSet(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) FilterTargetDepositSet(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerTargetDepositSetIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "TargetDepositSet", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTargetDepositSetIterator{contract: _Depositmanager.contract, event: "TargetDepositSet", logs: logs, sub: sub}, nil
}

// WatchTargetDepositSet is a free log subscription operation binding the contract event 0x6cb9e47917968530985729bd0855f1d40e32bfe78c5f738734c4f1f4e7e5e0ba.
//
// Solidity: event TargetDepositSet(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) WatchTargetDepositSet(opts *bind.WatchOpts, sink chan<- *DepositmanagerTargetDepositSet, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "TargetDepositSet", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerTargetDepositSet)
				if err := _Depositmanager.contract.UnpackLog(event, "TargetDepositSet", log); err != nil {
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

// ParseTargetDepositSet is a log parse operation binding the contract event 0x6cb9e47917968530985729bd0855f1d40e32bfe78c5f738734c4f1f4e7e5e0ba.
//
// Solidity: event TargetDepositSet(address indexed provider, uint256 amount)
func (_Depositmanager *DepositmanagerFilterer) ParseTargetDepositSet(log types.Log) (*DepositmanagerTargetDepositSet, error) {
	event := new(DepositmanagerTargetDepositSet)
	if err := _Depositmanager.contract.UnpackLog(event, "TargetDepositSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerTopUpReducedIterator is returned from FilterTopUpReduced and is used to iterate over the raw logs and unpacked data for TopUpReduced events raised by the Depositmanager contract.
type DepositmanagerTopUpReducedIterator struct {
	Event *DepositmanagerTopUpReduced // Event containing the contract specifics and raw log

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
func (it *DepositmanagerTopUpReducedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerTopUpReduced)
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
		it.Event = new(DepositmanagerTopUpReduced)
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
func (it *DepositmanagerTopUpReducedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerTopUpReducedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerTopUpReduced represents a TopUpReduced event raised by the Depositmanager contract.
type DepositmanagerTopUpReduced struct {
	Provider  common.Address
	Available *big.Int
	Needed    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTopUpReduced is a free log retrieval operation binding the contract event 0x3d3e22cac7f9bc8da17ff9b43a14bbb7b8dd50ec92f62faebf162eb8d832bcc2.
//
// Solidity: event TopUpReduced(address indexed provider, uint256 available, uint256 needed)
func (_Depositmanager *DepositmanagerFilterer) FilterTopUpReduced(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerTopUpReducedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "TopUpReduced", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTopUpReducedIterator{contract: _Depositmanager.contract, event: "TopUpReduced", logs: logs, sub: sub}, nil
}

// WatchTopUpReduced is a free log subscription operation binding the contract event 0x3d3e22cac7f9bc8da17ff9b43a14bbb7b8dd50ec92f62faebf162eb8d832bcc2.
//
// Solidity: event TopUpReduced(address indexed provider, uint256 available, uint256 needed)
func (_Depositmanager *DepositmanagerFilterer) WatchTopUpReduced(opts *bind.WatchOpts, sink chan<- *DepositmanagerTopUpReduced, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "TopUpReduced", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerTopUpReduced)
				if err := _Depositmanager.contract.UnpackLog(event, "TopUpReduced", log); err != nil {
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

// ParseTopUpReduced is a log parse operation binding the contract event 0x3d3e22cac7f9bc8da17ff9b43a14bbb7b8dd50ec92f62faebf162eb8d832bcc2.
//
// Solidity: event TopUpReduced(address indexed provider, uint256 available, uint256 needed)
func (_Depositmanager *DepositmanagerFilterer) ParseTopUpReduced(log types.Log) (*DepositmanagerTopUpReduced, error) {
	event := new(DepositmanagerTopUpReduced)
	if err := _Depositmanager.contract.UnpackLog(event, "TopUpReduced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositmanagerWithdrawalRequestExistsIterator is returned from FilterWithdrawalRequestExists and is used to iterate over the raw logs and unpacked data for WithdrawalRequestExists events raised by the Depositmanager contract.
type DepositmanagerWithdrawalRequestExistsIterator struct {
	Event *DepositmanagerWithdrawalRequestExists // Event containing the contract specifics and raw log

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
func (it *DepositmanagerWithdrawalRequestExistsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerWithdrawalRequestExists)
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
		it.Event = new(DepositmanagerWithdrawalRequestExists)
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
func (it *DepositmanagerWithdrawalRequestExistsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerWithdrawalRequestExistsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerWithdrawalRequestExists represents a WithdrawalRequestExists event raised by the Depositmanager contract.
type DepositmanagerWithdrawalRequestExists struct {
	Provider common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRequestExists is a free log retrieval operation binding the contract event 0xd172492ccc8c62ae9847c209253d7ecc901a1e00e4752533ae796fe3a606b4c8.
//
// Solidity: event WithdrawalRequestExists(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) FilterWithdrawalRequestExists(opts *bind.FilterOpts, provider []common.Address) (*DepositmanagerWithdrawalRequestExistsIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "WithdrawalRequestExists", providerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerWithdrawalRequestExistsIterator{contract: _Depositmanager.contract, event: "WithdrawalRequestExists", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRequestExists is a free log subscription operation binding the contract event 0xd172492ccc8c62ae9847c209253d7ecc901a1e00e4752533ae796fe3a606b4c8.
//
// Solidity: event WithdrawalRequestExists(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) WatchWithdrawalRequestExists(opts *bind.WatchOpts, sink chan<- *DepositmanagerWithdrawalRequestExists, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "WithdrawalRequestExists", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerWithdrawalRequestExists)
				if err := _Depositmanager.contract.UnpackLog(event, "WithdrawalRequestExists", log); err != nil {
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

// ParseWithdrawalRequestExists is a log parse operation binding the contract event 0xd172492ccc8c62ae9847c209253d7ecc901a1e00e4752533ae796fe3a606b4c8.
//
// Solidity: event WithdrawalRequestExists(address indexed provider)
func (_Depositmanager *DepositmanagerFilterer) ParseWithdrawalRequestExists(log types.Log) (*DepositmanagerWithdrawalRequestExists, error) {
	event := new(DepositmanagerWithdrawalRequestExists)
	if err := _Depositmanager.contract.UnpackLog(event, "WithdrawalRequestExists", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
