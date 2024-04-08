// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import {IBidderRegistry} from "./interfaces/IBidderRegistry.sol";

/// @title Bidder Registry
/// @author Kartik Chopra
/// @notice This contract is for bidder registry and staking.
contract BidderRegistry is IBidderRegistry, Ownable, ReentrancyGuard {
    /// @dev For improved precision
    uint256 constant PRECISION = 10 ** 25;
    uint256 constant PERCENT = 100 * PRECISION;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Minimum prepay required for registration
    uint256 public minAllowance;

    /// @dev Amount assigned to feeRecipient
    uint256 public feeRecipientAmount;

    /// @dev protocol fee, left over amount when there is no fee recipient assigned
    uint256 public protocolFeeAmount;

    /// @dev Address of the pre-confirmations contract
    address public preConfirmationsContract;

    /// @dev Fee recipient
    address public feeRecipient;

    /// @dev Mapping for if bidder is registered
    mapping(address => bool) public bidderRegistered;

    /// @dev Mapping from bidder addresses to their prepayed amount
    mapping(address => uint256) public bidderPrepaidBalances;

    /// @dev Mapping from bidder addresses to their locked amount based on bidID (commitmentDigest)
    mapping(bytes32 => BidState) public BidPayment;

    /// @dev Amount assigned to bidders
    mapping(address => uint256) public providerAmount;

    /// @dev Event emitted when a bidder is registered with their prepayed amount
    event BidderRegistered(address indexed bidder, uint256 prepaidAmount);

    /// @dev Event emitted when funds are retrieved from a bidder's prepay
    event FundsRetrieved(bytes32 indexed commitmentDigest, uint256 amount);

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Receive function registers bidders and takes their prepay
     * Should be removed from here in case the prepay function becomes more complex
     */
    receive() external payable {
        prepay();
    }

    /**
     * @dev Constructor to initialize the contract with a minimum prepay requirement.
     * @param _minAllowance The minimum prepay required for bidder registration.
     * @param _feeRecipient The address that receives fee
     * @param _feePercent The fee percentage for protocol
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     */
    constructor(
        uint256 _minAllowance,
        address _feeRecipient,
        uint16 _feePercent,
        address _owner 
    ) {
        minAllowance = _minAllowance;
        feeRecipient = _feeRecipient;
        feePercent = _feePercent;
        _transferOwnership(_owner);
    }

    /**
     * @dev Modifier to restrict a function to only be callable by the pre-confirmations contract.
     */
    modifier onlyPreConfirmationEngine() {
        require(
            msg.sender == preConfirmationsContract,
            "Only the pre-confirmations contract can call this function"
        );
        _;
    }

    /**
     * @dev Sets the pre-confirmations contract address. Can only be called by the owner.
     * @param contractAddress The address of the pre-confirmations contract.
     */
    function setPreconfirmationsContract(
        address contractAddress
    ) external onlyOwner {
        require(
            preConfirmationsContract == address(0),
            "Preconfirmations Contract is already set and cannot be changed."
        );
        preConfirmationsContract = contractAddress;
    }

    /**
     * @dev Get the amount assigned to a provider.
     */
    function getProviderAmount(address provider) external view returns (uint256) {
        return providerAmount[provider];
    }
    
    /**
     * @dev Get the amount assigned to the fee recipient (treasury).
     */
    function getFeeRecipientAmount() external onlyOwner view returns (uint256) {
        return feeRecipientAmount;
    }

    /**
     * @dev Internal function for bidder registration and staking.
     */
    function prepay() public payable {
        require(msg.value >= minAllowance, "Insufficient prepay");

        bidderPrepaidBalances[msg.sender] += msg.value;
        bidderRegistered[msg.sender] = true;

        emit BidderRegistered(msg.sender, bidderPrepaidBalances[msg.sender]);
    }

    /**
     * @dev Check the prepay of a bidder.
     * @param bidder The address of the bidder.
     * @return The prepayed amount for the bidder.
     */
    function getAllowance(address bidder) external view returns (uint256) {
        return bidderPrepaidBalances[bidder];
    }

    function LockBidFunds(bytes32 commitmentDigest, uint64 bid, address bidder) external onlyPreConfirmationEngine(){
        BidState memory bidState = BidPayment[commitmentDigest];
        if (bidState.state == State.Undefined) {
            BidPayment[commitmentDigest] = BidState({
                bidAmt: bid,
                state: State.PreConfirmed,
                bidder: bidder
            });
            bidderPrepaidBalances[bidder] -= bid;
        }
    }

    /**
     * @dev Retrieve funds from a bidder's prepay (only callable by the pre-confirmations contract).
     * @dev reenterancy not necessary but still putting here for precaution
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and prepayment
     * @param provider The address to transfer the retrieved funds to.
     */
    function retrieveFunds(
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external nonReentrant onlyPreConfirmationEngine {

        BidState memory bidState = BidPayment[commitmentDigest];
        require(bidState.state == State.PreConfirmed, "The bid was not preconfirmed");
        uint256 decayedAmt = ( bidState.bidAmt * residualBidPercentAfterDecay * PRECISION) / PERCENT;

        uint256 feeAmt = (decayedAmt * uint256(feePercent) * PRECISION) / PERCENT;
        uint256 amtMinusFeeAndDecay = decayedAmt - feeAmt;

        if (feeRecipient != address(0)) {
            feeRecipientAmount += feeAmt;
        } else {
            protocolFeeAmount += feeAmt;
        }

        providerAmount[provider] += amtMinusFeeAndDecay;

        // Ensures the bidder gets back the bid amount - decayed reward given to provider and protocol
        bidderPrepaidBalances[bidState.bidder] += bidState.bidAmt - decayedAmt;

        BidPayment[commitmentDigest].state = State.Withdrawn;
        BidPayment[commitmentDigest].bidAmt = 0;

        emit FundsRetrieved(commitmentDigest, decayedAmt);
    }

    /**
     * @dev Return funds to a bidder's prepay (only callable by the pre-confirmations contract).
     * @dev reenterancy not necessary but still putting here for precaution
     * @param bidID is the Bid ID that allows us to identify the bid, and prepayment
     */
    function unlockFunds(bytes32 bidID) external nonReentrant onlyPreConfirmationEngine() {
        BidState memory bidState = BidPayment[bidID];
        require(bidState.state == State.PreConfirmed, "The bid was not preconfirmed");
        uint256 amt = bidState.bidAmt;
        bidderPrepaidBalances[bidState.bidder] += amt;


        BidPayment[bidID].state = State.Withdrawn;
        BidPayment[bidID].bidAmt = 0;
        
        emit FundsRetrieved(bidID, amt);
    }

    /**
     * @notice Sets the new fee recipient
     * @dev onlyOwner restriction
     * @param newFeeRecipient The address to transfer the slashed funds to.
     */
    function setNewFeeRecipient(address newFeeRecipient) external onlyOwner {
        feeRecipient = newFeeRecipient;
    }

    /**
     * @notice Sets the new fee recipient
     * @dev onlyOwner restriction
     * @param newFeePercent this is the new fee percent
     */
    function setNewFeePercent(uint16 newFeePercent) external onlyOwner {
        feePercent = newFeePercent;
    }

    function withdrawFeeRecipientAmount() external nonReentrant {
        uint256 amount = feeRecipientAmount;
        feeRecipientAmount = 0;
        require(amount > 0, "fee recipient amount Amount is zero");
        (bool successFee, ) = feeRecipient.call{value: amount}("");
        require(successFee, "couldn't transfer to fee Recipient");
    }

    function withdrawProviderAmount(
        address payable provider
    ) external nonReentrant {
        uint256 amount = providerAmount[provider];
        providerAmount[provider] = 0;

        require(amount > 0, "provider Amount is zero");
        (bool success, ) = provider.call{value: amount}("");
        require(success, "couldn't transfer to provider");
    }

    function withdrawPrepaidAmount(address payable bidder) external nonReentrant {
        uint256 prepaidAmount = bidderPrepaidBalances[bidder];
        bidderPrepaidBalances[bidder] = 0;
        require(msg.sender == bidder, "only bidder can unprepay");
        require(prepaidAmount > 0, "bidder prepaid Amount is zero");

        (bool success, ) = bidder.call{value: prepaidAmount}("");
        require(success, "couldn't transfer prepay to bidder");
    }

    function withdrawProtocolFee(
        address payable bidder
    ) external onlyOwner nonReentrant {
        uint256 _protocolFeeAmount = protocolFeeAmount;
        protocolFeeAmount = 0;
        require(_protocolFeeAmount > 0, "insufficient protocol fee amount");

        (bool success, ) = bidder.call{value: _protocolFeeAmount}("");
        require(success, "couldn't transfer prepay to bidder");
    }
}
