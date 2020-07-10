package keyvault

import (
	"context"

	"github.com/prysmaticlabs/prysm/shared/bls"
)

// Store defines a struct which has capabilities of retrieving
// BLS secrets and public keys.
type Store interface {
	GetSecretKey(context.Context, bls.PublicKey) (bls.SecretKey, error)
	GetPublicKeys(context.Context) ([]bls.PublicKey, error)
}
