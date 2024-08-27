// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {EventHeightLib} from "../utils/EventHeight.sol";

interface IMevCommitMiddleware {

    struct ValidatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        address operator;
        uint256 priorityIndex;
    }

    struct OperatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        uint256 priorityIndexCounter;
    }

    /// @notice Emmitted when an operator is registered
    event OperatorRegistered(address indexed operator);

    /// @notice Emmitted when an operator requests deregistration
    event OperatorDeregistrationRequested(address indexed operator);

    /// @notice Emmitted when an operator is deregistered
    event OperatorDeregistered(address indexed operator);

    /// @notice Emmitted when a validator record is added to state
    event ValRecordAdded(bytes indexed blsPubkey, address indexed operator,
        uint256 indexed priorityIndex);

    /// @notice Emmitted when a validator record is replaced for a certain priority index
    event ValRecordReplaced(bytes indexed oldBlsPubkey, bytes indexed newBlsPubkey,
        address indexed operator, uint256 priorityIndex);

    /// @notice Emmitted when two validator records swap priority indexes
    event ValRecordsSwapped(bytes indexed blsPubkey1, bytes indexed blsPubkey2, address indexed operator,
        uint256 newPriorityIndex1, uint256 newPriorityIndex2);
    
    /// @notice Emmitted when validator deregistration is requested
    event ValidatorDeregistrationRequested(bytes indexed blsPubkey, address indexed operator,
        uint256 indexed priorityIndex);

    /// @notice Emmitted when a validator is slashed
    event ValidatorSlashed(bytes indexed blsPubkey, address indexed operator, uint256 indexed priorityIndex);

    /// @notice Emmitted when the operator deregistration period in blocks is set
    event OperatorDeregPeriodBlocksSet(uint256 operatorDeregPeriodBlocks);

    /// @notice Emmitted when the validator deregistration period in blocks is set
    event ValidatorDeregPeriodBlocksSet(uint256 validatorDeregPeriodBlocks);

    /// @notice Emmitted when the slash oracle is set
    event SlashOracleSet(address slashOracle);
}
