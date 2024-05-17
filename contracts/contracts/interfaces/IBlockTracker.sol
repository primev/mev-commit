// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

/// @title IBlockTracker interface for BlockTracker contract
interface IBlockTracker {
    /// @notice Retrieves the builder's address corresponding to the given name.
    /// @param builderNameGrafiti The name of the block builder.
    /// @return The Ethereum address of the builder.
    function getBuilder(string calldata builderNameGrafiti) external view returns (address);

    /// @notice Gets the current window number.
    /// @return The current window number.
    function getCurrentWindow() external view returns (uint256);

    /// @notice Gets the window number for a given block number.
    /// @param blockNumber The block number.
    /// @return The window number.
    function getWindowFromBlockNumber(uint256 blockNumber) external view returns (uint256);

    /// @notice Retrieves the number of blocks per window.
    /// @return The number of blocks per window.
    function getBlocksPerWindow() external view returns (uint256);

    /// @notice Retrieves the winner of a specific L1 block.
    /// @param _blockNumber The block number of the L1 block.
    /// @return The address of the winner of the L1 block.
    function getBlockWinner(uint256 _blockNumber) external view returns (address);
    
    /// @notice Records a new L1 block with its winner.
    /// @param _blockNumber The block number of the new L1 block.
    /// @param _winnerGrafitti The graffiti of the winner of the new L1 block.
    function recordL1Block(uint256 _blockNumber, string calldata _winnerGrafitti) external;

    /// @notice Emitted when a new L1 block is recorded.
    /// @param blockNumber The block number of the new L1 block.
    /// @param winner The address of the winner of the new L1 block.
    event NewL1Block(uint256 indexed blockNumber, address indexed winner);

    /// @notice Emitted when entering a new window.
    /// @param window The new window number.
    event NewWindow(uint256 indexed window);

    /// @notice Emitted when the number of blocks per window is updated.
    /// @param blocksPerWindow The new number of blocks per window.
    event NewBlocksPerWindow(uint256 blocksPerWindow);
}
