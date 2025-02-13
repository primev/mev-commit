// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package providerregistry

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

// ProviderregistryMetaData contains all meta data concerning the Providerregistry contract.
var ProviderregistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"ONE_HUNDRED_PERCENT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addVerifiedBLSKey\",\"inputs\":[{\"name\":\"blsPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidderSlashedAmount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockBuilderBLSKeyToAddress\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegateRegisterAndStake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"delegateStake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"eoaToBlsPubkeys\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feePercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAccumulatedPenaltyFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBLSKeys\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEoaFromBLSKey\",\"inputs\":[{\"name\":\"blsKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProviderStake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_penaltyFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feePercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_withdrawalDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_penaltyFeePayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isProviderValid\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"manuallyWithdrawPenaltyFee\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"overrideAddBLSKey\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blsPublicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"penaltyFeeTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accumulatedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastPayoutBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"preconfManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerStakes\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerAndStake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeePayoutPeriodBlocks\",\"inputs\":[{\"name\":\"_feePayoutPeriodBlocks\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"_minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewFeePercent\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNewPenaltyFeeRecipient\",\"inputs\":[{\"name\":\"newFeeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPreconfManager\",\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWithdrawalDelay\",\"inputs\":[{\"name\":\"_withdrawalDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slash\",\"inputs\":[{\"name\":\"amt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"verifySignature\",\"inputs\":[{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"message\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawSlashedAmount\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawalDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawalRequests\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BLSKeyAdded\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"blsPublicKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BidderWithdrawSlashedAmount\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePayoutPeriodBlocksUpdated\",\"inputs\":[{\"name\":\"newFeePayoutPeriodBlocks\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePercentUpdated\",\"inputs\":[{\"name\":\"newFeePercent\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeTransfer\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsDeposited\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsSlashed\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsufficientFundsToSlash\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"providerStake\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"residualAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"penaltyFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeUpdated\",\"inputs\":[{\"name\":\"newMinStake\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PenaltyFeeRecipientUpdated\",\"inputs\":[{\"name\":\"newPenaltyFeeRecipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PreconfManagerUpdated\",\"inputs\":[{\"name\":\"newPreconfManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderRegistered\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"stakedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferToBidderFailed\",\"inputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unstake\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalDelayUpdated\",\"inputs\":[{\"name\":\"newWithdrawalDelay\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AtLeastOneBLSKeyRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BLSSignatureInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BidderAmountIsZero\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BidderWithdrawalTransferFailed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"DelayNotPassed\",\"inputs\":[{\"name\":\"withdrawalRequestTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawalDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"currentBlockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeRecipientIsZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientStake\",\"inputs\":[{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minStake\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBLSPublicKeyLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expectedLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoStakeToWithdraw\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoUnstakeRequest\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPreconfContract\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"preconfManager\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PayoutPeriodMustBePositive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingWithdrawalRequest\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PreconfManagerNotSet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ProviderAlreadyRegistered\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProviderCommitmentsPending\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"numPending\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ProviderNotRegistered\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PublicKeyLengthInvalid\",\"inputs\":[{\"name\":\"exp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SignatureLengthInvalid\",\"inputs\":[{\"name\":\"exp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"got\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StakeTransferFailed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TransferToRecipientFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnstakeRequestExists\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// ProviderregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ProviderregistryMetaData.ABI instead.
var ProviderregistryABI = ProviderregistryMetaData.ABI

// Providerregistry is an auto generated Go binding around an Ethereum contract.
type Providerregistry struct {
	ProviderregistryCaller     // Read-only binding to the contract
	ProviderregistryTransactor // Write-only binding to the contract
	ProviderregistryFilterer   // Log filterer for contract events
}

// ProviderregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProviderregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProviderregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProviderregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProviderregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProviderregistrySession struct {
	Contract     *Providerregistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProviderregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProviderregistryCallerSession struct {
	Contract *ProviderregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ProviderregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProviderregistryTransactorSession struct {
	Contract     *ProviderregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ProviderregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProviderregistryRaw struct {
	Contract *Providerregistry // Generic contract binding to access the raw methods on
}

// ProviderregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProviderregistryCallerRaw struct {
	Contract *ProviderregistryCaller // Generic read-only contract binding to access the raw methods on
}

// ProviderregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProviderregistryTransactorRaw struct {
	Contract *ProviderregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProviderregistry creates a new instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistry(address common.Address, backend bind.ContractBackend) (*Providerregistry, error) {
	contract, err := bindProviderregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Providerregistry{ProviderregistryCaller: ProviderregistryCaller{contract: contract}, ProviderregistryTransactor: ProviderregistryTransactor{contract: contract}, ProviderregistryFilterer: ProviderregistryFilterer{contract: contract}}, nil
}

// NewProviderregistryCaller creates a new read-only instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryCaller(address common.Address, caller bind.ContractCaller) (*ProviderregistryCaller, error) {
	contract, err := bindProviderregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryCaller{contract: contract}, nil
}

// NewProviderregistryTransactor creates a new write-only instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ProviderregistryTransactor, error) {
	contract, err := bindProviderregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryTransactor{contract: contract}, nil
}

// NewProviderregistryFilterer creates a new log filterer instance of Providerregistry, bound to a specific deployed contract.
func NewProviderregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ProviderregistryFilterer, error) {
	contract, err := bindProviderregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFilterer{contract: contract}, nil
}

// bindProviderregistry binds a generic wrapper to an already deployed contract.
func bindProviderregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProviderregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Providerregistry *ProviderregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Providerregistry.Contract.ProviderregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Providerregistry *ProviderregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.Contract.ProviderregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Providerregistry *ProviderregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Providerregistry.Contract.ProviderregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Providerregistry *ProviderregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Providerregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Providerregistry *ProviderregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Providerregistry *ProviderregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Providerregistry.Contract.contract.Transact(opts, method, params...)
}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) ONEHUNDREDPERCENT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "ONE_HUNDRED_PERCENT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) ONEHUNDREDPERCENT() (*big.Int, error) {
	return _Providerregistry.Contract.ONEHUNDREDPERCENT(&_Providerregistry.CallOpts)
}

// ONEHUNDREDPERCENT is a free data retrieval call binding the contract method 0xdd0081c7.
//
// Solidity: function ONE_HUNDRED_PERCENT() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) ONEHUNDREDPERCENT() (*big.Int, error) {
	return _Providerregistry.Contract.ONEHUNDREDPERCENT(&_Providerregistry.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) PRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) PRECISION() (*big.Int, error) {
	return _Providerregistry.Contract.PRECISION(&_Providerregistry.CallOpts)
}

// PRECISION is a free data retrieval call binding the contract method 0xaaf5eb68.
//
// Solidity: function PRECISION() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) PRECISION() (*big.Int, error) {
	return _Providerregistry.Contract.PRECISION(&_Providerregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Providerregistry *ProviderregistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Providerregistry *ProviderregistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Providerregistry.Contract.UPGRADEINTERFACEVERSION(&_Providerregistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Providerregistry *ProviderregistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Providerregistry.Contract.UPGRADEINTERFACEVERSION(&_Providerregistry.CallOpts)
}

// BidderSlashedAmount is a free data retrieval call binding the contract method 0x3ab6fc1a.
//
// Solidity: function bidderSlashedAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) BidderSlashedAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "bidderSlashedAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BidderSlashedAmount is a free data retrieval call binding the contract method 0x3ab6fc1a.
//
// Solidity: function bidderSlashedAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) BidderSlashedAmount(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.BidderSlashedAmount(&_Providerregistry.CallOpts, arg0)
}

// BidderSlashedAmount is a free data retrieval call binding the contract method 0x3ab6fc1a.
//
// Solidity: function bidderSlashedAmount(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) BidderSlashedAmount(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.BidderSlashedAmount(&_Providerregistry.CallOpts, arg0)
}

// BlockBuilderBLSKeyToAddress is a free data retrieval call binding the contract method 0x929b63c4.
//
// Solidity: function blockBuilderBLSKeyToAddress(bytes ) view returns(address)
func (_Providerregistry *ProviderregistryCaller) BlockBuilderBLSKeyToAddress(opts *bind.CallOpts, arg0 []byte) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "blockBuilderBLSKeyToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockBuilderBLSKeyToAddress is a free data retrieval call binding the contract method 0x929b63c4.
//
// Solidity: function blockBuilderBLSKeyToAddress(bytes ) view returns(address)
func (_Providerregistry *ProviderregistrySession) BlockBuilderBLSKeyToAddress(arg0 []byte) (common.Address, error) {
	return _Providerregistry.Contract.BlockBuilderBLSKeyToAddress(&_Providerregistry.CallOpts, arg0)
}

// BlockBuilderBLSKeyToAddress is a free data retrieval call binding the contract method 0x929b63c4.
//
// Solidity: function blockBuilderBLSKeyToAddress(bytes ) view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) BlockBuilderBLSKeyToAddress(arg0 []byte) (common.Address, error) {
	return _Providerregistry.Contract.BlockBuilderBLSKeyToAddress(&_Providerregistry.CallOpts, arg0)
}

// EoaToBlsPubkeys is a free data retrieval call binding the contract method 0x1129ce1f.
//
// Solidity: function eoaToBlsPubkeys(address , uint256 ) view returns(bytes)
func (_Providerregistry *ProviderregistryCaller) EoaToBlsPubkeys(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "eoaToBlsPubkeys", arg0, arg1)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EoaToBlsPubkeys is a free data retrieval call binding the contract method 0x1129ce1f.
//
// Solidity: function eoaToBlsPubkeys(address , uint256 ) view returns(bytes)
func (_Providerregistry *ProviderregistrySession) EoaToBlsPubkeys(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return _Providerregistry.Contract.EoaToBlsPubkeys(&_Providerregistry.CallOpts, arg0, arg1)
}

// EoaToBlsPubkeys is a free data retrieval call binding the contract method 0x1129ce1f.
//
// Solidity: function eoaToBlsPubkeys(address , uint256 ) view returns(bytes)
func (_Providerregistry *ProviderregistryCallerSession) EoaToBlsPubkeys(arg0 common.Address, arg1 *big.Int) ([]byte, error) {
	return _Providerregistry.Contract.EoaToBlsPubkeys(&_Providerregistry.CallOpts, arg0, arg1)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) FeePercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "feePercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) FeePercent() (*big.Int, error) {
	return _Providerregistry.Contract.FeePercent(&_Providerregistry.CallOpts)
}

// FeePercent is a free data retrieval call binding the contract method 0x7fd6f15c.
//
// Solidity: function feePercent() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) FeePercent() (*big.Int, error) {
	return _Providerregistry.Contract.FeePercent(&_Providerregistry.CallOpts)
}

// GetAccumulatedPenaltyFee is a free data retrieval call binding the contract method 0xe4506e7c.
//
// Solidity: function getAccumulatedPenaltyFee() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) GetAccumulatedPenaltyFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "getAccumulatedPenaltyFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAccumulatedPenaltyFee is a free data retrieval call binding the contract method 0xe4506e7c.
//
// Solidity: function getAccumulatedPenaltyFee() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) GetAccumulatedPenaltyFee() (*big.Int, error) {
	return _Providerregistry.Contract.GetAccumulatedPenaltyFee(&_Providerregistry.CallOpts)
}

// GetAccumulatedPenaltyFee is a free data retrieval call binding the contract method 0xe4506e7c.
//
// Solidity: function getAccumulatedPenaltyFee() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) GetAccumulatedPenaltyFee() (*big.Int, error) {
	return _Providerregistry.Contract.GetAccumulatedPenaltyFee(&_Providerregistry.CallOpts)
}

// GetBLSKeys is a free data retrieval call binding the contract method 0xc50b59df.
//
// Solidity: function getBLSKeys(address provider) view returns(bytes[])
func (_Providerregistry *ProviderregistryCaller) GetBLSKeys(opts *bind.CallOpts, provider common.Address) ([][]byte, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "getBLSKeys", provider)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// GetBLSKeys is a free data retrieval call binding the contract method 0xc50b59df.
//
// Solidity: function getBLSKeys(address provider) view returns(bytes[])
func (_Providerregistry *ProviderregistrySession) GetBLSKeys(provider common.Address) ([][]byte, error) {
	return _Providerregistry.Contract.GetBLSKeys(&_Providerregistry.CallOpts, provider)
}

// GetBLSKeys is a free data retrieval call binding the contract method 0xc50b59df.
//
// Solidity: function getBLSKeys(address provider) view returns(bytes[])
func (_Providerregistry *ProviderregistryCallerSession) GetBLSKeys(provider common.Address) ([][]byte, error) {
	return _Providerregistry.Contract.GetBLSKeys(&_Providerregistry.CallOpts, provider)
}

// GetEoaFromBLSKey is a free data retrieval call binding the contract method 0xea3b275d.
//
// Solidity: function getEoaFromBLSKey(bytes blsKey) view returns(address)
func (_Providerregistry *ProviderregistryCaller) GetEoaFromBLSKey(opts *bind.CallOpts, blsKey []byte) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "getEoaFromBLSKey", blsKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetEoaFromBLSKey is a free data retrieval call binding the contract method 0xea3b275d.
//
// Solidity: function getEoaFromBLSKey(bytes blsKey) view returns(address)
func (_Providerregistry *ProviderregistrySession) GetEoaFromBLSKey(blsKey []byte) (common.Address, error) {
	return _Providerregistry.Contract.GetEoaFromBLSKey(&_Providerregistry.CallOpts, blsKey)
}

// GetEoaFromBLSKey is a free data retrieval call binding the contract method 0xea3b275d.
//
// Solidity: function getEoaFromBLSKey(bytes blsKey) view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) GetEoaFromBLSKey(blsKey []byte) (common.Address, error) {
	return _Providerregistry.Contract.GetEoaFromBLSKey(&_Providerregistry.CallOpts, blsKey)
}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) GetProviderStake(opts *bind.CallOpts, provider common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "getProviderStake", provider)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) GetProviderStake(provider common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.GetProviderStake(&_Providerregistry.CallOpts, provider)
}

// GetProviderStake is a free data retrieval call binding the contract method 0xbfebc370.
//
// Solidity: function getProviderStake(address provider) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) GetProviderStake(provider common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.GetProviderStake(&_Providerregistry.CallOpts, provider)
}

// IsProviderValid is a free data retrieval call binding the contract method 0xb066d50d.
//
// Solidity: function isProviderValid(address provider) view returns()
func (_Providerregistry *ProviderregistryCaller) IsProviderValid(opts *bind.CallOpts, provider common.Address) error {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "isProviderValid", provider)

	if err != nil {
		return err
	}

	return err

}

// IsProviderValid is a free data retrieval call binding the contract method 0xb066d50d.
//
// Solidity: function isProviderValid(address provider) view returns()
func (_Providerregistry *ProviderregistrySession) IsProviderValid(provider common.Address) error {
	return _Providerregistry.Contract.IsProviderValid(&_Providerregistry.CallOpts, provider)
}

// IsProviderValid is a free data retrieval call binding the contract method 0xb066d50d.
//
// Solidity: function isProviderValid(address provider) view returns()
func (_Providerregistry *ProviderregistryCallerSession) IsProviderValid(provider common.Address) error {
	return _Providerregistry.Contract.IsProviderValid(&_Providerregistry.CallOpts, provider)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "minStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) MinStake() (*big.Int, error) {
	return _Providerregistry.Contract.MinStake(&_Providerregistry.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) MinStake() (*big.Int, error) {
	return _Providerregistry.Contract.MinStake(&_Providerregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistrySession) Owner() (common.Address, error) {
	return _Providerregistry.Contract.Owner(&_Providerregistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) Owner() (common.Address, error) {
	return _Providerregistry.Contract.Owner(&_Providerregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Providerregistry *ProviderregistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Providerregistry *ProviderregistrySession) Paused() (bool, error) {
	return _Providerregistry.Contract.Paused(&_Providerregistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Providerregistry *ProviderregistryCallerSession) Paused() (bool, error) {
	return _Providerregistry.Contract.Paused(&_Providerregistry.CallOpts)
}

// PenaltyFeeTracker is a free data retrieval call binding the contract method 0xf1aab9ee.
//
// Solidity: function penaltyFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Providerregistry *ProviderregistryCaller) PenaltyFeeTracker(opts *bind.CallOpts) (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "penaltyFeeTracker")

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

// PenaltyFeeTracker is a free data retrieval call binding the contract method 0xf1aab9ee.
//
// Solidity: function penaltyFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Providerregistry *ProviderregistrySession) PenaltyFeeTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Providerregistry.Contract.PenaltyFeeTracker(&_Providerregistry.CallOpts)
}

// PenaltyFeeTracker is a free data retrieval call binding the contract method 0xf1aab9ee.
//
// Solidity: function penaltyFeeTracker() view returns(address recipient, uint256 accumulatedAmount, uint256 lastPayoutBlock, uint256 payoutPeriodBlocks)
func (_Providerregistry *ProviderregistryCallerSession) PenaltyFeeTracker() (struct {
	Recipient          common.Address
	AccumulatedAmount  *big.Int
	LastPayoutBlock    *big.Int
	PayoutPeriodBlocks *big.Int
}, error) {
	return _Providerregistry.Contract.PenaltyFeeTracker(&_Providerregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Providerregistry *ProviderregistryCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Providerregistry *ProviderregistrySession) PendingOwner() (common.Address, error) {
	return _Providerregistry.Contract.PendingOwner(&_Providerregistry.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) PendingOwner() (common.Address, error) {
	return _Providerregistry.Contract.PendingOwner(&_Providerregistry.CallOpts)
}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Providerregistry *ProviderregistryCaller) PreconfManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "preconfManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Providerregistry *ProviderregistrySession) PreconfManager() (common.Address, error) {
	return _Providerregistry.Contract.PreconfManager(&_Providerregistry.CallOpts)
}

// PreconfManager is a free data retrieval call binding the contract method 0x94a87500.
//
// Solidity: function preconfManager() view returns(address)
func (_Providerregistry *ProviderregistryCallerSession) PreconfManager() (common.Address, error) {
	return _Providerregistry.Contract.PreconfManager(&_Providerregistry.CallOpts)
}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistryCaller) ProviderRegistered(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "providerRegistered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistrySession) ProviderRegistered(arg0 common.Address) (bool, error) {
	return _Providerregistry.Contract.ProviderRegistered(&_Providerregistry.CallOpts, arg0)
}

// ProviderRegistered is a free data retrieval call binding the contract method 0xab255b41.
//
// Solidity: function providerRegistered(address ) view returns(bool)
func (_Providerregistry *ProviderregistryCallerSession) ProviderRegistered(arg0 common.Address) (bool, error) {
	return _Providerregistry.Contract.ProviderRegistered(&_Providerregistry.CallOpts, arg0)
}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) ProviderStakes(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "providerStakes", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) ProviderStakes(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.ProviderStakes(&_Providerregistry.CallOpts, arg0)
}

// ProviderStakes is a free data retrieval call binding the contract method 0x0d6b4c9f.
//
// Solidity: function providerStakes(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) ProviderStakes(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.ProviderStakes(&_Providerregistry.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Providerregistry *ProviderregistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Providerregistry *ProviderregistrySession) ProxiableUUID() ([32]byte, error) {
	return _Providerregistry.Contract.ProxiableUUID(&_Providerregistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Providerregistry *ProviderregistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Providerregistry.Contract.ProxiableUUID(&_Providerregistry.CallOpts)
}

// VerifySignature is a free data retrieval call binding the contract method 0x2222e36f.
//
// Solidity: function verifySignature(bytes pubKey, bytes32 message, bytes signature) view returns(bool)
func (_Providerregistry *ProviderregistryCaller) VerifySignature(opts *bind.CallOpts, pubKey []byte, message [32]byte, signature []byte) (bool, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "verifySignature", pubKey, message, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifySignature is a free data retrieval call binding the contract method 0x2222e36f.
//
// Solidity: function verifySignature(bytes pubKey, bytes32 message, bytes signature) view returns(bool)
func (_Providerregistry *ProviderregistrySession) VerifySignature(pubKey []byte, message [32]byte, signature []byte) (bool, error) {
	return _Providerregistry.Contract.VerifySignature(&_Providerregistry.CallOpts, pubKey, message, signature)
}

// VerifySignature is a free data retrieval call binding the contract method 0x2222e36f.
//
// Solidity: function verifySignature(bytes pubKey, bytes32 message, bytes signature) view returns(bool)
func (_Providerregistry *ProviderregistryCallerSession) VerifySignature(pubKey []byte, message [32]byte, signature []byte) (bool, error) {
	return _Providerregistry.Contract.VerifySignature(&_Providerregistry.CallOpts, pubKey, message, signature)
}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) WithdrawalDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "withdrawalDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_Providerregistry *ProviderregistrySession) WithdrawalDelay() (*big.Int, error) {
	return _Providerregistry.Contract.WithdrawalDelay(&_Providerregistry.CallOpts)
}

// WithdrawalDelay is a free data retrieval call binding the contract method 0xa7ab6961.
//
// Solidity: function withdrawalDelay() view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) WithdrawalDelay() (*big.Int, error) {
	return _Providerregistry.Contract.WithdrawalDelay(&_Providerregistry.CallOpts)
}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x27b380f3.
//
// Solidity: function withdrawalRequests(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCaller) WithdrawalRequests(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Providerregistry.contract.Call(opts, &out, "withdrawalRequests", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x27b380f3.
//
// Solidity: function withdrawalRequests(address ) view returns(uint256)
func (_Providerregistry *ProviderregistrySession) WithdrawalRequests(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.WithdrawalRequests(&_Providerregistry.CallOpts, arg0)
}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x27b380f3.
//
// Solidity: function withdrawalRequests(address ) view returns(uint256)
func (_Providerregistry *ProviderregistryCallerSession) WithdrawalRequests(arg0 common.Address) (*big.Int, error) {
	return _Providerregistry.Contract.WithdrawalRequests(&_Providerregistry.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Providerregistry *ProviderregistryTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Providerregistry *ProviderregistrySession) AcceptOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.AcceptOwnership(&_Providerregistry.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Providerregistry *ProviderregistryTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.AcceptOwnership(&_Providerregistry.TransactOpts)
}

// AddVerifiedBLSKey is a paid mutator transaction binding the contract method 0x7fd358e4.
//
// Solidity: function addVerifiedBLSKey(bytes blsPublicKey, bytes signature) returns()
func (_Providerregistry *ProviderregistryTransactor) AddVerifiedBLSKey(opts *bind.TransactOpts, blsPublicKey []byte, signature []byte) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "addVerifiedBLSKey", blsPublicKey, signature)
}

// AddVerifiedBLSKey is a paid mutator transaction binding the contract method 0x7fd358e4.
//
// Solidity: function addVerifiedBLSKey(bytes blsPublicKey, bytes signature) returns()
func (_Providerregistry *ProviderregistrySession) AddVerifiedBLSKey(blsPublicKey []byte, signature []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.AddVerifiedBLSKey(&_Providerregistry.TransactOpts, blsPublicKey, signature)
}

// AddVerifiedBLSKey is a paid mutator transaction binding the contract method 0x7fd358e4.
//
// Solidity: function addVerifiedBLSKey(bytes blsPublicKey, bytes signature) returns()
func (_Providerregistry *ProviderregistryTransactorSession) AddVerifiedBLSKey(blsPublicKey []byte, signature []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.AddVerifiedBLSKey(&_Providerregistry.TransactOpts, blsPublicKey, signature)
}

// DelegateRegisterAndStake is a paid mutator transaction binding the contract method 0x3d75c0ba.
//
// Solidity: function delegateRegisterAndStake(address provider) payable returns()
func (_Providerregistry *ProviderregistryTransactor) DelegateRegisterAndStake(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "delegateRegisterAndStake", provider)
}

// DelegateRegisterAndStake is a paid mutator transaction binding the contract method 0x3d75c0ba.
//
// Solidity: function delegateRegisterAndStake(address provider) payable returns()
func (_Providerregistry *ProviderregistrySession) DelegateRegisterAndStake(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.DelegateRegisterAndStake(&_Providerregistry.TransactOpts, provider)
}

// DelegateRegisterAndStake is a paid mutator transaction binding the contract method 0x3d75c0ba.
//
// Solidity: function delegateRegisterAndStake(address provider) payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) DelegateRegisterAndStake(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.DelegateRegisterAndStake(&_Providerregistry.TransactOpts, provider)
}

// DelegateStake is a paid mutator transaction binding the contract method 0xf094cc39.
//
// Solidity: function delegateStake(address provider) payable returns()
func (_Providerregistry *ProviderregistryTransactor) DelegateStake(opts *bind.TransactOpts, provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "delegateStake", provider)
}

// DelegateStake is a paid mutator transaction binding the contract method 0xf094cc39.
//
// Solidity: function delegateStake(address provider) payable returns()
func (_Providerregistry *ProviderregistrySession) DelegateStake(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.DelegateStake(&_Providerregistry.TransactOpts, provider)
}

// DelegateStake is a paid mutator transaction binding the contract method 0xf094cc39.
//
// Solidity: function delegateStake(address provider) payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) DelegateStake(provider common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.DelegateStake(&_Providerregistry.TransactOpts, provider)
}

// Initialize is a paid mutator transaction binding the contract method 0xecfa5f52.
//
// Solidity: function initialize(uint256 _minStake, address _penaltyFeeRecipient, uint256 _feePercent, address _owner, uint256 _withdrawalDelay, uint256 _penaltyFeePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistryTransactor) Initialize(opts *bind.TransactOpts, _minStake *big.Int, _penaltyFeeRecipient common.Address, _feePercent *big.Int, _owner common.Address, _withdrawalDelay *big.Int, _penaltyFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "initialize", _minStake, _penaltyFeeRecipient, _feePercent, _owner, _withdrawalDelay, _penaltyFeePayoutPeriodBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0xecfa5f52.
//
// Solidity: function initialize(uint256 _minStake, address _penaltyFeeRecipient, uint256 _feePercent, address _owner, uint256 _withdrawalDelay, uint256 _penaltyFeePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistrySession) Initialize(_minStake *big.Int, _penaltyFeeRecipient common.Address, _feePercent *big.Int, _owner common.Address, _withdrawalDelay *big.Int, _penaltyFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Initialize(&_Providerregistry.TransactOpts, _minStake, _penaltyFeeRecipient, _feePercent, _owner, _withdrawalDelay, _penaltyFeePayoutPeriodBlocks)
}

// Initialize is a paid mutator transaction binding the contract method 0xecfa5f52.
//
// Solidity: function initialize(uint256 _minStake, address _penaltyFeeRecipient, uint256 _feePercent, address _owner, uint256 _withdrawalDelay, uint256 _penaltyFeePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistryTransactorSession) Initialize(_minStake *big.Int, _penaltyFeeRecipient common.Address, _feePercent *big.Int, _owner common.Address, _withdrawalDelay *big.Int, _penaltyFeePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Initialize(&_Providerregistry.TransactOpts, _minStake, _penaltyFeeRecipient, _feePercent, _owner, _withdrawalDelay, _penaltyFeePayoutPeriodBlocks)
}

// ManuallyWithdrawPenaltyFee is a paid mutator transaction binding the contract method 0x7df61dc1.
//
// Solidity: function manuallyWithdrawPenaltyFee() returns()
func (_Providerregistry *ProviderregistryTransactor) ManuallyWithdrawPenaltyFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "manuallyWithdrawPenaltyFee")
}

// ManuallyWithdrawPenaltyFee is a paid mutator transaction binding the contract method 0x7df61dc1.
//
// Solidity: function manuallyWithdrawPenaltyFee() returns()
func (_Providerregistry *ProviderregistrySession) ManuallyWithdrawPenaltyFee() (*types.Transaction, error) {
	return _Providerregistry.Contract.ManuallyWithdrawPenaltyFee(&_Providerregistry.TransactOpts)
}

// ManuallyWithdrawPenaltyFee is a paid mutator transaction binding the contract method 0x7df61dc1.
//
// Solidity: function manuallyWithdrawPenaltyFee() returns()
func (_Providerregistry *ProviderregistryTransactorSession) ManuallyWithdrawPenaltyFee() (*types.Transaction, error) {
	return _Providerregistry.Contract.ManuallyWithdrawPenaltyFee(&_Providerregistry.TransactOpts)
}

// OverrideAddBLSKey is a paid mutator transaction binding the contract method 0xed5219de.
//
// Solidity: function overrideAddBLSKey(address provider, bytes blsPublicKey) returns()
func (_Providerregistry *ProviderregistryTransactor) OverrideAddBLSKey(opts *bind.TransactOpts, provider common.Address, blsPublicKey []byte) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "overrideAddBLSKey", provider, blsPublicKey)
}

// OverrideAddBLSKey is a paid mutator transaction binding the contract method 0xed5219de.
//
// Solidity: function overrideAddBLSKey(address provider, bytes blsPublicKey) returns()
func (_Providerregistry *ProviderregistrySession) OverrideAddBLSKey(provider common.Address, blsPublicKey []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.OverrideAddBLSKey(&_Providerregistry.TransactOpts, provider, blsPublicKey)
}

// OverrideAddBLSKey is a paid mutator transaction binding the contract method 0xed5219de.
//
// Solidity: function overrideAddBLSKey(address provider, bytes blsPublicKey) returns()
func (_Providerregistry *ProviderregistryTransactorSession) OverrideAddBLSKey(provider common.Address, blsPublicKey []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.OverrideAddBLSKey(&_Providerregistry.TransactOpts, provider, blsPublicKey)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Providerregistry *ProviderregistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Providerregistry *ProviderregistrySession) Pause() (*types.Transaction, error) {
	return _Providerregistry.Contract.Pause(&_Providerregistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Providerregistry *ProviderregistryTransactorSession) Pause() (*types.Transaction, error) {
	return _Providerregistry.Contract.Pause(&_Providerregistry.TransactOpts)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistryTransactor) RegisterAndStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "registerAndStake")
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistrySession) RegisterAndStake() (*types.Transaction, error) {
	return _Providerregistry.Contract.RegisterAndStake(&_Providerregistry.TransactOpts)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x84d180ee.
//
// Solidity: function registerAndStake() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) RegisterAndStake() (*types.Transaction, error) {
	return _Providerregistry.Contract.RegisterAndStake(&_Providerregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.RenounceOwnership(&_Providerregistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Providerregistry *ProviderregistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Providerregistry.Contract.RenounceOwnership(&_Providerregistry.TransactOpts)
}

// SetFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0x7cbf9f6e.
//
// Solidity: function setFeePayoutPeriodBlocks(uint256 _feePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistryTransactor) SetFeePayoutPeriodBlocks(opts *bind.TransactOpts, _feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setFeePayoutPeriodBlocks", _feePayoutPeriodBlocks)
}

// SetFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0x7cbf9f6e.
//
// Solidity: function setFeePayoutPeriodBlocks(uint256 _feePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistrySession) SetFeePayoutPeriodBlocks(_feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetFeePayoutPeriodBlocks(&_Providerregistry.TransactOpts, _feePayoutPeriodBlocks)
}

// SetFeePayoutPeriodBlocks is a paid mutator transaction binding the contract method 0x7cbf9f6e.
//
// Solidity: function setFeePayoutPeriodBlocks(uint256 _feePayoutPeriodBlocks) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetFeePayoutPeriodBlocks(_feePayoutPeriodBlocks *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetFeePayoutPeriodBlocks(&_Providerregistry.TransactOpts, _feePayoutPeriodBlocks)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 _minStake) returns()
func (_Providerregistry *ProviderregistryTransactor) SetMinStake(opts *bind.TransactOpts, _minStake *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setMinStake", _minStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 _minStake) returns()
func (_Providerregistry *ProviderregistrySession) SetMinStake(_minStake *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetMinStake(&_Providerregistry.TransactOpts, _minStake)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x8c80fd90.
//
// Solidity: function setMinStake(uint256 _minStake) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetMinStake(_minStake *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetMinStake(&_Providerregistry.TransactOpts, _minStake)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0x3221f6db.
//
// Solidity: function setNewFeePercent(uint256 newFeePercent) returns()
func (_Providerregistry *ProviderregistryTransactor) SetNewFeePercent(opts *bind.TransactOpts, newFeePercent *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setNewFeePercent", newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0x3221f6db.
//
// Solidity: function setNewFeePercent(uint256 newFeePercent) returns()
func (_Providerregistry *ProviderregistrySession) SetNewFeePercent(newFeePercent *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeePercent(&_Providerregistry.TransactOpts, newFeePercent)
}

// SetNewFeePercent is a paid mutator transaction binding the contract method 0x3221f6db.
//
// Solidity: function setNewFeePercent(uint256 newFeePercent) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetNewFeePercent(newFeePercent *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewFeePercent(&_Providerregistry.TransactOpts, newFeePercent)
}

// SetNewPenaltyFeeRecipient is a paid mutator transaction binding the contract method 0xc7cc6f80.
//
// Solidity: function setNewPenaltyFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistryTransactor) SetNewPenaltyFeeRecipient(opts *bind.TransactOpts, newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setNewPenaltyFeeRecipient", newFeeRecipient)
}

// SetNewPenaltyFeeRecipient is a paid mutator transaction binding the contract method 0xc7cc6f80.
//
// Solidity: function setNewPenaltyFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistrySession) SetNewPenaltyFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewPenaltyFeeRecipient(&_Providerregistry.TransactOpts, newFeeRecipient)
}

// SetNewPenaltyFeeRecipient is a paid mutator transaction binding the contract method 0xc7cc6f80.
//
// Solidity: function setNewPenaltyFeeRecipient(address newFeeRecipient) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetNewPenaltyFeeRecipient(newFeeRecipient common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetNewPenaltyFeeRecipient(&_Providerregistry.TransactOpts, newFeeRecipient)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Providerregistry *ProviderregistryTransactor) SetPreconfManager(opts *bind.TransactOpts, contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setPreconfManager", contractAddress)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Providerregistry *ProviderregistrySession) SetPreconfManager(contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetPreconfManager(&_Providerregistry.TransactOpts, contractAddress)
}

// SetPreconfManager is a paid mutator transaction binding the contract method 0x3b79297c.
//
// Solidity: function setPreconfManager(address contractAddress) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetPreconfManager(contractAddress common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetPreconfManager(&_Providerregistry.TransactOpts, contractAddress)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _withdrawalDelay) returns()
func (_Providerregistry *ProviderregistryTransactor) SetWithdrawalDelay(opts *bind.TransactOpts, _withdrawalDelay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "setWithdrawalDelay", _withdrawalDelay)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _withdrawalDelay) returns()
func (_Providerregistry *ProviderregistrySession) SetWithdrawalDelay(_withdrawalDelay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetWithdrawalDelay(&_Providerregistry.TransactOpts, _withdrawalDelay)
}

// SetWithdrawalDelay is a paid mutator transaction binding the contract method 0xd2c13da5.
//
// Solidity: function setWithdrawalDelay(uint256 _withdrawalDelay) returns()
func (_Providerregistry *ProviderregistryTransactorSession) SetWithdrawalDelay(_withdrawalDelay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.SetWithdrawalDelay(&_Providerregistry.TransactOpts, _withdrawalDelay)
}

// Slash is a paid mutator transaction binding the contract method 0x76ddeafc.
//
// Solidity: function slash(uint256 amt, uint256 slashAmt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistryTransactor) Slash(opts *bind.TransactOpts, amt *big.Int, slashAmt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "slash", amt, slashAmt, provider, bidder, residualBidPercentAfterDecay)
}

// Slash is a paid mutator transaction binding the contract method 0x76ddeafc.
//
// Solidity: function slash(uint256 amt, uint256 slashAmt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistrySession) Slash(amt *big.Int, slashAmt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Slash(&_Providerregistry.TransactOpts, amt, slashAmt, provider, bidder, residualBidPercentAfterDecay)
}

// Slash is a paid mutator transaction binding the contract method 0x76ddeafc.
//
// Solidity: function slash(uint256 amt, uint256 slashAmt, address provider, address bidder, uint256 residualBidPercentAfterDecay) returns()
func (_Providerregistry *ProviderregistryTransactorSession) Slash(amt *big.Int, slashAmt *big.Int, provider common.Address, bidder common.Address, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Providerregistry.Contract.Slash(&_Providerregistry.TransactOpts, amt, slashAmt, provider, bidder, residualBidPercentAfterDecay)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_Providerregistry *ProviderregistryTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_Providerregistry *ProviderregistrySession) Stake() (*types.Transaction, error) {
	return _Providerregistry.Contract.Stake(&_Providerregistry.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) Stake() (*types.Transaction, error) {
	return _Providerregistry.Contract.Stake(&_Providerregistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.TransferOwnership(&_Providerregistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Providerregistry *ProviderregistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Providerregistry.Contract.TransferOwnership(&_Providerregistry.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Providerregistry *ProviderregistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Providerregistry *ProviderregistrySession) Unpause() (*types.Transaction, error) {
	return _Providerregistry.Contract.Unpause(&_Providerregistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Providerregistry *ProviderregistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _Providerregistry.Contract.Unpause(&_Providerregistry.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Providerregistry *ProviderregistryTransactor) Unstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "unstake")
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Providerregistry *ProviderregistrySession) Unstake() (*types.Transaction, error) {
	return _Providerregistry.Contract.Unstake(&_Providerregistry.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Providerregistry *ProviderregistryTransactorSession) Unstake() (*types.Transaction, error) {
	return _Providerregistry.Contract.Unstake(&_Providerregistry.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Providerregistry *ProviderregistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Providerregistry *ProviderregistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.UpgradeToAndCall(&_Providerregistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.UpgradeToAndCall(&_Providerregistry.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Providerregistry *ProviderregistryTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Providerregistry *ProviderregistrySession) Withdraw() (*types.Transaction, error) {
	return _Providerregistry.Contract.Withdraw(&_Providerregistry.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_Providerregistry *ProviderregistryTransactorSession) Withdraw() (*types.Transaction, error) {
	return _Providerregistry.Contract.Withdraw(&_Providerregistry.TransactOpts)
}

// WithdrawSlashedAmount is a paid mutator transaction binding the contract method 0x70d6092b.
//
// Solidity: function withdrawSlashedAmount() returns()
func (_Providerregistry *ProviderregistryTransactor) WithdrawSlashedAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.Transact(opts, "withdrawSlashedAmount")
}

// WithdrawSlashedAmount is a paid mutator transaction binding the contract method 0x70d6092b.
//
// Solidity: function withdrawSlashedAmount() returns()
func (_Providerregistry *ProviderregistrySession) WithdrawSlashedAmount() (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawSlashedAmount(&_Providerregistry.TransactOpts)
}

// WithdrawSlashedAmount is a paid mutator transaction binding the contract method 0x70d6092b.
//
// Solidity: function withdrawSlashedAmount() returns()
func (_Providerregistry *ProviderregistryTransactorSession) WithdrawSlashedAmount() (*types.Transaction, error) {
	return _Providerregistry.Contract.WithdrawSlashedAmount(&_Providerregistry.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistryTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistrySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.Fallback(&_Providerregistry.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Providerregistry.Contract.Fallback(&_Providerregistry.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistryTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Providerregistry.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistrySession) Receive() (*types.Transaction, error) {
	return _Providerregistry.Contract.Receive(&_Providerregistry.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Providerregistry *ProviderregistryTransactorSession) Receive() (*types.Transaction, error) {
	return _Providerregistry.Contract.Receive(&_Providerregistry.TransactOpts)
}

// ProviderregistryBLSKeyAddedIterator is returned from FilterBLSKeyAdded and is used to iterate over the raw logs and unpacked data for BLSKeyAdded events raised by the Providerregistry contract.
type ProviderregistryBLSKeyAddedIterator struct {
	Event *ProviderregistryBLSKeyAdded // Event containing the contract specifics and raw log

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
func (it *ProviderregistryBLSKeyAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryBLSKeyAdded)
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
		it.Event = new(ProviderregistryBLSKeyAdded)
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
func (it *ProviderregistryBLSKeyAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryBLSKeyAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryBLSKeyAdded represents a BLSKeyAdded event raised by the Providerregistry contract.
type ProviderregistryBLSKeyAdded struct {
	Provider     common.Address
	BlsPublicKey []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBLSKeyAdded is a free log retrieval operation binding the contract event 0xc81315c03024fb67ddfe902ce0153b3d56572c0569de9564fcb90cc174a960bf.
//
// Solidity: event BLSKeyAdded(address indexed provider, bytes blsPublicKey)
func (_Providerregistry *ProviderregistryFilterer) FilterBLSKeyAdded(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryBLSKeyAddedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "BLSKeyAdded", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryBLSKeyAddedIterator{contract: _Providerregistry.contract, event: "BLSKeyAdded", logs: logs, sub: sub}, nil
}

// WatchBLSKeyAdded is a free log subscription operation binding the contract event 0xc81315c03024fb67ddfe902ce0153b3d56572c0569de9564fcb90cc174a960bf.
//
// Solidity: event BLSKeyAdded(address indexed provider, bytes blsPublicKey)
func (_Providerregistry *ProviderregistryFilterer) WatchBLSKeyAdded(opts *bind.WatchOpts, sink chan<- *ProviderregistryBLSKeyAdded, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "BLSKeyAdded", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryBLSKeyAdded)
				if err := _Providerregistry.contract.UnpackLog(event, "BLSKeyAdded", log); err != nil {
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

// ParseBLSKeyAdded is a log parse operation binding the contract event 0xc81315c03024fb67ddfe902ce0153b3d56572c0569de9564fcb90cc174a960bf.
//
// Solidity: event BLSKeyAdded(address indexed provider, bytes blsPublicKey)
func (_Providerregistry *ProviderregistryFilterer) ParseBLSKeyAdded(log types.Log) (*ProviderregistryBLSKeyAdded, error) {
	event := new(ProviderregistryBLSKeyAdded)
	if err := _Providerregistry.contract.UnpackLog(event, "BLSKeyAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryBidderWithdrawSlashedAmountIterator is returned from FilterBidderWithdrawSlashedAmount and is used to iterate over the raw logs and unpacked data for BidderWithdrawSlashedAmount events raised by the Providerregistry contract.
type ProviderregistryBidderWithdrawSlashedAmountIterator struct {
	Event *ProviderregistryBidderWithdrawSlashedAmount // Event containing the contract specifics and raw log

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
func (it *ProviderregistryBidderWithdrawSlashedAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryBidderWithdrawSlashedAmount)
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
		it.Event = new(ProviderregistryBidderWithdrawSlashedAmount)
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
func (it *ProviderregistryBidderWithdrawSlashedAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryBidderWithdrawSlashedAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryBidderWithdrawSlashedAmount represents a BidderWithdrawSlashedAmount event raised by the Providerregistry contract.
type ProviderregistryBidderWithdrawSlashedAmount struct {
	Bidder common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBidderWithdrawSlashedAmount is a free log retrieval operation binding the contract event 0x37f4dc41eefcb0d2f96341475d50bc699b0752135621b1fabba90b41e62fc68d.
//
// Solidity: event BidderWithdrawSlashedAmount(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterBidderWithdrawSlashedAmount(opts *bind.FilterOpts) (*ProviderregistryBidderWithdrawSlashedAmountIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "BidderWithdrawSlashedAmount")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryBidderWithdrawSlashedAmountIterator{contract: _Providerregistry.contract, event: "BidderWithdrawSlashedAmount", logs: logs, sub: sub}, nil
}

// WatchBidderWithdrawSlashedAmount is a free log subscription operation binding the contract event 0x37f4dc41eefcb0d2f96341475d50bc699b0752135621b1fabba90b41e62fc68d.
//
// Solidity: event BidderWithdrawSlashedAmount(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchBidderWithdrawSlashedAmount(opts *bind.WatchOpts, sink chan<- *ProviderregistryBidderWithdrawSlashedAmount) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "BidderWithdrawSlashedAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryBidderWithdrawSlashedAmount)
				if err := _Providerregistry.contract.UnpackLog(event, "BidderWithdrawSlashedAmount", log); err != nil {
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

// ParseBidderWithdrawSlashedAmount is a log parse operation binding the contract event 0x37f4dc41eefcb0d2f96341475d50bc699b0752135621b1fabba90b41e62fc68d.
//
// Solidity: event BidderWithdrawSlashedAmount(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseBidderWithdrawSlashedAmount(log types.Log) (*ProviderregistryBidderWithdrawSlashedAmount, error) {
	event := new(ProviderregistryBidderWithdrawSlashedAmount)
	if err := _Providerregistry.contract.UnpackLog(event, "BidderWithdrawSlashedAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFeePayoutPeriodBlocksUpdatedIterator is returned from FilterFeePayoutPeriodBlocksUpdated and is used to iterate over the raw logs and unpacked data for FeePayoutPeriodBlocksUpdated events raised by the Providerregistry contract.
type ProviderregistryFeePayoutPeriodBlocksUpdatedIterator struct {
	Event *ProviderregistryFeePayoutPeriodBlocksUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFeePayoutPeriodBlocksUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFeePayoutPeriodBlocksUpdated)
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
		it.Event = new(ProviderregistryFeePayoutPeriodBlocksUpdated)
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
func (it *ProviderregistryFeePayoutPeriodBlocksUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFeePayoutPeriodBlocksUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFeePayoutPeriodBlocksUpdated represents a FeePayoutPeriodBlocksUpdated event raised by the Providerregistry contract.
type ProviderregistryFeePayoutPeriodBlocksUpdated struct {
	NewFeePayoutPeriodBlocks *big.Int
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterFeePayoutPeriodBlocksUpdated is a free log retrieval operation binding the contract event 0x1b8b3f7fd7594ce5b7155b4c56b19bd6a1eac8c1ec5d941635acf104c8db3571.
//
// Solidity: event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks)
func (_Providerregistry *ProviderregistryFilterer) FilterFeePayoutPeriodBlocksUpdated(opts *bind.FilterOpts, newFeePayoutPeriodBlocks []*big.Int) (*ProviderregistryFeePayoutPeriodBlocksUpdatedIterator, error) {

	var newFeePayoutPeriodBlocksRule []interface{}
	for _, newFeePayoutPeriodBlocksItem := range newFeePayoutPeriodBlocks {
		newFeePayoutPeriodBlocksRule = append(newFeePayoutPeriodBlocksRule, newFeePayoutPeriodBlocksItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FeePayoutPeriodBlocksUpdated", newFeePayoutPeriodBlocksRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFeePayoutPeriodBlocksUpdatedIterator{contract: _Providerregistry.contract, event: "FeePayoutPeriodBlocksUpdated", logs: logs, sub: sub}, nil
}

// WatchFeePayoutPeriodBlocksUpdated is a free log subscription operation binding the contract event 0x1b8b3f7fd7594ce5b7155b4c56b19bd6a1eac8c1ec5d941635acf104c8db3571.
//
// Solidity: event FeePayoutPeriodBlocksUpdated(uint256 indexed newFeePayoutPeriodBlocks)
func (_Providerregistry *ProviderregistryFilterer) WatchFeePayoutPeriodBlocksUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryFeePayoutPeriodBlocksUpdated, newFeePayoutPeriodBlocks []*big.Int) (event.Subscription, error) {

	var newFeePayoutPeriodBlocksRule []interface{}
	for _, newFeePayoutPeriodBlocksItem := range newFeePayoutPeriodBlocks {
		newFeePayoutPeriodBlocksRule = append(newFeePayoutPeriodBlocksRule, newFeePayoutPeriodBlocksItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FeePayoutPeriodBlocksUpdated", newFeePayoutPeriodBlocksRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFeePayoutPeriodBlocksUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "FeePayoutPeriodBlocksUpdated", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseFeePayoutPeriodBlocksUpdated(log types.Log) (*ProviderregistryFeePayoutPeriodBlocksUpdated, error) {
	event := new(ProviderregistryFeePayoutPeriodBlocksUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "FeePayoutPeriodBlocksUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFeePercentUpdatedIterator is returned from FilterFeePercentUpdated and is used to iterate over the raw logs and unpacked data for FeePercentUpdated events raised by the Providerregistry contract.
type ProviderregistryFeePercentUpdatedIterator struct {
	Event *ProviderregistryFeePercentUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFeePercentUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFeePercentUpdated)
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
		it.Event = new(ProviderregistryFeePercentUpdated)
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
func (it *ProviderregistryFeePercentUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFeePercentUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFeePercentUpdated represents a FeePercentUpdated event raised by the Providerregistry contract.
type ProviderregistryFeePercentUpdated struct {
	NewFeePercent *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterFeePercentUpdated is a free log retrieval operation binding the contract event 0x64d1887078a96d281ed60dd69ba75bfb6b5cd2cb4c2d2538b2eb7816a4c646ea.
//
// Solidity: event FeePercentUpdated(uint256 indexed newFeePercent)
func (_Providerregistry *ProviderregistryFilterer) FilterFeePercentUpdated(opts *bind.FilterOpts, newFeePercent []*big.Int) (*ProviderregistryFeePercentUpdatedIterator, error) {

	var newFeePercentRule []interface{}
	for _, newFeePercentItem := range newFeePercent {
		newFeePercentRule = append(newFeePercentRule, newFeePercentItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FeePercentUpdated", newFeePercentRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFeePercentUpdatedIterator{contract: _Providerregistry.contract, event: "FeePercentUpdated", logs: logs, sub: sub}, nil
}

// WatchFeePercentUpdated is a free log subscription operation binding the contract event 0x64d1887078a96d281ed60dd69ba75bfb6b5cd2cb4c2d2538b2eb7816a4c646ea.
//
// Solidity: event FeePercentUpdated(uint256 indexed newFeePercent)
func (_Providerregistry *ProviderregistryFilterer) WatchFeePercentUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryFeePercentUpdated, newFeePercent []*big.Int) (event.Subscription, error) {

	var newFeePercentRule []interface{}
	for _, newFeePercentItem := range newFeePercent {
		newFeePercentRule = append(newFeePercentRule, newFeePercentItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FeePercentUpdated", newFeePercentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFeePercentUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "FeePercentUpdated", log); err != nil {
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

// ParseFeePercentUpdated is a log parse operation binding the contract event 0x64d1887078a96d281ed60dd69ba75bfb6b5cd2cb4c2d2538b2eb7816a4c646ea.
//
// Solidity: event FeePercentUpdated(uint256 indexed newFeePercent)
func (_Providerregistry *ProviderregistryFilterer) ParseFeePercentUpdated(log types.Log) (*ProviderregistryFeePercentUpdated, error) {
	event := new(ProviderregistryFeePercentUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "FeePercentUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFeeTransferIterator is returned from FilterFeeTransfer and is used to iterate over the raw logs and unpacked data for FeeTransfer events raised by the Providerregistry contract.
type ProviderregistryFeeTransferIterator struct {
	Event *ProviderregistryFeeTransfer // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFeeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFeeTransfer)
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
		it.Event = new(ProviderregistryFeeTransfer)
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
func (it *ProviderregistryFeeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFeeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFeeTransfer represents a FeeTransfer event raised by the Providerregistry contract.
type ProviderregistryFeeTransfer struct {
	Amount    *big.Int
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeTransfer is a free log retrieval operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Providerregistry *ProviderregistryFilterer) FilterFeeTransfer(opts *bind.FilterOpts, recipient []common.Address) (*ProviderregistryFeeTransferIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFeeTransferIterator{contract: _Providerregistry.contract, event: "FeeTransfer", logs: logs, sub: sub}, nil
}

// WatchFeeTransfer is a free log subscription operation binding the contract event 0x445bb6587d6cd09e272a0d1e5179e772b547dbf1041b6163f86bb62e86f25031.
//
// Solidity: event FeeTransfer(uint256 amount, address indexed recipient)
func (_Providerregistry *ProviderregistryFilterer) WatchFeeTransfer(opts *bind.WatchOpts, sink chan<- *ProviderregistryFeeTransfer, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FeeTransfer", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFeeTransfer)
				if err := _Providerregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseFeeTransfer(log types.Log) (*ProviderregistryFeeTransfer, error) {
	event := new(ProviderregistryFeeTransfer)
	if err := _Providerregistry.contract.UnpackLog(event, "FeeTransfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFundsDepositedIterator is returned from FilterFundsDeposited and is used to iterate over the raw logs and unpacked data for FundsDeposited events raised by the Providerregistry contract.
type ProviderregistryFundsDepositedIterator struct {
	Event *ProviderregistryFundsDeposited // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFundsDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFundsDeposited)
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
		it.Event = new(ProviderregistryFundsDeposited)
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
func (it *ProviderregistryFundsDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFundsDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFundsDeposited represents a FundsDeposited event raised by the Providerregistry contract.
type ProviderregistryFundsDeposited struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFundsDeposited is a free log retrieval operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterFundsDeposited(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryFundsDepositedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FundsDeposited", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFundsDepositedIterator{contract: _Providerregistry.contract, event: "FundsDeposited", logs: logs, sub: sub}, nil
}

// WatchFundsDeposited is a free log subscription operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchFundsDeposited(opts *bind.WatchOpts, sink chan<- *ProviderregistryFundsDeposited, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FundsDeposited", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFundsDeposited)
				if err := _Providerregistry.contract.UnpackLog(event, "FundsDeposited", log); err != nil {
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

// ParseFundsDeposited is a log parse operation binding the contract event 0x543ba50a5eec5e6178218e364b1d0f396157b3c8fa278522c2cb7fd99407d474.
//
// Solidity: event FundsDeposited(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseFundsDeposited(log types.Log) (*ProviderregistryFundsDeposited, error) {
	event := new(ProviderregistryFundsDeposited)
	if err := _Providerregistry.contract.UnpackLog(event, "FundsDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryFundsSlashedIterator is returned from FilterFundsSlashed and is used to iterate over the raw logs and unpacked data for FundsSlashed events raised by the Providerregistry contract.
type ProviderregistryFundsSlashedIterator struct {
	Event *ProviderregistryFundsSlashed // Event containing the contract specifics and raw log

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
func (it *ProviderregistryFundsSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryFundsSlashed)
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
		it.Event = new(ProviderregistryFundsSlashed)
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
func (it *ProviderregistryFundsSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryFundsSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryFundsSlashed represents a FundsSlashed event raised by the Providerregistry contract.
type ProviderregistryFundsSlashed struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFundsSlashed is a free log retrieval operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterFundsSlashed(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryFundsSlashedIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "FundsSlashed", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryFundsSlashedIterator{contract: _Providerregistry.contract, event: "FundsSlashed", logs: logs, sub: sub}, nil
}

// WatchFundsSlashed is a free log subscription operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchFundsSlashed(opts *bind.WatchOpts, sink chan<- *ProviderregistryFundsSlashed, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "FundsSlashed", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryFundsSlashed)
				if err := _Providerregistry.contract.UnpackLog(event, "FundsSlashed", log); err != nil {
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

// ParseFundsSlashed is a log parse operation binding the contract event 0x4a00481d3f7b0802643df0bdfb9bfc491a24ffca3eb1becc9fe8b525e0427a74.
//
// Solidity: event FundsSlashed(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseFundsSlashed(log types.Log) (*ProviderregistryFundsSlashed, error) {
	event := new(ProviderregistryFundsSlashed)
	if err := _Providerregistry.contract.UnpackLog(event, "FundsSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Providerregistry contract.
type ProviderregistryInitializedIterator struct {
	Event *ProviderregistryInitialized // Event containing the contract specifics and raw log

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
func (it *ProviderregistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryInitialized)
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
		it.Event = new(ProviderregistryInitialized)
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
func (it *ProviderregistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryInitialized represents a Initialized event raised by the Providerregistry contract.
type ProviderregistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Providerregistry *ProviderregistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*ProviderregistryInitializedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryInitializedIterator{contract: _Providerregistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Providerregistry *ProviderregistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ProviderregistryInitialized) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryInitialized)
				if err := _Providerregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseInitialized(log types.Log) (*ProviderregistryInitialized, error) {
	event := new(ProviderregistryInitialized)
	if err := _Providerregistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryInsufficientFundsToSlashIterator is returned from FilterInsufficientFundsToSlash and is used to iterate over the raw logs and unpacked data for InsufficientFundsToSlash events raised by the Providerregistry contract.
type ProviderregistryInsufficientFundsToSlashIterator struct {
	Event *ProviderregistryInsufficientFundsToSlash // Event containing the contract specifics and raw log

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
func (it *ProviderregistryInsufficientFundsToSlashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryInsufficientFundsToSlash)
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
		it.Event = new(ProviderregistryInsufficientFundsToSlash)
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
func (it *ProviderregistryInsufficientFundsToSlashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryInsufficientFundsToSlashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryInsufficientFundsToSlash represents a InsufficientFundsToSlash event raised by the Providerregistry contract.
type ProviderregistryInsufficientFundsToSlash struct {
	Provider       common.Address
	ProviderStake  *big.Int
	ResidualAmount *big.Int
	PenaltyFee     *big.Int
	SlashAmt       *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterInsufficientFundsToSlash is a free log retrieval operation binding the contract event 0x358857ec44573677cd31f4c53a96a9a253bfaea0ae485b613fd33a99bacd79dd.
//
// Solidity: event InsufficientFundsToSlash(address indexed provider, uint256 providerStake, uint256 residualAmount, uint256 penaltyFee, uint256 slashAmt)
func (_Providerregistry *ProviderregistryFilterer) FilterInsufficientFundsToSlash(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryInsufficientFundsToSlashIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "InsufficientFundsToSlash", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryInsufficientFundsToSlashIterator{contract: _Providerregistry.contract, event: "InsufficientFundsToSlash", logs: logs, sub: sub}, nil
}

// WatchInsufficientFundsToSlash is a free log subscription operation binding the contract event 0x358857ec44573677cd31f4c53a96a9a253bfaea0ae485b613fd33a99bacd79dd.
//
// Solidity: event InsufficientFundsToSlash(address indexed provider, uint256 providerStake, uint256 residualAmount, uint256 penaltyFee, uint256 slashAmt)
func (_Providerregistry *ProviderregistryFilterer) WatchInsufficientFundsToSlash(opts *bind.WatchOpts, sink chan<- *ProviderregistryInsufficientFundsToSlash, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "InsufficientFundsToSlash", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryInsufficientFundsToSlash)
				if err := _Providerregistry.contract.UnpackLog(event, "InsufficientFundsToSlash", log); err != nil {
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

// ParseInsufficientFundsToSlash is a log parse operation binding the contract event 0x358857ec44573677cd31f4c53a96a9a253bfaea0ae485b613fd33a99bacd79dd.
//
// Solidity: event InsufficientFundsToSlash(address indexed provider, uint256 providerStake, uint256 residualAmount, uint256 penaltyFee, uint256 slashAmt)
func (_Providerregistry *ProviderregistryFilterer) ParseInsufficientFundsToSlash(log types.Log) (*ProviderregistryInsufficientFundsToSlash, error) {
	event := new(ProviderregistryInsufficientFundsToSlash)
	if err := _Providerregistry.contract.UnpackLog(event, "InsufficientFundsToSlash", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryMinStakeUpdatedIterator is returned from FilterMinStakeUpdated and is used to iterate over the raw logs and unpacked data for MinStakeUpdated events raised by the Providerregistry contract.
type ProviderregistryMinStakeUpdatedIterator struct {
	Event *ProviderregistryMinStakeUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryMinStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryMinStakeUpdated)
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
		it.Event = new(ProviderregistryMinStakeUpdated)
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
func (it *ProviderregistryMinStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryMinStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryMinStakeUpdated represents a MinStakeUpdated event raised by the Providerregistry contract.
type ProviderregistryMinStakeUpdated struct {
	NewMinStake *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMinStakeUpdated is a free log retrieval operation binding the contract event 0x47ab46f2c8d4258304a2f5551c1cbdb6981be49631365d1ba7191288a73f39ef.
//
// Solidity: event MinStakeUpdated(uint256 indexed newMinStake)
func (_Providerregistry *ProviderregistryFilterer) FilterMinStakeUpdated(opts *bind.FilterOpts, newMinStake []*big.Int) (*ProviderregistryMinStakeUpdatedIterator, error) {

	var newMinStakeRule []interface{}
	for _, newMinStakeItem := range newMinStake {
		newMinStakeRule = append(newMinStakeRule, newMinStakeItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "MinStakeUpdated", newMinStakeRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryMinStakeUpdatedIterator{contract: _Providerregistry.contract, event: "MinStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinStakeUpdated is a free log subscription operation binding the contract event 0x47ab46f2c8d4258304a2f5551c1cbdb6981be49631365d1ba7191288a73f39ef.
//
// Solidity: event MinStakeUpdated(uint256 indexed newMinStake)
func (_Providerregistry *ProviderregistryFilterer) WatchMinStakeUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryMinStakeUpdated, newMinStake []*big.Int) (event.Subscription, error) {

	var newMinStakeRule []interface{}
	for _, newMinStakeItem := range newMinStake {
		newMinStakeRule = append(newMinStakeRule, newMinStakeItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "MinStakeUpdated", newMinStakeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryMinStakeUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
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

// ParseMinStakeUpdated is a log parse operation binding the contract event 0x47ab46f2c8d4258304a2f5551c1cbdb6981be49631365d1ba7191288a73f39ef.
//
// Solidity: event MinStakeUpdated(uint256 indexed newMinStake)
func (_Providerregistry *ProviderregistryFilterer) ParseMinStakeUpdated(log types.Log) (*ProviderregistryMinStakeUpdated, error) {
	event := new(ProviderregistryMinStakeUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferStartedIterator struct {
	Event *ProviderregistryOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ProviderregistryOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryOwnershipTransferStarted)
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
		it.Event = new(ProviderregistryOwnershipTransferStarted)
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
func (it *ProviderregistryOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProviderregistryOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryOwnershipTransferStartedIterator{contract: _Providerregistry.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ProviderregistryOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryOwnershipTransferStarted)
				if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseOwnershipTransferStarted(log types.Log) (*ProviderregistryOwnershipTransferStarted, error) {
	event := new(ProviderregistryOwnershipTransferStarted)
	if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferredIterator struct {
	Event *ProviderregistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ProviderregistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryOwnershipTransferred)
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
		it.Event = new(ProviderregistryOwnershipTransferred)
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
func (it *ProviderregistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryOwnershipTransferred represents a OwnershipTransferred event raised by the Providerregistry contract.
type ProviderregistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ProviderregistryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryOwnershipTransferredIterator{contract: _Providerregistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Providerregistry *ProviderregistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ProviderregistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryOwnershipTransferred)
				if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseOwnershipTransferred(log types.Log) (*ProviderregistryOwnershipTransferred, error) {
	event := new(ProviderregistryOwnershipTransferred)
	if err := _Providerregistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Providerregistry contract.
type ProviderregistryPausedIterator struct {
	Event *ProviderregistryPaused // Event containing the contract specifics and raw log

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
func (it *ProviderregistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryPaused)
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
		it.Event = new(ProviderregistryPaused)
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
func (it *ProviderregistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryPaused represents a Paused event raised by the Providerregistry contract.
type ProviderregistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Providerregistry *ProviderregistryFilterer) FilterPaused(opts *bind.FilterOpts) (*ProviderregistryPausedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryPausedIterator{contract: _Providerregistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Providerregistry *ProviderregistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ProviderregistryPaused) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryPaused)
				if err := _Providerregistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParsePaused(log types.Log) (*ProviderregistryPaused, error) {
	event := new(ProviderregistryPaused)
	if err := _Providerregistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryPenaltyFeeRecipientUpdatedIterator is returned from FilterPenaltyFeeRecipientUpdated and is used to iterate over the raw logs and unpacked data for PenaltyFeeRecipientUpdated events raised by the Providerregistry contract.
type ProviderregistryPenaltyFeeRecipientUpdatedIterator struct {
	Event *ProviderregistryPenaltyFeeRecipientUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryPenaltyFeeRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryPenaltyFeeRecipientUpdated)
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
		it.Event = new(ProviderregistryPenaltyFeeRecipientUpdated)
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
func (it *ProviderregistryPenaltyFeeRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryPenaltyFeeRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryPenaltyFeeRecipientUpdated represents a PenaltyFeeRecipientUpdated event raised by the Providerregistry contract.
type ProviderregistryPenaltyFeeRecipientUpdated struct {
	NewPenaltyFeeRecipient common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterPenaltyFeeRecipientUpdated is a free log retrieval operation binding the contract event 0xb5d4f6e9d9b31eb58c205f8c3bb44d4f0094605caa42167107a9c3e91ebf8c85.
//
// Solidity: event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient)
func (_Providerregistry *ProviderregistryFilterer) FilterPenaltyFeeRecipientUpdated(opts *bind.FilterOpts, newPenaltyFeeRecipient []common.Address) (*ProviderregistryPenaltyFeeRecipientUpdatedIterator, error) {

	var newPenaltyFeeRecipientRule []interface{}
	for _, newPenaltyFeeRecipientItem := range newPenaltyFeeRecipient {
		newPenaltyFeeRecipientRule = append(newPenaltyFeeRecipientRule, newPenaltyFeeRecipientItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "PenaltyFeeRecipientUpdated", newPenaltyFeeRecipientRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryPenaltyFeeRecipientUpdatedIterator{contract: _Providerregistry.contract, event: "PenaltyFeeRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchPenaltyFeeRecipientUpdated is a free log subscription operation binding the contract event 0xb5d4f6e9d9b31eb58c205f8c3bb44d4f0094605caa42167107a9c3e91ebf8c85.
//
// Solidity: event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient)
func (_Providerregistry *ProviderregistryFilterer) WatchPenaltyFeeRecipientUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryPenaltyFeeRecipientUpdated, newPenaltyFeeRecipient []common.Address) (event.Subscription, error) {

	var newPenaltyFeeRecipientRule []interface{}
	for _, newPenaltyFeeRecipientItem := range newPenaltyFeeRecipient {
		newPenaltyFeeRecipientRule = append(newPenaltyFeeRecipientRule, newPenaltyFeeRecipientItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "PenaltyFeeRecipientUpdated", newPenaltyFeeRecipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryPenaltyFeeRecipientUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "PenaltyFeeRecipientUpdated", log); err != nil {
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

// ParsePenaltyFeeRecipientUpdated is a log parse operation binding the contract event 0xb5d4f6e9d9b31eb58c205f8c3bb44d4f0094605caa42167107a9c3e91ebf8c85.
//
// Solidity: event PenaltyFeeRecipientUpdated(address indexed newPenaltyFeeRecipient)
func (_Providerregistry *ProviderregistryFilterer) ParsePenaltyFeeRecipientUpdated(log types.Log) (*ProviderregistryPenaltyFeeRecipientUpdated, error) {
	event := new(ProviderregistryPenaltyFeeRecipientUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "PenaltyFeeRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryPreconfManagerUpdatedIterator is returned from FilterPreconfManagerUpdated and is used to iterate over the raw logs and unpacked data for PreconfManagerUpdated events raised by the Providerregistry contract.
type ProviderregistryPreconfManagerUpdatedIterator struct {
	Event *ProviderregistryPreconfManagerUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryPreconfManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryPreconfManagerUpdated)
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
		it.Event = new(ProviderregistryPreconfManagerUpdated)
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
func (it *ProviderregistryPreconfManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryPreconfManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryPreconfManagerUpdated represents a PreconfManagerUpdated event raised by the Providerregistry contract.
type ProviderregistryPreconfManagerUpdated struct {
	NewPreconfManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterPreconfManagerUpdated is a free log retrieval operation binding the contract event 0xba8b46bd4ed6a7ec49dc1a4758a5ffd0cbae99c172bbe007676fcb74fefa310f.
//
// Solidity: event PreconfManagerUpdated(address indexed newPreconfManager)
func (_Providerregistry *ProviderregistryFilterer) FilterPreconfManagerUpdated(opts *bind.FilterOpts, newPreconfManager []common.Address) (*ProviderregistryPreconfManagerUpdatedIterator, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "PreconfManagerUpdated", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryPreconfManagerUpdatedIterator{contract: _Providerregistry.contract, event: "PreconfManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchPreconfManagerUpdated is a free log subscription operation binding the contract event 0xba8b46bd4ed6a7ec49dc1a4758a5ffd0cbae99c172bbe007676fcb74fefa310f.
//
// Solidity: event PreconfManagerUpdated(address indexed newPreconfManager)
func (_Providerregistry *ProviderregistryFilterer) WatchPreconfManagerUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryPreconfManagerUpdated, newPreconfManager []common.Address) (event.Subscription, error) {

	var newPreconfManagerRule []interface{}
	for _, newPreconfManagerItem := range newPreconfManager {
		newPreconfManagerRule = append(newPreconfManagerRule, newPreconfManagerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "PreconfManagerUpdated", newPreconfManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryPreconfManagerUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "PreconfManagerUpdated", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParsePreconfManagerUpdated(log types.Log) (*ProviderregistryPreconfManagerUpdated, error) {
	event := new(ProviderregistryPreconfManagerUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "PreconfManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryProviderRegisteredIterator is returned from FilterProviderRegistered and is used to iterate over the raw logs and unpacked data for ProviderRegistered events raised by the Providerregistry contract.
type ProviderregistryProviderRegisteredIterator struct {
	Event *ProviderregistryProviderRegistered // Event containing the contract specifics and raw log

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
func (it *ProviderregistryProviderRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryProviderRegistered)
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
		it.Event = new(ProviderregistryProviderRegistered)
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
func (it *ProviderregistryProviderRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryProviderRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryProviderRegistered represents a ProviderRegistered event raised by the Providerregistry contract.
type ProviderregistryProviderRegistered struct {
	Provider     common.Address
	StakedAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterProviderRegistered is a free log retrieval operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) FilterProviderRegistered(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryProviderRegisteredIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryProviderRegisteredIterator{contract: _Providerregistry.contract, event: "ProviderRegistered", logs: logs, sub: sub}, nil
}

// WatchProviderRegistered is a free log subscription operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) WatchProviderRegistered(opts *bind.WatchOpts, sink chan<- *ProviderregistryProviderRegistered, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "ProviderRegistered", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryProviderRegistered)
				if err := _Providerregistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
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

// ParseProviderRegistered is a log parse operation binding the contract event 0x90c9734131c1e4fb36cde2d71e6feb93fb258f71be8a85411c173d25e1516e80.
//
// Solidity: event ProviderRegistered(address indexed provider, uint256 stakedAmount)
func (_Providerregistry *ProviderregistryFilterer) ParseProviderRegistered(log types.Log) (*ProviderregistryProviderRegistered, error) {
	event := new(ProviderregistryProviderRegistered)
	if err := _Providerregistry.contract.UnpackLog(event, "ProviderRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryTransferToBidderFailedIterator is returned from FilterTransferToBidderFailed and is used to iterate over the raw logs and unpacked data for TransferToBidderFailed events raised by the Providerregistry contract.
type ProviderregistryTransferToBidderFailedIterator struct {
	Event *ProviderregistryTransferToBidderFailed // Event containing the contract specifics and raw log

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
func (it *ProviderregistryTransferToBidderFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryTransferToBidderFailed)
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
		it.Event = new(ProviderregistryTransferToBidderFailed)
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
func (it *ProviderregistryTransferToBidderFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryTransferToBidderFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryTransferToBidderFailed represents a TransferToBidderFailed event raised by the Providerregistry contract.
type ProviderregistryTransferToBidderFailed struct {
	Bidder common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransferToBidderFailed is a free log retrieval operation binding the contract event 0xd4bd5c1c0f198fbafd25af26c36f4c115af31d522d0f520abc017845e225aca6.
//
// Solidity: event TransferToBidderFailed(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterTransferToBidderFailed(opts *bind.FilterOpts) (*ProviderregistryTransferToBidderFailedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "TransferToBidderFailed")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryTransferToBidderFailedIterator{contract: _Providerregistry.contract, event: "TransferToBidderFailed", logs: logs, sub: sub}, nil
}

// WatchTransferToBidderFailed is a free log subscription operation binding the contract event 0xd4bd5c1c0f198fbafd25af26c36f4c115af31d522d0f520abc017845e225aca6.
//
// Solidity: event TransferToBidderFailed(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchTransferToBidderFailed(opts *bind.WatchOpts, sink chan<- *ProviderregistryTransferToBidderFailed) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "TransferToBidderFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryTransferToBidderFailed)
				if err := _Providerregistry.contract.UnpackLog(event, "TransferToBidderFailed", log); err != nil {
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

// ParseTransferToBidderFailed is a log parse operation binding the contract event 0xd4bd5c1c0f198fbafd25af26c36f4c115af31d522d0f520abc017845e225aca6.
//
// Solidity: event TransferToBidderFailed(address bidder, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseTransferToBidderFailed(log types.Log) (*ProviderregistryTransferToBidderFailed, error) {
	event := new(ProviderregistryTransferToBidderFailed)
	if err := _Providerregistry.contract.UnpackLog(event, "TransferToBidderFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Providerregistry contract.
type ProviderregistryUnpausedIterator struct {
	Event *ProviderregistryUnpaused // Event containing the contract specifics and raw log

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
func (it *ProviderregistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryUnpaused)
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
		it.Event = new(ProviderregistryUnpaused)
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
func (it *ProviderregistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryUnpaused represents a Unpaused event raised by the Providerregistry contract.
type ProviderregistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Providerregistry *ProviderregistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ProviderregistryUnpausedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryUnpausedIterator{contract: _Providerregistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Providerregistry *ProviderregistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ProviderregistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryUnpaused)
				if err := _Providerregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseUnpaused(log types.Log) (*ProviderregistryUnpaused, error) {
	event := new(ProviderregistryUnpaused)
	if err := _Providerregistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryUnstakeIterator is returned from FilterUnstake and is used to iterate over the raw logs and unpacked data for Unstake events raised by the Providerregistry contract.
type ProviderregistryUnstakeIterator struct {
	Event *ProviderregistryUnstake // Event containing the contract specifics and raw log

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
func (it *ProviderregistryUnstakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryUnstake)
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
		it.Event = new(ProviderregistryUnstake)
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
func (it *ProviderregistryUnstakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryUnstakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryUnstake represents a Unstake event raised by the Providerregistry contract.
type ProviderregistryUnstake struct {
	Provider  common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnstake is a free log retrieval operation binding the contract event 0x85082129d87b2fe11527cb1b3b7a520aeb5aa6913f88a3d8757fe40d1db02fdd.
//
// Solidity: event Unstake(address indexed provider, uint256 timestamp)
func (_Providerregistry *ProviderregistryFilterer) FilterUnstake(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryUnstakeIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Unstake", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryUnstakeIterator{contract: _Providerregistry.contract, event: "Unstake", logs: logs, sub: sub}, nil
}

// WatchUnstake is a free log subscription operation binding the contract event 0x85082129d87b2fe11527cb1b3b7a520aeb5aa6913f88a3d8757fe40d1db02fdd.
//
// Solidity: event Unstake(address indexed provider, uint256 timestamp)
func (_Providerregistry *ProviderregistryFilterer) WatchUnstake(opts *bind.WatchOpts, sink chan<- *ProviderregistryUnstake, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Unstake", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryUnstake)
				if err := _Providerregistry.contract.UnpackLog(event, "Unstake", log); err != nil {
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

// ParseUnstake is a log parse operation binding the contract event 0x85082129d87b2fe11527cb1b3b7a520aeb5aa6913f88a3d8757fe40d1db02fdd.
//
// Solidity: event Unstake(address indexed provider, uint256 timestamp)
func (_Providerregistry *ProviderregistryFilterer) ParseUnstake(log types.Log) (*ProviderregistryUnstake, error) {
	event := new(ProviderregistryUnstake)
	if err := _Providerregistry.contract.UnpackLog(event, "Unstake", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Providerregistry contract.
type ProviderregistryUpgradedIterator struct {
	Event *ProviderregistryUpgraded // Event containing the contract specifics and raw log

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
func (it *ProviderregistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryUpgraded)
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
		it.Event = new(ProviderregistryUpgraded)
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
func (it *ProviderregistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryUpgraded represents a Upgraded event raised by the Providerregistry contract.
type ProviderregistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Providerregistry *ProviderregistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ProviderregistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryUpgradedIterator{contract: _Providerregistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Providerregistry *ProviderregistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ProviderregistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryUpgraded)
				if err := _Providerregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Providerregistry *ProviderregistryFilterer) ParseUpgraded(log types.Log) (*ProviderregistryUpgraded, error) {
	event := new(ProviderregistryUpgraded)
	if err := _Providerregistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Providerregistry contract.
type ProviderregistryWithdrawIterator struct {
	Event *ProviderregistryWithdraw // Event containing the contract specifics and raw log

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
func (it *ProviderregistryWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryWithdraw)
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
		it.Event = new(ProviderregistryWithdraw)
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
func (it *ProviderregistryWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryWithdraw represents a Withdraw event raised by the Providerregistry contract.
type ProviderregistryWithdraw struct {
	Provider common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) FilterWithdraw(opts *bind.FilterOpts, provider []common.Address) (*ProviderregistryWithdrawIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "Withdraw", providerRule)
	if err != nil {
		return nil, err
	}
	return &ProviderregistryWithdrawIterator{contract: _Providerregistry.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *ProviderregistryWithdraw, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "Withdraw", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryWithdraw)
				if err := _Providerregistry.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed provider, uint256 amount)
func (_Providerregistry *ProviderregistryFilterer) ParseWithdraw(log types.Log) (*ProviderregistryWithdraw, error) {
	event := new(ProviderregistryWithdraw)
	if err := _Providerregistry.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProviderregistryWithdrawalDelayUpdatedIterator is returned from FilterWithdrawalDelayUpdated and is used to iterate over the raw logs and unpacked data for WithdrawalDelayUpdated events raised by the Providerregistry contract.
type ProviderregistryWithdrawalDelayUpdatedIterator struct {
	Event *ProviderregistryWithdrawalDelayUpdated // Event containing the contract specifics and raw log

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
func (it *ProviderregistryWithdrawalDelayUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProviderregistryWithdrawalDelayUpdated)
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
		it.Event = new(ProviderregistryWithdrawalDelayUpdated)
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
func (it *ProviderregistryWithdrawalDelayUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProviderregistryWithdrawalDelayUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProviderregistryWithdrawalDelayUpdated represents a WithdrawalDelayUpdated event raised by the Providerregistry contract.
type ProviderregistryWithdrawalDelayUpdated struct {
	NewWithdrawalDelay *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalDelayUpdated is a free log retrieval operation binding the contract event 0xb34b8e54226bc5e8f4a9d846e344e0c42d09eeb1c79798df7ce7460c16071e2d.
//
// Solidity: event WithdrawalDelayUpdated(uint256 newWithdrawalDelay)
func (_Providerregistry *ProviderregistryFilterer) FilterWithdrawalDelayUpdated(opts *bind.FilterOpts) (*ProviderregistryWithdrawalDelayUpdatedIterator, error) {

	logs, sub, err := _Providerregistry.contract.FilterLogs(opts, "WithdrawalDelayUpdated")
	if err != nil {
		return nil, err
	}
	return &ProviderregistryWithdrawalDelayUpdatedIterator{contract: _Providerregistry.contract, event: "WithdrawalDelayUpdated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalDelayUpdated is a free log subscription operation binding the contract event 0xb34b8e54226bc5e8f4a9d846e344e0c42d09eeb1c79798df7ce7460c16071e2d.
//
// Solidity: event WithdrawalDelayUpdated(uint256 newWithdrawalDelay)
func (_Providerregistry *ProviderregistryFilterer) WatchWithdrawalDelayUpdated(opts *bind.WatchOpts, sink chan<- *ProviderregistryWithdrawalDelayUpdated) (event.Subscription, error) {

	logs, sub, err := _Providerregistry.contract.WatchLogs(opts, "WithdrawalDelayUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProviderregistryWithdrawalDelayUpdated)
				if err := _Providerregistry.contract.UnpackLog(event, "WithdrawalDelayUpdated", log); err != nil {
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

// ParseWithdrawalDelayUpdated is a log parse operation binding the contract event 0xb34b8e54226bc5e8f4a9d846e344e0c42d09eeb1c79798df7ce7460c16071e2d.
//
// Solidity: event WithdrawalDelayUpdated(uint256 newWithdrawalDelay)
func (_Providerregistry *ProviderregistryFilterer) ParseWithdrawalDelayUpdated(log types.Log) (*ProviderregistryWithdrawalDelayUpdated, error) {
	event := new(ProviderregistryWithdrawalDelayUpdated)
	if err := _Providerregistry.contract.UnpackLog(event, "WithdrawalDelayUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
