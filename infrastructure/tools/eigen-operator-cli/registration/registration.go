package registration

import (
	"fmt"
	"log/slog"

	eigenclitypes "github.com/Layr-Labs/eigenlayer-cli/pkg/types"
	eigenecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

type Command struct {
	Logger         *slog.Logger
	OperatorConfig eigenclitypes.OperatorConfig
}

func (c *Command) RegisterOperator(ctx *cli.Context) error {
	c.Logger.Info("Registering operator...")

	password := ctx.String("password")

	privKey, err := eigenecdsa.ReadKey(c.OperatorConfig.PrivateKeyStorePath, password)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	client, err := ethclient.Dial(ctx.String("eth-node-url"))
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	c.Logger.Info("Client", "client", client)

	operatorAddress := ctx.String("operator-address")
	signature := ctx.String("signature")

	// Add your registration logic here, e.g., sending a transaction to the Ethereum network
	c.Logger.Info("Operator registered", "address", operatorAddress, "signature", signature)
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
