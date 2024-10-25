package main

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primev/mev-commit/bridge/standard/pkg/transfer"
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
		Name:     "l1-rpc-url",
		Usage:    "URL for L1 RPC",
		EnvVars:  []string{"L1_RPC_URL"},
		Required: true,
	}
	optionSettlementRPCUrl = &cli.StringFlag{
		Name:     "settlement-rpc-url",
		Usage:    "URL for settlement RPC",
		EnvVars:  []string{"SETTLEMENT_RPC_URL"},
		Required: true,
	}
	optionL1ContractAddr = &cli.StringFlag{
		Name:     "l1-contract-addr",
		Usage:    "address of the L1 gateway contract",
		EnvVars:  []string{"L1_CONTRACT_ADDR"},
		Required: true,
	}
	optionSettlementContractAddr = &cli.StringFlag{
		Name:     "settlement-contract-addr",
		Usage:    "address of the settlement gateway contract",
		EnvVars:  []string{"SETTLEMENT_CONTRACT_ADDR"},
		Required: true,
	}
)

func main() {
	app := &cli.App{
		Name:  "bridge-cli",
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
	spinner, err := createSpinner()
	if err != nil {
		return fmt.Errorf("failed to create spinner: %w", err)
	}
	spinner.Message("Starting transfer from L1 to Settlement chain")
	if err := spinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}
	statusC := t.Do(c.Context)
	for status := range statusC {
		if status.Error != nil {
			spinner.StopFailMessage(fmt.Sprintf("%s: Error: %s", status.Message, status.Error))
			return fmt.Errorf("failed to start transfer to settlement: %w", status.Error)
		}
		spinner.Message(status.Message)
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
	spinner, err := createSpinner()
	if err != nil {
		return fmt.Errorf("failed to create spinner: %w", err)
	}
	spinner.Message("Starting transfer from Settlement chain to L1")
	if err := spinner.Start(); err != nil {
		return fmt.Errorf("failed to start spinner: %w", err)
	}
	statusC := t.Do(c.Context)
	for status := range statusC {
		if status.Error != nil {
			spinner.StopFailMessage(fmt.Sprintf("%s: Error: %s", status.Message, status.Error))
			return fmt.Errorf("failed to start transfer to L1: %w", status.Error)
		}
		spinner.Message(status.Message)
	}
	return spinner.Stop()
}

func createSpinner() (*yacspin.Spinner, error) {
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
		StopMessage:       "done",
		StopFailCharacter: "✗",
		StopFailColors:    []string{"fgRed"},
		StopFailMessage:   "failed",
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to make spinner from struct: %w", err)
	}

	return s, nil
}
