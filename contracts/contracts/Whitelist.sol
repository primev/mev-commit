// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

// Contract that allows an admin to add/remove addresses from the whitelist,
// and allows whitelisted addresses to mint native tokens.
//
// The whitelist contract's create2 address must be funded on genesis.
contract Whitelist is OwnableUpgradeable {

    mapping(address => bool) public whitelistedAddresses;

    function initialize(address _owner) external initializer {
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function addToWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = true;
    }

    function removeFromWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = false;
    }

    function isWhitelisted(address _address) public view returns (bool) {
        return whitelistedAddresses[_address];
    }

    // "Mints" native tokens (transfer ether from this contract) if the sender is whitelisted.
    function mint(address _mintTo, uint256 _amount) external {
        require(isWhitelisted(msg.sender), "Sender is not whitelisted");
        require(address(this).balance >= _amount, "Insufficient contract balance");
        payable(_mintTo).transfer(_amount);
    }

    // Receiver for native tokens to be "burnt"
    receive() external payable {}
}
