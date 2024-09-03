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
// TODO: You're prob able to remove some of the dereg logic for vaults etc. and piggyback off symbiotic core "vault epochs" etc. 
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
        uint256 _operatorDeregPeriodBlocks,
        uint256 _validatorDeregPeriodBlocks,
        uint256 _vaultDeregPeriodBlocks,
        address _slashOracle,
        address _owner
    ) public initializer {
        _setNetworkRegistry(_networkRegistry);
        _setOperatorRegistry(_operatorRegistry);
        _setVaultFactory(_vaultFactory);
        _setNetwork(_network);
        _setOperatorDeregPeriodBlocks(_operatorDeregPeriodBlocks);
        _setValidatorDeregPeriodBlocks(_validatorDeregPeriodBlocks);
        _setVaultDeregPeriodBlocks(_vaultDeregPeriodBlocks);
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
    // TODO: enforce that validator would be slashable (enough funds + high enough prio) to allow the registration.
    // Idea here is we need to enforce a newly registered validator is immediately opted-in.
    // TODO: require that provided vault is registered and has enough funds.
    function registerValidators(bytes[][] calldata blsPubkeys, address[] calldata vaults) external whenNotPaused {
        uint256 vaultLen = vaults.length;
        require(vaultLen == blsPubkeys.length, "invalid array lengths");
        for (uint256 i = 0; i < vaultLen; ++i) {
            uint256 keyLen = blsPubkeys[i].length;
            for (uint256 j = 0; j < keyLen; ++j) {
                _addValRecord(blsPubkeys[i][j], vaults[i]);
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
    // TODO: IMPORTANT, this FUNCTION SHOULD NOT BE ONLY OWNER, BUT ALSO OPERATOR.
    // TODO: OWNER can only delete if operator is blacklisted, prob make this separate function.
    function deregisterValidators(bytes[] calldata blsPubkeys) external onlyOwner {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _deregisterValidator(blsPubkeys[i]);
        }
    }

    function registerVaults(address[] calldata vaults, address[] calldata operators, uint256[] calldata slashAmounts) external onlyOwner {
        uint256 vLen = vaults.length;
        require(vLen == operators.length && vLen == slashAmounts.length, "invalid length");
        for (uint256 i = 0; i < vLen; i++) {
            _registerVault(vaults[i], operators[i], slashAmounts[i]);
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

    /// @dev Sets the operator deregistration period in blocks, restricted to contract owner.
    function setOperatorDeregPeriodBlocks(uint256 operatorDeregPeriodBlocks_) external onlyOwner {
        _setOperatorDeregPeriodBlocks(operatorDeregPeriodBlocks_);
    }

    /// @dev Sets the validator deregistration period in blocks, restricted to contract owner.
    function setValidatorDeregPeriodBlocks(uint256 validatorDeregPeriodBlocks_) external onlyOwner {
        _setValidatorDeregPeriodBlocks(validatorDeregPeriodBlocks_);
    }

    /// @dev Sets the vault deregistration period in blocks, restricted to contract owner.
    function setVaultDeregPeriodBlocks(uint256 vaultDeregPeriodBlocks_) external onlyOwner {
        _setVaultDeregPeriodBlocks(vaultDeregPeriodBlocks_);
    }

    /// @dev Sets the slash oracle, restricted to contract owner.
    function setSlashOracle(address slashOracle_) external onlyOwner {
        _setSlashOracle(slashOracle_);
    }

    function isValidatorOptedIn(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorOptedIn(blsPubkey);
    }

    function isValidatorSlashable(bytes calldata blsPubkey) external view returns (bool) {
        return _isValidatorSlashable(blsPubkey);
    }

    function vaultCollateralizesAllValidators(address vault) external view returns (bool) {
        return _vaultCollteralizesAllValidators(vault);
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
        require(operatorRecords[operator].exists, "operator dereg not requested");
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

    function _setValRecord(bytes calldata blsPubkey, address vault) internal {
        validatorRecords[blsPubkey] = ValidatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            vault: vault
        });
        _vaultToValidatorSet[vault].add(blsPubkey);
    }

    // TODO: Need to add more requires here and below s.t. we don't allow operators who
    // are deregistered or req deregistered, from adding val records.
    // DO a full sweep comparison of MevCommitAVS to see which checks exist for each function.
    function _addValRecord(bytes calldata blsPubkey, address vault) internal {
        require(!validatorRecords[blsPubkey].exists, "val record already exists");

        require(operatorRecords[msg.sender].exists, "operator not registered");
        require(!operatorRecords[msg.sender].deregRequestHeight.exists, "operator dereg req exists");
        require(!operatorRecords[msg.sender].isBlacklisted, "operator is blacklisted");

        require(vaultRecords[vault].exists, "vault not registered");
        require(!vaultRecords[vault].deregRequestHeight.exists, "vault dereg req exists");

        // TODO: check vault has enough funds

        require(vaultRecords[vault].operator == msg.sender, "vault operator mismatch");

        _setValRecord(blsPubkey, vault);
        emit ValRecordAdded(blsPubkey, msg.sender, _getPositionInValset(blsPubkey));
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        require(_getOperatorFromValRecord(blsPubkey) == msg.sender, "sender is not operator");
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, _getPositionInValset(blsPubkey));
    }

    function _deregisterValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing val record");
        address operator = _getOperatorFromValRecord(blsPubkey);
        require(operatorRecords[operator].exists, "operator not registered");
        require(!operatorRecords[operator].deregRequestHeight.exists, "operator dereg request exists");
        require(operatorRecords[operator].isBlacklisted, "operator is blacklisted");
        delete validatorRecords[blsPubkey];
        address vault = validatorRecords[blsPubkey].vault;
        _vaultToValidatorSet[vault].remove(blsPubkey);
        emit ValRecordDeleted(blsPubkey, operator);
    }

    function _setVaultRecord(address vault, address operator, uint256 slashAmount) internal {
        vaultRecords[vault] = VaultRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            operator: operator,
            slashAmount: slashAmount
        });
    }

    function _registerVault(address vault, address operator, uint256 slashAmount) internal {
        require(!vaultRecords[vault].exists, "vault already registered");
        require(vaultFactory.isEntity(vault), "vault not registered");
        require(operatorRegistry.isEntity(operator), "operator not registered");
        require(slashAmount != 0, "slash amount must be non-zero");

        // Check all slashable stake is allocated to single operator
        IVaultStorage vaultContract = IVaultStorage(vault);
        IBaseDelegator delegator = IBaseDelegator(vaultContract.delegator());
        uint256 stake = delegator.stake(_getSubnetwork(), operator);
        require(stake != 0, "operator must have vault stake");
        uint256 maxNetworkLimit = delegator.maxNetworkLimit(_getSubnetwork());
        require(stake == maxNetworkLimit, "oper stake != network limit");

        // TODO: Ensure vault epoch duration is long enough that oracle can submit slash tx in time. 
        // Maybe equiv to L1 epoch? 
        uint256 vaultEpochDuration = vaultContract.epochDuration();
        require(vaultEpochDuration > 33, "invalid vault epoch");

        _setVaultRecord(vault, operator, slashAmount);
        emit VaultRegistered(vault, operator, slashAmount);
    }

    function _updateSlashAmount(address vault, uint256 slashAmount) internal {
        require(vaultRecords[vault].exists, "vault not registered");
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

    /// @dev Internal function to set the vault deregistration period in blocks.
    function _setVaultDeregPeriodBlocks(uint256 vaultDeregPeriodBlocks_) internal {
        vaultDeregPeriodBlocks = vaultDeregPeriodBlocks_;
        emit VaultDeregPeriodBlocksSet(vaultDeregPeriodBlocks_);
    }

    /// @dev Internal function to set the slash oracle.
    function _setSlashOracle(address slashOracle_) internal {
        slashOracle = slashOracle_;
        emit SlashOracleSet(slashOracle_);
    }

    /// @dev Authorizes contract upgrades, restricted to contract owner.
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // TODO: confirm you can call this with zero valued stuff
    function _getOperatorFromValRecord(bytes calldata blsPubkey) internal view returns (address) {
        return vaultRecords[validatorRecords[blsPubkey].vault].operator;
    }

    // TODO: This and above function are optimistic on existence of records.
    // TODO: Either a. confirm the calling functions have neccessary requires, or add them here.
    function _getPositionInValset(bytes calldata blsPubkey) internal view returns (uint256) {
        return _vaultToValidatorSet[validatorRecords[blsPubkey].vault].position(blsPubkey);
    }

    function _isValidatorReadyToDeregister(bytes calldata blsPubkey) internal view returns (bool) {
        return validatorRecords[blsPubkey].deregRequestHeight.exists && 
            block.number > validatorDeregPeriodBlocks + validatorRecords[blsPubkey].deregRequestHeight.blockHeight;
    }

    function _isOperatorReadyToDeregister(address operator) internal view returns (bool) {
        return operatorRecords[operator].deregRequestHeight.exists && 
            block.number > operatorDeregPeriodBlocks + operatorRecords[operator].deregRequestHeight.blockHeight;
    }

    function _isVaultReadyToDeregister(address vault) internal view returns (bool) {
        return vaultRecords[vault].deregRequestHeight.exists && 
            block.number > vaultDeregPeriodBlocks + vaultRecords[vault].deregRequestHeight.blockHeight;
    }

    function _getSubnetwork() internal view returns (bytes32) {
        return Subnetwork.subnetwork(network, SUBNETWORK_ID);
    }

    function _getAllocatedStake(address vault) internal view returns (uint256) {
        IBaseDelegator delegator = IBaseDelegator(IVault(vault).delegator());
        bytes32 subnetwork = _getSubnetwork();
        address operator = vaultRecords[vault].operator;
        return delegator.stake(subnetwork, operator);
    }

    function _vaultCollteralizesAllValidators(address vault) internal view returns (bool) {
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 numVals = _vaultToValidatorSet[vault].length();
        uint256 allocatedStake = _getAllocatedStake(vault);
        return allocatedStake > slashAmount * numVals;
    }

    function _isValidatorSlashable(bytes calldata blsPubkey) internal view returns (bool) {
        address vault = validatorRecords[blsPubkey].vault;
        uint256 allocatedStake = _getAllocatedStake(vault);
        uint256 slashAmount = vaultRecords[vault].slashAmount;
        uint256 position = _getPositionInValset(blsPubkey);
        uint256 slashableVals = allocatedStake / slashAmount;
        return position < slashableVals;
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
        // TODO: Check symbiotic core vault registry?
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
        if (!_isValidatorSlashable(blsPubkey)) {
            return false;
        }
        return true;
    }
}
