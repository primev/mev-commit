#!/usr/bin/env bash

# Script to validate contract upgrade safety using OpenZeppelin's upgrades-core
# This script wraps the npx @openzeppelin/upgrades-core validate command

set -e

contract=""
reference=""
build_info_path="out/build-info"
skip_build=false

show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --contract <CONTRACT>       Name of the new contract to validate (e.g., MevCommitAVSV2)"
    echo "  --reference <REFERENCE>     Name of the reference/old contract (e.g., MevCommitAVS)"
    echo "  --build-info <PATH>         Path to build-info directory (default: out/build-info)"
    echo "  --skip-build                Skip building contracts before validation"
    echo "  --help                      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --contract MevCommitAVSV2 --reference MevCommitAVS"
    echo "  $0 --contract ProviderRegistryV2 --reference ProviderRegistry --skip-build"
    exit 0
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --contract)
            contract="$2"
            shift 2
            ;;
        --reference)
            reference="$2"
            shift 2
            ;;
        --build-info)
            build_info_path="$2"
            shift 2
            ;;
        --skip-build)
            skip_build=true
            shift
            ;;
        --help|-h)
            show_help
            ;;
        *)
            echo "Error: Unknown option '$1'"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Validate required arguments
if [[ -z "$contract" ]]; then
    echo "Error: --contract is required"
    echo "Use --help for usage information"
    exit 1
fi

if [[ -z "$reference" ]]; then
    echo "Error: --reference is required"
    echo "Use --help for usage information"
    exit 1
fi

# Build contracts if not skipped
if [[ "$skip_build" != true ]]; then
    echo "Building contracts..."
    if ! (forge clean && forge build); then
        echo "Error: Failed to build contracts"
        exit 1
    fi
    echo ""
fi

# Check if build-info directory exists
if [[ ! -d "$build_info_path" ]]; then
    echo "Error: Build info directory not found: $build_info_path"
    echo "Make sure contracts are built (run 'forge clean && forge build')"
    exit 1
fi

# Run validation
echo "Validating upgrade safety..."
echo "Contract: $contract"
echo "Reference: $reference"
echo ""

if npx @openzeppelin/upgrades-core validate "$build_info_path" --contract "$contract" --reference "$reference"; then
    echo ""
    echo "✓ Upgrade validation passed!"
    exit 0
else
    echo ""
    echo "✗ Upgrade validation failed!"
    echo "Please review the errors above and fix the issues before proceeding with the upgrade."
    exit 1
fi

