// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";

import {IBidderRegistry} from "./interfaces/IBidderRegistry.sol";
import {IBlockTracker} from "./interfaces/IBlockTracker.sol";
import {WindowFromBlockNumber} from "./utils/WindowFromBlockNumber.sol";

/// @title Bidder Registry
/// @author Kartik Chopra
/// @notice This contract is for bidder registry and staking.
contract BidderRegistry is
    IBidderRegistry,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable
{
    /// @dev For improved precision
    uint256 public constant PRECISION = 10 ** 25;
    uint256 public constant PERCENT = 100 * PRECISION;

    /// @dev Amount assigned to feeRecipient
    uint256 public feeRecipientAmount;

    /// @dev protocol fee, left over amount when there is no fee recipient assigned
    uint256 public protocolFeeAmount;

    /// @dev Address of the pre-confirmations contract
    address public preConfirmationsContract;

    /// @dev Fee percent that would be taken by protocol when provider is slashed
    uint16 public feePercent;

    /// @dev Block tracker contract
    IBlockTracker public blockTrackerContract;

    /// @dev Fee recipient
    address public feeRecipient;

    /// @dev Mapping for if bidder is registered
    mapping(address => bool) public bidderRegistered;

    // Mapping from bidder addresses and window numbers to their locked funds
    mapping(address => mapping(uint256 => uint256)) public lockedFunds;

    // Mapping from bidder addresses and blocks to their used funds
    mapping(address => mapping(uint64 => uint256)) public usedFunds;

    /// Mapping from bidder addresses and window numbers to their funds per window
    mapping(address => mapping(uint256 => uint256)) public maxBidPerBlock;

    /// @dev Mapping from bidder addresses to their locked amount based on bidID (commitmentDigest)
    mapping(bytes32 => BidState) public BidPayment;

    /// @dev Amount assigned to bidders
    mapping(address => uint256) public providerAmount;

    /// @dev Amount assigned to bidders
    uint256 public blocksPerWindow;

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
     * @dev Initializes the contract with a minimum deposit requirement.
     * @param _feeRecipient The address that receives fee
     * @param _feePercent The fee percentage for protocol
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _blockTracker The address of the block tracker contract.
     * @param _blocksPerWindow The number of blocks per window.
     */
    function initialize(
        address _feeRecipient,
        uint16 _feePercent,
        address _owner,
        address _blockTracker,
        uint256 _blocksPerWindow
    ) external initializer {
        feeRecipient = _feeRecipient;
        feePercent = _feePercent;
        blockTrackerContract = IBlockTracker(_blockTracker);
        blocksPerWindow = _blocksPerWindow;
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Receive function registers bidders and takes their deposit
     * Should be removed from here in case the deposit function becomes more complex
     */
    receive() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
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
     * @dev Deposit for a specific window.
     * @param window The window for which the deposit is being made.
     */
    function depositForWindow(uint256 window) external payable {
        if (!bidderRegistered[msg.sender]) {
            bidderRegistered[msg.sender] = true;
        }

        uint256 newLockedFunds = lockedFunds[msg.sender][window] + msg.value;
        lockedFunds[msg.sender][window] = newLockedFunds;

        // Calculate the maximum bid per block for the given window
        maxBidPerBlock[msg.sender][window] = newLockedFunds / blocksPerWindow;

        emit BidderRegistered(msg.sender, newLockedFunds, window);
    }

    /**
     * @dev Deposit for multiple windows.
     * @param windows The windows for which the deposits are being made.
     */
    function depositForWindows(uint256[] calldata windows) external payable {
        if (!bidderRegistered[msg.sender]) {
            bidderRegistered[msg.sender] = true;
        }

        uint256 amountToDeposit = msg.value / windows.length;
        uint256 remainingAmount = msg.value % windows.length; // to handle rounding issues

        uint256 len = windows.length;
        for (uint16 i = 0; i < len; ++i) {
            uint256 window = windows[i];

            uint256 currentLockedFunds = lockedFunds[msg.sender][window];

            uint256 newLockedFunds = currentLockedFunds + amountToDeposit;
            if (i == windows.length - 1) {
                newLockedFunds += remainingAmount; // Add the remainder to the last window
            }

            lockedFunds[msg.sender][window] = newLockedFunds;
            maxBidPerBlock[msg.sender][window] =
                newLockedFunds /
                blocksPerWindow;

            emit BidderRegistered(msg.sender, newLockedFunds, window);
        }
    }

    /**
     * @dev Withdraw from specific windows.
     * @param windows The windows from which the deposit is being withdrawn.
     */
    function withdrawFromWindows(
        uint256[] calldata windows
    ) external nonReentrant {
        uint256 currentWindow = blockTrackerContract.getCurrentWindow();
        uint256 totalAmount;

        uint256 len = windows.length;
        for (uint256 i = 0; i < len; ++i) {
            uint256 window = windows[i];
            require(
                window < currentWindow,
                "withdraw after window settled"
            );

            uint256 amount = lockedFunds[msg.sender][window];
            if (amount == 0) {
                continue;
            }

            lockedFunds[msg.sender][window] = 0;
            maxBidPerBlock[msg.sender][window] = 0;
            emit BidderWithdrawal(msg.sender, window, amount);

            totalAmount += amount;
        }

        if (totalAmount == 0) {
            return;
        }

        (bool success, ) = msg.sender.call{value: totalAmount}("");
        require(success, "transfer to bidder failed");
    }

    /**
     * @dev Retrieve funds from a bidder's deposit (only callable by the pre-confirmations contract).
     * @dev reenterancy not necessary but still putting here for precaution
     * @param windowToSettle The window for which the funds are being retrieved.
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param provider The address to transfer the retrieved funds to.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function retrieveFunds(
        uint256 windowToSettle,
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external nonReentrant onlyPreConfirmationEngine {
        BidState storage bidState = BidPayment[commitmentDigest];
        require(
            bidState.state == State.PreConfirmed,
            "bid not preconfirmed"
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

        // Transfer funds back to the bidder wallet
        uint256 fundsToReturn = bidState.bidAmt - decayedAmt;
        if (fundsToReturn > 0) {
            (bool success, ) = payable(bidState.bidder).call{value: (fundsToReturn)}("");
            require(success, "transfer to bidder failed");
        }

        bidState.state = State.Withdrawn;
        bidState.bidAmt = 0;

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
     * @param window The window for which the funds are being retrieved.
     * @param bidID is the Bid ID that allows us to identify the bid, and deposit
     */
    function unlockFunds(
        uint256 window,
        bytes32 bidID
    ) external nonReentrant onlyPreConfirmationEngine {
        BidState storage bidState = BidPayment[bidID];
        require(
            bidState.state == State.PreConfirmed,
            "The bid was not preconfirmed"
        );
        uint256 amt = bidState.bidAmt;
        bidState.state = State.Withdrawn;
        bidState.bidAmt = 0;

        (bool success, ) = payable(bidState.bidder).call{value: amt}("");
        require(success, "couldn't transfer to bidder");

        emit FundsRetrieved(bidID, bidState.bidder, window, amt);
    }

    /**
     * @dev Open a bid and update the used funds for the block (only callable by the pre-confirmations contract).
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param bid The bid amount.
     * @param bidder The address of the bidder.
     * @param blockNumber The block number.
     */
    function OpenBid(
        bytes32 commitmentDigest,
        uint256 bid,
        address bidder,
        uint64 blockNumber
    ) external onlyPreConfirmationEngine {
        BidState storage bidState = BidPayment[commitmentDigest];
        if (bidState.state != State.Undefined) {
            return;
        }
        uint256 currentWindow = WindowFromBlockNumber.getWindowFromBlockNumber(
            blockNumber,
            blocksPerWindow
        );

        uint256 windowAmount = maxBidPerBlock[bidder][currentWindow];
        uint256 usedAmount = usedFunds[bidder][blockNumber];

        // Calculate the available amount for this block
        uint256 availableAmount = windowAmount > usedAmount
            ? windowAmount - usedAmount
            : 0;

        // Check if bid exceeds the available amount for the block
        if (availableAmount < bid) {
            // todo: burn it, until oracle will do the calculation and transfers
            bid = uint64(availableAmount);
        }

        // Update the used funds for the block and locked funds if bid is greater than 0
        if (bid > 0) {
            usedFunds[bidder][blockNumber] += bid;
            lockedFunds[bidder][currentWindow] -= bid;
        }

        bidState.state = State.PreConfirmed;
        bidState.bidder = bidder;
        bidState.bidAmt = bid;
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
        require(amount != 0, "fee amount is zero");
        (bool successFee, ) = feeRecipient.call{value: amount}("");
        require(successFee, "fee transfer failed");
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

        require(amount != 0, "provider amount is zero");
        (bool success, ) = provider.call{value: amount}("");
        require(success, "couldn't transfer to provider");
    }

    /**
     * @dev Withdraw funds to the bidder.
     * @param bidder The address of the bidder.
     * @param window The window for which the funds are being withdrawn.
     */
    function withdrawBidderAmountFromWindow(
        address payable bidder,
        uint256 window
    ) external nonReentrant {
        require(
            msg.sender == bidder,
            "only bidder can withdraw"
        );
        uint256 currentWindow = blockTrackerContract.getCurrentWindow();
        // withdraw is enabled only when closed and settled
        require(
            window < currentWindow,
            "window not settled"
        );
        uint256 amount = lockedFunds[bidder][window];
        require(amount != 0, "bidder amount is zero");

        lockedFunds[bidder][window] = 0;
        maxBidPerBlock[bidder][window] = 0;

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
        require(_protocolFeeAmount != 0, "insufficient protocol fee amount");

        (bool success, ) = bidder.call{value: _protocolFeeAmount}("");
        require(success, "deposit transfer failed");
    }

    /**
     * @dev Get the amount assigned to a provider.
     * @param provider The address of the provider.
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
     * @dev Check the deposit of a bidder.
     * @param bidder The address of the bidder.
     * @param window The window for which the deposit is being checked.
     * @return The deposited amount for the bidder.
     */
    function getDeposit(
        address bidder,
        uint256 window
    ) external view returns (uint256) {
        return lockedFunds[bidder][window];
    }

    function _authorizeUpgrade(address) internal override onlyOwner {} // solhint-disable no-empty-blocks
}
