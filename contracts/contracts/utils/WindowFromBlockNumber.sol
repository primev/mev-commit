// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

/**
 * @title WindowFromBlockNumber
 * @dev A library that calculates the window number for a given block number.
 */
library WindowFromBlockNumber {

    /// @dev The number of blocks per window.
    uint256 public constant BLOCKS_PER_WINDOW = 10;

    /**
     * @dev Retrieves the window number for a given block number.
     * @param blockNumber The block number.
     * @return The window number.
     */
    function getWindowFromBlockNumber(uint256 blockNumber) internal pure returns (uint256) {
        return (blockNumber - 1) / BLOCKS_PER_WINDOW + 1;
    }

    /**
     * @dev Retrieves the start and end block numbers for a given window.
     * @param window The window number.
     * @return startBlock The starting block number of the window.
     * @return endBlock The ending block number of the window.
     */
    function getBlockNumbersFromWindow(uint256 window) internal pure returns (uint256 startBlock, uint256 endBlock) {
        startBlock = (window - 1) * BLOCKS_PER_WINDOW + 1;
        endBlock = window * BLOCKS_PER_WINDOW;
    }
}