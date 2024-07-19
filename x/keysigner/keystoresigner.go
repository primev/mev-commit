package keysigner

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"runtime"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type KeystoreSigner struct {
	keystore *keystore.KeyStore
	password string
	account  accounts.Account
}

func NewKeystoreSigner(path, password string) (*KeystoreSigner, error) {
	// lightscripts are using 4MB memory and taking approximately 100ms CPU time on a modern processor to decrypt
	keystore := keystore.NewKeyStore(path, keystore.LightScryptN, keystore.LightScryptP)
	ksAccounts := keystore.Accounts()

	var account accounts.Account
	if len(ksAccounts) == 0 {
		var err error
		account, err = keystore.NewAccount(password)
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
	} else {
		account = ksAccounts[0]
	}

	if err := keystore.Unlock(account, password); err != nil {
		return nil, err
	}

	return &KeystoreSigner{
		keystore: keystore,
		password: password,
		account:  account,
	}, nil
}

func (kss *KeystoreSigner) SignHash(hash []byte) ([]byte, error) {
	return kss.keystore.SignHash(kss.account, hash)
}

func (kss *KeystoreSigner) SignTx(tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return kss.keystore.SignTx(kss.account, tx, chainID)
}

func (kss *KeystoreSigner) GetAddress() common.Address {
	return kss.account.Address
}

func (kss *KeystoreSigner) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return extractPrivateKey(kss.account.URL.Path, kss.password)
}

func (kss *KeystoreSigner) ZeroPrivateKey(key *ecdsa.PrivateKey) {
	b := key.D.Bits()
	for i := range b {
		b[i] = 0
	}
	// Force garbage collection to remove the key from memory
	runtime.GC()
}

func (kss *KeystoreSigner) String() string {
	return kss.account.URL.String()
}

func (kss *KeystoreSigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyStoreTransactorWithChainID(kss.keystore, kss.account, chainID)
	if err != nil {
		return nil, err
	}

	opts.GasLimit = 1_000_000
	opts.GasTipCap = big.NewInt(1_000_000_000)
	opts.GasFeeCap = big.NewInt(2_000_000_000)
	return opts, nil
}

func (kss *KeystoreSigner) GetAuthWithCtx(ctx context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	opts, err := kss.GetAuth(chainID)
	if err != nil {
		return nil, err
	}

	opts.Context = ctx
	return opts, nil
}
