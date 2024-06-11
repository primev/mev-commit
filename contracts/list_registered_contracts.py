import json
import argparse
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

def load_account_from_keystore(keystore_file, password):
    with open(keystore_file) as keyfile:
        encrypted_key = keyfile.read()
        private_key = Account.decrypt(encrypted_key, password)
        return Account.from_key(private_key)

def list_registered_contracts(w3, registry_contract, registry_abi):
    contract = w3.eth.contract(address=registry_contract, abi=registry_abi)
    all_contracts = contract.functions.getAllContracts().call()
    print(f"Total registered contracts: {len(all_contracts)}")
    for contract_info in all_contracts:
        contract_name = contract_info[0]
        contract_address = contract_info[1]
        print(f"Contract Name: {contract_name}, Contract Address: {contract_address}")

parser = argparse.ArgumentParser(description="List registered Ethereum contracts.")
parser.add_argument("--keystore", help="The keystore file for the Ethereum account.")
parser.add_argument("--password", help="The password for the keystore file.")
parser.add_argument("--private-key", help="The private key for the Ethereum account.")
parser.add_argument("--rpc-url", default="http://127.0.0.1:8545", help="The RPC URL of the Ethereum node.")
parser.add_argument("--registry-contract", required=True, help="Address of the registry contract.")
parser.add_argument("--registry-abi", default="out/ContractRegistry.sol/ContractRegistry.json", help="Path to the ABI JSON file of the registry contract.")

args = parser.parse_args()

KEYSTORE = args.keystore
PASSWORD = args.password
PRIVATE_KEY = args.private_key
RPC_URL = args.rpc_url
REGISTRY_CONTRACT = args.registry_contract
REGISTRY_ABI_PATH = args.registry_abi

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

# Load the registry contract ABI
with open(REGISTRY_ABI_PATH) as f:
    registry_contract_json = json.load(f)
    registry_abi = registry_contract_json["abi"]

# List registered contracts
list_registered_contracts(w3, REGISTRY_CONTRACT, registry_abi)
