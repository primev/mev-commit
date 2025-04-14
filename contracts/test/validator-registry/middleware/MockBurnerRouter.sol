// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

contract MockBurnerRouter {
    uint256 public delay;
    mapping(address network => address receiver) public networkReceiver;
    mapping(address network => mapping(address operator => address receiver)) public operatorNetworkReceiver;

    constructor(uint256 delay_) {
        delay = delay_;
    }

    function setDelay(uint256 delay_) external {
        delay = delay_;
    }

    function setNetworkReceiver(address network, address receiver) external {
        networkReceiver[network] = receiver;
    }

    function setOperatorNetworkReceiver(address network, address operator, address receiver) external {
        operatorNetworkReceiver[network][operator] = receiver;
    }
}
