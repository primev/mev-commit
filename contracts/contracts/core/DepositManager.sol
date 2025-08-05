// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";
import {Errors} from "../utils/Errors.sol";

contract DepositManager {

    mapping(address => uint256) public targetDeposits;
    address public immutable bidderRegistry;
    uint256 public immutable minBalance;
    error NotThisEOA(address msgSender, address thisAddress);

    event TargetDepositSet(address indexed provider, uint256 amount);
    event TargetDepositDoesNotExist(address indexed provider);
    event CurrentDepositIsSufficient(address indexed provider, uint256 currentDeposit, uint256 targetDeposit);
    event CurrentBalanceAtOrBelowMin(address indexed provider, uint256 currentBalance, uint256 minBalance);
    event NotEnoughEOABalance(address indexed provider, uint256 available, uint256 needed);
    event TopUpReduced(address indexed provider, uint256 available, uint256 needed);
    event DepositToppedUp(address indexed provider, uint256 amount);

    constructor(address _registry, uint256 _minBalance) {
        bidderRegistry = _registry;
        minBalance = _minBalance;
    }

    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA(msg.sender, address(this)));
        _;
    }

    function setTargetDeposit(address provider, uint256 amount) external onlyThisEOA {
        targetDeposits[provider] = amount;
        emit TargetDepositSet(provider, amount);
    }

    /// @notice Top-up deposits if needed, as configured by this EOA.
    /// @param provider to top-up the deposit for.
    /// @dev This function will be called automatically by external addresses.
    function topUpDeposit(address provider) external {
        uint256 target = targetDeposits[provider];
        if (target == 0) {
            emit TargetDepositDoesNotExist(provider);
            return;
        }

        uint256 currentDeposit = IBidderRegistry(bidderRegistry).getDeposit(address(this), provider);
        if (currentDeposit >= target) {
            emit CurrentDepositIsSufficient(provider, currentDeposit, target);
            return;
        }
        uint256 needed = target - currentDeposit; // No underflow/overflow, target must be greater than current deposit

        uint256 currentBalance = address(this).balance;
        if (currentBalance <= minBalance) {
            emit CurrentBalanceAtOrBelowMin(provider, currentBalance, minBalance);
            return;
        }

        uint256 available = currentBalance - minBalance; // No underflow/overflow, currentBalance must be greater than minBalance
        if (available < needed) {
            emit TopUpReduced(provider, available, needed);
            needed = available;
        }
        IBidderRegistry(bidderRegistry).depositAsBidder{value: needed}(provider);
        emit DepositToppedUp(provider, needed);
    }

    fallback() external payable {
        revert Errors.InvalidFallback();
    }

    receive() external payable { 
        // Eth transfers allowed.
    }
}
