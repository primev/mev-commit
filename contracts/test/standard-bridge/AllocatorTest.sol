// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

import {Test} from "forge-std/Test.sol";
import {Allocator} from "../../contracts/standard-bridge/Allocator.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";

// Tests the allocator contract.
// Note precompile interactions to mint/burn must be tested manually. 
contract AllocatorTest is Test {

    address admin;
    address normalBidder;
    address addressInstance;
    Allocator allocator;

    function setUp() public {
        admin = address(this); // Original contract deployer as admin
        normalBidder = address(0x100);
        addressInstance = address(0x200);

        address allocatorProxy = Upgrades.deployUUPSProxy(
            "Allocator.sol",
            abi.encodeCall(Allocator.initialize, (admin))
        ); 
        allocator = Allocator(payable(allocatorProxy));
    }

    function test_IsWhitelisted() public {
        assertFalse(allocator.isWhitelisted(addressInstance));
        vm.prank(admin);
        allocator.addToWhitelist(addressInstance);
        assertTrue(allocator.isWhitelisted(addressInstance));
    }

    function test_AdminAddToWhitelist() public {
        vm.prank(admin);
        allocator.addToWhitelist(addressInstance);
        assertTrue(allocator.isWhitelisted(addressInstance));
    }

    function test_AdminRemoveFromWhitelist() public {
        vm.prank(admin);
        allocator.addToWhitelist(addressInstance);
        assertTrue(allocator.isWhitelisted(addressInstance));
        vm.prank(admin);
        allocator.removeFromWhitelist(addressInstance);
        assertFalse(allocator.isWhitelisted(addressInstance));
    }

    function test_RevertNormalBidderAddToWhitelist() public {
        vm.prank(normalBidder);
        vm.expectRevert(); // Only owner can add to whitelist
        allocator.addToWhitelist(addressInstance);
    }

    function test_RevertNormalBidderRemoveFromWhitelist() public {
        vm.prank(admin);
        allocator.addToWhitelist(addressInstance);
        assertTrue(allocator.isWhitelisted(addressInstance));
        vm.prank(normalBidder);
        vm.expectRevert(); // Only owner can remove from whitelist
        allocator.removeFromWhitelist(addressInstance);
    }
}
