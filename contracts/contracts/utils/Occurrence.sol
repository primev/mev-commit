// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

library OccurrenceLib {
    struct BlockHeightOccurrence {
        bool exists;
        uint256 blockHeight;
    }

    struct TimestampOccurrence {
        bool exists;
        uint256 timestamp;
    }

    function captureOccurrence(BlockHeightOccurrence storage self) internal {
        self.exists = true;
        self.blockHeight = block.number;
    }

    function del(BlockHeightOccurrence storage self) internal {
        self.exists = false;
        self.blockHeight = 0;
    }

    function captureOccurrence(TimestampOccurrence storage self) internal {
        self.exists = true;
        self.timestamp = block.timestamp;
    }

    function del(TimestampOccurrence storage self) internal {
        self.exists = false;
        self.timestamp = 0;
    }
}
