// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract RevertingReceiver {
    bool internal shouldRevert;

    function setShouldRevert(bool _shouldRevert) external {
        shouldRevert = _shouldRevert;
    }

    receive() external payable {
        if (shouldRevert) revert("RevertingReceiver: Revert on receiving Ether");
    }
}
