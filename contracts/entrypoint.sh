#!/bin/sh

# Default to core deploy type if not specified
DEPLOY_TYPE=${DEPLOY_TYPE:-core}

# Set the forge binary path, default to 'forge' if not provided
FORGE_BIN_PATH=${FORGE_BIN_PATH:-forge}

# Define the script path prefix, default to 'scripts/' if not provided
SCRIPT_PATH_PREFIX=${SCRIPT_PATH_PREFIX:-scripts/}

# Check if CONTRACT_REPO_ROOT_PATH is set, if so, prepare the --root option
ROOT_OPTION=""
if [ -n "$CONTRACT_REPO_ROOT_PATH" ]; then
    ROOT_OPTION="--root $CONTRACT_REPO_ROOT_PATH"
fi

if [ "$DEPLOY_TYPE" = "core" ]; then
    echo "Deploying core contracts"
    $FORGE_BIN_PATH script ${SCRIPT_PATH_PREFIX}DeployScripts.s.sol:DeployScript --rpc-url "$RPC_URL" --private-key "$PRIVATE_KEY" --broadcast --chain-id "$CHAIN_ID" -vvvv --use 0.8.23 $ROOT_OPTION --via-ir

elif [ "$DEPLOY_TYPE" = "whitelist" ]; then
    if [ -z "$HYP_ERC20_ADDR" ]; then
        echo "HYP_ERC20_ADDR not specified"
        exit 1
    fi
    echo "Deploying whitelist contract"
    HYP_ERC20_ADDR="$HYP_ERC20_ADDR" $FORGE_BIN_PATH script ${SCRIPT_PATH_PREFIX}DeployScripts.s.sol:DeployWhitelist --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast --chain-id $CHAIN_ID -vvvv --use 0.8.23 $ROOT_OPTION

elif [ "$DEPLOY_TYPE" = "settlement-gateway" ]; then
    if [ -z "$RELAYER_ADDR" ]; then
        echo "RELAYER_ADDR not specified"
        exit 1
    fi
    echo "Deploying gateway contract on settlement chain"
    RELAYER_ADDR="$RELAYER_ADDR" $FORGE_BIN_PATH script ${SCRIPT_PATH_PREFIX}DeployStandardBridge.s.sol:DeploySettlementGateway --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast --chain-id $CHAIN_ID -vvvv --use 0.8.23 $ROOT_OPTION

elif [ "$DEPLOY_TYPE" = "l1-gateway" ]; then
    if [ -z "$RELAYER_ADDR" ]; then
        echo "RELAYER_ADDR not specified"
        exit 1
    fi
    echo "Deploying gateway contract on L1"
    RELAYER_ADDR="$RELAYER_ADDR" $FORGE_BIN_PATH script ${SCRIPT_PATH_PREFIX}DeployStandardBridge.s.sol:DeployL1Gateway --rpc-url $RPC_URL --private-key $PRIVATE_KEY --broadcast --chain-id $CHAIN_ID -vvvv --use 0.8.23 $ROOT_OPTION
fi
