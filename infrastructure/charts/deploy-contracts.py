#!/usr/bin/env python3
# deploy_contracts.py
#
# Flow:
#   1. download keystores, verify, print Oracle & Sender addresses
#   2. git clone mev-commit
#   3. forge install (at repo root)
#   4. forge clean / build (in contracts/)
#   5. port-forward & deploy contracts, parse logs → JSONs

import os, re, sys, json, subprocess, signal, time, socket, argparse
from datetime import datetime
from pathlib import Path
from tempfile import NamedTemporaryFile

# ---------- Defaults ----------
REPO_URL = "https://github.com/primev/mev-commit.git"
REPO_BRANCH = "main"

ORACLE_KEYSTORE_URL = "https://storage.googleapis.com/devnet-artifacts/keystores/erigon-keystores/keystore2/UTC--2025-06-24T18-20-44.554617000Z--1c533735c11dd317bc816629f86e00f479d097a3"
SENDER_KEYSTORE_URL = "https://storage.googleapis.com/devnet-artifacts/keystores/erigon-keystores/keystore1/UTC--2025-06-24T18-20-43.890647000Z--421657a89f467ac04c542e46645a7752e199b5e6"

KUBE_POD = "erigon-snode-0"
RPC_HOST = "127.0.0.1"
RPC_PORT = 8545
RPC_URL = f"http://{RPC_HOST}:{RPC_PORT}"
KUBECTL_BIN = "kubectl"

CHAIN_ID = "141414"
PRIORITY_GAS_PRICE = "2000000000"
GAS_PRICE = "5000000000"
SOLC = "0.8.26"

CORE_KEYS = [
    "BidderRegistry", "BlockTracker", "Oracle",
    "PreconfManager", "ProviderRegistry", "ValidatorOptInRouter"
]

# ---------- CLI ----------
ap = argparse.ArgumentParser()
ap.add_argument("--password", "-p", required=True, help="Password for keystores")
ap.add_argument("--work-dir", default=f"/tmp/mev-commit-{datetime.now().strftime('%Y%m%d_%H%M%S')}")
args = ap.parse_args()

WORK_DIR = Path(args.work_dir).resolve()
WORK_DIR.mkdir(parents=True, exist_ok=True)

# ---------- Helpers ----------
def run(cmd, cwd=None, log_file=None):
    if log_file:
        with open(log_file, "ab") as lf:
            lf.write(("\n$ " + " ".join(cmd) + "\n").encode())
            return subprocess.run(cmd, cwd=cwd, stdout=lf, stderr=subprocess.STDOUT).returncode
    return subprocess.run(cmd, cwd=cwd).returncode

def wait_for_port(host, port, timeout=25):
    deadline = time.time() + timeout
    while time.time() < deadline:
        try:
            with socket.create_connection((host, port), timeout=1):
                return True
        except Exception:
            time.sleep(0.5)
    return False

def download_keystore(url: str, dest_dir: Path) -> Path:
    fname = url.split("/")[-1]
    fpath = dest_dir / fname
    if not fpath.exists():
        subprocess.check_call(["wget", "-q", "-O", str(fpath), url])
    return fpath

def read_keystore_address(fpath: Path) -> str:
    data = json.loads(fpath.read_text())
    return "0x" + data["address"]

# ---------- Step 1. Keystores ----------
ks_dir = WORK_DIR / "keystores"
ks_dir.mkdir(parents=True, exist_ok=True)

oracle_keystore = download_keystore(ORACLE_KEYSTORE_URL, ks_dir)
oracle_addr = read_keystore_address(oracle_keystore)
os.environ["ORACLE_KEYSTORE_ADDRESS"] = oracle_addr

sender_keystore = download_keystore(SENDER_KEYSTORE_URL, ks_dir)
sender_addr = read_keystore_address(sender_keystore)

print(f"🔑 Oracle address: {oracle_addr}")
print(f"👤 Sender address: {sender_addr}")

# ---------- Step 2. Clone repo ----------
repo_dir = WORK_DIR / "mev-commit"
print(f"🔹 Cloning {REPO_URL}@{REPO_BRANCH} → {repo_dir}")
subprocess.check_call(["git", "clone", "--branch", REPO_BRANCH, "--single-branch", REPO_URL, str(repo_dir)])
contracts_dir = repo_dir / "contracts"

log_dir = contracts_dir / "deploy-logs"
out_dir = contracts_dir / "deploy-summaries"
for d in (log_dir, out_dir): d.mkdir(parents=True, exist_ok=True)

ts = datetime.now().strftime("%Y%m%d_%H%M%S")
log_path = log_dir / f"deploy_core_{ts}.log"
full_json_path = out_dir / f"contracts_{ts}.json"
core_json_path = out_dir / "core-contracts.json"

# ---------- Step 3. Forge install (repo root) ----------
print("🔹 forge install …")
for pkg in ("OpenZeppelin/openzeppelin-contracts-upgradeable", "OpenZeppelin/openzeppelin-contracts"):
    if run(["forge", "install", pkg], cwd=repo_dir, log_file=log_path) != 0:
        sys.exit(1)

# ---------- Step 4. Clean/build (contracts/) ----------
print("🔹 forge clean && forge build …")
if run(["forge", "clean"], cwd=contracts_dir, log_file=log_path) != 0: sys.exit(1)
if run(["forge", "build"], cwd=contracts_dir, log_file=log_path) != 0: sys.exit(1)

# ---------- Step 5. Port-forward + Deploy ----------
print(f"🔹 Port-forward {KUBE_POD} → localhost:{RPC_PORT}")
pf = subprocess.Popen([KUBECTL_BIN, "port-forward", f"pod/{KUBE_POD}", f"{RPC_PORT}:{RPC_PORT}"],
                      preexec_fn=os.setsid)
import atexit; atexit.register(lambda: os.killpg(os.getpgid(pf.pid), signal.SIGTERM))
if not wait_for_port(RPC_HOST, RPC_PORT): sys.exit("RPC not up")

print("🔹 Deploying …")
deploy_cmd = [
    "forge", "script", "scripts/core/DeployCore.s.sol:DeployCore",
    "--priority-gas-price", PRIORITY_GAS_PRICE,
    "--with-gas-price", GAS_PRICE,
    "--chain-id", CHAIN_ID,
    "--rpc-url", RPC_URL,
    "--keystores", str(sender_keystore),
    "--password", args.password,
    "--sender", sender_addr,
    "--use", SOLC,
    "--broadcast", "--json", "--via-ir",
]
rc = run(deploy_cmd, cwd=contracts_dir, log_file=log_path)

# ---------- Parse logs ----------
text = log_path.read_text(errors="ignore")
pairs = {}
for name in CORE_KEYS:
    m = re.search(re.escape(name)+r"[^0-9a-fA-F]{0,200}(0x[0-9a-fA-F]{40})", text, re.I|re.S|re.M)
    if m: pairs[name] = m.group(1)
for m in re.finditer(r"([A-Za-z][A-Za-z0-9:_ ]{1,64})[^0-9a-fA-F]{0,40}(0x[0-9a-fA-F]{40})", text):
    name = re.sub(r"[:\s]+$", "", m.group(1).strip())
    if len(name) >= 2: pairs.setdefault(name, m.group(2))

with full_json_path.open("w") as jf: json.dump(pairs, jf, indent=2, sort_keys=True)
core_map = {k: pairs[k] for k in CORE_KEYS if k in pairs}
with NamedTemporaryFile("w", delete=False, dir=str(core_json_path.parent), prefix=".tmp_", suffix=".json") as tmp:
    json.dump(core_map, tmp, indent=2); tmp.flush(); os.fsync(tmp.fileno())
    os.replace(tmp.name, core_json_path)

# ---------- Result ----------
if rc == 0: print("✅ Success.")
else: sys.exit(f"❌ Deploy failed (exit {rc}). See log: {log_path}")

print(f"📄 Log: {log_path}")
print(f"🧾 Full JSON: {full_json_path}")
print(f"🎯 Core JSON: {core_json_path}")

