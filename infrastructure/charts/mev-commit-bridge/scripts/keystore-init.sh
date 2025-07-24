#!/bin/sh

set -e

# Configuration from environment variables
KEYSTORE_PATH="${KEYSTORE_PATH:-/keystore}"
MAX_RETRIES="${KEYSTORE_RETRIES:-3}"

echo "=== MEV-Commit Keystore Initialization ==="
echo "Keystore path: $KEYSTORE_PATH"

# Create keystore directory if it doesn't exist
mkdir -p "$KEYSTORE_PATH"

# Function to extract address from filename using POSIX shell
extract_address_from_filename() {
    local filename="$1"
    
    # Check if filename matches UTC--<timestamp>--<address> format
    case "$filename" in
        UTC--????-??-??T??-??-??.??????*Z--*)
            # Extract the address part after the last --
            address_part="${filename##*--}"
            # Validate address is 40 hex characters
            case "$address_part" in
                *[!0-9a-fA-F]* | ????????????????????????????????????)
                    echo "Invalid address format: $address_part" >&2
                    return 1
                    ;;
                ????????????????????????????????????????)
                    echo "$address_part"
                    return 0
                    ;;
                *)
                    echo "Invalid address length: $address_part" >&2
                    return 1
                    ;;
            esac
            ;;
        *)
            echo "Invalid keystore filename format: $filename" >&2
            return 1
            ;;
    esac
}

# Function to extract address from keystore JSON content
extract_address_from_content() {
    local file_path="$1"
    
    if command -v jq >/dev/null 2>&1; then
        jq -r '.address' "$file_path" 2>/dev/null || echo ""
    else
        # Extract address using grep and sed (fallback)
        grep -o '"address"[[:space:]]*:[[:space:]]*"[^"]*"' "$file_path" 2>/dev/null | \
        sed 's/.*"address"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/' || echo ""
    fi
}

# Function to normalize address (remove 0x prefix and convert to lowercase)
normalize_address() {
    local addr="$1"
    # Remove 0x prefix if present and convert to lowercase
    echo "${addr#0x}" | tr '[:upper:]' '[:lower:]'
}

# Function to validate keystore file
validate_keystore() {
    local file_path="$1"
    local filename="$2"
    
    echo "Validating keystore file: $file_path"
    
    if [ ! -f "$file_path" ]; then
        echo "Keystore file does not exist"
        return 1
    fi
    
    # Check if file is readable and not empty
    if [ ! -s "$file_path" ]; then
        echo "Keystore file is empty"
        return 1
    fi
    
    # Extract address from filename
    filename_address=$(extract_address_from_filename "$filename")
    if [ $? -ne 0 ]; then
        echo "✗ Invalid filename format"
        return 1
    fi
    
    # Normalize filename address
    filename_address=$(normalize_address "$filename_address")
    echo "Address from filename: $filename_address"
    
    # Check basic JSON structure first
    if command -v jq >/dev/null 2>&1; then
        echo "Using jq for JSON validation"
        
        # Validate JSON structure
        if ! jq -e '.address and .crypto and .crypto.cipher and .crypto.ciphertext and .crypto.kdf' "$file_path" >/dev/null 2>&1; then
            echo "✗ Invalid keystore JSON structure"
            return 1
        fi
        
        # Extract address from content
        content_address=$(extract_address_from_content "$file_path")
        if [ -z "$content_address" ]; then
            echo "✗ Could not extract address from keystore content"
            return 1
        fi
        
        # Normalize content address
        content_address=$(normalize_address "$content_address")
        echo "Address from content: $content_address"
        
        # Compare addresses (case-insensitive)
        if [ "$filename_address" = "$content_address" ]; then
            echo "✓ Keystore validation successful - addresses match: $filename_address"
            return 0
        else
            echo "✗ Address mismatch - filename: $filename_address, content: $content_address"
            return 1
        fi
    else
        echo "Using basic validation (jq not available)"
        
        # Basic validation without jq
        if ! grep -q '"crypto"' "$file_path" || \
           ! grep -q '"cipher"' "$file_path" || \
           ! grep -q '"ciphertext"' "$file_path" || \
           ! grep -q '"kdf"' "$file_path"; then
            echo "✗ Basic keystore structure validation failed"
            return 1
        fi
        
        # Extract address from content using basic tools
        content_address=$(extract_address_from_content "$file_path")
        if [ -z "$content_address" ]; then
            echo "✗ Could not extract address from keystore content"
            return 1
        fi
        
        # Normalize content address
        content_address=$(normalize_address "$content_address")
        echo "Address from content: $content_address"
        
        # Compare addresses (case-insensitive)
        if [ "$filename_address" = "$content_address" ]; then
            echo "✓ Basic keystore validation successful - addresses match: $filename_address"
            return 0
        else
            echo "✗ Address mismatch - filename: $filename_address, content: $content_address"
            return 1
        fi
    fi
}

# Main logic - Read keystore content and filename from External Secret
KEYSTORE_CONTENT_FILE="/secrets/temp_keystore.json"
FILENAME_FILE="/secrets/filename.txt"

echo "Reading keystore data from External Secret..."

# Validate that secret files exist
if [ ! -f "$KEYSTORE_CONTENT_FILE" ]; then
    echo "✗ Keystore content file not found: $KEYSTORE_CONTENT_FILE"
    exit 1
fi

if [ ! -f "$FILENAME_FILE" ]; then
    echo "✗ Filename file not found: $FILENAME_FILE"
    exit 1
fi

# Read the expected filename
EXPECTED_FILENAME=$(cat "$FILENAME_FILE")
KEYSTORE_FILE_PATH="$KEYSTORE_PATH/$EXPECTED_FILENAME"

echo "Expected filename: $EXPECTED_FILENAME"

# Validate filename format before proceeding
if ! extract_address_from_filename "$EXPECTED_FILENAME" >/dev/null 2>&1; then
    echo "✗ Invalid filename format: $EXPECTED_FILENAME"
    echo "Expected format: UTC--<timestamp>--<address>"
    exit 1
fi

# Check if keystore already exists and is valid
if [ -f "$KEYSTORE_FILE_PATH" ]; then
    if validate_keystore "$KEYSTORE_FILE_PATH" "$EXPECTED_FILENAME"; then
        echo "✓ Valid keystore already exists - initialization complete"
        exit 0
    else
        echo "Removing invalid keystore file..."
        rm -f "$KEYSTORE_FILE_PATH"
    fi
fi

# Copy keystore content to final location
echo "Setting up keystore from External Secret..."
cp "$KEYSTORE_CONTENT_FILE" "$KEYSTORE_FILE_PATH"

# Validate the keystore
if ! validate_keystore "$KEYSTORE_FILE_PATH" "$EXPECTED_FILENAME"; then
    echo "✗ Keystore validation failed"
    rm -f "$KEYSTORE_FILE_PATH"
    exit 1
fi

# Set proper permissions
chmod 600 "$KEYSTORE_FILE_PATH"

# Create a checksum file for integrity checking
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$KEYSTORE_FILE_PATH" > "$KEYSTORE_PATH/.keystore.checksum"
    echo "✓ Checksum file created"
fi

# Extract final address for logging
FINAL_ADDRESS=$(extract_address_from_filename "$EXPECTED_FILENAME")
echo "✓ Keystore initialization completed successfully"
echo "✓ Keystore file: $KEYSTORE_FILE_PATH"
echo "✓ Address: $FINAL_ADDRESS"
