#!/bin/bash

# Configuration
RPC_URL=${RPC_URL:-"http://localhost:8545"}
THRESHOLD_SECONDS=${THRESHOLD_SECONDS:-3600}

# Get latest block number
block_number_hex=$(curl -s -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    -H "Content-Type: application/json" ${RPC_URL} | jq -r '.result')

if [ -z "$block_number_hex" ] || [ "$block_number_hex" == "null" ]; then
    echo "❌ Failed to fetch latest block number"
    exit 1
fi

# Convert block number to decimal
block_number=$((block_number_hex))

# Get block data
block_data=$(curl -s -X POST --data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"$block_number_hex\", false],\"id\":1}" \
    -H "Content-Type: application/json" ${RPC_URL})

block_timestamp_hex=$(echo "$block_data" | jq -r '.result.timestamp')

if [ -z "$block_timestamp_hex" ] || [ "$block_timestamp_hex" == "null" ]; then
    echo "❌ Failed to fetch block timestamp for block $block_number"
    exit 1
fi

# Convert timestamp to decimal
block_timestamp=$((block_timestamp_hex))

# Get current timestamp
current_timestamp=$(date +%s)

# Calculate lag
lag=$((current_timestamp - block_timestamp))

echo "🧮 Block number: $block_number"
echo "⏱️  Block timestamp: $block_timestamp (UTC)"
echo "⏳ Current timestamp: $current_timestamp (UTC)"
echo "📊 Lag: $lag seconds"

# Check if lag exceeds threshold
if [ "$lag" -le "$THRESHOLD_SECONDS" ]; then
    echo "✅ Lag ($lag sec) is within threshold ($THRESHOLD_SECONDS sec). Marking pod as READY."
    exit 0
else
    echo "❌ Lag ($lag sec) exceeds threshold ($THRESHOLD_SECONDS sec). Marking pod as NOT READY."
    exit 1
fi
