// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IGateway} from "../interfaces/IGateway.sol";
import {GatewayStorage} from "./GatewayStorage.sol";

abstract contract Gateway is IGateway, GatewayStorage,
    Ownable2StepUpgradeable, UUPSUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable {   

    modifier onlyRelayer() {
        require(msg.sender == relayer, SenderNotRelayer(msg.sender, relayer));
        _;
    }

    /// @dev Initiates a transfer from the source chain gateway to its counterparty gateway on another chain.
    /// @notice The _recipient is transferred eth on the destination chain via solidity's send function (with built-in gas limit).
    /// Therefore the _recipient MUST be an EOA on the other chain to gauruntee a successful transfer.
    /// @notice If _recipient is a contract, manual withdrawal may be required on the counterparty chain.
    /// @notice The caller of this function takes responsiblity for whatever address is specified as the _recipient.
    /// That is, if _recipient is a contract with an immutable receiver that reverts, the user would be at fault for loss of funds.
    /// @param _recipient The address to receive the tokens.
    /// @param _amount The amount of Ether to transfer in wei.
    /// @return returnIdx The index of the initiated transfer.
    function initiateTransfer(address _recipient, uint256 _amount) 
        external payable whenNotPaused nonReentrant returns (uint256 returnIdx) {
        require(_amount >= counterpartyFinalizationFee, AmountTooSmall(_amount, counterpartyFinalizationFee));
        _decrementMsgSender(_amount);
        ++transferInitiatedIdx;
        emit TransferInitiated(msg.sender, _recipient, _amount, transferInitiatedIdx, counterpartyFinalizationFee);
        return transferInitiatedIdx;
    }

    /// @dev Finalizes a transfer as the destination chain gateway.
    /// @dev The inheriting contract MUST implement eth transfer failure handling, and the retry capability.
    /// @param _recipient The address to receive the tokens.
    /// @param _amount The amount of Ether to transfer in wei.
    /// @param _counterpartyIdx The index of the counterparty gateway contract.
    /// @param _finalizationFee The finalization fee (wei) paid to the relayer by the this contract.
    /// @notice The relayer is responsible for ensuring that the _finalizationFee value is the same as what was specified 
    /// as `counterpartyFinalizationFee` in the corresponding TransferInitiated event.
    function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx, uint256 _finalizationFee) 
        external onlyRelayer whenNotPaused nonReentrant {
        require(_amount >= _finalizationFee, AmountTooSmall(_amount, _finalizationFee));
        require(_counterpartyIdx == transferFinalizedIdx, InvalidCounterpartyIndex(_counterpartyIdx, transferFinalizedIdx));
        uint256 amountAfterFee = _amount - _finalizationFee;
        _fund(amountAfterFee, _recipient);
        _fund(_finalizationFee, relayer);
        ++transferFinalizedIdx;
        emit TransferFinalized(_recipient, _amount, _counterpartyIdx);
    }

    /// @dev Allows owner to pause the contract.
    function pause() external onlyOwner { _pause(); }

    /// @dev Allows owner to unpause the contract.
    function unpause() external onlyOwner { _unpause(); }

    /// @dev Allows owner to set a new relayer account.
    function setRelayer(address _relayer) external onlyOwner {
        require(_relayer != address(0), RelayerCannotBeZeroAddress());
        relayer = _relayer;
        emit RelayerSet(_relayer);
    }

    /// @dev Allows owner to set a new counterparty finalization fee.
    function setCounterpartyFinalizationFee(uint256 _counterpartyFinalizationFee) external onlyOwner {
        require(_counterpartyFinalizationFee > 0, CounterpartyFinalizationFeeTooSmall(_counterpartyFinalizationFee));
        counterpartyFinalizationFee = _counterpartyFinalizationFee;
        emit CounterpartyFinalizationFeeSet(_counterpartyFinalizationFee);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // @dev where _decrementMsgSender is implemented by inheriting contract.
    function _decrementMsgSender(uint256 _amount) internal virtual;

    // @dev where _fund is implemented by inheriting contract.
    function _fund(uint256 _amount, address _toFund) internal virtual;
}
