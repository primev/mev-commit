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

    error FinalizationFeeTooSmall(uint256 _finalizationFee);
    error CounterpartyFeeTooSmall(uint256 _counterpartyFee);

    modifier onlyRelayer() {
        require(msg.sender == relayer, SenderNotRelayer(msg.sender, relayer));
        _;
    }

    function initiateTransfer(address _recipient, uint256 _amount) 
        external payable whenNotPaused nonReentrant returns (uint256 returnIdx) {
        require(_amount >= counterpartyFee, AmountTooSmall(_amount, counterpartyFee));
        _decrementMsgSender(_amount);
        ++transferInitiatedIdx;
        emit TransferInitiated(msg.sender, _recipient, _amount, transferInitiatedIdx);
        return transferInitiatedIdx;
    }

    function finalizeTransfer(address _recipient, uint256 _amount, uint256 _counterpartyIdx) 
        external onlyRelayer whenNotPaused nonReentrant {
        require(_amount >= finalizationFee, AmountTooSmall(_amount, finalizationFee));
        require(_counterpartyIdx == transferFinalizedIdx, InvalidCounterpartyIndex(_counterpartyIdx, transferFinalizedIdx));
        uint256 amountAfterFee = _amount - finalizationFee;
        _fund(amountAfterFee, _recipient);
        _fund(finalizationFee, relayer);
        ++transferFinalizedIdx;
        emit TransferFinalized(_recipient, _amount, _counterpartyIdx);
    }

    /// @dev Allows owner to pause the contract.
    function pause() external onlyOwner { _pause(); }

    /// @dev Allows owner to unpause the contract.
    function unpause() external onlyOwner { _unpause(); }

    /// @dev Allows owner to set a new finalization fee.
    /// @notice If using this function, ensure the same value is set as the `counterpartyFee` in the counterparty contract.
    function setFinalizationFee(uint256 _finalizationFee) external onlyOwner {
        require(_finalizationFee > 0, FinalizationFeeTooSmall(_finalizationFee));
        finalizationFee = _finalizationFee;
        emit FinalizationFeeSet(_finalizationFee);
    }

    /// @dev Allows owner to set a new counterparty fee.
    /// @notice If using this function, ensure the same value is set as the `finalizationFee` in the counterparty contract.
    function setCounterpartyFee(uint256 _counterpartyFee) external onlyOwner {
        require(_counterpartyFee > 0, CounterpartyFeeTooSmall(_counterpartyFee));
        counterpartyFee = _counterpartyFee;
        emit CounterpartyFeeSet(_counterpartyFee);
    }

    // solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    // @dev where _decrementMsgSender is implemented by inheriting contract.
    function _decrementMsgSender(uint256 _amount) internal virtual;

    // @dev where _fund is implemented by inheriting contract.
    function _fund(uint256 _amount, address _toFund) internal virtual;
}
