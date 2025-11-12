// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {IGasTankManager} from "../interfaces/IGasTankManager.sol";
import {GasTankManagerStorage} from "./GasTankManagerStorage.sol";
import {Errors} from "../utils/Errors.sol";

/// @title GasTankManager
/// @notice Coordinates on-demand ETH Transfers to the RPC Service for EOA custodial gas tanks.
/// @dev Flow overview:
/// - EOA (Prerequisites for use)
///   - Authorizes and sets this contract as a delegate against its own EOA address. (ERC-7702 compliant)
///   - Sends the initial minimum deposit via `initializeGasTank` to the RPC Service.
/// - RPC Service
///   - Triggers `topUpGasTank`, transferring the `minDeposit` when the gas tank requires funding.
///   - These funds are then transferred to the RPC Service's custodial gas tank.
contract GasTankManager is IGasTankManager, GasTankManagerStorage, Ownable2StepUpgradeable, UUPSUpgradeable {
    /// @notice Restricts calls to those triggered internally via `onlyThisEOA`.
    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    modifier isValidCaller() {
        require(msg.sender == address(this) || msg.sender == owner(), NotValidCaller(msg.sender));
        _;
    }

    /// @notice Locks the implementation upon deployment.
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /// @notice Accepts direct ETH deposits.
    receive() external payable {}

    /// @notice Reverts on any call data to keep the interface surface narrow.
    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    /// @notice Initializes ownership and the minimum deposit requirement.
    /// @param _owner EOA managed by the RPC Service.
    /// @param _minDeposit Minimum deposit requirement.
    function initialize(address _owner, uint256 _minDeposit) external initializer {
        minDeposit = _minDeposit;
        __Ownable_init(_owner);
        __UUPSUpgradeable_init();
    }

    /// @inheritdoc IGasTankManager
    function setMinimumDeposit(uint256 _minDeposit) external onlyOwner {
        minDeposit = _minDeposit;
        emit MinimumDepositSet(_minDeposit);
    }

    /// @inheritdoc IGasTankManager
    function topUpGasTank() external onlyOwner {
        _sendFundsToProvider(minDeposit);
    }

    /// @inheritdoc IGasTankManager
    function initializeGasTank() external onlyThisEOA {
        _sendFundsToProvider(minDeposit);
    }

    /// @inheritdoc IGasTankManager
    function fundGasTank() external payable isValidCaller {
        _sendFundsToProvider(msg.value);
    }

    /// solhint-disable-next-line no-empty-blocks
    function _authorizeUpgrade(address) internal override onlyOwner {}

    /// @notice Forwards ETH to the rpcService.
    /// @param _amount Gas tank top-up amount.
    function _sendFundsToProvider(uint256 _amount) internal {
        address rpcService = owner();

        require(rpcService != address(0), ProviderNotSet(rpcService));
        require(_amount > 0 && address(this).balance >= _amount, InvalidAmount());

        (bool success,) = rpcService.call{value: _amount}("");
        require(success, GasTankTopUpFailed(rpcService, _amount));

        emit GasTankToppedUp(address(this), msg.sender, _amount);
    }
}
