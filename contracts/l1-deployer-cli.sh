#!/usr/bin/env bash

deploy_all_flag=false
deploy_vanilla_flag=false
deploy_avs_flag=false
deploy_middleware_flag=false
deploy_router_flag=false
verify_bridge_flag=false
chain="holesky"
chain_id=17000

help() {
    echo "Usage:"
    echo "  $0 [--chain <chain>] <command>"
    echo
    echo "Commands:"
    echo "  deploy-all          Deploy all components (vanilla, AVS, middleware, router, verify bridge)."
    echo "  deploy-vanilla      Deploy and verify the VanillaRegistry contract to L1."
    echo "  deploy-avs          Deploy and verify the MevCommitAVS contract to L1."
    echo "  deploy-middleware   Deploy and verify the MevCommitMiddleware contract to L1."
    echo "  deploy-router       Deploy and verify the ValidatorOptInRouter contract to L1."
    echo "  verify-bridge       Verify the L1Gateway contract with etherscan."
    echo
    echo "Options:"
    echo "  --chain, -c <chain>    Specify the chain to deploy to ('mainnet' or 'holesky'). Default is 'holesky'."
    echo "  --help                 Display this help message."
    echo
    exit 1
}

usage() {
    echo "Usage:"
    echo "  $0 [--chain <chain>] <command>"
    echo
    echo "Use '$0 --help' to see available commands and options."
    exit 1
}

check_dependencies() {
    local missing_utils=()
    local required_utilities=(
        git
        forge
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
            deploy-avs)
                deploy_avs_flag=true
                shift
                ;;
            deploy-middleware)
                deploy_middleware_flag=true
                shift
                ;;
            deploy-router)
                deploy_router_flag=true
                shift
                ;;
            verify-bridge)
                verify_bridge_flag=true
                shift
                ;;
            --chain|-c)
                if [[ -z "$2" ]]; then
                    echo "Error: --chain requires an argument."
                    exit 1
                fi
                chain="$2"
                if [[ "$chain" != "mainnet" && "$chain" != "holesky" ]]; then
                    echo "Error: Unknown chain '$chain'. Valid options are 'mainnet' or 'holesky'."
                    exit 1
                fi
                if [[ "$chain" == "mainnet" ]]; then
                    chain_id=1
                elif [[ "$chain" == "holesky" ]]; then
                    chain_id=17000
                fi
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

    commands_specified=0
    for flag in deploy_all_flag deploy_vanilla_flag deploy_avs_flag deploy_middleware_flag deploy_router_flag verify_bridge_flag; do
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

check_git_status() {
    if ! git describe --tags --exact-match > /dev/null 2>&1; then
        echo "Error: Current commit is not tagged. Please ensure the commit is tagged before deploying."
        exit 1
    fi

    if [[ -n "$(git status --porcelain)" ]]; then
        echo "Error: There are uncommitted changes. Please commit or stash them before deploying."
        exit 1
    fi
}

get_deploy_contract() {
    if [[ "$chain" == "mainnet" ]]; then
        deploy_contract_suffix="Mainnet"
    elif [[ "$chain" == "holesky" ]]; then
        deploy_contract_suffix="Holesky"
    fi
    deploy_contract="Deploy${deploy_contract_suffix}"
}

main() {
    check_dependencies
    parse_args "$@"

    get_deploy_contract

    if [[ "${deploy_all_flag}" == true ]]; then
        check_git_status
        echo "Deploying all contracts to $chain..."
        deploy_vanilla
        deploy_avs
        deploy_middleware
        deploy_router
    elif [[ "${deploy_vanilla_flag}" == true ]]; then
        check_git_status
        deploy_vanilla
    elif [[ "${deploy_avs_flag}" == true ]]; then
        check_git_status
        deploy_avs
    elif [[ "${deploy_middleware_flag}" == true ]]; then
        check_git_status
        deploy_middleware
    elif [[ "${deploy_router_flag}" == true ]]; then
        check_git_status
        deploy_router
    elif [[ "${verify_bridge_flag}" == true ]]; then
        verify_bridge
    else
        usage
    fi
}

deploy_vanilla() {
    echo "Deploying VanillaRegistry contract to $chain..."
    echo "Using $deploy_contract contract for deployment."
    # Actual deployment commands go here, using $deploy_contract and $chain_id
    # Example:
    # forge script script/${deploy_contract}.sol:VanillaRegistry --chain-id $chain_id --other-options
    echo $chain_id
}

deploy_avs() {
    echo "Deploying MevCommitAVS contract to $chain..."
    echo "Using $deploy_contract contract for deployment."
    # Actual deployment commands go here, using $deploy_contract and $chain_id
}

deploy_middleware() {
    echo "Deploying MevCommitMiddleware contract to $chain..."
    echo "Using $deploy_contract contract for deployment."
    # Actual deployment commands go here, using $deploy_contract and $chain_id
}

deploy_router() {
    echo "Deploying ValidatorOptInRouter contract to $chain..."
    echo "Using $deploy_contract contract for deployment."
    # Actual deployment commands go here, using $deploy_contract and $chain_id
}

verify_bridge() {
    echo "Verifying L1Gateway contract with etherscan..."
    # Add actual verification commands here
}

main "$@"
