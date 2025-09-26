// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {RocketMinipoolInterface} from "rocketpool/contracts/interface/minipool/RocketMinipoolInterface.sol";
import {MinipoolStatus} from "rocketpool/contracts/types/MinipoolStatus.sol";
import {RocketStorageInterface} from "rocketpool/contracts/interface/RocketStorageInterface.sol";

import {Errors} from "../../utils/Errors.sol";
import {IRocketMinipoolRegistry} from "../../interfaces/IRocketMinipoolRegistry.sol";
import {RocketMinipoolRegistryStorage} from "./RocketMinipoolRegistryStorage.sol";

/// @title RocketMinipoolRegistry
/// @notice This contract serves as the entrypoint for operators to register with
/// the mev-commit protocol via Rocketpool minipools.
contract RocketMinipoolRegistry is IRocketMinipoolRegistry, RocketMinipoolRegistryStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, IRocketMinipoolRegistry.OnlyFreezeOracle());
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are valid length.
    modifier onlyValidBLSPubKeys(bytes[] calldata valPubKeys) {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            require(valPubKeys[i].length == 48, IRocketMinipoolRegistry.InvalidBLSPubKeyLength(48, valPubKeys[i].length));
        }
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Receive function to prevent unintended contract interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to prevent unintended contract interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }
    
    function initialize(address owner, address freezeOracle, address unfreezeReceiver, address rocketStorage, uint256 unfreezeFee, uint64 deregistrationPeriod) external initializer {
        __Ownable_init(owner);
        __Pausable_init();
        __UUPSUpgradeable_init();
        _setFreezeOracle(freezeOracle);
        _setUnfreezeReceiver(unfreezeReceiver);
        _setUnfreezeFee(unfreezeFee);
        _setRocketStorage(rocketStorage);
        _setDeregistrationPeriod(deregistrationPeriod);
    }

    function registerValidators(bytes[] calldata valPubKeys) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _registerValidator(valPubKeys[i]);
            
        }
    }

    function requestValidatorDeregistration(bytes[] calldata valPubKeys) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestValidatorDeregistration(valPubKeys[i]);
        }
    }

    /// @dev Deregister validators. Can only be called once the deregistration period has passed from time of request.
    function deregisterValidators(bytes[] calldata valPubKeys) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterValidator(valPubKeys[i]);
        }
    }

    /// @dev Allows the freeze oracle account to freeze validators which disobey the mev-commit protocol.
    function freeze(bytes[] calldata valPubKeys) external whenNotPaused onlyFreezeOracle() {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _freeze(valPubKeys[i]);
        }
    }

    /// @dev Allows any account to unfreeze validators which have been frozen, for a fee.
    function unfreeze(bytes[] calldata valPubKeys) external payable whenNotPaused {
        uint256 requiredFee = unfreezeFee * valPubKeys.length;
        require(msg.value >= requiredFee, UnfreezeFeeRequired(requiredFee));
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _unfreeze(valPubKeys[i]);
        }
        (bool success, ) = unfreezeReceiver.call{value: requiredFee}("");
        require(success, UnfreezeTransferFailed());
        uint256 excessFee = msg.value - requiredFee;
        if (excessFee != 0) {
            (bool successRefund, ) = msg.sender.call{value: excessFee}("");
            require(successRefund, RefundFailed());
        }
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner whenNotPaused {
        _pause();
    }

    /// @dev Unpauses the contract, restricted to contract owner.   
    function unpause() external onlyOwner whenPaused {
        _unpause();
    }

    /// @dev Sets the deregistration period, restricted to contract owner.
    function setDeregistrationPeriod(uint64 newDeregistrationPeriod) external onlyOwner {
        _setDeregistrationPeriod(newDeregistrationPeriod);
    }

    /// @dev Sets the rocket storage, restricted to contract owner.    
    function setRocketStorage(address newRocketStorage) external onlyOwner {
        _setRocketStorage(newRocketStorage);
    }

    /// @dev Sets the freeze oracle, restricted to contract owner.
    function setFreezeOracle(address newFreezeOracle) external onlyOwner {
        _setFreezeOracle(newFreezeOracle);
    }

    /// @dev Sets the unfreeze receiver, restricted to contract owner.
    function setUnfreezeReceiver(address newUnfreezeReceiver) external onlyOwner {
        _setUnfreezeReceiver(newUnfreezeReceiver);
    }

    /// @dev Sets the unfreeze fee, restricted to contract owner.
    function setUnfreezeFee(uint256 newUnfreezeFee) external onlyOwner {
        _setUnfreezeFee(newUnfreezeFee);
    }

    /// @dev Unfreezes validators, restricted to contract owner.
    function ownerUnfreeze(bytes[] calldata valPubKeys) external onlyOwner {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _unfreeze(valPubKeys[i]);
        }
    }

    /// @dev Returns the time at which a validator can be deregistered.
    function getEligibleTimeForDeregistration(bytes calldata validatorPubkey) external view returns (uint64) {
        uint64 deregistrationTime = validatorRegistrations[validatorPubkey].deregTimestamp;
        return deregistrationTime != 0 ? (deregistrationTime + deregistrationPeriod) : 0;
    }

    /// @dev Returns both the node address and the withdrawal address of the key's minipool, as these addresses both have minipool permissions.
    function getValidOperatorsForKey(bytes calldata validatorPubkey) external view returns (address, address) {
        address minipool = getMinipoolFromPubkey(validatorPubkey);
        address nodeAddress = getNodeAddressFromMinipool(minipool);
        return (nodeAddress, rocketStorage.getNodeWithdrawalAddress(nodeAddress));
    }

    /// @dev Returns validator registration info.
    function getValidatorRegInfo(bytes calldata valPubKey) external view returns (ValidatorRegistration memory) {
        return validatorRegistrations[valPubKey];
    }

    /// @dev Checks if a validator is opted-in.
    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valPubKey);
    }

        /// @dev Returns the minipool for a validator.
    function getMinipoolFromPubkey(bytes calldata validatorPubkey) public view returns (address) {
        return rocketStorage.getAddress(keccak256(abi.encodePacked("validator.minipool", validatorPubkey)));
    }

    /// @dev Returns the node address from a validator's minipool.
    function getNodeAddressFromPubkey(bytes calldata validatorPubkey) public view returns (address) {
        address minipool = getMinipoolFromPubkey(validatorPubkey);
        return getNodeAddressFromMinipool(minipool);
    }

    /// @dev Returns true if a minipool is active.
    function isMinipoolActive(address minipool) public view returns (bool) {
        return RocketMinipoolInterface(minipool).getStatus() == MinipoolStatus.Staking;
    }

    /// @dev Fetches the minipool from a validator's pubkey and returns true if caller is either the minipool's node address or node's withdrawal address.
    function isOperatorValidForKey(bytes calldata validatorPubkey) public view returns (bool) {
        address minipool = getMinipoolFromPubkey(validatorPubkey);
        address nodeAddress = getNodeAddressFromMinipool(minipool);
        return (nodeAddress == msg.sender || rocketStorage.getNodeWithdrawalAddress(nodeAddress) == msg.sender);
    }

    function isValidatorRegistered(bytes calldata validatorPubkey) public view returns (bool) {
        return validatorRegistrations[validatorPubkey].exists;
    }

    /// @dev Returns the node address of a minipool.
    function getNodeAddressFromMinipool(address minipool) public view returns (address) {
        return RocketMinipoolInterface(minipool).getNodeAddress();
    }

    /// @dev Registers a validator.
    function _registerValidator(bytes calldata valPubKey) internal {
        address minipool = getMinipoolFromPubkey(valPubKey);
        require(minipool != address(0), NoMinipoolForKey(valPubKey));
        require(_isOperatorValid(minipool), NotMinipoolOperator(valPubKey));
        require(isMinipoolActive(minipool), MinipoolNotActive(valPubKey));
        require(!isValidatorRegistered(valPubKey), ValidatorAlreadyRegistered(valPubKey));
        address nodeAddress = getNodeAddressFromMinipool(minipool);
        IRocketMinipoolRegistry.ValidatorRegistration storage reg = validatorRegistrations[valPubKey];
        reg.exists = true;
        emit ValidatorRegistered(valPubKey, nodeAddress);
    }

    /// @dev Requests deregistration for a validator.
    function _requestValidatorDeregistration(bytes calldata valPubKey) internal {
        address nodeAddress = getNodeAddressFromPubkey(valPubKey);
        require(_isOperatorValid(nodeAddress), NotMinipoolOperator(valPubKey));
        require(isValidatorRegistered(valPubKey), ValidatorNotRegistered(valPubKey));
        IRocketMinipoolRegistry.ValidatorRegistration storage reg = validatorRegistrations[valPubKey];
        require(reg.deregTimestamp == 0, DeregRequestAlreadyExists(valPubKey));
        reg.deregTimestamp = uint64(block.timestamp);
        emit ValidatorDeregistrationRequested(valPubKey, nodeAddress);
    }
    
    function _deregisterValidator(bytes calldata valPubKey) internal {
        address nodeAddress = getNodeAddressFromPubkey(valPubKey);
        require(_isOperatorValid(nodeAddress), NotMinipoolOperator(valPubKey));
        IRocketMinipoolRegistry.ValidatorRegistration storage reg = validatorRegistrations[valPubKey];
        require(reg.deregTimestamp != 0, DeregRequestDoesNotExist(valPubKey));
        require(uint64(block.timestamp) > reg.deregTimestamp + deregistrationPeriod, DeregistrationTooSoon(valPubKey));
        require(reg.freezeTimestamp == 0, FrozenValidatorCannotDeregister(valPubKey));
        delete validatorRegistrations[valPubKey];
        emit ValidatorDeregistered(valPubKey, nodeAddress);
    }

    function _freeze(bytes calldata valPubKey) internal {
        require(isValidatorRegistered(valPubKey), ValidatorNotRegistered(valPubKey));
        require(validatorRegistrations[valPubKey].freezeTimestamp == 0, ValidatorAlreadyFrozen(valPubKey));
        validatorRegistrations[valPubKey].freezeTimestamp = uint64(block.timestamp);
        emit ValidatorFrozen(valPubKey);
    }

    function _unfreeze(bytes calldata valPubKey) internal {
        IRocketMinipoolRegistry.ValidatorRegistration storage regInfo = validatorRegistrations[valPubKey];
        require(regInfo.freezeTimestamp != 0, ValidatorNotFrozen(valPubKey)); 
        regInfo.freezeTimestamp = 0;
        emit ValidatorUnfrozen(valPubKey);
    }

    function _setFreezeOracle(address newFreezeOracle) internal {
        freezeOracle = newFreezeOracle;
    }

    function _setUnfreezeReceiver(address newUnfreezeReceiver) internal {
        unfreezeReceiver = newUnfreezeReceiver;
    }

    function _setUnfreezeFee(uint256 newUnfreezeFee) internal {
        unfreezeFee = newUnfreezeFee;
    }

    function _setRocketStorage(address newRocketStorage) internal {
        rocketStorage = RocketStorageInterface(newRocketStorage);
    }

    function _setDeregistrationPeriod(uint64 newDeregistrationPeriod) internal {
        deregistrationPeriod = newDeregistrationPeriod;
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        if (!isValidatorRegistered(valPubKey)) return false;
        if (validatorRegistrations[valPubKey].freezeTimestamp != 0) return false;
        if (validatorRegistrations[valPubKey].deregTimestamp != 0) return false;
        if (getMinipoolFromPubkey(valPubKey) == address(0)) return false;
        if (!isMinipoolActive(getMinipoolFromPubkey(valPubKey))) return false;
        return true;
    }

    /// @dev Returns true if caller is either the minipool's node address or node'swithdrawal address.
    function _isOperatorValid(address operator) internal view returns (bool) {
        return (operator == msg.sender || rocketStorage.getNodeWithdrawalAddress(operator) == msg.sender);
    }
}