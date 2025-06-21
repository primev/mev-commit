#!/usr/bin/env python3
import json
import argparse
import requests
from pathlib import Path

# Default filename for the migrated genesis output
DEFAULT_OUT_NAME = "out_genesis.json"

def rpc_call(rpc_url: str, method: str, params: list) -> dict:
    payload = {
        "jsonrpc": "2.0",
        "id": 1,
        "method": method,
        "params": params
    }
    resp = requests.post(rpc_url, json=payload)
    resp.raise_for_status()
    data = resp.json()
    if "error" in data:
        raise RuntimeError(f"RPC error ({method}): {data['error']}")
    return data["result"]

def rpc_get_block_number(rpc_url: str) -> int:
    """Fetch the latest block number (as an int)."""
    hex_bn = rpc_call(rpc_url, "eth_blockNumber", [])
    return int(hex_bn, 16)

def rpc_debug_dump_block(rpc_url: str, block_param: str) -> dict:
    """Call debug_dumpBlock at the given block tag or hex number."""
    return rpc_call(rpc_url, "debug_dumpBlock", [block_param])

def to_hex(x: str) -> str:
    """Convert a decimal string to a hex string (0x-prefixed)."""
    return hex(int(x))

def build_alloc(accounts: dict) -> dict:
    """
    Given the 'accounts' map from debug_dumpBlock, return an 'alloc'
    dictionary suitable for a genesis file: address → { balance, code?, storage? }.
    """
    alloc = {}
    for addr, acct in accounts.items():
        entry = {"balance": to_hex(acct["balance"])}
        if acct.get("code"):
            entry["code"] = acct["code"]
        if acct.get("storage"):
            entry["storage"] = acct["storage"]
        alloc[addr] = entry
    return alloc

def main():
    p = argparse.ArgumentParser(
        description="Snapshot world-state via debug_dumpBlock and merge into a genesis template"
    )
    p.add_argument(
        "--rpc",
        default="http://127.0.0.1:8545",
        help="Geth RPC endpoint"
    )
    p.add_argument(
        "--input-genesis",
        required=True,
        help="Path to your source chain genesis.json"
    )
    p.add_argument(
        "--block", "-b",
        type=int,
        help="Block number to snapshot (defaults to latest)"
    )
    p.add_argument(
        "--output", "-o",
        help=(
            "Path or directory for output genesis JSON "
            f"(default: ./{DEFAULT_OUT_NAME})"
        )
    )
    args = p.parse_args()

    # load template
    tmpl_path = Path(args.input_genesis)
    if not tmpl_path.is_file():
        raise SystemExit(f"❌ Input genesis not found: {tmpl_path}")
    genesis_tpl = json.loads(tmpl_path.read_text())

    # decide which block to dump
    if args.block is None:
        # fetch actual latest block number
        block_no = rpc_get_block_number(args.rpc)
        block_param = "latest"
        print(f"⛓  Dumping world-state at latest (block {block_no})…")
    else:
        block_no = args.block
        block_param = hex(block_no)
        print(f"⛓  Dumping world-state at block {block_no}…")

    # fetch the dump
    dump = rpc_debug_dump_block(args.rpc, block_param)
    print(f"  ↳ {len(dump['accounts'])} accounts loaded")

    # build alloc → merge → write out
    latest_alloc = build_alloc(dump["accounts"])
    new_gen = genesis_tpl.copy()
    orig_alloc = new_gen.get("alloc", {})
    new_gen["alloc"] = {**orig_alloc, **latest_alloc}

    # determine output path
    if args.output:
        out_path = Path(args.output)
        if out_path.is_dir():
            out_path = out_path / DEFAULT_OUT_NAME
    else:
        out_path = Path.cwd() / DEFAULT_OUT_NAME

    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(json.dumps(new_gen, indent=2))
    print(f"✅  Wrote merged genesis → {out_path}")

if __name__ == "__main__":
    main()
