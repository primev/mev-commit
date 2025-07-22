solhint:
	if !command -v solhint &>/dev/null; then \
		echo "Install solhint with 'npm install -g solhint'"; \
		exit 1; \
	fi
	solhint '**/*.sol'

deploy-avs-with-mock-eigen:
	forge clean
	forge script scripts/validator-registry/avs/DeployAVSWithMockEigen.s.sol:DeployAVSWithMockEigen --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --via-ir --broadcast

run-stake-example:
	forge script scripts/validator-registry/ValidatorExampleScript.s.sol:StakeExample --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --via-ir --broadcast

run-unstake-example:
	forge script scripts/validator-registry/ValidatorExampleScript.s.sol:UnstakeExample --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --via-ir --broadcast

run-withdraw-example:
	forge script scripts/validator-registry/ValidatorExampleScript.s.sol:WithdrawExample --rpc-url http://localhost:8545 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 --via-ir --broadcast

run-slash-example:
	forge script scripts/validator-registry/ValidatorExampleScript.s.sol:SlashExample --rpc-url http://localhost:8545 --private-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d --via-ir --broadcast

get-vanilla-registry-params:
	forge script scripts/validator-registry/GetParams.s.sol:GetVanillaRegistryParams --via-ir --rpc-url https://eth.llamarpc.com

get-mev-commit-avs-params:
	forge script scripts/validator-registry/GetParams.s.sol:GetMevCommitAVSParams --via-ir --rpc-url https://eth.llamarpc.com

setup-vault:
	forge script scripts/validator-registry/middleware/ExampleSetup.s.sol:SetupVault \
		--via-ir \
		--rpc-url https://1rpc.io/holesky \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--broadcast

deposit-to-vault:
	forge script scripts/validator-registry/middleware/ExampleSetup.s.sol:DepositToVault \
		--via-ir \
		--rpc-url https://1rpc.io/holesky \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--broadcast

operator-actions:
	forge script scripts/validator-registry/middleware/ExampleSetup.s.sol:OperatorActions \
		--via-ir \
		--rpc-url https://1rpc.io/holesky \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--broadcast

primev-team-actions:
	forge script scripts/validator-registry/middleware/ExampleSetup.s.sol:PrimevTeamActions \
		--via-ir \
		--rpc-url https://1rpc.io/holesky \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--broadcast

vault-actions:
	forge script scripts/validator-registry/middleware/ExampleSetup.s.sol:VaultActions \
		--via-ir \
		--rpc-url https://1rpc.io/holesky \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--broadcast

# Note if running upgrade on old geth chain, you need to set evm_version = "london" in foundry.toml
# PUSH0 opcode is not supported with older versions of geth
upgrade-registries:
	forge clean
	forge script scripts/core/UpgradeRegistries.s.sol:UpgradeRegistries \
		--via-ir \
		--rpc-url https://chainrpc.mev-commit.xyz/ \
		--keystore $$KEYSTORE \
		--password $$PASSWORD \
		--sender $$SENDER \
		--priority-gas-price 50000000 \
        --with-gas-price 60000000 \
		--verbosity 5 \
		--broadcast
