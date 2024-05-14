// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../contracts/utils/EnumerableMap.sol";

contract EnumerableMapTest is Test {
    using EnumerableMap for EnumerableMap.BytesToUint256Map;

    EnumerableMap.BytesToUint256Map private map;

    function setUp() public {
        while (map.length() > 0) {
            (bytes memory key,) = map.at(map.length() - 1);
            map.remove(key);
        }
        require(map.length() == 0, "Map should be empty after setup.");
        require(map.keys().length == 0, "Map should have no keys after setup.");
    }

    function testSetAndGet() public {
        bool added = map.set("key1", 123);
        assertTrue(added, "Key should be added.");
        assertTrue(map.contains("key1"), "Key should exist in the map.");

        uint256 retrievedValue = map.get("key1");
        assertEq(retrievedValue, 123, "Retrieved value should match the set value.");
    }

    function testUpdateValue() public {
        bool added = map.set("key1", 123);
        assertTrue(added, "Key should be added.");
        assertEq(map.get("key1"), 123, "Initial value should be set.");

        bool updated = map.set("key1", 456);
        assertFalse(updated, "Key should not be added again, should be updated.");

        uint256 retrievedValue = map.get("key1");
        assertEq(retrievedValue, 456, "Retrieved value should match the updated value.");
    }

    function testRemoveKey() public {
        map.set("key1", 123);
        assertTrue(map.contains("key1"), "Key should exist before removal.");

        bool removed = map.remove("key1");
        assertTrue(removed, "Key should be removed.");
        assertFalse(map.contains("key1"), "Key should no longer exist in the map.");
    }

    function testNonexistentKey() public {
        bytes memory key = bytes("nonexistent");
        bool exists;
        uint256 value;

        (exists, value) = map.tryGet(key);
        assertFalse(exists, "Key should not exist.");

        bool removed = map.remove(key);
        assertFalse(removed, "Key should not be removed.");
    }

    function testGetAllKeys() public {
        map.set("key1", 123);
        map.set("key2", 456);
        map.set("key3", 789);
        map.set("key4", 101112);
        map.set("key5", 131415);

        bytes[] memory keys = map.keys();

        bytes[] memory expected = new bytes[](5);
        expected[0] = "key5";
        expected[1] = "key4";
        expected[2] = "key3";
        expected[3] = "key2";
        expected[4] = "key1";
        assertContainsAll(keys, expected);

        map.remove("key1");
        map.remove("key3");

        keys = map.keys();

        expected = new bytes[](3);
        expected[0] = "key5";
        expected[1] = "key4";
        expected[2] = "key2"; 
        assertContainsAll(keys, expected);
    }

    function testGetKeySubset() public {
        map.set("key1", 123);
        map.set("key2", 456);
        map.set("key3", 789);
        map.set("key4", 101);

        bytes[] memory keys = map.keys();
        assertEq(keys.length, 4, "Map should have 4 keys.");
        
        bytes[] memory firstThree = map.keySubset(0, 3);
        assertEq(firstThree.length, 3, "Subset should have 3 keys.");

        bytes[] memory lastOne = map.keySubset(3, 4);
        assertEq(lastOne.length, 1, "Subset should have 1 key.");

        bytes[] memory combo = new bytes[](4);
        combo[0] = lastOne[0];
        combo[1] = firstThree[0];
        combo[2] = firstThree[1];
        combo[3] = firstThree[2];

        bytes[] memory expectedCombo = new bytes[](4);
        expectedCombo[0] = "key4";
        expectedCombo[1] = "key1";
        expectedCombo[2] = "key2";
        expectedCombo[3] = "key3";

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
