#!/usr/bin/env python3
import argparse
import json
import sys
from web3 import Web3
from web3.middleware import ExtraDataToPOAMiddleware
from web3.exceptions import Web3RPCError


def parse_quantity(q: str) -> int:
    """
    Convert a hex or decimal string to an integer.
    """
    return int(q, 16) if isinstance(q, str) and q.startswith("0x") else int(q)


def load_alloc(path: str) -> dict:
    """
    Load the `alloc` section from a genesis JSON file.
    """
    with open(path, 'r') as f:
        data = json.load(f)
    return data.get("alloc") or data.get("Alloc") or {}


# Addresses to exclude (lowercase, no '0x' prefix)
EXCLUDE_ADDRS = {
    '00000000219ab540356cbb839cbe05303d7705fa',
    '00000961ef480eb55e80d19ad83579a64c007002',
    '0000bbddc7ce488642fb579f8b00f3a590007251',
    '0000f90827f1c53a10cb7a02335b175320002935',
    'e666c8fc7be06a4a88145cd0c626fae4d40fcda7',
    '3ba6a3318a7d55c73f743529e2ca69ccf112d538',
    'badc9848e1e87e5017e7790be7c4b5d35c304fc1',
    'c286bef43cea547545d5b7179aef6747f63ac8aa',
    '072f6d4d7a1f7af547d47d927beaf38e01fcb33b'
}


def normalize(addr: str) -> str:
    """
    Normalize address to lowercase without '0x' prefix for comparison.
    """
    a = addr.lower()
    return a[2:] if a.startswith('0x') else a


def check_nonces(w3_source: Web3, w3_target: Web3, src_id, tgt_id, alloc: dict) -> None:
    print("🔎 Checking nonces for all addresses")
    for addr, entry in alloc.items():
        if normalize(addr) in EXCLUDE_ADDRS:
            continue
        cs = Web3.to_checksum_address(addr)
        gen_nonce = parse_quantity(entry.get("nonce", "0x0"))
        try:
            n_src = w3_source.eth.get_transaction_count(cs, block_identifier=src_id)
        except Web3RPCError:
            n_src = gen_nonce
        try:
            n_tgt = w3_target.eth.get_transaction_count(cs, block_identifier=tgt_id)
        except Web3RPCError:
            n_tgt = gen_nonce
        status = "✅" if (n_src == n_tgt == gen_nonce) else "⛔"
        print(f"Address {addr}: {status} nonce src={n_src} tgt={n_tgt} genesis={gen_nonce}")
    print()


def check_balances(w3_source: Web3, w3_target: Web3, src_id, tgt_id, alloc: dict) -> None:
    print("🔎 Checking balances for all addresses")
    for addr, entry in alloc.items():
        if normalize(addr) in EXCLUDE_ADDRS:
            continue
        cs = Web3.to_checksum_address(addr)
        gen_bal = parse_quantity(entry.get("balance", hex(0)))
        try:
            b_src = w3_source.eth.get_balance(cs, block_identifier=src_id)
        except Web3RPCError:
            b_src = gen_bal
        try:
            b_tgt = w3_target.eth.get_balance(cs, block_identifier=tgt_id)
        except Web3RPCError:
            b_tgt = gen_bal
        status = "✅" if (b_src == b_tgt == gen_bal) else "⛔"
        print(f"Address {addr}: {status} balance src={b_src} tgt={b_tgt} genesis={gen_bal}")
    print()


def check_code(w3_source: Web3, w3_target: Web3, src_id, tgt_id, alloc: dict) -> None:
    print("🔎 Checking code for all contracts (genesis vs source vs target)")
    for addr, entry in alloc.items():
        if normalize(addr) in EXCLUDE_ADDRS:
            continue
        gen_code = entry.get("code", "").lower()
        if gen_code in ("", "0x"):
            continue
        cs = Web3.to_checksum_address(addr)
        try:
            code_src = "0x" + w3_source.eth.get_code(cs, block_identifier=src_id).hex()
        except Web3RPCError:
            code_src = None
        try:
            code_tgt = "0x" + w3_target.eth.get_code(cs, block_identifier=tgt_id).hex()
        except Web3RPCError:
            code_tgt = None
        if code_src == code_tgt == gen_code:
            print(f"Address {addr}: ✅ code matches (genesis=source=target)")
        else:
            print(f"Address {addr}: ⛔ code mismatch")
            print(f"  genesis: {gen_code}")
            print(f"  source:  {code_src}")
            print(f"  target:  {code_tgt}")
    print()


def check_storage_roots(w3_source: Web3, w3_target: Web3, src_id, tgt_id, alloc: dict) -> None:
    print("🔎 Checking storage‐roots for all contracts")
    for addr, entry in alloc.items():
        if normalize(addr) in EXCLUDE_ADDRS:
            continue
        gen_code = entry.get("code", "").lower()
        if gen_code in ("", "0x"):
            continue

        cs = Web3.to_checksum_address(addr)

        try:
            proof_src = w3_source.eth.get_proof(cs, [], block_identifier=src_id)
            hex_src = Web3.to_hex(proof_src["storageHash"])
        except Web3RPCError:
            hex_src = None

        try:
            proof_tgt = w3_target.eth.get_proof(cs, [], block_identifier=tgt_id)
            hex_tgt = Web3.to_hex(proof_tgt["storageHash"])
        except Web3RPCError:
            hex_tgt = None

        if hex_src is not None and hex_src == hex_tgt:
            print(f"Address {addr}: ✅ storage‐root matches (value={hex_src})")
        else:
            print(f"Address {addr}: ⛔ storage‐root src={hex_src} tgt={hex_tgt}")
    print()


def main() -> None:
    parser = argparse.ArgumentParser(description="Verify EVM state migration using a genesis alloc file")
    parser.add_argument("--genesis", required=True)
    parser.add_argument("--source-rpc-url", required=True)
    parser.add_argument(
        "--source-block-number",
        type=str,
        default="latest",
        help="source block identifier (hex, dec) or 'latest' (default)"
    )
    parser.add_argument("--target-rpc-url", required=True)
    parser.add_argument(
        "--target-block-number",
        type=str,
        default="latest",
        help="target block identifier (hex, dec) or 'latest' (default)"
    )
    args = parser.parse_args()

    alloc = load_alloc(args.genesis)
    if not alloc:
        print("No alloc found.")
        sys.exit(1)

    # Determine source identifier
    src_id = args.source_block_number.lower()
    if src_id != 'latest':
        try:
            src_id = parse_quantity(src_id)
        except ValueError:
            print(f"Invalid source block: {args.source_block_number}")
            sys.exit(1)

    # Determine target identifier
    tgt_id = args.target_block_number.lower()
    if tgt_id != 'latest':
        try:
            tgt_id = parse_quantity(tgt_id)
        except ValueError:
            print(f"Invalid target block: {args.target_block_number}")
            sys.exit(1)

    w3_src = Web3(Web3.HTTPProvider(args.source_rpc_url))
    w3_src.middleware_onion.inject(ExtraDataToPOAMiddleware, layer=0)
    w3_tgt = Web3(Web3.HTTPProvider(args.target_rpc_url))
    w3_tgt.middleware_onion.inject(ExtraDataToPOAMiddleware, layer=0)

    check_nonces(w3_src, w3_tgt, src_id, tgt_id, alloc)
    check_balances(w3_src, w3_tgt, src_id, tgt_id, alloc)
    check_code(w3_src, w3_tgt, src_id, tgt_id, alloc)
    check_storage_roots(w3_src, w3_tgt, src_id, tgt_id, alloc)


if __name__ == "__main__":
    main()
