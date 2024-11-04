// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatorregistryv1

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

// BlockHeightOccurrenceOccurrence is an auto generated low-level Go binding around an user-defined struct.
type BlockHeightOccurrenceOccurrence struct {
	Exists      bool
	BlockHeight *big.Int
}

// IVanillaRegistryStakedValidator is an auto generated low-level Go binding around an user-defined struct.
type IVanillaRegistryStakedValidator struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}

// Validatorregistryv1MetaData contains all meta data concerning the Validatorregistryv1 contract.
var Validatorregistryv1MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"claimForceWithdrawnFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"forceWithdrawalAsOwner\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forceWithdrawnFunds\",\"inputs\":[{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amountToClaim\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccumulatedSlashingFunds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlocksTillWithdrawAllowed\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedValidator\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVanillaRegistry.StakedValidator\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashOracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_slashReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_unstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSlashingPayoutDue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyTransferSlashingFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"newMinStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashOracle\",\"inputs\":[{\"name\":\"newSlashOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashReceiver\",\"inputs\":[{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashingPayoutPeriodBlocks\",\"inputs\":[{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnstakePeriodBlocks\",\"inputs\":[{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slash\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"payoutIfDue\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashingFundsTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accumulatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastPayoutBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakedValidators\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeeTransfer\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMinStake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashOracleSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashOracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashReceiverSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Slashed\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashReceiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashingPayoutPeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeAdded\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TotalStakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakePeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AtLeastOneRecipientRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeRecipientIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MinStakeMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustUnstakeToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoFundsToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotSlashOracle\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slashOracle\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotWithdrawalAddress\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlashAmountMustBeLessThanMinStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashAmountMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashOracleMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashReceiverMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingPayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakeTooLowForNumberOfKeys\",\"inputs\":[{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TransferToRecipientFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnstakePeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorCannotBeUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustNotExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMismatch\",\"inputs\":[{\"name\":\"actualWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawalFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawingTooSoon\",\"inputs\":[]}]",
}

// Validatorregistryv1ABI is the input ABI used to generate the binding from.
// Deprecated: Use Validatorregistryv1MetaData.ABI instead.
var Validatorregistryv1ABI = Validatorregistryv1MetaData.ABI

// Validatorregistryv1 is an auto generated Go binding around an Ethereum contract.
type Validatorregistryv1 struct {
	Validatorregistryv1Caller     // Read-only binding to the contract
	Validatorregistryv1Transactor // Write-only binding to the contract
	Validatorregistryv1Filterer   // Log filterer for contract events
}

// Validatorregistryv1Caller is an auto generated read-only Go binding around an Ethereum contract.
type Validatorregistryv1Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Validatorregistryv1Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Validatorregistryv1Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Validatorregistryv1Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Validatorregistryv1Session struct {
	Contract     *Validatorregistryv1 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// Validatorregistryv1CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Validatorregistryv1CallerSession struct {
	Contract *Validatorregistryv1Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// Validatorregistryv1TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Validatorregistryv1TransactorSession struct {
	Contract     *Validatorregistryv1Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// Validatorregistryv1Raw is an auto generated low-level Go binding around an Ethereum contract.
type Validatorregistryv1Raw struct {
	Contract *Validatorregistryv1 // Generic contract binding to access the raw methods on
}

// Validatorregistryv1CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Validatorregistryv1CallerRaw struct {
	Contract *Validatorregistryv1Caller // Generic read-only contract binding to access the raw methods on
}

// Validatorregistryv1TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Validatorregistryv1TransactorRaw struct {
	Contract *Validatorregistryv1Transactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorregistryv1 creates a new instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1(address common.Address, backend bind.ContractBackend) (*Validatorregistryv1, error) {
	contract, err := bindValidatorregistryv1(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1{Validatorregistryv1Caller: Validatorregistryv1Caller{contract: contract}, Validatorregistryv1Transactor: Validatorregistryv1Transactor{contract: contract}, Validatorregistryv1Filterer: Validatorregistryv1Filterer{contract: contract}}, nil
}

// NewValidatorregistryv1Caller creates a new read-only instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Caller(address common.Address, caller bind.ContractCaller) (*Validatorregistryv1Caller, error) {
	contract, err := bindValidatorregistryv1(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Caller{contract: contract}, nil
}

// NewValidatorregistryv1Transactor creates a new write-only instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Transactor(address common.Address, transactor bind.ContractTransactor) (*Validatorregistryv1Transactor, error) {
	contract, err := bindValidatorregistryv1(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Transactor{contract: contract}, nil
}

// NewValidatorregistryv1Filterer creates a new log filterer instance of Validatorregistryv1, bound to a specific deployed contract.
func NewValidatorregistryv1Filterer(address common.Address, filterer bind.ContractFilterer) (*Validatorregistryv1Filterer, error) {
	contract, err := bindValidatorregistryv1(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1Filterer{contract: contract}, nil
}

// bindValidatorregistryv1 binds a generic wrapper to an already deployed contract.
func bindValidatorregistryv1(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Validatorregistryv1MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistryv1.Contract.Validatorregistryv1Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Validatorregistryv1Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistryv1 *Validatorregistryv1Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Validatorregistryv1Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Validatorregistryv1 *Validatorregistryv1CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Validatorregistryv1.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Validatorregistryv1 *Validatorregistryv1TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Validatorregistryv1 *Validatorregistryv1TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatorregistryv1.Contract.UPGRADEINTERFACEVERSION(&_Validatorregistryv1.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Validatorregistryv1.Contract.UPGRADEINTERFACEVERSION(&_Validatorregistryv1.CallOpts)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Validatorregistryv1 *Validatorregistryv1Caller) ForceWithdrawnFunds(opts *bind.CallOpts, withdrawalAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "forceWithdrawnFunds", withdrawalAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Validatorregistryv1 *Validatorregistryv1Session) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Validatorregistryv1.Contract.ForceWithdrawnFunds(&_Validatorregistryv1.CallOpts, withdrawalAddress)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Validatorregistryv1.Contract.ForceWithdrawnFunds(&_Validatorregistryv1.CallOpts, withdrawalAddress)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetAccumulatedSlashingFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getAccumulatedSlashingFunds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetAccumulatedSlashingFunds(&_Validatorregistryv1.CallOpts)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetAccumulatedSlashingFunds(&_Validatorregistryv1.CallOpts)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetBlocksTillWithdrawAllowed(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getBlocksTillWithdrawAllowed", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetBlocksTillWithdrawAllowed(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetStakedAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getStakedAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Validatorregistryv1.Contract.GetStakedAmount(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Validatorregistryv1 *Validatorregistryv1Caller) GetStakedValidator(opts *bind.CallOpts, valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "getStakedValidator", valBLSPubKey)

	if err != nil {
		return *new(IVanillaRegistryStakedValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(IVanillaRegistryStakedValidator)).(*IVanillaRegistryStakedValidator)

	return out0, err

}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Validatorregistryv1 *Validatorregistryv1Session) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	return _Validatorregistryv1.Contract.GetStakedValidator(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	return _Validatorregistryv1.Contract.GetStakedValidator(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Caller) IsSlashingPayoutDue(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "isSlashingPayoutDue")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Session) IsSlashingPayoutDue() (bool, error) {
	return _Validatorregistryv1.Contract.IsSlashingPayoutDue(&_Validatorregistryv1.CallOpts)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) IsSlashingPayoutDue() (bool, error) {
	return _Validatorregistryv1.Contract.IsSlashingPayoutDue(&_Validatorregistryv1.CallOpts)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Caller) IsUnstaking(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "isUnstaking", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Session) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsUnstaking(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsUnstaking(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Caller) IsValidatorOptedIn(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "isValidatorOptedIn", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Session) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsValidatorOptedIn(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Validatorregistryv1.Contract.IsValidatorOptedIn(&_Validatorregistryv1.CallOpts, valBLSPubKey)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) MinStake() (*big.Int, error) {
	return _Validatorregistryv1.Contract.MinStake(&_Validatorregistryv1.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) MinStake() (*big.Int, error) {
	return _Validatorregistryv1.Contract.MinStake(&_Validatorregistryv1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Session) Owner() (common.Address, error) {
	return _Validatorregistryv1.Contract.Owner(&_Validatorregistryv1.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) Owner() (common.Address, error) {
	return _Validatorregistryv1.Contract.Owner(&_Validatorregistryv1.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1Session) Paused() (bool, error) {
	return _Validatorregistryv1.Contract.Paused(&_Validatorregistryv1.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) Paused() (bool, error) {
	return _Validatorregistryv1.Contract.Paused(&_Validatorregistryv1.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Caller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Session) PendingOwner() (common.Address, error) {
	return _Validatorregistryv1.Contract.PendingOwner(&_Validatorregistryv1.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) PendingOwner() (common.Address, error) {
	return _Validatorregistryv1.Contract.PendingOwner(&_Validatorregistryv1.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1Session) ProxiableUUID() ([32]byte, error) {
	return _Validatorregistryv1.Contract.ProxiableUUID(&_Validatorregistryv1.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) ProxiableUUID() ([32]byte, error) {
	return _Validatorregistryv1.Contract.ProxiableUUID(&_Validatorregistryv1.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Caller) SlashOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "slashOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1Session) SlashOracle() (common.Address, error) {
	return _Validatorregistryv1.Contract.SlashOracle(&_Validatorregistryv1.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) SlashOracle() (common.Address, error) {
	return _Validatorregistryv1.Contract.SlashOracle(&_Validatorregistryv1.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Caller) SlashingFundsTracker(opts *bind.CallOpts) (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "slashingFundsTracker")

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

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Session) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Validatorregistryv1.Contract.SlashingFundsTracker(&_Validatorregistryv1.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Validatorregistryv1.Contract.SlashingFundsTracker(&_Validatorregistryv1.CallOpts)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Validatorregistryv1 *Validatorregistryv1Caller) StakedValidators(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "stakedValidators", arg0)

	outstruct := new(struct {
		Exists            bool
		WithdrawalAddress common.Address
		Balance           *big.Int
		UnstakeOccurrence BlockHeightOccurrenceOccurrence
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.WithdrawalAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Balance = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.UnstakeOccurrence = *abi.ConvertType(out[3], new(BlockHeightOccurrenceOccurrence)).(*BlockHeightOccurrenceOccurrence)

	return *outstruct, err

}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Validatorregistryv1 *Validatorregistryv1Session) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Validatorregistryv1.Contract.StakedValidators(&_Validatorregistryv1.CallOpts, arg0)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Validatorregistryv1.Contract.StakedValidators(&_Validatorregistryv1.CallOpts, arg0)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Caller) UnstakePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Validatorregistryv1.contract.Call(opts, &out, "unstakePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1Session) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakePeriodBlocks(&_Validatorregistryv1.CallOpts)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Validatorregistryv1 *Validatorregistryv1CallerSession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Validatorregistryv1.Contract.UnstakePeriodBlocks(&_Validatorregistryv1.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) AcceptOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.AcceptOwnership(&_Validatorregistryv1.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.AcceptOwnership(&_Validatorregistryv1.TransactOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) AddStake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "addStake", blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.AddStake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.AddStake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) ClaimForceWithdrawnFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "claimForceWithdrawnFunds")
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ClaimForceWithdrawnFunds(&_Validatorregistryv1.TransactOpts)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ClaimForceWithdrawnFunds(&_Validatorregistryv1.TransactOpts)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) DelegateStake(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "delegateStake", blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.DelegateStake(&_Validatorregistryv1.TransactOpts, blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.DelegateStake(&_Validatorregistryv1.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) ForceWithdrawalAsOwner(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "forceWithdrawalAsOwner", blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ForceWithdrawalAsOwner(&_Validatorregistryv1.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ForceWithdrawalAsOwner(&_Validatorregistryv1.TransactOpts, blsPubKeys, withdrawalAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "initialize", _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Initialize(&_Validatorregistryv1.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Initialize(&_Validatorregistryv1.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) ManuallyTransferSlashingFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "manuallyTransferSlashingFunds")
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ManuallyTransferSlashingFunds(&_Validatorregistryv1.TransactOpts)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.ManuallyTransferSlashingFunds(&_Validatorregistryv1.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Pause() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Pause(&_Validatorregistryv1.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Pause() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Pause(&_Validatorregistryv1.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.RenounceOwnership(&_Validatorregistryv1.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.RenounceOwnership(&_Validatorregistryv1.TransactOpts)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) SetMinStake(opts *bind.TransactOpts, newMinStake *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "setMinStake", newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetMinStake(&_Validatorregistryv1.TransactOpts, newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetMinStake(&_Validatorregistryv1.TransactOpts, newMinStake)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) SetSlashOracle(opts *bind.TransactOpts, newSlashOracle common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "setSlashOracle", newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashOracle(&_Validatorregistryv1.TransactOpts, newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashOracle(&_Validatorregistryv1.TransactOpts, newSlashOracle)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) SetSlashReceiver(opts *bind.TransactOpts, newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "setSlashReceiver", newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashReceiver(&_Validatorregistryv1.TransactOpts, newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashReceiver(&_Validatorregistryv1.TransactOpts, newSlashReceiver)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) SetSlashingPayoutPeriodBlocks(opts *bind.TransactOpts, newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "setSlashingPayoutPeriodBlocks", newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashingPayoutPeriodBlocks(&_Validatorregistryv1.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetSlashingPayoutPeriodBlocks(&_Validatorregistryv1.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) SetUnstakePeriodBlocks(opts *bind.TransactOpts, newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "setUnstakePeriodBlocks", newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetUnstakePeriodBlocks(&_Validatorregistryv1.TransactOpts, newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.SetUnstakePeriodBlocks(&_Validatorregistryv1.TransactOpts, newUnstakePeriodBlocks)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Slash(opts *bind.TransactOpts, blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "slash", blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Slash(&_Validatorregistryv1.TransactOpts, blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Slash(&_Validatorregistryv1.TransactOpts, blsPubKeys, payoutIfDue)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Stake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "stake", blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Stake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Stake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.TransferOwnership(&_Validatorregistryv1.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.TransferOwnership(&_Validatorregistryv1.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Unpause() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unpause(&_Validatorregistryv1.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Unpause() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unpause(&_Validatorregistryv1.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Unstake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "unstake", blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unstake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Unstake(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.UpgradeToAndCall(&_Validatorregistryv1.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.UpgradeToAndCall(&_Validatorregistryv1.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Withdraw(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.Transact(opts, "withdraw", blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Withdraw(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Withdraw(&_Validatorregistryv1.TransactOpts, blsPubKeys)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Fallback(&_Validatorregistryv1.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Fallback(&_Validatorregistryv1.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Validatorregistryv1.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1Session) Receive() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Receive(&_Validatorregistryv1.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Validatorregistryv1 *Validatorregistryv1TransactorSession) Receive() (*types.Transaction, error) {
	return _Validatorregistryv1.Contract.Receive(&_Validatorregistryv1.TransactOpts)
}

// Validatorregistryv1FeeTransferIterator is returned from FilterFeeTransfer and is used to iterate over the raw logs and unpacked data for FeeTransfer events raised by the Validatorregistryv1 contract.
type Validatorregistryv1FeeTransferIterator struct {
	Event *Validatorregistryv1FeeTransfer // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1FeeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1FeeTransfer)
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
		it.Event = new(Validatorregistryv1FeeTransfer)
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
func (it *Validatorregistryv1FeeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1FeeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1FeeTransfer represents a FeeTransfer event raised by the Validatorregistryv1 contract.
type Validatorregistryv1FeeTransfer struct {
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeTransfer is a free log retrieval operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterFeeTransfer(opts *bind.FilterOpts, recipient []common.Address) (*Validatorregistryv1FeeTransferIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1FeeTransferIterator{contract: _Validatorregistryv1.contract, event: "FeeTransfer", logs: logs, sub: sub}, nil
}

// WatchFeeTransfer is a free log subscription operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchFeeTransfer(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1FeeTransfer, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1FeeTransfer)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseFeeTransfer(log types.Log) (*Validatorregistryv1FeeTransfer, error) {
	event := new(Validatorregistryv1FeeTransfer)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Validatorregistryv1 contract.
type Validatorregistryv1InitializedIterator struct {
	Event *Validatorregistryv1Initialized // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Initialized)
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
		it.Event = new(Validatorregistryv1Initialized)
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
func (it *Validatorregistryv1InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Initialized represents a Initialized event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterInitialized(opts *bind.FilterOpts) (*Validatorregistryv1InitializedIterator, error) {

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1InitializedIterator{contract: _Validatorregistryv1.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Initialized) (event.Subscription, error) {

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Initialized)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseInitialized(log types.Log) (*Validatorregistryv1Initialized, error) {
	event := new(Validatorregistryv1Initialized)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1MinStakeSetIterator is returned from FilterMinStakeSet and is used to iterate over the raw logs and unpacked data for MinStakeSet events raised by the Validatorregistryv1 contract.
type Validatorregistryv1MinStakeSetIterator struct {
	Event *Validatorregistryv1MinStakeSet // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1MinStakeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1MinStakeSet)
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
		it.Event = new(Validatorregistryv1MinStakeSet)
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
func (it *Validatorregistryv1MinStakeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1MinStakeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1MinStakeSet represents a MinStakeSet event raised by the Validatorregistryv1 contract.
type Validatorregistryv1MinStakeSet struct {
	MsgSender   common.Address
	NewMinStake *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMinStakeSet is a free log retrieval operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterMinStakeSet(opts *bind.FilterOpts, msgSender []common.Address) (*Validatorregistryv1MinStakeSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1MinStakeSetIterator{contract: _Validatorregistryv1.contract, event: "MinStakeSet", logs: logs, sub: sub}, nil
}

// WatchMinStakeSet is a free log subscription operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchMinStakeSet(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1MinStakeSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1MinStakeSet)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
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

// ParseMinStakeSet is a log parse operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseMinStakeSet(log types.Log) (*Validatorregistryv1MinStakeSet, error) {
	event := new(Validatorregistryv1MinStakeSet)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1OwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferStartedIterator struct {
	Event *Validatorregistryv1OwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1OwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1OwnershipTransferStarted)
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
		it.Event = new(Validatorregistryv1OwnershipTransferStarted)
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
func (it *Validatorregistryv1OwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1OwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1OwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Validatorregistryv1OwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1OwnershipTransferStartedIterator{contract: _Validatorregistryv1.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1OwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1OwnershipTransferStarted)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseOwnershipTransferStarted(log types.Log) (*Validatorregistryv1OwnershipTransferStarted, error) {
	event := new(Validatorregistryv1OwnershipTransferStarted)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferredIterator struct {
	Event *Validatorregistryv1OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1OwnershipTransferred)
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
		it.Event = new(Validatorregistryv1OwnershipTransferred)
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
func (it *Validatorregistryv1OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1OwnershipTransferred represents a OwnershipTransferred event raised by the Validatorregistryv1 contract.
type Validatorregistryv1OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Validatorregistryv1OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1OwnershipTransferredIterator{contract: _Validatorregistryv1.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1OwnershipTransferred)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseOwnershipTransferred(log types.Log) (*Validatorregistryv1OwnershipTransferred, error) {
	event := new(Validatorregistryv1OwnershipTransferred)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Validatorregistryv1 contract.
type Validatorregistryv1PausedIterator struct {
	Event *Validatorregistryv1Paused // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Paused)
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
		it.Event = new(Validatorregistryv1Paused)
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
func (it *Validatorregistryv1PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Paused represents a Paused event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterPaused(opts *bind.FilterOpts) (*Validatorregistryv1PausedIterator, error) {

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1PausedIterator{contract: _Validatorregistryv1.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Paused) (event.Subscription, error) {

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Paused)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParsePaused(log types.Log) (*Validatorregistryv1Paused, error) {
	event := new(Validatorregistryv1Paused)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1SlashOracleSetIterator is returned from FilterSlashOracleSet and is used to iterate over the raw logs and unpacked data for SlashOracleSet events raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashOracleSetIterator struct {
	Event *Validatorregistryv1SlashOracleSet // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1SlashOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1SlashOracleSet)
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
		it.Event = new(Validatorregistryv1SlashOracleSet)
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
func (it *Validatorregistryv1SlashOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1SlashOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1SlashOracleSet represents a SlashOracleSet event raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashOracleSet struct {
	MsgSender      common.Address
	NewSlashOracle common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSlashOracleSet is a free log retrieval operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterSlashOracleSet(opts *bind.FilterOpts, msgSender []common.Address) (*Validatorregistryv1SlashOracleSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1SlashOracleSetIterator{contract: _Validatorregistryv1.contract, event: "SlashOracleSet", logs: logs, sub: sub}, nil
}

// WatchSlashOracleSet is a free log subscription operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchSlashOracleSet(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1SlashOracleSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1SlashOracleSet)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
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

// ParseSlashOracleSet is a log parse operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseSlashOracleSet(log types.Log) (*Validatorregistryv1SlashOracleSet, error) {
	event := new(Validatorregistryv1SlashOracleSet)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1SlashReceiverSetIterator is returned from FilterSlashReceiverSet and is used to iterate over the raw logs and unpacked data for SlashReceiverSet events raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashReceiverSetIterator struct {
	Event *Validatorregistryv1SlashReceiverSet // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1SlashReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1SlashReceiverSet)
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
		it.Event = new(Validatorregistryv1SlashReceiverSet)
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
func (it *Validatorregistryv1SlashReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1SlashReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1SlashReceiverSet represents a SlashReceiverSet event raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashReceiverSet struct {
	MsgSender        common.Address
	NewSlashReceiver common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSlashReceiverSet is a free log retrieval operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterSlashReceiverSet(opts *bind.FilterOpts, msgSender []common.Address) (*Validatorregistryv1SlashReceiverSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1SlashReceiverSetIterator{contract: _Validatorregistryv1.contract, event: "SlashReceiverSet", logs: logs, sub: sub}, nil
}

// WatchSlashReceiverSet is a free log subscription operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchSlashReceiverSet(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1SlashReceiverSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1SlashReceiverSet)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
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

// ParseSlashReceiverSet is a log parse operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseSlashReceiverSet(log types.Log) (*Validatorregistryv1SlashReceiverSet, error) {
	event := new(Validatorregistryv1SlashReceiverSet)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1SlashedIterator is returned from FilterSlashed and is used to iterate over the raw logs and unpacked data for Slashed events raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashedIterator struct {
	Event *Validatorregistryv1Slashed // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1SlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Slashed)
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
		it.Event = new(Validatorregistryv1Slashed)
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
func (it *Validatorregistryv1SlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1SlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Slashed represents a Slashed event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Slashed struct {
	MsgSender         common.Address
	SlashReceiver     common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSlashed is a free log retrieval operation binding the contract event 0xf15b8630ce764d5dbcfaaa9843c3e5fcdb460aaaa46d7dc3ff4f19ca4096fc07.
//
// Solidity: event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterSlashed(opts *bind.FilterOpts, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1SlashedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var slashReceiverRule []interface{}
	for _, slashReceiverItem := range slashReceiver {
		slashReceiverRule = append(slashReceiverRule, slashReceiverItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1SlashedIterator{contract: _Validatorregistryv1.contract, event: "Slashed", logs: logs, sub: sub}, nil
}

// WatchSlashed is a free log subscription operation binding the contract event 0xf15b8630ce764d5dbcfaaa9843c3e5fcdb460aaaa46d7dc3ff4f19ca4096fc07.
//
// Solidity: event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchSlashed(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Slashed, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var slashReceiverRule []interface{}
	for _, slashReceiverItem := range slashReceiver {
		slashReceiverRule = append(slashReceiverRule, slashReceiverItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Slashed)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Slashed", log); err != nil {
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

// ParseSlashed is a log parse operation binding the contract event 0xf15b8630ce764d5dbcfaaa9843c3e5fcdb460aaaa46d7dc3ff4f19ca4096fc07.
//
// Solidity: event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseSlashed(log types.Log) (*Validatorregistryv1Slashed, error) {
	event := new(Validatorregistryv1Slashed)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Slashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator is returned from FilterSlashingPayoutPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for SlashingPayoutPeriodBlocksSet events raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator struct {
	Event *Validatorregistryv1SlashingPayoutPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1SlashingPayoutPeriodBlocksSet)
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
		it.Event = new(Validatorregistryv1SlashingPayoutPeriodBlocksSet)
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
func (it *Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1SlashingPayoutPeriodBlocksSet represents a SlashingPayoutPeriodBlocksSet event raised by the Validatorregistryv1 contract.
type Validatorregistryv1SlashingPayoutPeriodBlocksSet struct {
	MsgSender                     common.Address
	NewSlashingPayoutPeriodBlocks *big.Int
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterSlashingPayoutPeriodBlocksSet is a free log retrieval operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterSlashingPayoutPeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1SlashingPayoutPeriodBlocksSetIterator{contract: _Validatorregistryv1.contract, event: "SlashingPayoutPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchSlashingPayoutPeriodBlocksSet is a free log subscription operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchSlashingPayoutPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1SlashingPayoutPeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1SlashingPayoutPeriodBlocksSet)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
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

// ParseSlashingPayoutPeriodBlocksSet is a log parse operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseSlashingPayoutPeriodBlocksSet(log types.Log) (*Validatorregistryv1SlashingPayoutPeriodBlocksSet, error) {
	event := new(Validatorregistryv1SlashingPayoutPeriodBlocksSet)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1StakeAddedIterator is returned from FilterStakeAdded and is used to iterate over the raw logs and unpacked data for StakeAdded events raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeAddedIterator struct {
	Event *Validatorregistryv1StakeAdded // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1StakeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1StakeAdded)
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
		it.Event = new(Validatorregistryv1StakeAdded)
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
func (it *Validatorregistryv1StakeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1StakeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1StakeAdded represents a StakeAdded event raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeAdded struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	NewBalance        *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStakeAdded is a free log retrieval operation binding the contract event 0xb01516cc7ddda8b10127c714474503b38a75b9afa8a4e4b9da306e61181980c7.
//
// Solidity: event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterStakeAdded(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1StakeAddedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1StakeAddedIterator{contract: _Validatorregistryv1.contract, event: "StakeAdded", logs: logs, sub: sub}, nil
}

// WatchStakeAdded is a free log subscription operation binding the contract event 0xb01516cc7ddda8b10127c714474503b38a75b9afa8a4e4b9da306e61181980c7.
//
// Solidity: event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchStakeAdded(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1StakeAdded, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1StakeAdded)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeAdded", log); err != nil {
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

// ParseStakeAdded is a log parse operation binding the contract event 0xb01516cc7ddda8b10127c714474503b38a75b9afa8a4e4b9da306e61181980c7.
//
// Solidity: event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseStakeAdded(log types.Log) (*Validatorregistryv1StakeAdded, error) {
	event := new(Validatorregistryv1StakeAdded)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1StakeWithdrawnIterator is returned from FilterStakeWithdrawn and is used to iterate over the raw logs and unpacked data for StakeWithdrawn events raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeWithdrawnIterator struct {
	Event *Validatorregistryv1StakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1StakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1StakeWithdrawn)
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
		it.Event = new(Validatorregistryv1StakeWithdrawn)
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
func (it *Validatorregistryv1StakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1StakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1StakeWithdrawn represents a StakeWithdrawn event raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdrawn is a free log retrieval operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1StakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1StakeWithdrawnIterator{contract: _Validatorregistryv1.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchStakeWithdrawn is a free log subscription operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1StakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1StakeWithdrawn)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
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

// ParseStakeWithdrawn is a log parse operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseStakeWithdrawn(log types.Log) (*Validatorregistryv1StakeWithdrawn, error) {
	event := new(Validatorregistryv1StakeWithdrawn)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1StakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Validatorregistryv1 contract.
type Validatorregistryv1StakedIterator struct {
	Event *Validatorregistryv1Staked // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1StakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Staked)
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
		it.Event = new(Validatorregistryv1Staked)
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
func (it *Validatorregistryv1StakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1StakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Staked represents a Staked event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Staked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterStaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1StakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1StakedIterator{contract: _Validatorregistryv1.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Staked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Staked)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseStaked(log types.Log) (*Validatorregistryv1Staked, error) {
	event := new(Validatorregistryv1Staked)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1TotalStakeWithdrawnIterator is returned from FilterTotalStakeWithdrawn and is used to iterate over the raw logs and unpacked data for TotalStakeWithdrawn events raised by the Validatorregistryv1 contract.
type Validatorregistryv1TotalStakeWithdrawnIterator struct {
	Event *Validatorregistryv1TotalStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1TotalStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1TotalStakeWithdrawn)
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
		it.Event = new(Validatorregistryv1TotalStakeWithdrawn)
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
func (it *Validatorregistryv1TotalStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1TotalStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1TotalStakeWithdrawn represents a TotalStakeWithdrawn event raised by the Validatorregistryv1 contract.
type Validatorregistryv1TotalStakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	TotalAmount       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTotalStakeWithdrawn is a free log retrieval operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterTotalStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1TotalStakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1TotalStakeWithdrawnIterator{contract: _Validatorregistryv1.contract, event: "TotalStakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTotalStakeWithdrawn is a free log subscription operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchTotalStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1TotalStakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1TotalStakeWithdrawn)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
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

// ParseTotalStakeWithdrawn is a log parse operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseTotalStakeWithdrawn(log types.Log) (*Validatorregistryv1TotalStakeWithdrawn, error) {
	event := new(Validatorregistryv1TotalStakeWithdrawn)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UnpausedIterator struct {
	Event *Validatorregistryv1Unpaused // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Unpaused)
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
		it.Event = new(Validatorregistryv1Unpaused)
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
func (it *Validatorregistryv1UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Unpaused represents a Unpaused event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUnpaused(opts *bind.FilterOpts) (*Validatorregistryv1UnpausedIterator, error) {

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UnpausedIterator{contract: _Validatorregistryv1.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Unpaused) (event.Subscription, error) {

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Unpaused)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUnpaused(log types.Log) (*Validatorregistryv1Unpaused, error) {
	event := new(Validatorregistryv1Unpaused)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UnstakePeriodBlocksSetIterator is returned from FilterUnstakePeriodBlocksSet and is used to iterate over the raw logs and unpacked data for UnstakePeriodBlocksSet events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UnstakePeriodBlocksSetIterator struct {
	Event *Validatorregistryv1UnstakePeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UnstakePeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1UnstakePeriodBlocksSet)
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
		it.Event = new(Validatorregistryv1UnstakePeriodBlocksSet)
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
func (it *Validatorregistryv1UnstakePeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UnstakePeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1UnstakePeriodBlocksSet represents a UnstakePeriodBlocksSet event raised by the Validatorregistryv1 contract.
type Validatorregistryv1UnstakePeriodBlocksSet struct {
	MsgSender              common.Address
	NewUnstakePeriodBlocks *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUnstakePeriodBlocksSet is a free log retrieval operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUnstakePeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*Validatorregistryv1UnstakePeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UnstakePeriodBlocksSetIterator{contract: _Validatorregistryv1.contract, event: "UnstakePeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchUnstakePeriodBlocksSet is a free log subscription operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUnstakePeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1UnstakePeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1UnstakePeriodBlocksSet)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
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

// ParseUnstakePeriodBlocksSet is a log parse operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUnstakePeriodBlocksSet(log types.Log) (*Validatorregistryv1UnstakePeriodBlocksSet, error) {
	event := new(Validatorregistryv1UnstakePeriodBlocksSet)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UnstakedIterator struct {
	Event *Validatorregistryv1Unstaked // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Unstaked)
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
		it.Event = new(Validatorregistryv1Unstaked)
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
func (it *Validatorregistryv1UnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Unstaked represents a Unstaked event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Unstaked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUnstaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Validatorregistryv1UnstakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UnstakedIterator{contract: _Validatorregistryv1.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Unstaked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Unstaked)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUnstaked(log types.Log) (*Validatorregistryv1Unstaked, error) {
	event := new(Validatorregistryv1Unstaked)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Validatorregistryv1UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Validatorregistryv1 contract.
type Validatorregistryv1UpgradedIterator struct {
	Event *Validatorregistryv1Upgraded // Event containing the contract specifics and raw log

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
func (it *Validatorregistryv1UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Validatorregistryv1Upgraded)
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
		it.Event = new(Validatorregistryv1Upgraded)
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
func (it *Validatorregistryv1UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Validatorregistryv1UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Validatorregistryv1Upgraded represents a Upgraded event raised by the Validatorregistryv1 contract.
type Validatorregistryv1Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*Validatorregistryv1UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &Validatorregistryv1UpgradedIterator{contract: _Validatorregistryv1.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Validatorregistryv1 *Validatorregistryv1Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *Validatorregistryv1Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Validatorregistryv1.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Validatorregistryv1Upgraded)
				if err := _Validatorregistryv1.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Validatorregistryv1 *Validatorregistryv1Filterer) ParseUpgraded(log types.Log) (*Validatorregistryv1Upgraded, error) {
	event := new(Validatorregistryv1Upgraded)
	if err := _Validatorregistryv1.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
