// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

import {IRegistry} from "symbiotic-core/interfaces/common/IRegistry.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

contract RegistryMock is IRegistry {
    using EnumerableSet for EnumerableSet.AddressSet;

    EnumerableSet.AddressSet private _entities;

    modifier checkEntity(
        address account
    ) {
        _checkEntity(account);
        _;
    }

    function register() external returns (address) {
        _addEntity(msg.sender);
        return msg.sender;
    }

    function isEntity(
        address entity_
    ) public view returns (bool) {
        return _entities.contains(entity_);
    }

    function totalEntities() public view returns (uint256) {
        return _entities.length();
    }

    function entity(
        uint256 index
    ) public view returns (address) {
        return _entities.at(index);
    }

    function _addEntity(
        address entity_
    ) internal {
        _entities.add(entity_);

        emit AddEntity(entity_);
    }

    function _checkEntity(
        address account
    ) internal view {
        if (!isEntity(account)) {
            revert EntityNotExist();
        }
    }
}
