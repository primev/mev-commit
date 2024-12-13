package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cloudflare/circl/sign/bls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli/v2"
)

var (
	optionPrivateKey = &cli.StringFlag{
		Name:  "private-key",
		Usage: "BLS private key as hex encoded string (with optional 0x prefix). Must be a valid hex string representing a BLS private key.",
	}

	optionEthAddress = &cli.StringFlag{
		Name:  "eth-address",
		Usage: "Ethereum address as hex string (with optional 0x prefix)",
	}
)

func main() {
	app := &cli.App{
		Name:  "bls-signature-creator",
		Usage: "Create BLS signatures",
		Flags: []cli.Flag{
			optionPrivateKey,
			optionEthAddress,
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	blsPrivKeyHex := c.String("private-key")
	ethAddress := c.String("eth-address")

	if blsPrivKeyHex == "" {
		return fmt.Errorf("--private-key flag is required")
	}

	if ethAddress == "" {
		return fmt.Errorf("--eth-address flag is required")
	}

	// Strip 0x prefix if present
	blsPrivKeyHex = strings.TrimPrefix(blsPrivKeyHex, "0x")

	// Validate hex string
	privKeyBytes, err := hexutil.Decode("0x" + blsPrivKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex string: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create BLS signature
	hashedMessage := crypto.Keccak256(common.HexToAddress(ethAddress).Bytes())
	privateKey := new(bls.PrivateKey[bls.G1])
	if err := privateKey.UnmarshalBinary(privKeyBytes); err != nil {
		logger.Error("failed to unmarshal private key", "error", err)
		return err
	}

	publicKey := privateKey.PublicKey()
	signature := bls.Sign(privateKey, hashedMessage)

	// Verify the signature
	if !bls.Verify(publicKey, hashedMessage, signature) {
		logger.Error("failed to verify generated BLS signature")
		return fmt.Errorf("failed to verify generated BLS signature")
	}

	pubkeyb, err := publicKey.MarshalBinary()
	if err != nil {
		logger.Error("failed to marshal public key", "error", err)
		return err
	}

	logger.Info("generated BLS signature",
		"eth_address", ethAddress,
		"public_key", hex.EncodeToString(pubkeyb),
		"signature", hex.EncodeToString(signature))

	return nil
}
