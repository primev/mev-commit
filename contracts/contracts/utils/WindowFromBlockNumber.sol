// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/**
 * @title WindowFromBlockNumber
 * @dev A library that calculates the window number for a given block number.
 */
library WindowFromBlockNumber {

    /// @dev The number of blocks per window.
    uint256 constant BLOCKS_PER_WINDOW = 10;

    /**
     * @dev Retrieves the window number for a given block number.
     * @param blockNumber The block number.
     * @return The window number.
     */
    function getWindowFromBlockNumber(uint256 blockNumber) internal pure returns (uint256) {
        return (blockNumber - 1) / BLOCKS_PER_WINDOW + 1;
    }
}