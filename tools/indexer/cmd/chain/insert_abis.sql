-- Insert BlockTracker v0.8 ABI
-- You'll need to calculate MD5 hash of the ABI JSON and insert
-- Format: INSERT INTO contract_abis (abi_hash, address, name, version, abi) VALUES (MD5('abi_json'), 'address', 'name', 'version', parse_json('abi_json'));

-- BlockTracker v0.8
-- Address: 0x0da2a367c51f2a34465acd6ae5d8a48385e9cb03
-- You need to insert the full ABI JSON here

-- BlockTracker v1.1.0
-- Address: 0x0da2a367c51f2a34465acd6ae5d8a48385e9cb03
-- You need to insert the full ABI JSON here

-- BidderRegistry v0.8
-- Address: 0xc973d09e51a20c9ab0214c439e4b34dbac52ad67

-- BidderRegistry v1.1.0
-- Address: 0xc973d09e51a20c9ab0214c439e4b34dbac52ad67

-- Template for manual insertion:
-- INSERT INTO contract_abis (abi_hash, address, name, version, abi)
-- VALUES (
--   MD5('your_abi_json_here'),
--   '0xcontract_address',
--   'ContractName',
--   'v0.8',
--   parse_json('your_abi_json_here')
-- );
