// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {BidderRegistryStorage} from "./BidderRegistryStorage.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {WindowFromBlockNumber} from "../utils/WindowFromBlockNumber.sol";
import {FeePayout} from "../utils/FeePayout.sol";

/// @title Bidder Registry
/// @notice This contract is for bidder registry and staking.
contract BidderRegistry is
    IBidderRegistry,
    BidderRegistryStorage,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable,
    PausableUpgradeable
{
    /**
     * @dev Modifier to restrict a function to only be callable by the preconfManager contract.
     */
    modifier onlyPreconfManager() {
        require(msg.sender == preconfManager, SenderIsNotPreconfManager(msg.sender, preconfManager));
        _;
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
     * @dev Initializes the contract with a minimum deposit requirement.
     * @param _protocolFeeRecipient The address that accumulates protocol fees
     * @param _feePercent The fee percentage for protocol
     * @param _owner Owner of the contract, explicitly needed since contract is deployed w/ create2 factory.
     * @param _blockTracker The address of the block tracker contract.
     * @param _feePayoutPeriod The number of seconds or ms on the mev-commit chain for the fee payout period
     */
    function initialize(
        address _protocolFeeRecipient,
        uint256 _feePercent,
        address _owner,
        address _blockTracker,
        uint256 _feePayoutPeriod
    ) external initializer {
        FeePayout.initTimestampTracker(protocolFeeTracker, _protocolFeeRecipient, _feePayoutPeriod);
        feePercent = _feePercent;
        blockTrackerContract = IBlockTracker(_blockTracker);
        __ReentrancyGuard_init();
        __Ownable_init(_owner);
        __Pausable_init();
    }

    /**
     * @dev Deposit for a specific window.
     * @param window The window for which the deposit is being made.
     */
    function depositForWindow(uint256 window) external payable whenNotPaused {
        require(msg.value != 0, DepositAmountIsZero());

        uint256 newLockedFunds = lockedFunds[msg.sender][window] + msg.value;
        lockedFunds[msg.sender][window] = newLockedFunds;

        // Calculate the maximum bid per block for the given window
        maxBidPerBlock[msg.sender][window] = newLockedFunds / WindowFromBlockNumber.BLOCKS_PER_WINDOW;

        emit BidderRegistered(msg.sender, newLockedFunds, window);
    }

    /**
     * @dev Deposit for multiple windows.
     * @param windows The windows for which the deposits are being made.
     */
    function depositForWindows(uint256[] calldata windows) external payable whenNotPaused {
        require(msg.value != 0, DepositAmountIsZero());

        uint256 amountToDeposit = msg.value / windows.length;
        uint256 remainingAmount = msg.value % windows.length; // to handle rounding issues

        uint256 len = windows.length;
        for (uint16 i = 0; i < len; ++i) {
            uint256 window = windows[i];

            uint256 currentLockedFunds = lockedFunds[msg.sender][window];

            uint256 newLockedFunds = currentLockedFunds + amountToDeposit;
            if (i == len - 1) {
                newLockedFunds += remainingAmount; // Add the remainder to the last window
            }

            lockedFunds[msg.sender][window] = newLockedFunds;
            maxBidPerBlock[msg.sender][window] =
                newLockedFunds /
                WindowFromBlockNumber.BLOCKS_PER_WINDOW;

            emit BidderRegistered(msg.sender, newLockedFunds, window);
        }
    }

    /**
     * @dev Withdraw from specific windows.
     * @param windows The windows from which the deposit is being withdrawn.
     */
    function withdrawFromWindows(
        uint256[] calldata windows
    ) external nonReentrant whenNotPaused {
        uint256 currentWindow = blockTrackerContract.getCurrentWindow();
        uint256 totalAmount;

        uint256 len = windows.length;
        for (uint256 i = 0; i < len; ++i) {
            uint256 window = windows[i];
            require(window < currentWindow, WithdrawAfterWindowSettled(window, currentWindow));

            uint256 amount = lockedFunds[msg.sender][window];

            lockedFunds[msg.sender][window] = 0;
            maxBidPerBlock[msg.sender][window] = 0;

            (uint256 startBlock, uint256 endBlock) = WindowFromBlockNumber.getBlockNumbersFromWindow(window);

            for (uint256 blockNumber = startBlock; blockNumber <= endBlock; ++blockNumber) {
                usedFunds[msg.sender][uint64(blockNumber)] = 0;
            }

            emit BidderWithdrawal(msg.sender, window, amount);

            totalAmount += amount;
        }

        (bool success, ) = msg.sender.call{value: totalAmount}("");
        require(success, BidderWithdrawalTransferFailed(msg.sender, totalAmount));
    }

    /**
     * @dev Converts bidder's deposited funds into withdrawable eth (reward) for the provider.
     * @dev This function is only callable from the pre-confirmations contract during reward settlement.
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
    ) external nonReentrant onlyPreconfManager whenNotPaused {
        BidState storage bidState = bidPayment[commitmentDigest];
        require(bidState.state == State.PreConfirmed, BidNotPreConfirmed(commitmentDigest, bidState.state, State.PreConfirmed));
        
        uint256 decayedAmt = (bidState.bidAmt *
            residualBidPercentAfterDecay) / ONE_HUNDRED_PERCENT;

        uint256 feeAmt = (decayedAmt * feePercent) /
            ONE_HUNDRED_PERCENT;
        uint256 amtMinusFeeAndDecay = decayedAmt - feeAmt;

        protocolFeeTracker.accumulatedAmount += feeAmt;
        if (FeePayout.isPayoutDueByTimestamp(protocolFeeTracker)) {
            FeePayout.transferToRecipientByTimestamp(protocolFeeTracker);
        }

        providerAmount[provider] += amtMinusFeeAndDecay;

        // Transfer non-rewarded funds back to the bidder wallet
        uint256 fundsToReturn = bidState.bidAmt - decayedAmt;
        if (fundsToReturn > 0) {
            if (!payable(bidState.bidder).send(fundsToReturn)) {
                // edge case, when bidder is rejecting transfer
                emit TransferToBidderFailed(bidState.bidder, fundsToReturn);
                lockedFunds[bidState.bidder][windowToSettle] += fundsToReturn;
            }
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
     * @dev Returns escrowed funds to the bidder, since the provider is being slashed and didn't earn it.
     * @dev This function is only callable from the pre-confirmations contract during slashing.
     * @dev reenterancy not necessary but still putting here for precaution
     * @param window The window for which the funds are being retrieved.
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     */
    function unlockFunds(
        uint256 window,
        bytes32 commitmentDigest
    ) external nonReentrant onlyPreconfManager whenNotPaused {
        BidState storage bidState = bidPayment[commitmentDigest];
        require(bidState.state == State.PreConfirmed, BidNotPreConfirmed(commitmentDigest, bidState.state, State.PreConfirmed));
        
        uint256 amt = bidState.bidAmt;
        bidState.state = State.Withdrawn;
        bidState.bidAmt = 0;

        if (!payable(bidState.bidder).send(amt)) {
            emit TransferToBidderFailed(bidState.bidder, amt);
            lockedFunds[bidState.bidder][window] += amt;
        }

        emit FundsRetrieved(commitmentDigest, bidState.bidder, window, amt);
    }

    /**
     * @dev Open a bid and update the used funds for the block (only callable by the pre-confirmations contract).
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param bidAmt The bid amount.
     * @param bidder The address of the bidder.
     * @param blockNumber The block number.
     */
    function openBid(
        bytes32 commitmentDigest,
        uint256 bidAmt,
        address bidder,
        uint64 blockNumber
    ) external onlyPreconfManager whenNotPaused returns (uint256) {
        BidState storage bidState = bidPayment[commitmentDigest];
        if (bidState.state != State.Undefined) {
            return bidAmt;
        }
        uint256 currentWindow = WindowFromBlockNumber.getWindowFromBlockNumber(
            blockNumber
        );

        uint256 windowAmount = maxBidPerBlock[bidder][currentWindow];
        uint256 usedAmount = usedFunds[bidder][blockNumber];

        // Calculate the available amount for this block
        uint256 availableAmount = windowAmount > usedAmount
            ? windowAmount - usedAmount
            : 0;

        // Check if bid exceeds the available amount for the block
        if (availableAmount < bidAmt) {
            // This operation shouldn't happen in normal flow. See provider node's CheckAndDeductDeposit function
            // which checks if a bidder's deposit for the block covers the bid amount.
            bidAmt = availableAmount;
        }

        // Update the used funds for the block and locked funds if bid is greater than 0
        if (bidAmt > 0) {
            usedFunds[bidder][blockNumber] += bidAmt;
            lockedFunds[bidder][currentWindow] -= bidAmt;
        }

        bidState.state = State.PreConfirmed;
        bidState.bidder = bidder;
        bidState.bidAmt = bidAmt;

        return bidAmt;
    }

    /**
     * @dev Sets the preconfManager contract address. Can only be called by the owner.
     * @param contractAddress The address of the preconfManager contract.
     */
    function setPreconfManager(
        address contractAddress
    ) external onlyOwner {
        preconfManager = contractAddress;
        emit PreconfManagerUpdated(contractAddress);
    }

    /**
     * @notice Sets the new fee percent
     * @dev onlyOwner restriction
     * @param newFeePercent this is the new fee percent
     */
    function setNewFeePercent(uint256 newFeePercent) external onlyOwner {
        feePercent = newFeePercent;
        emit FeePercentUpdated(newFeePercent);
    }
    
    function setBlockTrackerContract(address newBlockTrackerContract) external onlyOwner {
        blockTrackerContract = IBlockTracker(newBlockTrackerContract);
        emit BlockTrackerUpdated(newBlockTrackerContract);
    }
    
    /**
     * @notice Sets the new fee recipient
     * @dev onlyOwner restriction
     * @param newProtocolFeeRecipient The new address to accumulate protocol fees
     */
    function setNewProtocolFeeRecipient(address newProtocolFeeRecipient) external onlyOwner {
        protocolFeeTracker.recipient = newProtocolFeeRecipient;
        emit ProtocolFeeRecipientUpdated(newProtocolFeeRecipient);
    }

    /** 
     * @notice Sets the new fee payout period in seconds or ms on the mev-commit chain
     * @dev onlyOwner restriction
     * @param newFeePayoutPeriod The new fee payout period in seconds or ms on the mev-commit chain
     */
    function setNewFeePayoutPeriod(uint256 newFeePayoutPeriod) external onlyOwner {
        protocolFeeTracker.payoutTimePeriod = newFeePayoutPeriod;
        emit FeePayoutPeriodUpdated(newFeePayoutPeriod);
    }

    /**
     * @dev Withdraw funds rewarded to the provider for fulfilling commitments.
     * @param provider The address of the provider.
     */
    function withdrawProviderAmount(
        address payable provider
    ) external nonReentrant whenNotPaused {
        uint256 amount = providerAmount[provider];
        providerAmount[provider] = 0;

        require(amount != 0, ProviderAmountIsZero(provider));

        (bool success, ) = provider.call{value: amount}("");
        require(success, TransferToProviderFailed(provider, amount));
    }

    /**
     * @dev Withdraw funds to the bidder.
     * @param bidder The address of the bidder.
     * @param window The window for which the funds are being withdrawn.
     */
    function withdrawBidderAmountFromWindow(
        address payable bidder,
        uint256 window
    ) external nonReentrant whenNotPaused {
        require(msg.sender == bidder, OnlyBidderCanWithdraw(msg.sender, bidder));
        uint256 currentWindow = blockTrackerContract.getCurrentWindow();
        // withdraw is enabled only when closed and settled
        require(window < currentWindow, WindowNotSettled());
        uint256 amount = lockedFunds[bidder][window];

        lockedFunds[bidder][window] = 0;
        maxBidPerBlock[bidder][window] = 0;

        (uint256 startBlock, uint256 endBlock) = WindowFromBlockNumber.getBlockNumbersFromWindow(window);

        for (uint256 blockNumber = startBlock; blockNumber <= endBlock; ++blockNumber) {
            usedFunds[bidder][uint64(blockNumber)] = 0;
        }

        (bool success, ) = bidder.call{value: amount}("");
        require(success, BidderWithdrawalTransferFailed(bidder, amount));

        emit BidderWithdrawal(bidder, window, amount);
    }

    /**
     * @dev Manually withdraws accumulated protocol fees to the recipient
     * to cover the edge case that oracle doesn't slash/reward, and funds still need to be withdrawn.
     */
    function manuallyWithdrawProtocolFee() external onlyOwner {
        FeePayout.transferToRecipientByTimestamp(protocolFeeTracker);
    }

    /// @dev Allows owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Allows owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Get the amount of funds rewarded to a provider for fulfilling commitments
     * @param provider The address of the provider.
     */
    function getProviderAmount(
        address provider
    ) external view returns (uint256) {
        return providerAmount[provider];
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

    /// @return protocolFee amount not yet transferred to recipient
    function getAccumulatedProtocolFee() external view returns (uint256) {
        return protocolFeeTracker.accumulatedAmount;
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
