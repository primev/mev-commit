// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

import {ValidatorOptInRouterStorage} from "./ValidatorOptInRouterStorage.sol";
import {IValidatorOptInRouter} from "../interfaces/IValidatorOptInRouter.sol";
import {IVanillaRegistry} from "../interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Errors} from "../utils/Errors.sol";
import {IMevCommitMiddleware} from "../interfaces/IMevCommitMiddleware.sol";

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
        address _vanillaRegistry,
        address _mevCommitAVS,
        address _mevCommitMiddleware,
        address _owner
    ) external initializer {
        vanillaRegistry = IVanillaRegistry(_vanillaRegistry);
        mevCommitAVS = IMevCommitAVS(_mevCommitAVS);
        mevCommitMiddleware = IMevCommitMiddleware(_mevCommitMiddleware);
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    /// @notice Allows the owner to set the vanilla registry contract.
    function setVanillaRegistry(IVanillaRegistry _vanillaRegistry) external onlyOwner {
        address oldContract = address(vanillaRegistry);
        vanillaRegistry = _vanillaRegistry;
        emit VanillaRegistrySet(oldContract, address(vanillaRegistry));
    }

    /// @notice Allows the owner to set the mev-commit AVS contract.
    function setMevCommitAVS(IMevCommitAVS _mevCommitAVS) external onlyOwner {
        address oldContract = address(mevCommitAVS);
        mevCommitAVS = _mevCommitAVS;
        emit MevCommitAVSSet(oldContract, address(mevCommitAVS));
    }

    /// @notice Allows the owner to set the mev-commit middleware contract.
    function setMevCommitMiddleware(IMevCommitMiddleware _mevCommitMiddleware) external onlyOwner {
        address oldContract = address(mevCommitMiddleware);
        mevCommitMiddleware = _mevCommitMiddleware;
        emit MevCommitMiddlewareSet(oldContract, address(mevCommitMiddleware));
    }

    /// @notice Returns an array of OptInStatus structs indicating whether each validator pubkey is opted in to mev-commit.
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external view returns (OptInStatus[] memory) {
        uint256 len = valBLSPubKeys.length;
        OptInStatus[] memory optInStatuses = new OptInStatus[](len);
        for (uint256 i = 0; i < len; ++i) {
            optInStatuses[i] = _isValidatorOptedIn(valBLSPubKeys[i]);
        }
        return optInStatuses;
    }

    /*
     * @dev implements _authorizeUpgrade from UUPSUpgradeable to enable only
     * the owner to upgrade the implementation contract.
     */
    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @notice Internal function to check if a validator is opted in to mev-commit with either simple staking or restaking.
    /// @return OptInStatus struct indicating whether the validator is opted in to vanilla registry, 
    /// eigen avs registry, and/or symbiotic middleware registry.
    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal view returns (OptInStatus memory) {
        OptInStatus memory optInStatus;
        optInStatus.isVanillaOptedIn = vanillaRegistry.isValidatorOptedIn(valBLSPubKey);
        optInStatus.isAvsOptedIn = mevCommitAVS.isValidatorOptedIn(valBLSPubKey);
        optInStatus.isMiddlewareOptedIn = mevCommitMiddleware.isValidatorOptedIn(valBLSPubKey);
        return optInStatus;
    }
}
