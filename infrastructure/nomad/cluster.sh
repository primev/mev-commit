#!/usr/bin/env bash

export OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES

init_flag=false
deploy_flag=false
destroy_flag=false
debug_flag=false
no_logs_collection_flag=false
force_build_templates_flag=false
skip_certificates_setup_flag=false
release_flag=false
deploy_version="HEAD"
environment_name="devenv"
profile_name="devnet"
datadog_key=""
l1_rpc_urls="mock"
etherscan_api_key=""
otel_collector_endpoint_url=""
genesis_file_url=""
geth_bootnode_url=""
oracle_relay_urls=""
settlement_rpc_url=""
contracts_json_url=""

help() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--l1-rpc-urls <urls>] [--oracle-relay-urls <urls>] [--etherscan-api-key <key>] [--otel-collector-endpoint-url <url>] [--genesis-file-url <url>] [--geth-bootnode-url <url>] [--settlement-rpc-url <url>] [--contracts-json-url <url>] [--release] [--debug]]"
    echo "$0 [destroy [--environment <name=devenv>] [--debug]]"
    echo "$0 --help"
    echo
    echo "Parameters:"
    echo "  init                           Initialize the environment."
    echo "    --environment <name=devenv>  Specify the environment to use (default is devenv)."
    echo "    --skip-certificates-setup    Skip the certificates installation and setup."
    echo "    --debug                      Enable debug mode for detailed output."
    echo
    echo "  deploy [version=HEAD]"
    echo "    --environment <name=devenv>          Specify the environment (default is devenv)."
    echo "    --profile <name=devnet>              Specify the profile (default is devnet)."
    echo "    --force-build-templates              Force the build of all job templates before deployment."
    echo "    --no-logs-collection                 Disable the collection of logs from deployed jobs."
    echo "    --datadog-key <key>                  Datadog API key, cannot be empty."
    echo "    --l1-rpc-urls <urls>                 Comma-separated list of L1 RPC URLs, cannot be empty."
    echo "    --oracle-relay-urls <urls>           Comma-separated list of Oracle Relay URLs, cannot be empty."
    echo "    --etherscan-api-key <key>            Etherscan API key, cannot be empty."
    echo "    --otel-collector-endpoint-url <url>  OpenTelemetry Collector Endpoint URL, cannot be empty."
    echo "    --genesis-file-url <url>             URL to the genesis file, cannot be empty."
    echo "    --geth-bootnode-url <url>            URL to the Geth bootnode, cannot be empty."
    echo "    --settlement-rpc-url <url>           URL for the settlement RPC, cannot be empty."
    echo "    --contracts-json-url <url>           URL to the contracts JSON file, cannot be empty."
    echo "    --release                            Ignore the specified deployment version and use the current HEAD tag as the build version."
    echo "    --debug                              Enable debug mode for detailed output."
    echo
    echo "  destroy [--environment <name=devenv>] [--debug]"
    echo "    Destroy the whole cluster."
    echo "    --environment <name=devenv>  Specify the environment to use (default is devenv)."
    echo "    --debug                      Enable debug mode for detailed output."
    echo
    echo "  --help  Display this help message."
    echo
    echo "Examples:"
    echo "  Destroy all jobs with a specific environment in debug mode:"
    echo "    $0 destroy --environment devenv --debug"
    exit 1
}

usage() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--l1-rpc-urls <urls>] [--oracle-relay-urls <urls>] [--etherscan-api-key <key>] [--otel-collector-endpoint-url <url>] [--genesis-file-url <url>] [--geth-bootnode-url <url>] [--settlement-rpc-url <url>] [--contracts-json-url <url>] [--release] [--debug]]"
    echo "$0 [destroy [--environment <name=devenv>] [--debug]]"
    echo "$0 --help"
    exit 1
}

check_deps() {
    local missing_utils=()
    local required_utilities=(
        go
        yq
        aws
        flock
        ansible
        bootnode
        remarshal
        goreleaser
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

    local missing_ansible_collections=()
    local required_ansible_collections=(
        community.aws
        community.general
    )
    for collection in "${required_ansible_collections[@]}"; do
        if ! ansible-galaxy collection list | grep -q "${collection}"; then
            missing_ansible_collections+=("${collection}")
        fi
    done
    if [[ ${#missing_ansible_collections[@]} -ne 0 ]]; then
        echo "Error: The following required Ansible collections are not installed: ${missing_ansible_collections[*]}."
        exit 1
    fi

    if ! aws sts get-caller-identity &> /dev/null; then
        echo "Error: AWS is not configured properly. Please run 'aws configure' to set up your credentials."
        exit 1
    fi

    if [[ ! -f ansible.cfg ]]; then
        echo "Error: ansible.cfg file not found."
        exit 1
    fi

    if [[ ! -f hosts.ini ]]; then
        echo "Error: hosts.ini file not found."
        exit 1
    fi
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        if $init_flag || $deploy_flag || $destroy_flag; then
            echo "Error: Only one of 'init', 'deploy', or 'destroy' can be specified."
            usage
        fi

        key="$1"
        case $key in
            init)
                init_flag=true
                shift
                if [[ $# -gt 0 && $1 == "--environment" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        environment_name="$2"
                        shift 2
                    else
                        echo "Error: --environment requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--skip-certificates-setup" ]]; then
                    skip_certificates_setup_flag=true
                    shift
                fi
                if [[ $# -gt 0 && $1 == "--debug" ]]; then
                    debug_flag=true
                    shift
                fi
                ;;
            deploy)
                deploy_flag=true
                if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                    deploy_version="$2"
                    shift
                fi
                shift
                if [[ $# -gt 0 && $1 == "--environment" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        environment_name="$2"
                        shift 2
                    else
                        echo "Error: --environment requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--profile" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        profile_name="$2"
                        shift 2
                    else
                        echo "Error: --profile requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--force-build-templates" ]]; then
                    force_build_templates_flag=true
                    shift
                fi
                if [[ $# -gt 0 && $1 == "--no-logs-collection" ]]; then
                    no_logs_collection_flag=true
                    shift
                fi
                if [[ $# -gt 0 && $1 == "--datadog-key" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        datadog_key="$2"
                        shift 2
                    else
                        echo "Error: --datadog-key requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--l1-rpc-urls" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        l1_rpc_urls="$2"
                        shift 2
                    else
                        echo "Error: --l1-rpc-urls requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--oracle-relay-urls" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        oracle_relay_urls="$2"
                        shift 2
                    else
                        echo "Error: --oracle-relay-urls requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--etherscan-api-key" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        etherscan_api_key="$2"
                        shift 2
                    else
                        echo "Error: --etherscan-api-key requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--otel-collector-endpoint-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        otel_collector_endpoint_url="$2"
                        shift 2
                    else
                        echo "Error: --otel-collector-endpoint-url requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--genesis-file-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        genesis_file_url="$2"
                        shift 2
                    else
                        echo "Error: --genesis-file-url requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--geth-bootnode-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        geth_bootnode_url="$2"
                        shift 2
                    else
                        echo "Error: --geth-bootnode-url requires a value."
                        usage
                    fi
                fi
                # Added flag: settlement_rpc_url
                if [[ $# -gt 0 && $1 == "--settlement-rpc-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        settlement_rpc_url="$2"
                        shift 2
                    else
                        echo "Error: --settlement-rpc-url requires a value."
                        usage
                    fi
                fi
                # Added flag: contracts_json_url
                if [[ $# -gt 0 && $1 == "--contracts-json-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        contracts_json_url="$2"
                        shift 2
                    else
                        echo "Error: --contracts-json-url requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--release" ]]; then
                    release_flag=true
                    shift
                    if [[ "$deploy_version" != "HEAD" ]]; then
                        echo "Warning: deploy version ignored, using current HEAD."
                        deploy_version="HEAD"
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--debug" ]]; then
                    debug_flag=true
                    shift
                fi
                ;;
            destroy)
                destroy_flag=true
                shift
                if [[ $# -gt 0 && $1 == "--environment" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        environment_name="$2"
                        shift 2
                    else
                        echo "Error: --environment requires a value."
                        usage
                    fi
                fi
                if [[ $# -gt 0 && $1 == "--debug" ]]; then
                    debug_flag=true
                    shift
                fi
                ;;
            --help)
                help
                ;;
            *)
                echo "Error: Unknown flag '$1'."
                usage
                ;;
        esac
    done
}

main() {
    check_deps
    parse_args "$@"
    rm -rf /tmp/dist &> /dev/null

    local playbook="playbooks/"
    local flags=(
        "--extra-vars" "env=${environment_name}"
        "--extra-vars" "profile=${profile_name}"
    )
    [[ "${debug_flag}" == true ]] && flags+=("-vvv")

    case true in
        "${init_flag}")
            playbook+="init.yml"
            [[ "${skip_certificates_setup_flag}" == true ]] && flags+=("--skip-tags" "certs")
            ;;
        "${deploy_flag}")
            playbook+="deploy.yml"
            [[ "${deploy_version}" != "HEAD" ]] && flags+=("--extra-vars" "version=${deploy_version}")
            [[ "${no_logs_collection_flag}" == true ]] && flags+=("--extra-vars" "no_logs_collection=true")
            [[ "${force_build_templates_flag}" == true ]] && flags+=("--extra-vars" "build_templates=true")
            [[ -n "${datadog_key}" ]] && flags+=("--extra-vars" "datadog_key=${datadog_key}")
            [[ -n "${l1_rpc_urls}" ]] && flags+=("--extra-vars" "l1_rpc_urls=${l1_rpc_urls}")
            [[ -n "${oracle_relay_urls}" ]] && flags+=("--extra-vars" "oracle_relay_urls=${oracle_relay_urls}")
            [[ -n "${etherscan_api_key}" ]] && flags+=("--extra-vars" "etherscan_api_key=${etherscan_api_key}")
            [[ -n "${otel_collector_endpoint_url}" ]] && flags+=("--extra-vars" "otel_collector_endpoint_url=${otel_collector_endpoint_url}")
            [[ -n "${genesis_file_url}" ]] && flags+=("--extra-vars" "genesis_file_url=${genesis_file_url}")
            [[ -n "${geth_bootnode_url}" ]] && flags+=("--extra-vars" "geth_bootnode_url=${geth_bootnode_url}")
            [[ -n "${settlement_rpc_url}" ]] && flags+=("--extra-vars" "settlement_rpc_url=${settlement_rpc_url}")
            [[ -n "${contracts_json_url}" ]] && flags+=("--extra-vars" "contracts_json_url=${contracts_json_url}")
            [[ "${release_flag}" == true ]] && flags+=("--extra-vars" "release=true")
            ;;
        "${destroy_flag}")
            playbook+="destroy.yml"
            ;;
        *)
            usage
            ;;
    esac

    ansible-playbook "${playbook}" "${flags[@]}"
}

main "$@"
