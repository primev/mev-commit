// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

library FeePayout {

    struct Tracker {
        /// @dev Address that accumulates fees
        address recipient;
        /// @dev Accumulated fees since last payout
        uint256 accumulatedAmount;
        /// @dev Block number when the last fee payout was made
        uint256 lastPayoutBlock;
        /// @dev Min number of blocks between payouts
        uint256 payoutPeriodBlocks;
    }

    /// @dev Event emitted when fees are transferred to the recipient.
    event FeeTransfer(uint256 amount, address indexed recipient);

    /// @dev Event emitted when funds are accumulated for the treasury during commitment resolution
    event FundsAccumulatedForTreasury(
        bytes32 indexed commitmentDigest,
        uint256 amount,
        uint256 totalAccumulated,
        address recipient,
        address indexed bidder,
        address indexed provider
    );

    error FeeRecipientIsZero();
    error PayoutPeriodMustBePositive();
    error TransferToRecipientFailed();

    /// @dev Initialize a new fee tracker in storage
    function init(Tracker storage self, address _recipient, uint256 _payoutPeriodBlocks) internal {
        require(_recipient != address(0), FeeRecipientIsZero());
        require(_payoutPeriodBlocks != 0, PayoutPeriodMustBePositive());
        self.recipient = _recipient;
        self.accumulatedAmount = 0;
        self.lastPayoutBlock = block.number;
        self.payoutPeriodBlocks = _payoutPeriodBlocks;
    }

    /// @dev Transfers the accumulated fees to the recipient and resets the tracker
    /// @param tracker The FeePayout.Tracker struct
    function transferToRecipient(Tracker storage tracker) internal {
        uint256 amountToPay = tracker.accumulatedAmount;
        tracker.accumulatedAmount = 0;
        tracker.lastPayoutBlock = block.number;
        (bool success, ) = payable(tracker.recipient).call{value: amountToPay}("");
        require(success, TransferToRecipientFailed());
        emit FeeTransfer(amountToPay, tracker.recipient);
    }

    /// @dev Checks if a fee payout is due
    /// @param tracker The FeePayout.Tracker struct
    /// @return true if a payout is due, false otherwise
    function isPayoutDue(Tracker storage tracker) internal view returns (bool) {
        return block.number > tracker.lastPayoutBlock + tracker.payoutPeriodBlocks;
    }
}
