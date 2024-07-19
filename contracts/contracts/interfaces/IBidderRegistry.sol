// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.20;

interface IBidderRegistry {
    struct PreConfCommitment {
        string txnHash;
        uint256 bid;
        uint64 blockNumber;
        string bidHash;
        string bidSignature;
        string commitmentHash;
        string commitmentSignature;
    }

    struct BidState {
        address bidder;
        uint256 bidAmt;
        State state;
    }

    enum State {
        Undefined,
        PreConfirmed,
        Withdrawn
    }

    function OpenBid(bytes32 commitmentDigest, uint256 bid, address bidder, uint64 blockNumber) external;

    function depositForWindow(uint256 window) external payable;

    function retrieveFunds(
        uint256 windowToSettle,
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external;

    function unlockFunds(uint256 windowToSettle, bytes32 bidID) external;

    function getDeposit(address bidder, uint256 window) external view returns (uint256);
}
