// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vanillaregistry

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

// VanillaregistryMetaData contains all meta data concerning the Vanillaregistry contract.
var VanillaregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"claimForceWithdrawnFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"forceWithdrawalAsOwner\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forceWithdrawnFunds\",\"inputs\":[{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amountToClaim\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccumulatedSlashingFunds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlocksTillWithdrawAllowed\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedValidator\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVanillaRegistry.StakedValidator\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashOracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_slashReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_unstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSlashingPayoutDue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyTransferSlashingFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeWhitelistedStakers\",\"inputs\":[{\"name\":\"stakers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"newMinStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashOracle\",\"inputs\":[{\"name\":\"newSlashOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashReceiver\",\"inputs\":[{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashingPayoutPeriodBlocks\",\"inputs\":[{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnstakePeriodBlocks\",\"inputs\":[{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slash\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"payoutIfDue\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashingFundsTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accumulatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastPayoutBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakedValidators\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"whitelistStakers\",\"inputs\":[{\"name\":\"stakers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"whitelistedStakers\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"whitelisted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeeTransfer\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMinStake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashOracleSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashOracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashReceiverSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Slashed\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashReceiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashingPayoutPeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeAdded\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerRemovedFromWhitelist\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerWhitelisted\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TotalStakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakePeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AtLeastOneRecipientRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeRecipientIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MinStakeMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustUnstakeToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoFundsToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotSlashOracle\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slashOracle\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotWhitelistedStaker\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotWithdrawalAddress\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlashAmountMustBeLessThanMinStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashAmountMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashOracleMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashReceiverMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingPayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakeTooLowForNumberOfKeys\",\"inputs\":[{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StakerAlreadyWhitelisted\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StakerNotWhitelisted\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferToRecipientFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnstakePeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorCannotBeUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustNotExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMismatch\",\"inputs\":[{\"name\":\"actualWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawalFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawingTooSoon\",\"inputs\":[]}]",
}

// VanillaregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use VanillaregistryMetaData.ABI instead.
var VanillaregistryABI = VanillaregistryMetaData.ABI

// Vanillaregistry is an auto generated Go binding around an Ethereum contract.
type Vanillaregistry struct {
	VanillaregistryCaller     // Read-only binding to the contract
	VanillaregistryTransactor // Write-only binding to the contract
	VanillaregistryFilterer   // Log filterer for contract events
}

// VanillaregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type VanillaregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VanillaregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VanillaregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VanillaregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VanillaregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VanillaregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VanillaregistrySession struct {
	Contract     *Vanillaregistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VanillaregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VanillaregistryCallerSession struct {
	Contract *VanillaregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// VanillaregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VanillaregistryTransactorSession struct {
	Contract     *VanillaregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// VanillaregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type VanillaregistryRaw struct {
	Contract *Vanillaregistry // Generic contract binding to access the raw methods on
}

// VanillaregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VanillaregistryCallerRaw struct {
	Contract *VanillaregistryCaller // Generic read-only contract binding to access the raw methods on
}

// VanillaregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VanillaregistryTransactorRaw struct {
	Contract *VanillaregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVanillaregistry creates a new instance of Vanillaregistry, bound to a specific deployed contract.
func NewVanillaregistry(address common.Address, backend bind.ContractBackend) (*Vanillaregistry, error) {
	contract, err := bindVanillaregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistry{VanillaregistryCaller: VanillaregistryCaller{contract: contract}, VanillaregistryTransactor: VanillaregistryTransactor{contract: contract}, VanillaregistryFilterer: VanillaregistryFilterer{contract: contract}}, nil
}

// NewVanillaregistryCaller creates a new read-only instance of Vanillaregistry, bound to a specific deployed contract.
func NewVanillaregistryCaller(address common.Address, caller bind.ContractCaller) (*VanillaregistryCaller, error) {
	contract, err := bindVanillaregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryCaller{contract: contract}, nil
}

// NewVanillaregistryTransactor creates a new write-only instance of Vanillaregistry, bound to a specific deployed contract.
func NewVanillaregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*VanillaregistryTransactor, error) {
	contract, err := bindVanillaregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryTransactor{contract: contract}, nil
}

// NewVanillaregistryFilterer creates a new log filterer instance of Vanillaregistry, bound to a specific deployed contract.
func NewVanillaregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*VanillaregistryFilterer, error) {
	contract, err := bindVanillaregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryFilterer{contract: contract}, nil
}

// bindVanillaregistry binds a generic wrapper to an already deployed contract.
func bindVanillaregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VanillaregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vanillaregistry *VanillaregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vanillaregistry.Contract.VanillaregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vanillaregistry *VanillaregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.VanillaregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vanillaregistry *VanillaregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.VanillaregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vanillaregistry *VanillaregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vanillaregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vanillaregistry *VanillaregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vanillaregistry *VanillaregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistry *VanillaregistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistry *VanillaregistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Vanillaregistry.Contract.UPGRADEINTERFACEVERSION(&_Vanillaregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistry *VanillaregistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Vanillaregistry.Contract.UPGRADEINTERFACEVERSION(&_Vanillaregistry.CallOpts)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistry *VanillaregistryCaller) ForceWithdrawnFunds(opts *bind.CallOpts, withdrawalAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "forceWithdrawnFunds", withdrawalAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistry *VanillaregistrySession) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Vanillaregistry.Contract.ForceWithdrawnFunds(&_Vanillaregistry.CallOpts, withdrawalAddress)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistry *VanillaregistryCallerSession) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Vanillaregistry.Contract.ForceWithdrawnFunds(&_Vanillaregistry.CallOpts, withdrawalAddress)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCaller) GetAccumulatedSlashingFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "getAccumulatedSlashingFunds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistry *VanillaregistrySession) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Vanillaregistry.Contract.GetAccumulatedSlashingFunds(&_Vanillaregistry.CallOpts)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCallerSession) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Vanillaregistry.Contract.GetAccumulatedSlashingFunds(&_Vanillaregistry.CallOpts)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistryCaller) GetBlocksTillWithdrawAllowed(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "getBlocksTillWithdrawAllowed", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistrySession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistry.Contract.GetBlocksTillWithdrawAllowed(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistryCallerSession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistry.Contract.GetBlocksTillWithdrawAllowed(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistryCaller) GetStakedAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "getStakedAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistrySession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistry.Contract.GetStakedAmount(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistry *VanillaregistryCallerSession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistry.Contract.GetStakedAmount(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistry *VanillaregistryCaller) GetStakedValidator(opts *bind.CallOpts, valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "getStakedValidator", valBLSPubKey)

	if err != nil {
		return *new(IVanillaRegistryStakedValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(IVanillaRegistryStakedValidator)).(*IVanillaRegistryStakedValidator)

	return out0, err

}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistry *VanillaregistrySession) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	return _Vanillaregistry.Contract.GetStakedValidator(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistry *VanillaregistryCallerSession) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryStakedValidator, error) {
	return _Vanillaregistry.Contract.GetStakedValidator(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistry *VanillaregistryCaller) IsSlashingPayoutDue(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "isSlashingPayoutDue")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistry *VanillaregistrySession) IsSlashingPayoutDue() (bool, error) {
	return _Vanillaregistry.Contract.IsSlashingPayoutDue(&_Vanillaregistry.CallOpts)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistry *VanillaregistryCallerSession) IsSlashingPayoutDue() (bool, error) {
	return _Vanillaregistry.Contract.IsSlashingPayoutDue(&_Vanillaregistry.CallOpts)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistryCaller) IsUnstaking(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "isUnstaking", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistrySession) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistry.Contract.IsUnstaking(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistryCallerSession) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistry.Contract.IsUnstaking(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistryCaller) IsValidatorOptedIn(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "isValidatorOptedIn", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistrySession) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistry.Contract.IsValidatorOptedIn(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistry *VanillaregistryCallerSession) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistry.Contract.IsValidatorOptedIn(&_Vanillaregistry.CallOpts, valBLSPubKey)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistry *VanillaregistrySession) MinStake() (*big.Int, error) {
	return _Vanillaregistry.Contract.MinStake(&_Vanillaregistry.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCallerSession) MinStake() (*big.Int, error) {
	return _Vanillaregistry.Contract.MinStake(&_Vanillaregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistry *VanillaregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistry *VanillaregistrySession) Owner() (common.Address, error) {
	return _Vanillaregistry.Contract.Owner(&_Vanillaregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistry *VanillaregistryCallerSession) Owner() (common.Address, error) {
	return _Vanillaregistry.Contract.Owner(&_Vanillaregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistry *VanillaregistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistry *VanillaregistrySession) Paused() (bool, error) {
	return _Vanillaregistry.Contract.Paused(&_Vanillaregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistry *VanillaregistryCallerSession) Paused() (bool, error) {
	return _Vanillaregistry.Contract.Paused(&_Vanillaregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistry *VanillaregistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistry *VanillaregistrySession) PendingOwner() (common.Address, error) {
	return _Vanillaregistry.Contract.PendingOwner(&_Vanillaregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistry *VanillaregistryCallerSession) PendingOwner() (common.Address, error) {
	return _Vanillaregistry.Contract.PendingOwner(&_Vanillaregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistry *VanillaregistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistry *VanillaregistrySession) ProxiableUUID() ([32]byte, error) {
	return _Vanillaregistry.Contract.ProxiableUUID(&_Vanillaregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistry *VanillaregistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Vanillaregistry.Contract.ProxiableUUID(&_Vanillaregistry.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistry *VanillaregistryCaller) SlashOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "slashOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistry *VanillaregistrySession) SlashOracle() (common.Address, error) {
	return _Vanillaregistry.Contract.SlashOracle(&_Vanillaregistry.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistry *VanillaregistryCallerSession) SlashOracle() (common.Address, error) {
	return _Vanillaregistry.Contract.SlashOracle(&_Vanillaregistry.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Vanillaregistry *VanillaregistryCaller) SlashingFundsTracker(opts *bind.CallOpts) (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "slashingFundsTracker")

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
func (_Vanillaregistry *VanillaregistrySession) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Vanillaregistry.Contract.SlashingFundsTracker(&_Vanillaregistry.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Vanillaregistry *VanillaregistryCallerSession) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Vanillaregistry.Contract.SlashingFundsTracker(&_Vanillaregistry.CallOpts)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Vanillaregistry *VanillaregistryCaller) StakedValidators(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "stakedValidators", arg0)

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
func (_Vanillaregistry *VanillaregistrySession) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Vanillaregistry.Contract.StakedValidators(&_Vanillaregistry.CallOpts, arg0)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Vanillaregistry *VanillaregistryCallerSession) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Vanillaregistry.Contract.StakedValidators(&_Vanillaregistry.CallOpts, arg0)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCaller) UnstakePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "unstakePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistry *VanillaregistrySession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Vanillaregistry.Contract.UnstakePeriodBlocks(&_Vanillaregistry.CallOpts)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistry *VanillaregistryCallerSession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Vanillaregistry.Contract.UnstakePeriodBlocks(&_Vanillaregistry.CallOpts)
}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistry *VanillaregistryCaller) WhitelistedStakers(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _Vanillaregistry.contract.Call(opts, &out, "whitelistedStakers", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistry *VanillaregistrySession) WhitelistedStakers(staker common.Address) (bool, error) {
	return _Vanillaregistry.Contract.WhitelistedStakers(&_Vanillaregistry.CallOpts, staker)
}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistry *VanillaregistryCallerSession) WhitelistedStakers(staker common.Address) (bool, error) {
	return _Vanillaregistry.Contract.WhitelistedStakers(&_Vanillaregistry.CallOpts, staker)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistry *VanillaregistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistry *VanillaregistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.AcceptOwnership(&_Vanillaregistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.AcceptOwnership(&_Vanillaregistry.TransactOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) AddStake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "addStake", blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistrySession) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.AddStake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.AddStake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistry *VanillaregistryTransactor) ClaimForceWithdrawnFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "claimForceWithdrawnFunds")
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistry *VanillaregistrySession) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ClaimForceWithdrawnFunds(&_Vanillaregistry.TransactOpts)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ClaimForceWithdrawnFunds(&_Vanillaregistry.TransactOpts)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) DelegateStake(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "delegateStake", blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistry *VanillaregistrySession) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.DelegateStake(&_Vanillaregistry.TransactOpts, blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.DelegateStake(&_Vanillaregistry.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistry *VanillaregistryTransactor) ForceWithdrawalAsOwner(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "forceWithdrawalAsOwner", blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistry *VanillaregistrySession) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ForceWithdrawalAsOwner(&_Vanillaregistry.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ForceWithdrawalAsOwner(&_Vanillaregistry.TransactOpts, blsPubKeys, withdrawalAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistry *VanillaregistryTransactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "initialize", _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistry *VanillaregistrySession) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Initialize(&_Vanillaregistry.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Initialize(&_Vanillaregistry.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistry *VanillaregistryTransactor) ManuallyTransferSlashingFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "manuallyTransferSlashingFunds")
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistry *VanillaregistrySession) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ManuallyTransferSlashingFunds(&_Vanillaregistry.TransactOpts)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.ManuallyTransferSlashingFunds(&_Vanillaregistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistry *VanillaregistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistry *VanillaregistrySession) Pause() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Pause(&_Vanillaregistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Pause() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Pause(&_Vanillaregistry.TransactOpts)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistryTransactor) RemoveWhitelistedStakers(opts *bind.TransactOpts, stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "removeWhitelistedStakers", stakers)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistrySession) RemoveWhitelistedStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.RemoveWhitelistedStakers(&_Vanillaregistry.TransactOpts, stakers)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) RemoveWhitelistedStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.RemoveWhitelistedStakers(&_Vanillaregistry.TransactOpts, stakers)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistry *VanillaregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistry *VanillaregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.RenounceOwnership(&_Vanillaregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.RenounceOwnership(&_Vanillaregistry.TransactOpts)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistry *VanillaregistryTransactor) SetMinStake(opts *bind.TransactOpts, newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "setMinStake", newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistry *VanillaregistrySession) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetMinStake(&_Vanillaregistry.TransactOpts, newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetMinStake(&_Vanillaregistry.TransactOpts, newMinStake)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistry *VanillaregistryTransactor) SetSlashOracle(opts *bind.TransactOpts, newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "setSlashOracle", newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistry *VanillaregistrySession) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashOracle(&_Vanillaregistry.TransactOpts, newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashOracle(&_Vanillaregistry.TransactOpts, newSlashOracle)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistry *VanillaregistryTransactor) SetSlashReceiver(opts *bind.TransactOpts, newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "setSlashReceiver", newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistry *VanillaregistrySession) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashReceiver(&_Vanillaregistry.TransactOpts, newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashReceiver(&_Vanillaregistry.TransactOpts, newSlashReceiver)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistryTransactor) SetSlashingPayoutPeriodBlocks(opts *bind.TransactOpts, newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "setSlashingPayoutPeriodBlocks", newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistrySession) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashingPayoutPeriodBlocks(&_Vanillaregistry.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetSlashingPayoutPeriodBlocks(&_Vanillaregistry.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistryTransactor) SetUnstakePeriodBlocks(opts *bind.TransactOpts, newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "setUnstakePeriodBlocks", newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistrySession) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetUnstakePeriodBlocks(&_Vanillaregistry.TransactOpts, newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.SetUnstakePeriodBlocks(&_Vanillaregistry.TransactOpts, newUnstakePeriodBlocks)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistry *VanillaregistryTransactor) Slash(opts *bind.TransactOpts, blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "slash", blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistry *VanillaregistrySession) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Slash(&_Vanillaregistry.TransactOpts, blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Slash(&_Vanillaregistry.TransactOpts, blsPubKeys, payoutIfDue)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) Stake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "stake", blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistrySession) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Stake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Stake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistry *VanillaregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistry *VanillaregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.TransferOwnership(&_Vanillaregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.TransferOwnership(&_Vanillaregistry.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistry *VanillaregistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistry *VanillaregistrySession) Unpause() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Unpause(&_Vanillaregistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Unpause(&_Vanillaregistry.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistryTransactor) Unstake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "unstake", blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistrySession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Unstake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Unstake(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistry *VanillaregistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.UpgradeToAndCall(&_Vanillaregistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.UpgradeToAndCall(&_Vanillaregistry.TransactOpts, newImplementation, data)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistryTransactor) WhitelistStakers(opts *bind.TransactOpts, stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "whitelistStakers", stakers)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistrySession) WhitelistStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.WhitelistStakers(&_Vanillaregistry.TransactOpts, stakers)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) WhitelistStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.WhitelistStakers(&_Vanillaregistry.TransactOpts, stakers)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistryTransactor) Withdraw(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.Transact(opts, "withdraw", blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistrySession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Withdraw(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Withdraw(&_Vanillaregistry.TransactOpts, blsPubKeys)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistry.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistry *VanillaregistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Fallback(&_Vanillaregistry.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Fallback(&_Vanillaregistry.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistry *VanillaregistryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistry.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistry *VanillaregistrySession) Receive() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Receive(&_Vanillaregistry.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistry *VanillaregistryTransactorSession) Receive() (*types.Transaction, error) {
	return _Vanillaregistry.Contract.Receive(&_Vanillaregistry.TransactOpts)
}

// VanillaregistryFeeTransferIterator is returned from FilterFeeTransfer and is used to iterate over the raw logs and unpacked data for FeeTransfer events raised by the Vanillaregistry contract.
type VanillaregistryFeeTransferIterator struct {
	Event *VanillaregistryFeeTransfer // Event containing the contract specifics and raw log

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
func (it *VanillaregistryFeeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryFeeTransfer)
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
		it.Event = new(VanillaregistryFeeTransfer)
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
func (it *VanillaregistryFeeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryFeeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryFeeTransfer represents a FeeTransfer event raised by the Vanillaregistry contract.
type VanillaregistryFeeTransfer struct {
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeTransfer is a free log retrieval operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Vanillaregistry *VanillaregistryFilterer) FilterFeeTransfer(opts *bind.FilterOpts, recipient []common.Address) (*VanillaregistryFeeTransferIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryFeeTransferIterator{contract: _Vanillaregistry.contract, event: "FeeTransfer", logs: logs, sub: sub}, nil
}

// WatchFeeTransfer is a free log subscription operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Vanillaregistry *VanillaregistryFilterer) WatchFeeTransfer(opts *bind.WatchOpts, sink chan<- *VanillaregistryFeeTransfer, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryFeeTransfer)
				if err := _Vanillaregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseFeeTransfer(log types.Log) (*VanillaregistryFeeTransfer, error) {
	event := new(VanillaregistryFeeTransfer)
	if err := _Vanillaregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Vanillaregistry contract.
type VanillaregistryInitializedIterator struct {
	Event *VanillaregistryInitialized // Event containing the contract specifics and raw log

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
func (it *VanillaregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryInitialized)
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
		it.Event = new(VanillaregistryInitialized)
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
func (it *VanillaregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryInitialized represents a Initialized event raised by the Vanillaregistry contract.
type VanillaregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vanillaregistry *VanillaregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*VanillaregistryInitializedIterator, error) {

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VanillaregistryInitializedIterator{contract: _Vanillaregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vanillaregistry *VanillaregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VanillaregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryInitialized)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseInitialized(log types.Log) (*VanillaregistryInitialized, error) {
	event := new(VanillaregistryInitialized)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryMinStakeSetIterator is returned from FilterMinStakeSet and is used to iterate over the raw logs and unpacked data for MinStakeSet events raised by the Vanillaregistry contract.
type VanillaregistryMinStakeSetIterator struct {
	Event *VanillaregistryMinStakeSet // Event containing the contract specifics and raw log

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
func (it *VanillaregistryMinStakeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryMinStakeSet)
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
		it.Event = new(VanillaregistryMinStakeSet)
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
func (it *VanillaregistryMinStakeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryMinStakeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryMinStakeSet represents a MinStakeSet event raised by the Vanillaregistry contract.
type VanillaregistryMinStakeSet struct {
	MsgSender   common.Address
	NewMinStake *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMinStakeSet is a free log retrieval operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Vanillaregistry *VanillaregistryFilterer) FilterMinStakeSet(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistryMinStakeSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryMinStakeSetIterator{contract: _Vanillaregistry.contract, event: "MinStakeSet", logs: logs, sub: sub}, nil
}

// WatchMinStakeSet is a free log subscription operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Vanillaregistry *VanillaregistryFilterer) WatchMinStakeSet(opts *bind.WatchOpts, sink chan<- *VanillaregistryMinStakeSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryMinStakeSet)
				if err := _Vanillaregistry.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseMinStakeSet(log types.Log) (*VanillaregistryMinStakeSet, error) {
	event := new(VanillaregistryMinStakeSet)
	if err := _Vanillaregistry.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Vanillaregistry contract.
type VanillaregistryOwnershipTransferStartedIterator struct {
	Event *VanillaregistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *VanillaregistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryOwnershipTransferStarted)
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
		it.Event = new(VanillaregistryOwnershipTransferStarted)
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
func (it *VanillaregistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Vanillaregistry contract.
type VanillaregistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistry *VanillaregistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VanillaregistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryOwnershipTransferStartedIterator{contract: _Vanillaregistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistry *VanillaregistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *VanillaregistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryOwnershipTransferStarted)
				if err := _Vanillaregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*VanillaregistryOwnershipTransferStarted, error) {
	event := new(VanillaregistryOwnershipTransferStarted)
	if err := _Vanillaregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Vanillaregistry contract.
type VanillaregistryOwnershipTransferredIterator struct {
	Event *VanillaregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VanillaregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryOwnershipTransferred)
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
		it.Event = new(VanillaregistryOwnershipTransferred)
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
func (it *VanillaregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Vanillaregistry contract.
type VanillaregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistry *VanillaregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VanillaregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryOwnershipTransferredIterator{contract: _Vanillaregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistry *VanillaregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VanillaregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryOwnershipTransferred)
				if err := _Vanillaregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseOwnershipTransferred(log types.Log) (*VanillaregistryOwnershipTransferred, error) {
	event := new(VanillaregistryOwnershipTransferred)
	if err := _Vanillaregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Vanillaregistry contract.
type VanillaregistryPausedIterator struct {
	Event *VanillaregistryPaused // Event containing the contract specifics and raw log

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
func (it *VanillaregistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryPaused)
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
		it.Event = new(VanillaregistryPaused)
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
func (it *VanillaregistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryPaused represents a Paused event raised by the Vanillaregistry contract.
type VanillaregistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Vanillaregistry *VanillaregistryFilterer) FilterPaused(opts *bind.FilterOpts) (*VanillaregistryPausedIterator, error) {

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &VanillaregistryPausedIterator{contract: _Vanillaregistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Vanillaregistry *VanillaregistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *VanillaregistryPaused) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryPaused)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParsePaused(log types.Log) (*VanillaregistryPaused, error) {
	event := new(VanillaregistryPaused)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistrySlashOracleSetIterator is returned from FilterSlashOracleSet and is used to iterate over the raw logs and unpacked data for SlashOracleSet events raised by the Vanillaregistry contract.
type VanillaregistrySlashOracleSetIterator struct {
	Event *VanillaregistrySlashOracleSet // Event containing the contract specifics and raw log

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
func (it *VanillaregistrySlashOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistrySlashOracleSet)
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
		it.Event = new(VanillaregistrySlashOracleSet)
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
func (it *VanillaregistrySlashOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistrySlashOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistrySlashOracleSet represents a SlashOracleSet event raised by the Vanillaregistry contract.
type VanillaregistrySlashOracleSet struct {
	MsgSender      common.Address
	NewSlashOracle common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSlashOracleSet is a free log retrieval operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Vanillaregistry *VanillaregistryFilterer) FilterSlashOracleSet(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistrySlashOracleSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistrySlashOracleSetIterator{contract: _Vanillaregistry.contract, event: "SlashOracleSet", logs: logs, sub: sub}, nil
}

// WatchSlashOracleSet is a free log subscription operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Vanillaregistry *VanillaregistryFilterer) WatchSlashOracleSet(opts *bind.WatchOpts, sink chan<- *VanillaregistrySlashOracleSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistrySlashOracleSet)
				if err := _Vanillaregistry.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseSlashOracleSet(log types.Log) (*VanillaregistrySlashOracleSet, error) {
	event := new(VanillaregistrySlashOracleSet)
	if err := _Vanillaregistry.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistrySlashReceiverSetIterator is returned from FilterSlashReceiverSet and is used to iterate over the raw logs and unpacked data for SlashReceiverSet events raised by the Vanillaregistry contract.
type VanillaregistrySlashReceiverSetIterator struct {
	Event *VanillaregistrySlashReceiverSet // Event containing the contract specifics and raw log

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
func (it *VanillaregistrySlashReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistrySlashReceiverSet)
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
		it.Event = new(VanillaregistrySlashReceiverSet)
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
func (it *VanillaregistrySlashReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistrySlashReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistrySlashReceiverSet represents a SlashReceiverSet event raised by the Vanillaregistry contract.
type VanillaregistrySlashReceiverSet struct {
	MsgSender        common.Address
	NewSlashReceiver common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSlashReceiverSet is a free log retrieval operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Vanillaregistry *VanillaregistryFilterer) FilterSlashReceiverSet(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistrySlashReceiverSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistrySlashReceiverSetIterator{contract: _Vanillaregistry.contract, event: "SlashReceiverSet", logs: logs, sub: sub}, nil
}

// WatchSlashReceiverSet is a free log subscription operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Vanillaregistry *VanillaregistryFilterer) WatchSlashReceiverSet(opts *bind.WatchOpts, sink chan<- *VanillaregistrySlashReceiverSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistrySlashReceiverSet)
				if err := _Vanillaregistry.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseSlashReceiverSet(log types.Log) (*VanillaregistrySlashReceiverSet, error) {
	event := new(VanillaregistrySlashReceiverSet)
	if err := _Vanillaregistry.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistrySlashedIterator is returned from FilterSlashed and is used to iterate over the raw logs and unpacked data for Slashed events raised by the Vanillaregistry contract.
type VanillaregistrySlashedIterator struct {
	Event *VanillaregistrySlashed // Event containing the contract specifics and raw log

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
func (it *VanillaregistrySlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistrySlashed)
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
		it.Event = new(VanillaregistrySlashed)
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
func (it *VanillaregistrySlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistrySlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistrySlashed represents a Slashed event raised by the Vanillaregistry contract.
type VanillaregistrySlashed struct {
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
func (_Vanillaregistry *VanillaregistryFilterer) FilterSlashed(opts *bind.FilterOpts, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (*VanillaregistrySlashedIterator, error) {

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

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistrySlashedIterator{contract: _Vanillaregistry.contract, event: "Slashed", logs: logs, sub: sub}, nil
}

// WatchSlashed is a free log subscription operation binding the contract event 0xf15b8630ce764d5dbcfaaa9843c3e5fcdb460aaaa46d7dc3ff4f19ca4096fc07.
//
// Solidity: event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) WatchSlashed(opts *bind.WatchOpts, sink chan<- *VanillaregistrySlashed, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistrySlashed)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Slashed", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseSlashed(log types.Log) (*VanillaregistrySlashed, error) {
	event := new(VanillaregistrySlashed)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Slashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistrySlashingPayoutPeriodBlocksSetIterator is returned from FilterSlashingPayoutPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for SlashingPayoutPeriodBlocksSet events raised by the Vanillaregistry contract.
type VanillaregistrySlashingPayoutPeriodBlocksSetIterator struct {
	Event *VanillaregistrySlashingPayoutPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *VanillaregistrySlashingPayoutPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistrySlashingPayoutPeriodBlocksSet)
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
		it.Event = new(VanillaregistrySlashingPayoutPeriodBlocksSet)
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
func (it *VanillaregistrySlashingPayoutPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistrySlashingPayoutPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistrySlashingPayoutPeriodBlocksSet represents a SlashingPayoutPeriodBlocksSet event raised by the Vanillaregistry contract.
type VanillaregistrySlashingPayoutPeriodBlocksSet struct {
	MsgSender                     common.Address
	NewSlashingPayoutPeriodBlocks *big.Int
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterSlashingPayoutPeriodBlocksSet is a free log retrieval operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Vanillaregistry *VanillaregistryFilterer) FilterSlashingPayoutPeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistrySlashingPayoutPeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistrySlashingPayoutPeriodBlocksSetIterator{contract: _Vanillaregistry.contract, event: "SlashingPayoutPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchSlashingPayoutPeriodBlocksSet is a free log subscription operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Vanillaregistry *VanillaregistryFilterer) WatchSlashingPayoutPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *VanillaregistrySlashingPayoutPeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistrySlashingPayoutPeriodBlocksSet)
				if err := _Vanillaregistry.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseSlashingPayoutPeriodBlocksSet(log types.Log) (*VanillaregistrySlashingPayoutPeriodBlocksSet, error) {
	event := new(VanillaregistrySlashingPayoutPeriodBlocksSet)
	if err := _Vanillaregistry.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryStakeAddedIterator is returned from FilterStakeAdded and is used to iterate over the raw logs and unpacked data for StakeAdded events raised by the Vanillaregistry contract.
type VanillaregistryStakeAddedIterator struct {
	Event *VanillaregistryStakeAdded // Event containing the contract specifics and raw log

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
func (it *VanillaregistryStakeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryStakeAdded)
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
		it.Event = new(VanillaregistryStakeAdded)
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
func (it *VanillaregistryStakeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryStakeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryStakeAdded represents a StakeAdded event raised by the Vanillaregistry contract.
type VanillaregistryStakeAdded struct {
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
func (_Vanillaregistry *VanillaregistryFilterer) FilterStakeAdded(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*VanillaregistryStakeAddedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryStakeAddedIterator{contract: _Vanillaregistry.contract, event: "StakeAdded", logs: logs, sub: sub}, nil
}

// WatchStakeAdded is a free log subscription operation binding the contract event 0xb01516cc7ddda8b10127c714474503b38a75b9afa8a4e4b9da306e61181980c7.
//
// Solidity: event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance)
func (_Vanillaregistry *VanillaregistryFilterer) WatchStakeAdded(opts *bind.WatchOpts, sink chan<- *VanillaregistryStakeAdded, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryStakeAdded)
				if err := _Vanillaregistry.contract.UnpackLog(event, "StakeAdded", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseStakeAdded(log types.Log) (*VanillaregistryStakeAdded, error) {
	event := new(VanillaregistryStakeAdded)
	if err := _Vanillaregistry.contract.UnpackLog(event, "StakeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryStakeWithdrawnIterator is returned from FilterStakeWithdrawn and is used to iterate over the raw logs and unpacked data for StakeWithdrawn events raised by the Vanillaregistry contract.
type VanillaregistryStakeWithdrawnIterator struct {
	Event *VanillaregistryStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *VanillaregistryStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryStakeWithdrawn)
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
		it.Event = new(VanillaregistryStakeWithdrawn)
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
func (it *VanillaregistryStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryStakeWithdrawn represents a StakeWithdrawn event raised by the Vanillaregistry contract.
type VanillaregistryStakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdrawn is a free log retrieval operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) FilterStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*VanillaregistryStakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryStakeWithdrawnIterator{contract: _Vanillaregistry.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchStakeWithdrawn is a free log subscription operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *VanillaregistryStakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryStakeWithdrawn)
				if err := _Vanillaregistry.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseStakeWithdrawn(log types.Log) (*VanillaregistryStakeWithdrawn, error) {
	event := new(VanillaregistryStakeWithdrawn)
	if err := _Vanillaregistry.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Vanillaregistry contract.
type VanillaregistryStakedIterator struct {
	Event *VanillaregistryStaked // Event containing the contract specifics and raw log

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
func (it *VanillaregistryStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryStaked)
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
		it.Event = new(VanillaregistryStaked)
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
func (it *VanillaregistryStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryStaked represents a Staked event raised by the Vanillaregistry contract.
type VanillaregistryStaked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) FilterStaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*VanillaregistryStakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryStakedIterator{contract: _Vanillaregistry.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *VanillaregistryStaked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryStaked)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseStaked(log types.Log) (*VanillaregistryStaked, error) {
	event := new(VanillaregistryStaked)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryStakerRemovedFromWhitelistIterator is returned from FilterStakerRemovedFromWhitelist and is used to iterate over the raw logs and unpacked data for StakerRemovedFromWhitelist events raised by the Vanillaregistry contract.
type VanillaregistryStakerRemovedFromWhitelistIterator struct {
	Event *VanillaregistryStakerRemovedFromWhitelist // Event containing the contract specifics and raw log

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
func (it *VanillaregistryStakerRemovedFromWhitelistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryStakerRemovedFromWhitelist)
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
		it.Event = new(VanillaregistryStakerRemovedFromWhitelist)
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
func (it *VanillaregistryStakerRemovedFromWhitelistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryStakerRemovedFromWhitelistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryStakerRemovedFromWhitelist represents a StakerRemovedFromWhitelist event raised by the Vanillaregistry contract.
type VanillaregistryStakerRemovedFromWhitelist struct {
	MsgSender common.Address
	Staker    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStakerRemovedFromWhitelist is a free log retrieval operation binding the contract event 0x79fccea9ca325c6c715589de68546a63a35425178d1a8b1436bef7f7a76087a4.
//
// Solidity: event StakerRemovedFromWhitelist(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) FilterStakerRemovedFromWhitelist(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistryStakerRemovedFromWhitelistIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "StakerRemovedFromWhitelist", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryStakerRemovedFromWhitelistIterator{contract: _Vanillaregistry.contract, event: "StakerRemovedFromWhitelist", logs: logs, sub: sub}, nil
}

// WatchStakerRemovedFromWhitelist is a free log subscription operation binding the contract event 0x79fccea9ca325c6c715589de68546a63a35425178d1a8b1436bef7f7a76087a4.
//
// Solidity: event StakerRemovedFromWhitelist(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) WatchStakerRemovedFromWhitelist(opts *bind.WatchOpts, sink chan<- *VanillaregistryStakerRemovedFromWhitelist, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "StakerRemovedFromWhitelist", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryStakerRemovedFromWhitelist)
				if err := _Vanillaregistry.contract.UnpackLog(event, "StakerRemovedFromWhitelist", log); err != nil {
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

// ParseStakerRemovedFromWhitelist is a log parse operation binding the contract event 0x79fccea9ca325c6c715589de68546a63a35425178d1a8b1436bef7f7a76087a4.
//
// Solidity: event StakerRemovedFromWhitelist(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) ParseStakerRemovedFromWhitelist(log types.Log) (*VanillaregistryStakerRemovedFromWhitelist, error) {
	event := new(VanillaregistryStakerRemovedFromWhitelist)
	if err := _Vanillaregistry.contract.UnpackLog(event, "StakerRemovedFromWhitelist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryStakerWhitelistedIterator is returned from FilterStakerWhitelisted and is used to iterate over the raw logs and unpacked data for StakerWhitelisted events raised by the Vanillaregistry contract.
type VanillaregistryStakerWhitelistedIterator struct {
	Event *VanillaregistryStakerWhitelisted // Event containing the contract specifics and raw log

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
func (it *VanillaregistryStakerWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryStakerWhitelisted)
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
		it.Event = new(VanillaregistryStakerWhitelisted)
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
func (it *VanillaregistryStakerWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryStakerWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryStakerWhitelisted represents a StakerWhitelisted event raised by the Vanillaregistry contract.
type VanillaregistryStakerWhitelisted struct {
	MsgSender common.Address
	Staker    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStakerWhitelisted is a free log retrieval operation binding the contract event 0xc94be990bcb425c41f84d9928ad9082894fb4f1f1d2508bc088d3a6a3e4059db.
//
// Solidity: event StakerWhitelisted(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) FilterStakerWhitelisted(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistryStakerWhitelistedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "StakerWhitelisted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryStakerWhitelistedIterator{contract: _Vanillaregistry.contract, event: "StakerWhitelisted", logs: logs, sub: sub}, nil
}

// WatchStakerWhitelisted is a free log subscription operation binding the contract event 0xc94be990bcb425c41f84d9928ad9082894fb4f1f1d2508bc088d3a6a3e4059db.
//
// Solidity: event StakerWhitelisted(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) WatchStakerWhitelisted(opts *bind.WatchOpts, sink chan<- *VanillaregistryStakerWhitelisted, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "StakerWhitelisted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryStakerWhitelisted)
				if err := _Vanillaregistry.contract.UnpackLog(event, "StakerWhitelisted", log); err != nil {
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

// ParseStakerWhitelisted is a log parse operation binding the contract event 0xc94be990bcb425c41f84d9928ad9082894fb4f1f1d2508bc088d3a6a3e4059db.
//
// Solidity: event StakerWhitelisted(address indexed msgSender, address staker)
func (_Vanillaregistry *VanillaregistryFilterer) ParseStakerWhitelisted(log types.Log) (*VanillaregistryStakerWhitelisted, error) {
	event := new(VanillaregistryStakerWhitelisted)
	if err := _Vanillaregistry.contract.UnpackLog(event, "StakerWhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryTotalStakeWithdrawnIterator is returned from FilterTotalStakeWithdrawn and is used to iterate over the raw logs and unpacked data for TotalStakeWithdrawn events raised by the Vanillaregistry contract.
type VanillaregistryTotalStakeWithdrawnIterator struct {
	Event *VanillaregistryTotalStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *VanillaregistryTotalStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryTotalStakeWithdrawn)
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
		it.Event = new(VanillaregistryTotalStakeWithdrawn)
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
func (it *VanillaregistryTotalStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryTotalStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryTotalStakeWithdrawn represents a TotalStakeWithdrawn event raised by the Vanillaregistry contract.
type VanillaregistryTotalStakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	TotalAmount       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTotalStakeWithdrawn is a free log retrieval operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Vanillaregistry *VanillaregistryFilterer) FilterTotalStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*VanillaregistryTotalStakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryTotalStakeWithdrawnIterator{contract: _Vanillaregistry.contract, event: "TotalStakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTotalStakeWithdrawn is a free log subscription operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Vanillaregistry *VanillaregistryFilterer) WatchTotalStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *VanillaregistryTotalStakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryTotalStakeWithdrawn)
				if err := _Vanillaregistry.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseTotalStakeWithdrawn(log types.Log) (*VanillaregistryTotalStakeWithdrawn, error) {
	event := new(VanillaregistryTotalStakeWithdrawn)
	if err := _Vanillaregistry.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Vanillaregistry contract.
type VanillaregistryUnpausedIterator struct {
	Event *VanillaregistryUnpaused // Event containing the contract specifics and raw log

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
func (it *VanillaregistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryUnpaused)
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
		it.Event = new(VanillaregistryUnpaused)
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
func (it *VanillaregistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryUnpaused represents a Unpaused event raised by the Vanillaregistry contract.
type VanillaregistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Vanillaregistry *VanillaregistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*VanillaregistryUnpausedIterator, error) {

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &VanillaregistryUnpausedIterator{contract: _Vanillaregistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Vanillaregistry *VanillaregistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *VanillaregistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryUnpaused)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseUnpaused(log types.Log) (*VanillaregistryUnpaused, error) {
	event := new(VanillaregistryUnpaused)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryUnstakePeriodBlocksSetIterator is returned from FilterUnstakePeriodBlocksSet and is used to iterate over the raw logs and unpacked data for UnstakePeriodBlocksSet events raised by the Vanillaregistry contract.
type VanillaregistryUnstakePeriodBlocksSetIterator struct {
	Event *VanillaregistryUnstakePeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *VanillaregistryUnstakePeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryUnstakePeriodBlocksSet)
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
		it.Event = new(VanillaregistryUnstakePeriodBlocksSet)
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
func (it *VanillaregistryUnstakePeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryUnstakePeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryUnstakePeriodBlocksSet represents a UnstakePeriodBlocksSet event raised by the Vanillaregistry contract.
type VanillaregistryUnstakePeriodBlocksSet struct {
	MsgSender              common.Address
	NewUnstakePeriodBlocks *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUnstakePeriodBlocksSet is a free log retrieval operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Vanillaregistry *VanillaregistryFilterer) FilterUnstakePeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*VanillaregistryUnstakePeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryUnstakePeriodBlocksSetIterator{contract: _Vanillaregistry.contract, event: "UnstakePeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchUnstakePeriodBlocksSet is a free log subscription operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Vanillaregistry *VanillaregistryFilterer) WatchUnstakePeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *VanillaregistryUnstakePeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryUnstakePeriodBlocksSet)
				if err := _Vanillaregistry.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseUnstakePeriodBlocksSet(log types.Log) (*VanillaregistryUnstakePeriodBlocksSet, error) {
	event := new(VanillaregistryUnstakePeriodBlocksSet)
	if err := _Vanillaregistry.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Vanillaregistry contract.
type VanillaregistryUnstakedIterator struct {
	Event *VanillaregistryUnstaked // Event containing the contract specifics and raw log

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
func (it *VanillaregistryUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryUnstaked)
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
		it.Event = new(VanillaregistryUnstaked)
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
func (it *VanillaregistryUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryUnstaked represents a Unstaked event raised by the Vanillaregistry contract.
type VanillaregistryUnstaked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) FilterUnstaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*VanillaregistryUnstakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryUnstakedIterator{contract: _Vanillaregistry.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistry *VanillaregistryFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *VanillaregistryUnstaked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryUnstaked)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Unstaked", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseUnstaked(log types.Log) (*VanillaregistryUnstaked, error) {
	event := new(VanillaregistryUnstaked)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VanillaregistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Vanillaregistry contract.
type VanillaregistryUpgradedIterator struct {
	Event *VanillaregistryUpgraded // Event containing the contract specifics and raw log

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
func (it *VanillaregistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VanillaregistryUpgraded)
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
		it.Event = new(VanillaregistryUpgraded)
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
func (it *VanillaregistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VanillaregistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VanillaregistryUpgraded represents a Upgraded event raised by the Vanillaregistry contract.
type VanillaregistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Vanillaregistry *VanillaregistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*VanillaregistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Vanillaregistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &VanillaregistryUpgradedIterator{contract: _Vanillaregistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Vanillaregistry *VanillaregistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *VanillaregistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Vanillaregistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VanillaregistryUpgraded)
				if err := _Vanillaregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Vanillaregistry *VanillaregistryFilterer) ParseUpgraded(log types.Log) (*VanillaregistryUpgraded, error) {
	event := new(VanillaregistryUpgraded)
	if err := _Vanillaregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
