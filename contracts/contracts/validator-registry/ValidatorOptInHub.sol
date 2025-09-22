// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {ValidatorOptInHubStorage} from "./ValidatorOptInHubStorage.sol";
import {IValidatorOptInHub} from "../interfaces/IValidatorOptInHub.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Errors} from "../utils/Errors.sol";
import {IRegistry} from "../interfaces/IRegistry.sol";

/// @title ValidatorOptInHub
/// @notice This contract acts as the top level source of truth for whether a validator 
/// is opted in to mev-commit from any of the registries.
contract ValidatorOptInHub is IValidatorOptInHub, ValidatorOptInHubStorage,
    Ownable2StepUpgradeable, UUPSUpgradeable {

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Receive function is disabled for this contract to prevent unintended interactions.
    receive() external payable {
        revert Errors.InvalidReceive();
    }

    /// @dev Fallback function to revert all calls, ensuring no unintended interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Initializes the contract with the validator registry and mev-commit AVS contracts.
    function initialize(
        address[] calldata _registries,
        address _owner
    ) external initializer {
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
        uint256 len = _registries.length;
        for (uint256 i = 0; i < len; ++i) {
            _validateRegistry(_registries[i]);
            registries.push(IRegistry(_registries[i]));
        }
    }

    // --- admin ---
    function addRegistry(address registry) external onlyOwner {
        _validateRegistry(registry);
        registries.push(IRegistry(registry));
        emit RegistryAdded(registries.length - 1, registry);
    }

    /// Pass in index and old registry address for safer replacement.
    function updateRegistry(uint256 index, address oldRegistry, address newRegistry) external onlyOwner {
        require(index < registries.length, InvalidIndex());
        require(oldRegistry != address(0), ZeroAddress());
        _validateRegistry(newRegistry);
        if (address(registries[index]) == oldRegistry) {
            registries[index] = IRegistry(newRegistry);
            emit RegistryReplaced(index, oldRegistry, newRegistry);
            return;
        }
        revert IndexRegistryMismatch();
    }

    /// Pass in index and registry address for safer removal.
    function removeRegistry(uint256 index, address registry) external onlyOwner {
        require(registry != address(0), ZeroAddress());
        require(index < registries.length, InvalidIndex());
        if (address(registries[index]) == registry) {
            registries[index] = IRegistry(address(0));
            emit RegistryRemoved(index, registry);
            return;
        }
        revert IndexRegistryMismatch();
    }

    /// @notice Returns an array of bool lists indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedInList(bytes[] calldata valBLSPubKeys) external view returns (bool[][] memory) {
        uint256 len = valBLSPubKeys.length;
        bool[][] memory _optInStatuses = new bool[][](len);
        for (uint256 i = 0; i < len; ++i) {
            _optInStatuses[i] = _isValidatorOptedInList(valBLSPubKeys[i]);
        }
        return _optInStatuses;
    }

    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory) {
        uint256 len = valBLSPubKeys.length;
        bool[] memory _optInStatuses = new bool[](len);
        for (uint256 i = 0; i < len; ++i) {
            _optInStatuses[i] = _isValidatorOptedIn(valBLSPubKeys[i]);
        }
        return _optInStatuses;
    }

    function isValidatorOptedInList(bytes calldata valPubKey) external view returns (bool[] memory) {
        return _isValidatorOptedInList(valPubKey);
    }

    function isValidatorOptedIn(bytes calldata valPubKey) external view returns (bool) {
        return _isValidatorOptedIn(valPubKey);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    //Sanity check for owner operations
    function _validateRegistry(address registry) internal view {
        require(registry != address(0), ZeroAddress());
        if (registry == address(0)) revert IValidatorOptInHub.ZeroAddress();

        // 1) Must be a contract
        uint256 size;
        assembly { size := extcodesize(registry) }
        if (size == 0) revert IValidatorOptInHub.InvalidRegistry();

        // 2) Must implement isValidatorOptedIn(bytes) -> bool and be STATICCALL-safe
        (bool ok, bytes memory returnData) = registry.staticcall(
            abi.encodeWithSelector(IRegistry.isValidatorOptedIn.selector, bytes(""))
        );
        // ok must be true and return data must be exactly 32 bytes (bool)
        if (!ok || returnData.length != 32) revert IValidatorOptInHub.InvalidRegistry();
    }

    /// @notice Internal function to check if a validator is opted in to mev-commit with any of the registries.
    /// @return bool list indicating whether the validator is opted in to each registry.
    function _isValidatorOptedInList(bytes calldata valPubKey) internal view returns (bool[] memory) {
        bool[] memory _optInStatuses = new bool[](registries.length);
        uint256 len = registries.length;
        for (uint256 i = 0; i < len; ++i) {
            if (address(registries[i]) == address(0)) {
                _optInStatuses[i] = false;
            } else {
                _optInStatuses[i] = registries[i].isValidatorOptedIn(valPubKey);
            }
        }
        return _optInStatuses;
    }

    function _isValidatorOptedIn(bytes calldata valPubKey) internal view returns (bool) {
        uint256 len = registries.length;
        for (uint256 i = 0; i < len; ++i) {
            if (address(registries[i]) != address(0)) {
                if (registries[i].isValidatorOptedIn(valPubKey)) {
                    return true;
                }
            }
        }
        return false;
    }
}
