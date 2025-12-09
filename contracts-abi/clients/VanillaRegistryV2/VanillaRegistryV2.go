// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vanillaregistryv2

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

// IVanillaRegistryV2StakedValidator is an auto generated low-level Go binding around an user-defined struct.
type IVanillaRegistryV2StakedValidator struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}

// Vanillaregistryv2MetaData contains all meta data concerning the Vanillaregistryv2 contract.
var Vanillaregistryv2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"claimForceWithdrawnFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delegateStake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"forceWithdrawalAsOwner\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"forceWithdrawnFunds\",\"inputs\":[{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amountToClaim\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccumulatedSlashingFunds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBlocksTillWithdrawAllowed\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedAmount\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStakedValidator\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVanillaRegistryV2.StakedValidator\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashOracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_slashReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_unstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSlashingPayoutDue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyTransferSlashingFunds\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeWhitelistedStakers\",\"inputs\":[{\"name\":\"stakers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"newMinStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashOracle\",\"inputs\":[{\"name\":\"newSlashOracle\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashReceiver\",\"inputs\":[{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashingPayoutPeriodBlocks\",\"inputs\":[{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnstakePeriodBlocks\",\"inputs\":[{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slash\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"payoutIfDue\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashingFundsTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accumulatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastPayoutBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"stakedValidators\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unstakeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstakePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"whitelistStakers\",\"inputs\":[{\"name\":\"stakers\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"whitelistedStakers\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"whitelisted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"blsPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"FeeTransfer\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMinStake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashOracleSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashOracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashReceiverSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Slashed\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashReceiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashingPayoutPeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newSlashingPayoutPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeAdded\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newBalance\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Staked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerRemovedFromWhitelist\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakerWhitelisted\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"staker\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TotalStakeWithdrawn\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnstakePeriodBlocksSet\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newUnstakePeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstaked\",\"inputs\":[{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AtLeastOneRecipientRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeRecipientIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustUnstakeToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoFundsToWithdraw\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotSlashOracle\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slashOracle\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotWhitelistedStaker\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotWithdrawalAddress\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlashAmountMustBeLessThanMinStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashAmountMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashOracleMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashReceiverMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingPayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlashingTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakeTooLowForNumberOfKeys\",\"inputs\":[{\"name\":\"msgValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StakerAlreadyWhitelisted\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StakerNotWhitelisted\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TransferToRecipientFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnstakePeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorCannotBeUnstaking\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordMustNotExist\",\"inputs\":[{\"name\":\"valBLSPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMismatch\",\"inputs\":[{\"name\":\"actualWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"expectedWithdrawalAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WithdrawalAddressMustBeSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawalFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WithdrawingTooSoon\",\"inputs\":[]}]",
}

// Vanillaregistryv2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Vanillaregistryv2MetaData.ABI instead.
var Vanillaregistryv2ABI = Vanillaregistryv2MetaData.ABI

// Vanillaregistryv2 is an auto generated Go binding around an Ethereum contract.
type Vanillaregistryv2 struct {
	Vanillaregistryv2Caller     // Read-only binding to the contract
	Vanillaregistryv2Transactor // Write-only binding to the contract
	Vanillaregistryv2Filterer   // Log filterer for contract events
}

// Vanillaregistryv2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Vanillaregistryv2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Vanillaregistryv2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Vanillaregistryv2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Vanillaregistryv2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Vanillaregistryv2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Vanillaregistryv2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Vanillaregistryv2Session struct {
	Contract     *Vanillaregistryv2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Vanillaregistryv2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Vanillaregistryv2CallerSession struct {
	Contract *Vanillaregistryv2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// Vanillaregistryv2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Vanillaregistryv2TransactorSession struct {
	Contract     *Vanillaregistryv2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// Vanillaregistryv2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Vanillaregistryv2Raw struct {
	Contract *Vanillaregistryv2 // Generic contract binding to access the raw methods on
}

// Vanillaregistryv2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Vanillaregistryv2CallerRaw struct {
	Contract *Vanillaregistryv2Caller // Generic read-only contract binding to access the raw methods on
}

// Vanillaregistryv2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Vanillaregistryv2TransactorRaw struct {
	Contract *Vanillaregistryv2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewVanillaregistryv2 creates a new instance of Vanillaregistryv2, bound to a specific deployed contract.
func NewVanillaregistryv2(address common.Address, backend bind.ContractBackend) (*Vanillaregistryv2, error) {
	contract, err := bindVanillaregistryv2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2{Vanillaregistryv2Caller: Vanillaregistryv2Caller{contract: contract}, Vanillaregistryv2Transactor: Vanillaregistryv2Transactor{contract: contract}, Vanillaregistryv2Filterer: Vanillaregistryv2Filterer{contract: contract}}, nil
}

// NewVanillaregistryv2Caller creates a new read-only instance of Vanillaregistryv2, bound to a specific deployed contract.
func NewVanillaregistryv2Caller(address common.Address, caller bind.ContractCaller) (*Vanillaregistryv2Caller, error) {
	contract, err := bindVanillaregistryv2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2Caller{contract: contract}, nil
}

// NewVanillaregistryv2Transactor creates a new write-only instance of Vanillaregistryv2, bound to a specific deployed contract.
func NewVanillaregistryv2Transactor(address common.Address, transactor bind.ContractTransactor) (*Vanillaregistryv2Transactor, error) {
	contract, err := bindVanillaregistryv2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2Transactor{contract: contract}, nil
}

// NewVanillaregistryv2Filterer creates a new log filterer instance of Vanillaregistryv2, bound to a specific deployed contract.
func NewVanillaregistryv2Filterer(address common.Address, filterer bind.ContractFilterer) (*Vanillaregistryv2Filterer, error) {
	contract, err := bindVanillaregistryv2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2Filterer{contract: contract}, nil
}

// bindVanillaregistryv2 binds a generic wrapper to an already deployed contract.
func bindVanillaregistryv2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Vanillaregistryv2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vanillaregistryv2 *Vanillaregistryv2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vanillaregistryv2.Contract.Vanillaregistryv2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vanillaregistryv2 *Vanillaregistryv2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Vanillaregistryv2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vanillaregistryv2 *Vanillaregistryv2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Vanillaregistryv2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vanillaregistryv2 *Vanillaregistryv2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vanillaregistryv2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _Vanillaregistryv2.Contract.UPGRADEINTERFACEVERSION(&_Vanillaregistryv2.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Vanillaregistryv2.Contract.UPGRADEINTERFACEVERSION(&_Vanillaregistryv2.CallOpts)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) ForceWithdrawnFunds(opts *bind.CallOpts, withdrawalAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "forceWithdrawnFunds", withdrawalAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.ForceWithdrawnFunds(&_Vanillaregistryv2.CallOpts, withdrawalAddress)
}

// ForceWithdrawnFunds is a free data retrieval call binding the contract method 0x3de24562.
//
// Solidity: function forceWithdrawnFunds(address withdrawalAddress) view returns(uint256 amountToClaim)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) ForceWithdrawnFunds(withdrawalAddress common.Address) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.ForceWithdrawnFunds(&_Vanillaregistryv2.CallOpts, withdrawalAddress)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) GetAccumulatedSlashingFunds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "getAccumulatedSlashingFunds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetAccumulatedSlashingFunds(&_Vanillaregistryv2.CallOpts)
}

// GetAccumulatedSlashingFunds is a free data retrieval call binding the contract method 0x5ddae85d.
//
// Solidity: function getAccumulatedSlashingFunds() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) GetAccumulatedSlashingFunds() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetAccumulatedSlashingFunds(&_Vanillaregistryv2.CallOpts)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) GetBlocksTillWithdrawAllowed(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "getBlocksTillWithdrawAllowed", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetBlocksTillWithdrawAllowed(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// GetBlocksTillWithdrawAllowed is a free data retrieval call binding the contract method 0x14699cb9.
//
// Solidity: function getBlocksTillWithdrawAllowed(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) GetBlocksTillWithdrawAllowed(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetBlocksTillWithdrawAllowed(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) GetStakedAmount(opts *bind.CallOpts, valBLSPubKey []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "getStakedAmount", valBLSPubKey)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetStakedAmount(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// GetStakedAmount is a free data retrieval call binding the contract method 0xb2a453e6.
//
// Solidity: function getStakedAmount(bytes valBLSPubKey) view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) GetStakedAmount(valBLSPubKey []byte) (*big.Int, error) {
	return _Vanillaregistryv2.Contract.GetStakedAmount(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) GetStakedValidator(opts *bind.CallOpts, valBLSPubKey []byte) (IVanillaRegistryV2StakedValidator, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "getStakedValidator", valBLSPubKey)

	if err != nil {
		return *new(IVanillaRegistryV2StakedValidator), err
	}

	out0 := *abi.ConvertType(out[0], new(IVanillaRegistryV2StakedValidator)).(*IVanillaRegistryV2StakedValidator)

	return out0, err

}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistryv2 *Vanillaregistryv2Session) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryV2StakedValidator, error) {
	return _Vanillaregistryv2.Contract.GetStakedValidator(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// GetStakedValidator is a free data retrieval call binding the contract method 0x1fc7c7c8.
//
// Solidity: function getStakedValidator(bytes valBLSPubKey) view returns((bool,address,uint256,(bool,uint256)))
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) GetStakedValidator(valBLSPubKey []byte) (IVanillaRegistryV2StakedValidator, error) {
	return _Vanillaregistryv2.Contract.GetStakedValidator(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) IsSlashingPayoutDue(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "isSlashingPayoutDue")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) IsSlashingPayoutDue() (bool, error) {
	return _Vanillaregistryv2.Contract.IsSlashingPayoutDue(&_Vanillaregistryv2.CallOpts)
}

// IsSlashingPayoutDue is a free data retrieval call binding the contract method 0x35fe201b.
//
// Solidity: function isSlashingPayoutDue() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) IsSlashingPayoutDue() (bool, error) {
	return _Vanillaregistryv2.Contract.IsSlashingPayoutDue(&_Vanillaregistryv2.CallOpts)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) IsUnstaking(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "isUnstaking", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistryv2.Contract.IsUnstaking(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// IsUnstaking is a free data retrieval call binding the contract method 0x388a7968.
//
// Solidity: function isUnstaking(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) IsUnstaking(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistryv2.Contract.IsUnstaking(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) IsValidatorOptedIn(opts *bind.CallOpts, valBLSPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "isValidatorOptedIn", valBLSPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistryv2.Contract.IsValidatorOptedIn(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valBLSPubKey) view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) IsValidatorOptedIn(valBLSPubKey []byte) (bool, error) {
	return _Vanillaregistryv2.Contract.IsValidatorOptedIn(&_Vanillaregistryv2.CallOpts, valBLSPubKey)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) MinStake() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.MinStake(&_Vanillaregistryv2.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) MinStake() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.MinStake(&_Vanillaregistryv2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Owner() (common.Address, error) {
	return _Vanillaregistryv2.Contract.Owner(&_Vanillaregistryv2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) Owner() (common.Address, error) {
	return _Vanillaregistryv2.Contract.Owner(&_Vanillaregistryv2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Paused() (bool, error) {
	return _Vanillaregistryv2.Contract.Paused(&_Vanillaregistryv2.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) Paused() (bool, error) {
	return _Vanillaregistryv2.Contract.Paused(&_Vanillaregistryv2.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) PendingOwner() (common.Address, error) {
	return _Vanillaregistryv2.Contract.PendingOwner(&_Vanillaregistryv2.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) PendingOwner() (common.Address, error) {
	return _Vanillaregistryv2.Contract.PendingOwner(&_Vanillaregistryv2.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) ProxiableUUID() ([32]byte, error) {
	return _Vanillaregistryv2.Contract.ProxiableUUID(&_Vanillaregistryv2.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) ProxiableUUID() ([32]byte, error) {
	return _Vanillaregistryv2.Contract.ProxiableUUID(&_Vanillaregistryv2.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) SlashOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "slashOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SlashOracle() (common.Address, error) {
	return _Vanillaregistryv2.Contract.SlashOracle(&_Vanillaregistryv2.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) SlashOracle() (common.Address, error) {
	return _Vanillaregistryv2.Contract.SlashOracle(&_Vanillaregistryv2.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) SlashingFundsTracker(opts *bind.CallOpts) (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "slashingFundsTracker")

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
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Vanillaregistryv2.Contract.SlashingFundsTracker(&_Vanillaregistryv2.CallOpts)
}

// SlashingFundsTracker is a free data retrieval call binding the contract method 0x6f0301bd.
//
// Solidity: function slashingFundsTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) SlashingFundsTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Vanillaregistryv2.Contract.SlashingFundsTracker(&_Vanillaregistryv2.CallOpts)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) StakedValidators(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "stakedValidators", arg0)

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
func (_Vanillaregistryv2 *Vanillaregistryv2Session) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Vanillaregistryv2.Contract.StakedValidators(&_Vanillaregistryv2.CallOpts, arg0)
}

// StakedValidators is a free data retrieval call binding the contract method 0xfced6425.
//
// Solidity: function stakedValidators(bytes ) view returns(bool exists, address withdrawalAddress, uint256 balance, (bool,uint256) unstakeOccurrence)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) StakedValidators(arg0 []byte) (struct {
	Exists            bool
	WithdrawalAddress common.Address
	Balance           *big.Int
	UnstakeOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Vanillaregistryv2.Contract.StakedValidators(&_Vanillaregistryv2.CallOpts, arg0)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) UnstakePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "unstakePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) UnstakePeriodBlocks() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.UnstakePeriodBlocks(&_Vanillaregistryv2.CallOpts)
}

// UnstakePeriodBlocks is a free data retrieval call binding the contract method 0xc253f765.
//
// Solidity: function unstakePeriodBlocks() view returns(uint256)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) UnstakePeriodBlocks() (*big.Int, error) {
	return _Vanillaregistryv2.Contract.UnstakePeriodBlocks(&_Vanillaregistryv2.CallOpts)
}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistryv2 *Vanillaregistryv2Caller) WhitelistedStakers(opts *bind.CallOpts, staker common.Address) (bool, error) {
	var out []interface{}
	err := _Vanillaregistryv2.contract.Call(opts, &out, "whitelistedStakers", staker)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistryv2 *Vanillaregistryv2Session) WhitelistedStakers(staker common.Address) (bool, error) {
	return _Vanillaregistryv2.Contract.WhitelistedStakers(&_Vanillaregistryv2.CallOpts, staker)
}

// WhitelistedStakers is a free data retrieval call binding the contract method 0xfdaf17f0.
//
// Solidity: function whitelistedStakers(address staker) view returns(bool whitelisted)
func (_Vanillaregistryv2 *Vanillaregistryv2CallerSession) WhitelistedStakers(staker common.Address) (bool, error) {
	return _Vanillaregistryv2.Contract.WhitelistedStakers(&_Vanillaregistryv2.CallOpts, staker)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) AcceptOwnership() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.AcceptOwnership(&_Vanillaregistryv2.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.AcceptOwnership(&_Vanillaregistryv2.TransactOpts)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) AddStake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "addStake", blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.AddStake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// AddStake is a paid mutator transaction binding the contract method 0x92afedf6.
//
// Solidity: function addStake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) AddStake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.AddStake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) ClaimForceWithdrawnFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "claimForceWithdrawnFunds")
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ClaimForceWithdrawnFunds(&_Vanillaregistryv2.TransactOpts)
}

// ClaimForceWithdrawnFunds is a paid mutator transaction binding the contract method 0xf55690fd.
//
// Solidity: function claimForceWithdrawnFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) ClaimForceWithdrawnFunds() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ClaimForceWithdrawnFunds(&_Vanillaregistryv2.TransactOpts)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) DelegateStake(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "delegateStake", blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.DelegateStake(&_Vanillaregistryv2.TransactOpts, blsPubKeys, withdrawalAddress)
}

// DelegateStake is a paid mutator transaction binding the contract method 0x4b7952b3.
//
// Solidity: function delegateStake(bytes[] blsPubKeys, address withdrawalAddress) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) DelegateStake(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.DelegateStake(&_Vanillaregistryv2.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) ForceWithdrawalAsOwner(opts *bind.TransactOpts, blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "forceWithdrawalAsOwner", blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ForceWithdrawalAsOwner(&_Vanillaregistryv2.TransactOpts, blsPubKeys, withdrawalAddress)
}

// ForceWithdrawalAsOwner is a paid mutator transaction binding the contract method 0x7cadea98.
//
// Solidity: function forceWithdrawalAsOwner(bytes[] blsPubKeys, address withdrawalAddress) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) ForceWithdrawalAsOwner(blsPubKeys [][]byte, withdrawalAddress common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ForceWithdrawalAsOwner(&_Vanillaregistryv2.TransactOpts, blsPubKeys, withdrawalAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "initialize", _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Initialize(&_Vanillaregistryv2.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xacfb89fd.
//
// Solidity: function initialize(uint256 _minStake, address _slashOracle, address _slashReceiver, uint256 _unstakePeriodBlocks, uint256 _slashingPayoutPeriodBlocks, address _owner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Initialize(_minStake *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _unstakePeriodBlocks *big.Int, _slashingPayoutPeriodBlocks *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Initialize(&_Vanillaregistryv2.TransactOpts, _minStake, _slashOracle, _slashReceiver, _unstakePeriodBlocks, _slashingPayoutPeriodBlocks, _owner)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) ManuallyTransferSlashingFunds(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "manuallyTransferSlashingFunds")
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ManuallyTransferSlashingFunds(&_Vanillaregistryv2.TransactOpts)
}

// ManuallyTransferSlashingFunds is a paid mutator transaction binding the contract method 0xa1d694eb.
//
// Solidity: function manuallyTransferSlashingFunds() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) ManuallyTransferSlashingFunds() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.ManuallyTransferSlashingFunds(&_Vanillaregistryv2.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Pause() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Pause(&_Vanillaregistryv2.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Pause() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Pause(&_Vanillaregistryv2.TransactOpts)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) RemoveWhitelistedStakers(opts *bind.TransactOpts, stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "removeWhitelistedStakers", stakers)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) RemoveWhitelistedStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.RemoveWhitelistedStakers(&_Vanillaregistryv2.TransactOpts, stakers)
}

// RemoveWhitelistedStakers is a paid mutator transaction binding the contract method 0x5158c9fa.
//
// Solidity: function removeWhitelistedStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) RemoveWhitelistedStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.RemoveWhitelistedStakers(&_Vanillaregistryv2.TransactOpts, stakers)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) RenounceOwnership() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.RenounceOwnership(&_Vanillaregistryv2.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.RenounceOwnership(&_Vanillaregistryv2.TransactOpts)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) SetMinStake(opts *bind.TransactOpts, newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "setMinStake", newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetMinStake(&_Vanillaregistryv2.TransactOpts, newMinStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 newMinStake) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) SetMinStake(newMinStake *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetMinStake(&_Vanillaregistryv2.TransactOpts, newMinStake)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) SetSlashOracle(opts *bind.TransactOpts, newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "setSlashOracle", newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashOracle(&_Vanillaregistryv2.TransactOpts, newSlashOracle)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address newSlashOracle) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) SetSlashOracle(newSlashOracle common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashOracle(&_Vanillaregistryv2.TransactOpts, newSlashOracle)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) SetSlashReceiver(opts *bind.TransactOpts, newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "setSlashReceiver", newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashReceiver(&_Vanillaregistryv2.TransactOpts, newSlashReceiver)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address newSlashReceiver) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) SetSlashReceiver(newSlashReceiver common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashReceiver(&_Vanillaregistryv2.TransactOpts, newSlashReceiver)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) SetSlashingPayoutPeriodBlocks(opts *bind.TransactOpts, newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "setSlashingPayoutPeriodBlocks", newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashingPayoutPeriodBlocks(&_Vanillaregistryv2.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetSlashingPayoutPeriodBlocks is a paid mutator transaction binding the contract method 0xc4828f6b.
//
// Solidity: function setSlashingPayoutPeriodBlocks(uint256 newSlashingPayoutPeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) SetSlashingPayoutPeriodBlocks(newSlashingPayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetSlashingPayoutPeriodBlocks(&_Vanillaregistryv2.TransactOpts, newSlashingPayoutPeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) SetUnstakePeriodBlocks(opts *bind.TransactOpts, newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "setUnstakePeriodBlocks", newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetUnstakePeriodBlocks(&_Vanillaregistryv2.TransactOpts, newUnstakePeriodBlocks)
}

// SetUnstakePeriodBlocks is a paid mutator transaction binding the contract method 0xbc325c59.
//
// Solidity: function setUnstakePeriodBlocks(uint256 newUnstakePeriodBlocks) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) SetUnstakePeriodBlocks(newUnstakePeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.SetUnstakePeriodBlocks(&_Vanillaregistryv2.TransactOpts, newUnstakePeriodBlocks)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Slash(opts *bind.TransactOpts, blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "slash", blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Slash(&_Vanillaregistryv2.TransactOpts, blsPubKeys, payoutIfDue)
}

// Slash is a paid mutator transaction binding the contract method 0x7aa7dc14.
//
// Solidity: function slash(bytes[] blsPubKeys, bool payoutIfDue) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Slash(blsPubKeys [][]byte, payoutIfDue bool) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Slash(&_Vanillaregistryv2.TransactOpts, blsPubKeys, payoutIfDue)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Stake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "stake", blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Stake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// Stake is a paid mutator transaction binding the contract method 0x7299e0e6.
//
// Solidity: function stake(bytes[] blsPubKeys) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Stake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Stake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.TransferOwnership(&_Vanillaregistryv2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.TransferOwnership(&_Vanillaregistryv2.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Unpause() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Unpause(&_Vanillaregistryv2.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Unpause() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Unpause(&_Vanillaregistryv2.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Unstake(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "unstake", blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Unstake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// Unstake is a paid mutator transaction binding the contract method 0xc08a2081.
//
// Solidity: function unstake(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Unstake(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Unstake(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.UpgradeToAndCall(&_Vanillaregistryv2.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.UpgradeToAndCall(&_Vanillaregistryv2.TransactOpts, newImplementation, data)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) WhitelistStakers(opts *bind.TransactOpts, stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "whitelistStakers", stakers)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) WhitelistStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.WhitelistStakers(&_Vanillaregistryv2.TransactOpts, stakers)
}

// WhitelistStakers is a paid mutator transaction binding the contract method 0x95ebffdb.
//
// Solidity: function whitelistStakers(address[] stakers) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) WhitelistStakers(stakers []common.Address) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.WhitelistStakers(&_Vanillaregistryv2.TransactOpts, stakers)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Withdraw(opts *bind.TransactOpts, blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.Transact(opts, "withdraw", blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Withdraw(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// Withdraw is a paid mutator transaction binding the contract method 0xdcb1edcb.
//
// Solidity: function withdraw(bytes[] blsPubKeys) returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Withdraw(blsPubKeys [][]byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Withdraw(&_Vanillaregistryv2.TransactOpts, blsPubKeys)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Fallback(&_Vanillaregistryv2.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Fallback(&_Vanillaregistryv2.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vanillaregistryv2.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2Session) Receive() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Receive(&_Vanillaregistryv2.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Vanillaregistryv2 *Vanillaregistryv2TransactorSession) Receive() (*types.Transaction, error) {
	return _Vanillaregistryv2.Contract.Receive(&_Vanillaregistryv2.TransactOpts)
}

// Vanillaregistryv2FeeTransferIterator is returned from FilterFeeTransfer and is used to iterate over the raw logs and unpacked data for FeeTransfer events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2FeeTransferIterator struct {
	Event *Vanillaregistryv2FeeTransfer // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2FeeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2FeeTransfer)
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
		it.Event = new(Vanillaregistryv2FeeTransfer)
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
func (it *Vanillaregistryv2FeeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2FeeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2FeeTransfer represents a FeeTransfer event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2FeeTransfer struct {
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeTransfer is a free log retrieval operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterFeeTransfer(opts *bind.FilterOpts, recipient []common.Address) (*Vanillaregistryv2FeeTransferIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2FeeTransferIterator{contract: _Vanillaregistryv2.contract, event: "FeeTransfer", logs: logs, sub: sub}, nil
}

// WatchFeeTransfer is a free log subscription operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchFeeTransfer(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2FeeTransfer, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2FeeTransfer)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseFeeTransfer(log types.Log) (*Vanillaregistryv2FeeTransfer, error) {
	event := new(Vanillaregistryv2FeeTransfer)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2InitializedIterator struct {
	Event *Vanillaregistryv2Initialized // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Initialized)
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
		it.Event = new(Vanillaregistryv2Initialized)
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
func (it *Vanillaregistryv2InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Initialized represents a Initialized event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterInitialized(opts *bind.FilterOpts) (*Vanillaregistryv2InitializedIterator, error) {

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2InitializedIterator{contract: _Vanillaregistryv2.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Initialized) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Initialized)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseInitialized(log types.Log) (*Vanillaregistryv2Initialized, error) {
	event := new(Vanillaregistryv2Initialized)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2MinStakeSetIterator is returned from FilterMinStakeSet and is used to iterate over the raw logs and unpacked data for MinStakeSet events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2MinStakeSetIterator struct {
	Event *Vanillaregistryv2MinStakeSet // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2MinStakeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2MinStakeSet)
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
		it.Event = new(Vanillaregistryv2MinStakeSet)
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
func (it *Vanillaregistryv2MinStakeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2MinStakeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2MinStakeSet represents a MinStakeSet event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2MinStakeSet struct {
	MsgSender   common.Address
	NewMinStake *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMinStakeSet is a free log retrieval operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterMinStakeSet(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2MinStakeSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2MinStakeSetIterator{contract: _Vanillaregistryv2.contract, event: "MinStakeSet", logs: logs, sub: sub}, nil
}

// WatchMinStakeSet is a free log subscription operation binding the contract event 0xbd0f06c543aec7980853f7cb191dff311f0ef977570d34683aacc97e33b3f301.
//
// Solidity: event MinStakeSet(address indexed msgSender, uint256 newMinStake)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchMinStakeSet(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2MinStakeSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "MinStakeSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2MinStakeSet)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseMinStakeSet(log types.Log) (*Vanillaregistryv2MinStakeSet, error) {
	event := new(Vanillaregistryv2MinStakeSet)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "MinStakeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2OwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2OwnershipTransferStartedIterator struct {
	Event *Vanillaregistryv2OwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2OwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2OwnershipTransferStarted)
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
		it.Event = new(Vanillaregistryv2OwnershipTransferStarted)
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
func (it *Vanillaregistryv2OwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2OwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2OwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2OwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Vanillaregistryv2OwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2OwnershipTransferStartedIterator{contract: _Vanillaregistryv2.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2OwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2OwnershipTransferStarted)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseOwnershipTransferStarted(log types.Log) (*Vanillaregistryv2OwnershipTransferStarted, error) {
	event := new(Vanillaregistryv2OwnershipTransferStarted)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2OwnershipTransferredIterator struct {
	Event *Vanillaregistryv2OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2OwnershipTransferred)
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
		it.Event = new(Vanillaregistryv2OwnershipTransferred)
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
func (it *Vanillaregistryv2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2OwnershipTransferred represents a OwnershipTransferred event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Vanillaregistryv2OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2OwnershipTransferredIterator{contract: _Vanillaregistryv2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2OwnershipTransferred)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseOwnershipTransferred(log types.Log) (*Vanillaregistryv2OwnershipTransferred, error) {
	event := new(Vanillaregistryv2OwnershipTransferred)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2PausedIterator struct {
	Event *Vanillaregistryv2Paused // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Paused)
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
		it.Event = new(Vanillaregistryv2Paused)
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
func (it *Vanillaregistryv2PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Paused represents a Paused event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterPaused(opts *bind.FilterOpts) (*Vanillaregistryv2PausedIterator, error) {

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2PausedIterator{contract: _Vanillaregistryv2.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Paused) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Paused)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParsePaused(log types.Log) (*Vanillaregistryv2Paused, error) {
	event := new(Vanillaregistryv2Paused)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2SlashOracleSetIterator is returned from FilterSlashOracleSet and is used to iterate over the raw logs and unpacked data for SlashOracleSet events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashOracleSetIterator struct {
	Event *Vanillaregistryv2SlashOracleSet // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2SlashOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2SlashOracleSet)
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
		it.Event = new(Vanillaregistryv2SlashOracleSet)
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
func (it *Vanillaregistryv2SlashOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2SlashOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2SlashOracleSet represents a SlashOracleSet event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashOracleSet struct {
	MsgSender      common.Address
	NewSlashOracle common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSlashOracleSet is a free log retrieval operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterSlashOracleSet(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2SlashOracleSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2SlashOracleSetIterator{contract: _Vanillaregistryv2.contract, event: "SlashOracleSet", logs: logs, sub: sub}, nil
}

// WatchSlashOracleSet is a free log subscription operation binding the contract event 0x5b8cc95f72c2f7fba20ba3e60c77062f56cc5a2f3cba5aeaddee4c51812d27ea.
//
// Solidity: event SlashOracleSet(address indexed msgSender, address newSlashOracle)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchSlashOracleSet(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2SlashOracleSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "SlashOracleSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2SlashOracleSet)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseSlashOracleSet(log types.Log) (*Vanillaregistryv2SlashOracleSet, error) {
	event := new(Vanillaregistryv2SlashOracleSet)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2SlashReceiverSetIterator is returned from FilterSlashReceiverSet and is used to iterate over the raw logs and unpacked data for SlashReceiverSet events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashReceiverSetIterator struct {
	Event *Vanillaregistryv2SlashReceiverSet // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2SlashReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2SlashReceiverSet)
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
		it.Event = new(Vanillaregistryv2SlashReceiverSet)
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
func (it *Vanillaregistryv2SlashReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2SlashReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2SlashReceiverSet represents a SlashReceiverSet event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashReceiverSet struct {
	MsgSender        common.Address
	NewSlashReceiver common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSlashReceiverSet is a free log retrieval operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterSlashReceiverSet(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2SlashReceiverSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2SlashReceiverSetIterator{contract: _Vanillaregistryv2.contract, event: "SlashReceiverSet", logs: logs, sub: sub}, nil
}

// WatchSlashReceiverSet is a free log subscription operation binding the contract event 0xf7f99ea479b331e341a35cdf347f232a35dd611f889867759df261eeb540770a.
//
// Solidity: event SlashReceiverSet(address indexed msgSender, address newSlashReceiver)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchSlashReceiverSet(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2SlashReceiverSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "SlashReceiverSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2SlashReceiverSet)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseSlashReceiverSet(log types.Log) (*Vanillaregistryv2SlashReceiverSet, error) {
	event := new(Vanillaregistryv2SlashReceiverSet)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2SlashedIterator is returned from FilterSlashed and is used to iterate over the raw logs and unpacked data for Slashed events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashedIterator struct {
	Event *Vanillaregistryv2Slashed // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2SlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Slashed)
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
		it.Event = new(Vanillaregistryv2Slashed)
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
func (it *Vanillaregistryv2SlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2SlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Slashed represents a Slashed event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Slashed struct {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterSlashed(opts *bind.FilterOpts, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2SlashedIterator, error) {

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

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2SlashedIterator{contract: _Vanillaregistryv2.contract, event: "Slashed", logs: logs, sub: sub}, nil
}

// WatchSlashed is a free log subscription operation binding the contract event 0xf15b8630ce764d5dbcfaaa9843c3e5fcdb460aaaa46d7dc3ff4f19ca4096fc07.
//
// Solidity: event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchSlashed(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Slashed, msgSender []common.Address, slashReceiver []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Slashed", msgSenderRule, slashReceiverRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Slashed)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Slashed", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseSlashed(log types.Log) (*Vanillaregistryv2Slashed, error) {
	event := new(Vanillaregistryv2Slashed)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Slashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator is returned from FilterSlashingPayoutPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for SlashingPayoutPeriodBlocksSet events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator struct {
	Event *Vanillaregistryv2SlashingPayoutPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2SlashingPayoutPeriodBlocksSet)
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
		it.Event = new(Vanillaregistryv2SlashingPayoutPeriodBlocksSet)
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
func (it *Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2SlashingPayoutPeriodBlocksSet represents a SlashingPayoutPeriodBlocksSet event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2SlashingPayoutPeriodBlocksSet struct {
	MsgSender                     common.Address
	NewSlashingPayoutPeriodBlocks *big.Int
	Raw                           types.Log // Blockchain specific contextual infos
}

// FilterSlashingPayoutPeriodBlocksSet is a free log retrieval operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterSlashingPayoutPeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2SlashingPayoutPeriodBlocksSetIterator{contract: _Vanillaregistryv2.contract, event: "SlashingPayoutPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchSlashingPayoutPeriodBlocksSet is a free log subscription operation binding the contract event 0x537af662b191583a2538de843160914f46fd6033598c645efd5f9ac3cb3f650e.
//
// Solidity: event SlashingPayoutPeriodBlocksSet(address indexed msgSender, uint256 newSlashingPayoutPeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchSlashingPayoutPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2SlashingPayoutPeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "SlashingPayoutPeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2SlashingPayoutPeriodBlocksSet)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseSlashingPayoutPeriodBlocksSet(log types.Log) (*Vanillaregistryv2SlashingPayoutPeriodBlocksSet, error) {
	event := new(Vanillaregistryv2SlashingPayoutPeriodBlocksSet)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "SlashingPayoutPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2StakeAddedIterator is returned from FilterStakeAdded and is used to iterate over the raw logs and unpacked data for StakeAdded events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakeAddedIterator struct {
	Event *Vanillaregistryv2StakeAdded // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2StakeAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2StakeAdded)
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
		it.Event = new(Vanillaregistryv2StakeAdded)
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
func (it *Vanillaregistryv2StakeAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2StakeAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2StakeAdded represents a StakeAdded event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakeAdded struct {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterStakeAdded(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2StakeAddedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2StakeAddedIterator{contract: _Vanillaregistryv2.contract, event: "StakeAdded", logs: logs, sub: sub}, nil
}

// WatchStakeAdded is a free log subscription operation binding the contract event 0xb01516cc7ddda8b10127c714474503b38a75b9afa8a4e4b9da306e61181980c7.
//
// Solidity: event StakeAdded(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount, uint256 newBalance)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchStakeAdded(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2StakeAdded, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "StakeAdded", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2StakeAdded)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakeAdded", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseStakeAdded(log types.Log) (*Vanillaregistryv2StakeAdded, error) {
	event := new(Vanillaregistryv2StakeAdded)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakeAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2StakeWithdrawnIterator is returned from FilterStakeWithdrawn and is used to iterate over the raw logs and unpacked data for StakeWithdrawn events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakeWithdrawnIterator struct {
	Event *Vanillaregistryv2StakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2StakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2StakeWithdrawn)
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
		it.Event = new(Vanillaregistryv2StakeWithdrawn)
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
func (it *Vanillaregistryv2StakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2StakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2StakeWithdrawn represents a StakeWithdrawn event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdrawn is a free log retrieval operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2StakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2StakeWithdrawnIterator{contract: _Vanillaregistryv2.contract, event: "StakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchStakeWithdrawn is a free log subscription operation binding the contract event 0x3ff0f1758b0b95c72d1f781b732306588b99dabb298fec793499eb8803b05465.
//
// Solidity: event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2StakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "StakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2StakeWithdrawn)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseStakeWithdrawn(log types.Log) (*Vanillaregistryv2StakeWithdrawn, error) {
	event := new(Vanillaregistryv2StakeWithdrawn)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2StakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakedIterator struct {
	Event *Vanillaregistryv2Staked // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2StakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Staked)
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
		it.Event = new(Vanillaregistryv2Staked)
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
func (it *Vanillaregistryv2StakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2StakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Staked represents a Staked event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Staked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterStaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2StakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2StakedIterator{contract: _Vanillaregistryv2.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1c9a8e1c32f2ea144885ec1a1398b5d51d627f9532fb2614516322a0b8087de5.
//
// Solidity: event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Staked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Staked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Staked)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Staked", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseStaked(log types.Log) (*Vanillaregistryv2Staked, error) {
	event := new(Vanillaregistryv2Staked)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2StakerRemovedFromWhitelistIterator is returned from FilterStakerRemovedFromWhitelist and is used to iterate over the raw logs and unpacked data for StakerRemovedFromWhitelist events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakerRemovedFromWhitelistIterator struct {
	Event *Vanillaregistryv2StakerRemovedFromWhitelist // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2StakerRemovedFromWhitelistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2StakerRemovedFromWhitelist)
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
		it.Event = new(Vanillaregistryv2StakerRemovedFromWhitelist)
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
func (it *Vanillaregistryv2StakerRemovedFromWhitelistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2StakerRemovedFromWhitelistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2StakerRemovedFromWhitelist represents a StakerRemovedFromWhitelist event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakerRemovedFromWhitelist struct {
	MsgSender common.Address
	Staker    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStakerRemovedFromWhitelist is a free log retrieval operation binding the contract event 0x79fccea9ca325c6c715589de68546a63a35425178d1a8b1436bef7f7a76087a4.
//
// Solidity: event StakerRemovedFromWhitelist(address indexed msgSender, address staker)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterStakerRemovedFromWhitelist(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2StakerRemovedFromWhitelistIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "StakerRemovedFromWhitelist", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2StakerRemovedFromWhitelistIterator{contract: _Vanillaregistryv2.contract, event: "StakerRemovedFromWhitelist", logs: logs, sub: sub}, nil
}

// WatchStakerRemovedFromWhitelist is a free log subscription operation binding the contract event 0x79fccea9ca325c6c715589de68546a63a35425178d1a8b1436bef7f7a76087a4.
//
// Solidity: event StakerRemovedFromWhitelist(address indexed msgSender, address staker)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchStakerRemovedFromWhitelist(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2StakerRemovedFromWhitelist, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "StakerRemovedFromWhitelist", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2StakerRemovedFromWhitelist)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakerRemovedFromWhitelist", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseStakerRemovedFromWhitelist(log types.Log) (*Vanillaregistryv2StakerRemovedFromWhitelist, error) {
	event := new(Vanillaregistryv2StakerRemovedFromWhitelist)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakerRemovedFromWhitelist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2StakerWhitelistedIterator is returned from FilterStakerWhitelisted and is used to iterate over the raw logs and unpacked data for StakerWhitelisted events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakerWhitelistedIterator struct {
	Event *Vanillaregistryv2StakerWhitelisted // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2StakerWhitelistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2StakerWhitelisted)
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
		it.Event = new(Vanillaregistryv2StakerWhitelisted)
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
func (it *Vanillaregistryv2StakerWhitelistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2StakerWhitelistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2StakerWhitelisted represents a StakerWhitelisted event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2StakerWhitelisted struct {
	MsgSender common.Address
	Staker    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStakerWhitelisted is a free log retrieval operation binding the contract event 0xc94be990bcb425c41f84d9928ad9082894fb4f1f1d2508bc088d3a6a3e4059db.
//
// Solidity: event StakerWhitelisted(address indexed msgSender, address staker)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterStakerWhitelisted(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2StakerWhitelistedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "StakerWhitelisted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2StakerWhitelistedIterator{contract: _Vanillaregistryv2.contract, event: "StakerWhitelisted", logs: logs, sub: sub}, nil
}

// WatchStakerWhitelisted is a free log subscription operation binding the contract event 0xc94be990bcb425c41f84d9928ad9082894fb4f1f1d2508bc088d3a6a3e4059db.
//
// Solidity: event StakerWhitelisted(address indexed msgSender, address staker)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchStakerWhitelisted(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2StakerWhitelisted, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "StakerWhitelisted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2StakerWhitelisted)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakerWhitelisted", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseStakerWhitelisted(log types.Log) (*Vanillaregistryv2StakerWhitelisted, error) {
	event := new(Vanillaregistryv2StakerWhitelisted)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "StakerWhitelisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2TotalStakeWithdrawnIterator is returned from FilterTotalStakeWithdrawn and is used to iterate over the raw logs and unpacked data for TotalStakeWithdrawn events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2TotalStakeWithdrawnIterator struct {
	Event *Vanillaregistryv2TotalStakeWithdrawn // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2TotalStakeWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2TotalStakeWithdrawn)
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
		it.Event = new(Vanillaregistryv2TotalStakeWithdrawn)
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
func (it *Vanillaregistryv2TotalStakeWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2TotalStakeWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2TotalStakeWithdrawn represents a TotalStakeWithdrawn event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2TotalStakeWithdrawn struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	TotalAmount       *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTotalStakeWithdrawn is a free log retrieval operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterTotalStakeWithdrawn(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2TotalStakeWithdrawnIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2TotalStakeWithdrawnIterator{contract: _Vanillaregistryv2.contract, event: "TotalStakeWithdrawn", logs: logs, sub: sub}, nil
}

// WatchTotalStakeWithdrawn is a free log subscription operation binding the contract event 0xf0ac877c24c32b3466c3766ad66c170058d5f4ae8347c93bc5adc21b10c14cbe.
//
// Solidity: event TotalStakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, uint256 totalAmount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchTotalStakeWithdrawn(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2TotalStakeWithdrawn, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "TotalStakeWithdrawn", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2TotalStakeWithdrawn)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseTotalStakeWithdrawn(log types.Log) (*Vanillaregistryv2TotalStakeWithdrawn, error) {
	event := new(Vanillaregistryv2TotalStakeWithdrawn)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "TotalStakeWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2UnpausedIterator struct {
	Event *Vanillaregistryv2Unpaused // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Unpaused)
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
		it.Event = new(Vanillaregistryv2Unpaused)
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
func (it *Vanillaregistryv2UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Unpaused represents a Unpaused event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterUnpaused(opts *bind.FilterOpts) (*Vanillaregistryv2UnpausedIterator, error) {

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2UnpausedIterator{contract: _Vanillaregistryv2.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Unpaused) (event.Subscription, error) {

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Unpaused)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseUnpaused(log types.Log) (*Vanillaregistryv2Unpaused, error) {
	event := new(Vanillaregistryv2Unpaused)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2UnstakePeriodBlocksSetIterator is returned from FilterUnstakePeriodBlocksSet and is used to iterate over the raw logs and unpacked data for UnstakePeriodBlocksSet events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2UnstakePeriodBlocksSetIterator struct {
	Event *Vanillaregistryv2UnstakePeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2UnstakePeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2UnstakePeriodBlocksSet)
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
		it.Event = new(Vanillaregistryv2UnstakePeriodBlocksSet)
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
func (it *Vanillaregistryv2UnstakePeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2UnstakePeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2UnstakePeriodBlocksSet represents a UnstakePeriodBlocksSet event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2UnstakePeriodBlocksSet struct {
	MsgSender              common.Address
	NewUnstakePeriodBlocks *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUnstakePeriodBlocksSet is a free log retrieval operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterUnstakePeriodBlocksSet(opts *bind.FilterOpts, msgSender []common.Address) (*Vanillaregistryv2UnstakePeriodBlocksSetIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2UnstakePeriodBlocksSetIterator{contract: _Vanillaregistryv2.contract, event: "UnstakePeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchUnstakePeriodBlocksSet is a free log subscription operation binding the contract event 0x1c7b684565a5bbbb1e7647588e4e6cf72ffa21a25545a4385f2074132aa51613.
//
// Solidity: event UnstakePeriodBlocksSet(address indexed msgSender, uint256 newUnstakePeriodBlocks)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchUnstakePeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2UnstakePeriodBlocksSet, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "UnstakePeriodBlocksSet", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2UnstakePeriodBlocksSet)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseUnstakePeriodBlocksSet(log types.Log) (*Vanillaregistryv2UnstakePeriodBlocksSet, error) {
	event := new(Vanillaregistryv2UnstakePeriodBlocksSet)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "UnstakePeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2UnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2UnstakedIterator struct {
	Event *Vanillaregistryv2Unstaked // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2UnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Unstaked)
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
		it.Event = new(Vanillaregistryv2Unstaked)
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
func (it *Vanillaregistryv2UnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2UnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Unstaked represents a Unstaked event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Unstaked struct {
	MsgSender         common.Address
	WithdrawalAddress common.Address
	ValBLSPubKey      []byte
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterUnstaked(opts *bind.FilterOpts, msgSender []common.Address, withdrawalAddress []common.Address) (*Vanillaregistryv2UnstakedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2UnstakedIterator{contract: _Vanillaregistryv2.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x104975b81462e4e991f38b5b4158bf402b1528cd36ac80b123aef9d06dd0e1a9.
//
// Solidity: event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Unstaked, msgSender []common.Address, withdrawalAddress []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var withdrawalAddressRule []interface{}
	for _, withdrawalAddressItem := range withdrawalAddress {
		withdrawalAddressRule = append(withdrawalAddressRule, withdrawalAddressItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Unstaked", msgSenderRule, withdrawalAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Unstaked)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Unstaked", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseUnstaked(log types.Log) (*Vanillaregistryv2Unstaked, error) {
	event := new(Vanillaregistryv2Unstaked)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Vanillaregistryv2UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2UpgradedIterator struct {
	Event *Vanillaregistryv2Upgraded // Event containing the contract specifics and raw log

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
func (it *Vanillaregistryv2UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Vanillaregistryv2Upgraded)
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
		it.Event = new(Vanillaregistryv2Upgraded)
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
func (it *Vanillaregistryv2UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Vanillaregistryv2UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Vanillaregistryv2Upgraded represents a Upgraded event raised by the Vanillaregistryv2 contract.
type Vanillaregistryv2Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*Vanillaregistryv2UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &Vanillaregistryv2UpgradedIterator{contract: _Vanillaregistryv2.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *Vanillaregistryv2Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Vanillaregistryv2.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Vanillaregistryv2Upgraded)
				if err := _Vanillaregistryv2.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Vanillaregistryv2 *Vanillaregistryv2Filterer) ParseUpgraded(log types.Log) (*Vanillaregistryv2Upgraded, error) {
	event := new(Vanillaregistryv2Upgraded)
	if err := _Vanillaregistryv2.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
