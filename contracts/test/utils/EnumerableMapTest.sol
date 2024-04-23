// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

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
        bytes memory key = bytes("key1");
        uint256 value = 123;

        // Test setting a new key
        bool added = map.set(key, value);
        assertTrue(added, "Key should be added.");
        assertTrue(map.contains(key), "Key should exist in the map.");

        // Test retrieving a value
        uint256 retrievedValue = map.get(key);
        assertEq(retrievedValue, value, "Retrieved value should match the set value.");
    }

    function testUpdateValue() public {
        bytes memory key = bytes("key1");
        uint256 initialValue = 123;
        uint256 updatedValue = 456;

        map.set(key, initialValue);
        bool updated = map.set(key, updatedValue);
        assertFalse(updated, "Key should not be added again, should be updated.");

        uint256 retrievedValue = map.get(key);
        assertEq(retrievedValue, updatedValue, "Retrieved value should match the updated value.");
    }

    function testRemoveKey() public {
        bytes memory key = bytes("key1");
        map.set(key, 123);

        assertTrue(map.contains(key), "Key should exist before removal.");

        bool removed = map.remove(key);
        assertTrue(removed, "Key should be removed.");
        assertFalse(map.contains(key), "Key should no longer exist in the map.");
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

    function testIteration() public {
        bytes memory key1 = bytes("key1");
        bytes memory key2 = bytes("key2");
        map.set(key1, 100);
        map.set(key2, 200);

        assertEq(map.length(), 2, "Map should contain two entries.");

        (bytes memory iterKey1, uint256 iterValue1) = map.at(0);
        assertTrue(iterKey1.length > 0, "First key should be non-empty.");
        assertTrue(iterValue1 > 0, "First value should be non-zero.");

        (bytes memory iterKey2, uint256 iterValue2) = map.at(1);
        assertTrue(iterKey2.length > 0, "Second key should be non-empty.");
        assertTrue(iterValue2 > 0, "Second value should be non-zero.");
    }

    function testKeyArray() public {
        map.set(bytes("key1"), 100);
        map.set(bytes("key2"), 200);
        map.set(bytes("key3"), 300);

        bytes[] memory keys = map.keys();
        assertEq(keys.length, 3, "There should be three keys in the map.");
    }
}
