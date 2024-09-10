// SPDX-License-Identifier: BSL 1.1

pragma solidity 0.8.26;

import {EventHeightLib} from "../../utils/EventHeight.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareStorage} from "./MevCommitMiddlewareStorage.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";
import {IVault} from "symbiotic-core/interfaces/vault/IVault.sol";
import {IVaultStorage} from "symbiotic-core/interfaces/vault/IVaultStorage.sol";
import {IBaseDelegator} from "symbiotic-core/interfaces/delegator/IBaseDelegator.sol";
import {IEntity} from "symbiotic-core/interfaces/common/IEntity.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Subnetwork} from "symbiotic-core/contracts/libraries/Subnetwork.sol";
import {ISlasher} from "symbiotic-core/interfaces/slasher/ISlasher.sol";
import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";
import {Errors} from "../../utils/Errors.sol";

contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    using EnumerableSet for EnumerableSet.BytesSet;

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, OnlySlashOracle(slashOracle));
        _;
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        IRegistry _networkRegistry,
        IRegistry _operatorRegistry,
        IRegistry _vaultFactory,
        address _network,
        uint256 _slashPeriodBlocks,
        address _slashOracle,
        address _owner
    ) public initializer {
        _setNetworkRegistry(_networkRegistry);
        _setOperatorRegistry(_operatorRegistry);
        _setVaultFactory(_vaultFactory);
        _setNetwork(_network);
        _setSlashPeriodBlocks(_slashPeriodBlocks);
        _setSlashOracle(_slashOracle);
        __Pausable_init();
        __UUPSUpgradeable_init();
        __Ownable_init(_owner);
    }

    /// @dev Receive function to prevent unintended contract interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to prevent unintended contract interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    function registerOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _registerOperator(operators[i]);
        }
    }

    function requestOperatorDeregistrations(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestOperatorDeregistration(operators[i]);
        }
    }

    function deregisterOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterOperator(operators[i]);
        }
    }

    function blacklistOperators(address[] calldata operators) external onlyOwner {
        uint256 len = operators.length;
        for (uint256 i = 0; i < len; ++i) {
            _blacklistOperator(operators[i]);
        }
    }

    function registerValidators(bytes[][] calldata blsPubkeys, address[] calldata vaults) external whenNotPaused {
        uint256 vaultLen = vaults.length;
        require(vaultLen == blsPubkeys.length, InvalidArrayLengths(vaultLen, blsPubkeys.length));
        address operator = msg.sender;
        _checkOperator(operator);
        for (uint256 i = 0; i < vaultLen; ++i) {
            uint256 keyLen = blsPubkeys[i].length;
            _checkVault(vaults[i]);
            uint256 potentialSlashableVals = _potentialSlashableVals(vaults[i], operator);
            require(keyLen < potentialSlashableVals + 1,
                ValidatorsNotSlashable(vaults[i], operator, keyLen, potentialSlashableVals));
            for (uint256 j = 0; j < keyLen; ++j) {
                _addValRecord(blsPubkeys[i][j], vaults[i], operator);
            }
        }
    }

    function requestValDeregistrations(bytes[] calldata blsPubkeys) external whenNotPaused {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestValDeregistration(blsPubkeys[i]);
        }
    }

    /// @dev Deletes validator records, only if the associated operator is blacklisted.
    /// Restricted to contract owner.
    /// @notice This function allows the contract owner to combat a greifing scenario where an operator
    /// registers a validator pubkey that it does not control, own, or otherwise manage.
    function deregisterValidators(bytes[] calldata blsPubkeys) external {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterValidator(blsPubkeys[i]);
        }
    }

    function registerVaults(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, InvalidArrayLengths(vLen, slashAmounts.length));
        for (uint256 i = 0; i < vLen; ++i) {
            _registerVault(vaults[i], slashAmounts[i]);
        }
    }

    function updateSlashAmounts(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, InvalidArrayLengths(vLen, slashAmounts.length));
        for (uint256 i = 0; i < vLen; ++i) {
            _updateSlashAmount(vaults[i], slashAmounts[i]);
        }
    }

    function requestVaultDeregistrations(address[] calldata vaults) external onlyOwner {
        uint256 len = vaults.length;
        for (uint256 i = 0; i < len; ++i) {
            _requestVaultDeregistration(vaults[i]);
        }
    }

    function deregisterVaults(address[] calldata vaults) external onlyOwner {
        uint256 len = vaults.length;
        for (uint256 i = 0; i < len; ++i) {
            _deregisterVault(vaults[i]);
        }
    }

    /// @dev Slashes validators and marks them for deregistration.
    /// @param blsPubkeys The L1 validator BLS public keys to slash.
    /// @param infractionTimestamps The block.timestamps for blocks during which each infraction occurred.
    function slashValidators(bytes[] calldata blsPubkeys, uint256[] calldata infractionTimestamps) external onlySlashOracle {
        uint256 len = blsPubkeys.length;
        for (uint256 i = 0; i < len; ++i) {
            _slashValidator(blsPubkeys[i], infractionTimestamps[i]);
        }
    }

    /// @dev Pauses the contract, restricted to contract owner.
    function pause() external onlyOwner { _pause(); }

    /// @dev Unpauses the contract, restricted to contract owner.
    function unpause() external onlyOwner { _unpause(); }

    /// @dev Sets the network registry, restricted to contract owner.
    function setNetworkRegistry(IRegistry _networkRegistry) external onlyOwner {
        _setNetworkRegistry(_networkRegistry);
    }

    /// @dev Sets the operator registry, restricted to contract owner.
    function setOperatorRegistry(IRegistry _operatorRegistry) external onlyOwner {
        _setOperatorRegistry(_operatorRegistry);
    }

    /// @dev Sets the vault factory, restricted to contract owner.
    function setVaultFactory(IRegistry _vaultFactory) external onlyOwner {
        _setVaultFactory(_vaultFactory);
    }

    /// @dev Sets the network address, restricted to contract owner.
    function setNetwork(address _network) external onlyOwner {
        _setNetwork(_network);
    }

    /// @dev Sets the slash period in blocks, restricted to contract owner.
    function setSlashPeriodBlocks(uint256 slashPeriodBlocks_) external onlyOwner {
        _setSlashPeriodBlocks(slashPeriodBlocks_);
    }

    /// @dev Sets the slash oracle, restricted to contract owner.
    function setSlashOracle(address slashOracle_) external onlyOwner {
        _setSlashOracle(slashOracle_);
    }

    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    function isValidatorSlashable(bytes calldata blsPubkey) external view returns (bool) {
        require(validatorRecords[blsPubkey].exists, MissingValRecord(blsPubkey));
        _checkVault(validatorRecords[blsPubkey].vault);
        _checkOperator(validatorRecords[blsPubkey].operator);
        return _isValidatorSlashable(blsPubkey,
            validatorRecords[blsPubkey].vault, validatorRecords[blsPubkey].operator);
    }

    function potentialSlashableValidators(address vault, address operator) external view returns (uint256) {
        return _potentialSlashableVals(vault, operator);
    }

    function allValidatorsAreSlashable(address vault, address operator) external view returns (bool) {
        return _allValidatorsAreSlashable(vault, operator);
    }

    function _setOperatorRecord(address operator) internal {
        operatorRecords[operator] = OperatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            isBlacklisted: false
        });
    }

    function _registerOperator(address operator) internal {
        require(!operatorRecords[operator].exists, OperatorAlreadyRegistered(operator));
        require(operatorRegistry.isEntity(operator), OperatorNotEntity(operator));
        _setOperatorRecord(operator);
        emit OperatorRegistered(operator);
    }

    function _requestOperatorDeregistration(address operator) internal {
        require(operatorRecords[operator].exists, OperatorNotRegistered(operator));
        require(!operatorRecords[operator].isBlacklisted, OperatorIsBlacklisted(operator));
        require(!operatorRecords[operator].deregRequestHeight.exists, OperatorDeregRequestExists(operator));
        EventHeightLib.set(operatorRecords[operator].deregRequestHeight, block.number);
        emit OperatorDeregistrationRequested(operator);
    }

    function _deregisterOperator(address operator) internal {
        require(operatorRecords[operator].exists, OperatorNotRegistered(operator));
        require(_isOperatorReadyToDeregister(operator), OperatorNotReadyToDeregister(
            operator, block.number, operatorRecords[operator].deregRequestHeight.blockHeight));
        require(!operatorRecords[operator].isBlacklisted, OperatorIsBlacklisted(operator));
        delete operatorRecords[operator];
        emit OperatorDeregistered(operator);
    }

    function _blacklistOperator(address operator) internal {
        if (!operatorRecords[operator].exists) {
            _setOperatorRecord(operator);
        }
        require(!operatorRecords[operator].isBlacklisted, OperatorAlreadyBlacklisted(operator));
        operatorRecords[operator].isBlacklisted = true;
        emit OperatorBlacklisted(operator);
    }

    function _setValRecord(bytes calldata blsPubkey, address vault, address operator) internal {
        validatorRecords[blsPubkey] = ValidatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            vault: vault,
            operator: operator
        });
        _vaultAndOperatorToValset[vault][operator].add(blsPubkey);
    }

    function _addValRecord(bytes calldata blsPubkey, address vault, address operator) internal {
        require(!validatorRecords[blsPubkey].exists, ValidatorRecordAlreadyExists(blsPubkey));
        _setValRecord(blsPubkey, vault, operator);
        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        emit ValRecordAdded(blsPubkey, msg.sender, position);
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, MissingValidatorRecord(blsPubkey));
        if (msg.sender != owner()) {
            _checkCallingOperator(validatorRecords[blsPubkey].operator);
        }
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);
        uint256 position = _getPositionInValset(blsPubkey, validatorRecords[blsPubkey].vault,
            validatorRecords[blsPubkey].operator);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, position);
    }

    function _deregisterValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, MissingValidatorRecord(blsPubkey));
        require(_isValidatorReadyToDeregister(blsPubkey), ValidatorNotReadyToDeregister(
            blsPubkey, block.number, validatorRecords[blsPubkey].deregRequestHeight.blockHeight));
        if (msg.sender != owner()) {
            _checkCallingOperator(validatorRecords[blsPubkey].operator);
        }
        address vault = validatorRecords[blsPubkey].vault;
        address operator = validatorRecords[blsPubkey].operator;
        _vaultAndOperatorToValset[vault][operator].remove(blsPubkey);
        delete validatorRecords[blsPubkey];
        emit ValRecordDeleted(blsPubkey, msg.sender);
    }

    function _setVaultRecord(address vault, uint256 slashAmount) internal {
        vaultRecords[vault] = VaultRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            slashAmount: slashAmount
        });
    }

    function _registerVault(address vault,uint256 slashAmount) internal {
        require(!vaultRecords[vault].exists, VaultAlreadyRegistered(vault));
        require(vaultFactory.isEntity(vault), VaultNotEntity(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));

        IVaultStorage vaultContract = IVaultStorage(vault);
        uint256 vaultEpochDuration = vaultContract.epochDuration();
        require(vaultEpochDuration > slashPeriodBlocks,
            InvalidVaultEpochDuration(vault, vaultEpochDuration, slashPeriodBlocks));
        
        IEntity delegator = IEntity(IVault(vault).delegator());
        if (delegator.TYPE() == FULL_RESTAKE_DELEGATOR_TYPE) {
            revert FullRestakeDelegatorNotSupported(vault);
        } else if (delegator.TYPE() != NETWORK_RESTAKE_DELEGATOR_TYPE) {
            revert UnknownDelegatorType(vault, delegator.TYPE());
        }

        address slasher = IVault(vault).slasher();
        require(slasher != address(0), SlasherNotSetForVault(vault));
        uint256 slasherType = IEntity(slasher).TYPE();
        if (slasherType == VETO_SLASHER_TYPE) {
            revert VetoSlasherNotSupported(vault);
        } else if (slasherType != INSTANT_SLASHER_TYPE) {
            revert UnknownSlasherType(vault, slasherType);
        }

        _setVaultRecord(vault, slashAmount);
        emit VaultRegistered(vault, slashAmount);
    }

    function _updateSlashAmount(address vault, uint256 slashAmount) internal {
        require(vaultRecords[vault].exists, VaultNotRegistered(vault));
        require(slashAmount != 0, SlashAmountMustBeNonZero(vault));
        vaultRecords[vault].slashAmount = slashAmount;
        emit VaultSlashAmountUpdated(vault, slashAmount);
    }

    function _requestVaultDeregistration(address vault) internal {
        require(vaultRecords[vault].exists, VaultNotRegistered(vault));
        require(!vaultRecords[vault].deregRequestHeight.exists, VaultDeregRequestExists(vault));
        EventHeightLib.set(vaultRecords[vault].deregRequestHeight, block.number);
        emit VaultDeregistrationRequested(vault);
    }

    function _deregisterVault(address vault) internal {
        require(vaultRecords[vault].exists, VaultNotRegistered(vault));
        require(_isVaultReadyToDeregister(vault), VaultNotReadyToDeregister(vault, block.number,
            vaultRecords[vault].deregRequestHeight.blockHeight));
        delete vaultRecords[vault];
        emit VaultDeregistered(vault);
    }

    /// @dev Slashes a validator and marks it for deregistration.
    /// @param blsPubkey The L1 validator BLS public key to slash.
    /// @param infractionTimestamp The block.timestamp for the block during which the infraction occurred.
    function _slashValidator(bytes calldata blsPubkey, uint256 infractionTimestamp) internal {
        require(validatorRecords[blsPubkey].exists, MissingValidatorRecord(blsPubkey));
        address vault = validatorRecords[blsPubkey].vault;
        require(vaultRecords[vault].exists, MissingVaultRecord(vault));
        address operator = validatorRecords[blsPubkey].operator;
        uint256 amount = vaultRecords[vault].slashAmount;

        ISlasher(IVault(vault).slasher()).slash(
            _getSubnetwork(), operator, amount, SafeCast.toUint48(infractionTimestamp), new bytes(0));

        // Set dereg request height so validator is no longer opted-in.
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);

        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        emit ValidatorSlashed(blsPubkey, operator, position);
    }

    /// @dev Internal function to set the network registry.
    function _setNetworkRegistry(IRegistry _networkRegistry) internal {
        require(_networkRegistry != IRegistry(address(0)), ZeroAddressNotAllowed());
        networkRegistry = _networkRegistry;
        emit NetworkRegistrySet(address(_networkRegistry));
    }

    /// @dev Internal function to set the operator registry.
    function _setOperatorRegistry(IRegistry _operatorRegistry) internal {
        require(_operatorRegistry != IRegistry(address(0)), ZeroAddressNotAllowed());
        operatorRegistry = _operatorRegistry;
        emit OperatorRegistrySet(address(_operatorRegistry));
    }

    function _setVaultFactory(IRegistry _vaultFactory) internal {
        require(_vaultFactory != IRegistry(address(0)), ZeroAddressNotAllowed());
        vaultFactory = _vaultFactory;
        emit VaultFactorySet(address(_vaultFactory));
    }

    /// @dev Internal function to set the network address, which must have registered with the NETWORK_REGISTRY.
    function _setNetwork(address _network) internal {
        require(networkRegistry.isEntity(_network), NetworkNotEntity(_network));
        network = _network;
        emit NetworkSet(_network);
    }

    /// @dev Internal function to set the slash period in blocks.
    function _setSlashPeriodBlocks(uint256 slashPeriodBlocks_) internal {
        slashPeriodBlocks = slashPeriodBlocks_;
        emit SlashPeriodBlocksSet(slashPeriodBlocks_);
    }

    /// @dev Internal function to set the slash oracle.
    function _setSlashOracle(address slashOracle_) internal {
        slashOracle = slashOracle_;
        emit SlashOracleSet(slashOracle_);
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    function _checkOperator(address operator) internal view {
        require(operatorRegistry.isEntity(operator), OperatorNotEntity(operator));
        require(operatorRecords[operator].exists, OperatorNotRegistered(operator));
        require(!operatorRecords[operator].deregRequestHeight.exists, OperatorDeregRequestExists(operator));
        require(!operatorRecords[operator].isBlacklisted, OperatorIsBlacklisted(operator));
    }

    function _checkCallingOperator(address operator) internal view {
        require(msg.sender == operator, OnlyOperator(operator));
        _checkOperator(operator);
    }

    function _checkVault(address vault) internal view {
        require(vaultFactory.isEntity(vault), VaultNotEntity(vault));
        require(vaultRecords[vault].exists, VaultNotRegistered(vault));
        require(!vaultRecords[vault].deregRequestHeight.exists, VaultDeregRequestExists(vault));
    }

    function _getPositionInValset(bytes calldata blsPubkey,
        address vault, address operator) internal view returns (uint256) {
        return _vaultAndOperatorToValset[vault][operator].position(blsPubkey);
    }

    function _isValidatorReadyToDeregister(bytes calldata blsPubkey) internal view returns (bool) {
        return validatorRecords[blsPubkey].deregRequestHeight.exists && 
            block.number > slashPeriodBlocks + validatorRecords[blsPubkey].deregRequestHeight.blockHeight;
    }

    function _isOperatorReadyToDeregister(address operator) internal view returns (bool) {
        return operatorRecords[operator].deregRequestHeight.exists && 
            block.number > slashPeriodBlocks + operatorRecords[operator].deregRequestHeight.blockHeight;
    }

    function _isVaultReadyToDeregister(address vault) internal view returns (bool) {
        return vaultRecords[vault].deregRequestHeight.exists && 
            block.number > slashPeriodBlocks + vaultRecords[vault].deregRequestHeight.blockHeight;
    }

    function _getSubnetwork() internal view returns (bytes32) {
        return Subnetwork.subnetwork(network, SUBNETWORK_ID);
    }

    function _getAllocatedStake(address vault, address operator) internal view returns (uint256) {
        IBaseDelegator delegator = IBaseDelegator(IVault(vault).delegator());
        bytes32 subnetwork = _getSubnetwork();
        return delegator.stake(subnetwork, operator);
    }

    function _allValidatorsAreSlashable(address vault, address operator) internal view returns (bool) {
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 numVals = _vaultAndOperatorToValset[vault][operator].length();
        uint256 allocatedStake = _getAllocatedStake(vault, operator);
        return allocatedStake > slashAmount * numVals;
    }

    function _isValidatorSlashable(bytes calldata blsPubkey, address vault, address operator) internal view returns (bool) {
        uint256 allocatedStake = _getAllocatedStake(vault, operator);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 position = _getPositionInValset(blsPubkey, vault, operator);
        uint256 slashableVals = allocatedStake / slashAmount;
        return position < slashableVals;
    }

    function _potentialSlashableVals(address vault, address operator) internal view returns (uint256) {
        uint256 allocatedStake = _getAllocatedStake(vault, operator);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 alreadyCollateralized = _vaultAndOperatorToValset[vault][operator].length();
        uint256 slashableVals = allocatedStake / slashAmount;
        return slashableVals - alreadyCollateralized;
    }
    
    function _isValidatorOptedIn(bytes calldata blsPubkey) internal view returns (bool) {
        if (!validatorRecords[blsPubkey].exists) {
            return false;
        }
        if (validatorRecords[blsPubkey].deregRequestHeight.exists) {
            return false;
        }
        if (!vaultRecords[validatorRecords[blsPubkey].vault].exists) {
            return false;
        }
        if (vaultRecords[validatorRecords[blsPubkey].vault].deregRequestHeight.exists) {
            return false;
        }
        if (!vaultFactory.isEntity(validatorRecords[blsPubkey].vault)) {
            return false;
        }
        address operator = validatorRecords[blsPubkey].operator;
        if (!operatorRecords[operator].exists) {
            return false;
        }
        if (operatorRecords[operator].deregRequestHeight.exists) {
            return false;
        }
        if (operatorRecords[operator].isBlacklisted) {
            return false;
        }
        if (!operatorRegistry.isEntity(operator)) {
            return false;
        }
        if (!_isValidatorSlashable(blsPubkey, validatorRecords[blsPubkey].vault,
            validatorRecords[blsPubkey].operator)) {
            return false;
        }
        return true;
    }
}
