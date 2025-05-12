#!/bin/sh
# init-keystore.sh

# Script: init-keystore.sh
# Purpose: Initialize signer node using either keystore file or private key import

# Configuration with defaults
DATADIR="${DATADIR:-/geth/data}"
KEYSTORE_DIR="$DATADIR/keystore"
PASSWORD_FILE="${PASSWORD_FILE:-/geth/password/password.txt}"
SIGNER_KEY_FILE="${SIGNER_KEY_FILE:-/geth/keys/signer-key}"

# Required parameters
NODE_BOOTSTRAP_METHOD="${NODE_BOOTSTRAP_METHOD:-}"
KEYSTORE_URL="${KEYSTORE_URL:-}"
SIGNER_ADDRESS="${SIGNER_ADDRESS:-}"  # Get signer address from environment

echo "🚀 Initializing signer node..."
echo "📋 Configuration:"
echo "   - Bootstrap method: $NODE_BOOTSTRAP_METHOD"
echo "   - Data directory: $DATADIR"
echo "   - Keystore directory: $KEYSTORE_DIR"
echo "   - Password file: $PASSWORD_FILE"
echo "   - Signer address: $SIGNER_ADDRESS"

# Ensure keystore directory exists
echo "📁 Ensuring keystore directory exists..."
mkdir -p "$KEYSTORE_DIR"

# Function to check if account already exists
account_exists() {
    local address=$1
    # Remove 0x prefix if present
    address=$(echo "$address" | sed 's/^0x//i')
    
    # Check if any keystore file contains this address
    for keyfile in "$KEYSTORE_DIR"/UTC--*; do
        if [ -f "$keyfile" ]; then
            # Extract address from filename (last part after --)
            file_address=$(basename "$keyfile" | sed 's/.*--//')
            if [ "$(echo "$file_address" | tr '[:upper:]' '[:lower:]')" = "$(echo "$address" | tr '[:upper:]' '[:lower:]')" ]; then
                echo "   ✅ Account $address already exists in keystore"
                return 0
            fi
        fi
    done
    
    return 1
}

# Process based on bootstrap method
case "$NODE_BOOTSTRAP_METHOD" in
  "signerKeystore")
    echo "📥 Using signer keystore method..."
    echo "   - Keystore URL: $KEYSTORE_URL"

    # Validate required environment variable
    if [ -z "$KEYSTORE_URL" ]; then
      echo "❌ ERROR: KEYSTORE_URL is required for signerKeystore method!"
      exit 1
    fi

    # Check if keystore already exists
    if [ -n "$SIGNER_ADDRESS" ] && account_exists "$SIGNER_ADDRESS"; then
      echo "⏭️ Keystore already exists, skipping download"
    else
      # Download keystore file
      echo "⬇️  Downloading keystore..."
      FILENAME=$(basename "$KEYSTORE_URL")
      TARGET_FILE="$KEYSTORE_DIR/$FILENAME"

      if command -v wget > /dev/null 2>&1; then
        wget -O "$TARGET_FILE" "$KEYSTORE_URL"
      elif command -v curl > /dev/null 2>&1; then
        curl -sSL -o "$TARGET_FILE" "$KEYSTORE_URL"
      else
        echo "❌ ERROR: Neither wget nor curl found!"
        exit 1
      fi

      if [ $? -eq 0 ]; then
        echo "✅ Keystore successfully downloaded to $TARGET_FILE"
        echo "📄 File details: $(ls -la "$TARGET_FILE")"
      else
        echo "❌ Failed to download keystore file!"
        exit 1
      fi
    fi
    ;;

  "signerPrivateKey")
    echo "🔑 Using signer private key method..."
    echo "   - Private key file: $SIGNER_KEY_FILE"

    # Validate required files
    if [ ! -f "$SIGNER_KEY_FILE" ]; then
      echo "❌ ERROR: Signer key not found at $SIGNER_KEY_FILE!"
      exit 1
    fi

    if [ ! -f "$PASSWORD_FILE" ]; then
      echo "❌ ERROR: Password file not found at $PASSWORD_FILE!"
      exit 1
    fi

    # Check if account already exists
    if [ -n "$SIGNER_ADDRESS" ] && account_exists "$SIGNER_ADDRESS"; then
      echo "⏭️ Account already imported, skipping import"
    else
      # Import private key
      echo "🔐 Importing private key..."
      geth --datadir "$DATADIR" account import --password "$PASSWORD_FILE" "$SIGNER_KEY_FILE"

      if [ $? -eq 0 ]; then
        echo "✅ Private key successfully imported"
        echo "📂 Keystore contents:"
        ls -la "$KEYSTORE_DIR"
      else
        echo "❌ Failed to import private key!"
        exit 1
      fi
    fi
    ;;

  *)
    echo "❌ ERROR: Invalid or missing NODE_BOOTSTRAP_METHOD"
    echo "   Valid options: 'signerKeystore' or 'signerPrivateKey'"
    echo "   Current value: '$NODE_BOOTSTRAP_METHOD'"
    exit 1
    ;;
esac

echo "🎉 Signer node initialization complete!"
