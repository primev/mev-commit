// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {PreConfCommitmentStore} from "./PreConfCommitmentStore.sol";
import {IProviderRegistry} from "./interfaces/IProviderRegistry.sol";
import {IPreConfCommitmentStore} from './interfaces/IPreConfCommitmentStore.sol';
import {IBidderRegistry} from './interfaces/IBidderRegistry.sol';
import {IBlockTracker} from "./interfaces/IBlockTracker.sol";

/// @title Oracle Contract
/// @author Kartik Chopra
/// @notice This contract is for settling commitments made by providers.

/**
 * @title Oracle - A contract for Fetching L1 Block Builder Info and Block Data.
 * @dev This contract serves as an oracle to fetch and process Ethereum Layer 1 block data.
 */
contract Oracle is Ownable2StepUpgradeable, UUPSUpgradeable {
    /// @dev Maps builder names to their respective Ethereum addresses.
    mapping(string => address) public blockBuilderNameToAddress;

    /// @dev Permissioned address of the oracle account.
    address public oracleAccount;

    /// @dev Modifier to ensure that the sender is the oracle account.
    modifier onlyOracle() {
        require(msg.sender == oracleAccount, "sender isn't oracle account");
        _;
    }

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
     * @param preConfContract_ The address of the pre-confirmations contract.
     * @param blockTrackerContract_ The address of the block tracker contract.
     * @param oracleAccount_ The address of the oracle account.
     * @param owner_ Owner of the contract, explicitly needed since contract is deployed with create2 factory.
     */
    function initialize(
        address preConfContract_,
        address blockTrackerContract_,
        address oracleAccount_,
        address owner_
    ) external initializer {
        preConfContract = IPreConfCommitmentStore(preConfContract_);
        blockTrackerContract = IBlockTracker(blockTrackerContract_);
        _setOracleAccount(oracleAccount_);
        __Ownable_init(owner_);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Event emitted when a commitment is processed.
    event CommitmentProcessed(bytes32 indexed commitmentIndex, bool isSlash);

    /// @dev Event emitted when the oracle account is set.
    event OracleAccountSet(address indexed oldOracleAccount, address indexed newOracleAccount);

    // Function to receive and process the block data (this would be automated in a real-world scenario)
    /**
     * @dev Processes a builder's commitment for a specific block number.
     * @param commitmentIndex The id of the commitment in the PreConfCommitmentStore.
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
    ) external onlyOracle {
        require(
            blockTrackerContract.getBlockWinner(blockNumber) == builder,
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

    /// @dev Allows the owner to set the oracle account.
    function setOracleAccount(address newOracleAccount) external onlyOwner {
        _setOracleAccount(newOracleAccount);
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
     * @dev Internal function to process a commitment, either slashing or rewarding based on the commitment's state.
     * @param commitmentIndex The id of the commitment to be processed.
     * @param isSlash Determines if the commitment should be slashed or rewarded.
     * @param residualBidPercentAfterDecay The residual bid percent after decay.
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
