pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin-contracts/contracts/access/Ownable.sol";

/**
 * @title BlockTracker
 * @dev A contract that tracks Ethereum blocks and their winners.
 */
contract BlockTracker is Ownable {
    /// @dev Event emitted when a new L1 block is tracked.
    event NewL1Block(
        uint256 indexed blockNumber,
        address indexed winner,
        uint256 indexed window
    );

    /// @dev Event emitted when a new window is created.
    event NewWindow(uint256 indexed window);

    /// @dev Event emitted when the number of blocks per window is updated.
    event NewBlocksPerWindow(uint256 blocksPerWindow);

    uint256 public currentWindow = 1;
    uint256 public blocksPerWindow = 64;

    // Mapping from block number to the winner's address
    mapping(uint256 => address) public blockWinners;

     /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    /**
     * @dev Initializes the BlockTracker contract with the specified owner.
     * @param _owner The address of the contract owner.
     */
    constructor(address _owner) Ownable() {
        _transferOwnership(_owner);
    }

    /**
     * @dev Retrieves the current window number.
     * @return The current window number.
     */
    function getCurrentWindow() external view returns (uint256) {
        return currentWindow;
    }

    /**
     * @dev Allows the owner to add a new builder address.
     * @param builderName The name of the block builder as it appears on extra data.
     * @param builderAddress The Ethereum address of the builder.
     */
    function addBuilderAddress(
        string memory builderName,
        address builderAddress
    ) external onlyOwner {
        blockBuilderNameToAddress[builderName] = builderAddress;
    }

    /**
     * @dev Returns the builder's address corresponding to the given name.
     * @param builderNameGrafiti The name (or graffiti) of the block builder.
     */
    function getBuilder(
        string calldata builderNameGrafiti
    ) external view returns (address) {
        return blockBuilderNameToAddress[builderNameGrafiti];
    }

    /**
     * @dev Returns the window number corresponding to a given block number.
     * @param blockNumber The block number.
     * @return The window number.
     */
    function getWindowFromBlockNumber(uint256 blockNumber) external view returns (uint256) {
        return (blockNumber - 1) / blocksPerWindow + 1;
    }

    /**
     * @dev Returns the number of blocks per window.
     * @return The number of blocks per window.
     */
    function getBlocksPerWindow() external view returns (uint256) {
        return blocksPerWindow;
    }

    /**
     * @dev Records a new L1 block and its winner.
     * @param _blockNumber The number of the new L1 block.
     * @param _winnerGraffiti The graffiti of the winner of the new L1 block.
     */
    function recordL1Block(
        uint256 _blockNumber,
        string calldata _winnerGraffiti
    ) external onlyOwner {
        address _winner = blockBuilderNameToAddress[_winnerGraffiti];
        recordBlockWinner(_blockNumber, _winner);
        uint256 newWindow = (_blockNumber - 1) / blocksPerWindow + 1;
        if (newWindow > currentWindow) {
            // We've entered a new window
            currentWindow = newWindow;
            emit NewWindow(currentWindow);
        }
        emit NewL1Block(_blockNumber, _winner, currentWindow);
    }

    // Function to record a new block winner
    function recordBlockWinner(uint256 blockNumber, address winner) internal {
        // Check if the block number is valid (not 0)
        require(blockNumber != 0, "Invalid block number");

        // Check if the winner address is valid (not the zero address)
        require(winner != address(0), "Invalid winner address");

        blockWinners[blockNumber] = winner;
    }

    // Function to get the winner of a specific block
    function getBlockWinner(
        uint256 blockNumber
    ) external view returns (address) {
        return blockWinners[blockNumber];
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Receive function is disabled for this contract to prevent unintended interactions.
     * Should be removed from here in case the registerAndStake function becomes more complex
     */
    receive() external payable {
        revert("Invalid call");
    }
}
