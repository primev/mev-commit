// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mevcommitmiddleware

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

// CheckpointsCheckpoint160 is an auto generated low-level Go binding around an user-defined struct.
type CheckpointsCheckpoint160 struct {
	Key   *big.Int
	Value *big.Int
}

// CheckpointsTrace160 is an auto generated low-level Go binding around an user-defined struct.
type CheckpointsTrace160 struct {
	Checkpoints []CheckpointsCheckpoint160
}

// TimestampOccurrenceOccurrence is an auto generated low-level Go binding around an user-defined struct.
type TimestampOccurrenceOccurrence struct {
	Exists    bool
	Timestamp *big.Int
}

// MevcommitmiddlewareMetaData contains all meta data concerning the Mevcommitmiddleware contract.
var MevcommitmiddlewareMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"blacklistOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnerRouterFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegatorFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterValidators\",\"inputs\":[{\"name\":\"blsPubkeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterVaults\",\"inputs\":[{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getLatestSlashAmount\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint160\",\"internalType\":\"uint160\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNumSlashableVals\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPositionInValset\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlashAmountAt\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blockTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint160\",\"internalType\":\"uint160\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_networkRegistry\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_operatorRegistry\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_vaultFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_delegatorFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_slasherFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_burnerRouterFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"},{\"name\":\"_network\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_slashPeriodSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_slashOracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_slashReceiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_minBurnerRouterDelay\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidatorSlashable\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isVaultBurnerValid\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isVaultBurnerValidAgainstOperator\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minBurnerRouterDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"network\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"networkRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorRecords\",\"inputs\":[{\"name\":\"operatorAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structTimestampOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isBlacklisted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"potentialSlashableValidators\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pubkeyAtPositionInValset\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidators\",\"inputs\":[{\"name\":\"blsPubkeys\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerVaults\",\"inputs\":[{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"slashAmounts\",\"type\":\"uint160[]\",\"internalType\":\"uint160[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestOperatorDeregistrations\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestValDeregistrations\",\"inputs\":[{\"name\":\"blsPubkeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestVaultDeregistrations\",\"inputs\":[{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBurnerRouterFactory\",\"inputs\":[{\"name\":\"_burnerRouterFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegatorFactory\",\"inputs\":[{\"name\":\"_delegatorFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinBurnerRouterDelay\",\"inputs\":[{\"name\":\"minBurnerRouterDelay_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNetwork\",\"inputs\":[{\"name\":\"_network\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setNetworkRegistry\",\"inputs\":[{\"name\":\"_networkRegistry\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorRegistry\",\"inputs\":[{\"name\":\"_operatorRegistry\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashOracle\",\"inputs\":[{\"name\":\"slashOracle_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashPeriodSeconds\",\"inputs\":[{\"name\":\"slashPeriodSeconds_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlashReceiver\",\"inputs\":[{\"name\":\"slashReceiver_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSlasherFactory\",\"inputs\":[{\"name\":\"_slasherFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVaultFactory\",\"inputs\":[{\"name\":\"_vaultFactory\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashPeriodSeconds\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashReceiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashRecords\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"numSlashed\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"numRegistered\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slashValidators\",\"inputs\":[{\"name\":\"blsPubkeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"captureTimestamps\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slasherFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unblacklistOperators\",\"inputs\":[{\"name\":\"operators\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSlashAmounts\",\"inputs\":[{\"name\":\"vaults\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"slashAmounts\",\"type\":\"uint160[]\",\"internalType\":\"uint160[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validatorRecords\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structTimestampOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"valsetLength\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaultFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vaultRecords\",\"inputs\":[{\"name\":\"vaultAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structTimestampOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"slashAmountHistory\",\"type\":\"tuple\",\"internalType\":\"structCheckpoints.Trace160\",\"components\":[{\"name\":\"_checkpoints\",\"type\":\"tuple[]\",\"internalType\":\"structCheckpoints.Checkpoint160[]\",\"components\":[{\"name\":\"_key\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"_value\",\"type\":\"uint160\",\"internalType\":\"uint160\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"wouldVaultBeValidWith\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"potentialSLashPeriodSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BurnerRouterFactorySet\",\"inputs\":[{\"name\":\"burnerRouterFactory\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegatorFactorySet\",\"inputs\":[{\"name\":\"delegatorFactory\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinBurnerRouterDelaySet\",\"inputs\":[{\"name\":\"minBurnerRouterDelay\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NetworkRegistrySet\",\"inputs\":[{\"name\":\"networkRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NetworkSet\",\"inputs\":[{\"name\":\"network\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorBlacklisted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistrationRequested\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistrySet\",\"inputs\":[{\"name\":\"operatorRegistry\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUnblacklisted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashOracleSet\",\"inputs\":[{\"name\":\"slashOracle\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashPeriodBlocksSet\",\"inputs\":[{\"name\":\"slashPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashPeriodSecondsSet\",\"inputs\":[{\"name\":\"slashPeriodSeconds\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlashReceiverSet\",\"inputs\":[{\"name\":\"slashReceiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SlasherFactorySet\",\"inputs\":[{\"name\":\"slasherFactory\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValRecordAdded\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValRecordDeleted\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistrationRequested\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"msgSender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"position\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorPositionsSwapped\",\"inputs\":[{\"name\":\"blsPubkeys\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"vaults\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"operators\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"newPositions\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSlashed\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultDeregistered\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultDeregistrationRequested\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultFactorySet\",\"inputs\":[{\"name\":\"vaultFactory\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultRegistered\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashAmount\",\"type\":\"uint160\",\"indexed\":false,\"internalType\":\"uint160\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VaultSlashAmountUpdated\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"slashAmount\",\"type\":\"uint160\",\"indexed\":false,\"internalType\":\"uint160\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BurnerHookNotSetForVault\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"CaptureTimestampMustBeNonZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CheckpointUnorderedInsertion\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DelegatorNotEntity\",\"inputs\":[{\"name\":\"delegator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatorFactory\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedToAddValidatorToValset\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FullRestakeDelegatorNotSupported\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"FutureTimestampDisallowed\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidArrayLengths\",\"inputs\":[{\"name\":\"vaultLen\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pubkeyLen\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidBLSPubKeyLength\",\"inputs\":[{\"name\":\"expectedLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actualLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidVaultBurner\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidVaultBurnerConsideringOperator\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidVaultEpochDuration\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vaultEpochDurationSec\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashPeriodSec\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"MissingOperatorRecord\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"MissingValRecord\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MissingValidatorRecord\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"MissingVaultRecord\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NetworkNotEntity\",\"inputs\":[{\"name\":\"network\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoRegisteredValidators\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NoSlashAmountAtTimestamp\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlySlashOracle\",\"inputs\":[{\"name\":\"slashOracle\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OnlyVetoSlashersRequireExecution\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slasherType\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OperatorAlreadyBlacklisted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorAlreadyRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorDeregRequestExists\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorIsBlacklisted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorNotBlacklisted\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorNotEntity\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OperatorNotReadyToDeregister\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"OperatorNotRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SlashAmountMustBeNonZero\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlasherNotEntity\",\"inputs\":[{\"name\":\"slasher\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slasherFactory\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SlasherNotSetForVault\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnknownDelegatorType\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegatorType\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownSlasherType\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"slasherType\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorDeregRequestExists\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotInValset\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotReadyToDeregister\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"currentTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotRemovedFromValset\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotSlashable\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ValidatorRecordAlreadyExists\",\"inputs\":[{\"name\":\"blsPubkey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorsNotSlashable\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"numRequested\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"potentialSlashableVals\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"VaultAlreadyRegistered\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultDeregNotRequested\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultDeregRequestExists\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultNotEntity\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VaultNotReadyToDeregister\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestTimestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"VaultNotRegistered\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"VetoDurationTooShort\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vetoDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"VetoSlasherMustHaveZeroResolver\",\"inputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddressNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroUintNotAllowed\",\"inputs\":[]}]",
}

// MevcommitmiddlewareABI is the input ABI used to generate the binding from.
// Deprecated: Use MevcommitmiddlewareMetaData.ABI instead.
var MevcommitmiddlewareABI = MevcommitmiddlewareMetaData.ABI

// Mevcommitmiddleware is an auto generated Go binding around an Ethereum contract.
type Mevcommitmiddleware struct {
	MevcommitmiddlewareCaller     // Read-only binding to the contract
	MevcommitmiddlewareTransactor // Write-only binding to the contract
	MevcommitmiddlewareFilterer   // Log filterer for contract events
}

// MevcommitmiddlewareCaller is an auto generated read-only Go binding around an Ethereum contract.
type MevcommitmiddlewareCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitmiddlewareTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MevcommitmiddlewareTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitmiddlewareFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MevcommitmiddlewareFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitmiddlewareSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MevcommitmiddlewareSession struct {
	Contract     *Mevcommitmiddleware // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MevcommitmiddlewareCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MevcommitmiddlewareCallerSession struct {
	Contract *MevcommitmiddlewareCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// MevcommitmiddlewareTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MevcommitmiddlewareTransactorSession struct {
	Contract     *MevcommitmiddlewareTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// MevcommitmiddlewareRaw is an auto generated low-level Go binding around an Ethereum contract.
type MevcommitmiddlewareRaw struct {
	Contract *Mevcommitmiddleware // Generic contract binding to access the raw methods on
}

// MevcommitmiddlewareCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MevcommitmiddlewareCallerRaw struct {
	Contract *MevcommitmiddlewareCaller // Generic read-only contract binding to access the raw methods on
}

// MevcommitmiddlewareTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MevcommitmiddlewareTransactorRaw struct {
	Contract *MevcommitmiddlewareTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMevcommitmiddleware creates a new instance of Mevcommitmiddleware, bound to a specific deployed contract.
func NewMevcommitmiddleware(address common.Address, backend bind.ContractBackend) (*Mevcommitmiddleware, error) {
	contract, err := bindMevcommitmiddleware(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mevcommitmiddleware{MevcommitmiddlewareCaller: MevcommitmiddlewareCaller{contract: contract}, MevcommitmiddlewareTransactor: MevcommitmiddlewareTransactor{contract: contract}, MevcommitmiddlewareFilterer: MevcommitmiddlewareFilterer{contract: contract}}, nil
}

// NewMevcommitmiddlewareCaller creates a new read-only instance of Mevcommitmiddleware, bound to a specific deployed contract.
func NewMevcommitmiddlewareCaller(address common.Address, caller bind.ContractCaller) (*MevcommitmiddlewareCaller, error) {
	contract, err := bindMevcommitmiddleware(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareCaller{contract: contract}, nil
}

// NewMevcommitmiddlewareTransactor creates a new write-only instance of Mevcommitmiddleware, bound to a specific deployed contract.
func NewMevcommitmiddlewareTransactor(address common.Address, transactor bind.ContractTransactor) (*MevcommitmiddlewareTransactor, error) {
	contract, err := bindMevcommitmiddleware(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareTransactor{contract: contract}, nil
}

// NewMevcommitmiddlewareFilterer creates a new log filterer instance of Mevcommitmiddleware, bound to a specific deployed contract.
func NewMevcommitmiddlewareFilterer(address common.Address, filterer bind.ContractFilterer) (*MevcommitmiddlewareFilterer, error) {
	contract, err := bindMevcommitmiddleware(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareFilterer{contract: contract}, nil
}

// bindMevcommitmiddleware binds a generic wrapper to an already deployed contract.
func bindMevcommitmiddleware(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MevcommitmiddlewareMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitmiddleware *MevcommitmiddlewareRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitmiddleware.Contract.MevcommitmiddlewareCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitmiddleware *MevcommitmiddlewareRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.MevcommitmiddlewareTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitmiddleware *MevcommitmiddlewareRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.MevcommitmiddlewareTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitmiddleware.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitmiddleware.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitmiddleware.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitmiddleware.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitmiddleware.CallOpts)
}

// BurnerRouterFactory is a free data retrieval call binding the contract method 0xc70a3c5f.
//
// Solidity: function burnerRouterFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) BurnerRouterFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "burnerRouterFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BurnerRouterFactory is a free data retrieval call binding the contract method 0xc70a3c5f.
//
// Solidity: function burnerRouterFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) BurnerRouterFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.BurnerRouterFactory(&_Mevcommitmiddleware.CallOpts)
}

// BurnerRouterFactory is a free data retrieval call binding the contract method 0xc70a3c5f.
//
// Solidity: function burnerRouterFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) BurnerRouterFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.BurnerRouterFactory(&_Mevcommitmiddleware.CallOpts)
}

// DelegatorFactory is a free data retrieval call binding the contract method 0x079faad4.
//
// Solidity: function delegatorFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) DelegatorFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "delegatorFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DelegatorFactory is a free data retrieval call binding the contract method 0x079faad4.
//
// Solidity: function delegatorFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) DelegatorFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.DelegatorFactory(&_Mevcommitmiddleware.CallOpts)
}

// DelegatorFactory is a free data retrieval call binding the contract method 0x079faad4.
//
// Solidity: function delegatorFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) DelegatorFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.DelegatorFactory(&_Mevcommitmiddleware.CallOpts)
}

// GetLatestSlashAmount is a free data retrieval call binding the contract method 0xb39edc0f.
//
// Solidity: function getLatestSlashAmount(address vault) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) GetLatestSlashAmount(opts *bind.CallOpts, vault common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "getLatestSlashAmount", vault)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLatestSlashAmount is a free data retrieval call binding the contract method 0xb39edc0f.
//
// Solidity: function getLatestSlashAmount(address vault) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) GetLatestSlashAmount(vault common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetLatestSlashAmount(&_Mevcommitmiddleware.CallOpts, vault)
}

// GetLatestSlashAmount is a free data retrieval call binding the contract method 0xb39edc0f.
//
// Solidity: function getLatestSlashAmount(address vault) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) GetLatestSlashAmount(vault common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetLatestSlashAmount(&_Mevcommitmiddleware.CallOpts, vault)
}

// GetNumSlashableVals is a free data retrieval call binding the contract method 0xf9a9184d.
//
// Solidity: function getNumSlashableVals(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) GetNumSlashableVals(opts *bind.CallOpts, vault common.Address, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "getNumSlashableVals", vault, operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNumSlashableVals is a free data retrieval call binding the contract method 0xf9a9184d.
//
// Solidity: function getNumSlashableVals(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) GetNumSlashableVals(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetNumSlashableVals(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// GetNumSlashableVals is a free data retrieval call binding the contract method 0xf9a9184d.
//
// Solidity: function getNumSlashableVals(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) GetNumSlashableVals(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetNumSlashableVals(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// GetPositionInValset is a free data retrieval call binding the contract method 0x2f205b6c.
//
// Solidity: function getPositionInValset(bytes blsPubkey, address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) GetPositionInValset(opts *bind.CallOpts, blsPubkey []byte, vault common.Address, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "getPositionInValset", blsPubkey, vault, operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPositionInValset is a free data retrieval call binding the contract method 0x2f205b6c.
//
// Solidity: function getPositionInValset(bytes blsPubkey, address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) GetPositionInValset(blsPubkey []byte, vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetPositionInValset(&_Mevcommitmiddleware.CallOpts, blsPubkey, vault, operator)
}

// GetPositionInValset is a free data retrieval call binding the contract method 0x2f205b6c.
//
// Solidity: function getPositionInValset(bytes blsPubkey, address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) GetPositionInValset(blsPubkey []byte, vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetPositionInValset(&_Mevcommitmiddleware.CallOpts, blsPubkey, vault, operator)
}

// GetSlashAmountAt is a free data retrieval call binding the contract method 0xec8ece5d.
//
// Solidity: function getSlashAmountAt(address vault, uint256 blockTimestamp) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) GetSlashAmountAt(opts *bind.CallOpts, vault common.Address, blockTimestamp *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "getSlashAmountAt", vault, blockTimestamp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlashAmountAt is a free data retrieval call binding the contract method 0xec8ece5d.
//
// Solidity: function getSlashAmountAt(address vault, uint256 blockTimestamp) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) GetSlashAmountAt(vault common.Address, blockTimestamp *big.Int) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetSlashAmountAt(&_Mevcommitmiddleware.CallOpts, vault, blockTimestamp)
}

// GetSlashAmountAt is a free data retrieval call binding the contract method 0xec8ece5d.
//
// Solidity: function getSlashAmountAt(address vault, uint256 blockTimestamp) view returns(uint160)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) GetSlashAmountAt(vault common.Address, blockTimestamp *big.Int) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.GetSlashAmountAt(&_Mevcommitmiddleware.CallOpts, vault, blockTimestamp)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) IsValidatorOptedIn(opts *bind.CallOpts, blsPubkey []byte) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "isValidatorOptedIn", blsPubkey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) IsValidatorOptedIn(blsPubkey []byte) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsValidatorOptedIn(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) IsValidatorOptedIn(blsPubkey []byte) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsValidatorOptedIn(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// IsValidatorSlashable is a free data retrieval call binding the contract method 0x6fe0b852.
//
// Solidity: function isValidatorSlashable(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) IsValidatorSlashable(opts *bind.CallOpts, blsPubkey []byte) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "isValidatorSlashable", blsPubkey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorSlashable is a free data retrieval call binding the contract method 0x6fe0b852.
//
// Solidity: function isValidatorSlashable(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) IsValidatorSlashable(blsPubkey []byte) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsValidatorSlashable(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// IsValidatorSlashable is a free data retrieval call binding the contract method 0x6fe0b852.
//
// Solidity: function isValidatorSlashable(bytes blsPubkey) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) IsValidatorSlashable(blsPubkey []byte) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsValidatorSlashable(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// IsVaultBurnerValid is a free data retrieval call binding the contract method 0x80bd4e93.
//
// Solidity: function isVaultBurnerValid(address vault) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) IsVaultBurnerValid(opts *bind.CallOpts, vault common.Address) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "isVaultBurnerValid", vault)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsVaultBurnerValid is a free data retrieval call binding the contract method 0x80bd4e93.
//
// Solidity: function isVaultBurnerValid(address vault) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) IsVaultBurnerValid(vault common.Address) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsVaultBurnerValid(&_Mevcommitmiddleware.CallOpts, vault)
}

// IsVaultBurnerValid is a free data retrieval call binding the contract method 0x80bd4e93.
//
// Solidity: function isVaultBurnerValid(address vault) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) IsVaultBurnerValid(vault common.Address) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsVaultBurnerValid(&_Mevcommitmiddleware.CallOpts, vault)
}

// IsVaultBurnerValidAgainstOperator is a free data retrieval call binding the contract method 0x7d194fdf.
//
// Solidity: function isVaultBurnerValidAgainstOperator(address vault, address operator) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) IsVaultBurnerValidAgainstOperator(opts *bind.CallOpts, vault common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "isVaultBurnerValidAgainstOperator", vault, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsVaultBurnerValidAgainstOperator is a free data retrieval call binding the contract method 0x7d194fdf.
//
// Solidity: function isVaultBurnerValidAgainstOperator(address vault, address operator) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) IsVaultBurnerValidAgainstOperator(vault common.Address, operator common.Address) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsVaultBurnerValidAgainstOperator(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// IsVaultBurnerValidAgainstOperator is a free data retrieval call binding the contract method 0x7d194fdf.
//
// Solidity: function isVaultBurnerValidAgainstOperator(address vault, address operator) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) IsVaultBurnerValidAgainstOperator(vault common.Address, operator common.Address) (bool, error) {
	return _Mevcommitmiddleware.Contract.IsVaultBurnerValidAgainstOperator(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// MinBurnerRouterDelay is a free data retrieval call binding the contract method 0x8f55b4f0.
//
// Solidity: function minBurnerRouterDelay() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) MinBurnerRouterDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "minBurnerRouterDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinBurnerRouterDelay is a free data retrieval call binding the contract method 0x8f55b4f0.
//
// Solidity: function minBurnerRouterDelay() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) MinBurnerRouterDelay() (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.MinBurnerRouterDelay(&_Mevcommitmiddleware.CallOpts)
}

// MinBurnerRouterDelay is a free data retrieval call binding the contract method 0x8f55b4f0.
//
// Solidity: function minBurnerRouterDelay() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) MinBurnerRouterDelay() (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.MinBurnerRouterDelay(&_Mevcommitmiddleware.CallOpts)
}

// Network is a free data retrieval call binding the contract method 0x6739afca.
//
// Solidity: function network() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) Network(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "network")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Network is a free data retrieval call binding the contract method 0x6739afca.
//
// Solidity: function network() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Network() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.Network(&_Mevcommitmiddleware.CallOpts)
}

// Network is a free data retrieval call binding the contract method 0x6739afca.
//
// Solidity: function network() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) Network() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.Network(&_Mevcommitmiddleware.CallOpts)
}

// NetworkRegistry is a free data retrieval call binding the contract method 0xe45f40be.
//
// Solidity: function networkRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) NetworkRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "networkRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NetworkRegistry is a free data retrieval call binding the contract method 0xe45f40be.
//
// Solidity: function networkRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) NetworkRegistry() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.NetworkRegistry(&_Mevcommitmiddleware.CallOpts)
}

// NetworkRegistry is a free data retrieval call binding the contract method 0xe45f40be.
//
// Solidity: function networkRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) NetworkRegistry() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.NetworkRegistry(&_Mevcommitmiddleware.CallOpts)
}

// OperatorRecords is a free data retrieval call binding the contract method 0xd0b8643f.
//
// Solidity: function operatorRecords(address operatorAddress) view returns((bool,uint256) deregRequestOccurrence, bool exists, bool isBlacklisted)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) OperatorRecords(opts *bind.CallOpts, operatorAddress common.Address) (struct {
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	Exists                 bool
	IsBlacklisted          bool
}, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "operatorRecords", operatorAddress)

	outstruct := new(struct {
		DeregRequestOccurrence TimestampOccurrenceOccurrence
		Exists                 bool
		IsBlacklisted          bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[0], new(TimestampOccurrenceOccurrence)).(*TimestampOccurrenceOccurrence)
	outstruct.Exists = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.IsBlacklisted = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// OperatorRecords is a free data retrieval call binding the contract method 0xd0b8643f.
//
// Solidity: function operatorRecords(address operatorAddress) view returns((bool,uint256) deregRequestOccurrence, bool exists, bool isBlacklisted)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) OperatorRecords(operatorAddress common.Address) (struct {
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	Exists                 bool
	IsBlacklisted          bool
}, error) {
	return _Mevcommitmiddleware.Contract.OperatorRecords(&_Mevcommitmiddleware.CallOpts, operatorAddress)
}

// OperatorRecords is a free data retrieval call binding the contract method 0xd0b8643f.
//
// Solidity: function operatorRecords(address operatorAddress) view returns((bool,uint256) deregRequestOccurrence, bool exists, bool isBlacklisted)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) OperatorRecords(operatorAddress common.Address) (struct {
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	Exists                 bool
	IsBlacklisted          bool
}, error) {
	return _Mevcommitmiddleware.Contract.OperatorRecords(&_Mevcommitmiddleware.CallOpts, operatorAddress)
}

// OperatorRegistry is a free data retrieval call binding the contract method 0x58c2225b.
//
// Solidity: function operatorRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) OperatorRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "operatorRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OperatorRegistry is a free data retrieval call binding the contract method 0x58c2225b.
//
// Solidity: function operatorRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) OperatorRegistry() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.OperatorRegistry(&_Mevcommitmiddleware.CallOpts)
}

// OperatorRegistry is a free data retrieval call binding the contract method 0x58c2225b.
//
// Solidity: function operatorRegistry() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) OperatorRegistry() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.OperatorRegistry(&_Mevcommitmiddleware.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Owner() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.Owner(&_Mevcommitmiddleware.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) Owner() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.Owner(&_Mevcommitmiddleware.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Paused() (bool, error) {
	return _Mevcommitmiddleware.Contract.Paused(&_Mevcommitmiddleware.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) Paused() (bool, error) {
	return _Mevcommitmiddleware.Contract.Paused(&_Mevcommitmiddleware.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) PendingOwner() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.PendingOwner(&_Mevcommitmiddleware.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) PendingOwner() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.PendingOwner(&_Mevcommitmiddleware.CallOpts)
}

// PotentialSlashableValidators is a free data retrieval call binding the contract method 0x608bbd64.
//
// Solidity: function potentialSlashableValidators(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) PotentialSlashableValidators(opts *bind.CallOpts, vault common.Address, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "potentialSlashableValidators", vault, operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PotentialSlashableValidators is a free data retrieval call binding the contract method 0x608bbd64.
//
// Solidity: function potentialSlashableValidators(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) PotentialSlashableValidators(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.PotentialSlashableValidators(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// PotentialSlashableValidators is a free data retrieval call binding the contract method 0x608bbd64.
//
// Solidity: function potentialSlashableValidators(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) PotentialSlashableValidators(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.PotentialSlashableValidators(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitmiddleware.Contract.ProxiableUUID(&_Mevcommitmiddleware.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitmiddleware.Contract.ProxiableUUID(&_Mevcommitmiddleware.CallOpts)
}

// PubkeyAtPositionInValset is a free data retrieval call binding the contract method 0x09aa2431.
//
// Solidity: function pubkeyAtPositionInValset(uint256 index, address vault, address operator) view returns(bytes)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) PubkeyAtPositionInValset(opts *bind.CallOpts, index *big.Int, vault common.Address, operator common.Address) ([]byte, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "pubkeyAtPositionInValset", index, vault, operator)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// PubkeyAtPositionInValset is a free data retrieval call binding the contract method 0x09aa2431.
//
// Solidity: function pubkeyAtPositionInValset(uint256 index, address vault, address operator) view returns(bytes)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) PubkeyAtPositionInValset(index *big.Int, vault common.Address, operator common.Address) ([]byte, error) {
	return _Mevcommitmiddleware.Contract.PubkeyAtPositionInValset(&_Mevcommitmiddleware.CallOpts, index, vault, operator)
}

// PubkeyAtPositionInValset is a free data retrieval call binding the contract method 0x09aa2431.
//
// Solidity: function pubkeyAtPositionInValset(uint256 index, address vault, address operator) view returns(bytes)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) PubkeyAtPositionInValset(index *big.Int, vault common.Address, operator common.Address) ([]byte, error) {
	return _Mevcommitmiddleware.Contract.PubkeyAtPositionInValset(&_Mevcommitmiddleware.CallOpts, index, vault, operator)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) SlashOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "slashOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlashOracle() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlashOracle(&_Mevcommitmiddleware.CallOpts)
}

// SlashOracle is a free data retrieval call binding the contract method 0x38063b54.
//
// Solidity: function slashOracle() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) SlashOracle() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlashOracle(&_Mevcommitmiddleware.CallOpts)
}

// SlashPeriodSeconds is a free data retrieval call binding the contract method 0x61793ef9.
//
// Solidity: function slashPeriodSeconds() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) SlashPeriodSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "slashPeriodSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashPeriodSeconds is a free data retrieval call binding the contract method 0x61793ef9.
//
// Solidity: function slashPeriodSeconds() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlashPeriodSeconds() (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.SlashPeriodSeconds(&_Mevcommitmiddleware.CallOpts)
}

// SlashPeriodSeconds is a free data retrieval call binding the contract method 0x61793ef9.
//
// Solidity: function slashPeriodSeconds() view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) SlashPeriodSeconds() (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.SlashPeriodSeconds(&_Mevcommitmiddleware.CallOpts)
}

// SlashReceiver is a free data retrieval call binding the contract method 0x1bc4e5fb.
//
// Solidity: function slashReceiver() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) SlashReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "slashReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlashReceiver is a free data retrieval call binding the contract method 0x1bc4e5fb.
//
// Solidity: function slashReceiver() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlashReceiver() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlashReceiver(&_Mevcommitmiddleware.CallOpts)
}

// SlashReceiver is a free data retrieval call binding the contract method 0x1bc4e5fb.
//
// Solidity: function slashReceiver() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) SlashReceiver() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlashReceiver(&_Mevcommitmiddleware.CallOpts)
}

// SlashRecords is a free data retrieval call binding the contract method 0x1da9f192.
//
// Solidity: function slashRecords(address vault, address operator, uint256 blockNumber) view returns(bool exists, uint256 numSlashed, uint256 numRegistered)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) SlashRecords(opts *bind.CallOpts, vault common.Address, operator common.Address, blockNumber *big.Int) (struct {
	Exists        bool
	NumSlashed    *big.Int
	NumRegistered *big.Int
}, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "slashRecords", vault, operator, blockNumber)

	outstruct := new(struct {
		Exists        bool
		NumSlashed    *big.Int
		NumRegistered *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.NumSlashed = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NumRegistered = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SlashRecords is a free data retrieval call binding the contract method 0x1da9f192.
//
// Solidity: function slashRecords(address vault, address operator, uint256 blockNumber) view returns(bool exists, uint256 numSlashed, uint256 numRegistered)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlashRecords(vault common.Address, operator common.Address, blockNumber *big.Int) (struct {
	Exists        bool
	NumSlashed    *big.Int
	NumRegistered *big.Int
}, error) {
	return _Mevcommitmiddleware.Contract.SlashRecords(&_Mevcommitmiddleware.CallOpts, vault, operator, blockNumber)
}

// SlashRecords is a free data retrieval call binding the contract method 0x1da9f192.
//
// Solidity: function slashRecords(address vault, address operator, uint256 blockNumber) view returns(bool exists, uint256 numSlashed, uint256 numRegistered)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) SlashRecords(vault common.Address, operator common.Address, blockNumber *big.Int) (struct {
	Exists        bool
	NumSlashed    *big.Int
	NumRegistered *big.Int
}, error) {
	return _Mevcommitmiddleware.Contract.SlashRecords(&_Mevcommitmiddleware.CallOpts, vault, operator, blockNumber)
}

// SlasherFactory is a free data retrieval call binding the contract method 0x6a3f8b5f.
//
// Solidity: function slasherFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) SlasherFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "slasherFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SlasherFactory is a free data retrieval call binding the contract method 0x6a3f8b5f.
//
// Solidity: function slasherFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlasherFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlasherFactory(&_Mevcommitmiddleware.CallOpts)
}

// SlasherFactory is a free data retrieval call binding the contract method 0x6a3f8b5f.
//
// Solidity: function slasherFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) SlasherFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.SlasherFactory(&_Mevcommitmiddleware.CallOpts)
}

// ValidatorRecords is a free data retrieval call binding the contract method 0x052bbca0.
//
// Solidity: function validatorRecords(bytes blsPubkey) view returns(address vault, address operator, bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) ValidatorRecords(opts *bind.CallOpts, blsPubkey []byte) (struct {
	Vault                  common.Address
	Operator               common.Address
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "validatorRecords", blsPubkey)

	outstruct := new(struct {
		Vault                  common.Address
		Operator               common.Address
		Exists                 bool
		DeregRequestOccurrence TimestampOccurrenceOccurrence
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Vault = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Operator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Exists = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[3], new(TimestampOccurrenceOccurrence)).(*TimestampOccurrenceOccurrence)

	return *outstruct, err

}

// ValidatorRecords is a free data retrieval call binding the contract method 0x052bbca0.
//
// Solidity: function validatorRecords(bytes blsPubkey) view returns(address vault, address operator, bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) ValidatorRecords(blsPubkey []byte) (struct {
	Vault                  common.Address
	Operator               common.Address
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
}, error) {
	return _Mevcommitmiddleware.Contract.ValidatorRecords(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// ValidatorRecords is a free data retrieval call binding the contract method 0x052bbca0.
//
// Solidity: function validatorRecords(bytes blsPubkey) view returns(address vault, address operator, bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) ValidatorRecords(blsPubkey []byte) (struct {
	Vault                  common.Address
	Operator               common.Address
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
}, error) {
	return _Mevcommitmiddleware.Contract.ValidatorRecords(&_Mevcommitmiddleware.CallOpts, blsPubkey)
}

// ValsetLength is a free data retrieval call binding the contract method 0x284759ba.
//
// Solidity: function valsetLength(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) ValsetLength(opts *bind.CallOpts, vault common.Address, operator common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "valsetLength", vault, operator)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValsetLength is a free data retrieval call binding the contract method 0x284759ba.
//
// Solidity: function valsetLength(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) ValsetLength(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.ValsetLength(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// ValsetLength is a free data retrieval call binding the contract method 0x284759ba.
//
// Solidity: function valsetLength(address vault, address operator) view returns(uint256)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) ValsetLength(vault common.Address, operator common.Address) (*big.Int, error) {
	return _Mevcommitmiddleware.Contract.ValsetLength(&_Mevcommitmiddleware.CallOpts, vault, operator)
}

// VaultFactory is a free data retrieval call binding the contract method 0xd8a06f73.
//
// Solidity: function vaultFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) VaultFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "vaultFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VaultFactory is a free data retrieval call binding the contract method 0xd8a06f73.
//
// Solidity: function vaultFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) VaultFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.VaultFactory(&_Mevcommitmiddleware.CallOpts)
}

// VaultFactory is a free data retrieval call binding the contract method 0xd8a06f73.
//
// Solidity: function vaultFactory() view returns(address)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) VaultFactory() (common.Address, error) {
	return _Mevcommitmiddleware.Contract.VaultFactory(&_Mevcommitmiddleware.CallOpts)
}

// VaultRecords is a free data retrieval call binding the contract method 0xdef52ca5.
//
// Solidity: function vaultRecords(address vaultAddress) view returns(bool exists, (bool,uint256) deregRequestOccurrence, ((uint96,uint160)[]) slashAmountHistory)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) VaultRecords(opts *bind.CallOpts, vaultAddress common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	SlashAmountHistory     CheckpointsTrace160
}, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "vaultRecords", vaultAddress)

	outstruct := new(struct {
		Exists                 bool
		DeregRequestOccurrence TimestampOccurrenceOccurrence
		SlashAmountHistory     CheckpointsTrace160
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[1], new(TimestampOccurrenceOccurrence)).(*TimestampOccurrenceOccurrence)
	outstruct.SlashAmountHistory = *abi.ConvertType(out[2], new(CheckpointsTrace160)).(*CheckpointsTrace160)

	return *outstruct, err

}

// VaultRecords is a free data retrieval call binding the contract method 0xdef52ca5.
//
// Solidity: function vaultRecords(address vaultAddress) view returns(bool exists, (bool,uint256) deregRequestOccurrence, ((uint96,uint160)[]) slashAmountHistory)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) VaultRecords(vaultAddress common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	SlashAmountHistory     CheckpointsTrace160
}, error) {
	return _Mevcommitmiddleware.Contract.VaultRecords(&_Mevcommitmiddleware.CallOpts, vaultAddress)
}

// VaultRecords is a free data retrieval call binding the contract method 0xdef52ca5.
//
// Solidity: function vaultRecords(address vaultAddress) view returns(bool exists, (bool,uint256) deregRequestOccurrence, ((uint96,uint160)[]) slashAmountHistory)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) VaultRecords(vaultAddress common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence TimestampOccurrenceOccurrence
	SlashAmountHistory     CheckpointsTrace160
}, error) {
	return _Mevcommitmiddleware.Contract.VaultRecords(&_Mevcommitmiddleware.CallOpts, vaultAddress)
}

// WouldVaultBeValidWith is a free data retrieval call binding the contract method 0x6b0dc2d3.
//
// Solidity: function wouldVaultBeValidWith(address vault, uint256 potentialSLashPeriodSeconds) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCaller) WouldVaultBeValidWith(opts *bind.CallOpts, vault common.Address, potentialSLashPeriodSeconds *big.Int) (bool, error) {
	var out []interface{}
	err := _Mevcommitmiddleware.contract.Call(opts, &out, "wouldVaultBeValidWith", vault, potentialSLashPeriodSeconds)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WouldVaultBeValidWith is a free data retrieval call binding the contract method 0x6b0dc2d3.
//
// Solidity: function wouldVaultBeValidWith(address vault, uint256 potentialSLashPeriodSeconds) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) WouldVaultBeValidWith(vault common.Address, potentialSLashPeriodSeconds *big.Int) (bool, error) {
	return _Mevcommitmiddleware.Contract.WouldVaultBeValidWith(&_Mevcommitmiddleware.CallOpts, vault, potentialSLashPeriodSeconds)
}

// WouldVaultBeValidWith is a free data retrieval call binding the contract method 0x6b0dc2d3.
//
// Solidity: function wouldVaultBeValidWith(address vault, uint256 potentialSLashPeriodSeconds) view returns(bool)
func (_Mevcommitmiddleware *MevcommitmiddlewareCallerSession) WouldVaultBeValidWith(vault common.Address, potentialSLashPeriodSeconds *big.Int) (bool, error) {
	return _Mevcommitmiddleware.Contract.WouldVaultBeValidWith(&_Mevcommitmiddleware.CallOpts, vault, potentialSLashPeriodSeconds)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.AcceptOwnership(&_Mevcommitmiddleware.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.AcceptOwnership(&_Mevcommitmiddleware.TransactOpts)
}

// BlacklistOperators is a paid mutator transaction binding the contract method 0xc999a751.
//
// Solidity: function blacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) BlacklistOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "blacklistOperators", operators)
}

// BlacklistOperators is a paid mutator transaction binding the contract method 0xc999a751.
//
// Solidity: function blacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) BlacklistOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.BlacklistOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// BlacklistOperators is a paid mutator transaction binding the contract method 0xc999a751.
//
// Solidity: function blacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) BlacklistOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.BlacklistOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// DeregisterOperators is a paid mutator transaction binding the contract method 0x4de702b7.
//
// Solidity: function deregisterOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) DeregisterOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "deregisterOperators", operators)
}

// DeregisterOperators is a paid mutator transaction binding the contract method 0x4de702b7.
//
// Solidity: function deregisterOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) DeregisterOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// DeregisterOperators is a paid mutator transaction binding the contract method 0x4de702b7.
//
// Solidity: function deregisterOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) DeregisterOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) DeregisterValidators(opts *bind.TransactOpts, blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "deregisterValidators", blsPubkeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) DeregisterValidators(blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) DeregisterValidators(blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys)
}

// DeregisterVaults is a paid mutator transaction binding the contract method 0xbb18f271.
//
// Solidity: function deregisterVaults(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) DeregisterVaults(opts *bind.TransactOpts, vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "deregisterVaults", vaults)
}

// DeregisterVaults is a paid mutator transaction binding the contract method 0xbb18f271.
//
// Solidity: function deregisterVaults(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) DeregisterVaults(vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterVaults(&_Mevcommitmiddleware.TransactOpts, vaults)
}

// DeregisterVaults is a paid mutator transaction binding the contract method 0xbb18f271.
//
// Solidity: function deregisterVaults(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) DeregisterVaults(vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.DeregisterVaults(&_Mevcommitmiddleware.TransactOpts, vaults)
}

// Initialize is a paid mutator transaction binding the contract method 0x4895b0a8.
//
// Solidity: function initialize(address _networkRegistry, address _operatorRegistry, address _vaultFactory, address _delegatorFactory, address _slasherFactory, address _burnerRouterFactory, address _network, uint256 _slashPeriodSeconds, address _slashOracle, address _slashReceiver, uint256 _minBurnerRouterDelay, address _owner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) Initialize(opts *bind.TransactOpts, _networkRegistry common.Address, _operatorRegistry common.Address, _vaultFactory common.Address, _delegatorFactory common.Address, _slasherFactory common.Address, _burnerRouterFactory common.Address, _network common.Address, _slashPeriodSeconds *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _minBurnerRouterDelay *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "initialize", _networkRegistry, _operatorRegistry, _vaultFactory, _delegatorFactory, _slasherFactory, _burnerRouterFactory, _network, _slashPeriodSeconds, _slashOracle, _slashReceiver, _minBurnerRouterDelay, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x4895b0a8.
//
// Solidity: function initialize(address _networkRegistry, address _operatorRegistry, address _vaultFactory, address _delegatorFactory, address _slasherFactory, address _burnerRouterFactory, address _network, uint256 _slashPeriodSeconds, address _slashOracle, address _slashReceiver, uint256 _minBurnerRouterDelay, address _owner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Initialize(_networkRegistry common.Address, _operatorRegistry common.Address, _vaultFactory common.Address, _delegatorFactory common.Address, _slasherFactory common.Address, _burnerRouterFactory common.Address, _network common.Address, _slashPeriodSeconds *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _minBurnerRouterDelay *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Initialize(&_Mevcommitmiddleware.TransactOpts, _networkRegistry, _operatorRegistry, _vaultFactory, _delegatorFactory, _slasherFactory, _burnerRouterFactory, _network, _slashPeriodSeconds, _slashOracle, _slashReceiver, _minBurnerRouterDelay, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x4895b0a8.
//
// Solidity: function initialize(address _networkRegistry, address _operatorRegistry, address _vaultFactory, address _delegatorFactory, address _slasherFactory, address _burnerRouterFactory, address _network, uint256 _slashPeriodSeconds, address _slashOracle, address _slashReceiver, uint256 _minBurnerRouterDelay, address _owner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) Initialize(_networkRegistry common.Address, _operatorRegistry common.Address, _vaultFactory common.Address, _delegatorFactory common.Address, _slasherFactory common.Address, _burnerRouterFactory common.Address, _network common.Address, _slashPeriodSeconds *big.Int, _slashOracle common.Address, _slashReceiver common.Address, _minBurnerRouterDelay *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Initialize(&_Mevcommitmiddleware.TransactOpts, _networkRegistry, _operatorRegistry, _vaultFactory, _delegatorFactory, _slasherFactory, _burnerRouterFactory, _network, _slashPeriodSeconds, _slashOracle, _slashReceiver, _minBurnerRouterDelay, _owner)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Pause() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Pause(&_Mevcommitmiddleware.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) Pause() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Pause(&_Mevcommitmiddleware.TransactOpts)
}

// RegisterOperators is a paid mutator transaction binding the contract method 0x3ec1418e.
//
// Solidity: function registerOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RegisterOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "registerOperators", operators)
}

// RegisterOperators is a paid mutator transaction binding the contract method 0x3ec1418e.
//
// Solidity: function registerOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RegisterOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// RegisterOperators is a paid mutator transaction binding the contract method 0x3ec1418e.
//
// Solidity: function registerOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RegisterOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0x8b7a8ea8.
//
// Solidity: function registerValidators(bytes[][] blsPubkeys, address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RegisterValidators(opts *bind.TransactOpts, blsPubkeys [][][]byte, vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "registerValidators", blsPubkeys, vaults)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0x8b7a8ea8.
//
// Solidity: function registerValidators(bytes[][] blsPubkeys, address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RegisterValidators(blsPubkeys [][][]byte, vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys, vaults)
}

// RegisterValidators is a paid mutator transaction binding the contract method 0x8b7a8ea8.
//
// Solidity: function registerValidators(bytes[][] blsPubkeys, address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RegisterValidators(blsPubkeys [][][]byte, vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys, vaults)
}

// RegisterVaults is a paid mutator transaction binding the contract method 0xebb0d875.
//
// Solidity: function registerVaults(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RegisterVaults(opts *bind.TransactOpts, vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "registerVaults", vaults, slashAmounts)
}

// RegisterVaults is a paid mutator transaction binding the contract method 0xebb0d875.
//
// Solidity: function registerVaults(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RegisterVaults(vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterVaults(&_Mevcommitmiddleware.TransactOpts, vaults, slashAmounts)
}

// RegisterVaults is a paid mutator transaction binding the contract method 0xebb0d875.
//
// Solidity: function registerVaults(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RegisterVaults(vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RegisterVaults(&_Mevcommitmiddleware.TransactOpts, vaults, slashAmounts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RenounceOwnership(&_Mevcommitmiddleware.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RenounceOwnership(&_Mevcommitmiddleware.TransactOpts)
}

// RequestOperatorDeregistrations is a paid mutator transaction binding the contract method 0x1e902bc4.
//
// Solidity: function requestOperatorDeregistrations(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RequestOperatorDeregistrations(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "requestOperatorDeregistrations", operators)
}

// RequestOperatorDeregistrations is a paid mutator transaction binding the contract method 0x1e902bc4.
//
// Solidity: function requestOperatorDeregistrations(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RequestOperatorDeregistrations(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestOperatorDeregistrations(&_Mevcommitmiddleware.TransactOpts, operators)
}

// RequestOperatorDeregistrations is a paid mutator transaction binding the contract method 0x1e902bc4.
//
// Solidity: function requestOperatorDeregistrations(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RequestOperatorDeregistrations(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestOperatorDeregistrations(&_Mevcommitmiddleware.TransactOpts, operators)
}

// RequestValDeregistrations is a paid mutator transaction binding the contract method 0x44119b2f.
//
// Solidity: function requestValDeregistrations(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RequestValDeregistrations(opts *bind.TransactOpts, blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "requestValDeregistrations", blsPubkeys)
}

// RequestValDeregistrations is a paid mutator transaction binding the contract method 0x44119b2f.
//
// Solidity: function requestValDeregistrations(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RequestValDeregistrations(blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestValDeregistrations(&_Mevcommitmiddleware.TransactOpts, blsPubkeys)
}

// RequestValDeregistrations is a paid mutator transaction binding the contract method 0x44119b2f.
//
// Solidity: function requestValDeregistrations(bytes[] blsPubkeys) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RequestValDeregistrations(blsPubkeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestValDeregistrations(&_Mevcommitmiddleware.TransactOpts, blsPubkeys)
}

// RequestVaultDeregistrations is a paid mutator transaction binding the contract method 0x0fdba156.
//
// Solidity: function requestVaultDeregistrations(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) RequestVaultDeregistrations(opts *bind.TransactOpts, vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "requestVaultDeregistrations", vaults)
}

// RequestVaultDeregistrations is a paid mutator transaction binding the contract method 0x0fdba156.
//
// Solidity: function requestVaultDeregistrations(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) RequestVaultDeregistrations(vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestVaultDeregistrations(&_Mevcommitmiddleware.TransactOpts, vaults)
}

// RequestVaultDeregistrations is a paid mutator transaction binding the contract method 0x0fdba156.
//
// Solidity: function requestVaultDeregistrations(address[] vaults) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) RequestVaultDeregistrations(vaults []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.RequestVaultDeregistrations(&_Mevcommitmiddleware.TransactOpts, vaults)
}

// SetBurnerRouterFactory is a paid mutator transaction binding the contract method 0xd0352521.
//
// Solidity: function setBurnerRouterFactory(address _burnerRouterFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetBurnerRouterFactory(opts *bind.TransactOpts, _burnerRouterFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setBurnerRouterFactory", _burnerRouterFactory)
}

// SetBurnerRouterFactory is a paid mutator transaction binding the contract method 0xd0352521.
//
// Solidity: function setBurnerRouterFactory(address _burnerRouterFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetBurnerRouterFactory(_burnerRouterFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetBurnerRouterFactory(&_Mevcommitmiddleware.TransactOpts, _burnerRouterFactory)
}

// SetBurnerRouterFactory is a paid mutator transaction binding the contract method 0xd0352521.
//
// Solidity: function setBurnerRouterFactory(address _burnerRouterFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetBurnerRouterFactory(_burnerRouterFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetBurnerRouterFactory(&_Mevcommitmiddleware.TransactOpts, _burnerRouterFactory)
}

// SetDelegatorFactory is a paid mutator transaction binding the contract method 0xfd44e64f.
//
// Solidity: function setDelegatorFactory(address _delegatorFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetDelegatorFactory(opts *bind.TransactOpts, _delegatorFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setDelegatorFactory", _delegatorFactory)
}

// SetDelegatorFactory is a paid mutator transaction binding the contract method 0xfd44e64f.
//
// Solidity: function setDelegatorFactory(address _delegatorFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetDelegatorFactory(_delegatorFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetDelegatorFactory(&_Mevcommitmiddleware.TransactOpts, _delegatorFactory)
}

// SetDelegatorFactory is a paid mutator transaction binding the contract method 0xfd44e64f.
//
// Solidity: function setDelegatorFactory(address _delegatorFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetDelegatorFactory(_delegatorFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetDelegatorFactory(&_Mevcommitmiddleware.TransactOpts, _delegatorFactory)
}

// SetMinBurnerRouterDelay is a paid mutator transaction binding the contract method 0x9c8c3022.
//
// Solidity: function setMinBurnerRouterDelay(uint256 minBurnerRouterDelay_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetMinBurnerRouterDelay(opts *bind.TransactOpts, minBurnerRouterDelay_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setMinBurnerRouterDelay", minBurnerRouterDelay_)
}

// SetMinBurnerRouterDelay is a paid mutator transaction binding the contract method 0x9c8c3022.
//
// Solidity: function setMinBurnerRouterDelay(uint256 minBurnerRouterDelay_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetMinBurnerRouterDelay(minBurnerRouterDelay_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetMinBurnerRouterDelay(&_Mevcommitmiddleware.TransactOpts, minBurnerRouterDelay_)
}

// SetMinBurnerRouterDelay is a paid mutator transaction binding the contract method 0x9c8c3022.
//
// Solidity: function setMinBurnerRouterDelay(uint256 minBurnerRouterDelay_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetMinBurnerRouterDelay(minBurnerRouterDelay_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetMinBurnerRouterDelay(&_Mevcommitmiddleware.TransactOpts, minBurnerRouterDelay_)
}

// SetNetwork is a paid mutator transaction binding the contract method 0xa1d71142.
//
// Solidity: function setNetwork(address _network) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetNetwork(opts *bind.TransactOpts, _network common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setNetwork", _network)
}

// SetNetwork is a paid mutator transaction binding the contract method 0xa1d71142.
//
// Solidity: function setNetwork(address _network) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetNetwork(_network common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetNetwork(&_Mevcommitmiddleware.TransactOpts, _network)
}

// SetNetwork is a paid mutator transaction binding the contract method 0xa1d71142.
//
// Solidity: function setNetwork(address _network) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetNetwork(_network common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetNetwork(&_Mevcommitmiddleware.TransactOpts, _network)
}

// SetNetworkRegistry is a paid mutator transaction binding the contract method 0x1a28acdd.
//
// Solidity: function setNetworkRegistry(address _networkRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetNetworkRegistry(opts *bind.TransactOpts, _networkRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setNetworkRegistry", _networkRegistry)
}

// SetNetworkRegistry is a paid mutator transaction binding the contract method 0x1a28acdd.
//
// Solidity: function setNetworkRegistry(address _networkRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetNetworkRegistry(_networkRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetNetworkRegistry(&_Mevcommitmiddleware.TransactOpts, _networkRegistry)
}

// SetNetworkRegistry is a paid mutator transaction binding the contract method 0x1a28acdd.
//
// Solidity: function setNetworkRegistry(address _networkRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetNetworkRegistry(_networkRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetNetworkRegistry(&_Mevcommitmiddleware.TransactOpts, _networkRegistry)
}

// SetOperatorRegistry is a paid mutator transaction binding the contract method 0x9d28fb86.
//
// Solidity: function setOperatorRegistry(address _operatorRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetOperatorRegistry(opts *bind.TransactOpts, _operatorRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setOperatorRegistry", _operatorRegistry)
}

// SetOperatorRegistry is a paid mutator transaction binding the contract method 0x9d28fb86.
//
// Solidity: function setOperatorRegistry(address _operatorRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetOperatorRegistry(_operatorRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetOperatorRegistry(&_Mevcommitmiddleware.TransactOpts, _operatorRegistry)
}

// SetOperatorRegistry is a paid mutator transaction binding the contract method 0x9d28fb86.
//
// Solidity: function setOperatorRegistry(address _operatorRegistry) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetOperatorRegistry(_operatorRegistry common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetOperatorRegistry(&_Mevcommitmiddleware.TransactOpts, _operatorRegistry)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address slashOracle_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetSlashOracle(opts *bind.TransactOpts, slashOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setSlashOracle", slashOracle_)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address slashOracle_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetSlashOracle(slashOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashOracle(&_Mevcommitmiddleware.TransactOpts, slashOracle_)
}

// SetSlashOracle is a paid mutator transaction binding the contract method 0x370baff6.
//
// Solidity: function setSlashOracle(address slashOracle_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetSlashOracle(slashOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashOracle(&_Mevcommitmiddleware.TransactOpts, slashOracle_)
}

// SetSlashPeriodSeconds is a paid mutator transaction binding the contract method 0x0dd92a24.
//
// Solidity: function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetSlashPeriodSeconds(opts *bind.TransactOpts, slashPeriodSeconds_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setSlashPeriodSeconds", slashPeriodSeconds_)
}

// SetSlashPeriodSeconds is a paid mutator transaction binding the contract method 0x0dd92a24.
//
// Solidity: function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetSlashPeriodSeconds(slashPeriodSeconds_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashPeriodSeconds(&_Mevcommitmiddleware.TransactOpts, slashPeriodSeconds_)
}

// SetSlashPeriodSeconds is a paid mutator transaction binding the contract method 0x0dd92a24.
//
// Solidity: function setSlashPeriodSeconds(uint256 slashPeriodSeconds_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetSlashPeriodSeconds(slashPeriodSeconds_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashPeriodSeconds(&_Mevcommitmiddleware.TransactOpts, slashPeriodSeconds_)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address slashReceiver_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetSlashReceiver(opts *bind.TransactOpts, slashReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setSlashReceiver", slashReceiver_)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address slashReceiver_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetSlashReceiver(slashReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashReceiver(&_Mevcommitmiddleware.TransactOpts, slashReceiver_)
}

// SetSlashReceiver is a paid mutator transaction binding the contract method 0x1a6933d5.
//
// Solidity: function setSlashReceiver(address slashReceiver_) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetSlashReceiver(slashReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlashReceiver(&_Mevcommitmiddleware.TransactOpts, slashReceiver_)
}

// SetSlasherFactory is a paid mutator transaction binding the contract method 0x7a8c82ea.
//
// Solidity: function setSlasherFactory(address _slasherFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetSlasherFactory(opts *bind.TransactOpts, _slasherFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setSlasherFactory", _slasherFactory)
}

// SetSlasherFactory is a paid mutator transaction binding the contract method 0x7a8c82ea.
//
// Solidity: function setSlasherFactory(address _slasherFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetSlasherFactory(_slasherFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlasherFactory(&_Mevcommitmiddleware.TransactOpts, _slasherFactory)
}

// SetSlasherFactory is a paid mutator transaction binding the contract method 0x7a8c82ea.
//
// Solidity: function setSlasherFactory(address _slasherFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetSlasherFactory(_slasherFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetSlasherFactory(&_Mevcommitmiddleware.TransactOpts, _slasherFactory)
}

// SetVaultFactory is a paid mutator transaction binding the contract method 0x3ea7fbdb.
//
// Solidity: function setVaultFactory(address _vaultFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SetVaultFactory(opts *bind.TransactOpts, _vaultFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "setVaultFactory", _vaultFactory)
}

// SetVaultFactory is a paid mutator transaction binding the contract method 0x3ea7fbdb.
//
// Solidity: function setVaultFactory(address _vaultFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SetVaultFactory(_vaultFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetVaultFactory(&_Mevcommitmiddleware.TransactOpts, _vaultFactory)
}

// SetVaultFactory is a paid mutator transaction binding the contract method 0x3ea7fbdb.
//
// Solidity: function setVaultFactory(address _vaultFactory) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SetVaultFactory(_vaultFactory common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SetVaultFactory(&_Mevcommitmiddleware.TransactOpts, _vaultFactory)
}

// SlashValidators is a paid mutator transaction binding the contract method 0xf62f76fa.
//
// Solidity: function slashValidators(bytes[] blsPubkeys, uint256[] captureTimestamps) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) SlashValidators(opts *bind.TransactOpts, blsPubkeys [][]byte, captureTimestamps []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "slashValidators", blsPubkeys, captureTimestamps)
}

// SlashValidators is a paid mutator transaction binding the contract method 0xf62f76fa.
//
// Solidity: function slashValidators(bytes[] blsPubkeys, uint256[] captureTimestamps) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) SlashValidators(blsPubkeys [][]byte, captureTimestamps []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SlashValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys, captureTimestamps)
}

// SlashValidators is a paid mutator transaction binding the contract method 0xf62f76fa.
//
// Solidity: function slashValidators(bytes[] blsPubkeys, uint256[] captureTimestamps) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) SlashValidators(blsPubkeys [][]byte, captureTimestamps []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.SlashValidators(&_Mevcommitmiddleware.TransactOpts, blsPubkeys, captureTimestamps)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.TransferOwnership(&_Mevcommitmiddleware.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.TransferOwnership(&_Mevcommitmiddleware.TransactOpts, newOwner)
}

// UnblacklistOperators is a paid mutator transaction binding the contract method 0x1a2b30d8.
//
// Solidity: function unblacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) UnblacklistOperators(opts *bind.TransactOpts, operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "unblacklistOperators", operators)
}

// UnblacklistOperators is a paid mutator transaction binding the contract method 0x1a2b30d8.
//
// Solidity: function unblacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) UnblacklistOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UnblacklistOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// UnblacklistOperators is a paid mutator transaction binding the contract method 0x1a2b30d8.
//
// Solidity: function unblacklistOperators(address[] operators) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) UnblacklistOperators(operators []common.Address) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UnblacklistOperators(&_Mevcommitmiddleware.TransactOpts, operators)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Unpause() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Unpause(&_Mevcommitmiddleware.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) Unpause() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Unpause(&_Mevcommitmiddleware.TransactOpts)
}

// UpdateSlashAmounts is a paid mutator transaction binding the contract method 0xafd7913d.
//
// Solidity: function updateSlashAmounts(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) UpdateSlashAmounts(opts *bind.TransactOpts, vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "updateSlashAmounts", vaults, slashAmounts)
}

// UpdateSlashAmounts is a paid mutator transaction binding the contract method 0xafd7913d.
//
// Solidity: function updateSlashAmounts(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) UpdateSlashAmounts(vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UpdateSlashAmounts(&_Mevcommitmiddleware.TransactOpts, vaults, slashAmounts)
}

// UpdateSlashAmounts is a paid mutator transaction binding the contract method 0xafd7913d.
//
// Solidity: function updateSlashAmounts(address[] vaults, uint160[] slashAmounts) returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) UpdateSlashAmounts(vaults []common.Address, slashAmounts []*big.Int) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UpdateSlashAmounts(&_Mevcommitmiddleware.TransactOpts, vaults, slashAmounts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UpgradeToAndCall(&_Mevcommitmiddleware.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.UpgradeToAndCall(&_Mevcommitmiddleware.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Fallback(&_Mevcommitmiddleware.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Fallback(&_Mevcommitmiddleware.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitmiddleware.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareSession) Receive() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Receive(&_Mevcommitmiddleware.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitmiddleware *MevcommitmiddlewareTransactorSession) Receive() (*types.Transaction, error) {
	return _Mevcommitmiddleware.Contract.Receive(&_Mevcommitmiddleware.TransactOpts)
}

// MevcommitmiddlewareBurnerRouterFactorySetIterator is returned from FilterBurnerRouterFactorySet and is used to iterate over the raw logs and unpacked data for BurnerRouterFactorySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareBurnerRouterFactorySetIterator struct {
	Event *MevcommitmiddlewareBurnerRouterFactorySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareBurnerRouterFactorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareBurnerRouterFactorySet)
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
		it.Event = new(MevcommitmiddlewareBurnerRouterFactorySet)
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
func (it *MevcommitmiddlewareBurnerRouterFactorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareBurnerRouterFactorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareBurnerRouterFactorySet represents a BurnerRouterFactorySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareBurnerRouterFactorySet struct {
	BurnerRouterFactory common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterBurnerRouterFactorySet is a free log retrieval operation binding the contract event 0x9c214079845899d94b3bd881a14e996ebd153ef99fdc98ee4681eacf19c62f38.
//
// Solidity: event BurnerRouterFactorySet(address burnerRouterFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterBurnerRouterFactorySet(opts *bind.FilterOpts) (*MevcommitmiddlewareBurnerRouterFactorySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "BurnerRouterFactorySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareBurnerRouterFactorySetIterator{contract: _Mevcommitmiddleware.contract, event: "BurnerRouterFactorySet", logs: logs, sub: sub}, nil
}

// WatchBurnerRouterFactorySet is a free log subscription operation binding the contract event 0x9c214079845899d94b3bd881a14e996ebd153ef99fdc98ee4681eacf19c62f38.
//
// Solidity: event BurnerRouterFactorySet(address burnerRouterFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchBurnerRouterFactorySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareBurnerRouterFactorySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "BurnerRouterFactorySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareBurnerRouterFactorySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "BurnerRouterFactorySet", log); err != nil {
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

// ParseBurnerRouterFactorySet is a log parse operation binding the contract event 0x9c214079845899d94b3bd881a14e996ebd153ef99fdc98ee4681eacf19c62f38.
//
// Solidity: event BurnerRouterFactorySet(address burnerRouterFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseBurnerRouterFactorySet(log types.Log) (*MevcommitmiddlewareBurnerRouterFactorySet, error) {
	event := new(MevcommitmiddlewareBurnerRouterFactorySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "BurnerRouterFactorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareDelegatorFactorySetIterator is returned from FilterDelegatorFactorySet and is used to iterate over the raw logs and unpacked data for DelegatorFactorySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareDelegatorFactorySetIterator struct {
	Event *MevcommitmiddlewareDelegatorFactorySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareDelegatorFactorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareDelegatorFactorySet)
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
		it.Event = new(MevcommitmiddlewareDelegatorFactorySet)
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
func (it *MevcommitmiddlewareDelegatorFactorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareDelegatorFactorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareDelegatorFactorySet represents a DelegatorFactorySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareDelegatorFactorySet struct {
	DelegatorFactory common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDelegatorFactorySet is a free log retrieval operation binding the contract event 0xdaafeb874ae63f130db762ae2fa1223a8384825723dbe67d3647071a8d499c39.
//
// Solidity: event DelegatorFactorySet(address delegatorFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterDelegatorFactorySet(opts *bind.FilterOpts) (*MevcommitmiddlewareDelegatorFactorySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "DelegatorFactorySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareDelegatorFactorySetIterator{contract: _Mevcommitmiddleware.contract, event: "DelegatorFactorySet", logs: logs, sub: sub}, nil
}

// WatchDelegatorFactorySet is a free log subscription operation binding the contract event 0xdaafeb874ae63f130db762ae2fa1223a8384825723dbe67d3647071a8d499c39.
//
// Solidity: event DelegatorFactorySet(address delegatorFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchDelegatorFactorySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareDelegatorFactorySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "DelegatorFactorySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareDelegatorFactorySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "DelegatorFactorySet", log); err != nil {
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

// ParseDelegatorFactorySet is a log parse operation binding the contract event 0xdaafeb874ae63f130db762ae2fa1223a8384825723dbe67d3647071a8d499c39.
//
// Solidity: event DelegatorFactorySet(address delegatorFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseDelegatorFactorySet(log types.Log) (*MevcommitmiddlewareDelegatorFactorySet, error) {
	event := new(MevcommitmiddlewareDelegatorFactorySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "DelegatorFactorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareInitializedIterator struct {
	Event *MevcommitmiddlewareInitialized // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareInitialized)
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
		it.Event = new(MevcommitmiddlewareInitialized)
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
func (it *MevcommitmiddlewareInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareInitialized represents a Initialized event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterInitialized(opts *bind.FilterOpts) (*MevcommitmiddlewareInitializedIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareInitializedIterator{contract: _Mevcommitmiddleware.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareInitialized) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareInitialized)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseInitialized(log types.Log) (*MevcommitmiddlewareInitialized, error) {
	event := new(MevcommitmiddlewareInitialized)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareMinBurnerRouterDelaySetIterator is returned from FilterMinBurnerRouterDelaySet and is used to iterate over the raw logs and unpacked data for MinBurnerRouterDelaySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareMinBurnerRouterDelaySetIterator struct {
	Event *MevcommitmiddlewareMinBurnerRouterDelaySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareMinBurnerRouterDelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareMinBurnerRouterDelaySet)
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
		it.Event = new(MevcommitmiddlewareMinBurnerRouterDelaySet)
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
func (it *MevcommitmiddlewareMinBurnerRouterDelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareMinBurnerRouterDelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareMinBurnerRouterDelaySet represents a MinBurnerRouterDelaySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareMinBurnerRouterDelaySet struct {
	MinBurnerRouterDelay *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterMinBurnerRouterDelaySet is a free log retrieval operation binding the contract event 0x1692318fc65bb45681dfd5556f2d21be7659c6809c725f503ad03240998b4b18.
//
// Solidity: event MinBurnerRouterDelaySet(uint256 minBurnerRouterDelay)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterMinBurnerRouterDelaySet(opts *bind.FilterOpts) (*MevcommitmiddlewareMinBurnerRouterDelaySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "MinBurnerRouterDelaySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareMinBurnerRouterDelaySetIterator{contract: _Mevcommitmiddleware.contract, event: "MinBurnerRouterDelaySet", logs: logs, sub: sub}, nil
}

// WatchMinBurnerRouterDelaySet is a free log subscription operation binding the contract event 0x1692318fc65bb45681dfd5556f2d21be7659c6809c725f503ad03240998b4b18.
//
// Solidity: event MinBurnerRouterDelaySet(uint256 minBurnerRouterDelay)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchMinBurnerRouterDelaySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareMinBurnerRouterDelaySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "MinBurnerRouterDelaySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareMinBurnerRouterDelaySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "MinBurnerRouterDelaySet", log); err != nil {
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

// ParseMinBurnerRouterDelaySet is a log parse operation binding the contract event 0x1692318fc65bb45681dfd5556f2d21be7659c6809c725f503ad03240998b4b18.
//
// Solidity: event MinBurnerRouterDelaySet(uint256 minBurnerRouterDelay)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseMinBurnerRouterDelaySet(log types.Log) (*MevcommitmiddlewareMinBurnerRouterDelaySet, error) {
	event := new(MevcommitmiddlewareMinBurnerRouterDelaySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "MinBurnerRouterDelaySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareNetworkRegistrySetIterator is returned from FilterNetworkRegistrySet and is used to iterate over the raw logs and unpacked data for NetworkRegistrySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareNetworkRegistrySetIterator struct {
	Event *MevcommitmiddlewareNetworkRegistrySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareNetworkRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareNetworkRegistrySet)
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
		it.Event = new(MevcommitmiddlewareNetworkRegistrySet)
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
func (it *MevcommitmiddlewareNetworkRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareNetworkRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareNetworkRegistrySet represents a NetworkRegistrySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareNetworkRegistrySet struct {
	NetworkRegistry common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNetworkRegistrySet is a free log retrieval operation binding the contract event 0xa1dc9b46f16ec4a543a99a0417fbd53dfb61ea14a947e937059e7ea877df1127.
//
// Solidity: event NetworkRegistrySet(address networkRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterNetworkRegistrySet(opts *bind.FilterOpts) (*MevcommitmiddlewareNetworkRegistrySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "NetworkRegistrySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareNetworkRegistrySetIterator{contract: _Mevcommitmiddleware.contract, event: "NetworkRegistrySet", logs: logs, sub: sub}, nil
}

// WatchNetworkRegistrySet is a free log subscription operation binding the contract event 0xa1dc9b46f16ec4a543a99a0417fbd53dfb61ea14a947e937059e7ea877df1127.
//
// Solidity: event NetworkRegistrySet(address networkRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchNetworkRegistrySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareNetworkRegistrySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "NetworkRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareNetworkRegistrySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "NetworkRegistrySet", log); err != nil {
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

// ParseNetworkRegistrySet is a log parse operation binding the contract event 0xa1dc9b46f16ec4a543a99a0417fbd53dfb61ea14a947e937059e7ea877df1127.
//
// Solidity: event NetworkRegistrySet(address networkRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseNetworkRegistrySet(log types.Log) (*MevcommitmiddlewareNetworkRegistrySet, error) {
	event := new(MevcommitmiddlewareNetworkRegistrySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "NetworkRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareNetworkSetIterator is returned from FilterNetworkSet and is used to iterate over the raw logs and unpacked data for NetworkSet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareNetworkSetIterator struct {
	Event *MevcommitmiddlewareNetworkSet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareNetworkSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareNetworkSet)
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
		it.Event = new(MevcommitmiddlewareNetworkSet)
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
func (it *MevcommitmiddlewareNetworkSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareNetworkSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareNetworkSet represents a NetworkSet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareNetworkSet struct {
	Network common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterNetworkSet is a free log retrieval operation binding the contract event 0x1dc2d95ebfe3aad1b48e66ae0f9aa8bfd5bef4ed9d321eba19879734ceb7b2a8.
//
// Solidity: event NetworkSet(address network)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterNetworkSet(opts *bind.FilterOpts) (*MevcommitmiddlewareNetworkSetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "NetworkSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareNetworkSetIterator{contract: _Mevcommitmiddleware.contract, event: "NetworkSet", logs: logs, sub: sub}, nil
}

// WatchNetworkSet is a free log subscription operation binding the contract event 0x1dc2d95ebfe3aad1b48e66ae0f9aa8bfd5bef4ed9d321eba19879734ceb7b2a8.
//
// Solidity: event NetworkSet(address network)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchNetworkSet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareNetworkSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "NetworkSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareNetworkSet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "NetworkSet", log); err != nil {
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

// ParseNetworkSet is a log parse operation binding the contract event 0x1dc2d95ebfe3aad1b48e66ae0f9aa8bfd5bef4ed9d321eba19879734ceb7b2a8.
//
// Solidity: event NetworkSet(address network)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseNetworkSet(log types.Log) (*MevcommitmiddlewareNetworkSet, error) {
	event := new(MevcommitmiddlewareNetworkSet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "NetworkSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorBlacklistedIterator is returned from FilterOperatorBlacklisted and is used to iterate over the raw logs and unpacked data for OperatorBlacklisted events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorBlacklistedIterator struct {
	Event *MevcommitmiddlewareOperatorBlacklisted // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorBlacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorBlacklisted)
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
		it.Event = new(MevcommitmiddlewareOperatorBlacklisted)
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
func (it *MevcommitmiddlewareOperatorBlacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorBlacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorBlacklisted represents a OperatorBlacklisted event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorBlacklisted struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorBlacklisted is a free log retrieval operation binding the contract event 0xe8a03158453769005cd5febcdae802c94cf43866ea101c3cf42cf42c140a4895.
//
// Solidity: event OperatorBlacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorBlacklisted(opts *bind.FilterOpts, operator []common.Address) (*MevcommitmiddlewareOperatorBlacklistedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorBlacklisted", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorBlacklistedIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorBlacklisted", logs: logs, sub: sub}, nil
}

// WatchOperatorBlacklisted is a free log subscription operation binding the contract event 0xe8a03158453769005cd5febcdae802c94cf43866ea101c3cf42cf42c140a4895.
//
// Solidity: event OperatorBlacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorBlacklisted(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorBlacklisted, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorBlacklisted", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorBlacklisted)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorBlacklisted", log); err != nil {
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

// ParseOperatorBlacklisted is a log parse operation binding the contract event 0xe8a03158453769005cd5febcdae802c94cf43866ea101c3cf42cf42c140a4895.
//
// Solidity: event OperatorBlacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorBlacklisted(log types.Log) (*MevcommitmiddlewareOperatorBlacklisted, error) {
	event := new(MevcommitmiddlewareOperatorBlacklisted)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorBlacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorDeregisteredIterator is returned from FilterOperatorDeregistered and is used to iterate over the raw logs and unpacked data for OperatorDeregistered events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorDeregisteredIterator struct {
	Event *MevcommitmiddlewareOperatorDeregistered // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorDeregistered)
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
		it.Event = new(MevcommitmiddlewareOperatorDeregistered)
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
func (it *MevcommitmiddlewareOperatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorDeregistered represents a OperatorDeregistered event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorDeregistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistered is a free log retrieval operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorDeregistered(opts *bind.FilterOpts, operator []common.Address) (*MevcommitmiddlewareOperatorDeregisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorDeregisteredIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistered is a free log subscription operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorDeregistered(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorDeregistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorDeregistered)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
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

// ParseOperatorDeregistered is a log parse operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorDeregistered(log types.Log) (*MevcommitmiddlewareOperatorDeregistered, error) {
	event := new(MevcommitmiddlewareOperatorDeregistered)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorDeregistrationRequestedIterator is returned from FilterOperatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for OperatorDeregistrationRequested events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorDeregistrationRequestedIterator struct {
	Event *MevcommitmiddlewareOperatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorDeregistrationRequested)
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
		it.Event = new(MevcommitmiddlewareOperatorDeregistrationRequested)
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
func (it *MevcommitmiddlewareOperatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorDeregistrationRequested represents a OperatorDeregistrationRequested event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorDeregistrationRequested struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorDeregistrationRequested(opts *bind.FilterOpts, operator []common.Address) (*MevcommitmiddlewareOperatorDeregistrationRequestedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorDeregistrationRequestedIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistrationRequested is a free log subscription operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorDeregistrationRequested, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorDeregistrationRequested)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
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

// ParseOperatorDeregistrationRequested is a log parse operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorDeregistrationRequested(log types.Log) (*MevcommitmiddlewareOperatorDeregistrationRequested, error) {
	event := new(MevcommitmiddlewareOperatorDeregistrationRequested)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorRegisteredIterator struct {
	Event *MevcommitmiddlewareOperatorRegistered // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorRegistered)
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
		it.Event = new(MevcommitmiddlewareOperatorRegistered)
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
func (it *MevcommitmiddlewareOperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorRegistered represents a OperatorRegistered event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorRegistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address) (*MevcommitmiddlewareOperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorRegisteredIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorRegistered)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
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

// ParseOperatorRegistered is a log parse operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorRegistered(log types.Log) (*MevcommitmiddlewareOperatorRegistered, error) {
	event := new(MevcommitmiddlewareOperatorRegistered)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorRegistrySetIterator is returned from FilterOperatorRegistrySet and is used to iterate over the raw logs and unpacked data for OperatorRegistrySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorRegistrySetIterator struct {
	Event *MevcommitmiddlewareOperatorRegistrySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorRegistrySet)
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
		it.Event = new(MevcommitmiddlewareOperatorRegistrySet)
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
func (it *MevcommitmiddlewareOperatorRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorRegistrySet represents a OperatorRegistrySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorRegistrySet struct {
	OperatorRegistry common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistrySet is a free log retrieval operation binding the contract event 0xc6df119c56c99171b170652a3c4750ba46dcaacbdb3b7ab4847a9fa339659bd4.
//
// Solidity: event OperatorRegistrySet(address operatorRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorRegistrySet(opts *bind.FilterOpts) (*MevcommitmiddlewareOperatorRegistrySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorRegistrySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorRegistrySetIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorRegistrySet", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistrySet is a free log subscription operation binding the contract event 0xc6df119c56c99171b170652a3c4750ba46dcaacbdb3b7ab4847a9fa339659bd4.
//
// Solidity: event OperatorRegistrySet(address operatorRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorRegistrySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorRegistrySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorRegistrySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorRegistrySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorRegistrySet", log); err != nil {
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

// ParseOperatorRegistrySet is a log parse operation binding the contract event 0xc6df119c56c99171b170652a3c4750ba46dcaacbdb3b7ab4847a9fa339659bd4.
//
// Solidity: event OperatorRegistrySet(address operatorRegistry)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorRegistrySet(log types.Log) (*MevcommitmiddlewareOperatorRegistrySet, error) {
	event := new(MevcommitmiddlewareOperatorRegistrySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorRegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOperatorUnblacklistedIterator is returned from FilterOperatorUnblacklisted and is used to iterate over the raw logs and unpacked data for OperatorUnblacklisted events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorUnblacklistedIterator struct {
	Event *MevcommitmiddlewareOperatorUnblacklisted // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOperatorUnblacklistedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOperatorUnblacklisted)
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
		it.Event = new(MevcommitmiddlewareOperatorUnblacklisted)
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
func (it *MevcommitmiddlewareOperatorUnblacklistedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOperatorUnblacklistedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOperatorUnblacklisted represents a OperatorUnblacklisted event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOperatorUnblacklisted struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorUnblacklisted is a free log retrieval operation binding the contract event 0x30f82d25b70616233efbf2847113c6c8f9e3af6ac8cd01bbeaf2096ad7a9304f.
//
// Solidity: event OperatorUnblacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOperatorUnblacklisted(opts *bind.FilterOpts, operator []common.Address) (*MevcommitmiddlewareOperatorUnblacklistedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OperatorUnblacklisted", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOperatorUnblacklistedIterator{contract: _Mevcommitmiddleware.contract, event: "OperatorUnblacklisted", logs: logs, sub: sub}, nil
}

// WatchOperatorUnblacklisted is a free log subscription operation binding the contract event 0x30f82d25b70616233efbf2847113c6c8f9e3af6ac8cd01bbeaf2096ad7a9304f.
//
// Solidity: event OperatorUnblacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOperatorUnblacklisted(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOperatorUnblacklisted, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OperatorUnblacklisted", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOperatorUnblacklisted)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorUnblacklisted", log); err != nil {
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

// ParseOperatorUnblacklisted is a log parse operation binding the contract event 0x30f82d25b70616233efbf2847113c6c8f9e3af6ac8cd01bbeaf2096ad7a9304f.
//
// Solidity: event OperatorUnblacklisted(address indexed operator)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOperatorUnblacklisted(log types.Log) (*MevcommitmiddlewareOperatorUnblacklisted, error) {
	event := new(MevcommitmiddlewareOperatorUnblacklisted)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OperatorUnblacklisted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOwnershipTransferStartedIterator struct {
	Event *MevcommitmiddlewareOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOwnershipTransferStarted)
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
		it.Event = new(MevcommitmiddlewareOwnershipTransferStarted)
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
func (it *MevcommitmiddlewareOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MevcommitmiddlewareOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOwnershipTransferStartedIterator{contract: _Mevcommitmiddleware.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOwnershipTransferStarted)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOwnershipTransferStarted(log types.Log) (*MevcommitmiddlewareOwnershipTransferStarted, error) {
	event := new(MevcommitmiddlewareOwnershipTransferStarted)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOwnershipTransferredIterator struct {
	Event *MevcommitmiddlewareOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareOwnershipTransferred)
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
		it.Event = new(MevcommitmiddlewareOwnershipTransferred)
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
func (it *MevcommitmiddlewareOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareOwnershipTransferred represents a OwnershipTransferred event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MevcommitmiddlewareOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareOwnershipTransferredIterator{contract: _Mevcommitmiddleware.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareOwnershipTransferred)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseOwnershipTransferred(log types.Log) (*MevcommitmiddlewareOwnershipTransferred, error) {
	event := new(MevcommitmiddlewareOwnershipTransferred)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewarePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewarePausedIterator struct {
	Event *MevcommitmiddlewarePaused // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewarePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewarePaused)
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
		it.Event = new(MevcommitmiddlewarePaused)
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
func (it *MevcommitmiddlewarePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewarePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewarePaused represents a Paused event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewarePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterPaused(opts *bind.FilterOpts) (*MevcommitmiddlewarePausedIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewarePausedIterator{contract: _Mevcommitmiddleware.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewarePaused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewarePaused)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParsePaused(log types.Log) (*MevcommitmiddlewarePaused, error) {
	event := new(MevcommitmiddlewarePaused)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareSlashOracleSetIterator is returned from FilterSlashOracleSet and is used to iterate over the raw logs and unpacked data for SlashOracleSet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashOracleSetIterator struct {
	Event *MevcommitmiddlewareSlashOracleSet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareSlashOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareSlashOracleSet)
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
		it.Event = new(MevcommitmiddlewareSlashOracleSet)
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
func (it *MevcommitmiddlewareSlashOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareSlashOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareSlashOracleSet represents a SlashOracleSet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashOracleSet struct {
	SlashOracle common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSlashOracleSet is a free log retrieval operation binding the contract event 0x96ffdc22b636403256c918edb15322ea5048fc70e5819a3a482db893e22e7cd1.
//
// Solidity: event SlashOracleSet(address slashOracle)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterSlashOracleSet(opts *bind.FilterOpts) (*MevcommitmiddlewareSlashOracleSetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "SlashOracleSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareSlashOracleSetIterator{contract: _Mevcommitmiddleware.contract, event: "SlashOracleSet", logs: logs, sub: sub}, nil
}

// WatchSlashOracleSet is a free log subscription operation binding the contract event 0x96ffdc22b636403256c918edb15322ea5048fc70e5819a3a482db893e22e7cd1.
//
// Solidity: event SlashOracleSet(address slashOracle)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchSlashOracleSet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareSlashOracleSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "SlashOracleSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareSlashOracleSet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
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

// ParseSlashOracleSet is a log parse operation binding the contract event 0x96ffdc22b636403256c918edb15322ea5048fc70e5819a3a482db893e22e7cd1.
//
// Solidity: event SlashOracleSet(address slashOracle)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseSlashOracleSet(log types.Log) (*MevcommitmiddlewareSlashOracleSet, error) {
	event := new(MevcommitmiddlewareSlashOracleSet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareSlashPeriodBlocksSetIterator is returned from FilterSlashPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for SlashPeriodBlocksSet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashPeriodBlocksSetIterator struct {
	Event *MevcommitmiddlewareSlashPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareSlashPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareSlashPeriodBlocksSet)
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
		it.Event = new(MevcommitmiddlewareSlashPeriodBlocksSet)
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
func (it *MevcommitmiddlewareSlashPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareSlashPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareSlashPeriodBlocksSet represents a SlashPeriodBlocksSet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashPeriodBlocksSet struct {
	SlashPeriodBlocks *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSlashPeriodBlocksSet is a free log retrieval operation binding the contract event 0x8836130fb6014984997af7d9e5f19c1d219b67679fe6b868e5a207bc989254ad.
//
// Solidity: event SlashPeriodBlocksSet(uint256 slashPeriodBlocks)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterSlashPeriodBlocksSet(opts *bind.FilterOpts) (*MevcommitmiddlewareSlashPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "SlashPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareSlashPeriodBlocksSetIterator{contract: _Mevcommitmiddleware.contract, event: "SlashPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchSlashPeriodBlocksSet is a free log subscription operation binding the contract event 0x8836130fb6014984997af7d9e5f19c1d219b67679fe6b868e5a207bc989254ad.
//
// Solidity: event SlashPeriodBlocksSet(uint256 slashPeriodBlocks)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchSlashPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareSlashPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "SlashPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareSlashPeriodBlocksSet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashPeriodBlocksSet", log); err != nil {
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

// ParseSlashPeriodBlocksSet is a log parse operation binding the contract event 0x8836130fb6014984997af7d9e5f19c1d219b67679fe6b868e5a207bc989254ad.
//
// Solidity: event SlashPeriodBlocksSet(uint256 slashPeriodBlocks)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseSlashPeriodBlocksSet(log types.Log) (*MevcommitmiddlewareSlashPeriodBlocksSet, error) {
	event := new(MevcommitmiddlewareSlashPeriodBlocksSet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareSlashPeriodSecondsSetIterator is returned from FilterSlashPeriodSecondsSet and is used to iterate over the raw logs and unpacked data for SlashPeriodSecondsSet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashPeriodSecondsSetIterator struct {
	Event *MevcommitmiddlewareSlashPeriodSecondsSet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareSlashPeriodSecondsSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareSlashPeriodSecondsSet)
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
		it.Event = new(MevcommitmiddlewareSlashPeriodSecondsSet)
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
func (it *MevcommitmiddlewareSlashPeriodSecondsSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareSlashPeriodSecondsSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareSlashPeriodSecondsSet represents a SlashPeriodSecondsSet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashPeriodSecondsSet struct {
	SlashPeriodSeconds *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSlashPeriodSecondsSet is a free log retrieval operation binding the contract event 0xac5cea39ce565630c66a23d9bc870dea5693cdf6971096bf680ebaa05a4fbbc2.
//
// Solidity: event SlashPeriodSecondsSet(uint256 slashPeriodSeconds)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterSlashPeriodSecondsSet(opts *bind.FilterOpts) (*MevcommitmiddlewareSlashPeriodSecondsSetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "SlashPeriodSecondsSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareSlashPeriodSecondsSetIterator{contract: _Mevcommitmiddleware.contract, event: "SlashPeriodSecondsSet", logs: logs, sub: sub}, nil
}

// WatchSlashPeriodSecondsSet is a free log subscription operation binding the contract event 0xac5cea39ce565630c66a23d9bc870dea5693cdf6971096bf680ebaa05a4fbbc2.
//
// Solidity: event SlashPeriodSecondsSet(uint256 slashPeriodSeconds)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchSlashPeriodSecondsSet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareSlashPeriodSecondsSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "SlashPeriodSecondsSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareSlashPeriodSecondsSet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashPeriodSecondsSet", log); err != nil {
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

// ParseSlashPeriodSecondsSet is a log parse operation binding the contract event 0xac5cea39ce565630c66a23d9bc870dea5693cdf6971096bf680ebaa05a4fbbc2.
//
// Solidity: event SlashPeriodSecondsSet(uint256 slashPeriodSeconds)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseSlashPeriodSecondsSet(log types.Log) (*MevcommitmiddlewareSlashPeriodSecondsSet, error) {
	event := new(MevcommitmiddlewareSlashPeriodSecondsSet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashPeriodSecondsSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareSlashReceiverSetIterator is returned from FilterSlashReceiverSet and is used to iterate over the raw logs and unpacked data for SlashReceiverSet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashReceiverSetIterator struct {
	Event *MevcommitmiddlewareSlashReceiverSet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareSlashReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareSlashReceiverSet)
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
		it.Event = new(MevcommitmiddlewareSlashReceiverSet)
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
func (it *MevcommitmiddlewareSlashReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareSlashReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareSlashReceiverSet represents a SlashReceiverSet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlashReceiverSet struct {
	SlashReceiver common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSlashReceiverSet is a free log retrieval operation binding the contract event 0x1299adc355601236099496a2d2b5de5dadc58ddf4628bd4e9ca3d2560931f9a6.
//
// Solidity: event SlashReceiverSet(address slashReceiver)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterSlashReceiverSet(opts *bind.FilterOpts) (*MevcommitmiddlewareSlashReceiverSetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "SlashReceiverSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareSlashReceiverSetIterator{contract: _Mevcommitmiddleware.contract, event: "SlashReceiverSet", logs: logs, sub: sub}, nil
}

// WatchSlashReceiverSet is a free log subscription operation binding the contract event 0x1299adc355601236099496a2d2b5de5dadc58ddf4628bd4e9ca3d2560931f9a6.
//
// Solidity: event SlashReceiverSet(address slashReceiver)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchSlashReceiverSet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareSlashReceiverSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "SlashReceiverSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareSlashReceiverSet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
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

// ParseSlashReceiverSet is a log parse operation binding the contract event 0x1299adc355601236099496a2d2b5de5dadc58ddf4628bd4e9ca3d2560931f9a6.
//
// Solidity: event SlashReceiverSet(address slashReceiver)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseSlashReceiverSet(log types.Log) (*MevcommitmiddlewareSlashReceiverSet, error) {
	event := new(MevcommitmiddlewareSlashReceiverSet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlashReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareSlasherFactorySetIterator is returned from FilterSlasherFactorySet and is used to iterate over the raw logs and unpacked data for SlasherFactorySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlasherFactorySetIterator struct {
	Event *MevcommitmiddlewareSlasherFactorySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareSlasherFactorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareSlasherFactorySet)
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
		it.Event = new(MevcommitmiddlewareSlasherFactorySet)
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
func (it *MevcommitmiddlewareSlasherFactorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareSlasherFactorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareSlasherFactorySet represents a SlasherFactorySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareSlasherFactorySet struct {
	SlasherFactory common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSlasherFactorySet is a free log retrieval operation binding the contract event 0xb936d33bb7fde9011d48a05202e13ff9e6fb19ea53ba9650ac7a6fd84235cd80.
//
// Solidity: event SlasherFactorySet(address slasherFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterSlasherFactorySet(opts *bind.FilterOpts) (*MevcommitmiddlewareSlasherFactorySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "SlasherFactorySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareSlasherFactorySetIterator{contract: _Mevcommitmiddleware.contract, event: "SlasherFactorySet", logs: logs, sub: sub}, nil
}

// WatchSlasherFactorySet is a free log subscription operation binding the contract event 0xb936d33bb7fde9011d48a05202e13ff9e6fb19ea53ba9650ac7a6fd84235cd80.
//
// Solidity: event SlasherFactorySet(address slasherFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchSlasherFactorySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareSlasherFactorySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "SlasherFactorySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareSlasherFactorySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlasherFactorySet", log); err != nil {
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

// ParseSlasherFactorySet is a log parse operation binding the contract event 0xb936d33bb7fde9011d48a05202e13ff9e6fb19ea53ba9650ac7a6fd84235cd80.
//
// Solidity: event SlasherFactorySet(address slasherFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseSlasherFactorySet(log types.Log) (*MevcommitmiddlewareSlasherFactorySet, error) {
	event := new(MevcommitmiddlewareSlasherFactorySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "SlasherFactorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareUnpausedIterator struct {
	Event *MevcommitmiddlewareUnpaused // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareUnpaused)
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
		it.Event = new(MevcommitmiddlewareUnpaused)
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
func (it *MevcommitmiddlewareUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareUnpaused represents a Unpaused event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MevcommitmiddlewareUnpausedIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareUnpausedIterator{contract: _Mevcommitmiddleware.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareUnpaused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareUnpaused)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseUnpaused(log types.Log) (*MevcommitmiddlewareUnpaused, error) {
	event := new(MevcommitmiddlewareUnpaused)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareUpgradedIterator struct {
	Event *MevcommitmiddlewareUpgraded // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareUpgraded)
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
		it.Event = new(MevcommitmiddlewareUpgraded)
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
func (it *MevcommitmiddlewareUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareUpgraded represents a Upgraded event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MevcommitmiddlewareUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareUpgradedIterator{contract: _Mevcommitmiddleware.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareUpgraded)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseUpgraded(log types.Log) (*MevcommitmiddlewareUpgraded, error) {
	event := new(MevcommitmiddlewareUpgraded)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareValRecordAddedIterator is returned from FilterValRecordAdded and is used to iterate over the raw logs and unpacked data for ValRecordAdded events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValRecordAddedIterator struct {
	Event *MevcommitmiddlewareValRecordAdded // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareValRecordAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareValRecordAdded)
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
		it.Event = new(MevcommitmiddlewareValRecordAdded)
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
func (it *MevcommitmiddlewareValRecordAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareValRecordAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareValRecordAdded represents a ValRecordAdded event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValRecordAdded struct {
	BlsPubkey []byte
	Operator  common.Address
	Vault     common.Address
	Position  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValRecordAdded is a free log retrieval operation binding the contract event 0x292d0cd97c126d79e15cbccfedd21071c6230c185049f19a1c84f92a85ec0413.
//
// Solidity: event ValRecordAdded(bytes blsPubkey, address indexed operator, address indexed vault, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterValRecordAdded(opts *bind.FilterOpts, operator []common.Address, vault []common.Address, position []*big.Int) (*MevcommitmiddlewareValRecordAddedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}
	var positionRule []interface{}
	for _, positionItem := range position {
		positionRule = append(positionRule, positionItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "ValRecordAdded", operatorRule, vaultRule, positionRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareValRecordAddedIterator{contract: _Mevcommitmiddleware.contract, event: "ValRecordAdded", logs: logs, sub: sub}, nil
}

// WatchValRecordAdded is a free log subscription operation binding the contract event 0x292d0cd97c126d79e15cbccfedd21071c6230c185049f19a1c84f92a85ec0413.
//
// Solidity: event ValRecordAdded(bytes blsPubkey, address indexed operator, address indexed vault, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchValRecordAdded(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareValRecordAdded, operator []common.Address, vault []common.Address, position []*big.Int) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}
	var positionRule []interface{}
	for _, positionItem := range position {
		positionRule = append(positionRule, positionItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "ValRecordAdded", operatorRule, vaultRule, positionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareValRecordAdded)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValRecordAdded", log); err != nil {
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

// ParseValRecordAdded is a log parse operation binding the contract event 0x292d0cd97c126d79e15cbccfedd21071c6230c185049f19a1c84f92a85ec0413.
//
// Solidity: event ValRecordAdded(bytes blsPubkey, address indexed operator, address indexed vault, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseValRecordAdded(log types.Log) (*MevcommitmiddlewareValRecordAdded, error) {
	event := new(MevcommitmiddlewareValRecordAdded)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValRecordAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareValRecordDeletedIterator is returned from FilterValRecordDeleted and is used to iterate over the raw logs and unpacked data for ValRecordDeleted events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValRecordDeletedIterator struct {
	Event *MevcommitmiddlewareValRecordDeleted // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareValRecordDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareValRecordDeleted)
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
		it.Event = new(MevcommitmiddlewareValRecordDeleted)
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
func (it *MevcommitmiddlewareValRecordDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareValRecordDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareValRecordDeleted represents a ValRecordDeleted event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValRecordDeleted struct {
	BlsPubkey []byte
	MsgSender common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValRecordDeleted is a free log retrieval operation binding the contract event 0xf6c1aacd8c44b2fd50a983cd37e4c74ffd298976b52e9ac8fff914c9abb96127.
//
// Solidity: event ValRecordDeleted(bytes blsPubkey, address indexed msgSender)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterValRecordDeleted(opts *bind.FilterOpts, msgSender []common.Address) (*MevcommitmiddlewareValRecordDeletedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "ValRecordDeleted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareValRecordDeletedIterator{contract: _Mevcommitmiddleware.contract, event: "ValRecordDeleted", logs: logs, sub: sub}, nil
}

// WatchValRecordDeleted is a free log subscription operation binding the contract event 0xf6c1aacd8c44b2fd50a983cd37e4c74ffd298976b52e9ac8fff914c9abb96127.
//
// Solidity: event ValRecordDeleted(bytes blsPubkey, address indexed msgSender)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchValRecordDeleted(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareValRecordDeleted, msgSender []common.Address) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "ValRecordDeleted", msgSenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareValRecordDeleted)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValRecordDeleted", log); err != nil {
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

// ParseValRecordDeleted is a log parse operation binding the contract event 0xf6c1aacd8c44b2fd50a983cd37e4c74ffd298976b52e9ac8fff914c9abb96127.
//
// Solidity: event ValRecordDeleted(bytes blsPubkey, address indexed msgSender)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseValRecordDeleted(log types.Log) (*MevcommitmiddlewareValRecordDeleted, error) {
	event := new(MevcommitmiddlewareValRecordDeleted)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValRecordDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareValidatorDeregistrationRequestedIterator is returned from FilterValidatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for ValidatorDeregistrationRequested events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorDeregistrationRequestedIterator struct {
	Event *MevcommitmiddlewareValidatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareValidatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareValidatorDeregistrationRequested)
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
		it.Event = new(MevcommitmiddlewareValidatorDeregistrationRequested)
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
func (it *MevcommitmiddlewareValidatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareValidatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareValidatorDeregistrationRequested represents a ValidatorDeregistrationRequested event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorDeregistrationRequested struct {
	BlsPubkey []byte
	MsgSender common.Address
	Position  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistrationRequested is a free log retrieval operation binding the contract event 0xe07b973733e5b57105d104ee2b05426819359b0c2f731930917b1ac3d763aad0.
//
// Solidity: event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterValidatorDeregistrationRequested(opts *bind.FilterOpts, msgSender []common.Address, position []*big.Int) (*MevcommitmiddlewareValidatorDeregistrationRequestedIterator, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var positionRule []interface{}
	for _, positionItem := range position {
		positionRule = append(positionRule, positionItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "ValidatorDeregistrationRequested", msgSenderRule, positionRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareValidatorDeregistrationRequestedIterator{contract: _Mevcommitmiddleware.contract, event: "ValidatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistrationRequested is a free log subscription operation binding the contract event 0xe07b973733e5b57105d104ee2b05426819359b0c2f731930917b1ac3d763aad0.
//
// Solidity: event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchValidatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareValidatorDeregistrationRequested, msgSender []common.Address, position []*big.Int) (event.Subscription, error) {

	var msgSenderRule []interface{}
	for _, msgSenderItem := range msgSender {
		msgSenderRule = append(msgSenderRule, msgSenderItem)
	}
	var positionRule []interface{}
	for _, positionItem := range position {
		positionRule = append(positionRule, positionItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "ValidatorDeregistrationRequested", msgSenderRule, positionRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareValidatorDeregistrationRequested)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
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

// ParseValidatorDeregistrationRequested is a log parse operation binding the contract event 0xe07b973733e5b57105d104ee2b05426819359b0c2f731930917b1ac3d763aad0.
//
// Solidity: event ValidatorDeregistrationRequested(bytes blsPubkey, address indexed msgSender, uint256 indexed position)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseValidatorDeregistrationRequested(log types.Log) (*MevcommitmiddlewareValidatorDeregistrationRequested, error) {
	event := new(MevcommitmiddlewareValidatorDeregistrationRequested)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareValidatorPositionsSwappedIterator is returned from FilterValidatorPositionsSwapped and is used to iterate over the raw logs and unpacked data for ValidatorPositionsSwapped events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorPositionsSwappedIterator struct {
	Event *MevcommitmiddlewareValidatorPositionsSwapped // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareValidatorPositionsSwappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareValidatorPositionsSwapped)
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
		it.Event = new(MevcommitmiddlewareValidatorPositionsSwapped)
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
func (it *MevcommitmiddlewareValidatorPositionsSwappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareValidatorPositionsSwappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareValidatorPositionsSwapped represents a ValidatorPositionsSwapped event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorPositionsSwapped struct {
	BlsPubkeys   [][]byte
	Vaults       []common.Address
	Operators    []common.Address
	NewPositions []*big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterValidatorPositionsSwapped is a free log retrieval operation binding the contract event 0xe0cafa4f48344c7bc84dcc427b7f80feae065a8ef4dd485cbd07f5e1f954cd97.
//
// Solidity: event ValidatorPositionsSwapped(bytes[] blsPubkeys, address[] vaults, address[] operators, uint256[] newPositions)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterValidatorPositionsSwapped(opts *bind.FilterOpts) (*MevcommitmiddlewareValidatorPositionsSwappedIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "ValidatorPositionsSwapped")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareValidatorPositionsSwappedIterator{contract: _Mevcommitmiddleware.contract, event: "ValidatorPositionsSwapped", logs: logs, sub: sub}, nil
}

// WatchValidatorPositionsSwapped is a free log subscription operation binding the contract event 0xe0cafa4f48344c7bc84dcc427b7f80feae065a8ef4dd485cbd07f5e1f954cd97.
//
// Solidity: event ValidatorPositionsSwapped(bytes[] blsPubkeys, address[] vaults, address[] operators, uint256[] newPositions)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchValidatorPositionsSwapped(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareValidatorPositionsSwapped) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "ValidatorPositionsSwapped")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareValidatorPositionsSwapped)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorPositionsSwapped", log); err != nil {
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

// ParseValidatorPositionsSwapped is a log parse operation binding the contract event 0xe0cafa4f48344c7bc84dcc427b7f80feae065a8ef4dd485cbd07f5e1f954cd97.
//
// Solidity: event ValidatorPositionsSwapped(bytes[] blsPubkeys, address[] vaults, address[] operators, uint256[] newPositions)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseValidatorPositionsSwapped(log types.Log) (*MevcommitmiddlewareValidatorPositionsSwapped, error) {
	event := new(MevcommitmiddlewareValidatorPositionsSwapped)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorPositionsSwapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareValidatorSlashedIterator is returned from FilterValidatorSlashed and is used to iterate over the raw logs and unpacked data for ValidatorSlashed events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorSlashedIterator struct {
	Event *MevcommitmiddlewareValidatorSlashed // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareValidatorSlashedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareValidatorSlashed)
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
		it.Event = new(MevcommitmiddlewareValidatorSlashed)
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
func (it *MevcommitmiddlewareValidatorSlashedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareValidatorSlashedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareValidatorSlashed represents a ValidatorSlashed event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareValidatorSlashed struct {
	BlsPubkey     []byte
	Operator      common.Address
	Vault         common.Address
	SlashedAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterValidatorSlashed is a free log retrieval operation binding the contract event 0x944e6261ff5ee2ea0ba8b612f46a1681cf88c282f88b6a9b95cb5ca0e2369bfc.
//
// Solidity: event ValidatorSlashed(bytes blsPubkey, address indexed operator, address indexed vault, uint256 slashedAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterValidatorSlashed(opts *bind.FilterOpts, operator []common.Address, vault []common.Address) (*MevcommitmiddlewareValidatorSlashedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "ValidatorSlashed", operatorRule, vaultRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareValidatorSlashedIterator{contract: _Mevcommitmiddleware.contract, event: "ValidatorSlashed", logs: logs, sub: sub}, nil
}

// WatchValidatorSlashed is a free log subscription operation binding the contract event 0x944e6261ff5ee2ea0ba8b612f46a1681cf88c282f88b6a9b95cb5ca0e2369bfc.
//
// Solidity: event ValidatorSlashed(bytes blsPubkey, address indexed operator, address indexed vault, uint256 slashedAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchValidatorSlashed(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareValidatorSlashed, operator []common.Address, vault []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "ValidatorSlashed", operatorRule, vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareValidatorSlashed)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorSlashed", log); err != nil {
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

// ParseValidatorSlashed is a log parse operation binding the contract event 0x944e6261ff5ee2ea0ba8b612f46a1681cf88c282f88b6a9b95cb5ca0e2369bfc.
//
// Solidity: event ValidatorSlashed(bytes blsPubkey, address indexed operator, address indexed vault, uint256 slashedAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseValidatorSlashed(log types.Log) (*MevcommitmiddlewareValidatorSlashed, error) {
	event := new(MevcommitmiddlewareValidatorSlashed)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "ValidatorSlashed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareVaultDeregisteredIterator is returned from FilterVaultDeregistered and is used to iterate over the raw logs and unpacked data for VaultDeregistered events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultDeregisteredIterator struct {
	Event *MevcommitmiddlewareVaultDeregistered // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareVaultDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareVaultDeregistered)
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
		it.Event = new(MevcommitmiddlewareVaultDeregistered)
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
func (it *MevcommitmiddlewareVaultDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareVaultDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareVaultDeregistered represents a VaultDeregistered event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultDeregistered struct {
	Vault common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVaultDeregistered is a free log retrieval operation binding the contract event 0xf8769b01493238f5c26a42a7b690cb1ff2b53a7d89d9a57e6332458703db8b04.
//
// Solidity: event VaultDeregistered(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterVaultDeregistered(opts *bind.FilterOpts, vault []common.Address) (*MevcommitmiddlewareVaultDeregisteredIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "VaultDeregistered", vaultRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareVaultDeregisteredIterator{contract: _Mevcommitmiddleware.contract, event: "VaultDeregistered", logs: logs, sub: sub}, nil
}

// WatchVaultDeregistered is a free log subscription operation binding the contract event 0xf8769b01493238f5c26a42a7b690cb1ff2b53a7d89d9a57e6332458703db8b04.
//
// Solidity: event VaultDeregistered(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchVaultDeregistered(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareVaultDeregistered, vault []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "VaultDeregistered", vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareVaultDeregistered)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultDeregistered", log); err != nil {
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

// ParseVaultDeregistered is a log parse operation binding the contract event 0xf8769b01493238f5c26a42a7b690cb1ff2b53a7d89d9a57e6332458703db8b04.
//
// Solidity: event VaultDeregistered(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseVaultDeregistered(log types.Log) (*MevcommitmiddlewareVaultDeregistered, error) {
	event := new(MevcommitmiddlewareVaultDeregistered)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareVaultDeregistrationRequestedIterator is returned from FilterVaultDeregistrationRequested and is used to iterate over the raw logs and unpacked data for VaultDeregistrationRequested events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultDeregistrationRequestedIterator struct {
	Event *MevcommitmiddlewareVaultDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareVaultDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareVaultDeregistrationRequested)
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
		it.Event = new(MevcommitmiddlewareVaultDeregistrationRequested)
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
func (it *MevcommitmiddlewareVaultDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareVaultDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareVaultDeregistrationRequested represents a VaultDeregistrationRequested event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultDeregistrationRequested struct {
	Vault common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterVaultDeregistrationRequested is a free log retrieval operation binding the contract event 0x5d945f2e26a066596d1298a6cc395038a2c99f1de44bc24f9b15b2fbc5c5be0f.
//
// Solidity: event VaultDeregistrationRequested(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterVaultDeregistrationRequested(opts *bind.FilterOpts, vault []common.Address) (*MevcommitmiddlewareVaultDeregistrationRequestedIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "VaultDeregistrationRequested", vaultRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareVaultDeregistrationRequestedIterator{contract: _Mevcommitmiddleware.contract, event: "VaultDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchVaultDeregistrationRequested is a free log subscription operation binding the contract event 0x5d945f2e26a066596d1298a6cc395038a2c99f1de44bc24f9b15b2fbc5c5be0f.
//
// Solidity: event VaultDeregistrationRequested(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchVaultDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareVaultDeregistrationRequested, vault []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "VaultDeregistrationRequested", vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareVaultDeregistrationRequested)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultDeregistrationRequested", log); err != nil {
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

// ParseVaultDeregistrationRequested is a log parse operation binding the contract event 0x5d945f2e26a066596d1298a6cc395038a2c99f1de44bc24f9b15b2fbc5c5be0f.
//
// Solidity: event VaultDeregistrationRequested(address indexed vault)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseVaultDeregistrationRequested(log types.Log) (*MevcommitmiddlewareVaultDeregistrationRequested, error) {
	event := new(MevcommitmiddlewareVaultDeregistrationRequested)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareVaultFactorySetIterator is returned from FilterVaultFactorySet and is used to iterate over the raw logs and unpacked data for VaultFactorySet events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultFactorySetIterator struct {
	Event *MevcommitmiddlewareVaultFactorySet // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareVaultFactorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareVaultFactorySet)
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
		it.Event = new(MevcommitmiddlewareVaultFactorySet)
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
func (it *MevcommitmiddlewareVaultFactorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareVaultFactorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareVaultFactorySet represents a VaultFactorySet event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultFactorySet struct {
	VaultFactory common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterVaultFactorySet is a free log retrieval operation binding the contract event 0x57a7b93fb32405c4ce3505a4c1fbb4b06d58e4c10ebd4cfad8b948f2e5d5ffdf.
//
// Solidity: event VaultFactorySet(address vaultFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterVaultFactorySet(opts *bind.FilterOpts) (*MevcommitmiddlewareVaultFactorySetIterator, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "VaultFactorySet")
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareVaultFactorySetIterator{contract: _Mevcommitmiddleware.contract, event: "VaultFactorySet", logs: logs, sub: sub}, nil
}

// WatchVaultFactorySet is a free log subscription operation binding the contract event 0x57a7b93fb32405c4ce3505a4c1fbb4b06d58e4c10ebd4cfad8b948f2e5d5ffdf.
//
// Solidity: event VaultFactorySet(address vaultFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchVaultFactorySet(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareVaultFactorySet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "VaultFactorySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareVaultFactorySet)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultFactorySet", log); err != nil {
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

// ParseVaultFactorySet is a log parse operation binding the contract event 0x57a7b93fb32405c4ce3505a4c1fbb4b06d58e4c10ebd4cfad8b948f2e5d5ffdf.
//
// Solidity: event VaultFactorySet(address vaultFactory)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseVaultFactorySet(log types.Log) (*MevcommitmiddlewareVaultFactorySet, error) {
	event := new(MevcommitmiddlewareVaultFactorySet)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultFactorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareVaultRegisteredIterator is returned from FilterVaultRegistered and is used to iterate over the raw logs and unpacked data for VaultRegistered events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultRegisteredIterator struct {
	Event *MevcommitmiddlewareVaultRegistered // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareVaultRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareVaultRegistered)
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
		it.Event = new(MevcommitmiddlewareVaultRegistered)
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
func (it *MevcommitmiddlewareVaultRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareVaultRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareVaultRegistered represents a VaultRegistered event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultRegistered struct {
	Vault       common.Address
	SlashAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVaultRegistered is a free log retrieval operation binding the contract event 0x9166afffc11e38358e14e3e0563931dc8f4fd340152eaef1ea729a70093279a4.
//
// Solidity: event VaultRegistered(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterVaultRegistered(opts *bind.FilterOpts, vault []common.Address) (*MevcommitmiddlewareVaultRegisteredIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "VaultRegistered", vaultRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareVaultRegisteredIterator{contract: _Mevcommitmiddleware.contract, event: "VaultRegistered", logs: logs, sub: sub}, nil
}

// WatchVaultRegistered is a free log subscription operation binding the contract event 0x9166afffc11e38358e14e3e0563931dc8f4fd340152eaef1ea729a70093279a4.
//
// Solidity: event VaultRegistered(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchVaultRegistered(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareVaultRegistered, vault []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "VaultRegistered", vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareVaultRegistered)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultRegistered", log); err != nil {
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

// ParseVaultRegistered is a log parse operation binding the contract event 0x9166afffc11e38358e14e3e0563931dc8f4fd340152eaef1ea729a70093279a4.
//
// Solidity: event VaultRegistered(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseVaultRegistered(log types.Log) (*MevcommitmiddlewareVaultRegistered, error) {
	event := new(MevcommitmiddlewareVaultRegistered)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitmiddlewareVaultSlashAmountUpdatedIterator is returned from FilterVaultSlashAmountUpdated and is used to iterate over the raw logs and unpacked data for VaultSlashAmountUpdated events raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultSlashAmountUpdatedIterator struct {
	Event *MevcommitmiddlewareVaultSlashAmountUpdated // Event containing the contract specifics and raw log

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
func (it *MevcommitmiddlewareVaultSlashAmountUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitmiddlewareVaultSlashAmountUpdated)
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
		it.Event = new(MevcommitmiddlewareVaultSlashAmountUpdated)
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
func (it *MevcommitmiddlewareVaultSlashAmountUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitmiddlewareVaultSlashAmountUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitmiddlewareVaultSlashAmountUpdated represents a VaultSlashAmountUpdated event raised by the Mevcommitmiddleware contract.
type MevcommitmiddlewareVaultSlashAmountUpdated struct {
	Vault       common.Address
	SlashAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVaultSlashAmountUpdated is a free log retrieval operation binding the contract event 0x373931bc95a57530dae76c99d7fa3efc056899775f8cfa82867d5c2a92b06c10.
//
// Solidity: event VaultSlashAmountUpdated(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) FilterVaultSlashAmountUpdated(opts *bind.FilterOpts, vault []common.Address) (*MevcommitmiddlewareVaultSlashAmountUpdatedIterator, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.FilterLogs(opts, "VaultSlashAmountUpdated", vaultRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitmiddlewareVaultSlashAmountUpdatedIterator{contract: _Mevcommitmiddleware.contract, event: "VaultSlashAmountUpdated", logs: logs, sub: sub}, nil
}

// WatchVaultSlashAmountUpdated is a free log subscription operation binding the contract event 0x373931bc95a57530dae76c99d7fa3efc056899775f8cfa82867d5c2a92b06c10.
//
// Solidity: event VaultSlashAmountUpdated(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) WatchVaultSlashAmountUpdated(opts *bind.WatchOpts, sink chan<- *MevcommitmiddlewareVaultSlashAmountUpdated, vault []common.Address) (event.Subscription, error) {

	var vaultRule []interface{}
	for _, vaultItem := range vault {
		vaultRule = append(vaultRule, vaultItem)
	}

	logs, sub, err := _Mevcommitmiddleware.contract.WatchLogs(opts, "VaultSlashAmountUpdated", vaultRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitmiddlewareVaultSlashAmountUpdated)
				if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultSlashAmountUpdated", log); err != nil {
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

// ParseVaultSlashAmountUpdated is a log parse operation binding the contract event 0x373931bc95a57530dae76c99d7fa3efc056899775f8cfa82867d5c2a92b06c10.
//
// Solidity: event VaultSlashAmountUpdated(address indexed vault, uint160 slashAmount)
func (_Mevcommitmiddleware *MevcommitmiddlewareFilterer) ParseVaultSlashAmountUpdated(log types.Log) (*MevcommitmiddlewareVaultSlashAmountUpdated, error) {
	event := new(MevcommitmiddlewareVaultSlashAmountUpdated)
	if err := _Mevcommitmiddleware.contract.UnpackLog(event, "VaultSlashAmountUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
