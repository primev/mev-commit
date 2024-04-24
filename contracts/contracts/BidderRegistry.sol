// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {Ownable} from "@openzeppelin-contracts/contracts/access/Ownable.sol";
import {ReentrancyGuard} from "@openzeppelin-contracts/contracts/security/ReentrancyGuard.sol";
import {IBidderRegistry} from "./interfaces/IBidderRegistry.sol";
import {IBlockTracker} from "./interfaces/IBlockTracker.sol";

/// @title Bidder Registry
/// @author Kartik Chopra
/// @notice This contract is for bidder registry and staking.
contract BidderRegistry is IBidderRegistry, Ownable, ReentrancyGuard {
    /// @dev For improved precision
    uint256 constant PRECISION = 10 ** 25;
    uint256 constant PERCENT = 100 * PRECISION;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Minimum deposit required for registration
    uint256 public minDeposit;

    /// @dev Amount assigned to feeRecipient
    uint256 public feeRecipientAmount;

    /// @dev protocol fee, left over amount when there is no fee recipient assigned
    uint256 public protocolFeeAmount;

    /// @dev Address of the pre-confirmations contract
    address public preConfirmationsContract;

    /// @dev Block tracker contract
    IBlockTracker public blockTrackerContract;

    /// @dev Fee recipient
    address public feeRecipient;

    /// @dev Mapping for if bidder is registered
    mapping(address => bool) public bidderRegistered;

    // Mapping from bidder addresses and window numbers to their locked funds
    mapping(address => mapping(uint256 => uint256)) public lockedFunds;

    /// @dev Mapping from bidder addresses to their locked amount based on bidID (commitmentDigest)
    mapping(bytes32 => BidState) public BidPayment;

    /// @dev Amount assigned to bidders
    mapping(address => uint256) public providerAmount;

    /// @dev Event emitted when a bidder is registered with their deposited amount
    event BidderRegistered(
        address indexed bidder,
        uint256 depositedAmount,
        uint256 windowNumber
    );

    /// @dev Event emitted when funds are retrieved from a bidder's deposit
    event FundsRetrieved(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        uint256 window,
        uint256 amount
    );

    /// @dev Event emitted when funds are retrieved from a bidder's deposit
    event FundsRewarded(
        bytes32 indexed commitmentDigest,
        address indexed bidder,
        address indexed provider,
        uint256 window,
        uint256 amount
    );

    /// @dev Event emitted when a bidder withdraws their deposit
    event BidderWithdrawal(
        address indexed bidder,
        uint256 window,
        uint256 amount
    );

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Receive function registers bidders and takes their deposit
     * Should be removed from here in case the deposit function becomes more complex
     */
    receive() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Constructor to initialize the contract with a minimum deposit requirement.
     * @param _minDeposit The minimum deposit required for bidder registration.
     * @param _feeRecipient The address that receives fee
     * @param _feePercent The fee percentage for protocol
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     */
    constructor(
        uint256 _minDeposit,
        address _feeRecipient,
        uint16 _feePercent,
        address _owner,
        address _blockTracker
    ) {
        minDeposit = _minDeposit;
        feeRecipient = _feeRecipient;
        feePercent = _feePercent;
        blockTrackerContract = IBlockTracker(_blockTracker);
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
    function getProviderAmount(
        address provider
    ) external view returns (uint256) {
        return providerAmount[provider];
    }

    /**
     * @dev Get the amount assigned to the fee recipient (treasury).
     */
    function getFeeRecipientAmount() external view onlyOwner returns (uint256) {
        return feeRecipientAmount;
    }

    /**
     * @dev Deposit for a specific window.
     * @param window The window for which the deposit is being made.
     */
    function depositForSpecificWindow(uint256 window) external payable {
        require(msg.value >= minDeposit, "Insufficient deposit");

        bidderRegistered[msg.sender] = true;
        lockedFunds[msg.sender][window] += msg.value;

        emit BidderRegistered(msg.sender, lockedFunds[msg.sender][window], window);
    }

    /**
     * @dev Check the deposit of a bidder.
     * @param bidder The address of the bidder.
     * @return The deposited amount for the bidder.
     */
    function getDeposit(
        address bidder,
        uint256 window
    ) external view returns (uint256) {
        return lockedFunds[bidder][window];
    }

    /**
     * @dev Retrieve funds from a bidder's deposit (only callable by the pre-confirmations contract).
     * @dev reenterancy not necessary but still putting here for precaution
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param provider The address to transfer the retrieved funds to.
     */
    function retrieveFunds(
        uint256 windowToSettle,
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external nonReentrant onlyPreConfirmationEngine {
        BidState memory bidState = BidPayment[commitmentDigest];
        require(
            bidState.state == State.PreConfirmed,
            "The bid was not preconfirmed"
        );
        uint256 decayedAmt = (bidState.bidAmt *
            residualBidPercentAfterDecay *
            PRECISION) / PERCENT;

        uint256 feeAmt = (decayedAmt * uint256(feePercent) * PRECISION) /
            PERCENT;
        uint256 amtMinusFeeAndDecay = decayedAmt - feeAmt;

        if (feeRecipient != address(0)) {
            feeRecipientAmount += feeAmt;
        } else {
            protocolFeeAmount += feeAmt;
        }

        providerAmount[provider] += amtMinusFeeAndDecay;

        // Ensures the bidder gets back the bid amount - decayed reward given to provider and protocol
        lockedFunds[bidState.bidder][windowToSettle] +=
            bidState.bidAmt -
            decayedAmt;

        BidPayment[commitmentDigest].state = State.Withdrawn;
        BidPayment[commitmentDigest].bidAmt = 0;

        emit FundsRewarded(
            commitmentDigest,
            bidState.bidder,
            provider,
            windowToSettle,
            decayedAmt
        );
    }

    /**
     * @dev Return funds to a bidder's deposit (only callable by the pre-confirmations contract).
     * @dev reenterancy not necessary but still putting here for precaution
     * @param bidID is the Bid ID that allows us to identify the bid, and deposit
     */
    function unlockFunds(uint256 window, bytes32 bidID) external nonReentrant onlyPreConfirmationEngine() {
        BidState memory bidState = BidPayment[bidID];
        require(bidState.state == State.PreConfirmed, "The bid was not preconfirmed");
        uint256 amt = bidState.bidAmt;
        lockedFunds[bidState.bidder][window] += amt;

        BidPayment[bidID].state = State.Withdrawn;
        BidPayment[bidID].bidAmt = 0;

        emit FundsRetrieved(bidID, bidState.bidder, window, amt);
    }

    /**
     * @dev Open a bid (only callable by the pre-confirmations contract).
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param bid The bid amount.
     * @param bidder The address of the bidder.
     */
    function OpenBid(
        bytes32 commitmentDigest,
        uint64 bid,
        address bidder,
        uint64 blockNumber
    ) external onlyPreConfirmationEngine {
        BidState memory bidState = BidPayment[commitmentDigest];
        if (bidState.state == State.Undefined) {
            uint256 currentWindow = blockTrackerContract.getWindowFromBlockNumber(blockNumber);
            // @todo delete this, when oracle will do the calculation
            // bidder cannot bid more than allowed for the round
            uint256 numberOfRounds = blockTrackerContract.getBlocksPerWindow();
            uint256 windowAmount = lockedFunds[bidder][currentWindow] / numberOfRounds;
            if (windowAmount < bid) {
                bid = uint64(windowAmount);
            }
            BidPayment[commitmentDigest] = BidState({
                state: State.PreConfirmed,
                bidder: bidder,
                bidAmt: bid
            });
            lockedFunds[bidder][currentWindow] -= bid;
        }
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

    /**
     * @dev Withdraw funds to the fee recipient.
     */
    function withdrawFeeRecipientAmount() external nonReentrant {
        uint256 amount = feeRecipientAmount;
        feeRecipientAmount = 0;
        require(amount > 0, "fee recipient amount Amount is zero");
        (bool successFee, ) = feeRecipient.call{value: amount}("");
        require(successFee, "couldn't transfer to fee Recipient");
    }

    /**
     * @dev Withdraw funds to the provider.
     * @param provider The address of the provider.
     */
    function withdrawProviderAmount(
        address payable provider
    ) external nonReentrant {
        uint256 amount = providerAmount[provider];
        providerAmount[provider] = 0;

        require(amount > 0, "provider Amount is zero");
        (bool success, ) = provider.call{value: amount}("");
        require(success, "couldn't transfer to provider");
    }

    /**
     * @dev Withdraw funds to the bidder.
     * @param bidder The address of the bidder.
     */
    function withdrawBidderAmountFromWindow(
        address payable bidder,
        uint256 window
    ) external nonReentrant {
        require(
            msg.sender == bidder,
            "only bidder can withdraw funds from window"
        );
        uint256 currentWindow = blockTrackerContract.getCurrentWindow();
        // withdraw is enabled only when closed and settled
        require(
            window < currentWindow,
            "funds can only be withdrawn after the window is settled"
        );
        uint256 amount = lockedFunds[bidder][window];
        lockedFunds[bidder][window] = 0;
        require(amount > 0, "bidder Amount is zero");

        (bool success, ) = bidder.call{value: amount}("");
        require(success, "couldn't transfer to bidder");

        emit BidderWithdrawal(bidder, window, amount);
    }

    /**
     * @dev Withdraw protocol fee.
     * @param bidder The address of the bidder.
     */
    function withdrawProtocolFee(
        address payable bidder
    ) external onlyOwner nonReentrant {
        uint256 _protocolFeeAmount = protocolFeeAmount;
        protocolFeeAmount = 0;
        require(_protocolFeeAmount > 0, "insufficient protocol fee amount");

        (bool success, ) = bidder.call{value: _protocolFeeAmount}("");
        require(success, "couldn't transfer deposit to bidder");
    }
}
