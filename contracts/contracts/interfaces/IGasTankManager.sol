// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

/**
 * @title IGasTankManager
 * @dev Interface for GasTankManager
 */
interface IGasTankManager {
    /// @dev Event to log successful update of the minimum deposit requirement.
    event MinimumDepositSet(uint256 indexed minDeposit);

    /// @dev Event to log successful gas tank top up.
    event GasTankToppedUp(uint256 indexed transferAmount);

    /// @notice Raised when a call expected from the contract itself is not.
    /// @param msgSender Address that invoked the call.
    /// @param thisAddress Address of the contract guard.
    error NotThisEOA(address msgSender, address thisAddress);

    /// @notice Raised when the contract balance cannot satisfy a refill.
    /// @param eoaAddress L1 EOA address.
    /// @param available Wei currently available on the L1 EOA address.
    /// @param needed Wei required to meet the minimum deposit requirement.
    error InsufficientEOABalance(address eoaAddress, uint256 available, uint256 needed);

    /// @notice Raised when forwarding funds to the provider reverts.
    /// @param rpcProvider Destination address for the transfer.
    /// @param transferAmount Amount of wei attempted.
    error GasTankTopUpFailed(address rpcProvider, uint256 transferAmount);

    /// @notice Raised when a gas tank balance already meets the minimum deposit requirement.
    /// @param currentBalance Gas tank balance reported by the provider.
    /// @param minDeposit Minimum deposit requirement.
    error GasTankBalanceIsSufficient(uint256 currentBalance, uint256 minDeposit);

    /// @notice Raised when the initial minimum deposit is not sufficient.
    /// @param sentAmount Amount of wei sent.
    /// @param requiredAmount Amount of wei required.
    error InsufficientFundsSent(uint256 sentAmount, uint256 requiredAmount);

    /// @notice Raised when the provider address is not set.
    /// @param provider Provider address.
    error ProviderNotSet(address provider);

    /// @notice Updates the minimum deposit requirement.
    /// @param minDeposit New minimum balance expressed in wei.
    /// @dev Only the owner can call this function.
    function setMinimumDeposit(uint256 minDeposit) external;

    /// @notice Requests a top-up when the gas tank's balance is below the minimum deposit requirement.
    /// @param gasTankBalance Balance that the provider currently holds on the gas tank.
    /// @dev Only the owner can call this function.
    /// @dev Reverts if the current gas tank balance is greater than or equal to the minimum deposit requirement.
    /// @dev Reverts if the contract balance is less than the amount needed to reach the minimum deposit requirement.
    /// @dev Always transfer the difference between the minimum deposit requirement and the current gas tank balance.
    function topUpGasTank(uint256 gasTankBalance) external;

    /// @notice Allows anyone to contribute funds. Forwards them to the provider immediately.
    function fundGasTank() external payable;

    /// @notice Sets the initial minimum deposit requirement. Forwards them to the provider immediately.
    /// @dev Reverts if the amount sent is less than the baseline deposit.
    function sendMinimumDeposit() external payable;
}
