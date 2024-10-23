// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {IProviderRegistry} from "../interfaces/IProviderRegistry.sol";
import {BlockTrackerStorage} from "./BlockTrackerStorage.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {Errors} from "../utils/Errors.sol";
import {WindowFromBlockNumber} from "../utils/WindowFromBlockNumber.sol";

/**
 * @title BlockTracker
 * @dev A contract that tracks Ethereum blocks and their winners.
 */
contract BlockTracker is IBlockTracker, BlockTrackerStorage,
    Ownable2StepUpgradeable, UUPSUpgradeable, PausableUpgradeable {

    /// @dev Modifier to ensure that the sender is the oracle account.
    modifier onlyOracle() {
        require(msg.sender == oracleAccount, NotOracleAccount(msg.sender, oracleAccount));
        _;
    }

    /**
     * @dev Initializes the BlockTracker contract with the specified owner.
     * @param oracleAccount_ Address of the permissoined oracle account.
     * @param owner_ Address of the contract owner.
     */
    function initialize(address oracleAccount_, address owner_) external initializer {
        currentWindow = 1;
        _setOracleAccount(oracleAccount_);
        __Ownable_init(owner_);
        __Pausable_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev Receive function is disabled for this contract to prevent unintended interactions.
     * Should be removed from here in case the registerAndStake function becomes more complex
     */
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /**
     * @dev Allows the oracle account to add a new builder address mapping.
     * @param builderName The name of the block builder as it appears on extra data.
     * @param builderAddress The Ethereum address of the builder.
     */
    function addBuilderAddress(
        string calldata builderName,
        address builderAddress
    ) external onlyOracle whenNotPaused {
        blockBuilderNameToAddress[builderName] = builderAddress;
        emit BuilderAddressAdded(builderName, builderAddress);
    }

    /**
     * @dev Records a new L1 block and its winner.
     * @param _blockNumber The number of the new L1 block.
     * @param _winnerBLSKey The BLS key of the winner of the new L1 block.
     */
    function recordL1Block(
        uint256 _blockNumber,
        bytes calldata _winnerBLSKey
    ) external onlyOracle whenNotPaused {
        address _winner = _providerRegistry.getEoaFromBLSKey(_winnerBLSKey);
        _recordBlockWinner(_blockNumber, _winner);
        uint256 newWindow = (_blockNumber - 1) / blocksPerWindow + 1;
        if (newWindow > currentWindow) {
            // We've entered a new window
            currentWindow = newWindow;
            emit NewWindow(currentWindow);
        }
        emit NewL1Block(_blockNumber, _winner, currentWindow);
    }

    /// @dev Allows the owner to set the oracle account.
    function setOracleAccount(address newOracleAccount) external onlyOwner {
        _setOracleAccount(newOracleAccount);
    }

    /// @dev Allows the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Allows the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    /**
     * @dev Retrieves the current window number.
     * @return currentWindow The current window number.
     */
    function getCurrentWindow() external view returns (uint256) {
        return currentWindow;
    }

    /**
     * @dev Function to get the winner of a specific block.
     * @param blockNumber The number of the block.
     * @return The address of the block winner.
     */
    function getBlockWinner(uint256 blockNumber) external view returns (address) {
        return blockWinners[blockNumber];
    }

    /**
     * @dev Returns the builder's address corresponding to the given name.
     * @param builderNameGraffiti The name (or graffiti) of the block builder.
     * @return The Ethereum address of the builder.
     */
    function getBuilder(
        string calldata builderNameGraffiti
    ) external view returns (address) {
        return blockBuilderNameToAddress[builderNameGraffiti];
    }

    /**
     * @dev Returns the number of blocks per window.
     * @return The number of blocks per window.
     */
    function getBlocksPerWindow() external pure returns (uint256) {
        return WindowFromBlockNumber.BLOCKS_PER_WINDOW;
    }

    /**
     * @dev Internal function to set the oracle account.
     * @param newOracleAccount The new address of the oracle account.
     */
    function _setOracleAccount(address newOracleAccount) internal {
        address oldOracleAccount = oracleAccount;
        oracleAccount = newOracleAccount;
        emit OracleAccountSet(oldOracleAccount, newOracleAccount);
    }

    /// @dev Allows the owner to set the provider registry.
    function setProviderRegistry(address newProviderRegistry) external onlyOwner {
        _providerRegistry = IProviderRegistry(newProviderRegistry);
    }

    /**
     * @dev Internal function to record a new block winner.
     * @param blockNumber The number of the block.
     * @param winner The address of the block winner.
     */
    function _recordBlockWinner(uint256 blockNumber, address winner) internal {
        // Check if the block number is valid (not 0)
        require(blockNumber != 0, BlockNumberIsZero());

        blockWinners[blockNumber] = winner;
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
