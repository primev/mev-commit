// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.29;

import {Gateway} from "./Gateway.sol";
import {L1GatewayStorage} from "./L1GatewayStorage.sol";
import {Errors} from "../utils/Errors.sol";

/// @title L1Gateway
/// @notice Gateway contract deployed on L1 enabling the mev-commit standard bridge.
/// @dev This contract will escrow locked ETH, while a corresponding amount is minted from the SettlementGateway on the mev-commit chain.
contract L1Gateway is L1GatewayStorage, Gateway {

    /// @dev Emitted when a transfer needs withdrawal.
    event TransferNeedsWithdrawal(address indexed recipient, uint256 amount);

    /// @dev Emitted when a transfer is successful.
    event TransferSuccess(address indexed recipient, uint256 amount);

    error IncorrectEtherValueSent(uint256 msgValue, uint256 amountExpected);
    error InsufficientContractBalance(uint256 thisContractBalance, uint256 amountRequested);
    error NoFundsNeedingWithdrawal(address recipient);
    error TransferFailed(address recipient);

    /// @dev See https://docs.openzeppelin.com/upgrades-plugins/1.x/writing-upgradeable#initializing_the_implementation_contract
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @dev Receiver for native ETH.
    receive() external payable { }

    /// @dev Fallback function is disabled for this contract to prevent unintended interactions.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    function initialize(
        address _owner, 
        address _relayer, 
        uint256 _counterpartyFinalizationFee
    ) external initializer {
        relayer = _relayer;
        counterpartyFinalizationFee = _counterpartyFinalizationFee;
        transferInitiatedIdx = 0;
        transferFinalizedIdx = 1; // First expected transfer index is 1
        __Ownable_init(_owner);
        __Pausable_init();
        __ReentrancyGuard_init();
    }

    /// @dev Allows any account to manually withdraw funds that failed to be transferred by the relayer.
    /// @dev The relayer should NEVER call this function.
    function withdraw(address _recipient) external whenNotPaused nonReentrant {
        uint256 amount = transferredFundsNeedingWithdrawal[_recipient];
        require(amount > 0, NoFundsNeedingWithdrawal(_recipient));
        transferredFundsNeedingWithdrawal[_recipient] = 0;
        (bool success, ) = _recipient.call{value: amount}("");
        require(success, TransferFailed(_recipient));
        emit TransferSuccess(_recipient, amount);
    }

    function _decrementMsgSender(uint256 _amount) internal override {
        require(msg.value == _amount, IncorrectEtherValueSent(msg.value, _amount));
        // Wrapping function initiateTransfer is payable. Ether is escrowed in contract balance
    }

    function _fund(uint256 _amount, address _toFund) internal override {
        require(address(this).balance >= _amount, InsufficientContractBalance(address(this).balance, _amount));
        if (!payable(_toFund).send(_amount)) {
            transferredFundsNeedingWithdrawal[_toFund] += _amount;
            emit TransferNeedsWithdrawal(_toFund, _amount);
            return;
        } 
        emit TransferSuccess(_toFund, _amount);
    }
}

