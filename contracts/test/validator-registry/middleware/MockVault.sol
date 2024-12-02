// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

contract MockVault {
    address public delegator;
    address public slasher;
    address public burner;
    uint48 private _epochDuration;

    constructor(address _delegator, address _slasher, address _burner, uint48 epochDuration_) {
        delegator = _delegator;
        slasher = _slasher;
        burner = _burner;
        _epochDuration = epochDuration_;
    }

    function setSlasher(address _slasher) external {
        slasher = _slasher;
    }

    function setEpochDuration(uint48 epochDuration_) external {
        _epochDuration = epochDuration_;
    }

    function setBurner(address _burner) external {
        burner = _burner;
    }

    function epochDuration() external view returns (uint48) {
        return _epochDuration;
    }
}
