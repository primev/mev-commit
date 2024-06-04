// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../contracts/validator-registry/DelegationValReg.sol";
import "../../contracts/validator-registry/ReputationValReg.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MockERC20 is ERC20 {
    constructor() ERC20("MockERC20", "MERC20") {
        _mint(msg.sender, 1000000 * (10 ** uint256(decimals())));
    }
}

contract DelegationValRegTest is Test {
    DelegationValReg public delegationValReg;
    ReputationValReg public reputationValReg;
    ERC20 public mockERC20;
    address public owner;
    address public user1;
    address public user2;
    address public user3;
    address public delegator1;
    address public delegator2;

    uint256 constant MAX_CONS_ADDRS_PER_EOA = 5;
    uint256 constant MIN_FREEZE_BLOCKS = 100;
    uint256 constant UNFREEZE_FEE = 1 ether;
    uint256 constant WITHDRAW_PERIOD = 25;

    bytes public constant exampleConsAddr1 = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public constant exampleConsAddr2 = hex"a5c99dfdfc69791937ac5efc5d33316cd4e0698be24ef149bbc18f0f25ad92e5e11aafd39701dcdab6d3205ad38c307b";

    event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount);
    event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount);
    event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount);

    function setUp() public {
        owner = address(0x789);
        user1 = address(0x123);
        user2 = address(0x456);
        user3 = address(0x777);
        delegator1 = address(0x776);
        delegator2 = address(0x888);

        address reputationProxy = Upgrades.deployUUPSProxy(
            "ReputationValReg.sol",
            abi.encodeCall(ReputationValReg.initialize, 
            (owner, MAX_CONS_ADDRS_PER_EOA, MIN_FREEZE_BLOCKS, UNFREEZE_FEE))
        );
        reputationValReg = ReputationValReg(payable(reputationProxy));

        vm.startPrank(owner);
        reputationValReg.addWhitelistedEOA(user1, "moniker_user1");
        reputationValReg.addWhitelistedEOA(user2, "moniker_user2");
        vm.stopPrank();

        bytes[] memory consAddrs = new bytes[](1);
        consAddrs[0] = exampleConsAddr1;
        vm.prank(user1);
        reputationValReg.storeConsAddrs(consAddrs);

        consAddrs = new bytes[](1);
        consAddrs[0] = exampleConsAddr2;
        vm.prank(user2);
        reputationValReg.storeConsAddrs(consAddrs);

        mockERC20 = new MockERC20();
        mockERC20.approve(owner, 100 ether);
        mockERC20.transfer(owner, 100 ether);

        address delegationProxy = Upgrades.deployUUPSProxy(
            "DelegationValReg.sol",
            abi.encodeCall(DelegationValReg.initialize, 
            (owner, reputationProxy, WITHDRAW_PERIOD, address(mockERC20)))
        );
        delegationValReg = DelegationValReg(payable(delegationProxy));
    }

    function testSecondInitialize() public {
        vm.prank(owner);
        vm.expectRevert();
        delegationValReg.initialize(owner, address(reputationValReg), WITHDRAW_PERIOD, address(mockERC20));
        vm.stopPrank();
    }

    function testDelegate() public {

        vm.prank(delegator1);
        vm.expectRevert("Validator EOA must be whitelisted");
        delegationValReg.delegate(user3, 1 ether);

        vm.prank(delegator1);
        vm.expectRevert("Amount must be greater than 0");
        delegationValReg.delegate(user1, 0 ether);

        vm.prank(owner);
        mockERC20.transfer(delegator1, 10 ether);

        DelegationValReg.DelegationInfo memory delegationInfo = delegationValReg.getDelegationInfo(delegator1);
        assertEq(uint256(delegationInfo.state), uint256(DelegationValReg.State.nonExistant));
        assertEq(delegationInfo.validatorEOA, address(0));
        assertEq(delegationInfo.amount, 0);
        assertEq(delegationInfo.withdrawHeight, 0);

        vm.startPrank(delegator1);
        mockERC20.approve(address(delegationValReg), 10 ether);
        vm.expectEmit(true, true, true, true);
        emit Delegated(delegator1, user1, 10 ether);
        delegationValReg.delegate(user1, 10 ether);
        vm.stopPrank();

        delegationInfo = delegationValReg.getDelegationInfo(delegator1);
        assertEq(uint256(delegationInfo.state), uint256(DelegationValReg.State.active));
        assertEq(delegationInfo.validatorEOA, user1);
        assertEq(delegationInfo.amount, 10 ether);
        assertEq(delegationInfo.withdrawHeight, 0);

        vm.prank(delegator1);
        vm.expectRevert("Delegation must not exist for sender");
        delegationValReg.delegate(user2, 1 ether);
    }

    function testChangeDelegation() public {
        testDelegate();

        vm.prank(delegator2); 
        vm.expectRevert("Active delegation must exist for sender");
        delegationValReg.changeDelegation(user1);

        vm.prank(delegator1);
        vm.expectRevert("New validator EOA must be whitelisted");
        delegationValReg.changeDelegation(user3);

        DelegationValReg.DelegationInfo memory delegationInfo = delegationValReg.getDelegationInfo(delegator1);
        assertEq(uint256(delegationInfo.state), uint256(DelegationValReg.State.active));
        assertEq(delegationInfo.validatorEOA, user1);
        assertEq(delegationInfo.amount, 10 ether);
        assertEq(delegationInfo.withdrawHeight, 0);

        vm.prank(delegator1);
        vm.expectEmit(true, true, true, true);
        emit DelegationChanged(delegator1, user1, user2, 10 ether);
        delegationValReg.changeDelegation(user2);

        delegationInfo = delegationValReg.getDelegationInfo(delegator1);
        assertEq(uint256(delegationInfo.state), uint256(DelegationValReg.State.active));
        assertEq(delegationInfo.validatorEOA, user2);
        assertEq(delegationInfo.amount, 10 ether);
        assertEq(delegationInfo.withdrawHeight, 0);
    }
}
