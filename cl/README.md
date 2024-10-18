# README.md

## Introduction

This project sets up a local Ethereum environment with two Geth nodes, a Redis instance, and a consensus client. The setup is useful for testing and development purposes, allowing you to simulate a blockchain network locally.

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

Redis is configured in `redis-cluster` folder withing `docker-compose.yml`

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

## Additional Notes

- **Multiple Instances**: You can run multiple instances of the consensus client by changing the `--instance-id` and `--eth-client-url` parameters.

## Conclusion

You now have a local Ethereum environment with Geth nodes, Redis, and a consensus client.
