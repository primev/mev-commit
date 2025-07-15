// SPDX-License-Identifier: BSL 1.1

// solhint-disable no-console
// solhint-disable one-contract-per-file

pragma solidity 0.8.26;
import {RewardManager} from "../../../contracts/validator-registry/rewards/RewardManager.sol";
import {Script} from "forge-std/Script.sol";
import {Upgrades} from "openzeppelin-foundry-upgrades/Upgrades.sol";
import {MainnetConstants} from "../../MainnetConstants.sol";

contract DeployMainnet is Script {
    address constant public VANILLA_REGISTRY = 0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9;
    address constant public MEV_COMMIT_AVS = 0xBc77233855e3274E1903771675Eb71E602D9DC2e;
    address constant public MEV_COMMIT_MIDDLEWARE = 0x21fD239311B050bbeE7F32850d99ADc224761382;
    uint256 constant public AUTO_CLAIM_GAS_LIMIT = 250_000; 
    address constant public OWNER = MainnetConstants.PRIMEV_TEAM_MULTISIG;
    function run() public {
        require(block.chainid == 1, "must deploy on mainnet");
        vm.startBroadcast();

        address proxy = Upgrades.deployUUPSProxy(
            "RewardManager.sol",
            abi.encodeCall(
                RewardManager.initialize,
                (VANILLA_REGISTRY,
                MEV_COMMIT_AVS,
                MEV_COMMIT_MIDDLEWARE,
                AUTO_CLAIM_GAS_LIMIT,
                OWNER)
            )
        );
        console.log("RewardManager UUPS proxy deployed to:", address(proxy));
      vm.stopBroadcast();
    }
}
