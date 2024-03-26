// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

/**
 * @title IPreConfCommitmentStore
 * @dev Interface for PreConfCommitmentStore
 */
interface IPreConfCommitmentStore {
    // Structs, events, and errors can also be included in the interface if they are used in the external functions

    /// @dev Struct for all the information around preconfirmations commitment
    struct PreConfCommitment {
        bool commitmentUsed;
        address bidder;
        address commiter;
        uint64 bid;
        uint64 blockNumber;
        bytes32 bidHash;
        uint64 decayStartTimeStamp;
        uint64 decayEndTimeStamp;
        string txnHash;
        bytes32 commitmentHash;
        bytes bidSignature;
        bytes commitmentSignature;
        uint256 blockCommitedAt;
    }



    event SignatureVerified(
        address indexed signer,
        string txnHash,
        uint64 indexed bid,
        uint64 blockNumber
    );

    // External functions that need to be implemented

    function getBidHash(
        string memory _txnHash,
        uint64 _bid,
        uint64 _blockNumber
    ) external view returns (bytes32);

    function getPreConfHash(
        string memory _txnHash,
        uint64 _bid,
        uint64 _blockNumber,
        bytes32 _bidHash,
        string memory _bidSignature
    ) external view returns (bytes32);

    function retreiveCommitments() external view returns (PreConfCommitment[] memory);

    function retreiveCommitment() external view returns (PreConfCommitment memory);

    function verifyBid(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        bytes calldata bidSignature
    ) external view returns (bytes32 messageDigest, address recoveredAddress, uint256 stake);

    function storeCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory commitmentHash,
        bytes calldata bidSignature,
        bytes memory commitmentSignature
    ) external returns (uint256);

    function getCommitmentsByBlockNumber(uint256 blockNumber) external view returns (bytes32[] memory);


    function getCommitment(bytes32 commitmentIndex) external view returns (PreConfCommitment memory);

    function initiateSlash(bytes32 commitmentIndex, uint256 residualDecayedBid) external;

    function initiateReward(bytes32 commitmentIndex, uint256 residualDecayedBid) external;
    
    function unlockBidFunds(bytes32 commitmentDigest) external;

    function updateOracle(address newOracle) external;

    function updateProviderRegistry(address newProviderRegistry) external;

    function updateBidderRegistry(address newBidderRegistry) external;

    // Public functions that can be included if they are meant to be called from other contracts

    function _bytes32ToHexString(bytes32 _bytes32) external pure returns (string memory);

    function _bytesToHexString(bytes memory _bytes) external pure returns (string memory);
}
