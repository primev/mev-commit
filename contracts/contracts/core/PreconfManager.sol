// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IPreconfManager} from "../interfaces/IPreconfManager.sol";
import {PreconfManagerStorage} from "./PreconfManagerStorage.sol";
import {WindowFromBlockNumber} from "../utils/WindowFromBlockNumber.sol";
import {Errors} from "../utils/Errors.sol";
import {BN128} from "../utils/BN128.sol";
import {Strings} from "@openzeppelin/contracts/utils/Strings.sol";

/**
 * @title PreconfManager - A contract for managing preconfirmation commitments and bids.
 * @notice This contract allows bidders to make precommitments and bids and provides a mechanism for the oracle to verify and process them.
 */
contract PreconfManager is
    IPreconfManager,
    PreconfManagerStorage,
    Ownable2StepUpgradeable,
    UUPSUpgradeable,
    PausableUpgradeable
{
    using ECDSA for bytes32;

    /// @dev EIP-712 Type Hash for preconfirmation commitment
    bytes32 public constant EIP712_COMMITMENT_TYPEHASH =
        keccak256(
            "OpenedCommitment(bytes32 bidHash,string signature,uint256 sharedKeyX,uint256 sharedKeyY)"
        );

    /// @dev EIP-712 Type Hash for preconfirmation bid
    bytes32 public constant EIP712_BID_TYPEHASH =
        keccak256(
            "PreConfBid(string txnHash,string revertingTxHashes,uint256 bidAmt,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,uint256 bidderPKx,uint256 bidderPKy)"
        );

    // EIP-712 domain separator
    bytes32 public domainSeparatorPreconf;
    bytes32 public domainSeparatorBid;

    // Hex characters
    bytes public constant HEXCHARS = "0123456789abcdef";

    // zk proof related contstants
    bytes32 public zkContextHash;
    uint256 private constant _GX = 1;
    uint256 private constant _GY = 2;
    uint256 private constant _BN254_MASK_253 = (1 << 253) - 1;

    /**
     * @dev Makes sure transaction sender is oracle contract
     */
    modifier onlyOracleContract() {
        require(
            msg.sender == oracleContract,
            SenderIsNotOracleContract(msg.sender, oracleContract)
        );
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
     */
    function initialize(
        address _providerRegistry,
        address _bidderRegistry,
        address _oracleContract,
        address _owner,
        address _blockTracker,
        uint64 _commitmentDispatchWindow
    ) external initializer {
        providerRegistry = IProviderRegistry(_providerRegistry);
        bidderRegistry = IBidderRegistry(_bidderRegistry);
        oracleContract = _oracleContract;
        __Ownable_init(_owner);
        blockTracker = IBlockTracker(_blockTracker);
        commitmentDispatchWindow = _commitmentDispatchWindow;
        __Pausable_init();

        // Compute the domain separators
        uint256 chainId = block.chainid;
        domainSeparatorPreconf = keccak256(
            abi.encode(
                keccak256(
                    "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
                ),
                "OpenedCommitment",
                "1",
                chainId,
                address(this)
            )
        );
        domainSeparatorBid = keccak256(
            abi.encode(
                keccak256(
                    "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
                ),
                "PreConfBid",
                "1",
                chainId,
                address(this)
            )
        );

        zkContextHash = keccak256(
            abi.encodePacked("mev-commit opening ", Strings.toString(chainId))
        );
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
        revert Errors.InvalidReceive();
    }

    /**
     * @dev fallback to revert all the calls.
     */
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /**
     * @dev Updates the commitment dispatch window to a new value. This function can only be called by the contract owner.
     * @param newDispatchWindow The new dispatch window value to be set.
     */
    function updateCommitmentDispatchWindow(
        uint64 newDispatchWindow
    ) external onlyOwner {
        commitmentDispatchWindow = newDispatchWindow;
        emit CommitmentDispatchWindowUpdated(newDispatchWindow);
    }

    /**
     * @dev Updates the address of the oracle contract.
     * @param newOracleContract The new oracle contract address.
     */
    function updateOracleContract(
        address newOracleContract
    ) external onlyOwner {
        oracleContract = newOracleContract;
        emit OracleContractUpdated(newOracleContract);
    }

    /**
     * @dev Updates the address of the provider registry.
     * @param newProviderRegistry The new provider registry address.
     */
    function updateProviderRegistry(
        address newProviderRegistry
    ) external onlyOwner {
        providerRegistry = IProviderRegistry(newProviderRegistry);
        emit ProviderRegistryUpdated(newProviderRegistry);
    }

    /**
     * @dev Updates the address of the bidder registry.
     * @param newBidderRegistry The new bidder registry address.
     */
    function updateBidderRegistry(
        address newBidderRegistry
    ) external onlyOwner {
        bidderRegistry = IBidderRegistry(newBidderRegistry);
        emit BidderRegistryUpdated(newBidderRegistry);
    }

    /**
     * @dev Updates the address of the block tracker.
     * @param newBlockTracker The new block tracker address.
     */
    function updateBlockTracker(address newBlockTracker) external onlyOwner {
        blockTracker = IBlockTracker(newBlockTracker);
        emit BlockTrackerUpdated(newBlockTracker);
    }

    /// @dev Allows the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Allows the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Open a commitment
     * @param unopenedCommitmentIndex The index of the unopened commitment
     * @param bidAmt The bid amount
     * @param blockNumber The block number
     * @param txnHash The transaction hash
     * @param revertingTxHashes The reverting transaction hashes
     * @param decayStartTimeStamp The start time of the decay
     * @param decayEndTimeStamp The end time of the decay
     * @param bidSignature The signature of the bid
     * @param zkProof The zk proof
     * @return commitmentIndex The index of the stored commitment
     */
    function openCommitment(
        bytes32 unopenedCommitmentIndex,
        uint256 bidAmt,
        uint64 blockNumber,
        string memory txnHash,
        string memory revertingTxHashes,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        bytes calldata bidSignature,
        uint256[] calldata zkProof
    ) public whenNotPaused returns (bytes32 commitmentIndex) {
        if (decayStartTimeStamp >= decayEndTimeStamp) {
            revert InvalidDecayTime(decayStartTimeStamp, decayEndTimeStamp);
        }

        (bytes32 bHash, address bidderAddress) = verifyBid(
            bidAmt,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            txnHash,
            revertingTxHashes,
            bidSignature,
            zkProof
        );

        bytes32 txnHashBidderBlockNumber = keccak256(
            abi.encode(txnHash, bidderAddress, blockNumber)
        );

        require(
            processedTxnHashes[txnHashBidderBlockNumber] == false,
            TxnHashAlreadyProcessed(txnHash, bidderAddress)
        );

        bytes32 commitmentDigest = getPreConfHash(bHash, bidSignature, zkProof);

        UnopenedCommitment storage unopenedCommitment = unopenedCommitments[
            unopenedCommitmentIndex
        ];

        if (unopenedCommitment.isOpened) {
            revert CommitmentAlreadyOpened(unopenedCommitmentIndex);
        }

        if (unopenedCommitment.commitmentDigest != commitmentDigest) {
            revert InvalidCommitmentDigest(
                unopenedCommitment.commitmentDigest,
                commitmentDigest
            );
        }

        address committerAddress = commitmentDigest.recover(
            unopenedCommitment.commitmentSignature
        );

        address winner = blockTracker.getBlockWinner(blockNumber);

        if (winner != committerAddress) {
            revert WinnerIsNotCommitter(committerAddress, winner);
        }

        if (msg.sender == winner) {
            require(
                _verifyZKProof(zkProof),
                ProviderZKProofInvalid(msg.sender, commitmentDigest)
            );
        }

        if (msg.sender != winner && msg.sender != bidderAddress) {
            revert UnauthorizedOpenCommitment(
                committerAddress,
                bidderAddress,
                msg.sender
            );
        }

        OpenedCommitment memory newCommitment = OpenedCommitment(
            bidderAddress,
            false,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            unopenedCommitment.dispatchTimestamp,
            committerAddress,
            bidAmt,
            commitmentDigest,
            unopenedCommitment.commitmentSignature,
            txnHash,
            revertingTxHashes
        );

        commitmentIndex = getOpenedCommitmentIndex(newCommitment);

        uint256 updatedBidAmt = bidderRegistry.openBid(
            commitmentDigest,
            bidAmt,
            bidderAddress,
            blockNumber
        );

        newCommitment.bidAmt = updatedBidAmt;

        // Store the new commitment
        openedCommitments[commitmentIndex] = newCommitment;
        // Mark the unopened commitment as opened
        unopenedCommitment.isOpened = true;

        ++commitmentsCount[committerAddress];

        processedTxnHashes[txnHashBidderBlockNumber] = true;

        _emitOpenedCommitmentStored(commitmentIndex, newCommitment);
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
    ) public whenNotPaused returns (bytes32 commitmentIndex) {
        // Calculate the minimum valid timestamp for dispatching the commitment
        uint256 minTime = block.timestamp - commitmentDispatchWindow;
        // Check if the dispatch timestamp is within the allowed dispatch window
        if (dispatchTimestamp <= minTime) {
            revert InvalidDispatchTimestamp(minTime, dispatchTimestamp);
        }

        address committerAddress = commitmentDigest.recover(
            commitmentSignature
        );

        if (committerAddress != msg.sender) {
            revert SenderIsNotCommitter(committerAddress, msg.sender);
        }

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

        require(
            unopenedCommitments[commitmentIndex].committer == address(0),
            UnopenedCommitmentAlreadyExists(commitmentIndex)
        );

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
    ) public onlyOracleContract whenNotPaused {
        OpenedCommitment storage commitment = openedCommitments[
            commitmentIndex
        ];
        require(
            !commitment.isSettled,
            CommitmentAlreadySettled(commitmentIndex)
        );

        commitment.isSettled = true;
        --commitmentsCount[commitment.committer];

        uint256 windowToSettle = WindowFromBlockNumber.getWindowFromBlockNumber(
            commitment.blockNumber
        );

        providerRegistry.slash(
            commitment.bidAmt,
            commitment.committer,
            payable(commitment.bidder),
            residualBidPercentAfterDecay
        );

        bidderRegistry.unlockFunds(windowToSettle, commitment.commitmentDigest);
    }

    /**
     * @dev Initiate a reward for a commitment.
     * @param commitmentIndex The hash of the commitment to be rewarded.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function initiateReward(
        bytes32 commitmentIndex,
        uint256 residualBidPercentAfterDecay
    ) public onlyOracleContract whenNotPaused {
        OpenedCommitment storage commitment = openedCommitments[
            commitmentIndex
        ];
        require(
            !commitment.isSettled,
            CommitmentAlreadySettled(commitmentIndex)
        );

        uint256 windowToSettle = WindowFromBlockNumber.getWindowFromBlockNumber(
            commitment.blockNumber
        );

        commitment.isSettled = true;
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
     * @param _bidAmt bid amount.
     * @param _blockNumber block number
     * @param _revertingTxHashes reverting transaction hashes.
     * @param _decayStartTimeStamp decay start time.
     * @param _decayEndTimeStamp decay end time.
     * @param _zkProof zk proof with bidder public key.
     * @return digest it returns a digest that can be used for signing bids
     */
    function getBidHash(
        string memory _txnHash,
        string memory _revertingTxHashes,
        uint256 _bidAmt,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp,
        uint256[] calldata _zkProof
    ) public view returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                domainSeparatorBid,
                keccak256(
                    abi.encode(
                        EIP712_BID_TYPEHASH,
                        keccak256(bytes(_txnHash)),
                        keccak256(bytes(_revertingTxHashes)),
                        _bidAmt,
                        _blockNumber,
                        _decayStartTimeStamp,
                        _decayEndTimeStamp,
                        _zkProof[2], // _bidderPKx,
                        _zkProof[3] // _bidderPKy
                    )
                )
            );
    }

    /**
     * @dev Gives digest to be signed for pre confirmation
     * @param _bidHash hash of the bid.
     * @param _bidSignature signature of the bid.
     * @param _zkProof zk proof.
     * @return digest it returns a digest that can be used for signing bids.
     */
    function getPreConfHash(
        bytes32 _bidHash,
        bytes memory _bidSignature,
        uint256[] calldata _zkProof
    ) public view returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                domainSeparatorPreconf,
                keccak256(
                    abi.encode(
                        EIP712_COMMITMENT_TYPEHASH,
                        _bidHash,
                        keccak256(_bidSignature),
                        _zkProof[4], // sharedKeyX
                        _zkProof[5] // sharedKeyY
                    )
                )
            );
    }

    /**
     * @dev Internal function to verify a bid
     * @param bidAmt bid amount.
     * @param blockNumber block number.
     * @param decayStartTimeStamp decay start time.
     * @param decayEndTimeStamp decay end time.
     * @param txnHash transaction Hash.
     * @param revertingTxHashes reverting transaction hashes.
     * @param bidSignature bid signature.
     * @param zkProof zk proof.
     * @return messageDigest returns the bid hash for given bid info.
     * @return recoveredAddress the address from the bid hash.
     */
    function verifyBid(
        uint256 bidAmt,
        uint64 blockNumber,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        string memory txnHash,
        string memory revertingTxHashes,
        bytes calldata bidSignature,
        uint256[] calldata zkProof
    ) public view returns (bytes32 messageDigest, address recoveredAddress) {
        messageDigest = getBidHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            zkProof
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
        CommitmentParams calldata params
    ) public view returns (bytes32 preConfHash, address committerAddress) {
        preConfHash = _getPreConfHash(params);
        committerAddress = preConfHash.recover(params.commitmentSignature);
    }

    /**
     * @dev Get the index of an opened commitment.
     * @param commitment The commitment to get the index for.
     * @return The index of the commitment.
     */
    function getOpenedCommitmentIndex(
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

    /**
     * @dev Emit OpenedCommitmentStored event
     * @param commitmentIndex The index of the stored commitment
     * @param newCommitment The commitment to be stored
     */
    function _emitOpenedCommitmentStored(
        bytes32 commitmentIndex,
        OpenedCommitment memory newCommitment
    ) internal {
        emit OpenedCommitmentStored(
            commitmentIndex,
            newCommitment.bidder,
            newCommitment.committer,
            newCommitment.bidAmt,
            newCommitment.blockNumber,
            newCommitment.decayStartTimeStamp,
            newCommitment.decayEndTimeStamp,
            newCommitment.txnHash,
            newCommitment.revertingTxHashes,
            newCommitment.commitmentDigest,
            newCommitment.dispatchTimestamp
        );
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _getPreConfHash(
        CommitmentParams calldata params
    ) internal view returns (bytes32) {
        return
            getPreConfHash(params.bidHash, params.bidSignature, params.zkProof);
    }

    /**
     * @notice Verifies the sigma-proof that:
     *  (1) providerPub = g^sk
     *  (2) sharedSecret = bidPub^sk
     * using the non-interactive challenge c = H(...).
     *
     * The proof is {c, z}(zkProof[6,7]), plus we have inputs
     * {providerPub(zkProof[0,1]), bidPub(zkProof[2,3]), sharedSecret(zkProof[4,5])}.
     */
    function _verifyZKProof(
        uint256[] calldata zkProof
    ) internal view returns (bool) {
        // 1. Recompute a = g^z * (providerPub)^c
        (uint256 gzX, uint256 gzY) = BN128._ecMul(_GX, _GY, zkProof[7]); // zkProof[7] = z
        (uint256 acX, uint256 acY) = BN128._ecMul(
            zkProof[0], // zkProof[0] = providerPubX
            zkProof[1], // zkProof[1] = providerPubY
            zkProof[6] // zkProof[6] = c
        );
        (uint256 aX, uint256 aY) = BN128._ecAdd(gzX, gzY, acX, acY);

        // 2. Recompute a' = B^z * C^c
        (uint256 bzX, uint256 bzY) = BN128._ecMul(
            zkProof[2], // zkProof[2] = bidPubX
            zkProof[3], // zkProof[3] = bidPubY
            zkProof[7] // zkProof[7] = z
        );
        (uint256 ccX, uint256 ccY) = BN128._ecMul(
            zkProof[4], // zkProof[4] = sharedSecretX
            zkProof[5], // zkProof[5] = sharedSecretY
            zkProof[6] // zkProof[6] = c
        );
        (uint256 aX2, uint256 aY2) = BN128._ecAdd(bzX, bzY, ccX, ccY);

        // 3. Recompute c' by hashing the context + all relevant points
        bytes32 computedChallenge = keccak256(
            abi.encodePacked(
                zkContextHash,
                zkProof[0], // providerPubX
                zkProof[1], // providerPubY
                zkProof[2], // bidPubX
                zkProof[3], // bidPubY
                zkProof[4], // sharedSecretX
                zkProof[5], // sharedSecretY
                aX,
                aY,
                aX2,
                aY2
            )
        );

        // Compare the numeric value of computedChallenge vs the given c
        uint256 computedChallengeInt = uint256(computedChallenge) &
            _BN254_MASK_253;
        return (computedChallengeInt == zkProof[6]); // zkProof[6] = c
    }
}
