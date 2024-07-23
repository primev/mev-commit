package registration

import (
	"fmt"
	"log/slog"

	eigenclitypes "github.com/Layr-Labs/eigenlayer-cli/pkg/types"
	eigencliutils "github.com/Layr-Labs/eigenlayer-cli/pkg/utils"
	dm "github.com/Layr-Labs/eigensdk-go/contracts/bindings/DelegationManager"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	avs "github.com/primev/mev-commit/contracts-abi/clients/MevCommitAVS"
	"github.com/primev/mev-commit/x/contracts/transactor"
	"github.com/primev/mev-commit/x/contracts/txmonitor"
	ks "github.com/primev/mev-commit/x/keysigner"
	"github.com/urfave/cli/v2"
)

type avsContractWithSesh struct {
	avs  *avs.Mevcommitavs
	sesh *avs.MevcommitavsTransactorSession
}

// TODO: re-eval if all these fields need to be stored
type Command struct {
	Logger           *slog.Logger
	OperatorConfig   *eigenclitypes.OperatorConfig
	KeystorePassword string
	signer           *ks.KeystoreSigner
	ethClient        *ethclient.Client

	AVSAddress          string
	avsContractWithSesh *avsContractWithSesh

	DelegationManagerAddress string
}

func (c *Command) initialize(ctx *cli.Context) error {
	ethClient, err := ethclient.Dial(c.OperatorConfig.EthRPCUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	c.ethClient = ethClient

	chainID, err := ethClient.ChainID(ctx.Context)
	if err != nil {
		c.Logger.Error("failed to get chain ID", "error", err)
		return err
	}

	if chainID.Cmp(&c.OperatorConfig.ChainId) != 0 {
		return fmt.Errorf("chain ID from rpc url doesn't match operator config: %s != %s",
			chainID.String(), c.OperatorConfig.ChainId.String())
	}
	c.Logger.Info("Chain ID", "chainID", chainID)

	if c.KeystorePassword == "" {
		prompter := eigencliutils.NewPrompter()
		keystorePwd, err := prompter.InputHiddenString(
			fmt.Sprintf("Enter password to decrypt ecdsa keystore for %s:", c.OperatorConfig.PrivateKeyStorePath), "",
			func(string) error {
				return nil
			},
		)
		if err != nil {
			return fmt.Errorf("failed to read keystore password: %w", err)
		}
		c.KeystorePassword = keystorePwd
	}

	signer, err := ks.NewKeystoreSigner(c.OperatorConfig.PrivateKeyStorePath, c.KeystorePassword)
	if err != nil {
		return fmt.Errorf("failed to create keystore signer: %w", err)
	}
	c.signer = signer

	monitor := txmonitor.New(
		signer.GetAddress(),
		ethClient,
		txmonitor.NewEVMHelperWithLogger(ethClient.Client(), c.Logger),
		nil, // TOOD: re-eval if you need saver/store
		c.Logger.With("component", "txmonitor"),
		1, // TODO: re-eval max pending
	)
	transactor := transactor.NewTransactor(
		ethClient,
		monitor,
	)

	avsAddress := common.HexToAddress(c.AVSAddress)
	mevCommitAVS, err := avs.NewMevcommitavs(avsAddress, transactor)
	if err != nil {
		return fmt.Errorf("failed to create mev-commit avs: %w", err)
	}

	tOpts, err := c.signer.GetAuth(chainID)
	if err != nil {
		c.Logger.Error("failed to get auth", "error", err)
		return err
	}

	sesh := &avs.MevcommitavsTransactorSession{
		Contract:     &mevCommitAVS.MevcommitavsTransactor,
		TransactOpts: *tOpts,
	}

	c.avsContractWithSesh = &avsContractWithSesh{
		avs:  mevCommitAVS,
		sesh: sesh,
	}

	delegationManager, err := dm.NewContractDelegationManagerCaller(common.HexToAddress(""), nil) // TODO
	if err != nil {
		return fmt.Errorf("failed to create delegation manager: %w", err)
	}
	fmt.Println(delegationManager)

	return nil
}

func (c *Command) RegisterOperator(ctx *cli.Context) error {

	c.Logger.Info("Registering operator...")
	c.initialize(ctx)

	operatorRegInfo, err := c.avsContractWithSesh.avs.GetOperatorRegInfo(
		&bind.CallOpts{Context: ctx.Context}, c.signer.GetAddress())
	if err != nil {
		return fmt.Errorf("failed to get operator reg info: %w", err)
	}
	if operatorRegInfo.Exists {
		return fmt.Errorf("signing operator already registered")
	}

	// TODO: also query EL's delegation manager

	// TODO: generate actual sig
	operatorSig := avs.ISignatureUtilsSignatureWithSaltAndExpiry{}

	tx, err := c.avsContractWithSesh.sesh.RegisterOperator(operatorSig)
	if err != nil {
		return fmt.Errorf("failed to register operator: %w", err)
	}

	c.Logger.Info("waiting for tx to be mined", "txHash", tx.Hash().Hex(), "nonce", tx.Nonce())
	rec, err := bind.WaitMined(ctx.Context, c.ethClient, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for tx to be mined: %w", err)
	} else if rec.Status != ethtypes.ReceiptStatusSuccessful {
		return fmt.Errorf("receipt status unsuccessful: %d", rec.Status)
	}
	// TODO: Confirm we don't need to set gas params manually anymore?

	// TODO: Determine if we want to support fee bump, cancelling, etc. Do some testing here..
	return nil
}

func (c *Command) RequestOperatorDeregistration(ctx *cli.Context) error {
	c.Logger.Info("Requesting operator deregistration...")

	client, err := ethclient.Dial(ctx.String("eth-node-url"))
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	c.Logger.Info("Client", "client", client)

	operatorAddress := ctx.String("operator-address")

	// Add your request deregistration logic here, e.g., sending a transaction to the Ethereum network
	c.Logger.Info("Operator deregistration requested", "address", operatorAddress)
	return nil
}

func (c *Command) DeregisterOperator(ctx *cli.Context) error {
	c.Logger.Info("Deregistering operator...")

	client, err := ethclient.Dial(ctx.String("eth-node-url"))
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	c.Logger.Info("Client", "client", client)

	operatorAddress := ctx.String("operator-address")

	// Add your deregistration logic here, e.g., sending a transaction to the Ethereum network
	c.Logger.Info("Operator deregistered", "address", operatorAddress)
	return nil
}
