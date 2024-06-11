from web3 import Web3
import argparse

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

def fetch_contract_code(w3, address):
    try:
        checksum_address = Web3.to_checksum_address(address)
        if checksum_address != address:
            print(f"Address {address} is not in checksum format")
            return
        code = w3.eth.get_code(checksum_address)
        if code != b'':
            print(f"Contract code at address {checksum_address}: {code.hex()}")
        else:
            print(f"No contract deployed at address {checksum_address}")
    except ValueError as e:
        print(f"Invalid address {address}: {e}")

# Command-line argument parsing
parser = argparse.ArgumentParser(description="Fetch contract code from Ethereum node.")
parser.add_argument("--rpc-url", default="http://127.0.0.1:8545", help="The RPC URL of the Ethereum node.")
parser.add_argument("--addresses", nargs='+', required=True, help="List of contract addresses to fetch code from.")

args = parser.parse_args()

RPC_URL = args.rpc_url
ADDRESSES = args.addresses

# Connect to the Ethereum node
w3 = Web3(Web3.HTTPProvider(RPC_URL))
if not w3.is_connected():
    print("Failed to connect to the Ethereum node.")
    exit(1)

# Check if Geth is running
check_geth_running(w3)

# Fetch and print the contract code for each address
for address in ADDRESSES:
    fetch_contract_code(w3, address)
