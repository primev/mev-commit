// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

interface IAllocator {
    function addToWhitelist(address _address) external;
    function removeFromWhitelist(address _address) external;
    function mint(address _mintTo, uint256 _amount) external;
    function pause() external;
    function unpause() external;
}
