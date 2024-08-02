// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {IProviderRegistry} from "./interfaces/IProviderRegistry.sol";
import {IBidderRegistry} from "./interfaces/IBidderRegistry.sol";
import {IBlockTracker} from "./interfaces/IBlockTracker.sol";
import {IPreConfCommitmentStore} from "./interfaces/IPreConfCommitmentStore.sol";
import {WindowFromBlockNumber} from "./utils/WindowFromBlockNumber.sol";

/**
 * @title PreConfCommitmentStore - A contract for managing preconfirmation commitments and bids.
 * @notice This contract allows bidders to make precommitments and bids and provides a mechanism for the oracle to verify and process them.
 */
contract PreConfCommitmentStore is
    IPreConfCommitmentStore,
    Ownable2StepUpgradeable,
    UUPSUpgradeable
{
    using ECDSA for bytes32;

    /// @dev EIP-712 Type Hash for preconfirmation commitment
    bytes32 public constant EIP712_COMMITMENT_TYPEHASH =
        keccak256(
            "OpenedCommitment(string txnHash,string revertingTxHashes,uint256 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,bytes32 bidHash,string signature,string sharedSecretKey)"
        );

    /// @dev EIP-712 Type Hash for preconfirmation bid
    bytes32 public constant EIP712_BID_TYPEHASH =
        keccak256(
            "PreConfBid(string txnHash,string revertingTxHashes,uint256 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp)"
        );

    // EIP-712 domain separator
    bytes32 public constant DOMAIN_SEPARATOR_PRECONF =
        keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version)"),
                keccak256("OpenedCommitment"),
                keccak256("1")
            )
        );

    // EIP-712 domain separator
    bytes32 public constant DOMAIN_SEPARATOR_BID =
        keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version)"),
                keccak256("PreConfBid"),
                keccak256("1")
            )
        );

    // Hex characters
    bytes public constant HEXCHARS = "0123456789abcdef";

    // Represents the dispatch window in milliseconds
    uint64 public commitmentDispatchWindow;

    /// @dev Address of the oracle contract
    address public oracleContract;

    /// @dev The number of blocks per window
    uint256 public blocksPerWindow;

    /// @dev Address of provider registry
    IProviderRegistry public providerRegistry;

    /// @dev Address of bidderRegistry
    IBidderRegistry public bidderRegistry;

    /// @dev Address of blockTracker
    IBlockTracker public blockTracker;

    /// @dev Mapping from provider to commitments count
    mapping(address => uint256) public commitmentsCount;

    /// @dev Commitment Hash -> Opened Commitemnt
    /// @dev Only stores valid commitments
    mapping(bytes32 => OpenedCommitment) public openedCommitments;

    /// @dev Unopened Commitment Hash -> Unopened Commitment
    /// @dev Only stores valid unopened commitments
    mapping(bytes32 => UnopenedCommitment) public unopenedCommitments;

    /**
     * @dev Makes sure transaction sender is oracle contract
     */
    modifier onlyOracleContract() {
        require(msg.sender == oracleContract, "sender is not oracle contract");
        _;
    }

    /**
     * @dev Initializes the contract with the specified registry addresses, oracle, name, and version.
     * @param _providerRegistry The address of the provider registry.
     * @param _bidderRegistry The address of the bidder registry.
     * @param _oracleContract The address of the oracle contract.
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _blockTracker The address of the block tracker.
     * @param _commitmentDispatchWindow The dispatch window for commitments.
     * @param _blocksPerWindow The number of blocks per window.
     */
    function initialize(
        address _providerRegistry,
        address _bidderRegistry,
        address _oracleContract,
        address _owner,
        address _blockTracker,
        uint64 _commitmentDispatchWindow,
        uint256 _blocksPerWindow
    ) external initializer {
        providerRegistry = IProviderRegistry(_providerRegistry);
        bidderRegistry = IBidderRegistry(_bidderRegistry);
        oracleContract = _oracleContract;
        __Ownable_init(_owner);
        blockTracker = IBlockTracker(_blockTracker);
        commitmentDispatchWindow = _commitmentDispatchWindow;
        blocksPerWindow = _blocksPerWindow;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Revert if eth sent to this contract
     */
    receive() external payable {
        revert("Invalid call");
    }

    /**
     * @dev fallback to revert all the calls.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Updates the commitment dispatch window to a new value. This function can only be called by the contract owner.
     * @param newDispatchWindow The new dispatch window value to be set.
     */
    function updateCommitmentDispatchWindow(
        uint64 newDispatchWindow
    ) external onlyOwner {
        commitmentDispatchWindow = newDispatchWindow;
    }

    /**
     * @dev Updates the address of the oracle contract.
     * @param newOracleContract The new oracle contract address.
     */
    function updateOracleContract(
        address newOracleContract
    ) external onlyOwner {
        oracleContract = newOracleContract;
    }

    /**
     * @dev Updates the address of the bidder registry.
     * @param newBidderRegistry The new bidder registry address.
     */
    function updateBidderRegistry(
        address newBidderRegistry
    ) external onlyOwner {
        bidderRegistry = IBidderRegistry(newBidderRegistry);
    }

    /**
     * @dev Updates the address of the provider registry.
     * @param newProviderRegistry The new provider registry address.
     */
    function updateProviderRegistry(
        address newProviderRegistry
    ) external onlyOwner {
        providerRegistry = IProviderRegistry(newProviderRegistry);
    }

    /**
     * @dev Open a commitment
     * @param unopenedCommitmentIndex The index of the unopened commitment
     * @param bid The bid amount
     * @param blockNumber The block number
     * @param txnHash The transaction hash
     * @param revertingTxHashes The reverting transaction hashes
     * @param decayStartTimeStamp The start time of the decay
     * @param decayEndTimeStamp The end time of the decay
     * @param bidSignature The signature of the bid
     * @param commitmentSignature The signature of the commitment
     * @param sharedSecretKey The shared secret key
     * @return commitmentIndex The index of the stored commitment
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
    ) public returns (bytes32 commitmentIndex) {
        require(decayStartTimeStamp < decayEndTimeStamp, "Invalid decay time");

        (bytes32 bHash, address bidderAddress) = verifyBid(
            bid,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            txnHash,
            revertingTxHashes,
            bidSignature
        );

        bytes32 commitmentDigest = getPreConfHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            bHash,
            _bytesToHexString(bidSignature),
            _bytesToHexString(sharedSecretKey)
        );

        UnopenedCommitment storage unopenedCommitment = unopenedCommitments[
            unopenedCommitmentIndex
        ];

        require(!unopenedCommitment.isUsed, "Commitment already used");
        require(
            unopenedCommitment.commitmentDigest == commitmentDigest,
            "Invalid commitment digest"
        );

        address committerAddress = commitmentDigest.recover(
            commitmentSignature
        );

        address winner = blockTracker.getBlockWinner(blockNumber);
        require(
            winner == committerAddress &&
                (msg.sender == winner || msg.sender == bidderAddress),
            "invalid sender"
        );

        OpenedCommitment memory newCommitment = OpenedCommitment(
            bidderAddress,
            false,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            unopenedCommitment.dispatchTimestamp,
            committerAddress,
            bid,
            bHash,
            commitmentDigest,
            bidSignature,
            commitmentSignature,
            sharedSecretKey,
            txnHash,
            revertingTxHashes
        );

        commitmentIndex = getCommitmentIndex(newCommitment);

        // Store the new commitment
        openedCommitments[commitmentIndex] = newCommitment;
        // Mark the unopened commitment as used
        unopenedCommitment.isUsed = true;

        bidderRegistry.openBid(
            commitmentDigest,
            bid,
            bidderAddress,
            blockNumber
        );

        ++commitmentsCount[committerAddress];

        emit OpenedCommitmentStored(
            commitmentIndex,
            bidderAddress,
            committerAddress,
            bid,
            blockNumber,
            bHash,
            decayStartTimeStamp,
            decayEndTimeStamp,
            txnHash,
            revertingTxHashes,
            commitmentDigest,
            bidSignature,
            commitmentSignature,
            unopenedCommitment.dispatchTimestamp,
            sharedSecretKey
        );
        return commitmentIndex;
    }

    /**
     * @dev Store an unopened commitment.
     * @param commitmentDigest The digest of the commitment.
     * @param commitmentSignature The signature of the commitment.
     * @param dispatchTimestamp The timestamp at which the commitment is dispatched.
     * @return commitmentIndex The index of the stored commitment
     */
    function storeUnopenedCommitment(
        bytes32 commitmentDigest,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp
    ) public returns (bytes32 commitmentIndex) {
        // Calculate the minimum valid timestamp for dispatching the commitment
        uint256 minTime = block.timestamp - commitmentDispatchWindow;
        // Check if the dispatch timestamp is within the allowed dispatch window
        require(dispatchTimestamp > minTime, "Invalid dispatch timestamp");

        address committerAddress = commitmentDigest.recover(
            commitmentSignature
        );

        require(committerAddress == msg.sender, "sender is not committer");

        // Ensure the provider's balance is greater than minStake and no pending withdrawal
        providerRegistry.isProviderValid(committerAddress);

        UnopenedCommitment memory newCommitment = UnopenedCommitment(
            false,
            committerAddress,
            dispatchTimestamp,
            commitmentDigest,
            commitmentSignature
        );

        commitmentIndex = getUnopenedCommitmentIndex(newCommitment);

        unopenedCommitments[commitmentIndex] = newCommitment;

        emit UnopenedCommitmentStored(
            commitmentIndex,
            committerAddress,
            commitmentDigest,
            commitmentSignature,
            dispatchTimestamp
        );

        return commitmentIndex;
    }

    /**
     * @dev Initiate a slash for a commitment.
     * @param commitmentIndex The hash of the commitment to be slashed.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function initiateSlash(
        bytes32 commitmentIndex,
        uint256 residualBidPercentAfterDecay
    ) public onlyOracleContract {
        OpenedCommitment storage commitment = openedCommitments[commitmentIndex];
        require(!commitment.isUsed, "Commitment already used");

        commitment.isUsed = true;
        --commitmentsCount[commitment.committer];

        uint256 windowToSettle = WindowFromBlockNumber.getWindowFromBlockNumber(
            commitment.blockNumber,
            blocksPerWindow
        );

        providerRegistry.slash(
            commitment.bid,
            commitment.committer,
            payable(commitment.bidder),
            residualBidPercentAfterDecay
        );

        bidderRegistry.unlockFunds(windowToSettle, commitment.commitmentDigest);
    }

    /**
     * @dev Initiate a reward for a commitment.
     * @param commitmentIndex The hash of the commitment to be rewarded.
     */
    function initiateReward(
        bytes32 commitmentIndex,
        uint256 residualBidPercentAfterDecay
    ) public onlyOracleContract {
        OpenedCommitment storage commitment = openedCommitments[commitmentIndex];
        require(!commitment.isUsed, "Commitment already used");

        uint256 windowToSettle = WindowFromBlockNumber.getWindowFromBlockNumber(
            commitment.blockNumber,
            blocksPerWindow
        );

        commitment.isUsed = true;
        --commitmentsCount[commitment.committer];

        bidderRegistry.retrieveFunds(
            windowToSettle,
            commitment.commitmentDigest,
            payable(commitment.committer),
            residualBidPercentAfterDecay
        );
    }

    /**
     * @dev Get a commitments' enclosed transaction by its commitmentIndex.
     * @param commitmentIndex The index of the commitment.
     * @return txnHash The transaction hash.
     */
    function getTxnHashFromCommitment(
        bytes32 commitmentIndex
    ) public view returns (string memory txnHash) {
        return openedCommitments[commitmentIndex].txnHash;
    }

    /**
     * @dev Get a commitment by its commitmentIndex.
     * @param commitmentIndex The index of the commitment.
     * @return A OpenedCommitment structure representing the commitment.
     */
    function getCommitment(
        bytes32 commitmentIndex
    ) public view returns (OpenedCommitment memory) {
        return openedCommitments[commitmentIndex];
    }

    /**
     * @dev Get a commitments' enclosed transaction by its commitmentIndex.
     * @param commitmentIndex The index of the commitment.
     * @return txnHash The transaction hash.
     */
    function getUnopenedCommitment(
        bytes32 commitmentIndex
    ) public view returns (UnopenedCommitment memory) {
        return unopenedCommitments[commitmentIndex];
    }

    /**
     * @dev Gives digest to be signed for bids
     * @param _txnHash transaction Hash.
     * @param _bid bid id.
     * @param _blockNumber block number
     * @param _revertingTxHashes reverting transaction hashes.
     * @return digest it returns a digest that can be used for signing bids
     */
    function getBidHash(
        string memory _txnHash,
        string memory _revertingTxHashes,
        uint256 _bid,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp
    ) public pure returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                DOMAIN_SEPARATOR_BID,
                keccak256(
                    abi.encode(
                        EIP712_BID_TYPEHASH,
                        keccak256(abi.encodePacked(_txnHash)),
                        keccak256(abi.encodePacked(_revertingTxHashes)),
                        _bid,
                        _blockNumber,
                        _decayStartTimeStamp,
                        _decayEndTimeStamp
                    )
                )
            );
    }

    /**
     * @dev Gives digest to be signed for pre confirmation
     * @param _txnHash transaction Hash.
     * @param _bid bid id.
     * @param _blockNumber block number.
     * @param _revertingTxHashes reverting transaction hashes.
     * @param _bidHash hash of the bid.
     * @param _bidSignature signature of the bid.
     * @param _sharedSecretKey shared secret key.
     * @return digest it returns a digest that can be used for signing bids.
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
    ) public pure returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                DOMAIN_SEPARATOR_PRECONF,
                keccak256(
                    abi.encode(
                        EIP712_COMMITMENT_TYPEHASH,
                        keccak256(abi.encodePacked(_txnHash)),
                        keccak256(abi.encodePacked(_revertingTxHashes)),
                        _bid,
                        _blockNumber,
                        _decayStartTimeStamp,
                        _decayEndTimeStamp,
                        keccak256(
                            abi.encodePacked(_bytes32ToHexString(_bidHash))
                        ),
                        keccak256(abi.encodePacked(_bidSignature)),
                        keccak256(abi.encodePacked(_sharedSecretKey))
                    )
                )
            );
    }

    /**
     * @dev Internal function to verify a bid
     * @param bid bid id.
     * @param blockNumber block number.
     * @param decayStartTimeStamp decay start time.
     * @param decayEndTimeStamp decay end time.
     * @param txnHash transaction Hash.
     * @param revertingTxHashes reverting transaction hashes.
     * @param bidSignature bid signature.
     * @return messageDigest returns the bid hash for given bid id.
     * @return recoveredAddress the address from the bid hash.
     */
    function verifyBid(
        uint256 bid,
        uint64 blockNumber,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes calldata bidSignature
    ) public pure returns (bytes32 messageDigest, address recoveredAddress) {
        messageDigest = getBidHash(
            txnHash,
            revertingTxHashes,
            bid,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp
        );
        recoveredAddress = messageDigest.recover(bidSignature);
    }

    /**
     * @dev Verifies a pre-confirmation commitment by computing the hash and recovering the committer's address.
     * @param params The commitment params associated with the commitment.
     * @return preConfHash The hash of the pre-confirmation commitment.
     * @return committerAddress The address of the committer recovered from the commitment signature.
     */
    function verifyPreConfCommitment(
        CommitmentParams memory params
    ) public pure returns (bytes32 preConfHash, address committerAddress) {
        preConfHash = _getPreConfHash(params);
        committerAddress = preConfHash.recover(params.commitmentSignature);
    }

    /**
     * @dev Get the index of a commitment.
     * @param commitment The commitment to get the index for.
     * @return The index of the commitment.
     */
    function getCommitmentIndex(
        OpenedCommitment memory commitment
    ) public pure returns (bytes32) {
        return
            keccak256(
                abi.encodePacked(
                    commitment.commitmentDigest,
                    commitment.commitmentSignature
                )
            );
    }

    /**
     * @dev Get the index of an unopened commitment.
     * @param commitment The commitment to get the index for.
     * @return The index of the commitment.
     */
    function getUnopenedCommitmentIndex(
        UnopenedCommitment memory commitment
    ) public pure returns (bytes32) {
        return
            keccak256(
                abi.encodePacked(
                    commitment.commitmentDigest,
                    commitment.commitmentSignature
                )
            );
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _getPreConfHash(
        CommitmentParams memory params
    ) internal pure returns (bytes32) {
        return
            getPreConfHash(
                params.txnHash,
                params.revertingTxHashes,
                params.bid,
                params.blockNumber,
                params.decayStartTimeStamp,
                params.decayEndTimeStamp,
                params.bidHash,
                _bytesToHexString(params.bidSignature),
                _bytesToHexString(params.sharedSecretKey)
            );
    }

    /**
     * @dev Internal Function to convert bytes32 to hex string without 0x
     * @param _bytes32 the byte array to convert to string
     * @return hex string from the byte 32 array
     */
    function _bytes32ToHexString(
        bytes32 _bytes32
    ) internal pure returns (string memory) {
        bytes memory _string = new bytes(64);
        for (uint8 i = 0; i < 32; ++i) {
            _string[i * 2] = HEXCHARS[uint8(_bytes32[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes32[i] & 0x0f)];
        }
        return string(_string);
    }

    /**
     * @dev Internal Function to convert bytes array to hex string without 0x
     * @param _bytes the byte array to convert to string
     * @return hex string from the bytes array
     */
    function _bytesToHexString(
        bytes memory _bytes
    ) internal pure returns (string memory) {
        bytes memory _string = new bytes(_bytes.length * 2);
        uint256 len = _bytes.length;
        for (uint256 i = 0; i < len; ++i) {
            _string[i * 2] = HEXCHARS[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }
}
