#!/bin/sh
set -exu

sleep 10

# Construct .dogrc file from env vars
cat > /.dogrc <<EOF
[Connection]
apikey = $DD_API_KEY
appkey = $DD_APP_KEY
EOF

echo ".dogrc file created successfully."

# Fail script if no warp deployment file is found
if [ ! -f /deploy-artifacts/warp-deployment.json ]; then
    echo "Error: warp-deployment.json not found. Please deploy bridge."
    exit 1
fi

# Init bridge client
SEPOLIA_ROUTER=$(cat /deploy-artifacts/warp-deployment.json | jq -r '.sepolia.router')
MEV_COMMIT_CHAIN_ROUTER=$(cat /deploy-artifacts/warp-deployment.json | jq -r '.mevcommitsettlement.router')
SEPOLIA_CHAIN_ID=11155111
MEV_COMMIT_CHAIN_ID=17864
SEPOLIA_URL=${SEPOLIA_RPC_URL}
MEV_COMMIT_CHAIN_URL=${SETTLEMENT_RPC_URL}

# Ensure balances on both chains are above 1 ETH
L1_BALANCE=$(cast balance --rpc-url $SEPOLIA_URL $EMULATOR_ADDRESS)
MEV_COMMIT_BALANCE=$(cast balance --rpc-url $MEV_COMMIT_CHAIN_URL $EMULATOR_ADDRESS)
MIN_BALANCE="1000000000000000000"  # 1.0 ether in wei
if [ "$(echo "$L1_BALANCE < $MIN_BALANCE" | bc)" -eq 1 ]; then
    echo "$EMULATOR_ADDRESS must be funded with at least 1.0 ether on Sepolia."
    exit 1
fi
if [ "$(echo "$MEV_COMMIT_BALANCE < $MIN_BALANCE" | bc)" -eq 1 ]; then
    echo "$EMULATOR_ADDRESS must be funded with at least 1.0 ether on mev-commit chain."
    exit 1
fi

bridge-cli init \
    ${SEPOLIA_ROUTER} ${MEV_COMMIT_CHAIN_ROUTER} \
    ${SEPOLIA_CHAIN_ID} ${MEV_COMMIT_CHAIN_ID} \
    ${SEPOLIA_URL} ${MEV_COMMIT_CHAIN_URL} \
    --yes

function bridge_and_post_metric() {
    SUB_CMD=$1
    CHAIN_ID=$2
    AMOUNT=$3

    start_time=$(date +%s.%N)
    output=$(bridge-cli $SUB_CMD $AMOUNT $EMULATOR_ADDRESS $EMULATOR_PRIVATE_KEY --yes 2>&1)
    end_time=$(date +%s.%N)
    elapsed_time=$(echo "$end_time - $start_time" | bc)

    if echo "$output" | grep -q "SUCCESS"; then
        echo "Bridged $AMOUNT ether to Chain $CHAIN_ID successfully."
        dog --config /.dogrc metric post bridging.success $elapsed_time --tags="account_addr:$EMULATOR_ADDRESS,to_chain_id:$CHAIN_ID"
    elif echo "$output" | grep -q "FAILURE"; then
        echo "Failed to bridge $AMOUNT ether to Chain $CHAIN_ID."
        dog --config /.dogrc metric post bridging.failure $elapsed_time --tags="account_addr:$EMULATOR_ADDRESS,to_chain_id:$CHAIN_ID"
    elif echo "$output" | grep -q "EXPIRED"; then
        echo "Failed to bridge $AMOUNT ether to Chain $CHAIN_ID, source transaction is stuck pending."
        dog --config /.dogrc metric post bridging.expired $elapsed_time --tags="account_addr:$EMULATOR_ADDRESS,to_chain_id:$CHAIN_ID"
    else
        echo "Unknown bridge result: $output"
        exit 1
    fi
    echo "Elapsed time: $elapsed_time"
    echo "Metric posted at: $(date)"
}

while true; do
    # Generate a random amount between 0 and 1,000,000,000,000 wei
    RANDOM_AMOUNT=$(( RANDOM % 1000000000001 ))

    bridge_and_post_metric "bridge-to-mev-commit" $MEV_COMMIT_CHAIN_ID $RANDOM_AMOUNT
    bridge_and_post_metric "bridge-to-l1" $SEPOLIA_CHAIN_ID $RANDOM_AMOUNT
    sleep 10
done
