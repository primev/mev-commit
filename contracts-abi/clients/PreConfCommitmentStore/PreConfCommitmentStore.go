// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package preconfcommitmentstore

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

// IPreConfCommitmentStoreCommitmentParams is an auto generated low-level Go binding around an user-defined struct.
type IPreConfCommitmentStoreCommitmentParams struct {
	TxnHash             string
	RevertingTxHashes   string
	Bid                 *big.Int
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	BidHash             [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	SharedSecretKey     []byte
}

// IPreConfCommitmentStoreEncrPreConfCommitment is an auto generated low-level Go binding around an user-defined struct.
type IPreConfCommitmentStoreEncrPreConfCommitment struct {
	IsUsed              bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}

// IPreConfCommitmentStorePreConfCommitment is an auto generated low-level Go binding around an user-defined struct.
type IPreConfCommitmentStorePreConfCommitment struct {
	Bidder              common.Address
	IsUsed              bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	Bid                 *big.Int
	BidHash             [32]byte
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	SharedSecretKey     []byte
	TxnHash             string
	RevertingTxHashes   string
}

// PreconfcommitmentstoreMetaData contains all meta data concerning the Preconfcommitmentstore contract.
var PreconfcommitmentstoreMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_BID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR_PRECONF\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_BID_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EIP712_COMMITMENT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"HEXCHARS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bidderRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBidderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blockTracker\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIBlockTracker\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"blocksPerWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentDispatchWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commitmentsCount\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"encryptedCommitments\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBidHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIPreConfCommitmentStore.PreConfCommitment\",\"components\":[{\"name\":\"bidder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getEncryptedCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIPreConfCommitmentStore.EncrPreConfCommitment\",\"components\":[{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEncryptedCommitmentIndex\",\"inputs\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structIPreConfCommitmentStore.EncrPreConfCommitment\",\"components\":[{\"name\":\"isUsed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"committer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getPreConfHash\",\"inputs\":[{\"name\":\"_txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"_bidSignature\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_sharedSecretKey\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getTxnHashFromCommitment\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_providerRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_bidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_oracleContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_blockTracker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_commitmentDispatchWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_blocksPerWindow\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateReward\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initiateSlash\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"residualBidPercentAfterDecay\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"openCommitment\",\"inputs\":[{\"name\":\"encryptedCommitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"oracleContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"providerRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIProviderRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"storeEncryptedCommitment\",\"inputs\":[{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateBidderRegistry\",\"inputs\":[{\"name\":\"newBidderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCommitmentDispatchWindow\",\"inputs\":[{\"name\":\"newDispatchWindow\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOracleContract\",\"inputs\":[{\"name\":\"newOracleContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateProviderRegistry\",\"inputs\":[{\"name\":\"newProviderRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"verifyBid\",\"inputs\":[{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"messageDigest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"recoveredAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"verifyPreConfCommitment\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIPreConfCommitmentStore.CommitmentParams\",\"components\":[{\"name\":\"txnHash\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"preConfHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"committerAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"CommitmentStored\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"bidder\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"committer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"bid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"bidHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"decayStartTimeStamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"decayEndTimeStamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"commitmentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"bidSignature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"sharedSecretKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EncryptedCommitmentStored\",\"inputs\":[{\"name\":\"commitmentIndex\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"committer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"commitmentDigest\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"commitmentSignature\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"dispatchTimestamp\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SignatureVerified\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"txnHash\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"revertingTxHashes\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"bid\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// PreconfcommitmentstoreABI is the input ABI used to generate the binding from.
// Deprecated: Use PreconfcommitmentstoreMetaData.ABI instead.
var PreconfcommitmentstoreABI = PreconfcommitmentstoreMetaData.ABI

// Preconfcommitmentstore is an auto generated Go binding around an Ethereum contract.
type Preconfcommitmentstore struct {
	PreconfcommitmentstoreCaller     // Read-only binding to the contract
	PreconfcommitmentstoreTransactor // Write-only binding to the contract
	PreconfcommitmentstoreFilterer   // Log filterer for contract events
}

// PreconfcommitmentstoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreconfcommitmentstoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreconfcommitmentstoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreconfcommitmentstoreSession struct {
	Contract     *Preconfcommitmentstore // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PreconfcommitmentstoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreconfcommitmentstoreCallerSession struct {
	Contract *PreconfcommitmentstoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// PreconfcommitmentstoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreconfcommitmentstoreTransactorSession struct {
	Contract     *PreconfcommitmentstoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// PreconfcommitmentstoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreconfcommitmentstoreRaw struct {
	Contract *Preconfcommitmentstore // Generic contract binding to access the raw methods on
}

// PreconfcommitmentstoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreCallerRaw struct {
	Contract *PreconfcommitmentstoreCaller // Generic read-only contract binding to access the raw methods on
}

// PreconfcommitmentstoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreconfcommitmentstoreTransactorRaw struct {
	Contract *PreconfcommitmentstoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreconfcommitmentstore creates a new instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstore(address common.Address, backend bind.ContractBackend) (*Preconfcommitmentstore, error) {
	contract, err := bindPreconfcommitmentstore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Preconfcommitmentstore{PreconfcommitmentstoreCaller: PreconfcommitmentstoreCaller{contract: contract}, PreconfcommitmentstoreTransactor: PreconfcommitmentstoreTransactor{contract: contract}, PreconfcommitmentstoreFilterer: PreconfcommitmentstoreFilterer{contract: contract}}, nil
}

// NewPreconfcommitmentstoreCaller creates a new read-only instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreCaller(address common.Address, caller bind.ContractCaller) (*PreconfcommitmentstoreCaller, error) {
	contract, err := bindPreconfcommitmentstore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreCaller{contract: contract}, nil
}

// NewPreconfcommitmentstoreTransactor creates a new write-only instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreTransactor(address common.Address, transactor bind.ContractTransactor) (*PreconfcommitmentstoreTransactor, error) {
	contract, err := bindPreconfcommitmentstore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreTransactor{contract: contract}, nil
}

// NewPreconfcommitmentstoreFilterer creates a new log filterer instance of Preconfcommitmentstore, bound to a specific deployed contract.
func NewPreconfcommitmentstoreFilterer(address common.Address, filterer bind.ContractFilterer) (*PreconfcommitmentstoreFilterer, error) {
	contract, err := bindPreconfcommitmentstore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreFilterer{contract: contract}, nil
}

// bindPreconfcommitmentstore binds a generic wrapper to an already deployed contract.
func bindPreconfcommitmentstore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreconfcommitmentstoreMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfcommitmentstore *PreconfcommitmentstoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.PreconfcommitmentstoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Preconfcommitmentstore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) DOMAINSEPARATORBID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "DOMAIN_SEPARATOR_BID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORBID(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORBID is a free data retrieval call binding the contract method 0x940b5765.
//
// Solidity: function DOMAIN_SEPARATOR_BID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) DOMAINSEPARATORBID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORBID(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) DOMAINSEPARATORPRECONF(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "DOMAIN_SEPARATOR_PRECONF")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORPRECONF(&_Preconfcommitmentstore.CallOpts)
}

// DOMAINSEPARATORPRECONF is a free data retrieval call binding the contract method 0xe5ae370f.
//
// Solidity: function DOMAIN_SEPARATOR_PRECONF() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) DOMAINSEPARATORPRECONF() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.DOMAINSEPARATORPRECONF(&_Preconfcommitmentstore.CallOpts)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) EIP712BIDTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "EIP712_BID_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712BIDTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712BIDTYPEHASH is a free data retrieval call binding the contract method 0x517aa8b7.
//
// Solidity: function EIP712_BID_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) EIP712BIDTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712BIDTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) EIP712COMMITMENTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "EIP712_COMMITMENT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// EIP712COMMITMENTTYPEHASH is a free data retrieval call binding the contract method 0x10ce6471.
//
// Solidity: function EIP712_COMMITMENT_TYPEHASH() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) EIP712COMMITMENTTYPEHASH() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.EIP712COMMITMENTTYPEHASH(&_Preconfcommitmentstore.CallOpts)
}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) HEXCHARS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "HEXCHARS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) HEXCHARS() ([]byte, error) {
	return _Preconfcommitmentstore.Contract.HEXCHARS(&_Preconfcommitmentstore.CallOpts)
}

// HEXCHARS is a free data retrieval call binding the contract method 0x05c51da6.
//
// Solidity: function HEXCHARS() view returns(bytes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) HEXCHARS() ([]byte, error) {
	return _Preconfcommitmentstore.Contract.HEXCHARS(&_Preconfcommitmentstore.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Preconfcommitmentstore.Contract.UPGRADEINTERFACEVERSION(&_Preconfcommitmentstore.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Preconfcommitmentstore.Contract.UPGRADEINTERFACEVERSION(&_Preconfcommitmentstore.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BidderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "bidderRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BidderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BidderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// BidderRegistry is a free data retrieval call binding the contract method 0x909e54e2.
//
// Solidity: function bidderRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BidderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BidderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BlockTracker(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "blockTracker")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BlockTracker() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BlockTracker(&_Preconfcommitmentstore.CallOpts)
}

// BlockTracker is a free data retrieval call binding the contract method 0x381c1d6c.
//
// Solidity: function blockTracker() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BlockTracker() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.BlockTracker(&_Preconfcommitmentstore.CallOpts)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) BlocksPerWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "blocksPerWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) BlocksPerWindow() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.BlocksPerWindow(&_Preconfcommitmentstore.CallOpts)
}

// BlocksPerWindow is a free data retrieval call binding the contract method 0x6347609e.
//
// Solidity: function blocksPerWindow() view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) BlocksPerWindow() (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.BlocksPerWindow(&_Preconfcommitmentstore.CallOpts)
}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) CommitmentDispatchWindow(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitmentDispatchWindow")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) CommitmentDispatchWindow() (uint64, error) {
	return _Preconfcommitmentstore.Contract.CommitmentDispatchWindow(&_Preconfcommitmentstore.CallOpts)
}

// CommitmentDispatchWindow is a free data retrieval call binding the contract method 0xf2357c03.
//
// Solidity: function commitmentDispatchWindow() view returns(uint64)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) CommitmentDispatchWindow() (uint64, error) {
	return _Preconfcommitmentstore.Contract.CommitmentDispatchWindow(&_Preconfcommitmentstore.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address bidder, bool isUsed, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bid, bytes32 bidHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey, string txnHash, string revertingTxHashes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) Commitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Bidder              common.Address
	IsUsed              bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	Bid                 *big.Int
	BidHash             [32]byte
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	SharedSecretKey     []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitments", arg0)

	outstruct := new(struct {
		Bidder              common.Address
		IsUsed              bool
		BlockNumber         uint64
		DecayStartTimeStamp uint64
		DecayEndTimeStamp   uint64
		DispatchTimestamp   uint64
		Committer           common.Address
		Bid                 *big.Int
		BidHash             [32]byte
		CommitmentHash      [32]byte
		BidSignature        []byte
		CommitmentSignature []byte
		SharedSecretKey     []byte
		TxnHash             string
		RevertingTxHashes   string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Bidder = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.IsUsed = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.BlockNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.DecayStartTimeStamp = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.DecayEndTimeStamp = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.DispatchTimestamp = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.Committer = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.Bid = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.BidHash = *abi.ConvertType(out[8], new([32]byte)).(*[32]byte)
	outstruct.CommitmentHash = *abi.ConvertType(out[9], new([32]byte)).(*[32]byte)
	outstruct.BidSignature = *abi.ConvertType(out[10], new([]byte)).(*[]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[11], new([]byte)).(*[]byte)
	outstruct.SharedSecretKey = *abi.ConvertType(out[12], new([]byte)).(*[]byte)
	outstruct.TxnHash = *abi.ConvertType(out[13], new(string)).(*string)
	outstruct.RevertingTxHashes = *abi.ConvertType(out[14], new(string)).(*string)

	return *outstruct, err

}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address bidder, bool isUsed, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bid, bytes32 bidHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey, string txnHash, string revertingTxHashes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Commitments(arg0 [32]byte) (struct {
	Bidder              common.Address
	IsUsed              bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	Bid                 *big.Int
	BidHash             [32]byte
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	SharedSecretKey     []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	return _Preconfcommitmentstore.Contract.Commitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x839df945.
//
// Solidity: function commitments(bytes32 ) view returns(address bidder, bool isUsed, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, uint64 dispatchTimestamp, address committer, uint256 bid, bytes32 bidHash, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey, string txnHash, string revertingTxHashes)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) Commitments(arg0 [32]byte) (struct {
	Bidder              common.Address
	IsUsed              bool
	BlockNumber         uint64
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	DispatchTimestamp   uint64
	Committer           common.Address
	Bid                 *big.Int
	BidHash             [32]byte
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	SharedSecretKey     []byte
	TxnHash             string
	RevertingTxHashes   string
}, error) {
	return _Preconfcommitmentstore.Contract.Commitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) CommitmentsCount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "commitmentsCount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentsCount(&_Preconfcommitmentstore.CallOpts, arg0)
}

// CommitmentsCount is a free data retrieval call binding the contract method 0x25f5cf21.
//
// Solidity: function commitmentsCount(address ) view returns(uint256)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) CommitmentsCount(arg0 common.Address) (*big.Int, error) {
	return _Preconfcommitmentstore.Contract.CommitmentsCount(&_Preconfcommitmentstore.CallOpts, arg0)
}

// EncryptedCommitments is a free data retrieval call binding the contract method 0x566f998c.
//
// Solidity: function encryptedCommitments(bytes32 ) view returns(bool isUsed, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) EncryptedCommitments(opts *bind.CallOpts, arg0 [32]byte) (struct {
	IsUsed              bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "encryptedCommitments", arg0)

	outstruct := new(struct {
		IsUsed              bool
		Committer           common.Address
		DispatchTimestamp   uint64
		CommitmentDigest    [32]byte
		CommitmentSignature []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsUsed = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Committer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.DispatchTimestamp = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.CommitmentDigest = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.CommitmentSignature = *abi.ConvertType(out[4], new([]byte)).(*[]byte)

	return *outstruct, err

}

// EncryptedCommitments is a free data retrieval call binding the contract method 0x566f998c.
//
// Solidity: function encryptedCommitments(bytes32 ) view returns(bool isUsed, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) EncryptedCommitments(arg0 [32]byte) (struct {
	IsUsed              bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	return _Preconfcommitmentstore.Contract.EncryptedCommitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// EncryptedCommitments is a free data retrieval call binding the contract method 0x566f998c.
//
// Solidity: function encryptedCommitments(bytes32 ) view returns(bool isUsed, address committer, uint64 dispatchTimestamp, bytes32 commitmentDigest, bytes commitmentSignature)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) EncryptedCommitments(arg0 [32]byte) (struct {
	IsUsed              bool
	Committer           common.Address
	DispatchTimestamp   uint64
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
}, error) {
	return _Preconfcommitmentstore.Contract.EncryptedCommitments(&_Preconfcommitmentstore.CallOpts, arg0)
}

// GetBidHash is a free data retrieval call binding the contract method 0xdbd007ab.
//
// Solidity: function getBidHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetBidHash(opts *bind.CallOpts, _txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getBidHash", _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBidHash is a free data retrieval call binding the contract method 0xdbd007ab.
//
// Solidity: function getBidHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetBidHash(_txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetBidHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)
}

// GetBidHash is a free data retrieval call binding the contract method 0xdbd007ab.
//
// Solidity: function getBidHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetBidHash(_txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetBidHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (IPreConfCommitmentStorePreConfCommitment, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitment", commitmentIndex)

	if err != nil {
		return *new(IPreConfCommitmentStorePreConfCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(IPreConfCommitmentStorePreConfCommitment)).(*IPreConfCommitmentStorePreConfCommitment)

	return out0, err

}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string))
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitment(commitmentIndex [32]byte) (IPreConfCommitmentStorePreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetCommitment is a free data retrieval call binding the contract method 0x7795820c.
//
// Solidity: function getCommitment(bytes32 commitmentIndex) view returns((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitment(commitmentIndex [32]byte) (IPreConfCommitmentStorePreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xb2200dbc.
//
// Solidity: function getCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetCommitmentIndex(opts *bind.CallOpts, commitment IPreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xb2200dbc.
//
// Solidity: function getCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetCommitmentIndex(commitment IPreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetCommitmentIndex is a free data retrieval call binding the contract method 0xb2200dbc.
//
// Solidity: function getCommitmentIndex((address,bool,uint64,uint64,uint64,uint64,address,uint256,bytes32,bytes32,bytes,bytes,bytes,string,string) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetCommitmentIndex(commitment IPreConfCommitmentStorePreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetEncryptedCommitment is a free data retrieval call binding the contract method 0x1725b4a7.
//
// Solidity: function getEncryptedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetEncryptedCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (IPreConfCommitmentStoreEncrPreConfCommitment, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getEncryptedCommitment", commitmentIndex)

	if err != nil {
		return *new(IPreConfCommitmentStoreEncrPreConfCommitment), err
	}

	out0 := *abi.ConvertType(out[0], new(IPreConfCommitmentStoreEncrPreConfCommitment)).(*IPreConfCommitmentStoreEncrPreConfCommitment)

	return out0, err

}

// GetEncryptedCommitment is a free data retrieval call binding the contract method 0x1725b4a7.
//
// Solidity: function getEncryptedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetEncryptedCommitment(commitmentIndex [32]byte) (IPreConfCommitmentStoreEncrPreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetEncryptedCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetEncryptedCommitment is a free data retrieval call binding the contract method 0x1725b4a7.
//
// Solidity: function getEncryptedCommitment(bytes32 commitmentIndex) view returns((bool,address,uint64,bytes32,bytes))
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetEncryptedCommitment(commitmentIndex [32]byte) (IPreConfCommitmentStoreEncrPreConfCommitment, error) {
	return _Preconfcommitmentstore.Contract.GetEncryptedCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetEncryptedCommitmentIndex is a free data retrieval call binding the contract method 0xfd25c65a.
//
// Solidity: function getEncryptedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetEncryptedCommitmentIndex(opts *bind.CallOpts, commitment IPreConfCommitmentStoreEncrPreConfCommitment) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getEncryptedCommitmentIndex", commitment)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetEncryptedCommitmentIndex is a free data retrieval call binding the contract method 0xfd25c65a.
//
// Solidity: function getEncryptedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetEncryptedCommitmentIndex(commitment IPreConfCommitmentStoreEncrPreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetEncryptedCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetEncryptedCommitmentIndex is a free data retrieval call binding the contract method 0xfd25c65a.
//
// Solidity: function getEncryptedCommitmentIndex((bool,address,uint64,bytes32,bytes) commitment) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetEncryptedCommitmentIndex(commitment IPreConfCommitmentStoreEncrPreConfCommitment) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetEncryptedCommitmentIndex(&_Preconfcommitmentstore.CallOpts, commitment)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xd3dff2ae.
//
// Solidity: function getPreConfHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature, string _sharedSecretKey) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetPreConfHash(opts *bind.CallOpts, _txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string, _sharedSecretKey string) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getPreConfHash", _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature, _sharedSecretKey)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPreConfHash is a free data retrieval call binding the contract method 0xd3dff2ae.
//
// Solidity: function getPreConfHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature, string _sharedSecretKey) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetPreConfHash(_txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string, _sharedSecretKey string) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetPreConfHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature, _sharedSecretKey)
}

// GetPreConfHash is a free data retrieval call binding the contract method 0xd3dff2ae.
//
// Solidity: function getPreConfHash(string _txnHash, string _revertingTxHashes, uint256 _bid, uint64 _blockNumber, uint64 _decayStartTimeStamp, uint64 _decayEndTimeStamp, bytes32 _bidHash, string _bidSignature, string _sharedSecretKey) pure returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetPreConfHash(_txnHash string, _revertingTxHashes string, _bid *big.Int, _blockNumber uint64, _decayStartTimeStamp uint64, _decayEndTimeStamp uint64, _bidHash [32]byte, _bidSignature string, _sharedSecretKey string) ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.GetPreConfHash(&_Preconfcommitmentstore.CallOpts, _txnHash, _revertingTxHashes, _bid, _blockNumber, _decayStartTimeStamp, _decayEndTimeStamp, _bidHash, _bidSignature, _sharedSecretKey)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "getTxnHashFromCommitment", commitmentIndex)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfcommitmentstore.Contract.GetTxnHashFromCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// GetTxnHashFromCommitment is a free data retrieval call binding the contract method 0xfc4fbe32.
//
// Solidity: function getTxnHashFromCommitment(bytes32 commitmentIndex) view returns(string txnHash)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) GetTxnHashFromCommitment(commitmentIndex [32]byte) (string, error) {
	return _Preconfcommitmentstore.Contract.GetTxnHashFromCommitment(&_Preconfcommitmentstore.CallOpts, commitmentIndex)
}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) OracleContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "oracleContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) OracleContract() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.OracleContract(&_Preconfcommitmentstore.CallOpts)
}

// OracleContract is a free data retrieval call binding the contract method 0xbece7532.
//
// Solidity: function oracleContract() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) OracleContract() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.OracleContract(&_Preconfcommitmentstore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Owner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Owner(&_Preconfcommitmentstore.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) Owner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.Owner(&_Preconfcommitmentstore.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) PendingOwner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.PendingOwner(&_Preconfcommitmentstore.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) PendingOwner() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.PendingOwner(&_Preconfcommitmentstore.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) ProviderRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "providerRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) ProviderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.ProviderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// ProviderRegistry is a free data retrieval call binding the contract method 0x545921d9.
//
// Solidity: function providerRegistry() view returns(address)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) ProviderRegistry() (common.Address, error) {
	return _Preconfcommitmentstore.Contract.ProviderRegistry(&_Preconfcommitmentstore.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) ProxiableUUID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.ProxiableUUID(&_Preconfcommitmentstore.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Preconfcommitmentstore.Contract.ProxiableUUID(&_Preconfcommitmentstore.CallOpts)
}

// VerifyBid is a free data retrieval call binding the contract method 0x20ee734c.
//
// Solidity: function verifyBid(uint256 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes bidSignature) pure returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) VerifyBid(opts *bind.CallOpts, bid *big.Int, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, revertingTxHashes string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "verifyBid", bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, revertingTxHashes, bidSignature)

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

// VerifyBid is a free data retrieval call binding the contract method 0x20ee734c.
//
// Solidity: function verifyBid(uint256 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes bidSignature) pure returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) VerifyBid(bid *big.Int, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, revertingTxHashes string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyBid(&_Preconfcommitmentstore.CallOpts, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, revertingTxHashes, bidSignature)
}

// VerifyBid is a free data retrieval call binding the contract method 0x20ee734c.
//
// Solidity: function verifyBid(uint256 bid, uint64 blockNumber, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes bidSignature) pure returns(bytes32 messageDigest, address recoveredAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) VerifyBid(bid *big.Int, blockNumber uint64, decayStartTimeStamp uint64, decayEndTimeStamp uint64, txnHash string, revertingTxHashes string, bidSignature []byte) (struct {
	MessageDigest    [32]byte
	RecoveredAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyBid(&_Preconfcommitmentstore.CallOpts, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp, txnHash, revertingTxHashes, bidSignature)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xe93bbf26.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,bytes) params) pure returns(bytes32 preConfHash, address committerAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCaller) VerifyPreConfCommitment(opts *bind.CallOpts, params IPreConfCommitmentStoreCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	var out []interface{}
	err := _Preconfcommitmentstore.contract.Call(opts, &out, "verifyPreConfCommitment", params)

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

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xe93bbf26.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,bytes) params) pure returns(bytes32 preConfHash, address committerAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) VerifyPreConfCommitment(params IPreConfCommitmentStoreCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyPreConfCommitment(&_Preconfcommitmentstore.CallOpts, params)
}

// VerifyPreConfCommitment is a free data retrieval call binding the contract method 0xe93bbf26.
//
// Solidity: function verifyPreConfCommitment((string,string,uint256,uint64,uint64,uint64,bytes32,bytes,bytes,bytes) params) pure returns(bytes32 preConfHash, address committerAddress)
func (_Preconfcommitmentstore *PreconfcommitmentstoreCallerSession) VerifyPreConfCommitment(params IPreConfCommitmentStoreCommitmentParams) (struct {
	PreConfHash      [32]byte
	CommitterAddress common.Address
}, error) {
	return _Preconfcommitmentstore.Contract.VerifyPreConfCommitment(&_Preconfcommitmentstore.CallOpts, params)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) AcceptOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.AcceptOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.AcceptOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xd8d0cbc1.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow, uint256 _blocksPerWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) Initialize(opts *bind.TransactOpts, _providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64, _blocksPerWindow *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "initialize", _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow, _blocksPerWindow)
}

// Initialize is a paid mutator transaction binding the contract method 0xd8d0cbc1.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow, uint256 _blocksPerWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Initialize(_providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64, _blocksPerWindow *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Initialize(&_Preconfcommitmentstore.TransactOpts, _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow, _blocksPerWindow)
}

// Initialize is a paid mutator transaction binding the contract method 0xd8d0cbc1.
//
// Solidity: function initialize(address _providerRegistry, address _bidderRegistry, address _oracleContract, address _owner, address _blockTracker, uint64 _commitmentDispatchWindow, uint256 _blocksPerWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) Initialize(_providerRegistry common.Address, _bidderRegistry common.Address, _oracleContract common.Address, _owner common.Address, _blockTracker common.Address, _commitmentDispatchWindow uint64, _blocksPerWindow *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Initialize(&_Preconfcommitmentstore.TransactOpts, _providerRegistry, _bidderRegistry, _oracleContract, _owner, _blockTracker, _commitmentDispatchWindow, _blocksPerWindow)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) InitiateReward(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "initiateReward", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateReward(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateReward is a paid mutator transaction binding the contract method 0x03faf979.
//
// Solidity: function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) InitiateReward(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateReward(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) InitiateSlash(opts *bind.TransactOpts, commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "initiateSlash", commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateSlash(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// InitiateSlash is a paid mutator transaction binding the contract method 0x30778c78.
//
// Solidity: function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) InitiateSlash(commitmentIndex [32]byte, residualBidPercentAfterDecay *big.Int) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.InitiateSlash(&_Preconfcommitmentstore.TransactOpts, commitmentIndex, residualBidPercentAfterDecay)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x05880c6a.
//
// Solidity: function openCommitment(bytes32 encryptedCommitmentIndex, uint256 bid, uint64 blockNumber, string txnHash, string revertingTxHashes, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) OpenCommitment(opts *bind.TransactOpts, encryptedCommitmentIndex [32]byte, bid *big.Int, blockNumber uint64, txnHash string, revertingTxHashes string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte, sharedSecretKey []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "openCommitment", encryptedCommitmentIndex, bid, blockNumber, txnHash, revertingTxHashes, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature, sharedSecretKey)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x05880c6a.
//
// Solidity: function openCommitment(bytes32 encryptedCommitmentIndex, uint256 bid, uint64 blockNumber, string txnHash, string revertingTxHashes, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) OpenCommitment(encryptedCommitmentIndex [32]byte, bid *big.Int, blockNumber uint64, txnHash string, revertingTxHashes string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte, sharedSecretKey []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.OpenCommitment(&_Preconfcommitmentstore.TransactOpts, encryptedCommitmentIndex, bid, blockNumber, txnHash, revertingTxHashes, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature, sharedSecretKey)
}

// OpenCommitment is a paid mutator transaction binding the contract method 0x05880c6a.
//
// Solidity: function openCommitment(bytes32 encryptedCommitmentIndex, uint256 bid, uint64 blockNumber, string txnHash, string revertingTxHashes, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, bytes bidSignature, bytes commitmentSignature, bytes sharedSecretKey) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) OpenCommitment(encryptedCommitmentIndex [32]byte, bid *big.Int, blockNumber uint64, txnHash string, revertingTxHashes string, decayStartTimeStamp uint64, decayEndTimeStamp uint64, bidSignature []byte, commitmentSignature []byte, sharedSecretKey []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.OpenCommitment(&_Preconfcommitmentstore.TransactOpts, encryptedCommitmentIndex, bid, blockNumber, txnHash, revertingTxHashes, decayStartTimeStamp, decayEndTimeStamp, bidSignature, commitmentSignature, sharedSecretKey)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.RenounceOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.RenounceOwnership(&_Preconfcommitmentstore.TransactOpts)
}

// StoreEncryptedCommitment is a paid mutator transaction binding the contract method 0x88b74730.
//
// Solidity: function storeEncryptedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) StoreEncryptedCommitment(opts *bind.TransactOpts, commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "storeEncryptedCommitment", commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// StoreEncryptedCommitment is a paid mutator transaction binding the contract method 0x88b74730.
//
// Solidity: function storeEncryptedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) StoreEncryptedCommitment(commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.StoreEncryptedCommitment(&_Preconfcommitmentstore.TransactOpts, commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// StoreEncryptedCommitment is a paid mutator transaction binding the contract method 0x88b74730.
//
// Solidity: function storeEncryptedCommitment(bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp) returns(bytes32 commitmentIndex)
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) StoreEncryptedCommitment(commitmentDigest [32]byte, commitmentSignature []byte, dispatchTimestamp uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.StoreEncryptedCommitment(&_Preconfcommitmentstore.TransactOpts, commitmentDigest, commitmentSignature, dispatchTimestamp)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.TransferOwnership(&_Preconfcommitmentstore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.TransferOwnership(&_Preconfcommitmentstore.TransactOpts, newOwner)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateBidderRegistry(opts *bind.TransactOpts, newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateBidderRegistry", newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateBidderRegistry(&_Preconfcommitmentstore.TransactOpts, newBidderRegistry)
}

// UpdateBidderRegistry is a paid mutator transaction binding the contract method 0x66544c41.
//
// Solidity: function updateBidderRegistry(address newBidderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateBidderRegistry(newBidderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateBidderRegistry(&_Preconfcommitmentstore.TransactOpts, newBidderRegistry)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateCommitmentDispatchWindow(opts *bind.TransactOpts, newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateCommitmentDispatchWindow", newDispatchWindow)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateCommitmentDispatchWindow(newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateCommitmentDispatchWindow(&_Preconfcommitmentstore.TransactOpts, newDispatchWindow)
}

// UpdateCommitmentDispatchWindow is a paid mutator transaction binding the contract method 0x03800560.
//
// Solidity: function updateCommitmentDispatchWindow(uint64 newDispatchWindow) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateCommitmentDispatchWindow(newDispatchWindow uint64) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateCommitmentDispatchWindow(&_Preconfcommitmentstore.TransactOpts, newDispatchWindow)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateOracleContract(opts *bind.TransactOpts, newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateOracleContract", newOracleContract)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateOracleContract(newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateOracleContract(&_Preconfcommitmentstore.TransactOpts, newOracleContract)
}

// UpdateOracleContract is a paid mutator transaction binding the contract method 0xd3bab58f.
//
// Solidity: function updateOracleContract(address newOracleContract) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateOracleContract(newOracleContract common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateOracleContract(&_Preconfcommitmentstore.TransactOpts, newOracleContract)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpdateProviderRegistry(opts *bind.TransactOpts, newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "updateProviderRegistry", newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateProviderRegistry(&_Preconfcommitmentstore.TransactOpts, newProviderRegistry)
}

// UpdateProviderRegistry is a paid mutator transaction binding the contract method 0x92d2e3e7.
//
// Solidity: function updateProviderRegistry(address newProviderRegistry) returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpdateProviderRegistry(newProviderRegistry common.Address) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpdateProviderRegistry(&_Preconfcommitmentstore.TransactOpts, newProviderRegistry)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpgradeToAndCall(&_Preconfcommitmentstore.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.UpgradeToAndCall(&_Preconfcommitmentstore.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Fallback(&_Preconfcommitmentstore.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Fallback(&_Preconfcommitmentstore.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Preconfcommitmentstore.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreSession) Receive() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Receive(&_Preconfcommitmentstore.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Preconfcommitmentstore *PreconfcommitmentstoreTransactorSession) Receive() (*types.Transaction, error) {
	return _Preconfcommitmentstore.Contract.Receive(&_Preconfcommitmentstore.TransactOpts)
}

// PreconfcommitmentstoreCommitmentStoredIterator is returned from FilterCommitmentStored and is used to iterate over the raw logs and unpacked data for CommitmentStored events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreCommitmentStoredIterator struct {
	Event *PreconfcommitmentstoreCommitmentStored // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreCommitmentStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreCommitmentStored)
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
		it.Event = new(PreconfcommitmentstoreCommitmentStored)
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
func (it *PreconfcommitmentstoreCommitmentStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreCommitmentStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreCommitmentStored represents a CommitmentStored event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreCommitmentStored struct {
	CommitmentIndex     [32]byte
	Bidder              common.Address
	Committer           common.Address
	Bid                 *big.Int
	BlockNumber         uint64
	BidHash             [32]byte
	DecayStartTimeStamp uint64
	DecayEndTimeStamp   uint64
	TxnHash             string
	RevertingTxHashes   string
	CommitmentHash      [32]byte
	BidSignature        []byte
	CommitmentSignature []byte
	DispatchTimestamp   uint64
	SharedSecretKey     []byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCommitmentStored is a free log retrieval operation binding the contract event 0xe44dd4d002deb2c79cf08ce285a9d80c69753f31ca65c8e49f0a60d27ed9fea3.
//
// Solidity: event CommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint64 dispatchTimestamp, bytes sharedSecretKey)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterCommitmentStored(opts *bind.FilterOpts, commitmentIndex [][32]byte) (*PreconfcommitmentstoreCommitmentStoredIterator, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "CommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreCommitmentStoredIterator{contract: _Preconfcommitmentstore.contract, event: "CommitmentStored", logs: logs, sub: sub}, nil
}

// WatchCommitmentStored is a free log subscription operation binding the contract event 0xe44dd4d002deb2c79cf08ce285a9d80c69753f31ca65c8e49f0a60d27ed9fea3.
//
// Solidity: event CommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint64 dispatchTimestamp, bytes sharedSecretKey)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchCommitmentStored(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreCommitmentStored, commitmentIndex [][32]byte) (event.Subscription, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "CommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreCommitmentStored)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "CommitmentStored", log); err != nil {
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

// ParseCommitmentStored is a log parse operation binding the contract event 0xe44dd4d002deb2c79cf08ce285a9d80c69753f31ca65c8e49f0a60d27ed9fea3.
//
// Solidity: event CommitmentStored(bytes32 indexed commitmentIndex, address bidder, address committer, uint256 bid, uint64 blockNumber, bytes32 bidHash, uint64 decayStartTimeStamp, uint64 decayEndTimeStamp, string txnHash, string revertingTxHashes, bytes32 commitmentHash, bytes bidSignature, bytes commitmentSignature, uint64 dispatchTimestamp, bytes sharedSecretKey)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseCommitmentStored(log types.Log) (*PreconfcommitmentstoreCommitmentStored, error) {
	event := new(PreconfcommitmentstoreCommitmentStored)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "CommitmentStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreEncryptedCommitmentStoredIterator is returned from FilterEncryptedCommitmentStored and is used to iterate over the raw logs and unpacked data for EncryptedCommitmentStored events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreEncryptedCommitmentStoredIterator struct {
	Event *PreconfcommitmentstoreEncryptedCommitmentStored // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreEncryptedCommitmentStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreEncryptedCommitmentStored)
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
		it.Event = new(PreconfcommitmentstoreEncryptedCommitmentStored)
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
func (it *PreconfcommitmentstoreEncryptedCommitmentStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreEncryptedCommitmentStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreEncryptedCommitmentStored represents a EncryptedCommitmentStored event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreEncryptedCommitmentStored struct {
	CommitmentIndex     [32]byte
	Committer           common.Address
	CommitmentDigest    [32]byte
	CommitmentSignature []byte
	DispatchTimestamp   uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterEncryptedCommitmentStored is a free log retrieval operation binding the contract event 0x3b9ebc34b9c72a41ddaf26db939c901a27a144dfbebbe80c3b105a7684a617f2.
//
// Solidity: event EncryptedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterEncryptedCommitmentStored(opts *bind.FilterOpts, commitmentIndex [][32]byte) (*PreconfcommitmentstoreEncryptedCommitmentStoredIterator, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "EncryptedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreEncryptedCommitmentStoredIterator{contract: _Preconfcommitmentstore.contract, event: "EncryptedCommitmentStored", logs: logs, sub: sub}, nil
}

// WatchEncryptedCommitmentStored is a free log subscription operation binding the contract event 0x3b9ebc34b9c72a41ddaf26db939c901a27a144dfbebbe80c3b105a7684a617f2.
//
// Solidity: event EncryptedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchEncryptedCommitmentStored(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreEncryptedCommitmentStored, commitmentIndex [][32]byte) (event.Subscription, error) {

	var commitmentIndexRule []interface{}
	for _, commitmentIndexItem := range commitmentIndex {
		commitmentIndexRule = append(commitmentIndexRule, commitmentIndexItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "EncryptedCommitmentStored", commitmentIndexRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreEncryptedCommitmentStored)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "EncryptedCommitmentStored", log); err != nil {
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

// ParseEncryptedCommitmentStored is a log parse operation binding the contract event 0x3b9ebc34b9c72a41ddaf26db939c901a27a144dfbebbe80c3b105a7684a617f2.
//
// Solidity: event EncryptedCommitmentStored(bytes32 indexed commitmentIndex, address committer, bytes32 commitmentDigest, bytes commitmentSignature, uint64 dispatchTimestamp)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseEncryptedCommitmentStored(log types.Log) (*PreconfcommitmentstoreEncryptedCommitmentStored, error) {
	event := new(PreconfcommitmentstoreEncryptedCommitmentStored)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "EncryptedCommitmentStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreInitializedIterator struct {
	Event *PreconfcommitmentstoreInitialized // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreInitialized)
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
		it.Event = new(PreconfcommitmentstoreInitialized)
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
func (it *PreconfcommitmentstoreInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreInitialized represents a Initialized event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterInitialized(opts *bind.FilterOpts) (*PreconfcommitmentstoreInitializedIterator, error) {

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreInitializedIterator{contract: _Preconfcommitmentstore.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreInitialized) (event.Subscription, error) {

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreInitialized)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseInitialized(log types.Log) (*PreconfcommitmentstoreInitialized, error) {
	event := new(PreconfcommitmentstoreInitialized)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferStartedIterator struct {
	Event *PreconfcommitmentstoreOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreOwnershipTransferStarted)
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
		it.Event = new(PreconfcommitmentstoreOwnershipTransferStarted)
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
func (it *PreconfcommitmentstoreOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfcommitmentstoreOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreOwnershipTransferStartedIterator{contract: _Preconfcommitmentstore.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreOwnershipTransferStarted)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseOwnershipTransferStarted(log types.Log) (*PreconfcommitmentstoreOwnershipTransferStarted, error) {
	event := new(PreconfcommitmentstoreOwnershipTransferStarted)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferredIterator struct {
	Event *PreconfcommitmentstoreOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreOwnershipTransferred)
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
		it.Event = new(PreconfcommitmentstoreOwnershipTransferred)
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
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreOwnershipTransferred represents a OwnershipTransferred event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreconfcommitmentstoreOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreOwnershipTransferredIterator{contract: _Preconfcommitmentstore.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreOwnershipTransferred)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseOwnershipTransferred(log types.Log) (*PreconfcommitmentstoreOwnershipTransferred, error) {
	event := new(PreconfcommitmentstoreOwnershipTransferred)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreSignatureVerifiedIterator is returned from FilterSignatureVerified and is used to iterate over the raw logs and unpacked data for SignatureVerified events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreSignatureVerifiedIterator struct {
	Event *PreconfcommitmentstoreSignatureVerified // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreSignatureVerified)
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
		it.Event = new(PreconfcommitmentstoreSignatureVerified)
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
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreSignatureVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreSignatureVerified represents a SignatureVerified event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreSignatureVerified struct {
	Signer            common.Address
	TxnHash           string
	RevertingTxHashes string
	Bid               *big.Int
	BlockNumber       uint64
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSignatureVerified is a free log retrieval operation binding the contract event 0xf5fa3f7902916f0e4e2af218fa9a77809c54e64b77ec4061a205ec3b7ce5c35e.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterSignatureVerified(opts *bind.FilterOpts, signer []common.Address, bid []*big.Int) (*PreconfcommitmentstoreSignatureVerifiedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreSignatureVerifiedIterator{contract: _Preconfcommitmentstore.contract, event: "SignatureVerified", logs: logs, sub: sub}, nil
}

// WatchSignatureVerified is a free log subscription operation binding the contract event 0xf5fa3f7902916f0e4e2af218fa9a77809c54e64b77ec4061a205ec3b7ce5c35e.
//
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchSignatureVerified(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreSignatureVerified, signer []common.Address, bid []*big.Int) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	var bidRule []interface{}
	for _, bidItem := range bid {
		bidRule = append(bidRule, bidItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "SignatureVerified", signerRule, bidRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreSignatureVerified)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
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
// Solidity: event SignatureVerified(address indexed signer, string txnHash, string revertingTxHashes, uint256 indexed bid, uint64 blockNumber)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseSignatureVerified(log types.Log) (*PreconfcommitmentstoreSignatureVerified, error) {
	event := new(PreconfcommitmentstoreSignatureVerified)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "SignatureVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreconfcommitmentstoreUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreUpgradedIterator struct {
	Event *PreconfcommitmentstoreUpgraded // Event containing the contract specifics and raw log

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
func (it *PreconfcommitmentstoreUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreconfcommitmentstoreUpgraded)
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
		it.Event = new(PreconfcommitmentstoreUpgraded)
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
func (it *PreconfcommitmentstoreUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreconfcommitmentstoreUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreconfcommitmentstoreUpgraded represents a Upgraded event raised by the Preconfcommitmentstore contract.
type PreconfcommitmentstoreUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PreconfcommitmentstoreUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PreconfcommitmentstoreUpgradedIterator{contract: _Preconfcommitmentstore.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PreconfcommitmentstoreUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Preconfcommitmentstore.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreconfcommitmentstoreUpgraded)
				if err := _Preconfcommitmentstore.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Preconfcommitmentstore *PreconfcommitmentstoreFilterer) ParseUpgraded(log types.Log) (*PreconfcommitmentstoreUpgraded, error) {
	event := new(PreconfcommitmentstoreUpgraded)
	if err := _Preconfcommitmentstore.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
