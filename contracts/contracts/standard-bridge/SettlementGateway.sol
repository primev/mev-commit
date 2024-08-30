// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Gateway} from "./Gateway.sol";
import {IWhitelist} from "../interfaces/IWhitelist.sol";

contract SettlementGateway is Gateway {

    // Assuming deployer is 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,
    // whitelist's create2 addr should be 0x57508f0B0f3426758F1f3D63ad4935a7c9383620.
    // This variable is not hardcoded for testing purposes.
    address public whitelistAddr;
    
    function initialize(
        address _whitelistAddr,
        address _owner,
        address _relayer,
        uint256 _finalizationFee,
        uint256 _counterpartyFee
    ) external initializer {
        whitelistAddr = _whitelistAddr;
        relayer = _relayer;
        finalizationFee = _finalizationFee;
        counterpartyFee = _counterpartyFee;
        transferInitiatedIdx = 0;
        transferFinalizedIdx = 1; // First expected transfer index is 1
        __Ownable_init(_owner);
        __Pausable_init();
    }

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    // Burns native ether on settlement chain by sending it to the whitelist contract,
    // there should be equiv ether on L1 which will be UNLOCKED during finalization.
    function _decrementMsgSender(uint256 _amount) internal override {
        require(msg.value == _amount, "Incorrect Ether value sent");
        (bool success, ) = whitelistAddr.call{value: msg.value}("");
        require(success, "Failed to send Ether");
    }

    // Mints native ether on settlement chain via whitelist contract,
    // there should be equiv ether on L1 which remains LOCKED.
    function _fund(uint256 _amount, address _toFund) internal override {
        IWhitelist(whitelistAddr).mint(_toFund, _amount);
    }
}
