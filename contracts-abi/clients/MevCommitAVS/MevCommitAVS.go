// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mevcommitavs

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

// IMevCommitAVSLSTRestakerRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSLSTRestakerRegistrationInfo struct {
	Exists                 bool
	ChosenValidators       [][]byte
	NumChosen              *big.Int
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}

// IMevCommitAVSOperatorRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSOperatorRegistrationInfo struct {
	Exists                 bool
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}

// IMevCommitAVSValidatorRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSValidatorRegistrationInfo struct {
	Exists                 bool
	PodOwner               common.Address
	FreezeOccurrence       BlockHeightOccurrenceOccurrence
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// MevcommitavsMetaData contains all meta data concerning the Mevcommitavs contract.
var MevcommitavsMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterLSTRestaker\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterValidators\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLSTRestakerRegInfo\",\"inputs\":[{\"name\":\"lstRestaker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVS.LSTRestakerRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"chosenValidators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRegInfo\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVS.OperatorRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatorRegInfo\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVS.ValidatorRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freezeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationManager_\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"},{\"name\":\"eigenPodManager_\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"},{\"name\":\"strategyManager_\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"avsDirectory_\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"restakeableStrategies_\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"freezeOracle_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezeFee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unfreezeReceiver_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezePeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"operatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lstRestakerDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lstRestakerDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lstRestakerRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerLSTRestaker\",\"inputs\":[{\"name\":\"chosenValidators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidatorsByPodOwners\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"podOwners\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestLSTRestakerDeregistration\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestOperatorDeregistration\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestValidatorsDeregistration\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"restakeableStrategies\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAVSDirectory\",\"inputs\":[{\"name\":\"avsDirectory_\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegationManager\",\"inputs\":[{\"name\":\"delegationManager_\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEigenPodManager\",\"inputs\":[{\"name\":\"eigenPodManager_\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFreezeOracle\",\"inputs\":[{\"name\":\"freezeOracle_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLstRestakerDeregPeriodBlocks\",\"inputs\":[{\"name\":\"lstRestakerDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorDeregPeriodBlocks\",\"inputs\":[{\"name\":\"operatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRestakeableStrategies\",\"inputs\":[{\"name\":\"restakeableStrategies_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStrategyManager\",\"inputs\":[{\"name\":\"strategyManager_\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeFee\",\"inputs\":[{\"name\":\"unfreezeFee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezePeriodBlocks\",\"inputs\":[{\"name\":\"unfreezePeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeReceiver\",\"inputs\":[{\"name\":\"unfreezeReceiver_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setValidatorDeregPeriodBlocks\",\"inputs\":[{\"name\":\"validatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unfreezeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezeReceiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadataURI\",\"inputs\":[{\"name\":\"metadataURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validatorDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freezeOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"deregRequestOccurrence\",\"type\":\"tuple\",\"internalType\":\"structBlockHeightOccurrence.Occurrence\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AVSDirectorySet\",\"inputs\":[{\"name\":\"avsDirectory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegationManagerSet\",\"inputs\":[{\"name\":\"delegationManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EigenPodManagerSet\",\"inputs\":[{\"name\":\"eigenPodManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FreezeOracleSet\",\"inputs\":[{\"name\":\"freezeOracle\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"lstRestakerDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregistered\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregistrationRequested\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerRegistered\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"operatorDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistrationRequested\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RestakeableStrategiesSet\",\"inputs\":[{\"name\":\"restakeableStrategies\",\"type\":\"address[]\",\"indexed\":true,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyManagerSet\",\"inputs\":[{\"name\":\"strategyManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezeFeeSet\",\"inputs\":[{\"name\":\"unfreezeFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezePeriodBlocksSet\",\"inputs\":[{\"name\":\"unfreezePeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezeReceiverSet\",\"inputs\":[{\"name\":\"unfreezeReceiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"validatorDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistrationRequested\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorFrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRegistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnfrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"DeregistrationAlreadyRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeregistrationNotRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeregistrationTooSoon\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FrozenValidatorCannotDeregister\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LstRestakerIsRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LstRestakerNotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NeedChosenValidators\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoDelegationToRegisteredOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoEigenStrategyDeposits\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoPodExists\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperatorDeregAlreadyRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OperatorNotRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"RefundFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotEigenCoreOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotFreezeOracle\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderIsNotSpecifiedOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsRegisteredOperator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SenderNotPodOwnerOrOperator\",\"inputs\":[{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderNotPodOwnerOrOperatorOfValidator\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnfreezeFeeRequired\",\"inputs\":[{\"name\":\"requiredFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnfreezeTooSoon\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnfreezeTransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorAlreadyFrozen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorDeregAlreadyRequested\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ValidatorIsRegistered\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotActiveWithEigenCore\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotFrozen\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"type\":\"error\",\"name\":\"ValidatorNotRegistered\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]",
}

// MevcommitavsABI is the input ABI used to generate the binding from.
// Deprecated: Use MevcommitavsMetaData.ABI instead.
var MevcommitavsABI = MevcommitavsMetaData.ABI

// Mevcommitavs is an auto generated Go binding around an Ethereum contract.
type Mevcommitavs struct {
	MevcommitavsCaller     // Read-only binding to the contract
	MevcommitavsTransactor // Write-only binding to the contract
	MevcommitavsFilterer   // Log filterer for contract events
}

// MevcommitavsCaller is an auto generated read-only Go binding around an Ethereum contract.
type MevcommitavsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitavsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MevcommitavsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitavsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MevcommitavsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevcommitavsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MevcommitavsSession struct {
	Contract     *Mevcommitavs     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MevcommitavsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MevcommitavsCallerSession struct {
	Contract *MevcommitavsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// MevcommitavsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MevcommitavsTransactorSession struct {
	Contract     *MevcommitavsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// MevcommitavsRaw is an auto generated low-level Go binding around an Ethereum contract.
type MevcommitavsRaw struct {
	Contract *Mevcommitavs // Generic contract binding to access the raw methods on
}

// MevcommitavsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MevcommitavsCallerRaw struct {
	Contract *MevcommitavsCaller // Generic read-only contract binding to access the raw methods on
}

// MevcommitavsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MevcommitavsTransactorRaw struct {
	Contract *MevcommitavsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMevcommitavs creates a new instance of Mevcommitavs, bound to a specific deployed contract.
func NewMevcommitavs(address common.Address, backend bind.ContractBackend) (*Mevcommitavs, error) {
	contract, err := bindMevcommitavs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavs{MevcommitavsCaller: MevcommitavsCaller{contract: contract}, MevcommitavsTransactor: MevcommitavsTransactor{contract: contract}, MevcommitavsFilterer: MevcommitavsFilterer{contract: contract}}, nil
}

// NewMevcommitavsCaller creates a new read-only instance of Mevcommitavs, bound to a specific deployed contract.
func NewMevcommitavsCaller(address common.Address, caller bind.ContractCaller) (*MevcommitavsCaller, error) {
	contract, err := bindMevcommitavs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsCaller{contract: contract}, nil
}

// NewMevcommitavsTransactor creates a new write-only instance of Mevcommitavs, bound to a specific deployed contract.
func NewMevcommitavsTransactor(address common.Address, transactor bind.ContractTransactor) (*MevcommitavsTransactor, error) {
	contract, err := bindMevcommitavs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsTransactor{contract: contract}, nil
}

// NewMevcommitavsFilterer creates a new log filterer instance of Mevcommitavs, bound to a specific deployed contract.
func NewMevcommitavsFilterer(address common.Address, filterer bind.ContractFilterer) (*MevcommitavsFilterer, error) {
	contract, err := bindMevcommitavs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsFilterer{contract: contract}, nil
}

// bindMevcommitavs binds a generic wrapper to an already deployed contract.
func bindMevcommitavs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MevcommitavsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitavs *MevcommitavsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitavs.Contract.MevcommitavsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitavs *MevcommitavsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.MevcommitavsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitavs *MevcommitavsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.MevcommitavsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitavs *MevcommitavsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitavs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitavs *MevcommitavsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitavs *MevcommitavsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavs *MevcommitavsCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavs *MevcommitavsSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitavs.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitavs.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavs *MevcommitavsCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitavs.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitavs.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavs *MevcommitavsSession) AvsDirectory() (common.Address, error) {
	return _Mevcommitavs.Contract.AvsDirectory(&_Mevcommitavs.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) AvsDirectory() (common.Address, error) {
	return _Mevcommitavs.Contract.AvsDirectory(&_Mevcommitavs.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) FreezeOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "freezeOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavs *MevcommitavsSession) FreezeOracle() (common.Address, error) {
	return _Mevcommitavs.Contract.FreezeOracle(&_Mevcommitavs.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) FreezeOracle() (common.Address, error) {
	return _Mevcommitavs.Contract.FreezeOracle(&_Mevcommitavs.CallOpts)
}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCaller) GetLSTRestakerRegInfo(opts *bind.CallOpts, lstRestaker common.Address) (IMevCommitAVSLSTRestakerRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "getLSTRestakerRegInfo", lstRestaker)

	if err != nil {
		return *new(IMevCommitAVSLSTRestakerRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSLSTRestakerRegistrationInfo)).(*IMevCommitAVSLSTRestakerRegistrationInfo)

	return out0, err

}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsSession) GetLSTRestakerRegInfo(lstRestaker common.Address) (IMevCommitAVSLSTRestakerRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetLSTRestakerRegInfo(&_Mevcommitavs.CallOpts, lstRestaker)
}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCallerSession) GetLSTRestakerRegInfo(lstRestaker common.Address) (IMevCommitAVSLSTRestakerRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetLSTRestakerRegInfo(&_Mevcommitavs.CallOpts, lstRestaker)
}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCaller) GetOperatorRegInfo(opts *bind.CallOpts, operator common.Address) (IMevCommitAVSOperatorRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "getOperatorRegInfo", operator)

	if err != nil {
		return *new(IMevCommitAVSOperatorRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSOperatorRegistrationInfo)).(*IMevCommitAVSOperatorRegistrationInfo)

	return out0, err

}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsSession) GetOperatorRegInfo(operator common.Address) (IMevCommitAVSOperatorRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetOperatorRegInfo(&_Mevcommitavs.CallOpts, operator)
}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCallerSession) GetOperatorRegInfo(operator common.Address) (IMevCommitAVSOperatorRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetOperatorRegInfo(&_Mevcommitavs.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavs *MevcommitavsCaller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavs *MevcommitavsSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _Mevcommitavs.Contract.GetOperatorRestakedStrategies(&_Mevcommitavs.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavs *MevcommitavsCallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _Mevcommitavs.Contract.GetOperatorRestakedStrategies(&_Mevcommitavs.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavs *MevcommitavsCaller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavs *MevcommitavsSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _Mevcommitavs.Contract.GetRestakeableStrategies(&_Mevcommitavs.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavs *MevcommitavsCallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _Mevcommitavs.Contract.GetRestakeableStrategies(&_Mevcommitavs.CallOpts)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCaller) GetValidatorRegInfo(opts *bind.CallOpts, valPubKey []byte) (IMevCommitAVSValidatorRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "getValidatorRegInfo", valPubKey)

	if err != nil {
		return *new(IMevCommitAVSValidatorRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSValidatorRegistrationInfo)).(*IMevCommitAVSValidatorRegistrationInfo)

	return out0, err

}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavs *MevcommitavsSession) GetValidatorRegInfo(valPubKey []byte) (IMevCommitAVSValidatorRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetValidatorRegInfo(&_Mevcommitavs.CallOpts, valPubKey)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavs *MevcommitavsCallerSession) GetValidatorRegInfo(valPubKey []byte) (IMevCommitAVSValidatorRegistrationInfo, error) {
	return _Mevcommitavs.Contract.GetValidatorRegInfo(&_Mevcommitavs.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavs *MevcommitavsCaller) IsValidatorOptedIn(opts *bind.CallOpts, valPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "isValidatorOptedIn", valPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavs *MevcommitavsSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Mevcommitavs.Contract.IsValidatorOptedIn(&_Mevcommitavs.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavs *MevcommitavsCallerSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Mevcommitavs.Contract.IsValidatorOptedIn(&_Mevcommitavs.CallOpts, valPubKey)
}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCaller) LstRestakerDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "lstRestakerDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsSession) LstRestakerDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.LstRestakerDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCallerSession) LstRestakerDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.LstRestakerDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCaller) LstRestakerRegistrations(opts *bind.CallOpts, arg0 common.Address) (struct {
	Exists                 bool
	NumChosen              *big.Int
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "lstRestakerRegistrations", arg0)

	outstruct := new(struct {
		Exists                 bool
		NumChosen              *big.Int
		DeregRequestOccurrence BlockHeightOccurrenceOccurrence
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.NumChosen = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[2], new(BlockHeightOccurrenceOccurrence)).(*BlockHeightOccurrenceOccurrence)

	return *outstruct, err

}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsSession) LstRestakerRegistrations(arg0 common.Address) (struct {
	Exists                 bool
	NumChosen              *big.Int
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.LstRestakerRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCallerSession) LstRestakerRegistrations(arg0 common.Address) (struct {
	Exists                 bool
	NumChosen              *big.Int
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.LstRestakerRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCaller) OperatorDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "operatorDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsSession) OperatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.OperatorDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCallerSession) OperatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.OperatorDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCaller) OperatorRegistrations(opts *bind.CallOpts, arg0 common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "operatorRegistrations", arg0)

	outstruct := new(struct {
		Exists                 bool
		DeregRequestOccurrence BlockHeightOccurrenceOccurrence
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[1], new(BlockHeightOccurrenceOccurrence)).(*BlockHeightOccurrenceOccurrence)

	return *outstruct, err

}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsSession) OperatorRegistrations(arg0 common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.OperatorRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCallerSession) OperatorRegistrations(arg0 common.Address) (struct {
	Exists                 bool
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.OperatorRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavs *MevcommitavsSession) Owner() (common.Address, error) {
	return _Mevcommitavs.Contract.Owner(&_Mevcommitavs.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) Owner() (common.Address, error) {
	return _Mevcommitavs.Contract.Owner(&_Mevcommitavs.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavs *MevcommitavsCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavs *MevcommitavsSession) Paused() (bool, error) {
	return _Mevcommitavs.Contract.Paused(&_Mevcommitavs.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavs *MevcommitavsCallerSession) Paused() (bool, error) {
	return _Mevcommitavs.Contract.Paused(&_Mevcommitavs.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavs *MevcommitavsSession) PendingOwner() (common.Address, error) {
	return _Mevcommitavs.Contract.PendingOwner(&_Mevcommitavs.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) PendingOwner() (common.Address, error) {
	return _Mevcommitavs.Contract.PendingOwner(&_Mevcommitavs.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavs *MevcommitavsCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavs *MevcommitavsSession) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitavs.Contract.ProxiableUUID(&_Mevcommitavs.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavs *MevcommitavsCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitavs.Contract.ProxiableUUID(&_Mevcommitavs.CallOpts)
}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) RestakeableStrategies(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "restakeableStrategies", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavs *MevcommitavsSession) RestakeableStrategies(arg0 *big.Int) (common.Address, error) {
	return _Mevcommitavs.Contract.RestakeableStrategies(&_Mevcommitavs.CallOpts, arg0)
}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) RestakeableStrategies(arg0 *big.Int) (common.Address, error) {
	return _Mevcommitavs.Contract.RestakeableStrategies(&_Mevcommitavs.CallOpts, arg0)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCaller) UnfreezeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "unfreezeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavs *MevcommitavsSession) UnfreezeFee() (*big.Int, error) {
	return _Mevcommitavs.Contract.UnfreezeFee(&_Mevcommitavs.CallOpts)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCallerSession) UnfreezeFee() (*big.Int, error) {
	return _Mevcommitavs.Contract.UnfreezeFee(&_Mevcommitavs.CallOpts)
}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCaller) UnfreezePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "unfreezePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsSession) UnfreezePeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.UnfreezePeriodBlocks(&_Mevcommitavs.CallOpts)
}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCallerSession) UnfreezePeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.UnfreezePeriodBlocks(&_Mevcommitavs.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavs *MevcommitavsCaller) UnfreezeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "unfreezeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavs *MevcommitavsSession) UnfreezeReceiver() (common.Address, error) {
	return _Mevcommitavs.Contract.UnfreezeReceiver(&_Mevcommitavs.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavs *MevcommitavsCallerSession) UnfreezeReceiver() (common.Address, error) {
	return _Mevcommitavs.Contract.UnfreezeReceiver(&_Mevcommitavs.CallOpts)
}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCaller) ValidatorDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "validatorDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsSession) ValidatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.ValidatorDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavs *MevcommitavsCallerSession) ValidatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavs.Contract.ValidatorDeregPeriodBlocks(&_Mevcommitavs.CallOpts)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeOccurrence, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCaller) ValidatorRegistrations(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists                 bool
	PodOwner               common.Address
	FreezeOccurrence       BlockHeightOccurrenceOccurrence
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	var out []interface{}
	err := _Mevcommitavs.contract.Call(opts, &out, "validatorRegistrations", arg0)

	outstruct := new(struct {
		Exists                 bool
		PodOwner               common.Address
		FreezeOccurrence       BlockHeightOccurrenceOccurrence
		DeregRequestOccurrence BlockHeightOccurrenceOccurrence
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PodOwner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.FreezeOccurrence = *abi.ConvertType(out[2], new(BlockHeightOccurrenceOccurrence)).(*BlockHeightOccurrenceOccurrence)
	outstruct.DeregRequestOccurrence = *abi.ConvertType(out[3], new(BlockHeightOccurrenceOccurrence)).(*BlockHeightOccurrenceOccurrence)

	return *outstruct, err

}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeOccurrence, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsSession) ValidatorRegistrations(arg0 []byte) (struct {
	Exists                 bool
	PodOwner               common.Address
	FreezeOccurrence       BlockHeightOccurrenceOccurrence
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.ValidatorRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeOccurrence, (bool,uint256) deregRequestOccurrence)
func (_Mevcommitavs *MevcommitavsCallerSession) ValidatorRegistrations(arg0 []byte) (struct {
	Exists                 bool
	PodOwner               common.Address
	FreezeOccurrence       BlockHeightOccurrenceOccurrence
	DeregRequestOccurrence BlockHeightOccurrenceOccurrence
}, error) {
	return _Mevcommitavs.Contract.ValidatorRegistrations(&_Mevcommitavs.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavs *MevcommitavsTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavs *MevcommitavsSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.AcceptOwnership(&_Mevcommitavs.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.AcceptOwnership(&_Mevcommitavs.TransactOpts)
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavs *MevcommitavsTransactor) DeregisterLSTRestaker(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "deregisterLSTRestaker")
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavs *MevcommitavsSession) DeregisterLSTRestaker() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterLSTRestaker(&_Mevcommitavs.TransactOpts)
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) DeregisterLSTRestaker() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterLSTRestaker(&_Mevcommitavs.TransactOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavs *MevcommitavsTransactor) DeregisterOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "deregisterOperator", operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavs *MevcommitavsSession) DeregisterOperator(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterOperator(&_Mevcommitavs.TransactOpts, operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) DeregisterOperator(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterOperator(&_Mevcommitavs.TransactOpts, operator)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactor) DeregisterValidators(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "deregisterValidators", valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsSession) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterValidators(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.DeregisterValidators(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactor) Freeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "freeze", valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsSession) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Freeze(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Freeze(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) Initialize(opts *bind.TransactOpts, owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "initialize", owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsSession) Initialize(owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Initialize(&_Mevcommitavs.TransactOpts, owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Initialize(owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Initialize(&_Mevcommitavs.TransactOpts, owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavs *MevcommitavsTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavs *MevcommitavsSession) Pause() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Pause(&_Mevcommitavs.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Pause() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Pause(&_Mevcommitavs.TransactOpts)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavs *MevcommitavsTransactor) RegisterLSTRestaker(opts *bind.TransactOpts, chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "registerLSTRestaker", chosenValidators)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavs *MevcommitavsSession) RegisterLSTRestaker(chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterLSTRestaker(&_Mevcommitavs.TransactOpts, chosenValidators)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RegisterLSTRestaker(chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterLSTRestaker(&_Mevcommitavs.TransactOpts, chosenValidators)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavs *MevcommitavsTransactor) RegisterOperator(opts *bind.TransactOpts, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "registerOperator", operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavs *MevcommitavsSession) RegisterOperator(operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterOperator(&_Mevcommitavs.TransactOpts, operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RegisterOperator(operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterOperator(&_Mevcommitavs.TransactOpts, operatorSignature)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavs *MevcommitavsTransactor) RegisterValidatorsByPodOwners(opts *bind.TransactOpts, valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "registerValidatorsByPodOwners", valPubKeys, podOwners)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavs *MevcommitavsSession) RegisterValidatorsByPodOwners(valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterValidatorsByPodOwners(&_Mevcommitavs.TransactOpts, valPubKeys, podOwners)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RegisterValidatorsByPodOwners(valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RegisterValidatorsByPodOwners(&_Mevcommitavs.TransactOpts, valPubKeys, podOwners)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavs *MevcommitavsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavs *MevcommitavsSession) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RenounceOwnership(&_Mevcommitavs.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RenounceOwnership(&_Mevcommitavs.TransactOpts)
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavs *MevcommitavsTransactor) RequestLSTRestakerDeregistration(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "requestLSTRestakerDeregistration")
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavs *MevcommitavsSession) RequestLSTRestakerDeregistration() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestLSTRestakerDeregistration(&_Mevcommitavs.TransactOpts)
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RequestLSTRestakerDeregistration() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestLSTRestakerDeregistration(&_Mevcommitavs.TransactOpts)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavs *MevcommitavsTransactor) RequestOperatorDeregistration(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "requestOperatorDeregistration", operator)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavs *MevcommitavsSession) RequestOperatorDeregistration(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestOperatorDeregistration(&_Mevcommitavs.TransactOpts, operator)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RequestOperatorDeregistration(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestOperatorDeregistration(&_Mevcommitavs.TransactOpts, operator)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactor) RequestValidatorsDeregistration(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "requestValidatorsDeregistration", valPubKeys)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsSession) RequestValidatorsDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestValidatorsDeregistration(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) RequestValidatorsDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.RequestValidatorsDeregistration(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetAVSDirectory(opts *bind.TransactOpts, avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setAVSDirectory", avsDirectory_)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetAVSDirectory(avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetAVSDirectory(&_Mevcommitavs.TransactOpts, avsDirectory_)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetAVSDirectory(avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetAVSDirectory(&_Mevcommitavs.TransactOpts, avsDirectory_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetDelegationManager(opts *bind.TransactOpts, delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setDelegationManager", delegationManager_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetDelegationManager(delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetDelegationManager(&_Mevcommitavs.TransactOpts, delegationManager_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetDelegationManager(delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetDelegationManager(&_Mevcommitavs.TransactOpts, delegationManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetEigenPodManager(opts *bind.TransactOpts, eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setEigenPodManager", eigenPodManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetEigenPodManager(eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetEigenPodManager(&_Mevcommitavs.TransactOpts, eigenPodManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetEigenPodManager(eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetEigenPodManager(&_Mevcommitavs.TransactOpts, eigenPodManager_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetFreezeOracle(opts *bind.TransactOpts, freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setFreezeOracle", freezeOracle_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetFreezeOracle(freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetFreezeOracle(&_Mevcommitavs.TransactOpts, freezeOracle_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetFreezeOracle(freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetFreezeOracle(&_Mevcommitavs.TransactOpts, freezeOracle_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetLstRestakerDeregPeriodBlocks(opts *bind.TransactOpts, lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setLstRestakerDeregPeriodBlocks", lstRestakerDeregPeriodBlocks_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetLstRestakerDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, lstRestakerDeregPeriodBlocks_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetLstRestakerDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, lstRestakerDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetOperatorDeregPeriodBlocks(opts *bind.TransactOpts, operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setOperatorDeregPeriodBlocks", operatorDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetOperatorDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, operatorDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetOperatorDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, operatorDeregPeriodBlocks_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetRestakeableStrategies(opts *bind.TransactOpts, restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setRestakeableStrategies", restakeableStrategies_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetRestakeableStrategies(restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetRestakeableStrategies(&_Mevcommitavs.TransactOpts, restakeableStrategies_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetRestakeableStrategies(restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetRestakeableStrategies(&_Mevcommitavs.TransactOpts, restakeableStrategies_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetStrategyManager(opts *bind.TransactOpts, strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setStrategyManager", strategyManager_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetStrategyManager(strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetStrategyManager(&_Mevcommitavs.TransactOpts, strategyManager_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetStrategyManager(strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetStrategyManager(&_Mevcommitavs.TransactOpts, strategyManager_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetUnfreezeFee(opts *bind.TransactOpts, unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setUnfreezeFee", unfreezeFee_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetUnfreezeFee(unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezeFee(&_Mevcommitavs.TransactOpts, unfreezeFee_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetUnfreezeFee(unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezeFee(&_Mevcommitavs.TransactOpts, unfreezeFee_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetUnfreezePeriodBlocks(opts *bind.TransactOpts, unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setUnfreezePeriodBlocks", unfreezePeriodBlocks_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetUnfreezePeriodBlocks(unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezePeriodBlocks(&_Mevcommitavs.TransactOpts, unfreezePeriodBlocks_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetUnfreezePeriodBlocks(unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezePeriodBlocks(&_Mevcommitavs.TransactOpts, unfreezePeriodBlocks_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetUnfreezeReceiver(opts *bind.TransactOpts, unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setUnfreezeReceiver", unfreezeReceiver_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetUnfreezeReceiver(unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezeReceiver(&_Mevcommitavs.TransactOpts, unfreezeReceiver_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetUnfreezeReceiver(unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetUnfreezeReceiver(&_Mevcommitavs.TransactOpts, unfreezeReceiver_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) SetValidatorDeregPeriodBlocks(opts *bind.TransactOpts, validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "setValidatorDeregPeriodBlocks", validatorDeregPeriodBlocks_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsSession) SetValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetValidatorDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, validatorDeregPeriodBlocks_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) SetValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.SetValidatorDeregPeriodBlocks(&_Mevcommitavs.TransactOpts, validatorDeregPeriodBlocks_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavs *MevcommitavsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavs *MevcommitavsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.TransferOwnership(&_Mevcommitavs.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.TransferOwnership(&_Mevcommitavs.TransactOpts, newOwner)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Mevcommitavs *MevcommitavsTransactor) Unfreeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "unfreeze", valPubKeys)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Mevcommitavs *MevcommitavsSession) Unfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Unfreeze(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKeys) payable returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Unfreeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Unfreeze(&_Mevcommitavs.TransactOpts, valPubKeys)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavs *MevcommitavsTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavs *MevcommitavsSession) Unpause() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Unpause(&_Mevcommitavs.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Unpause() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Unpause(&_Mevcommitavs.TransactOpts)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsTransactor) UpdateMetadataURI(opts *bind.TransactOpts, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "updateMetadataURI", metadataURI_)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsSession) UpdateMetadataURI(metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.UpdateMetadataURI(&_Mevcommitavs.TransactOpts, metadataURI_)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) UpdateMetadataURI(metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.UpdateMetadataURI(&_Mevcommitavs.TransactOpts, metadataURI_)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavs *MevcommitavsTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavs *MevcommitavsSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.UpgradeToAndCall(&_Mevcommitavs.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.UpgradeToAndCall(&_Mevcommitavs.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavs *MevcommitavsTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavs.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavs *MevcommitavsSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Fallback(&_Mevcommitavs.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Fallback(&_Mevcommitavs.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavs *MevcommitavsTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavs.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavs *MevcommitavsSession) Receive() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Receive(&_Mevcommitavs.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavs *MevcommitavsTransactorSession) Receive() (*types.Transaction, error) {
	return _Mevcommitavs.Contract.Receive(&_Mevcommitavs.TransactOpts)
}

// MevcommitavsAVSDirectorySetIterator is returned from FilterAVSDirectorySet and is used to iterate over the raw logs and unpacked data for AVSDirectorySet events raised by the Mevcommitavs contract.
type MevcommitavsAVSDirectorySetIterator struct {
	Event *MevcommitavsAVSDirectorySet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsAVSDirectorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsAVSDirectorySet)
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
		it.Event = new(MevcommitavsAVSDirectorySet)
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
func (it *MevcommitavsAVSDirectorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsAVSDirectorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsAVSDirectorySet represents a AVSDirectorySet event raised by the Mevcommitavs contract.
type MevcommitavsAVSDirectorySet struct {
	AvsDirectory common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAVSDirectorySet is a free log retrieval operation binding the contract event 0x934223b20c24d569ff89796ae10a6997d43e2b3df0c3677fb6ca1f6e37ce344b.
//
// Solidity: event AVSDirectorySet(address indexed avsDirectory)
func (_Mevcommitavs *MevcommitavsFilterer) FilterAVSDirectorySet(opts *bind.FilterOpts, avsDirectory []common.Address) (*MevcommitavsAVSDirectorySetIterator, error) {

	var avsDirectoryRule []interface{}
	for _, avsDirectoryItem := range avsDirectory {
		avsDirectoryRule = append(avsDirectoryRule, avsDirectoryItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "AVSDirectorySet", avsDirectoryRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsAVSDirectorySetIterator{contract: _Mevcommitavs.contract, event: "AVSDirectorySet", logs: logs, sub: sub}, nil
}

// WatchAVSDirectorySet is a free log subscription operation binding the contract event 0x934223b20c24d569ff89796ae10a6997d43e2b3df0c3677fb6ca1f6e37ce344b.
//
// Solidity: event AVSDirectorySet(address indexed avsDirectory)
func (_Mevcommitavs *MevcommitavsFilterer) WatchAVSDirectorySet(opts *bind.WatchOpts, sink chan<- *MevcommitavsAVSDirectorySet, avsDirectory []common.Address) (event.Subscription, error) {

	var avsDirectoryRule []interface{}
	for _, avsDirectoryItem := range avsDirectory {
		avsDirectoryRule = append(avsDirectoryRule, avsDirectoryItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "AVSDirectorySet", avsDirectoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsAVSDirectorySet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "AVSDirectorySet", log); err != nil {
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

// ParseAVSDirectorySet is a log parse operation binding the contract event 0x934223b20c24d569ff89796ae10a6997d43e2b3df0c3677fb6ca1f6e37ce344b.
//
// Solidity: event AVSDirectorySet(address indexed avsDirectory)
func (_Mevcommitavs *MevcommitavsFilterer) ParseAVSDirectorySet(log types.Log) (*MevcommitavsAVSDirectorySet, error) {
	event := new(MevcommitavsAVSDirectorySet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "AVSDirectorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsDelegationManagerSetIterator is returned from FilterDelegationManagerSet and is used to iterate over the raw logs and unpacked data for DelegationManagerSet events raised by the Mevcommitavs contract.
type MevcommitavsDelegationManagerSetIterator struct {
	Event *MevcommitavsDelegationManagerSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsDelegationManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsDelegationManagerSet)
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
		it.Event = new(MevcommitavsDelegationManagerSet)
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
func (it *MevcommitavsDelegationManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsDelegationManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsDelegationManagerSet represents a DelegationManagerSet event raised by the Mevcommitavs contract.
type MevcommitavsDelegationManagerSet struct {
	DelegationManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterDelegationManagerSet is a free log retrieval operation binding the contract event 0x2296e6d8aebb5c81250fd381a114c2ec346fc44bc4582ba95cdcac0f09df6cd9.
//
// Solidity: event DelegationManagerSet(address indexed delegationManager)
func (_Mevcommitavs *MevcommitavsFilterer) FilterDelegationManagerSet(opts *bind.FilterOpts, delegationManager []common.Address) (*MevcommitavsDelegationManagerSetIterator, error) {

	var delegationManagerRule []interface{}
	for _, delegationManagerItem := range delegationManager {
		delegationManagerRule = append(delegationManagerRule, delegationManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "DelegationManagerSet", delegationManagerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsDelegationManagerSetIterator{contract: _Mevcommitavs.contract, event: "DelegationManagerSet", logs: logs, sub: sub}, nil
}

// WatchDelegationManagerSet is a free log subscription operation binding the contract event 0x2296e6d8aebb5c81250fd381a114c2ec346fc44bc4582ba95cdcac0f09df6cd9.
//
// Solidity: event DelegationManagerSet(address indexed delegationManager)
func (_Mevcommitavs *MevcommitavsFilterer) WatchDelegationManagerSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsDelegationManagerSet, delegationManager []common.Address) (event.Subscription, error) {

	var delegationManagerRule []interface{}
	for _, delegationManagerItem := range delegationManager {
		delegationManagerRule = append(delegationManagerRule, delegationManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "DelegationManagerSet", delegationManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsDelegationManagerSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "DelegationManagerSet", log); err != nil {
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

// ParseDelegationManagerSet is a log parse operation binding the contract event 0x2296e6d8aebb5c81250fd381a114c2ec346fc44bc4582ba95cdcac0f09df6cd9.
//
// Solidity: event DelegationManagerSet(address indexed delegationManager)
func (_Mevcommitavs *MevcommitavsFilterer) ParseDelegationManagerSet(log types.Log) (*MevcommitavsDelegationManagerSet, error) {
	event := new(MevcommitavsDelegationManagerSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "DelegationManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsEigenPodManagerSetIterator is returned from FilterEigenPodManagerSet and is used to iterate over the raw logs and unpacked data for EigenPodManagerSet events raised by the Mevcommitavs contract.
type MevcommitavsEigenPodManagerSetIterator struct {
	Event *MevcommitavsEigenPodManagerSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsEigenPodManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsEigenPodManagerSet)
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
		it.Event = new(MevcommitavsEigenPodManagerSet)
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
func (it *MevcommitavsEigenPodManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsEigenPodManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsEigenPodManagerSet represents a EigenPodManagerSet event raised by the Mevcommitavs contract.
type MevcommitavsEigenPodManagerSet struct {
	EigenPodManager common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterEigenPodManagerSet is a free log retrieval operation binding the contract event 0x42070ca05aa3fae96a2bb90f36887ecc4894e2e33e748efeb2721962c11fd801.
//
// Solidity: event EigenPodManagerSet(address indexed eigenPodManager)
func (_Mevcommitavs *MevcommitavsFilterer) FilterEigenPodManagerSet(opts *bind.FilterOpts, eigenPodManager []common.Address) (*MevcommitavsEigenPodManagerSetIterator, error) {

	var eigenPodManagerRule []interface{}
	for _, eigenPodManagerItem := range eigenPodManager {
		eigenPodManagerRule = append(eigenPodManagerRule, eigenPodManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "EigenPodManagerSet", eigenPodManagerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsEigenPodManagerSetIterator{contract: _Mevcommitavs.contract, event: "EigenPodManagerSet", logs: logs, sub: sub}, nil
}

// WatchEigenPodManagerSet is a free log subscription operation binding the contract event 0x42070ca05aa3fae96a2bb90f36887ecc4894e2e33e748efeb2721962c11fd801.
//
// Solidity: event EigenPodManagerSet(address indexed eigenPodManager)
func (_Mevcommitavs *MevcommitavsFilterer) WatchEigenPodManagerSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsEigenPodManagerSet, eigenPodManager []common.Address) (event.Subscription, error) {

	var eigenPodManagerRule []interface{}
	for _, eigenPodManagerItem := range eigenPodManager {
		eigenPodManagerRule = append(eigenPodManagerRule, eigenPodManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "EigenPodManagerSet", eigenPodManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsEigenPodManagerSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "EigenPodManagerSet", log); err != nil {
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

// ParseEigenPodManagerSet is a log parse operation binding the contract event 0x42070ca05aa3fae96a2bb90f36887ecc4894e2e33e748efeb2721962c11fd801.
//
// Solidity: event EigenPodManagerSet(address indexed eigenPodManager)
func (_Mevcommitavs *MevcommitavsFilterer) ParseEigenPodManagerSet(log types.Log) (*MevcommitavsEigenPodManagerSet, error) {
	event := new(MevcommitavsEigenPodManagerSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "EigenPodManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsFreezeOracleSetIterator is returned from FilterFreezeOracleSet and is used to iterate over the raw logs and unpacked data for FreezeOracleSet events raised by the Mevcommitavs contract.
type MevcommitavsFreezeOracleSetIterator struct {
	Event *MevcommitavsFreezeOracleSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsFreezeOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsFreezeOracleSet)
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
		it.Event = new(MevcommitavsFreezeOracleSet)
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
func (it *MevcommitavsFreezeOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsFreezeOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsFreezeOracleSet represents a FreezeOracleSet event raised by the Mevcommitavs contract.
type MevcommitavsFreezeOracleSet struct {
	FreezeOracle common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFreezeOracleSet is a free log retrieval operation binding the contract event 0xa33f3723c675820c785c70cde43f95aea5a4a0da3a5443a6cc129e14fcc9455a.
//
// Solidity: event FreezeOracleSet(address indexed freezeOracle)
func (_Mevcommitavs *MevcommitavsFilterer) FilterFreezeOracleSet(opts *bind.FilterOpts, freezeOracle []common.Address) (*MevcommitavsFreezeOracleSetIterator, error) {

	var freezeOracleRule []interface{}
	for _, freezeOracleItem := range freezeOracle {
		freezeOracleRule = append(freezeOracleRule, freezeOracleItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "FreezeOracleSet", freezeOracleRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsFreezeOracleSetIterator{contract: _Mevcommitavs.contract, event: "FreezeOracleSet", logs: logs, sub: sub}, nil
}

// WatchFreezeOracleSet is a free log subscription operation binding the contract event 0xa33f3723c675820c785c70cde43f95aea5a4a0da3a5443a6cc129e14fcc9455a.
//
// Solidity: event FreezeOracleSet(address indexed freezeOracle)
func (_Mevcommitavs *MevcommitavsFilterer) WatchFreezeOracleSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsFreezeOracleSet, freezeOracle []common.Address) (event.Subscription, error) {

	var freezeOracleRule []interface{}
	for _, freezeOracleItem := range freezeOracle {
		freezeOracleRule = append(freezeOracleRule, freezeOracleItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "FreezeOracleSet", freezeOracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsFreezeOracleSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "FreezeOracleSet", log); err != nil {
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

// ParseFreezeOracleSet is a log parse operation binding the contract event 0xa33f3723c675820c785c70cde43f95aea5a4a0da3a5443a6cc129e14fcc9455a.
//
// Solidity: event FreezeOracleSet(address indexed freezeOracle)
func (_Mevcommitavs *MevcommitavsFilterer) ParseFreezeOracleSet(log types.Log) (*MevcommitavsFreezeOracleSet, error) {
	event := new(MevcommitavsFreezeOracleSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "FreezeOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Mevcommitavs contract.
type MevcommitavsInitializedIterator struct {
	Event *MevcommitavsInitialized // Event containing the contract specifics and raw log

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
func (it *MevcommitavsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsInitialized)
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
		it.Event = new(MevcommitavsInitialized)
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
func (it *MevcommitavsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsInitialized represents a Initialized event raised by the Mevcommitavs contract.
type MevcommitavsInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitavs *MevcommitavsFilterer) FilterInitialized(opts *bind.FilterOpts) (*MevcommitavsInitializedIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsInitializedIterator{contract: _Mevcommitavs.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitavs *MevcommitavsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MevcommitavsInitialized) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsInitialized)
				if err := _Mevcommitavs.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseInitialized(log types.Log) (*MevcommitavsInitialized, error) {
	event := new(MevcommitavsInitialized)
	if err := _Mevcommitavs.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator is returned from FilterLSTRestakerDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregPeriodBlocksSet events raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator struct {
	Event *MevcommitavsLSTRestakerDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsLSTRestakerDeregPeriodBlocksSet)
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
		it.Event = new(MevcommitavsLSTRestakerDeregPeriodBlocksSet)
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
func (it *MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsLSTRestakerDeregPeriodBlocksSet represents a LSTRestakerDeregPeriodBlocksSet event raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregPeriodBlocksSet struct {
	LstRestakerDeregPeriodBlocks *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0x7bd82fe806a0299d04c1dc6b928d934d0515dd3dfb8e6b0a0ca02267da5ec181.
//
// Solidity: event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) FilterLSTRestakerDeregPeriodBlocksSet(opts *bind.FilterOpts) (*MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "LSTRestakerDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsLSTRestakerDeregPeriodBlocksSetIterator{contract: _Mevcommitavs.contract, event: "LSTRestakerDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0x7bd82fe806a0299d04c1dc6b928d934d0515dd3dfb8e6b0a0ca02267da5ec181.
//
// Solidity: event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) WatchLSTRestakerDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsLSTRestakerDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "LSTRestakerDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsLSTRestakerDeregPeriodBlocksSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregPeriodBlocksSet", log); err != nil {
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

// ParseLSTRestakerDeregPeriodBlocksSet is a log parse operation binding the contract event 0x7bd82fe806a0299d04c1dc6b928d934d0515dd3dfb8e6b0a0ca02267da5ec181.
//
// Solidity: event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) ParseLSTRestakerDeregPeriodBlocksSet(log types.Log) (*MevcommitavsLSTRestakerDeregPeriodBlocksSet, error) {
	event := new(MevcommitavsLSTRestakerDeregPeriodBlocksSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsLSTRestakerDeregisteredIterator is returned from FilterLSTRestakerDeregistered and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregistered events raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregisteredIterator struct {
	Event *MevcommitavsLSTRestakerDeregistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsLSTRestakerDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsLSTRestakerDeregistered)
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
		it.Event = new(MevcommitavsLSTRestakerDeregistered)
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
func (it *MevcommitavsLSTRestakerDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsLSTRestakerDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsLSTRestakerDeregistered represents a LSTRestakerDeregistered event raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregistered struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregistered is a free log retrieval operation binding the contract event 0xaf3d14c12fe6a17d9cab68818354270a52de8de87f44c614c8cf7c35e96086fd.
//
// Solidity: event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) FilterLSTRestakerDeregistered(opts *bind.FilterOpts, lstRestaker []common.Address) (*MevcommitavsLSTRestakerDeregisteredIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "LSTRestakerDeregistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsLSTRestakerDeregisteredIterator{contract: _Mevcommitavs.contract, event: "LSTRestakerDeregistered", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregistered is a free log subscription operation binding the contract event 0xaf3d14c12fe6a17d9cab68818354270a52de8de87f44c614c8cf7c35e96086fd.
//
// Solidity: event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) WatchLSTRestakerDeregistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsLSTRestakerDeregistered, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "LSTRestakerDeregistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsLSTRestakerDeregistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregistered", log); err != nil {
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

// ParseLSTRestakerDeregistered is a log parse operation binding the contract event 0xaf3d14c12fe6a17d9cab68818354270a52de8de87f44c614c8cf7c35e96086fd.
//
// Solidity: event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) ParseLSTRestakerDeregistered(log types.Log) (*MevcommitavsLSTRestakerDeregistered, error) {
	event := new(MevcommitavsLSTRestakerDeregistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsLSTRestakerDeregistrationRequestedIterator is returned from FilterLSTRestakerDeregistrationRequested and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregistrationRequested events raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregistrationRequestedIterator struct {
	Event *MevcommitavsLSTRestakerDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitavsLSTRestakerDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsLSTRestakerDeregistrationRequested)
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
		it.Event = new(MevcommitavsLSTRestakerDeregistrationRequested)
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
func (it *MevcommitavsLSTRestakerDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsLSTRestakerDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsLSTRestakerDeregistrationRequested represents a LSTRestakerDeregistrationRequested event raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerDeregistrationRequested struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregistrationRequested is a free log retrieval operation binding the contract event 0x6a7ec4f90fd0baec2ae0d35e6c4731b6c833f9f082a3fe8f5b30fe1be4af1c3b.
//
// Solidity: event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) FilterLSTRestakerDeregistrationRequested(opts *bind.FilterOpts, lstRestaker []common.Address) (*MevcommitavsLSTRestakerDeregistrationRequestedIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "LSTRestakerDeregistrationRequested", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsLSTRestakerDeregistrationRequestedIterator{contract: _Mevcommitavs.contract, event: "LSTRestakerDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregistrationRequested is a free log subscription operation binding the contract event 0x6a7ec4f90fd0baec2ae0d35e6c4731b6c833f9f082a3fe8f5b30fe1be4af1c3b.
//
// Solidity: event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) WatchLSTRestakerDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitavsLSTRestakerDeregistrationRequested, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "LSTRestakerDeregistrationRequested", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsLSTRestakerDeregistrationRequested)
				if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregistrationRequested", log); err != nil {
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

// ParseLSTRestakerDeregistrationRequested is a log parse operation binding the contract event 0x6a7ec4f90fd0baec2ae0d35e6c4731b6c833f9f082a3fe8f5b30fe1be4af1c3b.
//
// Solidity: event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) ParseLSTRestakerDeregistrationRequested(log types.Log) (*MevcommitavsLSTRestakerDeregistrationRequested, error) {
	event := new(MevcommitavsLSTRestakerDeregistrationRequested)
	if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsLSTRestakerRegisteredIterator is returned from FilterLSTRestakerRegistered and is used to iterate over the raw logs and unpacked data for LSTRestakerRegistered events raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerRegisteredIterator struct {
	Event *MevcommitavsLSTRestakerRegistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsLSTRestakerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsLSTRestakerRegistered)
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
		it.Event = new(MevcommitavsLSTRestakerRegistered)
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
func (it *MevcommitavsLSTRestakerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsLSTRestakerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsLSTRestakerRegistered represents a LSTRestakerRegistered event raised by the Mevcommitavs contract.
type MevcommitavsLSTRestakerRegistered struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerRegistered is a free log retrieval operation binding the contract event 0xdecaa7bc9a78e41b524d10750b5a52f0c8c1144cd5ab2991c9ca75ff380e011a.
//
// Solidity: event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) FilterLSTRestakerRegistered(opts *bind.FilterOpts, lstRestaker []common.Address) (*MevcommitavsLSTRestakerRegisteredIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "LSTRestakerRegistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsLSTRestakerRegisteredIterator{contract: _Mevcommitavs.contract, event: "LSTRestakerRegistered", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerRegistered is a free log subscription operation binding the contract event 0xdecaa7bc9a78e41b524d10750b5a52f0c8c1144cd5ab2991c9ca75ff380e011a.
//
// Solidity: event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) WatchLSTRestakerRegistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsLSTRestakerRegistered, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "LSTRestakerRegistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsLSTRestakerRegistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerRegistered", log); err != nil {
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

// ParseLSTRestakerRegistered is a log parse operation binding the contract event 0xdecaa7bc9a78e41b524d10750b5a52f0c8c1144cd5ab2991c9ca75ff380e011a.
//
// Solidity: event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavs *MevcommitavsFilterer) ParseLSTRestakerRegistered(log types.Log) (*MevcommitavsLSTRestakerRegistered, error) {
	event := new(MevcommitavsLSTRestakerRegistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "LSTRestakerRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOperatorDeregPeriodBlocksSetIterator is returned from FilterOperatorDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for OperatorDeregPeriodBlocksSet events raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregPeriodBlocksSetIterator struct {
	Event *MevcommitavsOperatorDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOperatorDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOperatorDeregPeriodBlocksSet)
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
		it.Event = new(MevcommitavsOperatorDeregPeriodBlocksSet)
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
func (it *MevcommitavsOperatorDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOperatorDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOperatorDeregPeriodBlocksSet represents a OperatorDeregPeriodBlocksSet event raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregPeriodBlocksSet struct {
	OperatorDeregPeriodBlocks *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0xfea621e39fd1186d690d8fa903a946ee52fee14c9c1f1c7295173b2e623b517e.
//
// Solidity: event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOperatorDeregPeriodBlocksSet(opts *bind.FilterOpts) (*MevcommitavsOperatorDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OperatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOperatorDeregPeriodBlocksSetIterator{contract: _Mevcommitavs.contract, event: "OperatorDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0xfea621e39fd1186d690d8fa903a946ee52fee14c9c1f1c7295173b2e623b517e.
//
// Solidity: event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOperatorDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsOperatorDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OperatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOperatorDeregPeriodBlocksSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregPeriodBlocksSet", log); err != nil {
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

// ParseOperatorDeregPeriodBlocksSet is a log parse operation binding the contract event 0xfea621e39fd1186d690d8fa903a946ee52fee14c9c1f1c7295173b2e623b517e.
//
// Solidity: event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) ParseOperatorDeregPeriodBlocksSet(log types.Log) (*MevcommitavsOperatorDeregPeriodBlocksSet, error) {
	event := new(MevcommitavsOperatorDeregPeriodBlocksSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOperatorDeregisteredIterator is returned from FilterOperatorDeregistered and is used to iterate over the raw logs and unpacked data for OperatorDeregistered events raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregisteredIterator struct {
	Event *MevcommitavsOperatorDeregistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOperatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOperatorDeregistered)
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
		it.Event = new(MevcommitavsOperatorDeregistered)
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
func (it *MevcommitavsOperatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOperatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOperatorDeregistered represents a OperatorDeregistered event raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistered is a free log retrieval operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOperatorDeregistered(opts *bind.FilterOpts, operator []common.Address) (*MevcommitavsOperatorDeregisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOperatorDeregisteredIterator{contract: _Mevcommitavs.contract, event: "OperatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistered is a free log subscription operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOperatorDeregistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsOperatorDeregistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOperatorDeregistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseOperatorDeregistered(log types.Log) (*MevcommitavsOperatorDeregistered, error) {
	event := new(MevcommitavsOperatorDeregistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOperatorDeregistrationRequestedIterator is returned from FilterOperatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for OperatorDeregistrationRequested events raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregistrationRequestedIterator struct {
	Event *MevcommitavsOperatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOperatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOperatorDeregistrationRequested)
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
		it.Event = new(MevcommitavsOperatorDeregistrationRequested)
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
func (it *MevcommitavsOperatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOperatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOperatorDeregistrationRequested represents a OperatorDeregistrationRequested event raised by the Mevcommitavs contract.
type MevcommitavsOperatorDeregistrationRequested struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOperatorDeregistrationRequested(opts *bind.FilterOpts, operator []common.Address) (*MevcommitavsOperatorDeregistrationRequestedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOperatorDeregistrationRequestedIterator{contract: _Mevcommitavs.contract, event: "OperatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistrationRequested is a free log subscription operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOperatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitavsOperatorDeregistrationRequested, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOperatorDeregistrationRequested)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseOperatorDeregistrationRequested(log types.Log) (*MevcommitavsOperatorDeregistrationRequested, error) {
	event := new(MevcommitavsOperatorDeregistrationRequested)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the Mevcommitavs contract.
type MevcommitavsOperatorRegisteredIterator struct {
	Event *MevcommitavsOperatorRegistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOperatorRegistered)
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
		it.Event = new(MevcommitavsOperatorRegistered)
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
func (it *MevcommitavsOperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOperatorRegistered represents a OperatorRegistered event raised by the Mevcommitavs contract.
type MevcommitavsOperatorRegistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address) (*MevcommitavsOperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOperatorRegisteredIterator{contract: _Mevcommitavs.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsOperatorRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOperatorRegistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseOperatorRegistered(log types.Log) (*MevcommitavsOperatorRegistered, error) {
	event := new(MevcommitavsOperatorRegistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Mevcommitavs contract.
type MevcommitavsOwnershipTransferStartedIterator struct {
	Event *MevcommitavsOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOwnershipTransferStarted)
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
		it.Event = new(MevcommitavsOwnershipTransferStarted)
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
func (it *MevcommitavsOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Mevcommitavs contract.
type MevcommitavsOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MevcommitavsOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOwnershipTransferStartedIterator{contract: _Mevcommitavs.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *MevcommitavsOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOwnershipTransferStarted)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseOwnershipTransferStarted(log types.Log) (*MevcommitavsOwnershipTransferStarted, error) {
	event := new(MevcommitavsOwnershipTransferStarted)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Mevcommitavs contract.
type MevcommitavsOwnershipTransferredIterator struct {
	Event *MevcommitavsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MevcommitavsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsOwnershipTransferred)
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
		it.Event = new(MevcommitavsOwnershipTransferred)
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
func (it *MevcommitavsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsOwnershipTransferred represents a OwnershipTransferred event raised by the Mevcommitavs contract.
type MevcommitavsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MevcommitavsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsOwnershipTransferredIterator{contract: _Mevcommitavs.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MevcommitavsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsOwnershipTransferred)
				if err := _Mevcommitavs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseOwnershipTransferred(log types.Log) (*MevcommitavsOwnershipTransferred, error) {
	event := new(MevcommitavsOwnershipTransferred)
	if err := _Mevcommitavs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Mevcommitavs contract.
type MevcommitavsPausedIterator struct {
	Event *MevcommitavsPaused // Event containing the contract specifics and raw log

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
func (it *MevcommitavsPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsPaused)
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
		it.Event = new(MevcommitavsPaused)
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
func (it *MevcommitavsPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsPaused represents a Paused event raised by the Mevcommitavs contract.
type MevcommitavsPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitavs *MevcommitavsFilterer) FilterPaused(opts *bind.FilterOpts) (*MevcommitavsPausedIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsPausedIterator{contract: _Mevcommitavs.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitavs *MevcommitavsFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MevcommitavsPaused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsPaused)
				if err := _Mevcommitavs.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParsePaused(log types.Log) (*MevcommitavsPaused, error) {
	event := new(MevcommitavsPaused)
	if err := _Mevcommitavs.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsRestakeableStrategiesSetIterator is returned from FilterRestakeableStrategiesSet and is used to iterate over the raw logs and unpacked data for RestakeableStrategiesSet events raised by the Mevcommitavs contract.
type MevcommitavsRestakeableStrategiesSetIterator struct {
	Event *MevcommitavsRestakeableStrategiesSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsRestakeableStrategiesSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsRestakeableStrategiesSet)
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
		it.Event = new(MevcommitavsRestakeableStrategiesSet)
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
func (it *MevcommitavsRestakeableStrategiesSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsRestakeableStrategiesSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsRestakeableStrategiesSet represents a RestakeableStrategiesSet event raised by the Mevcommitavs contract.
type MevcommitavsRestakeableStrategiesSet struct {
	RestakeableStrategies []common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterRestakeableStrategiesSet is a free log retrieval operation binding the contract event 0xda4dd29046c55387af7520737a1e06033ae31f610dde3d0851458dffe13a0c0f.
//
// Solidity: event RestakeableStrategiesSet(address[] indexed restakeableStrategies)
func (_Mevcommitavs *MevcommitavsFilterer) FilterRestakeableStrategiesSet(opts *bind.FilterOpts, restakeableStrategies [][]common.Address) (*MevcommitavsRestakeableStrategiesSetIterator, error) {

	var restakeableStrategiesRule []interface{}
	for _, restakeableStrategiesItem := range restakeableStrategies {
		restakeableStrategiesRule = append(restakeableStrategiesRule, restakeableStrategiesItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "RestakeableStrategiesSet", restakeableStrategiesRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsRestakeableStrategiesSetIterator{contract: _Mevcommitavs.contract, event: "RestakeableStrategiesSet", logs: logs, sub: sub}, nil
}

// WatchRestakeableStrategiesSet is a free log subscription operation binding the contract event 0xda4dd29046c55387af7520737a1e06033ae31f610dde3d0851458dffe13a0c0f.
//
// Solidity: event RestakeableStrategiesSet(address[] indexed restakeableStrategies)
func (_Mevcommitavs *MevcommitavsFilterer) WatchRestakeableStrategiesSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsRestakeableStrategiesSet, restakeableStrategies [][]common.Address) (event.Subscription, error) {

	var restakeableStrategiesRule []interface{}
	for _, restakeableStrategiesItem := range restakeableStrategies {
		restakeableStrategiesRule = append(restakeableStrategiesRule, restakeableStrategiesItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "RestakeableStrategiesSet", restakeableStrategiesRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsRestakeableStrategiesSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "RestakeableStrategiesSet", log); err != nil {
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

// ParseRestakeableStrategiesSet is a log parse operation binding the contract event 0xda4dd29046c55387af7520737a1e06033ae31f610dde3d0851458dffe13a0c0f.
//
// Solidity: event RestakeableStrategiesSet(address[] indexed restakeableStrategies)
func (_Mevcommitavs *MevcommitavsFilterer) ParseRestakeableStrategiesSet(log types.Log) (*MevcommitavsRestakeableStrategiesSet, error) {
	event := new(MevcommitavsRestakeableStrategiesSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "RestakeableStrategiesSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsStrategyManagerSetIterator is returned from FilterStrategyManagerSet and is used to iterate over the raw logs and unpacked data for StrategyManagerSet events raised by the Mevcommitavs contract.
type MevcommitavsStrategyManagerSetIterator struct {
	Event *MevcommitavsStrategyManagerSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsStrategyManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsStrategyManagerSet)
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
		it.Event = new(MevcommitavsStrategyManagerSet)
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
func (it *MevcommitavsStrategyManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsStrategyManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsStrategyManagerSet represents a StrategyManagerSet event raised by the Mevcommitavs contract.
type MevcommitavsStrategyManagerSet struct {
	StrategyManager common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStrategyManagerSet is a free log retrieval operation binding the contract event 0x76fe640b216c20f563ab2d807634271e9d772e92c8a3752325cb2bc924e9e514.
//
// Solidity: event StrategyManagerSet(address indexed strategyManager)
func (_Mevcommitavs *MevcommitavsFilterer) FilterStrategyManagerSet(opts *bind.FilterOpts, strategyManager []common.Address) (*MevcommitavsStrategyManagerSetIterator, error) {

	var strategyManagerRule []interface{}
	for _, strategyManagerItem := range strategyManager {
		strategyManagerRule = append(strategyManagerRule, strategyManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "StrategyManagerSet", strategyManagerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsStrategyManagerSetIterator{contract: _Mevcommitavs.contract, event: "StrategyManagerSet", logs: logs, sub: sub}, nil
}

// WatchStrategyManagerSet is a free log subscription operation binding the contract event 0x76fe640b216c20f563ab2d807634271e9d772e92c8a3752325cb2bc924e9e514.
//
// Solidity: event StrategyManagerSet(address indexed strategyManager)
func (_Mevcommitavs *MevcommitavsFilterer) WatchStrategyManagerSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsStrategyManagerSet, strategyManager []common.Address) (event.Subscription, error) {

	var strategyManagerRule []interface{}
	for _, strategyManagerItem := range strategyManager {
		strategyManagerRule = append(strategyManagerRule, strategyManagerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "StrategyManagerSet", strategyManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsStrategyManagerSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "StrategyManagerSet", log); err != nil {
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

// ParseStrategyManagerSet is a log parse operation binding the contract event 0x76fe640b216c20f563ab2d807634271e9d772e92c8a3752325cb2bc924e9e514.
//
// Solidity: event StrategyManagerSet(address indexed strategyManager)
func (_Mevcommitavs *MevcommitavsFilterer) ParseStrategyManagerSet(log types.Log) (*MevcommitavsStrategyManagerSet, error) {
	event := new(MevcommitavsStrategyManagerSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "StrategyManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsUnfreezeFeeSetIterator is returned from FilterUnfreezeFeeSet and is used to iterate over the raw logs and unpacked data for UnfreezeFeeSet events raised by the Mevcommitavs contract.
type MevcommitavsUnfreezeFeeSetIterator struct {
	Event *MevcommitavsUnfreezeFeeSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsUnfreezeFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsUnfreezeFeeSet)
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
		it.Event = new(MevcommitavsUnfreezeFeeSet)
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
func (it *MevcommitavsUnfreezeFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsUnfreezeFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsUnfreezeFeeSet represents a UnfreezeFeeSet event raised by the Mevcommitavs contract.
type MevcommitavsUnfreezeFeeSet struct {
	UnfreezeFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnfreezeFeeSet is a free log retrieval operation binding the contract event 0x6c0cf79356801bf6665a4c6cc85d35896ca003fe22c8c92c2b0e1d563b384c9d.
//
// Solidity: event UnfreezeFeeSet(uint256 unfreezeFee)
func (_Mevcommitavs *MevcommitavsFilterer) FilterUnfreezeFeeSet(opts *bind.FilterOpts) (*MevcommitavsUnfreezeFeeSetIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "UnfreezeFeeSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsUnfreezeFeeSetIterator{contract: _Mevcommitavs.contract, event: "UnfreezeFeeSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezeFeeSet is a free log subscription operation binding the contract event 0x6c0cf79356801bf6665a4c6cc85d35896ca003fe22c8c92c2b0e1d563b384c9d.
//
// Solidity: event UnfreezeFeeSet(uint256 unfreezeFee)
func (_Mevcommitavs *MevcommitavsFilterer) WatchUnfreezeFeeSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsUnfreezeFeeSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "UnfreezeFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsUnfreezeFeeSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezeFeeSet", log); err != nil {
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

// ParseUnfreezeFeeSet is a log parse operation binding the contract event 0x6c0cf79356801bf6665a4c6cc85d35896ca003fe22c8c92c2b0e1d563b384c9d.
//
// Solidity: event UnfreezeFeeSet(uint256 unfreezeFee)
func (_Mevcommitavs *MevcommitavsFilterer) ParseUnfreezeFeeSet(log types.Log) (*MevcommitavsUnfreezeFeeSet, error) {
	event := new(MevcommitavsUnfreezeFeeSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezeFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsUnfreezePeriodBlocksSetIterator is returned from FilterUnfreezePeriodBlocksSet and is used to iterate over the raw logs and unpacked data for UnfreezePeriodBlocksSet events raised by the Mevcommitavs contract.
type MevcommitavsUnfreezePeriodBlocksSetIterator struct {
	Event *MevcommitavsUnfreezePeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsUnfreezePeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsUnfreezePeriodBlocksSet)
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
		it.Event = new(MevcommitavsUnfreezePeriodBlocksSet)
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
func (it *MevcommitavsUnfreezePeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsUnfreezePeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsUnfreezePeriodBlocksSet represents a UnfreezePeriodBlocksSet event raised by the Mevcommitavs contract.
type MevcommitavsUnfreezePeriodBlocksSet struct {
	UnfreezePeriodBlocks *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterUnfreezePeriodBlocksSet is a free log retrieval operation binding the contract event 0xdf42b0e9b907c6972f5ad1b757b3ab4e0eeca44f9b9dc7c8d97f2f40c1a042dd.
//
// Solidity: event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) FilterUnfreezePeriodBlocksSet(opts *bind.FilterOpts) (*MevcommitavsUnfreezePeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "UnfreezePeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsUnfreezePeriodBlocksSetIterator{contract: _Mevcommitavs.contract, event: "UnfreezePeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezePeriodBlocksSet is a free log subscription operation binding the contract event 0xdf42b0e9b907c6972f5ad1b757b3ab4e0eeca44f9b9dc7c8d97f2f40c1a042dd.
//
// Solidity: event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) WatchUnfreezePeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsUnfreezePeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "UnfreezePeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsUnfreezePeriodBlocksSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezePeriodBlocksSet", log); err != nil {
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

// ParseUnfreezePeriodBlocksSet is a log parse operation binding the contract event 0xdf42b0e9b907c6972f5ad1b757b3ab4e0eeca44f9b9dc7c8d97f2f40c1a042dd.
//
// Solidity: event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) ParseUnfreezePeriodBlocksSet(log types.Log) (*MevcommitavsUnfreezePeriodBlocksSet, error) {
	event := new(MevcommitavsUnfreezePeriodBlocksSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezePeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsUnfreezeReceiverSetIterator is returned from FilterUnfreezeReceiverSet and is used to iterate over the raw logs and unpacked data for UnfreezeReceiverSet events raised by the Mevcommitavs contract.
type MevcommitavsUnfreezeReceiverSetIterator struct {
	Event *MevcommitavsUnfreezeReceiverSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsUnfreezeReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsUnfreezeReceiverSet)
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
		it.Event = new(MevcommitavsUnfreezeReceiverSet)
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
func (it *MevcommitavsUnfreezeReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsUnfreezeReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsUnfreezeReceiverSet represents a UnfreezeReceiverSet event raised by the Mevcommitavs contract.
type MevcommitavsUnfreezeReceiverSet struct {
	UnfreezeReceiver common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUnfreezeReceiverSet is a free log retrieval operation binding the contract event 0x7b992f6c09b22e43cad1f0f9a1d7e41949ac1c53a4eebb7f4b2374605049d2bf.
//
// Solidity: event UnfreezeReceiverSet(address indexed unfreezeReceiver)
func (_Mevcommitavs *MevcommitavsFilterer) FilterUnfreezeReceiverSet(opts *bind.FilterOpts, unfreezeReceiver []common.Address) (*MevcommitavsUnfreezeReceiverSetIterator, error) {

	var unfreezeReceiverRule []interface{}
	for _, unfreezeReceiverItem := range unfreezeReceiver {
		unfreezeReceiverRule = append(unfreezeReceiverRule, unfreezeReceiverItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "UnfreezeReceiverSet", unfreezeReceiverRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsUnfreezeReceiverSetIterator{contract: _Mevcommitavs.contract, event: "UnfreezeReceiverSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezeReceiverSet is a free log subscription operation binding the contract event 0x7b992f6c09b22e43cad1f0f9a1d7e41949ac1c53a4eebb7f4b2374605049d2bf.
//
// Solidity: event UnfreezeReceiverSet(address indexed unfreezeReceiver)
func (_Mevcommitavs *MevcommitavsFilterer) WatchUnfreezeReceiverSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsUnfreezeReceiverSet, unfreezeReceiver []common.Address) (event.Subscription, error) {

	var unfreezeReceiverRule []interface{}
	for _, unfreezeReceiverItem := range unfreezeReceiver {
		unfreezeReceiverRule = append(unfreezeReceiverRule, unfreezeReceiverItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "UnfreezeReceiverSet", unfreezeReceiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsUnfreezeReceiverSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezeReceiverSet", log); err != nil {
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

// ParseUnfreezeReceiverSet is a log parse operation binding the contract event 0x7b992f6c09b22e43cad1f0f9a1d7e41949ac1c53a4eebb7f4b2374605049d2bf.
//
// Solidity: event UnfreezeReceiverSet(address indexed unfreezeReceiver)
func (_Mevcommitavs *MevcommitavsFilterer) ParseUnfreezeReceiverSet(log types.Log) (*MevcommitavsUnfreezeReceiverSet, error) {
	event := new(MevcommitavsUnfreezeReceiverSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "UnfreezeReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Mevcommitavs contract.
type MevcommitavsUnpausedIterator struct {
	Event *MevcommitavsUnpaused // Event containing the contract specifics and raw log

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
func (it *MevcommitavsUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsUnpaused)
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
		it.Event = new(MevcommitavsUnpaused)
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
func (it *MevcommitavsUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsUnpaused represents a Unpaused event raised by the Mevcommitavs contract.
type MevcommitavsUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitavs *MevcommitavsFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MevcommitavsUnpausedIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsUnpausedIterator{contract: _Mevcommitavs.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitavs *MevcommitavsFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MevcommitavsUnpaused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsUnpaused)
				if err := _Mevcommitavs.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseUnpaused(log types.Log) (*MevcommitavsUnpaused, error) {
	event := new(MevcommitavsUnpaused)
	if err := _Mevcommitavs.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Mevcommitavs contract.
type MevcommitavsUpgradedIterator struct {
	Event *MevcommitavsUpgraded // Event containing the contract specifics and raw log

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
func (it *MevcommitavsUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsUpgraded)
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
		it.Event = new(MevcommitavsUpgraded)
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
func (it *MevcommitavsUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsUpgraded represents a Upgraded event raised by the Mevcommitavs contract.
type MevcommitavsUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitavs *MevcommitavsFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*MevcommitavsUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsUpgradedIterator{contract: _Mevcommitavs.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitavs *MevcommitavsFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *MevcommitavsUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsUpgraded)
				if err := _Mevcommitavs.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Mevcommitavs *MevcommitavsFilterer) ParseUpgraded(log types.Log) (*MevcommitavsUpgraded, error) {
	event := new(MevcommitavsUpgraded)
	if err := _Mevcommitavs.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorDeregPeriodBlocksSetIterator is returned from FilterValidatorDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for ValidatorDeregPeriodBlocksSet events raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregPeriodBlocksSetIterator struct {
	Event *MevcommitavsValidatorDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorDeregPeriodBlocksSet)
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
		it.Event = new(MevcommitavsValidatorDeregPeriodBlocksSet)
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
func (it *MevcommitavsValidatorDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorDeregPeriodBlocksSet represents a ValidatorDeregPeriodBlocksSet event raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregPeriodBlocksSet struct {
	ValidatorDeregPeriodBlocks *big.Int
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0x38755fd87ce522d770f5627d05f541406a9358d853b64009737f5c4d8913eae7.
//
// Solidity: event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorDeregPeriodBlocksSet(opts *bind.FilterOpts) (*MevcommitavsValidatorDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorDeregPeriodBlocksSetIterator{contract: _Mevcommitavs.contract, event: "ValidatorDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0x38755fd87ce522d770f5627d05f541406a9358d853b64009737f5c4d8913eae7.
//
// Solidity: event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorDeregPeriodBlocksSet)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregPeriodBlocksSet", log); err != nil {
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

// ParseValidatorDeregPeriodBlocksSet is a log parse operation binding the contract event 0x38755fd87ce522d770f5627d05f541406a9358d853b64009737f5c4d8913eae7.
//
// Solidity: event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorDeregPeriodBlocksSet(log types.Log) (*MevcommitavsValidatorDeregPeriodBlocksSet, error) {
	event := new(MevcommitavsValidatorDeregPeriodBlocksSet)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorDeregisteredIterator is returned from FilterValidatorDeregistered and is used to iterate over the raw logs and unpacked data for ValidatorDeregistered events raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregisteredIterator struct {
	Event *MevcommitavsValidatorDeregistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorDeregistered)
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
		it.Event = new(MevcommitavsValidatorDeregistered)
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
func (it *MevcommitavsValidatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorDeregistered represents a ValidatorDeregistered event raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregistered struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistered is a free log retrieval operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorDeregistered(opts *bind.FilterOpts, podOwner []common.Address) (*MevcommitavsValidatorDeregisteredIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorDeregistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorDeregisteredIterator{contract: _Mevcommitavs.contract, event: "ValidatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistered is a free log subscription operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorDeregistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorDeregistered, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorDeregistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorDeregistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
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

// ParseValidatorDeregistered is a log parse operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorDeregistered(log types.Log) (*MevcommitavsValidatorDeregistered, error) {
	event := new(MevcommitavsValidatorDeregistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorDeregistrationRequestedIterator is returned from FilterValidatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for ValidatorDeregistrationRequested events raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregistrationRequestedIterator struct {
	Event *MevcommitavsValidatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorDeregistrationRequested)
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
		it.Event = new(MevcommitavsValidatorDeregistrationRequested)
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
func (it *MevcommitavsValidatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorDeregistrationRequested represents a ValidatorDeregistrationRequested event raised by the Mevcommitavs contract.
type MevcommitavsValidatorDeregistrationRequested struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorDeregistrationRequested(opts *bind.FilterOpts, podOwner []common.Address) (*MevcommitavsValidatorDeregistrationRequestedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorDeregistrationRequested", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorDeregistrationRequestedIterator{contract: _Mevcommitavs.contract, event: "ValidatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistrationRequested is a free log subscription operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorDeregistrationRequested, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorDeregistrationRequested", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorDeregistrationRequested)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
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

// ParseValidatorDeregistrationRequested is a log parse operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorDeregistrationRequested(log types.Log) (*MevcommitavsValidatorDeregistrationRequested, error) {
	event := new(MevcommitavsValidatorDeregistrationRequested)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorFrozenIterator is returned from FilterValidatorFrozen and is used to iterate over the raw logs and unpacked data for ValidatorFrozen events raised by the Mevcommitavs contract.
type MevcommitavsValidatorFrozenIterator struct {
	Event *MevcommitavsValidatorFrozen // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorFrozen)
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
		it.Event = new(MevcommitavsValidatorFrozen)
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
func (it *MevcommitavsValidatorFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorFrozen represents a ValidatorFrozen event raised by the Mevcommitavs contract.
type MevcommitavsValidatorFrozen struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorFrozen is a free log retrieval operation binding the contract event 0x5f565c1dd6cf6dc33dbdfd22c94b541af6ee1390251c8975b0c84106c58654bf.
//
// Solidity: event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorFrozen(opts *bind.FilterOpts, podOwner []common.Address) (*MevcommitavsValidatorFrozenIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorFrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorFrozenIterator{contract: _Mevcommitavs.contract, event: "ValidatorFrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorFrozen is a free log subscription operation binding the contract event 0x5f565c1dd6cf6dc33dbdfd22c94b541af6ee1390251c8975b0c84106c58654bf.
//
// Solidity: event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorFrozen(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorFrozen, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorFrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorFrozen)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
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

// ParseValidatorFrozen is a log parse operation binding the contract event 0x5f565c1dd6cf6dc33dbdfd22c94b541af6ee1390251c8975b0c84106c58654bf.
//
// Solidity: event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorFrozen(log types.Log) (*MevcommitavsValidatorFrozen, error) {
	event := new(MevcommitavsValidatorFrozen)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorRegisteredIterator is returned from FilterValidatorRegistered and is used to iterate over the raw logs and unpacked data for ValidatorRegistered events raised by the Mevcommitavs contract.
type MevcommitavsValidatorRegisteredIterator struct {
	Event *MevcommitavsValidatorRegistered // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorRegistered)
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
		it.Event = new(MevcommitavsValidatorRegistered)
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
func (it *MevcommitavsValidatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorRegistered represents a ValidatorRegistered event raised by the Mevcommitavs contract.
type MevcommitavsValidatorRegistered struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorRegistered is a free log retrieval operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorRegistered(opts *bind.FilterOpts, podOwner []common.Address) (*MevcommitavsValidatorRegisteredIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorRegistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorRegisteredIterator{contract: _Mevcommitavs.contract, event: "ValidatorRegistered", logs: logs, sub: sub}, nil
}

// WatchValidatorRegistered is a free log subscription operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorRegistered(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorRegistered, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorRegistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorRegistered)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
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

// ParseValidatorRegistered is a log parse operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorRegistered(log types.Log) (*MevcommitavsValidatorRegistered, error) {
	event := new(MevcommitavsValidatorRegistered)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MevcommitavsValidatorUnfrozenIterator is returned from FilterValidatorUnfrozen and is used to iterate over the raw logs and unpacked data for ValidatorUnfrozen events raised by the Mevcommitavs contract.
type MevcommitavsValidatorUnfrozenIterator struct {
	Event *MevcommitavsValidatorUnfrozen // Event containing the contract specifics and raw log

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
func (it *MevcommitavsValidatorUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MevcommitavsValidatorUnfrozen)
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
		it.Event = new(MevcommitavsValidatorUnfrozen)
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
func (it *MevcommitavsValidatorUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MevcommitavsValidatorUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MevcommitavsValidatorUnfrozen represents a ValidatorUnfrozen event raised by the Mevcommitavs contract.
type MevcommitavsValidatorUnfrozen struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorUnfrozen is a free log retrieval operation binding the contract event 0x8b4f548363e5182887ee88395b0bacbd44ef955d9b9c1ace7aad0da43af0de40.
//
// Solidity: event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) FilterValidatorUnfrozen(opts *bind.FilterOpts, podOwner []common.Address) (*MevcommitavsValidatorUnfrozenIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.FilterLogs(opts, "ValidatorUnfrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MevcommitavsValidatorUnfrozenIterator{contract: _Mevcommitavs.contract, event: "ValidatorUnfrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorUnfrozen is a free log subscription operation binding the contract event 0x8b4f548363e5182887ee88395b0bacbd44ef955d9b9c1ace7aad0da43af0de40.
//
// Solidity: event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) WatchValidatorUnfrozen(opts *bind.WatchOpts, sink chan<- *MevcommitavsValidatorUnfrozen, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavs.contract.WatchLogs(opts, "ValidatorUnfrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MevcommitavsValidatorUnfrozen)
				if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
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

// ParseValidatorUnfrozen is a log parse operation binding the contract event 0x8b4f548363e5182887ee88395b0bacbd44ef955d9b9c1ace7aad0da43af0de40.
//
// Solidity: event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavs *MevcommitavsFilterer) ParseValidatorUnfrozen(log types.Log) (*MevcommitavsValidatorUnfrozen, error) {
	event := new(MevcommitavsValidatorUnfrozen)
	if err := _Mevcommitavs.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
