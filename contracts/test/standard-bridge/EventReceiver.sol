// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract EventReceiver {
    event Received(address indexed sender, uint256 amount);

    receive() external payable {
        for (uint256 i = 0; i < 5; i++) {
            emit Received(msg.sender, msg.value);
        }
    }
}
