// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

import {MevCommitBappStorage} from "./MevCommitBappStorage.sol";
import {TimestampOccurrence} from "../../utils/Occurrence.sol";
import {OwnableBasedApp} from "@ssv/src/middleware/modules/core+roles/OwnableBasedApp.sol";

contract MevCommitBapp is MevCommitBappStorage, OwnableBasedApp {

    event ValidatorRegistered(bytes indexed pubkey, address indexed registrar);
    event ValidatorDeregistrationRequested(bytes indexed pubkey, address indexed deregistrar);
    event ValidatorDeregistered(bytes indexed pubkey, address indexed deregistrar);
    event ValidatorFrozen(bytes indexed pubkey, address indexed freezer);
    event ValidatorUnfrozen(bytes indexed pubkey, address indexed unfreezer);

    error AlreadyWhitelisted();
    error ZeroAddress();
    error NotWhitelisted();
    error NonWhitelistedCaller();
    error NotOptedInCaller();
    error ValidatorAlreadyRegistered(bytes pubkey);
    error ValidatorNotRegistered(bytes pubkey);
    error ValidatorAlreadyRequestedDeregistration(bytes pubkey);
    error ValidatorNotFrozen(bytes pubkey);
    error UnfreezeTooSoon();
    error UnfreezeTransferFailed();
    error RefundFailed();
    error UnfreezeFeeRequired(uint256 requiredFee);

    constructor(address _basedAppManager, address owner) OwnableBasedApp(_basedAppManager, owner) {}

    function addWhitelisted(address account) external onlyOwner {
        if (isWhitelisted[account]) revert AlreadyWhitelisted();
        if (account == address(0)) revert ZeroAddress();
        isWhitelisted[account] = true;
    }

    function removeWhitelisted(address account) external onlyOwner {
        if (!isWhitelisted[account]) revert NotWhitelisted();
        delete isWhitelisted[account];
    }

    function pause() external onlyOwner {
        _pause();
    }
    function unpause() external onlyOwner {
        _unpause();
    }

    function optInToBApp(
        uint32, /*strategyId*/
        address[] calldata, /*tokens*/
        uint32[] calldata, /*obligationPercentages*/
        bytes calldata /*data*/
    ) external override onlySSVBasedAppManager whenNotPaused returns (bool success) {
        if (!isWhitelisted[msg.sender]) revert NonWhitelistedCaller();
        isOptedIn[msg.sender] = true;
        // Before slashing is enabled, strategies are irrelevant to mev-commit.
        return true;
    }

    /// After opting-in, EOAs must register L1 validator pubkeys that attest to the rules of mev-commit.
    function registerValidatorPubkeys(bytes[] calldata pubkeys) external whenNotPaused {
        if (!isWhitelisted[msg.sender]) revert NonWhitelistedCaller();
        if (!isOptedIn[msg.sender]) revert NotOptedInCaller();
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes calldata pubkey = pubkeys[i];
            if (validatorRecords[pubkey].exists) revert ValidatorAlreadyRegistered(pubkey);
            validatorRecords[pubkey] = ValidatorRecord({
                exists: true,
                registrar: msg.sender,
                freezeOccurrence: TimestampOccurrence.Occurrence({
                    exists: false,
                    timestamp: 0
                }),
                deregRequestOccurrence: TimestampOccurrence.Occurrence({
                    exists: false,
                    timestamp: 0
                })
            });
            emit ValidatorRegistered(pubkey, msg.sender);
        }
    }

    function requestDeregistrations(bytes[] calldata pubkeys) external whenNotPaused {
        if (!isWhitelisted[msg.sender]) revert NonWhitelistedCaller();
        if (!isOptedIn[msg.sender]) revert NotOptedInCaller();
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes calldata pubkey = pubkeys[i];
            if (!validatorRecords[pubkey].exists) revert ValidatorNotRegistered(pubkey);
            if (validatorRecords[pubkey].deregRequestOccurrence.exists) revert ValidatorAlreadyRequestedDeregistration(pubkey);
            TimestampOccurrence.captureOccurrence(validatorRecords[pubkey].deregRequestOccurrence);
            emit ValidatorDeregistrationRequested(pubkey, msg.sender);
        }
    }

    function deregisterValidators(bytes[] calldata pubkeys) external whenNotPaused {
        if (!isWhitelisted[msg.sender]) revert NonWhitelistedCaller();
        if (!isOptedIn[msg.sender]) revert NotOptedInCaller();
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes calldata pubkey = pubkeys[i];
            delete validatorRecords[pubkey];
            emit ValidatorDeregistered(pubkey, msg.sender);
        }
    }

    function freezeValidators(bytes[] calldata pubkeys) external onlyOwner {
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes calldata pubkey = pubkeys[i];
            if (!validatorRecords[pubkey].exists) revert ValidatorNotRegistered(pubkey);
            TimestampOccurrence.captureOccurrence(validatorRecords[pubkey].freezeOccurrence);
            emit ValidatorFrozen(pubkey, msg.sender);
        }
    }

    function unfreeze(bytes[] calldata valPubKeys) external payable whenNotPaused() {
        uint256 requiredFee = unfreezeFee * valPubKeys.length;
        require(msg.value >= requiredFee, UnfreezeFeeRequired(requiredFee));
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            bytes calldata pubkey = valPubKeys[i];
            ValidatorRecord storage record = validatorRecords[pubkey];
            require(record.exists, ValidatorNotRegistered(pubkey));
            require(record.freezeOccurrence.exists, ValidatorNotFrozen(pubkey)); 
            require(block.timestamp > record.freezeOccurrence.timestamp + unfreezePeriod, UnfreezeTooSoon());
            TimestampOccurrence.del(record.freezeOccurrence);
            emit ValidatorUnfrozen(pubkey, record.registrar);
        }
        (bool success, ) = unfreezeReceiver.call{value: requiredFee}("");
        require(success, UnfreezeTransferFailed());
        uint256 excessFee = msg.value - requiredFee;
        if (excessFee != 0) {
            (bool successRefund, ) = msg.sender.call{value: excessFee}("");
            require(successRefund, RefundFailed());
        }
    }

    function isValidatorOptedIn(bytes calldata pubkey) public view returns (bool) {
        ValidatorRecord storage record = validatorRecords[pubkey];
        if (!record.exists) return false;
        if (!isWhitelisted[record.registrar]) return false;
        if (!isOptedIn[record.registrar]) return false;
        if (record.freezeOccurrence.exists) return false;
        if (record.deregRequestOccurrence.exists) return false;
        return true;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
} 
