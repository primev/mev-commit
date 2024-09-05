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
backup_flag=false
deploy_version="HEAD"
environment_name="devenv"
profile_name="devnet"
datadog_key=""
l1_rpc_url=""
otel_collector_endpoint_url=""
genesis_file_url=""
geth_bootnode_url=""

help() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--l1-rpc-url <url>] [--otel-collector-endpoint-url <url>] [--genesis-file-url <url>] [--geth-bootnode-url <url>] [--release] [--debug]]"
    echo "$0 [destroy [--backup] [--debug]]"
    echo "$0 --help"
    echo
    echo "Parameters:"
    echo "  init                           Initialize the environment."
    echo "    --environment <name=devenv>  Specify the environment to use (default is devenv)."
    echo "    --skip-certificates-setup    Skip the certificates installation and setup."
    echo "    --debug                      Enable debug mode for detailed output."
    echo
    echo "  deploy [version=HEAD]                  Deploy the specified artifact version (a git commit hash or an existing AWS S3 tag). If not specified or set to HEAD, a local build is triggered."
    echo "    --environment <name=devenv>]         Specify the environment to use (default is devenv)."
    echo "    --profile <name=devnet>]             Specify the profile to use (default is devnet)."
    echo "    --force-build-templates              Force the build of all job templates before deployment."
    echo "    --no-logs-collection                 Disable the collection of logs from deployed jobs."
    echo "    --datadog-key <key>]                 Datadog API key, cannot be empty."
    echo "    --l1-rpc-url <url>]                  L1 RPC URL, cannot be empty."
    echo "    --otel-collector-endpoint-url <url>] OpenTelemetry Collector Endpoint URL, cannot be empty."
    echo "    --genesis-file-url <url>]            URL to the genesis file, cannot be empty."
    echo "    --geth-bootnode-url <url>]           URL to the Geth bootnode, cannot be empty."
    echo "    --release                            It will ignore the specified deployment version and use the current HEAD tag as the build version."
    echo "    --debug                              Enable debug mode for detailed output."
    echo
    echo "  destroy    Destroy the whole cluster."
    echo "    --backup  Create a backup before destroying the environment."
    echo "    --debug   Enable debug mode for detailed output."
    echo
    echo "  --help  Display this help message."
    echo
    echo "Examples:"
    echo "  Initialize with default environment and profile:"
    echo "    $0 init"
    echo
    echo "  Initialize with a specific environment and profile:"
    echo "    $0 init --environment devenv --profile testnet"
    echo
    echo "  Initialize with a specific environment, profile and skip certificates setup:"
    echo "    $0 init --environment devenv --profile testnet --skip-certificates-setup"
    echo
    echo "  Initialize with a specific environment, profile in debug mode:"
    echo "    $0 init --environment devenv --profile testnet --debug"
    echo
    echo "  Deploy the current vcs version, environment and profile:"
    echo "    $0 deploy"
    echo
    echo "  Deploy with a specific version, environment and profile:"
    echo "    $0 deploy v0.1.0 --environment devenv --profile testnet"
    echo
    echo "  Deploy with a specific version, environment, profile and force to build all job templates:"
    echo "    $0 deploy v0.1.0 --environment devenv --profile testnet --force-build-templates"
    echo
    echo "  Deploy with a specific version, environment, profile in debug mode with disabled logs collection, Datadog API key, L1 RPC URL, OpenTelemetry Collector Endpoint URL, genesis file URL, and Geth bootnode URL:"
    echo "    $0 deploy v0.1.0 --environment devenv --profile testnet --no-logs-collection --datadog-key your_datadog_key --l1-rpc-url your_rpc_url --otel-collector-endpoint-url your_otel_url --genesis-file-url your_genesis_file_url --geth-bootnode-url your_geth_bootnode_url --debug"
    echo
    echo "  Destroy all jobs but backup before do so:"
    echo "    $0 destroy --backup --debug"
    exit 1
}

usage() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--l1-rpc-url <url>] [--otel-collector-endpoint-url <url>] [--genesis-file-url <url>] [--geth-bootnode-url <url>] [--release] [--debug]]"
    echo "$0 [destroy [--backup] [--debug]]"
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
                if [[ $# -gt 0 && $1 == "--l1-rpc-url" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        l1_rpc_url="$2"
                        shift 2
                    else
                        echo "Error: --l1-rpc-url requires a value."
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
                if [[ $# -gt 0 && $1 == "--backup" ]]; then
                    backup_flag=true
                    shift
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
            [[ -n "${l1_rpc_url}" ]] && flags+=("--extra-vars" "l1_rpc_url=${l1_rpc_url}")
            [[ -n "${otel_collector_endpoint_url}" ]] && flags+=("--extra-vars" "otel_collector_endpoint_url=${otel_collector_endpoint_url}")
            [[ -n "${genesis_file_url}" ]] && flags+=("--extra-vars" "genesis_file_url=${genesis_file_url}")
            [[ -n "${geth_bootnode_url}" ]] && flags+=("--extra-vars" "geth_bootnode_url=${geth_bootnode_url}")
            [[ "${release_flag}" == true ]] && flags+=("--extra-vars" "release=true")
            ;;
        "${destroy_flag}")
            playbook+="destroy.yml"
            [[ "${backup_flag}" == true ]] && flags+=("--extra-vars" "backup=true")
            ;;
        *)
            usage
            ;;
    esac

    ansible-playbook "${playbook}" "${flags[@]}"
}

main "$@"
