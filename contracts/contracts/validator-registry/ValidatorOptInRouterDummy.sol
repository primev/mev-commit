// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IValidatorOptInRouter} from "../interfaces/IValidatorOptInRouter.sol";
import {IVanillaRegistry} from "../interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../interfaces/IMevCommitAVS.sol";
import {IMevCommitMiddleware} from "../interfaces/IMevCommitMiddleware.sol";

/// @title ValidatorOptInRouterDummy
/// @notice A dummy implementation of ValidatorOptInRouter that returns true 25% of the time
/// for testing purposes.
contract ValidatorOptInRouterDummy is IValidatorOptInRouter {

    /// @notice Initializes the contract (no-op for dummy implementation)
    function initialize(
        address _vanillaRegistry,
        address _mevCommitAVS,
        address _mevCommitMiddleware,
        address _owner
    ) external {
        // No-op
    }

    /// @notice Dummy implementation that always reverts
    function setVanillaRegistry(IVanillaRegistry _vanillaRegistry) external {
        revert("Not implemented");
    }

    /// @notice Dummy implementation that always reverts
    function setMevCommitAVS(IMevCommitAVS _mevCommitAVS) external {
        revert("Not implemented");
    }

    /// @notice Dummy implementation that always reverts
    function setMevCommitMiddleware(IMevCommitMiddleware _mevCommitMiddleware) external {
        revert("Not implemented");
    }

    /// @notice Returns an array of OptInStatus structs with 25% probability of being opted in
    function areValidatorsOptedIn(bytes[] calldata valBLSPubKeys) external pure returns (OptInStatus[] memory) {
        uint256 len = valBLSPubKeys.length;
        OptInStatus[] memory optInStatuses = new OptInStatus[](len);
        for (uint256 i = 0; i < len; ++i) {
            optInStatuses[i] = _isValidatorOptedIn(valBLSPubKeys[i]);
        }
        return optInStatuses;
    }

    /// @notice Internal function that returns true 25% of the time based on the first byte of the pubkey
    function _isValidatorOptedIn(bytes calldata valBLSPubKey) internal pure returns (OptInStatus memory) {
        OptInStatus memory optInStatus;
        
        // Use first byte of pubkey to determine opt-in status
        // If first byte is < 64 (25% of 256), return true
        bool isOptedIn = uint8(valBLSPubKey[0]) < 64;
        
        optInStatus.isVanillaOptedIn = isOptedIn;
        optInStatus.isAvsOptedIn = isOptedIn;
        optInStatus.isMiddlewareOptedIn = isOptedIn;
        
        return optInStatus;
    }
}
