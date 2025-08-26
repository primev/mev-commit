#!/bin/sh

set -e

echo "Starting account generation..."

# Create keystore directory if it doesn't exist
mkdir -p /keystore

# Create accounts file to store generated addresses
ACCOUNTS_FILE="/shared/generated_accounts.txt"
> "$ACCOUNTS_FILE"  # Clear the file

echo "Generating $ACCOUNT_COUNT accounts..."

for i in $(seq 1 $ACCOUNT_COUNT); do
    echo "Generating account $i of $ACCOUNT_COUNT..."
    
    # Generate account using geth
    ACCOUNT_OUTPUT=$(geth account new \
        --keystore /keystore \
        --password /password/password.txt 2>&1)
    
    # Extract address from output
    # Geth outputs: "Your new account is locked with a password. Please give a password. Do not forget this password.
    # Your new key was generated
    # Public address of the key:   0x..."
    ADDRESS=$(echo "$ACCOUNT_OUTPUT" | grep -i "Public address of the key:" | awk '{print $NF}' | tr -d '\n\r')
    
    if [ -z "$ADDRESS" ]; then
        # Fallback: try to extract from different output format
        ADDRESS=$(echo "$ACCOUNT_OUTPUT" | grep -o "0x[a-fA-F0-9]\{40\}" | head -1)
    fi
    
    if [ -z "$ADDRESS" ]; then
        echo "Error: Could not extract address from geth output"
        echo "Geth output was: $ACCOUNT_OUTPUT"
        exit 1
    fi
    
    echo "Generated account: $ADDRESS"
    echo "$ADDRESS" >> "$ACCOUNTS_FILE"
done

echo "Generated $ACCOUNT_COUNT accounts successfully!"
echo "Accounts stored in: $ACCOUNTS_FILE"

# List generated keystore files
echo "Keystore files:"
ls -la /keystore/

# Show generated accounts
echo "Generated account addresses:"
cat "$ACCOUNTS_FILE"
