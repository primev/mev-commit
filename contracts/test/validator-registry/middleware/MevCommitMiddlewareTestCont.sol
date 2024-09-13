// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable func-name-mixedcase

import {IMevCommitMiddleware} from "../../../contracts/interfaces/IMevCommitMiddleware.sol";
import {MevCommitMiddlewareTest} from "./MevCommitMiddlewareTest.sol";
import {MockVetoSlasher} from "./MockVetoSlasher.sol";
import {MockInstantSlasher} from "./MockInstantSlasher.sol";

contract MevCommitMiddlewareTestCont is MevCommitMiddlewareTest {

    function setUp() public override {
        super.setUp();
    }

    function test_registerValidatorsVaultReverts() public {
        test_registerOperators();
        address operator1 = vm.addr(0x1117);

        bytes[][] memory blsPubkeys = new bytes[][](2);
        blsPubkeys[0] = new bytes[](2);
        blsPubkeys[0][0] = sampleValPubkey1;
        blsPubkeys[0][1] = sampleValPubkey2;
        blsPubkeys[1] = new bytes[](1);
        blsPubkeys[1][0] = sampleValPubkey3;

        address[] memory vaults = new address[](2);
        vaults[0] = address(vault1);
        vaults[1] = address(vault2);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotEntity.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        vm.prank(address(vault1));
        vaultFactoryMock.register();

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultNotRegistered.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);

        uint256[] memory slashAmounts = new uint256[](2);
        slashAmounts[0] = 10;
        slashAmounts[1] = 20;

        mockDelegator1.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());
        mockDelegator2.setType(mevCommitMiddleware.NETWORK_RESTAKE_DELEGATOR_TYPE());

        MockInstantSlasher mockSlasher1 = new MockInstantSlasher(mevCommitMiddleware.INSTANT_SLASHER_TYPE());
        MockVetoSlasher mockSlasher2 = new MockVetoSlasher(mevCommitMiddleware.VETO_SLASHER_TYPE(), address(0));

        vault1.setSlasher(address(mockSlasher1));
        vault2.setSlasher(address(mockSlasher2));

        vault1.setEpochDuration(151);
        vault2.setEpochDuration(152);

        vm.prank(address(vault1));
        vaultFactoryMock.register();
        vm.prank(address(vault2));
        vaultFactoryMock.register();

        vm.prank(owner);
        mevCommitMiddleware.registerVaults(vaults, slashAmounts);

        vm.prank(owner);
        mevCommitMiddleware.requestVaultDeregistrations(vaults);

        vm.prank(operator1);
        vm.expectRevert(
            abi.encodeWithSelector(IMevCommitMiddleware.VaultDeregRequestExists.selector, vault1)
        );
        mevCommitMiddleware.registerValidators(blsPubkeys, vaults);
    }

    // TODO: val reg cycle
}
