// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {BlockTrackerStorage} from "./BlockTrackerStorage.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";

/**
 * @title BlockTracker
 * @dev A contract that tracks Ethereum blocks and their winners.
 */
contract BlockTracker is IBlockTracker, BlockTrackerStorage, Ownable2StepUpgradeable, UUPSUpgradeable {

    /// @dev Modifier to ensure that the sender is the oracle account.
    modifier onlyOracle() {
        require(msg.sender == oracleAccount, "sender isn't oracle account");
        _;
    }

    /**
     * @dev Initializes the BlockTracker contract with the specified owner.
     * @param blocksPerWindow_ Number of blocks per window.
     * @param oracleAccount_ Address of the permissoined oracle account.
     * @param owner_ Address of the contract owner.
     */
    function initialize(uint256 blocksPerWindow_, address oracleAccount_, address owner_) external initializer {
        currentWindow = 1;
        blocksPerWindow = blocksPerWindow_;
        _setOracleAccount(oracleAccount_);
        __Ownable_init(owner_);
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
        revert("Invalid call");
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /**
     * @dev Allows the oracle account to add a new builder address mapping.
     * @param builderName The name of the block builder as it appears on extra data.
     * @param builderAddress The Ethereum address of the builder.
     */
    function addBuilderAddress(
        string calldata builderName,
        address builderAddress
    ) external onlyOracle {
        blockBuilderNameToAddress[builderName] = builderAddress;
    }

    /**
     * @dev Records a new L1 block and its winner.
     * @param _blockNumber The number of the new L1 block.
     * @param _winnerGraffiti The graffiti of the winner of the new L1 block.
     */
    function recordL1Block(
        uint256 _blockNumber,
        string calldata _winnerGraffiti
    ) external onlyOracle {
        address _winner = blockBuilderNameToAddress[_winnerGraffiti];
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

    /**
     * @dev Retrieves the current window number.
     * @return currentWindow The current window number.
     */
    function getCurrentWindow() external view returns (uint256) {
        return currentWindow;
    }

    /**
     * @dev Returns the builder's address corresponding to the given name.
     * @param builderNameGrafiti The name (or graffiti) of the block builder.
     * @return The Ethereum address of the builder.
     */
    function getBuilder(
        string calldata builderNameGrafiti
    ) external view returns (address) {
        return blockBuilderNameToAddress[builderNameGrafiti];
    }

    /**
     * @dev Returns the number of blocks per window.
     * @return The number of blocks per window.
     */
    function getBlocksPerWindow() external view returns (uint256) {
        return blocksPerWindow;
    }

    // Function to get the winner of a specific block
    function getBlockWinner(
        uint256 blockNumber
    ) external view returns (address) {
        return blockWinners[blockNumber];
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

    /**
    * @dev Internal function to record a new block winner
    * @param blockNumber The number of the block
    * @param winner The address of the block winner
    */
    function _recordBlockWinner(uint256 blockNumber, address winner) internal {
        // Check if the block number is valid (not 0)
        require(blockNumber != 0, "Invalid block number");

        blockWinners[blockNumber] = winner;
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
