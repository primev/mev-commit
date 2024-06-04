package keysigner

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func NewKeystoreSigner(path, password string) (*PrivateKeySigner, error) {
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

	pk, err := extractPrivateKey(account.URL.Path, password)
	if err != nil {
		return nil, err
	}

	return &PrivateKeySigner{
		privKey: pk,
		path:    account.URL.Path,
	}, nil
}
