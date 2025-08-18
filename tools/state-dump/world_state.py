#!/usr/bin/env python3
"""
World-state → genesis.alloc via Erigon-style debug_accountRange (6-arg, 'next'),
pinned to a single block number to avoid 'missing trie node ... loc: diff' errors.

- Prints and pins the exact block number (hex) for ALL calls.
- Lightweight account paging (no inline code/storage).
- Optional per-contract code and full storage paging.
- Adaptive page sizes for both accounts and storage.
"""

import argparse
import json
import time
from pathlib import Path
from typing import Dict, Any, Optional

import requests

DEFAULT_OUT = "out_genesis.json"


class RpcError(RuntimeError):
    pass


def rpc_call(url: str, method: str, params: list, *, timeout: float, retries: int = 0, backoff: float = 1.6) -> Any:
    last = None
    for i in range(retries + 1):
        try:
            r = requests.post(url, json={"jsonrpc": "2.0", "id": 1, "method": method, "params": params}, timeout=timeout)
            r.raise_for_status()
            data = r.json()
            if "error" in data:
                raise RpcError(f"{method}: {data['error']}")
            return data["result"]
        except (requests.Timeout, requests.ConnectionError, RpcError) as e:
            last = e
            if i < retries:
                time.sleep(backoff ** i)
            else:
                raise
    raise last  # pragma: no cover


def to_hex_any(x) -> str:
    if x is None:
        return "0x0"
    if isinstance(x, int):
        return hex(x)
    if isinstance(x, str):
        s = x.strip()
        if s.startswith("0x") or s.startswith("0X"):
            return s
        return hex(int(s))
    return hex(int(x))


def get_latest_block(url: str, timeout: float) -> int:
    h = rpc_call(url, "eth_blockNumber", [], timeout=timeout)
    return int(h, 16)


def account_range(url: str, block_hex: str, start_token: str, n: int, *, timeout: float) -> Dict[str, Any]:
    # Erigon 6-arg: nocode=True, nostorage=True, incompletes=True (lightweight)
    return rpc_call(url, "debug_accountRange", [block_hex, start_token, n, True, True, True], timeout=timeout)


def storage_range_at(url: str, block_hex: str, addr: str, start_key: str, n: int, *, timeout: float) -> Dict[str, Any]:
    # Pin to the same block number (hex). txIndex=None → post-state.
    return rpc_call(url, "debug_storageRangeAt", [block_hex, None, addr, start_key, n], timeout=timeout)


def eth_get_code(url: str, addr: str, block_hex: str, *, timeout: float) -> str:
    # Use the same pinned block number (hex) here as well.
    return rpc_call(url, "eth_getCode", [addr, block_hex], timeout=timeout)


def page_full_storage(
    url: str,
    block_hex: str,
    addr: str,
    *,
    timeout: float,
    initial_page_size: int = 2048,
    min_page_size: int = 128,
) -> Dict[str, str]:
    storage: Dict[str, str] = {}
    sk = "0x"
    page = initial_page_size
    while True:
        try:
            sres = storage_range_at(url, block_hex, addr, sk, page, timeout=timeout)
        except RpcError as e:
            msg = str(e).lower()
            # Adaptive shrink on timeouts
            if "timed out" in msg and page > min_page_size:
                page = max(min_page_size, page // 2)
                print(f"⚠️  storageRangeAt timed out for {addr}; reducing storage page-size to {page} and retrying…")
                continue
            # Most important: missing trie node usually comes from unpinned or pruned state.
            if "missing trie node" in msg:
                raise
            raise

        s_map = sres.get("storage") or {}
        for k, v in s_map.items():
            val = v.get("value", "0x0")
            if val not in (None, "0x", "0x0", "0", 0):
                storage[k] = val
        nxt = sres.get("nextKey")
        if not nxt:
            break
        sk = nxt
    return storage


def build_alloc(
    url: str,
    block_hex: str,
    *,
    initial_page_size: int,
    min_page_size: int,
    include_nonces: bool,
    include_code: bool,
    include_storage: bool,
    timeout: float,
    storage_page_initial: int,
    storage_page_min: int,
) -> dict:
    alloc: dict = {}
    start = "0x"
    total = 0
    page_size = initial_page_size

    while True:
        # Adaptive account page
        while True:
            try:
                res = account_range(url, block_hex, start, page_size, timeout=timeout)
                break
            except RpcError as e:
                if "timed out" in str(e).lower() and page_size > min_page_size:
                    page_size = max(min_page_size, page_size // 2)
                    print(f"⚠️  accountRange timed out; reducing account page-size to {page_size} and retrying…")
                    continue
                raise

        accounts = res.get("accounts") or {}

        for addr, meta in accounts.items():
            entry = {"balance": to_hex_any(meta.get("balance", "0"))}
            if include_nonces and (meta.get("nonce") is not None):
                entry["nonce"] = to_hex_any(meta["nonce"])

            # Code / storage (per-contract)
            if include_code or include_storage:
                code = eth_get_code(url, addr, block_hex, timeout=timeout)
                if code and code != "0x":
                    if include_code:
                        entry["code"] = code
                    if include_storage:
                        try:
                            storage = page_full_storage(
                                url,
                                block_hex,
                                addr,
                                timeout=timeout,
                                initial_page_size=storage_page_initial,
                                min_page_size=storage_page_min,
                            )
                            if storage:
                                entry["storage"] = storage
                        except RpcError as e:
                            # If still missing trie nodes at a fixed block, node likely pruned / missing state.
                            if "missing trie node" in str(e).lower():
                                print(f"❌ storageRangeAt failed for {addr} due to missing trie node at pinned block. "
                                      f"Your node likely lacks full state for storage paging (needs archive-like state).")
                                print("   → Options: (1) re-run without --include-storage, "
                                      "or (2) use an archive/fully-synced node with state available at this block.")
                                # Continue with other accounts; keep code/balance/nonce
                            else:
                                raise

            alloc[addr] = entry
            total += 1

        nxt = res.get("next")
        if not nxt:
            break
        start = nxt

    print(f"  ↳ paged: {total} accounts")
    return alloc


def main():
    ap = argparse.ArgumentParser(description="World-state → genesis.alloc (Erigon accountRange) pinned to a block, with adaptive paging.")
    ap.add_argument("--rpc", required=True)
    ap.add_argument("--input-genesis", required=True)
    ap.add_argument("--output", "-o", help=f"Output path or directory (default ./{DEFAULT_OUT})")
    ap.add_argument("--block", "-b", type=int, help="Block number (default: latest)")
    ap.add_argument("--exclude-nonces", action="store_true")
    ap.add_argument("--page-size", type=int, default=2048)
    ap.add_argument("--min-page-size", type=int, default=256)
    ap.add_argument("--include-code", action="store_true")
    ap.add_argument("--include-storage", action="store_true")  # implies include-code if contract present
    ap.add_argument("--rpc-timeout", type=float, default=300.0)
    ap.add_argument("--storage-page-size", type=int, default=2048, help="Initial storage page size")
    ap.add_argument("--storage-min-page-size", type=int, default=128, help="Minimum storage page size")
    args = ap.parse_args()

    include_nonces = not args.exclude_nonces
    include_code = bool(args.include_code or args.include_storage)
    include_storage = bool(args.include_storage)

    # Load template
    tpl_path = Path(args.input_genesis)
    if not tpl_path.is_file():
        raise SystemExit(f"❌ Input genesis not found: {tpl_path}")
    genesis_tpl = json.loads(tpl_path.read_text())

    # Resolve and PIN the block to a constant hex number
    if args.block is None:
        bn = get_latest_block(args.rpc, args.rpc_timeout)
    else:
        bn = args.block
    block_hex = hex(bn)

    print(f"⛓  Scanning state pinned at block {bn} (tag {block_hex}, account page {args.page_size})…")
    if include_storage:
        print("   • Including contract storage (per-account paging). This may take a while.")

    # Build alloc
    alloc = build_alloc(
        args.rpc,
        block_hex,
        initial_page_size=args.page_size,
        min_page_size=args.min_page_size,
        include_nonces=include_nonces,
        include_code=include_code,
        include_storage=include_storage,
        timeout=args.rpc_timeout,
        storage_page_initial=args.storage_page_size,
        storage_page_min=args.storage_min_page_size,
    )

    # Merge & write
    new_gen = dict(genesis_tpl)
    base_alloc = dict(new_gen.get("alloc") or {})
    base_alloc.update(alloc)
    new_gen["alloc"] = base_alloc

    out_path = Path(args.output) if args.output else Path.cwd() / DEFAULT_OUT
    if out_path.is_dir():
        out_path = out_path / DEFAULT_OUT
    out_path.parent.mkdir(parents=True, exist_ok=True)
    out_path.write_text(json.dumps(new_gen, indent=2))
    print(f"✅  Wrote merged genesis → {out_path}")


if __name__ == "__main__":
    main()
