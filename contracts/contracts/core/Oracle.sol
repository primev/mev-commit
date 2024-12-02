// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IPreconfManager} from "../interfaces/IPreconfManager.sol";
import {IBlockTracker} from "../interfaces/IBlockTracker.sol";
import {OracleStorage} from "./OracleStorage.sol";
import {IOracle} from "../interfaces/IOracle.sol";
import {Errors} from "../utils/Errors.sol";

/**
 * @title Oracle
 * @notice A contract for Fetching L1 Block Builder Info and Block Data.
 * @author Kartik Chopra
 * @dev This contract serves as an oracle to fetch and process Ethereum Layer 1 block data.
 */
contract Oracle is OracleStorage, IOracle, Ownable2StepUpgradeable, UUPSUpgradeable, PausableUpgradeable {

    /// @dev Modifier to ensure that the sender is the oracle account.
    modifier onlyOracle() {
        require(msg.sender == oracleAccount, NotOracleAccount(msg.sender, oracleAccount));
        _;
    }

    /**
     * @dev Initializes the contract with a PreConfirmations contract.
     * @param preconfManager_ The address of the preconf manager contract.
     * @param blockTrackerContract_ The address of the block tracker contract.
     * @param oracleAccount_ The address of the oracle account.
     * @param owner_ Owner of the contract, explicitly needed since contract is deployed with create2 factory.
     */
    function initialize(
        address preconfManager_,
        address blockTrackerContract_,
        address oracleAccount_,
        address owner_
    ) external initializer {
        _setPreconfManager(preconfManager_);
        _setBlockTracker(blockTrackerContract_);
        _setOracleAccount(oracleAccount_);
        __Ownable_init(owner_);
        __Pausable_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Empty receive function to prevent unintended interactions.
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
     * @dev Processes a builder's commitment for a specific block number.
     * @param commitmentIndex The id of the commitment in the PreconfManager.
     * @param blockNumber The block number to be processed.
     * @param builder The address of the builder.
     * @param isSlash Determines whether the commitment should be slashed or rewarded.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function processBuilderCommitmentForBlockNumber(
        bytes32 commitmentIndex,
        uint256 blockNumber,
        address builder,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) external onlyOracle whenNotPaused {
        address blockWinner = _blockTrackerContract.getBlockWinner(blockNumber);
        require(blockWinner == builder, BuilderNotBlockWinner(blockWinner, builder));
        require(residualBidPercentAfterDecay <= 1e18, ResidualBidPercentAfterDecayExceedsMax(residualBidPercentAfterDecay));

        IPreconfManager.OpenedCommitment
            memory commitment = _preconfManager.getCommitment(commitmentIndex);
        if (
            commitment.committer == builder &&
            commitment.blockNumber == blockNumber
        ) {
            _processCommitment(
                commitmentIndex,
                isSlash,
                residualBidPercentAfterDecay
            );
        }
    }

    /// @dev Allows the owner to set the oracle account.
    function setOracleAccount(address newOracleAccount) external onlyOwner {
        _setOracleAccount(newOracleAccount);
    }

    /// @dev Allows the owner to set the preconf manager.
    function setPreconfManager(address newPreconfManager) external onlyOwner {
        _setPreconfManager(newPreconfManager);
    }

    /// @dev Allows the owner to set the block tracker.
    function setBlockTracker(address newBlockTracker) external onlyOwner {
        _setBlockTracker(newBlockTracker);
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
     * @dev Internal function to set the oracle account.
     * @param newOracleAccount The new address of the oracle account.
     */
    function _setOracleAccount(address newOracleAccount) internal {
        address oldOracleAccount = oracleAccount;
        oracleAccount = newOracleAccount;
        emit OracleAccountSet(oldOracleAccount, newOracleAccount);
    }

    /**
     * @dev Internal function to set the preconf manager.
     * @param newPreconfManager The new address of the preconf manager.
     */
    function _setPreconfManager(address newPreconfManager) internal {
        _preconfManager = IPreconfManager(newPreconfManager);
        emit PreconfManagerSet(newPreconfManager);
    }

    /**
     * @dev Internal function to set the block tracker.
     * @param newBlockTracker The new address of the block tracker.
     */
    function _setBlockTracker(address newBlockTracker) internal {
        _blockTrackerContract = IBlockTracker(newBlockTracker);
        emit BlockTrackerSet(newBlockTracker);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /**
     * @dev Internal function to process a commitment, either slashing or rewarding based on the commitment's state.
     * @param commitmentIndex The id of the commitment to be processed.
     * @param isSlash Determines if the commitment should be slashed or rewarded.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
     */
    function _processCommitment(
        bytes32 commitmentIndex,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) private {
        if (isSlash) {
            _preconfManager.initiateSlash(
                commitmentIndex,
                residualBidPercentAfterDecay
            );
        } else {
            _preconfManager.initiateReward(
                commitmentIndex,
                residualBidPercentAfterDecay
            );
        }
        // Emit an event that a commitment has been processed
        emit CommitmentProcessed(commitmentIndex, isSlash);
    }
}
