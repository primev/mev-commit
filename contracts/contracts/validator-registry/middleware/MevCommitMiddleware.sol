// SPDX-License-Identifier: BSL 1.1

pragma solidity 0.8.26;

import {EventHeightLib} from "../../utils/EventHeight.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareStorage} from "./MevCommitMiddlewareStorage.sol";
import {EnumerableSet} from "../../utils/EnumerableSet.sol";
import {IVault} from "symbiotic-core/interfaces/Vault/IVault.sol";
import {IVaultStorage} from "symbiotic-core/interfaces/Vault/IVaultStorage.sol";
import {IBaseDelegator} from "symbiotic-core/interfaces/Delegator/IBaseDelegator.sol";
import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {Subnetwork} from "symbiotic-core/contracts/libraries/Subnetwork.sol";

// TODO: add symbiotic core integration via lifecycle: https://docs.symbiotic.fi/core-modules/networks#staking-lifecycle
// TODO: determine if you need timestamping similar to cosmos sdk example. Edit yes you will for slashing. See "captureTimestamp". 
// TODO: Parse through MevCommitAVS and make sure translatable reg/dreg functions have the same operators / check the same things. 
// TODO: for example you need to add requires s.t. a validator MUST be opted-in right after registering. 
// TODO: add function for a validator to "chage vault used for collateral", which involves a delete + new reg. 
// TODO: attempt to make storage more fsm like with enum. See if this can lessen the amount of requires needed
// TODO: Get through full Handbook for Networks page and confirm you follow all rules for slashing logic, network epoch, slashing epochs etc. 
// TODO: Use custom errors since our clients are compatible with this now.
// TODO: Accept BOTH vaults that have slashing via resolver or not. Oracle account can be resolver.
contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {

    using EnumerableSet for EnumerableSet.BytesSet;

    uint96 public constant SUBNETWORK_ID = 1;

    // TODO: more modifiers similar to MevCommitAVS

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "only slash oracle");
        _;
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

    constructor() {
        _disableInitializers();
    }

    function registerOperators(address[] calldata operators) external onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _registerOperator(operators[i]);
        }
    }

    function requestOperatorDeregistrations(address[] calldata operators) external onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _requestOperatorDeregistration(operators[i]);
        }
    }

    function deregisterOperators(address[] calldata operators) external onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _deregisterOperator(operators[i]);
        }
    }

    function blacklistOperators(address[] calldata operators) external onlyOwner {
        for (uint256 i = 0; i < operators.length; i++) {
            _blacklistOperator(operators[i]);
        }
    }

    // TODO: confirm this and other external functions can handle empty arrays
    // TODO: confirm only operator can edit their own records. Does contract owner need access as well?
    // Be consistent with MevCommitAVS.
    function registerValidators(bytes[][] calldata blsPubkeys, address[] calldata vaults) external whenNotPaused {
        uint256 vaultLen = vaults.length;
        require(vaultLen == blsPubkeys.length, "invalid array lengths");
        address operator = msg.sender;
        _checkOperator(operator);
        for (uint256 i = 0; i < vaultLen; ++i) {
            uint256 keyLen = blsPubkeys[i].length;
            _checkVault(vaults[i]);
            require(keyLen < _potentialSlashableVals(vaults[i], operator) + 1, "validators not slashable");
            for (uint256 j = 0; j < keyLen; ++j) {
                _addValRecord(blsPubkeys[i][j], vaults[i], operator);
            }
        }
    }

    function requestValDeregistrations(bytes[] calldata blsPubkeys) external whenNotPaused {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _requestValDeregistration(blsPubkeys[i]);
        }
    }

    /// @dev Deletes validator records, only if the associated operator is blacklisted.
    /// Restricted to contract owner.
    /// @notice This function allows the contract owner to combat a greifing scenario where an operator
    /// registers a validator pubkey that it does not control, own, or otherwise manage.
    ///
    /// TODO: Write test for scenario where operator greifs another, and contract owner
    /// has to blacklist that operator, then delete the greifed validator records.
    function deregisterValidators(bytes[] calldata blsPubkeys) external {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _deregisterValidator(blsPubkeys[i]);
        }
    }

    function registerVaults(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, "invalid length");
        for (uint256 i = 0; i < vLen; i++) {
            _registerVault(vaults[i], slashAmounts[i]);
        }
    }

    function updateSlashAmounts(address[] calldata vaults, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == slashAmounts.length, "invalid length");
        for (uint256 i = 0; i < vLen; i++) {
            _updateSlashAmount(vaults[i], slashAmounts[i]);
        }
    }

    function requestVaultDeregistrations(address[] calldata vaults) external onlyOwner {
        for (uint256 i = 0; i < vaults.length; i++) {
            _requestVaultDeregistration(vaults[i]);
        }
    }

    function deregisterVaults(address[] calldata vaults) external onlyOwner {
        for (uint256 i = 0; i < vaults.length; i++) {
            _deregisterVault(vaults[i]);
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
        require(validatorRecords[blsPubkey].exists, "missing val record");
        _checkVault(validatorRecords[blsPubkey].vault);
        _checkOperator(validatorRecords[blsPubkey].operator);
        return _isValidatorSlashable(blsPubkey,
            validatorRecords[blsPubkey].vault, validatorRecords[blsPubkey].operator);
    }

    function potentialSlashableValidators(address vault, address operator) external view returns (uint256) {
        return _potentialSlashableVals(vault, operator);
    }

    // TODO: Use this for unit tests
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
        require(!operatorRecords[operator].exists, "operator already registered");
        require(operatorRegistry.isEntity(operator), "operator not reg with core");
        _setOperatorRecord(operator);
        emit OperatorRegistered(operator);
    }

    function _requestOperatorDeregistration(address operator) internal {
        require(operatorRecords[operator].exists, "operator not registered");
        require(!operatorRecords[operator].isBlacklisted, "operator is blacklisted");
        EventHeightLib.set(operatorRecords[operator].deregRequestHeight, block.number);
        emit OperatorDeregistrationRequested(operator);
    }

    function _deregisterOperator(address operator) internal {
        require(operatorRecords[operator].exists, "operator not registered");
        require(_isOperatorReadyToDeregister(operator), "not ready to dereg");
        require(!operatorRecords[operator].isBlacklisted, "operator is blacklisted");
        delete operatorRecords[operator];
        emit OperatorDeregistered(operator);
    }

    // TODO: confirm validator can ALWAYS be blacklisted from any previous state,
    // and that no other operations can be performed on the operator record after being blacklisted.
    function _blacklistOperator(address operator) internal {
        if (!operatorRecords[operator].exists) {
            _setOperatorRecord(operator);
        }
        require(!operatorRecords[operator].isBlacklisted, "operator already blacklisted");
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

    // TODO: DO a full sweep comparison of MevCommitAVS to see which checks exist for each function.
    function _addValRecord(bytes calldata blsPubkey, address vault, address operator) internal {
        require(!validatorRecords[blsPubkey].exists, "val record already exists");
        _setValRecord(blsPubkey, vault, operator);
        emit ValRecordAdded(blsPubkey, msg.sender, _getPositionInValset(blsPubkey));
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing val record");
        if (msg.sender != owner()) {
            _checkCallingOperator(_getOperatorFromValRecord(blsPubkey));
        }
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, _getPositionInValset(blsPubkey));
    }

    function _deregisterValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing val record");
        require(_isValidatorReadyToDeregister(blsPubkey), "not ready to dereg");
        if (msg.sender != owner()) {
            _checkCallingOperator(_getOperatorFromValRecord(blsPubkey));
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
        require(!vaultRecords[vault].exists, "vault already registered");
        require(vaultFactory.isEntity(vault), "vault not entity");
        require(slashAmount != 0, "slash amount must be non-zero");

        IVaultStorage vaultContract = IVaultStorage(vault);
        uint256 vaultEpochDuration = vaultContract.epochDuration();
        require(vaultEpochDuration > slashPeriodBlocks, "invalid vault epoch duration");

        _setVaultRecord(vault, slashAmount);
        emit VaultRegistered(vault, slashAmount);
    }

    function _updateSlashAmount(address vault, uint256 slashAmount) internal {
        require(vaultRecords[vault].exists, "vault not registered");
        require(slashAmount != 0, "slash amount must be non-zero");
        vaultRecords[vault].slashAmount = slashAmount;
        emit VaultSlashAmountUpdated(vault, slashAmount);
    }

    function _requestVaultDeregistration(address vault) internal {
        require(vaultRecords[vault].exists, "vault not registered");
        require(!vaultRecords[vault].deregRequestHeight.exists, "vault dereg request already made");
        EventHeightLib.set(vaultRecords[vault].deregRequestHeight, block.number);
        emit VaultDeregistrationRequested(vault);
    }

    function _deregisterVault(address vault) internal {
        require(vaultRecords[vault].exists, "vault dereg not requested");
        require(_isVaultReadyToDeregister(vault), "dereg too soon");
        delete vaultRecords[vault];
        emit VaultDeregistered(vault);
    }

    // TODO: Feedback from meeting: Look into using historical state for slashing. 
    function _slashValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        address operator = _getOperatorFromValRecord(blsPubkey);
        // address vault = validatorRecords[blsPubkey].vault;
        // address slasher = IVault(vault).slasher();
        // uint256 slasherType = IEntity(slasher).TYPE();
        // if (slasherType == INSTANT_SLASHER_TYPE) {
        //     ISlasher(slasher).slash(subnetwork, operator, amount, timestamp, new bytes(0));
        // } else if (slasherType == VETO_SLASHER_TYPE) {
        //     IVetoSlasher(slasher).requestSlash(subnetwork, operator, amount, timestamp, new bytes(0));
        // } else {
        //     revert UnknownSlasherType();
        // }
        _requestValDeregistration(blsPubkey); // TODO: determine if validator should be deregistered
        emit ValidatorSlashed(blsPubkey, operator, _getPositionInValset(blsPubkey));
    }

    /// @dev Internal function to set the network registry.
    function _setNetworkRegistry(IRegistry _networkRegistry) internal {
        require(_networkRegistry != IRegistry(address(0)), "zero address not allowed");
        networkRegistry = _networkRegistry;
        emit NetworkRegistrySet(address(_networkRegistry));
    }

    /// @dev Internal function to set the operator registry.
    function _setOperatorRegistry(IRegistry _operatorRegistry) internal {
        require(_operatorRegistry != IRegistry(address(0)), "zero address not allowed");
        operatorRegistry = _operatorRegistry;
        emit OperatorRegistrySet(address(_operatorRegistry));
    }

    function _setVaultFactory(IRegistry _vaultFactory) internal {
        require(_vaultFactory != IRegistry(address(0)), "zero address not allowed");
        vaultFactory = _vaultFactory;
        emit VaultFactorySet(address(_vaultFactory));
    }

    /// @dev Internal function to set the network address, which must have registered with the NETWORK_REGISTRY.
    function _setNetwork(address _network) internal {
        require(networkRegistry.isEntity(_network), "network not registered");
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
        require(operatorRegistry.isEntity(operator), "operator not registered");
        require(operatorRecords[operator].exists, "operator not registered");
        require(!operatorRecords[operator].deregRequestHeight.exists, "operator dereg request exists");
        require(!operatorRecords[operator].isBlacklisted, "operator is blacklisted");
    }

    function _checkCallingOperator(address operator) internal view {
        require(msg.sender == operator, "only operator");
        _checkOperator(operator);
    }

    function _checkVault(address vault) internal view {
        require(vaultFactory.isEntity(vault), "vault not registered");
        require(vaultRecords[vault].exists, "vault not registered");
        require(!vaultRecords[vault].deregRequestHeight.exists, "vault dereg request exists");
    }

    function _getOperatorFromValRecord(bytes calldata blsPubkey) internal view returns (address) {
        return validatorRecords[blsPubkey].operator;
    }

    function _getPositionInValset(bytes calldata blsPubkey) internal view returns (uint256) {
        return _vaultAndOperatorToValset[
            validatorRecords[blsPubkey].vault][
                validatorRecords[blsPubkey].operator].position(blsPubkey);
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
        address vault = validatorRecords[blsPubkey].vault;
        address operator = validatorRecords[blsPubkey].operator;
        uint256 allocatedStake = _getAllocatedStake(vault, operator);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 position = _getPositionInValset(blsPubkey);
        uint256 slashableVals = allocatedStake / slashAmount;
        return position < slashableVals;
    }

    // TODO: def need unit tests, in conjuction with related functions
    function _potentialSlashableVals(address vault, address operator) internal view returns (uint256) {
        uint256 allocatedStake = _getAllocatedStake(vault, operator);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 alreadyCollateralized = _vaultAndOperatorToValset[vault][operator].length();
        uint256 slashableVals = allocatedStake / slashAmount;
        return slashableVals - alreadyCollateralized;
    }
    
    // TODO: Unit tests around confirming a validator who's newly registered MUST be opted in.
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
        address operator = _getOperatorFromValRecord(blsPubkey);
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
