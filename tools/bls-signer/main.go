package main

import (
	"encoding/hex"
	"fmt"
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
		Name:     "private-key",
		Usage:    "BLS private key as hex encoded string (with optional 0x prefix). Must be a valid hex string representing a BLS private key.",
		Required: true,
	}

	optionPayload = &cli.StringFlag{
		Name:     "payload",
		Usage:    "Payload as hex string (with optional 0x prefix)",
		Required: true,
	}
)

func main() {
	app := &cli.App{
		Name:  "bls-signer",
		Usage: "Create BLS signatures",
		Flags: []cli.Flag{
			optionPrivateKey,
			optionPayload,
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
	payload := c.String("payload")

	// Strip 0x prefix if present
	blsPrivKeyHex = strings.TrimPrefix(blsPrivKeyHex, "0x")

	// Validate hex string
	privKeyBytes, err := hexutil.Decode("0x" + blsPrivKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex string: %v", err)
	}

	// Create BLS signature
	hashedMessage := crypto.Keccak256(common.HexToAddress(payload).Bytes())
	privateKey := new(bls.PrivateKey[bls.G1])
	if err := privateKey.UnmarshalBinary(privKeyBytes); err != nil {
		//nolint:errcheck
		fmt.Fprintf(c.App.Writer, "Failed to unmarshal private key: %v\n", err)
		return err
	}

	publicKey := privateKey.PublicKey()
	signature := bls.Sign(privateKey, hashedMessage)

	// Verify the signature
	if !bls.Verify(publicKey, hashedMessage, signature) {
		//nolint:errcheck
		fmt.Fprintln(c.App.Writer, "Failed to verify generated BLS signature")
		return fmt.Errorf("failed to verify generated BLS signature")
	}

	pubkeyb, err := publicKey.MarshalBinary()
	if err != nil {
		//nolint:errcheck
		fmt.Fprintf(c.App.Writer, "Failed to marshal public key: %v\n", err)
		return err
	}

	//nolint:errcheck
	fmt.Fprintf(c.App.Writer, "Generated BLS signature:\n")
	//nolint:errcheck
	fmt.Fprintf(c.App.Writer, "Payload: %s\n", payload)
	//nolint:errcheck
	fmt.Fprintf(c.App.Writer, "Public key: %s\n", hex.EncodeToString(pubkeyb))
	//nolint:errcheck
	fmt.Fprintf(c.App.Writer, "Signature: %s\n", hex.EncodeToString(signature))

	return nil
}
