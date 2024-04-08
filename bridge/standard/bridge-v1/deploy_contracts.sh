#!/bin/sh

L1_CHAIN_ID=${L1_CHAIN_ID:-"17000"} # Holesky
STANDARD_BRIDGE_RELAYER_L1_RPC_URL=${STANDARD_BRIDGE_RELAYER_L1_RPC_URL:-"https://ethereum-holesky.publicnode.com"}
SETTLEMENT_CHAIN_ID=${SETTLEMENT_CHAIN_ID:-"17864"}
SETTLEMENT_DEPLOYER_PRIVKEY=${DEPLOYER_PRIVKEY:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"} # Same default deployer as core contracts

if [ -n "$FORGE_BIN_PATH" ]; then
    FORGE_BIN_PATH=$(realpath "$FORGE_BIN_PATH")
else
    FORGE_BIN_PATH="forge"
fi
if [ -n "$CAST_BIN_PATH" ]; then
    CAST_BIN_PATH=$(realpath "$CAST_BIN_PATH")
else
    CAST_BIN_PATH="cast"
fi

CONTRACTS_PATH=${CONTRACTS_PATH:-"$HOME/.primev/contracts"}
if [ ! -d "$CONTRACTS_PATH" ]; then
    echo "Error: Contracts path not found at $CONTRACTS_PATH. Please ensure the contracts are installed and the path is correct."
    exit 1
fi

ARTIFACT_OUT_PATH=${ARTIFACT_OUT_PATH:-"$HOME/.primev/contracts"}
ARTIFACT_OUT_PATH=$(realpath "$ARTIFACT_OUT_PATH")

fail_if_not_set() {
    if [ -z "$1" ]; then
        echo "Error: Required environment variable not set (one of STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL, RELAYER_PRIVKEY)"
        exit 1
    fi
}
fail_if_not_set "${STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL}"
fail_if_not_set "${RELAYER_PRIVKEY}"

RELAYER_ADDR=$("$CAST_BIN_PATH" wallet address "$RELAYER_PRIVKEY")

check_chain_id() {
    RPC_URL="$1"
    EXPECTED_CHAIN_ID="$2"
    RETRIEVED_CHAIN_ID=$("$CAST_BIN_PATH" chain-id --rpc-url "$RPC_URL")
    if [ "$RETRIEVED_CHAIN_ID" -ne "$EXPECTED_CHAIN_ID" ]; then
        echo "Error: Chain ID mismatch for $RPC_URL. Expected: $EXPECTED_CHAIN_ID, Got: $RETRIEVED_CHAIN_ID"
        exit 1
    else
        echo "Success: Chain ID for $RPC_URL matches the expected ID: $EXPECTED_CHAIN_ID"
    fi
}

check_create2() {
    RPC_URL="$1"
    CREATE2_AADR="0x4e59b44847b379578588920ca78fbf26c0b4956c"
    CODE=$("$CAST_BIN_PATH" code --rpc-url "$RPC_URL" $CREATE2_AADR)
    if [ -z "$CODE" ] || [ "$CODE" = "0x" ]; then
        echo "Create2 proxy not deployed on $RPC_URL"
        exit 1
    else
        echo "Create2 proxy deployed on $RPC_URL"
    fi
}

check_balance() {
    RPC_URL="$1"
    ADDR="$2"
    BALANCE_WEI=$("$CAST_BIN_PATH" balance "$ADDR" --rpc-url "$RPC_URL")
    ONE_ETH_WEI="1000000000000000000"

    SUFFICIENT=$(echo "$BALANCE_WEI >= $ONE_ETH_WEI" | bc)
    if [ "$SUFFICIENT" -eq 0 ]; then
        echo "Error: $ADDR has insufficient balance on chain with RPC URL $RPC_URL. Balance: $BALANCE_WEI wei"
        exit 1
    else
        echo "Confirmed: $ADDR has sufficient balance (>= 1 ETH) on chain with RPC URL $RPC_URL. Balance: $BALANCE_WEI wei"
    fi
}

check_chain_id "$STANDARD_BRIDGE_RELAYER_L1_RPC_URL" "$L1_CHAIN_ID"
check_chain_id "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" "$SETTLEMENT_CHAIN_ID"

check_create2 "$STANDARD_BRIDGE_RELAYER_L1_RPC_URL"
check_create2 "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL"

SETTLEMENT_DEPLOYER_ADDR=$("$CAST_BIN_PATH" wallet address "$SETTLEMENT_DEPLOYER_PRIVKEY")
check_balance "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" "$SETTLEMENT_DEPLOYER_ADDR"

"$CAST_BIN_PATH" send \
    --rpc-url "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" \
    --private-key "$SETTLEMENT_DEPLOYER_PRIVKEY" \
    "$RELAYER_ADDR" \
    --value 100ether

check_balance "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" "$RELAYER_ADDR"
check_balance "$STANDARD_BRIDGE_RELAYER_L1_RPC_URL" "$RELAYER_ADDR"

# Create/fund a new L1 deployer to avoid L1Gateway contract addr collision on Holeksy
L1_DEPLOYER_PRIVKEY=$($CAST_BIN_PATH wallet new | grep 'Private key' | awk '{ print $NF }')
L1_DEPLOYER_ADDR=$($CAST_BIN_PATH wallet address "$L1_DEPLOYER_PRIVKEY")
echo "New L1 deployer to be funded by relayer: $L1_DEPLOYER_ADDR"
$CAST_BIN_PATH send \
    --rpc-url "$STANDARD_BRIDGE_RELAYER_L1_RPC_URL" \
    --private-key "$RELAYER_PRIVKEY" \
    "$L1_DEPLOYER_ADDR" \
    --value 0.5ether

EXPECTED_WHITELIST_ADDR="0x57508f0B0f3426758F1f3D63ad4935a7c9383620"
check_balance "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" "$EXPECTED_WHITELIST_ADDR"

echo "changing directory to $CONTRACTS_PATH and running deploy scripts for standard bridge"
cd "$CONTRACTS_PATH" || exit

RELAYER_ADDR="$RELAYER_ADDR" $FORGE_BIN_PATH script \
    "scripts/DeployStandardBridge.s.sol:DeploySettlementGateway" \
    --rpc-url "$STANDARD_BRIDGE_RELAYER_SETTLEMENT_RPC_URL" \
    --private-key "$SETTLEMENT_DEPLOYER_PRIVKEY" \
    --broadcast \
    --chain-id "$SETTLEMENT_CHAIN_ID" \
    -vvvv \
    --use 0.8.23 | tee deploy_sg_output.txt

awk -F"JSON_DEPLOY_ARTIFACT: " '/JSON_DEPLOY_ARTIFACT:/ {print $2}' deploy_sg_output.txt | sed '/^$/d' > SettlementGatewayArtifact.json
mv SettlementGatewayArtifact.json "$ARTIFACT_OUT_PATH"

RELAYER_ADDR="$RELAYER_ADDR" $FORGE_BIN_PATH script \
    "scripts/DeployStandardBridge.s.sol:DeployL1Gateway" \
    --rpc-url "$STANDARD_BRIDGE_RELAYER_L1_RPC_URL" \
    --private-key "$L1_DEPLOYER_PRIVKEY" \
    --broadcast \
    --chain-id "$L1_CHAIN_ID" \
    -vvvv \
    --use 0.8.23 | tee deploy_l1g_output.txt

awk -F"JSON_DEPLOY_ARTIFACT: " '/JSON_DEPLOY_ARTIFACT:/ {print $2}' deploy_l1g_output.txt | sed '/^$/d' > L1GatewayArtifact.json
mv L1GatewayArtifact.json "$ARTIFACT_OUT_PATH"

rm deploy_sg_output.txt deploy_l1g_output.txt
