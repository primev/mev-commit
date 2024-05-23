#!/usr/bin/env bash

export OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES

init_flag=false
deploy_flag=false
destroy_flag=false
debug_flag=false
no_logs_collection_flag=false
force_build_templates_flag=false
skip_certificates_setup_flag=false
deploy_version="HEAD"
profile_name="devnet"

help() {
    echo "Usage:"
    echo "$0 [init [--profile <name=devnet>] [--skip-certificates-setup] [--debug]]"
    echo "$0 [deploy [version=HEAD] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--debug]]"
    echo "$0 [destroy [--debug]] [--help]"
    echo "$0 --help"
    echo
    echo "Parameters:"
    echo "  init                            Initialize the environment."
    echo "    --profile <name=devnet>       Specify the profile to use (default is devnet)."
    echo "    --skip-certificates-setup     Skip the certificates installation and setup."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  deploy [version=HEAD]           Deploy the specified artifact version (a git commit hash or an existing AWS S3 tag). If not specified or set to HEAD, a local build is triggered."
    echo "    --profile <name=devnet>       Specify the profile to use (default is devnet)."
    echo "    --force-build-templates       Force the build of all job templates before deployment."
    echo "    --no-logs-collection          Disable the collection of logs from deployed jobs."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  destroy                         Destroy the environment."
    echo "    --debug                       Enable debug mode for detailed output."
    echo
    echo "  --help                          Display this help message."
    echo
    echo "Examples:"
    echo "  Initialize with default profile:"
    echo "    $0 init"
    echo
    echo "  Initialize with a specific profile:"
    echo "    $0 init --profile testnet"
    echo
    echo "  Initialize with a specific profile and skip certificates setup:"
    echo "    $0 init --profile testnet --skip-certificates-setup"
    echo
    echo "  Initialize with a specific profile in debug mode:"
    echo "    $0 init --profile testnet --debug"
    echo
    echo "  Deploy the current vcs version and profile:"
    echo "    $0 deploy"
    echo
    echo "  Deploy with a specific version:"
    echo "    $0 deploy 5266b68"
    echo
    echo "  Deploy with a specific version and profile:"
    echo "    $0 deploy v0.1.0 --profile testnet"
    echo
    echo "  Deploy with a specific version and profile and force to build all job templates:"
    echo "    $0 deploy v0.1.0 --profile testnet --force-build-templates"
    echo
    echo "  Deploy with a specific version and profile in debug mode with disabled logs collection:"
    echo "    $0 deploy v0.1.0 --profile testnet --no-logs-collection --debug"
    echo
    echo "  Destroy with debug mode:"
    echo "    $0 destroy --debug"
    exit 1
}

usage() {
        echo "Usage:"
        echo "$0 [init [--profile <name=devnet>] [--skip-certificates-setup] [--debug]]"
        echo "$0 [deploy [version=HEAD] [--profile <name=devnet>] [--force-build-templates] [--no-logs-collection] [--debug]]"
        echo "$0 [destroy [--debug]] [--help]"
        echo "$0 --help"
    exit 1
}

check_deps() {
    local missing_utils=()
    local required_utilities=(
        aws
        ansible
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

    if [[ ! -f hosts.ini ]]; then
        echo "Error: hosts.ini file not found."
        exit 1
    fi
}

parse_args() {
    while [[ $# -gt 0 ]]; do
        key="$1"
        case $key in
            init)
                if $init_flag || $deploy_flag || $destroy_flag; then
                    echo "Error: Only one of 'init', 'deploy', or 'destroy' can be specified."
                    usage
                fi
                init_flag=true
                shift
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
                if $init_flag || $deploy_flag || $destroy_flag; then
                    echo "Error: Only one of 'init', 'deploy', or 'destroy' can be specified."
                    usage
                fi
                deploy_flag=true
                if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                    deploy_version="$2"
                    shift
                fi
                shift
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
                if [[ $# -gt 0 && $1 == "--debug" ]]; then
                    debug_flag=true
                    shift
                fi
                ;;
            destroy)
                if $init_flag || $deploy_flag || $destroy_flag; then
                    echo "Error: Only one of 'init', 'deploy', or 'destroy' can be specified."
                    usage
                fi
                destroy_flag=true
                shift
                if [[ $# -gt 0 && $1 == "--debug" ]]; then
                    debug_flag=true
                    shift
                fi
                ;;
            --profile)
                if [[ $# -gt 1 && ! $2 =~ ^-- ]]; then
                    profile_name="$2"
                    shift
                else
                    echo "Error: --profile requires a value."
                    usage
                fi
                shift
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
    local flags=("--extra-vars" "profile=${profile_name}")
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
            ;;
        "${destroy_flag}")
            playbook+="destroy.yml"
            ;;
        *)
            usage
            ;;
    esac

    ansible-playbook --inventory hosts.ini "${playbook}" "${flags[@]}"
}

main "$@"
