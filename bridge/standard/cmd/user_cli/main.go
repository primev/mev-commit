package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/bridge/standard/pkg/transfer"
	"github.com/primev/mev-commit/contracts-abi/config"
	"github.com/primev/mev-commit/x/keysigner"
	"github.com/theckman/yacspin"
	"github.com/urfave/cli/v2"
)

var (
	optionPrivKey = &cli.StringFlag{
		Name:    "account-key",
		Usage:   "private key of the account to use for signing transactions in hex encoding",
		EnvVars: []string{"ACCOUNT_KEY"},
	}
	optionKeystorePath = &cli.StringFlag{
		Name:    "keystore-path",
		Usage:   "path to keystore location",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PATH"},
	}
	optionKeystorePassword = &cli.StringFlag{
		Name:    "keystore-password",
		Usage:   "use to access keystore",
		EnvVars: []string{"MEV_ORACLE_KEYSTORE_PASSWORD"},
	}
	optionAmount = &cli.StringFlag{
		Name:     "amount",
		Usage:    "amount of ether to bridge in wei",
		EnvVars:  []string{"AMOUNT"},
		Required: true,
	}
	optionDestAddr = &cli.StringFlag{
		Name:     "dest-addr",
		Usage:    "destination address on the mev-commit (settlement) chain",
		EnvVars:  []string{"DEST_ADDR"},
		Required: true,
	}
	optionL1RPCUrl = &cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL for L1 RPC",
		EnvVars: []string{"L1_RPC_URL"},
		Value:   "https://ethereum-holesky-rpc.publicnode.com",
	}
	optionSettlementRPCUrl = &cli.StringFlag{
		Name:    "settlement-rpc-url",
		Usage:   "URL for settlement RPC",
		EnvVars: []string{"SETTLEMENT_RPC_URL"},
		Value:   "https://chainrpc.testnet.mev-commit.xyz",
	}
	optionL1ContractAddr = &cli.StringFlag{
		Name:    "l1-contract-addr",
		Usage:   "address of the L1 gateway contract",
		EnvVars: []string{"L1_CONTRACT_ADDR"},
		Value:   config.EthereumContracts.L1Gateway,
	}
	optionSettlementContractAddr = &cli.StringFlag{
		Name:    "settlement-contract-addr",
		Usage:   "address of the settlement gateway contract",
		EnvVars: []string{"SETTLEMENT_CONTRACT_ADDR"},
		Value:   config.MainnetContracts.SettlementGateway,
	}
	optionSilent = &cli.BoolFlag{
		Name:    "silent",
		Usage:   "disable spinner",
		EnvVars: []string{"SILENT"},
		Value:   false,
	}
)

func main() {
	app := &cli.App{
		Name:  "mev-commit-bridge-user-cli",
		Usage: "CLI for interacting with a custom bridge between L1 and the mev-commit (settlement) chain",
		Commands: []*cli.Command{
			{
				Name:  "bridge-to-settlement",
				Usage: "Submit a transaction to bridge ether to the settlement chain",
				Flags: []cli.Flag{
					optionPrivKey,
					optionKeystorePath,
					optionKeystorePassword,
					optionAmount,
					optionDestAddr,
					optionL1RPCUrl,
					optionSettlementRPCUrl,
					optionL1ContractAddr,
					optionSettlementContractAddr,
					optionSilent,
				},
				Action: bridgeToSettlement,
			},
			{
				Name:  "bridge-to-l1",
				Usage: "Submit a transaction to bridge ether back to L1",
				Flags: []cli.Flag{
					optionPrivKey,
					optionKeystorePath,
					optionKeystorePassword,
					optionAmount,
					optionDestAddr,
					optionL1RPCUrl,
					optionSettlementRPCUrl,
					optionL1ContractAddr,
					optionSettlementContractAddr,
					optionSilent,
				},
				Action: bridgeToL1,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getSigner(c *cli.Context) (keysigner.KeySigner, error) {
	if c.String("account-key") != "" {
		return keysigner.NewPrivateKeySignerFromHex(c.String("account-key"))
	}
	if c.String("keystore-path") != "" && c.String("keystore-password") != "" {
		return keysigner.NewKeystoreSigner(c.String("keystore-path"), c.String("keystore-password"))
	}
	return nil, errors.New("either account-key or keystore-path and keystore-password must be provided")
}

func bridgeToSettlement(c *cli.Context) error {
	amount, ok := big.NewInt(0).SetString(c.String("amount"), 10)
	if !ok {
		return errors.New("failed to parse amount")
	}
	signer, err := getSigner(c)
	if err != nil {
		return fmt.Errorf("failed to get signer: %w", err)
	}
	t, err := transfer.NewTransferToSettlement(
		amount,
		common.HexToAddress(c.String("dest-addr")),
		signer,
		c.String("settlement-rpc-url"),
		c.String("l1-rpc-url"),
		common.HexToAddress(c.String("l1-contract-addr")),
		common.HexToAddress(c.String("settlement-contract-addr")),
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer to settlement: %w", err)
	}
	silent := c.Bool(optionSilent.Name)
	spinner, err := createSpinner("Starting transfer from L1 to Settlement chain", silent)
	if err != nil {
		return fmt.Errorf("failed to create spinner: %w", err)
	}
	if err := spinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}

	ctx, cancel := context.WithTimeout(c.Context, 30*time.Minute)
	defer cancel()

	statusC := t.Do(ctx)
	for status := range statusC {
		if status.Error != nil {
			spinner.StopFailMessage(fmt.Sprintf("%s: Error: %s", status.Message, status.Error))
			return fmt.Errorf("failed to start transfer to settlement: %w", status.Error)
		}
		if err := spinner.Stop(); err != nil {
			return fmt.Errorf("failed to stop spinner: %w", err)
		}
		spinner, err = createSpinner(status.Message, silent)
		if err != nil {
			return fmt.Errorf("failed to create spinner: %w", err)
		}
		if err := spinner.Start(); err != nil {
			return fmt.Errorf("failed to start spinner: %w", err)
		}
	}
	return spinner.Stop()
}

func bridgeToL1(c *cli.Context) error {
	amount, ok := big.NewInt(0).SetString(c.String("amount"), 10)
	if !ok {
		return errors.New("failed to parse amount")
	}
	signer, err := getSigner(c)
	if err != nil {
		return fmt.Errorf("failed to get signer: %w", err)
	}
	t, err := transfer.NewTransferToL1(
		amount,
		common.HexToAddress(c.String("dest-addr")),
		signer,
		c.String("settlement-rpc-url"),
		c.String("l1-rpc-url"),
		common.HexToAddress(c.String("l1-contract-addr")),
		common.HexToAddress(c.String("settlement-contract-addr")),
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer to L1: %w", err)
	}
	silent := c.Bool(optionSilent.Name)
	spinner, err := createSpinner("Starting transfer from Settlement chain to L1", silent)
	if err != nil {
		return fmt.Errorf("failed to create spinner: %w", err)
	}
	if err := spinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}

	ctx, cancel := context.WithTimeout(c.Context, 30*time.Minute)
	defer cancel()

	statusC := t.Do(ctx)
	for status := range statusC {
		if status.Error != nil {
			spinner.StopFailMessage(fmt.Sprintf("%s: Error: %s", status.Message, status.Error))
			return fmt.Errorf("failed to start transfer to L1: %w", status.Error)
		}
		if err := spinner.Stop(); err != nil {
			return fmt.Errorf("failed to stop spinner: %w", err)
		}
		spinner, err = createSpinner(status.Message, silent)
		if err != nil {
			return fmt.Errorf("failed to create spinner: %w", err)
		}
		if err := spinner.Start(); err != nil {
			return fmt.Errorf("failed to start spinner: %w", err)
		}
	}
	return spinner.Stop()
}

func createSpinner(msg string, silent bool) (*yacspin.Spinner, error) {
	// build the configuration, each field is documented
	cfg := yacspin.Config{
		Frequency:         100 * time.Millisecond,
		CharSet:           yacspin.CharSets[11],
		Suffix:            " ", // puts a least one space between the animating spinner and the Message
		SuffixAutoColon:   true,
		ColorAll:          true,
		Colors:            []string{"fgYellow"},
		StopCharacter:     "✓",
		StopColors:        []string{"fgGreen"},
		Message:           msg,
		StopMessage:       msg,
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	if silent {
		cfg.Writer = io.Discard
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to make spinner from struct: %w", err)
	}

	return s, nil
}
