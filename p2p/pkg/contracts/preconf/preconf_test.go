package preconfcontract_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	preconfcontract "github.com/primevprotocol/mev-commit/p2p/pkg/contracts/preconf"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
	mockevmclient "github.com/primevprotocol/mev-commit/p2p/pkg/evmclient/mock"
	"github.com/primevprotocol/mev-commit/x/util"
)

func TestPreconfContract(t *testing.T) {
	t.Parallel()

	t.Run("StoreEncryptedCommitment", func(t *testing.T) {
		preConfContract := common.HexToAddress("abcd")
		txHash := common.HexToHash("abcdef")
		commitment := [32]byte([]byte("abcdefabcdefabcdefabcdefabcdefaa"))
		commitmentSignature := []byte("abcdef")

		expCallData, err := preconfcontract.PreConfABI().Pack(
			"storeEncryptedCommitment",
			commitment,
			commitmentSignature,
		)

		if err != nil {
			t.Fatal(err)
		}

		mockClient := mockevmclient.New(
			mockevmclient.WithSendFunc(
				func(ctx context.Context, req *evmclient.TxRequest) (common.Hash, error) {
					if req.To.Cmp(preConfContract) != 0 {
						t.Fatalf(
							"expected to address to be %s, got %s",
							preConfContract.Hex(), req.To.Hex(),
						)
					}
					if !bytes.Equal(req.CallData, expCallData) {
						t.Fatalf("expected call data to be %x, got %x", expCallData, req.CallData)
					}
					return txHash, nil
				},
			),
			mockevmclient.WithWaitForReceiptFunc(
				func(ctx context.Context, txnHash common.Hash) (*types.Receipt, error) {
					if txnHash != txHash {
						t.Fatalf("expected tx hash to be %s, got %s", txHash.Hex(), txnHash.Hex())
					}
					return &types.Receipt{
						Status: 1,
					}, nil
				},
			),
		)

		preConfContractClient := preconfcontract.New(
			preConfContract,
			mockClient,
			util.NewTestLogger(os.Stdout),
		)

		_, err = preConfContractClient.StoreEncryptedCommitment(
			context.Background(),
			commitment[:],
			commitmentSignature,
		)
		if err != nil {
			t.Fatal(err)
		}
	})
}
