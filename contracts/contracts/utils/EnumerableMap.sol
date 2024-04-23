// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

// TODO: talk about open zep and lack of generic types

import {EnumerableSet} from "./EnumerableSet.sol";

library EnumerableMap {
    using EnumerableSet for EnumerableSet.BytesSet;

    error EnumerableMapNonexistentKey(bytes key);

    struct BytesToUint256Map {
        EnumerableSet.BytesSet _keys;
        mapping(bytes key => uint256) _values;
    }

    /**
     * @dev Adds a key-value pair to a map, or updates the value for an existing
     * key. O(1).
     *
     * Returns true if the key was added to the map, that is if it was not
     * already present.
     */
    function set(BytesToUint256Map storage map, bytes memory key, uint256 value) internal returns (bool) {
        map._values[key] = value;
        return map._keys.add(key);
    }

    /**
     * @dev Removes a key-value pair from a map. O(1).
     *
     * Returns true if the key was removed from the map, that is if it was present.
     */
    function remove(BytesToUint256Map storage map, bytes memory key) internal returns (bool) {
        delete map._values[key];
        return map._keys.remove(key);
    }

    /**
     * @dev Returns true if the key is in the map. O(1).
     */
    function contains(BytesToUint256Map storage map, bytes memory key) internal view returns (bool) {
        return map._keys.contains(key);
    }

    /**
     * @dev Returns the number of key-value pairs in the map. O(1).
     */
    function length(BytesToUint256Map storage map) internal view returns (uint256) {
        return map._keys.length();
    }

    /**
     * @dev Returns the key-value pair stored at position `index` in the map. O(1).
     *
     * Note that there are no guarantees on the ordering of entries inside the
     * array, and it may change when more entries are added or removed.
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function at(BytesToUint256Map storage map, uint256 index) internal view returns (bytes memory, uint256) {
        bytes storage key = map._keys.at(index);
        return (key, map._values[key]);
    }

    /**
     * @dev Tries to returns the value associated with `key`. O(1).
     * Does not revert if `key` is not in the map.
     */
    function tryGet(BytesToUint256Map storage map, bytes memory key) internal view returns (bool, uint256) {
        uint256 value = map._values[key];
        if (value == uint256(0)) {
            return (contains(map, key), uint256(0));
        } else {
            return (true, value);
        }
    }

    /**
     * @dev Returns the value associated with `key`. O(1).
     *
     * Requirements:
     *
     * - `key` must be in the map.
     */
    function get(BytesToUint256Map storage map, bytes memory key) internal view returns (uint256) {
        uint256 value = map._values[key];
        if (value == 0 && !contains(map, key)) {
            revert EnumerableMapNonexistentKey(key);
        }
        return value;
    }

    /**
     * @dev Return the an array containing all the keys
     *
     * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
     * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
     * this function has an unbounded cost, and using it as part of a state-changing function may render the function
     * uncallable if the map grows to a point where copying to memory consumes too much gas to fit in a block.
     */
    function keys(BytesToUint256Map storage map) internal view returns (bytes[] memory) {
        return map._keys.values();
    }

    /**
    * @dev Return a subset of keys
     */
    function keySubset(BytesToUint256Map storage map, uint256 start, uint256 end) internal view returns (bytes[] memory) {
        return map._keys.valueSubset(start, end);
    }
}
