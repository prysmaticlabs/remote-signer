package keyvault

import "github.com/prysmaticlabs/remote-signer/keyvault/deterministic"

var _ = Store(&deterministic.Store{})
