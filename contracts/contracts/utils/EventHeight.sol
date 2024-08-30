// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

library EventHeightLib {
    /// @title EventHeight
    /// @notice A struct to store the block height of an event, where the uint256 height value 
    /// is only relevant when the struct has been explicitly set.
    struct EventHeight {
        bool exists;
        uint256 blockHeight;
    }

    /// @notice Sets the block height of an event
    function set(EventHeight storage self, uint256 height) internal {
        self.exists = true;
        self.blockHeight = height;
    }

    /// @notice Deletes the event struct
    function del(EventHeight storage self) internal {
        self.exists = false;
        self.blockHeight = 0;
    }

    /// @notice Gets the existance and possible block height of an event
    function get(EventHeight storage self) internal view returns (bool, uint256) {
        return (self.exists, self.blockHeight);
    }
}
