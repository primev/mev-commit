// SPDX-License-Identifier: BSL 1.1

// solhint-disable one-contract-per-file
pragma solidity 0.8.29;

library BlockHeightOccurrence {
    struct Occurrence {
        bool exists;
        uint256 blockHeight;
    }

    function captureOccurrence(Occurrence storage self) internal {
        self.exists = true;
        self.blockHeight = block.number;
    }

    function del(Occurrence storage self) internal {
        self.exists = false;
        self.blockHeight = 0;
    }
}

library TimestampOccurrence {
    struct Occurrence {
        bool exists;
        uint256 timestamp;
    }

    function captureOccurrence(Occurrence storage self) internal {
        self.exists = true;
        self.timestamp = block.timestamp;
    }

    function del(Occurrence storage self) internal {
        self.exists = false;
        self.timestamp = 0;
    }
}
