// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {VanillaRegistry} from "../../contracts/validator-registry/VanillaRegistry.sol";
import {ValidatorOptInRouter} from "../../contracts/validator-registry/ValidatorOptInRouter.sol";
import {MevCommitAVS} from "../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {VanillaRegistryTest} from "./VanillaRegistryTest.sol";
import {MevCommitAVSTest} from "./avs/MevCommitAVSTest.sol";
import {IVanillaRegistry} from "../../contracts/interfaces/IVanillaRegistry.sol";
import {IMevCommitAVS} from "../../contracts/interfaces/IMevCommitAVS.sol";
import {IMevCommitMiddleware} from "../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddleware} from "../../contracts/validator-registry/middleware/MevCommitMiddleware.sol";
import {MevCommitMiddlewareTestCont} from "./middleware/MevCommitMiddlewareTestCont.sol";

contract ValidatorOptInRouterTest is Test {
    ValidatorOptInRouter public validatorOptInRouter;

    VanillaRegistry public vanillaRegistry;
    VanillaRegistryTest public vanillaRegistryTest;
    MevCommitAVS public mevCommitAVS;
    MevCommitAVSTest public mevCommitAVSTest;
    MevCommitMiddleware public mevCommitMiddleware;
    MevCommitMiddlewareTestCont public mevCommitMiddlewareTest;

    address public owner;
    address public user1;
    address public user2;

    event VanillaRegistrySet(address oldContract, address newContract);
    event MevCommitAVSSet(address oldContract, address newContract);
    event MevCommitMiddlewareSet(address oldContract, address newContract);

    function setUp() public {
        owner = address(0x123456);
        user1 = address(0x123);
        user2 = address(0x456);

        vanillaRegistryTest = new VanillaRegistryTest();
        vanillaRegistryTest.setUp();
        vanillaRegistry = vanillaRegistryTest.validatorRegistry();

        mevCommitAVSTest = new MevCommitAVSTest();
        mevCommitAVSTest.setUp();
        mevCommitAVS = mevCommitAVSTest.mevCommitAVS();

        mevCommitMiddlewareTest = new MevCommitMiddlewareTestCont();
        mevCommitMiddlewareTest.setUp();
        mevCommitMiddleware = mevCommitMiddlewareTest.mevCommitMiddleware();

        address routerProxy = Upgrades.deployUUPSProxy(
            "ValidatorOptInRouter.sol",
            abi.encodeCall(ValidatorOptInRouter.initialize,
            (
                address(vanillaRegistry),
                address(mevCommitAVS),
                address(mevCommitMiddleware),
                owner
            ))
        );
        validatorOptInRouter = ValidatorOptInRouter(payable(routerProxy));
    }

    function testSetters() public {
        IVanillaRegistry newValReg = new VanillaRegistry();
        IVanillaRegistry oldValReg = vanillaRegistry;
        vm.prank(owner);
        vm.expectEmit();
        emit VanillaRegistrySet(address(oldValReg), address(newValReg));
        validatorOptInRouter.setVanillaRegistry(newValReg);
        assertEq(address(validatorOptInRouter.vanillaRegistry()), address(newValReg));

        IMevCommitAVS newMevCommitAVS = new MevCommitAVS();
        IMevCommitAVS oldMevCommitAVS = mevCommitAVS;
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitAVSSet(address(oldMevCommitAVS), address(newMevCommitAVS));
        validatorOptInRouter.setMevCommitAVS(newMevCommitAVS);
        assertEq(address(validatorOptInRouter.mevCommitAVS()), address(newMevCommitAVS));

        IMevCommitMiddleware newMevCommitMiddleware = new MevCommitMiddleware();
        IMevCommitMiddleware oldMevCommitMiddleware = mevCommitMiddleware;
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitMiddlewareSet(address(oldMevCommitMiddleware), address(newMevCommitMiddleware));
        validatorOptInRouter.setMevCommitMiddleware(newMevCommitMiddleware);
        assertEq(address(validatorOptInRouter.mevCommitMiddleware()), address(newMevCommitMiddleware));
    }

    function testAreValidatorsOptedInViaRestaking() public {
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = mevCommitAVSTest.sampleValPubkey1();
        valPubkeys[1] = mevCommitAVSTest.sampleValPubkey2();

        assertTrue(mevCommitAVS.isValidatorOptedIn(valPubkeys[0]));
        assertTrue(mevCommitAVS.isValidatorOptedIn(valPubkeys[1]));
        
        bool[] memory areOptedIn = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        assertEq(areOptedIn.length, 2);
        for (uint256 i = 0; i < areOptedIn.length; ++i) {
            assertTrue(areOptedIn[i]);
        }
    }

    function testAreValidatorsOptedInViaVanillaStaking() public {
        vanillaRegistryTest.testMultiStake();

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = vanillaRegistryTest.user1BLSKey();
        valPubkeys[1] = vanillaRegistryTest.user2BLSKey();

        assertTrue(vanillaRegistry.isValidatorOptedIn(valPubkeys[0]));
        assertTrue(vanillaRegistry.isValidatorOptedIn(valPubkeys[1]));

        bool[] memory areOptedIn = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        assertEq(areOptedIn.length, 2);
        for (uint256 i = 0; i < areOptedIn.length; ++i) {
            assertTrue(areOptedIn[i]);
        }
    }

    function testAreValidatorsOptedInForMevCommitMiddleware() public {
        mevCommitMiddlewareTest.test_registerValidators();

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = mevCommitMiddlewareTest.sampleValPubkey1();
        valPubkeys[1] = mevCommitMiddlewareTest.sampleValPubkey2();

        bool[] memory areOptedIn = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        assertEq(areOptedIn.length, 2);
        for (uint256 i = 0; i < areOptedIn.length; ++i) {
            assertTrue(areOptedIn[i]);
        }
    }
}
