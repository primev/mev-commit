-- init-db.sql
-- This script will be automatically executed when PostgreSQL starts for the first time

-- Create the execution_payloads table with proper indexing
CREATE TABLE IF NOT EXISTS execution_payloads (
    id SERIAL PRIMARY KEY,
    payload_id VARCHAR(66) UNIQUE NOT NULL, -- e.g., 0x... (32 bytes hex + 0x prefix)
    raw_execution_payload TEXT NOT NULL,
    block_height BIGINT NOT NULL UNIQUE,
    inserted_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_execution_payloads_block_height ON execution_payloads(block_height);
CREATE INDEX IF NOT EXISTS idx_execution_payloads_inserted_at ON execution_payloads(inserted_at);
CREATE INDEX IF NOT EXISTS idx_execution_payloads_payload_id ON execution_payloads(payload_id);

-- Create a partial index for recent payloads (optimization for common queries)
CREATE INDEX IF NOT EXISTS idx_execution_payloads_recent 
ON execution_payloads(block_height DESC) 
WHERE inserted_at > NOW() - INTERVAL '24 hours';
