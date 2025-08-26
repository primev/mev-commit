package setcode_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/primev/mev-commit/p2p/pkg/setcode"
	"github.com/primev/mev-commit/x/keysigner"
)

func TestSetCode(t *testing.T) {

	priv, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	ks, err := keysigner.NewPrivateKeySignerFromHex(hex.EncodeToString(crypto.FromECDSA(priv)))
	if err != nil {
		t.Fatalf("failed to create keysigner: %v", err)
	}
	sender := ks.GetAddress()
	genesisAlloc := types.GenesisAlloc{
		sender: {Balance: big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(100))}, // 100 ETH
	}
	sim := simulated.NewBackend(genesisAlloc)
	defer func() {
		if err := sim.Close(); err != nil {
			t.Fatalf("failed to close simulated backend: %v", err)
		}
	}()

	tx := types.NewTransaction(0, sender, big.NewInt(1e18), 21000, big.NewInt(1e9), nil)
	tx, err = ks.SignTx(tx, big.NewInt(1337))
	if err != nil {
		t.Fatalf("failed to sign tx: %v", err)
	}
	err = sim.Client().SendTransaction(context.Background(), tx)
	if err != nil {
		t.Fatalf("failed to send tx: %v", err)
	}
	nonce, err := sim.Client().PendingNonceAt(context.Background(), sender)
	if err != nil {
		t.Fatalf("failed to get pending nonce: %v", err)
	}
	if nonce != 1 {
		t.Fatalf("nonce is not incremented: %v", nonce)
	}

	_ = sim.Commit()

	code, err := sim.Client().CodeAt(context.Background(), sender, nil)
	if err != nil {
		t.Fatalf("failed to get code: %v", err)
	}
	if len(code) != 0 {
		t.Fatalf("code is not empty before setcode: %v", code)
	}

	contractAddr, err := deployMinimalContract(
		context.Background(),
		sim.Client(),
		ks,
		big.NewInt(1337),
	)
	if err != nil {
		t.Fatalf("failed to deploy minimal contract: %v", err)
	}

	nonce, err = sim.Client().PendingNonceAt(context.Background(), sender)
	if err != nil {
		t.Fatalf("failed to get pending nonce: %v", err)
	}
	if nonce != 2 {
		t.Fatalf("nonce is not incremented: %v", nonce)
	}

	_ = sim.Commit()

	testLogger := newTestLogger(t, os.Stdout)
	setCodeHelper := setcode.NewSetCodeHelper(testLogger, ks, sim.Client(), big.NewInt(1337))

	opts, err := ks.GetAuth(big.NewInt(1337))
	if err != nil {
		t.Fatalf("failed to get auth: %v", err)
	}
	opts.GasFeeCap = big.NewInt(1e9)
	opts.GasTipCap = big.NewInt(1e9)
	opts.GasLimit = 2000000

	tx, err = setCodeHelper.SetCode(context.Background(), opts, contractAddr)
	if err != nil {
		t.Fatalf("failed to set code: %v", err)
	}

	_ = sim.Commit()

	receipt, err := bind.WaitMined(context.Background(), sim.Client(), tx)
	if err != nil {
		t.Fatalf("failed to wait for receipt: %v", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("tx failed: %v", receipt.Status)
	}

	nonce, err = sim.Client().PendingNonceAt(context.Background(), sender)
	if err != nil {
		t.Fatalf("failed to get pending nonce: %v", err)
	}
	// Nonce incremented twice, the setcode tx and the auth itself
	if nonce != 4 {
		t.Fatalf("nonce is not incremented: %v", nonce)
	}

	code, err = sim.Client().CodeAt(context.Background(), sender, nil)
	if err != nil {
		t.Fatalf("failed to get code: %v", err)
	}
	if len(code) == 0 {
		t.Fatalf("code is empty")
	}

	codehash := crypto.Keccak256Hash(code)
	expectedCodehash := crypto.Keccak256Hash(common.FromHex("0xef0100"), contractAddr.Bytes())
	if codehash != expectedCodehash {
		t.Fatalf("codehash is not correct. Actual: %v, Expected: %v", codehash, expectedCodehash)
	}

	zeroAddr := common.Address{}
	tx, err = setCodeHelper.SetCode(context.Background(), opts, zeroAddr)
	if err != nil {
		t.Fatalf("failed to set code: %v", err)
	}

	_ = sim.Commit()

	receipt, err = bind.WaitMined(context.Background(), sim.Client(), tx)
	if err != nil {
		t.Fatalf("failed to wait for receipt: %v", err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatalf("tx failed: %v", receipt.Status)
	}

	code, err = sim.Client().CodeAt(context.Background(), sender, nil)
	if err != nil {
		t.Fatalf("failed to get code: %v", err)
	}
	if len(code) != 0 {
		t.Fatalf("code is not empty after setcode to zero address: %v", code)
	}
}

func newTestLogger(t *testing.T, w io.Writer) *slog.Logger {
	t.Helper()
	testLogger := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(testLogger)
}

func deployMinimalContract(
	ctx context.Context,
	backend bind.ContractBackend,
	ks keysigner.KeySigner,
	chainID *big.Int,
) (common.Address, error) {
	from := ks.GetAddress()
	nonce, err := backend.PendingNonceAt(ctx, from)
	if err != nil {
		return common.Address{}, fmt.Errorf("get nonce: %w", err)
	}
	initcode := common.FromHex("0x6001600c60003960016000f300")
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(1e9),
		GasFeeCap: big.NewInt(1e9),
		Gas:       200_000,
		To:        nil,
		Value:     big.NewInt(0),
		Data:      initcode,
	})
	signed, err := ks.SignTx(tx, chainID)
	if err != nil {
		return common.Address{}, fmt.Errorf("sign tx: %w", err)
	}
	if err := backend.SendTransaction(ctx, signed); err != nil {
		return common.Address{}, fmt.Errorf("send tx: %w", err)
	}
	contractAddr := crypto.CreateAddress(from, nonce)
	return contractAddr, nil
}
