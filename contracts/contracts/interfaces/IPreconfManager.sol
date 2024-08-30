// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

/**
 * @title IPreconfManager
 * @dev Interface for PreconfManager
 */
interface IPreconfManager {
    /// @dev Struct for all the information around preconfirmations commitment
    struct OpenedCommitment {
        address bidder;
        bool isSettled; // Flag to check if the commitment is settled with slashing or rewarding
        uint64 blockNumber;
        uint64 decayStartTimeStamp;
        uint64 decayEndTimeStamp;
        uint64 dispatchTimestamp;
        address committer;
        uint256 bid;
        bytes32 bidHash;
        bytes32 commitmentDigest;
        bytes bidSignature;
        bytes commitmentSignature;
        bytes sharedSecretKey;
        string txnHash;
        string revertingTxHashes;
    }

    /// @dev Struct for all the commitment params to avoid too deep in the stack error
    struct CommitmentParams {
        string txnHash;
        string revertingTxHashes;
        uint256 bid;
        uint64 blockNumber;
        uint64 decayStartTimeStamp;
        uint64 decayEndTimeStamp;
        bytes32 bidHash;
        bytes bidSignature;
        bytes commitmentSignature;
        bytes sharedSecretKey;
    }

    /// @dev Struct for all the information around unopened preconfirmations commitment
    struct UnopenedCommitment {
        bool isOpened; // Flag to check if the commitment is opened already
        address committer;
        uint64 dispatchTimestamp;
        bytes32 commitmentDigest;
        bytes commitmentSignature;
    }

    /// @dev Event to log successful opened commitment storage
    event OpenedCommitmentStored(
        bytes32 indexed commitmentIndex,
        address bidder,
        address committer,
        uint256 bid,
        uint64 blockNumber,
        bytes32 bidHash,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        string txnHash,
        string revertingTxHashes,
        bytes32 commitmentDigest,
        bytes bidSignature,
        bytes commitmentSignature,
        uint64 dispatchTimestamp,
        bytes sharedSecretKey
    );

    /// @dev Event to log successful unopened commitment storage
    event UnopenedCommitmentStored(
        bytes32 indexed commitmentIndex,
        address committer,
        bytes32 commitmentDigest,
        bytes commitmentSignature,
        uint64 dispatchTimestamp
    );

    /// @dev Event to log successful verifications
    event SignatureVerified(
        address indexed signer,
        string txnHash,
        string revertingTxHashes,
        uint256 indexed bid,
        uint64 blockNumber
    );

    /// @dev Error if dispatch timestamp is invalid
    error InvalidDispatchTimestamp(uint256 minTime, uint64 dispatchTimestamp);

    /// @dev Error if decay parameters are invalid
    error InvalidDecayTime(uint64 startTime, uint64 endTime);

    /// @dev Error if commitment is already opened
    error CommitmentAlreadyOpened(bytes32 commitmentIndex);

    /// @dev Error if commitment index is invalid
    error InvalidCommitmentDigest(bytes32 commitmentDigest, bytes32 computedDigest);

    /// @dev Error if commitment is not by the winner
    error WinnerIsNotCommitter(address committer, address winner);

    /// @dev Error if commitment is not opened by the committer or the bidder
    error UnauthorizedOpenCommitment(address committer, address bidder, address sender);

    /// @dev Error if encrypted commitment is sent by the committer
    error SenderIsNotCommitter(address expected, address actual);

    /**
     * @dev Initializes the contract with the specified registry addresses, oracle, name, and version.
     * @param _providerRegistry The address of the provider registry.
     * @param _bidderRegistry The address of the bidder registry.
     * @param _blockTracker The address of the block tracker.
     * @param _oracle The address of the oracle.
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _commitmentDispatchWindow The dispatch window for commitments.
     * @param _blocksPerWindow The number of blocks per window.
     */
    function initialize(
        address _providerRegistry,
        address _bidderRegistry,
        address _oracle,
        address _owner,
        address _blockTracker,
        uint64 _commitmentDispatchWindow,
        uint256 _blocksPerWindow
    ) external;

    /**
     * @dev Updates the commitment dispatch window to a new value.
     * @param newDispatchWindow The new dispatch window value to be set.
     */
    function updateCommitmentDispatchWindow(uint64 newDispatchWindow) external;

    /**
     * @dev Updates the address of the oracle contract.
     * @param newOracleContract The new oracle contract address.
     */
    function updateOracleContract(address newOracleContract) external;

    /**
     * @dev Updates the address of the bidder registry.
     * @param newBidderRegistry The new bidder registry address.
     */
    function updateBidderRegistry(address newBidderRegistry) external;

    /**
     * @dev Updates the address of the provider registry.
     * @param newProviderRegistry The new provider registry address.
     */
    function updateProviderRegistry(address newProviderRegistry) external;

    /**
     * @dev Opens a commitment.
     * @param unopenedCommitmentIndex The index of the unopened commitment.
     * @param bid The bid amount.
     * @param blockNumber The block number.
     * @param txnHash The transaction hash.
     * @param revertingTxHashes The reverting transaction hashes.
     * @param decayStartTimeStamp The start time of the decay.
     * @param decayEndTimeStamp The end time of the decay.
     * @param bidSignature The signature of the bid.
     * @param commitmentSignature The signature of the commitment.
     * @param sharedSecretKey The shared secret key.
     * @return commitmentIndex The index of the stored commitment.
     */
    function openCommitment(
        bytes32 unopenedCommitmentIndex,
        uint256 bid,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        bytes calldata bidSignature,
        bytes memory commitmentSignature,
        bytes memory sharedSecretKey
    ) external returns (bytes32 commitmentIndex);

    /**
     * @dev Stores an unopened commitment.
     * @param commitmentDigest The digest of the commitment.
     * @param commitmentSignature The signature of the commitment.
     * @param dispatchTimestamp The timestamp at which the commitment is dispatched.
     * @return commitmentIndex The index of the stored commitment.
     */
    function storeUnopenedCommitment(
        bytes32 commitmentDigest,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp
    ) external returns (bytes32 commitmentIndex);

    /**
     * @dev Initiates a slash for a commitment.
     * @param commitmentIndex The hash of the commitment to be slashed.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function initiateSlash(
        bytes32 commitmentIndex,
        uint256 residualBidPercentAfterDecay
    ) external;

    /**
     * @dev Initiates a reward for a commitment.
     * @param commitmentIndex The hash of the commitment to be rewarded.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function initiateReward(
        bytes32 commitmentIndex,
        uint256 residualBidPercentAfterDecay
    ) external;

    /**
     * @dev Gets the transaction hash from a commitment.
     * @param commitmentIndex The index of the commitment.
     * @return txnHash The transaction hash.
     */
    function getTxnHashFromCommitment(
        bytes32 commitmentIndex
    ) external view returns (string memory txnHash);

    /**
     * @dev Gets a commitment by its index.
     * @param commitmentIndex The index of the commitment.
     * @return A OpenedCommitment structure representing the commitment.
     */
    function getCommitment(
        bytes32 commitmentIndex
    ) external view returns (OpenedCommitment memory);

    /**
     * @dev Gets an unopened commitment by its index.
     * @param commitmentIndex The index of the unopened commitment.
     * @return An UnopenedCommitment structure representing the unopened commitment.
     */
    function getUnopenedCommitment(
        bytes32 commitmentIndex
    ) external view returns (UnopenedCommitment memory);

    /**
     * @dev Computes the bid hash for a given set of parameters.
     * @param _txnHash The transaction hash.
     * @param _revertingTxHashes The reverting transaction hashes.
     * @param _bid The bid amount.
     * @param _blockNumber The block number.
     * @param _decayStartTimeStamp The start time of the decay.
     * @param _decayEndTimeStamp The end time of the decay.
     * @return The computed bid hash.
     */
    function getBidHash(
        string memory _txnHash,
        string memory _revertingTxHashes,
        uint256 _bid,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp
    ) external pure returns (bytes32);

    /**
     * @dev Computes the pre-confirmation hash for a given set of parameters.
     * @param _txnHash The transaction hash.
     * @param _revertingTxHashes The reverting transaction hashes.
     * @param _bid The bid amount.
     * @param _blockNumber The block number.
     * @param _decayStartTimeStamp The start time of the decay.
     * @param _decayEndTimeStamp The end time of the decay.
     * @param _bidHash The bid hash.
     * @param _bidSignature The bid signature.
     * @param _sharedSecretKey The shared secret key.
     * @return The computed pre-confirmation hash.
     */
    function getPreConfHash(
        string memory _txnHash,
        string memory _revertingTxHashes,
        uint256 _bid,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp,
        bytes32 _bidHash,
        string memory _bidSignature,
        string memory _sharedSecretKey
    ) external pure returns (bytes32);

    /**
     * @dev Verifies a bid by computing the hash and recovering the signer's address.
     * @param bid The bid amount.
     * @param blockNumber The block number.
     * @param decayStartTimeStamp The start time of the decay.
     * @param decayEndTimeStamp The end time of the decay.
     * @param txnHash The transaction hash.
     * @param revertingTxHashes The reverting transaction hashes.
     * @param bidSignature The bid signature.
     * @return messageDigest The computed bid hash.
     * @return recoveredAddress The address recovered from the bid signature.
     */
    function verifyBid(
        uint256 bid,
        uint64 blockNumber,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes calldata bidSignature
    ) external pure returns (bytes32 messageDigest, address recoveredAddress);

    /**
     * @dev Verifies a pre-confirmation commitment by computing the hash and recovering the committer's address.
     * @param params The commitment params associated with the commitment.
     * @return preConfHash The hash of the pre-confirmation commitment.
     * @return committerAddress The address of the committer recovered from the commitment signature.
     */
    function verifyPreConfCommitment(
        CommitmentParams memory params
    ) external pure returns (bytes32 preConfHash, address committerAddress);

    /**
     * @dev Computes the index of an opened commitment.
     * @param commitment The commitment to compute the index for.
     * @return The computed index of the commitment.
     */
    function getOpenedCommitmentIndex(
        OpenedCommitment memory commitment
    ) external pure returns (bytes32);

    /**
     * @dev Computes the index of an unopened commitment.
     * @param commitment The unopened commitment to compute the index for.
     * @return The computed index of the unopened commitment.
     */
    function getUnopenedCommitmentIndex(
        UnopenedCommitment memory commitment
    ) external pure returns (bytes32);
}
