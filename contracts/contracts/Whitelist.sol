// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

// Contract that allows an admin to add/remove addresses from the whitelist,
// and allows whitelisted addresses to mint native tokens.
//
// The whitelist contract's create2 address must be funded on genesis.
contract Whitelist is Ownable2StepUpgradeable, UUPSUpgradeable {

    mapping(address => bool) public whitelistedAddresses;

    function initialize(address _owner) external initializer {
        __Ownable_init(_owner);
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // Receiver for native tokens to be "burnt"
    receive() external payable {}

    function addToWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = true;
    }

    function removeFromWhitelist(address _address) external onlyOwner {
        whitelistedAddresses[_address] = false;
    }

    // "Mints" native tokens (transfer ether from this contract) if the sender is whitelisted.
    function mint(address _mintTo, uint256 _amount) external {
        require(isWhitelisted(msg.sender), "Sender is not whitelisted");
        require(address(this).balance >= _amount, "Insufficient contract balance");
        (bool success, ) = _mintTo.call{value: _amount}("");
        require(success, "Transfer to _mintTo failed");
    }

    function isWhitelisted(address _address) public view returns (bool) {
        return whitelistedAddresses[_address];
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}
}
