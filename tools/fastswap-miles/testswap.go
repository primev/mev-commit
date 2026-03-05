package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// runTestSwap performs a single FastSwap sweep using the same code path as
// production, then exits. Useful for validating Permit2 signing, Barter
// quotes, and the FastSwap API with a small amount before production deploy.
func runTestSwap(
	keystorePath, passphrase string,
	fastswapURL, l1RPC string,
	barterURL, barterKey string,
	fundsRecipient, contractAddr string,
	testInputToken, testInputAmount string,
	maxGasGwei uint64,
) {
	signer, err := loadKeystoreFile(keystorePath, passphrase)
	if err != nil {
		log.Fatalf("failed to load keystore: %v", err)
	}
	walletAddr := signer.GetAddress()
	log.Printf("TEST-SWAP: wallet=%s token=%s amount=%s recipient=%s",
		walletAddr.Hex(), testInputToken, testInputAmount, fundsRecipient)

	client, err := ethclient.Dial(fastswapURL)
	if err != nil {
		log.Fatalf("ethclient.Dial (fastswap): %v", err)
	}

	// Separate L1 client for sending the approve tx (no mev-commit balance needed).
	l1Client, err := ethclient.Dial(l1RPC)
	if err != nil {
		log.Fatalf("ethclient.Dial (l1): %v", err)
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	ctx := context.Background()

	tokenAddr := common.HexToAddress(testInputToken)
	inputAmt, ok := new(big.Int).SetString(testInputAmount, 10)
	if !ok {
		log.Fatalf("invalid -test-input-amount: %s", testInputAmount)
	}

	// Get a Barter quote for the swap to determine minReturn.
	barterReq := barterRequest{
		Source:            tokenAddr.Hex(),
		Target:            "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // ETH
		SellAmount:        inputAmt.String(),
		Recipient:         fundsRecipient,
		Origin:            walletAddr.Hex(),
		MinReturnFraction: 0.98,
		Deadline:          fmt.Sprintf("%d", time.Now().Add(10*time.Minute).Unix()),
	}
	barterResp, err := callBarter(ctx, httpClient, barterURL, barterKey, barterReq)
	if err != nil {
		log.Fatalf("Barter quote failed: %v", err)
	}
	log.Printf("Barter quote: minReturn=%s gasLimit=%s", barterResp.MinReturn, barterResp.GasLimit)

	settlement := common.HexToAddress(contractAddr)
	recipient := common.HexToAddress(fundsRecipient)

	ethReturn, gasCost, err := submitFastSwapSweep(
		ctx, client, l1Client, httpClient, signer, walletAddr,
		tokenAddr, inputAmt,
		fastswapURL, recipient, settlement,
		barterResp, maxGasGwei,
	)
	if err != nil {
		log.Fatalf("TEST-SWAP FAILED: %v", err)
	}
	log.Printf("TEST-SWAP SUCCESS: expected ETH return=%.6f estimated gas=%.6f",
		weiToEth(ethReturn), weiToEth(gasCost))
}
