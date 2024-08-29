// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {EventHeightLib} from "../../utils/EventHeight.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IMevCommitMiddleware} from "../../interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareStorage} from "./MevCommitMiddlewareStorage.sol";

// TODO: add symbiotic core integration via lifecycle: https://docs.symbiotic.fi/core-modules/networks#staking-lifecycle
// TODO: determine if you need timestamping similar to cosmos sdk example. Edit yes you will for slashing. See "captureTimestamp". 
// TODO: Parse through MevCommitAVS and make sure translatable reg/dreg functions have the same operators / check the same things. 
// TODO: for example you need to add requires s.t. a validator MUST be opted-in right after registering. 
// TODO: Implement contract owner setting minStake when a vault is registered. Also impl vault registration itself.
// TODO: add function for a validator to "chage vault used for collateral", which involves a delete + new reg. 
contract MevCommitMiddleware is IMevCommitMiddleware, MevCommitMiddlewareStorage,
    Ownable2StepUpgradeable, PausableUpgradeable, UUPSUpgradeable {
    
    // TODO: more modifiers similar to MevCommitAVS

    modifier onlySlashOracle() {
        require(msg.sender == slashOracle, "only slash oracle");
        _;
    }

    // TODO: Define integration with individual vaults, and how you decide on "min stake" per validator
    // for each denom. Price oracle or hardcoded minStake? 

    // TODO: make some sort of integration or fuzz test for two main invariants defined in notion. 

    // TODO: Add things like network epoch duration, ref to core contracts, etc. 
    function initialize(
        uint256 _operatorDeregPeriodBlocks,
        uint256 _validatorDeregPeriodBlocks,
        uint256 _vaultDeregPeriodBlocks,
        address _slashOracle,
        address _owner
    ) public initializer {
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

    /// @notice This function allows a validator to swap its registration with another validator's registration.
    /// This is only possible if both validators are registered to the same operator and vault.
    function swapValRegistrations(bytes[] calldata blsPubkeys1, bytes[] calldata blsPubkeys2) external whenNotPaused {
        require(blsPubkeys1.length == blsPubkeys2.length, "invalid length");
        for (uint256 i = 0; i < blsPubkeys1.length; i++) {
            _swapValRecords(blsPubkeys1[i], blsPubkeys2[i]);
        }
    }

    /// @dev Deletes validator records, only if the associated operator is blacklisted.
    /// Restricted to contract owner.
    /// @notice This function allows the contract owner to combat a greifing scenario where an operator
    /// registers a validator pubkey that it does not control, own, or otherwise manage.
    ///
    /// TODO: Write test for scenario where operator greifs another, and contract owner
    /// has to blacklist that operator, then delete the greifed validator records.
    function deleteValRecords(bytes[] calldata blsPubkeys) external onlyOwner {
        for (uint256 i = 0; i < blsPubkeys.length; i++) {
            _deleteValRecord(blsPubkeys[i]);
        }
    }

    // TODO: IMPORTANT, ABOVE FUNCTION SHOULD NOT BE ONLY OWNER, BUT ALSO OPERATOR.
    // WITH AN ENUMERABLE SET, WE CAN ALLOW ARBITRARY DELETION WHILE MAINTAINING INVARIANTS.
    // TODO: WILL NEED TO DECREMENT COUNTER WHEN REMOVING A RECORD.
    // TODO: ALSO NEED TO SWAP LAST VALIDATOR IN THE VAULT WITH THE DELETED VALIDATOR, PRIOR TO DELETION.
    // TODO: Also remove this section from the notion doc.

    function registerVaults(address[] calldata vaults, address[] calldata operators) external onlyOwner {
        require(vaults.length == operators.length, "invalid length");
        for (uint256 i = 0; i < vaults.length; i++) {
            _registerVault(vaults[i], operators[i]);
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

    // TODO: hook these into symbiotic core
    function _registerOperator(address operator) internal {
        require(!operatorRecords[operator].exists, "operator already registered");
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

    function _setValRecord(bytes calldata blsPubkey, address vault, uint256 priorityIndex) internal {
        validatorRecords[blsPubkey] = ValidatorRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            vault: vault,
            priorityIndex: priorityIndex
        });
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

        _setValRecord(blsPubkey, vault, vaultRecords[vault].priorityIndexCounter);
        emit ValRecordAdded(blsPubkey, msg.sender, vaultRecords[vault].priorityIndexCounter);
        ++vaultRecords[vault].priorityIndexCounter;
    }

    function _requestValDeregistration(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        require(_getOperatorFromValRecord(blsPubkey) == msg.sender, "sender is not operator");
        EventHeightLib.set(validatorRecords[blsPubkey].deregRequestHeight, block.number);
        emit ValidatorDeregistrationRequested(blsPubkey, msg.sender, validatorRecords[blsPubkey].priorityIndex);
    }

    // TODO: test newBlsPubkey could be the same as oldBlsPubkey
    function _swapValRecords(bytes calldata blsPubkey1, bytes calldata blsPubkey2) internal {
        require(validatorRecords[blsPubkey1].exists, "missing val record 1");
        require(validatorRecords[blsPubkey2].exists, "missing val record 2");

        require(validatorRecords[blsPubkey1].vault == validatorRecords[blsPubkey2].vault,
            "vaults do not match");
        
        address vault = validatorRecords[blsPubkey1].vault;
        require(vaultRecords[vault].exists, "vault not registered");
        require(!vaultRecords[vault].deregRequestHeight.exists, "vault dereg request exists");

        require(_getOperatorFromValRecord(blsPubkey1) == msg.sender &&
            _getOperatorFromValRecord(blsPubkey2) == msg.sender, "sender is not operator");
        
        require(operatorRecords[msg.sender].exists, "operator not registered");
        require(!operatorRecords[msg.sender].deregRequestHeight.exists, "operator dereg request exists");
        require(!operatorRecords[msg.sender].isBlacklisted, "operator is blacklisted");

        require(_isValidatorReadyToDeregister(blsPubkey1), "not ready to dereg");
        require(_isValidatorReadyToDeregister(blsPubkey2), "not ready to dereg");
            
        // swap priorities, reset dereg request heights
        uint256 priorityIndex1 = validatorRecords[blsPubkey1].priorityIndex;
        _setValRecord(blsPubkey1, vault, validatorRecords[blsPubkey2].priorityIndex);
        _setValRecord(blsPubkey2, vault, priorityIndex1);
        emit ValRecordsSwapped(blsPubkey1, blsPubkey2, msg.sender,
            // Log new stored priority indexes
            validatorRecords[blsPubkey1].priorityIndex, validatorRecords[blsPubkey2].priorityIndex);
    }

    function _deleteValRecord(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing val record");
        address operator = _getOperatorFromValRecord(blsPubkey);
        require(operatorRecords[operator].exists, "operator not registered");
        require(!operatorRecords[operator].deregRequestHeight.exists, "operator dereg request exists");
        require(operatorRecords[operator].isBlacklisted, "operator is blacklisted");
        uint256 priorityIndex = validatorRecords[blsPubkey].priorityIndex;
        delete validatorRecords[blsPubkey];
        // TODO: subtract from vault priority index counter
        emit ValRecordDeleted(blsPubkey, operator, priorityIndex);
    }

    function _setVaultRecord(address vault, address operator) internal {
        vaultRecords[vault] = VaultRecord({
            exists: true,
            deregRequestHeight: EventHeightLib.EventHeight({
                exists: false,
                blockHeight: 0
            }),
            priorityIndexCounter: 0,
            operator: operator
        });
    }

    function _registerVault(address vault, address operator) internal {
        require(!vaultRecords[vault].exists, "vault already registered");
        _setVaultRecord(vault, operator);
        emit VaultRegistered(vault, operator);
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

    function _slashValidator(bytes calldata blsPubkey) internal {
        require(validatorRecords[blsPubkey].exists, "missing validator record");
        // TODO: slash operator with core
        address operator = _getOperatorFromValRecord(blsPubkey);
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
        // TODO: check liquidity exists to slash if needed
        if (validatorRecords[blsPubkey].priorityIndex > 71) { // where 71 is some threshold defined by amount of liquidity
            return false;
        }
        return true;
    }

    function _isValidatorSlashable(bytes calldata blsPubkey) internal view returns (bool) {
        // USE stakeAt or stake from IBaseDelegator
    }
}
