// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {BidderRegistryStorage} from "./BidderRegistryStorage.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {FeePayout} from "../utils/FeePayout.sol";
import {TimestampOccurrence} from "../utils/Occurrence.sol";
import {DepositManager} from "./DepositManager.sol";


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

    modifier depositManagerIsSet() {
        require(depositManagerImpl != address(0), DepositManagerNotSet());
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
     * @param _bidderWithdrawalPeriodMs bidder withdrawal period in milliseconds (mev-commit chain uses ms timestamps)
     */
    function initialize(
        address _protocolFeeRecipient,
        uint256 _feePercent,
        address _owner,
        address _blockTracker,
        uint256 _feePayoutPeriod,
        uint256 _bidderWithdrawalPeriodMs
    ) external initializer {
        FeePayout.initTimestampTracker(protocolFeeTracker, _protocolFeeRecipient, _feePayoutPeriod);
        feePercent = _feePercent;
        blockTrackerContract = IBlockTracker(_blockTracker);
        bidderWithdrawalPeriodMs = _bidderWithdrawalPeriodMs;
        __ReentrancyGuard_init();
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
        __Pausable_init();
    }

    /**
     * @dev Enables a bidder to deposit for a specific provider.
     * @param provider The provider for which the deposit is being made.
     */
    function depositAsBidder(address provider) external payable whenNotPaused {
        require(msg.value != 0, DepositAmountIsZero());
        require(provider != address(0), ProviderIsZeroAddress());
        _depositAsBidder(provider, msg.value);
    }

    /**
     * @dev Enables a bidder to deposit eth evenly to multiple providers.
     * @param providers The providers for which the deposits are being made.
     */
    function depositEvenlyAsBidder(address[] calldata providers) external payable whenNotPaused {
        uint256 len = providers.length;
        require(len > 0, NoProviders());
        require(msg.value >= len, DepositAmountIsLessThanProviders(msg.value, len));

        uint256 amountToDeposit = msg.value / len;
        uint256 remainingAmount = msg.value % len; // to handle rounding issues

        for (uint16 i = 0; i < len; ++i) {
            address provider = providers[i];
            require(provider != address(0), ProviderIsZeroAddress());
            uint256 amount = amountToDeposit;
            if (i == len - 1) {
                amount += remainingAmount; // Add the remainder to the last provider
            }
            _depositAsBidder(provider, amount);
        }
    }

    /**
     * @dev Enables a bidder to request a withdrawal from specific providers.
     * @param providers Providers to request a withdrawal from.
     */
    function requestWithdrawalsAsBidder(address[] calldata providers) external nonReentrant whenNotPaused {
        address bidder = msg.sender;
        uint256 len = providers.length;
        require(len > 0, NoProviders());

        for (uint256 i = 0; i < len; ++i) {
            address provider = providers[i];
            Deposit storage deposit = deposits[bidder][provider];
            require(deposit.exists, DepositDoesNotExist(bidder, provider));
            require(!deposit.withdrawalRequestOccurrence.exists, WithdrawalRequestAlreadyExists(bidder, provider));
            TimestampOccurrence.captureOccurrence(deposit.withdrawalRequestOccurrence);
            emit WithdrawalRequested(bidder, provider, deposit.availableAmount, deposit.escrowedAmount,
                deposit.withdrawalRequestOccurrence.timestamp);
        }
    }

    /**
     * @dev Enables a bidder to withdraw from specific providers.
     * @param providers Providers to withdraw from.
     */
    function withdrawAsBidder(address[] calldata providers) external nonReentrant whenNotPaused {
        address bidder = msg.sender;
        uint256 totalAmount;

        uint256 len = providers.length;
        require(len > 0, NoProviders());

        for (uint256 i = 0; i < len; ++i) {
            address provider = providers[i];
            Deposit storage deposit = deposits[bidder][provider];
            require(deposit.exists, DepositDoesNotExist(bidder, provider));
            require(deposit.withdrawalRequestOccurrence.exists, WithdrawalRequestDoesNotExist(bidder, provider));
            require(deposit.withdrawalRequestOccurrence.timestamp + bidderWithdrawalPeriodMs < block.timestamp,
                WithdrawalPeriodNotElapsed(block.timestamp, deposit.withdrawalRequestOccurrence.timestamp, bidderWithdrawalPeriodMs));

            // Note deposit.escrowedAmount isn't withdrawn here as it still needs to be settled.
            // In the normal flow of the protocol it'd be zero anyways.

            uint256 availableAmount = deposit.availableAmount;
            deposit.availableAmount = 0;
            totalAmount += availableAmount;
            TimestampOccurrence.del(deposit.withdrawalRequestOccurrence);
            emit BidderWithdrawal(msg.sender, provider, availableAmount, deposit.escrowedAmount);
        }

        (bool success, ) = msg.sender.call{value: totalAmount}("");
        require(success, BidderWithdrawalTransferFailed(msg.sender, totalAmount));
    }

    /**
     * @dev Converts bidder's escrowed funds into withdrawable eth (reward) for the provider.
     * @dev This function is only callable from the pre-confirmations contract during reward settlement.
     * @dev reenterancy not necessary but still putting here for precaution
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param provider The address to transfer the retrieved funds to.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function convertFundsToProviderReward(
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external nonReentrant onlyPreconfManager whenNotPaused {
        BidState storage bidState = bidPayment[commitmentDigest];
        require(bidState.state == State.PreConfirmed, BidNotPreConfirmed(commitmentDigest, bidState.state, State.PreConfirmed));

        address bidder = bidState.bidder;
        uint256 bidAmt = bidState.bidAmt;
        bidState.state = State.Settled;
        bidState.bidAmt = 0;

        Deposit storage deposit = deposits[bidder][provider];
        deposit.escrowedAmount -= bidAmt;

        uint256 decayedAmt = (bidAmt *
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
        uint256 fundsToReturn = bidAmt - decayedAmt;
        if (fundsToReturn > 0) {
            if (!payable(bidder).send(fundsToReturn)) {
                // edge case, when bidder is rejecting transfer
                emit TransferToBidderFailed(bidder, fundsToReturn);
                deposit.availableAmount += fundsToReturn;
            }
        }

        emit FundsRewarded(
            commitmentDigest,
            bidder,
            provider,
            decayedAmt
        );
    }

    /**
     * @dev Returns escrowed funds to the bidder, since the provider is being slashed and didn't earn it.
     * @dev This function is only callable from the pre-confirmations contract during slashing.
     * @dev reenterancy not necessary but still putting here for precaution
     * @param provider that committed
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     */
    function unlockFunds(
        address provider,
        bytes32 commitmentDigest
    ) external nonReentrant onlyPreconfManager whenNotPaused {
        BidState storage bidState = bidPayment[commitmentDigest];
        require(bidState.state == State.PreConfirmed, BidNotPreConfirmed(commitmentDigest, bidState.state, State.PreConfirmed));
        
        address bidder = bidState.bidder;
        uint256 bidAmt = bidState.bidAmt;
        bidState.state = State.Settled;
        bidState.bidAmt = 0;

        Deposit storage deposit = deposits[bidder][provider];
        deposit.escrowedAmount -= bidAmt;

        if (!payable(bidder).send(bidAmt)) {
            emit TransferToBidderFailed(bidder, bidAmt);
            deposit.availableAmount += bidAmt;
        }

        emit FundsUnlocked(commitmentDigest, bidder, provider, bidAmt);
    }

    /**
     * @dev Opens a bid and escrows funds equivalent to the bid amount.
     * @param commitmentDigest is the Bid ID that allows us to identify the bid, and deposit
     * @param bidAmt The bid amount.
     * @param bidder The address of the bidder.
     * @param provider The address of the provider.
     */
    function openBid(
        bytes32 commitmentDigest,
        uint256 bidAmt,
        address bidder,
        address provider
    ) external onlyPreconfManager whenNotPaused nonReentrant depositManagerIsSet returns (uint256) {
        BidState storage bidState = bidPayment[commitmentDigest];
        if (bidState.state != State.Undefined) {
            return bidAmt;
        }

        Deposit storage deposit = deposits[bidder][provider];

        // Check if bid exceeds the available amount w.r.t bidder->provider deposit
        if (deposit.availableAmount < bidAmt) {
            // This operation shouldn't happen in normal flow. See provider node's CheckAndDeductDeposit function
            // which checks if a bidder's deposit covers the bid amount.
            bidAmt = deposit.availableAmount;
            emit BidAmountReduced(bidder, provider, bidAmt);
        }

        if (bidAmt > 0) {
            deposit.escrowedAmount += bidAmt;
            deposit.availableAmount -= bidAmt;
            if (bidder.code.length == 23 && bidder.codehash == depositManagerHash) {
                try DepositManager(payable(bidder)).topUpDeposit(provider) {
                } catch {
                    // Revert shouldn't happen, but is gracefully caught for safety
                    emit TopUpFailed(bidder, provider); 
                }
            }
        }

        bidState.state = State.PreConfirmed;
        bidState.bidder = bidder;
        bidState.bidAmt = bidAmt;

        return bidAmt;
    }

    /**
     * @dev Sets the deposit manager implementation address. Can only be called by the owner.
     * @param _depositManagerImpl The address of the deposit manager implementation.
     */
    function setDepositManagerImpl(address _depositManagerImpl) external onlyOwner {
        depositManagerImpl = _depositManagerImpl;
        depositManagerHash = keccak256(abi.encodePacked(hex"ef0100", depositManagerImpl));
        emit DepositManagerImplUpdated(depositManagerImpl);
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
     * @param provider The address of the provider.
     * @return The available deposited amount for the bidder.
     */
    function getDeposit(
        address bidder,
        address provider
    ) external view returns (uint256) {
        return deposits[bidder][provider].availableAmount;
    }

    function getDepositConsideringWithdrawalRequest(
        address bidder,
        address provider
    ) external view returns (uint256) {
        Deposit storage deposit = deposits[bidder][provider];
        if (!deposit.exists || deposit.withdrawalRequestOccurrence.exists) {
            return 0;
        }
        return deposit.availableAmount;
    }

    function getEscrowedAmount(
        address bidder,
        address provider
    ) external view returns (uint256) {
        return deposits[bidder][provider].escrowedAmount;
    }

    function withdrawalRequestExists(
        address bidder,
        address provider
    ) external view returns (bool) {
        return deposits[bidder][provider].withdrawalRequestOccurrence.exists;
    }

    /// @return protocolFee amount not yet transferred to recipient
    function getAccumulatedProtocolFee() external view returns (uint256) {
        return protocolFeeTracker.accumulatedAmount;
    }

    function _depositAsBidder(address provider, uint256 amount) internal {
        address bidder = msg.sender;
        Deposit storage deposit = deposits[bidder][provider];
        if (deposit.exists) {
            require(!deposit.withdrawalRequestOccurrence.exists, 
                WithdrawalOccurrenceExists(bidder, provider, deposit.withdrawalRequestOccurrence.timestamp));
            deposit.availableAmount += amount;
        } else {
            deposits[bidder][provider] = Deposit({
                exists: true,
                availableAmount: amount,
                escrowedAmount: 0,
                withdrawalRequestOccurrence: TimestampOccurrence.Occurrence({
                    exists: false,
                    timestamp: 0})
            });
        }
        emit BidderDeposited(bidder, provider, amount);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
