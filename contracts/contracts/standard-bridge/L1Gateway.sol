// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Gateway} from "./Gateway.sol";

/// @title L1Gateway
/// @notice Gateway contract deployed on L1 enabling the mev-commit standard bridge.
/// @dev This contract will escrow locked ETH, while a corresponding amount is minted from the SettlementGateway on the mev-commit chain.
contract L1Gateway is Gateway {

    function initialize(
        address _owner, 
        address _relayer, 
        uint256 _finalizationFee,
        uint256 _counterpartyFee
    ) external initializer {
        relayer = _relayer;
        finalizationFee = _finalizationFee;
        counterpartyFee = _counterpartyFee;
        transferInitiatedIdx = 0;
        transferFinalizedIdx = 1; // First expected transfer index is 1
        __Ownable_init(_owner);
        __Pausable_init();
        __ReentrancyGuard_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    receive() external payable {}

    function _decrementMsgSender(uint256 _amount) internal override {
        require(msg.value == _amount, "Incorrect Ether value sent");
        // Wrapping function initiateTransfer is payable. Ether is escrowed in contract balance
    }

    function _fund(uint256 _amount, address _toFund) internal override {
        require(address(this).balance >= _amount, "Insufficient contract balance");
        (bool success, ) = _toFund.call{value: _amount}("");
        require(success, "Transfer to _toFund failed");
    }
}

