// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IWhitelist {
    function mint(address _mintTo, uint256 _amount) external;
    function burn(address _burnFrom, uint256 _amount) external;
}
