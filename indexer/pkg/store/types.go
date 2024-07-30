package store

type IndexBlock struct {
	Number       uint64 `json:"number"`
	Hash         string `json:"hash"`
	ParentHash   string `json:"parentHash"`
	Root         string `json:"root"`
	Nonce        uint64 `json:"nonce"`
	Timestamp    string `json:"timestamp"`
	Transactions int    `json:"transactions"`
	BaseFee      uint64 `json:"baseFee"`
	GasLimit     uint64 `json:"gasLimit"`
	GasUsed      uint64 `json:"gasUsed"`
	Difficulty   uint64 `json:"difficulty"`
	ExtraData    string `json:"extraData"`
}

type IndexTransaction struct {
	Hash               string `json:"hash"`
	From               string `json:"from"`
	To                 string `json:"to"`
	Gas                uint64 `json:"gas"`
	GasPrice           uint64 `json:"gasPrice"`
	GasTipCap          uint64 `json:"gasTipCap"`
	GasFeeCap          uint64 `json:"gasFeeCap"`
	Value              string `json:"value"`
	Nonce              uint64 `json:"nonce"`
	BlockHash          string `json:"blockHash"`
	BlockNumber        uint64 `json:"blockNumber"`
	ChainId            string `json:"chainId"`
	V                  string `json:"v"`
	R                  string `json:"r"`
	S                  string `json:"s"`
	Input              string `json:"input"`
	Timestamp          string `json:"timestamp"`
	Status             uint64 `json:"status"`
	GasUsed            uint64 `json:"gasUsed"`
	CumulativeGasUsed  uint64 `json:"cumulativeGasUsed"`
	ContractAddress    string `json:"contractAddress"`
	TransactionIndex   uint   `json:"transactionIndex"`
	ReceiptBlockHash   string `json:"receiptBlockHash"`
	ReceiptBlockNumber uint64 `json:"receiptBlockNumber"`
}

type AccountBalance struct {
	Address     string `json:"address"`
	Balance     string `json:"balance"`
	Timestamp   string `json:"timestamp"`
	BlockNumber uint64 `json:"blockNumber"`
}
