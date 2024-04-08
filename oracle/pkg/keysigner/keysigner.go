package keysigner

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KeySigner interface {
	fmt.Stringer

	GetAddress() common.Address
	GetAuth(chainID *big.Int) (*bind.TransactOpts, error)
	GetPrivateKey() (*ecdsa.PrivateKey, error)
}

type PrivateKeySigner struct {
	path    string
	privKey *ecdsa.PrivateKey
}

func NewPrivateKeySigner(path string) (*PrivateKeySigner, error) {
	privKeyFile, err := resolveFilePath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key file path: %w", err)
	}

	if err := createKeyIfNotExists(privKeyFile); err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}

	privKey, err := crypto.LoadECDSA(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key from file '%s': %w", privKeyFile, err)
	}

	return &PrivateKeySigner{
		path:    privKeyFile,
		privKey: privKey,
	}, nil
}

func (pks *PrivateKeySigner) GetAddress() common.Address {
	return crypto.PubkeyToAddress(pks.privKey.PublicKey)
}

func (pks *PrivateKeySigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	return bind.NewKeyedTransactorWithChainID(pks.privKey, chainID)
}

func (pks *PrivateKeySigner) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return pks.privKey, nil
}

func (pks *PrivateKeySigner) String() string {
	return pks.path
}

type KeystoreSigner struct {
	keystore *keystore.KeyStore
	password string
	account  accounts.Account
}

func NewKeystoreSigner(path, password string) (*KeystoreSigner, error) {
	ks := keystore.NewKeyStore(path, keystore.LightScryptN, keystore.LightScryptP)
	ksAccounts := ks.Accounts()

	var account accounts.Account
	if len(ksAccounts) == 0 {
		var err error
		account, err = ks.NewAccount(password)
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
	} else {
		account = ksAccounts[0]
	}

	return &KeystoreSigner{
		keystore: ks,
		password: password,
		account:  account,
	}, nil
}

func (kss *KeystoreSigner) GetAddress() common.Address {
	return kss.account.Address
}

func (kss *KeystoreSigner) GetAuth(chainID *big.Int) (*bind.TransactOpts, error) {
	if err := kss.keystore.Unlock(kss.account, kss.password); err != nil {
		return nil, err
	}

	return bind.NewKeyStoreTransactorWithChainID(kss.keystore, kss.account, chainID)
}

func (kss *KeystoreSigner) GetPrivateKey() (*ecdsa.PrivateKey, error) {
	return extractPrivateKey(kss.account.URL.Path, kss.password)
}

func (kss *KeystoreSigner) String() string {
	return kss.account.URL.String()
}

func resolveFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}

func createKeyIfNotExists(path string) error {
	// check if key already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	// check if parent directory exists
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		// create parent directory
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}

	return crypto.SaveECDSA(path, privKey)
}

func extractPrivateKey(keystoreFile, passphrase string) (*ecdsa.PrivateKey, error) {
	keyjson, err := os.ReadFile(keystoreFile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return nil, err
	}

	// Overwrite the keyjson slice with zeros to wipe the sensitive data from memory.
	// This is a security measure to reduce the risk of the encrypted key being extracted from memory.
	for i := range keyjson {
		keyjson[i] = 0
	}

	return key.PrivateKey, nil
}
