package blockbuilder

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/primev/mev-commit/cl/mocks"
	"github.com/primev/mev-commit/cl/redisapp/state"
	"github.com/primev/mev-commit/cl/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	etypes "github.com/ethereum/go-ethereum/core/types"

	redismock "github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

var handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
})
var stLog = slog.New(handler)
var buildDelay = time.Duration(1 * time.Second)

type MockEngineClient struct {
	mock.Mock
}

func (m *MockEngineClient) ForkchoiceUpdatedV3(ctx context.Context, fcs engine.ForkchoiceStateV1, attrs *engine.PayloadAttributes) (engine.ForkChoiceResponse, error) {
	args := m.Called(ctx, fcs, attrs)
	return args.Get(0).(engine.ForkChoiceResponse), args.Error(1)
}

func (m *MockEngineClient) GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error) {
	args := m.Called(ctx, payloadID)
	return args.Get(0).(*engine.ExecutionPayloadEnvelope), args.Error(1)
}

func (m *MockEngineClient) NewPayloadV4(ctx context.Context, executionPayload engine.ExecutableData, versionHashes []common.Hash, randao *common.Hash, executionRequests []hexutil.Bytes) (engine.PayloadStatusV1, error) {
	args := m.Called(ctx, executionPayload, versionHashes, randao, executionRequests)
	return args.Get(0).(engine.PayloadStatusV1), args.Error(1)
}

func (m *MockEngineClient) HeaderByNumber(ctx context.Context, number *big.Int) (*etypes.Header, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*etypes.Header), args.Error(1)
}

func createMockRPCClient() rpcClient {
	mockRPC := &MockRPCClient{}
	return mockRPC
}

type MockRPCClient struct{}

func (m *MockRPCClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	if method == "txpool_status" && result != nil {
		if mempoolStatus, ok := result.(*MempoolStatus); ok {
			mempoolStatus.Pending = 1
			mempoolStatus.Queued = 0
		}
	}
	return nil
}

func TestBlockBuilder_startBuild(t *testing.T) {
	ctx := context.Background()

	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x01, 0x02, 0x03},
		BlockHeight: 100,
		BlockTime:   uint64(time.Now().UnixMilli()) - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}
	timestamp := time.Now()

	hash := common.BytesToHash(executionHead.BlockHash)
	expectedFCS := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: engine.VALID,
		},
		PayloadID: &engine.PayloadID{0x01, 0x02, 0x03},
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, expectedFCS, mock.MatchedBy(matchPayloadAttributes(hash, executionHead.BlockTime))).Return(forkChoiceResponse, nil)

	resp, err := blockBuilder.startBuild(ctx, uint64(timestamp.UnixMilli()))

	require.NoError(t, err)
	assert.Equal(t, forkChoiceResponse, resp)

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_getPayload(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := mocks.NewMockRedisClient(ctrl)
	mockPipeliner := mocks.NewMockPipeliner(ctrl)

	timestamp := time.Now()
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x01, 0x02, 0x03},
		BlockHeight: 100,
		BlockTime:   uint64(timestamp.UnixMilli()),
	}

	mockRedisClient.EXPECT().
		XGroupCreateMkStream(gomock.Any(), "mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").Return(redis.NewStatusCmd(ctx))

	mockRedisClient.EXPECT().TxPipeline().Return(mockPipeliner)

	mockPipeliner.EXPECT().Set(ctx, "blockBuildState:instanceID123", gomock.Any(), time.Duration(0)).Return(redis.NewStatusCmd(ctx))
	mockPipeliner.EXPECT().XAdd(ctx, gomock.Any()).Return(redis.NewStringCmd(ctx, "result"))

	mockPipeliner.EXPECT().Exec(ctx).Return([]redis.Cmder{}, nil)

	stateManager, err := state.NewRedisCoordinator("instanceID123", mockRedisClient, nil)
	require.NoError(t, err)

	mockEngineClient := new(MockEngineClient)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	hash := common.BytesToHash(executionHead.BlockHash)
	ts := uint64(time.Now().UnixMilli())
	expectedFCS := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	payloadID := &engine.PayloadID{0x01, 0x02, 0x03}
	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: engine.VALID,
		},
		PayloadID: payloadID,
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, expectedFCS, mock.MatchedBy(matchPayloadAttributes(hash, executionHead.BlockTime))).Return(forkChoiceResponse, nil)

	executionPayload := &engine.ExecutionPayloadEnvelope{
		ExecutionPayload: &engine.ExecutableData{
			BlockHash:    hash,
			Number:       101,
			Timestamp:    ts,
			Random:       hash,
			ParentHash:   hash,
			ReceiptsRoot: hash,
		},
	}
	mockEngineClient.On("GetPayloadV4", mock.Anything, *payloadID).Return(executionPayload, nil)

	blockBuilder.executionHead = executionHead
	err = blockBuilder.GetPayload(ctx)

	require.NoError(t, err)

	mockEngineClient.AssertExpectations(t)
}

func TestBlockBuilder_FinalizeBlock(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedisClient := mocks.NewMockRedisClient(ctrl)

	mockRedisClient.EXPECT().XGroupCreateMkStream(ctx, "mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").Return(redis.NewStatusCmd(ctx))

	timestamp := uint64(1728051707) // 0x66fff9fb
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0xb, 0xf3, 0x9b, 0xc1, 0x8b, 0xe0, 0x59, 0xc1, 0xdc, 0xb8, 0x72, 0xac, 0x8c, 0xb, 0xc, 0x84, 0x56, 0x55, 0xa0, 0x1c, 0x2b, 0x7d, 0x8f, 0xd0, 0x1c, 0x4b, 0xec, 0xde, 0x6b, 0x3f, 0x93, 0xd7},
		BlockHeight: 2,
		BlockTime:   timestamp - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", mockRedisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	payloadIDStr := "payloadID123"

	executionPayloadStr := `{
		"parentHash":"0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7",
		"feeRecipient":"0x0000000000000000000000000000000000000000",
		"stateRoot":"0xcdc166a6c2e7f8b873889a7256873144e61121f9fc1f027d79b8fa310b91ff0f",
		"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"prevRandao":"0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7",
		"blockNumber":"0x3",
		"gasLimit":"0x1c9c380",
		"gasUsed":"0x0",
		"timestamp":"0x66fff9fb",
		"extraData":"0xd983010e08846765746888676f312e32322e368664617277696e",
		"baseFeePerGas":"0x27ee3253",
		"blockHash":"0x9a9b2f7e98934f8544c22cdcb00526f48886170b15c4e4e96bd43af189b5aac4",
		"transactions":[],
		"withdrawals":[],
		"blobGasUsed":"0x0",
		"excessBlobGas":"0x0"
	}`

	msgID := ""

	var executionPayload engine.ExecutableData
	err = json.Unmarshal([]byte(executionPayloadStr), &executionPayload)
	require.NoError(t, err)

	// Marshal struct to msgpack
	msgpackData, err := msgpack.Marshal(executionPayload)
	require.NoError(t, err)

	encodedPayload := base64.StdEncoding.EncodeToString(msgpackData)

	payloadStatus := engine.PayloadStatusV1{
		Status: engine.VALID,
	}

	mockEngineClient.On("NewPayloadV4", mock.Anything, executionPayload, []common.Hash{}, mock.Anything, mock.Anything).Return(payloadStatus, nil)

	hash := executionPayload.BlockHash
	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}
	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: payloadStatus,
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, fcs, (*engine.PayloadAttributes)(nil)).Return(forkChoiceResponse, nil)

	blockBuilder.executionHead = executionHead
	err = blockBuilder.FinalizeBlock(ctx, payloadIDStr, encodedPayload, msgID)

	require.NoError(t, err)

	mockEngineClient.AssertExpectations(t)
}

func TestBlockBuilder_startBuild_ForkchoiceUpdatedError(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x01, 0x02, 0x03},
		BlockHeight: 100,
		BlockTime:   uint64(time.Now().UnixMilli()) - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	timestamp := time.Now()

	hash := common.BytesToHash(executionHead.BlockHash)
	expectedFCS := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, expectedFCS, mock.MatchedBy(matchPayloadAttributes(hash, executionHead.BlockTime))).Return(engine.ForkChoiceResponse{}, errors.New("engine error"))

	resp, err := blockBuilder.startBuild(ctx, uint64(timestamp.UnixMilli()))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "forkchoice update")
	assert.Equal(t, engine.ForkChoiceResponse{}, resp)

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_startBuild_InvalidPayloadStatus(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x01, 0x02, 0x03},
		BlockHeight: 100,
		BlockTime:   uint64(time.Now().UnixMilli()) - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)

	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	timestamp := time.Now()

	hash := common.BytesToHash(executionHead.BlockHash)
	expectedFCS := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: "INVALID_STATUS",
		},
		PayloadID: nil,
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, expectedFCS, mock.MatchedBy(matchPayloadAttributes(hash, executionHead.BlockTime))).Return(forkChoiceResponse, nil)

	resp, err := blockBuilder.startBuild(ctx, uint64(timestamp.UnixMilli()))

	require.NoError(t, err)
	assert.Equal(t, forkChoiceResponse, resp)

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_getPayload_GetPayloadUnknownPayload(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	timestamp := time.Now()
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x01, 0x02, 0x03},
		BlockHeight: 100,
		BlockTime:   uint64(timestamp.UnixMilli()) - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)
	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   time.Duration(1 * time.Second),
		logger:       stLog,
	}

	hash := common.BytesToHash(executionHead.BlockHash)

	expectedFCS := engine.ForkchoiceStateV1{
		HeadBlockHash:      hash,
		SafeBlockHash:      hash,
		FinalizedBlockHash: hash,
	}

	payloadID := &engine.PayloadID{0x01, 0x02, 0x03}
	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: engine.VALID,
		},
		PayloadID: payloadID,
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", mock.Anything, expectedFCS, mock.MatchedBy(matchPayloadAttributes(hash, executionHead.BlockTime))).Return(forkChoiceResponse, nil)

	mockEngineClient.On("GetPayloadV4", mock.Anything, *payloadID).Return(&engine.ExecutionPayloadEnvelope{}, errors.New("Unknown payload"))

	blockBuilder.executionHead = executionHead
	err = blockBuilder.GetPayload(ctx)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get payload")

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_FinalizeBlock_InvalidBlockHeight(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	timestamp := time.Now()
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0x00, 0x00, 0x00},
		BlockHeight: 100,
		BlockTime:   uint64(timestamp.UnixMilli()) - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)
	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	payloadIDStr := "payloadID123"
	executionPayload := &engine.ExecutableData{
		ParentHash:    common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		FeeRecipient:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		StateRoot:     common.HexToHash("0xcdc166a6c2e7f8b873889a7256873144e61121f9fc1f027d79b8fa310b91ff0f"),
		ReceiptsRoot:  common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		LogsBloom:     common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		Random:        common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		Number:        3,
		GasLimit:      30000000,
		GasUsed:       0,
		Timestamp:     0x66fff9fb,
		ExtraData:     common.FromHex("0xd983010e08846765746888676f312e32322e368664617277696e"),
		BaseFeePerGas: big.NewInt(0x27ee3253),
		BlockHash:     common.HexToHash("0x9a9b2f7e98934f8544c22cdcb00526f48886170b15c4e4e96bd43af189b5aac4"),
		Transactions:  [][]byte{},
		Withdrawals:   []*etypes.Withdrawal{},
		BlobGasUsed:   new(uint64),
		ExcessBlobGas: new(uint64),
	}
	executionPayloadData, _ := msgpack.Marshal(executionPayload)

	executionPayloadEncoded := base64.StdEncoding.EncodeToString(executionPayloadData)
	blockBuilder.executionHead = executionHead
	err = blockBuilder.FinalizeBlock(ctx, payloadIDStr, executionPayloadEncoded, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid block height")

	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_FinalizeBlock_NewPayloadInvalidStatus(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	timestamp := uint64(1728051707) // 0x66fff9fb
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0xb, 0xf3, 0x9b, 0xc1, 0x8b, 0xe0, 0x59, 0xc1, 0xdc, 0xb8, 0x72, 0xac, 0x8c, 0xb, 0xc, 0x84, 0x56, 0x55, 0xa0, 0x1c, 0x2b, 0x7d, 0x8f, 0xd0, 0x1c, 0x4b, 0xec, 0xde, 0x6b, 0x3f, 0x93, 0xd7},
		BlockHeight: 2,
		BlockTime:   timestamp - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)
	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	payloadIDStr := "payloadID123"
	executionPayload := engine.ExecutableData{
		ParentHash:    common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		FeeRecipient:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		StateRoot:     common.HexToHash("0xcdc166a6c2e7f8b873889a7256873144e61121f9fc1f027d79b8fa310b91ff0f"),
		ReceiptsRoot:  common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		LogsBloom:     common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		Random:        common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		Number:        3,
		GasLimit:      30000000,
		GasUsed:       0,
		Timestamp:     0x66fff9fb,
		ExtraData:     common.FromHex("0xd983010e08846765746888676f312e32322e368664617277696e"),
		BaseFeePerGas: big.NewInt(0x27ee3253),
		BlockHash:     common.HexToHash("0x9a9b2f7e98934f8544c22cdcb00526f48886170b15c4e4e96bd43af189b5aac4"),
		Transactions:  [][]byte{},
		Withdrawals:   []*etypes.Withdrawal{},
		BlobGasUsed:   new(uint64),
		ExcessBlobGas: new(uint64),
	}

	executionPayloadData, _ := msgpack.Marshal(executionPayload)
	executionPayloadEncoded := base64.StdEncoding.EncodeToString(executionPayloadData)
	payloadStatus := engine.PayloadStatusV1{
		Status: "INVALID",
	}
	mockEngineClient.On("NewPayloadV4", mock.Anything, executionPayload, []common.Hash{}, mock.Anything, mock.Anything).Return(payloadStatus, nil)
	blockBuilder.executionHead = executionHead
	err = blockBuilder.FinalizeBlock(ctx, payloadIDStr, executionPayloadEncoded, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to push new payload")

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func TestBlockBuilder_FinalizeBlock_ForkchoiceUpdatedInvalidStatus(t *testing.T) {
	ctx := context.Background()
	redisClient, redisMock := redismock.NewClientMock()
	redisMock.ExpectXGroupCreateMkStream("mevcommit_block_stream", "mevcommit_consumer_group:instanceID123", "0").SetVal("OK")

	timestamp := uint64(1728051707) // 0x66fff9fb
	executionHead := &types.ExecutionHead{
		BlockHash:   []byte{0xb, 0xf3, 0x9b, 0xc1, 0x8b, 0xe0, 0x59, 0xc1, 0xdc, 0xb8, 0x72, 0xac, 0x8c, 0xb, 0xc, 0x84, 0x56, 0x55, 0xa0, 0x1c, 0x2b, 0x7d, 0x8f, 0xd0, 0x1c, 0x4b, 0xec, 0xde, 0x6b, 0x3f, 0x93, 0xd7},
		BlockHeight: 2,
		BlockTime:   timestamp - 10,
	}

	stateManager, err := state.NewRedisCoordinator("instanceID123", redisClient, nil)
	require.NoError(t, err)
	mockEngineClient := new(MockEngineClient)
	blockBuilder := &BlockBuilder{
		stateManager: stateManager,
		engineCl:     mockEngineClient,
		rpcClient:    createMockRPCClient(),
		buildDelay:   buildDelay,
		logger:       stLog,
	}

	payloadIDStr := "payloadID123"
	executionPayload := engine.ExecutableData{
		ParentHash:    common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		FeeRecipient:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		StateRoot:     common.HexToHash("0xcdc166a6c2e7f8b873889a7256873144e61121f9fc1f027d79b8fa310b91ff0f"),
		ReceiptsRoot:  common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		LogsBloom:     common.FromHex("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		Random:        common.HexToHash("0x0bf39bc18be059c1dcb872ac8c0b0c845655a01c2b7d8fd01c4becde6b3f93d7"),
		Number:        3,
		GasLimit:      30000000,
		GasUsed:       0,
		Timestamp:     0x66fff9fb,
		ExtraData:     common.FromHex("0xd983010e08846765746888676f312e32322e368664617277696e"),
		BaseFeePerGas: big.NewInt(0x27ee3253),
		BlockHash:     common.HexToHash("0x9a9b2f7e98934f8544c22cdcb00526f48886170b15c4e4e96bd43af189b5aac4"),
		Transactions:  [][]byte{},
		Withdrawals:   []*etypes.Withdrawal{},
		BlobGasUsed:   new(uint64),
		ExcessBlobGas: new(uint64),
	}
	executionPayloadData, _ := msgpack.Marshal(executionPayload)
	executionPayloadEncoded := base64.StdEncoding.EncodeToString(executionPayloadData)
	payloadStatus := engine.PayloadStatusV1{
		Status: engine.VALID,
	}
	mockEngineClient.On("NewPayloadV4", mock.Anything, executionPayload, []common.Hash{}, mock.Anything, mock.Anything).Return(payloadStatus, nil)

	fcs := engine.ForkchoiceStateV1{
		HeadBlockHash:      executionPayload.BlockHash,
		SafeBlockHash:      executionPayload.BlockHash,
		FinalizedBlockHash: executionPayload.BlockHash,
	}
	forkChoiceResponse := engine.ForkChoiceResponse{
		PayloadStatus: engine.PayloadStatusV1{
			Status: "INVALID",
		},
	}
	mockEngineClient.On("ForkchoiceUpdatedV3", ctx, fcs, (*engine.PayloadAttributes)(nil)).Return(forkChoiceResponse, nil)

	blockBuilder.executionHead = executionHead
	err = blockBuilder.FinalizeBlock(ctx, payloadIDStr, executionPayloadEncoded, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to finalize fork choice update")

	mockEngineClient.AssertExpectations(t)
	require.NoError(t, redisMock.ExpectationsWereMet())
}

func matchPayloadAttributes(expectedHash common.Hash, executionHeadTime uint64) func(*engine.PayloadAttributes) bool {
	return func(attrs *engine.PayloadAttributes) bool {
		if attrs == nil {
			return false
		}
		if attrs.Random != expectedHash {
			return false
		}
		if attrs.SuggestedFeeRecipient != (common.Address{}) {
			return false
		}
		if attrs.BeaconRoot == nil || *attrs.BeaconRoot != expectedHash {
			return false
		}
		if len(attrs.Withdrawals) != 0 {
			return false
		}
		if attrs.Timestamp <= executionHeadTime {
			return false
		}
		return true
	}
}
