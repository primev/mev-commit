// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {EventHeightLib} from "../../utils/EventHeight.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

// TODO: split out to storage and interface, also test file.
// TODO: Don't modularize, instead just copy relevant reg/dereg logic from MevCommitAVS.
contract MevCommitMiddleware is Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    struct ValidatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
        address operator;
        uint256 priorityIndex;
    }

    struct OperatorRecord {
        bool exists;
        EventHeightLib.EventHeight deregRequestHeight;
    }

    mapping(bytes => ValidatorRecord) public blsPubkeyToValRecord;

    mapping(address => OperatorRecord) public operatorRecords;

    // TODO: invariant here is that no two validator records have the same priority for the same operator, 
    // and that operatorToPriorityIndexCounter[operator] number of records exist at any given time for an operator.

    mapping(address => uint256) public operatorToPriorityIndexCounter;

    uint256 public validatorDeregPeriodBlocks;
    uint256 public operatorDeregPeriodBlocks;

    function initialize() public initializer {
        __Ownable2Step_init();
        __Pausable_init();
        __UUPSUpgradeable_init();
    }

    constructor() {
        _disableInitializers();
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    // TODO: events.... and index everything

    // TODO: onlyOwner for now, but we can also whitelist potential operators
    function registerOperators(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _registerOperator(operators[i]);
        }
    }

    // TODO: onlyOwner for now, but we can also whitelist potential operators
    function requestOperatorDeregistrations(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _requestOperatorDeregistration(operators[i]);
        }
    }

    // TODO: onlyOwner for now, but we can also whitelist potential operators
    function deregisterOperators(address[] calldata operators) external whenNotPaused onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _deregisterOperator(operators[i]);
        }
    }

    // TODO: confirm this and other external functions can handle empty arrays
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

    function slashValidators(bytes[] calldata blsPubkeys) external onlyOwner {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _slashValidator(blsPubkeys[i]);
        }
    }

    // TODO: Param setters

    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    function _setOperatorRecord(address operator) internal {
        operatorRecords[operator] = OperatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            })
        });
    }

    // TODO: hook these into symbiotic core

    function _registerOperator(address operator) internal {
        require(!operatorRecords[operator].exists, "operator already registered");
        _setOperatorRecord(operator);
        // TODO: emit event
    }

    function _requestOperatorDeregistration(address operator) internal {
        require(operatorRecords[operator].exists, "operator not registered");
        EventHeightLib.set(operatorRecords[operator].deregRequestHeight, block.number);
        // TODO: emit event
    }

    function _deregisterOperator(address operator) internal {
        require(operatorRecords[operator].exists, "operator dereg not requested");
        require(_isOperatorDeregistered(operator), "dereg too soon");
        delete operatorRecords[operator];
        // TODO: emit event
    }

    function _setValRecord(bytes calldata blsPubkey, uint256 priorityIndex) internal {
        blsPubkeyToValRecord[blsPubkey] = ValidatorRecord({
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
        require(!blsPubkeyToValRecord[blsPubkey].exists, "val record already exists");
        _setValRecord(blsPubkey, operatorToPriorityIndexCounter[msg.sender]);
        ++operatorToPriorityIndexCounter[msg.sender];
    }

    function _replaceValRecord(bytes calldata newBlsPubkey, bytes calldata oldBlsPubkey) internal {
        require(blsPubkeyToValRecord[oldBlsPubkey].exists, "missing val record");
        require(blsPubkeyToValRecord[oldBlsPubkey].operator == msg.sender, "sender is not operator");
        require(_isValidatorDeregistered(oldBlsPubkey), "val record not deregistered");
        require(!blsPubkeyToValRecord[newBlsPubkey].exists, "val record already exists");

        uint256 priorityIndex = blsPubkeyToValRecord[oldBlsPubkey].priorityIndex;
        delete blsPubkeyToValRecord[oldBlsPubkey];
        _setValRecord(newBlsPubkey, priorityIndex);
    }

    // TODO: test newBlsPubkey could be the same as oldBlsPubkey
    function _swapValRecords(bytes calldata blsPubkey1, bytes calldata blsPubkey2) internal {
        require(blsPubkeyToValRecord[blsPubkey1].exists, "missing val record 1");
        require(blsPubkeyToValRecord[blsPubkey2].exists, "missing val record 2");

        require(msg.sender == blsPubkeyToValRecord[blsPubkey1].operator &&
            msg.sender == blsPubkeyToValRecord[blsPubkey2].operator, "sender is not operator");

        require(_isValidatorDeregistered(blsPubkey1), "val record 1 not deregistered");
        require(_isValidatorDeregistered(blsPubkey2), "val record 2 not deregistered");
            
        // swap priorities, reset dereg request heights
        uint256 priorityIndex1 = blsPubkeyToValRecord[blsPubkey1].priorityIndex;
        _setValRecord(blsPubkey1, blsPubkeyToValRecord[blsPubkey2].priorityIndex);
        _setValRecord(blsPubkey2, priorityIndex1);
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(blsPubkeyToValRecord[blsPubkey].exists, "missing validator record");
        require(blsPubkeyToValRecord[blsPubkey].operator == msg.sender, "sender is not operator");
        EventHeightLib.set(blsPubkeyToValRecord[blsPubkey].deregRequestHeight, block.number);
    }

    function _slashValidator(bytes calldata blsPubkey) internal {
        require(blsPubkeyToValRecord[blsPubkey].exists, "missing validator record");
        // TODO: slash operator with core
        // address operator = blsPubkeyToValRecord[blsPubkey].operator;
        _requestValDeregistration(blsPubkey); // TODO: determine if validator should be deregistered
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    function _isValidatorDeregistered(bytes calldata blsPubkey) internal view returns (bool) {
        return blsPubkeyToValRecord[blsPubkey].deregRequestHeight.exists && 
            block.number > validatorDeregPeriodBlocks + blsPubkeyToValRecord[blsPubkey].deregRequestHeight.blockHeight;
    }

    function _isOperatorDeregistered(address operator) internal view returns (bool) {
        return operatorRecords[operator].deregRequestHeight.exists && 
            block.number > operatorDeregPeriodBlocks + operatorRecords[operator].deregRequestHeight.blockHeight;
    }

    function _isValidatorOptedIn(bytes calldata blsPubkey) internal view returns (bool) {
        if (!blsPubkeyToValRecord[blsPubkey].exists) {
            return false;
        }
        if (blsPubkeyToValRecord[blsPubkey].deregRequestHeight.exists) {
            return false;
        }
        if (!operatorRecords[blsPubkeyToValRecord[blsPubkey].operator].exists) {
            return false;
        }
        if (operatorRecords[blsPubkeyToValRecord[blsPubkey].operator].deregRequestHeight.exists) {
            return false;
        }
        // TODO: check liquidity exists to slash if needed
        if (blsPubkeyToValRecord[blsPubkey].priorityIndex > 71) { // where 71 is some threshold defined by amount of liquidity
            return false;
        }
        return true;
    }
}
