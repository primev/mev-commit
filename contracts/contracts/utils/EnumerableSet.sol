// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/// @title Implements an enumerable set of bytes arrays.
/// @notice Adapted from OpenZeppelin's EnumerableSet.sol implementation. 
/// See https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.0/contracts/utils/structs/EnumerableSet.sol
///
/// The openzeppelin EnumerableSet implementations are only compatible with values that are 32 bytes.
/// Hence we were required to alter their source code to enable 48 byte BLS pubkey storage.
/// This implementation is streamlined to only support a set of "bytes" type.
library EnumerableSet {

    // Represents a set of byte array values
    struct BytesSet {
        bytes[] _values;
        mapping(bytes value => uint256) _positions;
    }

    error StartMustBeLessThanEnd();
    error EndTooLarge();

    /**
     * @dev Add a value to a set. O(1).
     *
     * Returns true if the value was added to the set, that is if it was not
     * already present.
     */
    function add(BytesSet storage set, bytes memory value) internal returns (bool) {
        return _add(set, value);
    }

    /**
     * @dev Removes a value from a set. O(1).
     *
     * Returns true if the value was removed from the set, that is if it was
     * present.
     */
    function remove(BytesSet storage set, bytes memory value) internal returns (bool) {
        return _remove(set, value);
    }

    /**
     * @dev Swaps a value with another value at a target (1-indexed) position. O(1).
     *
     * Returns true if the swap was successful.
     */
    function swapWithPosition(BytesSet storage set, bytes memory originalValue, uint256 targetPosition) internal returns (bool) {
        return _swapWithPosition(set, originalValue, targetPosition);
    }

    /**
     * @dev Returns true if the value is in the set. O(1).
     */
    function contains(BytesSet storage set, bytes memory value) internal view returns (bool) {
        return _contains(set, value);
    }

    /**
     * @dev Returns the position of the value in the set. O(1).
     */
    function position(BytesSet storage set, bytes memory value) internal view returns (uint256) {
        return _position(set, value);
    }

    /**
     * @dev Returns the number of values in the set. O(1).
     */
    function length(BytesSet storage set) internal view returns (uint256) {
        return _length(set);
    }

    /**
     * @dev Returns the value stored at position `index` in the set. O(1).
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function at(BytesSet storage set, uint256 index) internal view returns (bytes storage) {
        return _at(set, index);
    }

    /**
     * @dev Return the entire set in an array
     *
     * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
     * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
     * this function has an unbounded cost, and using it as part of a state-changing function may render the function
     * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
     */
    function values(BytesSet storage set) internal view returns (bytes[] memory) {
        bytes[] memory store = _values(set);
        bytes[] memory result;

        /// @solidity memory-safe-assembly
        // solhint-disable-next-line no-inline-assembly
        assembly {
            result := store
        }

        return result;
    }

    /**
     * @dev Returns an end-exclusive subset of the values in the set. O(end - start).
     *
     * Requirements:
     *
     * - `start` must be less than `end`.
     * - `end` must be less than or equal to the length of the set.
     */
    function valueSubset(BytesSet storage set, uint256 start, uint256 end) internal view returns (bytes[] memory) {
        require(start < end, StartMustBeLessThanEnd());
        require(end <= set._values.length, EndTooLarge());

        bytes[] memory result = new bytes[](end - start);

        for (uint256 i = start; i < end; ++i) {
            result[i - start] = set._values[i];
        }

        return result;
    }

    /**
     * @dev Add a value to a set. O(1).
     *
     * Returns true if the value was added to the set, that is if it was not
     * already present.
     */
    function _add(BytesSet storage set, bytes memory value) private returns (bool) {
        if (!_contains(set, value)) {
            set._values.push(value);
            // The value is stored at length-1, but we add 1 to all indexes
            // and use 0 as a sentinel value
            set._positions[value] = set._values.length;
            return true;
        } else {
            return false;
        }
    }

    /**
     * @dev Removes a value from a set. O(1).
     *
     * Returns true if the value was removed from the set, that is if it was
     * present.
     */
    function _remove(BytesSet storage set, bytes memory value) private returns (bool) {
        // We cache the value's position to prevent multiple reads from the same storage slot
        uint256 pos = set._positions[value];

        if (pos != 0) {
            // Equivalent to contains(set, value)
            // To delete an element from the _values array in O(1), we swap the element to delete with the last one in
            // the array, and then remove the last element (sometimes called as 'swap and pop').
            // This modifies the order of the array, as noted in {at}.

            uint256 valueIndex = pos - 1;
            uint256 lastIndex = set._values.length - 1;

            if (valueIndex != lastIndex) {
                bytes storage lastValue = set._values[lastIndex];

                // Move the lastValue to the index where the value to delete is
                set._values[valueIndex] = lastValue;
                // Update the tracked position of the lastValue (that was just moved)
                set._positions[lastValue] = pos;
            }

            // Delete the slot where the moved value was stored
            set._values.pop();

            // Delete the tracked position for the deleted slot
            delete set._positions[value];

            return true;
        } else {
            return false;
        }
    }

    /**
     * @dev Swaps any value with the value at the specified (1-indexed) position in the array.
     * @return boolean indicating if the swap was successful
     */
    function _swapWithPosition(BytesSet storage set, bytes memory originalValue, uint256 targetPosition) private returns (bool) {
        uint256 len = set._values.length;

        if (targetPosition == 0 || targetPosition > len) {
            return false;
        }

        uint256 originalPosition = set._positions[originalValue];
        if (originalPosition == 0) {
            return false;
        }

        if (originalPosition == targetPosition) {
            return true;
        }

        uint256 originalIndex = originalPosition - 1;
        uint256 targetIndex = targetPosition - 1;

        bytes memory targetValue = set._values[targetIndex];

        set._values[originalIndex] = targetValue;
        set._values[targetIndex] = originalValue;

        set._positions[targetValue] = originalPosition;
        set._positions[originalValue] = targetPosition;

        return true;
    }

    /**
     * @dev Returns true if the value is in the set. O(1).
     */
    function _contains(BytesSet storage set, bytes memory value) private view returns (bool) {
        return set._positions[value] != 0;
    }

    /**
     * @dev Returns the one-indexed position of the value in the set. O(1).
     */
    function _position(BytesSet storage set, bytes memory value) private view returns (uint256) {
        return set._positions[value];
    }

    /**
     * @dev Returns the number of values on the set. O(1).
     */
    function _length(BytesSet storage set) private view returns (uint256) {
        return set._values.length;
    }

    /**
     * @dev Returns the value stored at position `index` in the set. O(1).
     *
     * Requirements:
     *
     * - `index` must be strictly less than {length}.
     */
    function _at(BytesSet storage set, uint256 index) private view returns (bytes storage) {
        return set._values[index];
    }

    /**
     * @dev Return the entire set in an array
     *
     * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
     * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
     * this function has an unbounded cost, and using it as part of a state-changing function may render the function
     * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
     */
    function _values(BytesSet storage set) private view returns (bytes[] memory) {
        return set._values;
    }
}
