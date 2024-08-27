// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {EventHeightLib} from "../../utils/EventHeight.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareStorage} from "./MevCommitMiddlewareStorage.sol";

// TODO: See if reputational val reg PR: https://github.com/primev/mev-commit/pull/131/files serves any inspiration in operator whitelisting. 
// TODO: add symbiotic core integration via lifecycle: https://docs.symbiotic.fi/core-modules/networks#staking-lifecycle
// TODO: determine if you need timestamping similar to cosmos sdk example. Edit yes you will for slashing. See "captureTimestamp". 
contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    
    // TODO: more modifiers similar to MevCommitAVS

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "only slash oracle");
        _;
    }

    // TODO: Define integration with individual vaults, and how you decide on "min stake" per validator
    // for each denom. Price oracle or hardcoded minStake? 

    // TODO: invariant here is that no two validator records have the same priority for the same operator, 
    // and that operatorToPriorityIndexCounter[operator] number of records exist at any given time for an operator.

    // TODO: Add things like network epoch duration, ref to core contracts, etc. 
    function initialize(
        uint256 _operatorDeregPeriodBlocks,
        uint256 _validatorDeregPeriodBlocks,
        address _slashOracle,
        address _owner
    ) public initializer {
        _setOperatorDeregPeriodBlocks(_operatorDeregPeriodBlocks);
        _setValidatorDeregPeriodBlocks(_validatorDeregPeriodBlocks);
        _setSlashOracle(_slashOracle);
        __Pausable_init();
        __UUPSUpgradeable_init();
        __Ownable2Step_init();
        transferOwnership(_owner);
    }

    constructor() {
        _disableInitializers();
    }

    // TODO: Make this whitelist instead and operators register themselves
    function registerOperators(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _registerOperator(operators[i]);
        }
    }

    // TODO: Make this whitelist instead and operators register themselves
    function requestOperatorDeregistrations(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _requestOperatorDeregistration(operators[i]);
        }
    }

    // TODO: Make this whitelist instead and operators register themselves
    function deregisterOperators(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _deregisterOperator(operators[i]);
        }
    }

    // TODO: confirm this and other external functions can handle empty arrays
    // TODO: confirm only operator can edit their own records. Does contract owner need access as well?
    // Be consistent with MevCommitAVS.
    function registerValidators(bytes[] calldata blsPubkeys) external whenNotPaused {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _addValRecord(blsPubkeys[i]);
        }
    }

    function requestValDeregistrations(bytes[] calldata blsPubkeys) external whenNotPaused {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _requestValDeregistration(blsPubkeys[i]);
        }
    }

    function replaceValRegistrations(bytes[] calldata newBlsPubkeys, bytes[] calldata oldBlsPubkeys) external whenNotPaused {
        require(newBlsPubkeys.length == oldBlsPubkeys.length, "invalid length");
        for (uint256 i = 0; i < newBlsPubkeys.length; i++) {
            _replaceValRecord(newBlsPubkeys[i], oldBlsPubkeys[i]);
        }
    }

    function swapValRegistrations(bytes[] calldata blsPubkeys1, bytes[] calldata blsPubkeys2) external whenNotPaused {
        require(blsPubkeys1.length == blsPubkeys2.length, "invalid length");
        for (uint256 i = 0; i < blsPubkeys1.length; i++) {
            _swapValRecords(blsPubkeys1[i], blsPubkeys2[i]);
        }
    }

    function slashValidators(bytes[] calldata blsPubkeys) external onlySlashOracle {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _slashValidator(blsPubkeys[i]);
        }
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner { _pause(); }

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external onlyOwner { _unpause(); }

    /// @dev Sets the operator deregistration period in blocks, restricted to contract owner.
    function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) external onlyOwner {
        _setOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_);
    }

    /// @dev Sets the validator deregistration period in blocks, restricted to contract owner.
    function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) external onlyOwner {
        _setValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_);
    }

    /// @dev Sets the slash oracle, restricted to contract owner.
    function setSlashOracle(address slashOracle_) external onlyOwner {
        _setSlashOracle(slashOracle_);
    }

    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    // TODO: hook these into symbiotic core
    function _registerOperator(address operator) internal {
        require(!operatorRecords[operator].exists, "operator already registered");
        operatorRecords[operator] = OperatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            priorityIndexCounter: 0
        });
        emit OperatorRegistered(operator);
    }

    function _requestOperatorDeregistration(address operator) internal {
        require(operatorRecords[operator].exists, "operator not registered");
        EventHeightLib.set(operatorRecords[operator].deregRequestHeight, block.number);
        emit OperatorDeregistrationRequested(operator);
    }

    function _deregisterOperator(address operator) internal {
        require(operatorRecords[operator].exists, "operator dereg not requested");
        require(_isOperatorDeregistered(operator), "dereg too soon");
        delete operatorRecords[operator];
        emit OperatorDeregistered(operator);
    }

    function _setValRecord(bytes calldata blsPubkey, uint256 priorityIndex) internal {
        validatorRecords[blsPubkey] = ValidatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            operator: msg.sender,
            priorityIndex: priorityIndex
        });
    }

    function _addValRecord(bytes calldata blsPubkey) internal {
        require(!validatorRecords[blsPubkey].exists, "val record already exists");
        _setValRecord(blsPubkey, operatorRecords[msg.sender].priorityIndexCounter);
        emit ValRecordAdded(blsPubkey, msg.sender, operatorRecords[msg.sender].priorityIndexCounter);
        ++operatorRecords[msg.sender].priorityIndexCounter;
    }

    function _replaceValRecord(bytes calldata newBlsPubkey, bytes calldata oldBlsPubkey) internal {
        require(validatorRecords[oldBlsPubkey].exists, "missing val record");
        require(validatorRecords[oldBlsPubkey].operator == msg.sender, "sender is not operator");
        require(_isValidatorDeregistered(oldBlsPubkey), "val record not deregistered");
        require(!validatorRecords[newBlsPubkey].exists, "val record already exists");

        uint256 priorityIndex = validatorRecords[oldBlsPubkey].priorityIndex;
        delete validatorRecords[oldBlsPubkey];
        _setValRecord(newBlsPubkey, priorityIndex);
        emit ValRecordReplaced(oldBlsPubkey, newBlsPubkey, msg.sender, priorityIndex);
    }

    // TODO: test newBlsPubkey could be the same as oldBlsPubkey
    function _swapValRecords(bytes calldata blsPubkey1, bytes calldata blsPubkey2) internal {
        require(validatorRecords[blsPubkey1].exists, "missing val record 1");
        require(validatorRecords[blsPubkey2].exists, "missing val record 2");

        require(msg.sender == validatorRecords[blsPubkey1].operator &&
            msg.sender == validatorRecords[blsPubkey2].operator, "sender is not operator");

        require(_isValidatorDeregistered(blsPubkey1), "val record 1 not deregistered");
        require(_isValidatorDeregistered(blsPubkey2), "val record 2 not deregistered");
            
        // swap priorities, reset dereg request heights
        uint256 priorityIndex1 = validatorRecords[blsPubkey1].priorityIndex;
        _setValRecord(blsPubkey1, validatorRecords[blsPubkey2].priorityIndex);
        _setValRecord(blsPubkey2, priorityIndex1);
        emit ValRecordsSwapped(blsPubkey1, blsPubkey2, msg.sender,
            // Log new stored priority indexes
            validatorRecords[blsPubkey1].priorityIndex, validatorRecords[blsPubkey2].priorityIndex);
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        require(validatorRecords[blsPubkey].operator == msg.sender, "sender is not operator");
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, validatorRecords[blsPubkey].priorityIndex);
    }

    function _slashValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        // TODO: slash operator with core
        address operator = validatorRecords[blsPubkey].operator;
        _requestValDeregistration(blsPubkey); // TODO: determine if validator should be deregistered
        emit ValidatorSlashed(blsPubkey, operator, validatorRecords[blsPubkey].priorityIndex);
    }

    /// @dev Internal function to set the operator deregistration period in blocks.
    function _setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) internal {
        operatorDeregPeriodBlocks = operatorDeregPeriodBlocks_;
        emit OperatorDeregPeriodBlocksSet(operatorDeregPeriodBlocks_);
    }

    /// @dev Internal function to set the validator deregistration period in blocks.
    function _setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) internal {
        validatorDeregPeriodBlocks = validatorDeregPeriodBlocks_;
        emit ValidatorDeregPeriodBlocksSet(validatorDeregPeriodBlocks_);
    }

    /// @dev Internal function to set the slash oracle.
    function _setSlashOracle(address slashOracle_) internal {
        slashOracle = slashOracle_;
        emit SlashOracleSet(slashOracle_);
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _isValidatorDeregistered(bytes calldata blsPubkey) internal view returns (bool) {
        return validatorRecords[blsPubkey].deregRequestHeight.exists && 
            block.number > validatorDeregPeriodBlocks + validatorRecords[blsPubkey].deregRequestHeight.blockHeight;
    }

    function _isOperatorDeregistered(address operator) internal view returns (bool) {
        return operatorRecords[operator].deregRequestHeight.exists && 
            block.number > operatorDeregPeriodBlocks + operatorRecords[operator].deregRequestHeight.blockHeight;
    }

    function _isValidatorOptedIn(bytes calldata blsPubkey) internal view returns (bool) {
        if (!validatorRecords[blsPubkey].exists) {
            return false;
        }
        if (validatorRecords[blsPubkey].deregRequestHeight.exists) {
            return false;
        }
        if (!operatorRecords[validatorRecords[blsPubkey].operator].exists) {
            return false;
        }
        if (operatorRecords[validatorRecords[blsPubkey].operator].deregRequestHeight.exists) {
            return false;
        }
        // TODO: check liquidity exists to slash if needed
        if (validatorRecords[blsPubkey].priorityIndex > 71) { // where 71 is some threshold defined by amount of liquidity
            return false;
        }
        return true;
    }
}
