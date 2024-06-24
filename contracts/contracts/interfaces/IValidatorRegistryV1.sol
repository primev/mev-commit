// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

contract IValidatorRegistryV1 {

    /// @dev Event emitted when a validator is staked.
    event Staked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator is unstaked.
    event Unstaked(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator's stake is withdrawn.
    event StakeWithdrawn(address indexed msgSender, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Event emitted when a validator is slashed.
    event Slashed(address indexed msgSender, address indexed slashReceiver, address indexed withdrawalAddress, bytes valBLSPubKey, uint256 amount);

    /// @dev Struct representing a validator staked with the registry.
    struct StakedValidator {
        uint256 balance;
        address withdrawalAddress;
        uint256 unstakeBlockNum;
    }
}
