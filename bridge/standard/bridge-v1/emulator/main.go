package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	mathrand "math/rand"
	"os"
	"time"

	"github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/primevprotocol/mev-commit/bridge/standard/bridge-v1/pkg/shared"
	"github.com/primevprotocol/mev-commit/bridge/standard/bridge-v1/pkg/transfer"
	"github.com/primevprotocol/mev-commit/bridge/standard/bridge-v1/pkg/util"
)

const (
	settlementRPCUrl             = "http://sl-bootnode:8545"
	l1RPCUrl                     = "http://l1-bootnode:8545"
	l1ContractAddrString         = "0xe7B3D66A143Bfb85Cd980d1577Ba52f68F8f5809"
	settlementContractAddrString = "0xC157631f4aaaaabAB864544d743c5726a68C85a4"
)

func main() {
	logger, err := util.NewLogger(slog.LevelDebug.String(), "text", "", os.Stdout)
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		os.Exit(1)
	}

	privateKeyString := os.Getenv("PRIVATE_KEY")
	if privateKeyString == "" {
		logger.Error("PRIVATE_KEY env var is required")
		os.Exit(1)
	}
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		logger.Error("failed to parse private key")
		os.Exit(1)
	}

	transferAddressString := os.Getenv("ACCOUNT_ADDR")
	if transferAddressString == "" {
		logger.Error("ACCOUNT_ADDR env var is required")
		os.Exit(1)
	}
	if !common.IsHexAddress(transferAddressString) {
		logger.Error("ACCOUNT_ADDR is not a valid address")
		os.Exit(1)
	}
	transferAddr := common.HexToAddress(transferAddressString)

	l1ContractAddr := common.HexToAddress(l1ContractAddrString)
	settlementContractAddr := common.HexToAddress(settlementContractAddrString)

	// DD setup
	ctx := context.WithValue(context.Background(), datadog.ContextAPIKeys, map[string]datadog.APIKey{
		"apiKeyAuth": {
			Key: os.Getenv("DD_API_KEY"),
		},
		"appKeyAuth": {
			Key: os.Getenv("DD_APP_KEY"),
		},
	})

	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)

	// Construct two eth clients and cancel all pending txes on both chains
	l1Client, err := ethclient.Dial(l1RPCUrl)
	if err != nil {
		logger.Error("failed to dial l1 rpc", "error", err)
		os.Exit(1)
	}
	settlementClient, err := ethclient.Dial(settlementRPCUrl)
	if err != nil {
		logger.Error("failed to dial settlement rpc", "error", err)
		os.Exit(1)
	}
	if err := shared.NewETHClient(logger.With("component", "l1_eth_client"), l1Client).CancelPendingTxes(ctx, privateKey); err != nil {
		logger.Error("failed to cancel pending L1 transactions")
		os.Exit(1)
	}
	if err := shared.NewETHClient(logger.With("component", "settlement_eth_client"), settlementClient).CancelPendingTxes(ctx, privateKey); err != nil {
		logger.Error("failed to cancel pending settlement transactions")
		os.Exit(1)
	}

	for {
		// Generate a random amount of wei in [0.01, 10] ETH
		maxWei := new(big.Int).Mul(big.NewInt(10), big.NewInt(params.Ether))
		randWeiValue, err := rand.Int(rand.Reader, maxWei)
		if err != nil {
			logger.Error("failed to generate random value", "error", err)
			os.Exit(1)
		}
		if randWeiValue.Cmp(big.NewInt(params.Ether/100)) < 0 {
			// Enforce minimum value of 0.01 ETH
			randWeiValue = big.NewInt(params.Ether / 100)
		}

		// TODO: support persistent connections

		// Create and start the transfer to the settlement chain
		tSettlement, err := transfer.NewTransferToSettlement(
			logger.With("component", "settlement_transfer"),
			randWeiValue,
			transferAddr,
			privateKey,
			settlementRPCUrl,
			l1RPCUrl,
			l1ContractAddr,
			settlementContractAddr,
		)
		if err != nil {
			tags := []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "17864"}
			if err := postMetricToDatadog(ctx, apiClient, "bridging.failure", 0, tags); err != nil {
				logger.Error("failed to post metric", "error", err)
			}
			time.Sleep(time.Minute)
			continue
		}
		startTime := time.Now()
		err = tSettlement.Start(ctx)
		if err != nil {
			tags := []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "17864"}
			if err := postMetricToDatadog(ctx, apiClient, "bridging.failure", time.Since(startTime).Seconds(), tags); err != nil {
				logger.Error("failed to post metric", "error", err)
			}
			time.Sleep(time.Minute)
			continue
		}
		completionTimeSec := time.Since(startTime).Seconds()
		tags := []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "17864"}
		if err := postMetricToDatadog(ctx, apiClient, "bridging.success", completionTimeSec, tags); err != nil {
			logger.Error("failed to post metric", "error", err)
		}

		// Sleep for random interval between 0 and 5 seconds
		time.Sleep(time.Duration(mathrand.Intn(6)) * time.Second)

		// Bridge back same amount minus 0.009 ETH for fees
		pZZNineEth := big.NewInt(9 * params.Ether / 1000)
		amountBack := new(big.Int).Sub(randWeiValue, pZZNineEth)

		// Create and start the transfer back to L1 with the same amount
		tL1, err := transfer.NewTransferToL1(
			logger.With("component", "l1_transfer"),
			amountBack,
			transferAddr,
			privateKey,
			settlementRPCUrl,
			l1RPCUrl,
			l1ContractAddr,
			settlementContractAddr,
		)
		if err != nil {
			tags := []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "39999"}
			if err := postMetricToDatadog(ctx, apiClient, "bridging.failure", 0, tags); err != nil {
				logger.Error("failed to post metric", "error", err)
			}
			time.Sleep(time.Minute)
			continue
		}
		startTime = time.Now()
		err = tL1.Start(ctx)
		if err != nil {
			tags := []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "39999"}
			if err := postMetricToDatadog(ctx, apiClient, "bridging.failure", time.Since(startTime).Seconds(), tags); err != nil {
				logger.Error("failed to post metric", "error", err)
			}
			time.Sleep(time.Minute)
			continue
		}
		completionTimeSec = time.Since(startTime).Seconds()
		tags = []string{"environment:bridge_test", "account_addr:" + transferAddressString, "to_chain_id:" + "39999"}
		if err := postMetricToDatadog(ctx, apiClient, "bridging.success", completionTimeSec, tags); err != nil {
			logger.Error("failed to post metric", "error", err)
		}

		// Sleep for random interval between 0 and 5 seconds
		time.Sleep(time.Duration(mathrand.Intn(6)) * time.Second)
	}
}

func postMetricToDatadog(ctx context.Context, client *datadog.APIClient, metricName string, value float64, tags []string) error {
	now := time.Now().Unix()
	point := datadog.MetricPoint{
		Timestamp: datadog.PtrInt64(now),
		Value:     datadog.PtrFloat64(value),
	}
	series := datadog.MetricSeries{
		Metric: metricName,
		Type:   datadog.METRICINTAKETYPE_GAUGE.Ptr(),
		Points: []datadog.MetricPoint{point},
		Tags:   tags,
	}
	payload := datadog.MetricPayload{
		Series: []datadog.MetricSeries{series},
	}
	_, _, err := client.MetricsApi.SubmitMetrics(ctx, payload)
	return err
}
