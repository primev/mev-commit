package bridge

import (
	"context"
	"fmt"
	"math/big"
	"math/rand/v2"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/primev/mev-commit/bridge/standard/pkg/transfer"
	"github.com/primev/mev-commit/testing/pkg/orchestrator"
	"github.com/primev/mev-commit/x/keysigner"
)

type BridgeTestConfig struct {
	BridgeAccount          keysigner.KeySigner
	L1RPCURL               string
	SettlementRPCURL       string
	L1ContractAddr         common.Address
	SettlementContractAddr common.Address
}

func RunBridge(ctx context.Context, cluster orchestrator.Orchestrator, cfg any) error {
	bridgeTestConf, ok := cfg.(BridgeTestConfig)
	if !ok {
		return fmt.Errorf("unexpected config type")
	}

	logger := cluster.Logger().With("test", "bridge")

	minWeiValue := big.NewInt(params.Ether / 100)         // Enforce minimum value of 0.01 ETH.
	randWeiValue := big.NewInt(rand.Int64N(params.Ether)) // Generate a random amount of wei in [0.01, 1] ETH
	if randWeiValue.Cmp(minWeiValue) < 0 {
		randWeiValue = minWeiValue
	}

	// Create and start the transfer to the settlement chain
	tSettlement, err := transfer.NewTransferToSettlement(
		randWeiValue,
		bridgeTestConf.BridgeAccount.GetAddress(),
		bridgeTestConf.BridgeAccount,
		bridgeTestConf.SettlementRPCURL,
		bridgeTestConf.L1RPCURL,
		bridgeTestConf.L1ContractAddr,
		bridgeTestConf.SettlementContractAddr,
	)
	if err != nil {
		logger.Error("failed to create transfer to settlement", "error", err)
		return err
	}
	cctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	statusC := tSettlement.Do(cctx)
	for status := range statusC {
		if status.Error != nil {
			logger.Error("failed transfer to settlement", "error", status.Error)
			return status.Error
		}
		logger.Info("transfer to settlement status", "message", status.Message)
	}
	logger.Info("completed settlement transfer",
		"amount", randWeiValue.String(),
		"address", bridgeTestConf.BridgeAccount.GetAddress().String(),
	)

	// Sleep for random interval between 0 and 5 seconds
	time.Sleep(time.Duration(rand.IntN(6)) * time.Second)

	// Bridge back same amount minus 0.009 ETH for fees
	pZZNineEth := big.NewInt(9 * params.Ether / 1000)
	amountBack := new(big.Int).Sub(randWeiValue, pZZNineEth)

	// Create and start the transfer back to L1 with the same amount
	tL1, err := transfer.NewTransferToL1(
		amountBack,
		bridgeTestConf.BridgeAccount.GetAddress(),
		bridgeTestConf.BridgeAccount,
		bridgeTestConf.SettlementRPCURL,
		bridgeTestConf.L1RPCURL,
		bridgeTestConf.L1ContractAddr,
		bridgeTestConf.SettlementContractAddr,
	)
	if err != nil {
		logger.Error("failed to create transfer to L1", "error", err)
		return err
	}
	cctx, cancel = context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	statusC = tL1.Do(ctx)
	for status := range statusC {
		if status.Error != nil {
			logger.Error("failed transfer to L1", "error", status.Error)
			return status.Error
		}
		logger.Info("transfer to L1 status", "message", status.Message)
	}
	logger.Info("completed L1 transfer",
		"amount", amountBack.String(),
		"address", bridgeTestConf.BridgeAccount.GetAddress().String(),
	)
	return nil
}
