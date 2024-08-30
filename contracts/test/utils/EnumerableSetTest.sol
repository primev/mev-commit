// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import "../../contracts/utils/EnumerableSet.sol";

contract EnumerableSetTest is Test {

    using EnumerableSet for EnumerableSet.BytesSet;
    EnumerableSet.BytesSet private set;

    function setUp() public { 
        for (uint256 i = 0; i < set.length(); i++) {
            set.remove(set.at(i));
        }
        require(set.length() == 0, "Set should be empty.");
        require(set.values().length == 0, "Set should have no values.");
    }

    function testAddElement() public {
        bytes memory tempData = "data";
        bool added = set.add(tempData);
        assertTrue(added, "Element should be added.");
        assertTrue(set.contains(tempData), "Element should be in the set.");
    }

    function testAddDuplicateElement() public {
        bytes memory tempData = "data";
        bool added = set.add(tempData);
        assertTrue(added, "Element should be added.");
        bool addedAgain = set.add(tempData);
        assertFalse(addedAgain, "Duplicate element should not be added again.");
    }

    function testRemoveElement() public {
        bytes memory tempData = "data";
        bool added = set.add(tempData);
        assertTrue(added, "Element should be added.");
        assertTrue(set.contains(tempData), "Element should be in the set before removal.");

        bool removed = set.remove(tempData);
        assertTrue(removed, "Element should be removed.");
        assertFalse(set.contains(tempData), "Element should not be in the set after removal.");
    }

    function testRemoveNonexistentElement() public {
        bytes memory tempData = "nonexistent";
        assertEq(set.length(), 0, "Set should be empty.");
        bool removed = set.remove(tempData);
        assertFalse(removed, "Nonexistent element should not be removed.");
        assertEq(set.length(), 0, "Set should still be empty.");

        set.add("data");
        assertEq(set.length(), 1, "Set should have 1 element.");
        assertTrue(set.contains("data"), "Element should be in the set.");
        bool removed2 = set.remove(tempData);
        assertFalse(removed2, "Nonexistent element should not be removed.");
        assertEq(set.length(), 1, "Set should still have 1 element.");
    }

    function testAt() public {
        set.add("data");
        assertEq(set.at(0), "data", "Element at index 0 should be the added element.");

        set.add("data2");
        set.add("data3");
        set.add("data4");

        assertEq(set.at(0), "data", "Element at index 0 should be the first added element.");
        assertEq(set.at(1), "data2", "Element at index 1 should be the second added element.");
        assertEq(set.at(2), "data3", "Element at index 2 should be the third added element.");
        assertEq(set.at(3), "data4", "Element at index 3 should be the fourth added element.");
    }

    function testGetAllValues() public {
        bytes[] memory values = new bytes[](3);
        values[0] = "data1";
        values[1] = "data2";
        values[2] = "data3";
        for (uint256 i = 0; i < values.length; i++) {
            set.add(values[i]);
        }

        assertTrue(set.contains("data1"), "Element 1 should be in the set.");
        assertTrue(set.contains("data2"), "Element 2 should be in the set.");
        assertTrue(set.contains("data3"), "Element 3 should be in the set.");

        set.remove("data2");
        assertFalse(set.contains("data2"), "Element 2 should not be in the set after removal.");
        assertTrue(set.contains("data1"), "Element 1 should still be in the set.");
        assertTrue(set.contains("data3"), "Element 3 should still be in the set.");

        bytes[] memory queriedValues = set.values();
        assertEq(queriedValues.length, 2, "There should be 2 elements in the set.");
        assertEq(queriedValues[0], "data1", "Element 1 should be in the set.");
        assertEq(queriedValues[1], "data3", "Element 3 should be in the set.");
    }

    function testGetValueBatch() public {
        bytes[] memory values = new bytes[](3);
        values[0] = "data1";
        values[1] = "data2";
        values[2] = "data3";

        for (uint256 i = 0; i < values.length; i++) {
            set.add(values[i]);
        }

        assertTrue(set.contains("data1"), "Element 1 should be in the set.");
        assertTrue(set.contains("data2"), "Element 2 should be in the set.");
        assertTrue(set.contains("data3"), "Element 3 should be in the set.");

        set.remove("data1");
        assertFalse(set.contains("data1"), "Element 1 should not be in the set after removal.");
        assertTrue(set.contains("data2"), "Element 2 should still be in the set.");
        assertTrue(set.contains("data3"), "Element 3 should still be in the set.");

        bytes[] memory allQuery = set.values();
        assertEq(allQuery.length, 2, "There should be 2 elements in the full set");

        // Assert expected values without asserting order 
        bytes[] memory expected = new bytes[](2);
        expected[0] = "data2";
        expected[1] = "data3";
        assertContainsAll(allQuery, expected);

        bytes[] memory batchQuery = set.valueSubset(0, 1);
        assertEq(batchQuery.length, 1, "There should be 1 element in the batch.");
        // require element is either data2 or data3
        require(keccak256(abi.encodePacked(batchQuery[0])) == keccak256(abi.encodePacked("data2")) || keccak256(abi.encodePacked(batchQuery[0])) == keccak256(abi.encodePacked("data3")), "Element 2 or 3 should be in the batch.");

        bytes[] memory batchQuery2 = set.valueSubset(1, 2);
        assertEq(batchQuery2.length, 1, "There should be 1 element in the batch.");
        require(keccak256(abi.encodePacked(batchQuery2[0])) == keccak256(abi.encodePacked("data2")) || keccak256(abi.encodePacked(batchQuery2[0])) == keccak256(abi.encodePacked("data3")), "Element 2 or 3 should be in the batch.");

        set.add("data4");
        set.add("data5");
        set.add("data6");
        set.add("data7");

        bytes[] memory lastThree = set.valueSubset(3, 6);
        assertEq(lastThree.length, 3, "There should be 3 elements in the batch.");

        bytes[] memory firstThree = set.valueSubset(0, 3);
        assertEq(firstThree.length, 3, "There should be 3 elements in the batch.");

        bytes[] memory combo = new bytes[](6);
        for (uint256 i = 0; i < 3; i++) {
            combo[i] = firstThree[i];
        }
        for (uint256 i = 0; i < 3; i++) {
            combo[i + 3] = lastThree[i];
        }
        bytes[] memory expectedCombo = new bytes[](6);
        expectedCombo[0] = "data7";
        expectedCombo[1] = "data6";
        expectedCombo[2] = "data5";
        expectedCombo[3] = "data4";
        expectedCombo[4] = "data3";
        expectedCombo[5] = "data2";
        assertContainsAll(combo, expectedCombo);
    }

    // Asserts an array contains all expected elements without asserting order
    function assertContainsAll(bytes[] memory actual, bytes[] memory expected) internal pure {
        require(actual.length == expected.length, "Arrays have different lengths");

        uint foundCount = 0;

        for (uint i = 0; i < expected.length; i++) {
            for (uint j = 0; j < actual.length; j++) {
                if (keccak256(abi.encodePacked(actual[j])) == keccak256(abi.encodePacked(expected[i]))) {
                    foundCount++;
                    break;
                }
            }
        }

        require(foundCount == expected.length, "Not all expected elements found");
    }
}
