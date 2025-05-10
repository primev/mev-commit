#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

NODE_TYPE=${NODE_TYPE:-"unknown"}
KEYSTORE_DIR=${KEYSTORE_DIR:-"/keystore"}
KEYSTORE_PASSWORD=${KEYSTORE_PASSWORD:-"changeme"}

echo -e "${BLUE}Initializing account for MEV-Commit ${NODE_TYPE} node${NC}"

# Make sure directory exists
mkdir -p "$KEYSTORE_DIR"

# Check if keystore already exists
existing_files=$(find "$KEYSTORE_DIR" -type f -name "UTC--*" | wc -l)

if [ "$existing_files" -gt 0 ]; then
    echo -e "${YELLOW}Keystore already exists in $KEYSTORE_DIR. Skipping account creation.${NC}"
    # Print the address from the existing keystore
    KEYSTORE_FILE=$(find "$KEYSTORE_DIR" -type f -name "UTC--*" | head -n 1)
    if [ -n "$KEYSTORE_FILE" ]; then
        ADDRESS=$(grep -o '"address":"[^"]*' "$KEYSTORE_FILE" | cut -d'"' -f4)
        echo -e "${GREEN}Using existing account with address: 0x${ADDRESS}${NC}"
    else
        echo -e "${RED}Failed to read address from existing keystore.${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}No existing keystore found. Creating new account...${NC}"
    
    # Save password to a temporary file with secure permissions
    echo "$KEYSTORE_PASSWORD" > /tmp/password.txt
    chmod 600 /tmp/password.txt
    
    # Generate a new account using Geth
    echo -e "${YELLOW}Generating new Ethereum account...${NC}"
    geth account new --keystore "$KEYSTORE_DIR" --password /tmp/password.txt

    # Verify keystore was created
    KEYSTORE_FILE=$(find "$KEYSTORE_DIR" -type f -name "UTC--*" | head -n 1)
    if [ -z "$KEYSTORE_FILE" ]; then
        echo -e "${RED}Failed to generate keystore.${NC}"
        rm -f /tmp/password.txt
        exit 1
    fi

    # Extract the address from the file
    ADDRESS=$(grep -o '"address":"[^"]*' "$KEYSTORE_FILE" | cut -d'"' -f4)
    echo -e "${GREEN}Generated account with address: 0x${ADDRESS}${NC}"
    
    # Clean up
    rm -f /tmp/password.txt
fi

# Set permissions on keystore
chmod -R 600 "$KEYSTORE_DIR"/*
echo -e "${GREEN}Keystore permissions set.${NC}"

echo -e "${GREEN}Account initialization complete!${NC}"
echo "ADDRESS=0x${ADDRESS}" > /account-info/address.env
