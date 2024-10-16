#!/bin/sh

# Default to core deploy type if not specified
DEPLOY_TYPE=${DEPLOY_TYPE:-core}

# Set the forge binary path, default to 'forge' if not provided
FORGE_BIN_PATH=${FORGE_BIN_PATH:-forge}

# Define the script path prefix, default to 'scripts/' if not provided
SCRIPT_PATH_PREFIX=${SCRIPT_PATH_PREFIX:-scripts/}

KEYSTORE_DIR=${KEYSTORE_DIR:-"$PWD/deployer_keystore"}
KEYSTORE_FILENAME=${KEYSTORE_FILENAME:-*}
KEYSTORE_PASSWORD=${KEYSTORE_PASSWORD:-"pwd"}
CONTRACT_REPO_ROOT_PATH=${CONTRACT_REPO_ROOT_PATH:-$PWD}

if [ "${DEPLOY_TYPE}" = "core" ]; then
    if [ -z "$ORACLE_KEYSTORE_ADDRESS" ]; then
        echo "ORACLE_KEYSTORE_ADDRESS not specified"
        exit 1
    fi
    echo "Deploying core contracts in testnet environment"
    ORACLE_KEYSTORE_ADDRESS="$ORACLE_KEYSTORE_ADDRESS" \
    "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}"core/DeployCore.s.sol:DeployTestnet \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --chain-id "${CHAIN_ID}" \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --skip-simulation \
        --use 0.8.26 \
        --broadcast \
        --force \
        --json \
        --via-ir
fi 
