// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MevCommitAVS} from "../../../contracts/validator-registry/avs/MevCommitAVS.sol";
import {IDelegationManager} from "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import {IEigenPodManager} from "eigenlayer-contracts/src/contracts/interfaces/IEigenPodManager.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IAVSDirectory} from "eigenlayer-contracts/src/contracts/interfaces/IAVSDirectory.sol";
import {EigenHoodiReleaseConsts, EigenHoleskyReleaseConsts, EigenMainnetReleaseConsts} from "./ReleaseAddrConsts.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

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

contract DeployMainnet is BaseDeploy {
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    IDelegationManager constant public DELEGATION_MANAGER = IDelegationManager(EigenMainnetReleaseConsts.DELEGATION_MANAGER);
    IEigenPodManager constant public EIGENPOD_MANAGER = IEigenPodManager(EigenMainnetReleaseConsts.EIGENPOD_MANAGER);
    IStrategyManager constant public STRATEGY_MANAGER = IStrategyManager(EigenMainnetReleaseConsts.STRATEGY_MANAGER);
    IAVSDirectory constant public AVS_DIRECTORY = IAVSDirectory(EigenMainnetReleaseConsts.AVS_DIRECTORY);
    address constant public FREEZE_ORACLE = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    uint256 constant public UNFREEZE_FEE = 1 ether;
    address constant public UNFREEZE_RECEIVER = MainnetConstants.COMMITMENT_HOLDINGS_MULTISIG;
    uint256 constant public UNFREEZE_PERIOD_BLOCKS = 12000; // ~ 1 day
    uint256 constant public OPERATOR_DEREG_PERIOD_BLOCKS = 50000; // 50k * 12s ~= 1 week, which suffices for short-term manual slashing.
    uint256 constant public VALIDATOR_DEREG_PERIOD_BLOCKS = 50000; // 50k * 12s ~= 1 week, which suffices for short-term manual slashing.
    uint256 constant public LST_RESTARKER_DEREG_PERIOD_BLOCKS = 50000; // 50k * 12s ~= 1 week, which suffices for short-term manual slashing.

    function run() external {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();
        address[] memory restakeableStrategies = new address[](13);
        restakeableStrategies[0] = EigenMainnetReleaseConsts.STRATEGY_BASE_CBETH;
        restakeableStrategies[1] = EigenMainnetReleaseConsts.STRATEGY_BASE_STETH;
        restakeableStrategies[2] = EigenMainnetReleaseConsts.STRATEGY_BASE_RETH;
        restakeableStrategies[3] = EigenMainnetReleaseConsts.STRATEGY_BASE_ETHX;
        restakeableStrategies[4] = EigenMainnetReleaseConsts.STRATEGY_BASE_ANKRETH;
        restakeableStrategies[5] = EigenMainnetReleaseConsts.STRATEGY_BASE_OETH;
        restakeableStrategies[6] = EigenMainnetReleaseConsts.STRATEGY_BASE_OSETH;
        restakeableStrategies[7] = EigenMainnetReleaseConsts.STRATEGY_BASE_SWETH;
        restakeableStrategies[8] = EigenMainnetReleaseConsts.STRATEGY_BASE_WBETH;
        restakeableStrategies[9] = EigenMainnetReleaseConsts.STRATEGY_BASE_SFRXETH;
        restakeableStrategies[10] = EigenMainnetReleaseConsts.STRATEGY_BASE_LSETH;
        restakeableStrategies[11] = EigenMainnetReleaseConsts.STRATEGY_BASE_METH;
        restakeableStrategies[12] = EigenMainnetReleaseConsts.BEACON_CHAIN_ETH;

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

contract DeployHoodi is BaseDeploy {
    // This is the most important field. On mainnet it'll be the primev multisig.
    address constant public OWNER = 0x1623fE21185c92BB43bD83741E226288B516134a;

    IDelegationManager constant public DELEGATION_MANAGER = IDelegationManager(EigenHoodiReleaseConsts.DELEGATION_MANAGER);
    IEigenPodManager constant public EIGENPOD_MANAGER = IEigenPodManager(EigenHoodiReleaseConsts.EIGENPOD_MANAGER);
    IStrategyManager constant public STRATEGY_MANAGER = IStrategyManager(EigenHoodiReleaseConsts.STRATEGY_MANAGER);
    IAVSDirectory constant public AVS_DIRECTORY = IAVSDirectory(EigenHoodiReleaseConsts.AVS_DIRECTORY);
    address constant public FREEZE_ORACLE = 0x1623fE21185c92BB43bD83741E226288B516134a; // Temporary freeze oracle
    uint256 constant public UNFREEZE_FEE = 0.1 ether;
    address constant public UNFREEZE_RECEIVER = 0x1623fE21185c92BB43bD83741E226288B516134a; // Temporary unfreezeReceiver
    uint256 constant public UNFREEZE_PERIOD_BLOCKS = 12000; // ~ 1 day
    uint256 constant public OPERATOR_DEREG_PERIOD_BLOCKS = 12000; // ~ 1 day
    uint256 constant public VALIDATOR_DEREG_PERIOD_BLOCKS = 32 * 3; // 2 epoch finalization time + settlement buffer
    uint256 constant public LST_RESTARKER_DEREG_PERIOD_BLOCKS = 12000; // ~ 1 day

    function run() external {
        require(block.chainid == 23118, "must deploy on Hoodi");

        address[] memory restakeableStrategies = new address[](11);
        restakeableStrategies[0] = EigenHoodiReleaseConsts.STRATEGY_BASE_STETH;
        restakeableStrategies[1] = EigenHoodiReleaseConsts.STRATEGY_BASE_WETH;
        restakeableStrategies[3] = EigenHoodiReleaseConsts.STRATEGY_BASE_EIGEN;
        restakeableStrategies[4] = EigenHoodiReleaseConsts.BEACON_CHAIN_ETH;

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
