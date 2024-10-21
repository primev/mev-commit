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
