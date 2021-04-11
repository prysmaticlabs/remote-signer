/*
Allows to create a Store that contains a set of public a private keys
derived from a mnemonic.

Important note: Not meant for production, mnemonic is exposed as a raw text.
*/
package mnemonic

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/sirupsen/logrus"

	"github.com/prysmaticlabs/prysm/validator/accounts/wallet"
	"github.com/prysmaticlabs/prysm/validator/keymanager/derived"
)

var log = logrus.WithField("prefix", "mnemonic-keyvault")

// Store defines a mnemonic keyvault, written for demonstrative purposes.
type Store struct {
	pubKeysToSecretKeys map[[48]byte]bls.SecretKey
	pubKeys             []bls.PublicKey
}

// NewStore instantiates a mnemonic keyvault using a set number of keys.
func NewStore(
	mnemonicPhrase string,
	mnemonicPassword string,
	startIndex int,
	numKeys int) (*Store, error) {

	log.Infof("Generating keys from mnemonic")

	ctx := context.Background()

	// Create empty wallet
	wallet := &wallet.Wallet{}
	km, err := derived.NewKeymanager(ctx, &derived.SetupConfig{
		Wallet: wallet,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a new keymanager")
	}

	// If a startIndex is provided, keys in range [0, startIndex] are also generated but not used.
	// This is done for convenience as hoisting the code from RecoverAccountsFromMnemonic can't
	// be done due to Go package constraints.
	err = km.RecoverAccountsFromMnemonic(ctx, mnemonicPhrase, mnemonicPassword, startIndex+numKeys)

	publicKeys, err := km.FetchValidatingPublicKeys(ctx)
	privateKeys, err := km.FetchValidatingPrivateKeys(ctx)

	mnemonicPubKeys := make([]bls.PublicKey, numKeys)
	mnemonicPrivKeys := make([]bls.SecretKey, numKeys)
	pubKeysToSecretKeys := make(map[[48]byte]bls.SecretKey)

	// Copy only the keys that we are interested in
	for i := startIndex; i < (numKeys + startIndex); i++ {
		log.Infof("Key from mnemonic at index %d, %x", i, publicKeys[i])
		log.Infof("Private Key from mnemonic %d, %x", i, privateKeys[i])

		blsPrivate, err := bls.SecretKeyFromBytes(privateKeys[i][:])
		if err != nil {
			return nil, errors.Wrapf(err, "could not create bls secret key at index %d from raw bytes", i)
		}

		exportIndex := i - startIndex

		mnemonicPrivKeys[exportIndex] = blsPrivate
		mnemonicPubKeys[exportIndex] = blsPrivate.PublicKey()

		pubKeysToSecretKeys[publicKeys[i]] = mnemonicPrivKeys[exportIndex]
	}

	s := &Store{
		pubKeys:             mnemonicPubKeys,
		pubKeysToSecretKeys: pubKeysToSecretKeys,
	}

	return s, nil
}

// GetSecretKey returns the corresponding secret key for a BLS12-381 public key.
func (s *Store) GetSecretKey(ctx context.Context, pubKey bls.PublicKey) (bls.SecretKey, error) {
	key := bytesutil.ToBytes48(pubKey.Marshal())
	secretKey, ok := s.pubKeysToSecretKeys[key]
	if !ok {
		return nil, fmt.Errorf("could not find secret key for public key %#x", key)
	}
	return secretKey, nil
}

// GetPublicKeys returns all available BLS12-381 public keys in the mnemonic keyvault.
func (s *Store) GetPublicKeys(context.Context) ([]bls.PublicKey, error) {
	return s.pubKeys, nil
}
