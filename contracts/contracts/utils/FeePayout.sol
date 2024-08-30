// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

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

    /// @dev Initialize a new fee tracker in storage
    function init(Tracker storage self, address _recipient, uint256 _payoutPeriodBlocks) internal {
        require(_recipient != address(0), "fee recipient is zero");
        require(_payoutPeriodBlocks != 0, "pay period must be positive");
        self.recipient = _recipient;
        self.accumulatedAmount = 0;
        self.lastPayoutBlock = block.number;
        self.payoutPeriodBlocks = _payoutPeriodBlocks;
    }

    /// @dev Transfers the accumulated fees to the recipient and resets the tracker
    /// @param tracker The FeePayout.Tracker struct
    function transferToRecipient(Tracker storage tracker) internal {
        uint256 amountToPay = tracker.accumulatedAmount;
        (bool success, ) = payable(tracker.recipient).call{value: amountToPay}("");
        require(success, "transfer to recipient failed");
        tracker.accumulatedAmount = 0;
        tracker.lastPayoutBlock = block.number;
        emit FeeTransfer(amountToPay, tracker.recipient);
    }

    /// @dev Checks if a fee payout is due
    /// @param tracker The FeePayout.Tracker struct
    /// @return true if a payout is due, false otherwise
    function isPayoutDue(Tracker storage tracker) internal view returns (bool) {
        return block.number > tracker.lastPayoutBlock + tracker.payoutPeriodBlocks;
    }
}
