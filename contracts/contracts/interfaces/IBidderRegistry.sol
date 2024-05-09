// SPDX-License-Identifier: BSL 1.1
pragma solidity ^0.8.15;

interface IBidderRegistry {
    struct PreConfCommitment {
        string txnHash;
        uint64 bid;
        uint64 blockNumber;
        string bidHash;
        string bidSignature;
        string commitmentHash;
        string commitmentSignature;
    }

    struct BidState {
        address bidder;
        uint64 bidAmt;
        State state;
    }

    enum State {
        Undefined,
        PreConfirmed,
        Withdrawn
    }

    function OpenBid(bytes32 commitmentDigest, uint64 bid, address bidder, uint64 blockNumber) external;

    function getDeposit(address bidder, uint256 window) external view returns (uint256);

    function depositForSpecificWindow(uint256 window) external payable;

    function retrieveFunds(
        uint256 windowToSettle,
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external;

    function unlockFunds(uint256 windowToSettle, bytes32 bidID) external;
}
