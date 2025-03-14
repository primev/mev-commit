// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.28;

interface IAllocator {

    /// @dev Emitted when a transfer needs withdrawal.
    event TransferNeedsWithdrawal(address indexed recipient, uint256 amount);

    /// @dev Emitted when a transfer is successful.
    event TransferSuccess(address indexed recipient, uint256 amount);

    error SenderNotWhitelisted(address sender);
    error InsufficientContractBalance(uint256 contractBalance, uint256 amountRequested);
    error NoFundsNeedingWithdrawal(address recipient);
    error TransferFailed(address recipient);

    function addToWhitelist(address _address) external;
    function removeFromWhitelist(address _address) external;
    function mint(address _mintTo, uint256 _amount) external;
    function pause() external;
    function unpause() external;
}
