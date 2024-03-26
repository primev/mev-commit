#!/bin/bash
set -ex

config_file=".bridge_config"

show_usage() {
    echo "Usage: $0 [command] [arguments] [options]"
    echo ""
    echo "Commands:"
    echo "  bridge-to-mev-commit <amount in wei> <dest_addr> <private_key>"
    echo "    Bridge tokens to mev-commit chain. Requires the amount to bridge, destination account, and private key."
    echo "    Example: $0 bridge-to-mev-commit 100 0x123... 0xABC..."
    echo ""
    echo "  bridge-to-l1 <amount in wei> <dest_addr> <private_key>"
    echo "    Bridge tokens to L1. Requires the amount to bridge, destination account, and private key."
    echo "    Example: $0 bridge-to-l1 100 0x456... 0xDEF..."
    echo ""
    echo "  init <L1 Router> <mev-commit chain Router> <L1 Chain ID> <mev-commit chain ID> <L1 URL> <MEV-Commit URL>"
    echo "    Initialize configuration with specified hyperlane router addresses, chain IDs, and URLs."
    echo "    Example: $0 init 0xc20B3C7852FA81f36130313220890eA7Ea5F5B0e 0x4b2DC8A5C4da51f821390AbD2B6fe8122BC6fA97 11155111 17864 https://ethereum-sepolia.publicnode.com http://localhost:8545"
    echo ""
    echo "Options:"
    echo "  -y, --yes   Automatically answer 'yes' to all prompts"
    echo "    Example: $0 bridge-to-mev-commit 100 0x123... 0xABC... -y"
    echo ""
    echo "Note: This script requires foundry and jq to be installed."
}


bridge_confirmation() {
    if [ "$skip_confirmation" = false ]; then
        local source_chain_name=$1
        local dest_chain_name=$2
        local source_chain_id=$3
        local dest_chain_id=$4
        local source_url=$5
        local dest_url=$6
        local source_router=$7
        local dest_router=$8
        local amount=$9
        local dest_address=${10} # Arguments after $9 must be accessed with braces

        echo "You are about to bridge..."
        echo "From $source_chain_name (ID: $source_chain_id, URL: $source_url, Router: $source_router)"
        echo "To $dest_chain_name (ID: $dest_chain_id, URL: $dest_url, Router: $dest_router)"
        echo "Amount to bridge: $amount wei"
        echo "Destination address: $dest_address"
        read -p "Are you sure you want to proceed with the bridging operation? (y/n): " answer
        if [ "$answer" != "y" ]; then
            echo "Operation cancelled."
            exit 1
        fi
    fi
}

check_config_loaded() {
    local config_vars=("l1_router" "mev_commit_chain_router" "l1_chain_id" "mev_commit_chain_id" "l1_url" "mev_commit_url")

    for var in "${config_vars[@]}"; do
        if [ -z "${!var}" ]; then
            echo "Error: Configuration for '$var' not loaded."
            echo "Please run the init command to set up the necessary configuration."
            exit 1
        fi
    done
}

# TODO: Support more secure forms of private key management: https://book.getfoundry.sh/tutorials/best-practices#private-key-management
bridge_transfer() {
    local source_chain_name=$1
    local dest_chain_name=$2
    local source_chain_id=$3
    local dest_chain_id=$4
    local source_url=$5
    local dest_url=$6
    local source_router=$7
    local dest_router=$8
    local amount=$9
    local dest_address=${10}
    local private_key=${11}

    check_config_loaded

    if ! [[ $amount =~ ^[0-9]+$ ]]; then
        echo "Error: Amount of wei is not a valid number."
        return 1
    fi

    bridge_confirmation \
        "$source_chain_name" \
        "$dest_chain_name" \
        "$source_chain_id" \
        "$dest_chain_id" \
        "$source_url" \
        "$dest_url" \
        "$source_router" \
        "$dest_router" \
        "$amount" \
        "$dest_address"

    echo "Bridging to $dest_chain_name..."
    echo "Amount: $amount"
    echo "Destination Address: $dest_address"
    echo "Using $source_chain_name Router: $source_router"
    echo "Using $dest_chain_name Router: $dest_router"
    echo "$source_chain_name Chain ID: $source_chain_id"
    echo "$dest_chain_name Chain ID: $dest_chain_id"
    echo "$source_chain_name URL: $source_url"
    echo "$dest_chain_name URL: $dest_url"

    local dest_address_clean=${dest_address#0x} # Remove prefix
    local padded_dest_address="0x000000000000000000000000$dest_address_clean"

    local gas_payment_hex=$(cast call --rpc-url $source_url $source_router "quoteGasPayment(uint32)" $dest_chain_id)
    local gas_payment_hex_clean=${gas_payment_hex#0x} # Remove prefix
    local gas_payment_dec=$((16#$gas_payment_hex_clean))
    local total_amount_dec=$(($amount + $gas_payment_dec))

    dest_account_init_balance=$(cast balance --rpc-url $dest_url $dest_address)

    src_address=$(cast wallet address $private_key)
    src_account_init_balance=$(cast balance --rpc-url $source_url $src_address)

    # TODO: Tune transaction parameters and/or allow for user to inject custom config with sourced env vars.
    # See https://book.getfoundry.sh/reference/cli/cast/send
    cast send \
        --rpc-url $source_url \
        --private-key $private_key \
        $source_router "transferRemote(uint32,bytes32,uint256)" \
        $dest_chain_id \
        $padded_dest_address \
        $amount \
        --value $total_amount_dec"wei"
        # --nonce 366

    # After tx is submitted to src chain, a simple FSM ensues. Bridge invocation follows:
    #
    # [PENDING] -- Timeout --> [EXPIRED]
    #    |
    #    | tx committed to source chain
    #    |
    # [COMMITTED] -- Timeout --> [FAILURE]
    #    |
    #    | balance incremented (tx committed to dest chain)
    #    v
    # [SUCCESS]
    #
    # Pending: tx has been submitted to source chain, but is still pending (src balance not decremented)
    # Expired: 30 minute timeout has been reached while source chain tx is still pending
    # Committed: tx has been committed to source chain, 30 minute timeout is NOT reached
    # Failure: tx has been committed to source chain, 30 minute timeout has elapsed
    # Success: destination balance incremented, bridge invocation complete

    # Iterate until source balance changes. Timeout after 30 minutes. 
    max_retries=180
    sleep_time=10
    retry_count=0
    # TODO: Use transaction receipt instead of balance polling
    while [ "$(cast balance --rpc-url "$source_url" "$src_address")" = "$src_account_init_balance" ]
    do 
        echo "$((retry_count * 10)) seconds passed since bridge invocation. Waiting for source balance to change..."
        sleep $sleep_time
        retry_count=$((retry_count + 1))
        if [ "$retry_count" -ge "$max_retries" ]; then
            echo "Maximum retries reached. 30 minutes have passed and source balance has not changed."
            echo "EXPIRED"
            # TODO: If expired, cancel tx here
            return 0
        fi
    done

    # Iterate until destination balance changes. Timeout after 30 minutes.
    max_retries=180
    sleep_time=10
    retry_count=0
    while [ "$(cast balance --rpc-url "$dest_url" "$dest_address")" = "$dest_account_init_balance" ]
    do
        echo "$((retry_count * 10)) seconds passed since bridge invocation. Waiting for destination balance to change..."
        sleep $sleep_time
        retry_count=$((retry_count + 1))
        if [ "$retry_count" -ge "$max_retries" ]; then
            echo "Maximum retries reached. 30 minutes have passed and destination balance has not changed."
            echo "FAILURE"
            return 0
        fi
    done

    echo "Source and destination balances have changed. Bridge operation successful."
    echo "If in production, confirm destination balance was not incremented by irrelevant transaction."
    echo "SUCCESS"
    return 0
}

bridge_to_mev_commit() {
    bridge_transfer \
        "L1" \
        "mev-commit chain" \
        "$l1_chain_id" \
        "$mev_commit_chain_id" \
        "$l1_url" \
        "$mev_commit_url" \
        "$l1_router" \
        "$mev_commit_chain_router" \
        "$1" \
        "$2" \
        "$3"
}

bridge_to_l1() {
    bridge_transfer \
        "mev-commit chain" \
        "L1" \
        "$mev_commit_chain_id" \
        "$l1_chain_id" \
        "$mev_commit_url" \
        "$l1_url" \
        "$mev_commit_chain_router" \
        "$l1_router" \
        "$1" \
        "$2" \
        "$3"
}

# Function to initialize and save configuration
init_config() {
    local l1_router=$1
    local mev_commit_chain_router=$2
    local l1_chain_id=$3
    local mev_commit_chain_id=$4
    local l1_url=$5
    local mev_commit_url=$6

    # Confirmation prompt
    if [ "$skip_confirmation" = false ]; then
        echo "You are about to initialize the configuration with the following settings:"
        echo "L1 Router: $l1_router"
        echo "mev-commit chain Router: $mev_commit_chain_router"
        echo "L1 Chain ID: $l1_chain_id"
        echo "mev-commit chain ID: $mev_commit_chain_id"
        echo "L1 URL: $l1_url"
        echo "MEV-Commit URL: $mev_commit_url"
        read -p "Are you sure you want to proceed? (y/n): " answer
        if [ "$answer" != "y" ]; then
            echo "Operation cancelled."
            exit 1
        fi
    fi

    # Create JSON config file
    jq -n \
        --arg l1_router "$l1_router" \
        --arg mev_commit_chain_router "$mev_commit_chain_router" \
        --arg l1_chain_id "$l1_chain_id" \
        --arg mev_commit_chain_id "$mev_commit_chain_id" \
        --arg l1_url "$l1_url" \
        --arg mev_commit_url "$mev_commit_url" \
        '{l1_router: $l1_router, mev_commit_chain_router: $mev_commit_chain_router, l1_chain_id: $l1_chain_id, mev_commit_chain_id: $mev_commit_chain_id, l1_url: $l1_url, mev_commit_url: $mev_commit_url}' \
        > "$config_file"

    echo "Configuration initialized and saved."
}

# Loads configuration from JSON
load_config() {
    if [ -f "$config_file" ]; then
        l1_router=$(jq -r '.l1_router' "$config_file")
        mev_commit_chain_router=$(jq -r '.mev_commit_chain_router' "$config_file")
        l1_chain_id=$(jq -r '.l1_chain_id' "$config_file")
        mev_commit_chain_id=$(jq -r '.mev_commit_chain_id' "$config_file")
        l1_url=$(jq -r '.l1_url' "$config_file")
        mev_commit_url=$(jq -r '.mev_commit_url' "$config_file")
    else
        echo "Error: Configuration file not found. Please run the init command first."
        exit 1
    fi
}

# Check for help command before loading config
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_usage
    exit 0
fi

# If first arg is not "init", load configuration.
if [[ "$1" != "init" ]]; then
    load_config
fi

# Check if last argument is --yes or -y, set flag accordingly
skip_confirmation=false
if [[ "${@: -1}" == "--yes" || "${@: -1}" == "-y" ]]; then
    skip_confirmation=true
    set -- "${@:1:$#-1}"  # Remove the last argument
fi

command=$1
shift  # Shift to get the next set of parameters after the command

case "$command" in
    init)
        if [ $# -ne 6 ]; then
            echo "Error: Incorrect number of arguments for init command."
            show_usage
            exit 1
        fi
        init_config "$1" "$2" "$3" "$4" "$5" "$6"
        ;;
    bridge-to-mev-commit)
        if [ $# -ne 3 ]; then
            echo "Error: Incorrect number of arguments for bridging to mev-commit chain."
            show_usage
            exit 1
        fi
        bridge_to_mev_commit "$1" "$2" "$3"
        ;;
    bridge-to-l1)
        if [ $# -ne 3 ]; then
            echo "Error: Incorrect number of arguments for bridging to L1."
            show_usage
            exit 1
        fi
        bridge_to_l1 "$1" "$2" "$3"
        ;;
    -h|--help)
        show_usage
        ;;
    *)
        echo "Unknown command: $command"
        show_usage
        exit 1
        ;;
esac
