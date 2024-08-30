// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.25;

import {Test} from "forge-std/Test.sol";
import {Whitelist} from "../../contracts/standard-bridge/Whitelist.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

// Tests the Whitelist contract.
// Note precompile interactions to mint/burn must be tested manually. 
contract WhitelistTest is Test {

    address admin;
    address normalBidder;
    address addressInstance;
    Whitelist whitelist;

    function setUp() public {
        admin = address(this); // Original contract deployer as admin
        normalBidder = address(0x100);
        addressInstance = address(0x200);

        address whitelistProxy = Upgrades.deployUUPSProxy(
            "Whitelist.sol",
            abi.encodeCall(Whitelist.initialize, (admin))
        ); 
        whitelist = Whitelist(payable(whitelistProxy));
    }

    function test_IsWhitelisted() public {
        assertFalse(whitelist.isWhitelisted(addressInstance));
        vm.prank(admin);
        whitelist.addToWhitelist(addressInstance);
        assertTrue(whitelist.isWhitelisted(addressInstance));
    }

    function test_AdminAddToWhitelist() public {
        vm.prank(admin);
        whitelist.addToWhitelist(addressInstance);
        assertTrue(whitelist.isWhitelisted(addressInstance));
    }

    function test_AdminRemoveFromWhitelist() public {
        vm.prank(admin);
        whitelist.addToWhitelist(addressInstance);
        assertTrue(whitelist.isWhitelisted(addressInstance));
        vm.prank(admin);
        whitelist.removeFromWhitelist(addressInstance);
        assertFalse(whitelist.isWhitelisted(addressInstance));
    }

    function test_RevertNormalBidderAddToWhitelist() public {
        vm.prank(normalBidder);
        vm.expectRevert(); // Only owner can add to whitelist
        whitelist.addToWhitelist(addressInstance);
    }

    function test_RevertNormalBidderRemoveFromWhitelist() public {
        vm.prank(admin);
        whitelist.addToWhitelist(addressInstance);
        assertTrue(whitelist.isWhitelisted(addressInstance));
        vm.prank(normalBidder);
        vm.expectRevert(); // Only owner can remove from whitelist
        whitelist.removeFromWhitelist(addressInstance);
    }
}
