package preconfcontract

import (
	"context"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	preconfcommitmentstore "github.com/primevprotocol/mev-commit/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-commit/p2p/pkg/evmclient"
)

var preconfABI = func() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(preconfcommitmentstore.PreconfcommitmentstoreMetaData.ABI))
	if err != nil {
		panic(err)
	}
	return abi
}

var defaultWaitTimeout = 10 * time.Second

type Interface interface {
	StoreCommitment(
		ctx context.Context,
		bid *big.Int,
		blockNumber uint64,
		txHash string,
		decayStartTimeStamp uint64,
		decayEndTimeStamp uint64,
		bidSignature []byte,
		commitmentSignature []byte,
		decayDispatchTimestamp uint64,
	) error
}

type preconfContract struct {
	preconfABI          abi.ABI
	preconfContractAddr common.Address
	client              evmclient.Interface
	logger              *slog.Logger
}

func New(
	preconfContractAddr common.Address,
	client evmclient.Interface,
	logger *slog.Logger,
) Interface {
	return &preconfContract{
		preconfABI:          preconfABI(),
		preconfContractAddr: preconfContractAddr,
		client:              client,
		logger:              logger,
	}
}

func (p *preconfContract) StoreCommitment(
	ctx context.Context,
	bid *big.Int,
	blockNumber uint64,
	txHash string,
	deacyStartTimeStamp uint64,
	decayEndTimeStamp uint64,
	bidSignature []byte,
	commitmentSignature []byte,
	decayDispatchTimestamp uint64,
) error {

	callData, err := p.preconfABI.Pack(
		"storeCommitment",
		uint64(bid.Int64()),
		blockNumber,
		txHash,
		deacyStartTimeStamp,
		decayEndTimeStamp,
		bidSignature,
		commitmentSignature,
		decayDispatchTimestamp,
	)
	if err != nil {
		p.logger.Error("preconf contract storeCommitment pack error", "err", err)
		return err
	}

	txnHash, err := p.client.Send(ctx, &evmclient.TxRequest{
		To:       &p.preconfContractAddr,
		CallData: callData,
	})
	if err != nil {
		return err
	}

	p.logger.Info("preconf contract storeCommitment successful", "txnHash", txnHash)

	return nil
}
