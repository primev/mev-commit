import json
import argparse
import time
from web3 import Web3
from web3.middleware import geth_poa_middleware
from eth_account import Account

def check_geth_running(w3):
    w3.middleware_onion.inject(geth_poa_middleware, layer=0)
    while True:
        try:
            latest_block = w3.eth.get_block('latest')
            if latest_block:
                print(f"Geth is running. Latest block number: {latest_block['number']}")
                return True
        except Exception as e:
            print(f"Error connecting to Geth: {e}. Retrying in 5 seconds...")
            time.sleep(5)

def print_revert_reason(w3, tx_receipt):
    try:
        tx = w3.eth.get_transaction(tx_receipt['transactionHash'])
        tx_input = tx['input']
        to_address = tx_receipt['contractAddress'] if tx_receipt['contractAddress'] else tx_receipt['to']
        try:
            w3.eth.call({
                'to': to_address,
                'from': tx['from'],
                'data': tx_input
            }, tx_receipt['blockNumber'])
        except Exception as e:
            print(f"Revert reason: {e}")
    except Exception as e:
        print(f"Failed to fetch revert reason: {e}")

def deploy_create2(w3, account, chain_id):
    address = "0x4e59b44847b379578588920ca78fbf26c0b4956c"
    checksum_address = Web3.to_checksum_address(address)
    code = w3.eth.get_code(checksum_address)
    if code != b'':
        print(f"Contract already deployed at {checksum_address}")
        return
    else:
        print(f"No contract deployed at {checksum_address}. Deploying...")

    transaction = "0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222"

    try:
        tx = {
            'to': checksum_address,
            'from': account.address,
            'data': transaction,
            'gasPrice': w3.eth.gas_price,
            'chainId': chain_id
        }
        
        try:
            estimated_gas = w3.eth.estimate_gas(tx)
            tx['gas'] = estimated_gas
        except Exception as e:
            print(f"Simulation failed for CREATE2 deployment: {e}")
            return

        tx_hash = w3.eth.send_raw_transaction(transaction)
        print(f"CREATE2 deployment transaction hash: {tx_hash.hex()}")
        tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
        if tx_receipt['status'] == 0:
            print(f"CREATE2 contract deployment failed.")
            print_revert_reason(w3, tx_receipt)
        else:
            print(f"CREATE2 contract deployed at address: {checksum_address}")
    except Exception as e:
        print(f"Failed to deploy CREATE2 contract: {e}")

    address_balance = "0x3fab184622dc19b6109349b94811493bf2a45362"
    checksum_balance_address = Web3.to_checksum_address(address_balance)
    balance = w3.eth.get_balance(checksum_balance_address)
    if balance != 0:
        print(f"WARNING: Deployment signer ({checksum_balance_address}) has leftover balance of {balance} wei")

def deploy_contract(w3, account, contract_file, out_dir, gas_limit, chain_id, owner_addresses, registry_contract=None, registry_abi=None):
    contract_name = contract_file.split('.')[0]
    json_file = f"{out_dir}/{contract_file}/{contract_name}.json"

    with open(json_file) as f:
        contract_json = json.load(f)
        abi = contract_json["abi"]
        bytecode = contract_json["bytecode"]["object"]

    if not bytecode.startswith("0x"):
        bytecode = "0x" + bytecode

    contract = w3.eth.contract(abi=abi, bytecode=bytecode)
    construct_txn = contract.constructor().build_transaction({
        'from': account.address,
        'nonce': w3.eth.get_transaction_count(account.address),
        'gasPrice': w3.eth.gas_price,
        'chainId': chain_id
    })

    try:
        estimated_gas = w3.eth.estimate_gas(construct_txn)
        construct_txn['gas'] = estimated_gas
    except Exception as e:
        print(f"Gas estimation failed for contract {contract_name}: {e}")
        return

    signed_txn = w3.eth.account.sign_transaction(construct_txn, private_key=account.key)

    try:
        tx_hash = w3.eth.send_raw_transaction(signed_txn.rawTransaction)
        print(f"Transaction hash for {contract_name}: {tx_hash.hex()}")
        tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
        if tx_receipt['status'] == 0:
            print(f"Contract deployment failed for {contract_name}.")
            print_revert_reason(w3, tx_receipt)
        else:
            contract_address = tx_receipt.contractAddress
            print(f"Contract {contract_name} successfully deployed at address: {contract_address}")
            if contract_name in owner_addresses:
                transfer_ownership(w3, account, abi, contract_address, owner_addresses[contract_name], chain_id)
            if registry_contract:
                register_contract(w3, account, registry_contract, contract_name, contract_address, registry_abi, chain_id)
    except Exception as e:
        print(f"Failed to deploy contract {contract_name}: {e}")

def transfer_ownership(w3, account, abi, contract_address, new_owner, chain_id):
    contract = w3.eth.contract(address=contract_address, abi=abi)
    transfer_txn = contract.functions.transferOwnership(new_owner).build_transaction({
        'from': account.address,
        'nonce': w3.eth.get_transaction_count(account.address),
        'gasPrice': w3.eth.gas_price,
        'chainId': chain_id
    })

    try:
        estimated_gas = w3.eth.estimate_gas(transfer_txn)
        transfer_txn['gas'] = estimated_gas
    except Exception as e:
        print(f"Gas estimation failed for transfer ownership: {e}")
        return

    signed_txn = w3.eth.account.sign_transaction(transfer_txn, private_key=account.key)

    try:
        tx_hash = w3.eth.send_raw_transaction(signed_txn.rawTransaction)
        print(f"Ownership transfer transaction hash: {tx_hash.hex()}")
        tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
        if tx_receipt['status'] == 0:
            print(f"Ownership transfer failed for contract at address: {contract_address}")
            print_revert_reason(w3, tx_receipt)
        else:
            print(f"Ownership successfully transferred to {new_owner} for contract at address: {contract_address}")
    except Exception as e:
        print(f"Failed to transfer ownership for contract at address: {contract_address}: {e}")

def register_contract(w3, account, registry_contract, contract_name, contract_address, registry_abi, chain_id):
    contract = w3.eth.contract(address=registry_contract, abi=registry_abi)
    register_txn = contract.functions.addContract(contract_name, contract_address).build_transaction({
        'from': account.address,
        'nonce': w3.eth.get_transaction_count(account.address),
        'gasPrice': w3.eth.gas_price,
        'chainId': chain_id
    })

    try:
        estimated_gas = w3.eth.estimate_gas(register_txn)
        register_txn['gas'] = estimated_gas
    except Exception as e:
        print(f"Gas estimation failed for register contract: {e}")
        return

    signed_txn = w3.eth.account.sign_transaction(register_txn, private_key=account.key)

    try:
        tx_hash = w3.eth.send_raw_transaction(signed_txn.rawTransaction)
        print(f"Register contract transaction hash for {contract_name}: {tx_hash.hex()}")
        tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
        if tx_receipt['status'] == 0:
            print(f"Registering contract failed for {contract_name}")
            print_revert_reason(w3, tx_receipt)
        else:
            print(f"Contract {contract_name} successfully registered in registry")
    except Exception as e:
        print(f"Failed to register contract {contract_name}: {e}")

def load_account_from_keystore(keystore_file, password):
    with open(keystore_file) as keyfile:
        encrypted_key = keyfile.read()
        private_key = Account.decrypt(encrypted_key, password)
        return Account.from_key(private_key)

def parse_owner_addresses(owner_addresses_str):
    owner_addresses = {}
    if owner_addresses_str:
        pairs = owner_addresses_str.split(',')
        for pair in pairs:
            contract, address = pair.split(':')
            owner_addresses[contract] = address
    return owner_addresses

def deploy_registry_contract(w3, account, out_dir, gas_limit, chain_id):
    registry_name = "ContractRegistry"
    json_file = f"{out_dir}/{registry_name}.sol/{registry_name}.json"

    with open(json_file) as f:
        contract_json = json.load(f)
        abi = contract_json["abi"]
        bytecode = contract_json["bytecode"]["object"]

    if not bytecode.startswith("0x"):
        bytecode = "0x" + bytecode

    contract = w3.eth.contract(abi=abi, bytecode=bytecode)
    construct_txn = contract.constructor().build_transaction({
        'from': account.address,
        'nonce': w3.eth.get_transaction_count(account.address),
        'gasPrice': w3.eth.gas_price,
        'chainId': chain_id
    })

    try:
        estimated_gas = w3.eth.estimate_gas(construct_txn)
        construct_txn['gas'] = estimated_gas
    except Exception as e:
        print(f"Gas estimation failed for registry contract: {e}")
        return None, None

    signed_txn = w3.eth.account.sign_transaction(construct_txn, private_key=account.key)

    try:
        tx_hash = w3.eth.send_raw_transaction(signed_txn.rawTransaction)
        print(f"Transaction hash for {registry_name}: {tx_hash.hex()}")
        tx_receipt = w3.eth.wait_for_transaction_receipt(tx_hash)
        if tx_receipt['status'] == 0:
            print(f"Registry contract deployment failed.")
            print_revert_reason(w3, tx_receipt)
            return None, None
        else:
            contract_address = tx_receipt.contractAddress
            print(f"Registry contract successfully deployed at address: {contract_address}")
            return contract_address, abi
    except Exception as e:
        print(f"Failed to deploy registry contract: {e}")
        return None, None

parser = argparse.ArgumentParser(description="Deploy Ethereum contracts.")
parser.add_argument("--keystore", help="The keystore file for the Ethereum account.")
parser.add_argument("--password", help="The password for the keystore file.")
parser.add_argument("--private-key", help="The private key for the Ethereum account.")
parser.add_argument("--rpc-url", default="http://127.0.0.1:8545", help="The RPC URL of the Ethereum node.")
parser.add_argument("--contracts", nargs='+', default=["BidderRegistry.sol", "PreConfCommitmentStore.sol", "BlockTracker.sol", "Oracle.sol", "ProviderRegistry.sol"], help="List of contract files to deploy.")
parser.add_argument("--gas-limit", type=int, default=3000000, help="The gas limit for the transactions.")
parser.add_argument("--chain-id", type=int, help="The chain ID of the Ethereum network.")
parser.add_argument("--out-dir", default="out", help="The directory containing the compiled contract outputs.")
parser.add_argument("--skip-deploy-create2", action='store_true', help="Skip deploying the CREATE2 proxy contract.")
parser.add_argument("--owner-addresses", help="Comma-separated list of contract:owner_address pairs for ownership transfer.")
parser.add_argument("--registry-contract", help="Address of the registry contract to register deployed contracts.")

args = parser.parse_args()

KEYSTORE = args.keystore
PASSWORD = args.password
PRIVATE_KEY = args.private_key
RPC_URL = args.rpc_url
CONTRACTS = args.contracts
GAS_LIMIT = args.gas_limit
OUT_DIR = args.out_dir
SKIP_DEPLOY_CREATE2 = args.skip_deploy_create2
OWNER_ADDRESSES_STR = args.owner_addresses
REGISTRY_CONTRACT = args.registry_contract

OWNER_ADDRESSES = parse_owner_addresses(OWNER_ADDRESSES_STR)

if (KEYSTORE and PRIVATE_KEY) or (not KEYSTORE and not PRIVATE_KEY):
    print("You must provide either a keystore file and password, or a private key, but not both.")
    exit(1)

w3 = Web3(Web3.HTTPProvider(RPC_URL))
if not w3.is_connected():
    print("Failed to connect to the Ethereum node.")
    exit(1)

check_geth_running(w3)

if KEYSTORE:
    account = load_account_from_keystore(KEYSTORE, PASSWORD)
else:
    account = w3.eth.account.from_key(PRIVATE_KEY)

# Get chain ID from Geth if not specified
CHAIN_ID = args.chain_id if args.chain_id else w3.eth.chain_id

if not SKIP_DEPLOY_CREATE2:
    deploy_create2(w3, account, CHAIN_ID)

if not REGISTRY_CONTRACT:
    REGISTRY_CONTRACT, registry_abi = deploy_registry_contract(w3, account, OUT_DIR, GAS_LIMIT, CHAIN_ID)
    if not REGISTRY_CONTRACT:
        print("Failed to deploy or locate the registry contract.")
        exit(1)
else:
    registry_abi = None
    # If registry_contract is provided, load its ABI
    registry_json_file = f"{OUT_DIR}/ContractRegistry.sol/ContractRegistry.json"
    with open(registry_json_file) as f:
        registry_contract_json = json.load(f)
        registry_abi = registry_contract_json["abi"]

for contract in CONTRACTS:
    deploy_contract(w3, account, contract, OUT_DIR, GAS_LIMIT, CHAIN_ID, OWNER_ADDRESSES, REGISTRY_CONTRACT, registry_abi)

print("All contracts deployed!")
