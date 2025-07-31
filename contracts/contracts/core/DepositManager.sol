// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IBidderRegistry} from "../interfaces/IBidderRegistry.sol";

contract DepositManager {

    mapping(address => uint256) public targetDeposits;
    address public immutable bidderRegistry;
    uint256 public immutable minBalance;
    error NotThisEOA();

    event TargetDepositSet(address indexed provider, uint256 amount);
    event TargetDepositDoesNotExist(address indexed provider);
    event NotEnoughEOABalance(uint256 balance, uint256 minBalance);
    event TopUpReduced(address indexed provider, uint256 needed, uint256 available);
    event CurrentDepositIsSufficient(address indexed provider);
    event DepositToppedUp(address indexed provider, uint256 amount);

    constructor(address _registry, uint256 _minBalance) {
        bidderRegistry = _registry;
        minBalance = _minBalance;
    }

    modifier onlyThisEOA() {
        require(msg.sender == address(this), NotThisEOA());
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
            emit CurrentDepositIsSufficient(provider);
            return;
        }
        uint256 needed = target - currentDeposit;

        uint256 balance = address(this).balance;
        if (balance <= minBalance) {
            emit NotEnoughEOABalance(balance, minBalance);
            return;
        }
        uint256 available = balance - minBalance;

        if (needed > available) {
            emit TopUpReduced(provider, needed, available);
            needed = available;
        }
        IBidderRegistry(bidderRegistry).depositAsBidder{value: needed}(provider);
        emit DepositToppedUp(provider, needed);
    }

    receive() external payable { 
        // Eth transfers allowed.
    }
}
