// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

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

// VaultMetaData contains all meta data concerning the Vault contract.
var VaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"delegatorFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slasherFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vaultFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DELEGATOR_FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEPOSITOR_WHITELIST_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEPOSIT_LIMIT_SET_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEPOSIT_WHITELIST_SET_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"IS_DEPOSIT_LIMIT_SET_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SLASHER_FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeBalanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeBalanceOfAt\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"hints\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeShares\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeSharesAt\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"hint\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeSharesOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeSharesOfAt\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"hint\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeStakeAt\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"hint\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimBatch\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"epochs\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collateral\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentEpochStart\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"onBehalfOf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"depositedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mintedShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositWhitelist\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochAt\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochDurationInit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initialVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isDelegatorInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDepositLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDepositorWhitelisted\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSlasherInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isWithdrawalsClaimed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[{\"name\":\"newVersion\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nextEpochStart\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onSlash\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"captureTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[{\"name\":\"slashedAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previousEpochStart\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"redeem\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"shares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"withdrawnAssets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mintedShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegator\",\"inputs\":[{\"name\":\"delegator_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositLimit\",\"inputs\":[{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositWhitelist\",\"inputs\":[{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDepositorWhitelistStatus\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setIsDepositLimit\",\"inputs\":[{\"name\":\"status\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlasher\",\"inputs\":[{\"name\":\"slasher_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashableBalanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slasher\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staticDelegateCall\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalStake\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"burnedShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"mintedShares\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawalShares\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawalSharesOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawals\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawalsOf\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Claim\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ClaimBatch\",\"inputs\":[{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"epochs\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"onBehalfOf\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"shares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OnSlash\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"captureTimestamp\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"},{\"name\":\"slashedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetDelegator\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetDepositLimit\",\"inputs\":[{\"name\":\"limit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetDepositWhitelist\",\"inputs\":[{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetDepositorWhitelistStatus\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetIsDepositLimit\",\"inputs\":[{\"name\":\"status\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetSlasher\",\"inputs\":[{\"name\":\"slasher\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"withdrawer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"claimer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"burnedShares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"mintedShares\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AddressInsufficientBalance\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckpointUnorderedInsertion\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DelegatorAlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DepositLimitReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientRedemption\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientWithdrawal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAccount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCaptureEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidClaimer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCollateral\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDelegator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEpochDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidLengthEpochs\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOnBehalfOf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSlasher\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTimestamp\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MathOverflowedMulDiv\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MissingRoles\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoPreviousEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotDelegator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotFactory\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSlasher\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotWhitelistedDepositor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlasherAlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooMuchRedeem\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooMuchWithdraw\",\"inputs\":[]}]",
}

// VaultABI is the input ABI used to generate the binding from.
// Deprecated: Use VaultMetaData.ABI instead.
var VaultABI = VaultMetaData.ABI

// Vault is an auto generated Go binding around an Ethereum contract.
type Vault struct {
	VaultCaller     // Read-only binding to the contract
	VaultTransactor // Write-only binding to the contract
	VaultFilterer   // Log filterer for contract events
}

// VaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type VaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VaultSession struct {
	Contract     *Vault            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VaultCallerSession struct {
	Contract *VaultCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VaultTransactorSession struct {
	Contract     *VaultTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type VaultRaw struct {
	Contract *Vault // Generic contract binding to access the raw methods on
}

// VaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VaultCallerRaw struct {
	Contract *VaultCaller // Generic read-only contract binding to access the raw methods on
}

// VaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VaultTransactorRaw struct {
	Contract *VaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVault creates a new instance of Vault, bound to a specific deployed contract.
func NewVault(address common.Address, backend bind.ContractBackend) (*Vault, error) {
	contract, err := bindVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vault{VaultCaller: VaultCaller{contract: contract}, VaultTransactor: VaultTransactor{contract: contract}, VaultFilterer: VaultFilterer{contract: contract}}, nil
}

// NewVaultCaller creates a new read-only instance of Vault, bound to a specific deployed contract.
func NewVaultCaller(address common.Address, caller bind.ContractCaller) (*VaultCaller, error) {
	contract, err := bindVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VaultCaller{contract: contract}, nil
}

// NewVaultTransactor creates a new write-only instance of Vault, bound to a specific deployed contract.
func NewVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*VaultTransactor, error) {
	contract, err := bindVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VaultTransactor{contract: contract}, nil
}

// NewVaultFilterer creates a new log filterer instance of Vault, bound to a specific deployed contract.
func NewVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*VaultFilterer, error) {
	contract, err := bindVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VaultFilterer{contract: contract}, nil
}

// bindVault binds a generic wrapper to an already deployed contract.
func bindVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.VaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Vault *VaultCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Vault *VaultSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Vault.Contract.DEFAULTADMINROLE(&_Vault.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Vault *VaultCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Vault.Contract.DEFAULTADMINROLE(&_Vault.CallOpts)
}

// DELEGATORFACTORY is a free data retrieval call binding the contract method 0x6da3e06d.
//
// Solidity: function DELEGATOR_FACTORY() view returns(address)
func (_Vault *VaultCaller) DELEGATORFACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "DELEGATOR_FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DELEGATORFACTORY is a free data retrieval call binding the contract method 0x6da3e06d.
//
// Solidity: function DELEGATOR_FACTORY() view returns(address)
func (_Vault *VaultSession) DELEGATORFACTORY() (common.Address, error) {
	return _Vault.Contract.DELEGATORFACTORY(&_Vault.CallOpts)
}

// DELEGATORFACTORY is a free data retrieval call binding the contract method 0x6da3e06d.
//
// Solidity: function DELEGATOR_FACTORY() view returns(address)
func (_Vault *VaultCallerSession) DELEGATORFACTORY() (common.Address, error) {
	return _Vault.Contract.DELEGATORFACTORY(&_Vault.CallOpts)
}

// DEPOSITORWHITELISTROLE is a free data retrieval call binding the contract method 0x1b66c9e1.
//
// Solidity: function DEPOSITOR_WHITELIST_ROLE() view returns(bytes32)
func (_Vault *VaultCaller) DEPOSITORWHITELISTROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "DEPOSITOR_WHITELIST_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEPOSITORWHITELISTROLE is a free data retrieval call binding the contract method 0x1b66c9e1.
//
// Solidity: function DEPOSITOR_WHITELIST_ROLE() view returns(bytes32)
func (_Vault *VaultSession) DEPOSITORWHITELISTROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITORWHITELISTROLE(&_Vault.CallOpts)
}

// DEPOSITORWHITELISTROLE is a free data retrieval call binding the contract method 0x1b66c9e1.
//
// Solidity: function DEPOSITOR_WHITELIST_ROLE() view returns(bytes32)
func (_Vault *VaultCallerSession) DEPOSITORWHITELISTROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITORWHITELISTROLE(&_Vault.CallOpts)
}

// DEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0xa21a1df9.
//
// Solidity: function DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCaller) DEPOSITLIMITSETROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "DEPOSIT_LIMIT_SET_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0xa21a1df9.
//
// Solidity: function DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultSession) DEPOSITLIMITSETROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITLIMITSETROLE(&_Vault.CallOpts)
}

// DEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0xa21a1df9.
//
// Solidity: function DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCallerSession) DEPOSITLIMITSETROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITLIMITSETROLE(&_Vault.CallOpts)
}

// DEPOSITWHITELISTSETROLE is a free data retrieval call binding the contract method 0xdb388715.
//
// Solidity: function DEPOSIT_WHITELIST_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCaller) DEPOSITWHITELISTSETROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "DEPOSIT_WHITELIST_SET_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEPOSITWHITELISTSETROLE is a free data retrieval call binding the contract method 0xdb388715.
//
// Solidity: function DEPOSIT_WHITELIST_SET_ROLE() view returns(bytes32)
func (_Vault *VaultSession) DEPOSITWHITELISTSETROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITWHITELISTSETROLE(&_Vault.CallOpts)
}

// DEPOSITWHITELISTSETROLE is a free data retrieval call binding the contract method 0xdb388715.
//
// Solidity: function DEPOSIT_WHITELIST_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCallerSession) DEPOSITWHITELISTSETROLE() ([32]byte, error) {
	return _Vault.Contract.DEPOSITWHITELISTSETROLE(&_Vault.CallOpts)
}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Vault *VaultCaller) FACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Vault *VaultSession) FACTORY() (common.Address, error) {
	return _Vault.Contract.FACTORY(&_Vault.CallOpts)
}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Vault *VaultCallerSession) FACTORY() (common.Address, error) {
	return _Vault.Contract.FACTORY(&_Vault.CallOpts)
}

// ISDEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0x1415519b.
//
// Solidity: function IS_DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCaller) ISDEPOSITLIMITSETROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "IS_DEPOSIT_LIMIT_SET_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ISDEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0x1415519b.
//
// Solidity: function IS_DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultSession) ISDEPOSITLIMITSETROLE() ([32]byte, error) {
	return _Vault.Contract.ISDEPOSITLIMITSETROLE(&_Vault.CallOpts)
}

// ISDEPOSITLIMITSETROLE is a free data retrieval call binding the contract method 0x1415519b.
//
// Solidity: function IS_DEPOSIT_LIMIT_SET_ROLE() view returns(bytes32)
func (_Vault *VaultCallerSession) ISDEPOSITLIMITSETROLE() ([32]byte, error) {
	return _Vault.Contract.ISDEPOSITLIMITSETROLE(&_Vault.CallOpts)
}

// SLASHERFACTORY is a free data retrieval call binding the contract method 0x87df0788.
//
// Solidity: function SLASHER_FACTORY() view returns(address)
func (_Vault *VaultCaller) SLASHERFACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "SLASHER_FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SLASHERFACTORY is a free data retrieval call binding the contract method 0x87df0788.
//
// Solidity: function SLASHER_FACTORY() view returns(address)
func (_Vault *VaultSession) SLASHERFACTORY() (common.Address, error) {
	return _Vault.Contract.SLASHERFACTORY(&_Vault.CallOpts)
}

// SLASHERFACTORY is a free data retrieval call binding the contract method 0x87df0788.
//
// Solidity: function SLASHER_FACTORY() view returns(address)
func (_Vault *VaultCallerSession) SLASHERFACTORY() (common.Address, error) {
	return _Vault.Contract.SLASHERFACTORY(&_Vault.CallOpts)
}

// ActiveBalanceOf is a free data retrieval call binding the contract method 0x59f769a9.
//
// Solidity: function activeBalanceOf(address account) view returns(uint256)
func (_Vault *VaultCaller) ActiveBalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeBalanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveBalanceOf is a free data retrieval call binding the contract method 0x59f769a9.
//
// Solidity: function activeBalanceOf(address account) view returns(uint256)
func (_Vault *VaultSession) ActiveBalanceOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.ActiveBalanceOf(&_Vault.CallOpts, account)
}

// ActiveBalanceOf is a free data retrieval call binding the contract method 0x59f769a9.
//
// Solidity: function activeBalanceOf(address account) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveBalanceOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.ActiveBalanceOf(&_Vault.CallOpts, account)
}

// ActiveBalanceOfAt is a free data retrieval call binding the contract method 0xefb559d6.
//
// Solidity: function activeBalanceOfAt(address account, uint48 timestamp, bytes hints) view returns(uint256)
func (_Vault *VaultCaller) ActiveBalanceOfAt(opts *bind.CallOpts, account common.Address, timestamp *big.Int, hints []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeBalanceOfAt", account, timestamp, hints)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveBalanceOfAt is a free data retrieval call binding the contract method 0xefb559d6.
//
// Solidity: function activeBalanceOfAt(address account, uint48 timestamp, bytes hints) view returns(uint256)
func (_Vault *VaultSession) ActiveBalanceOfAt(account common.Address, timestamp *big.Int, hints []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveBalanceOfAt(&_Vault.CallOpts, account, timestamp, hints)
}

// ActiveBalanceOfAt is a free data retrieval call binding the contract method 0xefb559d6.
//
// Solidity: function activeBalanceOfAt(address account, uint48 timestamp, bytes hints) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveBalanceOfAt(account common.Address, timestamp *big.Int, hints []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveBalanceOfAt(&_Vault.CallOpts, account, timestamp, hints)
}

// ActiveShares is a free data retrieval call binding the contract method 0xbfefcd7b.
//
// Solidity: function activeShares() view returns(uint256)
func (_Vault *VaultCaller) ActiveShares(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeShares")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveShares is a free data retrieval call binding the contract method 0xbfefcd7b.
//
// Solidity: function activeShares() view returns(uint256)
func (_Vault *VaultSession) ActiveShares() (*big.Int, error) {
	return _Vault.Contract.ActiveShares(&_Vault.CallOpts)
}

// ActiveShares is a free data retrieval call binding the contract method 0xbfefcd7b.
//
// Solidity: function activeShares() view returns(uint256)
func (_Vault *VaultCallerSession) ActiveShares() (*big.Int, error) {
	return _Vault.Contract.ActiveShares(&_Vault.CallOpts)
}

// ActiveSharesAt is a free data retrieval call binding the contract method 0x50f22068.
//
// Solidity: function activeSharesAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCaller) ActiveSharesAt(opts *bind.CallOpts, timestamp *big.Int, hint []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeSharesAt", timestamp, hint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveSharesAt is a free data retrieval call binding the contract method 0x50f22068.
//
// Solidity: function activeSharesAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultSession) ActiveSharesAt(timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesAt(&_Vault.CallOpts, timestamp, hint)
}

// ActiveSharesAt is a free data retrieval call binding the contract method 0x50f22068.
//
// Solidity: function activeSharesAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveSharesAt(timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesAt(&_Vault.CallOpts, timestamp, hint)
}

// ActiveSharesOf is a free data retrieval call binding the contract method 0x9d66201b.
//
// Solidity: function activeSharesOf(address account) view returns(uint256)
func (_Vault *VaultCaller) ActiveSharesOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeSharesOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveSharesOf is a free data retrieval call binding the contract method 0x9d66201b.
//
// Solidity: function activeSharesOf(address account) view returns(uint256)
func (_Vault *VaultSession) ActiveSharesOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesOf(&_Vault.CallOpts, account)
}

// ActiveSharesOf is a free data retrieval call binding the contract method 0x9d66201b.
//
// Solidity: function activeSharesOf(address account) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveSharesOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesOf(&_Vault.CallOpts, account)
}

// ActiveSharesOfAt is a free data retrieval call binding the contract method 0x2d73c69c.
//
// Solidity: function activeSharesOfAt(address account, uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCaller) ActiveSharesOfAt(opts *bind.CallOpts, account common.Address, timestamp *big.Int, hint []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeSharesOfAt", account, timestamp, hint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveSharesOfAt is a free data retrieval call binding the contract method 0x2d73c69c.
//
// Solidity: function activeSharesOfAt(address account, uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultSession) ActiveSharesOfAt(account common.Address, timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesOfAt(&_Vault.CallOpts, account, timestamp, hint)
}

// ActiveSharesOfAt is a free data retrieval call binding the contract method 0x2d73c69c.
//
// Solidity: function activeSharesOfAt(address account, uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveSharesOfAt(account common.Address, timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveSharesOfAt(&_Vault.CallOpts, account, timestamp, hint)
}

// ActiveStake is a free data retrieval call binding the contract method 0xbd49c35f.
//
// Solidity: function activeStake() view returns(uint256)
func (_Vault *VaultCaller) ActiveStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveStake is a free data retrieval call binding the contract method 0xbd49c35f.
//
// Solidity: function activeStake() view returns(uint256)
func (_Vault *VaultSession) ActiveStake() (*big.Int, error) {
	return _Vault.Contract.ActiveStake(&_Vault.CallOpts)
}

// ActiveStake is a free data retrieval call binding the contract method 0xbd49c35f.
//
// Solidity: function activeStake() view returns(uint256)
func (_Vault *VaultCallerSession) ActiveStake() (*big.Int, error) {
	return _Vault.Contract.ActiveStake(&_Vault.CallOpts)
}

// ActiveStakeAt is a free data retrieval call binding the contract method 0x810da75d.
//
// Solidity: function activeStakeAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCaller) ActiveStakeAt(opts *bind.CallOpts, timestamp *big.Int, hint []byte) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "activeStakeAt", timestamp, hint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveStakeAt is a free data retrieval call binding the contract method 0x810da75d.
//
// Solidity: function activeStakeAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultSession) ActiveStakeAt(timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveStakeAt(&_Vault.CallOpts, timestamp, hint)
}

// ActiveStakeAt is a free data retrieval call binding the contract method 0x810da75d.
//
// Solidity: function activeStakeAt(uint48 timestamp, bytes hint) view returns(uint256)
func (_Vault *VaultCallerSession) ActiveStakeAt(timestamp *big.Int, hint []byte) (*big.Int, error) {
	return _Vault.Contract.ActiveStakeAt(&_Vault.CallOpts, timestamp, hint)
}

// Burner is a free data retrieval call binding the contract method 0x27810b6e.
//
// Solidity: function burner() view returns(address)
func (_Vault *VaultCaller) Burner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "burner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Burner is a free data retrieval call binding the contract method 0x27810b6e.
//
// Solidity: function burner() view returns(address)
func (_Vault *VaultSession) Burner() (common.Address, error) {
	return _Vault.Contract.Burner(&_Vault.CallOpts)
}

// Burner is a free data retrieval call binding the contract method 0x27810b6e.
//
// Solidity: function burner() view returns(address)
func (_Vault *VaultCallerSession) Burner() (common.Address, error) {
	return _Vault.Contract.Burner(&_Vault.CallOpts)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Vault *VaultCaller) Collateral(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "collateral")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Vault *VaultSession) Collateral() (common.Address, error) {
	return _Vault.Contract.Collateral(&_Vault.CallOpts)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_Vault *VaultCallerSession) Collateral() (common.Address, error) {
	return _Vault.Contract.Collateral(&_Vault.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Vault *VaultCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Vault *VaultSession) CurrentEpoch() (*big.Int, error) {
	return _Vault.Contract.CurrentEpoch(&_Vault.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Vault *VaultCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Vault.Contract.CurrentEpoch(&_Vault.CallOpts)
}

// CurrentEpochStart is a free data retrieval call binding the contract method 0x61a8c8c4.
//
// Solidity: function currentEpochStart() view returns(uint48)
func (_Vault *VaultCaller) CurrentEpochStart(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "currentEpochStart")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpochStart is a free data retrieval call binding the contract method 0x61a8c8c4.
//
// Solidity: function currentEpochStart() view returns(uint48)
func (_Vault *VaultSession) CurrentEpochStart() (*big.Int, error) {
	return _Vault.Contract.CurrentEpochStart(&_Vault.CallOpts)
}

// CurrentEpochStart is a free data retrieval call binding the contract method 0x61a8c8c4.
//
// Solidity: function currentEpochStart() view returns(uint48)
func (_Vault *VaultCallerSession) CurrentEpochStart() (*big.Int, error) {
	return _Vault.Contract.CurrentEpochStart(&_Vault.CallOpts)
}

// Delegator is a free data retrieval call binding the contract method 0xce9b7930.
//
// Solidity: function delegator() view returns(address)
func (_Vault *VaultCaller) Delegator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "delegator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Delegator is a free data retrieval call binding the contract method 0xce9b7930.
//
// Solidity: function delegator() view returns(address)
func (_Vault *VaultSession) Delegator() (common.Address, error) {
	return _Vault.Contract.Delegator(&_Vault.CallOpts)
}

// Delegator is a free data retrieval call binding the contract method 0xce9b7930.
//
// Solidity: function delegator() view returns(address)
func (_Vault *VaultCallerSession) Delegator() (common.Address, error) {
	return _Vault.Contract.Delegator(&_Vault.CallOpts)
}

// DepositLimit is a free data retrieval call binding the contract method 0xecf70858.
//
// Solidity: function depositLimit() view returns(uint256)
func (_Vault *VaultCaller) DepositLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "depositLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositLimit is a free data retrieval call binding the contract method 0xecf70858.
//
// Solidity: function depositLimit() view returns(uint256)
func (_Vault *VaultSession) DepositLimit() (*big.Int, error) {
	return _Vault.Contract.DepositLimit(&_Vault.CallOpts)
}

// DepositLimit is a free data retrieval call binding the contract method 0xecf70858.
//
// Solidity: function depositLimit() view returns(uint256)
func (_Vault *VaultCallerSession) DepositLimit() (*big.Int, error) {
	return _Vault.Contract.DepositLimit(&_Vault.CallOpts)
}

// DepositWhitelist is a free data retrieval call binding the contract method 0x48d3b775.
//
// Solidity: function depositWhitelist() view returns(bool)
func (_Vault *VaultCaller) DepositWhitelist(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "depositWhitelist")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// DepositWhitelist is a free data retrieval call binding the contract method 0x48d3b775.
//
// Solidity: function depositWhitelist() view returns(bool)
func (_Vault *VaultSession) DepositWhitelist() (bool, error) {
	return _Vault.Contract.DepositWhitelist(&_Vault.CallOpts)
}

// DepositWhitelist is a free data retrieval call binding the contract method 0x48d3b775.
//
// Solidity: function depositWhitelist() view returns(bool)
func (_Vault *VaultCallerSession) DepositWhitelist() (bool, error) {
	return _Vault.Contract.DepositWhitelist(&_Vault.CallOpts)
}

// EpochAt is a free data retrieval call binding the contract method 0x7953b33b.
//
// Solidity: function epochAt(uint48 timestamp) view returns(uint256)
func (_Vault *VaultCaller) EpochAt(opts *bind.CallOpts, timestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "epochAt", timestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochAt is a free data retrieval call binding the contract method 0x7953b33b.
//
// Solidity: function epochAt(uint48 timestamp) view returns(uint256)
func (_Vault *VaultSession) EpochAt(timestamp *big.Int) (*big.Int, error) {
	return _Vault.Contract.EpochAt(&_Vault.CallOpts, timestamp)
}

// EpochAt is a free data retrieval call binding the contract method 0x7953b33b.
//
// Solidity: function epochAt(uint48 timestamp) view returns(uint256)
func (_Vault *VaultCallerSession) EpochAt(timestamp *big.Int) (*big.Int, error) {
	return _Vault.Contract.EpochAt(&_Vault.CallOpts, timestamp)
}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint48)
func (_Vault *VaultCaller) EpochDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "epochDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint48)
func (_Vault *VaultSession) EpochDuration() (*big.Int, error) {
	return _Vault.Contract.EpochDuration(&_Vault.CallOpts)
}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint48)
func (_Vault *VaultCallerSession) EpochDuration() (*big.Int, error) {
	return _Vault.Contract.EpochDuration(&_Vault.CallOpts)
}

// EpochDurationInit is a free data retrieval call binding the contract method 0x46361671.
//
// Solidity: function epochDurationInit() view returns(uint48)
func (_Vault *VaultCaller) EpochDurationInit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "epochDurationInit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochDurationInit is a free data retrieval call binding the contract method 0x46361671.
//
// Solidity: function epochDurationInit() view returns(uint48)
func (_Vault *VaultSession) EpochDurationInit() (*big.Int, error) {
	return _Vault.Contract.EpochDurationInit(&_Vault.CallOpts)
}

// EpochDurationInit is a free data retrieval call binding the contract method 0x46361671.
//
// Solidity: function epochDurationInit() view returns(uint48)
func (_Vault *VaultCallerSession) EpochDurationInit() (*big.Int, error) {
	return _Vault.Contract.EpochDurationInit(&_Vault.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Vault *VaultCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Vault *VaultSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Vault.Contract.GetRoleAdmin(&_Vault.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Vault *VaultCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Vault.Contract.GetRoleAdmin(&_Vault.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Vault *VaultCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Vault *VaultSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Vault.Contract.HasRole(&_Vault.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Vault *VaultCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Vault.Contract.HasRole(&_Vault.CallOpts, role, account)
}

// IsDelegatorInitialized is a free data retrieval call binding the contract method 0x50861adc.
//
// Solidity: function isDelegatorInitialized() view returns(bool)
func (_Vault *VaultCaller) IsDelegatorInitialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isDelegatorInitialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDelegatorInitialized is a free data retrieval call binding the contract method 0x50861adc.
//
// Solidity: function isDelegatorInitialized() view returns(bool)
func (_Vault *VaultSession) IsDelegatorInitialized() (bool, error) {
	return _Vault.Contract.IsDelegatorInitialized(&_Vault.CallOpts)
}

// IsDelegatorInitialized is a free data retrieval call binding the contract method 0x50861adc.
//
// Solidity: function isDelegatorInitialized() view returns(bool)
func (_Vault *VaultCallerSession) IsDelegatorInitialized() (bool, error) {
	return _Vault.Contract.IsDelegatorInitialized(&_Vault.CallOpts)
}

// IsDepositLimit is a free data retrieval call binding the contract method 0xa1b12202.
//
// Solidity: function isDepositLimit() view returns(bool)
func (_Vault *VaultCaller) IsDepositLimit(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isDepositLimit")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDepositLimit is a free data retrieval call binding the contract method 0xa1b12202.
//
// Solidity: function isDepositLimit() view returns(bool)
func (_Vault *VaultSession) IsDepositLimit() (bool, error) {
	return _Vault.Contract.IsDepositLimit(&_Vault.CallOpts)
}

// IsDepositLimit is a free data retrieval call binding the contract method 0xa1b12202.
//
// Solidity: function isDepositLimit() view returns(bool)
func (_Vault *VaultCallerSession) IsDepositLimit() (bool, error) {
	return _Vault.Contract.IsDepositLimit(&_Vault.CallOpts)
}

// IsDepositorWhitelisted is a free data retrieval call binding the contract method 0x794b15b7.
//
// Solidity: function isDepositorWhitelisted(address account) view returns(bool value)
func (_Vault *VaultCaller) IsDepositorWhitelisted(opts *bind.CallOpts, account common.Address) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isDepositorWhitelisted", account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDepositorWhitelisted is a free data retrieval call binding the contract method 0x794b15b7.
//
// Solidity: function isDepositorWhitelisted(address account) view returns(bool value)
func (_Vault *VaultSession) IsDepositorWhitelisted(account common.Address) (bool, error) {
	return _Vault.Contract.IsDepositorWhitelisted(&_Vault.CallOpts, account)
}

// IsDepositorWhitelisted is a free data retrieval call binding the contract method 0x794b15b7.
//
// Solidity: function isDepositorWhitelisted(address account) view returns(bool value)
func (_Vault *VaultCallerSession) IsDepositorWhitelisted(account common.Address) (bool, error) {
	return _Vault.Contract.IsDepositorWhitelisted(&_Vault.CallOpts, account)
}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(bool)
func (_Vault *VaultCaller) IsInitialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isInitialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(bool)
func (_Vault *VaultSession) IsInitialized() (bool, error) {
	return _Vault.Contract.IsInitialized(&_Vault.CallOpts)
}

// IsInitialized is a free data retrieval call binding the contract method 0x392e53cd.
//
// Solidity: function isInitialized() view returns(bool)
func (_Vault *VaultCallerSession) IsInitialized() (bool, error) {
	return _Vault.Contract.IsInitialized(&_Vault.CallOpts)
}

// IsSlasherInitialized is a free data retrieval call binding the contract method 0x6ec1e3f8.
//
// Solidity: function isSlasherInitialized() view returns(bool)
func (_Vault *VaultCaller) IsSlasherInitialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isSlasherInitialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSlasherInitialized is a free data retrieval call binding the contract method 0x6ec1e3f8.
//
// Solidity: function isSlasherInitialized() view returns(bool)
func (_Vault *VaultSession) IsSlasherInitialized() (bool, error) {
	return _Vault.Contract.IsSlasherInitialized(&_Vault.CallOpts)
}

// IsSlasherInitialized is a free data retrieval call binding the contract method 0x6ec1e3f8.
//
// Solidity: function isSlasherInitialized() view returns(bool)
func (_Vault *VaultCallerSession) IsSlasherInitialized() (bool, error) {
	return _Vault.Contract.IsSlasherInitialized(&_Vault.CallOpts)
}

// IsWithdrawalsClaimed is a free data retrieval call binding the contract method 0xa5d03223.
//
// Solidity: function isWithdrawalsClaimed(uint256 epoch, address account) view returns(bool value)
func (_Vault *VaultCaller) IsWithdrawalsClaimed(opts *bind.CallOpts, epoch *big.Int, account common.Address) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isWithdrawalsClaimed", epoch, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWithdrawalsClaimed is a free data retrieval call binding the contract method 0xa5d03223.
//
// Solidity: function isWithdrawalsClaimed(uint256 epoch, address account) view returns(bool value)
func (_Vault *VaultSession) IsWithdrawalsClaimed(epoch *big.Int, account common.Address) (bool, error) {
	return _Vault.Contract.IsWithdrawalsClaimed(&_Vault.CallOpts, epoch, account)
}

// IsWithdrawalsClaimed is a free data retrieval call binding the contract method 0xa5d03223.
//
// Solidity: function isWithdrawalsClaimed(uint256 epoch, address account) view returns(bool value)
func (_Vault *VaultCallerSession) IsWithdrawalsClaimed(epoch *big.Int, account common.Address) (bool, error) {
	return _Vault.Contract.IsWithdrawalsClaimed(&_Vault.CallOpts, epoch, account)
}

// NextEpochStart is a free data retrieval call binding the contract method 0x73790ab3.
//
// Solidity: function nextEpochStart() view returns(uint48)
func (_Vault *VaultCaller) NextEpochStart(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "nextEpochStart")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextEpochStart is a free data retrieval call binding the contract method 0x73790ab3.
//
// Solidity: function nextEpochStart() view returns(uint48)
func (_Vault *VaultSession) NextEpochStart() (*big.Int, error) {
	return _Vault.Contract.NextEpochStart(&_Vault.CallOpts)
}

// NextEpochStart is a free data retrieval call binding the contract method 0x73790ab3.
//
// Solidity: function nextEpochStart() view returns(uint48)
func (_Vault *VaultCallerSession) NextEpochStart() (*big.Int, error) {
	return _Vault.Contract.NextEpochStart(&_Vault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultSession) Owner() (common.Address, error) {
	return _Vault.Contract.Owner(&_Vault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultCallerSession) Owner() (common.Address, error) {
	return _Vault.Contract.Owner(&_Vault.CallOpts)
}

// PreviousEpochStart is a free data retrieval call binding the contract method 0x281f5752.
//
// Solidity: function previousEpochStart() view returns(uint48)
func (_Vault *VaultCaller) PreviousEpochStart(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "previousEpochStart")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviousEpochStart is a free data retrieval call binding the contract method 0x281f5752.
//
// Solidity: function previousEpochStart() view returns(uint48)
func (_Vault *VaultSession) PreviousEpochStart() (*big.Int, error) {
	return _Vault.Contract.PreviousEpochStart(&_Vault.CallOpts)
}

// PreviousEpochStart is a free data retrieval call binding the contract method 0x281f5752.
//
// Solidity: function previousEpochStart() view returns(uint48)
func (_Vault *VaultCallerSession) PreviousEpochStart() (*big.Int, error) {
	return _Vault.Contract.PreviousEpochStart(&_Vault.CallOpts)
}

// SlashableBalanceOf is a free data retrieval call binding the contract method 0xc31e8dd7.
//
// Solidity: function slashableBalanceOf(address account) view returns(uint256)
func (_Vault *VaultCaller) SlashableBalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "slashableBalanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashableBalanceOf is a free data retrieval call binding the contract method 0xc31e8dd7.
//
// Solidity: function slashableBalanceOf(address account) view returns(uint256)
func (_Vault *VaultSession) SlashableBalanceOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.SlashableBalanceOf(&_Vault.CallOpts, account)
}

// SlashableBalanceOf is a free data retrieval call binding the contract method 0xc31e8dd7.
//
// Solidity: function slashableBalanceOf(address account) view returns(uint256)
func (_Vault *VaultCallerSession) SlashableBalanceOf(account common.Address) (*big.Int, error) {
	return _Vault.Contract.SlashableBalanceOf(&_Vault.CallOpts, account)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_Vault *VaultCaller) Slasher(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "slasher")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_Vault *VaultSession) Slasher() (common.Address, error) {
	return _Vault.Contract.Slasher(&_Vault.CallOpts)
}

// Slasher is a free data retrieval call binding the contract method 0xb1344271.
//
// Solidity: function slasher() view returns(address)
func (_Vault *VaultCallerSession) Slasher() (common.Address, error) {
	return _Vault.Contract.Slasher(&_Vault.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Vault *VaultCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Vault *VaultSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Vault.Contract.SupportsInterface(&_Vault.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Vault *VaultCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Vault.Contract.SupportsInterface(&_Vault.CallOpts, interfaceId)
}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_Vault *VaultCaller) TotalStake(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "totalStake")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_Vault *VaultSession) TotalStake() (*big.Int, error) {
	return _Vault.Contract.TotalStake(&_Vault.CallOpts)
}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_Vault *VaultCallerSession) TotalStake() (*big.Int, error) {
	return _Vault.Contract.TotalStake(&_Vault.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint64)
func (_Vault *VaultCaller) Version(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint64)
func (_Vault *VaultSession) Version() (uint64, error) {
	return _Vault.Contract.Version(&_Vault.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint64)
func (_Vault *VaultCallerSession) Version() (uint64, error) {
	return _Vault.Contract.Version(&_Vault.CallOpts)
}

// WithdrawalShares is a free data retrieval call binding the contract method 0xafba70ad.
//
// Solidity: function withdrawalShares(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultCaller) WithdrawalShares(opts *bind.CallOpts, epoch *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "withdrawalShares", epoch)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalShares is a free data retrieval call binding the contract method 0xafba70ad.
//
// Solidity: function withdrawalShares(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultSession) WithdrawalShares(epoch *big.Int) (*big.Int, error) {
	return _Vault.Contract.WithdrawalShares(&_Vault.CallOpts, epoch)
}

// WithdrawalShares is a free data retrieval call binding the contract method 0xafba70ad.
//
// Solidity: function withdrawalShares(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultCallerSession) WithdrawalShares(epoch *big.Int) (*big.Int, error) {
	return _Vault.Contract.WithdrawalShares(&_Vault.CallOpts, epoch)
}

// WithdrawalSharesOf is a free data retrieval call binding the contract method 0xa3b54172.
//
// Solidity: function withdrawalSharesOf(uint256 epoch, address account) view returns(uint256 amount)
func (_Vault *VaultCaller) WithdrawalSharesOf(opts *bind.CallOpts, epoch *big.Int, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "withdrawalSharesOf", epoch, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalSharesOf is a free data retrieval call binding the contract method 0xa3b54172.
//
// Solidity: function withdrawalSharesOf(uint256 epoch, address account) view returns(uint256 amount)
func (_Vault *VaultSession) WithdrawalSharesOf(epoch *big.Int, account common.Address) (*big.Int, error) {
	return _Vault.Contract.WithdrawalSharesOf(&_Vault.CallOpts, epoch, account)
}

// WithdrawalSharesOf is a free data retrieval call binding the contract method 0xa3b54172.
//
// Solidity: function withdrawalSharesOf(uint256 epoch, address account) view returns(uint256 amount)
func (_Vault *VaultCallerSession) WithdrawalSharesOf(epoch *big.Int, account common.Address) (*big.Int, error) {
	return _Vault.Contract.WithdrawalSharesOf(&_Vault.CallOpts, epoch, account)
}

// Withdrawals is a free data retrieval call binding the contract method 0x5cc07076.
//
// Solidity: function withdrawals(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultCaller) Withdrawals(opts *bind.CallOpts, epoch *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "withdrawals", epoch)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Withdrawals is a free data retrieval call binding the contract method 0x5cc07076.
//
// Solidity: function withdrawals(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultSession) Withdrawals(epoch *big.Int) (*big.Int, error) {
	return _Vault.Contract.Withdrawals(&_Vault.CallOpts, epoch)
}

// Withdrawals is a free data retrieval call binding the contract method 0x5cc07076.
//
// Solidity: function withdrawals(uint256 epoch) view returns(uint256 amount)
func (_Vault *VaultCallerSession) Withdrawals(epoch *big.Int) (*big.Int, error) {
	return _Vault.Contract.Withdrawals(&_Vault.CallOpts, epoch)
}

// WithdrawalsOf is a free data retrieval call binding the contract method 0xf5e7ee0f.
//
// Solidity: function withdrawalsOf(uint256 epoch, address account) view returns(uint256)
func (_Vault *VaultCaller) WithdrawalsOf(opts *bind.CallOpts, epoch *big.Int, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "withdrawalsOf", epoch, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalsOf is a free data retrieval call binding the contract method 0xf5e7ee0f.
//
// Solidity: function withdrawalsOf(uint256 epoch, address account) view returns(uint256)
func (_Vault *VaultSession) WithdrawalsOf(epoch *big.Int, account common.Address) (*big.Int, error) {
	return _Vault.Contract.WithdrawalsOf(&_Vault.CallOpts, epoch, account)
}

// WithdrawalsOf is a free data retrieval call binding the contract method 0xf5e7ee0f.
//
// Solidity: function withdrawalsOf(uint256 epoch, address account) view returns(uint256)
func (_Vault *VaultCallerSession) WithdrawalsOf(epoch *big.Int, account common.Address) (*big.Int, error) {
	return _Vault.Contract.WithdrawalsOf(&_Vault.CallOpts, epoch, account)
}

// Claim is a paid mutator transaction binding the contract method 0xaad3ec96.
//
// Solidity: function claim(address recipient, uint256 epoch) returns(uint256 amount)
func (_Vault *VaultTransactor) Claim(opts *bind.TransactOpts, recipient common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "claim", recipient, epoch)
}

// Claim is a paid mutator transaction binding the contract method 0xaad3ec96.
//
// Solidity: function claim(address recipient, uint256 epoch) returns(uint256 amount)
func (_Vault *VaultSession) Claim(recipient common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Claim(&_Vault.TransactOpts, recipient, epoch)
}

// Claim is a paid mutator transaction binding the contract method 0xaad3ec96.
//
// Solidity: function claim(address recipient, uint256 epoch) returns(uint256 amount)
func (_Vault *VaultTransactorSession) Claim(recipient common.Address, epoch *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Claim(&_Vault.TransactOpts, recipient, epoch)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0x7c04c80a.
//
// Solidity: function claimBatch(address recipient, uint256[] epochs) returns(uint256 amount)
func (_Vault *VaultTransactor) ClaimBatch(opts *bind.TransactOpts, recipient common.Address, epochs []*big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "claimBatch", recipient, epochs)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0x7c04c80a.
//
// Solidity: function claimBatch(address recipient, uint256[] epochs) returns(uint256 amount)
func (_Vault *VaultSession) ClaimBatch(recipient common.Address, epochs []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.ClaimBatch(&_Vault.TransactOpts, recipient, epochs)
}

// ClaimBatch is a paid mutator transaction binding the contract method 0x7c04c80a.
//
// Solidity: function claimBatch(address recipient, uint256[] epochs) returns(uint256 amount)
func (_Vault *VaultTransactorSession) ClaimBatch(recipient common.Address, epochs []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.ClaimBatch(&_Vault.TransactOpts, recipient, epochs)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address onBehalfOf, uint256 amount) returns(uint256 depositedAmount, uint256 mintedShares)
func (_Vault *VaultTransactor) Deposit(opts *bind.TransactOpts, onBehalfOf common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "deposit", onBehalfOf, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address onBehalfOf, uint256 amount) returns(uint256 depositedAmount, uint256 mintedShares)
func (_Vault *VaultSession) Deposit(onBehalfOf common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Deposit(&_Vault.TransactOpts, onBehalfOf, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address onBehalfOf, uint256 amount) returns(uint256 depositedAmount, uint256 mintedShares)
func (_Vault *VaultTransactorSession) Deposit(onBehalfOf common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Deposit(&_Vault.TransactOpts, onBehalfOf, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Vault *VaultTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Vault *VaultSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.Contract.GrantRole(&_Vault.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Vault *VaultTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.Contract.GrantRole(&_Vault.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x57ec83cc.
//
// Solidity: function initialize(uint64 initialVersion, address owner_, bytes data) returns()
func (_Vault *VaultTransactor) Initialize(opts *bind.TransactOpts, initialVersion uint64, owner_ common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "initialize", initialVersion, owner_, data)
}

// Initialize is a paid mutator transaction binding the contract method 0x57ec83cc.
//
// Solidity: function initialize(uint64 initialVersion, address owner_, bytes data) returns()
func (_Vault *VaultSession) Initialize(initialVersion uint64, owner_ common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.Initialize(&_Vault.TransactOpts, initialVersion, owner_, data)
}

// Initialize is a paid mutator transaction binding the contract method 0x57ec83cc.
//
// Solidity: function initialize(uint64 initialVersion, address owner_, bytes data) returns()
func (_Vault *VaultTransactorSession) Initialize(initialVersion uint64, owner_ common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.Initialize(&_Vault.TransactOpts, initialVersion, owner_, data)
}

// Migrate is a paid mutator transaction binding the contract method 0x2abe3048.
//
// Solidity: function migrate(uint64 newVersion, bytes data) returns()
func (_Vault *VaultTransactor) Migrate(opts *bind.TransactOpts, newVersion uint64, data []byte) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "migrate", newVersion, data)
}

// Migrate is a paid mutator transaction binding the contract method 0x2abe3048.
//
// Solidity: function migrate(uint64 newVersion, bytes data) returns()
func (_Vault *VaultSession) Migrate(newVersion uint64, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.Migrate(&_Vault.TransactOpts, newVersion, data)
}

// Migrate is a paid mutator transaction binding the contract method 0x2abe3048.
//
// Solidity: function migrate(uint64 newVersion, bytes data) returns()
func (_Vault *VaultTransactorSession) Migrate(newVersion uint64, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.Migrate(&_Vault.TransactOpts, newVersion, data)
}

// OnSlash is a paid mutator transaction binding the contract method 0x7278e31c.
//
// Solidity: function onSlash(uint256 amount, uint48 captureTimestamp) returns(uint256 slashedAmount)
func (_Vault *VaultTransactor) OnSlash(opts *bind.TransactOpts, amount *big.Int, captureTimestamp *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "onSlash", amount, captureTimestamp)
}

// OnSlash is a paid mutator transaction binding the contract method 0x7278e31c.
//
// Solidity: function onSlash(uint256 amount, uint48 captureTimestamp) returns(uint256 slashedAmount)
func (_Vault *VaultSession) OnSlash(amount *big.Int, captureTimestamp *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.OnSlash(&_Vault.TransactOpts, amount, captureTimestamp)
}

// OnSlash is a paid mutator transaction binding the contract method 0x7278e31c.
//
// Solidity: function onSlash(uint256 amount, uint48 captureTimestamp) returns(uint256 slashedAmount)
func (_Vault *VaultTransactorSession) OnSlash(amount *big.Int, captureTimestamp *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.OnSlash(&_Vault.TransactOpts, amount, captureTimestamp)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address claimer, uint256 shares) returns(uint256 withdrawnAssets, uint256 mintedShares)
func (_Vault *VaultTransactor) Redeem(opts *bind.TransactOpts, claimer common.Address, shares *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "redeem", claimer, shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address claimer, uint256 shares) returns(uint256 withdrawnAssets, uint256 mintedShares)
func (_Vault *VaultSession) Redeem(claimer common.Address, shares *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Redeem(&_Vault.TransactOpts, claimer, shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address claimer, uint256 shares) returns(uint256 withdrawnAssets, uint256 mintedShares)
func (_Vault *VaultTransactorSession) Redeem(claimer common.Address, shares *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Redeem(&_Vault.TransactOpts, claimer, shares)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vault *VaultTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vault *VaultSession) RenounceOwnership() (*types.Transaction, error) {
	return _Vault.Contract.RenounceOwnership(&_Vault.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Vault *VaultTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Vault.Contract.RenounceOwnership(&_Vault.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Vault *VaultTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Vault *VaultSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RenounceRole(&_Vault.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_Vault *VaultTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RenounceRole(&_Vault.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Vault *VaultTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Vault *VaultSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RevokeRole(&_Vault.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Vault *VaultTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RevokeRole(&_Vault.TransactOpts, role, account)
}

// SetDelegator is a paid mutator transaction binding the contract method 0x83cd9cc3.
//
// Solidity: function setDelegator(address delegator_) returns()
func (_Vault *VaultTransactor) SetDelegator(opts *bind.TransactOpts, delegator_ common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setDelegator", delegator_)
}

// SetDelegator is a paid mutator transaction binding the contract method 0x83cd9cc3.
//
// Solidity: function setDelegator(address delegator_) returns()
func (_Vault *VaultSession) SetDelegator(delegator_ common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetDelegator(&_Vault.TransactOpts, delegator_)
}

// SetDelegator is a paid mutator transaction binding the contract method 0x83cd9cc3.
//
// Solidity: function setDelegator(address delegator_) returns()
func (_Vault *VaultTransactorSession) SetDelegator(delegator_ common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetDelegator(&_Vault.TransactOpts, delegator_)
}

// SetDepositLimit is a paid mutator transaction binding the contract method 0xbdc8144b.
//
// Solidity: function setDepositLimit(uint256 limit) returns()
func (_Vault *VaultTransactor) SetDepositLimit(opts *bind.TransactOpts, limit *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setDepositLimit", limit)
}

// SetDepositLimit is a paid mutator transaction binding the contract method 0xbdc8144b.
//
// Solidity: function setDepositLimit(uint256 limit) returns()
func (_Vault *VaultSession) SetDepositLimit(limit *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositLimit(&_Vault.TransactOpts, limit)
}

// SetDepositLimit is a paid mutator transaction binding the contract method 0xbdc8144b.
//
// Solidity: function setDepositLimit(uint256 limit) returns()
func (_Vault *VaultTransactorSession) SetDepositLimit(limit *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositLimit(&_Vault.TransactOpts, limit)
}

// SetDepositWhitelist is a paid mutator transaction binding the contract method 0x4105a7dd.
//
// Solidity: function setDepositWhitelist(bool status) returns()
func (_Vault *VaultTransactor) SetDepositWhitelist(opts *bind.TransactOpts, status bool) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setDepositWhitelist", status)
}

// SetDepositWhitelist is a paid mutator transaction binding the contract method 0x4105a7dd.
//
// Solidity: function setDepositWhitelist(bool status) returns()
func (_Vault *VaultSession) SetDepositWhitelist(status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositWhitelist(&_Vault.TransactOpts, status)
}

// SetDepositWhitelist is a paid mutator transaction binding the contract method 0x4105a7dd.
//
// Solidity: function setDepositWhitelist(bool status) returns()
func (_Vault *VaultTransactorSession) SetDepositWhitelist(status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositWhitelist(&_Vault.TransactOpts, status)
}

// SetDepositorWhitelistStatus is a paid mutator transaction binding the contract method 0xa2861466.
//
// Solidity: function setDepositorWhitelistStatus(address account, bool status) returns()
func (_Vault *VaultTransactor) SetDepositorWhitelistStatus(opts *bind.TransactOpts, account common.Address, status bool) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setDepositorWhitelistStatus", account, status)
}

// SetDepositorWhitelistStatus is a paid mutator transaction binding the contract method 0xa2861466.
//
// Solidity: function setDepositorWhitelistStatus(address account, bool status) returns()
func (_Vault *VaultSession) SetDepositorWhitelistStatus(account common.Address, status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositorWhitelistStatus(&_Vault.TransactOpts, account, status)
}

// SetDepositorWhitelistStatus is a paid mutator transaction binding the contract method 0xa2861466.
//
// Solidity: function setDepositorWhitelistStatus(address account, bool status) returns()
func (_Vault *VaultTransactorSession) SetDepositorWhitelistStatus(account common.Address, status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetDepositorWhitelistStatus(&_Vault.TransactOpts, account, status)
}

// SetIsDepositLimit is a paid mutator transaction binding the contract method 0x5346e34f.
//
// Solidity: function setIsDepositLimit(bool status) returns()
func (_Vault *VaultTransactor) SetIsDepositLimit(opts *bind.TransactOpts, status bool) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setIsDepositLimit", status)
}

// SetIsDepositLimit is a paid mutator transaction binding the contract method 0x5346e34f.
//
// Solidity: function setIsDepositLimit(bool status) returns()
func (_Vault *VaultSession) SetIsDepositLimit(status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetIsDepositLimit(&_Vault.TransactOpts, status)
}

// SetIsDepositLimit is a paid mutator transaction binding the contract method 0x5346e34f.
//
// Solidity: function setIsDepositLimit(bool status) returns()
func (_Vault *VaultTransactorSession) SetIsDepositLimit(status bool) (*types.Transaction, error) {
	return _Vault.Contract.SetIsDepositLimit(&_Vault.TransactOpts, status)
}

// SetSlasher is a paid mutator transaction binding the contract method 0xaabc2496.
//
// Solidity: function setSlasher(address slasher_) returns()
func (_Vault *VaultTransactor) SetSlasher(opts *bind.TransactOpts, slasher_ common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setSlasher", slasher_)
}

// SetSlasher is a paid mutator transaction binding the contract method 0xaabc2496.
//
// Solidity: function setSlasher(address slasher_) returns()
func (_Vault *VaultSession) SetSlasher(slasher_ common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetSlasher(&_Vault.TransactOpts, slasher_)
}

// SetSlasher is a paid mutator transaction binding the contract method 0xaabc2496.
//
// Solidity: function setSlasher(address slasher_) returns()
func (_Vault *VaultTransactorSession) SetSlasher(slasher_ common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetSlasher(&_Vault.TransactOpts, slasher_)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_Vault *VaultTransactor) StaticDelegateCall(opts *bind.TransactOpts, target common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "staticDelegateCall", target, data)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_Vault *VaultSession) StaticDelegateCall(target common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.StaticDelegateCall(&_Vault.TransactOpts, target, data)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_Vault *VaultTransactorSession) StaticDelegateCall(target common.Address, data []byte) (*types.Transaction, error) {
	return _Vault.Contract.StaticDelegateCall(&_Vault.TransactOpts, target, data)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vault *VaultTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vault *VaultSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vault.Contract.TransferOwnership(&_Vault.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vault *VaultTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vault.Contract.TransferOwnership(&_Vault.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address claimer, uint256 amount) returns(uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultTransactor) Withdraw(opts *bind.TransactOpts, claimer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "withdraw", claimer, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address claimer, uint256 amount) returns(uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultSession) Withdraw(claimer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Withdraw(&_Vault.TransactOpts, claimer, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address claimer, uint256 amount) returns(uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultTransactorSession) Withdraw(claimer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Withdraw(&_Vault.TransactOpts, claimer, amount)
}

// VaultClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the Vault contract.
type VaultClaimIterator struct {
	Event *VaultClaim // Event containing the contract specifics and raw log

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
func (it *VaultClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultClaim)
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
		it.Event = new(VaultClaim)
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
func (it *VaultClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultClaim represents a Claim event raised by the Vault contract.
type VaultClaim struct {
	Claimer   common.Address
	Recipient common.Address
	Epoch     *big.Int
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed claimer, address indexed recipient, uint256 epoch, uint256 amount)
func (_Vault *VaultFilterer) FilterClaim(opts *bind.FilterOpts, claimer []common.Address, recipient []common.Address) (*VaultClaimIterator, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Claim", claimerRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &VaultClaimIterator{contract: _Vault.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed claimer, address indexed recipient, uint256 epoch, uint256 amount)
func (_Vault *VaultFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *VaultClaim, claimer []common.Address, recipient []common.Address) (event.Subscription, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Claim", claimerRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultClaim)
				if err := _Vault.contract.UnpackLog(event, "Claim", log); err != nil {
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

// ParseClaim is a log parse operation binding the contract event 0x865ca08d59f5cb456e85cd2f7ef63664ea4f73327414e9d8152c4158b0e94645.
//
// Solidity: event Claim(address indexed claimer, address indexed recipient, uint256 epoch, uint256 amount)
func (_Vault *VaultFilterer) ParseClaim(log types.Log) (*VaultClaim, error) {
	event := new(VaultClaim)
	if err := _Vault.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultClaimBatchIterator is returned from FilterClaimBatch and is used to iterate over the raw logs and unpacked data for ClaimBatch events raised by the Vault contract.
type VaultClaimBatchIterator struct {
	Event *VaultClaimBatch // Event containing the contract specifics and raw log

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
func (it *VaultClaimBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultClaimBatch)
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
		it.Event = new(VaultClaimBatch)
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
func (it *VaultClaimBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultClaimBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultClaimBatch represents a ClaimBatch event raised by the Vault contract.
type VaultClaimBatch struct {
	Claimer   common.Address
	Recipient common.Address
	Epochs    []*big.Int
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimBatch is a free log retrieval operation binding the contract event 0x326b6aff1cd2fb1c8234de4f9dcfb9047c5c36eb9ef2eb34af5121e969a75d27.
//
// Solidity: event ClaimBatch(address indexed claimer, address indexed recipient, uint256[] epochs, uint256 amount)
func (_Vault *VaultFilterer) FilterClaimBatch(opts *bind.FilterOpts, claimer []common.Address, recipient []common.Address) (*VaultClaimBatchIterator, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "ClaimBatch", claimerRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &VaultClaimBatchIterator{contract: _Vault.contract, event: "ClaimBatch", logs: logs, sub: sub}, nil
}

// WatchClaimBatch is a free log subscription operation binding the contract event 0x326b6aff1cd2fb1c8234de4f9dcfb9047c5c36eb9ef2eb34af5121e969a75d27.
//
// Solidity: event ClaimBatch(address indexed claimer, address indexed recipient, uint256[] epochs, uint256 amount)
func (_Vault *VaultFilterer) WatchClaimBatch(opts *bind.WatchOpts, sink chan<- *VaultClaimBatch, claimer []common.Address, recipient []common.Address) (event.Subscription, error) {

	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "ClaimBatch", claimerRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultClaimBatch)
				if err := _Vault.contract.UnpackLog(event, "ClaimBatch", log); err != nil {
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

// ParseClaimBatch is a log parse operation binding the contract event 0x326b6aff1cd2fb1c8234de4f9dcfb9047c5c36eb9ef2eb34af5121e969a75d27.
//
// Solidity: event ClaimBatch(address indexed claimer, address indexed recipient, uint256[] epochs, uint256 amount)
func (_Vault *VaultFilterer) ParseClaimBatch(log types.Log) (*VaultClaimBatch, error) {
	event := new(VaultClaimBatch)
	if err := _Vault.contract.UnpackLog(event, "ClaimBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Vault contract.
type VaultDepositIterator struct {
	Event *VaultDeposit // Event containing the contract specifics and raw log

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
func (it *VaultDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultDeposit)
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
		it.Event = new(VaultDeposit)
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
func (it *VaultDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultDeposit represents a Deposit event raised by the Vault contract.
type VaultDeposit struct {
	Depositor  common.Address
	OnBehalfOf common.Address
	Amount     *big.Int
	Shares     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed depositor, address indexed onBehalfOf, uint256 amount, uint256 shares)
func (_Vault *VaultFilterer) FilterDeposit(opts *bind.FilterOpts, depositor []common.Address, onBehalfOf []common.Address) (*VaultDepositIterator, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Deposit", depositorRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return &VaultDepositIterator{contract: _Vault.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed depositor, address indexed onBehalfOf, uint256 amount, uint256 shares)
func (_Vault *VaultFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *VaultDeposit, depositor []common.Address, onBehalfOf []common.Address) (event.Subscription, error) {

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}
	var onBehalfOfRule []interface{}
	for _, onBehalfOfItem := range onBehalfOf {
		onBehalfOfRule = append(onBehalfOfRule, onBehalfOfItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Deposit", depositorRule, onBehalfOfRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultDeposit)
				if err := _Vault.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed depositor, address indexed onBehalfOf, uint256 amount, uint256 shares)
func (_Vault *VaultFilterer) ParseDeposit(log types.Log) (*VaultDeposit, error) {
	event := new(VaultDeposit)
	if err := _Vault.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Vault contract.
type VaultInitializedIterator struct {
	Event *VaultInitialized // Event containing the contract specifics and raw log

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
func (it *VaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultInitialized)
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
		it.Event = new(VaultInitialized)
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
func (it *VaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultInitialized represents a Initialized event raised by the Vault contract.
type VaultInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vault *VaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*VaultInitializedIterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &VaultInitializedIterator{contract: _Vault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Vault *VaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *VaultInitialized) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultInitialized)
				if err := _Vault.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Vault *VaultFilterer) ParseInitialized(log types.Log) (*VaultInitialized, error) {
	event := new(VaultInitialized)
	if err := _Vault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultOnSlashIterator is returned from FilterOnSlash and is used to iterate over the raw logs and unpacked data for OnSlash events raised by the Vault contract.
type VaultOnSlashIterator struct {
	Event *VaultOnSlash // Event containing the contract specifics and raw log

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
func (it *VaultOnSlashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultOnSlash)
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
		it.Event = new(VaultOnSlash)
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
func (it *VaultOnSlashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultOnSlashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultOnSlash represents a OnSlash event raised by the Vault contract.
type VaultOnSlash struct {
	Amount           *big.Int
	CaptureTimestamp *big.Int
	SlashedAmount    *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOnSlash is a free log retrieval operation binding the contract event 0xf9d090c096f71cd1659861a9ce5b6f384bceb4e2fa4e4a19edf6489a9b8d56c7.
//
// Solidity: event OnSlash(uint256 amount, uint48 captureTimestamp, uint256 slashedAmount)
func (_Vault *VaultFilterer) FilterOnSlash(opts *bind.FilterOpts) (*VaultOnSlashIterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "OnSlash")
	if err != nil {
		return nil, err
	}
	return &VaultOnSlashIterator{contract: _Vault.contract, event: "OnSlash", logs: logs, sub: sub}, nil
}

// WatchOnSlash is a free log subscription operation binding the contract event 0xf9d090c096f71cd1659861a9ce5b6f384bceb4e2fa4e4a19edf6489a9b8d56c7.
//
// Solidity: event OnSlash(uint256 amount, uint48 captureTimestamp, uint256 slashedAmount)
func (_Vault *VaultFilterer) WatchOnSlash(opts *bind.WatchOpts, sink chan<- *VaultOnSlash) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "OnSlash")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultOnSlash)
				if err := _Vault.contract.UnpackLog(event, "OnSlash", log); err != nil {
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

// ParseOnSlash is a log parse operation binding the contract event 0xf9d090c096f71cd1659861a9ce5b6f384bceb4e2fa4e4a19edf6489a9b8d56c7.
//
// Solidity: event OnSlash(uint256 amount, uint48 captureTimestamp, uint256 slashedAmount)
func (_Vault *VaultFilterer) ParseOnSlash(log types.Log) (*VaultOnSlash, error) {
	event := new(VaultOnSlash)
	if err := _Vault.contract.UnpackLog(event, "OnSlash", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Vault contract.
type VaultOwnershipTransferredIterator struct {
	Event *VaultOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VaultOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultOwnershipTransferred)
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
		it.Event = new(VaultOwnershipTransferred)
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
func (it *VaultOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultOwnershipTransferred represents a OwnershipTransferred event raised by the Vault contract.
type VaultOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vault *VaultFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VaultOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VaultOwnershipTransferredIterator{contract: _Vault.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vault *VaultFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VaultOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultOwnershipTransferred)
				if err := _Vault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Vault *VaultFilterer) ParseOwnershipTransferred(log types.Log) (*VaultOwnershipTransferred, error) {
	event := new(VaultOwnershipTransferred)
	if err := _Vault.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Vault contract.
type VaultRoleAdminChangedIterator struct {
	Event *VaultRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *VaultRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultRoleAdminChanged)
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
		it.Event = new(VaultRoleAdminChanged)
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
func (it *VaultRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultRoleAdminChanged represents a RoleAdminChanged event raised by the Vault contract.
type VaultRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Vault *VaultFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*VaultRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &VaultRoleAdminChangedIterator{contract: _Vault.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Vault *VaultFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *VaultRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultRoleAdminChanged)
				if err := _Vault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Vault *VaultFilterer) ParseRoleAdminChanged(log types.Log) (*VaultRoleAdminChanged, error) {
	event := new(VaultRoleAdminChanged)
	if err := _Vault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Vault contract.
type VaultRoleGrantedIterator struct {
	Event *VaultRoleGranted // Event containing the contract specifics and raw log

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
func (it *VaultRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultRoleGranted)
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
		it.Event = new(VaultRoleGranted)
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
func (it *VaultRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultRoleGranted represents a RoleGranted event raised by the Vault contract.
type VaultRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VaultRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VaultRoleGrantedIterator{contract: _Vault.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *VaultRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultRoleGranted)
				if err := _Vault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) ParseRoleGranted(log types.Log) (*VaultRoleGranted, error) {
	event := new(VaultRoleGranted)
	if err := _Vault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Vault contract.
type VaultRoleRevokedIterator struct {
	Event *VaultRoleRevoked // Event containing the contract specifics and raw log

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
func (it *VaultRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultRoleRevoked)
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
		it.Event = new(VaultRoleRevoked)
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
func (it *VaultRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultRoleRevoked represents a RoleRevoked event raised by the Vault contract.
type VaultRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VaultRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VaultRoleRevokedIterator{contract: _Vault.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *VaultRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultRoleRevoked)
				if err := _Vault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Vault *VaultFilterer) ParseRoleRevoked(log types.Log) (*VaultRoleRevoked, error) {
	event := new(VaultRoleRevoked)
	if err := _Vault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetDelegatorIterator is returned from FilterSetDelegator and is used to iterate over the raw logs and unpacked data for SetDelegator events raised by the Vault contract.
type VaultSetDelegatorIterator struct {
	Event *VaultSetDelegator // Event containing the contract specifics and raw log

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
func (it *VaultSetDelegatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetDelegator)
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
		it.Event = new(VaultSetDelegator)
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
func (it *VaultSetDelegatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetDelegatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetDelegator represents a SetDelegator event raised by the Vault contract.
type VaultSetDelegator struct {
	Delegator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetDelegator is a free log retrieval operation binding the contract event 0xdb2160616f776a37b24808115554e79439bf26cccbbd4438190cc6d28e80ecd1.
//
// Solidity: event SetDelegator(address indexed delegator)
func (_Vault *VaultFilterer) FilterSetDelegator(opts *bind.FilterOpts, delegator []common.Address) (*VaultSetDelegatorIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetDelegator", delegatorRule)
	if err != nil {
		return nil, err
	}
	return &VaultSetDelegatorIterator{contract: _Vault.contract, event: "SetDelegator", logs: logs, sub: sub}, nil
}

// WatchSetDelegator is a free log subscription operation binding the contract event 0xdb2160616f776a37b24808115554e79439bf26cccbbd4438190cc6d28e80ecd1.
//
// Solidity: event SetDelegator(address indexed delegator)
func (_Vault *VaultFilterer) WatchSetDelegator(opts *bind.WatchOpts, sink chan<- *VaultSetDelegator, delegator []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetDelegator", delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetDelegator)
				if err := _Vault.contract.UnpackLog(event, "SetDelegator", log); err != nil {
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

// ParseSetDelegator is a log parse operation binding the contract event 0xdb2160616f776a37b24808115554e79439bf26cccbbd4438190cc6d28e80ecd1.
//
// Solidity: event SetDelegator(address indexed delegator)
func (_Vault *VaultFilterer) ParseSetDelegator(log types.Log) (*VaultSetDelegator, error) {
	event := new(VaultSetDelegator)
	if err := _Vault.contract.UnpackLog(event, "SetDelegator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetDepositLimitIterator is returned from FilterSetDepositLimit and is used to iterate over the raw logs and unpacked data for SetDepositLimit events raised by the Vault contract.
type VaultSetDepositLimitIterator struct {
	Event *VaultSetDepositLimit // Event containing the contract specifics and raw log

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
func (it *VaultSetDepositLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetDepositLimit)
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
		it.Event = new(VaultSetDepositLimit)
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
func (it *VaultSetDepositLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetDepositLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetDepositLimit represents a SetDepositLimit event raised by the Vault contract.
type VaultSetDepositLimit struct {
	Limit *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSetDepositLimit is a free log retrieval operation binding the contract event 0x854df3eb95564502c8bc871ebdd15310ee26270f955f6c6bd8cea68e75045bc0.
//
// Solidity: event SetDepositLimit(uint256 limit)
func (_Vault *VaultFilterer) FilterSetDepositLimit(opts *bind.FilterOpts) (*VaultSetDepositLimitIterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetDepositLimit")
	if err != nil {
		return nil, err
	}
	return &VaultSetDepositLimitIterator{contract: _Vault.contract, event: "SetDepositLimit", logs: logs, sub: sub}, nil
}

// WatchSetDepositLimit is a free log subscription operation binding the contract event 0x854df3eb95564502c8bc871ebdd15310ee26270f955f6c6bd8cea68e75045bc0.
//
// Solidity: event SetDepositLimit(uint256 limit)
func (_Vault *VaultFilterer) WatchSetDepositLimit(opts *bind.WatchOpts, sink chan<- *VaultSetDepositLimit) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetDepositLimit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetDepositLimit)
				if err := _Vault.contract.UnpackLog(event, "SetDepositLimit", log); err != nil {
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

// ParseSetDepositLimit is a log parse operation binding the contract event 0x854df3eb95564502c8bc871ebdd15310ee26270f955f6c6bd8cea68e75045bc0.
//
// Solidity: event SetDepositLimit(uint256 limit)
func (_Vault *VaultFilterer) ParseSetDepositLimit(log types.Log) (*VaultSetDepositLimit, error) {
	event := new(VaultSetDepositLimit)
	if err := _Vault.contract.UnpackLog(event, "SetDepositLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetDepositWhitelistIterator is returned from FilterSetDepositWhitelist and is used to iterate over the raw logs and unpacked data for SetDepositWhitelist events raised by the Vault contract.
type VaultSetDepositWhitelistIterator struct {
	Event *VaultSetDepositWhitelist // Event containing the contract specifics and raw log

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
func (it *VaultSetDepositWhitelistIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetDepositWhitelist)
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
		it.Event = new(VaultSetDepositWhitelist)
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
func (it *VaultSetDepositWhitelistIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetDepositWhitelistIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetDepositWhitelist represents a SetDepositWhitelist event raised by the Vault contract.
type VaultSetDepositWhitelist struct {
	Status bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSetDepositWhitelist is a free log retrieval operation binding the contract event 0x3e12b7b36c75ac9609a3f58609b331210428e1a85909132638955ba0301eec33.
//
// Solidity: event SetDepositWhitelist(bool status)
func (_Vault *VaultFilterer) FilterSetDepositWhitelist(opts *bind.FilterOpts) (*VaultSetDepositWhitelistIterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetDepositWhitelist")
	if err != nil {
		return nil, err
	}
	return &VaultSetDepositWhitelistIterator{contract: _Vault.contract, event: "SetDepositWhitelist", logs: logs, sub: sub}, nil
}

// WatchSetDepositWhitelist is a free log subscription operation binding the contract event 0x3e12b7b36c75ac9609a3f58609b331210428e1a85909132638955ba0301eec33.
//
// Solidity: event SetDepositWhitelist(bool status)
func (_Vault *VaultFilterer) WatchSetDepositWhitelist(opts *bind.WatchOpts, sink chan<- *VaultSetDepositWhitelist) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetDepositWhitelist")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetDepositWhitelist)
				if err := _Vault.contract.UnpackLog(event, "SetDepositWhitelist", log); err != nil {
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

// ParseSetDepositWhitelist is a log parse operation binding the contract event 0x3e12b7b36c75ac9609a3f58609b331210428e1a85909132638955ba0301eec33.
//
// Solidity: event SetDepositWhitelist(bool status)
func (_Vault *VaultFilterer) ParseSetDepositWhitelist(log types.Log) (*VaultSetDepositWhitelist, error) {
	event := new(VaultSetDepositWhitelist)
	if err := _Vault.contract.UnpackLog(event, "SetDepositWhitelist", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetDepositorWhitelistStatusIterator is returned from FilterSetDepositorWhitelistStatus and is used to iterate over the raw logs and unpacked data for SetDepositorWhitelistStatus events raised by the Vault contract.
type VaultSetDepositorWhitelistStatusIterator struct {
	Event *VaultSetDepositorWhitelistStatus // Event containing the contract specifics and raw log

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
func (it *VaultSetDepositorWhitelistStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetDepositorWhitelistStatus)
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
		it.Event = new(VaultSetDepositorWhitelistStatus)
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
func (it *VaultSetDepositorWhitelistStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetDepositorWhitelistStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetDepositorWhitelistStatus represents a SetDepositorWhitelistStatus event raised by the Vault contract.
type VaultSetDepositorWhitelistStatus struct {
	Account common.Address
	Status  bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSetDepositorWhitelistStatus is a free log retrieval operation binding the contract event 0xf991b1ecfb5115cbb36a2b2e2240c058406d2acc2fcc6e9e2dc99d845ff70a62.
//
// Solidity: event SetDepositorWhitelistStatus(address indexed account, bool status)
func (_Vault *VaultFilterer) FilterSetDepositorWhitelistStatus(opts *bind.FilterOpts, account []common.Address) (*VaultSetDepositorWhitelistStatusIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetDepositorWhitelistStatus", accountRule)
	if err != nil {
		return nil, err
	}
	return &VaultSetDepositorWhitelistStatusIterator{contract: _Vault.contract, event: "SetDepositorWhitelistStatus", logs: logs, sub: sub}, nil
}

// WatchSetDepositorWhitelistStatus is a free log subscription operation binding the contract event 0xf991b1ecfb5115cbb36a2b2e2240c058406d2acc2fcc6e9e2dc99d845ff70a62.
//
// Solidity: event SetDepositorWhitelistStatus(address indexed account, bool status)
func (_Vault *VaultFilterer) WatchSetDepositorWhitelistStatus(opts *bind.WatchOpts, sink chan<- *VaultSetDepositorWhitelistStatus, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetDepositorWhitelistStatus", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetDepositorWhitelistStatus)
				if err := _Vault.contract.UnpackLog(event, "SetDepositorWhitelistStatus", log); err != nil {
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

// ParseSetDepositorWhitelistStatus is a log parse operation binding the contract event 0xf991b1ecfb5115cbb36a2b2e2240c058406d2acc2fcc6e9e2dc99d845ff70a62.
//
// Solidity: event SetDepositorWhitelistStatus(address indexed account, bool status)
func (_Vault *VaultFilterer) ParseSetDepositorWhitelistStatus(log types.Log) (*VaultSetDepositorWhitelistStatus, error) {
	event := new(VaultSetDepositorWhitelistStatus)
	if err := _Vault.contract.UnpackLog(event, "SetDepositorWhitelistStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetIsDepositLimitIterator is returned from FilterSetIsDepositLimit and is used to iterate over the raw logs and unpacked data for SetIsDepositLimit events raised by the Vault contract.
type VaultSetIsDepositLimitIterator struct {
	Event *VaultSetIsDepositLimit // Event containing the contract specifics and raw log

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
func (it *VaultSetIsDepositLimitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetIsDepositLimit)
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
		it.Event = new(VaultSetIsDepositLimit)
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
func (it *VaultSetIsDepositLimitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetIsDepositLimitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetIsDepositLimit represents a SetIsDepositLimit event raised by the Vault contract.
type VaultSetIsDepositLimit struct {
	Status bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSetIsDepositLimit is a free log retrieval operation binding the contract event 0xfa7a25a0b611d4ba3c0ea990e90dc23d484a5dd7a1be4733fef2946ba74530c6.
//
// Solidity: event SetIsDepositLimit(bool status)
func (_Vault *VaultFilterer) FilterSetIsDepositLimit(opts *bind.FilterOpts) (*VaultSetIsDepositLimitIterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetIsDepositLimit")
	if err != nil {
		return nil, err
	}
	return &VaultSetIsDepositLimitIterator{contract: _Vault.contract, event: "SetIsDepositLimit", logs: logs, sub: sub}, nil
}

// WatchSetIsDepositLimit is a free log subscription operation binding the contract event 0xfa7a25a0b611d4ba3c0ea990e90dc23d484a5dd7a1be4733fef2946ba74530c6.
//
// Solidity: event SetIsDepositLimit(bool status)
func (_Vault *VaultFilterer) WatchSetIsDepositLimit(opts *bind.WatchOpts, sink chan<- *VaultSetIsDepositLimit) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetIsDepositLimit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetIsDepositLimit)
				if err := _Vault.contract.UnpackLog(event, "SetIsDepositLimit", log); err != nil {
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

// ParseSetIsDepositLimit is a log parse operation binding the contract event 0xfa7a25a0b611d4ba3c0ea990e90dc23d484a5dd7a1be4733fef2946ba74530c6.
//
// Solidity: event SetIsDepositLimit(bool status)
func (_Vault *VaultFilterer) ParseSetIsDepositLimit(log types.Log) (*VaultSetIsDepositLimit, error) {
	event := new(VaultSetIsDepositLimit)
	if err := _Vault.contract.UnpackLog(event, "SetIsDepositLimit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultSetSlasherIterator is returned from FilterSetSlasher and is used to iterate over the raw logs and unpacked data for SetSlasher events raised by the Vault contract.
type VaultSetSlasherIterator struct {
	Event *VaultSetSlasher // Event containing the contract specifics and raw log

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
func (it *VaultSetSlasherIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultSetSlasher)
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
		it.Event = new(VaultSetSlasher)
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
func (it *VaultSetSlasherIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultSetSlasherIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultSetSlasher represents a SetSlasher event raised by the Vault contract.
type VaultSetSlasher struct {
	Slasher common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSetSlasher is a free log retrieval operation binding the contract event 0xe7e4c932e03abddfe20f83af42c33627e816115c7ec2b168441f65dc14bfc3ba.
//
// Solidity: event SetSlasher(address indexed slasher)
func (_Vault *VaultFilterer) FilterSetSlasher(opts *bind.FilterOpts, slasher []common.Address) (*VaultSetSlasherIterator, error) {

	var slasherRule []interface{}
	for _, slasherItem := range slasher {
		slasherRule = append(slasherRule, slasherItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "SetSlasher", slasherRule)
	if err != nil {
		return nil, err
	}
	return &VaultSetSlasherIterator{contract: _Vault.contract, event: "SetSlasher", logs: logs, sub: sub}, nil
}

// WatchSetSlasher is a free log subscription operation binding the contract event 0xe7e4c932e03abddfe20f83af42c33627e816115c7ec2b168441f65dc14bfc3ba.
//
// Solidity: event SetSlasher(address indexed slasher)
func (_Vault *VaultFilterer) WatchSetSlasher(opts *bind.WatchOpts, sink chan<- *VaultSetSlasher, slasher []common.Address) (event.Subscription, error) {

	var slasherRule []interface{}
	for _, slasherItem := range slasher {
		slasherRule = append(slasherRule, slasherItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "SetSlasher", slasherRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultSetSlasher)
				if err := _Vault.contract.UnpackLog(event, "SetSlasher", log); err != nil {
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

// ParseSetSlasher is a log parse operation binding the contract event 0xe7e4c932e03abddfe20f83af42c33627e816115c7ec2b168441f65dc14bfc3ba.
//
// Solidity: event SetSlasher(address indexed slasher)
func (_Vault *VaultFilterer) ParseSetSlasher(log types.Log) (*VaultSetSlasher, error) {
	event := new(VaultSetSlasher)
	if err := _Vault.contract.UnpackLog(event, "SetSlasher", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Vault contract.
type VaultWithdrawIterator struct {
	Event *VaultWithdraw // Event containing the contract specifics and raw log

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
func (it *VaultWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultWithdraw)
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
		it.Event = new(VaultWithdraw)
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
func (it *VaultWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultWithdraw represents a Withdraw event raised by the Vault contract.
type VaultWithdraw struct {
	Withdrawer   common.Address
	Claimer      common.Address
	Amount       *big.Int
	BurnedShares *big.Int
	MintedShares *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xebff2602b3f468259e1e99f613fed6691f3a6526effe6ef3e768ba7ae7a36c4f.
//
// Solidity: event Withdraw(address indexed withdrawer, address indexed claimer, uint256 amount, uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultFilterer) FilterWithdraw(opts *bind.FilterOpts, withdrawer []common.Address, claimer []common.Address) (*VaultWithdrawIterator, error) {

	var withdrawerRule []interface{}
	for _, withdrawerItem := range withdrawer {
		withdrawerRule = append(withdrawerRule, withdrawerItem)
	}
	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Withdraw", withdrawerRule, claimerRule)
	if err != nil {
		return nil, err
	}
	return &VaultWithdrawIterator{contract: _Vault.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xebff2602b3f468259e1e99f613fed6691f3a6526effe6ef3e768ba7ae7a36c4f.
//
// Solidity: event Withdraw(address indexed withdrawer, address indexed claimer, uint256 amount, uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *VaultWithdraw, withdrawer []common.Address, claimer []common.Address) (event.Subscription, error) {

	var withdrawerRule []interface{}
	for _, withdrawerItem := range withdrawer {
		withdrawerRule = append(withdrawerRule, withdrawerItem)
	}
	var claimerRule []interface{}
	for _, claimerItem := range claimer {
		claimerRule = append(claimerRule, claimerItem)
	}

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Withdraw", withdrawerRule, claimerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultWithdraw)
				if err := _Vault.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xebff2602b3f468259e1e99f613fed6691f3a6526effe6ef3e768ba7ae7a36c4f.
//
// Solidity: event Withdraw(address indexed withdrawer, address indexed claimer, uint256 amount, uint256 burnedShares, uint256 mintedShares)
func (_Vault *VaultFilterer) ParseWithdraw(log types.Log) (*VaultWithdraw, error) {
	event := new(VaultWithdraw)
	if err := _Vault.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
