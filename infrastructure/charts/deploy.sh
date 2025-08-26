#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

# Defaults
NAMESPACE="default"
TIMEOUT="600"
RPC_PORT="8545"
DEPLOY_SCRIPT="deploy-contracts.py"
DRY_RUN=false
CLEANUP=false
SKIP_CONTRACTS=false
PASSWORD=""
WORK_DIR=""
SELECTED_CHARTS=""

# Parse args
while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run) DRY_RUN=true; shift ;;
        --cleanup) CLEANUP=true; shift ;;
        --skip-contracts) SKIP_CONTRACTS=true; shift ;;
        --namespace) NAMESPACE="$2"; shift 2 ;;
        --password) PASSWORD="$2"; shift 2 ;;
        --work-dir) WORK_DIR="$2"; shift 2 ;;
        --charts) SELECTED_CHARTS="$2"; shift 2 ;;
        --help) 
            echo "Usage: $0 [--dry-run] [--cleanup] [--skip-contracts] [--password PASS] [--charts CHARTS]"
            echo "Charts: mock-l1,erigon-snode,relay-emulator,dashboard"
            exit 0 ;;
        *) print_error "Unknown option: $1"; exit 1 ;;
    esac
done

# Cleanup function
cleanup_all() {
    print_info "Cleaning up all releases..."
    
    for release in "mev-commit-dashboard" "mev-commit-relay-emulator" "erigon-snode" "erigon-mev-commit-mock-l1"; do
        if helm list -n $NAMESPACE | grep -q "^$release"; then
            print_info "Deleting $release"
            helm uninstall "$release" -n "$NAMESPACE" || true
        fi
    done
    
    sleep 5
    print_info "Deleting PVCs..."
    kubectl get pvc -n "$NAMESPACE" --no-headers 2>/dev/null | grep -E "(erigon|geth|mock)" | awk '{print $1}' | xargs -r kubectl delete pvc -n "$NAMESPACE" || true
    
    print_success "Cleanup completed!"
}

if [[ "$CLEANUP" == true ]]; then
    cleanup_all
    exit 0
fi

# Find pod by app label
find_pod() {
    APP_NAME="$1"
    
    # Handle P2P nodes which all use mev-commit-p2p as app name
    case "$APP_NAME" in
        erigon-mev-commit-bootnode) 
            kubectl get pods -l "app.kubernetes.io/name=mev-commit-p2p,app.kubernetes.io/component=bootnode" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
            ;;
        erigon-mev-commit-bidder)
            kubectl get pods -l "app.kubernetes.io/name=mev-commit-p2p,app.kubernetes.io/component=bidder" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
            ;;
        erigon-mev-commit-provider)
            kubectl get pods -l "app.kubernetes.io/name=mev-commit-p2p,app.kubernetes.io/component=provider" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
            ;;
        mev-commit-emulator-bt)
            kubectl get pods -l "app.kubernetes.io/name=bidder-emulator" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
            ;;
        *)
            kubectl get pods -l "app.kubernetes.io/name=$APP_NAME" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
            ;;
    esac
}

# Find service by app label
find_service() {
    kubectl get svc -l "app.kubernetes.io/name=$1" -n "$NAMESPACE" --no-headers 2>/dev/null | awk 'NR==1{print $1}' || true
}

# Wait for pod ready
wait_pod_ready() {
    APP_NAME="$1"
    MAX_WAIT=30
    POD_NAME=""
    
    print_info "Waiting for $APP_NAME pod..."
    for i in $(seq 1 $MAX_WAIT); do
        POD_NAME=$(find_pod "$APP_NAME")
        [[ -n "$POD_NAME" ]] && break
        sleep 2
    done
    
    if [[ -z "$POD_NAME" ]]; then
        print_error "Pod for $APP_NAME not found"
        return 1
    fi
    
    print_info "Waiting for $POD_NAME to be ready..."
    if kubectl wait pod/"$POD_NAME" --for=condition=Ready --timeout="${TIMEOUT}s" -n "$NAMESPACE"; then
        print_success "$POD_NAME ready!"
    else
        print_error "$POD_NAME failed to be ready"
        return 1
    fi
}

# Deploy chart
deploy_chart() {
    CHART_NAME="$1"
    RELEASE_NAME="$2"
    CHART_PATH="$3"
    shift 3
    EXTRA_ARGS=("$@")
    
    print_info "Deploying $CHART_NAME..."
    
    HELM_CMD=(helm install "$RELEASE_NAME" "$CHART_PATH" -n "$NAMESPACE")
    [[ "$NAMESPACE" != "default" ]] && HELM_CMD+=(--create-namespace)
    [[ "$DRY_RUN" == true ]] && HELM_CMD+=(--dry-run)
    [[ ${#EXTRA_ARGS[@]} -gt 0 ]] && HELM_CMD+=("${EXTRA_ARGS[@]}")
    
    if [[ "$DRY_RUN" == true ]]; then
        if "${HELM_CMD[@]}" >/dev/null 2>&1; then
            print_success "$CHART_NAME dry-run OK"
        else
            print_error "$CHART_NAME dry-run failed"
            return 1
        fi
    else
        print_info "Running: ${HELM_CMD[*]}"
        if "${HELM_CMD[@]}"; then
            print_success "$CHART_NAME deployed"
            wait_pod_ready "$CHART_NAME"
        else
            print_error "$CHART_NAME deployment failed"
            return 1
        fi
    fi
}

# RPC check
check_rpc() {
    print_info "Checking RPC..."
    RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://127.0.0.1:$RPC_PORT" 2>/dev/null)
    
    BLOCK_HEX=$(echo "$RESPONSE" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
    if [[ -n "$BLOCK_HEX" && "$BLOCK_HEX" != "null" ]]; then
        BLOCK_NUMBER=$((16#${BLOCK_HEX#0x}))
        print_success "RPC working! Block: $BLOCK_NUMBER"
        return 0
    fi
    return 1
}

# Find newest contracts JSON
find_contracts_json() {
    find /tmp /private/tmp -name "core-contracts.json" -path "*/contracts/deploy-summaries/*" 2>/dev/null | xargs ls -t 2>/dev/null | head -1
}

# Extract contract addresses
extract_contracts() {
    JSON_FILE="$1"
    if [[ ! -f "$JSON_FILE" ]]; then
        print_error "Contracts JSON not found: $JSON_FILE"
        return 1
    fi
    
    if command -v jq >/dev/null 2>&1; then
        BIDDER_REGISTRY=$(jq -r '.BidderRegistry // ""' "$JSON_FILE")
        BLOCK_TRACKER=$(jq -r '.BlockTracker // ""' "$JSON_FILE")
        ORACLE=$(jq -r '.Oracle // ""' "$JSON_FILE")
        PRECONF_MANAGER=$(jq -r '.PreconfManager // ""' "$JSON_FILE")
        PROVIDER_REGISTRY=$(jq -r '.ProviderRegistry // ""' "$JSON_FILE")
    else
        BIDDER_REGISTRY=$(python3 -c "import json; print(json.load(open('$JSON_FILE')).get('BidderRegistry', ''))")
        BLOCK_TRACKER=$(python3 -c "import json; print(json.load(open('$JSON_FILE')).get('BlockTracker', ''))")
        ORACLE=$(python3 -c "import json; print(json.load(open('$JSON_FILE')).get('Oracle', ''))")
        PRECONF_MANAGER=$(python3 -c "import json; print(json.load(open('$JSON_FILE')).get('PreconfManager', ''))")
        PROVIDER_REGISTRY=$(python3 -c "import json; print(json.load(open('$JSON_FILE')).get('ProviderRegistry', ''))")
    fi
    
    if [[ -z "$BIDDER_REGISTRY" || -z "$ORACLE" ]]; then
        print_error "Failed to extract contract addresses"
        return 1
    fi
    
    print_success "Contracts extracted"
}

# Get bootnode connection info
get_bootnode_connection() {
    BOOTNODE_SERVICE="erigon-mev-commit-bootnode-mev-commit-p2p-bootnode"
    MAX_ATTEMPTS=30
    
    print_info "Getting bootnode connection info..."
    
    # Wait for LoadBalancer to get external IP
    print_info "Waiting for LoadBalancer external IP..."
    for i in $(seq 1 $MAX_ATTEMPTS); do
        BOOTNODE_IP=$(kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null)
        if [[ -n "$BOOTNODE_IP" && "$BOOTNODE_IP" != "<nil>" ]]; then
            break
        fi
        print_info "Attempt $i/$MAX_ATTEMPTS: Waiting for LoadBalancer IP..."
        sleep 10
    done
    
    if [[ -z "$BOOTNODE_IP" || "$BOOTNODE_IP" == "<nil>" ]]; then
        print_error "Failed to get bootnode external IP after $((MAX_ATTEMPTS * 10)) seconds"
        kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" || true
        return 1
    fi
    
    print_success "Bootnode IP: $BOOTNODE_IP"
    
    # Get peer ID from bootnode API
    print_info "Getting peer ID from bootnode API..."
    for i in $(seq 1 20); do
        print_info "Attempt $i/20: Calling https://$BOOTNODE_IP:13723/v1/debug/topology"
        
        # Simple curl call matching the manual version
        PEER_ID=$(curl -sk "https://$BOOTNODE_IP:13723/v1/debug/topology" | jq -r '.topology.self.Underlay' 2>/dev/null)
        
        if [[ -n "$PEER_ID" && "$PEER_ID" != "null" ]]; then
            BOOTNODE_CONNECTION="/ip4/$BOOTNODE_IP/tcp/13522/p2p/$PEER_ID"
            print_success "Bootnode connection: $BOOTNODE_CONNECTION"
            return 0
        fi
        
        print_warning "No valid peer ID received, retrying in 10 seconds..."
        
        sleep 15
    done
    
    print_error "Failed to get bootnode peer ID after 20 attempts"
    print_info "Final attempt to check bootnode service and pod status:"
    kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" || true
    kubectl get pods -l "app.kubernetes.io/name=mev-commit-p2p,app.kubernetes.io/component=bootnode" -n "$NAMESPACE" || true
    return 1
}

# Port forward cleanup
cleanup_pf() {
    [[ -n "$PF_PID" ]] && { kill $PF_PID 2>/dev/null || true; wait $PF_PID 2>/dev/null || true; }
}
trap cleanup_pf EXIT

# Determine charts to deploy
if [[ -n "$SELECTED_CHARTS" ]]; then
    IFS=',' read -ra CHARTS <<< "$SELECTED_CHARTS"
else
    CHARTS=("mock-l1" "erigon-snode" "dashboard" "relay-emulator" "oracle" "bootnode" "bidder" "bidder-emulator" "provider" "provider-emulator")
fi

# Validation
for CHART in "${CHARTS[@]}"; do
    case "$CHART" in
        mock-l1) CHART_PATH="./mev-commit-geth-l1" ;;
        erigon-snode) CHART_PATH="./mev-commit-erigon" ;;
        relay-emulator) CHART_PATH="./mev-commit-relay-emulator" ;;
        oracle) CHART_PATH="./mev-commit-oracle" ;;
        dashboard) CHART_PATH="./mev-commit-dashboard" ;;
        bootnode|bidder|provider) CHART_PATH="./mev-commit-p2p" ;;
        bidder-emulator) CHART_PATH="./mev-commit-emulator-bt" ;;
        provider-emulator) CHART_PATH="./mev-commit-emulator" ;;
        *) print_error "Unknown chart: $CHART"; exit 1 ;;
    esac
    [[ ! -d "$CHART_PATH" ]] && { print_error "Chart not found: $CHART_PATH"; exit 1; }
done

# Check password requirement
if [[ "$SKIP_CONTRACTS" == false && "$DRY_RUN" == false && -z "$PASSWORD" ]]; then
    for CHART in "${CHARTS[@]}"; do
        [[ "$CHART" == "erigon-snode" ]] && { print_error "Password required for contracts"; exit 1; }
    done
fi

print_info "Deploying charts: ${CHARTS[*]}"

# Initialize contract variables
BIDDER_REGISTRY=""
BLOCK_TRACKER=""
ORACLE=""
PRECONF_MANAGER=""
PROVIDER_REGISTRY=""

# Deploy charts
for CHART in "${CHARTS[@]}"; do
    case "$CHART" in
        mock-l1)
            deploy_chart "mev-commit-mock-l1" "erigon-mev-commit-mock-l1" "./mev-commit-geth-l1"
            ;;
        erigon-snode)
            deploy_chart "erigon-snode" "erigon-snode" "./mev-commit-erigon"
            
            # Deploy contracts if not dry-run and not skipping
            if [[ "$SKIP_CONTRACTS" == false && "$DRY_RUN" == false ]]; then
                print_info "Deploying contracts..."
                
                # Port forward
                ERIGON_POD=$(find_pod "erigon-snode")
                [[ -z "$ERIGON_POD" ]] && { print_error "Erigon pod not found"; exit 1; }
                
                kubectl port-forward pod/"$ERIGON_POD" "$RPC_PORT:$RPC_PORT" -n "$NAMESPACE" &
                PF_PID=$!
                sleep 5
                
                # Check RPC
                if ! check_rpc; then
                    print_error "RPC check failed"
                    exit 1
                fi
                
                # Deploy contracts
                cleanup_pf  # Stop our PF, Python script will create its own
                
                PYTHON_ARGS=("$DEPLOY_SCRIPT" "--password" "$PASSWORD")
                [[ -n "$WORK_DIR" ]] && PYTHON_ARGS+=("--work-dir" "$WORK_DIR")
                
                if python3 "${PYTHON_ARGS[@]}"; then
                    print_success "Contracts deployed!"
                    
                    # Extract addresses
                    CONTRACTS_JSON=$(find_contracts_json)
                    if [[ -n "$CONTRACTS_JSON" ]] && extract_contracts "$CONTRACTS_JSON"; then
                        print_info "Contract addresses ready for dashboard"
                    else
                        print_error "Failed to extract contract addresses"
                        exit 1
                    fi
                else
                    print_error "Contract deployment failed"
                    exit 1
                fi
            fi
            ;;
        relay-emulator)
            # Get L1 service
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            deploy_chart "mev-commit-emulator" "mev-commit-relay-emulator" "./mev-commit-relay-emulator" \
                --set "job.l1RpcUrl=$L1_URL"
            ;;
        oracle)
            # Oracle needs all services and contract addresses
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            ERIGON_SERVICE=$(find_service "erigon-snode")
            ERIGON_HTTP="http://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8545"
            ERIGON_WS="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            RELAY_SERVICE=$(find_service "mev-commit-emulator")
            RELAY_URL="http://${RELAY_SERVICE:-mev-commit-relay-emulator-mev-commit-emulator}.${NAMESPACE}.svc.cluster.local:8080"
            
            ORACLE_ARGS=(
                --set "network.l1RpcUrls[0]=$L1_URL"
                --set "network.settlementRpcHttp=$ERIGON_HTTP"
                --set "network.settlementRpcWs=$ERIGON_WS"
                --set "network.relayUrls[0]=$RELAY_URL"
            )
            
            # Add contract addresses if available
            if [[ -n "$BIDDER_REGISTRY" ]]; then
                ORACLE_ARGS+=(
                    --set "contracts.bidderRegistry=$BIDDER_REGISTRY"
                    --set "contracts.blockTracker=$BLOCK_TRACKER"
                    --set "contracts.oracle=$ORACLE"
                    --set "contracts.preconf=$PRECONF_MANAGER"
                    --set "contracts.providerRegistry=$PROVIDER_REGISTRY"
                )
            fi
            
            deploy_chart "erigon-oracle" "erigon-oracle" "./mev-commit-oracle" "${ORACLE_ARGS[@]}"
            ;;
        bootnode)
            # Bootnode needs contract addresses and service URLs
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            ERIGON_SERVICE=$(find_service "erigon-snode")
            ERIGON_HTTP="http://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8545"
            ERIGON_WS="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            BOOTNODE_ARGS=(
                --set "global.rpc.l1Endpoint=$L1_URL"
                --set "global.rpc.settlementEndpoint=$ERIGON_HTTP"
                --set "global.rpc.settlementWsEndpoint=$ERIGON_WS"
            )
            
            # Add contract addresses if available
            if [[ -n "$BIDDER_REGISTRY" ]]; then
                BOOTNODE_ARGS+=(
                    --set "global.contracts.bidderRegistry=$BIDDER_REGISTRY"
                    --set "global.contracts.blockTracker=$BLOCK_TRACKER"
                    --set "global.contracts.preconfStore=$PRECONF_MANAGER"
                    --set "global.contracts.providerRegistry=$PROVIDER_REGISTRY"
                )
            fi
            
            deploy_chart "erigon-mev-commit-bootnode" "erigon-mev-commit-bootnode" "./mev-commit-p2p" "${BOOTNODE_ARGS[@]}" -f "./mev-commit-p2p/bootnode-values.yaml"
            
            # Get bootnode connection string for other P2P nodes
            if [[ "$DRY_RUN" == false ]]; then
                get_bootnode_connection
            fi
            ;;
        bidder)
            # Bidder needs bootnode connection string and RPC endpoints
            if [[ -z "$BOOTNODE_CONNECTION" && "$DRY_RUN" == false ]]; then
                print_error "Bootnode connection string not available for bidder"
                exit 1
            fi
            
            # Get service URLs
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            ERIGON_SERVICE=$(find_service "erigon-snode")
            ERIGON_HTTP="http://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8545"
            ERIGON_WS="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            BIDDER_ARGS=(
                --set "global.rpc.l1Endpoint=$L1_URL"
                --set "global.rpc.settlementEndpoint=$ERIGON_HTTP"
                --set "global.rpc.settlementWsEndpoint=$ERIGON_WS"
            )
            
            [[ -n "$BOOTNODE_CONNECTION" ]] && BIDDER_ARGS+=(--set "node.bootnodeConnectionString=$BOOTNODE_CONNECTION")
            
            # Add contract addresses if available
            if [[ -n "$BIDDER_REGISTRY" ]]; then
                BIDDER_ARGS+=(
                    --set "global.contracts.bidderRegistry=$BIDDER_REGISTRY"
                    --set "global.contracts.blockTracker=$BLOCK_TRACKER"
                    --set "global.contracts.preconfStore=$PRECONF_MANAGER"
                    --set "global.contracts.providerRegistry=$PROVIDER_REGISTRY"
                )
            fi
            
            deploy_chart "erigon-mev-commit-bidder" "erigon-mev-commit-bidder" "./mev-commit-p2p" "${BIDDER_ARGS[@]}" -f "./mev-commit-p2p/values-bidder.yaml"
            ;;
        bidder-emulator)
            # Bidder emulator needs Mock L1 and Bidder service URLs
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            # For bidder service, we need to handle the P2P service naming
            BIDDER_SERVICE="erigon-mev-commit-bidder-mev-commit-p2p-bidder"
            BIDDER_RPC_URL="${BIDDER_SERVICE}.${NAMESPACE}.svc.cluster.local:13724"
            
            BIDDER_EMU_ARGS=(
                --set "bidderEmulator.l1RpcUrl=$L1_URL"
                --set "bidderEmulator.bidderRpcUrl=$BIDDER_RPC_URL"
            )
            
            deploy_chart "mev-commit-emulator-bt" "mev-commit-emulator-bt" "./mev-commit-emulator-bt" "${BIDDER_EMU_ARGS[@]}"
            ;;
        provider)
            # Provider needs bootnode connection string, RPC endpoints, and contract addresses
            if [[ -z "$BOOTNODE_CONNECTION" && "$DRY_RUN" == false ]]; then
                print_error "Bootnode connection string not available for provider"
                exit 1
            fi
            
            # Get service URLs
            L1_SERVICE=$(find_service "mev-commit-mock-l1")
            L1_URL="http://${L1_SERVICE:-erigon-mev-commit-mock-l1}.${NAMESPACE}.svc.cluster.local:8545"
            
            ERIGON_SERVICE=$(find_service "erigon-snode")
            ERIGON_HTTP="http://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8545"
            ERIGON_WS="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            PROVIDER_ARGS=(
                --set "global.rpc.l1Endpoint=$L1_URL"
                --set "global.rpc.settlementEndpoint=$ERIGON_HTTP"
                --set "global.rpc.settlementWsEndpoint=$ERIGON_WS"
            )
            
            [[ -n "$BOOTNODE_CONNECTION" ]] && PROVIDER_ARGS+=(--set "node.bootnodeConnectionString=$BOOTNODE_CONNECTION")
            
            # Add contract addresses if available
            if [[ -n "$BIDDER_REGISTRY" ]]; then
                PROVIDER_ARGS+=(
                    --set "global.contracts.bidderRegistry=$BIDDER_REGISTRY"
                    --set "global.contracts.blockTracker=$BLOCK_TRACKER"
                    --set "global.contracts.preconfStore=$PRECONF_MANAGER"
                    --set "global.contracts.providerRegistry=$PROVIDER_REGISTRY"
                )
            fi
            
            deploy_chart "erigon-mev-commit-provider" "erigon-mev-commit-provider" "./mev-commit-p2p" "${PROVIDER_ARGS[@]}" -f "./mev-commit-p2p/values-provider.yaml"
            ;;
        provider-emulator)
            # Provider emulator needs Provider and Relay service URLs
            PROVIDER_SERVICE="erigon-mev-commit-provider-mev-commit-p2p-provider"
            PROVIDER_URL="${PROVIDER_SERVICE}.${NAMESPACE}.svc.cluster.local:13624"
            
            RELAY_SERVICE=$(find_service "mev-commit-emulator")
            RELAY_URL="http://${RELAY_SERVICE:-mev-commit-relay-emulator-mev-commit-emulator}.${NAMESPACE}.svc.cluster.local:8080"
            
            PROVIDER_EMU_ARGS=(
                --set "job.serverAddr=$PROVIDER_URL"
                --set "job.relay=$RELAY_URL"
            )
            
            deploy_chart "mev-commit-emulator" "mev-commit-emulator" "./mev-commit-emulator" "${PROVIDER_EMU_ARGS[@]}" -f "./mev-commit-emulator/provider-emulator-values.yaml"
            ;;
        dashboard)
            # Get erigon service and contract addresses
            ERIGON_SERVICE=$(find_service "erigon-snode")
            RPC_URL="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            DASHBOARD_ARGS=(
                --set "config.rpcUrl=$RPC_URL"
            )
            
            # Add contract addresses if available
            if [[ -n "$BIDDER_REGISTRY" ]]; then
                DASHBOARD_ARGS+=(
                    --set "config.bidderregistryContractAddr=$BIDDER_REGISTRY"
                    --set "config.blocktrackerContractAddr=$BLOCK_TRACKER"
                    --set "config.oracleContractAddr=$ORACLE"
                    --set "config.preconfContractAddr=$PRECONF_MANAGER"
                    --set "config.providerregistryContractAddr=$PROVIDER_REGISTRY"
                )
            fi
            
            deploy_chart "mev-commit-dashboard" "mev-commit-dashboard" "./mev-commit-dashboard" "${DASHBOARD_ARGS[@]}"
            ;;
    esac
done

[[ "$DRY_RUN" == true ]] && { print_success "Dry-run completed!"; exit 0; }

print_success "All deployments completed!"

# Summary
print_info "=== Summary ==="
for CHART in "${CHARTS[@]}"; do
    case "$CHART" in
        mock-l1) print_info "Mock L1: $(find_service "mev-commit-mock-l1").$NAMESPACE.svc.cluster.local:8545" ;;
        erigon-snode) print_info "Erigon: $(find_service "erigon-snode").$NAMESPACE.svc.cluster.local:8545" ;;
        relay-emulator) print_info "Relay: $(find_service "mev-commit-relay-emulator").$NAMESPACE.svc.cluster.local:8080" ;;
        dashboard) print_info "Dashboard: $(find_service "mev-commit-dashboard").$NAMESPACE.svc.cluster.local:8080" ;;
    esac
done
[[ -n "$BIDDER_REGISTRY" ]] && print_info "Contracts deployed ✅"
