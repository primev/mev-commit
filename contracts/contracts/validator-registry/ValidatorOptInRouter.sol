// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {ValidatorOptInRouterStorage} from "./ValidatorOptInRouterStorage.sol";
import {IValidatorOptInRouter} from "../interfaces/IValidatorOptInRouter.sol";
import {IValidatorRegistryV1} from "../interfaces/IValidatorRegistryV1.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title ValidatorOptInRouter
/// @notice This contract acts as the top level source of truth for whether a validator 
/// is opted in to mev-commit from either the v1 validator registry or the mev-commit AVS.
contract ValidatorOptInRouter is IValidatorOptInRouter, ValidatorOptInRouterStorage,
    Ownable2StepUpgradeable, UUPSUpgradeable {

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Initializes the contract with the validator registry and mev-commit AVS contracts.
    function initialize(
        address _validatorRegistry,
        address _mevCommitAVS,
        address _owner
    ) external initializer {
        validatorRegistryV1 = IValidatorRegistryV1(_validatorRegistry);
        mevCommitAVS = IMevCommitAVS(_mevCommitAVS);
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    /*
     * @dev implements _authorizeUpgrade from UUPSUpgradeable to enable only
     * the owner to upgrade the implementation contract.
     */
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @notice Allows the owner to set the validator registry V1 contract.
    function setValidatorRegistryV1(IValidatorRegistryV1 _validatorRegistry) external onlyOwner {
        address oldContract = address(validatorRegistryV1);
        validatorRegistryV1 = _validatorRegistry;
        emit ValidatorRegistryV1Set(oldContract, address(validatorRegistryV1));
    }

    /// @notice Allows the owner to set the mev-commit AVS contract.
    function setMevCommitAVS(IMevCommitAVS _mevCommitAVS) external onlyOwner {
        address oldContract = address(mevCommitAVS);
        mevCommitAVS = _mevCommitAVS;
        emit MevCommitAVSSet(oldContract, address(mevCommitAVS));
    }

    /// @notice Returns an array of bools indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (bool[] memory) {
        bool[] memory optedIn = new bool[](valBLSPubKeys.length);
        for (uint256 i = 0; i < valBLSPubKeys.length; i++) {
            optedIn[i] = _isValidatorOptedIn(valBLSPubKeys[i]);
        }
        return optedIn;
    }

    /// @notice Internal function to check if a validator is opted in to mev-commit with either simple staking or restaking.
    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal view returns (bool) {
        if (validatorRegistryV1.isValidatorOptedIn(valBLSPubKey)) {
            return true;
        }
        return mevCommitAVS.isValidatorOptedIn(valBLSPubKey);
    }

    /// @dev Fallback function to revert all calls, ensuring no unintended interactions.
    fallback() external payable {
        revert("Invalid call");
    }

    /// @dev Receive function is disabled for this contract to prevent unintended interactions.
    receive() external payable {
        revert("Invalid call");
    }
}
