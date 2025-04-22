package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	validatorrouter "github.com/primev/mev-commit/contracts-abi/clients/ValidatorOptInRouter"
)

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	Close()
}

type RouterContract interface {
	AreValidatorsOptedIn(opts *bind.CallOpts, valBLSPubKeys [][]byte) ([]validatorrouter.IValidatorOptInRouterOptInStatus, error)
}

type ValidatorOptInChecker struct {
	client         EthClient
	routerContract RouterContract
	optsGetter     func() (*bind.CallOpts, error)
}

// NewValidatorOptInChecker creates a new ValidatorOptInChecker
func NewValidatorOptInChecker(
	rpcURL,
	contractAddress string,
	laggardMode *big.Int,
) (*ValidatorOptInChecker, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %s, %v", rpcURL, err)
	}

	callOptsGetter := func() (*bind.CallOpts, error) {
		blkNum, err := client.BlockNumber(context.Background())
		if err != nil {
			return nil, err
		}
		currentBlkNum := big.NewInt(0).SetUint64(blkNum)
		queryBlkNum := big.NewInt(0).Sub(currentBlkNum, laggardMode)
		return &bind.CallOpts{
			BlockNumber: queryBlkNum,
		}, nil
	}

	address := common.HexToAddress(contractAddress)
	routerContract, err := validatorrouter.NewValidatoroptinrouter(address, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create contract binding: %v", err)
	}

	return &ValidatorOptInChecker{
		client:         client,
		routerContract: routerContract,
		optsGetter:     callOptsGetter,
	}, nil
}

// CheckValidatorsOptedIn checks which validators in a batch are opted in
func (c *ValidatorOptInChecker) CheckValidatorsOptedIn(
	ctx context.Context,
	pubkeys []string,
) ([]validatorrouter.IValidatorOptInRouterOptInStatus, error) {
	blsPubKeys := make([][]byte, len(pubkeys))
	for j, pubkey := range pubkeys {
		pubkeyStr := strings.TrimPrefix(pubkey, "0x")
		pubkeyBytes, err := hex.DecodeString(pubkeyStr)
		if err != nil {
			return nil, fmt.Errorf("error decoding pubkey %s: %w", pubkey, err)
		}
		blsPubKeys[j] = pubkeyBytes
	}

	opts, err := c.optsGetter()
	if err != nil {
		return nil, fmt.Errorf("getting call opts: %w", err)
	}

	optInStatuses, err := c.routerContract.AreValidatorsOptedIn(opts, blsPubKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	return optInStatuses, nil
}

// Close closes the underlying ethclient connection
func (c *ValidatorOptInChecker) Close() {
	c.client.Close()
}
