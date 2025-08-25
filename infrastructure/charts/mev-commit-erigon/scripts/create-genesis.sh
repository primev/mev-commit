#!/bin/sh

set -e

echo "Creating genesis.json with generated accounts..."

# Check if accounts were generated
ACCOUNTS_FILE="/shared/generated_accounts.txt"
if [ ! -f "$ACCOUNTS_FILE" ]; then
    echo "Error: No generated accounts found at $ACCOUNTS_FILE"
    exit 1
fi

# Read the genesis template
GENESIS_TEMPLATE="/genesis-template/genesis-template.json"
if [ ! -f "$GENESIS_TEMPLATE" ]; then
    echo "Error: Genesis template not found at $GENESIS_TEMPLATE"
    exit 1
fi

# Copy template to working location
cp "$GENESIS_TEMPLATE" /shared/genesis.json

echo "Adding generated accounts to genesis..."

# Install jq for JSON manipulation
apk add --no-cache jq

# Read accounts and add them to genesis
while IFS= read -r address; do
    if [ -n "$address" ]; then
        echo "Adding account $address with balance $ACCOUNT_BALANCE"
        
        # Add account to alloc section
        jq --arg addr "$address" --arg balance "$ACCOUNT_BALANCE" \
           '.alloc[$addr] = {"balance": $balance}' \
           /shared/genesis.json > /shared/genesis_temp.json && \
           mv /shared/genesis_temp.json /shared/genesis.json
    fi
done < "$ACCOUNTS_FILE"

echo "Genesis file created successfully!"
echo "Genesis file location: /shared/genesis.json"

# Show the final genesis file (optional, for debugging)
echo "Final genesis.json alloc section:"
cat /shared/genesis.json | jq '.alloc'
