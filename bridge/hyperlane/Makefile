up-dev:
	@if [ ! -f .env ]; then echo "Error: .env file not found. Please populate the .env file before running this command."; exit 1; fi
	AGENT_BASE_IMAGE=shaspitz/hyperlane-agent-mac:6985064 SETTLEMENT_RPC_URL=http://172.29.0.98:8545 PUBLIC_SETTLEMENT_RPC_URL=nil docker compose --profile bridge up -d --build

up-prod:
	@if [ ! -f .env ]; then echo "Error: .env file not found. Please populate the .env file before running this command."; exit 1; fi
	AGENT_BASE_IMAGE=gcr.io/abacus-labs-dev/hyperlane-agent@sha256:854f92966eac6b49e5132e152cc58168ecdddc76c2d390e657b81bdaf1396af0 SETTLEMENT_RPC_URL=http://172.29.0.98:8545 PUBLIC_SETTLEMENT_RPC_URL=nil docker compose --profile bridge up -d --build

down:
	AGENT_BASE_IMAGE=nil SETTLEMENT_RPC_URL=nil PUBLIC_SETTLEMENT_RPC_URL=nil docker compose --profile bridge down

up-emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=http://sl-bootnode:8545 docker compose --profile 1-emulators up -d --build

up-2emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=http://sl-bootnode:8545 docker compose --profile 2-emulators up -d --build

up-3emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=http://sl-bootnode:8545 docker compose --profile 3-emulators up -d --build

up-4emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=http://sl-bootnode:8545 docker compose --profile 4-emulators up -d --build

up-5emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=http://sl-bootnode:8545 docker compose --profile 5-emulators up -d --build

down-emulator:
	AGENT_BASE_IMAGE=nil PUBLIC_SETTLEMENT_RPC_URL=nil SETTLEMENT_RPC_URL=nil docker compose --profile 5-emulators down

clean-dbs:
	@read -p "WARNING: This command will wipe all persistent disk data relevant to the containers. Press enter to continue or Ctrl+C to cancel." _ 
	-docker compose --profile bridge down --rmi all --volumes
	-docker compose --profile bridge rm -fv
	docker image prune -f

print-warp-deploy:
	@echo "Printing hyperlane warp deployment artifact"
	docker run --rm -v hyperlane_hyperlane-deploy-artifacts:/data alpine cat /data/warp-deployment.json