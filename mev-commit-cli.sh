#!/bin/bash

# Default RPC URL and Paths
L1_RPC_BASE_URL=https://sepolia.infura.io/v3
DEFAULT_RPC_URL="http://sl-bootnode:8545"
PRIMEV_DIR="$HOME/.primev"
DEFAULT_CHAIN_ID="17864"

# Default Docker network name
DOCKER_NETWORK_NAME="primev_net"

# Default values for optional arguments
rpc_url=$DEFAULT_RPC_URL
datadog_key=""
command=""

# Function to initialize the environment
initialize_environment() {
    create_docker_network
}

# Function to create a Docker network
create_docker_network() {
    echo "Creating Docker network: $DOCKER_NETWORK_NAME..."
    if ! docker network inspect $DOCKER_NETWORK_NAME >/dev/null 2>&1; then
        docker network create --driver bridge --subnet 172.29.0.0/16 $DOCKER_NETWORK_NAME
    else
        echo "Network $DOCKER_NETWORK_NAME already exists."
    fi
}

start_settlement_layer() {
    local datadog_key=$1

    cat > "external/geth/geth-poa/.env" <<EOF
    HYPERLANE_DEPLOYER_PRIVATE_KEY=0xc065f4c9a6dda0785e2224f5af8e473614de1c029acf094f03d5830e2dd5b0ea
    NODE1_PRIVATE_KEY=0xe82a054e06f89598485134b4f2ce8a612ce7f7f7e14e650f9f20b30efddd0e57
    NODE2_PRIVATE_KEY=0xb17b77fe56797c1a6c236f628d25ede823496af371b3fec858a7a6beff07696b
    RELAYER_PRIVATE_KEY=0xa0d74f611ee519f3fd4a84236ee24b955df2a3f40632f404ca46e0b17f696df3
    NEXT_PUBLIC_WALLET_CONNECT_ID=0x074ac60cba235536b25b262f66dee686
    DD_KEY=${datadog_key}
EOF

    export AGENT_BASE_IMAGE=nil
    export L2_NODE_URL=nil

    DD_KEY=nil docker compose --profile settlement -f "external/geth/geth-poa/docker-compose.yml" up -d --build

    # Wait for settlement layer to be up before deploying create2
    sleep 10

    deploy_create2
}

stop_settlement_layer() {
    DD_KEY=nil docker compose --profile settlement -f "external/geth/geth-poa/docker-compose.yml" down
}

start_mev_commit_minimal() {
    echo "Starting MEV-Commit..."
    docker compose --profile minimal-setup -f "p2p/integration-compose.yml" up --build -d
}

start_mev_commit_e2e() {
    local datadog_key=""
    local sepolia_key=""
    echo "Starting MEV-Commit..."

    # Loop through arguments and process them
    for arg in "$@"
    do
        case $arg in
            --datadog-key=*)
            datadog_key="${arg#*=}"
            shift # Remove --datadog-key= from processing
            ;;
            --sepolia-key=*)
            sepolia_key="${arg#*=}"
            shift # Remove --sepolia-key= from processing
            ;;
            *)
            # Unknown option
            ;;
        esac
    done
    echo "Setting .env file ..."

        # Create or overwrite the .env file
    cat > "p2p/integrationtest/.env" <<EOF
    BIDDER_REGISTRY=0x02CcEcB19c6D7EFe583C8b97022cB4b4C0B65608
    PROVIDER_REGISTRY=0x070cE6161AD79a3BC7aEa222FdfC6AD171Ca83F3
    PRECONF_CONTRACT=0x4DfF34f74aE5C48a5050eb54e7cEDAb9DEF03715

    RPC_URL=${rpc_url}
    PRIVATE_KEY=ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
    L1_RPC_URL="${L1_RPC_BASE_URL}/${sepolia_key}"
EOF


    # Check if datadog_key is empty
    if [ -z "$datadog_key" ]; then
        echo "DD_KEY is empty, so no agents will be started."
        # Run Docker Compose without --profile agent
        docker compose --profile e2etest -f "p2p/e2e-compose.yml" up --build -d
    else
        # Run Docker Compose with --profile agent
        DD_KEY="$datadog_key" docker compose --profile e2etest --profile agent -f "p2p/e2e-compose.yml" up --build -d
    fi
}

start_mev_commit() {
    local datadog_key=$1

    echo "Starting MEV-Commit..."

    # Check if datadog_key is empty
    if [ -z "$datadog_key" ]; then
        echo "DD_KEY is empty, so no agents will be started."
        # Run Docker Compose without --profile agent
        docker compose --profile integration-test -f "p2p/integration-compose.yml" up --build -d
    else
        # Run Docker Compose with --profile agent
        DD_KEY="$datadog_key" docker compose --profile integration-test --profile agent -f "p2p/integration-compose.yml" up --build -d
    fi
}

# Builds contract-deployer image only if it doesn't already exist
build_contract_deployer() {
    docker build -t contract-deployer "contracts"
}

deploy_create2() {
    local rpc_url=${1:-$DEFAULT_RPC_URL}
    local network_name=${2:-"$DOCKER_NETWORK_NAME"}
    chmod +x "external/geth/geth-poa/util/deploy_create2.sh"
    docker run \
        --rm \
        --network "$network_name" \
        -v "$(pwd)/external/geth/geth-poa/util/deploy_create2.sh:/deploy_create2.sh" \
        alpine /bin/sh -c \
        "apk add --no-cache curl jq \
        && /deploy_create2.sh ${rpc_url}"
}

deploy_contracts() {
    local rpc_url=${1:-$DEFAULT_RPC_URL}
    local chain_id=${2:-$DEFAULT_CHAIN_ID}
    local private_key=${3:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"} # Default private key

    # Wait for the Geth POA network to be up and running
    echo "Waiting for Geth POA network to be fully up..."
    sleep 10

    build_contract_deployer

    # Run the Docker container to deploy the contracts
    echo "Deploying Contracts with RPC URL: $rpc_url, Chain ID: $chain_id, and Private Key: [HIDDEN]"
    docker run --rm --network "$DOCKER_NETWORK_NAME" \
        -e RPC_URL="$rpc_url" \
        -e CHAIN_ID="$chain_id" \
        -e PRIVATE_KEY="$private_key" \
        contract-deployer
}

start_oracle(){
    local l1_url=${L1_RPC_BASE_URL}/$1
    local datadog_key=$2
    # Run Docker Compose
    L1_URL="$l1_url" DD_KEY="$datadog_key" docker compose -f "oracle/integration-compose.yml" up -d --build
}

stop_oracle(){
    # Run Docker Compose
    docker compose -f "oracle/integration-compose.yml" down
}

start_hyperlane(){
    local public_rpc_url=${1:-$DEFAULT_RPC_URL}
    local rpc_url=${2:-$DEFAULT_RPC_URL}
    local chain_id=${3:-$DEFAULT_CHAIN_ID}
    local private_key=${4:-"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"}
    AGENT_BASE_IMAGE=gcr.io/abacus-labs-dev/hyperlane-agent@sha256:854f92966eac6b49e5132e152cc58168ecdddc76c2d390e657b81bdaf1396af0 \
        PUBLIC_SETTLEMENT_RPC_URL="$public_rpc_url" \
        SETTLEMENT_RPC_URL="$rpc_url" \
        docker compose -f "bridge/hyperlane/docker-compose.yml" --profile bridge up -d --build

    # Run Alpine container which:
    # 1. Install jq
    # 2. Prints warp-deployment.json from docker volume
    # 3. Parses the JSON to get hyperlane ERC20 contract address "router" deployed on settlement layer
    HYP_ERC20_ADDR=$(docker run --rm -v hyperlane_hyperlane-deploy-artifacts:/data alpine /bin/sh -c "apk add --no-cache -q jq && cat /data/warp-deployment.json | jq -r '.mevcommitsettlement.router'")

    echo "HYP_ERC20_ADDR: $HYP_ERC20_ADDR"

    build_contract_deployer

    # Deploy whitelist contract 
    docker run --rm --network "$DOCKER_NETWORK_NAME" \
        -e RPC_URL="$rpc_url" \
        -e CHAIN_ID="$chain_id" \
        -e PRIVATE_KEY="$private_key" \
        -e DEPLOY_TYPE="whitelist" \
        -e HYP_ERC20_ADDR="$HYP_ERC20_ADDR" \
        contract-deployer
}

stop_hyperlane(){
    AGENT_BASE_IMAGE=gcr.io/abacus-labs-dev/hyperlane-agent@sha256:854f92966eac6b49e5132e152cc58168ecdddc76c2d390e657b81bdaf1396af0 PUBLIC_SETTLEMENT_RPC_URL="$public_rpc_url" SETTLEMENT_RPC_URL="$rpc_url" docker compose -f "bridge/hyperlane/docker-compose.yml" --profile bridge down
}

start_local_l1() {
    DD_KEY=nil docker compose --profile local_l1 -f "external/geth/geth-poa/docker-compose.yml" up -d --build
    # wait for l1 to be up before deploying create2
    sleep 10
    deploy_create2 "http://l1-bootnode:8545" "geth-poa_l1_net"
}

stop_local_l1() {
    DD_KEY=nil docker compose --profile local_l1 -f "external/geth/geth-poa/docker-compose.yml" down
}

deploy_standard_bridge_contracts() {
    local rpc_url=$DEFAULT_RPC_URL
    local chain_id=$DEFAULT_CHAIN_ID
    # Pre-funded account for deployment on settlement layer
    local private_key="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

    build_contract_deployer

    docker run --rm --network "$DOCKER_NETWORK_NAME" \
        -e RPC_URL="$rpc_url" \
        -e CHAIN_ID="$chain_id" \
        -e PRIVATE_KEY="$private_key" \
        -e DEPLOY_TYPE="settlement-gateway" \
        -e RELAYER_ADDR="0x0DCaa27B9E4Db92F820189345792f8eC5Ef148F6" \
        contract-deployer 

    local settlement_gateway_addr="0xc1f93bE11D7472c9B9a4d87B41dD0a491F1fbc75"
    docker run --rm --network "$DOCKER_NETWORK_NAME" \
        -e RPC_URL="$rpc_url" \
        -e CHAIN_ID="$chain_id" \
        -e PRIVATE_KEY="$private_key" \
        -e DEPLOY_TYPE="whitelist" \
        -e HYP_ERC20_ADDR="$settlement_gateway_addr" \
        contract-deployer

    docker run --rm --network "geth-poa_l1_net" \
        -e RPC_URL="http://l1-bootnode:8545" \
        -e CHAIN_ID="39999" \
        -e PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" \
        -e DEPLOY_TYPE="l1-gateway" \
        -e RELAYER_ADDR="0x0DCaa27B9E4Db92F820189345792f8eC5Ef148F6" \
        contract-deployer
    
    echo "Standard bridge gateway contracts finished deploying on both chains"
}

start_standard_bridge() {
    echo "Starting standard bridge..."
    docker compose -f "bridge/standard/bridge-v1/docker-compose.yml" up --build -d
}

stop_standard_bridge() {
    echo "Stopping standard bridge..."
    docker compose -f "bridge/standard/bridge-v1/docker-compose.yml" down
}

clean() {
    echo "Cleaning up..."
    # Docker cleanup script
    echo "Stopping all Docker containers..."
    docker ps -aq | xargs -r docker stop

    echo "Removing all Docker containers..."
    docker ps -aq | xargs -r docker rm

    echo "Removing all Docker images..."
    docker images -q | xargs -r docker rmi

    echo "Removing all Docker volumes..."
    docker volume ls -q | xargs -r docker volume rm

    echo "Removing all Docker networks..."
    docker network ls | grep "bridge\|none\|host" -v | awk '{if(NR>1)print $1}' | xargs -r -I {} docker network rm {}

    echo "Pruning Docker system..."
    docker system prune -a -f --volumes

    echo "Docker cleanup complete."
}

stop_services() {
    service=$1
    echo "Stopping Docker Compose services..."

    case $service in
        "sl")
            stop_settlement_layer
            ;;
        "oracle")
            stop_oracle  # Assuming stop_oracle is a function you've defined elsewhere
            ;;
        "hyperlane")
            stop_hyperlane
            ;;
        "mev-commit")
            docker compose -f "p2p/integration-compose.yml" down
            ;;
        "local_l1")
            stop_local_l1
            ;;
        "standard_bridge")
            stop_standard_bridge
            ;;
        "all")
            stop_settlement_layer
            stop_oracle
            stop_hyperlane
            stop_local_l1
            stop_standard_bridge
            docker compose -f "p2p/integration-compose.yml" down
            ;;
        *)
            echo "Invalid service: $service"
            echo "Valid services: sl, oracle, mev-commit, hyperlane, local_l1, standard_bridge, all"
            return 1
    esac

    echo "Service(s) stopped."
}


start_service() {
    local service_name=$1
    case $service_name in
        "all")
            start_settlement_layer "$datadog_key"
            deploy_contracts "$rpc_url"
            start_mev_commit "$datadog_key"
            start_oracle "$sepolia_key" "$datadog_key"
            start_hyperlane "$public_rpc_url"
            ;;
        "e2e")
            initialize_environment
            start_settlement_layer "$datadog_key"
            deploy_contracts "$rpc_url"
            start_mev_commit_e2e "--sepolia-key=$sepolia_key" "--datadog-key=$datadog_key"
            sleep 12
            start_oracle "$sepolia_key" "$datadog_key"
            start_hyperlane "$public_rpc_url"
            ;;
        "mev-commit")
            initialize_environment
            start_mev_commit "$datadog_key"
            ;;
        "oracle")
            start_oracle "$sepolia_key" "$datadog_key"
            ;;
        "sl")
            start_settlement_layer "$datadog_key"
            ;;
        "hyperlane")
            start_hyperlane "$public_rpc_url"
            ;;
        "minimal")
            initialize_environment
            start_settlement_layer "$datadog_key"
            deploy_contracts "$rpc_url"
            start_mev_commit_minimal
            ;;
        "local_l1")
            start_local_l1
            ;;
        "standard_bridge")
            create_docker_network
            start_settlement_layer
            start_local_l1
            deploy_standard_bridge_contracts
            start_standard_bridge
            ;;
        *)
            echo "Invalid service name: $service_name"
            echo "Valid services: all, e2e, oracle, sl, hyperlane, minimal, local_l1, standard_bridge"
            return 1
            ;;
    esac
}

# Function to display help
show_help() {
    echo "Usage: $0 [command] [service(s)] [options]"
    echo ""
    echo "Commands:"
    echo "  deploy_contracts       Deploy contracts"
    echo "  start [services]       Start specified services. Available services: all, e2e, mev-commit, oracle, sl, hyperlane, minimal, local_l1, standard_bridge"
    echo "  stop [service]         Stop specified service. Available services: sl, mev-commit, hyperlane, local_l1, standard_bridge, all"
    echo "  update                 Update repositories"
    echo "  clean                  Cleanup Docker"
    echo ""
    echo "Options:"
    echo "  -h, --help             Show this help message"
    echo "  --rpc-url URL          Set the internal RPC URL for mev-commit-geth"
    echo "  --public-rpc-url URL   Set the public RPC URL for mev-commit-geth"
    echo "  --datadog-key KEY      Set the Datadog key"
    echo "  --sepolia-key KEY      Set the Sepolia key"
    echo ""
    echo "Examples:"
    echo "  $0 start all --rpc-url http://localhost:8545  Start all services with a specific RPC URL"
    echo "  $0 start e2e --datadog-key abc123             Start only the e2e service with a Datadog key"
    echo "  $0 start oracle                               Start only the oracle service"
    echo "  $0 start sl                                   Start only the settlement layer service"
    echo "  $0 stop sl                                    Stop the settlement layer service"
    echo ""
}

# Parse command line options
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        --rpc-url)
            rpc_url="$2"
            shift 2
            ;;
        --public-rpc-url)
            public_rpc_url="$2"
            shift 2
            ;;
        --datadog-key)
            datadog_key="$2"
            shift 2
            ;;
        --sepolia-key)
            sepolia_key="$2"
            shift 2
            ;;
        start|stop|deploy_contracts|update|clean|init_network)
            command="$1"
            shift
            # If additional arguments are present after the command, they are captured as service names or additional options
            service_names=()
            while [[ "$#" -gt 0 ]] && [[ "$1" != "--"* ]]; do
                service_names+=("$1")
                shift
            done
            ;;
        *)
            echo "Invalid option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Check if a command has been specified
if [[ -z "$command" ]]; then
    echo "No command specified."
    show_help
    exit 1
fi

# Main script logic based on the command variable
case "$command" in
    init_network)
        create_docker_network
        ;;
    start)
        if [ ${#service_names[@]} -eq 0 ]; then
            echo "No service specified. Starting all services."
            start_service "all"
        else
            for service_name in "${service_names[@]}"; do
                start_service "$service_name"
            done
        fi
        ;;
    deploy_contracts)
        deploy_contracts "$rpc_url"
        ;;
    stop)
        if [ -z "${service_names[0]}" ]; then
            echo "No service specified for stopping."
            exit 1
        else
            stop_services "${service_names[0]}"
        fi
        ;;
    update)
        update_repos
        ;;
    clean)
        clean
        ;;
esac

exit 0
