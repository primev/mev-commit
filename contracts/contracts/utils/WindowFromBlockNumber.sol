// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

/**
 * @title WindowFromBlockNumber
 * @dev A library that calculates the window number for a given block number.
 */
library WindowFromBlockNumber {
    /**
     * @dev Retrieves the window number for a given block number.
     * @param blockNumber The block number.
     * @param blocksPerWindow The number of blocks per window.
     * @return The window number.
     */
    function getWindowFromBlockNumber(uint256 blockNumber, uint256 blocksPerWindow) internal pure returns (uint256) {
        return (blockNumber - 1) / blocksPerWindow + 1;
    }
}