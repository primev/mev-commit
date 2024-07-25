package debugapi_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	debugapiv1 "github.com/primev/mev-commit/p2p/gen/go/debugapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/p2p"
	debugapi "github.com/primev/mev-commit/p2p/pkg/rpc/debug"
	"github.com/primev/mev-commit/p2p/pkg/topology"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	"github.com/stretchr/testify/assert"
)

type mockStore struct{}

func (m *mockStore) PendingTxns() ([]*txmonitor.TxnDetails, error) {
	return []*txmonitor.TxnDetails{
		{
			Hash:    common.HexToHash("0x00001"),
			Nonce:   1,
			Created: time.Now().Unix(),
		},
		{
			Hash:    common.HexToHash("0x00002"),
			Nonce:   2,
			Created: time.Now().Unix(),
		},
	}, nil
}

type mockCanceller struct{}

func (m *mockCanceller) CancelTx(ctx context.Context, txHash common.Hash) (common.Hash, error) {
	return txHash, nil
}

type mockP2PService struct{}

func (m *mockP2PService) Self() map[string]interface{} {
	return map[string]interface{}{
		"self_key": "self_value",
	}
}

func (m *mockP2PService) BlockedPeers() []p2p.BlockedPeerInfo {
	return []p2p.BlockedPeerInfo{
		{
			Peer:   common.HexToAddress("0xab"),
			Reason: "reason1",
		},
		{
			Peer:   common.HexToAddress("0xcd"),
			Reason: "reason2",
		},
	}
}

type mockTopology struct{}

func (m *mockTopology) GetPeers(q topology.Query) []p2p.Peer {
	if q.Type == p2p.PeerTypeProvider {
		return []p2p.Peer{
			{
				EthAddress: common.HexToAddress("0x11111"),
			},
			{
				EthAddress: common.HexToAddress("0x22222"),
			},
		}
	} else if q.Type == p2p.PeerTypeBidder {
		return []p2p.Peer{
			{
				EthAddress: common.HexToAddress("0x33333"),
			},
			{
				EthAddress: common.HexToAddress("0x44444"),
			},
		}
	}
	return nil
}

func TestService_GetTopology(t *testing.T) {
	service := debugapi.NewService(&mockStore{}, &mockCanceller{}, &mockP2PService{}, &mockTopology{})

	ctx := context.Background()
	req := &debugapiv1.EmptyMessage{}

	resp, err := service.GetTopology(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Topology)

	self, ok := resp.Topology.Fields["self"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(self.GetStructValue().Fields))
	assert.Equal(t, "self_value", self.GetStructValue().Fields["self_key"].GetStringValue())

	connectedProviders, ok := resp.Topology.Fields["connected_providers"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(connectedProviders.GetListValue().Values))
	assert.Equal(t, common.HexToAddress("0x11111").String(), connectedProviders.GetListValue().Values[0].GetStringValue())
	assert.Equal(t, common.HexToAddress("0x22222").String(), connectedProviders.GetListValue().Values[1].GetStringValue())

	connectedBidders, ok := resp.Topology.Fields["connected_bidders"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(connectedBidders.GetListValue().Values))
	assert.Equal(t, common.HexToAddress("0x33333").String(), connectedBidders.GetListValue().Values[0].GetStringValue())
	assert.Equal(t, common.HexToAddress("0x44444").String(), connectedBidders.GetListValue().Values[1].GetStringValue())

	blockedPeers, ok := resp.Topology.Fields["blocked_peers"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(blockedPeers.GetListValue().Values))
	for _, v := range blockedPeers.GetListValue().Values {
		assert.Equal(t, 2, len(v.GetStructValue().Fields))
		if v.GetStructValue().Fields["peer"].GetStringValue() == common.HexToAddress("0xab").String() {
			assert.Equal(t, "reason1", v.GetStructValue().Fields["reason"].GetStringValue())
		} else if v.GetStructValue().Fields["peer"].GetStringValue() == common.HexToAddress("0xcd").String() {
			assert.Equal(t, "reason2", v.GetStructValue().Fields["reason"].GetStringValue())
		} else {
			assert.Fail(t, "unexpected peer")
		}
	}
}

func TestService_GetPendingTransactions(t *testing.T) {
	service := debugapi.NewService(&mockStore{}, &mockCanceller{}, &mockP2PService{}, &mockTopology{})

	ctx := context.Background()
	req := &debugapiv1.EmptyMessage{}

	resp, err := service.GetPendingTransactions(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.PendingTransactions))
	assert.Equal(t, common.HexToHash("0x00001").String(), resp.PendingTransactions[0].TxHash)
	assert.Equal(t, int64(1), resp.PendingTransactions[0].Nonce)
	assert.NotEmpty(t, resp.PendingTransactions[0].Created)
	assert.Equal(t, common.HexToHash("0x00002").String(), resp.PendingTransactions[1].TxHash)
	assert.Equal(t, int64(2), resp.PendingTransactions[1].Nonce)
	assert.NotEmpty(t, resp.PendingTransactions[1].Created)
}

func TestService_CancelTransaction(t *testing.T) {
	service := debugapi.NewService(&mockStore{}, &mockCanceller{}, &mockP2PService{}, &mockTopology{})

	ctx := context.Background()
	req := &debugapiv1.CancelTransactionReq{
		TxHash: common.HexToHash("0x12345").String(),
	}

	resp, err := service.CancelTransaction(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, common.HexToHash("0x12345").String(), resp.TxHash)
}
