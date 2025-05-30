# README.md

## Introduction

This project sets up a local Ethereum environment with two Geth nodes, a Redis instance, and both a consensus client and a single node application (snode). The setup is useful for testing and development purposes, allowing you to simulate a blockchain network locally with different consensus options.

## Prerequisites

- **Go Programming Language**: Ensure you have Go installed on your system.
- **Geth (Go Ethereum)**: Install Geth from the [official website](https://geth.ethereum.org/downloads/).
- **Docker and Docker Compose**: Install Docker and Docker Compose from the [official website](https://docs.docker.com/get-docker/).
- **Redis**: We will use Redis via Docker Compose.
- **Git**: To clone the repository if needed.
- **Foundry (Optional)**: For interacting with the blockchain using `cast`.

## Setup Directory Structure

Ensure you have a directory named `geth-setup` in `cl` containing `genesis.json` and `jwt.hex` files:

```bash
geth-setup/
├── genesis.json
└── jwt.hex
```

## Installation

### Set Up JWT Secret

If you don't already have a `jwt.hex` file, create one:

```bash
echo "your_jwt_secret" > geth-setup/jwt.hex
```

**Note**: Replace `"your_jwt_secret"` with your actual JWT secret. In a production environment, ensure this secret is kept secure.

## Running Geth Nodes

We will initialize and run two Geth nodes.

### Initialize Geth Nodes

First, remove any existing data directories for clean setups:

```bash
rm -rf client1 && rm -rf client2
```

Initialize the nodes with the `genesis.json` file located in `geth-setup`:

```bash
geth init --datadir ./client1 geth-setup/genesis.json
geth init --datadir ./client2 geth-setup/genesis.json
```

### Run Geth Nodes

#### Node 1

```bash
geth --verbosity 5 \
  --datadir ./client1 \
  --nodiscover \
  --http \
  --http.port 8545 \
  --port 30303 \
  --authrpc.jwtsecret ./geth-setup/jwt.hex \
  --authrpc.port 8551 \
  --networkid 141414 \
  --http.api "admin,eth,net,web3,engine" \
  --syncmode full \
  --miner.recommit 900ms
```

#### Node 2

```bash
geth --verbosity 5 \
  --datadir ./client2 \
  --nodiscover \
  --http \
  --http.port 8546 \
  --port 30304 \
  --authrpc.jwtsecret ./geth-setup/jwt.hex \
  --authrpc.port 8552 \
  --networkid 141414 \
  --http.api "admin,eth,net,web3,engine" \
  --syncmode full \
  --miner.recommit 900ms
```

### Explanation of Geth Flags

- `--verbosity 5`: Sets the logging verbosity.
- `--datadir`: Specifies the data directory for the node.
- `--nodiscover`: Disables the peer discovery mechanism.
- `--http`: Enables the HTTP-RPC server.
- `--http.port`: Port for the HTTP-RPC server.
- `--port`: Network listening port.
- `--authrpc.jwtsecret`: Path to the JWT secret file.
- `--authrpc.port`: Port for authenticated RPC endpoints.
- `--networkid`: Network identifier for the Ethereum network.
- `--http.api`: APIs exposed over the HTTP-RPC interface.
- `--syncmode full`: Synchronization mode.
- `--miner.recommit`: Frequency of miner recommit.

### Obtaining the Genesis Block Hash

You can obtain the genesis block hash by querying the latest block after initializing your node:

```bash
cast block latest -r http://localhost:8545
```

Look for the `hash` field in the output, which represents the latest block hash. Since the chain is just initialized, this will be the genesis block hash.

Alternatively, you can use `curl` to get the genesis block hash:

```bash
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x0", false],"id":1}' -H "Content-Type: application/json" http://localhost:8545
```

Extract the `hash` value from the response and use it without `0x`.

## Running Redis

We will use Docker Compose to run Redis.

### Docker Compose Configuration

Redis is configured in `redis-cluster` folder within `docker-compose.yml`

### Start Redis

Stop any existing containers and remove volumes:

```bash
docker compose down -v
```

Start Redis in detached mode:

```bash
docker compose up -d
```

### Verify Redis is Running

You can verify that Redis is running by connecting to it:

```bash
redis-cli -p 7001
```

## Running the Consensus Client

The consensus client is a Go application that interacts with the Geth nodes and Redis.

### Build the Consensus Client

Ensure all dependencies are installed and build the application:

```bash
go mod tidy
go build -o consensus-client main.go
```

### Configuration

The consensus client can be configured via command-line flags, environment variables, or a YAML configuration file.

#### Command-Line Flags

- `--instance-id`: **(Required)** Unique instance ID for this node.
- `--eth-client-url`: Ethereum client URL (default: `http://localhost:8551`).
- `--jwt-secret`: JWT secret for Ethereum client.
- `--genesis-block-hash`: Genesis block hash.
- `--redis-addr`: Redis address (default: `127.0.0.1:7001`).
- `--evm-build-delay`: EVM build delay (default: `1s`).
- `--config`: Path to a YAML configuration file.

#### Environment Variables

- `RAPP_INSTANCE_ID`
- `RAPP_ETH_CLIENT_URL`
- `RAPP_JWT_SECRET`
- `RAPP_GENESIS_BLOCK_HASH`
- `RAPP_REDIS_ADDR`
- `RAPP_EVM_BUILD_DELAY`
- `RAPP_CONFIG`

### Run the Consensus Client

Run the client using command-line flags:

```bash
./consensus-client start \
  --instance-id "node1" \
  --eth-client-url "http://localhost:8551" \
  --jwt-secret "your_jwt_secret" \
  --genesis-block-hash "your_genesis_block_hash" \
  --redis-addr "127.0.0.1:7001" \
  --evm-build-delay "1s"
```

**Note**:

- Replace `"your_jwt_secret"` with the actual JWT secret you used earlier.
- Replace `"your_genesis_block_hash"` with the genesis block hash obtained earlier.

### Using a Configuration File

Create a `config.yaml` file:

```yaml
instance-id: "node1"
eth-client-url: "http://localhost:8551"
jwt-secret: "your_jwt_secret"
genesis-block-hash: "your_genesis_block_hash"
redis-addr: "127.0.0.1:7001"
evm-build-delay: "1s"
```

Run the client with the configuration file:

```bash
./consensus-client start --config config.yaml
```

## Running the Single Node Application (snode)

The single node application provides a simplified MEV-commit setup that doesn't require Redis, but using Postgres to save payloads, so member nodes could request that payload later on.

## Architecture Overview

The application supports two operational modes:

1. **Leader Node**: Produces blocks and serves payloads to member nodes via API
2. **Member Node**: Follows a leader node by polling for and processing payloads sequentially

## Running Postgres

We will use Docker Compose to run Postgres

### Postgres Docker Compose Configuration

Postgres is configured in `postgres` folder within `docker-compose.yml`

### Start Postgres

Stop any existing containers and remove volumes:

```bash
docker compose down -v
```

Start Postgres in detached mode:

```bash
docker compose up -d
```

### Build the Single Node Application

```bash
go mod tidy
go build -o snode main.go
```

### SNode Configuration

The snode application can be configured via command-line flags, environment variables, or a YAML configuration file.

### Common Configuration Flags

- `--instance-id`: **(Required)** Unique instance ID for this node
- `--eth-client-url`: Ethereum Execution client Engine API URL (default: `http://localhost:8551`)
- `--jwt-secret`: Hex-encoded JWT secret for Ethereum Execution client Engine API (default: `13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c`)
- `--health-addr`: Address for health check endpoint (default: `:8080`)
- `--config`: Path to a YAML configuration file
- `--log-fmt`: Log format ('text' or 'json') (default: `text`)
- `--log-level`: Log level ('debug', 'info', 'warn', 'error') (default: `info`)
- `--log-tags`: Comma-separated <name:value> log tags (e.g., `env:prod,service:snode`)

### Leader Node Specific Flags

- `--priority-fee-recipient`: **(Required)** Ethereum address for receiving priority fees (block proposer fee)
- `--evm-build-delay`: Delay after initiating payload construction before calling getPayload (default: `100ms`)
- `--evm-build-delay-empty-block`: Minimum time since last block to build an empty block (default: `2s`, 0 to disable skipping)
- `--postgres-dsn`: PostgreSQL DSN for storing payloads (optional, e.g., `postgres://user:pass@host:port/dbname?sslmode=disable`)
- `--api-addr`: Address for member node API endpoint (default: `:9090`, empty to disable)

### Member Node Specific Flags

- `--leader-api-url`: **(Required)** Leader node API URL for member nodes (e.g., `http://leader:9090`)
- `--poll-interval`: Interval for polling leader node for new payloads (default: `1s`)

### SNode Environment Variables

**Common:**

- `SNODE_INSTANCE_ID`
- `SNODE_ETH_CLIENT_URL`
- `SNODE_JWT_SECRET`
- `SNODE_HEALTH_ADDR`
- `SNODE_CONFIG`
- `MEV_COMMIT_LOG_FMT`
- `MEV_COMMIT_LOG_LEVEL`
- `MEV_COMMIT_LOG_TAGS`

**Leader Node:**

- `SNODE_PRIORITY_FEE_RECIPIENT`
- `SNODE_EVM_BUILD_DELAY`
- `SNODE_EVM_BUILD_DELAY_EMPTY_BLOCK`
- `SNODE_POSTGRES_DSN`
- `SNODE_API_ADDR`

**Member Node:**

- `MEMBER_LEADER_API_URL`
- `MEMBER_POLL_INTERVAL`

## Running the Application

### Leader Node

Run as a leader node (produces blocks and serves API for member nodes):

```bash
./snode leader \
  --instance-id "leader1" \
  --eth-client-url "http://localhost:8551" \
  --jwt-secret "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c" \
  --priority-fee-recipient "0xYourEthereumAddress" \
  --evm-build-delay "100ms" \
  --evm-build-delay-empty-block "2s" \
  --api-addr ":9090" \
  --log-level "info"
```

### Member Node

Run as a member node (follows leader by polling for payloads):

```bash
./snode member \
  --instance-id "member1" \
  --eth-client-url "http://localhost:8552" \
  --jwt-secret "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c" \
  --leader-api-url "http://localhost:9090" \
  --poll-interval "1s" \
  --log-level "info"
```

### Backward Compatibility

The legacy `start` command is still supported and runs as a leader node:

```bash
./snode start \
  --instance-id "snode1" \
  --priority-fee-recipient "0xYourEthereumAddress" \
  # ... other flags
```

**Note**:

- Replace `"0xYourEthereumAddress"` with a valid Ethereum address for receiving priority fees.
- The JWT secret should be a 64-character hex string (32 bytes).

## Configuration Files

### Leader Node Configuration

Create a `leader-config.yaml` file:

```yaml
instance-id: "leader1"
eth-client-url: "http://localhost:8551"
jwt-secret: "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c"
priority-fee-recipient: "0xYourEthereumAddress"
evm-build-delay: "100ms"
evm-build-delay-empty-block: "2s"
api-addr: ":9090"
postgres-dsn: "postgres://user:pass@localhost:5432/mevcommit?sslmode=disable"
health-addr: ":8080"
log-fmt: "text"
log-level: "info"
log-tags: "env:dev,service:leader"
```

Run with configuration file:

```bash
./snode leader --config leader-config.yaml
```

### Member Node Configuration

Create a `member-config.yaml` file:

```yaml
instance-id: "member1"
eth-client-url: "http://localhost:8552"
jwt-secret: "13373d9a0257983ad150392d7ddb2f9172c9396b4c450e26af469d123c7aaa5c"
leader-api-url: "http://localhost:9090"
poll-interval: "1s"
health-addr: ":8081"
log-fmt: "text"
log-level: "info"
log-tags: "env:dev,service:member"
```

Run with configuration file:

```bash
./snode member --config member-config.yaml
```

### Health Endpoints

Both node types provide health check endpoints:

- **Leader**: Returns 200 OK when operational, 503 when Ethereum client unavailable
- **Member**: Returns 200 OK when operational and leader available, 503 otherwise

Access health endpoints at: `http://localhost:8080/health` (or configured port)

## Multi-Node Setup Example

For a complete leader-follower setup:

1. **Start Leader Node**:

   ```bash
   ./snode leader --instance-id "leader" --priority-fee-recipient "0xYourAddress" --api-addr ":9090"
   ```

2. **Start Member Node(s)**:

   ```bash
   ./snode member --instance-id "member1" --leader-api-url "http://localhost:9090" --eth-client-url "http://localhost:8552" --health-addr ":8081"
   ```

Each member node should connect to its own Geth instance and configure unique health endpoints to avoid port conflicts.

## Additional Notes

- **Graceful Shutdown**: Both applications support graceful shutdown via SIGTERM or Ctrl+C.

## Conclusion

You now have a local Ethereum environment with Geth nodes offering two consensus options: a Redis-based leader-follower consensus setup or a simplified single node consensus.
