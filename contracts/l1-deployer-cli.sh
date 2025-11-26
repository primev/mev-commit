#!/usr/bin/env bash

deploy_all_flag=false
deploy_vanilla_flag=false
deploy_vanilla_rep_flag=false
deploy_avs_flag=false
deploy_middleware_flag=false
deploy_rocketpool_flag=false
deploy_opt_in_hub_flag=false
deploy_block_rewards_flag=false
deploy_reward_distributor_flag=false
skip_release_verification_flag=false
resume_flag=false
wallet_type=""
private_key=""
chain=""
chain_id=0
deploy_contract=""
priority_gas_price=""
with_gas_price=""

help() {
    echo "Usage:"
    echo "  $0 <command> --chain <chain> <wallet option> [optional options]"
    echo
    echo "Commands (one required):"
    echo "  deploy-all          Deploy all components (vanilla, AVS, middleware, opt-in-hub)."
    echo "  deploy-vanilla      Deploy and verify the VanillaRegistry contract to L1."
    echo "  deploy-vanilla-rep  Deploy and verify the Reputational VanillaRegistry contract to L1."
    echo "  deploy-avs          Deploy and verify the MevCommitAVS contract to L1."
    echo "  deploy-middleware   Deploy and verify the MevCommitMiddleware contract to L1."
    echo "  deploy-rocketpool   Deploy and verify the RocketMinipoolRegistry contract to L1."
    echo "  deploy-opt-in-hub       Deploy and verify the ValidatorOptInHub contract to L1."
    echo "  deploy-block-rewards      Deploy and verify the BlockRewardManager contract to L1."
    echo "  deploy-reward-distributor      Deploy and verify the RewardDistributor contract to L1."
    echo
    echo "Required Options:"
    echo "  --chain, -c <chain>                Specify the chain to deploy to ('mainnet', 'holesky', 'hoodi', or 'anvil')."
    echo
    echo "Wallet Options (one required, except for anvil where --private-key is recommended):"
    echo "  --keystore                         Use a keystore for deployment."
    echo "  --ledger                           Use a Ledger hardware wallet for deployment."
    echo "  --trezor                           Use a Trezor hardware wallet for deployment."
    echo "  --private-key <KEY>               Use a private key for deployment (useful for anvil/local testing)."
    echo
    echo "Optional Options:"
    echo "  --skip-release-verification        Skip the GitHub release verification step."
    echo "  --resume                           Resume the deployment process if previously interrupted."
    echo "  --priority-gas-price <price>       Sets the priority gas price for EIP1559 transactions. Useful for when gas prices are volatile and you want to get your transaction included."
    echo "  --with-gas-price <price>           Sets the gas price for broadcasted legacy transactions, or the max fee for broadcasted EIP1559 transactions."
    echo "  --help                             Display this help message."
    echo
    echo "Notes:"
    echo "  - Exactly one command must be specified."
    echo "  - One wallet option must be specified (except for anvil where --private-key is recommended)."
    echo "  - Options and commands can be specified in any order."
    echo "  - Required arguments must immediately follow their options."
    echo
    echo "Environment Variables Required:"
    echo "  For Keystore:"
    echo "    KEYSTORES          Path(s) to keystore(s) passed to forge script as --keystores flag."
    echo "    KEYSTORE_PASSWORD  Password(s) for keystore(s) passed to forge script as --password flag."
    echo "    SENDER             Address of the sender."
    echo "    RPC_URL            RPC URL for the deployment chain."
    echo "    ETHERSCAN_API_KEY  API key for etherscan contract verification."
    echo
    echo "  For Ledger or Trezor:"
    echo "    HD_PATHS           Derivation path(s) for hardware wallets passed to forge script as --hd-paths flag."
    echo "    SENDER             Address of the sender."
    echo "    RPC_URL            RPC URL for the deployment chain."
    echo
    echo "  For Private Key (--private-key option):"
    echo "    SENDER             Address of the sender (optional, derived from private key if not set)."
    echo "    RPC_URL            RPC URL for the deployment chain."
    echo
    echo "Examples:"
    echo "  $0 deploy-all --chain mainnet --keystore --priority-gas-price 2000000000 --with-gas-price 5000000000"
    echo "  $0 --ledger deploy-avs --chain holesky --priority-gas-price 2000000000 --with-gas-price 5000000000"
    echo "  $0 --chain holesky deploy-middleware --trezor --priority-gas-price 2000000000 --with-gas-price 5000000000"
    echo "  $0 deploy-avs --chain anvil --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
    exit 1
}

usage() {
    echo "Usage:"
    echo "  $0 <command> --chain <chain> [wallet option] [options]"
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
            deploy-all)
                deploy_all_flag=true
                shift
                ;;
            deploy-vanilla)
                deploy_vanilla_flag=true
                shift
                ;;
            deploy-vanilla-rep)
                deploy_vanilla_rep_flag=true
                shift
                ;;
            deploy-avs)
                deploy_avs_flag=true
                shift
                ;;
            deploy-middleware)
                deploy_middleware_flag=true
                shift
                ;;
            deploy-rocketpool)
                deploy_rocketpool_flag=true
                shift
                ;;
            deploy-opt-in-hub)
                deploy_opt_in_hub_flag=true
                shift
                ;;
            deploy-block-rewards)
                deploy_block_rewards_flag=true
                shift
                ;;
            deploy-reward-distributor)
                deploy_reward_distributor_flag=true
                shift
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
            --skip-release-verification)
                skip_release_verification_flag=true
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

    commands_specified=0
    for flag in deploy_all_flag deploy_vanilla_flag deploy_vanilla_rep_flag deploy_avs_flag deploy_middleware_flag deploy_rocketpool_flag deploy_opt_in_hub_flag deploy_block_rewards_flag deploy_reward_distributor_flag; do
        if [[ "${!flag}" == true ]]; then
            ((commands_specified++))
        fi
    done

    if [[ $commands_specified -eq 0 ]]; then
        echo "Error: No command specified."
        usage
    elif [[ $commands_specified -gt 1 ]]; then
        echo "Error: Multiple commands specified. Please specify only one command at a time."
        usage
    fi
}

check_env_variables() {
    local missing_vars=()
    local required_vars=("RPC_URL")

    # ETHERSCAN_API_KEY is only required for non-anvil chains (for verification)
    if [[ "$chain" != "anvil" ]]; then
        required_vars+=("ETHERSCAN_API_KEY")
    fi

    if [[ "$wallet_type" == "keystore" ]]; then
        required_vars+=("KEYSTORES" "KEYSTORE_PASSWORD" "SENDER")
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
        deploy_contract="DeployMainnet"
    elif [[ "$chain" == "holesky" ]]; then
        chain_id=17000
        deploy_contract="DeployHolesky"
    elif [[ "$chain" == "hoodi" ]]; then
        chain_id=560048
        deploy_contract="DeployHoodi"
    elif [[ "$chain" == "anvil" ]]; then
        chain_id=31337
        # For anvil, use DeployAVSWithMockEigen which doesn't have chain-specific variants
        # It's a single contract that works for anvil
        deploy_contract="DeployAVSWithMockEigen"
    fi
}

check_git_status() {
    # Skip git checks for anvil (local testing)
    if [[ "$chain" == "anvil" ]]; then
        return
    fi
    
    if [[ ${chain_id:-0} -eq 1 ]]; then
        if ! current_tag=$(git describe --tags --exact-match 2>/dev/null); then
            echo "Error: Current commit is not tagged. Please ensure the commit is tagged before deploying."
            exit 1
        fi

        if [[ -n "$(git status --porcelain)" ]]; then
            echo "Error: There are uncommitted changes. Please commit or stash them before deploying."
            exit 1
        fi
    fi

    if [[ "$skip_release_verification_flag" != true ]]; then
        releases_url="https://api.github.com/repos/primev/mev-commit/releases?per_page=100"
        releases_json=$(curl -s "$releases_url")

        if [[ -z "$releases_json" ]]; then
            echo "Error: Unable to fetch releases from GitHub."
            exit 1
        fi

        release_tags=$(echo "$releases_json" | jq -r '.[].tag_name')

        if echo "$release_tags" | grep -q "^$current_tag$"; then
            echo "Tag '$current_tag' is associated with a release on GitHub."
        else
            echo "Error: Tag '$current_tag' is not associated with any release on GitHub. Please create a release for this tag before deploying."
            exit 1
        fi
    else
        echo "Skipping release verification as per user request."
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

    if [[ "$RPC_URL" != *"alchemy"* && "$RPC_URL" != *"infura"* ]]; then
        echo "Are you using a public rate-limited RPC URL? If so, contract verification may fail."
        read -p "Do you want to continue? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Exiting script."
            exit 1
        fi
    fi
}

check_etherscan_api_key() {
    # Skip etherscan check for anvil (not needed for local testing)
    if [[ "$chain" == "anvil" ]]; then
        return
    fi
    
    response=$(curl -s "https://api.etherscan.io/v2/api?chainid=${chain_id}&module=account&action=balance&address=${SENDER}&tag=latest&apikey=${ETHERSCAN_API_KEY}")

    status=$(echo "$response" | grep -o '"status":"[0-9]"' | cut -d':' -f2 | tr -d '"')

    if [[ "$status" != "1" ]]; then
        echo "Error: Etherscan API call failed or invalid etherscan API key."
        exit 1
    fi
}

deploy_contract_generic() {
    local script_path="$1"

    forge clean

    declare -a forge_args=()
    forge_args+=("script" "${script_path}:${deploy_contract}")
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

    local wallet_desc="${wallet_type:-private-key (anvil default)}"
    
    if forge "${forge_args[@]}"; then
        echo "Successfully ran ${script_path} on chain ID ${chain_id} using ${wallet_desc}."
        echo "Remember to update documentation with new contract addresses!"
    else
        echo "Error: Failed to run ${script_path} on chain ID ${chain_id} using ${wallet_desc}."
        exit 1
    fi
}

deploy_vanilla() {
    deploy_contract_generic "scripts/validator-registry/DeployVanillaRegistry.s.sol"
}

deploy_vanilla_rep() {
    deploy_contract_generic "scripts/validator-registry/DeployVanillaReputational.s.sol"
}

deploy_avs() {
    local script_path
    if [[ "$chain" == "anvil" ]]; then
        script_path="scripts/validator-registry/avs/DeployAVSWithMockEigen.s.sol"
    else
        script_path="scripts/validator-registry/avs/DeployAVS.s.sol"
    fi
    deploy_contract_generic "$script_path"
}

deploy_middleware() {
    deploy_contract_generic "scripts/validator-registry/middleware/DeployMiddleware.s.sol"
}

deploy_rocketpool() {
    deploy_contract_generic "scripts/validator-registry/rocketpool/DeployRocketMinipoolRegistry.s.sol"
}

deploy_opt_in_hub() {
    deploy_contract_generic "scripts/validator-registry/DeployValidatorOptInHub.s.sol"
}

deploy_block_rewards() {
    deploy_contract_generic "scripts/validator-registry/rewards/DeployBlockRewardManager.s.sol"
}

deploy_reward_distributor() {
    deploy_contract_generic "scripts/validator-registry/rewards/DeployRewardDistributor.s.sol"
}

main() {
    check_dependencies
    parse_args "$@"
    check_env_variables
    get_chain_params
    check_git_status
    check_rpc_url
    check_etherscan_api_key

    if [[ "${deploy_all_flag}" == true ]]; then
        echo "Deploying all contracts to $chain using $wallet_type..."
        deploy_vanilla
        deploy_vanilla_rep
        deploy_avs
        deploy_middleware
        deploy_rocketpool
        deploy_opt_in_hub
    elif [[ "${deploy_vanilla_flag}" == true ]]; then
        deploy_vanilla
    elif [[ "${deploy_vanilla_rep_flag}" == true ]]; then
        deploy_vanilla_rep
    elif [[ "${deploy_avs_flag}" == true ]]; then
        deploy_avs
    elif [[ "${deploy_middleware_flag}" == true ]]; then
        deploy_middleware
    elif [[ "${deploy_rocketpool_flag}" == true ]]; then
        deploy_rocketpool
    elif [[ "${deploy_opt_in_hub_flag}" == true ]]; then
        deploy_opt_in_hub
    elif [[ "${deploy_block_rewards_flag}" == true ]]; then
        deploy_block_rewards
    elif [[ "${deploy_reward_distributor_flag}" == true ]]; then
        deploy_reward_distributor
    else
        usage
    fi
}

main "$@"
