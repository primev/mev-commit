// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Gateway} from "./Gateway.sol";
import {IAllocator} from "../interfaces/IAllocator.sol";

/// @title SettlementGateway
/// @notice Gateway contract deployed on the mev-commit chain enabling the mev-commit standard bridge.
contract SettlementGateway is Gateway {

    address public allocatorAddr;

    error IncorrectEtherValueSent(uint256 msgValue, uint256 amountExpected);
    error TransferFailed(address recipient, uint256 amount);
    
    function initialize(
        address _allocatorAddr,
        address _owner,
        address _relayer,
        uint256 _finalizationFee,
        uint256 _counterpartyFee
    ) external initializer {
        allocatorAddr = _allocatorAddr;
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

    // Burns native ether on settlement chain by sending it to the allocator contract,
    // there should be equiv ether on L1 which will be UNLOCKED during finalization.
    function _decrementMsgSender(uint256 _amount) internal override {
        require(msg.value == _amount, IncorrectEtherValueSent(msg.value, _amount));
        (bool success, ) = allocatorAddr.call{value: msg.value}("");
        require(success, TransferFailed(allocatorAddr, msg.value));
    }

    // Mints native ether on settlement chain via allocator contract,
    // there should be equiv ether on L1 which remains LOCKED.
    function _fund(uint256 _amount, address _toFund) internal override {
        IAllocator(allocatorAddr).mint(_toFund, _amount);
    }
}
