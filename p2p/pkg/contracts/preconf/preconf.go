package preconfcontract

import (
	"context"
	"fmt"
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
	StoreEncryptedCommitment(
		ctx context.Context,
		commitmentDigest []byte,
		commitmentSignature []byte,
		decayDispatchTimestamp uint64,
	) (common.Hash, error)
	OpenCommitment(
		ctx context.Context,
		encryptedCommitmentIndex []byte,
		bid string,
		blockNumber int64,
		txnHash string,
		decayStartTimeStamp int64,
		decayEndTimeStamp int64,
		bidSignature []byte,
		commitmentSignature []byte,
		sharedSecretKey []byte,
	) (common.Hash, error)
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

func (p *preconfContract) StoreEncryptedCommitment(
	ctx context.Context,
	commitmentDigest []byte,
	commitmentSignature []byte,
	decayDispatchTimestamp uint64,
) (common.Hash, error) {
	callData, err := p.preconfABI.Pack(
		"storeEncryptedCommitment",
		[32]byte(commitmentDigest),
		commitmentSignature,
		decayDispatchTimestamp,
	)
	if err != nil {
		p.logger.Error("preconf contract storeEncryptedCommitment pack error", "err", err)
		return common.Hash{}, err
	}

	txnHash, err := p.client.Send(ctx, &evmclient.TxRequest{
		To:       &p.preconfContractAddr,
		CallData: callData,
	})
	if err != nil {
		return common.Hash{}, err
	}

	return txnHash, nil
}

func (p *preconfContract) OpenCommitment(
	ctx context.Context,
	encryptedCommitmentIndex []byte,
	bid string,
	blockNumber int64,
	txnHash string,
	decayStartTimeStamp int64,
	decayEndTimeStamp int64,
	bidSignature []byte,
	commitmentSignature []byte,
	sharedSecretKey []byte,
) (common.Hash, error) {
	bidAmt, ok := new(big.Int).SetString(bid, 10)
	if !ok {
		p.logger.Error("Error converting bid to big.Int", "bid", bid)
		return common.Hash{}, fmt.Errorf("error converting bid to big.Int, bid: %s", bid)
	}

	var eciBytes [32]byte
	copy(eciBytes[:], encryptedCommitmentIndex)

	callData, err := p.preconfABI.Pack(
		"openCommitment",
		eciBytes,
		bidAmt.Uint64(),
		uint64(blockNumber),
		txnHash,
		uint64(decayStartTimeStamp),
		uint64(decayEndTimeStamp),
		bidSignature,
		commitmentSignature,
		sharedSecretKey,
	)
	if err != nil {
		p.logger.Error("Error packing call data for openCommitment", "error", err)
		return common.Hash{}, err
	}

	txHash, err := p.client.Send(ctx, &evmclient.TxRequest{
		To:       &p.preconfContractAddr,
		CallData: callData,
	})
	if err != nil {
		return common.Hash{}, err
	}

	p.logger.Info("preconf contract openCommitment successful", "txHash", txHash.String())

	return txHash, nil
}
