// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

import {Gateway} from "./Gateway.sol";

contract L1Gateway is Gateway {

    constructor(address _owner, address _relayer, uint256 _finalizationFee, uint256 _counterpartyFee
        ) Gateway(_owner, _relayer, _finalizationFee, _counterpartyFee) {}

    function _decrementMsgSender(uint256 _amount) internal override {
        require(msg.value == _amount, "Incorrect Ether value sent");
        // Wrapping function initiateTransfer is payable. Ether is escrowed in contract balance
    }

    function _fund(uint256 _amount, address _toFund) internal override {
        require(address(this).balance >= _amount, "Insufficient contract balance");
        payable(_toFund).transfer(_amount);
    }

    receive() external payable {}
}

