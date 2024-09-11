// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.9;

import "forge-std/Test.sol";
import {EigenPodMock} from "./EigenPodMock.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IEigenPod} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPod.sol";
import {IStrategy} from "eigenlayer-contracts/src/contracts/interfaces/IStrategy.sol";
import "@openzeppelin/contracts/proxy/beacon/IBeacon.sol";
import {IBeaconChainOracle} from "eigenlayer-contracts/src/contracts/interfaces/IBeaconChainOracle.sol";
import {IETHPOSDeposit} from "eigenlayer-contracts/src/contracts/interfaces/IETHPOSDeposit.sol";
import {ISlasher} from "eigenlayer-contracts/src/contracts/interfaces/ISlasher.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IPauserRegistry} from "eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";

// Similar to eigenlayer core's EigenPodManagerMock but their mocks don't use virtual functions.
contract EigenPodManagerMock is IEigenPodManager, Test {

    mapping(address => EigenPodMock) public pods;

    function setMockPod(address podOwner, EigenPodMock pod) external {
        pods[podOwner] = pod;
    }

    function getPod(address podOwner) external view returns(IEigenPod) {
        return pods[podOwner];
    }

    IStrategy public constant beaconChainETHStrategy = IStrategy(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0);

    IBeacon public eigenPodBeacon;

    IETHPOSDeposit public ethPOS;

    mapping(address => int256) public podShares;

    function slasher() external view returns(ISlasher) {}

    function createPod() external returns(address) {}

    function stake(bytes calldata /*pubkey*/, bytes calldata /*signature*/, bytes32 /*depositDataRoot*/) external payable {}

    function recordBeaconChainETHBalanceUpdate(address /*podOwner*/, int256 /*sharesDelta*/) external pure {}
    
    function updateBeaconChainOracle(IBeaconChainOracle /*newBeaconChainOracle*/) external pure {}

    function ownerToPod(address /*podOwner*/) external pure returns(IEigenPod) {
        return IEigenPod(address(0));
    }

    function beaconChainOracle() external pure returns(IBeaconChainOracle) {
        return IBeaconChainOracle(address(0));
    }   

    function getBlockRootAtTimestamp(uint64 /*timestamp*/) external pure returns(bytes32) {
        return bytes32(0);
    }

    function strategyManager() external pure returns(IStrategyManager) {
        return IStrategyManager(address(0));
    }
    
    function hasPod(address podOwner) external view returns (bool) {
        return pods[podOwner] != EigenPodMock(address(0));
    }

    function pause(uint256 /*newPausedStatus*/) external{}

    function pauseAll() external{}

    function paused() external pure returns (uint256) {
        return 0;
    }
    
    function paused(uint8 /*index*/) external pure returns (bool) {
        return false;
    }

    function setPauserRegistry(IPauserRegistry /*newPauserRegistry*/) external {}

    function pauserRegistry() external pure returns (IPauserRegistry) {
        return IPauserRegistry(address(0));
    }

    function unpause(uint256 /*newPausedStatus*/) external{}

    function podOwnerShares(address podOwner) external view returns (int256) {
        return podShares[podOwner];
    }

    function setPodOwnerShares(address podOwner, int256 shares) external {
        podShares[podOwner] = shares;
    }

    function addShares(address /*podOwner*/, uint256 shares) external pure returns (uint256) {
        // this is the "increase in delegateable tokens"
        return (shares);
    }

    function withdrawSharesAsTokens(address podOwner, address destination, uint256 shares) external {}

    function removeShares(address podOwner, uint256 shares) external {}

    function numPods() external view returns (uint256) {}

    function denebForkTimestamp() external pure returns (uint64) {
        return type(uint64).max;
    }

    function setDenebForkTimestamp(uint64 timestamp) external{}
}