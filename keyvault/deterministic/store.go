/*
Package deterministic defines an implementation of an insecure,
purely deterministic keyvault used by this remote signer server
for demonstration purposes.

WARN: This is a toy implementation that returns deterministic BLS12-381
secret+public key pairs, NOT MEANT FOR PRODUCTION.
*/
package deterministic

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/runtime/interop"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "deterministic-keyvault")

// Store defines a deterministic keyvault, written for demonstrative purposes.
type Store struct {
	pubKeysToSecretKeys map[[48]byte]bls.SecretKey
	pubKeys             []bls.PublicKey
}

// NewStore instantiates a deterministic keyvault using a set number of keys.
func NewStore(numKeys int) (*Store, error) {
	log.Infof("Generating %d determistic keys...", numKeys)
	secretKeys, pubKeys, err := interop.DeterministicallyGenerateKeys(
		0 /* start idx */, uint64(numKeys),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not deterministically generate %d keys", numKeys)
	}
	pubKeysToSecretKeys := make(map[[48]byte]bls.SecretKey)
	for i := 0; i < len(pubKeys); i++ {
		pubKey := bytesutil.ToBytes48(pubKeys[i].Marshal())
		pubKeysToSecretKeys[pubKey] = secretKeys[i]
	}
	log.WithField(
		"numKeys", numKeys,
	).Info("Initialized deterministic keyvault")
	s := &Store{
		pubKeys:             pubKeys,
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

// GetPublicKeys returns all available BLS12-381 public keys in the deterministic keyvault.
func (s *Store) GetPublicKeys(context.Context) ([]bls.PublicKey, error) {
	return s.pubKeys, nil
}
