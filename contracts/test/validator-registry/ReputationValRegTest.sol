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

    function testWhitelistCycle() public {
        testDeleteWhitelistedEOA();
        testAddWhitelistedEOA();
    }

    function testFreeze() public {
        testAddWhitelistedEOA();

        vm.startPrank(user3);
        vm.expectRevert(abi.encodeWithSelector(OwnableUpgradeable.OwnableUnauthorizedAccount.selector, user3));
        reputationValReg.freeze(exampleConsAddr1);
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectRevert("Validator consensus address must be stored");
        reputationValReg.freeze(exampleConsAddr1);
        vm.stopPrank();

        vm.startPrank(user1);
        bytes[] memory consAddrs = new bytes[](1);
        consAddrs[0] = exampleConsAddr1;
        reputationValReg.storeConsAddrs(consAddrs);
        vm.stopPrank();

        vm.roll(5);
        assertEq(block.number, 5);
        assertTrue(reputationValReg.isEOAWhitelisted(user1));

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit EOAFrozen(user1, "bob");
        reputationValReg.freeze(exampleConsAddr1);
        vm.stopPrank();

        (ReputationValReg.State state, , uint256 freezeHeight, ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state), uint256(ReputationValReg.State.Frozen));
        assertEq(freezeHeight, 5);

        vm.startPrank(owner);
        vm.expectRevert("EOA representing validator must be active");
        reputationValReg.freeze(exampleConsAddr1);
        vm.stopPrank();
    }

    function testUnfreeze() public {
        testFreeze();

        vm.startPrank(user2);
        vm.expectRevert("Sender must be frozen");
        reputationValReg.unfreeze();
        vm.stopPrank();

        vm.startPrank(owner);
        vm.expectRevert("Sender must be frozen");
        reputationValReg.unfreeze();
        vm.stopPrank();

        vm.startPrank(user1);
        vm.expectRevert("Freeze period has not elapsed");
        reputationValReg.unfreeze();
        vm.stopPrank();

        vm.roll(95);
        vm.startPrank(user1);
        vm.expectRevert("Freeze period has not elapsed");
        reputationValReg.unfreeze();
        vm.stopPrank();

        vm.roll(110);
        vm.startPrank(user1);
        vm.expectRevert("Insufficient unfreeze fee");
        reputationValReg.unfreeze();
        vm.stopPrank();

        assertEq(reputationValReg.isEOAWhitelisted(user1), true); // frozen EOAs are still whitelisted
        (ReputationValReg.State state, uint256 numConsAddrsStored, 
            uint256 freezeHeight, string memory moniker) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state), uint256(ReputationValReg.State.Frozen));
        assertEq(numConsAddrsStored, 1);
        assertEq(freezeHeight, 5);
        assertEq(moniker, "bob");

        vm.deal(user1, 2 ether);
        vm.startPrank(user1);
        vm.expectEmit(true, true, true, true);
        emit EOAUnfrozen(user1, "bob");
        reputationValReg.unfreeze{value: 1 ether}();
        vm.stopPrank();

        (ReputationValReg.State state2, uint256 numConsAddrsStored2, 
            uint256 freezeHeight2, string memory moniker2) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state2), uint256(ReputationValReg.State.Active));
        assertEq(numConsAddrsStored2, 1);
        assertEq(freezeHeight2, 0);
        assertEq(moniker2, "bob");
    }

    function testFreezeCycle() public {
        testUnfreeze();

        vm.startPrank(owner);
        vm.expectEmit(true, true, true, true);
        emit EOAFrozen(user1, "bob");
        reputationValReg.freeze(exampleConsAddr1);
        vm.stopPrank();

        assertEq(reputationValReg.isEOAWhitelisted(user1), true); // frozen EOAs are still whitelisted
        (ReputationValReg.State state, uint256 numConsAddrsStored, 
            uint256 freezeHeight, string memory moniker) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(uint256(state), uint256(ReputationValReg.State.Frozen));
        assertEq(numConsAddrsStored, 1);
        assertEq(freezeHeight, 110);
        assertEq(moniker, "bob");
    }

    function testStoreConsAddrs() public {
        testAddWhitelistedEOA();
        
        bytes[] memory consAddrs = new bytes[](101);
        for (uint256 i = 0; i < 101; i++) {
            bytes memory consAddr = exampleConsAddr1;
            consAddr[consAddr.length - 1] = bytes1(uint8(i % 256));
            consAddrs[i] = consAddr;
        }
        vm.expectRevert("Too many cons addrs in request. Try batching");
        reputationValReg.storeConsAddrs(consAddrs);

        bytes[] memory consAddrsSmaller = new bytes[](1);
        consAddrsSmaller[0] = exampleConsAddr1;
        vm.prank(user3);
        vm.expectRevert("Sender must be whitelisted");
        reputationValReg.storeConsAddrs(consAddrsSmaller);
        vm.stopPrank();

        bytes[] memory consAddrsDuplicate = new bytes[](2);
        consAddrsDuplicate[0] = exampleConsAddr1;
        consAddrsDuplicate[1] = exampleConsAddr1;
        vm.prank(user1);
        vm.expectRevert("Duplicate consensus address is already stored");
        reputationValReg.storeConsAddrs(consAddrsDuplicate);
        vm.stopPrank();

        bytes[] memory consAddrsExceedingMax = new bytes[](MAX_CONS_ADDRS_PER_EOA + 1);
        for (uint256 i = 0; i < MAX_CONS_ADDRS_PER_EOA + 1; i++) {
            consAddrsExceedingMax[i] = exampleConsAddr1;
            consAddrsExceedingMax[i][consAddrsExceedingMax[i].length - 1] = bytes1(uint8(i % 256));
        }
        vm.prank(user1);
        vm.expectRevert("EOA must not store more than max allowed cons addrs");
        reputationValReg.storeConsAddrs(consAddrsExceedingMax);
        vm.stopPrank();

        bytes[] memory consAddrsValid = new bytes[](MAX_CONS_ADDRS_PER_EOA);
        for (uint256 i = 0; i < MAX_CONS_ADDRS_PER_EOA; i++) {
            consAddrsValid[i] = exampleConsAddr1;
            consAddrsValid[i][consAddrsValid[i].length - 1] = bytes1(uint8(i % 256));
        }

        (, uint256 numConsAddrsStored, , ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(numConsAddrsStored, 0);

        vm.prank(user1);
        vm.expectEmit(true, true, true, true);
        emit ConsAddrStored(consAddrsValid[0], user1, "bob");
        vm.expectEmit(true, true, true, true);
        emit ConsAddrStored(consAddrsValid[1], user1, "bob");
        vm.expectEmit(true, true, true, true);
        emit ConsAddrStored(consAddrsValid[2], user1, "bob");
        vm.expectEmit(true, true, true, true);
        emit ConsAddrStored(consAddrsValid[3], user1, "bob");
        vm.expectEmit(true, true, true, true);
        emit ConsAddrStored(consAddrsValid[4], user1, "bob");
        reputationValReg.storeConsAddrs(consAddrsValid);
        vm.stopPrank();

        (, uint256 numConsAddrsStored2, , ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(numConsAddrsStored2, 5);
    }

    function testDeleteConsAddrs() public {
        testStoreConsAddrs();

        bytes[] memory consAddrs = new bytes[](101);
        for (uint256 i = 0; i < 101; i++) {
            consAddrs[i] = exampleConsAddr1;
            consAddrs[i][consAddrs[i].length - 1] = bytes1(uint8(i % 256));
        }
        vm.expectRevert("Too many cons addrs in request. Try batching");
        reputationValReg.deleteConsAddrs(consAddrs);

        bytes[] memory consAddrsValid = new bytes[](MAX_CONS_ADDRS_PER_EOA);
        for (uint256 i = 0; i < MAX_CONS_ADDRS_PER_EOA; i++) {
            consAddrsValid[i] = exampleConsAddr1;
            consAddrsValid[i][consAddrsValid[i].length - 1] = bytes1(uint8(i % 256));
        }

        vm.prank(user3);
        vm.expectRevert("Sender must be whitelisted");
        reputationValReg.deleteConsAddrs(consAddrsValid);
        vm.stopPrank();

        bytes[] memory consAddrNotStored = new bytes[](1);
        consAddrNotStored[0] = exampleConsAddr2;
        vm.prank(user1);
        vm.expectRevert("Consensus address must be stored by sender");
        reputationValReg.deleteConsAddrs(consAddrNotStored);
        vm.stopPrank();

        (, uint256 numConsAddrsStored, , ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(numConsAddrsStored, 5);

        bytes[] memory consAddrSubset = new bytes[](3);
        for (uint256 i = 0; i < 3; i++) {
            consAddrSubset[i] = consAddrsValid[i];
        }
        vm.prank(user1);
        reputationValReg.deleteConsAddrs(consAddrSubset);
        vm.stopPrank();

        (, uint256 numConsAddrsStoredAfterDeletion, , ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(numConsAddrsStoredAfterDeletion, 2);

        bytes[] memory remainingConsAddrs = new bytes[](2);
        remainingConsAddrs[0] = consAddrsValid[3];
        remainingConsAddrs[1] = consAddrsValid[4];

        vm.prank(user1);
        reputationValReg.deleteConsAddrs(remainingConsAddrs);
        vm.stopPrank();

        (, uint256 numConsAddrsStoredAfterDeletion2, , ) = reputationValReg.getWhitelistedEOAInfo(user1);
        assertEq(numConsAddrsStoredAfterDeletion2, 0);
    }

    // TODO: test on add cons addr -> remove -> add again

    function testAreValidatorsOptedIn() public {
        testStoreConsAddrs();

        bytes[] memory consAddrsValid = new bytes[](MAX_CONS_ADDRS_PER_EOA);
        for (uint256 i = 0; i < MAX_CONS_ADDRS_PER_EOA; i++) {
            consAddrsValid[i] = exampleConsAddr1;
            consAddrsValid[i][consAddrsValid[i].length - 1] = bytes1(uint8(i % 256));
        }
        bool[] memory optedIn = reputationValReg.areValidatorsOptedIn(consAddrsValid);
        for (uint256 i = 0; i < 5; i++) {
            assertTrue(optedIn[i]);
        }

        bytes[] memory consAddrsSubset = new bytes[](3);
        for (uint256 i = 0; i < 3; i++) {
            consAddrsSubset[i] = consAddrsValid[i];
        }
        bool[] memory optedInSubset = reputationValReg.areValidatorsOptedIn(consAddrsSubset);
        for (uint256 i = 0; i < 3; i++) {
            assertTrue(optedInSubset[i]);
        }
        
        vm.prank(owner);
        reputationValReg.freeze(consAddrsValid[0]);
        vm.stopPrank();
        optedIn = reputationValReg.areValidatorsOptedIn(consAddrsValid);
        for (uint256 i = 0; i < MAX_CONS_ADDRS_PER_EOA; i++) {
            assertFalse(optedIn[i]);
        }
    }
    
    



}
