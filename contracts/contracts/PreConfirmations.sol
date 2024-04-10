// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";
import {Ownable} from "@openzeppelin-contracts/contracts/access/Ownable.sol";

import {IProviderRegistry} from "./interfaces/IProviderRegistry.sol";
import {IBidderRegistry} from "./interfaces/IBidderRegistry.sol";
import {ECDSA} from "@openzeppelin-contracts/contracts/utils/cryptography/ECDSA.sol";

import "forge-std/console.sol";

/**
 * @title PreConfCommitmentStore - A contract for managing preconfirmation commitments and bids.
 * @notice This contract allows bidders to make precommitments and bids and provides a mechanism for the oracle to verify and process them.
 */
contract PreConfCommitmentStore is Ownable {
    using ECDSA for bytes32;

    /// @dev EIP-712 Type Hash for preconfirmation commitment
    bytes32 public constant EIP712_COMMITMENT_TYPEHASH =
        keccak256(
            "PreConfCommitment(string txnHash,uint64 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp,string bidHash,string signature)"
        );

    /// @dev EIP-712 Type Hash for preconfirmation bid
    bytes32 public constant EIP712_BID_TYPEHASH =
        keccak256("PreConfBid(string txnHash,uint64 bid,uint64 blockNumber,uint64 decayStartTimeStamp,uint64 decayEndTimeStamp)");

    // Represents the dispatch window in milliseconds
    uint64 public COMMITMENT_DISPATCH_WINDOW;

    /// @dev commitment counter
    uint256 public commitmentCount;

    /// @dev Address of the oracle
    address public oracle;

    /// @dev The last block that was processed by the Oracle
    uint256 public lastProcessedBlock;

    // EIP-712 Domain Separator
    bytes32 public DOMAIN_SEPARATOR_PRECONF;

    // EIP-712 Domain Separator
    bytes32 public DOMAIN_SEPARATOR_BID;

    /// @dev Address of provider registry
    IProviderRegistry public providerRegistry;

    /// @dev Address of bidderRegistry
    IBidderRegistry public bidderRegistry;

    /// @dev Mapping from provider to commitments count
    mapping(address => uint256) public commitmentsCount;

    /// @dev Mapping from address to commitmentss list
    mapping(address => bytes32[]) public providerCommitments;

    /// @dev Mapping for blocknumber to list of hash of commitments
    mapping(uint256 => bytes32[]) public blockCommitments;

    /// @dev Commitment Hash -> Commitemnt
    /// @dev Only stores valid commitments
    mapping(bytes32 => PreConfCommitment) public commitments;

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
        uint64 dispatchTimestamp;
    }

    /// @dev Event to log successful verifications
    event SignatureVerified(
        address indexed signer,
        string txnHash,
        uint64 indexed bid,
        uint64 blockNumber
    );

    /**
     * @dev fallback to revert all the calls.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Revert if eth sent to this contract
     */
    receive() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Makes sure transaction sender is oracle
     */
    modifier onlyOracle() {
        require(msg.sender == oracle, "Only the oracle can call this function");
        _;
    }

    /**
     * @dev Initializes the contract with the specified registry addresses, oracle, name, and version.
     * @param _providerRegistry The address of the provider registry.
     * @param _bidderRegistry The address of the bidder registry.
     * @param _oracle The address of the oracle.
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     */
    constructor(
        address _providerRegistry,
        address _bidderRegistry,
        address _oracle, 
        address _owner,
        uint64 _commitment_dispatch_window
    ) {
        oracle = _oracle;
        providerRegistry = IProviderRegistry(_providerRegistry);
        bidderRegistry = IBidderRegistry(_bidderRegistry);
        _transferOwnership(_owner);

        // EIP-712 domain separator
        DOMAIN_SEPARATOR_PRECONF = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version)"),
                keccak256("PreConfCommitment"),
                keccak256("1")
            )
        );
        // EIP-712 domain separator
        DOMAIN_SEPARATOR_BID = keccak256(
            abi.encode(
                keccak256("EIP712Domain(string name,string version)"),
                keccak256("PreConfBid"),
                keccak256("1")
            )
        );
        COMMITMENT_DISPATCH_WINDOW = _commitment_dispatch_window;
    }

    /**
     * @dev Updates the commitment dispatch window to a new value. This function can only be called by the contract owner.
     * @param newDispatchWindow The new dispatch window value to be set.
     */
    function updateCommitmentDispatchWindow(uint64 newDispatchWindow) external onlyOwner {
        COMMITMENT_DISPATCH_WINDOW = newDispatchWindow;
    }

    /**
     * @dev Gives digest to be signed for bids
     * @param _txnHash transaction Hash.
     * @param _bid bid id.
     * @param _blockNumber block number
     * @return digest it returns a digest that can be used for signing bids
     */
    function getBidHash(
        string memory _txnHash,
        uint64 _bid,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp
    ) public view returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                DOMAIN_SEPARATOR_BID,
                keccak256(
                    abi.encode(
                        EIP712_BID_TYPEHASH,
                        keccak256(abi.encodePacked(_txnHash)),
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
     * @param _bidHash hash of the bid.
     * @return digest it returns a digest that can be used for signing bids.
     */
    function getPreConfHash(
        string memory _txnHash,
        uint64 _bid,
        uint64 _blockNumber,
        uint64 _decayStartTimeStamp,
        uint64 _decayEndTimeStamp,
        bytes32 _bidHash,
        string memory _bidSignature
    ) public view returns (bytes32) {
        return
            ECDSA.toTypedDataHash(
                DOMAIN_SEPARATOR_PRECONF,
                keccak256(
                    abi.encode(
                        EIP712_COMMITMENT_TYPEHASH,
                        keccak256(abi.encodePacked(_txnHash)),
                        _bid,
                        _blockNumber,
                        _decayStartTimeStamp,
                        _decayEndTimeStamp,
                        keccak256(
                            abi.encodePacked(_bytes32ToHexString(_bidHash))
                        ),
                        keccak256(abi.encodePacked(_bidSignature))
                    )
                )
            );
    }


    /**
     * @dev Internal function to verify a bid
     * @param bid bid id.
     * @param blockNumber block number.
     * @param txnHash transaction Hash.
     * @param bidSignature bid signature.
     * @return messageDigest returns the bid hash for given bid id.
     * @return recoveredAddress the address from the bid hash.
     * @return stake the stake amount of the address for bid id bidder.
     */
    function verifyBid(
        uint64 bid,
        uint64 blockNumber,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        string memory txnHash,
        bytes calldata bidSignature
    )
        public
        view
        returns (bytes32 messageDigest, address recoveredAddress, uint256 stake)
    {
        messageDigest = getBidHash(txnHash, bid, blockNumber, decayStartTimeStamp, decayEndTimeStamp);
        recoveredAddress = messageDigest.recover(bidSignature);
        stake = bidderRegistry.getAllowance(recoveredAddress);
        require(stake > (10 * bid), "Invalid bid");
    }

    /**
     * @dev Verifies a pre-confirmation commitment by computing the hash and recovering the committer's address.
     * @param txnHash The transaction hash associated with the commitment.
     * @param bid The bid amount.
     * @param blockNumber The block number at the time of the bid.
     * @param bidHash The hash of the bid details.
     * @param bidSignature The signature of the bid.
     * @param commitmentSignature The signature of the commitment.
     * @return preConfHash The hash of the pre-confirmation commitment.
     * @return commiterAddress The address of the committer recovered from the commitment signature.
     */
    function verifyPreConfCommitment(
        string memory txnHash,
        uint64 bid,
        uint64 blockNumber,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        bytes32 bidHash,
        bytes memory bidSignature,
        bytes memory commitmentSignature
    )
        public
        view
        returns (bytes32 preConfHash, address commiterAddress)
    {
        preConfHash = getPreConfHash(
            txnHash,
            bid,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            bidHash,
            _bytesToHexString(bidSignature)
        );

        commiterAddress = preConfHash.recover(commitmentSignature);
    }

    function getCommitmentIndex(
        PreConfCommitment memory commitment
    )  public pure returns (bytes32){
        return keccak256(
            abi.encodePacked(
                commitment.commitmentHash,
                commitment.commitmentSignature
            )
        );
    }



    /**
     * @dev Store a commitment.
     * @param bid The bid amount.
     * @param blockNumber The block number.
     * @param txnHash The transaction hash.
     * @param bidSignature The signature of the bid.
     * @param commitmentSignature The signature of the commitment.
     * @param dispatchTimestamp The timestamp at which the commitment is dispatched
     * @return commitmentIndex The index of the stored commitment
     */
    function storeCommitment(
        uint64 bid,
        uint64 blockNumber,
        string memory txnHash,
        uint64 decayStartTimeStamp,
        uint64 decayEndTimeStamp,
        bytes calldata bidSignature,
        bytes memory commitmentSignature,
        uint64 dispatchTimestamp
    ) public returns (bytes32 commitmentIndex) {
        (bytes32 bHash, address bidderAddress, uint256 stake) = verifyBid(
            bid,
            blockNumber,
            decayStartTimeStamp,
            decayEndTimeStamp,
            txnHash,
            bidSignature
        );

        require(block.timestamp - dispatchTimestamp < COMMITMENT_DISPATCH_WINDOW, "Invalid dispatch timestamp, block.timestamp - dispatchTimestamp < COMMITMENT_DISPATCH_WINDOW");
        
        // This helps in avoiding stack too deep
        {
            bytes32 commitmentDigest = getPreConfHash(
                txnHash,
                bid,
                blockNumber,
                decayStartTimeStamp,
                decayEndTimeStamp,
                bHash,
                _bytesToHexString(bidSignature)
            );

            address commiterAddress = commitmentDigest.recover(commitmentSignature);

            require(stake > (10 * bid), "Stake too low");
            require(decayStartTimeStamp < decayEndTimeStamp, "Invalid decay time");
            
            PreConfCommitment memory newCommitment =  PreConfCommitment(
                false,
                bidderAddress,
                commiterAddress,
                bid,
                blockNumber,
                bHash,
                decayStartTimeStamp,
                decayEndTimeStamp,
                txnHash,
                commitmentDigest,
                bidSignature,
                commitmentSignature,
                dispatchTimestamp
            );

            commitmentIndex = getCommitmentIndex(newCommitment);

            // Store commitment
            commitments[commitmentIndex] = newCommitment;

            // Push pointers to other mappings
            providerCommitments[commiterAddress].push(commitmentIndex);
            blockCommitments[blockNumber].push(commitmentIndex);
            
            commitmentCount++;
            commitmentsCount[commiterAddress] += 1;

            // Check if Bid has bid-amt stored
            bidderRegistry.LockBidFunds(commitmentDigest, bid, bidderAddress);

        }

        return commitmentIndex;
    }

        /**
     * @dev Retrieves the list of commitments for a given committer.
     * @param commiter The address of the committer.
     * @return A list of PreConfCommitment structures for the specified committer.
     */
    function getCommitmentsByCommitter(address commiter)
        public
        view
        returns (bytes32[] memory)
    {
        return providerCommitments[commiter];
    }


    /** 
     * @dev Retrieves the list of commitments for a given block number.
    * @param blockNumber The block number.
    * @return A list of indexes referencing preconfimration structures for the specified block number.
    */
    function getCommitmentsByBlockNumber(uint256 blockNumber)
        public
        view
        returns (bytes32[] memory)
    {
        return blockCommitments[blockNumber];
    }

    /**
     * @dev Get a commitments' enclosed transaction by its commitmentIndex.
     * @param commitmentIndex The index of the commitment.
     * @return txnHash The transaction hash.
     */
    function getTxnHashFromCommitment(bytes32 commitmentIndex) public view returns (string memory txnHash)
    {
        return commitments[commitmentIndex].txnHash;
    }

    /**
     * @dev Get a commitment by its commitmentIndex.
     * @param commitmentIndex The index of the commitment.
     * @return A PreConfCommitment structure representing the commitment.
     */
    function getCommitment(
        bytes32 commitmentIndex
    ) public view returns (PreConfCommitment memory) {
        return commitments[commitmentIndex];
    }

    /**
     * @dev Initiate a slash for a commitment.
     * @param commitmentIndex The hash of the commitment to be slashed.
     */
    function initiateSlash(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) public onlyOracle {
        PreConfCommitment memory commitment = commitments[commitmentIndex];
        require(
            !commitments[commitmentIndex].commitmentUsed,
            "Commitment already used"
        );

        // Mark this commitment as used to prevent replays
        commitments[commitmentIndex].commitmentUsed = true;
        commitmentsCount[commitment.commiter] -= 1;

        providerRegistry.slash(
            commitment.bid,
            commitment.commiter,
            payable(commitment.bidder),
            residualBidPercentAfterDecay
        );

        bidderRegistry.unlockFunds(commitment.commitmentHash);
    }

    /**
        * @dev Initiate a return of funds for a bid that was not slashed.
        * @param commitmentDigest The hash of the bid to be unlocked.
     */
     function unlockBidFunds(bytes32 commitmentDigest) public onlyOracle {
        bidderRegistry.unlockFunds(commitmentDigest);
     }

    /**
     * @dev Initiate a reward for a commitment.
     * @param commitmentIndex The hash of the commitment to be rewarded.
     */
    function initiateReward(bytes32 commitmentIndex, uint256 residualBidPercentAfterDecay) public onlyOracle {
        PreConfCommitment memory commitment = commitments[commitmentIndex];
        require(
            !commitments[commitmentIndex].commitmentUsed,
            "Commitment already used"
        );

        // Mark this commitment as used to prevent replays
        commitments[commitmentIndex].commitmentUsed = true;
        commitmentsCount[commitment.commiter] -= 1;

        bidderRegistry.retrieveFunds(
            commitment.commitmentHash,
            payable(commitment.commiter),
            residualBidPercentAfterDecay
        );
    }

    /**
     * @dev Updates the address of the oracle.
     * @param newOracle The new oracle address.
     */
    function updateOracle(address newOracle) external onlyOwner {
        oracle = newOracle;
    }

    /**
     * @dev Updates the address of the provider registry.
     * @param newProviderRegistry The new provider registry address.
     */
    function updateProviderRegistry(
        address newProviderRegistry
    ) public onlyOwner {
        providerRegistry = IProviderRegistry(newProviderRegistry);
    }

    /**
     * @dev Updates the address of the bidder registry.
     * @param newBidderRegistry The new bidder registry address.
     */
    function updateBidderRegistry(address newBidderRegistry) external onlyOwner {
        bidderRegistry = IBidderRegistry(newBidderRegistry);
    }

    /**
     * @dev Internal Function to convert bytes32 to hex string without 0x
     * @param _bytes32 the byte array to convert to string
     * @return hex string from the byte 32 array
     */
    function _bytes32ToHexString(
        bytes32 _bytes32
    ) internal pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(64);
        for (uint8 i = 0; i < 32; i++) {
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
    ) public pure returns (string memory) {
        bytes memory HEXCHARS = "0123456789abcdef";
        bytes memory _string = new bytes(_bytes.length * 2);
        for (uint256 i = 0; i < _bytes.length; i++) {
            _string[i * 2] = HEXCHARS[uint8(_bytes[i] >> 4)];
            _string[1 + i * 2] = HEXCHARS[uint8(_bytes[i] & 0x0f)];
        }
        return string(_string);
    }
}
