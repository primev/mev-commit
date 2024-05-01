// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.20;

interface IProviderRegistry {
    function registerAndStake() external payable;

    function checkStake(address provider) external view returns (uint256);

    function depositFunds() external payable;

    function slash(
        uint256 amt,
        address provider,
        address payable bidder,
        uint256 residualBidPercentAfterDecay
    ) external;
}
