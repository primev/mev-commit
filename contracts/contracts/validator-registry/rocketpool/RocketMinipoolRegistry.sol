// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {EIP712Upgradeable} from "@openzeppelin/contracts-upgradeable/utils/cryptography/EIP712Upgradeable.sol";
import {SignatureChecker} from "@openzeppelin/contracts/utils/cryptography/SignatureChecker.sol";

import {RocketMinipoolInterface} from "rocketpool/contracts/interface/minipool/RocketMinipoolInterface.sol";
import {MinipoolStatus} from "rocketpool/contracts/types/MinipoolStatus.sol";
import {RocketStorageInterface} from "rocketpool/contracts/interface/RocketStorageInterface.sol";

import {Errors} from "../../utils/Errors.sol";
import {IRocketMinipoolRegistry} from "../../interfaces/IRocketMinipoolRegistry.sol";
import {RocketMinipoolRegistryStorage} from "./RocketMinipoolRegistryStorage.sol";

/// @title RocketMinipoolRegistry
/// @notice This contract serves as the entrypoint for operators to register with
/// the mev-commit protocol via Rocketpool minipools.
contract RocketMinipoolRegistry is
    IRocketMinipoolRegistry,
    RocketMinipoolRegistryStorage,
    Ownable2StepUpgradeable,
    PausableUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable,
    EIP712Upgradeable
{
    // We keep payloads compact by hashing the pubkeys array to a bytes32.
    bytes32 private constant _REGISTER_TYPEHASH =
        keccak256("Register(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");
    bytes32 private constant _DEREG_REQ_TYPEHASH =
        keccak256("DeregRequest(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");
    bytes32 private constant _DEREG_TYPEHASH =
        keccak256("Deregister(bytes32 pubkeysHash,address executor,uint256 nonce,uint256 deadline)");

    modifier onlyFreezeOracle() {
        require(msg.sender == freezeOracle, IRocketMinipoolRegistry.OnlyFreezeOracle());
        _;
    }

    /// @dev Modifier to confirm all provided BLS pubkeys are valid length.
    modifier onlyValidBLSPubKeys(bytes[] calldata valPubKeys) {
        uint256 len = valPubKeys.length;
        require(len > 0, NoKeysProvided());
        for (uint256 i = 0; i < len; ++i) {
            require(
                valPubKeys[i].length == 48, IRocketMinipoolRegistry.InvalidBLSPubKeyLength(48, valPubKeys[i].length)
            );
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

    function initialize(
        address owner,
        address freezeOracle,
        address unfreezeReceiver,
        address rocketStorage,
        uint256 unfreezeFee,
        uint64 deregistrationPeriod
    ) external initializer {
        __Ownable_init(owner);
        __Pausable_init();
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();
        __EIP712_init("RocketMinipoolRegistry", "1");
        _setFreezeOracle(freezeOracle);
        _setUnfreezeReceiver(unfreezeReceiver);
        _setUnfreezeFee(unfreezeFee);
        _setRocketStorage(rocketStorage);
        _setDeregistrationPeriod(deregistrationPeriod);
    }

    /// @notice Allows Minipool's withdrawal address to register validators.
    function registerValidators(bytes[] calldata valPubKeys) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        address nodeAddress0 = _getFirstNodeAddrAndValidate(valPubKeys[0]);
        _registerValidator(valPubKeys[0]);
        emit ValidatorRegistered(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _registerValidator(valPubKeys[i]);
            emit ValidatorRegistered(valPubKeys[i], nodeAddress0);
        }
    }

    /// @notice Allows a user to register using a signature.
    function registerValidatorsWithSig(bytes[] calldata valPubKeys, bytes calldata signature, uint256 deadline)
        external
        onlyValidBLSPubKeys(valPubKeys)
        whenNotPaused
    {
        bytes32 pkHash = pubkeysHash(valPubKeys);
        address nodeAddress0 = getNodeAddressFromPubkey(valPubKeys[0]);

        // 2) Pull nonce & check deadline
        uint256 nonce = nonces[nodeAddress0];
        require(block.timestamp <= deadline, ExpiredSignature());

        bytes32 digest =
            _hashTypedDataV4(keccak256(abi.encode(_REGISTER_TYPEHASH, pkHash, msg.sender, nonce, deadline)));

        // 6) Verify signature came from the node
        require(SignatureChecker.isValidSignatureNow(nodeAddress0, digest, signature), InvalidSignature());
        unchecked {
            nonces[nodeAddress0] = nonce + 1;
        }

        _registerValidator(valPubKeys[0]);
        emit ValidatorRegistered(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            //still need to check that the node address is the same for all validators in the batch
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _registerValidator(valPubKeys[i]);
            emit ValidatorRegistered(valPubKeys[i], nodeAddress0);
        }
    }

    /// @notice Allows Minipool's withdrawal address to request deregistration for validators.
    function requestValidatorDeregistration(bytes[] calldata valPubKeys)
        external
        onlyValidBLSPubKeys(valPubKeys)
        whenNotPaused
    {
        address nodeAddress0 = _getFirstNodeAddrAndValidate(valPubKeys[0]);
        _requestValidatorDeregistration(valPubKeys[0]);
        emit ValidatorDeregistrationRequested(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _requestValidatorDeregistration(valPubKeys[i]);
            emit ValidatorDeregistrationRequested(valPubKeys[i], nodeAddress0);
        }
    }

    function requestValidatorDeregistrationWithSig(
        bytes[] calldata valPubKeys,
        bytes calldata signature,
        uint256 deadline
    ) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        bytes32 pkHash = pubkeysHash(valPubKeys);
        address nodeAddress0 = getNodeAddressFromPubkey(valPubKeys[0]);

        uint256 nonce = nonces[nodeAddress0];
        require(block.timestamp <= deadline, ExpiredSignature());

        bytes32 digest =
            _hashTypedDataV4(keccak256(abi.encode(_DEREG_REQ_TYPEHASH, pkHash, msg.sender, nonce, deadline)));

        require(SignatureChecker.isValidSignatureNow(nodeAddress0, digest, signature), InvalidSignature());
        unchecked {
            nonces[nodeAddress0] = nonce + 1;
        }

        _requestValidatorDeregistration(valPubKeys[0]);
        emit ValidatorDeregistrationRequested(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _requestValidatorDeregistration(valPubKeys[i]);
            emit ValidatorDeregistrationRequested(valPubKeys[i], nodeAddress0);
        }
    }

    /// @notice Allows Minipool's withdrawal address to eregister validators. Can only be called once the deregistration period has passed from time of request.
    function deregisterValidators(bytes[] calldata valPubKeys) external onlyValidBLSPubKeys(valPubKeys) whenNotPaused {
        address nodeAddress0 = _getFirstNodeAddrAndValidate(valPubKeys[0]);
        _deregisterValidator(valPubKeys[0]);
        emit ValidatorDeregistered(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _deregisterValidator(valPubKeys[i]);
            emit ValidatorDeregistered(valPubKeys[i], nodeAddress0);
        }
    }

    function deregisterValidatorsWithSig(bytes[] calldata valPubKeys, bytes calldata signature, uint256 deadline)
        external
        onlyValidBLSPubKeys(valPubKeys)
        whenNotPaused
    {
        bytes32 pkHash = pubkeysHash(valPubKeys);
        address nodeAddress0 = getNodeAddressFromPubkey(valPubKeys[0]);

        uint256 nonce = nonces[nodeAddress0];
        require(block.timestamp <= deadline, ExpiredSignature());

        bytes32 digest = _hashTypedDataV4(keccak256(abi.encode(_DEREG_TYPEHASH, pkHash, msg.sender, nonce, deadline)));

        require(SignatureChecker.isValidSignatureNow(nodeAddress0, digest, signature), InvalidSignature());
        unchecked {
            nonces[nodeAddress0] = nonce + 1;
        }

        _deregisterValidator(valPubKeys[0]);
        emit ValidatorDeregistered(valPubKeys[0], nodeAddress0);

        uint256 len = valPubKeys.length;
        for (uint256 i = 1; i < len; ++i) {
            require(getNodeAddressFromPubkey(valPubKeys[i]) == nodeAddress0, MixedNodeBatch(valPubKeys[i]));
            _deregisterValidator(valPubKeys[i]);
            emit ValidatorDeregistered(valPubKeys[i], nodeAddress0);
        }
    }

    /// @dev Allows the freeze oracle account to freeze validators which disobey the mev-commit protocol.
    function freeze(bytes[] calldata valPubKeys) external whenNotPaused onlyFreezeOracle {
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _freeze(valPubKeys[i]);
        }
    }

    /// @dev Allows any account to unfreeze validators which have been frozen, for a fee.
    function unfreeze(bytes[] calldata valPubKeys) external payable whenNotPaused nonReentrant {
        uint256 requiredFee = unfreezeFee * valPubKeys.length;
        require(msg.value >= requiredFee, UnfreezeFeeRequired(requiredFee));
        uint256 len = valPubKeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _unfreeze(valPubKeys[i]);
        }
        (bool success,) = unfreezeReceiver.call{value: requiredFee}("");
        require(success, UnfreezeTransferFailed());
        uint256 excessFee = msg.value - requiredFee;
        if (excessFee != 0) {
            (bool successRefund,) = msg.sender.call{value: excessFee}("");
            require(successRefund, RefundFailed());
        }
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external onlyOwner {
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

    function isValidatorRegistered(bytes calldata validatorPubkey) public view returns (bool) {
        return validatorRegistrations[validatorPubkey].exists;
    }

    /// @dev Returns the node address of a minipool.
    function getNodeAddressFromMinipool(address minipool) public view returns (address) {
        address nodeAddress = RocketMinipoolInterface(minipool).getNodeAddress();
        require(nodeAddress != address(0), NoMinipoolForKey());
        return nodeAddress;
    }

    function pubkeysHash(bytes[] calldata pubkeys) public pure returns (bytes32) {
        uint256 count = pubkeys.length;

        // Contiguous buffer to hold all pubkeys back-to-back
        bytes memory concatenated = new bytes(48 * count);

        // Pointer to the next write position in `concatenated`
        uint256 writePtr;
        assembly {
            writePtr := add(concatenated, 32)
        }

        for (uint256 i = 0; i < count; ++i) {
            bytes calldata pubkey = pubkeys[i];
            assembly {
                let len := 48
                calldatacopy(writePtr, pubkey.offset, len)
                writePtr := add(writePtr, len)
            }
        }
        return keccak256(concatenated);
    }

    /// @dev Registers a validator.
    function _registerValidator(bytes calldata valPubKey) internal {
        require(!isValidatorRegistered(valPubKey), ValidatorAlreadyRegistered(valPubKey));
        validatorRegistrations[valPubKey].exists = true;
    }

    /// @dev Requests deregistration for a validator.
    function _requestValidatorDeregistration(bytes calldata valPubKey) internal {
        require(isValidatorRegistered(valPubKey), ValidatorNotRegistered(valPubKey));
        IRocketMinipoolRegistry.ValidatorRegistration storage reg = validatorRegistrations[valPubKey];
        require(reg.deregTimestamp == 0, DeregRequestAlreadyExists(valPubKey));
        reg.deregTimestamp = uint64(block.timestamp);
    }

    function _deregisterValidator(bytes calldata valPubKey) internal {
        IRocketMinipoolRegistry.ValidatorRegistration storage reg = validatorRegistrations[valPubKey];
        require(reg.freezeTimestamp == 0, FrozenValidatorCannotDeregister(valPubKey));
        require(reg.deregTimestamp != 0, DeregRequestDoesNotExist(valPubKey));
        require(uint64(block.timestamp) > reg.deregTimestamp + deregistrationPeriod, DeregistrationTooSoon(valPubKey));
        delete validatorRegistrations[valPubKey];
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
        require(newFreezeOracle != address(0), ZeroParam());
        freezeOracle = newFreezeOracle;
    }

    function _setUnfreezeReceiver(address newUnfreezeReceiver) internal {
        require(newUnfreezeReceiver != address(0), ZeroParam());
        unfreezeReceiver = newUnfreezeReceiver;
    }

    function _setUnfreezeFee(uint256 newUnfreezeFee) internal {
        require(newUnfreezeFee != 0, ZeroParam());
        unfreezeFee = newUnfreezeFee;
    }

    function _setRocketStorage(address newRocketStorage) internal {
        require(newRocketStorage != address(0), ZeroParam());
        rocketStorage = RocketStorageInterface(newRocketStorage);
    }

    function _setDeregistrationPeriod(uint64 newDeregistrationPeriod) internal {
        deregistrationPeriod = newDeregistrationPeriod;
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _getFirstNodeAddrAndValidate(bytes calldata firstPubKey)
        internal
        view
        returns (address firstNodeAddress)
    {
        firstNodeAddress = getNodeAddressFromMinipool(getMinipoolFromPubkey(firstPubKey));
        address withdrawalAddress = rocketStorage.getNodeWithdrawalAddress(firstNodeAddress);
        require(withdrawalAddress == msg.sender, NotWithdrawalAddress(withdrawalAddress));
    }

    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        if (!isValidatorRegistered(valPubKey)) return false;
        if (validatorRegistrations[valPubKey].freezeTimestamp != 0) return false;
        if (validatorRegistrations[valPubKey].deregTimestamp != 0) return false;
        address minipool = getMinipoolFromPubkey(valPubKey);
        if (minipool == address(0)) return false;
        if (!isMinipoolActive(minipool)) return false;
        return true;
    }

    /// @dev Deterministically hash a list of pubkeys to avoid large EIP-712 arrays.
    function _hashPubkeysAndFunction(bytes memory functionSelector, bytes[] calldata pubkeys)
        internal
        pure
        returns (bytes32)
    {
        bytes32[] memory leaves = new bytes32[](pubkeys.length);
        uint256 len = pubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            leaves[i] = keccak256(pubkeys[i]);
        }
        return keccak256(abi.encodePacked(functionSelector, leaves));
    }
}
