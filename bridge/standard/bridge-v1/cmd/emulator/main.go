package main

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"path"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/primev/mev-commit/bridge/standard/bridge-v1/pkg/transfer"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/primev/mev-commit/x/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

var (
	optionKeystorePathPassword = &cli.StringSliceFlag{
		Name:    "keystore-path-password",
		Usage:   "Path to the keystore file and password in the format path:password",
		EnvVars: []string{"EMULATOR_KEYSTORE_PATH_PASSWORD"},
		Action: func(c *cli.Context, keystores []string) error {
			for _, kp := range keystores {
				parts := strings.Split(kp, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid keystore-path-password format: %s", kp)
				}
			}
			return nil
		},
	}

	optionL1RPCURL = &cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL of the L1 RPC server",
		EnvVars: []string{"EMULATOR_L1_RPC_URL"},
	}

	optionSettlementRPCEndpoint = &cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL of the settlement RPC server",
		EnvVars: []string{"EMULATOR_SETTLEMENT_RPC_URL"},
	}

	optionL1GatewayContractAddr = &cli.StringFlag{
		Name:    "l1-gateway-contract-addr",
		Usage:   "L1 gateway contract address",
		EnvVars: []string{"EMULATOR_L1_GATEWAY_CONTRACT_ADDR"},
	}

	optionSettlementGatewayContractAddr = &cli.StringFlag{
		Name:    "settlement-gateway-contract-addr",
		Usage:   "Settlement gateway contract address",
		EnvVars: []string{"EMULATOR_SETTLEMENT_GATEWAY_CONTRACT_ADDR"},
	}

	optionHTTPPort = &cli.IntFlag{
		Name:    "http-port",
		Usage:   "HTTP port to listen on",
		EnvVars: []string{"EMULATOR_HTTP_PORT"},
		Value:   8080,
	}

	optionLogFmt = &cli.StringFlag{
		Name:    "log-fmt",
		Usage:   "log format to use, options are 'text' or 'json'",
		EnvVars: []string{"EMULATOR_LOG_FMT"},
		Value:   "text",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"text", "json"}, s) {
				return fmt.Errorf("invalid log-fmt, expecting 'text' or 'json'")
			}
			return nil
		},
	}

	optionLogLevel = &cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level to use, options are 'debug', 'info', 'warn', 'error'",
		EnvVars: []string{"EMULATOR_LOG_LEVEL"},
		Value:   "info",
		Action: func(ctx *cli.Context, s string) error {
			if !slices.Contains([]string{"debug", "info", "warn", "error"}, s) {
				return fmt.Errorf("invalid log-level, expecting 'debug', 'info', 'warn', 'error'")
			}
			return nil
		},
	}

	optionLogTags = &cli.StringFlag{
		Name:    "log-tags",
		Usage:   "log tags is a comma-separated list of <name:value> pairs that will be inserted into each log line",
		EnvVars: []string{"EMULATOR_LOG_TAGS"},
		Action: func(ctx *cli.Context, s string) error {
			for i, p := range strings.Split(s, ",") {
				if len(strings.Split(p, ":")) != 2 {
					return fmt.Errorf("invalid log-tags at index %d, expecting <name:value>", i)
				}
			}
			return nil
		},
	}
)

var (
	bridgeSuccessDurations = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bridge_success_duration",
		Help: "Duration of successful bridge transactions",
	}, []string{"account", "direction"})

	bridgeFailureDurations = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bridge_failure_duration",
		Help: "Duration of failed bridge transactions",
	}, []string{"account", "direction"})
)

func main() {
	app := &cli.App{
		Name:  "bridge-emulator",
		Usage: "Issue random bridge transactions between L1 and settlement chain",
		Flags: []cli.Flag{
			optionKeystorePathPassword,
			optionL1RPCURL,
			optionSettlementRPCEndpoint,
			optionL1GatewayContractAddr,
			optionSettlementGatewayContractAddr,
			optionLogFmt,
			optionLogLevel,
			optionLogTags,
		},
		Action: func(c *cli.Context) error {
			logger, err := util.NewLogger(
				c.String(optionLogLevel.Name),
				c.String(optionLogFmt.Name),
				c.String(optionLogTags.Name),
				c.App.Writer,
			)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			type txtor struct {
				pk   *ecdsa.PrivateKey
				addr common.Address
			}

			txtors := make([]*txtor, 0)
			for _, kp := range c.StringSlice(optionKeystorePathPassword.Name) {
				parts := strings.Split(kp, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid keystore-path-password format: %s", kp)
				}
				keySigner, err := keysigner.NewKeystoreSigner(path.Dir(parts[0]), parts[1])
				if err != nil {
					return fmt.Errorf("failed creating key signer: %w", err)
				}
				pk, err := keySigner.GetPrivateKey()
				if err != nil {
					return fmt.Errorf("failed to get private key: %w", err)
				}
				txtors = append(txtors, &txtor{
					pk:   pk,
					addr: keySigner.GetAddress(),
				})
			}

			registry := prometheus.NewRegistry()
			registry.MustRegister(bridgeSuccessDurations, bridgeFailureDurations)

			router := http.NewServeMux()
			router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

			server := &http.Server{
				Addr:    fmt.Sprintf(":%d", c.Int(optionHTTPPort.Name)),
				Handler: router,
			}

			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("failed to start server", "error", err)
				}
			}()

			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

			for {
				select {
				case <-sigc:
					return server.Shutdown(c.Context)
				default:
				}

				rand.Shuffle(len(txtors), func(i, j int) {
					txtors[i], txtors[j] = txtors[j], txtors[i]
				})

				// Generate a random amount of wei in [0.01, 10] ETH
				randWeiValue := big.NewInt(rand.Int64N(int64(params.Ether)))
				if randWeiValue.Cmp(big.NewInt(params.Ether/1000)) < 0 {
					// Enforce minimum value of 0.01 ETH
					randWeiValue = big.NewInt(params.Ether / 1000)
				}
				randWeiValue = new(big.Int).Mul(randWeiValue, big.NewInt(10))

				// Create and start the transfer to the settlement chain
				tSettlement, err := transfer.NewTransferToSettlement(
					logger.With("component", "settlement_transfer"),
					randWeiValue,
					txtors[0].addr,
					txtors[0].pk,
					c.String(optionSettlementRPCEndpoint.Name),
					c.String(optionL1RPCURL.Name),
					common.HexToAddress(c.String(optionL1GatewayContractAddr.Name)),
					common.HexToAddress(c.String(optionSettlementGatewayContractAddr.Name)),
				)
				if err != nil {
					logger.Error("failed to create transfer to settlement", "error", err)
					time.Sleep(time.Minute)
					continue
				}
				startTime := time.Now()
				err = tSettlement.Start(c.Context)
				if err != nil {
					logger.Error("failed to start transfer to settlement", "error", err)
					bridgeFailureDurations.WithLabelValues(
						txtors[0].addr.String(),
						"L1->Settlement",
					).Set(time.Since(startTime).Seconds())
					time.Sleep(time.Minute)
					continue
				}
				completionTimeSec := time.Since(startTime).Seconds()
				logger.Info("completed settlement transfer",
					"time", completionTimeSec,
					"amount", randWeiValue.String(),
					"address", txtors[0].addr.String(),
				)
				bridgeSuccessDurations.WithLabelValues(
					txtors[0].addr.String(),
					"L1->Settlement",
				).Set(completionTimeSec)

				// Sleep for random interval between 0 and 5 seconds
				time.Sleep(time.Duration(rand.IntN(6)) * time.Second)

				// Bridge back same amount minus 0.009 ETH for fees
				pZZNineEth := big.NewInt(9 * params.Ether / 1000)
				amountBack := new(big.Int).Sub(randWeiValue, pZZNineEth)

				// Create and start the transfer back to L1 with the same amount
				tL1, err := transfer.NewTransferToL1(
					logger.With("component", "l1_transfer"),
					amountBack,
					txtors[0].addr,
					txtors[0].pk,
					c.String(optionSettlementRPCEndpoint.Name),
					c.String(optionL1RPCURL.Name),
					common.HexToAddress(c.String(optionL1GatewayContractAddr.Name)),
					common.HexToAddress(c.String(optionSettlementGatewayContractAddr.Name)),
				)
				if err != nil {
					logger.Error("failed to create transfer to L1", "error", err)
					time.Sleep(time.Minute)
					continue
				}
				startTime = time.Now()
				err = tL1.Start(c.Context)
				if err != nil {
					logger.Error("failed to start transfer to L1", "error", err)
					bridgeFailureDurations.WithLabelValues(
						txtors[0].addr.String(),
						"Settlement->L1",
					).Set(time.Since(startTime).Seconds())
					time.Sleep(time.Minute)
					continue
				}
				completionTimeSec = time.Since(startTime).Seconds()
				logger.Info("completed L1 transfer",
					"time", completionTimeSec,
					"amount", amountBack.String(),
					"address", txtors[0].addr.String(),
				)
				bridgeSuccessDurations.WithLabelValues(
					txtors[0].addr.String(),
					"Settlement->L1",
				).Set(completionTimeSec)

				// Sleep for random interval between 0 and 5 seconds
				time.Sleep(time.Duration(rand.IntN(6)) * time.Second)
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
