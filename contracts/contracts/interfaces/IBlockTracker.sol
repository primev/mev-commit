// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title IBlockTracker interface for BlockTracker contract
interface IBlockTracker {    

    /// @dev Event emitted when a new L1 block is tracked.
    event NewL1Block(
        uint256 indexed blockNumber,
        address indexed winner,
        uint256 indexed window
    );

    /// @notice Emitted when entering a new window.
    /// @param window The new window number.
    event NewWindow(uint256 indexed window);

    /// @dev Event emitted when the oracle account is set.
    event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount);
    
    /// @dev Event emitted when the provider registry is updated.
    event ProviderRegistrySet(address indexed oldProviderRegistry, address indexed newProviderRegistry);
    
    /// @dev Event emitted when a builder's address is added.
    event BuilderAddressAdded(string indexed builderName, address indexed builderAddress);

    error NotOracleAccount(address sender, address oracleAccount);

    error BlockNumberIsZero();
    
    /// @notice Records a new L1 block with its winner.
    /// @param _blockNumber The block number of the new L1 block.
    /// @param _winnerBLSKey The BLS public key of the winner of the new L1 block.
    function recordL1Block(uint256 _blockNumber, bytes calldata _winnerBLSKey) external;

      /// @notice Sets the provider registry.
    /// @param newProviderRegistry The address of the new provider registry.
    function setProviderRegistry(address newProviderRegistry) external;

    /// @notice Retrieves the builder's address corresponding to the given name.
    /// @param builderNameGrafiti The name of the block builder.
    /// @return The Ethereum address of the builder.
    function getBuilder(string calldata builderNameGrafiti) external view returns (address);

    /// @notice Gets the current window number.
    /// @return The current window number.
    function getCurrentWindow() external view returns (uint256);

    /// @notice Retrieves the winner of a specific L1 block.
    /// @param _blockNumber The block number of the L1 block.
    /// @return The address of the winner of the L1 block.
    function getBlockWinner(uint256 _blockNumber) external view returns (address);

    /// @notice Retrieves the number of blocks per window.
    /// @return The number of blocks per window.
    function getBlocksPerWindow() external pure returns (uint256);
}
