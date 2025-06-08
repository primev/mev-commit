#!/bin/bash
set -e  # Remove the 'x' flag to reduce verbosity

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration from environment
KEYSTORE_DIR=${KEYSTORE_DIR:-"/keystore"}
KEYSTORE_PASSWORD=${KEYSTORE_PASSWORD:-""}
KEYSTORE_SECRET=${KEYSTORE_SECRET:-""}
OUTPUT_DIR=${OUTPUT_DIR:-"/shared"}

echo -e "${BLUE}Initializing MEV Oracle account...${NC}"

# Validate environment
if [ -z "$KEYSTORE_PASSWORD" ]; then
    echo -e "${RED}Error: KEYSTORE_PASSWORD environment variable must be set${NC}"
    exit 1
fi

# Create directories
mkdir -p "$KEYSTORE_DIR"
mkdir -p "$OUTPUT_DIR"

# Function to extract address from keystore
extract_address() {
    local keystore_file=$1
    grep -o '"address":"[^"]*' "$keystore_file" | cut -d'"' -f4
}

# Check if we're using an existing secret
if [ -n "$KEYSTORE_SECRET" ] && [ -d "/existing-keystore" ]; then
    echo -e "${YELLOW}Using existing keystore from secret...${NC}"
    
    # Copy keystore files from secret mount to working directory
    cp /existing-keystore/* "$KEYSTORE_DIR/" 2>/dev/null || true
    
    # Find the keystore file
    KEYSTORE_FILE=$(find "$KEYSTORE_DIR" -type f -name "UTC--*" | head -n 1)
    
    if [ -z "$KEYSTORE_FILE" ]; then
        echo -e "${RED}Error: No keystore file found in secret${NC}"
        exit 1
    fi
    
    ADDRESS=$(extract_address "$KEYSTORE_FILE")
    echo -e "${GREEN}Using existing account: 0x${ADDRESS}${NC}"
else
    echo -e "${YELLOW}Generating new keystore...${NC}"
    
    # Check if geth is available
    if ! command -v geth &> /dev/null; then
        echo -e "${RED}Error: geth is not installed in init container${NC}"
        exit 1
    fi
    
    # Create password file
    echo "$KEYSTORE_PASSWORD" > /tmp/password.txt
    chmod 600 /tmp/password.txt
    
    # Generate new account
    echo -e "${YELLOW}Creating new Ethereum account...${NC}"
    geth account new --keystore "$KEYSTORE_DIR" --password /tmp/password.txt > /tmp/geth_output.txt 2>&1
    
    # Clean up password file
    rm -f /tmp/password.txt
    
    # Find the generated keystore file
    KEYSTORE_FILE=$(find "$KEYSTORE_DIR" -type f -name "UTC--*" | head -n 1)
    
    if [ -z "$KEYSTORE_FILE" ]; then
        echo -e "${RED}Error: Failed to generate keystore${NC}"
        exit 1
    fi
    
    ADDRESS=$(extract_address "$KEYSTORE_FILE")
    echo -e "${GREEN}Generated new account: 0x${ADDRESS}${NC}"
fi

# Set permissions
chmod -R 600 "$KEYSTORE_DIR"/*

# Share keystore with main container
if [ "$OUTPUT_DIR" != "$KEYSTORE_DIR" ]; then
    echo -e "${YELLOW}Copying keystore to shared volume...${NC}"
    cp -r "$KEYSTORE_DIR"/* "$OUTPUT_DIR/"
fi

# Create info file for main container
cat > "$OUTPUT_DIR/account-info.txt" <<EOF
ADDRESS=0x${ADDRESS}
KEYSTORE_FILE=$(basename "$KEYSTORE_FILE")
INITIALIZED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
EOF

echo -e "${GREEN}Account initialization complete!${NC}"
echo -e "${BLUE}Address:${NC} 0x${ADDRESS}"
echo -e "${BLUE}Keystore:${NC} $(basename "$KEYSTORE_FILE")"

# Keep the init container running briefly to ensure files are written
sleep 2
