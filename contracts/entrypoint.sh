#!/bin/sh

# Default to core deploy type if not specified
DEPLOY_TYPE=${DEPLOY_TYPE:-core}

# Set the forge binary path, default to 'forge' if not provided
FORGE_BIN_PATH=${FORGE_BIN_PATH:-forge}

# Define the script path prefix, default to 'scripts/' if not provided
SCRIPT_PATH_PREFIX=${SCRIPT_PATH_PREFIX:-scripts}

KEYSTORE_DIR=${KEYSTORE_DIR:-"$PWD/deployer_keystore"}
KEYSTORE_FILENAME=${KEYSTORE_FILENAME:-*}
KEYSTORE_PASSWORD=${KEYSTORE_PASSWORD:-"pwd"}
CONTRACT_REPO_ROOT_PATH=${CONTRACT_REPO_ROOT_PATH:-$PWD}

if [ -n "${ETHERSCAN_API_KEY}" ]; then
    VERIFY_OPTION="--verify"
fi

if [ "${DEPLOY_TYPE}" = "core" ]; then
    if [ -z "$ORACLE_KEYSTORE_ADDRESS" ]; then
        echo "ORACLE_KEYSTORE_ADDRESS not specified"
        exit 1
    fi
    echo "Deploying core contracts in testnet environment"
    ORACLE_KEYSTORE_ADDRESS="$ORACLE_KEYSTORE_ADDRESS" \
    "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}/core/DeployCore.s.sol:DeployTestnet" \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --chain-id "${CHAIN_ID}" \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --use 0.8.26 \
        --broadcast \
        --json \
        --via-ir
elif [ "${DEPLOY_TYPE}" = "settlement-gateway" ]; then
    echo "Deploying gateway contract on settlement chain"
    "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}/standard-bridge/DeployStandardBridge.s.sol:DeploySettlementGateway" \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --chain-id "${CHAIN_ID}" \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --use 0.8.26 \
        --broadcast \
        --json \
        --via-ir
elif [ "${DEPLOY_TYPE}" = "l1-gateway" ]; then
    "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}/standard-bridge/DeployStandardBridge.s.sol:DeployL1Gateway" \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --rpc-url "https://eth-holesky.g.alchemy.com/v2/a4-kNJRrkOqfQza9FNIayasNghgQk5E6" \
        --chain-id "17000" \
        --private-key "0x120e26406fe7504f4877cc55d77fec99a26e8b01275fe72d46b9ccdd4f4ca27f" \
        --sender "0x88195Fe694bf1EA76cE8d3bc02D6abC2E7c17c77" \
        --use 0.8.26 \
        --broadcast \
        --json \
        ${VERIFY_OPTION} \
        --via-ir
fi
