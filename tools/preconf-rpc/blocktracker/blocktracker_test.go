package blocktracker_test


func TestBlockTracker(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &mockEthClient{
		blockNumber: 100,
		blocks: map[uint64]*types.Block{
			100: {
				Number: big.NewInt(100),
				Hash:   common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				Transactions: []*types.Transaction{},
			},
			101: {
				Number: big.NewInt(101),
				Hash:   common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
				Transactions: []*types.Transaction{},
			},
		},
	}

	tracker := &blockTracker{

