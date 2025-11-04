#!/usr/bin/env bash

# MUST READ - /mev-commit/contracts/scripts/upgrades/README.md

upgrade_contract_flag=false
old_contract=""
new_contract=""
proxy_address=""
wallet_type=""
private_key=""
chain=""
chain_id=0
upgrade_script=""
priority_gas_price=""
with_gas_price=""
resume_flag=false
multisig_flag=false

help() {
    echo "Usage:"
    echo "  $0 upgrade --old-contract <OLD_CONTRACT> --new-contract <NEW_CONTRACT> --proxy-address <PROXY_ADDRESS> --chain <chain> <wallet option> [optional options]"
    echo
    echo "Commands (one required):"
    echo "  upgrade              Upgrade a contract from old version to new version."
    echo
    echo "Required Options:"
    echo "  --old-contract, -o <OLD_CONTRACT>      Name of the old contract version (e.g., MevCommitAVS)."
    echo "  --new-contract, -n <NEW_CONTRACT>     Name of the new contract version (e.g., MevCommitAVSV2)."
    echo "  --proxy-address, -p <PROXY_ADDRESS>    Address of the proxy contract to upgrade (not required with --multisig)."
    echo "  --chain, -c <chain>                    Specify the chain to upgrade on ('mainnet', 'holesky', 'hoodi', or 'anvil')."
    echo
    echo "Wallet Options (one required, except for anvil where --private-key is recommended):"
    echo "  --keystore                             Use a keystore for upgrade."
    echo "  --ledger                               Use a Ledger hardware wallet for upgrade."
    echo "  --trezor                               Use a Trezor hardware wallet for upgrade."
    echo "  --private-key <KEY>                    Use a private key for upgrade (useful for anvil/local testing)."
    echo
    echo "Optional Options:"
    echo "  --multisig                             Deploy implementation contract only (for multisig upgrades)."
    echo "                                         When used, proxy upgrade is skipped and only the implementation is deployed."
    echo "                                         Proxy address is not required when using this flag."
    echo "  --resume                               Resume the upgrade process if previously interrupted."
    echo "  --priority-gas-price <price>           Sets the priority gas price for EIP1559 transactions. Useful for when gas prices are volatile."
    echo "  --with-gas-price <price>               Sets the gas price for broadcasted legacy transactions, or the max fee for broadcasted EIP1559 transactions."
    echo "  --help                                 Display this help message."
    echo
    echo "Notes:"
    echo "  - Exactly one command must be specified."
    echo "  - Options and commands can be specified in any order."
    echo "  - Required arguments must immediately follow their options."
    echo
    echo "Environment Variables Required:"
    echo "  For Keystore:"
    echo "    KEYSTORES          Path(s) to keystore(s) passed to forge script as --keystores flag."
    echo "    KEYSTORE_PASSWORD  Password(s) for keystore(s) passed to forge script as --password flag."
    echo "    SENDER             Address of the sender."
    echo "    RPC_URL            RPC URL for the upgrade chain."
    echo "    ETHERSCAN_API_KEY  API key for etherscan contract verification."
    echo
    echo "  For Ledger or Trezor:"
    echo "    HD_PATHS           Derivation path(s) for hardware wallets passed to forge script as --hd-paths flag."
    echo "    SENDER             Address of the sender."
    echo "    RPC_URL            RPC URL for the upgrade chain."
    echo
    echo "  For Private Key (--private-key option):"
    echo "    SENDER             Address of the sender (optional, derived from private key if not set)."
    echo "    RPC_URL            RPC URL for the upgrade chain."
    echo
    echo "  Optional:"
    echo "    PROXY_ADDRESS      Proxy address (can be provided via --proxy-address instead)."
    echo
    echo "Examples:"
    echo "  $0 upgrade --old-contract MevCommitAVS --new-contract MevCommitAVSV2 --proxy-address 0x1234... --chain mainnet --keystore"
    echo "  $0 upgrade --old-contract ProviderRegistry --new-contract ProviderRegistryV2 --proxy-address 0x5678... --chain holesky --ledger --priority-gas-price 2000000000"
    echo "  $0 upgrade --old-contract MevCommitAVS --new-contract MevCommitAVSV2 --proxy-address 0x1234... --chain anvil --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    echo "  $0 upgrade --old-contract MevCommitAVS --new-contract MevCommitAVSV2 --chain mainnet --multisig --keystore"
    exit 1
}

usage() {
    echo "Usage:"
    echo "  $0 upgrade --old-contract <OLD_CONTRACT> --new-contract <NEW_CONTRACT> --proxy-address <PROXY_ADDRESS> --chain <chain> [wallet option] [options]"
    echo
    echo "Use '$0 --help' to see available commands and options."
    exit 1
}

check_dependencies() {
    local missing_utils=()
    local required_utilities=(
        git
        forge
        cast
        curl
        jq
    )
    for util in "${required_utilities[@]}"; do
        if ! command -v "${util}" &> /dev/null; then
            missing_utils+=("${util}")
        fi
    done
    if [[ ${#missing_utils[@]} -ne 0 ]]; then
        echo "Error: The following required utilities are not installed: ${missing_utils[*]}."
        exit 1
    fi
}

parse_args() {
    if [[ $# -eq 0 ]]; then
        usage
    fi

    while [[ $# -gt 0 ]]; do
        key="$1"
        case $key in
            upgrade)
                upgrade_contract_flag=true
                shift
                ;;
            --old-contract|-o)
                if [[ -z "$2" ]]; then
                    echo "Error: --old-contract requires an argument."
                    exit 1
                fi
                old_contract="$2"
                shift 2
                ;;
            --new-contract|-n)
                if [[ -z "$2" ]]; then
                    echo "Error: --new-contract requires an argument."
                    exit 1
                fi
                new_contract="$2"
                shift 2
                ;;
            --proxy-address|-p)
                if [[ -z "$2" ]]; then
                    echo "Error: --proxy-address requires an argument."
                    exit 1
                fi
                proxy_address="$2"
                shift 2
                ;;
            --chain|-c)
                if [[ -z "$2" ]]; then
                    echo "Error: --chain requires an argument."
                    exit 1
                fi
                chain="$2"
                if [[ "$chain" != "mainnet" && "$chain" != "holesky" && "$chain" != "hoodi" && "$chain" != "anvil" ]]; then
                    echo "Error: Unknown chain '$chain'. Valid options are 'mainnet', 'holesky', 'hoodi', or 'anvil'."
                    exit 1
                fi
                shift 2
                ;;
            --multisig)
                multisig_flag=true
                shift
                ;;
            --resume)
                resume_flag=true
                shift
                ;;
            --keystore)
                if [[ -n "$wallet_type" ]]; then
                    echo "Error: Multiple wallet types specified. Please specify only one wallet option."
                    exit 1
                fi
                wallet_type="keystore"
                shift
                ;;
            --ledger)
                if [[ -n "$wallet_type" ]]; then
                    echo "Error: Multiple wallet types specified. Please specify only one wallet option."
                    exit 1
                fi
                wallet_type="ledger"
                shift
                ;;
            --trezor)
                if [[ -n "$wallet_type" ]]; then
                    echo "Error: Multiple wallet types specified. Please specify only one wallet option."
                    exit 1
                fi
                wallet_type="trezor"
                shift
                ;;
            --private-key)
                if [[ -z "$2" ]]; then
                    echo "Error: --private-key requires an argument."
                    exit 1
                fi
                if [[ -n "$wallet_type" ]]; then
                    echo "Error: Multiple wallet types specified. Please specify only one wallet option."
                    exit 1
                fi
                wallet_type="private-key"
                private_key="$2"
                shift 2
                ;;
            --priority-gas-price)
                if [[ -z "$2" ]]; then
                    echo "Error: --priority-gas-price requires an argument."
                    exit 1
                fi
                priority_gas_price="$2"
                shift 2
                ;;
            --with-gas-price)
                if [[ -z "$2" ]]; then
                    echo "Error: --with-gas-price requires an argument."
                    exit 1
                fi
                with_gas_price="$2"
                shift 2
                ;;
            --help)
                help
                ;;
            *)
                echo "Error: Unknown command or option '$1'."
                usage
                ;;
        esac
    done

    if [[ -z "$chain" ]]; then
        echo "Error: The --chain option is required."
        usage
    fi

    if [[ -z "$wallet_type" && "$chain" != "anvil" ]]; then
        echo "Error: A wallet option is required. Please specify one of --keystore, --ledger, --trezor, or --private-key."
        echo "Note: For anvil, --private-key is recommended but not required."
        usage
    fi

    if [[ "$upgrade_contract_flag" != true ]]; then
        echo "Error: No command specified. Use 'upgrade' command."
        usage
    fi

    if [[ -z "$old_contract" ]]; then
        echo "Error: The --old-contract option is required."
        usage
    fi

    if [[ -z "$new_contract" ]]; then
        echo "Error: The --new-contract option is required."
        usage
    fi

    # Proxy address is only required for direct upgrades, not for multisig deployments
    if [[ "$multisig_flag" != true ]]; then
        if [[ -z "$proxy_address" ]]; then
            # Try to get from environment variable
            if [[ -n "${PROXY_ADDRESS}" ]]; then
                proxy_address="${PROXY_ADDRESS}"
            else
                echo "Error: The --proxy-address option is required (or set PROXY_ADDRESS environment variable)."
                echo "Note: If deploying for multisig upgrade, use --multisig flag (proxy address not needed)."
                usage
            fi
        fi
    fi
}

check_env_variables() {
    local missing_vars=()
    local required_vars=("RPC_URL")

    if [[ "$wallet_type" == "keystore" ]]; then
        required_vars+=("KEYSTORES" "KEYSTORE_PASSWORD" "SENDER")
        # ETHERSCAN_API_KEY is optional for upgrades but recommended for verification
    elif [[ "$wallet_type" == "ledger" || "$wallet_type" == "trezor" ]]; then
        required_vars+=("HD_PATHS" "SENDER")
    elif [[ "$wallet_type" == "private-key" ]]; then
        # SENDER is optional for private-key, can be derived from the key
        # But we'll still require it for consistency unless on anvil
        if [[ "$chain" != "anvil" ]]; then
            required_vars+=("SENDER")
        fi
    elif [[ -z "$wallet_type" && "$chain" == "anvil" ]]; then
        # For anvil without explicit wallet, we'll use private-key from env or default
        # SENDER is optional
        :
    fi

    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_vars+=("${var}")
        fi
    done

    if [[ ${#missing_vars[@]} -ne 0 ]]; then
        echo "Error: The following environment variables are not set: ${missing_vars[*]}."
        echo "Please set them before running the script."
        exit 1
    fi
}

get_chain_params() {
    if [[ "$chain" == "mainnet" ]]; then
        chain_id=1
        if [[ "$multisig_flag" == true ]]; then
            upgrade_script="DeployMultisigImplMainnet"
        else
            upgrade_script="UpgradeContractMainnet"
        fi
    elif [[ "$chain" == "holesky" ]]; then
        chain_id=17000
        if [[ "$multisig_flag" == true ]]; then
            upgrade_script="DeployMultisigImplHolesky"
        else
            upgrade_script="UpgradeContractHolesky"
        fi
    elif [[ "$chain" == "hoodi" ]]; then
        chain_id=560048
        if [[ "$multisig_flag" == true ]]; then
            upgrade_script="DeployMultisigImplHoodi"
        else
            upgrade_script="UpgradeContractHoodi"
        fi
    elif [[ "$chain" == "anvil" ]]; then
        chain_id=31337
        if [[ "$multisig_flag" == true ]]; then
            upgrade_script="DeployMultisigImplAnvil"
        else
            upgrade_script="UpgradeContractAnvil"
        fi
    fi
}

check_git_status() {
    # Skip git checks for anvil (local testing)
    if [[ "$chain" == "anvil" ]]; then
        return
    fi
    
    if [[ ${chain_id:-0} -eq 1 ]]; then
        if ! current_tag=$(git describe --tags --exact-match 2>/dev/null); then
            echo "Error: Current commit is not tagged. Please ensure the commit is tagged before upgrading on mainnet."
            exit 1
        fi

        if [[ -n "$(git status --porcelain)" ]]; then
            echo "Error: There are uncommitted changes. Please commit or stash them before upgrading on mainnet."
            exit 1
        fi
    fi
}

check_rpc_url() {
    queried_chain_id=$(cast chain-id --rpc-url "$RPC_URL" 2>/dev/null)
    cast_exit_code=$?
    if [[ $cast_exit_code -ne 0 ]]; then
        echo "Error: Failed to retrieve chain ID using the provided RPC URL."
        echo "Please ensure the RPC URL is correct and accessible."
        exit 1
    fi
    if [[ "$queried_chain_id" -ne "$chain_id" ]]; then
        echo "Error: RPC URL does not match the expected chain ID."
        echo "Expected chain ID: $chain_id, but got: $queried_chain_id."
        exit 1
    fi

    # Skip RPC URL warning for anvil (local testing)
    if [[ "$chain" != "anvil" && "$RPC_URL" != *"alchemy"* && "$RPC_URL" != *"infura"* ]]; then
        echo "Warning: You may be using a public rate-limited RPC URL. Contract verification may fail."
        read -p "Do you want to continue? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Exiting script."
            exit 1
        fi
    fi
}

find_contract_path() {
    local contract_name="$1"
    local search_path="$2"
    
    # Try to find the contract file
    local contract_file=$(find contracts -name "${contract_name}.sol" -type f 2>/dev/null | head -n 1)
    
    if [[ -n "$contract_file" ]]; then
        echo "$contract_file"
        return 0
    fi
    
    # If not found in contracts, try upgrades folder
    contract_file=$(find contracts/upgrades -name "${contract_name}.sol" -type f 2>/dev/null | head -n 1)
    
    if [[ -n "$contract_file" ]]; then
        echo "$contract_file"
        return 0
    fi
    
    return 1
}

upgrade_contract() {
    # Find the contract paths
    local new_contract_path=$(find_contract_path "$new_contract" "contracts")
    local old_contract_path=$(find_contract_path "$old_contract" "contracts/upgrades")
    
    if [[ -z "$new_contract_path" ]]; then
        echo "Error: Could not find contract file for $new_contract"
        echo "Please ensure the contract exists in the contracts/ directory."
        exit 1
    fi
    
    # Extract just the filename for deployment
    local new_contract_filename=$(basename "$new_contract_path")
    
    if [[ "$multisig_flag" == true ]]; then
        echo "Deploying implementation contract for multisig upgrade..."
    else
        echo "Upgrading contract..."
    fi
    echo ""
    
    forge clean

    declare -a forge_args=()
    
    # Select the appropriate script based on multisig flag
    if [[ "$multisig_flag" == true ]]; then
        forge_args+=("script" "scripts/upgrades/GenericMultisigUpgrade.s.sol:${upgrade_script}")
    else
        forge_args+=("script" "scripts/upgrades/GenericUpgrade.s.sol:${upgrade_script}")
    fi
    
    forge_args+=("--rpc-url" "${RPC_URL}")
    forge_args+=("--via-ir")
    forge_args+=("--chain-id" "${chain_id}")
    forge_args+=("--use" "0.8.26")
    forge_args+=("--broadcast")
    
    # Add verification if ETHERSCAN_API_KEY is set (and not anvil)
    if [[ -n "${ETHERSCAN_API_KEY:-}" && "$chain" != "anvil" ]]; then
        forge_args+=("--verify")
    fi
    
    if [[ -n "$priority_gas_price" ]]; then
        forge_args+=("--priority-gas-price" "${priority_gas_price}")
    fi

    if [[ -n "$with_gas_price" ]]; then
        forge_args+=("--with-gas-price" "${with_gas_price}")
    fi

    if [[ "$resume_flag" == true ]]; then
        forge_args+=("--resume")
    fi

    if [[ "$wallet_type" == "keystore" ]]; then
        forge_args+=("--keystores" "${KEYSTORES}")
        forge_args+=("--password" "${KEYSTORE_PASSWORD}")
        forge_args+=("--sender" "${SENDER}")
    elif [[ "$wallet_type" == "ledger" ]]; then
        forge_args+=("--ledger")
        forge_args+=("--hd-paths" "${HD_PATHS}")
        forge_args+=("--sender" "${SENDER}")
    elif [[ "$wallet_type" == "trezor" ]]; then
        forge_args+=("--trezor")
        forge_args+=("--hd-paths" "${HD_PATHS}")
        forge_args+=("--sender" "${SENDER}")
    elif [[ "$wallet_type" == "private-key" ]]; then
        forge_args+=("--private-key" "${private_key}")
        if [[ -n "${SENDER:-}" ]]; then
            forge_args+=("--sender" "${SENDER}")
        fi
    elif [[ -z "$wallet_type" && "$chain" == "anvil" ]]; then
        # For anvil without explicit wallet, try to use private key from env or default anvil key
        if [[ -n "${PRIVATE_KEY:-}" ]]; then
            forge_args+=("--private-key" "${PRIVATE_KEY}")
        elif [[ -n "${private_key:-}" ]]; then
            forge_args+=("--private-key" "${private_key}")
        else
            # Use default anvil private key (first account)
            forge_args+=("--private-key" "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
        fi
        if [[ -n "${SENDER:-}" ]]; then
            forge_args+=("--sender" "${SENDER}")
        fi
    fi
    
    # Pass contract names and paths as environment variables
    export OLD_CONTRACT_NAME="$old_contract"
    export NEW_CONTRACT_NAME="$new_contract"
    export NEW_CONTRACT_PATH="$new_contract_filename"
    
    # Proxy address only needed for direct upgrades
    if [[ "$multisig_flag" != true ]]; then
        export PROXY_ADDRESS="$proxy_address"
    fi
    
    local wallet_desc="${wallet_type:-private-key (anvil default)}"
    
    if forge "${forge_args[@]}"; then
        if [[ "$multisig_flag" == true ]]; then
            echo ""
            echo "✓ Implementation contract deployed successfully!"
            echo "  Contract: $old_contract -> $new_contract"
            echo "  Chain ID: ${chain_id}"
        else
            echo ""
            echo "✓ Successfully upgraded $old_contract to $new_contract"
            echo "  Chain ID: ${chain_id}"
            echo "  Proxy: $proxy_address"
        fi
    else
        if [[ "$multisig_flag" == true ]]; then
            echo "Error: Failed to deploy implementation contract for $new_contract on chain ID ${chain_id} using ${wallet_desc}."
        else
            echo "Error: Failed to upgrade $old_contract to $new_contract on chain ID ${chain_id} using ${wallet_desc}."
        fi
        exit 1
    fi
}

main() {
    check_dependencies
    parse_args "$@"
    check_env_variables
    get_chain_params
    check_git_status
    check_rpc_url
    upgrade_contract
}

main "$@"

