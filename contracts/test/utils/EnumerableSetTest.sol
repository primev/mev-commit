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

    // Other test methods would similarly convert storage to memory for calldata compatibility...
}
