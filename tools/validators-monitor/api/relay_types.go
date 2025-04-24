package api

// RelayResult represents the result of querying a relay
type RelayResult struct {
	Relay      string `json:"relay"`
	StatusCode int    `json:"status_code"`
	Response   any    `json:"response,omitempty"`
	Error      string `json:"error,omitempty"`
}

// BidTrace represents a bid trace from a relay
type BidTrace struct {
	Slot                 string `json:"slot"`
	ParentHash           string `json:"parent_hash"`
	BlockHash            string `json:"block_hash"`
	BuilderPubkey        string `json:"builder_pubkey"`
	ProposerPubkey       string `json:"proposer_pubkey"`
	ProposerFeeRecipient string `json:"proposer_fee_recipient"`
	GasLimit             string `json:"gas_limit"`
	GasUsed              string `json:"gas_used"`
	Value                string `json:"value"`
	NumTx                string `json:"num_tx"`
	BlockNumber          string `json:"block_number"`
}
