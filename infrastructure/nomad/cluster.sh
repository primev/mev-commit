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

help() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--profile <name=devnet>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--release] [--debug]]"
    echo "$0 [destroy [--debug]] [--help]"
    echo "$0 --help"
    echo
    echo "Parameters:"
    echo "  init                            Initialize the environment."
    echo "    --environment <name=devenv>   Specify the environment to use (default is devenv)."
    echo "    --profile <name=devnet>       Specify the profile to use (default is devnet)."
    echo "    --skip-certificates-setup     Skip the certificates installation and setup."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  deploy [version=HEAD]           Deploy the specified artifact version (a git commit hash or an existing AWS S3 tag). If not specified or set to HEAD, a local build is triggered."
    echo "    --environment <name=devenv>   Specify the environment to use (default is devenv)."
    echo "    --profile <name=devnet>       Specify the profile to use (default is devnet)."
    echo "    --force-build-templates       Force the build of all job templates before deployment."
    echo "    --no-logs-collection          Disable the collection of logs from deployed jobs."
    echo "    --datadog-key <key>           Datadog API key."
    echo "    --release                     It will ignore the specified deployment version and use the current HEAD tag as the build version."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  destroy                         Destroy the whole cluster."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  --help                          Display this help message."
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
    echo "  Deploy with a specific version, environment, profile in debug mode with disabled logs collection and Datadog API key:"
    echo "    $0 deploy v0.1.0 --environment devenv --profile testnet --no-logs-collection --datadog-key your_datadog_key --debug"
    echo
    echo "  Destroy with specific environment and debug mode:"
    echo "    $0 destroy --environment devenv --debug"
    exit 1
}

usage() {
    echo "Usage:"
    echo "$0 [init [--environment <name=devenv>] [--profile <name=devnet>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--environment <name=devenv>] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--datadog-key <key>] [--release] [--debug]]"
    echo "$0 [destroy [--debug]] [--help]"
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
                if [[ $# -gt 0 && $1 == "--profile" ]]; then
                    if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                        profile_name="$2"
                        shift 2
                    else
                        echo "Error: --profile requires a value."
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
