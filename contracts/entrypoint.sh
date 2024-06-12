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
    echo "Deploying core contracts"
    "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}"DeployScripts.s.sol:DeployScript \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --chain-id "${CHAIN_ID}" \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --skip-simulation \
        --use 0.8.23 \
        --broadcast \
        --force \
        --json \
        --via-ir

elif [ "${DEPLOY_TYPE}" = "transfer-ownership" ]; then
    if [ -z "$ORACLE_KEYSTORE_ADDRESS" ]; then
        echo "ORACLE_KEYSTORE_ADDRESS not specified"
        exit 1
    fi
    if [ -z "$BLOCK_TRACKER_ADDRESS" ]; then
        echo "BLOCK_TRACKER_ADDRESS not specified"
        exit 1
    fi
    if [ -z "$ORACLE_ADDRESS" ]; then
        echo "ORACLE_ADDRESS not specified"
        exit 1
    fi
    echo "Transferring ownership to ${ORACLE_KEYSTORE_ADDRESS}"
    ORACLE_KEYSTORE_ADDRESS="$ORACLE_KEYSTORE_ADDRESS" \
    BLOCK_TRACKER_ADDRESS="$BLOCK_TRACKER_ADDRESS" \
    ORACLE_ADDRESS="$ORACLE_ADDRESS" \
    "${FORGE_BIN_PATH}" script \
    "${SCRIPT_PATH_PREFIX}"DeployScripts.s.sol:TransferOwnership \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        --chain-id "${CHAIN_ID}" \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --skip-simulation \
        --use 0.8.23 \
        --broadcast \
        --force \
        --json \
        --via-ir

elif [ "${DEPLOY_TYPE}" = "whitelist" ]; then
    if [ -z "$HYP_ERC20_ADDR" ]; then
        echo "HYP_ERC20_ADDR not specified"
        exit 1
    fi
    echo "Deploying whitelist contract"
    HYP_ERC20_ADDR="$HYP_ERC20_ADDR" "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}"DeployScripts.s.sol:DeployWhitelist \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --broadcast \
        --chain-id "${CHAIN_ID}" \
        -vvvv \
        --use 0.8.23 \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --via-ir

elif [ "${DEPLOY_TYPE}" = "settlement-gateway" ]; then
    if [ -z "$RELAYER_ADDR" ]; then
        echo "RELAYER_ADDR not specified"
        exit 1
    fi
    echo "Deploying gateway contract on settlement chain"
    RELAYER_ADDR="$RELAYER_ADDR" "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}"DeployStandardBridge.s.sol:DeploySettlementGateway \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --broadcast \
        --chain-id "${CHAIN_ID}" \
        -vvvv \
        --use 0.8.23 \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --via-ir

elif [ "${DEPLOY_TYPE}" = "l1-gateway" ]; then
    if [ -z "$RELAYER_ADDR" ]; then
        echo "RELAYER_ADDR not specified"
        exit 1
    fi
    echo "Deploying gateway contract on L1"
    RELAYER_ADDR="$RELAYER_ADDR" "${FORGE_BIN_PATH}" script \
        "${SCRIPT_PATH_PREFIX}"DeployStandardBridge.s.sol:DeployL1Gateway \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --broadcast \
        --chain-id "${CHAIN_ID}" \
        -vvvv \
        --use 0.8.23 \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --via-ir

elif [ "${DEPLOY_TYPE}" = "validator-registry" ]; then
    echo "Deploying validator registry contract"
    # Setting gas params manually ensures inclusion to mev-commit chain, recent bug in forge sets priority fee to 0.
    "${FORGE_BIN_PATH}" script \
        --priority-gas-price 2000000000 \
        --with-gas-price 5000000000 \
        "${SCRIPT_PATH_PREFIX}"DeployScripts.s.sol:DeployValidatorRegistry \
        --rpc-url "${RPC_URL}" \
        --keystores "${KEYSTORE_DIR}/${KEYSTORE_FILENAME}" \
        --password "${KEYSTORE_PASSWORD}" \
        --sender "${SENDER}" \
        --broadcast \
        --chain-id "${CHAIN_ID}" \
        -vvvv \
        --use 0.8.23 \
        --root "${CONTRACT_REPO_ROOT_PATH}" \
        --via-ir
fi 
