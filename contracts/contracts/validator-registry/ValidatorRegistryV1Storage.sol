// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {IValidatorRegistryV1} from "../interfaces/IValidatorRegistryV1.sol";

contract ValidatorRegistryV1Storage { 
    /// @dev Minimum stake required for validators. 
    uint256 public minStake;
    
    /// @dev Amount of ETH to slash per validator pubkey when a slash is invoked.
    uint256 public slashAmount;

    /// @dev Permissioned account that is able to invoke slashes.
    address public slashOracle; 

    /// @dev Account to receive all slashed ETH.
    address public slashReceiver;

    /// @dev Number of blocks required between unstake initiation and withdrawal.
    uint256 public unstakePeriodBlocks;

    mapping(bytes => IValidatorRegistryV1.StakedValidator) public stakedValidators;
}
