package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"

	vanillaregistry "github.com/primev/mev-commit/contracts-abi/clients/VanillaRegistry"
	"github.com/primev/mev-commit/x/keysigner"
)

var (
	optionKeystorePath = &cli.StringFlag{
		Name:     "keystore-dir",
		Usage:    "directory where keystore file is stored",
		EnvVars:  []string{"KEYSTORE_DIR"},
		Required: true,
	}
	optionKeystorePassword = &cli.StringFlag{
		Name:     "keystore-password",
		Usage:    "use to access keystore",
		EnvVars:  []string{"KEYSTORE_PASSWORD"},
		Required: true,
	}
	optionL1RPCURL = &cli.StringFlag{
		Name:     "l1-rpc-url",
		Usage:    "URL of the L1 RPC server",
		EnvVars:  []string{"L1_RPC_URL"},
		Required: true,
	}
	optionPubkeyFilePath = &cli.StringFlag{
		Name:     "pubkey-file-path",
		Usage:    "path to the file containing the public keys",
		EnvVars:  []string{"PUBKEY_FILE_PATH"},
		Required: true,
	}
	optionVanillaRegistryAddress = &cli.StringFlag{
		Name:    "vanilla-registry-address",
		Usage:   "address of the vanilla registry contract",
		EnvVars: []string{"VANILLA_REGISTRY_ADDRESS"},
		Value:   "0x47afdcB2B089C16CEe354811EA1Bbe0DB7c335E9",
	}
)

func main() {

	flags := []cli.Flag{
		optionKeystorePath,
		optionKeystorePassword,
		optionL1RPCURL,
		optionPubkeyFilePath,
	}

	app := &cli.App{
		Name:  "vanill-stake",
		Usage: "Stake validators programmatically with the mev-commit vanilla registry",
		Flags: flags,
		Action: func(c *cli.Context) error {
			return stakeVanilla(c)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func stakeVanilla(c *cli.Context) error {
	keystoreDir := c.String(optionKeystorePath.Name)
	keystorePassword := c.String(optionKeystorePassword.Name)
	l1RpcUrl := c.String(optionL1RPCURL.Name)
	pubkeyFilePath := c.String(optionPubkeyFilePath.Name)
	vanillaRegistryAddress := c.String(optionVanillaRegistryAddress.Name)

	signer, err := keysigner.NewKeystoreSigner(keystoreDir, keystorePassword)
	if err != nil {
		return fmt.Errorf("failed to create signer: %w", err)
	}

	client, err := ethclient.Dial(l1RpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	vrt, err := vanillaregistry.NewVanillaregistryTransactor(common.HexToAddress(vanillaRegistryAddress), client)
	if err != nil {
		return fmt.Errorf("failed to create Vanilla Registry transactor: %w", err)
	}

	pksAsBytes, err := readBLSPublicKeysFromFile(pubkeyFilePath)
	if err != nil {
		return fmt.Errorf("failed to read public keys from file: %w", err)
	}

	batchSize := 20
	type Batch struct {
		pubKeys [][]byte
	}
	batches := make([]Batch, 0)
	for i := 0; i < len(pksAsBytes); i += batchSize {
		end := i + batchSize
		if end > len(pksAsBytes) {
			end = len(pksAsBytes)
		}
		batches = append(batches, Batch{pubKeys: pksAsBytes[i:end]})
	}

	for idx, batch := range batches {

		opts, err := signer.GetAuth(chainID)
		if err != nil {
			return fmt.Errorf("failed to create transact opts: %w", err)
		}

		amountPerValidator := new(big.Int)
		amountPerValidator.SetString("100000000000000", 10) // TODO: make configurable
		totalAmount := new(big.Int).Mul(amountPerValidator, big.NewInt(int64(batchSize)))
		opts.Value = totalAmount

		tx, err := vrt.Stake(opts, batch.pubKeys)
		if err != nil {
			return fmt.Errorf("failed to stake: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			return fmt.Errorf("failed to wait for stake tx to be mined: %w", err)
		}

		if receipt.Status != types.ReceiptStatusSuccessful {
			return fmt.Errorf("stake tx included, but failed")
		}

		fmt.Println("-------------------")
		fmt.Printf("Batch %d completed\n", idx+1)
		fmt.Println("-------------------")
	}
	fmt.Println("All staking batches completed!")
	return nil
}

func readBLSPublicKeysFromFile(filePath string) ([][]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keys [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := scanner.Text()
		if len(key) > 2 && key[:2] == "0x" {
			key = key[2:]
		}
		keyBytes := common.Hex2Bytes(key)
		keys = append(keys, keyBytes)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}
