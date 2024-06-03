import json
import argparse
import time
from web3 import Web3
from eth_account import Account

def check_geth_running(w3):
    while True:
        try:
            latest_block = w3.eth.get_block('latest')
            if latest_block:
                print(f"Geth is running. Latest block number: {latest_block['number']}")
                return True
        except Exception as e:
            print(f"Error connecting to Geth: {e}. Retrying in 5 seconds...")
            time.sleep(5)

def deploy_create2(w3, account):
    # Check if contract already deployed
    address = "0x4e59b44847b379578588920ca78fbf26c0b4956c"
    checksum_address = Web3.to_checksum_address(address)
    code = w3.eth.get_code(checksum_address)
    if code != b'':
        print(f"Contract already deployed at {checksum_address}")
        return
    else:
        print(f"No contract deployed at {checksum_address}. Deploying...")

    # Presigned transaction
    transaction = "0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222"

    # Deploy contract
    tx_hash = w3.eth.send_raw_transaction(transaction)
    print(f"CREATE2 deployment transaction hash: {tx_hash.hex()}")

    # Wait for the transaction receipt
    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    print(f"CREATE2 contract deployed at address: {checksum_address}")

    # Check for leftover balance
    address_balance = "0x3fab184622dc19b6109349b94811493bf2a45362"
    checksum_balance_address = Web3.to_checksum_address(address_balance)
    balance = w3.eth.get_balance(checksum_balance_address)
    if balance != 0:
        print(f"WARNING: Deployment signer ({checksum_balance_address}) has leftover balance of {balance} wei")

def deploy_contract_create2(w3, account, contract_file, out_dir, gas_limit, chain_id):
    contract_name = contract_file.split('.')[0]
    json_file = f"{out_dir}/{contract_file}/{contract_name}.json"

    # Load the contract ABI and bytecode
    with open(json_file) as f:
        contract_json = json.load(f)
        abi = contract_json["abi"]
        bytecode = contract_json["bytecode"]["object"]

    # Ensure the bytecode is a valid hexadecimal string
    if not bytecode.startswith("0x"):
        bytecode = "0x" + bytecode

    # Create2 address generation
    salt = Web3.keccak(text=contract_name)  # Use a consistent salt
    create2_address = Web3.to_checksum_address(Web3.solidity_keccak(
        ['bytes1', 'address', 'bytes32', 'bytes32'],
        ['0xff', account.address, salt, Web3.keccak(hexstr=bytecode)]
    )[12:])

    # Check if the contract is already deployed
    if w3.eth.get_code(create2_address) != b'':
        print(f"Contract {contract_name} already deployed at {create2_address}")
        return

    # Prepare CREATE2 deployment via proxy contract
    proxy_contract_address = Web3.to_checksum_address("0x4e59b44847b379578588920ca78fbf26c0b4956c")
    proxy_contract = w3.eth.contract(address=proxy_contract_address, abi=[{
        "constant": False,
        "inputs": [
            {"name": "_salt", "type": "bytes32"},
            {"name": "_code", "type": "bytes"}
        ],
        "name": "deploy",
        "outputs": [{"name": "addr", "type": "address"}],
        "payable": False,
        "stateMutability": "nonpayable",
        "type": "function"
    }])

    # Build the proxy deployment transaction
    construct_txn = proxy_contract.functions.deploy(salt, bytecode).build_transaction({
        'from': account.address,
        'nonce': w3.eth.get_transaction_count(account.address),
        'gas': gas_limit,
        'gasPrice': w3.eth.gas_price,
        'chainId': chain_id
    })

    # Sign the transaction with the private key
    signed_txn = w3.eth.account.sign_transaction(construct_txn, private_key=account.key)

    # Send the transaction and wait for the receipt
    tx_hash = w3.eth.send_raw_transaction(signed_txn.rawTransaction)
    print(f"Transaction hash for {contract_name}: {tx_hash.hex()}")

    tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
    print(f"Contract {contract_name} deployed at address: {create2_address}")

def load_account_from_keystore(keystore_file, password):
    with open(keystore_file) as keyfile:
        encrypted_key = keyfile.read()
        private_key = Account.decrypt(encrypted_key, password)
        return Account.from_key(private_key)

# Command-line argument parsing
parser = argparse.ArgumentParser(description="Deploy Ethereum contracts.")
parser.add_argument("--keystore", help="The keystore file for the Ethereum account.")
parser.add_argument("--password", help="The password for the keystore file.")
parser.add_argument("--private-key", help="The private key for the Ethereum account.")
parser.add_argument("--rpc-url", default="http://127.0.0.1:8545", help="The RPC URL of the Ethereum node.")
parser.add_argument("--contracts", nargs='+', default=["BidderRegistry.sol", "PreConfCommitmentStore.sol", "BlockTracker.sol", "Oracle.sol", "ProviderRegistry.sol"], help="List of contract files to deploy.")
parser.add_argument("--gas-limit", type=int, default=3000000, help="The gas limit for the transactions.")
parser.add_argument("--chain-id", type=int, default=31337, help="The chain ID of the Ethereum network.")
parser.add_argument("--out-dir", default="out", help="The directory containing the compiled contract outputs.")
parser.add_argument("--skip-deploy-create2", action='store_true', help="Skip deploying the CREATE2 proxy contract.")

args = parser.parse_args()

# Configuration
KEYSTORE = args.keystore
PASSWORD = args.password
PRIVATE_KEY = args.private_key
RPC_URL = args.rpc_url
CONTRACTS = args.contracts
GAS_LIMIT = args.gas_limit
CHAIN_ID = args.chain_id
OUT_DIR = args.out_dir
SKIP_DEPLOY_CREATE2 = args.skip_deploy_create2

# Validate input arguments
if (KEYSTORE and PRIVATE_KEY) or (not KEYSTORE and not PRIVATE_KEY):
    print("You must provide either a keystore file and password, or a private key, but not both.")
    exit(1)

# Connect to the Ethereum node
w3 = Web3(Web3.HTTPProvider(RPC_URL))
if not w3.is_connected():
    print("Failed to connect to the Ethereum node.")
    exit(1)

# Check if Geth is running
check_geth_running(w3)

# Get the account from the keystore file and password or from the private key
if KEYSTORE:
    account = load_account_from_keystore(KEYSTORE, PASSWORD)
else:
    account = w3.eth.account.from_key(PRIVATE_KEY)

# Deploy CREATE2 proxy if not skipped
if not SKIP_DEPLOY_CREATE2:
    deploy_create2(w3, account)

# Loop through each contract and deploy it using CREATE2
for contract in CONTRACTS:
    deploy_contract_create2(w3, account, contract, OUT_DIR, GAS_LIMIT, CHAIN_ID)

print("All contracts deployed!")