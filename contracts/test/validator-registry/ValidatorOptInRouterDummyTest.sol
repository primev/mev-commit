// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {ValidatorOptInRouterDummy} from "../../contracts/validator-registry/ValidatorOptInRouterDummy.sol";
import {IValidatorOptInRouter} from "../../contracts/interfaces/IValidatorOptInRouter.sol";
import {IVanillaRegistry} from "../../contracts/interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../../contracts/interfaces/IMevCommitAVS.sol";
import {IMevCommitMiddleware} from "../../contracts/interfaces/IMevCommitMiddleware.sol";

contract ValidatorOptInRouterDummyTest is Test {
    ValidatorOptInRouterDummy public validatorOptInRouter;

    address public owner;
    address public user1;
    address public user2;

    function setUp() public {
        owner = address(0x123456);
        user1 = address(0x123);
        user2 = address(0x456);

        validatorOptInRouter = new ValidatorOptInRouterDummy();
        validatorOptInRouter.initialize(address(0), address(0), address(0), owner);
    }

    function testSettersRevert() public {
        vm.expectRevert("Not implemented");
        validatorOptInRouter.setVanillaRegistry(IVanillaRegistry(address(0)));

        vm.expectRevert("Not implemented");
        validatorOptInRouter.setMevCommitAVS(IMevCommitAVS(address(0)));

        vm.expectRevert("Not implemented");
        validatorOptInRouter.setMevCommitMiddleware(IMevCommitMiddleware(address(0)));
    }

    function testAreValidatorsOptedIn() public {
        // Create test pubkeys where first byte is < 64 (opted in) and > 64 (not opted in)
        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = abi.encodePacked(bytes1(0x01)); // First byte 1 - should be opted in
        valPubkeys[1] = abi.encodePacked(bytes1(0xFF)); // First byte 255 - should not be opted in

        IValidatorOptInRouter.OptInStatus[] memory optInStatuses = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        
        assertEq(optInStatuses.length, 2);

        // First validator should be opted in (first byte < 64)
        assertTrue(optInStatuses[0].isVanillaOptedIn);
        assertTrue(optInStatuses[0].isAvsOptedIn); 
        assertTrue(optInStatuses[0].isMiddlewareOptedIn);

        // Second validator should not be opted in (first byte > 64)
        assertFalse(optInStatuses[1].isVanillaOptedIn);
        assertFalse(optInStatuses[1].isAvsOptedIn);
        assertFalse(optInStatuses[1].isMiddlewareOptedIn);
    }

    function testAreValidatorsOptedInBoundary() public {
        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = abi.encodePacked(bytes1(0x3F)); // First byte 63 - should be opted in
        valPubkeys[1] = abi.encodePacked(bytes1(0x40)); // First byte 64 - should not be opted in

        IValidatorOptInRouter.OptInStatus[] memory optInStatuses = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        
        assertEq(optInStatuses.length, 2);

        // First validator should be opted in (63 < 64)
        assertTrue(optInStatuses[0].isVanillaOptedIn);
        assertTrue(optInStatuses[0].isAvsOptedIn);
        assertTrue(optInStatuses[0].isMiddlewareOptedIn);

        // Second validator should not be opted in (64 >= 64) 
        assertFalse(optInStatuses[1].isVanillaOptedIn);
        assertFalse(optInStatuses[1].isAvsOptedIn);
        assertFalse(optInStatuses[1].isMiddlewareOptedIn);
    }
}
