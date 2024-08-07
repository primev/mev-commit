// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.20;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {EigenHoleskyReleaseConsts} from "./ReleaseAddrConsts.sol";

contract BaseDeploy is Script {
    function deployMevCommitAVS(
        address owner,
        IDelegationManager delegationManager,
        IEigenPodManager eigenPodManager,
        IStrategyManager strategyManager,
        IAVSDirectory avsDirectory,
        address[] memory restakeableStrategies,
        address freezeOracle,
        uint256 unfreezeFee,
        address unfreezeReceiver,
        uint256 unfreezePeriodBlocks,
        uint256 operatorDeregPeriodBlocks,
        uint256 validatorDeregPeriodBlocks,
        uint256 lstRestakerDeregPeriodBlocks
    ) public returns (address) {
        console.log("Deploying MevCommitAVS on chain:", block.chainid);
        string memory metadataUrl = "https://raw.githubusercontent.com/primev/mev-commit/main/static/avs-metadata.json";
        address proxy = Upgrades.deployUUPSProxy(
            "MevCommitAVS.sol",
            abi.encodeCall(
                MevCommitAVS.initialize,
                (
                    owner,
                    delegationManager,
                    eigenPodManager,
                    strategyManager,
                    avsDirectory,
                    restakeableStrategies,
                    freezeOracle,
                    unfreezeFee,
                    unfreezeReceiver,
                    unfreezePeriodBlocks,
                    operatorDeregPeriodBlocks,
                    validatorDeregPeriodBlocks,
                    lstRestakerDeregPeriodBlocks,
                    metadataUrl
                )
            )
        );
        console.log("MevCommitAVS UUPS proxy deployed to:", address(proxy));
        MevCommitAVS mevCommitAVS = MevCommitAVS(payable(proxy));
        console.log("MevCommitAVS owner:", mevCommitAVS.owner());
        return proxy;
    }
}

contract DeployHolesky is BaseDeploy {
    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1;

    IDelegationManager constant public DELEGATION_MANAGER = IDelegationManager(EigenHoleskyReleaseConsts.DELEGATION_MANAGER);
    IEigenPodManager constant public EIGENPOD_MANAGER = IEigenPodManager(EigenHoleskyReleaseConsts.EIGENPOD_MANAGER);
    IStrategyManager constant public STRATEGY_MANAGER = IStrategyManager(EigenHoleskyReleaseConsts.STRATEGY_MANAGER);
    IAVSDirectory constant public AVS_DIRECTORY = IAVSDirectory(EigenHoleskyReleaseConsts.AVS_DIRECTORY);
    address constant public FREEZE_ORACLE = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1; // Temporary freeze oracle
    uint256 constant public UNFREEZE_FEE = 0.1 ether;
    address constant public UNFREEZE_RECEIVER = 0x4535bd6fF24860b5fd2889857651a85fb3d3C6b1; // Temporary unfreezeReceiver
    uint256 constant public UNFREEZE_PERIOD_BLOCKS = 12000; // ~ 1 day
    uint256 constant public OPERATOR_DEREG_PERIOD_BLOCKS = 12000; // ~ 1 day
    uint256 constant public VALIDATOR_DEREG_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer
    uint256 constant public LST_RESTARKER_DEREG_PERIOD_BLOCKS = 12000; // ~ 1 day

    function run() external {
        require(block.chainid == 17000, "must deploy on Holesky");

        address[] memory restakeableStrategies = new address[](11);
        restakeableStrategies[0] = EigenHoleskyReleaseConsts.STRATEGY_BASE_STETH;
        restakeableStrategies[1] = EigenHoleskyReleaseConsts.STRATEGY_BASE_RETH;
        restakeableStrategies[2] = EigenHoleskyReleaseConsts.STRATEGY_BASE_WETH;
        restakeableStrategies[3] = EigenHoleskyReleaseConsts.STRATEGY_BASE_LSETH;
        restakeableStrategies[4] = EigenHoleskyReleaseConsts.STRATEGY_BASE_SFRXETH;
        restakeableStrategies[5] = EigenHoleskyReleaseConsts.STRATEGY_BASE_ETHX;
        restakeableStrategies[6] = EigenHoleskyReleaseConsts.STRATEGY_BASE_OSETH;
        restakeableStrategies[7] = EigenHoleskyReleaseConsts.STRATEGY_BASE_CBETH;
        restakeableStrategies[8] = EigenHoleskyReleaseConsts.STRATEGY_BASE_METH;
        restakeableStrategies[9] = EigenHoleskyReleaseConsts.STRATEGY_BASE_ANKRETH;
        restakeableStrategies[10] = EigenHoleskyReleaseConsts.BEACON_CHAIN_ETH;

        vm.startBroadcast();
        deployMevCommitAVS(
            OWNER,
            DELEGATION_MANAGER,
            EIGENPOD_MANAGER,
            STRATEGY_MANAGER,
            AVS_DIRECTORY,
            restakeableStrategies,
            FREEZE_ORACLE, 
            UNFREEZE_FEE, 
            UNFREEZE_RECEIVER, 
            UNFREEZE_PERIOD_BLOCKS, 
            OPERATOR_DEREG_PERIOD_BLOCKS, 
            VALIDATOR_DEREG_PERIOD_BLOCKS, 
            LST_RESTARKER_DEREG_PERIOD_BLOCKS
        );
        vm.stopBroadcast();
    }
}
