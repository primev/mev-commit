// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Errors} from "../utils/Errors.sol";

/// @title GasTankDepositor
/// @notice Coordinates on-demand ETH Transfers to the RPC Service for EOA custodial gas tanks.
/// @dev This contract implicitly trusts the RPC_SERVICE address.
contract GasTankDepositor {
    address public immutable RPC_SERVICE;
    uint256 public immutable MAXIMUM_DEPOSIT;
    address public immutable GAS_TANK_ADDRESS;

    event GasTankFunded(address indexed smartAccount, address indexed caller, uint256 indexed amount);

    error FailedToFundGasTank(address rpcProvider, uint256 transferAmount);
    error RPCServiceNotSet(address provider);
    error NotRPCService(address caller);
    error InsufficientFunds(uint256 currentBalance, uint256 requiredBalance);
    error NotThisEOA(address msgSender, address thisAddress);
    error MaximumDepositNotMet(uint256 amountToTransfer, uint256 maximumDeposit);

    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    modifier onlyRPCService() {
        require(msg.sender == RPC_SERVICE, NotRPCService(msg.sender));
        _;
    }

    /// @dev Writes the variables into the contract bytecode.
    /// @dev No storage is used in this contract.
    constructor(address rpcService, uint256 _maxDeposit) {
        require(rpcService != address(0), RPCServiceNotSet(rpcService));
        require(_maxDeposit > 0, MaximumDepositNotMet(0, _maxDeposit));
        RPC_SERVICE = rpcService;
        MAXIMUM_DEPOSIT = _maxDeposit;
        GAS_TANK_ADDRESS = address(this);
    }

    receive() external payable {
        if (address(this) == GAS_TANK_ADDRESS) {
            revert Errors.InvalidReceive();
        }
    }

    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Transfers ETH from the EOA's balance to the Gas RPC Service.
    /// @param _amount The amount to fund the gas tank with.
    /// @dev Only the EOA can call this function.
    function fundGasTank(uint256 _amount) external onlyThisEOA {
        _fundGasTank(_amount);
    }

    /// @notice Transfers the maximum deposit amount of ETH from the EOA's balance to the Gas RPC Service.
    /// @dev Only the RPC Service can call this function.
    function fundGasTank() external onlyRPCService {
        _fundGasTank(MAXIMUM_DEPOSIT);
    }

    function _fundGasTank(uint256 _amountToTransfer) internal {
        require(address(this).balance >= _amountToTransfer, InsufficientFunds(address(this).balance, _amountToTransfer));

        (bool success,) = RPC_SERVICE.call{value: _amountToTransfer}("");
        if (!success) {
            revert FailedToFundGasTank(RPC_SERVICE, _amountToTransfer);
        }

        emit GasTankFunded(address(this), msg.sender, _amountToTransfer);
    }
}
