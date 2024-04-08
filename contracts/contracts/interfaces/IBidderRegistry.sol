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

    function prepay() external payable;

    function LockBidFunds(bytes32 commitmentDigest, uint64 bid, address bidder) external;

    function getAllowance(address bidder) external view returns (uint256);

    function retrieveFunds(
        bytes32 commitmentDigest,
        address payable provider,
        uint256 residualBidPercentAfterDecay
    ) external;

    function unlockFunds(bytes32 bidID) external;
}
