#!/bin/bash

# Define the base directory where the .json files are located
BASE_DIR="../contracts"

# Create the abi directory if it doesn't exist
ABI_DIR="./abi"
mkdir -p "$ABI_DIR"

forge compile --root $BASE_DIR --via-ir

# Function to extract and save the ABI
extract_and_save_abi() {
    local json_file="$1"
    local abi_file="$2"
    jq .abi "$json_file" > "$abi_file"
}

# Extract ABI for BidderRegistry.json
extract_and_save_abi "$BASE_DIR/out/BidderRegistry.sol/BidderRegistry.json" "$ABI_DIR/BidderRegistry.abi"

# Extract ABI for ProviderRegistry.json
extract_and_save_abi "$BASE_DIR/out/ProviderRegistry.sol/ProviderRegistry.json" "$ABI_DIR/ProviderRegistry.abi"

# Extract ABI for Oracle.json
extract_and_save_abi "$BASE_DIR/out/Oracle.sol/Oracle.json" "$ABI_DIR/Oracle.abi"

# Extract ABI for PreconfManager.json
extract_and_save_abi "$BASE_DIR/out/PreconfManager.sol/PreconfManager.json" "$ABI_DIR/PreconfManager.abi"

# Extract ABI for SettlementGateway.json
extract_and_save_abi "$BASE_DIR/out/SettlementGateway.sol/SettlementGateway.json" "$ABI_DIR/SettlementGateway.abi"

# Extract ABI for L1Gateway.json
extract_and_save_abi "$BASE_DIR/out/L1Gateway.sol/L1Gateway.json" "$ABI_DIR/L1Gateway.abi"

extract_and_save_abi "$BASE_DIR/out/ValidatorRegistryV1.sol/ValidatorRegistryV1.json" "$ABI_DIR/ValidatorRegistryV1.abi"

extract_and_save_abi "$BASE_DIR/out/BlockTracker.sol/BlockTracker.json" "$ABI_DIR/BlockTracker.abi"

extract_and_save_abi "$BASE_DIR/out/MevCommitAVS.sol/MevCommitAVS.json" "$ABI_DIR/MevCommitAVS.abi"

extract_and_save_abi "$BASE_DIR/out/ValidatorOptInRouter.sol/ValidatorOptInRouter.json" "$ABI_DIR/ValidatorOptInRouter.abi"

echo "ABI files extracted successfully."


ABI_DIR="./abi"
GO_CODE_BASE_DIR="./clients"

# Create the Go code base directory if it doesn't exist
mkdir -p "$GO_CODE_BASE_DIR"

# Function to generate Go code from ABI and place it in a separate folder
generate_go_code() {
    local abi_file="$1"
    local contract_name="$2"
    local pkg_name="$3"

    # Create a directory for the contract
    local contract_dir="$GO_CODE_BASE_DIR/$contract_name"
    mkdir -p "$contract_dir"

    # Run abigen and output the Go code in the contract's directory
    abigen --abi "$abi_file" --pkg "$pkg_name" --out "$contract_dir/$contract_name.go"
}

# Generate Go code for BidderRegistry.abi
generate_go_code "$ABI_DIR/BidderRegistry.abi" "BidderRegistry" "bidderregistry"

# Generate Go code for ProviderRegistry.abi
generate_go_code "$ABI_DIR/ProviderRegistry.abi" "ProviderRegistry" "providerregistry"

# Generate Go code for Oracle.abi
generate_go_code "$ABI_DIR/Oracle.abi" "Oracle" "oracle"

# Generate Go code for PreconfManager.abi
generate_go_code "$ABI_DIR/PreconfManager.abi" "PreconfManager" "preconfmanager"

# Generate Go code for SettlementGateway.abi
generate_go_code "$ABI_DIR/SettlementGateway.abi" "SettlementGateway" "settlementgateway"

# Generate Go code for L1Gateway.abi
generate_go_code "$ABI_DIR/L1Gateway.abi" "L1Gateway" "l1gateway"

generate_go_code "$ABI_DIR/ValidatorRegistryV1.abi" "ValidatorRegistryV1" "validatorregistryv1"

generate_go_code "$ABI_DIR/BlockTracker.abi" "BlockTracker" "blocktracker"

generate_go_code "$ABI_DIR/MevCommitAVS.abi" "MevCommitAVS" "mevcommitavs"

generate_go_code "$ABI_DIR/ValidatorOptInRouter.abi" "ValidatorOptInRouter" "validatoroptinrouter"

echo "Go code generated successfully in separate folders."

