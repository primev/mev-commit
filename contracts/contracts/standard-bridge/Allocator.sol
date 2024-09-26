// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {IAllocator} from "../interfaces/IAllocator.sol";
import {AllocatorStorage} from "./AllocatorStorage.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {Errors} from "../utils/Errors.sol";

/// @title Allocator
/// @notice Contract that allows an admin to add/remove addresses from a whitelist,
/// which can "mint" native ETH on the mev-commit chain, enabling native ETH bridging.
/// @dev This contract must be funded (ideally on genesis) prior to being an effective minting entity.
/// @dev Contracts which "mint" should implement "burning" eth as transferring it to this contract with no data.
contract Allocator is AllocatorStorage, IAllocator,
    Ownable2StepUpgradeable, UUPSUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable {

    function initialize(address _owner) external initializer {
        __Ownable_init(_owner);
        __Pausable_init();
        __ReentrancyGuard_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Receiver for native eth to be "burnt"
    receive() external payable {}

    /// @dev Fallback function is disabled for this contract to prevent unintended interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    function addToWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = true;
    }

    function removeFromWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = false;
    }

    // "Mints" native eth (transfer ether from this contract) if the sender is whitelisted.
    function mint(address _mintTo, uint256 _amount) external whenNotPaused nonReentrant {
        require(isWhitelisted(msg.sender), SenderNotWhitelisted(msg.sender));
        require(address(this).balance >= _amount, InsufficientContractBalance(address(this).balance, _amount));
        (bool success, ) = _mintTo.call{value: _amount}("");
        require(success, TransferFailed(_mintTo, _amount));
    }

    /// @dev Allows the owner to pause the contract.
    function pause() external onlyOwner {
        _pause();
    }

    /// @dev Allows the owner to unpause the contract.
    function unpause() external onlyOwner {
        _unpause();
    }

    function isWhitelisted(address _address) public view returns (bool) {
        return whitelistedAddresses[_address];
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
