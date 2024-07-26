// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

interface IProviderRegistry {
    function registerAndStake(bytes calldata blsPublicKey) external payable;

    function depositFunds() external payable;

    function slash(
        uint256 amt,
        address provider,
        address payable bidder,
        uint256 residualBidPercentAfterDecay
    ) external;

    function checkStake(address provider) external view returns (uint256);

    function minStake() external view returns (uint256);

    function withdrawalRequests(address provider) external view returns (uint256);
}
