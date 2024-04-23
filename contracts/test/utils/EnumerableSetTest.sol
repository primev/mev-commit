// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import "forge-std/Test.sol";
import "../../contracts/utils/EnumerableSet.sol";

contract EnumerableSetTest is Test {
    using EnumerableSet for EnumerableSet.BytesSet;

    EnumerableSet.BytesSet private set;

    // bytes storage variables are not used since we need to pass calldata directly

    function setUp() public { 
        // Optional: Initialize anything here if needed
    }

    function testAddElement() public {
        bytes memory tempData = "data"; // Use memory for temporary storage
        bool added = set.add(tempData); // Explicit conversion to calldata occurs here
        assertTrue(added, "Element should be added.");
        assertTrue(set.contains(tempData), "Element should be in the set.");
    }

    function testAddDuplicateElement() public {
        bytes memory tempData = "data"; // Reuse the memory variable
        set.add(tempData);  // First addition
        bool addedAgain = set.add(tempData);  // Attempt to add the same element again
        assertFalse(addedAgain, "Duplicate element should not be added again.");
    }

    function testRemoveElement() public {
        bytes memory tempData = "data";
        set.add(tempData);  // Ensure the element is added first
        assertTrue(set.contains(tempData), "Element should be in the set before removal.");

        bool removed = set.remove(tempData);  // Try to remove the element
        assertTrue(removed, "Element should be removed.");
        assertFalse(set.contains(tempData), "Element should not be in the set after removal.");
    }

    function testRemoveNonexistentElement() public {
        bytes memory tempData = "nonexistent";  // Use memory for temporary data
        bool removed = set.remove(tempData);  // Try to remove a non-existent element
        assertFalse(removed, "Nonexistent element should not be removed.");
    }
    
    // TODO: mas CRUD tests

    function testGetAllValues() public {
        bytes[] memory values = new bytes[](3);
        values[0] = "data1";
        values[1] = "data2";
        values[2] = "data3";

        for (uint256 i = 0; i < values.length; i++) {
            set.add(values[i]);
        }

        assertTrue(set.contains(values[0]), "Element 1 should be in the set.");
        assertTrue(set.contains(values[1]), "Element 2 should be in the set.");
        assertTrue(set.contains(values[2]), "Element 3 should be in the set.");

        set.remove(values[0]);
        assertFalse(set.contains(values[0]), "Element 1 should not be in the set after removal.");
        assertTrue(set.contains(values[1]), "Element 2 should still be in the set.");
        assertTrue(set.contains(values[2]), "Element 3 should still be in the set.");

        bytes[] memory setValues = set.values();
        assertEq(setValues.length, 2, "There should be 2 elements in the set.");
        assertTrue(set.contains(setValues[0]), "Element 2 should be in the set.");
        assertTrue(set.contains(setValues[1]), "Element 3 should be in the set.");
    }

    function testGetValueBatch() public {
        bytes[] memory values = new bytes[](3);
        values[0] = "data1";
        values[1] = "data2";
        values[2] = "data3";

        for (uint256 i = 0; i < values.length; i++) {
            set.add(values[i]);
        }

        assertTrue(set.contains(values[0]), "Element 1 should be in the set.");
        assertTrue(set.contains(values[1]), "Element 2 should be in the set.");
        assertTrue(set.contains(values[2]), "Element 3 should be in the set.");

        set.remove(values[0]);
        assertFalse(set.contains(values[0]), "Element 1 should not be in the set after removal.");
        assertTrue(set.contains(values[1]), "Element 2 should still be in the set.");
        assertTrue(set.contains(values[2]), "Element 3 should still be in the set.");

        bytes[] memory setValues = set.values(0, 1);
        assertEq(setValues.length, 1, "There should be 1 element in the set.");
        assertTrue(set.contains(setValues[0]), "Element 2 should be in the set.");

        setValues = set.values(1, 2);
        assertEq(setValues.length, 1, "There should be 1 element in the set.");
        assertTrue(set.contains(setValues[0]), "Element 3 should be in the set.");
    }

}
