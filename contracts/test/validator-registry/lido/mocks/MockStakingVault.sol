// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

contract MockStakingVault {
    address public nodeOperator_;

    constructor(address _nodeOperator) {
        nodeOperator_ = _nodeOperator;
    }

    function setNodeOperator(address n) external { nodeOperator_ = n; }

    // ---- IStakingVault ----
    function nodeOperator() external view returns (address) {
        return nodeOperator_;
    }

    // allow receiving ETH (not required, but handy)
    receive() external payable {}
}
