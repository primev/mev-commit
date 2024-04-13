pragma solidity ^0.8.19;

contract Timestamp {
    uint256[] public blockTimestamps;

    function getBlockTimestamp() public view returns (uint256) {
        return block.timestamp;
    }

    function storeBlockTimestamp() public {
        blockTimestamps.push(block.timestamp);
    }

    function getBlockTimestampAtIndex(uint256 index) public view returns (uint256) {
        require(index < blockTimestamps.length, "Index out of bounds");
        return blockTimestamps[index];
    }
}

