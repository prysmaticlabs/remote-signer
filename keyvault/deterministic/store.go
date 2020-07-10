package deterministic

import (
	"context"

	"github.com/prysmaticlabs/prysm/shared/bls"
)

type Store struct{}

func (s *Store) GetSecretKey(context.Context, bls.PublicKey) (bls.SecretKey, error) {
	panic("implement me")
}

func (s *Store) GetPublicKeys(context.Context) ([]bls.PublicKey, error) {
	panic("implement me")
}
