/*
Package keyvault defines a common interface for a struct
which is capable of retrieving BLS12-381 eth2 validator keys
from a secure source, such as Hashicorp Vault or other enclave. The
kind of keyvault can be dynamically chosen at runtime, and as long as
the implementation satisfies the interface, the remote signer server
will be able to use it to successfully complete signing request.
*/
package keyvault

import (
	"context"

	"github.com/prysmaticlabs/prysm/crypto/bls"
)

// Store defines a struct which has capabilities of retrieving
// BLS12-381 eth2 secret keys and public keys from a secure source.
type Store interface {
	GetSecretKey(context.Context, bls.PublicKey) (bls.SecretKey, error)
	GetPublicKeys(context.Context) ([]bls.PublicKey, error)
}
