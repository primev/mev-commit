// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package preconfmanager

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

// IPreconfManagerCommitmentParams is an auto generated low-level Go binding around an user-defined struct.
type IPreconfManagerCommitmentParams struct {
	TxnHash             string
	RevertingTxHashes   string
	BidAmt              *big.Int
	SlashAmt            *big.Int
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	BidHash             [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	ZkProof             []*big.Int
}

// IPreconfManagerOpenCommitmentParams is an auto generated low-level Go binding around an user-defined struct.
type IPreconfManagerOpenCommitmentParams struct {
	UnopenedCommitmentIndex [32]byte
	BidAmt                  *big.Int
	SlashAmt                *big.Int
	BlockNumber             uint64
	DecayStartTimeStamp     uint64
	DecayEndTimeStamp       uint64
	TxnHash                 string
	RevertingTxHashes       string
	BidSignature            []byte
	ZkProof                 []*big.Int
}

// IPreconfManagerOpenedCommitment is an auto generated low-level Go binding around an user-defined struct.
type IPreconfManagerOpenedCommitment struct {
	Bidder              common.Address
	IsSettled           bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	BidAmt              *big.Int
	SlashAmt            *big.Int
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	TxnHash             string
	RevertingTxHashes   string
}

// IPreconfManagerUnopenedCommitment is an auto generated low-level Go binding around an user-defined struct.
type IPreconfManagerUnopenedCommitment struct {
	IsOpened            bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}

// PreconfmanagerMetaData contains all meta data concerning the Preconfmanager contract.
var PreconfmanagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"EIP712_BID_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_COMMITMENT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"HEXCHARS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidderRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBidderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBlockTracker\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentDispatchWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentsCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"domainSeparatorBid\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"domainSeparatorPreconf\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBidHash\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.OpenCommitmentParams\",\"components\":[{\"name\":\"unopenedCommitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"zkProof\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.OpenedCommitment\",\"components\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOpenedCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.OpenedCommitment\",\"components\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPreConfHash\",\"inputs\":[{\"name\":\"_bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_zkProof\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTxnHashFromCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnopenedCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.UnopenedCommitment\",\"components\":[{\"name\":\"isOpened\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnopenedCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.UnopenedCommitment\",\"components\":[{\"name\":\"isOpened\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_providerRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_bidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_oracleContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blockTracker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_commitmentDispatchWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateReward\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateSlash\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"openCommitment\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.OpenCommitmentParams\",\"components\":[{\"name\":\"unopenedCommitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"zkProof\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"openedCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSettled\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracleContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processedTxnHashes\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProviderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeUnopenedCommitment\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unopenedCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"isOpened\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBidderRegistry\",\"inputs\":[{\"name\":\"newBidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBlockTracker\",\"inputs\":[{\"name\":\"newBlockTracker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCommitmentDispatchWindow\",\"inputs\":[{\"name\":\"newDispatchWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOracleContract\",\"inputs\":[{\"name\":\"newOracleContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProviderRegistry\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"verifyBid\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.OpenCommitmentParams\",\"components\":[{\"name\":\"unopenedCommitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"zkProof\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"messageDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recoveredAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"verifyPreConfCommitment\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIPreconfManager.CommitmentParams\",\"components\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"zkProof\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}]}],\"outputs\":[{\"name\":\"preConfHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"committerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"zkContextHash\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BidderRegistryUpdated\",\"inputs\":[{\"name\":\"newBidderRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BlockTrackerUpdated\",\"inputs\":[{\"name\":\"newBlockTracker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CommitmentDispatchWindowUpdated\",\"inputs\":[{\"name\":\"newDispatchWindow\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OpenedCommitmentStored\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"committer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"slashAmt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleContractUpdated\",\"inputs\":[{\"name\":\"newOracleContract\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProviderRegistryUpdated\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignatureVerified\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"bidAmt\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UnopenedCommitmentStored\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"committer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"BN128AddFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BN128MulFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CommitmentAlreadyOpened\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CommitmentAlreadySettled\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCommitmentDigest\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"computedDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidDecayTime\",\"inputs\":[{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"endTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidDispatchTimestamp\",\"inputs\":[{\"name\":\"minTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"InvalidFallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReceive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ProviderZKProofInvalid\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotCommitter\",\"inputs\":[{\"name\":\"expected\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"actual\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SenderIsNotOracleContract\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracleContract\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TxnHashAlreadyProcessed\",\"inputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidderAddress\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedOpenCommitment\",\"inputs\":[{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnopenedCommitmentAlreadyExists\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"WinnerIsNotCommitter\",\"inputs\":[{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"winner\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// PreconfmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use PreconfmanagerMetaData.ABI instead.
var PreconfmanagerABI = PreconfmanagerMetaData.ABI

// Preconfmanager is an auto generated Go binding around an Ethereum contract.
type Preconfmanager struct {
	PreconfmanagerCaller     // Read-only binding to the contract
	PreconfmanagerTransactor // Write-only binding to the contract
	PreconfmanagerFilterer   // Log filterer for contract events
}

// PreconfmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreconfmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreconfmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreconfmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreconfmanagerSession struct {
	Contract     *Preconfmanager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PreconfmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreconfmanagerCallerSession struct {
	Contract *PreconfmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PreconfmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreconfmanagerTransactorSession struct {
	Contract     *PreconfmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PreconfmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreconfmanagerRaw struct {
	Contract *Preconfmanager // Generic contract binding to access the raw methods on
}

// PreconfmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreconfmanagerCallerRaw struct {
	Contract *PreconfmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// PreconfmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreconfmanagerTransactorRaw struct {
	Contract *PreconfmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreconfmanager creates a new instance of Preconfmanager, bound to a specific deployed contract.
func NewPreconfmanager(address common.Address, backend bind.ContractBackend) (*Preconfmanager, error) {
	contract, err := bindPreconfmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Preconfmanager{PreconfmanagerCaller: PreconfmanagerCaller{contract: contract}, PreconfmanagerTransactor: PreconfmanagerTransactor{contract: contract}, PreconfmanagerFilterer: PreconfmanagerFilterer{contract: contract}}, nil
}

// NewPreconfmanagerCaller creates a new read-only instance of Preconfmanager, bound to a specific deployed contract.
func NewPreconfmanagerCaller(address common.Address, caller bind.ContractCaller) (*PreconfmanagerCaller, error) {
	contract, err := bindPreconfmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerCaller{contract: contract}, nil
}

// NewPreconfmanagerTransactor creates a new write-only instance of Preconfmanager, bound to a specific deployed contract.
func NewPreconfmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*PreconfmanagerTransactor, error) {
	contract, err := bindPreconfmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerTransactor{contract: contract}, nil
}

// NewPreconfmanagerFilterer creates a new log filterer instance of Preconfmanager, bound to a specific deployed contract.
func NewPreconfmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*PreconfmanagerFilterer, error) {
	contract, err := bindPreconfmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerFilterer{contract: contract}, nil
}

// bindPreconfmanager binds a generic wrapper to an already deployed contract.
func bindPreconfmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreconfmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfmanager *PreconfmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfmanager.Contract.PreconfmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfmanager *PreconfmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.Contract.PreconfmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfmanager *PreconfmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfmanager.Contract.PreconfmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfmanager *PreconfmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfmanager *PreconfmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfmanager *PreconfmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfmanager.Contract.contract.Transact(opts, method, params...)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) EIP712BIDTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "EIP712_BID_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfmanager.Contract.EIP712BIDTYPEHASH(&_Preconfmanager.CallOpts)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfmanager.Contract.EIP712BIDTYPEHASH(&_Preconfmanager.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) EIP712COMMITMENTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "EIP712_COMMITMENT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfmanager.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfmanager.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfmanager.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfmanager.CallOpts)
}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfmanager *PreconfmanagerCaller) HEXCHARS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "HEXCHARS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfmanager *PreconfmanagerSession) HEXCHARS() ([]byte, error) {
	return _Preconfmanager.Contract.HEXCHARS(&_Preconfmanager.CallOpts)
}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfmanager *PreconfmanagerCallerSession) HEXCHARS() ([]byte, error) {
	return _Preconfmanager.Contract.HEXCHARS(&_Preconfmanager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfmanager *PreconfmanagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfmanager *PreconfmanagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Preconfmanager.Contract.UPGRADEINTERFACEVERSION(&_Preconfmanager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfmanager *PreconfmanagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Preconfmanager.Contract.UPGRADEINTERFACEVERSION(&_Preconfmanager.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) BidderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "bidderRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) BidderRegistry() (common.Address, error) {
	return _Preconfmanager.Contract.BidderRegistry(&_Preconfmanager.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) BidderRegistry() (common.Address, error) {
	return _Preconfmanager.Contract.BidderRegistry(&_Preconfmanager.CallOpts)
}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) BlockTracker(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "blockTracker")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) BlockTracker() (common.Address, error) {
	return _Preconfmanager.Contract.BlockTracker(&_Preconfmanager.CallOpts)
}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) BlockTracker() (common.Address, error) {
	return _Preconfmanager.Contract.BlockTracker(&_Preconfmanager.CallOpts)
}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfmanager *PreconfmanagerCaller) CommitmentDispatchWindow(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "commitmentDispatchWindow")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfmanager *PreconfmanagerSession) CommitmentDispatchWindow() (uint64, error) {
	return _Preconfmanager.Contract.CommitmentDispatchWindow(&_Preconfmanager.CallOpts)
}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfmanager *PreconfmanagerCallerSession) CommitmentDispatchWindow() (uint64, error) {
	return _Preconfmanager.Contract.CommitmentDispatchWindow(&_Preconfmanager.CallOpts)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfmanager *PreconfmanagerCaller) CommitmentsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "commitmentsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfmanager *PreconfmanagerSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfmanager.Contract.CommitmentsCount(&_Preconfmanager.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfmanager *PreconfmanagerCallerSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfmanager.Contract.CommitmentsCount(&_Preconfmanager.CallOpts, arg0)
}

// DomainSeparatorBid is a free data retrieval call binding the contract method 0x3c281734.
//
// Solidity: function domainSeparatorBid() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) DomainSeparatorBid(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "domainSeparatorBid")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparatorBid is a free data retrieval call binding the contract method 0x3c281734.
//
// Solidity: function domainSeparatorBid() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) DomainSeparatorBid() ([32]byte, error) {
	return _Preconfmanager.Contract.DomainSeparatorBid(&_Preconfmanager.CallOpts)
}

// DomainSeparatorBid is a free data retrieval call binding the contract method 0x3c281734.
//
// Solidity: function domainSeparatorBid() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) DomainSeparatorBid() ([32]byte, error) {
	return _Preconfmanager.Contract.DomainSeparatorBid(&_Preconfmanager.CallOpts)
}

// DomainSeparatorPreconf is a free data retrieval call binding the contract method 0x3b8c3072.
//
// Solidity: function domainSeparatorPreconf() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) DomainSeparatorPreconf(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "domainSeparatorPreconf")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparatorPreconf is a free data retrieval call binding the contract method 0x3b8c3072.
//
// Solidity: function domainSeparatorPreconf() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) DomainSeparatorPreconf() ([32]byte, error) {
	return _Preconfmanager.Contract.DomainSeparatorPreconf(&_Preconfmanager.CallOpts)
}

// DomainSeparatorPreconf is a free data retrieval call binding the contract method 0x3b8c3072.
//
// Solidity: function domainSeparatorPreconf() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) DomainSeparatorPreconf() ([32]byte, error) {
	return _Preconfmanager.Contract.DomainSeparatorPreconf(&_Preconfmanager.CallOpts)
}

// GetBidHash is a free data retrieval call binding the contract method 0xc9e66dc5.
//
// Solidity: function getBidHash((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) GetBidHash(opts *bind.CallOpts, params IPreconfManagerOpenCommitmentParams) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getBidHash", params)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBidHash is a free data retrieval call binding the contract method 0xc9e66dc5.
//
// Solidity: function getBidHash((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) GetBidHash(params IPreconfManagerOpenCommitmentParams) ([32]byte, error) {
	return _Preconfmanager.Contract.GetBidHash(&_Preconfmanager.CallOpts, params)
}

// GetBidHash is a free data retrieval call binding the contract method 0xc9e66dc5.
//
// Solidity: function getBidHash((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) GetBidHash(params IPreconfManagerOpenCommitmentParams) ([32]byte, error) {
	return _Preconfmanager.Contract.GetBidHash(&_Preconfmanager.CallOpts, params)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string))
func (_Preconfmanager *PreconfmanagerCaller) GetCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (IPreconfManagerOpenedCommitment, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getCommitment", commitmentIndex)

	if err != nil {
		return *new(IPreconfManagerOpenedCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(IPreconfManagerOpenedCommitment)).(*IPreconfManagerOpenedCommitment)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string))
func (_Preconfmanager *PreconfmanagerSession) GetCommitment(commitmentIndex [32]byte) (IPreconfManagerOpenedCommitment, error) {
	return _Preconfmanager.Contract.GetCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string))
func (_Preconfmanager *PreconfmanagerCallerSession) GetCommitment(commitmentIndex [32]byte) (IPreconfManagerOpenedCommitment, error) {
	return _Preconfmanager.Contract.GetCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetOpenedCommitmentIndex is a free data retrieval call binding the contract method 0xe6158c02.
//
// Solidity: function getOpenedCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) GetOpenedCommitmentIndex(opts *bind.CallOpts, commitment IPreconfManagerOpenedCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getOpenedCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetOpenedCommitmentIndex is a free data retrieval call binding the contract method 0xe6158c02.
//
// Solidity: function getOpenedCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) GetOpenedCommitmentIndex(commitment IPreconfManagerOpenedCommitment) ([32]byte, error) {
	return _Preconfmanager.Contract.GetOpenedCommitmentIndex(&_Preconfmanager.CallOpts, commitment)
}

// GetOpenedCommitmentIndex is a free data retrieval call binding the contract method 0xe6158c02.
//
// Solidity: function getOpenedCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,uint256,bytes32,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) GetOpenedCommitmentIndex(commitment IPreconfManagerOpenedCommitment) ([32]byte, error) {
	return _Preconfmanager.Contract.GetOpenedCommitmentIndex(&_Preconfmanager.CallOpts, commitment)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xde44a719.
//
// Solidity: function getPreConfHash(bytes32 _bidHash, bytes _bidSignature, uint256[] _zkProof) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) GetPreConfHash(opts *bind.CallOpts, _bidHash [32]byte, _bidSignature []byte, _zkProof []*big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getPreConfHash", _bidHash, _bidSignature, _zkProof)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPreConfHash is a free data retrieval call binding the contract method 0xde44a719.
//
// Solidity: function getPreConfHash(bytes32 _bidHash, bytes _bidSignature, uint256[] _zkProof) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) GetPreConfHash(_bidHash [32]byte, _bidSignature []byte, _zkProof []*big.Int) ([32]byte, error) {
	return _Preconfmanager.Contract.GetPreConfHash(&_Preconfmanager.CallOpts, _bidHash, _bidSignature, _zkProof)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xde44a719.
//
// Solidity: function getPreConfHash(bytes32 _bidHash, bytes _bidSignature, uint256[] _zkProof) view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) GetPreConfHash(_bidHash [32]byte, _bidSignature []byte, _zkProof []*big.Int) ([32]byte, error) {
	return _Preconfmanager.Contract.GetPreConfHash(&_Preconfmanager.CallOpts, _bidHash, _bidSignature, _zkProof)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfmanager *PreconfmanagerCaller) GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getTxnHashFromCommitment", commitmentIndex)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfmanager *PreconfmanagerSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfmanager.Contract.GetTxnHashFromCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfmanager *PreconfmanagerCallerSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfmanager.Contract.GetTxnHashFromCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetUnopenedCommitment is a free data retrieval call binding the contract method 0x86a768e6.
//
// Solidity: function getUnopenedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfmanager *PreconfmanagerCaller) GetUnopenedCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (IPreconfManagerUnopenedCommitment, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getUnopenedCommitment", commitmentIndex)

	if err != nil {
		return *new(IPreconfManagerUnopenedCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(IPreconfManagerUnopenedCommitment)).(*IPreconfManagerUnopenedCommitment)

	return out0, err

}

// GetUnopenedCommitment is a free data retrieval call binding the contract method 0x86a768e6.
//
// Solidity: function getUnopenedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfmanager *PreconfmanagerSession) GetUnopenedCommitment(commitmentIndex [32]byte) (IPreconfManagerUnopenedCommitment, error) {
	return _Preconfmanager.Contract.GetUnopenedCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetUnopenedCommitment is a free data retrieval call binding the contract method 0x86a768e6.
//
// Solidity: function getUnopenedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfmanager *PreconfmanagerCallerSession) GetUnopenedCommitment(commitmentIndex [32]byte) (IPreconfManagerUnopenedCommitment, error) {
	return _Preconfmanager.Contract.GetUnopenedCommitment(&_Preconfmanager.CallOpts, commitmentIndex)
}

// GetUnopenedCommitmentIndex is a free data retrieval call binding the contract method 0xf9b1349f.
//
// Solidity: function getUnopenedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) GetUnopenedCommitmentIndex(opts *bind.CallOpts, commitment IPreconfManagerUnopenedCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "getUnopenedCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetUnopenedCommitmentIndex is a free data retrieval call binding the contract method 0xf9b1349f.
//
// Solidity: function getUnopenedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) GetUnopenedCommitmentIndex(commitment IPreconfManagerUnopenedCommitment) ([32]byte, error) {
	return _Preconfmanager.Contract.GetUnopenedCommitmentIndex(&_Preconfmanager.CallOpts, commitment)
}

// GetUnopenedCommitmentIndex is a free data retrieval call binding the contract method 0xf9b1349f.
//
// Solidity: function getUnopenedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) GetUnopenedCommitmentIndex(commitment IPreconfManagerUnopenedCommitment) ([32]byte, error) {
	return _Preconfmanager.Contract.GetUnopenedCommitmentIndex(&_Preconfmanager.CallOpts, commitment)
}

// OpenedCommitments is a free data retrieval call binding the contract method 0x5e0c7e10.
//
// Solidity: function openedCommitments(bytes32 ) view returns(address bidder, bool isSettled, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bidAmt, uint256 slashAmt, bytes32 commitmentDigest, bytes commitmentSignature, string txnHash, string revertingTxHashes)
func (_Preconfmanager *PreconfmanagerCaller) OpenedCommitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Bidder              common.Address
	IsSettled           bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	BidAmt              *big.Int
	SlashAmt            *big.Int
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "openedCommitments", arg0)

	outstruct := new(struct {
		Bidder              common.Address
		IsSettled           bool
		BlockNumber         uint64
		DecayStartTimeStamp uint64
		DecayEndTimeStamp   uint64
		DispatchTimestamp   uint64
		Committer           common.Address
		BidAmt              *big.Int
		SlashAmt            *big.Int
		CommitmentDigest    [32]byte
		CommitmentSignature []byte
		TxnHash             string
		RevertingTxHashes   string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Bidder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.IsSettled = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.BlockNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.DecayStartTimeStamp = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.DecayEndTimeStamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.DispatchTimestamp = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.Committer = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.BidAmt = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.SlashAmt = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.CommitmentDigest = *abi.ConvertType(out[9], new([32]byte)).(*[32]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[10], new([]byte)).(*[]byte)
	outstruct.TxnHash = *abi.ConvertType(out[11], new(string)).(*string)
	outstruct.RevertingTxHashes = *abi.ConvertType(out[12], new(string)).(*string)

	return *outstruct, err

}

// OpenedCommitments is a free data retrieval call binding the contract method 0x5e0c7e10.
//
// Solidity: function openedCommitments(bytes32 ) view returns(address bidder, bool isSettled, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bidAmt, uint256 slashAmt, bytes32 commitmentDigest, bytes commitmentSignature, string txnHash, string revertingTxHashes)
func (_Preconfmanager *PreconfmanagerSession) OpenedCommitments(arg0 [32]byte) (struct {
	Bidder              common.Address
	IsSettled           bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	BidAmt              *big.Int
	SlashAmt            *big.Int
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	return _Preconfmanager.Contract.OpenedCommitments(&_Preconfmanager.CallOpts, arg0)
}

// OpenedCommitments is a free data retrieval call binding the contract method 0x5e0c7e10.
//
// Solidity: function openedCommitments(bytes32 ) view returns(address bidder, bool isSettled, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bidAmt, uint256 slashAmt, bytes32 commitmentDigest, bytes commitmentSignature, string txnHash, string revertingTxHashes)
func (_Preconfmanager *PreconfmanagerCallerSession) OpenedCommitments(arg0 [32]byte) (struct {
	Bidder              common.Address
	IsSettled           bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	BidAmt              *big.Int
	SlashAmt            *big.Int
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	return _Preconfmanager.Contract.OpenedCommitments(&_Preconfmanager.CallOpts, arg0)
}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) OracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "oracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) OracleContract() (common.Address, error) {
	return _Preconfmanager.Contract.OracleContract(&_Preconfmanager.CallOpts)
}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) OracleContract() (common.Address, error) {
	return _Preconfmanager.Contract.OracleContract(&_Preconfmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) Owner() (common.Address, error) {
	return _Preconfmanager.Contract.Owner(&_Preconfmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) Owner() (common.Address, error) {
	return _Preconfmanager.Contract.Owner(&_Preconfmanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Preconfmanager *PreconfmanagerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Preconfmanager *PreconfmanagerSession) Paused() (bool, error) {
	return _Preconfmanager.Contract.Paused(&_Preconfmanager.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Preconfmanager *PreconfmanagerCallerSession) Paused() (bool, error) {
	return _Preconfmanager.Contract.Paused(&_Preconfmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) PendingOwner() (common.Address, error) {
	return _Preconfmanager.Contract.PendingOwner(&_Preconfmanager.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) PendingOwner() (common.Address, error) {
	return _Preconfmanager.Contract.PendingOwner(&_Preconfmanager.CallOpts)
}

// ProcessedTxnHashes is a free data retrieval call binding the contract method 0x58160b50.
//
// Solidity: function processedTxnHashes(bytes32 ) view returns(bool)
func (_Preconfmanager *PreconfmanagerCaller) ProcessedTxnHashes(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "processedTxnHashes", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProcessedTxnHashes is a free data retrieval call binding the contract method 0x58160b50.
//
// Solidity: function processedTxnHashes(bytes32 ) view returns(bool)
func (_Preconfmanager *PreconfmanagerSession) ProcessedTxnHashes(arg0 [32]byte) (bool, error) {
	return _Preconfmanager.Contract.ProcessedTxnHashes(&_Preconfmanager.CallOpts, arg0)
}

// ProcessedTxnHashes is a free data retrieval call binding the contract method 0x58160b50.
//
// Solidity: function processedTxnHashes(bytes32 ) view returns(bool)
func (_Preconfmanager *PreconfmanagerCallerSession) ProcessedTxnHashes(arg0 [32]byte) (bool, error) {
	return _Preconfmanager.Contract.ProcessedTxnHashes(&_Preconfmanager.CallOpts, arg0)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerCaller) ProviderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "providerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerSession) ProviderRegistry() (common.Address, error) {
	return _Preconfmanager.Contract.ProviderRegistry(&_Preconfmanager.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfmanager *PreconfmanagerCallerSession) ProviderRegistry() (common.Address, error) {
	return _Preconfmanager.Contract.ProviderRegistry(&_Preconfmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) ProxiableUUID() ([32]byte, error) {
	return _Preconfmanager.Contract.ProxiableUUID(&_Preconfmanager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Preconfmanager.Contract.ProxiableUUID(&_Preconfmanager.CallOpts)
}

// UnopenedCommitments is a free data retrieval call binding the contract method 0x0fb7e259.
//
// Solidity: function unopenedCommitments(bytes32 ) view returns(bool isOpened, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfmanager *PreconfmanagerCaller) UnopenedCommitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	IsOpened            bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "unopenedCommitments", arg0)

	outstruct := new(struct {
		IsOpened            bool
		Committer           common.Address
		DispatchTimestamp   uint64
		CommitmentDigest    [32]byte
		CommitmentSignature []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsOpened = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Committer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DispatchTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.CommitmentDigest = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[4], new([]byte)).(*[]byte)

	return *outstruct, err

}

// UnopenedCommitments is a free data retrieval call binding the contract method 0x0fb7e259.
//
// Solidity: function unopenedCommitments(bytes32 ) view returns(bool isOpened, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfmanager *PreconfmanagerSession) UnopenedCommitments(arg0 [32]byte) (struct {
	IsOpened            bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	return _Preconfmanager.Contract.UnopenedCommitments(&_Preconfmanager.CallOpts, arg0)
}

// UnopenedCommitments is a free data retrieval call binding the contract method 0x0fb7e259.
//
// Solidity: function unopenedCommitments(bytes32 ) view returns(bool isOpened, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfmanager *PreconfmanagerCallerSession) UnopenedCommitments(arg0 [32]byte) (struct {
	IsOpened            bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	return _Preconfmanager.Contract.UnopenedCommitments(&_Preconfmanager.CallOpts, arg0)
}

// VerifyBid is a free data retrieval call binding the contract method 0x9d243a6a.
//
// Solidity: function verifyBid((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfmanager *PreconfmanagerCaller) VerifyBid(opts *bind.CallOpts, params IPreconfManagerOpenCommitmentParams) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "verifyBid", params)

	outstruct := new(struct {
		MessageDigest    [32]byte
		RecoveredAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MessageDigest = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.RecoveredAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// VerifyBid is a free data retrieval call binding the contract method 0x9d243a6a.
//
// Solidity: function verifyBid((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfmanager *PreconfmanagerSession) VerifyBid(params IPreconfManagerOpenCommitmentParams) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	return _Preconfmanager.Contract.VerifyBid(&_Preconfmanager.CallOpts, params)
}

// VerifyBid is a free data retrieval call binding the contract method 0x9d243a6a.
//
// Solidity: function verifyBid((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) view returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfmanager *PreconfmanagerCallerSession) VerifyBid(params IPreconfManagerOpenCommitmentParams) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	return _Preconfmanager.Contract.VerifyBid(&_Preconfmanager.CallOpts, params)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xb0a33f4e.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,uint256[]) params) view returns(bytes32 preConfHash, address committerAddress)
func (_Preconfmanager *PreconfmanagerCaller) VerifyPreConfCommitment(opts *bind.CallOpts, params IPreconfManagerCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "verifyPreConfCommitment", params)

	outstruct := new(struct {
		PreConfHash      [32]byte
		CommitterAddress common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PreConfHash = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.CommitterAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xb0a33f4e.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,uint256[]) params) view returns(bytes32 preConfHash, address committerAddress)
func (_Preconfmanager *PreconfmanagerSession) VerifyPreConfCommitment(params IPreconfManagerCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	return _Preconfmanager.Contract.VerifyPreConfCommitment(&_Preconfmanager.CallOpts, params)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xb0a33f4e.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,uint256[]) params) view returns(bytes32 preConfHash, address committerAddress)
func (_Preconfmanager *PreconfmanagerCallerSession) VerifyPreConfCommitment(params IPreconfManagerCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	return _Preconfmanager.Contract.VerifyPreConfCommitment(&_Preconfmanager.CallOpts, params)
}

// ZkContextHash is a free data retrieval call binding the contract method 0x32a5c16c.
//
// Solidity: function zkContextHash() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCaller) ZkContextHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfmanager.contract.Call(opts, &out, "zkContextHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ZkContextHash is a free data retrieval call binding the contract method 0x32a5c16c.
//
// Solidity: function zkContextHash() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerSession) ZkContextHash() ([32]byte, error) {
	return _Preconfmanager.Contract.ZkContextHash(&_Preconfmanager.CallOpts)
}

// ZkContextHash is a free data retrieval call binding the contract method 0x32a5c16c.
//
// Solidity: function zkContextHash() view returns(bytes32)
func (_Preconfmanager *PreconfmanagerCallerSession) ZkContextHash() ([32]byte, error) {
	return _Preconfmanager.Contract.ZkContextHash(&_Preconfmanager.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfmanager *PreconfmanagerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfmanager *PreconfmanagerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Preconfmanager.Contract.AcceptOwnership(&_Preconfmanager.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Preconfmanager.Contract.AcceptOwnership(&_Preconfmanager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x5a606b41.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerTransactor) Initialize(opts *bind.TransactOpts, _providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "initialize", _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow)
}

// Initialize is a paid mutator transaction binding the contract method 0x5a606b41.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerSession) Initialize(_providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.Initialize(&_Preconfmanager.TransactOpts, _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow)
}

// Initialize is a paid mutator transaction binding the contract method 0x5a606b41.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) Initialize(_providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.Initialize(&_Preconfmanager.TransactOpts, _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerTransactor) InitiateReward(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "initiateReward", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.Contract.InitiateReward(&_Preconfmanager.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.Contract.InitiateReward(&_Preconfmanager.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerTransactor) InitiateSlash(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "initiateSlash", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.Contract.InitiateSlash(&_Preconfmanager.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfmanager.Contract.InitiateSlash(&_Preconfmanager.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x82a9ba0f.
//
// Solidity: function openCommitment((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerTransactor) OpenCommitment(opts *bind.TransactOpts, params IPreconfManagerOpenCommitmentParams) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "openCommitment", params)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x82a9ba0f.
//
// Solidity: function openCommitment((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerSession) OpenCommitment(params IPreconfManagerOpenCommitmentParams) (*types.Transaction, error) {
	return _Preconfmanager.Contract.OpenCommitment(&_Preconfmanager.TransactOpts, params)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x82a9ba0f.
//
// Solidity: function openCommitment((bytes32,uint256,uint256,uint64,uint64,uint64,string,string,bytes,uint256[]) params) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerTransactorSession) OpenCommitment(params IPreconfManagerOpenCommitmentParams) (*types.Transaction, error) {
	return _Preconfmanager.Contract.OpenCommitment(&_Preconfmanager.TransactOpts, params)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Preconfmanager *PreconfmanagerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Preconfmanager *PreconfmanagerSession) Pause() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Pause(&_Preconfmanager.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) Pause() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Pause(&_Preconfmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfmanager *PreconfmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfmanager *PreconfmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfmanager.Contract.RenounceOwnership(&_Preconfmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfmanager.Contract.RenounceOwnership(&_Preconfmanager.TransactOpts)
}

// StoreUnopenedCommitment is a paid mutator transaction binding the contract method 0x692c575d.
//
// Solidity: function storeUnopenedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerTransactor) StoreUnopenedCommitment(opts *bind.TransactOpts, commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "storeUnopenedCommitment", commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// StoreUnopenedCommitment is a paid mutator transaction binding the contract method 0x692c575d.
//
// Solidity: function storeUnopenedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerSession) StoreUnopenedCommitment(commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.StoreUnopenedCommitment(&_Preconfmanager.TransactOpts, commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// StoreUnopenedCommitment is a paid mutator transaction binding the contract method 0x692c575d.
//
// Solidity: function storeUnopenedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfmanager *PreconfmanagerTransactorSession) StoreUnopenedCommitment(commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.StoreUnopenedCommitment(&_Preconfmanager.TransactOpts, commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfmanager *PreconfmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfmanager *PreconfmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.TransferOwnership(&_Preconfmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.TransferOwnership(&_Preconfmanager.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Preconfmanager *PreconfmanagerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Preconfmanager *PreconfmanagerSession) Unpause() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Unpause(&_Preconfmanager.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Unpause(&_Preconfmanager.TransactOpts)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpdateBidderRegistry(opts *bind.TransactOpts, newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "updateBidderRegistry", newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfmanager *PreconfmanagerSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateBidderRegistry(&_Preconfmanager.TransactOpts, newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateBidderRegistry(&_Preconfmanager.TransactOpts, newBidderRegistry)
}

// UpdateBlockTracker is a paid mutator transaction binding the contract method 0x8cb8da1d.
//
// Solidity: function updateBlockTracker(address newBlockTracker) returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpdateBlockTracker(opts *bind.TransactOpts, newBlockTracker common.Address) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "updateBlockTracker", newBlockTracker)
}

// UpdateBlockTracker is a paid mutator transaction binding the contract method 0x8cb8da1d.
//
// Solidity: function updateBlockTracker(address newBlockTracker) returns()
func (_Preconfmanager *PreconfmanagerSession) UpdateBlockTracker(newBlockTracker common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateBlockTracker(&_Preconfmanager.TransactOpts, newBlockTracker)
}

// UpdateBlockTracker is a paid mutator transaction binding the contract method 0x8cb8da1d.
//
// Solidity: function updateBlockTracker(address newBlockTracker) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpdateBlockTracker(newBlockTracker common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateBlockTracker(&_Preconfmanager.TransactOpts, newBlockTracker)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpdateCommitmentDispatchWindow(opts *bind.TransactOpts, newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "updateCommitmentDispatchWindow", newDispatchWindow)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerSession) UpdateCommitmentDispatchWindow(newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateCommitmentDispatchWindow(&_Preconfmanager.TransactOpts, newDispatchWindow)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpdateCommitmentDispatchWindow(newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateCommitmentDispatchWindow(&_Preconfmanager.TransactOpts, newDispatchWindow)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpdateOracleContract(opts *bind.TransactOpts, newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "updateOracleContract", newOracleContract)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfmanager *PreconfmanagerSession) UpdateOracleContract(newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateOracleContract(&_Preconfmanager.TransactOpts, newOracleContract)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpdateOracleContract(newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateOracleContract(&_Preconfmanager.TransactOpts, newOracleContract)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpdateProviderRegistry(opts *bind.TransactOpts, newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "updateProviderRegistry", newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfmanager *PreconfmanagerSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateProviderRegistry(&_Preconfmanager.TransactOpts, newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpdateProviderRegistry(&_Preconfmanager.TransactOpts, newProviderRegistry)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfmanager *PreconfmanagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfmanager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfmanager *PreconfmanagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpgradeToAndCall(&_Preconfmanager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfmanager.Contract.UpgradeToAndCall(&_Preconfmanager.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfmanager *PreconfmanagerTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Preconfmanager.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfmanager *PreconfmanagerSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfmanager.Contract.Fallback(&_Preconfmanager.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfmanager.Contract.Fallback(&_Preconfmanager.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfmanager *PreconfmanagerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfmanager.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfmanager *PreconfmanagerSession) Receive() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Receive(&_Preconfmanager.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfmanager *PreconfmanagerTransactorSession) Receive() (*types.Transaction, error) {
	return _Preconfmanager.Contract.Receive(&_Preconfmanager.TransactOpts)
}

// PreconfmanagerBidderRegistryUpdatedIterator is returned from FilterBidderRegistryUpdated and is used to iterate over the raw logs and unpacked data for BidderRegistryUpdated events raised by the Preconfmanager contract.
type PreconfmanagerBidderRegistryUpdatedIterator struct {
	Event *PreconfmanagerBidderRegistryUpdated // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerBidderRegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerBidderRegistryUpdated)
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
		it.Event = new(PreconfmanagerBidderRegistryUpdated)
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
func (it *PreconfmanagerBidderRegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerBidderRegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerBidderRegistryUpdated represents a BidderRegistryUpdated event raised by the Preconfmanager contract.
type PreconfmanagerBidderRegistryUpdated struct {
	NewBidderRegistry common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterBidderRegistryUpdated is a free log retrieval operation binding the contract event 0xa1e7ee70f66ed8d94bde04c89c39cad2377817914a330274e559ba8e6d5cdca5.
//
// Solidity: event BidderRegistryUpdated(address indexed newBidderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) FilterBidderRegistryUpdated(opts *bind.FilterOpts, newBidderRegistry []common.Address) (*PreconfmanagerBidderRegistryUpdatedIterator, error) {

	var newBidderRegistryRule []interface{}
	for _, newBidderRegistryItem := range newBidderRegistry {
		newBidderRegistryRule = append(newBidderRegistryRule, newBidderRegistryItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "BidderRegistryUpdated", newBidderRegistryRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerBidderRegistryUpdatedIterator{contract: _Preconfmanager.contract, event: "BidderRegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchBidderRegistryUpdated is a free log subscription operation binding the contract event 0xa1e7ee70f66ed8d94bde04c89c39cad2377817914a330274e559ba8e6d5cdca5.
//
// Solidity: event BidderRegistryUpdated(address indexed newBidderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) WatchBidderRegistryUpdated(opts *bind.WatchOpts, sink chan<- *PreconfmanagerBidderRegistryUpdated, newBidderRegistry []common.Address) (event.Subscription, error) {

	var newBidderRegistryRule []interface{}
	for _, newBidderRegistryItem := range newBidderRegistry {
		newBidderRegistryRule = append(newBidderRegistryRule, newBidderRegistryItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "BidderRegistryUpdated", newBidderRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerBidderRegistryUpdated)
				if err := _Preconfmanager.contract.UnpackLog(event, "BidderRegistryUpdated", log); err != nil {
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

// ParseBidderRegistryUpdated is a log parse operation binding the contract event 0xa1e7ee70f66ed8d94bde04c89c39cad2377817914a330274e559ba8e6d5cdca5.
//
// Solidity: event BidderRegistryUpdated(address indexed newBidderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) ParseBidderRegistryUpdated(log types.Log) (*PreconfmanagerBidderRegistryUpdated, error) {
	event := new(PreconfmanagerBidderRegistryUpdated)
	if err := _Preconfmanager.contract.UnpackLog(event, "BidderRegistryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerBlockTrackerUpdatedIterator is returned from FilterBlockTrackerUpdated and is used to iterate over the raw logs and unpacked data for BlockTrackerUpdated events raised by the Preconfmanager contract.
type PreconfmanagerBlockTrackerUpdatedIterator struct {
	Event *PreconfmanagerBlockTrackerUpdated // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerBlockTrackerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerBlockTrackerUpdated)
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
		it.Event = new(PreconfmanagerBlockTrackerUpdated)
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
func (it *PreconfmanagerBlockTrackerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerBlockTrackerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerBlockTrackerUpdated represents a BlockTrackerUpdated event raised by the Preconfmanager contract.
type PreconfmanagerBlockTrackerUpdated struct {
	NewBlockTracker common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBlockTrackerUpdated is a free log retrieval operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Preconfmanager *PreconfmanagerFilterer) FilterBlockTrackerUpdated(opts *bind.FilterOpts, newBlockTracker []common.Address) (*PreconfmanagerBlockTrackerUpdatedIterator, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "BlockTrackerUpdated", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerBlockTrackerUpdatedIterator{contract: _Preconfmanager.contract, event: "BlockTrackerUpdated", logs: logs, sub: sub}, nil
}

// WatchBlockTrackerUpdated is a free log subscription operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Preconfmanager *PreconfmanagerFilterer) WatchBlockTrackerUpdated(opts *bind.WatchOpts, sink chan<- *PreconfmanagerBlockTrackerUpdated, newBlockTracker []common.Address) (event.Subscription, error) {

	var newBlockTrackerRule []interface{}
	for _, newBlockTrackerItem := range newBlockTracker {
		newBlockTrackerRule = append(newBlockTrackerRule, newBlockTrackerItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "BlockTrackerUpdated", newBlockTrackerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerBlockTrackerUpdated)
				if err := _Preconfmanager.contract.UnpackLog(event, "BlockTrackerUpdated", log); err != nil {
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

// ParseBlockTrackerUpdated is a log parse operation binding the contract event 0x3a013345829f05d7c43406984b75f593c7bf0f77fc18042ff8de7f70314935f0.
//
// Solidity: event BlockTrackerUpdated(address indexed newBlockTracker)
func (_Preconfmanager *PreconfmanagerFilterer) ParseBlockTrackerUpdated(log types.Log) (*PreconfmanagerBlockTrackerUpdated, error) {
	event := new(PreconfmanagerBlockTrackerUpdated)
	if err := _Preconfmanager.contract.UnpackLog(event, "BlockTrackerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerCommitmentDispatchWindowUpdatedIterator is returned from FilterCommitmentDispatchWindowUpdated and is used to iterate over the raw logs and unpacked data for CommitmentDispatchWindowUpdated events raised by the Preconfmanager contract.
type PreconfmanagerCommitmentDispatchWindowUpdatedIterator struct {
	Event *PreconfmanagerCommitmentDispatchWindowUpdated // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerCommitmentDispatchWindowUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerCommitmentDispatchWindowUpdated)
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
		it.Event = new(PreconfmanagerCommitmentDispatchWindowUpdated)
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
func (it *PreconfmanagerCommitmentDispatchWindowUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerCommitmentDispatchWindowUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerCommitmentDispatchWindowUpdated represents a CommitmentDispatchWindowUpdated event raised by the Preconfmanager contract.
type PreconfmanagerCommitmentDispatchWindowUpdated struct {
	NewDispatchWindow uint64
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterCommitmentDispatchWindowUpdated is a free log retrieval operation binding the contract event 0xfc32b165cb91c9802d41143cf3b5ebb294250107d79ef586a36c7d3f370bc107.
//
// Solidity: event CommitmentDispatchWindowUpdated(uint64 newDispatchWindow)
func (_Preconfmanager *PreconfmanagerFilterer) FilterCommitmentDispatchWindowUpdated(opts *bind.FilterOpts) (*PreconfmanagerCommitmentDispatchWindowUpdatedIterator, error) {

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "CommitmentDispatchWindowUpdated")
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerCommitmentDispatchWindowUpdatedIterator{contract: _Preconfmanager.contract, event: "CommitmentDispatchWindowUpdated", logs: logs, sub: sub}, nil
}

// WatchCommitmentDispatchWindowUpdated is a free log subscription operation binding the contract event 0xfc32b165cb91c9802d41143cf3b5ebb294250107d79ef586a36c7d3f370bc107.
//
// Solidity: event CommitmentDispatchWindowUpdated(uint64 newDispatchWindow)
func (_Preconfmanager *PreconfmanagerFilterer) WatchCommitmentDispatchWindowUpdated(opts *bind.WatchOpts, sink chan<- *PreconfmanagerCommitmentDispatchWindowUpdated) (event.Subscription, error) {

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "CommitmentDispatchWindowUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerCommitmentDispatchWindowUpdated)
				if err := _Preconfmanager.contract.UnpackLog(event, "CommitmentDispatchWindowUpdated", log); err != nil {
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

// ParseCommitmentDispatchWindowUpdated is a log parse operation binding the contract event 0xfc32b165cb91c9802d41143cf3b5ebb294250107d79ef586a36c7d3f370bc107.
//
// Solidity: event CommitmentDispatchWindowUpdated(uint64 newDispatchWindow)
func (_Preconfmanager *PreconfmanagerFilterer) ParseCommitmentDispatchWindowUpdated(log types.Log) (*PreconfmanagerCommitmentDispatchWindowUpdated, error) {
	event := new(PreconfmanagerCommitmentDispatchWindowUpdated)
	if err := _Preconfmanager.contract.UnpackLog(event, "CommitmentDispatchWindowUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Preconfmanager contract.
type PreconfmanagerInitializedIterator struct {
	Event *PreconfmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerInitialized)
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
		it.Event = new(PreconfmanagerInitialized)
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
func (it *PreconfmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerInitialized represents a Initialized event raised by the Preconfmanager contract.
type PreconfmanagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Preconfmanager *PreconfmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*PreconfmanagerInitializedIterator, error) {

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerInitializedIterator{contract: _Preconfmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Preconfmanager *PreconfmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PreconfmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerInitialized)
				if err := _Preconfmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParseInitialized(log types.Log) (*PreconfmanagerInitialized, error) {
	event := new(PreconfmanagerInitialized)
	if err := _Preconfmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerOpenedCommitmentStoredIterator is returned from FilterOpenedCommitmentStored and is used to iterate over the raw logs and unpacked data for OpenedCommitmentStored events raised by the Preconfmanager contract.
type PreconfmanagerOpenedCommitmentStoredIterator struct {
	Event *PreconfmanagerOpenedCommitmentStored // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerOpenedCommitmentStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerOpenedCommitmentStored)
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
		it.Event = new(PreconfmanagerOpenedCommitmentStored)
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
func (it *PreconfmanagerOpenedCommitmentStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerOpenedCommitmentStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerOpenedCommitmentStored represents a OpenedCommitmentStored event raised by the Preconfmanager contract.
type PreconfmanagerOpenedCommitmentStored struct {
	CommitmentIndex     [32]byte
	Bidder              common.Address
	Committer           common.Address
	BidAmt              *big.Int
	SlashAmt            *big.Int
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	RevertingTxHashes   string
	CommitmentDigest    [32]byte
	DispatchTimestamp   uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterOpenedCommitmentStored is a free log retrieval operation binding the contract event 0x4e798091cca241fe6b4417890663555a768ee632a5ab90275a1c7ce4512e58d8.
//
// Solidity: event OpenedCommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bidAmt, uint256 slashAmt, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentDigest, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) FilterOpenedCommitmentStored(opts *bind.FilterOpts, commitmentIndex [][32]byte) (*PreconfmanagerOpenedCommitmentStoredIterator, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "OpenedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerOpenedCommitmentStoredIterator{contract: _Preconfmanager.contract, event: "OpenedCommitmentStored", logs: logs, sub: sub}, nil
}

// WatchOpenedCommitmentStored is a free log subscription operation binding the contract event 0x4e798091cca241fe6b4417890663555a768ee632a5ab90275a1c7ce4512e58d8.
//
// Solidity: event OpenedCommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bidAmt, uint256 slashAmt, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentDigest, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) WatchOpenedCommitmentStored(opts *bind.WatchOpts, sink chan<- *PreconfmanagerOpenedCommitmentStored, commitmentIndex [][32]byte) (event.Subscription, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "OpenedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerOpenedCommitmentStored)
				if err := _Preconfmanager.contract.UnpackLog(event, "OpenedCommitmentStored", log); err != nil {
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

// ParseOpenedCommitmentStored is a log parse operation binding the contract event 0x4e798091cca241fe6b4417890663555a768ee632a5ab90275a1c7ce4512e58d8.
//
// Solidity: event OpenedCommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bidAmt, uint256 slashAmt, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentDigest, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) ParseOpenedCommitmentStored(log types.Log) (*PreconfmanagerOpenedCommitmentStored, error) {
	event := new(PreconfmanagerOpenedCommitmentStored)
	if err := _Preconfmanager.contract.UnpackLog(event, "OpenedCommitmentStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerOracleContractUpdatedIterator is returned from FilterOracleContractUpdated and is used to iterate over the raw logs and unpacked data for OracleContractUpdated events raised by the Preconfmanager contract.
type PreconfmanagerOracleContractUpdatedIterator struct {
	Event *PreconfmanagerOracleContractUpdated // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerOracleContractUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerOracleContractUpdated)
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
		it.Event = new(PreconfmanagerOracleContractUpdated)
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
func (it *PreconfmanagerOracleContractUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerOracleContractUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerOracleContractUpdated represents a OracleContractUpdated event raised by the Preconfmanager contract.
type PreconfmanagerOracleContractUpdated struct {
	NewOracleContract common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterOracleContractUpdated is a free log retrieval operation binding the contract event 0x9f17f88b8f07aebd1a366b1589656e38eddd6fbd5a1fa121158b240d4795a0cb.
//
// Solidity: event OracleContractUpdated(address indexed newOracleContract)
func (_Preconfmanager *PreconfmanagerFilterer) FilterOracleContractUpdated(opts *bind.FilterOpts, newOracleContract []common.Address) (*PreconfmanagerOracleContractUpdatedIterator, error) {

	var newOracleContractRule []interface{}
	for _, newOracleContractItem := range newOracleContract {
		newOracleContractRule = append(newOracleContractRule, newOracleContractItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "OracleContractUpdated", newOracleContractRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerOracleContractUpdatedIterator{contract: _Preconfmanager.contract, event: "OracleContractUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleContractUpdated is a free log subscription operation binding the contract event 0x9f17f88b8f07aebd1a366b1589656e38eddd6fbd5a1fa121158b240d4795a0cb.
//
// Solidity: event OracleContractUpdated(address indexed newOracleContract)
func (_Preconfmanager *PreconfmanagerFilterer) WatchOracleContractUpdated(opts *bind.WatchOpts, sink chan<- *PreconfmanagerOracleContractUpdated, newOracleContract []common.Address) (event.Subscription, error) {

	var newOracleContractRule []interface{}
	for _, newOracleContractItem := range newOracleContract {
		newOracleContractRule = append(newOracleContractRule, newOracleContractItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "OracleContractUpdated", newOracleContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerOracleContractUpdated)
				if err := _Preconfmanager.contract.UnpackLog(event, "OracleContractUpdated", log); err != nil {
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

// ParseOracleContractUpdated is a log parse operation binding the contract event 0x9f17f88b8f07aebd1a366b1589656e38eddd6fbd5a1fa121158b240d4795a0cb.
//
// Solidity: event OracleContractUpdated(address indexed newOracleContract)
func (_Preconfmanager *PreconfmanagerFilterer) ParseOracleContractUpdated(log types.Log) (*PreconfmanagerOracleContractUpdated, error) {
	event := new(PreconfmanagerOracleContractUpdated)
	if err := _Preconfmanager.contract.UnpackLog(event, "OracleContractUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Preconfmanager contract.
type PreconfmanagerOwnershipTransferStartedIterator struct {
	Event *PreconfmanagerOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerOwnershipTransferStarted)
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
		it.Event = new(PreconfmanagerOwnershipTransferStarted)
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
func (it *PreconfmanagerOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Preconfmanager contract.
type PreconfmanagerOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Preconfmanager *PreconfmanagerFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfmanagerOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerOwnershipTransferStartedIterator{contract: _Preconfmanager.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Preconfmanager *PreconfmanagerFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *PreconfmanagerOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerOwnershipTransferStarted)
				if err := _Preconfmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParseOwnershipTransferStarted(log types.Log) (*PreconfmanagerOwnershipTransferStarted, error) {
	event := new(PreconfmanagerOwnershipTransferStarted)
	if err := _Preconfmanager.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Preconfmanager contract.
type PreconfmanagerOwnershipTransferredIterator struct {
	Event *PreconfmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerOwnershipTransferred)
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
		it.Event = new(PreconfmanagerOwnershipTransferred)
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
func (it *PreconfmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Preconfmanager contract.
type PreconfmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfmanager *PreconfmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerOwnershipTransferredIterator{contract: _Preconfmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfmanager *PreconfmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PreconfmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerOwnershipTransferred)
				if err := _Preconfmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParseOwnershipTransferred(log types.Log) (*PreconfmanagerOwnershipTransferred, error) {
	event := new(PreconfmanagerOwnershipTransferred)
	if err := _Preconfmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Preconfmanager contract.
type PreconfmanagerPausedIterator struct {
	Event *PreconfmanagerPaused // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerPaused)
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
		it.Event = new(PreconfmanagerPaused)
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
func (it *PreconfmanagerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerPaused represents a Paused event raised by the Preconfmanager contract.
type PreconfmanagerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Preconfmanager *PreconfmanagerFilterer) FilterPaused(opts *bind.FilterOpts) (*PreconfmanagerPausedIterator, error) {

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerPausedIterator{contract: _Preconfmanager.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Preconfmanager *PreconfmanagerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *PreconfmanagerPaused) (event.Subscription, error) {

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerPaused)
				if err := _Preconfmanager.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParsePaused(log types.Log) (*PreconfmanagerPaused, error) {
	event := new(PreconfmanagerPaused)
	if err := _Preconfmanager.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerProviderRegistryUpdatedIterator is returned from FilterProviderRegistryUpdated and is used to iterate over the raw logs and unpacked data for ProviderRegistryUpdated events raised by the Preconfmanager contract.
type PreconfmanagerProviderRegistryUpdatedIterator struct {
	Event *PreconfmanagerProviderRegistryUpdated // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerProviderRegistryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerProviderRegistryUpdated)
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
		it.Event = new(PreconfmanagerProviderRegistryUpdated)
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
func (it *PreconfmanagerProviderRegistryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerProviderRegistryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerProviderRegistryUpdated represents a ProviderRegistryUpdated event raised by the Preconfmanager contract.
type PreconfmanagerProviderRegistryUpdated struct {
	NewProviderRegistry common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterProviderRegistryUpdated is a free log retrieval operation binding the contract event 0xb8f8a6768be726978fef2f5342ee7b1d0eead3b5cbb0c5d55c68f56271ca0eeb.
//
// Solidity: event ProviderRegistryUpdated(address indexed newProviderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) FilterProviderRegistryUpdated(opts *bind.FilterOpts, newProviderRegistry []common.Address) (*PreconfmanagerProviderRegistryUpdatedIterator, error) {

	var newProviderRegistryRule []interface{}
	for _, newProviderRegistryItem := range newProviderRegistry {
		newProviderRegistryRule = append(newProviderRegistryRule, newProviderRegistryItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "ProviderRegistryUpdated", newProviderRegistryRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerProviderRegistryUpdatedIterator{contract: _Preconfmanager.contract, event: "ProviderRegistryUpdated", logs: logs, sub: sub}, nil
}

// WatchProviderRegistryUpdated is a free log subscription operation binding the contract event 0xb8f8a6768be726978fef2f5342ee7b1d0eead3b5cbb0c5d55c68f56271ca0eeb.
//
// Solidity: event ProviderRegistryUpdated(address indexed newProviderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) WatchProviderRegistryUpdated(opts *bind.WatchOpts, sink chan<- *PreconfmanagerProviderRegistryUpdated, newProviderRegistry []common.Address) (event.Subscription, error) {

	var newProviderRegistryRule []interface{}
	for _, newProviderRegistryItem := range newProviderRegistry {
		newProviderRegistryRule = append(newProviderRegistryRule, newProviderRegistryItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "ProviderRegistryUpdated", newProviderRegistryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerProviderRegistryUpdated)
				if err := _Preconfmanager.contract.UnpackLog(event, "ProviderRegistryUpdated", log); err != nil {
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

// ParseProviderRegistryUpdated is a log parse operation binding the contract event 0xb8f8a6768be726978fef2f5342ee7b1d0eead3b5cbb0c5d55c68f56271ca0eeb.
//
// Solidity: event ProviderRegistryUpdated(address indexed newProviderRegistry)
func (_Preconfmanager *PreconfmanagerFilterer) ParseProviderRegistryUpdated(log types.Log) (*PreconfmanagerProviderRegistryUpdated, error) {
	event := new(PreconfmanagerProviderRegistryUpdated)
	if err := _Preconfmanager.contract.UnpackLog(event, "ProviderRegistryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerSignatureVerifiedIterator is returned from FilterSignatureVerified and is used to iterate over the raw logs and unpacked data for SignatureVerified events raised by the Preconfmanager contract.
type PreconfmanagerSignatureVerifiedIterator struct {
	Event *PreconfmanagerSignatureVerified // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerSignatureVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerSignatureVerified)
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
		it.Event = new(PreconfmanagerSignatureVerified)
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
func (it *PreconfmanagerSignatureVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerSignatureVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerSignatureVerified represents a SignatureVerified event raised by the Preconfmanager contract.
type PreconfmanagerSignatureVerified struct {
	Signer            common.Address
	TxnHash           string
	RevertingTxHashes string
	BidAmt            *big.Int
	BlockNumber       uint64
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSignatureVerified is a free log retrieval operation binding the contract event 0xf5fa3f7902916f0e4e2af218fa9a77809c54e64b77ec4061a205ec3b7ce5c35e.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bidAmt, uint64 blockNumber)
func (_Preconfmanager *PreconfmanagerFilterer) FilterSignatureVerified(opts *bind.FilterOpts, signer []common.Address, bidAmt []*big.Int) (*PreconfmanagerSignatureVerifiedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidAmtRule []interface{}
	for _, bidAmtItem := range bidAmt {
		bidAmtRule = append(bidAmtRule, bidAmtItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "SignatureVerified", signerRule, bidAmtRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerSignatureVerifiedIterator{contract: _Preconfmanager.contract, event: "SignatureVerified", logs: logs, sub: sub}, nil
}

// WatchSignatureVerified is a free log subscription operation binding the contract event 0xf5fa3f7902916f0e4e2af218fa9a77809c54e64b77ec4061a205ec3b7ce5c35e.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bidAmt, uint64 blockNumber)
func (_Preconfmanager *PreconfmanagerFilterer) WatchSignatureVerified(opts *bind.WatchOpts, sink chan<- *PreconfmanagerSignatureVerified, signer []common.Address, bidAmt []*big.Int) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidAmtRule []interface{}
	for _, bidAmtItem := range bidAmt {
		bidAmtRule = append(bidAmtRule, bidAmtItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "SignatureVerified", signerRule, bidAmtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerSignatureVerified)
				if err := _Preconfmanager.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
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

// ParseSignatureVerified is a log parse operation binding the contract event 0xf5fa3f7902916f0e4e2af218fa9a77809c54e64b77ec4061a205ec3b7ce5c35e.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bidAmt, uint64 blockNumber)
func (_Preconfmanager *PreconfmanagerFilterer) ParseSignatureVerified(log types.Log) (*PreconfmanagerSignatureVerified, error) {
	event := new(PreconfmanagerSignatureVerified)
	if err := _Preconfmanager.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerUnopenedCommitmentStoredIterator is returned from FilterUnopenedCommitmentStored and is used to iterate over the raw logs and unpacked data for UnopenedCommitmentStored events raised by the Preconfmanager contract.
type PreconfmanagerUnopenedCommitmentStoredIterator struct {
	Event *PreconfmanagerUnopenedCommitmentStored // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerUnopenedCommitmentStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerUnopenedCommitmentStored)
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
		it.Event = new(PreconfmanagerUnopenedCommitmentStored)
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
func (it *PreconfmanagerUnopenedCommitmentStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerUnopenedCommitmentStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerUnopenedCommitmentStored represents a UnopenedCommitmentStored event raised by the Preconfmanager contract.
type PreconfmanagerUnopenedCommitmentStored struct {
	CommitmentIndex     [32]byte
	Committer           common.Address
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	DispatchTimestamp   uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterUnopenedCommitmentStored is a free log retrieval operation binding the contract event 0xbe650dcd1894b46cca996156910a6f10e521dfcb65da07e885f41ae7d2db3c78.
//
// Solidity: event UnopenedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) FilterUnopenedCommitmentStored(opts *bind.FilterOpts, commitmentIndex [][32]byte) (*PreconfmanagerUnopenedCommitmentStoredIterator, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "UnopenedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerUnopenedCommitmentStoredIterator{contract: _Preconfmanager.contract, event: "UnopenedCommitmentStored", logs: logs, sub: sub}, nil
}

// WatchUnopenedCommitmentStored is a free log subscription operation binding the contract event 0xbe650dcd1894b46cca996156910a6f10e521dfcb65da07e885f41ae7d2db3c78.
//
// Solidity: event UnopenedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) WatchUnopenedCommitmentStored(opts *bind.WatchOpts, sink chan<- *PreconfmanagerUnopenedCommitmentStored, commitmentIndex [][32]byte) (event.Subscription, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "UnopenedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerUnopenedCommitmentStored)
				if err := _Preconfmanager.contract.UnpackLog(event, "UnopenedCommitmentStored", log); err != nil {
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

// ParseUnopenedCommitmentStored is a log parse operation binding the contract event 0xbe650dcd1894b46cca996156910a6f10e521dfcb65da07e885f41ae7d2db3c78.
//
// Solidity: event UnopenedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfmanager *PreconfmanagerFilterer) ParseUnopenedCommitmentStored(log types.Log) (*PreconfmanagerUnopenedCommitmentStored, error) {
	event := new(PreconfmanagerUnopenedCommitmentStored)
	if err := _Preconfmanager.contract.UnpackLog(event, "UnopenedCommitmentStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Preconfmanager contract.
type PreconfmanagerUnpausedIterator struct {
	Event *PreconfmanagerUnpaused // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerUnpaused)
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
		it.Event = new(PreconfmanagerUnpaused)
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
func (it *PreconfmanagerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerUnpaused represents a Unpaused event raised by the Preconfmanager contract.
type PreconfmanagerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Preconfmanager *PreconfmanagerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*PreconfmanagerUnpausedIterator, error) {

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerUnpausedIterator{contract: _Preconfmanager.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Preconfmanager *PreconfmanagerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *PreconfmanagerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerUnpaused)
				if err := _Preconfmanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParseUnpaused(log types.Log) (*PreconfmanagerUnpaused, error) {
	event := new(PreconfmanagerUnpaused)
	if err := _Preconfmanager.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfmanagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Preconfmanager contract.
type PreconfmanagerUpgradedIterator struct {
	Event *PreconfmanagerUpgraded // Event containing the contract specifics and raw log

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
func (it *PreconfmanagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfmanagerUpgraded)
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
		it.Event = new(PreconfmanagerUpgraded)
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
func (it *PreconfmanagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfmanagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfmanagerUpgraded represents a Upgraded event raised by the Preconfmanager contract.
type PreconfmanagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Preconfmanager *PreconfmanagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PreconfmanagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Preconfmanager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PreconfmanagerUpgradedIterator{contract: _Preconfmanager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Preconfmanager *PreconfmanagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PreconfmanagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Preconfmanager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfmanagerUpgraded)
				if err := _Preconfmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Preconfmanager *PreconfmanagerFilterer) ParseUpgraded(log types.Log) (*PreconfmanagerUpgraded, error) {
	event := new(PreconfmanagerUpgraded)
	if err := _Preconfmanager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
