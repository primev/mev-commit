// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {PreConfCommitmentStore} from "./PreConfCommitmentStore.sol";
import {IProviderRegistry} from "./interfaces/IProviderRegistry.sol";
import {IPreConfCommitmentStore} from './interfaces/IPreConfCommitmentStore.sol';
import {IBidderRegistry} from './interfaces/IBidderRegistry.sol';
import {IBlockTracker} from "./interfaces/IBlockTracker.sol";

/// @title Oracle Contract
/// @author Kartik Chopra
/// @notice This contract is for fetching L1 Ethereum Block Data

/**
 * @title Oracle - A contract for Fetching L1 Block Builder Info and Block Data.
 * @dev This contract serves as an oracle to fetch and process Ethereum Layer 1 block data.
 */
contract Oracle is OwnableUpgradeable, UUPSUpgradeable {
    /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    // To shutup the compiler
    /// @dev Empty receive function to silence compiler warnings about missing payable functions.
    receive() external payable {
        // Empty receive function
    }

    /**
     * @dev Fallback function to revert all calls, ensuring no unintended interactions.
     */
    fallback() external payable {
        revert("Invalid call");
    }

    /// @dev Reference to the PreConfCommitmentStore contract interface.
    IPreConfCommitmentStore private preConfContract;

    /// @dev Reference to the BlockTracker contract interface.
    IBlockTracker private blockTrackerContract;

    function _authorizeUpgrade(address) internal override onlyOwner {}

    /**
     * @dev Initializes the contract with a PreConfirmations contract.
     * @param _preConfContract The address of the pre-confirmations contract.
     * @param _owner Owner of the contract, explicitly needed since contract is deployed with create2 factory.
     */
    function initialize(
        address _preConfContract,
        address _blockTrackerContract,
        address _owner
    ) external initializer {
        preConfContract = IPreConfCommitmentStore(_preConfContract);
        blockTrackerContract = IBlockTracker(_blockTrackerContract);
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Event emitted when a commitment is processed.
    event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash);

    // Function to receive and process the block data (this would be automated in a real-world scenario)
    /**
     * @dev Processes a builder's commitment for a specific block number.
     * @param commitmentIndex The id of the commitment in the PreConfCommitmentStore.
     * @param blockNumber The block number to be processed.
     * @param builder The address of the builder.
     * @param isSlash Determines whether the commitment should be slashed or rewarded.
     */
    function processBuilderCommitmentForBlockNumber(
        bytes32 commitmentIndex,
        uint256 blockNumber,
        address builder,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) external onlyOwner {
        address winner = blockTrackerContract.getBlockWinner(blockNumber);
        require(
            winner == builder,
            "Builder is not the winner of the block"
        );
        require(
            residualBidPercentAfterDecay <= 100,
            "Residual bid after decay cannot be greater than 100 percent"
        );
        IPreConfCommitmentStore.PreConfCommitment
            memory commitment = preConfContract.getCommitment(commitmentIndex);
        if (
            commitment.commiter == builder &&
            commitment.blockNumber == blockNumber
        ) {
            processCommitment(
                commitmentIndex,
                isSlash,
                residualBidPercentAfterDecay
            );
        }
    }

    /**
     * @dev Internal function to process a commitment, either slashing or rewarding based on the commitment's state.
     * @param commitmentIndex The id of the commitment to be processed.
     * @param isSlash Determines if the commitment should be slashed or rewarded.
     */
    function processCommitment(
        bytes32 commitmentIndex,
        bool isSlash,
        uint256 residualBidPercentAfterDecay
    ) private {
        if (isSlash) {
            preConfContract.initiateSlash(
                commitmentIndex,
                residualBidPercentAfterDecay
            );
        } else {
            preConfContract.initiateReward(
                commitmentIndex,
                residualBidPercentAfterDecay
            );
        }
        // Emit an event that a commitment has been processed
        emit CommitmentProcessed(commitmentIndex, isSlash);
    }
}
