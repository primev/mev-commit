// SPDX-License-Identifier: BSL
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../contracts/validator-registry/ReputationValReg.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";


contract ReputationValRegTest is Test {
    ReputationValReg public reputationValReg;
    address public owner;
    address public user1;
    address public user2;
    address public user3;

    uint256 public constant MAX_CONS_ADDRS_PER_EOA = 5;
    uint256 public constant MIN_FREEZE_BLOCKS = 100;
    uint256 public constant UNFREEZE_FEE = 1 ether;

    bytes public constant exampleConsAddr1 = hex"96db1884af7bf7a1b57c77222723286a8ce3ef9a16ab6c5542ec5160662d450a1b396b22fc519679adae6ad741547268";
    bytes public constant exampleConsAddr2 = hex"a5c99dfdfc69791937ac5efc5d33316cd4e0698be24ef149bbc18f0f25ad92e5e11aafd39701dcdab6d3205ad38c307b";

    event WhitelistedEOAAdded(address indexed eoa, string moniker);
    event WhitelistedEOADeleted(address indexed eoa, string moniker);
    event EOAFrozen(address indexed eoa, string moniker);
    event EOAUnfrozen(address indexed eoa, string moniker);
    event ConsAddrStored(bytes consAddr, address indexed eoa, string moniker);
    event ConsAddrDeleted(bytes consAddr, address indexed eoa, string moniker);

    function setUp() public {
        owner = address(0x789);
        user1 = address(0x123);
        user2 = address(0x456);
        user3 = address(0x777);
        
        assertEq(exampleConsAddr1.length, 48);
        assertEq(exampleConsAddr2.length, 48);
        
        address proxy = Upgrades.deployUUPSProxy(
            "ReputationValReg.sol",
            abi.encodeCall(ReputationValReg.initialize, 
            (owner, MAX_CONS_ADDRS_PER_EOA, MIN_FREEZE_BLOCKS, UNFREEZE_FEE))
        );
        reputationValReg = ReputationValReg(payable(proxy));
    }

    function testSecondInitialize() public {
        vm.prank(owner);
        vm.expectRevert();
        reputationValReg.initialize(owner, MAX_CONS_ADDRS_PER_EOA, MIN_FREEZE_BLOCKS, UNFREEZE_FEE);
        vm.stopPrank();
    }
    
    function testNonWhitelistedEOA() public view {
        assertFalse(reputationValReg.isEOAWhitelisted(user1));
        (ReputationValReg.State state, uint256 numConsAddrsStored, 
            uint256 freezeHeight, string memory moniker) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state), uint256(ReputationValReg.State.NotWhitelisted));
        assertEq(numConsAddrsStored, 0);
        assertEq(freezeHeight, 0);
        assertEq(moniker, "");
    }

    function testNonExistantConsAddr() public view{
        // TODO
    }

    function testAddWhitelistedEOA() public {

        vm.startPrank(user1);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user1));
        reputationValReg.addWhitelistedEOA(address(0), "bob");
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectRevert("Invalid address");
        reputationValReg.addWhitelistedEOA(address(0), "bob");
        vm.stopPrank();

        uint256 initialOwnerBalance = owner.balance;
        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit WhitelistedEOAAdded(user1, "bob");
        reputationValReg.addWhitelistedEOA(user1, "bob");
        vm.stopPrank();

        assertEq(owner.balance, initialOwnerBalance);
        assertTrue(reputationValReg.isEOAWhitelisted(user1));
        (ReputationValReg.State state, uint256 numConsAddrsStored, 
            uint256 freezeHeight, string memory moniker) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state), uint256(ReputationValReg.State.Active));
        assertEq(numConsAddrsStored, 0);
        assertEq(freezeHeight, 0);
        assertEq(moniker, "bob");

        vm.startPrank(owner);
        vm.expectRevert("EOA must not already be whitelisted");
        reputationValReg.addWhitelistedEOA(user1, "bob");
        vm.stopPrank();

        initialOwnerBalance = owner.balance;
        vm.startPrank(owner);
        reputationValReg.addWhitelistedEOA(user2, "alice");
        vm.stopPrank();
        assertEq(owner.balance, initialOwnerBalance);

        assertTrue(reputationValReg.isEOAWhitelisted(user1));
        assertTrue(reputationValReg.isEOAWhitelisted(user2));
        (ReputationValReg.State state2, uint256 numConsAddrsStored2, 
            uint256 freezeHeight2, string memory moniker2) = reputationValReg.getWhitelistedEOAInfo(user2);
        assertEq(uint256(state2), uint256(ReputationValReg.State.Active));
        assertEq(numConsAddrsStored2, 0);
        assertEq(freezeHeight2, 0);
        assertEq(moniker2, "alice");
    }

    function testDeleteWhitelistedEOA() public {
        testAddWhitelistedEOA();

        vm.startPrank(user3);
        vm.expectRevert("Only owner or EOA itself can delete whitelisted EOA");
        reputationValReg.deleteWhitelistedEOA(user1);
        vm.expectRevert("EOA must be whitelisted");
        reputationValReg.deleteWhitelistedEOA(user3);
        vm.stopPrank();

        assertTrue(reputationValReg.isEOAWhitelisted(user1));
        assertTrue(reputationValReg.isEOAWhitelisted(user2));

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit WhitelistedEOADeleted(user1, "bob");
        reputationValReg.deleteWhitelistedEOA(user1);
        vm.stopPrank();

        assertFalse(reputationValReg.isEOAWhitelisted(user1));
        assertTrue(reputationValReg.isEOAWhitelisted(user2));
        (ReputationValReg.State state1, uint256 numConsAddrsStored1, 
            uint256 freezeHeight1, string memory moniker1) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state1), uint256(ReputationValReg.State.NotWhitelisted));
        assertEq(numConsAddrsStored1, 0);
        assertEq(freezeHeight1, 0);
        assertEq(moniker1, "");
        (ReputationValReg.State state2, uint256 numConsAddrsStored2, 
            uint256 freezeHeight2, string memory moniker2) = reputationValReg.getWhitelistedEOAInfo(user2);
        assertEq(uint256(state2), uint256(ReputationValReg.State.Active));
        assertEq(numConsAddrsStored2, 0);
        assertEq(freezeHeight2, 0);
        assertEq(moniker2, "alice");

        vm.startPrank(user2);
        vm.expectEmit(true, true, true, true);
        emit WhitelistedEOADeleted(user2, "alice");
        reputationValReg.deleteWhitelistedEOA(user2);
        vm.stopPrank();

        assertFalse(reputationValReg.isEOAWhitelisted(user2));
    }

}
