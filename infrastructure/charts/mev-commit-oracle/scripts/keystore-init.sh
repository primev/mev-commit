#!/bin/sh

set -e

# Configuration from environment variables
KEYSTORE_PATH="${KEYSTORE_PATH:-/keystore}"
KEYSTORE_URL="${KEYSTORE_DOWNLOAD_URL}"
MAX_RETRIES="${KEYSTORE_RETRIES:-3}"
KEYSTORE_SOURCE="${KEYSTORE_SOURCE:-url}"  # 'url' or 'aws'

echo "=== MEV Oracle Keystore Initialization ==="
echo "Keystore path: $KEYSTORE_PATH"
echo "Keystore source: $KEYSTORE_SOURCE"

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

# Function to download keystore with retries (URL mode)
download_keystore() {
    local url="$1"
    local output_path="$2"
    local retries=0
    
    echo "Downloading keystore from: $url"
    
    # Download with retries
    while [ $retries -lt $MAX_RETRIES ]; do
        echo "Download attempt $((retries + 1))/$MAX_RETRIES"
        
        if curl -fsSL --connect-timeout 30 --max-time 300 -o "$output_path" "$url"; then
            echo "✓ Download successful"
            return 0
        else
            retries=$((retries + 1))
            if [ $retries -lt $MAX_RETRIES ]; then
                echo "Download failed, retrying in 5 seconds..."
                sleep 5
            else
                echo "✗ Download failed after $MAX_RETRIES attempts"
                return 1
            fi
        fi
    done
}

# Function to get filename from URL
get_filename_from_url() {
    local url="$1"
    basename "$url"
}

# Function to handle AWS Secrets Manager source
handle_aws_source() {
    echo "Using AWS Secrets Manager source..."
    
    # Read keystore content and filename from External Secret
    local keystore_content_file="/secrets/temp_keystore.json"
    local filename_file="/secrets/filename.txt"
    
    # Validate that secret files exist
    if [ ! -f "$keystore_content_file" ]; then
        echo "✗ Keystore content file not found: $keystore_content_file"
        exit 1
    fi
    
    if [ ! -f "$filename_file" ]; then
        echo "✗ Filename file not found: $filename_file"
        exit 1
    fi
    
    # Read the expected filename
    local expected_filename=$(cat "$filename_file")
    local keystore_file_path="$KEYSTORE_PATH/$expected_filename"
    
    echo "Expected filename: $expected_filename"
    
    # Validate filename format before proceeding
    if ! extract_address_from_filename "$expected_filename" >/dev/null 2>&1; then
        echo "✗ Invalid filename format: $expected_filename"
        echo "Expected format: UTC--<timestamp>--<address>"
        exit 1
    fi
    
    # Check if keystore already exists and is valid
    if [ -f "$keystore_file_path" ]; then
        if validate_keystore "$keystore_file_path" "$expected_filename"; then
            echo "✓ Valid keystore already exists - initialization complete"
            return 0
        else
            echo "Removing invalid keystore file..."
            rm -f "$keystore_file_path"
        fi
    fi
    
    # Copy keystore content to final location
    echo "Setting up keystore from External Secret..."
    cp "$keystore_content_file" "$keystore_file_path"
    
    # Validate the keystore
    if ! validate_keystore "$keystore_file_path" "$expected_filename"; then
        echo "✗ Keystore validation failed"
        rm -f "$keystore_file_path"
        exit 1
    fi
    
    # Set proper permissions
    chmod 600 "$keystore_file_path"
    
    # Create a checksum file for integrity checking
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$keystore_file_path" > "$KEYSTORE_PATH/.keystore.checksum"
        echo "✓ Checksum file created"
    fi
    
    # Extract final address for logging
    local final_address=$(extract_address_from_filename "$expected_filename")
    echo "✓ AWS keystore initialization completed successfully"
    echo "✓ Keystore file: $keystore_file_path"
    echo "✓ Address: $final_address"
}

# Function to handle URL download source
handle_url_source() {
    echo "Using URL download source..."
    
    # Validate required environment variables
    if [ -z "$KEYSTORE_URL" ]; then
        echo "Error: KEYSTORE_DOWNLOAD_URL is required for URL source"
        exit 1
    fi
    
    echo "Download URL: $KEYSTORE_URL"
    
    local temp_download="/tmp/keystore_temp.json"
    
    # Get expected filename from URL
    local expected_filename=$(get_filename_from_url "$KEYSTORE_URL")
    local keystore_file_path="$KEYSTORE_PATH/$expected_filename"
    
    echo "Expected filename: $expected_filename"
    
    # Validate filename format before proceeding
    if ! extract_address_from_filename "$expected_filename" >/dev/null 2>&1; then
        echo "✗ Invalid filename format in URL: $expected_filename"
        echo "Expected format: UTC--<timestamp>--<address>"
        exit 1
    fi
    
    # Check if keystore already exists and is valid
    if [ -f "$keystore_file_path" ]; then
        if validate_keystore "$keystore_file_path" "$expected_filename"; then
            echo "✓ Valid keystore already exists - initialization complete"
            return 0
        else
            echo "Removing invalid keystore file..."
            rm -f "$keystore_file_path"
        fi
    fi
    
    # Download keystore to temporary location first
    if ! download_keystore "$KEYSTORE_URL" "$temp_download"; then
        echo "✗ Failed to download keystore"
        exit 1
    fi
    
    # Validate downloaded keystore before moving it to final location
    if ! validate_keystore "$temp_download" "$expected_filename"; then
        echo "✗ Downloaded keystore validation failed"
        rm -f "$temp_download"
        exit 1
    fi
    
    # Move validated keystore to final location
    mv "$temp_download" "$keystore_file_path"
    
    # Set proper permissions
    chmod 600 "$keystore_file_path"
    
    # Create a checksum file for integrity checking
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$keystore_file_path" > "$KEYSTORE_PATH/.keystore.checksum"
        echo "✓ Checksum file created"
    fi
    
    # Extract final address for logging
    local final_address=$(extract_address_from_filename "$expected_filename")
    echo "✓ URL keystore initialization completed successfully"
    echo "✓ Keystore file: $keystore_file_path"
    echo "✓ Address: $final_address"
}

# Main logic - route to appropriate handler based on source
case "$KEYSTORE_SOURCE" in
    "aws")
        handle_aws_source
        ;;
    "url")
        handle_url_source
        ;;
    *)
        echo "✗ Invalid KEYSTORE_SOURCE: $KEYSTORE_SOURCE"
        echo "Valid values: 'aws' or 'url'"
        exit 1
        ;;
esac
