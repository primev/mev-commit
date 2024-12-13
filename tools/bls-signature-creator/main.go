package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/circl/sign/bls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/urfave/cli/v2"
)

type Topology struct {
	Self struct {
		EthereumAddress string `json:"Ethereum Address"`
	} `json:"self"`
}

var (
	optionPrivateKey = &cli.StringFlag{
		Name:  "private-key",
		Usage: "BLS private key as hex encoded string (with optional 0x prefix). Must be a valid hex string representing a BLS private key.",
	}

	optionServerAddr = &cli.StringFlag{
		Name:  "server",
		Usage: "Provider server address",
		Value: "localhost:8545",
	}
)

func main() {
	app := &cli.App{
		Name:  "bls-signature-creator",
		Usage: "Create BLS signatures",
		Flags: []cli.Flag{
			optionPrivateKey,
			optionServerAddr,
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
	serverAddr := c.String("server")

	if blsPrivKeyHex == "" {
		return fmt.Errorf("--private-key flag is required")
	}

	// Strip 0x prefix if present
	blsPrivKeyHex = strings.TrimPrefix(blsPrivKeyHex, "0x")

	// Validate hex string
	privKeyBytes, err := hexutil.Decode("0x" + blsPrivKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex string: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get topology from debug endpoint
	resp, err := http.Get(fmt.Sprintf("http://%s/v1/debug/topology", serverAddr))
	if err != nil {
		logger.Error("failed to get topology", "error", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("failed to read response body", "error", err)
		return err
	}

	var topology Topology
	if err := json.Unmarshal(body, &topology); err != nil {
		logger.Error("failed to unmarshal topology", "error", err)
		return err
	}

	ethAddress := topology.Self.EthereumAddress

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
