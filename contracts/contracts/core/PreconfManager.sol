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
            "OpenedCommitment(string txnHash,string revertingTxHashes,uint256 bidAmt,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,bytes32 bidHash,string signature,string sharedSecretKey)"
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
    bytes32 public constant ZK_CONTEXT_HASH =
        keccak256("mev-commit opening, mainnet, v1.0");
    uint256 constant GX = 1;
    uint256 constant GY = 2;
    uint256 constant BN254_MASK_253 = (1 << 253) - 1;

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
     * @param sharedSecretKey The shared secret key
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
        bytes memory sharedSecretKey,
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

        bytes32 commitmentDigest = getPreConfHash(
            txnHash,
            revertingTxHashes,
            bidAmt,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            bHash,
            bidSignature,
            sharedSecretKey
        );

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
            _providerVerifyZKProof(zkProof);
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
            // bHash,
            commitmentDigest,
            // bidSignature,
            unopenedCommitment.commitmentSignature,
            // sharedSecretKey,
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

        // emit OpenedCommitmentStored(
        //     commitmentIndex,
        //     bidderAddress,
        //     committerAddress,
        //     updatedBidAmt,
        //     blockNumber,
        //     bHash,
        //     decayStartTimeStamp,
        //     decayEndTimeStamp,
        //     txnHash,
        //     revertingTxHashes,
        //     commitmentDigest,
        //     bidSignature,
        //     unopenedCommitment.commitmentSignature,
        //     unopenedCommitment.dispatchTimestamp,
        //     sharedSecretKey
        // );
        emitOpenedCommitmentStored(commitmentIndex, newCommitment);
        // emit OpenedCommitmentStored(
        //     commitmentIndex,
        //     newCommitment.bidder,
        //     newCommitment.committer,
        //     newCommitment.blockNumber,
        //     newCommitment.bidAmt,
        //     newCommitment.commitmentDigest
        // );
        return commitmentIndex;
    }

    function emitOpenedCommitmentStored(
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

    // /**
    //  * @dev Get a commitments' enclosed transaction by its commitmentIndex.
    //  * @param commitmentIndex The index of the commitment.
    //  * @return txnHash The transaction hash.
    //  */
    // function getTxnHashFromCommitment(
    //     bytes32 commitmentIndex
    // ) public view returns (string memory txnHash) {
    //     return openedCommitments[commitmentIndex].txnHash;
    // }

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
     * @return digest it returns a digest that can be used for signing bids
     */
    function getBidHash(
        string memory _txnHash,
        string memory _revertingTxHashes,
        uint256 _bidAmt,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp,
        uint256[] calldata zkProof
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
                        zkProof[2], // _bidderPKx,
                        zkProof[3] // _bidderPKy
                    )
                )
            );
    }

    /**
     * @dev Gives digest to be signed for pre confirmation
     * @param _txnHash transaction Hash.
     * @param _bidAmt bid amount.
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
        uint256 _bidAmt,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp,
        bytes32 _bidHash,
        bytes memory _bidSignature,
        bytes memory _sharedSecretKey
    )
        public
        view
        returns (
            // uint256 _bidderPKx,
            // uint256 _bidderPKy
            bytes32
        )
    {
        return
            ECDSA.toTypedDataHash(
                domainSeparatorPreconf,
                keccak256(
                    abi.encode(
                        EIP712_COMMITMENT_TYPEHASH,
                        keccak256(bytes(_txnHash)),
                        keccak256(bytes(_revertingTxHashes)),
                        _bidAmt,
                        _blockNumber,
                        _decayStartTimeStamp,
                        _decayEndTimeStamp,
                        _bidHash,
                        keccak256(_bidSignature),
                        keccak256(_sharedSecretKey)
                        // _bidderPKx,
                        // _bidderPKy
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
            // bidderPKx,
            // bidderPKy
        );
        recoveredAddress = messageDigest.recover(bidSignature);
    }

    /**
     * @dev Verifies a pre-confirmation commitment by computing the hash and recovering the committer's address.
     * @param params The commitment params associated with the commitment.
     * @return preConfHash The hash of the pre-confirmation commitment.
     * @return committerAddress The address of the committer recovered from the commitment signature.
     */
    // function verifyPreConfCommitment(
    //     CommitmentParams memory params
    // ) public view returns (bytes32 preConfHash, address committerAddress) {
    //     preConfHash = _getPreConfHash(params);
    //     committerAddress = preConfHash.recover(params.commitmentSignature);
    // }

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

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // function _getPreConfHash(
    //     CommitmentParams memory params
    // ) internal view returns (bytes32) {
    //     return
    //         getPreConfHash(
    //             params.txnHash,
    //             params.revertingTxHashes,
    //             params.bidAmt,
    //             params.blockNumber,
    //             params.decayStartTimeStamp,
    //             params.decayEndTimeStamp,
    //             params.bidHash,
    //             params.bidSignature,
    //             params.sharedSecretKey
    //             // 1,
    //             // 2
    //         );
    // }

    function _providerVerifyZKProof(uint256[] calldata zkProof) internal view {
        require(_verifyZKProof(zkProof), "Provider's ZK proof invalid");
    }

    /**
     * @notice Verifies the sigma-proof that:
     *  (1) providerPub = g^sk
     *  (2) sharedSecret = bidPub^sk
     * using the non-interactive challenge c = H(...).
     *
     * The proof is {c, z}, plus we have inputs {providerPub, bidPub, sharedSecret}.
     */
    function _verifyZKProof(
        // // The provider's registered BN254 public key (A = g^a).
        // uint256 providerPubX,
        // uint256 providerPubY,
        // // The bidder's ephemeral BN254 public key (B = g^b).
        // uint256 bidPubX,
        // uint256 bidPubY,
        // // The final shared secret in the commitment (C = B^a).
        // uint256 sharedSecX,
        // uint256 sharedSecY,
        // // The proof data
        // uint256 c,
        // uint256 z
        uint256[] calldata zkProof
    ) internal view returns (bool) {
        // 1. Recompute a = g^z * (providerPub)^c
        (uint256 gzX, uint256 gzY) = _ecMul(GX, GY, zkProof[7]);
        (uint256 AcX, uint256 AcY) = _ecMul(zkProof[0], zkProof[1], zkProof[6]);
        (uint256 aX, uint256 aY) = _ecAdd(gzX, gzY, AcX, AcY);

        // 2. Recompute a' = B^z * C^c
        (uint256 BzX, uint256 BzY) = _ecMul(zkProof[2], zkProof[3], zkProof[7]);
        (uint256 CcX, uint256 CcY) = _ecMul(zkProof[4], zkProof[5], zkProof[6]);
        (uint256 aX2, uint256 aY2) = _ecAdd(BzX, BzY, CcX, CcY);

        // 3. Recompute c' by hashing the context + all relevant points
        bytes32 computedChallenge = keccak256(
            abi.encodePacked(
                ZK_CONTEXT_HASH,
                zkProof[0],
                zkProof[1],
                zkProof[2],
                zkProof[3],
                zkProof[4],
                zkProof[5],
                aX,
                aY,
                aX2,
                aY2
            )
        );

        // Compare the numeric value of computedChallenge vs the given c
        uint256 computedChallengeInt = uint256(computedChallenge) &
            BN254_MASK_253;
        return (computedChallengeInt == zkProof[6]);
    }

    /**
     * @dev BN128 addition precompile call:
     *       (x3, y3) = (x1, y1) + (x2, y2)
     */
    function _ecAdd(
        uint256 x1,
        uint256 y1,
        uint256 x2,
        uint256 y2
    ) internal view returns (uint256 x3, uint256 y3) {
        // 0x06 = bn128Add precompile
        // Inputs are 4 * 32 bytes = x1, y1, x2, y2
        // Output is 2 * 32 bytes = (x3, y3)
        bool success;
        assembly {
            // free memory pointer
            let memPtr := mload(0x40)
            mstore(memPtr, x1)
            mstore(add(memPtr, 0x20), y1)
            mstore(add(memPtr, 0x40), x2)
            mstore(add(memPtr, 0x60), y2)
            // call precompile
            if iszero(staticcall(gas(), 0x06, memPtr, 0x80, memPtr, 0x40)) {
                revert(0, 0)
            }
            x3 := mload(memPtr)
            y3 := mload(add(memPtr, 0x20))
            success := true
        }
        require(success, "bn128Add failed");
    }

    /**
     * @dev BN128 multiplication precompile call:
     *       (x3, y3) = scalar * (x1, y1)
     */
    function _ecMul(
        uint256 x1,
        uint256 y1,
        uint256 scalar
    ) internal view returns (uint256 x2, uint256 y2) {
        bool success;
        assembly {
            let memPtr := mload(0x40)
            mstore(memPtr, x1)
            mstore(add(memPtr, 0x20), y1)
            mstore(add(memPtr, 0x40), scalar)
            // call precompile at 0x07
            if iszero(staticcall(gas(), 0x07, memPtr, 0x60, memPtr, 0x40)) {
                revert(0, 0)
            }
            x2 := mload(memPtr)
            y2 := mload(add(memPtr, 0x20))
            success := true
        }
        require(success, "bn128Mul failed");
    }
}
