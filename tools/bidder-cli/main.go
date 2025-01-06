package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	pb "github.com/primev/mev-commit/p2p/gen/go/bidderapi/v1"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const basefeeWiggleMultiplier = 2

var (
	optionRPCURL = &cli.StringFlag{
		Name:     "rpc-url",
		Usage:    "URL of the gRPC server",
		Required: true,
		EnvVars:  []string{"BIDDER_CLI_RPC_URL"},
	}
	optionL1RPCURL = &cli.StringFlag{
		Name:    "l1-rpc-url",
		Usage:   "URL of the L1 RPC",
		EnvVars: []string{"BIDDER_CLI_L1_RPC_URL"},
		Value:   "https://ethereum-holesky-rpc.publicnode.com",
	}
	optionVerbose = &cli.BoolFlag{
		Name:    "verbose",
		Usage:   "enable verbose logging",
		EnvVars: []string{"BIDDER_CLI_VERBOSE"},
		Value:   false,
	}
	optionTxHashes = &cli.StringSliceFlag{
		Name:     "tx-hashes",
		Usage:    "transaction hashes",
		Required: true,
		Action: func(c *cli.Context, txHashes []string) error {
			if len(txHashes) == 0 {
				return fmt.Errorf("no transaction hashes provided")
			}
			for _, txHash := range txHashes {
				if len(txHash) != 66 && len(txHash) != 64 {
					return fmt.Errorf("invalid transaction hash: %s", txHash)
				}
			}
			return nil
		},
	}
	optionBidAmount = &cli.StringFlag{
		Name:     "bid-amount",
		Usage:    "bid amount",
		Required: true,
		Action: func(c *cli.Context, amount string) error {
			_, ok := big.NewInt(0).SetString(amount, 10)
			if !ok {
				return fmt.Errorf("invalid bid amount: %s", amount)
			}
			return nil
		},
	}
	optionBlockNumber = &cli.Int64Flag{
		Name:  "block-number",
		Usage: "block number",
	}
	optionDecayDuration = &cli.StringFlag{
		Name:  "decay-duration",
		Usage: "decay duration",
		Action: func(c *cli.Context, duration string) error {
			_, err := time.ParseDuration(duration)
			if err != nil {
				return fmt.Errorf("invalid decay duration: %w", err)
			}
			return nil
		},
		Value: "10m",
	}
	optionRevertingTxHashes = &cli.StringSliceFlag{
		Name:  "reverting-tx-hashes",
		Usage: "reverting transaction hashes",
		Action: func(c *cli.Context, txHashes []string) error {
			for _, txHash := range txHashes {
				if len(txHash) != 66 && len(txHash) != 64 {
					return fmt.Errorf("invalid transaction hash: %s", txHash)
				}
			}
			return nil
		},
	}
	optionFrom = &cli.StringFlag{
		Name:     "from",
		Usage:    "private key of the account sending the transaction",
		Required: true,
	}
	optionTo = &cli.StringFlag{
		Name:     "to",
		Usage:    "to address",
		Required: true,
		Action: func(c *cli.Context, to string) error {
			if !common.IsHexAddress(to) {
				return fmt.Errorf("invalid to address: %s", to)
			}
			return nil
		},
	}
	optionAmount = &cli.StringFlag{
		Name:     "value",
		Usage:    "amount of the transaction",
		Required: true,
		Action: func(c *cli.Context, value string) error {
			_, ok := big.NewInt(0).SetString(value, 10)
			if !ok {
				return fmt.Errorf("invalid value: %s", value)
			}
			return nil
		},
	}
	optionTipBoost = &cli.Int64Flag{
		Name:  "tip-boost",
		Usage: "tip boost %",
		Value: 10,
	}
)

func prettyPrintMsg(msg proto.Message) string {
	marshaler := protojson.MarshalOptions{
		Multiline: true, // Enables pretty-printing with indentation
		Indent:    "  ", // Customize the indentation (2 spaces here)
	}

	jsonBytes, err := marshaler.Marshal(msg)
	if err != nil {
		return fmt.Sprintf("failed to marshal proto message: %v", err)
	}
	return string(jsonBytes)
}

func main() {
	app := cli.NewApp()
	app.Name = "bidder-cli"
	app.Usage = "A CLI tool for interacting with a gRPC bidder server"
	app.Version = "1.0.0"

	app.Commands = []*cli.Command{
		{
			Name:  "send-tx-hash",
			Usage: "Send a bid to the gRPC bidder server",
			Flags: []cli.Flag{
				optionRPCURL,
				optionL1RPCURL,
				optionVerbose,
				optionTxHashes,
				optionBidAmount,
				optionBlockNumber,
				optionDecayDuration,
				optionRevertingTxHashes,
			},
			Action: func(c *cli.Context) error {
				decay, err := time.ParseDuration(c.String(optionDecayDuration.Name))
				if err != nil {
					return fmt.Errorf("invalid decay duration: %w", err)
				}

				creds := insecure.NewCredentials()
				conn, err := grpc.NewClient(
					c.String(optionRPCURL.Name),
					grpc.WithTransportCredentials(creds),
				)
				if err != nil {
					return err
				}
				defer conn.Close()

				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "connected to %s\n", c.String(optionRPCURL.Name))
				}

				client := pb.NewBidderClient(conn)

				blkNum := c.Int64(optionBlockNumber.Name)
				if blkNum == 0 {
					client, err := ethclient.Dial(c.String(optionL1RPCURL.Name))
					if err != nil {
						return fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
					}

					bNo, err := client.BlockNumber(c.Context)
					if err != nil {
						return fmt.Errorf("failed to get block number: %w", err)
					}
					blkNum = int64(bNo) + 1
					if c.Bool(optionVerbose.Name) {
						fmt.Fprintf(app.Writer, "using latest block number: %d\n", blkNum)
					}
				}

				bid := &pb.Bid{
					TxHashes:            c.StringSlice(optionTxHashes.Name),
					Amount:              c.String(optionBidAmount.Name),
					BlockNumber:         blkNum,
					DecayStartTimestamp: time.Now().UnixMilli(),
					DecayEndTimestamp:   time.Now().Add(decay).UnixMilli(),
					RevertingTxHashes:   c.StringSlice(optionRevertingTxHashes.Name),
				}

				stream, err := client.SendBid(c.Context, bid)
				if err != nil {
					return err
				}
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintln(app.Writer, "sent bid", "bid", prettyPrintMsg(bid))
				}

				for {
					preConfirmation, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						return fmt.Errorf("failed to receive preconfirmation: %w", err)
					}
					fmt.Fprintln(app.Writer, prettyPrintMsg(preConfirmation))
				}
			},
		},
		{
			Name:  "send-tx",
			Usage: "Send a transaction to the gRPC bidder server",
			Flags: []cli.Flag{
				optionRPCURL,
				optionL1RPCURL,
				optionVerbose,
				optionFrom,
				optionTo,
				optionAmount,
				optionBlockNumber,
				optionDecayDuration,
				optionBidAmount,
				optionTipBoost,
			},
			Action: func(c *cli.Context) error {
				privateKey, err := crypto.HexToECDSA(c.String(optionFrom.Name))
				if err != nil {
					return fmt.Errorf("invalid private key: %w", err)
				}

				client, err := ethclient.Dial(c.String(optionL1RPCURL.Name))
				if err != nil {
					return fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
				}

				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "connected to %s\n", c.String(optionL1RPCURL.Name))
				}

				// Get the public address of the sender
				publicKey := privateKey.Public()
				publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
				if !ok {
					return fmt.Errorf("error casting public key to ECDSA")
				}
				fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "from address: %s\n", fromAddress.Hex())
				}

				// Define transaction parameters
				nonce, err := client.PendingNonceAt(c.Context, fromAddress)
				if err != nil {
					return fmt.Errorf("failed getting account nonce: %w", err)
				}
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "nonce: %d\n", nonce)
				}

				gasLimit := uint64(21000) // basic gas limit for Ether transfer

				head, err := client.HeaderByNumber(c.Context, nil)
				if err != nil {
					return fmt.Errorf("failed getting the latest block: %w", err)
				}

				decay, err := time.ParseDuration(c.String(optionDecayDuration.Name))
				if err != nil {
					return fmt.Errorf("invalid decay duration: %w", err)
				}

				tip, err := client.SuggestGasTipCap(c.Context)
				if err != nil {
					return fmt.Errorf("failed getting gas tip cap: %w", err)
				}

				// Add boost % to tip
				tip = new(big.Int).Div(new(big.Int).Mul(big.NewInt(c.Int64(optionTipBoost.Name)+100), tip), big.NewInt(100))
				feeCap := new(big.Int).Add(
					tip,
					new(big.Int).Mul(head.BaseFee, big.NewInt(basefeeWiggleMultiplier)),
				)
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "tip: %s\n", tip.String())
					fmt.Fprintf(app.Writer, "fee cap: %s\n", feeCap.String())
				}

				value, ok := big.NewInt(0).SetString(c.String(optionAmount.Name), 10)
				if !ok {
					return fmt.Errorf("invalid amount: %s", c.String(optionAmount.Name))
				}

				to := common.HexToAddress(c.String(optionTo.Name))
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "to address: %s\n", to.Hex())
				}

				txData := &types.DynamicFeeTx{
					To:        &to,
					Nonce:     nonce,
					GasFeeCap: feeCap,
					GasTipCap: tip,
					Gas:       gasLimit,
					Value:     value,
				}

				tx := types.NewTx(txData)

				chainID, err := client.ChainID(c.Context)
				if err != nil {
					return fmt.Errorf("failed to get chain ID: %w", err)
				}

				// Sign the transaction
				signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
				if err != nil {
					return fmt.Errorf("failed to sign transaction payload: %w", err)
				}

				// Encode the transaction to a hex string
				txHex, err := signedTx.MarshalBinary()
				if err != nil {
					return fmt.Errorf("failed to encode signed transaction payload: %w", err)
				}
				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "signed transaction: %s\n", hex.EncodeToString(txHex))
				}

				creds := insecure.NewCredentials()
				conn, err := grpc.NewClient(c.String(optionRPCURL.Name), grpc.WithTransportCredentials(creds))
				if err != nil {
					return err
				}
				defer conn.Close()

				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "connected to bidder node %s\n", c.String(optionRPCURL.Name))
				}

				rpcClient := pb.NewBidderClient(conn)

				blkNum := c.Int64(optionBlockNumber.Name)
				if blkNum == 0 {
					bNo, err := client.BlockNumber(c.Context)
					if err != nil {
						return fmt.Errorf("failed to get block number: %w", err)
					}

					blkNum = int64(bNo) + 1
					if c.Bool(optionVerbose.Name) {
						fmt.Fprintf(app.Writer, "using latest block number: %d\n", blkNum)
					}
				}

				bid := &pb.Bid{
					RawTransactions:     []string{hex.EncodeToString(txHex)},
					Amount:              c.String(optionBidAmount.Name),
					BlockNumber:         blkNum,
					DecayStartTimestamp: time.Now().UnixMilli(),
					DecayEndTimestamp:   time.Now().Add(decay).UnixMilli(),
				}

				ctx := context.Background()
				stream, err := rpcClient.SendBid(ctx, bid)
				if err != nil {
					return err
				}

				if c.Bool(optionVerbose.Name) {
					fmt.Fprintf(app.Writer, "sent bid: %s\n", prettyPrintMsg(bid))
				}

				for {
					preConfirmation, err := stream.Recv()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						return fmt.Errorf("failed to receive preconfirmation: %w", err)
					}
					fmt.Fprintln(app.Writer, prettyPrintMsg(preConfirmation))
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(app.Writer, "exited with error: %v\n", err)
	}
}
