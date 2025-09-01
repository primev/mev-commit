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
BRANCH=""
CONTRACTS_FILE=""
LOCAL=false

# Image variables
DASHBOARD_IMAGE=""
P2P_IMAGE=""
RELAY_IMAGE=""
BIDDER_EMULATOR_IMAGE=""
PROVIDER_EMULATOR_IMAGE=""
ORACLE_IMAGE=""

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
        --branch) BRANCH="$2"; shift 2 ;;
        --contracts-file) CONTRACTS_FILE="$2"; shift 2 ;;
        --local) LOCAL="$2"; shift 2 ;;
        --help) 
            echo "Usage: $0 [--dry-run] [--cleanup] [--skip-contracts] [--password PASS] [--charts CHARTS] [--branch BRANCH] [--contracts-file FILE] [--local true/false]"
            echo "Charts: mock-l1,erigon-snode,relay-emulator,dashboard,oracle,bootnode,bidder,bidder-emulator,provider,provider-emulator"
            echo "Branch: Filter Docker images by branch label (required for image discovery)"
            echo "Contracts File: JSON file with contract addresses (used with --skip-contracts)"
            echo "Local: true=use ClusterIP & IfNotPresent policy (minikube), false=use LoadBalancer & Always policy (cloud)"
            exit 0 ;;
        *) print_error "Unknown option: $1"; exit 1 ;;
    esac
done

# Cleanup function
cleanup_all() {
    print_info "Cleaning up all releases..."
    
    for release in "mev-commit-dashboard" "mev-commit-relay-emulator" "erigon-snode" "erigon-mev-commit-mock-l1" "erigon-oracle" "erigon-mev-commit-bootnode" "erigon-mev-commit-bidder" "mev-commit-emulator-bt" "erigon-mev-commit-provider" "mev-commit-emulator"; do
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

# Discover Docker images by branch and component
discover_images() {
    if [[ -z "$BRANCH" ]]; then
        print_warning "No branch specified, using hardcoded image values"
        return 0
    fi
    
    print_info "Discovering Docker images for branch: $BRANCH"
    
    # Get all images filtered by branch
    BRANCH_IMAGES=$(docker images --filter "label=branch=$BRANCH" --format "table {{.Repository}}:{{.Tag}}\t{{.ID}}" 2>/dev/null || true)
    
    if [[ -z "$BRANCH_IMAGES" ]]; then
        print_warning "No Docker images found with branch label: $BRANCH"
        print_warning "Using hardcoded image values from charts"
        return 0
    fi
    
    print_info "Found images for branch $BRANCH:"
    echo "$BRANCH_IMAGES"
    echo
    
    # Extract images by component label
    print_info "Matching images by component labels..."
    
    # Dashboard image (component=dashboard)
    DASHBOARD_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=dashboard" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    # P2P image (component=p2p) - used by bootnode, bidder, provider
    P2P_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=p2p" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    # Relay emulator image (component=relay-emulator)
    RELAY_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=relay-emulator" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    # Bidder emulator image (component=bidder-emulator)
    BIDDER_EMULATOR_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=bidder-emulator" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    # Provider emulator image (component=provider-emulator)
    PROVIDER_EMULATOR_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=provider-emulator" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    # Oracle image (component=oracle)
    ORACLE_IMAGE=$(docker images --filter "label=branch=$BRANCH" --filter "label=component=oracle" --format "{{.Repository}}:{{.Tag}}" 2>/dev/null | head -1 || true)
    
    print_info "=== Image Discovery Results ==="
    [[ -n "$DASHBOARD_IMAGE" ]] && print_success "Dashboard: $DASHBOARD_IMAGE" || print_warning "Dashboard: Using chart default"
    [[ -n "$P2P_IMAGE" ]] && print_success "P2P (bootnode/bidder/provider): $P2P_IMAGE" || print_warning "P2P: Using chart default"
    [[ -n "$RELAY_IMAGE" ]] && print_success "Relay Emulator: $RELAY_IMAGE" || print_warning "Relay Emulator: Using chart default"
    [[ -n "$BIDDER_EMULATOR_IMAGE" ]] && print_success "Bidder Emulator: $BIDDER_EMULATOR_IMAGE" || print_warning "Bidder Emulator: Using chart default"
    [[ -n "$PROVIDER_EMULATOR_IMAGE" ]] && print_success "Provider Emulator: $PROVIDER_EMULATOR_IMAGE" || print_warning "Provider Emulator: Using chart default"
    [[ -n "$ORACLE_IMAGE" ]] && print_success "Oracle: $ORACLE_IMAGE" || print_warning "Oracle: Using chart default"
    echo
}

# Parse discovered image into repository and tag
parse_image() {
    local IMAGE="$1"
    if [[ -n "$IMAGE" && "$IMAGE" == *":"* ]]; then
        REPO="${IMAGE%:*}"
        TAG="${IMAGE##*:}"
        echo "$REPO" "$TAG"
    else
        echo "" ""
    fi
}

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

# Get bootnode connection info - behavior depends on --local flag
get_bootnode_connection() {
    BOOTNODE_SERVICE="erigon-mev-commit-bootnode-mev-commit-p2p-bootnode"
    MAX_ATTEMPTS=20
    
    if [[ "$LOCAL" == "true" ]]; then
        print_info "Getting bootnode connection info for local deployment..."
        
        # Get ClusterIP for the connection string (internal cluster communication)
        BOOTNODE_CLUSTER_IP=$(kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" -o jsonpath='{.spec.clusterIP}' 2>/dev/null)
        
        if [[ -z "$BOOTNODE_CLUSTER_IP" || "$BOOTNODE_CLUSTER_IP" == "<nil>" ]]; then
            print_error "Failed to get bootnode ClusterIP"
            kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" || true
            return 1
        fi
        
        print_success "Bootnode ClusterIP: $BOOTNODE_CLUSTER_IP (for connection string)"
        
        # For local deployment, use 127.0.0.1 for API access (minikube exposes via localhost)
        print_info "Using 127.0.0.1 for API access (minikube/local setup)"
        API_HOST="127.0.0.1"
        API_URL="https://$API_HOST:13723/v1/debug/topology"
        
        # Connection string will use ClusterIP for internal cluster communication
        CONNECTION_IP="$BOOTNODE_CLUSTER_IP"
        
    else
        print_info "Getting bootnode connection info for cloud deployment..."
        
        # Wait for LoadBalancer to get external IP
        print_info "Waiting for LoadBalancer external IP..."
        for i in $(seq 1 30); do
            BOOTNODE_EXTERNAL_IP=$(kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null)
            if [[ -n "$BOOTNODE_EXTERNAL_IP" && "$BOOTNODE_EXTERNAL_IP" != "<nil>" ]]; then
                break
            fi
            print_info "Attempt $i/30: Waiting for LoadBalancer IP..."
            sleep 10
        done
        
        if [[ -z "$BOOTNODE_EXTERNAL_IP" || "$BOOTNODE_EXTERNAL_IP" == "<nil>" ]]; then
            print_error "Failed to get bootnode external IP after 300 seconds"
            kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" || true
            return 1
        fi
        
        print_success "Bootnode LoadBalancer IP: $BOOTNODE_EXTERNAL_IP"
        
        # For cloud deployment, use external IP for both API and connection string
        API_URL="https://$BOOTNODE_EXTERNAL_IP:13723/v1/debug/topology"
        CONNECTION_IP="$BOOTNODE_EXTERNAL_IP"
    fi
    
    # Get peer ID from bootnode API
    print_info "Getting peer ID from bootnode API..."
    for i in $(seq 1 $MAX_ATTEMPTS); do
        print_info "Attempt $i/$MAX_ATTEMPTS: Calling $API_URL"
        
        # Use HTTPS for API call
        PEER_ID=$(curl -sk "$API_URL" | jq -r '.topology.self.Underlay' 2>/dev/null)
        
        if [[ -n "$PEER_ID" && "$PEER_ID" != "null" ]]; then
            # Connection string uses the appropriate IP based on deployment type
            BOOTNODE_CONNECTION="/ip4/$CONNECTION_IP/tcp/13522/p2p/$PEER_ID"
            print_success "Bootnode connection: $BOOTNODE_CONNECTION"
            
            if [[ "$LOCAL" == "true" ]]; then
                print_info "Local deployment: API via 127.0.0.1, P2P connection via ClusterIP"
            else
                print_info "Cloud deployment: Both API and P2P connection via LoadBalancer IP"
            fi
            return 0
        fi
        
        print_warning "No valid peer ID received, retrying in 5 seconds..."
        sleep 5
    done
    
    print_error "Failed to get bootnode peer ID after $MAX_ATTEMPTS attempts"
    
    if [[ "$LOCAL" == "true" ]]; then
        print_info "For local deployment, make sure:"
        print_info "1. Bootnode service is accessible via https://127.0.0.1:13723"
        print_info "2. You may need to run: kubectl port-forward svc/$BOOTNODE_SERVICE 13723:13723 -n $NAMESPACE"
        print_info "3. Or use: minikube service $BOOTNODE_SERVICE --url -n $NAMESPACE"
    else
        print_info "For cloud deployment, make sure LoadBalancer service is properly configured"
    fi
    
    print_info "Checking bootnode service and pod status:"
    kubectl get svc "$BOOTNODE_SERVICE" -n "$NAMESPACE" || true
    kubectl get pods -l "app.kubernetes.io/name=mev-commit-p2p,app.kubernetes.io/component=bootnode" -n "$NAMESPACE" || true
    return 1
}

# Load contracts from file (when --skip-contracts is used)
load_contracts_from_file() {
    if [[ "$SKIP_CONTRACTS" == false ]]; then
        return 0  # Not skipping contracts, will deploy them
    fi
    
    print_info "Loading contracts from file (--skip-contracts enabled)..."
    
    # Determine contracts file to use
    CONTRACTS_JSON=""
    
    if [[ -n "$CONTRACTS_FILE" ]]; then
        # Use specified file
        if [[ -f "$CONTRACTS_FILE" ]]; then
            CONTRACTS_JSON="$CONTRACTS_FILE"
            print_info "Using specified contracts file: $CONTRACTS_FILE"
        else
            print_error "Specified contracts file not found: $CONTRACTS_FILE"
            return 1
        fi
    else
        # Auto-discover latest contracts file
        CONTRACTS_JSON=$(find_contracts_json)
        if [[ -n "$CONTRACTS_JSON" ]]; then
            print_info "Auto-discovered contracts file: $CONTRACTS_JSON"
        else
            print_error "No contracts file found. Provide one with --contracts-file or deploy contracts first"
            print_info "Expected JSON format:"
            cat << 'EOF'
{
  "BidderRegistry": "0x1234567890123456789012345678901234567890",
  "BlockTracker": "0x2345678901234567890123456789012345678901",
  "Oracle": "0x3456789012345678901234567890123456789012",
  "PreconfManager": "0x4567890123456789012345678901234567890123",
  "ProviderRegistry": "0x5678901234567890123456789012345678901234"
}
EOF
            return 1
        fi
    fi
    
    # Load contracts from file
    if extract_contracts "$CONTRACTS_JSON"; then
        print_success "Contracts loaded from file!"
        print_info "BidderRegistry: $BIDDER_REGISTRY"
        print_info "BlockTracker: $BLOCK_TRACKER" 
        print_info "Oracle: $ORACLE"
        print_info "PreconfManager: $PRECONF_MANAGER"
        print_info "ProviderRegistry: $PROVIDER_REGISTRY"
        return 0
    else
        return 1
    fi
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

# Check password requirement (only when deploying contracts)
if [[ "$SKIP_CONTRACTS" == false && "$DRY_RUN" == false && -z "$PASSWORD" ]]; then
    for CHART in "${CHARTS[@]}"; do
        [[ "$CHART" == "erigon-snode" ]] && { print_error "Password required for contracts"; exit 1; }
    done
fi

# Discover images before deployment
discover_images

# Load contracts from file if skipping deployment
load_contracts_from_file

print_info "Deploying charts: ${CHARTS[*]}"

# Initialize contract variables (will be set by load_contracts_from_file if skipping)
if [[ "$SKIP_CONTRACTS" == false ]]; then
    BIDDER_REGISTRY=""
    BLOCK_TRACKER=""
    ORACLE=""
    PRECONF_MANAGER=""
    PROVIDER_REGISTRY=""
fi

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
            
            # Prepare image args
            RELAY_ARGS=(--set "job.l1RpcUrl=$L1_URL")
            
            if [[ -n "$RELAY_IMAGE" ]]; then
                read RELAY_REPO RELAY_TAG <<< "$(parse_image "$RELAY_IMAGE")"
                if [[ -n "$RELAY_REPO" && -n "$RELAY_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    RELAY_ARGS+=(--set "image.repository=$RELAY_REPO" --set "image.tag=$RELAY_TAG" --set "image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered relay image: $RELAY_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
            deploy_chart "mev-commit-emulator" "mev-commit-relay-emulator" "./mev-commit-relay-emulator" "${RELAY_ARGS[@]}"
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
            
            # Add discovered oracle image
            if [[ -n "$ORACLE_IMAGE" ]]; then
                read ORACLE_REPO ORACLE_TAG <<< "$(parse_image "$ORACLE_IMAGE")"
                if [[ -n "$ORACLE_REPO" && -n "$ORACLE_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    ORACLE_ARGS+=(--set "image.repository=$ORACLE_REPO" --set "image.tag=$ORACLE_TAG" --set "image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered oracle image: $ORACLE_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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
            
            # Add discovered P2P image
            if [[ -n "$P2P_IMAGE" ]]; then
                read P2P_REPO P2P_TAG <<< "$(parse_image "$P2P_IMAGE")"
                if [[ -n "$P2P_REPO" && -n "$P2P_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    BOOTNODE_ARGS+=(--set "global.image.repository=$P2P_REPO" --set "global.image.tag=$P2P_TAG" --set "global.image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered P2P image: $P2P_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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
            
            # Add discovered P2P image
            if [[ -n "$P2P_IMAGE" ]]; then
                read P2P_REPO P2P_TAG <<< "$(parse_image "$P2P_IMAGE")"
                if [[ -n "$P2P_REPO" && -n "$P2P_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    BIDDER_ARGS+=(--set "global.image.repository=$P2P_REPO" --set "global.image.tag=$P2P_TAG" --set "global.image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered P2P image: $P2P_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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
            
            # Add discovered bidder emulator image
            if [[ -n "$BIDDER_EMULATOR_IMAGE" ]]; then
                read BIDDER_EMU_REPO BIDDER_EMU_TAG <<< "$(parse_image "$BIDDER_EMULATOR_IMAGE")"
                if [[ -n "$BIDDER_EMU_REPO" && -n "$BIDDER_EMU_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    BIDDER_EMU_ARGS+=(--set "image.repository=$BIDDER_EMU_REPO" --set "image.tag=$BIDDER_EMU_TAG" --set "image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered bidder emulator image: $BIDDER_EMULATOR_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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
            
            # Add discovered P2P image
            if [[ -n "$P2P_IMAGE" ]]; then
                read P2P_REPO P2P_TAG <<< "$(parse_image "$P2P_IMAGE")"
                if [[ -n "$P2P_REPO" && -n "$P2P_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    PROVIDER_ARGS+=(--set "global.image.repository=$P2P_REPO" --set "global.image.tag=$P2P_TAG" --set "global.image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered P2P image: $P2P_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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
            
            # Add discovered provider emulator image
            if [[ -n "$PROVIDER_EMULATOR_IMAGE" ]]; then
                read PROVIDER_EMU_REPO PROVIDER_EMU_TAG <<< "$(parse_image "$PROVIDER_EMULATOR_IMAGE")"
                if [[ -n "$PROVIDER_EMU_REPO" && -n "$PROVIDER_EMU_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    PROVIDER_EMU_ARGS+=(--set "image.repository=$PROVIDER_EMU_REPO" --set "image.tag=$PROVIDER_EMU_TAG" --set "image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered provider emulator image: $PROVIDER_EMULATOR_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
            deploy_chart "mev-commit-emulator" "mev-commit-emulator" "./mev-commit-emulator" "${PROVIDER_EMU_ARGS[@]}" -f "./mev-commit-emulator/provider-emulator-values.yaml"
            ;;
        dashboard)
            # Get erigon service and contract addresses
            ERIGON_SERVICE=$(find_service "erigon-snode")
            RPC_URL="ws://${ERIGON_SERVICE:-erigon-snode-erigon}.${NAMESPACE}.svc.cluster.local:8546"
            
            DASHBOARD_ARGS=(
                --set "config.rpcUrl=$RPC_URL"
            )
            
            # Add discovered dashboard image
            if [[ -n "$DASHBOARD_IMAGE" ]]; then
                read DASHBOARD_REPO DASHBOARD_TAG <<< "$(parse_image "$DASHBOARD_IMAGE")"
                if [[ -n "$DASHBOARD_REPO" && -n "$DASHBOARD_TAG" ]]; then
                    PULL_POLICY="Always"
                    [[ "$LOCAL" == "true" ]] && PULL_POLICY="IfNotPresent"
                    DASHBOARD_ARGS+=(--set "image.repository=$DASHBOARD_REPO" --set "image.tag=$DASHBOARD_TAG" --set "image.pullPolicy=$PULL_POLICY")
                    print_info "Using discovered dashboard image: $DASHBOARD_IMAGE (pullPolicy: $PULL_POLICY)"
                fi
            fi
            
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

# Image summary
if [[ -n "$BRANCH" ]]; then
    print_info "=== Images Used ==="
    [[ -n "$DASHBOARD_IMAGE" ]] && print_info "Dashboard: $DASHBOARD_IMAGE" || print_info "Dashboard: Chart default"
    [[ -n "$P2P_IMAGE" ]] && print_info "P2P (bootnode/bidder/provider): $P2P_IMAGE" || print_info "P2P: Chart default"
    [[ -n "$RELAY_IMAGE" ]] && print_info "Relay Emulator: $RELAY_IMAGE" || print_info "Relay Emulator: Chart default"
    [[ -n "$BIDDER_EMULATOR_IMAGE" ]] && print_info "Bidder Emulator: $BIDDER_EMULATOR_IMAGE" || print_info "Bidder Emulator: Chart default"
    [[ -n "$PROVIDER_EMULATOR_IMAGE" ]] && print_info "Provider Emulator: $PROVIDER_EMULATOR_IMAGE" || print_info "Provider Emulator: Chart default"
    [[ -n "$ORACLE_IMAGE" ]] && print_info "Oracle: $ORACLE_IMAGE" || print_info "Oracle: Chart default"
fi
