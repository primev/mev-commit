#!/usr/bin/env bash

deploy_all_flag=false
deploy_vanilla_flag=false
deploy_avs_flag=false
deploy_middleware_flag=false
deploy_router_flag=false
verify_bridge_flag=false

help() {
    echo "Usage:"
    echo "  $0 deploy-all"
    echo "  $0 deploy-vanilla"
    echo "  $0 deploy-avs"
    echo "  $0 deploy-middleware"
    echo "  $0 deploy-router"
    echo "  $0 verify-bridge"
    echo "  $0 --help"
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
    echo "  --help              Display this help message."
    echo
    exit 1
}

usage() {
    echo "Usage:"
    echo "  $0 <command>"
    echo
    echo "Use '$0 --help' to see available commands."
    exit 1
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

main() {
    parse_args "$@"

    if [[ "${deploy_all_flag}" == true ]]; then
        echo "Deploying all..."
        deploy_vanilla
        deploy_avs
        deploy_middleware
        deploy_router
    elif [[ "${deploy_vanilla_flag}" == true ]]; then
        deploy_vanilla
    elif [[ "${deploy_avs_flag}" == true ]]; then
        deploy_avs
    elif [[ "${deploy_middleware_flag}" == true ]]; then
        deploy_middleware
    elif [[ "${deploy_router_flag}" == true ]]; then
        deploy_router
    elif [[ "${verify_bridge_flag}" == true ]]; then
        verify_bridge
    else
        usage
    fi
}

deploy_vanilla() {
    echo "Deploying VanillaRegistry contract..."
}

deploy_avs() {
    echo "Deploying MevCommitAVS contract..."
}

deploy_middleware() {
    echo "Deploying MevCommitMiddleware contract..."
}

deploy_router() {
    echo "Deploying ValidatorOptInRouter contract..."
}

verify_bridge() {
    echo "Verifying L1Gateway contract with etherscan..."
}

main "$@"
