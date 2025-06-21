# Geth World State Snapshot

This repository provides a utility to snapshot the Ethereum world state at a given block and merge it into a new genesis file. The main script, `world_state.py`, connects to an archive node via JSON-RPC, retrieves account and storage data, and produces an updated genesis alloc.

---

## Prerequisites

- Python 3.8+
- An Ethereum archive node exposing the JSON-RPC interface (e.g. Geth with `--gcmode=archive`).

---

## Repository Structure

```bash
mev-commit/tools/state-dump
├── genesis_test.json
├── initial_genesis.json
├── requirements.txt
└── world_state.py
```

---

## Setup & Usage

```bash
# 1. Create and activate a virtual environment
python3 -m venv geth-world-state
source geth-world-state/bin/activate

# 2. Upgrade pip and install dependencies
pip install --upgrade pip
pip install -r requirements.txt

# 3. Run the snapshot script and merge into out_genesis.json
python3 world_state.py \
  --rpc http://<ARCHIVE_NODE_IP>:8545 \
  --input-genesis initial_genesis.json \
  --output out_genesis.json
```

**Example:**

```bash
python3 world_state.py \
  --rpc http://34.75.194.46:8545 \
  --input-genesis initial_genesis.json \
  --output out_genesis.json
```

- `--rpc`  
  JSON-RPC endpoint of your archive node (e.g. `http://127.0.0.1:8545`).

- `--input-genesis`  
  Path to your existing genesis template (e.g. `initial_genesis.json`).

- `--output`  
  Path where the merged genesis file will be written (e.g. `out_genesis.json`).
