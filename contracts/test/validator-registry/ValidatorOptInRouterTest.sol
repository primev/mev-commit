// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import "../../contracts/validator-registry/ValidatorRegistryV1.sol";
import "../../contracts/validator-registry/ValidatorOptInRouter.sol";
import "../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import "./ValidatorRegistryV1Test.sol";
import "./avs/MevCommitAVSTest.sol";

contract ValidatorOptInRouterTest is Test {
    ValidatorOptInRouter public validatorOptInRouter;

    ValidatorRegistryV1 public validatorRegistry;
    ValidatorRegistryV1Test public validatorRegistryTest;
    MevCommitAVS public mevCommitAVS;
    MevCommitAVSTest public mevCommitAVSTest;

    address public owner;
    address public user1;
    address public user2;

    event ValidatorRegistryV1Set(address oldContract, address newContract);
    event MevCommitAVSSet(address oldContract, address newContract);

    function setUp() public {
        owner = address(0x123456);
        user1 = address(0x123);
        user2 = address(0x456);

        validatorRegistryTest = new ValidatorRegistryV1Test();
        validatorRegistryTest.setUp();
        validatorRegistry = validatorRegistryTest.validatorRegistry();

        mevCommitAVSTest = new MevCommitAVSTest();
        mevCommitAVSTest.setUp();
        mevCommitAVS = mevCommitAVSTest.mevCommitAVS();

        address routerProxy = Upgrades.deployUUPSProxy(
            "ValidatorOptInRouter.sol",
            abi.encodeCall(ValidatorOptInRouter.initialize, (address(validatorRegistry), address(mevCommitAVS), owner))
        );
        validatorOptInRouter = ValidatorOptInRouter(payable(routerProxy));
    }

    function testSetters() public {
        IValidatorRegistryV1 newValRegV1 = new ValidatorRegistryV1();
        IValidatorRegistryV1 oldValRegV1 = validatorRegistry;
        vm.prank(owner);
        vm.expectEmit();
        emit ValidatorRegistryV1Set(address(oldValRegV1), address(newValRegV1));
        validatorOptInRouter.setValidatorRegistryV1(newValRegV1);
        assertEq(address(validatorOptInRouter.validatorRegistryV1()), address(newValRegV1));

        IMevCommitAVS newMevCommitAVS = new MevCommitAVS();
        IMevCommitAVS oldMevCommitAVS = mevCommitAVS;
        vm.prank(owner);
        vm.expectEmit();
        emit MevCommitAVSSet(address(oldMevCommitAVS), address(newMevCommitAVS));
        validatorOptInRouter.setMevCommitAVS(newMevCommitAVS);
        assertEq(address(validatorOptInRouter.mevCommitAVS()), address(newMevCommitAVS));
    }

    function testAreValidatorsOptedInViaRestaking() public {
        mevCommitAVSTest.testRegisterValidatorsByPodOwners();

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = bytes("valPubkey1");
        valPubkeys[1] = bytes("valPubkey2");

        assertTrue(mevCommitAVS.isValidatorOptedIn(valPubkeys[0]));
        assertTrue(mevCommitAVS.isValidatorOptedIn(valPubkeys[1]));
        
        bool[] memory areOptedIn = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        assertEq(areOptedIn.length, 2);
        for (uint256 i = 0; i < areOptedIn.length; i++) {
            assertTrue(areOptedIn[i]);
        }
    }

    function testAreValidatorsOptedInViaSimpleStaking() public {
        validatorRegistryTest.testMultiStake();

        bytes[] memory valPubkeys = new bytes[](2);
        valPubkeys[0] = validatorRegistryTest.user1BLSKey();
        valPubkeys[1] = validatorRegistryTest.user2BLSKey();

        assertTrue(validatorRegistry.isValidatorOptedIn(valPubkeys[0]));
        assertTrue(validatorRegistry.isValidatorOptedIn(valPubkeys[1]));

        bool[] memory areOptedIn = validatorOptInRouter.areValidatorsOptedIn(valPubkeys);
        assertEq(areOptedIn.length, 2);
        for (uint256 i = 0; i < areOptedIn.length; i++) {
            assertTrue(areOptedIn[i]);
        }
    }
}
