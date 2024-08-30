// SPDX-License-Identifier: BSL 1.1
pragma solidity 0.8.26;

// solhint-disable no-console
// solhint-disable one-contract-per-file

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {StrategyManagerMock} from "eigenlayer-contracts/src/test/mocks/StrategyManagerMock.sol";
import {DelegationManagerMock} from "eigenlayer-contracts/src/test/mocks/DelegationManagerMock.sol";
import {EigenPodManagerMock} from "../../../test/validator-registry/avs/EigenPodManagerMock.sol";
import {AVSDirectoryMock} from "../../../test/validator-registry/avs/AVSDirectoryMock.sol";

contract DeployAVSWithMockEigen is Script {
    function run() external {
        require(block.chainid == 31337, "must deploy on anvil");
        vm.startBroadcast();

        StrategyManagerMock strategyManagerMock = new StrategyManagerMock();
        DelegationManagerMock delegationManagerMock = new DelegationManagerMock();
        EigenPodManagerMock eigenPodManagerMock = new EigenPodManagerMock();
        AVSDirectoryMock avsDirectoryMock = new AVSDirectoryMock();

        address[] memory restakeableStrategies = new address[](3);
        restakeableStrategies[0] = address(0x1);
        restakeableStrategies[1] = address(0x2);
        restakeableStrategies[2] = address(0x3);

        address freezeOracle = address(0x5);
        uint256 unfreezeFee = 1 ether;
        address unfreezeReceiver = address(0x6);
        uint256 unfreezePeriodBlocks = 100;
        uint256 operatorDeregPeriodBlocks = 200;
        uint256 validatorDeregPeriodBlocks = 300;
        uint256 lstRestakerDeregPeriodBlocks = 400;
        string memory metadataUrl = "https://raw.githubusercontent.com/primev/mev-commit/main/static/avs-metadata.json";

        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitAVS.sol",
            abi.encodeCall(MevCommitAVS.initialize, (
                msg.sender,
                delegationManagerMock, 
                eigenPodManagerMock, 
                strategyManagerMock, 
                avsDirectoryMock, 
                restakeableStrategies,
                freezeOracle,
                unfreezeFee,
                unfreezeReceiver,
                unfreezePeriodBlocks,
                operatorDeregPeriodBlocks,
                validatorDeregPeriodBlocks,
                lstRestakerDeregPeriodBlocks,
                metadataUrl
            ))
        );
        MevCommitAVS mevCommitAVS = MevCommitAVS(payable(proxy));

        console.log("StrategyManagerMock deployed at:", address(strategyManagerMock));
        console.log("DelegationManagerMock deployed at:", address(delegationManagerMock));
        console.log("EigenPodManagerMock deployed at:", address(eigenPodManagerMock));
        console.log("AVSDirectoryMock deployed at:", address(avsDirectoryMock));
        console.log("MevCommitAVS deployed at:", address(mevCommitAVS));

        delegationManagerMock.setIsOperator(msg.sender, true);

        vm.stopBroadcast();
    }
}
