// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

library BN128 {
    /// @dev Error if bn128 addition failed during zk proof validation
    error BN128AddFailed();

    /// @dev Error if bn128 multiplication failed during zk proof validation
    error BN128MulFailed();

    /**
     * @dev BN128 addition precompile call:
     *       (x3, y3) = (x1, y1) + (x2, y2)
     */
    function _ecAdd(
        uint256 x1,
        uint256 y1,
        uint256 x2,
        uint256 y2
    ) internal view returns (uint256 x3, uint256 y3) {
        // 0x06 = bn128Add precompile
        // Inputs are 4 * 32 bytes = x1, y1, x2, y2
        // Output is 2 * 32 bytes = (x3, y3)
        bool success;
        // solhint-disable-next-line no-inline-assembly
        assembly {
            // free memory pointer
            let memPtr := mload(0x40)
            mstore(memPtr, x1)
            mstore(add(memPtr, 0x20), y1)
            mstore(add(memPtr, 0x40), x2)
            mstore(add(memPtr, 0x60), y2)
            // call precompile
            if iszero(staticcall(gas(), 0x06, memPtr, 0x80, memPtr, 0x40)) {
                revert(0, 0)
            }
            x3 := mload(memPtr)
            y3 := mload(add(memPtr, 0x20))
            success := true
        }
        require(success, BN128AddFailed());
    }

    /**
     * @dev BN128 multiplication precompile call:
     *       (x3, y3) = scalar * (x1, y1)
     */
    function _ecMul(
        uint256 x1,
        uint256 y1,
        uint256 scalar
    ) internal view returns (uint256 x2, uint256 y2) {
        bool success;
        // solhint-disable-next-line no-inline-assembly
        assembly {
            let memPtr := mload(0x40)
            mstore(memPtr, x1)
            mstore(add(memPtr, 0x20), y1)
            mstore(add(memPtr, 0x40), scalar)
            // call precompile at 0x07
            if iszero(staticcall(gas(), 0x07, memPtr, 0x60, memPtr, 0x40)) {
                revert(0, 0)
            }
            x2 := mload(memPtr)
            y2 := mload(add(memPtr, 0x20))
            success := true
        }
        require(success, BN128MulFailed());
    }
}
