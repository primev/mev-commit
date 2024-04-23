// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "../../contracts/utils/EnumerableSet.sol";

contract EnumerableSetTest is Test {
    using EnumerableSet for EnumerableSet.BytesSet;

    EnumerableSet.BytesSet private set;
    bytes dataElement = "data";

    function setUp() public { }

    function testAddElement() public {
        bool added = set.add(dataElement);
        assertTrue(added, "Element should be added.");
        assertTrue(set.contains(dataElement), "Element should be in the set.");
    }

    function testAddDuplicateElement() public {
        set.add(dataElement);
        bool addedAgain = set.add(dataElement);
        assertFalse(addedAgain, "Duplicate element should not be added again.");
    }

    function testRemoveElement() public {
        set.add(dataElement);
        assertTrue(set.contains(dataElement), "Element should be in the set before removal.");

        bool removed = set.remove(dataElement);
        assertTrue(removed, "Element should be removed.");
        assertFalse(set.contains(dataElement), "Element should not be in the set after removal.");
    }

    // function testRemoveNonexistentElement() public {
    //     bytes memory dataNonexistent = "nonexistent";
    //     bool removed = set.remove(dataNonexistent);
    //     assertFalse(removed, "Nonexistent element should not be removed.");
    // }

    // // Test counting elements in the set
    // function testCountElements() public {
    //     bytes memory dataOne = "one";
    //     bytes memory dataTwo = "two";
    //     set.add(dataOne);
    //     set.add(dataTwo);
    //     uint256 count = set.length();
    //     assertEq(count, 2, "There should be two elements in the set.");
    // }

    // // Test retrieving elements by index
    // function testGetElementByIndex() public {
    //     bytes memory dataIndexed = "indexed";
    //     set.add(dataIndexed);
    //     bytes storage retrievedElement = set.at(0);
    //     assertEq(retrievedElement, dataIndexed, "Element retrieved by index should match the added element.");
    // }

    // // Test that listing all values works as expected
    // function testListAllValues() public {
    //     bytes memory dataFirst = "first";
    //     bytes memory dataSecond = "second";
    //     set.add(dataFirst);
    //     set.add(dataSecond);

    //     bytes[] memory allValues = set.values();
    //     assertEq(allValues.length, 2, "All values should be listed.");
    //     assertEq(allValues[0], dataFirst, "First element should match.");
    //     assertEq(allValues[1], dataSecond, "Second element should match.");
    // }
}
