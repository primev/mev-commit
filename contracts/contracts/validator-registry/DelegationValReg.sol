// SPDX-License-Identifier: MIT

pragma solidity 0.8.20;

import "./ReputationalValReg.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";

contract DelegationValReg is OwnableUpgradeable, UUPSUpgradeable {
    ReputationalValReg public reputationalValReg;

    struct Delegation {
        address validatorEOA;
        uint256 amount;
    }

    mapping(address => Delegation) public delegations;

    event Delegated(address indexed delegator, address indexed validatorEOA, uint256 amount);
    event DelegationChanged(address indexed delegator, address indexed oldValidatorEOA, address indexed newValidatorEOA, uint256 amount);
    event Withdrawn(address indexed delegator, address indexed validatorEOA, uint256 amount);

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // TODO: Determine upgrade logic and test process
    }

    function initialize(
        address _owner,
        address _reputationalValReg
    ) external initializer {
        __Ownable_init(_owner);
        reputationalValReg = ReputationalValReg(payable(_reputationalValReg));
    }

    function delegate(address validatorEOA, uint256 amount) external {
        require(reputationalValReg.isEOAWhitelisted(validatorEOA), "Validator EOA must be whitelisted");
        require(delegations[msg.sender].amount == 0, "Already delegated");

        delegations[msg.sender] = Delegation({
            validatorEOA: validatorEOA,
            amount: amount
        });

        emit Delegated(msg.sender, validatorEOA, amount);
    }

    function changeDelegation(address newValidatorEOA) external {
        require(reputationalValReg.isEOAWhitelisted(newValidatorEOA), "New validator EOA must be whitelisted");
        Delegation storage delegation = delegations[msg.sender];
        require(delegation.amount > 0, "No active delegation");

        address oldValidatorEOA = delegation.validatorEOA;
        delegation.validatorEOA = newValidatorEOA;

        emit DelegationChanged(msg.sender, oldValidatorEOA, newValidatorEOA, delegation.amount);
    }

    function withdraw() external {
        Delegation storage delegation = delegations[msg.sender];
        require(delegation.amount > 0, "No active delegation");

        uint256 amount = delegation.amount;
        address validatorEOA = delegation.validatorEOA;

        delete delegations[msg.sender];

        emit Withdrawn(msg.sender, validatorEOA, amount);
    }

    fallback() external payable {
        revert("Invalid call");
    }

    receive() external payable {
        revert("Invalid call");
    }
}
// withdraw if the EOA becomes frozen or unwhitelisted, withdraw time, etc. 
