// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mevcommitavsv3

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

// EventHeightLibEventHeight is an auto generated low-level Go binding around an user-defined struct.
type EventHeightLibEventHeight struct {
	Exists      bool
	BlockHeight *big.Int
}

// IMevCommitAVSV3LSTRestakerRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSV3LSTRestakerRegistrationInfo struct {
	Exists             bool
	ChosenValidators   [][]byte
	NumChosen          *big.Int
	DeregRequestHeight EventHeightLibEventHeight
}

// IMevCommitAVSV3OperatorRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSV3OperatorRegistrationInfo struct {
	Exists             bool
	DeregRequestHeight EventHeightLibEventHeight
}

// IMevCommitAVSV3ValidatorRegistrationInfo is an auto generated low-level Go binding around an user-defined struct.
type IMevCommitAVSV3ValidatorRegistrationInfo struct {
	Exists             bool
	PodOwner           common.Address
	FreezeHeight       EventHeightLibEventHeight
	DeregRequestHeight EventHeightLibEventHeight
}

// ISignatureUtilsSignatureWithSaltAndExpiry is an auto generated low-level Go binding around an user-defined struct.
type ISignatureUtilsSignatureWithSaltAndExpiry struct {
	Signature []byte
	Salt      [32]byte
	Expiry    *big.Int
}

// Mevcommitavsv3MetaData contains all meta data concerning the Mevcommitavsv3 contract.
var Mevcommitavsv3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"avsDirectory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterLSTRestaker\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deregisterValidators\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freeze\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeOracle\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLSTRestakerRegInfo\",\"inputs\":[{\"name\":\"lstRestaker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVSV3.LSTRestakerRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"chosenValidators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRegInfo\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVSV3.OperatorRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOperatorRestakedStrategies\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRestakeableStrategies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValidatorRegInfo\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIMevCommitAVSV3.ValidatorRegistrationInfo\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freezeHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"delegationManager_\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"},{\"name\":\"eigenPodManager_\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"},{\"name\":\"strategyManager_\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"},{\"name\":\"avsDirectory_\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"},{\"name\":\"restakeableStrategies_\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"freezeOracle_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezeFee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unfreezeReceiver_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"unfreezePeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"operatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"validatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lstRestakerDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isValidatorOptedIn\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lstRestakerDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lstRestakerRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"operatorRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerLSTRestaker\",\"inputs\":[{\"name\":\"chosenValidators\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOperator\",\"inputs\":[{\"name\":\"operatorSignature\",\"type\":\"tuple\",\"internalType\":\"structISignatureUtils.SignatureWithSaltAndExpiry\",\"components\":[{\"name\":\"signature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerValidatorsByPodOwners\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[][]\",\"internalType\":\"bytes[][]\"},{\"name\":\"podOwners\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestLSTRestakerDeregistration\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestOperatorDeregistration\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestValidatorsDeregistration\",\"inputs\":[{\"name\":\"valPubKeys\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"restakeableStrategies\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setAVSDirectory\",\"inputs\":[{\"name\":\"avsDirectory_\",\"type\":\"address\",\"internalType\":\"contractIAVSDirectory\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegationManager\",\"inputs\":[{\"name\":\"delegationManager_\",\"type\":\"address\",\"internalType\":\"contractIDelegationManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEigenPodManager\",\"inputs\":[{\"name\":\"eigenPodManager_\",\"type\":\"address\",\"internalType\":\"contractIEigenPodManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFreezeOracle\",\"inputs\":[{\"name\":\"freezeOracle_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setLstRestakerDeregPeriodBlocks\",\"inputs\":[{\"name\":\"lstRestakerDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setOperatorDeregPeriodBlocks\",\"inputs\":[{\"name\":\"operatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRestakeableStrategies\",\"inputs\":[{\"name\":\"restakeableStrategies_\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStrategyManager\",\"inputs\":[{\"name\":\"strategyManager_\",\"type\":\"address\",\"internalType\":\"contractIStrategyManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeFee\",\"inputs\":[{\"name\":\"unfreezeFee_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezePeriodBlocks\",\"inputs\":[{\"name\":\"unfreezePeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setUnfreezeReceiver\",\"inputs\":[{\"name\":\"unfreezeReceiver_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setValidatorDeregPeriodBlocks\",\"inputs\":[{\"name\":\"validatorDeregPeriodBlocks_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unfreeze\",\"inputs\":[{\"name\":\"valPubKey\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"unfreezeFee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezePeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unfreezeReceiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadataURI\",\"inputs\":[{\"name\":\"metadataURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validatorDeregPeriodBlocks\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorRegistrations\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"podOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"freezeHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"deregRequestHeight\",\"type\":\"tuple\",\"internalType\":\"structEventHeightLib.EventHeight\",\"components\":[{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockHeight\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AVSDirectorySet\",\"inputs\":[{\"name\":\"avsDirectory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegationManagerSet\",\"inputs\":[{\"name\":\"delegationManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EigenPodManagerSet\",\"inputs\":[{\"name\":\"eigenPodManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FreezeOracleSet\",\"inputs\":[{\"name\":\"freezeOracle\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"lstRestakerDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregistered\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerDeregistrationRequested\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LSTRestakerRegistered\",\"inputs\":[{\"name\":\"chosenValidator\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"numChosen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lstRestaker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"operatorDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorDeregistrationRequested\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorRegistered\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RestakeableStrategiesSet\",\"inputs\":[{\"name\":\"restakeableStrategies\",\"type\":\"address[]\",\"indexed\":true,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StrategyManagerSet\",\"inputs\":[{\"name\":\"strategyManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezeFeeSet\",\"inputs\":[{\"name\":\"unfreezeFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezePeriodBlocksSet\",\"inputs\":[{\"name\":\"unfreezePeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnfreezeReceiverSet\",\"inputs\":[{\"name\":\"unfreezeReceiver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregPeriodBlocksSet\",\"inputs\":[{\"name\":\"validatorDeregPeriodBlocks\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorDeregistrationRequested\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorFrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorRegistered\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnfrozen\",\"inputs\":[{\"name\":\"validatorPubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"podOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// Mevcommitavsv3ABI is the input ABI used to generate the binding from.
// Deprecated: Use Mevcommitavsv3MetaData.ABI instead.
var Mevcommitavsv3ABI = Mevcommitavsv3MetaData.ABI

// Mevcommitavsv3 is an auto generated Go binding around an Ethereum contract.
type Mevcommitavsv3 struct {
	Mevcommitavsv3Caller     // Read-only binding to the contract
	Mevcommitavsv3Transactor // Write-only binding to the contract
	Mevcommitavsv3Filterer   // Log filterer for contract events
}

// Mevcommitavsv3Caller is an auto generated read-only Go binding around an Ethereum contract.
type Mevcommitavsv3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Mevcommitavsv3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Mevcommitavsv3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Mevcommitavsv3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Mevcommitavsv3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Mevcommitavsv3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Mevcommitavsv3Session struct {
	Contract     *Mevcommitavsv3   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Mevcommitavsv3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Mevcommitavsv3CallerSession struct {
	Contract *Mevcommitavsv3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// Mevcommitavsv3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Mevcommitavsv3TransactorSession struct {
	Contract     *Mevcommitavsv3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// Mevcommitavsv3Raw is an auto generated low-level Go binding around an Ethereum contract.
type Mevcommitavsv3Raw struct {
	Contract *Mevcommitavsv3 // Generic contract binding to access the raw methods on
}

// Mevcommitavsv3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Mevcommitavsv3CallerRaw struct {
	Contract *Mevcommitavsv3Caller // Generic read-only contract binding to access the raw methods on
}

// Mevcommitavsv3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Mevcommitavsv3TransactorRaw struct {
	Contract *Mevcommitavsv3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewMevcommitavsv3 creates a new instance of Mevcommitavsv3, bound to a specific deployed contract.
func NewMevcommitavsv3(address common.Address, backend bind.ContractBackend) (*Mevcommitavsv3, error) {
	contract, err := bindMevcommitavsv3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3{Mevcommitavsv3Caller: Mevcommitavsv3Caller{contract: contract}, Mevcommitavsv3Transactor: Mevcommitavsv3Transactor{contract: contract}, Mevcommitavsv3Filterer: Mevcommitavsv3Filterer{contract: contract}}, nil
}

// NewMevcommitavsv3Caller creates a new read-only instance of Mevcommitavsv3, bound to a specific deployed contract.
func NewMevcommitavsv3Caller(address common.Address, caller bind.ContractCaller) (*Mevcommitavsv3Caller, error) {
	contract, err := bindMevcommitavsv3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3Caller{contract: contract}, nil
}

// NewMevcommitavsv3Transactor creates a new write-only instance of Mevcommitavsv3, bound to a specific deployed contract.
func NewMevcommitavsv3Transactor(address common.Address, transactor bind.ContractTransactor) (*Mevcommitavsv3Transactor, error) {
	contract, err := bindMevcommitavsv3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3Transactor{contract: contract}, nil
}

// NewMevcommitavsv3Filterer creates a new log filterer instance of Mevcommitavsv3, bound to a specific deployed contract.
func NewMevcommitavsv3Filterer(address common.Address, filterer bind.ContractFilterer) (*Mevcommitavsv3Filterer, error) {
	contract, err := bindMevcommitavsv3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3Filterer{contract: contract}, nil
}

// bindMevcommitavsv3 binds a generic wrapper to an already deployed contract.
func bindMevcommitavsv3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Mevcommitavsv3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitavsv3 *Mevcommitavsv3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitavsv3.Contract.Mevcommitavsv3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitavsv3 *Mevcommitavsv3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Mevcommitavsv3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitavsv3 *Mevcommitavsv3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Mevcommitavsv3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mevcommitavsv3 *Mevcommitavsv3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mevcommitavsv3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitavsv3.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitavsv3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Mevcommitavsv3.Contract.UPGRADEINTERFACEVERSION(&_Mevcommitavsv3.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) AvsDirectory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "avsDirectory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) AvsDirectory() (common.Address, error) {
	return _Mevcommitavsv3.Contract.AvsDirectory(&_Mevcommitavsv3.CallOpts)
}

// AvsDirectory is a free data retrieval call binding the contract method 0x6b3aa72e.
//
// Solidity: function avsDirectory() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) AvsDirectory() (common.Address, error) {
	return _Mevcommitavsv3.Contract.AvsDirectory(&_Mevcommitavsv3.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) FreezeOracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "freezeOracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) FreezeOracle() (common.Address, error) {
	return _Mevcommitavsv3.Contract.FreezeOracle(&_Mevcommitavsv3.CallOpts)
}

// FreezeOracle is a free data retrieval call binding the contract method 0xaf91e0bf.
//
// Solidity: function freezeOracle() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) FreezeOracle() (common.Address, error) {
	return _Mevcommitavsv3.Contract.FreezeOracle(&_Mevcommitavsv3.CallOpts)
}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) GetLSTRestakerRegInfo(opts *bind.CallOpts, lstRestaker common.Address) (IMevCommitAVSV3LSTRestakerRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "getLSTRestakerRegInfo", lstRestaker)

	if err != nil {
		return *new(IMevCommitAVSV3LSTRestakerRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSV3LSTRestakerRegistrationInfo)).(*IMevCommitAVSV3LSTRestakerRegistrationInfo)

	return out0, err

}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Session) GetLSTRestakerRegInfo(lstRestaker common.Address) (IMevCommitAVSV3LSTRestakerRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetLSTRestakerRegInfo(&_Mevcommitavsv3.CallOpts, lstRestaker)
}

// GetLSTRestakerRegInfo is a free data retrieval call binding the contract method 0xeaeb9c88.
//
// Solidity: function getLSTRestakerRegInfo(address lstRestaker) view returns((bool,bytes[],uint256,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) GetLSTRestakerRegInfo(lstRestaker common.Address) (IMevCommitAVSV3LSTRestakerRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetLSTRestakerRegInfo(&_Mevcommitavsv3.CallOpts, lstRestaker)
}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) GetOperatorRegInfo(opts *bind.CallOpts, operator common.Address) (IMevCommitAVSV3OperatorRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "getOperatorRegInfo", operator)

	if err != nil {
		return *new(IMevCommitAVSV3OperatorRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSV3OperatorRegistrationInfo)).(*IMevCommitAVSV3OperatorRegistrationInfo)

	return out0, err

}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Session) GetOperatorRegInfo(operator common.Address) (IMevCommitAVSV3OperatorRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetOperatorRegInfo(&_Mevcommitavsv3.CallOpts, operator)
}

// GetOperatorRegInfo is a free data retrieval call binding the contract method 0x2c249e6c.
//
// Solidity: function getOperatorRegInfo(address operator) view returns((bool,(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) GetOperatorRegInfo(operator common.Address) (IMevCommitAVSV3OperatorRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetOperatorRegInfo(&_Mevcommitavsv3.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) GetOperatorRestakedStrategies(opts *bind.CallOpts, operator common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "getOperatorRestakedStrategies", operator)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3Session) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _Mevcommitavsv3.Contract.GetOperatorRestakedStrategies(&_Mevcommitavsv3.CallOpts, operator)
}

// GetOperatorRestakedStrategies is a free data retrieval call binding the contract method 0x33cfb7b7.
//
// Solidity: function getOperatorRestakedStrategies(address operator) view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) GetOperatorRestakedStrategies(operator common.Address) ([]common.Address, error) {
	return _Mevcommitavsv3.Contract.GetOperatorRestakedStrategies(&_Mevcommitavsv3.CallOpts, operator)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) GetRestakeableStrategies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "getRestakeableStrategies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3Session) GetRestakeableStrategies() ([]common.Address, error) {
	return _Mevcommitavsv3.Contract.GetRestakeableStrategies(&_Mevcommitavsv3.CallOpts)
}

// GetRestakeableStrategies is a free data retrieval call binding the contract method 0xe481af9d.
//
// Solidity: function getRestakeableStrategies() view returns(address[])
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) GetRestakeableStrategies() ([]common.Address, error) {
	return _Mevcommitavsv3.Contract.GetRestakeableStrategies(&_Mevcommitavsv3.CallOpts)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) GetValidatorRegInfo(opts *bind.CallOpts, valPubKey []byte) (IMevCommitAVSV3ValidatorRegistrationInfo, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "getValidatorRegInfo", valPubKey)

	if err != nil {
		return *new(IMevCommitAVSV3ValidatorRegistrationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IMevCommitAVSV3ValidatorRegistrationInfo)).(*IMevCommitAVSV3ValidatorRegistrationInfo)

	return out0, err

}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3Session) GetValidatorRegInfo(valPubKey []byte) (IMevCommitAVSV3ValidatorRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetValidatorRegInfo(&_Mevcommitavsv3.CallOpts, valPubKey)
}

// GetValidatorRegInfo is a free data retrieval call binding the contract method 0x972ac83c.
//
// Solidity: function getValidatorRegInfo(bytes valPubKey) view returns((bool,address,(bool,uint256),(bool,uint256)))
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) GetValidatorRegInfo(valPubKey []byte) (IMevCommitAVSV3ValidatorRegistrationInfo, error) {
	return _Mevcommitavsv3.Contract.GetValidatorRegInfo(&_Mevcommitavsv3.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) IsValidatorOptedIn(opts *bind.CallOpts, valPubKey []byte) (bool, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "isValidatorOptedIn", valPubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Mevcommitavsv3.Contract.IsValidatorOptedIn(&_Mevcommitavsv3.CallOpts, valPubKey)
}

// IsValidatorOptedIn is a free data retrieval call binding the contract method 0x470b690f.
//
// Solidity: function isValidatorOptedIn(bytes valPubKey) view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) IsValidatorOptedIn(valPubKey []byte) (bool, error) {
	return _Mevcommitavsv3.Contract.IsValidatorOptedIn(&_Mevcommitavsv3.CallOpts, valPubKey)
}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) LstRestakerDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "lstRestakerDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) LstRestakerDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.LstRestakerDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// LstRestakerDeregPeriodBlocks is a free data retrieval call binding the contract method 0xb0282b23.
//
// Solidity: function lstRestakerDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) LstRestakerDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.LstRestakerDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) LstRestakerRegistrations(opts *bind.CallOpts, arg0 common.Address) (struct {
	Exists             bool
	NumChosen          *big.Int
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "lstRestakerRegistrations", arg0)

	outstruct := new(struct {
		Exists             bool
		NumChosen          *big.Int
		DeregRequestHeight EventHeightLibEventHeight
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.NumChosen = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.DeregRequestHeight = *abi.ConvertType(out[2], new(EventHeightLibEventHeight)).(*EventHeightLibEventHeight)

	return *outstruct, err

}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) LstRestakerRegistrations(arg0 common.Address) (struct {
	Exists             bool
	NumChosen          *big.Int
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.LstRestakerRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// LstRestakerRegistrations is a free data retrieval call binding the contract method 0x25911aba.
//
// Solidity: function lstRestakerRegistrations(address ) view returns(bool exists, uint256 numChosen, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) LstRestakerRegistrations(arg0 common.Address) (struct {
	Exists             bool
	NumChosen          *big.Int
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.LstRestakerRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) OperatorDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "operatorDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) OperatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.OperatorDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// OperatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x14be85bd.
//
// Solidity: function operatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) OperatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.OperatorDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) OperatorRegistrations(opts *bind.CallOpts, arg0 common.Address) (struct {
	Exists             bool
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "operatorRegistrations", arg0)

	outstruct := new(struct {
		Exists             bool
		DeregRequestHeight EventHeightLibEventHeight
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.DeregRequestHeight = *abi.ConvertType(out[1], new(EventHeightLibEventHeight)).(*EventHeightLibEventHeight)

	return *outstruct, err

}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) OperatorRegistrations(arg0 common.Address) (struct {
	Exists             bool
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.OperatorRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// OperatorRegistrations is a free data retrieval call binding the contract method 0xfe07a836.
//
// Solidity: function operatorRegistrations(address ) view returns(bool exists, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) OperatorRegistrations(arg0 common.Address) (struct {
	Exists             bool
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.OperatorRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Owner() (common.Address, error) {
	return _Mevcommitavsv3.Contract.Owner(&_Mevcommitavsv3.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) Owner() (common.Address, error) {
	return _Mevcommitavsv3.Contract.Owner(&_Mevcommitavsv3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Paused() (bool, error) {
	return _Mevcommitavsv3.Contract.Paused(&_Mevcommitavsv3.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) Paused() (bool, error) {
	return _Mevcommitavsv3.Contract.Paused(&_Mevcommitavsv3.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) PendingOwner() (common.Address, error) {
	return _Mevcommitavsv3.Contract.PendingOwner(&_Mevcommitavsv3.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) PendingOwner() (common.Address, error) {
	return _Mevcommitavsv3.Contract.PendingOwner(&_Mevcommitavsv3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitavsv3.Contract.ProxiableUUID(&_Mevcommitavsv3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) ProxiableUUID() ([32]byte, error) {
	return _Mevcommitavsv3.Contract.ProxiableUUID(&_Mevcommitavsv3.CallOpts)
}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) RestakeableStrategies(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "restakeableStrategies", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RestakeableStrategies(arg0 *big.Int) (common.Address, error) {
	return _Mevcommitavsv3.Contract.RestakeableStrategies(&_Mevcommitavsv3.CallOpts, arg0)
}

// RestakeableStrategies is a free data retrieval call binding the contract method 0x94eef385.
//
// Solidity: function restakeableStrategies(uint256 ) view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) RestakeableStrategies(arg0 *big.Int) (common.Address, error) {
	return _Mevcommitavsv3.Contract.RestakeableStrategies(&_Mevcommitavsv3.CallOpts, arg0)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) UnfreezeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "unfreezeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UnfreezeFee() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.UnfreezeFee(&_Mevcommitavsv3.CallOpts)
}

// UnfreezeFee is a free data retrieval call binding the contract method 0x90d0c8c2.
//
// Solidity: function unfreezeFee() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) UnfreezeFee() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.UnfreezeFee(&_Mevcommitavsv3.CallOpts)
}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) UnfreezePeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "unfreezePeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UnfreezePeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.UnfreezePeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// UnfreezePeriodBlocks is a free data retrieval call binding the contract method 0x735ca5dd.
//
// Solidity: function unfreezePeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) UnfreezePeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.UnfreezePeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) UnfreezeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "unfreezeReceiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UnfreezeReceiver() (common.Address, error) {
	return _Mevcommitavsv3.Contract.UnfreezeReceiver(&_Mevcommitavsv3.CallOpts)
}

// UnfreezeReceiver is a free data retrieval call binding the contract method 0xc9207afb.
//
// Solidity: function unfreezeReceiver() view returns(address)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) UnfreezeReceiver() (common.Address, error) {
	return _Mevcommitavsv3.Contract.UnfreezeReceiver(&_Mevcommitavsv3.CallOpts)
}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) ValidatorDeregPeriodBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "validatorDeregPeriodBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) ValidatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.ValidatorDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// ValidatorDeregPeriodBlocks is a free data retrieval call binding the contract method 0x41a364a8.
//
// Solidity: function validatorDeregPeriodBlocks() view returns(uint256)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) ValidatorDeregPeriodBlocks() (*big.Int, error) {
	return _Mevcommitavsv3.Contract.ValidatorDeregPeriodBlocks(&_Mevcommitavsv3.CallOpts)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeHeight, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Caller) ValidatorRegistrations(opts *bind.CallOpts, arg0 []byte) (struct {
	Exists             bool
	PodOwner           common.Address
	FreezeHeight       EventHeightLibEventHeight
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	var out []interface{}
	err := _Mevcommitavsv3.contract.Call(opts, &out, "validatorRegistrations", arg0)

	outstruct := new(struct {
		Exists             bool
		PodOwner           common.Address
		FreezeHeight       EventHeightLibEventHeight
		DeregRequestHeight EventHeightLibEventHeight
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Exists = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PodOwner = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.FreezeHeight = *abi.ConvertType(out[2], new(EventHeightLibEventHeight)).(*EventHeightLibEventHeight)
	outstruct.DeregRequestHeight = *abi.ConvertType(out[3], new(EventHeightLibEventHeight)).(*EventHeightLibEventHeight)

	return *outstruct, err

}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeHeight, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3Session) ValidatorRegistrations(arg0 []byte) (struct {
	Exists             bool
	PodOwner           common.Address
	FreezeHeight       EventHeightLibEventHeight
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.ValidatorRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// ValidatorRegistrations is a free data retrieval call binding the contract method 0x8cdaf000.
//
// Solidity: function validatorRegistrations(bytes ) view returns(bool exists, address podOwner, (bool,uint256) freezeHeight, (bool,uint256) deregRequestHeight)
func (_Mevcommitavsv3 *Mevcommitavsv3CallerSession) ValidatorRegistrations(arg0 []byte) (struct {
	Exists             bool
	PodOwner           common.Address
	FreezeHeight       EventHeightLibEventHeight
	DeregRequestHeight EventHeightLibEventHeight
}, error) {
	return _Mevcommitavsv3.Contract.ValidatorRegistrations(&_Mevcommitavsv3.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.AcceptOwnership(&_Mevcommitavsv3.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.AcceptOwnership(&_Mevcommitavsv3.TransactOpts)
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) DeregisterLSTRestaker(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "deregisterLSTRestaker")
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) DeregisterLSTRestaker() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterLSTRestaker(&_Mevcommitavsv3.TransactOpts)
}

// DeregisterLSTRestaker is a paid mutator transaction binding the contract method 0x4ad29427.
//
// Solidity: function deregisterLSTRestaker() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) DeregisterLSTRestaker() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterLSTRestaker(&_Mevcommitavsv3.TransactOpts)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) DeregisterOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "deregisterOperator", operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) DeregisterOperator(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterOperator(&_Mevcommitavsv3.TransactOpts, operator)
}

// DeregisterOperator is a paid mutator transaction binding the contract method 0xd8cf98ca.
//
// Solidity: function deregisterOperator(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) DeregisterOperator(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterOperator(&_Mevcommitavsv3.TransactOpts, operator)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) DeregisterValidators(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "deregisterValidators", valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterValidators(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// DeregisterValidators is a paid mutator transaction binding the contract method 0xc6c6a657.
//
// Solidity: function deregisterValidators(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) DeregisterValidators(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.DeregisterValidators(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Freeze(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "freeze", valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Freeze(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// Freeze is a paid mutator transaction binding the contract method 0xa694d33f.
//
// Solidity: function freeze(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Freeze(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Freeze(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Initialize(opts *bind.TransactOpts, owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "initialize", owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Initialize(owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Initialize(&_Mevcommitavsv3.TransactOpts, owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Initialize is a paid mutator transaction binding the contract method 0x78a69cd4.
//
// Solidity: function initialize(address owner_, address delegationManager_, address eigenPodManager_, address strategyManager_, address avsDirectory_, address[] restakeableStrategies_, address freezeOracle_, uint256 unfreezeFee_, address unfreezeReceiver_, uint256 unfreezePeriodBlocks_, uint256 operatorDeregPeriodBlocks_, uint256 validatorDeregPeriodBlocks_, uint256 lstRestakerDeregPeriodBlocks_, string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Initialize(owner_ common.Address, delegationManager_ common.Address, eigenPodManager_ common.Address, strategyManager_ common.Address, avsDirectory_ common.Address, restakeableStrategies_ []common.Address, freezeOracle_ common.Address, unfreezeFee_ *big.Int, unfreezeReceiver_ common.Address, unfreezePeriodBlocks_ *big.Int, operatorDeregPeriodBlocks_ *big.Int, validatorDeregPeriodBlocks_ *big.Int, lstRestakerDeregPeriodBlocks_ *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Initialize(&_Mevcommitavsv3.TransactOpts, owner_, delegationManager_, eigenPodManager_, strategyManager_, avsDirectory_, restakeableStrategies_, freezeOracle_, unfreezeFee_, unfreezeReceiver_, unfreezePeriodBlocks_, operatorDeregPeriodBlocks_, validatorDeregPeriodBlocks_, lstRestakerDeregPeriodBlocks_, metadataURI_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Pause() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Pause(&_Mevcommitavsv3.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Pause() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Pause(&_Mevcommitavsv3.TransactOpts)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RegisterLSTRestaker(opts *bind.TransactOpts, chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "registerLSTRestaker", chosenValidators)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RegisterLSTRestaker(chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterLSTRestaker(&_Mevcommitavsv3.TransactOpts, chosenValidators)
}

// RegisterLSTRestaker is a paid mutator transaction binding the contract method 0xa807a70e.
//
// Solidity: function registerLSTRestaker(bytes[] chosenValidators) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RegisterLSTRestaker(chosenValidators [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterLSTRestaker(&_Mevcommitavsv3.TransactOpts, chosenValidators)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RegisterOperator(opts *bind.TransactOpts, operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "registerOperator", operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RegisterOperator(operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterOperator(&_Mevcommitavsv3.TransactOpts, operatorSignature)
}

// RegisterOperator is a paid mutator transaction binding the contract method 0x8317781d.
//
// Solidity: function registerOperator((bytes,bytes32,uint256) operatorSignature) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RegisterOperator(operatorSignature ISignatureUtilsSignatureWithSaltAndExpiry) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterOperator(&_Mevcommitavsv3.TransactOpts, operatorSignature)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RegisterValidatorsByPodOwners(opts *bind.TransactOpts, valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "registerValidatorsByPodOwners", valPubKeys, podOwners)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RegisterValidatorsByPodOwners(valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterValidatorsByPodOwners(&_Mevcommitavsv3.TransactOpts, valPubKeys, podOwners)
}

// RegisterValidatorsByPodOwners is a paid mutator transaction binding the contract method 0x86566f96.
//
// Solidity: function registerValidatorsByPodOwners(bytes[][] valPubKeys, address[] podOwners) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RegisterValidatorsByPodOwners(valPubKeys [][][]byte, podOwners []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RegisterValidatorsByPodOwners(&_Mevcommitavsv3.TransactOpts, valPubKeys, podOwners)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RenounceOwnership(&_Mevcommitavsv3.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RenounceOwnership(&_Mevcommitavsv3.TransactOpts)
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RequestLSTRestakerDeregistration(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "requestLSTRestakerDeregistration")
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RequestLSTRestakerDeregistration() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestLSTRestakerDeregistration(&_Mevcommitavsv3.TransactOpts)
}

// RequestLSTRestakerDeregistration is a paid mutator transaction binding the contract method 0x6e7a0e1f.
//
// Solidity: function requestLSTRestakerDeregistration() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RequestLSTRestakerDeregistration() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestLSTRestakerDeregistration(&_Mevcommitavsv3.TransactOpts)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RequestOperatorDeregistration(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "requestOperatorDeregistration", operator)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RequestOperatorDeregistration(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestOperatorDeregistration(&_Mevcommitavsv3.TransactOpts, operator)
}

// RequestOperatorDeregistration is a paid mutator transaction binding the contract method 0x95f8451b.
//
// Solidity: function requestOperatorDeregistration(address operator) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RequestOperatorDeregistration(operator common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestOperatorDeregistration(&_Mevcommitavsv3.TransactOpts, operator)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) RequestValidatorsDeregistration(opts *bind.TransactOpts, valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "requestValidatorsDeregistration", valPubKeys)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) RequestValidatorsDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestValidatorsDeregistration(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// RequestValidatorsDeregistration is a paid mutator transaction binding the contract method 0xeb35369b.
//
// Solidity: function requestValidatorsDeregistration(bytes[] valPubKeys) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) RequestValidatorsDeregistration(valPubKeys [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.RequestValidatorsDeregistration(&_Mevcommitavsv3.TransactOpts, valPubKeys)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetAVSDirectory(opts *bind.TransactOpts, avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setAVSDirectory", avsDirectory_)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetAVSDirectory(avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetAVSDirectory(&_Mevcommitavsv3.TransactOpts, avsDirectory_)
}

// SetAVSDirectory is a paid mutator transaction binding the contract method 0x862621ef.
//
// Solidity: function setAVSDirectory(address avsDirectory_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetAVSDirectory(avsDirectory_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetAVSDirectory(&_Mevcommitavsv3.TransactOpts, avsDirectory_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetDelegationManager(opts *bind.TransactOpts, delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setDelegationManager", delegationManager_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetDelegationManager(delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetDelegationManager(&_Mevcommitavsv3.TransactOpts, delegationManager_)
}

// SetDelegationManager is a paid mutator transaction binding the contract method 0x1a8d0de2.
//
// Solidity: function setDelegationManager(address delegationManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetDelegationManager(delegationManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetDelegationManager(&_Mevcommitavsv3.TransactOpts, delegationManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetEigenPodManager(opts *bind.TransactOpts, eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setEigenPodManager", eigenPodManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetEigenPodManager(eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetEigenPodManager(&_Mevcommitavsv3.TransactOpts, eigenPodManager_)
}

// SetEigenPodManager is a paid mutator transaction binding the contract method 0x3c2adfde.
//
// Solidity: function setEigenPodManager(address eigenPodManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetEigenPodManager(eigenPodManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetEigenPodManager(&_Mevcommitavsv3.TransactOpts, eigenPodManager_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetFreezeOracle(opts *bind.TransactOpts, freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setFreezeOracle", freezeOracle_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetFreezeOracle(freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetFreezeOracle(&_Mevcommitavsv3.TransactOpts, freezeOracle_)
}

// SetFreezeOracle is a paid mutator transaction binding the contract method 0x65a49071.
//
// Solidity: function setFreezeOracle(address freezeOracle_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetFreezeOracle(freezeOracle_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetFreezeOracle(&_Mevcommitavsv3.TransactOpts, freezeOracle_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetLstRestakerDeregPeriodBlocks(opts *bind.TransactOpts, lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setLstRestakerDeregPeriodBlocks", lstRestakerDeregPeriodBlocks_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetLstRestakerDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, lstRestakerDeregPeriodBlocks_)
}

// SetLstRestakerDeregPeriodBlocks is a paid mutator transaction binding the contract method 0x62f3dedb.
//
// Solidity: function setLstRestakerDeregPeriodBlocks(uint256 lstRestakerDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetLstRestakerDeregPeriodBlocks(lstRestakerDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetLstRestakerDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, lstRestakerDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetOperatorDeregPeriodBlocks(opts *bind.TransactOpts, operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setOperatorDeregPeriodBlocks", operatorDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetOperatorDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, operatorDeregPeriodBlocks_)
}

// SetOperatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xedfb9d0d.
//
// Solidity: function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetOperatorDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, operatorDeregPeriodBlocks_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetRestakeableStrategies(opts *bind.TransactOpts, restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setRestakeableStrategies", restakeableStrategies_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetRestakeableStrategies(restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetRestakeableStrategies(&_Mevcommitavsv3.TransactOpts, restakeableStrategies_)
}

// SetRestakeableStrategies is a paid mutator transaction binding the contract method 0xd871d570.
//
// Solidity: function setRestakeableStrategies(address[] restakeableStrategies_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetRestakeableStrategies(restakeableStrategies_ []common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetRestakeableStrategies(&_Mevcommitavsv3.TransactOpts, restakeableStrategies_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetStrategyManager(opts *bind.TransactOpts, strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setStrategyManager", strategyManager_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetStrategyManager(strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetStrategyManager(&_Mevcommitavsv3.TransactOpts, strategyManager_)
}

// SetStrategyManager is a paid mutator transaction binding the contract method 0x5c966646.
//
// Solidity: function setStrategyManager(address strategyManager_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetStrategyManager(strategyManager_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetStrategyManager(&_Mevcommitavsv3.TransactOpts, strategyManager_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetUnfreezeFee(opts *bind.TransactOpts, unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setUnfreezeFee", unfreezeFee_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetUnfreezeFee(unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezeFee(&_Mevcommitavsv3.TransactOpts, unfreezeFee_)
}

// SetUnfreezeFee is a paid mutator transaction binding the contract method 0x80e7751c.
//
// Solidity: function setUnfreezeFee(uint256 unfreezeFee_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetUnfreezeFee(unfreezeFee_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezeFee(&_Mevcommitavsv3.TransactOpts, unfreezeFee_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetUnfreezePeriodBlocks(opts *bind.TransactOpts, unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setUnfreezePeriodBlocks", unfreezePeriodBlocks_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetUnfreezePeriodBlocks(unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezePeriodBlocks(&_Mevcommitavsv3.TransactOpts, unfreezePeriodBlocks_)
}

// SetUnfreezePeriodBlocks is a paid mutator transaction binding the contract method 0x86c823e0.
//
// Solidity: function setUnfreezePeriodBlocks(uint256 unfreezePeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetUnfreezePeriodBlocks(unfreezePeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezePeriodBlocks(&_Mevcommitavsv3.TransactOpts, unfreezePeriodBlocks_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetUnfreezeReceiver(opts *bind.TransactOpts, unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setUnfreezeReceiver", unfreezeReceiver_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetUnfreezeReceiver(unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezeReceiver(&_Mevcommitavsv3.TransactOpts, unfreezeReceiver_)
}

// SetUnfreezeReceiver is a paid mutator transaction binding the contract method 0x7d0b802d.
//
// Solidity: function setUnfreezeReceiver(address unfreezeReceiver_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetUnfreezeReceiver(unfreezeReceiver_ common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetUnfreezeReceiver(&_Mevcommitavsv3.TransactOpts, unfreezeReceiver_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) SetValidatorDeregPeriodBlocks(opts *bind.TransactOpts, validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "setValidatorDeregPeriodBlocks", validatorDeregPeriodBlocks_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) SetValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetValidatorDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, validatorDeregPeriodBlocks_)
}

// SetValidatorDeregPeriodBlocks is a paid mutator transaction binding the contract method 0xb20bbf0a.
//
// Solidity: function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) SetValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_ *big.Int) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.SetValidatorDeregPeriodBlocks(&_Mevcommitavsv3.TransactOpts, validatorDeregPeriodBlocks_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.TransferOwnership(&_Mevcommitavsv3.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.TransferOwnership(&_Mevcommitavsv3.TransactOpts, newOwner)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKey) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Unfreeze(opts *bind.TransactOpts, valPubKey [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "unfreeze", valPubKey)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKey) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Unfreeze(valPubKey [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Unfreeze(&_Mevcommitavsv3.TransactOpts, valPubKey)
}

// Unfreeze is a paid mutator transaction binding the contract method 0xb764d33c.
//
// Solidity: function unfreeze(bytes[] valPubKey) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Unfreeze(valPubKey [][]byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Unfreeze(&_Mevcommitavsv3.TransactOpts, valPubKey)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Unpause() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Unpause(&_Mevcommitavsv3.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Unpause() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Unpause(&_Mevcommitavsv3.TransactOpts)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) UpdateMetadataURI(opts *bind.TransactOpts, metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "updateMetadataURI", metadataURI_)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UpdateMetadataURI(metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.UpdateMetadataURI(&_Mevcommitavsv3.TransactOpts, metadataURI_)
}

// UpdateMetadataURI is a paid mutator transaction binding the contract method 0x53fd3e81.
//
// Solidity: function updateMetadataURI(string metadataURI_) returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) UpdateMetadataURI(metadataURI_ string) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.UpdateMetadataURI(&_Mevcommitavsv3.TransactOpts, metadataURI_)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.UpgradeToAndCall(&_Mevcommitavsv3.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.UpgradeToAndCall(&_Mevcommitavsv3.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Fallback(&_Mevcommitavsv3.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Fallback(&_Mevcommitavsv3.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Transactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mevcommitavsv3.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3Session) Receive() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Receive(&_Mevcommitavsv3.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Mevcommitavsv3 *Mevcommitavsv3TransactorSession) Receive() (*types.Transaction, error) {
	return _Mevcommitavsv3.Contract.Receive(&_Mevcommitavsv3.TransactOpts)
}

// Mevcommitavsv3AVSDirectorySetIterator is returned from FilterAVSDirectorySet and is used to iterate over the raw logs and unpacked data for AVSDirectorySet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3AVSDirectorySetIterator struct {
	Event *Mevcommitavsv3AVSDirectorySet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3AVSDirectorySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3AVSDirectorySet)
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
		it.Event = new(Mevcommitavsv3AVSDirectorySet)
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
func (it *Mevcommitavsv3AVSDirectorySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3AVSDirectorySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3AVSDirectorySet represents a AVSDirectorySet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3AVSDirectorySet struct {
	AvsDirectory common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAVSDirectorySet is a free log retrieval operation binding the contract event 0x934223b20c24d569ff89796ae10a6997d43e2b3df0c3677fb6ca1f6e37ce344b.
//
// Solidity: event AVSDirectorySet(address indexed avsDirectory)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterAVSDirectorySet(opts *bind.FilterOpts, avsDirectory []common.Address) (*Mevcommitavsv3AVSDirectorySetIterator, error) {

	var avsDirectoryRule []interface{}
	for _, avsDirectoryItem := range avsDirectory {
		avsDirectoryRule = append(avsDirectoryRule, avsDirectoryItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "AVSDirectorySet", avsDirectoryRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3AVSDirectorySetIterator{contract: _Mevcommitavsv3.contract, event: "AVSDirectorySet", logs: logs, sub: sub}, nil
}

// WatchAVSDirectorySet is a free log subscription operation binding the contract event 0x934223b20c24d569ff89796ae10a6997d43e2b3df0c3677fb6ca1f6e37ce344b.
//
// Solidity: event AVSDirectorySet(address indexed avsDirectory)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchAVSDirectorySet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3AVSDirectorySet, avsDirectory []common.Address) (event.Subscription, error) {

	var avsDirectoryRule []interface{}
	for _, avsDirectoryItem := range avsDirectory {
		avsDirectoryRule = append(avsDirectoryRule, avsDirectoryItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "AVSDirectorySet", avsDirectoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3AVSDirectorySet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "AVSDirectorySet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseAVSDirectorySet(log types.Log) (*Mevcommitavsv3AVSDirectorySet, error) {
	event := new(Mevcommitavsv3AVSDirectorySet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "AVSDirectorySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3DelegationManagerSetIterator is returned from FilterDelegationManagerSet and is used to iterate over the raw logs and unpacked data for DelegationManagerSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3DelegationManagerSetIterator struct {
	Event *Mevcommitavsv3DelegationManagerSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3DelegationManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3DelegationManagerSet)
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
		it.Event = new(Mevcommitavsv3DelegationManagerSet)
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
func (it *Mevcommitavsv3DelegationManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3DelegationManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3DelegationManagerSet represents a DelegationManagerSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3DelegationManagerSet struct {
	DelegationManager common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterDelegationManagerSet is a free log retrieval operation binding the contract event 0x2296e6d8aebb5c81250fd381a114c2ec346fc44bc4582ba95cdcac0f09df6cd9.
//
// Solidity: event DelegationManagerSet(address indexed delegationManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterDelegationManagerSet(opts *bind.FilterOpts, delegationManager []common.Address) (*Mevcommitavsv3DelegationManagerSetIterator, error) {

	var delegationManagerRule []interface{}
	for _, delegationManagerItem := range delegationManager {
		delegationManagerRule = append(delegationManagerRule, delegationManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "DelegationManagerSet", delegationManagerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3DelegationManagerSetIterator{contract: _Mevcommitavsv3.contract, event: "DelegationManagerSet", logs: logs, sub: sub}, nil
}

// WatchDelegationManagerSet is a free log subscription operation binding the contract event 0x2296e6d8aebb5c81250fd381a114c2ec346fc44bc4582ba95cdcac0f09df6cd9.
//
// Solidity: event DelegationManagerSet(address indexed delegationManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchDelegationManagerSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3DelegationManagerSet, delegationManager []common.Address) (event.Subscription, error) {

	var delegationManagerRule []interface{}
	for _, delegationManagerItem := range delegationManager {
		delegationManagerRule = append(delegationManagerRule, delegationManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "DelegationManagerSet", delegationManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3DelegationManagerSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "DelegationManagerSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseDelegationManagerSet(log types.Log) (*Mevcommitavsv3DelegationManagerSet, error) {
	event := new(Mevcommitavsv3DelegationManagerSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "DelegationManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3EigenPodManagerSetIterator is returned from FilterEigenPodManagerSet and is used to iterate over the raw logs and unpacked data for EigenPodManagerSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3EigenPodManagerSetIterator struct {
	Event *Mevcommitavsv3EigenPodManagerSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3EigenPodManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3EigenPodManagerSet)
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
		it.Event = new(Mevcommitavsv3EigenPodManagerSet)
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
func (it *Mevcommitavsv3EigenPodManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3EigenPodManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3EigenPodManagerSet represents a EigenPodManagerSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3EigenPodManagerSet struct {
	EigenPodManager common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterEigenPodManagerSet is a free log retrieval operation binding the contract event 0x42070ca05aa3fae96a2bb90f36887ecc4894e2e33e748efeb2721962c11fd801.
//
// Solidity: event EigenPodManagerSet(address indexed eigenPodManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterEigenPodManagerSet(opts *bind.FilterOpts, eigenPodManager []common.Address) (*Mevcommitavsv3EigenPodManagerSetIterator, error) {

	var eigenPodManagerRule []interface{}
	for _, eigenPodManagerItem := range eigenPodManager {
		eigenPodManagerRule = append(eigenPodManagerRule, eigenPodManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "EigenPodManagerSet", eigenPodManagerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3EigenPodManagerSetIterator{contract: _Mevcommitavsv3.contract, event: "EigenPodManagerSet", logs: logs, sub: sub}, nil
}

// WatchEigenPodManagerSet is a free log subscription operation binding the contract event 0x42070ca05aa3fae96a2bb90f36887ecc4894e2e33e748efeb2721962c11fd801.
//
// Solidity: event EigenPodManagerSet(address indexed eigenPodManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchEigenPodManagerSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3EigenPodManagerSet, eigenPodManager []common.Address) (event.Subscription, error) {

	var eigenPodManagerRule []interface{}
	for _, eigenPodManagerItem := range eigenPodManager {
		eigenPodManagerRule = append(eigenPodManagerRule, eigenPodManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "EigenPodManagerSet", eigenPodManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3EigenPodManagerSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "EigenPodManagerSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseEigenPodManagerSet(log types.Log) (*Mevcommitavsv3EigenPodManagerSet, error) {
	event := new(Mevcommitavsv3EigenPodManagerSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "EigenPodManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3FreezeOracleSetIterator is returned from FilterFreezeOracleSet and is used to iterate over the raw logs and unpacked data for FreezeOracleSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3FreezeOracleSetIterator struct {
	Event *Mevcommitavsv3FreezeOracleSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3FreezeOracleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3FreezeOracleSet)
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
		it.Event = new(Mevcommitavsv3FreezeOracleSet)
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
func (it *Mevcommitavsv3FreezeOracleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3FreezeOracleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3FreezeOracleSet represents a FreezeOracleSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3FreezeOracleSet struct {
	FreezeOracle common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFreezeOracleSet is a free log retrieval operation binding the contract event 0xa33f3723c675820c785c70cde43f95aea5a4a0da3a5443a6cc129e14fcc9455a.
//
// Solidity: event FreezeOracleSet(address indexed freezeOracle)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterFreezeOracleSet(opts *bind.FilterOpts, freezeOracle []common.Address) (*Mevcommitavsv3FreezeOracleSetIterator, error) {

	var freezeOracleRule []interface{}
	for _, freezeOracleItem := range freezeOracle {
		freezeOracleRule = append(freezeOracleRule, freezeOracleItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "FreezeOracleSet", freezeOracleRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3FreezeOracleSetIterator{contract: _Mevcommitavsv3.contract, event: "FreezeOracleSet", logs: logs, sub: sub}, nil
}

// WatchFreezeOracleSet is a free log subscription operation binding the contract event 0xa33f3723c675820c785c70cde43f95aea5a4a0da3a5443a6cc129e14fcc9455a.
//
// Solidity: event FreezeOracleSet(address indexed freezeOracle)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchFreezeOracleSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3FreezeOracleSet, freezeOracle []common.Address) (event.Subscription, error) {

	var freezeOracleRule []interface{}
	for _, freezeOracleItem := range freezeOracle {
		freezeOracleRule = append(freezeOracleRule, freezeOracleItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "FreezeOracleSet", freezeOracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3FreezeOracleSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "FreezeOracleSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseFreezeOracleSet(log types.Log) (*Mevcommitavsv3FreezeOracleSet, error) {
	event := new(Mevcommitavsv3FreezeOracleSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "FreezeOracleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3InitializedIterator struct {
	Event *Mevcommitavsv3Initialized // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3Initialized)
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
		it.Event = new(Mevcommitavsv3Initialized)
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
func (it *Mevcommitavsv3InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3Initialized represents a Initialized event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterInitialized(opts *bind.FilterOpts) (*Mevcommitavsv3InitializedIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3InitializedIterator{contract: _Mevcommitavsv3.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3Initialized) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3Initialized)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseInitialized(log types.Log) (*Mevcommitavsv3Initialized, error) {
	event := new(Mevcommitavsv3Initialized)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator is returned from FilterLSTRestakerDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregPeriodBlocksSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator struct {
	Event *Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet)
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
		it.Event = new(Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet)
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
func (it *Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet represents a LSTRestakerDeregPeriodBlocksSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet struct {
	LstRestakerDeregPeriodBlocks *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0x7bd82fe806a0299d04c1dc6b928d934d0515dd3dfb8e6b0a0ca02267da5ec181.
//
// Solidity: event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterLSTRestakerDeregPeriodBlocksSet(opts *bind.FilterOpts) (*Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "LSTRestakerDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3LSTRestakerDeregPeriodBlocksSetIterator{contract: _Mevcommitavsv3.contract, event: "LSTRestakerDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0x7bd82fe806a0299d04c1dc6b928d934d0515dd3dfb8e6b0a0ca02267da5ec181.
//
// Solidity: event LSTRestakerDeregPeriodBlocksSet(uint256 lstRestakerDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchLSTRestakerDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "LSTRestakerDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregPeriodBlocksSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseLSTRestakerDeregPeriodBlocksSet(log types.Log) (*Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet, error) {
	event := new(Mevcommitavsv3LSTRestakerDeregPeriodBlocksSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3LSTRestakerDeregisteredIterator is returned from FilterLSTRestakerDeregistered and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregisteredIterator struct {
	Event *Mevcommitavsv3LSTRestakerDeregistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3LSTRestakerDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3LSTRestakerDeregistered)
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
		it.Event = new(Mevcommitavsv3LSTRestakerDeregistered)
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
func (it *Mevcommitavsv3LSTRestakerDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3LSTRestakerDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3LSTRestakerDeregistered represents a LSTRestakerDeregistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregistered struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregistered is a free log retrieval operation binding the contract event 0xaf3d14c12fe6a17d9cab68818354270a52de8de87f44c614c8cf7c35e96086fd.
//
// Solidity: event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterLSTRestakerDeregistered(opts *bind.FilterOpts, lstRestaker []common.Address) (*Mevcommitavsv3LSTRestakerDeregisteredIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "LSTRestakerDeregistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3LSTRestakerDeregisteredIterator{contract: _Mevcommitavsv3.contract, event: "LSTRestakerDeregistered", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregistered is a free log subscription operation binding the contract event 0xaf3d14c12fe6a17d9cab68818354270a52de8de87f44c614c8cf7c35e96086fd.
//
// Solidity: event LSTRestakerDeregistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchLSTRestakerDeregistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3LSTRestakerDeregistered, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "LSTRestakerDeregistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3LSTRestakerDeregistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseLSTRestakerDeregistered(log types.Log) (*Mevcommitavsv3LSTRestakerDeregistered, error) {
	event := new(Mevcommitavsv3LSTRestakerDeregistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator is returned from FilterLSTRestakerDeregistrationRequested and is used to iterate over the raw logs and unpacked data for LSTRestakerDeregistrationRequested events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator struct {
	Event *Mevcommitavsv3LSTRestakerDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3LSTRestakerDeregistrationRequested)
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
		it.Event = new(Mevcommitavsv3LSTRestakerDeregistrationRequested)
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
func (it *Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3LSTRestakerDeregistrationRequested represents a LSTRestakerDeregistrationRequested event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerDeregistrationRequested struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerDeregistrationRequested is a free log retrieval operation binding the contract event 0x6a7ec4f90fd0baec2ae0d35e6c4731b6c833f9f082a3fe8f5b30fe1be4af1c3b.
//
// Solidity: event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterLSTRestakerDeregistrationRequested(opts *bind.FilterOpts, lstRestaker []common.Address) (*Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "LSTRestakerDeregistrationRequested", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3LSTRestakerDeregistrationRequestedIterator{contract: _Mevcommitavsv3.contract, event: "LSTRestakerDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerDeregistrationRequested is a free log subscription operation binding the contract event 0x6a7ec4f90fd0baec2ae0d35e6c4731b6c833f9f082a3fe8f5b30fe1be4af1c3b.
//
// Solidity: event LSTRestakerDeregistrationRequested(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchLSTRestakerDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3LSTRestakerDeregistrationRequested, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "LSTRestakerDeregistrationRequested", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3LSTRestakerDeregistrationRequested)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregistrationRequested", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseLSTRestakerDeregistrationRequested(log types.Log) (*Mevcommitavsv3LSTRestakerDeregistrationRequested, error) {
	event := new(Mevcommitavsv3LSTRestakerDeregistrationRequested)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3LSTRestakerRegisteredIterator is returned from FilterLSTRestakerRegistered and is used to iterate over the raw logs and unpacked data for LSTRestakerRegistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerRegisteredIterator struct {
	Event *Mevcommitavsv3LSTRestakerRegistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3LSTRestakerRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3LSTRestakerRegistered)
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
		it.Event = new(Mevcommitavsv3LSTRestakerRegistered)
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
func (it *Mevcommitavsv3LSTRestakerRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3LSTRestakerRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3LSTRestakerRegistered represents a LSTRestakerRegistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3LSTRestakerRegistered struct {
	ChosenValidator []byte
	NumChosen       *big.Int
	LstRestaker     common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterLSTRestakerRegistered is a free log retrieval operation binding the contract event 0xdecaa7bc9a78e41b524d10750b5a52f0c8c1144cd5ab2991c9ca75ff380e011a.
//
// Solidity: event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterLSTRestakerRegistered(opts *bind.FilterOpts, lstRestaker []common.Address) (*Mevcommitavsv3LSTRestakerRegisteredIterator, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "LSTRestakerRegistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3LSTRestakerRegisteredIterator{contract: _Mevcommitavsv3.contract, event: "LSTRestakerRegistered", logs: logs, sub: sub}, nil
}

// WatchLSTRestakerRegistered is a free log subscription operation binding the contract event 0xdecaa7bc9a78e41b524d10750b5a52f0c8c1144cd5ab2991c9ca75ff380e011a.
//
// Solidity: event LSTRestakerRegistered(bytes chosenValidator, uint256 numChosen, address indexed lstRestaker)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchLSTRestakerRegistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3LSTRestakerRegistered, lstRestaker []common.Address) (event.Subscription, error) {

	var lstRestakerRule []interface{}
	for _, lstRestakerItem := range lstRestaker {
		lstRestakerRule = append(lstRestakerRule, lstRestakerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "LSTRestakerRegistered", lstRestakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3LSTRestakerRegistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerRegistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseLSTRestakerRegistered(log types.Log) (*Mevcommitavsv3LSTRestakerRegistered, error) {
	event := new(Mevcommitavsv3LSTRestakerRegistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "LSTRestakerRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator is returned from FilterOperatorDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for OperatorDeregPeriodBlocksSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator struct {
	Event *Mevcommitavsv3OperatorDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OperatorDeregPeriodBlocksSet)
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
		it.Event = new(Mevcommitavsv3OperatorDeregPeriodBlocksSet)
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
func (it *Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OperatorDeregPeriodBlocksSet represents a OperatorDeregPeriodBlocksSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregPeriodBlocksSet struct {
	OperatorDeregPeriodBlocks *big.Int
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0xfea621e39fd1186d690d8fa903a946ee52fee14c9c1f1c7295173b2e623b517e.
//
// Solidity: event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOperatorDeregPeriodBlocksSet(opts *bind.FilterOpts) (*Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OperatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OperatorDeregPeriodBlocksSetIterator{contract: _Mevcommitavsv3.contract, event: "OperatorDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0xfea621e39fd1186d690d8fa903a946ee52fee14c9c1f1c7295173b2e623b517e.
//
// Solidity: event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOperatorDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OperatorDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OperatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OperatorDeregPeriodBlocksSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregPeriodBlocksSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOperatorDeregPeriodBlocksSet(log types.Log) (*Mevcommitavsv3OperatorDeregPeriodBlocksSet, error) {
	event := new(Mevcommitavsv3OperatorDeregPeriodBlocksSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OperatorDeregisteredIterator is returned from FilterOperatorDeregistered and is used to iterate over the raw logs and unpacked data for OperatorDeregistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregisteredIterator struct {
	Event *Mevcommitavsv3OperatorDeregistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OperatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OperatorDeregistered)
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
		it.Event = new(Mevcommitavsv3OperatorDeregistered)
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
func (it *Mevcommitavsv3OperatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OperatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OperatorDeregistered represents a OperatorDeregistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistered is a free log retrieval operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOperatorDeregistered(opts *bind.FilterOpts, operator []common.Address) (*Mevcommitavsv3OperatorDeregisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OperatorDeregisteredIterator{contract: _Mevcommitavsv3.contract, event: "OperatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistered is a free log subscription operation binding the contract event 0x6dd4ca66565fb3dee8076c654634c6c4ad949022d809d0394308617d6791218d.
//
// Solidity: event OperatorDeregistered(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOperatorDeregistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OperatorDeregistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OperatorDeregistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OperatorDeregistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOperatorDeregistered(log types.Log) (*Mevcommitavsv3OperatorDeregistered, error) {
	event := new(Mevcommitavsv3OperatorDeregistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OperatorDeregistrationRequestedIterator is returned from FilterOperatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for OperatorDeregistrationRequested events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregistrationRequestedIterator struct {
	Event *Mevcommitavsv3OperatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OperatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OperatorDeregistrationRequested)
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
		it.Event = new(Mevcommitavsv3OperatorDeregistrationRequested)
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
func (it *Mevcommitavsv3OperatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OperatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OperatorDeregistrationRequested represents a OperatorDeregistrationRequested event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorDeregistrationRequested struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOperatorDeregistrationRequested(opts *bind.FilterOpts, operator []common.Address) (*Mevcommitavsv3OperatorDeregistrationRequestedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OperatorDeregistrationRequestedIterator{contract: _Mevcommitavsv3.contract, event: "OperatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchOperatorDeregistrationRequested is a free log subscription operation binding the contract event 0x4df522c04c21ddaed6db450a1c41907201a3daa6e80a58d12962062860a20d02.
//
// Solidity: event OperatorDeregistrationRequested(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOperatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OperatorDeregistrationRequested, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OperatorDeregistrationRequested", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OperatorDeregistrationRequested)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOperatorDeregistrationRequested(log types.Log) (*Mevcommitavsv3OperatorDeregistrationRequested, error) {
	event := new(Mevcommitavsv3OperatorDeregistrationRequested)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OperatorRegisteredIterator is returned from FilterOperatorRegistered and is used to iterate over the raw logs and unpacked data for OperatorRegistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorRegisteredIterator struct {
	Event *Mevcommitavsv3OperatorRegistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OperatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OperatorRegistered)
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
		it.Event = new(Mevcommitavsv3OperatorRegistered)
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
func (it *Mevcommitavsv3OperatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OperatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OperatorRegistered represents a OperatorRegistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OperatorRegistered struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRegistered is a free log retrieval operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOperatorRegistered(opts *bind.FilterOpts, operator []common.Address) (*Mevcommitavsv3OperatorRegisteredIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OperatorRegisteredIterator{contract: _Mevcommitavsv3.contract, event: "OperatorRegistered", logs: logs, sub: sub}, nil
}

// WatchOperatorRegistered is a free log subscription operation binding the contract event 0x4d0eb1f4bac8744fd2be119845e23b3befc88094b42bcda1204c65694a00f9e5.
//
// Solidity: event OperatorRegistered(address indexed operator)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOperatorRegistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OperatorRegistered, operator []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OperatorRegistered", operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OperatorRegistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOperatorRegistered(log types.Log) (*Mevcommitavsv3OperatorRegistered, error) {
	event := new(Mevcommitavsv3OperatorRegistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OperatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OwnershipTransferStartedIterator struct {
	Event *Mevcommitavsv3OwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OwnershipTransferStarted)
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
		it.Event = new(Mevcommitavsv3OwnershipTransferStarted)
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
func (it *Mevcommitavsv3OwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Mevcommitavsv3OwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OwnershipTransferStartedIterator{contract: _Mevcommitavsv3.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OwnershipTransferStarted)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOwnershipTransferStarted(log types.Log) (*Mevcommitavsv3OwnershipTransferStarted, error) {
	event := new(Mevcommitavsv3OwnershipTransferStarted)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OwnershipTransferredIterator struct {
	Event *Mevcommitavsv3OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3OwnershipTransferred)
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
		it.Event = new(Mevcommitavsv3OwnershipTransferred)
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
func (it *Mevcommitavsv3OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3OwnershipTransferred represents a OwnershipTransferred event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Mevcommitavsv3OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3OwnershipTransferredIterator{contract: _Mevcommitavsv3.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3OwnershipTransferred)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseOwnershipTransferred(log types.Log) (*Mevcommitavsv3OwnershipTransferred, error) {
	event := new(Mevcommitavsv3OwnershipTransferred)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3PausedIterator struct {
	Event *Mevcommitavsv3Paused // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3Paused)
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
		it.Event = new(Mevcommitavsv3Paused)
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
func (it *Mevcommitavsv3PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3Paused represents a Paused event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterPaused(opts *bind.FilterOpts) (*Mevcommitavsv3PausedIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3PausedIterator{contract: _Mevcommitavsv3.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3Paused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3Paused)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParsePaused(log types.Log) (*Mevcommitavsv3Paused, error) {
	event := new(Mevcommitavsv3Paused)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3RestakeableStrategiesSetIterator is returned from FilterRestakeableStrategiesSet and is used to iterate over the raw logs and unpacked data for RestakeableStrategiesSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3RestakeableStrategiesSetIterator struct {
	Event *Mevcommitavsv3RestakeableStrategiesSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3RestakeableStrategiesSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3RestakeableStrategiesSet)
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
		it.Event = new(Mevcommitavsv3RestakeableStrategiesSet)
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
func (it *Mevcommitavsv3RestakeableStrategiesSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3RestakeableStrategiesSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3RestakeableStrategiesSet represents a RestakeableStrategiesSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3RestakeableStrategiesSet struct {
	RestakeableStrategies []common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterRestakeableStrategiesSet is a free log retrieval operation binding the contract event 0xda4dd29046c55387af7520737a1e06033ae31f610dde3d0851458dffe13a0c0f.
//
// Solidity: event RestakeableStrategiesSet(address[] indexed restakeableStrategies)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterRestakeableStrategiesSet(opts *bind.FilterOpts, restakeableStrategies [][]common.Address) (*Mevcommitavsv3RestakeableStrategiesSetIterator, error) {

	var restakeableStrategiesRule []interface{}
	for _, restakeableStrategiesItem := range restakeableStrategies {
		restakeableStrategiesRule = append(restakeableStrategiesRule, restakeableStrategiesItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "RestakeableStrategiesSet", restakeableStrategiesRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3RestakeableStrategiesSetIterator{contract: _Mevcommitavsv3.contract, event: "RestakeableStrategiesSet", logs: logs, sub: sub}, nil
}

// WatchRestakeableStrategiesSet is a free log subscription operation binding the contract event 0xda4dd29046c55387af7520737a1e06033ae31f610dde3d0851458dffe13a0c0f.
//
// Solidity: event RestakeableStrategiesSet(address[] indexed restakeableStrategies)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchRestakeableStrategiesSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3RestakeableStrategiesSet, restakeableStrategies [][]common.Address) (event.Subscription, error) {

	var restakeableStrategiesRule []interface{}
	for _, restakeableStrategiesItem := range restakeableStrategies {
		restakeableStrategiesRule = append(restakeableStrategiesRule, restakeableStrategiesItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "RestakeableStrategiesSet", restakeableStrategiesRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3RestakeableStrategiesSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "RestakeableStrategiesSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseRestakeableStrategiesSet(log types.Log) (*Mevcommitavsv3RestakeableStrategiesSet, error) {
	event := new(Mevcommitavsv3RestakeableStrategiesSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "RestakeableStrategiesSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3StrategyManagerSetIterator is returned from FilterStrategyManagerSet and is used to iterate over the raw logs and unpacked data for StrategyManagerSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3StrategyManagerSetIterator struct {
	Event *Mevcommitavsv3StrategyManagerSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3StrategyManagerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3StrategyManagerSet)
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
		it.Event = new(Mevcommitavsv3StrategyManagerSet)
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
func (it *Mevcommitavsv3StrategyManagerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3StrategyManagerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3StrategyManagerSet represents a StrategyManagerSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3StrategyManagerSet struct {
	StrategyManager common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStrategyManagerSet is a free log retrieval operation binding the contract event 0x76fe640b216c20f563ab2d807634271e9d772e92c8a3752325cb2bc924e9e514.
//
// Solidity: event StrategyManagerSet(address indexed strategyManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterStrategyManagerSet(opts *bind.FilterOpts, strategyManager []common.Address) (*Mevcommitavsv3StrategyManagerSetIterator, error) {

	var strategyManagerRule []interface{}
	for _, strategyManagerItem := range strategyManager {
		strategyManagerRule = append(strategyManagerRule, strategyManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "StrategyManagerSet", strategyManagerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3StrategyManagerSetIterator{contract: _Mevcommitavsv3.contract, event: "StrategyManagerSet", logs: logs, sub: sub}, nil
}

// WatchStrategyManagerSet is a free log subscription operation binding the contract event 0x76fe640b216c20f563ab2d807634271e9d772e92c8a3752325cb2bc924e9e514.
//
// Solidity: event StrategyManagerSet(address indexed strategyManager)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchStrategyManagerSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3StrategyManagerSet, strategyManager []common.Address) (event.Subscription, error) {

	var strategyManagerRule []interface{}
	for _, strategyManagerItem := range strategyManager {
		strategyManagerRule = append(strategyManagerRule, strategyManagerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "StrategyManagerSet", strategyManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3StrategyManagerSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "StrategyManagerSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseStrategyManagerSet(log types.Log) (*Mevcommitavsv3StrategyManagerSet, error) {
	event := new(Mevcommitavsv3StrategyManagerSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "StrategyManagerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3UnfreezeFeeSetIterator is returned from FilterUnfreezeFeeSet and is used to iterate over the raw logs and unpacked data for UnfreezeFeeSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezeFeeSetIterator struct {
	Event *Mevcommitavsv3UnfreezeFeeSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3UnfreezeFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3UnfreezeFeeSet)
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
		it.Event = new(Mevcommitavsv3UnfreezeFeeSet)
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
func (it *Mevcommitavsv3UnfreezeFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3UnfreezeFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3UnfreezeFeeSet represents a UnfreezeFeeSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezeFeeSet struct {
	UnfreezeFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnfreezeFeeSet is a free log retrieval operation binding the contract event 0x6c0cf79356801bf6665a4c6cc85d35896ca003fe22c8c92c2b0e1d563b384c9d.
//
// Solidity: event UnfreezeFeeSet(uint256 unfreezeFee)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterUnfreezeFeeSet(opts *bind.FilterOpts) (*Mevcommitavsv3UnfreezeFeeSetIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "UnfreezeFeeSet")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3UnfreezeFeeSetIterator{contract: _Mevcommitavsv3.contract, event: "UnfreezeFeeSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezeFeeSet is a free log subscription operation binding the contract event 0x6c0cf79356801bf6665a4c6cc85d35896ca003fe22c8c92c2b0e1d563b384c9d.
//
// Solidity: event UnfreezeFeeSet(uint256 unfreezeFee)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchUnfreezeFeeSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3UnfreezeFeeSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "UnfreezeFeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3UnfreezeFeeSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezeFeeSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseUnfreezeFeeSet(log types.Log) (*Mevcommitavsv3UnfreezeFeeSet, error) {
	event := new(Mevcommitavsv3UnfreezeFeeSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezeFeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3UnfreezePeriodBlocksSetIterator is returned from FilterUnfreezePeriodBlocksSet and is used to iterate over the raw logs and unpacked data for UnfreezePeriodBlocksSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezePeriodBlocksSetIterator struct {
	Event *Mevcommitavsv3UnfreezePeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3UnfreezePeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3UnfreezePeriodBlocksSet)
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
		it.Event = new(Mevcommitavsv3UnfreezePeriodBlocksSet)
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
func (it *Mevcommitavsv3UnfreezePeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3UnfreezePeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3UnfreezePeriodBlocksSet represents a UnfreezePeriodBlocksSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezePeriodBlocksSet struct {
	UnfreezePeriodBlocks *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterUnfreezePeriodBlocksSet is a free log retrieval operation binding the contract event 0xdf42b0e9b907c6972f5ad1b757b3ab4e0eeca44f9b9dc7c8d97f2f40c1a042dd.
//
// Solidity: event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterUnfreezePeriodBlocksSet(opts *bind.FilterOpts) (*Mevcommitavsv3UnfreezePeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "UnfreezePeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3UnfreezePeriodBlocksSetIterator{contract: _Mevcommitavsv3.contract, event: "UnfreezePeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezePeriodBlocksSet is a free log subscription operation binding the contract event 0xdf42b0e9b907c6972f5ad1b757b3ab4e0eeca44f9b9dc7c8d97f2f40c1a042dd.
//
// Solidity: event UnfreezePeriodBlocksSet(uint256 unfreezePeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchUnfreezePeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3UnfreezePeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "UnfreezePeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3UnfreezePeriodBlocksSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezePeriodBlocksSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseUnfreezePeriodBlocksSet(log types.Log) (*Mevcommitavsv3UnfreezePeriodBlocksSet, error) {
	event := new(Mevcommitavsv3UnfreezePeriodBlocksSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezePeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3UnfreezeReceiverSetIterator is returned from FilterUnfreezeReceiverSet and is used to iterate over the raw logs and unpacked data for UnfreezeReceiverSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezeReceiverSetIterator struct {
	Event *Mevcommitavsv3UnfreezeReceiverSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3UnfreezeReceiverSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3UnfreezeReceiverSet)
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
		it.Event = new(Mevcommitavsv3UnfreezeReceiverSet)
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
func (it *Mevcommitavsv3UnfreezeReceiverSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3UnfreezeReceiverSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3UnfreezeReceiverSet represents a UnfreezeReceiverSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnfreezeReceiverSet struct {
	UnfreezeReceiver common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUnfreezeReceiverSet is a free log retrieval operation binding the contract event 0x7b992f6c09b22e43cad1f0f9a1d7e41949ac1c53a4eebb7f4b2374605049d2bf.
//
// Solidity: event UnfreezeReceiverSet(address indexed unfreezeReceiver)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterUnfreezeReceiverSet(opts *bind.FilterOpts, unfreezeReceiver []common.Address) (*Mevcommitavsv3UnfreezeReceiverSetIterator, error) {

	var unfreezeReceiverRule []interface{}
	for _, unfreezeReceiverItem := range unfreezeReceiver {
		unfreezeReceiverRule = append(unfreezeReceiverRule, unfreezeReceiverItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "UnfreezeReceiverSet", unfreezeReceiverRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3UnfreezeReceiverSetIterator{contract: _Mevcommitavsv3.contract, event: "UnfreezeReceiverSet", logs: logs, sub: sub}, nil
}

// WatchUnfreezeReceiverSet is a free log subscription operation binding the contract event 0x7b992f6c09b22e43cad1f0f9a1d7e41949ac1c53a4eebb7f4b2374605049d2bf.
//
// Solidity: event UnfreezeReceiverSet(address indexed unfreezeReceiver)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchUnfreezeReceiverSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3UnfreezeReceiverSet, unfreezeReceiver []common.Address) (event.Subscription, error) {

	var unfreezeReceiverRule []interface{}
	for _, unfreezeReceiverItem := range unfreezeReceiver {
		unfreezeReceiverRule = append(unfreezeReceiverRule, unfreezeReceiverItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "UnfreezeReceiverSet", unfreezeReceiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3UnfreezeReceiverSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezeReceiverSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseUnfreezeReceiverSet(log types.Log) (*Mevcommitavsv3UnfreezeReceiverSet, error) {
	event := new(Mevcommitavsv3UnfreezeReceiverSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "UnfreezeReceiverSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UnpausedIterator struct {
	Event *Mevcommitavsv3Unpaused // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3Unpaused)
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
		it.Event = new(Mevcommitavsv3Unpaused)
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
func (it *Mevcommitavsv3UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3Unpaused represents a Unpaused event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterUnpaused(opts *bind.FilterOpts) (*Mevcommitavsv3UnpausedIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3UnpausedIterator{contract: _Mevcommitavsv3.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3Unpaused) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3Unpaused)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseUnpaused(log types.Log) (*Mevcommitavsv3Unpaused, error) {
	event := new(Mevcommitavsv3Unpaused)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3UpgradedIterator struct {
	Event *Mevcommitavsv3Upgraded // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3Upgraded)
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
		it.Event = new(Mevcommitavsv3Upgraded)
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
func (it *Mevcommitavsv3UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3Upgraded represents a Upgraded event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*Mevcommitavsv3UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3UpgradedIterator{contract: _Mevcommitavsv3.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3Upgraded)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseUpgraded(log types.Log) (*Mevcommitavsv3Upgraded, error) {
	event := new(Mevcommitavsv3Upgraded)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator is returned from FilterValidatorDeregPeriodBlocksSet and is used to iterate over the raw logs and unpacked data for ValidatorDeregPeriodBlocksSet events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator struct {
	Event *Mevcommitavsv3ValidatorDeregPeriodBlocksSet // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorDeregPeriodBlocksSet)
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
		it.Event = new(Mevcommitavsv3ValidatorDeregPeriodBlocksSet)
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
func (it *Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorDeregPeriodBlocksSet represents a ValidatorDeregPeriodBlocksSet event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregPeriodBlocksSet struct {
	ValidatorDeregPeriodBlocks *big.Int
	Raw                        types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregPeriodBlocksSet is a free log retrieval operation binding the contract event 0x38755fd87ce522d770f5627d05f541406a9358d853b64009737f5c4d8913eae7.
//
// Solidity: event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorDeregPeriodBlocksSet(opts *bind.FilterOpts) (*Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator, error) {

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorDeregPeriodBlocksSetIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorDeregPeriodBlocksSet", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregPeriodBlocksSet is a free log subscription operation binding the contract event 0x38755fd87ce522d770f5627d05f541406a9358d853b64009737f5c4d8913eae7.
//
// Solidity: event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorDeregPeriodBlocksSet(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorDeregPeriodBlocksSet) (event.Subscription, error) {

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorDeregPeriodBlocksSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorDeregPeriodBlocksSet)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregPeriodBlocksSet", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorDeregPeriodBlocksSet(log types.Log) (*Mevcommitavsv3ValidatorDeregPeriodBlocksSet, error) {
	event := new(Mevcommitavsv3ValidatorDeregPeriodBlocksSet)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregPeriodBlocksSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorDeregisteredIterator is returned from FilterValidatorDeregistered and is used to iterate over the raw logs and unpacked data for ValidatorDeregistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregisteredIterator struct {
	Event *Mevcommitavsv3ValidatorDeregistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorDeregistered)
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
		it.Event = new(Mevcommitavsv3ValidatorDeregistered)
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
func (it *Mevcommitavsv3ValidatorDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorDeregistered represents a ValidatorDeregistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregistered struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistered is a free log retrieval operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorDeregistered(opts *bind.FilterOpts, podOwner []common.Address) (*Mevcommitavsv3ValidatorDeregisteredIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorDeregistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorDeregisteredIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorDeregistered", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistered is a free log subscription operation binding the contract event 0x10ec0bb1533e599e504516d6b49226d8a637ea19cbadfc6f7ff14a01bede3170.
//
// Solidity: event ValidatorDeregistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorDeregistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorDeregistered, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorDeregistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorDeregistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorDeregistered(log types.Log) (*Mevcommitavsv3ValidatorDeregistered, error) {
	event := new(Mevcommitavsv3ValidatorDeregistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorDeregistrationRequestedIterator is returned from FilterValidatorDeregistrationRequested and is used to iterate over the raw logs and unpacked data for ValidatorDeregistrationRequested events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregistrationRequestedIterator struct {
	Event *Mevcommitavsv3ValidatorDeregistrationRequested // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorDeregistrationRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorDeregistrationRequested)
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
		it.Event = new(Mevcommitavsv3ValidatorDeregistrationRequested)
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
func (it *Mevcommitavsv3ValidatorDeregistrationRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorDeregistrationRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorDeregistrationRequested represents a ValidatorDeregistrationRequested event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorDeregistrationRequested struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorDeregistrationRequested is a free log retrieval operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorDeregistrationRequested(opts *bind.FilterOpts, podOwner []common.Address) (*Mevcommitavsv3ValidatorDeregistrationRequestedIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorDeregistrationRequested", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorDeregistrationRequestedIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorDeregistrationRequested", logs: logs, sub: sub}, nil
}

// WatchValidatorDeregistrationRequested is a free log subscription operation binding the contract event 0x13b70fd48d462f71863cae24350d77b0dc4115a7e928b39dd0f0f60b701ffed3.
//
// Solidity: event ValidatorDeregistrationRequested(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorDeregistrationRequested(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorDeregistrationRequested, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorDeregistrationRequested", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorDeregistrationRequested)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorDeregistrationRequested(log types.Log) (*Mevcommitavsv3ValidatorDeregistrationRequested, error) {
	event := new(Mevcommitavsv3ValidatorDeregistrationRequested)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorDeregistrationRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorFrozenIterator is returned from FilterValidatorFrozen and is used to iterate over the raw logs and unpacked data for ValidatorFrozen events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorFrozenIterator struct {
	Event *Mevcommitavsv3ValidatorFrozen // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorFrozen)
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
		it.Event = new(Mevcommitavsv3ValidatorFrozen)
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
func (it *Mevcommitavsv3ValidatorFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorFrozen represents a ValidatorFrozen event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorFrozen struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorFrozen is a free log retrieval operation binding the contract event 0x5f565c1dd6cf6dc33dbdfd22c94b541af6ee1390251c8975b0c84106c58654bf.
//
// Solidity: event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorFrozen(opts *bind.FilterOpts, podOwner []common.Address) (*Mevcommitavsv3ValidatorFrozenIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorFrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorFrozenIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorFrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorFrozen is a free log subscription operation binding the contract event 0x5f565c1dd6cf6dc33dbdfd22c94b541af6ee1390251c8975b0c84106c58654bf.
//
// Solidity: event ValidatorFrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorFrozen(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorFrozen, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorFrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorFrozen)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorFrozen(log types.Log) (*Mevcommitavsv3ValidatorFrozen, error) {
	event := new(Mevcommitavsv3ValidatorFrozen)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorRegisteredIterator is returned from FilterValidatorRegistered and is used to iterate over the raw logs and unpacked data for ValidatorRegistered events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorRegisteredIterator struct {
	Event *Mevcommitavsv3ValidatorRegistered // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorRegistered)
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
		it.Event = new(Mevcommitavsv3ValidatorRegistered)
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
func (it *Mevcommitavsv3ValidatorRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorRegistered represents a ValidatorRegistered event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorRegistered struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorRegistered is a free log retrieval operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorRegistered(opts *bind.FilterOpts, podOwner []common.Address) (*Mevcommitavsv3ValidatorRegisteredIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorRegistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorRegisteredIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorRegistered", logs: logs, sub: sub}, nil
}

// WatchValidatorRegistered is a free log subscription operation binding the contract event 0x7cb7aef9bd2e5ee3f6073019691bb332fe3ef290465065aca1b9983f3dc66c56.
//
// Solidity: event ValidatorRegistered(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorRegistered(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorRegistered, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorRegistered", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorRegistered)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorRegistered(log types.Log) (*Mevcommitavsv3ValidatorRegistered, error) {
	event := new(Mevcommitavsv3ValidatorRegistered)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Mevcommitavsv3ValidatorUnfrozenIterator is returned from FilterValidatorUnfrozen and is used to iterate over the raw logs and unpacked data for ValidatorUnfrozen events raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorUnfrozenIterator struct {
	Event *Mevcommitavsv3ValidatorUnfrozen // Event containing the contract specifics and raw log

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
func (it *Mevcommitavsv3ValidatorUnfrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Mevcommitavsv3ValidatorUnfrozen)
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
		it.Event = new(Mevcommitavsv3ValidatorUnfrozen)
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
func (it *Mevcommitavsv3ValidatorUnfrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Mevcommitavsv3ValidatorUnfrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Mevcommitavsv3ValidatorUnfrozen represents a ValidatorUnfrozen event raised by the Mevcommitavsv3 contract.
type Mevcommitavsv3ValidatorUnfrozen struct {
	ValidatorPubKey []byte
	PodOwner        common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterValidatorUnfrozen is a free log retrieval operation binding the contract event 0x8b4f548363e5182887ee88395b0bacbd44ef955d9b9c1ace7aad0da43af0de40.
//
// Solidity: event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) FilterValidatorUnfrozen(opts *bind.FilterOpts, podOwner []common.Address) (*Mevcommitavsv3ValidatorUnfrozenIterator, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.FilterLogs(opts, "ValidatorUnfrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Mevcommitavsv3ValidatorUnfrozenIterator{contract: _Mevcommitavsv3.contract, event: "ValidatorUnfrozen", logs: logs, sub: sub}, nil
}

// WatchValidatorUnfrozen is a free log subscription operation binding the contract event 0x8b4f548363e5182887ee88395b0bacbd44ef955d9b9c1ace7aad0da43af0de40.
//
// Solidity: event ValidatorUnfrozen(bytes validatorPubKey, address indexed podOwner)
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) WatchValidatorUnfrozen(opts *bind.WatchOpts, sink chan<- *Mevcommitavsv3ValidatorUnfrozen, podOwner []common.Address) (event.Subscription, error) {

	var podOwnerRule []interface{}
	for _, podOwnerItem := range podOwner {
		podOwnerRule = append(podOwnerRule, podOwnerItem)
	}

	logs, sub, err := _Mevcommitavsv3.contract.WatchLogs(opts, "ValidatorUnfrozen", podOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Mevcommitavsv3ValidatorUnfrozen)
				if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
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
func (_Mevcommitavsv3 *Mevcommitavsv3Filterer) ParseValidatorUnfrozen(log types.Log) (*Mevcommitavsv3ValidatorUnfrozen, error) {
	event := new(Mevcommitavsv3ValidatorUnfrozen)
	if err := _Mevcommitavsv3.contract.UnpackLog(event, "ValidatorUnfrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
