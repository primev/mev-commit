solhint:
	if !command -v solhint &>/dev/null; then \
		echo "Install solhint with 'npm install -g solhint'"; \
		exit 1; \
	fi
	solhint '**/*.sol'

deploy-val-reg:
	forge clean
	forge script scripts/validator-registry/DeployVanillaRegistry.s.sol:DeployHolesky \
		--rpc-url http://localhost:8545 \
		--private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
		--broadcast \
		--chain-id 31337 \
		-vvvv \
		--use 0.8.26 \
		--via-ir \

deploy-opt-in-router-holesky:
	forge clean
	forge script scripts/validator-registry/DeployValidatorOptInRouter.s.sol:DeployHolesky \
		--rpc-url https://ethereum-holesky-rpc.publicnode.com \
		--keystores "${KEYSTORES}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
		--via-ir \
		--chain-id 17000 \
		--use 0.8.26 \
		--broadcast

deploy-core:
	forge clean
	DEPLOY_TYPE="core" RPC_URL="http://localhost:8545" PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" CHAIN_ID="17864" ./entrypoint.sh

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
